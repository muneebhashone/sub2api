package service

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseGatewayRequest(t *testing.T) {
	body := []byte(`{"model":"claude-3-7-sonnet","stream":true,"metadata":{"user_id":"session_123e4567-e89b-12d3-a456-426614174000"placeholder,"system":[{"type":"text","text":"hello","cache_control":{"type":"ephemeral"placeholderplaceholder],"messages":[{"content":"hi"placeholder]placeholder`)
	parsed, err := ParseGatewayRequest(body)
placeholder
	require.Equal(t, "claude-3-7-sonnet", parsed.Model)
	require.True(t, parsed.Stream)
	require.Equal(t, "session_123e4567-e89b-12d3-a456-426614174000", parsed.MetadataUserID)
	require.True(t, parsed.HasSystem)
	require.NotNil(t, parsed.System)
	require.Len(t, parsed.Messages, 1)
placeholder

func TestParseGatewayRequest_SystemNull(t *testing.T) {
	body := []byte(`{"model":"claude-3","system":nullplaceholder`)
	parsed, err := ParseGatewayRequest(body)
placeholder
	// 显式传入 system:null 也应视为“字段已存在”，避免默认 system 被注入。
	require.True(t, parsed.HasSystem)
	require.Nil(t, parsed.System)
placeholder

func TestParseGatewayRequest_InvalidModelType(t *testing.T) {
	body := []byte(`{"model":placeholder`)
	_, err := ParseGatewayRequest(body)
placeholder
placeholder

func TestParseGatewayRequest_InvalidStreamType(t *testing.T) {
	body := []byte(`{"stream":"true"placeholder`)
	_, err := ParseGatewayRequest(body)
placeholder
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
