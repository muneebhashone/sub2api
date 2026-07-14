package repository

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/enttest"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
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

	return &apiKeyRepository{client: client, sql: dbplaceholder, client
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

func mustCreateAPIKeyRepoAccount(t *testing.T, ctx context.Context, client *dbent.Client, name string) int64 {
placeholder
	a, err := client.Account.Create().
		SetName(name).
		SetPlatform(service.PlatformOpenAI).
		SetType(service.AccountTypeAPIKey).
		SetStatus(service.StatusActive).
		SetCredentials(map[string]any{"api_key": "sk-test"placeholder).
		Save(ctx)
placeholder
	return a.ID
placeholder

func mustCreateAPIKeyRepoUsageLog(t *testing.T, ctx context.Context, client *dbent.Client, userID, apiKeyID, accountID int64, requestID string, createdAt time.Time, ipAddress *string) {
placeholder
	builder := client.UsageLog.Create().
		SetUserID(userID).
		SetAPIKeyID(apiKeyID).
		SetAccountID(accountID).
		SetRequestID(requestID).
		SetModel("gpt-5").
		SetCreatedAt(createdAt)
	if ipAddress != nil {
		builder.SetIPAddress(*ipAddress)
placeholder
	_, err := builder.Save(ctx)
placeholder
placeholder

func TestAPIKeyRepositoryListByUserIDAttachesLastUsedIP(t *testing.T) {
	repo, client := newAPIKeyRepoSQLite(t)
	ctx := context.Background()
	user := mustCreateAPIKeyRepoUser(t, ctx, client, "list-last-used-ip@test.com")
	accountID := mustCreateAPIKeyRepoAccount(t, ctx, client, "acc-list-last-used-ip")

	withLogs := &service.APIKey{
		UserID: user.ID,
		Key:    "sk-list-last-used-ip-logs",
		Name:   "With Logs",
		Status: service.StatusActive,
placeholder
	emptyOnly := &service.APIKey{
		UserID: user.ID,
		Key:    "sk-list-last-used-ip-empty",
		Name:   "Empty Only",
		Status: service.StatusActive,
placeholder
	noLogs := &service.APIKey{
		UserID: user.ID,
		Key:    "sk-list-last-used-ip-none",
		Name:   "No Logs",
		Status: service.StatusActive,
placeholder
	require.NoError(t, repo.Create(ctx, withLogs))
	require.NoError(t, repo.Create(ctx, emptyOnly))
	require.NoError(t, repo.Create(ctx, noLogs))

	olderIP := "198.51.100.10"
	newerEmptyIP := ""
	newestIP := "203.0.113.20"
	base := time.Now().UTC().Add(-3 * time.Hour).Truncate(time.Second)
	mustCreateAPIKeyRepoUsageLog(t, ctx, client, user.ID, withLogs.ID, accountID, "req-last-ip-older", base, &olderIP)
	mustCreateAPIKeyRepoUsageLog(t, ctx, client, user.ID, withLogs.ID, accountID, "req-last-ip-empty", base.Add(time.Hour), &newerEmptyIP)
	mustCreateAPIKeyRepoUsageLog(t, ctx, client, user.ID, withLogs.ID, accountID, "req-last-ip-newest", base.Add(2*time.Hour), &newestIP)
	mustCreateAPIKeyRepoUsageLog(t, ctx, client, user.ID, emptyOnly.ID, accountID, "req-empty-ip", base.Add(3*time.Hour), &newerEmptyIP)

	keys, _, err := repo.ListByUserID(ctx, user.ID, pagination.PaginationParams{Page: 1, PageSize: 10placeholder, service.APIKeyListFilters{placeholder)
placeholder

	byID := make(map[int64]service.APIKey, len(keys))
	for _, key := range keys {
		byID[key.ID] = key
placeholder
	require.NotNil(t, byID[withLogs.ID].LastUsedIP)
	require.Equal(t, newestIP, *byID[withLogs.ID].LastUsedIP)
	require.Nil(t, byID[emptyOnly.ID].LastUsedIP)
	require.Nil(t, byID[noLogs.ID].LastUsedIP)
placeholder

func TestLatestUsageLogIPsQueryPostgresUsesPerKeyLateralLookup(t *testing.T) {
	query, args := latestUsageLogIPsQuery([]int64{11, 22placeholder, dialect.Postgres)
	normalizedQuery := strings.Join(strings.Fields(query), " ")

	require.Contains(t, normalizedQuery, "FROM unnest($1::bigint[]) AS requested(api_key_id)")
	require.Contains(t, normalizedQuery, "CROSS JOIN LATERAL")
	require.Contains(t, normalizedQuery, "WHERE ul.api_key_id = requested.api_key_id")
	require.Contains(t, normalizedQuery, "AND ul.ip_address IS NOT NULL")
	require.Contains(t, normalizedQuery, "AND ul.ip_address <> ''")
	require.Contains(t, normalizedQuery, "ORDER BY ul.created_at DESC, ul.id DESC LIMIT 1")
	require.NotContains(t, normalizedQuery, "ROW_NUMBER")
	require.Len(t, args, 1)
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
