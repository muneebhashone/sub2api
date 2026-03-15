//go:build unit

package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestParseGatewayRequest(t *testing.T) {
	body := []byte(`{"model":"claude-3-7-sonnet","stream":true,"metadata":{"user_id":"session_123e4567-e89b-12d3-a456-426614174000"placeholder,"system":[{"type":"text","text":"hello","cache_control":{"type":"ephemeral"placeholderplaceholder],"messages":[{"content":"hi"placeholder]placeholder`)
	parsed, err := ParseGatewayRequest(body, "")
placeholder
	require.Equal(t, "claude-3-7-sonnet", parsed.Model)
	require.True(t, parsed.Stream)
	require.Equal(t, "session_123e4567-e89b-12d3-a456-426614174000", parsed.MetadataUserID)
	require.True(t, parsed.HasSystem)
	require.NotNil(t, parsed.System)
	require.Len(t, parsed.Messages, 1)
	require.False(t, parsed.ThinkingEnabled)
placeholder

func TestParseGatewayRequest_ThinkingEnabled(t *testing.T) {
	body := []byte(`{"model":"claude-sonnet-4-5","thinking":{"type":"enabled"placeholder,"messages":[{"content":"hi"placeholder]placeholder`)
	parsed, err := ParseGatewayRequest(body, "")
placeholder
	require.Equal(t, "claude-sonnet-4-5", parsed.Model)
	require.True(t, parsed.ThinkingEnabled)
placeholder

func TestParseGatewayRequest_ThinkingAdaptiveEnabled(t *testing.T) {
	body := []byte(`{"model":"claude-sonnet-4-5","thinking":{"type":"adaptive"placeholder,"messages":[{"content":"hi"placeholder]placeholder`)
	parsed, err := ParseGatewayRequest(body, "")
placeholder
	require.Equal(t, "claude-sonnet-4-5", parsed.Model)
	require.True(t, parsed.ThinkingEnabled)
placeholder

func TestParseGatewayRequest_MaxTokens(t *testing.T) {
	body := []byte(`{"model":"claude-haiku-4-5","max_tokens":1placeholder`)
	parsed, err := ParseGatewayRequest(body, "")
placeholder
	require.Equal(t, 1, parsed.MaxTokens)
placeholder

func TestParseGatewayRequest_MaxTokensNonIntegralIgnored(t *testing.T) {
	body := []byte(`{"model":"claude-haiku-4-5","max_tokens":placeholder`)
	parsed, err := ParseGatewayRequest(body, "")
placeholder
	require.Equal(t, 0, parsed.MaxTokens)
placeholder

func TestParseGatewayRequest_SystemNull(t *testing.T) {
	body := []byte(`{"model":"claude-3","system":nullplaceholder`)
	parsed, err := ParseGatewayRequest(body, "")
placeholder
	// 显式传入 system:null 也应视为“字段已存在”，避免默认 system 被注入。
	require.True(t, parsed.HasSystem)
	require.Nil(t, parsed.System)
placeholder

func TestParseGatewayRequest_InvalidModelType(t *testing.T) {
	body := []byte(`{"model":placeholder`)
	_, err := ParseGatewayRequest(body, "")
placeholder
placeholder

func TestParseGatewayRequest_InvalidStreamType(t *testing.T) {
	body := []byte(`{"stream":"true"placeholder`)
	_, err := ParseGatewayRequest(body, "")
placeholder
placeholder

// ============ Gemini 原生格式解析测试 ============

func TestParseGatewayRequest_GeminiContents(t *testing.T) {
	body := []byte(`{
		"contents": [
			{"role": "user", "parts": [{"text": "Hello"placeholder]placeholder,
			{"role": "model", "parts": [{"text": "Hi there"placeholder]placeholder,
			{"role": "user", "parts": [{"text": "How are you?"placeholder]placeholder
		]
placeholder`)
	parsed, err := ParseGatewayRequest(body, domain.PlatformGemini)
placeholder
	require.Len(t, parsed.Messages, 3, "should parse contents as Messages")
	require.False(t, parsed.HasSystem, "Gemini format should not set HasSystem")
	require.Nil(t, parsed.System, "no systemInstruction means nil System")
placeholder

func TestParseGatewayRequest_GeminiSystemInstruction(t *testing.T) {
	body := []byte(`{
		"systemInstruction": {
			"parts": [{"text": "You are a helpful assistant."placeholder]
	placeholder,
		"contents": [
			{"role": "user", "parts": [{"text": "Hello"placeholder]placeholder
		]
placeholder`)
	parsed, err := ParseGatewayRequest(body, domain.PlatformGemini)
placeholder
	require.NotNil(t, parsed.System, "should parse systemInstruction.parts as System")
	parts, ok := parsed.System.([]any)
	require.True(t, ok)
	require.Len(t, parts, 1)
	partMap, ok := parts[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "You are a helpful assistant.", partMap["text"])
	require.Len(t, parsed.Messages, 1)
placeholder

func TestParseGatewayRequest_GeminiWithModel(t *testing.T) {
	body := []byte(`{
		"model": "gemini-2.5-pro",
		"contents": [{"role": "user", "parts": [{"text": "test"placeholder]placeholder]
placeholder`)
	parsed, err := ParseGatewayRequest(body, domain.PlatformGemini)
placeholder
	require.Equal(t, "gemini-2.5-pro", parsed.Model)
	require.Len(t, parsed.Messages, 1)
placeholder

func TestParseGatewayRequest_GeminiIgnoresAnthropicFields(t *testing.T) {
	// Gemini 格式下 system/messages 字段应被忽略
	body := []byte(`{
		"system": "should be ignored",
		"messages": [{"role": "user", "content": "ignored"placeholder],
		"contents": [{"role": "user", "parts": [{"text": "real content"placeholder]placeholder]
placeholder`)
	parsed, err := ParseGatewayRequest(body, domain.PlatformGemini)
placeholder
	require.False(t, parsed.HasSystem, "Gemini protocol should not parse Anthropic system field")
	require.Nil(t, parsed.System, "no systemInstruction = nil System")
	require.Len(t, parsed.Messages, 1, "should use contents, not messages")
placeholder

func TestParseGatewayRequest_GeminiEmptyContents(t *testing.T) {
	body := []byte(`{"contents": []placeholder`)
	parsed, err := ParseGatewayRequest(body, domain.PlatformGemini)
placeholder
	require.Empty(t, parsed.Messages)
placeholder

func TestParseGatewayRequest_GeminiNoContents(t *testing.T) {
	body := []byte(`{"model": "gemini-2.5-flash"placeholder`)
	parsed, err := ParseGatewayRequest(body, domain.PlatformGemini)
placeholder
	require.Nil(t, parsed.Messages)
	require.Equal(t, "gemini-2.5-flash", parsed.Model)
placeholder

func TestParseGatewayRequest_AnthropicIgnoresGeminiFields(t *testing.T) {
	// Anthropic 格式下 contents/systemInstruction 字段应被忽略
	body := []byte(`{
		"system": "real system",
		"messages": [{"role": "user", "content": "real content"placeholder],
		"contents": [{"role": "user", "parts": [{"text": "ignored"placeholder]placeholder],
		"systemInstruction": {"parts": [{"text": "ignored"placeholder]placeholder
placeholder`)
	parsed, err := ParseGatewayRequest(body, domain.PlatformAnthropic)
placeholder
	require.True(t, parsed.HasSystem)
	require.Equal(t, "real system", parsed.System)
	require.Len(t, parsed.Messages, 1)
	msg, ok := parsed.Messages[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "real content", msg["content"])
placeholder

func TestFilterThinkingBlocks(t *testing.T) {
	containsThinkingBlock := func(body []byte) bool {
		var req map[string]any
		if err := json.Unmarshal(body, &req); err != nil {
			return false
	placeholder
		messages, ok := req["messages"].([]any)
		if !ok {
			return false
	placeholder
		for _, msg := range messages {
			msgMap, ok := msg.(map[string]any)
			if !ok {
				continue
		placeholder
			content, ok := msgMap["content"].([]any)
			if !ok {
				continue
		placeholder
			for _, block := range content {
				blockMap, ok := block.(map[string]any)
				if !ok {
					continue
			placeholder
				blockType, _ := blockMap["type"].(string)
				if blockType == "thinking" {
					return true
			placeholder
				if blockType == "" {
					if _, hasThinking := blockMap["thinking"]; hasThinking {
						return true
				placeholder
			placeholder
		placeholder
	placeholder
		return false
placeholder

	tests := []struct {
		name         string
		input        string
		shouldFilter bool
		expectError  bool
placeholder{
		{
			name:         "filters thinking blocks",
			input:        `{"model":"claude-3-5-sonnet-20241022","messages":[{"role":"user","content":[{"type":"text","text":"Hello"placeholder,{"type":"thinking","thinking":"internal","signature":"invalid"placeholder,{"type":"text","text":"World"placeholder]placeholder]placeholder`,
			shouldFilter: true,
	placeholder,
		{
			name:         "does not filter signed thinking blocks when thinking adaptive",
			input:        `{"thinking":{"type":"adaptive"placeholder,"messages":[{"role":"assistant","content":[{"type":"thinking","thinking":"ok","signature":"sig_real_123"placeholder,{"type":"text","text":"B"placeholder]placeholder]placeholder`,
			shouldFilter: false,
	placeholder,
		{
			name:         "filters unsigned thinking blocks when thinking adaptive",
			input:        `{"thinking":{"type":"adaptive"placeholder,"messages":[{"role":"assistant","content":[{"type":"thinking","thinking":"internal","signature":""placeholder,{"type":"text","text":"B"placeholder]placeholder]placeholder`,
			shouldFilter: true,
	placeholder,
		{
			name:         "handles no thinking blocks",
			input:        `{"model":"claude-3-5-sonnet-20241022","messages":[{"role":"user","content":[{"type":"text","text":"Hello"placeholder]placeholder]placeholder`,
			shouldFilter: false,
	placeholder,
		{
			name:         "handles invalid JSON gracefully",
			input:        `{invalid json`,
			shouldFilter: false,
			expectError:  true,
	placeholder,
		{
			name:         "handles multiple messages with thinking blocks",
			input:        `{"messages":[{"role":"user","content":[{"type":"text","text":"A"placeholder]placeholder,{"role":"assistant","content":[{"type":"thinking","thinking":"think"placeholder,{"type":"text","text":"B"placeholder]placeholder]placeholder`,
			shouldFilter: true,
	placeholder,
		{
			name:         "filters thinking blocks without type discriminator",
			input:        `{"messages":[{"role":"assistant","content":[{"thinking":{"text":"internal"placeholderplaceholder,{"type":"text","text":"B"placeholder]placeholder]placeholder`,
			shouldFilter: true,
	placeholder,
		{
			name:         "does not filter tool_use input fields named thinking",
			input:        `{"messages":[{"role":"user","content":[{"type":"tool_use","id":"t1","name":"foo","input":{"thinking":"keepme","x":1placeholderplaceholder,{"type":"text","text":"Hello"placeholder]placeholder]placeholder`,
			shouldFilter: false,
	placeholder,
		{
			name:         "handles empty messages array",
			input:        `{"messages":[]placeholder`,
			shouldFilter: false,
	placeholder,
		{
			name:         "handles missing messages field",
			input:        `{"model":"claude-3"placeholder`,
			shouldFilter: false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterThinkingBlocks([]byte(tt.input))

			if tt.expectError {
				// For invalid JSON, should return original
				require.Equal(t, tt.input, string(result))
				return
		placeholder

			if tt.shouldFilter {
				require.False(t, containsThinkingBlock(result))
		placeholder else {
				// Ensure we don't rewrite JSON when no filtering is needed.
				require.Equal(t, tt.input, string(result))
		placeholder

			// Verify valid JSON returned (unless input was invalid)
			var parsed map[string]any
			err := json.Unmarshal(result, &parsed)
		placeholder
	placeholder)
placeholder
placeholder

func TestFilterThinkingBlocksForRetry_DisablesThinkingAndPreservesAsText(t *testing.T) {
	input := []byte(`{
		"model":"claude-3-5-sonnet-20241022",
		"thinking":{"type":"enabled","budget_tokens":placeholder,
		"messages":[
			{"role":"user","content":[{"type":"text","text":"Hi"placeholder]placeholder,
			{"role":"assistant","content":[
				{"type":"thinking","thinking":"Let me think...","signature":"bad_sig"placeholder,
				{"type":"text","text":"Answer"placeholder
			]placeholder
		]
placeholder`)

	out := FilterThinkingBlocksForRetry(input)

	var req map[string]any
	require.NoError(t, json.Unmarshal(out, &req))
	_, hasThinking := req["thinking"]
	require.False(t, hasThinking)

	msgs, ok := req["messages"].([]any)
	require.True(t, ok)
	require.Len(t, msgs, 2)

	assistant, ok := msgs[1].(map[string]any)
	require.True(t, ok)
	content, ok := assistant["content"].([]any)
	require.True(t, ok)
	require.Len(t, content, 2)

	first, ok := content[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "text", first["type"])
	require.Equal(t, "Let me think...", first["text"])
placeholder

func TestFilterThinkingBlocksForRetry_DisablesThinkingEvenWithoutThinkingBlocks(t *testing.T) {
	input := []byte(`{
		"model":"claude-3-5-sonnet-20241022",
		"thinking":{"type":"enabled","budget_tokens":placeholder,
		"messages":[
			{"role":"user","content":[{"type":"text","text":"Hi"placeholder]placeholder,
			{"role":"assistant","content":[{"type":"text","text":"Prefill"placeholder]placeholder
		]
placeholder`)

	out := FilterThinkingBlocksForRetry(input)

	var req map[string]any
	require.NoError(t, json.Unmarshal(out, &req))
	_, hasThinking := req["thinking"]
	require.False(t, hasThinking)
placeholder

func TestFilterThinkingBlocksForRetry_RemovesRedactedThinkingAndKeepsValidContent(t *testing.T) {
	input := []byte(`{
		"thinking":{"type":"enabled","budget_tokens":placeholder,
		"messages":[
			{"role":"assistant","content":[
				{"type":"redacted_thinking","data":"..."placeholder,
				{"type":"text","text":"Visible"placeholder
			]placeholder
		]
placeholder`)

	out := FilterThinkingBlocksForRetry(input)

	var req map[string]any
	require.NoError(t, json.Unmarshal(out, &req))
	_, hasThinking := req["thinking"]
	require.False(t, hasThinking)

	msgs, ok := req["messages"].([]any)
	require.True(t, ok)
	msg0, ok := msgs[0].(map[string]any)
	require.True(t, ok)
	content, ok := msg0["content"].([]any)
	require.True(t, ok)
	require.Len(t, content, 1)
	content0, ok := content[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "text", content0["type"])
	require.Equal(t, "Visible", content0["text"])
placeholder

func TestFilterThinkingBlocksForRetry_EmptyContentGetsPlaceholder(t *testing.T) {
	input := []byte(`{
		"thinking":{"type":"enabled"placeholder,
		"messages":[
			{"role":"assistant","content":[{"type":"redacted_thinking","data":"..."placeholder]placeholder
		]
placeholder`)

	out := FilterThinkingBlocksForRetry(input)

	var req map[string]any
	require.NoError(t, json.Unmarshal(out, &req))
	msgs, ok := req["messages"].([]any)
	require.True(t, ok)
	msg0, ok := msgs[0].(map[string]any)
	require.True(t, ok)
	content, ok := msg0["content"].([]any)
	require.True(t, ok)
	require.Len(t, content, 1)
	content0, ok := content[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "text", content0["type"])
	require.NotEmpty(t, content0["text"])
placeholder

func TestFilterSignatureSensitiveBlocksForRetry_DowngradesTools(t *testing.T) {
	input := []byte(`{
		"thinking":{"type":"enabled","budget_tokens":placeholder,
		"messages":[
			{"role":"assistant","content":[
				{"type":"tool_use","id":"t1","name":"Bash","input":{"command":"ls"placeholderplaceholder,
				{"type":"tool_result","tool_use_id":"t1","content":"ok","is_error":falseplaceholder
			]placeholder
		]
placeholder`)

	out := FilterSignatureSensitiveBlocksForRetry(input)

	var req map[string]any
	require.NoError(t, json.Unmarshal(out, &req))
	_, hasThinking := req["thinking"]
	require.False(t, hasThinking)

	msgs, ok := req["messages"].([]any)
	require.True(t, ok)
	msg0, ok := msgs[0].(map[string]any)
	require.True(t, ok)
	content, ok := msg0["content"].([]any)
	require.True(t, ok)
	require.Len(t, content, 2)
	content0, ok := content[0].(map[string]any)
	require.True(t, ok)
	content1, ok := content[1].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "text", content0["type"])
	require.Equal(t, "text", content1["type"])
	require.Contains(t, content0["text"], "tool_use")
	require.Contains(t, content1["text"], "tool_result")
placeholder

// ============ Group 6b: context_management.edits 清理测试 ============

// removeThinkingDependentContextStrategies — 边界用例

func TestRemoveThinkingDependentContextStrategies_NoContextManagement(t *testing.T) {
	input := []byte(`{"thinking":{"type":"enabled"placeholder,"messages":[]placeholder`)
	out := removeThinkingDependentContextStrategies(input)
	require.Equal(t, input, out, "无 context_management 字段时应原样返回")
placeholder

func TestRemoveThinkingDependentContextStrategies_EmptyEdits(t *testing.T) {
	input := []byte(`{"context_management":{"edits":[]placeholder,"messages":[]placeholder`)
	out := removeThinkingDependentContextStrategies(input)
	require.Equal(t, input, out, "edits 为空数组时应原样返回")
placeholder

func TestRemoveThinkingDependentContextStrategies_NoClearThinkingEntry(t *testing.T) {
	input := []byte(`{"context_management":{"edits":[{"type":"other_strategy"placeholder]placeholder,"messages":[]placeholder`)
	out := removeThinkingDependentContextStrategies(input)
	require.Equal(t, input, out, "edits 中无 clear_thinking_20251015 时应原样返回")
placeholder

func TestRemoveThinkingDependentContextStrategies_RemovesSingleEntry(t *testing.T) {
	input := []byte(`{"context_management":{"edits":[{"type":"clear_thinking_20251015"placeholder]placeholder,"messages":[]placeholder`)
	out := removeThinkingDependentContextStrategies(input)

	var req map[string]any
	require.NoError(t, json.Unmarshal(out, &req))
	cm, ok := req["context_management"].(map[string]any)
	require.True(t, ok)
	_, hasEdits := cm["edits"]
	require.False(t, hasEdits, "所有 edits 均为 clear_thinking_20251015 时应删除 edits 键")
placeholder

func TestRemoveThinkingDependentContextStrategies_MixedEntries(t *testing.T) {
	input := []byte(`{"context_management":{"edits":[{"type":"clear_thinking_20251015"placeholder,{"type":"other_strategy","param":1placeholder]placeholder,"messages":[]placeholder`)
	out := removeThinkingDependentContextStrategies(input)

	var req map[string]any
	require.NoError(t, json.Unmarshal(out, &req))
	cm, ok := req["context_management"].(map[string]any)
	require.True(t, ok)
	edits, ok := cm["edits"].([]any)
	require.True(t, ok)
	require.Len(t, edits, 1, "仅移除 clear_thinking_20251015，保留其他条目")
	edit0, ok := edits[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "other_strategy", edit0["type"])
placeholder

// FilterThinkingBlocksForRetry — 包含 context_management 的场景

func TestFilterThinkingBlocksForRetry_RemovesClearThinkingStrategy_FastPath(t *testing.T) {
	// 快速路径：messages 中无 thinking 块，仅有顶层 thinking 字段
	// 这条路径曾因提前 return 跳过 removeThinkingDependentContextStrategies 而存在 bug
	input := []byte(`{
		"thinking":{"type":"enabled","budget_tokens":placeholder,
		"context_management":{"edits":[{"type":"clear_thinking_20251015"placeholder]placeholder,
		"messages":[
			{"role":"user","content":[{"type":"text","text":"Hello"placeholder]placeholder
		]
placeholder`)

	out := FilterThinkingBlocksForRetry(input)

	var req map[string]any
	require.NoError(t, json.Unmarshal(out, &req))
	_, hasThinking := req["thinking"]
	require.False(t, hasThinking, "顶层 thinking 应被移除")

	cm, ok := req["context_management"].(map[string]any)
	require.True(t, ok)
	_, hasEdits := cm["edits"]
	require.False(t, hasEdits, "fast path 下 clear_thinking_20251015 应被移除，edits 键应被删除")
placeholder

func TestFilterThinkingBlocksForRetry_RemovesClearThinkingStrategy_WithThinkingBlocks(t *testing.T) {
	// 完整路径：messages 中有 thinking 块（非 fast path）
	input := []byte(`{
		"thinking":{"type":"enabled","budget_tokens":placeholder,
		"context_management":{"edits":[{"type":"clear_thinking_20251015"placeholder,{"type":"keep_this"placeholder]placeholder,
		"messages":[
			{"role":"assistant","content":[
				{"type":"thinking","thinking":"some thought","signature":"sig"placeholder,
				{"type":"text","text":"Answer"placeholder
			]placeholder
		]
placeholder`)

	out := FilterThinkingBlocksForRetry(input)

	var req map[string]any
	require.NoError(t, json.Unmarshal(out, &req))
	_, hasThinking := req["thinking"]
	require.False(t, hasThinking, "顶层 thinking 应被移除")

	cm, ok := req["context_management"].(map[string]any)
	require.True(t, ok)
	edits, ok := cm["edits"].([]any)
	require.True(t, ok)
	require.Len(t, edits, 1, "仅移除 clear_thinking_20251015，保留 keep_this")
	edit0, ok := edits[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "keep_this", edit0["type"])
placeholder

func TestFilterThinkingBlocksForRetry_NoContextManagement_Unaffected(t *testing.T) {
	// 无 context_management 时不应报错，且 thinking 正常被移除
	input := []byte(`{
		"thinking":{"type":"enabled"placeholder,
		"messages":[{"role":"user","content":[{"type":"text","text":"Hi"placeholder]placeholder]
placeholder`)

	out := FilterThinkingBlocksForRetry(input)

	var req map[string]any
	require.NoError(t, json.Unmarshal(out, &req))
	_, hasThinking := req["thinking"]
	require.False(t, hasThinking)
	_, hasCM := req["context_management"]
	require.False(t, hasCM)
placeholder

// FilterSignatureSensitiveBlocksForRetry — 包含 context_management 的场景

func TestFilterSignatureSensitiveBlocksForRetry_RemovesClearThinkingStrategy(t *testing.T) {
	input := []byte(`{
		"thinking":{"type":"enabled","budget_tokens":placeholder,
		"context_management":{"edits":[{"type":"clear_thinking_20251015"placeholder]placeholder,
		"messages":[
			{"role":"assistant","content":[
				{"type":"thinking","thinking":"thought","signature":"sig"placeholder
			]placeholder
		]
placeholder`)

	out := FilterSignatureSensitiveBlocksForRetry(input)

	var req map[string]any
	require.NoError(t, json.Unmarshal(out, &req))
	_, hasThinking := req["thinking"]
	require.False(t, hasThinking, "顶层 thinking 应被移除")

	cm, ok := req["context_management"].(map[string]any)
	require.True(t, ok)
	if rawEdits, hasEdits := cm["edits"]; hasEdits {
		edits, ok := rawEdits.([]any)
		require.True(t, ok)
		for _, e := range edits {
			em, ok := e.(map[string]any)
			require.True(t, ok)
			require.NotEqual(t, "clear_thinking_20251015", em["type"], "clear_thinking_20251015 应被移除")
	placeholder
placeholder
placeholder

func TestFilterSignatureSensitiveBlocksForRetry_PreservesNonThinkingStrategies(t *testing.T) {
	input := []byte(`{
		"thinking":{"type":"enabled"placeholder,
		"context_management":{"edits":[{"type":"clear_thinking_20251015"placeholder,{"type":"other_edit"placeholder]placeholder,
		"messages":[
			{"role":"assistant","content":[
				{"type":"thinking","thinking":"t","signature":"s"placeholder
			]placeholder
		]
placeholder`)

	out := FilterSignatureSensitiveBlocksForRetry(input)

	var req map[string]any
	require.NoError(t, json.Unmarshal(out, &req))

	cm, ok := req["context_management"].(map[string]any)
	require.True(t, ok)
	edits, ok := cm["edits"].([]any)
	require.True(t, ok)
	require.Len(t, edits, 1, "仅移除 clear_thinking_20251015，保留 other_edit")
	edit0, ok := edits[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "other_edit", edit0["type"])
placeholder

func TestFilterSignatureSensitiveBlocksForRetry_NoThinkingField_ContextManagementUntouched(t *testing.T) {
	// 没有顶层 thinking 字段时，context_management 不应被修改
	input := []byte(`{
		"context_management":{"edits":[{"type":"clear_thinking_20251015"placeholder]placeholder,
		"messages":[
			{"role":"assistant","content":[
				{"type":"thinking","thinking":"t","signature":"s"placeholder
			]placeholder
		]
placeholder`)

	out := FilterSignatureSensitiveBlocksForRetry(input)

	var req map[string]any
	require.NoError(t, json.Unmarshal(out, &req))
	cm, ok := req["context_management"].(map[string]any)
	require.True(t, ok)
	edits, ok := cm["edits"].([]any)
	require.True(t, ok)
	require.Len(t, edits, 1, "无顶层 thinking 时 context_management 不应被修改")
placeholder

// ============ Group 7: ParseGatewayRequest 补充单元测试 ============

// Task 7.1 — 类型校验边界测试
func TestParseGatewayRequest_TypeValidation(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		wantErr   bool
		errSubstr string // 期望的错误信息子串（为空则不检查）
placeholder{
		{
			name:      "model 为 int",
			body:      `{"model":placeholder`,
			wantErr:   true,
			errSubstr: "invalid model field type",
	placeholder,
		{
			name:      "model 为 array",
			body:      `{"model":[]placeholder`,
			wantErr:   true,
			errSubstr: "invalid model field type",
	placeholder,
		{
			name:      "model 为 bool",
			body:      `{"model":trueplaceholder`,
			wantErr:   true,
			errSubstr: "invalid model field type",
	placeholder,
		{
			name:      "model 为 null — gjson Null 类型触发类型校验错误",
			body:      `{"model":nullplaceholder`,
			wantErr:   true, // gjson: Exists()=true, Type=Null != String → 返回错误
			errSubstr: "invalid model field type",
	placeholder,
		{
			name:      "stream 为 string",
			body:      `{"stream":"true"placeholder`,
			wantErr:   true,
			errSubstr: "invalid stream field type",
	placeholder,
		{
			name:      "stream 为 int",
			body:      `{"stream":1placeholder`,
			wantErr:   true,
			errSubstr: "invalid stream field type",
	placeholder,
		{
			name:      "stream 为 null — gjson Null 类型触发类型校验错误",
			body:      `{"stream":nullplaceholder`,
			wantErr:   true, // gjson: Exists()=true, Type=Null != True && != False → 返回错误
			errSubstr: "invalid stream field type",
	placeholder,
		{
			name:      "model 为 object",
			body:      `{"model":{placeholderplaceholder`,
			wantErr:   true,
			errSubstr: "invalid model field type",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseGatewayRequest([]byte(tt.body), "")
			if tt.wantErr {
			placeholder
				if tt.errSubstr != "" {
					require.Contains(t, err.Error(), tt.errSubstr)
			placeholder
		placeholder else {
			placeholder
		placeholder
	placeholder)
placeholder
placeholder

// Task 7.2 — 可选字段缺失测试
func TestParseGatewayRequest_OptionalFieldsMissing(t *testing.T) {
	tests := []struct {
		name            string
		body            string
		wantModel       string
		wantStream      bool
		wantMetadataUID string
		wantHasSystem   bool
		wantThinking    bool
		wantMaxTokens   int
		wantMessagesNil bool
		wantMessagesLen int
placeholder{
		{
			name:            "完全空 JSON — 所有字段零值",
			body:            `{placeholder`,
			wantModel:       "",
			wantStream:      false,
			wantMetadataUID: "",
			wantHasSystem:   false,
			wantThinking:    false,
			wantMaxTokens:   0,
			wantMessagesNil: true,
	placeholder,
		{
			name:            "metadata 无 user_id",
			body:            `{"model":"test"placeholder`,
			wantModel:       "test",
			wantMetadataUID: "",
			wantHasSystem:   false,
			wantThinking:    false,
	placeholder,
		{
			name:         "thinking 非 enabled（type=disabled）",
			body:         `{"model":"test","thinking":{"type":"disabled"placeholderplaceholder`,
			wantModel:    "test",
			wantThinking: false,
	placeholder,
		{
			name:         "thinking 字段缺失",
			body:         `{"model":"test"placeholder`,
			wantModel:    "test",
			wantThinking: false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := ParseGatewayRequest([]byte(tt.body), "")
		placeholder

			require.Equal(t, tt.wantModel, parsed.Model)
			require.Equal(t, tt.wantStream, parsed.Stream)
			require.Equal(t, tt.wantMetadataUID, parsed.MetadataUserID)
			require.Equal(t, tt.wantHasSystem, parsed.HasSystem)
			require.Equal(t, tt.wantThinking, parsed.ThinkingEnabled)
			require.Equal(t, tt.wantMaxTokens, parsed.MaxTokens)

			if tt.wantMessagesNil {
				require.Nil(t, parsed.Messages)
		placeholder
			if tt.wantMessagesLen > 0 {
				require.Len(t, parsed.Messages, tt.wantMessagesLen)
		placeholder
	placeholder)
placeholder
placeholder

// Task 7.3 — Gemini 协议分支测试
// 已有测试覆盖：
// - TestParseGatewayRequest_GeminiSystemInstruction: 正常 systemInstruction+contents
// - TestParseGatewayRequest_GeminiNoContents: 缺失 contents
// - TestParseGatewayRequest_GeminiContents: 正常 contents（无 systemInstruction）
// 因此跳过。

// Task 7.4 — max_tokens 边界测试
func TestParseGatewayRequest_MaxTokensBoundary(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		wantMaxTokens int
		wantErr       bool
placeholder{
		{
			name:          "正常整数",
			body:          `{"max_tokens":placeholder`,
			wantMaxTokens: 1024,
	placeholder,
		{
			name:          "浮点数（非整数）被忽略",
			body:          `{"max_tokens":placeholder`,
			wantMaxTokens: 0,
	placeholder,
		{
			name:          "负整数可以通过",
			body:          `{"max_tokens":-1placeholder`,
			wantMaxTokens: -1,
	placeholder,
		{
			name:          "超大值不 panic",
			body:          `{"max_tokens":9999999999999999placeholder`,
			wantMaxTokens: 10000000000000000, // float64 精度导致 9999999999999999 → 1e16
	placeholder,
		{
			name:          "null 值被忽略",
			body:          `{"max_tokens":nullplaceholder`,
			wantMaxTokens: 0, // gjson Type=Null != Number → 条件不满足，跳过
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := ParseGatewayRequest([]byte(tt.body), "")
			if tt.wantErr {
			placeholder
				return
		placeholder
		placeholder
			require.Equal(t, tt.wantMaxTokens, parsed.MaxTokens)
	placeholder)
placeholder
placeholder

// ============ Task 7.5: Benchmark 测试 ============

// parseGatewayRequestOld 是基于完整 json.Unmarshal 的旧实现，用于 benchmark 对比基线。
// 核心路径：先 Unmarshal 到 map[string]any，再逐字段提取。
func parseGatewayRequestOld(body []byte, protocol string) (*ParsedRequest, error) {
	parsed := &ParsedRequest{
		Body: body,
placeholder

	var req map[string]any
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
placeholder

	// model
	if raw, ok := req["model"]; ok {
		s, ok := raw.(string)
		if !ok {
			return nil, fmt.Errorf("invalid model field type")
	placeholder
		parsed.Model = s
placeholder

	// stream
	if raw, ok := req["stream"]; ok {
		b, ok := raw.(bool)
		if !ok {
			return nil, fmt.Errorf("invalid stream field type")
	placeholder
		parsed.Stream = b
placeholder

	// metadata.user_id
	if meta, ok := req["metadata"].(map[string]any); ok {
		if uid, ok := meta["user_id"].(string); ok {
			parsed.MetadataUserID = uid
	placeholder
placeholder

	// thinking.type
	if thinking, ok := req["thinking"].(map[string]any); ok {
		if thinkType, ok := thinking["type"].(string); ok && thinkType == "enabled" {
			parsed.ThinkingEnabled = true
	placeholder
placeholder

	// max_tokens
	if raw, ok := req["max_tokens"]; ok {
		if n, ok := parseIntegralNumber(raw); ok {
			parsed.MaxTokens = n
	placeholder
placeholder

	// system / messages（按协议分支）
	switch protocol {
	case domain.PlatformGemini:
		if sysInst, ok := req["systemInstruction"].(map[string]any); ok {
			if parts, ok := sysInst["parts"].([]any); ok {
				parsed.System = parts
		placeholder
	placeholder
		if contents, ok := req["contents"].([]any); ok {
			parsed.Messages = contents
	placeholder
	default:
		if system, ok := req["system"]; ok {
			parsed.HasSystem = true
			parsed.System = system
	placeholder
		if messages, ok := req["messages"].([]any); ok {
			parsed.Messages = messages
	placeholder
placeholder

	return parsed, nil
placeholder

// buildSmallJSON 构建 ~500B 的小型测试 JSON
func buildSmallJSON() []byte {
	return []byte(`{"model":"claude-sonnet-4-5","stream":true,"max_tokens":4096,"metadata":{"user_id":"user-abc123"placeholder,"thinking":{"type":"enabled","budget_tokens":placeholder,"system":"You are a helpful assistant.","messages":[{"role":"user","content":"What is the meaning of life?"placeholder,{"role":"assistant","content":"The meaning of life is a philosophical question."placeholder,{"role":"user","content":"Can you elaborate?"placeholder]placeholder`)
placeholder

// buildLargeJSON 构建 ~50KB 的大型测试 JSON（大量 messages）
func buildLargeJSON() []byte {
	var b strings.Builder
	b.WriteString(`{"model":"claude-sonnet-4-5","stream":true,"max_tokens":8192,"metadata":{"user_id":"user-xyz789"placeholder,"system":[{"type":"text","text":"You are a detailed assistant.","cache_control":{"type":"ephemeral"placeholderplaceholder],"messages":[`)

	msgCount := 200
	for i := 0; i < msgCount; i++ {
		if i > 0 {
			b.WriteByte(',')
	placeholder
		if i%2 == 0 {
			b.WriteString(fmt.Sprintf(`{"role":"user","content":"This is user message number %d with some extra padding text to make the message reasonably long for benchmarking purposes. Lorem ipsum dolor sit amet."placeholder`, i))
	placeholder else {
			b.WriteString(fmt.Sprintf(`{"role":"assistant","content":[{"type":"text","text":"This is assistant response number %d. I will provide a detailed answer with multiple sentences to simulate real conversation content for benchmark testing."placeholder]placeholder`, i))
	placeholder
placeholder

	b.WriteString(`]placeholder`)
	return []byte(b.String())
placeholder

func BenchmarkParseGatewayRequest_Old_Small(b *testing.B) {
	data := buildSmallJSON()
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parseGatewayRequestOld(data, "")
placeholder
placeholder

func BenchmarkParseGatewayRequest_New_Small(b *testing.B) {
	data := buildSmallJSON()
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseGatewayRequest(data, "")
placeholder
placeholder

func BenchmarkParseGatewayRequest_Old_Large(b *testing.B) {
	data := buildLargeJSON()
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parseGatewayRequestOld(data, "")
placeholder
placeholder

func TestParseGatewayRequest_OutputEffort(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantEffort string
placeholder{
		{
			name:       "output_config.effort present",
			body:       `{"model":"claude-opus-4-6","output_config":{"effort":"medium"placeholder,"messages":[]placeholder`,
			wantEffort: "medium",
	placeholder,
		{
			name:       "output_config.effort max",
			body:       `{"model":"claude-opus-4-6","output_config":{"effort":"max"placeholder,"messages":[]placeholder`,
			wantEffort: "max",
	placeholder,
		{
			name:       "output_config without effort",
			body:       `{"model":"claude-opus-4-6","output_config":{placeholder,"messages":[]placeholder`,
			wantEffort: "",
	placeholder,
		{
			name:       "no output_config",
			body:       `{"model":"claude-opus-4-6","messages":[]placeholder`,
			wantEffort: "",
	placeholder,
		{
			name:       "effort with whitespace trimmed",
			body:       `{"model":"claude-opus-4-6","output_config":{"effort":" high "placeholder,"messages":[]placeholder`,
			wantEffort: "high",
	placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := ParseGatewayRequest([]byte(tt.body), "")
		placeholder
			require.Equal(t, tt.wantEffort, parsed.OutputEffort)
	placeholder)
placeholder
placeholder

func TestNormalizeClaudeOutputEffort(t *testing.T) {
	tests := []struct {
		input string
		want  *string
placeholder{
		{"low", strPtr("low")placeholder,
		{"medium", strPtr("medium")placeholder,
		{"high", strPtr("high")placeholder,
		{"max", strPtr("max")placeholder,
		{"LOW", strPtr("low")placeholder,
		{"Max", strPtr("max")placeholder,
		{" medium ", strPtr("medium")placeholder,
		{"", nilplaceholder,
		{"unknown", nilplaceholder,
		{"xhigh", nilplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := NormalizeClaudeOutputEffort(tt.input)
			if tt.want == nil {
				require.Nil(t, got)
		placeholder else {
				require.NotNil(t, got)
				require.Equal(t, *tt.want, *got)
		placeholder
	placeholder)
placeholder
placeholder

func BenchmarkParseGatewayRequest_New_Large(b *testing.B) {
	data := buildLargeJSON()
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseGatewayRequest(data, "")
placeholder
placeholder
