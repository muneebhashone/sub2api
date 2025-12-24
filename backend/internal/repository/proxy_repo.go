package repository

import (
	"context"
	"github.com/Wei-Shaw/sub2api/internal/model"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"

	"gorm.io/gorm"
)

type ProxyRepository struct {
	db *gorm.DB
placeholder

func NewProxyRepository(db *gorm.DB) *ProxyRepository {
	return &ProxyRepository{db: dbplaceholder
placeholder

func (r *ProxyRepository) Create(ctx context.Context, proxy *model.Proxy) error {
	return r.db.WithContext(ctx).Create(proxy).Error
placeholder

func (r *ProxyRepository) GetByID(ctx context.Context, id int64) (*model.Proxy, error) {
	var proxy model.Proxy
	err := r.db.WithContext(ctx).First(&proxy, id).Error
	if err != nil {
		return nil, err
placeholder
	return &proxy, nil
placeholder

func (r *ProxyRepository) Update(ctx context.Context, proxy *model.Proxy) error {
	return r.db.WithContext(ctx).Save(proxy).Error
placeholder

func (r *ProxyRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Proxy{placeholder, id).Error
placeholder

func (r *ProxyRepository) List(ctx context.Context, params pagination.PaginationParams) ([]model.Proxy, *pagination.PaginationResult, error) {
	return r.ListWithFilters(ctx, params, "", "", "")
placeholder

// ListWithFilters lists proxies with optional filtering by protocol, status, and search query
func (r *ProxyRepository) ListWithFilters(ctx context.Context, params pagination.PaginationParams, protocol, status, search string) ([]model.Proxy, *pagination.PaginationResult, error) {
	var proxies []model.Proxy
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Proxy{placeholder)

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

	pages := int(total) / params.Limit()
	if int(total)%params.Limit() > 0 {
		pages++
placeholder

	return proxies, &pagination.PaginationResult{
		Total:    total,
		Page:     params.Page,
		PageSize: params.Limit(),
		Pages:    pages,
placeholder, nil
placeholder

func (r *ProxyRepository) ListActive(ctx context.Context) ([]model.Proxy, error) {
	var proxies []model.Proxy
	err := r.db.WithContext(ctx).Where("status = ?", model.StatusActive).Find(&proxies).Error
	return proxies, err
placeholder

// ExistsByHostPortAuth checks if a proxy with the same host, port, username, and password exists
func (r *ProxyRepository) ExistsByHostPortAuth(ctx context.Context, host string, port int, username, password string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Proxy{placeholder).
		Where("host = ? AND port = ? AND username = ? AND password = ?", host, port, username, password).
		Count(&count).Error
	if err != nil {
		return false, err
placeholder
	return count > 0, nil
placeholder

// CountAccountsByProxyID returns the number of accounts using a specific proxy
func (r *ProxyRepository) CountAccountsByProxyID(ctx context.Context, proxyID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Account{placeholder).
		Where("proxy_id = ?", proxyID).
		Count(&count).Error
	return count, err
placeholder

// GetAccountCountsForProxies returns a map of proxy ID to account count for all proxies
func (r *ProxyRepository) GetAccountCountsForProxies(ctx context.Context) (map[int64]int64, error) {
	type result struct {
		ProxyID int64 `gorm:"column:proxy_id"`
		Count   int64 `gorm:"column:count"`
placeholder
	var results []result
	err := r.db.WithContext(ctx).
		Model(&model.Account{placeholder).
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
func (r *ProxyRepository) ListActiveWithAccountCount(ctx context.Context) ([]model.ProxyWithAccountCount, error) {
	var proxies []model.Proxy
	err := r.db.WithContext(ctx).
		Where("status = ?", model.StatusActive).
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
	result := make([]model.ProxyWithAccountCount, len(proxies))
	for i, proxy := range proxies {
		result[i] = model.ProxyWithAccountCount{
			Proxy:        proxy,
			AccountCount: counts[proxy.ID],
	placeholder
placeholder

	return result, nil
placeholder
