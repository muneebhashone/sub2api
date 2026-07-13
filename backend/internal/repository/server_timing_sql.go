package repository

import (
	"context"
	"database/sql/driver"
	"errors"
	"io"
	"reflect"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/servertiming"
)

type serverTimingConnector struct {
	base driver.Connector
placeholder

func newServerTimingConnector(base driver.Connector) driver.Connector {
	return &serverTimingConnector{base: baseplaceholder
placeholder

func (c *serverTimingConnector) Connect(ctx context.Context) (driver.Conn, error) {
	startedAt := time.Now()
	conn, err := c.base.Connect(ctx)
	servertiming.RecordInterval(ctx, servertiming.MetricDatabase, startedAt, time.Now())
	if err != nil {
		return nil, err
placeholder
	return &serverTimingConn{Conn: connplaceholder, nil
placeholder

func (c *serverTimingConnector) Driver() driver.Driver {
	return c.base.Driver()
placeholder

type serverTimingConn struct {
	driver.Conn
placeholder

func (c *serverTimingConn) Prepare(query string) (driver.Stmt, error) {
	stmt, err := c.Conn.Prepare(query)
	if err != nil {
		return nil, err
placeholder
	return &serverTimingStmt{Stmt: stmtplaceholder, nil
placeholder

func (c *serverTimingConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	startedAt := time.Now()
	var (
		stmt driver.Stmt
		err  error
	)
	if preparer, ok := c.Conn.(driver.ConnPrepareContext); ok {
		stmt, err = preparer.PrepareContext(ctx, query)
placeholder else {
		stmt, err = c.Conn.Prepare(query)
placeholder
	servertiming.Record(ctx, servertiming.MetricDatabase, startedAt, time.Now(), 1)
	if err != nil {
		return nil, err
placeholder
	return &serverTimingStmt{Stmt: stmtplaceholder, nil
placeholder

func (c *serverTimingConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	execer, ok := c.Conn.(driver.ExecerContext)
	if !ok {
		return nil, driver.ErrSkip
placeholder
	startedAt := time.Now()
	result, err := execer.ExecContext(ctx, query, args)
	servertiming.Record(ctx, servertiming.MetricDatabase, startedAt, time.Now(), 1)
	return result, err
placeholder

func (c *serverTimingConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	queryer, ok := c.Conn.(driver.QueryerContext)
	if !ok {
		return nil, driver.ErrSkip
placeholder
	startedAt := time.Now()
	rows, err := queryer.QueryContext(ctx, query, args)
	servertiming.Record(ctx, servertiming.MetricDatabase, startedAt, time.Now(), 1)
	if err != nil || rows == nil {
		return rows, err
placeholder
	return newServerTimingRows(ctx, rows), nil
placeholder

func (c *serverTimingConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	startedAt := time.Now()
	var (
		tx  driver.Tx
		err error
	)
	if beginner, ok := c.Conn.(driver.ConnBeginTx); ok {
		tx, err = beginner.BeginTx(ctx, opts)
placeholder else {
		if opts.Isolation != driver.IsolationLevel(0) {
			return nil, errors.New("driver does not support non-default isolation")
	placeholder
		if opts.ReadOnly {
			return nil, errors.New("driver does not support read-only transactions")
	placeholder
		tx, err = c.Conn.Begin()
placeholder
	servertiming.RecordInterval(ctx, servertiming.MetricDatabase, startedAt, time.Now())
	if err != nil || tx == nil {
		return tx, err
placeholder
	return &serverTimingTx{Tx: tx, ctx: ctxplaceholder, nil
placeholder

func (c *serverTimingConn) Ping(ctx context.Context) error {
	if pinger, ok := c.Conn.(driver.Pinger); ok {
		startedAt := time.Now()
		err := pinger.Ping(ctx)
		servertiming.RecordInterval(ctx, servertiming.MetricDatabase, startedAt, time.Now())
		return err
placeholder
	return nil
placeholder

func (c *serverTimingConn) ResetSession(ctx context.Context) error {
	if resetter, ok := c.Conn.(driver.SessionResetter); ok {
		startedAt := time.Now()
		err := resetter.ResetSession(ctx)
		servertiming.RecordInterval(ctx, servertiming.MetricDatabase, startedAt, time.Now())
		return err
placeholder
	return nil
placeholder

func (c *serverTimingConn) IsValid() bool {
	if validator, ok := c.Conn.(driver.Validator); ok {
		return validator.IsValid()
placeholder
	return true
placeholder

func (c *serverTimingConn) CheckNamedValue(value *driver.NamedValue) error {
	if checker, ok := c.Conn.(driver.NamedValueChecker); ok {
		return checker.CheckNamedValue(value)
placeholder
	return driver.ErrSkip
placeholder

type serverTimingStmt struct {
	driver.Stmt
placeholder

func (s *serverTimingStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	startedAt := time.Now()
	var (
		result driver.Result
		err    error
	)
	if execer, ok := s.Stmt.(driver.StmtExecContext); ok {
		result, err = execer.ExecContext(ctx, args)
placeholder else {
		var values []driver.Value
		values, err = namedValues(args)
		if err == nil {
			result, err = s.Stmt.Exec(values)
	placeholder
placeholder
	servertiming.Record(ctx, servertiming.MetricDatabase, startedAt, time.Now(), 1)
	return result, err
placeholder

func (s *serverTimingStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	startedAt := time.Now()
	var (
		rows driver.Rows
		err  error
	)
	if queryer, ok := s.Stmt.(driver.StmtQueryContext); ok {
		rows, err = queryer.QueryContext(ctx, args)
placeholder else {
		var values []driver.Value
		values, err = namedValues(args)
		if err == nil {
			rows, err = s.Stmt.Query(values)
	placeholder
placeholder
	servertiming.Record(ctx, servertiming.MetricDatabase, startedAt, time.Now(), 1)
	if err != nil || rows == nil {
		return rows, err
placeholder
	return newServerTimingRows(ctx, rows), nil
placeholder

func (s *serverTimingStmt) CheckNamedValue(value *driver.NamedValue) error {
	if checker, ok := s.Stmt.(driver.NamedValueChecker); ok {
		return checker.CheckNamedValue(value)
placeholder
	return driver.ErrSkip
placeholder

func (s *serverTimingStmt) ColumnConverter(index int) driver.ValueConverter {
	if converter, ok := s.Stmt.(driver.ColumnConverter); ok {
		return converter.ColumnConverter(index)
placeholder
	return driver.DefaultParameterConverter
placeholder

func namedValues(args []driver.NamedValue) ([]driver.Value, error) {
	values := make([]driver.Value, len(args))
	for i, arg := range args {
		if arg.Name != "" {
			return nil, errors.New("named parameters are not supported")
	placeholder
		values[i] = arg.Value
placeholder
	return values, nil
placeholder

type serverTimingRows struct {
	driver.Rows
	ctx context.Context
placeholder

func newServerTimingRows(ctx context.Context, rows driver.Rows) *serverTimingRows {
	return &serverTimingRows{Rows: rows, ctx: ctxplaceholder
placeholder

func (r *serverTimingRows) Close() error {
	startedAt := time.Now()
	err := r.Rows.Close()
	servertiming.RecordInterval(r.ctx, servertiming.MetricDatabase, startedAt, time.Now())
	return err
placeholder

func (r *serverTimingRows) Next(dest []driver.Value) error {
	startedAt := time.Now()
	err := r.Rows.Next(dest)
	servertiming.RecordInterval(r.ctx, servertiming.MetricDatabase, startedAt, time.Now())
	return err
placeholder

func (r *serverTimingRows) HasNextResultSet() bool {
	if rows, ok := r.Rows.(driver.RowsNextResultSet); ok {
		return rows.HasNextResultSet()
placeholder
	return false
placeholder

func (r *serverTimingRows) NextResultSet() error {
	rows, ok := r.Rows.(driver.RowsNextResultSet)
	if !ok {
		return io.EOF
placeholder
	startedAt := time.Now()
	err := rows.NextResultSet()
	servertiming.RecordInterval(r.ctx, servertiming.MetricDatabase, startedAt, time.Now())
	return err
placeholder

func (r *serverTimingRows) ColumnTypeScanType(index int) reflect.Type {
	if rows, ok := r.Rows.(driver.RowsColumnTypeScanType); ok {
		return rows.ColumnTypeScanType(index)
placeholder
	return reflect.TypeOf(new(any)).Elem()
placeholder

func (r *serverTimingRows) ColumnTypeDatabaseTypeName(index int) string {
	if rows, ok := r.Rows.(driver.RowsColumnTypeDatabaseTypeName); ok {
		return rows.ColumnTypeDatabaseTypeName(index)
placeholder
	return ""
placeholder

func (r *serverTimingRows) ColumnTypeLength(index int) (int64, bool) {
	if rows, ok := r.Rows.(driver.RowsColumnTypeLength); ok {
		return rows.ColumnTypeLength(index)
placeholder
	return 0, false
placeholder

func (r *serverTimingRows) ColumnTypeNullable(index int) (bool, bool) {
	if rows, ok := r.Rows.(driver.RowsColumnTypeNullable); ok {
		return rows.ColumnTypeNullable(index)
placeholder
	return false, false
placeholder

func (r *serverTimingRows) ColumnTypePrecisionScale(index int) (int64, int64, bool) {
	if rows, ok := r.Rows.(driver.RowsColumnTypePrecisionScale); ok {
		return rows.ColumnTypePrecisionScale(index)
placeholder
	return 0, 0, false
placeholder

type serverTimingTx struct {
	driver.Tx
	ctx context.Context
placeholder

func (t *serverTimingTx) Commit() error {
	startedAt := time.Now()
	err := t.Tx.Commit()
	servertiming.RecordInterval(t.ctx, servertiming.MetricDatabase, startedAt, time.Now())
	return err
placeholder

func (t *serverTimingTx) Rollback() error {
	startedAt := time.Now()
	err := t.Tx.Rollback()
	servertiming.RecordInterval(t.ctx, servertiming.MetricDatabase, startedAt, time.Now())
	return err
placeholder
