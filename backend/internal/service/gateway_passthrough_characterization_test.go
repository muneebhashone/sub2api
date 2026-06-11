//go:build unit

// 本文件是 Phase-0 TASK-002 的特征化测试（characterization tests）。
// 它把"客户端最终收到的内容"固化为基线：固定 mock 上游，逐字节/逐事件断言
// 状态码、响应体与响应头。任何后续重构若改变这些外部可观测行为，测试会立即失败。
//
// 覆盖的不变量（见 .claude/plugin-refactor/INVARIANTS.md）：
//   - I-1.1 anthropic 非流式 2xx body/状态码逐字节透传（升级自 JSONEq 级断言）
//   - I-1.4 响应 header 白名单过滤（默认白名单 / additional_allowed / force_remove）
//   - I-1.7 上游 4xx/5xx 错误按 ErrorPassthroughService 规则透传或转换
package service

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// passCharSplitSSEEvents 按 SSE 事件边界（空行 "\n\n"）分帧。
// 特征化测试只按事件比对，不按网络 chunk 比对；帧首尾多余的换行
// （合法的 SSE 分帧差异，如额外空行）会被剥离，但帧内内容保持原文。
func passCharSplitSSEEvents(body string) []string {
	frames := strings.Split(body, "\n\n")
	out := make([]string, 0, len(frames))
	for _, f := range frames {
		f = strings.Trim(f, "\r\n")
		if f == "" {
			continue
	placeholder
		out = append(out, f)
placeholder
	return out
placeholder

// passCharAnthropicAccount 返回 anthropic API Key 透传账号夹具。
func passCharAnthropicAccount(id int64) *Account {
placeholder
		ID:          id,
		Name:        "pass-char-anthropic",
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

// passCharGinContext 构造带 /v1/messages POST 请求的 gin 测试上下文。
func passCharGinContext(t *testing.T, path string) (*gin.Context, *httptest.ResponseRecorder) {
placeholder
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, path, nil)
	return c, rec
placeholder

func passCharGatewayService(cfg *config.Config, upstream HTTPUpstream) *GatewayService {
	return &GatewayService{
		cfg:                  cfg,
		responseHeaderFilter: compileResponseHeaderFilter(cfg),
		httpUpstream:         upstream,
		rateLimitService:     &RateLimitService{placeholder,
		deferredService:      &DeferredService{placeholder,
placeholder
placeholder

// TestGatewayCharacterization_AnthropicNonStreamPassthroughByteExact 固化 I-1.1：
// anthropic 透传账号非流式转发时，上游 2xx 响应体逐字节透传、状态码透传。
// 上游 body 故意包含字段顺序、多余空白、Unicode、HTML 字符与未知字段，
// 任何 JSON 重新序列化都会破坏逐字节相等。
func TestGatewayCharacterization_AnthropicNonStreamPassthroughByteExact(t *testing.T) {
	upstreamJSON := "{\"id\":\"msg_char_01\",  \"type\":\"message\",\n" +
		"\"unknown_field\":{\"nested\":[1,2.50,\"<&>\"]placeholder,\n" +
		"\"content\":[{\"type\":\"text\",\"text\":\"héllo ✓ <b>&amp;</b>\"placeholder],\n" +
		"\"usage\":{\"input_tokens\":12,\"output_tokens\":7placeholderplaceholder"

	for _, status := range []int{http.StatusOK, http.StatusAcceptedplaceholder {
		upstream := &anthropicHTTPUpstreamRecorder{
			resp: &http.Response{
				StatusCode: status,
				Header: http.Header{
					"Content-Type": []string{"application/json; charset=utf-8"placeholder,
					"x-request-id": []string{"rid-char-nonstream"placeholder,
			placeholder,
				Body: io.NopCloser(strings.NewReader(upstreamJSON)),
		placeholder,
	placeholder
		cfg := &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder
		svc := passCharGatewayService(cfg, upstream)

		c, rec := passCharGinContext(t, "/v1/messages")
		parsed := &ParsedRequest{
			Body:  NewRequestBodyRef([]byte(`{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"placeholder]placeholder`)),
			Model: "claude-sonnet-4-20250514",
	placeholder

		result, err := svc.Forward(context.Background(), c, passCharAnthropicAccount(9001), parsed)
	placeholder
		require.NotNil(t, result)

		require.Equal(t, status, rec.Code, "上游 2xx 状态码必须透传")
		require.Equal(t, upstreamJSON, rec.Body.String(), "上游 2xx body 必须逐字节透传")
		require.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"), "Content-Type 必须透传")
		require.Equal(t, "rid-char-nonstream", rec.Header().Get("x-request-id"), "x-request-id 必须透传")
		require.Equal(t, 12, result.Usage.InputTokens)
		require.Equal(t, 7, result.Usage.OutputTokens)
placeholder
placeholder

// TestGatewayCharacterization_AnthropicStreamPassthroughEventSequence 固化 I-1.1/I-1.2 流式侧：
// anthropic 透传账号流式转发时，客户端收到的 SSE 事件序列（按 \n\n 分帧）与上游完全一致，
// 含 event: 行、usage 事件与 cached_tokens 兼容字段，且网关不改写事件内容。
func TestGatewayCharacterization_AnthropicStreamPassthroughEventSequence(t *testing.T) {
	upstreamEvents := []string{
		"event: message_start\ndata: {\"type\":\"message_start\",\"message\":{\"id\":\"msg_char_02\",\"usage\":{\"input_tokens\":9,\"cached_tokens\":3placeholderplaceholderplaceholder",
		"event: content_block_delta\ndata: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"text_delta\",\"text\":\"Hel\"placeholderplaceholder",
		"event: content_block_delta\ndata: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"text_delta\",\"text\":\"lo ✓\"placeholderplaceholder",
		"event: message_delta\ndata: {\"type\":\"message_delta\",\"delta\":{\"stop_reason\":\"end_turn\"placeholder,\"usage\":{\"output_tokens\":5placeholderplaceholder",
		"event: message_stop\ndata: {\"type\":\"message_stop\"placeholder",
placeholder
	upstreamSSE := strings.Join(upstreamEvents, "\n\n") + "\n\n"

	upstream := &anthropicHTTPUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header: http.Header{
				"Content-Type": []string{"text/event-stream"placeholder,
				"x-request-id": []string{"rid-char-stream"placeholder,
		placeholder,
			Body: io.NopCloser(strings.NewReader(upstreamSSE)),
	placeholder,
placeholder
	cfg := &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder
	svc := passCharGatewayService(cfg, upstream)

	c, rec := passCharGinContext(t, "/v1/messages")
	parsed := &ParsedRequest{
		Body:   NewRequestBodyRef([]byte(`{"model":"claude-sonnet-4-20250514","stream":true,"messages":[{"role":"user","content":"hi"placeholder]placeholder`)),
		Model:  "claude-sonnet-4-20250514",
		Stream: true,
placeholder

	result, err := svc.Forward(context.Background(), c, passCharAnthropicAccount(9002), parsed)
placeholder
	require.NotNil(t, result)
	require.True(t, result.Stream)

	gotEvents := passCharSplitSSEEvents(rec.Body.String())
	require.Equal(t, upstreamEvents, gotEvents, "客户端收到的 SSE 事件序列必须与上游逐事件一致")

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "text/event-stream", rec.Header().Get("Content-Type"))
	require.Equal(t, "rid-char-stream", rec.Header().Get("x-request-id"))
	require.Equal(t, "no-cache", rec.Header().Get("Cache-Control"))
	require.Equal(t, "no", rec.Header().Get("X-Accel-Buffering"))

	require.Equal(t, 9, result.Usage.InputTokens)
	require.Equal(t, 5, result.Usage.OutputTokens)
	require.Equal(t, 3, result.Usage.CacheReadInputTokens, "cached_tokens 应被解析进计费 usage")
placeholder

// TestGatewayCharacterization_ResponseHeaderFilter 固化 I-1.4：
// 响应 header 白名单过滤语义。
//   - 关闭（security.response_headers.enabled=false，默认）：只用内置白名单，
//     additional_allowed / force_remove 均不生效；
//   - 开启：additional_allowed 追加放行，force_remove 即使是默认白名单键也强制移除。
func TestGatewayCharacterization_ResponseHeaderFilter(t *testing.T) {
	upstreamHeaders := func() http.Header {
		return http.Header{
			"Content-Type":                   []string{"application/json"placeholder,
			"Cache-Control":                  []string{"no-store"placeholder,
			"x-request-id":                   []string{"rid-char-headers"placeholder,
			"X-Ratelimit-Remaining-Requests": []string{"99"placeholder,
			"Retry-After":                    []string{"17"placeholder,
			"Www-Authenticate":               []string{"Bearer realm=api"placeholder,
			"Set-Cookie":                     []string{"secret=upstream"placeholder,
			"X-Custom-Upstream":              []string{"custom-value"placeholder,
	placeholder
placeholder

	tests := []struct {
		name           string
		headersCfg     config.ResponseHeaderConfig
		wantPresent    map[string]string
		wantAbsentKeys []string
placeholder{
		{
			name: "默认关闭_只用内置白名单且additional与force_remove不生效",
			headersCfg: config.ResponseHeaderConfig{
				Enabled:           false,
				AdditionalAllowed: []string{"x-custom-upstream"placeholder, // 关闭时不生效
				ForceRemove:       []string{"x-request-id"placeholder,      // 关闭时不生效
		placeholder,
			wantPresent: map[string]string{
				"Content-Type":                   "application/json",
				"Cache-Control":                  "no-store",
				"x-request-id":                   "rid-char-headers",
				"X-Ratelimit-Remaining-Requests": "99",
				"Retry-After":                    "17",
				"Www-Authenticate":               "Bearer realm=api",
		placeholder,
			wantAbsentKeys: []string{"Set-Cookie", "X-Custom-Upstream"placeholder,
	placeholder,
		{
			name: "开启_additional追加放行且force_remove覆盖默认白名单",
			headersCfg: config.ResponseHeaderConfig{
				Enabled:           true,
				AdditionalAllowed: []string{"x-custom-upstream"placeholder,
				ForceRemove:       []string{"x-request-id"placeholder,
		placeholder,
			wantPresent: map[string]string{
				"Content-Type":      "application/json",
				"X-Custom-Upstream": "custom-value",
				"Retry-After":       "17",
		placeholder,
			wantAbsentKeys: []string{"Set-Cookie", "x-request-id"placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upstream := &anthropicHTTPUpstreamRecorder{
				resp: &http.Response{
					StatusCode: http.StatusOK,
					Header:     upstreamHeaders(),
					Body:       io.NopCloser(strings.NewReader(`{"id":"msg_h","usage":{"input_tokens":1,"output_tokens":1placeholderplaceholder`)),
			placeholder,
		placeholder
			cfg := &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder
			cfg.Security.ResponseHeaders = tt.headersCfg
			svc := passCharGatewayService(cfg, upstream)

			c, rec := passCharGinContext(t, "/v1/messages")
			parsed := &ParsedRequest{
				Body:  NewRequestBodyRef([]byte(`{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"placeholder]placeholder`)),
				Model: "claude-sonnet-4-20250514",
		placeholder

			_, err := svc.Forward(context.Background(), c, passCharAnthropicAccount(9003), parsed)
		placeholder
			require.Equal(t, http.StatusOK, rec.Code)

			for key, want := range tt.wantPresent {
				require.Equal(t, want, rec.Header().Get(key), "响应头 %s 应保留", key)
		placeholder
			for _, key := range tt.wantAbsentKeys {
				require.Empty(t, rec.Header().Get(key), "响应头 %s 应被过滤", key)
		placeholder
	placeholder)
placeholder
placeholder

func passCharIntPtr(v int) *int       { return &v placeholder
func passCharStrPtr(v string) *string { return &v placeholder

// passCharErrorPassthroughService 构造带固定规则的错误透传服务（不连缓存）。
func passCharErrorPassthroughService(rules []*model.ErrorPassthroughRule) *ErrorPassthroughService {
	return NewErrorPassthroughService(&mockErrorPassthroughRepo{rules: rulesplaceholder, nil)
placeholder

// TestGatewayCharacterization_UpstreamErrorPassthroughRules 固化 I-1.7：
// 上游非 failover 错误（如 404）经 ErrorPassthroughService.MatchRule 决定透传或转换：
//   - 规则命中 + passthrough_code/passthrough_body：透传上游状态码 + 提取的上游 message；
//   - 规则命中 + 自定义 response_code/custom_message：使用自定义值；
//   - 规则未命中：默认 502 + "Upstream request failed"。
//
// 三种情况下错误体均为 {"type":"error","error":{"type":"upstream_error","message":...placeholderplaceholder。
func TestGatewayCharacterization_UpstreamErrorPassthroughRules(t *testing.T) {
	passthroughRule := &model.ErrorPassthroughRule{
		ID:              1,
		Name:            "404-keyword-passthrough",
		Enabled:         true,
		Priority:        1,
		ErrorCodes:      []int{404placeholder,
		Keywords:        []string{"quota_exceeded_marker"placeholder,
		MatchMode:       model.MatchModeAll,
		Platforms:       []string{model.PlatformAnthropicplaceholder,
		PassthroughCode: true,
		PassthroughBody: true,
placeholder
	overrideRule := &model.ErrorPassthroughRule{
		ID:              2,
		Name:            "404-keyword-override",
		Enabled:         true,
		Priority:        1,
		ErrorCodes:      []int{404placeholder,
		Keywords:        []string{"quota_exceeded_marker"placeholder,
		MatchMode:       model.MatchModeAll,
		Platforms:       []string{model.PlatformAnthropicplaceholder,
		PassthroughCode: false,
		ResponseCode:    passCharIntPtr(http.StatusTooManyRequests),
		PassthroughBody: false,
		CustomMessage:   passCharStrPtr("custom upstream busy"),
placeholder

	upstreamHitBody := `{"type":"error","error":{"type":"not_found_error","message":"quota_exceeded_marker: model not found"placeholderplaceholder`
	upstreamMissBody := `{"type":"error","error":{"type":"not_found_error","message":"plain not found"placeholderplaceholder`

	tests := []struct {
		name         string
		rules        []*model.ErrorPassthroughRule
		upstreamBody string
		wantStatus   int
		wantMessage  string
placeholder{
		{
			name:         "规则命中_透传上游状态码与消息",
			rules:        []*model.ErrorPassthroughRule{passthroughRuleplaceholder,
			upstreamBody: upstreamHitBody,
			wantStatus:   http.StatusNotFound,
			wantMessage:  "quota_exceeded_marker: model not found",
	placeholder,
		{
			name:         "规则命中_自定义状态码与消息",
			rules:        []*model.ErrorPassthroughRule{overrideRuleplaceholder,
			upstreamBody: upstreamHitBody,
			wantStatus:   http.StatusTooManyRequests,
			wantMessage:  "custom upstream busy",
	placeholder,
		{
			name:         "规则未命中_默认502转换",
			rules:        []*model.ErrorPassthroughRule{passthroughRuleplaceholder,
			upstreamBody: upstreamMissBody,
			wantStatus:   http.StatusBadGateway,
			wantMessage:  "Upstream request failed",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upstream := &anthropicHTTPUpstreamRecorder{
				resp: &http.Response{
					StatusCode: http.StatusNotFound,
					Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
					Body:       io.NopCloser(strings.NewReader(tt.upstreamBody)),
			placeholder,
		placeholder
			cfg := &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder
			svc := &GatewayService{
				cfg:                  cfg,
				responseHeaderFilter: compileResponseHeaderFilter(cfg),
				httpUpstream:         upstream,
		placeholder

			c, rec := passCharGinContext(t, "/v1/messages")
			BindErrorPassthroughService(c, passCharErrorPassthroughService(tt.rules))

			parsed := &ParsedRequest{
				Body:  NewRequestBodyRef([]byte(`{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"placeholder]placeholder`)),
				Model: "claude-sonnet-4-20250514",
		placeholder

			result, err := svc.Forward(context.Background(), c, passCharAnthropicAccount(9004), parsed)
			require.Error(t, err, "上游错误时 Forward 必须返回 error（响应已写给客户端）")
			require.Nil(t, result)

			require.Equal(t, tt.wantStatus, rec.Code)
			require.JSONEq(t,
				`{"type":"error","error":{"type":"upstream_error","message":"`+tt.wantMessage+`"placeholderplaceholder`,
				rec.Body.String())
	placeholder)
placeholder
placeholder

// TestGatewayCharacterization_Upstream400BodyPassthrough 固化 handleErrorResponse 对 400 的
// 特殊语义：未命中透传规则时，上游 400 响应体原样透传给客户端（状态码 400 + 原 body）。
func TestGatewayCharacterization_Upstream400BodyPassthrough(t *testing.T) {
	upstreamBody := `{"type":"error","error":{"type":"invalid_request_error","message":"max_tokens: required"placeholderplaceholder`
	upstream := &anthropicHTTPUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusBadRequest,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(upstreamBody)),
	placeholder,
placeholder
	cfg := &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder
	svc := &GatewayService{
		cfg:                  cfg,
		responseHeaderFilter: compileResponseHeaderFilter(cfg),
		httpUpstream:         upstream,
placeholder

	c, rec := passCharGinContext(t, "/v1/messages")
	parsed := &ParsedRequest{
		Body:  NewRequestBodyRef([]byte(`{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"placeholder]placeholder`)),
		Model: "claude-sonnet-4-20250514",
placeholder

	result, err := svc.Forward(context.Background(), c, passCharAnthropicAccount(9005), parsed)
placeholder
	require.Nil(t, result)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, upstreamBody, rec.Body.String(), "上游 400 body 应原样透传")
placeholder
