package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type cleanupRepoStub struct {
	mu         sync.Mutex
	created    []*service.UsageCleanupTask
	listTasks  []service.UsageCleanupTask
	listResult *pagination.PaginationResult
	listErr    error
	statusByID map[int64]string
placeholder

func (s *cleanupRepoStub) CreateTask(ctx context.Context, task *service.UsageCleanupTask) error {
	if task == nil {
		return nil
placeholder
	s.mu.Lock()
	defer s.mu.Unlock()
	if task.ID == 0 {
		task.ID = int64(len(s.created) + 1)
placeholder
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now().UTC()
placeholder
	task.UpdatedAt = task.CreatedAt
	clone := *task
	s.created = append(s.created, &clone)
	return nil
placeholder

func (s *cleanupRepoStub) ListTasks(ctx context.Context, params pagination.PaginationParams) ([]service.UsageCleanupTask, *pagination.PaginationResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.listTasks, s.listResult, s.listErr
placeholder

func (s *cleanupRepoStub) ClaimNextPendingTask(ctx context.Context, staleRunningAfterSeconds int64) (*service.UsageCleanupTask, error) {
	return nil, nil
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
	return nil
placeholder

func (s *cleanupRepoStub) CancelTask(ctx context.Context, taskID int64, canceledBy int64) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.statusByID == nil {
		s.statusByID = map[int64]string{placeholder
placeholder
	status := s.statusByID[taskID]
	if status != service.UsageCleanupStatusPending && status != service.UsageCleanupStatusRunning {
		return false, nil
placeholder
	s.statusByID[taskID] = service.UsageCleanupStatusCanceled
	return true, nil
placeholder

func (s *cleanupRepoStub) MarkTaskSucceeded(ctx context.Context, taskID int64, deletedRows int64) error {
	return nil
placeholder

func (s *cleanupRepoStub) MarkTaskFailed(ctx context.Context, taskID int64, deletedRows int64, errorMsg string) error {
	return nil
placeholder

func (s *cleanupRepoStub) DeleteUsageLogsBatch(ctx context.Context, filters service.UsageCleanupFilters, limit int) (int64, error) {
	return 0, nil
placeholder

var _ service.UsageCleanupRepository = (*cleanupRepoStub)(nil)

func setupCleanupRouter(cleanupService *service.UsageCleanupService, userID int64) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	if userID > 0 {
		router.Use(func(c *gin.Context) {
			c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: userIDplaceholder)
			c.Next()
	placeholder)
placeholder

	handler := NewUsageHandler(nil, nil, nil, cleanupService)
	router.POST("/api/v1/admin/usage/cleanup-tasks", handler.CreateCleanupTask)
	router.GET("/api/v1/admin/usage/cleanup-tasks", handler.ListCleanupTasks)
	router.POST("/api/v1/admin/usage/cleanup-tasks/:id/cancel", handler.CancelCleanupTask)
	return router
placeholder

func TestUsageHandlerCreateCleanupTaskUnauthorized(t *testing.T) {
	repo := &cleanupRepoStub{placeholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: true, MaxRangeDays: 31placeholderplaceholder
	cleanupService := service.NewUsageCleanupService(repo, nil, nil, cfg)
	router := setupCleanupRouter(cleanupService, 0)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/usage/cleanup-tasks", bytes.NewBufferString(`{placeholder`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusUnauthorized, recorder.Code)
placeholder

func TestUsageHandlerCreateCleanupTaskUnavailable(t *testing.T) {
	router := setupCleanupRouter(nil, 1)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/usage/cleanup-tasks", bytes.NewBufferString(`{placeholder`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusServiceUnavailable, recorder.Code)
placeholder

func TestUsageHandlerCreateCleanupTaskBindError(t *testing.T) {
	repo := &cleanupRepoStub{placeholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: true, MaxRangeDays: 31placeholderplaceholder
	cleanupService := service.NewUsageCleanupService(repo, nil, nil, cfg)
	router := setupCleanupRouter(cleanupService, 88)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/usage/cleanup-tasks", bytes.NewBufferString("{bad-json"))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
placeholder

func TestUsageHandlerCreateCleanupTaskMissingRange(t *testing.T) {
	repo := &cleanupRepoStub{placeholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: true, MaxRangeDays: 31placeholderplaceholder
	cleanupService := service.NewUsageCleanupService(repo, nil, nil, cfg)
	router := setupCleanupRouter(cleanupService, 88)

	payload := map[string]any{
		"start_date": "2024-01-01",
		"timezone":   "UTC",
placeholder
	body, err := json.Marshal(payload)
placeholder

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/usage/cleanup-tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
placeholder

func TestUsageHandlerCreateCleanupTaskInvalidDate(t *testing.T) {
	repo := &cleanupRepoStub{placeholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: true, MaxRangeDays: 31placeholderplaceholder
	cleanupService := service.NewUsageCleanupService(repo, nil, nil, cfg)
	router := setupCleanupRouter(cleanupService, 88)

	payload := map[string]any{
		"start_date": "2024-13-01",
		"end_date":   "2024-01-02",
		"timezone":   "UTC",
placeholder
	body, err := json.Marshal(payload)
placeholder

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/usage/cleanup-tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
placeholder

func TestUsageHandlerCreateCleanupTaskInvalidEndDate(t *testing.T) {
	repo := &cleanupRepoStub{placeholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: true, MaxRangeDays: 31placeholderplaceholder
	cleanupService := service.NewUsageCleanupService(repo, nil, nil, cfg)
	router := setupCleanupRouter(cleanupService, 88)

	payload := map[string]any{
		"start_date": "2024-01-01",
		"end_date":   "2024-02-40",
		"timezone":   "UTC",
placeholder
	body, err := json.Marshal(payload)
placeholder

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/usage/cleanup-tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
placeholder

func TestUsageHandlerCreateCleanupTaskSuccess(t *testing.T) {
	repo := &cleanupRepoStub{placeholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: true, MaxRangeDays: 31placeholderplaceholder
	cleanupService := service.NewUsageCleanupService(repo, nil, nil, cfg)
	router := setupCleanupRouter(cleanupService, 99)

	payload := map[string]any{
		"start_date": " 2024-01-01 ",
		"end_date":   "2024-01-02",
		"timezone":   "UTC",
		"model":      "gpt-4",
placeholder
	body, err := json.Marshal(payload)
placeholder

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/usage/cleanup-tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)

	var resp response.Response
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	require.Equal(t, 0, resp.Code)

	repo.mu.Lock()
	defer repo.mu.Unlock()
	require.Len(t, repo.created, 1)
	created := repo.created[0]
	require.Equal(t, int64(99), created.CreatedBy)
	require.NotNil(t, created.Filters.Model)
	require.Equal(t, "gpt-4", *created.Filters.Model)

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Add(24*time.Hour - time.Nanosecond)
	require.True(t, created.Filters.StartTime.Equal(start))
	require.True(t, created.Filters.EndTime.Equal(end))
placeholder

func TestUsageHandlerListCleanupTasksUnavailable(t *testing.T) {
	router := setupCleanupRouter(nil, 0)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/usage/cleanup-tasks", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusServiceUnavailable, recorder.Code)
placeholder

func TestUsageHandlerListCleanupTasksSuccess(t *testing.T) {
	repo := &cleanupRepoStub{placeholder
	repo.listTasks = []service.UsageCleanupTask{
		{
			ID:        7,
			Status:    service.UsageCleanupStatusSucceeded,
			CreatedBy: 4,
	placeholder,
placeholder
	repo.listResult = &pagination.PaginationResult{Total: 1, Page: 1, PageSize: 20, Pages: 1placeholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: true, MaxRangeDays: 31placeholderplaceholder
	cleanupService := service.NewUsageCleanupService(repo, nil, nil, cfg)
	router := setupCleanupRouter(cleanupService, 1)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/usage/cleanup-tasks", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)

	var resp struct {
		Code int `json:"code"`
		Data struct {
			Items []dto.UsageCleanupTask `json:"items"`
			Total int64                  `json:"total"`
			Page  int                    `json:"page"`
	placeholder `json:"data"`
placeholder
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	require.Equal(t, 0, resp.Code)
	require.Len(t, resp.Data.Items, 1)
	require.Equal(t, int64(7), resp.Data.Items[0].ID)
	require.Equal(t, int64(1), resp.Data.Total)
	require.Equal(t, 1, resp.Data.Page)
placeholder

func TestUsageHandlerListCleanupTasksError(t *testing.T) {
	repo := &cleanupRepoStub{listErr: errors.New("boom")placeholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: true, MaxRangeDays: 31placeholderplaceholder
	cleanupService := service.NewUsageCleanupService(repo, nil, nil, cfg)
	router := setupCleanupRouter(cleanupService, 1)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/usage/cleanup-tasks", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
placeholder

func TestUsageHandlerCancelCleanupTaskUnauthorized(t *testing.T) {
	repo := &cleanupRepoStub{placeholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: trueplaceholderplaceholder
	cleanupService := service.NewUsageCleanupService(repo, nil, nil, cfg)
	router := setupCleanupRouter(cleanupService, 0)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/usage/cleanup-tasks/1/cancel", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
placeholder

func TestUsageHandlerCancelCleanupTaskNotFound(t *testing.T) {
	repo := &cleanupRepoStub{placeholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: trueplaceholderplaceholder
	cleanupService := service.NewUsageCleanupService(repo, nil, nil, cfg)
	router := setupCleanupRouter(cleanupService, 1)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/usage/cleanup-tasks/999/cancel", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
placeholder

func TestUsageHandlerCancelCleanupTaskConflict(t *testing.T) {
	repo := &cleanupRepoStub{statusByID: map[int64]string{2: service.UsageCleanupStatusSucceededplaceholderplaceholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: trueplaceholderplaceholder
	cleanupService := service.NewUsageCleanupService(repo, nil, nil, cfg)
	router := setupCleanupRouter(cleanupService, 1)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/usage/cleanup-tasks/2/cancel", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusConflict, rec.Code)
placeholder

func TestUsageHandlerCancelCleanupTaskSuccess(t *testing.T) {
	repo := &cleanupRepoStub{statusByID: map[int64]string{3: service.UsageCleanupStatusPendingplaceholderplaceholder
	cfg := &config.Config{UsageCleanup: config.UsageCleanupConfig{Enabled: trueplaceholderplaceholder
	cleanupService := service.NewUsageCleanupService(repo, nil, nil, cfg)
	router := setupCleanupRouter(cleanupService, 1)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/usage/cleanup-tasks/3/cancel", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
placeholder
