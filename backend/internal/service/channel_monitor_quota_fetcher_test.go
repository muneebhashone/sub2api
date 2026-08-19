//go:build unit

package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/domain"
	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/stretchr/testify/require"
)

// --- fetcher 依赖 stub ---

type stubMonitorUsageSource struct {
	usage *UsageInfo
	err   error
	// block 非 nil 时 GetUsageForAccount 阻塞在该 channel 上，用于并发/singleflight 测试。
	block chan struct{placeholder

	mu          sync.Mutex
	calls       int
	lastCtx     context.Context
	lastAccount *Account
placeholder

func (s *stubMonitorUsageSource) GetUsageForAccount(ctx context.Context, account *Account, force ...bool) (*UsageInfo, error) {
	s.mu.Lock()
	s.calls++
	s.lastCtx = ctx
	s.lastAccount = account
	s.mu.Unlock()
	if s.block != nil {
		<-s.block
placeholder
	return s.usage, s.err
placeholder

func (s *stubMonitorUsageSource) getCalls() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.calls
placeholder

func (s *stubMonitorUsageSource) getLastAccount() *Account {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lastAccount
placeholder

type stubMonitorCNQuotaSource struct {
	result      *CNProviderQuotaProbeResult
	err         error
	calls       int
	lastAccount *Account
placeholder

func (s *stubMonitorCNQuotaSource) QueryUsageForAccount(ctx context.Context, account *Account) (*CNProviderQuotaProbeResult, error) {
	s.calls++
	s.lastAccount = account
	return s.result, s.err
placeholder

type stubMonitorCNBalanceSource struct {
	result      *CNProviderBalanceResult
	err         error
	calls       int
	lastAccount *Account
placeholder

func (s *stubMonitorCNBalanceSource) QueryBalanceForAccount(ctx context.Context, account *Account) (*CNProviderBalanceResult, error) {
	s.calls++
	s.lastAccount = account
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
		usage:            usage,
		cnQuota:          cnQuota,
		cnBalance:        cnBalance,
		accounts:         accounts,
		balanceThreshold: monitorBalanceThreshold(nil),
		cache:            make(map[int64]monitorQuotaCacheEntry),
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
	require.Equal(t, 1, usage.getCalls())
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
		Success:   true,
		Available: true,
		Balance:   12.34,
		Currency:  "CNY",
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
	require.False(t, snapshot.BalanceLow)
	require.Empty(t, snapshot.Error)
placeholder

// P2-6：fetchUncached 只 GetByID 一次，已加载的 account 指针直传数据源，
// 三条路由都不能让下游重载账号。
func TestQuotaFetcher_LoadsAccountOnceAndPassesItThrough(t *testing.T) {
	t.Run("overseas usage", func(t *testing.T) {
		fetcher, usage, _, _, accounts := newQuotaFetcherTestSetup(t)
		acc := &Account{ID: 21, Platform: domain.PlatformAnthropicplaceholder
		accounts.accounts[21] = acc
		usage.usage = &UsageInfo{placeholder

		fetcher.Fetch(context.Background(), 21)

		require.Equal(t, 1, accounts.calls)
		require.Same(t, acc, usage.getLastAccount())
		require.Equal(t, 1, usage.getCalls())
placeholder)

	t.Run("cn coding plan", func(t *testing.T) {
		fetcher, _, cnQuota, _, accounts := newQuotaFetcherTestSetup(t)
		acc := &Account{ID: 22, Platform: domain.PlatformKimi, Credentials: map[string]any{"account_mode": AccountModeCodingplaceholderplaceholder
		accounts.accounts[22] = acc
		cnQuota.result = &CNProviderQuotaProbeResult{Success: trueplaceholder

		fetcher.Fetch(context.Background(), 22)

		require.Equal(t, 1, accounts.calls)
		require.Same(t, acc, cnQuota.lastAccount)
		require.Equal(t, 1, cnQuota.calls)
placeholder)

	t.Run("cn payg", func(t *testing.T) {
		fetcher, _, _, cnBalance, accounts := newQuotaFetcherTestSetup(t)
		acc := &Account{ID: 23, Platform: domain.PlatformDeepseek, Credentials: map[string]any{"account_mode": AccountModePayGplaceholderplaceholder
		accounts.accounts[23] = acc
		cnBalance.result = &CNProviderBalanceResult{Success: true, Available: true, Balance: 1, Currency: "CNY"placeholder

		fetcher.Fetch(context.Background(), 23)

		require.Equal(t, 1, accounts.calls)
		require.Same(t, acc, cnBalance.lastAccount)
		require.Equal(t, 1, cnBalance.calls)
placeholder)
placeholder

// --- 失败路径（Fetch 永不返回 error） ---

func TestQuotaFetcher_AccountMissingYieldsLinkedAccountSnapshot(t *testing.T) {
	fetcher, usage, _, _, accounts := newQuotaFetcherTestSetup(t)
	accounts.err = errors.New("not found")

	snapshot := fetcher.Fetch(context.Background(), 404)

	require.False(t, snapshot.Success)
	require.Equal(t, "linked account not found", snapshot.Error)
	require.Equal(t, 0, usage.getCalls()) // 未走到数据源
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

// 值通道失败：antigravity/grok 等平台 err==nil 但错误降级在 UsageInfo 字段里，
// 必须识别为失败快照，否则会被误判为 operational。
func TestQuotaFetcher_UsageValueChannelFailureYieldsFailureSnapshot(t *testing.T) {
	fetcher, usage, _, _, accounts := newQuotaFetcherTestSetup(t)

	// 凭据失效（401 语义）→ failed。
	accounts.accounts[3] = &Account{ID: 3, Platform: domain.PlatformAnthropicplaceholder
	usage.usage = &UsageInfo{Error: "usage API error: HTTP 401", ErrorCode: errorCodeUnauthenticated, NeedsReauth: trueplaceholder
	snapshot := fetcher.Fetch(context.Background(), 3)
	require.False(t, snapshot.Success)
	require.True(t, snapshot.CredentialInvalid)
	require.Contains(t, snapshot.Error, "401")
	require.Equal(t, MonitorStatusFailed, deriveQuotaCheckResult(snapshot, "quota", time.Now()).Status)

	// 限流等非凭据失败 → error（而非 operational）。
	accounts.accounts[13] = &Account{ID: 13, Platform: domain.PlatformAnthropicplaceholder
	usage.usage = &UsageInfo{Error: "usage API error: HTTP 429", ErrorCode: errorCodeRateLimitedplaceholder
	snapshot = fetcher.Fetch(context.Background(), 13)
	require.False(t, snapshot.Success)
	require.False(t, snapshot.CredentialInvalid)
	require.Contains(t, snapshot.Error, "429")
	require.Equal(t, MonitorStatusError, deriveQuotaCheckResult(snapshot, "quota", time.Now()).Status)

	// grok 已知未知态（尚未观测到计费/限流头）不算失败。
	accounts.accounts[14] = &Account{ID: 14, Platform: domain.PlatformGrokplaceholder
	usage.usage = &UsageInfo{ErrorCode: "quota_unknown", Error: "Grok quota is unknown until billing is probed"placeholder
	snapshot = fetcher.Fetch(context.Background(), 14)
	require.True(t, snapshot.Success)
	require.Empty(t, snapshot.Error)
	require.Empty(t, snapshot.Tiers)
	require.Equal(t, MonitorStatusOperational, deriveQuotaCheckResult(snapshot, "quota", time.Now()).Status)
placeholder

func TestUsageFailureInfo_ClassificationMatrix(t *testing.T) {
	cases := []struct {
		name              string
		usage             *UsageInfo
		failed            bool
		credentialInvalid bool
		msg               string
placeholder{
		{name: "nil usage", usage: nilplaceholder,
		{name: "healthy empty", usage: &UsageInfo{placeholderplaceholder,
		{name: "error text only", usage: &UsageInfo{Error: "boom"placeholder, failed: true, msg: "boom"placeholder,
		{name: "needs reauth", usage: &UsageInfo{NeedsReauth: trueplaceholder, failed: true, credentialInvalid: true, msg: "usage fetch failed"placeholder,
		{name: "banned", usage: &UsageInfo{IsBanned: trueplaceholder, failed: true, credentialInvalid: true, msg: "usage fetch failed"placeholder,
		{name: "forbidden with reason", usage: &UsageInfo{IsForbidden: true, ForbiddenReason: "usage limited"placeholder, failed: true, credentialInvalid: true, msg: "usage limited"placeholder,
		{name: "error code unauthenticated", usage: &UsageInfo{ErrorCode: errorCodeUnauthenticatedplaceholder, failed: true, credentialInvalid: true, msg: errorCodeUnauthenticatedplaceholder,
		{name: "error code forbidden", usage: &UsageInfo{ErrorCode: errorCodeForbiddenplaceholder, failed: true, credentialInvalid: true, msg: errorCodeForbiddenplaceholder,
		{name: "error code rate limited", usage: &UsageInfo{ErrorCode: errorCodeRateLimitedplaceholder, failed: true, msg: errorCodeRateLimitedplaceholder,
		{name: "error code network error", usage: &UsageInfo{ErrorCode: errorCodeNetworkErrorplaceholder, failed: true, msg: errorCodeNetworkErrorplaceholder,
		{name: "grok quota unknown exempted", usage: &UsageInfo{ErrorCode: "quota_unknown", Error: "Grok quota is unknown until billing is probed"placeholderplaceholder,
placeholder
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			failed, credentialInvalid, msg := usageFailureInfo(tc.usage)
			require.Equal(t, tc.failed, failed)
			require.Equal(t, tc.credentialInvalid, credentialInvalid)
			if tc.msg != "" {
				require.Equal(t, tc.msg, msg)
		placeholder
	placeholder)
placeholder
placeholder

// 凭据失效只认 401/403（与 fetchCNBalance 口径一致）：CN quota 服务的
// CredentialValid 仅成功路径置 true，500/429/智谱业务错误须推导为 error 而非 failed。
func TestQuotaFetcher_CNQuotaCredentialInvalidByStatusCode(t *testing.T) {
	cases := []struct {
		name           string
		accountID      int64
		statusCode     int
		credentialBad  bool
		expectedStatus string
placeholder{
		{name: "401 unauthorized", accountID: 5, statusCode: 401, credentialBad: true, expectedStatus: MonitorStatusFailedplaceholder,
		{name: "403 forbidden", accountID: 15, statusCode: 403, credentialBad: true, expectedStatus: MonitorStatusFailedplaceholder,
		{name: "500 server error", accountID: 16, statusCode: 500, expectedStatus: MonitorStatusErrorplaceholder,
		{name: "429 rate limited", accountID: 17, statusCode: 429, expectedStatus: MonitorStatusErrorplaceholder,
		// 智谱 2xx 但业务级失败：StatusCode=200，非凭据问题。
		{name: "200 business error", accountID: 18, statusCode: 200, expectedStatus: MonitorStatusErrorplaceholder,
placeholder
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fetcher, _, cnQuota, _, accounts := newQuotaFetcherTestSetup(t)
			accounts.accounts[tc.accountID] = &Account{
				ID:          tc.accountID,
				Platform:    domain.PlatformZhipu,
		placeholder"account_mode": AccountModeCodingplaceholder,
		placeholder
			cnQuota.result = &CNProviderQuotaProbeResult{
				Success:    false,
				StatusCode: tc.statusCode,
				Error:      "api key expired",
		placeholder

			snapshot := fetcher.Fetch(context.Background(), tc.accountID)

			require.False(t, snapshot.Success)
			require.Equal(t, tc.credentialBad, snapshot.CredentialInvalid)
			require.Equal(t, tc.expectedStatus, deriveQuotaCheckResult(snapshot, "quota", time.Now()).Status)
	placeholder)
placeholder
placeholder

func TestQuotaFetcher_CNBalanceHTTP403MarksCredentialInvalid(t *testing.T) {
	fetcher, _, _, cnBalance, accounts := newQuotaFetcherTestSetup(t)
	accounts.accounts[6] = &Account{ID: 6, Platform: domain.PlatformKimiplaceholder
	cnBalance.result = &CNProviderBalanceResult{Success: false, StatusCode: 403, Error: "forbidden"placeholder

	snapshot := fetcher.Fetch(context.Background(), 6)

	require.False(t, snapshot.Success)
	require.True(t, snapshot.CredentialInvalid)
placeholder

// 余额告警口径与账号停调（CNProviderBalanceCheckService.checkOne）一致：
// 上游标记不可用或全部币种低于阈值 → BalanceLow → degraded；任一币种达标即健康。
func TestQuotaFetcher_CNBalanceLowMarksDegraded(t *testing.T) {
	cases := []struct {
		name        string
		accountID   int64
		result      *CNProviderBalanceResult
		balanceLow  bool
		wantStatus  string
		wantMessage string
placeholder{
		{
			// 审查例：余额 5/阈值 10 的账号调度器已停调，监控不能仍绿灯。
			name:       "balance below threshold",
			accountID:  21,
			result:     &CNProviderBalanceResult{Success: true, Available: true, Balance: 5, Currency: "CNY"placeholder,
			balanceLow: true,
			wantStatus: MonitorStatusDegraded, wantMessage: "balance low: 5 CNY",
	placeholder,
		{
			name:       "upstream marked unavailable",
			accountID:  22,
			result:     &CNProviderBalanceResult{Success: true, Available: false, Balance: 20, Currency: "CNY"placeholder,
			balanceLow: true,
			wantStatus: MonitorStatusDegraded, wantMessage: "balance low: 20 CNY",
	placeholder,
		{
			// deepseek 双币种：任一币种（USD 20）达标即健康。
			name:      "any currency above threshold is healthy",
			accountID: 23,
			result: &CNProviderBalanceResult{
				Success: true, Available: true, Balance: 5, Currency: "CNY",
				Balances: []CNProviderBalanceEntry{{Currency: "CNY", Balance: 5placeholder, {Currency: "USD", Balance: 20placeholderplaceholder,
		placeholder,
			wantStatus: MonitorStatusOperational,
	placeholder,
		{
			name:       "single currency above threshold",
			accountID:  24,
			result:     &CNProviderBalanceResult{Success: true, Available: true, Balance: 20, Currency: "CNY"placeholder,
			wantStatus: MonitorStatusOperational,
	placeholder,
placeholder
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fetcher, _, _, cnBalance, accounts := newQuotaFetcherTestSetup(t)
			fetcher.balanceThreshold = 10
			accounts.accounts[tc.accountID] = &Account{
				ID:          tc.accountID,
				Platform:    domain.PlatformKimi,
		placeholder"account_mode": AccountModePayGplaceholder,
		placeholder
			cnBalance.result = tc.result

			snapshot := fetcher.Fetch(context.Background(), tc.accountID)

			require.True(t, snapshot.Success)
			require.Equal(t, tc.balanceLow, snapshot.BalanceLow)
			res := deriveQuotaCheckResult(snapshot, "quota", time.Now())
			require.Equal(t, tc.wantStatus, res.Status)
			if tc.wantMessage != "" {
				require.Contains(t, res.Message, tc.wantMessage)
		placeholder else {
				require.Empty(t, res.Message)
		placeholder
	placeholder)
placeholder
placeholder

func TestNewChannelMonitorQuotaFetcher_ThresholdFromConfig(t *testing.T) {
	require.InDelta(t, 0.5, NewChannelMonitorQuotaFetcher(nil, nil, nil, nil, nil).balanceThreshold, 0.0001)

	cfg10 := &config.Config{Gateway: config.GatewayConfig{CNProviders: config.GatewayCNProvidersConfig{BalanceThreshold: 10placeholderplaceholderplaceholder
	require.InDelta(t, 10, NewChannelMonitorQuotaFetcher(nil, nil, nil, nil, cfg10).balanceThreshold, 0.0001)

	// 非正值（含显式 0）回退默认，避免 0 阈值下「余额=0 也不告警」。
	cfg0 := &config.Config{Gateway: config.GatewayConfig{CNProviders: config.GatewayCNProvidersConfig{BalanceThreshold: 0placeholderplaceholderplaceholder
	require.InDelta(t, 0.5, NewChannelMonitorQuotaFetcher(nil, nil, nil, nil, cfg0).balanceThreshold, 0.0001)
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
	require.Equal(t, 1, usage.getCalls(), "success snapshots should be served from cache")

	// 缓存过期后重新拉取。
	fetcher.mu.Lock()
	entry := fetcher.cache[8]
	entry.expiry = time.Now().Add(-time.Second)
	fetcher.cache[8] = entry
	fetcher.mu.Unlock()

	_ = fetcher.Fetch(context.Background(), 8)
	require.Equal(t, 2, usage.getCalls())
placeholder

func TestQuotaFetcher_CachesFailureSnapshotWithShortTTL(t *testing.T) {
	fetcher, usage, _, _, accounts := newQuotaFetcherTestSetup(t)
	accounts.accounts[4] = &Account{ID: 4, Platform: domain.PlatformOpenAIplaceholder
	usage.err = errors.New("boom")

	for i := 0; i < 2; i++ {
		snapshot := fetcher.Fetch(context.Background(), 4)
		require.False(t, snapshot.Success)
placeholder
	require.Equal(t, 1, usage.getCalls(), "failure snapshots should be served from the short negative cache")

	// 失败快照的 TTL 是负缓存时长（而非成功 TTL）。
	fetcher.mu.Lock()
	entry := fetcher.cache[4]
	require.WithinDuration(t, entry.snapshot.FetchedAt.Add(monitorQuotaErrorCacheTTL), entry.expiry, time.Second)
	entry.expiry = time.Now().Add(-time.Second)
	fetcher.cache[4] = entry
	fetcher.mu.Unlock()

	_ = fetcher.Fetch(context.Background(), 4)
	require.Equal(t, 2, usage.getCalls(), "expired negative cache should refetch")
placeholder

func TestQuotaFetcher_ConcurrentFetchesShareSingleFlight(t *testing.T) {
	fetcher, usage, _, _, accounts := newQuotaFetcherTestSetup(t)
	accounts.accounts[12] = &Account{ID: 12, Platform: domain.PlatformOpenAIplaceholder
	usage.usage = &UsageInfo{FiveHour: &UsageProgress{Utilization: 10placeholderplaceholder
	usage.block = make(chan struct{placeholder)

	var wg sync.WaitGroup
	snapshots := make([]*domain.MonitorQuotaSnapshot, 5)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			snapshots[idx] = fetcher.Fetch(context.Background(), 12)
	placeholder(i)
placeholder

	// 上游被 block 卡住时，5 个并发 Fetch 应只产生 1 次真实查询。
	require.Eventually(t, func() bool { return usage.getCalls() == 1 placeholder,
		5*time.Second, 10*time.Millisecond)
	close(usage.block)
	wg.Wait()

	for _, snapshot := range snapshots {
		require.NotNil(t, snapshot)
		require.True(t, snapshot.Success)
placeholder
	require.Equal(t, 1, usage.getCalls())

	// 成功快照已缓存：再取一次仍不打上游。
	_ = fetcher.Fetch(context.Background(), 12)
	require.Equal(t, 1, usage.getCalls())
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
	lowBalance := &domain.MonitorQuotaSnapshot{Success: true, BalanceLow: true, Balance: &balance, Currency: "CNY"placeholder
	res = deriveQuotaCheckResult(lowBalance, "quota", now)
	require.Equal(t, MonitorStatusDegraded, res.Status)
	require.Contains(t, res.Message, "balance low")

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
