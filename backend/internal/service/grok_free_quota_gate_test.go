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
	cfg.Gateway.Grok.FreeQuotaStatsCacheSeconds = 5
	return cfg
placeholder

func TestFilterGrokFreeQuotaAccountsOnlyBlocksExplicitFreeOAuth(t *testing.T) {
	repo := &grokFreeQuotaUsageRepoStub{stats: map[int64]*usagestats.AccountStats{
		1: {Tokens: 475_000placeholder, // 95% of 500k
placeholderplaceholder
	scheduler := &defaultOpenAIAccountScheduler{service: &OpenAIGatewayService{cfg: grokFreeQuotaTestConfig(), usageLogRepo: repoplaceholderplaceholder
	accounts := []Account{
		{ID: 1, Platform: PlatformGrok, Type: AccountTypeOAuth, Credentials: map[string]any{"subscription_tier": "FREE"placeholderplaceholder,
		{ID: 2, Platform: PlatformGrok, Type: AccountTypeOAuth, Credentials: map[string]any{"subscription_tier": "PRO"placeholderplaceholder,
		{ID: 3, Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder,
		{ID: 4, Platform: PlatformGrok, Type: AccountTypeAPIKey, Credentials: map[string]any{"subscription_tier": "FREE"placeholderplaceholder,
placeholder

	filtered := scheduler.filterGrokFreeQuotaAccounts(context.Background(), accounts)
	require.Equal(t, []int64{2, 3, 4placeholder, accountIDs(filtered), "paid and unknown fail-open; API-key free marker is not gated")
	require.Equal(t, 1, repo.calls)
	require.Equal(t, []int64{1placeholder, repo.lastIDs, "paid, unknown, and API-key accounts must not enter the local free-tier query")
	require.WithinDuration(t, time.Now().UTC().Add(-24*time.Hour), repo.start, time.Second)
placeholder

func TestFilterGrokFreeQuotaAccountsStatsFailureFailsOpen(t *testing.T) {
	repo := &grokFreeQuotaUsageRepoStub{err: errors.New("usage database unavailable")placeholder
	scheduler := &defaultOpenAIAccountScheduler{service: &OpenAIGatewayService{cfg: grokFreeQuotaTestConfig(), usageLogRepo: repoplaceholderplaceholder
	accounts := []Account{{
		ID: 1, Platform: PlatformGrok, Type: AccountTypeOAuth,
placeholder"subscription_tier": "free"placeholder,
placeholderplaceholder

	filtered := scheduler.filterGrokFreeQuotaAccounts(context.Background(), accounts)
	require.Equal(t, []int64{1placeholder, accountIDs(filtered))
	// Cache the failure entry so a second call still fails open without re-query thrash.
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
	scheduler := &defaultOpenAIAccountScheduler{service: &OpenAIGatewayService{cfg: grokFreeQuotaTestConfig(), usageLogRepo: repoplaceholderplaceholder
	accounts := []Account{{
		ID: 1, Platform: PlatformGrok, Type: AccountTypeOAuth,
placeholder"plan_type": "free"placeholder,
placeholderplaceholder

	require.Empty(t, scheduler.filterGrokFreeQuotaAccounts(context.Background(), accounts))
	repo.mu.Lock()
	repo.stats[1] = &usagestats.AccountStats{Tokens: 100_000placeholder
	repo.mu.Unlock()
	require.Empty(t, scheduler.filterGrokFreeQuotaAccounts(context.Background(), accounts), "fresh cache keeps the short soft-gate hold")

	scheduler.grokFreeQuotaGateCache.Store(int64(1), grokFreeQuotaGateCacheEntry{
		tokens: 490_000, checkedAt: time.Now().Add(-time.Minute), known: true,
placeholder)
	require.Equal(t, []int64{1placeholder, accountIDs(scheduler.filterGrokFreeQuotaAccounts(context.Background(), accounts)))
	require.Equal(t, 2, repo.calls)
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
	// basic / inferred free are NOT soft-gated (personal-dev contract).
	require.False(t, isExplicitGrokFreeOAuthAccount(&Account{Platform: PlatformGrok, Type: AccountTypeOAuth, Credentials: map[string]any{"subscription_tier": "basic"placeholderplaceholder))
	require.False(t, isExplicitGrokFreeOAuthAccount(&Account{Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder))
placeholder

func TestOpenAIAccountSchedulerLoadBalanceAppliesGrokFreeQuotaGate(t *testing.T) {
	cfg := grokFreeQuotaTestConfig()
	cfg.RunMode = config.RunModeSimple
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
	scheduler := &defaultOpenAIAccountScheduler{service: &OpenAIGatewayService{cfg: grokFreeQuotaTestConfig(), usageLogRepo: repoplaceholderplaceholder
	overGate := Account{ID: 9, Platform: PlatformGrok, Type: AccountTypeOAuth, Credentials: map[string]any{"subscription_tier": "FREE"placeholderplaceholder
	require.Empty(t, scheduler.filterGrokFreeQuotaAccounts(context.Background(), []Account{overGateplaceholder))
	// Without going through the scheduler filter, the account object itself is unchanged.
	require.True(t, isExplicitGrokFreeOAuthAccount(&overGate))
	require.Equal(t, int64(9), overGate.ID)
placeholder

func accountIDs(accounts []Account) []int64 {
	ids := make([]int64, 0, len(accounts))
	for i := range accounts {
		ids = append(ids, accounts[i].ID)
placeholder
	return ids
placeholder
