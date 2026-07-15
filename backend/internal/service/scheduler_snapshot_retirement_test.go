//go:build unit

package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type retirementRaceCache struct {
	SchedulerCache

	mu          sync.Mutex
	epochs      map[string]int64
	retired     map[string]bool
	listBuckets []SchedulerBucket
	captures    []SchedulerBucket
	reopens     []SchedulerBucket
	setAttempts map[string]int
	published   map[string]int
	versions    map[string]int
	beforeSet   func()
placeholder

func newRetirementRaceCache(buckets ...SchedulerBucket) *retirementRaceCache {
	return &retirementRaceCache{
		epochs:      make(map[string]int64),
		retired:     make(map[string]bool),
		listBuckets: buckets,
		setAttempts: make(map[string]int),
		published:   make(map[string]int),
		versions:    make(map[string]int),
placeholder
placeholder

func (c *retirementRaceCache) GetSnapshot(context.Context, SchedulerBucket) ([]*Account, bool, error) {
	return nil, false, nil
placeholder

func (c *retirementRaceCache) CaptureBucketWriteToken(_ context.Context, bucket SchedulerBucket) (SchedulerBucketWriteToken, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := bucket.String()
	c.captures = append(c.captures, bucket)
	if c.retired[key] {
		return SchedulerBucketWriteToken{placeholder, ErrSchedulerBucketRetired
placeholder
	if c.epochs[key] == 0 {
		c.epochs[key] = 1
placeholder
	return SchedulerBucketWriteToken{Bucket: bucket, Epoch: c.epochs[key]placeholder, nil
placeholder

func (c *retirementRaceCache) SetSnapshot(_ context.Context, bucket SchedulerBucket, token SchedulerBucketWriteToken, _ []Account) error {
	if c.beforeSet != nil {
		c.beforeSet()
placeholder
	c.mu.Lock()
	defer c.mu.Unlock()
	key := bucket.String()
	c.setAttempts[key]++
	if !token.ValidFor(bucket) {
		return ErrSchedulerBucketWriteFenced
placeholder
	if c.retired[key] {
		return ErrSchedulerBucketRetired
placeholder
	if c.epochs[key] != token.Epoch {
		return ErrSchedulerBucketWriteFenced
placeholder
	c.versions[key]++
	c.published[key]++
	return nil
placeholder

func (c *retirementRaceCache) RetireBucket(_ context.Context, bucket SchedulerBucket) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := bucket.String()
	if !c.retired[key] {
		c.epochs[key]++
		if c.epochs[key] < 1 {
			c.epochs[key] = 1
	placeholder
		c.retired[key] = true
placeholder
	return nil
placeholder

func (c *retirementRaceCache) ReopenBucket(_ context.Context, bucket SchedulerBucket) (SchedulerBucketWriteToken, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := bucket.String()
	if c.epochs[key] == 0 {
		c.epochs[key] = 1
placeholder
	delete(c.retired, key)
	c.reopens = append(c.reopens, bucket)
	return SchedulerBucketWriteToken{Bucket: bucket, Epoch: c.epochs[key]placeholder, nil
placeholder

func (c *retirementRaceCache) TryLockBucket(context.Context, SchedulerBucket, time.Duration) (bool, error) {
	return true, nil
placeholder

func (c *retirementRaceCache) UnlockBucket(context.Context, SchedulerBucket) error {
	return nil
placeholder

func (c *retirementRaceCache) ListBuckets(context.Context) ([]SchedulerBucket, error) {
	return append([]SchedulerBucket(nil), c.listBuckets...), nil
placeholder

func (c *retirementRaceCache) counts(bucket SchedulerBucket) (setAttempts, published int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.setAttempts[bucket.String()], c.published[bucket.String()]
placeholder

func (c *retirementRaceCache) captureAndReopenCounts() (captures, reopens int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.captures), len(c.reopens)
placeholder

func (c *retirementRaceCache) version(bucket SchedulerBucket) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.versions[bucket.String()]
placeholder

type retirementGroupRepo struct {
	GroupRepository
	groups []Group
	err    error
placeholder

func (r *retirementGroupRepo) ListActive(context.Context) ([]Group, error) {
	return r.groups, r.err
placeholder

func TestSchedulerFullRebuildCapturesAllRegistryTokensBeforeDBLoad(t *testing.T) {
	first := SchedulerBucket{GroupID: 61, Platform: PlatformOpenAI, Mode: SchedulerModeSingleplaceholder
	queued := SchedulerBucket{GroupID: 61, Platform: PlatformOpenAI, Mode: SchedulerModeForcedplaceholder
	cache := newRetirementRaceCache(first, queued)
	dbStarted := make(chan struct{placeholder)
	releaseDB := make(chan struct{placeholder)
	var firstDB sync.Once
	repo := &mockAccountRepoForPlatform{
		listPlatformFunc: func(context.Context, string) ([]Account, error) {
			firstDB.Do(func() {
				close(dbStarted)
				<-releaseDB
		placeholder)
			return []Account{{ID: 6101, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: trueplaceholderplaceholder, nil
	placeholder,
placeholder
	svc := NewSchedulerSnapshotService(cache, nil, repo, &retirementGroupRepo{groups: []Group{{ID: 61, Status: StatusActiveplaceholderplaceholderplaceholder, &config.Config{
		RunMode: config.RunModeStandard,
		Gateway: config.GatewayConfig{Scheduling: config.GatewaySchedulingConfig{
			DbFallbackEnabled: true,
placeholder
placeholder)

	result := make(chan error, 1)
	go func() { result <- svc.triggerFullRebuild("retirement_race_a") placeholder()
	select {
	case <-dbStarted:
	case <-time.After(time.Second):
		t.Fatal("first DB load did not start")
placeholder

	captures, reopens := cache.captureAndReopenCounts()
	require.Equal(t, 24, captures, "group0 and active-group canonical tokens must be captured before the first DB load")
	require.Zero(t, reopens)
	require.NoError(t, cache.RetireBucket(context.Background(), queued))
	_, err := cache.ReopenBucket(context.Background(), queued)
placeholder
	close(releaseDB)
	require.NoError(t, <-result)

	_, firstPublished := cache.counts(first)
	queuedAttempts, queuedPublished := cache.counts(queued)
	require.Equal(t, 1, firstPublished)
	require.Equal(t, 1, queuedAttempts)
	require.Zero(t, queuedPublished, "queued registry task must not adopt the reopened epoch")
placeholder

func TestSchedulerRebuildRetireAfterDBLoadFencesPublish(t *testing.T) {
	bucket := SchedulerBucket{GroupID: 62, Platform: PlatformOpenAI, Mode: SchedulerModeSingleplaceholder
	cache := newRetirementRaceCache()
	dbReturned := make(chan struct{placeholder)
	setEntered := make(chan struct{placeholder)
	releaseSet := make(chan struct{placeholder)
	cache.beforeSet = func() {
		close(setEntered)
		<-releaseSet
placeholder
	repo := &mockAccountRepoForPlatform{
		listPlatformFunc: func(context.Context, string) ([]Account, error) {
			close(dbReturned)
			return []Account{{ID: 6201, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: trueplaceholderplaceholder, nil
	placeholder,
placeholder
	svc := NewSchedulerSnapshotService(cache, nil, repo, nil, &config.Config{
		RunMode: config.RunModeStandard,
		Gateway: config.GatewayConfig{Scheduling: config.GatewaySchedulingConfig{
			DbFallbackEnabled: true,
placeholder
placeholder)

	result := make(chan error, 1)
	go func() {
		result <- svc.rebuildBuckets(context.Background(), []SchedulerBucket{bucketplaceholder, "retirement_race_b")
placeholder()
	select {
	case <-dbReturned:
	case <-time.After(time.Second):
		t.Fatal("DB load did not return")
placeholder
	select {
	case <-setEntered:
	case <-time.After(time.Second):
		t.Fatal("snapshot writer did not reach allocation boundary")
placeholder
	require.NoError(t, cache.RetireBucket(context.Background(), bucket))
	close(releaseSet)
	require.NoError(t, <-result)

	setAttempts, published := cache.counts(bucket)
	require.Equal(t, 1, setAttempts)
	require.Zero(t, published)
	require.Zero(t, cache.version(bucket), "retirement before allocation must not advance the snapshot version")
placeholder

func TestSchedulerFallbackReturnsDBAccountsWhenBucketRetired(t *testing.T) {
	bucket := SchedulerBucket{GroupID: 63, Platform: PlatformOpenAI, Mode: SchedulerModeSingleplaceholder
	cache := newRetirementRaceCache()
	require.NoError(t, cache.RetireBucket(context.Background(), bucket))
	repo := &mockAccountRepoForPlatform{
		accounts: []Account{{ID: 6301, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: trueplaceholderplaceholder,
placeholder
	svc := NewSchedulerSnapshotService(cache, nil, repo, nil, &config.Config{
		RunMode: config.RunModeStandard,
		Gateway: config.GatewayConfig{Scheduling: config.GatewaySchedulingConfig{
			DbFallbackEnabled: true,
placeholder
placeholder)
	groupID := bucket.GroupID

	accounts, useMixed, err := svc.ListSchedulableAccounts(context.Background(), &groupID, bucket.Platform, false)
placeholder
	require.False(t, useMixed)
	require.Len(t, accounts, 1)
	setAttempts, published := cache.counts(bucket)
	require.Zero(t, setAttempts)
	require.Zero(t, published)
placeholder
