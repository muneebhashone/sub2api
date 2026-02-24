package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/enttest"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "modernc.org/sqlite"
)

func newAPIKeyRepoSQLite(t *testing.T) (*apiKeyRepository, *dbent.Client) {
placeholder

	db, err := sql.Open("sqlite", "file:api_key_repo_last_used?mode=memory&cache=shared")
placeholder
	t.Cleanup(func() { _ = db.Close() placeholder)

	_, err = db.Exec("PRAGMA foreign_keys = ON")
placeholder

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := enttest.NewClient(t, enttest.WithOptions(dbent.Driver(drv)))
	t.Cleanup(func() { _ = client.Close() placeholder)

	return &apiKeyRepository{client: clientplaceholder, client
placeholder

func mustCreateAPIKeyRepoUser(t *testing.T, ctx context.Context, client *dbent.Client, email string) *service.User {
placeholder
	u, err := client.User.Create().
		SetEmail(email).
		SetPasswordHash("test-password-hash").
		SetRole(service.RoleUser).
		SetStatus(service.StatusActive).
		Save(ctx)
placeholder
	return userEntityToService(u)
placeholder

func TestAPIKeyRepository_CreateWithLastUsedAt(t *testing.T) {
	repo, client := newAPIKeyRepoSQLite(t)
	ctx := context.Background()
	user := mustCreateAPIKeyRepoUser(t, ctx, client, "create-last-used@test.com")

	lastUsed := time.Now().UTC().Add(-time.Hour).Truncate(time.Second)
	key := &service.APIKey{
		UserID:     user.ID,
		Key:        "sk-create-last-used",
		Name:       "CreateWithLastUsed",
		Status:     service.StatusActive,
		LastUsedAt: &lastUsed,
placeholder

	require.NoError(t, repo.Create(ctx, key))
	require.NotNil(t, key.LastUsedAt)
	require.WithinDuration(t, lastUsed, *key.LastUsedAt, time.Second)

	got, err := repo.GetByID(ctx, key.ID)
placeholder
	require.NotNil(t, got.LastUsedAt)
	require.WithinDuration(t, lastUsed, *got.LastUsedAt, time.Second)
placeholder

func TestAPIKeyRepository_UpdateLastUsed(t *testing.T) {
	repo, client := newAPIKeyRepoSQLite(t)
	ctx := context.Background()
	user := mustCreateAPIKeyRepoUser(t, ctx, client, "update-last-used@test.com")

	key := &service.APIKey{
		UserID: user.ID,
		Key:    "sk-update-last-used",
		Name:   "UpdateLastUsed",
		Status: service.StatusActive,
placeholder
	require.NoError(t, repo.Create(ctx, key))

	before, err := repo.GetByID(ctx, key.ID)
placeholder
	require.Nil(t, before.LastUsedAt)

	target := time.Now().UTC().Add(2 * time.Minute).Truncate(time.Second)
	require.NoError(t, repo.UpdateLastUsed(ctx, key.ID, target))

	after, err := repo.GetByID(ctx, key.ID)
placeholder
	require.NotNil(t, after.LastUsedAt)
	require.WithinDuration(t, target, *after.LastUsedAt, time.Second)
	require.WithinDuration(t, target, after.UpdatedAt, time.Second)
placeholder

func TestAPIKeyRepository_UpdateLastUsedDeletedKey(t *testing.T) {
	repo, client := newAPIKeyRepoSQLite(t)
	ctx := context.Background()
	user := mustCreateAPIKeyRepoUser(t, ctx, client, "deleted-last-used@test.com")

	key := &service.APIKey{
		UserID: user.ID,
		Key:    "sk-update-last-used-deleted",
		Name:   "UpdateLastUsedDeleted",
		Status: service.StatusActive,
placeholder
	require.NoError(t, repo.Create(ctx, key))
	require.NoError(t, repo.Delete(ctx, key.ID))

	err := repo.UpdateLastUsed(ctx, key.ID, time.Now().UTC())
	require.ErrorIs(t, err, service.ErrAPIKeyNotFound)
placeholder

func TestAPIKeyRepository_UpdateLastUsedDBError(t *testing.T) {
	repo, client := newAPIKeyRepoSQLite(t)
	ctx := context.Background()
	user := mustCreateAPIKeyRepoUser(t, ctx, client, "db-error-last-used@test.com")

	key := &service.APIKey{
		UserID: user.ID,
		Key:    "sk-update-last-used-db-error",
		Name:   "UpdateLastUsedDBError",
		Status: service.StatusActive,
placeholder
	require.NoError(t, repo.Create(ctx, key))

	require.NoError(t, client.Close())
	err := repo.UpdateLastUsed(ctx, key.ID, time.Now().UTC())
placeholder
placeholder

func TestAPIKeyRepository_CreateDuplicateKey(t *testing.T) {
	repo, client := newAPIKeyRepoSQLite(t)
	ctx := context.Background()
	user := mustCreateAPIKeyRepoUser(t, ctx, client, "duplicate-key@test.com")

	first := &service.APIKey{
		UserID: user.ID,
		Key:    "sk-duplicate",
		Name:   "first",
		Status: service.StatusActive,
placeholder
	second := &service.APIKey{
		UserID: user.ID,
		Key:    "sk-duplicate",
		Name:   "second",
		Status: service.StatusActive,
placeholder

	require.NoError(t, repo.Create(ctx, first))
	err := repo.Create(ctx, second)
	require.ErrorIs(t, err, service.ErrAPIKeyExists)
placeholder
