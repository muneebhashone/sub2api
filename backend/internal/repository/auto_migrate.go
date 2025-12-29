package repository

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// MaxExpiresAt is the maximum allowed expiration date for subscriptions (year 2099)
// This prevents time.Time JSON serialization errors (RFC 3339 requires year <= 9999)
var maxExpiresAt = time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC)

// AutoMigrate runs schema migrations for all repository persistence models.
// Persistence models are defined within individual `*_repo.go` files.
func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&userModel{placeholder,
		&apiKeyModel{placeholder,
		&groupModel{placeholder,
		&accountModel{placeholder,
		&accountGroupModel{placeholder,
		&proxyModel{placeholder,
		&redeemCodeModel{placeholder,
		&usageLogModel{placeholder,
		&settingModel{placeholder,
		&userSubscriptionModel{placeholder,
	)
	if err != nil {
		return err
placeholder

	// 创建默认分组(简易模式支持)
	if err := ensureDefaultGroups(db); err != nil {
		return err
placeholder

	// 修复无效的过期时间（年份超过 2099 会导致 JSON 序列化失败）
	return fixInvalidExpiresAt(db)
placeholder

// fixInvalidExpiresAt 修复 user_subscriptions 表中无效的过期时间
func fixInvalidExpiresAt(db *gorm.DB) error {
	result := db.Model(&userSubscriptionModel{placeholder).
		Where("expires_at > ?", maxExpiresAt).
		Update("expires_at", maxExpiresAt)
	if result.Error != nil {
		return result.Error
placeholder
	if result.RowsAffected > 0 {
		log.Printf("[AutoMigrate] Fixed %d subscriptions with invalid expires_at (year > 2099)", result.RowsAffected)
placeholder
	return nil
placeholder

// ensureDefaultGroups 确保默认分组存在(简易模式支持)
// 为每个平台创建一个默认分组,配置最大权限以确保简易模式下不受限制
func ensureDefaultGroups(db *gorm.DB) error {
	defaultGroups := []struct {
		name        string
		platform    string
		description string
placeholder{
		{
			name:        "anthropic-default",
			platform:    "anthropic",
			description: "Default group for Anthropic accounts (Simple Mode)",
	placeholder,
		{
			name:        "openai-default",
			platform:    "openai",
			description: "Default group for OpenAI accounts (Simple Mode)",
	placeholder,
		{
			name:        "gemini-default",
			platform:    "gemini",
			description: "Default group for Gemini accounts (Simple Mode)",
	placeholder,
placeholder

	for _, dg := range defaultGroups {
		var count int64
		if err := db.Model(&groupModel{placeholder).Where("name = ?", dg.name).Count(&count).Error; err != nil {
			return err
	placeholder

		if count == 0 {
			group := &groupModel{
				Name:             dg.name,
				Description:      dg.description,
				Platform:         dg.platform,
				RateMultiplier:   1.0,
				IsExclusive:      false,
				Status:           "active",
				SubscriptionType: "standard",
		placeholder
			if err := db.Create(group).Error; err != nil {
				log.Printf("[AutoMigrate] Failed to create default group %s: %v", dg.name, err)
				return err
		placeholder
			log.Printf("[AutoMigrate] Created default group: %s (platform: %s)", dg.name, dg.platform)
	placeholder
placeholder

	return nil
placeholder
