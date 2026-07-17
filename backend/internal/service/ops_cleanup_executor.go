package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

const (
	opsCleanupDefaultSchedule  = "0 2 * * *"
	opsCleanupBatchSize        = 5000
	opsCleanupCronStopTimeout  = 3 * time.Second
	opsCleanupRunTimeout       = 30 * time.Minute
	opsCleanupHeartbeatTimeout = 2 * time.Second
)

type opsCleanupTarget struct {
	retentionDays int
	table         string
	timeCol       string
	castDate      bool
	counter       *int64
placeholder

type opsCleanupDeletedCounts struct {
	errorLogs      int64
	ingressRejects int64
	alertEvents    int64
	systemLogs     int64
	logAudits      int64
	systemMetrics  int64
	hourlyPreagg   int64
	dailyPreagg    int64
placeholder

func (c opsCleanupDeletedCounts) String() string {
	return fmt.Sprintf(
		"error_logs=%d ingress_rejects=%d alert_events=%d system_logs=%d log_audits=%d system_metrics=%d hourly_preagg=%d daily_preagg=%d",
		c.errorLogs,
		c.ingressRejects,
		c.alertEvents,
		c.systemLogs,
		c.logAudits,
		c.systemMetrics,
		c.hourlyPreagg,
		c.dailyPreagg,
	)
placeholder

// opsCleanupPlan 把"保留天数"翻译成具体的清理动作。
//   - days < 0  → 跳过该项清理（ok=false），保留兼容老数据
//   - days == 0 → TRUNCATE TABLE（O(1) 全清），truncate=true
//   - days > 0  → 批量 DELETE 早于 now-N天 的行，cutoff = now - N 天
func opsCleanupPlan(now time.Time, days int) (cutoff time.Time, truncate, ok bool) {
	if days < 0 {
		return time.Time{placeholder, false, false
placeholder
	if days == 0 {
		return time.Time{placeholder, true, true
placeholder
	return now.AddDate(0, 0, -days), false, true
placeholder

func opsCleanupRunOne(
	ctx context.Context,
	db *sql.DB,
	truncate bool,
	cutoff time.Time,
	table, timeCol string,
	castDate bool,
	batchSize int,
) (int64, error) {
	if truncate {
		return truncateOpsTable(ctx, db, table)
placeholder
	return deleteOldRowsByID(ctx, db, table, timeCol, cutoff, batchSize, castDate)
placeholder

func deleteOldRowsByID(
	ctx context.Context,
	db *sql.DB,
	table string,
	timeColumn string,
	cutoff time.Time,
	batchSize int,
	castCutoffToDate bool,
) (int64, error) {
	if db == nil {
		return 0, nil
placeholder
	if batchSize <= 0 {
		batchSize = opsCleanupBatchSize
placeholder

	where := fmt.Sprintf("%s < $1", timeColumn)
	if castCutoffToDate {
		where = fmt.Sprintf("%s < $1::date", timeColumn)
placeholder

	q := fmt.Sprintf(`
WITH batch AS (
  SELECT id FROM %s
  WHERE %s
  ORDER BY id
  LIMIT $2
)
DELETE FROM %s
WHERE id IN (SELECT id FROM batch)
`, table, where, table)

	var total int64
	for {
		res, err := db.ExecContext(ctx, q, cutoff, batchSize)
		if err != nil {
			if isMissingRelationError(err) {
				return total, nil
		placeholder
			return total, err
	placeholder
		affected, err := res.RowsAffected()
		if err != nil {
			return total, err
	placeholder
		total += affected
		if affected == 0 {
			break
	placeholder
placeholder
	return total, nil
placeholder

// truncateOpsTable 用 TRUNCATE TABLE 清空指定表，先 SELECT COUNT(*) 取得清空前行数用于 heartbeat。
func truncateOpsTable(ctx context.Context, db *sql.DB, table string) (int64, error) {
	if db == nil {
		return 0, nil
placeholder
	var count int64
	if err := db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count); err != nil {
		if isMissingRelationError(err) {
			return 0, nil
	placeholder
		return 0, fmt.Errorf("count %s: %w", table, err)
placeholder
	if count == 0 {
		return 0, nil
placeholder
	if _, err := db.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s", table)); err != nil {
		if isMissingRelationError(err) {
			return 0, nil
	placeholder
		return 0, fmt.Errorf("truncate %s: %w", table, err)
placeholder
	return count, nil
placeholder

func isMissingRelationError(err error) bool {
	if err == nil {
		return false
placeholder
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "does not exist") && strings.Contains(s, "relation")
placeholder
