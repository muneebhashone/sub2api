package repository

import (
	"gorm.io/gorm"
)

// Repositories 所有仓库的集合
type Repositories struct {
	User             *UserRepository
	ApiKey           *ApiKeyRepository
	Group            *GroupRepository
	Account          *AccountRepository
	Proxy            *ProxyRepository
	RedeemCode       *RedeemCodeRepository
	UsageLog         *UsageLogRepository
	Setting          *SettingRepository
	UserSubscription *UserSubscriptionRepository
placeholder

// NewRepositories 创建所有仓库
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:             NewUserRepository(db),
		ApiKey:           NewApiKeyRepository(db),
		Group:            NewGroupRepository(db),
		Account:          NewAccountRepository(db),
		Proxy:            NewProxyRepository(db),
		RedeemCode:       NewRedeemCodeRepository(db),
		UsageLog:         NewUsageLogRepository(db),
		Setting:          NewSettingRepository(db),
		UserSubscription: NewUserSubscriptionRepository(db),
placeholder
placeholder

// PaginationParams 分页参数
type PaginationParams struct {
	Page     int
	PageSize int
placeholder

// PaginationResult 分页结果
type PaginationResult struct {
	Total    int64
	Page     int
	PageSize int
	Pages    int
placeholder

// DefaultPagination 默认分页参数
func DefaultPagination() PaginationParams {
	return PaginationParams{
		Page:     1,
		PageSize: 20,
placeholder
placeholder

// Offset 计算偏移量
func (p PaginationParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
placeholder
	return (p.Page - 1) * p.PageSize
placeholder

// Limit 获取限制数
func (p PaginationParams) Limit() int {
	if p.PageSize < 1 {
		return 20
placeholder
	if p.PageSize > 100 {
		return 100
placeholder
	return p.PageSize
placeholder
