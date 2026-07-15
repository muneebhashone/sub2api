package service

import (
	"context"
	"errors"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type upstreamCostTrackingConcurrencyCache struct {
	ConcurrencyCache
	loadMap       map[int64]*AccountLoadInfo
	acquireLimits map[int64][]int
	releases      map[int64]int
	rejectAcquire bool
placeholder

func (c *upstreamCostTrackingConcurrencyCache) AcquireAccountSlot(_ context.Context, accountID int64, maxConcurrency int, _ string) (bool, error) {
	if c.acquireLimits == nil {
		c.acquireLimits = make(map[int64][]int)
placeholder
	c.acquireLimits[accountID] = append(c.acquireLimits[accountID], maxConcurrency)
	return !c.rejectAcquire, nil
placeholder

func (c *upstreamCostTrackingConcurrencyCache) ReleaseAccountSlot(_ context.Context, accountID int64, _ string) error {
	if c.releases == nil {
		c.releases = make(map[int64]int)
placeholder
	c.releases[accountID]++
	return nil
placeholder

func (c *upstreamCostTrackingConcurrencyCache) GetAccountsLoadBatch(_ context.Context, accounts []AccountWithConcurrency) (map[int64]*AccountLoadInfo, error) {
	out := make(map[int64]*AccountLoadInfo, len(accounts))
	for _, account := range accounts {
		if load := c.loadMap[account.ID]; load != nil {
			copied := *load
			out[account.ID] = &copied
	placeholder
placeholder
	return out, nil
placeholder

func (c *upstreamCostTrackingConcurrencyCache) limits(accountID int64) []int {
	return append([]int(nil), c.acquireLimits[accountID]...)
placeholder

func (c *upstreamCostTrackingConcurrencyCache) releaseCount(accountID int64) int {
	return c.releases[accountID]
placeholder

func (c *upstreamCostTrackingConcurrencyCache) totalAcquires() int {
	total := 0
	for _, limits := range c.acquireLimits {
		total += len(limits)
placeholder
	return total
placeholder

type upstreamCostCountingAccountRepo struct {
	AccountRepository
	accounts map[int64]*Account
	getCalls int
placeholder

func (r *upstreamCostCountingAccountRepo) GetByID(_ context.Context, accountID int64) (*Account, error) {
	r.getCalls++
	account := r.accounts[accountID]
	if account == nil {
		return nil, errors.New("account not found")
placeholder
	cloned := *account
	return &cloned, nil
placeholder

func (r *upstreamCostCountingAccountRepo) calls() int {
	return r.getCalls
placeholder

func upstreamCostTestAccount(id int64, status string, rate float64, receivedAt time.Time, interval time.Duration) *Account {
placeholder
		ID:       id,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Extra: map[string]any{
			UpstreamBillingProbeExtraKey: map[string]any{
				"status": status,
				"data": map[string]any{
					"billing_scope":             "token",
					"resolved_rate_multiplier":  rate,
					"peak_rate_enabled":         false,
					"effective_rate_multiplier": rate,
			placeholder,
				"received_at":     receivedAt.UTC().Format(time.RFC3339Nano),
				"fresh_until":     receivedAt.Add(2 * interval).UTC().Format(time.RFC3339Nano),
				"last_attempt_at": receivedAt.UTC().Format(time.RFC3339Nano),
				"next_probe_at":   receivedAt.Add(interval).UTC().Format(time.RFC3339Nano),
		placeholder,
	placeholder,
placeholder
placeholder

func upstreamCostTestOAuthAccount(id int64) *Account {
placeholderID: id, Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder
placeholder

func TestAdvancedCostSchedulerUsesTopKOverflowWhenPreferredAccountIsKnownFull(t *testing.T) {
	resetOpenAIAdvancedSchedulerSettingCacheForTest()
	defer resetOpenAIAdvancedSchedulerSettingCacheForTest()

	now := time.Now()
	cheap := upstreamCostTestAccount(1, UpstreamBillingProbeStatusOK, 0.03, now.Add(-time.Minute), 30*time.Minute)
	expensive := upstreamCostTestAccount(2, UpstreamBillingProbeStatusOK, 0.8, now.Add(-time.Minute), 30*time.Minute)
	for _, account := range []*Account{cheap, expensiveplaceholder {
		account.Status = StatusActive
		account.Schedulable = true
		account.Concurrency = 1
placeholder
	cache := &upstreamCostTrackingConcurrencyCache{loadMap: map[int64]*AccountLoadInfo{
		cheap.ID:     {AccountID: cheap.ID, CurrentConcurrency: 1, LoadRate: 100placeholder,
		expensive.ID: {AccountID: expensive.IDplaceholder,
placeholderplaceholder
	cfg := &config.Config{placeholder
	cfg.Gateway.OpenAIWS.LBTopK = 1
	cfg.Gateway.OpenAIWS.SchedulerScoreWeights.UpstreamCost = 1
	svc := &OpenAIGatewayService{
		accountRepo:        schedulerTestOpenAIAccountRepo{accounts: []Account{*cheap, *expensiveplaceholderplaceholder,
		cfg:                cfg,
		rateLimitService:   newOpenAIAdvancedSchedulerRateLimitService("true"),
		concurrencyService: NewConcurrencyService(cache),
placeholder
	groupID := int64(1)

	selection, _, err := svc.SelectAccountWithScheduler(context.Background(), &groupID, "", "", "gpt-test", nil, OpenAIUpstreamTransportAny, false)
placeholder
	require.Equal(t, expensive.ID, selection.Account.ID)
	require.Empty(t, cache.limits(cheap.ID))
	require.Equal(t, []int{1placeholder, cache.limits(expensive.ID))
	selection.ReleaseFunc()
placeholder

func TestAdvancedSchedulerCapsRejectedCostOverflowAcquires(t *testing.T) {
	selectionOrder := make([]openAIAccountCandidateScore, 0, 15_000)
	for id := int64(1); id <= 15_000; id++ {
		account := &Account{ID: id, Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Status: StatusActive, Schedulable: true, Concurrency: 1placeholder
		selectionOrder = append(selectionOrder, openAIAccountCandidateScore{
			account: account, loadInfo: &AccountLoadInfo{AccountID: idplaceholder, loadKnown: false,
	placeholder)
placeholder
	cache := &upstreamCostTrackingConcurrencyCache{rejectAcquire: trueplaceholder
	scheduler := &defaultOpenAIAccountScheduler{service: &OpenAIGatewayService{
		concurrencyService: NewConcurrencyService(cache),
placeholderplaceholder

	selection, _, err := scheduler.tryAcquireOpenAISelectionOrder(
		context.Background(), OpenAIAccountScheduleRequest{Platform: PlatformOpenAIplaceholder, selectionOrder,
	)

placeholder
	require.Nil(t, selection)
	require.Equal(t, openAIAccountSelectionProbeLimit, cache.totalAcquires())
placeholder

func TestOpenAICostOverflowExpandedOnlyWhenCostAddsCandidates(t *testing.T) {
	candidates := []openAIAccountCandidateScore{
		{account: &Account{ID: 1, Extra: map[string]any{"openai_compact_supported": trueplaceholderplaceholderplaceholder,
		{account: &Account{ID: 2placeholderplaceholder,
placeholder
	plan := openAIAccountLoadPlan{candidates: candidates, topK: 1, includeOverflowFallback: trueplaceholder
	require.True(t, openAICostOverflowExpanded(OpenAIAccountScheduleRequest{placeholder, plan))
	require.False(t, openAICostOverflowExpanded(OpenAIAccountScheduleRequest{RequireCompact: trueplaceholder, plan),
		"one candidate per compact tier does not expand either tier's top-k")
	plan.topK = len(candidates)
	require.False(t, openAICostOverflowExpanded(OpenAIAccountScheduleRequest{placeholder, plan))
	plan.includeOverflowFallback = false
	plan.topK = 1
	require.False(t, openAICostOverflowExpanded(OpenAIAccountScheduleRequest{placeholder, plan))
placeholder

func TestAdvancedSchedulerKnownFullOverflowStillFindsAvailableAccount(t *testing.T) {
	selectionOrder := make([]openAIAccountCandidateScore, 0, openAIAccountSelectionProbeLimit+2)
	for id := int64(1); id <= openAIAccountSelectionProbeLimit+1; id++ {
		account := &Account{ID: id, Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Status: StatusActive, Schedulable: true, Concurrency: 1placeholder
		selectionOrder = append(selectionOrder, openAIAccountCandidateScore{
			account:   account,
			loadInfo:  &AccountLoadInfo{AccountID: id, CurrentConcurrency: 1, LoadRate: 100placeholder,
			loadKnown: true,
	placeholder)
placeholder
	availableID := int64(openAIAccountSelectionProbeLimit + 2)
	available := &Account{ID: availableID, Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Status: StatusActive, Schedulable: true, Concurrency: 1placeholder
	selectionOrder = append(selectionOrder, openAIAccountCandidateScore{
		account: available, loadInfo: &AccountLoadInfo{AccountID: availableIDplaceholder, loadKnown: true,
placeholder)
	cache := &upstreamCostTrackingConcurrencyCache{placeholder
	scheduler := &defaultOpenAIAccountScheduler{service: &OpenAIGatewayService{
		concurrencyService: NewConcurrencyService(cache),
placeholderplaceholder

	selection, _, err := scheduler.tryAcquireOpenAISelectionOrder(
		context.Background(), OpenAIAccountScheduleRequest{Platform: PlatformOpenAIplaceholder, selectionOrder,
	)

placeholder
	require.NotNil(t, selection)
	require.Equal(t, availableID, selection.Account.ID)
	require.Equal(t, 1, cache.totalAcquires())
	selection.ReleaseFunc()
placeholder

func TestAdvancedSchedulerSharesProbeBudgetWithFallbackDBRechecks(t *testing.T) {
	const size = 15_000
	latestAccounts := make(map[int64]*Account, size)
	snapshotAccounts := make(map[int64]*Account, size)
	selectionOrder := make([]openAIAccountCandidateScore, 0, size)
	for id := int64(1); id <= size; id++ {
		stale := &Account{ID: id, Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Status: StatusActive, Schedulable: true, Concurrency: 1placeholder
		latest := *stale
		latest.Status = StatusDisabled
		snapshotAccounts[id] = stale
		latestAccounts[id] = &latest
		selectionOrder = append(selectionOrder, openAIAccountCandidateScore{
			account: stale, loadInfo: &AccountLoadInfo{AccountID: idplaceholder, loadKnown: false,
	placeholder)
placeholder
	repo := &upstreamCostCountingAccountRepo{accounts: latestAccountsplaceholder
	cache := &upstreamCostTrackingConcurrencyCache{placeholder
	scheduler := &defaultOpenAIAccountScheduler{service: &OpenAIGatewayService{
		accountRepo:        repo,
		schedulerSnapshot:  &SchedulerSnapshotService{cache: &openAISnapshotCacheStub{accountsByID: snapshotAccountsplaceholderplaceholder,
		concurrencyService: NewConcurrencyService(cache),
placeholderplaceholder
	budget := newOpenAISelectionProbeBudget()
	budget.enableLimit()
	req := OpenAIAccountScheduleRequest{Platform: PlatformOpenAIplaceholder

	selection, _, err := scheduler.tryAcquireOpenAISelectionOrderWithBudget(context.Background(), req, selectionOrder, budget)
placeholder
	require.Nil(t, selection)
	selection, _, _, _, err = scheduler.finishLoadBalanceSelectionFallback(
		context.Background(), req, openAIAccountLoadSelectionAttempt{selectionOrder: selectionOrderplaceholder, budget,
	)

placeholder
	require.Nil(t, selection)
	require.Equal(t, openAIAccountSelectionProbeLimit, cache.totalAcquires())
	require.Equal(t, openAIAccountSelectionProbeLimit, repo.calls())
placeholder

func TestAdvancedCostSchedulerKeepsCompactSupportedOverflowAheadOfUnknown(t *testing.T) {
	resetOpenAIAdvancedSchedulerSettingCacheForTest()
	defer resetOpenAIAdvancedSchedulerSettingCacheForTest()

	now := time.Now()
	preferred := upstreamCostTestAccount(11, UpstreamBillingProbeStatusOK, 0.01, now.Add(-time.Minute), 30*time.Minute)
	overflow := upstreamCostTestAccount(12, UpstreamBillingProbeStatusOK, 0.1, now.Add(-time.Minute), 30*time.Minute)
	unknown := upstreamCostTestAccount(13, UpstreamBillingProbeStatusOK, 0.001, now.Add(-time.Minute), 30*time.Minute)
	preferred.Extra["openai_compact_supported"] = true
	overflow.Extra["openai_compact_supported"] = true
	for _, account := range []*Account{preferred, overflow, unknownplaceholder {
		account.Status = StatusActive
		account.Schedulable = true
		account.Concurrency = 1
placeholder
	cache := &upstreamCostTrackingConcurrencyCache{loadMap: map[int64]*AccountLoadInfo{
		preferred.ID: {AccountID: preferred.ID, CurrentConcurrency: 1, LoadRate: 100placeholder,
		overflow.ID:  {AccountID: overflow.IDplaceholder,
		unknown.ID:   {AccountID: unknown.IDplaceholder,
placeholderplaceholder
	cfg := &config.Config{placeholder
	cfg.Gateway.OpenAIWS.LBTopK = 1
	cfg.Gateway.OpenAIWS.SchedulerScoreWeights.UpstreamCost = 1
	svc := &OpenAIGatewayService{
		accountRepo:        schedulerTestOpenAIAccountRepo{accounts: []Account{*preferred, *overflow, *unknownplaceholderplaceholder,
		cfg:                cfg,
		rateLimitService:   newOpenAIAdvancedSchedulerRateLimitService("true"),
		concurrencyService: NewConcurrencyService(cache),
placeholder
	groupID := int64(1)

	selection, _, err := svc.SelectAccountWithScheduler(context.Background(), &groupID, "", "", "gpt-test", nil, OpenAIUpstreamTransportAny, true)
placeholder
	require.Equal(t, overflow.ID, selection.Account.ID)
	require.Empty(t, cache.limits(preferred.ID))
	require.Equal(t, []int{1placeholder, cache.limits(overflow.ID))
	require.Empty(t, cache.limits(unknown.ID))
	selection.ReleaseFunc()
placeholder

func TestAdvancedSchedulerUnknownLoadFailsOpen(t *testing.T) {
	account := &Account{ID: 21, Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Status: StatusActive, Schedulable: true, Concurrency: 1placeholder
	cache := &upstreamCostTrackingConcurrencyCache{placeholder
	scheduler := &defaultOpenAIAccountScheduler{service: &OpenAIGatewayService{concurrencyService: NewConcurrencyService(cache)placeholderplaceholder

	selection, _, err := scheduler.tryAcquireOpenAISelectionOrder(context.Background(), OpenAIAccountScheduleRequest{Platform: PlatformOpenAIplaceholder, []openAIAccountCandidateScore{{
		account: account, loadInfo: &AccountLoadInfo{AccountID: account.ID, CurrentConcurrency: 99placeholder, loadKnown: false,
placeholderplaceholder)
placeholder
	require.NotNil(t, selection)
	require.Equal(t, []int{1placeholder, cache.limits(account.ID))
	selection.ReleaseFunc()
placeholder

func TestAdvancedSchedulerReleasesSlotWhenDBDisablesCandidate(t *testing.T) {
	stale := &Account{ID: 31, Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Status: StatusActive, Schedulable: true, Concurrency: 1placeholder
	backup := &Account{ID: 32, Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Status: StatusActive, Schedulable: true, Concurrency: 1placeholder
	disabled := *stale
	disabled.Status = StatusDisabled
	repo := &upstreamCostCountingAccountRepo{accounts: map[int64]*Account{stale.ID: &disabled, backup.ID: backupplaceholderplaceholder
	snapshot := &openAISnapshotCacheStub{accountsByID: map[int64]*Account{stale.ID: stale, backup.ID: backupplaceholderplaceholder
	cache := &upstreamCostTrackingConcurrencyCache{placeholder
	scheduler := &defaultOpenAIAccountScheduler{service: &OpenAIGatewayService{
		accountRepo:        repo,
		schedulerSnapshot:  &SchedulerSnapshotService{cache: snapshotplaceholder,
		concurrencyService: NewConcurrencyService(cache),
placeholderplaceholder

	selection, _, err := scheduler.tryAcquireOpenAISelectionOrder(context.Background(), OpenAIAccountScheduleRequest{Platform: PlatformOpenAIplaceholder, []openAIAccountCandidateScore{
		{account: stale, loadInfo: &AccountLoadInfo{AccountID: stale.IDplaceholder, loadKnown: trueplaceholder,
		{account: backup, loadInfo: &AccountLoadInfo{AccountID: backup.IDplaceholder, loadKnown: trueplaceholder,
placeholder)
placeholder
	require.Equal(t, backup.ID, selection.Account.ID)
	require.Equal(t, 1, cache.releaseCount(stale.ID))
	selection.ReleaseFunc()
placeholder

func TestAdvancedSchedulerReacquiresOnceWhenDBConcurrencyChanges(t *testing.T) {
	stale := &Account{ID: 41, Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Status: StatusActive, Schedulable: true, Concurrency: 10placeholder
	latest := *stale
	latest.Concurrency = 1
	repo := &upstreamCostCountingAccountRepo{accounts: map[int64]*Account{stale.ID: &latestplaceholderplaceholder
	snapshot := &openAISnapshotCacheStub{accountsByID: map[int64]*Account{stale.ID: staleplaceholderplaceholder
	cache := &upstreamCostTrackingConcurrencyCache{placeholder
	scheduler := &defaultOpenAIAccountScheduler{service: &OpenAIGatewayService{
		accountRepo:        repo,
		schedulerSnapshot:  &SchedulerSnapshotService{cache: snapshotplaceholder,
		concurrencyService: NewConcurrencyService(cache),
placeholderplaceholder

	selection, _, err := scheduler.tryAcquireOpenAISelectionOrder(context.Background(), OpenAIAccountScheduleRequest{Platform: PlatformOpenAIplaceholder, []openAIAccountCandidateScore{{
		account: stale, loadInfo: &AccountLoadInfo{AccountID: stale.IDplaceholder, loadKnown: true,
placeholderplaceholder)
placeholder
	require.Equal(t, 1, selection.Account.Concurrency)
	require.Equal(t, []int{10, 1placeholder, cache.limits(stale.ID))
	require.Equal(t, 1, cache.releaseCount(stale.ID))
	selection.ReleaseFunc()
placeholder

func TestAdvancedSchedulerKnownFullPoolsDoNotRecheckDB(t *testing.T) {
	for _, size := range []int{100, 15_000placeholder {
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			accounts := make(map[int64]*Account, size)
			selectionOrder := make([]openAIAccountCandidateScore, 0, size)
			for i := 1; i <= size; i++ {
				account := &Account{ID: int64(i), Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Status: StatusActive, Schedulable: true, Concurrency: 1placeholder
				accounts[account.ID] = account
				selectionOrder = append(selectionOrder, openAIAccountCandidateScore{
					account:   account,
					loadInfo:  &AccountLoadInfo{AccountID: account.ID, CurrentConcurrency: 1, LoadRate: 100placeholder,
					loadKnown: true,
			placeholder)
		placeholder
			repo := &upstreamCostCountingAccountRepo{accounts: accountsplaceholder
			cache := &upstreamCostTrackingConcurrencyCache{placeholder
			scheduler := &defaultOpenAIAccountScheduler{service: &OpenAIGatewayService{
				accountRepo:        repo,
				schedulerSnapshot:  &SchedulerSnapshotService{cache: &openAISnapshotCacheStub{accountsByID: accountsplaceholderplaceholder,
				concurrencyService: NewConcurrencyService(cache),
		placeholderplaceholder

			selection, _, err := scheduler.tryAcquireOpenAISelectionOrder(context.Background(), OpenAIAccountScheduleRequest{Platform: PlatformOpenAIplaceholder, selectionOrder)
		placeholder
			require.Nil(t, selection)
			require.Zero(t, repo.calls())
			require.Zero(t, cache.totalAcquires())
	placeholder)
placeholder
placeholder

func TestOpenAIFreshUpstreamBillingRateRecomputesPeakAtSelectionTime(t *testing.T) {
	receivedAt := time.Date(2026, 7, 13, 17, 30, 0, 0, time.UTC)
	account := upstreamCostTestAccount(1, UpstreamBillingProbeStatusOK, 0.4, receivedAt, time.Hour)
	snapshot, ok := account.Extra[UpstreamBillingProbeExtraKey].(map[string]any)
	require.True(t, ok)
	snapshot["data"] = map[string]any{
		"billing_scope":             "token",
		"resolved_rate_multiplier":  0.4,
		"peak_rate_enabled":         true,
		"peak_start":                "09:00",
		"peak_end":                  "18:00",
		"peak_rate_multiplier":      2.0,
		"applied_peak_multiplier":   2.0,
		"effective_rate_multiplier": 0.8,
		"timezone":                  "UTC",
placeholder

	duringPeak, ok := openAIFreshUpstreamBillingRate(account, time.Date(2026, 7, 13, 17, 59, 0, 0, time.UTC))
	require.True(t, ok)
	require.Equal(t, 0.8, duringPeak)

	afterPeak, ok := openAIFreshUpstreamBillingRate(account, time.Date(2026, 7, 13, 18, 1, 0, 0, time.UTC))
	require.True(t, ok)
	require.Equal(t, 0.4, afterPeak)
placeholder

func TestOpenAIUpstreamCostFactorsSparseProbeIsNeutral(t *testing.T) {
	now := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	accounts := make([]*Account, 0, 10)
	accounts = append(accounts, upstreamCostTestAccount(1, UpstreamBillingProbeStatusOK, 1, now.Add(-time.Minute), 30*time.Minute))
	for id := int64(2); id <= 10; id++ {
		accounts = append(accounts, &Account{
			ID:       id,
			Platform: PlatformOpenAI,
			Type:     AccountTypeAPIKey,
			Extra: map[string]any{
				UpstreamBillingProbeExtraKey: map[string]any{
					"status":          UpstreamBillingProbeStatusFailed,
					"last_attempt_at": now.UTC().Format(time.RFC3339Nano),
					"next_probe_at":   now.Add(time.Hour).UTC().Format(time.RFC3339Nano),
			placeholder,
		placeholder,
	placeholder)
placeholder

	factors := openAIUpstreamCostFactors(accounts, now, defaultOpenAIOAuthSchedulingRateMultiplier)
	for id := int64(1); id <= 10; id++ {
		require.Equal(t, openAIUpstreamCostNeutralFactor, factors[id])
placeholder
placeholder

func TestOpenAIUpstreamCostFactorsCoverageShrinksSparseSignal(t *testing.T) {
	now := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	accounts := []*Account{
		upstreamCostTestAccount(1, UpstreamBillingProbeStatusOK, 0.03, now.Add(-time.Minute), 30*time.Minute),
		upstreamCostTestAccount(2, UpstreamBillingProbeStatusOK, 0.8, now.Add(-time.Minute), 30*time.Minute),
placeholder
	for id := int64(3); id <= 10; id++ {
		accounts = append(accounts, &Account{ID: id, Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder)
placeholder

	factors := openAIUpstreamCostFactors(accounts, now, defaultOpenAIOAuthSchedulingRateMultiplier)
	center := math.Sqrt(0.03 * 0.8)
	require.InDelta(t, 0.5+0.2*(1/(1+0.03/center)-0.5), factors[1], 1e-12)
	require.InDelta(t, 0.5+0.2*(1/(1+0.8/center)-0.5), factors[2], 1e-12)
	require.Equal(t, openAIUpstreamCostNeutralFactor, factors[3])
placeholder

func TestOpenAIUpstreamCostFactorsUseMedianAgainstOutlier(t *testing.T) {
	now := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	accounts := []*Account{
		upstreamCostTestAccount(1, UpstreamBillingProbeStatusOK, 0.1, now.Add(-time.Minute), 30*time.Minute),
		upstreamCostTestAccount(2, UpstreamBillingProbeStatusOK, 0.2, now.Add(-time.Minute), 30*time.Minute),
		upstreamCostTestAccount(3, UpstreamBillingProbeStatusOK, 100, now.Add(-time.Minute), 30*time.Minute),
placeholder

	factors := openAIUpstreamCostFactors(accounts, now, defaultOpenAIOAuthSchedulingRateMultiplier)
	require.InDelta(t, 2.0/3.0, factors[1], 1e-12)
	require.InDelta(t, 0.5, factors[2], 1e-12)
	require.InDelta(t, 1/(1+100/0.2), factors[3], 1e-12)
placeholder

func TestOpenAILegacyUpstreamRateOrderRequiresComparableRates(t *testing.T) {
	now := time.Now()
	oneKnown := newOpenAILegacyUpstreamRateOrder([]*Account{
		upstreamCostTestAccount(1, UpstreamBillingProbeStatusOK, 0.03, now.Add(-time.Minute), 30*time.Minute),
		{ID: 2, Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder,
placeholder, now, defaultOpenAIOAuthSchedulingRateMultiplier)
	require.False(t, oneKnown.enabled)

	allEqual := newOpenAILegacyUpstreamRateOrder([]*Account{
		upstreamCostTestAccount(1, UpstreamBillingProbeStatusOK, 0.3, now.Add(-time.Minute), 30*time.Minute),
		upstreamCostTestAccount(2, UpstreamBillingProbeStatusOK, 0.3, now.Add(-time.Minute), 30*time.Minute),
placeholder, now, defaultOpenAIOAuthSchedulingRateMultiplier)
	require.False(t, allEqual.enabled)

	distinct := newOpenAILegacyUpstreamRateOrder([]*Account{
		upstreamCostTestAccount(1, UpstreamBillingProbeStatusOK, 0.03, now.Add(-time.Minute), 30*time.Minute),
		upstreamCostTestAccount(2, UpstreamBillingProbeStatusOK, 0.8, now.Add(-time.Minute), 30*time.Minute),
		{ID: 3, Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder,
placeholder, now, defaultOpenAIOAuthSchedulingRateMultiplier)
	require.True(t, distinct.enabled)
	require.Negative(t, distinct.compare(&Account{ID: 1placeholder, &Account{ID: 2placeholder))
	require.Negative(t, distinct.compare(&Account{ID: 2placeholder, &Account{ID: 3placeholder))
placeholder

func TestOpenAISchedulingRatePlacesOAuthAtConfiguredReference(t *testing.T) {
	now := time.Now()
	cheap := upstreamCostTestAccount(1, UpstreamBillingProbeStatusOK, 0.02, now.Add(-time.Minute), 30*time.Minute)
	oauth := upstreamCostTestOAuthAccount(2)
	expensive := upstreamCostTestAccount(3, UpstreamBillingProbeStatusOK, 0.12, now.Add(-time.Minute), 30*time.Minute)

	order := newOpenAILegacyUpstreamRateOrder([]*Account{cheap, oauth, expensiveplaceholder, now, 0.05)
	require.True(t, order.enabled)
	require.Negative(t, order.compare(cheap, oauth))
	require.Negative(t, order.compare(oauth, expensive))

	factors := openAIUpstreamCostFactors([]*Account{cheap, oauth, expensiveplaceholder, now, 0.05)
	require.Greater(t, factors[cheap.ID], factors[oauth.ID])
	require.Greater(t, factors[oauth.ID], factors[expensive.ID])
placeholder

func TestOpenAIGatewayServiceLegacyLowRatePriorityUsesConfiguredOAuthReference(t *testing.T) {
	resetOpenAIAdvancedSchedulerSettingCacheForTest()
	defer resetOpenAIAdvancedSchedulerSettingCacheForTest()

	now := time.Now()
	cheap := upstreamCostTestAccount(1, UpstreamBillingProbeStatusOK, 0.02, now.Add(-time.Minute), 30*time.Minute)
	oauth := upstreamCostTestOAuthAccount(2)
	expensive := upstreamCostTestAccount(3, UpstreamBillingProbeStatusOK, 0.12, now.Add(-time.Minute), 30*time.Minute)
	for _, account := range []*Account{cheap, oauth, expensiveplaceholder {
		account.Status = StatusActive
		account.Schedulable = true
		account.Concurrency = 1
placeholder
	cheap.Priority, oauth.Priority, expensive.Priority = 20, 10, 0

	settings := &openAIAdvancedSchedulerSettingRepoStub{values: map[string]string{
		openAIAdvancedSchedulerSettingKey:              "false",
		SettingKeyOpenAILowUpstreamRatePriorityEnabled: "true",
		SettingKeyOpenAIOAuthSchedulingRateMultiplier:  "0.05",
placeholderplaceholder
	cfg := &config.Config{placeholder
	svc := &OpenAIGatewayService{
		accountRepo:      schedulerTestOpenAIAccountRepo{accounts: []Account{*cheap, *oauth, *expensiveplaceholderplaceholder,
		cache:            &schedulerTestGatewayCache{placeholder,
		cfg:              cfg,
		rateLimitService: &RateLimitService{settingService: NewSettingService(settings, cfg)placeholder,
placeholder
	groupID := int64(1)

	first, _, err := svc.SelectAccountWithScheduler(context.Background(), &groupID, "", "", "gpt-test", nil, OpenAIUpstreamTransportAny, false)
placeholder
	require.Equal(t, cheap.ID, first.Account.ID)

	second, _, err := svc.SelectAccountWithScheduler(context.Background(), &groupID, "", "", "gpt-test", map[int64]struct{placeholder{cheap.ID: {placeholderplaceholder, OpenAIUpstreamTransportAny, false)
placeholder
	require.Equal(t, oauth.ID, second.Account.ID)
placeholder

func TestOpenAIModelsSelectionIgnoresTokenCostSignal(t *testing.T) {
	resetOpenAIAdvancedSchedulerSettingCacheForTest()
	defer resetOpenAIAdvancedSchedulerSettingCacheForTest()

	now := time.Now()
	cheap := upstreamCostTestAccount(51, UpstreamBillingProbeStatusOK, 0.02, now.Add(-time.Minute), 30*time.Minute)
	expensive := upstreamCostTestAccount(52, UpstreamBillingProbeStatusOK, 0.8, now.Add(-time.Minute), 30*time.Minute)
	for _, account := range []*Account{cheap, expensiveplaceholder {
		account.Status = StatusActive
		account.Schedulable = true
		account.Concurrency = 1
placeholder
	cheap.Priority = 10
	expensive.Priority = 0
	settings := &openAIAdvancedSchedulerSettingRepoStub{values: map[string]string{
		SettingKeyOpenAILowUpstreamRatePriorityEnabled: "true",
placeholderplaceholder
	cfg := &config.Config{placeholder
	svc := &OpenAIGatewayService{
		accountRepo:      schedulerTestOpenAIAccountRepo{accounts: []Account{*cheap, *expensiveplaceholderplaceholder,
		cfg:              cfg,
		rateLimitService: &RateLimitService{settingService: NewSettingService(settings, cfg)placeholder,
placeholder

	account, err := svc.SelectAccountForModelWithExclusions(context.Background(), nil, "", "", nil)
placeholder
	require.Equal(t, expensive.ID, account.ID)
placeholder

func TestOpenAIGatewayServiceLegacyLowRatePriorityIsIndependentFromAdvancedScheduler(t *testing.T) {
	now := time.Now()
	cheap := upstreamCostTestAccount(1, UpstreamBillingProbeStatusOK, 0.03, now.Add(-time.Minute), 30*time.Minute)
	cheap.Status, cheap.Schedulable, cheap.Concurrency, cheap.Priority = StatusActive, true, 1, 10
	expensive := upstreamCostTestAccount(2, UpstreamBillingProbeStatusOK, 0.8, now.Add(-time.Minute), 30*time.Minute)
	expensive.Status, expensive.Schedulable, expensive.Concurrency, expensive.Priority = StatusActive, true, 1, 0
	accounts := []Account{*cheap, *expensiveplaceholder
	groupID := int64(1)

	tests := []struct {
		name      string
		enabled   bool
		loadBatch bool
		loadErr   error
		wantID    int64
placeholder{
		{name: "switch off keeps priority first", loadBatch: true, wantID: 2placeholder,
		{name: "load batch", enabled: true, loadBatch: true, wantID: 1placeholder,
		{name: "load batch disabled", enabled: true, wantID: 1placeholder,
		{name: "load lookup failure", enabled: true, loadBatch: true, loadErr: errors.New("load unavailable"), wantID: 1placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetOpenAIAdvancedSchedulerSettingCacheForTest()
			settings := &openAIAdvancedSchedulerSettingRepoStub{values: map[string]string{
				openAIAdvancedSchedulerSettingKey:              "false",
				SettingKeyOpenAILowUpstreamRatePriorityEnabled: strconv.FormatBool(tt.enabled),
		placeholderplaceholder
			cfg := &config.Config{placeholder
			cfg.Gateway.Scheduling.LoadBatchEnabled = tt.loadBatch
			svc := &OpenAIGatewayService{
				accountRepo:      schedulerTestOpenAIAccountRepo{accounts: accountsplaceholder,
				cache:            &schedulerTestGatewayCache{placeholder,
				cfg:              cfg,
				rateLimitService: &RateLimitService{settingService: NewSettingService(settings, cfg)placeholder,
				concurrencyService: NewConcurrencyService(schedulerTestConcurrencyCache{
					loadBatchErr: tt.loadErr,
					loadMap: map[int64]*AccountLoadInfo{
						1: {AccountID: 1, LoadRate: 90placeholder,
						2: {AccountID: 2, LoadRate: 10placeholder,
				placeholder,
			placeholder),
		placeholder

			selection, _, err := svc.SelectAccountWithScheduler(context.Background(), &groupID, "", "", "gpt-test", nil, OpenAIUpstreamTransportAny, false)
		placeholder
			require.Equal(t, tt.wantID, selection.Account.ID)
			if selection.ReleaseFunc != nil {
				selection.ReleaseFunc()
		placeholder
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayServiceAdvancedSchedulerIgnoresLegacyLowRateSwitch(t *testing.T) {
	resetOpenAIAdvancedSchedulerSettingCacheForTest()
	defer resetOpenAIAdvancedSchedulerSettingCacheForTest()
	now := time.Now()
	cheap := upstreamCostTestAccount(1, UpstreamBillingProbeStatusOK, 0.03, now.Add(-time.Minute), 30*time.Minute)
	cheap.Status, cheap.Schedulable, cheap.Concurrency, cheap.Priority = StatusActive, true, 1, 10
	expensive := upstreamCostTestAccount(2, UpstreamBillingProbeStatusOK, 0.8, now.Add(-time.Minute), 30*time.Minute)
	expensive.Status, expensive.Schedulable, expensive.Concurrency, expensive.Priority = StatusActive, true, 1, 0
	settings := &openAIAdvancedSchedulerSettingRepoStub{values: map[string]string{
		openAIAdvancedSchedulerSettingKey:              "true",
		SettingKeyOpenAILowUpstreamRatePriorityEnabled: "true",
placeholderplaceholder
	cfg := &config.Config{placeholder
	cfg.Gateway.OpenAIWS.LBTopK = 1
	cfg.Gateway.OpenAIWS.SchedulerScoreWeights.Priority = 1
	svc := &OpenAIGatewayService{
		accountRepo:        schedulerTestOpenAIAccountRepo{accounts: []Account{*cheap, *expensiveplaceholderplaceholder,
		cache:              &schedulerTestGatewayCache{placeholder,
		cfg:                cfg,
		rateLimitService:   &RateLimitService{settingService: NewSettingService(settings, cfg)placeholder,
		concurrencyService: NewConcurrencyService(schedulerTestConcurrencyCache{placeholder),
placeholder
	groupID := int64(1)

	selection, _, err := svc.SelectAccountWithScheduler(context.Background(), &groupID, "", "", "gpt-test", nil, OpenAIUpstreamTransportAny, false)
placeholder
	require.Equal(t, int64(2), selection.Account.ID)
	if selection.ReleaseFunc != nil {
		selection.ReleaseFunc()
placeholder
placeholder

func TestOpenAIGatewayServiceLegacyLowRatePrioritySkipsCooledDownAccount(t *testing.T) {
	resetOpenAIAdvancedSchedulerSettingCacheForTest()
	defer resetOpenAIAdvancedSchedulerSettingCacheForTest()

	now := time.Now()
	cheap := upstreamCostTestAccount(1, UpstreamBillingProbeStatusOK, 0.03, now.Add(-time.Minute), 30*time.Minute)
	cheap.Status, cheap.Schedulable, cheap.Concurrency, cheap.Priority = StatusActive, true, 1, 10
	cooldownUntil := now.Add(time.Minute)
	cheap.TempUnschedulableUntil = &cooldownUntil
	expensive := upstreamCostTestAccount(2, UpstreamBillingProbeStatusOK, 0.8, now.Add(-time.Minute), 30*time.Minute)
	expensive.Status, expensive.Schedulable, expensive.Concurrency, expensive.Priority = StatusActive, true, 1, 0
	settings := &openAIAdvancedSchedulerSettingRepoStub{values: map[string]string{
		openAIAdvancedSchedulerSettingKey:              "false",
		SettingKeyOpenAILowUpstreamRatePriorityEnabled: "true",
placeholderplaceholder
	cfg := &config.Config{placeholder
	cfg.Gateway.Scheduling.LoadBatchEnabled = true
	svc := &OpenAIGatewayService{
		accountRepo:      schedulerTestOpenAIAccountRepo{accounts: []Account{*cheap, *expensiveplaceholderplaceholder,
		cache:            &schedulerTestGatewayCache{placeholder,
		cfg:              cfg,
		rateLimitService: &RateLimitService{settingService: NewSettingService(settings, cfg)placeholder,
		concurrencyService: NewConcurrencyService(schedulerTestConcurrencyCache{loadMap: map[int64]*AccountLoadInfo{
			1: {AccountID: 1placeholder,
			2: {AccountID: 2placeholder,
	placeholderplaceholder),
placeholder
	groupID := int64(1)

	selection, _, err := svc.SelectAccountWithScheduler(context.Background(), &groupID, "", "", "gpt-test", nil, OpenAIUpstreamTransportAny, false)
placeholder
	require.Equal(t, int64(2), selection.Account.ID)
	if selection.ReleaseFunc != nil {
		selection.ReleaseFunc()
placeholder
placeholder

func TestOpenAIFreshUpstreamBillingRateUsesFreshCachedSuccessOnly(t *testing.T) {
	now := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name    string
		account *Account
		wantOK  bool
placeholder{
		{name: "fresh", account: upstreamCostTestAccount(1, UpstreamBillingProbeStatusOK, 0.3, now.Add(-time.Minute), 30*time.Minute), wantOK: trueplaceholder,
		{name: "zero rate", account: upstreamCostTestAccount(2, UpstreamBillingProbeStatusOK, 0, now.Add(-time.Minute), 30*time.Minute), wantOK: trueplaceholder,
		{name: "transient failure with fresh cache", account: upstreamCostTestAccount(3, UpstreamBillingProbeStatusFailed, 0.3, now.Add(-time.Minute), 30*time.Minute), wantOK: trueplaceholder,
		{name: "stale", account: upstreamCostTestAccount(4, UpstreamBillingProbeStatusOK, 0.3, now.Add(-61*time.Minute), 30*time.Minute)placeholder,
		{name: "future", account: upstreamCostTestAccount(5, UpstreamBillingProbeStatusOK, 0.3, now.Add(time.Minute), 30*time.Minute)placeholder,
		{name: "unsupported", account: upstreamCostTestAccount(6, UpstreamBillingProbeStatusUnsupported, 0.3, now.Add(-time.Minute), 30*time.Minute)placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := openAIFreshUpstreamBillingRate(tt.account, now)
			require.Equal(t, tt.wantOK, ok)
	placeholder)
placeholder
placeholder

func TestBuildOpenAISelectionOrderIncludesOverflowOnlyForCostScheduling(t *testing.T) {
	scheduler := &defaultOpenAIAccountScheduler{placeholder
	candidates := []openAIAccountCandidateScore{
		{account: &Account{ID: 1placeholder, loadInfo: &AccountLoadInfo{placeholder, score: 3placeholder,
		{account: &Account{ID: 2placeholder, loadInfo: &AccountLoadInfo{placeholder, score: 2placeholder,
		{account: &Account{ID: 3placeholder, loadInfo: &AccountLoadInfo{placeholder, score: 1placeholder,
placeholder

	legacy := scheduler.buildOpenAISelectionOrder(OpenAIAccountScheduleRequest{placeholder, openAIAccountLoadPlan{
		candidates: candidates,
		topK:       1,
placeholder)
	require.Len(t, legacy, 1)

	costAware := scheduler.buildOpenAISelectionOrder(OpenAIAccountScheduleRequest{placeholder, openAIAccountLoadPlan{
		candidates:              candidates,
		topK:                    1,
		includeOverflowFallback: true,
placeholder)
	require.Equal(t, []int64{1, 2, 3placeholder, []int64{
		costAware[0].account.ID,
		costAware[1].account.ID,
		costAware[2].account.ID,
placeholder)
placeholder

func TestBuildOpenAIAccountLoadPlanUsesCostOnlyForTokenScope(t *testing.T) {
	resetOpenAIAdvancedSchedulerSettingCacheForTest()
	defer resetOpenAIAdvancedSchedulerSettingCacheForTest()

	now := time.Now()
	accounts := []*Account{
		upstreamCostTestAccount(1, UpstreamBillingProbeStatusOK, 0.03, now.Add(-time.Minute), 30*time.Minute),
		upstreamCostTestOAuthAccount(2),
		upstreamCostTestAccount(3, UpstreamBillingProbeStatusOK, 0.8, now.Add(-time.Minute), 30*time.Minute),
placeholder
	cfg := &config.Config{placeholder
	cfg.Gateway.OpenAIWS.LBTopK = 1
	cfg.Gateway.OpenAIWS.SchedulerScoreWeights.UpstreamCost = 1.5
	settings := &openAIAdvancedSchedulerSettingRepoStub{values: map[string]string{
		SettingKeyOpenAIOAuthSchedulingRateMultiplier: "0.05",
placeholderplaceholder
	scheduler := &defaultOpenAIAccountScheduler{service: &OpenAIGatewayService{
		cfg:              cfg,
		rateLimitService: &RateLimitService{settingService: NewSettingService(settings, cfg)placeholder,
placeholderplaceholder
	loadMap := map[int64]*AccountLoadInfo{
		1: {AccountID: 1placeholder,
		2: {AccountID: 2placeholder,
		3: {AccountID: 3placeholder,
placeholder

	tokenPlan := scheduler.buildOpenAIAccountLoadPlan(context.Background(), OpenAIAccountScheduleRequest{UseUpstreamTokenCost: trueplaceholder, accounts, loadMap)
	require.Greater(t, tokenPlan.candidates[0].score, tokenPlan.candidates[1].score)
	require.Greater(t, tokenPlan.candidates[1].score, tokenPlan.candidates[2].score)
	require.True(t, tokenPlan.includeOverflowFallback)

	otherPlan := scheduler.buildOpenAIAccountLoadPlan(context.Background(), OpenAIAccountScheduleRequest{placeholder, accounts, loadMap)
	require.Equal(t, otherPlan.candidates[0].score, otherPlan.candidates[1].score)
	require.Equal(t, otherPlan.candidates[1].score, otherPlan.candidates[2].score)
	require.False(t, otherPlan.includeOverflowFallback)
placeholder

func TestBuildOpenAIAccountSchedulerScoreSnapshotUpstreamCostIsExactNoOpWithoutSignal(t *testing.T) {
	accounts := []*Account{
		{ID: 1, Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder,
		{ID: 2, Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder,
placeholder
	loadMap := map[int64]*AccountLoadInfo{
		1: {AccountID: 1, LoadRate: 20placeholder,
		2: {AccountID: 2, LoadRate: 80placeholder,
placeholder
	weights := GatewayOpenAIWSSchedulerScoreWeightsView{Priority: 1, Load: 1, Queue: 0.7, ErrorRate: 0.8, TTFT: 0.5placeholder
	baseline := buildOpenAIAccountSchedulerScoreSnapshot(accounts, loadMap, weights, false, defaultOpenAIOAuthSchedulingRateMultiplier)
	weights.UpstreamCost = 1.5
	withCost := buildOpenAIAccountSchedulerScoreSnapshot(accounts, loadMap, weights, false, defaultOpenAIOAuthSchedulingRateMultiplier)

	require.Equal(t, baseline, withCost)
placeholder

func TestBuildOpenAIAccountSchedulerScoreSnapshotUsesUpstreamCostSignal(t *testing.T) {
	now := time.Now()
	accounts := []*Account{
		upstreamCostTestAccount(1, UpstreamBillingProbeStatusOK, 0.03, now.Add(-time.Minute), 30*time.Minute),
		upstreamCostTestAccount(2, UpstreamBillingProbeStatusOK, 0.8, now.Add(-time.Minute), 30*time.Minute),
placeholder
	weights := GatewayOpenAIWSSchedulerScoreWeightsView{UpstreamCost: placeholder
	scores := buildOpenAIAccountSchedulerScoreSnapshot(accounts, nil, weights, false, defaultOpenAIOAuthSchedulingRateMultiplier)

	require.Greater(t, scores[1].BaseScore, scores[2].BaseScore)
placeholder
