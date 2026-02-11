package service

import (
	"encoding/json"
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
