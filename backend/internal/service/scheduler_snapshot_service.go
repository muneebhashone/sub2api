package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
)

var (
	ErrSchedulerCacheNotReady           = errors.New("scheduler cache not ready")
	ErrSchedulerFallbackLimited         = errors.New("scheduler db fallback limited")
	ErrSchedulerGroupLifecycleLeaseBusy = errors.New("scheduler group lifecycle lease busy")
	ErrSchedulerBucketRebuildBusy       = errors.New("scheduler bucket rebuild busy")
)

const (
	outboxEventTimeout                    = 2 * time.Minute
	schedulerOutboxCleanupBatch           = 5000
	schedulerGroupLifecycleTimeout        = 30 * time.Second
	schedulerGroupLifecycleLeaseTTL       = 60 * time.Second
	schedulerGroupLifecycleReleaseTimeout = 2 * time.Second
)

// batchSeenKey tracks completed per-platform rebuilds and group lifecycle work
// within one pollOutbox call.
type batchSeenKey struct {
	groupID   int64
	platform  string
	lifecycle bool
placeholder

type schedulerBucketWriteTask struct {
	bucket SchedulerBucket
	token  SchedulerBucketWriteToken
placeholder

type schedulerAccountQueryKey struct {
	groupID  int64
	platform string
placeholder

type schedulerAccountQueryCache struct {
	remaining map[schedulerAccountQueryKey]int
	accounts  map[schedulerAccountQueryKey][]Account
placeholder

func newSchedulerAccountQueryCache(taskSets ...[]schedulerBucketWriteTask) *schedulerAccountQueryCache {
	queries := &schedulerAccountQueryCache{
		remaining: make(map[schedulerAccountQueryKey]int),
		accounts:  make(map[schedulerAccountQueryKey][]Account),
placeholder
	for _, tasks := range taskSets {
		for _, task := range tasks {
			if key, ok := schedulerAccountQueryKeyForBucket(task.bucket); ok {
				queries.remaining[key]++
		placeholder
	placeholder
placeholder
	return queries
placeholder

func schedulerAccountQueryKeyForBucket(bucket SchedulerBucket) (schedulerAccountQueryKey, bool) {
	if bucket.Mode != SchedulerModeSingle && bucket.Mode != SchedulerModeForced {
		return schedulerAccountQueryKey{placeholder, false
placeholder
	return schedulerAccountQueryKey{groupID: bucket.GroupID, platform: bucket.Platformplaceholder, true
placeholder

func (c *schedulerAccountQueryCache) release(bucket SchedulerBucket) {
	if c == nil {
		return
placeholder
	key, ok := schedulerAccountQueryKeyForBucket(bucket)
	if !ok {
		return
placeholder
	remaining := c.remaining[key] - 1
	if remaining <= 0 {
		delete(c.remaining, key)
		delete(c.accounts, key)
		return
placeholder
	c.remaining[key] = remaining
placeholder

type schedulerGroupLifecyclePlan struct {
	active bool
	tasks  []schedulerBucketWriteTask
placeholder

type schedulerActiveGroupIDLister interface {
	ListActiveIDs(ctx context.Context) ([]int64, error)
placeholder

type SchedulerSnapshotService struct {
	cache         SchedulerCache
	outboxRepo    SchedulerOutboxRepository
	accountRepo   AccountRepository
	groupRepo     GroupRepository
	cfg           *config.Config
	stopCh        chan struct{placeholder
	stopOnce      sync.Once
	wg            sync.WaitGroup
	fallbackLimit *fallbackLimiter
	lagMu         sync.Mutex
	lagFailures   int

	fullRebuildRunMu     sync.Mutex
	fullRebuildStateMu   sync.Mutex
	fullRebuildRequested uint64
	fullRebuildCompleted uint64
	fullRebuildLastErr   error
placeholder

func NewSchedulerSnapshotService(
	cache SchedulerCache,
	outboxRepo SchedulerOutboxRepository,
	accountRepo AccountRepository,
	groupRepo GroupRepository,
	cfg *config.Config,
) *SchedulerSnapshotService {
	maxQPS := 0
	if cfg != nil {
		maxQPS = cfg.Gateway.Scheduling.DbFallbackMaxQPS
placeholder
	return &SchedulerSnapshotService{
		cache:         cache,
		outboxRepo:    outboxRepo,
		accountRepo:   accountRepo,
		groupRepo:     groupRepo,
		cfg:           cfg,
		stopCh:        make(chan struct{placeholder),
		fallbackLimit: newFallbackLimiter(maxQPS),
placeholder
placeholder

func (s *SchedulerSnapshotService) Start() {
	if s == nil || s.cache == nil {
		return
placeholder

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.runInitialRebuild()
placeholder()

	interval := s.outboxPollInterval()
	if s.outboxRepo != nil && interval > 0 {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.runOutboxWorker(interval)
	placeholder()
placeholder

	fullInterval := s.fullRebuildInterval()
	if fullInterval > 0 {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.runFullRebuildWorker(fullInterval)
	placeholder()
placeholder
placeholder

func (s *SchedulerSnapshotService) Stop() {
	if s == nil {
		return
placeholder
	s.stopOnce.Do(func() {
		close(s.stopCh)
placeholder)
	s.wg.Wait()
placeholder

func (s *SchedulerSnapshotService) ListSchedulableAccounts(ctx context.Context, groupID *int64, platform string, hasForcePlatform bool) ([]Account, bool, error) {
	useMixed := (platform == PlatformAnthropic || platform == PlatformGemini) && !hasForcePlatform
	mode := s.resolveMode(platform, hasForcePlatform)
	bucket := s.bucketFor(groupID, platform, mode)
	var writeToken SchedulerBucketWriteToken
	canPublish := false

	if s.cache != nil {
		cached, hit, err := s.cache.GetSnapshot(ctx, bucket)
		if err != nil {
			logger.LegacyPrintf("service.scheduler_snapshot", "[Scheduler] cache read failed: bucket=%s err=%v", bucket.String(), err)
	placeholder else if hit {
			return derefAccounts(cached), useMixed, nil
	placeholder
		token, err := s.cache.CaptureBucketWriteToken(ctx, bucket)
		if err != nil {
			if errors.Is(err, ErrSchedulerBucketRetired) || errors.Is(err, ErrSchedulerBucketWriteFenced) {
				slog.Debug("[Scheduler] cache publish fenced", "bucket", bucket.String())
		placeholder else {
				logger.LegacyPrintf("service.scheduler_snapshot", "[Scheduler] cache publish token failed: bucket=%s err=%v", bucket.String(), err)
		placeholder
	placeholder else {
			writeToken = token
			canPublish = true
	placeholder
placeholder

	if err := s.guardFallback(ctx); err != nil {
		return nil, useMixed, err
placeholder

	fallbackCtx, cancel := s.withFallbackTimeout(ctx)
	defer cancel()

	accounts, err := s.loadAccountsFromDB(fallbackCtx, bucket, useMixed)
	if err != nil {
		return nil, useMixed, err
placeholder

	if s.cache != nil && canPublish {
		if err := s.cache.SetSnapshot(fallbackCtx, bucket, writeToken, accounts); err != nil {
			if errors.Is(err, ErrSchedulerBucketRetired) || errors.Is(err, ErrSchedulerBucketWriteFenced) {
				slog.Debug("[Scheduler] cache publish fenced", "bucket", bucket.String())
		placeholder else {
				logger.LegacyPrintf("service.scheduler_snapshot", "[Scheduler] cache write failed: bucket=%s err=%v", bucket.String(), err)
		placeholder
	placeholder
placeholder

	return accounts, useMixed, nil
placeholder

func (s *SchedulerSnapshotService) GetAccount(ctx context.Context, accountID int64) (*Account, error) {
	if accountID <= 0 {
		return nil, nil
placeholder
	if s.cache != nil {
		account, err := s.cache.GetAccount(ctx, accountID)
		if err != nil {
			logger.LegacyPrintf("service.scheduler_snapshot", "[Scheduler] account cache read failed: id=%d err=%v", accountID, err)
	placeholder else if account != nil {
			return account, nil
	placeholder
placeholder

	if err := s.guardFallback(ctx); err != nil {
		return nil, err
placeholder
	fallbackCtx, cancel := s.withFallbackTimeout(ctx)
	defer cancel()
	return s.accountRepo.GetByID(fallbackCtx, accountID)
placeholder

// GetGroupByID 获取分组信息（供调度器使用）
func (s *SchedulerSnapshotService) GetGroupByID(ctx context.Context, groupID int64) (*Group, error) {
	if s.groupRepo == nil {
		return nil, nil
placeholder
	return s.groupRepo.GetByID(ctx, groupID)
placeholder

// UpdateAccountInCache 立即更新 Redis 中单个账号的数据（用于模型限流后立即生效）
func (s *SchedulerSnapshotService) UpdateAccountInCache(ctx context.Context, account *Account) error {
	if s.cache == nil || account == nil {
		return nil
placeholder
	return s.cache.SetAccount(ctx, account)
placeholder

func (s *SchedulerSnapshotService) runInitialRebuild() {
	if s.cache == nil {
		return
placeholder
	_ = s.coalesceFullRebuild(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		if err := s.rebuildFullSnapshot(ctx, "startup"); err != nil {
			logger.LegacyPrintf("service.scheduler_snapshot", "[Scheduler] rebuild startup failed: %v", err)
			return err
	placeholder
		return nil
placeholder)
placeholder

func (s *SchedulerSnapshotService) runOutboxWorker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	s.pollOutbox()
	for {
		select {
		case <-ticker.C:
			s.pollOutbox()
		case <-s.stopCh:
			return
	placeholder
placeholder
placeholder

func (s *SchedulerSnapshotService) runFullRebuildWorker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.triggerFullRebuild("interval"); err != nil {
				logger.LegacyPrintf("service.scheduler_snapshot", "[Scheduler] full rebuild failed: %v", err)
		placeholder
		case <-s.stopCh:
			return
	placeholder
placeholder
placeholder

func (s *SchedulerSnapshotService) pollOutbox() {
	if s.outboxRepo == nil || s.cache == nil {
		return
placeholder
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	watermark, err := s.cache.GetOutboxWatermark(ctx)
	if err != nil {
		logger.LegacyPrintf("service.scheduler_snapshot", "[Scheduler] outbox watermark read failed: %v", err)
		return
placeholder

	events, err := s.outboxRepo.ListAfterAndReleaseDedup(ctx, watermark, 200)
	if err != nil {
		logger.LegacyPrintf("service.scheduler_snapshot", "[Scheduler] outbox poll failed: %v", err)
		return
placeholder
	if len(events) == 0 {
		return
placeholder

	seen := make(map[batchSeenKey]struct{placeholder)
	for _, event := range events {
		eventCtx, cancel := context.WithTimeout(context.Background(), outboxEventTimeout)
		err := s.handleOutboxEvent(eventCtx, event, seen)
		cancel()
		if err != nil {
			logger.LegacyPrintf("service.scheduler_snapshot", "[Scheduler] outbox handle failed: id=%d type=%s err=%v", event.ID, event.EventType, err)
			return
	placeholder
placeholder

	lastID := events[len(events)-1].ID
	var wmErr error
	for i := range 3 {
		wmCtx, wmCancel := context.WithTimeout(context.Background(), 5*time.Second)
		wmErr = s.cache.SetOutboxWatermark(wmCtx, lastID)
		wmCancel()
		if wmErr == nil {
			break
	placeholder
		if i < 2 {
			time.Sleep(200 * time.Millisecond)
	placeholder
placeholder
	if wmErr != nil {
		logger.LegacyPrintf("service.scheduler_snapshot", "[Scheduler] outbox watermark write failed: %v", wmErr)
		return
placeholder
	s.cleanupConsumedOutbox(lastID)

	// 只有 watermark 成功推进后，当前批次才算已消费。延迟必须按下一条待消费事件计算，
	// 否则本批次处理越慢，越容易误触发一次更慢的全量重建，形成正反馈。
	lagCtx, lagCancel := context.WithTimeout(context.Background(), 5*time.Second)
	s.checkOutboxLag(lagCtx, lastID)
	lagCancel()
placeholder

func (s *SchedulerSnapshotService) cleanupConsumedOutbox(watermark int64) {
	if s == nil || s.outboxRepo == nil || watermark <= 0 {
		return
placeholder

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	lease, acquired, err := s.outboxRepo.TryAcquireCleanupLock(ctx)
	if err != nil {
		logger.LegacyPrintf("service.scheduler_snapshot", "[Scheduler] outbox cleanup lock failed: %v", err)
		return
placeholder
	if !acquired {
		return
placeholder
	defer lease.Release()

	for {
		deleted, err := s.outboxRepo.DeleteConsumedUpTo(ctx, watermark, schedulerOutboxCleanupBatch)
		if err != nil {
			logger.LegacyPrintf("service.scheduler_snapshot", "[Scheduler] outbox cleanup failed: watermark=%d err=%v", watermark, err)
			return
	placeholder
		if deleted == 0 || deleted < schedulerOutboxCleanupBatch {
			return
	placeholder
placeholder
placeholder

func (s *SchedulerSnapshotService) handleOutboxEvent(ctx context.Context, event SchedulerOutboxEvent, seen map[batchSeenKey]struct{placeholder) error {
	switch event.EventType {
	case SchedulerOutboxEventAccountLastUsed:
		return s.handleLastUsedEvent(ctx, event.Payload)
	case SchedulerOutboxEventAccountBulkChanged:
		return s.handleBulkAccountEvent(ctx, event.Payload, seen)
	case SchedulerOutboxEventAccountGroupsChanged:
		return s.handleAccountEvent(ctx, event.AccountID, event.Payload, seen)
	case SchedulerOutboxEventAccountChanged:
		return s.handleAccountEvent(ctx, event.AccountID, event.Payload, seen)
	case SchedulerOutboxEventGroupChanged:
		return s.handleGroupEvent(ctx, event.GroupID, seen)
	case SchedulerOutboxEventFullRebuild:
		return s.triggerFullRebuild("outbox")
	default:
		return nil
placeholder
placeholder

func (s *SchedulerSnapshotService) handleLastUsedEvent(ctx context.Context, payload map[string]any) error {
	if s.cache == nil || payload == nil {
		return nil
placeholder
	raw, ok := payload["last_used"].(map[string]any)
	if !ok || len(raw) == 0 {
		return nil
placeholder
	updates := make(map[int64]time.Time, len(raw))
	for key, value := range raw {
		id, err := strconv.ParseInt(key, 10, 64)
		if err != nil || id <= 0 {
			continue
	placeholder
		sec, ok := toInt64(value)
		if !ok || sec <= 0 {
			continue
	placeholder
		updates[id] = time.Unix(sec, 0)
placeholder
	if len(updates) == 0 {
		return nil
placeholder
	return s.cache.UpdateLastUsed(ctx, updates)
placeholder

func (s *SchedulerSnapshotService) handleBulkAccountEvent(ctx context.Context, payload map[string]any, seen map[batchSeenKey]struct{placeholder) error {
	if payload == nil {
		return nil
placeholder
	if s.accountRepo == nil {
		return nil
placeholder

	rawIDs := parseInt64Slice(payload["account_ids"])
	if len(rawIDs) == 0 {
		return nil
placeholder

	ids := make([]int64, 0, len(rawIDs))
	seenIDs := make(map[int64]struct{placeholder, len(rawIDs))
	for _, id := range rawIDs {
		if id <= 0 {
			continue
	placeholder
		if _, exists := seenIDs[id]; exists {
			continue
	placeholder
		seenIDs[id] = struct{placeholder{placeholder
		ids = append(ids, id)
placeholder
	if len(ids) == 0 {
		return nil
placeholder

	preloadGroupIDs := parseInt64Slice(payload["group_ids"])
	accounts, err := s.accountRepo.GetByIDs(ctx, ids)
	if err != nil {
		return err
placeholder

	found := make(map[int64]struct{placeholder, len(accounts))
	rebuildGroupSet := make(map[int64]struct{placeholder, len(preloadGroupIDs))
	for _, gid := range preloadGroupIDs {
		if gid > 0 {
			rebuildGroupSet[gid] = struct{placeholder{placeholder
	placeholder
placeholder

	for _, account := range accounts {
		if account == nil || account.ID <= 0 {
			continue
	placeholder
		found[account.ID] = struct{placeholder{placeholder
		if s.cache != nil {
			if err := s.cache.SetAccount(ctx, account); err != nil {
				return err
		placeholder
	placeholder
		for _, gid := range account.GroupIDs {
			if gid > 0 {
				rebuildGroupSet[gid] = struct{placeholder{placeholder
		placeholder
	placeholder
placeholder

	if s.cache != nil {
		for _, id := range ids {
			if _, ok := found[id]; ok {
				continue
		placeholder
			if err := s.cache.DeleteAccount(ctx, id); err != nil {
				return err
		placeholder
	placeholder
placeholder

	rebuildGroupIDs := make([]int64, 0, len(rebuildGroupSet))
	for gid := range rebuildGroupSet {
		rebuildGroupIDs = append(rebuildGroupIDs, gid)
placeholder
	return s.rebuildByGroupIDs(ctx, rebuildGroupIDs, "account_bulk_change", seen)
placeholder

func (s *SchedulerSnapshotService) handleAccountEvent(ctx context.Context, accountID *int64, payload map[string]any, seen map[batchSeenKey]struct{placeholder) error {
	if accountID == nil || *accountID <= 0 {
		return nil
placeholder
	if s.accountRepo == nil {
		return nil
placeholder

	var groupIDs []int64
	if payload != nil {
		groupIDs = parseInt64Slice(payload["group_ids"])
placeholder

	account, err := s.accountRepo.GetByID(ctx, *accountID)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			if s.cache != nil {
				if err := s.cache.DeleteAccount(ctx, *accountID); err != nil {
					return err
			placeholder
		placeholder
			return s.rebuildByGroupIDs(ctx, groupIDs, "account_miss", seen)
	placeholder
		return err
placeholder
	if s.cache != nil {
		if err := s.cache.SetAccount(ctx, account); err != nil {
			return err
	placeholder
placeholder
	if len(groupIDs) == 0 {
		groupIDs = account.GroupIDs
placeholder
	return s.rebuildByAccount(ctx, account, groupIDs, "account_change", seen)
placeholder

func (s *SchedulerSnapshotService) handleGroupEvent(ctx context.Context, groupID *int64, seen map[batchSeenKey]struct{placeholder) error {
	if groupID == nil || *groupID <= 0 || s.isRunModeSimple() {
		return nil
placeholder
	if seen != nil {
		if _, ok := seen[batchSeenKey{groupID: *groupID, lifecycle: trueplaceholder]; ok {
			return nil
	placeholder
placeholder
	return s.reconcileGroupLifecycle(ctx, *groupID, seen)
placeholder

func (s *SchedulerSnapshotService) reconcileGroupLifecycle(ctx context.Context, groupID int64, seen map[batchSeenKey]struct{placeholder) error {
	plan, err := s.prepareGroupLifecycle(ctx, groupID, nil)
	if err != nil {
		return err
placeholder
	if plan.active {
		queries := newSchedulerAccountQueryCache(plan.tasks)
		for _, task := range plan.tasks {
			if err := s.rebuildBucketWithTokenPolicyAndQueryCache(ctx, task, "group_change", true, queries); err != nil {
				return err
		placeholder
	placeholder
placeholder
	markGroupLifecycleSeen(seen, groupID)
	return nil
placeholder

func (s *SchedulerSnapshotService) prepareGroupLifecycle(ctx context.Context, groupID int64, knownHistorical []SchedulerBucket) (plan schedulerGroupLifecyclePlan, retErr error) {
	if groupID <= 0 || s.isRunModeSimple() {
		return schedulerGroupLifecyclePlan{placeholder, nil
placeholder
	if s.cache == nil || s.groupRepo == nil {
		return schedulerGroupLifecyclePlan{placeholder, ErrSchedulerCacheNotReady
placeholder

	lifecycleCtx, cancel := context.WithTimeout(ctx, schedulerGroupLifecycleTimeout)
	defer cancel()
	lease, acquired, err := s.cache.TryAcquireGroupLifecycleLease(lifecycleCtx, groupID, schedulerGroupLifecycleLeaseTTL)
	if err != nil {
		return schedulerGroupLifecyclePlan{placeholder, err
placeholder
	if !acquired {
		return schedulerGroupLifecyclePlan{placeholder, fmt.Errorf("%w: group=%d", ErrSchedulerGroupLifecycleLeaseBusy, groupID)
placeholder
	leaseHeld := true
	defer func() {
		if leaseHeld {
			retErr = errors.Join(retErr, s.releaseGroupLifecycleLease(lease))
	placeholder
placeholder()

	group, err := s.groupRepo.GetByIDLite(lifecycleCtx, groupID)
	missing := errors.Is(err, ErrGroupNotFound)
	if err != nil && !missing {
		return schedulerGroupLifecyclePlan{placeholder, err
placeholder
	if err == nil && (group == nil || group.ID != groupID || !group.Hydrated) {
		return schedulerGroupLifecyclePlan{placeholder, fmt.Errorf("untrusted scheduler group lifecycle state: group=%d", groupID)
placeholder

	plan = schedulerGroupLifecyclePlan{active: !missing && group.IsActive()placeholder
	if plan.active {
		buckets := schedulerBucketsForGroup(groupID)
		plan.tasks = make([]schedulerBucketWriteTask, 0, len(buckets))
		for _, bucket := range buckets {
			token, err := s.cache.ReopenBucket(lifecycleCtx, bucket)
			if err != nil {
				return schedulerGroupLifecyclePlan{placeholder, err
		placeholder
			plan.tasks = append(plan.tasks, schedulerBucketWriteTask{bucket: bucket, token: tokenplaceholder)
	placeholder
placeholder else {
		registered := knownHistorical
		if registered == nil {
			registered, err = s.cache.ListBuckets(lifecycleCtx)
			if err != nil {
				return schedulerGroupLifecyclePlan{placeholder, err
		placeholder
	placeholder
		buckets := schedulerBucketsForGroup(groupID)
		for _, bucket := range registered {
			if bucket.GroupID == groupID {
				buckets = append(buckets, bucket)
		placeholder
	placeholder
		for _, bucket := range dedupeBuckets(buckets) {
			if err := s.cache.RetireBucket(lifecycleCtx, bucket); err != nil {
				return schedulerGroupLifecyclePlan{placeholder, err
		placeholder
	placeholder
placeholder

	releaseErr := s.releaseGroupLifecycleLease(lease)
	leaseHeld = false
	if releaseErr != nil {
		return schedulerGroupLifecyclePlan{placeholder, releaseErr
placeholder
	return plan, nil
placeholder

func (s *SchedulerSnapshotService) releaseGroupLifecycleLease(lease SchedulerGroupLifecycleLease) error {
	releaseCtx, cancel := context.WithTimeout(context.Background(), schedulerGroupLifecycleReleaseTimeout)
	defer cancel()
	return s.cache.ReleaseGroupLifecycleLease(releaseCtx, lease)
placeholder

func markGroupLifecycleSeen(seen map[batchSeenKey]struct{placeholder, groupID int64) {
	if seen == nil {
		return
placeholder
	seen[batchSeenKey{groupID: groupID, lifecycle: trueplaceholder] = struct{placeholder{placeholder
	for _, platform := range schedulerSnapshotPlatforms() {
		seen[batchSeenKey{groupID: groupID, platform: platformplaceholder] = struct{placeholder{placeholder
placeholder
placeholder

func (s *SchedulerSnapshotService) rebuildByAccount(ctx context.Context, account *Account, groupIDs []int64, reason string, seen map[batchSeenKey]struct{placeholder) error {
	if account == nil {
		return nil
placeholder
	groupIDs = s.normalizeGroupIDs(groupIDs)
	if len(groupIDs) == 0 {
		return nil
placeholder

	buckets := s.bucketsForPlatform(account.Platform, groupIDs, seen)
	if account.Platform == PlatformAntigravity && account.IsMixedSchedulingEnabled() {
		buckets = append(buckets, s.bucketsForPlatform(PlatformAnthropic, groupIDs, seen)...)
		buckets = append(buckets, s.bucketsForPlatform(PlatformGemini, groupIDs, seen)...)
placeholder
	return s.rebuildBuckets(ctx, buckets, reason)
placeholder

func schedulerSnapshotPlatforms() [5]string {
	return [5]string{PlatformAnthropic, PlatformGemini, PlatformOpenAI, PlatformAntigravity, PlatformGrokplaceholder
placeholder

func schedulerBucketsForGroup(groupID int64) []SchedulerBucket {
	if groupID <= 0 {
		return nil
placeholder
	return schedulerCanonicalBuckets(groupID)
placeholder

func schedulerCanonicalBuckets(groupID int64) []SchedulerBucket {
	buckets := make([]SchedulerBucket, 0, 12)
	for _, platform := range schedulerSnapshotPlatforms() {
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

func (s *SchedulerSnapshotService) rebuildByGroupIDs(ctx context.Context, groupIDs []int64, reason string, seen map[batchSeenKey]struct{placeholder) error {
	groupIDs = s.normalizeGroupIDs(groupIDs)
	if len(groupIDs) == 0 {
		return nil
placeholder
	buckets := make([]SchedulerBucket, 0, len(groupIDs)*12)
	for _, platform := range schedulerSnapshotPlatforms() {
		buckets = append(buckets, s.bucketsForPlatform(platform, groupIDs, seen)...)
placeholder
	return s.rebuildBuckets(ctx, buckets, reason)
placeholder

func (s *SchedulerSnapshotService) bucketsForPlatform(platform string, groupIDs []int64, seen map[batchSeenKey]struct{placeholder) []SchedulerBucket {
	if platform == "" {
		return nil
placeholder
	buckets := make([]SchedulerBucket, 0, len(groupIDs)*3)
	for _, gid := range groupIDs {
		// Within a single poll batch, skip (groupID, platform) pairs that were
		// already rebuilt. The first rebuild loads fresh DB data for all accounts
		// in the group, so subsequent rebuilds for the same group+platform within
		// the same batch are redundant.
		if seen != nil {
			key := batchSeenKey{groupID: gid, platform: platformplaceholder
			if _, exists := seen[key]; exists {
				continue
		placeholder
			seen[key] = struct{placeholder{placeholder
	placeholder
		buckets = append(buckets, SchedulerBucket{GroupID: gid, Platform: platform, Mode: SchedulerModeSingleplaceholder)
		buckets = append(buckets, SchedulerBucket{GroupID: gid, Platform: platform, Mode: SchedulerModeForcedplaceholder)
		if platform == PlatformAnthropic || platform == PlatformGemini {
			buckets = append(buckets, SchedulerBucket{GroupID: gid, Platform: platform, Mode: SchedulerModeMixedplaceholder)
	placeholder
placeholder
	return buckets
placeholder

func (s *SchedulerSnapshotService) rebuildBuckets(ctx context.Context, buckets []SchedulerBucket, reason string) error {
	tasks, firstErr := s.prepareBucketWriteTasks(ctx, buckets)
	queries := newSchedulerAccountQueryCache(tasks)
	if err := s.rebuildPreparedBucketTasks(ctx, tasks, reason, false, queries); err != nil && firstErr == nil {
		firstErr = err
placeholder
	return firstErr
placeholder

func (s *SchedulerSnapshotService) prepareBucketWriteTasks(ctx context.Context, buckets []SchedulerBucket) ([]schedulerBucketWriteTask, error) {
	if s.cache == nil {
		return nil, ErrSchedulerCacheNotReady
placeholder
	tasks := make([]schedulerBucketWriteTask, 0, len(buckets))
	var firstErr error
	for _, bucket := range buckets {
		token, err := s.cache.CaptureBucketWriteToken(ctx, bucket)
		if err != nil {
			if errors.Is(err, ErrSchedulerBucketRetired) || errors.Is(err, ErrSchedulerBucketWriteFenced) {
				continue
		placeholder
			if firstErr == nil {
				firstErr = err
		placeholder
			continue
	placeholder
		tasks = append(tasks, schedulerBucketWriteTask{bucket: bucket, token: tokenplaceholder)
placeholder
	return tasks, firstErr
placeholder

func (s *SchedulerSnapshotService) rebuildPreparedBucketTasks(
	ctx context.Context,
	tasks []schedulerBucketWriteTask,
	reason string,
	strict bool,
	queries *schedulerAccountQueryCache,
) error {
	var firstErr error
	for _, task := range tasks {
		if err := s.rebuildBucketWithTokenPolicyAndQueryCache(ctx, task, reason, strict, queries); err != nil && firstErr == nil {
			firstErr = err
	placeholder
placeholder
	return firstErr
placeholder

func (s *SchedulerSnapshotService) rebuildBucketWithTokenPolicyAndQueryCache(
	ctx context.Context,
	task schedulerBucketWriteTask,
	reason string,
	strict bool,
	queries *schedulerAccountQueryCache,
) error {
	if queries != nil {
		defer queries.release(task.bucket)
placeholder
	if s.cache == nil {
		return ErrSchedulerCacheNotReady
placeholder
	bucket := task.bucket
	ok, err := s.cache.TryLockBucket(ctx, bucket, 30*time.Second)
	if err != nil {
		return err
placeholder
	if !ok {
		if strict {
			return fmt.Errorf("%w: bucket=%s", ErrSchedulerBucketRebuildBusy, bucket.String())
	placeholder
		return nil
placeholder
	defer func() {
		_ = s.cache.UnlockBucket(ctx, bucket)
placeholder()

	rebuildCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	accounts, err := s.loadAccountsForRebuild(rebuildCtx, bucket, queries)
	if err != nil {
		logger.LegacyPrintf("service.scheduler_snapshot", "[Scheduler] rebuild failed: bucket=%s reason=%s err=%v", bucket.String(), reason, err)
		return err
placeholder
	if err := s.cache.SetSnapshot(rebuildCtx, bucket, task.token, accounts); err != nil {
		if errors.Is(err, ErrSchedulerBucketRetired) || errors.Is(err, ErrSchedulerBucketWriteFenced) {
			slog.Debug("[Scheduler] rebuild fenced", "bucket", bucket.String(), "reason", reason)
			if strict {
				return err
		placeholder
			return nil
	placeholder
		logger.LegacyPrintf("service.scheduler_snapshot", "[Scheduler] rebuild cache failed: bucket=%s reason=%s err=%v", bucket.String(), reason, err)
		return err
placeholder
	slog.Debug("[Scheduler] rebuild ok", "bucket", bucket.String(), "reason", reason, "size", len(accounts))
	return nil
placeholder

func (s *SchedulerSnapshotService) triggerFullRebuild(reason string) error {
	if s.cache == nil {
		return ErrSchedulerCacheNotReady
placeholder
	return s.coalesceFullRebuild(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		return s.rebuildFullSnapshot(ctx, reason)
placeholder)
placeholder

func (s *SchedulerSnapshotService) rebuildFullSnapshot(ctx context.Context, reason string) error {
	if s.cache == nil {
		return ErrSchedulerCacheNotReady
placeholder

	registered, err := s.cache.ListBuckets(ctx)
	if err != nil {
		return err
placeholder
	registered = dedupeBuckets(registered)

	if s.isRunModeSimple() {
		canonical := schedulerCanonicalBuckets(0)
		captured, err := s.captureFullRebuildCanonicalTasks(ctx, canonical)
		if err != nil {
			return err
	placeholder
		ordinary := appendBucketsExcept(nil, registered, canonical)
		return s.prepareAndRebuildFullSnapshot(ctx, captured, nil, ordinary, reason)
placeholder

	activeGroupIDs, err := s.listActiveSchedulerGroupIDs(ctx)
	if err != nil {
		return err
placeholder
	activeGroups := make(map[int64]struct{placeholder, len(activeGroupIDs))
	for _, groupID := range activeGroupIDs {
		activeGroups[groupID] = struct{placeholder{placeholder
placeholder

	registeredByGroup := make(map[int64][]SchedulerBucket)
	for _, bucket := range registered {
		registeredByGroup[bucket.GroupID] = append(registeredByGroup[bucket.GroupID], bucket)
placeholder

	groupZeroCanonical := schedulerCanonicalBuckets(0)
	capturedTasks, err := s.captureFullRebuildCanonicalTasks(ctx, groupZeroCanonical)
	if err != nil {
		return err
placeholder
	ordinaryBuckets := appendBucketsExcept(nil, registeredByGroup[0], groupZeroCanonical)
	for groupID, buckets := range registeredByGroup {
		if groupID < 0 {
			ordinaryBuckets = append(ordinaryBuckets, buckets...)
	placeholder
placeholder

	reopenedTasks := make([]schedulerBucketWriteTask, 0)
	for _, groupID := range activeGroupIDs {
		canonical := schedulerBucketsForGroup(groupID)
		canonicalTasks, captureErr := s.captureFullRebuildCanonicalTasks(ctx, canonical)
		if captureErr == nil {
			capturedTasks = append(capturedTasks, canonicalTasks...)
			ordinaryBuckets = appendBucketsExcept(ordinaryBuckets, registeredByGroup[groupID], canonical)
			continue
	placeholder
		if !errors.Is(captureErr, ErrSchedulerBucketRetired) && !errors.Is(captureErr, ErrSchedulerBucketWriteFenced) {
			return captureErr
	placeholder

		// A prior full_rebuild event can observe the active state committed for a
		// later group_changed event. Recover here under fresh authority so the
		// earlier event cannot block the outbox watermark before that event runs.
		knownHistorical := registeredByGroup[groupID]
		if knownHistorical == nil {
			knownHistorical = []SchedulerBucket{placeholder
	placeholder
		plan, err := s.prepareGroupLifecycle(ctx, groupID, knownHistorical)
		if err != nil {
			return err
	placeholder
		if plan.active {
			reopenedTasks = append(reopenedTasks, plan.tasks...)
			ordinaryBuckets = appendBucketsExcept(ordinaryBuckets, registeredByGroup[groupID], canonical)
	placeholder
placeholder

	staleGroupIDs := make([]int64, 0)
	for groupID := range registeredByGroup {
		if groupID <= 0 {
			continue
	placeholder
		if _, active := activeGroups[groupID]; !active {
			staleGroupIDs = append(staleGroupIDs, groupID)
	placeholder
placeholder
	sort.Slice(staleGroupIDs, func(i, j int) bool { return staleGroupIDs[i] < staleGroupIDs[j] placeholder)

	for _, groupID := range staleGroupIDs {
		plan, err := s.prepareGroupLifecycle(ctx, groupID, registeredByGroup[groupID])
		if err != nil {
			return err
	placeholder
		if plan.active {
			reopenedTasks = append(reopenedTasks, plan.tasks...)
			ordinaryBuckets = appendBucketsExcept(ordinaryBuckets, registeredByGroup[groupID], schedulerBucketsForGroup(groupID))
	placeholder
placeholder

	return s.prepareAndRebuildFullSnapshot(ctx, capturedTasks, reopenedTasks, ordinaryBuckets, reason)
placeholder

func (s *SchedulerSnapshotService) listActiveSchedulerGroupIDs(ctx context.Context) ([]int64, error) {
	if s.groupRepo == nil {
		return nil, ErrSchedulerCacheNotReady
placeholder

	var groupIDs []int64
	if lister, ok := s.groupRepo.(schedulerActiveGroupIDLister); ok {
		ids, err := lister.ListActiveIDs(ctx)
		if err != nil {
			return nil, err
	placeholder
		groupIDs = ids
placeholder else {
		groups, err := s.groupRepo.ListActive(ctx)
		if err != nil {
			return nil, err
	placeholder
		groupIDs = make([]int64, 0, len(groups))
		for _, group := range groups {
			groupIDs = append(groupIDs, group.ID)
	placeholder
placeholder

	seen := make(map[int64]struct{placeholder, len(groupIDs))
	normalized := make([]int64, 0, len(groupIDs))
	for _, groupID := range groupIDs {
		if groupID <= 0 {
			continue
	placeholder
		if _, ok := seen[groupID]; ok {
			continue
	placeholder
		seen[groupID] = struct{placeholder{placeholder
		normalized = append(normalized, groupID)
placeholder
	sort.Slice(normalized, func(i, j int) bool { return normalized[i] < normalized[j] placeholder)
	return normalized, nil
placeholder

func (s *SchedulerSnapshotService) prepareAndRebuildFullSnapshot(
	ctx context.Context,
	captured []schedulerBucketWriteTask,
	reopened []schedulerBucketWriteTask,
	ordinaryBuckets []SchedulerBucket,
	reason string,
) error {
	preparedBuckets := make(map[SchedulerBucket]struct{placeholder, len(captured)+len(reopened))
	for _, task := range captured {
		preparedBuckets[task.bucket] = struct{placeholder{placeholder
placeholder
	for _, task := range reopened {
		preparedBuckets[task.bucket] = struct{placeholder{placeholder
placeholder

	ordinaryBuckets = dedupeBuckets(ordinaryBuckets)
	toCapture := make([]SchedulerBucket, 0, len(ordinaryBuckets))
	for _, bucket := range ordinaryBuckets {
		if _, ok := preparedBuckets[bucket]; !ok {
			toCapture = append(toCapture, bucket)
	placeholder
placeholder
	ordinary, firstErr := s.prepareBucketWriteTasks(ctx, toCapture)
	if firstErr != nil {
		return firstErr
placeholder
	captured = append(captured, ordinary...)
	queries := newSchedulerAccountQueryCache(reopened, captured)
	if err := s.rebuildPreparedBucketTasks(ctx, reopened, reason, true, queries); err != nil {
		firstErr = err
placeholder
	if err := s.rebuildPreparedBucketTasks(ctx, captured, reason, false, queries); err != nil && firstErr == nil {
		firstErr = err
placeholder
	return firstErr
placeholder

func (s *SchedulerSnapshotService) captureFullRebuildCanonicalTasks(ctx context.Context, buckets []SchedulerBucket) ([]schedulerBucketWriteTask, error) {
	if s.cache == nil {
		return nil, ErrSchedulerCacheNotReady
placeholder
	tasks := make([]schedulerBucketWriteTask, 0, len(buckets))
	for _, bucket := range buckets {
		token, err := s.cache.CaptureBucketWriteToken(ctx, bucket)
		if err != nil {
			return nil, err
	placeholder
		tasks = append(tasks, schedulerBucketWriteTask{bucket: bucket, token: tokenplaceholder)
placeholder
	return tasks, nil
placeholder

func appendBucketsExcept(dst, buckets, excluded []SchedulerBucket) []SchedulerBucket {
	excludedKeys := make(map[SchedulerBucket]struct{placeholder, len(excluded))
	for _, bucket := range excluded {
		excludedKeys[bucket] = struct{placeholder{placeholder
placeholder
	for _, bucket := range buckets {
		if _, ok := excludedKeys[bucket]; !ok {
			dst = append(dst, bucket)
	placeholder
placeholder
	return dst
placeholder

func (s *SchedulerSnapshotService) coalesceFullRebuild(run func() error) error {
	s.fullRebuildStateMu.Lock()
	s.fullRebuildRequested++
	requestID := s.fullRebuildRequested
	s.fullRebuildStateMu.Unlock()

	s.fullRebuildRunMu.Lock()
	defer s.fullRebuildRunMu.Unlock()

	s.fullRebuildStateMu.Lock()
	if s.fullRebuildCompleted >= requestID {
		err := s.fullRebuildLastErr
		s.fullRebuildStateMu.Unlock()
		return err
placeholder
	// 当前轮重建可能早于新 outbox 事件对应事务的提交，不能让后到请求直接复用当前轮。
	// 每轮开始前记录可覆盖的请求代次，执行期间登记的请求统一合并到下一轮。
	coveredThrough := s.fullRebuildRequested
	s.fullRebuildStateMu.Unlock()

	err := run()

	s.fullRebuildStateMu.Lock()
	s.fullRebuildCompleted = coveredThrough
	s.fullRebuildLastErr = err
	s.fullRebuildStateMu.Unlock()
	return err
placeholder

func (s *SchedulerSnapshotService) checkOutboxLag(ctx context.Context, watermark int64) {
	if s.cfg == nil || s.outboxRepo == nil {
		return
placeholder
	oldestCreatedAt, ok, err := s.outboxRepo.FirstCreatedAtAfter(ctx, watermark)
	if err != nil {
		logger.LegacyPrintf("service.scheduler_snapshot", "[Scheduler] outbox pending event read failed: %v", err)
		return
placeholder
	if !ok || oldestCreatedAt.IsZero() {
		s.lagMu.Lock()
		s.lagFailures = 0
		s.lagMu.Unlock()
		return
placeholder

	lag := time.Since(oldestCreatedAt)
	if lagSeconds := int(lag.Seconds()); lagSeconds >= s.cfg.Gateway.Scheduling.OutboxLagWarnSeconds && s.cfg.Gateway.Scheduling.OutboxLagWarnSeconds > 0 {
		logger.LegacyPrintf("service.scheduler_snapshot", "[Scheduler] outbox lag warning: %ds", lagSeconds)
placeholder

	if s.cfg.Gateway.Scheduling.OutboxLagRebuildSeconds > 0 && int(lag.Seconds()) >= s.cfg.Gateway.Scheduling.OutboxLagRebuildSeconds {
		s.lagMu.Lock()
		s.lagFailures++
		failures := s.lagFailures
		s.lagMu.Unlock()

		if failures >= s.cfg.Gateway.Scheduling.OutboxLagRebuildFailures {
			logger.LegacyPrintf("service.scheduler_snapshot", "[Scheduler] outbox lag rebuild triggered: lag=%s failures=%d", lag, failures)
			s.lagMu.Lock()
			s.lagFailures = 0
			s.lagMu.Unlock()
			if err := s.triggerFullRebuild("outbox_lag"); err != nil {
				logger.LegacyPrintf("service.scheduler_snapshot", "[Scheduler] outbox lag rebuild failed: %v", err)
		placeholder
	placeholder
placeholder else {
		s.lagMu.Lock()
		s.lagFailures = 0
		s.lagMu.Unlock()
placeholder

	threshold := s.cfg.Gateway.Scheduling.OutboxBacklogRebuildRows
	if threshold <= 0 {
		return
placeholder
	maxID, err := s.outboxRepo.MaxID(ctx)
	if err != nil {
		return
placeholder
	if maxID-watermark >= int64(threshold) {
		logger.LegacyPrintf("service.scheduler_snapshot", "[Scheduler] outbox backlog rebuild triggered: backlog=%d", maxID-watermark)
		if err := s.triggerFullRebuild("outbox_backlog"); err != nil {
			logger.LegacyPrintf("service.scheduler_snapshot", "[Scheduler] outbox backlog rebuild failed: %v", err)
	placeholder
placeholder
placeholder

func (s *SchedulerSnapshotService) loadAccountsFromDB(ctx context.Context, bucket SchedulerBucket, useMixed bool) ([]Account, error) {
	if s.accountRepo == nil {
		return nil, ErrSchedulerCacheNotReady
placeholder
	groupID := bucket.GroupID
	if s.isRunModeSimple() {
		groupID = 0
placeholder

	if useMixed {
		platforms := []string{bucket.Platform, PlatformAntigravityplaceholder
		var accounts []Account
		var err error
		if groupID > 0 {
			accounts, err = s.accountRepo.ListSchedulableByGroupIDAndPlatforms(ctx, groupID, platforms)
	placeholder else if s.isRunModeSimple() {
			accounts, err = s.accountRepo.ListSchedulableByPlatforms(ctx, platforms)
	placeholder else {
			accounts, err = s.accountRepo.ListSchedulableUngroupedByPlatforms(ctx, platforms)
	placeholder
		if err != nil {
			return nil, err
	placeholder
		filtered := make([]Account, 0, len(accounts))
		for _, acc := range accounts {
			if acc.Platform == PlatformAntigravity && !acc.IsMixedSchedulingEnabled() {
				continue
		placeholder
			filtered = append(filtered, acc)
	placeholder
		return filtered, nil
placeholder

	if groupID > 0 {
		return s.accountRepo.ListSchedulableByGroupIDAndPlatform(ctx, groupID, bucket.Platform)
placeholder
	if s.isRunModeSimple() {
		return s.accountRepo.ListSchedulableByPlatform(ctx, bucket.Platform)
placeholder
	return s.accountRepo.ListSchedulableUngroupedByPlatform(ctx, bucket.Platform)
placeholder

func (s *SchedulerSnapshotService) loadAccountsForRebuild(
	ctx context.Context,
	bucket SchedulerBucket,
	queries *schedulerAccountQueryCache,
) ([]Account, error) {
	key, cacheable := schedulerAccountQueryKeyForBucket(bucket)
	if queries == nil || !cacheable {
		return s.loadAccountsFromDB(ctx, bucket, bucket.Mode == SchedulerModeMixed)
placeholder

	if accounts, ok := queries.accounts[key]; ok {
		return accounts, nil
placeholder
	if queries.remaining[key] <= 1 {
		return s.loadAccountsFromDB(ctx, bucket, false)
placeholder
	accounts, err := s.loadAccountsFromDB(ctx, bucket, false)
	if err != nil {
		return nil, err
placeholder
	queries.accounts[key] = accounts
	return accounts, nil
placeholder

func (s *SchedulerSnapshotService) bucketFor(groupID *int64, platform string, mode string) SchedulerBucket {
	return SchedulerBucket{
		GroupID:  s.normalizeGroupID(groupID),
		Platform: platform,
		Mode:     mode,
placeholder
placeholder

func (s *SchedulerSnapshotService) normalizeGroupID(groupID *int64) int64 {
	if s.isRunModeSimple() {
		return 0
placeholder
	if groupID == nil || *groupID <= 0 {
		return 0
placeholder
	return *groupID
placeholder

func (s *SchedulerSnapshotService) normalizeGroupIDs(groupIDs []int64) []int64 {
	if s.isRunModeSimple() {
		return []int64{0placeholder
placeholder
	if len(groupIDs) == 0 {
		return []int64{0placeholder
placeholder
	seen := make(map[int64]struct{placeholder, len(groupIDs))
	out := make([]int64, 0, len(groupIDs))
	for _, id := range groupIDs {
		if id <= 0 {
			continue
	placeholder
		if _, ok := seen[id]; ok {
			continue
	placeholder
		seen[id] = struct{placeholder{placeholder
		out = append(out, id)
placeholder
	if len(out) == 0 {
		return []int64{0placeholder
placeholder
	return out
placeholder

func (s *SchedulerSnapshotService) resolveMode(platform string, hasForcePlatform bool) string {
	if hasForcePlatform {
		return SchedulerModeForced
placeholder
	if platform == PlatformAnthropic || platform == PlatformGemini {
		return SchedulerModeMixed
placeholder
	return SchedulerModeSingle
placeholder

func (s *SchedulerSnapshotService) guardFallback(ctx context.Context) error {
	if s.cfg == nil || s.cfg.Gateway.Scheduling.DbFallbackEnabled {
		if s.fallbackLimit == nil || s.fallbackLimit.Allow() {
			return nil
	placeholder
		return ErrSchedulerFallbackLimited
placeholder
	return ErrSchedulerCacheNotReady
placeholder

func (s *SchedulerSnapshotService) withFallbackTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if s.cfg == nil || s.cfg.Gateway.Scheduling.DbFallbackTimeoutSeconds <= 0 {
		return context.WithCancel(ctx)
placeholder
	timeout := time.Duration(s.cfg.Gateway.Scheduling.DbFallbackTimeoutSeconds) * time.Second
	if deadline, ok := ctx.Deadline(); ok {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return context.WithCancel(ctx)
	placeholder
		if remaining < timeout {
			timeout = remaining
	placeholder
placeholder
	return context.WithTimeout(ctx, timeout)
placeholder

func (s *SchedulerSnapshotService) isRunModeSimple() bool {
	return s.cfg != nil && s.cfg.RunMode == config.RunModeSimple
placeholder

func (s *SchedulerSnapshotService) outboxPollInterval() time.Duration {
	if s.cfg == nil {
		return time.Second
placeholder
	sec := s.cfg.Gateway.Scheduling.OutboxPollIntervalSeconds
	if sec <= 0 {
		return time.Second
placeholder
	return time.Duration(sec) * time.Second
placeholder

func (s *SchedulerSnapshotService) fullRebuildInterval() time.Duration {
	if s.cfg == nil {
		return 0
placeholder
	sec := s.cfg.Gateway.Scheduling.FullRebuildIntervalSeconds
	if sec <= 0 {
		return 0
placeholder
	return time.Duration(sec) * time.Second
placeholder

func dedupeBuckets(in []SchedulerBucket) []SchedulerBucket {
	seen := make(map[string]struct{placeholder, len(in))
	out := make([]SchedulerBucket, 0, len(in))
	for _, bucket := range in {
		key := bucket.String()
		if _, ok := seen[key]; ok {
			continue
	placeholder
		seen[key] = struct{placeholder{placeholder
		out = append(out, bucket)
placeholder
	return out
placeholder

func derefAccounts(accounts []*Account) []Account {
	if len(accounts) == 0 {
		return []Account{placeholder
placeholder
	out := make([]Account, 0, len(accounts))
	for _, account := range accounts {
		if account == nil {
			continue
	placeholder
		out = append(out, *account)
placeholder
	return out
placeholder

func parseInt64Slice(value any) []int64 {
	raw, ok := value.([]any)
	if !ok {
		return nil
placeholder
	out := make([]int64, 0, len(raw))
	for _, item := range raw {
		if v, ok := toInt64(item); ok && v > 0 {
			out = append(out, v)
	placeholder
placeholder
	return out
placeholder

func toInt64(value any) (int64, bool) {
	switch v := value.(type) {
	case float64:
		return int64(v), true
	case int64:
		return v, true
	case int:
		return int64(v), true
	case json.Number:
		parsed, err := strconv.ParseInt(v.String(), 10, 64)
		return parsed, err == nil
	default:
		return 0, false
placeholder
placeholder

type fallbackLimiter struct {
	maxQPS int
	mu     sync.Mutex
	window time.Time
	count  int
placeholder

func newFallbackLimiter(maxQPS int) *fallbackLimiter {
	if maxQPS <= 0 {
		return nil
placeholder
	return &fallbackLimiter{
		maxQPS: maxQPS,
		window: time.Now(),
placeholder
placeholder

func (l *fallbackLimiter) Allow() bool {
	if l == nil || l.maxQPS <= 0 {
		return true
placeholder
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if now.Sub(l.window) >= time.Second {
		l.window = now
		l.count = 0
placeholder
	if l.count >= l.maxQPS {
		return false
placeholder
	l.count++
	return true
placeholder
