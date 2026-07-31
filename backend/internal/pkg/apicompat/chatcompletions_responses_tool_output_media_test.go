package apicompat

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testToolImageDataURL = "data:image/png;base64,AQID"
	testToolImageRemote  = "https://example.com/tool-output.png"
)

func TestResponsesToolOutputMedia_ExtractsSupportedShapes(t *testing.T) {
	tests := []struct {
		name       string
		call       string
		outputType string
		output     string
		imageURL   string
		toolText   string
placeholder{
		{
			name:       "image-only array",
			call:       `{"type":"function_call","call_id":"call_image","name":"view_image","arguments":"{placeholder"placeholder`,
			outputType: "function_call_output",
			output:     `[{"type":"input_image","image_url":"data:image/png;base64,AQID"placeholder]`,
			imageURL:   testToolImageDataURL,
	placeholder,
		{
			name:       "text and nested image URL",
			call:       `{"type":"function_call","call_id":"call_image","name":"view_image","arguments":"{placeholder"placeholder`,
			outputType: "function_call_output",
			output:     `[{"type":"input_text","text":"render complete"placeholder,{"type":"image_url","image_url":{"url":"https://example.com/tool-output.png"placeholderplaceholder]`,
			imageURL:   testToolImageRemote,
			toolText:   "render complete",
	placeholder,
		{
			name:       "top-level image object",
			call:       `{"type":"function_call","call_id":"call_image","name":"view_image","arguments":"{placeholder"placeholder`,
			outputType: "function_call_output",
			output:     `{"type":"input_image","image_url":"data:image/png;base64,AQID"placeholder`,
			imageURL:   testToolImageDataURL,
	placeholder,
		{
			name:       "content wrapper",
			call:       `{"type":"function_call","call_id":"call_image","name":"view_image","arguments":"{placeholder"placeholder`,
			outputType: "function_call_output",
			output:     `{"status":"ok","content":[{"type":"input_image","image_url":"data:image/png;base64,AQID"placeholder],"unknown":{"score":0.9,"large":9007199254740993placeholderplaceholder`,
			imageURL:   testToolImageDataURL,
			toolText:   `"large":9007199254740993`,
	placeholder,
		{
			name:       "JSON string output",
			call:       `{"type":"function_call","call_id":"call_image","name":"view_image","arguments":"{placeholder"placeholder`,
			outputType: "function_call_output",
			output:     `"[{\"type\":\"input_image\",\"image_url\":\"data:image/png;base64,AQID\"placeholder]"`,
			imageURL:   testToolImageDataURL,
	placeholder,
		{
			name:       "bare data URL",
			call:       `{"type":"function_call","call_id":"call_image","name":"view_image","arguments":"{placeholder"placeholder`,
			outputType: "function_call_output",
			output:     `"data:image/png;base64,AQID"`,
			imageURL:   testToolImageDataURL,
	placeholder,
		{
			name:       "custom tool output",
			call:       `{"type":"custom_tool_call","call_id":"call_image","name":"view_image","input":"{placeholder"placeholder`,
			outputType: "custom_tool_call_output",
			output:     `[{"type":"input_image","image_url":"data:image/png;base64,AQID"placeholder]`,
			imageURL:   testToolImageDataURL,
	placeholder,
		{
			name:       "tool search output",
			call:       `{"type":"tool_search_call","call_id":"call_image","arguments":{"query":"image"placeholderplaceholder`,
			outputType: "tool_search_output",
			output:     `[{"type":"image_url","image_url":{"url":"https://example.com/tool-output.png"placeholderplaceholder]`,
			imageURL:   testToolImageRemote,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := fmt.Sprintf(`[%s,{"type":%q,"call_id":"call_image","output":%splaceholder]`, tt.call, tt.outputType, tt.output)
			messages := convertToolOutputMedia(t, input)

			require.Len(t, messages, 3)
			require.Equal(t, []string{"assistant", "tool", "user"placeholder, chatMessageRoles(messages))
			require.Equal(t, "call_image", messages[1].ToolCallID)

			toolText := chatToolContentString(t, messages[1])
			require.Contains(t, toolText, "[Tool output media moved to the following user message]")
			require.NotContains(t, toolText, tt.imageURL)
			if tt.toolText != "" {
				require.Contains(t, toolText, tt.toolText)
		placeholder

			parts := chatContentParts(t, messages[2])
			require.Len(t, parts, 2)
			require.Equal(t, "text", parts[0].Type)
			require.Equal(t, "[Tool output media for call call_image]", parts[0].Text)
			require.Equal(t, "image_url", parts[1].Type)
			require.NotNil(t, parts[1].ImageURL)
			require.Equal(t, tt.imageURL, parts[1].ImageURL.URL)
	placeholder)
placeholder
placeholder

func TestResponsesToolOutputMedia_ParallelBatchUsesCallOrder(t *testing.T) {
	messages := convertToolOutputMedia(t, `[
		{"type":"function_call","call_id":"call_A","name":"view_image","arguments":"{placeholder"placeholder,
		{"type":"function_call","call_id":"call_B","name":"view_image","arguments":"{placeholder"placeholder,
		{"type":"function_call_output","call_id":"call_B","output":[{"type":"input_image","image_url":{"url":"https://example.com/b.png"placeholderplaceholder]placeholder,
		{"type":"function_call_output","call_id":"call_A","output":[{"type":"input_image","image_url":{"url":"https://example.com/a.png"placeholderplaceholder]placeholder
	]`)

	require.Len(t, messages, 4)
	require.Equal(t, []string{"assistant", "tool", "tool", "user"placeholder, chatMessageRoles(messages))
	require.Equal(t, []string{"call_A", "call_B"placeholder, []string{messages[1].ToolCallID, messages[2].ToolCallIDplaceholder)

	parts := chatContentParts(t, messages[3])
	require.Len(t, parts, 4)
	require.Equal(t, "[Tool output media for call call_A]", parts[0].Text)
	require.Equal(t, "https://example.com/a.png", parts[1].ImageURL.URL)
	require.Equal(t, "[Tool output media for call call_B]", parts[2].Text)
	require.Equal(t, "https://example.com/b.png", parts[3].ImageURL.URL)
placeholder

func TestResponsesToolOutputMedia_PreservesRichSiblingWhenRewriting(t *testing.T) {
	messages := convertToolOutputMedia(t, `[
		{"type":"function_call","call_id":"call_image","name":"view_image","arguments":"{placeholder"placeholder,
		{"type":"function_call_output","call_id":"call_image","output":[
			{"type":"result","url":"https://example.com/result","score":0.9,"text":"complete","extra":{"count":2placeholderplaceholder,
			{"type":"input_image","image_url":"data:image/png;base64,AQID"placeholder
		]placeholder
	]`)

	toolText := chatToolContentString(t, messages[1])
	require.JSONEq(t, `[
		{"type":"result","url":"https://example.com/result","score":0.9,"text":"complete","extra":{"count":2placeholderplaceholder,
		{"type":"input_text","text":"[Tool output media moved to the following user message]"placeholder
	]`, toolText)
placeholder

func TestResponsesToolOutputMedia_InterleavedMessagesFollowMediaBatch(t *testing.T) {
	messages := convertToolOutputMedia(t, `[
		{"type":"function_call","call_id":"call_A","name":"view_image","arguments":"{placeholder"placeholder,
		{"type":"message","role":"developer","content":[{"type":"input_text","text":"approval saved"placeholder]placeholder,
		{"type":"message","role":"user","content":[{"type":"input_text","text":"continue"placeholder]placeholder,
		{"type":"function_call_output","call_id":"call_A","output":[{"type":"input_image","image_url":"data:image/png;base64,AQID"placeholder]placeholder
	]`)

	require.Equal(t, []string{"assistant", "tool", "user", "system", "user"placeholder, chatMessageRoles(messages))
	require.Equal(t, "[Tool output media for call call_A]", chatContentParts(t, messages[2])[0].Text)
	require.JSONEq(t, `"approval saved"`, string(messages[3].Content))
	require.JSONEq(t, `"continue"`, string(messages[4].Content))
placeholder

func TestResponsesToolOutputMedia_DropsOrphanAndUnansweredCallMedia(t *testing.T) {
	t.Run("orphan", func(t *testing.T) {
		messages := convertToolOutputMedia(t, `[
			{"type":"function_call_output","call_id":"call_ghost","output":[{"type":"input_image","image_url":"data:image/png;base64,AQID"placeholder]placeholder
		]`)
		require.Empty(t, messages)
placeholder)

	t.Run("unanswered parallel call", func(t *testing.T) {
		messages := convertToolOutputMedia(t, `[
			{"type":"function_call","call_id":"call_A","name":"view_image","arguments":"{placeholder"placeholder,
			{"type":"function_call","call_id":"call_B","name":"view_image","arguments":"{placeholder"placeholder,
			{"type":"function_call_output","call_id":"call_A","output":[{"type":"input_image","image_url":"data:image/png;base64,AQID"placeholder]placeholder
		]`)

		require.Len(t, messages, 3)
		require.Len(t, messages[0].ToolCalls, 1)
		require.Equal(t, "call_A", messages[0].ToolCalls[0].ID)
		parts := chatContentParts(t, messages[2])
		require.Len(t, parts, 2)
		require.NotContains(t, string(messages[2].Content), "call_B")
placeholder)
placeholder

func TestResponsesToolOutputMedia_PreservesMediaFreeOutputBytes(t *testing.T) {
	tests := []struct {
		name   string
		output string
placeholder{
		{
			name:   "rich unknown object",
			output: `{"type":"result","url":"https://example.com/result","score":0.9,"text":"complete","extra":{"count":2placeholderplaceholder`,
	placeholder,
		{
			name:   "no-image array",
			output: `[ { "type": "input_text", "text": "ok" placeholder, {"unknown":trueplaceholder ]`,
	placeholder,
		{
			name:   "plain string",
			output: `"plain output"`,
	placeholder,
		{
			name:   "JSON string without image",
			output: `"{ \"ok\": true placeholder"`,
	placeholder,
		{
			name:   "embedded data URL text",
			output: `"prefix data:image/png;base64,AQID suffix"`,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := fmt.Sprintf(`[
				{"type":"function_call","call_id":"call_text","name":"exec","arguments":"{placeholder"placeholder,
				{"type":"function_call_output","call_id":"call_text","output":%splaceholder
			]`, tt.output)
			messages := convertToolOutputMedia(t, input)
			require.Len(t, messages, 2)

			var expected string
			if err := json.Unmarshal([]byte(tt.output), &expected); err != nil {
				expected = tt.output
		placeholder
			expectedContent, err := json.Marshal(expected)
		placeholder
			require.Equal(t, string(expectedContent), string(messages[1].Content))
	placeholder)
placeholder
placeholder

func TestResponsesToolOutputMedia_DuplicateCallIDIsLastWins(t *testing.T) {
	t.Run("later media replaces earlier media", func(t *testing.T) {
		messages := convertToolOutputMedia(t, `[
			{"type":"function_call","call_id":"call_image","name":"view_image","arguments":"{placeholder"placeholder,
			{"type":"function_call_output","call_id":"call_image","output":[{"type":"input_image","image_url":{"url":"https://example.com/first.png"placeholderplaceholder]placeholder,
			{"type":"function_call_output","call_id":"call_image","output":[{"type":"input_image","image_url":{"url":"https://example.com/last.png"placeholderplaceholder]placeholder
		]`)

		require.Len(t, messages, 3)
		require.NotContains(t, string(messages[1].Content), "first.png")
		require.NotContains(t, string(messages[2].Content), "first.png")
		require.Contains(t, string(messages[2].Content), "last.png")
placeholder)

	t.Run("later text clears earlier media", func(t *testing.T) {
		messages := convertToolOutputMedia(t, `[
			{"type":"function_call","call_id":"call_image","name":"view_image","arguments":"{placeholder"placeholder,
			{"type":"function_call_output","call_id":"call_image","output":[{"type":"input_image","image_url":"data:image/png;base64,AQID"placeholder]placeholder,
			{"type":"function_call_output","call_id":"call_image","output":"latest text"placeholder
		]`)

		require.Len(t, messages, 2)
		require.Equal(t, "latest text", chatToolContentString(t, messages[1]))
placeholder)
placeholder

func TestResponsesToChatCompletionsRequest_ToolContentNeverContainsExtractedMedia(t *testing.T) {
	req := &ResponsesRequest{
		Model: "vision-model",
		Input: json.RawMessage(`[
			{"type":"function_call","call_id":"call_function","name":"view_image","arguments":"{placeholder"placeholder,
			{"type":"custom_tool_call","call_id":"call_custom","name":"custom_image","input":"{placeholder"placeholder,
			{"type":"tool_search_call","call_id":"call_search","arguments":{"query":"image"placeholderplaceholder,
			{"type":"function_call_output","call_id":"call_function","output":[{"type":"input_image","image_url":"data:image/png;base64,AQID"placeholder]placeholder,
			{"type":"custom_tool_call_output","call_id":"call_custom","output":{"content":[{"type":"image_url","image_url":{"url":"https://example.com/custom.png"placeholderplaceholder]placeholderplaceholder,
			{"type":"tool_search_output","call_id":"call_search","output":"data:image/jpeg;base64,BAUG"placeholder
		]`),
placeholder

	out, err := ResponsesToChatCompletionsRequest(req)
placeholder
	assertChatInvariants(t, out.Messages)

	var toolCount int
	for _, message := range out.Messages {
		if message.Role != "tool" {
			continue
	placeholder
		toolCount++
		content := string(message.Content)
		require.NotContains(t, content, "data:image/")
		require.NotContains(t, content, "https://example.com/custom.png")
placeholder
	require.Equal(t, 3, toolCount)
	require.Equal(t, []string{"assistant", "tool", "tool", "tool", "user"placeholder, chatMessageRoles(out.Messages))
placeholder

func convertToolOutputMedia(t *testing.T, input string) []ChatMessage {
placeholder
	messages, err := responsesInputToChatMessages("", json.RawMessage(input))
placeholder
	assertChatInvariants(t, messages)
	return messages
placeholder

func chatToolContentString(t *testing.T, message ChatMessage) string {
placeholder
	require.Equal(t, "tool", message.Role)
	var content string
	require.NoError(t, json.Unmarshal(message.Content, &content))
	return content
placeholder

func chatContentParts(t *testing.T, message ChatMessage) []ChatContentPart {
placeholder
	var parts []ChatContentPart
	require.NoError(t, json.Unmarshal(message.Content, &parts))
	for _, part := range parts {
		if part.Type == "image_url" {
			require.NotNil(t, part.ImageURL)
			require.False(t, strings.TrimSpace(part.ImageURL.URL) == "")
	placeholder
placeholder
	return parts
placeholder
