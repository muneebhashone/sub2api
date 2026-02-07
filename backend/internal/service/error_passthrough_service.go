package service

import (
	"context"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/model"
)

// ErrorPassthroughRepository 定义错误透传规则的数据访问接口
type ErrorPassthroughRepository interface {
	// List 获取所有规则
	List(ctx context.Context) ([]*model.ErrorPassthroughRule, error)
	// GetByID 根据 ID 获取规则
	GetByID(ctx context.Context, id int64) (*model.ErrorPassthroughRule, error)
	// Create 创建规则
	Create(ctx context.Context, rule *model.ErrorPassthroughRule) (*model.ErrorPassthroughRule, error)
	// Update 更新规则
	Update(ctx context.Context, rule *model.ErrorPassthroughRule) (*model.ErrorPassthroughRule, error)
	// Delete 删除规则
	Delete(ctx context.Context, id int64) error
placeholder

// ErrorPassthroughCache 定义错误透传规则的缓存接口
type ErrorPassthroughCache interface {
	// Get 从缓存获取规则列表
	Get(ctx context.Context) ([]*model.ErrorPassthroughRule, bool)
	// Set 设置缓存
	Set(ctx context.Context, rules []*model.ErrorPassthroughRule) error
	// Invalidate 使缓存失效
	Invalidate(ctx context.Context) error
	// NotifyUpdate 通知其他实例刷新缓存
	NotifyUpdate(ctx context.Context) error
	// SubscribeUpdates 订阅缓存更新通知
	SubscribeUpdates(ctx context.Context, handler func())
placeholder

// ErrorPassthroughService 错误透传规则服务
type ErrorPassthroughService struct {
	repo  ErrorPassthroughRepository
	cache ErrorPassthroughCache

	// 本地内存缓存，用于快速匹配
	localCache   []*model.ErrorPassthroughRule
	localCacheMu sync.RWMutex
placeholder

// NewErrorPassthroughService 创建错误透传规则服务
func NewErrorPassthroughService(
	repo ErrorPassthroughRepository,
	cache ErrorPassthroughCache,
) *ErrorPassthroughService {
	svc := &ErrorPassthroughService{
		repo:  repo,
		cache: cache,
placeholder

	// 启动时加载规则到本地缓存
	ctx := context.Background()
	if err := svc.reloadRulesFromDB(ctx); err != nil {
		log.Printf("[ErrorPassthroughService] Failed to load rules from DB on startup: %v", err)
		if fallbackErr := svc.refreshLocalCache(ctx); fallbackErr != nil {
			log.Printf("[ErrorPassthroughService] Failed to load rules from cache fallback on startup: %v", fallbackErr)
	placeholder
placeholder

	// 订阅缓存更新通知
	if cache != nil {
		cache.SubscribeUpdates(ctx, func() {
			if err := svc.refreshLocalCache(context.Background()); err != nil {
				log.Printf("[ErrorPassthroughService] Failed to refresh cache on notification: %v", err)
		placeholder
	placeholder)
placeholder

	return svc
placeholder

// List 获取所有规则
func (s *ErrorPassthroughService) List(ctx context.Context) ([]*model.ErrorPassthroughRule, error) {
	return s.repo.List(ctx)
placeholder

// GetByID 根据 ID 获取规则
func (s *ErrorPassthroughService) GetByID(ctx context.Context, id int64) (*model.ErrorPassthroughRule, error) {
	return s.repo.GetByID(ctx, id)
placeholder

// Create 创建规则
func (s *ErrorPassthroughService) Create(ctx context.Context, rule *model.ErrorPassthroughRule) (*model.ErrorPassthroughRule, error) {
	if err := rule.Validate(); err != nil {
		return nil, err
placeholder

	created, err := s.repo.Create(ctx, rule)
	if err != nil {
		return nil, err
placeholder

	// 刷新缓存
	refreshCtx, cancel := s.newCacheRefreshContext()
	defer cancel()
	s.invalidateAndNotify(refreshCtx)

	return created, nil
placeholder

// Update 更新规则
func (s *ErrorPassthroughService) Update(ctx context.Context, rule *model.ErrorPassthroughRule) (*model.ErrorPassthroughRule, error) {
	if err := rule.Validate(); err != nil {
		return nil, err
placeholder

	updated, err := s.repo.Update(ctx, rule)
	if err != nil {
		return nil, err
placeholder

	// 刷新缓存
	refreshCtx, cancel := s.newCacheRefreshContext()
	defer cancel()
	s.invalidateAndNotify(refreshCtx)

	return updated, nil
placeholder

// Delete 删除规则
func (s *ErrorPassthroughService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
placeholder

	// 刷新缓存
	refreshCtx, cancel := s.newCacheRefreshContext()
	defer cancel()
	s.invalidateAndNotify(refreshCtx)

	return nil
placeholder

// MatchRule 匹配透传规则
// 返回第一个匹配的规则，如果没有匹配则返回 nil
func (s *ErrorPassthroughService) MatchRule(platform string, statusCode int, body []byte) *model.ErrorPassthroughRule {
	rules := s.getCachedRules()
	if len(rules) == 0 {
		return nil
placeholder

	bodyStr := strings.ToLower(string(body))

	for _, rule := range rules {
		if !rule.Enabled {
			continue
	placeholder
		if !s.platformMatches(rule, platform) {
			continue
	placeholder
		if s.ruleMatches(rule, statusCode, bodyStr) {
			return rule
	placeholder
placeholder

	return nil
placeholder

// getCachedRules 获取缓存的规则列表（按优先级排序）
func (s *ErrorPassthroughService) getCachedRules() []*model.ErrorPassthroughRule {
	s.localCacheMu.RLock()
	rules := s.localCache
	s.localCacheMu.RUnlock()

	if rules != nil {
		return rules
placeholder

	// 如果本地缓存为空，尝试刷新
	ctx := context.Background()
	if err := s.refreshLocalCache(ctx); err != nil {
		log.Printf("[ErrorPassthroughService] Failed to refresh cache: %v", err)
		return nil
placeholder

	s.localCacheMu.RLock()
	defer s.localCacheMu.RUnlock()
	return s.localCache
placeholder

// refreshLocalCache 刷新本地缓存
func (s *ErrorPassthroughService) refreshLocalCache(ctx context.Context) error {
	// 先尝试从 Redis 缓存获取
	if s.cache != nil {
		if rules, ok := s.cache.Get(ctx); ok {
			s.setLocalCache(rules)
			return nil
	placeholder
placeholder

	return s.reloadRulesFromDB(ctx)
placeholder

// 从数据库加载（repo.List 已按 priority 排序）
// 注意：该方法会绕过 cache.Get，确保拿到数据库最新值。
func (s *ErrorPassthroughService) reloadRulesFromDB(ctx context.Context) error {
	rules, err := s.repo.List(ctx)
	if err != nil {
		return err
placeholder

	// 更新 Redis 缓存
	if s.cache != nil {
		if err := s.cache.Set(ctx, rules); err != nil {
			log.Printf("[ErrorPassthroughService] Failed to set cache: %v", err)
	placeholder
placeholder

	// 更新本地缓存（setLocalCache 内部会确保排序）
	s.setLocalCache(rules)

	return nil
placeholder

// setLocalCache 设置本地缓存
func (s *ErrorPassthroughService) setLocalCache(rules []*model.ErrorPassthroughRule) {
	// 按优先级排序
	sorted := make([]*model.ErrorPassthroughRule, len(rules))
	copy(sorted, rules)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Priority < sorted[j].Priority
placeholder)

	s.localCacheMu.Lock()
	s.localCache = sorted
	s.localCacheMu.Unlock()
placeholder

// clearLocalCache 清空本地缓存，避免刷新失败时继续命中陈旧规则。
func (s *ErrorPassthroughService) clearLocalCache() {
	s.localCacheMu.Lock()
	s.localCache = nil
	s.localCacheMu.Unlock()
placeholder

// newCacheRefreshContext 为写路径缓存同步创建独立上下文，避免受请求取消影响。
func (s *ErrorPassthroughService) newCacheRefreshContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
placeholder

// invalidateAndNotify 使缓存失效并通知其他实例
func (s *ErrorPassthroughService) invalidateAndNotify(ctx context.Context) {
	// 先失效缓存，避免后续刷新读到陈旧规则。
	if s.cache != nil {
		if err := s.cache.Invalidate(ctx); err != nil {
			log.Printf("[ErrorPassthroughService] Failed to invalidate cache: %v", err)
	placeholder
placeholder

	// 刷新本地缓存
	if err := s.reloadRulesFromDB(ctx); err != nil {
		log.Printf("[ErrorPassthroughService] Failed to refresh local cache: %v", err)
		// 刷新失败时清空本地缓存，避免继续使用陈旧规则。
		s.clearLocalCache()
placeholder

	// 通知其他实例
	if s.cache != nil {
		if err := s.cache.NotifyUpdate(ctx); err != nil {
			log.Printf("[ErrorPassthroughService] Failed to notify cache update: %v", err)
	placeholder
placeholder
placeholder

// platformMatches 检查平台是否匹配
func (s *ErrorPassthroughService) platformMatches(rule *model.ErrorPassthroughRule, platform string) bool {
	// 如果没有配置平台限制，则匹配所有平台
	if len(rule.Platforms) == 0 {
		return true
placeholder

	platform = strings.ToLower(platform)
	for _, p := range rule.Platforms {
		if strings.ToLower(p) == platform {
			return true
	placeholder
placeholder

	return false
placeholder

// ruleMatches 检查规则是否匹配
func (s *ErrorPassthroughService) ruleMatches(rule *model.ErrorPassthroughRule, statusCode int, bodyLower string) bool {
	hasErrorCodes := len(rule.ErrorCodes) > 0
	hasKeywords := len(rule.Keywords) > 0

	// 如果没有配置任何条件，不匹配
	if !hasErrorCodes && !hasKeywords {
		return false
placeholder

	codeMatch := !hasErrorCodes || s.containsInt(rule.ErrorCodes, statusCode)
	keywordMatch := !hasKeywords || s.containsAnyKeyword(bodyLower, rule.Keywords)

	if rule.MatchMode == model.MatchModeAll {
		// "all" 模式：所有配置的条件都必须满足
		return codeMatch && keywordMatch
placeholder

	// "any" 模式：任一条件满足即可
	if hasErrorCodes && hasKeywords {
		return codeMatch || keywordMatch
placeholder
	return codeMatch && keywordMatch
placeholder

// containsInt 检查切片是否包含指定整数
func (s *ErrorPassthroughService) containsInt(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
	placeholder
placeholder
	return false
placeholder

// containsAnyKeyword 检查字符串是否包含任一关键词（不区分大小写）
func (s *ErrorPassthroughService) containsAnyKeyword(bodyLower string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.Contains(bodyLower, strings.ToLower(kw)) {
			return true
	placeholder
placeholder
	return false
placeholder
