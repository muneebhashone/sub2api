//go:build unit

package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type groupLifecycleTestCache struct {
	*retirementRaceCache

	stateMu sync.Mutex

	leaseHeld       bool
	lease           SchedulerGroupLifecycleLease
	leaseSequence   int
	leaseBusy       bool
	leaseAcquireErr error
	leaseReleaseErr error
	acquireCalls    int
	releaseCalls    int
	acquireTTL      time.Duration
	acquireDeadline bool
	releaseDeadline bool
	releaseCtxErr   error

	listErr   error
	listCalls int

	retireCalls  []SchedulerBucket
	reopenTokens []SchedulerBucketWriteToken
	retireHeld   []bool
	reopenHeld   []bool
	retireErr    error
	retireErrAt  int
	reopenErr    error
	reopenErrAt  int

	bucketLockBusy bool
	bucketLockErr  error
	bucketLockTTLs []time.Duration
	unlockCalls    int
	setErr         error
placeholder

func newGroupLifecycleTestCache(buckets ...SchedulerBucket) *groupLifecycleTestCache {
	return &groupLifecycleTestCache{retirementRaceCache: newRetirementRaceCache(buckets...)placeholder
placeholder

func (c *groupLifecycleTestCache) TryAcquireGroupLifecycleLease(ctx context.Context, groupID int64, ttl time.Duration) (SchedulerGroupLifecycleLease, bool, error) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.acquireCalls++
	c.acquireTTL = ttl
	_, c.acquireDeadline = ctx.Deadline()
	if c.leaseAcquireErr != nil {
		return SchedulerGroupLifecycleLease{placeholder, false, c.leaseAcquireErr
placeholder
	if c.leaseBusy || c.leaseHeld {
		return SchedulerGroupLifecycleLease{placeholder, false, nil
placeholder
	c.leaseSequence++
	c.lease = SchedulerGroupLifecycleLease{GroupID: groupID, OwnerToken: fmt.Sprintf("owner-%d", c.leaseSequence)placeholder
	c.leaseHeld = true
	return c.lease, true, nil
placeholder

func (c *groupLifecycleTestCache) ReleaseGroupLifecycleLease(ctx context.Context, lease SchedulerGroupLifecycleLease) error {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.releaseCalls++
	_, c.releaseDeadline = ctx.Deadline()
	c.releaseCtxErr = ctx.Err()
	if c.leaseReleaseErr != nil {
		return c.leaseReleaseErr
placeholder
	if !c.leaseHeld || lease != c.lease {
		return ErrSchedulerGroupLifecycleLeaseLost
placeholder
	c.leaseHeld = false
	return nil
placeholder

func (c *groupLifecycleTestCache) RetireBucket(ctx context.Context, bucket SchedulerBucket) error {
	c.stateMu.Lock()
	c.retireCalls = append(c.retireCalls, bucket)
	c.retireHeld = append(c.retireHeld, c.leaseHeld)
	held := c.leaseHeld
	call := len(c.retireCalls)
	err := c.retireErr
	errAt := c.retireErrAt
	c.stateMu.Unlock()
	if !held {
		return errors.New("retire called outside group lifecycle lease")
placeholder
	if err != nil && (errAt <= 0 || call == errAt) {
		return err
placeholder
	return c.retirementRaceCache.RetireBucket(ctx, bucket)
placeholder

func (c *groupLifecycleTestCache) ReopenBucket(ctx context.Context, bucket SchedulerBucket) (SchedulerBucketWriteToken, error) {
	if err := ctx.Err(); err != nil {
		return SchedulerBucketWriteToken{placeholder, err
placeholder
	c.stateMu.Lock()
	c.reopenHeld = append(c.reopenHeld, c.leaseHeld)
	held := c.leaseHeld
	call := len(c.reopenHeld)
	reopenErr := c.reopenErr
	reopenErrAt := c.reopenErrAt
	c.stateMu.Unlock()
	if !held {
		return SchedulerBucketWriteToken{placeholder, errors.New("reopen called outside group lifecycle lease")
placeholder
	if reopenErr != nil && (reopenErrAt <= 0 || call == reopenErrAt) {
		return SchedulerBucketWriteToken{placeholder, reopenErr
placeholder
	token, err := c.retirementRaceCache.ReopenBucket(ctx, bucket)
	if err != nil {
		return SchedulerBucketWriteToken{placeholder, err
placeholder
	c.stateMu.Lock()
	c.reopenTokens = append(c.reopenTokens, token)
	c.stateMu.Unlock()
	return token, nil
placeholder

func (c *groupLifecycleTestCache) ListBuckets(ctx context.Context) ([]SchedulerBucket, error) {
	c.stateMu.Lock()
	c.listCalls++
	err := c.listErr
	c.stateMu.Unlock()
	if err != nil {
		return nil, err
placeholder
	return c.retirementRaceCache.ListBuckets(ctx)
placeholder

func (c *groupLifecycleTestCache) TryLockBucket(_ context.Context, _ SchedulerBucket, ttl time.Duration) (bool, error) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.bucketLockTTLs = append(c.bucketLockTTLs, ttl)
	if c.bucketLockErr != nil {
		return false, c.bucketLockErr
placeholder
	return !c.bucketLockBusy, nil
placeholder

func (c *groupLifecycleTestCache) UnlockBucket(context.Context, SchedulerBucket) error {
	c.stateMu.Lock()
	c.unlockCalls++
	c.stateMu.Unlock()
	return nil
placeholder

func (c *groupLifecycleTestCache) SetSnapshot(ctx context.Context, bucket SchedulerBucket, token SchedulerBucketWriteToken, accounts []Account) error {
	c.stateMu.Lock()
	err := c.setErr
	c.stateMu.Unlock()
	if err != nil {
		return err
placeholder
	return c.retirementRaceCache.SetSnapshot(ctx, bucket, token, accounts)
placeholder

func (c *groupLifecycleTestCache) lifecycleCounts() (acquires, releases, listCalls int) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	return c.acquireCalls, c.releaseCalls, c.listCalls
placeholder

func (c *groupLifecycleTestCache) retiredBuckets() []SchedulerBucket {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	return append([]SchedulerBucket(nil), c.retireCalls...)
placeholder

func (c *groupLifecycleTestCache) tokens() []SchedulerBucketWriteToken {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	return append([]SchedulerBucketWriteToken(nil), c.reopenTokens...)
placeholder

func (c *groupLifecycleTestCache) leaseHeldAndTokenCount() (bool, int) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	return c.leaseHeld, len(c.reopenTokens)
placeholder

func (c *groupLifecycleTestCache) lockStats() ([]time.Duration, int) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	return append([]time.Duration(nil), c.bucketLockTTLs...), c.unlockCalls
placeholder

func (c *groupLifecycleTestCache) lifecycleMutationLeaseStates() (retire, reopen []bool) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	return append([]bool(nil), c.retireHeld...), append([]bool(nil), c.reopenHeld...)
placeholder

type groupLifecycleTestGroupRepo struct {
	GroupRepository

	mu       sync.Mutex
	group    *Group
	err      error
	calls    int
	afterGet func()
placeholder

func (r *groupLifecycleTestGroupRepo) GetByIDLite(context.Context, int64) (*Group, error) {
	r.mu.Lock()
	r.calls++
	if r.err != nil {
		err := r.err
		r.mu.Unlock()
		return nil, err
placeholder
	if r.group == nil {
		r.mu.Unlock()
		return nil, ErrGroupNotFound
placeholder
	copyGroup := *r.group
	afterGet := r.afterGet
	r.mu.Unlock()
	if afterGet != nil {
		afterGet()
placeholder
	return &copyGroup, nil
placeholder

func (r *groupLifecycleTestGroupRepo) set(group *Group, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.group = group
	r.err = err
placeholder

func (r *groupLifecycleTestGroupRepo) callCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.calls
placeholder

type groupLifecycleTestAccountRepo struct {
	AccountRepository

	mu              sync.Mutex
	calls           int
	callsByPlatform map[string]int
	err             error
	started         chan struct{placeholder
	release         chan struct{placeholder
	once            sync.Once
	beforeLoad      func()
	beforeLoadOnce  sync.Once
placeholder

func (r *groupLifecycleTestAccountRepo) load(ctx context.Context, platform string) ([]Account, error) {
	r.mu.Lock()
	r.calls++
	if r.callsByPlatform == nil {
		r.callsByPlatform = make(map[string]int)
placeholder
	r.callsByPlatform[platform]++
	err := r.err
	started := r.started
	release := r.release
	r.mu.Unlock()
	if started != nil {
		r.once.Do(func() { close(started) placeholder)
placeholder
	if r.beforeLoad != nil {
		r.beforeLoadOnce.Do(r.beforeLoad)
placeholder
	if release != nil {
		select {
		case <-release:
		case <-ctx.Done():
			return nil, ctx.Err()
	placeholder
placeholder
	if err != nil {
		return nil, err
placeholder
	return []Account{{ID: 9001, Platform: platform, Status: StatusActive, Schedulable: trueplaceholderplaceholder, nil
placeholder

func (r *groupLifecycleTestAccountRepo) ListSchedulableByGroupIDAndPlatform(ctx context.Context, _ int64, platform string) ([]Account, error) {
	return r.load(ctx, platform)
placeholder

func (r *groupLifecycleTestAccountRepo) ListSchedulableByGroupIDAndPlatforms(ctx context.Context, _ int64, platforms []string) ([]Account, error) {
	platform := "mixed"
	if len(platforms) > 0 {
		platform = platforms[0]
placeholder
	return r.load(ctx, platform)
placeholder

func (r *groupLifecycleTestAccountRepo) callCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.calls
placeholder

func (r *groupLifecycleTestAccountRepo) platformCallCount(platform string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.callsByPlatform[platform]
placeholder

func newGroupLifecycleTestService(cache SchedulerCache, accounts AccountRepository, groups GroupRepository, runMode string) *SchedulerSnapshotService {
	return NewSchedulerSnapshotService(cache, nil, accounts, groups, &config.Config{RunMode: runModeplaceholder)
placeholder

func expectedGroupLifecycleBuckets(groupID int64) []SchedulerBucket {
	platforms := schedulerSnapshotPlatforms()
	buckets := make([]SchedulerBucket, 0, 18)
	for _, platform := range platforms {
		buckets = append(buckets,
			SchedulerBucket{GroupID: groupID, Platform: platform, Mode: SchedulerModeSingleplaceholder,
			SchedulerBucket{GroupID: groupID, Platform: platform, Mode: SchedulerModeForcedplaceholder,
		)
		if platform == PlatformAnthropic || platform == PlatformGemini {
			buckets = append(buckets, SchedulerBucket{GroupID: groupID, Platform: platform, Mode: SchedulerModeMixedplaceholder)
	placeholder
placeholder
	return buckets
placeholder

func bucketStrings(buckets []SchedulerBucket) map[string]struct{placeholder {
	out := make(map[string]struct{placeholder, len(buckets))
	for _, bucket := range buckets {
		out[bucket.String()] = struct{placeholder{placeholder
placeholder
	return out
placeholder

func requireLifecycleSeen(t *testing.T, seen map[batchSeenKey]struct{placeholder, groupID int64) {
placeholder
	_, ok := seen[batchSeenKey{groupID: groupID, lifecycle: trueplaceholder]
	require.True(t, ok)
	for _, platform := range schedulerSnapshotPlatforms() {
		_, ok = seen[batchSeenKey{groupID: groupID, platform: platformplaceholder]
		require.True(t, ok)
placeholder
placeholder

func requireLifecycleNotSeen(t *testing.T, seen map[batchSeenKey]struct{placeholder, groupID int64) {
placeholder
	_, ok := seen[batchSeenKey{groupID: groupID, lifecycle: trueplaceholder]
	require.False(t, ok)
	for _, platform := range schedulerSnapshotPlatforms() {
		_, ok = seen[batchSeenKey{groupID: groupID, platform: platformplaceholder]
		require.False(t, ok)
placeholder
placeholder

func TestSchedulerGroupLifecycleInactiveAndMissingRetireAllHistoricalBucketsWithoutAccountReads(t *testing.T) {
	for _, tc := range []struct {
		name  string
		group *Group
		err   error
placeholder{
		{name: "inactive", group: &Group{ID: 81, Status: StatusDisabled, Hydrated: trueplaceholderplaceholder,
		{name: "missing", err: ErrGroupNotFoundplaceholder,
placeholder {
		t.Run(tc.name, func(t *testing.T) {
			const groupID int64 = 81
			current := expectedGroupLifecycleBuckets(groupID)
			historical := SchedulerBucket{GroupID: groupID, Platform: "legacy", Mode: "obsolete"placeholder
			other := SchedulerBucket{GroupID: groupID + 1, Platform: PlatformOpenAI, Mode: SchedulerModeForcedplaceholder
			groupZero := SchedulerBucket{GroupID: 0, Platform: PlatformOpenAI, Mode: SchedulerModeForcedplaceholder
			cache := newGroupLifecycleTestCache(current[0], historical, other, groupZero)
			groups := &groupLifecycleTestGroupRepo{group: tc.group, err: tc.errplaceholder
			accounts := &groupLifecycleTestAccountRepo{placeholder
			svc := newGroupLifecycleTestService(cache, accounts, groups, config.RunModeStandard)
			seen := make(map[batchSeenKey]struct{placeholder)

			require.NoError(t, svc.handleGroupEvent(context.Background(), ptrInt64(groupID), seen))

			expected := bucketStrings(append(current, historical))
			got := bucketStrings(cache.retiredBuckets())
			require.Equal(t, expected, got)
			retireHeld, _ := cache.lifecycleMutationLeaseStates()
			require.Len(t, retireHeld, len(expected))
			for _, held := range retireHeld {
				require.True(t, held)
		placeholder
			require.NotContains(t, got, other.String())
			require.NotContains(t, got, groupZero.String())
			require.Zero(t, accounts.callCount())
			require.Equal(t, 1, groups.callCount())
			_, _, listCalls := cache.lifecycleCounts()
			require.Equal(t, 1, listCalls)
			requireLifecycleSeen(t, seen, groupID)
	placeholder)
placeholder
placeholder

func TestSchedulerPrepareGroupLifecycleUsesKnownHistoricalBucketsWithoutListingRegistry(t *testing.T) {
	const groupID int64 = 811
	historical := SchedulerBucket{GroupID: groupID, Platform: "legacy", Mode: "obsolete"placeholder
	cache := newGroupLifecycleTestCache()
	cache.listErr = errors.New("registry must not be listed")
	groups := &groupLifecycleTestGroupRepo{group: &Group{ID: groupID, Status: StatusDisabled, Hydrated: trueplaceholderplaceholder
	accounts := &groupLifecycleTestAccountRepo{placeholder
	svc := newGroupLifecycleTestService(cache, accounts, groups, config.RunModeStandard)

	plan, err := svc.prepareGroupLifecycle(context.Background(), groupID, []SchedulerBucket{historicalplaceholder)
placeholder
	require.False(t, plan.active)
	require.Empty(t, plan.tasks)
	_, _, listCalls := cache.lifecycleCounts()
	require.Zero(t, listCalls)
	require.Contains(t, bucketStrings(cache.retiredBuckets()), historical.String())
	require.Zero(t, accounts.callCount())
placeholder

func TestSchedulerGroupLifecycleActiveReopensAndRebuildsAllCurrentBuckets(t *testing.T) {
	const groupID int64 = 82
	current := expectedGroupLifecycleBuckets(groupID)
	historical := SchedulerBucket{GroupID: groupID, Platform: "legacy", Mode: "obsolete"placeholder
	cache := newGroupLifecycleTestCache(historical)
	for _, bucket := range current {
		require.NoError(t, cache.retirementRaceCache.RetireBucket(context.Background(), bucket))
placeholder
	groups := &groupLifecycleTestGroupRepo{group: &Group{ID: groupID, Platform: PlatformOpenAI, Status: StatusActive, Hydrated: trueplaceholderplaceholder
	accounts := &groupLifecycleTestAccountRepo{placeholder
	accounts.beforeLoad = func() {
		held, tokenCount := cache.leaseHeldAndTokenCount()
		require.False(t, held, "the group lifecycle lease must be released before the first account query")
		require.Equal(t, 18, tokenCount, "all reopen tokens must be prepared before the first account query")
placeholder
	svc := newGroupLifecycleTestService(cache, accounts, groups, config.RunModeStandard)
	seen := make(map[batchSeenKey]struct{placeholder)

	require.NoError(t, svc.handleGroupEvent(context.Background(), ptrInt64(groupID), seen))

	require.Equal(t, bucketStrings(current), bucketStrings(cache.reopens))
	require.Empty(t, cache.retiredBuckets())
	registered, err := cache.retirementRaceCache.ListBuckets(context.Background())
placeholder
	require.Contains(t, bucketStrings(registered), historical.String())
	require.Len(t, cache.tokens(), 18)
	require.Equal(t, 10, accounts.callCount())
	require.Equal(t, 1, accounts.platformCallCount(PlatformOpenAI))
	for _, bucket := range current {
		_, published := cache.counts(bucket)
		require.Equal(t, 1, published, bucket.String())
placeholder
	require.Contains(t, bucketStrings(current), SchedulerBucket{GroupID: groupID, Platform: PlatformAntigravity, Mode: SchedulerModeForcedplaceholder.String())
	require.Contains(t, bucketStrings(current), SchedulerBucket{GroupID: groupID, Platform: PlatformAnthropic, Mode: SchedulerModeMixedplaceholder.String())
	require.Contains(t, bucketStrings(current), SchedulerBucket{GroupID: groupID, Platform: PlatformGemini, Mode: SchedulerModeMixedplaceholder.String())
	acquires, releases, listCalls := cache.lifecycleCounts()
	require.Equal(t, 1, acquires)
	require.Equal(t, 1, releases)
	require.Zero(t, listCalls)
	require.Equal(t, schedulerGroupLifecycleLeaseTTL, cache.acquireTTL)
	require.True(t, cache.acquireDeadline)
	require.True(t, cache.releaseDeadline)
	require.NoError(t, cache.releaseCtxErr)
	_, reopenHeld := cache.lifecycleMutationLeaseStates()
	require.Len(t, reopenHeld, 18)
	for _, held := range reopenHeld {
		require.True(t, held)
placeholder
	lockTTLs, unlockCalls := cache.lockStats()
	require.Len(t, lockTTLs, 18)
	for _, ttl := range lockTTLs {
		require.Equal(t, 30*time.Second, ttl)
placeholder
	require.Equal(t, 18, unlockCalls)
	requireLifecycleSeen(t, seen, groupID)
placeholder

func TestSchedulerGroupLifecycleInactiveThenActiveAuthoritativelyReopens(t *testing.T) {
	const groupID int64 = 83
	cache := newGroupLifecycleTestCache()
	groups := &groupLifecycleTestGroupRepo{group: &Group{ID: groupID, Status: StatusDisabled, Hydrated: trueplaceholderplaceholder
	accounts := &groupLifecycleTestAccountRepo{placeholder
	svc := newGroupLifecycleTestService(cache, accounts, groups, config.RunModeStandard)

	require.NoError(t, svc.handleGroupEvent(context.Background(), ptrInt64(groupID), make(map[batchSeenKey]struct{placeholder)))
	require.Zero(t, accounts.callCount())
	groups.set(&Group{ID: groupID, Status: StatusActive, Hydrated: trueplaceholder, nil)
	require.NoError(t, svc.handleGroupEvent(context.Background(), ptrInt64(groupID), make(map[batchSeenKey]struct{placeholder)))

	require.Len(t, cache.tokens(), 18)
	require.Equal(t, 10, accounts.callCount())
	for _, bucket := range expectedGroupLifecycleBuckets(groupID) {
		_, published := cache.counts(bucket)
		require.Equal(t, 1, published, bucket.String())
placeholder
placeholder

func TestSchedulerGroupLifecycleLaterInactiveFencesLongActiveRebuild(t *testing.T) {
	const groupID int64 = 84
	cache := newGroupLifecycleTestCache()
	groups := &groupLifecycleTestGroupRepo{group: &Group{ID: groupID, Status: StatusActive, Hydrated: trueplaceholderplaceholder
	started := make(chan struct{placeholder)
	release := make(chan struct{placeholder)
	accounts := &groupLifecycleTestAccountRepo{started: started, release: releaseplaceholder
	svc := newGroupLifecycleTestService(cache, accounts, groups, config.RunModeStandard)
	activeSeen := make(map[batchSeenKey]struct{placeholder)
	inactiveSeen := make(map[batchSeenKey]struct{placeholder)
	activeResult := make(chan error, 1)

	go func() {
		activeResult <- svc.handleGroupEvent(context.Background(), ptrInt64(groupID), activeSeen)
placeholder()
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("active rebuild did not reach the account load")
placeholder

	groups.set(&Group{ID: groupID, Status: StatusDisabled, Hydrated: trueplaceholder, nil)
	require.NoError(t, svc.handleGroupEvent(context.Background(), ptrInt64(groupID), inactiveSeen))
	close(release)
	err := <-activeResult
	require.ErrorIs(t, err, ErrSchedulerBucketRetired)
	requireLifecycleNotSeen(t, activeSeen, groupID)
	requireLifecycleSeen(t, inactiveSeen, groupID)
placeholder

func TestSchedulerGroupLifecycleEpochPreventsABA(t *testing.T) {
	const groupID int64 = 85
	cache := newGroupLifecycleTestCache()
	groups := &groupLifecycleTestGroupRepo{group: &Group{ID: groupID, Status: StatusDisabled, Hydrated: trueplaceholderplaceholder
	accounts := &groupLifecycleTestAccountRepo{placeholder
	svc := newGroupLifecycleTestService(cache, accounts, groups, config.RunModeStandard)

	require.NoError(t, svc.handleGroupEvent(context.Background(), ptrInt64(groupID), make(map[batchSeenKey]struct{placeholder)))
	groups.set(&Group{ID: groupID, Status: StatusActive, Hydrated: trueplaceholder, nil)
	require.NoError(t, svc.handleGroupEvent(context.Background(), ptrInt64(groupID), make(map[batchSeenKey]struct{placeholder)))
	firstActiveTokens := cache.tokens()
	require.Len(t, firstActiveTokens, 18)

	groups.set(&Group{ID: groupID, Status: StatusDisabled, Hydrated: trueplaceholder, nil)
	require.NoError(t, svc.handleGroupEvent(context.Background(), ptrInt64(groupID), make(map[batchSeenKey]struct{placeholder)))
	groups.set(&Group{ID: groupID, Status: StatusActive, Hydrated: trueplaceholder, nil)
	require.NoError(t, svc.handleGroupEvent(context.Background(), ptrInt64(groupID), make(map[batchSeenKey]struct{placeholder)))
	allTokens := cache.tokens()
	require.Len(t, allTokens, 36)
	require.Greater(t, allTokens[18].Epoch, firstActiveTokens[0].Epoch)
	require.ErrorIs(t, cache.SetSnapshot(context.Background(), firstActiveTokens[0].Bucket, firstActiveTokens[0], nil), ErrSchedulerBucketWriteFenced)
placeholder

func TestSchedulerGroupLifecycleSeenIsIndependentAndDeduplicatesGroupEvents(t *testing.T) {
	const groupID int64 = 86
	cache := newGroupLifecycleTestCache()
	groups := &groupLifecycleTestGroupRepo{group: &Group{ID: groupID, Status: StatusActive, Hydrated: trueplaceholderplaceholder
	accounts := &groupLifecycleTestAccountRepo{placeholder
	svc := newGroupLifecycleTestService(cache, accounts, groups, config.RunModeStandard)
	seen := make(map[batchSeenKey]struct{placeholder)
	for _, platform := range schedulerSnapshotPlatforms() {
		seen[batchSeenKey{groupID: groupID, platform: platformplaceholder] = struct{placeholder{placeholder
placeholder

	require.NoError(t, svc.handleGroupEvent(context.Background(), ptrInt64(groupID), seen))
	require.Equal(t, 1, groups.callCount())
	require.Equal(t, 10, accounts.callCount())
	requireLifecycleSeen(t, seen, groupID)
	require.NoError(t, svc.handleGroupEvent(context.Background(), ptrInt64(groupID), seen))
	require.Equal(t, 1, groups.callCount())
	require.Equal(t, 10, accounts.callCount())
placeholder

func TestSchedulerGroupLifecycleFailuresDoNotMarkSeen(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(*groupLifecycleTestCache, *groupLifecycleTestGroupRepo, *groupLifecycleTestAccountRepo)
		check   func(*testing.T, error)
placeholder{
		{
			name: "lease busy",
			prepare: func(cache *groupLifecycleTestCache, _ *groupLifecycleTestGroupRepo, _ *groupLifecycleTestAccountRepo) {
				cache.leaseBusy = true
		placeholder,
			check: func(t *testing.T, err error) { require.ErrorIs(t, err, ErrSchedulerGroupLifecycleLeaseBusy) placeholder,
	placeholder,
		{
			name: "lease error",
			prepare: func(cache *groupLifecycleTestCache, _ *groupLifecycleTestGroupRepo, _ *groupLifecycleTestAccountRepo) {
				cache.leaseAcquireErr = errors.New("lease failed")
		placeholder,
			check: func(t *testing.T, err error) { require.EqualError(t, err, "lease failed") placeholder,
	placeholder,
		{
			name: "release lost",
			prepare: func(cache *groupLifecycleTestCache, _ *groupLifecycleTestGroupRepo, _ *groupLifecycleTestAccountRepo) {
				cache.leaseReleaseErr = ErrSchedulerGroupLifecycleLeaseLost
		placeholder,
			check: func(t *testing.T, err error) { require.ErrorIs(t, err, ErrSchedulerGroupLifecycleLeaseLost) placeholder,
	placeholder,
		{
			name: "release error",
			prepare: func(cache *groupLifecycleTestCache, _ *groupLifecycleTestGroupRepo, _ *groupLifecycleTestAccountRepo) {
				cache.leaseReleaseErr = errors.New("release failed")
		placeholder,
			check: func(t *testing.T, err error) { require.EqualError(t, err, "release failed") placeholder,
	placeholder,
		{
			name: "group query error",
			prepare: func(_ *groupLifecycleTestCache, groups *groupLifecycleTestGroupRepo, _ *groupLifecycleTestAccountRepo) {
				groups.err = errors.New("group query failed")
		placeholder,
			check: func(t *testing.T, err error) { require.EqualError(t, err, "group query failed") placeholder,
	placeholder,
		{
			name: "list buckets error",
			prepare: func(cache *groupLifecycleTestCache, groups *groupLifecycleTestGroupRepo, _ *groupLifecycleTestAccountRepo) {
				groups.group.Status = StatusDisabled
				cache.listErr = errors.New("list buckets failed")
		placeholder,
			check: func(t *testing.T, err error) { require.EqualError(t, err, "list buckets failed") placeholder,
	placeholder,
		{
			name: "retire bucket error",
			prepare: func(cache *groupLifecycleTestCache, groups *groupLifecycleTestGroupRepo, _ *groupLifecycleTestAccountRepo) {
				groups.group.Status = StatusDisabled
				cache.retireErr = errors.New("retire bucket failed")
				cache.retireErrAt = 2
		placeholder,
			check: func(t *testing.T, err error) { require.EqualError(t, err, "retire bucket failed") placeholder,
	placeholder,
		{
			name: "reopen bucket error",
			prepare: func(cache *groupLifecycleTestCache, _ *groupLifecycleTestGroupRepo, _ *groupLifecycleTestAccountRepo) {
				cache.reopenErr = errors.New("reopen bucket failed")
				cache.reopenErrAt = 2
		placeholder,
			check: func(t *testing.T, err error) { require.EqualError(t, err, "reopen bucket failed") placeholder,
	placeholder,
		{
			name: "account rebuild error",
			prepare: func(_ *groupLifecycleTestCache, _ *groupLifecycleTestGroupRepo, accounts *groupLifecycleTestAccountRepo) {
				accounts.err = errors.New("account load failed")
		placeholder,
			check: func(t *testing.T, err error) { require.EqualError(t, err, "account load failed") placeholder,
	placeholder,
		{
			name: "bucket lock busy",
			prepare: func(cache *groupLifecycleTestCache, _ *groupLifecycleTestGroupRepo, _ *groupLifecycleTestAccountRepo) {
				cache.bucketLockBusy = true
		placeholder,
			check: func(t *testing.T, err error) { require.ErrorIs(t, err, ErrSchedulerBucketRebuildBusy) placeholder,
	placeholder,
		{
			name: "bucket lock error",
			prepare: func(cache *groupLifecycleTestCache, _ *groupLifecycleTestGroupRepo, _ *groupLifecycleTestAccountRepo) {
				cache.bucketLockErr = errors.New("bucket lock failed")
		placeholder,
			check: func(t *testing.T, err error) { require.EqualError(t, err, "bucket lock failed") placeholder,
	placeholder,
		{
			name: "set snapshot error",
			prepare: func(cache *groupLifecycleTestCache, _ *groupLifecycleTestGroupRepo, _ *groupLifecycleTestAccountRepo) {
				cache.setErr = errors.New("set snapshot failed")
		placeholder,
			check: func(t *testing.T, err error) { require.EqualError(t, err, "set snapshot failed") placeholder,
	placeholder,
placeholder

	for index, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			groupID := int64(870 + index)
			cache := newGroupLifecycleTestCache()
			groups := &groupLifecycleTestGroupRepo{group: &Group{ID: groupID, Status: StatusActive, Hydrated: trueplaceholderplaceholder
			accounts := &groupLifecycleTestAccountRepo{placeholder
			tc.prepare(cache, groups, accounts)
			svc := newGroupLifecycleTestService(cache, accounts, groups, config.RunModeStandard)
			seen := make(map[batchSeenKey]struct{placeholder)

			err := svc.handleGroupEvent(context.Background(), ptrInt64(groupID), seen)
			tc.check(t, err)
			requireLifecycleNotSeen(t, seen, groupID)
			if tc.name == "release lost" || tc.name == "release error" {
				require.Zero(t, accounts.callCount())
		placeholder
			if tc.name == "retire bucket error" || tc.name == "reopen bucket error" {
				_, releases, _ := cache.lifecycleCounts()
				require.Equal(t, 1, releases)
				require.Zero(t, accounts.callCount())
		placeholder
			if tc.name == "account rebuild error" || tc.name == "set snapshot error" {
				lockTTLs, unlockCalls := cache.lockStats()
				require.Len(t, lockTTLs, 1)
				require.Equal(t, 1, unlockCalls)
				require.Equal(t, 1, accounts.callCount())
		placeholder
	placeholder)
placeholder
placeholder

func TestSchedulerGroupLifecycleOperationAndReleaseErrorsPreserveBothCauses(t *testing.T) {
	const groupID int64 = 880
	operationErr := errors.New("group query failed")
	cache := newGroupLifecycleTestCache()
	cache.leaseReleaseErr = ErrSchedulerGroupLifecycleLeaseLost
	groups := &groupLifecycleTestGroupRepo{err: operationErrplaceholder
	accounts := &groupLifecycleTestAccountRepo{placeholder
	svc := newGroupLifecycleTestService(cache, accounts, groups, config.RunModeStandard)
	seen := make(map[batchSeenKey]struct{placeholder)

	err := svc.handleGroupEvent(context.Background(), ptrInt64(groupID), seen)
	require.ErrorIs(t, err, operationErr)
	require.ErrorIs(t, err, ErrSchedulerGroupLifecycleLeaseLost)
	requireLifecycleNotSeen(t, seen, groupID)
	require.Zero(t, accounts.callCount())
placeholder

func TestSchedulerGroupLifecycleUntrustedGroupStateFailsClosed(t *testing.T) {
	for _, tc := range []struct {
		name  string
		group *Group
placeholder{
		{name: "not hydrated", group: &Group{ID: 88, Status: StatusActiveplaceholderplaceholder,
		{name: "mismatched id", group: &Group{ID: 89, Status: StatusActive, Hydrated: trueplaceholderplaceholder,
placeholder {
		t.Run(tc.name, func(t *testing.T) {
			const eventGroupID int64 = 88
			cache := newGroupLifecycleTestCache()
			groups := &groupLifecycleTestGroupRepo{group: tc.groupplaceholder
			accounts := &groupLifecycleTestAccountRepo{placeholder
			svc := newGroupLifecycleTestService(cache, accounts, groups, config.RunModeStandard)
			seen := make(map[batchSeenKey]struct{placeholder)

			err := svc.handleGroupEvent(context.Background(), ptrInt64(eventGroupID), seen)
		placeholder
			require.Empty(t, cache.retiredBuckets())
			require.Empty(t, cache.tokens())
			require.Zero(t, accounts.callCount())
			requireLifecycleNotSeen(t, seen, eventGroupID)
			acquires, releases, listCalls := cache.lifecycleCounts()
			require.Equal(t, 1, acquires)
			require.Equal(t, 1, releases)
			require.Zero(t, listCalls)
	placeholder)
placeholder
placeholder

func TestSchedulerGroupLifecycleCanceledAfterFreshQueryUsesIndependentReleaseContext(t *testing.T) {
	const groupID int64 = 89
	ctx, cancel := context.WithCancel(context.Background())
	cache := newGroupLifecycleTestCache()
	groups := &groupLifecycleTestGroupRepo{
		group:    &Group{ID: groupID, Status: StatusActive, Hydrated: trueplaceholder,
		afterGet: cancel,
placeholder
	accounts := &groupLifecycleTestAccountRepo{placeholder
	svc := newGroupLifecycleTestService(cache, accounts, groups, config.RunModeStandard)
	seen := make(map[batchSeenKey]struct{placeholder)

	err := svc.handleGroupEvent(ctx, ptrInt64(groupID), seen)
	require.ErrorIs(t, err, context.Canceled)
	requireLifecycleNotSeen(t, seen, groupID)
	require.Empty(t, cache.tokens())
	require.Zero(t, accounts.callCount())
	acquires, releases, _ := cache.lifecycleCounts()
	require.Equal(t, 1, acquires)
	require.Equal(t, 1, releases)
	require.True(t, cache.releaseDeadline)
	require.NoError(t, cache.releaseCtxErr)
placeholder

func TestSchedulerGroupLifecycleGroupZeroAndSimpleModeAreNoOps(t *testing.T) {
	cache := newGroupLifecycleTestCache()
	groups := &groupLifecycleTestGroupRepo{group: &Group{ID: 88, Status: StatusActive, Hydrated: trueplaceholderplaceholder
	accounts := &groupLifecycleTestAccountRepo{placeholder
	standard := newGroupLifecycleTestService(cache, accounts, groups, config.RunModeStandard)
	simple := newGroupLifecycleTestService(cache, accounts, groups, config.RunModeSimple)

	require.NoError(t, standard.handleGroupEvent(context.Background(), nil, make(map[batchSeenKey]struct{placeholder)))
	require.NoError(t, standard.handleGroupEvent(context.Background(), ptrInt64(0), make(map[batchSeenKey]struct{placeholder)))
	require.NoError(t, simple.handleGroupEvent(context.Background(), ptrInt64(88), make(map[batchSeenKey]struct{placeholder)))

	acquires, releases, listCalls := cache.lifecycleCounts()
	require.Zero(t, acquires)
	require.Zero(t, releases)
	require.Zero(t, listCalls)
	require.Zero(t, groups.callCount())
	require.Zero(t, accounts.callCount())
placeholder
