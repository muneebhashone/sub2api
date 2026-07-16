package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"
)

// auditLogRepository 审计日志仓储（raw SQL，append-only）。
// 刻意不实现单条删除：审计日志只允许追加、按保留期批量清理、以及带 2FA 的全量清空。
type auditLogRepository struct {
	db *sql.DB
placeholder

// NewAuditLogRepository 创建审计日志仓储。
func NewAuditLogRepository(db *sql.DB) service.AuditLogRepository {
	return &auditLogRepository{db: dbplaceholder
placeholder

const auditLogInsertColumns = `created_at, actor_user_id, actor_email, actor_role, auth_method,
credential_masked, action, method, path, request_id, client_ip, user_agent,
request_body, status_code, latency_ms, extra`

func auditLogInsertValues(log *service.AuditLog) []any {
	createdAt := log.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
placeholder
	extraJSON := "{placeholder"
	if len(log.Extra) > 0 {
		if encoded, err := json.Marshal(log.Extra); err == nil {
			extraJSON = string(encoded)
	placeholder
placeholder
	return []any{
		createdAt.UTC(),
		nullInt64Ptr(log.ActorUserID),
		truncateString(log.ActorEmail, 255),
		truncateString(log.ActorRole, 32),
		truncateString(log.AuthMethod, 32),
		truncateString(log.CredentialMasked, 160),
		truncateString(log.Action, 128),
		truncateString(log.Method, 16),
		truncateString(log.Path, 512),
		truncateString(log.RequestID, 64),
		truncateString(log.ClientIP, 64),
		truncateString(log.UserAgent, 512),
		log.RequestBody,
		log.StatusCode,
		log.LatencyMs,
		extraJSON,
placeholder
placeholder

func (r *auditLogRepository) BatchInsert(ctx context.Context, logs []*service.AuditLog) (int64, error) {
	if r == nil || r.db == nil {
		return 0, fmt.Errorf("nil audit log repository")
placeholder
	if len(logs) == 0 {
		return 0, nil
placeholder

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
placeholder
	stmt, err := tx.PrepareContext(ctx, pq.CopyIn(
		"audit_logs",
		"created_at", "actor_user_id", "actor_email", "actor_role", "auth_method",
		"credential_masked", "action", "method", "path", "request_id", "client_ip", "user_agent",
		"request_body", "status_code", "latency_ms", "extra",
	))
	if err != nil {
		_ = tx.Rollback()
		return 0, err
placeholder

	var inserted int64
	for _, log := range logs {
		if log == nil {
			continue
	placeholder
		if _, err := stmt.ExecContext(ctx, auditLogInsertValues(log)...); err != nil {
			_ = stmt.Close()
			_ = tx.Rollback()
			return inserted, err
	placeholder
		inserted++
placeholder

	if _, err := stmt.ExecContext(ctx); err != nil {
		_ = stmt.Close()
		_ = tx.Rollback()
		return inserted, err
placeholder
	if err := stmt.Close(); err != nil {
		_ = tx.Rollback()
		return inserted, err
placeholder
	if err := tx.Commit(); err != nil {
		return inserted, err
placeholder
	return inserted, nil
placeholder

func (r *auditLogRepository) Insert(ctx context.Context, log *service.AuditLog) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("nil audit log repository")
placeholder
	if log == nil {
		return fmt.Errorf("nil audit log")
placeholder
	query := `INSERT INTO audit_logs (` + auditLogInsertColumns + `)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`
	_, err := r.db.ExecContext(ctx, query, auditLogInsertValues(log)...)
	return err
placeholder

func buildAuditLogsWhere(filter *service.AuditLogFilter) (string, []any) {
	clauses := make([]string, 0, 10)
	args := make([]any, 0, 10)
	clauses = append(clauses, "1=1")

	if filter.StartTime != nil {
		args = append(args, filter.StartTime.UTC())
		clauses = append(clauses, "l.created_at >= $"+itoa(len(args)))
placeholder
	if filter.EndTime != nil {
		args = append(args, filter.EndTime.UTC())
		clauses = append(clauses, "l.created_at <= $"+itoa(len(args)))
placeholder
	if filter.ActorUserID != nil {
		args = append(args, *filter.ActorUserID)
		clauses = append(clauses, "l.actor_user_id = $"+itoa(len(args)))
placeholder
	if v := strings.TrimSpace(filter.ActorEmail); v != "" {
		args = append(args, "%"+escapeLikePattern(v)+"%")
		clauses = append(clauses, "l.actor_email ILIKE $"+itoa(len(args)))
placeholder
	if v := strings.TrimSpace(filter.AuthMethod); v != "" {
		args = append(args, v)
		clauses = append(clauses, "l.auth_method = $"+itoa(len(args)))
placeholder
	if v := strings.TrimSpace(filter.Action); v != "" {
		args = append(args, "%"+escapeLikePattern(v)+"%")
		clauses = append(clauses, "l.action ILIKE $"+itoa(len(args)))
placeholder
	if v := strings.TrimSpace(filter.Method); v != "" {
		args = append(args, strings.ToUpper(v))
		clauses = append(clauses, "l.method = $"+itoa(len(args)))
placeholder
	if v := strings.TrimSpace(filter.ClientIP); v != "" {
		args = append(args, v)
		clauses = append(clauses, "l.client_ip = $"+itoa(len(args)))
placeholder
	if filter.Success != nil {
		if *filter.Success {
			clauses = append(clauses, "l.status_code < 400")
	placeholder else {
			clauses = append(clauses, "l.status_code >= 400")
	placeholder
placeholder
	if v := strings.TrimSpace(filter.Query); v != "" {
		args = append(args, "%"+escapeLikePattern(v)+"%")
		idx := itoa(len(args))
		clauses = append(clauses, "(l.path ILIKE $"+idx+" OR l.action ILIKE $"+idx+" OR l.actor_email ILIKE $"+idx+")")
placeholder

	return "WHERE " + strings.Join(clauses, " AND "), args
placeholder

const auditLogSelectColumns = `
  l.id,
  l.created_at,
  l.actor_user_id,
  COALESCE(l.actor_email, ''),
  COALESCE(l.actor_role, ''),
  COALESCE(l.auth_method, ''),
  COALESCE(l.credential_masked, ''),
  COALESCE(l.action, ''),
  COALESCE(l.method, ''),
  COALESCE(l.path, ''),
  COALESCE(l.request_id, ''),
  COALESCE(l.client_ip, ''),
  COALESCE(l.user_agent, ''),
  COALESCE(l.request_body, ''),
  l.status_code,
  l.latency_ms,
  COALESCE(l.extra::text, '{placeholder')`

func scanAuditLogRow(scan func(dest ...any) error) (*service.AuditLog, error) {
	item := &service.AuditLog{placeholder
	var actorUserID sql.NullInt64
	var extraRaw string
	if err := scan(
		&item.ID,
		&item.CreatedAt,
		&actorUserID,
		&item.ActorEmail,
		&item.ActorRole,
		&item.AuthMethod,
		&item.CredentialMasked,
		&item.Action,
		&item.Method,
		&item.Path,
		&item.RequestID,
		&item.ClientIP,
		&item.UserAgent,
		&item.RequestBody,
		&item.StatusCode,
		&item.LatencyMs,
		&extraRaw,
	); err != nil {
		return nil, err
placeholder
	if actorUserID.Valid {
		v := actorUserID.Int64
		item.ActorUserID = &v
placeholder
	extraRaw = strings.TrimSpace(extraRaw)
	if extraRaw != "" && extraRaw != "null" && extraRaw != "{placeholder" {
		extra := make(map[string]any)
		if err := json.Unmarshal([]byte(extraRaw), &extra); err == nil {
			item.Extra = extra
	placeholder
placeholder
	return item, nil
placeholder

func (r *auditLogRepository) List(ctx context.Context, filter *service.AuditLogFilter) (*service.AuditLogList, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("nil audit log repository")
placeholder
	if filter == nil {
		filter = &service.AuditLogFilter{placeholder
placeholder

	page := filter.Page
	if page <= 0 {
		page = 1
placeholder
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 50
placeholder
	if pageSize > 200 {
		pageSize = 200
placeholder

	where, args := buildAuditLogsWhere(filter)
	countSQL := "SELECT COUNT(*) FROM audit_logs l " + where
	var total int
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, err
placeholder

	offset := (page - 1) * pageSize
	argsWithLimit := append(args, pageSize, offset)
	query := "SELECT" + auditLogSelectColumns + "\nFROM audit_logs l\n" + where + `
ORDER BY l.created_at DESC, l.id DESC
LIMIT $` + itoa(len(args)+1) + ` OFFSET $` + itoa(len(args)+2)

	rows, err := r.db.QueryContext(ctx, query, argsWithLimit...)
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = rows.Close() placeholder()

	logs := make([]*service.AuditLog, 0, pageSize)
	for rows.Next() {
		item, err := scanAuditLogRow(rows.Scan)
		if err != nil {
			return nil, err
	placeholder
		// 列表页不返回 body，降低载荷；详情接口返回完整记录。
		item.RequestBody = ""
		logs = append(logs, item)
placeholder
	if err := rows.Err(); err != nil {
		return nil, err
placeholder

	return &service.AuditLogList{
		Logs:     logs,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
placeholder, nil
placeholder

func (r *auditLogRepository) GetByID(ctx context.Context, id int64) (*service.AuditLog, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("nil audit log repository")
placeholder
	query := "SELECT" + auditLogSelectColumns + "\nFROM audit_logs l WHERE l.id = $1"
	row := r.db.QueryRowContext(ctx, query, id)
	item, err := scanAuditLogRow(row.Scan)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrAuditLogNotFound
	placeholder
		return nil, err
placeholder
	return item, nil
placeholder

func (r *auditLogRepository) Count(ctx context.Context) (int64, error) {
	if r == nil || r.db == nil {
		return 0, fmt.Errorf("nil audit log repository")
placeholder
	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM audit_logs").Scan(&total); err != nil {
		return 0, err
placeholder
	return total, nil
placeholder

func (r *auditLogRepository) TruncateAll(ctx context.Context) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("nil audit log repository")
placeholder
	_, err := r.db.ExecContext(ctx, "TRUNCATE TABLE audit_logs")
	return err
placeholder

func (r *auditLogRepository) DeleteBefore(ctx context.Context, cutoff time.Time, batchSize int) (int64, error) {
	if r == nil || r.db == nil {
		return 0, fmt.Errorf("nil audit log repository")
placeholder
	if batchSize <= 0 {
		batchSize = 5000
placeholder
	res, err := r.db.ExecContext(ctx, `
WITH batch AS (
  SELECT id FROM audit_logs WHERE created_at < $1 ORDER BY id LIMIT $2
)
DELETE FROM audit_logs WHERE id IN (SELECT id FROM batch)`, cutoff.UTC(), batchSize)
	if err != nil {
		return 0, err
placeholder
	return res.RowsAffected()
placeholder

func nullInt64Ptr(v *int64) any {
	if v == nil || *v <= 0 {
		return nil
placeholder
	return *v
placeholder

func truncateString(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
placeholder
	// 按字节截断可能切断多字节字符，按 rune 处理。
	runes := []rune(s)
	for len(string(runes)) > max && len(runes) > 0 {
		runes = runes[:len(runes)-1]
placeholder
	return string(runes)
placeholder
