package securityaudit

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/redis/go-redis/v9"
)

type activeConfigSnapshot struct {
	storage  storageConfig
	active   ActiveConfig
	loadedAt time.Time
placeholder

type ConfigManager struct {
	db        *sql.DB
	settings  service.SettingRepository
	redis     *redis.Client
	encryptor SecretEncryptor
	clock     Clock
	// encryptionKeyConfigured mirrors cfg.Totp.EncryptionKeyConfigured. With an
	// auto-generated (per-boot) key, newly saved endpoint tokens would become
	// undecryptable after the next restart, so Save rejects them (issue #4887).
	encryptionKeyConfigured bool

	snapshot atomic.Pointer[activeConfigSnapshot]
	expected atomic.Int64
	// expectedBlocking records the last storage intent that could be decoded,
	// independently of whether endpoint credentials or the full config could be
	// activated. A config version alone cannot distinguish async from blocking.
	expectedBlocking atomic.Bool
	// configUntrusted is set when a load/reload fails before a trustworthy
	// snapshot is installed. Combined with expectedBlocking, EffectiveMode
	// fails closed so a persisted blocking policy cannot be silently skipped
	// after startup or invalidation errors. Without blocking intent, untrusted
	// alone must not force ModeBlocking—Prompt Audit is default-off and must
	// not take the gateway down for every API request (see issue #4560).
	configUntrusted atomic.Bool

	stateMu       sync.RWMutex
	lastLoadError string
	lastErrorAt   *time.Time

	lifecycleMu sync.Mutex
	cancel      context.CancelFunc
	wg          sync.WaitGroup
placeholder

func NewConfigManager(db *sql.DB, settings service.SettingRepository, redisClient *redis.Client, encryptor service.SecretEncryptor, cfg *config.Config) *ConfigManager {
	return &ConfigManager{
		db: db, settings: settings, redis: redisClient, encryptor: encryptor, clock: realClock{placeholder,
		encryptionKeyConfigured: cfg != nil && cfg.Totp.EncryptionKeyConfigured,
placeholder
placeholder

func (m *ConfigManager) Start(ctx context.Context) error {
	if m == nil {
		return errors.New("prompt audit config manager unavailable")
placeholder
	m.lifecycleMu.Lock()
	if m.cancel != nil {
		m.lifecycleMu.Unlock()
		return nil
placeholder
	runCtx, cancel := context.WithCancel(ctx)
	m.cancel = cancel
	m.lifecycleMu.Unlock()
	loadErr := m.Reload(runCtx)
	if loadErr != nil {
		m.markConfigUntrusted()
placeholder
	m.wg.Add(1)
	go m.refreshLoop(runCtx)
	if m.redis != nil {
		m.wg.Add(1)
		go m.subscribeLoop(runCtx)
placeholder
	return loadErr
placeholder

func (m *ConfigManager) Shutdown(_ context.Context) error {
	if m == nil {
		return nil
placeholder
	m.lifecycleMu.Lock()
	cancel := m.cancel
	m.cancel = nil
	m.lifecycleMu.Unlock()
	if cancel != nil {
		cancel()
placeholder
	m.wg.Wait()
	return nil
placeholder

func (m *ConfigManager) Reload(ctx context.Context) error {
	if m == nil || m.settings == nil {
		m.markUntrustedIfNoActiveSnapshot()
		return errors.New("prompt audit setting repository unavailable")
placeholder
	values, err := m.settings.GetMultiple(ctx, []string{SettingKeyPromptAuditConfig, SettingKeyRiskControlplaceholder)
	if err != nil {
		m.recordLoadError(err)
		m.markUntrustedIfNoActiveSnapshot()
		return err
placeholder
	m.observeExpectedState(values[SettingKeyPromptAuditConfig], values[SettingKeyRiskControl] == "true")
	storage, err := ParseStorageConfig(values[SettingKeyPromptAuditConfig])
	if err != nil {
		m.recordLoadError(err)
		m.markUntrustedIfNoActiveSnapshot()
		return err
placeholder
	m.expected.Store(storage.ConfigVersion)
	m.expectedBlocking.Store(values[SettingKeyRiskControl] == "true" && storage.Enabled && storage.BlockingEnabled)
	active, err := ActiveFromStorage(storage, values[SettingKeyRiskControl] == "true", m.encryptor)
	if err != nil {
		m.recordLoadError(err)
		// expectedBlocking may already require fail-closed via BlockingActivationDegraded.
		m.markUntrustedIfNoActiveSnapshot()
		return err
placeholder
	now := m.clock.Now()
	previous := m.snapshot.Load()
	m.snapshot.Store(&activeConfigSnapshot{storage: cloneStorageConfig(storage), active: cloneActiveConfig(active), loadedAt: nowplaceholder)
	m.configUntrusted.Store(false)
	m.clearLoadError()
	m.logInvalidTokenEndpoints(previous, active)
	LogInfo(EventConfigLoaded, map[string]any{
		"config_version": storage.ConfigVersion, "status": "loaded",
placeholder)
	return nil
placeholder

// logInvalidTokenEndpoints warns once per change (not on every 5s refresh)
// when stored endpoint tokens cannot be decrypted with the current key.
func (m *ConfigManager) logInvalidTokenEndpoints(previous *activeConfigSnapshot, active ActiveConfig) {
	invalid := active.InvalidTokenEndpointIDs()
	if len(invalid) == 0 {
		return
placeholder
	if previous != nil {
		prior := previous.active.InvalidTokenEndpointIDs()
		if len(prior) == len(invalid) {
			same := true
			for i := range invalid {
				if prior[i] != invalid[i] {
					same = false
					break
			placeholder
		placeholder
			if same && previous.active.ConfigVersion == active.ConfigVersion {
				return
		placeholder
	placeholder
placeholder
	LogWarn(EventConfigTokenInvalid, map[string]any{
		"config_version": active.ConfigVersion, "status": "degraded",
		"error_code": "endpoint_token_undecryptable", "guard_endpoint_id": strings.Join(invalid, ","),
placeholder)
placeholder

func (m *ConfigManager) Active() (ActiveConfig, bool) {
	if m == nil {
		return ActiveConfig{placeholder, false
placeholder
	snapshot := m.snapshot.Load()
	if snapshot == nil {
		return ActiveConfig{placeholder, false
placeholder
	return cloneActiveConfig(snapshot.active), true
placeholder

func (m *ConfigManager) BlockingActivationDegraded() bool {
	if m == nil {
		return false
placeholder
	// Fail closed only when storage intent requires blocking. Untrusted config
	// without blocking intent must remain ModeOff so administrators can still
	// operate the gateway and turn Prompt Audit off after a failed reload.
	if !m.expectedBlocking.Load() {
		return false
placeholder
	if m.configUntrusted.Load() {
		return true
placeholder
	active, ok := m.Active()
	if !ok {
		return true
placeholder
	// A still-active weaker snapshot after a failed blocking activation must not
	// keep serving allow decisions under the old off/async mode.
	return active.EffectiveMode() != ModeBlocking
placeholder

func (m *ConfigManager) EffectiveMode() Mode {
	if m != nil && m.BlockingActivationDegraded() {
		return ModeBlocking
placeholder
	active, ok := m.Active()
	if !ok {
		return ModeOff
placeholder
	return active.EffectiveMode()
placeholder

func (m *ConfigManager) markConfigUntrusted() {
	if m == nil {
		return
placeholder
	m.configUntrusted.Store(true)
placeholder

func (m *ConfigManager) markUntrustedIfNoActiveSnapshot() {
	if m == nil {
		return
placeholder
	if _, ok := m.Active(); !ok {
		m.markConfigUntrusted()
placeholder
placeholder

func (m *ConfigManager) Public() (PublicConfig, error) {
	if m == nil {
		return PublicConfig{placeholder, infraerrors.ServiceUnavailable(ErrorCodeConfigUnavailable, "提示词审计配置暂不可用")
placeholder
	snapshot := m.snapshot.Load()
	if snapshot == nil {
		return PublicConfig{placeholder, infraerrors.ServiceUnavailable(ErrorCodeConfigUnavailable, "提示词审计配置暂不可用")
placeholder
	return PublicFromStorage(cloneStorageConfig(snapshot.storage), snapshot.active.RiskControlEnabled, snapshot.active.InvalidTokenEndpointIDs()), nil
placeholder

func (m *ConfigManager) Save(ctx context.Context, req UpdateConfigRequest, actorID int64) (PublicConfig, error) {
	if m == nil || m.db == nil || m.encryptor == nil {
		return PublicConfig{placeholder, errors.New("prompt audit config persistence unavailable")
placeholder
	if req.ExpectedConfigVersion < 1 {
		return PublicConfig{placeholder, infraerrors.BadRequest("prompt_audit_expected_config_version_required", "必须提供有效的配置版本")
placeholder
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommittedplaceholder)
	if err != nil {
		return PublicConfig{placeholder, err
placeholder
	defer func() { _ = tx.Rollback() placeholder()
	if _, err := tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock($1)`, promptAuditConfigLockKey); err != nil {
		return PublicConfig{placeholder, err
placeholder
	current := DefaultStorageConfig()
	var raw string
	err = tx.QueryRowContext(ctx, `SELECT value FROM settings WHERE key=$1 FOR UPDATE`, SettingKeyPromptAuditConfig).Scan(&raw)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return PublicConfig{placeholder, err
placeholder
	if err == nil {
		current, err = ParseStorageConfig(raw)
		if err != nil {
			return PublicConfig{placeholder, err
	placeholder
placeholder
	if current.ConfigVersion != req.ExpectedConfigVersion {
		return PublicConfig{placeholder, infraerrors.Conflict(ErrorCodeConfigConflict, "提示词审计配置已被其他管理员更新")
placeholder
	next, err := m.buildNextStorage(current, req, actorID)
	if err != nil {
		return PublicConfig{placeholder, err
placeholder
	next.ConfigVersion = current.ConfigVersion + 1
	next.UpdatedAt = m.clock.Now()
	next.UpdatedBy = actorID
	next.ChangeSummary = changeSummary(next)
	rawNext, err := json.Marshal(next)
	if err != nil {
		return PublicConfig{placeholder, err
placeholder
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO settings (key,value,updated_at) VALUES ($1,$2,NOW())
		ON CONFLICT (key) DO UPDATE SET value=EXCLUDED.value, updated_at=EXCLUDED.updated_at`,
		SettingKeyPromptAuditConfig, string(rawNext)); err != nil {
		return PublicConfig{placeholder, err
placeholder
	if err := tx.Commit(); err != nil {
		return PublicConfig{placeholder, err
placeholder
	// Install the snapshot with the current global gate, not merely the value
	// cached when this process last reloaded Prompt Audit configuration.
	riskControlEnabled := m.currentRiskControlEnabled()
	if values, getErr := m.settings.GetMultiple(ctx, []string{SettingKeyRiskControlplaceholder); getErr == nil {
		riskControlEnabled = values[SettingKeyRiskControl] == "true"
placeholder
	active, err := ActiveFromStorage(next, riskControlEnabled, m.encryptor)
	if err != nil {
		return PublicConfig{placeholder, err
placeholder
	m.expected.Store(next.ConfigVersion)
	m.expectedBlocking.Store(active.RiskControlEnabled && next.Enabled && next.BlockingEnabled)
	previous := m.snapshot.Load()
	m.snapshot.Store(&activeConfigSnapshot{storage: cloneStorageConfig(next), active: cloneActiveConfig(active), loadedAt: m.clock.Now()placeholder)
	// A successful admin save installs a trustworthy snapshot; clear any prior
	// fail-closed degradation so disabling audit actually takes effect.
	m.configUntrusted.Store(false)
	m.clearLoadError()
	m.logInvalidTokenEndpoints(previous, active)
	LogInfo(EventConfigUpdated, map[string]any{
		"config_version": next.ConfigVersion, "status": "updated",
placeholder)
	if m.redis != nil {
		if err := m.redis.Publish(ctx, ConfigInvalidationChannel, strconv.FormatInt(next.ConfigVersion, 10)).Err(); err != nil {
			LogWarn(EventConfigReloadDegraded, map[string]any{
				"config_version": next.ConfigVersion, "status": "degraded", "error_code": "config_invalidation_publish_failed",
		placeholder)
	placeholder
placeholder
	return PublicFromStorage(next, active.RiskControlEnabled, active.InvalidTokenEndpointIDs()), nil
placeholder

func (m *ConfigManager) buildNextStorage(current storageConfig, req UpdateConfigRequest, actorID int64) (storageConfig, error) {
	if err := validateUpdateConfigRequest(req); err != nil {
		return storageConfig{placeholder, err
placeholder
	currentByID := make(map[string]StorageEndpoint, len(current.Endpoints))
	for _, endpoint := range current.Endpoints {
		currentByID[endpoint.ID] = endpoint
placeholder
	next := storageConfig{
		Enabled: req.Enabled, BlockingEnabled: req.BlockingEnabled, StorePassEvents: req.StorePassEvents,
		Strategy: strings.TrimSpace(req.Strategy), WorkerCount: req.WorkerCount,
		QueueCapacity: req.QueueCapacity, Scanners: append([]string(nil), req.Scanners...),
		AllGroups: req.AllGroups, GroupIDs: append([]int64(nil), req.GroupIDs...),
		ConfigVersion: current.ConfigVersion, UpdatedBy: actorID,
		Endpoints: make([]StorageEndpoint, 0, len(req.Endpoints)),
placeholder
	for _, endpoint := range req.Endpoints {
		baseURL, err := NormalizeBaseURL(endpoint.BaseURL)
		if err != nil {
			return storageConfig{placeholder, err
	placeholder
		stored := StorageEndpoint{
			ID: strings.TrimSpace(endpoint.ID), Name: strings.TrimSpace(endpoint.Name),
			Protocol: strings.TrimSpace(endpoint.Protocol), BaseURL: baseURL, Model: strings.TrimSpace(endpoint.Model),
			TimeoutMS: endpoint.TimeoutMS, InputLimit: endpoint.InputLimit, Enabled: endpoint.Enabled,
	placeholder
		old, hadOld := currentByID[stored.ID]
		switch {
		case endpoint.ClearToken:
			stored.TokenCiphertext = ""
		case strings.TrimSpace(endpoint.Token) != "":
			if !m.encryptionKeyConfigured {
				return storageConfig{placeholder, infraerrors.BadRequest(ErrorCodeEncryptionKeyRequired,
					"未配置固定加密密钥，审计节点 Token 将在服务重启后失效。请先设置 TOTP_ENCRYPTION_KEY 环境变量（64 位十六进制）并重启服务")
		placeholder
			ciphertext, err := m.encryptor.Encrypt(strings.TrimSpace(endpoint.Token))
			if err != nil {
				return storageConfig{placeholder, fmt.Errorf("encrypt prompt audit endpoint token: %w", err)
		placeholder
			stored.TokenCiphertext = ciphertext
		case hadOld:
			stored.TokenCiphertext = old.TokenCiphertext
	placeholder
		next.Endpoints = append(next.Endpoints, stored)
placeholder
	normalizeStorageConfig(&next)
	if err := validateStorageConfig(next); err != nil {
		return storageConfig{placeholder, err
placeholder
	return next, nil
placeholder

func (m *ConfigManager) RuntimeState() (expected int64, active int64, loadedAt *time.Time, loadError string) {
	if m == nil {
		return 1, 0, nil, "config_manager_unavailable"
placeholder
	expected = m.expected.Load()
	if expected < 1 {
		expected = 1
placeholder
	if snapshot := m.snapshot.Load(); snapshot != nil {
		active = snapshot.active.ConfigVersion
		value := snapshot.loadedAt
		loadedAt = &value
placeholder
	m.stateMu.RLock()
	loadError = m.lastLoadError
	m.stateMu.RUnlock()
	return
placeholder

func (m *ConfigManager) Encrypt(value string) (string, error) { return m.encryptor.Encrypt(value) placeholder
func (m *ConfigManager) Decrypt(value string) (string, error) { return m.encryptor.Decrypt(value) placeholder

func (m *ConfigManager) currentRiskControlEnabled() bool {
	if snapshot := m.snapshot.Load(); snapshot != nil {
		return snapshot.active.RiskControlEnabled
placeholder
	return false
placeholder

func (m *ConfigManager) observeExpectedState(raw string, riskControlEnabled bool) {
	if m == nil {
		return
placeholder
	if strings.TrimSpace(raw) == "" {
		m.expected.Store(1)
		m.expectedBlocking.Store(false)
		return
placeholder
	var intent struct {
		Enabled         bool  `json:"enabled"`
		BlockingEnabled bool  `json:"blocking_enabled"`
		ConfigVersion   int64 `json:"config_version"`
placeholder
	if err := json.Unmarshal([]byte(raw), &intent); err != nil {
		return
placeholder
	if intent.ConfigVersion < 1 {
		intent.ConfigVersion = 1
placeholder
	m.expected.Store(intent.ConfigVersion)
	m.expectedBlocking.Store(riskControlEnabled && intent.Enabled && intent.BlockingEnabled)
placeholder

func (m *ConfigManager) refreshLoop(ctx context.Context) {
	defer m.wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := m.Reload(ctx); err != nil {
				LogWarn(EventConfigReloadDegraded, map[string]any{"status": "degraded", "error_code": "config_ttl_reload_failed"placeholder)
		placeholder
	placeholder
placeholder
placeholder

func (m *ConfigManager) subscribeLoop(ctx context.Context) {
	defer m.wg.Done()
	pubsub := m.redis.Subscribe(ctx, ConfigInvalidationChannel)
	defer func() { _ = pubsub.Close() placeholder()
	channel := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case message, ok := <-channel:
			if !ok {
				return
		placeholder
			version, err := strconv.ParseInt(strings.TrimSpace(message.Payload), 10, 64)
			if err != nil || version < 1 {
				continue
		placeholder
			m.expected.Store(version)
			if err := m.Reload(ctx); err != nil {
				// A newer published version failed to activate. Until reload
				// succeeds, do not keep serving a potentially stale weaker mode.
				if active, ok := m.Active(); !ok || active.ConfigVersion < version {
					m.markConfigUntrusted()
			placeholder
				LogWarn(EventConfigReloadDegraded, map[string]any{
					"config_version": version, "status": "degraded", "error_code": "config_invalidation_reload_failed",
			placeholder)
		placeholder
	placeholder
placeholder
placeholder

func (m *ConfigManager) recordLoadError(_ error) {
	if m == nil {
		return
placeholder
	now := m.clock.Now()
	m.stateMu.Lock()
	m.lastLoadError = stableErrorMessage("config_load_failed")
	m.lastErrorAt = &now
	m.stateMu.Unlock()
placeholder

func (m *ConfigManager) clearLoadError() {
	m.stateMu.Lock()
	m.lastLoadError = ""
	m.lastErrorAt = nil
	m.stateMu.Unlock()
placeholder

func cloneStorageConfig(cfg storageConfig) storageConfig {
	cfg.Scanners = append([]string(nil), cfg.Scanners...)
	cfg.GroupIDs = append([]int64(nil), cfg.GroupIDs...)
	cfg.Endpoints = append([]StorageEndpoint(nil), cfg.Endpoints...)
	return cfg
placeholder

func cloneActiveConfig(cfg ActiveConfig) ActiveConfig {
	cfg.Scanners = append([]string(nil), cfg.Scanners...)
	cfg.GroupIDs = append([]int64(nil), cfg.GroupIDs...)
	cfg.Endpoints = append([]ActiveEndpoint(nil), cfg.Endpoints...)
	return cfg
placeholder
