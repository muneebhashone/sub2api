//go:build unit

package service

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	"github.com/stretchr/testify/require"
)

func newTestValidator() *ClaudeCodeValidator {
	return NewClaudeCodeValidator()
placeholder

// validClaudeCodeBody 构造一个完整有效的 Claude Code 请求体
func validClaudeCodeBody() map[string]any {
	return map[string]any{
		"model": "claude-sonnet-4-20250514",
		"system": []any{
			map[string]any{
				"type": "text",
				"text": "You are Claude Code, Anthropic's official CLI for Claude.",
		placeholder,
	placeholder,
		"metadata": map[string]any{
			"user_id": "user_" + "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2" + "_account__session_" + "12345678-1234-1234-1234-123456789abc",
	placeholder,
placeholder
placeholder

func TestValidate_ClaudeCLIUserAgent(t *testing.T) {
	v := newTestValidator()

	tests := []struct {
		name string
		ua   string
		want bool
placeholder{
		{"标准版本号", "claude-cli/1.0.0", trueplaceholder,
		{"多位版本号", "claude-cli/12.34.56", trueplaceholder,
		{"大写开头", "Claude-CLI/1.0.0", trueplaceholder,
		{"非 claude-cli", "curl/7.64.1", falseplaceholder,
		{"空 User-Agent", "", falseplaceholder,
		{"部分匹配", "not-claude-cli/1.0.0", falseplaceholder,
		{"缺少版本号", "claude-cli/", falseplaceholder,
		{"版本格式不对", "claude-cli/1.0", falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, v.ValidateUserAgent(tt.ua), "UA: %q", tt.ua)
	placeholder)
placeholder
placeholder

func TestValidate_NonMessagesPath_UAOnly(t *testing.T) {
	v := newTestValidator()

	// 非 messages 路径只检查 UA
	req := httptest.NewRequest("GET", "/v1/models", nil)
	req.Header.Set("User-Agent", "claude-cli/1.0.0")

	result := v.Validate(req, nil)
	require.True(t, result, "非 messages 路径只需 UA 匹配")
placeholder

func TestValidate_NonMessagesPath_InvalidUA(t *testing.T) {
	v := newTestValidator()

	req := httptest.NewRequest("GET", "/v1/models", nil)
	req.Header.Set("User-Agent", "curl/7.64.1")

	result := v.Validate(req, nil)
	require.False(t, result, "UA 不匹配时应返回 false")
placeholder

func TestValidate_MessagesPath_FullValid(t *testing.T) {
	v := newTestValidator()

	req := httptest.NewRequest("POST", "/v1/messages", nil)
	req.Header.Set("User-Agent", "claude-cli/1.0.0")
	req.Header.Set("X-App", "claude-code")
	req.Header.Set("anthropic-beta", "max-tokens-3-5-sonnet-2024-07-15")
	req.Header.Set("anthropic-version", "2023-06-01")

	result := v.Validate(req, validClaudeCodeBody())
	require.True(t, result, "完整有效请求应通过")
placeholder

func TestValidate_MessagesPath_MissingHeaders(t *testing.T) {
	v := newTestValidator()
	body := validClaudeCodeBody()

	tests := []struct {
		name          string
		missingHeader string
placeholder{
		{"缺少 X-App", "X-App"placeholder,
		{"缺少 anthropic-beta", "anthropic-beta"placeholder,
		{"缺少 anthropic-version", "anthropic-version"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/v1/messages", nil)
			req.Header.Set("User-Agent", "claude-cli/1.0.0")
			req.Header.Set("X-App", "claude-code")
			req.Header.Set("anthropic-beta", "beta")
			req.Header.Set("anthropic-version", "2023-06-01")
			req.Header.Del(tt.missingHeader)

			result := v.Validate(req, body)
			require.False(t, result, "缺少 %s 应返回 false", tt.missingHeader)
	placeholder)
placeholder
placeholder

func TestValidate_MessagesPath_InvalidMetadataUserID(t *testing.T) {
	v := newTestValidator()

	tests := []struct {
		name     string
		metadata map[string]any
placeholder{
		{"缺少 metadata", nilplaceholder,
		{"缺少 user_id", map[string]any{"other": "value"placeholderplaceholder,
		{"空 user_id", map[string]any{"user_id": ""placeholderplaceholder,
		{"格式错误", map[string]any{"user_id": "invalid-format"placeholderplaceholder,
		{"hex 长度不足", map[string]any{"user_id": "user_abc_account__session_uuid"placeholderplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/v1/messages", nil)
			req.Header.Set("User-Agent", "claude-cli/1.0.0")
			req.Header.Set("X-App", "claude-code")
			req.Header.Set("anthropic-beta", "beta")
			req.Header.Set("anthropic-version", "2023-06-01")

			body := map[string]any{
				"model": "claude-sonnet-4",
				"system": []any{
					map[string]any{
						"type": "text",
						"text": "You are Claude Code, Anthropic's official CLI for Claude.",
				placeholder,
			placeholder,
		placeholder
			if tt.metadata != nil {
				body["metadata"] = tt.metadata
		placeholder

			result := v.Validate(req, body)
			require.False(t, result, "metadata.user_id: %v", tt.metadata)
	placeholder)
placeholder
placeholder

func TestValidate_MessagesPath_InvalidSystemPrompt(t *testing.T) {
	v := newTestValidator()

	req := httptest.NewRequest("POST", "/v1/messages", nil)
	req.Header.Set("User-Agent", "claude-cli/1.0.0")
	req.Header.Set("X-App", "claude-code")
	req.Header.Set("anthropic-beta", "beta")
	req.Header.Set("anthropic-version", "2023-06-01")

	body := map[string]any{
		"model": "claude-sonnet-4",
		"system": []any{
			map[string]any{
				"type": "text",
				"text": "Generate JSON data for testing database migrations.",
		placeholder,
	placeholder,
		"metadata": map[string]any{
			"user_id": "user_" + "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2" + "_account__session_12345678-1234-1234-1234-123456789abc",
	placeholder,
placeholder

	result := v.Validate(req, body)
	require.False(t, result, "无关系统提示词应返回 false")
placeholder

func TestValidate_MaxTokensOneHaikuBypass(t *testing.T) {
	v := newTestValidator()

	req := httptest.NewRequest("POST", "/v1/messages", nil)
	req.Header.Set("User-Agent", "claude-cli/1.0.0")
	// 不设置 X-App 等头，通过 context 标记为 haiku 探测请求
	ctx := context.WithValue(req.Context(), ctxkey.IsMaxTokensOneHaikuRequest, true)
	req = req.WithContext(ctx)

	// 即使 body 不包含 system prompt，也应通过
	result := v.Validate(req, map[string]any{"model": "claude-3-haiku", "max_tokens": 1placeholder)
	require.True(t, result, "max_tokens=1+haiku 探测请求应绕过严格验证")
placeholder

func TestSystemPromptSimilarity(t *testing.T) {
	v := newTestValidator()

	tests := []struct {
		name   string
		prompt string
		want   bool
placeholder{
		{"精确匹配", "You are Claude Code, Anthropic's official CLI for Claude.", trueplaceholder,
		{"带多余空格", "You  are  Claude  Code,  Anthropic's  official  CLI  for  Claude.", trueplaceholder,
		{"Agent SDK 模板", "You are a Claude agent, built on Anthropic's Claude Agent SDK.", trueplaceholder,
		{"文件搜索专家模板", "You are a file search specialist for Claude Code, Anthropic's official CLI for Claude.", trueplaceholder,
		{"对话摘要模板", "You are a helpful AI assistant tasked with summarizing conversations.", trueplaceholder,
		{"交互式 CLI 模板", "You are an interactive CLI tool that helps users", trueplaceholder,
		{"无关文本", "Write me a poem about cats", falseplaceholder,
		{"空文本", "", falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := map[string]any{
				"model": "claude-sonnet-4",
				"system": []any{
					map[string]any{"type": "text", "text": tt.promptplaceholder,
			placeholder,
		placeholder
			result := v.IncludesClaudeCodeSystemPrompt(body)
			require.Equal(t, tt.want, result, "提示词: %q", tt.prompt)
	placeholder)
placeholder
placeholder

func TestDiceCoefficient(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want float64
		tol  float64
placeholder{
		{"相同字符串", "hello", "hello", 1.0, 0.001placeholder,
		{"完全不同", "abc", "xyz", 0.0, 0.001placeholder,
		{"空字符串", "", "hello", 0.0, 0.001placeholder,
		{"单字符", "a", "b", 0.0, 0.001placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := diceCoefficient(tt.a, tt.b)
			require.InDelta(t, tt.want, result, tt.tol)
	placeholder)
placeholder
placeholder

func TestIsClaudeCodeClient_Context(t *testing.T) {
	ctx := context.Background()

	// 默认应为 false
	require.False(t, IsClaudeCodeClient(ctx))

	// 设置为 true
	ctx = SetClaudeCodeClient(ctx, true)
	require.True(t, IsClaudeCodeClient(ctx))

	// 设置为 false
	ctx = SetClaudeCodeClient(ctx, false)
	require.False(t, IsClaudeCodeClient(ctx))
placeholder

func TestValidate_NilBody_MessagesPath(t *testing.T) {
	v := newTestValidator()

	req := httptest.NewRequest("POST", "/v1/messages", nil)
	req.Header.Set("User-Agent", "claude-cli/1.0.0")
	req.Header.Set("X-App", "claude-code")
	req.Header.Set("anthropic-beta", "beta")
	req.Header.Set("anthropic-version", "2023-06-01")

	result := v.Validate(req, nil)
	require.False(t, result, "nil body 的 messages 请求应返回 false")
placeholder
