package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/enttest"
	dbusagecleanuptask "github.com/Wei-Shaw/sub2api/ent/usagecleanuptask"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "modernc.org/sqlite"
)

func newUsageCleanupEntRepo(t *testing.T) (*usageCleanupRepository, *dbent.Client) {
placeholder
	db, err := sql.Open("sqlite", "file:usage_cleanup?mode=memory&cache=shared")
placeholder
	t.Cleanup(func() { _ = db.Close() placeholder)
	_, err = db.Exec("PRAGMA foreign_keys = ON")
placeholder

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := enttest.NewClient(t, enttest.WithOptions(dbent.Driver(drv)))
	t.Cleanup(func() { _ = client.Close() placeholder)

	repo := &usageCleanupRepository{client: client, sql: dbplaceholder
	return repo, client
placeholder

func TestUsageCleanupRepositoryEntCreateAndList(t *testing.T) {
	repo, _ := newUsageCleanupEntRepo(t)

	start := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)
	task := &service.UsageCleanupTask{
		Status:    service.UsageCleanupStatusPending,
		Filters:   service.UsageCleanupFilters{StartTime: start, EndTime: endplaceholder,
		CreatedBy: 9,
placeholder
	require.NoError(t, repo.CreateTask(context.Background(), task))
	require.NotZero(t, task.ID)

	task2 := &service.UsageCleanupTask{
		Status:    service.UsageCleanupStatusRunning,
		Filters:   service.UsageCleanupFilters{StartTime: start.Add(-24 * time.Hour), EndTime: end.Add(-24 * time.Hour)placeholder,
		CreatedBy: 10,
placeholder
	require.NoError(t, repo.CreateTask(context.Background(), task2))

	tasks, result, err := repo.ListTasks(context.Background(), pagination.PaginationParams{Page: 1, PageSize: 10placeholder)
placeholder
	require.Len(t, tasks, 2)
	require.Equal(t, int64(2), result.Total)
	require.Greater(t, tasks[0].ID, tasks[1].ID)
	require.Equal(t, start, tasks[1].Filters.StartTime)
	require.Equal(t, end, tasks[1].Filters.EndTime)
placeholder

func TestUsageCleanupRepositoryEntListEmpty(t *testing.T) {
	repo, _ := newUsageCleanupEntRepo(t)

	tasks, result, err := repo.ListTasks(context.Background(), pagination.PaginationParams{Page: 1, PageSize: 10placeholder)
placeholder
	require.Empty(t, tasks)
	require.Equal(t, int64(0), result.Total)
placeholder

func TestUsageCleanupRepositoryEntGetStatusAndProgress(t *testing.T) {
	repo, client := newUsageCleanupEntRepo(t)

	task := &service.UsageCleanupTask{
		Status:    service.UsageCleanupStatusPending,
		Filters:   service.UsageCleanupFilters{StartTime: time.Now().UTC(), EndTime: time.Now().UTC().Add(time.Hour)placeholder,
		CreatedBy: 3,
placeholder
	require.NoError(t, repo.CreateTask(context.Background(), task))

	status, err := repo.GetTaskStatus(context.Background(), task.ID)
placeholder
	require.Equal(t, service.UsageCleanupStatusPending, status)

	_, err = repo.GetTaskStatus(context.Background(), task.ID+99)
	require.ErrorIs(t, err, sql.ErrNoRows)

	require.NoError(t, repo.UpdateTaskProgress(context.Background(), task.ID, 42))
	loaded, err := client.UsageCleanupTask.Get(context.Background(), task.ID)
placeholder
	require.Equal(t, int64(42), loaded.DeletedRows)
placeholder

func TestUsageCleanupRepositoryEntCancelAndFinish(t *testing.T) {
	repo, client := newUsageCleanupEntRepo(t)

	task := &service.UsageCleanupTask{
		Status:    service.UsageCleanupStatusPending,
		Filters:   service.UsageCleanupFilters{StartTime: time.Now().UTC(), EndTime: time.Now().UTC().Add(time.Hour)placeholder,
		CreatedBy: 5,
placeholder
	require.NoError(t, repo.CreateTask(context.Background(), task))

	ok, err := repo.CancelTask(context.Background(), task.ID, 7)
placeholder
	require.True(t, ok)

	loaded, err := client.UsageCleanupTask.Get(context.Background(), task.ID)
placeholder
	require.Equal(t, service.UsageCleanupStatusCanceled, loaded.Status)
	require.NotNil(t, loaded.CanceledBy)
	require.NotNil(t, loaded.CanceledAt)
	require.NotNil(t, loaded.FinishedAt)

	loaded.Status = service.UsageCleanupStatusSucceeded
	_, err = client.UsageCleanupTask.Update().Where(dbusagecleanuptask.IDEQ(task.ID)).SetStatus(loaded.Status).Save(context.Background())
placeholder

	ok, err = repo.CancelTask(context.Background(), task.ID, 7)
placeholder
	require.False(t, ok)
placeholder

func TestUsageCleanupRepositoryEntCancelError(t *testing.T) {
	repo, client := newUsageCleanupEntRepo(t)

	task := &service.UsageCleanupTask{
		Status:    service.UsageCleanupStatusPending,
		Filters:   service.UsageCleanupFilters{StartTime: time.Now().UTC(), EndTime: time.Now().UTC().Add(time.Hour)placeholder,
		CreatedBy: 5,
placeholder
	require.NoError(t, repo.CreateTask(context.Background(), task))

	require.NoError(t, client.Close())
	_, err := repo.CancelTask(context.Background(), task.ID, 7)
placeholder
placeholder

func TestUsageCleanupRepositoryEntMarkResults(t *testing.T) {
	repo, client := newUsageCleanupEntRepo(t)

	task := &service.UsageCleanupTask{
		Status:    service.UsageCleanupStatusRunning,
		Filters:   service.UsageCleanupFilters{StartTime: time.Now().UTC(), EndTime: time.Now().UTC().Add(time.Hour)placeholder,
		CreatedBy: 12,
placeholder
	require.NoError(t, repo.CreateTask(context.Background(), task))

	require.NoError(t, repo.MarkTaskSucceeded(context.Background(), task.ID, 6))
	loaded, err := client.UsageCleanupTask.Get(context.Background(), task.ID)
placeholder
	require.Equal(t, service.UsageCleanupStatusSucceeded, loaded.Status)
	require.Equal(t, int64(6), loaded.DeletedRows)
	require.NotNil(t, loaded.FinishedAt)

	task2 := &service.UsageCleanupTask{
		Status:    service.UsageCleanupStatusRunning,
		Filters:   service.UsageCleanupFilters{StartTime: time.Now().UTC(), EndTime: time.Now().UTC().Add(time.Hour)placeholder,
		CreatedBy: 12,
placeholder
	require.NoError(t, repo.CreateTask(context.Background(), task2))

	require.NoError(t, repo.MarkTaskFailed(context.Background(), task2.ID, 4, "boom"))
	loaded2, err := client.UsageCleanupTask.Get(context.Background(), task2.ID)
placeholder
	require.Equal(t, service.UsageCleanupStatusFailed, loaded2.Status)
	require.Equal(t, "boom", *loaded2.ErrorMessage)
placeholder

func TestUsageCleanupRepositoryEntInvalidStatus(t *testing.T) {
	repo, _ := newUsageCleanupEntRepo(t)

	task := &service.UsageCleanupTask{
		Status:    "invalid",
		Filters:   service.UsageCleanupFilters{StartTime: time.Now().UTC(), EndTime: time.Now().UTC().Add(time.Hour)placeholder,
		CreatedBy: 1,
placeholder
	require.Error(t, repo.CreateTask(context.Background(), task))
placeholder

func TestUsageCleanupRepositoryEntListInvalidFilters(t *testing.T) {
	repo, client := newUsageCleanupEntRepo(t)

	now := time.Now().UTC()
	driver, ok := client.Driver().(*entsql.Driver)
	require.True(t, ok)
	_, err := driver.DB().ExecContext(
		context.Background(),
		`INSERT INTO usage_cleanup_tasks (status, filters, created_by, deleted_rows, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		service.UsageCleanupStatusPending,
		[]byte("invalid-json"),
		int64(1),
		int64(0),
		now,
		now,
	)
placeholder

	_, _, err = repo.ListTasks(context.Background(), pagination.PaginationParams{Page: 1, PageSize: 10placeholder)
placeholder
placeholder

func TestUsageCleanupTaskFromEntFull(t *testing.T) {
	start := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)
	errMsg := "failed"
	canceledBy := int64(2)
	canceledAt := start.Add(time.Minute)
	startedAt := start.Add(2 * time.Minute)
	finishedAt := start.Add(3 * time.Minute)
	filters := service.UsageCleanupFilters{StartTime: start, EndTime: endplaceholder
	filtersJSON, err := json.Marshal(filters)
placeholder

	task, err := usageCleanupTaskFromEnt(&dbent.UsageCleanupTask{
		ID:           10,
		Status:       service.UsageCleanupStatusFailed,
		Filters:      filtersJSON,
		CreatedBy:    11,
		DeletedRows:  7,
		ErrorMessage: &errMsg,
		CanceledBy:   &canceledBy,
		CanceledAt:   &canceledAt,
		StartedAt:    &startedAt,
		FinishedAt:   &finishedAt,
		CreatedAt:    start,
		UpdatedAt:    end,
placeholder)
placeholder
	require.Equal(t, int64(10), task.ID)
	require.Equal(t, service.UsageCleanupStatusFailed, task.Status)
	require.NotNil(t, task.ErrorMsg)
	require.NotNil(t, task.CanceledBy)
	require.NotNil(t, task.CanceledAt)
	require.NotNil(t, task.StartedAt)
	require.NotNil(t, task.FinishedAt)
placeholder

func TestUsageCleanupTaskFromEntInvalidFilters(t *testing.T) {
	task, err := usageCleanupTaskFromEnt(&dbent.UsageCleanupTask{
		Filters: json.RawMessage("invalid-json"),
placeholder)
placeholder
	require.Empty(t, task)
placeholder
