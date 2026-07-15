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

type batchAccountQueryKey struct {
	groupID  int64
	platform string
	mixed    bool
placeholder

type batchAccountQueryResult struct {
	accounts []Account
	err      error
placeholder

type batchAccountQueryRepo struct {
	AccountRepository

	mu        sync.Mutex
	calls     map[batchAccountQueryKey]int
	results   map[batchAccountQueryKey][]batchAccountQueryResult
	beforeRun func(batchAccountQueryKey)
placeholder

func newBatchAccountQueryRepo() *batchAccountQueryRepo {
	return &batchAccountQueryRepo{
		calls:   make(map[batchAccountQueryKey]int),
		results: make(map[batchAccountQueryKey][]batchAccountQueryResult),
placeholder
placeholder

func (r *batchAccountQueryRepo) ListSchedulableByGroupIDAndPlatform(_ context.Context, groupID int64, platform string) ([]Account, error) {
	return r.run(batchAccountQueryKey{groupID: groupID, platform: platformplaceholder)
placeholder

func (r *batchAccountQueryRepo) ListSchedulableByGroupIDAndPlatforms(_ context.Context, groupID int64, platforms []string) ([]Account, error) {
	return r.run(batchAccountQueryKey{groupID: groupID, platform: platforms[0], mixed: trueplaceholder)
placeholder

func (r *batchAccountQueryRepo) ListSchedulableUngroupedByPlatform(_ context.Context, platform string) ([]Account, error) {
	return r.run(batchAccountQueryKey{platform: platformplaceholder)
placeholder

func (r *batchAccountQueryRepo) ListSchedulableUngroupedByPlatforms(_ context.Context, platforms []string) ([]Account, error) {
	return r.run(batchAccountQueryKey{platform: platforms[0], mixed: trueplaceholder)
placeholder

func (r *batchAccountQueryRepo) ListSchedulableByPlatform(_ context.Context, platform string) ([]Account, error) {
	return r.run(batchAccountQueryKey{platform: platformplaceholder)
placeholder

func (r *batchAccountQueryRepo) ListSchedulableByPlatforms(_ context.Context, platforms []string) ([]Account, error) {
	return r.run(batchAccountQueryKey{platform: platforms[0], mixed: trueplaceholder)
placeholder

func (r *batchAccountQueryRepo) run(key batchAccountQueryKey) ([]Account, error) {
	r.mu.Lock()
	r.calls[key]++
	call := r.calls[key]
	results := r.results[key]
	beforeRun := r.beforeRun
	r.mu.Unlock()

	if beforeRun != nil {
		beforeRun(key)
placeholder
	if call <= len(results) {
		result := results[call-1]
		return append([]Account(nil), result.accounts...), result.err
placeholder
	return []Account{{
		ID:          int64(call),
		Name:        "source",
		Platform:    key.platform,
		Status:      StatusActive,
		Schedulable: true,
placeholderplaceholder, nil
placeholder

func (r *batchAccountQueryRepo) callCount(key batchAccountQueryKey) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.calls[key]
placeholder

type batchSnapshotWrite struct {
	token    SchedulerBucketWriteToken
	accounts []Account
placeholder

type batchSnapshotCache struct {
	SchedulerCache

	mu          sync.Mutex
	nextEpoch   int64
	captures    []SchedulerBucket
	captured    map[SchedulerBucket]SchedulerBucketWriteToken
	locks       map[SchedulerBucket]int
	lockBusy    map[SchedulerBucket]bool
	lockErrors  map[SchedulerBucket]error
	setErrors   map[SchedulerBucket]error
	setAttempts map[SchedulerBucket]int
	writes      map[SchedulerBucket][]batchSnapshotWrite
	versions    map[SchedulerBucket]int
	beforeSet   func()
placeholder

func newBatchSnapshotCache() *batchSnapshotCache {
	return &batchSnapshotCache{
		captured:    make(map[SchedulerBucket]SchedulerBucketWriteToken),
		locks:       make(map[SchedulerBucket]int),
		lockBusy:    make(map[SchedulerBucket]bool),
		lockErrors:  make(map[SchedulerBucket]error),
		setErrors:   make(map[SchedulerBucket]error),
		setAttempts: make(map[SchedulerBucket]int),
		writes:      make(map[SchedulerBucket][]batchSnapshotWrite),
		versions:    make(map[SchedulerBucket]int),
placeholder
placeholder

func (c *batchSnapshotCache) CaptureBucketWriteToken(_ context.Context, bucket SchedulerBucket) (SchedulerBucketWriteToken, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nextEpoch++
	token := SchedulerBucketWriteToken{Bucket: bucket, Epoch: c.nextEpochplaceholder
	c.captures = append(c.captures, bucket)
	c.captured[bucket] = token
	return token, nil
placeholder

func (c *batchSnapshotCache) TryLockBucket(_ context.Context, bucket SchedulerBucket, _ time.Duration) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.locks[bucket]++
	if err := c.lockErrors[bucket]; err != nil {
		return false, err
placeholder
	return !c.lockBusy[bucket], nil
placeholder

func (c *batchSnapshotCache) UnlockBucket(context.Context, SchedulerBucket) error {
	return nil
placeholder

func (c *batchSnapshotCache) SetSnapshot(_ context.Context, bucket SchedulerBucket, token SchedulerBucketWriteToken, accounts []Account) error {
	if c.beforeSet != nil {
		c.beforeSet()
placeholder
	c.mu.Lock()
	defer c.mu.Unlock()
	c.setAttempts[bucket]++
	if token != c.captured[bucket] || !token.ValidFor(bucket) {
		return ErrSchedulerBucketWriteFenced
placeholder
	if err := c.setErrors[bucket]; err != nil {
		return err
placeholder
	c.versions[bucket]++
	c.writes[bucket] = append(c.writes[bucket], batchSnapshotWrite{
		token:    token,
		accounts: append([]Account(nil), accounts...),
placeholder)
	return nil
placeholder

func (c *batchSnapshotCache) captureCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.captures)
placeholder

func (c *batchSnapshotCache) bucketState(bucket SchedulerBucket) (locks, attempts, version int, writes []batchSnapshotWrite) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.locks[bucket], c.setAttempts[bucket], c.versions[bucket], append([]batchSnapshotWrite(nil), c.writes[bucket]...)
placeholder

func newBatchQueryTestService(cache SchedulerCache, accounts AccountRepository, runMode string) *SchedulerSnapshotService {
	return NewSchedulerSnapshotService(cache, nil, accounts, nil, &config.Config{RunMode: runModeplaceholder)
placeholder

func TestSchedulerRebuildBatchReusesSingleForcedQueryAndKeepsSnapshotsIndependent(t *testing.T) {
	const groupID int64 = 201
	single := SchedulerBucket{GroupID: groupID, Platform: PlatformOpenAI, Mode: SchedulerModeSingleplaceholder
	forced := SchedulerBucket{GroupID: groupID, Platform: PlatformOpenAI, Mode: SchedulerModeForcedplaceholder
	cache := newBatchSnapshotCache()
	repo := newBatchAccountQueryRepo()
	wantCaptures := 2
	repo.beforeRun = func(batchAccountQueryKey) {
		require.Equal(t, wantCaptures, cache.captureCount(), "all tokens must be prepared before the first DB query")
		wantCaptures += 2
placeholder
	svc := newBatchQueryTestService(cache, repo, config.RunModeStandard)

	require.NoError(t, svc.rebuildBuckets(context.Background(), []SchedulerBucket{single, forcedplaceholder, "first"))
	queryKey := batchAccountQueryKey{groupID: groupID, platform: PlatformOpenAIplaceholder
	require.Equal(t, 1, repo.callCount(queryKey))
	for _, bucket := range []SchedulerBucket{single, forcedplaceholder {
		locks, attempts, version, writes := cache.bucketState(bucket)
		require.Equal(t, 1, locks, bucket.String())
		require.Equal(t, 1, attempts, bucket.String())
		require.Equal(t, 1, version, bucket.String())
		require.Len(t, writes, 1, bucket.String())
		require.Equal(t, "source", writes[0].accounts[0].Name, bucket.String())
		require.Equal(t, bucket, writes[0].token.Bucket)
placeholder
	_, _, _, singleWrites := cache.bucketState(single)
	_, _, _, forcedWrites := cache.bucketState(forced)
	require.NotEqual(t, singleWrites[0].token.Epoch, forcedWrites[0].token.Epoch)

	require.NoError(t, svc.rebuildBuckets(context.Background(), []SchedulerBucket{single, forcedplaceholder, "second"))
	require.Equal(t, 2, repo.callCount(queryKey), "successful results must not be cached across rebuild batches")
	for _, bucket := range []SchedulerBucket{single, forcedplaceholder {
		locks, attempts, version, writes := cache.bucketState(bucket)
		require.Equal(t, 2, locks, bucket.String())
		require.Equal(t, 2, attempts, bucket.String())
		require.Equal(t, 2, version, bucket.String())
		require.Len(t, writes, 2, bucket.String())
		require.Equal(t, "source", writes[1].accounts[0].Name, bucket.String())
placeholder
placeholder

func TestSchedulerRebuildBatchKeepsMixedAndDifferentKeysIndependent(t *testing.T) {
	const groupID int64 = 202
	buckets := []SchedulerBucket{
		{GroupID: groupID, Platform: PlatformAnthropic, Mode: SchedulerModeSingleplaceholder,
		{GroupID: groupID, Platform: PlatformAnthropic, Mode: SchedulerModeForcedplaceholder,
		{GroupID: groupID, Platform: PlatformAnthropic, Mode: SchedulerModeMixedplaceholder,
		{GroupID: groupID + 1, Platform: PlatformAnthropic, Mode: SchedulerModeSingleplaceholder,
		{GroupID: groupID, Platform: PlatformGemini, Mode: SchedulerModeForcedplaceholder,
		{GroupID: 0, Platform: PlatformOpenAI, Mode: SchedulerModeSingleplaceholder,
		{GroupID: -1, Platform: PlatformOpenAI, Mode: SchedulerModeForcedplaceholder,
placeholder
	cache := newBatchSnapshotCache()
	repo := newBatchAccountQueryRepo()
	svc := newBatchQueryTestService(cache, repo, config.RunModeStandard)

	require.NoError(t, svc.rebuildBuckets(context.Background(), buckets, "test"))
	require.Equal(t, 1, repo.callCount(batchAccountQueryKey{groupID: groupID, platform: PlatformAnthropicplaceholder))
	require.Equal(t, 1, repo.callCount(batchAccountQueryKey{groupID: groupID, platform: PlatformAnthropic, mixed: trueplaceholder))
	require.Equal(t, 1, repo.callCount(batchAccountQueryKey{groupID: groupID + 1, platform: PlatformAnthropicplaceholder))
	require.Equal(t, 1, repo.callCount(batchAccountQueryKey{groupID: groupID, platform: PlatformGeminiplaceholder))
	require.Equal(t, 2, repo.callCount(batchAccountQueryKey{platform: PlatformOpenAIplaceholder), "group0 and a negative historical group must not share")
	for _, bucket := range buckets {
		locks, attempts, version, _ := cache.bucketState(bucket)
		require.Equal(t, 1, locks, bucket.String())
		require.Equal(t, 1, attempts, bucket.String())
		require.Equal(t, 1, version, bucket.String())
placeholder
placeholder

func TestSchedulerRebuildBatchKeepsSimpleModeBucketGroupsIndependent(t *testing.T) {
	single := SchedulerBucket{GroupID: 204, Platform: PlatformOpenAI, Mode: SchedulerModeSingleplaceholder
	forced := SchedulerBucket{GroupID: 0, Platform: PlatformOpenAI, Mode: SchedulerModeForcedplaceholder
	cache := newBatchSnapshotCache()
	repo := newBatchAccountQueryRepo()
	svc := newBatchQueryTestService(cache, repo, config.RunModeSimple)

	require.NoError(t, svc.rebuildBuckets(context.Background(), []SchedulerBucket{single, forcedplaceholder, "test"))
	require.Equal(t, 2, repo.callCount(batchAccountQueryKey{platform: PlatformOpenAIplaceholder))
placeholder

func TestSchedulerRebuildBatchDoesNotCacheMixedOrHistoricalQueries(t *testing.T) {
	for _, tc := range []struct {
		name   string
		bucket SchedulerBucket
		key    batchAccountQueryKey
placeholder{
		{
			name:   "mixed",
			bucket: SchedulerBucket{GroupID: 204, Platform: PlatformAnthropic, Mode: SchedulerModeMixedplaceholder,
			key:    batchAccountQueryKey{groupID: 204, platform: PlatformAnthropic, mixed: trueplaceholder,
	placeholder,
		{
			name:   "historical",
			bucket: SchedulerBucket{GroupID: 204, Platform: PlatformOpenAI, Mode: "unknown"placeholder,
			key:    batchAccountQueryKey{groupID: 204, platform: PlatformOpenAIplaceholder,
	placeholder,
placeholder {
		t.Run(tc.name, func(t *testing.T) {
			cache := newBatchSnapshotCache()
			token, err := cache.CaptureBucketWriteToken(context.Background(), tc.bucket)
		placeholder
			repo := newBatchAccountQueryRepo()
			svc := newBatchQueryTestService(cache, repo, config.RunModeStandard)
			tasks := []schedulerBucketWriteTask{
				{bucket: tc.bucket, token: tokenplaceholder,
				{bucket: tc.bucket, token: tokenplaceholder,
		placeholder
			queries := newSchedulerAccountQueryCache(tasks)

			require.NoError(t, svc.rebuildPreparedBucketTasks(context.Background(), tasks, "test", false, queries))
			require.Equal(t, 2, repo.callCount(tc.key))
			require.Empty(t, queries.accounts)
			locks, attempts, version, _ := cache.bucketState(tc.bucket)
			require.Equal(t, 2, locks)
			require.Equal(t, 2, attempts)
			require.Equal(t, 2, version)
	placeholder)
placeholder
placeholder

func TestSchedulerRebuildBatchRetriesQueryFailureForFollowingBucket(t *testing.T) {
	const groupID int64 = 205
	single := SchedulerBucket{GroupID: groupID, Platform: PlatformOpenAI, Mode: SchedulerModeSingleplaceholder
	forced := SchedulerBucket{GroupID: groupID, Platform: PlatformOpenAI, Mode: SchedulerModeForcedplaceholder
	wantErr := errors.New("first query failed")
	key := batchAccountQueryKey{groupID: groupID, platform: PlatformOpenAIplaceholder
	repo := newBatchAccountQueryRepo()
	repo.results[key] = []batchAccountQueryResult{
		{err: wantErrplaceholder,
		{accounts: []Account{{ID: 2051, Name: "retry", Platform: PlatformOpenAIplaceholderplaceholderplaceholder,
placeholder
	cache := newBatchSnapshotCache()
	svc := newBatchQueryTestService(cache, repo, config.RunModeStandard)

	err := svc.rebuildBuckets(context.Background(), []SchedulerBucket{single, forcedplaceholder, "test")
	require.ErrorIs(t, err, wantErr)
	require.Equal(t, 2, repo.callCount(key), "failed queries must not enter the batch cache")
	_, singleAttempts, singleVersion, _ := cache.bucketState(single)
	_, forcedAttempts, forcedVersion, forcedWrites := cache.bucketState(forced)
	require.Zero(t, singleAttempts)
	require.Zero(t, singleVersion)
	require.Equal(t, 1, forcedAttempts)
	require.Equal(t, 1, forcedVersion)
	require.Equal(t, "retry", forcedWrites[0].accounts[0].Name)
placeholder

func TestSchedulerFullRebuildSharesSuccessfulQueryAcrossStrictAndOrdinarySegments(t *testing.T) {
	const groupID int64 = 206
	single := SchedulerBucket{GroupID: groupID, Platform: PlatformOpenAI, Mode: SchedulerModeSingleplaceholder
	forced := SchedulerBucket{GroupID: groupID, Platform: PlatformOpenAI, Mode: SchedulerModeForcedplaceholder
	cache := newBatchSnapshotCache()
	cache.setErrors[single] = ErrSchedulerBucketWriteFenced
	singleToken, err := cache.CaptureBucketWriteToken(context.Background(), single)
placeholder
	forcedToken, err := cache.CaptureBucketWriteToken(context.Background(), forced)
placeholder
	repo := newBatchAccountQueryRepo()
	svc := newBatchQueryTestService(cache, repo, config.RunModeStandard)

	err = svc.prepareAndRebuildFullSnapshot(
		context.Background(),
		[]schedulerBucketWriteTask{{bucket: forced, token: forcedTokenplaceholderplaceholder,
		[]schedulerBucketWriteTask{{bucket: single, token: singleTokenplaceholderplaceholder,
		nil,
		"test",
	)
	require.ErrorIs(t, err, ErrSchedulerBucketWriteFenced)
	require.Equal(t, 1, repo.callCount(batchAccountQueryKey{groupID: groupID, platform: PlatformOpenAIplaceholder), "SetSnapshot failure must not discard a successful query")
	_, singleAttempts, singleVersion, _ := cache.bucketState(single)
	_, forcedAttempts, forcedVersion, _ := cache.bucketState(forced)
	require.Equal(t, 1, singleAttempts)
	require.Zero(t, singleVersion)
	require.Equal(t, 1, forcedAttempts)
	require.Equal(t, 1, forcedVersion)
placeholder

func TestSchedulerRebuildBatchPreservesLockBusyAndFencingPolicy(t *testing.T) {
	const groupID int64 = 207
	single := SchedulerBucket{GroupID: groupID, Platform: PlatformOpenAI, Mode: SchedulerModeSingleplaceholder
	forced := SchedulerBucket{GroupID: groupID, Platform: PlatformOpenAI, Mode: SchedulerModeForcedplaceholder

	t.Run("ordinary lock busy skips only that bucket", func(t *testing.T) {
		cache := newBatchSnapshotCache()
		cache.lockBusy[single] = true
		repo := newBatchAccountQueryRepo()
		svc := newBatchQueryTestService(cache, repo, config.RunModeStandard)

		require.NoError(t, svc.rebuildBuckets(context.Background(), []SchedulerBucket{single, forcedplaceholder, "test"))
		require.Equal(t, 1, repo.callCount(batchAccountQueryKey{groupID: groupID, platform: PlatformOpenAIplaceholder))
		_, singleAttempts, _, _ := cache.bucketState(single)
		_, forcedAttempts, forcedVersion, _ := cache.bucketState(forced)
		require.Zero(t, singleAttempts)
		require.Equal(t, 1, forcedAttempts)
		require.Equal(t, 1, forcedVersion)
placeholder)

	t.Run("strict lock busy is returned while ordinary work continues", func(t *testing.T) {
		cache := newBatchSnapshotCache()
		cache.lockBusy[single] = true
		singleToken, err := cache.CaptureBucketWriteToken(context.Background(), single)
	placeholder
		forcedToken, err := cache.CaptureBucketWriteToken(context.Background(), forced)
	placeholder
		repo := newBatchAccountQueryRepo()
		svc := newBatchQueryTestService(cache, repo, config.RunModeStandard)

		err = svc.prepareAndRebuildFullSnapshot(
			context.Background(),
			[]schedulerBucketWriteTask{{bucket: forced, token: forcedTokenplaceholderplaceholder,
			[]schedulerBucketWriteTask{{bucket: single, token: singleTokenplaceholderplaceholder,
			nil,
			"test",
		)
		require.ErrorIs(t, err, ErrSchedulerBucketRebuildBusy)
		require.Equal(t, 1, repo.callCount(batchAccountQueryKey{groupID: groupID, platform: PlatformOpenAIplaceholder))
		_, forcedAttempts, forcedVersion, _ := cache.bucketState(forced)
		require.Equal(t, 1, forcedAttempts)
		require.Equal(t, 1, forcedVersion)
placeholder)

	t.Run("ordinary fencing stays non-fatal", func(t *testing.T) {
		cache := newBatchSnapshotCache()
		cache.setErrors[single] = ErrSchedulerBucketWriteFenced
		repo := newBatchAccountQueryRepo()
		svc := newBatchQueryTestService(cache, repo, config.RunModeStandard)

		require.NoError(t, svc.rebuildBuckets(context.Background(), []SchedulerBucket{single, forcedplaceholder, "test"))
		require.Equal(t, 1, repo.callCount(batchAccountQueryKey{groupID: groupID, platform: PlatformOpenAIplaceholder))
		_, singleAttempts, singleVersion, _ := cache.bucketState(single)
		_, forcedAttempts, forcedVersion, _ := cache.bucketState(forced)
		require.Equal(t, 1, singleAttempts)
		require.Zero(t, singleVersion)
		require.Equal(t, 1, forcedAttempts)
		require.Equal(t, 1, forcedVersion)
placeholder)
placeholder

func TestSchedulerRebuildBatchReleasesResultsAfterLastConsumer(t *testing.T) {
	const groups = 128
	cache := newBatchSnapshotCache()
	repo := newBatchAccountQueryRepo()
	tasks := make([]schedulerBucketWriteTask, 0, groups*2)
	wantLockErr := errors.New("lock failed")
	for i := 1; i <= groups; i++ {
		groupID := int64(300 + i)
		single := SchedulerBucket{GroupID: groupID, Platform: PlatformOpenAI, Mode: SchedulerModeSingleplaceholder
		forced := SchedulerBucket{GroupID: groupID, Platform: PlatformOpenAI, Mode: SchedulerModeForcedplaceholder
		if i == 1 {
			cache.lockBusy[single] = true
	placeholder
		if i == 2 {
			cache.lockErrors[single] = wantLockErr
	placeholder
		for _, bucket := range []SchedulerBucket{single, forcedplaceholder {
			token, err := cache.CaptureBucketWriteToken(context.Background(), bucket)
		placeholder
			tasks = append(tasks, schedulerBucketWriteTask{bucket: bucket, token: tokenplaceholder)
	placeholder
placeholder
	queries := newSchedulerAccountQueryCache(tasks)
	maxResident := 0
	cache.beforeSet = func() {
		if resident := len(queries.accounts); resident > maxResident {
			maxResident = resident
	placeholder
placeholder
	svc := newBatchQueryTestService(cache, repo, config.RunModeStandard)

	err := svc.rebuildPreparedBucketTasks(context.Background(), tasks, "test", false, queries)
	require.ErrorIs(t, err, wantLockErr)
	require.LessOrEqual(t, maxResident, 1, "adjacent single/forced pairs must not accumulate full-batch results")
	require.Empty(t, queries.accounts)
	require.Empty(t, queries.remaining)
	for i := 1; i <= groups; i++ {
		key := batchAccountQueryKey{groupID: int64(300 + i), platform: PlatformOpenAIplaceholder
		require.Equal(t, 1, repo.callCount(key), key)
placeholder
placeholder

type batchQueryBenchmarkRepo struct {
	AccountRepository
	accounts []Account
placeholder

func (r *batchQueryBenchmarkRepo) ListSchedulableByGroupIDAndPlatform(context.Context, int64, string) ([]Account, error) {
	return r.accounts, nil
placeholder

type batchQueryBenchmarkCache struct {
	SchedulerCache
placeholder

func (c *batchQueryBenchmarkCache) CaptureBucketWriteToken(_ context.Context, bucket SchedulerBucket) (SchedulerBucketWriteToken, error) {
	return SchedulerBucketWriteToken{Bucket: bucket, Epoch: 1placeholder, nil
placeholder

func (c *batchQueryBenchmarkCache) TryLockBucket(context.Context, SchedulerBucket, time.Duration) (bool, error) {
	return true, nil
placeholder

func (c *batchQueryBenchmarkCache) UnlockBucket(context.Context, SchedulerBucket) error {
	return nil
placeholder

var batchQueryBenchmarkAccountCount int

func (c *batchQueryBenchmarkCache) SetSnapshot(_ context.Context, _ SchedulerBucket, _ SchedulerBucketWriteToken, accounts []Account) error {
	batchQueryBenchmarkAccountCount = len(accounts)
	return nil
placeholder

func BenchmarkSchedulerRebuildBatchQueryReuse(b *testing.B) {
	const groupID int64 = 208
	buckets := []SchedulerBucket{
		{GroupID: groupID, Platform: PlatformOpenAI, Mode: SchedulerModeSingleplaceholder,
		{GroupID: groupID, Platform: PlatformOpenAI, Mode: SchedulerModeForcedplaceholder,
placeholder
	for _, tc := range []struct {
		name string
		size int
placeholder{
		{name: "1_account", size: 1placeholder,
		{name: "10000_accounts", size: 10_000placeholder,
placeholder {
		b.Run(tc.name, func(b *testing.B) {
			accounts := make([]Account, tc.size)
			svc := newBatchQueryTestService(
				&batchQueryBenchmarkCache{placeholder,
				&batchQueryBenchmarkRepo{accounts: accountsplaceholder,
				config.RunModeStandard,
			)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := svc.rebuildBuckets(context.Background(), buckets, "benchmark"); err != nil {
					b.Fatal(err)
			placeholder
		placeholder
	placeholder)
placeholder
placeholder
