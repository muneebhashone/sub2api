package repository

import (
	"context"
	"database/sql"
	"regexp"
	"strings"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestAccountRepository_SetTempUnschedulable_NoRowsAffectedDoesNotWriteOutbox(t *testing.T) {
	exec := &recordingSQLExecutor{result: rowsAffectedResult(0)placeholder
	repo := newAccountRepositoryWithSQL(nil, exec, nil)
	until := time.Now().Add(10 * time.Minute)

	err := repo.SetTempUnschedulable(context.Background(), 42, until, "retry")
placeholder
	require.Len(t, exec.execQueries, 1)
	require.Contains(t, exec.execQueries[0], "UPDATE accounts")
	require.NotContains(t, strings.Join(exec.execQueries, "\n"), "scheduler_outbox")
placeholder

func TestAccountRepository_GrokCredentialConditionalMutationsAreEligibleAndAtomicallyPropagated(t *testing.T) {
	proxyID := int64(77)
	snapshot := service.GrokCredentialMutationSnapshot{
		CredentialsJSON: `{"access_token":"access","refresh_token":"refresh","_token_version":placeholder`,
		ProxyID:         &proxyID,
placeholder

	t.Run("permanent", func(t *testing.T) {
		exec := &recordingSQLExecutor{result: rowsAffectedResult(0)placeholder
		repo := newAccountRepositoryWithSQL(nil, exec, nil)

		updated, err := repo.SetGrokCredentialErrorIfMatch(context.Background(), 42, snapshot, "revoked")

	placeholder
		require.False(t, updated)
		require.Len(t, exec.execQueries, 1)
		normalized := normalizeSQLWhitespace(exec.execQueries[0])
		require.Contains(t, normalized, "WITH updated AS ( UPDATE accounts AS a")
		require.Contains(t, normalized, "a.schedulable IS TRUE")
		require.Contains(t, normalized, "a.temp_unschedulable_until IS NULL OR a.temp_unschedulable_until <= NOW()")
		require.Contains(t, normalized, "a.rate_limit_reset_at IS NULL OR a.rate_limit_reset_at <= NOW()")
		require.Contains(t, normalized, "a.overload_until IS NULL OR a.overload_until <= NOW()")
		require.Contains(t, normalized, "a.credentials = $7::jsonb")
		require.Contains(t, normalized, "a.proxy_id IS NOT DISTINCT FROM $8")
		require.Contains(t, normalized, "NOT EXISTS ( SELECT 1 FROM proxies p")
		require.Contains(t, normalized, "INSERT INTO scheduler_outbox")
		require.Len(t, exec.execArgs[0], 10)
		require.Equal(t, snapshot.CredentialsJSON, exec.execArgs[0][6])
		require.Equal(t, &proxyID, exec.execArgs[0][7])
		require.Equal(t, string(service.GrokCredentialReasonProxyInvalid), exec.execArgs[0][8])
		require.Equal(t, service.SchedulerOutboxEventAccountChanged, exec.execArgs[0][9])
placeholder)

	t.Run("transient", func(t *testing.T) {
		exec := &recordingSQLExecutor{result: rowsAffectedResult(0)placeholder
		repo := newAccountRepositoryWithSQL(nil, exec, nil)

		updated, err := repo.SetGrokCredentialTempUnschedulableIfMatch(
			context.Background(), 42, snapshot, time.Now().Add(time.Minute), "temporary",
		)

	placeholder
		require.False(t, updated)
		require.Len(t, exec.execQueries, 1)
		normalized := normalizeSQLWhitespace(exec.execQueries[0])
		require.Contains(t, normalized, "WITH updated AS ( UPDATE accounts AS a")
		require.Contains(t, normalized, "a.schedulable IS TRUE")
		require.Contains(t, normalized, "a.temp_unschedulable_until IS NULL OR a.temp_unschedulable_until <= NOW()")
		require.Contains(t, normalized, "a.rate_limit_reset_at IS NULL OR a.rate_limit_reset_at <= NOW()")
		require.Contains(t, normalized, "a.overload_until IS NULL OR a.overload_until <= NOW()")
		require.Contains(t, normalized, "a.credentials = $7::jsonb")
		require.Contains(t, normalized, "a.proxy_id IS NOT DISTINCT FROM $8")
		require.Contains(t, normalized, "INSERT INTO scheduler_outbox")
		require.Len(t, exec.execArgs[0], 9)
		require.Equal(t, snapshot.CredentialsJSON, exec.execArgs[0][6])
		require.Equal(t, &proxyID, exec.execArgs[0][7])
		require.Equal(t, service.SchedulerOutboxEventAccountChanged, exec.execArgs[0][8])
placeholder)
placeholder

func TestAccountRepository_GrokCredentialCommitCarriesOutboxAcrossCallerCancellation(t *testing.T) {
	snapshot := service.GrokCredentialMutationSnapshot{CredentialsJSON: `{"access_token":"access","refresh_token":"refresh"placeholder`placeholder
	tests := []struct {
		name   string
		mutate func(context.Context, *accountRepository) (bool, error)
placeholder{
		{
			name: "permanent",
			mutate: func(ctx context.Context, repo *accountRepository) (bool, error) {
				return repo.SetGrokCredentialErrorIfMatch(ctx, 42, snapshot, string(service.GrokCredentialReasonRevoked))
		placeholder,
	placeholder,
		{
			name: "transient",
			mutate: func(ctx context.Context, repo *accountRepository) (bool, error) {
				return repo.SetGrokCredentialTempUnschedulableIfMatch(ctx, 42, snapshot, time.Now().Add(time.Minute), string(service.GrokCredentialReasonRefreshTransient))
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			exec := &recordingSQLExecutor{result: rowsAffectedResult(1), afterExec: cancelplaceholder
			repo := newAccountRepositoryWithSQL(nil, exec, nil)

			updated, err := tt.mutate(ctx, repo)

		placeholder
			require.True(t, updated)
			require.ErrorIs(t, ctx.Err(), context.Canceled)
			require.Len(t, exec.execQueries, 1, "state update and scheduler outbox must share one atomic SQL statement")
			require.Contains(t, normalizeSQLWhitespace(exec.execQueries[0]), "INSERT INTO scheduler_outbox")
	placeholder)
placeholder
placeholder

func TestAccountRepository_ListOAuthRefreshCandidates_SQLFilter(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
placeholder
	defer func() { _ = db.Close() placeholder()

	var capturedSQL string
	mock.ExpectQuery("SELECT id").
		WillReturnRows(sqlmock.NewRows([]string{"id"placeholder)).
		WillDelayFor(0)

	repo := newAccountRepositoryWithSQL(nil, captureQuerySQL{db: db, captured: &capturedSQLplaceholder, nil)

	accounts, err := repo.ListOAuthRefreshCandidates(context.Background())
placeholder
	require.Empty(t, accounts)

	normalized := normalizeSQLWhitespace(capturedSQL)
	require.Contains(t, normalized, "deleted_at IS NULL")
	require.Contains(t, normalized, "status = 'active'")
	// setup-token 的 access_token 同为 8h 短期令牌，必须与 oauth 一起纳入后台刷新候选
	require.Contains(t, normalized, "type IN ('oauth', 'setup-token')")
	require.Contains(t, normalized, "platform IN ('anthropic', 'openai', 'gemini', 'antigravity')")
	require.Contains(t, normalized, "credentials ? 'refresh_token'")
	require.Contains(t, normalized, "btrim(credentials->>'refresh_token') <> ''")
	require.Contains(t, normalized, "temp_unschedulable_until > NOW()")
	require.Contains(t, normalized, "temp_unschedulable_reason LIKE 'token refresh retry exhausted:%'")
	require.Contains(t, normalized, "IS NOT TRUE",
		"must use IS NOT TRUE so accounts with NULL temp_unschedulable_until are not silently excluded by PG 3-valued logic")
	require.NotContains(t, normalized, "AND NOT (",
		"plain NOT (...) excludes NULL temp_unschedulable_until rows (the common healthy case)")
	require.Contains(t, normalized, "ORDER BY priority ASC, id ASC")
	require.NotContains(t, normalized, "credentials->>'expires_at'")
	require.NoError(t, mock.ExpectationsWereMet())
placeholder

type captureQuerySQL struct {
	db       *sql.DB
	captured *string
placeholder

func (c captureQuerySQL) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.db.ExecContext(ctx, query, args...)
placeholder

func (c captureQuerySQL) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if c.captured != nil {
		*c.captured = query
placeholder
	return c.db.QueryContext(ctx, query, args...)
placeholder

func normalizeSQLWhitespace(sql string) string {
	return strings.Join(regexp.MustCompile(`\s+`).Split(strings.TrimSpace(sql), -1), " ")
placeholder

type rowsAffectedResult int64

func (r rowsAffectedResult) LastInsertId() (int64, error) { return 0, nil placeholder
func (r rowsAffectedResult) RowsAffected() (int64, error) { return int64(r), nil placeholder

type recordingSQLExecutor struct {
	result      sql.Result
	err         error
	afterExec   func()
	execQueries []string
	execArgs    [][]any
placeholder

func (e *recordingSQLExecutor) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	e.execQueries = append(e.execQueries, query)
	e.execArgs = append(e.execArgs, append([]any(nil), args...))
	if e.err != nil {
		return nil, e.err
placeholder
	if e.afterExec != nil {
		e.afterExec()
placeholder
	return e.result, nil
placeholder

func (e *recordingSQLExecutor) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return nil, sql.ErrNoRows
placeholder
