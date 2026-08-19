//go:build unit

package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type fullRebuildLifecycleCache struct {
	*groupLifecycleTestCache

	mu              sync.Mutex
	captureAttempts []SchedulerBucket
	captureErrors   map[string]error
	lockBusyOnce    map[string]bool
	watermark       int64
	watermarkWrites []int64
placeholder

func newFullRebuildLifecycleCache(buckets ...SchedulerBucket) *fullRebuildLifecycleCache {
	return &fullRebuildLifecycleCache{
		groupLifecycleTestCache: newGroupLifecycleTestCache(buckets...),
		captureErrors:           make(map[string]error),
		lockBusyOnce:            make(map[string]bool),
placeholder
placeholder

func (c *fullRebuildLifecycleCache) CaptureBucketWriteToken(ctx context.Context, bucket SchedulerBucket) (SchedulerBucketWriteToken, error) {
	c.mu.Lock()
	c.captureAttempts = append(c.captureAttempts, bucket)
	err := c.captureErrors[bucket.String()]
	c.mu.Unlock()
	if err != nil {
		return SchedulerBucketWriteToken{placeholder, err
placeholder
	return c.retirementRaceCache.CaptureBucketWriteToken(ctx, bucket)
placeholder

func (c *fullRebuildLifecycleCache) ListBuckets(ctx context.Context) ([]SchedulerBucket, error) {
	buckets, err := c.groupLifecycleTestCache.ListBuckets(ctx)
	if err != nil {
		return nil, err
placeholder
	c.retirementRaceCache.mu.Lock()
	defer c.retirementRaceCache.mu.Unlock()
	registered := make([]SchedulerBucket, 0, len(buckets))
	for _, bucket := range buckets {
		if !c.retirementRaceCache.retired[bucket.String()] {
			registered = append(registered, bucket)
	placeholder
placeholder
	return registered, nil
placeholder

func (c *fullRebuildLifecycleCache) TryLockBucket(ctx context.Context, bucket SchedulerBucket, ttl time.Duration) (bool, error) {
	c.mu.Lock()
	busy := c.lockBusyOnce[bucket.String()]
	if busy {
		delete(c.lockBusyOnce, bucket.String())
placeholder
	c.mu.Unlock()
	if busy {
		return false, nil
placeholder
	return c.groupLifecycleTestCache.TryLockBucket(ctx, bucket, ttl)
placeholder

func (c *fullRebuildLifecycleCache) GetOutboxWatermark(context.Context) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.watermark, nil
placeholder

func (c *fullRebuildLifecycleCache) SetOutboxWatermark(_ context.Context, id int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.watermark = id
	c.watermarkWrites = append(c.watermarkWrites, id)
	return nil
placeholder

func (c *fullRebuildLifecycleCache) captureAttemptCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.captureAttempts)
placeholder

func (c *fullRebuildLifecycleCache) currentWatermark() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.watermark
placeholder

func (c *fullRebuildLifecycleCache) totalSetAttempts() int {
	c.retirementRaceCache.mu.Lock()
	defer c.retirementRaceCache.mu.Unlock()
	var total int
	for _, attempts := range c.retirementRaceCache.setAttempts {
		total += attempts
placeholder
	return total
placeholder

type fullRebuildLifecycleGroupRepo struct {
	GroupRepository

	mu              sync.Mutex
	activeIDs       []int64
	activeIDsErr    error
	listActiveErr   error
	fresh           map[int64]*Group
	freshErr        map[int64]error
	activeIDCalls   int
	listActiveCalls int
	freshCalls      []int64
placeholder

func (r *fullRebuildLifecycleGroupRepo) ListActiveIDs(context.Context) ([]int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.activeIDCalls++
	return append([]int64(nil), r.activeIDs...), r.activeIDsErr
placeholder

func (r *fullRebuildLifecycleGroupRepo) ListActive(context.Context) ([]Group, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.listActiveCalls++
	return nil, r.listActiveErr
placeholder

func (r *fullRebuildLifecycleGroupRepo) GetByIDLite(_ context.Context, id int64) (*Group, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.freshCalls = append(r.freshCalls, id)
	if err := r.freshErr[id]; err != nil {
		return nil, err
placeholder
	group := r.fresh[id]
	if group == nil {
		return nil, ErrGroupNotFound
placeholder
	copyGroup := *group
	return &copyGroup, nil
placeholder

func (r *fullRebuildLifecycleGroupRepo) stats() (activeIDs, listActive int, fresh []int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.activeIDCalls, r.listActiveCalls, append([]int64(nil), r.freshCalls...)
placeholder

type fullRebuildFallbackGroupRepo struct {
	GroupRepository

	mu        sync.Mutex
	groups    []Group
	err       error
	listCalls int
placeholder

func (r *fullRebuildFallbackGroupRepo) ListActive(context.Context) ([]Group, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.listCalls++
	return append([]Group(nil), r.groups...), r.err
placeholder

type fullRebuildAccountCall struct {
	groupID  int64
	platform string
placeholder

type fullRebuildAccountRepo struct {
	AccountRepository

	mu          sync.Mutex
	calls       []fullRebuildAccountCall
	beforeFirst func()
	once        sync.Once
placeholder

func (r *fullRebuildAccountRepo) record(groupID int64, platform string) ([]Account, error) {
	r.mu.Lock()
	r.calls = append(r.calls, fullRebuildAccountCall{groupID: groupID, platform: platformplaceholder)
	beforeFirst := r.beforeFirst
	r.mu.Unlock()
	if beforeFirst != nil {
		r.once.Do(beforeFirst)
placeholder
	return []Account{{ID: 1, Platform: platform, Status: StatusActive, Schedulable: trueplaceholderplaceholder, nil
placeholder

func (r *fullRebuildAccountRepo) ListSchedulableByGroupIDAndPlatform(_ context.Context, groupID int64, platform string) ([]Account, error) {
	return r.record(groupID, platform)
placeholder

func (r *fullRebuildAccountRepo) ListSchedulableByGroupIDAndPlatforms(_ context.Context, groupID int64, platforms []string) ([]Account, error) {
	return r.record(groupID, firstPlatform(platforms))
placeholder

func (r *fullRebuildAccountRepo) ListSchedulableUngroupedByPlatform(_ context.Context, platform string) ([]Account, error) {
	return r.record(0, platform)
placeholder

func (r *fullRebuildAccountRepo) ListSchedulableUngroupedByPlatforms(_ context.Context, platforms []string) ([]Account, error) {
	return r.record(0, firstPlatform(platforms))
placeholder

func (r *fullRebuildAccountRepo) ListModelAvailabilityCandidates(context.Context, *int64, []string, bool) ([]Account, error) {
	panic("unexpected ListModelAvailabilityCandidates call")
placeholder

func (r *fullRebuildAccountRepo) ListSchedulableByPlatform(_ context.Context, platform string) ([]Account, error) {
	return r.record(0, platform)
placeholder

func (r *fullRebuildAccountRepo) ListSchedulableByPlatforms(_ context.Context, platforms []string) ([]Account, error) {
	return r.record(0, firstPlatform(platforms))
placeholder

func (r *fullRebuildAccountRepo) callCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.calls)
placeholder

func (r *fullRebuildAccountRepo) groupCallCount(groupID int64) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	var count int
	for _, call := range r.calls {
		if call.groupID == groupID {
			count++
	placeholder
placeholder
	return count
placeholder

func firstPlatform(platforms []string) string {
	if len(platforms) == 0 {
		return ""
placeholder
	return platforms[0]
placeholder

func newFullRebuildLifecycleService(
	cache SchedulerCache,
	outbox SchedulerOutboxRepository,
	accounts AccountRepository,
	groups GroupRepository,
	runMode string,
) *SchedulerSnapshotService {
	return NewSchedulerSnapshotService(cache, outbox, accounts, groups, &config.Config{RunMode: runModeplaceholder)
placeholder

func TestSchedulerFullRebuildActiveTombstoneDoesNotBlockFollowingGroupEvent(t *testing.T) {
	const groupID int64 = 101
	canonical := schedulerBucketsForGroup(groupID)
	cache := newFullRebuildLifecycleCache()
	require.NoError(t, cache.retirementRaceCache.RetireBucket(context.Background(), canonical[0]))
	groups := &fullRebuildLifecycleGroupRepo{
		activeIDs: []int64{groupIDplaceholder,
		fresh: map[int64]*Group{
			groupID: {ID: groupID, Status: StatusActive, Hydrated: trueplaceholder,
	placeholder,
		freshErr: make(map[int64]error),
placeholder
	accounts := &fullRebuildAccountRepo{placeholder
	outbox := &outboxCleanupRepo{events: []SchedulerOutboxEvent{
		{ID: 1, EventType: SchedulerOutboxEventFullRebuildplaceholder,
		{ID: 2, EventType: SchedulerOutboxEventGroupChanged, GroupID: ptrInt64(groupID)placeholder,
placeholderplaceholder
	svc := newFullRebuildLifecycleService(cache, outbox, accounts, groups, config.RunModeStandard)

	svc.pollOutbox()

	require.Equal(t, int64(2), cache.currentWatermark())
	cache.mu.Lock()
	require.Equal(t, []int64{2placeholder, cache.watermarkWrites)
	cache.mu.Unlock()
	activeCalls, fallbackCalls, freshCalls := groups.stats()
	require.Equal(t, 1, activeCalls)
	require.Zero(t, fallbackCalls)
	require.Equal(t, []int64{groupID, groupIDplaceholder, freshCalls)
	require.Len(t, cache.tokens(), 36, "full rebuild and the following group event must each run fresh authority")
	_, reopenHeld := cache.lifecycleMutationLeaseStates()
	require.Len(t, reopenHeld, 36)
	for _, held := range reopenHeld {
		require.True(t, held)
placeholder
	require.Equal(t, 30, accounts.callCount())
placeholder

func TestSchedulerFullRebuildGlobalReadErrorsFailBeforeMutationOrDB(t *testing.T) {
	t.Run("list buckets", func(t *testing.T) {
		cache := newFullRebuildLifecycleCache()
		cache.listErr = errors.New("registry failed")
		groups := &fullRebuildLifecycleGroupRepo{fresh: make(map[int64]*Group), freshErr: make(map[int64]error)placeholder
		accounts := &fullRebuildAccountRepo{placeholder
		svc := newFullRebuildLifecycleService(cache, nil, accounts, groups, config.RunModeStandard)

		err := svc.rebuildFullSnapshot(context.Background(), "test")
		require.ErrorIs(t, err, cache.listErr)
		activeCalls, fallbackCalls, freshCalls := groups.stats()
		require.Zero(t, activeCalls)
		require.Zero(t, fallbackCalls)
		require.Empty(t, freshCalls)
		requireFullRebuildNoMutationOrDB(t, cache, accounts)
placeholder)

	t.Run("list active ids without fallback", func(t *testing.T) {
		cache := newFullRebuildLifecycleCache()
		groups := &fullRebuildLifecycleGroupRepo{
			activeIDsErr:  errors.New("active ids failed"),
			listActiveErr: errors.New("fallback must not run"),
			fresh:         make(map[int64]*Group),
			freshErr:      make(map[int64]error),
	placeholder
		accounts := &fullRebuildAccountRepo{placeholder
		svc := newFullRebuildLifecycleService(cache, nil, accounts, groups, config.RunModeStandard)

		err := svc.rebuildFullSnapshot(context.Background(), "test")
		require.ErrorIs(t, err, groups.activeIDsErr)
		activeCalls, fallbackCalls, freshCalls := groups.stats()
		require.Equal(t, 1, activeCalls)
		require.Zero(t, fallbackCalls)
		require.Empty(t, freshCalls)
		requireFullRebuildNoMutationOrDB(t, cache, accounts)
placeholder)

	t.Run("list active fallback", func(t *testing.T) {
		cache := newFullRebuildLifecycleCache()
		groups := &fullRebuildFallbackGroupRepo{err: errors.New("active groups failed")placeholder
		accounts := &fullRebuildAccountRepo{placeholder
		svc := newFullRebuildLifecycleService(cache, nil, accounts, groups, config.RunModeStandard)

		err := svc.rebuildFullSnapshot(context.Background(), "test")
		require.ErrorIs(t, err, groups.err)
		groups.mu.Lock()
		require.Equal(t, 1, groups.listCalls)
		groups.mu.Unlock()
		requireFullRebuildNoMutationOrDB(t, cache, accounts)
placeholder)
placeholder

func TestSchedulerFullRebuildFreshActivePreparesEveryTokenBeforeFirstDB(t *testing.T) {
	const groupID int64 = 102
	historical := SchedulerBucket{GroupID: groupID, Platform: "legacy", Mode: "unknown"placeholder
	cache := newFullRebuildLifecycleCache(historical)
	groups := &fullRebuildLifecycleGroupRepo{
		fresh: map[int64]*Group{
			groupID: {ID: groupID, Status: StatusActive, Hydrated: trueplaceholder,
	placeholder,
		freshErr: make(map[int64]error),
placeholder
	accounts := &fullRebuildAccountRepo{placeholder
	var capturesAtFirstDB int
	accounts.beforeFirst = func() {
		capturesAtFirstDB = cache.captureAttemptCount()
		held, reopenCount := cache.leaseHeldAndTokenCount()
		require.False(t, held)
		require.Equal(t, 18, reopenCount)
		require.Equal(t, 19, capturesAtFirstDB, "C(0) and the historical bucket must be captured before DB")
placeholder
	svc := newFullRebuildLifecycleService(cache, nil, accounts, groups, config.RunModeStandard)

	require.NoError(t, svc.rebuildFullSnapshot(context.Background(), "test"))
	require.Equal(t, capturesAtFirstDB, cache.captureAttemptCount())
	require.Equal(t, 21, accounts.callCount())
	require.Equal(t, 11, accounts.groupCallCount(groupID))
	_, historicalPublished := cache.counts(historical)
	require.Equal(t, 1, historicalPublished)
	activeCalls, fallbackCalls, freshCalls := groups.stats()
	require.Equal(t, 1, activeCalls)
	require.Zero(t, fallbackCalls)
	require.Equal(t, []int64{groupIDplaceholder, freshCalls)
	_, _, listCalls := cache.lifecycleCounts()
	require.Equal(t, 1, listCalls, "known Rg must prevent a second registry read")
placeholder

func TestSchedulerFullRebuildOrdinaryCaptureErrorReturnsBeforeFirstDB(t *testing.T) {
	first := SchedulerBucket{GroupID: 0, Platform: "legacy-a", Mode: "unknown"placeholder
	last := SchedulerBucket{GroupID: 0, Platform: "legacy-b", Mode: "unknown"placeholder
	cache := newFullRebuildLifecycleCache(first, last)
	wantErr := errors.New("capture failed")
	cache.captureErrors[last.String()] = wantErr
	groups := &fullRebuildLifecycleGroupRepo{fresh: make(map[int64]*Group), freshErr: make(map[int64]error)placeholder
	accounts := &fullRebuildAccountRepo{placeholder
	svc := newFullRebuildLifecycleService(cache, nil, accounts, groups, config.RunModeStandard)

	err := svc.rebuildFullSnapshot(context.Background(), "test")
	require.ErrorIs(t, err, wantErr)
	require.Equal(t, 20, cache.captureAttemptCount(), "all canonical and ordinary captures must be attempted before returning")
	require.Zero(t, accounts.callCount())
	require.Zero(t, cache.totalSetAttempts())
placeholder

func TestSchedulerFullRebuildPreservesGroupZeroActiveHistoricalAndInvalidRegistryBuckets(t *testing.T) {
	const groupID int64 = 103
	groupZeroHistorical := SchedulerBucket{GroupID: 0, Platform: "legacy-zero", Mode: "unknown"placeholder
	activeHistorical := SchedulerBucket{GroupID: groupID, Platform: "legacy-active", Mode: "unknown"placeholder
	activeForced := SchedulerBucket{GroupID: groupID, Platform: PlatformOpenAI, Mode: SchedulerModeForcedplaceholder
	activeMixed := SchedulerBucket{GroupID: groupID, Platform: PlatformAnthropic, Mode: SchedulerModeMixedplaceholder
	invalidHistorical := SchedulerBucket{GroupID: -7, Platform: "legacy-invalid", Mode: "unknown"placeholder
	cache := newFullRebuildLifecycleCache(groupZeroHistorical, activeHistorical, activeForced, activeMixed, invalidHistorical)
	groups := &fullRebuildFallbackGroupRepo{groups: []Group{{ID: groupID, Status: StatusActiveplaceholderplaceholderplaceholder
	accounts := &fullRebuildAccountRepo{placeholder
	svc := newFullRebuildLifecycleService(cache, nil, accounts, groups, config.RunModeStandard)

	require.NoError(t, svc.rebuildFullSnapshot(context.Background(), "test"))
	require.Equal(t, 39, cache.captureAttemptCount())
	require.Equal(t, 23, accounts.callCount())
	groups.mu.Lock()
	require.Equal(t, 1, groups.listCalls)
	groups.mu.Unlock()
	require.Empty(t, cache.retiredBuckets())
	require.Empty(t, cache.tokens())
	for _, bucket := range []SchedulerBucket{groupZeroHistorical, activeHistorical, activeForced, activeMixed, invalidHistoricalplaceholder {
		_, published := cache.counts(bucket)
		require.Equal(t, 1, published, bucket.String())
placeholder
placeholder

func TestSchedulerFullRebuildActiveTombstoneFreshInactiveOrMissingFiltersAllGroupTasks(t *testing.T) {
	for _, tc := range []struct {
		name  string
		group *Group
placeholder{
		{name: "inactive", group: &Group{ID: 104, Status: StatusDisabled, Hydrated: trueplaceholderplaceholder,
		{name: "missing"placeholder,
placeholder {
		t.Run(tc.name, func(t *testing.T) {
			const groupID int64 = 104
			canonical := schedulerBucketsForGroup(groupID)
			historical := SchedulerBucket{GroupID: groupID, Platform: "legacy", Mode: "unknown"placeholder
			cache := newFullRebuildLifecycleCache(historical)
			require.NoError(t, cache.retirementRaceCache.RetireBucket(context.Background(), canonical[5]))
			groups := &fullRebuildLifecycleGroupRepo{
				activeIDs: []int64{groupIDplaceholder,
				fresh:     make(map[int64]*Group),
				freshErr:  make(map[int64]error),
		placeholder
			if tc.group != nil {
				groups.fresh[groupID] = tc.group
		placeholder
			accounts := &fullRebuildAccountRepo{placeholder
			svc := newFullRebuildLifecycleService(cache, nil, accounts, groups, config.RunModeStandard)

			require.NoError(t, svc.rebuildFullSnapshot(context.Background(), "test"))
			require.Zero(t, accounts.groupCallCount(groupID))
			require.Equal(t, 10, accounts.groupCallCount(0))
			require.Empty(t, cache.tokens())
			require.Equal(t, bucketStrings(append(canonical, historical)), bucketStrings(cache.retiredBuckets()))
			for _, bucket := range append(canonical, historical) {
				attempts, published := cache.counts(bucket)
				require.Zero(t, attempts, bucket.String())
				require.Zero(t, published, bucket.String())
		placeholder
			_, _, listCalls := cache.lifecycleCounts()
			require.Equal(t, 1, listCalls)
	placeholder)
placeholder
placeholder

func TestSchedulerFullRebuildStaleCandidatesAreSortedAndNeverReopenedAcrossRounds(t *testing.T) {
	historical := []SchedulerBucket{
		{GroupID: 3, Platform: "legacy", Mode: "unknown"placeholder,
		{GroupID: 1, Platform: "legacy", Mode: "unknown"placeholder,
		{GroupID: 2, Platform: "legacy", Mode: "unknown"placeholder,
placeholder
	cache := newFullRebuildLifecycleCache(historical...)
	groups := &fullRebuildLifecycleGroupRepo{
		fresh: map[int64]*Group{
			1: {ID: 1, Status: StatusDisabled, Hydrated: trueplaceholder,
			3: {ID: 3, Status: StatusDisabled, Hydrated: trueplaceholder,
	placeholder,
		freshErr: make(map[int64]error),
placeholder
	accounts := &fullRebuildAccountRepo{placeholder
	svc := newFullRebuildLifecycleService(cache, nil, accounts, groups, config.RunModeStandard)

	require.NoError(t, svc.rebuildFullSnapshot(context.Background(), "first"))
	retireCallsAfterFirst := len(cache.retiredBuckets())
	require.NoError(t, svc.rebuildFullSnapshot(context.Background(), "second"))
	_, _, freshCalls := groups.stats()
	require.Equal(t, []int64{1, 2, 3placeholder, freshCalls)
	require.Equal(t, retireCallsAfterFirst, len(cache.retiredBuckets()))
	require.Empty(t, cache.tokens())
	require.Zero(t, accounts.groupCallCount(1))
	require.Zero(t, accounts.groupCallCount(2))
	require.Zero(t, accounts.groupCallCount(3))
	_, _, listCalls := cache.lifecycleCounts()
	require.Equal(t, 2, listCalls, "each round must use only its one global registry snapshot")
	for _, bucket := range historical {
		require.Contains(t, bucketStrings(cache.retiredBuckets()), bucket.String())
placeholder
placeholder

func TestSchedulerFullRebuildPartialLifecycleFailureReturnsBeforeDBAndRetries(t *testing.T) {
	historical := []SchedulerBucket{
		{GroupID: 1, Platform: "legacy", Mode: "unknown"placeholder,
		{GroupID: 2, Platform: "legacy", Mode: "unknown"placeholder,
		{GroupID: 3, Platform: "legacy", Mode: "unknown"placeholder,
placeholder
	wantErr := errors.New("fresh group query failed")
	cache := newFullRebuildLifecycleCache(historical...)
	groups := &fullRebuildLifecycleGroupRepo{
		fresh: map[int64]*Group{
			1: {ID: 1, Status: StatusDisabled, Hydrated: trueplaceholder,
			3: {ID: 3, Status: StatusDisabled, Hydrated: trueplaceholder,
	placeholder,
		freshErr: map[int64]error{2: wantErrplaceholder,
placeholder
	accounts := &fullRebuildAccountRepo{placeholder
	svc := newFullRebuildLifecycleService(cache, nil, accounts, groups, config.RunModeStandard)

	err := svc.triggerFullRebuild("first")
	require.ErrorIs(t, err, wantErr)
	_, _, freshCalls := groups.stats()
	require.Equal(t, []int64{1, 2placeholder, freshCalls)
	require.Zero(t, accounts.callCount())
	require.Zero(t, cache.totalSetAttempts())
	require.Equal(t, 19, len(cache.retiredBuckets()))

	groups.mu.Lock()
	delete(groups.freshErr, 2)
	groups.fresh[2] = &Group{ID: 2, Status: StatusDisabled, Hydrated: trueplaceholder
	groups.mu.Unlock()
	require.NoError(t, svc.triggerFullRebuild("retry"))
	_, _, freshCalls = groups.stats()
	require.Equal(t, []int64{1, 2, 2, 3placeholder, freshCalls)
	require.Equal(t, 57, len(cache.retiredBuckets()))
	require.Equal(t, 10, accounts.callCount())
	require.Empty(t, cache.tokens())
placeholder

func TestSchedulerFullRebuildActiveTombstoneLazyRecoveryDiscardsPartialCaptureTasks(t *testing.T) {
	const groupID int64 = 105
	canonical := schedulerBucketsForGroup(groupID)
	historical := SchedulerBucket{GroupID: groupID, Platform: "legacy", Mode: "unknown"placeholder
	cache := newFullRebuildLifecycleCache(canonical[0], canonical[4], historical)
	require.NoError(t, cache.retirementRaceCache.RetireBucket(context.Background(), canonical[5]))
	groups := &fullRebuildLifecycleGroupRepo{
		activeIDs: []int64{groupIDplaceholder,
		fresh: map[int64]*Group{
			groupID: {ID: groupID, Status: StatusActive, Hydrated: trueplaceholder,
	placeholder,
		freshErr: make(map[int64]error),
placeholder
	accounts := &fullRebuildAccountRepo{placeholder
	var capturesAtFirstDB int
	accounts.beforeFirst = func() {
		capturesAtFirstDB = cache.captureAttemptCount()
		held, reopenCount := cache.leaseHeldAndTokenCount()
		require.False(t, held)
		require.Equal(t, 18, reopenCount)
		require.Equal(t, 25, capturesAtFirstDB)
placeholder
	svc := newFullRebuildLifecycleService(cache, nil, accounts, groups, config.RunModeStandard)

	require.NoError(t, svc.rebuildFullSnapshot(context.Background(), "test"))
	require.Equal(t, capturesAtFirstDB, cache.captureAttemptCount())
	require.Equal(t, 21, accounts.callCount())
	for _, bucket := range canonical {
		attempts, published := cache.counts(bucket)
		require.Equal(t, 1, attempts, "discarded pre-recovery tokens must never publish: %s", bucket.String())
		require.Equal(t, 1, published, bucket.String())
placeholder
	_, published := cache.counts(historical)
	require.Equal(t, 1, published)
	_, _, listCalls := cache.lifecycleCounts()
	require.Equal(t, 1, listCalls)
placeholder

func TestSchedulerFullRebuildSimpleModePreservesRegistryWithoutLifecycleAuthority(t *testing.T) {
	registered := []SchedulerBucket{
		{GroupID: 0, Platform: "legacy-zero", Mode: "unknown"placeholder,
		{GroupID: 106, Platform: "legacy-positive", Mode: "unknown"placeholder,
		{GroupID: -8, Platform: "legacy-negative", Mode: "unknown"placeholder,
placeholder
	cache := newFullRebuildLifecycleCache(registered...)
	groups := &fullRebuildLifecycleGroupRepo{
		activeIDsErr: errors.New("simple mode must not query groups"),
		fresh:        make(map[int64]*Group),
		freshErr:     make(map[int64]error),
placeholder
	accounts := &fullRebuildAccountRepo{placeholder
	svc := newFullRebuildLifecycleService(cache, nil, accounts, groups, config.RunModeSimple)

	require.NoError(t, svc.rebuildFullSnapshot(context.Background(), "test"))
	activeCalls, fallbackCalls, freshCalls := groups.stats()
	require.Zero(t, activeCalls)
	require.Zero(t, fallbackCalls)
	require.Empty(t, freshCalls)
	require.Equal(t, 21, cache.captureAttemptCount())
	require.Equal(t, 13, accounts.callCount())
	require.Equal(t, 13, accounts.groupCallCount(0))
	require.Empty(t, cache.retiredBuckets())
	require.Empty(t, cache.tokens())
	for _, bucket := range registered {
		_, published := cache.counts(bucket)
		require.Equal(t, 1, published, bucket.String())
placeholder
placeholder

func TestSchedulerFullRebuildFreshReopenLockBusyRetriesWithoutBlockingOrdinaryTasks(t *testing.T) {
	const groupID int64 = 107
	canonical := schedulerBucketsForGroup(groupID)
	cache := newFullRebuildLifecycleCache()
	require.NoError(t, cache.retirementRaceCache.RetireBucket(context.Background(), canonical[0]))
	cache.lockBusyOnce[canonical[0].String()] = true
	groups := &fullRebuildLifecycleGroupRepo{
		activeIDs: []int64{groupIDplaceholder,
		fresh: map[int64]*Group{
			groupID: {ID: groupID, Status: StatusActive, Hydrated: trueplaceholder,
	placeholder,
		freshErr: make(map[int64]error),
placeholder
	accounts := &fullRebuildAccountRepo{placeholder
	outbox := &outboxCleanupRepo{events: []SchedulerOutboxEvent{{ID: 1, EventType: SchedulerOutboxEventFullRebuildplaceholderplaceholderplaceholder
	svc := newFullRebuildLifecycleService(cache, outbox, accounts, groups, config.RunModeStandard)

	svc.pollOutbox()
	require.Zero(t, cache.currentWatermark())
	_, groupZeroPublished := cache.counts(schedulerCanonicalBuckets(0)[0])
	require.Equal(t, 1, groupZeroPublished, "ordinary tasks must still run when one strict Reopen task is busy")
	require.Equal(t, 20, accounts.callCount())

	svc.pollOutbox()
	require.Equal(t, int64(1), cache.currentWatermark())
	require.Equal(t, 40, accounts.callCount())
	_, busyBucketPublished := cache.counts(canonical[0])
	require.Equal(t, 1, busyBucketPublished)
	activeCalls, fallbackCalls, freshCalls := groups.stats()
	require.Equal(t, 2, activeCalls)
	require.Zero(t, fallbackCalls)
	require.Equal(t, []int64{groupIDplaceholder, freshCalls)
placeholder

func TestSchedulerFullRebuildOrdinaryLockBusyKeepsExistingSkipSemantics(t *testing.T) {
	busyBucket := schedulerCanonicalBuckets(0)[0]
	cache := newFullRebuildLifecycleCache()
	cache.lockBusyOnce[busyBucket.String()] = true
	groups := &fullRebuildLifecycleGroupRepo{fresh: make(map[int64]*Group), freshErr: make(map[int64]error)placeholder
	accounts := &fullRebuildAccountRepo{placeholder
	svc := newFullRebuildLifecycleService(cache, nil, accounts, groups, config.RunModeStandard)

	require.NoError(t, svc.rebuildFullSnapshot(context.Background(), "test"))
	require.Equal(t, 10, accounts.callCount())
	attempts, published := cache.counts(busyBucket)
	require.Zero(t, attempts)
	require.Zero(t, published)
placeholder

func requireFullRebuildNoMutationOrDB(t *testing.T, cache *fullRebuildLifecycleCache, accounts *fullRebuildAccountRepo) {
placeholder
	require.Zero(t, cache.captureAttemptCount())
	require.Empty(t, cache.retiredBuckets())
	require.Empty(t, cache.tokens())
	require.Zero(t, accounts.callCount())
	require.Zero(t, cache.totalSetAttempts())
placeholder
