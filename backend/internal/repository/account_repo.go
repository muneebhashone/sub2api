package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type accountRepository struct {
	db *gorm.DB
placeholder

func NewAccountRepository(db *gorm.DB) service.AccountRepository {
	return &accountRepository{db: dbplaceholder
placeholder

func (r *accountRepository) Create(ctx context.Context, account *service.Account) error {
	m := accountModelFromService(account)
	err := r.db.WithContext(ctx).Create(m).Error
	if err == nil {
		applyAccountModelToService(account, m)
placeholder
	return err
placeholder

func (r *accountRepository) GetByID(ctx context.Context, id int64) (*service.Account, error) {
	var m accountModel
	err := r.db.WithContext(ctx).Preload("Proxy").Preload("AccountGroups.Group").First(&m, id).Error
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrAccountNotFound, nil)
placeholder
	return accountModelToService(&m), nil
placeholder

func (r *accountRepository) GetByCRSAccountID(ctx context.Context, crsAccountID string) (*service.Account, error) {
	if crsAccountID == "" {
		return nil, nil
placeholder

	var m accountModel
	err := r.db.WithContext(ctx).Where("extra->>'crs_account_id' = ?", crsAccountID).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
	placeholder
		return nil, err
placeholder
	return accountModelToService(&m), nil
placeholder

func (r *accountRepository) Update(ctx context.Context, account *service.Account) error {
	m := accountModelFromService(account)
	err := r.db.WithContext(ctx).Save(m).Error
	if err == nil {
		applyAccountModelToService(account, m)
placeholder
	return err
placeholder

func (r *accountRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Where("account_id = ?", id).Delete(&accountGroupModel{placeholder).Error; err != nil {
		return err
placeholder
	return r.db.WithContext(ctx).Delete(&accountModel{placeholder, id).Error
placeholder

func (r *accountRepository) List(ctx context.Context, params pagination.PaginationParams) ([]service.Account, *pagination.PaginationResult, error) {
	return r.ListWithFilters(ctx, params, "", "", "", "")
placeholder

func (r *accountRepository) ListWithFilters(ctx context.Context, params pagination.PaginationParams, platform, accountType, status, search string) ([]service.Account, *pagination.PaginationResult, error) {
	var accounts []accountModel
	var total int64

	db := r.db.WithContext(ctx).Model(&accountModel{placeholder)

	if platform != "" {
		db = db.Where("platform = ?", platform)
placeholder
	if accountType != "" {
		db = db.Where("type = ?", accountType)
placeholder
	if status != "" {
		db = db.Where("status = ?", status)
placeholder
	if search != "" {
		searchPattern := "%" + search + "%"
		db = db.Where("name ILIKE ?", searchPattern)
placeholder

	if err := db.Count(&total).Error; err != nil {
		return nil, nil, err
placeholder

	if err := db.Preload("Proxy").Preload("AccountGroups.Group").Offset(params.Offset()).Limit(params.Limit()).Order("id DESC").Find(&accounts).Error; err != nil {
		return nil, nil, err
placeholder

	outAccounts := make([]service.Account, 0, len(accounts))
	for i := range accounts {
		outAccounts = append(outAccounts, *accountModelToService(&accounts[i]))
placeholder

	return outAccounts, paginationResultFromTotal(total, params), nil
placeholder

func (r *accountRepository) ListByGroup(ctx context.Context, groupID int64) ([]service.Account, error) {
	var accounts []accountModel
	err := r.db.WithContext(ctx).
		Joins("JOIN account_groups ON account_groups.account_id = accounts.id").
		Where("account_groups.group_id = ? AND accounts.status = ?", groupID, service.StatusActive).
		Preload("Proxy").
		Order("account_groups.priority ASC, accounts.priority ASC").
		Find(&accounts).Error
	if err != nil {
		return nil, err
placeholder

	outAccounts := make([]service.Account, 0, len(accounts))
	for i := range accounts {
		outAccounts = append(outAccounts, *accountModelToService(&accounts[i]))
placeholder
	return outAccounts, nil
placeholder

func (r *accountRepository) ListActive(ctx context.Context) ([]service.Account, error) {
	var accounts []accountModel
	err := r.db.WithContext(ctx).
		Where("status = ?", service.StatusActive).
		Preload("Proxy").
		Order("priority ASC").
		Find(&accounts).Error
	if err != nil {
		return nil, err
placeholder

	outAccounts := make([]service.Account, 0, len(accounts))
	for i := range accounts {
		outAccounts = append(outAccounts, *accountModelToService(&accounts[i]))
placeholder
	return outAccounts, nil
placeholder

func (r *accountRepository) ListByPlatform(ctx context.Context, platform string) ([]service.Account, error) {
	var accounts []accountModel
	err := r.db.WithContext(ctx).
		Where("platform = ? AND status = ?", platform, service.StatusActive).
		Preload("Proxy").
		Order("priority ASC").
		Find(&accounts).Error
	if err != nil {
		return nil, err
placeholder

	outAccounts := make([]service.Account, 0, len(accounts))
	for i := range accounts {
		outAccounts = append(outAccounts, *accountModelToService(&accounts[i]))
placeholder
	return outAccounts, nil
placeholder

func (r *accountRepository) UpdateLastUsed(ctx context.Context, id int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&accountModel{placeholder).Where("id = ?", id).Update("last_used_at", now).Error
placeholder

func (r *accountRepository) SetError(ctx context.Context, id int64, errorMsg string) error {
	return r.db.WithContext(ctx).Model(&accountModel{placeholder).Where("id = ?", id).
		Updates(map[string]any{
			"status":        service.StatusError,
			"error_message": errorMsg,
	placeholder).Error
placeholder

func (r *accountRepository) AddToGroup(ctx context.Context, accountID, groupID int64, priority int) error {
	ag := &accountGroupModel{
		AccountID: accountID,
		GroupID:   groupID,
		Priority:  priority,
placeholder
	return r.db.WithContext(ctx).Create(ag).Error
placeholder

func (r *accountRepository) RemoveFromGroup(ctx context.Context, accountID, groupID int64) error {
	return r.db.WithContext(ctx).Where("account_id = ? AND group_id = ?", accountID, groupID).
		Delete(&accountGroupModel{placeholder).Error
placeholder

func (r *accountRepository) GetGroups(ctx context.Context, accountID int64) ([]service.Group, error) {
	var groups []groupModel
	err := r.db.WithContext(ctx).
		Joins("JOIN account_groups ON account_groups.group_id = groups.id").
		Where("account_groups.account_id = ?", accountID).
		Find(&groups).Error
	if err != nil {
		return nil, err
placeholder

	outGroups := make([]service.Group, 0, len(groups))
	for i := range groups {
		outGroups = append(outGroups, *groupModelToService(&groups[i]))
placeholder
	return outGroups, nil
placeholder

func (r *accountRepository) BindGroups(ctx context.Context, accountID int64, groupIDs []int64) error {
	if err := r.db.WithContext(ctx).Where("account_id = ?", accountID).Delete(&accountGroupModel{placeholder).Error; err != nil {
		return err
placeholder

	if len(groupIDs) == 0 {
		return nil
placeholder

	accountGroups := make([]accountGroupModel, 0, len(groupIDs))
	for i, groupID := range groupIDs {
		accountGroups = append(accountGroups, accountGroupModel{
			AccountID: accountID,
			GroupID:   groupID,
			Priority:  i + 1,
	placeholder)
placeholder
	return r.db.WithContext(ctx).Create(&accountGroups).Error
placeholder

func (r *accountRepository) ListSchedulable(ctx context.Context) ([]service.Account, error) {
	var accounts []accountModel
	now := time.Now()
	err := r.db.WithContext(ctx).
		Where("status = ? AND schedulable = ?", service.StatusActive, true).
		Where("(overload_until IS NULL OR overload_until <= ?)", now).
		Where("(rate_limit_reset_at IS NULL OR rate_limit_reset_at <= ?)", now).
		Preload("Proxy").
		Order("priority ASC").
		Find(&accounts).Error
	if err != nil {
		return nil, err
placeholder
	outAccounts := make([]service.Account, 0, len(accounts))
	for i := range accounts {
		outAccounts = append(outAccounts, *accountModelToService(&accounts[i]))
placeholder
	return outAccounts, nil
placeholder

func (r *accountRepository) ListSchedulableByGroupID(ctx context.Context, groupID int64) ([]service.Account, error) {
	var accounts []accountModel
	now := time.Now()
	err := r.db.WithContext(ctx).
		Joins("JOIN account_groups ON account_groups.account_id = accounts.id").
		Where("account_groups.group_id = ?", groupID).
		Where("accounts.status = ? AND accounts.schedulable = ?", service.StatusActive, true).
		Where("(accounts.overload_until IS NULL OR accounts.overload_until <= ?)", now).
		Where("(accounts.rate_limit_reset_at IS NULL OR accounts.rate_limit_reset_at <= ?)", now).
		Preload("Proxy").
		Order("account_groups.priority ASC, accounts.priority ASC").
		Find(&accounts).Error
	if err != nil {
		return nil, err
placeholder
	outAccounts := make([]service.Account, 0, len(accounts))
	for i := range accounts {
		outAccounts = append(outAccounts, *accountModelToService(&accounts[i]))
placeholder
	return outAccounts, nil
placeholder

func (r *accountRepository) ListSchedulableByPlatform(ctx context.Context, platform string) ([]service.Account, error) {
	var accounts []accountModel
	now := time.Now()
	err := r.db.WithContext(ctx).
		Where("platform = ?", platform).
		Where("status = ? AND schedulable = ?", service.StatusActive, true).
		Where("(overload_until IS NULL OR overload_until <= ?)", now).
		Where("(rate_limit_reset_at IS NULL OR rate_limit_reset_at <= ?)", now).
		Preload("Proxy").
		Order("priority ASC").
		Find(&accounts).Error
	if err != nil {
		return nil, err
placeholder
	outAccounts := make([]service.Account, 0, len(accounts))
	for i := range accounts {
		outAccounts = append(outAccounts, *accountModelToService(&accounts[i]))
placeholder
	return outAccounts, nil
placeholder

func (r *accountRepository) ListSchedulableByGroupIDAndPlatform(ctx context.Context, groupID int64, platform string) ([]service.Account, error) {
	var accounts []accountModel
	now := time.Now()
	err := r.db.WithContext(ctx).
		Joins("JOIN account_groups ON account_groups.account_id = accounts.id").
		Where("account_groups.group_id = ?", groupID).
		Where("accounts.platform = ?", platform).
		Where("accounts.status = ? AND accounts.schedulable = ?", service.StatusActive, true).
		Where("(accounts.overload_until IS NULL OR accounts.overload_until <= ?)", now).
		Where("(accounts.rate_limit_reset_at IS NULL OR accounts.rate_limit_reset_at <= ?)", now).
		Preload("Proxy").
		Order("account_groups.priority ASC, accounts.priority ASC").
		Find(&accounts).Error
	if err != nil {
		return nil, err
placeholder
	outAccounts := make([]service.Account, 0, len(accounts))
	for i := range accounts {
		outAccounts = append(outAccounts, *accountModelToService(&accounts[i]))
placeholder
	return outAccounts, nil
placeholder

func (r *accountRepository) SetRateLimited(ctx context.Context, id int64, resetAt time.Time) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&accountModel{placeholder).Where("id = ?", id).
		Updates(map[string]any{
			"rate_limited_at":     now,
			"rate_limit_reset_at": resetAt,
	placeholder).Error
placeholder

func (r *accountRepository) SetOverloaded(ctx context.Context, id int64, until time.Time) error {
	return r.db.WithContext(ctx).Model(&accountModel{placeholder).Where("id = ?", id).
		Update("overload_until", until).Error
placeholder

func (r *accountRepository) ClearRateLimit(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&accountModel{placeholder).Where("id = ?", id).
		Updates(map[string]any{
			"rate_limited_at":     nil,
			"rate_limit_reset_at": nil,
			"overload_until":      nil,
	placeholder).Error
placeholder

func (r *accountRepository) UpdateSessionWindow(ctx context.Context, id int64, start, end *time.Time, status string) error {
	updates := map[string]any{
		"session_window_status": status,
placeholder
	if start != nil {
		updates["session_window_start"] = start
placeholder
	if end != nil {
		updates["session_window_end"] = end
placeholder
	return r.db.WithContext(ctx).Model(&accountModel{placeholder).Where("id = ?", id).Updates(updates).Error
placeholder

func (r *accountRepository) SetSchedulable(ctx context.Context, id int64, schedulable bool) error {
	return r.db.WithContext(ctx).Model(&accountModel{placeholder).Where("id = ?", id).
		Update("schedulable", schedulable).Error
placeholder

func (r *accountRepository) UpdateExtra(ctx context.Context, id int64, updates map[string]any) error {
	if len(updates) == 0 {
		return nil
placeholder

	var account accountModel
	if err := r.db.WithContext(ctx).Select("extra").Where("id = ?", id).First(&account).Error; err != nil {
		return err
placeholder

	if account.Extra == nil {
		account.Extra = datatypes.JSONMap{placeholder
placeholder
	for k, v := range updates {
		account.Extra[k] = v
placeholder

	return r.db.WithContext(ctx).Model(&accountModel{placeholder).Where("id = ?", id).
		Update("extra", account.Extra).Error
placeholder

func (r *accountRepository) BulkUpdate(ctx context.Context, ids []int64, updates service.AccountBulkUpdate) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
placeholder

	updateMap := map[string]any{placeholder

	if updates.Name != nil {
		updateMap["name"] = *updates.Name
placeholder
	if updates.ProxyID != nil {
		updateMap["proxy_id"] = updates.ProxyID
placeholder
	if updates.Concurrency != nil {
		updateMap["concurrency"] = *updates.Concurrency
placeholder
	if updates.Priority != nil {
		updateMap["priority"] = *updates.Priority
placeholder
	if updates.Status != nil {
		updateMap["status"] = *updates.Status
placeholder
	if len(updates.Credentials) > 0 {
		updateMap["credentials"] = gorm.Expr("COALESCE(credentials,'{placeholder') || ?", datatypes.JSONMap(updates.Credentials))
placeholder
	if len(updates.Extra) > 0 {
		updateMap["extra"] = gorm.Expr("COALESCE(extra,'{placeholder') || ?", datatypes.JSONMap(updates.Extra))
placeholder

	if len(updateMap) == 0 {
		return 0, nil
placeholder

	result := r.db.WithContext(ctx).
		Model(&accountModel{placeholder).
		Where("id IN ?", ids).
		Clauses(clause.Returning{placeholder).
		Updates(updateMap)

	return result.RowsAffected, result.Error
placeholder

type accountModel struct {
	ID           int64             `gorm:"primaryKey"`
	Name         string            `gorm:"size:100;not null"`
	Platform     string            `gorm:"size:50;not null"`
	Type         string            `gorm:"size:20;not null"`
	Credentials  datatypes.JSONMap `gorm:"type:jsonb;default:'{placeholder'"`
	Extra        datatypes.JSONMap `gorm:"type:jsonb;default:'{placeholder'"`
	ProxyID      *int64            `gorm:"index"`
	Concurrency  int               `gorm:"default:3;not null"`
	Priority     int               `gorm:"default:50;not null"`
	Status       string            `gorm:"size:20;default:active;not null"`
	ErrorMessage string            `gorm:"type:text"`
	LastUsedAt   *time.Time        `gorm:"index"`
	CreatedAt    time.Time         `gorm:"not null"`
	UpdatedAt    time.Time         `gorm:"not null"`
	DeletedAt    gorm.DeletedAt    `gorm:"index"`

	Schedulable bool `gorm:"default:true;not null"`

	RateLimitedAt    *time.Time `gorm:"index"`
	RateLimitResetAt *time.Time `gorm:"index"`
	OverloadUntil    *time.Time `gorm:"index"`

	SessionWindowStart  *time.Time
	SessionWindowEnd    *time.Time
	SessionWindowStatus string `gorm:"size:20"`

	Proxy         *proxyModel         `gorm:"foreignKey:ProxyID"`
	AccountGroups []accountGroupModel `gorm:"foreignKey:AccountID"`
placeholder

func (accountModel) TableName() string { return "accounts" placeholder

type accountGroupModel struct {
	AccountID int64     `gorm:"primaryKey"`
	GroupID   int64     `gorm:"primaryKey"`
	Priority  int       `gorm:"default:50;not null"`
	CreatedAt time.Time `gorm:"not null"`

	Account *accountModel `gorm:"foreignKey:AccountID"`
	Group   *groupModel   `gorm:"foreignKey:GroupID"`
placeholder

func (accountGroupModel) TableName() string { return "account_groups" placeholder

func accountGroupModelToService(m *accountGroupModel) *service.AccountGroup {
	if m == nil {
		return nil
placeholder
	return &service.AccountGroup{
		AccountID: m.AccountID,
		GroupID:   m.GroupID,
		Priority:  m.Priority,
		CreatedAt: m.CreatedAt,
		Account:   accountModelToService(m.Account),
		Group:     groupModelToService(m.Group),
placeholder
placeholder

func accountModelToService(m *accountModel) *service.Account {
	if m == nil {
		return nil
placeholder

	var credentials map[string]any
	if m.Credentials != nil {
		credentials = map[string]any(m.Credentials)
placeholder

	var extra map[string]any
	if m.Extra != nil {
		extra = map[string]any(m.Extra)
placeholder

	account := &service.Account{
		ID:                  m.ID,
		Name:                m.Name,
		Platform:            m.Platform,
		Type:                m.Type,
		Credentials:         credentials,
		Extra:               extra,
		ProxyID:             m.ProxyID,
		Concurrency:         m.Concurrency,
		Priority:            m.Priority,
		Status:              m.Status,
		ErrorMessage:        m.ErrorMessage,
		LastUsedAt:          m.LastUsedAt,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
		Schedulable:         m.Schedulable,
		RateLimitedAt:       m.RateLimitedAt,
		RateLimitResetAt:    m.RateLimitResetAt,
		OverloadUntil:       m.OverloadUntil,
		SessionWindowStart:  m.SessionWindowStart,
		SessionWindowEnd:    m.SessionWindowEnd,
		SessionWindowStatus: m.SessionWindowStatus,
		Proxy:               proxyModelToService(m.Proxy),
placeholder

	if len(m.AccountGroups) > 0 {
		account.AccountGroups = make([]service.AccountGroup, 0, len(m.AccountGroups))
		account.GroupIDs = make([]int64, 0, len(m.AccountGroups))
		account.Groups = make([]*service.Group, 0, len(m.AccountGroups))
		for i := range m.AccountGroups {
			ag := accountGroupModelToService(&m.AccountGroups[i])
			if ag == nil {
				continue
		placeholder
			account.AccountGroups = append(account.AccountGroups, *ag)
			account.GroupIDs = append(account.GroupIDs, ag.GroupID)
			if ag.Group != nil {
				account.Groups = append(account.Groups, ag.Group)
		placeholder
	placeholder
placeholder

	return account
placeholder

func accountModelFromService(a *service.Account) *accountModel {
	if a == nil {
		return nil
placeholder

	var credentials datatypes.JSONMap
	if a.Credentials != nil {
		credentials = datatypes.JSONMap(a.Credentials)
placeholder

	var extra datatypes.JSONMap
	if a.Extra != nil {
		extra = datatypes.JSONMap(a.Extra)
placeholder

	return &accountModel{
		ID:                  a.ID,
		Name:                a.Name,
		Platform:            a.Platform,
		Type:                a.Type,
		Credentials:         credentials,
		Extra:               extra,
		ProxyID:             a.ProxyID,
		Concurrency:         a.Concurrency,
		Priority:            a.Priority,
		Status:              a.Status,
		ErrorMessage:        a.ErrorMessage,
		LastUsedAt:          a.LastUsedAt,
		CreatedAt:           a.CreatedAt,
		UpdatedAt:           a.UpdatedAt,
		Schedulable:         a.Schedulable,
		RateLimitedAt:       a.RateLimitedAt,
		RateLimitResetAt:    a.RateLimitResetAt,
		OverloadUntil:       a.OverloadUntil,
		SessionWindowStart:  a.SessionWindowStart,
		SessionWindowEnd:    a.SessionWindowEnd,
		SessionWindowStatus: a.SessionWindowStatus,
placeholder
placeholder

func applyAccountModelToService(account *service.Account, m *accountModel) {
	if account == nil || m == nil {
		return
placeholder
	account.ID = m.ID
	account.CreatedAt = m.CreatedAt
	account.UpdatedAt = m.UpdatedAt
placeholder
