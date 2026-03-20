package service

import "testing"

func TestResolveOpenAIForwardModel(t *testing.T) {
	tests := []struct {
		name               string
		account            *Account
		requestedModel     string
		defaultMappedModel string
		expectedModel      string
placeholder{
		{
			name: "falls back to group default when account has no mapping",
			account: &Account{
		placeholderplaceholder,
		placeholder,
			requestedModel:     "gpt-5.4",
			defaultMappedModel: "gpt-4o-mini",
			expectedModel:      "gpt-4o-mini",
	placeholder,
		{
			name: "preserves exact passthrough mapping instead of group default",
			account: &Account{
		placeholder
					"model_mapping": map[string]any{
						"gpt-5.4": "gpt-5.4",
				placeholder,
			placeholder,
		placeholder,
			requestedModel:     "gpt-5.4",
			defaultMappedModel: "gpt-4o-mini",
			expectedModel:      "gpt-5.4",
	placeholder,
		{
			name: "preserves wildcard passthrough mapping instead of group default",
			account: &Account{
		placeholder
					"model_mapping": map[string]any{
						"gpt-*": "gpt-5.4",
				placeholder,
			placeholder,
		placeholder,
			requestedModel:     "gpt-5.4",
			defaultMappedModel: "gpt-4o-mini",
			expectedModel:      "gpt-5.4",
	placeholder,
		{
			name: "uses account remap when explicit target differs",
			account: &Account{
		placeholder
					"model_mapping": map[string]any{
						"gpt-5": "gpt-5.4",
				placeholder,
			placeholder,
		placeholder,
			requestedModel:     "gpt-5",
			defaultMappedModel: "gpt-4o-mini",
			expectedModel:      "gpt-5.4",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveOpenAIForwardModel(tt.account, tt.requestedModel, tt.defaultMappedModel); got != tt.expectedModel {
				t.Fatalf("resolveOpenAIForwardModel(...) = %q, want %q", got, tt.expectedModel)
		placeholder
	placeholder)
placeholder
placeholder

func TestResolveOpenAIForwardModel_PreventsClaudeModelFromFallingBackToGpt51(t *testing.T) {
	account := &Account{
placeholderplaceholder,
placeholder

	withoutDefault := resolveOpenAIForwardModel(account, "claude-opus-4-6", "")
	if got := normalizeCodexModel(withoutDefault); got != "gpt-5.1" {
		t.Fatalf("normalizeCodexModel(%q) = %q, want %q", withoutDefault, got, "gpt-5.1")
placeholder

	withDefault := resolveOpenAIForwardModel(account, "claude-opus-4-6", "gpt-5.4")
	if got := normalizeCodexModel(withDefault); got != "gpt-5.4" {
		t.Fatalf("normalizeCodexModel(%q) = %q, want %q", withDefault, got, "gpt-5.4")
placeholder
placeholder
