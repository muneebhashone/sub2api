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
