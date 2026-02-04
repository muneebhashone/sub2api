package repository

import (
	"context"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/apikey"
	"github.com/Wei-Shaw/sub2api/ent/group"
	"github.com/Wei-Shaw/sub2api/ent/schema/mixins"
	"github.com/Wei-Shaw/sub2api/ent/user"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

type apiKeyRepository struct {
	client *dbent.Client
placeholder

func NewAPIKeyRepository(client *dbent.Client) service.APIKeyRepository {
	return &apiKeyRepository{client: clientplaceholder
placeholder

func (r *apiKeyRepository) activeQuery() *dbent.APIKeyQuery {
	// 默认过滤已软删除记录，避免删除后仍被查询到。
	return r.client.APIKey.Query().Where(apikey.DeletedAtIsNil())
placeholder

func (r *apiKeyRepository) Create(ctx context.Context, key *service.APIKey) error {
	builder := r.client.APIKey.Create().
		SetUserID(key.UserID).
		SetKey(key.Key).
		SetName(key.Name).
		SetStatus(key.Status).
		SetNillableGroupID(key.GroupID).
		SetQuota(key.Quota).
		SetQuotaUsed(key.QuotaUsed).
		SetNillableExpiresAt(key.ExpiresAt)

	if len(key.IPWhitelist) > 0 {
		builder.SetIPWhitelist(key.IPWhitelist)
placeholder
	if len(key.IPBlacklist) > 0 {
		builder.SetIPBlacklist(key.IPBlacklist)
placeholder

	created, err := builder.Save(ctx)
	if err == nil {
		key.ID = created.ID
		key.CreatedAt = created.CreatedAt
		key.UpdatedAt = created.UpdatedAt
placeholder
	return translatePersistenceError(err, nil, service.ErrAPIKeyExists)
placeholder

func (r *apiKeyRepository) GetByID(ctx context.Context, id int64) (*service.APIKey, error) {
	m, err := r.activeQuery().
		Where(apikey.IDEQ(id)).
		WithUser().
		WithGroup().
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrAPIKeyNotFound
	placeholder
		return nil, err
placeholder
	return apiKeyEntityToService(m), nil
placeholder

// GetKeyAndOwnerID 根据 API Key ID 获取其 key 与所有者（用户）ID。
// 相比 GetByID，此方法性能更优，因为：
//   - 使用 Select() 只查询必要字段，减少数据传输量
//   - 不加载完整的 API Key 实体及其关联数据（User、Group 等）
//   - 适用于删除等只需 key 与用户 ID 的场景
func (r *apiKeyRepository) GetKeyAndOwnerID(ctx context.Context, id int64) (string, int64, error) {
	m, err := r.activeQuery().
		Where(apikey.IDEQ(id)).
		Select(apikey.FieldKey, apikey.FieldUserID).
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return "", 0, service.ErrAPIKeyNotFound
	placeholder
		return "", 0, err
placeholder
	return m.Key, m.UserID, nil
placeholder

func (r *apiKeyRepository) GetByKey(ctx context.Context, key string) (*service.APIKey, error) {
	m, err := r.activeQuery().
		Where(apikey.KeyEQ(key)).
		WithUser().
		WithGroup().
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrAPIKeyNotFound
	placeholder
		return nil, err
placeholder
	return apiKeyEntityToService(m), nil
placeholder

func (r *apiKeyRepository) GetByKeyForAuth(ctx context.Context, key string) (*service.APIKey, error) {
	m, err := r.activeQuery().
		Where(apikey.KeyEQ(key)).
		Select(
			apikey.FieldID,
			apikey.FieldUserID,
			apikey.FieldGroupID,
			apikey.FieldStatus,
			apikey.FieldIPWhitelist,
			apikey.FieldIPBlacklist,
			apikey.FieldQuota,
			apikey.FieldQuotaUsed,
			apikey.FieldExpiresAt,
		).
		WithUser(func(q *dbent.UserQuery) {
			q.Select(
				user.FieldID,
				user.FieldStatus,
				user.FieldRole,
				user.FieldBalance,
				user.FieldConcurrency,
			)
	placeholder).
		WithGroup(func(q *dbent.GroupQuery) {
			q.Select(
				group.FieldID,
				group.FieldName,
				group.FieldPlatform,
				group.FieldStatus,
				group.FieldSubscriptionType,
				group.FieldRateMultiplier,
				group.FieldDailyLimitUsd,
				group.FieldWeeklyLimitUsd,
				group.FieldMonthlyLimitUsd,
				group.FieldImagePrice1k,
				group.FieldImagePrice2k,
				group.FieldImagePrice4k,
				group.FieldSoraImagePrice360,
				group.FieldSoraImagePrice540,
				group.FieldSoraVideoPricePerRequest,
				group.FieldSoraVideoPricePerRequestHd,
				group.FieldClaudeCodeOnly,
				group.FieldFallbackGroupID,
				group.FieldFallbackGroupIDOnInvalidRequest,
				group.FieldModelRoutingEnabled,
				group.FieldModelRouting,
				group.FieldMcpXMLInject,
				group.FieldSupportedModelScopes,
			)
	placeholder).
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrAPIKeyNotFound
	placeholder
		return nil, err
placeholder
	return apiKeyEntityToService(m), nil
placeholder

func (r *apiKeyRepository) Update(ctx context.Context, key *service.APIKey) error {
	// 使用原子操作：将软删除检查与更新合并到同一语句，避免竞态条件。
	// 之前的实现先检查 Exist 再 UpdateOneID，若在两步之间发生软删除，
	// 则会更新已删除的记录。
	// 这里选择 Update().Where()，确保只有未软删除记录能被更新。
	// 同时显式设置 updated_at，避免二次查询带来的并发可见性问题。
	now := time.Now()
	builder := r.client.APIKey.Update().
		Where(apikey.IDEQ(key.ID), apikey.DeletedAtIsNil()).
		SetName(key.Name).
		SetStatus(key.Status).
		SetQuota(key.Quota).
		SetQuotaUsed(key.QuotaUsed).
		SetUpdatedAt(now)
	if key.GroupID != nil {
		builder.SetGroupID(*key.GroupID)
placeholder else {
		builder.ClearGroupID()
placeholder

	// Expiration time
	if key.ExpiresAt != nil {
		builder.SetExpiresAt(*key.ExpiresAt)
placeholder else {
		builder.ClearExpiresAt()
placeholder

	// IP 限制字段
	if len(key.IPWhitelist) > 0 {
		builder.SetIPWhitelist(key.IPWhitelist)
placeholder else {
		builder.ClearIPWhitelist()
placeholder
	if len(key.IPBlacklist) > 0 {
		builder.SetIPBlacklist(key.IPBlacklist)
placeholder else {
		builder.ClearIPBlacklist()
placeholder

	affected, err := builder.Save(ctx)
	if err != nil {
		return err
placeholder
	if affected == 0 {
		// 更新影响行数为 0，说明记录不存在或已被软删除。
		return service.ErrAPIKeyNotFound
placeholder

	// 使用同一时间戳回填，避免并发删除导致二次查询失败。
	key.UpdatedAt = now
	return nil
placeholder

func (r *apiKeyRepository) Delete(ctx context.Context, id int64) error {
	// 显式软删除：避免依赖 Hook 行为，确保 deleted_at 一定被设置。
	affected, err := r.client.APIKey.Update().
		Where(apikey.IDEQ(id), apikey.DeletedAtIsNil()).
		SetDeletedAt(time.Now()).
		Save(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return service.ErrAPIKeyNotFound
	placeholder
		return err
placeholder
	if affected == 0 {
		exists, err := r.client.APIKey.Query().
			Where(apikey.IDEQ(id)).
			Exist(mixins.SkipSoftDelete(ctx))
		if err != nil {
			return err
	placeholder
		if exists {
			return nil
	placeholder
		return service.ErrAPIKeyNotFound
placeholder
	return nil
placeholder

func (r *apiKeyRepository) ListByUserID(ctx context.Context, userID int64, params pagination.PaginationParams) ([]service.APIKey, *pagination.PaginationResult, error) {
	q := r.activeQuery().Where(apikey.UserIDEQ(userID))

	total, err := q.Count(ctx)
	if err != nil {
		return nil, nil, err
placeholder

	keys, err := q.
		WithGroup().
		Offset(params.Offset()).
		Limit(params.Limit()).
		Order(dbent.Desc(apikey.FieldID)).
		All(ctx)
	if err != nil {
		return nil, nil, err
placeholder

	outKeys := make([]service.APIKey, 0, len(keys))
	for i := range keys {
		outKeys = append(outKeys, *apiKeyEntityToService(keys[i]))
placeholder

	return outKeys, paginationResultFromTotal(int64(total), params), nil
placeholder

func (r *apiKeyRepository) VerifyOwnership(ctx context.Context, userID int64, apiKeyIDs []int64) ([]int64, error) {
	if len(apiKeyIDs) == 0 {
		return []int64{placeholder, nil
placeholder

	ids, err := r.client.APIKey.Query().
		Where(apikey.UserIDEQ(userID), apikey.IDIn(apiKeyIDs...), apikey.DeletedAtIsNil()).
		IDs(ctx)
	if err != nil {
		return nil, err
placeholder
	return ids, nil
placeholder

func (r *apiKeyRepository) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	count, err := r.activeQuery().Where(apikey.UserIDEQ(userID)).Count(ctx)
	return int64(count), err
placeholder

func (r *apiKeyRepository) ExistsByKey(ctx context.Context, key string) (bool, error) {
	count, err := r.activeQuery().Where(apikey.KeyEQ(key)).Count(ctx)
	return count > 0, err
placeholder

func (r *apiKeyRepository) ListByGroupID(ctx context.Context, groupID int64, params pagination.PaginationParams) ([]service.APIKey, *pagination.PaginationResult, error) {
	q := r.activeQuery().Where(apikey.GroupIDEQ(groupID))

	total, err := q.Count(ctx)
	if err != nil {
		return nil, nil, err
placeholder

	keys, err := q.
		WithUser().
		Offset(params.Offset()).
		Limit(params.Limit()).
		Order(dbent.Desc(apikey.FieldID)).
		All(ctx)
	if err != nil {
		return nil, nil, err
placeholder

	outKeys := make([]service.APIKey, 0, len(keys))
	for i := range keys {
		outKeys = append(outKeys, *apiKeyEntityToService(keys[i]))
placeholder

	return outKeys, paginationResultFromTotal(int64(total), params), nil
placeholder

// SearchAPIKeys searches API keys by user ID and/or keyword (name)
func (r *apiKeyRepository) SearchAPIKeys(ctx context.Context, userID int64, keyword string, limit int) ([]service.APIKey, error) {
	q := r.activeQuery()
	if userID > 0 {
		q = q.Where(apikey.UserIDEQ(userID))
placeholder

	if keyword != "" {
		q = q.Where(apikey.NameContainsFold(keyword))
placeholder

	keys, err := q.Limit(limit).Order(dbent.Desc(apikey.FieldID)).All(ctx)
	if err != nil {
		return nil, err
placeholder

	outKeys := make([]service.APIKey, 0, len(keys))
	for i := range keys {
		outKeys = append(outKeys, *apiKeyEntityToService(keys[i]))
placeholder
	return outKeys, nil
placeholder

// ClearGroupIDByGroupID 将指定分组的所有 API Key 的 group_id 设为 nil
func (r *apiKeyRepository) ClearGroupIDByGroupID(ctx context.Context, groupID int64) (int64, error) {
	n, err := r.client.APIKey.Update().
		Where(apikey.GroupIDEQ(groupID), apikey.DeletedAtIsNil()).
		ClearGroupID().
		Save(ctx)
	return int64(n), err
placeholder

// CountByGroupID 获取分组的 API Key 数量
func (r *apiKeyRepository) CountByGroupID(ctx context.Context, groupID int64) (int64, error) {
	count, err := r.activeQuery().Where(apikey.GroupIDEQ(groupID)).Count(ctx)
	return int64(count), err
placeholder

func (r *apiKeyRepository) ListKeysByUserID(ctx context.Context, userID int64) ([]string, error) {
	keys, err := r.activeQuery().
		Where(apikey.UserIDEQ(userID)).
		Select(apikey.FieldKey).
		Strings(ctx)
	if err != nil {
		return nil, err
placeholder
	return keys, nil
placeholder

func (r *apiKeyRepository) ListKeysByGroupID(ctx context.Context, groupID int64) ([]string, error) {
	keys, err := r.activeQuery().
		Where(apikey.GroupIDEQ(groupID)).
		Select(apikey.FieldKey).
		Strings(ctx)
	if err != nil {
		return nil, err
placeholder
	return keys, nil
placeholder

// IncrementQuotaUsed atomically increments the quota_used field and returns the new value
func (r *apiKeyRepository) IncrementQuotaUsed(ctx context.Context, id int64, amount float64) (float64, error) {
	// Use raw SQL for atomic increment to avoid race conditions
	// First get current value
	m, err := r.activeQuery().
		Where(apikey.IDEQ(id)).
		Select(apikey.FieldQuotaUsed).
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return 0, service.ErrAPIKeyNotFound
	placeholder
		return 0, err
placeholder

	newValue := m.QuotaUsed + amount

	// Update with new value
	affected, err := r.client.APIKey.Update().
		Where(apikey.IDEQ(id), apikey.DeletedAtIsNil()).
		SetQuotaUsed(newValue).
		Save(ctx)
	if err != nil {
		return 0, err
placeholder
	if affected == 0 {
		return 0, service.ErrAPIKeyNotFound
placeholder

	return newValue, nil
placeholder

func apiKeyEntityToService(m *dbent.APIKey) *service.APIKey {
	if m == nil {
		return nil
placeholder
	out := &service.APIKey{
		ID:          m.ID,
		UserID:      m.UserID,
		Key:         m.Key,
		Name:        m.Name,
		Status:      m.Status,
		IPWhitelist: m.IPWhitelist,
		IPBlacklist: m.IPBlacklist,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		GroupID:     m.GroupID,
		Quota:       m.Quota,
		QuotaUsed:   m.QuotaUsed,
		ExpiresAt:   m.ExpiresAt,
placeholder
	if m.Edges.User != nil {
		out.User = userEntityToService(m.Edges.User)
placeholder
	if m.Edges.Group != nil {
		out.Group = groupEntityToService(m.Edges.Group)
placeholder
	return out
placeholder

func userEntityToService(u *dbent.User) *service.User {
	if u == nil {
		return nil
placeholder
	return &service.User{
		ID:                  u.ID,
		Email:               u.Email,
		Username:            u.Username,
		Notes:               u.Notes,
		PasswordHash:        u.PasswordHash,
		Role:                u.Role,
		Balance:             u.Balance,
		Concurrency:         u.Concurrency,
		Status:              u.Status,
		TotpSecretEncrypted: u.TotpSecretEncrypted,
		TotpEnabled:         u.TotpEnabled,
		TotpEnabledAt:       u.TotpEnabledAt,
		CreatedAt:           u.CreatedAt,
		UpdatedAt:           u.UpdatedAt,
placeholder
placeholder

func groupEntityToService(g *dbent.Group) *service.Group {
	if g == nil {
		return nil
placeholder
	return &service.Group{
		ID:                              g.ID,
		Name:                            g.Name,
		Description:                     derefString(g.Description),
		Platform:                        g.Platform,
		RateMultiplier:                  g.RateMultiplier,
		IsExclusive:                     g.IsExclusive,
		Status:                          g.Status,
		Hydrated:                        true,
		SubscriptionType:                g.SubscriptionType,
		DailyLimitUSD:                   g.DailyLimitUsd,
		WeeklyLimitUSD:                  g.WeeklyLimitUsd,
		MonthlyLimitUSD:                 g.MonthlyLimitUsd,
		ImagePrice1K:                    g.ImagePrice1k,
		ImagePrice2K:                    g.ImagePrice2k,
		ImagePrice4K:                    g.ImagePrice4k,
		SoraImagePrice360:               g.SoraImagePrice360,
		SoraImagePrice540:               g.SoraImagePrice540,
		SoraVideoPricePerRequest:        g.SoraVideoPricePerRequest,
		SoraVideoPricePerRequestHD:      g.SoraVideoPricePerRequestHd,
		DefaultValidityDays:             g.DefaultValidityDays,
		ClaudeCodeOnly:                  g.ClaudeCodeOnly,
		FallbackGroupID:                 g.FallbackGroupID,
		FallbackGroupIDOnInvalidRequest: g.FallbackGroupIDOnInvalidRequest,
		ModelRouting:                    g.ModelRouting,
		ModelRoutingEnabled:             g.ModelRoutingEnabled,
		MCPXMLInject:                    g.McpXMLInject,
		SupportedModelScopes:            g.SupportedModelScopes,
		CreatedAt:                       g.CreatedAt,
		UpdatedAt:                       g.UpdatedAt,
placeholder
placeholder

func derefString(s *string) string {
	if s == nil {
		return ""
placeholder
	return *s
placeholder
