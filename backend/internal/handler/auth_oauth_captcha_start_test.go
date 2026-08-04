//go:build unit

package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type oauthCaptchaSettingRepo struct {
	values map[string]string
placeholder

func (r *oauthCaptchaSettingRepo) Get(context.Context, string) (*service.Setting, error) {
	return nil, service.ErrSettingNotFound
placeholder
func (r *oauthCaptchaSettingRepo) GetValue(_ context.Context, key string) (string, error) {
	value, ok := r.values[key]
	if !ok {
		return "", service.ErrSettingNotFound
placeholder
	return value, nil
placeholder
func (r *oauthCaptchaSettingRepo) Set(context.Context, string, string) error { return nil placeholder
func (r *oauthCaptchaSettingRepo) GetMultiple(_ context.Context, keys []string) (map[string]string, error) {
	values := make(map[string]string, len(keys))
	for _, key := range keys {
		if value, ok := r.values[key]; ok {
			values[key] = value
	placeholder
placeholder
	return values, nil
placeholder
func (r *oauthCaptchaSettingRepo) SetMultiple(context.Context, map[string]string) error {
	return nil
placeholder
func (r *oauthCaptchaSettingRepo) GetAll(context.Context) (map[string]string, error) {
	return r.values, nil
placeholder
func (r *oauthCaptchaSettingRepo) Delete(context.Context, string) error { return nil placeholder

type oauthCaptchaVerifier struct {
	calls int
	proof service.TencentCaptchaProof
placeholder

func (v *oauthCaptchaVerifier) VerifyTicket(_ context.Context, _ service.TencentCaptchaCredentials, proof service.TencentCaptchaProof, _ string) (*service.TencentCaptchaVerifyResponse, error) {
	v.calls++
	v.proof = proof
	return &service.TencentCaptchaVerifyResponse{CaptchaCode: 1placeholder, nil
placeholder

func newOAuthCaptchaTestHandler(enabled bool) (*AuthHandler, *oauthCaptchaVerifier) {
	values := map[string]string{placeholder
	if enabled {
		values = map[string]string{
			service.SettingKeyTencentCaptchaEnabled:        "true",
			service.SettingKeyTencentCaptchaAppID:          "123456789",
			service.SettingKeyTencentCaptchaAppSecretKey:   "app-secret",
			service.SettingKeyTencentCaptchaCloudSecretID:  "cloud-secret-id",
			service.SettingKeyTencentCaptchaCloudSecretKey: "cloud-secret-key",
	placeholder
placeholder
	cfg := &config.Config{placeholder
	settings := service.NewSettingService(&oauthCaptchaSettingRepo{values: valuesplaceholder, cfg)
	verifier := &oauthCaptchaVerifier{placeholder
	authService := service.NewAuthService(nil, nil, nil, nil, cfg, settings, nil, nil, nil, nil, nil, nil, nil)
	authService.SetTencentCaptchaService(service.NewTencentCaptchaService(settings, verifier))
	return &AuthHandler{authService: authService, settingSvc: settings, cfg: cfgplaceholder, verifier
placeholder

func oauthStartHandlers() map[string]func(*AuthHandler, *gin.Context) {
	return map[string]func(*AuthHandler, *gin.Context){
		"github":   func(h *AuthHandler, c *gin.Context) { h.GitHubOAuthStart(c) placeholder,
		"google":   func(h *AuthHandler, c *gin.Context) { h.GoogleOAuthStart(c) placeholder,
		"linuxdo":  func(h *AuthHandler, c *gin.Context) { h.LinuxDoOAuthStart(c) placeholder,
		"dingtalk": func(h *AuthHandler, c *gin.Context) { h.DingTalkOAuthStart(c) placeholder,
		"wechat":   func(h *AuthHandler, c *gin.Context) { h.WeChatOAuthStart(c) placeholder,
		"oidc":     func(h *AuthHandler, c *gin.Context) { h.OIDCOAuthStart(c) placeholder,
placeholder
placeholder

func TestOAuthStartGetRejectsAnonymousLoginWhenTencentEnabledWithoutSideEffects(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for provider, start := range oauthStartHandlers() {
		t.Run(provider, func(t *testing.T) {
			handler, verifier := newOAuthCaptchaTestHandler(true)
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/auth/oauth/"+provider+"/start?intent=bind_current_user", nil)

			start(handler, c)

			require.Equal(t, http.StatusBadRequest, recorder.Code)
			require.Contains(t, recorder.Body.String(), "TENCENT_CAPTCHA_VERIFICATION_FAILED")
			require.Empty(t, recorder.Header().Get("Location"))
			require.Empty(t, recorder.Header().Values("Set-Cookie"))
			require.Zero(t, verifier.calls)
	placeholder)
placeholder
placeholder

func TestOAuthStartPostReturnsAuthorizeURLAfterTencentVerification(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for provider := range oauthStartHandlers() {
		t.Run(provider, func(t *testing.T) {
			handler, verifier := newOAuthCaptchaTestHandler(true)
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request = httptest.NewRequest(
				http.MethodPost,
				"/api/v1/auth/oauth/"+provider+"/start",
				strings.NewReader(`{"tencent_captcha_ticket":"ticket-value","tencent_captcha_randstr":"@rand-value"placeholder`),
			)
			c.Request.Header.Set("Content-Type", "application/json")

			require.True(t, handler.requireActionCaptchaForOAuthLoginStart(c))
			respondOAuthStart(c, "https://provider.example/authorize")

			require.Equal(t, http.StatusOK, recorder.Code)
			require.Contains(t, recorder.Body.String(), `"authorize_url":"https://provider.example/authorize"`)
			require.Equal(t, 1, verifier.calls)
			require.Equal(t, service.TencentCaptchaProof{Ticket: "ticket-value", Randstr: "@rand-value"placeholder, verifier.proof)
	placeholder)
placeholder
placeholder

func TestOAuthStartPostRequiresTencentProofWhenEnabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for provider := range oauthStartHandlers() {
		t.Run(provider, func(t *testing.T) {
			handler, verifier := newOAuthCaptchaTestHandler(true)
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/oauth/"+provider+"/start", strings.NewReader(`{placeholder`))
			c.Request.Header.Set("Content-Type", "application/json")

			require.False(t, handler.requireActionCaptchaForOAuthLoginStart(c))
			require.Equal(t, http.StatusBadRequest, recorder.Code)
			require.Contains(t, recorder.Body.String(), "TENCENT_CAPTCHA_VERIFICATION_FAILED")
			require.Zero(t, verifier.calls)
	placeholder)
placeholder
placeholder

func TestOAuthBindingPathRemainsOutsideTencentGate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := &AuthHandler{placeholder
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/auth/oauth/oidc/bind/start", nil)

	require.True(t, handler.requireActionCaptchaForOAuthLoginStart(c))
	require.Equal(t, http.StatusOK, recorder.Code)
placeholder

func TestOAuthStartGetRemainsCompatibleWhenTencentDisabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, verifier := newOAuthCaptchaTestHandler(false)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/auth/oauth/github/start", nil)

	require.True(t, handler.requireActionCaptchaForOAuthLoginStart(c))
	respondOAuthStart(c, "https://provider.example/authorize")

	require.Equal(t, http.StatusFound, recorder.Code)
	require.Equal(t, "https://provider.example/authorize", recorder.Header().Get("Location"))
	require.Zero(t, verifier.calls)
placeholder
