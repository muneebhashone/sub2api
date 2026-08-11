package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/claude"
	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

type anthropicHTTPUpstreamRecorder struct {
	lastReq  *http.Request
	lastBody []byte
	resp     *http.Response
	err      error
placeholder

func newAnthropicAPIKeyAccountForTest() *Account {
placeholder
		ID:          201,
		Name:        "anthropic-apikey-pass-test",
		Platform:    PlatformAnthropic,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "upstream-anthropic-key",
			"base_url": "https://api.anthropic.com",
	placeholder,
		Extra: map[string]any{
			"anthropic_passthrough": true,
	placeholder,
		Status:      StatusActive,
		Schedulable: true,
placeholder
placeholder

func (u *anthropicHTTPUpstreamRecorder) Do(req *http.Request, proxyURL string, accountID int64, accountConcurrency int) (*http.Response, error) {
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

func (u *anthropicHTTPUpstreamRecorder) DoWithTLS(req *http.Request, proxyURL string, accountID int64, accountConcurrency int, profile *tlsfingerprint.Profile) (*http.Response, error) {
	return u.Do(req, proxyURL, accountID, accountConcurrency)
placeholder

type streamReadCloser struct {
	payload []byte
	sent    bool
	err     error
placeholder

func (r *streamReadCloser) Read(p []byte) (int, error) {
	if !r.sent {
		r.sent = true
		n := copy(p, r.payload)
		return n, nil
placeholder
	if r.err != nil {
		return 0, r.err
placeholder
	return 0, io.EOF
placeholder

func (r *streamReadCloser) Close() error { return nil placeholder

type failWriteResponseWriter struct {
	gin.ResponseWriter
placeholder

func (w *failWriteResponseWriter) Write(data []byte) (int, error) {
	return 0, errors.New("client disconnected")
placeholder

func (w *failWriteResponseWriter) WriteString(_ string) (int, error) {
	return 0, errors.New("client disconnected")
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_ForwardStreamPreservesBodyAndAuthReplacement(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
	c.Request.Header.Set("User-Agent", "claude-cli/1.0.0")
	c.Request.Header.Set("Authorization", "Bearer inbound-token")
	c.Request.Header.Set("X-Api-Key", "inbound-api-key")
	c.Request.Header.Set("X-Goog-Api-Key", "inbound-goog-key")
	c.Request.Header.Set("Cookie", "secret=1")
	c.Request.Header.Set("Anthropic-Beta", "placeholder")

	body := []byte(`{"model":"claude-3-7-sonnet-20250219","stream":true,"system":[{"type":"text","text":"x-anthropic-billing-header keep"placeholder],"messages":[{"role":"user","content":[{"type":"text","text":"hello"placeholder]placeholder]placeholder`)
	parsed := &ParsedRequest{
		Body:   NewRequestBodyRef(body),
		Model:  "claude-3-7-sonnet-20250219",
		Stream: true,
placeholder

	upstreamSSE := strings.Join([]string{
		`data: {"type":"message_start","message":{"usage":{"input_tokens":9,"cached_tokens":7placeholderplaceholderplaceholder`,
		"",
		`data: {"type":"message_delta","usage":{"output_tokens":3placeholderplaceholder`,
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	upstream := &anthropicHTTPUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header: http.Header{
				"Content-Type": []string{"text/event-stream"placeholder,
				"x-request-id": []string{"rid-anthropic-pass"placeholder,
				"Set-Cookie":   []string{"secret=upstream"placeholder,
		placeholder,
			Body: io.NopCloser(strings.NewReader(upstreamSSE)),
	placeholder,
placeholder

	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			MaxLineSize: defaultMaxLineSize,
	placeholder,
placeholder
	svc := &GatewayService{
		cfg:                  cfg,
		responseHeaderFilter: compileResponseHeaderFilter(cfg),
		httpUpstream:         upstream,
		rateLimitService:     &RateLimitService{placeholder,
		deferredService:      &DeferredService{placeholder,
		billingCacheService:  nil,
placeholder

	account := &Account{
		ID:          101,
		Name:        "anthropic-apikey-pass",
		Platform:    PlatformAnthropic,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":       "upstream-anthropic-key",
			"base_url":      "https://api.anthropic.com",
			"model_mapping": map[string]any{"claude-3-7-sonnet-20250219": "claude-3-haiku-20240307"placeholder,
	placeholder,
		Extra: map[string]any{
			"anthropic_passthrough": true,
	placeholder,
		Status:      StatusActive,
		Schedulable: true,
placeholder

	result, err := svc.Forward(context.Background(), c, account, parsed)
placeholder
	require.NotNil(t, result)
	require.True(t, result.Stream)

	require.Equal(t, "claude-3-haiku-20240307", gjson.GetBytes(upstream.lastBody, "model").String(), "透传模式应应用账号级模型映射")

	require.Equal(t, "upstream-anthropic-key", getHeaderRaw(upstream.lastReq.Header, "x-api-key"))
	require.Empty(t, getHeaderRaw(upstream.lastReq.Header, "authorization"))
	require.Empty(t, getHeaderRaw(upstream.lastReq.Header, "x-goog-api-key"))
	require.Empty(t, getHeaderRaw(upstream.lastReq.Header, "cookie"))
	require.Equal(t, "2023-06-01", getHeaderRaw(upstream.lastReq.Header, "anthropic-version"))
	require.Equal(t, "placeholder", getHeaderRaw(upstream.lastReq.Header, "anthropic-beta"))
	require.Empty(t, getHeaderRaw(upstream.lastReq.Header, "x-stainless-lang"), "API Key 透传不应注入 OAuth 指纹头")

	require.Contains(t, rec.Body.String(), `"cached_tokens":7`)
	require.NotContains(t, rec.Body.String(), `"cache_read_input_tokens":7`, "透传输出不应被网关改写")
	require.Equal(t, 7, result.Usage.CacheReadInputTokens, "计费 usage 解析应保留 cached_tokens 兼容")
	require.Empty(t, rec.Header().Get("Set-Cookie"), "响应头应经过安全过滤")
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_ForwardCountTokensPreservesBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages/count_tokens", nil)
	c.Request.Header.Set("Authorization", "Bearer inbound-token")
	c.Request.Header.Set("X-Api-Key", "inbound-api-key")
	c.Request.Header.Set("Cookie", "secret=1")

	body := []byte(`{"model":"claude-3-5-sonnet-latest","messages":[{"role":"user","content":[{"type":"text","text":"hello"placeholder]placeholder],"thinking":{"type":"enabled"placeholderplaceholder`)
	parsed := &ParsedRequest{
		Body:  NewRequestBodyRef(body),
		Model: "claude-3-5-sonnet-latest",
placeholder

	upstreamRespBody := `{"input_tokens":42placeholder`
	upstream := &anthropicHTTPUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header: http.Header{
				"Content-Type": []string{"application/json"placeholder,
				"x-request-id": []string{"rid-count"placeholder,
				"Set-Cookie":   []string{"secret=upstream"placeholder,
		placeholder,
			Body: io.NopCloser(strings.NewReader(upstreamRespBody)),
	placeholder,
placeholder

	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			MaxLineSize: defaultMaxLineSize,
	placeholder,
placeholder
	svc := &GatewayService{
		cfg:                  cfg,
		responseHeaderFilter: compileResponseHeaderFilter(cfg),
		httpUpstream:         upstream,
		rateLimitService:     &RateLimitService{placeholder,
placeholder

	account := &Account{
		ID:          102,
		Name:        "anthropic-apikey-pass-count",
		Platform:    PlatformAnthropic,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":       "upstream-anthropic-key",
			"base_url":      "https://api.anthropic.com",
			"model_mapping": map[string]any{"claude-3-5-sonnet-latest": "claude-3-opus-20240229"placeholder,
	placeholder,
		Extra: map[string]any{
			"anthropic_passthrough": true,
	placeholder,
		Status:      StatusActive,
		Schedulable: true,
placeholder

	err := svc.ForwardCountTokens(context.Background(), c, account, parsed)
placeholder

	require.Equal(t, "claude-3-opus-20240229", gjson.GetBytes(upstream.lastBody, "model").String(), "count_tokens 透传模式应应用账号级模型映射")
	require.Equal(t, "upstream-anthropic-key", getHeaderRaw(upstream.lastReq.Header, "x-api-key"))
	require.Empty(t, getHeaderRaw(upstream.lastReq.Header, "authorization"))
	require.Empty(t, getHeaderRaw(upstream.lastReq.Header, "cookie"))
	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, upstreamRespBody, rec.Body.String())
	require.Empty(t, rec.Header().Get("Set-Cookie"))
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_BearerAuthScheme(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
	c.Request.Header.Set("Authorization", "Bearer inbound-token")
	c.Request.Header.Set("X-Api-Key", "inbound-api-key")
	c.Request.Header.Set("Cookie", "secret=1")

	svc := &GatewayService{
		cfg: &config.Config{
			Security: config.SecurityConfig{
				URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholder,
		placeholder,
	placeholder,
placeholder
	account := &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeAPIKey,
placeholder
			"api_key":  "ollama-key",
			"base_url": "https://ollama.com",
	placeholder,
		Extra: map[string]any{
			"anthropic_passthrough":        true,
			"anthropic_apikey_auth_scheme": AnthropicAPIKeyAuthSchemeAuthorizationBearer,
	placeholder,
placeholder

	msgReq, wireBody, err := svc.buildUpstreamRequestAnthropicAPIKeyPassthrough(
		context.Background(), c, account, []byte(`{"model":"gpt-oss:20b","messages":[]placeholder`), "ollama-key",
	)
placeholder
	require.Equal(t, "https://ollama.com/v1/messages?beta=true", msgReq.URL.String())
	require.JSONEq(t, `{"model":"gpt-oss:20b","messages":[]placeholder`, string(wireBody))
	require.Equal(t, "Bearer ollama-key", getHeaderRaw(msgReq.Header, "authorization"))
	require.Empty(t, getHeaderRaw(msgReq.Header, "x-api-key"))
	require.Empty(t, getHeaderRaw(msgReq.Header, "cookie"))

	countReq, err := svc.buildCountTokensRequestAnthropicAPIKeyPassthrough(
		context.Background(), c, account, []byte(`{"model":"gpt-oss:20b","messages":[]placeholder`), "ollama-key",
	)
placeholder
	require.Equal(t, "https://ollama.com/v1/messages/count_tokens?beta=true", countReq.URL.String())
	require.Equal(t, "Bearer ollama-key", getHeaderRaw(countReq.Header, "authorization"))
	require.Empty(t, getHeaderRaw(countReq.Header, "x-api-key"))
	require.Empty(t, getHeaderRaw(countReq.Header, "cookie"))
placeholder

// TestGatewayService_AnthropicAPIKeyPassthrough_ModelMappingEdgeCases 覆盖透传模式下模型映射的各种边界情况
func TestGatewayService_AnthropicAPIKeyPassthrough_ModelMappingEdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		model         string
		modelMapping  map[string]any // nil = 不配置映射
		expectedModel string
		endpoint      string // "messages" or "count_tokens"
placeholder{
		{
			name:          "Forward: 无映射配置时不改写模型",
			model:         "claude-sonnet-4-20250514",
			modelMapping:  nil,
			expectedModel: "claude-sonnet-4-20250514",
			endpoint:      "messages",
	placeholder,
		{
			name:          "Forward: 空映射配置时不改写模型",
			model:         "claude-sonnet-4-20250514",
			modelMapping:  map[string]any{placeholder,
			expectedModel: "claude-sonnet-4-20250514",
			endpoint:      "messages",
	placeholder,
		{
			name:          "Forward: 模型不在映射表中时不改写",
			model:         "claude-sonnet-4-20250514",
			modelMapping:  map[string]any{"claude-3-haiku-20240307": "claude-3-opus-20240229"placeholder,
			expectedModel: "claude-sonnet-4-20250514",
			endpoint:      "messages",
	placeholder,
		{
			name:          "Forward: 精确匹配映射应改写模型",
			model:         "claude-sonnet-4-20250514",
			modelMapping:  map[string]any{"claude-sonnet-4-20250514": "claude-sonnet-4-5-20241022"placeholder,
			expectedModel: "claude-sonnet-4-5-20241022",
			endpoint:      "messages",
	placeholder,
		{
			name:          "Forward: 通配符映射应改写模型",
			model:         "claude-sonnet-4-20250514",
			modelMapping:  map[string]any{"claude-sonnet-4-*": "claude-sonnet-4-5-20241022"placeholder,
			expectedModel: "claude-sonnet-4-5-20241022",
			endpoint:      "messages",
	placeholder,
		{
			name:          "CountTokens: 无映射配置时不改写模型",
			model:         "claude-sonnet-4-20250514",
			modelMapping:  nil,
			expectedModel: "claude-sonnet-4-20250514",
			endpoint:      "count_tokens",
	placeholder,
		{
			name:          "CountTokens: 模型不在映射表中时不改写",
			model:         "claude-sonnet-4-20250514",
			modelMapping:  map[string]any{"claude-3-haiku-20240307": "claude-3-opus-20240229"placeholder,
			expectedModel: "claude-sonnet-4-20250514",
			endpoint:      "count_tokens",
	placeholder,
		{
			name:          "CountTokens: 精确匹配映射应改写模型",
			model:         "claude-sonnet-4-20250514",
			modelMapping:  map[string]any{"claude-sonnet-4-20250514": "claude-sonnet-4-5-20241022"placeholder,
			expectedModel: "claude-sonnet-4-5-20241022",
			endpoint:      "count_tokens",
	placeholder,
		{
			name:          "CountTokens: 通配符映射应改写模型",
			model:         "claude-sonnet-4-20250514",
			modelMapping:  map[string]any{"claude-sonnet-4-*": "claude-sonnet-4-5-20241022"placeholder,
			expectedModel: "claude-sonnet-4-5-20241022",
			endpoint:      "count_tokens",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)

			body := []byte(`{"model":"` + tt.model + `","messages":[{"role":"user","content":[{"type":"text","text":"hello"placeholder]placeholder]placeholder`)
			parsed := &ParsedRequest{
				Body:  NewRequestBodyRef(body),
				Model: tt.model,
		placeholder

			credentials := map[string]any{
				"api_key":  "upstream-key",
				"base_url": "https://api.anthropic.com",
		placeholder
			if tt.modelMapping != nil {
				credentials["model_mapping"] = tt.modelMapping
		placeholder

			account := &Account{
				ID:          300,
				Name:        "edge-case-test",
				Platform:    PlatformAnthropic,
				Type:        AccountTypeAPIKey,
				Concurrency: 1,
				Credentials: credentials,
				Extra:       map[string]any{"anthropic_passthrough": trueplaceholder,
				Status:      StatusActive,
				Schedulable: true,
		placeholder

			if tt.endpoint == "messages" {
				c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
				parsed.Stream = false

				upstreamJSON := `{"id":"msg_1","type":"message","usage":{"input_tokens":5,"output_tokens":3placeholderplaceholder`
				upstream := &anthropicHTTPUpstreamRecorder{
					resp: &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
						Body:       io.NopCloser(strings.NewReader(upstreamJSON)),
				placeholder,
			placeholder
				svc := &GatewayService{
					cfg:              &config.Config{placeholder,
					httpUpstream:     upstream,
					rateLimitService: &RateLimitService{placeholder,
			placeholder

				result, err := svc.Forward(context.Background(), c, account, parsed)
			placeholder
				require.NotNil(t, result)
				require.Equal(t, tt.expectedModel, gjson.GetBytes(upstream.lastBody, "model").String(),
					"Forward 上游请求体中的模型应为: %s", tt.expectedModel)
		placeholder else {
				c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages/count_tokens", nil)

				upstreamRespBody := `{"input_tokens":42placeholder`
				upstream := &anthropicHTTPUpstreamRecorder{
					resp: &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
						Body:       io.NopCloser(strings.NewReader(upstreamRespBody)),
				placeholder,
			placeholder
				svc := &GatewayService{
					cfg:              &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder,
					httpUpstream:     upstream,
					rateLimitService: &RateLimitService{placeholder,
			placeholder

				err := svc.ForwardCountTokens(context.Background(), c, account, parsed)
			placeholder
				require.Equal(t, tt.expectedModel, gjson.GetBytes(upstream.lastBody, "model").String(),
					"CountTokens 上游请求体中的模型应为: %s", tt.expectedModel)
		placeholder
	placeholder)
placeholder
placeholder

// TestGatewayService_AnthropicAPIKeyPassthrough_ModelMappingPreservesOtherFields
// 确保模型映射只替换 model 字段，不影响请求体中的其他字段
func TestGatewayService_AnthropicAPIKeyPassthrough_ModelMappingPreservesOtherFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages/count_tokens", nil)

	// 包含复杂字段的请求体：system、thinking、messages
	body := []byte(`{"model":"claude-sonnet-4-20250514","system":[{"type":"text","text":"You are a helpful assistant."placeholder],"messages":[{"role":"user","content":[{"type":"text","text":"hello world"placeholder]placeholder],"thinking":{"type":"enabled","budget_tokens":5000placeholder,"max_tokens":placeholder`)
	parsed := &ParsedRequest{
		Body:  NewRequestBodyRef(body),
		Model: "claude-sonnet-4-20250514",
placeholder

	upstreamRespBody := `{"input_tokens":42placeholder`
	upstream := &anthropicHTTPUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(upstreamRespBody)),
	placeholder,
placeholder

	svc := &GatewayService{
		cfg:              &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder,
		httpUpstream:     upstream,
		rateLimitService: &RateLimitService{placeholder,
placeholder

	account := &Account{
		ID:          301,
		Name:        "preserve-fields-test",
		Platform:    PlatformAnthropic,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":       "upstream-key",
			"base_url":      "https://api.anthropic.com",
			"model_mapping": map[string]any{"claude-sonnet-4-20250514": "claude-sonnet-4-5-20241022"placeholder,
	placeholder,
		Extra:       map[string]any{"anthropic_passthrough": trueplaceholder,
		Status:      StatusActive,
		Schedulable: true,
placeholder

	err := svc.ForwardCountTokens(context.Background(), c, account, parsed)
placeholder

	sentBody := upstream.lastBody
	require.Equal(t, "claude-sonnet-4-5-20241022", gjson.GetBytes(sentBody, "model").String(), "model 应被映射")
	require.Equal(t, "You are a helpful assistant.", gjson.GetBytes(sentBody, "system.0.text").String(), "system 字段不应被修改")
	require.Equal(t, "hello world", gjson.GetBytes(sentBody, "messages.0.content.0.text").String(), "messages 字段不应被修改")
	require.Equal(t, "enabled", gjson.GetBytes(sentBody, "thinking.type").String(), "thinking 字段不应被修改")
	require.Equal(t, int64(5000), gjson.GetBytes(sentBody, "thinking.budget_tokens").Int(), "thinking.budget_tokens 不应被修改")
	require.False(t, gjson.GetBytes(sentBody, "max_tokens").Exists(),
		"max_tokens 作为生成参数应被 count_tokens 过滤剥离")
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_CountTokensFiltersGenerationFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages/count_tokens", nil)

	body := []byte(`{"model":"claude-sonnet-4-20250514","system":[{"type":"text","text":"sys"placeholder],"messages":[{"role":"user","content":"hello"placeholder],"tools":[{"name":"tool","input_schema":{"type":"object"placeholderplaceholder],"temperature":0.7,"top_p":0.9,"top_k":40,"stream":true,"stop_sequences":["END"],"max_tokens":1024,"thinking":{"type":"enabled","budget_tokens":5000placeholderplaceholder`)
	parsed := &ParsedRequest{
		Body:  NewRequestBodyRef(body),
		Model: "claude-sonnet-4-20250514",
placeholder

	upstreamRespBody := `{"input_tokens":42placeholder`
	upstream := &anthropicHTTPUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(upstreamRespBody)),
	placeholder,
placeholder

	svc := &GatewayService{
		cfg:              &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder,
		httpUpstream:     upstream,
		rateLimitService: &RateLimitService{placeholder,
placeholder

	account := &Account{
		ID:          302,
		Name:        "count-token-filter-test",
		Platform:    PlatformAnthropic,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "upstream-key",
			"base_url": "https://api.anthropic.com",
	placeholder,
		Extra:       map[string]any{"anthropic_passthrough": trueplaceholder,
		Status:      StatusActive,
		Schedulable: true,
placeholder

	err := svc.ForwardCountTokens(context.Background(), c, account, parsed)
placeholder

	sentBody := upstream.lastBody
	require.False(t, gjson.GetBytes(sentBody, "temperature").Exists())
	require.False(t, gjson.GetBytes(sentBody, "top_p").Exists())
	require.False(t, gjson.GetBytes(sentBody, "top_k").Exists())
	require.False(t, gjson.GetBytes(sentBody, "stream").Exists())
	require.False(t, gjson.GetBytes(sentBody, "stop_sequences").Exists())
	require.Equal(t, "claude-sonnet-4-20250514", gjson.GetBytes(sentBody, "model").String())
	require.Equal(t, "sys", gjson.GetBytes(sentBody, "system.0.text").String())
	require.Equal(t, "hello", gjson.GetBytes(sentBody, "messages.0.content").String())
	require.Equal(t, "tool", gjson.GetBytes(sentBody, "tools.0.name").String())
	require.False(t, gjson.GetBytes(sentBody, "max_tokens").Exists(),
		"count_tokens 请求不得携带生成参数 max_tokens")
	require.Equal(t, "enabled", gjson.GetBytes(sentBody, "thinking.type").String())
placeholder

// TestGatewayService_AnthropicAPIKeyPassthrough_EmptyModelSkipsMapping
// 确保空模型名不会触发映射逻辑
func TestGatewayService_AnthropicAPIKeyPassthrough_EmptyModelSkipsMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages/count_tokens", nil)

	body := []byte(`{"messages":[{"role":"user","content":"hello"placeholder]placeholder`)
	parsed := &ParsedRequest{
		Body:  NewRequestBodyRef(body),
		Model: "", // 空模型
placeholder

	upstreamRespBody := `{"input_tokens":10placeholder`
	upstream := &anthropicHTTPUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(upstreamRespBody)),
	placeholder,
placeholder

	svc := &GatewayService{
		cfg:              &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder,
		httpUpstream:     upstream,
		rateLimitService: &RateLimitService{placeholder,
placeholder

	account := &Account{
		ID:          302,
		Name:        "empty-model-test",
		Platform:    PlatformAnthropic,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":       "upstream-key",
			"base_url":      "https://api.anthropic.com",
			"model_mapping": map[string]any{"*": "claude-3-opus-20240229"placeholder,
	placeholder,
		Extra:       map[string]any{"anthropic_passthrough": trueplaceholder,
		Status:      StatusActive,
		Schedulable: true,
placeholder

	err := svc.ForwardCountTokens(context.Background(), c, account, parsed)
placeholder
	// 空模型名时，body 应原样透传，不应触发映射
	require.Equal(t, body, upstream.lastBody, "空模型名时请求体不应被修改")
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_CountTokens404PassthroughNotError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		statusCode      int
		respBody        string
		wantPassthrough bool
placeholder{
		{
			name:            "404 endpoint not found passes through as 404",
			statusCode:      http.StatusNotFound,
			respBody:        `{"error":{"message":"Not found: /v1/messages/count_tokens","type":"not_found_error"placeholderplaceholder`,
			wantPassthrough: true,
	placeholder,
		{
			name:            "404 generic not found does not passthrough",
			statusCode:      http.StatusNotFound,
			respBody:        `{"error":{"message":"resource not found","type":"not_found_error"placeholderplaceholder`,
			wantPassthrough: false,
	placeholder,
		{
			name:            "400 Invalid URL does not passthrough",
			statusCode:      http.StatusBadRequest,
			respBody:        `{"error":{"message":"Invalid URL (POST /v1/messages/count_tokens)","type":"invalid_request_error"placeholderplaceholder`,
			wantPassthrough: false,
	placeholder,
		{
			name:            "400 model error does not passthrough",
			statusCode:      http.StatusBadRequest,
			respBody:        `{"error":{"message":"model not found: claude-unknown","type":"invalid_request_error"placeholderplaceholder`,
			wantPassthrough: false,
	placeholder,
		{
			name:            "500 internal error does not passthrough",
			statusCode:      http.StatusInternalServerError,
			respBody:        `{"error":{"message":"internal error","type":"api_error"placeholderplaceholder`,
			wantPassthrough: false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages/count_tokens", nil)

			body := []byte(`{"model":"claude-sonnet-4-5-20250929","messages":[{"role":"user","content":"hi"placeholder]placeholder`)
			parsed := &ParsedRequest{Body: NewRequestBodyRef(body), Model: "claude-sonnet-4-5-20250929"placeholder

			upstream := &anthropicHTTPUpstreamRecorder{
				resp: &http.Response{
					StatusCode: tt.statusCode,
					Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
					Body:       io.NopCloser(strings.NewReader(tt.respBody)),
			placeholder,
		placeholder

			svc := &GatewayService{
				cfg: &config.Config{
					Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
			placeholder,
				httpUpstream:     upstream,
				rateLimitService: nil,
		placeholder

			account := &Account{
				ID:          200,
				Name:        "proxy-acc",
				Platform:    PlatformAnthropic,
				Type:        AccountTypeAPIKey,
				Concurrency: 1,
		placeholder
					"api_key":  "sk-proxy",
					"base_url": "https://proxy.example.com",
			placeholder,
				Extra:       map[string]any{"anthropic_passthrough": trueplaceholder,
				Status:      StatusActive,
				Schedulable: true,
		placeholder

			err := svc.ForwardCountTokens(context.Background(), c, account, parsed)

			if tt.wantPassthrough {
				// 返回 nil（不记录为错误），HTTP 状态码 404 + Anthropic 错误体
			placeholder
				require.Equal(t, http.StatusNotFound, rec.Code)
				var errResp map[string]any
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &errResp))
				require.Equal(t, "error", errResp["type"])
				errObj, ok := errResp["error"].(map[string]any)
				require.True(t, ok)
				require.Equal(t, "not_found_error", errObj["type"])
		placeholder else {
			placeholder
				require.Equal(t, tt.statusCode, rec.Code)
		placeholder
	placeholder)
placeholder
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_BuildRequestRejectsInvalidBaseURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	svc := &GatewayService{
		cfg: &config.Config{
			Security: config.SecurityConfig{
				URLAllowlist: config.URLAllowlistConfig{
					Enabled: false,
			placeholder,
		placeholder,
	placeholder,
placeholder
	account := &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeAPIKey,
placeholder
			"api_key":  "k",
			"base_url": "://invalid-url",
	placeholder,
placeholder

	_, _, err := svc.buildUpstreamRequestAnthropicAPIKeyPassthrough(context.Background(), c, account, []byte(`{placeholder`), "k")
placeholder
placeholder

func TestGatewayService_AnthropicOAuth_NotAffectedByAPIKeyPassthroughToggle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	svc := &GatewayService{
		cfg: &config.Config{
			Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
	placeholder,
placeholder
	account := &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeOAuth,
		Extra: map[string]any{
			"anthropic_passthrough": true,
	placeholder,
placeholder

	require.False(t, account.IsAnthropicAPIKeyPassthroughEnabled())

	req, _, err := svc.buildUpstreamRequest(context.Background(), c, account, []byte(`{"model":"claude-3-7-sonnet-20250219"placeholder`), "oauth-token", "oauth", "claude-3-7-sonnet-20250219", true, false)
placeholder
	require.Equal(t, "Bearer oauth-token", getHeaderRaw(req.Header, "authorization"))
	require.Contains(t, getHeaderRaw(req.Header, "anthropic-beta"), claude.BetaOAuth, "OAuth 链路仍应按原逻辑补齐 oauth beta")
placeholder

func TestGatewayService_AnthropicOAuthMimic_RewritesSystemWithBillingBlock(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name                       string
		body                       string
		wantModel                  string
		wantOriginalSystem         string
		wantOriginalSystemCacheTTL string
		wantMetadataUserID         string
placeholder{
		{
			name:               "sonnet system array",
			body:               `{"model":"claude-3-5-sonnet-latest","system":[{"type":"text","text":"x-anthropic-billing-header keep"placeholder],"messages":[{"role":"user","content":[{"type":"text","text":"hello"placeholder]placeholder]placeholder`,
			wantModel:          "claude-3-5-sonnet-latest",
			wantOriginalSystem: "x-anthropic-billing-header keep",
	placeholder,
		{
			name:               "sonnet system string",
			body:               `{"model":"claude-3-5-sonnet-latest","system":"x-anthropic-billing-header keep","messages":[{"role":"user","content":[{"type":"text","text":"hello"placeholder]placeholder]placeholder`,
			wantModel:          "claude-3-5-sonnet-latest",
			wantOriginalSystem: "x-anthropic-billing-header keep",
	placeholder,
		{
			name:                       "haiku full mimicry",
			body:                       `{"model":"claude-haiku-4-5","metadata":{"user_id":"pi-session-metadata"placeholder,"system":[{"type":"text","text":"Pi project instructions","cache_control":{"type":"ephemeral","ttl":"1h"placeholderplaceholder],"thinking":{"type":"enabled","budget_tokens":placeholder,"messages":[{"role":"user","content":[{"type":"text","text":"hello"placeholder]placeholder]placeholder`,
			wantModel:                  "placeholder",
			wantOriginalSystem:         "Pi project instructions",
			wantOriginalSystemCacheTTL: "1h",
			wantMetadataUserID:         "pi-session-metadata",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
			c.Request.Header.Set("User-Agent", "pi/0.51.0")
			c.Request.Header.Set("Anthropic-Beta", "client-only-beta")

			parsed, err := ParseGatewayRequest(NewRequestBodyRef([]byte(tt.body)), PlatformAnthropic)
		placeholder

			upstream := &anthropicHTTPUpstreamRecorder{
				resp: &http.Response{
					StatusCode: http.StatusOK,
					Header: http.Header{
						"Content-Type": []string{"application/json"placeholder,
						"x-request-id": []string{"rid-oauth-mimic"placeholder,
				placeholder,
					Body: io.NopCloser(strings.NewReader(`{"id":"msg_1","type":"message","role":"assistant","model":"placeholder","content":[{"type":"text","text":"ok"placeholder],"usage":{"input_tokens":12,"output_tokens":7placeholderplaceholder`)),
			placeholder,
		placeholder

			cfg := &config.Config{
				Gateway: config.GatewayConfig{
					MaxLineSize: defaultMaxLineSize,
			placeholder,
		placeholder
			svc := &GatewayService{
				cfg:                  cfg,
				responseHeaderFilter: compileResponseHeaderFilter(cfg),
				httpUpstream:         upstream,
				rateLimitService:     &RateLimitService{placeholder,
				deferredService:      &DeferredService{placeholder,
		placeholder

			account := &Account{
				ID:          301,
				Name:        "anthropic-oauth-mimic",
				Platform:    PlatformAnthropic,
				Type:        AccountTypeOAuth,
				Concurrency: 1,
		placeholder
					"access_token": "oauth-token",
			placeholder,
				Status:      StatusActive,
				Schedulable: true,
		placeholder

			result, err := svc.Forward(context.Background(), c, account, parsed)
		placeholder
			require.NotNil(t, result)
			require.NotNil(t, upstream.lastReq)
			require.Equal(t, "Bearer oauth-token", getHeaderRaw(upstream.lastReq.Header, "authorization"))
			finalBeta := getHeaderRaw(upstream.lastReq.Header, "anthropic-beta")
			for _, beta := range claude.FullClaudeCodeMimicryBetas() {
				require.Truef(t, anthropicBetaTokensContains(finalBeta, beta), "missing mimic beta %s", beta)
		placeholder
			require.False(t, anthropicBetaTokensContains(finalBeta, "client-only-beta"))
			for key, value := range claude.DefaultHeaders {
				require.Equal(t, value, getHeaderRaw(upstream.lastReq.Header, key), "mimic fingerprint header %s", key)
		placeholder
			require.NotEmpty(t, getHeaderRaw(upstream.lastReq.Header, "x-client-request-id"))

			require.Equal(t, tt.wantModel, gjson.GetBytes(upstream.lastBody, "model").String())
			system := gjson.GetBytes(upstream.lastBody, "system")
			require.True(t, system.Exists())
			require.True(t, system.IsArray(), "system should be an array")
			arr := system.Array()
			require.Len(t, arr, 3, "system array should have billing block + cc prompt block + expansion block")

			billingText := arr[0].Get("text").String()
			require.Contains(t, billingText, "x-anthropic-billing-header:")
			require.Contains(t, billingText, "cc_version="+claude.CLICurrentVersion+".")
			require.Contains(t, billingText, "cc_entrypoint=cli;")

			require.Equal(t, claudeCodeSystemPrompt, arr[1].Get("text").String())
			require.False(t, arr[1].Get("cache_control").Exists(), "身份前缀 block 不应带 cache_control")

			require.Equal(t, claudeCodeSystemPromptExpansion, arr[2].Get("text").String())
			require.Equal(t, "ephemeral", arr[2].Get("cache_control.type").String())

			// 原始 system prompt 应迁移至 messages 中。
			messages := gjson.GetBytes(upstream.lastBody, "messages")
			require.True(t, messages.IsArray())
			firstMsg := messages.Array()[0]
			require.Equal(t, "user", firstMsg.Get("role").String())
			require.Contains(t, firstMsg.Get("content.0.text").String(), tt.wantOriginalSystem)
			if tt.wantOriginalSystemCacheTTL != "" {
				require.Equal(t, "ephemeral", firstMsg.Get("content.0.cache_control.type").String())
				require.Equal(t, tt.wantOriginalSystemCacheTTL, firstMsg.Get("content.0.cache_control.ttl").String())
		placeholder else {
				require.False(t, firstMsg.Get("content.0.cache_control").Exists())
		placeholder

			if tt.wantMetadataUserID != "" {
				require.Equal(t, tt.wantMetadataUserID, gjson.GetBytes(upstream.lastBody, "metadata.user_id").String())
				require.True(t, gjson.GetBytes(upstream.lastBody, "context_management").Exists())
		placeholder
	placeholder)
placeholder
placeholder

func TestGatewayService_AnthropicOAuthRealClaudeCodeHaiku_PreservesClientHeadersAndBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	metadataUserID := FormatMetadataUserID(
		"a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
		"550e8400-e29b-41d4-a716-446655440000",
		"123e4567-e89b-42d3-a456-426614174000",
		claude.CLICurrentVersion,
	)
	body := []byte(`{"model":"placeholder","metadata":{"user_id":` + strconvQuote(metadataUserID) + `placeholder,"system":[{"type":"text","text":"Client-owned Claude Code system","cache_control":{"type":"ephemeral"placeholderplaceholder],"context_management":{"edits":[{"type":"clear_thinking_20251015","keep":"all"placeholder]placeholder,"messages":[{"role":"user","content":[{"type":"text","text":"hello"placeholder]placeholder]placeholder`)
	parsed, err := ParseGatewayRequest(NewRequestBodyRef(body), PlatformAnthropic)
placeholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
	c.Request.Header.Set("User-Agent", "claude-cli/"+claude.CLICurrentVersion+" (external, cli)")
	c.Request.Header.Set("X-Stainless-Package-Version", "real-client-package")
	clientBeta := strings.Join([]string{
		claude.BetaClaudeCode,
		claude.BetaOAuth,
		claude.BetaInterleavedThinking,
		claude.BetaContextManagement,
placeholder, ",")
	c.Request.Header.Set("Anthropic-Beta", clientBeta)

	upstream := &anthropicHTTPUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"id":"msg_real_cc","type":"message","role":"assistant","model":"placeholder","content":[{"type":"text","text":"ok"placeholder],"usage":{"input_tokens":12,"output_tokens":7placeholderplaceholder`)),
placeholderplaceholder
	cfg := &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder
	svc := &GatewayService{
		cfg:                  cfg,
		responseHeaderFilter: compileResponseHeaderFilter(cfg),
		httpUpstream:         upstream,
		rateLimitService:     &RateLimitService{placeholder,
		deferredService:      &DeferredService{placeholder,
placeholder
	account := &Account{
		ID: 302, Name: "anthropic-real-cc", Platform: PlatformAnthropic, Type: AccountTypeOAuth, Concurrency: 1,
placeholder"access_token": "oauth-token"placeholder, Status: StatusActive, Schedulable: true,
placeholder

	result, err := svc.Forward(context.Background(), c, account, parsed)
placeholder
	require.NotNil(t, result)
	require.NotNil(t, upstream.lastReq)
	require.Equal(t, c.Request.Header.Get("User-Agent"), getHeaderRaw(upstream.lastReq.Header, "User-Agent"))
	require.Equal(t, "real-client-package", getHeaderRaw(upstream.lastReq.Header, "X-Stainless-Package-Version"))
	require.Equal(t, clientBeta, getHeaderRaw(upstream.lastReq.Header, "anthropic-beta"))
	require.Empty(t, getHeaderRaw(upstream.lastReq.Header, "x-client-request-id"), "真实 CC 不应被强制写入 mimic request id")
	require.Equal(t, gjson.GetBytes(body, "system").Raw, gjson.GetBytes(upstream.lastBody, "system").Raw)
	require.Equal(t, gjson.GetBytes(body, "messages").Raw, gjson.GetBytes(upstream.lastBody, "messages").Raw)
	require.Equal(t, metadataUserID, gjson.GetBytes(upstream.lastBody, "metadata.user_id").String())
	require.True(t, gjson.GetBytes(upstream.lastBody, "context_management").Exists())
	require.NotContains(t, string(upstream.lastBody), "x-anthropic-billing-header:")
placeholder

func TestGatewayService_AnthropicOAuth_SystemPromptInjectionCanBeDisabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resetGatewayForwardingSettingsCacheForTest(t)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	body := []byte(`{"model":"claude-3-5-sonnet-latest","system":"Original system prompt","messages":[{"role":"user","content":[{"type":"text","text":"hello"placeholder]placeholder]placeholder`)
	parsed, err := ParseGatewayRequest(NewRequestBodyRef(body), PlatformAnthropic)
placeholder

	upstream := &anthropicHTTPUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header: http.Header{
				"Content-Type": []string{"application/json"placeholder,
				"x-request-id": []string{"rid-oauth-no-system-injection"placeholder,
		placeholder,
			Body: io.NopCloser(strings.NewReader(`{"id":"msg_1","type":"message","role":"assistant","model":"claude-3-5-sonnet-20241022","content":[{"type":"text","text":"ok"placeholder],"usage":{"input_tokens":12,"output_tokens":7placeholderplaceholder`)),
	placeholder,
placeholder

	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			MaxLineSize: defaultMaxLineSize,
	placeholder,
placeholder
	settingService := NewSettingService(&gatewayTTLSettingRepo{data: map[string]string{
		SettingKeyEnableClaudeOAuthSystemPromptInjection: "false",
placeholderplaceholder, cfg)
	svc := &GatewayService{
		cfg:                  cfg,
		responseHeaderFilter: compileResponseHeaderFilter(cfg),
		httpUpstream:         upstream,
		rateLimitService:     &RateLimitService{placeholder,
		deferredService:      &DeferredService{placeholder,
		settingService:       settingService,
placeholder

	account := &Account{
		ID:          302,
		Name:        "anthropic-oauth-no-system-injection",
		Platform:    PlatformAnthropic,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token": "oauth-token",
	placeholder,
		Status:      StatusActive,
		Schedulable: true,
placeholder

	result, err := svc.Forward(context.Background(), c, account, parsed)
placeholder
	require.NotNil(t, result)

	system := gjson.GetBytes(upstream.lastBody, "system")
	require.True(t, system.Exists())
	require.Equal(t, "Original system prompt", system.String())
	require.NotContains(t, string(upstream.lastBody), "x-anthropic-billing-header:")
	require.NotContains(t, string(upstream.lastBody), "[System Instructions]")
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_StreamingStillCollectsUsageAfterClientDisconnect(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Use a canceled context recorder to simulate client disconnect behavior.
	req := httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
	ctx, cancel := context.WithCancel(req.Context())
	cancel()
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	svc := &GatewayService{
		cfg: &config.Config{
			Gateway: config.GatewayConfig{
				MaxLineSize: defaultMaxLineSize,
		placeholder,
	placeholder,
		rateLimitService: &RateLimitService{placeholder,
placeholder

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(strings.Join([]string{
			`data: {"type":"message_start","message":{"usage":{"input_tokens":11placeholderplaceholderplaceholder`,
			"",
			`data: {"type":"message_delta","usage":{"output_tokens":5placeholderplaceholder`,
			"",
			"data: [DONE]",
			"",
	placeholder, "\n"))),
placeholder

	result, err := svc.handleStreamingResponseAnthropicAPIKeyPassthrough(context.Background(), resp, c, &Account{ID: 1placeholder, time.Now(), "claude-3-7-sonnet-20250219")
placeholder
	require.NotNil(t, result)
	require.NotNil(t, result.usage)
	require.Equal(t, 11, result.usage.InputTokens)
	require.Equal(t, 5, result.usage.OutputTokens)
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_MissingTerminalEventReturnsError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	svc := &GatewayService{
		cfg: &config.Config{
			Gateway: config.GatewayConfig{
				MaxLineSize: defaultMaxLineSize,
		placeholder,
	placeholder,
		rateLimitService: &RateLimitService{placeholder,
placeholder

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(strings.Join([]string{
			`data: {"type":"message_start","message":{"usage":{"input_tokens":11placeholderplaceholderplaceholder`,
			"",
			`data: {"type":"message_delta","usage":{"output_tokens":5placeholderplaceholder`,
			"",
	placeholder, "\n"))),
placeholder

	result, err := svc.handleStreamingResponseAnthropicAPIKeyPassthrough(context.Background(), resp, c, &Account{ID: 1placeholder, time.Now(), "claude-3-7-sonnet-20250219")
placeholder
	require.Contains(t, err.Error(), "missing terminal event")
	require.NotNil(t, result)
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_ForwardDirect_NonStreamingSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	body := []byte(`{"model":"claude-3-5-sonnet-latest","messages":[{"role":"user","content":[{"type":"text","text":"hello"placeholder]placeholder]placeholder`)
	upstreamJSON := `{"id":"msg_1","type":"message","usage":{"input_tokens":12,"output_tokens":7,"cache_creation":{"ephemeral_5m_input_tokens":2,"ephemeral_1h_input_tokens":3placeholder,"cached_tokens":4placeholderplaceholder`
	upstream := &anthropicHTTPUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header: http.Header{
				"Content-Type": []string{"application/json"placeholder,
				"x-request-id": []string{"rid-nonstream"placeholder,
		placeholder,
			Body: io.NopCloser(strings.NewReader(upstreamJSON)),
	placeholder,
placeholder
	svc := &GatewayService{
		cfg:              &config.Config{placeholder,
		httpUpstream:     upstream,
		rateLimitService: &RateLimitService{placeholder,
placeholder

	result, err := svc.forwardAnthropicAPIKeyPassthrough(context.Background(), c, newAnthropicAPIKeyAccountForTest(), body, "claude-3-5-sonnet-latest", "claude-3-5-sonnet-latest", false, time.Now())
placeholder
	require.NotNil(t, result)
	require.Equal(t, 12, result.Usage.InputTokens)
	require.Equal(t, 7, result.Usage.OutputTokens)
	require.Equal(t, 5, result.Usage.CacheCreationInputTokens)
	require.Equal(t, 4, result.Usage.CacheReadInputTokens)
	require.Equal(t, upstreamJSON, rec.Body.String())
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_ForwardDirect_InvalidTokenType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	account := &Account{
		ID:       202,
		Name:     "anthropic-oauth",
		Platform: PlatformAnthropic,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "oauth-token",
	placeholder,
placeholder
	svc := &GatewayService{placeholder

	result, err := svc.forwardAnthropicAPIKeyPassthrough(context.Background(), c, account, []byte(`{placeholder`), "claude-3-5-sonnet-latest", "claude-3-5-sonnet-latest", false, time.Now())
	require.Nil(t, result)
placeholder
	require.Contains(t, err.Error(), "requires apikey token")
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_ForwardDirect_UpstreamRequestError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	upstream := &anthropicHTTPUpstreamRecorder{
		err: errors.New("dial tcp timeout"),
placeholder
	svc := &GatewayService{
		cfg: &config.Config{
			Security: config.SecurityConfig{
				URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholder,
		placeholder,
	placeholder,
		httpUpstream: upstream,
placeholder
	account := newAnthropicAPIKeyAccountForTest()

	result, err := svc.forwardAnthropicAPIKeyPassthrough(context.Background(), c, account, []byte(`{"model":"x"placeholder`), "x", "x", false, time.Now())
	require.Nil(t, result)
placeholder
	require.Contains(t, err.Error(), "upstream request failed")
	require.Equal(t, http.StatusBadGateway, rec.Code)
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_ForwardDirect_EmptyResponseBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	upstream := &anthropicHTTPUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"x-request-id": []string{"rid-empty-body"placeholderplaceholder,
			Body:       nil,
	placeholder,
placeholder
	svc := &GatewayService{
		cfg: &config.Config{
			Security: config.SecurityConfig{
				URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholder,
		placeholder,
	placeholder,
		httpUpstream: upstream,
placeholder

	result, err := svc.forwardAnthropicAPIKeyPassthrough(context.Background(), c, newAnthropicAPIKeyAccountForTest(), []byte(`{"model":"x"placeholder`), "x", "x", false, time.Now())
	require.Nil(t, result)
placeholder
	require.Contains(t, err.Error(), "empty response")
placeholder

func TestExtractAnthropicSSEDataLine(t *testing.T) {
	t.Run("valid data line with spaces", func(t *testing.T) {
		data, ok := extractAnthropicSSEDataLine("data:   {\"type\":\"message_start\"placeholder")
		require.True(t, ok)
		require.Equal(t, `{"type":"message_start"placeholder`, data)
placeholder)

	t.Run("non data line", func(t *testing.T) {
		data, ok := extractAnthropicSSEDataLine("event: message_start")
		require.False(t, ok)
		require.Empty(t, data)
placeholder)
placeholder

func TestGatewayService_ParseSSEUsagePassthrough_MessageStartFallbacks(t *testing.T) {
	svc := &GatewayService{placeholder
	usage := &ClaudeUsage{placeholder
	data := `{"type":"message_start","message":{"usage":{"input_tokens":12,"cache_creation_input_tokens":0,"cache_read_input_tokens":0,"cached_tokens":9,"cache_creation":{"ephemeral_5m_input_tokens":3,"ephemeral_1h_input_tokens":4placeholderplaceholderplaceholderplaceholder`

	svc.parseSSEUsagePassthrough(data, usage)

	require.Equal(t, 12, usage.InputTokens)
	require.Equal(t, 9, usage.CacheReadInputTokens, "应兼容 cached_tokens 字段")
	require.Equal(t, 7, usage.CacheCreationInputTokens, "聚合字段为空时应从 5m/1h 明细回填")
	require.Equal(t, 3, usage.CacheCreation5mTokens)
	require.Equal(t, 4, usage.CacheCreation1hTokens)
placeholder

func TestGatewayService_ParseSSEUsagePassthrough_MessageDeltaSelectiveOverwrite(t *testing.T) {
	svc := &GatewayService{placeholder
	usage := &ClaudeUsage{
		InputTokens:           10,
		CacheCreation5mTokens: 2,
		CacheCreation1hTokens: 6,
placeholder
	data := `{"type":"message_delta","usage":{"input_tokens":0,"output_tokens":5,"cache_creation_input_tokens":8,"cache_read_input_tokens":0,"cached_tokens":11,"cache_creation":{"ephemeral_5m_input_tokens":1,"ephemeral_1h_input_tokens":0placeholderplaceholderplaceholder`

	svc.parseSSEUsagePassthrough(data, usage)

	require.Equal(t, 10, usage.InputTokens, "message_delta 中 0 值不应覆盖已有 input_tokens")
	require.Equal(t, 5, usage.OutputTokens)
	require.Equal(t, 8, usage.CacheCreationInputTokens)
	require.Equal(t, 11, usage.CacheReadInputTokens, "cache_read_input_tokens 为空时应回退到 cached_tokens")
	require.Equal(t, 1, usage.CacheCreation5mTokens)
	require.Equal(t, 6, usage.CacheCreation1hTokens, "message_delta 中 0 值不应覆盖已有 1h 明细")
placeholder

func TestGatewayService_ParseSSEUsagePassthrough_NoopCases(t *testing.T) {
	svc := &GatewayService{placeholder

	usage := &ClaudeUsage{InputTokens: 3placeholder
	svc.parseSSEUsagePassthrough("", usage)
	require.Equal(t, 3, usage.InputTokens)

	svc.parseSSEUsagePassthrough("[DONE]", usage)
	require.Equal(t, 3, usage.InputTokens)

	svc.parseSSEUsagePassthrough("not-json", usage)
	require.Equal(t, 3, usage.InputTokens)

	// nil usage 不应 panic
	svc.parseSSEUsagePassthrough(`{"type":"message_start"placeholder`, nil)
placeholder

func TestGatewayService_ParseSSEUsagePassthrough_FallbackFromUsageNode(t *testing.T) {
	svc := &GatewayService{placeholder
	usage := &ClaudeUsage{placeholder
	data := `{"type":"content_block_delta","usage":{"cached_tokens":6,"cache_creation":{"ephemeral_5m_input_tokens":2,"ephemeral_1h_input_tokens":1placeholderplaceholderplaceholder`

	svc.parseSSEUsagePassthrough(data, usage)

	require.Equal(t, 6, usage.CacheReadInputTokens)
	require.Equal(t, 3, usage.CacheCreationInputTokens)
placeholder

func TestParseClaudeUsageFromResponseBody(t *testing.T) {
	t.Run("empty or missing usage", func(t *testing.T) {
		got := parseClaudeUsageFromResponseBody(nil)
		require.NotNil(t, got)
		require.Equal(t, 0, got.InputTokens)

		got = parseClaudeUsageFromResponseBody([]byte(`{"id":"x"placeholder`))
		require.NotNil(t, got)
		require.Equal(t, 0, got.OutputTokens)
placeholder)

	t.Run("parse all usage fields and fallback", func(t *testing.T) {
		body := []byte(`{"usage":{"input_tokens":21,"output_tokens":34,"cache_creation_input_tokens":0,"cache_read_input_tokens":0,"cached_tokens":13,"cache_creation":{"ephemeral_5m_input_tokens":5,"ephemeral_1h_input_tokens":8placeholderplaceholderplaceholder`)
		got := parseClaudeUsageFromResponseBody(body)
		require.Equal(t, 21, got.InputTokens)
		require.Equal(t, 34, got.OutputTokens)
		require.Equal(t, 13, got.CacheReadInputTokens, "cache_read_input_tokens 为空时应回退 cached_tokens")
		require.Equal(t, 13, got.CacheCreationInputTokens, "聚合字段为空时应由 5m/1h 回填")
		require.Equal(t, 5, got.CacheCreation5mTokens)
		require.Equal(t, 8, got.CacheCreation1hTokens)
placeholder)

	t.Run("keep explicit aggregate values", func(t *testing.T) {
		body := []byte(`{"usage":{"input_tokens":1,"output_tokens":2,"cache_creation_input_tokens":9,"cache_read_input_tokens":7,"cached_tokens":99,"cache_creation":{"ephemeral_5m_input_tokens":4,"ephemeral_1h_input_tokens":5placeholderplaceholderplaceholder`)
		got := parseClaudeUsageFromResponseBody(body)
		require.Equal(t, 9, got.CacheCreationInputTokens, "已显式提供聚合字段时不应被明细覆盖")
		require.Equal(t, 7, got.CacheReadInputTokens, "已显式提供 cache_read_input_tokens 时不应回退 cached_tokens")
placeholder)
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_StreamingErrTooLong(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	svc := &GatewayService{
		cfg: &config.Config{
			Gateway: config.GatewayConfig{
				MaxLineSize: 32,
		placeholder,
	placeholder,
placeholder

	// Scanner 初始缓冲为 64KB，构造更长单行触发 bufio.ErrTooLong。
	longLine := "data: " + strings.Repeat("x", 80*1024)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(longLine)),
placeholder

	result, err := svc.handleStreamingResponseAnthropicAPIKeyPassthrough(context.Background(), resp, c, &Account{ID: 2placeholder, time.Now(), "claude-3-7-sonnet-20250219")
placeholder
	require.ErrorIs(t, err, bufio.ErrTooLong)
	require.NotNil(t, result)
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_StreamingDataIntervalTimeout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	svc := &GatewayService{
		cfg: &config.Config{
			Gateway: config.GatewayConfig{
				StreamDataIntervalTimeout: 1,
				MaxLineSize:               defaultMaxLineSize,
		placeholder,
	placeholder,
		rateLimitService: &RateLimitService{placeholder,
placeholder

	pr, pw := io.Pipe()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       pr,
placeholder

	result, err := svc.handleStreamingResponseAnthropicAPIKeyPassthrough(context.Background(), resp, c, &Account{ID: 5placeholder, time.Now(), "claude-3-7-sonnet-20250219")
	_ = pw.Close()
	_ = pr.Close()

placeholder
	require.Contains(t, err.Error(), "stream data interval timeout")
	require.NotNil(t, result)
	require.False(t, result.clientDisconnect)
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_StreamingSendsKeepaliveDuringIdle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	svc := &GatewayService{
		cfg: &config.Config{
			Gateway: config.GatewayConfig{
				StreamKeepaliveInterval: 1,
				MaxLineSize:             defaultMaxLineSize,
		placeholder,
	placeholder,
		rateLimitService: &RateLimitService{placeholder,
placeholder

	pr, pw := io.Pipe()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       pr,
placeholder

	done := make(chan struct{placeholder)
	go func() {
		defer close(done)
		time.Sleep(1200 * time.Millisecond)
		_, _ = pw.Write([]byte(strings.Join([]string{
			`data: {"type":"message_start","message":{"usage":{"input_tokens":3placeholderplaceholderplaceholder`,
			"",
			`data: {"type":"message_delta","usage":{"output_tokens":2placeholderplaceholder`,
			"",
			"data: [DONE]",
			"",
	placeholder, "\n")))
		_ = pw.Close()
placeholder()

	result, err := svc.handleStreamingResponseAnthropicAPIKeyPassthrough(context.Background(), resp, c, &Account{ID: 8placeholder, time.Now(), "claude-3-7-sonnet-20250219")
	_ = pr.Close()
	<-done

placeholder
	require.NotNil(t, result)
	require.Contains(t, rec.Body.String(), "event: ping\ndata: {\"type\": \"ping\"placeholder\n\n")
	require.Contains(t, rec.Body.String(), "data: [DONE]")
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_StreamingKeepaliveDoesNotInterleavePartialEvent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	svc := &GatewayService{
		cfg: &config.Config{
			Gateway: config.GatewayConfig{
				StreamKeepaliveInterval: 1,
				MaxLineSize:             defaultMaxLineSize,
		placeholder,
	placeholder,
		rateLimitService: &RateLimitService{placeholder,
placeholder

	pr, pw := io.Pipe()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       pr,
placeholder

	done := make(chan struct{placeholder)
	go func() {
		defer close(done)
		_, _ = pw.Write([]byte(`data: {"type":"message_start","message":{"usage":{"input_tokens":4placeholderplaceholderplaceholder` + "\n"))
		time.Sleep(1200 * time.Millisecond)
		_, _ = pw.Write([]byte("\n"))
		_, _ = pw.Write([]byte("data: [DONE]\n\n"))
		_ = pw.Close()
placeholder()

	result, err := svc.handleStreamingResponseAnthropicAPIKeyPassthrough(context.Background(), resp, c, &Account{ID: 9placeholder, time.Now(), "claude-3-7-sonnet-20250219")
	_ = pr.Close()
	<-done

placeholder
	require.NotNil(t, result)
	body := rec.Body.String()
	require.NotContains(t, body, `data: {"type":"message_start","message":{"usage":{"input_tokens":4placeholderplaceholderplaceholder`+"\n"+"event: ping")
	require.NotContains(t, body, "event: ping")
	require.Contains(t, body, "data: [DONE]")
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_StreamingReadError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	svc := &GatewayService{
		cfg: &config.Config{
			Gateway: config.GatewayConfig{
				MaxLineSize: defaultMaxLineSize,
		placeholder,
	placeholder,
placeholder

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body: &streamReadCloser{
			err: io.ErrUnexpectedEOF,
	placeholder,
placeholder

	result, err := svc.handleStreamingResponseAnthropicAPIKeyPassthrough(context.Background(), resp, c, &Account{ID: 6placeholder, time.Now(), "claude-3-7-sonnet-20250219")
placeholder
	require.Contains(t, err.Error(), "stream read error")
	require.NotNil(t, result)
	require.False(t, result.clientDisconnect)
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_StreamingTimeoutAfterClientDisconnect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
	c.Writer = &failWriteResponseWriter{ResponseWriter: c.Writerplaceholder

	svc := &GatewayService{
		cfg: &config.Config{
			Gateway: config.GatewayConfig{
				StreamDataIntervalTimeout: 1,
				MaxLineSize:               defaultMaxLineSize,
		placeholder,
	placeholder,
		rateLimitService: &RateLimitService{placeholder,
placeholder

	pr, pw := io.Pipe()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       pr,
placeholder

	done := make(chan struct{placeholder)
	go func() {
		defer close(done)
		_, _ = pw.Write([]byte(`data: {"type":"message_start","message":{"usage":{"input_tokens":9placeholderplaceholderplaceholder` + "\n"))
		// 保持上游连接静默，触发数据间隔超时分支。
		time.Sleep(1500 * time.Millisecond)
		_ = pw.Close()
placeholder()

	result, err := svc.handleStreamingResponseAnthropicAPIKeyPassthrough(context.Background(), resp, c, &Account{ID: 7placeholder, time.Now(), "claude-3-7-sonnet-20250219")
	_ = pr.Close()
	<-done

placeholder
	require.Contains(t, err.Error(), "stream usage incomplete after timeout")
	require.NotNil(t, result)
	require.True(t, result.clientDisconnect)
	require.Equal(t, 9, result.usage.InputTokens)
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_StreamingContextCanceled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	svc := &GatewayService{
		cfg: &config.Config{
			Gateway: config.GatewayConfig{
				MaxLineSize: defaultMaxLineSize,
		placeholder,
	placeholder,
placeholder

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body: &streamReadCloser{
			err: context.Canceled,
	placeholder,
placeholder

	result, err := svc.handleStreamingResponseAnthropicAPIKeyPassthrough(context.Background(), resp, c, &Account{ID: 3placeholder, time.Now(), "claude-3-7-sonnet-20250219")
placeholder
	require.Contains(t, err.Error(), "stream usage incomplete")
	require.NotNil(t, result)
	require.True(t, result.clientDisconnect)
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_StreamingUpstreamReadErrorAfterClientDisconnect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
	c.Writer = &failWriteResponseWriter{ResponseWriter: c.Writerplaceholder

	svc := &GatewayService{
		cfg: &config.Config{
			Gateway: config.GatewayConfig{
				MaxLineSize: defaultMaxLineSize,
		placeholder,
	placeholder,
placeholder

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body: &streamReadCloser{
			payload: []byte(`data: {"type":"message_start","message":{"usage":{"input_tokens":8placeholderplaceholderplaceholder` + "\n\n"),
			err:     io.ErrUnexpectedEOF,
	placeholder,
placeholder

	result, err := svc.handleStreamingResponseAnthropicAPIKeyPassthrough(context.Background(), resp, c, &Account{ID: 4placeholder, time.Now(), "claude-3-7-sonnet-20250219")
placeholder
	require.Contains(t, err.Error(), "stream usage incomplete after disconnect")
	require.NotNil(t, result)
	require.True(t, result.clientDisconnect)
	require.Equal(t, 8, result.usage.InputTokens)
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_TransportErrorRecordsOllamaActivity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	deferred := NewDeferredService(nil, nil, time.Second)
	upstream := &anthropicHTTPUpstreamRecorder{err: errors.New("dial tcp timeout")placeholder
	svc := &GatewayService{
		cfg: &config.Config{
			Security: config.SecurityConfig{
				URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholder,
		placeholder,
	placeholder,
		httpUpstream:    upstream,
		deferredService: deferred,
placeholder

	ollama := &Account{
		ID: 601, Name: "ollama-anthropic", Platform: PlatformAnthropic, Type: AccountTypeAPIKey,
		Concurrency: 1,
placeholder"api_key": "k-ollama", "base_url": "https://ollama.com"placeholder,
		Extra:       map[string]any{"anthropic_passthrough": trueplaceholder,
		Status:      StatusActive, Schedulable: true,
placeholder
	other := newAnthropicAPIKeyAccountForTest()
	other.ID = 602

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
	_, err := svc.forwardAnthropicAPIKeyPassthrough(context.Background(), c, ollama, []byte(`{"model":"x"placeholder`), "x", "x", false, time.Now())
placeholder

	rec2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(rec2)
	c2.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
	_, err = svc.forwardAnthropicAPIKeyPassthrough(context.Background(), c2, other, []byte(`{"model":"x"placeholder`), "x", "x", false, time.Now())
placeholder

	_, ok := deferred.lastUsedUpdates.Load(int64(601))
	require.True(t, ok, "Anthropic passthrough transport error on Ollama account must record activity")
	_, ok = deferred.lastUsedUpdates.Load(int64(602))
	require.False(t, ok, "non-Ollama Anthropic passthrough transport error must not record Ollama activity")
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_ContextCanceledSkipsOllamaActivity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	deferred := NewDeferredService(nil, nil, time.Second)
	upstream := &anthropicHTTPUpstreamRecorder{err: context.Canceledplaceholder
	svc := &GatewayService{
		cfg: &config.Config{
			Security: config.SecurityConfig{
				URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholder,
		placeholder,
	placeholder,
		httpUpstream:    upstream,
		deferredService: deferred,
placeholder
	ollama := &Account{
		ID: 603, Name: "ollama-canceled", Platform: PlatformAnthropic, Type: AccountTypeAPIKey,
		Concurrency: 1,
placeholder"api_key": "k-ollama", "base_url": "https://ollama.com"placeholder,
		Extra:       map[string]any{"anthropic_passthrough": trueplaceholder,
		Status:      StatusActive, Schedulable: true,
placeholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	_, err := svc.forwardAnthropicAPIKeyPassthrough(context.Background(), c, ollama, []byte(`{"model":"x"placeholder`), "x", "x", false, time.Now())

placeholder
	_, ok := deferred.lastUsedUpdates.Load(int64(603))
	require.False(t, ok, "context.Canceled on Anthropic passthrough must not count as Ollama activity")
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_Non2xxRecordsOllamaActivity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	deferred := NewDeferredService(nil, nil, time.Second)
	// 400 is non-retryable / non-failover for default API-key accounts, so it reaches handleErrorResponse.
	upstream := &anthropicHTTPUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusBadRequest,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"type":"error","error":{"type":"invalid_request_error","message":"bad"placeholderplaceholder`)),
	placeholder,
placeholder
	svc := &GatewayService{
		cfg: &config.Config{
			Security: config.SecurityConfig{
				URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholder,
		placeholder,
	placeholder,
		httpUpstream:     upstream,
		deferredService:  deferred,
		rateLimitService: &RateLimitService{placeholder,
placeholder
	ollama := &Account{
		ID: 604, Name: "ollama-400", Platform: PlatformAnthropic, Type: AccountTypeAPIKey,
		Concurrency: 1,
placeholder"api_key": "k-ollama", "base_url": "https://ollama.com"placeholder,
		Extra:       map[string]any{"anthropic_passthrough": trueplaceholder,
		Status:      StatusActive, Schedulable: true,
placeholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	_, _ = svc.forwardAnthropicAPIKeyPassthrough(context.Background(), c, ollama, []byte(`{"model":"x"placeholder`), "x", "x", false, time.Now())

	_, ok := deferred.lastUsedUpdates.Load(int64(604))
	require.True(t, ok, "Anthropic passthrough non-2xx on Ollama account must record activity via handleErrorResponse")
placeholder
