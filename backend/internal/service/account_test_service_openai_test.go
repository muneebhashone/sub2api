//go:build unit

package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai_compat"
	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

// --- shared test helpers ---

type queuedHTTPUpstream struct {
	responses []*http.Response
	requests  []*http.Request
	tlsFlags  []bool
placeholder

func (u *queuedHTTPUpstream) Do(_ *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
	return nil, fmt.Errorf("unexpected Do call")
placeholder

func (u *queuedHTTPUpstream) DoWithTLS(req *http.Request, _ string, _ int64, _ int, profile *tlsfingerprint.Profile) (*http.Response, error) {
	u.requests = append(u.requests, req)
	u.tlsFlags = append(u.tlsFlags, profile != nil)
	if len(u.responses) == 0 {
		return nil, fmt.Errorf("no mocked response")
placeholder
	resp := u.responses[0]
	u.responses = u.responses[1:]
	return resp, nil
placeholder

func newJSONResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
placeholder
placeholder

// --- test functions ---

func newTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/accounts/1/test", nil)
	return c, rec
placeholder

type openAIAccountTestRepo struct {
	mockAccountRepoForGemini
	updatedExtra       map[string]any
	bulkUpdatedIDs     []int64
	bulkUpdatedPayload AccountBulkUpdate
	rateLimitedID      int64
	rateLimitedAt      *time.Time
	clearedErrorID     int64
	setErrorID         int64
	setErrorMsg        string
placeholder

func (r *openAIAccountTestRepo) UpdateExtra(_ context.Context, _ int64, updates map[string]any) error {
	r.updatedExtra = updates
	return nil
placeholder

func (r *openAIAccountTestRepo) BulkUpdate(_ context.Context, ids []int64, updates AccountBulkUpdate) (int64, error) {
	r.bulkUpdatedIDs = append([]int64(nil), ids...)
	r.bulkUpdatedPayload = updates
	return int64(len(ids)), nil
placeholder

func (r *openAIAccountTestRepo) SetRateLimited(_ context.Context, id int64, resetAt time.Time) error {
	r.rateLimitedID = id
	r.rateLimitedAt = &resetAt
	return nil
placeholder

func (r *openAIAccountTestRepo) ClearError(_ context.Context, id int64) error {
	r.clearedErrorID = id
	return nil
placeholder

func (r *openAIAccountTestRepo) SetError(_ context.Context, id int64, errorMsg string) error {
	r.setErrorID = id
	r.setErrorMsg = errorMsg
	return nil
placeholder

func TestAccountTestService_OpenAISuccessPersistsSnapshotFromHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, recorder := newTestContext()

	resp := newJSONResponse(http.StatusOK, "")
	resp.Body = io.NopCloser(strings.NewReader(`data: {"type":"response.completed"placeholder

`))
	resp.Header.Set("x-codex-primary-used-percent", "88")
	resp.Header.Set("x-codex-primary-reset-after-seconds", "604800")
	resp.Header.Set("x-codex-primary-window-minutes", "10080")
	resp.Header.Set("x-codex-secondary-used-percent", "42")
	resp.Header.Set("x-codex-secondary-reset-after-seconds", "18000")
	resp.Header.Set("x-codex-secondary-window-minutes", "300")

	repo := &openAIAccountTestRepo{placeholder
	upstream := &queuedHTTPUpstream{responses: []*http.Response{respplaceholderplaceholder
	svc := &AccountTestService{accountRepo: repo, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          89,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder"access_token": "test-token"placeholder,
placeholder

	err := svc.testOpenAIAccountConnection(ctx, account, "gpt-5.4", "", "")
placeholder
	require.Len(t, upstream.requests, 1)
	require.Equal(t, HTTPUpstreamProfileOpenAI, HTTPUpstreamProfileFromContext(upstream.requests[0].Context()))
	require.NotEmpty(t, repo.updatedExtra)
	require.Equal(t, 42.0, repo.updatedExtra["codex_5h_used_percent"])
	require.Equal(t, 88.0, repo.updatedExtra["codex_7d_used_percent"])
	require.Contains(t, recorder.Body.String(), "test_complete")
placeholder

func TestAccountTestService_OpenAIOAuthTestNormalizesGPT56Alias(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := newTestContext()

	resp := newJSONResponse(http.StatusOK, "")
	resp.Body = io.NopCloser(strings.NewReader(`data: {"type":"response.completed"placeholder

`))

	upstream := &queuedHTTPUpstream{responses: []*http.Response{respplaceholderplaceholder
	svc := &AccountTestService{httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          90,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder"access_token": "test-token"placeholder,
placeholder

	err := svc.testOpenAIAccountConnection(ctx, account, "gpt-5.6", "", "")
placeholder
	require.Len(t, upstream.requests, 1)

	body, err := io.ReadAll(upstream.requests[0].Body)
placeholder
	require.Equal(t, "gpt-5.6-sol", gjson.GetBytes(body, "model").String())
placeholder

func TestAccountTestService_OpenAIShadowUsesParentCredentialsAndShadowModel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, recorder := newTestContext()

	resp := newJSONResponse(http.StatusOK, "")
	resp.Body = io.NopCloser(strings.NewReader(`data: {"type":"response.completed"placeholder

`))

	parentID := int64(100)
	parent := &Account{
		ID:       parentID,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
		Status:   StatusActive,
placeholder
			"access_token":       "parent-token",
			"chatgpt_account_id": "org-parent",
	placeholder,
placeholder
	shadow := &Account{
		ID:              200,
		Platform:        PlatformOpenAI,
		Type:            AccountTypeOAuth,
		Status:          StatusActive,
		ParentAccountID: &parentID,
		QuotaDimension:  QuotaDimensionSpark,
		Concurrency:     2,
placeholder
			"model_mapping": map[string]any{
				"gpt-5.3-codex-spark": "gpt-5.3-codex-spark",
		placeholder,
	placeholder,
placeholder

	repo := &openAIAccountTestRepo{
		mockAccountRepoForGemini: mockAccountRepoForGemini{
			accountsByID: map[int64]*Account{
				parentID: parent,
				200:      shadow,
		placeholder,
	placeholder,
placeholder
	upstream := &queuedHTTPUpstream{responses: []*http.Response{respplaceholderplaceholder
	svc := &AccountTestService{accountRepo: repo, httpUpstream: upstreamplaceholder

	err := svc.TestAccountConnection(ctx, shadow.ID, "gpt-5.3-codex-spark", "", "")
placeholder
	require.Len(t, upstream.requests, 1)
	req := upstream.requests[0]
	require.Equal(t, "Bearer parent-token", req.Header.Get("Authorization"))
	require.Equal(t, "org-parent", req.Header.Get("chatgpt-account-id"))
	body, err := io.ReadAll(req.Body)
placeholder
	require.Equal(t, "gpt-5.3-codex-spark", gjson.GetBytes(body, "model").String())
	require.Contains(t, recorder.Body.String(), `"success":true`)
placeholder

func TestAccountTestService_OpenAIStreamEOFBeforeCompletedFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, recorder := newTestContext()

	resp := newJSONResponse(http.StatusOK, "")
	resp.Body = io.NopCloser(strings.NewReader(`data: {"type":"response.output_text.delta","delta":"hi"placeholder

`))

	upstream := &queuedHTTPUpstream{responses: []*http.Response{respplaceholderplaceholder
	svc := &AccountTestService{httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          90,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder"access_token": "test-token"placeholder,
placeholder

	err := svc.testOpenAIAccountConnection(ctx, account, "gpt-5.4", "", "")
placeholder
	require.Contains(t, recorder.Body.String(), "response.completed")
	require.NotContains(t, recorder.Body.String(), `"success":true`)
placeholder

func TestAccountTestService_DeepSeekCustomBaseURLUsesV1ResponsesPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := newTestContext()

	resp := newJSONResponse(http.StatusOK, "")
	resp.Body = io.NopCloser(strings.NewReader(`data: {"type":"response.completed"placeholder

`))
	upstream := &queuedHTTPUpstream{responses: []*http.Response{respplaceholderplaceholder
	svc := &AccountTestService{
		httpUpstream: upstream,
		cfg:          &config.Config{Security: config.SecurityConfig{URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholderplaceholderplaceholder,
placeholder
	account := &Account{
		ID:          91,
		Platform:    PlatformDeepseek,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":      "sk-test",
			"base_url":     "https://relay.example.com/v1",
			"api_protocol": APIProtocolResponses,
	placeholder,
		Extra: map[string]any{
			openai_compat.ExtraKeyResponsesSupported: true,
	placeholder,
placeholder

	err := svc.testOpenAIAccountConnection(ctx, account, "gpt-5.4", "", "")
placeholder
	require.Len(t, upstream.requests, 1)
	require.Equal(t, "https://relay.example.com/v1/responses", upstream.requests[0].URL.String())
placeholder

func TestAccountTestService_DeepSeekResponsesRoutesToOpenAIProbe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := newTestContext()

	resp := newJSONResponse(http.StatusOK, "")
	resp.Body = io.NopCloser(strings.NewReader(`data: {"type":"response.completed"placeholder

`))
	upstream := &queuedHTTPUpstream{responses: []*http.Response{respplaceholderplaceholder
	svc := &AccountTestService{
		httpUpstream: upstream,
		cfg:          &config.Config{Security: config.SecurityConfig{URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholderplaceholderplaceholder,
placeholder
	account := &Account{
		ID:          93,
		Platform:    PlatformDeepseek,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":      "sk-test",
			"base_url":     "https://relay.example.com/v1",
			"api_protocol": APIProtocolResponses,
	placeholder,
		Extra: map[string]any{
			openai_compat.ExtraKeyResponsesSupported: true,
	placeholder,
placeholder
	repo := &openAIAccountTestRepo{
		mockAccountRepoForGemini: mockAccountRepoForGemini{
			accountsByID: map[int64]*Account{93: accountplaceholder,
	placeholder,
placeholder
	svc.accountRepo = repo

	err := svc.TestAccountConnection(ctx, account.ID, "gpt-5.4", "", "")
placeholder
	require.Len(t, upstream.requests, 1)
	require.Equal(t, "https://relay.example.com/v1/responses", upstream.requests[0].URL.String())
placeholder

func TestAccountTestService_DeepSeekDefaultBaseURLUsesNativeResponsesPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := newTestContext()

	resp := newJSONResponse(http.StatusOK, "")
	resp.Body = io.NopCloser(strings.NewReader(`data: {"type":"response.completed"placeholder

`))
	upstream := &queuedHTTPUpstream{responses: []*http.Response{respplaceholderplaceholder
	svc := &AccountTestService{
		httpUpstream: upstream,
		cfg:          &config.Config{Security: config.SecurityConfig{URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholderplaceholderplaceholder,
placeholder
	account := &Account{
		ID:          92,
		Platform:    PlatformDeepseek,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":      "sk-test",
			"api_protocol": APIProtocolResponses,
	placeholder,
		Extra: map[string]any{
			openai_compat.ExtraKeyResponsesSupported: true,
	placeholder,
placeholder

	err := svc.testOpenAIAccountConnection(ctx, account, "gpt-5.4", "", "")
placeholder
	require.Len(t, upstream.requests, 1)
	require.Equal(t, "https://api.deepseek.com/responses", upstream.requests[0].URL.String())
placeholder

func TestAccountTestService_OpenAI429PersistsSnapshotAndRateLimitState(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := newTestContext()

	resp := newJSONResponse(http.StatusTooManyRequests, `{"error":{"type":"usage_limit_reached","message":"limit reached","resets_at":1777283883placeholderplaceholder`)
	resp.Header.Set("x-codex-primary-used-percent", "100")
	resp.Header.Set("x-codex-primary-reset-after-seconds", "604800")
	resp.Header.Set("x-codex-primary-window-minutes", "10080")
	resp.Header.Set("x-codex-secondary-used-percent", "100")
	resp.Header.Set("x-codex-secondary-reset-after-seconds", "18000")
	resp.Header.Set("x-codex-secondary-window-minutes", "300")

	repo := &openAIAccountTestRepo{placeholder
	upstream := &queuedHTTPUpstream{responses: []*http.Response{respplaceholderplaceholder
	svc := &AccountTestService{accountRepo: repo, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          88,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Status:      StatusError,
		Concurrency: 1,
placeholder"access_token": "test-token"placeholder,
placeholder

	err := svc.testOpenAIAccountConnection(ctx, account, "gpt-5.4", "", "")
placeholder
	require.NotEmpty(t, repo.updatedExtra)
	require.Equal(t, 100.0, repo.updatedExtra["codex_5h_used_percent"])
	require.Equal(t, account.ID, repo.rateLimitedID)
	require.NotNil(t, repo.rateLimitedAt)
	require.Equal(t, account.ID, repo.clearedErrorID)
	require.Equal(t, StatusActive, account.Status)
	require.Empty(t, account.ErrorMessage)
	require.NotNil(t, account.RateLimitResetAt)
placeholder

func TestAccountTestService_OpenAI429BodyOnlyPersistsRateLimitAndClearsStaleError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := newTestContext()

	resp := newJSONResponse(http.StatusTooManyRequests, `{"error":{"type":"usage_limit_reached","message":"limit reached","resets_at":"1777283883"placeholderplaceholder`)

	repo := &openAIAccountTestRepo{placeholder
	upstream := &queuedHTTPUpstream{responses: []*http.Response{respplaceholderplaceholder
	svc := &AccountTestService{accountRepo: repo, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:           77,
		Platform:     PlatformOpenAI,
		Type:         AccountTypeOAuth,
		Status:       StatusError,
		ErrorMessage: "Access forbidden (403): account may be suspended or lack permissions",
		Concurrency:  1,
		Credentials:  map[string]any{"access_token": "test-token"placeholder,
placeholder

	err := svc.testOpenAIAccountConnection(ctx, account, "gpt-5.4", "", "")
placeholder
	require.Equal(t, account.ID, repo.rateLimitedID)
	require.NotNil(t, repo.rateLimitedAt)
	require.Equal(t, account.ID, repo.clearedErrorID)
	require.Equal(t, StatusActive, account.Status)
	require.Empty(t, account.ErrorMessage)
	require.NotNil(t, account.RateLimitResetAt)
	require.Empty(t, repo.updatedExtra)
placeholder

func TestAccountTestService_OpenAI429SyncsObservedPlanType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := newTestContext()

	resp := newJSONResponse(http.StatusTooManyRequests, `{"error":{"type":"usage_limit_reached","message":"limit reached","plan_type":"free","resets_at":1777283883placeholderplaceholder`)

	repo := &openAIAccountTestRepo{placeholder
	upstream := &queuedHTTPUpstream{responses: []*http.Response{respplaceholderplaceholder
	svc := &AccountTestService{accountRepo: repo, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          81,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
placeholder"access_token": "test-token", "plan_type": "plus"placeholder,
placeholder

	err := svc.testOpenAIAccountConnection(ctx, account, "gpt-5.4", "", "")
placeholder
	require.Equal(t, []int64{account.IDplaceholder, repo.bulkUpdatedIDs)
	require.Equal(t, "free", repo.bulkUpdatedPayload.Credentials["plan_type"])
	require.Equal(t, "free", account.Credentials["plan_type"])
	require.Equal(t, account.ID, repo.rateLimitedID)
	require.NotNil(t, account.RateLimitResetAt)
placeholder

func TestAccountTestService_OpenAI429ActiveAccountDoesNotClearError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := newTestContext()

	resp := newJSONResponse(http.StatusTooManyRequests, `{"error":{"type":"usage_limit_reached","message":"limit reached","resets_in_seconds":3600placeholderplaceholder`)

	repo := &openAIAccountTestRepo{placeholder
	upstream := &queuedHTTPUpstream{responses: []*http.Response{respplaceholderplaceholder
	svc := &AccountTestService{accountRepo: repo, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          78,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
placeholder"access_token": "test-token"placeholder,
placeholder

	err := svc.testOpenAIAccountConnection(ctx, account, "gpt-5.4", "", "")
placeholder
	require.Equal(t, account.ID, repo.rateLimitedID)
	require.NotNil(t, repo.rateLimitedAt)
	require.Zero(t, repo.clearedErrorID)
	require.Equal(t, StatusActive, account.Status)
	require.NotNil(t, account.RateLimitResetAt)
placeholder

func TestAccountTestService_OpenAI429WithoutResetSignalDoesNotMutateRuntimeState(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := newTestContext()

	resp := newJSONResponse(http.StatusTooManyRequests, `{"error":{"type":"usage_limit_reached","message":"limit reached"placeholderplaceholder`)

	repo := &openAIAccountTestRepo{placeholder
	upstream := &queuedHTTPUpstream{responses: []*http.Response{respplaceholderplaceholder
	svc := &AccountTestService{accountRepo: repo, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:           79,
		Platform:     PlatformOpenAI,
		Type:         AccountTypeOAuth,
		Status:       StatusError,
		ErrorMessage: "stale 403",
		Concurrency:  1,
		Credentials:  map[string]any{"access_token": "test-token"placeholder,
placeholder

	err := svc.testOpenAIAccountConnection(ctx, account, "gpt-5.4", "", "")
placeholder
	require.Zero(t, repo.rateLimitedID)
	require.Nil(t, repo.rateLimitedAt)
	require.Zero(t, repo.clearedErrorID)
	require.Equal(t, StatusError, account.Status)
	require.Equal(t, "stale 403", account.ErrorMessage)
	require.Nil(t, account.RateLimitResetAt)
placeholder

func TestAccountTestService_OpenAI401SetsPermanentErrorOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := newTestContext()

	resp := newJSONResponse(http.StatusUnauthorized, `{"error":"bad token"placeholder`)

	repo := &openAIAccountTestRepo{placeholder
	upstream := &queuedHTTPUpstream{responses: []*http.Response{respplaceholderplaceholder
	svc := &AccountTestService{accountRepo: repo, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          80,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
placeholder"access_token": "test-token"placeholder,
placeholder

	err := svc.testOpenAIAccountConnection(ctx, account, "gpt-5.4", "", "")
placeholder
	require.Equal(t, account.ID, repo.setErrorID)
	require.Contains(t, repo.setErrorMsg, "Authentication failed (401)")
	require.Zero(t, repo.rateLimitedID)
	require.Zero(t, repo.clearedErrorID)
	require.Nil(t, account.RateLimitResetAt)
placeholder

func TestAccountTestService_OpenAIAPIKeyResponsesUsesCodexProbeHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := newTestContext()

	resp := newJSONResponse(http.StatusOK, "")
	resp.Body = io.NopCloser(strings.NewReader("data: {\"type\":\"response.completed\"placeholder\n\n"))
	upstream := &queuedHTTPUpstream{responses: []*http.Response{respplaceholderplaceholder
	svc := &AccountTestService{
		httpUpstream: upstream,
		cfg:          &config.Config{Security: config.SecurityConfig{URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholderplaceholderplaceholder,
placeholder
	account := &Account{
		ID:          95,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "https://compat-upstream.example/v1",
	placeholder,
		Extra: map[string]any{openai_compat.ExtraKeyResponsesSupported: trueplaceholder,
placeholder

	err := svc.testOpenAIAccountConnection(ctx, account, "gpt-5.4", "", "")
placeholder
	require.Len(t, upstream.requests, 1)
	req := upstream.requests[0]
	require.Equal(t, "https://compat-upstream.example/v1/responses", req.URL.String())
	requireOpenAICodexProbeHeaders(t, req.Header)
placeholder

func TestAccountTestService_OpenAIAPIKeyResponsesUnsupportedUsesChatCompletionsPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, recorder := newTestContext()

	upstreamBody := strings.Join([]string{
		`data: {"id":"chatcmpl_test","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":"pong"placeholder,"finish_reason":nullplaceholder]placeholder`,
		"",
		`data: {"id":"chatcmpl_test","object":"chat.completion.chunk","choices":[{"index":0,"delta":{placeholder,"finish_reason":"stop"placeholder]placeholder`,
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamBody)),
placeholderplaceholder
	svc := &AccountTestService{
		httpUpstream: upstream,
		cfg:          &config.Config{Security: config.SecurityConfig{URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholderplaceholderplaceholder,
placeholder
	account := &Account{
		ID:          91,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "https://compat-upstream.example/v1",
	placeholder,
		Extra: map[string]any{openai_compat.ExtraKeyResponsesSupported: falseplaceholder,
placeholder

	err := svc.testOpenAIAccountConnection(ctx, account, "gpt-5.4", "hello", "")
placeholder
	require.NotNil(t, upstream.lastReq)
	require.Equal(t, HTTPUpstreamProfileOpenAI, HTTPUpstreamProfileFromContext(upstream.lastReq.Context()))
	require.Equal(t, "https://compat-upstream.example/v1/chat/completions", upstream.lastReq.URL.String())
	require.Equal(t, "Bearer sk-test", upstream.lastReq.Header.Get("Authorization"))
	require.Equal(t, "text/event-stream", upstream.lastReq.Header.Get("Accept"))
	require.Equal(t, "gpt-5.4", gjson.GetBytes(upstream.lastBody, "model").String())
	require.True(t, gjson.GetBytes(upstream.lastBody, "stream").Bool())
	require.Equal(t, "hello", gjson.GetBytes(upstream.lastBody, "messages.0.content").String())
	require.False(t, gjson.GetBytes(upstream.lastBody, "input").Exists())
	body := recorder.Body.String()
	require.Contains(t, body, "pong")
	require.Contains(t, body, "已通过 /v1/chat/completions 验证")
	require.Contains(t, body, `"success":true`)
	require.NotContains(t, body, "当前测试接口仅支持 Responses API 路径")
placeholder

func TestAccountTestService_OpenAIChatCompletionsPathReturns4xx(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, recorder := newTestContext()

	upstream := &httpUpstreamRecorder{resp: newJSONResponse(http.StatusBadRequest, `{"error":{"message":"bad request"placeholderplaceholder`)placeholder
	svc := &AccountTestService{
		httpUpstream: upstream,
		cfg:          &config.Config{Security: config.SecurityConfig{URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholderplaceholderplaceholder,
placeholder
	account := &Account{
		ID:          92,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "https://compat-upstream.example",
	placeholder,
		Extra: map[string]any{openai_compat.ExtraKeyResponsesSupported: falseplaceholder,
placeholder

	err := svc.testOpenAIAccountConnection(ctx, account, "gpt-5.4", "", "")
placeholder
	require.Equal(t, "https://compat-upstream.example/v1/chat/completions", upstream.lastReq.URL.String())
	require.Contains(t, err.Error(), "Chat Completions API (/v1/chat/completions) returned 400")
	require.Contains(t, recorder.Body.String(), "/v1/chat/completions")
	require.NotContains(t, recorder.Body.String(), `"success":true`)
placeholder

func TestAccountTestService_OpenAIChatCompletionsPathTimeout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, recorder := newTestContext()

	upstream := &httpUpstreamRecorder{err: context.DeadlineExceededplaceholder
	svc := &AccountTestService{
		httpUpstream: upstream,
		cfg:          &config.Config{Security: config.SecurityConfig{URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholderplaceholderplaceholder,
placeholder
	account := &Account{
		ID:          93,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "https://compat-upstream.example",
	placeholder,
		Extra: map[string]any{openai_compat.ExtraKeyResponsesSupported: falseplaceholder,
placeholder

	err := svc.testOpenAIAccountConnection(ctx, account, "gpt-5.4", "", "")
placeholder
	require.Equal(t, "https://compat-upstream.example/v1/chat/completions", upstream.lastReq.URL.String())
	require.Contains(t, err.Error(), "Chat Completions API (/v1/chat/completions) request failed")
	require.Contains(t, err.Error(), context.DeadlineExceeded.Error())
	require.Contains(t, recorder.Body.String(), "/v1/chat/completions")
	require.NotContains(t, recorder.Body.String(), `"success":true`)
placeholder

func TestAccountTestService_OpenAIChatCompletionsPathRejectsNonJSONStream(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, recorder := newTestContext()

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader("data: not-json\n\n")),
placeholderplaceholder
	svc := &AccountTestService{
		httpUpstream: upstream,
		cfg:          &config.Config{Security: config.SecurityConfig{URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholderplaceholderplaceholder,
placeholder
	account := &Account{
		ID:          94,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "https://compat-upstream.example",
	placeholder,
		Extra: map[string]any{openai_compat.ExtraKeyResponsesSupported: falseplaceholder,
placeholder

	err := svc.testOpenAIAccountConnection(ctx, account, "gpt-5.4", "", "")
placeholder
	require.Equal(t, "https://compat-upstream.example/v1/chat/completions", upstream.lastReq.URL.String())
	require.Contains(t, err.Error(), "Invalid Chat Completions response from /v1/chat/completions")
	require.Contains(t, recorder.Body.String(), "/v1/chat/completions")
	require.NotContains(t, recorder.Body.String(), `"success":true`)
placeholder
