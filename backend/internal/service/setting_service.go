package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/imroc/req/v3"
	"golang.org/x/sync/singleflight"
)

var (
	ErrRegistrationDisabled   = infraerrors.Forbidden("REGISTRATION_DISABLED", "registration is currently disabled")
	ErrSettingNotFound        = infraerrors.NotFound("SETTING_NOT_FOUND", "setting not found")
	ErrDefaultSubGroupInvalid = infraerrors.BadRequest(
		"DEFAULT_SUBSCRIPTION_GROUP_INVALID",
		"default subscription group must exist and be subscription type",
	)
	ErrDefaultSubGroupDuplicate = infraerrors.BadRequest(
		"DEFAULT_SUBSCRIPTION_GROUP_DUPLICATE",
		"default subscription group cannot be duplicated",
	)
)

type SettingRepository interface {
	Get(ctx context.Context, key string) (*Setting, error)
	GetValue(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	GetMultiple(ctx context.Context, keys []string) (map[string]string, error)
	SetMultiple(ctx context.Context, settings map[string]string) error
	GetAll(ctx context.Context) (map[string]string, error)
	Delete(ctx context.Context, key string) error
placeholder

// cachedVersionBounds 缓存 Claude Code 版本号上下限（进程内缓存，60s TTL）
type cachedVersionBounds struct {
	min       string // 空字符串 = 不检查
	max       string // 空字符串 = 不检查
	expiresAt int64  // unix nano
placeholder

// versionBoundsCache 版本号上下限进程内缓存
var versionBoundsCache atomic.Value // *cachedVersionBounds

// versionBoundsSF 防止缓存过期时 thundering herd
var versionBoundsSF singleflight.Group

// versionBoundsCacheTTL 缓存有效期
const versionBoundsCacheTTL = 60 * time.Second

// versionBoundsErrorTTL DB 错误时的短缓存，快速重试
const versionBoundsErrorTTL = 5 * time.Second

// versionBoundsDBTimeout singleflight 内 DB 查询超时，独立于请求 context
const versionBoundsDBTimeout = 5 * time.Second

// cachedBackendMode Backend Mode cache (in-process, 60s TTL)
type cachedBackendMode struct {
	value     bool
	expiresAt int64 // unix nano
placeholder

var backendModeCache atomic.Value // *cachedBackendMode
var backendModeSF singleflight.Group

const backendModeCacheTTL = 60 * time.Second
const backendModeErrorTTL = 5 * time.Second
const backendModeDBTimeout = 5 * time.Second

// cachedGatewayForwardingSettings 缓存网关转发行为设置（进程内缓存，60s TTL）
type cachedGatewayForwardingSettings struct {
	fingerprintUnification bool
	metadataPassthrough    bool
	cchSigning             bool
	expiresAt              int64 // unix nano
placeholder

var gatewayForwardingCache atomic.Value // *cachedGatewayForwardingSettings
var gatewayForwardingSF singleflight.Group

const gatewayForwardingCacheTTL = 60 * time.Second
const gatewayForwardingErrorTTL = 5 * time.Second
const gatewayForwardingDBTimeout = 5 * time.Second

// DefaultSubscriptionGroupReader validates group references used by default subscriptions.
type DefaultSubscriptionGroupReader interface {
	GetByID(ctx context.Context, id int64) (*Group, error)
placeholder

// SettingService 系统设置服务
type SettingService struct {
	settingRepo           SettingRepository
	defaultSubGroupReader DefaultSubscriptionGroupReader
	cfg                   *config.Config
	onUpdate              func() // Callback when settings are updated (for cache invalidation)
	version               string // Application version
placeholder

// NewSettingService 创建系统设置服务实例
func NewSettingService(settingRepo SettingRepository, cfg *config.Config) *SettingService {
	return &SettingService{
		settingRepo: settingRepo,
		cfg:         cfg,
placeholder
placeholder

// SetDefaultSubscriptionGroupReader injects an optional group reader for default subscription validation.
func (s *SettingService) SetDefaultSubscriptionGroupReader(reader DefaultSubscriptionGroupReader) {
	s.defaultSubGroupReader = reader
placeholder

// GetAllSettings 获取所有系统设置
func (s *SettingService) GetAllSettings(ctx context.Context) (*SystemSettings, error) {
	settings, err := s.settingRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all settings: %w", err)
placeholder

	return s.parseSettings(settings), nil
placeholder

// GetFrontendURL 获取前端基础URL（数据库优先，fallback 到配置文件）
func (s *SettingService) GetFrontendURL(ctx context.Context) string {
	val, err := s.settingRepo.GetValue(ctx, SettingKeyFrontendURL)
	if err == nil && strings.TrimSpace(val) != "" {
		return strings.TrimSpace(val)
placeholder
	return s.cfg.Server.FrontendURL
placeholder

// GetPublicSettings 获取公开设置（无需登录）
func (s *SettingService) GetPublicSettings(ctx context.Context) (*PublicSettings, error) {
	keys := []string{
		SettingKeyRegistrationEnabled,
		SettingKeyEmailVerifyEnabled,
		SettingKeyRegistrationEmailSuffixWhitelist,
		SettingKeyPromoCodeEnabled,
		SettingKeyPasswordResetEnabled,
		SettingKeyInvitationCodeEnabled,
		SettingKeyTotpEnabled,
		SettingKeyTurnstileEnabled,
		SettingKeyTurnstileSiteKey,
		SettingKeySiteName,
		SettingKeySiteLogo,
		SettingKeySiteSubtitle,
		SettingKeyAPIBaseURL,
		SettingKeyContactInfo,
		SettingKeyDocURL,
		SettingKeyHomeContent,
		SettingKeyHideCcsImportButton,
		SettingKeyPurchaseSubscriptionEnabled,
		SettingKeyPurchaseSubscriptionURL,
		SettingKeyTableDefaultPageSize,
		SettingKeyTablePageSizeOptions,
		SettingKeyCustomMenuItems,
		SettingKeyCustomEndpoints,
		SettingKeyLinuxDoConnectEnabled,
		SettingKeyBackendModeEnabled,
		SettingPaymentEnabled,
		SettingKeyOIDCConnectEnabled,
		SettingKeyOIDCConnectProviderName,
placeholder

	settings, err := s.settingRepo.GetMultiple(ctx, keys)
	if err != nil {
		return nil, fmt.Errorf("get public settings: %w", err)
placeholder

	linuxDoEnabled := false
	if raw, ok := settings[SettingKeyLinuxDoConnectEnabled]; ok {
		linuxDoEnabled = raw == "true"
placeholder else {
		linuxDoEnabled = s.cfg != nil && s.cfg.LinuxDo.Enabled
placeholder
	oidcEnabled := false
	if raw, ok := settings[SettingKeyOIDCConnectEnabled]; ok {
		oidcEnabled = raw == "true"
placeholder else {
		oidcEnabled = s.cfg != nil && s.cfg.OIDC.Enabled
placeholder
	oidcProviderName := strings.TrimSpace(settings[SettingKeyOIDCConnectProviderName])
	if oidcProviderName == "" && s.cfg != nil {
		oidcProviderName = strings.TrimSpace(s.cfg.OIDC.ProviderName)
placeholder
	if oidcProviderName == "" {
		oidcProviderName = "OIDC"
placeholder

	// Password reset requires email verification to be enabled
	emailVerifyEnabled := settings[SettingKeyEmailVerifyEnabled] == "true"
	passwordResetEnabled := emailVerifyEnabled && settings[SettingKeyPasswordResetEnabled] == "true"
	registrationEmailSuffixWhitelist := ParseRegistrationEmailSuffixWhitelist(
		settings[SettingKeyRegistrationEmailSuffixWhitelist],
	)
	tableDefaultPageSize, tablePageSizeOptions := parseTablePreferences(
		settings[SettingKeyTableDefaultPageSize],
		settings[SettingKeyTablePageSizeOptions],
	)

	return &PublicSettings{
		RegistrationEnabled:              settings[SettingKeyRegistrationEnabled] == "true",
		EmailVerifyEnabled:               emailVerifyEnabled,
		RegistrationEmailSuffixWhitelist: registrationEmailSuffixWhitelist,
		PromoCodeEnabled:                 settings[SettingKeyPromoCodeEnabled] != "false", // 默认启用
		PasswordResetEnabled:             passwordResetEnabled,
		InvitationCodeEnabled:            settings[SettingKeyInvitationCodeEnabled] == "true",
		TotpEnabled:                      settings[SettingKeyTotpEnabled] == "true",
		TurnstileEnabled:                 settings[SettingKeyTurnstileEnabled] == "true",
		TurnstileSiteKey:                 settings[SettingKeyTurnstileSiteKey],
		SiteName:                         s.getStringOrDefault(settings, SettingKeySiteName, "Sub2API"),
		SiteLogo:                         settings[SettingKeySiteLogo],
		SiteSubtitle:                     s.getStringOrDefault(settings, SettingKeySiteSubtitle, "Subscription to API Conversion Platform"),
		APIBaseURL:                       settings[SettingKeyAPIBaseURL],
		ContactInfo:                      settings[SettingKeyContactInfo],
		DocURL:                           settings[SettingKeyDocURL],
		HomeContent:                      settings[SettingKeyHomeContent],
		HideCcsImportButton:              settings[SettingKeyHideCcsImportButton] == "true",
		PurchaseSubscriptionEnabled:      settings[SettingKeyPurchaseSubscriptionEnabled] == "true",
		PurchaseSubscriptionURL:          strings.TrimSpace(settings[SettingKeyPurchaseSubscriptionURL]),
		TableDefaultPageSize:             tableDefaultPageSize,
		TablePageSizeOptions:             tablePageSizeOptions,
		CustomMenuItems:                  settings[SettingKeyCustomMenuItems],
		CustomEndpoints:                  settings[SettingKeyCustomEndpoints],
		LinuxDoOAuthEnabled:              linuxDoEnabled,
		BackendModeEnabled:               settings[SettingKeyBackendModeEnabled] == "true",
		PaymentEnabled:                   settings[SettingPaymentEnabled] == "true",
		OIDCOAuthEnabled:                 oidcEnabled,
		OIDCOAuthProviderName:            oidcProviderName,
placeholder, nil
placeholder

// SetOnUpdateCallback sets a callback function to be called when settings are updated
// This is used for cache invalidation (e.g., HTML cache in frontend server)
func (s *SettingService) SetOnUpdateCallback(callback func()) {
	s.onUpdate = callback
placeholder

// SetVersion sets the application version for injection into public settings
func (s *SettingService) SetVersion(version string) {
	s.version = version
placeholder

// GetPublicSettingsForInjection returns public settings in a format suitable for HTML injection
// This implements the web.PublicSettingsProvider interface
func (s *SettingService) GetPublicSettingsForInjection(ctx context.Context) (any, error) {
	settings, err := s.GetPublicSettings(ctx)
	if err != nil {
		return nil, err
placeholder

	// Return a struct that matches the frontend's expected format
	return &struct {
		RegistrationEnabled              bool            `json:"registration_enabled"`
		EmailVerifyEnabled               bool            `json:"email_verify_enabled"`
		RegistrationEmailSuffixWhitelist []string        `json:"registration_email_suffix_whitelist"`
		PromoCodeEnabled                 bool            `json:"promo_code_enabled"`
		PasswordResetEnabled             bool            `json:"password_reset_enabled"`
		InvitationCodeEnabled            bool            `json:"invitation_code_enabled"`
		TotpEnabled                      bool            `json:"totp_enabled"`
		TurnstileEnabled                 bool            `json:"turnstile_enabled"`
		TurnstileSiteKey                 string          `json:"turnstile_site_key,omitempty"`
		SiteName                         string          `json:"site_name"`
		SiteLogo                         string          `json:"site_logo,omitempty"`
		SiteSubtitle                     string          `json:"site_subtitle,omitempty"`
		APIBaseURL                       string          `json:"api_base_url,omitempty"`
		ContactInfo                      string          `json:"contact_info,omitempty"`
		DocURL                           string          `json:"doc_url,omitempty"`
		HomeContent                      string          `json:"home_content,omitempty"`
		HideCcsImportButton              bool            `json:"hide_ccs_import_button"`
		PurchaseSubscriptionEnabled      bool            `json:"purchase_subscription_enabled"`
		PurchaseSubscriptionURL          string          `json:"purchase_subscription_url,omitempty"`
		TableDefaultPageSize             int             `json:"table_default_page_size"`
		TablePageSizeOptions             []int           `json:"table_page_size_options"`
		CustomMenuItems                  json.RawMessage `json:"custom_menu_items"`
		CustomEndpoints                  json.RawMessage `json:"custom_endpoints"`
		LinuxDoOAuthEnabled              bool            `json:"linuxdo_oauth_enabled"`
		BackendModeEnabled               bool            `json:"backend_mode_enabled"`
		PaymentEnabled                   bool            `json:"payment_enabled"`
		OIDCOAuthEnabled                 bool            `json:"oidc_oauth_enabled"`
		OIDCOAuthProviderName            string          `json:"oidc_oauth_provider_name"`
		Version                          string          `json:"version,omitempty"`
placeholder{
		RegistrationEnabled:              settings.RegistrationEnabled,
		EmailVerifyEnabled:               settings.EmailVerifyEnabled,
		RegistrationEmailSuffixWhitelist: settings.RegistrationEmailSuffixWhitelist,
		PromoCodeEnabled:                 settings.PromoCodeEnabled,
		PasswordResetEnabled:             settings.PasswordResetEnabled,
		InvitationCodeEnabled:            settings.InvitationCodeEnabled,
		TotpEnabled:                      settings.TotpEnabled,
		TurnstileEnabled:                 settings.TurnstileEnabled,
		TurnstileSiteKey:                 settings.TurnstileSiteKey,
		SiteName:                         settings.SiteName,
		SiteLogo:                         settings.SiteLogo,
		SiteSubtitle:                     settings.SiteSubtitle,
		APIBaseURL:                       settings.APIBaseURL,
		ContactInfo:                      settings.ContactInfo,
		DocURL:                           settings.DocURL,
		HomeContent:                      settings.HomeContent,
		HideCcsImportButton:              settings.HideCcsImportButton,
		PurchaseSubscriptionEnabled:      settings.PurchaseSubscriptionEnabled,
		PurchaseSubscriptionURL:          settings.PurchaseSubscriptionURL,
		TableDefaultPageSize:             settings.TableDefaultPageSize,
		TablePageSizeOptions:             settings.TablePageSizeOptions,
		CustomMenuItems:                  filterUserVisibleMenuItems(settings.CustomMenuItems),
		CustomEndpoints:                  safeRawJSONArray(settings.CustomEndpoints),
		LinuxDoOAuthEnabled:              settings.LinuxDoOAuthEnabled,
		BackendModeEnabled:               settings.BackendModeEnabled,
		PaymentEnabled:                   settings.PaymentEnabled,
		OIDCOAuthEnabled:                 settings.OIDCOAuthEnabled,
		OIDCOAuthProviderName:            settings.OIDCOAuthProviderName,
		Version:                          s.version,
placeholder, nil
placeholder

// filterUserVisibleMenuItems filters out admin-only menu items from a raw JSON
// array string, returning only items with visibility != "admin".
func filterUserVisibleMenuItems(raw string) json.RawMessage {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "[]" {
		return json.RawMessage("[]")
placeholder
	var items []struct {
		Visibility string `json:"visibility"`
placeholder
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return json.RawMessage("[]")
placeholder

	// Parse full items to preserve all fields
	var fullItems []json.RawMessage
	if err := json.Unmarshal([]byte(raw), &fullItems); err != nil {
		return json.RawMessage("[]")
placeholder

	var filtered []json.RawMessage
	for i, item := range items {
		if item.Visibility != "admin" {
			filtered = append(filtered, fullItems[i])
	placeholder
placeholder
	if len(filtered) == 0 {
		return json.RawMessage("[]")
placeholder
	result, err := json.Marshal(filtered)
	if err != nil {
		return json.RawMessage("[]")
placeholder
	return result
placeholder

// safeRawJSONArray returns raw as json.RawMessage if it's valid JSON, otherwise "[]".
func safeRawJSONArray(raw string) json.RawMessage {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return json.RawMessage("[]")
placeholder
	if json.Valid([]byte(raw)) {
		return json.RawMessage(raw)
placeholder
	return json.RawMessage("[]")
placeholder

// GetFrameSrcOrigins returns deduplicated http(s) origins from home_content URL,
// purchase_subscription_url, and all custom_menu_items URLs. Used by the router layer for CSP frame-src injection.
func (s *SettingService) GetFrameSrcOrigins(ctx context.Context) ([]string, error) {
	settings, err := s.GetPublicSettings(ctx)
	if err != nil {
		return nil, err
placeholder

	seen := make(map[string]struct{placeholder)
	var origins []string

	addOrigin := func(rawURL string) {
		if origin := extractOriginFromURL(rawURL); origin != "" {
			if _, ok := seen[origin]; !ok {
				seen[origin] = struct{placeholder{placeholder
				origins = append(origins, origin)
		placeholder
	placeholder
placeholder

	// home content URL (when home_content is set to a URL for iframe embedding)
	addOrigin(settings.HomeContent)

	// purchase subscription URL
	if settings.PurchaseSubscriptionEnabled {
		addOrigin(settings.PurchaseSubscriptionURL)
placeholder

	// all custom menu items (including admin-only, since CSP must allow all iframes)
	for _, item := range parseCustomMenuItemURLs(settings.CustomMenuItems) {
		addOrigin(item)
placeholder

	return origins, nil
placeholder

// extractOriginFromURL returns the scheme+host origin from rawURL.
// Only http and https schemes are accepted.
func extractOriginFromURL(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
placeholder
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		return ""
placeholder
	if u.Scheme != "http" && u.Scheme != "https" {
		return ""
placeholder
	return u.Scheme + "://" + u.Host
placeholder

// parseCustomMenuItemURLs extracts URLs from a raw JSON array of custom menu items.
func parseCustomMenuItemURLs(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "[]" {
		return nil
placeholder
	var items []struct {
		URL string `json:"url"`
placeholder
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return nil
placeholder
	urls := make([]string, 0, len(items))
	for _, item := range items {
		if item.URL != "" {
			urls = append(urls, item.URL)
	placeholder
placeholder
	return urls
placeholder

// UpdateSettings 更新系统设置
func (s *SettingService) UpdateSettings(ctx context.Context, settings *SystemSettings) error {
	if err := s.validateDefaultSubscriptionGroups(ctx, settings.DefaultSubscriptions); err != nil {
		return err
placeholder
	normalizedWhitelist, err := NormalizeRegistrationEmailSuffixWhitelist(settings.RegistrationEmailSuffixWhitelist)
	if err != nil {
		return infraerrors.BadRequest("INVALID_REGISTRATION_EMAIL_SUFFIX_WHITELIST", err.Error())
placeholder
	if normalizedWhitelist == nil {
		normalizedWhitelist = []string{placeholder
placeholder
	settings.RegistrationEmailSuffixWhitelist = normalizedWhitelist

	updates := make(map[string]string)

	// 注册设置
	updates[SettingKeyRegistrationEnabled] = strconv.FormatBool(settings.RegistrationEnabled)
	updates[SettingKeyEmailVerifyEnabled] = strconv.FormatBool(settings.EmailVerifyEnabled)
	registrationEmailSuffixWhitelistJSON, err := json.Marshal(settings.RegistrationEmailSuffixWhitelist)
	if err != nil {
		return fmt.Errorf("marshal registration email suffix whitelist: %w", err)
placeholder
	updates[SettingKeyRegistrationEmailSuffixWhitelist] = string(registrationEmailSuffixWhitelistJSON)
	updates[SettingKeyPromoCodeEnabled] = strconv.FormatBool(settings.PromoCodeEnabled)
	updates[SettingKeyPasswordResetEnabled] = strconv.FormatBool(settings.PasswordResetEnabled)
	updates[SettingKeyFrontendURL] = settings.FrontendURL
	updates[SettingKeyInvitationCodeEnabled] = strconv.FormatBool(settings.InvitationCodeEnabled)
	updates[SettingKeyTotpEnabled] = strconv.FormatBool(settings.TotpEnabled)

	// 邮件服务设置（只有非空才更新密码）
	updates[SettingKeySMTPHost] = settings.SMTPHost
	updates[SettingKeySMTPPort] = strconv.Itoa(settings.SMTPPort)
	updates[SettingKeySMTPUsername] = settings.SMTPUsername
	if settings.SMTPPassword != "" {
		updates[SettingKeySMTPPassword] = settings.SMTPPassword
placeholder
	updates[SettingKeySMTPFrom] = settings.SMTPFrom
	updates[SettingKeySMTPFromName] = settings.SMTPFromName
	updates[SettingKeySMTPUseTLS] = strconv.FormatBool(settings.SMTPUseTLS)

	// Cloudflare Turnstile 设置（只有非空才更新密钥）
	updates[SettingKeyTurnstileEnabled] = strconv.FormatBool(settings.TurnstileEnabled)
	updates[SettingKeyTurnstileSiteKey] = settings.TurnstileSiteKey
	if settings.TurnstileSecretKey != "" {
		updates[SettingKeyTurnstileSecretKey] = settings.TurnstileSecretKey
placeholder

	// LinuxDo Connect OAuth 登录
	updates[SettingKeyLinuxDoConnectEnabled] = strconv.FormatBool(settings.LinuxDoConnectEnabled)
	updates[SettingKeyLinuxDoConnectClientID] = settings.LinuxDoConnectClientID
	updates[SettingKeyLinuxDoConnectRedirectURL] = settings.LinuxDoConnectRedirectURL
	if settings.LinuxDoConnectClientSecret != "" {
		updates[SettingKeyLinuxDoConnectClientSecret] = settings.LinuxDoConnectClientSecret
placeholder

	// Generic OIDC OAuth 登录
	updates[SettingKeyOIDCConnectEnabled] = strconv.FormatBool(settings.OIDCConnectEnabled)
	updates[SettingKeyOIDCConnectProviderName] = settings.OIDCConnectProviderName
	updates[SettingKeyOIDCConnectClientID] = settings.OIDCConnectClientID
	updates[SettingKeyOIDCConnectIssuerURL] = settings.OIDCConnectIssuerURL
	updates[SettingKeyOIDCConnectDiscoveryURL] = settings.OIDCConnectDiscoveryURL
	updates[SettingKeyOIDCConnectAuthorizeURL] = settings.OIDCConnectAuthorizeURL
	updates[SettingKeyOIDCConnectTokenURL] = settings.OIDCConnectTokenURL
	updates[SettingKeyOIDCConnectUserInfoURL] = settings.OIDCConnectUserInfoURL
	updates[SettingKeyOIDCConnectJWKSURL] = settings.OIDCConnectJWKSURL
	updates[SettingKeyOIDCConnectScopes] = settings.OIDCConnectScopes
	updates[SettingKeyOIDCConnectRedirectURL] = settings.OIDCConnectRedirectURL
	updates[SettingKeyOIDCConnectFrontendRedirectURL] = settings.OIDCConnectFrontendRedirectURL
	updates[SettingKeyOIDCConnectTokenAuthMethod] = settings.OIDCConnectTokenAuthMethod
	updates[SettingKeyOIDCConnectUsePKCE] = strconv.FormatBool(settings.OIDCConnectUsePKCE)
	updates[SettingKeyOIDCConnectValidateIDToken] = strconv.FormatBool(settings.OIDCConnectValidateIDToken)
	updates[SettingKeyOIDCConnectAllowedSigningAlgs] = settings.OIDCConnectAllowedSigningAlgs
	updates[SettingKeyOIDCConnectClockSkewSeconds] = strconv.Itoa(settings.OIDCConnectClockSkewSeconds)
	updates[SettingKeyOIDCConnectRequireEmailVerified] = strconv.FormatBool(settings.OIDCConnectRequireEmailVerified)
	updates[SettingKeyOIDCConnectUserInfoEmailPath] = settings.OIDCConnectUserInfoEmailPath
	updates[SettingKeyOIDCConnectUserInfoIDPath] = settings.OIDCConnectUserInfoIDPath
	updates[SettingKeyOIDCConnectUserInfoUsernamePath] = settings.OIDCConnectUserInfoUsernamePath
	if settings.OIDCConnectClientSecret != "" {
		updates[SettingKeyOIDCConnectClientSecret] = settings.OIDCConnectClientSecret
placeholder

	// OEM设置
	updates[SettingKeySiteName] = settings.SiteName
	updates[SettingKeySiteLogo] = settings.SiteLogo
	updates[SettingKeySiteSubtitle] = settings.SiteSubtitle
	updates[SettingKeyAPIBaseURL] = settings.APIBaseURL
	updates[SettingKeyContactInfo] = settings.ContactInfo
	updates[SettingKeyDocURL] = settings.DocURL
	updates[SettingKeyHomeContent] = settings.HomeContent
	updates[SettingKeyHideCcsImportButton] = strconv.FormatBool(settings.HideCcsImportButton)
	updates[SettingKeyPurchaseSubscriptionEnabled] = strconv.FormatBool(settings.PurchaseSubscriptionEnabled)
	updates[SettingKeyPurchaseSubscriptionURL] = strings.TrimSpace(settings.PurchaseSubscriptionURL)
	tableDefaultPageSize, tablePageSizeOptions := normalizeTablePreferences(
		settings.TableDefaultPageSize,
		settings.TablePageSizeOptions,
	)
	updates[SettingKeyTableDefaultPageSize] = strconv.Itoa(tableDefaultPageSize)
	tablePageSizeOptionsJSON, err := json.Marshal(tablePageSizeOptions)
	if err != nil {
		return fmt.Errorf("marshal table page size options: %w", err)
placeholder
	updates[SettingKeyTablePageSizeOptions] = string(tablePageSizeOptionsJSON)
	updates[SettingKeyCustomMenuItems] = settings.CustomMenuItems
	updates[SettingKeyCustomEndpoints] = settings.CustomEndpoints

	// 默认配置
	updates[SettingKeyDefaultConcurrency] = strconv.Itoa(settings.DefaultConcurrency)
	updates[SettingKeyDefaultBalance] = strconv.FormatFloat(settings.DefaultBalance, 'f', 8, 64)
	defaultSubsJSON, err := json.Marshal(settings.DefaultSubscriptions)
	if err != nil {
		return fmt.Errorf("marshal default subscriptions: %w", err)
placeholder
	updates[SettingKeyDefaultSubscriptions] = string(defaultSubsJSON)

	// Model fallback configuration
	updates[SettingKeyEnableModelFallback] = strconv.FormatBool(settings.EnableModelFallback)
	updates[SettingKeyFallbackModelAnthropic] = settings.FallbackModelAnthropic
	updates[SettingKeyFallbackModelOpenAI] = settings.FallbackModelOpenAI
	updates[SettingKeyFallbackModelGemini] = settings.FallbackModelGemini
	updates[SettingKeyFallbackModelAntigravity] = settings.FallbackModelAntigravity

	// Identity patch configuration (Claude -> Gemini)
	updates[SettingKeyEnableIdentityPatch] = strconv.FormatBool(settings.EnableIdentityPatch)
	updates[SettingKeyIdentityPatchPrompt] = settings.IdentityPatchPrompt

	// Ops monitoring (vNext)
	updates[SettingKeyOpsMonitoringEnabled] = strconv.FormatBool(settings.OpsMonitoringEnabled)
	updates[SettingKeyOpsRealtimeMonitoringEnabled] = strconv.FormatBool(settings.OpsRealtimeMonitoringEnabled)
	updates[SettingKeyOpsQueryModeDefault] = string(ParseOpsQueryMode(settings.OpsQueryModeDefault))
	if settings.OpsMetricsIntervalSeconds > 0 {
		updates[SettingKeyOpsMetricsIntervalSeconds] = strconv.Itoa(settings.OpsMetricsIntervalSeconds)
placeholder

	// Claude Code version check
	updates[SettingKeyMinClaudeCodeVersion] = settings.MinClaudeCodeVersion
	updates[SettingKeyMaxClaudeCodeVersion] = settings.MaxClaudeCodeVersion

	// 分组隔离
	updates[SettingKeyAllowUngroupedKeyScheduling] = strconv.FormatBool(settings.AllowUngroupedKeyScheduling)

	// Backend Mode
	updates[SettingKeyBackendModeEnabled] = strconv.FormatBool(settings.BackendModeEnabled)

	// Gateway forwarding behavior
	updates[SettingKeyEnableFingerprintUnification] = strconv.FormatBool(settings.EnableFingerprintUnification)
	updates[SettingKeyEnableMetadataPassthrough] = strconv.FormatBool(settings.EnableMetadataPassthrough)
	updates[SettingKeyEnableCCHSigning] = strconv.FormatBool(settings.EnableCCHSigning)

	// Balance low notification
	updates[SettingKeyBalanceLowNotifyEnabled] = strconv.FormatBool(settings.BalanceLowNotifyEnabled)
	updates[SettingKeyBalanceLowNotifyThreshold] = strconv.FormatFloat(settings.BalanceLowNotifyThreshold, 'f', 8, 64)
	accountQuotaNotifyEmailsJSON, err := json.Marshal(settings.AccountQuotaNotifyEmails)
	if err != nil {
		return fmt.Errorf("marshal account quota notify emails: %w", err)
placeholder
	updates[SettingKeyAccountQuotaNotifyEmails] = string(accountQuotaNotifyEmailsJSON)

	err = s.settingRepo.SetMultiple(ctx, updates)
	if err == nil {
		// 先使 inflight singleflight 失效，再刷新缓存，缩小旧值覆盖新值的竞态窗口
		versionBoundsSF.Forget("version_bounds")
		versionBoundsCache.Store(&cachedVersionBounds{
			min:       settings.MinClaudeCodeVersion,
			max:       settings.MaxClaudeCodeVersion,
			expiresAt: time.Now().Add(versionBoundsCacheTTL).UnixNano(),
	placeholder)
		backendModeSF.Forget("backend_mode")
		backendModeCache.Store(&cachedBackendMode{
			value:     settings.BackendModeEnabled,
			expiresAt: time.Now().Add(backendModeCacheTTL).UnixNano(),
	placeholder)
		gatewayForwardingSF.Forget("gateway_forwarding")
		gatewayForwardingCache.Store(&cachedGatewayForwardingSettings{
			fingerprintUnification: settings.EnableFingerprintUnification,
			metadataPassthrough:    settings.EnableMetadataPassthrough,
			cchSigning:             settings.EnableCCHSigning,
			expiresAt:              time.Now().Add(gatewayForwardingCacheTTL).UnixNano(),
	placeholder)
		if s.onUpdate != nil {
			s.onUpdate() // Invalidate cache after settings update
	placeholder
placeholder
	return err
placeholder

func (s *SettingService) validateDefaultSubscriptionGroups(ctx context.Context, items []DefaultSubscriptionSetting) error {
	if len(items) == 0 {
		return nil
placeholder

	checked := make(map[int64]struct{placeholder, len(items))
	for _, item := range items {
		if item.GroupID <= 0 {
			continue
	placeholder
		if _, ok := checked[item.GroupID]; ok {
			return ErrDefaultSubGroupDuplicate.WithMetadata(map[string]string{
				"group_id": strconv.FormatInt(item.GroupID, 10),
		placeholder)
	placeholder
		checked[item.GroupID] = struct{placeholder{placeholder
		if s.defaultSubGroupReader == nil {
			continue
	placeholder

		group, err := s.defaultSubGroupReader.GetByID(ctx, item.GroupID)
		if err != nil {
			if errors.Is(err, ErrGroupNotFound) {
				return ErrDefaultSubGroupInvalid.WithMetadata(map[string]string{
					"group_id": strconv.FormatInt(item.GroupID, 10),
			placeholder)
		placeholder
			return fmt.Errorf("get default subscription group %d: %w", item.GroupID, err)
	placeholder
		if !group.IsSubscriptionType() {
			return ErrDefaultSubGroupInvalid.WithMetadata(map[string]string{
				"group_id": strconv.FormatInt(item.GroupID, 10),
		placeholder)
	placeholder
placeholder

	return nil
placeholder

// IsRegistrationEnabled 检查是否开放注册
func (s *SettingService) IsRegistrationEnabled(ctx context.Context) bool {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyRegistrationEnabled)
	if err != nil {
		// 安全默认：如果设置不存在或查询出错，默认关闭注册
		return false
placeholder
	return value == "true"
placeholder

// IsBackendModeEnabled checks if backend mode is enabled
// Uses in-process atomic.Value cache with 60s TTL, zero-lock hot path
func (s *SettingService) IsBackendModeEnabled(ctx context.Context) bool {
	if cached, ok := backendModeCache.Load().(*cachedBackendMode); ok && cached != nil {
		if time.Now().UnixNano() < cached.expiresAt {
			return cached.value
	placeholder
placeholder
	result, _, _ := backendModeSF.Do("backend_mode", func() (any, error) {
		if cached, ok := backendModeCache.Load().(*cachedBackendMode); ok && cached != nil {
			if time.Now().UnixNano() < cached.expiresAt {
				return cached.value, nil
		placeholder
	placeholder
		dbCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), backendModeDBTimeout)
		defer cancel()
		value, err := s.settingRepo.GetValue(dbCtx, SettingKeyBackendModeEnabled)
		if err != nil {
			if errors.Is(err, ErrSettingNotFound) {
				// Setting not yet created (fresh install) - default to disabled with full TTL
				backendModeCache.Store(&cachedBackendMode{
					value:     false,
					expiresAt: time.Now().Add(backendModeCacheTTL).UnixNano(),
			placeholder)
				return false, nil
		placeholder
			slog.Warn("failed to get backend_mode_enabled setting", "error", err)
			backendModeCache.Store(&cachedBackendMode{
				value:     false,
				expiresAt: time.Now().Add(backendModeErrorTTL).UnixNano(),
		placeholder)
			return false, nil
	placeholder
		enabled := value == "true"
		backendModeCache.Store(&cachedBackendMode{
			value:     enabled,
			expiresAt: time.Now().Add(backendModeCacheTTL).UnixNano(),
	placeholder)
		return enabled, nil
placeholder)
	if val, ok := result.(bool); ok {
		return val
placeholder
	return false
placeholder

// GetGatewayForwardingSettings returns cached gateway forwarding settings.
// Uses in-process atomic.Value cache with 60s TTL, zero-lock hot path.
// Returns (fingerprintUnification, metadataPassthrough, cchSigning).
func (s *SettingService) GetGatewayForwardingSettings(ctx context.Context) (fingerprintUnification, metadataPassthrough, cchSigning bool) {
	if cached, ok := gatewayForwardingCache.Load().(*cachedGatewayForwardingSettings); ok && cached != nil {
		if time.Now().UnixNano() < cached.expiresAt {
			return cached.fingerprintUnification, cached.metadataPassthrough, cached.cchSigning
	placeholder
placeholder
	type gwfResult struct {
		fp, mp, cch bool
placeholder
	val, _, _ := gatewayForwardingSF.Do("gateway_forwarding", func() (any, error) {
		if cached, ok := gatewayForwardingCache.Load().(*cachedGatewayForwardingSettings); ok && cached != nil {
			if time.Now().UnixNano() < cached.expiresAt {
				return gwfResult{cached.fingerprintUnification, cached.metadataPassthrough, cached.cchSigningplaceholder, nil
		placeholder
	placeholder
		dbCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), gatewayForwardingDBTimeout)
		defer cancel()
		values, err := s.settingRepo.GetMultiple(dbCtx, []string{
			SettingKeyEnableFingerprintUnification,
			SettingKeyEnableMetadataPassthrough,
			SettingKeyEnableCCHSigning,
	placeholder)
		if err != nil {
			slog.Warn("failed to get gateway forwarding settings", "error", err)
			gatewayForwardingCache.Store(&cachedGatewayForwardingSettings{
				fingerprintUnification: true,
				metadataPassthrough:    false,
				cchSigning:             false,
				expiresAt:              time.Now().Add(gatewayForwardingErrorTTL).UnixNano(),
		placeholder)
			return gwfResult{true, false, falseplaceholder, nil
	placeholder
		fp := true
		if v, ok := values[SettingKeyEnableFingerprintUnification]; ok && v != "" {
			fp = v == "true"
	placeholder
		mp := values[SettingKeyEnableMetadataPassthrough] == "true"
		cch := values[SettingKeyEnableCCHSigning] == "true"
		gatewayForwardingCache.Store(&cachedGatewayForwardingSettings{
			fingerprintUnification: fp,
			metadataPassthrough:    mp,
			cchSigning:             cch,
			expiresAt:              time.Now().Add(gatewayForwardingCacheTTL).UnixNano(),
	placeholder)
		return gwfResult{fp, mp, cchplaceholder, nil
placeholder)
	if r, ok := val.(gwfResult); ok {
		return r.fp, r.mp, r.cch
placeholder
	return true, false, false // fail-open defaults
placeholder

// IsEmailVerifyEnabled 检查是否开启邮件验证
func (s *SettingService) IsEmailVerifyEnabled(ctx context.Context) bool {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyEmailVerifyEnabled)
	if err != nil {
		return false
placeholder
	return value == "true"
placeholder

// GetRegistrationEmailSuffixWhitelist returns normalized registration email suffix whitelist.
func (s *SettingService) GetRegistrationEmailSuffixWhitelist(ctx context.Context) []string {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyRegistrationEmailSuffixWhitelist)
	if err != nil {
		return []string{placeholder
placeholder
	return ParseRegistrationEmailSuffixWhitelist(value)
placeholder

// IsPromoCodeEnabled 检查是否启用优惠码功能
func (s *SettingService) IsPromoCodeEnabled(ctx context.Context) bool {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyPromoCodeEnabled)
	if err != nil {
		return true // 默认启用
placeholder
	return value != "false"
placeholder

// IsInvitationCodeEnabled 检查是否启用邀请码注册功能
func (s *SettingService) IsInvitationCodeEnabled(ctx context.Context) bool {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyInvitationCodeEnabled)
	if err != nil {
		return false // 默认关闭
placeholder
	return value == "true"
placeholder

// IsPasswordResetEnabled 检查是否启用密码重置功能
// 要求：必须同时开启邮件验证
func (s *SettingService) IsPasswordResetEnabled(ctx context.Context) bool {
	// Password reset requires email verification to be enabled
	if !s.IsEmailVerifyEnabled(ctx) {
		return false
placeholder
	value, err := s.settingRepo.GetValue(ctx, SettingKeyPasswordResetEnabled)
	if err != nil {
		return false // 默认关闭
placeholder
	return value == "true"
placeholder

// IsTotpEnabled 检查是否启用 TOTP 双因素认证功能
func (s *SettingService) IsTotpEnabled(ctx context.Context) bool {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyTotpEnabled)
	if err != nil {
		return false // 默认关闭
placeholder
	return value == "true"
placeholder

// IsTotpEncryptionKeyConfigured 检查 TOTP 加密密钥是否已手动配置
// 只有手动配置了密钥才允许在管理后台启用 TOTP 功能
func (s *SettingService) IsTotpEncryptionKeyConfigured() bool {
	return s.cfg.Totp.EncryptionKeyConfigured
placeholder

// GetSiteName 获取网站名称
func (s *SettingService) GetSiteName(ctx context.Context) string {
	value, err := s.settingRepo.GetValue(ctx, SettingKeySiteName)
	if err != nil || value == "" {
		return "Sub2API"
placeholder
	return value
placeholder

// GetDefaultConcurrency 获取默认并发量
func (s *SettingService) GetDefaultConcurrency(ctx context.Context) int {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyDefaultConcurrency)
	if err != nil {
		return s.cfg.Default.UserConcurrency
placeholder
	if v, err := strconv.Atoi(value); err == nil && v > 0 {
		return v
placeholder
	return s.cfg.Default.UserConcurrency
placeholder

// GetDefaultBalance 获取默认余额
func (s *SettingService) GetDefaultBalance(ctx context.Context) float64 {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyDefaultBalance)
	if err != nil {
		return s.cfg.Default.UserBalance
placeholder
	if v, err := strconv.ParseFloat(value, 64); err == nil && v >= 0 {
		return v
placeholder
	return s.cfg.Default.UserBalance
placeholder

// GetDefaultSubscriptions 获取新用户默认订阅配置列表。
func (s *SettingService) GetDefaultSubscriptions(ctx context.Context) []DefaultSubscriptionSetting {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyDefaultSubscriptions)
	if err != nil {
		return nil
placeholder
	return parseDefaultSubscriptions(value)
placeholder

// InitializeDefaultSettings 初始化默认设置
func (s *SettingService) InitializeDefaultSettings(ctx context.Context) error {
	// 检查是否已有设置
	_, err := s.settingRepo.GetValue(ctx, SettingKeyRegistrationEnabled)
	if err == nil {
		// 已有设置，不需要初始化
		return nil
placeholder
	if !errors.Is(err, ErrSettingNotFound) {
		return fmt.Errorf("check existing settings: %w", err)
placeholder

	// 初始化默认设置
	defaults := map[string]string{
		SettingKeyRegistrationEnabled:              "true",
		SettingKeyEmailVerifyEnabled:               "false",
		SettingKeyRegistrationEmailSuffixWhitelist: "[]",
		SettingKeyPromoCodeEnabled:                 "true", // 默认启用优惠码功能
		SettingKeySiteName:                         "Sub2API",
		SettingKeySiteLogo:                         "",
		SettingKeyPurchaseSubscriptionEnabled:      "false",
		SettingKeyPurchaseSubscriptionURL:          "",
		SettingKeyTableDefaultPageSize:             "20",
		SettingKeyTablePageSizeOptions:             "[10,20,50,100]",
		SettingKeyCustomMenuItems:                  "[]",
		SettingKeyCustomEndpoints:                  "[]",
		SettingKeyOIDCConnectEnabled:               "false",
		SettingKeyOIDCConnectProviderName:          "OIDC",
		SettingKeyDefaultConcurrency:               strconv.Itoa(s.cfg.Default.UserConcurrency),
		SettingKeyDefaultBalance:                   strconv.FormatFloat(s.cfg.Default.UserBalance, 'f', 8, 64),
		SettingKeyDefaultSubscriptions:             "[]",
		SettingKeySMTPPort:                         "587",
		SettingKeySMTPUseTLS:                       "false",
		// Model fallback defaults
		SettingKeyEnableModelFallback:      "false",
		SettingKeyFallbackModelAnthropic:   "claude-3-5-sonnet-20241022",
		SettingKeyFallbackModelOpenAI:      "gpt-4o",
		SettingKeyFallbackModelGemini:      "gemini-2.5-pro",
		SettingKeyFallbackModelAntigravity: "gemini-2.5-pro",
		// Identity patch defaults
		SettingKeyEnableIdentityPatch: "true",
		SettingKeyIdentityPatchPrompt: "",

		// Ops monitoring defaults (vNext)
		SettingKeyOpsMonitoringEnabled:         "true",
		SettingKeyOpsRealtimeMonitoringEnabled: "true",
		SettingKeyOpsQueryModeDefault:          "auto",
		SettingKeyOpsMetricsIntervalSeconds:    "60",

		// Claude Code version check (default: empty = disabled)
		SettingKeyMinClaudeCodeVersion: "",
		SettingKeyMaxClaudeCodeVersion: "",

		// 分组隔离（默认不允许未分组 Key 调度）
		SettingKeyAllowUngroupedKeyScheduling: "false",
placeholder

	return s.settingRepo.SetMultiple(ctx, defaults)
placeholder

// parseSettings 解析设置到结构体
func (s *SettingService) parseSettings(settings map[string]string) *SystemSettings {
	emailVerifyEnabled := settings[SettingKeyEmailVerifyEnabled] == "true"
	result := &SystemSettings{
		RegistrationEnabled:              settings[SettingKeyRegistrationEnabled] == "true",
		EmailVerifyEnabled:               emailVerifyEnabled,
		RegistrationEmailSuffixWhitelist: ParseRegistrationEmailSuffixWhitelist(settings[SettingKeyRegistrationEmailSuffixWhitelist]),
		PromoCodeEnabled:                 settings[SettingKeyPromoCodeEnabled] != "false", // 默认启用
		PasswordResetEnabled:             emailVerifyEnabled && settings[SettingKeyPasswordResetEnabled] == "true",
		FrontendURL:                      settings[SettingKeyFrontendURL],
		InvitationCodeEnabled:            settings[SettingKeyInvitationCodeEnabled] == "true",
		TotpEnabled:                      settings[SettingKeyTotpEnabled] == "true",
		SMTPHost:                         settings[SettingKeySMTPHost],
		SMTPUsername:                     settings[SettingKeySMTPUsername],
		SMTPFrom:                         settings[SettingKeySMTPFrom],
		SMTPFromName:                     settings[SettingKeySMTPFromName],
		SMTPUseTLS:                       settings[SettingKeySMTPUseTLS] == "true",
		SMTPPasswordConfigured:           settings[SettingKeySMTPPassword] != "",
		TurnstileEnabled:                 settings[SettingKeyTurnstileEnabled] == "true",
		TurnstileSiteKey:                 settings[SettingKeyTurnstileSiteKey],
		TurnstileSecretKeyConfigured:     settings[SettingKeyTurnstileSecretKey] != "",
		SiteName:                         s.getStringOrDefault(settings, SettingKeySiteName, "Sub2API"),
		SiteLogo:                         settings[SettingKeySiteLogo],
		SiteSubtitle:                     s.getStringOrDefault(settings, SettingKeySiteSubtitle, "Subscription to API Conversion Platform"),
		APIBaseURL:                       settings[SettingKeyAPIBaseURL],
		ContactInfo:                      settings[SettingKeyContactInfo],
		DocURL:                           settings[SettingKeyDocURL],
		HomeContent:                      settings[SettingKeyHomeContent],
		HideCcsImportButton:              settings[SettingKeyHideCcsImportButton] == "true",
		PurchaseSubscriptionEnabled:      settings[SettingKeyPurchaseSubscriptionEnabled] == "true",
		PurchaseSubscriptionURL:          strings.TrimSpace(settings[SettingKeyPurchaseSubscriptionURL]),
		CustomMenuItems:                  settings[SettingKeyCustomMenuItems],
		CustomEndpoints:                  settings[SettingKeyCustomEndpoints],
		BackendModeEnabled:               settings[SettingKeyBackendModeEnabled] == "true",
placeholder
	result.TableDefaultPageSize, result.TablePageSizeOptions = parseTablePreferences(
		settings[SettingKeyTableDefaultPageSize],
		settings[SettingKeyTablePageSizeOptions],
	)

	// 解析整数类型
	if port, err := strconv.Atoi(settings[SettingKeySMTPPort]); err == nil {
		result.SMTPPort = port
placeholder else {
		result.SMTPPort = 587
placeholder

	if concurrency, err := strconv.Atoi(settings[SettingKeyDefaultConcurrency]); err == nil {
		result.DefaultConcurrency = concurrency
placeholder else {
		result.DefaultConcurrency = s.cfg.Default.UserConcurrency
placeholder

	// 解析浮点数类型
	if balance, err := strconv.ParseFloat(settings[SettingKeyDefaultBalance], 64); err == nil {
		result.DefaultBalance = balance
placeholder else {
		result.DefaultBalance = s.cfg.Default.UserBalance
placeholder
	result.DefaultSubscriptions = parseDefaultSubscriptions(settings[SettingKeyDefaultSubscriptions])

	// 敏感信息直接返回，方便测试连接时使用
	result.SMTPPassword = settings[SettingKeySMTPPassword]
	result.TurnstileSecretKey = settings[SettingKeyTurnstileSecretKey]

	// LinuxDo Connect 设置：
	// - 兼容 config.yaml/env（避免老部署因为未迁移到数据库设置而被意外关闭）
	// - 支持在后台“系统设置”中覆盖并持久化（存储于 DB）
	linuxDoBase := config.LinuxDoConnectConfig{placeholder
	if s.cfg != nil {
		linuxDoBase = s.cfg.LinuxDo
placeholder

	if raw, ok := settings[SettingKeyLinuxDoConnectEnabled]; ok {
		result.LinuxDoConnectEnabled = raw == "true"
placeholder else {
		result.LinuxDoConnectEnabled = linuxDoBase.Enabled
placeholder

	if v, ok := settings[SettingKeyLinuxDoConnectClientID]; ok && strings.TrimSpace(v) != "" {
		result.LinuxDoConnectClientID = strings.TrimSpace(v)
placeholder else {
		result.LinuxDoConnectClientID = linuxDoBase.ClientID
placeholder

	if v, ok := settings[SettingKeyLinuxDoConnectRedirectURL]; ok && strings.TrimSpace(v) != "" {
		result.LinuxDoConnectRedirectURL = strings.TrimSpace(v)
placeholder else {
		result.LinuxDoConnectRedirectURL = linuxDoBase.RedirectURL
placeholder

	result.LinuxDoConnectClientSecret = strings.TrimSpace(settings[SettingKeyLinuxDoConnectClientSecret])
	if result.LinuxDoConnectClientSecret == "" {
		result.LinuxDoConnectClientSecret = strings.TrimSpace(linuxDoBase.ClientSecret)
placeholder
	result.LinuxDoConnectClientSecretConfigured = result.LinuxDoConnectClientSecret != ""

	// Generic OIDC 设置：
	// - 兼容 config.yaml/env
	// - 支持后台系统设置覆盖并持久化（存储于 DB）
	oidcBase := config.OIDCConnectConfig{placeholder
	if s.cfg != nil {
		oidcBase = s.cfg.OIDC
placeholder

	if raw, ok := settings[SettingKeyOIDCConnectEnabled]; ok {
		result.OIDCConnectEnabled = raw == "true"
placeholder else {
		result.OIDCConnectEnabled = oidcBase.Enabled
placeholder

	if v, ok := settings[SettingKeyOIDCConnectProviderName]; ok && strings.TrimSpace(v) != "" {
		result.OIDCConnectProviderName = strings.TrimSpace(v)
placeholder else {
		result.OIDCConnectProviderName = strings.TrimSpace(oidcBase.ProviderName)
placeholder
	if result.OIDCConnectProviderName == "" {
		result.OIDCConnectProviderName = "OIDC"
placeholder

	if v, ok := settings[SettingKeyOIDCConnectClientID]; ok && strings.TrimSpace(v) != "" {
		result.OIDCConnectClientID = strings.TrimSpace(v)
placeholder else {
		result.OIDCConnectClientID = strings.TrimSpace(oidcBase.ClientID)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectIssuerURL]; ok && strings.TrimSpace(v) != "" {
		result.OIDCConnectIssuerURL = strings.TrimSpace(v)
placeholder else {
		result.OIDCConnectIssuerURL = strings.TrimSpace(oidcBase.IssuerURL)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectDiscoveryURL]; ok && strings.TrimSpace(v) != "" {
		result.OIDCConnectDiscoveryURL = strings.TrimSpace(v)
placeholder else {
		result.OIDCConnectDiscoveryURL = strings.TrimSpace(oidcBase.DiscoveryURL)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectAuthorizeURL]; ok && strings.TrimSpace(v) != "" {
		result.OIDCConnectAuthorizeURL = strings.TrimSpace(v)
placeholder else {
		result.OIDCConnectAuthorizeURL = strings.TrimSpace(oidcBase.AuthorizeURL)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectTokenURL]; ok && strings.TrimSpace(v) != "" {
		result.OIDCConnectTokenURL = strings.TrimSpace(v)
placeholder else {
		result.OIDCConnectTokenURL = strings.TrimSpace(oidcBase.TokenURL)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectUserInfoURL]; ok && strings.TrimSpace(v) != "" {
		result.OIDCConnectUserInfoURL = strings.TrimSpace(v)
placeholder else {
		result.OIDCConnectUserInfoURL = strings.TrimSpace(oidcBase.UserInfoURL)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectJWKSURL]; ok && strings.TrimSpace(v) != "" {
		result.OIDCConnectJWKSURL = strings.TrimSpace(v)
placeholder else {
		result.OIDCConnectJWKSURL = strings.TrimSpace(oidcBase.JWKSURL)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectScopes]; ok && strings.TrimSpace(v) != "" {
		result.OIDCConnectScopes = strings.TrimSpace(v)
placeholder else {
		result.OIDCConnectScopes = strings.TrimSpace(oidcBase.Scopes)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectRedirectURL]; ok && strings.TrimSpace(v) != "" {
		result.OIDCConnectRedirectURL = strings.TrimSpace(v)
placeholder else {
		result.OIDCConnectRedirectURL = strings.TrimSpace(oidcBase.RedirectURL)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectFrontendRedirectURL]; ok && strings.TrimSpace(v) != "" {
		result.OIDCConnectFrontendRedirectURL = strings.TrimSpace(v)
placeholder else {
		result.OIDCConnectFrontendRedirectURL = strings.TrimSpace(oidcBase.FrontendRedirectURL)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectTokenAuthMethod]; ok && strings.TrimSpace(v) != "" {
		result.OIDCConnectTokenAuthMethod = strings.ToLower(strings.TrimSpace(v))
placeholder else {
		result.OIDCConnectTokenAuthMethod = strings.ToLower(strings.TrimSpace(oidcBase.TokenAuthMethod))
placeholder
	if raw, ok := settings[SettingKeyOIDCConnectUsePKCE]; ok {
		result.OIDCConnectUsePKCE = raw == "true"
placeholder else {
		result.OIDCConnectUsePKCE = oidcBase.UsePKCE
placeholder
	if raw, ok := settings[SettingKeyOIDCConnectValidateIDToken]; ok {
		result.OIDCConnectValidateIDToken = raw == "true"
placeholder else {
		result.OIDCConnectValidateIDToken = oidcBase.ValidateIDToken
placeholder
	if v, ok := settings[SettingKeyOIDCConnectAllowedSigningAlgs]; ok && strings.TrimSpace(v) != "" {
		result.OIDCConnectAllowedSigningAlgs = strings.TrimSpace(v)
placeholder else {
		result.OIDCConnectAllowedSigningAlgs = strings.TrimSpace(oidcBase.AllowedSigningAlgs)
placeholder
	clockSkewSet := false
	if raw, ok := settings[SettingKeyOIDCConnectClockSkewSeconds]; ok && strings.TrimSpace(raw) != "" {
		if parsed, err := strconv.Atoi(strings.TrimSpace(raw)); err == nil {
			result.OIDCConnectClockSkewSeconds = parsed
			clockSkewSet = true
	placeholder
placeholder
	if !clockSkewSet {
		result.OIDCConnectClockSkewSeconds = oidcBase.ClockSkewSeconds
placeholder
	if !clockSkewSet && result.OIDCConnectClockSkewSeconds == 0 {
		result.OIDCConnectClockSkewSeconds = 120
placeholder
	if raw, ok := settings[SettingKeyOIDCConnectRequireEmailVerified]; ok {
		result.OIDCConnectRequireEmailVerified = raw == "true"
placeholder else {
		result.OIDCConnectRequireEmailVerified = oidcBase.RequireEmailVerified
placeholder
	if v, ok := settings[SettingKeyOIDCConnectUserInfoEmailPath]; ok {
		result.OIDCConnectUserInfoEmailPath = strings.TrimSpace(v)
placeholder else {
		result.OIDCConnectUserInfoEmailPath = strings.TrimSpace(oidcBase.UserInfoEmailPath)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectUserInfoIDPath]; ok {
		result.OIDCConnectUserInfoIDPath = strings.TrimSpace(v)
placeholder else {
		result.OIDCConnectUserInfoIDPath = strings.TrimSpace(oidcBase.UserInfoIDPath)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectUserInfoUsernamePath]; ok {
		result.OIDCConnectUserInfoUsernamePath = strings.TrimSpace(v)
placeholder else {
		result.OIDCConnectUserInfoUsernamePath = strings.TrimSpace(oidcBase.UserInfoUsernamePath)
placeholder
	result.OIDCConnectClientSecret = strings.TrimSpace(settings[SettingKeyOIDCConnectClientSecret])
	if result.OIDCConnectClientSecret == "" {
		result.OIDCConnectClientSecret = strings.TrimSpace(oidcBase.ClientSecret)
placeholder
	result.OIDCConnectClientSecretConfigured = result.OIDCConnectClientSecret != ""

	// Model fallback settings
	result.EnableModelFallback = settings[SettingKeyEnableModelFallback] == "true"
	result.FallbackModelAnthropic = s.getStringOrDefault(settings, SettingKeyFallbackModelAnthropic, "claude-3-5-sonnet-20241022")
	result.FallbackModelOpenAI = s.getStringOrDefault(settings, SettingKeyFallbackModelOpenAI, "gpt-4o")
	result.FallbackModelGemini = s.getStringOrDefault(settings, SettingKeyFallbackModelGemini, "gemini-2.5-pro")
	result.FallbackModelAntigravity = s.getStringOrDefault(settings, SettingKeyFallbackModelAntigravity, "gemini-2.5-pro")

	// Identity patch settings (default: enabled, to preserve existing behavior)
	if v, ok := settings[SettingKeyEnableIdentityPatch]; ok && v != "" {
		result.EnableIdentityPatch = v == "true"
placeholder else {
		result.EnableIdentityPatch = true
placeholder
	result.IdentityPatchPrompt = settings[SettingKeyIdentityPatchPrompt]

	// Ops monitoring settings (default: enabled, fail-open)
	result.OpsMonitoringEnabled = !isFalseSettingValue(settings[SettingKeyOpsMonitoringEnabled])
	result.OpsRealtimeMonitoringEnabled = !isFalseSettingValue(settings[SettingKeyOpsRealtimeMonitoringEnabled])
	result.OpsQueryModeDefault = string(ParseOpsQueryMode(settings[SettingKeyOpsQueryModeDefault]))
	result.OpsMetricsIntervalSeconds = 60
	if raw := strings.TrimSpace(settings[SettingKeyOpsMetricsIntervalSeconds]); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil {
			if v < 60 {
				v = 60
		placeholder
			if v > 3600 {
				v = 3600
		placeholder
			result.OpsMetricsIntervalSeconds = v
	placeholder
placeholder

	// Claude Code version check
	result.MinClaudeCodeVersion = settings[SettingKeyMinClaudeCodeVersion]
	result.MaxClaudeCodeVersion = settings[SettingKeyMaxClaudeCodeVersion]

	// 分组隔离
	result.AllowUngroupedKeyScheduling = settings[SettingKeyAllowUngroupedKeyScheduling] == "true"

	// Gateway forwarding behavior (defaults: fingerprint=true, metadata_passthrough=false, cch_signing=false)
	if v, ok := settings[SettingKeyEnableFingerprintUnification]; ok && v != "" {
		result.EnableFingerprintUnification = v == "true"
placeholder else {
		result.EnableFingerprintUnification = true // default: enabled (current behavior)
placeholder
	result.EnableMetadataPassthrough = settings[SettingKeyEnableMetadataPassthrough] == "true"
	result.EnableCCHSigning = settings[SettingKeyEnableCCHSigning] == "true"

	// Balance low notification
	result.BalanceLowNotifyEnabled = settings[SettingKeyBalanceLowNotifyEnabled] == "true"
	if v, err := strconv.ParseFloat(settings[SettingKeyBalanceLowNotifyThreshold], 64); err == nil && v >= 0 {
		result.BalanceLowNotifyThreshold = v
placeholder

	// Account quota notification emails
	if raw := strings.TrimSpace(settings[SettingKeyAccountQuotaNotifyEmails]); raw != "" {
		var emails []string
		if err := json.Unmarshal([]byte(raw), &emails); err == nil {
			result.AccountQuotaNotifyEmails = emails
	placeholder
placeholder
	if result.AccountQuotaNotifyEmails == nil {
		result.AccountQuotaNotifyEmails = []string{placeholder
placeholder

	return result
placeholder

func isFalseSettingValue(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "false", "0", "off", "disabled":
		return true
	default:
		return false
placeholder
placeholder

func parseDefaultSubscriptions(raw string) []DefaultSubscriptionSetting {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
placeholder

	var items []DefaultSubscriptionSetting
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return nil
placeholder

	normalized := make([]DefaultSubscriptionSetting, 0, len(items))
	for _, item := range items {
		if item.GroupID <= 0 || item.ValidityDays <= 0 {
			continue
	placeholder
		if item.ValidityDays > MaxValidityDays {
			item.ValidityDays = MaxValidityDays
	placeholder
		normalized = append(normalized, item)
placeholder

	return normalized
placeholder

func parseTablePreferences(defaultPageSizeRaw, optionsRaw string) (int, []int) {
	defaultPageSize := 20
	if v, err := strconv.Atoi(strings.TrimSpace(defaultPageSizeRaw)); err == nil {
		defaultPageSize = v
placeholder

	var options []int
	if strings.TrimSpace(optionsRaw) != "" {
		_ = json.Unmarshal([]byte(optionsRaw), &options)
placeholder

	return normalizeTablePreferences(defaultPageSize, options)
placeholder

func normalizeTablePreferences(defaultPageSize int, options []int) (int, []int) {
	const minPageSize = 5
	const maxPageSize = 1000
	const fallbackPageSize = 20

	seen := make(map[int]struct{placeholder, len(options))
	normalizedOptions := make([]int, 0, len(options))
	for _, option := range options {
		if option < minPageSize || option > maxPageSize {
			continue
	placeholder
		if _, ok := seen[option]; ok {
			continue
	placeholder
		seen[option] = struct{placeholder{placeholder
		normalizedOptions = append(normalizedOptions, option)
placeholder
	sort.Ints(normalizedOptions)

	if defaultPageSize < minPageSize || defaultPageSize > maxPageSize {
		defaultPageSize = fallbackPageSize
placeholder

	if len(normalizedOptions) == 0 {
		normalizedOptions = []int{10, 20, 50placeholder
placeholder

	return defaultPageSize, normalizedOptions
placeholder

// getStringOrDefault 获取字符串值或默认值
func (s *SettingService) getStringOrDefault(settings map[string]string, key, defaultValue string) string {
	if value, ok := settings[key]; ok && value != "" {
		return value
placeholder
	return defaultValue
placeholder

// IsTurnstileEnabled 检查是否启用 Turnstile 验证
func (s *SettingService) IsTurnstileEnabled(ctx context.Context) bool {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyTurnstileEnabled)
	if err != nil {
		return false
placeholder
	return value == "true"
placeholder

// GetTurnstileSecretKey 获取 Turnstile Secret Key
func (s *SettingService) GetTurnstileSecretKey(ctx context.Context) string {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyTurnstileSecretKey)
	if err != nil {
		return ""
placeholder
	return value
placeholder

// IsIdentityPatchEnabled 检查是否启用身份补丁（Claude -> Gemini systemInstruction 注入）
func (s *SettingService) IsIdentityPatchEnabled(ctx context.Context) bool {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyEnableIdentityPatch)
	if err != nil {
		// 默认开启，保持兼容
		return true
placeholder
	return value == "true"
placeholder

// GetIdentityPatchPrompt 获取自定义身份补丁提示词（为空表示使用内置默认模板）
func (s *SettingService) GetIdentityPatchPrompt(ctx context.Context) string {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyIdentityPatchPrompt)
	if err != nil {
		return ""
placeholder
	return value
placeholder

// GenerateAdminAPIKey 生成新的管理员 API Key
func (s *SettingService) GenerateAdminAPIKey(ctx context.Context) (string, error) {
	// 生成 32 字节随机数 = 64 位十六进制字符
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate random bytes: %w", err)
placeholder

	key := AdminAPIKeyPrefix + hex.EncodeToString(bytes)

	// 存储到 settings 表
	if err := s.settingRepo.Set(ctx, SettingKeyAdminAPIKey, key); err != nil {
		return "", fmt.Errorf("save admin api key: %w", err)
placeholder

	return key, nil
placeholder

// GetAdminAPIKeyStatus 获取管理员 API Key 状态
// 返回脱敏的 key、是否存在、错误
func (s *SettingService) GetAdminAPIKeyStatus(ctx context.Context) (maskedKey string, exists bool, err error) {
	key, err := s.settingRepo.GetValue(ctx, SettingKeyAdminAPIKey)
	if err != nil {
		if errors.Is(err, ErrSettingNotFound) {
			return "", false, nil
	placeholder
		return "", false, err
placeholder
	if key == "" {
		return "", false, nil
placeholder

	// 脱敏：显示前 10 位和后 4 位
	if len(key) > 14 {
		maskedKey = key[:10] + "..." + key[len(key)-4:]
placeholder else {
		maskedKey = key
placeholder

	return maskedKey, true, nil
placeholder

// GetAdminAPIKey 获取完整的管理员 API Key（仅供内部验证使用）
// 如果未配置返回空字符串和 nil 错误，只有数据库错误时才返回 error
func (s *SettingService) GetAdminAPIKey(ctx context.Context) (string, error) {
	key, err := s.settingRepo.GetValue(ctx, SettingKeyAdminAPIKey)
	if err != nil {
		if errors.Is(err, ErrSettingNotFound) {
			return "", nil // 未配置，返回空字符串
	placeholder
		return "", err // 数据库错误
placeholder
	return key, nil
placeholder

// DeleteAdminAPIKey 删除管理员 API Key
func (s *SettingService) DeleteAdminAPIKey(ctx context.Context) error {
	return s.settingRepo.Delete(ctx, SettingKeyAdminAPIKey)
placeholder

// IsModelFallbackEnabled 检查是否启用模型兜底机制
func (s *SettingService) IsModelFallbackEnabled(ctx context.Context) bool {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyEnableModelFallback)
	if err != nil {
		return false // Default: disabled
placeholder
	return value == "true"
placeholder

// GetFallbackModel 获取指定平台的兜底模型
func (s *SettingService) GetFallbackModel(ctx context.Context, platform string) string {
	var key string
	var defaultModel string

	switch platform {
	case PlatformAnthropic:
		key = SettingKeyFallbackModelAnthropic
		defaultModel = "claude-3-5-sonnet-20241022"
	case PlatformOpenAI:
		key = SettingKeyFallbackModelOpenAI
		defaultModel = "gpt-4o"
	case PlatformGemini:
		key = SettingKeyFallbackModelGemini
		defaultModel = "gemini-2.5-pro"
	case PlatformAntigravity:
		key = SettingKeyFallbackModelAntigravity
		defaultModel = "gemini-2.5-pro"
	default:
		return ""
placeholder

	value, err := s.settingRepo.GetValue(ctx, key)
	if err != nil || value == "" {
		return defaultModel
placeholder
	return value
placeholder

// GetLinuxDoConnectOAuthConfig 返回用于登录的"最终生效" LinuxDo Connect 配置。
//
// 优先级：
// - 若对应系统设置键存在，则覆盖 config.yaml/env 的值
// - 否则回退到 config.yaml/env 的值
func (s *SettingService) GetLinuxDoConnectOAuthConfig(ctx context.Context) (config.LinuxDoConnectConfig, error) {
	if s == nil || s.cfg == nil {
		return config.LinuxDoConnectConfig{placeholder, infraerrors.ServiceUnavailable("CONFIG_NOT_READY", "config not loaded")
placeholder

	effective := s.cfg.LinuxDo

	keys := []string{
		SettingKeyLinuxDoConnectEnabled,
		SettingKeyLinuxDoConnectClientID,
		SettingKeyLinuxDoConnectClientSecret,
		SettingKeyLinuxDoConnectRedirectURL,
placeholder
	settings, err := s.settingRepo.GetMultiple(ctx, keys)
	if err != nil {
		return config.LinuxDoConnectConfig{placeholder, fmt.Errorf("get linuxdo connect settings: %w", err)
placeholder

	if raw, ok := settings[SettingKeyLinuxDoConnectEnabled]; ok {
		effective.Enabled = raw == "true"
placeholder
	if v, ok := settings[SettingKeyLinuxDoConnectClientID]; ok && strings.TrimSpace(v) != "" {
		effective.ClientID = strings.TrimSpace(v)
placeholder
	if v, ok := settings[SettingKeyLinuxDoConnectClientSecret]; ok && strings.TrimSpace(v) != "" {
		effective.ClientSecret = strings.TrimSpace(v)
placeholder
	if v, ok := settings[SettingKeyLinuxDoConnectRedirectURL]; ok && strings.TrimSpace(v) != "" {
		effective.RedirectURL = strings.TrimSpace(v)
placeholder

	if !effective.Enabled {
		return config.LinuxDoConnectConfig{placeholder, infraerrors.NotFound("OAUTH_DISABLED", "oauth login is disabled")
placeholder

	// 基础健壮性校验（避免把用户重定向到一个必然失败或不安全的 OAuth 流程里）。
	if strings.TrimSpace(effective.ClientID) == "" {
		return config.LinuxDoConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth client id not configured")
placeholder
	if strings.TrimSpace(effective.AuthorizeURL) == "" {
		return config.LinuxDoConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth authorize url not configured")
placeholder
	if strings.TrimSpace(effective.TokenURL) == "" {
		return config.LinuxDoConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth token url not configured")
placeholder
	if strings.TrimSpace(effective.UserInfoURL) == "" {
		return config.LinuxDoConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth userinfo url not configured")
placeholder
	if strings.TrimSpace(effective.RedirectURL) == "" {
		return config.LinuxDoConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth redirect url not configured")
placeholder
	if strings.TrimSpace(effective.FrontendRedirectURL) == "" {
		return config.LinuxDoConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth frontend redirect url not configured")
placeholder

	if err := config.ValidateAbsoluteHTTPURL(effective.AuthorizeURL); err != nil {
		return config.LinuxDoConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth authorize url invalid")
placeholder
	if err := config.ValidateAbsoluteHTTPURL(effective.TokenURL); err != nil {
		return config.LinuxDoConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth token url invalid")
placeholder
	if err := config.ValidateAbsoluteHTTPURL(effective.UserInfoURL); err != nil {
		return config.LinuxDoConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth userinfo url invalid")
placeholder
	if err := config.ValidateAbsoluteHTTPURL(effective.RedirectURL); err != nil {
		return config.LinuxDoConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth redirect url invalid")
placeholder
	if err := config.ValidateFrontendRedirectURL(effective.FrontendRedirectURL); err != nil {
		return config.LinuxDoConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth frontend redirect url invalid")
placeholder

	method := strings.ToLower(strings.TrimSpace(effective.TokenAuthMethod))
	switch method {
	case "", "client_secret_post", "client_secret_basic":
		if strings.TrimSpace(effective.ClientSecret) == "" {
			return config.LinuxDoConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth client secret not configured")
	placeholder
	case "none":
		if !effective.UsePKCE {
			return config.LinuxDoConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth pkce must be enabled when token_auth_method=none")
	placeholder
	default:
		return config.LinuxDoConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth token_auth_method invalid")
placeholder

	return effective, nil
placeholder

// GetOverloadCooldownSettings 获取529过载冷却配置
func (s *SettingService) GetOverloadCooldownSettings(ctx context.Context) (*OverloadCooldownSettings, error) {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyOverloadCooldownSettings)
	if err != nil {
		if errors.Is(err, ErrSettingNotFound) {
			return DefaultOverloadCooldownSettings(), nil
	placeholder
		return nil, fmt.Errorf("get overload cooldown settings: %w", err)
placeholder
	if value == "" {
		return DefaultOverloadCooldownSettings(), nil
placeholder

	var settings OverloadCooldownSettings
	if err := json.Unmarshal([]byte(value), &settings); err != nil {
		return DefaultOverloadCooldownSettings(), nil
placeholder

	// 修正配置值范围
	if settings.CooldownMinutes < 1 {
		settings.CooldownMinutes = 1
placeholder
	if settings.CooldownMinutes > 120 {
		settings.CooldownMinutes = 120
placeholder

	return &settings, nil
placeholder

// SetOverloadCooldownSettings 设置529过载冷却配置
func (s *SettingService) SetOverloadCooldownSettings(ctx context.Context, settings *OverloadCooldownSettings) error {
	if settings == nil {
		return fmt.Errorf("settings cannot be nil")
placeholder

	// 禁用时修正为合法值即可，不拒绝请求
	if settings.CooldownMinutes < 1 || settings.CooldownMinutes > 120 {
		if settings.Enabled {
			return fmt.Errorf("cooldown_minutes must be between 1-120")
	placeholder
		settings.CooldownMinutes = 10 // 禁用状态下归一化为默认值
placeholder

	data, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("marshal overload cooldown settings: %w", err)
placeholder

	return s.settingRepo.Set(ctx, SettingKeyOverloadCooldownSettings, string(data))
placeholder

// GetOIDCConnectOAuthConfig 返回用于登录的“最终生效” OIDC 配置。
//
// 优先级：
// - 若对应系统设置键存在，则覆盖 config.yaml/env 的值
// - 否则回退到 config.yaml/env 的值
func (s *SettingService) GetOIDCConnectOAuthConfig(ctx context.Context) (config.OIDCConnectConfig, error) {
	if s == nil || s.cfg == nil {
		return config.OIDCConnectConfig{placeholder, infraerrors.ServiceUnavailable("CONFIG_NOT_READY", "config not loaded")
placeholder

	effective := s.cfg.OIDC

	keys := []string{
		SettingKeyOIDCConnectEnabled,
		SettingKeyOIDCConnectProviderName,
		SettingKeyOIDCConnectClientID,
		SettingKeyOIDCConnectClientSecret,
		SettingKeyOIDCConnectIssuerURL,
		SettingKeyOIDCConnectDiscoveryURL,
		SettingKeyOIDCConnectAuthorizeURL,
		SettingKeyOIDCConnectTokenURL,
		SettingKeyOIDCConnectUserInfoURL,
		SettingKeyOIDCConnectJWKSURL,
		SettingKeyOIDCConnectScopes,
		SettingKeyOIDCConnectRedirectURL,
		SettingKeyOIDCConnectFrontendRedirectURL,
		SettingKeyOIDCConnectTokenAuthMethod,
		SettingKeyOIDCConnectUsePKCE,
		SettingKeyOIDCConnectValidateIDToken,
		SettingKeyOIDCConnectAllowedSigningAlgs,
		SettingKeyOIDCConnectClockSkewSeconds,
		SettingKeyOIDCConnectRequireEmailVerified,
		SettingKeyOIDCConnectUserInfoEmailPath,
		SettingKeyOIDCConnectUserInfoIDPath,
		SettingKeyOIDCConnectUserInfoUsernamePath,
placeholder
	settings, err := s.settingRepo.GetMultiple(ctx, keys)
	if err != nil {
		return config.OIDCConnectConfig{placeholder, fmt.Errorf("get oidc connect settings: %w", err)
placeholder

	if raw, ok := settings[SettingKeyOIDCConnectEnabled]; ok {
		effective.Enabled = raw == "true"
placeholder
	if v, ok := settings[SettingKeyOIDCConnectProviderName]; ok && strings.TrimSpace(v) != "" {
		effective.ProviderName = strings.TrimSpace(v)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectClientID]; ok && strings.TrimSpace(v) != "" {
		effective.ClientID = strings.TrimSpace(v)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectClientSecret]; ok && strings.TrimSpace(v) != "" {
		effective.ClientSecret = strings.TrimSpace(v)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectIssuerURL]; ok && strings.TrimSpace(v) != "" {
		effective.IssuerURL = strings.TrimSpace(v)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectDiscoveryURL]; ok && strings.TrimSpace(v) != "" {
		effective.DiscoveryURL = strings.TrimSpace(v)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectAuthorizeURL]; ok && strings.TrimSpace(v) != "" {
		effective.AuthorizeURL = strings.TrimSpace(v)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectTokenURL]; ok && strings.TrimSpace(v) != "" {
		effective.TokenURL = strings.TrimSpace(v)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectUserInfoURL]; ok && strings.TrimSpace(v) != "" {
		effective.UserInfoURL = strings.TrimSpace(v)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectJWKSURL]; ok && strings.TrimSpace(v) != "" {
		effective.JWKSURL = strings.TrimSpace(v)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectScopes]; ok && strings.TrimSpace(v) != "" {
		effective.Scopes = strings.TrimSpace(v)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectRedirectURL]; ok && strings.TrimSpace(v) != "" {
		effective.RedirectURL = strings.TrimSpace(v)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectFrontendRedirectURL]; ok && strings.TrimSpace(v) != "" {
		effective.FrontendRedirectURL = strings.TrimSpace(v)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectTokenAuthMethod]; ok && strings.TrimSpace(v) != "" {
		effective.TokenAuthMethod = strings.ToLower(strings.TrimSpace(v))
placeholder
	if raw, ok := settings[SettingKeyOIDCConnectUsePKCE]; ok {
		effective.UsePKCE = raw == "true"
placeholder
	if raw, ok := settings[SettingKeyOIDCConnectValidateIDToken]; ok {
		effective.ValidateIDToken = raw == "true"
placeholder
	if v, ok := settings[SettingKeyOIDCConnectAllowedSigningAlgs]; ok && strings.TrimSpace(v) != "" {
		effective.AllowedSigningAlgs = strings.TrimSpace(v)
placeholder
	if raw, ok := settings[SettingKeyOIDCConnectClockSkewSeconds]; ok && strings.TrimSpace(raw) != "" {
		if parsed, parseErr := strconv.Atoi(strings.TrimSpace(raw)); parseErr == nil {
			effective.ClockSkewSeconds = parsed
	placeholder
placeholder
	if raw, ok := settings[SettingKeyOIDCConnectRequireEmailVerified]; ok {
		effective.RequireEmailVerified = raw == "true"
placeholder
	if v, ok := settings[SettingKeyOIDCConnectUserInfoEmailPath]; ok {
		effective.UserInfoEmailPath = strings.TrimSpace(v)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectUserInfoIDPath]; ok {
		effective.UserInfoIDPath = strings.TrimSpace(v)
placeholder
	if v, ok := settings[SettingKeyOIDCConnectUserInfoUsernamePath]; ok {
		effective.UserInfoUsernamePath = strings.TrimSpace(v)
placeholder

	if !effective.Enabled {
		return config.OIDCConnectConfig{placeholder, infraerrors.NotFound("OAUTH_DISABLED", "oauth login is disabled")
placeholder
	if strings.TrimSpace(effective.ProviderName) == "" {
		effective.ProviderName = "OIDC"
placeholder
	if strings.TrimSpace(effective.ClientID) == "" {
		return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth client id not configured")
placeholder
	if strings.TrimSpace(effective.IssuerURL) == "" {
		return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth issuer url not configured")
placeholder
	if strings.TrimSpace(effective.RedirectURL) == "" {
		return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth redirect url not configured")
placeholder
	if strings.TrimSpace(effective.FrontendRedirectURL) == "" {
		return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth frontend redirect url not configured")
placeholder
	if !scopesContainOpenID(effective.Scopes) {
		return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth scopes must contain openid")
placeholder
	if effective.ClockSkewSeconds < 0 || effective.ClockSkewSeconds > 600 {
		return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth clock skew must be between 0 and 600")
placeholder

	if err := config.ValidateAbsoluteHTTPURL(effective.IssuerURL); err != nil {
		return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth issuer url invalid")
placeholder

	discoveryURL := strings.TrimSpace(effective.DiscoveryURL)
	if discoveryURL == "" {
		discoveryURL = oidcDefaultDiscoveryURL(effective.IssuerURL)
		effective.DiscoveryURL = discoveryURL
placeholder
	if discoveryURL != "" {
		if err := config.ValidateAbsoluteHTTPURL(discoveryURL); err != nil {
			return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth discovery url invalid")
	placeholder
placeholder

	needsDiscovery := strings.TrimSpace(effective.AuthorizeURL) == "" ||
		strings.TrimSpace(effective.TokenURL) == "" ||
		(effective.ValidateIDToken && strings.TrimSpace(effective.JWKSURL) == "")
	if needsDiscovery && discoveryURL != "" {
		metadata, resolveErr := oidcResolveProviderMetadata(ctx, discoveryURL)
		if resolveErr != nil {
			return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth discovery resolve failed").WithCause(resolveErr)
	placeholder
		if strings.TrimSpace(effective.AuthorizeURL) == "" {
			effective.AuthorizeURL = strings.TrimSpace(metadata.AuthorizationEndpoint)
	placeholder
		if strings.TrimSpace(effective.TokenURL) == "" {
			effective.TokenURL = strings.TrimSpace(metadata.TokenEndpoint)
	placeholder
		if strings.TrimSpace(effective.UserInfoURL) == "" {
			effective.UserInfoURL = strings.TrimSpace(metadata.UserInfoEndpoint)
	placeholder
		if strings.TrimSpace(effective.JWKSURL) == "" {
			effective.JWKSURL = strings.TrimSpace(metadata.JWKSURI)
	placeholder
placeholder

	if strings.TrimSpace(effective.AuthorizeURL) == "" {
		return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth authorize url not configured")
placeholder
	if strings.TrimSpace(effective.TokenURL) == "" {
		return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth token url not configured")
placeholder
	if err := config.ValidateAbsoluteHTTPURL(effective.AuthorizeURL); err != nil {
		return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth authorize url invalid")
placeholder
	if err := config.ValidateAbsoluteHTTPURL(effective.TokenURL); err != nil {
		return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth token url invalid")
placeholder
	if v := strings.TrimSpace(effective.UserInfoURL); v != "" {
		if err := config.ValidateAbsoluteHTTPURL(v); err != nil {
			return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth userinfo url invalid")
	placeholder
placeholder
	if effective.ValidateIDToken {
		if strings.TrimSpace(effective.JWKSURL) == "" {
			return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth jwks url not configured")
	placeholder
		if strings.TrimSpace(effective.AllowedSigningAlgs) == "" {
			return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth signing algs not configured")
	placeholder
placeholder
	if v := strings.TrimSpace(effective.JWKSURL); v != "" {
		if err := config.ValidateAbsoluteHTTPURL(v); err != nil {
			return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth jwks url invalid")
	placeholder
placeholder
	if err := config.ValidateAbsoluteHTTPURL(effective.RedirectURL); err != nil {
		return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth redirect url invalid")
placeholder
	if err := config.ValidateFrontendRedirectURL(effective.FrontendRedirectURL); err != nil {
		return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth frontend redirect url invalid")
placeholder

	method := strings.ToLower(strings.TrimSpace(effective.TokenAuthMethod))
	switch method {
	case "", "client_secret_post", "client_secret_basic":
		if strings.TrimSpace(effective.ClientSecret) == "" {
			return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth client secret not configured")
	placeholder
	case "none":
		if !effective.UsePKCE {
			return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth pkce must be enabled when token_auth_method=none")
	placeholder
	default:
		return config.OIDCConnectConfig{placeholder, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth token_auth_method invalid")
placeholder

	return effective, nil
placeholder

func scopesContainOpenID(scopes string) bool {
	for _, scope := range strings.Fields(strings.ToLower(strings.TrimSpace(scopes))) {
		if scope == "openid" {
			return true
	placeholder
placeholder
	return false
placeholder

type oidcProviderMetadata struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	UserInfoEndpoint      string `json:"userinfo_endpoint"`
	JWKSURI               string `json:"jwks_uri"`
placeholder

func oidcDefaultDiscoveryURL(issuerURL string) string {
	issuerURL = strings.TrimSpace(issuerURL)
	if issuerURL == "" {
		return ""
placeholder
	return strings.TrimRight(issuerURL, "/") + "/.well-known/openid-configuration"
placeholder

func oidcResolveProviderMetadata(ctx context.Context, discoveryURL string) (*oidcProviderMetadata, error) {
	discoveryURL = strings.TrimSpace(discoveryURL)
	if discoveryURL == "" {
		return nil, fmt.Errorf("discovery url is empty")
placeholder

	resp, err := req.C().
		SetTimeout(15*time.Second).
		R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		Get(discoveryURL)
	if err != nil {
		return nil, fmt.Errorf("request discovery document: %w", err)
placeholder
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("discovery request failed: status=%d", resp.StatusCode)
placeholder

	metadata := &oidcProviderMetadata{placeholder
	if err := json.Unmarshal(resp.Bytes(), metadata); err != nil {
		return nil, fmt.Errorf("parse discovery document: %w", err)
placeholder
	return metadata, nil
placeholder

// GetStreamTimeoutSettings 获取流超时处理配置
func (s *SettingService) GetStreamTimeoutSettings(ctx context.Context) (*StreamTimeoutSettings, error) {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyStreamTimeoutSettings)
	if err != nil {
		if errors.Is(err, ErrSettingNotFound) {
			return DefaultStreamTimeoutSettings(), nil
	placeholder
		return nil, fmt.Errorf("get stream timeout settings: %w", err)
placeholder
	if value == "" {
		return DefaultStreamTimeoutSettings(), nil
placeholder

	var settings StreamTimeoutSettings
	if err := json.Unmarshal([]byte(value), &settings); err != nil {
		return DefaultStreamTimeoutSettings(), nil
placeholder

	// 验证并修正配置值
	if settings.TempUnschedMinutes < 1 {
		settings.TempUnschedMinutes = 1
placeholder
	if settings.TempUnschedMinutes > 60 {
		settings.TempUnschedMinutes = 60
placeholder
	if settings.ThresholdCount < 1 {
		settings.ThresholdCount = 1
placeholder
	if settings.ThresholdCount > 10 {
		settings.ThresholdCount = 10
placeholder
	if settings.ThresholdWindowMinutes < 1 {
		settings.ThresholdWindowMinutes = 1
placeholder
	if settings.ThresholdWindowMinutes > 60 {
		settings.ThresholdWindowMinutes = 60
placeholder

	// 验证 action
	switch settings.Action {
	case StreamTimeoutActionTempUnsched, StreamTimeoutActionError, StreamTimeoutActionNone:
		// valid
	default:
		settings.Action = StreamTimeoutActionTempUnsched
placeholder

	return &settings, nil
placeholder

// IsUngroupedKeySchedulingAllowed 查询是否允许未分组 Key 调度
func (s *SettingService) IsUngroupedKeySchedulingAllowed(ctx context.Context) bool {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyAllowUngroupedKeyScheduling)
	if err != nil {
		return false // fail-closed: 查询失败时默认不允许
placeholder
	return value == "true"
placeholder

// GetClaudeCodeVersionBounds 获取 Claude Code 版本号上下限要求
// 使用进程内 atomic.Value 缓存，60 秒 TTL，热路径零锁开销
// singleflight 防止缓存过期时 thundering herd
// 返回空字符串表示不做对应方向的版本检查
func (s *SettingService) GetClaudeCodeVersionBounds(ctx context.Context) (min, max string) {
	if cached, ok := versionBoundsCache.Load().(*cachedVersionBounds); ok {
		if time.Now().UnixNano() < cached.expiresAt {
			return cached.min, cached.max
	placeholder
placeholder
	// singleflight: 同一时刻只有一个 goroutine 查询 DB，其余复用结果
	type bounds struct{ min, max string placeholder
	result, err, _ := versionBoundsSF.Do("version_bounds", func() (any, error) {
		// 二次检查，避免排队的 goroutine 重复查询
		if cached, ok := versionBoundsCache.Load().(*cachedVersionBounds); ok {
			if time.Now().UnixNano() < cached.expiresAt {
				return bounds{cached.min, cached.maxplaceholder, nil
		placeholder
	placeholder
		// 使用独立 context：断开请求取消链，避免客户端断连导致空值被长期缓存
		dbCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), versionBoundsDBTimeout)
		defer cancel()
		values, err := s.settingRepo.GetMultiple(dbCtx, []string{
			SettingKeyMinClaudeCodeVersion,
			SettingKeyMaxClaudeCodeVersion,
	placeholder)
		if err != nil {
			// fail-open: DB 错误时不阻塞请求，但记录日志并使用短 TTL 快速重试
			slog.Warn("failed to get claude code version bounds setting, skipping version check", "error", err)
			versionBoundsCache.Store(&cachedVersionBounds{
				min:       "",
				max:       "",
				expiresAt: time.Now().Add(versionBoundsErrorTTL).UnixNano(),
		placeholder)
			return bounds{"", ""placeholder, nil
	placeholder
		b := bounds{
			min: values[SettingKeyMinClaudeCodeVersion],
			max: values[SettingKeyMaxClaudeCodeVersion],
	placeholder
		versionBoundsCache.Store(&cachedVersionBounds{
			min:       b.min,
			max:       b.max,
			expiresAt: time.Now().Add(versionBoundsCacheTTL).UnixNano(),
	placeholder)
		return b, nil
placeholder)
	if err != nil {
		return "", ""
placeholder
	b, ok := result.(bounds)
	if !ok {
		return "", ""
placeholder
	return b.min, b.max
placeholder

// GetRectifierSettings 获取请求整流器配置
func (s *SettingService) GetRectifierSettings(ctx context.Context) (*RectifierSettings, error) {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyRectifierSettings)
	if err != nil {
		if errors.Is(err, ErrSettingNotFound) {
			return DefaultRectifierSettings(), nil
	placeholder
		return nil, fmt.Errorf("get rectifier settings: %w", err)
placeholder
	if value == "" {
		return DefaultRectifierSettings(), nil
placeholder

	var settings RectifierSettings
	if err := json.Unmarshal([]byte(value), &settings); err != nil {
		return DefaultRectifierSettings(), nil
placeholder

	return &settings, nil
placeholder

// SetRectifierSettings 设置请求整流器配置
func (s *SettingService) SetRectifierSettings(ctx context.Context, settings *RectifierSettings) error {
	if settings == nil {
		return fmt.Errorf("settings cannot be nil")
placeholder

	data, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("marshal rectifier settings: %w", err)
placeholder

	return s.settingRepo.Set(ctx, SettingKeyRectifierSettings, string(data))
placeholder

// IsSignatureRectifierEnabled 判断签名整流是否启用（总开关 && 签名子开关）
func (s *SettingService) IsSignatureRectifierEnabled(ctx context.Context) bool {
	settings, err := s.GetRectifierSettings(ctx)
	if err != nil {
		return true // fail-open: 查询失败时默认启用
placeholder
	return settings.Enabled && settings.ThinkingSignatureEnabled
placeholder

// IsBudgetRectifierEnabled 判断 Budget 整流是否启用（总开关 && Budget 子开关）
func (s *SettingService) IsBudgetRectifierEnabled(ctx context.Context) bool {
	settings, err := s.GetRectifierSettings(ctx)
	if err != nil {
		return true // fail-open: 查询失败时默认启用
placeholder
	return settings.Enabled && settings.ThinkingBudgetEnabled
placeholder

// GetBetaPolicySettings 获取 Beta 策略配置
func (s *SettingService) GetBetaPolicySettings(ctx context.Context) (*BetaPolicySettings, error) {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyBetaPolicySettings)
	if err != nil {
		if errors.Is(err, ErrSettingNotFound) {
			return DefaultBetaPolicySettings(), nil
	placeholder
		return nil, fmt.Errorf("get beta policy settings: %w", err)
placeholder
	if value == "" {
		return DefaultBetaPolicySettings(), nil
placeholder

	var settings BetaPolicySettings
	if err := json.Unmarshal([]byte(value), &settings); err != nil {
		return DefaultBetaPolicySettings(), nil
placeholder

	return &settings, nil
placeholder

// SetBetaPolicySettings 设置 Beta 策略配置
func (s *SettingService) SetBetaPolicySettings(ctx context.Context, settings *BetaPolicySettings) error {
	if settings == nil {
		return fmt.Errorf("settings cannot be nil")
placeholder

	validActions := map[string]bool{
		BetaPolicyActionPass: true, BetaPolicyActionFilter: true, BetaPolicyActionBlock: true,
placeholder
	validScopes := map[string]bool{
		BetaPolicyScopeAll: true, BetaPolicyScopeOAuth: true, BetaPolicyScopeAPIKey: true, BetaPolicyScopeBedrock: true,
placeholder

	for i, rule := range settings.Rules {
		if rule.BetaToken == "" {
			return fmt.Errorf("rule[%d]: beta_token cannot be empty", i)
	placeholder
		if !validActions[rule.Action] {
			return fmt.Errorf("rule[%d]: invalid action %q", i, rule.Action)
	placeholder
		if !validScopes[rule.Scope] {
			return fmt.Errorf("rule[%d]: invalid scope %q", i, rule.Scope)
	placeholder
		// Validate model_whitelist patterns
		for j, pattern := range rule.ModelWhitelist {
			trimmed := strings.TrimSpace(pattern)
			if trimmed == "" {
				return fmt.Errorf("rule[%d]: model_whitelist[%d] cannot be empty", i, j)
		placeholder
			settings.Rules[i].ModelWhitelist[j] = trimmed
	placeholder
		// Validate fallback_action
		if rule.FallbackAction != "" && !validActions[rule.FallbackAction] {
			return fmt.Errorf("rule[%d]: invalid fallback_action %q", i, rule.FallbackAction)
	placeholder
placeholder

	data, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("marshal beta policy settings: %w", err)
placeholder

	return s.settingRepo.Set(ctx, SettingKeyBetaPolicySettings, string(data))
placeholder

// SetStreamTimeoutSettings 设置流超时处理配置
func (s *SettingService) SetStreamTimeoutSettings(ctx context.Context, settings *StreamTimeoutSettings) error {
	if settings == nil {
		return fmt.Errorf("settings cannot be nil")
placeholder

	// 验证配置值
	if settings.TempUnschedMinutes < 1 || settings.TempUnschedMinutes > 60 {
		return fmt.Errorf("temp_unsched_minutes must be between 1-60")
placeholder
	if settings.ThresholdCount < 1 || settings.ThresholdCount > 10 {
		return fmt.Errorf("threshold_count must be between 1-10")
placeholder
	if settings.ThresholdWindowMinutes < 1 || settings.ThresholdWindowMinutes > 60 {
		return fmt.Errorf("threshold_window_minutes must be between 1-60")
placeholder

	switch settings.Action {
	case StreamTimeoutActionTempUnsched, StreamTimeoutActionError, StreamTimeoutActionNone:
		// valid
	default:
		return fmt.Errorf("invalid action: %s", settings.Action)
placeholder

	data, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("marshal stream timeout settings: %w", err)
placeholder

	return s.settingRepo.Set(ctx, SettingKeyStreamTimeoutSettings, string(data))
placeholder
