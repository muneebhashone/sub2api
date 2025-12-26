package repository

import "gorm.io/gorm"

// AutoMigrate runs schema migrations for all repository persistence models.
// Persistence models are defined within individual `*_repo.go` files.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
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
placeholder
