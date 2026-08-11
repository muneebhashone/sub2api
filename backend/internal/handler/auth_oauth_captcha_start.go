package handler

import (
	"net/http"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ip"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type oauthStartCaptchaRequest struct {
	// TurnstileToken 承载阿里云验证码的 captchaVerifyParam（复用既有请求字段名）
	TurnstileToken        string `json:"turnstile_token"`
	TencentCaptchaTicket  string `json:"tencent_captcha_ticket"`
	TencentCaptchaRandstr string `json:"tencent_captcha_randstr"`
placeholder

type oauthStartResponse struct {
	AuthorizeURL string `json:"authorize_url"`
placeholder

func (h *AuthHandler) requireActionCaptchaForOAuthLoginStart(c *gin.Context) bool {
	if strings.HasSuffix(strings.TrimRight(c.Request.URL.Path, "/"), "/bind/start") {
		return true
placeholder

	var req oauthStartCaptchaRequest
	if c.Request.Method == http.MethodPost {
		_ = c.ShouldBindJSON(&req)
placeholder
	if err := h.authService.VerifyActionCaptchaIfEnabled(c.Request.Context(), service.CaptchaProof{
		TurnstileToken: req.TurnstileToken,
		TencentTicket:  req.TencentCaptchaTicket,
		TencentRandstr: req.TencentCaptchaRandstr,
placeholder, ip.GetClientIP(c)); err != nil {
		response.ErrorFrom(c, err)
		return false
placeholder
	return true
placeholder

func respondOAuthStart(c *gin.Context, authorizeURL string) {
	if c.Request.Method == http.MethodPost {
		response.Success(c, oauthStartResponse{AuthorizeURL: authorizeURLplaceholder)
		return
placeholder
	c.Redirect(http.StatusFound, authorizeURL)
placeholder
