package service

import "testing"

func TestResolveThinkingProtocol(t *testing.T) {
	tests := []struct {
		name    string
		modelID string
		want    ThinkingProtocol
placeholder{
		// Anthropic 官方
		{"claude-sonnet-4-5", "claude-sonnet-4-5", ThinkingProtocolAnthropicStrictplaceholder,
		{"claude-opus-4-5", "claude-opus-4-5-20251101", ThinkingProtocolAnthropicStrictplaceholder,
		{"claude-haiku full id", "placeholder", ThinkingProtocolAnthropicStrictplaceholder,
		{"opus short", "opus-4-5", ThinkingProtocolAnthropicStrictplaceholder,
		{"sonnet short", "sonnet-4-5", ThinkingProtocolAnthropicStrictplaceholder,
		{"haiku short", "haiku-4-5", ThinkingProtocolAnthropicStrictplaceholder,
		{"upper case Claude", "Claude-Sonnet-4-5", ThinkingProtocolAnthropicStrictplaceholder,

		// 第三方兼容上游
		{"deepseek-v4-pro", "deepseek-v4-pro", ThinkingProtocolPassbackRequiredplaceholder,
		{"deepseek-r2-thinking", "deepseek-r2-thinking", ThinkingProtocolPassbackRequiredplaceholder,
		{"kimi-coding", "kimi-coding-v2", ThinkingProtocolPassbackRequiredplaceholder,
		{"kimi-k2-thinking", "kimi-k2-thinking", ThinkingProtocolPassbackRequiredplaceholder,
		{"kimi-k3 platform", "kimi-k3", ThinkingProtocolPassbackRequiredplaceholder,
		{"kimi-k3 anthropic 1m", "kimi-k3[1m]", ThinkingProtocolPassbackRequiredplaceholder,
		{"kimi code bare k3", "k3", ThinkingProtocolPassbackRequiredplaceholder,
		{"kimi code bare k3-256k", "k3-256k", ThinkingProtocolPassbackRequiredplaceholder,
		{"moonshot-v1", "moonshot-v1-32k", ThinkingProtocolPassbackRequiredplaceholder,
		{"glm-5.1", "glm-5.1", ThinkingProtocolPassbackRequiredplaceholder,
		{"qwen-2 thinking variant", "qwen-2-72b-thinking", ThinkingProtocolPassbackRequiredplaceholder,
		{"qwen3 thinking (real Alibaba naming)", "qwen3-235b-a22b-thinking-2507", ThinkingProtocolPassbackRequiredplaceholder,
		{"qwen3-next thinking", "qwen3-next-80b-a3b-thinking", ThinkingProtocolPassbackRequiredplaceholder,
		{"upper case Deepseek", "DeepSeek-V4-Pro", ThinkingProtocolPassbackRequiredplaceholder,

		// MiniMax M 系列（Anthropic 兼容端点要求 thinking round-trip）
		{"MiniMax-M2 (case-sensitive original)", "MiniMax-M2", ThinkingProtocolPassbackRequiredplaceholder,
		{"MiniMax-M2.1", "MiniMax-M2.1", ThinkingProtocolPassbackRequiredplaceholder,
		{"MiniMax-M2.5", "MiniMax-M2.5", ThinkingProtocolPassbackRequiredplaceholder,
		{"MiniMax-M2.7", "MiniMax-M2.7", ThinkingProtocolPassbackRequiredplaceholder,
		{"MiniMax-M2.7-highspeed", "MiniMax-M2.7-highspeed", ThinkingProtocolPassbackRequiredplaceholder,
		{"minimax-m2 lowercase", "minimax-m2", ThinkingProtocolPassbackRequiredplaceholder,

		// 未知 / 保守
		{"empty", "", ThinkingProtocolUnknownplaceholder,
		{"gpt-5", "gpt-5.1", ThinkingProtocolUnknownplaceholder,
		{"gemini", "gemini-3-pro-preview", ThinkingProtocolUnknownplaceholder,
		{"qwen3 non-thinking", "qwen3-32b", ThinkingProtocolUnknownplaceholder,
		{"qwen2 non-thinking", "qwen-2-72b", ThinkingProtocolUnknownplaceholder,
		{"random vendor", "yi-large", ThinkingProtocolUnknownplaceholder,
		// 相似但未知的 k3 型号：不得因含 k3 被宽泛匹配为 passback-required
		{"k3-like unknown", "foo-k3-bar", ThinkingProtocolUnknownplaceholder,
		// MiniMax 非 M 系列（如 abab、speech 等其他产品线）—— unknown
		{"minimax abab non-M", "abab6.5-chat", ThinkingProtocolUnknownplaceholder,
		// Doubao 走 OpenAI 协议，不属于本网关 Anthropic 路径——归 unknown
		{"doubao goes via openai", "doubao-1-5-thinking-vision-pro-250428", ThinkingProtocolUnknownplaceholder,
		// Hunyuan T1 未暴露 Anthropic 端点——归 unknown
		{"hunyuan t1 no anthropic endpoint", "hunyuan-t1", ThinkingProtocolUnknownplaceholder,
		{"hy-t1 short alias", "hy-t1", ThinkingProtocolUnknownplaceholder,
		// claude-something 但不是 anthropic 官方命名风格——也归 strict（前缀匹配优先）
		{"weird claude prefix", "claude-experimental-fork", ThinkingProtocolAnthropicStrictplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveThinkingProtocol(tt.modelID)
			if got != tt.want {
				t.Errorf("ResolveThinkingProtocol(%q) = %v, want %v", tt.modelID, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestShouldPreFilterThinkingBlocks(t *testing.T) {
	tests := []struct {
		modelID string
		want    bool
placeholder{
		{"claude-sonnet-4-5", trueplaceholder,
		{"deepseek-v4-pro", falseplaceholder,
		{"kimi-coding", falseplaceholder,
		{"glm-5.1", falseplaceholder,
		{"gpt-5.1", falseplaceholder,
		{"", falseplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.modelID, func(t *testing.T) {
			if got := ShouldPreFilterThinkingBlocks(tt.modelID); got != tt.want {
				t.Errorf("ShouldPreFilterThinkingBlocks(%q) = %v, want %v", tt.modelID, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestShouldRectifyThinkingSignatureError(t *testing.T) {
	if !ShouldRectifyThinkingSignatureError("claude-sonnet-4-5") {
		t.Error("anthropic-strict should rectify signature error")
placeholder
	if ShouldRectifyThinkingSignatureError("deepseek-v4-pro") {
		t.Error("passback-required must NOT rectify (would break protocol contract)")
placeholder
	if ShouldRectifyThinkingSignatureError("gpt-5.1") {
		t.Error("unknown should NOT rectify (conservative default)")
placeholder
	if ShouldRectifyThinkingSignatureError("") {
		t.Error("empty model id should NOT rectify")
placeholder
placeholder

// ShouldApplyRetryFilters 与 ShouldPreFilterThinkingBlocks 必须语义一致：
// 仅 anthropic-strict 走变形，避免预过滤跳过但 retry 路径反而剥离的语义裂缝。
func TestShouldApplyRetryFiltersMirrorsPreFilter(t *testing.T) {
	models := []string{
		"claude-sonnet-4-5", "claude-opus-4-5-20251101", "haiku-4-5",
		"deepseek-v4-pro", "kimi-coding", "glm-5.1",
		"qwen3-235b-a22b-thinking-2507", "qwen3-32b",
		"gpt-5.1", "gemini-3-pro-preview", "yi-large", "",
placeholder
	for _, m := range models {
		t.Run(m, func(t *testing.T) {
			if got := ShouldApplyRetryFilters(m); got != ShouldPreFilterThinkingBlocks(m) {
				t.Errorf("ShouldApplyRetryFilters(%q)=%v but ShouldPreFilterThinkingBlocks=%v — must match",
					m, got, ShouldPreFilterThinkingBlocks(m))
		placeholder
	placeholder)
placeholder
placeholder
