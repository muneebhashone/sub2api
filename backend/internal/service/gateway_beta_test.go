package service

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/claude"

	"github.com/stretchr/testify/require"
)

func TestMergeAnthropicBeta(t *testing.T) {
	got := mergeAnthropicBeta(
		[]string{"oauth-2025-04-20", "placeholder"placeholder,
		"foo, oauth-2025-04-20,bar, foo",
	)
	require.Equal(t, "oauth-2025-04-20,placeholder,foo,bar", got)
placeholder

func TestMergeAnthropicBeta_EmptyIncoming(t *testing.T) {
	got := mergeAnthropicBeta(
		[]string{"oauth-2025-04-20", "placeholder"placeholder,
		"",
	)
	require.Equal(t, "oauth-2025-04-20,placeholder", got)
placeholder

func TestStripBetaTokens(t *testing.T) {
	tests := []struct {
		name   string
		header string
		tokens []string
		want   string
placeholder{
		{
			name:   "single token in middle",
			header: "oauth-2025-04-20,context-1m-2025-08-07,placeholder",
			tokens: []string{"context-1m-2025-08-07"placeholder,
			want:   "oauth-2025-04-20,placeholder",
	placeholder,
		{
			name:   "single token at start",
			header: "context-1m-2025-08-07,oauth-2025-04-20,placeholder",
			tokens: []string{"context-1m-2025-08-07"placeholder,
			want:   "oauth-2025-04-20,placeholder",
	placeholder,
		{
			name:   "single token at end",
			header: "oauth-2025-04-20,placeholder,context-1m-2025-08-07",
			tokens: []string{"context-1m-2025-08-07"placeholder,
			want:   "oauth-2025-04-20,placeholder",
	placeholder,
		{
			name:   "token not present",
			header: "oauth-2025-04-20,placeholder",
			tokens: []string{"context-1m-2025-08-07"placeholder,
			want:   "oauth-2025-04-20,placeholder",
	placeholder,
		{
			name:   "empty header",
			header: "",
			tokens: []string{"context-1m-2025-08-07"placeholder,
			want:   "",
	placeholder,
		{
			name:   "with spaces",
			header: "oauth-2025-04-20, context-1m-2025-08-07 , placeholder",
			tokens: []string{"context-1m-2025-08-07"placeholder,
			want:   "oauth-2025-04-20,placeholder",
	placeholder,
		{
			name:   "only token",
			header: "context-1m-2025-08-07",
			tokens: []string{"context-1m-2025-08-07"placeholder,
			want:   "",
	placeholder,
		{
			name:   "nil tokens",
			header: "oauth-2025-04-20,placeholder",
			tokens: nil,
			want:   "oauth-2025-04-20,placeholder",
	placeholder,
		{
			name:   "multiple tokens removed",
			header: "oauth-2025-04-20,context-1m-2025-08-07,placeholder,fast-mode-2026-02-01",
			tokens: []string{"context-1m-2025-08-07", "fast-mode-2026-02-01"placeholder,
			want:   "oauth-2025-04-20,placeholder",
	placeholder,
		{
			name:   "DroppedBetas is empty (filtering moved to configurable beta policy)",
			header: "oauth-2025-04-20,context-1m-2025-08-07,fast-mode-2026-02-01,placeholder",
			tokens: claude.DroppedBetas,
			want:   "oauth-2025-04-20,context-1m-2025-08-07,fast-mode-2026-02-01,placeholder",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripBetaTokens(tt.header, tt.tokens)
			require.Equal(t, tt.want, got)
	placeholder)
placeholder
placeholder

func TestMergeAnthropicBetaDropping_Context1M(t *testing.T) {
	required := []string{"oauth-2025-04-20", "placeholder"placeholder
	incoming := "context-1m-2025-08-07,foo-beta,oauth-2025-04-20"
	drop := map[string]struct{placeholder{"context-1m-2025-08-07": {placeholderplaceholder

	got := mergeAnthropicBetaDropping(required, incoming, drop)
	require.Equal(t, "oauth-2025-04-20,placeholder,foo-beta", got)
	require.NotContains(t, got, "context-1m-2025-08-07")
placeholder

func TestMergeAnthropicBetaDropping_DroppedBetas(t *testing.T) {
	required := []string{"oauth-2025-04-20", "placeholder"placeholder
	incoming := "context-1m-2025-08-07,fast-mode-2026-02-01,foo-beta,oauth-2025-04-20"
	// DroppedBetas is now empty — filtering moved to configurable beta policy.
	// Without a policy filter set, nothing gets dropped from the static set.
	drop := droppedBetaSet()

	got := mergeAnthropicBetaDropping(required, incoming, drop)
	require.Equal(t, "oauth-2025-04-20,placeholder,context-1m-2025-08-07,fast-mode-2026-02-01,foo-beta", got)
	require.Contains(t, got, "context-1m-2025-08-07")
	require.Contains(t, got, "fast-mode-2026-02-01")
placeholder

func TestFullClaudeCodeMimicryBetas_DoesNotDefaultRedactThinking(t *testing.T) {
	required := claude.FullClaudeCodeMimicryBetas()

	require.NotContains(t, required, claude.BetaRedactThinking)
	require.Contains(t, required, claude.BetaClaudeCode)
	require.Contains(t, required, claude.BetaOAuth)
	require.Contains(t, required, claude.BetaInterleavedThinking)
placeholder

func TestMergeAnthropicBetaDropping_PreservesIncomingRedactThinking(t *testing.T) {
	required := claude.FullClaudeCodeMimicryBetas()
	incoming := claude.BetaRedactThinking

	got := mergeAnthropicBetaDropping(required, incoming, droppedBetaSet())

	require.Contains(t, got, claude.BetaRedactThinking)
placeholder

func TestDroppedBetaSet(t *testing.T) {
	// Base set contains DroppedBetas (now empty — filtering moved to configurable beta policy)
	base := droppedBetaSet()
	require.Len(t, base, len(claude.DroppedBetas))

	// With extra tokens
	extended := droppedBetaSet(claude.BetaClaudeCode)
	require.Contains(t, extended, claude.BetaClaudeCode)
	require.Len(t, extended, len(claude.DroppedBetas)+1)
placeholder

func TestBuildBetaTokenSet(t *testing.T) {
	got := buildBetaTokenSet([]string{"foo", "", "bar", "foo"placeholder)
	require.Len(t, got, 2)
	require.Contains(t, got, "foo")
	require.Contains(t, got, "bar")
	require.NotContains(t, got, "")

	empty := buildBetaTokenSet(nil)
	require.Empty(t, empty)
placeholder

func TestContainsBetaToken(t *testing.T) {
	tests := []struct {
		name   string
		header string
		token  string
		want   bool
placeholder{
		{"present in middle", "oauth-2025-04-20,fast-mode-2026-02-01,placeholder", "fast-mode-2026-02-01", trueplaceholder,
		{"present at start", "fast-mode-2026-02-01,oauth-2025-04-20", "fast-mode-2026-02-01", trueplaceholder,
		{"present at end", "oauth-2025-04-20,fast-mode-2026-02-01", "fast-mode-2026-02-01", trueplaceholder,
		{"only token", "fast-mode-2026-02-01", "fast-mode-2026-02-01", trueplaceholder,
		{"not present", "oauth-2025-04-20,placeholder", "fast-mode-2026-02-01", falseplaceholder,
		{"with spaces", "oauth-2025-04-20, fast-mode-2026-02-01 , placeholder", "fast-mode-2026-02-01", trueplaceholder,
		{"empty header", "", "fast-mode-2026-02-01", falseplaceholder,
		{"empty token", "fast-mode-2026-02-01", "", falseplaceholder,
		{"partial match", "fast-mode-2026-02-01-extra", "fast-mode-2026-02-01", falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsBetaToken(tt.header, tt.token)
			require.Equal(t, tt.want, got)
	placeholder)
placeholder
placeholder

func TestStripBetaTokensWithSet_EmptyDropSet(t *testing.T) {
	header := "oauth-2025-04-20,placeholder"
	got := stripBetaTokensWithSet(header, map[string]struct{placeholder{placeholder)
	require.Equal(t, header, got)
placeholder

func TestIsCountTokensUnsupported404(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		want       bool
placeholder{
		{
			name:       "exact endpoint not found",
			statusCode: 404,
			body:       `{"error":{"message":"Not found: /v1/messages/count_tokens","type":"not_found_error"placeholderplaceholder`,
			want:       true,
	placeholder,
		{
			name:       "contains count_tokens and not found",
			statusCode: 404,
			body:       `{"error":{"message":"count_tokens route not found","type":"not_found_error"placeholderplaceholder`,
			want:       true,
	placeholder,
		{
			name:       "generic 404",
			statusCode: 404,
			body:       `{"error":{"message":"resource not found","type":"not_found_error"placeholderplaceholder`,
			want:       false,
	placeholder,
		{
			name:       "404 with empty error message",
			statusCode: 404,
			body:       `{"error":{"message":"","type":"not_found_error"placeholderplaceholder`,
			want:       false,
	placeholder,
		{
			name:       "non-404 status",
			statusCode: 400,
			body:       `{"error":{"message":"Not found: /v1/messages/count_tokens","type":"invalid_request_error"placeholderplaceholder`,
			want:       false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isCountTokensUnsupported404(tt.statusCode, []byte(tt.body))
			require.Equal(t, tt.want, got)
	placeholder)
placeholder
placeholder

// TestDefaultBetaPolicy_Context1M_Sonnet5Whitelist 验证默认策略下 context-1m-2025-08-07 的分模型行为：
//   - claude-sonnet-5 及后续版本：pass（放行），保留 1M 上下文能力
//   - 其他 sonnet 版本（4.x 及以下）、opus、haiku：filter（过滤），因为上游不支持
func TestDefaultBetaPolicy_Context1M_Sonnet5Whitelist(t *testing.T) {
	settings := DefaultBetaPolicySettings()

	// 找到 context-1m-2025-08-07 规则
	var rule *BetaPolicyRule
	for i := range settings.Rules {
		if settings.Rules[i].BetaToken == "context-1m-2025-08-07" {
			rule = &settings.Rules[i]
			break
	placeholder
placeholder
	require.NotNil(t, rule, "default policy must include context-1m-2025-08-07 rule")
	require.Equal(t, BetaPolicyActionPass, rule.Action, "primary action for whitelisted models is pass")
	require.Equal(t, BetaPolicyActionFilter, rule.FallbackAction, "non-whitelisted models must be filtered")
	require.NotEmpty(t, rule.ModelWhitelist, "context-1m must be scoped to sonnet-5+ via whitelist")

	// 表驱动：模型 → 期望 action
	// 覆盖每种上游路径下的模型 ID 变形：直连 Anthropic API、Vertex AI（"@YYYYMMDD" 后缀）、
	// AWS Bedrock 跨区域推理（us./eu./apac./jp./au./us-gov./global./anthropic. 前缀）。
	cases := []struct {
		model      string
		wantAction string
		desc       string
placeholder{
		// —— 直连 Anthropic API —— sonnet-5 系列应放行
		{"claude-sonnet-5", BetaPolicyActionPass, "sonnet-5 canonical"placeholder,
		{"claude-sonnet-5-20260701", BetaPolicyActionPass, "sonnet-5 dated variant matches wildcard"placeholder,
		{"claude-sonnet-5-thinking", BetaPolicyActionPass, "sonnet-5 thinking variant matches wildcard"placeholder,
		// —— Vertex AI 归一化后的 sonnet-5 —— 也应放行
		{"claude-sonnet-5@20260701", BetaPolicyActionPass, "sonnet-5 Vertex-normalized dated form"placeholder,
		// —— AWS Bedrock 各跨区域前缀 sonnet-5 —— 也应放行
		{"us.anthropic.claude-sonnet-5-v1", BetaPolicyActionPass, "bedrock us. sonnet-5"placeholder,
		{"eu.anthropic.claude-sonnet-5-20260701-v1:0", BetaPolicyActionPass, "bedrock eu. sonnet-5 dated"placeholder,
		{"apac.anthropic.claude-sonnet-5-v1", BetaPolicyActionPass, "bedrock apac. sonnet-5"placeholder,
		{"jp.anthropic.claude-sonnet-5-v1", BetaPolicyActionPass, "bedrock jp. sonnet-5"placeholder,
		{"au.anthropic.claude-sonnet-5-v1", BetaPolicyActionPass, "bedrock au. sonnet-5"placeholder,
		{"us-gov.anthropic.claude-sonnet-5-v1", BetaPolicyActionPass, "bedrock us-gov. sonnet-5"placeholder,
		{"global.anthropic.claude-sonnet-5-v1", BetaPolicyActionPass, "bedrock global. sonnet-5"placeholder,
		{"anthropic.claude-sonnet-5-v1", BetaPolicyActionPass, "bedrock no-region sonnet-5"placeholder,

		// —— sonnet-4.x 及以下必须过滤 ——
		{"claude-sonnet-4-6", BetaPolicyActionFilter, "sonnet-4.6 must be filtered"placeholder,
		{"claude-sonnet-4-5-20250929", BetaPolicyActionFilter, "sonnet-4.5 dated must be filtered"placeholder,
		{"claude-sonnet-4", BetaPolicyActionFilter, "sonnet-4 must be filtered"placeholder,
		{"claude-sonnet-4-5@20250929", BetaPolicyActionFilter, "sonnet-4.5 Vertex format must be filtered"placeholder,
		{"us.anthropic.claude-sonnet-4-6", BetaPolicyActionFilter, "bedrock us. sonnet-4.6 must be filtered"placeholder,
		{"us.anthropic.claude-sonnet-4-5-20250929-v1:0", BetaPolicyActionFilter, "bedrock us. sonnet-4.5 must be filtered"placeholder,
		// —— Opus / Haiku 必须过滤（无 1M） ——
		{"claude-opus-4-8", BetaPolicyActionFilter, "opus must be filtered"placeholder,
		{"claude-opus-4-7", BetaPolicyActionFilter, "opus 4.7 must be filtered"placeholder,
		{"us.anthropic.claude-opus-4-8-v1", BetaPolicyActionFilter, "bedrock opus 4.8 must be filtered"placeholder,
		{"claude-haiku-4-5", BetaPolicyActionFilter, "haiku must be filtered"placeholder,
		{"us.anthropic.placeholder-v1:0", BetaPolicyActionFilter, "bedrock haiku must be filtered"placeholder,
		{"claude-3-5-sonnet-20241022", BetaPolicyActionFilter, "legacy sonnet 3.5 must be filtered"placeholder,
		// —— 特殊边界：不应把 "claude-sonnet-50" / "claude-sonnet-5.1" 之类意外命名误放行 ——
		{"claude-sonnet-50", BetaPolicyActionFilter, "must not over-match a hypothetical sonnet-50"placeholder,
placeholder

	for _, tc := range cases {
		t.Run(tc.model, func(t *testing.T) {
			action, _ := resolveRuleAction(*rule, tc.model)
			require.Equal(t, tc.wantAction, action, tc.desc)
	placeholder)
placeholder
placeholder
