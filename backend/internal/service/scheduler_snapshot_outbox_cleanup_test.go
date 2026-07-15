package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
)

type outboxCleanupCache struct {
	watermark       int64
	setWatermarks   []int64
	updateErr       error
	listBucketCalls int
placeholder

func (c *outboxCleanupCache) GetSnapshot(ctx context.Context, bucket SchedulerBucket) ([]*Account, bool, error) {
	return nil, false, nil
placeholder

func (c *outboxCleanupCache) CaptureBucketWriteToken(ctx context.Context, bucket SchedulerBucket) (SchedulerBucketWriteToken, error) {
	return SchedulerBucketWriteToken{Bucket: bucket, Epoch: 1placeholder, nil
placeholder

func (c *outboxCleanupCache) SetSnapshot(ctx context.Context, bucket SchedulerBucket, token SchedulerBucketWriteToken, accounts []Account) error {
	return nil
placeholder

func (c *outboxCleanupCache) RetireBucket(ctx context.Context, bucket SchedulerBucket) error {
	return nil
placeholder

func (c *outboxCleanupCache) ReopenBucket(ctx context.Context, bucket SchedulerBucket) (SchedulerBucketWriteToken, error) {
	return SchedulerBucketWriteToken{Bucket: bucket, Epoch: 1placeholder, nil
placeholder

func (c *outboxCleanupCache) TryAcquireGroupLifecycleLease(context.Context, int64, time.Duration) (SchedulerGroupLifecycleLease, bool, error) {
	return SchedulerGroupLifecycleLease{placeholder, false, nil
placeholder

func (c *outboxCleanupCache) ReleaseGroupLifecycleLease(context.Context, SchedulerGroupLifecycleLease) error {
	return nil
placeholder

func (c *outboxCleanupCache) GetAccount(ctx context.Context, accountID int64) (*Account, error) {
	return nil, nil
placeholder

func (c *outboxCleanupCache) SetAccount(ctx context.Context, account *Account) error {
	return nil
placeholder

func (c *outboxCleanupCache) DeleteAccount(ctx context.Context, accountID int64) error {
	return nil
placeholder

func (c *outboxCleanupCache) UpdateLastUsed(ctx context.Context, updates map[int64]time.Time) error {
	return c.updateErr
placeholder

func (c *outboxCleanupCache) TryLockBucket(ctx context.Context, bucket SchedulerBucket, ttl time.Duration) (bool, error) {
	return true, nil
placeholder

func (c *outboxCleanupCache) UnlockBucket(ctx context.Context, bucket SchedulerBucket) error {
	return nil
placeholder

func (c *outboxCleanupCache) ListBuckets(ctx context.Context) ([]SchedulerBucket, error) {
	c.listBucketCalls++
	return nil, nil
placeholder

func (c *outboxCleanupCache) GetOutboxWatermark(ctx context.Context) (int64, error) {
	return c.watermark, nil
placeholder

func (c *outboxCleanupCache) SetOutboxWatermark(ctx context.Context, id int64) error {
	c.watermark = id
	c.setWatermarks = append(c.setWatermarks, id)
	return nil
placeholder

type outboxCleanupDeleteCall struct {
	watermark int64
	limit     int
placeholder

type outboxCleanupRepo struct {
	events              []SchedulerOutboxEvent
	rows                []int64
	lockAcquired        bool
	lockAttempts        int
	releaseCount        int
	deleteCalls         []outboxCleanupDeleteCall
	firstCreatedAfterID []int64
placeholder

func (r *outboxCleanupRepo) ListAfterAndReleaseDedup(ctx context.Context, afterID int64, limit int) ([]SchedulerOutboxEvent, error) {
	events := make([]SchedulerOutboxEvent, 0, len(r.events))
	for _, event := range r.events {
		if event.ID <= afterID {
			continue
	placeholder
		events = append(events, event)
		if limit > 0 && len(events) >= limit {
			break
	placeholder
placeholder
	return events, nil
placeholder

func (r *outboxCleanupRepo) FirstCreatedAtAfter(ctx context.Context, afterID int64) (time.Time, bool, error) {
	r.firstCreatedAfterID = append(r.firstCreatedAfterID, afterID)
	for _, event := range r.events {
		if event.ID > afterID {
			return event.CreatedAt, true, nil
	placeholder
placeholder
	return time.Time{placeholder, false, nil
placeholder

func (r *outboxCleanupRepo) MaxID(ctx context.Context) (int64, error) {
	var maxID int64
	for _, id := range r.rows {
		if id > maxID {
			maxID = id
	placeholder
placeholder
	return maxID, nil
placeholder

func (r *outboxCleanupRepo) DeleteConsumedUpTo(ctx context.Context, watermark int64, limit int) (int64, error) {
	r.deleteCalls = append(r.deleteCalls, outboxCleanupDeleteCall{
		watermark: watermark,
		limit:     limit,
placeholder)
	if watermark <= 0 || limit <= 0 {
		return 0, nil
placeholder

	deleted := int64(0)
	kept := make([]int64, 0, len(r.rows))
	for _, id := range r.rows {
		if id <= watermark && deleted < int64(limit) {
			deleted++
			continue
	placeholder
		kept = append(kept, id)
placeholder
	r.rows = kept
	return deleted, nil
placeholder

func (r *outboxCleanupRepo) TryAcquireCleanupLock(ctx context.Context) (SchedulerOutboxCleanupLease, bool, error) {
	r.lockAttempts++
	if !r.lockAcquired {
		return nil, false, nil
placeholder
	return outboxCleanupLease{release: func() {
		r.releaseCount++
placeholderplaceholder, true, nil
placeholder

type outboxCleanupLease struct {
	release func()
placeholder

func (l outboxCleanupLease) Release() {
	if l.release != nil {
		l.release()
placeholder
placeholder

func TestSchedulerSnapshotServicePollOutboxCleansConsumedRowsAfterWatermark(t *testing.T) {
	cache := &outboxCleanupCache{placeholder
	repo := &outboxCleanupRepo{
		events: []SchedulerOutboxEvent{
			{ID: 10000, EventType: SchedulerOutboxEventAccountLastUsedplaceholder,
	placeholder,
		rows:         int64Range(1, 10003),
		lockAcquired: true,
placeholder
	svc := NewSchedulerSnapshotService(cache, repo, nil, nil, nil)

	svc.pollOutbox()

	if cache.watermark != 10000 {
		t.Fatalf("expected watermark 10000, got %d", cache.watermark)
placeholder
	if !reflect.DeepEqual(cache.setWatermarks, []int64{10000placeholder) {
		t.Fatalf("unexpected watermark writes: %#v", cache.setWatermarks)
placeholder
	if !reflect.DeepEqual(repo.rows, []int64{10001, 10002, 10003placeholder) {
		t.Fatalf("expected rows above watermark to remain, got %#v", repo.rows)
placeholder
	if repo.lockAttempts != 1 || repo.releaseCount != 1 {
		t.Fatalf("expected one lock acquire/release, got acquire=%d release=%d", repo.lockAttempts, repo.releaseCount)
placeholder
	if len(repo.deleteCalls) != 3 {
		t.Fatalf("expected cleanup to loop until a short batch, got %d calls", len(repo.deleteCalls))
placeholder
	for _, call := range repo.deleteCalls {
		if call.watermark != 10000 || call.limit != schedulerOutboxCleanupBatch {
			t.Fatalf("unexpected cleanup call: %#v", call)
	placeholder
placeholder
placeholder

func TestSchedulerSnapshotServicePollOutboxSkipsCleanupWhenLockUnavailable(t *testing.T) {
	cache := &outboxCleanupCache{placeholder
	repo := &outboxCleanupRepo{
		events: []SchedulerOutboxEvent{
			{ID: 3, EventType: SchedulerOutboxEventAccountLastUsedplaceholder,
	placeholder,
		rows:         []int64{1, 2, 3, 4placeholder,
		lockAcquired: false,
placeholder
	svc := NewSchedulerSnapshotService(cache, repo, nil, nil, nil)

	svc.pollOutbox()

	if cache.watermark != 3 {
		t.Fatalf("expected watermark 3, got %d", cache.watermark)
placeholder
	if !reflect.DeepEqual(repo.rows, []int64{1, 2, 3, 4placeholder) {
		t.Fatalf("expected cleanup to skip all rows, got %#v", repo.rows)
placeholder
	if repo.lockAttempts != 1 {
		t.Fatalf("expected one lock attempt, got %d", repo.lockAttempts)
placeholder
	if len(repo.deleteCalls) != 0 {
		t.Fatalf("expected no delete calls, got %#v", repo.deleteCalls)
placeholder
	if repo.releaseCount != 0 {
		t.Fatalf("expected no release without lock, got %d", repo.releaseCount)
placeholder
placeholder

func TestSchedulerSnapshotServicePollOutboxDoesNotCleanupOnHandleFailure(t *testing.T) {
	cache := &outboxCleanupCache{
		updateErr: errors.New("cache update failed"),
placeholder
	repo := &outboxCleanupRepo{
		events: []SchedulerOutboxEvent{
			{
				ID:        5,
				EventType: SchedulerOutboxEventAccountLastUsed,
				Payload: map[string]any{
					"last_used": map[string]any{"101": float64(123)placeholder,
			placeholder,
		placeholder,
	placeholder,
		rows:         []int64{1, 2, 3, 4, 5, 6placeholder,
		lockAcquired: true,
placeholder
	svc := NewSchedulerSnapshotService(cache, repo, nil, nil, nil)

	svc.pollOutbox()

	if len(cache.setWatermarks) != 0 {
		t.Fatalf("expected no watermark write on handle failure, got %#v", cache.setWatermarks)
placeholder
	if repo.lockAttempts != 0 {
		t.Fatalf("expected cleanup lock not to be attempted, got %d", repo.lockAttempts)
placeholder
	if len(repo.deleteCalls) != 0 {
		t.Fatalf("expected no delete calls, got %#v", repo.deleteCalls)
placeholder
	if !reflect.DeepEqual(repo.rows, []int64{1, 2, 3, 4, 5, 6placeholder) {
		t.Fatalf("expected rows unchanged, got %#v", repo.rows)
placeholder
placeholder

func TestSchedulerSnapshotServicePollOutboxDoesNotUseConsumedEventForLag(t *testing.T) {
	cache := &outboxCleanupCache{placeholder
	repo := &outboxCleanupRepo{
		events: []SchedulerOutboxEvent{
			{
				ID:        7,
				EventType: SchedulerOutboxEventAccountLastUsed,
				CreatedAt: time.Now().Add(-time.Hour),
		placeholder,
	placeholder,
placeholder
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			Scheduling: config.GatewaySchedulingConfig{
				OutboxLagWarnSeconds:     1,
				OutboxLagRebuildSeconds:  1,
				OutboxLagRebuildFailures: 1,
		placeholder,
	placeholder,
placeholder
	svc := NewSchedulerSnapshotService(cache, repo, nil, nil, cfg)

	svc.pollOutbox()

	if cache.watermark != 7 {
		t.Fatalf("expected watermark 7, got %d", cache.watermark)
placeholder
	if !reflect.DeepEqual(repo.firstCreatedAfterID, []int64{7placeholder) {
		t.Fatalf("expected lag check after consumed watermark, got %#v", repo.firstCreatedAfterID)
placeholder
	if cache.listBucketCalls != 0 {
		t.Fatalf("expected consumed event not to trigger full rebuild, got %d attempts", cache.listBucketCalls)
placeholder
	if svc.lagFailures != 0 {
		t.Fatalf("expected lag failures to remain reset, got %d", svc.lagFailures)
placeholder
placeholder

func TestSchedulerSnapshotServiceCleanupSkipsNonPositiveWatermark(t *testing.T) {
	repo := &outboxCleanupRepo{
		rows:         []int64{1, 2, 3placeholder,
		lockAcquired: true,
placeholder
	svc := NewSchedulerSnapshotService(&outboxCleanupCache{placeholder, repo, nil, nil, nil)

	svc.cleanupConsumedOutbox(0)

	if repo.lockAttempts != 0 {
		t.Fatalf("expected no lock attempt for non-positive watermark, got %d", repo.lockAttempts)
placeholder
	if len(repo.deleteCalls) != 0 {
		t.Fatalf("expected no delete calls, got %#v", repo.deleteCalls)
placeholder
	if !reflect.DeepEqual(repo.rows, []int64{1, 2, 3placeholder) {
		t.Fatalf("expected rows unchanged, got %#v", repo.rows)
placeholder
placeholder

func int64Range(start, end int64) []int64 {
	values := make([]int64, 0, end-start+1)
	for id := start; id <= end; id++ {
		values = append(values, id)
placeholder
	return values
placeholder
