package admin

import (
	"fmt"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

type GeminiOAuthHandler struct {
	geminiOAuthService *service.GeminiOAuthService
placeholder

func NewGeminiOAuthHandler(geminiOAuthService *service.GeminiOAuthService) *GeminiOAuthHandler {
	return &GeminiOAuthHandler{geminiOAuthService: geminiOAuthServiceplaceholder
placeholder

// GetCapabilities retrieves OAuth configuration capabilities.
// GET /api/v1/admin/gemini/oauth/capabilities
func (h *GeminiOAuthHandler) GetCapabilities(c *gin.Context) {
	cfg := h.geminiOAuthService.GetOAuthConfig()
	response.Success(c, cfg)
placeholder

type GeminiGenerateAuthURLRequest struct {
	ProxyID   *int64 `json:"proxy_id"`
	ProjectID string `json:"project_id"`
	// OAuth 类型: "code_assist" (需要 project_id) 或 "ai_studio" (不需要 project_id)
	// 默认为 "code_assist" 以保持向后兼容
	OAuthType string `json:"oauth_type"`
placeholder

// GenerateAuthURL generates Google OAuth authorization URL for Gemini.
// POST /api/v1/admin/gemini/oauth/auth-url
func (h *GeminiOAuthHandler) GenerateAuthURL(c *gin.Context) {
	var req GeminiGenerateAuthURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	// 默认使用 code_assist 以保持向后兼容
	oauthType := strings.TrimSpace(req.OAuthType)
	if oauthType == "" {
		oauthType = "code_assist"
placeholder
	if oauthType != "code_assist" && oauthType != "google_one" && oauthType != "ai_studio" {
		response.BadRequest(c, "Invalid oauth_type: must be 'code_assist', 'google_one', or 'ai_studio'")
		return
placeholder

	// Always pass the "hosted" callback URI; the OAuth service may override it depending on
	// oauth_type and whether the built-in Gemini CLI OAuth client is used.
	redirectURI := deriveGeminiRedirectURI(c)
	result, err := h.geminiOAuthService.GenerateAuthURL(c.Request.Context(), req.ProxyID, redirectURI, req.ProjectID, oauthType)
	if err != nil {
		msg := err.Error()
		// Treat missing/invalid OAuth client configuration as a user/config error.
		if strings.Contains(msg, "OAuth client not configured") || strings.Contains(msg, "requires your own OAuth Client") {
			response.BadRequest(c, "Failed to generate auth URL: "+msg)
			return
	placeholder
		response.InternalError(c, "Failed to generate auth URL: "+msg)
		return
placeholder

	response.Success(c, result)
placeholder

type GeminiExchangeCodeRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	State     string `json:"state" binding:"required"`
	Code      string `json:"code" binding:"required"`
	ProxyID   *int64 `json:"proxy_id"`
	// OAuth 类型: "code_assist" 或 "ai_studio"，需要与 GenerateAuthURL 时的类型一致
	OAuthType string `json:"oauth_type"`
placeholder

// ExchangeCode exchanges authorization code for tokens.
// POST /api/v1/admin/gemini/oauth/exchange-code
func (h *GeminiOAuthHandler) ExchangeCode(c *gin.Context) {
	var req GeminiExchangeCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	// 默认使用 code_assist 以保持向后兼容
	oauthType := strings.TrimSpace(req.OAuthType)
	if oauthType == "" {
		oauthType = "code_assist"
placeholder
	if oauthType != "code_assist" && oauthType != "google_one" && oauthType != "ai_studio" {
		response.BadRequest(c, "Invalid oauth_type: must be 'code_assist', 'google_one', or 'ai_studio'")
		return
placeholder

	tokenInfo, err := h.geminiOAuthService.ExchangeCode(c.Request.Context(), &service.GeminiExchangeCodeInput{
		SessionID: req.SessionID,
		State:     req.State,
		Code:      req.Code,
		ProxyID:   req.ProxyID,
		OAuthType: oauthType,
placeholder)
	if err != nil {
		response.BadRequest(c, "Failed to exchange code: "+err.Error())
		return
placeholder

	response.Success(c, tokenInfo)
placeholder

func deriveGeminiRedirectURI(c *gin.Context) string {
	origin := strings.TrimSpace(c.GetHeader("Origin"))
	if origin != "" {
		return strings.TrimRight(origin, "/") + "/auth/callback"
placeholder

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
placeholder
	if xfProto := strings.TrimSpace(c.GetHeader("X-Forwarded-Proto")); xfProto != "" {
		scheme = strings.TrimSpace(strings.Split(xfProto, ",")[0])
placeholder

	host := strings.TrimSpace(c.Request.Host)
	if xfHost := strings.TrimSpace(c.GetHeader("X-Forwarded-Host")); xfHost != "" {
		host = strings.TrimSpace(strings.Split(xfHost, ",")[0])
placeholder

	return fmt.Sprintf("%s://%s/auth/callback", scheme, host)
placeholder
