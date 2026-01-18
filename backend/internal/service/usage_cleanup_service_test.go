package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

type cleanupDeleteResponse struct {
	deleted int64
	err     error
placeholder

type cleanupDeleteCall struct {
	filters UsageCleanupFilters
	limit   int
placeholder

type cleanupMarkCall struct {
	taskID      int64
	deletedRows int64
	errMsg      string
placeholder

type cleanupRepoStub struct {
	mu            sync.Mutex
	created       []*UsageCleanupTask
	createErr     error
	listTasks     []UsageCleanupTask
	listResult    *pagination.PaginationResult
	listErr       error
	claimQueue    []*UsageCleanupTask
	claimErr      error
	deleteQueue   []cleanupDeleteResponse
	deleteCalls   []cleanupDeleteCall
	markSucceeded []cleanupMarkCall
	markFailed    []cleanupMarkCall
	statusByID    map[int64]string
	progressCalls []cleanupMarkCall
	cancelCalls   []int64
placeholder

func (s *cleanupRepoStub) CreateTask(ctx context.Context, task *UsageCleanupTask) error {
	if task == nil {
		return nil
placeholder
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.createErr != nil {
		return s.createErr
placeholder
	if task.ID == 0 {
		task.ID = int64(len(s.created) + 1)
placeholder
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now().UTC()
placeholder
	if task.UpdatedAt.IsZero() {
		task.UpdatedAt = task.CreatedAt
placeholder
	clone := *task
	s.created = append(s.created, &clone)
	return nil
placeholder

func (s *cleanupRepoStub) ListTasks(ctx context.Context, params pagination.PaginationParams) ([]UsageCleanupTask, *pagination.PaginationResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.listTasks, s.listResult, s.listErr
placeholder

func (s *cleanupRepoStub) ClaimNextPendingTask(ctx context.Context, staleRunningAfterSeconds int64) (*UsageCleanupTask, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.claimErr != nil {
		return nil, s.claimErr
placeholder
	if len(s.claimQueue) == 0 {
		return nil, nil
placeholder
	task := s.claimQueue[0]
	s.claimQueue = s.claimQueue[1:]
	if s.statusByID == nil {
		s.statusByID = map[int64]string{placeholder
placeholder
	s.statusByID[task.ID] = UsageCleanupStatusRunning
	return task, nil
placeholder

func (s *cleanupRepoStub) GetTaskStatus(ctx context.Context, taskID int64) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.statusByID == nil {
		return "", sql.ErrNoRows
placeholder
	status, ok := s.statusByID[taskID]
	if !ok {
		return "", sql.ErrNoRows
placeholder
	return status, nil
placeholder

func (s *cleanupRepoStub) UpdateTaskProgress(ctx context.Context, taskID int64, deletedRows int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.progressCalls = append(s.progressCalls, cleanupMarkCall{taskID: taskID, deletedRows: deletedRowsplaceholder)
	return nil
placeholder

func (s *cleanupRepoStub) CancelTask(ctx context.Context, taskID int64, canceledBy int64) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cancelCalls = append(s.cancelCalls, taskID)
	if s.statusByID == nil {
		s.statusByID = map[int64]string{placeholder
placeholder
	status := s.statusByID[taskID]
	if status != UsageCleanupStatusPending && status != UsageCleanupStatusRunning {
		return false, nil
placeholder
	s.statusByID[taskID] = UsageCleanupStatusCanceled
	return true, nil
placeholder

func (s *cleanupRepoStub) MarkTaskSucceeded(ctx context.Context, taskID int64, deletedRows int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.markSucceeded = append(s.markSucceeded, cleanupMarkCall{taskID: taskID, deletedRows: deletedRowsplaceholder)
	if s.statusByID == nil {
		s.statusByID = map[int64]string{placeholder
placeholder
	s.statusByID[taskID] = UsageCleanupStatusSucceeded
	return nil
placeholder

func (s *cleanupRepoStub) MarkTaskFailed(ctx context.Context, taskID int64, deletedRows int64, errorMsg string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.markFailed = append(s.markFailed, cleanupMarkCall{taskID: taskID, deletedRows: deletedRows, errMsg: errorMsgplaceholder)
	if s.statusByID == nil {
		s.statusByID = map[int64]string{placeholder
placeholder
	s.statusByID[taskID] = UsageCleanupStatusFailed
	return nil
placeholder

func (s *cleanupRepoStub) DeleteUsageLogsBatch(ctx context.Context, filters UsageCleanupFilters, limit int) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deleteCalls = append(s.deleteCalls, cleanupDeleteCall{filters: filters, limit: limitplaceholder)
	if len(s.deleteQueue) == 0 {
		return 0, nil
placeholder
	resp := s.deleteQueue[0]
	s.deleteQueue = s.deleteQueue[1:]
	return resp.deleted, resp.err
placeholder

func TestUsageCleanupServiceCreateTaskSanitizeFilters(t *testing.T) {
	repo := &cleanupRepoStub{placeholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: true, MaxRangeDays: 31placeholderplaceholder
	svc := NewUsageCleanupService(repo, nil, nil, cfg)

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)
	userID := int64(-1)
	apiKeyID := int64(10)
	model := "  gpt-4  "
	billingType := int8(-2)
	filters := UsageCleanupFilters{
		StartTime:   start,
		EndTime:     end,
		UserID:      &userID,
		APIKeyID:    &apiKeyID,
		Model:       &model,
		BillingType: &billingType,
placeholder

	task, err := svc.CreateTask(context.Background(), filters, 9)
placeholder
	require.Equal(t, UsageCleanupStatusPending, task.Status)
	require.Nil(t, task.Filters.UserID)
	require.NotNil(t, task.Filters.APIKeyID)
	require.Equal(t, apiKeyID, *task.Filters.APIKeyID)
	require.NotNil(t, task.Filters.Model)
	require.Equal(t, "gpt-4", *task.Filters.Model)
	require.Nil(t, task.Filters.BillingType)
	require.Equal(t, int64(9), task.CreatedBy)
placeholder

func TestUsageCleanupServiceCreateTaskInvalidCreator(t *testing.T) {
	repo := &cleanupRepoStub{placeholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: trueplaceholderplaceholder
	svc := NewUsageCleanupService(repo, nil, nil, cfg)

	filters := UsageCleanupFilters{
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
placeholder
	_, err := svc.CreateTask(context.Background(), filters, 0)
placeholder
	require.Equal(t, "USAGE_CLEANUP_INVALID_CREATOR", infraerrors.Reason(err))
placeholder

func TestUsageCleanupServiceCreateTaskDisabled(t *testing.T) {
	repo := &cleanupRepoStub{placeholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: falseplaceholderplaceholder
	svc := NewUsageCleanupService(repo, nil, nil, cfg)

	filters := UsageCleanupFilters{
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
placeholder
	_, err := svc.CreateTask(context.Background(), filters, 1)
placeholder
	require.Equal(t, http.StatusServiceUnavailable, infraerrors.Code(err))
	require.Equal(t, "USAGE_CLEANUP_DISABLED", infraerrors.Reason(err))
placeholder

func TestUsageCleanupServiceCreateTaskRangeTooLarge(t *testing.T) {
	repo := &cleanupRepoStub{placeholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: true, MaxRangeDays: 1placeholderplaceholder
	svc := NewUsageCleanupService(repo, nil, nil, cfg)

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(48 * time.Hour)
	filters := UsageCleanupFilters{StartTime: start, EndTime: endplaceholder

	_, err := svc.CreateTask(context.Background(), filters, 1)
placeholder
	require.Equal(t, "USAGE_CLEANUP_RANGE_TOO_LARGE", infraerrors.Reason(err))
placeholder

func TestUsageCleanupServiceCreateTaskMissingRange(t *testing.T) {
	repo := &cleanupRepoStub{placeholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: trueplaceholderplaceholder
	svc := NewUsageCleanupService(repo, nil, nil, cfg)

	_, err := svc.CreateTask(context.Background(), UsageCleanupFilters{placeholder, 1)
placeholder
	require.Equal(t, "USAGE_CLEANUP_MISSING_RANGE", infraerrors.Reason(err))
placeholder

func TestUsageCleanupServiceCreateTaskRepoError(t *testing.T) {
	repo := &cleanupRepoStub{createErr: errors.New("db down")placeholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: trueplaceholderplaceholder
	svc := NewUsageCleanupService(repo, nil, nil, cfg)

	filters := UsageCleanupFilters{
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
placeholder
	_, err := svc.CreateTask(context.Background(), filters, 1)
placeholder
	require.Contains(t, err.Error(), "create cleanup task")
placeholder

func TestUsageCleanupServiceRunOnceSuccess(t *testing.T) {
	repo := &cleanupRepoStub{
		claimQueue: []*UsageCleanupTask{
			{ID: 5, Filters: UsageCleanupFilters{StartTime: time.Now(), EndTime: time.Now().Add(2 * time.Hour)placeholderplaceholder,
	placeholder,
		deleteQueue: []cleanupDeleteResponse{
			{deleted: 2placeholder,
			{deleted: 2placeholder,
			{deleted: 1placeholder,
	placeholder,
placeholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: true, BatchSize: 2, TaskTimeoutSeconds: 30placeholderplaceholder
	svc := NewUsageCleanupService(repo, nil, nil, cfg)

	svc.runOnce()

	repo.mu.Lock()
	defer repo.mu.Unlock()
	require.Len(t, repo.deleteCalls, 3)
	require.Len(t, repo.markSucceeded, 1)
	require.Empty(t, repo.markFailed)
	require.Equal(t, int64(5), repo.markSucceeded[0].taskID)
	require.Equal(t, int64(5), repo.markSucceeded[0].deletedRows)
placeholder

func TestUsageCleanupServiceRunOnceClaimError(t *testing.T) {
	repo := &cleanupRepoStub{claimErr: errors.New("claim failed")placeholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: trueplaceholderplaceholder
	svc := NewUsageCleanupService(repo, nil, nil, cfg)
	svc.runOnce()

	repo.mu.Lock()
	defer repo.mu.Unlock()
	require.Empty(t, repo.markSucceeded)
	require.Empty(t, repo.markFailed)
placeholder

func TestUsageCleanupServiceRunOnceAlreadyRunning(t *testing.T) {
	repo := &cleanupRepoStub{placeholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: trueplaceholderplaceholder
	svc := NewUsageCleanupService(repo, nil, nil, cfg)
	svc.running = 1
	svc.runOnce()
placeholder

func TestUsageCleanupServiceExecuteTaskFailed(t *testing.T) {
	longMsg := strings.Repeat("x", 600)
	repo := &cleanupRepoStub{
		deleteQueue: []cleanupDeleteResponse{
			{err: errors.New(longMsg)placeholder,
	placeholder,
placeholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: true, BatchSize: 3placeholderplaceholder
	svc := NewUsageCleanupService(repo, nil, nil, cfg)
	task := &UsageCleanupTask{
		ID: 11,
		Filters: UsageCleanupFilters{
			StartTime: time.Now(),
			EndTime:   time.Now().Add(24 * time.Hour),
	placeholder,
placeholder

	svc.executeTask(context.Background(), task)

	repo.mu.Lock()
	defer repo.mu.Unlock()
	require.Len(t, repo.markFailed, 1)
	require.Equal(t, int64(11), repo.markFailed[0].taskID)
	require.Equal(t, 500, len(repo.markFailed[0].errMsg))
placeholder

func TestUsageCleanupServiceListTasks(t *testing.T) {
	repo := &cleanupRepoStub{
		listTasks: []UsageCleanupTask{{ID: 1placeholder, {ID: 2placeholderplaceholder,
		listResult: &pagination.PaginationResult{
			Total:    2,
			Page:     1,
			PageSize: 20,
			Pages:    1,
	placeholder,
placeholder
	svc := NewUsageCleanupService(repo, nil, nil, &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: trueplaceholderplaceholder)

	tasks, result, err := svc.ListTasks(context.Background(), pagination.PaginationParams{Page: 1, PageSize: 20placeholder)
placeholder
	require.Len(t, tasks, 2)
	require.Equal(t, int64(2), result.Total)
placeholder

func TestUsageCleanupServiceListTasksNotReady(t *testing.T) {
	var nilSvc *UsageCleanupService
	_, _, err := nilSvc.ListTasks(context.Background(), pagination.PaginationParams{Page: 1, PageSize: 20placeholder)
placeholder

	svc := NewUsageCleanupService(nil, nil, nil, &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: trueplaceholderplaceholder)
	_, _, err = svc.ListTasks(context.Background(), pagination.PaginationParams{Page: 1, PageSize: 20placeholder)
placeholder
placeholder

func TestUsageCleanupServiceDefaultsAndLifecycle(t *testing.T) {
	var nilSvc *UsageCleanupService
	require.Equal(t, 31, nilSvc.maxRangeDays())
	require.Equal(t, 5000, nilSvc.batchSize())
	require.Equal(t, 10*time.Second, nilSvc.workerInterval())
	require.Equal(t, 30*time.Minute, nilSvc.taskTimeout())
	nilSvc.Start()
	nilSvc.Stop()

	repo := &cleanupRepoStub{placeholder
	cfgDisabled := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: falseplaceholderplaceholder
	svcDisabled := NewUsageCleanupService(repo, nil, nil, cfgDisabled)
	svcDisabled.Start()
	svcDisabled.Stop()

	timingWheel, err := NewTimingWheelService()
placeholder

	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: true, WorkerIntervalSeconds: 5placeholderplaceholder
	svc := NewUsageCleanupService(repo, timingWheel, nil, cfg)
	require.Equal(t, 5*time.Second, svc.workerInterval())
	svc.Start()
	svc.Stop()

	cfgFallback := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: trueplaceholderplaceholder
	svcFallback := NewUsageCleanupService(repo, timingWheel, nil, cfgFallback)
	require.Equal(t, 31, svcFallback.maxRangeDays())
	require.Equal(t, 5000, svcFallback.batchSize())
	require.Equal(t, 10*time.Second, svcFallback.workerInterval())

	svcMissingDeps := NewUsageCleanupService(nil, nil, nil, cfgFallback)
	svcMissingDeps.Start()
placeholder

func TestSanitizeUsageCleanupFiltersModelEmpty(t *testing.T) {
	model := "   "
	apiKeyID := int64(-5)
	accountID := int64(-1)
	groupID := int64(-2)
	filters := UsageCleanupFilters{
		UserID:    &apiKeyID,
		APIKeyID:  &apiKeyID,
		AccountID: &accountID,
		GroupID:   &groupID,
		Model:     &model,
placeholder

	sanitizeUsageCleanupFilters(&filters)
	require.Nil(t, filters.UserID)
	require.Nil(t, filters.APIKeyID)
	require.Nil(t, filters.AccountID)
	require.Nil(t, filters.GroupID)
	require.Nil(t, filters.Model)
placeholder
