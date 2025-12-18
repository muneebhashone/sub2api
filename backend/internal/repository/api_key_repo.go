package repository

import (
	"context"
	"sub2api/internal/model"

	"gorm.io/gorm"
)

type ApiKeyRepository struct {
	db *gorm.DB
placeholder

func NewApiKeyRepository(db *gorm.DB) *ApiKeyRepository {
	return &ApiKeyRepository{db: dbplaceholder
placeholder

func (r *ApiKeyRepository) Create(ctx context.Context, key *model.ApiKey) error {
	return r.db.WithContext(ctx).Create(key).Error
placeholder

func (r *ApiKeyRepository) GetByID(ctx context.Context, id int64) (*model.ApiKey, error) {
	var key model.ApiKey
	err := r.db.WithContext(ctx).Preload("User").Preload("Group").First(&key, id).Error
	if err != nil {
		return nil, err
placeholder
	return &key, nil
placeholder

func (r *ApiKeyRepository) GetByKey(ctx context.Context, key string) (*model.ApiKey, error) {
	var apiKey model.ApiKey
	err := r.db.WithContext(ctx).Preload("User").Preload("Group").Where("key = ?", key).First(&apiKey).Error
	if err != nil {
		return nil, err
placeholder
	return &apiKey, nil
placeholder

func (r *ApiKeyRepository) Update(ctx context.Context, key *model.ApiKey) error {
	return r.db.WithContext(ctx).Model(key).Select("name", "group_id", "status", "updated_at").Updates(key).Error
placeholder

func (r *ApiKeyRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.ApiKey{placeholder, id).Error
placeholder

func (r *ApiKeyRepository) ListByUserID(ctx context.Context, userID int64, params PaginationParams) ([]model.ApiKey, *PaginationResult, error) {
	var keys []model.ApiKey
	var total int64

	db := r.db.WithContext(ctx).Model(&model.ApiKey{placeholder).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, nil, err
placeholder

	if err := db.Preload("Group").Offset(params.Offset()).Limit(params.Limit()).Order("id DESC").Find(&keys).Error; err != nil {
		return nil, nil, err
placeholder

	pages := int(total) / params.Limit()
	if int(total)%params.Limit() > 0 {
		pages++
placeholder

	return keys, &PaginationResult{
		Total:    total,
		Page:     params.Page,
		PageSize: params.Limit(),
		Pages:    pages,
placeholder, nil
placeholder

func (r *ApiKeyRepository) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ApiKey{placeholder).Where("user_id = ?", userID).Count(&count).Error
	return count, err
placeholder

func (r *ApiKeyRepository) ExistsByKey(ctx context.Context, key string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ApiKey{placeholder).Where("key = ?", key).Count(&count).Error
	return count > 0, err
placeholder

func (r *ApiKeyRepository) ListByGroupID(ctx context.Context, groupID int64, params PaginationParams) ([]model.ApiKey, *PaginationResult, error) {
	var keys []model.ApiKey
	var total int64

	db := r.db.WithContext(ctx).Model(&model.ApiKey{placeholder).Where("group_id = ?", groupID)

	if err := db.Count(&total).Error; err != nil {
		return nil, nil, err
placeholder

	if err := db.Preload("User").Offset(params.Offset()).Limit(params.Limit()).Order("id DESC").Find(&keys).Error; err != nil {
		return nil, nil, err
placeholder

	pages := int(total) / params.Limit()
	if int(total)%params.Limit() > 0 {
		pages++
placeholder

	return keys, &PaginationResult{
		Total:    total,
		Page:     params.Page,
		PageSize: params.Limit(),
		Pages:    pages,
placeholder, nil
placeholder

// SearchApiKeys searches API keys by user ID and/or keyword (name)
func (r *ApiKeyRepository) SearchApiKeys(ctx context.Context, userID int64, keyword string, limit int) ([]model.ApiKey, error) {
	var keys []model.ApiKey

	db := r.db.WithContext(ctx).Model(&model.ApiKey{placeholder)

	if userID > 0 {
		db = db.Where("user_id = ?", userID)
placeholder

	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		db = db.Where("name ILIKE ?", searchPattern)
placeholder

	if err := db.Limit(limit).Order("id DESC").Find(&keys).Error; err != nil {
		return nil, err
placeholder

	return keys, nil
placeholder

// ClearGroupIDByGroupID 将指定分组的所有 API Key 的 group_id 设为 nil
func (r *ApiKeyRepository) ClearGroupIDByGroupID(ctx context.Context, groupID int64) (int64, error) {
	result := r.db.WithContext(ctx).Model(&model.ApiKey{placeholder).
		Where("group_id = ?", groupID).
		Update("group_id", nil)
	return result.RowsAffected, result.Error
placeholder

// CountByGroupID 获取分组的 API Key 数量
func (r *ApiKeyRepository) CountByGroupID(ctx context.Context, groupID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ApiKey{placeholder).Where("group_id = ?", groupID).Count(&count).Error
	return count, err
placeholder
