//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/domain"
	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/stretchr/testify/require"
)

// --- fetcher 依赖 stub ---

type stubMonitorUsageSource struct {
	usage   *UsageInfo
	err     error
	calls   int
	lastCtx context.Context
placeholder

func (s *stubMonitorUsageSource) GetUsage(ctx context.Context, accountID int64, force ...bool) (*UsageInfo, error) {
	s.calls++
	s.lastCtx = ctx
	return s.usage, s.err
placeholder

type stubMonitorCNQuotaSource struct {
	result *CNProviderQuotaProbeResult
	err    error
	calls  int
placeholder

func (s *stubMonitorCNQuotaSource) QueryUsage(ctx context.Context, accountID int64) (*CNProviderQuotaProbeResult, error) {
	s.calls++
	return s.result, s.err
placeholder

type stubMonitorCNBalanceSource struct {
	result *CNProviderBalanceResult
	err    error
	calls  int
placeholder

func (s *stubMonitorCNBalanceSource) QueryBalance(ctx context.Context, accountID int64) (*CNProviderBalanceResult, error) {
	s.calls++
	return s.result, s.err
placeholder

type stubMonitorAccountSource struct {
	accounts map[int64]*Account
	err      error
	calls    int
placeholder

func (s *stubMonitorAccountSource) GetByID(ctx context.Context, id int64) (*Account, error) {
	s.calls++
	if s.err != nil {
		return nil, s.err
placeholder
	return s.accounts[id], nil
placeholder

func newQuotaFetcherTestSetup(t *testing.T) (*ChannelMonitorQuotaFetcher, *stubMonitorUsageSource, *stubMonitorCNQuotaSource, *stubMonitorCNBalanceSource, *stubMonitorAccountSource) {
placeholder
	usage := &stubMonitorUsageSource{placeholder
	cnQuota := &stubMonitorCNQuotaSource{placeholder
	cnBalance := &stubMonitorCNBalanceSource{placeholder
	accounts := &stubMonitorAccountSource{accounts: make(map[int64]*Account)placeholder
	fetcher := &ChannelMonitorQuotaFetcher{
		usage:     usage,
		cnQuota:   cnQuota,
		cnBalance: cnBalance,
		accounts:  accounts,
		cache:     make(map[int64]monitorQuotaCacheEntry),
placeholder
	return fetcher, usage, cnQuota, cnBalance, accounts
placeholder

// --- 分派 ---

func TestQuotaFetcher_OverseasAccountUsesUsageService(t *testing.T) {
	fetcher, usage, _, cnQuota, accounts := newQuotaFetcherTestSetup(t)
	accounts.accounts[7] = &Account{ID: 7, Platform: domain.PlatformAnthropicplaceholder
	resets := time.Now().Add(2 * time.Hour).UTC()
	usage.usage = &UsageInfo{
		FiveHour:         &UsageProgress{Utilization: 42.5, UsedRequests: 17, LimitRequests: 40, ResetsAt: &resetsplaceholder,
		SevenDay:         &UsageProgress{Utilization: 10placeholder,
		SubscriptionTier: "PRO",
placeholder

	snapshot := fetcher.Fetch(context.Background(), 7)

	require.True(t, snapshot.Success)
	require.Equal(t, "usage", snapshot.Source)
	require.Equal(t, "PRO", snapshot.PlanLevel)
	require.False(t, snapshot.CredentialInvalid)
	require.Empty(t, snapshot.Error)
	require.Len(t, snapshot.Tiers, 2)

	fiveHour := snapshot.Tiers[0]
	require.Equal(t, "5h", fiveHour.Window)
	require.Empty(t, fiveHour.Label)
	require.InDelta(t, 42.5, fiveHour.UsedPercent, 0.001)
	require.Equal(t, float64(17), fiveHour.Used)
	require.Equal(t, float64(40), fiveHour.Limit)
	require.NotEmpty(t, fiveHour.ResetAt)

	require.Equal(t, "7d", snapshot.Tiers[1].Window)
	require.Equal(t, 1, usage.calls)
	require.Equal(t, 0, cnQuota.calls)
placeholder

func TestQuotaFetcher_CodingPlanAccountUsesCNQuota(t *testing.T) {
	fetcher, _, cnQuota, cnBalance, accounts := newQuotaFetcherTestSetup(t)
	accounts.accounts[9] = &Account{
		ID:          9,
		Platform:    domain.PlatformKimi,
placeholder"account_mode": AccountModeCodingplaceholder,
placeholder
	cnQuota.result = &CNProviderQuotaProbeResult{
		Success:         true,
		CredentialValid: true,
		PlanLevel:       "",
		Tiers: []CNQuotaTier{
			{Window: "5h", UsedPercent: 33.3, ResetAt: "2026-08-18T06:00:00Z"placeholder,
			{Window: "weekly", UsedPercent: 12placeholder,
	placeholder,
placeholder

	snapshot := fetcher.Fetch(context.Background(), 9)

	require.True(t, snapshot.Success)
	require.Equal(t, "cn_quota", snapshot.Source)
	require.Len(t, snapshot.Tiers, 2)
	require.Equal(t, "5h", snapshot.Tiers[0].Window)
	require.InDelta(t, 33.3, snapshot.Tiers[0].UsedPercent, 0.001)
	require.Equal(t, "weekly", snapshot.Tiers[1].Window)
	require.Equal(t, 1, cnQuota.calls)
	require.Equal(t, 0, cnBalance.calls)
placeholder

func TestQuotaFetcher_PayGAccountUsesCNBalance(t *testing.T) {
	fetcher, _, _, cnBalance, accounts := newQuotaFetcherTestSetup(t)
	accounts.accounts[11] = &Account{
		ID:          11,
		Platform:    domain.PlatformDeepseek,
placeholder"account_mode": AccountModePayGplaceholder,
placeholder
	cnBalance.result = &CNProviderBalanceResult{
		Success:  true,
		Balance:  12.34,
		Currency: "CNY",
		Balances: []CNProviderBalanceEntry{
			{Currency: "CNY", Balance: 12.34placeholder,
			{Currency: "USD", Balance: placeholder,
	placeholder,
placeholder

	snapshot := fetcher.Fetch(context.Background(), 11)

	require.True(t, snapshot.Success)
	require.Equal(t, "cn_balance", snapshot.Source)
	require.NotNil(t, snapshot.Balance)
	require.InDelta(t, 12.34, *snapshot.Balance, 0.001)
	require.Equal(t, "CNY", snapshot.Currency)
	require.Len(t, snapshot.Balances, 2)
	require.Equal(t, "USD", snapshot.Balances[1].Currency)
	require.Empty(t, snapshot.Error)
placeholder

// --- 失败路径（Fetch 永不返回 error） ---

func TestQuotaFetcher_AccountMissingYieldsLinkedAccountSnapshot(t *testing.T) {
	fetcher, usage, _, _, accounts := newQuotaFetcherTestSetup(t)
	accounts.err = errors.New("not found")

	snapshot := fetcher.Fetch(context.Background(), 404)

	require.False(t, snapshot.Success)
	require.Equal(t, "linked account not found", snapshot.Error)
	require.Equal(t, 0, usage.calls) // 未走到数据源
placeholder

func TestQuotaFetcher_UsageAuthErrorMarksCredentialInvalid(t *testing.T) {
	fetcher, usage, _, _, accounts := newQuotaFetcherTestSetup(t)
	accounts.accounts[3] = &Account{ID: 3, Platform: domain.PlatformOpenAIplaceholder
	usage.err = errors.New("API returned 401: unauthorized")

	snapshot := fetcher.Fetch(context.Background(), 3)

	require.False(t, snapshot.Success)
	require.True(t, snapshot.CredentialInvalid)
	require.Contains(t, snapshot.Error, "401")
placeholder

func TestQuotaFetcher_CNQuotaCredentialInvalidFlagPropagates(t *testing.T) {
	fetcher, _, cnQuota, _, accounts := newQuotaFetcherTestSetup(t)
	accounts.accounts[5] = &Account{
		ID:          5,
		Platform:    domain.PlatformZhipu,
placeholder"account_mode": AccountModeCodingplaceholder,
placeholder
	cnQuota.result = &CNProviderQuotaProbeResult{Success: false, CredentialValid: false, Error: "api key expired"placeholder

	snapshot := fetcher.Fetch(context.Background(), 5)

	require.False(t, snapshot.Success)
	require.True(t, snapshot.CredentialInvalid)
	require.Equal(t, "api key expired", snapshot.Error)
placeholder

func TestQuotaFetcher_CNBalanceHTTP403MarksCredentialInvalid(t *testing.T) {
	fetcher, _, _, cnBalance, accounts := newQuotaFetcherTestSetup(t)
	accounts.accounts[6] = &Account{ID: 6, Platform: domain.PlatformKimiplaceholder
	cnBalance.result = &CNProviderBalanceResult{Success: false, StatusCode: 403, Error: "forbidden"placeholder

	snapshot := fetcher.Fetch(context.Background(), 6)

	require.False(t, snapshot.Success)
	require.True(t, snapshot.CredentialInvalid)
placeholder

func TestQuotaFetcher_NilDependenciesProduceErrorSnapshots(t *testing.T) {
	// fetcher 本体为 nil：直接降级为错误快照，不 panic。
	var nilFetcher *ChannelMonitorQuotaFetcher
	snapshot := nilFetcher.Fetch(context.Background(), 1)
	require.False(t, snapshot.Success)
	require.Equal(t, "quota fetcher is not configured", snapshot.Error)

	// 数据源缺失：账号能加载，但对应服务未注入。
	fetcher, _, _, _, accounts := newQuotaFetcherTestSetup(t)
	fetcher.usage = nil
	accounts.accounts[2] = &Account{ID: 2, Platform: domain.PlatformOpenAIplaceholder
	snapshot = fetcher.Fetch(context.Background(), 2)
	require.False(t, snapshot.Success)
	require.Contains(t, snapshot.Error, "not configured")
placeholder

// --- TTL 缓存 ---

func TestQuotaFetcher_CachesSuccessSnapshotPerAccount(t *testing.T) {
	fetcher, usage, _, _, accounts := newQuotaFetcherTestSetup(t)
	accounts.accounts[8] = &Account{ID: 8, Platform: domain.PlatformOpenAIplaceholder
	usage.usage = &UsageInfo{FiveHour: &UsageProgress{Utilization: 10placeholderplaceholder

	for i := 0; i < 3; i++ {
		snapshot := fetcher.Fetch(context.Background(), 8)
		require.True(t, snapshot.Success)
placeholder
	require.Equal(t, 1, usage.calls, "success snapshots should be served from cache")

	// 缓存过期后重新拉取。
	fetcher.mu.Lock()
	entry := fetcher.cache[8]
	entry.expiry = time.Now().Add(-time.Second)
	fetcher.cache[8] = entry
	fetcher.mu.Unlock()

	_ = fetcher.Fetch(context.Background(), 8)
	require.Equal(t, 2, usage.calls)
placeholder

func TestQuotaFetcher_DoesNotCacheFailures(t *testing.T) {
	fetcher, usage, _, _, accounts := newQuotaFetcherTestSetup(t)
	accounts.accounts[4] = &Account{ID: 4, Platform: domain.PlatformOpenAIplaceholder
	usage.err = errors.New("boom")

	_ = fetcher.Fetch(context.Background(), 4)
	_ = fetcher.Fetch(context.Background(), 4)

	require.Equal(t, 2, usage.calls, "failed snapshots must not be cached")
placeholder

// --- UsageInfo → tiers 归一 ---

func TestUsageQuotaTiers_MapsAllWindowKinds(t *testing.T) {
	limit := int64(1000)
	remaining := int64(400)
	resetUnix := int64(1777283883)
	usage := &UsageInfo{
		FiveHour:          &UsageProgress{Utilization: 50placeholder,
		SevenDay:          &UsageProgress{Utilization: 60placeholder,
		SevenDaySonnet:    &UsageProgress{Utilization: 70placeholder,
		SevenDayFable:     &UsageProgress{Utilization: 80placeholder,
		ThirtyDay:         &UsageProgress{Utilization: 20placeholder,
		GeminiSharedDaily: &UsageProgress{Utilization: 11placeholder,
		GeminiProDaily:    &UsageProgress{Utilization: 22placeholder,
		GeminiFlashDaily:  &UsageProgress{Utilization: 33placeholder,
		GrokRequestQuota:  &xai.QuotaWindow{Limit: &limit, Remaining: &remaining, ResetUnix: &resetUnixplaceholder,
		GrokTokenQuota:    &xai.QuotaWindow{Limit: &limit, Remaining: &remaining, ResetAt: "2026-08-19T00:00:00Z"placeholder,
		AntigravityQuota: map[string]*AntigravityModelQuota{
			"gemini-3-pro":   {Utilization: 45placeholder,
			"gemini-3-flash": {Utilization: 55placeholder,
	placeholder,
placeholder

	tiers := usageQuotaTiers(usage)

	// 5h/7d/7d-sonnet/7d-fable/30d + gemini×3 + grok×2 + antigravity×2
	require.Len(t, tiers, 12)

	byKey := make(map[string]domain.MonitorQuotaTier, len(tiers))
	for _, tier := range tiers {
		key := tier.Window
		if tier.Label != "" {
			key = tier.Window + "/" + tier.Label
	placeholder
		byKey[key] = tier
placeholder

	require.Contains(t, byKey, "5h")
	require.Contains(t, byKey, "7d")
	require.Contains(t, byKey, "7d-sonnet")
	require.Contains(t, byKey, "7d-fable")
	require.Contains(t, byKey, "30d")
	require.Contains(t, byKey, "daily/shared")
	require.Contains(t, byKey, "daily/pro")
	require.Contains(t, byKey, "daily/flash")
	require.Contains(t, byKey, "daily/requests")
	require.Contains(t, byKey, "daily/tokens")
	require.Contains(t, byKey, "total/gemini-3-pro")
	require.Contains(t, byKey, "total/gemini-3-flash")

	// grok requests 窗口：used = limit - remaining，百分比 60%。
	requests := byKey["daily/requests"]
	require.Equal(t, float64(600), requests.Used)
	require.Equal(t, float64(1000), requests.Limit)
	require.InDelta(t, 60.0, requests.UsedPercent, 0.001)
	require.NotEmpty(t, requests.ResetAt, "ResetUnix should fall back to RFC3339")

	tokens := byKey["daily/tokens"]
	require.Equal(t, "2026-08-19T00:00:00Z", tokens.ResetAt)
placeholder

func TestUsageQuotaTiers_NilAndEmptyInputs(t *testing.T) {
	require.Nil(t, usageQuotaTiers(nil))
	require.Nil(t, usageQuotaTiers(&UsageInfo{placeholder))

	// Grok 窗口 limit<=0 时跳过，避免除零。
	var zero int64
	tiers := usageQuotaTiers(&UsageInfo{
		GrokRequestQuota: &xai.QuotaWindow{Limit: &zero, Remaining: &zeroplaceholder,
placeholder)
	require.Nil(t, tiers)
placeholder

// --- 状态推导 ---

func TestDeriveQuotaCheckResult_StatusMatrix(t *testing.T) {
	now := time.Now()

	healthy := &domain.MonitorQuotaSnapshot{Success: true, Tiers: []domain.MonitorQuotaTier{{Window: "5h", UsedPercent: 40placeholderplaceholderplaceholder
	res := deriveQuotaCheckResult(healthy, "quota", now)
	require.Equal(t, MonitorStatusOperational, res.Status)
	require.Equal(t, "quota", res.Model)
	require.Empty(t, res.Message)

	highUsage := &domain.MonitorQuotaSnapshot{Success: true, Tiers: []domain.MonitorQuotaTier{
		{Window: "5h", UsedPercent: 30placeholder,
		{Window: "daily", Label: "pro", UsedPercent: 95placeholder,
placeholderplaceholder
	res = deriveQuotaCheckResult(highUsage, "quota", now)
	require.Equal(t, MonitorStatusDegraded, res.Status)
	require.Contains(t, res.Message, "pro/daily")
	require.Contains(t, res.Message, "95.0%")

	balance := -0.5
	depleted := &domain.MonitorQuotaSnapshot{Success: true, Balance: &balance, Currency: "CNY"placeholder
	res = deriveQuotaCheckResult(depleted, "quota", now)
	require.Equal(t, MonitorStatusDegraded, res.Status)
	require.Contains(t, res.Message, "balance depleted")

	invalid := &domain.MonitorQuotaSnapshot{Success: false, CredentialInvalid: true, Error: "401 unauthorized"placeholder
	res = deriveQuotaCheckResult(invalid, "quota", now)
	require.Equal(t, MonitorStatusFailed, res.Status)

	unlinked := &domain.MonitorQuotaSnapshot{Success: false, Error: "linked account not found"placeholder
	res = deriveQuotaCheckResult(unlinked, "quota", now)
	require.Equal(t, MonitorStatusDegraded, res.Status)

	other := &domain.MonitorQuotaSnapshot{Success: false, Error: "connection refused"placeholder
	res = deriveQuotaCheckResult(other, "quota", now)
	require.Equal(t, MonitorStatusError, res.Status)
	require.Equal(t, "connection refused", res.Message)

	res = deriveQuotaCheckResult(nil, "quota", now)
	require.Equal(t, MonitorStatusError, res.Status)
placeholder
