package repository

import (
	"context"
	"time"

	"sub2api/internal/model"

	"gorm.io/gorm"
)

// UserSubscriptionRepository 用户订阅仓库
type UserSubscriptionRepository struct {
	db *gorm.DB
placeholder

// NewUserSubscriptionRepository 创建用户订阅仓库
func NewUserSubscriptionRepository(db *gorm.DB) *UserSubscriptionRepository {
	return &UserSubscriptionRepository{db: dbplaceholder
placeholder

// Create 创建订阅
func (r *UserSubscriptionRepository) Create(ctx context.Context, sub *model.UserSubscription) error {
	return r.db.WithContext(ctx).Create(sub).Error
placeholder

// GetByID 根据ID获取订阅
func (r *UserSubscriptionRepository) GetByID(ctx context.Context, id int64) (*model.UserSubscription, error) {
	var sub model.UserSubscription
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Group").
		Preload("AssignedByUser").
		First(&sub, id).Error
	if err != nil {
		return nil, err
placeholder
	return &sub, nil
placeholder

// GetByUserIDAndGroupID 根据用户ID和分组ID获取订阅
func (r *UserSubscriptionRepository) GetByUserIDAndGroupID(ctx context.Context, userID, groupID int64) (*model.UserSubscription, error) {
	var sub model.UserSubscription
	err := r.db.WithContext(ctx).
		Preload("Group").
		Where("user_id = ? AND group_id = ?", userID, groupID).
		First(&sub).Error
	if err != nil {
		return nil, err
placeholder
	return &sub, nil
placeholder

// GetActiveByUserIDAndGroupID 获取用户对特定分组的有效订阅
func (r *UserSubscriptionRepository) GetActiveByUserIDAndGroupID(ctx context.Context, userID, groupID int64) (*model.UserSubscription, error) {
	var sub model.UserSubscription
	err := r.db.WithContext(ctx).
		Preload("Group").
		Where("user_id = ? AND group_id = ? AND status = ? AND expires_at > ?",
			userID, groupID, model.SubscriptionStatusActive, time.Now()).
		First(&sub).Error
	if err != nil {
		return nil, err
placeholder
	return &sub, nil
placeholder

// Update 更新订阅
func (r *UserSubscriptionRepository) Update(ctx context.Context, sub *model.UserSubscription) error {
	sub.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(sub).Error
placeholder

// Delete 删除订阅
func (r *UserSubscriptionRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.UserSubscription{placeholder, id).Error
placeholder

// ListByUserID 获取用户的所有订阅
func (r *UserSubscriptionRepository) ListByUserID(ctx context.Context, userID int64) ([]model.UserSubscription, error) {
	var subs []model.UserSubscription
	err := r.db.WithContext(ctx).
		Preload("Group").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&subs).Error
	return subs, err
placeholder

// ListActiveByUserID 获取用户的所有有效订阅
func (r *UserSubscriptionRepository) ListActiveByUserID(ctx context.Context, userID int64) ([]model.UserSubscription, error) {
	var subs []model.UserSubscription
	err := r.db.WithContext(ctx).
		Preload("Group").
		Where("user_id = ? AND status = ? AND expires_at > ?",
			userID, model.SubscriptionStatusActive, time.Now()).
		Order("created_at DESC").
		Find(&subs).Error
	return subs, err
placeholder

// ListByGroupID 获取分组的所有订阅（分页）
func (r *UserSubscriptionRepository) ListByGroupID(ctx context.Context, groupID int64, params PaginationParams) ([]model.UserSubscription, *PaginationResult, error) {
	var subs []model.UserSubscription
	var total int64

	query := r.db.WithContext(ctx).Model(&model.UserSubscription{placeholder).Where("group_id = ?", groupID)

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
placeholder

	err := query.
		Preload("User").
		Preload("Group").
		Order("created_at DESC").
		Offset(params.Offset()).
		Limit(params.Limit()).
		Find(&subs).Error
	if err != nil {
		return nil, nil, err
placeholder

	pages := int(total) / params.Limit()
	if int(total)%params.Limit() > 0 {
		pages++
placeholder

	return subs, &PaginationResult{
		Total:    total,
		Page:     params.Page,
		PageSize: params.Limit(),
		Pages:    pages,
placeholder, nil
placeholder

// List 获取所有订阅（分页，支持筛选）
func (r *UserSubscriptionRepository) List(ctx context.Context, params PaginationParams, userID, groupID *int64, status string) ([]model.UserSubscription, *PaginationResult, error) {
	var subs []model.UserSubscription
	var total int64

	query := r.db.WithContext(ctx).Model(&model.UserSubscription{placeholder)

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
placeholder
	if groupID != nil {
		query = query.Where("group_id = ?", *groupID)
placeholder
	if status != "" {
		query = query.Where("status = ?", status)
placeholder

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
placeholder

	err := query.
		Preload("User").
		Preload("Group").
		Preload("AssignedByUser").
		Order("created_at DESC").
		Offset(params.Offset()).
		Limit(params.Limit()).
		Find(&subs).Error
	if err != nil {
		return nil, nil, err
placeholder

	pages := int(total) / params.Limit()
	if int(total)%params.Limit() > 0 {
		pages++
placeholder

	return subs, &PaginationResult{
		Total:    total,
		Page:     params.Page,
		PageSize: params.Limit(),
		Pages:    pages,
placeholder, nil
placeholder

// IncrementUsage 增加使用量
func (r *UserSubscriptionRepository) IncrementUsage(ctx context.Context, id int64, costUSD float64) error {
	return r.db.WithContext(ctx).Model(&model.UserSubscription{placeholder).
		Where("id = ?", id).
		Updates(map[string]interface{placeholder{
			"daily_usage_usd":   gorm.Expr("daily_usage_usd + ?", costUSD),
			"weekly_usage_usd":  gorm.Expr("weekly_usage_usd + ?", costUSD),
			"monthly_usage_usd": gorm.Expr("monthly_usage_usd + ?", costUSD),
			"updated_at":        time.Now(),
	placeholder).Error
placeholder

// ResetDailyUsage 重置日使用量
func (r *UserSubscriptionRepository) ResetDailyUsage(ctx context.Context, id int64, newWindowStart time.Time) error {
	return r.db.WithContext(ctx).Model(&model.UserSubscription{placeholder).
		Where("id = ?", id).
		Updates(map[string]interface{placeholder{
			"daily_usage_usd":    0,
			"daily_window_start": newWindowStart,
			"updated_at":         time.Now(),
	placeholder).Error
placeholder

// ResetWeeklyUsage 重置周使用量
func (r *UserSubscriptionRepository) ResetWeeklyUsage(ctx context.Context, id int64, newWindowStart time.Time) error {
	return r.db.WithContext(ctx).Model(&model.UserSubscription{placeholder).
		Where("id = ?", id).
		Updates(map[string]interface{placeholder{
			"weekly_usage_usd":    0,
			"weekly_window_start": newWindowStart,
			"updated_at":          time.Now(),
	placeholder).Error
placeholder

// ResetMonthlyUsage 重置月使用量
func (r *UserSubscriptionRepository) ResetMonthlyUsage(ctx context.Context, id int64, newWindowStart time.Time) error {
	return r.db.WithContext(ctx).Model(&model.UserSubscription{placeholder).
		Where("id = ?", id).
		Updates(map[string]interface{placeholder{
			"monthly_usage_usd":    0,
			"monthly_window_start": newWindowStart,
			"updated_at":           time.Now(),
	placeholder).Error
placeholder

// ActivateWindows 激活所有窗口（首次使用时）
func (r *UserSubscriptionRepository) ActivateWindows(ctx context.Context, id int64, activateTime time.Time) error {
	return r.db.WithContext(ctx).Model(&model.UserSubscription{placeholder).
		Where("id = ?", id).
		Updates(map[string]interface{placeholder{
			"daily_window_start":   activateTime,
			"weekly_window_start":  activateTime,
			"monthly_window_start": activateTime,
			"updated_at":           time.Now(),
	placeholder).Error
placeholder

// UpdateStatus 更新订阅状态
func (r *UserSubscriptionRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&model.UserSubscription{placeholder).
		Where("id = ?", id).
		Updates(map[string]interface{placeholder{
			"status":     status,
			"updated_at": time.Now(),
	placeholder).Error
placeholder

// ExtendExpiry 延长订阅过期时间
func (r *UserSubscriptionRepository) ExtendExpiry(ctx context.Context, id int64, newExpiresAt time.Time) error {
	return r.db.WithContext(ctx).Model(&model.UserSubscription{placeholder).
		Where("id = ?", id).
		Updates(map[string]interface{placeholder{
			"expires_at": newExpiresAt,
			"updated_at": time.Now(),
	placeholder).Error
placeholder

// UpdateNotes 更新订阅备注
func (r *UserSubscriptionRepository) UpdateNotes(ctx context.Context, id int64, notes string) error {
	return r.db.WithContext(ctx).Model(&model.UserSubscription{placeholder).
		Where("id = ?", id).
		Updates(map[string]interface{placeholder{
			"notes":      notes,
			"updated_at": time.Now(),
	placeholder).Error
placeholder

// ListExpired 获取所有已过期但状态仍为active的订阅
func (r *UserSubscriptionRepository) ListExpired(ctx context.Context) ([]model.UserSubscription, error) {
	var subs []model.UserSubscription
	err := r.db.WithContext(ctx).
		Where("status = ? AND expires_at <= ?", model.SubscriptionStatusActive, time.Now()).
		Find(&subs).Error
	return subs, err
placeholder

// BatchUpdateExpiredStatus 批量更新过期订阅状态
func (r *UserSubscriptionRepository) BatchUpdateExpiredStatus(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).Model(&model.UserSubscription{placeholder).
		Where("status = ? AND expires_at <= ?", model.SubscriptionStatusActive, time.Now()).
		Updates(map[string]interface{placeholder{
			"status":     model.SubscriptionStatusExpired,
			"updated_at": time.Now(),
	placeholder)
	return result.RowsAffected, result.Error
placeholder

// ExistsByUserIDAndGroupID 检查用户是否已有该分组的订阅
func (r *UserSubscriptionRepository) ExistsByUserIDAndGroupID(ctx context.Context, userID, groupID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.UserSubscription{placeholder).
		Where("user_id = ? AND group_id = ?", userID, groupID).
		Count(&count).Error
	return count > 0, err
placeholder

// CountByGroupID 获取分组的订阅数量
func (r *UserSubscriptionRepository) CountByGroupID(ctx context.Context, groupID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.UserSubscription{placeholder).
		Where("group_id = ?", groupID).
		Count(&count).Error
	return count, err
placeholder

// CountActiveByGroupID 获取分组的有效订阅数量
func (r *UserSubscriptionRepository) CountActiveByGroupID(ctx context.Context, groupID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.UserSubscription{placeholder).
		Where("group_id = ? AND status = ? AND expires_at > ?",
			groupID, model.SubscriptionStatusActive, time.Now()).
		Count(&count).Error
	return count, err
placeholder

// DeleteByGroupID 删除分组相关的所有订阅记录
func (r *UserSubscriptionRepository) DeleteByGroupID(ctx context.Context, groupID int64) (int64, error) {
	result := r.db.WithContext(ctx).Where("group_id = ?", groupID).Delete(&model.UserSubscription{placeholder)
	return result.RowsAffected, result.Error
placeholder
