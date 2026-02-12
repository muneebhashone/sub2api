package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func f64p(v float64) *float64 { return &v placeholder

type httpUpstreamRecorder struct {
	lastReq  *http.Request
	lastBody []byte

	resp *http.Response
	err  error
placeholder

func (u *httpUpstreamRecorder) Do(req *http.Request, proxyURL string, accountID int64, accountConcurrency int) (*http.Response, error) {
	u.lastReq = req
	if req != nil && req.Body != nil {
		b, _ := io.ReadAll(req.Body)
		u.lastBody = b
		_ = req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(b))
placeholder
	if u.err != nil {
		return nil, u.err
placeholder
	return u.resp, nil
placeholder

func (u *httpUpstreamRecorder) DoWithTLS(req *http.Request, proxyURL string, accountID int64, accountConcurrency int, enableTLSFingerprint bool) (*http.Response, error) {
	return u.Do(req, proxyURL, accountID, accountConcurrency)
placeholder

var structuredLogCaptureMu sync.Mutex

type inMemoryLogSink struct {
	mu     sync.Mutex
	events []*logger.LogEvent
placeholder

func (s *inMemoryLogSink) WriteLogEvent(event *logger.LogEvent) {
	if event == nil {
		return
placeholder
	cloned := *event
	if event.Fields != nil {
		cloned.Fields = make(map[string]any, len(event.Fields))
		for k, v := range event.Fields {
			cloned.Fields[k] = v
	placeholder
placeholder
	s.mu.Lock()
	s.events = append(s.events, &cloned)
	s.mu.Unlock()
placeholder

func (s *inMemoryLogSink) ContainsMessage(substr string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, ev := range s.events {
		if ev != nil && strings.Contains(ev.Message, substr) {
			return true
	placeholder
placeholder
	return false
placeholder

func (s *inMemoryLogSink) ContainsFieldValue(field, substr string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, ev := range s.events {
		if ev == nil || ev.Fields == nil {
			continue
	placeholder
		if v, ok := ev.Fields[field]; ok && strings.Contains(fmt.Sprint(v), substr) {
			return true
	placeholder
placeholder
	return false
placeholder

func captureStructuredLog(t *testing.T) (*inMemoryLogSink, func()) {
placeholder
	structuredLogCaptureMu.Lock()

	err := logger.Init(logger.InitOptions{
		Level:       "debug",
		Format:      "json",
		ServiceName: "sub2api",
		Environment: "test",
		Output: logger.OutputOptions{
			ToStdout: true,
			ToFile:   false,
	placeholder,
		Sampling: logger.SamplingOptions{Enabled: falseplaceholder,
placeholder)
placeholder

	sink := &inMemoryLogSink{placeholder
	logger.SetSink(sink)
	return sink, func() {
		logger.SetSink(nil)
		structuredLogCaptureMu.Unlock()
placeholder
placeholder

func TestOpenAIGatewayService_OAuthPassthrough_StreamKeepsToolNameAndBodyUnchanged(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(nil))
	c.Request.Header.Set("User-Agent", "codex_cli_rs/0.1.0")
	c.Request.Header.Set("Authorization", "Bearer inbound-should-not-forward")
	c.Request.Header.Set("Cookie", "secret=1")
	c.Request.Header.Set("X-Api-Key", "sk-inbound")
	c.Request.Header.Set("X-Goog-Api-Key", "goog-inbound")
	c.Request.Header.Set("Accept-Encoding", "gzip")
	c.Request.Header.Set("Proxy-Authorization", "Basic abc")
	c.Request.Header.Set("X-Test", "keep")

	originalBody := []byte(`{"model":"gpt-5.2","stream":true,"store":true,"input":[{"type":"text","text":"hi"placeholder]placeholder`)

	upstreamSSE := strings.Join([]string{
		`data: {"type":"response.output_item.added","item":{"type":"tool_call","tool_calls":[{"function":{"name":"apply_patch"placeholderplaceholder]placeholderplaceholder`,
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamSSE)),
placeholder
	upstream := &httpUpstreamRecorder{resp: respplaceholder

	svc := &OpenAIGatewayService{
		cfg:          &config.Config{Gateway: config.GatewayConfig{ForceCodexCLI: falseplaceholderplaceholder,
		httpUpstream: upstream,
		openAITokenProvider: &OpenAITokenProvider{ // minimal: will be bypassed by nil cache/service, but GetAccessToken uses provider only if non-nil
			accountRepo: nil,
	placeholder,
placeholder

	account := &Account{
		ID:             123,
		Name:           "acc",
		Platform:       PlatformOpenAI,
		Type:           AccountTypeOAuth,
		Concurrency:    1,
		Credentials:    map[string]any{"access_token": "oauth-token", "chatgpt_account_id": "chatgpt-acc"placeholder,
		Extra:          map[string]any{"openai_passthrough": trueplaceholder,
		Status:         StatusActive,
		Schedulable:    true,
		RateMultiplier: f64p(1),
placeholder

	// Use the gateway method that reads token from credentials when provider is nil.
	svc.openAITokenProvider = nil

	result, err := svc.Forward(context.Background(), c, account, originalBody)
placeholder
	require.NotNil(t, result)
	require.True(t, result.Stream)

	// 1) upstream body is exactly unchanged
	require.Equal(t, originalBody, upstream.lastBody)

	// 2) only auth is replaced; inbound auth/cookie are not forwarded
	require.Equal(t, "Bearer oauth-token", upstream.lastReq.Header.Get("Authorization"))
	require.Equal(t, "codex_cli_rs/0.1.0", upstream.lastReq.Header.Get("User-Agent"))
	require.Empty(t, upstream.lastReq.Header.Get("Cookie"))
	require.Empty(t, upstream.lastReq.Header.Get("X-Api-Key"))
	require.Empty(t, upstream.lastReq.Header.Get("X-Goog-Api-Key"))
	require.Empty(t, upstream.lastReq.Header.Get("Accept-Encoding"))
	require.Empty(t, upstream.lastReq.Header.Get("Proxy-Authorization"))
	require.Empty(t, upstream.lastReq.Header.Get("X-Test"))

	// 3) required OAuth headers are present
	require.Equal(t, "chatgpt.com", upstream.lastReq.Host)
	require.Equal(t, "chatgpt-acc", upstream.lastReq.Header.Get("chatgpt-account-id"))

	// 4) downstream SSE keeps tool name (no toolCorrector)
	body := rec.Body.String()
	require.Contains(t, body, "apply_patch")
	require.NotContains(t, body, "\"name\":\"edit\"")
placeholder

func TestOpenAIGatewayService_OAuthPassthrough_DisabledUsesLegacyTransform(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(nil))
	c.Request.Header.Set("User-Agent", "codex_cli_rs/0.1.0")

	// store=true + stream=false should be forced to store=false + stream=true by applyCodexOAuthTransform (OAuth legacy path)
	inputBody := []byte(`{"model":"gpt-5.2","stream":false,"store":true,"input":[{"type":"text","text":"hi"placeholder]placeholder`)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader("data: [DONE]\n\n")),
placeholder
	upstream := &httpUpstreamRecorder{resp: respplaceholder

	svc := &OpenAIGatewayService{
		cfg:          &config.Config{Gateway: config.GatewayConfig{ForceCodexCLI: falseplaceholderplaceholder,
		httpUpstream: upstream,
placeholder

	account := &Account{
		ID:             123,
		Name:           "acc",
		Platform:       PlatformOpenAI,
		Type:           AccountTypeOAuth,
		Concurrency:    1,
		Credentials:    map[string]any{"access_token": "oauth-token", "chatgpt_account_id": "chatgpt-acc"placeholder,
		Extra:          map[string]any{"openai_passthrough": falseplaceholder,
		Status:         StatusActive,
		Schedulable:    true,
		RateMultiplier: f64p(1),
placeholder

	_, err := svc.Forward(context.Background(), c, account, inputBody)
placeholder

	// legacy path rewrites request body (not byte-equal)
	require.NotEqual(t, inputBody, upstream.lastBody)
	require.Contains(t, string(upstream.lastBody), `"store":false`)
	require.Contains(t, string(upstream.lastBody), `"stream":true`)
placeholder

func TestOpenAIGatewayService_OAuthPassthrough_ResponseHeadersAllowXCodex(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(nil))
	c.Request.Header.Set("User-Agent", "codex_cli_rs/0.1.0")

	originalBody := []byte(`{"model":"gpt-5.2","stream":false,"input":[{"type":"text","text":"hi"placeholder]placeholder`)

	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	headers.Set("x-request-id", "rid")
	headers.Set("x-codex-primary-used-percent", "12")
	headers.Set("x-codex-secondary-used-percent", "34")
	headers.Set("x-codex-primary-window-minutes", "300")
	headers.Set("x-codex-secondary-window-minutes", "10080")
	headers.Set("x-codex-primary-reset-after-seconds", "1")

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     headers,
		Body:       io.NopCloser(strings.NewReader(`{"output":[],"usage":{"input_tokens":1,"output_tokens":1,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholder`)),
placeholder
	upstream := &httpUpstreamRecorder{resp: respplaceholder

	svc := &OpenAIGatewayService{
		cfg:          &config.Config{Gateway: config.GatewayConfig{ForceCodexCLI: falseplaceholderplaceholder,
		httpUpstream: upstream,
placeholder

	account := &Account{
		ID:             123,
		Name:           "acc",
		Platform:       PlatformOpenAI,
		Type:           AccountTypeOAuth,
		Concurrency:    1,
		Credentials:    map[string]any{"access_token": "oauth-token", "chatgpt_account_id": "chatgpt-acc"placeholder,
		Extra:          map[string]any{"openai_passthrough": trueplaceholder,
		Status:         StatusActive,
		Schedulable:    true,
		RateMultiplier: f64p(1),
placeholder

	_, err := svc.Forward(context.Background(), c, account, originalBody)
placeholder

	require.Equal(t, "12", rec.Header().Get("x-codex-primary-used-percent"))
	require.Equal(t, "34", rec.Header().Get("x-codex-secondary-used-percent"))
placeholder

func TestOpenAIGatewayService_OAuthPassthrough_UpstreamErrorIncludesPassthroughFlag(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(nil))
	c.Request.Header.Set("User-Agent", "codex_cli_rs/0.1.0")

	originalBody := []byte(`{"model":"gpt-5.2","stream":false,"input":[{"type":"text","text":"hi"placeholder]placeholder`)

	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholder, "x-request-id": []string{"rid"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"bad"placeholderplaceholder`)),
placeholder
	upstream := &httpUpstreamRecorder{resp: respplaceholder

	svc := &OpenAIGatewayService{
		cfg:          &config.Config{Gateway: config.GatewayConfig{ForceCodexCLI: falseplaceholderplaceholder,
		httpUpstream: upstream,
placeholder

	account := &Account{
		ID:             123,
		Name:           "acc",
		Platform:       PlatformOpenAI,
		Type:           AccountTypeOAuth,
		Concurrency:    1,
		Credentials:    map[string]any{"access_token": "oauth-token", "chatgpt_account_id": "chatgpt-acc"placeholder,
		Extra:          map[string]any{"openai_passthrough": trueplaceholder,
		Status:         StatusActive,
		Schedulable:    true,
		RateMultiplier: f64p(1),
placeholder

	_, err := svc.Forward(context.Background(), c, account, originalBody)
placeholder

	// should append an upstream error event with passthrough=true
	v, ok := c.Get(OpsUpstreamErrorsKey)
	require.True(t, ok)
	arr, ok := v.([]*OpsUpstreamErrorEvent)
	require.True(t, ok)
	require.NotEmpty(t, arr)
	require.True(t, arr[len(arr)-1].Passthrough)
placeholder

func TestOpenAIGatewayService_OAuthPassthrough_NonCodexUAFallbackToCodexUA(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(nil))
	// Non-Codex UA
	c.Request.Header.Set("User-Agent", "curl/8.0")

	inputBody := []byte(`{"model":"gpt-5.2","stream":false,"store":true,"input":[{"type":"text","text":"hi"placeholder]placeholder`)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader("data: [DONE]\n\n")),
placeholder
	upstream := &httpUpstreamRecorder{resp: respplaceholder

	svc := &OpenAIGatewayService{
		cfg:          &config.Config{Gateway: config.GatewayConfig{ForceCodexCLI: falseplaceholderplaceholder,
		httpUpstream: upstream,
placeholder

	account := &Account{
		ID:             123,
		Name:           "acc",
		Platform:       PlatformOpenAI,
		Type:           AccountTypeOAuth,
		Concurrency:    1,
		Credentials:    map[string]any{"access_token": "oauth-token", "chatgpt_account_id": "chatgpt-acc"placeholder,
		Extra:          map[string]any{"openai_passthrough": trueplaceholder,
		Status:         StatusActive,
		Schedulable:    true,
		RateMultiplier: f64p(1),
placeholder

	_, err := svc.Forward(context.Background(), c, account, inputBody)
placeholder
	require.Equal(t, inputBody, upstream.lastBody)
	require.Equal(t, "codex_cli_rs/0.98.0", upstream.lastReq.Header.Get("User-Agent"))
placeholder

func TestOpenAIGatewayService_OAuthPassthrough_StreamingSetsFirstTokenMs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(nil))
	c.Request.Header.Set("User-Agent", "codex_cli_rs/0.1.0")

	originalBody := []byte(`{"model":"gpt-5.2","stream":true,"input":[{"type":"text","text":"hi"placeholder]placeholder`)

	upstreamSSE := strings.Join([]string{
		`data: {"type":"response.output_text.delta","delta":"h"placeholder`,
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamSSE)),
placeholder
	upstream := &httpUpstreamRecorder{resp: respplaceholder

	svc := &OpenAIGatewayService{
		cfg:          &config.Config{Gateway: config.GatewayConfig{ForceCodexCLI: falseplaceholderplaceholder,
		httpUpstream: upstream,
placeholder

	account := &Account{
		ID:             123,
		Name:           "acc",
		Platform:       PlatformOpenAI,
		Type:           AccountTypeOAuth,
		Concurrency:    1,
		Credentials:    map[string]any{"access_token": "oauth-token", "chatgpt_account_id": "chatgpt-acc"placeholder,
		Extra:          map[string]any{"openai_passthrough": trueplaceholder,
		Status:         StatusActive,
		Schedulable:    true,
		RateMultiplier: f64p(1),
placeholder

	start := time.Now()
	result, err := svc.Forward(context.Background(), c, account, originalBody)
placeholder
	// sanity: duration after start
	require.GreaterOrEqual(t, time.Since(start), time.Duration(0))
	require.NotNil(t, result.FirstTokenMs)
	require.GreaterOrEqual(t, *result.FirstTokenMs, 0)
placeholder

func TestOpenAIGatewayService_OAuthPassthrough_StreamClientDisconnectStillCollectsUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(nil))
	c.Request.Header.Set("User-Agent", "codex_cli_rs/0.1.0")
	// 首次写入成功，后续写入失败，模拟客户端中途断开。
	c.Writer = &failingGinWriter{ResponseWriter: c.Writer, failAfter: 1placeholder

	originalBody := []byte(`{"model":"gpt-5.2","stream":true,"input":[{"type":"text","text":"hi"placeholder]placeholder`)

	upstreamSSE := strings.Join([]string{
		`data: {"type":"response.output_text.delta","delta":"h"placeholder`,
		"",
		`data: {"type":"response.completed","response":{"usage":{"input_tokens":11,"output_tokens":7,"input_tokens_details":{"cached_tokens":3placeholderplaceholderplaceholderplaceholder`,
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamSSE)),
placeholder
	upstream := &httpUpstreamRecorder{resp: respplaceholder

	svc := &OpenAIGatewayService{
		cfg:          &config.Config{Gateway: config.GatewayConfig{ForceCodexCLI: falseplaceholderplaceholder,
		httpUpstream: upstream,
placeholder

	account := &Account{
		ID:             123,
		Name:           "acc",
		Platform:       PlatformOpenAI,
		Type:           AccountTypeOAuth,
		Concurrency:    1,
		Credentials:    map[string]any{"access_token": "oauth-token", "chatgpt_account_id": "chatgpt-acc"placeholder,
		Extra:          map[string]any{"openai_passthrough": trueplaceholder,
		Status:         StatusActive,
		Schedulable:    true,
		RateMultiplier: f64p(1),
placeholder

	result, err := svc.Forward(context.Background(), c, account, originalBody)
placeholder
	require.NotNil(t, result)
	require.True(t, result.Stream)
	require.NotNil(t, result.FirstTokenMs)
	require.Equal(t, 11, result.Usage.InputTokens)
	require.Equal(t, 7, result.Usage.OutputTokens)
	require.Equal(t, 3, result.Usage.CacheReadInputTokens)
placeholder

func TestOpenAIGatewayService_APIKeyPassthrough_PreservesBodyAndUsesResponsesEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(nil))
	c.Request.Header.Set("User-Agent", "curl/8.0")
	c.Request.Header.Set("X-Test", "keep")

	originalBody := []byte(`{"model":"gpt-5.2","stream":false,"max_output_tokens":128,"input":[{"type":"text","text":"hi"placeholder]placeholder`)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholder, "x-request-id": []string{"rid"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"output":[],"usage":{"input_tokens":1,"output_tokens":1,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholder`)),
placeholder
	upstream := &httpUpstreamRecorder{resp: respplaceholder

	svc := &OpenAIGatewayService{
		cfg:          &config.Config{Gateway: config.GatewayConfig{ForceCodexCLI: falseplaceholderplaceholder,
		httpUpstream: upstream,
placeholder

	account := &Account{
		ID:             456,
		Name:           "apikey-acc",
		Platform:       PlatformOpenAI,
		Type:           AccountTypeAPIKey,
		Concurrency:    1,
		Credentials:    map[string]any{"api_key": "sk-api-key", "base_url": "https://api.openai.com"placeholder,
		Extra:          map[string]any{"openai_passthrough": trueplaceholder,
		Status:         StatusActive,
		Schedulable:    true,
		RateMultiplier: f64p(1),
placeholder

	_, err := svc.Forward(context.Background(), c, account, originalBody)
placeholder
	require.NotNil(t, upstream.lastReq)
	require.Equal(t, originalBody, upstream.lastBody)
	require.Equal(t, "https://api.openai.com/v1/responses", upstream.lastReq.URL.String())
	require.Equal(t, "Bearer sk-api-key", upstream.lastReq.Header.Get("Authorization"))
	require.Equal(t, "curl/8.0", upstream.lastReq.Header.Get("User-Agent"))
	require.Empty(t, upstream.lastReq.Header.Get("X-Test"))
placeholder

func TestOpenAIGatewayService_OAuthPassthrough_WarnOnTimeoutHeadersForStream(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logSink, restore := captureStructuredLog(t)
	defer restore()

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(nil))
	c.Request.Header.Set("User-Agent", "codex_cli_rs/0.1.0")
	c.Request.Header.Set("x-stainless-timeout", "10000")

	originalBody := []byte(`{"model":"gpt-5.2","stream":true,"input":[{"type":"text","text":"hi"placeholder]placeholder`)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "X-Request-Id": []string{"rid-timeout"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader("data: [DONE]\n\n")),
placeholder
	upstream := &httpUpstreamRecorder{resp: respplaceholder
	svc := &OpenAIGatewayService{
		cfg:          &config.Config{Gateway: config.GatewayConfig{ForceCodexCLI: falseplaceholderplaceholder,
		httpUpstream: upstream,
placeholder
	account := &Account{
		ID:             321,
		Name:           "acc",
		Platform:       PlatformOpenAI,
		Type:           AccountTypeOAuth,
		Concurrency:    1,
		Credentials:    map[string]any{"access_token": "oauth-token", "chatgpt_account_id": "chatgpt-acc"placeholder,
		Extra:          map[string]any{"openai_passthrough": trueplaceholder,
		Status:         StatusActive,
		Schedulable:    true,
		RateMultiplier: f64p(1),
placeholder

	_, err := svc.Forward(context.Background(), c, account, originalBody)
placeholder
	require.True(t, logSink.ContainsMessage("检测到超时相关请求头，将按配置过滤以降低断流风险"))
	require.True(t, logSink.ContainsFieldValue("timeout_headers", "x-stainless-timeout=10000"))
placeholder

func TestOpenAIGatewayService_OAuthPassthrough_WarnWhenStreamEndsWithoutDone(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logSink, restore := captureStructuredLog(t)
	defer restore()

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(nil))
	c.Request.Header.Set("User-Agent", "codex_cli_rs/0.1.0")

	originalBody := []byte(`{"model":"gpt-5.2","stream":true,"input":[{"type":"text","text":"hi"placeholder]placeholder`)
	// 注意：刻意不发送 [DONE]，模拟上游中途断流。
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "X-Request-Id": []string{"rid-truncate"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader("data: {\"type\":\"response.output_text.delta\",\"delta\":\"h\"placeholder\n\n")),
placeholder
	upstream := &httpUpstreamRecorder{resp: respplaceholder
	svc := &OpenAIGatewayService{
		cfg:          &config.Config{Gateway: config.GatewayConfig{ForceCodexCLI: falseplaceholderplaceholder,
		httpUpstream: upstream,
placeholder
	account := &Account{
		ID:             654,
		Name:           "acc",
		Platform:       PlatformOpenAI,
		Type:           AccountTypeOAuth,
		Concurrency:    1,
		Credentials:    map[string]any{"access_token": "oauth-token", "chatgpt_account_id": "chatgpt-acc"placeholder,
		Extra:          map[string]any{"openai_passthrough": trueplaceholder,
		Status:         StatusActive,
		Schedulable:    true,
		RateMultiplier: f64p(1),
placeholder

	_, err := svc.Forward(context.Background(), c, account, originalBody)
placeholder
	require.True(t, logSink.ContainsMessage("上游流在未收到 [DONE] 时结束，疑似断流"))
	require.True(t, logSink.ContainsFieldValue("upstream_request_id", "rid-truncate"))
placeholder

func TestOpenAIGatewayService_OAuthPassthrough_DefaultFiltersTimeoutHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(nil))
	c.Request.Header.Set("User-Agent", "codex_cli_rs/0.1.0")
	c.Request.Header.Set("x-stainless-timeout", "120000")
	c.Request.Header.Set("X-Test", "keep")

	originalBody := []byte(`{"model":"gpt-5.2","stream":false,"input":[{"type":"text","text":"hi"placeholder]placeholder`)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholder, "X-Request-Id": []string{"rid-filter-default"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"output":[],"usage":{"input_tokens":1,"output_tokens":1,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholder`)),
placeholder
	upstream := &httpUpstreamRecorder{resp: respplaceholder
	svc := &OpenAIGatewayService{
		cfg:          &config.Config{Gateway: config.GatewayConfig{ForceCodexCLI: falseplaceholderplaceholder,
		httpUpstream: upstream,
placeholder
	account := &Account{
		ID:             111,
		Name:           "acc",
		Platform:       PlatformOpenAI,
		Type:           AccountTypeOAuth,
		Concurrency:    1,
		Credentials:    map[string]any{"access_token": "oauth-token", "chatgpt_account_id": "chatgpt-acc"placeholder,
		Extra:          map[string]any{"openai_passthrough": trueplaceholder,
		Status:         StatusActive,
		Schedulable:    true,
		RateMultiplier: f64p(1),
placeholder

	_, err := svc.Forward(context.Background(), c, account, originalBody)
placeholder
	require.NotNil(t, upstream.lastReq)
	require.Empty(t, upstream.lastReq.Header.Get("x-stainless-timeout"))
	require.Empty(t, upstream.lastReq.Header.Get("X-Test"))
placeholder

func TestOpenAIGatewayService_OAuthPassthrough_AllowTimeoutHeadersWhenConfigured(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(nil))
	c.Request.Header.Set("User-Agent", "codex_cli_rs/0.1.0")
	c.Request.Header.Set("x-stainless-timeout", "120000")
	c.Request.Header.Set("X-Test", "keep")

	originalBody := []byte(`{"model":"gpt-5.2","stream":false,"input":[{"type":"text","text":"hi"placeholder]placeholder`)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholder, "X-Request-Id": []string{"rid-filter-allow"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"output":[],"usage":{"input_tokens":1,"output_tokens":1,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholder`)),
placeholder
	upstream := &httpUpstreamRecorder{resp: respplaceholder
	svc := &OpenAIGatewayService{
		cfg: &config.Config{Gateway: config.GatewayConfig{
			ForceCodexCLI:                        false,
			OpenAIPassthroughAllowTimeoutHeaders: true,
placeholder
		httpUpstream: upstream,
placeholder
	account := &Account{
		ID:             222,
		Name:           "acc",
		Platform:       PlatformOpenAI,
		Type:           AccountTypeOAuth,
		Concurrency:    1,
		Credentials:    map[string]any{"access_token": "oauth-token", "chatgpt_account_id": "chatgpt-acc"placeholder,
		Extra:          map[string]any{"openai_passthrough": trueplaceholder,
		Status:         StatusActive,
		Schedulable:    true,
		RateMultiplier: f64p(1),
placeholder

	_, err := svc.Forward(context.Background(), c, account, originalBody)
placeholder
	require.NotNil(t, upstream.lastReq)
	require.Equal(t, "120000", upstream.lastReq.Header.Get("x-stainless-timeout"))
	require.Empty(t, upstream.lastReq.Header.Get("X-Test"))
placeholder
