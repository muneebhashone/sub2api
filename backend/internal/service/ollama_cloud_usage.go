package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/google/uuid"
	"golang.org/x/net/http/httpguts"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/singleflight"
)

const (
	OllamaCloudUsageSessionExtraKey     = "ollama_cloud_usage_session"
	OllamaCloudUsageAutoRefreshExtraKey = "ollama_cloud_usage_auto_refresh"
	OllamaCloudUsageSnapshotExtraKey    = "ollama_cloud_usage_snapshot"

	ollamaCloudUsageSettingsURL            = "https://ollama.com/settings"
	ollamaCloudUsageDefaultIntervalMinutes = 60
	ollamaCloudUsageMinIntervalMinutes     = 15
	ollamaCloudUsageMaxIntervalMinutes     = 24 * 60
	ollamaCloudUsageDefaultDebounceMinutes = 1
	ollamaCloudUsageMinDebounceMinutes     = 1
	ollamaCloudUsageMaxDebounceMinutes     = 60
	ollamaCloudUsageCycleInterval          = time.Minute
	ollamaCloudUsageManualRefreshInterval  = 30 * time.Second
	ollamaCloudUsageRequestTimeout         = 15 * time.Second
	ollamaCloudUsageMaxBodyBytes           = 512 * 1024
	ollamaCloudUsageMaxSessionBytes        = 16 * 1024
	ollamaCloudUsageMaxPerCycle            = 20
	ollamaCloudUsageConcurrency            = 4
	ollamaCloudUsageMaxDelay               = 24 * time.Hour
	ollamaCloudUsageLeaderLockKey          = "ollama:cloud:usage:leader"
	ollamaCloudUsageLeaderLockTTL          = 2 * time.Minute
)

var (
	ErrOllamaCloudUsageUnavailable = infraerrors.ServiceUnavailable(
		"OLLAMA_CLOUD_USAGE_UNAVAILABLE", "Ollama Cloud usage is unavailable",
	)
	ErrOllamaCloudUsageAccountInvalid = infraerrors.BadRequest(
		"OLLAMA_CLOUD_USAGE_ACCOUNT_INVALID", "account must be an OpenAI or Anthropic API key account using https://ollama.com",
	)
	ErrOllamaCloudUsageSessionRequired = infraerrors.BadRequest(
		"OLLAMA_CLOUD_USAGE_SESSION_REQUIRED", "an Ollama web session must be configured first",
	)
	ErrOllamaCloudUsageEncryptionKey = infraerrors.BadRequest(
		"OLLAMA_CLOUD_USAGE_ENCRYPTION_KEY_NOT_CONFIGURED", "cannot store an Ollama web session without a fixed TOTP_ENCRYPTION_KEY",
	)
	ErrOllamaCloudUsageIdentityChanged = infraerrors.Conflict(
		"OLLAMA_CLOUD_USAGE_IDENTITY_CHANGED", "account identity or Ollama web session changed during refresh; retry",
	)
	ErrOllamaCloudUsageRefreshRateLimited = infraerrors.TooManyRequests(
		"OLLAMA_CLOUD_USAGE_REFRESH_RATE_LIMITED", "Ollama Cloud usage can be refreshed manually once every 30 seconds",
	)
	errOllamaCloudUsageUnauthorizedHTML = errors.New("settings HTML is a sign-in page")
)

const (
	OllamaCloudUsageStatusOK           = "ok"
	OllamaCloudUsageStatusUnauthorized = "unauthorized"
	OllamaCloudUsageStatusFailed       = "failed"
)

// OllamaCloudUsageSettings controls the opt-in request-driven refresh runner.
//
// IntervalMinutes is the max-wait bound: when model requests keep arriving and
// the trailing debounce keeps sliding, a refresh is forced after this long.
// DebounceMinutes is the quiet period after the latest request in a group.
type OllamaCloudUsageSettings struct {
	Enabled         bool `json:"enabled"`
	IntervalMinutes int  `json:"interval_minutes"` // max wait while requests continue
	DebounceMinutes int  `json:"debounce_minutes"` // trailing quiet period after last request
placeholder

// OllamaCloudUsageWindow is a narrow, sanitized view of one official usage window.
type OllamaCloudUsageWindow struct {
	UsedPercent float64    `json:"used_percent"`
	ResetAt     *time.Time `json:"reset_at,omitempty"`
	ResetText   string     `json:"reset_text,omitempty"`
placeholder

// OllamaCloudUsageModelWindow identifies the official window for a model count.
type OllamaCloudUsageModelWindow string

const (
	OllamaCloudUsageModelWindowFiveHour OllamaCloudUsageModelWindow = "five_hour"
	OllamaCloudUsageModelWindowSevenDay OllamaCloudUsageModelWindow = "seven_day"
)

// OllamaCloudUsageModel is the window-scoped model/request pair exposed by Ollama's usage DOM.
type OllamaCloudUsageModel struct {
	Model    string                      `json:"model"`
	Window   OllamaCloudUsageModelWindow `json:"window"`
	Requests int64                       `json:"requests"`
placeholder

// OllamaCloudUsageData intentionally excludes raw HTML and browser-session data.
type OllamaCloudUsageData struct {
	Plan     string                  `json:"plan,omitempty"`
	FiveHour *OllamaCloudUsageWindow `json:"five_hour,omitempty"`
	SevenDay *OllamaCloudUsageWindow `json:"seven_day,omitempty"`
	Balance  string                  `json:"balance,omitempty"`
	Models   []OllamaCloudUsageModel `json:"models,omitempty"`
placeholder

// OllamaCloudUsageSnapshot is the only usage observation persisted in account extra.
//
// NextRefreshAt remains a persisted compatibility field. For status=ok it is a
// max-wait horizon marker only; automatic success refreshes are driven by model
// request activity (group last_used_at + debounce/max-wait), not by this field
// alone. For failed/unauthorized snapshots it is the failure not-before time
// (Retry-After / exponential backoff) and is enforced as max(activityDue, NextRefreshAt).
type OllamaCloudUsageSnapshot struct {
	Status        string                `json:"status"`
	Data          *OllamaCloudUsageData `json:"data,omitempty"`
	FetchedAt     *time.Time            `json:"fetched_at,omitempty"`
	LastAttemptAt time.Time             `json:"last_attempt_at"`
	NextRefreshAt time.Time             `json:"next_refresh_at"`
	FailureCount  int                   `json:"failure_count,omitempty"`
	HTTPStatus    int                   `json:"http_status,omitempty"`
	LastError     string                `json:"last_error,omitempty"`
placeholder

// OllamaCloudUsageState is the dedicated DTO exposed to administrators.
type OllamaCloudUsageState struct {
	AccountID               int64                     `json:"account_id"`
	Eligible                bool                      `json:"eligible"`
	Configured              bool                      `json:"configured"`
	AutoRefreshEnabled      bool                      `json:"auto_refresh_enabled"`
	EncryptionKeyConfigured bool                      `json:"encryption_key_configured"`
	Snapshot                *OllamaCloudUsageSnapshot `json:"snapshot,omitempty"`
placeholder

type ollamaCloudUsageRepository interface {
	ListOllamaCloudUsageGroupAccounts(context.Context, []*Account) ([]Account, error)
	SaveOllamaCloudUsageSession(context.Context, *Account, string, bool) error
	DeleteOllamaCloudUsageSession(context.Context, *Account) error
	SetOllamaCloudUsageAutoRefresh(context.Context, *Account, bool) error
	UpdateOllamaCloudUsageSnapshot(context.Context, *Account, *OllamaCloudUsageSnapshot) error
	DisableOllamaCloudUsageAutoRefresh(context.Context, *Account) error
	ListDueOllamaCloudUsageAccounts(context.Context, time.Time, int) ([]Account, error)
placeholder

// GetOllamaCloudUsageSettings returns fail-safe defaults when the setting is absent.
func (s *SettingService) GetOllamaCloudUsageSettings(ctx context.Context) (*OllamaCloudUsageSettings, error) {
	defaults := defaultOllamaCloudUsageSettings()
	if s == nil || s.settingRepo == nil {
		return defaults, nil
placeholder
	raw, err := s.settingRepo.GetValue(ctx, SettingKeyOllamaCloudUsageSettings)
	if err != nil {
		if errors.Is(err, ErrSettingNotFound) {
			return defaults, nil
	placeholder
		return nil, fmt.Errorf("get Ollama Cloud usage settings: %w", err)
placeholder
	if strings.TrimSpace(raw) == "" {
		return defaults, nil
placeholder
	settings := *defaults
	if err := json.Unmarshal([]byte(raw), &settings); err != nil {
		return nil, fmt.Errorf("parse Ollama Cloud usage settings: %w", err)
placeholder
	if settings.IntervalMinutes == 0 {
		settings.IntervalMinutes = defaults.IntervalMinutes
placeholder
	if settings.DebounceMinutes == 0 {
		settings.DebounceMinutes = defaults.DebounceMinutes
placeholder
	normalizeOllamaCloudUsageSettings(&settings)
	return &settings, nil
placeholder

func (s *SettingService) SetOllamaCloudUsageSettings(ctx context.Context, settings *OllamaCloudUsageSettings) error {
	if s == nil || s.settingRepo == nil {
		return ErrOllamaCloudUsageUnavailable
placeholder
	if settings == nil {
		return infraerrors.BadRequest("INVALID_OLLAMA_CLOUD_USAGE_SETTINGS", "settings cannot be nil")
placeholder
	if settings.DebounceMinutes == 0 {
		// Legacy clients that omit debounce_minutes keep the fail-safe default.
		settings.DebounceMinutes = ollamaCloudUsageDefaultDebounceMinutes
placeholder
	if settings.IntervalMinutes < ollamaCloudUsageMinIntervalMinutes || settings.IntervalMinutes > ollamaCloudUsageMaxIntervalMinutes {
		return infraerrors.BadRequest(
			"INVALID_OLLAMA_CLOUD_USAGE_INTERVAL",
			fmt.Sprintf("interval_minutes must be between %d and %d", ollamaCloudUsageMinIntervalMinutes, ollamaCloudUsageMaxIntervalMinutes),
		)
placeholder
	if settings.DebounceMinutes < ollamaCloudUsageMinDebounceMinutes || settings.DebounceMinutes > ollamaCloudUsageMaxDebounceMinutes {
		return infraerrors.BadRequest(
			"INVALID_OLLAMA_CLOUD_USAGE_DEBOUNCE",
			fmt.Sprintf("debounce_minutes must be between %d and %d", ollamaCloudUsageMinDebounceMinutes, ollamaCloudUsageMaxDebounceMinutes),
		)
placeholder
	normalizeOllamaCloudUsageSettings(settings)
	data, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("marshal Ollama Cloud usage settings: %w", err)
placeholder
	return s.settingRepo.Set(ctx, SettingKeyOllamaCloudUsageSettings, string(data))
placeholder

func defaultOllamaCloudUsageSettings() *OllamaCloudUsageSettings {
	return &OllamaCloudUsageSettings{
		Enabled:         false,
		IntervalMinutes: ollamaCloudUsageDefaultIntervalMinutes,
		DebounceMinutes: ollamaCloudUsageDefaultDebounceMinutes,
placeholder
placeholder

func normalizeOllamaCloudUsageSettings(settings *OllamaCloudUsageSettings) {
	if settings.IntervalMinutes < ollamaCloudUsageMinIntervalMinutes {
		settings.IntervalMinutes = ollamaCloudUsageMinIntervalMinutes
placeholder
	if settings.IntervalMinutes > ollamaCloudUsageMaxIntervalMinutes {
		settings.IntervalMinutes = ollamaCloudUsageMaxIntervalMinutes
placeholder
	if settings.DebounceMinutes <= 0 {
		settings.DebounceMinutes = ollamaCloudUsageDefaultDebounceMinutes
placeholder
	if settings.DebounceMinutes < ollamaCloudUsageMinDebounceMinutes {
		settings.DebounceMinutes = ollamaCloudUsageMinDebounceMinutes
placeholder
	if settings.DebounceMinutes > ollamaCloudUsageMaxDebounceMinutes {
		settings.DebounceMinutes = ollamaCloudUsageMaxDebounceMinutes
placeholder
placeholder

func ollamaCloudUsageDurations(settings *OllamaCloudUsageSettings) (debounce, maxWait time.Duration) {
	normalized := defaultOllamaCloudUsageSettings()
	if settings != nil {
		*normalized = *settings
placeholder
	normalizeOllamaCloudUsageSettings(normalized)
	return time.Duration(normalized.DebounceMinutes) * time.Minute,
		time.Duration(normalized.IntervalMinutes) * time.Minute
placeholder

// ollamaCloudUsageIsAutoRefreshDue decides whether a configured auto-refresh
// group should fetch now. groupLastUsedAt must be MAX(last_used_at) across the
// exact api_key group so shared multi-platform accounts do not miss activity.
//
// Success: a request must be newer than fetched_at; dueAt = min(lastUsed+debounce, fetchedAt+maxWait).
// Failure: a request must be newer than last_attempt_at; activity due uses the same min formula,
// then dueAt = max(activityDue, next_refresh_at) so Retry-After / exponential backoff win.
// Missing or invalid snapshots fail open to a first fetch.
func ollamaCloudUsageIsAutoRefreshDue(
	snapshot *OllamaCloudUsageSnapshot,
	groupLastUsedAt *time.Time,
	now time.Time,
	debounce, maxWait time.Duration,
) bool {
	dueAt, ok := ollamaCloudUsageAutoRefreshDueAt(snapshot, groupLastUsedAt, debounce, maxWait)
	if !ok {
		return false
placeholder
	return !now.Before(dueAt)
placeholder

func ollamaCloudUsageAutoRefreshDueAt(
	snapshot *OllamaCloudUsageSnapshot,
	groupLastUsedAt *time.Time,
	debounce, maxWait time.Duration,
) (time.Time, bool) {
	if debounce <= 0 {
		debounce = time.Duration(ollamaCloudUsageDefaultDebounceMinutes) * time.Minute
placeholder
	if maxWait <= 0 {
		maxWait = time.Duration(ollamaCloudUsageDefaultIntervalMinutes) * time.Minute
placeholder
	if snapshot == nil {
		return time.Time{placeholder, true
placeholder
	switch snapshot.Status {
	case OllamaCloudUsageStatusOK:
		if snapshot.FetchedAt == nil || snapshot.FetchedAt.IsZero() {
			return time.Time{placeholder, true
	placeholder
		fetchedAt := snapshot.FetchedAt.UTC()
		if groupLastUsedAt == nil || !groupLastUsedAt.After(fetchedAt) {
			return time.Time{placeholder, false
	placeholder
		lastUsed := groupLastUsedAt.UTC()
		return minTime(lastUsed.Add(debounce), fetchedAt.Add(maxWait)), true
	case OllamaCloudUsageStatusFailed, OllamaCloudUsageStatusUnauthorized:
		if snapshot.LastAttemptAt.IsZero() {
			return time.Time{placeholder, true
	placeholder
		lastAttempt := snapshot.LastAttemptAt.UTC()
		if groupLastUsedAt == nil || !groupLastUsedAt.After(lastAttempt) {
			return time.Time{placeholder, false
	placeholder
		lastUsed := groupLastUsedAt.UTC()
		activityDue := minTime(lastUsed.Add(debounce), lastAttempt.Add(maxWait))
		if !snapshot.NextRefreshAt.IsZero() && snapshot.NextRefreshAt.UTC().After(activityDue) {
			return snapshot.NextRefreshAt.UTC(), true
	placeholder
		return activityDue, true
	default:
		return time.Time{placeholder, true
placeholder
placeholder

// maxOllamaCloudUsageGroupLastUsed returns the newest last_used_at among group members.
func maxOllamaCloudUsageGroupLastUsed(accounts []Account) *time.Time {
	var latest *time.Time
	for i := range accounts {
		candidate := accounts[i].LastUsedAt
		if candidate == nil || candidate.IsZero() {
			continue
	placeholder
		if latest == nil || candidate.After(*latest) {
			ts := candidate.UTC()
			latest = &ts
	placeholder
placeholder
	return latest
placeholder

// scheduleOllamaCloudUsageActivity records that an Ollama Cloud API-key account
// actually attempted an upstream model request (including 429/5xx/transport errors).
// Local auth/validation failures must not call this. DeferredService dedupes writes.
func scheduleOllamaCloudUsageActivity(deferred *DeferredService, account *Account) {
	if deferred == nil || account == nil || !IsOllamaCloudUsageAccount(account) {
		return
placeholder
	deferred.ScheduleLastUsedUpdate(account.ID)
placeholder

// OllamaCloudUsageService refreshes the official settings HTML without affecting routing state.
type OllamaCloudUsageService struct {
	accountRepo             AccountRepository
	httpUpstream            HTTPUpstream
	settingService          *SettingService
	encryptor               SecretEncryptor
	encryptionKeyConfigured bool

	parentCtx    context.Context
	parentCancel context.CancelFunc
	wg           sync.WaitGroup
	mu           sync.Mutex
	started      bool
	stopped      bool
	cycleMu      sync.Mutex
	refreshGroup singleflight.Group
	refreshSlots chan struct{placeholder
	now          func() time.Time
	lockCache    LeaderLockCache
	db           *sql.DB
	instanceID   string
placeholder

func NewOllamaCloudUsageService(
	accountRepo AccountRepository,
	httpUpstream HTTPUpstream,
	settingService *SettingService,
	encryptor SecretEncryptor,
	encryptionKeyConfigured bool,
) *OllamaCloudUsageService {
	ctx, cancel := context.WithCancel(context.Background())
	return &OllamaCloudUsageService{
		accountRepo:             accountRepo,
		httpUpstream:            httpUpstream,
		settingService:          settingService,
		encryptor:               encryptor,
		encryptionKeyConfigured: encryptionKeyConfigured,
		parentCtx:               ctx,
		parentCancel:            cancel,
		refreshSlots:            make(chan struct{placeholder, ollamaCloudUsageConcurrency),
		now:                     time.Now,
		instanceID:              uuid.NewString(),
placeholder
placeholder

func ProvideOllamaCloudUsageService(
	accountRepo AccountRepository,
	httpUpstream HTTPUpstream,
	settingService *SettingService,
	encryptor SecretEncryptor,
	cfg *config.Config,
	lockCache LeaderLockCache,
	db *sql.DB,
) *OllamaCloudUsageService {
	keyConfigured := cfg != nil && cfg.Totp.EncryptionKeyConfigured
	svc := NewOllamaCloudUsageService(accountRepo, httpUpstream, settingService, encryptor, keyConfigured)
	svc.lockCache = lockCache
	svc.db = db
	svc.Start()
	return svc
placeholder

func (s *OllamaCloudUsageService) Start() {
	if s == nil {
		return
placeholder
	s.mu.Lock()
	if s.started || s.stopped {
		s.mu.Unlock()
		return
placeholder
	s.started = true
	s.wg.Add(1)
	s.mu.Unlock()
	go s.runLoop()
placeholder

func (s *OllamaCloudUsageService) Stop() {
	if s == nil {
		return
placeholder
	s.mu.Lock()
	if s.stopped {
		s.mu.Unlock()
		return
placeholder
	s.stopped = true
	s.parentCancel()
	s.mu.Unlock()
	s.wg.Wait()
placeholder

func (s *OllamaCloudUsageService) runLoop() {
	defer s.wg.Done()
	_ = s.RunDue(s.parentCtx)
	ticker := time.NewTicker(ollamaCloudUsageCycleInterval)
	defer ticker.Stop()
	for {
		select {
		case <-s.parentCtx.Done():
			return
		case <-ticker.C:
			if err := s.RunDue(s.parentCtx); err != nil {
				logger.LegacyPrintf("service.ollama_cloud_usage", "run_due_failed: err=%v", err)
		placeholder
	placeholder
placeholder
placeholder

func (s *OllamaCloudUsageService) GetSettings(ctx context.Context) (*OllamaCloudUsageSettings, error) {
	if s == nil || s.settingService == nil {
		return defaultOllamaCloudUsageSettings(), nil
placeholder
	return s.settingService.GetOllamaCloudUsageSettings(ctx)
placeholder

func (s *OllamaCloudUsageService) UpdateSettings(ctx context.Context, settings *OllamaCloudUsageSettings) error {
	if s == nil || s.settingService == nil {
		return ErrOllamaCloudUsageUnavailable
placeholder
	return s.settingService.SetOllamaCloudUsageSettings(ctx, settings)
placeholder

func (s *OllamaCloudUsageService) GetState(ctx context.Context, accountID int64) (*OllamaCloudUsageState, error) {
	if s == nil || s.accountRepo == nil {
		return nil, ErrOllamaCloudUsageUnavailable
placeholder
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
placeholder
	if err := s.ResolveAccounts(ctx, []*Account{accountplaceholder); err != nil {
		return nil, err
placeholder
	state := OllamaCloudUsageStateFromAccount(account)
	s.EnrichState(state)
	return state, nil
placeholder

// ResolveAccounts overlays group-owned managed state onto the supplied account
// objects. The repository resolves all matching siblings in one bounded query,
// so account-list responses do not issue one query per row.
func (s *OllamaCloudUsageService) ResolveAccounts(ctx context.Context, accounts []*Account) error {
	if s == nil || s.accountRepo == nil || len(accounts) == 0 {
		return nil
placeholder
	writer, ok := s.accountRepo.(ollamaCloudUsageRepository)
	if !ok {
		return nil
placeholder
	eligible := make([]*Account, 0, len(accounts))
	for _, account := range accounts {
		if _, ok := ollamaCloudUsageGroupFingerprint(account); ok {
			eligible = append(eligible, account)
	placeholder
placeholder
	if len(eligible) == 0 {
		return nil
placeholder
	siblings, err := writer.ListOllamaCloudUsageGroupAccounts(ctx, eligible)
	if err != nil {
		return fmt.Errorf("resolve Ollama Cloud usage groups: %w", err)
placeholder
	sources := make(map[string]*Account)
	for index := range siblings {
		candidate := &siblings[index]
		fingerprint, valid := ollamaCloudUsageGroupFingerprint(candidate)
		if !valid || !ollamaCloudUsageConfigured(candidate) {
			continue
	placeholder
		current := sources[fingerprint]
		if current == nil || candidate.UpdatedAt.After(current.UpdatedAt) ||
			(candidate.UpdatedAt.Equal(current.UpdatedAt) && candidate.ID < current.ID) {
			sources[fingerprint] = candidate
	placeholder
placeholder
	resolvedSources := make(map[string]*Account, len(sources))
	for fingerprint, source := range sources {
		clone := *source
		clone.Extra = make(map[string]any, len(source.Extra))
		maps.Copy(clone.Extra, source.Extra)
		resolvedSources[fingerprint] = &clone
placeholder
	for index := range siblings {
		candidate := &siblings[index]
		fingerprint, valid := ollamaCloudUsageGroupFingerprint(candidate)
		source := resolvedSources[fingerprint]
		if !valid || source == nil || !sameOllamaCloudUsageSession(source, candidate) {
			continue
	placeholder
		candidateSnapshot := decodeOllamaCloudUsageSnapshot(candidate.Extra)
		currentSnapshot := decodeOllamaCloudUsageSnapshot(source.Extra)
		if candidateSnapshot != nil && (currentSnapshot == nil || candidateSnapshot.LastAttemptAt.After(currentSnapshot.LastAttemptAt)) {
			source.Extra[OllamaCloudUsageSnapshotExtraKey] = candidate.Extra[OllamaCloudUsageSnapshotExtraKey]
	placeholder
placeholder
	for _, account := range eligible {
		fingerprint, _ := ollamaCloudUsageGroupFingerprint(account)
		applyOllamaCloudUsageManagedExtra(account, resolvedSources[fingerprint])
placeholder
	return nil
placeholder

func sameOllamaCloudUsageSession(left, right *Account) bool {
	if left == nil || right == nil || left.Extra == nil || right.Extra == nil {
		return false
placeholder
	leftSession, leftOK := left.Extra[OllamaCloudUsageSessionExtraKey].(string)
	rightSession, rightOK := right.Extra[OllamaCloudUsageSessionExtraKey].(string)
	return leftOK && rightOK && leftSession != "" && leftSession == rightSession
placeholder

func applyOllamaCloudUsageManagedExtra(target, source *Account) {
	if target == nil {
		return
placeholder
	if target.Extra == nil {
		target.Extra = make(map[string]any)
placeholder
	for _, key := range []string{
		OllamaCloudUsageSessionExtraKey,
		OllamaCloudUsageAutoRefreshExtraKey,
		OllamaCloudUsageSnapshotExtraKey,
placeholder {
		delete(target.Extra, key)
		if source != nil && source.Extra != nil {
			if value, ok := source.Extra[key]; ok {
				target.Extra[key] = value
		placeholder
	placeholder
placeholder
placeholder

func (s *OllamaCloudUsageService) SaveSession(ctx context.Context, accountID int64, session string) (*OllamaCloudUsageState, error) {
	if s == nil || s.accountRepo == nil || s.encryptor == nil {
		return nil, ErrOllamaCloudUsageUnavailable
placeholder
	if !s.encryptionKeyConfigured {
		return nil, ErrOllamaCloudUsageEncryptionKey
placeholder
	normalized, err := normalizeOllamaCloudUsageCookie(session)
	if err != nil {
		return nil, infraerrors.BadRequest("INVALID_OLLAMA_CLOUD_USAGE_SESSION", err.Error())
placeholder
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
placeholder
	if !IsOllamaCloudUsageAccount(account) {
		return nil, ErrOllamaCloudUsageAccountInvalid
placeholder
	if err := s.ResolveAccounts(ctx, []*Account{accountplaceholder); err != nil {
		return nil, err
placeholder
	ciphertext, err := s.encryptor.Encrypt(normalized)
	if err != nil {
		return nil, fmt.Errorf("encrypt Ollama web session: %w", err)
placeholder
	writer, ok := s.accountRepo.(ollamaCloudUsageRepository)
	if !ok {
		return nil, ErrOllamaCloudUsageUnavailable
placeholder
	preserveAutoRefresh := ollamaCloudUsageConfigured(account) && ollamaCloudUsageAutoRefreshEnabled(account)
	if err := writer.SaveOllamaCloudUsageSession(ctx, account, ciphertext, preserveAutoRefresh); err != nil {
		return nil, err
placeholder
	return s.GetState(ctx, accountID)
placeholder

func (s *OllamaCloudUsageService) DeleteSession(ctx context.Context, accountID int64) (*OllamaCloudUsageState, error) {
	if s == nil || s.accountRepo == nil {
		return nil, ErrOllamaCloudUsageUnavailable
placeholder
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
placeholder
	if !IsOllamaCloudUsageAccount(account) {
		return nil, ErrOllamaCloudUsageAccountInvalid
placeholder
	if err := s.ResolveAccounts(ctx, []*Account{accountplaceholder); err != nil {
		return nil, err
placeholder
	writer, ok := s.accountRepo.(ollamaCloudUsageRepository)
	if !ok {
		return nil, ErrOllamaCloudUsageUnavailable
placeholder
	if err := writer.DeleteOllamaCloudUsageSession(ctx, account); err != nil {
		return nil, err
placeholder
	return s.GetState(ctx, accountID)
placeholder

func (s *OllamaCloudUsageService) SetAutoRefresh(ctx context.Context, accountID int64, enabled bool) (*OllamaCloudUsageState, error) {
	if s == nil || s.accountRepo == nil {
		return nil, ErrOllamaCloudUsageUnavailable
placeholder
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
placeholder
	if !IsOllamaCloudUsageAccount(account) {
		return nil, ErrOllamaCloudUsageAccountInvalid
placeholder
	if err := s.ResolveAccounts(ctx, []*Account{accountplaceholder); err != nil {
		return nil, err
placeholder
	if enabled && !ollamaCloudUsageConfigured(account) {
		return nil, ErrOllamaCloudUsageSessionRequired
placeholder
	writer, ok := s.accountRepo.(ollamaCloudUsageRepository)
	if !ok {
		return nil, ErrOllamaCloudUsageUnavailable
placeholder
	if err := writer.SetOllamaCloudUsageAutoRefresh(ctx, account, enabled); err != nil {
		return nil, err
placeholder
	return s.GetState(ctx, accountID)
placeholder

func (s *OllamaCloudUsageService) Refresh(ctx context.Context, accountID int64) (*OllamaCloudUsageState, error) {
	settings, err := s.GetSettings(ctx)
	if err != nil {
		return nil, err
placeholder
	if _, err := s.refreshAccount(ctx, accountID, settings, false); err != nil {
		return nil, err
placeholder
	return s.GetState(ctx, accountID)
placeholder

func (s *OllamaCloudUsageService) RunDue(ctx context.Context) error {
	if s == nil || s.accountRepo == nil {
		return nil
placeholder
	s.cycleMu.Lock()
	defer s.cycleMu.Unlock()
	settings, err := s.GetSettings(ctx)
	if err != nil {
		return err
placeholder
	if !settings.Enabled {
		return nil
placeholder
	release, acquired := tryAcquireSingletonLeaderLock(ctx, s.lockCache, s.db, ollamaCloudUsageLeaderLockKey, s.instanceID, ollamaCloudUsageLeaderLockTTL)
	if !acquired {
		return nil
placeholder
	defer release()

	writer, ok := s.accountRepo.(ollamaCloudUsageRepository)
	if !ok {
		return ErrOllamaCloudUsageUnavailable
placeholder
	now := s.currentTime()
	debounce, maxWait := ollamaCloudUsageDurations(settings)
	accounts, err := writer.ListDueOllamaCloudUsageAccounts(ctx, now, ollamaCloudUsageMaxPerCycle)
	if err != nil {
		return fmt.Errorf("list due Ollama Cloud usage accounts: %w", err)
placeholder
	var group errgroup.Group
	seenGroups := make(map[string]struct{placeholder, len(accounts))
	for index := range accounts {
		account := accounts[index]
		fingerprint, valid := ollamaCloudUsageGroupFingerprint(&account)
		if !valid || !account.IsActive() || !ollamaCloudUsageConfigured(&account) || !ollamaCloudUsageAutoRefreshEnabled(&account) {
			continue
	placeholder
		if _, duplicate := seenGroups[fingerprint]; duplicate {
			continue
	placeholder
		seenGroups[fingerprint] = struct{placeholder{placeholder
		snapshot := decodeOllamaCloudUsageSnapshot(account.Extra)
		// ListDue stamps Account.LastUsedAt with the api_key group MAX(last_used_at).
		if !ollamaCloudUsageIsAutoRefreshDue(snapshot, account.LastUsedAt, now, debounce, maxWait) {
			continue
	placeholder
		accountID := account.ID
		expected := account
		group.Go(func() error {
			if _, refreshErr := s.refreshAccount(ctx, accountID, settings, true); refreshErr != nil {
				if errors.Is(refreshErr, ErrOllamaCloudUsageIdentityChanged) {
					if disableErr := writer.DisableOllamaCloudUsageAutoRefresh(ctx, &expected); disableErr != nil {
						logger.LegacyPrintf("service.ollama_cloud_usage", "disable_auto_refresh_failed: account_id=%d err=%v", accountID, disableErr)
				placeholder
					return nil
			placeholder
				logger.LegacyPrintf("service.ollama_cloud_usage", "refresh_due_failed: account_id=%d err=%v", accountID, refreshErr)
		placeholder
			return nil
	placeholder)
placeholder
	return group.Wait()
placeholder

func (s *OllamaCloudUsageService) refreshAccount(ctx context.Context, accountID int64, settings *OllamaCloudUsageSettings, requireEnabled bool) (*OllamaCloudUsageSnapshot, error) {
	if s == nil || s.accountRepo == nil {
		return nil, ErrOllamaCloudUsageUnavailable
placeholder
	if settings == nil {
		settings = defaultOllamaCloudUsageSettings()
placeholder
	intervalMinutes := settings.IntervalMinutes
	debounce, maxWait := ollamaCloudUsageDurations(settings)
	anchor, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
placeholder
	key, valid := ollamaCloudUsageGroupFingerprint(anchor)
	if !valid {
		return nil, ErrOllamaCloudUsageAccountInvalid
placeholder
	value, err, _ := s.refreshGroup.Do(key, func() (any, error) {
		select {
		case s.refreshSlots <- struct{placeholder{placeholder:
			defer func() { <-s.refreshSlots placeholder()
		case <-ctx.Done():
			return nil, ctx.Err()
	placeholder
		account, loadErr := s.accountRepo.GetByID(ctx, accountID)
		if loadErr != nil {
			return nil, loadErr
	placeholder
		currentKey, currentValid := ollamaCloudUsageGroupFingerprint(account)
		if !currentValid {
			return nil, ErrOllamaCloudUsageAccountInvalid
	placeholder
		if currentKey != key {
			return nil, ErrOllamaCloudUsageIdentityChanged
	placeholder
		if err := s.ResolveAccounts(ctx, []*Account{accountplaceholder); err != nil {
			return nil, err
	placeholder
		if !ollamaCloudUsageConfigured(account) {
			return nil, ErrOllamaCloudUsageSessionRequired
	placeholder
		if !requireEnabled {
			if snapshot := decodeOllamaCloudUsageSnapshot(account.Extra); snapshot != nil && !snapshot.LastAttemptAt.IsZero() {
				retryAt := snapshot.LastAttemptAt.Add(ollamaCloudUsageManualRefreshInterval)
				if now := s.currentTime(); now.Before(retryAt) {
					remaining := retryAt.Sub(now)
					seconds := int((remaining + time.Second - 1) / time.Second)
					return nil, ErrOllamaCloudUsageRefreshRateLimited.WithMetadata(map[string]string{
						"retry_after_seconds": strconv.Itoa(seconds),
				placeholder)
			placeholder
		placeholder
	placeholder
		if requireEnabled {
			if !account.IsActive() || !ollamaCloudUsageAutoRefreshEnabled(account) {
				return nil, nil
		placeholder
			groupLastUsed := account.LastUsedAt
			if writer, ok := s.accountRepo.(ollamaCloudUsageRepository); ok {
				if siblings, listErr := writer.ListOllamaCloudUsageGroupAccounts(ctx, []*Account{accountplaceholder); listErr == nil {
					groupLastUsed = maxOllamaCloudUsageGroupLastUsed(siblings)
			placeholder
		placeholder
			if !ollamaCloudUsageIsAutoRefreshDue(decodeOllamaCloudUsageSnapshot(account.Extra), groupLastUsed, s.currentTime(), debounce, maxWait) {
				return nil, nil
		placeholder
	placeholder
		return s.refreshLoadedAccount(ctx, account, intervalMinutes)
placeholder)
	if err != nil || value == nil {
		return nil, err
placeholder
	snapshot, ok := value.(*OllamaCloudUsageSnapshot)
	if !ok {
		return nil, fmt.Errorf("invalid Ollama Cloud usage refresh result")
placeholder
	return snapshot, nil
placeholder

func (s *OllamaCloudUsageService) refreshLoadedAccount(ctx context.Context, account *Account, intervalMinutes int) (*OllamaCloudUsageSnapshot, error) {
	now := s.currentTime().UTC()
	ciphertext, _ := account.Extra[OllamaCloudUsageSessionExtraKey].(string)
	if ciphertext == "" {
		return nil, ErrOllamaCloudUsageSessionRequired
placeholder
	if !s.encryptionKeyConfigured || s.encryptor == nil {
		return nil, ErrOllamaCloudUsageEncryptionKey
placeholder
	cookie, err := s.encryptor.Decrypt(ciphertext)
	if err != nil {
		return nil, infraerrors.ServiceUnavailable("OLLAMA_CLOUD_USAGE_SESSION_DECRYPT_FAILED", "stored Ollama web session cannot be decrypted")
placeholder
	cookie, err = normalizeOllamaCloudUsageCookie(cookie)
	if err != nil {
		return nil, infraerrors.ServiceUnavailable("OLLAMA_CLOUD_USAGE_SESSION_INVALID", "stored Ollama web session is invalid")
placeholder
	if s.httpUpstream == nil {
		return nil, ErrOllamaCloudUsageUnavailable
placeholder
	proxyURL := ""
	if account.ProxyID != nil {
		if account.Proxy == nil || account.Proxy.ID != *account.ProxyID {
			return nil, ErrOllamaCloudUsageIdentityChanged
	placeholder
		proxyURL = account.Proxy.URL()
placeholder
	requestCtx, cancel := context.WithTimeout(WithHTTPUpstreamRedirectsDisabled(ctx), ollamaCloudUsageRequestTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(requestCtx, http.MethodGet, ollamaCloudUsageSettingsURL, nil)
	if err != nil || !isExactOllamaCloudSettingsURL(req.URL) {
		return nil, ErrOllamaCloudUsageUnavailable
placeholder
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Cookie", cookie)
	req.Header.Set("User-Agent", "sub2api-ollama-usage/1")
	resp, err := s.httpUpstream.Do(req, proxyURL, account.ID, account.Concurrency)
	if err != nil {
		return s.persistFailure(ctx, account, intervalMinutes, now, 0, "request_failed", 0, false)
placeholder
	if resp == nil || resp.Body == nil {
		return s.persistFailure(ctx, account, intervalMinutes, now, 0, "empty_response", 0, false)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	if resp.Request != nil && !isExactOllamaCloudSettingsURL(resp.Request.URL) {
		return s.persistFailure(ctx, account, intervalMinutes, now, resp.StatusCode, "response_host_mismatch", 0, false)
placeholder
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		return s.persistFailure(ctx, account, intervalMinutes, now, resp.StatusCode, "redirect_blocked", retryAfter(resp.Header, now), false)
placeholder
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return s.persistFailure(ctx, account, intervalMinutes, now, resp.StatusCode, "unauthorized", retryAfter(resp.Header, now), true)
placeholder
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return s.persistFailure(ctx, account, intervalMinutes, now, resp.StatusCode, "http_error", retryAfter(resp.Header, now), false)
placeholder
	body, readErr := io.ReadAll(io.LimitReader(resp.Body, ollamaCloudUsageMaxBodyBytes+1))
	if readErr != nil {
		return s.persistFailure(ctx, account, intervalMinutes, now, resp.StatusCode, "response_read_failed", 0, false)
placeholder
	if len(body) > ollamaCloudUsageMaxBodyBytes {
		return s.persistFailure(ctx, account, intervalMinutes, now, resp.StatusCode, "response_too_large", 0, false)
placeholder
	data, parseErr := parseOllamaCloudUsageHTML(body)
	if errors.Is(parseErr, errOllamaCloudUsageUnauthorizedHTML) {
		return s.persistFailure(ctx, account, intervalMinutes, now, resp.StatusCode, "unauthorized", 0, true)
placeholder
	if parseErr != nil {
		return s.persistFailure(ctx, account, intervalMinutes, now, resp.StatusCode, "invalid_html", 0, false)
placeholder
	snapshot := &OllamaCloudUsageSnapshot{
		Status:        OllamaCloudUsageStatusOK,
		Data:          data,
		FetchedAt:     &now,
		LastAttemptAt: now,
		NextRefreshAt: now.Add(nextOllamaCloudUsageDelay(intervalMinutes, 0, 0)),
		HTTPStatus:    resp.StatusCode,
placeholder
	if err := s.updateSnapshot(ctx, account, snapshot); err != nil {
		return nil, err
placeholder
	return snapshot, nil
placeholder

func (s *OllamaCloudUsageService) persistFailure(
	ctx context.Context,
	account *Account,
	intervalMinutes int,
	now time.Time,
	httpStatus int,
	reason string,
	retryAfterDuration time.Duration,
	unauthorized bool,
) (*OllamaCloudUsageSnapshot, error) {
	previous := decodeOllamaCloudUsageSnapshot(account.Extra)
	failureCount := 1
	if previous != nil {
		failureCount = previous.FailureCount + 1
placeholder
	status := OllamaCloudUsageStatusFailed
	if unauthorized {
		status = OllamaCloudUsageStatusUnauthorized
placeholder
	snapshot := &OllamaCloudUsageSnapshot{
		Status:        status,
		LastAttemptAt: now,
		NextRefreshAt: now.Add(nextOllamaCloudUsageDelay(intervalMinutes, failureCount, retryAfterDuration)),
		FailureCount:  failureCount,
		HTTPStatus:    httpStatus,
		LastError:     reason,
placeholder
	if previous != nil {
		snapshot.Data = previous.Data
		snapshot.FetchedAt = previous.FetchedAt
placeholder
	if err := s.updateSnapshot(ctx, account, snapshot); err != nil {
		return nil, err
placeholder
	return snapshot, nil
placeholder

func (s *OllamaCloudUsageService) updateSnapshot(ctx context.Context, account *Account, snapshot *OllamaCloudUsageSnapshot) error {
	writer, ok := s.accountRepo.(ollamaCloudUsageRepository)
	if !ok {
		return ErrOllamaCloudUsageUnavailable
placeholder
	return writer.UpdateOllamaCloudUsageSnapshot(ctx, account, snapshot)
placeholder

// EnrichState adds service-owned runtime configuration to an account-derived state.
func (s *OllamaCloudUsageService) EnrichState(state *OllamaCloudUsageState) {
	if state == nil {
		return
placeholder
	state.EncryptionKeyConfigured = s != nil && s.encryptionKeyConfigured
placeholder

func OllamaCloudUsageStateFromAccount(account *Account) *OllamaCloudUsageState {
	state := &OllamaCloudUsageState{placeholder
	if account == nil {
		return state
placeholder
	state.AccountID = account.ID
	state.Eligible = IsOllamaCloudUsageAccount(account)
	if !state.Eligible {
		return state
placeholder
	state.Configured = ollamaCloudUsageConfigured(account)
	state.AutoRefreshEnabled = state.Configured && ollamaCloudUsageAutoRefreshEnabled(account)
	state.Snapshot = decodeOllamaCloudUsageSnapshot(account.Extra)
	return state
placeholder

func IsOllamaCloudUsageAccount(account *Account) bool {
	if account == nil || account.Type != AccountTypeAPIKey || (account.Platform != PlatformOpenAI && account.Platform != PlatformAnthropic) {
		return false
placeholder
	baseURL, _ := account.Credentials["base_url"].(string)
	return isOllamaCloudBaseURL(baseURL)
placeholder

func isOllamaCloudBaseURL(raw string) bool {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.ContainsAny(raw, "?#") {
		return false
placeholder
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Opaque != "" || !strings.EqualFold(parsed.Scheme, "https") || parsed.User != nil || parsed.ForceQuery || parsed.RawQuery != "" || parsed.Fragment != "" || parsed.RawFragment != "" {
		return false
placeholder
	hostname := strings.ToLower(parsed.Hostname())
	if hostname != "ollama.com" && hostname != "www.ollama.com" {
		return false
placeholder
	authority := strings.ToLower(parsed.Host)
	if authority != hostname && authority != hostname+":443" {
		return false
placeholder
	if parsed.RawPath != "" {
		return false
placeholder
	return parsed.Path == "" || parsed.Path == "/v1"
placeholder

func ollamaCloudUsageIdentity(account *Account) map[string]any {
	if !IsOllamaCloudUsageAccount(account) {
		return nil
placeholder
	apiKey, ok := account.Credentials["api_key"].(string)
	if !ok || apiKey == "" {
		return nil
placeholder
	return map[string]any{"host": "ollama.com", "api_key": apiKeyplaceholder
placeholder

func ollamaCloudUsageGroupFingerprint(account *Account) (string, bool) {
	identity := ollamaCloudUsageIdentity(account)
	if identity == nil {
		return "", false
placeholder
	apiKey, _ := identity["api_key"].(string)
	sum := sha256.Sum256([]byte("ollama.com\x00" + apiKey))
	return hex.EncodeToString(sum[:]), true
placeholder

func isExactOllamaCloudSettingsURL(parsed *url.URL) bool {
	return parsed != nil && parsed.Scheme == "https" && parsed.Host == "ollama.com" && parsed.Path == "/settings" &&
		parsed.User == nil && parsed.RawQuery == "" && parsed.Fragment == "" && parsed.RawPath == ""
placeholder

func normalizeOllamaCloudUsageCookie(raw string) (string, error) {
	if len(raw) > ollamaCloudUsageMaxSessionBytes {
		return "", errors.New("session is too large")
placeholder
	raw = strings.TrimSpace(raw)
	if strings.ContainsAny(raw, "\r\n") {
		return "", errors.New("session contains invalid header characters")
placeholder
	if raw == "" {
		return "", errors.New("session cannot be empty")
placeholder
	if !httpguts.ValidHeaderFieldValue(raw) {
		return "", errors.New("session contains invalid header characters")
placeholder
	blockedAttributes := map[string]struct{placeholder{
		"domain": {placeholder, "path": {placeholder, "expires": {placeholder, "max-age": {placeholder, "samesite": {placeholder, "secure": {placeholder, "httponly": {placeholder, "partitioned": {placeholder,
placeholder
	parts := strings.Split(raw, ";")
	normalized := make([]string, 0, len(parts))
	seen := make(map[string]struct{placeholder, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		name, value, ok := strings.Cut(part, "=")
		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)
		if !ok || name == "" || value == "" || !httpguts.ValidHeaderFieldName(name) || strings.HasPrefix(name, "$") {
			return "", errors.New("session must be a Cookie header containing name=value pairs")
	placeholder
		lowerName := strings.ToLower(name)
		if _, blocked := blockedAttributes[lowerName]; blocked {
			return "", errors.New("paste a Cookie header, not a Set-Cookie value with attributes")
	placeholder
		if _, duplicate := seen[lowerName]; duplicate {
			return "", errors.New("session contains duplicate cookie names")
	placeholder
		if strings.ContainsAny(value, ";\r\n") {
			return "", errors.New("session contains an invalid cookie value")
	placeholder
		seen[lowerName] = struct{placeholder{placeholder
		if isAllowedOllamaCloudSessionCookie(name) {
			normalized = append(normalized, name+"="+value)
	placeholder
placeholder
	if len(normalized) == 0 {
		return "", errors.New("session does not contain an allowed Ollama session cookie")
placeholder
	return strings.Join(normalized, "; "), nil
placeholder

func isAllowedOllamaCloudSessionCookie(name string) bool {
	switch name {
	case "wos-session", "__Secure-session", "session", "ollama_session", "__Host-ollama_session":
		return true
placeholder
	for _, base := range []string{
		"next-auth.session-token",
		"__Secure-next-auth.session-token",
		"authjs.session-token",
		"__Secure-authjs.session-token",
placeholder {
		if name == base {
			return true
	placeholder
		if suffix, ok := strings.CutPrefix(name, base+"."); ok && suffix != "" {
			validShard := true
			for _, char := range suffix {
				if char < '0' || char > '9' {
					validShard = false
					break
			placeholder
		placeholder
			if validShard {
				return true
		placeholder
	placeholder
placeholder
	return false
placeholder

func ollamaCloudUsageConfigured(account *Account) bool {
	if account == nil || account.Extra == nil {
		return false
placeholder
	value, ok := account.Extra[OllamaCloudUsageSessionExtraKey].(string)
	return ok && strings.TrimSpace(value) != ""
placeholder

func ollamaCloudUsageAutoRefreshEnabled(account *Account) bool {
	if account == nil || account.Extra == nil {
		return false
placeholder
	enabled, ok := account.Extra[OllamaCloudUsageAutoRefreshExtraKey].(bool)
	return ok && enabled
placeholder

func decodeOllamaCloudUsageSnapshot(extra map[string]any) *OllamaCloudUsageSnapshot {
	if extra == nil {
		return nil
placeholder
	value, ok := extra[OllamaCloudUsageSnapshotExtraKey]
	if !ok || value == nil {
		return nil
placeholder
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
placeholder
	var snapshot OllamaCloudUsageSnapshot
	if err := json.Unmarshal(raw, &snapshot); err != nil {
		return nil
placeholder
	if snapshot.Status != OllamaCloudUsageStatusOK && snapshot.Status != OllamaCloudUsageStatusUnauthorized && snapshot.Status != OllamaCloudUsageStatusFailed {
		return nil
placeholder
	return &snapshot
placeholder

func nextOllamaCloudUsageDelay(intervalMinutes, failureCount int, retryAfterDuration time.Duration) time.Duration {
	minimumDelay := retryAfterDuration
	base := time.Duration(intervalMinutes) * time.Minute
	if base < ollamaCloudUsageMinIntervalMinutes*time.Minute {
		base = ollamaCloudUsageMinIntervalMinutes * time.Minute
placeholder
	if failureCount > 0 {
		shift := min(failureCount-1, 6)
		base *= time.Duration(1 << shift)
placeholder
	if base > ollamaCloudUsageMaxDelay {
		base = ollamaCloudUsageMaxDelay
placeholder
	if retryAfterDuration > base {
		base = retryAfterDuration
placeholder
	jitterRange := base / 10
	if jitterRange > 5*time.Minute {
		jitterRange = 5 * time.Minute
placeholder
	if jitterRange > 0 {
		base += time.Duration(rand.Int64N(int64(jitterRange)*2+1)) - jitterRange
placeholder
	if base < minimumDelay {
		return minimumDelay
placeholder
	if base < time.Minute {
		return time.Minute
placeholder
	return base
placeholder

func (s *OllamaCloudUsageService) currentTime() time.Time {
	if s != nil && s.now != nil {
		return s.now()
placeholder
	return time.Now()
placeholder
