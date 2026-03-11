package service

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertChatCompletionsToResponses(t *testing.T) {
	req := map[string]any{
		"model": "gpt-4o",
		"messages": []any{
			map[string]any{
				"role":    "user",
				"content": "hello",
		placeholder,
			map[string]any{
				"role": "assistant",
				"tool_calls": []any{
					map[string]any{
						"id":   "call_1",
						"type": "function",
						"function": map[string]any{
							"name":      "ping",
							"arguments": "{placeholder",
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			map[string]any{
				"role":          "tool",
				"tool_call_id":  "call_1",
				"content":       "ok",
				"response":      "ignored",
				"response_time": 1,
		placeholder,
	placeholder,
		"functions": []any{
			map[string]any{
				"name":        "ping",
				"description": "ping tool",
				"parameters":  map[string]any{"type": "object"placeholder,
		placeholder,
	placeholder,
		"function_call": map[string]any{"name": "ping"placeholder,
placeholder

	converted, err := ConvertChatCompletionsToResponses(req)
placeholder
	require.Equal(t, "gpt-4o", converted["model"])

	input, ok := converted["input"].([]any)
	require.True(t, ok)
	require.Len(t, input, 3)

	toolCall := findInputItemByType(input, "tool_call")
	require.NotNil(t, toolCall)
	require.Equal(t, "call_1", toolCall["call_id"])

	toolOutput := findInputItemByType(input, "function_call_output")
	require.NotNil(t, toolOutput)
	require.Equal(t, "call_1", toolOutput["call_id"])

	tools, ok := converted["tools"].([]any)
	require.True(t, ok)
	require.Len(t, tools, 1)

	require.Equal(t, map[string]any{"name": "ping"placeholder, converted["tool_choice"])
placeholder

func TestConvertResponsesToChatCompletion(t *testing.T) {
	resp := map[string]any{
		"id":         "resp_123",
		"model":      "gpt-4o",
		"created_at": 1700000000,
		"output": []any{
			map[string]any{
				"type": "message",
				"role": "assistant",
				"content": []any{
					map[string]any{
						"type": "output_text",
						"text": "hi",
				placeholder,
			placeholder,
		placeholder,
	placeholder,
		"usage": map[string]any{
			"input_tokens":  2,
			"output_tokens": 3,
	placeholder,
placeholder
	body, err := json.Marshal(resp)
placeholder

	converted, err := ConvertResponsesToChatCompletion(body)
placeholder

	var chat map[string]any
	require.NoError(t, json.Unmarshal(converted, &chat))
	require.Equal(t, "chat.completion", chat["object"])

	choices, ok := chat["choices"].([]any)
	require.True(t, ok)
	require.Len(t, choices, 1)

	choice, ok := choices[0].(map[string]any)
	require.True(t, ok)
	message, ok := choice["message"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "hi", message["content"])

	usage, ok := chat["usage"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, float64(2), usage["prompt_tokens"])
	require.Equal(t, float64(3), usage["completion_tokens"])
	require.Equal(t, float64(5), usage["total_tokens"])
placeholder

func findInputItemByType(items []any, itemType string) map[string]any {
	for _, item := range items {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
	placeholder
		if itemMap["type"] == itemType {
			return itemMap
	placeholder
placeholder
	return nil
placeholder
