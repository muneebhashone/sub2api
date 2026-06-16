package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"sync"
	"testing"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	_ "github.com/Wei-Shaw/sub2api/ent/runtime"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
)

const parameterLimitTestDriverName = "sub2api_param_limit_test"

var registerParameterLimitTestDriverOnce sync.Once

func TestAccountsToService_LargeActiveAccountSetDoesNotExceedPostgresParameterLimit(t *testing.T) {
	repo := newParameterLimitAccountRepo(t)

	accounts := make([]*dbent.Account, 0, 65536)
	for i := range 65536 {
		accounts = append(accounts, &dbent.Account{
			ID:          int64(i + 1),
			Name:        "large-active",
			Platform:    service.PlatformOpenAI,
			Type:        service.AccountTypeOAuth,
	placeholderplaceholder,
			Extra:       map[string]any{placeholder,
			Status:      service.StatusActive,
			Schedulable: true,
	placeholder)
placeholder

	got, err := repo.accountsToService(context.Background(), accounts)
placeholder
	require.Len(t, got, len(accounts))
placeholder

func newParameterLimitAccountRepo(t *testing.T) *accountRepository {
placeholder

	registerParameterLimitTestDriverOnce.Do(func() {
		sql.Register(parameterLimitTestDriverName, parameterLimitDriver{placeholder)
placeholder)

	db, err := sql.Open(parameterLimitTestDriverName, "")
placeholder
	t.Cleanup(func() { _ = db.Close() placeholder)

	drv := entsql.OpenDB(dialect.Postgres, db)
	client := dbent.NewClient(dbent.Driver(drv))
	t.Cleanup(func() { _ = client.Close() placeholder)

	return newAccountRepositoryWithSQL(client, nil, nil)
placeholder

type parameterLimitDriver struct{placeholder

func (parameterLimitDriver) Open(string) (driver.Conn, error) {
	return parameterLimitConn{placeholder, nil
placeholder

type parameterLimitConn struct{placeholder

func (parameterLimitConn) Prepare(query string) (driver.Stmt, error) {
	return parameterLimitStmt{query: queryplaceholder, nil
placeholder

func (parameterLimitConn) Close() error {
	return nil
placeholder

func (parameterLimitConn) Begin() (driver.Tx, error) {
	return parameterLimitTx{placeholder, nil
placeholder

func (parameterLimitConn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return queryWithParameterLimit(query, args)
placeholder

type parameterLimitStmt struct {
	query string
placeholder

func (s parameterLimitStmt) Close() error {
	return nil
placeholder

func (s parameterLimitStmt) NumInput() int {
	return -1
placeholder

func (s parameterLimitStmt) Exec(args []driver.Value) (driver.Result, error) {
	return driver.RowsAffected(0), parameterLimitError(len(args))
placeholder

func (s parameterLimitStmt) Query(args []driver.Value) (driver.Rows, error) {
	namedArgs := make([]driver.NamedValue, len(args))
	for i, arg := range args {
		namedArgs[i] = driver.NamedValue{Ordinal: i + 1, Value: argplaceholder
placeholder
	return queryWithParameterLimit(s.query, namedArgs)
placeholder

type parameterLimitTx struct{placeholder

func (parameterLimitTx) Commit() error {
	return nil
placeholder

func (parameterLimitTx) Rollback() error {
	return nil
placeholder

func queryWithParameterLimit(query string, args []driver.NamedValue) (driver.Rows, error) {
	if err := parameterLimitError(len(args)); err != nil {
		return nil, err
placeholder
	return parameterLimitRows{columns: columnsForParameterLimitQuery(query)placeholder, nil
placeholder

func parameterLimitError(paramCount int) error {
	if paramCount <= 65535 {
		return nil
placeholder
	return fmt.Errorf("pq: got %d parameters but PostgreSQL only supports 65535 parameters", paramCount)
placeholder

func columnsForParameterLimitQuery(query string) []string {
	if query == "" {
		return nil
placeholder
	return []string{"account_id", "group_id", "priority", "created_at"placeholder
placeholder

type parameterLimitRows struct {
	columns []string
placeholder

func (r parameterLimitRows) Columns() []string {
	return r.columns
placeholder

func (parameterLimitRows) Close() error {
	return nil
placeholder

func (parameterLimitRows) Next([]driver.Value) error {
	return io.EOF
placeholder
