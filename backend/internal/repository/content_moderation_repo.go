package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type contentModerationRepository struct {
	db *sql.DB
placeholder

func NewContentModerationRepository(db *sql.DB) service.ContentModerationRepository {
	return &contentModerationRepository{db: dbplaceholder
placeholder

func (r *contentModerationRepository) CreateLog(ctx context.Context, log *service.ContentModerationLog) error {
	if log == nil {
		return nil
placeholder
	categoryScores, err := json.Marshal(log.CategoryScores)
	if err != nil {
		return fmt.Errorf("marshal moderation category scores: %w", err)
placeholder
	thresholdSnapshot, err := json.Marshal(log.ThresholdSnapshot)
	if err != nil {
		return fmt.Errorf("marshal moderation thresholds: %w", err)
placeholder
	var userID any
	if log.UserID != nil {
		userID = *log.UserID
placeholder
	var apiKeyID any
	if log.APIKeyID != nil {
		apiKeyID = *log.APIKeyID
placeholder
	var groupID any
	if log.GroupID != nil {
		groupID = *log.GroupID
placeholder
	var latency any
	if log.UpstreamLatencyMS != nil {
		latency = *log.UpstreamLatencyMS
placeholder
	err = r.db.QueryRowContext(ctx, `
INSERT INTO content_moderation_logs (
    request_id, user_id, user_email, api_key_id, api_key_name, group_id, group_name,
    endpoint, provider, model, mode, action, flagged, highest_category, highest_score,
    category_scores, threshold_snapshot, input_excerpt, upstream_latency_ms, error,
    violation_count, auto_banned, email_sent, queue_delay_ms
) VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11, $12, $13, $14, $15,
    $16::jsonb, $17::jsonb, $18, $19, $20,
    $21, $22, $23, $24
) RETURNING id, created_at`,
		log.RequestID, userID, log.UserEmail, apiKeyID, log.APIKeyName, groupID, log.GroupName,
		log.Endpoint, log.Provider, log.Model, log.Mode, log.Action, log.Flagged, log.HighestCategory, log.HighestScore,
		string(categoryScores), string(thresholdSnapshot), log.InputExcerpt, latency, log.Error,
		log.ViolationCount, log.AutoBanned, log.EmailSent, nullableIntPtr(log.QueueDelayMS),
	).Scan(&log.ID, &log.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert content moderation log: %w", err)
placeholder
	return nil
placeholder

func (r *contentModerationRepository) ListLogs(ctx context.Context, filter service.ContentModerationLogFilter) ([]service.ContentModerationLog, *pagination.PaginationResult, error) {
	where, args := buildContentModerationLogWhere(filter)
	whereSQL := "WHERE " + strings.Join(where, " AND ")

	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM content_moderation_logs l "+whereSQL, args...).Scan(&total); err != nil {
		return nil, nil, fmt.Errorf("count content moderation logs: %w", err)
placeholder

	params := filter.Pagination
	if params.Page <= 0 {
		params.Page = 1
placeholder
	if params.PageSize <= 0 {
		params.PageSize = 20
placeholder
	if params.PageSize > 100 {
		params.PageSize = 100
placeholder
	queryArgs := append([]any{placeholder, args...)
	queryArgs = append(queryArgs, params.Limit(), params.Offset())
	rows, err := r.db.QueryContext(ctx, `
SELECT
    l.id, l.request_id, l.user_id, l.user_email, l.api_key_id, l.api_key_name, l.group_id, l.group_name,
    l.endpoint, l.provider, l.model, l.mode, l.action, l.flagged, l.highest_category, l.highest_score,
    l.category_scores, l.threshold_snapshot, l.input_excerpt, l.upstream_latency_ms, l.error,
    l.violation_count, l.auto_banned, l.email_sent, COALESCE(u.status, ''), l.queue_delay_ms, l.created_at
FROM content_moderation_logs l
LEFT JOIN users u ON u.id = l.user_id `+whereSQL+`
ORDER BY l.created_at DESC, l.id DESC
LIMIT $`+fmt.Sprint(len(queryArgs)-1)+` OFFSET $`+fmt.Sprint(len(queryArgs)),
		queryArgs...,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("list content moderation logs: %w", err)
placeholder
	defer func() { _ = rows.Close() placeholder()

	items := make([]service.ContentModerationLog, 0)
	for rows.Next() {
		var item service.ContentModerationLog
		var userID, apiKeyID, groupID, latency, queueDelay sql.NullInt64
		var scoresRaw, thresholdsRaw []byte
		if err := rows.Scan(
			&item.ID,
			&item.RequestID,
			&userID,
			&item.UserEmail,
			&apiKeyID,
			&item.APIKeyName,
			&groupID,
			&item.GroupName,
			&item.Endpoint,
			&item.Provider,
			&item.Model,
			&item.Mode,
			&item.Action,
			&item.Flagged,
			&item.HighestCategory,
			&item.HighestScore,
			&scoresRaw,
			&thresholdsRaw,
			&item.InputExcerpt,
			&latency,
			&item.Error,
			&item.ViolationCount,
			&item.AutoBanned,
			&item.EmailSent,
			&item.UserStatus,
			&queueDelay,
			&item.CreatedAt,
		); err != nil {
			return nil, nil, fmt.Errorf("scan content moderation log: %w", err)
	placeholder
		if userID.Valid {
			v := userID.Int64
			item.UserID = &v
	placeholder
		if apiKeyID.Valid {
			v := apiKeyID.Int64
			item.APIKeyID = &v
	placeholder
		if groupID.Valid {
			v := groupID.Int64
			item.GroupID = &v
	placeholder
		if latency.Valid {
			v := int(latency.Int64)
			item.UpstreamLatencyMS = &v
	placeholder
		if queueDelay.Valid {
			v := int(queueDelay.Int64)
			item.QueueDelayMS = &v
	placeholder
		item.CategoryScores = map[string]float64{placeholder
		_ = json.Unmarshal(scoresRaw, &item.CategoryScores)
		item.ThresholdSnapshot = map[string]float64{placeholder
		_ = json.Unmarshal(thresholdsRaw, &item.ThresholdSnapshot)
		items = append(items, item)
placeholder
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("iterate content moderation logs: %w", err)
placeholder
	return items, paginationResultFromTotal(total, params), nil
placeholder

func (r *contentModerationRepository) CountFlaggedByUserSince(ctx context.Context, userID int64, since time.Time, excludeCyberPolicy bool) (int, error) {
	if userID <= 0 {
		return 0, nil
placeholder
	// SQL 中的 'cyber_policy' 字面量须与 service.ContentModerationActionCyberPolicy 保持一致。
	var count int
	err := r.db.QueryRowContext(ctx, `
WITH last_auto_ban AS (
    SELECT MAX(created_at) AS at
    FROM content_moderation_logs
    WHERE user_id = $1 AND auto_banned = TRUE
)
SELECT COUNT(*)
FROM content_moderation_logs
WHERE user_id = $1
  AND flagged = TRUE
  AND action <> 'hash_block'
  AND ($3::bool IS FALSE OR action <> 'cyber_policy')
  AND created_at >= $2
  AND created_at > COALESCE((SELECT at FROM last_auto_ban), '-infinity'::timestamptz)
`, userID, since, excludeCyberPolicy).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count user content moderation flagged logs: %w", err)
placeholder
	return count, nil
placeholder

func (r *contentModerationRepository) UpdateLogEmailSent(ctx context.Context, id int64, sent bool) error {
	_, err := r.db.ExecContext(ctx, `UPDATE content_moderation_logs SET email_sent = $1 WHERE id = $2`, sent, id)
	if err != nil {
		return fmt.Errorf("update content moderation log email_sent: %w", err)
placeholder
	return nil
placeholder

func (r *contentModerationRepository) CleanupExpiredLogs(ctx context.Context, hitBefore time.Time, nonHitBefore time.Time) (*service.ContentModerationCleanupResult, error) {
	result := &service.ContentModerationCleanupResult{FinishedAt: time.Now()placeholder
	if r == nil || r.db == nil {
		return result, nil
placeholder
	hitExec, err := r.db.ExecContext(ctx, `
DELETE FROM content_moderation_logs
WHERE flagged = TRUE AND created_at < $1
`, hitBefore)
	if err != nil {
		return nil, fmt.Errorf("delete expired hit content moderation logs: %w", err)
placeholder
	result.DeletedHit, _ = hitExec.RowsAffected()

	nonHitExec, err := r.db.ExecContext(ctx, `
DELETE FROM content_moderation_logs
WHERE flagged = FALSE AND created_at < $1
`, nonHitBefore)
	if err != nil {
		return nil, fmt.Errorf("delete expired non-hit content moderation logs: %w", err)
placeholder
	result.DeletedNonHit, _ = nonHitExec.RowsAffected()

	result.FinishedAt = time.Now()
	return result, nil
placeholder

func nullableIntPtr(value *int) any {
	if value == nil {
		return nil
placeholder
	return *value
placeholder

func buildContentModerationLogWhere(filter service.ContentModerationLogFilter) ([]string, []any) {
	where := []string{"l.id IS NOT NULL"placeholder
	args := make([]any, 0)
	add := func(expr string, value any) {
		args = append(args, value)
		where = append(where, fmt.Sprintf(expr, len(args)))
placeholder
	switch strings.ToLower(strings.TrimSpace(filter.Result)) {
	case "hit", "flagged":
		where = append(where, "l.flagged = TRUE")
	case "blocked", "block":
		where = append(where, "l.action IN ('block', 'keyword_block', 'hash_block')")
	case "pass", "allow":
		where = append(where, "l.flagged = FALSE AND l.error = ''")
	case "error":
		where = append(where, "l.error <> ''")
placeholder
	if filter.GroupID != nil {
		add("l.group_id = $%d", *filter.GroupID)
placeholder
	if endpoint := strings.TrimSpace(filter.Endpoint); endpoint != "" {
		add("l.endpoint = $%d", endpoint)
placeholder
	if search := strings.TrimSpace(filter.Search); search != "" {
		like := "%" + search + "%"
		args = append(args, like, like, like, like, like)
		idx := len(args) - 4
		where = append(where, fmt.Sprintf("(l.request_id ILIKE $%d OR l.user_email ILIKE $%d OR l.api_key_name ILIKE $%d OR l.model ILIKE $%d OR l.input_excerpt ILIKE $%d)", idx, idx+1, idx+2, idx+3, idx+4))
placeholder
	if filter.From != nil && !filter.From.IsZero() {
		add("l.created_at >= $%d", *filter.From)
placeholder
	if filter.To != nil && !filter.To.IsZero() {
		add("l.created_at <= $%d", *filter.To)
placeholder
	return where, args
placeholder
