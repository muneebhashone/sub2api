//go:build unit

package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestBuildOpenAIChatCompletionsURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		base string
		want string
placeholder{
		// 已是 /chat/completions：原样返回
		{"already chat/completions", "https://api.openai.com/v1/chat/completions", "https://api.openai.com/v1/chat/completions"placeholder,
		// 以 /v1 结尾：追加 /chat/completions
		{"bare /v1", "https://api.openai.com/v1", "https://api.openai.com/v1/chat/completions"placeholder,
		// 其他情况：追加 /v1/chat/completions
		{"bare domain", "https://api.openai.com", "https://api.openai.com/v1/chat/completions"placeholder,
		{"domain with trailing slash", "https://api.openai.com/", "https://api.openai.com/v1/chat/completions"placeholder,
		// 第三方上游常见形式
		{"third-party bare domain", "https://api.deepseek.com", "https://api.deepseek.com/v1/chat/completions"placeholder,
		{"third-party with path prefix", "https://api.gptgod.online/api", "https://api.gptgod.online/api/v1/chat/completions"placeholder,
		{"third-party versioned path", "https://open.bigmodel.cn/api/paas/v4", "https://open.bigmodel.cn/api/paas/v4/chat/completions"placeholder,
		// 带空白字符
		{"whitespace trimmed", "  https://api.openai.com/v1  ", "https://api.openai.com/v1/chat/completions"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := buildOpenAIChatCompletionsURL(tt.base)
			require.Equal(t, tt.want, got)
	placeholder)
placeholder
placeholder

// TestBuildOpenAIResponsesURL_ProbeURL 锁定 probe/测试端点使用的 URL 构建逻辑，
// 确保 buildOpenAIResponsesURL 对标准 OpenAI base_url 格式均拼出 `/v1/responses`。
func TestBuildOpenAIResponsesURL_ProbeURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		base string
		want string
placeholder{
		{"bare domain", "https://api.openai.com", "https://api.openai.com/v1/responses"placeholder,
		{"domain trailing slash", "https://api.openai.com/", "https://api.openai.com/v1/responses"placeholder,
		{"bare /v1", "https://api.openai.com/v1", "https://api.openai.com/v1/responses"placeholder,
		{"already /responses", "https://api.openai.com/v1/responses", "https://api.openai.com/v1/responses"placeholder,
		{"third-party bare domain", "https://api.deepseek.com", "https://api.deepseek.com/v1/responses"placeholder,
		{"third-party versioned path", "https://open.bigmodel.cn/api/paas/v4", "https://open.bigmodel.cn/api/paas/v4/responses"placeholder,
		{"only domain, no scheme", "api.gptgod.online", "api.gptgod.online/v1/responses"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := buildOpenAIResponsesURL(tt.base)
			require.Equal(t, tt.want, got)
	placeholder)
placeholder
placeholder

func TestForwardAsRawChatCompletions_ForcesStreamUsageUpstreamAndPassesUsageDownstream(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"model":"gpt-5.4","messages":[{"role":"user","content":"hello"placeholder],"stream":trueplaceholder`)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstreamBody := strings.Join([]string{
		`data: {"id":"chatcmpl_1","object":"chat.completion.chunk","model":"gpt-5.4","choices":[{"index":0,"delta":{"content":"ok"placeholderplaceholder]placeholder`,
		"",
		`data: {"id":"chatcmpl_1","object":"chat.completion.chunk","model":"gpt-5.4","choices":[],"usage":{"prompt_tokens":9,"completion_tokens":4,"total_tokens":13,"prompt_tokens_details":{"cached_tokens":3placeholderplaceholderplaceholder`,
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid_raw_usage"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamBody)),
placeholderplaceholder

	svc := &OpenAIGatewayService{
		cfg:          rawChatCompletionsTestConfig(),
		httpUpstream: upstream,
placeholder
	account := rawChatCompletionsTestAccount()

	result, err := svc.forwardAsRawChatCompletions(context.Background(), c, account, body, "")
placeholder
	require.NotNil(t, result)
	require.Equal(t, 9, result.Usage.InputTokens)
	require.Equal(t, 4, result.Usage.OutputTokens)
	require.Equal(t, 3, result.Usage.CacheReadInputTokens)
	require.NotNil(t, upstream.lastReq)
	require.NoError(t, upstream.lastReq.Context().Err())
	require.Equal(t, HTTPUpstreamProfileOpenAI, HTTPUpstreamProfileFromContext(upstream.lastReq.Context()))
	require.True(t, gjson.GetBytes(upstream.lastBody, "stream_options.include_usage").Bool())
	require.Contains(t, rec.Body.String(), `"usage"`)
	require.Contains(t, rec.Body.String(), "data: [DONE]")
placeholder

func TestForwardAsRawChatCompletions_PreservesMappedGPT56MaxEffort(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"model":"sol","messages":[{"role":"user","content":"hello"placeholder],"reasoning_effort":"max","stream":falseplaceholder`)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(
			`{"id":"chatcmpl_max","object":"chat.completion","model":"gpt-5.6-sol","choices":[{"index":0,"message":{"role":"assistant","content":"ok"placeholder,"finish_reason":"stop"placeholder],"usage":{"prompt_tokens":3,"completion_tokens":2,"total_tokens":5placeholderplaceholder`,
		)),
placeholderplaceholder
	svc := &OpenAIGatewayService{
		cfg:          rawChatCompletionsTestConfig(),
		httpUpstream: upstream,
placeholder
	account := rawChatCompletionsTestAccount()
	account.Credentials["model_mapping"] = map[string]any{"sol": "gpt-5.6-sol"placeholder

	result, err := svc.forwardAsRawChatCompletions(context.Background(), c, account, body, "")

placeholder
	require.NotNil(t, result)
	require.Equal(t, "gpt-5.6-sol", gjson.GetBytes(upstream.lastBody, "model").String())
	require.Equal(t, "max", gjson.GetBytes(upstream.lastBody, "reasoning_effort").String())
	require.NotNil(t, result.ReasoningEffort)
	require.Equal(t, "max", *result.ReasoningEffort)
placeholder

func TestForwardAsRawChatCompletions_NonStreamingCapturesCacheWriteUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		usageJSON string
		wantWrite int
placeholder{
		{
			name:      "positive cache write",
			usageJSON: `{"prompt_tokens":12,"completion_tokens":3,"total_tokens":15,"prompt_tokens_details":{"cached_tokens":4,"cache_write_tokens":6placeholderplaceholder`,
			wantWrite: 6,
	placeholder,
		{
			name:      "nested zero overrides legacy alias",
			usageJSON: `{"prompt_tokens":12,"completion_tokens":3,"total_tokens":15,"cache_creation_input_tokens":19,"prompt_tokens_details":{"cached_tokens":4,"cache_write_tokens":0placeholderplaceholder`,
			wantWrite: 0,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := []byte(`{"model":"gpt-5.6","messages":[{"role":"user","content":"hello"placeholder],"stream":falseplaceholder`)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			upstream := &httpUpstreamRecorder{resp: &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
				Body: io.NopCloser(strings.NewReader(
					`{"id":"chatcmpl_cache","object":"chat.completion","model":"gpt-5.6","choices":[{"index":0,"message":{"role":"assistant","content":"ok"placeholder,"finish_reason":"stop"placeholder],"usage":` + tt.usageJSON + `placeholder`,
				)),
		placeholderplaceholder
			svc := &OpenAIGatewayService{
				cfg:          rawChatCompletionsTestConfig(),
				httpUpstream: upstream,
		placeholder

			result, err := svc.forwardAsRawChatCompletions(context.Background(), c, rawChatCompletionsTestAccount(), body, "")

		placeholder
			require.NotNil(t, result)
			require.Equal(t, 12, result.Usage.InputTokens)
			require.Equal(t, 4, result.Usage.CacheReadInputTokens)
			require.Equal(t, tt.wantWrite, result.Usage.CacheCreationInputTokens)
	placeholder)
placeholder
placeholder

func TestForwardAsRawChatCompletions_PreservesDeepSeekReasoningContentNonStreaming(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"model":"deepseek-reasoner","messages":[{"role":"user","content":"hello"placeholder],"stream":falseplaceholder`)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstreamJSON := `{"id":"chatcmpl_reasoning","object":"chat.completion","model":"deepseek-reasoner","choices":[{"index":0,"message":{"role":"assistant","reasoning_content":"think first","content":"final answer"placeholder,"finish_reason":"stop"placeholder],"usage":{"prompt_tokens":3,"completion_tokens":5,"total_tokens":8placeholderplaceholder`
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholder, "x-request-id": []string{"rid_deepseek_reasoning_json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamJSON)),
placeholderplaceholder

	svc := &OpenAIGatewayService{
		cfg:          rawChatCompletionsTestConfig(),
		httpUpstream: upstream,
placeholder
	account := rawChatCompletionsTestAccount()

	result, err := svc.forwardAsRawChatCompletions(context.Background(), c, account, body, "")
placeholder
	require.NotNil(t, result)
	require.Equal(t, 3, result.Usage.InputTokens)
	require.Equal(t, 5, result.Usage.OutputTokens)
	require.Equal(t, "think first", gjson.Get(rec.Body.String(), "choices.0.message.reasoning_content").String())
	require.Equal(t, "final answer", gjson.Get(rec.Body.String(), "choices.0.message.content").String())
placeholder

func TestForwardAsRawChatCompletions_PreservesDeepSeekReasoningContentStreaming(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"model":"deepseek-reasoner","messages":[{"role":"user","content":"hello"placeholder],"stream":trueplaceholder`)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstreamBody := strings.Join([]string{
		`data: {"id":"chatcmpl_reasoning","object":"chat.completion.chunk","model":"deepseek-reasoner","choices":[{"index":0,"delta":{"role":"assistant"placeholder,"finish_reason":nullplaceholder]placeholder`,
		"",
		`data: {"id":"chatcmpl_reasoning","object":"chat.completion.chunk","model":"deepseek-reasoner","choices":[{"index":0,"delta":{"reasoning_content":"think first"placeholder,"finish_reason":nullplaceholder]placeholder`,
		"",
		`data: {"id":"chatcmpl_reasoning","object":"chat.completion.chunk","model":"deepseek-reasoner","choices":[{"index":0,"delta":{"content":"final answer"placeholder,"finish_reason":nullplaceholder]placeholder`,
		"",
		`data: {"id":"chatcmpl_reasoning","object":"chat.completion.chunk","model":"deepseek-reasoner","choices":[],"usage":{"prompt_tokens":3,"completion_tokens":5,"total_tokens":8placeholderplaceholder`,
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid_deepseek_reasoning_stream"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamBody)),
placeholderplaceholder

	svc := &OpenAIGatewayService{
		cfg:          rawChatCompletionsTestConfig(),
		httpUpstream: upstream,
placeholder
	account := rawChatCompletionsTestAccount()

	result, err := svc.forwardAsRawChatCompletions(context.Background(), c, account, body, "")
placeholder
	require.NotNil(t, result)
	require.Equal(t, 3, result.Usage.InputTokens)
	require.Equal(t, 5, result.Usage.OutputTokens)
	require.Contains(t, rec.Body.String(), `"reasoning_content":"think first"`)
	require.Contains(t, rec.Body.String(), `"content":"final answer"`)
	require.Contains(t, rec.Body.String(), "data: [DONE]")
placeholder

func TestForwardAsRawChatCompletions_PreservesDeepSeekReasoningContentInRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"model":"deepseek-v4-pro","messages":[{"role":"user","content":"weather"placeholder,{"role":"assistant","reasoning_content":"need tool","content":"","tool_calls":[{"id":"call_1","type":"function","function":{"name":"get_weather","arguments":"{placeholder"placeholderplaceholder]placeholder,{"role":"tool","tool_call_id":"call_1","content":"cloudy"placeholder],"stream":falseplaceholder`)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholder, "x-request-id": []string{"rid_deepseek_reasoning_request"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"id":"chatcmpl_request","object":"chat.completion","model":"deepseek-v4-pro","choices":[{"index":0,"message":{"role":"assistant","content":"done"placeholder,"finish_reason":"stop"placeholder],"usage":{"prompt_tokens":4,"completion_tokens":2,"total_tokens":6placeholderplaceholder`)),
placeholderplaceholder

	svc := &OpenAIGatewayService{
		cfg:          rawChatCompletionsTestConfig(),
		httpUpstream: upstream,
placeholder
	account := rawChatCompletionsTestAccount()

	result, err := svc.forwardAsRawChatCompletions(context.Background(), c, account, body, "")
placeholder
	require.NotNil(t, result)
	require.Equal(t, "need tool", gjson.GetBytes(upstream.lastBody, "messages.1.reasoning_content").String())
	require.Equal(t, "get_weather", gjson.GetBytes(upstream.lastBody, "messages.1.tool_calls.0.function.name").String())
placeholder

func TestForwardAsRawChatCompletions_NormalizesGLMReasoningEffortForUpstream(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"model":"glm-5.2","messages":[{"role":"user","content":"hello"placeholder],"reasoning_effort":"xhigh","stream":falseplaceholder`)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholder, "x-request-id": []string{"rid_glm_effort"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"id":"chatcmpl_glm","object":"chat.completion","model":"glm-5.2","choices":[{"index":0,"message":{"role":"assistant","content":"ok"placeholder,"finish_reason":"stop"placeholder],"usage":{"prompt_tokens":1,"completion_tokens":2,"total_tokens":3placeholderplaceholder`)),
placeholderplaceholder

	svc := &OpenAIGatewayService{
		cfg:          rawChatCompletionsTestConfig(),
		httpUpstream: upstream,
placeholder
	account := rawChatCompletionsTestAccount()

	result, err := svc.forwardAsRawChatCompletions(context.Background(), c, account, body, "")
placeholder
	require.NotNil(t, result)
	require.Equal(t, "max", gjson.GetBytes(upstream.lastBody, "reasoning_effort").String())
placeholder

func TestForwardAsRawChatCompletions_SilentRefusalTriggersFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := largeRawChatCompletionsBody()
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstreamBody := strings.Join([]string{
		`data: {"id":"chatcmpl_silent","object":"chat.completion.chunk","model":"gpt-5.5","choices":[{"index":0,"delta":{"role":"assistant"placeholderplaceholder]placeholder`,
		"",
		`data: {"id":"chatcmpl_silent","object":"chat.completion.chunk","model":"gpt-5.5","choices":[{"index":0,"delta":{"content":""placeholder,"finish_reason":"stop"placeholder]placeholder`,
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid_silent"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamBody)),
placeholderplaceholder

	svc := &OpenAIGatewayService{
		cfg:          rawChatCompletionsTestConfig(),
		httpUpstream: upstream,
placeholder

	result, err := svc.forwardAsRawChatCompletions(context.Background(), c, rawChatCompletionsTestAccount(), body, "")
	require.Nil(t, result)
	var failoverErr *UpstreamFailoverError
	require.True(t, errors.As(err, &failoverErr))
	require.Equal(t, http.StatusBadGateway, failoverErr.StatusCode)
	require.True(t, IsOpenAISilentRefusalErrorBody(failoverErr.ResponseBody))
	require.False(t, c.Writer.Written(), "silent refusal must not commit a 200 response before failover")
	require.Empty(t, rec.Body.String())
placeholder

func TestForwardAsRawChatCompletions_SilentRefusalToolCallsExempt(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := largeRawChatCompletionsBody()
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstreamBody := strings.Join([]string{
		`data: {"id":"chatcmpl_tool","object":"chat.completion.chunk","model":"gpt-5.5","choices":[{"index":0,"delta":{"role":"assistant"placeholderplaceholder]placeholder`,
		"",
		`data: {"id":"chatcmpl_tool","object":"chat.completion.chunk","model":"gpt-5.5","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_1","type":"function","function":{"name":"lookup","arguments":""placeholderplaceholder]placeholderplaceholder]placeholder`,
		"",
		`data: {"id":"chatcmpl_tool","object":"chat.completion.chunk","model":"gpt-5.5","choices":[{"index":0,"delta":{placeholder,"finish_reason":"tool_calls"placeholder]placeholder`,
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid_tool"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamBody)),
placeholderplaceholder

	svc := &OpenAIGatewayService{
		cfg:          rawChatCompletionsTestConfig(),
		httpUpstream: upstream,
placeholder

	result, err := svc.forwardAsRawChatCompletions(context.Background(), c, rawChatCompletionsTestAccount(), body, "")
placeholder
	require.NotNil(t, result)
	require.Contains(t, rec.Body.String(), `"tool_calls"`)
	require.Contains(t, rec.Body.String(), `"finish_reason":"tool_calls"`)
placeholder

func TestHandleChatStreamingResponse_SilentRefusalReasoningSummaryExempt(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)

	upstreamBody := strings.Join([]string{
		`data: {"type":"response.created","response":{"id":"resp_reasoning","model":"gpt-5.5"placeholderplaceholder`,
		"",
		`data: {"type":"response.reasoning_summary_text.delta","delta":"thinking only"placeholder`,
		"",
		`data: {"type":"response.completed","response":{"id":"resp_reasoning","model":"gpt-5.5","status":"completed"placeholderplaceholder`,
		"",
placeholder, "\n")
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid_reasoning"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamBody)),
placeholder
	svc := &OpenAIGatewayService{cfg: rawChatCompletionsTestConfig()placeholder

	result, err := svc.handleChatStreamingResponse(
		resp,
		c,
		rawChatCompletionsTestAccount(),
		"gpt-5.5",
		"gpt-5.5",
		"gpt-5.5",
		time.Now(),
		openAISilentRefusalMinRequestBodyBytes,
	)
placeholder
	require.NotNil(t, result)
	require.Contains(t, rec.Body.String(), `"reasoning_content":"thinking only"`)
	require.Contains(t, rec.Body.String(), "data: [DONE]")
placeholder

func TestForwardAsRawChatCompletions_SilentRefusalNormalContentExempt(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := largeRawChatCompletionsBody()
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstreamBody := strings.Join([]string{
		`data: {"id":"chatcmpl_ok","object":"chat.completion.chunk","model":"gpt-5.5","choices":[{"index":0,"delta":{"role":"assistant"placeholderplaceholder]placeholder`,
		"",
		`data: {"id":"chatcmpl_ok","object":"chat.completion.chunk","model":"gpt-5.5","choices":[{"index":0,"delta":{"content":"ok"placeholderplaceholder]placeholder`,
		"",
		`data: {"id":"chatcmpl_ok","object":"chat.completion.chunk","model":"gpt-5.5","choices":[{"index":0,"delta":{"content":""placeholder,"finish_reason":"stop"placeholder]placeholder`,
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid_ok"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamBody)),
placeholderplaceholder

	svc := &OpenAIGatewayService{
		cfg:          rawChatCompletionsTestConfig(),
		httpUpstream: upstream,
placeholder

	result, err := svc.forwardAsRawChatCompletions(context.Background(), c, rawChatCompletionsTestAccount(), body, "")
placeholder
	require.NotNil(t, result)
	require.Contains(t, rec.Body.String(), `"content":"ok"`)
	require.Contains(t, rec.Body.String(), "data: [DONE]")
placeholder

func TestForwardAsRawChatCompletions_ClientDisconnectDrainsUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"model":"gpt-5.4","messages":[{"role":"user","content":"hello"placeholder],"stream":trueplaceholder`)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Writer = &openAIChatFailingWriter{ResponseWriter: c.Writer, failAfter: 0placeholder
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstreamBody := strings.Join([]string{
		`data: {"id":"chatcmpl_1","object":"chat.completion.chunk","model":"gpt-5.4","choices":[{"index":0,"delta":{"content":"ok"placeholderplaceholder]placeholder`,
		"",
		`data: {"id":"chatcmpl_1","object":"chat.completion.chunk","model":"gpt-5.4","choices":[],"usage":{"prompt_tokens":17,"completion_tokens":8,"total_tokens":25,"prompt_tokens_details":{"cached_tokens":6placeholderplaceholderplaceholder`,
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid_raw_disconnect"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamBody)),
placeholderplaceholder

	svc := &OpenAIGatewayService{
		cfg:          rawChatCompletionsTestConfig(),
		httpUpstream: upstream,
placeholder
	account := rawChatCompletionsTestAccount()

	result, err := svc.forwardAsRawChatCompletions(context.Background(), c, account, body, "")
placeholder
	require.NotNil(t, result)
	require.Equal(t, 17, result.Usage.InputTokens)
	require.Equal(t, 8, result.Usage.OutputTokens)
	require.Equal(t, 6, result.Usage.CacheReadInputTokens)
	require.True(t, gjson.GetBytes(upstream.lastBody, "stream_options.include_usage").Bool())
placeholder

func TestForwardAsRawChatCompletions_UpstreamRequestIgnoresClientCancel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	reqCtx, cancel := context.WithCancel(context.Background())
	body := []byte(`{"model":"gpt-5.4","messages":[{"role":"user","content":"hello"placeholder],"stream":trueplaceholder`)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body)).WithContext(reqCtx)
	c.Request.Header.Set("Content-Type", "application/json")
	cancel()

	upstreamBody := strings.Join([]string{
		`data: {"id":"chatcmpl_1","object":"chat.completion.chunk","model":"gpt-5.4","choices":[],"usage":{"prompt_tokens":5,"completion_tokens":2,"total_tokens":7placeholderplaceholder`,
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid_raw_ctx"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamBody)),
placeholderplaceholder

	svc := &OpenAIGatewayService{
		cfg:          rawChatCompletionsTestConfig(),
		httpUpstream: upstream,
placeholder
	account := rawChatCompletionsTestAccount()

	result, err := svc.forwardAsRawChatCompletions(reqCtx, c, account, body, "")
placeholder
	require.NotNil(t, result)
	require.NotNil(t, upstream.lastReq)
	require.NoError(t, upstream.lastReq.Context().Err())
placeholder

func TestForwardAsChatCompletions_UnknownResponsesSupportFallbackUsesVersionedChatURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"model":"glm-4.5-air","messages":[{"role":"user","content":"hello"placeholder],"stream":falseplaceholder`)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstream := &httpUpstreamRecorder{responses: []*http.Response{
		{
			StatusCode: http.StatusNotFound,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"not found"placeholderplaceholder`)),
	placeholder,
		{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholder, "x-request-id": []string{"rid_raw_fallback"placeholderplaceholder,
			Body: io.NopCloser(strings.NewReader(
				`{"id":"chatcmpl_1","object":"chat.completion","model":"glm-4.5-air","choices":[{"index":0,"message":{"role":"assistant","content":"ok"placeholder,"finish_reason":"stop"placeholder],"usage":{"prompt_tokens":1,"completion_tokens":2,"total_tokens":3placeholderplaceholder`,
			)),
	placeholder,
placeholderplaceholder

	svc := &OpenAIGatewayService{
		cfg:          rawChatCompletionsTestConfig(),
		httpUpstream: upstream,
placeholder
	account := rawChatCompletionsTestAccount()
	account.Credentials["base_url"] = "https://open.bigmodel.cn/api/paas/v4"

	result, err := svc.ForwardAsChatCompletions(context.Background(), c, account, body, "", "")
placeholder
	require.NotNil(t, result)
	require.Equal(t, 1, result.Usage.InputTokens)
	require.Equal(t, 2, result.Usage.OutputTokens)
	require.Len(t, upstream.requests, 2)
	require.Equal(t, "https://open.bigmodel.cn/api/paas/v4/responses", upstream.requests[0].URL.String())
	require.Equal(t, "https://open.bigmodel.cn/api/paas/v4/chat/completions", upstream.requests[1].URL.String())
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"content":"ok"`)
placeholder

func TestIsOpenAIChatUsageOnlyStreamChunk(t *testing.T) {
	t.Parallel()

	require.True(t, isOpenAIChatUsageOnlyStreamChunk(`{"choices":[],"usage":{"prompt_tokens":1,"completion_tokens":2placeholderplaceholder`))
	require.False(t, isOpenAIChatUsageOnlyStreamChunk(`{"choices":[{"index":0placeholder],"usage":{"prompt_tokens":1,"completion_tokens":2placeholderplaceholder`))
	require.False(t, isOpenAIChatUsageOnlyStreamChunk(`{"choices":[]placeholder`))
	require.False(t, isOpenAIChatUsageOnlyStreamChunk(``))
placeholder

func TestEnsureOpenAIChatStreamUsage(t *testing.T) {
	t.Parallel()

	body, err := ensureOpenAIChatStreamUsage([]byte(`{"model":"gpt-5.4"placeholder`))
placeholder
	require.True(t, gjson.GetBytes(body, "stream_options.include_usage").Bool())

	body, err = ensureOpenAIChatStreamUsage([]byte(`{"model":"gpt-5.4","stream_options":{"include_usage":falseplaceholderplaceholder`))
placeholder
	require.True(t, gjson.GetBytes(body, "stream_options.include_usage").Bool())
placeholder

func TestBufferRawChatCompletions_RejectsOversizedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader("toolong")),
placeholder
	svc := &OpenAIGatewayService{cfg: rawChatCompletionsTestConfig()placeholder
	svc.cfg.Gateway.UpstreamResponseReadMaxBytes = 3

	result, err := svc.bufferRawChatCompletions(c, resp, rawChatCompletionsTestAccount(), "gpt-5.4", "gpt-5.4", "gpt-5.4", nil, nil, time.Now())
	require.ErrorIs(t, err, ErrUpstreamResponseBodyTooLarge)
	require.Nil(t, result)
	require.Equal(t, http.StatusBadGateway, rec.Code)
placeholder

func rawChatCompletionsTestConfig() *config.Config {
	return &config.Config{
		Security: config.SecurityConfig{
			URLAllowlist: config.URLAllowlistConfig{
				Enabled:           false,
				AllowInsecureHTTP: true,
		placeholder,
	placeholder,
placeholder
placeholder

func rawChatCompletionsTestAccount() *Account {
placeholder
		ID:          101,
		Name:        "raw-openai-apikey",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "http://upstream.example",
	placeholder,
placeholder
placeholder

func largeRawChatCompletionsBody() []byte {
	return []byte(`{"model":"gpt-5.5","messages":[{"role":"user","content":"` +
		strings.Repeat("x", openAISilentRefusalMinRequestBodyBytes) +
		`"placeholder],"stream":trueplaceholder`)
placeholder
