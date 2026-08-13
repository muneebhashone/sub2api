package apicompat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChatCompletionsToResponsesPreservesXSearchTool(t *testing.T) {
	enabled := true
	req := &ChatCompletionsRequest{
		Model: "grok-4.5",
		Messages: []ChatMessage{
			{Role: "user", Content: json.RawMessage(`"latest xAI post"`)placeholder,
	placeholder,
		Tools: []ChatTool{{
			Type:                     "x_search",
			AllowedXHandles:          []string{"xai"placeholder,
			ExcludedXHandles:         []string{"spam"placeholder,
			FromDate:                 "2026-08-01",
			ToDate:                   "2026-08-10",
			EnableImageUnderstanding: &enabled,
			EnableVideoUnderstanding: &enabled,
placeholder
		ToolChoice: json.RawMessage(`{"type":"x_search"placeholder`),
placeholder

	resp, err := ChatCompletionsToResponses(req)
placeholder
	require.Len(t, resp.Tools, 1)
	require.Equal(t, "x_search", resp.Tools[0].Type)
	require.Equal(t, []string{"xai"placeholder, resp.Tools[0].AllowedXHandles)
	require.Equal(t, []string{"spam"placeholder, resp.Tools[0].ExcludedXHandles)
	require.Equal(t, "2026-08-01", resp.Tools[0].FromDate)
	require.Equal(t, "2026-08-10", resp.Tools[0].ToDate)
	require.NotNil(t, resp.Tools[0].EnableImageUnderstanding)
	require.True(t, *resp.Tools[0].EnableImageUnderstanding)
	require.NotNil(t, resp.Tools[0].EnableVideoUnderstanding)
	require.True(t, *resp.Tools[0].EnableVideoUnderstanding)
	require.JSONEq(t, `{"type":"x_search"placeholder`, string(resp.ToolChoice))
placeholder

func TestResponsesToChatCompletionsPreservesXSearchTool(t *testing.T) {
	enabled := true
	req := &ResponsesRequest{
		Model: "grok-4.5",
		Input: json.RawMessage(`"latest xAI post"`),
		Tools: []ResponsesTool{{
			Type:                     "x_search",
			AllowedXHandles:          []string{"xai"placeholder,
			ExcludedXHandles:         []string{"spam"placeholder,
			FromDate:                 "2026-08-01",
			ToDate:                   "2026-08-10",
			EnableImageUnderstanding: &enabled,
			EnableVideoUnderstanding: &enabled,
placeholder
		ToolChoice: json.RawMessage(`{"type":"x_search"placeholder`),
placeholder

	chat, err := ResponsesToChatCompletionsRequest(req)
placeholder
	require.Len(t, chat.Tools, 1)
	require.Equal(t, "x_search", chat.Tools[0].Type)
	require.Equal(t, []string{"xai"placeholder, chat.Tools[0].AllowedXHandles)
	require.Equal(t, []string{"spam"placeholder, chat.Tools[0].ExcludedXHandles)
	require.JSONEq(t, `{"type":"x_search"placeholder`, string(chat.ToolChoice))
placeholder

func TestResponsesToChatCompletionsXSearchToolChoiceString(t *testing.T) {
	chat, err := ResponsesToChatCompletionsRequest(&ResponsesRequest{
		Model:      "grok-4.5",
		Input:      json.RawMessage(`"latest xAI post"`),
		Tools:      []ResponsesTool{{Type: "x_search"placeholderplaceholder,
		ToolChoice: json.RawMessage(`"x_search"`),
placeholder)
placeholder
	require.JSONEq(t, `"x_search"`, string(chat.ToolChoice))
placeholder
