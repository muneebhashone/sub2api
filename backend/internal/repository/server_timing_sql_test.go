package repository

import (
	"context"
	"database/sql/driver"
	"io"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/servertiming"
)

const fakeDriverDelay = 2 * time.Millisecond

type timingFakeDriver struct{placeholder

func (timingFakeDriver) Open(string) (driver.Conn, error) { return newTimingFakeConn(), nil placeholder

type timingFakeConnector struct {
	conn driver.Conn
placeholder

func (c timingFakeConnector) Connect(context.Context) (driver.Conn, error) {
	time.Sleep(fakeDriverDelay)
	return c.conn, nil
placeholder

func (timingFakeConnector) Driver() driver.Driver { return timingFakeDriver{placeholder placeholder

type timingFakeConn struct{placeholder

func newTimingFakeConn() *timingFakeConn { return &timingFakeConn{placeholder placeholder

func (c *timingFakeConn) Prepare(string) (driver.Stmt, error) {
	time.Sleep(fakeDriverDelay)
	return &timingFakeStmt{placeholder, nil
placeholder

func (c *timingFakeConn) PrepareContext(context.Context, string) (driver.Stmt, error) {
	time.Sleep(fakeDriverDelay)
	return &timingFakeStmt{placeholder, nil
placeholder

func (c *timingFakeConn) Close() error { return nil placeholder

func (c *timingFakeConn) Begin() (driver.Tx, error) {
	time.Sleep(fakeDriverDelay)
	return &timingFakeTx{placeholder, nil
placeholder

func (c *timingFakeConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	time.Sleep(fakeDriverDelay)
	return &timingFakeTx{placeholder, nil
placeholder

func (c *timingFakeConn) ExecContext(context.Context, string, []driver.NamedValue) (driver.Result, error) {
	time.Sleep(fakeDriverDelay)
	return driver.RowsAffected(1), nil
placeholder

func (c *timingFakeConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	time.Sleep(fakeDriverDelay)
	return &timingFakeRows{values: [][]driver.Value{{"value"placeholderplaceholderplaceholder, nil
placeholder

func (c *timingFakeConn) Ping(context.Context) error {
	time.Sleep(fakeDriverDelay)
	return nil
placeholder

func (c *timingFakeConn) ResetSession(context.Context) error {
	time.Sleep(fakeDriverDelay)
	return nil
placeholder

type timingFakeStmt struct{placeholder

func (s *timingFakeStmt) Close() error  { return nil placeholder
func (s *timingFakeStmt) NumInput() int { return -1 placeholder

func (s *timingFakeStmt) Exec([]driver.Value) (driver.Result, error) {
	time.Sleep(fakeDriverDelay)
	return driver.RowsAffected(1), nil
placeholder

func (s *timingFakeStmt) Query([]driver.Value) (driver.Rows, error) {
	time.Sleep(fakeDriverDelay)
	return &timingFakeRows{values: [][]driver.Value{{"value"placeholderplaceholderplaceholder, nil
placeholder

func (s *timingFakeStmt) ExecContext(context.Context, []driver.NamedValue) (driver.Result, error) {
	time.Sleep(fakeDriverDelay)
	return driver.RowsAffected(1), nil
placeholder

func (s *timingFakeStmt) QueryContext(context.Context, []driver.NamedValue) (driver.Rows, error) {
	time.Sleep(fakeDriverDelay)
	return &timingFakeRows{values: [][]driver.Value{{"value"placeholderplaceholderplaceholder, nil
placeholder

type timingFakeRows struct {
	values [][]driver.Value
	index  int
placeholder

func (r *timingFakeRows) Columns() []string { return []string{"value"placeholder placeholder

func (r *timingFakeRows) Close() error {
	time.Sleep(fakeDriverDelay)
	return nil
placeholder

func (r *timingFakeRows) Next(dest []driver.Value) error {
	time.Sleep(fakeDriverDelay)
	if r.index >= len(r.values) {
		return io.EOF
placeholder
	copy(dest, r.values[r.index])
	r.index++
	return nil
placeholder

type timingFakeTx struct{placeholder

func (t *timingFakeTx) Commit() error {
	time.Sleep(fakeDriverDelay)
	return nil
placeholder

func (t *timingFakeTx) Rollback() error {
	time.Sleep(fakeDriverDelay)
	return nil
placeholder

func metricDuration(t *testing.T, header, metric string) float64 {
placeholder
	re := regexp.MustCompile(`(?:^|, )` + regexp.QuoteMeta(metric) + `;dur=([0-9]+(?:\.[0-9]+)?)`)
	match := re.FindStringSubmatch(header)
	if len(match) != 2 {
		t.Fatalf("metric %q missing from header %q", metric, header)
placeholder
	value, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		t.Fatalf("parse %s duration: %v", metric, err)
placeholder
	return value
placeholder

func TestServerTimingConnectorRecordsDriverCallsWithoutRowLifetime(t *testing.T) {
	startedAt := time.Now()
	collector := servertiming.New(startedAt)
	ctx := servertiming.WithCollector(context.Background(), collector)

	wrapped := newServerTimingConnector(timingFakeConnector{conn: newTimingFakeConn()placeholder)
	rawConn, err := wrapped.Connect(ctx)
	if err != nil {
		t.Fatal(err)
placeholder
	conn, ok := rawConn.(*serverTimingConn)
	if !ok {
		t.Fatalf("Connect() returned %T, want *serverTimingConn", rawConn)
placeholder

	if _, err := conn.ExecContext(ctx, "sensitive update", nil); err != nil {
		t.Fatal(err)
placeholder
	rows, err := conn.QueryContext(ctx, "sensitive select", nil)
	if err != nil {
		t.Fatal(err)
placeholder
	values := make([]driver.Value, 1)
	if err := rows.Next(values); err != nil {
		t.Fatal(err)
placeholder

	// Application work between row reads must remain app time.
	time.Sleep(30 * time.Millisecond)
	if err := rows.Next(values); err != io.EOF {
		t.Fatalf("rows.Next() = %v, want EOF", err)
placeholder
	if err := rows.Close(); err != nil {
		t.Fatal(err)
placeholder

	header := collector.HeaderValue(time.Now(), "bypass")
	if !strings.Contains(header, `queries=2`) {
		t.Fatalf("header %q does not report two SQL operations", header)
placeholder
	if strings.Contains(header, "sensitive") {
		t.Fatalf("SQL text leaked into header: %q", header)
placeholder
	if app, db := metricDuration(t, header, "app"), metricDuration(t, header, "db"); app <= db {
		t.Fatalf("row processing gap was counted as DB time: app=%.1fms db=%.1fms header=%q", app, db, header)
placeholder
placeholder

func TestServerTimingPreparedStatementsAndTransactions(t *testing.T) {
	collector := servertiming.New(time.Now())
	ctx := servertiming.WithCollector(context.Background(), collector)
	conn := &serverTimingConn{Conn: newTimingFakeConn()placeholder

	stmt, err := conn.PrepareContext(ctx, "prepare sensitive statement")
	if err != nil {
		t.Fatal(err)
placeholder
	timedStmt, ok := stmt.(*serverTimingStmt)
	if !ok {
		t.Fatalf("PrepareContext() returned %T, want *serverTimingStmt", stmt)
placeholder
	if _, err := timedStmt.ExecContext(ctx, nil); err != nil {
		t.Fatal(err)
placeholder
	rows, err := timedStmt.QueryContext(ctx, nil)
	if err != nil {
		t.Fatal(err)
placeholder
	if err := rows.Close(); err != nil {
		t.Fatal(err)
placeholder

	tx, err := conn.BeginTx(ctx, driver.TxOptions{placeholder)
	if err != nil {
		t.Fatal(err)
placeholder
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
placeholder
	if err := conn.Ping(ctx); err != nil {
		t.Fatal(err)
placeholder
	if err := conn.ResetSession(ctx); err != nil {
		t.Fatal(err)
placeholder

	header := collector.HeaderValue(time.Now(), "bypass")
	if !strings.Contains(header, `queries=3`) {
		t.Fatalf("header %q does not report prepare, exec, and query operations", header)
placeholder
	if metricDuration(t, header, "db") <= 0 {
		t.Fatalf("DB duration was not recorded: %q", header)
placeholder
placeholder

func TestNamedValuesRejectNamedParameters(t *testing.T) {
	if _, err := namedValues([]driver.NamedValue{{Name: "secret", Value: 1placeholderplaceholder); err == nil {
		t.Fatal("namedValues accepted a named parameter")
placeholder
	values, err := namedValues([]driver.NamedValue{{Ordinal: 1, Value: "value"placeholderplaceholder)
	if err != nil {
		t.Fatal(err)
placeholder
	if len(values) != 1 || values[0] != "value" {
		t.Fatalf("namedValues() = %#v", values)
placeholder
placeholder
