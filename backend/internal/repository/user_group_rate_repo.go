package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"
)

type userGroupRateRepository struct {
	sql sqlExecutor
placeholder

// NewUserGroupRateRepository 创建用户专属分组倍率仓储
func NewUserGroupRateRepository(sqlDB *sql.DB) service.UserGroupRateRepository {
	return &userGroupRateRepository{sql: sqlDBplaceholder
placeholder

// GetByUserID 获取用户的所有专属分组倍率
func (r *userGroupRateRepository) GetByUserID(ctx context.Context, userID int64) (map[int64]float64, error) {
	query := `SELECT group_id, rate_multiplier FROM user_group_rate_multipliers WHERE user_id = $1`
	rows, err := r.sql.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = rows.Close() placeholder()

	result := make(map[int64]float64)
	for rows.Next() {
		var groupID int64
		var rate float64
		if err := rows.Scan(&groupID, &rate); err != nil {
			return nil, err
	placeholder
		result[groupID] = rate
placeholder
	if err := rows.Err(); err != nil {
		return nil, err
placeholder
	return result, nil
placeholder

// GetByUserIDs 批量获取多个用户的专属分组倍率。
// 返回结构：map[userID]map[groupID]rate
func (r *userGroupRateRepository) GetByUserIDs(ctx context.Context, userIDs []int64) (map[int64]map[int64]float64, error) {
	result := make(map[int64]map[int64]float64, len(userIDs))
	if len(userIDs) == 0 {
		return result, nil
placeholder

	uniqueIDs := make([]int64, 0, len(userIDs))
	seen := make(map[int64]struct{placeholder, len(userIDs))
	for _, userID := range userIDs {
		if userID <= 0 {
			continue
	placeholder
		if _, exists := seen[userID]; exists {
			continue
	placeholder
		seen[userID] = struct{placeholder{placeholder
		uniqueIDs = append(uniqueIDs, userID)
		result[userID] = make(map[int64]float64)
placeholder
	if len(uniqueIDs) == 0 {
		return result, nil
placeholder

	rows, err := r.sql.QueryContext(ctx, `
		SELECT user_id, group_id, rate_multiplier
		FROM user_group_rate_multipliers
		WHERE user_id = ANY($1)
	`, pq.Array(uniqueIDs))
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = rows.Close() placeholder()

	for rows.Next() {
		var userID int64
		var groupID int64
		var rate float64
		if err := rows.Scan(&userID, &groupID, &rate); err != nil {
			return nil, err
	placeholder
		if _, ok := result[userID]; !ok {
			result[userID] = make(map[int64]float64)
	placeholder
		result[userID][groupID] = rate
placeholder
	if err := rows.Err(); err != nil {
		return nil, err
placeholder
	return result, nil
placeholder

// GetByGroupID 获取指定分组下所有用户的专属倍率
func (r *userGroupRateRepository) GetByGroupID(ctx context.Context, groupID int64) ([]service.UserGroupRateEntry, error) {
	query := `
		SELECT ugr.user_id, u.email, ugr.rate_multiplier
		FROM user_group_rate_multipliers ugr
		JOIN users u ON u.id = ugr.user_id
		WHERE ugr.group_id = $1
		ORDER BY ugr.user_id
	`
	rows, err := r.sql.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = rows.Close() placeholder()

	var result []service.UserGroupRateEntry
	for rows.Next() {
		var entry service.UserGroupRateEntry
		if err := rows.Scan(&entry.UserID, &entry.UserEmail, &entry.RateMultiplier); err != nil {
			return nil, err
	placeholder
		result = append(result, entry)
placeholder
	if err := rows.Err(); err != nil {
		return nil, err
placeholder
	return result, nil
placeholder

// GetByUserAndGroup 获取用户在特定分组的专属倍率
func (r *userGroupRateRepository) GetByUserAndGroup(ctx context.Context, userID, groupID int64) (*float64, error) {
	query := `SELECT rate_multiplier FROM user_group_rate_multipliers WHERE user_id = $1 AND group_id = $2`
	var rate float64
	err := scanSingleRow(ctx, r.sql, query, []any{userID, groupIDplaceholder, &rate)
	if err == sql.ErrNoRows {
		return nil, nil
placeholder
	if err != nil {
		return nil, err
placeholder
	return &rate, nil
placeholder

// SyncUserGroupRates 同步用户的分组专属倍率
func (r *userGroupRateRepository) SyncUserGroupRates(ctx context.Context, userID int64, rates map[int64]*float64) error {
	if len(rates) == 0 {
		// 如果传入空 map，删除该用户的所有专属倍率
		_, err := r.sql.ExecContext(ctx, `DELETE FROM user_group_rate_multipliers WHERE user_id = $1`, userID)
		return err
placeholder

	// 分离需要删除和需要 upsert 的记录
	var toDelete []int64
	upsertGroupIDs := make([]int64, 0, len(rates))
	upsertRates := make([]float64, 0, len(rates))
	for groupID, rate := range rates {
		if rate == nil {
			toDelete = append(toDelete, groupID)
	placeholder else {
			upsertGroupIDs = append(upsertGroupIDs, groupID)
			upsertRates = append(upsertRates, *rate)
	placeholder
placeholder

	// 删除指定的记录
	if len(toDelete) > 0 {
		if _, err := r.sql.ExecContext(ctx,
			`DELETE FROM user_group_rate_multipliers WHERE user_id = $1 AND group_id = ANY($2)`,
			userID, pq.Array(toDelete)); err != nil {
			return err
	placeholder
placeholder

	// Upsert 记录
	now := time.Now()
	if len(upsertGroupIDs) > 0 {
		_, err := r.sql.ExecContext(ctx, `
			INSERT INTO user_group_rate_multipliers (user_id, group_id, rate_multiplier, created_at, updated_at)
			SELECT
				$1::bigint,
				data.group_id,
				data.rate_multiplier,
				$2::timestamptz,
				$2::timestamptz
			FROM unnest($3::bigint[], $4::double precision[]) AS data(group_id, rate_multiplier)
			ON CONFLICT (user_id, group_id)
			DO UPDATE SET
				rate_multiplier = EXCLUDED.rate_multiplier,
				updated_at = EXCLUDED.updated_at
		`, userID, now, pq.Array(upsertGroupIDs), pq.Array(upsertRates))
		if err != nil {
			return err
	placeholder
placeholder

	return nil
placeholder

// DeleteByGroupID 删除指定分组的所有用户专属倍率
func (r *userGroupRateRepository) DeleteByGroupID(ctx context.Context, groupID int64) error {
	_, err := r.sql.ExecContext(ctx, `DELETE FROM user_group_rate_multipliers WHERE group_id = $1`, groupID)
	return err
placeholder

// DeleteByUserID 删除指定用户的所有专属倍率
func (r *userGroupRateRepository) DeleteByUserID(ctx context.Context, userID int64) error {
	_, err := r.sql.ExecContext(ctx, `DELETE FROM user_group_rate_multipliers WHERE user_id = $1`, userID)
	return err
placeholder
