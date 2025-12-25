package repository

import (
	"context"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/Wei-Shaw/sub2api/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SettingRepository 系统设置数据访问层
type settingRepository struct {
	db *gorm.DB
placeholder

// NewSettingRepository 创建系统设置仓库实例
func NewSettingRepository(db *gorm.DB) service.SettingRepository {
	return &settingRepository{db: dbplaceholder
placeholder

// Get 根据Key获取设置值
func (r *settingRepository) Get(ctx context.Context, key string) (*model.Setting, error) {
	var setting model.Setting
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&setting).Error
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrSettingNotFound, nil)
placeholder
	return &setting, nil
placeholder

// GetValue 获取设置值字符串
func (r *settingRepository) GetValue(ctx context.Context, key string) (string, error) {
	setting, err := r.Get(ctx, key)
	if err != nil {
		return "", err
placeholder
	return setting.Value, nil
placeholder

// Set 设置值（存在则更新，不存在则创建）
func (r *settingRepository) Set(ctx context.Context, key, value string) error {
	setting := &model.Setting{
		Key:       key,
		Value:     value,
		UpdatedAt: time.Now(),
placeholder

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"placeholderplaceholder,
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"placeholder),
placeholder).Create(setting).Error
placeholder

// GetMultiple 批量获取设置
func (r *settingRepository) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	var settings []model.Setting
	err := r.db.WithContext(ctx).Where("key IN ?", keys).Find(&settings).Error
	if err != nil {
		return nil, err
placeholder

	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
placeholder
	return result, nil
placeholder

// SetMultiple 批量设置值
func (r *settingRepository) SetMultiple(ctx context.Context, settings map[string]string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for key, value := range settings {
			setting := &model.Setting{
				Key:       key,
				Value:     value,
				UpdatedAt: time.Now(),
		placeholder
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "key"placeholderplaceholder,
				DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"placeholder),
		placeholder).Create(setting).Error; err != nil {
				return err
		placeholder
	placeholder
		return nil
placeholder)
placeholder

// GetAll 获取所有设置
func (r *settingRepository) GetAll(ctx context.Context) (map[string]string, error) {
	var settings []model.Setting
	err := r.db.WithContext(ctx).Find(&settings).Error
	if err != nil {
		return nil, err
placeholder

	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
placeholder
	return result, nil
placeholder

// Delete 删除设置
func (r *settingRepository) Delete(ctx context.Context, key string) error {
	return r.db.WithContext(ctx).Where("key = ?", key).Delete(&model.Setting{placeholder).Error
placeholder
