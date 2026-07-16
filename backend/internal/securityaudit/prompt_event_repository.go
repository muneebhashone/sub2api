package securityaudit

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

type EventFilter struct {
	Decision   string     `json:"decision,omitempty"`
	RiskLevel  string     `json:"risk_level,omitempty"`
	Endpoint   string     `json:"endpoint,omitempty"`
	GroupID    *int64     `json:"group_id,omitempty"`
	UserID     *int64     `json:"user_id,omitempty"`
	APIKeyID   *int64     `json:"api_key_id,omitempty"`
	RequestID  string     `json:"request_id,omitempty"`
	PromptHash string     `json:"prompt_hash,omitempty"`
	Keyword    string     `json:"keyword,omitempty"`
	StartAt    *time.Time `json:"start_at,omitempty"`
	EndAt      *time.Time `json:"end_at,omitempty"`
placeholder

type EventPage struct {
	Items    []*Event `json:"items"`
	Total    int64    `json:"total"`
	Page     int      `json:"page"`
	PageSize int      `json:"page_size"`
	Pages    int      `json:"pages"`
placeholder

type DeletePreview struct {
	MatchedCount      int64       `json:"matched_count"`
	FilterSummary     EventFilter `json:"filter_summary"`
	SnapshotMaxID     int64       `json:"snapshot_max_id"`
	FilterHash        string      `json:"filter_hash"`
	ConfirmationToken string      `json:"confirmation_token,omitempty"`
	ExpiresAt         time.Time   `json:"expires_at,omitempty"`
placeholder

type DeleteResult struct {
	DeletedEvents int64   `json:"deleted_events"`
	DeletedJobs   int64   `json:"deleted_jobs"`
	JobIDs        []int64 `json:"-"`
placeholder

type EventRepository interface {
	ListEvents(ctx context.Context, filter EventFilter, page, pageSize int) (*EventPage, error)
	GetEvent(ctx context.Context, id int64) (*Event, error)
	DeleteEvent(ctx context.Context, id int64) (*DeleteResult, error)
	DeleteEventsByIDs(ctx context.Context, ids []int64) (*DeleteResult, error)
	PreviewDelete(ctx context.Context, filter EventFilter) (*DeletePreview, error)
	DeleteEventsByFilter(ctx context.Context, filter EventFilter, snapshotMaxID int64, batchSize int) (*DeleteResult, error)
placeholder

func (r *PostgreSQLRepository) ListEvents(ctx context.Context, filter EventFilter, page, pageSize int) (*EventPage, error) {
	if page < 1 {
		page = 1
placeholder
	if pageSize < 1 {
		pageSize = 20
placeholder
	if pageSize > 100 {
		pageSize = 100
placeholder
	where, args := buildEventWhere(filter, 1)
	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM prompt_audit_events e`+where, args...).Scan(&total); err != nil {
		return nil, err
placeholder
	queryArgs := append([]any(nil), args...)
	limitIndex := len(queryArgs) + 1
	queryArgs = append(queryArgs, pageSize, (page-1)*pageSize)
	rows, err := r.db.QueryContext(ctx, `SELECT `+eventColumns("e")+` FROM prompt_audit_events e`+where+
		fmt.Sprintf(` ORDER BY e.created_at DESC, e.id DESC LIMIT $%d OFFSET $%d`, limitIndex, limitIndex+1), queryArgs...)
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = rows.Close() placeholder()
	items := make([]*Event, 0, pageSize)
	for rows.Next() {
		event, err := scanEvent(rows)
		if err != nil {
			return nil, err
	placeholder
		items = append(items, event)
placeholder
	if err := rows.Err(); err != nil {
		return nil, err
placeholder
	pages := 0
	if total > 0 {
		pages = int((total + int64(pageSize) - 1) / int64(pageSize))
placeholder
	return &EventPage{Items: items, Total: total, Page: page, PageSize: pageSize, Pages: pagesplaceholder, nil
placeholder

func (r *PostgreSQLRepository) GetEvent(ctx context.Context, id int64) (*Event, error) {
	event, err := scanEvent(r.db.QueryRowContext(ctx, `SELECT `+eventColumns("e")+` FROM prompt_audit_events e WHERE e.id=$1`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrEventNotFound
placeholder
	return event, err
placeholder

func (r *PostgreSQLRepository) DeleteEvent(ctx context.Context, id int64) (*DeleteResult, error) {
	return r.DeleteEventsByIDs(ctx, []int64{idplaceholder)
placeholder

func (r *PostgreSQLRepository) DeleteEventsByIDs(ctx context.Context, ids []int64) (*DeleteResult, error) {
	ids = canonicalInt64s(ids)
	if len(ids) == 0 {
		return &DeleteResult{placeholder, nil
placeholder
	if len(ids) > 500 {
		return nil, errors.New("prompt audit delete batch exceeds 500 events")
placeholder
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = tx.Rollback() placeholder()
	rows, err := tx.QueryContext(ctx, `DELETE FROM prompt_audit_events WHERE id=ANY($1) RETURNING job_id`, pq.Array(ids))
	if err != nil {
		return nil, err
placeholder
	jobIDs, err := scanReturnedJobIDs(rows)
	if err != nil {
		return nil, err
placeholder
	deletedJobs, err := deleteOrphanJobs(ctx, tx, jobIDs)
	if err != nil {
		return nil, err
placeholder
	if err := tx.Commit(); err != nil {
		return nil, err
placeholder
	return &DeleteResult{DeletedEvents: int64(len(jobIDs)), DeletedJobs: deletedJobs, JobIDs: canonicalInt64s(jobIDs)placeholder, nil
placeholder

func (r *PostgreSQLRepository) PreviewDelete(ctx context.Context, filter EventFilter) (*DeletePreview, error) {
	if err := validateDeleteFilter(filter); err != nil {
		return nil, err
placeholder
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: trueplaceholder)
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = tx.Rollback() placeholder()
	where, args := buildEventWhere(filter, 1)
	var count, maxID int64
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*), COALESCE(MAX(e.id),0) FROM prompt_audit_events e`+where, args...).Scan(&count, &maxID); err != nil {
		return nil, err
placeholder
	if err := tx.Commit(); err != nil {
		return nil, err
placeholder
	canonical := canonicalEventFilter(filter)
	return &DeletePreview{MatchedCount: count, FilterSummary: canonical, SnapshotMaxID: maxID, FilterHash: FilterHash(canonical, maxID)placeholder, nil
placeholder

func (r *PostgreSQLRepository) DeleteEventsByFilter(ctx context.Context, filter EventFilter, snapshotMaxID int64, batchSize int) (*DeleteResult, error) {
	if err := validateDeleteFilter(filter); err != nil {
		return nil, err
placeholder
	if snapshotMaxID <= 0 {
		return &DeleteResult{placeholder, nil
placeholder
	if batchSize < 1 || batchSize > 1000 {
		batchSize = 200
placeholder
	total := &DeleteResult{placeholder
	jobSet := map[int64]struct{placeholder{placeholder
	for {
		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
	placeholder
		where, args := buildEventWhere(filter, 1)
		maxIndex := len(args) + 1
		limitIndex := maxIndex + 1
		args = append(args, snapshotMaxID, batchSize)
		rows, err := tx.QueryContext(ctx, `
			WITH selected AS (
				SELECT e.id FROM prompt_audit_events e`+where+
			fmt.Sprintf(` AND e.id <= $%d ORDER BY e.id LIMIT $%d FOR UPDATE SKIP LOCKED`, maxIndex, limitIndex)+`
			), deleted AS (
				DELETE FROM prompt_audit_events e USING selected s WHERE e.id=s.id RETURNING e.job_id
			) SELECT job_id FROM deleted`, args...)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
	placeholder
		jobIDs, err := scanReturnedJobIDs(rows)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
	placeholder
		deletedJobs, err := deleteOrphanJobs(ctx, tx, jobIDs)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
	placeholder
		if err := tx.Commit(); err != nil {
			return nil, err
	placeholder
		total.DeletedEvents += int64(len(jobIDs))
		total.DeletedJobs += deletedJobs
		for _, id := range jobIDs {
			jobSet[id] = struct{placeholder{placeholder
	placeholder
		if len(jobIDs) < batchSize {
			break
	placeholder
placeholder
	for id := range jobSet {
		total.JobIDs = append(total.JobIDs, id)
placeholder
	total.JobIDs = canonicalInt64s(total.JobIDs)
	return total, nil
placeholder

func FilterHash(filter EventFilter, snapshotMaxID int64) string {
	payload := struct {
		Filter        EventFilter `json:"filter"`
		SnapshotMaxID int64       `json:"snapshot_max_id"`
placeholder{canonicalEventFilter(filter), snapshotMaxIDplaceholder
	raw, _ := json.Marshal(payload)
	digest := sha256.Sum256(raw)
	return hex.EncodeToString(digest[:])
placeholder

func validateDeleteFilter(filter EventFilter) error {
	if filter.StartAt == nil || filter.EndAt == nil || !filter.StartAt.Before(*filter.EndAt) {
		return errors.New("prompt audit filter delete requires a valid explicit time range")
placeholder
	return nil
placeholder

func canonicalEventFilter(filter EventFilter) EventFilter {
	filter.Decision = strings.TrimSpace(strings.ToLower(filter.Decision))
	filter.RiskLevel = strings.TrimSpace(strings.ToLower(filter.RiskLevel))
	filter.Endpoint = strings.TrimSpace(filter.Endpoint)
	filter.RequestID = strings.TrimSpace(filter.RequestID)
	filter.PromptHash = strings.ToLower(strings.TrimSpace(filter.PromptHash))
	filter.Keyword = strings.TrimSpace(filter.Keyword)
	if filter.StartAt != nil {
		value := filter.StartAt.UTC()
		filter.StartAt = &value
placeholder
	if filter.EndAt != nil {
		value := filter.EndAt.UTC()
		filter.EndAt = &value
placeholder
	return filter
placeholder

func buildEventWhere(filter EventFilter, firstIndex int) (string, []any) {
	filter = canonicalEventFilter(filter)
	clauses := []string{" WHERE TRUE"placeholder
	args := make([]any, 0, 12)
	add := func(clause string, value any) {
		clauses = append(clauses, fmt.Sprintf(clause, firstIndex+len(args)))
		args = append(args, value)
placeholder
	if filter.Decision != "" {
		add(" AND e.decision=$%d", filter.Decision)
placeholder
	if filter.RiskLevel != "" {
		add(" AND e.risk_level=$%d", filter.RiskLevel)
placeholder
	if filter.Endpoint != "" {
		add(" AND e.endpoint=$%d", filter.Endpoint)
placeholder
	if filter.GroupID != nil {
		add(" AND e.group_id=$%d", *filter.GroupID)
placeholder
	if filter.UserID != nil {
		add(" AND e.user_id=$%d", *filter.UserID)
placeholder
	if filter.APIKeyID != nil {
		add(" AND e.api_key_id=$%d", *filter.APIKeyID)
placeholder
	if filter.RequestID != "" {
		add(" AND e.request_id=$%d", filter.RequestID)
placeholder
	if filter.PromptHash != "" {
		add(" AND e.prompt_hash=$%d", filter.PromptHash)
placeholder
	if filter.Keyword != "" {
		add(` AND (e.request_id ILIKE $%d OR e.prompt_hash ILIKE $%d OR e.redacted_preview ILIKE $%d
			OR e.username_snapshot ILIKE $%d OR e.user_email_snapshot ILIKE $%d OR e.api_key_name_snapshot ILIKE $%d)`, "%"+TrimRunes(filter.Keyword, 128)+"%")
		// The clause has six placeholders but add only supplied one. Rebuild it with one shared placeholder.
		clauses[len(clauses)-1] = fmt.Sprintf(` AND (e.request_id ILIKE $%[1]d OR e.prompt_hash ILIKE $%[1]d OR e.redacted_preview ILIKE $%[1]d
			OR e.username_snapshot ILIKE $%[1]d OR e.user_email_snapshot ILIKE $%[1]d OR e.api_key_name_snapshot ILIKE $%[1]d)`, firstIndex+len(args)-1)
placeholder
	if filter.StartAt != nil {
		add(" AND e.created_at >= $%d", filter.StartAt.UTC())
placeholder
	if filter.EndAt != nil {
		add(" AND e.created_at <= $%d", filter.EndAt.UTC())
placeholder
	return strings.Join(clauses, ""), args
placeholder

func eventColumns(alias string) string {
	return fmt.Sprintf(`%[1]s.id,%[1]s.job_id,%[1]s.request_id,%[1]s.user_id,%[1]s.username_snapshot,
		%[1]s.user_email_snapshot,%[1]s.api_key_id,%[1]s.api_key_name_snapshot,%[1]s.group_id,%[1]s.group_name,
		%[1]s.provider,%[1]s.endpoint,%[1]s.protocol,%[1]s.model,%[1]s.prompt_hash,%[1]s.redacted_preview,
		%[1]s.decision,%[1]s.risk_level,%[1]s.action,%[1]s.categories,%[1]s.matched_scanners,
		%[1]s.scanner_scores,%[1]s.scanner_evidence,%[1]s.scanner_backend,%[1]s.scanner_version,
		%[1]s.guard_endpoint_id,%[1]s.policy_id,%[1]s.policy_version,%[1]s.config_version,
		%[1]s.chunk_total,%[1]s.latency_ms,%[1]s.created_at`, alias)
placeholder

func scanEvent(row rowScanner) (*Event, error) {
	event := &Event{placeholder
	var userID, apiKeyID, groupID sql.NullInt64
	var categories, matched, scores, evidence []byte
	err := row.Scan(&event.ID, &event.JobID, &event.Snapshot.RequestID, &userID,
		&event.Snapshot.UsernameSnapshot, &event.Snapshot.UserEmailSnapshot, &apiKeyID,
		&event.Snapshot.APIKeyNameSnapshot, &groupID, &event.Snapshot.GroupName,
		&event.Snapshot.Provider, &event.Snapshot.Endpoint, &event.Snapshot.Protocol, &event.Snapshot.Model,
		&event.Snapshot.PromptHash, &event.Snapshot.RedactedPreview, &event.Decision, &event.RiskLevel,
		&event.Action, &categories, &matched, &scores, &evidence, &event.ScannerBackend,
		&event.ScannerVersion, &event.GuardEndpointID, &event.PolicyID, &event.PolicyVersion,
		&event.ConfigVersion, &event.ChunkTotal, &event.LatencyMS, &event.CreatedAt)
	if err != nil {
		return nil, err
placeholder
	event.Snapshot.UserID = nullableInt64Value(userID)
	event.Snapshot.APIKeyID = nullableInt64Value(apiKeyID)
	event.Snapshot.GroupID = nullableInt64Ptr(groupID)
	_ = json.Unmarshal(categories, &event.Categories)
	_ = json.Unmarshal(matched, &event.MatchedScanners)
	_ = json.Unmarshal(scores, &event.ScannerScores)
	_ = json.Unmarshal(evidence, &event.ScannerEvidence)
	result := NormalizedResult{Decision: event.Decision, RiskLevel: event.RiskLevel, Action: event.Action,
		Categories: event.Categories, MatchedScanners: event.MatchedScanners, ScannerScores: event.ScannerScores,
		ScannerEvidence: event.ScannerEvidenceplaceholder
	event.IssueSummaries = BuildIssueSummaries(result)
	return event, nil
placeholder

func scanReturnedJobIDs(rows *sql.Rows) ([]int64, error) {
	defer func() { _ = rows.Close() placeholder()
	result := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
	placeholder
		result = append(result, id)
placeholder
	return result, rows.Err()
placeholder

func deleteOrphanJobs(ctx context.Context, tx *sql.Tx, jobIDs []int64) (int64, error) {
	jobIDs = canonicalInt64s(jobIDs)
	if len(jobIDs) == 0 {
		return 0, nil
placeholder
	result, err := tx.ExecContext(ctx, `DELETE FROM prompt_audit_jobs j
		WHERE j.id=ANY($1) AND j.status <> 'processing'
		AND NOT EXISTS (SELECT 1 FROM prompt_audit_events e WHERE e.job_id=j.id)`, pq.Array(jobIDs))
	if err != nil {
		return 0, err
placeholder
	return result.RowsAffected()
placeholder
