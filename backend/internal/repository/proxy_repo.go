package repository

import (
	"context"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"

	"gorm.io/gorm"
)

type proxyRepository struct {
	db *gorm.DB
placeholder

func NewProxyRepository(db *gorm.DB) service.ProxyRepository {
	return &proxyRepository{db: dbplaceholder
placeholder

func (r *proxyRepository) Create(ctx context.Context, proxy *service.Proxy) error {
	m := proxyModelFromService(proxy)
	err := r.db.WithContext(ctx).Create(m).Error
	if err == nil {
		applyProxyModelToService(proxy, m)
placeholder
	return err
placeholder

func (r *proxyRepository) GetByID(ctx context.Context, id int64) (*service.Proxy, error) {
	var m proxyModel
	err := r.db.WithContext(ctx).First(&m, id).Error
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrProxyNotFound, nil)
placeholder
	return proxyModelToService(&m), nil
placeholder

func (r *proxyRepository) Update(ctx context.Context, proxy *service.Proxy) error {
	m := proxyModelFromService(proxy)
	err := r.db.WithContext(ctx).Save(m).Error
	if err == nil {
		applyProxyModelToService(proxy, m)
placeholder
	return err
placeholder

func (r *proxyRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&proxyModel{placeholder, id).Error
placeholder

func (r *proxyRepository) List(ctx context.Context, params pagination.PaginationParams) ([]service.Proxy, *pagination.PaginationResult, error) {
	return r.ListWithFilters(ctx, params, "", "", "")
placeholder

// ListWithFilters lists proxies with optional filtering by protocol, status, and search query
func (r *proxyRepository) ListWithFilters(ctx context.Context, params pagination.PaginationParams, protocol, status, search string) ([]service.Proxy, *pagination.PaginationResult, error) {
	var proxies []proxyModel
	var total int64

	db := r.db.WithContext(ctx).Model(&proxyModel{placeholder)

	// Apply filters
	if protocol != "" {
		db = db.Where("protocol = ?", protocol)
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

	if err := db.Offset(params.Offset()).Limit(params.Limit()).Order("id DESC").Find(&proxies).Error; err != nil {
		return nil, nil, err
placeholder

	outProxies := make([]service.Proxy, 0, len(proxies))
	for i := range proxies {
		outProxies = append(outProxies, *proxyModelToService(&proxies[i]))
placeholder

	return outProxies, paginationResultFromTotal(total, params), nil
placeholder

func (r *proxyRepository) ListActive(ctx context.Context) ([]service.Proxy, error) {
	var proxies []proxyModel
	err := r.db.WithContext(ctx).Where("status = ?", service.StatusActive).Find(&proxies).Error
	if err != nil {
		return nil, err
placeholder
	outProxies := make([]service.Proxy, 0, len(proxies))
	for i := range proxies {
		outProxies = append(outProxies, *proxyModelToService(&proxies[i]))
placeholder
	return outProxies, nil
placeholder

// ExistsByHostPortAuth checks if a proxy with the same host, port, username, and password exists
func (r *proxyRepository) ExistsByHostPortAuth(ctx context.Context, host string, port int, username, password string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&proxyModel{placeholder).
		Where("host = ? AND port = ? AND username = ? AND password = ?", host, port, username, password).
		Count(&count).Error
	if err != nil {
		return false, err
placeholder
	return count > 0, nil
placeholder

// CountAccountsByProxyID returns the number of accounts using a specific proxy
func (r *proxyRepository) CountAccountsByProxyID(ctx context.Context, proxyID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("accounts").
		Where("proxy_id = ?", proxyID).
		Count(&count).Error
	return count, err
placeholder

// GetAccountCountsForProxies returns a map of proxy ID to account count for all proxies
func (r *proxyRepository) GetAccountCountsForProxies(ctx context.Context) (map[int64]int64, error) {
	type result struct {
		ProxyID int64 `gorm:"column:proxy_id"`
		Count   int64 `gorm:"column:count"`
placeholder
	var results []result
	err := r.db.WithContext(ctx).
		Table("accounts").
		Select("proxy_id, COUNT(*) as count").
		Where("proxy_id IS NOT NULL").
		Group("proxy_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
placeholder

	counts := make(map[int64]int64)
	for _, r := range results {
		counts[r.ProxyID] = r.Count
placeholder
	return counts, nil
placeholder

// ListActiveWithAccountCount returns all active proxies with account count, sorted by creation time descending
func (r *proxyRepository) ListActiveWithAccountCount(ctx context.Context) ([]service.ProxyWithAccountCount, error) {
	var proxies []proxyModel
	err := r.db.WithContext(ctx).
		Where("status = ?", service.StatusActive).
		Order("created_at DESC").
		Find(&proxies).Error
	if err != nil {
		return nil, err
placeholder

	// Get account counts
	counts, err := r.GetAccountCountsForProxies(ctx)
	if err != nil {
		return nil, err
placeholder

	// Build result with account counts
	result := make([]service.ProxyWithAccountCount, 0, len(proxies))
	for i := range proxies {
		proxy := proxyModelToService(&proxies[i])
		if proxy == nil {
			continue
	placeholder
		result = append(result, service.ProxyWithAccountCount{
			Proxy:        *proxy,
			AccountCount: counts[proxy.ID],
	placeholder)
placeholder

	return result, nil
placeholder

type proxyModel struct {
	ID        int64          `gorm:"primaryKey"`
	Name      string         `gorm:"size:100;not null"`
	Protocol  string         `gorm:"size:20;not null"`
	Host      string         `gorm:"size:255;not null"`
	Port      int            `gorm:"not null"`
	Username  string         `gorm:"size:100"`
	Password  string         `gorm:"size:100"`
	Status    string         `gorm:"size:20;default:active;not null"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
placeholder

func (proxyModel) TableName() string { return "proxies" placeholder

func proxyModelToService(m *proxyModel) *service.Proxy {
	if m == nil {
		return nil
placeholder
	return &service.Proxy{
		ID:        m.ID,
		Name:      m.Name,
		Protocol:  m.Protocol,
		Host:      m.Host,
		Port:      m.Port,
		Username:  m.Username,
		Password:  m.Password,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
placeholder
placeholder

func proxyModelFromService(p *service.Proxy) *proxyModel {
	if p == nil {
		return nil
placeholder
	return &proxyModel{
		ID:        p.ID,
		Name:      p.Name,
		Protocol:  p.Protocol,
		Host:      p.Host,
		Port:      p.Port,
		Username:  p.Username,
		Password:  p.Password,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
placeholder
placeholder

func applyProxyModelToService(proxy *service.Proxy, m *proxyModel) {
	if proxy == nil || m == nil {
		return
placeholder
	proxy.ID = m.ID
	proxy.CreatedAt = m.CreatedAt
	proxy.UpdatedAt = m.UpdatedAt
placeholder
