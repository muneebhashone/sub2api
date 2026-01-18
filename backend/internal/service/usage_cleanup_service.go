package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

const (
	usageCleanupWorkerName = "usage_cleanup_worker"
)

// UsageCleanupService 负责创建与执行使用记录清理任务
type UsageCleanupService struct {
	repo        UsageCleanupRepository
	timingWheel *TimingWheelService
	dashboard   *DashboardAggregationService
	cfg         *config.Config

	running   int32
	startOnce sync.Once
	stopOnce  sync.Once

	workerCtx    context.Context
	workerCancel context.CancelFunc
placeholder

func NewUsageCleanupService(repo UsageCleanupRepository, timingWheel *TimingWheelService, dashboard *DashboardAggregationService, cfg *config.Config) *UsageCleanupService {
	workerCtx, workerCancel := context.WithCancel(context.Background())
	return &UsageCleanupService{
		repo:         repo,
		timingWheel:  timingWheel,
		dashboard:    dashboard,
		cfg:          cfg,
		workerCtx:    workerCtx,
		workerCancel: workerCancel,
placeholder
placeholder

func describeUsageCleanupFilters(filters UsageCleanupFilters) string {
	var parts []string
	parts = append(parts, "start="+filters.StartTime.UTC().Format(time.RFC3339))
	parts = append(parts, "end="+filters.EndTime.UTC().Format(time.RFC3339))
	if filters.UserID != nil {
		parts = append(parts, fmt.Sprintf("user_id=%d", *filters.UserID))
placeholder
	if filters.APIKeyID != nil {
		parts = append(parts, fmt.Sprintf("api_key_id=%d", *filters.APIKeyID))
placeholder
	if filters.AccountID != nil {
		parts = append(parts, fmt.Sprintf("account_id=%d", *filters.AccountID))
placeholder
	if filters.GroupID != nil {
		parts = append(parts, fmt.Sprintf("group_id=%d", *filters.GroupID))
placeholder
	if filters.Model != nil {
		parts = append(parts, "model="+strings.TrimSpace(*filters.Model))
placeholder
	if filters.Stream != nil {
		parts = append(parts, fmt.Sprintf("stream=%t", *filters.Stream))
placeholder
	if filters.BillingType != nil {
		parts = append(parts, fmt.Sprintf("billing_type=%d", *filters.BillingType))
placeholder
	return strings.Join(parts, " ")
placeholder

func (s *UsageCleanupService) Start() {
	if s == nil {
		return
placeholder
	if s.cfg != nil && !s.cfg.UsageCleanup.Enabled {
		log.Printf("[UsageCleanup] not started (disabled)")
		return
placeholder
	if s.repo == nil || s.timingWheel == nil {
		log.Printf("[UsageCleanup] not started (missing deps)")
		return
placeholder

	interval := s.workerInterval()
	s.startOnce.Do(func() {
		s.timingWheel.ScheduleRecurring(usageCleanupWorkerName, interval, s.runOnce)
		log.Printf("[UsageCleanup] started (interval=%s max_range_days=%d batch_size=%d task_timeout=%s)", interval, s.maxRangeDays(), s.batchSize(), s.taskTimeout())
placeholder)
placeholder

func (s *UsageCleanupService) Stop() {
	if s == nil {
		return
placeholder
	s.stopOnce.Do(func() {
		if s.workerCancel != nil {
			s.workerCancel()
	placeholder
		if s.timingWheel != nil {
			s.timingWheel.Cancel(usageCleanupWorkerName)
	placeholder
		log.Printf("[UsageCleanup] stopped")
placeholder)
placeholder

func (s *UsageCleanupService) ListTasks(ctx context.Context, params pagination.PaginationParams) ([]UsageCleanupTask, *pagination.PaginationResult, error) {
	if s == nil || s.repo == nil {
		return nil, nil, fmt.Errorf("cleanup service not ready")
placeholder
	return s.repo.ListTasks(ctx, params)
placeholder

func (s *UsageCleanupService) CreateTask(ctx context.Context, filters UsageCleanupFilters, createdBy int64) (*UsageCleanupTask, error) {
	if s == nil || s.repo == nil {
		return nil, fmt.Errorf("cleanup service not ready")
placeholder
	if s.cfg != nil && !s.cfg.UsageCleanup.Enabled {
		return nil, infraerrors.New(http.StatusServiceUnavailable, "USAGE_CLEANUP_DISABLED", "usage cleanup is disabled")
placeholder
	if createdBy <= 0 {
		return nil, infraerrors.BadRequest("USAGE_CLEANUP_INVALID_CREATOR", "invalid creator")
placeholder

	log.Printf("[UsageCleanup] create_task requested: operator=%d %s", createdBy, describeUsageCleanupFilters(filters))
	sanitizeUsageCleanupFilters(&filters)
	if err := s.validateFilters(filters); err != nil {
		log.Printf("[UsageCleanup] create_task rejected: operator=%d err=%v %s", createdBy, err, describeUsageCleanupFilters(filters))
		return nil, err
placeholder

	task := &UsageCleanupTask{
		Status:    UsageCleanupStatusPending,
		Filters:   filters,
		CreatedBy: createdBy,
placeholder
	if err := s.repo.CreateTask(ctx, task); err != nil {
		log.Printf("[UsageCleanup] create_task persist failed: operator=%d err=%v %s", createdBy, err, describeUsageCleanupFilters(filters))
		return nil, fmt.Errorf("create cleanup task: %w", err)
placeholder
	log.Printf("[UsageCleanup] create_task persisted: task=%d operator=%d status=%s deleted_rows=%d %s", task.ID, createdBy, task.Status, task.DeletedRows, describeUsageCleanupFilters(filters))
	go s.runOnce()
	return task, nil
placeholder

func (s *UsageCleanupService) runOnce() {
	svc := s
	if svc == nil {
		return
placeholder
	if !atomic.CompareAndSwapInt32(&svc.running, 0, 1) {
		log.Printf("[UsageCleanup] run_once skipped: already_running=true")
		return
placeholder
	defer atomic.StoreInt32(&svc.running, 0)

	parent := context.Background()
	if svc.workerCtx != nil {
		parent = svc.workerCtx
placeholder
	ctx, cancel := context.WithTimeout(parent, svc.taskTimeout())
	defer cancel()

	task, err := svc.repo.ClaimNextPendingTask(ctx, int64(svc.taskTimeout().Seconds()))
	if err != nil {
		log.Printf("[UsageCleanup] claim pending task failed: %v", err)
		return
placeholder
	if task == nil {
		log.Printf("[UsageCleanup] run_once done: no_task=true")
		return
placeholder

	log.Printf("[UsageCleanup] task claimed: task=%d status=%s created_by=%d deleted_rows=%d %s", task.ID, task.Status, task.CreatedBy, task.DeletedRows, describeUsageCleanupFilters(task.Filters))
	svc.executeTask(ctx, task)
placeholder

func (s *UsageCleanupService) executeTask(ctx context.Context, task *UsageCleanupTask) {
	if task == nil {
		return
placeholder

	batchSize := s.batchSize()
	deletedTotal := task.DeletedRows
	start := time.Now()
	log.Printf("[UsageCleanup] task started: task=%d batch_size=%d deleted_rows=%d %s", task.ID, batchSize, deletedTotal, describeUsageCleanupFilters(task.Filters))
	var batchNum int

	for {
		if ctx != nil && ctx.Err() != nil {
			log.Printf("[UsageCleanup] task interrupted: task=%d err=%v", task.ID, ctx.Err())
			return
	placeholder
		canceled, err := s.isTaskCanceled(ctx, task.ID)
		if err != nil {
			s.markTaskFailed(task.ID, deletedTotal, err)
			return
	placeholder
		if canceled {
			log.Printf("[UsageCleanup] task canceled: task=%d deleted_rows=%d duration=%s", task.ID, deletedTotal, time.Since(start))
			return
	placeholder

		batchNum++
		deleted, err := s.repo.DeleteUsageLogsBatch(ctx, task.Filters, batchSize)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				// 任务被中断（例如服务停止/超时），保持 running 状态，后续通过 stale reclaim 续跑。
				log.Printf("[UsageCleanup] task interrupted: task=%d err=%v", task.ID, err)
				return
		placeholder
			s.markTaskFailed(task.ID, deletedTotal, err)
			return
	placeholder
		deletedTotal += deleted
		if deleted > 0 {
			updateCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			if err := s.repo.UpdateTaskProgress(updateCtx, task.ID, deletedTotal); err != nil {
				log.Printf("[UsageCleanup] task progress update failed: task=%d deleted_rows=%d err=%v", task.ID, deletedTotal, err)
		placeholder
			cancel()
	placeholder
		if batchNum <= 3 || batchNum%20 == 0 || deleted < int64(batchSize) {
			log.Printf("[UsageCleanup] task batch done: task=%d batch=%d deleted=%d deleted_total=%d", task.ID, batchNum, deleted, deletedTotal)
	placeholder
		if deleted == 0 || deleted < int64(batchSize) {
			break
	placeholder
placeholder

	updateCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.repo.MarkTaskSucceeded(updateCtx, task.ID, deletedTotal); err != nil {
		log.Printf("[UsageCleanup] update task succeeded failed: task=%d err=%v", task.ID, err)
placeholder else {
		log.Printf("[UsageCleanup] task succeeded: task=%d deleted_rows=%d duration=%s", task.ID, deletedTotal, time.Since(start))
placeholder

	if s.dashboard != nil {
		if err := s.dashboard.TriggerRecomputeRange(task.Filters.StartTime, task.Filters.EndTime); err != nil {
			log.Printf("[UsageCleanup] trigger dashboard recompute failed: task=%d err=%v", task.ID, err)
	placeholder else {
			log.Printf("[UsageCleanup] trigger dashboard recompute: task=%d start=%s end=%s", task.ID, task.Filters.StartTime.UTC().Format(time.RFC3339), task.Filters.EndTime.UTC().Format(time.RFC3339))
	placeholder
placeholder
placeholder

func (s *UsageCleanupService) markTaskFailed(taskID int64, deletedRows int64, err error) {
	msg := strings.TrimSpace(err.Error())
	if len(msg) > 500 {
		msg = msg[:500]
placeholder
	log.Printf("[UsageCleanup] task failed: task=%d deleted_rows=%d err=%s", taskID, deletedRows, msg)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if updateErr := s.repo.MarkTaskFailed(ctx, taskID, deletedRows, msg); updateErr != nil {
		log.Printf("[UsageCleanup] update task failed failed: task=%d err=%v", taskID, updateErr)
placeholder
placeholder

func (s *UsageCleanupService) isTaskCanceled(ctx context.Context, taskID int64) (bool, error) {
	if s == nil || s.repo == nil {
		return false, fmt.Errorf("cleanup service not ready")
placeholder
	checkCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	status, err := s.repo.GetTaskStatus(checkCtx, taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
	placeholder
		return false, err
placeholder
	if status == UsageCleanupStatusCanceled {
		log.Printf("[UsageCleanup] task cancel detected: task=%d", taskID)
placeholder
	return status == UsageCleanupStatusCanceled, nil
placeholder

func (s *UsageCleanupService) validateFilters(filters UsageCleanupFilters) error {
	if filters.StartTime.IsZero() || filters.EndTime.IsZero() {
		return infraerrors.BadRequest("USAGE_CLEANUP_MISSING_RANGE", "start_date and end_date are required")
placeholder
	if filters.EndTime.Before(filters.StartTime) {
		return infraerrors.BadRequest("USAGE_CLEANUP_INVALID_RANGE", "end_date must be after start_date")
placeholder
	maxDays := s.maxRangeDays()
	if maxDays > 0 {
		delta := filters.EndTime.Sub(filters.StartTime)
		if delta > time.Duration(maxDays)*24*time.Hour {
			return infraerrors.BadRequest("USAGE_CLEANUP_RANGE_TOO_LARGE", fmt.Sprintf("date range exceeds %d days", maxDays))
	placeholder
placeholder
	return nil
placeholder

func (s *UsageCleanupService) CancelTask(ctx context.Context, taskID int64, canceledBy int64) error {
	if s == nil || s.repo == nil {
		return fmt.Errorf("cleanup service not ready")
placeholder
	if s.cfg != nil && !s.cfg.UsageCleanup.Enabled {
		return infraerrors.New(http.StatusServiceUnavailable, "USAGE_CLEANUP_DISABLED", "usage cleanup is disabled")
placeholder
	if canceledBy <= 0 {
		return infraerrors.BadRequest("USAGE_CLEANUP_INVALID_CANCELLER", "invalid canceller")
placeholder
	status, err := s.repo.GetTaskStatus(ctx, taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return infraerrors.New(http.StatusNotFound, "USAGE_CLEANUP_TASK_NOT_FOUND", "cleanup task not found")
	placeholder
		return err
placeholder
	log.Printf("[UsageCleanup] cancel_task requested: task=%d operator=%d status=%s", taskID, canceledBy, status)
	if status != UsageCleanupStatusPending && status != UsageCleanupStatusRunning {
		return infraerrors.New(http.StatusConflict, "USAGE_CLEANUP_CANCEL_CONFLICT", "cleanup task cannot be canceled in current status")
placeholder
	ok, err := s.repo.CancelTask(ctx, taskID, canceledBy)
	if err != nil {
		return err
placeholder
	if !ok {
		// 状态可能并发改变
		return infraerrors.New(http.StatusConflict, "USAGE_CLEANUP_CANCEL_CONFLICT", "cleanup task cannot be canceled in current status")
placeholder
	log.Printf("[UsageCleanup] cancel_task done: task=%d operator=%d", taskID, canceledBy)
	return nil
placeholder

func sanitizeUsageCleanupFilters(filters *UsageCleanupFilters) {
	if filters == nil {
		return
placeholder
	if filters.UserID != nil && *filters.UserID <= 0 {
		filters.UserID = nil
placeholder
	if filters.APIKeyID != nil && *filters.APIKeyID <= 0 {
		filters.APIKeyID = nil
placeholder
	if filters.AccountID != nil && *filters.AccountID <= 0 {
		filters.AccountID = nil
placeholder
	if filters.GroupID != nil && *filters.GroupID <= 0 {
		filters.GroupID = nil
placeholder
	if filters.Model != nil {
		model := strings.TrimSpace(*filters.Model)
		if model == "" {
			filters.Model = nil
	placeholder else {
			filters.Model = &model
	placeholder
placeholder
	if filters.BillingType != nil && *filters.BillingType < 0 {
		filters.BillingType = nil
placeholder
placeholder

func (s *UsageCleanupService) maxRangeDays() int {
	if s == nil || s.cfg == nil {
		return 31
placeholder
	if s.cfg.UsageCleanup.MaxRangeDays > 0 {
		return s.cfg.UsageCleanup.MaxRangeDays
placeholder
	return 31
placeholder

func (s *UsageCleanupService) batchSize() int {
	if s == nil || s.cfg == nil {
		return 5000
placeholder
	if s.cfg.UsageCleanup.BatchSize > 0 {
		return s.cfg.UsageCleanup.BatchSize
placeholder
	return 5000
placeholder

func (s *UsageCleanupService) workerInterval() time.Duration {
	if s == nil || s.cfg == nil {
		return 10 * time.Second
placeholder
	if s.cfg.UsageCleanup.WorkerIntervalSeconds > 0 {
		return time.Duration(s.cfg.UsageCleanup.WorkerIntervalSeconds) * time.Second
placeholder
	return 10 * time.Second
placeholder

func (s *UsageCleanupService) taskTimeout() time.Duration {
	if s == nil || s.cfg == nil {
		return 30 * time.Minute
placeholder
	if s.cfg.UsageCleanup.TaskTimeoutSeconds > 0 {
		return time.Duration(s.cfg.UsageCleanup.TaskTimeoutSeconds) * time.Second
placeholder
	return 30 * time.Minute
placeholder
