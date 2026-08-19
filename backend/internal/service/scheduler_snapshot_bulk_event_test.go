//go:build unit

package service

import (
	"context"
	"sync"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type bulkEventAccountRepo struct {
	*batchAccountQueryRepo
	accounts []*Account
placeholder

func newBulkEventAccountRepo(accounts ...*Account) *bulkEventAccountRepo {
	return &bulkEventAccountRepo{
		batchAccountQueryRepo: newBatchAccountQueryRepo(),
		accounts:              accounts,
placeholder
placeholder

func (r *bulkEventAccountRepo) GetByIDs(context.Context, []int64) ([]*Account, error) {
	return append([]*Account(nil), r.accounts...), nil
placeholder

type bulkEventSnapshotCache struct {
	*batchSnapshotCache

	accountMu        sync.Mutex
	setAccountIDs    []int64
	deleteAccountIDs []int64
placeholder

func newBulkEventSnapshotCache() *bulkEventSnapshotCache {
	return &bulkEventSnapshotCache{batchSnapshotCache: newBatchSnapshotCache()placeholder
placeholder

func (c *bulkEventSnapshotCache) SetAccount(_ context.Context, account *Account) error {
	c.accountMu.Lock()
	defer c.accountMu.Unlock()
	c.setAccountIDs = append(c.setAccountIDs, account.ID)
	return nil
placeholder

func (c *bulkEventSnapshotCache) DeleteAccount(_ context.Context, accountID int64) error {
	c.accountMu.Lock()
	defer c.accountMu.Unlock()
	c.deleteAccountIDs = append(c.deleteAccountIDs, accountID)
	return nil
placeholder

func (c *bulkEventSnapshotCache) accountWrites() (set []int64, deleted []int64) {
	c.accountMu.Lock()
	defer c.accountMu.Unlock()
	return append([]int64(nil), c.setAccountIDs...), append([]int64(nil), c.deleteAccountIDs...)
placeholder

func (c *bulkEventSnapshotCache) capturedBuckets() []SchedulerBucket {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]SchedulerBucket(nil), c.captures...)
placeholder

func newBulkEventTestService(cache SchedulerCache, accounts AccountRepository) *SchedulerSnapshotService {
	return NewSchedulerSnapshotService(cache, nil, accounts, nil, &config.Config{RunMode: config.RunModeStandardplaceholder)
placeholder

func bulkEventPayload(accountIDs []int64, groupIDs []int64) map[string]any {
	accountValues := make([]any, 0, len(accountIDs))
	for _, id := range accountIDs {
		accountValues = append(accountValues, id)
placeholder
	groupValues := make([]any, 0, len(groupIDs))
	for _, id := range groupIDs {
		groupValues = append(groupValues, id)
placeholder
	return map[string]any{
		"account_ids": accountValues,
		"group_ids":   groupValues,
placeholder
placeholder

func schedulerBucketsForTest(groupIDs []int64, platforms ...string) []SchedulerBucket {
	buckets := make([]SchedulerBucket, 0, len(groupIDs)*len(platforms)*3)
	for _, platform := range platforms {
		for _, groupID := range groupIDs {
			buckets = append(buckets,
				SchedulerBucket{GroupID: groupID, Platform: platform, Mode: SchedulerModeSingleplaceholder,
				SchedulerBucket{GroupID: groupID, Platform: platform, Mode: SchedulerModeForcedplaceholder,
			)
			if platform == PlatformAnthropic || platform == PlatformGemini {
				buckets = append(buckets, SchedulerBucket{GroupID: groupID, Platform: platform, Mode: SchedulerModeMixedplaceholder)
		placeholder
	placeholder
placeholder
	return buckets
placeholder

func TestSchedulerBulkAccountEventScopesOpenAIRebuildToFreshPlatform(t *testing.T) {
	cache := newBulkEventSnapshotCache()
	repo := newBulkEventAccountRepo(&Account{ID: 1, Platform: PlatformOpenAI, GroupIDs: []int64{12placeholderplaceholder)
	svc := newBulkEventTestService(cache, repo)

	err := svc.handleBulkAccountEvent(context.Background(), bulkEventPayload([]int64{1placeholder, []int64{11placeholder), make(map[batchSeenKey]struct{placeholder))

placeholder
	require.ElementsMatch(t, schedulerBucketsForTest([]int64{11, 12placeholder, PlatformOpenAI), cache.capturedBuckets())
	set, deleted := cache.accountWrites()
	require.Equal(t, []int64{1placeholder, set)
	require.Empty(t, deleted)
placeholder

func TestSchedulerBulkAccountEventScopesCNRebuildToFreshPlatform(t *testing.T) {
	for _, platform := range []string{PlatformKimi, PlatformZhipu, PlatformDeepseekplaceholder {
		t.Run(platform, func(t *testing.T) {
			cache := newBulkEventSnapshotCache()
			repo := newBulkEventAccountRepo(&Account{ID: 1, Platform: platform, GroupIDs: []int64{12placeholderplaceholder)
			svc := newBulkEventTestService(cache, repo)

			err := svc.handleBulkAccountEvent(context.Background(), bulkEventPayload([]int64{1placeholder, []int64{11placeholder), make(map[batchSeenKey]struct{placeholder))

		placeholder
			require.ElementsMatch(t, schedulerBucketsForTest([]int64{11, 12placeholder, platform), cache.capturedBuckets())
	placeholder)
placeholder
placeholder

func TestSchedulerBulkAccountEventRebuildsOpenAIUngroupedBucket(t *testing.T) {
	cache := newBulkEventSnapshotCache()
	repo := newBulkEventAccountRepo(&Account{ID: 6, Platform: PlatformOpenAIplaceholder)
	svc := newBulkEventTestService(cache, repo)

	err := svc.handleBulkAccountEvent(context.Background(), bulkEventPayload([]int64{6placeholder, nil), make(map[batchSeenKey]struct{placeholder))

placeholder
	require.ElementsMatch(t, schedulerBucketsForTest([]int64{0placeholder, PlatformOpenAI), cache.capturedBuckets())
placeholder

func TestSchedulerBulkAccountEventKeepsGroupedAndUngroupedBuckets(t *testing.T) {
	cache := newBulkEventSnapshotCache()
	repo := newBulkEventAccountRepo(
		&Account{ID: 7, Platform: PlatformOpenAI, GroupIDs: []int64{51placeholderplaceholder,
		&Account{ID: 8, Platform: PlatformOpenAIplaceholder,
	)
	svc := newBulkEventTestService(cache, repo)

	err := svc.handleBulkAccountEvent(context.Background(), bulkEventPayload([]int64{7, 8placeholder, nil), make(map[batchSeenKey]struct{placeholder))

placeholder
	require.ElementsMatch(t, schedulerBucketsForTest([]int64{0, 51placeholder, PlatformOpenAI), cache.capturedBuckets())
placeholder

func TestSchedulerBulkAccountEventDoesNotCrossCurrentGroupsBetweenPlatforms(t *testing.T) {
	cache := newBulkEventSnapshotCache()
	repo := newBulkEventAccountRepo(
		&Account{ID: 9, Platform: PlatformOpenAI, GroupIDs: []int64{61placeholderplaceholder,
		&Account{ID: 10, Platform: PlatformGrok, GroupIDs: []int64{62placeholderplaceholder,
	)
	svc := newBulkEventTestService(cache, repo)

	err := svc.handleBulkAccountEvent(context.Background(), bulkEventPayload([]int64{9, 10placeholder, []int64{63placeholder), make(map[batchSeenKey]struct{placeholder))

placeholder
	want := append(
		schedulerBucketsForTest([]int64{61, 63placeholder, PlatformOpenAI),
		schedulerBucketsForTest([]int64{62, 63placeholder, PlatformGrok)...,
	)
	require.ElementsMatch(t, want, cache.capturedBuckets())
placeholder

func TestSchedulerBulkAccountEventUsesGroupZeroInSimpleMode(t *testing.T) {
	cache := newBulkEventSnapshotCache()
	repo := newBulkEventAccountRepo(&Account{ID: 11, Platform: PlatformOpenAI, GroupIDs: []int64{71placeholderplaceholder)
	svc := NewSchedulerSnapshotService(cache, nil, repo, nil, &config.Config{RunMode: config.RunModeSimpleplaceholder)

	err := svc.handleBulkAccountEvent(context.Background(), bulkEventPayload([]int64{11placeholder, []int64{72placeholder), make(map[batchSeenKey]struct{placeholder))

placeholder
	require.ElementsMatch(t, schedulerBucketsForTest([]int64{0placeholder, PlatformOpenAI), cache.capturedBuckets())
placeholder

func TestSchedulerBulkAccountEventConservativelyExpandsAntigravityPlatforms(t *testing.T) {
	cache := newBulkEventSnapshotCache()
	// fresh 值可能已经关闭 mixed_scheduling，兼容平台仍要重建以清理旧快照。
	repo := newBulkEventAccountRepo(&Account{ID: 2, Platform: PlatformAntigravity, GroupIDs: []int64{22placeholderplaceholder)
	svc := newBulkEventTestService(cache, repo)

	err := svc.handleBulkAccountEvent(context.Background(), bulkEventPayload([]int64{2placeholder, []int64{21placeholder), make(map[batchSeenKey]struct{placeholder))

placeholder
	require.ElementsMatch(t,
		schedulerBucketsForTest([]int64{21, 22placeholder, PlatformAnthropic, PlatformGemini, PlatformAntigravity),
		cache.capturedBuckets(),
	)
placeholder

func TestSchedulerBulkAccountEventMissingAccountFallsBackToAllPlatforms(t *testing.T) {
	cache := newBulkEventSnapshotCache()
	repo := newBulkEventAccountRepo(&Account{ID: 3, Platform: PlatformOpenAI, GroupIDs: []int64{32placeholderplaceholder)
	svc := newBulkEventTestService(cache, repo)

	err := svc.handleBulkAccountEvent(context.Background(), bulkEventPayload([]int64{3, 4placeholder, []int64{31placeholder), make(map[batchSeenKey]struct{placeholder))

placeholder
	platforms := schedulerSnapshotPlatforms()
	require.ElementsMatch(t, schedulerBucketsForTest([]int64{31, 32placeholder, platforms[:]...), cache.capturedBuckets())
	set, deleted := cache.accountWrites()
	require.Equal(t, []int64{3placeholder, set)
	require.Equal(t, []int64{4placeholder, deleted)
placeholder

func TestSchedulerBulkAccountEventUnknownPlatformFallsBackToAllPlatforms(t *testing.T) {
	cache := newBulkEventSnapshotCache()
	repo := newBulkEventAccountRepo(&Account{ID: 5, GroupIDs: []int64{42placeholderplaceholder)
	svc := newBulkEventTestService(cache, repo)

	err := svc.handleBulkAccountEvent(context.Background(), bulkEventPayload([]int64{5placeholder, []int64{41placeholder), make(map[batchSeenKey]struct{placeholder))

placeholder
	platforms := schedulerSnapshotPlatforms()
	require.ElementsMatch(t, schedulerBucketsForTest([]int64{41, 42placeholder, platforms[:]...), cache.capturedBuckets())
placeholder
