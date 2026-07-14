package service

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
)

type outboxCleanupCache struct {
	watermark       int64
	setWatermarks   []int64
	updateErr       error
	listBucketErr   error
	listBuckets     []SchedulerBucket
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
	return c.listBuckets, c.listBucketErr
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
	maxIDCalls          int
	maxIDErr            error
	lockAcquired        bool
	lockAttempts        int
	releaseCount        int
	deleteCalls         []outboxCleanupDeleteCall
	firstCreatedAfterID []int64
placeholder

type outboxCleanupAccountRepo struct {
	AccountRepository
placeholder

func (r *outboxCleanupAccountRepo) ListSchedulableUngroupedByPlatform(context.Context, string) ([]Account, error) {
	return nil, nil
placeholder

type blockingOutboxCleanupCache struct {
	*outboxCleanupCache
	mu      sync.Mutex
	calls   int
	started chan struct{placeholder
	release chan struct{placeholder
placeholder

func (c *blockingOutboxCleanupCache) ListBuckets(context.Context) ([]SchedulerBucket, error) {
	c.mu.Lock()
	c.calls++
	call := c.calls
	c.mu.Unlock()
	if call == 1 {
		close(c.started)
		<-c.release
placeholder
	return c.listBuckets, c.listBucketErr
placeholder

func (c *blockingOutboxCleanupCache) listCalls() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.calls
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
	r.maxIDCalls++
	if r.maxIDErr != nil {
		return 0, r.maxIDErr
placeholder
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

func TestSchedulerSnapshotServiceCheckOutboxLagLatchesPersistentDegradation(t *testing.T) {
	tests := []struct {
		name             string
		createdAt        time.Time
		rows             []int64
		lagSeconds       int
		backlogThreshold int
placeholder{
		{
			name:       "lag",
			createdAt:  time.Now().Add(-time.Hour),
			rows:       []int64{1placeholder,
			lagSeconds: 1,
	placeholder,
		{
			name:             "backlog",
			createdAt:        time.Now(),
			rows:             []int64{100placeholder,
			backlogThreshold: 50,
	placeholder,
		{
			name:             "lag_and_backlog",
			createdAt:        time.Now().Add(-time.Hour),
			rows:             []int64{100placeholder,
			lagSeconds:       1,
			backlogThreshold: 50,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := &outboxCleanupCache{listBuckets: []SchedulerBucket{{Platform: PlatformOpenAI, Mode: SchedulerModeSingleplaceholderplaceholderplaceholder
			repo := &outboxCleanupRepo{
				events: []SchedulerOutboxEvent{{ID: 1, CreatedAt: tt.createdAtplaceholderplaceholder,
				rows:   tt.rows,
		placeholder
			cfg := &config.Config{
				Gateway: config.GatewayConfig{
					Scheduling: config.GatewaySchedulingConfig{
						OutboxLagRebuildSeconds:  tt.lagSeconds,
						OutboxLagRebuildFailures: 1,
						OutboxBacklogRebuildRows: tt.backlogThreshold,
				placeholder,
			placeholder,
		placeholder
			svc := NewSchedulerSnapshotService(cache, repo, &outboxCleanupAccountRepo{placeholder, nil, cfg)

			for range 3 {
				svc.checkOutboxLag(context.Background(), 0)
		placeholder

			if cache.listBucketCalls != 1 {
				t.Fatalf("expected one rebuild attempt during a persistent degraded episode, got %d", cache.listBucketCalls)
		placeholder
	placeholder)
placeholder
placeholder

func TestSchedulerSnapshotServiceCheckOutboxLagFailedRebuildRearmsAfterRecovery(t *testing.T) {
	cache := &outboxCleanupCache{listBucketErr: errors.New("list buckets failed")placeholder
	repo := &outboxCleanupRepo{
		events: []SchedulerOutboxEvent{{ID: 1, CreatedAt: time.Now().Add(-time.Hour)placeholderplaceholder,
		rows:   []int64{1placeholder,
placeholder
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			Scheduling: config.GatewaySchedulingConfig{
				OutboxLagRebuildSeconds:  1,
				OutboxLagRebuildFailures: 1,
		placeholder,
	placeholder,
placeholder
	svc := NewSchedulerSnapshotService(cache, repo, nil, nil, cfg)

	svc.checkOutboxLag(context.Background(), 0)
	svc.checkOutboxLag(context.Background(), 0)
	if cache.listBucketCalls != 1 {
		t.Fatalf("expected a failed rebuild to stay bounded within the episode, got %d attempts", cache.listBucketCalls)
placeholder

	svc.checkOutboxLag(context.Background(), 1)
	repo.events = append(repo.events, SchedulerOutboxEvent{ID: 2, CreatedAt: time.Now().Add(-time.Hour)placeholder)
	repo.rows = []int64{2placeholder
	svc.checkOutboxLag(context.Background(), 1)

	if cache.listBucketCalls != 2 {
		t.Fatalf("expected recovery to rearm a failed rebuild for the next episode, got %d attempts", cache.listBucketCalls)
placeholder
placeholder

func TestSchedulerSnapshotServiceCheckOutboxLagFailedRebuildRetriesAfterCooldownWithoutRecovery(t *testing.T) {
	cache := &outboxCleanupCache{
		listBucketErr: errors.New("list buckets failed"),
		listBuckets:   []SchedulerBucket{{Platform: PlatformOpenAI, Mode: SchedulerModeSingleplaceholderplaceholder,
placeholder
	repo := &outboxCleanupRepo{
		events: []SchedulerOutboxEvent{{ID: 1, CreatedAt: time.Now().Add(-time.Hour)placeholderplaceholder,
		rows:   []int64{1placeholder,
placeholder
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			Scheduling: config.GatewaySchedulingConfig{
				OutboxLagRebuildSeconds:    1,
				OutboxLagRebuildFailures:   1,
				FullRebuildIntervalSeconds: 0,
		placeholder,
	placeholder,
placeholder
	svc := NewSchedulerSnapshotService(cache, repo, &outboxCleanupAccountRepo{placeholder, nil, cfg)

	svc.checkOutboxLag(context.Background(), 0)
	for range 3 {
		svc.checkOutboxLag(context.Background(), 0)
placeholder
	if cache.listBucketCalls != 1 {
		t.Fatalf("expected failed rebuild polls to be rate limited, got %d attempts", cache.listBucketCalls)
placeholder

	svc.lagMu.Lock()
	if !svc.outboxRebuildRetryAt.After(time.Now()) {
		t.Fatal("expected failed rebuild to schedule a future retry")
placeholder
	svc.outboxRebuildRetryAt = time.Now().Add(-time.Second)
	svc.lagMu.Unlock()
	svc.checkOutboxLag(context.Background(), 0)
	if cache.listBucketCalls != 2 {
		t.Fatalf("expected persistent degradation to retry after cooldown, got %d attempts", cache.listBucketCalls)
placeholder

	for range 3 {
		svc.checkOutboxLag(context.Background(), 0)
placeholder
	if cache.listBucketCalls != 2 {
		t.Fatalf("expected repeated rebuild failures to stay rate limited, got %d attempts", cache.listBucketCalls)
placeholder

	svc.lagMu.Lock()
	svc.outboxRebuildRetryAt = time.Now().Add(-time.Second)
	svc.lagMu.Unlock()
	cache.listBucketErr = nil
	svc.checkOutboxLag(context.Background(), 0)
	if cache.listBucketCalls != 3 {
		t.Fatalf("expected degraded episode to retry after cooldown, got %d attempts", cache.listBucketCalls)
placeholder

	for range 3 {
		svc.checkOutboxLag(context.Background(), 0)
placeholder
	if cache.listBucketCalls != 3 {
		t.Fatalf("expected successful retry to latch the degraded episode, got %d attempts", cache.listBucketCalls)
placeholder
placeholder

func TestSchedulerSnapshotServiceCheckOutboxLagBacklogRetryDoesNotBypassNewLagThreshold(t *testing.T) {
	cache := &outboxCleanupCache{listBucketErr: errors.New("list buckets failed")placeholder
	repo := &outboxCleanupRepo{
		events: []SchedulerOutboxEvent{{ID: 1, CreatedAt: time.Now()placeholderplaceholder,
		rows:   []int64{100placeholder,
placeholder
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			Scheduling: config.GatewaySchedulingConfig{
				OutboxLagRebuildSeconds:  1,
				OutboxLagRebuildFailures: 3,
				OutboxBacklogRebuildRows: 50,
		placeholder,
	placeholder,
placeholder
	svc := NewSchedulerSnapshotService(cache, repo, nil, nil, cfg)

	// Start with backlog-only degradation and leave its failed rebuild retry due.
	svc.checkOutboxLag(context.Background(), 0)
	if cache.listBucketCalls != 1 {
		t.Fatalf("expected the backlog degradation to attempt one rebuild, got %d", cache.listBucketCalls)
placeholder
	svc.lagMu.Lock()
	svc.outboxRebuildRetryAt = time.Now().Add(-time.Second)
	svc.lagMu.Unlock()

	// The backlog recovers while lag becomes newly degraded. The stale backlog
	// retry must not make the first lag observation bypass its failure threshold.
	repo.rows = []int64{1placeholder
	repo.events[0].CreatedAt = time.Now().Add(-time.Hour)
	svc.checkOutboxLag(context.Background(), 0)
	if cache.listBucketCalls != 1 {
		t.Fatalf("expected the new lag episode to start at its own threshold, got %d rebuild attempts", cache.listBucketCalls)
placeholder

	svc.checkOutboxLag(context.Background(), 0)
	svc.checkOutboxLag(context.Background(), 0)
	if cache.listBucketCalls != 2 {
		t.Fatalf("expected lag rebuild only after three lag observations, got %d attempts", cache.listBucketCalls)
placeholder
placeholder

func TestSchedulerSnapshotServiceCheckOutboxLagLagRetryDoesNotDelayOrEscalateNewBacklog(t *testing.T) {
	cache := &outboxCleanupCache{listBucketErr: errors.New("list buckets failed")placeholder
	repo := &outboxCleanupRepo{
		events: []SchedulerOutboxEvent{{ID: 1, CreatedAt: time.Now().Add(-time.Hour)placeholderplaceholder,
		rows:   []int64{1placeholder,
placeholder
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			Scheduling: config.GatewaySchedulingConfig{
				OutboxLagRebuildSeconds:  1,
				OutboxLagRebuildFailures: 1,
				OutboxBacklogRebuildRows: 50,
		placeholder,
	placeholder,
placeholder
	svc := NewSchedulerSnapshotService(cache, repo, nil, nil, cfg)

	// Start with lag-only degradation and a failed rebuild in cooldown.
	svc.checkOutboxLag(context.Background(), 0)
	if cache.listBucketCalls != 1 {
		t.Fatalf("expected the lag degradation to attempt one rebuild, got %d", cache.listBucketCalls)
placeholder

	// Lag recovers while backlog becomes newly degraded. It must start immediately
	// and its first failure must use the base retry generation, not lag's count.
	repo.events[0].CreatedAt = time.Now()
	repo.rows = []int64{100placeholder
	svc.checkOutboxLag(context.Background(), 0)
	if cache.listBucketCalls != 2 {
		t.Fatalf("expected the new backlog degradation not to inherit lag cooldown, got %d rebuild attempts", cache.listBucketCalls)
placeholder
	svc.lagMu.Lock()
	failures := svc.outboxRebuildFailures
	svc.lagMu.Unlock()
	if failures != 1 {
		t.Fatalf("expected backlog retry failures to restart at one, got %d", failures)
placeholder
placeholder

func TestSchedulerSnapshotServiceCheckOutboxLagBacklogRetrySurvivesUnknownBacklog(t *testing.T) {
	cache := &outboxCleanupCache{listBucketErr: errors.New("list buckets failed")placeholder
	repo := &outboxCleanupRepo{
		events: []SchedulerOutboxEvent{{ID: 1, CreatedAt: time.Now()placeholderplaceholder,
		rows:   []int64{100placeholder,
placeholder
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			Scheduling: config.GatewaySchedulingConfig{
				OutboxBacklogRebuildRows: 50,
		placeholder,
	placeholder,
placeholder
	svc := NewSchedulerSnapshotService(cache, repo, nil, nil, cfg)

	// A failed backlog rebuild starts a reason-scoped cooldown.
	svc.checkOutboxLag(context.Background(), 0)
	if cache.listBucketCalls != 1 {
		t.Fatalf("expected one initial backlog rebuild, got %d", cache.listBucketCalls)
placeholder
	svc.lagMu.Lock()
	retryAt := svc.outboxRebuildRetryAt
	svc.lagMu.Unlock()
	if !retryAt.After(time.Now()) {
		t.Fatalf("expected a future backlog retry, got %s", retryAt)
placeholder

	// A temporary MaxID failure makes backlog health unknown, not recovered.
	repo.maxIDErr = errors.New("max id unavailable")
	svc.checkOutboxLag(context.Background(), 0)
	svc.lagMu.Lock()
	retryReason := svc.outboxRebuildRetryReason
	failures := svc.outboxRebuildFailures
	retryAtAfterUnknown := svc.outboxRebuildRetryAt
	svc.lagMu.Unlock()
	if retryReason != "outbox_backlog" || failures != 1 || !retryAtAfterUnknown.Equal(retryAt) {
		t.Fatalf("expected unknown backlog to preserve retry state, got reason=%q failures=%d retry_at=%s", retryReason, failures, retryAtAfterUnknown)
placeholder

	// When MaxID recovers and backlog remains degraded, the original cooldown
	// still applies; only an expired cooldown may trigger the retry.
	repo.maxIDErr = nil
	svc.checkOutboxLag(context.Background(), 0)
	if cache.listBucketCalls != 1 {
		t.Fatalf("expected backlog recovery before cooldown to stay rate limited, got %d attempts", cache.listBucketCalls)
placeholder
	svc.lagMu.Lock()
	svc.outboxRebuildRetryAt = time.Now().Add(-time.Second)
	svc.lagMu.Unlock()
	svc.checkOutboxLag(context.Background(), 0)
	if cache.listBucketCalls != 2 {
		t.Fatalf("expected backlog retry after cooldown expiry, got %d attempts", cache.listBucketCalls)
placeholder
placeholder

func TestSchedulerSnapshotServiceCheckOutboxLagPreemptsUnknownBacklogRetryAtThreshold(t *testing.T) {
	cache := &outboxCleanupCache{listBucketErr: errors.New("list buckets failed")placeholder
	repo := &outboxCleanupRepo{
		events: []SchedulerOutboxEvent{{ID: 1, CreatedAt: time.Now()placeholderplaceholder,
		rows:   []int64{100placeholder,
placeholder
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			Scheduling: config.GatewaySchedulingConfig{
				OutboxLagRebuildSeconds:  1,
				OutboxLagRebuildFailures: 3,
				OutboxBacklogRebuildRows: 50,
		placeholder,
	placeholder,
placeholder
	svc := NewSchedulerSnapshotService(cache, repo, nil, nil, cfg)

	// Backlog starts the first failed rebuild generation and remains unknown.
	svc.checkOutboxLag(context.Background(), 0)
	if cache.listBucketCalls != 1 {
		t.Fatalf("expected one initial backlog rebuild, got %d", cache.listBucketCalls)
placeholder
	repo.maxIDErr = errors.New("max id unavailable")
	repo.events[0].CreatedAt = time.Now().Add(-time.Hour)

	// A known lag degradation must keep accumulating independently of the active
	// backlog cooldown and preempt it only after reaching its own threshold.
	for observation := 1; observation <= 2; observation++ {
		svc.checkOutboxLag(context.Background(), 0)
		if cache.listBucketCalls != 1 {
			t.Fatalf("expected lag observation %d to stay below threshold, got %d rebuild attempts", observation, cache.listBucketCalls)
	placeholder
placeholder
	svc.checkOutboxLag(context.Background(), 0)
	if cache.listBucketCalls != 2 {
		t.Fatalf("expected lag to preempt backlog cooldown at its threshold, got %d attempts", cache.listBucketCalls)
placeholder

	svc.lagMu.Lock()
	retryReason := svc.outboxRebuildRetryReason
	failures := svc.outboxRebuildFailures
	retryAt := svc.outboxRebuildRetryAt
	svc.lagMu.Unlock()
	if retryReason != "outbox_lag" || failures != 1 || !retryAt.After(time.Now()) {
		t.Fatalf("expected a fresh lag retry generation, got reason=%q failures=%d retry_at=%s", retryReason, failures, retryAt)
placeholder
placeholder

func TestOutboxRebuildRetryDelayIsExponentiallyBounded(t *testing.T) {
	previous := time.Duration(0)
	for failures := 1; failures <= 20; failures++ {
		delay := outboxRebuildRetryDelay(failures)
		if delay < previous {
			t.Fatalf("expected retry delay to be monotonic, failure %d produced %s after %s", failures, delay, previous)
	placeholder
		if delay > outboxRebuildRetryMaxDelay {
			t.Fatalf("expected retry delay to stay bounded, got %s", delay)
	placeholder
		previous = delay
placeholder
	if previous != outboxRebuildRetryMaxDelay {
		t.Fatalf("expected repeated failures to reach max delay %s, got %s", outboxRebuildRetryMaxDelay, previous)
placeholder
placeholder

func TestSchedulerSnapshotServicePollOutboxEmptyBatchClearsDegradedEpisode(t *testing.T) {
	cache := &outboxCleanupCache{listBuckets: []SchedulerBucket{{Platform: PlatformOpenAI, Mode: SchedulerModeSingleplaceholderplaceholderplaceholder
	repo := &outboxCleanupRepo{
		events: []SchedulerOutboxEvent{{ID: 1, CreatedAt: time.Now().Add(-time.Hour)placeholderplaceholder,
		rows:   []int64{1placeholder,
placeholder
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			Scheduling: config.GatewaySchedulingConfig{
				OutboxLagRebuildSeconds:  1,
				OutboxLagRebuildFailures: 1,
				OutboxBacklogRebuildRows: 1,
		placeholder,
	placeholder,
placeholder
	svc := NewSchedulerSnapshotService(cache, repo, &outboxCleanupAccountRepo{placeholder, nil, cfg)

	svc.checkOutboxLag(context.Background(), 0)
	cache.watermark = 1
	svc.pollOutbox()

	if !reflect.DeepEqual(repo.firstCreatedAfterID, []int64{0placeholder) {
		t.Fatalf("expected empty poll to use the empty batch as recovery evidence, got watermarks %#v", repo.firstCreatedAfterID)
placeholder
	if repo.maxIDCalls != 1 {
		t.Fatalf("expected empty poll to skip a redundant backlog query, got %d health checks", repo.maxIDCalls)
placeholder

	repo.events = append(repo.events, SchedulerOutboxEvent{ID: 2, CreatedAt: time.Now().Add(-time.Hour)placeholder)
	repo.rows = []int64{2placeholder
	svc.checkOutboxLag(context.Background(), 1)
	if cache.listBucketCalls != 2 {
		t.Fatalf("expected empty-poll recovery to rearm the next degraded episode, got %d attempts", cache.listBucketCalls)
placeholder
placeholder

func TestSchedulerSnapshotServiceOutboxLagWarningIsTransitionLimited(t *testing.T) {
	svc := NewSchedulerSnapshotService(nil, nil, nil, nil, nil)

	if !svc.shouldLogOutboxLagWarning(true) {
		t.Fatal("expected the initial degraded transition to log")
placeholder
	if svc.shouldLogOutboxLagWarning(true) {
		t.Fatal("expected persistent degradation to suppress repeated warnings")
placeholder
	if svc.shouldLogOutboxLagWarning(true) {
		t.Fatal("expected persistent degradation to suppress repeated warnings")
placeholder
	if svc.shouldLogOutboxLagWarning(false) {
		t.Fatal("expected recovery not to emit a lag warning")
placeholder
	if !svc.shouldLogOutboxLagWarning(true) {
		t.Fatal("expected renewed degradation to log after recovery")
placeholder
placeholder

func TestSchedulerSnapshotServiceCheckOutboxLagSamplesMaxIDErrors(t *testing.T) {
	svc := NewSchedulerSnapshotService(nil, nil, nil, nil, nil)
	now := time.Now()

	if !svc.shouldLogOutboxMaxIDError(now) {
		t.Fatal("expected the first MaxID error to log")
placeholder
	if svc.shouldLogOutboxMaxIDError(now.Add(outboxMaxIDErrorLogSampleInterval / 2)) {
		t.Fatal("expected MaxID errors inside the sample interval to be suppressed")
placeholder
	if !svc.shouldLogOutboxMaxIDError(now.Add(outboxMaxIDErrorLogSampleInterval)) {
		t.Fatal("expected MaxID error logging to rearm after the sample interval")
placeholder
placeholder

func TestSchedulerSnapshotServicePollOutboxHealthyEmptyBatchSkipsLagHealthQueries(t *testing.T) {
	cache := &outboxCleanupCache{placeholder
	repo := &outboxCleanupRepo{placeholder
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			Scheduling: config.GatewaySchedulingConfig{
				OutboxLagRebuildSeconds:  1,
				OutboxLagRebuildFailures: 1,
				OutboxBacklogRebuildRows: 1,
		placeholder,
	placeholder,
placeholder
	svc := NewSchedulerSnapshotService(cache, repo, nil, nil, cfg)

	svc.pollOutbox()

	if len(repo.firstCreatedAfterID) != 0 {
		t.Fatalf("expected healthy empty poll to skip lag query, got watermarks %#v", repo.firstCreatedAfterID)
placeholder
	if repo.maxIDCalls != 0 {
		t.Fatalf("expected healthy empty poll to skip backlog query, got %d calls", repo.maxIDCalls)
placeholder
placeholder

func TestSchedulerSnapshotServiceEmptyPollDoesNotReleaseRunningRebuild(t *testing.T) {
	baseCache := &outboxCleanupCache{
		watermark:   1,
		listBuckets: []SchedulerBucket{{Platform: PlatformOpenAI, Mode: SchedulerModeSingleplaceholderplaceholder,
placeholder
	cache := &blockingOutboxCleanupCache{
		outboxCleanupCache: baseCache,
		started:            make(chan struct{placeholder),
		release:            make(chan struct{placeholder),
placeholder
	repo := &outboxCleanupRepo{
		events: []SchedulerOutboxEvent{{ID: 1, CreatedAt: time.Now().Add(-time.Hour)placeholderplaceholder,
		rows:   []int64{1placeholder,
placeholder
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			Scheduling: config.GatewaySchedulingConfig{
				OutboxLagRebuildSeconds:  1,
				OutboxLagRebuildFailures: 1,
		placeholder,
	placeholder,
placeholder
	svc := NewSchedulerSnapshotService(cache, repo, &outboxCleanupAccountRepo{placeholder, nil, cfg)

	firstDone := make(chan struct{placeholder)
	go func() {
		svc.checkOutboxLag(context.Background(), 0)
		close(firstDone)
placeholder()
	select {
	case <-cache.started:
	case <-time.After(time.Second):
		t.Fatal("first rebuild did not start")
placeholder

	// The empty batch proves recovery for episode/retry state, but it must not
	// release ownership of the still-running rebuild.
	svc.pollOutbox()

	secondDone := make(chan struct{placeholder)
	go func() {
		svc.checkOutboxLag(context.Background(), 0)
		close(secondDone)
placeholder()
	select {
	case <-secondDone:
	case <-time.After(200 * time.Millisecond):
		close(cache.release)
		<-firstDone
		<-secondDone
		t.Fatal("second lag check queued another rebuild while the first was running")
placeholder

	close(cache.release)
	<-firstDone
	if calls := cache.listCalls(); calls != 1 {
		t.Fatalf("expected one rebuild generation, got %d", calls)
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
