package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/antigravity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

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
		"model": "claude-opus-4-5",
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

	result, err := svc.Forward(context.Background(), c, account, body)
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

func TestAntigravityMaxRetriesForModel_AfterSwitch(t *testing.T) {
	t.Setenv(antigravityMaxRetriesEnv, "4")
	t.Setenv(antigravityMaxRetriesAfterSwitchEnv, "7")
	t.Setenv(antigravityMaxRetriesClaudeEnv, "")
	t.Setenv(antigravityMaxRetriesGeminiTextEnv, "")
	t.Setenv(antigravityMaxRetriesGeminiImageEnv, "")

	got := antigravityMaxRetriesForModel("claude-sonnet-4-5", false)
	require.Equal(t, 4, got)

	got = antigravityMaxRetriesForModel("claude-sonnet-4-5", true)
	require.Equal(t, 7, got)
placeholder

func TestAntigravityMaxRetriesForModel_AfterSwitchFallback(t *testing.T) {
	t.Setenv(antigravityMaxRetriesEnv, "5")
	t.Setenv(antigravityMaxRetriesAfterSwitchEnv, "")
	t.Setenv(antigravityMaxRetriesClaudeEnv, "")
	t.Setenv(antigravityMaxRetriesGeminiTextEnv, "")
	t.Setenv(antigravityMaxRetriesGeminiImageEnv, "")

	got := antigravityMaxRetriesForModel("gemini-2.5-flash", true)
	require.Equal(t, 5, got)
placeholder

func TestAntigravityStreamUpstreamResponse_UsageAndFirstToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Header: http.Header{placeholder, Body: prplaceholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		_, _ = pw.Write([]byte("data: {\"usage\":{\"input_tokens\":1,\"output_tokens\":2,\"cache_read_input_tokens\":3,\"cache_creation_input_tokens\":4placeholderplaceholder\n"))
		_, _ = pw.Write([]byte("data: {\"usage\":{\"output_tokens\":5placeholderplaceholder\n"))
placeholder()

	svc := &AntigravityGatewayService{placeholder
	start := time.Now().Add(-10 * time.Millisecond)
	usage, firstTokenMs := svc.streamUpstreamResponse(c, resp, start)
	_ = pr.Close()

	require.NotNil(t, usage)
	require.Equal(t, 1, usage.InputTokens)
	// 第二次事件覆盖 output_tokens
	require.Equal(t, 5, usage.OutputTokens)
	require.Equal(t, 3, usage.CacheReadInputTokens)
	require.Equal(t, 4, usage.CacheCreationInputTokens)

	if firstTokenMs == nil {
		t.Fatalf("expected firstTokenMs to be set")
placeholder
	// 确保有透传输出
	require.True(t, strings.Contains(writer.Body.String(), "data:"))
placeholder
