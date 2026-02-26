package service

import "testing"

func TestProjectUsageToClaudeMax1H_Conservation(t *testing.T) {
	usage := &ClaudeUsage{
		InputTokens:              1200,
		CacheCreationInputTokens: 0,
		CacheCreation5mTokens:    0,
		CacheCreation1hTokens:    0,
placeholder
	parsed := &ParsedRequest{
		Model: "claude-sonnet-4-5",
		Messages: []any{
			map[string]any{
				"role":    "user",
				"content": "请帮我总结这段代码并给出优化建议",
		placeholder,
	placeholder,
placeholder

	changed := projectUsageToClaudeMax1H(usage, parsed)
	if !changed {
		t.Fatalf("expected usage to be projected")
placeholder

	total := usage.InputTokens + usage.CacheCreation5mTokens + usage.CacheCreation1hTokens
	if total != 1200 {
		t.Fatalf("total tokens changed: got=%d want=%d", total, 1200)
placeholder
	if usage.CacheCreation5mTokens != 0 {
		t.Fatalf("cache_creation_5m should be 0, got=%d", usage.CacheCreation5mTokens)
placeholder
	if usage.InputTokens <= 0 || usage.InputTokens >= 1200 {
		t.Fatalf("simulated input out of range, got=%d", usage.InputTokens)
placeholder
	if usage.CacheCreation1hTokens <= 0 {
		t.Fatalf("cache_creation_1h should be > 0, got=%d", usage.CacheCreation1hTokens)
placeholder
	if usage.CacheCreationInputTokens != usage.CacheCreation1hTokens {
		t.Fatalf("cache_creation_input_tokens mismatch: got=%d want=%d", usage.CacheCreationInputTokens, usage.CacheCreation1hTokens)
placeholder
placeholder

func TestComputeClaudeMaxSimulatedInputTokens_Deterministic(t *testing.T) {
	parsed := &ParsedRequest{
		Model: "claude-opus-4-5",
		Messages: []any{
			map[string]any{
				"role": "user",
				"content": []any{
					map[string]any{"type": "text", "text": "请整理以下日志并定位错误根因"placeholder,
					map[string]any{"type": "tool_use", "name": "grep_logs"placeholder,
			placeholder,
		placeholder,
	placeholder,
placeholder

	got1 := computeClaudeMaxSimulatedInputTokens(4096, parsed)
	got2 := computeClaudeMaxSimulatedInputTokens(4096, parsed)
	if got1 != got2 {
		t.Fatalf("non-deterministic input tokens: %d != %d", got1, got2)
placeholder
placeholder

func TestShouldSimulateClaudeMaxUsage(t *testing.T) {
	group := &Group{
		Platform:                 PlatformAnthropic,
		SimulateClaudeMaxEnabled: true,
placeholder
	input := &RecordUsageInput{
		Result: &ForwardResult{
			Model: "claude-sonnet-4-5",
			Usage: ClaudeUsage{
				InputTokens:              3000,
				CacheCreationInputTokens: 0,
				CacheCreation5mTokens:    0,
				CacheCreation1hTokens:    0,
		placeholder,
	placeholder,
		APIKey: &APIKey{Group: groupplaceholder,
placeholder

	if !shouldSimulateClaudeMaxUsage(input) {
		t.Fatalf("expected simulate=true for claude group without cache creation")
placeholder

	input.Result.Usage.CacheCreationInputTokens = 100
	if shouldSimulateClaudeMaxUsage(input) {
		t.Fatalf("expected simulate=false when cache creation already exists")
placeholder
placeholder
