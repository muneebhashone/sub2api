package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/websearch"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

// WebSearchEmulationConfig holds the global web search emulation configuration.
type WebSearchEmulationConfig struct {
	Enabled   bool                      `json:"enabled"`
	Providers []WebSearchProviderConfig `json:"providers"`
placeholder

// WebSearchProviderConfig describes a single search provider (Brave or Tavily).
type WebSearchProviderConfig struct {
	Type                 string `json:"type"`                   // websearch.ProviderTypeBrave | Tavily
	APIKey               string `json:"api_key,omitempty"`      // secret — omitted in API responses
	APIKeyConfigured     bool   `json:"api_key_configured"`     // read-only mask
	Priority             int    `json:"priority"`               // lower = higher priority
	QuotaLimit           int64  `json:"quota_limit"`            // 0 = unlimited
	QuotaRefreshInterval string `json:"quota_refresh_interval"` // websearch.QuotaRefresh*
	QuotaUsed            int64  `json:"quota_used,omitempty"`   // read-only: current period usage
	ProxyID              *int64 `json:"proxy_id"`               // optional proxy association
	ExpiresAt            *int64 `json:"expires_at,omitempty"`   // optional expiration timestamp
placeholder

// --- Validation ---

const maxWebSearchProviders = 10

var validProviderTypes = map[string]bool{
	websearch.ProviderTypeBrave:  true,
	websearch.ProviderTypeTavily: true,
placeholder

var validQuotaIntervals = map[string]bool{
	websearch.QuotaRefreshDaily:   true,
	websearch.QuotaRefreshWeekly:  true,
	websearch.QuotaRefreshMonthly: true,
	"":                            true, // defaults to monthly
placeholder

func validateWebSearchConfig(cfg *WebSearchEmulationConfig) error {
	if cfg == nil {
		return nil
placeholder
	if len(cfg.Providers) > maxWebSearchProviders {
		return fmt.Errorf("too many providers (max %d)", maxWebSearchProviders)
placeholder
	seen := make(map[string]bool, len(cfg.Providers))
	for i, p := range cfg.Providers {
		if !validProviderTypes[p.Type] {
			return fmt.Errorf("provider[%d]: invalid type %q", i, p.Type)
	placeholder
		if !validQuotaIntervals[p.QuotaRefreshInterval] {
			return fmt.Errorf("provider[%d]: invalid quota_refresh_interval %q", i, p.QuotaRefreshInterval)
	placeholder
		if p.QuotaLimit < 0 {
			return fmt.Errorf("provider[%d]: quota_limit must be >= 0", i)
	placeholder
		if seen[p.Type] {
			return fmt.Errorf("provider[%d]: duplicate type %q", i, p.Type)
	placeholder
		seen[p.Type] = true
placeholder
	return nil
placeholder

// --- In-process cache (same pattern as gateway forwarding settings) ---

const sfKeyWebSearchConfig = "web_search_emulation_config"

type cachedWebSearchEmulationConfig struct {
	config    *WebSearchEmulationConfig
	expiresAt int64 // unix nano
placeholder

var webSearchEmulationCache atomic.Value // *cachedWebSearchEmulationConfig
var webSearchEmulationSF singleflight.Group

const (
	webSearchEmulationCacheTTL  = 60 * time.Second
	webSearchEmulationErrorTTL  = 5 * time.Second
	webSearchEmulationDBTimeout = 5 * time.Second
)

// GetWebSearchEmulationConfig returns the configuration with in-process cache + singleflight.
func (s *SettingService) GetWebSearchEmulationConfig(ctx context.Context) (*WebSearchEmulationConfig, error) {
	if cached := webSearchEmulationCache.Load(); cached != nil {
		c := cached.(*cachedWebSearchEmulationConfig)
		if time.Now().UnixNano() < c.expiresAt {
			return c.config, nil
	placeholder
placeholder
	result, err, _ := webSearchEmulationSF.Do(sfKeyWebSearchConfig, func() (any, error) {
		return s.loadWebSearchConfigFromDB()
placeholder)
	if err != nil {
		return &WebSearchEmulationConfig{placeholder, err
placeholder
	return result.(*WebSearchEmulationConfig), nil
placeholder

func (s *SettingService) loadWebSearchConfigFromDB() (*WebSearchEmulationConfig, error) {
	dbCtx, cancel := context.WithTimeout(context.Background(), webSearchEmulationDBTimeout)
	defer cancel()

	raw, err := s.settingRepo.GetValue(dbCtx, SettingKeyWebSearchEmulationConfig)
	if err != nil {
		webSearchEmulationCache.Store(&cachedWebSearchEmulationConfig{
			config:    &WebSearchEmulationConfig{placeholder,
			expiresAt: time.Now().Add(webSearchEmulationErrorTTL).UnixNano(),
	placeholder)
		return &WebSearchEmulationConfig{placeholder, err
placeholder
	cfg := parseWebSearchConfigJSON(raw)
	webSearchEmulationCache.Store(&cachedWebSearchEmulationConfig{
		config:    cfg,
		expiresAt: time.Now().Add(webSearchEmulationCacheTTL).UnixNano(),
placeholder)
	return cfg, nil
placeholder

func parseWebSearchConfigJSON(raw string) *WebSearchEmulationConfig {
	cfg := &WebSearchEmulationConfig{placeholder
	if raw == "" {
		return cfg
placeholder
	if err := json.Unmarshal([]byte(raw), cfg); err != nil {
		slog.Warn("websearch: failed to parse config JSON", "error", err)
		return &WebSearchEmulationConfig{placeholder
placeholder
	return cfg
placeholder

// SaveWebSearchEmulationConfig validates and persists the configuration.
// Empty API keys in the input are preserved from the existing config.
func (s *SettingService) SaveWebSearchEmulationConfig(ctx context.Context, cfg *WebSearchEmulationConfig) error {
	if err := validateWebSearchConfig(cfg); err != nil {
		return infraerrors.BadRequest("INVALID_WEB_SEARCH_CONFIG", err.Error())
placeholder
	s.mergeExistingAPIKeys(ctx, cfg)

	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("websearch: marshal config: %w", err)
placeholder
	if err := s.settingRepo.Set(ctx, SettingKeyWebSearchEmulationConfig, string(data)); err != nil {
		return fmt.Errorf("websearch: save config: %w", err)
placeholder
	// Invalidate: forget singleflight first, then store new value
	webSearchEmulationSF.Forget(sfKeyWebSearchConfig)
	webSearchEmulationCache.Store(&cachedWebSearchEmulationConfig{
		config:    cfg,
		expiresAt: time.Now().Add(webSearchEmulationCacheTTL).UnixNano(),
placeholder)

	// Hot-reload: rebuild the global Manager with new config
	s.RebuildWebSearchManager(ctx)
	return nil
placeholder

// mergeExistingAPIKeys preserves API keys from the current config when incoming value is empty.
func (s *SettingService) mergeExistingAPIKeys(ctx context.Context, cfg *WebSearchEmulationConfig) {
	existing, _ := s.getWebSearchEmulationConfigRaw(ctx)
	if existing == nil || cfg == nil {
		return
placeholder
	existingByType := make(map[string]string, len(existing.Providers))
	for _, p := range existing.Providers {
		if p.APIKey != "" {
			existingByType[p.Type] = p.APIKey
	placeholder
placeholder
	for i := range cfg.Providers {
		if cfg.Providers[i].APIKey == "" {
			if key, ok := existingByType[cfg.Providers[i].Type]; ok {
				cfg.Providers[i].APIKey = key
		placeholder
	placeholder
placeholder
placeholder

func (s *SettingService) getWebSearchEmulationConfigRaw(ctx context.Context) (*WebSearchEmulationConfig, error) {
	raw, err := s.settingRepo.GetValue(ctx, SettingKeyWebSearchEmulationConfig)
	if err != nil {
		return nil, err
placeholder
	return parseWebSearchConfigJSON(raw), nil
placeholder

// IsWebSearchEmulationEnabled is a quick check for whether the global switch is on.
func (s *SettingService) IsWebSearchEmulationEnabled(ctx context.Context) bool {
	cfg, err := s.GetWebSearchEmulationConfig(ctx)
	if err != nil {
		return false
placeholder
	return cfg.Enabled && len(cfg.Providers) > 0
placeholder

// SetWebSearchRedisClient injects the Redis client used for quota tracking.
// Call after construction, before first use. Triggers initial Manager build.
func (s *SettingService) SetWebSearchRedisClient(ctx context.Context, redisClient *redis.Client) {
	s.webSearchRedis = redisClient
	s.RebuildWebSearchManager(ctx)
placeholder

// RebuildWebSearchManager reads the current config and (re)creates the global websearch.Manager.
// Called on startup and after SaveWebSearchEmulationConfig.
func (s *SettingService) RebuildWebSearchManager(ctx context.Context) {
	cfg, err := s.GetWebSearchEmulationConfig(ctx)
	if err != nil || !cfg.Enabled || len(cfg.Providers) == 0 {
		SetWebSearchManager(nil)
		return
placeholder
	providerConfigs := make([]websearch.ProviderConfig, 0, len(cfg.Providers))
	for _, p := range cfg.Providers {
		providerConfigs = append(providerConfigs, websearch.ProviderConfig{
			Type:                 p.Type,
			APIKey:               p.APIKey,
			Priority:             p.Priority,
			QuotaLimit:           p.QuotaLimit,
			QuotaRefreshInterval: p.QuotaRefreshInterval,
			ExpiresAt:            p.ExpiresAt,
	placeholder)
placeholder
	SetWebSearchManager(websearch.NewManager(providerConfigs, s.webSearchRedis))
	slog.Info("websearch: manager rebuilt", "provider_count", len(providerConfigs))
placeholder

// SanitizeWebSearchConfig returns a copy with api_key fields masked for API responses.
func SanitizeWebSearchConfig(cfg *WebSearchEmulationConfig) *WebSearchEmulationConfig {
	if cfg == nil {
		return nil
placeholder
	out := *cfg
	out.Providers = make([]WebSearchProviderConfig, len(cfg.Providers))
	for i, p := range cfg.Providers {
		out.Providers[i] = p
		out.Providers[i].APIKeyConfigured = p.APIKey != ""
		out.Providers[i].APIKey = "" // never return the secret
placeholder
	return &out
placeholder
