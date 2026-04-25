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
			requestedModel:     "claude-opus-4-6",
			defaultMappedModel: "gpt-4o-mini",
			expectedModel:      "gpt-4o-mini",
	placeholder,
		{
			name: "preserves explicit gpt-5.4 instead of group default",
			account: &Account{
		placeholderplaceholder,
		placeholder,
			requestedModel:     "gpt-5.4",
			defaultMappedModel: "gpt-4o-mini",
			expectedModel:      "gpt-5.4",
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
		{
			name: "preserves codex spark instead of group default",
			account: &Account{
		placeholderplaceholder,
		placeholder,
			requestedModel:     "gpt-5.3-codex-spark",
			defaultMappedModel: "gpt-5.4",
			expectedModel:      "gpt-5.3-codex-spark",
	placeholder,
		{
			name: "preserves gpt-5.5 instead of group default",
			account: &Account{
		placeholderplaceholder,
		placeholder,
			requestedModel:     "gpt-5.5",
			defaultMappedModel: "gpt-5.4",
			expectedModel:      "gpt-5.5",
	placeholder,
		{
			name: "preserves openai namespaced gpt-5.5 instead of group default",
			account: &Account{
		placeholderplaceholder,
		placeholder,
			requestedModel:     "openai/gpt-5.5",
			defaultMappedModel: "gpt-5.4",
			expectedModel:      "openai/gpt-5.5",
	placeholder,
		{
			name: "preserves compact gpt-5.5 instead of group default",
			account: &Account{
		placeholderplaceholder,
		placeholder,
			requestedModel:     "gpt-5.5-openai-compact",
			defaultMappedModel: "gpt-5.4",
			expectedModel:      "gpt-5.5-openai-compact",
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

func TestResolveOpenAIForwardModel_PreventsClaudeModelFromFallingBackToGpt54(t *testing.T) {
	account := &Account{
placeholderplaceholder,
placeholder

	withoutDefault := normalizeCodexModel(resolveOpenAIForwardModel(account, "claude-opus-4-6", ""))
	if withoutDefault != "gpt-5.4" {
		t.Fatalf("normalizeCodexModel(...) = %q, want %q", withoutDefault, "gpt-5.4")
placeholder

	withDefault := normalizeCodexModel(resolveOpenAIForwardModel(account, "claude-opus-4-6", "gpt-5.4"))
	if withDefault != "gpt-5.4" {
		t.Fatalf("normalizeCodexModel(...) = %q, want %q", withDefault, "gpt-5.4")
placeholder
placeholder

func TestResolveOpenAICompactForwardModel(t *testing.T) {
	tests := []struct {
		name          string
		account       *Account
		model         string
		expectedModel string
placeholder{
		{
			name:          "nil account keeps original model",
			account:       nil,
			model:         "gpt-5.4",
			expectedModel: "gpt-5.4",
	placeholder,
		{
			name: "missing compact mapping keeps original model",
			account: &Account{
		placeholderplaceholder,
		placeholder,
			model:         "gpt-5.4",
			expectedModel: "gpt-5.4",
	placeholder,
		{
			name: "exact compact mapping overrides model",
			account: &Account{
		placeholder
					"compact_model_mapping": map[string]any{
						"gpt-5.4": "gpt-5.4-openai-compact",
				placeholder,
			placeholder,
		placeholder,
			model:         "gpt-5.4",
			expectedModel: "gpt-5.4-openai-compact",
	placeholder,
		{
			name: "wildcard compact mapping overrides model",
			account: &Account{
		placeholder
					"compact_model_mapping": map[string]any{
						"gpt-5.*": "gpt-5-openai-compact",
				placeholder,
			placeholder,
		placeholder,
			model:         "gpt-5.4",
			expectedModel: "gpt-5-openai-compact",
	placeholder,
		{
			name: "passthrough compact mapping remains unchanged",
			account: &Account{
		placeholder
					"compact_model_mapping": map[string]any{
						"gpt-5.4": "gpt-5.4",
				placeholder,
			placeholder,
		placeholder,
			model:         "gpt-5.4",
			expectedModel: "gpt-5.4",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveOpenAICompactForwardModel(tt.account, tt.model); got != tt.expectedModel {
				t.Fatalf("resolveOpenAICompactForwardModel(...) = %q, want %q", got, tt.expectedModel)
		placeholder
	placeholder)
placeholder
placeholder

func TestNormalizeCodexModel(t *testing.T) {
	cases := map[string]string{
		"gpt-5.3-codex-spark":       "gpt-5.3-codex-spark",
		"gpt-5.3-codex-spark-high":  "gpt-5.3-codex-spark",
		"gpt-5.3-codex-spark-xhigh": "gpt-5.3-codex-spark",
		"gpt-5.3":                   "gpt-5.3-codex",
		"gpt-image-2":               "gpt-image-2",
placeholder

	for input, expected := range cases {
		if got := normalizeCodexModel(input); got != expected {
			t.Fatalf("normalizeCodexModel(%q) = %q, want %q", input, got, expected)
	placeholder
placeholder
placeholder

func TestNormalizeOpenAIModelForUpstream(t *testing.T) {
	tests := []struct {
		name    string
		account *Account
		model   string
		want    string
placeholder{
		{
			name:    "oauth keeps codex normalization behavior",
			account: &Account{Type: AccountTypeOAuthplaceholder,
			model:   "gemini-3-flash-preview",
			want:    "gpt-5.4",
	placeholder,
		{
			name:    "apikey preserves custom compatible model",
			account: &Account{Type: AccountTypeAPIKeyplaceholder,
			model:   "gemini-3-flash-preview",
			want:    "gemini-3-flash-preview",
	placeholder,
		{
			name:    "apikey preserves official non codex model",
			account: &Account{Type: AccountTypeAPIKeyplaceholder,
			model:   "gpt-4.1",
			want:    "gpt-4.1",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeOpenAIModelForUpstream(tt.account, tt.model); got != tt.want {
				t.Fatalf("normalizeOpenAIModelForUpstream(...) = %q, want %q", got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder
