package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"
)

type channelRepository struct {
	db *sql.DB
placeholder

// NewChannelRepository 创建渠道数据访问实例
func NewChannelRepository(db *sql.DB) service.ChannelRepository {
	return &channelRepository{db: dbplaceholder
placeholder

// runInTx 在事务中执行 fn，成功 commit，失败 rollback。
func (r *channelRepository) runInTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
placeholder
	defer func() { _ = tx.Rollback() placeholder()

	if err := fn(tx); err != nil {
		return err
placeholder
	return tx.Commit()
placeholder

func (r *channelRepository) Create(ctx context.Context, channel *service.Channel) error {
	return r.runInTx(ctx, func(tx *sql.Tx) error {
		modelMappingJSON, err := marshalModelMapping(channel.ModelMapping)
		if err != nil {
			return err
	placeholder
		err = tx.QueryRowContext(ctx,
			`INSERT INTO channels (name, description, status, model_mapping, billing_model_source, restrict_models, features) VALUES ($1, $2, $3, $4, $5, $6, $7)
			 RETURNING id, created_at, updated_at`,
			channel.Name, channel.Description, channel.Status, modelMappingJSON, channel.BillingModelSource, channel.RestrictModels, channel.Features,
		).Scan(&channel.ID, &channel.CreatedAt, &channel.UpdatedAt)
		if err != nil {
			if isUniqueViolation(err) {
				return service.ErrChannelExists
		placeholder
			return fmt.Errorf("insert channel: %w", err)
	placeholder

		// 设置分组关联
		if len(channel.GroupIDs) > 0 {
			if err := setGroupIDsTx(ctx, tx, channel.ID, channel.GroupIDs); err != nil {
				return err
		placeholder
	placeholder

		// 设置模型定价
		if len(channel.ModelPricing) > 0 {
			if err := replaceModelPricingTx(ctx, tx, channel.ID, channel.ModelPricing); err != nil {
				return err
		placeholder
	placeholder

		return nil
placeholder)
placeholder

func (r *channelRepository) GetByID(ctx context.Context, id int64) (*service.Channel, error) {
	ch := &service.Channel{placeholder
	var modelMappingJSON []byte
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, description, status, model_mapping, billing_model_source, restrict_models, features, created_at, updated_at
		 FROM channels WHERE id = $1`, id,
	).Scan(&ch.ID, &ch.Name, &ch.Description, &ch.Status, &modelMappingJSON, &ch.BillingModelSource, &ch.RestrictModels, &ch.Features, &ch.CreatedAt, &ch.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, service.ErrChannelNotFound
placeholder
	if err != nil {
		return nil, fmt.Errorf("get channel: %w", err)
placeholder
	ch.ModelMapping = unmarshalModelMapping(modelMappingJSON)

	groupIDs, err := r.GetGroupIDs(ctx, id)
	if err != nil {
		return nil, err
placeholder
	ch.GroupIDs = groupIDs

	pricing, err := r.ListModelPricing(ctx, id)
	if err != nil {
		return nil, err
placeholder
	ch.ModelPricing = pricing

	return ch, nil
placeholder

func (r *channelRepository) Update(ctx context.Context, channel *service.Channel) error {
	return r.runInTx(ctx, func(tx *sql.Tx) error {
		modelMappingJSON, err := marshalModelMapping(channel.ModelMapping)
		if err != nil {
			return err
	placeholder
		result, err := tx.ExecContext(ctx,
			`UPDATE channels SET name = $1, description = $2, status = $3, model_mapping = $4, billing_model_source = $5, restrict_models = $6, features = $7, updated_at = NOW()
			 WHERE id = $8`,
			channel.Name, channel.Description, channel.Status, modelMappingJSON, channel.BillingModelSource, channel.RestrictModels, channel.Features, channel.ID,
		)
		if err != nil {
			if isUniqueViolation(err) {
				return service.ErrChannelExists
		placeholder
			return fmt.Errorf("update channel: %w", err)
	placeholder
		rows, _ := result.RowsAffected()
		if rows == 0 {
			return service.ErrChannelNotFound
	placeholder

		// 更新分组关联
		if channel.GroupIDs != nil {
			if err := setGroupIDsTx(ctx, tx, channel.ID, channel.GroupIDs); err != nil {
				return err
		placeholder
	placeholder

		// 更新模型定价
		if channel.ModelPricing != nil {
			if err := replaceModelPricingTx(ctx, tx, channel.ID, channel.ModelPricing); err != nil {
				return err
		placeholder
	placeholder

		return nil
placeholder)
placeholder

func (r *channelRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM channels WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete channel: %w", err)
placeholder
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return service.ErrChannelNotFound
placeholder
	return nil
placeholder

func (r *channelRepository) List(ctx context.Context, params pagination.PaginationParams, status, search string) ([]service.Channel, *pagination.PaginationResult, error) {
	where := []string{"1=1"placeholder
	args := []any{placeholder
	argIdx := 1

	if status != "" {
		where = append(where, fmt.Sprintf("c.status = $%d", argIdx))
		args = append(args, status)
		argIdx++
placeholder
	if search != "" {
		where = append(where, fmt.Sprintf("(c.name ILIKE $%d OR c.description ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+escapeLike(search)+"%")
		argIdx++
placeholder

	whereClause := strings.Join(where, " AND ")

	// 计数
	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM channels c WHERE %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, nil, fmt.Errorf("count channels: %w", err)
placeholder

	pageSize := params.Limit() // 约束在 [1, 100]
	page := params.Page
	if page < 1 {
		page = 1
placeholder
	offset := (page - 1) * pageSize

	// 查询 channel 列表
	dataQuery := fmt.Sprintf(
		`SELECT c.id, c.name, c.description, c.status, c.model_mapping, c.billing_model_source, c.restrict_models, c.created_at, c.updated_at
		 FROM channels c WHERE %s ORDER BY %s LIMIT $%d OFFSET $%d`,
		whereClause, channelListOrderBy(params), argIdx, argIdx+1,
	)
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("query channels: %w", err)
placeholder
	defer func() { _ = rows.Close() placeholder()

	var channels []service.Channel
	var channelIDs []int64
	for rows.Next() {
		var ch service.Channel
		var modelMappingJSON []byte
		if err := rows.Scan(&ch.ID, &ch.Name, &ch.Description, &ch.Status, &modelMappingJSON, &ch.BillingModelSource, &ch.RestrictModels, &ch.Features, &ch.CreatedAt, &ch.UpdatedAt); err != nil {
			return nil, nil, fmt.Errorf("scan channel: %w", err)
	placeholder
		ch.ModelMapping = unmarshalModelMapping(modelMappingJSON)
		channels = append(channels, ch)
		channelIDs = append(channelIDs, ch.ID)
placeholder
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("iterate channels: %w", err)
placeholder

	// 批量加载分组 ID 和模型定价（避免 N+1）
	if len(channelIDs) > 0 {
		groupMap, err := r.batchLoadGroupIDs(ctx, channelIDs)
		if err != nil {
			return nil, nil, err
	placeholder
		pricingMap, err := r.batchLoadModelPricing(ctx, channelIDs)
		if err != nil {
			return nil, nil, err
	placeholder
		for i := range channels {
			channels[i].GroupIDs = groupMap[channels[i].ID]
			channels[i].ModelPricing = pricingMap[channels[i].ID]
	placeholder
placeholder

	pages := 0
	if total > 0 {
		pages = int((total + int64(pageSize) - 1) / int64(pageSize))
placeholder

	paginationResult := &pagination.PaginationResult{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
placeholder

	return channels, paginationResult, nil
placeholder

func channelListOrderBy(params pagination.PaginationParams) string {
	sortBy := strings.ToLower(strings.TrimSpace(params.SortBy))
	sortOrder := strings.ToUpper(params.NormalizedSortOrder(pagination.SortOrderAsc))

	var column string
	switch sortBy {
	case "":
		column = "c.id"
		sortOrder = "ASC"
	case "id":
		column = "c.id"
	case "name":
		column = "c.name"
	case "status":
		column = "c.status"
	case "created_at":
		column = "c.created_at"
	default:
		column = "c.id"
		sortOrder = "ASC"
placeholder

	return fmt.Sprintf("%s %s, c.id %s", column, sortOrder, sortOrder)
placeholder

func (r *channelRepository) ListAll(ctx context.Context) ([]service.Channel, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, description, status, model_mapping, billing_model_source, restrict_models, features, created_at, updated_at FROM channels ORDER BY id`,
	)
	if err != nil {
		return nil, fmt.Errorf("query all channels: %w", err)
placeholder
	defer func() { _ = rows.Close() placeholder()

	var channels []service.Channel
	var channelIDs []int64
	for rows.Next() {
		var ch service.Channel
		var modelMappingJSON []byte
		if err := rows.Scan(&ch.ID, &ch.Name, &ch.Description, &ch.Status, &modelMappingJSON, &ch.BillingModelSource, &ch.RestrictModels, &ch.Features, &ch.CreatedAt, &ch.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan channel: %w", err)
	placeholder
		ch.ModelMapping = unmarshalModelMapping(modelMappingJSON)
		channels = append(channels, ch)
		channelIDs = append(channelIDs, ch.ID)
placeholder
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate channels: %w", err)
placeholder

	if len(channelIDs) == 0 {
		return channels, nil
placeholder

	// 批量加载分组 ID
	groupMap, err := r.batchLoadGroupIDs(ctx, channelIDs)
	if err != nil {
		return nil, err
placeholder

	// 批量加载模型定价
	pricingMap, err := r.batchLoadModelPricing(ctx, channelIDs)
	if err != nil {
		return nil, err
placeholder

	for i := range channels {
		channels[i].GroupIDs = groupMap[channels[i].ID]
		channels[i].ModelPricing = pricingMap[channels[i].ID]
placeholder

	return channels, nil
placeholder

// --- 批量加载辅助方法 ---

// batchLoadGroupIDs 批量加载多个渠道的分组 ID
func (r *channelRepository) batchLoadGroupIDs(ctx context.Context, channelIDs []int64) (map[int64][]int64, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT channel_id, group_id FROM channel_groups
		 WHERE channel_id = ANY($1) ORDER BY channel_id, group_id`,
		pq.Array(channelIDs),
	)
	if err != nil {
		return nil, fmt.Errorf("batch load group ids: %w", err)
placeholder
	defer func() { _ = rows.Close() placeholder()

	groupMap := make(map[int64][]int64, len(channelIDs))
	for rows.Next() {
		var channelID, groupID int64
		if err := rows.Scan(&channelID, &groupID); err != nil {
			return nil, fmt.Errorf("scan group id: %w", err)
	placeholder
		groupMap[channelID] = append(groupMap[channelID], groupID)
placeholder
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate group ids: %w", err)
placeholder
	return groupMap, nil
placeholder

func (r *channelRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM channels WHERE name = $1)`, name,
	).Scan(&exists)
	return exists, err
placeholder

func (r *channelRepository) ExistsByNameExcluding(ctx context.Context, name string, excludeID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM channels WHERE name = $1 AND id != $2)`, name, excludeID,
	).Scan(&exists)
	return exists, err
placeholder

// --- 分组关联 ---

func (r *channelRepository) GetGroupIDs(ctx context.Context, channelID int64) ([]int64, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT group_id FROM channel_groups WHERE channel_id = $1 ORDER BY group_id`, channelID,
	)
	if err != nil {
		return nil, fmt.Errorf("get group ids: %w", err)
placeholder
	defer func() { _ = rows.Close() placeholder()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan group id: %w", err)
	placeholder
		ids = append(ids, id)
placeholder
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate group ids: %w", err)
placeholder
	return ids, nil
placeholder

func (r *channelRepository) SetGroupIDs(ctx context.Context, channelID int64, groupIDs []int64) error {
	return setGroupIDsTx(ctx, r.db, channelID, groupIDs)
placeholder

func (r *channelRepository) GetChannelIDByGroupID(ctx context.Context, groupID int64) (int64, error) {
	var channelID int64
	err := r.db.QueryRowContext(ctx,
		`SELECT channel_id FROM channel_groups WHERE group_id = $1`, groupID,
	).Scan(&channelID)
	if err == sql.ErrNoRows {
		return 0, nil
placeholder
	return channelID, err
placeholder

func (r *channelRepository) GetGroupsInOtherChannels(ctx context.Context, channelID int64, groupIDs []int64) ([]int64, error) {
	if len(groupIDs) == 0 {
		return nil, nil
placeholder
	rows, err := r.db.QueryContext(ctx,
		`SELECT group_id FROM channel_groups WHERE group_id = ANY($1) AND channel_id != $2`,
		pq.Array(groupIDs), channelID,
	)
	if err != nil {
		return nil, fmt.Errorf("get groups in other channels: %w", err)
placeholder
	defer func() { _ = rows.Close() placeholder()

	var conflicting []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan conflicting group id: %w", err)
	placeholder
		conflicting = append(conflicting, id)
placeholder
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate conflicting group ids: %w", err)
placeholder
	return conflicting, nil
placeholder

// marshalModelMapping 将 model mapping 序列化为嵌套 JSON 字节
// 格式：{"platform": {"src": "dst"placeholder, ...placeholder
func marshalModelMapping(m map[string]map[string]string) ([]byte, error) {
	if len(m) == 0 {
		return []byte("{placeholder"), nil
placeholder
	data, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("marshal model_mapping: %w", err)
placeholder
	return data, nil
placeholder

// unmarshalModelMapping 将 JSON 字节反序列化为嵌套 model mapping
func unmarshalModelMapping(data []byte) map[string]map[string]string {
	if len(data) == 0 {
		return nil
placeholder
	var m map[string]map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return nil
placeholder
	return m
placeholder

// GetGroupPlatforms 批量查询分组 ID 对应的平台
func (r *channelRepository) GetGroupPlatforms(ctx context.Context, groupIDs []int64) (map[int64]string, error) {
	if len(groupIDs) == 0 {
		return make(map[int64]string), nil
placeholder
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, platform FROM groups WHERE id = ANY($1)`,
		pq.Array(groupIDs),
	)
	if err != nil {
		return nil, fmt.Errorf("get group platforms: %w", err)
placeholder
	defer rows.Close() //nolint:errcheck

	result := make(map[int64]string, len(groupIDs))
	for rows.Next() {
		var id int64
		var platform string
		if err := rows.Scan(&id, &platform); err != nil {
			return nil, fmt.Errorf("scan group platform: %w", err)
	placeholder
		result[id] = platform
placeholder
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate group platforms: %w", err)
placeholder
	return result, nil
placeholder
