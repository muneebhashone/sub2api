package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	"github.com/stretchr/testify/require"
)

func TestClaudeCodeValidator_ProbeBypass(t *testing.T) {
	validator := NewClaudeCodeValidator()
	req := httptest.NewRequest(http.MethodPost, "http://example.com/v1/messages", nil)
	req.Header.Set("User-Agent", "claude-cli/1.2.3 (darwin; arm64)")
	req = req.WithContext(context.WithValue(req.Context(), ctxkey.IsMaxTokensOneHaikuRequest, true))

	ok := validator.Validate(req, map[string]any{
		"model":      "claude-haiku-4-5",
		"max_tokens": 1,
placeholder)
	require.True(t, ok)
placeholder

func TestClaudeCodeValidator_ProbeBypassRequiresUA(t *testing.T) {
	validator := NewClaudeCodeValidator()
	req := httptest.NewRequest(http.MethodPost, "http://example.com/v1/messages", nil)
	req.Header.Set("User-Agent", "curl/8.0.0")
	req = req.WithContext(context.WithValue(req.Context(), ctxkey.IsMaxTokensOneHaikuRequest, true))

	ok := validator.Validate(req, map[string]any{
		"model":      "claude-haiku-4-5",
		"max_tokens": 1,
placeholder)
	require.False(t, ok)
placeholder

func TestClaudeCodeValidator_MessagesWithoutProbeStillNeedStrictValidation(t *testing.T) {
	validator := NewClaudeCodeValidator()
	req := httptest.NewRequest(http.MethodPost, "http://example.com/v1/messages", nil)
	req.Header.Set("User-Agent", "claude-cli/1.2.3 (darwin; arm64)")

	ok := validator.Validate(req, map[string]any{
		"model":      "claude-haiku-4-5",
		"max_tokens": 1,
placeholder)
	require.False(t, ok)
placeholder

func TestClaudeCodeValidator_NonMessagesPathUAOnly(t *testing.T) {
	validator := NewClaudeCodeValidator()
	req := httptest.NewRequest(http.MethodPost, "http://example.com/v1/models", nil)
	req.Header.Set("User-Agent", "claude-cli/1.2.3 (darwin; arm64)")

	ok := validator.Validate(req, nil)
	require.True(t, ok)
placeholder

func TestExtractVersion(t *testing.T) {
	v := NewClaudeCodeValidator()
	tests := []struct {
		ua   string
		want string
placeholder{
		{"claude-cli/2.1.22 (darwin; arm64)", "2.1.22"placeholder,
		{"claude-cli/1.0.0", "1.0.0"placeholder,
		{"Claude-CLI/3.10.5 (linux; x86_64)", "3.10.5"placeholder,  // 大小写不敏感
		{"curl/8.0.0", ""placeholder,                                 // 非 Claude CLI
		{"", ""placeholder,                                           // 空字符串
		{"claude-cli/", ""placeholder,                                // 无版本号
		{"claude-cli/2.1.22-beta", "2.1.22"placeholder,               // 带后缀仍提取主版本号
placeholder
	for _, tt := range tests {
		got := v.ExtractVersion(tt.ua)
		require.Equal(t, tt.want, got, "ExtractVersion(%q)", tt.ua)
placeholder
placeholder

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		a, b string
		want int
placeholder{
		{"2.1.0", "2.1.0", 0placeholder,    // 相等
		{"2.1.1", "2.1.0", 1placeholder,    // patch 更大
		{"2.0.0", "2.1.0", -1placeholder,   // minor 更小
		{"3.0.0", "2.99.99", 1placeholder,  // major 更大
		{"1.0.0", "2.0.0", -1placeholder,   // major 更小
		{"0.0.1", "0.0.0", 1placeholder,    // patch 差异
		{"", "1.0.0", -1placeholder,        // 空字符串 vs 正常版本
		{"v2.1.0", "2.1.0", 0placeholder,   // v 前缀处理
placeholder
	for _, tt := range tests {
		got := CompareVersions(tt.a, tt.b)
		require.Equal(t, tt.want, got, "CompareVersions(%q, %q)", tt.a, tt.b)
placeholder
placeholder

func TestSetGetClaudeCodeVersion(t *testing.T) {
	ctx := context.Background()
	require.Equal(t, "", GetClaudeCodeVersion(ctx), "empty context should return empty string")

	ctx = SetClaudeCodeVersion(ctx, "2.1.63")
	require.Equal(t, "2.1.63", GetClaudeCodeVersion(ctx))
placeholder
