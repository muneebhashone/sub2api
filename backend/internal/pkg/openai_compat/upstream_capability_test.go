package openai_compat

import "testing"

func TestResolveResponsesSupport(t *testing.T) {
	tests := []struct {
		name  string
		extra map[string]any
		want  AccountResponsesSupport
placeholder{
		{"nil extra", nil, ResponsesSupportUnknownplaceholder,
		{"empty extra", map[string]any{placeholder, ResponsesSupportUnknownplaceholder,
		{"key missing", map[string]any{"other": "value"placeholder, ResponsesSupportUnknownplaceholder,
		{"value true", map[string]any{ExtraKeyResponsesSupported: trueplaceholder, ResponsesSupportYesplaceholder,
		{"value false", map[string]any{ExtraKeyResponsesSupported: falseplaceholder, ResponsesSupportNoplaceholder,
		{"value wrong type string", map[string]any{ExtraKeyResponsesSupported: "true"placeholder, ResponsesSupportUnknownplaceholder,
		{"value wrong type number", map[string]any{ExtraKeyResponsesSupported: 1placeholder, ResponsesSupportUnknownplaceholder,
		{"value nil", map[string]any{ExtraKeyResponsesSupported: nilplaceholder, ResponsesSupportUnknownplaceholder,
		{"force responses", map[string]any{ExtraKeyResponsesMode: string(ResponsesSupportModeForceResponses)placeholder, ResponsesSupportYesplaceholder,
		{"force chat completions", map[string]any{ExtraKeyResponsesMode: string(ResponsesSupportModeForceChatCompletions)placeholder, ResponsesSupportNoplaceholder,
		{"auto follows probe", map[string]any{ExtraKeyResponsesMode: string(ResponsesSupportModeAuto), ExtraKeyResponsesSupported: falseplaceholder, ResponsesSupportNoplaceholder,
		{"invalid mode follows probe", map[string]any{ExtraKeyResponsesMode: "bogus", ExtraKeyResponsesSupported: trueplaceholder, ResponsesSupportYesplaceholder,
		{"force responses overrides probe false", map[string]any{ExtraKeyResponsesMode: string(ResponsesSupportModeForceResponses), ExtraKeyResponsesSupported: falseplaceholder, ResponsesSupportYesplaceholder,
		{"force chat completions overrides probe true", map[string]any{ExtraKeyResponsesMode: string(ResponsesSupportModeForceChatCompletions), ExtraKeyResponsesSupported: trueplaceholder, ResponsesSupportNoplaceholder,
placeholder

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ResolveResponsesSupport(tc.extra)
			if got != tc.want {
				t.Errorf("ResolveResponsesSupport(%v) = %v, want %v", tc.extra, got, tc.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestShouldUseResponsesAPI(t *testing.T) {
	tests := []struct {
		name  string
		extra map[string]any
		want  bool
placeholder{
		// 关键不变量：未探测必须返回 true（保留旧行为）
		{"unknown defaults to true (preserve old behavior)", nil, trueplaceholder,
		{"unknown empty defaults to true", map[string]any{placeholder, trueplaceholder,
		{"unknown wrong type defaults to true", map[string]any{ExtraKeyResponsesSupported: "yes"placeholder, trueplaceholder,

		// 已探测：标记决定
		{"explicitly supported", map[string]any{ExtraKeyResponsesSupported: trueplaceholder, trueplaceholder,
		{"explicitly unsupported", map[string]any{ExtraKeyResponsesSupported: falseplaceholder, falseplaceholder,

		// 手动覆盖：覆盖自动探测结果
		{"force responses overrides unsupported probe", map[string]any{ExtraKeyResponsesMode: string(ResponsesSupportModeForceResponses), ExtraKeyResponsesSupported: falseplaceholder, trueplaceholder,
		{"force chat completions overrides supported probe", map[string]any{ExtraKeyResponsesMode: string(ResponsesSupportModeForceChatCompletions), ExtraKeyResponsesSupported: trueplaceholder, falseplaceholder,
placeholder

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ShouldUseResponsesAPI(tc.extra)
			if got != tc.want {
				t.Errorf("ShouldUseResponsesAPI(%v) = %v, want %v", tc.extra, got, tc.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestNormalizeResponsesSupportMode(t *testing.T) {
	tests := []struct {
		name string
		mode string
		want ResponsesSupportMode
placeholder{
		{"empty", "", ResponsesSupportModeAutoplaceholder,
		{"auto", "auto", ResponsesSupportModeAutoplaceholder,
		{"force responses", "force_responses", ResponsesSupportModeForceResponsesplaceholder,
		{"force chat completions", "force_chat_completions", ResponsesSupportModeForceChatCompletionsplaceholder,
		{"invalid", "enabled", ResponsesSupportModeAutoplaceholder,
placeholder

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := NormalizeResponsesSupportMode(tc.mode)
			if got != tc.want {
				t.Errorf("NormalizeResponsesSupportMode(%q) = %q, want %q", tc.mode, got, tc.want)
		placeholder
	placeholder)
placeholder
placeholder
