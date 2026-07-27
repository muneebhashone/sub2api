package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// PanelRateLimitSettings 面板 API 限流配置。
// 认证后的面板接口按「用户 ID」维度限流（与客户端 IP 无关，反向代理/共享出口
// 不会被误伤）；无需认证的公开接口按安全客户端 IP 维度限流（内网/回环地址跳过）。
type PanelRateLimitSettings struct {
	// Enabled 总开关
	Enabled bool `json:"enabled"`
	// UserRPM 每用户每分钟请求数上限（认证面板接口全量计数；0 = 不限制）
	UserRPM int `json:"user_rpm"`
	// HeavyRPM 每用户每分钟重查询上限（usage/dashboard 等聚合统计接口；0 = 不限制）
	HeavyRPM int `json:"heavy_rpm"`
	// ExemptAdmin 管理员账号是否豁免按用户限流
	ExemptAdmin bool `json:"exempt_admin"`
	// PublicIPRPM 无需认证的公开接口每 IP 每分钟上限（0 = 不限制）
	PublicIPRPM int `json:"public_ip_rpm"`
placeholder

// 面板限流 RPM 的取值上限，防止配置异常大的值失去意义。
const panelRateLimitRPMMax = 100000

const (
	panelRateLimitCacheTTL  = 60 * time.Second
	panelRateLimitErrorTTL  = 5 * time.Second
	panelRateLimitDBTimeout = 5 * time.Second
)

// cachedPanelRateLimitSettings 进程内缓存条目（60s TTL）。
type cachedPanelRateLimitSettings struct {
	settings  PanelRateLimitSettings
	expiresAt int64 // unix nano
placeholder

// DefaultPanelRateLimitSettings 返回默认面板限流配置。
// 默认启用但阈值宽松：正常前端交互远达不到，仅拦截脚本高频刷接口打爆数据库的行为。
func DefaultPanelRateLimitSettings() *PanelRateLimitSettings {
	return &PanelRateLimitSettings{
		Enabled:     true,
		UserRPM:     240,
		HeavyRPM:    60,
		ExemptAdmin: true,
		PublicIPRPM: 300,
placeholder
placeholder

// normalizePanelRateLimitSettings 修正非法取值（负数归零、超上限截断）。
func normalizePanelRateLimitSettings(s *PanelRateLimitSettings) {
	if s == nil {
		return
placeholder
	if s.UserRPM < 0 {
		s.UserRPM = 0
placeholder
	if s.HeavyRPM < 0 {
		s.HeavyRPM = 0
placeholder
	if s.PublicIPRPM < 0 {
		s.PublicIPRPM = 0
placeholder
	if s.UserRPM > panelRateLimitRPMMax {
		s.UserRPM = panelRateLimitRPMMax
placeholder
	if s.HeavyRPM > panelRateLimitRPMMax {
		s.HeavyRPM = panelRateLimitRPMMax
placeholder
	if s.PublicIPRPM > panelRateLimitRPMMax {
		s.PublicIPRPM = panelRateLimitRPMMax
placeholder
placeholder

// GetPanelRateLimitSettings 获取面板 API 限流配置（直读 DB，供管理端读写路径使用）。
// 缺失/空/解析失败 → 返回默认配置。
func (s *SettingService) GetPanelRateLimitSettings(ctx context.Context) (*PanelRateLimitSettings, error) {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyPanelRateLimitSettings)
	if err != nil {
		if errors.Is(err, ErrSettingNotFound) {
			return DefaultPanelRateLimitSettings(), nil
	placeholder
		return nil, fmt.Errorf("get panel rate limit settings: %w", err)
placeholder
	if strings.TrimSpace(value) == "" {
		return DefaultPanelRateLimitSettings(), nil
placeholder

	settings := &PanelRateLimitSettings{placeholder
	if err := json.Unmarshal([]byte(value), settings); err != nil {
		slog.Warn("failed to unmarshal panel rate limit settings, falling back to defaults",
			"error", err, "key", SettingKeyPanelRateLimitSettings)
		return DefaultPanelRateLimitSettings(), nil
placeholder
	normalizePanelRateLimitSettings(settings)
	return settings, nil
placeholder

// SetPanelRateLimitSettings 保存面板 API 限流配置，并立即刷新进程内缓存，
// 使当前节点的下一个请求即生效（多节点部署最迟 60s 内生效）。
func (s *SettingService) SetPanelRateLimitSettings(ctx context.Context, settings *PanelRateLimitSettings) error {
	if settings == nil {
		return fmt.Errorf("settings cannot be nil")
placeholder
	if settings.UserRPM < 0 || settings.HeavyRPM < 0 || settings.PublicIPRPM < 0 {
		return fmt.Errorf("rate limit values cannot be negative")
placeholder
	if settings.UserRPM > panelRateLimitRPMMax || settings.HeavyRPM > panelRateLimitRPMMax || settings.PublicIPRPM > panelRateLimitRPMMax {
		return fmt.Errorf("rate limit values must be at most %d", panelRateLimitRPMMax)
placeholder

	data, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("marshal panel rate limit settings: %w", err)
placeholder
	if err := s.settingRepo.Set(ctx, SettingKeyPanelRateLimitSettings, string(data)); err != nil {
		return err
placeholder

	s.storePanelRateLimitCache(*settings, panelRateLimitCacheTTL)
	return nil
placeholder

// GetPanelRateLimitSettingsCached 返回面板限流配置（进程内缓存，60s TTL）。
// 面板每个认证请求的热路径都会调用，绝不能每次访问 DB；
// DB 错误时返回最近一次已知值（无缓存则返回默认值），并以短 TTL 快速重试。
func (s *SettingService) GetPanelRateLimitSettingsCached(ctx context.Context) PanelRateLimitSettings {
	if s == nil || s.settingRepo == nil {
		return *DefaultPanelRateLimitSettings()
placeholder
	if cached, ok := s.panelRateLimitCache.Load().(*cachedPanelRateLimitSettings); ok && cached != nil {
		if time.Now().UnixNano() < cached.expiresAt {
			return cached.settings
	placeholder
placeholder

	result, _, _ := s.panelRateLimitSF.Do("panel_rate_limit_settings", func() (any, error) {
		// 二次检查，避免排队的 goroutine 重复查询
		if cached, ok := s.panelRateLimitCache.Load().(*cachedPanelRateLimitSettings); ok && cached != nil {
			if time.Now().UnixNano() < cached.expiresAt {
				return cached.settings, nil
		placeholder
	placeholder
		if ctx == nil {
			ctx = context.Background()
	placeholder
		// 独立 context：断开请求取消链，避免客户端断连污染缓存
		dbCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), panelRateLimitDBTimeout)
		defer cancel()

		settings, err := s.GetPanelRateLimitSettings(dbCtx)
		if err != nil {
			slog.Warn("failed to get panel rate limit settings", "error", err)
			// 保留最近一次已知值，短 TTL 快速重试
			fallback := *DefaultPanelRateLimitSettings()
			if prior, ok := s.panelRateLimitCache.Load().(*cachedPanelRateLimitSettings); ok && prior != nil {
				fallback = prior.settings
		placeholder
			s.storePanelRateLimitCache(fallback, panelRateLimitErrorTTL)
			return fallback, nil
	placeholder

		s.storePanelRateLimitCache(*settings, panelRateLimitCacheTTL)
		return *settings, nil
placeholder)
	if settings, ok := result.(PanelRateLimitSettings); ok {
		return settings
placeholder
	return *DefaultPanelRateLimitSettings()
placeholder

func (s *SettingService) storePanelRateLimitCache(settings PanelRateLimitSettings, ttl time.Duration) {
	s.panelRateLimitCache.Store(&cachedPanelRateLimitSettings{
		settings:  settings,
		expiresAt: time.Now().Add(ttl).UnixNano(),
placeholder)
placeholder
