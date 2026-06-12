package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"
)

type schedulerOutboxRepository struct {
	db *sql.DB
placeholder

func NewSchedulerOutboxRepository(db *sql.DB) service.SchedulerOutboxRepository {
	return &schedulerOutboxRepository{db: dbplaceholder
placeholder

func (r *schedulerOutboxRepository) ListAfter(ctx context.Context, afterID int64, limit int) ([]service.SchedulerOutboxEvent, error) {
	if limit <= 0 {
		limit = 100
placeholder
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, event_type, account_id, group_id, payload, created_at
		FROM scheduler_outbox
		WHERE id > $1
		ORDER BY id ASC
		LIMIT $2
	`, afterID, limit)
	if err != nil {
		return nil, err
placeholder
	defer func() {
		_ = rows.Close()
placeholder()

	events := make([]service.SchedulerOutboxEvent, 0, limit)
	for rows.Next() {
		var (
			payloadRaw []byte
			accountID  sql.NullInt64
			groupID    sql.NullInt64
			event      service.SchedulerOutboxEvent
		)
		if err := rows.Scan(&event.ID, &event.EventType, &accountID, &groupID, &payloadRaw, &event.CreatedAt); err != nil {
			return nil, err
	placeholder
		if accountID.Valid {
			v := accountID.Int64
			event.AccountID = &v
	placeholder
		if groupID.Valid {
			v := groupID.Int64
			event.GroupID = &v
	placeholder
		if len(payloadRaw) > 0 {
			var payload map[string]any
			if err := json.Unmarshal(payloadRaw, &payload); err != nil {
				return nil, err
		placeholder
			event.Payload = payload
	placeholder
		events = append(events, event)
placeholder
	if err := rows.Err(); err != nil {
		return nil, err
placeholder
	return events, nil
placeholder

func (r *schedulerOutboxRepository) MaxID(ctx context.Context) (int64, error) {
	var maxID int64
	if err := r.db.QueryRowContext(ctx, "SELECT COALESCE(MAX(id), 0) FROM scheduler_outbox").Scan(&maxID); err != nil {
		return 0, err
placeholder
	return maxID, nil
placeholder

func (r *schedulerOutboxRepository) MarkProcessed(ctx context.Context, eventIDs []int64) error {
	if len(eventIDs) == 0 {
		return nil
placeholder
	_, err := r.db.ExecContext(ctx, `
		UPDATE scheduler_outbox
		SET dedup_key = NULL
		WHERE id = ANY($1)
			AND dedup_key IS NOT NULL
	`, pq.Array(eventIDs))
	return err
placeholder

func enqueueSchedulerOutbox(ctx context.Context, exec sqlExecutor, eventType string, accountID *int64, groupID *int64, payload any) error {
	if exec == nil {
		return nil
placeholder
	var payloadArg any
	var payloadJSON []byte
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			return err
	placeholder
		payloadArg = encoded
		payloadJSON = encoded
placeholder
	query := `
		INSERT INTO scheduler_outbox (event_type, account_id, group_id, payload)
		VALUES ($1, $2, $3, $4)
	`
	args := []any{eventType, accountID, groupID, payloadArgplaceholder
	if schedulerOutboxEventSupportsDedup(eventType) {
		dedupKey := schedulerOutboxDedupKey(eventType, accountID, groupID, payloadJSON)
		query = `
			INSERT INTO scheduler_outbox (event_type, account_id, group_id, payload, dedup_key)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (dedup_key) WHERE dedup_key IS NOT NULL DO NOTHING
		`
		args = append(args, dedupKey)
placeholder
	_, err := exec.ExecContext(ctx, query, args...)
	return err
placeholder

func schedulerOutboxDedupKey(eventType string, accountID *int64, groupID *int64, payloadJSON []byte) string {
	h := sha256.New()
	_, _ = h.Write([]byte(eventType))
	_, _ = h.Write([]byte{0placeholder)
	if accountID != nil {
		_, _ = h.Write([]byte(strconv.FormatInt(*accountID, 10)))
placeholder
	_, _ = h.Write([]byte{0placeholder)
	if groupID != nil {
		_, _ = h.Write([]byte(strconv.FormatInt(*groupID, 10)))
placeholder
	_, _ = h.Write([]byte{0placeholder)
	_, _ = h.Write(payloadJSON)
	return fmt.Sprintf("scheduler_outbox:%s", hex.EncodeToString(h.Sum(nil)))
placeholder

func schedulerOutboxEventSupportsDedup(eventType string) bool {
	switch eventType {
	case service.SchedulerOutboxEventAccountChanged,
		service.SchedulerOutboxEventGroupChanged,
		service.SchedulerOutboxEventFullRebuild:
		return true
	default:
		return false
placeholder
placeholder
