package service

import "testing"

func TestResolveOpenAIForwardModel(t *testing.T) {
	tests := []struct {
		name                        string
		account                     *Account
		requestedModel              string
		messagesDispatchMappedModel string
		expectedModel               string
placeholder{
		{
			name: "uses messages dispatch model for known claude family",
			account: &Account{
		placeholderplaceholder,
		placeholder,
			requestedModel:              "claude-opus-4-6",
			messagesDispatchMappedModel: "gpt-4o-mini",
			expectedModel:               "gpt-4o-mini",
	placeholder,
		{
			name: "uses exact messages dispatch model for unknown claude family",
			account: &Account{
		placeholderplaceholder,
		placeholder,
			requestedModel:              "claude-fable-5",
			messagesDispatchMappedModel: " gpt-5.6-sol ",
			expectedModel:               "gpt-5.6-sol",
	placeholder,
		{
			name:                        "nil account uses messages dispatch model",
			requestedModel:              "claude-fable-5",
			messagesDispatchMappedModel: "gpt-5.6-sol",
			expectedModel:               "gpt-5.6-sol",
	placeholder,
		{
			name:           "nil account without messages dispatch keeps requested model",
			requestedModel: "claude-fable-5",
			expectedModel:  "claude-fable-5",
	placeholder,
		{
			name: "ordinary unknown gpt model has no messages dispatch fallback",
			account: &Account{
		placeholderplaceholder,
		placeholder,
			requestedModel: "gpt6",
			expectedModel:  "gpt6",
	placeholder,
		{
			name: "account exact mapping overrides messages dispatch model",
			account: &Account{
		placeholder
					"model_mapping": map[string]any{
						"claude-fable-5": "gpt-5.5",
				placeholder,
			placeholder,
		placeholder,
			requestedModel:              "claude-fable-5",
			messagesDispatchMappedModel: "gpt-5.6-sol",
			expectedModel:               "gpt-5.5",
	placeholder,
		{
			name: "account wildcard mapping overrides messages dispatch model",
			account: &Account{
		placeholder
					"model_mapping": map[string]any{
						"claude-*": "gpt-5.4",
				placeholder,
			placeholder,
		placeholder,
			requestedModel:              "claude-fable-5",
			messagesDispatchMappedModel: "gpt-5.6-sol",
			expectedModel:               "gpt-5.4",
	placeholder,
		{
			name: "account passthrough mapping overrides messages dispatch model",
			account: &Account{
		placeholder
					"model_mapping": map[string]any{
						"claude-fable-5": "claude-fable-5",
				placeholder,
			placeholder,
		placeholder,
			requestedModel:              "claude-fable-5",
			messagesDispatchMappedModel: "gpt-5.6-sol",
			expectedModel:               "claude-fable-5",
	placeholder,
		{
			name: "ordinary codex spark request keeps requested model",
			account: &Account{
		placeholderplaceholder,
		placeholder,
			requestedModel: "gpt-5.3-codex-spark",
			expectedModel:  "gpt-5.3-codex-spark",
	placeholder,
		{
			name: "ordinary gpt-5.5 request keeps requested model",
			account: &Account{
		placeholderplaceholder,
		placeholder,
			requestedModel: "gpt-5.5",
			expectedModel:  "gpt-5.5",
	placeholder,
		{
			name: "ordinary gpt-5.5-pro request keeps requested model",
			account: &Account{
		placeholderplaceholder,
		placeholder,
			requestedModel: "gpt-5.5-pro",
			expectedModel:  "gpt-5.5-pro",
	placeholder,
		{
			name: "ordinary compact-spelled gpt5.5 request keeps requested model",
			account: &Account{
		placeholderplaceholder,
		placeholder,
			requestedModel: "gpt5.5",
			expectedModel:  "gpt5.5",
	placeholder,
		{
			name: "ordinary namespaced gpt-5.5 request keeps requested model",
			account: &Account{
		placeholderplaceholder,
		placeholder,
			requestedModel: "openai/gpt-5.5",
			expectedModel:  "openai/gpt-5.5",
	placeholder,
		{
			name: "ordinary compact gpt-5.5 request keeps requested model",
			account: &Account{
		placeholderplaceholder,
		placeholder,
			requestedModel: "gpt-5.5-openai-compact",
			expectedModel:  "gpt-5.5-openai-compact",
	placeholder,
		{
			name: "whitespace-only messages dispatch model is ignored",
			account: &Account{
		placeholderplaceholder,
		placeholder,
			requestedModel:              "gpt-5.5",
			messagesDispatchMappedModel: "  ",
			expectedModel:               "gpt-5.5",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveOpenAIForwardModel(tt.account, tt.requestedModel, tt.messagesDispatchMappedModel); got != tt.expectedModel {
				t.Fatalf("resolveOpenAIForwardModel(...) = %q, want %q", got, tt.expectedModel)
		placeholder
	placeholder)
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

func TestResolveOpenAIForwardMappedModels_CompactMappingPrecedence(t *testing.T) {
	conflictingMappings := map[string]any{
		"model_mapping":         map[string]any{"gpt-5.5": "gpt-5.4"placeholder,
		"compact_model_mapping": map[string]any{"gpt-5.5": "gpt-5.5-openai-compact"placeholder,
placeholder
	mappedOnlyCompact := map[string]any{
		"model_mapping":         map[string]any{"gpt-5.5": "gpt-5.4"placeholder,
		"compact_model_mapping": map[string]any{"gpt-5.4": "gpt-5.4-openai-compact"placeholder,
placeholder
	tests := []struct {
		name           string
		account        *Account
		requireCompact bool
		wantBilling    string
		wantUpstream   string
placeholder{
		{
			name: "compact uses client-visible model before ordinary mapping",
			account: &Account{Platform: PlatformOpenAI, Type: AccountTypeOAuth,
				Credentials: conflictingMappingsplaceholder,
			requireCompact: true,
			wantBilling:    "gpt-5.4",
			wantUpstream:   "gpt-5.5-openai-compact",
	placeholder,
		{
			name: "non-compact uses ordinary mapping",
			account: &Account{Platform: PlatformOpenAI, Type: AccountTypeOAuth,
				Credentials: conflictingMappingsplaceholder,
			wantBilling:  "gpt-5.4",
			wantUpstream: "gpt-5.4",
	placeholder,
		{
			name: "compact falls back to ordinary mapped model",
			account: &Account{Platform: PlatformOpenAI, Type: AccountTypeOAuth,
				Credentials: mappedOnlyCompactplaceholder,
			requireCompact: true,
			wantBilling:    "gpt-5.4",
			wantUpstream:   "gpt-5.4-openai-compact",
	placeholder,
		{
			name: "passthrough ignores ordinary mapping",
			account: &Account{Platform: PlatformOpenAI, Type: AccountTypeOAuth,
				Credentials: conflictingMappings, Extra: map[string]any{"openai_passthrough": trueplaceholderplaceholder,
			requireCompact: true,
			wantBilling:    "gpt-5.5",
			wantUpstream:   "gpt-5.5-openai-compact",
	placeholder,
		{
			name: "raw chat fallback never applies compact mapping",
			account: &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKey,
				Credentials: conflictingMappings, Extra: map[string]any{"openai_responses_supported": falseplaceholderplaceholder,
			requireCompact: true,
			wantBilling:    "gpt-5.4",
			wantUpstream:   "gpt-5.4",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			billing, upstream := resolveOpenAIForwardMappedModels(tt.account, "gpt-5.5", tt.requireCompact)
			if billing != tt.wantBilling {
				t.Fatalf("billing model = %q, want %q", billing, tt.wantBilling)
		placeholder
			if upstream != tt.wantUpstream {
				t.Fatalf("upstream model = %q, want %q", upstream, tt.wantUpstream)
		placeholder
			if scheduler := resolveOpenAIAccountUpstreamModelForRequest(tt.account, "gpt-5.5", tt.requireCompact); scheduler != upstream {
				t.Fatalf("scheduler model %q disagrees with Forward model %q", scheduler, upstream)
		placeholder
	placeholder)
placeholder
placeholder

func TestCanonicalOpenAIAccountSchedulingModelMatchesForwardSemantics(t *testing.T) {
	tests := []struct {
		name    string
		account *Account
		model   string
		want    string
placeholder{
		{
			name:    "OpenAI OAuth applies Codex alias normalization",
			account: &Account{Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder,
			model:   "gpt-5.6",
			want:    "gpt-5.6-sol",
	placeholder,
		{
			name: "OpenAI passthrough ignores ordinary account mapping",
			account: &Account{Platform: PlatformOpenAI, Type: AccountTypeOAuth,
		placeholder"model_mapping": map[string]any{"public": "private"placeholderplaceholder,
				Extra:       map[string]any{"openai_passthrough": trueplaceholderplaceholder,
			model: "public",
			want:  "public",
	placeholder,
		{
			name:    "Grok OAuth does not inherit OpenAI Codex aliases",
			account: &Account{Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder,
			model:   "gpt-5.6",
			want:    "gpt-5.6",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canonicalOpenAIAccountSchedulingModel(tt.account, tt.model); got != tt.want {
				t.Fatalf("canonical scheduling model = %q, want %q", got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestResolveOpenAIErrorSchedulingModelPrefersActualUpstreamModel(t *testing.T) {
	if got := resolveOpenAIErrorSchedulingModel("gpt-5.4", "gpt-5.5-openai-compact"); got != "gpt-5.5-openai-compact" {
		t.Fatalf("error scheduling model = %q, want compact upstream model", got)
placeholder
	if got := resolveOpenAIErrorSchedulingModel("gpt-5.4", ""); got != "gpt-5.4" {
		t.Fatalf("empty upstream fallback = %q, want billing model", got)
placeholder
placeholder

func TestNormalizeCodexModel(t *testing.T) {
	cases := map[string]string{
		"gpt-5.3-codex-spark":       "gpt-5.3-codex-spark",
		"gpt-5.3-codex-spark-high":  "gpt-5.3-codex-spark",
		"gpt-5.3-codex-spark-xhigh": "gpt-5.3-codex-spark",
		"gpt-5.3":                   "gpt-5.3-codex",
		"gpt-image-2":               "gpt-image-2",
		"gpt-5.4-nano":              "gpt-5.4-nano",
		"gpt-5.4-nano-high":         "gpt-5.4-nano",
		"gpt6":                      "gpt6",
		"claude-opus-4-6":           "claude-opus-4-6",
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
			name:    "oauth routes bare GPT-5.6 alias to Sol",
			account: &Account{Type: AccountTypeOAuthplaceholder,
			model:   "gpt-5.6",
			want:    "gpt-5.6-sol",
	placeholder,
		{
			name:    "oauth routes provider-prefixed GPT-5.6 alias to Sol",
			account: &Account{Type: AccountTypeOAuthplaceholder,
			model:   "openai/gpt-5.6",
			want:    "gpt-5.6-sol",
	placeholder,
		{
			name:    "oauth preserves unknown non codex model",
			account: &Account{Type: AccountTypeOAuthplaceholder,
			model:   "gemini-3-flash-preview",
			want:    "gemini-3-flash-preview",
	placeholder,
		{
			name:    "oauth preserves invalid gpt model",
			account: &Account{Type: AccountTypeOAuthplaceholder,
			model:   "gpt6",
			want:    "gpt6",
	placeholder,
		{
			name:    "oauth normalizes known codex alias",
			account: &Account{Type: AccountTypeOAuthplaceholder,
			model:   "gpt-5.4-high",
			want:    "gpt-5.4",
	placeholder,
		{
			name:    "oauth preserves GPT-5.5 Pro model",
			account: &Account{Type: AccountTypeOAuthplaceholder,
			model:   "openai/gpt-5.5-pro",
			want:    "gpt-5.5-pro",
	placeholder,
		{
			name:    "oauth preserves codex auto review model",
			account: &Account{Type: AccountTypeOAuthplaceholder,
			model:   "codex-auto-review",
			want:    "codex-auto-review",
	placeholder,
		{
			name:    "apikey preserves official bare GPT-5.6 alias",
			account: &Account{Type: AccountTypeAPIKeyplaceholder,
			model:   "gpt-5.6",
			want:    "gpt-5.6",
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

func TestUsageBillingModelCandidatesPreserveCodexAutoReviewModel(t *testing.T) {
	candidates := usageBillingModelCandidates("codex-auto-review")

	expected := []string{"codex-auto-review"placeholder
	if len(candidates) != len(expected) {
		t.Fatalf("usageBillingModelCandidates(codex-auto-review) = %#v, want %#v", candidates, expected)
placeholder
	for i := range expected {
		if candidates[i] != expected[i] {
			t.Fatalf("usageBillingModelCandidates(codex-auto-review) = %#v, want %#v", candidates, expected)
	placeholder
placeholder
placeholder

func TestUsageBillingModelCandidatesPreserveGPT55ProModel(t *testing.T) {
	candidates := usageBillingModelCandidates("openai/gpt-5.5-pro")

	expected := []string{"openai/gpt-5.5-pro", "gpt-5.5-pro"placeholder
	if len(candidates) != len(expected) {
		t.Fatalf("usageBillingModelCandidates(openai/gpt-5.5-pro) = %#v, want %#v", candidates, expected)
placeholder
	for i := range expected {
		if candidates[i] != expected[i] {
			t.Fatalf("usageBillingModelCandidates(openai/gpt-5.5-pro) = %#v, want %#v", candidates, expected)
	placeholder
placeholder
placeholder
