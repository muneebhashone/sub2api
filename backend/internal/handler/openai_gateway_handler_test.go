package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	pkghttputil "github.com/Wei-Shaw/sub2api/internal/pkg/httputil"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	coderws "github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func TestOpenAIHandleStreamingAwareError_JSONEscaping(t *testing.T) {
	tests := []struct {
		name    string
		errType string
		message string
placeholder{
		{
			name:    "包含双引号的消息",
			errType: "server_error",
			message: `upstream returned "invalid" response`,
	placeholder,
		{
			name:    "包含反斜杠的消息",
			errType: "server_error",
			message: `path C:\Users\test\file.txt not found`,
	placeholder,
		{
			name:    "包含双引号和反斜杠的消息",
			errType: "upstream_error",
			message: `error parsing "key\value": unexpected token`,
	placeholder,
		{
			name:    "包含换行符的消息",
			errType: "server_error",
			message: "line1\nline2\ttab",
	placeholder,
		{
			name:    "普通消息",
			errType: "upstream_error",
			message: "Upstream service temporarily unavailable",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

			h := &OpenAIGatewayHandler{placeholder
			h.handleStreamingAwareError(c, http.StatusBadGateway, tt.errType, tt.message, true)

			body := w.Body.String()

			// 验证 SSE 格式：event: error\ndata: {JSONplaceholder\n\n
			assert.True(t, strings.HasPrefix(body, "event: error\n"), "应以 'event: error\\n' 开头")
			assert.True(t, strings.HasSuffix(body, "\n\n"), "应以 '\\n\\n' 结尾")

			// 提取 data 部分
			lines := strings.Split(strings.TrimSuffix(body, "\n\n"), "\n")
			require.Len(t, lines, 2, "应有 event 行和 data 行")
			dataLine := lines[1]
			require.True(t, strings.HasPrefix(dataLine, "data: "), "第二行应以 'data: ' 开头")
			jsonStr := strings.TrimPrefix(dataLine, "data: ")

			// 验证 JSON 合法性
			var parsed map[string]any
			err := json.Unmarshal([]byte(jsonStr), &parsed)
			require.NoError(t, err, "JSON 应能被成功解析，原始 JSON: %s", jsonStr)

			// 验证结构
			errorObj, ok := parsed["error"].(map[string]any)
			require.True(t, ok, "应包含 error 对象")
			assert.Equal(t, tt.errType, errorObj["type"])
			assert.Equal(t, tt.message, errorObj["message"])
	placeholder)
placeholder
placeholder

func TestOpenAIHandleStreamingAwareErrorWithCode_EmitsStableClassification(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)

	h := &OpenAIGatewayHandler{placeholder
	h.handleStreamingAwareErrorWithCode(
		c,
		http.StatusBadGateway,
		"upstream_error",
		service.OpenAIUpstreamHTTP2StreamErrorCode,
		"Upstream HTTP/2 stream failed",
		true,
		true,
	)

	body := w.Body.String()
	require.Contains(t, body, "event: error\n")
	require.Equal(t, "upstream_error", gjson.Get(body[strings.Index(body, "{"):], "error.type").String())
	require.Equal(t, service.OpenAIUpstreamHTTP2StreamErrorCode, gjson.Get(body[strings.Index(body, "{"):], "error.code").String())
	require.NotContains(t, body, "stream ID")

	streamErr, ok := service.GetOpsStreamError(c)
	require.True(t, ok)
	require.True(t, streamErr.CountTowardsSLA)
	require.Equal(t, http.StatusBadGateway, streamErr.IntendedStatus)
placeholder

func TestOpenAIForwardSucceededForScheduling(t *testing.T) {
	require.True(t, openAIForwardSucceededForScheduling(nil))
	require.True(t, openAIForwardSucceededForScheduling(&service.OpenAIForwardResult{placeholder))
	require.True(t, openAIForwardSucceededForScheduling(&service.OpenAIForwardResult{
		OpenAIWSMode:          true,
		UpstreamTerminalEvent: "response.completed",
placeholder))
	require.False(t, openAIForwardSucceededForScheduling(&service.OpenAIForwardResult{
		OpenAIWSMode:          true,
		UpstreamTerminalEvent: "response.failed",
placeholder))
placeholder

func TestOpenAIResponsesRequiredCapability(t *testing.T) {
	tests := []struct {
		name        string
		imageIntent bool
		platform    string
		want        service.OpenAIEndpointCapability
placeholder{
		{
			name:        "OpenAI explicit image intent requires Responses",
			imageIntent: true,
			platform:    service.PlatformOpenAI,
			want:        service.OpenAIEndpointCapabilityResponses,
	placeholder,
		{
			name:        "Grok explicit image intent keeps chat capability",
			imageIntent: true,
			platform:    service.PlatformGrok,
			want:        service.OpenAIEndpointCapabilityChatCompletions,
	placeholder,
		{
			name:     "non-image intent keeps chat capability",
			platform: service.PlatformOpenAI,
			want:     service.OpenAIEndpointCapabilityChatCompletions,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, openAIResponsesRequiredCapability(tt.imageIntent, tt.platform))
	placeholder)
placeholder
placeholder

func TestResolveOpenAIMessagesMetadataSession_DoesNotDerivePromptCacheKey(t *testing.T) {
	body := []byte(`{"model":"claude-sonnet-4-5","metadata":{"user_id":"claude-code-session"placeholder,"messages":[{"role":"user","content":"hello"placeholder]placeholder`)

	sessionHash, promptCacheKey := resolveOpenAIMessagesMetadataSession("", "", "claude-sonnet-4-5", body)

	require.NotEmpty(t, sessionHash)
	require.Empty(t, promptCacheKey)
placeholder

func TestResolveOpenAIMessagesMetadataSession_PreservesExplicitPromptCacheKey(t *testing.T) {
	body := []byte(`{"metadata":{"user_id":"claude-code-session"placeholderplaceholder`)

	sessionHash, promptCacheKey := resolveOpenAIMessagesMetadataSession("", "explicit-cache", "claude-sonnet-4-5", body)

	require.NotEmpty(t, sessionHash)
	require.Equal(t, "explicit-cache", promptCacheKey)
placeholder

func TestOpenAIHandleStreamingAwareError_NonStreaming(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	h := &OpenAIGatewayHandler{placeholder
	h.handleStreamingAwareError(c, http.StatusBadGateway, "upstream_error", "test error", false)

	// 非流式应返回 JSON 响应
	assert.Equal(t, http.StatusBadGateway, w.Code)

	var parsed map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &parsed)
placeholder
	errorObj, ok := parsed["error"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "upstream_error", errorObj["type"])
	assert.Equal(t, "test error", errorObj["message"])
placeholder

func TestReadRequestBodyWithPrealloc(t *testing.T) {
	payload := `{"model":"gpt-5","input":"hello"placeholder`
	req := httptest.NewRequest(http.MethodPost, "/v1/responses", strings.NewReader(payload))
	req.ContentLength = int64(len(payload))

	body, err := pkghttputil.ReadRequestBodyWithPrealloc(req)
placeholder
	require.Equal(t, payload, string(body))
placeholder

func TestReadRequestBodyWithPrealloc_MaxBytesError(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/responses", strings.NewReader(strings.Repeat("x", 8)))
	req.Body = http.MaxBytesReader(rec, req.Body, 4)

	_, err := pkghttputil.ReadRequestBodyWithPrealloc(req)
placeholder
	var maxErr *http.MaxBytesError
	require.ErrorAs(t, err, &maxErr)
placeholder

func TestOpenAIEnsureForwardErrorResponse_WritesFallbackWhenNotWritten(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	h := &OpenAIGatewayHandler{placeholder
	wrote := h.ensureForwardErrorResponse(c, false)

	require.True(t, wrote)
	require.Equal(t, http.StatusBadGateway, w.Code)

	var parsed map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &parsed)
placeholder
	errorObj, ok := parsed["error"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "upstream_error", errorObj["type"])
	assert.Equal(t, "Upstream request failed", errorObj["message"])
placeholder

// Writer 已写后 ensureForwardErrorResponse 必须仍然把错误信息以 SSE
// 形式追加给客户端（streamStarted 强制 true）。
// 这是 case B 修复：旧实现遇到 Writer.Written 直接 return false，
// 客户端只能拿到 silent EOF；Codex CLI 报 "stream closed before response.completed"。
func TestOpenAIEnsureForwardErrorResponse_AppendsSSEAfterWritten(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.String(http.StatusTeapot, "already written")

	h := &OpenAIGatewayHandler{placeholder
	wrote := h.ensureForwardErrorResponse(c, false)

	require.True(t, wrote, "must attempt to communicate the failure to the client via SSE")
	// 状态码改不了（headers 已 flush），但 body 应该追加 SSE 错误事件。
	require.Equal(t, http.StatusTeapot, w.Code)
	assert.Contains(t, w.Body.String(), "already written")
	// 非 /responses 路径走 legacy event: error 分支。
	assert.Contains(t, w.Body.String(), "event: error\n")
placeholder

// case B 回归测试：/responses 路径，Writer 已被写过（模拟 ping flushed），
// ensureForwardErrorResponse 必须发 response.failed，让 Codex 收到合规终止事件。
func TestOpenAIEnsureForwardErrorResponse_ResponsesRouteAfterWrittenEmitsResponseFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, EndpointResponses, nil)
	// 模拟 ping 已 flush 的状态：Writer 已写过 1 个字节
	_, _ = c.Writer.WriteString(":\n\n")

	h := &OpenAIGatewayHandler{placeholder
	wrote := h.ensureForwardErrorResponse(c, false)

	require.True(t, wrote)
	body := w.Body.String()
	assert.Contains(t, body, ":\n\n", "earlier ping bytes preserved")
	assert.Contains(t, body, "event: response.failed\n", "appended a Responses terminal event")
	assert.Contains(t, body, `"type":"response.failed"`)
	assert.Contains(t, body, `"code":"upstream_error"`)
	assert.Contains(t, body, "Upstream request failed")
placeholder

func TestOpenAIEnsureForwardErrorResponse_AfterDeltaAppendsSingleValidResponseFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, EndpointResponses, nil)

	delta := `{"type":"response.output_text.delta","delta":"ok","sequence_number":1placeholder`
	_, err := c.Writer.WriteString("event: response.output_text.delta\ndata: " + delta + "\n\n")
placeholder

	h := &OpenAIGatewayHandler{placeholder
	require.True(t, h.ensureForwardErrorResponse(c, true))

	frames := strings.Split(strings.TrimSuffix(w.Body.String(), "\n\n"), "\n\n")
	require.Len(t, frames, 2)
	errorEvents := 0
	for _, frame := range frames {
		lines := strings.Split(frame, "\n")
		require.Len(t, lines, 2)
		require.True(t, strings.HasPrefix(lines[0], "event: "))
		require.True(t, strings.HasPrefix(lines[1], "data: "))

		eventType := strings.TrimPrefix(lines[0], "event: ")
		data := strings.TrimPrefix(lines[1], "data: ")
		require.True(t, json.Valid([]byte(data)), "each downstream SSE frame must contain valid JSON")
		var event struct {
			Type string `json:"type"`
	placeholder
		require.NoError(t, json.Unmarshal([]byte(data), &event))
		require.Equal(t, eventType, event.Type)
		if eventType == "response.failed" {
			errorEvents++
	placeholder
placeholder
	require.Equal(t, 1, errorEvents)
placeholder

func TestOpenAIEnsureForwardErrorResponse_CompactKeepaliveOnlyWritesResponseFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, EndpointResponses, nil)
	service.MarkOpenAICompactClientStream(c)

	stop := service.StartOpenAICompactSSEKeepalive(c, 5*time.Millisecond)
	defer stop()
	before := service.OpenAICompactKeepaliveAdjustedWrittenSize(c)
	require.Eventually(t, c.Writer.Written, time.Second, time.Millisecond)
	require.Equal(t, before, service.OpenAICompactKeepaliveAdjustedWrittenSize(c))

	h := &OpenAIGatewayHandler{placeholder
	require.True(t, h.ensureForwardErrorResponse(c, false))
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "event: response.failed\n")
	require.NotContains(t, w.Body.String(), "event: error\n")
placeholder

func TestOpenAIEnsureForwardErrorResponse_ImageJSONKeepaliveWritesSingleJSONFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", nil)

	stop := service.StartOpenAIImagesJSONKeepalive(c, 5*time.Millisecond)
	defer stop()
	before := service.OpenAIImagesJSONKeepaliveAdjustedWrittenSize(c)
	require.Eventually(t, c.Writer.Written, time.Second, time.Millisecond)
	require.Equal(t, before, service.OpenAIImagesJSONKeepaliveAdjustedWrittenSize(c))
	require.False(t, openAIForwardErrorAlreadyCommunicated(c, before, errors.New("read upstream response: unexpected EOF")))

	h := &OpenAIGatewayHandler{placeholder
	wrote := h.ensureForwardErrorResponse(c, false)

	require.True(t, wrote)
	require.Equal(t, http.StatusOK, w.Code, "heartbeat already committed the status")
	require.True(t, json.Valid(w.Body.Bytes()), w.Body.String())
	require.NotContains(t, w.Body.String(), "event:")
	require.NotContains(t, w.Body.String(), "data:")

	decoder := json.NewDecoder(strings.NewReader(w.Body.String()))
	var payload map[string]any
	require.NoError(t, decoder.Decode(&payload))
	require.ErrorIs(t, decoder.Decode(&payload), io.EOF)
	require.Equal(t, "upstream_error", gjson.Get(w.Body.String(), "error.type").String())
	require.Equal(t, "Upstream request failed", gjson.Get(w.Body.String(), "error.message").String())
placeholder

func TestOpenAIEnsureForwardErrorResponse_ImageJSONKeepalivePreservesCompletedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", nil)

	stop := service.StartOpenAIImagesJSONKeepalive(c, 5*time.Millisecond)
	defer stop()
	before := service.OpenAIImagesJSONKeepaliveAdjustedWrittenSize(c)
	require.Eventually(t, c.Writer.Written, time.Second, time.Millisecond)
	c.JSON(http.StatusOK, gin.H{"data": []gin.H{{"b64_json": "aW1hZ2U="placeholderplaceholderplaceholder)
	completedBody := w.Body.String()
	require.True(t, json.Valid([]byte(completedBody)), completedBody)
	require.Greater(t, service.OpenAIImagesJSONKeepaliveAdjustedWrittenSize(c), before)
	require.False(t, openAIForwardErrorAlreadyCommunicated(c, before, errors.New("read upstream trailer: unexpected EOF")))

	h := &OpenAIGatewayHandler{placeholder
	wrote := h.ensureForwardErrorResponse(c, false)

	require.False(t, wrote, "the completed Images JSON already communicated the response")
	require.Equal(t, completedBody, w.Body.String())
	require.NotContains(t, w.Body.String(), "event:")
	require.NotContains(t, w.Body.String(), "data:")

	decoder := json.NewDecoder(strings.NewReader(w.Body.String()))
	var payload map[string]any
	require.NoError(t, decoder.Decode(&payload))
	require.ErrorIs(t, decoder.Decode(&payload), io.EOF)
	require.Equal(t, "aW1hZ2U=", gjson.Get(w.Body.String(), "data.0.b64_json").String())
placeholder

func TestOpenAIEnsureForwardErrorResponse_FastImageJSONKeepalivePreservesCompletedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", nil)

	stop := service.StartOpenAIImagesJSONKeepalive(c, time.Hour)
	defer stop()
	before := service.OpenAIImagesJSONKeepaliveAdjustedWrittenSize(c)
	c.JSON(http.StatusOK, gin.H{"data": []gin.H{{"b64_json": "ZmFzdC1pbWFnZQ=="placeholderplaceholderplaceholder)
	completedBody := w.Body.String()
	require.True(t, json.Valid([]byte(completedBody)), completedBody)
	require.Greater(t, service.OpenAIImagesJSONKeepaliveAdjustedWrittenSize(c), before)
	require.False(t, openAIForwardErrorAlreadyCommunicated(c, before, errors.New("read upstream trailer: unexpected EOF")))

	h := &OpenAIGatewayHandler{placeholder
	wrote := h.ensureForwardErrorResponse(c, false)

	require.False(t, wrote, "fast completed Images JSON already communicated the response")
	require.Equal(t, completedBody, w.Body.String())
	require.NotContains(t, w.Body.String(), "event:")
	require.NotContains(t, w.Body.String(), "data:")

	decoder := json.NewDecoder(strings.NewReader(w.Body.String()))
	var payload map[string]any
	require.NoError(t, decoder.Decode(&payload))
	require.ErrorIs(t, decoder.Decode(&payload), io.EOF)
	require.Equal(t, "ZmFzdC1pbWFnZQ==", gjson.Get(w.Body.String(), "data.0.b64_json").String())
placeholder

func TestShouldLogOpenAIForwardFailureAsWarn(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("fallback_written_should_not_downgrade", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		require.False(t, shouldLogOpenAIForwardFailureAsWarn(c, true))
placeholder)

	t.Run("context_nil_should_not_downgrade", func(t *testing.T) {
		require.False(t, shouldLogOpenAIForwardFailureAsWarn(nil, false))
placeholder)

	t.Run("response_not_written_should_not_downgrade", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		require.False(t, shouldLogOpenAIForwardFailureAsWarn(c, false))
placeholder)

	t.Run("response_already_written_should_downgrade", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		c.String(http.StatusForbidden, "already written")
		require.True(t, shouldLogOpenAIForwardFailureAsWarn(c, false))
placeholder)
placeholder

func TestOpenAIRecoverResponsesPanic_WritesFallbackResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)

	h := &OpenAIGatewayHandler{placeholder
	streamStarted := false
	require.NotPanics(t, func() {
		func() {
			defer h.recoverResponsesPanic(c, &streamStarted)
			panic("test panic")
	placeholder()
placeholder)

	require.Equal(t, http.StatusBadGateway, w.Code)

	var parsed map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &parsed)
placeholder

	errorObj, ok := parsed["error"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "upstream_error", errorObj["type"])
	assert.Equal(t, "Upstream request failed", errorObj["message"])
placeholder

func TestOpenAIRecoverResponsesPanic_NoPanicNoWrite(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)

	h := &OpenAIGatewayHandler{placeholder
	streamStarted := false
	require.NotPanics(t, func() {
		func() {
			defer h.recoverResponsesPanic(c, &streamStarted)
	placeholder()
placeholder)

	require.False(t, c.Writer.Written())
	assert.Equal(t, "", w.Body.String())
placeholder

// Panic 在已 flush 的 /v1/responses 流中：状态码无法改（已 written），
// 但 body 应追加 response.failed 让客户端识别为合规截断而不是 silent EOF。
func TestOpenAIRecoverResponsesPanic_AppendsResponseFailedAfterWritten(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	c.String(http.StatusTeapot, "already written")

	h := &OpenAIGatewayHandler{placeholder
	streamStarted := false
	require.NotPanics(t, func() {
		func() {
			defer h.recoverResponsesPanic(c, &streamStarted)
			panic("test panic")
	placeholder()
placeholder)

	require.Equal(t, http.StatusTeapot, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "already written")
	assert.Contains(t, body, "event: response.failed\n")
placeholder

func TestOpenAIMissingResponsesDependencies(t *testing.T) {
	t.Run("nil_handler", func(t *testing.T) {
		var h *OpenAIGatewayHandler
		require.Equal(t, []string{"handler"placeholder, h.missingResponsesDependencies())
placeholder)

	t.Run("all_dependencies_missing", func(t *testing.T) {
		h := &OpenAIGatewayHandler{placeholder
		require.Equal(t,
			[]string{"gatewayService", "billingCacheService", "apiKeyService", "concurrencyHelper"placeholder,
			h.missingResponsesDependencies(),
		)
placeholder)

	t.Run("all_dependencies_present", func(t *testing.T) {
		h := &OpenAIGatewayHandler{
			gatewayService:      &service.OpenAIGatewayService{placeholder,
			billingCacheService: &service.BillingCacheService{placeholder,
			apiKeyService:       &service.APIKeyService{placeholder,
			concurrencyHelper: &ConcurrencyHelper{
				concurrencyService: &service.ConcurrencyService{placeholder,
		placeholder,
	placeholder
		require.Empty(t, h.missingResponsesDependencies())
placeholder)
placeholder

func TestOpenAIEnsureResponsesDependencies(t *testing.T) {
	t.Run("missing_dependencies_returns_503", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)

		h := &OpenAIGatewayHandler{placeholder
		ok := h.ensureResponsesDependencies(c, nil)

		require.False(t, ok)
		require.Equal(t, http.StatusServiceUnavailable, w.Code)
		var parsed map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &parsed)
	placeholder
		errorObj, exists := parsed["error"].(map[string]any)
		require.True(t, exists)
		assert.Equal(t, "api_error", errorObj["type"])
		assert.Equal(t, "Service temporarily unavailable", errorObj["message"])
placeholder)

	t.Run("already_written_response_not_overridden", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
		c.String(http.StatusTeapot, "already written")

		h := &OpenAIGatewayHandler{placeholder
		ok := h.ensureResponsesDependencies(c, nil)

		require.False(t, ok)
		require.Equal(t, http.StatusTeapot, w.Code)
		assert.Equal(t, "already written", w.Body.String())
placeholder)

	t.Run("dependencies_ready_returns_true_and_no_write", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)

		h := &OpenAIGatewayHandler{
			gatewayService:      &service.OpenAIGatewayService{placeholder,
			billingCacheService: &service.BillingCacheService{placeholder,
			apiKeyService:       &service.APIKeyService{placeholder,
			concurrencyHelper: &ConcurrencyHelper{
				concurrencyService: &service.ConcurrencyService{placeholder,
		placeholder,
	placeholder
		ok := h.ensureResponsesDependencies(c, nil)

		require.True(t, ok)
		require.False(t, c.Writer.Written())
		assert.Equal(t, "", w.Body.String())
placeholder)
placeholder

func TestResolveOpenAIMessagesDispatchMappedModel(t *testing.T) {
	t.Run("exact_claude_model_override_wins", func(t *testing.T) {
		apiKey := &service.APIKey{
			Group: &service.Group{
				MessagesDispatchModelConfig: service.OpenAIMessagesDispatchModelConfig{
					SonnetMappedModel: "gpt-5.2",
					ExactModelMappings: map[string]string{
						"claude-sonnet-4-5-20250929": "gpt-5.4-mini-high",
						"claude-fable-5":             "gpt-5.6-sol",
				placeholder,
			placeholder,
		placeholder,
	placeholder
		require.Equal(t, "gpt-5.4-mini", resolveOpenAIMessagesDispatchMappedModel(nil, apiKey, "claude-sonnet-4-5-20250929"))
		require.Equal(t, "gpt-5.6-sol", resolveOpenAIMessagesDispatchMappedModel(nil, apiKey, "claude-fable-5"))
placeholder)

	t.Run("uses_family_default_when_no_override", func(t *testing.T) {
		apiKey := &service.APIKey{Group: &service.Group{placeholderplaceholder
		require.Equal(t, "gpt-5.4", resolveOpenAIMessagesDispatchMappedModel(nil, apiKey, "claude-opus-4-6"))
		require.Equal(t, "gpt-5.3-codex", resolveOpenAIMessagesDispatchMappedModel(nil, apiKey, "claude-sonnet-4-5-20250929"))
		require.Equal(t, "gpt-5.4-mini", resolveOpenAIMessagesDispatchMappedModel(nil, apiKey, "placeholder"))
placeholder)

	t.Run("returns_empty_for_non_claude_or_missing_group", func(t *testing.T) {
		require.Empty(t, resolveOpenAIMessagesDispatchMappedModel(nil, nil, "claude-sonnet-4-5-20250929"))
		require.Empty(t, resolveOpenAIMessagesDispatchMappedModel(nil, &service.APIKey{placeholder, "claude-sonnet-4-5-20250929"))
		require.Empty(t, resolveOpenAIMessagesDispatchMappedModel(nil, &service.APIKey{Group: &service.Group{placeholderplaceholder, "gpt-5.4"))
placeholder)

	t.Run("grok_group_maps_claude_cli_model_to_grok_default", func(t *testing.T) {
		original := xai.RuntimeModelMappingOptions()
		t.Cleanup(func() { xai.SetRuntimeModelMappingOptions(original) placeholder)
		xai.SetRuntimeModelMappingOptions(xai.ModelMappingOptions{EnableCrossClientMap: trueplaceholder)
		apiKey := &service.APIKey{
			Group: &service.Group{
				Platform: service.PlatformGrok,
		placeholder,
	placeholder
		require.Equal(t, "grok-4.6", resolveOpenAIMessagesDispatchMappedModel(nil, apiKey, "claude-sonnet-4-5"))
		require.Empty(t, resolveOpenAIMessagesDispatchMappedModel(nil, apiKey, "grok"))
placeholder)

	t.Run("does_not_fall_back_to_group_default_mapped_model", func(t *testing.T) {
		apiKey := &service.APIKey{
			Group: &service.Group{
				DefaultMappedModel: "gpt-5.4",
		placeholder,
	placeholder
		require.Empty(t, resolveOpenAIMessagesDispatchMappedModel(nil, apiKey, "gpt-5.4"))
		require.Equal(t, "gpt-5.3-codex", resolveOpenAIMessagesDispatchMappedModel(nil, apiKey, "claude-sonnet-4-5-20250929"))
placeholder)
placeholder

func TestOpenAIGatewayMessagesDispatchGateAllowsGrokGroups(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("openai_group_without_dispatch_flag_is_rejected", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(`{"model":"claude-sonnet-4-5","messages":[{"role":"user","content":"hi"placeholder]placeholder`))
		groupID := int64(4101)
		c.Set(string(middleware.ContextKeyAPIKey), &service.APIKey{
			ID:      5101,
			GroupID: &groupID,
			User:    &service.User{ID: 6101placeholder,
			Group: &service.Group{
				ID:                    groupID,
				Platform:              service.PlatformOpenAI,
				AllowMessagesDispatch: false,
		placeholder,
	placeholder)
		c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: 6101, Concurrency: 1placeholder)

		h := &OpenAIGatewayHandler{placeholder
		h.Messages(c)

		require.Equal(t, http.StatusForbidden, rec.Code)
		require.Equal(t, "permission_error", gjson.GetBytes(rec.Body.Bytes(), "error.type").String())
		require.Contains(t, rec.Body.String(), "This group does not allow /v1/messages dispatch")
placeholder)

	t.Run("grok_group_without_dispatch_flag_reaches_gateway_dependencies", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(`{"model":"grok-4.3","messages":[{"role":"user","content":"hi"placeholder]placeholder`))
		groupID := int64(4102)
		c.Set(string(middleware.ContextKeyAPIKey), &service.APIKey{
			ID:      5102,
			GroupID: &groupID,
			User:    &service.User{ID: 6102placeholder,
			Group: &service.Group{
				ID:                    groupID,
				Platform:              service.PlatformGrok,
				AllowMessagesDispatch: false,
		placeholder,
	placeholder)
		c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: 6102, Concurrency: 1placeholder)

		h := &OpenAIGatewayHandler{placeholder
		h.Messages(c)

		require.Equal(t, http.StatusServiceUnavailable, rec.Code)
		require.Equal(t, "api_error", gjson.GetBytes(rec.Body.Bytes(), "error.type").String())
		require.NotContains(t, rec.Body.String(), "This group does not allow /v1/messages dispatch")
placeholder)
placeholder

func TestOpenAIModelMappedBody(t *testing.T) {
	body := []byte(`{"model":"alias","input":"hello"placeholder`)
	calls := 0

	forwardBody := openAIModelMappedBody(body, true, "gpt-5.4", func(body []byte, newModel string) []byte {
		calls++
		return service.ReplaceModelInBody(body, newModel)
placeholder)

	require.Equal(t, 1, calls)
	require.Equal(t, "gpt-5.4", gjson.GetBytes(forwardBody, "model").String())
	require.Equal(t, "alias", gjson.GetBytes(body, "model").String())
placeholder

func TestOpenAIModelMappedBodyCache(t *testing.T) {
	body := []byte(`{"model":"alias","input":"hello"placeholder`)
	calls := 0
	mappedBody := newOpenAIModelMappedBodyCache(body, func(body []byte, newModel string) []byte {
		calls++
		return service.ReplaceModelInBody(body, newModel)
placeholder)

	first := mappedBody(true, "gpt-5.4")
	second := mappedBody(true, "gpt-5.4")
	third := mappedBody(true, "gpt-5.3-codex")
	unmapped := mappedBody(false, "ignored")

	require.Equal(t, 2, calls)
	require.Equal(t, "gpt-5.4", gjson.GetBytes(first, "model").String())
	require.Equal(t, "gpt-5.4", gjson.GetBytes(second, "model").String())
	require.Equal(t, "gpt-5.3-codex", gjson.GetBytes(third, "model").String())
	require.Equal(t, body, unmapped)
	require.Same(t, &first[0], &second[0])
placeholder

func TestOpenAIResponses_MissingDependencies_ReturnsServiceUnavailable(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", strings.NewReader(`{"model":"gpt-5","stream":falseplaceholder`))
	c.Request.Header.Set("Content-Type", "application/json")

	groupID := int64(2)
	c.Set(string(middleware.ContextKeyAPIKey), &service.APIKey{
		ID:      10,
		GroupID: &groupID,
placeholder)
	c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{
		UserID:      1,
		Concurrency: 1,
placeholder)

	// 故意使用未初始化依赖，验证快速失败而不是崩溃。
	h := &OpenAIGatewayHandler{placeholder
	require.NotPanics(t, func() {
		h.Responses(c)
placeholder)

	require.Equal(t, http.StatusServiceUnavailable, w.Code)

	var parsed map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &parsed)
placeholder

	errorObj, ok := parsed["error"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "api_error", errorObj["type"])
	assert.Equal(t, "Service temporarily unavailable", errorObj["message"])
placeholder

func TestOpenAIResponses_SetsClientTransportHTTP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", strings.NewReader(`{"model":"gpt-5"placeholder`))
	c.Request.Header.Set("Content-Type", "application/json")

	h := &OpenAIGatewayHandler{placeholder
	h.Responses(c)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.Equal(t, service.OpenAIClientTransportHTTP, service.GetOpenAIClientTransport(c))
placeholder

func TestOpenAIResponses_RejectsMessageIDAsPreviousResponseID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", strings.NewReader(
		`{"model":"gpt-5.1","stream":false,"previous_response_id":"msg_123456","input":[{"type":"input_text","text":"hello"placeholder]placeholder`,
	))
	c.Request.Header.Set("Content-Type", "application/json")

	groupID := int64(2)
	c.Set(string(middleware.ContextKeyAPIKey), &service.APIKey{
		ID:      101,
		GroupID: &groupID,
		User:    &service.User{ID: 1placeholder,
placeholder)
	c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{
		UserID:      1,
		Concurrency: 1,
placeholder)

	h := newOpenAIHandlerForPreviousResponseIDValidation(t, nil)
	h.Responses(c)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "previous_response_id must be a response.id")
placeholder

func TestOpenAIResponses_AcceptsHTTPContinuationPreviousResponseIDBeforeRouting(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", strings.NewReader(
		`{"model":"gpt-5.1","stream":false,"previous_response_id":"resp_123456","input":[{"type":"input_text","text":"hello"placeholder]placeholder`,
	))
	c.Request.Header.Set("Content-Type", "application/json")

	groupID := int64(2)
	c.Set(string(middleware.ContextKeyAPIKey), &service.APIKey{
		ID:      101,
		GroupID: &groupID,
		User:    &service.User{ID: 1placeholder,
placeholder)
	c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{
		UserID:      1,
		Concurrency: 1,
placeholder)

	h := newOpenAIHandlerForPreviousResponseIDValidation(t, nil)
	h.Responses(c)

	require.NotEqual(t, http.StatusBadRequest, w.Code)
	require.NotContains(t, w.Body.String(), "Responses WebSocket v2")
placeholder

func TestOpenAIResponses_FunctionCallOutputHTTPGuidanceDoesNotSuggestPreviousResponseReuse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", strings.NewReader(
		`{"model":"gpt-5.1","stream":false,"input":[{"type":"function_call_output","output":"{placeholder"placeholder]placeholder`,
	))
	c.Request.Header.Set("Content-Type", "application/json")

	groupID := int64(2)
	c.Set(string(middleware.ContextKeyAPIKey), &service.APIKey{
		ID:      101,
		GroupID: &groupID,
		User:    &service.User{ID: 1placeholder,
placeholder)
	c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{
		UserID:      1,
		Concurrency: 1,
placeholder)

	h := newOpenAIHandlerForPreviousResponseIDValidation(t, nil)
	h.Responses(c)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "Responses WebSocket v2")
	require.NotContains(t, w.Body.String(), "reuse previous_response_id")
placeholder

func TestOpenAIResponsesWebSocket_SetsClientTransportWSWhenUpgradeValid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/openai/v1/responses", nil)
	c.Request.Header.Set("Upgrade", "websocket")
	c.Request.Header.Set("Connection", "Upgrade")

	h := &OpenAIGatewayHandler{placeholder
	h.ResponsesWebSocket(c)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.Equal(t, service.OpenAIClientTransportWS, service.GetOpenAIClientTransport(c))
placeholder

func TestOpenAIResponsesWebSocket_InvalidUpgradeDoesNotSetTransport(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/openai/v1/responses", nil)

	h := &OpenAIGatewayHandler{placeholder
	h.ResponsesWebSocket(c)

	require.Equal(t, http.StatusUpgradeRequired, w.Code)
	require.Equal(t, service.OpenAIClientTransportUnknown, service.GetOpenAIClientTransport(c))
placeholder

func TestOpenAIResponsesWebSocket_IngressCapacityRejected(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := &concurrencyCacheMock{
		acquireIngressLeaseFn: func(context.Context, int64, int, string) (bool, error) {
			return false, nil
	placeholder,
placeholder
	h := newOpenAIHandlerForPreviousResponseIDValidation(t, cache)
	h.cfg = &config.Config{placeholder
	h.cfg.Gateway.OpenAIWS.MaxIngressConnectionsPerAPIKey = 1
	wsServer := newOpenAIWSHandlerTestServer(t, h, middleware.AuthSubject{UserID: 1, Concurrency: 1placeholder)
	defer wsServer.Close()

	dialCtx, cancelDial := context.WithTimeout(context.Background(), 3*time.Second)
	clientConn, response, err := coderws.Dial(dialCtx, "ws"+strings.TrimPrefix(wsServer.URL, "http")+"/openai/v1/responses", nil)
	cancelDial()
placeholder
	require.Nil(t, clientConn)
	require.NotNil(t, response)
	require.Equal(t, http.StatusTooManyRequests, response.StatusCode)
	_ = response.Body.Close()
placeholder

func TestOpenAIResponsesWebSocket_IngressLeaseBackendUnavailableBeforeUpgrade(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := &concurrencyCacheMock{
		acquireIngressLeaseFn: func(context.Context, int64, int, string) (bool, error) {
			return false, errors.New("redis unavailable")
	placeholder,
placeholder
	h := newOpenAIHandlerForPreviousResponseIDValidation(t, cache)
	h.cfg = &config.Config{placeholder
	h.cfg.Gateway.OpenAIWS.MaxIngressConnectionsPerAPIKey = 1
	wsServer := newOpenAIWSHandlerTestServer(t, h, middleware.AuthSubject{UserID: 1, Concurrency: 1placeholder)
	defer wsServer.Close()

	dialCtx, cancelDial := context.WithTimeout(context.Background(), 3*time.Second)
	clientConn, response, err := coderws.Dial(dialCtx, "ws"+strings.TrimPrefix(wsServer.URL, "http")+"/openai/v1/responses", nil)
	cancelDial()
placeholder
	require.Nil(t, clientConn)
	require.NotNil(t, response)
	require.Equal(t, http.StatusServiceUnavailable, response.StatusCode)
	_ = response.Body.Close()
placeholder

func TestOpenAIResponsesWebSocket_FirstMessageTimeoutUsesConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logSink, restore := captureHandlerStructuredLog(t)
	defer restore()

	h := newOpenAIHandlerForPreviousResponseIDValidation(t, nil)
	h.cfg = &config.Config{placeholder
	h.cfg.Gateway.OpenAIWS.ClientFirstMessageTimeoutSeconds = 1
	wsServer := newOpenAIWSHandlerTestServer(t, h, middleware.AuthSubject{UserID: 1, Concurrency: 1placeholder)
	defer wsServer.Close()

	dialCtx, cancelDial := context.WithTimeout(context.Background(), 3*time.Second)
	clientConn, _, err := coderws.Dial(dialCtx, "ws"+strings.TrimPrefix(wsServer.URL, "http")+"/openai/v1/responses", nil)
	cancelDial()
placeholder
	defer func() { _ = clientConn.CloseNow() placeholder()

	started := time.Now()
	readCtx, cancelRead := context.WithTimeout(context.Background(), 5*time.Second)
	_, _, err = clientConn.Read(readCtx)
	cancelRead()
	elapsed := time.Since(started)

placeholder
	require.NotErrorIs(t, err, context.DeadlineExceeded)
	require.GreaterOrEqual(t, elapsed, 500*time.Millisecond)
	require.Less(t, elapsed, 4*time.Second)
	require.Eventually(t, func() bool {
		readTimeout, ok := logSink.FieldValueForMessage("openai.websocket_read_first_message_failed", "read_timeout")
		return ok && readTimeout == time.Second &&
			logSink.ContainsMessageAtLevel("openai.websocket_read_first_message_failed", "warn")
placeholder, time.Second, 10*time.Millisecond)
placeholder

func TestOpenAIResponsesWebSocket_IngressLeaseReleasedOnEarlyReturn(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := &concurrencyCacheMock{
		acquireIngressLeaseFn: func(context.Context, int64, int, string) (bool, error) {
			return true, nil
	placeholder,
placeholder
	h := newOpenAIHandlerForPreviousResponseIDValidation(t, cache)
	h.cfg = &config.Config{placeholder
	h.cfg.Gateway.OpenAIWS.MaxIngressConnectionsPerAPIKey = 1
	wsServer := newOpenAIWSHandlerTestServer(t, h, middleware.AuthSubject{UserID: 1, Concurrency: 1placeholder)
	defer wsServer.Close()

	dialCtx, cancelDial := context.WithTimeout(context.Background(), 3*time.Second)
	clientConn, _, err := coderws.Dial(dialCtx, "ws"+strings.TrimPrefix(wsServer.URL, "http")+"/openai/v1/responses", nil)
	cancelDial()
placeholder
	defer func() { _ = clientConn.CloseNow() placeholder()

	writeCtx, cancelWrite := context.WithTimeout(context.Background(), 3*time.Second)
	err = clientConn.Write(writeCtx, coderws.MessageBinary, []byte("not a response.create frame"))
	cancelWrite()
placeholder

	readCtx, cancelRead := context.WithTimeout(context.Background(), 3*time.Second)
	_, _, err = clientConn.Read(readCtx)
	cancelRead()
	var closeErr coderws.CloseError
	require.ErrorAs(t, err, &closeErr)
	require.Equal(t, coderws.StatusPolicyViolation, closeErr.Code)
	require.Eventually(t, func() bool {
		return atomic.LoadInt32(&cache.releaseIngressCalled) == 1
placeholder, time.Second, 10*time.Millisecond)
placeholder

func TestOpenAIResponsesWebSocket_IngressLeaseReleasedWhenUpgradeFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := &concurrencyCacheMock{
		acquireIngressLeaseFn: func(context.Context, int64, int, string) (bool, error) {
			return true, nil
	placeholder,
placeholder
	h := newOpenAIHandlerForPreviousResponseIDValidation(t, cache)
	h.cfg = &config.Config{placeholder
	h.cfg.Gateway.OpenAIWS.MaxIngressConnectionsPerAPIKey = 1
	wsServer := newOpenAIWSHandlerTestServer(t, h, middleware.AuthSubject{UserID: 1, Concurrency: 1placeholder)
	defer wsServer.Close()

	req, err := http.NewRequest(http.MethodGet, wsServer.URL+"/openai/v1/responses", nil)
placeholder
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Sec-WebSocket-Version", "13")
	resp, err := http.DefaultClient.Do(req)
placeholder
	_ = resp.Body.Close()
	require.NotEqual(t, http.StatusSwitchingProtocols, resp.StatusCode)
	require.Eventually(t, func() bool {
		return atomic.LoadInt32(&cache.releaseIngressCalled) == 1
placeholder, time.Second, 10*time.Millisecond)
placeholder

func TestOpenAIResponsesWebSocket_RejectsMessageIDAsPreviousResponseID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := newOpenAIHandlerForPreviousResponseIDValidation(t, nil)
	wsServer := newOpenAIWSHandlerTestServer(t, h, middleware.AuthSubject{UserID: 1, Concurrency: 1placeholder)
	defer wsServer.Close()

	dialCtx, cancelDial := context.WithTimeout(context.Background(), 3*time.Second)
	clientConn, _, err := coderws.Dial(dialCtx, "ws"+strings.TrimPrefix(wsServer.URL, "http")+"/openai/v1/responses", nil)
	cancelDial()
placeholder
	defer func() {
		_ = clientConn.CloseNow()
placeholder()

	writeCtx, cancelWrite := context.WithTimeout(context.Background(), 3*time.Second)
	err = clientConn.Write(writeCtx, coderws.MessageText, []byte(
		`{"type":"response.create","model":"gpt-5.1","stream":false,"previous_response_id":"msg_abc123"placeholder`,
	))
	cancelWrite()
placeholder

	readCtx, cancelRead := context.WithTimeout(context.Background(), 3*time.Second)
	_, _, err = clientConn.Read(readCtx)
	cancelRead()
placeholder
	var closeErr coderws.CloseError
	require.ErrorAs(t, err, &closeErr)
	require.Equal(t, coderws.StatusPolicyViolation, closeErr.Code)
	require.Contains(t, strings.ToLower(closeErr.Reason), "previous_response_id")
placeholder

func TestOpenAIResponsesWebSocket_PreviousResponseIDKindLoggedBeforeAcquireFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := &concurrencyCacheMock{
		acquireUserSlotFn: func(ctx context.Context, userID int64, maxConcurrency int, requestID string) (bool, error) {
			return false, errors.New("user slot unavailable")
	placeholder,
placeholder
	h := newOpenAIHandlerForPreviousResponseIDValidation(t, cache)
	wsServer := newOpenAIWSHandlerTestServer(t, h, middleware.AuthSubject{UserID: 1, Concurrency: 1placeholder)
	defer wsServer.Close()

	dialCtx, cancelDial := context.WithTimeout(context.Background(), 3*time.Second)
	clientConn, _, err := coderws.Dial(dialCtx, "ws"+strings.TrimPrefix(wsServer.URL, "http")+"/openai/v1/responses", nil)
	cancelDial()
placeholder
	defer func() {
		_ = clientConn.CloseNow()
placeholder()

	writeCtx, cancelWrite := context.WithTimeout(context.Background(), 3*time.Second)
	err = clientConn.Write(writeCtx, coderws.MessageText, []byte(
		`{"type":"response.create","model":"gpt-5.1","stream":false,"previous_response_id":"resp_prev_123"placeholder`,
	))
	cancelWrite()
placeholder

	readCtx, cancelRead := context.WithTimeout(context.Background(), 3*time.Second)
	_, _, err = clientConn.Read(readCtx)
	cancelRead()
placeholder
	var closeErr coderws.CloseError
	require.ErrorAs(t, err, &closeErr)
	require.Equal(t, coderws.StatusInternalError, closeErr.Code)
	require.Contains(t, strings.ToLower(closeErr.Reason), "failed to acquire user concurrency slot")
placeholder

type contentModerationHandlerSettingRepo struct {
	values map[string]string
placeholder

func (r *contentModerationHandlerSettingRepo) Get(ctx context.Context, key string) (*service.Setting, error) {
	if value, ok := r.values[key]; ok {
		return &service.Setting{Key: key, Value: valueplaceholder, nil
placeholder
	return nil, service.ErrSettingNotFound
placeholder

func (r *contentModerationHandlerSettingRepo) GetValue(ctx context.Context, key string) (string, error) {
	if value, ok := r.values[key]; ok {
		return value, nil
placeholder
	return "", service.ErrSettingNotFound
placeholder

func (r *contentModerationHandlerSettingRepo) Set(ctx context.Context, key, value string) error {
	if r.values == nil {
		r.values = map[string]string{placeholder
placeholder
	r.values[key] = value
	return nil
placeholder

func (r *contentModerationHandlerSettingRepo) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	out := map[string]string{placeholder
	for _, key := range keys {
		if value, ok := r.values[key]; ok {
			out[key] = value
	placeholder
placeholder
	return out, nil
placeholder

func (r *contentModerationHandlerSettingRepo) SetMultiple(ctx context.Context, settings map[string]string) error {
	if r.values == nil {
		r.values = map[string]string{placeholder
placeholder
	for key, value := range settings {
		r.values[key] = value
placeholder
	return nil
placeholder

func (r *contentModerationHandlerSettingRepo) GetAll(ctx context.Context) (map[string]string, error) {
	out := make(map[string]string, len(r.values))
	for key, value := range r.values {
		out[key] = value
placeholder
	return out, nil
placeholder

func (r *contentModerationHandlerSettingRepo) Delete(ctx context.Context, key string) error {
	delete(r.values, key)
	return nil
placeholder

type contentModerationHandlerTestRepo struct {
	mu   sync.Mutex
	logs []service.ContentModerationLog
placeholder

func (r *contentModerationHandlerTestRepo) CreateLog(ctx context.Context, log *service.ContentModerationLog) error {
	if log != nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.logs = append(r.logs, *log)
placeholder
	return nil
placeholder

func (r *contentModerationHandlerTestRepo) resetLogs() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logs = nil
placeholder

func (r *contentModerationHandlerTestRepo) logSnapshot() []service.ContentModerationLog {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]service.ContentModerationLog(nil), r.logs...)
placeholder

func (r *contentModerationHandlerTestRepo) ListLogs(ctx context.Context, filter service.ContentModerationLogFilter) ([]service.ContentModerationLog, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder

func (r *contentModerationHandlerTestRepo) CountFlaggedByUserSince(ctx context.Context, userID int64, since time.Time, excludeCyberPolicy bool) (int, error) {
	return 0, nil
placeholder

func (r *contentModerationHandlerTestRepo) CleanupExpiredLogs(ctx context.Context, hitBefore time.Time, nonHitBefore time.Time) (*service.ContentModerationCleanupResult, error) {
	return &service.ContentModerationCleanupResult{placeholder, nil
placeholder

func (r *contentModerationHandlerTestRepo) UpdateLogEmailSent(ctx context.Context, id int64, sent bool) error {
	return nil
placeholder

func TestOpenAIResponsesWebSocket_ContentModerationBlocksFirstFrame(t *testing.T) {
	gin.SetMode(gin.TestMode)

	moderationServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v1/moderations", r.URL.Path)
		_, _ = w.Write([]byte(`{"results":[{"category_scores":{"sexual":0.9placeholderplaceholder]placeholder`))
placeholder))
	defer moderationServer.Close()

	cfg := &service.ContentModerationConfig{
		Enabled:      true,
		Mode:         service.ContentModerationModePreBlock,
		BaseURL:      moderationServer.URL,
		Model:        "omni-moderation-latest",
		APIKeys:      []string{"sk-test"placeholder,
		SampleRate:   100,
		AllGroups:    true,
		BlockMessage: "内容审计测试阻断",
placeholder
	rawCfg, err := json.Marshal(cfg)
placeholder

	repo := &contentModerationHandlerTestRepo{placeholder
	settingRepo := &contentModerationHandlerSettingRepo{values: map[string]string{
		service.SettingKeyRiskControlEnabled:      "true",
		service.SettingKeyContentModerationConfig: string(rawCfg),
placeholderplaceholder
	moderationSvc := service.NewContentModerationService(
		settingRepo,
		repo,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	decision, err := moderationSvc.Check(context.Background(), service.ContentModerationCheckInput{
		UserID:   1,
		Endpoint: "/v1/responses",
		Provider: "openai",
		Model:    "gpt-5.5",
		Protocol: service.ContentModerationProtocolOpenAIResponses,
		Body:     []byte(`{"model":"gpt-5.5","input":[{"type":"message","role":"user","content":[{"type":"input_text","text":"bad prompt"placeholder]placeholder]placeholder`),
placeholder)
placeholder
	require.True(t, decision.Blocked)
	require.Eventually(t, func() bool {
		return len(repo.logSnapshot()) == 1
placeholder, time.Second, 10*time.Millisecond)
	repo.resetLogs()
	h := &OpenAIGatewayHandler{
		gatewayService:           &service.OpenAIGatewayService{placeholder,
		billingCacheService:      &service.BillingCacheService{placeholder,
		apiKeyService:            &service.APIKeyService{placeholder,
		contentModerationService: moderationSvc,
		concurrencyHelper:        NewConcurrencyHelper(service.NewConcurrencyService(&concurrencyCacheMock{placeholder), SSEPingFormatNone, time.Second),
placeholder
	wsServer := newOpenAIWSHandlerTestServer(t, h, middleware.AuthSubject{UserID: 1, Concurrency: 1placeholder)
	defer wsServer.Close()

	dialCtx, cancelDial := context.WithTimeout(context.Background(), 3*time.Second)
	clientConn, _, err := coderws.Dial(dialCtx, "ws"+strings.TrimPrefix(wsServer.URL, "http")+"/openai/v1/responses", nil)
	cancelDial()
placeholder
	defer func() {
		_ = clientConn.CloseNow()
placeholder()

	writeCtx, cancelWrite := context.WithTimeout(context.Background(), 3*time.Second)
	err = clientConn.Write(writeCtx, coderws.MessageText, []byte(`{
		"type":"response.create",
		"model":"gpt-5.5",
		"input":[{"type":"message","role":"user","content":[{"type":"input_text","text":"bad prompt"placeholder]placeholder]
placeholder`))
	cancelWrite()
placeholder

	readCtx, cancelRead := context.WithTimeout(context.Background(), 3*time.Second)
	_, payload, readErr := clientConn.Read(readCtx)
	cancelRead()
	if readErr == nil {
		require.Contains(t, string(payload), "content_policy_violation")
		require.Contains(t, string(payload), "内容审计测试阻断")
placeholder else {
		var closeErr coderws.CloseError
		require.ErrorAs(t, readErr, &closeErr)
		require.Equal(t, coderws.StatusPolicyViolation, closeErr.Code)
		require.Contains(t, closeErr.Reason, "内容审计测试阻断")
placeholder
	var logs []service.ContentModerationLog
	require.Eventually(t, func() bool {
		logs = repo.logSnapshot()
		return len(logs) == 1
placeholder, time.Second, 10*time.Millisecond)
	require.True(t, logs[0].Flagged)
	require.Equal(t, service.ContentModerationActionBlock, logs[0].Action)
	require.Equal(t, "bad prompt", logs[0].InputExcerpt)
placeholder

func TestOpenAIResponsesWebSocket_PassthroughUsageLogPersistsUserAgentAndReasoningEffort(t *testing.T) {
	got := runOpenAIResponsesWebSocketUsageLogCase(t, openAIResponsesWSUsageLogCase{
		firstPayload: `{"type":"response.create","model":"gpt-5.4","stream":false,"reasoning":{"effort":"HIGH"placeholderplaceholder`,
		userAgent:    testStringPtr("codex_cli_rs/0.125.0 test"),
placeholder)

	require.NotNil(t, got.log.UserAgent)
	require.Equal(t, "codex_cli_rs/0.125.0 test", *got.log.UserAgent)
	require.NotNil(t, got.log.ReasoningEffort)
	require.Equal(t, "high", *got.log.ReasoningEffort)
	require.True(t, got.log.OpenAIWSMode)
placeholder

func TestOpenAIResponsesWebSocket_PassthroughUsageLogInfersReasoningFromInitialRequestModel(t *testing.T) {
	got := runOpenAIResponsesWebSocketUsageLogCase(t, openAIResponsesWSUsageLogCase{
		firstPayload: `{"type":"response.create","model":"gpt-5.4-xhigh","stream":falseplaceholder`,
		userAgent:    testStringPtr("codex_cli_rs/0.125.0 mapped"),
		channelMapping: map[string]string{
			"gpt-5.4-xhigh": "gpt-5.4",
	placeholder,
placeholder)

	require.Equal(t, "gpt-5.4", gjson.GetBytes(got.upstreamFirstPayload, "model").String(),
		"上游首帧应使用渠道映射后的模型")
	require.NotNil(t, got.log.ReasoningEffort)
	require.Equal(t, "xhigh", *got.log.ReasoningEffort,
		"usage log reasoning effort 必须使用渠道映射前首帧模型后缀推导")
placeholder

func TestOpenAIResponsesWebSocket_PassthroughUsageLogLeavesUserAgentNilWhenMissing(t *testing.T) {
	got := runOpenAIResponsesWebSocketUsageLogCase(t, openAIResponsesWSUsageLogCase{
		firstPayload: `{"type":"response.create","model":"gpt-5.4","stream":false,"reasoning":{"effort":"medium"placeholderplaceholder`,
		userAgent:    testStringPtr(""),
placeholder)

	require.Nil(t, got.log.UserAgent, "空入站 User-Agent 不应由上游握手 UA 或默认 UA 兜底")
	require.NotNil(t, got.log.ReasoningEffort)
	require.Equal(t, "medium", *got.log.ReasoningEffort)
placeholder

func TestOpenAIResponsesWebSocket_PassthroughTracksModelPerTurn(t *testing.T) {
	got := runOpenAIResponsesWebSocketUsageLogCase(t, openAIResponsesWSUsageLogCase{
		firstPayload:  `{"type":"response.create","model":"sol","stream":falseplaceholder`,
		secondPayload: `{"type":"response.create","model":"terra","stream":falseplaceholder`,
		channelMapping: map[string]string{
			"sol":   "sol-channel",
			"terra": "terra-channel",
	placeholder,
		accountModelMapping: map[string]any{
			"sol":           "gpt-5.6-sol",
			"terra":         "gpt-5.6-terra",
			"sol-channel":   "gpt-5.6-sol",
			"terra-channel": "gpt-5.6-terra",
	placeholder,
placeholder)

	require.Len(t, got.upstreamPayloads, 2)
	require.Equal(t, "sol-channel", gjson.GetBytes(got.upstreamPayloads[0], "model").String())
	require.Equal(t, "terra-channel", gjson.GetBytes(got.upstreamPayloads[1], "model").String())
	require.Len(t, got.clientEvents, 2)
	require.Equal(t, "sol", gjson.GetBytes(got.clientEvents[0], "response.model").String())
	require.Equal(t, "terra", gjson.GetBytes(got.clientEvents[1], "response.model").String())

	require.Len(t, got.logs, 2)
	require.Equal(t, "sol", got.logs[0].Model)
	require.Equal(t, "sol", got.logs[0].RequestedModel)
	require.NotNil(t, got.logs[0].UpstreamModel)
	require.Equal(t, "sol-channel", *got.logs[0].UpstreamModel)
	require.NotNil(t, got.logs[0].ModelMappingChain)
	require.Equal(t, "sol→sol-channel", *got.logs[0].ModelMappingChain)

	require.Equal(t, "terra", got.logs[1].Model)
	require.Equal(t, "terra", got.logs[1].RequestedModel)
	require.NotNil(t, got.logs[1].UpstreamModel)
	require.Equal(t, "terra-channel", *got.logs[1].UpstreamModel)
	require.NotNil(t, got.logs[1].ModelMappingChain)
	require.Equal(t, "terra→terra-channel", *got.logs[1].ModelMappingChain)
	require.InDelta(t, got.logs[1].TotalCost*2.5, got.logs[0].TotalCost, 1e-12,
		"each turn must be billed with its own channel-mapped model")
placeholder

func TestOpenAIResponsesWebSocket_UnchangedChannelTargetOutsideAccountMappingKeysRemainsValid(t *testing.T) {
	got := runOpenAIResponsesWebSocketUsageLogCase(t, openAIResponsesWSUsageLogCase{
		firstPayload:  `{"type":"response.create","model":"public-alias","stream":falseplaceholder`,
		secondPayload: `{"type":"response.create","stream":falseplaceholder`,
		channelMapping: map[string]string{
			"public-alias": "gpt-5.6-sol",
	placeholder,
		accountModelMapping: map[string]any{
			"public-alias": "gpt-5.6-terra",
	placeholder,
placeholder)

	require.Len(t, got.upstreamPayloads, 2)
	require.Equal(t, "gpt-5.6-sol", gjson.GetBytes(got.upstreamPayloads[0], "model").String())
	require.Equal(t, "gpt-5.6-sol", gjson.GetBytes(got.upstreamPayloads[1], "model").String())
	require.Len(t, got.clientEvents, 2)
	require.Equal(t, "public-alias", gjson.GetBytes(got.clientEvents[0], "response.model").String())
	require.Equal(t, "public-alias", gjson.GetBytes(got.clientEvents[1], "response.model").String())
	require.Len(t, got.logs, 2)
	for _, usageLog := range got.logs {
		require.Equal(t, "public-alias", usageLog.RequestedModel)
		require.NotNil(t, usageLog.UpstreamModel)
		require.Equal(t, "gpt-5.6-sol", *usageLog.UpstreamModel)
		require.NotNil(t, usageLog.ModelMappingChain)
		require.Equal(t, "public-alias→gpt-5.6-sol", *usageLog.ModelMappingChain)
placeholder
placeholder

func TestOpenAIResponsesWebSocket_PassthroughKeepsTurnMappingSnapshot(t *testing.T) {
	got := runOpenAIResponsesWebSocketUsageLogCase(t, openAIResponsesWSUsageLogCase{
		firstPayload:  `{"type":"response.create","model":"sol","stream":falseplaceholder`,
		secondPayload: `{"type":"response.create","model":"sol","stream":falseplaceholder`,
		channelMapping: map[string]string{
			"sol": "gpt-5.6-sol",
	placeholder,
		afterFirstUpstreamRequest: func(channelSvc *service.ChannelService) error {
			if channelSvc == nil {
				return errors.New("channel service is nil")
		placeholder
			_, err := channelSvc.Update(context.Background(), 7701, &service.UpdateChannelInput{
				ModelMapping: map[string]map[string]string{
					service.PlatformOpenAI: {"sol": "gpt-5.6-terra"placeholder,
			placeholder,
		placeholder)
			return err
	placeholder,
placeholder)

	require.Len(t, got.upstreamPayloads, 2)
	require.Equal(t, "gpt-5.6-sol", gjson.GetBytes(got.upstreamPayloads[0], "model").String())
	require.Equal(t, "gpt-5.6-terra", gjson.GetBytes(got.upstreamPayloads[1], "model").String())

	require.Len(t, got.logs, 2)
	require.Equal(t, "sol", got.logs[0].Model)
	require.NotNil(t, got.logs[0].UpstreamModel)
	require.Equal(t, "gpt-5.6-sol", *got.logs[0].UpstreamModel)
	require.NotNil(t, got.logs[0].ModelMappingChain)
	require.Equal(t, "sol→gpt-5.6-sol", *got.logs[0].ModelMappingChain)
	require.InDelta(t, 40e-6, got.logs[0].TotalCost, 1e-12,
		"the in-flight turn must retain the channel-mapped billing model used when it was sent")

	require.Equal(t, "sol", got.logs[1].Model)
	require.NotNil(t, got.logs[1].UpstreamModel)
	require.Equal(t, "gpt-5.6-terra", *got.logs[1].UpstreamModel)
	require.NotNil(t, got.logs[1].ModelMappingChain)
	require.Equal(t, "sol→gpt-5.6-terra", *got.logs[1].ModelMappingChain)
	require.InDelta(t, got.logs[1].TotalCost*2.5, got.logs[0].TotalCost, 1e-12,
		"the next turn must use the updated channel mapping")
placeholder

func TestOpenAIResponsesWebSocket_CtxPoolAppliesPerTurnMappingAndPreservesRequestedModel(t *testing.T) {
	got := runOpenAIResponsesWebSocketUsageLogCase(t, openAIResponsesWSUsageLogCase{
		firstPayload:       `{"type":"response.create","model":"gpt-5.6-sol","stream":falseplaceholder`,
		secondPayload:      `{"type":"response.create","model":"gpt-5.6-terra","stream":falseplaceholder`,
		ingressMode:        service.OpenAIWSIngressModeCtxPool,
		billingModelSource: service.BillingModelSourceRequested,
		channelMapping: map[string]string{
			"gpt-5.6-terra": "gpt-5.6-sol",
	placeholder,
placeholder)

	require.Len(t, got.upstreamPayloads, 2)
	require.Equal(t, "gpt-5.6-sol", gjson.GetBytes(got.upstreamPayloads[0], "model").String())
	require.Equal(t, "gpt-5.6-sol", gjson.GetBytes(got.upstreamPayloads[1], "model").String())
	require.Len(t, got.clientEvents, 2)
	require.Equal(t, "gpt-5.6-sol", gjson.GetBytes(got.clientEvents[0], "response.model").String())
	require.Equal(t, "gpt-5.6-terra", gjson.GetBytes(got.clientEvents[1], "response.model").String())

	require.Len(t, got.logs, 2)
	require.Equal(t, "gpt-5.6-sol", got.logs[0].RequestedModel)
	require.Nil(t, got.logs[0].ModelMappingChain)
	require.InDelta(t, 40e-6, got.logs[0].TotalCost, 1e-12)
	require.Equal(t, "gpt-5.6-terra", got.logs[1].RequestedModel)
	require.NotNil(t, got.logs[1].ModelMappingChain)
	require.Equal(t, "gpt-5.6-terra→gpt-5.6-sol", *got.logs[1].ModelMappingChain)
	require.InDelta(t, 16e-6, got.logs[1].TotalCost, 1e-12,
		"BillingModelSourceRequested must use the client model before channel mapping")
placeholder

func TestOpenAIWSTurnBillingModelPreservesImagePricingModel(t *testing.T) {
	tests := []struct {
		name             string
		resultModel      string
		mapping          service.ChannelMappingResult
		requestedModel   string
		upstreamModel    string
		wantBillingModel string
placeholder{
		{
			name:             "upstream billing preserves image model",
			resultModel:      "gpt-image-2",
			mapping:          service.ChannelMappingResult{BillingModelSource: service.BillingModelSourceUpstreamplaceholder,
			requestedModel:   "gpt-5.6-sol",
			upstreamModel:    "gpt-5.6-sol",
			wantBillingModel: "gpt-image-2",
	placeholder,
		{
			name:             "unmapped channel preserves image model",
			resultModel:      "gpt-image-2",
			mapping:          service.ChannelMappingResult{MappedModel: "gpt-5.6-sol", BillingModelSource: service.BillingModelSourceChannelMappedplaceholder,
			requestedModel:   "gpt-5.6-sol",
			upstreamModel:    "gpt-5.6-sol",
			wantBillingModel: "gpt-image-2",
	placeholder,
		{
			name:             "requested source overrides image model",
			resultModel:      "gpt-image-2",
			mapping:          service.ChannelMappingResult{BillingModelSource: service.BillingModelSourceRequestedplaceholder,
			requestedModel:   "public-image-alias",
			upstreamModel:    "gpt-5.6-sol",
			wantBillingModel: "public-image-alias",
	placeholder,
		{
			name:             "mapped channel source overrides image model",
			resultModel:      "gpt-image-2",
			mapping:          service.ChannelMappingResult{MappedModel: "priced-channel-model", BillingModelSource: service.BillingModelSourceChannelMappedplaceholder,
			requestedModel:   "public-image-alias",
			upstreamModel:    "gpt-5.6-sol",
			wantBillingModel: "priced-channel-model",
	placeholder,
		{
			name:             "text turn falls back to upstream model",
			mapping:          service.ChannelMappingResult{BillingModelSource: service.BillingModelSourceUpstreamplaceholder,
			requestedModel:   "public-alias",
			upstreamModel:    "gpt-5.6-sol",
			wantBillingModel: "gpt-5.6-sol",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &service.OpenAIForwardResult{BillingModel: tt.resultModelplaceholder
			require.Equal(t, tt.wantBillingModel, openAIWSTurnBillingModel(result, tt.mapping, tt.requestedModel, tt.upstreamModel))
	placeholder)
placeholder
placeholder

func TestOpenAIAccountScheduleModelUsesActualOrSharedResolver(t *testing.T) {
	account := &service.Account{
		Platform: service.PlatformOpenAI,
		Type:     service.AccountTypeOAuth,
placeholder
			"model_mapping":         map[string]any{"public": "billing"placeholder,
			"compact_model_mapping": map[string]any{"public": "compact-actual"placeholder,
	placeholder,
placeholder

	reported := &service.OpenAIForwardResult{UpstreamModel: "observed-actual"placeholder
	require.Equal(t, "observed-actual", openAIAccountScheduleModel(nil, account, "public", true, reported))
	require.Equal(t, "compact-actual", openAIAccountScheduleModel(nil, account, "public", true, nil))
	require.Equal(t, "billing", openAIAccountScheduleModel(nil, account, "public", false, nil))

	c, _ := gin.CreateTestContext(nil)
	service.SetOpsUpstreamModel(c, "attempt-actual")
	require.Equal(t, "attempt-actual", openAIAccountScheduleModel(c, account, "public", true, nil))

	setOpsSelectedAccount(c, account.ID, account.Platform)
	require.Equal(t, "attempt-actual", openAIAccountScheduleModel(c, account, "public", true, nil))
placeholder

func TestShouldReportOpenAIWSProxyAccountFailure(t *testing.T) {
	t.Run("unsupported client model switch does not penalize account", func(t *testing.T) {
		err := fmt.Errorf("wrapped ingress turn: %w", newOpenAIWSUnsupportedModelSwitchError("gpt-unsupported"))
		require.False(t, shouldReportOpenAIWSProxyAccountFailure(err))

		var closeErr *service.OpenAIWSClientCloseError
		require.ErrorAs(t, err, &closeErr)
		require.Equal(t, coderws.StatusPolicyViolation, closeErr.StatusCode())
		require.Equal(t, "model switch requires reconnect", closeErr.Reason())
placeholder)

	t.Run("upstream policy violation still penalizes account", func(t *testing.T) {
		err := service.NewOpenAIWSClientCloseError(
			coderws.StatusPolicyViolation,
			"upstream websocket authentication failed",
			errors.New("upstream rejected credentials"),
		)
		require.True(t, shouldReportOpenAIWSProxyAccountFailure(err))
placeholder)

	t.Run("generic proxy failure still penalizes account", func(t *testing.T) {
		require.True(t, shouldReportOpenAIWSProxyAccountFailure(errors.New("upstream websocket read failed")))
placeholder)
placeholder

func TestSetOpenAIClientTransportHTTP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	setOpenAIClientTransportHTTP(c)
	require.Equal(t, service.OpenAIClientTransportHTTP, service.GetOpenAIClientTransport(c))
placeholder

func TestSetOpenAIClientTransportWS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	setOpenAIClientTransportWS(c)
	require.Equal(t, service.OpenAIClientTransportWS, service.GetOpenAIClientTransport(c))
placeholder

// TestOpenAIHandler_GjsonExtraction 验证 gjson 从请求体中提取 model/stream 的正确性
func TestOpenAIHandler_GjsonExtraction(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantModel  string
		wantStream bool
placeholder{
		{"正常提取", `{"model":"gpt-4","stream":true,"input":"hello"placeholder`, "gpt-4", trueplaceholder,
		{"stream false", `{"model":"gpt-4","stream":falseplaceholder`, "gpt-4", falseplaceholder,
		{"无 stream 字段", `{"model":"gpt-4"placeholder`, "gpt-4", falseplaceholder,
		{"model 缺失", `{"stream":trueplaceholder`, "", trueplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := []byte(tt.body)
			modelResult := gjson.GetBytes(body, "model")
			model := ""
			if modelResult.Type == gjson.String {
				model = modelResult.String()
		placeholder
			stream := gjson.GetBytes(body, "stream").Bool()
			require.Equal(t, tt.wantModel, model)
			require.Equal(t, tt.wantStream, stream)
	placeholder)
placeholder
placeholder

// TestOpenAIHandler_GjsonValidation 验证修复后的 JSON 合法性和类型校验
func TestOpenAIHandler_GjsonValidation(t *testing.T) {
	// 非法 JSON 被 gjson.ValidBytes 拦截
	require.False(t, gjson.ValidBytes([]byte(`{invalid json`)))

	// model 为数字 → 类型不是 gjson.String，应被拒绝
	body := []byte(`{"model":placeholder`)
	modelResult := gjson.GetBytes(body, "model")
	require.True(t, modelResult.Exists())
	require.NotEqual(t, gjson.String, modelResult.Type)

	// model 为 null → 类型不是 gjson.String，应被拒绝
	body2 := []byte(`{"model":nullplaceholder`)
	modelResult2 := gjson.GetBytes(body2, "model")
	require.True(t, modelResult2.Exists())
	require.NotEqual(t, gjson.String, modelResult2.Type)

	// stream 为 string → 类型既不是 True 也不是 False，应被拒绝
	body3 := []byte(`{"model":"gpt-4","stream":"true"placeholder`)
	streamResult := gjson.GetBytes(body3, "stream")
	require.True(t, streamResult.Exists())
	require.NotEqual(t, gjson.True, streamResult.Type)
	require.NotEqual(t, gjson.False, streamResult.Type)

	// stream 为 int → 同上
	body4 := []byte(`{"model":"gpt-4","stream":1placeholder`)
	streamResult2 := gjson.GetBytes(body4, "stream")
	require.True(t, streamResult2.Exists())
	require.NotEqual(t, gjson.True, streamResult2.Type)
	require.NotEqual(t, gjson.False, streamResult2.Type)
placeholder

// TestOpenAIHandler_InstructionsInjection 验证 instructions 的 gjson/sjson 注入逻辑
func TestOpenAIHandler_InstructionsInjection(t *testing.T) {
	// 测试 1：无 instructions → 注入
	body := []byte(`{"model":"gpt-4"placeholder`)
	existing := gjson.GetBytes(body, "instructions").String()
	require.Empty(t, existing)
	newBody, err := sjson.SetBytes(body, "instructions", "test instruction")
placeholder
	require.Equal(t, "test instruction", gjson.GetBytes(newBody, "instructions").String())

	// 测试 2：已有 instructions → 不覆盖
	body2 := []byte(`{"model":"gpt-4","instructions":"existing"placeholder`)
	existing2 := gjson.GetBytes(body2, "instructions").String()
	require.Equal(t, "existing", existing2)

	// 测试 3：空白 instructions → 注入
	body3 := []byte(`{"model":"gpt-4","instructions":"   "placeholder`)
	existing3 := strings.TrimSpace(gjson.GetBytes(body3, "instructions").String())
	require.Empty(t, existing3)

	// 测试 4：sjson.SetBytes 返回错误时不应 panic
	// 正常 JSON 不会产生 sjson 错误，验证返回值被正确处理
	validBody := []byte(`{"model":"gpt-4"placeholder`)
	result, setErr := sjson.SetBytes(validBody, "instructions", "hello")
	require.NoError(t, setErr)
	require.True(t, gjson.ValidBytes(result))
placeholder

func newOpenAIHandlerForPreviousResponseIDValidation(t *testing.T, cache *concurrencyCacheMock) *OpenAIGatewayHandler {
placeholder
	if cache == nil {
		cache = &concurrencyCacheMock{
			acquireUserSlotFn: func(ctx context.Context, userID int64, maxConcurrency int, requestID string) (bool, error) {
				return true, nil
		placeholder,
			acquireAccountSlotFn: func(ctx context.Context, accountID int64, maxConcurrency int, requestID string) (bool, error) {
				return true, nil
		placeholder,
	placeholder
placeholder
	return &OpenAIGatewayHandler{
		gatewayService:      &service.OpenAIGatewayService{placeholder,
		billingCacheService: &service.BillingCacheService{placeholder,
		apiKeyService:       &service.APIKeyService{placeholder,
		concurrencyHelper:   NewConcurrencyHelper(service.NewConcurrencyService(cache), SSEPingFormatNone, time.Second),
placeholder
placeholder

func newOpenAIWSHandlerTestServer(t *testing.T, h *OpenAIGatewayHandler, subject middleware.AuthSubject) *httptest.Server {
placeholder
	groupID := int64(2)
	apiKey := &service.APIKey{
		ID:      101,
		GroupID: &groupID,
		User:    &service.User{ID: subject.UserIDplaceholder,
placeholder
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(string(middleware.ContextKeyAPIKey), apiKey)
		c.Set(string(middleware.ContextKeyUser), subject)
		c.Next()
placeholder)
	router.GET("/openai/v1/responses", h.ResponsesWebSocket)
	return httptest.NewServer(router)
placeholder

type openAIResponsesWSUsageLogCase struct {
	firstPayload              string
	secondPayload             string
	userAgent                 *string
	ingressMode               string
	channelMapping            map[string]string
	billingModelSource        string
	accountModelMapping       map[string]any
	afterFirstUpstreamRequest func(channelSvc *service.ChannelService) error
placeholder

type openAIResponsesWSUsageLogResult struct {
	log                  *service.UsageLog
	logs                 []*service.UsageLog
	upstreamFirstPayload []byte
	upstreamPayloads     [][]byte
	clientEvents         [][]byte
placeholder

type openAIWSUsageHandlerAccountRepoStub struct {
	service.AccountRepository
	account service.Account
placeholder

func (s *openAIWSUsageHandlerAccountRepoStub) ListSchedulableByPlatform(ctx context.Context, platform string) ([]service.Account, error) {
	if s.account.Platform != platform {
		return nil, nil
placeholder
	return []service.Account{s.accountplaceholder, nil
placeholder

func (s *openAIWSUsageHandlerAccountRepoStub) ListSchedulableByGroupIDAndPlatform(ctx context.Context, groupID int64, platform string) ([]service.Account, error) {
	return s.ListSchedulableByPlatform(ctx, platform)
placeholder

func (s *openAIWSUsageHandlerAccountRepoStub) GetByID(ctx context.Context, id int64) (*service.Account, error) {
	if s.account.ID != id {
		return nil, nil
placeholder
	account := s.account
	return &account, nil
placeholder

type openAIWSFailoverHandlerAccountRepoStub struct {
	service.AccountRepository
	accounts       []service.Account
	rateLimitedIDs []int64
placeholder

type openAIHTTPPassthroughFailoverUpstream struct {
	service.HTTPUpstream
	mu         sync.Mutex
	accountIDs []int64
placeholder

func (u *openAIHTTPPassthroughFailoverUpstream) Do(_ *http.Request, _ string, accountID int64, _ int) (*http.Response, error) {
	u.mu.Lock()
	u.accountIDs = append(u.accountIDs, accountID)
	u.mu.Unlock()
	return &http.Response{
		StatusCode: http.StatusBadGateway,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"temporary upstream failure"placeholderplaceholder`)),
placeholder, nil
placeholder

func (u *openAIHTTPPassthroughFailoverUpstream) calls() []int64 {
	u.mu.Lock()
	defer u.mu.Unlock()
	return append([]int64(nil), u.accountIDs...)
placeholder

type openAIHTTPPassthroughAuthFailoverUpstream struct {
	service.HTTPUpstream
	mu         sync.Mutex
	accountIDs []int64
	statusCode int
placeholder

func (u *openAIHTTPPassthroughAuthFailoverUpstream) Do(_ *http.Request, _ string, accountID int64, _ int) (*http.Response, error) {
	u.mu.Lock()
	u.accountIDs = append(u.accountIDs, accountID)
	u.mu.Unlock()
	if accountID == 9911 {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"id":"resp_healthy","object":"response","model":"gpt-5.2","status":"completed","usage":{"input_tokens":1,"output_tokens":1,"total_tokens":2placeholderplaceholder`)),
	placeholder, nil
placeholder
	return &http.Response{
		StatusCode: u.statusCode,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"upstream credential rejected"placeholderplaceholder`)),
placeholder, nil
placeholder

func (u *openAIHTTPPassthroughAuthFailoverUpstream) calls() []int64 {
	u.mu.Lock()
	defer u.mu.Unlock()
	return append([]int64(nil), u.accountIDs...)
placeholder

type openAIHTTPPassthroughSSERateLimitUpstream struct {
	service.HTTPUpstream
	mu         sync.Mutex
	accountIDs []int64
placeholder

func (u *openAIHTTPPassthroughSSERateLimitUpstream) Do(_ *http.Request, _ string, accountID int64, _ int) (*http.Response, error) {
	u.mu.Lock()
	u.accountIDs = append(u.accountIDs, accountID)
	u.mu.Unlock()
	body := strings.Join([]string{
		"event: response.created",
		`data: {"type":"response.created","response":{"id":"resp_rate_limited"placeholderplaceholder`,
		"",
		"event: response.failed",
		`data: {"type":"response.failed","response":{"id":"resp_rate_limited","status":"failed","error":{"code":"rate_limit_exceeded","message":"Concurrency limit exceeded for account, please retry later"placeholderplaceholderplaceholder`,
		"",
placeholder, "\n")
	return &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"text/event-stream"placeholder,
			"Retry-After":  []string{"1"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(body)),
placeholder, nil
placeholder

func (u *openAIHTTPPassthroughSSERateLimitUpstream) calls() []int64 {
	u.mu.Lock()
	defer u.mu.Unlock()
	return append([]int64(nil), u.accountIDs...)
placeholder

func (s *openAIWSFailoverHandlerAccountRepoStub) ListSchedulableByPlatform(ctx context.Context, platform string) ([]service.Account, error) {
	out := make([]service.Account, 0, len(s.accounts))
	for _, account := range s.accounts {
		if account.Platform == platform && account.IsSchedulable() {
			out = append(out, account)
	placeholder
placeholder
	return out, nil
placeholder

func (s *openAIWSFailoverHandlerAccountRepoStub) ListSchedulableByGroupIDAndPlatform(ctx context.Context, groupID int64, platform string) ([]service.Account, error) {
	return s.ListSchedulableByPlatform(ctx, platform)
placeholder

func (s *openAIWSFailoverHandlerAccountRepoStub) ListSchedulableUngroupedByPlatform(ctx context.Context, platform string) ([]service.Account, error) {
	return s.ListSchedulableByPlatform(ctx, platform)
placeholder

func (s *openAIWSFailoverHandlerAccountRepoStub) GetByID(ctx context.Context, id int64) (*service.Account, error) {
	for _, account := range s.accounts {
		if account.ID == id {
			acc := account
			return &acc, nil
	placeholder
placeholder
	return nil, nil
placeholder

func (s *openAIWSFailoverHandlerAccountRepoStub) SetRateLimited(ctx context.Context, id int64, resetAt time.Time) error {
	s.rateLimitedIDs = append(s.rateLimitedIDs, id)
	for i := range s.accounts {
		if s.accounts[i].ID == id {
			reset := resetAt
			s.accounts[i].RateLimitResetAt = &reset
			break
	placeholder
placeholder
	return nil
placeholder

type openAIWSUsageHandlerUsageLogRepoStub struct {
	service.UsageLogRepository
	created chan *service.UsageLog
placeholder

func (s *openAIWSUsageHandlerUsageLogRepoStub) Create(ctx context.Context, log *service.UsageLog) (bool, error) {
	if s.created != nil {
		s.created <- log
placeholder
	return true, nil
placeholder

type openAIWSUsageHandlerChannelRepoStub struct {
	service.ChannelRepository
	mu             sync.Mutex
	channels       []service.Channel
	groupPlatforms map[int64]string
placeholder

func (s *openAIWSUsageHandlerChannelRepoStub) ListAll(ctx context.Context) ([]service.Channel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]service.Channel, 0, len(s.channels))
	for i := range s.channels {
		out = append(out, *s.channels[i].Clone())
placeholder
	return out, nil
placeholder

func (s *openAIWSUsageHandlerChannelRepoStub) GetByID(ctx context.Context, id int64) (*service.Channel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.channels {
		if s.channels[i].ID == id {
			return s.channels[i].Clone(), nil
	placeholder
placeholder
	return nil, service.ErrChannelNotFound
placeholder

func (s *openAIWSUsageHandlerChannelRepoStub) Update(ctx context.Context, channel *service.Channel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.channels {
		if s.channels[i].ID == channel.ID {
			s.channels[i] = *channel.Clone()
			return nil
	placeholder
placeholder
	return service.ErrChannelNotFound
placeholder

func (s *openAIWSUsageHandlerChannelRepoStub) GetGroupIDs(ctx context.Context, channelID int64) ([]int64, error) {
	channel, err := s.GetByID(ctx, channelID)
	if err != nil {
		return nil, err
placeholder
	return append([]int64(nil), channel.GroupIDs...), nil
placeholder

func (s *openAIWSUsageHandlerChannelRepoStub) GetGroupPlatforms(ctx context.Context, groupIDs []int64) (map[int64]string, error) {
	out := make(map[int64]string, len(groupIDs))
	for _, groupID := range groupIDs {
		if platform := strings.TrimSpace(s.groupPlatforms[groupID]); platform != "" {
			out[groupID] = platform
	placeholder
placeholder
	return out, nil
placeholder

func TestOpenAIResponses_APIKeyPassthroughPool5xxRetriesThenExhaustsMaxSwitches(t *testing.T) {
	gin.SetMode(gin.TestMode)
	groupID := int64(4203)
	accounts := []service.Account{
		{
			ID: 9910, Name: "pool-api-key", Platform: service.PlatformOpenAI,
			Type: service.AccountTypeAPIKey, Status: service.StatusActive, Schedulable: true, Priority: 1,
	placeholder
				"api_key":                      "sk-pool",
				"base_url":                     "https://api.example.test",
				"pool_mode":                    true,
				"pool_mode_retry_count":        float64(1),
				"pool_mode_retry_status_codes": []any{float64(http.StatusBadGateway)placeholder,
		placeholder,
			Extra: map[string]any{"openai_passthrough": trueplaceholder,
	placeholder,
		{
			ID: 9911, Name: "fallback-api-key", Platform: service.PlatformOpenAI,
			Type: service.AccountTypeAPIKey, Status: service.StatusActive, Schedulable: true, Priority: 2,
	placeholder
				"api_key":  "sk-fallback",
				"base_url": "https://api.example.test",
		placeholder,
			Extra: map[string]any{"openai_passthrough": trueplaceholder,
	placeholder,
placeholder
	cfg := &config.Config{RunMode: config.RunModeSimpleplaceholder
	cfg.Default.RateMultiplier = 1
	cfg.Security.URLAllowlist.Enabled = false
	cfg.Gateway.MaxAccountSwitches = 1

	accountRepo := &openAIWSFailoverHandlerAccountRepoStub{accounts: accountsplaceholder
	upstream := &openAIHTTPPassthroughFailoverUpstream{placeholder
	billingCacheSvc := service.NewBillingCacheService(nil, nil, nil, nil, nil, nil, cfg, nil)
	t.Cleanup(billingCacheSvc.Stop)
	gatewaySvc := service.NewOpenAIGatewayService(
		accountRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		cfg,
		nil,
		nil,
		service.NewBillingService(cfg, nil),
		nil,
		billingCacheSvc,
		upstream,
		&service.DeferredService{placeholder,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	h := NewOpenAIGatewayHandler(
		gatewaySvc,
		service.NewConcurrencyService(nil),
		billingCacheSvc,
		service.NewAPIKeyService(nil, nil, nil, nil, nil, nil, cfg),
		nil,
		nil,
		nil,
		nil,
		cfg,
	)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", strings.NewReader(`{"model":"gpt-5.2","input":"hello","stream":falseplaceholder`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(string(middleware.ContextKeyAPIKey), &service.APIKey{
		ID: 1803, GroupID: &groupID,
		User:  &service.User{ID: 1703, Status: service.StatusActiveplaceholder,
		Group: &service.Group{ID: groupID, Platform: service.PlatformOpenAI, Status: service.StatusActiveplaceholder,
placeholder)
	c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: 1703, Concurrency: 0placeholder)

	h.Responses(c)

	require.Equal(t, []int64{9910, 9910, 9911placeholder, upstream.calls())
	require.Equal(t, http.StatusBadGateway, rec.Code)
	require.Equal(t, "upstream_error", gjson.GetBytes(rec.Body.Bytes(), "error.type").String())
	require.Equal(t, "Upstream service temporarily unavailable", gjson.GetBytes(rec.Body.Bytes(), "error.message").String())
placeholder

func TestOpenAIResponses_APIKeyPassthroughPoolAuthFailureRetriesThenSwitchesToHealthyAccount(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
placeholder{
		{name: "401", statusCode: http.StatusUnauthorizedplaceholder,
		{name: "403", statusCode: http.StatusForbiddenplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			groupID := int64(4203)
			accounts := []service.Account{
				{
					ID: 9910, Name: "pool-api-key", Platform: service.PlatformOpenAI,
					Type: service.AccountTypeAPIKey, Status: service.StatusActive, Schedulable: true, Priority: 1,
			placeholder
						"api_key":                      "sk-pool",
						"base_url":                     "https://api.example.test",
						"pool_mode":                    true,
						"pool_mode_retry_count":        float64(1),
						"pool_mode_retry_status_codes": []any{float64(tt.statusCode)placeholder,
				placeholder,
					Extra: map[string]any{"openai_passthrough": trueplaceholder,
			placeholder,
				{
					ID: 9911, Name: "fallback-api-key", Platform: service.PlatformOpenAI,
					Type: service.AccountTypeAPIKey, Status: service.StatusActive, Schedulable: true, Priority: 2,
			placeholder
						"api_key":  "sk-fallback",
						"base_url": "https://api.example.test",
				placeholder,
					Extra: map[string]any{"openai_passthrough": trueplaceholder,
			placeholder,
		placeholder
			cfg := &config.Config{RunMode: config.RunModeSimpleplaceholder
			cfg.Default.RateMultiplier = 1
			cfg.Security.URLAllowlist.Enabled = false
			cfg.Gateway.MaxAccountSwitches = 1

			accountRepo := &openAIWSFailoverHandlerAccountRepoStub{accounts: accountsplaceholder
			upstream := &openAIHTTPPassthroughAuthFailoverUpstream{statusCode: tt.statusCodeplaceholder
			rateLimitSvc := service.NewRateLimitService(accountRepo, nil, cfg, nil, nil)
			billingCacheSvc := service.NewBillingCacheService(nil, nil, nil, nil, nil, nil, cfg, nil)
			t.Cleanup(billingCacheSvc.Stop)
			gatewaySvc := service.NewOpenAIGatewayService(
				accountRepo,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				cfg,
				nil,
				nil,
				service.NewBillingService(cfg, nil),
				rateLimitSvc,
				billingCacheSvc,
				upstream,
				&service.DeferredService{placeholder,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
			)
			h := NewOpenAIGatewayHandler(
				gatewaySvc,
				service.NewConcurrencyService(nil),
				billingCacheSvc,
				service.NewAPIKeyService(nil, nil, nil, nil, nil, nil, cfg),
				nil,
				nil,
				nil,
				nil,
				cfg,
			)

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", strings.NewReader(`{"model":"gpt-5.2","input":"hello","stream":falseplaceholder`))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Set(string(middleware.ContextKeyAPIKey), &service.APIKey{
				ID: 1803, GroupID: &groupID,
				User:  &service.User{ID: 1703, Status: service.StatusActiveplaceholder,
				Group: &service.Group{ID: groupID, Platform: service.PlatformOpenAI, Status: service.StatusActiveplaceholder,
		placeholder)
			c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: 1703, Concurrency: 0placeholder)

			h.Responses(c)

			require.Equal(t, []int64{9910, 9910, 9911placeholder, upstream.calls())
			require.Equal(t, http.StatusOK, rec.Code)
			require.Equal(t, "resp_healthy", gjson.GetBytes(rec.Body.Bytes(), "id").String())
	placeholder)
placeholder
placeholder

func TestOpenAIResponses_APIKeyPassthroughSSERateLimitUsesConfiguredPoolRetry(t *testing.T) {
	gin.SetMode(gin.TestMode)
	groupID := int64(4204)
	accounts := []service.Account{
		{
			ID: 9912, Name: "pool-sse-rate-limit", Platform: service.PlatformOpenAI,
			Type: service.AccountTypeAPIKey, Status: service.StatusActive, Schedulable: true, Priority: 1,
	placeholder
				"api_key":                      "sk-pool",
				"base_url":                     "https://api.example.test",
				"pool_mode":                    true,
				"pool_mode_retry_count":        float64(1),
				"pool_mode_retry_status_codes": []any{float64(http.StatusTooManyRequests)placeholder,
		placeholder,
			Extra: map[string]any{"openai_passthrough": trueplaceholder,
	placeholder,
placeholder
	cfg := &config.Config{RunMode: config.RunModeSimpleplaceholder
	cfg.Default.RateMultiplier = 1
	cfg.Security.URLAllowlist.Enabled = false
	cfg.Gateway.MaxAccountSwitches = 1

	accountRepo := &openAIWSFailoverHandlerAccountRepoStub{accounts: accountsplaceholder
	upstream := &openAIHTTPPassthroughSSERateLimitUpstream{placeholder
	billingCacheSvc := service.NewBillingCacheService(nil, nil, nil, nil, nil, nil, cfg, nil)
	t.Cleanup(billingCacheSvc.Stop)
	gatewaySvc := service.NewOpenAIGatewayService(
		accountRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		cfg,
		nil,
		nil,
		service.NewBillingService(cfg, nil),
		nil,
		billingCacheSvc,
		upstream,
		&service.DeferredService{placeholder,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	h := NewOpenAIGatewayHandler(
		gatewaySvc,
		service.NewConcurrencyService(nil),
		billingCacheSvc,
		service.NewAPIKeyService(nil, nil, nil, nil, nil, nil, cfg),
		nil,
		nil,
		nil,
		nil,
		cfg,
	)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", strings.NewReader(`{"model":"gpt-5.6-sol","input":"hello","stream":trueplaceholder`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(string(middleware.ContextKeyAPIKey), &service.APIKey{
		ID: 1804, GroupID: &groupID,
		User:  &service.User{ID: 1704, Status: service.StatusActiveplaceholder,
		Group: &service.Group{ID: groupID, Platform: service.PlatformOpenAI, Status: service.StatusActiveplaceholder,
placeholder)
	c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: 1704, Concurrency: 0placeholder)

	h.Responses(c)

	require.Equal(t, []int64{9912, 9912placeholder, upstream.calls())
	require.Equal(t, http.StatusTooManyRequests, rec.Code)
	require.Equal(t, "1", rec.Header().Get("Retry-After"))
	require.Equal(t, "rate_limit_error", gjson.GetBytes(rec.Body.Bytes(), "error.type").String())
	require.Equal(t, "Upstream rate limit exceeded, please retry later", gjson.GetBytes(rec.Body.Bytes(), "error.message").String())
placeholder

func TestOpenAIResponsesWebSocket_FailoverOnUpstreamUsageLimitEvent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	firstHitCh := make(chan []byte, 1)
	secondHitCh := make(chan []byte, 1)

	firstUpstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := coderws.Accept(w, r, &coderws.AcceptOptions{CompressionMode: coderws.CompressionContextTakeoverplaceholder)
		if err != nil {
			return
	placeholder
		defer func() { _ = conn.CloseNow() placeholder()

		readCtx, cancelRead := context.WithTimeout(r.Context(), 3*time.Second)
		_, payload, readErr := conn.Read(readCtx)
		cancelRead()
		if readErr == nil {
			firstHitCh <- payload
	placeholder

		writeCtx, cancelWrite := context.WithTimeout(r.Context(), 3*time.Second)
		_ = conn.Write(writeCtx, coderws.MessageText, []byte(`{"type":"error","error":{"code":"rate_limit_exceeded","type":"usage_limit_reached","message":"The usage limit has been reached"placeholderplaceholder`))
		cancelWrite()
placeholder))
	defer firstUpstream.Close()

	secondUpstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := coderws.Accept(w, r, &coderws.AcceptOptions{CompressionMode: coderws.CompressionContextTakeoverplaceholder)
		if err != nil {
			return
	placeholder
		defer func() { _ = conn.CloseNow() placeholder()

		readCtx, cancelRead := context.WithTimeout(r.Context(), 3*time.Second)
		_, payload, readErr := conn.Read(readCtx)
		cancelRead()
		if readErr == nil {
			secondHitCh <- payload
	placeholder

		writeCtx, cancelWrite := context.WithTimeout(r.Context(), 3*time.Second)
		_ = conn.Write(writeCtx, coderws.MessageText, []byte(`{"type":"response.completed","response":{"id":"resp_ws_failover_ok","model":"gpt-5.1","usage":{"input_tokens":1,"output_tokens":1placeholderplaceholderplaceholder`))
		cancelWrite()
		_ = conn.Close(coderws.StatusNormalClosure, "done")
placeholder))
	defer secondUpstream.Close()

	groupID := int64(4202)
	accounts := []service.Account{
		{
			ID:          9902,
			Name:        "openai-ws-rate-limited",
			Platform:    service.PlatformOpenAI,
			Type:        service.AccountTypeAPIKey,
			Status:      service.StatusActive,
			Schedulable: true,
			Concurrency: 1,
			Priority:    1,
	placeholder
				"api_key":  "sk-first",
				"base_url": firstUpstream.URL,
		placeholder,
			Extra: map[string]any{
				"openai_apikey_responses_websockets_v2_enabled": true,
				"openai_apikey_responses_websockets_v2_mode":    service.OpenAIWSIngressModePassthrough,
		placeholder,
	placeholder,
		{
			ID:          9903,
			Name:        "openai-ws-healthy",
			Platform:    service.PlatformOpenAI,
			Type:        service.AccountTypeAPIKey,
			Status:      service.StatusActive,
			Schedulable: true,
			Concurrency: 1,
			Priority:    2,
	placeholder
				"api_key":  "sk-second",
				"base_url": secondUpstream.URL,
		placeholder,
			Extra: map[string]any{
				"openai_apikey_responses_websockets_v2_enabled": true,
				"openai_apikey_responses_websockets_v2_mode":    service.OpenAIWSIngressModePassthrough,
		placeholder,
	placeholder,
placeholder

	cfg := &config.Config{placeholder
	cfg.RunMode = config.RunModeSimple
	cfg.Default.RateMultiplier = 1
	cfg.Security.URLAllowlist.Enabled = false
	cfg.Security.URLAllowlist.AllowInsecureHTTP = true
	cfg.Gateway.OpenAIWS.Enabled = true
	cfg.Gateway.OpenAIWS.APIKeyEnabled = true
	cfg.Gateway.OpenAIWS.ResponsesWebsocketsV2 = true
	cfg.Gateway.OpenAIWS.ModeRouterV2Enabled = true
	cfg.Gateway.OpenAIWS.DialTimeoutSeconds = 3
	cfg.Gateway.OpenAIWS.ReadTimeoutSeconds = 3
	cfg.Gateway.OpenAIWS.WriteTimeoutSeconds = 3
	cfg.Gateway.MaxAccountSwitches = 3

	accountRepo := &openAIWSFailoverHandlerAccountRepoStub{accounts: accountsplaceholder
	rateLimitSvc := service.NewRateLimitService(accountRepo, nil, cfg, nil, nil)
	billingCacheSvc := service.NewBillingCacheService(nil, nil, nil, nil, nil, nil, cfg, nil)
	gatewaySvc := service.NewOpenAIGatewayService(
		accountRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		cfg,
		nil,
		nil,
		service.NewBillingService(cfg, nil),
		rateLimitSvc,
		billingCacheSvc,
		nil,
		&service.DeferredService{placeholder,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	cache := &concurrencyCacheMock{
		acquireUserSlotFn: func(ctx context.Context, userID int64, maxConcurrency int, requestID string) (bool, error) {
			return true, nil
	placeholder,
		acquireAccountSlotFn: func(ctx context.Context, accountID int64, maxConcurrency int, requestID string) (bool, error) {
			return true, nil
	placeholder,
placeholder
	h := &OpenAIGatewayHandler{
		gatewayService:      gatewaySvc,
		billingCacheService: billingCacheSvc,
		apiKeyService:       &service.APIKeyService{placeholder,
		concurrencyHelper:   NewConcurrencyHelper(service.NewConcurrencyService(cache), SSEPingFormatNone, time.Second),
		maxAccountSwitches:  3,
placeholder

	apiKey := &service.APIKey{
		ID:      1802,
		GroupID: &groupID,
		User:    &service.User{ID: 1702, Status: service.StatusActiveplaceholder,
		Group:   &service.Group{ID: groupID, Platform: service.PlatformOpenAI, Status: service.StatusActiveplaceholder,
placeholder
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(string(middleware.ContextKeyAPIKey), apiKey)
		c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: apiKey.User.ID, Concurrency: 1placeholder)
		c.Next()
placeholder)
	router.GET("/openai/v1/responses", h.ResponsesWebSocket)
	handlerServer := httptest.NewServer(router)
	defer handlerServer.Close()

	dialCtx, cancelDial := context.WithTimeout(context.Background(), 3*time.Second)
	clientConn, _, err := coderws.Dial(
		dialCtx,
		"ws"+strings.TrimPrefix(handlerServer.URL, "http")+"/openai/v1/responses",
		&coderws.DialOptions{CompressionMode: coderws.CompressionContextTakeoverplaceholder,
	)
	cancelDial()
placeholder
	defer func() { _ = clientConn.CloseNow() placeholder()

	writeCtx, cancelWrite := context.WithTimeout(context.Background(), 3*time.Second)
	err = clientConn.Write(writeCtx, coderws.MessageText, []byte(`{"type":"response.create","model":"gpt-5.1","stream":falseplaceholder`))
	cancelWrite()
placeholder

	readCtx, cancelRead := context.WithTimeout(context.Background(), 5*time.Second)
	_, event, err := clientConn.Read(readCtx)
	cancelRead()
placeholder
	require.Equal(t, "response.completed", gjson.GetBytes(event, "type").String())
	require.Equal(t, "resp_ws_failover_ok", gjson.GetBytes(event, "response.id").String())

	select {
	case <-firstHitCh:
	case <-time.After(3 * time.Second):
		t.Fatal("等待第一个上游收到首帧超时")
placeholder
	select {
	case <-secondHitCh:
	case <-time.After(3 * time.Second):
		t.Fatal("等待第二个上游收到重放首帧超时")
placeholder
	require.Equal(t, []int64{int64(9902)placeholder, accountRepo.rateLimitedIDs)
placeholder

func TestOpenAIResponsesWebSocket_FirstOutputTimeoutWithoutDownstreamReusesClientForOneFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)

	firstHitCh := make(chan []byte, 1)
	secondHitCh := make(chan []byte, 1)
	var firstConnections atomic.Int32
	var secondConnections atomic.Int32

	firstUpstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		firstConnections.Add(1)
		conn, err := coderws.Accept(w, r, &coderws.AcceptOptions{CompressionMode: coderws.CompressionContextTakeoverplaceholder)
		if err != nil {
			return
	placeholder
		defer func() { _ = conn.CloseNow() placeholder()

		readCtx, cancelRead := context.WithTimeout(r.Context(), 3*time.Second)
		_, payload, readErr := conn.Read(readCtx)
		cancelRead()
		if readErr == nil {
			firstHitCh <- payload
	placeholder

		select {
		case <-r.Context().Done():
		case <-time.After(3 * time.Second):
	placeholder
placeholder))
	defer firstUpstream.Close()

	secondUpstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secondConnections.Add(1)
		conn, err := coderws.Accept(w, r, &coderws.AcceptOptions{CompressionMode: coderws.CompressionContextTakeoverplaceholder)
		if err != nil {
			return
	placeholder
		defer func() { _ = conn.CloseNow() placeholder()

		readCtx, cancelRead := context.WithTimeout(r.Context(), 3*time.Second)
		_, payload, readErr := conn.Read(readCtx)
		cancelRead()
		if readErr == nil {
			secondHitCh <- payload
	placeholder

		for _, event := range []string{
			`{"type":"response.created","response":{"id":"resp_ws_timeout_b","model":"gpt-5.1"placeholderplaceholder`,
			`{"type":"response.output_text.delta","response_id":"resp_ws_timeout_b","delta":"recovered"placeholder`,
			`{"type":"response.completed","response":{"id":"resp_ws_timeout_b","model":"gpt-5.1","usage":{"input_tokens":1,"output_tokens":1placeholderplaceholderplaceholder`,
	placeholder {
			writeCtx, cancelWrite := context.WithTimeout(r.Context(), 3*time.Second)
			writeErr := conn.Write(writeCtx, coderws.MessageText, []byte(event))
			cancelWrite()
			if writeErr != nil {
				return
		placeholder
	placeholder
		readCtx, cancelRead = context.WithTimeout(r.Context(), 3*time.Second)
		_, _, _ = conn.Read(readCtx)
		cancelRead()
placeholder))
	defer secondUpstream.Close()

	groupID := int64(4212)
	accounts := []service.Account{
		{
			ID:          9912,
			Name:        "openai-ws-first-semantic-timeout",
			Platform:    service.PlatformOpenAI,
			Type:        service.AccountTypeAPIKey,
			Status:      service.StatusActive,
			Schedulable: true,
			Concurrency: 1,
			Priority:    1,
	placeholder"api_key": "sk-first", "base_url": firstUpstream.URLplaceholder,
			Extra: map[string]any{
				"openai_apikey_responses_websockets_v2_enabled": true,
				"openai_apikey_responses_websockets_v2_mode":    service.OpenAIWSIngressModePassthrough,
		placeholder,
	placeholder,
		{
			ID:          9913,
			Name:        "openai-ws-failover-healthy",
			Platform:    service.PlatformOpenAI,
			Type:        service.AccountTypeAPIKey,
			Status:      service.StatusActive,
			Schedulable: true,
			Concurrency: 1,
			Priority:    2,
	placeholder"api_key": "sk-second", "base_url": secondUpstream.URLplaceholder,
			Extra: map[string]any{
				"openai_apikey_responses_websockets_v2_enabled": true,
				"openai_apikey_responses_websockets_v2_mode":    service.OpenAIWSIngressModePassthrough,
		placeholder,
	placeholder,
placeholder

	cfg := &config.Config{placeholder
	cfg.RunMode = config.RunModeSimple
	cfg.Default.RateMultiplier = 1
	cfg.Security.URLAllowlist.Enabled = false
	cfg.Security.URLAllowlist.AllowInsecureHTTP = true
	cfg.Gateway.OpenAIFirstOutputTimeoutSeconds = 1
	cfg.Gateway.OpenAIWS.Enabled = true
	cfg.Gateway.OpenAIWS.APIKeyEnabled = true
	cfg.Gateway.OpenAIWS.ResponsesWebsocketsV2 = true
	cfg.Gateway.OpenAIWS.ModeRouterV2Enabled = true
	cfg.Gateway.OpenAIWS.DialTimeoutSeconds = 3
	cfg.Gateway.OpenAIWS.ReadTimeoutSeconds = 3
	cfg.Gateway.OpenAIWS.WriteTimeoutSeconds = 3
	cfg.Gateway.OpenAIWS.IngressInterTurnIdleTimeoutSeconds = 3
	cfg.Gateway.MaxAccountSwitches = 3

	accountRepo := &openAIWSFailoverHandlerAccountRepoStub{accounts: accountsplaceholder
	rateLimitSvc := service.NewRateLimitService(accountRepo, nil, cfg, nil, nil)
	billingCacheSvc := service.NewBillingCacheService(nil, nil, nil, nil, nil, nil, cfg, nil)
	gatewaySvc := service.NewOpenAIGatewayService(
		accountRepo, nil, nil, nil, nil, nil, nil, cfg, nil, nil,
		service.NewBillingService(cfg, nil), rateLimitSvc, billingCacheSvc,
		nil, &service.DeferredService{placeholder, nil, nil, nil, nil, nil, nil, nil,
	)
	cache := &concurrencyCacheMock{
		acquireUserSlotFn: func(context.Context, int64, int, string) (bool, error) { return true, nil placeholder,
		acquireAccountSlotFn: func(context.Context, int64, int, string) (bool, error) {
			return true, nil
	placeholder,
placeholder
	h := &OpenAIGatewayHandler{
		gatewayService:      gatewaySvc,
		billingCacheService: billingCacheSvc,
		apiKeyService:       &service.APIKeyService{placeholder,
		concurrencyHelper:   NewConcurrencyHelper(service.NewConcurrencyService(cache), SSEPingFormatNone, time.Second),
		maxAccountSwitches:  3,
placeholder

	apiKey := &service.APIKey{
		ID:      1812,
		GroupID: &groupID,
		User:    &service.User{ID: 1712, Status: service.StatusActiveplaceholder,
		Group:   &service.Group{ID: groupID, Platform: service.PlatformOpenAI, Status: service.StatusActiveplaceholder,
placeholder
	handlerDone := make(chan struct{placeholder)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(string(middleware.ContextKeyAPIKey), apiKey)
		c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: apiKey.User.ID, Concurrency: 1placeholder)
		c.Next()
placeholder)
	router.GET("/openai/v1/responses", func(c *gin.Context) {
		h.ResponsesWebSocket(c)
		close(handlerDone)
placeholder)
	handlerServer := httptest.NewServer(router)
	defer handlerServer.Close()

	dialCtx, cancelDial := context.WithTimeout(context.Background(), 3*time.Second)
	clientConn, _, err := coderws.Dial(
		dialCtx,
		"ws"+strings.TrimPrefix(handlerServer.URL, "http")+"/openai/v1/responses",
		&coderws.DialOptions{CompressionMode: coderws.CompressionContextTakeoverplaceholder,
	)
	cancelDial()
placeholder
	defer func() { _ = clientConn.CloseNow() placeholder()

	writeCtx, cancelWrite := context.WithTimeout(context.Background(), 3*time.Second)
	err = clientConn.Write(writeCtx, coderws.MessageText, []byte(`{"type":"response.create","model":"gpt-5.1","stream":falseplaceholder`))
	cancelWrite()
placeholder

	var eventTypes []string
	readCtx, cancelRead := context.WithTimeout(context.Background(), 6*time.Second)
	for {
		_, event, readErr := clientConn.Read(readCtx)
		require.NoError(t, readErr)
		eventType := gjson.GetBytes(event, "type").String()
		eventTypes = append(eventTypes, eventType)
		if eventType == "response.completed" {
			require.Equal(t, "resp_ws_timeout_b", gjson.GetBytes(event, "response.id").String())
			break
	placeholder
placeholder
	cancelRead()
	require.Contains(t, eventTypes, "response.output_text.delta")
	require.NoError(t, clientConn.Close(coderws.StatusNormalClosure, "done"))

	select {
	case <-handlerDone:
	case <-time.After(3 * time.Second):
		t.Fatal("websocket handler did not finish after healthy failover turn")
placeholder
	select {
	case <-firstHitCh:
	case <-time.After(3 * time.Second):
		t.Fatal("first upstream did not receive replayable request")
placeholder
	select {
	case <-secondHitCh:
	case <-time.After(3 * time.Second):
		t.Fatal("second upstream did not receive replayed request")
placeholder
	require.Equal(t, int32(1), firstConnections.Load())
	require.Equal(t, int32(1), secondConnections.Load())
	require.NotContains(t, accountRepo.rateLimitedIDs, int64(9913), "healthy failover account must not be penalized")
placeholder

func runOpenAIResponsesWebSocketUsageLogCase(t *testing.T, tc openAIResponsesWSUsageLogCase) openAIResponsesWSUsageLogResult {
placeholder
	gin.SetMode(gin.TestMode)

	turnCount := 1
	if strings.TrimSpace(tc.secondPayload) != "" {
		turnCount = 2
placeholder
	upstreamPayloadCh := make(chan []byte, turnCount)
	upstreamErrCh := make(chan error, 1)
	var channelSvc *service.ChannelService
	upstreamServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := coderws.Accept(w, r, &coderws.AcceptOptions{
			CompressionMode: coderws.CompressionContextTakeover,
	placeholder)
		if err != nil {
			upstreamErrCh <- err
			return
	placeholder
		defer func() {
			_ = conn.CloseNow()
	placeholder()

		for turn := 1; turn <= turnCount; turn++ {
			readCtx, cancelRead := context.WithTimeout(r.Context(), 3*time.Second)
			msgType, payload, readErr := conn.Read(readCtx)
			cancelRead()
			if readErr != nil {
				upstreamErrCh <- readErr
				return
		placeholder
			if msgType != coderws.MessageText && msgType != coderws.MessageBinary {
				upstreamErrCh <- errors.New("unexpected upstream websocket message type")
				return
		placeholder
			upstreamPayloadCh <- payload
			if turn == 1 && tc.afterFirstUpstreamRequest != nil {
				if callbackErr := tc.afterFirstUpstreamRequest(channelSvc); callbackErr != nil {
					upstreamErrCh <- callbackErr
					return
			placeholder
		placeholder

			response := fmt.Sprintf(
				`{"type":"response.completed","response":{"id":"resp_usage_e2e_%d","model":%q,"usage":{"input_tokens":2,"output_tokens":1placeholderplaceholderplaceholder`,
				turn,
				gjson.GetBytes(payload, "model").String(),
			)
			writeCtx, cancelWrite := context.WithTimeout(r.Context(), 3*time.Second)
			writeErr := conn.Write(writeCtx, coderws.MessageText, []byte(response))
			cancelWrite()
			if writeErr != nil {
				upstreamErrCh <- writeErr
				return
		placeholder
	placeholder
		upstreamErrCh <- nil
placeholder))
	defer upstreamServer.Close()

	groupID := int64(4201)
	account := service.Account{
		ID:          9901,
		Name:        "openai-ws-passthrough-usage-e2e",
		Platform:    service.PlatformOpenAI,
		Type:        service.AccountTypeAPIKey,
		Status:      service.StatusActive,
		Schedulable: true,
		Concurrency: 1,
placeholder
			"api_key":       "sk-test",
			"base_url":      upstreamServer.URL,
			"model_mapping": tc.accountModelMapping,
	placeholder,
		Extra: map[string]any{
			"openai_apikey_responses_websockets_v2_enabled": true,
			"openai_apikey_responses_websockets_v2_mode":    service.OpenAIWSIngressModePassthrough,
	placeholder,
placeholder
	if strings.TrimSpace(tc.ingressMode) != "" {
		account.Extra["openai_apikey_responses_websockets_v2_mode"] = tc.ingressMode
placeholder

	cfg := &config.Config{placeholder
	cfg.RunMode = config.RunModeSimple
	cfg.Default.RateMultiplier = 1
	cfg.Security.URLAllowlist.Enabled = false
	cfg.Security.URLAllowlist.AllowInsecureHTTP = true
	cfg.Gateway.OpenAIWS.Enabled = true
	cfg.Gateway.OpenAIWS.APIKeyEnabled = true
	cfg.Gateway.OpenAIWS.ResponsesWebsocketsV2 = true
	cfg.Gateway.OpenAIWS.ModeRouterV2Enabled = true
	cfg.Gateway.OpenAIWS.DialTimeoutSeconds = 3
	cfg.Gateway.OpenAIWS.ReadTimeoutSeconds = 3
	cfg.Gateway.OpenAIWS.WriteTimeoutSeconds = 3

	accountRepo := &openAIWSUsageHandlerAccountRepoStub{account: accountplaceholder
	usageRepo := &openAIWSUsageHandlerUsageLogRepoStub{created: make(chan *service.UsageLog, turnCount)placeholder

	if len(tc.channelMapping) > 0 {
		channelSvc = service.NewChannelService(&openAIWSUsageHandlerChannelRepoStub{
			channels: []service.Channel{{
				ID:                 7701,
				Name:               "openai-ws-e2e-channel",
				Status:             service.StatusActive,
				GroupIDs:           []int64{groupIDplaceholder,
				ModelMapping:       map[string]map[string]string{service.PlatformOpenAI: tc.channelMappingplaceholder,
				BillingModelSource: tc.billingModelSource,
	placeholder
			groupPlatforms: map[int64]string{groupID: service.PlatformOpenAIplaceholder,
	placeholder, nil, nil, nil)
placeholder

	billingCacheSvc := service.NewBillingCacheService(nil, nil, nil, nil, nil, nil, cfg, nil)
	gatewaySvc := service.NewOpenAIGatewayService(
		accountRepo,
		usageRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
		cfg,
		nil,
		nil,
		service.NewBillingService(cfg, nil),
		nil,
		billingCacheSvc,
		nil,
		&service.DeferredService{placeholder,
		nil,
		nil,
		nil,
		channelSvc,
		nil,
		nil,
		nil, // userPlatformQuotaRepo
	)

	cache := &concurrencyCacheMock{
		acquireUserSlotFn: func(ctx context.Context, userID int64, maxConcurrency int, requestID string) (bool, error) {
			return true, nil
	placeholder,
		acquireAccountSlotFn: func(ctx context.Context, accountID int64, maxConcurrency int, requestID string) (bool, error) {
			return true, nil
	placeholder,
placeholder
	h := &OpenAIGatewayHandler{
		gatewayService:      gatewaySvc,
		billingCacheService: billingCacheSvc,
		apiKeyService:       &service.APIKeyService{placeholder,
		concurrencyHelper:   NewConcurrencyHelper(service.NewConcurrencyService(cache), SSEPingFormatNone, time.Second),
placeholder

	apiKey := &service.APIKey{
		ID:      1801,
		GroupID: &groupID,
		User:    &service.User{ID: 1701, Status: service.StatusActiveplaceholder,
placeholder
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(string(middleware.ContextKeyAPIKey), apiKey)
		c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: apiKey.User.ID, Concurrency: 1placeholder)
		c.Next()
placeholder)
	router.GET("/openai/v1/responses", h.ResponsesWebSocket)
	handlerServer := httptest.NewServer(router)
	defer handlerServer.Close()

	headers := http.Header{placeholder
	if tc.userAgent != nil {
		headers.Set("User-Agent", *tc.userAgent)
placeholder
	dialCtx, cancelDial := context.WithTimeout(context.Background(), 3*time.Second)
	clientConn, _, err := coderws.Dial(
		dialCtx,
		"ws"+strings.TrimPrefix(handlerServer.URL, "http")+"/openai/v1/responses",
		&coderws.DialOptions{HTTPHeader: headers, CompressionMode: coderws.CompressionContextTakeoverplaceholder,
	)
	cancelDial()
placeholder
	defer func() {
		_ = clientConn.CloseNow()
placeholder()

	writeCtx, cancelWrite := context.WithTimeout(context.Background(), 3*time.Second)
	err = clientConn.Write(writeCtx, coderws.MessageText, []byte(tc.firstPayload))
	cancelWrite()
placeholder

	clientEvents := make([][]byte, 0, turnCount)
	readCompleted := func() {
		readCtx, cancelRead := context.WithTimeout(context.Background(), 3*time.Second)
		_, event, readErr := clientConn.Read(readCtx)
		cancelRead()
		require.NoError(t, readErr)
		require.Equal(t, "response.completed", gjson.GetBytes(event, "type").String())
		clientEvents = append(clientEvents, append([]byte(nil), event...))
placeholder
	readCompleted()
	if turnCount == 2 {
		writeCtx, cancelWrite = context.WithTimeout(context.Background(), 3*time.Second)
		err = clientConn.Write(writeCtx, coderws.MessageText, []byte(tc.secondPayload))
		cancelWrite()
	placeholder
		readCompleted()
placeholder
	_ = clientConn.Close(coderws.StatusNormalClosure, "done")

	usageLogs := make([]*service.UsageLog, 0, turnCount)
	for len(usageLogs) < turnCount {
		select {
		case usageLog := <-usageRepo.created:
			require.NotNil(t, usageLog)
			usageLogs = append(usageLogs, usageLog)
		case <-time.After(3 * time.Second):
			t.Fatal("等待 WebSocket usage log 写入超时")
	placeholder
placeholder

	upstreamPayloads := make([][]byte, 0, turnCount)
	for len(upstreamPayloads) < turnCount {
		select {
		case payload := <-upstreamPayloadCh:
			upstreamPayloads = append(upstreamPayloads, payload)
		case <-time.After(3 * time.Second):
			t.Fatal("等待上游 WebSocket 请求帧超时")
	placeholder
placeholder

	select {
	case upstreamErr := <-upstreamErrCh:
		require.NoError(t, upstreamErr)
	case <-time.After(3 * time.Second):
		t.Fatal("等待上游 WebSocket 结束超时")
placeholder

	return openAIResponsesWSUsageLogResult{
		log:                  usageLogs[0],
		logs:                 usageLogs,
		upstreamFirstPayload: upstreamPayloads[0],
		upstreamPayloads:     upstreamPayloads,
		clientEvents:         clientEvents,
placeholder
placeholder

func testStringPtr(v string) *string {
	return &v
placeholder

func TestOpenAIForwardErrorAlreadyCommunicated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("upstream response failed after write", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, EndpointResponses, nil)
		before := c.Writer.Size()
		_, _ = c.Writer.WriteString(`event: response.failed
data: {"type":"response.failed","error":{"message":"This content was flagged"placeholderplaceholder

`)

		reported := openAIForwardErrorAlreadyCommunicated(c, before, errors.New("upstream response failed: This content was flagged"))

		require.True(t, reported)
placeholder)

	t.Run("no write still needs fallback", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, EndpointResponses, nil)

		reported := openAIForwardErrorAlreadyCommunicated(c, c.Writer.Size(), errors.New("upstream response failed: This content was flagged"))

		require.False(t, reported)
placeholder)

	t.Run("generic error after write still needs fallback", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, EndpointResponses, nil)
		before := c.Writer.Size()
		_, _ = c.Writer.WriteString(":\n\n")

		reported := openAIForwardErrorAlreadyCommunicated(c, before, errors.New("stream read error: unexpected EOF"))

		require.False(t, reported)
placeholder)

	// H-2: cyber_policy 命中且响应已写出时，即便 err 前缀不在白名单（非流式 400 cyber
	// 返回 "openai cyber_policy:"、透传账号返回 "upstream error:"），也须判定已透传，避免
	// ensureForwardErrorResponse 在已写出的完整响应尾部追加 SSE 污染响应体。
	t.Run("cyber policy hit after write is already communicated", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, EndpointResponses, nil)
		service.MarkOpsCyberPolicy(c, service.CyberPolicyMark{Message: "blocked", UpstreamStatus: 400placeholder)
		before := c.Writer.Size()
		_, _ = c.Writer.WriteString(`{"error":{"code":"cyber_policy","message":"blocked"placeholderplaceholder`)

		require.True(t, openAIForwardErrorAlreadyCommunicated(c, before, errors.New("openai cyber_policy: blocked")))
placeholder)

	// Size 守卫优先于 cyber 短路：cyber 命中但未写出任何响应时仍需补写错误。
	t.Run("cyber policy without write still needs fallback", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, EndpointResponses, nil)
		service.MarkOpsCyberPolicy(c, service.CyberPolicyMark{Message: "blocked", UpstreamStatus: 400placeholder)

		require.False(t, openAIForwardErrorAlreadyCommunicated(c, c.Writer.Size(), errors.New("openai cyber_policy: blocked")))
placeholder)
placeholder
