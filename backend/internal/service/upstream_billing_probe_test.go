package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/stretchr/testify/require"
)

type upstreamBillingProbeAccountRepo struct {
	AccountRepository
	mu          sync.Mutex
	accounts    map[int64]*Account
	updates     map[int64][]map[string]any
	bulkUpdates []AccountBulkUpdate
placeholder

type staleDueUpstreamBillingProbeAccountRepo struct {
	*upstreamBillingProbeAccountRepo
	due []Account
placeholder

func (r *staleDueUpstreamBillingProbeAccountRepo) ListDueUpstreamBillingProbeAccounts(_ context.Context, _ time.Time, limit int) ([]Account, error) {
	if limit < len(r.due) {
		return append([]Account(nil), r.due[:limit]...), nil
placeholder
	return append([]Account(nil), r.due...), nil
placeholder

func (r *upstreamBillingProbeAccountRepo) Create(_ context.Context, account *Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.accounts == nil {
		r.accounts = make(map[int64]*Account)
placeholder
	if account.ID == 0 {
		account.ID = int64(len(r.accounts) + 1)
placeholder
	r.accounts[account.ID] = account
	return nil
placeholder

func (r *upstreamBillingProbeAccountRepo) Update(_ context.Context, account *Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.accounts[account.ID] = account
	return nil
placeholder

func (r *upstreamBillingProbeAccountRepo) BulkUpdate(_ context.Context, ids []int64, updates AccountBulkUpdate) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.bulkUpdates = append(r.bulkUpdates, updates)
	return int64(len(ids)), nil
placeholder

func (r *upstreamBillingProbeAccountRepo) GetByID(_ context.Context, id int64) (*Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	account := r.accounts[id]
	if account == nil {
		return nil, ErrAccountNotFound
placeholder
	clone := *account
	clone.Credentials = mergeMap(nil, account.Credentials)
	clone.Extra = mergeMap(nil, account.Extra)
	return &clone, nil
placeholder

func (r *upstreamBillingProbeAccountRepo) GetByIDs(_ context.Context, ids []int64) ([]*Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]*Account, 0, len(ids))
	for _, id := range ids {
		if account := r.accounts[id]; account != nil {
			result = append(result, account)
	placeholder
placeholder
	return result, nil
placeholder

func (r *upstreamBillingProbeAccountRepo) UpdateExtra(_ context.Context, id int64, updates map[string]any) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	account := r.accounts[id]
	if account == nil {
		return ErrAccountNotFound
placeholder
	if account.Extra == nil {
		account.Extra = make(map[string]any)
placeholder
	for key, value := range updates {
		account.Extra[key] = value
placeholder
	if r.updates == nil {
		r.updates = make(map[int64][]map[string]any)
placeholder
	r.updates[id] = append(r.updates[id], updates)
	return nil
placeholder

func (r *upstreamBillingProbeAccountRepo) UpdateUpstreamBillingProbeSnapshot(
	_ context.Context,
	expected *Account,
	snapshot *UpstreamBillingProbeSnapshot,
	rateMultiplier *float64,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	account := r.accounts[expected.ID]
	if account == nil || account.Platform != expected.Platform || account.Type != expected.Type || !reflect.DeepEqual(account.Credentials, expected.Credentials) {
		return ErrUpstreamBillingProbeIdentityChanged
placeholder
	if account.Extra == nil {
		account.Extra = make(map[string]any)
placeholder
	account.Extra[UpstreamBillingProbeExtraKey] = snapshot
	if snapshot.Status == UpstreamBillingProbeStatusOK &&
		rateMultiplier != nil &&
		upstreamBillingRateSyncEnabled(account) {
		value := *rateMultiplier
		account.RateMultiplier = &value
placeholder
	return nil
placeholder

func (r *upstreamBillingProbeAccountRepo) FindByExtraField(_ context.Context, key string, value any) ([]Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]Account, 0)
	for _, account := range r.accounts {
		if account.Extra != nil && account.Extra[key] == value {
			result = append(result, *account)
	placeholder
placeholder
	return result, nil
placeholder

type upstreamBillingProbeSettingRepo struct {
	SettingRepository
	mu     sync.Mutex
	values map[string]string
placeholder

type upstreamBillingProbeHTTPStub struct {
	calls          atomic.Int64
	active         atomic.Int64
	maxActive      atomic.Int64
	beforeResponse func()
placeholder

func (u *upstreamBillingProbeHTTPStub) Do(req *http.Request, proxyURL string, accountID int64, accountConcurrency int) (*http.Response, error) {
	u.calls.Add(1)
	active := u.active.Add(1)
	defer u.active.Add(-1)
	for {
		peak := u.maxActive.Load()
		if active <= peak || u.maxActive.CompareAndSwap(peak, active) {
			break
	placeholder
placeholder
	if u.beforeResponse != nil {
		u.beforeResponse()
placeholder
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(`{
			"object":"sub2api.key_billing",
			"schema_version":1,
			"billing_scope":"token",
			"group_rate_multiplier":0.8,
			"resolved_rate_multiplier":0.8,
			"peak_rate_enabled":false,
			"effective_rate_multiplier":0.8,
			"observed_at":"2026-07-13T01:00:00Z"
	placeholder`)),
placeholder, nil
placeholder

func (u *upstreamBillingProbeHTTPStub) DoWithTLS(req *http.Request, proxyURL string, accountID int64, accountConcurrency int, profile *tlsfingerprint.Profile) (*http.Response, error) {
	return u.Do(req, proxyURL, accountID, accountConcurrency)
placeholder

func (r *upstreamBillingProbeSettingRepo) GetValue(_ context.Context, key string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	value, ok := r.values[key]
	if !ok {
		return "", ErrSettingNotFound
placeholder
	return value, nil
placeholder

func (r *upstreamBillingProbeSettingRepo) Set(_ context.Context, key, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.values == nil {
		r.values = make(map[string]string)
placeholder
	r.values[key] = value
	return nil
placeholder

func newUpstreamBillingProbeTestService(
	repo AccountRepository,
	upstream HTTPUpstream,
	settingRepo SettingRepository,
) *UpstreamBillingProbeService {
	cfg := &config.Config{Security: config.SecurityConfig{URLAllowlist: config.URLAllowlistConfig{
		Enabled:           false,
		AllowInsecureHTTP: true,
placeholderplaceholderplaceholder
	accountTestService := &AccountTestService{accountRepo: repo, httpUpstream: upstream, cfg: cfgplaceholder
	return NewUpstreamBillingProbeService(repo, accountTestService, NewSettingService(settingRepo, cfg))
placeholder

func TestUpstreamBillingProbeSettingsDefaultsAndValidation(t *testing.T) {
	repo := &upstreamBillingProbeSettingRepo{placeholder
	settingsService := NewSettingService(repo, &config.Config{placeholder)

	settings, err := settingsService.GetUpstreamBillingProbeSettings(context.Background())
placeholder
	require.True(t, settings.Enabled)
	require.Equal(t, 30, settings.IntervalMinutes)

	err = settingsService.SetUpstreamBillingProbeSettings(context.Background(), &UpstreamBillingProbeSettings{
		Enabled:         false,
		IntervalMinutes: 4,
placeholder)
placeholder
	require.Contains(t, err.Error(), "interval_minutes must be between 5 and 1440")

	err = settingsService.SetUpstreamBillingProbeSettings(context.Background(), &UpstreamBillingProbeSettings{
		Enabled:         false,
		IntervalMinutes: 60,
placeholder)
placeholder
	settings, err = settingsService.GetUpstreamBillingProbeSettings(context.Background())
placeholder
	require.False(t, settings.Enabled)
	require.Equal(t, 60, settings.IntervalMinutes)

	repo.values[SettingKeyUpstreamBillingProbeSettings] = `{"interval_minutes":45placeholder`
	settings, err = settingsService.GetUpstreamBillingProbeSettings(context.Background())
placeholder
	require.True(t, settings.Enabled)
	require.Equal(t, 45, settings.IntervalMinutes)
	repo.values[SettingKeyUpstreamBillingProbeSettings] = `{"enabled":falseplaceholder`
	settings, err = settingsService.GetUpstreamBillingProbeSettings(context.Background())
placeholder
	require.False(t, settings.Enabled)
	require.Equal(t, 30, settings.IntervalMinutes)

	repo.values[SettingKeyUpstreamBillingProbeSettings] = `{"enabled":`
	settings, err = settingsService.GetUpstreamBillingProbeSettings(context.Background())
	require.ErrorContains(t, err, "parse upstream billing probe settings")
	require.Nil(t, settings)
placeholder

func TestUpstreamBillingProbeSuccessPersistsSanitizedSnapshot(t *testing.T) {
	initialRate := 0.25
	account := &Account{
		ID:          17,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Concurrency: 2,
placeholder
			"api_key":  "sk-sensitive",
			"base_url": "https://upstream.example/v1",
	placeholder,
		Extra: map[string]any{
			UpstreamBillingProbeEnabledExtraKey:    true,
			UpstreamBillingRateSyncEnabledExtraKey: true,
	placeholder,
		RateMultiplier: &initialRate,
placeholder
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{account.ID: accountplaceholderplaceholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(`{
			"object":"sub2api.key_billing",
			"schema_version":1,
			"billing_scope":"token",
			"group_rate_multiplier":0.8,
			"user_rate_multiplier":0.6,
			"resolved_rate_multiplier":0.6,
			"peak_rate_enabled":true,
			"peak_start":"09:00",
			"peak_end":"18:00",
			"peak_rate_multiplier":1.5,
			"applied_peak_multiplier":1.5,
			"effective_rate_multiplier":0.9,
			"timezone":"Asia/Shanghai",
			"observed_at":"2026-07-13T01:00:00Z",
			"unexpected_secret":"must-not-persist"
	placeholder`)),
placeholderplaceholder
	svc := newUpstreamBillingProbeTestService(repo, upstream, &upstreamBillingProbeSettingRepo{placeholder)
	fixedNow := time.Date(2026, time.July, 13, 2, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return fixedNow placeholder

	snapshot, err := svc.ProbeAccount(context.Background(), account.ID)
placeholder
	require.Equal(t, UpstreamBillingProbeStatusOK, snapshot.Status)
	require.Equal(t, 0.9, snapshot.Data["effective_rate_multiplier"])
	require.NotContains(t, snapshot.Data, "unexpected_secret")
	require.NotNil(t, snapshot.ReceivedAt)
	require.Equal(t, fixedNow, *snapshot.ReceivedAt)
	require.NotNil(t, snapshot.FreshUntil)
	require.Equal(t, fixedNow.Add(time.Hour), *snapshot.FreshUntil)
	require.False(t, snapshot.NextProbeAt.Before(fixedNow.Add(24*time.Minute)))
	require.False(t, snapshot.NextProbeAt.After(fixedNow.Add(36*time.Minute)))
	// 写回的是不含高峰因子的 resolved 倍率（0.6），不是探测那一刻含高峰的
	// effective 倍率（0.9）——否则一个探测周期的峰值会被冻结进静态列。
	require.NotNil(t, account.RateMultiplier)
	require.Equal(t, 0.6, *account.RateMultiplier)
	require.NotNil(t, snapshot.SyncedRateMultiplier)
	require.Equal(t, 0.6, *snapshot.SyncedRateMultiplier)
	require.Equal(t, "https://upstream.example/v1/sub2api/billing", upstream.lastReq.URL.String())
	require.Equal(t, http.MethodGet, upstream.lastReq.Method)
	require.Equal(t, "Bearer sk-sensitive", upstream.lastReq.Header.Get("Authorization"))
	require.True(t, HTTPUpstreamRedirectsDisabled(upstream.lastReq.Context()))

	persisted := decodeUpstreamBillingProbeSnapshot(account.Extra)
	require.NotNil(t, persisted)
	require.Equal(t, snapshot.Status, persisted.Status)
placeholder

func TestUpstreamBillingProbeAdaptiveCNUsesChatProtocolBaseURL(t *testing.T) {
	account := &Account{
		ID:          18,
		Platform:    PlatformKimi,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"api_key":      "sk-sensitive",
			"api_protocol": APIProtocolAdaptive,
			"base_url":     "https://legacy-relay.example/v1",
			"api_base_urls": map[string]any{
				APIProtocolChatCompletions: "https://chat-relay.example/v1",
		placeholder,
	placeholder,
		Extra: map[string]any{UpstreamBillingProbeEnabledExtraKey: trueplaceholder,
placeholder
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{account.ID: accountplaceholderplaceholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       upstreamBillingProbeValidBody(),
placeholderplaceholder
	svc := newUpstreamBillingProbeTestService(repo, upstream, &upstreamBillingProbeSettingRepo{placeholder)

	snapshot, err := svc.ProbeAccount(context.Background(), account.ID)

placeholder
	require.Equal(t, UpstreamBillingProbeStatusOK, snapshot.Status)
	require.Equal(t, "https://chat-relay.example/v1/sub2api/billing", upstream.lastReq.URL.String())
placeholder

func TestUpstreamBillingProbeSyncsResolvedRateForAllAPIKeyPlatforms(t *testing.T) {
	for _, platform := range []string{
		PlatformOpenAI,
		PlatformAnthropic,
		PlatformGemini,
		PlatformAntigravity,
		PlatformGrok,
placeholder {
		t.Run(platform, func(t *testing.T) {
			initialRate := 0.25
			account := &Account{
				ID:             17,
				Platform:       platform,
				Type:           AccountTypeAPIKey,
				Status:         StatusActive,
				Concurrency:    1,
				RateMultiplier: &initialRate,
		placeholder
					"api_key":  "sk-sensitive",
					"base_url": "https://upstream.example",
			placeholder,
				Extra: map[string]any{
					UpstreamBillingProbeEnabledExtraKey:    true,
					UpstreamBillingRateSyncEnabledExtraKey: true,
			placeholder,
		placeholder
			repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{account.ID: accountplaceholderplaceholder
			svc := newUpstreamBillingProbeTestService(repo, &upstreamBillingProbeHTTPStub{placeholder, &upstreamBillingProbeSettingRepo{placeholder)

			snapshot, err := svc.ProbeAccount(context.Background(), account.ID)

		placeholder
			require.Equal(t, UpstreamBillingProbeStatusOK, snapshot.Status)
			require.NotNil(t, account.RateMultiplier)
			require.Equal(t, 0.8, *account.RateMultiplier)
	placeholder)
placeholder
placeholder

func TestUpstreamBillingProbeOnlyDoesNotChangeAccountRate(t *testing.T) {
	initialRate := 0.25
	account := &Account{
		ID:             18,
		Platform:       PlatformGrok,
		Type:           AccountTypeAPIKey,
		Status:         StatusActive,
		Concurrency:    1,
		RateMultiplier: &initialRate,
placeholder
			"api_key":  "sk-sensitive",
			"base_url": "https://upstream.example",
	placeholder,
		Extra: map[string]any{UpstreamBillingProbeEnabledExtraKey: trueplaceholder,
placeholder
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{account.ID: accountplaceholderplaceholder
	svc := newUpstreamBillingProbeTestService(repo, &upstreamBillingProbeHTTPStub{placeholder, &upstreamBillingProbeSettingRepo{placeholder)

	snapshot, err := svc.ProbeAccount(context.Background(), account.ID)

placeholder
	require.Equal(t, UpstreamBillingProbeStatusOK, snapshot.Status)
	require.NotNil(t, account.RateMultiplier)
	require.Equal(t, initialRate, *account.RateMultiplier)
	require.Contains(t, account.Extra, UpstreamBillingProbeExtraKey)
placeholder

func TestUpstreamBillingProbeSyncRateRangeAndPrecision(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  float64
		ok    bool
placeholder{
		{name: "round to four decimals", value: 0.07654, want: 0.0765, ok: trueplaceholder,
		{name: "maximum", value: upstreamBillingRateSyncMaxMultiplier, want: upstreamBillingRateSyncMaxMultiplier, ok: trueplaceholder,
		// 0 会让 accountCost 恒为 0，账号配额与成本告警全部静默失效，
		// 自动写回一律拒绝（管理员手工设 0 仍然允许）。
		{name: "zero is rejected", value: 0, ok: falseplaceholder,
		{name: "positive below database precision rounds to zero", value: 0.00001, ok: falseplaceholder,
		{name: "just above the write-back ceiling", value: 100.0001, ok: falseplaceholder,
		{name: "column ceiling is far above the write-back ceiling", value: 999999.9999, ok: falseplaceholder,
		{name: "negative", value: -1, ok: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := upstreamBillingProbeSyncRate(map[string]any{"resolved_rate_multiplier": tt.valueplaceholder)
			require.Equal(t, tt.ok, ok)
			if tt.ok {
				require.Equal(t, tt.want, got)
		placeholder
	placeholder)
placeholder
placeholder

// 只读取 resolved（时间无关的基准倍率）：effective 含探测那一刻的高峰系数，
// 写回它会把一个探测周期的峰值/谷值冻结进静态列。
func TestUpstreamBillingProbeSyncRateIgnoresEffectiveRate(t *testing.T) {
	got, ok := upstreamBillingProbeSyncRate(map[string]any{
		"resolved_rate_multiplier":  0.6,
		"effective_rate_multiplier": 0.9,
placeholder)
	require.True(t, ok)
	require.Equal(t, 0.6, got)

	_, ok = upstreamBillingProbeSyncRate(map[string]any{"effective_rate_multiplier": 0.9placeholder)
	require.False(t, ok)
placeholder

// 上游声明超出自动写回值域时保持原倍率，但探测本身是成功的：
// 快照照常记 ok，不累计 failure_count、不进入退避。
func TestUpstreamBillingProbeKeepsRateWhenDeclarationOutOfSyncRange(t *testing.T) {
	for _, tt := range []struct {
		name     string
		declared string
placeholder{
		{name: "zero", declared: "0"placeholder,
		{name: "above ceiling", declared: "1000"placeholder,
placeholder {
		t.Run(tt.name, func(t *testing.T) {
			initialRate := 0.25
			account := &Account{
				ID:             21,
				Platform:       PlatformOpenAI,
				Type:           AccountTypeAPIKey,
				Status:         StatusActive,
				Concurrency:    1,
				RateMultiplier: &initialRate,
		placeholder
					"api_key":  "sk-sensitive",
					"base_url": "https://upstream.example",
			placeholder,
				Extra: map[string]any{
					UpstreamBillingProbeEnabledExtraKey:    true,
					UpstreamBillingRateSyncEnabledExtraKey: true,
			placeholder,
		placeholder
			repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{account.ID: accountplaceholderplaceholder
			upstream := &httpUpstreamRecorder{resp: &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
				Body: io.NopCloser(strings.NewReader(fmt.Sprintf(`{
					"object":"sub2api.key_billing",
					"schema_version":1,
					"billing_scope":"token",
					"group_rate_multiplier":%[1]s,
					"resolved_rate_multiplier":%[1]s,
					"peak_rate_enabled":false,
					"effective_rate_multiplier":%[1]s,
					"observed_at":"2026-07-13T01:00:00Z"
			placeholder`, tt.declared))),
		placeholderplaceholder
			svc := newUpstreamBillingProbeTestService(repo, upstream, &upstreamBillingProbeSettingRepo{placeholder)

			snapshot, err := svc.ProbeAccount(context.Background(), account.ID)

		placeholder
			require.Equal(t, UpstreamBillingProbeStatusOK, snapshot.Status)
			require.Zero(t, snapshot.FailureCount)
			require.Nil(t, snapshot.SyncedRateMultiplier)
			require.NotNil(t, account.RateMultiplier)
			require.Equal(t, initialRate, *account.RateMultiplier)
			// 原始声明仍进快照供展示。
			require.Equal(t, snapshot.Data["resolved_rate_multiplier"], snapshot.Data["effective_rate_multiplier"])
	placeholder)
placeholder
placeholder

// 未开启同步的账号只观察上游声明：声明值不适配 accounts.rate_multiplier
// 不得被记成探测失败（否则会累计 failure_count 并进入指数退避）。
func TestUpstreamBillingProbeWithoutSyncIgnoresUnusableDeclaredRate(t *testing.T) {
	account := &Account{
		ID:          22,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"api_key":  "sk-sensitive",
			"base_url": "https://upstream.example",
	placeholder,
		Extra: map[string]any{UpstreamBillingProbeEnabledExtraKey: trueplaceholder,
placeholder
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{account.ID: accountplaceholderplaceholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(`{
			"object":"sub2api.key_billing",
			"schema_version":1,
			"billing_scope":"token",
			"group_rate_multiplier":0,
			"resolved_rate_multiplier":0,
			"peak_rate_enabled":false,
			"effective_rate_multiplier":0,
			"observed_at":"2026-07-13T01:00:00Z"
	placeholder`)),
placeholderplaceholder
	svc := newUpstreamBillingProbeTestService(repo, upstream, &upstreamBillingProbeSettingRepo{placeholder)

	snapshot, err := svc.ProbeAccount(context.Background(), account.ID)

placeholder
	require.Equal(t, UpstreamBillingProbeStatusOK, snapshot.Status)
	require.Zero(t, snapshot.FailureCount)
	require.Empty(t, snapshot.LastError)
	require.Nil(t, snapshot.SyncedRateMultiplier)
	require.Equal(t, float64(0), snapshot.Data["resolved_rate_multiplier"])
	require.Nil(t, account.RateMultiplier)
placeholder

func TestUpstreamBillingProbeRejectsMissingRequiredMultiplier(t *testing.T) {
	_, err := parseUpstreamBillingProbeResponse([]byte(`{
		"object":"sub2api.key_billing",
		"schema_version":1,
		"billing_scope":"token",
		"group_rate_multiplier":0.8,
		"peak_rate_enabled":false,
		"effective_rate_multiplier":0.8,
		"observed_at":"2026-07-13T01:00:00Z"
placeholder`))

	require.ErrorContains(t, err, "incomplete billing response")
placeholder

func TestUpstreamBillingProbeDiscardsResultWhenIdentityChangesInFlight(t *testing.T) {
	account := &Account{
		ID:          19,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Concurrency: 1,
placeholder"api_key": "sk-old", "base_url": "https://upstream.example"placeholder,
		Extra:       map[string]any{UpstreamBillingProbeEnabledExtraKey: trueplaceholder,
placeholder
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{account.ID: accountplaceholderplaceholder
	upstream := &upstreamBillingProbeHTTPStub{beforeResponse: func() {
		repo.mu.Lock()
		defer repo.mu.Unlock()
		repo.accounts[account.ID].Credentials = map[string]any{"api_key": "sk-new", "base_url": "https://new.example"placeholder
placeholderplaceholder
	svc := newUpstreamBillingProbeTestService(repo, upstream, &upstreamBillingProbeSettingRepo{placeholder)

	snapshot, err := svc.ProbeAccount(context.Background(), account.ID)

	require.Nil(t, snapshot)
	require.ErrorIs(t, err, ErrUpstreamBillingProbeIdentityChanged)
	require.NotContains(t, repo.accounts[account.ID].Extra, UpstreamBillingProbeExtraKey)
placeholder

func TestUpstreamBillingProbeRejectsInvalidPeakConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		start    string
		end      string
		timezone string
placeholder{
		{name: "invalid start", start: "25:00", end: "18:00", timezone: "UTC"placeholder,
		{name: "cross midnight", start: "22:00", end: "02:00", timezone: "UTC"placeholder,
		{name: "invalid timezone", start: "09:00", end: "18:00", timezone: "Mars/Olympus"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := fmt.Sprintf(`{
				"object":"sub2api.key_billing",
				"schema_version":1,
				"billing_scope":"token",
				"group_rate_multiplier":0.8,
				"resolved_rate_multiplier":0.8,
				"peak_rate_enabled":true,
				"peak_start":%q,
				"peak_end":%q,
				"peak_rate_multiplier":1.5,
				"applied_peak_multiplier":1,
				"effective_rate_multiplier":0.8,
				"timezone":%q,
				"observed_at":"2026-07-13T01:00:00Z"
		placeholder`, tt.start, tt.end, tt.timezone)

			_, err := parseUpstreamBillingProbeResponse([]byte(body))
			require.ErrorContains(t, err, "invalid peak billing response")
	placeholder)
placeholder
placeholder

func TestUpstreamBillingProbeRejectsInconsistentMultipliers(t *testing.T) {
	tests := []struct {
		name string
		body string
placeholder{
		{
			name: "resolved does not use user override",
			body: `{
				"object":"sub2api.key_billing","schema_version":1,"billing_scope":"token",
				"group_rate_multiplier":0.8,"user_rate_multiplier":0.5,"resolved_rate_multiplier":0.8,
				"peak_rate_enabled":false,"effective_rate_multiplier":0.8,"observed_at":"2026-07-13T01:00:00Z"
		placeholder`,
	placeholder,
		{
			name: "effective rate does not match resolved rate",
			body: `{
				"object":"sub2api.key_billing","schema_version":1,"billing_scope":"token",
				"group_rate_multiplier":0.8,"resolved_rate_multiplier":0.8,
				"peak_rate_enabled":false,"effective_rate_multiplier":1.2,"observed_at":"2026-07-13T01:00:00Z"
		placeholder`,
	placeholder,
		{
			name: "applied peak does not match observed window",
			body: `{
				"object":"sub2api.key_billing","schema_version":1,"billing_scope":"token",
				"group_rate_multiplier":0.8,"resolved_rate_multiplier":0.8,
				"peak_rate_enabled":true,"peak_start":"09:00","peak_end":"18:00",
				"peak_rate_multiplier":1.5,"applied_peak_multiplier":1,
				"effective_rate_multiplier":0.8,"timezone":"Asia/Shanghai","observed_at":"2026-07-13T01:00:00Z"
		placeholder`,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseUpstreamBillingProbeResponse([]byte(tt.body))
			require.ErrorContains(t, err, "inconsistent")
	placeholder)
placeholder
placeholder

func TestUpstreamBillingRateAtHandlesDST(t *testing.T) {
	data := map[string]any{
		"billing_scope":            "token",
		"resolved_rate_multiplier": 1.0,
		"peak_rate_enabled":        true,
		"peak_start":               "02:00",
		"peak_end":                 "04:00",
		"peak_rate_multiplier":     2.0,
		"timezone":                 "America/New_York",
placeholder
	beforeJump := time.Date(2026, time.March, 8, 6, 30, 0, 0, time.UTC)
	afterJump := time.Date(2026, time.March, 8, 7, 30, 0, 0, time.UTC)

	rate, ok := upstreamBillingRateAt(data, beforeJump)
	require.True(t, ok)
	require.Equal(t, 1.0, rate)
	rate, ok = upstreamBillingRateAt(data, afterJump)
	require.True(t, ok)
	require.Equal(t, 2.0, rate)
placeholder

func TestUpstreamBillingProbeFailurePreservesLastSuccessAndRetryAfter(t *testing.T) {
	receivedAt := time.Date(2026, time.July, 12, 12, 0, 0, 0, time.UTC)
	initialRate := 0.35
	previous := &UpstreamBillingProbeSnapshot{
		Status:       UpstreamBillingProbeStatusOK,
		Data:         map[string]any{"effective_rate_multiplier": 0.5placeholder,
		ReceivedAt:   &receivedAt,
		FailureCount: 1,
placeholder
	account := &Account{
		ID:             18,
		Platform:       PlatformOpenAI,
		Type:           AccountTypeAPIKey,
		Status:         StatusActive,
		Concurrency:    1,
		RateMultiplier: &initialRate,
		Credentials:    map[string]any{"api_key": "sk-test", "base_url": "https://upstream.example"placeholder,
		Extra: map[string]any{
			UpstreamBillingProbeEnabledExtraKey: true,
			UpstreamBillingProbeExtraKey:        previous,
	placeholder,
placeholder
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{account.ID: accountplaceholderplaceholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{"Retry-After": []string{"14400"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"error":"do not persist this"placeholder`)),
placeholderplaceholder
	svc := newUpstreamBillingProbeTestService(repo, upstream, &upstreamBillingProbeSettingRepo{placeholder)
	fixedNow := time.Date(2026, time.July, 13, 2, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return fixedNow placeholder

	snapshot, err := svc.ProbeAccount(context.Background(), account.ID)
placeholder
	require.Equal(t, UpstreamBillingProbeStatusFailed, snapshot.Status)
	require.Equal(t, previous.Data, snapshot.Data)
	require.Equal(t, previous.ReceivedAt, snapshot.ReceivedAt)
	require.NotNil(t, snapshot.FreshUntil)
	require.Equal(t, receivedAt.Add(time.Hour), *snapshot.FreshUntil)
	require.Equal(t, 2, snapshot.FailureCount)
	require.Equal(t, "http_error", snapshot.LastError)
	require.False(t, snapshot.NextProbeAt.Before(fixedNow.Add(4*time.Hour)))
	require.NotContains(t, snapshot.LastError, "do not persist")
	require.NotNil(t, account.RateMultiplier)
	require.Equal(t, initialRate, *account.RateMultiplier)
placeholder

func TestUpstreamBillingProbeRetryAfterIsNotShortened(t *testing.T) {
	delay := nextProbeDelay(30, 48*time.Hour)
	require.Equal(t, 48*time.Hour, delay)
placeholder

// unsupported 的重探间隔明显长于普通失败，但始终有上界：上游后来接入 sub2api
// 时最迟一天内会被重新发现，且不会缩短上游 Retry-After 指令。
func TestUpstreamBillingProbeUnsupportedDelayIsStretchedAndBounded(t *testing.T) {
	// 默认 30 分钟 interval：普通失败 24~36 分钟，unsupported 为其 8 倍。
	stretched := unsupportedProbeDelay(30, 0)
	require.Greater(t, stretched, 36*time.Minute)
	require.GreaterOrEqual(t, stretched, 192*time.Minute)
	require.LessOrEqual(t, stretched, 288*time.Minute)

	// 永不超过封顶值，因此 unsupported 账号不会被永久排除在重探之外。
	require.LessOrEqual(t, unsupportedProbeDelay(upstreamBillingProbeMaxIntervalMinutes, 0), upstreamBillingProbeMaxDelay)
	require.Positive(t, unsupportedProbeDelay(upstreamBillingProbeMinIntervalMinutes, 0))

	// Retry-After 更长时原样保留，不被封顶缩短；更短时至少不早于该指令。
	require.Equal(t, 48*time.Hour, unsupportedProbeDelay(30, 48*time.Hour))
	require.GreaterOrEqual(t, unsupportedProbeDelay(30, time.Hour), time.Hour)
placeholder

// 加长退避只把 unsupported 账号移出周期性热队列，手动探测不受影响。
func TestUpstreamBillingProbeUnsupportedBackoffDefersRunnerButNotManualProbe(t *testing.T) {
	// Ollama Cloud 形态：官方域，不发请求直接落 unsupported。
	account := &Account{
		ID:          31,
		Platform:    PlatformAnthropic,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Concurrency: 1,
placeholder"api_key": "sk-ollama", "base_url": "https://ollama.com/v1"placeholder,
		Extra:       map[string]any{UpstreamBillingProbeEnabledExtraKey: trueplaceholder,
placeholder
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{account.ID: accountplaceholderplaceholder
	upstream := &upstreamBillingProbeHTTPStub{placeholder
	settingsRepo := &upstreamBillingProbeSettingRepo{values: map[string]string{
		SettingKeyUpstreamBillingProbeSettings: `{"enabled":true,"interval_minutes":30placeholder`,
placeholderplaceholder
	svc := newUpstreamBillingProbeTestService(repo, upstream, settingsRepo)
	start := time.Date(2026, time.July, 26, 2, 0, 0, 0, time.UTC)
	now := start
	svc.now = func() time.Time { return now placeholder

	require.NoError(t, svc.RunDue(context.Background()))
	first := decodeUpstreamBillingProbeSnapshot(account.Extra)
	require.NotNil(t, first)
	require.Equal(t, UpstreamBillingProbeStatusUnsupported, first.Status)
	require.Equal(t, start, first.LastAttemptAt)
	require.False(t, first.NextProbeAt.Before(start.Add(192*time.Minute)))
	require.Zero(t, upstream.calls.Load())

	// 一个普通失败早就该重探的时间点（远超 36 分钟），runner 仍跳过该账号。
	now = start.Add(90 * time.Minute)
	require.NoError(t, svc.RunDue(context.Background()))
	deferred := decodeUpstreamBillingProbeSnapshot(account.Extra)
	require.NotNil(t, deferred)
	require.Equal(t, first.LastAttemptAt, deferred.LastAttemptAt)

	// 手动探测无视退避窗口，管理员随时可以重试。
	manual, err := svc.ProbeAccount(context.Background(), account.ID)
placeholder
	require.Equal(t, UpstreamBillingProbeStatusUnsupported, manual.Status)
	require.Equal(t, now, manual.LastAttemptAt)
	require.Equal(t, 2, manual.FailureCount)
placeholder

func TestUpstreamBillingProbeEmptyResponseIsPersistedAsFailure(t *testing.T) {
	account := &Account{
		ID:          21,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Concurrency: 1,
placeholder"api_key": "sk-test", "base_url": "https://upstream.example"placeholder,
placeholder
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{account.ID: accountplaceholderplaceholder
	svc := newUpstreamBillingProbeTestService(repo, &httpUpstreamRecorder{placeholder, &upstreamBillingProbeSettingRepo{placeholder)
	fixedNow := time.Date(2026, time.July, 13, 2, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return fixedNow placeholder

	snapshot, err := svc.ProbeAccount(context.Background(), account.ID)
placeholder
	require.Equal(t, UpstreamBillingProbeStatusFailed, snapshot.Status)
	require.Equal(t, "empty_response", snapshot.LastError)
	require.False(t, snapshot.NextProbeAt.Before(fixedNow.Add(24*time.Minute)))
	require.False(t, snapshot.NextProbeAt.After(fixedNow.Add(36*time.Minute)))

	snapshot, err = svc.ProbeAccount(context.Background(), account.ID)
placeholder
	require.Equal(t, 2, snapshot.FailureCount)
	require.False(t, snapshot.NextProbeAt.Before(fixedNow.Add(24*time.Minute)))
	require.False(t, snapshot.NextProbeAt.After(fixedNow.Add(36*time.Minute)))
placeholder

func TestUpstreamBillingProbeUnsupportedAndAccountToggle(t *testing.T) {
	account := &Account{
		ID:          19,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Concurrency: 1,
placeholder"api_key": "sk-test", "base_url": "https://upstream.example"placeholder,
placeholder
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{account.ID: accountplaceholderplaceholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusNotFound,
		Header:     http.Header{placeholder,
		Body:       io.NopCloser(strings.NewReader("not found")),
placeholderplaceholder
	svc := newUpstreamBillingProbeTestService(repo, upstream, &upstreamBillingProbeSettingRepo{placeholder)
	fixedNow := time.Date(2026, time.July, 13, 2, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return fixedNow placeholder

	require.NoError(t, svc.SetAccountEnabled(context.Background(), account.ID, true))
	require.Equal(t, true, account.Extra[UpstreamBillingProbeEnabledExtraKey])
	account.Extra[UpstreamBillingRateSyncEnabledExtraKey] = true
	snapshot, err := svc.ProbeAccount(context.Background(), account.ID)
placeholder
	require.Equal(t, UpstreamBillingProbeStatusUnsupported, snapshot.Status)
	require.Equal(t, "unsupported", snapshot.LastError)
	// unsupported 走加长退避：默认 30 分钟 interval ⇒ (24~36) * 8 分钟。
	require.False(t, snapshot.NextProbeAt.Before(fixedNow.Add(192*time.Minute)))
	require.False(t, snapshot.NextProbeAt.After(fixedNow.Add(288*time.Minute)))

	snapshot, err = svc.ProbeAccount(context.Background(), account.ID)
placeholder
	require.Equal(t, 2, snapshot.FailureCount)
	require.False(t, snapshot.NextProbeAt.Before(fixedNow.Add(192*time.Minute)))
	require.False(t, snapshot.NextProbeAt.After(fixedNow.Add(288*time.Minute)))
	require.NoError(t, svc.SetAccountEnabled(context.Background(), account.ID, false))
	require.Equal(t, false, account.Extra[UpstreamBillingProbeEnabledExtraKey])
	require.Equal(t, false, account.Extra[UpstreamBillingRateSyncEnabledExtraKey])

	invalid := &Account{ID: 20, Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder
	repo.accounts[invalid.ID] = invalid
	err = svc.SetAccountEnabled(context.Background(), invalid.ID, true)
	require.True(t, errors.Is(err, ErrUpstreamBillingProbeAccountInvalid))
placeholder

func TestUpstreamBillingProbeRunnerIsBoundedAndManualProbeIgnoresSwitches(t *testing.T) {
	accounts := make(map[int64]*Account, 25)
	for id := int64(1); id <= 25; id++ {
		accounts[id] = &Account{
			ID:          id,
			Platform:    PlatformOpenAI,
			Type:        AccountTypeAPIKey,
			Status:      StatusActive,
			Concurrency: 1,
	placeholder"api_key": "sk-test", "base_url": "https://upstream.example"placeholder,
			Extra:       map[string]any{UpstreamBillingProbeEnabledExtraKey: trueplaceholder,
	placeholder
placeholder
	repo := &upstreamBillingProbeAccountRepo{accounts: accountsplaceholder
	settingsRepo := &upstreamBillingProbeSettingRepo{values: map[string]string{
		SettingKeyUpstreamBillingProbeSettings: `{"enabled":true,"interval_minutes":30placeholder`,
placeholderplaceholder
	upstream := &upstreamBillingProbeHTTPStub{placeholder
	svc := newUpstreamBillingProbeTestService(repo, upstream, settingsRepo)
	svc.now = func() time.Time { return time.Date(2026, time.July, 13, 2, 0, 0, 0, time.UTC) placeholder

	require.NoError(t, svc.RunDue(context.Background()))
	require.Equal(t, int64(20), upstream.calls.Load())

	settingsRepo.mu.Lock()
	settingsRepo.values[SettingKeyUpstreamBillingProbeSettings] = `{"enabled":false,"interval_minutes":30placeholder`
	settingsRepo.mu.Unlock()
	require.NoError(t, svc.RunDue(context.Background()))
	require.Equal(t, int64(20), upstream.calls.Load())

	accounts[25].Extra[UpstreamBillingProbeEnabledExtraKey] = false
	manualRate := 0.25
	accounts[25].RateMultiplier = &manualRate
	snapshot, err := svc.ProbeAccount(context.Background(), 25)
placeholder
	require.Equal(t, UpstreamBillingProbeStatusOK, snapshot.Status)
	require.Equal(t, int64(21), upstream.calls.Load())
	require.NotNil(t, accounts[25].RateMultiplier)
	require.Equal(t, manualRate, *accounts[25].RateMultiplier)
placeholder

func TestUpstreamBillingProbeRunnerRechecksEnabledAfterDueSelection(t *testing.T) {
	account := &Account{
		ID:          26,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Concurrency: 1,
placeholder"api_key": "sk-test", "base_url": "https://upstream.example"placeholder,
		Extra:       map[string]any{UpstreamBillingProbeEnabledExtraKey: falseplaceholder,
placeholder
	staleDue := *account
	staleDue.Extra = map[string]any{UpstreamBillingProbeEnabledExtraKey: trueplaceholder
	baseRepo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{account.ID: accountplaceholderplaceholder
	repo := &staleDueUpstreamBillingProbeAccountRepo{upstreamBillingProbeAccountRepo: baseRepo, due: []Account{staleDueplaceholderplaceholder
	settingsRepo := &upstreamBillingProbeSettingRepo{values: map[string]string{
		SettingKeyUpstreamBillingProbeSettings: `{"enabled":true,"interval_minutes":30placeholder`,
placeholderplaceholder
	upstream := &upstreamBillingProbeHTTPStub{placeholder
	svc := newUpstreamBillingProbeTestService(repo, upstream, settingsRepo)

	require.NoError(t, svc.RunDue(context.Background()))
	require.Zero(t, upstream.calls.Load())
	require.NotContains(t, account.Extra, UpstreamBillingProbeExtraKey)
placeholder

func TestUpstreamBillingProbeNeverDowngradesMissingConfiguredProxyToDirect(t *testing.T) {
	proxyID := int64(7)
	for _, tc := range []struct {
		name       string
		proxy      *Proxy
		wantReason string
		wantErr    error
placeholder{
		{name: "missing hydrated proxy", wantReason: "proxy_unavailable"placeholder,
		{name: "mismatched hydrated proxy", proxy: &Proxy{ID: 8, Protocol: "http", Host: "127.0.0.1", Port: 8080placeholder, wantErr: ErrUpstreamBillingProbeIdentityChangedplaceholder,
placeholder {
		t.Run(tc.name, func(t *testing.T) {
			account := &Account{
				ID:          27,
				Platform:    PlatformOpenAI,
				Type:        AccountTypeAPIKey,
				Status:      StatusActive,
				Concurrency: 1,
		placeholder"api_key": "sk-sensitive", "base_url": "https://upstream.example"placeholder,
				ProxyID:     &proxyID,
				Proxy:       tc.proxy,
		placeholder
			repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{account.ID: accountplaceholderplaceholder
			upstream := &upstreamBillingProbeHTTPStub{placeholder
			svc := newUpstreamBillingProbeTestService(repo, upstream, &upstreamBillingProbeSettingRepo{placeholder)

			snapshot, err := svc.ProbeAccount(context.Background(), account.ID)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				require.Nil(t, snapshot)
		placeholder else {
			placeholder
				require.Equal(t, UpstreamBillingProbeStatusFailed, snapshot.Status)
				require.Equal(t, tc.wantReason, snapshot.LastError)
		placeholder
			require.Zero(t, upstream.calls.Load())
			if tc.wantErr != nil {
				require.NotContains(t, account.Extra, UpstreamBillingProbeExtraKey)
		placeholder
	placeholder)
placeholder
placeholder

func TestUpstreamBillingProbeRunnerOnlyScansOnLeader(t *testing.T) {
	account := &Account{
		ID:          31,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Concurrency: 1,
placeholder"api_key": "sk-test", "base_url": "https://upstream.example"placeholder,
		Extra:       map[string]any{UpstreamBillingProbeEnabledExtraKey: trueplaceholder,
placeholder
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{account.ID: accountplaceholderplaceholder
	upstream := &upstreamBillingProbeHTTPStub{placeholder
	cache := &fakeLeaderLockCache{placeholder
	lockKey := upstreamBillingProbeLeaderLockKeyAt(time.Now())
	peer := newUpstreamBillingProbeTestService(repo, upstream, &upstreamBillingProbeSettingRepo{placeholder)
	peer.instanceID = "peer"
	peer.SetLeaderLock(cache, nil)
	_, acquired, err := peer.tryAcquireLeaderLock(context.Background(), lockKey)
placeholder
	require.True(t, acquired)
	svc := newUpstreamBillingProbeTestService(repo, upstream, &upstreamBillingProbeSettingRepo{placeholder)
	svc.SetLeaderLock(cache, nil)

	require.NoError(t, svc.RunDue(context.Background()))
	require.Zero(t, upstream.calls.Load())

	require.NoError(t, cache.ReleaseLeaderLock(context.Background(), lockKey, "peer"))
	require.NoError(t, svc.RunDue(context.Background()))
	require.Equal(t, int64(1), upstream.calls.Load())
placeholder

func TestUpstreamBillingProbeLeaderLockFailsClosedOnCacheError(t *testing.T) {
	svc := newUpstreamBillingProbeTestService(&upstreamBillingProbeAccountRepo{placeholder, &upstreamBillingProbeHTTPStub{placeholder, &upstreamBillingProbeSettingRepo{placeholder)
	svc.SetLeaderLock(&fakeLeaderLockCache{acquireErr: context.DeadlineExceededplaceholder, nil)

	release, acquired, err := svc.tryAcquireLeaderLock(context.Background(), upstreamBillingProbeLeaderLockKeyAt(time.Now()))

	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.False(t, acquired)
	require.Nil(t, release)
placeholder

func TestUpstreamBillingProbeLeaderLockUsesCadenceBuckets(t *testing.T) {
	cache := &fakeLeaderLockCache{placeholder
	first := newUpstreamBillingProbeTestService(&upstreamBillingProbeAccountRepo{placeholder, &upstreamBillingProbeHTTPStub{placeholder, &upstreamBillingProbeSettingRepo{placeholder)
	second := newUpstreamBillingProbeTestService(&upstreamBillingProbeAccountRepo{placeholder, &upstreamBillingProbeHTTPStub{placeholder, &upstreamBillingProbeSettingRepo{placeholder)
	first.SetLeaderLock(cache, nil)
	second.SetLeaderLock(cache, nil)
	beforeBoundary := time.Unix(59, 0)
	afterBoundary := beforeBoundary.Add(time.Second)

	releaseFirst, acquired, err := first.tryAcquireLeaderLock(context.Background(), upstreamBillingProbeLeaderLockKeyAt(beforeBoundary))
placeholder
	require.True(t, acquired)
	releaseSecond, acquired, err := second.tryAcquireLeaderLock(context.Background(), upstreamBillingProbeLeaderLockKeyAt(afterBoundary))
placeholder
	require.True(t, acquired, "the prior cadence lock must not suppress the next cadence")
	releaseFirst()
	releaseSecond()
placeholder

func TestUpstreamBillingProbeFiveInstancesRunOneConcurrentBatch(t *testing.T) {
	account := &Account{
		ID:          32,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Concurrency: 1,
placeholder"api_key": "sk-test", "base_url": "http://127.0.0.1:8080"placeholder,
		Extra:       map[string]any{UpstreamBillingProbeEnabledExtraKey: trueplaceholder,
placeholder
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{account.ID: accountplaceholderplaceholder
	settingsRepo := &upstreamBillingProbeSettingRepo{values: map[string]string{
		SettingKeyUpstreamBillingProbeSettings: `{"enabled":true,"interval_minutes":30placeholder`,
placeholderplaceholder
	cache := &fakeLeaderLockCache{placeholder
	entered := make(chan struct{placeholder)
	unblock := make(chan struct{placeholder)
	var enteredOnce sync.Once
	upstream := &upstreamBillingProbeHTTPStub{beforeResponse: func() {
		enteredOnce.Do(func() { close(entered) placeholder)
		<-unblock
placeholderplaceholder

	start := make(chan struct{placeholder)
	results := make(chan error, 5)
	for range 5 {
		svc := newUpstreamBillingProbeTestService(repo, upstream, settingsRepo)
		svc.SetLeaderLock(cache, nil)
		go func() {
			<-start
			results <- svc.RunDue(context.Background())
	placeholder()
placeholder
	close(start)

	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("leader did not start the probe batch")
placeholder
	for range 4 {
		select {
		case err := <-results:
		placeholder
		case <-time.After(time.Second):
			t.Fatal("non-leader instance did not skip the active batch")
	placeholder
placeholder
	require.Equal(t, int64(1), upstream.calls.Load())
	close(unblock)
	require.NoError(t, <-results)
	require.Equal(t, int64(1), upstream.calls.Load())
placeholder

func TestUpstreamBillingProbeManualBatchesShareConcurrencyLimit(t *testing.T) {
	accounts := make(map[int64]*Account, 12)
	for id := int64(1); id <= 12; id++ {
		accounts[id] = &Account{
			ID:          id,
			Platform:    PlatformOpenAI,
			Type:        AccountTypeAPIKey,
			Status:      StatusActive,
			Concurrency: 1,
	placeholder"api_key": "sk-test", "base_url": "http://127.0.0.1:8080"placeholder,
	placeholder
placeholder
	repo := &upstreamBillingProbeAccountRepo{accounts: accountsplaceholder
	settingsRepo := &upstreamBillingProbeSettingRepo{values: map[string]string{
		SettingKeyUpstreamBillingProbeSettings: `{"enabled":true,"interval_minutes":30placeholder`,
placeholderplaceholder
	entered := make(chan struct{placeholder, len(accounts))
	unblock := make(chan struct{placeholder)
	var unblockOnce sync.Once
	release := func() { unblockOnce.Do(func() { close(unblock) placeholder) placeholder
	t.Cleanup(release)
	upstream := &upstreamBillingProbeHTTPStub{beforeResponse: func() {
		entered <- struct{placeholder{placeholder
		<-unblock
placeholderplaceholder
	svc := newUpstreamBillingProbeTestService(repo, upstream, settingsRepo)

	results := make(chan []UpstreamBillingProbeResult, 3)
	for batch := 0; batch < 3; batch++ {
		firstID := int64(batch*4 + 1)
		ids := []int64{firstID, firstID + 1, firstID + 2, firstID + 3placeholder
		go func() { results <- svc.ProbeAccounts(context.Background(), ids) placeholder()
placeholder
	for range upstreamBillingProbeConcurrency {
		select {
		case <-entered:
		case <-time.After(time.Second):
			t.Fatal("shared probe slots did not fill")
	placeholder
placeholder
	select {
	case <-entered:
		release()
		t.Fatal("parallel manual batches exceeded the service-wide concurrency limit")
	case <-time.After(100 * time.Millisecond):
placeholder
	release()

	for range 3 {
		select {
		case batchResults := <-results:
			for _, result := range batchResults {
				require.Empty(t, result.Error)
				require.NotNil(t, result.Snapshot)
		placeholder
		case <-time.After(time.Second):
			t.Fatal("manual probe batch did not finish")
	placeholder
placeholder
	require.Equal(t, int64(upstreamBillingProbeConcurrency), upstream.maxActive.Load())
placeholder

func TestUpstreamBillingProbeManualAndScheduledRequestsShareOneNetworkProbe(t *testing.T) {
	account := &Account{
		ID:          46,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Concurrency: 1,
placeholder"api_key": "sk-test", "base_url": "https://upstream.example"placeholder,
		Extra:       map[string]any{UpstreamBillingProbeEnabledExtraKey: trueplaceholder,
placeholder
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{account.ID: accountplaceholderplaceholder
	started := make(chan struct{placeholder)
	unblock := make(chan struct{placeholder)
	var startedOnce sync.Once
	upstream := &upstreamBillingProbeHTTPStub{beforeResponse: func() {
		startedOnce.Do(func() { close(started) placeholder)
		<-unblock
placeholderplaceholder
	svc := newUpstreamBillingProbeTestService(repo, upstream, &upstreamBillingProbeSettingRepo{placeholder)

	errs := make(chan error, 2)
	go func() {
		_, err := svc.probeScheduledAccount(context.Background(), account.ID, 30)
		errs <- err
placeholder()
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("scheduled probe did not reach the upstream")
placeholder
	manualStarted := make(chan struct{placeholder)
	go func() {
		close(manualStarted)
		_, err := svc.ProbeAccount(context.Background(), account.ID)
		errs <- err
placeholder()
	<-manualStarted
	time.Sleep(20 * time.Millisecond)
	close(unblock)
	require.NoError(t, <-errs)
	require.NoError(t, <-errs)
	require.Equal(t, int64(1), upstream.calls.Load())
placeholder

func TestUpstreamBillingProbeScheduledRechecksAfterWaitingForSlot(t *testing.T) {
	account := &Account{
		ID:          47,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Concurrency: 1,
placeholder"api_key": "sk-test", "base_url": "https://upstream.example"placeholder,
		Extra:       map[string]any{UpstreamBillingProbeEnabledExtraKey: trueplaceholder,
placeholder
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{account.ID: accountplaceholderplaceholder
	upstream := &upstreamBillingProbeHTTPStub{placeholder
	svc := newUpstreamBillingProbeTestService(repo, upstream, &upstreamBillingProbeSettingRepo{placeholder)
	for range upstreamBillingProbeConcurrency {
		svc.probeSlots <- struct{placeholder{placeholder
placeholder
	result := make(chan error, 1)
	go func() {
		_, err := svc.probeScheduledAccount(context.Background(), account.ID, 30)
		result <- err
placeholder()
	time.Sleep(20 * time.Millisecond)
	repo.mu.Lock()
	account.Extra[UpstreamBillingProbeEnabledExtraKey] = false
	repo.mu.Unlock()
	<-svc.probeSlots

	require.NoError(t, <-result)
	require.Zero(t, upstream.calls.Load())
placeholder

func TestUpstreamBillingProbeLeaderLockCoversStaggeredInstancesInCadenceWindow(t *testing.T) {
	account := func(id int64) *Account {
	placeholder
			ID:          id,
			Platform:    PlatformOpenAI,
			Type:        AccountTypeAPIKey,
			Status:      StatusActive,
			Concurrency: 1,
	placeholder"api_key": "sk-test", "base_url": "http://127.0.0.1:8080"placeholder,
			Extra:       map[string]any{UpstreamBillingProbeEnabledExtraKey: trueplaceholder,
	placeholder
placeholder
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{41: account(41)placeholderplaceholder
	settingsRepo := &upstreamBillingProbeSettingRepo{values: map[string]string{
		SettingKeyUpstreamBillingProbeSettings: `{"enabled":true,"interval_minutes":30placeholder`,
placeholderplaceholder
	cache := &fakeLeaderLockCache{placeholder
	upstream := &upstreamBillingProbeHTTPStub{placeholder
	first := newUpstreamBillingProbeTestService(repo, upstream, settingsRepo)
	first.SetLeaderLock(cache, nil)

	require.NoError(t, first.RunDue(context.Background()))
	require.Equal(t, int64(1), upstream.calls.Load())
	require.Equal(t, first.instanceID, cache.heldBy(upstreamBillingProbeLeaderLockKeyAt(time.Now())))

	repo.mu.Lock()
	repo.accounts[42] = account(42)
	repo.mu.Unlock()
	staggered := newUpstreamBillingProbeTestService(repo, upstream, settingsRepo)
	staggered.SetLeaderLock(cache, nil)
	require.NoError(t, staggered.RunDue(context.Background()))
	require.Equal(t, int64(1), upstream.calls.Load(), "a staggered instance must not start a second batch inside the cadence window")
placeholder
