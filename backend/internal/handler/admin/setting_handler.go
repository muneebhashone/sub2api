package admin

import (
	"log"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// SettingHandler 系统设置处理器
type SettingHandler struct {
	settingService   *service.SettingService
	emailService     *service.EmailService
	turnstileService *service.TurnstileService
placeholder

// NewSettingHandler 创建系统设置处理器
func NewSettingHandler(settingService *service.SettingService, emailService *service.EmailService, turnstileService *service.TurnstileService) *SettingHandler {
	return &SettingHandler{
		settingService:   settingService,
		emailService:     emailService,
		turnstileService: turnstileService,
placeholder
placeholder

// GetSettings 获取所有系统设置
// GET /api/v1/admin/settings
func (h *SettingHandler) GetSettings(c *gin.Context) {
	settings, err := h.settingService.GetAllSettings(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, dto.SystemSettings{
		RegistrationEnabled:          settings.RegistrationEnabled,
		EmailVerifyEnabled:           settings.EmailVerifyEnabled,
		SMTPHost:                     settings.SMTPHost,
		SMTPPort:                     settings.SMTPPort,
		SMTPUsername:                 settings.SMTPUsername,
		SMTPPasswordConfigured:       settings.SMTPPasswordConfigured,
		SMTPFrom:                     settings.SMTPFrom,
		SMTPFromName:                 settings.SMTPFromName,
		SMTPUseTLS:                   settings.SMTPUseTLS,
		TurnstileEnabled:             settings.TurnstileEnabled,
		TurnstileSiteKey:             settings.TurnstileSiteKey,
		TurnstileSecretKeyConfigured: settings.TurnstileSecretKeyConfigured,
		SiteName:                     settings.SiteName,
		SiteLogo:                     settings.SiteLogo,
		SiteSubtitle:                 settings.SiteSubtitle,
		APIBaseURL:                   settings.APIBaseURL,
		ContactInfo:                  settings.ContactInfo,
		DocURL:                       settings.DocURL,
		DefaultConcurrency:           settings.DefaultConcurrency,
		DefaultBalance:               settings.DefaultBalance,
		EnableModelFallback:          settings.EnableModelFallback,
		FallbackModelAnthropic:       settings.FallbackModelAnthropic,
		FallbackModelOpenAI:          settings.FallbackModelOpenAI,
		FallbackModelGemini:          settings.FallbackModelGemini,
		FallbackModelAntigravity:     settings.FallbackModelAntigravity,
		EnableIdentityPatch:          settings.EnableIdentityPatch,
		IdentityPatchPrompt:          settings.IdentityPatchPrompt,
		OpsMonitoringEnabled:         settings.OpsMonitoringEnabled,
		OpsRealtimeMonitoringEnabled: settings.OpsRealtimeMonitoringEnabled,
		OpsQueryModeDefault:          settings.OpsQueryModeDefault,
		OpsMetricsIntervalSeconds:    settings.OpsMetricsIntervalSeconds,
placeholder)
placeholder

// UpdateSettingsRequest 更新设置请求
type UpdateSettingsRequest struct {
	// 注册设置
	RegistrationEnabled bool `json:"registration_enabled"`
	EmailVerifyEnabled  bool `json:"email_verify_enabled"`

	// 邮件服务设置
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUsername string `json:"smtp_username"`
	SMTPPassword string `json:"smtp_password"`
	SMTPFrom     string `json:"smtp_from_email"`
	SMTPFromName string `json:"smtp_from_name"`
	SMTPUseTLS   bool   `json:"smtp_use_tls"`

	// Cloudflare Turnstile 设置
	TurnstileEnabled   bool   `json:"turnstile_enabled"`
	TurnstileSiteKey   string `json:"turnstile_site_key"`
	TurnstileSecretKey string `json:"turnstile_secret_key"`

	// OEM设置
	SiteName     string `json:"site_name"`
	SiteLogo     string `json:"site_logo"`
	SiteSubtitle string `json:"site_subtitle"`
	APIBaseURL   string `json:"api_base_url"`
	ContactInfo  string `json:"contact_info"`
	DocURL       string `json:"doc_url"`

	// 默认配置
	DefaultConcurrency int     `json:"default_concurrency"`
	DefaultBalance     float64 `json:"default_balance"`

	// Model fallback configuration
	EnableModelFallback      bool   `json:"enable_model_fallback"`
	FallbackModelAnthropic   string `json:"fallback_model_anthropic"`
	FallbackModelOpenAI      string `json:"fallback_model_openai"`
	FallbackModelGemini      string `json:"fallback_model_gemini"`
	FallbackModelAntigravity string `json:"fallback_model_antigravity"`

	// Identity patch configuration (Claude -> Gemini)
	EnableIdentityPatch bool   `json:"enable_identity_patch"`
	IdentityPatchPrompt string `json:"identity_patch_prompt"`

	// Ops monitoring (vNext)
	OpsMonitoringEnabled         *bool   `json:"ops_monitoring_enabled"`
	OpsRealtimeMonitoringEnabled *bool   `json:"ops_realtime_monitoring_enabled"`
	OpsQueryModeDefault          *string `json:"ops_query_mode_default"`
	OpsMetricsIntervalSeconds    *int    `json:"ops_metrics_interval_seconds"`
placeholder

// UpdateSettings 更新系统设置
// PUT /api/v1/admin/settings
func (h *SettingHandler) UpdateSettings(c *gin.Context) {
	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	previousSettings, err := h.settingService.GetAllSettings(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	// 验证参数
	if req.DefaultConcurrency < 1 {
		req.DefaultConcurrency = 1
placeholder
	if req.DefaultBalance < 0 {
		req.DefaultBalance = 0
placeholder
	if req.SMTPPort <= 0 {
		req.SMTPPort = 587
placeholder

	// Turnstile 参数验证
	if req.TurnstileEnabled {
		// 检查必填字段
		if req.TurnstileSiteKey == "" {
			response.BadRequest(c, "Turnstile Site Key is required when enabled")
			return
	placeholder
		// 如果未提供 secret key，使用已保存的值（留空保留当前值）
		if req.TurnstileSecretKey == "" {
			if previousSettings.TurnstileSecretKey == "" {
				response.BadRequest(c, "Turnstile Secret Key is required when enabled")
				return
		placeholder
			req.TurnstileSecretKey = previousSettings.TurnstileSecretKey
	placeholder

		// 当 site_key 或 secret_key 任一变化时验证（避免配置错误导致无法登录）
		siteKeyChanged := previousSettings.TurnstileSiteKey != req.TurnstileSiteKey
		secretKeyChanged := previousSettings.TurnstileSecretKey != req.TurnstileSecretKey
		if siteKeyChanged || secretKeyChanged {
			if err := h.turnstileService.ValidateSecretKey(c.Request.Context(), req.TurnstileSecretKey); err != nil {
				response.ErrorFrom(c, err)
				return
		placeholder
	placeholder
placeholder

	// Ops metrics collector interval validation (seconds).
	if req.OpsMetricsIntervalSeconds != nil {
		v := *req.OpsMetricsIntervalSeconds
		if v < 60 {
			v = 60
	placeholder
		if v > 3600 {
			v = 3600
	placeholder
		req.OpsMetricsIntervalSeconds = &v
placeholder

	settings := &service.SystemSettings{
		RegistrationEnabled:      req.RegistrationEnabled,
		EmailVerifyEnabled:       req.EmailVerifyEnabled,
		SMTPHost:                 req.SMTPHost,
		SMTPPort:                 req.SMTPPort,
		SMTPUsername:             req.SMTPUsername,
		SMTPPassword:             req.SMTPPassword,
		SMTPFrom:                 req.SMTPFrom,
		SMTPFromName:             req.SMTPFromName,
		SMTPUseTLS:               req.SMTPUseTLS,
		TurnstileEnabled:         req.TurnstileEnabled,
		TurnstileSiteKey:         req.TurnstileSiteKey,
		TurnstileSecretKey:       req.TurnstileSecretKey,
		SiteName:                 req.SiteName,
		SiteLogo:                 req.SiteLogo,
		SiteSubtitle:             req.SiteSubtitle,
		APIBaseURL:               req.APIBaseURL,
		ContactInfo:              req.ContactInfo,
		DocURL:                   req.DocURL,
		DefaultConcurrency:       req.DefaultConcurrency,
		DefaultBalance:           req.DefaultBalance,
		EnableModelFallback:      req.EnableModelFallback,
		FallbackModelAnthropic:   req.FallbackModelAnthropic,
		FallbackModelOpenAI:      req.FallbackModelOpenAI,
		FallbackModelGemini:      req.FallbackModelGemini,
		FallbackModelAntigravity: req.FallbackModelAntigravity,
		EnableIdentityPatch:      req.EnableIdentityPatch,
		IdentityPatchPrompt:      req.IdentityPatchPrompt,
		OpsMonitoringEnabled: func() bool {
			if req.OpsMonitoringEnabled != nil {
				return *req.OpsMonitoringEnabled
		placeholder
			return previousSettings.OpsMonitoringEnabled
	placeholder(),
		OpsRealtimeMonitoringEnabled: func() bool {
			if req.OpsRealtimeMonitoringEnabled != nil {
				return *req.OpsRealtimeMonitoringEnabled
		placeholder
			return previousSettings.OpsRealtimeMonitoringEnabled
	placeholder(),
		OpsQueryModeDefault: func() string {
			if req.OpsQueryModeDefault != nil {
				return *req.OpsQueryModeDefault
		placeholder
			return previousSettings.OpsQueryModeDefault
	placeholder(),
		OpsMetricsIntervalSeconds: func() int {
			if req.OpsMetricsIntervalSeconds != nil {
				return *req.OpsMetricsIntervalSeconds
		placeholder
			return previousSettings.OpsMetricsIntervalSeconds
	placeholder(),
placeholder

	if err := h.settingService.UpdateSettings(c.Request.Context(), settings); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	h.auditSettingsUpdate(c, previousSettings, settings, req)

	// 重新获取设置返回
	updatedSettings, err := h.settingService.GetAllSettings(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, dto.SystemSettings{
		RegistrationEnabled:          updatedSettings.RegistrationEnabled,
		EmailVerifyEnabled:           updatedSettings.EmailVerifyEnabled,
		SMTPHost:                     updatedSettings.SMTPHost,
		SMTPPort:                     updatedSettings.SMTPPort,
		SMTPUsername:                 updatedSettings.SMTPUsername,
		SMTPPasswordConfigured:       updatedSettings.SMTPPasswordConfigured,
		SMTPFrom:                     updatedSettings.SMTPFrom,
		SMTPFromName:                 updatedSettings.SMTPFromName,
		SMTPUseTLS:                   updatedSettings.SMTPUseTLS,
		TurnstileEnabled:             updatedSettings.TurnstileEnabled,
		TurnstileSiteKey:             updatedSettings.TurnstileSiteKey,
		TurnstileSecretKeyConfigured: updatedSettings.TurnstileSecretKeyConfigured,
		SiteName:                     updatedSettings.SiteName,
		SiteLogo:                     updatedSettings.SiteLogo,
		SiteSubtitle:                 updatedSettings.SiteSubtitle,
		APIBaseURL:                   updatedSettings.APIBaseURL,
		ContactInfo:                  updatedSettings.ContactInfo,
		DocURL:                       updatedSettings.DocURL,
		DefaultConcurrency:           updatedSettings.DefaultConcurrency,
		DefaultBalance:               updatedSettings.DefaultBalance,
		EnableModelFallback:          updatedSettings.EnableModelFallback,
		FallbackModelAnthropic:       updatedSettings.FallbackModelAnthropic,
		FallbackModelOpenAI:          updatedSettings.FallbackModelOpenAI,
		FallbackModelGemini:          updatedSettings.FallbackModelGemini,
		FallbackModelAntigravity:     updatedSettings.FallbackModelAntigravity,
		EnableIdentityPatch:          updatedSettings.EnableIdentityPatch,
		IdentityPatchPrompt:          updatedSettings.IdentityPatchPrompt,
		OpsMonitoringEnabled:         updatedSettings.OpsMonitoringEnabled,
		OpsRealtimeMonitoringEnabled: updatedSettings.OpsRealtimeMonitoringEnabled,
		OpsQueryModeDefault:          updatedSettings.OpsQueryModeDefault,
		OpsMetricsIntervalSeconds:    updatedSettings.OpsMetricsIntervalSeconds,
placeholder)
placeholder

func (h *SettingHandler) auditSettingsUpdate(c *gin.Context, before *service.SystemSettings, after *service.SystemSettings, req UpdateSettingsRequest) {
	if before == nil || after == nil {
		return
placeholder

	changed := diffSettings(before, after, req)
	if len(changed) == 0 {
		return
placeholder

	subject, _ := middleware.GetAuthSubjectFromContext(c)
	role, _ := middleware.GetUserRoleFromContext(c)
	log.Printf("AUDIT: settings updated at=%s user_id=%d role=%s changed=%v",
		time.Now().UTC().Format(time.RFC3339),
		subject.UserID,
		role,
		changed,
	)
placeholder

func diffSettings(before *service.SystemSettings, after *service.SystemSettings, req UpdateSettingsRequest) []string {
	changed := make([]string, 0, 20)
	if before.RegistrationEnabled != after.RegistrationEnabled {
		changed = append(changed, "registration_enabled")
placeholder
	if before.EmailVerifyEnabled != after.EmailVerifyEnabled {
		changed = append(changed, "email_verify_enabled")
placeholder
	if before.SMTPHost != after.SMTPHost {
		changed = append(changed, "smtp_host")
placeholder
	if before.SMTPPort != after.SMTPPort {
		changed = append(changed, "smtp_port")
placeholder
	if before.SMTPUsername != after.SMTPUsername {
		changed = append(changed, "smtp_username")
placeholder
	if req.SMTPPassword != "" {
		changed = append(changed, "smtp_password")
placeholder
	if before.SMTPFrom != after.SMTPFrom {
		changed = append(changed, "smtp_from_email")
placeholder
	if before.SMTPFromName != after.SMTPFromName {
		changed = append(changed, "smtp_from_name")
placeholder
	if before.SMTPUseTLS != after.SMTPUseTLS {
		changed = append(changed, "smtp_use_tls")
placeholder
	if before.TurnstileEnabled != after.TurnstileEnabled {
		changed = append(changed, "turnstile_enabled")
placeholder
	if before.TurnstileSiteKey != after.TurnstileSiteKey {
		changed = append(changed, "turnstile_site_key")
placeholder
	if req.TurnstileSecretKey != "" {
		changed = append(changed, "turnstile_secret_key")
placeholder
	if before.SiteName != after.SiteName {
		changed = append(changed, "site_name")
placeholder
	if before.SiteLogo != after.SiteLogo {
		changed = append(changed, "site_logo")
placeholder
	if before.SiteSubtitle != after.SiteSubtitle {
		changed = append(changed, "site_subtitle")
placeholder
	if before.APIBaseURL != after.APIBaseURL {
		changed = append(changed, "api_base_url")
placeholder
	if before.ContactInfo != after.ContactInfo {
		changed = append(changed, "contact_info")
placeholder
	if before.DocURL != after.DocURL {
		changed = append(changed, "doc_url")
placeholder
	if before.DefaultConcurrency != after.DefaultConcurrency {
		changed = append(changed, "default_concurrency")
placeholder
	if before.DefaultBalance != after.DefaultBalance {
		changed = append(changed, "default_balance")
placeholder
	if before.EnableModelFallback != after.EnableModelFallback {
		changed = append(changed, "enable_model_fallback")
placeholder
	if before.FallbackModelAnthropic != after.FallbackModelAnthropic {
		changed = append(changed, "fallback_model_anthropic")
placeholder
	if before.FallbackModelOpenAI != after.FallbackModelOpenAI {
		changed = append(changed, "fallback_model_openai")
placeholder
	if before.FallbackModelGemini != after.FallbackModelGemini {
		changed = append(changed, "fallback_model_gemini")
placeholder
	if before.FallbackModelAntigravity != after.FallbackModelAntigravity {
		changed = append(changed, "fallback_model_antigravity")
placeholder
	if before.OpsMonitoringEnabled != after.OpsMonitoringEnabled {
		changed = append(changed, "ops_monitoring_enabled")
placeholder
	if before.OpsRealtimeMonitoringEnabled != after.OpsRealtimeMonitoringEnabled {
		changed = append(changed, "ops_realtime_monitoring_enabled")
placeholder
	if before.OpsQueryModeDefault != after.OpsQueryModeDefault {
		changed = append(changed, "ops_query_mode_default")
placeholder
	if before.OpsMetricsIntervalSeconds != after.OpsMetricsIntervalSeconds {
		changed = append(changed, "ops_metrics_interval_seconds")
placeholder
	return changed
placeholder

// TestSMTPRequest 测试SMTP连接请求
type TestSMTPRequest struct {
	SMTPHost     string `json:"smtp_host" binding:"required"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUsername string `json:"smtp_username"`
	SMTPPassword string `json:"smtp_password"`
	SMTPUseTLS   bool   `json:"smtp_use_tls"`
placeholder

// TestSMTPConnection 测试SMTP连接
// POST /api/v1/admin/settings/test-smtp
func (h *SettingHandler) TestSMTPConnection(c *gin.Context) {
	var req TestSMTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	if req.SMTPPort <= 0 {
		req.SMTPPort = 587
placeholder

	// 如果未提供密码，从数据库获取已保存的密码
	password := req.SMTPPassword
	if password == "" {
		savedConfig, err := h.emailService.GetSMTPConfig(c.Request.Context())
		if err == nil && savedConfig != nil {
			password = savedConfig.Password
	placeholder
placeholder

	config := &service.SMTPConfig{
		Host:     req.SMTPHost,
		Port:     req.SMTPPort,
		Username: req.SMTPUsername,
		Password: password,
		UseTLS:   req.SMTPUseTLS,
placeholder

	err := h.emailService.TestSMTPConnectionWithConfig(config)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, gin.H{"message": "SMTP connection successful"placeholder)
placeholder

// SendTestEmailRequest 发送测试邮件请求
type SendTestEmailRequest struct {
	Email        string `json:"email" binding:"required,email"`
	SMTPHost     string `json:"smtp_host" binding:"required"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUsername string `json:"smtp_username"`
	SMTPPassword string `json:"smtp_password"`
	SMTPFrom     string `json:"smtp_from_email"`
	SMTPFromName string `json:"smtp_from_name"`
	SMTPUseTLS   bool   `json:"smtp_use_tls"`
placeholder

// SendTestEmail 发送测试邮件
// POST /api/v1/admin/settings/send-test-email
func (h *SettingHandler) SendTestEmail(c *gin.Context) {
	var req SendTestEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	if req.SMTPPort <= 0 {
		req.SMTPPort = 587
placeholder

	// 如果未提供密码，从数据库获取已保存的密码
	password := req.SMTPPassword
	if password == "" {
		savedConfig, err := h.emailService.GetSMTPConfig(c.Request.Context())
		if err == nil && savedConfig != nil {
			password = savedConfig.Password
	placeholder
placeholder

	config := &service.SMTPConfig{
		Host:     req.SMTPHost,
		Port:     req.SMTPPort,
		Username: req.SMTPUsername,
		Password: password,
		From:     req.SMTPFrom,
		FromName: req.SMTPFromName,
		UseTLS:   req.SMTPUseTLS,
placeholder

	siteName := h.settingService.GetSiteName(c.Request.Context())
	subject := "[" + siteName + "] Test Email"
	body := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f5f5f5; margin: 0; padding: 20px; placeholder
        .container { max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.1); placeholder
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; text-align: center; placeholder
        .content { padding: 40px 30px; text-align: center; placeholder
        .success { color: #10b981; font-size: 48px; margin-bottom: 20px; placeholder
        .footer { background-color: #f8f9fa; padding: 20px; text-align: center; color: #999; font-size: 12px; placeholder
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>` + siteName + `</h1>
        </div>
        <div class="content">
            <div class="success">✓</div>
            <h2>Email Configuration Successful!</h2>
            <p>This is a test email to verify your SMTP settings are working correctly.</p>
        </div>
        <div class="footer">
            <p>This is an automated test message.</p>
        </div>
    </div>
</body>
</html>
`

	if err := h.emailService.SendEmailWithConfig(config, req.Email, subject, body); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, gin.H{"message": "Test email sent successfully"placeholder)
placeholder

// GetAdminAPIKey 获取管理员 API Key 状态
// GET /api/v1/admin/settings/admin-api-key
func (h *SettingHandler) GetAdminAPIKey(c *gin.Context) {
	maskedKey, exists, err := h.settingService.GetAdminAPIKeyStatus(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, gin.H{
		"exists":     exists,
		"masked_key": maskedKey,
placeholder)
placeholder

// RegenerateAdminAPIKey 生成/重新生成管理员 API Key
// POST /api/v1/admin/settings/admin-api-key/regenerate
func (h *SettingHandler) RegenerateAdminAPIKey(c *gin.Context) {
	key, err := h.settingService.GenerateAdminAPIKey(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, gin.H{
		"key": key, // 完整 key 只在生成时返回一次
placeholder)
placeholder

// DeleteAdminAPIKey 删除管理员 API Key
// DELETE /api/v1/admin/settings/admin-api-key
func (h *SettingHandler) DeleteAdminAPIKey(c *gin.Context) {
	if err := h.settingService.DeleteAdminAPIKey(c.Request.Context()); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, gin.H{"message": "Admin API key deleted"placeholder)
placeholder
