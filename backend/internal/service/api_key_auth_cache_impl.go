package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/dgraph-io/ristretto"
)

type apiKeyAuthCacheConfig struct {
	l1Size        int
	l1TTL         time.Duration
	l2TTL         time.Duration
	negativeTTL   time.Duration
	jitterPercent int
	singleflight  bool
placeholder

var (
	jitterRandMu sync.Mutex
	// 认证缓存抖动使用独立随机源，避免全局 Seed
	jitterRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func newAPIKeyAuthCacheConfig(cfg *config.Config) apiKeyAuthCacheConfig {
	if cfg == nil {
		return apiKeyAuthCacheConfig{placeholder
placeholder
	auth := cfg.APIKeyAuth
	return apiKeyAuthCacheConfig{
		l1Size:        auth.L1Size,
		l1TTL:         time.Duration(auth.L1TTLSeconds) * time.Second,
		l2TTL:         time.Duration(auth.L2TTLSeconds) * time.Second,
		negativeTTL:   time.Duration(auth.NegativeTTLSeconds) * time.Second,
		jitterPercent: auth.JitterPercent,
		singleflight:  auth.Singleflight,
placeholder
placeholder

func (c apiKeyAuthCacheConfig) l1Enabled() bool {
	return c.l1Size > 0 && c.l1TTL > 0
placeholder

func (c apiKeyAuthCacheConfig) l2Enabled() bool {
	return c.l2TTL > 0
placeholder

func (c apiKeyAuthCacheConfig) negativeEnabled() bool {
	return c.negativeTTL > 0
placeholder

func (c apiKeyAuthCacheConfig) jitterTTL(ttl time.Duration) time.Duration {
	if ttl <= 0 {
		return ttl
placeholder
	if c.jitterPercent <= 0 {
		return ttl
placeholder
	percent := c.jitterPercent
	if percent > 100 {
		percent = 100
placeholder
	delta := float64(percent) / 100
	jitterRandMu.Lock()
	randVal := jitterRand.Float64()
	jitterRandMu.Unlock()
	factor := 1 - delta + randVal*(2*delta)
	if factor <= 0 {
		return ttl
placeholder
	return time.Duration(float64(ttl) * factor)
placeholder

func (s *APIKeyService) initAuthCache(cfg *config.Config) {
	s.authCfg = newAPIKeyAuthCacheConfig(cfg)
	if !s.authCfg.l1Enabled() {
		return
placeholder
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: int64(s.authCfg.l1Size) * 10,
		MaxCost:     int64(s.authCfg.l1Size),
		BufferItems: 64,
placeholder)
	if err != nil {
		return
placeholder
	s.authCacheL1 = cache
placeholder

func (s *APIKeyService) authCacheKey(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
placeholder

func (s *APIKeyService) getAuthCacheEntry(ctx context.Context, cacheKey string) (*APIKeyAuthCacheEntry, bool) {
	if s.authCacheL1 != nil {
		if val, ok := s.authCacheL1.Get(cacheKey); ok {
			if entry, ok := val.(*APIKeyAuthCacheEntry); ok {
				return entry, true
		placeholder
	placeholder
placeholder
	if s.cache == nil || !s.authCfg.l2Enabled() {
		return nil, false
placeholder
	entry, err := s.cache.GetAuthCache(ctx, cacheKey)
	if err != nil {
		return nil, false
placeholder
	s.setAuthCacheL1(cacheKey, entry)
	return entry, true
placeholder

func (s *APIKeyService) setAuthCacheL1(cacheKey string, entry *APIKeyAuthCacheEntry) {
	if s.authCacheL1 == nil || entry == nil {
		return
placeholder
	ttl := s.authCfg.l1TTL
	if entry.NotFound && s.authCfg.negativeTTL > 0 && s.authCfg.negativeTTL < ttl {
		ttl = s.authCfg.negativeTTL
placeholder
	ttl = s.authCfg.jitterTTL(ttl)
	_ = s.authCacheL1.SetWithTTL(cacheKey, entry, 1, ttl)
placeholder

func (s *APIKeyService) setAuthCacheEntry(ctx context.Context, cacheKey string, entry *APIKeyAuthCacheEntry, ttl time.Duration) {
	if entry == nil {
		return
placeholder
	s.setAuthCacheL1(cacheKey, entry)
	if s.cache == nil || !s.authCfg.l2Enabled() {
		return
placeholder
	_ = s.cache.SetAuthCache(ctx, cacheKey, entry, s.authCfg.jitterTTL(ttl))
placeholder

func (s *APIKeyService) deleteAuthCache(ctx context.Context, cacheKey string) {
	if s.authCacheL1 != nil {
		s.authCacheL1.Del(cacheKey)
placeholder
	if s.cache == nil {
		return
placeholder
	_ = s.cache.DeleteAuthCache(ctx, cacheKey)
placeholder

func (s *APIKeyService) loadAuthCacheEntry(ctx context.Context, key, cacheKey string) (*APIKeyAuthCacheEntry, error) {
	apiKey, err := s.apiKeyRepo.GetByKeyForAuth(ctx, key)
	if err != nil {
		if errors.Is(err, ErrAPIKeyNotFound) {
			entry := &APIKeyAuthCacheEntry{NotFound: trueplaceholder
			if s.authCfg.negativeEnabled() {
				s.setAuthCacheEntry(ctx, cacheKey, entry, s.authCfg.negativeTTL)
		placeholder
			return entry, nil
	placeholder
		return nil, fmt.Errorf("get api key: %w", err)
placeholder
	apiKey.Key = key
	snapshot := s.snapshotFromAPIKey(apiKey)
	if snapshot == nil {
		return nil, fmt.Errorf("get api key: %w", ErrAPIKeyNotFound)
placeholder
	entry := &APIKeyAuthCacheEntry{Snapshot: snapshotplaceholder
	s.setAuthCacheEntry(ctx, cacheKey, entry, s.authCfg.l2TTL)
	return entry, nil
placeholder

func (s *APIKeyService) applyAuthCacheEntry(key string, entry *APIKeyAuthCacheEntry) (*APIKey, bool, error) {
	if entry == nil {
		return nil, false, nil
placeholder
	if entry.NotFound {
		return nil, true, ErrAPIKeyNotFound
placeholder
	if entry.Snapshot == nil {
		return nil, false, nil
placeholder
	return s.snapshotToAPIKey(key, entry.Snapshot), true, nil
placeholder

func (s *APIKeyService) snapshotFromAPIKey(apiKey *APIKey) *APIKeyAuthSnapshot {
	if apiKey == nil || apiKey.User == nil {
		return nil
placeholder
	snapshot := &APIKeyAuthSnapshot{
		APIKeyID:    apiKey.ID,
		UserID:      apiKey.UserID,
		GroupID:     apiKey.GroupID,
		Status:      apiKey.Status,
		IPWhitelist: apiKey.IPWhitelist,
		IPBlacklist: apiKey.IPBlacklist,
		User: APIKeyAuthUserSnapshot{
			ID:          apiKey.User.ID,
			Status:      apiKey.User.Status,
			Role:        apiKey.User.Role,
			Balance:     apiKey.User.Balance,
			Concurrency: apiKey.User.Concurrency,
	placeholder,
placeholder
	if apiKey.Group != nil {
		snapshot.Group = &APIKeyAuthGroupSnapshot{
			ID:               apiKey.Group.ID,
			Name:             apiKey.Group.Name,
			Platform:         apiKey.Group.Platform,
			Status:           apiKey.Group.Status,
			SubscriptionType: apiKey.Group.SubscriptionType,
			RateMultiplier:   apiKey.Group.RateMultiplier,
			DailyLimitUSD:    apiKey.Group.DailyLimitUSD,
			WeeklyLimitUSD:   apiKey.Group.WeeklyLimitUSD,
			MonthlyLimitUSD:  apiKey.Group.MonthlyLimitUSD,
			ImagePrice1K:     apiKey.Group.ImagePrice1K,
			ImagePrice2K:     apiKey.Group.ImagePrice2K,
			ImagePrice4K:     apiKey.Group.ImagePrice4K,
			ClaudeCodeOnly:   apiKey.Group.ClaudeCodeOnly,
			FallbackGroupID:  apiKey.Group.FallbackGroupID,
	placeholder
placeholder
	return snapshot
placeholder

func (s *APIKeyService) snapshotToAPIKey(key string, snapshot *APIKeyAuthSnapshot) *APIKey {
	if snapshot == nil {
		return nil
placeholder
	apiKey := &APIKey{
		ID:          snapshot.APIKeyID,
		UserID:      snapshot.UserID,
		GroupID:     snapshot.GroupID,
		Key:         key,
		Status:      snapshot.Status,
		IPWhitelist: snapshot.IPWhitelist,
		IPBlacklist: snapshot.IPBlacklist,
		User: &User{
			ID:          snapshot.User.ID,
			Status:      snapshot.User.Status,
			Role:        snapshot.User.Role,
			Balance:     snapshot.User.Balance,
			Concurrency: snapshot.User.Concurrency,
	placeholder,
placeholder
	if snapshot.Group != nil {
		apiKey.Group = &Group{
			ID:               snapshot.Group.ID,
			Name:             snapshot.Group.Name,
			Platform:         snapshot.Group.Platform,
			Status:           snapshot.Group.Status,
			Hydrated:         true,
			SubscriptionType: snapshot.Group.SubscriptionType,
			RateMultiplier:   snapshot.Group.RateMultiplier,
			DailyLimitUSD:    snapshot.Group.DailyLimitUSD,
			WeeklyLimitUSD:   snapshot.Group.WeeklyLimitUSD,
			MonthlyLimitUSD:  snapshot.Group.MonthlyLimitUSD,
			ImagePrice1K:     snapshot.Group.ImagePrice1K,
			ImagePrice2K:     snapshot.Group.ImagePrice2K,
			ImagePrice4K:     snapshot.Group.ImagePrice4K,
			ClaudeCodeOnly:   snapshot.Group.ClaudeCodeOnly,
			FallbackGroupID:  snapshot.Group.FallbackGroupID,
	placeholder
placeholder
	return apiKey
placeholder
