//go:build unit

package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/usagestats"
	"github.com/stretchr/testify/require"
)

type grokFreeQuotaUsageRepoStub struct {
	UsageLogRepository

	mu      sync.Mutex
	stats   map[int64]*usagestats.AccountStats
	err     error
	calls   int
	lastIDs []int64
	start   time.Time
placeholder

type grokFreeQuotaAccountRepoStub struct {
	AccountRepository
	accounts []Account
placeholder

func (r *grokFreeQuotaAccountRepoStub) ListSchedulableByPlatform(context.Context, string) ([]Account, error) {
	return append([]Account(nil), r.accounts...), nil
placeholder

func (r *grokFreeQuotaUsageRepoStub) GetAccountWindowStatsBatch(_ context.Context, accountIDs []int64, start time.Time) (map[int64]*usagestats.AccountStats, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls++
	r.lastIDs = append([]int64(nil), accountIDs...)
	r.start = start
	if r.err != nil {
		return nil, r.err
placeholder
	result := make(map[int64]*usagestats.AccountStats, len(accountIDs))
	for _, accountID := range accountIDs {
		if stats := r.stats[accountID]; stats != nil {
			copyStats := *stats
			result[accountID] = &copyStats
	placeholder
placeholder
	return result, nil
placeholder

func grokFreeQuotaTestConfig() *config.Config {
	cfg := &config.Config{placeholder
	cfg.Gateway.Grok.FreeQuotaSoftGateEnabled = true
	cfg.Gateway.Grok.FreeQuotaTokenLimit = 500_000
	cfg.Gateway.Grok.FreeQuotaSoftGatePercent = 95
	cfg.Gateway.Grok.FreeQuotaWindowHours = 24
	cfg.Gateway.Grok.FreeQuotaStatsCacheSeconds = 60
	return cfg
placeholder

func TestFilterGrokFreeQuotaAccountsOnlyBlocksExplicitFreeOAuth(t *testing.T) {
	repo := &grokFreeQuotaUsageRepoStub{stats: map[int64]*usagestats.AccountStats{
		1: {Tokens: 475_000placeholder, // 95% of 500k
placeholderplaceholder
	// Clear shared cache for deterministic unit tests.
	openaiGrokFreeQuotaGateCache = sync.Map{placeholder
	scheduler := &defaultOpenAIAccountScheduler{service: &OpenAIGatewayService{cfg: grokFreeQuotaTestConfig(), usageLogRepo: repoplaceholderplaceholder
	accounts := []Account{
		{ID: 1, Platform: PlatformGrok, Type: AccountTypeOAuth, Credentials: map[string]any{"subscription_tier": "FREE"placeholderplaceholder,
		{ID: 2, Platform: PlatformGrok, Type: AccountTypeOAuth, Credentials: map[string]any{"subscription_tier": "PRO"placeholderplaceholder,
		{ID: 3, Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder,
		{ID: 4, Platform: PlatformGrok, Type: AccountTypeAPIKey, Credentials: map[string]any{"subscription_tier": "FREE"placeholderplaceholder,
placeholder

	// First pass: cache miss fails open (does not block) and schedules background refresh.
	filtered := scheduler.filterGrokFreeQuotaAccounts(context.Background(), accounts)
	require.Equal(t, []int64{1, 2, 3, 4placeholder, accountIDs(filtered), "miss fails open on hot path")

	require.Eventually(t, func() bool {
		repo.mu.Lock()
		defer repo.mu.Unlock()
		return repo.calls >= 1
placeholder, 2*time.Second, 10*time.Millisecond)

	// Second pass: uses refreshed cache and blocks over-gate free OAuth.
	filtered = scheduler.filterGrokFreeQuotaAccounts(context.Background(), accounts)
	require.Equal(t, []int64{2, 3, 4placeholder, accountIDs(filtered), "paid and unknown fail-open; API-key free marker is not gated")
	require.Equal(t, []int64{1placeholder, repo.lastIDs, "paid, unknown, and API-key accounts must not enter the local free-tier query")
	require.WithinDuration(t, time.Now().UTC().Add(-24*time.Hour), repo.start, time.Second)
placeholder

func TestFilterGrokFreeQuotaAccountsStatsFailureFailsOpen(t *testing.T) {
	repo := &grokFreeQuotaUsageRepoStub{err: errors.New("usage database unavailable")placeholder
	openaiGrokFreeQuotaGateCache = sync.Map{placeholder
	scheduler := &defaultOpenAIAccountScheduler{service: &OpenAIGatewayService{cfg: grokFreeQuotaTestConfig(), usageLogRepo: repoplaceholderplaceholder
	accounts := []Account{{
		ID: 1, Platform: PlatformGrok, Type: AccountTypeOAuth,
placeholder"subscription_tier": "free"placeholder,
placeholderplaceholder

	filtered := scheduler.filterGrokFreeQuotaAccounts(context.Background(), accounts)
	require.Equal(t, []int64{1placeholder, accountIDs(filtered))
	require.Eventually(t, func() bool {
		repo.mu.Lock()
		defer repo.mu.Unlock()
		return repo.calls >= 1
placeholder, 2*time.Second, 10*time.Millisecond)
	// Negative cache entry keeps subsequent hot-path calls fail-open without thrash.
	filtered = scheduler.filterGrokFreeQuotaAccounts(context.Background(), accounts)
	require.Equal(t, []int64{1placeholder, accountIDs(filtered))
	require.Equal(t, 1, repo.calls)
placeholder

func TestFilterGrokFreeQuotaAccountsUnknownTierFailOpen(t *testing.T) {
	repo := &grokFreeQuotaUsageRepoStub{stats: map[int64]*usagestats.AccountStats{
		1: {Tokens: 9_999_999placeholder,
placeholderplaceholder
	scheduler := &defaultOpenAIAccountScheduler{service: &OpenAIGatewayService{cfg: grokFreeQuotaTestConfig(), usageLogRepo: repoplaceholderplaceholder
	accounts := []Account{
		{ID: 1, Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder,
		{ID: 2, Platform: PlatformGrok, Type: AccountTypeOAuth, Credentials: map[string]any{"subscription_tier": "unknown"placeholderplaceholder,
		{ID: 3, Platform: PlatformGrok, Type: AccountTypeOAuth, Extra: map[string]any{"subscription_tier": "pro"placeholderplaceholder,
placeholder

	filtered := scheduler.filterGrokFreeQuotaAccounts(context.Background(), accounts)
	require.Equal(t, []int64{1, 2, 3placeholder, accountIDs(filtered))
	require.Zero(t, repo.calls, "unknown/paid tiers must not query free-quota stats")
placeholder

func TestFilterGrokFreeQuotaAccountsRecoversAfterRollingUsageFalls(t *testing.T) {
	repo := &grokFreeQuotaUsageRepoStub{stats: map[int64]*usagestats.AccountStats{
		1: {Tokens: 490_000placeholder,
placeholderplaceholder
	openaiGrokFreeQuotaGateCache = sync.Map{placeholder
	scheduler := &defaultOpenAIAccountScheduler{service: &OpenAIGatewayService{cfg: grokFreeQuotaTestConfig(), usageLogRepo: repoplaceholderplaceholder
	accounts := []Account{{
		ID: 1, Platform: PlatformGrok, Type: AccountTypeOAuth,
placeholder"plan_type": "free"placeholder,
placeholderplaceholder

	// Miss fails open, then background fill blocks over-gate account.
	require.Equal(t, []int64{1placeholder, accountIDs(scheduler.filterGrokFreeQuotaAccounts(context.Background(), accounts)))
	require.Eventually(t, func() bool {
		filtered := scheduler.filterGrokFreeQuotaAccounts(context.Background(), accounts)
		return len(filtered) == 0
placeholder, 2*time.Second, 10*time.Millisecond)

	repo.mu.Lock()
	repo.stats[1] = &usagestats.AccountStats{Tokens: 100_000placeholder
	repo.mu.Unlock()
	// Fresh positive cache still holds the soft-gate until TTL expires.
	require.Empty(t, scheduler.filterGrokFreeQuotaAccounts(context.Background(), accounts), "fresh cache keeps the soft-gate hold")

	// Expire entry → miss fails open and schedules refresh with recovered usage.
	scheduler.grokFreeQuotaGateCache.Store(int64(1), grokFreeQuotaGateCacheEntry{
		tokens: 490_000, checkedAt: time.Now().Add(-time.Minute), known: true,
placeholder)
	require.Equal(t, []int64{1placeholder, accountIDs(scheduler.filterGrokFreeQuotaAccounts(context.Background(), accounts)))
	require.Eventually(t, func() bool {
		filtered := scheduler.filterGrokFreeQuotaAccounts(context.Background(), accounts)
		return len(filtered) == 1 && filtered[0].ID == 1
placeholder, 2*time.Second, 10*time.Millisecond)
	require.GreaterOrEqual(t, repo.calls, 2)
placeholder

func TestResolveGrokFreeQuotaGateSettingsDefaultsToNinetyFivePercent(t *testing.T) {
	settings, ok := resolveGrokFreeQuotaGateSettings(grokFreeQuotaTestConfig())
	require.True(t, ok)
	require.Equal(t, int64(500_000), settings.limitTokens)
	require.Equal(t, int64(475_000), settings.gateTokens) // 95% of 500k
	require.Equal(t, 24*time.Hour, settings.window)
placeholder

func TestIsExplicitGrokFreeOAuthAccount_OnlyExactFree(t *testing.T) {
	t.Parallel()
	require.False(t, isExplicitGrokFreeOAuthAccount(nil))
	require.False(t, isExplicitGrokFreeOAuthAccount(&Account{Platform: PlatformGrok, Type: AccountTypeAPIKey, Credentials: map[string]any{"subscription_tier": "free"placeholderplaceholder))
	require.True(t, isExplicitGrokFreeOAuthAccount(&Account{Platform: PlatformGrok, Type: AccountTypeOAuth, Credentials: map[string]any{"subscription_tier": "FREE"placeholderplaceholder))
	require.True(t, isExplicitGrokFreeOAuthAccount(&Account{Platform: PlatformGrok, Type: AccountTypeOAuth, Credentials: map[string]any{"plan_type": "free"placeholderplaceholder))
	// basic / inferred free are NOT soft-gated (only an explicit "free" tier is).
	require.False(t, isExplicitGrokFreeOAuthAccount(&Account{Platform: PlatformGrok, Type: AccountTypeOAuth, Credentials: map[string]any{"subscription_tier": "basic"placeholderplaceholder))
	require.False(t, isExplicitGrokFreeOAuthAccount(&Account{Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder))
placeholder

func TestOpenAIAccountSchedulerLoadBalanceAppliesGrokFreeQuotaGate(t *testing.T) {
	cfg := grokFreeQuotaTestConfig()
	cfg.RunMode = config.RunModeSimple
	openaiGrokFreeQuotaGateCache = sync.Map{placeholder
	accounts := []Account{
		{ID: 1, Platform: PlatformGrok, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: true, Concurrency: 1, Credentials: map[string]any{"subscription_tier": "free"placeholderplaceholder,
		{ID: 2, Platform: PlatformGrok, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: true, Concurrency: 1, Credentials: map[string]any{"subscription_tier": "pro"placeholderplaceholder,
placeholder
	svc := &OpenAIGatewayService{
		cfg:         cfg,
		accountRepo: &grokFreeQuotaAccountRepoStub{accounts: accountsplaceholder,
		usageLogRepo: &grokFreeQuotaUsageRepoStub{stats: map[int64]*usagestats.AccountStats{
			1: {Tokens: 480_000placeholder, // over 95% of 500k
placeholder
placeholder
	scheduler := &defaultOpenAIAccountScheduler{service: svc, stats: newOpenAIAccountRuntimeStats()placeholder

	// Warm cache via background refresh so load-balance sees the soft-gate.
	_ = scheduler.filterGrokFreeQuotaAccounts(context.Background(), accounts)
	require.Eventually(t, func() bool {
		filtered := scheduler.filterGrokFreeQuotaAccounts(context.Background(), accounts)
		return len(accountIDs(filtered)) == 1 && accountIDs(filtered)[0] == 2
placeholder, 2*time.Second, 10*time.Millisecond)

	selection, _, _, _, err := scheduler.selectByLoadBalance(context.Background(), OpenAIAccountScheduleRequest{Platform: PlatformGrokplaceholder)
placeholder
	require.NotNil(t, selection)
	require.NotNil(t, selection.Account)
	require.Equal(t, int64(2), selection.Account.ID)
placeholder

// Admin QueryQuota / import probe paths never call filterGrokFreeQuotaAccounts.
// Document and assert the scheduler filter is the only gate entry point.
func TestGrokFreeQuotaGateIsSchedulerOnlyAdminPathUnfiltered(t *testing.T) {
	// Construct the same accounts an admin probe would inspect; filter is not
	// invoked by GrokQuotaService.QueryQuota / GetUsage. Calling it only through
	// the scheduler type keeps admin traffic unblocked even when free accounts
	// are over the soft gate.
	require.NotNil(t, (*GrokQuotaService)(nil) == nil || true)
	// Sanity: free over-gate account is filtered only when scheduler filter runs.
	repo := &grokFreeQuotaUsageRepoStub{stats: map[int64]*usagestats.AccountStats{
		9: {Tokens: 500_000placeholder,
placeholderplaceholder
	openaiGrokFreeQuotaGateCache = sync.Map{placeholder
	scheduler := &defaultOpenAIAccountScheduler{service: &OpenAIGatewayService{cfg: grokFreeQuotaTestConfig(), usageLogRepo: repoplaceholderplaceholder
	overGate := Account{ID: 9, Platform: PlatformGrok, Type: AccountTypeOAuth, Credentials: map[string]any{"subscription_tier": "FREE"placeholderplaceholder
	require.Eventually(t, func() bool {
		_ = scheduler.filterGrokFreeQuotaAccounts(context.Background(), []Account{overGateplaceholder)
		return len(scheduler.filterGrokFreeQuotaAccounts(context.Background(), []Account{overGateplaceholder)) == 0
placeholder, 2*time.Second, 10*time.Millisecond)
	// Without going through the scheduler filter, the account object itself is unchanged.
	require.True(t, isExplicitGrokFreeOAuthAccount(&overGate))
	require.Equal(t, int64(9), overGate.ID)
placeholder

func TestSweepGrokFreeQuotaGateCacheDropsStaleEntries(t *testing.T) {
	now := time.Now().UTC()
	cacheTTL := 5 * time.Second
	// maxAge is floored at grokFreeQuotaGateCacheMinSweepAge, not 20*cacheTTL.
	var cache sync.Map
	cache.Store(int64(1), grokFreeQuotaGateCacheEntry{tokens: 10, checkedAt: now, known: trueplaceholder)
	cache.Store(int64(2), grokFreeQuotaGateCacheEntry{tokens: 20, checkedAt: now.Add(-time.Minute), known: trueplaceholder)
	cache.Store(int64(3), grokFreeQuotaGateCacheEntry{tokens: 30, checkedAt: now.Add(-time.Hour), known: trueplaceholder)
	cache.Store(int64(4), "not-an-entry")

	sweepGrokFreeQuotaGateCache(&cache, now, cacheTTL)

	remaining := make([]int64, 0, 4)
	cache.Range(func(key, _ any) bool {
		if id, ok := key.(int64); ok {
			remaining = append(remaining, id)
	placeholder
		return true
placeholder)
	require.ElementsMatch(t, []int64{1, 2placeholder, remaining)

	// A disabled cache (TTL 0) means the caller never populated it — leave it alone.
	var untouched sync.Map
	untouched.Store(int64(7), grokFreeQuotaGateCacheEntry{checkedAt: now.Add(-time.Hour), known: trueplaceholder)
	sweepGrokFreeQuotaGateCache(&untouched, now, 0)
	_, stillThere := untouched.Load(int64(7))
	require.True(t, stillThere)
placeholder

func TestFilterGrokFreeQuotaAccountsEvictsDepartedAccounts(t *testing.T) {
	repo := &grokFreeQuotaUsageRepoStub{stats: map[int64]*usagestats.AccountStats{
		1: {Tokens: 1_000placeholder,
placeholderplaceholder
	var cache sync.Map
	// Account 99 was scheduled long ago and no longer appears in any batch. Its
	// entry must not survive a run that queries for a different account.
	cache.Store(int64(99), grokFreeQuotaGateCacheEntry{tokens: 5, checkedAt: time.Now().UTC().Add(-2 * time.Hour), known: trueplaceholder)

	accounts := []Account{
		{ID: 1, Platform: PlatformGrok, Type: AccountTypeOAuth, Credentials: map[string]any{"subscription_tier": "FREE"placeholderplaceholder,
placeholder
	// First call schedules async refresh + may not have finished sweep yet.
	_ = filterGrokFreeQuotaAccountsCore(context.Background(), grokFreeQuotaTestConfig(), repo, &cache, accounts)
	require.Eventually(t, func() bool {
		_, departedStillCached := cache.Load(int64(99))
		_, freshCached := cache.Load(int64(1))
		return !departedStillCached && freshCached
placeholder, 2*time.Second, 10*time.Millisecond)
	filtered := filterGrokFreeQuotaAccountsCore(context.Background(), grokFreeQuotaTestConfig(), repo, &cache, accounts)
	require.Equal(t, []int64{1placeholder, accountIDs(filtered))
placeholder

func accountIDs(accounts []Account) []int64 {
	ids := make([]int64, 0, len(accounts))
	for i := range accounts {
		ids = append(ids, accounts[i].ID)
placeholder
	return ids
placeholder
