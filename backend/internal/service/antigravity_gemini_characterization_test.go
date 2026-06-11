//go:build unit

// Phase-0 TASK-002 特征化测试：gemini/antigravity 路径流式/非流式透传（INVARIANTS I-1.6）。
// ForwardGemini 的外部可观测语义：
//   - 流式：上游 v1internal SSE（data: {"response":{...placeholderplaceholder）被解包后逐事件转发给客户端；
//   - 非流式：上游流式响应被收集合并为单个 JSON（文本片段拼接）后返回；
//   - 非 failover 上游错误（如 404）：解包后的错误体 + 上游状态码原样返回。
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// passCharAntigravityService 构造最小可运行的 AntigravityGatewayService。
func passCharAntigravityService(upstream HTTPUpstream) *AntigravityGatewayService {
	return &AntigravityGatewayService{
		settingService: NewSettingService(&antigravitySettingRepoStub{placeholder, &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder),
		tokenProvider:  &AntigravityTokenProvider{placeholder,
		httpUpstream:   upstream,
placeholder
placeholder

// passCharAntigravityAccount 返回带模型映射的 antigravity OAuth 账号夹具。
func passCharAntigravityAccount(id int64, mapping map[string]any) *Account {
placeholder
		ID:          id,
		Name:        "pass-char-antigravity",
		Platform:    PlatformAntigravity,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
		Schedulable: true,
placeholder
			"access_token":  "ag-token",
			"model_mapping": mapping,
	placeholder,
placeholder
placeholder

func passCharGeminiBody(t *testing.T) []byte {
placeholder
	body, err := json.Marshal(map[string]any{
		"contents": []map[string]any{
			{"role": "user", "parts": []map[string]any{{"text": "hello"placeholderplaceholderplaceholder,
	placeholder,
placeholder)
placeholder
	return body
placeholder

// TestGatewayCharacterization_GeminiStreamUnwrapsV1Internal 固化 I-1.6 流式侧：
// 上游 data 行的 v1internal 包裹（{"response":{...placeholderplaceholder）被解包，客户端按事件收到内层 JSON 原文。
func TestGatewayCharacterization_GeminiStreamUnwrapsV1Internal(t *testing.T) {
	gin.SetMode(gin.TestMode)

	innerChunk1 := `{"candidates":[{"content":{"role":"model","parts":[{"text":"Hel"placeholder]placeholderplaceholder]placeholder`
	innerChunk2 := `{"candidates":[{"content":{"role":"model","parts":[{"text":"lo ✓"placeholder]placeholder,"finishReason":"STOP"placeholder],"usageMetadata":{"promptTokenCount":8,"candidatesTokenCount":3,"cachedContentTokenCount":2placeholderplaceholder`
	upstreamSSE := "data: {\"response\":" + innerChunk1 + "placeholder\n\n" +
		"data: {\"response\":" + innerChunk2 + "placeholder\n\n"

	upstream := &httpUpstreamStub{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"text/event-stream"placeholder,
			"X-Request-Id": []string{"rid-char-gemini-stream"placeholder,
	placeholder,
		Body: io.NopCloser(bytes.NewReader([]byte(upstreamSSE))),
placeholderplaceholder
	svc := passCharAntigravityService(upstream)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1beta/models/gemini-2.5-flash:streamGenerateContent", bytes.NewReader(passCharGeminiBody(t)))

	account := passCharAntigravityAccount(9201, map[string]any{"gemini-2.5-flash": "gemini-3-pro-high"placeholder)

	result, err := svc.ForwardGemini(context.Background(), c, account, "gemini-2.5-flash", "streamGenerateContent", true, passCharGeminiBody(t), false)
placeholder
	require.NotNil(t, result)
	require.True(t, result.Stream)

	wantEvents := []string{
		"data: " + innerChunk1,
		"data: " + innerChunk2,
placeholder
	require.Equal(t, wantEvents, passCharSplitSSEEvents(rec.Body.String()),
		"客户端应按事件收到解包后的内层 JSON 原文")

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "text/event-stream", rec.Header().Get("Content-Type"))
	require.Equal(t, "rid-char-gemini-stream", rec.Header().Get("x-request-id"))

	require.Equal(t, "gemini-2.5-flash", result.Model)
	require.Equal(t, "gemini-3-pro-high", result.UpstreamModel, "计费模型使用映射后的模型")
	require.Equal(t, 6, result.Usage.InputTokens, "input = promptTokenCount - cachedContentTokenCount")
	require.Equal(t, 3, result.Usage.OutputTokens)
	require.Equal(t, 2, result.Usage.CacheReadInputTokens)
placeholder

// TestGatewayCharacterization_GeminiNonStreamCollectsStream 固化 I-1.6 非流式侧：
// 客户端请求非流式时，网关收集上游流式 chunk，将文本片段拼接进最后一个含 parts 的
// 响应中，作为单个 JSON（HTTP 200, application/json）返回。
func TestGatewayCharacterization_GeminiNonStreamCollectsStream(t *testing.T) {
	gin.SetMode(gin.TestMode)

	upstreamSSE := `data: {"response":{"candidates":[{"content":{"role":"model","parts":[{"text":"Hel"placeholder]placeholderplaceholder]placeholderplaceholder` + "\n\n" +
		`data: {"response":{"candidates":[{"content":{"role":"model","parts":[{"text":"lo ✓"placeholder]placeholder,"finishReason":"STOP"placeholder],"usageMetadata":{"promptTokenCount":8,"candidatesTokenCount":3,"cachedContentTokenCount":2placeholderplaceholderplaceholder` + "\n\n"

	upstream := &httpUpstreamStub{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"text/event-stream"placeholder,
			"X-Request-Id": []string{"rid-char-gemini-collect"placeholder,
	placeholder,
		Body: io.NopCloser(bytes.NewReader([]byte(upstreamSSE))),
placeholderplaceholder
	svc := passCharAntigravityService(upstream)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1beta/models/gemini-2.5-flash:generateContent", bytes.NewReader(passCharGeminiBody(t)))

	account := passCharAntigravityAccount(9202, map[string]any{"gemini-2.5-flash": "gemini-3-pro-high"placeholder)

	result, err := svc.ForwardGemini(context.Background(), c, account, "gemini-2.5-flash", "generateContent", false, passCharGeminiBody(t), false)
placeholder
	require.NotNil(t, result)
	require.False(t, result.Stream)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	require.Equal(t, "rid-char-gemini-collect", rec.Header().Get("x-request-id"))

	// 合并语义：以最后一个含 parts 的 chunk 为基底，第一个 text part 替换为全部文本片段的拼接。
	require.JSONEq(t, `{
		"candidates":[{
			"content":{"role":"model","parts":[{"text":"Hello ✓"placeholder]placeholder,
			"finishReason":"STOP"
	placeholder],
		"usageMetadata":{"promptTokenCount":8,"candidatesTokenCount":3,"cachedContentTokenCount":2placeholder
placeholder`, rec.Body.String())

	require.Equal(t, 6, result.Usage.InputTokens)
	require.Equal(t, 3, result.Usage.OutputTokens)
	require.Equal(t, 2, result.Usage.CacheReadInputTokens)
placeholder

// TestGatewayCharacterization_GeminiUpstreamErrorUnwrappedPassthrough 固化 I-1.6 错误侧：
// 非 failover 上游错误（404）时，客户端收到上游状态码 + 解包后的错误体原文。
func TestGatewayCharacterization_GeminiUpstreamErrorUnwrappedPassthrough(t *testing.T) {
	gin.SetMode(gin.TestMode)

	innerErr := `{"error":{"code":404,"message":"model not found: gemini-3-pro-high","status":"NOT_FOUND"placeholderplaceholder`
	upstream := &httpUpstreamStub{resp: &http.Response{
		StatusCode: http.StatusNotFound,
		Header: http.Header{
			"Content-Type": []string{"application/json"placeholder,
			"X-Request-Id": []string{"rid-char-gemini-err"placeholder,
	placeholder,
		Body: io.NopCloser(bytes.NewReader([]byte(`{"response":` + innerErr + `placeholder`))),
placeholderplaceholder
	svc := passCharAntigravityService(upstream)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1beta/models/gemini-2.5-flash:generateContent", bytes.NewReader(passCharGeminiBody(t)))

	account := passCharAntigravityAccount(9203, map[string]any{"gemini-2.5-flash": "gemini-3-pro-high"placeholder)

	result, err := svc.ForwardGemini(context.Background(), c, account, "gemini-2.5-flash", "generateContent", false, passCharGeminiBody(t), false)
placeholder
	require.Nil(t, result)

	require.Equal(t, http.StatusNotFound, rec.Code, "非 failover 上游错误状态码原样返回")
	require.Equal(t, innerErr, rec.Body.String(), "错误体应为解包后的内层 JSON 原文")
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	require.Equal(t, "rid-char-gemini-err", rec.Header().Get("x-request-id"))
placeholder
