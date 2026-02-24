package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/antigravity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// antigravityFailingWriter 模拟客户端断开连接的 gin.ResponseWriter
type antigravityFailingWriter struct {
	gin.ResponseWriter
	failAfter int // 允许成功写入的次数，之后所有写入返回错误
	writes    int
placeholder

func (w *antigravityFailingWriter) Write(p []byte) (int, error) {
	if w.writes >= w.failAfter {
		return 0, errors.New("write failed: client disconnected")
placeholder
	w.writes++
	return w.ResponseWriter.Write(p)
placeholder

// newAntigravityTestService 创建用于流式测试的 AntigravityGatewayService
func newAntigravityTestService(cfg *config.Config) *AntigravityGatewayService {
	return &AntigravityGatewayService{
		settingService: &SettingService{cfg: cfgplaceholder,
placeholder
placeholder

func TestStripSignatureSensitiveBlocksFromClaudeRequest(t *testing.T) {
	req := &antigravity.ClaudeRequest{
		Model: "claude-sonnet-4-5",
		Thinking: &antigravity.ThinkingConfig{
			Type:         "enabled",
			BudgetTokens: 1024,
	placeholder,
		Messages: []antigravity.ClaudeMessage{
			{
				Role: "assistant",
				Content: json.RawMessage(`[
					{"type":"thinking","thinking":"secret plan","signature":""placeholder,
					{"type":"tool_use","id":"t1","name":"Bash","input":{"command":"ls"placeholderplaceholder
				]`),
		placeholder,
			{
				Role: "user",
				Content: json.RawMessage(`[
					{"type":"tool_result","tool_use_id":"t1","content":"ok","is_error":falseplaceholder,
					{"type":"redacted_thinking","data":"..."placeholder
				]`),
		placeholder,
	placeholder,
placeholder

	changed, err := stripSignatureSensitiveBlocksFromClaudeRequest(req)
placeholder
	require.True(t, changed)
	require.Nil(t, req.Thinking)

	require.Len(t, req.Messages, 2)

	var blocks0 []map[string]any
	require.NoError(t, json.Unmarshal(req.Messages[0].Content, &blocks0))
	require.Len(t, blocks0, 2)
	require.Equal(t, "text", blocks0[0]["type"])
	require.Equal(t, "secret plan", blocks0[0]["text"])
	require.Equal(t, "text", blocks0[1]["type"])

	var blocks1 []map[string]any
	require.NoError(t, json.Unmarshal(req.Messages[1].Content, &blocks1))
	require.Len(t, blocks1, 1)
	require.Equal(t, "text", blocks1[0]["type"])
	require.NotEmpty(t, blocks1[0]["text"])
placeholder

func TestStripThinkingFromClaudeRequest_DoesNotDowngradeTools(t *testing.T) {
	req := &antigravity.ClaudeRequest{
		Model: "claude-sonnet-4-5",
		Thinking: &antigravity.ThinkingConfig{
			Type:         "enabled",
			BudgetTokens: 1024,
	placeholder,
		Messages: []antigravity.ClaudeMessage{
			{
				Role:    "assistant",
				Content: json.RawMessage(`[{"type":"thinking","thinking":"secret plan"placeholder,{"type":"tool_use","id":"t1","name":"Bash","input":{"command":"ls"placeholderplaceholder]`),
		placeholder,
	placeholder,
placeholder

	changed, err := stripThinkingFromClaudeRequest(req)
placeholder
	require.True(t, changed)
	require.Nil(t, req.Thinking)

	var blocks []map[string]any
	require.NoError(t, json.Unmarshal(req.Messages[0].Content, &blocks))
	require.Len(t, blocks, 2)
	require.Equal(t, "text", blocks[0]["type"])
	require.Equal(t, "secret plan", blocks[0]["text"])
	require.Equal(t, "tool_use", blocks[1]["type"])
placeholder

func TestIsPromptTooLongError(t *testing.T) {
	require.True(t, isPromptTooLongError([]byte(`{"error":{"message":"Prompt is too long"placeholderplaceholder`)))
	require.True(t, isPromptTooLongError([]byte(`{"message":"Prompt is too long"placeholder`)))
	require.False(t, isPromptTooLongError([]byte(`{"error":{"message":"other"placeholderplaceholder`)))
placeholder

type httpUpstreamStub struct {
	resp *http.Response
	err  error
placeholder

func (s *httpUpstreamStub) Do(_ *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
	return s.resp, s.err
placeholder

func (s *httpUpstreamStub) DoWithTLS(_ *http.Request, _ string, _ int64, _ int, _ bool) (*http.Response, error) {
	return s.resp, s.err
placeholder

func TestAntigravityGatewayService_Forward_PromptTooLong(t *testing.T) {
	gin.SetMode(gin.TestMode)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)

	body, err := json.Marshal(map[string]any{
		"model": "claude-opus-4-6",
		"messages": []map[string]any{
			{"role": "user", "content": "hi"placeholder,
	placeholder,
		"max_tokens": 1,
		"stream":     false,
placeholder)
placeholder

	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	c.Request = req

	respBody := []byte(`{"error":{"message":"Prompt is too long"placeholderplaceholder`)
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     http.Header{"X-Request-Id": []string{"req-1"placeholderplaceholder,
		Body:       io.NopCloser(bytes.NewReader(respBody)),
placeholder

	svc := &AntigravityGatewayService{
		tokenProvider: &AntigravityTokenProvider{placeholder,
		httpUpstream:  &httpUpstreamStub{resp: respplaceholder,
placeholder

	account := &Account{
		ID:          1,
		Name:        "acc-1",
		Platform:    PlatformAntigravity,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"access_token": "token",
	placeholder,
placeholder

	result, err := svc.Forward(context.Background(), c, account, body, false)
	require.Nil(t, result)

	var promptErr *PromptTooLongError
	require.ErrorAs(t, err, &promptErr)
	require.Equal(t, http.StatusBadRequest, promptErr.StatusCode)
	require.Equal(t, "req-1", promptErr.RequestID)
	require.NotEmpty(t, promptErr.Body)

	raw, ok := c.Get(OpsUpstreamErrorsKey)
	require.True(t, ok)
	events, ok := raw.([]*OpsUpstreamErrorEvent)
	require.True(t, ok)
	require.Len(t, events, 1)
	require.Equal(t, "prompt_too_long", events[0].Kind)
placeholder

// TestAntigravityGatewayService_Forward_ModelRateLimitTriggersFailover
// 验证：当账号存在模型限流且剩余时间 >= antigravityRateLimitThreshold 时，
// Forward 方法应返回 UpstreamFailoverError，触发 Handler 切换账号
func TestAntigravityGatewayService_Forward_ModelRateLimitTriggersFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)

	body, err := json.Marshal(map[string]any{
		"model": "claude-opus-4-6",
		"messages": []map[string]any{
			{"role": "user", "content": "hi"placeholder,
	placeholder,
		"max_tokens": 1,
		"stream":     false,
placeholder)
placeholder

	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	c.Request = req

	// 不需要真正调用上游，因为预检查会直接返回切换信号
	svc := &AntigravityGatewayService{
		tokenProvider: &AntigravityTokenProvider{placeholder,
		httpUpstream:  &httpUpstreamStub{resp: nil, err: nilplaceholder,
placeholder

	// 设置模型限流：剩余时间 30 秒（> antigravityRateLimitThreshold 7s）
	futureResetAt := time.Now().Add(30 * time.Second).Format(time.RFC3339)
	account := &Account{
		ID:          1,
		Name:        "acc-rate-limited",
		Platform:    PlatformAntigravity,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"access_token": "token",
	placeholder,
		Extra: map[string]any{
			modelRateLimitsKey: map[string]any{
				"claude-opus-4-6-thinking": map[string]any{
					"rate_limit_reset_at": futureResetAt,
			placeholder,
		placeholder,
	placeholder,
placeholder

	result, err := svc.Forward(context.Background(), c, account, body, false)
	require.Nil(t, result, "Forward should not return result when model rate limited")
	require.NotNil(t, err, "Forward should return error")

	// 核心验证：错误应该是 UpstreamFailoverError，而不是普通 502 错误
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr, "error should be UpstreamFailoverError to trigger account switch")
	require.Equal(t, http.StatusServiceUnavailable, failoverErr.StatusCode)
	// 非粘性会话请求，ForceCacheBilling 应为 false
	require.False(t, failoverErr.ForceCacheBilling, "ForceCacheBilling should be false for non-sticky session")
placeholder

// TestAntigravityGatewayService_ForwardGemini_ModelRateLimitTriggersFailover
// 验证：ForwardGemini 方法同样能正确将 AntigravityAccountSwitchError 转换为 UpstreamFailoverError
func TestAntigravityGatewayService_ForwardGemini_ModelRateLimitTriggersFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)

	body, err := json.Marshal(map[string]any{
		"contents": []map[string]any{
			{"role": "user", "parts": []map[string]any{{"text": "hi"placeholderplaceholderplaceholder,
	placeholder,
placeholder)
placeholder

	req := httptest.NewRequest(http.MethodPost, "/v1beta/models/gemini-2.5-flash:generateContent", bytes.NewReader(body))
	c.Request = req

	// 不需要真正调用上游，因为预检查会直接返回切换信号
	svc := &AntigravityGatewayService{
		tokenProvider: &AntigravityTokenProvider{placeholder,
		httpUpstream:  &httpUpstreamStub{resp: nil, err: nilplaceholder,
placeholder

	// 设置模型限流：剩余时间 30 秒（> antigravityRateLimitThreshold 7s）
	futureResetAt := time.Now().Add(30 * time.Second).Format(time.RFC3339)
	account := &Account{
		ID:          2,
		Name:        "acc-gemini-rate-limited",
		Platform:    PlatformAntigravity,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"access_token": "token",
	placeholder,
		Extra: map[string]any{
			modelRateLimitsKey: map[string]any{
				"gemini-2.5-flash": map[string]any{
					"rate_limit_reset_at": futureResetAt,
			placeholder,
		placeholder,
	placeholder,
placeholder

	result, err := svc.ForwardGemini(context.Background(), c, account, "gemini-2.5-flash", "generateContent", false, body, false)
	require.Nil(t, result, "ForwardGemini should not return result when model rate limited")
	require.NotNil(t, err, "ForwardGemini should return error")

	// 核心验证：错误应该是 UpstreamFailoverError，而不是普通 502 错误
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr, "error should be UpstreamFailoverError to trigger account switch")
	require.Equal(t, http.StatusServiceUnavailable, failoverErr.StatusCode)
	// 非粘性会话请求，ForceCacheBilling 应为 false
	require.False(t, failoverErr.ForceCacheBilling, "ForceCacheBilling should be false for non-sticky session")
placeholder

// TestAntigravityGatewayService_Forward_StickySessionForceCacheBilling
// 验证：粘性会话切换时，UpstreamFailoverError.ForceCacheBilling 应为 true
func TestAntigravityGatewayService_Forward_StickySessionForceCacheBilling(t *testing.T) {
	gin.SetMode(gin.TestMode)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)

	body, err := json.Marshal(map[string]any{
		"model":    "claude-opus-4-6",
		"messages": []map[string]string{{"role": "user", "content": "hello"placeholderplaceholder,
placeholder)
placeholder

	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	c.Request = req

	svc := &AntigravityGatewayService{
		tokenProvider: &AntigravityTokenProvider{placeholder,
		httpUpstream:  &httpUpstreamStub{resp: nil, err: nilplaceholder,
placeholder

	// 设置模型限流：剩余时间 30 秒（> antigravityRateLimitThreshold 7s）
	futureResetAt := time.Now().Add(30 * time.Second).Format(time.RFC3339)
	account := &Account{
		ID:          3,
		Name:        "acc-sticky-rate-limited",
		Platform:    PlatformAntigravity,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"access_token": "token",
	placeholder,
		Extra: map[string]any{
			modelRateLimitsKey: map[string]any{
				"claude-opus-4-6-thinking": map[string]any{
					"rate_limit_reset_at": futureResetAt,
			placeholder,
		placeholder,
	placeholder,
placeholder

	// 传入 isStickySession = true
	result, err := svc.Forward(context.Background(), c, account, body, true)
	require.Nil(t, result, "Forward should not return result when model rate limited")
	require.NotNil(t, err, "Forward should return error")

	// 核心验证：粘性会话切换时，ForceCacheBilling 应为 true
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr, "error should be UpstreamFailoverError to trigger account switch")
	require.Equal(t, http.StatusServiceUnavailable, failoverErr.StatusCode)
	require.True(t, failoverErr.ForceCacheBilling, "ForceCacheBilling should be true for sticky session switch")
placeholder

// TestAntigravityGatewayService_ForwardGemini_StickySessionForceCacheBilling verifies
// that ForwardGemini sets ForceCacheBilling=true for sticky session switch.
func TestAntigravityGatewayService_ForwardGemini_StickySessionForceCacheBilling(t *testing.T) {
	gin.SetMode(gin.TestMode)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)

	body, err := json.Marshal(map[string]any{
		"contents": []map[string]any{
			{"role": "user", "parts": []map[string]any{{"text": "hi"placeholderplaceholderplaceholder,
	placeholder,
placeholder)
placeholder

	req := httptest.NewRequest(http.MethodPost, "/v1beta/models/gemini-2.5-flash:generateContent", bytes.NewReader(body))
	c.Request = req

	svc := &AntigravityGatewayService{
		tokenProvider: &AntigravityTokenProvider{placeholder,
		httpUpstream:  &httpUpstreamStub{resp: nil, err: nilplaceholder,
placeholder

	// 设置模型限流：剩余时间 30 秒（> antigravityRateLimitThreshold 7s）
	futureResetAt := time.Now().Add(30 * time.Second).Format(time.RFC3339)
	account := &Account{
		ID:          4,
		Name:        "acc-gemini-sticky-rate-limited",
		Platform:    PlatformAntigravity,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"access_token": "token",
	placeholder,
		Extra: map[string]any{
			modelRateLimitsKey: map[string]any{
				"gemini-2.5-flash": map[string]any{
					"rate_limit_reset_at": futureResetAt,
			placeholder,
		placeholder,
	placeholder,
placeholder

	// 传入 isStickySession = true
	result, err := svc.ForwardGemini(context.Background(), c, account, "gemini-2.5-flash", "generateContent", false, body, true)
	require.Nil(t, result, "ForwardGemini should not return result when model rate limited")
	require.NotNil(t, err, "ForwardGemini should return error")

	// 核心验证：粘性会话切换时，ForceCacheBilling 应为 true
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr, "error should be UpstreamFailoverError to trigger account switch")
	require.Equal(t, http.StatusServiceUnavailable, failoverErr.StatusCode)
	require.True(t, failoverErr.ForceCacheBilling, "ForceCacheBilling should be true for sticky session switch")
placeholder

// TestStreamUpstreamResponse_UsageAndFirstToken
// 验证：usage 字段可被累积/覆盖更新，并且能记录首 token 时间
func TestStreamUpstreamResponse_UsageAndFirstToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Header: http.Header{placeholder, Body: prplaceholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		fmt.Fprintln(pw, `data: {"usage":{"input_tokens":1,"output_tokens":2,"cache_read_input_tokens":3,"cache_creation_input_tokens":4placeholderplaceholder`)
		fmt.Fprintln(pw, `data: {"usage":{"output_tokens":5placeholderplaceholder`)
placeholder()

	start := time.Now().Add(-10 * time.Millisecond)
	result := svc.streamUpstreamResponse(c, resp, start)
	_ = pr.Close()

	require.NotNil(t, result)
	require.NotNil(t, result.usage)
	require.Equal(t, 1, result.usage.InputTokens)
	// 第二次事件覆盖 output_tokens
	require.Equal(t, 5, result.usage.OutputTokens)
	require.Equal(t, 3, result.usage.CacheReadInputTokens)
	require.Equal(t, 4, result.usage.CacheCreationInputTokens)
	require.NotNil(t, result.firstTokenMs)

	// 确保有透传输出
	require.Contains(t, rec.Body.String(), "data:")
placeholder

// --- 流式 happy path 测试 ---

// TestStreamUpstreamResponse_NormalComplete
// 验证：正常流式转发完成时，数据正确透传、usage 正确收集、clientDisconnect=false
func TestStreamUpstreamResponse_NormalComplete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		fmt.Fprintln(pw, `event: message_start`)
		fmt.Fprintln(pw, `data: {"type":"message_start","message":{"usage":{"input_tokens":10placeholderplaceholderplaceholder`)
		fmt.Fprintln(pw, "")
		fmt.Fprintln(pw, `event: content_block_delta`)
		fmt.Fprintln(pw, `data: {"type":"content_block_delta","delta":{"text":"hello"placeholderplaceholder`)
		fmt.Fprintln(pw, "")
		fmt.Fprintln(pw, `event: message_delta`)
		fmt.Fprintln(pw, `data: {"type":"message_delta","usage":{"output_tokens":5placeholderplaceholder`)
		fmt.Fprintln(pw, "")
placeholder()

	result := svc.streamUpstreamResponse(c, resp, time.Now())
	_ = pr.Close()

	require.NotNil(t, result)
	require.False(t, result.clientDisconnect, "normal completion should not set clientDisconnect")
	require.NotNil(t, result.usage)
	require.Equal(t, 5, result.usage.OutputTokens, "should collect output_tokens from message_delta")
	require.NotNil(t, result.firstTokenMs, "should record first token time")

	// 验证数据被透传到客户端
	body := rec.Body.String()
	require.Contains(t, body, "event: message_start")
	require.Contains(t, body, "content_block_delta")
	require.Contains(t, body, "message_delta")
placeholder

// TestHandleGeminiStreamingResponse_NormalComplete
// 验证：正常 Gemini 流式转发，数据正确透传、usage 正确收集
func TestHandleGeminiStreamingResponse_NormalComplete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		// 第一个 chunk（部分内容）
		fmt.Fprintln(pw, `data: {"candidates":[{"content":{"parts":[{"text":"Hello"placeholder]placeholderplaceholder],"usageMetadata":{"promptTokenCount":10,"candidatesTokenCount":3placeholderplaceholder`)
		fmt.Fprintln(pw, "")
		// 第二个 chunk（最终内容+完整 usage）
		fmt.Fprintln(pw, `data: {"candidates":[{"content":{"parts":[{"text":" world"placeholder]placeholder,"finishReason":"STOP"placeholder],"usageMetadata":{"promptTokenCount":10,"candidatesTokenCount":8,"cachedContentTokenCount":2placeholderplaceholder`)
		fmt.Fprintln(pw, "")
placeholder()

	result, err := svc.handleGeminiStreamingResponse(c, resp, time.Now())
	_ = pr.Close()

placeholder
	require.NotNil(t, result)
	require.False(t, result.clientDisconnect, "normal completion should not set clientDisconnect")
	require.NotNil(t, result.usage)
	// Gemini usage: promptTokenCount=10, candidatesTokenCount=8, cachedContentTokenCount=2
	// → InputTokens=10-2=8, OutputTokens=8, CacheReadInputTokens=2
	require.Equal(t, 8, result.usage.InputTokens)
	require.Equal(t, 8, result.usage.OutputTokens)
	require.Equal(t, 2, result.usage.CacheReadInputTokens)
	require.NotNil(t, result.firstTokenMs, "should record first token time")

	// 验证数据被透传到客户端
	body := rec.Body.String()
	require.Contains(t, body, "Hello")
	require.Contains(t, body, "world")
	// 不应包含错误事件
	require.NotContains(t, body, "event: error")
placeholder

// TestHandleClaudeStreamingResponse_NormalComplete
// 验证：正常 Claude 流式转发（Gemini→Claude 转换），数据正确转换并输出
func TestHandleClaudeStreamingResponse_NormalComplete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		// v1internal 包装格式：Gemini 数据嵌套在 "response" 字段下
		// ProcessLine 先尝试反序列化为 V1InternalResponse，裸格式会导致 Response.UsageMetadata 为空
		fmt.Fprintln(pw, `data: {"response":{"candidates":[{"content":{"parts":[{"text":"Hi there"placeholder]placeholder,"finishReason":"STOP"placeholder],"usageMetadata":{"promptTokenCount":5,"candidatesTokenCount":3placeholderplaceholderplaceholder`)
		fmt.Fprintln(pw, "")
placeholder()

	result, err := svc.handleClaudeStreamingResponse(c, resp, time.Now(), "claude-sonnet-4-5")
	_ = pr.Close()

placeholder
	require.NotNil(t, result)
	require.False(t, result.clientDisconnect, "normal completion should not set clientDisconnect")
	require.NotNil(t, result.usage)
	// Gemini→Claude 转换的 usage：promptTokenCount=5→InputTokens=5, candidatesTokenCount=3→OutputTokens=3
	require.Equal(t, 5, result.usage.InputTokens)
	require.Equal(t, 3, result.usage.OutputTokens)
	require.NotNil(t, result.firstTokenMs, "should record first token time")

	// 验证输出是 Claude SSE 格式（processor 会转换）
	body := rec.Body.String()
	require.Contains(t, body, "event: message_start", "should contain Claude message_start event")
	require.Contains(t, body, "event: message_stop", "should contain Claude message_stop event")
	// 不应包含错误事件
	require.NotContains(t, body, "event: error")
placeholder

// TestHandleGeminiStreamingResponse_ThoughtsTokenCount
// 验证：Gemini 流式转发时 thoughtsTokenCount 被计入 OutputTokens
func TestHandleGeminiStreamingResponse_ThoughtsTokenCount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		fmt.Fprintln(pw, `data: {"candidates":[{"content":{"parts":[{"text":"Hello"placeholder]placeholderplaceholder],"usageMetadata":{"promptTokenCount":100,"candidatesTokenCount":20,"thoughtsTokenCount":50placeholderplaceholder`)
		fmt.Fprintln(pw, "")
		fmt.Fprintln(pw, `data: {"candidates":[{"content":{"parts":[{"text":" world"placeholder]placeholder,"finishReason":"STOP"placeholder],"usageMetadata":{"promptTokenCount":100,"candidatesTokenCount":30,"thoughtsTokenCount":80,"cachedContentTokenCount":10placeholderplaceholder`)
		fmt.Fprintln(pw, "")
placeholder()

	result, err := svc.handleGeminiStreamingResponse(c, resp, time.Now())
	_ = pr.Close()

placeholder
	require.NotNil(t, result)
	require.NotNil(t, result.usage)
	// promptTokenCount=100, cachedContentTokenCount=10 → InputTokens=90
	require.Equal(t, 90, result.usage.InputTokens)
	// candidatesTokenCount=30 + thoughtsTokenCount=80 → OutputTokens=110
	require.Equal(t, 110, result.usage.OutputTokens)
	require.Equal(t, 10, result.usage.CacheReadInputTokens)
placeholder

// TestHandleClaudeStreamingResponse_ThoughtsTokenCount
// 验证：Gemini→Claude 流式转换时 thoughtsTokenCount 被计入 OutputTokens
func TestHandleClaudeStreamingResponse_ThoughtsTokenCount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		fmt.Fprintln(pw, `data: {"response":{"candidates":[{"content":{"parts":[{"text":"Hi"placeholder]placeholder,"finishReason":"STOP"placeholder],"usageMetadata":{"promptTokenCount":50,"candidatesTokenCount":10,"thoughtsTokenCount":25placeholderplaceholderplaceholder`)
		fmt.Fprintln(pw, "")
placeholder()

	result, err := svc.handleClaudeStreamingResponse(c, resp, time.Now(), "gemini-2.5-pro")
	_ = pr.Close()

placeholder
	require.NotNil(t, result)
	require.NotNil(t, result.usage)
	// promptTokenCount=50 → InputTokens=50
	require.Equal(t, 50, result.usage.InputTokens)
	// candidatesTokenCount=10 + thoughtsTokenCount=25 → OutputTokens=35
	require.Equal(t, 35, result.usage.OutputTokens)
placeholder

// --- 流式客户端断开检测测试 ---

// TestStreamUpstreamResponse_ClientDisconnectDrainsUsage
// 验证：客户端写入失败后，streamUpstreamResponse 继续读取上游以收集 usage
func TestStreamUpstreamResponse_ClientDisconnectDrainsUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	c.Writer = &antigravityFailingWriter{ResponseWriter: c.Writer, failAfter: 0placeholder

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		fmt.Fprintln(pw, `event: message_start`)
		fmt.Fprintln(pw, `data: {"type":"message_start","message":{"usage":{"input_tokens":10placeholderplaceholderplaceholder`)
		fmt.Fprintln(pw, "")
		fmt.Fprintln(pw, `event: message_delta`)
		fmt.Fprintln(pw, `data: {"type":"message_delta","usage":{"output_tokens":20placeholderplaceholder`)
		fmt.Fprintln(pw, "")
placeholder()

	result := svc.streamUpstreamResponse(c, resp, time.Now())
	_ = pr.Close()

	require.NotNil(t, result)
	require.True(t, result.clientDisconnect)
	require.NotNil(t, result.usage)
	require.Equal(t, 20, result.usage.OutputTokens)
placeholder

// TestStreamUpstreamResponse_ContextCanceled
// 验证：context 取消时返回 usage 且标记 clientDisconnect
func TestStreamUpstreamResponse_ContextCanceled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil).WithContext(ctx)

	resp := &http.Response{StatusCode: http.StatusOK, Body: cancelReadCloser{placeholder, Header: http.Header{placeholderplaceholder

	result := svc.streamUpstreamResponse(c, resp, time.Now())

	require.NotNil(t, result)
	require.True(t, result.clientDisconnect)
	require.NotContains(t, rec.Body.String(), "event: error")
placeholder

// TestStreamUpstreamResponse_Timeout
// 验证：上游超时时返回已收集的 usage
func TestStreamUpstreamResponse_Timeout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{StreamDataIntervalTimeout: 1, MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	result := svc.streamUpstreamResponse(c, resp, time.Now())
	_ = pw.Close()
	_ = pr.Close()

	require.NotNil(t, result)
	require.False(t, result.clientDisconnect)
placeholder

// TestStreamUpstreamResponse_TimeoutAfterClientDisconnect
// 验证：客户端断开后上游超时，返回 usage 并标记 clientDisconnect
func TestStreamUpstreamResponse_TimeoutAfterClientDisconnect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{StreamDataIntervalTimeout: 1, MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	c.Writer = &antigravityFailingWriter{ResponseWriter: c.Writer, failAfter: 0placeholder

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	go func() {
		fmt.Fprintln(pw, `data: {"type":"message_start","message":{"usage":{"input_tokens":5placeholderplaceholderplaceholder`)
		fmt.Fprintln(pw, "")
		// 不关闭 pw → 等待超时
placeholder()

	result := svc.streamUpstreamResponse(c, resp, time.Now())
	_ = pw.Close()
	_ = pr.Close()

	require.NotNil(t, result)
	require.True(t, result.clientDisconnect)
placeholder

// TestHandleGeminiStreamingResponse_ClientDisconnect
// 验证：Gemini 流式转发中客户端断开后继续 drain 上游
func TestHandleGeminiStreamingResponse_ClientDisconnect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	c.Writer = &antigravityFailingWriter{ResponseWriter: c.Writer, failAfter: 0placeholder

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		fmt.Fprintln(pw, `data: {"candidates":[{"content":{"parts":[{"text":"hi"placeholder]placeholderplaceholder],"usageMetadata":{"promptTokenCount":5,"candidatesTokenCount":10placeholderplaceholder`)
		fmt.Fprintln(pw, "")
placeholder()

	result, err := svc.handleGeminiStreamingResponse(c, resp, time.Now())
	_ = pr.Close()

placeholder
	require.NotNil(t, result)
	require.True(t, result.clientDisconnect)
	require.NotContains(t, rec.Body.String(), "write_failed")
placeholder

// TestHandleGeminiStreamingResponse_ContextCanceled
// 验证：context 取消时不注入错误事件
func TestHandleGeminiStreamingResponse_ContextCanceled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil).WithContext(ctx)

	resp := &http.Response{StatusCode: http.StatusOK, Body: cancelReadCloser{placeholder, Header: http.Header{placeholderplaceholder

	result, err := svc.handleGeminiStreamingResponse(c, resp, time.Now())

placeholder
	require.NotNil(t, result)
	require.True(t, result.clientDisconnect)
	require.NotContains(t, rec.Body.String(), "event: error")
placeholder

// TestHandleClaudeStreamingResponse_ClientDisconnect
// 验证：Claude 流式转发中客户端断开后继续 drain 上游
func TestHandleClaudeStreamingResponse_ClientDisconnect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	c.Writer = &antigravityFailingWriter{ResponseWriter: c.Writer, failAfter: 0placeholder

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		// v1internal 包装格式
		fmt.Fprintln(pw, `data: {"response":{"candidates":[{"content":{"parts":[{"text":"hello"placeholder]placeholder,"finishReason":"STOP"placeholder],"usageMetadata":{"promptTokenCount":8,"candidatesTokenCount":15placeholderplaceholderplaceholder`)
		fmt.Fprintln(pw, "")
placeholder()

	result, err := svc.handleClaudeStreamingResponse(c, resp, time.Now(), "claude-sonnet-4-5")
	_ = pr.Close()

placeholder
	require.NotNil(t, result)
	require.True(t, result.clientDisconnect)
placeholder

// TestHandleClaudeStreamingResponse_ContextCanceled
// 验证：context 取消时不注入错误事件
func TestHandleClaudeStreamingResponse_ContextCanceled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil).WithContext(ctx)

	resp := &http.Response{StatusCode: http.StatusOK, Body: cancelReadCloser{placeholder, Header: http.Header{placeholderplaceholder

	result, err := svc.handleClaudeStreamingResponse(c, resp, time.Now(), "claude-sonnet-4-5")

placeholder
	require.NotNil(t, result)
	require.True(t, result.clientDisconnect)
	require.NotContains(t, rec.Body.String(), "event: error")
placeholder

// TestExtractSSEUsage 验证 extractSSEUsage 从 SSE data 行正确提取 usage
func TestExtractSSEUsage(t *testing.T) {
	svc := &AntigravityGatewayService{placeholder
	tests := []struct {
		name     string
		line     string
		expected ClaudeUsage
placeholder{
		{
			name:     "message_delta with output_tokens",
			line:     `data: {"type":"message_delta","usage":{"output_tokens":42placeholderplaceholder`,
			expected: ClaudeUsage{OutputTokens: 42placeholder,
	placeholder,
		{
			name:     "non-data line ignored",
			line:     `event: message_start`,
			expected: ClaudeUsage{placeholder,
	placeholder,
		{
			name:     "top-level usage with all fields",
			line:     `data: {"usage":{"input_tokens":10,"output_tokens":20,"cache_read_input_tokens":5,"cache_creation_input_tokens":3placeholderplaceholder`,
			expected: ClaudeUsage{InputTokens: 10, OutputTokens: 20, CacheReadInputTokens: 5, CacheCreationInputTokens: 3placeholder,
	placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usage := &ClaudeUsage{placeholder
			svc.extractSSEUsage(tt.line, usage)
			require.Equal(t, tt.expected, *usage)
	placeholder)
placeholder
placeholder

// TestAntigravityClientWriter 验证 antigravityClientWriter 的断开检测
func TestAntigravityClientWriter(t *testing.T) {
	t.Run("normal write succeeds", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		flusher, _ := c.Writer.(http.Flusher)
		cw := newAntigravityClientWriter(c.Writer, flusher, "test")

		ok := cw.Write([]byte("hello"))
		require.True(t, ok)
		require.False(t, cw.Disconnected())
		require.Contains(t, rec.Body.String(), "hello")
placeholder)

	t.Run("write failure marks disconnected", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		fw := &antigravityFailingWriter{ResponseWriter: c.Writer, failAfter: 0placeholder
		flusher, _ := c.Writer.(http.Flusher)
		cw := newAntigravityClientWriter(fw, flusher, "test")

		ok := cw.Write([]byte("hello"))
		require.False(t, ok)
		require.True(t, cw.Disconnected())
placeholder)

	t.Run("subsequent writes are no-op", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		fw := &antigravityFailingWriter{ResponseWriter: c.Writer, failAfter: 0placeholder
		flusher, _ := c.Writer.(http.Flusher)
		cw := newAntigravityClientWriter(fw, flusher, "test")

		cw.Write([]byte("first"))
		ok := cw.Fprintf("second %d", 2)
		require.False(t, ok)
		require.True(t, cw.Disconnected())
placeholder)
placeholder

// TestUnwrapV1InternalResponse 测试 unwrapV1InternalResponse 的各种输入场景
func TestUnwrapV1InternalResponse(t *testing.T) {
	svc := &AntigravityGatewayService{placeholder

	// 构造 >50KB 的大型 JSON
	largePadding := strings.Repeat("x", 50*1024)
	largeInput := []byte(fmt.Sprintf(`{"response":{"id":"big","pad":"%s"placeholderplaceholder`, largePadding))
	largeExpected := fmt.Sprintf(`{"id":"big","pad":"%s"placeholder`, largePadding)

	tests := []struct {
		name     string
		input    []byte
		expected string
		wantErr  bool
placeholder{
		{
			name:     "正常 response 包装",
			input:    []byte(`{"response":{"id":"123","content":"hello"placeholderplaceholder`),
			expected: `{"id":"123","content":"hello"placeholder`,
	placeholder,
		{
			name:     "无 response 透传",
			input:    []byte(`{"id":"456"placeholder`),
			expected: `{"id":"456"placeholder`,
	placeholder,
		{
			name:     "空 JSON",
			input:    []byte(`{placeholder`),
			expected: `{placeholder`,
	placeholder,
		{
			name:     "response 为 null",
			input:    []byte(`{"response":nullplaceholder`),
			expected: `null`,
	placeholder,
		{
			name:     "response 为基础类型 string",
			input:    []byte(`{"response":"hello"placeholder`),
			expected: `"hello"`,
	placeholder,
		{
			name:     "非法 JSON",
			input:    []byte(`not json`),
			expected: `not json`,
	placeholder,
		{
			name:     "嵌套 response 只解一层",
			input:    []byte(`{"response":{"response":{"inner":trueplaceholderplaceholderplaceholder`),
			expected: `{"response":{"inner":trueplaceholderplaceholder`,
	placeholder,
		{
			name:     "大型 JSON >50KB",
			input:    largeInput,
			expected: largeExpected,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.unwrapV1InternalResponse(tt.input)
			if tt.wantErr {
			placeholder
				return
		placeholder
		placeholder
			require.Equal(t, tt.expected, strings.TrimSpace(string(got)))
	placeholder)
placeholder
placeholder

// --- unwrapV1InternalResponse benchmark 对照组 ---

// unwrapV1InternalResponseOld 旧实现：Unmarshal+Marshal 双重开销（仅用于 benchmark 对照）
func unwrapV1InternalResponseOld(body []byte) ([]byte, error) {
	var outer map[string]any
	if err := json.Unmarshal(body, &outer); err != nil {
		return nil, err
placeholder
	if resp, ok := outer["response"]; ok {
		return json.Marshal(resp)
placeholder
	return body, nil
placeholder

func BenchmarkUnwrapV1Internal_Old_Small(b *testing.B) {
	body := []byte(`{"response":{"candidates":[{"content":{"parts":[{"text":"hello world"placeholder]placeholderplaceholder],"usageMetadata":{"promptTokenCount":10,"candidatesTokenCount":5placeholderplaceholderplaceholder`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = unwrapV1InternalResponseOld(body)
placeholder
placeholder

func BenchmarkUnwrapV1Internal_New_Small(b *testing.B) {
	body := []byte(`{"response":{"candidates":[{"content":{"parts":[{"text":"hello world"placeholder]placeholderplaceholder],"usageMetadata":{"promptTokenCount":10,"candidatesTokenCount":5placeholderplaceholderplaceholder`)
	svc := &AntigravityGatewayService{placeholder
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.unwrapV1InternalResponse(body)
placeholder
placeholder

func BenchmarkUnwrapV1Internal_Old_Large(b *testing.B) {
	body := generateLargeUnwrapJSON(10 * 1024) // ~10KB
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = unwrapV1InternalResponseOld(body)
placeholder
placeholder

func BenchmarkUnwrapV1Internal_New_Large(b *testing.B) {
	body := generateLargeUnwrapJSON(10 * 1024) // ~10KB
	svc := &AntigravityGatewayService{placeholder
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.unwrapV1InternalResponse(body)
placeholder
placeholder

// generateLargeUnwrapJSON 生成指定最小大小的包含 response 包装的 JSON
func generateLargeUnwrapJSON(minSize int) []byte {
	parts := make([]map[string]string, 0)
	current := 0
	for current < minSize {
		text := fmt.Sprintf("这是第 %d 段内容，用于填充 JSON 到目标大小。", len(parts)+1)
		parts = append(parts, map[string]string{"text": textplaceholder)
		current += len(text) + 20 // 估算 JSON 编码开销
placeholder
	inner := map[string]any{
		"candidates": []map[string]any{
			{"content": map[string]any{"parts": partsplaceholderplaceholder,
	placeholder,
		"usageMetadata": map[string]any{
			"promptTokenCount":     100,
			"candidatesTokenCount": 50,
	placeholder,
placeholder
	outer := map[string]any{"response": innerplaceholder
	b, _ := json.Marshal(outer)
	return b
placeholder
