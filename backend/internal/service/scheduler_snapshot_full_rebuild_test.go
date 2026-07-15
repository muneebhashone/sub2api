package service

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type schedulerFullRebuildTestCache struct {
	SchedulerCache

	mu        sync.Mutex
	listErr   error
	listCalls int
	captures  int
	lockCalls int
placeholder

func (c *schedulerFullRebuildTestCache) ListBuckets(context.Context) ([]SchedulerBucket, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.listCalls++
	return nil, c.listErr
placeholder

func (c *schedulerFullRebuildTestCache) TryLockBucket(context.Context, SchedulerBucket, time.Duration) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lockCalls++
	return false, nil
placeholder

func (c *schedulerFullRebuildTestCache) CaptureBucketWriteToken(_ context.Context, bucket SchedulerBucket) (SchedulerBucketWriteToken, error) {
	c.mu.Lock()
	c.captures++
	c.mu.Unlock()
	return SchedulerBucketWriteToken{Bucket: bucket, Epoch: 1placeholder, nil
placeholder

func (c *schedulerFullRebuildTestCache) ReopenBucket(_ context.Context, bucket SchedulerBucket) (SchedulerBucketWriteToken, error) {
	return SchedulerBucketWriteToken{Bucket: bucket, Epoch: 1placeholder, nil
placeholder

func TestSchedulerSnapshotServiceFullRebuildCoalescesConcurrentRequestsIntoTrailingRun(t *testing.T) {
	svc := &SchedulerSnapshotService{placeholder
	wantTrailingErr := errors.New("trailing rebuild failed")
	firstStarted := make(chan struct{placeholder)
	releaseFirst := make(chan struct{placeholder)
	var releaseOnce sync.Once
	release := func() {
		releaseOnce.Do(func() { close(releaseFirst) placeholder)
placeholder
	defer release()

	var calls atomic.Int32
	var active atomic.Int32
	var maxActive atomic.Int32
	run := func() error {
		call := calls.Add(1)
		currentActive := active.Add(1)
		defer active.Add(-1)
		for {
			previousMax := maxActive.Load()
			if currentActive <= previousMax || maxActive.CompareAndSwap(previousMax, currentActive) {
				break
		placeholder
	placeholder
		if call == 1 {
			close(firstStarted)
			<-releaseFirst
			return nil
	placeholder
		return wantTrailingErr
placeholder

	firstResult := make(chan error, 1)
	go func() {
		firstResult <- svc.coalesceFullRebuild(run)
placeholder()

	select {
	case <-firstStarted:
	case <-time.After(time.Second):
		t.Fatal("first rebuild did not start")
placeholder

	const followers = 20
	followerResults := make(chan error, followers)
	for range followers {
		go func() {
			followerResults <- svc.coalesceFullRebuild(run)
	placeholder()
placeholder

	require.Eventually(t, func() bool {
		requested, _ := schedulerFullRebuildState(svc)
		return requested == followers+1
placeholder, time.Second, time.Millisecond)
	release()

	require.NoError(t, <-firstResult)
	for range followers {
		require.ErrorIs(t, <-followerResults, wantTrailingErr)
placeholder
	require.EqualValues(t, 2, calls.Load())
	require.EqualValues(t, 1, maxActive.Load())
	requested, completed := schedulerFullRebuildState(svc)
	require.EqualValues(t, followers+1, requested)
	require.Equal(t, requested, completed)
placeholder

func TestSchedulerSnapshotServiceFullRebuildRunsAgainForSequentialRequest(t *testing.T) {
	svc := &SchedulerSnapshotService{placeholder
	wantSecondErr := errors.New("second rebuild failed")
	var calls atomic.Int32
	run := func() error {
		if calls.Add(1) == 2 {
			return wantSecondErr
	placeholder
		return nil
placeholder

	require.NoError(t, svc.coalesceFullRebuild(run))
	require.ErrorIs(t, svc.coalesceFullRebuild(run), wantSecondErr)
	require.EqualValues(t, 2, calls.Load())
	requested, completed := schedulerFullRebuildState(svc)
	require.EqualValues(t, 2, requested)
	require.Equal(t, requested, completed)
placeholder

func TestSchedulerSnapshotServiceInitialFullRebuildFailsClosedWhenListBucketsFails(t *testing.T) {
	cache := &schedulerFullRebuildTestCache{listErr: errors.New("list buckets failed")placeholder
	svc := NewSchedulerSnapshotService(cache, nil, nil, nil, nil)

	svc.runInitialRebuild()

	cache.mu.Lock()
	listCalls := cache.listCalls
	captures := cache.captures
	lockCalls := cache.lockCalls
	cache.mu.Unlock()
	require.Equal(t, 1, listCalls)
	require.Zero(t, captures)
	require.Zero(t, lockCalls)
	requested, completed := schedulerFullRebuildState(svc)
	require.EqualValues(t, 1, requested)
	require.Equal(t, requested, completed)
	svc.fullRebuildStateMu.Lock()
	require.ErrorIs(t, svc.fullRebuildLastErr, cache.listErr)
	svc.fullRebuildStateMu.Unlock()
placeholder

func schedulerFullRebuildState(svc *SchedulerSnapshotService) (requested uint64, completed uint64) {
	svc.fullRebuildStateMu.Lock()
	defer svc.fullRebuildStateMu.Unlock()
	return svc.fullRebuildRequested, svc.fullRebuildCompleted
placeholder
