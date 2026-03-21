package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type geminiCompatHTTPUpstreamStub struct {
	response *http.Response
	err      error
	calls    int
	lastReq  *http.Request
placeholder

func (s *geminiCompatHTTPUpstreamStub) Do(req *http.Request, proxyURL string, accountID int64, accountConcurrency int) (*http.Response, error) {
	s.calls++
	s.lastReq = req
	if s.err != nil {
		return nil, s.err
placeholder
	if s.response == nil {
		return nil, fmt.Errorf("missing stub response")
placeholder
	resp := *s.response
	return &resp, nil
placeholder

func (s *geminiCompatHTTPUpstreamStub) DoWithTLS(req *http.Request, proxyURL string, accountID int64, accountConcurrency int, enableTLSFingerprint bool) (*http.Response, error) {
	return s.Do(req, proxyURL, accountID, accountConcurrency)
placeholder

// TestConvertClaudeToolsToGeminiTools_CustomType 测试custom类型工具转换
func TestConvertClaudeToolsToGeminiTools_CustomType(t *testing.T) {
	tests := []struct {
		name        string
		tools       any
		expectedLen int
		description string
placeholder{
		{
			name: "Standard tools",
			tools: []any{
				map[string]any{
					"name":         "get_weather",
					"description":  "Get weather info",
					"input_schema": map[string]any{"type": "object"placeholder,
			placeholder,
		placeholder,
			expectedLen: 1,
			description: "标准工具格式应该正常转换",
	placeholder,
		{
			name: "Custom type tool (MCP format)",
			tools: []any{
				map[string]any{
					"type": "custom",
					"name": "mcp_tool",
					"custom": map[string]any{
						"description":  "MCP tool description",
						"input_schema": map[string]any{"type": "object"placeholder,
				placeholder,
			placeholder,
		placeholder,
			expectedLen: 1,
			description: "Custom类型工具应该从custom字段读取",
	placeholder,
		{
			name: "Mixed standard and custom tools",
			tools: []any{
				map[string]any{
					"name":         "standard_tool",
					"description":  "Standard",
					"input_schema": map[string]any{"type": "object"placeholder,
			placeholder,
				map[string]any{
					"type": "custom",
					"name": "custom_tool",
					"custom": map[string]any{
						"description":  "Custom",
						"input_schema": map[string]any{"type": "object"placeholder,
				placeholder,
			placeholder,
		placeholder,
			expectedLen: 1,
			description: "混合工具应该都能正确转换",
	placeholder,
		{
			name: "Custom tool without custom field",
			tools: []any{
				map[string]any{
					"type": "custom",
					"name": "invalid_custom",
					// 缺少 custom 字段
			placeholder,
		placeholder,
			expectedLen: 0, // 应该被跳过
			description: "缺少custom字段的custom工具应该被跳过",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertClaudeToolsToGeminiTools(tt.tools)

			if tt.expectedLen == 0 {
				if result != nil {
					t.Errorf("%s: expected nil result, got %v", tt.description, result)
			placeholder
				return
		placeholder

			if result == nil {
				t.Fatalf("%s: expected non-nil result", tt.description)
		placeholder

			if len(result) != 1 {
				t.Errorf("%s: expected 1 tool declaration, got %d", tt.description, len(result))
				return
		placeholder

			toolDecl, ok := result[0].(map[string]any)
			if !ok {
				t.Fatalf("%s: result[0] is not map[string]any", tt.description)
		placeholder

			funcDecls, ok := toolDecl["functionDeclarations"].([]any)
			if !ok {
				t.Fatalf("%s: functionDeclarations is not []any", tt.description)
		placeholder

			toolsArr, _ := tt.tools.([]any)
			expectedFuncCount := 0
			for _, tool := range toolsArr {
				toolMap, _ := tool.(map[string]any)
				if toolMap["name"] != "" {
					// 检查是否为有效的custom工具
					if toolMap["type"] == "custom" {
						if toolMap["custom"] != nil {
							expectedFuncCount++
					placeholder
				placeholder else {
						expectedFuncCount++
				placeholder
			placeholder
		placeholder

			if len(funcDecls) != expectedFuncCount {
				t.Errorf("%s: expected %d function declarations, got %d",
					tt.description, expectedFuncCount, len(funcDecls))
		placeholder
	placeholder)
placeholder
placeholder

func TestGeminiHandleNativeNonStreamingResponse_DebugDisabledDoesNotEmitHeaderLogs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logSink, restore := captureStructuredLog(t)
	defer restore()

	svc := &GeminiMessagesCompatService{
		cfg: &config.Config{
			Gateway: config.GatewayConfig{
				GeminiDebugResponseHeaders: false,
		placeholder,
	placeholder,
placeholder

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type":      []string{"application/json"placeholder,
			"X-RateLimit-Limit": []string{"60"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(`{"usageMetadata":{"promptTokenCount":10,"candidatesTokenCount":2placeholderplaceholder`)),
placeholder

	usage, err := svc.handleNativeNonStreamingResponse(c, resp, false)
placeholder
	require.NotNil(t, usage)
	require.False(t, logSink.ContainsMessage("[GeminiAPI]"), "debug 关闭时不应输出 Gemini 响应头日志")
placeholder

func TestGeminiMessagesCompatServiceForward_PreservesRequestedModelAndMappedUpstreamModel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	httpStub := &geminiCompatHTTPUpstreamStub{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"x-request-id": []string{"gemini-req-1"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"candidates":[{"content":{"parts":[{"text":"hello"placeholder]placeholderplaceholder],"usageMetadata":{"promptTokenCount":10,"candidatesTokenCount":5placeholderplaceholder`)),
	placeholder,
placeholder
	svc := &GeminiMessagesCompatService{httpUpstream: httpStub, cfg: &config.Config{placeholderplaceholder
	account := &Account{
		ID:   1,
		Type: AccountTypeAPIKey,
placeholder
			"api_key": "test-key",
			"model_mapping": map[string]any{
				"claude-sonnet-4": "claude-sonnet-4-20250514",
		placeholder,
	placeholder,
placeholder
	body := []byte(`{"model":"claude-sonnet-4","max_tokens":16,"messages":[{"role":"user","content":"hello"placeholder]placeholder`)

	result, err := svc.Forward(context.Background(), c, account, body)
placeholder
	require.NotNil(t, result)
	require.Equal(t, "claude-sonnet-4", result.Model)
	require.Equal(t, "claude-sonnet-4-20250514", result.UpstreamModel)
	require.Equal(t, 1, httpStub.calls)
	require.NotNil(t, httpStub.lastReq)
	require.Contains(t, httpStub.lastReq.URL.String(), "/models/claude-sonnet-4-20250514:")
placeholder

func TestConvertClaudeMessagesToGeminiGenerateContent_AddsThoughtSignatureForToolUse(t *testing.T) {
	claudeReq := map[string]any{
		"model":      "placeholder",
		"max_tokens": 10,
		"messages": []any{
			map[string]any{
				"role": "user",
				"content": []any{
					map[string]any{"type": "text", "text": "hi"placeholder,
			placeholder,
		placeholder,
			map[string]any{
				"role": "assistant",
				"content": []any{
					map[string]any{"type": "text", "text": "ok"placeholder,
					map[string]any{
						"type":  "tool_use",
						"id":    "toolu_123",
						"name":  "default_api:write_file",
						"input": map[string]any{"path": "a.txt", "content": "x"placeholder,
						// no signature on purpose
				placeholder,
			placeholder,
		placeholder,
	placeholder,
		"tools": []any{
			map[string]any{
				"name":        "default_api:write_file",
				"description": "write file",
				"input_schema": map[string]any{
					"type":       "object",
					"properties": map[string]any{"path": map[string]any{"type": "string"placeholderplaceholder,
			placeholder,
		placeholder,
	placeholder,
placeholder
	b, _ := json.Marshal(claudeReq)

	out, err := convertClaudeMessagesToGeminiGenerateContent(b)
	if err != nil {
		t.Fatalf("convert failed: %v", err)
placeholder
	s := string(out)
	if !strings.Contains(s, "\"functionCall\"") {
		t.Fatalf("expected functionCall in output, got: %s", s)
placeholder
	if !strings.Contains(s, "\"thoughtSignature\":\""+geminiDummyThoughtSignature+"\"") {
		t.Fatalf("expected injected thoughtSignature %q, got: %s", geminiDummyThoughtSignature, s)
placeholder
placeholder

func TestEnsureGeminiFunctionCallThoughtSignatures_InsertsWhenMissing(t *testing.T) {
	geminiReq := map[string]any{
		"contents": []any{
			map[string]any{
				"role": "user",
				"parts": []any{
					map[string]any{
						"functionCall": map[string]any{
							"name": "default_api:write_file",
							"args": map[string]any{"path": "a.txt"placeholder,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
	placeholder,
placeholder
	b, _ := json.Marshal(geminiReq)
	out := ensureGeminiFunctionCallThoughtSignatures(b)
	s := string(out)
	if !strings.Contains(s, "\"thoughtSignature\":\""+geminiDummyThoughtSignature+"\"") {
		t.Fatalf("expected injected thoughtSignature %q, got: %s", geminiDummyThoughtSignature, s)
placeholder
placeholder

// TestUnwrapGeminiResponse 测试 unwrapGeminiResponse 的各种输入场景
// 关键区别：只有 response 为 JSON 对象/数组时才解包
func TestUnwrapGeminiResponse(t *testing.T) {
	// 构造 >50KB 的大型 JSON 对象
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
			name:     "正常 response 包装（JSON 对象）",
			input:    []byte(`{"response":{"key":"val"placeholderplaceholder`),
			expected: `{"key":"val"placeholder`,
	placeholder,
		{
			name:     "无包装直接返回",
			input:    []byte(`{"key":"val"placeholder`),
			expected: `{"key":"val"placeholder`,
	placeholder,
		{
			name:     "空 JSON",
			input:    []byte(`{placeholder`),
			expected: `{placeholder`,
	placeholder,
		{
			name:     "null response 返回原始 body",
			input:    []byte(`{"response":nullplaceholder`),
			expected: `{"response":nullplaceholder`,
	placeholder,
		{
			name:     "非法 JSON 返回原始 body",
			input:    []byte(`not json`),
			expected: `not json`,
	placeholder,
		{
			name:     "response 为基础类型 string 返回原始 body",
			input:    []byte(`{"response":"hello"placeholder`),
			expected: `{"response":"hello"placeholder`,
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
			got, err := unwrapGeminiResponse(tt.input)
			if tt.wantErr {
			placeholder
				return
		placeholder
		placeholder
			require.Equal(t, tt.expected, strings.TrimSpace(string(got)))
	placeholder)
placeholder
placeholder

// ---------------------------------------------------------------------------
// Task 8.1 — extractGeminiUsage 测试
// ---------------------------------------------------------------------------

func TestExtractGeminiUsage(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantNil   bool
		wantUsage *ClaudeUsage
placeholder{
		{
			name:    "完整 usageMetadata",
			input:   `{"usageMetadata":{"promptTokenCount":100,"candidatesTokenCount":50,"cachedContentTokenCount":20placeholderplaceholder`,
			wantNil: false,
			wantUsage: &ClaudeUsage{
				InputTokens:          80,
				OutputTokens:         50,
				CacheReadInputTokens: 20,
		placeholder,
	placeholder,
		{
			name:    "包含 thoughtsTokenCount",
			input:   `{"usageMetadata":{"promptTokenCount":100,"candidatesTokenCount":20,"thoughtsTokenCount":50placeholderplaceholder`,
			wantNil: false,
			wantUsage: &ClaudeUsage{
				InputTokens:          100,
				OutputTokens:         70,
				CacheReadInputTokens: 0,
		placeholder,
	placeholder,
		{
			name:    "包含 thoughtsTokenCount 与缓存",
			input:   `{"usageMetadata":{"promptTokenCount":100,"candidatesTokenCount":20,"cachedContentTokenCount":30,"thoughtsTokenCount":50placeholderplaceholder`,
			wantNil: false,
			wantUsage: &ClaudeUsage{
				InputTokens:          70,
				OutputTokens:         70,
				CacheReadInputTokens: 30,
		placeholder,
	placeholder,
		{
			name:    "缺失 cachedContentTokenCount",
			input:   `{"usageMetadata":{"promptTokenCount":100,"candidatesTokenCount":50placeholderplaceholder`,
			wantNil: false,
			wantUsage: &ClaudeUsage{
				InputTokens:          100,
				OutputTokens:         50,
				CacheReadInputTokens: 0,
		placeholder,
	placeholder,
		{
			name:    "无 usageMetadata",
			input:   `{"candidates":[]placeholder`,
			wantNil: true,
	placeholder,
		{
			// gjson 对 null 返回 Exists()=true，因此函数不会返回 nil，
			// 而是返回全零的 ClaudeUsage。
			name:    "null usageMetadata — gjson Exists 为 true",
			input:   `{"usageMetadata":nullplaceholder`,
			wantNil: false,
			wantUsage: &ClaudeUsage{
				InputTokens:          0,
				OutputTokens:         0,
				CacheReadInputTokens: 0,
		placeholder,
	placeholder,
		{
			name:    "零值字段",
			input:   `{"usageMetadata":{"promptTokenCount":0,"candidatesTokenCount":0,"cachedContentTokenCount":0placeholderplaceholder`,
			wantNil: false,
			wantUsage: &ClaudeUsage{
				InputTokens:          0,
				OutputTokens:         0,
				CacheReadInputTokens: 0,
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractGeminiUsage([]byte(tt.input))
			if tt.wantNil {
				if got != nil {
					t.Fatalf("期望返回 nil，实际返回 %+v", got)
			placeholder
				return
		placeholder
			if got == nil {
				t.Fatalf("期望返回非 nil，实际返回 nil")
		placeholder
			if got.InputTokens != tt.wantUsage.InputTokens {
				t.Errorf("InputTokens: 期望 %d，实际 %d", tt.wantUsage.InputTokens, got.InputTokens)
		placeholder
			if got.OutputTokens != tt.wantUsage.OutputTokens {
				t.Errorf("OutputTokens: 期望 %d，实际 %d", tt.wantUsage.OutputTokens, got.OutputTokens)
		placeholder
			if got.CacheReadInputTokens != tt.wantUsage.CacheReadInputTokens {
				t.Errorf("CacheReadInputTokens: 期望 %d，实际 %d", tt.wantUsage.CacheReadInputTokens, got.CacheReadInputTokens)
		placeholder
	placeholder)
placeholder
placeholder

// ---------------------------------------------------------------------------
// Task 8.2 — estimateGeminiCountTokens 测试
// ---------------------------------------------------------------------------

func TestEstimateGeminiCountTokens(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantGt0   bool // 期望结果 > 0
		wantExact *int // 如果非 nil，期望精确匹配
placeholder{
		{
			name: "含 systemInstruction 和 contents",
			input: `{
				"systemInstruction":{"parts":[{"text":"You are a helpful assistant."placeholder]placeholder,
				"contents":[{"parts":[{"text":"Hello, how are you?"placeholder]placeholder]
		placeholder`,
			wantGt0: true,
	placeholder,
		{
			name: "仅 contents，无 systemInstruction",
			input: `{
				"contents":[{"parts":[{"text":"Hello, how are you?"placeholder]placeholder]
		placeholder`,
			wantGt0: true,
	placeholder,
		{
			name:      "空 parts",
			input:     `{"contents":[{"parts":[]placeholder]placeholder`,
			wantGt0:   false,
			wantExact: intPtr(0),
	placeholder,
		{
			name:      "非文本 parts（inlineData）",
			input:     `{"contents":[{"parts":[{"inlineData":{"mimeType":"image/png"placeholderplaceholder]placeholder]placeholder`,
			wantGt0:   false,
			wantExact: intPtr(0),
	placeholder,
		{
			name:      "空白文本",
			input:     `{"contents":[{"parts":[{"text":"   "placeholder]placeholder]placeholder`,
			wantGt0:   false,
			wantExact: intPtr(0),
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := estimateGeminiCountTokens([]byte(tt.input))
			if tt.wantExact != nil {
				if got != *tt.wantExact {
					t.Errorf("期望精确值 %d，实际 %d", *tt.wantExact, got)
			placeholder
				return
		placeholder
			if tt.wantGt0 && got <= 0 {
				t.Errorf("期望返回 > 0，实际 %d", got)
		placeholder
			if !tt.wantGt0 && got != 0 {
				t.Errorf("期望返回 0，实际 %d", got)
		placeholder
	placeholder)
placeholder
placeholder

// ---------------------------------------------------------------------------
// Task 8.3 — ParseGeminiRateLimitResetTime 测试
// ---------------------------------------------------------------------------

func TestParseGeminiRateLimitResetTime(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantNil     bool
		approxDelta int64 // 预期的 (返回值 - now) 大约是多少秒
placeholder{
		{
			name:        "正常 quotaResetDelay",
			input:       `{"error":{"details":[{"metadata":{"quotaResetDelay":"12.345s"placeholderplaceholder]placeholderplaceholder`,
			wantNil:     false,
			approxDelta: 13, // 向上取整 12.345 -> 13
	placeholder,
		{
			name:        "daily quota",
			input:       `{"error":{"message":"quota per day exceeded"placeholderplaceholder`,
			wantNil:     false,
			approxDelta: -1, // 不检查精确 delta，仅检查非 nil
	placeholder,
		{
			name:    "无 details 且无 regex 匹配",
			input:   `{"error":{"message":"rate limit"placeholderplaceholder`,
			wantNil: true,
	placeholder,
		{
			name:        "regex 回退匹配",
			input:       `Please retry in 30s`,
			wantNil:     false,
			approxDelta: 30,
	placeholder,
		{
			name:    "完全无匹配",
			input:   `{"error":{"code":429placeholderplaceholder`,
			wantNil: true,
	placeholder,
		{
			name:        "非法 JSON 但 regex 回退仍工作",
			input:       `not json but Please retry in 10s`,
			wantNil:     false,
			approxDelta: 10,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now().Unix()
			got := ParseGeminiRateLimitResetTime([]byte(tt.input))

			if tt.wantNil {
				if got != nil {
					t.Fatalf("期望返回 nil，实际返回 %d", *got)
			placeholder
				return
		placeholder

			if got == nil {
				t.Fatalf("期望返回非 nil，实际返回 nil")
		placeholder

			// approxDelta == -1 表示只检查非 nil，不检查具体值（如 daily quota 场景）
			if tt.approxDelta == -1 {
				// 仅验证返回的时间戳在合理范围内（未来的某个时间）
				if *got < now {
					t.Errorf("期望返回的时间戳 >= now(%d)，实际 %d", now, *got)
			placeholder
				return
		placeholder

			// 使用 +/-2 秒容差进行范围检查
			delta := *got - now
			if delta < tt.approxDelta-2 || delta > tt.approxDelta+2 {
				t.Errorf("期望 delta 约为 %d 秒（+/-2），实际 delta = %d 秒（返回值=%d, now=%d）",
					tt.approxDelta, delta, *got, now)
		placeholder
	placeholder)
placeholder
placeholder
