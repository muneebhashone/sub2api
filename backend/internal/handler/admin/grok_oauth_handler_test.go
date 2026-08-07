//go:build unit

package admin

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type grokQuotaHandlerAccountRepo struct {
	service.AccountRepository
	account *service.Account
	updates map[int64]map[string]any
placeholder

func (r *grokQuotaHandlerAccountRepo) GetByID(_ context.Context, id int64) (*service.Account, error) {
	if r.account != nil && r.account.ID == id {
		return r.account, nil
placeholder
	return nil, service.ErrAccountNotFound
placeholder

func (r *grokQuotaHandlerAccountRepo) UpdateExtra(_ context.Context, id int64, updates map[string]any) error {
	if r.updates == nil {
		r.updates = make(map[int64]map[string]any)
placeholder
	r.updates[id] = updates
	return nil
placeholder

type grokQuotaHandlerUpstream struct {
	mu       sync.Mutex
	requests []*http.Request
	bodies   [][]byte
placeholder

type grokOAuthReconcilerStub struct {
	input  service.GrokOAuthReconcileInput
	calls  int
	result *service.GrokOAuthReconcileResult
	err    error
placeholder

func (s *grokOAuthReconcilerStub) ReconcileGrokOAuth(_ context.Context, input service.GrokOAuthReconcileInput) (*service.GrokOAuthReconcileResult, error) {
	s.calls++
	s.input = input
	return s.result, s.err
placeholder

func (u *grokQuotaHandlerUpstream) Do(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
placeholder
	u.mu.Lock()
	u.requests = append(u.requests, req)
	u.bodies = append(u.bodies, body)
	u.mu.Unlock()
	if req.URL.Path == "/v1/responses" {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header: http.Header{
				"X-Ratelimit-Limit-Requests":     []string{"10"placeholder,
				"X-Ratelimit-Remaining-Requests": []string{"8"placeholder,
		placeholder,
			Body: io.NopCloser(strings.NewReader(`{"id":"resp_probe"placeholder`)),
	placeholder, nil
placeholder
	payload := `{"config":{"billingPeriodStart":"2026-07-01T00:00:00Z","billingPeriodEnd":"2026-08-01T00:00:00Z"placeholderplaceholder`
	if req.URL.RawQuery == "format=credits" {
		payload = `{"config":{"currentPeriod":{"type":"WEEKLY","start":"2026-07-09T03:25:00Z","end":"2026-07-16T03:25:00Z"placeholderplaceholderplaceholder`
placeholder
	return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(payload))placeholder, nil
placeholder

func (u *grokQuotaHandlerUpstream) DoWithTLS(
	req *http.Request,
	proxyURL string,
	accountID int64,
	accountConcurrency int,
	_ *tlsfingerprint.Profile,
) (*http.Response, error) {
	return u.Do(req, proxyURL, accountID, accountConcurrency)
placeholder

func TestGrokOAuthHandlerQueryQuotaProbesUpstream(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &grokQuotaHandlerAccountRepo{account: &service.Account{
		ID:          42,
		Platform:    service.PlatformGrok,
		Type:        service.AccountTypeOAuth,
		Status:      service.StatusActive,
		Schedulable: true,
		Concurrency: 1,
placeholder
			"access_token":  "access-token",
			"refresh_token": "refresh-token",
			"expires_at":    time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339),
	placeholder,
placeholderplaceholder
	upstream := &grokQuotaHandlerUpstream{placeholder
	quotaService := service.NewGrokQuotaService(repo, nil, service.NewGrokTokenProvider(repo, nil), upstream, nil)
	handler := NewGrokOAuthHandler(nil, nil, quotaService, nil)

	router := gin.New()
	router.GET("/api/v1/admin/grok/accounts/:id/quota", handler.QueryQuota)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/grok/accounts/42/quota", nil)
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"source":"hybrid_probe"`)
	require.Contains(t, rec.Body.String(), `"billing":`)
	require.Contains(t, rec.Body.String(), `"snapshot":`)
	require.Contains(t, rec.Body.String(), `"headers_observed":true`)
	require.NotContains(t, rec.Body.String(), "access-token")
	require.Eventually(t, func() bool {
		upstream.mu.Lock()
		defer upstream.mu.Unlock()
		return len(upstream.requests) == 4
placeholder, time.Second, 10*time.Millisecond)
	upstream.mu.Lock()
	requests := append([]*http.Request(nil), upstream.requests...)
	bodies := append([][]byte(nil), upstream.bodies...)
	upstream.mu.Unlock()
	require.Len(t, requests, 4)
	responsesProbeSeen := false
	modelsSyncSeen := false
	for i, upstreamReq := range requests {
		require.Equal(t, "Bearer access-token", upstreamReq.Header.Get("Authorization"))
		if upstreamReq.URL.String() == xai.DefaultCLIBaseURL+"/responses" {
			responsesProbeSeen = true
			require.Equal(t, "application/json, text/event-stream", upstreamReq.Header.Get("Accept"))
			require.Contains(t, string(bodies[i]), `"model":"grok-4.5"`)
			require.Contains(t, string(bodies[i]), `"input":"hi"`)
			require.Contains(t, string(bodies[i]), `"stream":true`)
			require.NotContains(t, string(bodies[i]), `"max_output_tokens"`)
			require.NotContains(t, string(bodies[i]), `"store"`)
	placeholder
		if upstreamReq.URL.String() == xai.DefaultCLIBaseURL+"/models" {
			modelsSyncSeen = true
	placeholder
placeholder
	require.True(t, responsesProbeSeen)
	require.True(t, modelsSyncSeen)
	require.NotNil(t, repo.updates[42])
placeholder

func TestGrokOAuthHandlerResetQuotaReturnsUnsupported(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &grokQuotaHandlerAccountRepo{account: &service.Account{
		ID:       43,
		Platform: service.PlatformGrok,
		Type:     service.AccountTypeOAuth,
placeholderplaceholder
	quotaService := service.NewGrokQuotaService(repo, nil, nil, nil, nil)
	handler := NewGrokOAuthHandler(nil, nil, quotaService, nil)

	router := gin.New()
	router.POST("/api/v1/admin/grok/accounts/:id/reset-quota", handler.ResetQuota)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/grok/accounts/43/reset-quota", nil)
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotImplemented, rec.Code)
	require.Contains(t, rec.Body.String(), `"reason":"GROK_QUOTA_RESET_UNSUPPORTED"`)
	require.NotContains(t, rec.Body.String(), "access-token")
placeholder

func TestGrokOAuthHandlerRuntimeSanityDoesNotExposeSecrets(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv(xai.EnvBaseURL, "http://127.0.0.1:8080/v1?access_token=secret")
	t.Setenv(xai.EnvClientID, "client-secret-like-value")

	handler := NewGrokOAuthHandler(nil, nil, nil, nil)
	router := gin.New()
	router.GET("/api/v1/admin/grok/runtime-sanity", handler.RuntimeSanity)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/grok/runtime-sanity", nil)
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"public_gateway_scope":"responses_only"`)
	require.Contains(t, rec.Body.String(), `"valid":false`)
	require.NotContains(t, rec.Body.String(), "access_token")
	require.NotContains(t, rec.Body.String(), "secret")
	require.NotContains(t, rec.Body.String(), "client-secret-like-value")
placeholder

type grokOAuthHandlerClient struct{placeholder

func (c *grokOAuthHandlerClient) ExchangeCode(context.Context, string, string, string, string, string) (*xai.TokenResponse, error) {
	return nil, errors.New("unexpected exchange")
placeholder

func (c *grokOAuthHandlerClient) RefreshToken(context.Context, string, string, string) (*xai.TokenResponse, error) {
	return &xai.TokenResponse{AccessToken: "access-token", RefreshToken: "refresh-token", ExpiresIn: 3600placeholder, nil
placeholder

func (c *grokOAuthHandlerClient) LoginWithPassword(_ context.Context, email, _ string, _ string) (*service.GrokPasswordLoginResult, error) {
	return &service.GrokPasswordLoginResult{
		Email:    email,
		SSOToken: "sso-from-password",
placeholder, nil
placeholder

func (c *grokOAuthHandlerClient) ConvertSSOToBuild(context.Context, string, string) (*xai.TokenResponse, error) {
	return &xai.TokenResponse{AccessToken: "access-token", RefreshToken: "refresh-token", ExpiresIn: 3600placeholder, nil
placeholder

func TestGrokOAuthHandlerValidateSSOTokenReturnsTokenInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	oauthClient := &grokOAuthHandlerClient{placeholder
	oauthService := service.NewGrokOAuthService(nil, oauthClient)
	defer oauthService.Stop()
	handler := NewGrokOAuthHandler(oauthService, nil, nil, nil)

	router := gin.New()
	router.POST("/api/v1/admin/grok/oauth/sso-token", handler.ValidateSSOToken)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/grok/oauth/sso-token", strings.NewReader(`{"sso_token":"sso-token"placeholder`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"access_token":"access-token"`)
	require.NotContains(t, rec.Body.String(), `"sso_token"`)
placeholder

func TestGrokOAuthHandlerAuthorizePasswordReturnsTokenInfoWithoutPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	oauthClient := &grokOAuthHandlerClient{placeholder
	cfg := &config.Config{placeholder
	cfg.Gateway.Grok.PasswordAuthEnabled = true
	oauthService := service.NewGrokOAuthService(nil, oauthClient, cfg)
	defer oauthService.Stop()
	handler := NewGrokOAuthHandler(oauthService, nil, nil, nil)

	router := gin.New()
	router.POST("/api/v1/admin/grok/oauth/password", handler.AuthorizePassword)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/grok/oauth/password", strings.NewReader(`{"email":"user@example.com","password":"super-secret"placeholder`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"access_token":"access-token"`)
	require.NotContains(t, rec.Body.String(), "super-secret")
placeholder

func TestGrokOAuthHandlerPasswordCapabilityDefaultsToDisabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	oauthService := service.NewGrokOAuthService(nil, &grokOAuthHandlerClient{placeholder)
	defer oauthService.Stop()
	handler := NewGrokOAuthHandler(oauthService, nil, nil, nil)

	router := gin.New()
	router.GET("/api/v1/admin/grok/oauth/capabilities", handler.GetCapabilities)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/admin/grok/oauth/capabilities", nil))

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"password_auth_enabled":false`)
placeholder

func TestGrokSSOImportExpiryUsesTokenExpiryWithoutRefreshToken(t *testing.T) {
	tokenExpiry := time.Now().Add(6 * time.Hour).Unix()
	expiresAt, autoPause := grokSSOImportExpiry(nil, nil, &service.GrokTokenInfo{
		ExpiresAt: tokenExpiry,
placeholder)

	require.NotNil(t, expiresAt)
	require.Equal(t, tokenExpiry, *expiresAt)
	require.NotNil(t, autoPause)
	require.True(t, *autoPause)
placeholder

func TestGrokSSOImportExpiryUsesEarlierRequestedExpiryWithoutRefreshToken(t *testing.T) {
	requestedExpiry := time.Now().Add(2 * time.Hour).Unix()
	tokenExpiry := time.Now().Add(6 * time.Hour).Unix()
	requestedAutoPause := false
	expiresAt, autoPause := grokSSOImportExpiry(&requestedExpiry, &requestedAutoPause, &service.GrokTokenInfo{
		ExpiresAt: tokenExpiry,
placeholder)

	require.NotNil(t, expiresAt)
	require.Equal(t, requestedExpiry, *expiresAt)
	require.NotNil(t, autoPause)
	require.True(t, *autoPause)
placeholder

func TestGrokSSOImportExpiryPreservesRequestSettingsWithRefreshToken(t *testing.T) {
	requestedExpiry := time.Now().Add(2 * time.Hour).Unix()
	requestedAutoPause := false
	expiresAt, autoPause := grokSSOImportExpiry(&requestedExpiry, &requestedAutoPause, &service.GrokTokenInfo{
		RefreshToken: "refresh-token",
		ExpiresAt:    time.Now().Add(6 * time.Hour).Unix(),
placeholder)

	require.Same(t, &requestedExpiry, expiresAt)
	require.Same(t, &requestedAutoPause, autoPause)
placeholder

func TestGrokSSOImportCredentialsPreservesRequestedBaseURL(t *testing.T) {
	built := map[string]any{
		"access_token": "at-1",
		"base_url":     xai.DefaultCLIBaseURL,
placeholder
	reqCredentials := map[string]any{
		"base_url":                "https://relay.example.com/v1",
		"header_override_enabled": true,
		"header_overrides":        map[string]any{"x-relay-key": "k"placeholder,
placeholder

	credentials := grokSSOImportCredentials(built, reqCredentials)

	// token 字段以兑换结果为准；base_url 是运营侧配置，必须保留请求里的自定义地址
	require.Equal(t, "at-1", credentials["access_token"])
	require.Equal(t, "https://relay.example.com/v1", credentials["base_url"])
	require.Equal(t, true, credentials["header_override_enabled"])
	require.Equal(t, map[string]any{"x-relay-key": "k"placeholder, credentials["header_overrides"])
	// 入参不被污染（req.Credentials 会被多个 worker 并发读取）
	require.Equal(t, "https://relay.example.com/v1", reqCredentials["base_url"])
placeholder

func TestGrokSSOImportCredentialsDefaultsToOfficialBaseURL(t *testing.T) {
	built := map[string]any{
		"access_token": "at-1",
		"base_url":     xai.DefaultCLIBaseURL,
placeholder

	credentials := grokSSOImportCredentials(built, nil)
	require.Equal(t, xai.DefaultCLIBaseURL, credentials["base_url"])

	credentials = grokSSOImportCredentials(map[string]any{
		"access_token": "at-2",
		"base_url":     xai.DefaultCLIBaseURL,
placeholder, map[string]any{"base_url": "   "placeholder)
	require.Equal(t, xai.DefaultCLIBaseURL, credentials["base_url"])
	require.Equal(t, "at-2", credentials["access_token"])
placeholder

func TestGrokSSOImportWorkerHandlesMissingOAuthService(t *testing.T) {
	h := &GrokOAuthHandler{placeholder
	result := h.safeCreateAccountFromSSOToken(context.Background(), GrokSSOToOAuthRequest{placeholder, "token", 2, 3)
	require.False(t, result.created)
	require.Equal(t, 2, result.item.Index)
	require.Contains(t, result.item.Error, "GROK_OAUTH_CLIENT_NOT_CONFIGURED")
placeholder

func TestGrokOAuthHandlerReconcileDefaultsToDryRun(t *testing.T) {
	gin.SetMode(gin.TestMode)
	reconciler := &grokOAuthReconcilerStub{result: &service.GrokOAuthReconcileResult{
		DryRun:      true,
		Scanned:     2,
		Actionable:  1,
		WouldBlock:  1,
		Items:       []service.GrokOAuthReconcileItem{{AccountID: 42, Reason: service.GrokOAuthReconcileReasonMissingRefreshToken, Action: service.GrokOAuthReconcileActionBlock, Outcome: service.GrokOAuthReconcileOutcomePlannedplaceholderplaceholder,
		NextAfterID: 0,
placeholderplaceholder
	handler := NewGrokOAuthHandler(nil, nil, nil, reconciler)
	router := gin.New()
	router.POST("/api/v1/admin/grok/oauth/reconcile", handler.ReconcileOAuthAccounts)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/grok/oauth/reconcile", strings.NewReader(`{placeholder`))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, 1, reconciler.calls)
	require.True(t, reconciler.input.DryRun)
	require.False(t, reconciler.input.Apply)
	require.Contains(t, rec.Body.String(), `"reason":"missing_refresh_token"`)
	require.NotContains(t, rec.Body.String(), `"refresh_token":`)
	require.NotContains(t, rec.Body.String(), `"access_token":`)
placeholder

func TestGrokOAuthHandlerReconcileRequiresExplicitApply(t *testing.T) {
	gin.SetMode(gin.TestMode)
	reconciler := &grokOAuthReconcilerStub{placeholder
	handler := NewGrokOAuthHandler(nil, nil, nil, reconciler)
	router := gin.New()
	router.POST("/api/v1/admin/grok/oauth/reconcile", handler.ReconcileOAuthAccounts)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/grok/oauth/reconcile", strings.NewReader(`{"dry_run":falseplaceholder`))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Zero(t, reconciler.calls)
	require.NotContains(t, rec.Body.String(), "credentials")
placeholder

func TestGrokOAuthHandlerReconcileExplicitApply(t *testing.T) {
	gin.SetMode(gin.TestMode)
	reconciler := &grokOAuthReconcilerStub{result: &service.GrokOAuthReconcileResult{DryRun: false, Refreshed: 1placeholderplaceholder
	handler := NewGrokOAuthHandler(nil, nil, nil, reconciler)
	router := gin.New()
	router.POST("/api/v1/admin/grok/oauth/reconcile", handler.ReconcileOAuthAccounts)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/grok/oauth/reconcile", strings.NewReader(`{"apply":true,"dry_run":false,"after_id":10,"limit":25,"refresh_window_seconds":3600placeholder`))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, 1, reconciler.calls)
	require.True(t, reconciler.input.Apply)
	require.False(t, reconciler.input.DryRun)
	require.Equal(t, int64(10), reconciler.input.AfterID)
	require.Equal(t, 25, reconciler.input.Limit)
	require.Equal(t, time.Hour, reconciler.input.RefreshWindow)
placeholder
