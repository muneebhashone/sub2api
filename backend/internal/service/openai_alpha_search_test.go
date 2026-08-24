package service

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

type alphaSearchAccountStateRepo struct {
	AccountRepository
	setErrorCalls int
	lastError     string
placeholder

func (r *alphaSearchAccountStateRepo) SetError(_ context.Context, _ int64, errorMsg string) error {
	r.setErrorCalls++
	r.lastError = errorMsg
	return nil
placeholder

func alphaSearchResponsesSSE(output string) string {
	return "event: response.output_text.delta\n" +
		`data: {"type":"response.output_text.delta","delta":` + strconv.Quote(output) + `placeholder` + "\n\n" +
		"event: response.output_text.annotation.added\n" +
		`data: {"type":"response.output_text.annotation.added","annotation":{"type":"url_citation","url":"https://example.com/news","title":"Example News"placeholderplaceholder` + "\n\n" +
		"event: response.completed\n" +
		`data: {"type":"response.completed","response":{"output":[{"type":"message","content":[{"type":"output_text","text":` + strconv.Quote(output) + `placeholder]placeholder]placeholderplaceholder` + "\n\n"
placeholder

func TestForwardAlphaSearchOAuthPreservesWire(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{
		"id":"search-session",
		"model":"gpt-5.6-sol",
		"reasoning":{"effort":"max","context":"all_turns"placeholder,
		"input":[{"type":"message","role":"user","content":[{"type":"input_text","text":"latest news"placeholder]placeholder],
		"commands":{"search_query":[{"q":"OpenAI news","recency":1placeholder]placeholder,
		"settings":{"allowed_callers":["direct"],"external_web_access":trueplaceholder,
		"max_output_tokens":2000,
		"future_field":{"keep":trueplaceholder
placeholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/alpha/search?feature=standalone", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("User-Agent", codexCLIUserAgent)
	c.Request.Header.Set("Originator", "codex_cli_rs")
	c.Request.Header.Set("Version", "0.144.1")
	c.Request.Header.Set("X-Codex-Turn-Metadata", `{"session_id":"search-session","turn_id":"search-turn"placeholder`)

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"encrypted_output":"ciphertext","output":"search result"placeholder`)),
placeholderplaceholder
	service := &OpenAIGatewayService{cfg: &config.Config{placeholder, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          42,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token":       "oauth-token",
			"chatgpt_account_id": "chatgpt-account",
	placeholder,
placeholder

	result, err := service.ForwardAlphaSearch(context.Background(), c, account, body)

placeholder
	require.NotNil(t, result)
	require.Equal(t, 1, result.WebSearchCalls)
	require.Equal(t, "gpt-5.6-sol", result.Model)
	require.Equal(t, http.StatusOK, recorder.Code)
	require.JSONEq(t, `{"encrypted_output":"ciphertext","output":"search result"placeholder`, recorder.Body.String())
	require.Equal(t, chatgptCodexAlphaSearchURL+"?feature=standalone", upstream.lastReq.URL.String())
	require.Equal(t, "chatgpt.com", upstream.lastReq.Host)
	require.Equal(t, "Bearer oauth-token", upstream.lastReq.Header.Get("Authorization"))
	require.Equal(t, "chatgpt-account", upstream.lastReq.Header.Get("chatgpt-account-id"))
	require.Equal(t, "application/json", upstream.lastReq.Header.Get("Accept"))
	require.Equal(t, codexCLIVersion, upstream.lastReq.Header.Get("Version"))
	require.Empty(t, upstream.lastReq.Header.Get("OpenAI-Beta"))
	require.Equal(t,
		scopeCodexAccountIdentityValue(account, 0, "session", "search-session"),
		gjson.Get(upstream.lastReq.Header.Get("X-Codex-Turn-Metadata"), "session_id").String(),
	)
	require.Equal(t,
		scopeCodexAccountIdentityValue(account, 0, "turn", "search-turn"),
		gjson.Get(upstream.lastReq.Header.Get("X-Codex-Turn-Metadata"), "turn_id").String(),
	)
	require.JSONEq(t, string(body), string(upstream.lastBody))
placeholder

func TestForwardAlphaSearchPATUsesResponsesWebSearchFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{
		"id":"search-session",
		"model":"gpt-5.6-sol",
		"commands":{"search_query":[{"q":"OpenAI news"placeholder]placeholder,
		"prompt_cache_key":"responses-cache-key",
		"prompt_cache_retention":"24h"
placeholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/alpha/search", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("User-Agent", codexCLIUserAgent)
	c.Request.Header.Set("Originator", "codex_cli_rs")
	c.Request.Header.Set("Version", "0.144.1")
	c.Request.Header.Set("OpenAI-Beta", "responses=experimental")
	c.Request.Header.Set("Accept-Language", "zh-CN")
	c.Request.Header.Set("Authorization", "Bearer client-token")
	c.Request.Header.Set("Session_ID", "session-client")
	c.Request.Header.Set("Conversation_ID", "conversation-client")
	c.Request.Header.Set("X-Codex-Beta-Features", "feature-a")
	c.Request.Header.Set("X-Codex-Turn-State", "turn-state")
	c.Request.Header.Set(responsesLiteHeaderKey, "true")
	c.Request.Header.Set("X-Codex-Turn-Metadata", `{"turn_id":"turn-1"placeholder`)

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"req-search"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(alphaSearchResponsesSSE("search result"))),
placeholderplaceholder
	service := &OpenAIGatewayService{cfg: &config.Config{placeholder, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          43,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token":               "at-test-token",
			"auth_mode":                  OpenAIAuthModePersonalAccessToken,
			"chatgpt_account_id":         "chatgpt-account",
			"chatgpt_account_is_fedramp": true,
	placeholder,
placeholder

	result, err := service.ForwardAlphaSearch(context.Background(), c, account, body)

placeholder
	require.NotNil(t, result)
	require.Equal(t, 1, result.WebSearchCalls)
	require.Equal(t, "/v1/responses", result.UpstreamEndpoint)
	require.Equal(t, http.StatusOK, recorder.Code)
	require.JSONEq(t, `{"output":"search result","results":[{"type":"text_result","ref_id":"turn0search0","url":"https://example.com/news","title":"Example News"placeholder]placeholder`, recorder.Body.String())
	require.Equal(t, chatgptCodexURL, upstream.lastReq.URL.String())
	require.Equal(t, "Bearer at-test-token", upstream.lastReq.Header.Get("Authorization"))
	require.Equal(t, "chatgpt-account", upstream.lastReq.Header.Get("ChatGPT-Account-ID"))
	require.Equal(t, "true", upstream.lastReq.Header.Get("X-OpenAI-Fedramp"))
	require.Equal(t, "application/json", upstream.lastReq.Header.Get("Content-Type"))
	require.Equal(t, "text/event-stream", upstream.lastReq.Header.Get("Accept"))
	require.Equal(t, "responses=experimental", upstream.lastReq.Header.Get("OpenAI-Beta"))
	require.Equal(t, codexCLIVersion, upstream.lastReq.Header.Get("Version"))
	require.Equal(t,
		scopeCodexAccountIdentityValue(account, 0, "turn", "turn-1"),
		gjson.Get(upstream.lastReq.Header.Get("X-Codex-Turn-Metadata"), "turn_id").String(),
	)
	require.Equal(t, openai.CodexDefaultOriginator, upstream.lastReq.Header.Get("Originator"))
	require.Empty(t, upstream.lastReq.Header.Get("X-Codex-Beta-Features"))
	require.Empty(t, upstream.lastReq.Header.Get("X-Codex-Turn-State"))
	require.Empty(t, upstream.lastReq.Header.Get(responsesLiteHeaderKey))
	require.Empty(t, upstream.lastReq.Header.Get("Accept-Language"))
	require.False(t, gjson.GetBytes(upstream.lastBody, "prompt_cache_key").Exists())
	require.False(t, gjson.GetBytes(upstream.lastBody, "prompt_cache_retention").Exists())
	require.Equal(t, "gpt-5.6-sol", gjson.GetBytes(upstream.lastBody, "model").String())
	require.True(t, gjson.GetBytes(upstream.lastBody, "stream").Bool())
	require.False(t, gjson.GetBytes(upstream.lastBody, "store").Bool())
	require.Equal(t, "web_search", gjson.GetBytes(upstream.lastBody, "tools.0.type").String())
	require.Contains(t, gjson.GetBytes(upstream.lastBody, "input.0.content.0.text").String(), `"search_query"`)
placeholder

func TestForwardAlphaSearchPATBackfillsMissingChatGPTAccountMetadata(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"id":"search-session","model":"gpt-5.6-sol","commands":{"search_query":[{"q":"OpenAI news"placeholder]placeholderplaceholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/alpha/search", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	var whoamiCalls int32
	whoamiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&whoamiCalls, 1)
		require.Equal(t, "Bearer at-test-token", r.Header.Get("Authorization"))
		require.Equal(t, "application/json", r.Header.Get("Accept"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"email":"pat@example.com",
			"chatgpt_user_id":"user-123",
			"chatgpt_account_id":"acct-123",
			"chatgpt_plan_type":"plus",
			"chatgpt_account_is_fedramp":true
	placeholder`))
placeholder))
	defer whoamiServer.Close()
	oldWhoamiURL := openAICodexPATWhoamiURL
	openAICodexPATWhoamiURL = whoamiServer.URL
	defer func() { openAICodexPATWhoamiURL = oldWhoamiURL placeholder()

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"output":"search result"placeholder`)),
placeholderplaceholder
	oauthService := NewOpenAIOAuthService(nil, nil)
	service := &OpenAIGatewayService{
		cfg:                 &config.Config{placeholder,
		httpUpstream:        upstream,
		openAITokenProvider: NewOpenAITokenProvider(nil, nil, oauthService),
placeholder
	account := &Account{
		ID:          45,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token": "at-test-token",
			"auth_mode":    OpenAIAuthModePersonalAccessToken,
	placeholder,
placeholder

	result, err := service.ForwardAlphaSearch(context.Background(), c, account, body)

placeholder
	require.NotNil(t, result)
	require.Equal(t, int32(1), atomic.LoadInt32(&whoamiCalls))
	require.Equal(t, "acct-123", upstream.lastReq.Header.Get("ChatGPT-Account-ID"))
	require.Equal(t, "true", upstream.lastReq.Header.Get("X-OpenAI-Fedramp"))
	require.Equal(t, "acct-123", account.Credentials["chatgpt_account_id"])
	require.Equal(t, "user-123", account.Credentials["chatgpt_user_id"])
	require.Equal(t, OpenAIAuthModePersonalAccessToken, account.Credentials["auth_mode"])
placeholder

func TestForwardAlphaSearchAPIKeyMapsModelAndPassesThroughError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"id":"search-session","model":"gpt-5.6-sol","commands":{"search_query":[{"q":"news"placeholder]placeholderplaceholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/alpha/search", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstreamBody := `{"error":{"type":"invalid_request_error","message":"bad search"placeholderplaceholder`
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamBody)),
placeholderplaceholder
	service := &OpenAIGatewayService{cfg: &config.Config{placeholder, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:       7,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
placeholder
			"api_key":  "sk-test",
			"base_url": "https://compat.example/v4",
			"model_mapping": map[string]any{
				"gpt-5.6-sol": "upstream-5.6",
		placeholder,
	placeholder,
placeholder

	result, err := service.ForwardAlphaSearch(context.Background(), c, account, body)

placeholder
	// 上游错误透传不是一次成功的搜索：不返回 result、不产生按次计费。
	require.Nil(t, result)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.JSONEq(t, upstreamBody, recorder.Body.String())
	require.Equal(t, "https://compat.example/v4/alpha/search", upstream.lastReq.URL.String())
	require.Equal(t, "Bearer sk-test", upstream.lastReq.Header.Get("Authorization"))
	require.Equal(t, "upstream-5.6", gjson.GetBytes(upstream.lastBody, "model").String())
	require.True(t, gjson.GetBytes(upstream.lastBody, "commands.search_query").IsArray())
placeholder

func TestForwardAlphaSearchReturnsFailoverBeforeWriting(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"id":"search-session","model":"gpt-5.6-sol","commands":{placeholderplaceholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/alpha/search", bytes.NewReader(body))

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"rate limited"placeholderplaceholder`)),
placeholderplaceholder
	service := &OpenAIGatewayService{cfg: &config.Config{placeholder, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:       8,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
placeholder
			"api_key": "sk-test",
	placeholder,
placeholder

	result, err := service.ForwardAlphaSearch(context.Background(), c, account, body)

	require.Nil(t, result)
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusTooManyRequests, failoverErr.StatusCode)
	require.Empty(t, failoverErr.Stage)
	require.Empty(t, failoverErr.Scope)
	require.Empty(t, failoverErr.Reason)
	require.Zero(t, failoverErr.ClientStatusCode)
	require.Empty(t, failoverErr.ClientMessage)
	require.Nil(t, failoverErr.ResponseHeaders)
	require.Equal(t, openAIPlatformAlphaSearchURL, upstream.lastReq.URL.String())
	require.False(t, c.Writer.Written())
	require.Empty(t, recorder.Body.String())
placeholder

func TestForwardAlphaSearchSetupToken429CarriesSameAccountRetryWindow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"id":"search-session","model":"gpt-5.6-sol","commands":{placeholderplaceholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/alpha/search", bytes.NewReader(body))

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header: http.Header{
			"Content-Type": []string{"application/json"placeholder,
			"Retry-After":  []string{"1"placeholder,
			"X-Request-Id": []string{"req_alpha_oauth_429"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(`{"error":{"type":"rate_limit_error","code":"rate_limit_exceeded","message":"rate limited"placeholderplaceholder`)),
placeholderplaceholder
	service := &OpenAIGatewayService{cfg: &config.Config{placeholder, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          81,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeSetupToken,
		Concurrency: 1,
placeholder
			"access_token":       "oauth-token",
			"chatgpt_account_id": "chatgpt-account",
	placeholder,
placeholder
	startedAt := time.Now()

	result, err := service.ForwardAlphaSearch(context.Background(), c, account, body)

	require.Nil(t, result)
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.True(t, failoverErr.RetryableOnSameAccount)
	require.Equal(t, time.Second, failoverErr.SameAccountRetryDelay)
	require.WithinDuration(t, startedAt.Add(openAIOAuth429RetryWindow), failoverErr.SameAccountRetryDeadline, time.Second)
	require.Equal(t, "req_alpha_oauth_429", failoverErr.ResponseHeaders.Get("x-request-id"))
	require.False(t, c.Writer.Written())
placeholder

func TestForwardAlphaSearchAccessStateUsesTypedFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"id":"search-session","model":"gpt-5.6-sol","commands":{placeholderplaceholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/alpha/search", bytes.NewReader(body))

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusForbidden,
		Header: http.Header{
			"Content-Type": []string{"application/json"placeholder,
			"X-Request-Id": []string{"req_alpha_access_state"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(`{"error":{"code":"deactivated_workspace","message":"Workspace is deactivated"placeholderplaceholder`)),
placeholderplaceholder
	service := &OpenAIGatewayService{cfg: &config.Config{placeholder, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          11,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token":       "oauth-token",
			"chatgpt_account_id": "chatgpt-account",
	placeholder,
placeholder

	result, err := service.ForwardAlphaSearch(context.Background(), c, account, body)

	require.Nil(t, result)
	assertOpenAIAlphaSearchAccessStateFailover(t, err, "req_alpha_access_state")
	require.Equal(t, chatgptCodexAlphaSearchURL, upstream.lastReq.URL.String())
	require.False(t, c.Writer.Written())
placeholder

func TestForwardAlphaSearchPATFallbackAccessStateUsesTypedFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"id":"search-session","model":"gpt-5.6-sol","commands":{"search_query":[{"q":"news"placeholder]placeholderplaceholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/alpha/search", bytes.NewReader(body))

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusForbidden,
		Header: http.Header{
			"Content-Type": []string{"application/json"placeholder,
			"X-Request-Id": []string{"req_alpha_pat_access_state"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(`{"error":{"code":"account_disabled","message":"Account is disabled"placeholderplaceholder`)),
placeholderplaceholder
	service := &OpenAIGatewayService{cfg: &config.Config{placeholder, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          12,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token":       "at-test-token",
			"auth_mode":          OpenAIAuthModePersonalAccessToken,
			"chatgpt_account_id": "chatgpt-account",
	placeholder,
placeholder

	result, err := service.ForwardAlphaSearch(context.Background(), c, account, body)

	require.Nil(t, result)
	assertOpenAIAlphaSearchAccessStateFailover(t, err, "req_alpha_pat_access_state")
	require.Equal(t, chatgptCodexURL, upstream.lastReq.URL.String())
	require.False(t, c.Writer.Written())
placeholder

func assertOpenAIAlphaSearchAccessStateFailover(t *testing.T, err error, requestID string) {
placeholder
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusForbidden, failoverErr.StatusCode)
	require.Equal(t, GatewayFailureStageAccountAuth, failoverErr.Stage)
	require.Equal(t, GatewayFailureScopeAccount, failoverErr.Scope)
	require.Equal(t, OpenAIUpstreamAccessStateReason, failoverErr.Reason)
	require.Equal(t, NextAccountRetry, failoverErr.NextAccountAction)
	require.Equal(t, http.StatusBadGateway, failoverErr.ClientStatusCode)
	require.Equal(t, openAIUpstreamAccessUnavailableClientMessage, failoverErr.ClientMessage)
	require.False(t, failoverErr.RetryableOnSameAccount)
	require.Equal(t, requestID, failoverErr.ResponseHeaders.Get("x-request-id"))
placeholder

func TestForwardAlphaSearchUnauthorizedDoesNotMarkAccountError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"id":"search-session","model":"gpt-5.6-sol","commands":{"search_query":[{"q":"news"placeholder]placeholderplaceholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/alpha/search", bytes.NewReader(body))

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusUnauthorized,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"detail":"Unauthorized"placeholder`)),
placeholderplaceholder
	repo := &alphaSearchAccountStateRepo{placeholder
	cfg := &config.Config{placeholder
	service := &OpenAIGatewayService{
		cfg:              cfg,
		httpUpstream:     upstream,
		accountRepo:      repo,
		rateLimitService: NewRateLimitService(repo, nil, cfg, nil, nil),
placeholder
	account := &Account{
		ID:          44,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			// 刻意不设置 auth_mode：覆盖历史上把 at- token 当普通 OAuth 导入的账号。
			"access_token":       "at-test-token",
			"chatgpt_account_id": "chatgpt-account",
	placeholder,
placeholder

	result, err := service.ForwardAlphaSearch(context.Background(), c, account, body)

	require.Nil(t, result)
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusUnauthorized, failoverErr.StatusCode)
	require.Zero(t, repo.setErrorCalls)
	require.Empty(t, repo.lastError)
	require.False(t, c.Writer.Written())
placeholder

func TestForwardAlphaSearchPATResponsesFallbackUnauthorizedDoesNotMarkAccountError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"id":"search-session","model":"gpt-5.6-sol","commands":{"search_query":[{"q":"news"placeholder]placeholderplaceholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/alpha/search", bytes.NewReader(body))

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusUnauthorized,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"detail":"Unauthorized"placeholder`)),
placeholderplaceholder
	repo := &alphaSearchAccountStateRepo{placeholder
	cfg := &config.Config{placeholder
	service := &OpenAIGatewayService{
		cfg:              cfg,
		httpUpstream:     upstream,
		accountRepo:      repo,
		rateLimitService: NewRateLimitService(repo, nil, cfg, nil, nil),
placeholder
	account := &Account{
		ID:          46,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token":       "at-test-token",
			"auth_mode":          OpenAIAuthModePersonalAccessToken,
			"chatgpt_account_id": "chatgpt-account",
	placeholder,
placeholder

	result, err := service.ForwardAlphaSearch(context.Background(), c, account, body)

	require.Nil(t, result)
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusUnauthorized, failoverErr.StatusCode)
	require.Equal(t, chatgptCodexURL, upstream.lastReq.URL.String())
	require.Equal(t, "text/event-stream", upstream.lastReq.Header.Get("Accept"))
	require.Equal(t, "responses=experimental", upstream.lastReq.Header.Get("OpenAI-Beta"))
	require.Zero(t, repo.setErrorCalls)
	require.Empty(t, repo.lastError)
	require.False(t, c.Writer.Written())
placeholder

// API key 上游（官方平台或第三方中转）不提供 /v1/alpha/search 时返回的
// 404/405 必须触发换号而不是把错误透传给客户端：混合分组里 OAuth 账号可以
// 承接搜索，请求不能死在先被选中的 API key 账号上。端点缺失也不能写账号
// 错误状态——账号本身是健康的。
func TestForwardAlphaSearchAPIKeyEndpointNotFoundFailsOver(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"id":"search-session","model":"gpt-5.6-sol","commands":{"search_query":[{"q":"news"placeholder]placeholderplaceholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/alpha/search", bytes.NewReader(body))

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusNotFound,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"Not Found"placeholderplaceholder`)),
placeholderplaceholder
	repo := &alphaSearchAccountStateRepo{placeholder
	cfg := &config.Config{placeholder
	service := &OpenAIGatewayService{
		cfg:              cfg,
		httpUpstream:     upstream,
		accountRepo:      repo,
		rateLimitService: NewRateLimitService(repo, nil, cfg, nil, nil),
placeholder
	account := &Account{
		ID:       9,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
placeholder
			"api_key":  "sk-test",
			"base_url": "https://relay.example",
	placeholder,
placeholder

	result, err := service.ForwardAlphaSearch(context.Background(), c, account, body)

	require.Nil(t, result)
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusNotFound, failoverErr.StatusCode)
	require.Zero(t, repo.setErrorCalls)
	require.Empty(t, repo.lastError)
	require.False(t, c.Writer.Written())
	require.Empty(t, recorder.Body.String())
placeholder

// OAuth 账号的 chatgpt.com 端点固定存在，404 保持原有透传行为不变。
func TestForwardAlphaSearchOAuthNotFoundPassesThrough(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"id":"search-session","model":"gpt-5.6-sol","commands":{"search_query":[{"q":"news"placeholder]placeholderplaceholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/alpha/search", bytes.NewReader(body))

	upstreamBody := `{"detail":"Not Found"placeholder`
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusNotFound,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamBody)),
placeholderplaceholder
	service := &OpenAIGatewayService{cfg: &config.Config{placeholder, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          10,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token":       "oauth-token",
			"chatgpt_account_id": "chatgpt-account",
	placeholder,
placeholder

	result, err := service.ForwardAlphaSearch(context.Background(), c, account, body)

placeholder
	require.Nil(t, result)
	require.Equal(t, http.StatusNotFound, recorder.Code)
	require.JSONEq(t, upstreamBody, recorder.Body.String())
placeholder

func TestShouldApplyOpenAIAlphaSearchAccountErrorSideEffects(t *testing.T) {
	require.False(t, shouldApplyOpenAIAlphaSearchAccountErrorSideEffects(http.StatusUnauthorized))
	require.False(t, shouldApplyOpenAIAlphaSearchAccountErrorSideEffects(http.StatusNotFound))
	require.False(t, shouldApplyOpenAIAlphaSearchAccountErrorSideEffects(http.StatusMethodNotAllowed))
	require.True(t, shouldApplyOpenAIAlphaSearchAccountErrorSideEffects(http.StatusForbidden))
	require.True(t, shouldApplyOpenAIAlphaSearchAccountErrorSideEffects(http.StatusTooManyRequests))
placeholder

func TestOpenAIAlphaSearchSchedulingModelUsesCanonicalAccountMapping(t *testing.T) {
	account := &Account{Credentials: map[string]any{
		"model_mapping": map[string]any{"client-visible": "canonical-upstream"placeholder,
placeholderplaceholder
	require.Equal(t, "canonical-upstream", openAIAlphaSearchSchedulingModel(account, "client-visible"))
	require.Equal(t, "unmapped", openAIAlphaSearchSchedulingModel(account, "unmapped"))
placeholder

func TestSanitizeOpenAIAlphaSearchBody_RemovesResponsesOnlyFields(t *testing.T) {
	body := []byte(`{"id":"search-session","store":false,"prompt_cache_key":"cache","commands":{"search_query":[{"q":"news"placeholder]placeholderplaceholder`)

	normalized, err := sanitizeOpenAIAlphaSearchBody(body)
placeholder
	require.False(t, gjson.GetBytes(normalized, "store").Exists())
	require.False(t, gjson.GetBytes(normalized, "prompt_cache_key").Exists())
	require.Equal(t, "news", gjson.GetBytes(normalized, "commands.search_query.0.q").String())
placeholder

func TestIsOpenAIAlphaSearchEndpointUnsupported(t *testing.T) {
	apiKey := &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder
	oauth := &Account{Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder

	require.True(t, isOpenAIAlphaSearchEndpointUnsupported(apiKey, http.StatusNotFound))
	require.True(t, isOpenAIAlphaSearchEndpointUnsupported(apiKey, http.StatusMethodNotAllowed))
	require.False(t, isOpenAIAlphaSearchEndpointUnsupported(apiKey, http.StatusBadRequest))
	require.False(t, isOpenAIAlphaSearchEndpointUnsupported(oauth, http.StatusNotFound))
	require.False(t, isOpenAIAlphaSearchEndpointUnsupported(nil, http.StatusNotFound))
placeholder
