package apicompat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// ChatCompletionsToResponses tests
// ---------------------------------------------------------------------------

func TestChatCompletionsToResponses_BasicText(t *testing.T) {
	req := &ChatCompletionsRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{Role: "user", Content: json.RawMessage(`"Hello"`)placeholder,
	placeholder,
placeholder

	resp, err := ChatCompletionsToResponses(req)
placeholder
	assert.Equal(t, "gpt-4o", resp.Model)
	assert.True(t, resp.Stream) // always forced true
	assert.False(t, *resp.Store)

	var items []ResponsesInputItem
	require.NoError(t, json.Unmarshal(resp.Input, &items))
	require.Len(t, items, 1)
	assert.Equal(t, "user", items[0].Role)
placeholder

func TestUsageConversionsPreserveCacheWriteTokens(t *testing.T) {
	var responsesUsage ResponsesUsage
	require.NoError(t, json.Unmarshal([]byte(`{
		"input_tokens":1000,
		"output_tokens":50,
		"input_tokens_details":{"cached_tokens":100,"cache_write_tokens":200placeholder
placeholder`), &responsesUsage))
	require.NotNil(t, responsesUsage.InputTokensDetails)
	require.Equal(t, 200, responsesUsage.InputTokensDetails.CacheWriteTokens)

	chatUsage := chatUsageFromResponsesUsage(&responsesUsage)
	require.NotNil(t, chatUsage.PromptTokensDetails)
	require.Equal(t, 100, chatUsage.PromptTokensDetails.CachedTokens)
	require.Equal(t, 200, chatUsage.PromptTokensDetails.CacheWriteTokens)

	roundTrip := ChatUsageToResponsesUsage(chatUsage)
	require.NotNil(t, roundTrip.InputTokensDetails)
	require.Equal(t, 200, roundTrip.CacheCreationInputTokens)
	require.Equal(t, 200, roundTrip.InputTokensDetails.CacheWriteTokens)
placeholder

func TestResponsesUsageNestedCacheWritePresenceOverridesTopLevelAlias(t *testing.T) {
	tests := []struct {
		name       string
		nestedJSON string
		want       int
placeholder{
		{name: "explicit zero", nestedJSON: `{"cache_write_tokens":0placeholder`, want: 0placeholder,
		{name: "nonzero", nestedJSON: `{"cache_write_tokens":7placeholder`, want: 7placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var usage ResponsesUsage
			payload := []byte(`{"input_tokens":20,"output_tokens":2,"cache_creation_input_tokens":19,"input_tokens_details":` + tt.nestedJSON + `placeholder`)
			require.NoError(t, json.Unmarshal(payload, &usage))
			require.Equal(t, tt.want, usage.CacheCreationInputTokens)
	placeholder)
placeholder
placeholder

func TestChatCompletionsToResponses_SystemMessage(t *testing.T) {
	req := &ChatCompletionsRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{Role: "system", Content: json.RawMessage(`"You are helpful."`)placeholder,
			{Role: "user", Content: json.RawMessage(`"Hi"`)placeholder,
	placeholder,
placeholder

	resp, err := ChatCompletionsToResponses(req)
placeholder

	var items []ResponsesInputItem
	require.NoError(t, json.Unmarshal(resp.Input, &items))
	require.Len(t, items, 2)
	assert.Equal(t, "system", items[0].Role)
	assert.Equal(t, "user", items[1].Role)
placeholder

func TestChatCompletionsToResponses_ToolCalls(t *testing.T) {
	req := &ChatCompletionsRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{Role: "user", Content: json.RawMessage(`"Call the function"`)placeholder,
			{
				Role: "assistant",
				ToolCalls: []ChatToolCall{
					{
						ID:   "call_1",
						Type: "function",
						Function: ChatFunctionCall{
							Name:      "ping",
							Arguments: `{"host":"example.com"placeholder`,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			{
				Role:       "tool",
				ToolCallID: "call_1",
				Content:    json.RawMessage(`"pong"`),
		placeholder,
	placeholder,
		Tools: []ChatTool{
			{
				Type: "function",
				Function: &ChatFunction{
					Name:        "ping",
					Description: "Ping a host",
					Parameters:  json.RawMessage(`{"type":"object"placeholder`),
			placeholder,
		placeholder,
	placeholder,
placeholder

	resp, err := ChatCompletionsToResponses(req)
placeholder

	var items []ResponsesInputItem
	require.NoError(t, json.Unmarshal(resp.Input, &items))
	// user + function_call + function_call_output = 3
	// (assistant message with empty content + tool_calls → only function_call items emitted)
	require.Len(t, items, 3)

	// Check function_call item
	assert.Equal(t, "function_call", items[1].Type)
	assert.Equal(t, "call_1", items[1].CallID)
	assert.Empty(t, items[1].ID)
	assert.Equal(t, "ping", items[1].Name)

	// Check function_call_output item
	assert.Equal(t, "function_call_output", items[2].Type)
	assert.Equal(t, "call_1", items[2].CallID)
	assert.Equal(t, "pong", items[2].Output)

	// Check tools
	require.Len(t, resp.Tools, 1)
	assert.Equal(t, "function", resp.Tools[0].Type)
	assert.Equal(t, "ping", resp.Tools[0].Name)
placeholder

func TestChatCompletionsToResponses_ToolStrict(t *testing.T) {
	strictTrue := true
	strictFalse := false
	tests := []struct {
		name   string
		strict *bool
		want   bool
placeholder{
		{name: "defaults omitted strict to false", want: falseplaceholder,
		{name: "preserves explicit true", strict: &strictTrue, want: trueplaceholder,
		{name: "preserves explicit false", strict: &strictFalse, want: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ChatCompletionsRequest{
				Model:    "gpt-4o",
				Messages: []ChatMessage{{Role: "user", Content: json.RawMessage(`"Hi"`)placeholderplaceholder,
				Tools: []ChatTool{{
					Type: "function",
					Function: &ChatFunction{
						Name:   "lookup",
						Strict: tt.strict,
				placeholder,
		placeholder
		placeholder

			resp, err := ChatCompletionsToResponses(req)
		placeholder
			require.Len(t, resp.Tools, 1)
			require.NotNil(t, resp.Tools[0].Strict)
			assert.Equal(t, tt.want, *resp.Tools[0].Strict)

			payload, err := json.Marshal(resp)
		placeholder

			var serialized struct {
				Tools []map[string]json.RawMessage `json:"tools"`
		placeholder
			require.NoError(t, json.Unmarshal(payload, &serialized))
			require.Len(t, serialized.Tools, 1)
			strictJSON, ok := serialized.Tools[0]["strict"]
			require.True(t, ok, "strict must be present in the Responses payload")
			assert.JSONEq(t, string(mustMarshalJSON(t, tt.want)), string(strictJSON))
	placeholder)
placeholder
placeholder

func TestChatCompletionsToResponses_LegacyFunctionDefaultsStrictFalse(t *testing.T) {
	req := &ChatCompletionsRequest{
		Model:    "gpt-4o",
		Messages: []ChatMessage{{Role: "user", Content: json.RawMessage(`"Hi"`)placeholderplaceholder,
		Functions: []ChatFunction{{
			Name: "lookup",
placeholder
placeholder

	resp, err := ChatCompletionsToResponses(req)
placeholder
	require.Len(t, resp.Tools, 1)
	require.NotNil(t, resp.Tools[0].Strict)
	assert.False(t, *resp.Tools[0].Strict)

	payload, err := json.Marshal(resp)
placeholder
	assert.Contains(t, string(payload), `"strict":false`)
placeholder

func TestResponsesTool_StrictFalseIsSerialized(t *testing.T) {
	strict := false
	payload, err := json.Marshal(ResponsesTool{
		Type:   "function",
		Strict: &strict,
placeholder)
placeholder
	assert.JSONEq(t, `{"type":"function","strict":falseplaceholder`, string(payload))
placeholder

func mustMarshalJSON(t *testing.T, value any) []byte {
placeholder
	data, err := json.Marshal(value)
placeholder
	return data
placeholder

func TestChatCompletionsToResponses_MaxTokens(t *testing.T) {
	t.Run("max_tokens", func(t *testing.T) {
		maxTokens := 100
		req := &ChatCompletionsRequest{
			Model:     "gpt-4o",
			MaxTokens: &maxTokens,
			Messages:  []ChatMessage{{Role: "user", Content: json.RawMessage(`"Hi"`)placeholderplaceholder,
	placeholder
		resp, err := ChatCompletionsToResponses(req)
	placeholder
		require.NotNil(t, resp.MaxOutputTokens)
		// Below minMaxOutputTokens (128), should be clamped
		assert.Equal(t, minMaxOutputTokens, *resp.MaxOutputTokens)
placeholder)

	t.Run("max_completion_tokens_preferred", func(t *testing.T) {
		maxTokens := 100
		maxCompletion := 500
		req := &ChatCompletionsRequest{
			Model:               "gpt-4o",
			MaxTokens:           &maxTokens,
			MaxCompletionTokens: &maxCompletion,
			Messages:            []ChatMessage{{Role: "user", Content: json.RawMessage(`"Hi"`)placeholderplaceholder,
	placeholder
		resp, err := ChatCompletionsToResponses(req)
	placeholder
		require.NotNil(t, resp.MaxOutputTokens)
		assert.Equal(t, 500, *resp.MaxOutputTokens)
placeholder)
placeholder

func TestChatCompletionsToResponses_ReasoningEffort(t *testing.T) {
	req := &ChatCompletionsRequest{
		Model:           "gpt-4o",
		ReasoningEffort: "high",
		Messages:        []ChatMessage{{Role: "user", Content: json.RawMessage(`"Hi"`)placeholderplaceholder,
placeholder
	resp, err := ChatCompletionsToResponses(req)
placeholder
	require.NotNil(t, resp.Reasoning)
	assert.Equal(t, "high", resp.Reasoning.Effort)
	assert.Equal(t, "auto", resp.Reasoning.Summary)
placeholder

func TestChatCompletionsToResponses_ResponseFormatJsonObject(t *testing.T) {
	req := &ChatCompletionsRequest{
		Model:          "gpt-4o",
		Messages:       []ChatMessage{{Role: "user", Content: json.RawMessage(`"Return JSON"`)placeholderplaceholder,
		ResponseFormat: json.RawMessage(`{"type":"json_object"placeholder`),
placeholder

	resp, err := ChatCompletionsToResponses(req)
placeholder
	require.NotNil(t, resp.Text)
	assert.JSONEq(t, `{"type":"json_object"placeholder`, string(resp.Text.Format))

	payload, err := json.Marshal(resp)
placeholder
	var serialized struct {
		Text ResponsesText `json:"text"`
placeholder
	require.NoError(t, json.Unmarshal(payload, &serialized))
	assert.JSONEq(t, `{"type":"json_object"placeholder`, string(serialized.Text.Format))
placeholder

func TestChatCompletionsToResponses_ResponseFormatJsonSchema(t *testing.T) {
	req := &ChatCompletionsRequest{
		Model:    "gpt-4o",
		Messages: []ChatMessage{{Role: "user", Content: json.RawMessage(`"Return structured JSON"`)placeholderplaceholder,
		ResponseFormat: json.RawMessage(`{
			"type":"json_schema",
			"json_schema":{
				"name":"answer",
				"schema":{
					"type":"object",
					"properties":{"ok":{"type":"boolean"placeholderplaceholder,
					"required":["ok"],
					"additionalProperties":false
			placeholder,
				"strict":true
		placeholder
	placeholder`),
placeholder

	resp, err := ChatCompletionsToResponses(req)
placeholder
	require.NotNil(t, resp.Text)
	assert.JSONEq(t, `{
		"type":"json_schema",
		"name":"answer",
		"schema":{
			"type":"object",
			"properties":{"ok":{"type":"boolean"placeholderplaceholder,
			"required":["ok"],
			"additionalProperties":false
	placeholder,
		"strict":true
placeholder`, string(resp.Text.Format))
placeholder

func TestChatCompletionsToResponses_ImageURL(t *testing.T) {
	content := `[{"type":"text","text":"Describe this"placeholder,{"type":"image_url","image_url":{"url":"data:image/png;base64,abc123"placeholderplaceholder]`
	req := &ChatCompletionsRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{Role: "user", Content: json.RawMessage(content)placeholder,
	placeholder,
placeholder
	resp, err := ChatCompletionsToResponses(req)
placeholder

	var items []ResponsesInputItem
	require.NoError(t, json.Unmarshal(resp.Input, &items))
	require.Len(t, items, 1)

	var parts []ResponsesContentPart
	require.NoError(t, json.Unmarshal(items[0].Content, &parts))
	require.Len(t, parts, 2)
	assert.Equal(t, "input_text", parts[0].Type)
	assert.Equal(t, "Describe this", parts[0].Text)
	assert.Equal(t, "input_image", parts[1].Type)
	assert.Equal(t, "data:image/png;base64,abc123", parts[1].ImageURL)
placeholder

func TestChatCompletionsToResponses_EmptyBase64ImageURLSkipped(t *testing.T) {
	content := `[{"type":"text","text":"Describe this"placeholder,{"type":"image_url","image_url":{"url":"data:image/png;base64,"placeholderplaceholder]`
	req := &ChatCompletionsRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{Role: "user", Content: json.RawMessage(content)placeholder,
	placeholder,
placeholder
	resp, err := ChatCompletionsToResponses(req)
placeholder

	var items []ResponsesInputItem
	require.NoError(t, json.Unmarshal(resp.Input, &items))
	require.Len(t, items, 1)

	var parts []ResponsesContentPart
	require.NoError(t, json.Unmarshal(items[0].Content, &parts))
	require.Len(t, parts, 1)
	assert.Equal(t, "input_text", parts[0].Type)
	assert.Equal(t, "Describe this", parts[0].Text)
placeholder

func TestChatCompletionsToResponses_WhitespaceOnlyBase64ImageURLSkipped(t *testing.T) {
	content := `[{"type":"text","text":"Describe this"placeholder,{"type":"image_url","image_url":{"url":"data:image/png;base64,   "placeholderplaceholder]`
	req := &ChatCompletionsRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{Role: "user", Content: json.RawMessage(content)placeholder,
	placeholder,
placeholder
	resp, err := ChatCompletionsToResponses(req)
placeholder

	var items []ResponsesInputItem
	require.NoError(t, json.Unmarshal(resp.Input, &items))
	require.Len(t, items, 1)

	var parts []ResponsesContentPart
	require.NoError(t, json.Unmarshal(items[0].Content, &parts))
	require.Len(t, parts, 1)
	assert.Equal(t, "input_text", parts[0].Type)
	assert.Equal(t, "Describe this", parts[0].Text)
placeholder

func TestChatCompletionsToResponses_FilePartFileData(t *testing.T) {
	content := `[{"type":"text","text":"Summarize the attached document"placeholder,{"type":"file","file":{"filename":"document.pdf","file_data":"data:application/pdf;base64,JVBERi0xLjQ="placeholderplaceholder]`
	req := &ChatCompletionsRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{Role: "user", Content: json.RawMessage(content)placeholder,
	placeholder,
placeholder
	resp, err := ChatCompletionsToResponses(req)
placeholder

	var items []ResponsesInputItem
	require.NoError(t, json.Unmarshal(resp.Input, &items))
	require.Len(t, items, 1)

	var parts []ResponsesContentPart
	require.NoError(t, json.Unmarshal(items[0].Content, &parts))
	require.Len(t, parts, 2)
	assert.Equal(t, "input_text", parts[0].Type)
	assert.Equal(t, "Summarize the attached document", parts[0].Text)
	assert.Equal(t, "input_file", parts[1].Type)
	assert.Equal(t, "document.pdf", parts[1].Filename)
	assert.Equal(t, "data:application/pdf;base64,JVBERi0xLjQ=", parts[1].FileData)
	assert.Empty(t, parts[1].FileID)
placeholder

func TestChatCompletionsToResponses_FilePartFileID(t *testing.T) {
	content := `[{"type":"file","file":{"file_id":"file-abc123"placeholderplaceholder]`
	req := &ChatCompletionsRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{Role: "user", Content: json.RawMessage(content)placeholder,
	placeholder,
placeholder
	resp, err := ChatCompletionsToResponses(req)
placeholder

	var items []ResponsesInputItem
	require.NoError(t, json.Unmarshal(resp.Input, &items))
	require.Len(t, items, 1)

	var parts []ResponsesContentPart
	require.NoError(t, json.Unmarshal(items[0].Content, &parts))
	require.Len(t, parts, 1)
	assert.Equal(t, "input_file", parts[0].Type)
	assert.Equal(t, "file-abc123", parts[0].FileID)
	assert.Empty(t, parts[0].FileData)
placeholder

func TestChatCompletionsToResponses_EmptyFilePartSkipped(t *testing.T) {
	// A file part with neither file_data nor file_id carries nothing the
	// Responses API can use; dropping it (like empty image URLs) avoids an
	// upstream 400 on an empty input_file part.
	content := `[{"type":"text","text":"Describe this"placeholder,{"type":"file","file":{"filename":"empty.pdf"placeholderplaceholder]`
	req := &ChatCompletionsRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{Role: "user", Content: json.RawMessage(content)placeholder,
	placeholder,
placeholder
	resp, err := ChatCompletionsToResponses(req)
placeholder

	var items []ResponsesInputItem
	require.NoError(t, json.Unmarshal(resp.Input, &items))
	require.Len(t, items, 1)

	var parts []ResponsesContentPart
	require.NoError(t, json.Unmarshal(items[0].Content, &parts))
	require.Len(t, parts, 1)
	assert.Equal(t, "input_text", parts[0].Type)
placeholder

func TestChatCompletionsToResponses_EmptyContentNeverNull(t *testing.T) {
	// Regression for #2515: the upstream Responses API rejects an input item
	// whose content field is JSON null. Any chat-completions message that
	// yields no usable content parts must serialize content as a string.
	cases := []struct {
		name    string
		content json.RawMessage
placeholder{
		{"null content", json.RawMessage(`null`)placeholder,
		{"empty array content", json.RawMessage(`[]`)placeholder,
		{"only empty text part", json.RawMessage(`[{"type":"text","text":""placeholder]`)placeholder,
		{"only empty base64 image part", json.RawMessage(`[{"type":"image_url","image_url":{"url":"data:image/png;base64,"placeholderplaceholder]`)placeholder,
placeholder
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := &ChatCompletionsRequest{
				Model: "gpt-5.5",
				Messages: []ChatMessage{
					{Role: "user", Content: tc.contentplaceholder,
			placeholder,
		placeholder
			resp, err := ChatCompletionsToResponses(req)
		placeholder
			assert.NotContains(t, string(resp.Input), `"content":null`,
				"converted input must not contain a null content field")

			var items []ResponsesInputItem
			require.NoError(t, json.Unmarshal(resp.Input, &items))
			require.Len(t, items, 1)
			assert.Equal(t, `""`, string(items[0].Content),
				"content must be an empty string, not null")
	placeholder)
placeholder
placeholder

func TestChatCompletionsResponseToResponses_DeepSeekReasoningOnlyFallsBackToMessageText(t *testing.T) {
	content := json.RawMessage(`""`)
	resp := &ChatCompletionsResponse{
		ID:     "chatcmpl_deepseek_reasoning_only",
		Object: "chat.completion",
		Model:  "deepseek-reasoner",
		Choices: []ChatChoice{{
			Index: 0,
			Message: ChatMessage{
				Role:             "assistant",
				Content:          content,
				ReasoningContent: "reasoning-only answer",
		placeholder,
			FinishReason: "stop",
placeholder
placeholder

	out := ChatCompletionsResponseToResponses(resp, "deepseek-reasoner", nil, false, nil)

	require.Len(t, out.Output, 2)
	require.Equal(t, "reasoning", out.Output[0].Type)
	require.Equal(t, "message", out.Output[1].Type)
	require.Len(t, out.Output[1].Content, 1)
	assert.Equal(t, "reasoning-only answer", out.Output[1].Content[0].Text)
placeholder

func TestChatCompletionsResponseToResponses_DeepSeekReasoningToolCallDoesNotFallbackToMessageText(t *testing.T) {
	content := json.RawMessage(`""`)
	resp := &ChatCompletionsResponse{
		ID:     "chatcmpl_deepseek_reasoning_tool",
		Object: "chat.completion",
		Model:  "deepseek-reasoner",
		Choices: []ChatChoice{{
			Index: 0,
			Message: ChatMessage{
				Role:             "assistant",
				Content:          content,
				ReasoningContent: "call a tool",
				ToolCalls: []ChatToolCall{{
					ID:   "call_a",
					Type: "function",
					Function: ChatFunctionCall{
						Name:      "exec",
						Arguments: `{placeholder`,
				placeholder,
		placeholder
		placeholder,
			FinishReason: "tool_calls",
placeholder
placeholder

	out := ChatCompletionsResponseToResponses(resp, "deepseek-reasoner", nil, false, nil)

	require.Len(t, out.Output, 2)
	require.Equal(t, "reasoning", out.Output[0].Type)
	require.Equal(t, "function_call", out.Output[1].Type)
	assert.Equal(t, "exec", out.Output[1].Name)
placeholder

func TestChatCompletionsToResponses_SystemArrayContent(t *testing.T) {
	req := &ChatCompletionsRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{Role: "system", Content: json.RawMessage(`[{"type":"text","text":"You are a careful visual assistant."placeholder]`)placeholder,
			{Role: "user", Content: json.RawMessage(`[{"type":"text","text":"Describe this image"placeholder,{"type":"image_url","image_url":{"url":"data:image/png;base64,abc123"placeholderplaceholder]`)placeholder,
	placeholder,
placeholder

	resp, err := ChatCompletionsToResponses(req)
placeholder

	var items []ResponsesInputItem
	require.NoError(t, json.Unmarshal(resp.Input, &items))
	require.Len(t, items, 2)

	var systemParts []ResponsesContentPart
	require.NoError(t, json.Unmarshal(items[0].Content, &systemParts))
	require.Len(t, systemParts, 1)
	assert.Equal(t, "input_text", systemParts[0].Type)
	assert.Equal(t, "You are a careful visual assistant.", systemParts[0].Text)

	var userParts []ResponsesContentPart
	require.NoError(t, json.Unmarshal(items[1].Content, &userParts))
	require.Len(t, userParts, 2)
	assert.Equal(t, "input_image", userParts[1].Type)
	assert.Equal(t, "data:image/png;base64,abc123", userParts[1].ImageURL)
placeholder

func TestChatCompletionsToResponses_LegacyFunctions(t *testing.T) {
	req := &ChatCompletionsRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{Role: "user", Content: json.RawMessage(`"Hi"`)placeholder,
	placeholder,
		Functions: []ChatFunction{
			{
				Name:        "get_weather",
				Description: "Get weather",
				Parameters:  json.RawMessage(`{"type":"object"placeholder`),
		placeholder,
	placeholder,
		FunctionCall: json.RawMessage(`{"name":"get_weather"placeholder`),
placeholder

	resp, err := ChatCompletionsToResponses(req)
placeholder
	require.Len(t, resp.Tools, 1)
	assert.Equal(t, "function", resp.Tools[0].Type)
	assert.Equal(t, "get_weather", resp.Tools[0].Name)

	// tool_choice should be converted
	require.NotNil(t, resp.ToolChoice)
	var tc map[string]any
	require.NoError(t, json.Unmarshal(resp.ToolChoice, &tc))
	assert.Equal(t, "function", tc["type"])
	assert.Equal(t, "get_weather", tc["name"])
	assert.NotContains(t, tc, "function")
placeholder

func TestChatCompletionsToResponses_ServiceTier(t *testing.T) {
	req := &ChatCompletionsRequest{
		Model:       "gpt-4o",
		ServiceTier: "flex",
		Messages:    []ChatMessage{{Role: "user", Content: json.RawMessage(`"Hi"`)placeholderplaceholder,
placeholder
	resp, err := ChatCompletionsToResponses(req)
placeholder
	assert.Equal(t, "flex", resp.ServiceTier)
placeholder

func TestChatCompletionsToResponses_ParallelToolCalls(t *testing.T) {
	for _, value := range []bool{false, trueplaceholder {
		req := &ChatCompletionsRequest{
			Model:             "gpt-4o",
			ParallelToolCalls: &value,
			Messages:          []ChatMessage{{Role: "user", Content: json.RawMessage(`"Hi"`)placeholderplaceholder,
	placeholder

		resp, err := ChatCompletionsToResponses(req)
	placeholder
		require.NotNil(t, resp.ParallelToolCalls)
		assert.Equal(t, value, *resp.ParallelToolCalls)

		payload, err := json.Marshal(resp)
	placeholder
		assert.Contains(t, string(payload), `"parallel_tool_calls":`+string(mustMarshalJSON(t, value)))
placeholder
placeholder

// ---------------------------------------------------------------------------
// temperature / top_p stripping for reasoning models
// ---------------------------------------------------------------------------

func TestChatCompletionsToResponses_TemperatureStrippedForReasoningModel(t *testing.T) {
	temp := 0.7
	req := &ChatCompletionsRequest{
		Model:       "gpt-5.2",
		Messages:    []ChatMessage{{Role: "user", Content: json.RawMessage(`"Hi"`)placeholderplaceholder,
		Temperature: &temp,
		TopP:        &temp,
placeholder

	resp, err := ChatCompletionsToResponses(req)
placeholder
	assert.Nil(t, resp.Temperature, "reasoning model: temperature must be stripped")
	assert.Nil(t, resp.TopP, "reasoning model: top_p must be stripped")

	// Must not appear in the serialised request body sent to the upstream.
	b, err := json.Marshal(resp)
placeholder
	assert.NotContains(t, string(b), `"temperature"`)
	assert.NotContains(t, string(b), `"top_p"`)
placeholder

func TestChatCompletionsToResponses_TemperaturePreservedForNonReasoningModel(t *testing.T) {
	temp := 0.7
	req := &ChatCompletionsRequest{
		Model:       "gpt-4o",
		Messages:    []ChatMessage{{Role: "user", Content: json.RawMessage(`"Hi"`)placeholderplaceholder,
		Temperature: &temp,
		TopP:        &temp,
placeholder

	resp, err := ChatCompletionsToResponses(req)
placeholder
	require.NotNil(t, resp.Temperature, "non-reasoning model: temperature must be preserved")
	assert.InDelta(t, 0.7, *resp.Temperature, 1e-9)
	require.NotNil(t, resp.TopP, "non-reasoning model: top_p must be preserved")
	assert.InDelta(t, 0.7, *resp.TopP, 1e-9)
placeholder

func TestChatCompletionsToResponses_AssistantWithTextAndToolCalls(t *testing.T) {
	req := &ChatCompletionsRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{Role: "user", Content: json.RawMessage(`"Do something"`)placeholder,
			{
				Role:    "assistant",
				Content: json.RawMessage(`"Let me call a function."`),
				ToolCalls: []ChatToolCall{
					{
						ID:   "call_abc",
						Type: "function",
						Function: ChatFunctionCall{
							Name:      "do_thing",
							Arguments: `{placeholder`,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
	placeholder,
placeholder

	resp, err := ChatCompletionsToResponses(req)
placeholder

	var items []ResponsesInputItem
	require.NoError(t, json.Unmarshal(resp.Input, &items))
	// user + assistant message (with text) + function_call
	require.Len(t, items, 3)
	assert.Equal(t, "user", items[0].Role)
	assert.Equal(t, "assistant", items[1].Role)
	assert.Equal(t, "function_call", items[2].Type)
	assert.Empty(t, items[2].ID)
placeholder

func TestChatCompletionsToResponses_AssistantArrayContentPreserved(t *testing.T) {
	req := &ChatCompletionsRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{Role: "user", Content: json.RawMessage(`"Hi"`)placeholder,
			{Role: "assistant", Content: json.RawMessage(`[{"type":"text","text":"A"placeholder,{"type":"text","text":"B"placeholder]`)placeholder,
	placeholder,
placeholder

	resp, err := ChatCompletionsToResponses(req)
placeholder

	var items []ResponsesInputItem
	require.NoError(t, json.Unmarshal(resp.Input, &items))
	require.Len(t, items, 2)
	assert.Equal(t, "assistant", items[1].Role)

	var parts []ResponsesContentPart
	require.NoError(t, json.Unmarshal(items[1].Content, &parts))
	require.Len(t, parts, 1)
	assert.Equal(t, "output_text", parts[0].Type)
	assert.Equal(t, "AB", parts[0].Text)
placeholder

func TestChatCompletionsToResponses_AssistantThinkingTagPreserved(t *testing.T) {
	req := &ChatCompletionsRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{Role: "user", Content: json.RawMessage(`"Hi"`)placeholder,
			{Role: "assistant", Content: json.RawMessage(`[{"type":"thinking","thinking":"internal plan"placeholder,{"type":"text","text":"final answer"placeholder]`)placeholder,
	placeholder,
placeholder

	resp, err := ChatCompletionsToResponses(req)
placeholder

	var items []ResponsesInputItem
	require.NoError(t, json.Unmarshal(resp.Input, &items))
	require.Len(t, items, 2)

	var parts []ResponsesContentPart
	require.NoError(t, json.Unmarshal(items[1].Content, &parts))
	require.Len(t, parts, 1)
	assert.Equal(t, "output_text", parts[0].Type)
	assert.Contains(t, parts[0].Text, "<thinking>internal plan</thinking>")
	assert.Contains(t, parts[0].Text, "final answer")
placeholder

func TestChatCompletionsToResponses_AssistantReasoningContentPreserved(t *testing.T) {
	req := &ChatCompletionsRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{Role: "user", Content: json.RawMessage(`"Hi"`)placeholder,
			{
				Role:             "assistant",
				ReasoningContent: "internal plan",
				Content:          json.RawMessage(`"final answer"`),
		placeholder,
	placeholder,
placeholder

	resp, err := ChatCompletionsToResponses(req)
placeholder

	var items []ResponsesInputItem
	require.NoError(t, json.Unmarshal(resp.Input, &items))
	require.Len(t, items, 2)

	var parts []ResponsesContentPart
	require.NoError(t, json.Unmarshal(items[1].Content, &parts))
	require.Len(t, parts, 1)
	assert.Equal(t, "output_text", parts[0].Type)
	assert.Contains(t, parts[0].Text, "<thinking>internal plan</thinking>")
	assert.Contains(t, parts[0].Text, "final answer")
placeholder

// ---------------------------------------------------------------------------
// ResponsesToChatCompletions tests
// ---------------------------------------------------------------------------

func TestResponsesToChatCompletions_BasicText(t *testing.T) {
	resp := &ResponsesResponse{
		ID:     "resp_123",
		Status: "completed",
		Output: []ResponsesOutput{
			{
				Type: "message",
				Content: []ResponsesContentPart{
					{Type: "output_text", Text: "Hello, world!"placeholder,
			placeholder,
		placeholder,
	placeholder,
		Usage: &ResponsesUsage{
			InputTokens:  10,
			OutputTokens: 5,
			TotalTokens:  15,
	placeholder,
placeholder

	chat := ResponsesToChatCompletions(resp, "gpt-4o")
	assert.Equal(t, "chat.completion", chat.Object)
	assert.Equal(t, "gpt-4o", chat.Model)
	require.Len(t, chat.Choices, 1)
	assert.Equal(t, "stop", chat.Choices[0].FinishReason)

	var content string
	require.NoError(t, json.Unmarshal(chat.Choices[0].Message.Content, &content))
	assert.Equal(t, "Hello, world!", content)

	require.NotNil(t, chat.Usage)
	assert.Equal(t, 10, chat.Usage.PromptTokens)
	assert.Equal(t, 5, chat.Usage.CompletionTokens)
	assert.Equal(t, 15, chat.Usage.TotalTokens)
placeholder

func TestResponsesToChatCompletions_ToolCalls(t *testing.T) {
	resp := &ResponsesResponse{
		ID:     "resp_456",
		Status: "completed",
		Output: []ResponsesOutput{
			{
				Type:      "function_call",
				CallID:    "call_xyz",
				Name:      "get_weather",
				Arguments: `{"city":"NYC"placeholder`,
		placeholder,
	placeholder,
placeholder

	chat := ResponsesToChatCompletions(resp, "gpt-4o")
	require.Len(t, chat.Choices, 1)
	assert.Equal(t, "tool_calls", chat.Choices[0].FinishReason)

	msg := chat.Choices[0].Message
	require.Len(t, msg.ToolCalls, 1)
	assert.Equal(t, "call_xyz", msg.ToolCalls[0].ID)
	assert.Equal(t, "function", msg.ToolCalls[0].Type)
	assert.Equal(t, "get_weather", msg.ToolCalls[0].Function.Name)
	assert.Equal(t, `{"city":"NYC"placeholder`, msg.ToolCalls[0].Function.Arguments)
placeholder

func TestResponsesToChatCompletions_Reasoning(t *testing.T) {
	resp := &ResponsesResponse{
		ID:     "resp_789",
		Status: "completed",
		Output: []ResponsesOutput{
			{
				Type: "reasoning",
				Summary: []ResponsesSummary{
					{Type: "summary_text", Text: "I thought about it."placeholder,
			placeholder,
		placeholder,
			{
				Type: "message",
				Content: []ResponsesContentPart{
					{Type: "output_text", Text: "The answer is 42."placeholder,
			placeholder,
		placeholder,
	placeholder,
placeholder

	chat := ResponsesToChatCompletions(resp, "gpt-4o")
	require.Len(t, chat.Choices, 1)

	var content string
	require.NoError(t, json.Unmarshal(chat.Choices[0].Message.Content, &content))
	assert.Equal(t, "The answer is 42.", content)
	assert.Equal(t, "I thought about it.", chat.Choices[0].Message.ReasoningContent)
placeholder

func TestChatCompletionsToResponses_ToolArrayContent(t *testing.T) {
	req := &ChatCompletionsRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{Role: "user", Content: json.RawMessage(`"Use the tool"`)placeholder,
			{
				Role: "assistant",
				ToolCalls: []ChatToolCall{
					{
						ID:   "call_1",
						Type: "function",
						Function: ChatFunctionCall{
							Name:      "inspect_image",
							Arguments: `{placeholder`,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			{
				Role:       "tool",
				ToolCallID: "call_1",
				Content: json.RawMessage(
					`[{"type":"text","text":"image width: 100"placeholder,{"type":"image_url","image_url":{"url":"data:image/png;base64,ignored"placeholderplaceholder,{"type":"text","text":"; image height: 200"placeholder]`,
				),
		placeholder,
	placeholder,
placeholder

	resp, err := ChatCompletionsToResponses(req)
placeholder

	var items []ResponsesInputItem
	require.NoError(t, json.Unmarshal(resp.Input, &items))
	require.Len(t, items, 3)
	assert.Equal(t, "function_call_output", items[2].Type)
	assert.Equal(t, "call_1", items[2].CallID)
	assert.Equal(t, "image width: 100; image height: 200", items[2].Output)
placeholder

func TestResponsesToChatCompletions_Incomplete(t *testing.T) {
	resp := &ResponsesResponse{
		ID:                "resp_inc",
		Status:            "incomplete",
		IncompleteDetails: &ResponsesIncompleteDetails{Reason: "max_output_tokens"placeholder,
		Output: []ResponsesOutput{
			{
				Type: "message",
				Content: []ResponsesContentPart{
					{Type: "output_text", Text: "partial..."placeholder,
			placeholder,
		placeholder,
	placeholder,
placeholder

	chat := ResponsesToChatCompletions(resp, "gpt-4o")
	require.Len(t, chat.Choices, 1)
	assert.Equal(t, "length", chat.Choices[0].FinishReason)
placeholder

func TestResponsesToChatCompletions_CachedTokens(t *testing.T) {
	resp := &ResponsesResponse{
		ID:     "resp_cache",
		Status: "completed",
		Output: []ResponsesOutput{
			{
				Type:    "message",
				Content: []ResponsesContentPart{{Type: "output_text", Text: "cached"placeholderplaceholder,
		placeholder,
	placeholder,
		Usage: &ResponsesUsage{
			InputTokens:  100,
			OutputTokens: 10,
			TotalTokens:  110,
			InputTokensDetails: &ResponsesInputTokensDetails{
				CachedTokens: 80,
		placeholder,
	placeholder,
placeholder

	chat := ResponsesToChatCompletions(resp, "gpt-4o")
	require.NotNil(t, chat.Usage)
	require.NotNil(t, chat.Usage.PromptTokensDetails)
	assert.Equal(t, 80, chat.Usage.PromptTokensDetails.CachedTokens)
placeholder

func TestResponsesToChatCompletions_ReasoningTokens(t *testing.T) {
	resp := &ResponsesResponse{
		ID:     "resp_reasoning",
		Status: "completed",
		Output: []ResponsesOutput{
			{
				Type:    "message",
				Content: []ResponsesContentPart{{Type: "output_text", Text: "ping"placeholderplaceholder,
		placeholder,
	placeholder,
		Usage: &ResponsesUsage{
			InputTokens:  24,
			OutputTokens: 33,
			TotalTokens:  57,
			OutputTokensDetails: &ResponsesOutputTokensDetails{
				ReasoningTokens: 32,
		placeholder,
	placeholder,
placeholder

	chat := ResponsesToChatCompletions(resp, "gpt-5.5")
	require.NotNil(t, chat.Usage)
	assert.Equal(t, 33, chat.Usage.CompletionTokens)
	require.NotNil(t, chat.Usage.CompletionTokensDetails)
	assert.Equal(t, 32, chat.Usage.CompletionTokensDetails.ReasoningTokens)
placeholder

func TestResponsesToChatCompletions_AllTokenDetailsPassThrough(t *testing.T) {
	// Covers the full OpenAI CompletionUsage detail field set so future audio
	// and prediction-outputs responses propagate without further changes.
	resp := &ResponsesResponse{
		ID:     "resp_full_details",
		Status: "completed",
		Output: []ResponsesOutput{
			{
				Type:    "message",
				Content: []ResponsesContentPart{{Type: "output_text", Text: "x"placeholderplaceholder,
		placeholder,
	placeholder,
		Usage: &ResponsesUsage{
			InputTokens:  100,
			OutputTokens: 50,
			TotalTokens:  150,
			InputTokensDetails: &ResponsesInputTokensDetails{
				CachedTokens: 60,
				AudioTokens:  4,
		placeholder,
			OutputTokensDetails: &ResponsesOutputTokensDetails{
				ReasoningTokens:          30,
				AudioTokens:              2,
				AcceptedPredictionTokens: 10,
				RejectedPredictionTokens: 3,
		placeholder,
	placeholder,
placeholder

	chat := ResponsesToChatCompletions(resp, "gpt-5.5")
	require.NotNil(t, chat.Usage)
	require.NotNil(t, chat.Usage.PromptTokensDetails)
	assert.Equal(t, 60, chat.Usage.PromptTokensDetails.CachedTokens)
	assert.Equal(t, 4, chat.Usage.PromptTokensDetails.AudioTokens)

	require.NotNil(t, chat.Usage.CompletionTokensDetails)
	assert.Equal(t, 30, chat.Usage.CompletionTokensDetails.ReasoningTokens)
	assert.Equal(t, 2, chat.Usage.CompletionTokensDetails.AudioTokens)
	assert.Equal(t, 10, chat.Usage.CompletionTokensDetails.AcceptedPredictionTokens)
	assert.Equal(t, 3, chat.Usage.CompletionTokensDetails.RejectedPredictionTokens)

	raw, err := json.Marshal(chat.Usage)
placeholder
	assert.Contains(t, string(raw), `"prompt_tokens_details"`)
	assert.Contains(t, string(raw), `"completion_tokens_details"`)
	assert.Contains(t, string(raw), `"reasoning_tokens":30`)
	assert.Contains(t, string(raw), `"accepted_prediction_tokens":10`)
placeholder

func TestResponsesToChatCompletions_NoReasoningTokensWhenZero(t *testing.T) {
	// Non-reasoning models do not return reasoning_tokens. The mapping must
	// omit completion_tokens_details entirely rather than emitting a zero-valued
	// field, so non-reasoning responses stay clean.
	resp := &ResponsesResponse{
		ID:     "resp_no_reasoning",
		Status: "completed",
		Output: []ResponsesOutput{
			{
				Type:    "message",
				Content: []ResponsesContentPart{{Type: "output_text", Text: "hi"placeholderplaceholder,
		placeholder,
	placeholder,
		Usage: &ResponsesUsage{
			InputTokens:  10,
			OutputTokens: 5,
			TotalTokens:  15,
			OutputTokensDetails: &ResponsesOutputTokensDetails{
				ReasoningTokens: 0,
		placeholder,
	placeholder,
placeholder

	chat := ResponsesToChatCompletions(resp, "gpt-4o")
	require.NotNil(t, chat.Usage)
	assert.Nil(t, chat.Usage.CompletionTokensDetails)

	raw, err := json.Marshal(chat.Usage)
placeholder
	assert.NotContains(t, string(raw), "completion_tokens_details")
	assert.NotContains(t, string(raw), "reasoning_tokens")
placeholder

func TestResponsesToChatCompletions_WebSearch(t *testing.T) {
	resp := &ResponsesResponse{
		ID:     "resp_ws",
		Status: "completed",
		Output: []ResponsesOutput{
			{
				Type:   "web_search_call",
				Action: &WebSearchAction{Type: "search", Query: "test"placeholder,
		placeholder,
			{
				Type:    "message",
				Content: []ResponsesContentPart{{Type: "output_text", Text: "search results"placeholderplaceholder,
		placeholder,
	placeholder,
placeholder

	chat := ResponsesToChatCompletions(resp, "gpt-4o")
	require.Len(t, chat.Choices, 1)
	assert.Equal(t, "stop", chat.Choices[0].FinishReason)

	var content string
	require.NoError(t, json.Unmarshal(chat.Choices[0].Message.Content, &content))
	assert.Equal(t, "search results", content)
placeholder

// ---------------------------------------------------------------------------
// Streaming: ResponsesEventToChatChunks tests
// ---------------------------------------------------------------------------

func TestResponsesEventToChatChunks_TextDelta(t *testing.T) {
	state := NewResponsesEventToChatState()
	state.Model = "gpt-4o"

	// response.created → role chunk
	chunks := ResponsesEventToChatChunks(&ResponsesStreamEvent{
		Type: "response.created",
		Response: &ResponsesResponse{
			ID: "resp_stream",
	placeholder,
placeholder, state)
	require.Len(t, chunks, 1)
	assert.Equal(t, "assistant", chunks[0].Choices[0].Delta.Role)
	assert.True(t, state.SentRole)

	// response.output_text.delta → content chunk
	chunks = ResponsesEventToChatChunks(&ResponsesStreamEvent{
		Type:  "response.output_text.delta",
		Delta: "Hello",
placeholder, state)
	require.Len(t, chunks, 1)
	require.NotNil(t, chunks[0].Choices[0].Delta.Content)
	assert.Equal(t, "Hello", *chunks[0].Choices[0].Delta.Content)
placeholder

func TestResponsesEventToChatChunks_ToolCallDelta(t *testing.T) {
	state := NewResponsesEventToChatState()
	state.Model = "gpt-4o"
	state.SentRole = true

	// response.output_item.added (function_call) — output_index=1 (e.g. after a message item at 0)
	chunks := ResponsesEventToChatChunks(&ResponsesStreamEvent{
		Type:        "response.output_item.added",
		OutputIndex: 1,
		Item: &ResponsesOutput{
			Type:   "function_call",
			CallID: "call_1",
			Name:   "get_weather",
	placeholder,
placeholder, state)
	require.Len(t, chunks, 1)
	require.Len(t, chunks[0].Choices[0].Delta.ToolCalls, 1)
	tc := chunks[0].Choices[0].Delta.ToolCalls[0]
	assert.Equal(t, "call_1", tc.ID)
	assert.Equal(t, "get_weather", tc.Function.Name)
	require.NotNil(t, tc.Index)
	assert.Equal(t, 0, *tc.Index)

	// response.function_call_arguments.delta — uses output_index (NOT call_id) to find tool
	chunks = ResponsesEventToChatChunks(&ResponsesStreamEvent{
		Type:        "response.function_call_arguments.delta",
		OutputIndex: 1, // matches the output_index from output_item.added above
		Delta:       `{"city":`,
placeholder, state)
	require.Len(t, chunks, 1)
	tc = chunks[0].Choices[0].Delta.ToolCalls[0]
	require.NotNil(t, tc.Index)
	assert.Equal(t, 0, *tc.Index, "argument delta must use same index as the tool call")
	assert.Equal(t, `{"city":`, tc.Function.Arguments)

	// Add a second function call at output_index=2
	chunks = ResponsesEventToChatChunks(&ResponsesStreamEvent{
		Type:        "response.output_item.added",
		OutputIndex: 2,
		Item: &ResponsesOutput{
			Type:   "function_call",
			CallID: "call_2",
			Name:   "get_time",
	placeholder,
placeholder, state)
	require.Len(t, chunks, 1)
	tc = chunks[0].Choices[0].Delta.ToolCalls[0]
	require.NotNil(t, tc.Index)
	assert.Equal(t, 1, *tc.Index, "second tool call should get index 1")

	// Argument delta for second tool call
	chunks = ResponsesEventToChatChunks(&ResponsesStreamEvent{
		Type:        "response.function_call_arguments.delta",
		OutputIndex: 2,
		Delta:       `{"tz":"UTC"placeholder`,
placeholder, state)
	require.Len(t, chunks, 1)
	tc = chunks[0].Choices[0].Delta.ToolCalls[0]
	require.NotNil(t, tc.Index)
	assert.Equal(t, 1, *tc.Index, "second tool arg delta must use index 1")

	// Argument delta for first tool call (interleaved)
	chunks = ResponsesEventToChatChunks(&ResponsesStreamEvent{
		Type:        "response.function_call_arguments.delta",
		OutputIndex: 1,
		Delta:       `"Tokyo"placeholder`,
placeholder, state)
	require.Len(t, chunks, 1)
	tc = chunks[0].Choices[0].Delta.ToolCalls[0]
	require.NotNil(t, tc.Index)
	assert.Equal(t, 0, *tc.Index, "first tool arg delta must still use index 0")
placeholder

func TestResponsesEventToChatChunks_Completed(t *testing.T) {
	state := NewResponsesEventToChatState()
	state.Model = "gpt-4o"
	state.IncludeUsage = true

	chunks := ResponsesEventToChatChunks(&ResponsesStreamEvent{
		Type: "response.completed",
		Response: &ResponsesResponse{
			Status: "completed",
			Usage: &ResponsesUsage{
				InputTokens:  50,
				OutputTokens: 20,
				TotalTokens:  70,
				InputTokensDetails: &ResponsesInputTokensDetails{
					CachedTokens: 30,
			placeholder,
		placeholder,
	placeholder,
placeholder, state)
	// finish chunk + usage chunk
	require.Len(t, chunks, 2)

	// First chunk: finish_reason
	require.NotNil(t, chunks[0].Choices[0].FinishReason)
	assert.Equal(t, "stop", *chunks[0].Choices[0].FinishReason)

	// Second chunk: usage
	require.NotNil(t, chunks[1].Usage)
	assert.Equal(t, 50, chunks[1].Usage.PromptTokens)
	assert.Equal(t, 20, chunks[1].Usage.CompletionTokens)
	assert.Equal(t, 70, chunks[1].Usage.TotalTokens)
	require.NotNil(t, chunks[1].Usage.PromptTokensDetails)
	assert.Equal(t, 30, chunks[1].Usage.PromptTokensDetails.CachedTokens)
placeholder

func TestResponsesEventToChatChunks_CompletedWithReasoningTokens(t *testing.T) {
	state := NewResponsesEventToChatState()
	state.Model = "gpt-5.5"
	state.IncludeUsage = true

	chunks := ResponsesEventToChatChunks(&ResponsesStreamEvent{
		Type: "response.completed",
		Response: &ResponsesResponse{
			Status: "completed",
			Usage: &ResponsesUsage{
				InputTokens:  24,
				OutputTokens: 33,
				TotalTokens:  57,
				OutputTokensDetails: &ResponsesOutputTokensDetails{
					ReasoningTokens: 32,
			placeholder,
		placeholder,
	placeholder,
placeholder, state)
	require.Len(t, chunks, 2)

	require.NotNil(t, chunks[1].Usage)
	require.NotNil(t, chunks[1].Usage.CompletionTokensDetails)
	assert.Equal(t, 32, chunks[1].Usage.CompletionTokensDetails.ReasoningTokens)
placeholder

func TestResponsesEventToChatChunks_ResponseDone(t *testing.T) {
	state := NewResponsesEventToChatState()
	state.Model = "gpt-4o"
	state.IncludeUsage = true

	chunks := ResponsesEventToChatChunks(&ResponsesStreamEvent{
		Type: "response.done",
		Response: &ResponsesResponse{
			Status: "completed",
			Usage:  &ResponsesUsage{InputTokens: 13, OutputTokens: 7placeholder,
	placeholder,
placeholder, state)
	require.Len(t, chunks, 2)
	require.NotNil(t, chunks[0].Choices[0].FinishReason)
	assert.Equal(t, "stop", *chunks[0].Choices[0].FinishReason)
	require.NotNil(t, chunks[1].Usage)
	assert.Equal(t, 13, chunks[1].Usage.PromptTokens)
	assert.Equal(t, 7, chunks[1].Usage.CompletionTokens)
	assert.Nil(t, FinalizeResponsesChatStream(state))
placeholder

func TestResponsesEventToChatChunks_TopLevelTerminalUsage(t *testing.T) {
	state := NewResponsesEventToChatState()
	state.Model = "gpt-4o"
	state.IncludeUsage = true

	chunks := ResponsesEventToChatChunks(&ResponsesStreamEvent{
		Type: "response.completed",
		Response: &ResponsesResponse{
			Status: "completed",
	placeholder,
		Usage: &ResponsesUsage{
			InputTokens:  21,
			OutputTokens: 9,
			InputTokensDetails: &ResponsesInputTokensDetails{
				CachedTokens: 4,
		placeholder,
	placeholder,
placeholder, state)

	require.Len(t, chunks, 2)
	require.NotNil(t, chunks[1].Usage)
	assert.Equal(t, 21, chunks[1].Usage.PromptTokens)
	assert.Equal(t, 9, chunks[1].Usage.CompletionTokens)
	require.NotNil(t, chunks[1].Usage.PromptTokensDetails)
	assert.Equal(t, 4, chunks[1].Usage.PromptTokensDetails.CachedTokens)
placeholder

func TestResponsesEventToChatChunks_ResponseDoneIncomplete(t *testing.T) {
	state := NewResponsesEventToChatState()
	state.Model = "gpt-4o"
	state.IncludeUsage = true

	chunks := ResponsesEventToChatChunks(&ResponsesStreamEvent{
		Type: "response.done",
		Response: &ResponsesResponse{
			Status:            "incomplete",
			IncompleteDetails: &ResponsesIncompleteDetails{Reason: "max_output_tokens"placeholder,
			Usage:             &ResponsesUsage{InputTokens: 13, OutputTokens: 7placeholder,
	placeholder,
placeholder, state)
	require.Len(t, chunks, 2)
	require.NotNil(t, chunks[0].Choices[0].FinishReason)
	assert.Equal(t, "length", *chunks[0].Choices[0].FinishReason)
	require.NotNil(t, chunks[1].Usage)
	assert.Equal(t, 13, chunks[1].Usage.PromptTokens)
	assert.Equal(t, 7, chunks[1].Usage.CompletionTokens)
	assert.Nil(t, FinalizeResponsesChatStream(state))
placeholder

func TestResponsesEventToChatChunks_CompletedWithToolCalls(t *testing.T) {
	state := NewResponsesEventToChatState()
	state.Model = "gpt-4o"
	state.SawToolCall = true

	chunks := ResponsesEventToChatChunks(&ResponsesStreamEvent{
		Type: "response.completed",
		Response: &ResponsesResponse{
			Status: "completed",
	placeholder,
placeholder, state)
	require.Len(t, chunks, 1)
	require.NotNil(t, chunks[0].Choices[0].FinishReason)
	assert.Equal(t, "tool_calls", *chunks[0].Choices[0].FinishReason)
placeholder

func TestResponsesEventToChatChunks_ReasoningDelta(t *testing.T) {
	state := NewResponsesEventToChatState()
	state.Model = "gpt-4o"
	state.SentRole = true

	chunks := ResponsesEventToChatChunks(&ResponsesStreamEvent{
		Type:  "response.reasoning_summary_text.delta",
		Delta: "Thinking...",
placeholder, state)
	require.Len(t, chunks, 1)
	require.NotNil(t, chunks[0].Choices[0].Delta.ReasoningContent)
	assert.Equal(t, "Thinking...", *chunks[0].Choices[0].Delta.ReasoningContent)

	chunks = ResponsesEventToChatChunks(&ResponsesStreamEvent{
		Type: "response.reasoning_summary_text.done",
placeholder, state)
	require.Len(t, chunks, 0)
placeholder

func TestResponsesEventToChatChunks_ReasoningThenTextAutoCloseTag(t *testing.T) {
	state := NewResponsesEventToChatState()
	state.Model = "gpt-4o"
	state.SentRole = true

	chunks := ResponsesEventToChatChunks(&ResponsesStreamEvent{
		Type:  "response.reasoning_summary_text.delta",
		Delta: "plan",
placeholder, state)
	require.Len(t, chunks, 1)
	require.NotNil(t, chunks[0].Choices[0].Delta.ReasoningContent)
	assert.Equal(t, "plan", *chunks[0].Choices[0].Delta.ReasoningContent)

	chunks = ResponsesEventToChatChunks(&ResponsesStreamEvent{
		Type:  "response.output_text.delta",
		Delta: "answer",
placeholder, state)
	require.Len(t, chunks, 1)
	require.NotNil(t, chunks[0].Choices[0].Delta.Content)
	assert.Equal(t, "answer", *chunks[0].Choices[0].Delta.Content)
placeholder

func TestFinalizeResponsesChatStream(t *testing.T) {
	state := NewResponsesEventToChatState()
	state.Model = "gpt-4o"
	state.IncludeUsage = true
	state.Usage = &ChatUsage{
		PromptTokens:     100,
		CompletionTokens: 50,
		TotalTokens:      150,
placeholder

	chunks := FinalizeResponsesChatStream(state)
	require.Len(t, chunks, 2)

	// Finish chunk
	require.NotNil(t, chunks[0].Choices[0].FinishReason)
	assert.Equal(t, "stop", *chunks[0].Choices[0].FinishReason)

	// Usage chunk
	require.NotNil(t, chunks[1].Usage)
	assert.Equal(t, 100, chunks[1].Usage.PromptTokens)

	// Idempotent: second call returns nil
	assert.Nil(t, FinalizeResponsesChatStream(state))
placeholder

func TestFinalizeResponsesChatStream_AfterCompleted(t *testing.T) {
	// If response.completed already emitted the finish chunk, FinalizeResponsesChatStream
	// must be a no-op (prevents double finish_reason being sent to the client).
	state := NewResponsesEventToChatState()
	state.Model = "gpt-4o"
	state.IncludeUsage = true

	// Simulate response.completed
	chunks := ResponsesEventToChatChunks(&ResponsesStreamEvent{
		Type: "response.completed",
		Response: &ResponsesResponse{
			Status: "completed",
			Usage: &ResponsesUsage{
				InputTokens:  10,
				OutputTokens: 5,
				TotalTokens:  15,
		placeholder,
	placeholder,
placeholder, state)
	require.NotEmpty(t, chunks) // finish + usage chunks

	// Now FinalizeResponsesChatStream should return nil — already finalized.
	assert.Nil(t, FinalizeResponsesChatStream(state))
placeholder

func TestChatChunkToSSE(t *testing.T) {
	chunk := ChatCompletionsChunk{
		ID:      "chatcmpl-test",
		Object:  "chat.completion.chunk",
		Created: 1700000000,
		Model:   "gpt-4o",
		Choices: []ChatChunkChoice{
			{
				Index:        0,
				Delta:        ChatDelta{Role: "assistant"placeholder,
				FinishReason: nil,
		placeholder,
	placeholder,
placeholder

	sse, err := ChatChunkToSSE(chunk)
placeholder
	assert.Contains(t, sse, "data: ")
	assert.Contains(t, sse, "chatcmpl-test")
	assert.Contains(t, sse, "assistant")
	assert.True(t, len(sse) > 10)
placeholder

// ---------------------------------------------------------------------------
// Stream round-trip test
// ---------------------------------------------------------------------------

func TestChatCompletionsStreamRoundTrip(t *testing.T) {
	// Simulate: client sends chat completions request, upstream returns Responses SSE events.
	// Verify that the streaming state machine produces correct chat completions chunks.

	state := NewResponsesEventToChatState()
	state.Model = "gpt-4o"
	state.IncludeUsage = true

	var allChunks []ChatCompletionsChunk

	// 1. response.created
	chunks := ResponsesEventToChatChunks(&ResponsesStreamEvent{
		Type:     "response.created",
		Response: &ResponsesResponse{ID: "resp_rt"placeholder,
placeholder, state)
	allChunks = append(allChunks, chunks...)

	// 2. text deltas
	for _, text := range []string{"Hello", ", ", "world", "!"placeholder {
		chunks = ResponsesEventToChatChunks(&ResponsesStreamEvent{
			Type:  "response.output_text.delta",
			Delta: text,
	placeholder, state)
		allChunks = append(allChunks, chunks...)
placeholder

	// 3. response.completed
	chunks = ResponsesEventToChatChunks(&ResponsesStreamEvent{
		Type: "response.completed",
		Response: &ResponsesResponse{
			Status: "completed",
			Usage: &ResponsesUsage{
				InputTokens:  10,
				OutputTokens: 4,
				TotalTokens:  14,
		placeholder,
	placeholder,
placeholder, state)
	allChunks = append(allChunks, chunks...)

	// Verify: role chunk + 4 text chunks + finish chunk + usage chunk = 7
	require.Len(t, allChunks, 7)

	// First chunk has role
	assert.Equal(t, "assistant", allChunks[0].Choices[0].Delta.Role)

	// Text chunks
	var fullText string
	for i := 1; i <= 4; i++ {
		require.NotNil(t, allChunks[i].Choices[0].Delta.Content)
		fullText += *allChunks[i].Choices[0].Delta.Content
placeholder
	assert.Equal(t, "Hello, world!", fullText)

	// Finish chunk
	require.NotNil(t, allChunks[5].Choices[0].FinishReason)
	assert.Equal(t, "stop", *allChunks[5].Choices[0].FinishReason)

	// Usage chunk
	require.NotNil(t, allChunks[6].Usage)
	assert.Equal(t, 10, allChunks[6].Usage.PromptTokens)
	assert.Equal(t, 4, allChunks[6].Usage.CompletionTokens)

	// All chunks share the same ID
	for _, c := range allChunks {
		assert.Equal(t, "resp_rt", c.ID)
placeholder
placeholder

// ---------------------------------------------------------------------------
// BufferedResponseAccumulator tests
// ---------------------------------------------------------------------------

func TestBufferedResponseAccumulator_TextOnly(t *testing.T) {
	acc := NewBufferedResponseAccumulator()

	acc.ProcessEvent(&ResponsesStreamEvent{Type: "response.output_text.delta", Delta: "Hello"placeholder)
	acc.ProcessEvent(&ResponsesStreamEvent{Type: "response.output_text.delta", Delta: ", world!"placeholder)

	assert.True(t, acc.HasContent())

	output := acc.BuildOutput()
	require.Len(t, output, 1)
	assert.Equal(t, "message", output[0].Type)
	assert.Equal(t, "assistant", output[0].Role)
	require.Len(t, output[0].Content, 1)
	assert.Equal(t, "output_text", output[0].Content[0].Type)
	assert.Equal(t, "Hello, world!", output[0].Content[0].Text)
placeholder

func TestBufferedResponseAccumulator_ToolCalls(t *testing.T) {
	acc := NewBufferedResponseAccumulator()

	// Add function call at output_index=1
	acc.ProcessEvent(&ResponsesStreamEvent{
		Type:        "response.output_item.added",
		OutputIndex: 1,
		Item: &ResponsesOutput{
			Type:   "function_call",
			CallID: "call_abc",
			Name:   "get_weather",
	placeholder,
placeholder)
	acc.ProcessEvent(&ResponsesStreamEvent{
		Type:        "response.function_call_arguments.delta",
		OutputIndex: 1,
		Delta:       `{"city":`,
placeholder)
	acc.ProcessEvent(&ResponsesStreamEvent{
		Type:        "response.function_call_arguments.delta",
		OutputIndex: 1,
		Delta:       `"NYC"placeholder`,
placeholder)

	assert.True(t, acc.HasContent())

	output := acc.BuildOutput()
	require.Len(t, output, 1)
	assert.Equal(t, "function_call", output[0].Type)
	assert.Equal(t, "call_abc", output[0].CallID)
	assert.Equal(t, "get_weather", output[0].Name)
	assert.Equal(t, `{"city":"NYC"placeholder`, output[0].Arguments)
placeholder

func TestBufferedResponseAccumulator_Reasoning(t *testing.T) {
	acc := NewBufferedResponseAccumulator()

	acc.ProcessEvent(&ResponsesStreamEvent{Type: "response.reasoning_summary_text.delta", Delta: "Step 1: "placeholder)
	acc.ProcessEvent(&ResponsesStreamEvent{Type: "response.reasoning_summary_text.delta", Delta: "think about it"placeholder)

	assert.True(t, acc.HasContent())

	output := acc.BuildOutput()
	require.Len(t, output, 1)
	assert.Equal(t, "reasoning", output[0].Type)
	require.Len(t, output[0].Summary, 1)
	assert.Equal(t, "summary_text", output[0].Summary[0].Type)
	assert.Equal(t, "Step 1: think about it", output[0].Summary[0].Text)
placeholder

func TestBufferedResponseAccumulator_Mixed(t *testing.T) {
	acc := NewBufferedResponseAccumulator()

	// Reasoning first
	acc.ProcessEvent(&ResponsesStreamEvent{Type: "response.reasoning_summary_text.delta", Delta: "I thought about it."placeholder)

	// Then text
	acc.ProcessEvent(&ResponsesStreamEvent{Type: "response.output_text.delta", Delta: "The answer is 42."placeholder)

	// Then a tool call
	acc.ProcessEvent(&ResponsesStreamEvent{
		Type:        "response.output_item.added",
		OutputIndex: 2,
		Item: &ResponsesOutput{
			Type:   "function_call",
			CallID: "call_1",
			Name:   "verify",
	placeholder,
placeholder)
	acc.ProcessEvent(&ResponsesStreamEvent{
		Type:        "response.function_call_arguments.delta",
		OutputIndex: 2,
		Delta:       `{placeholder`,
placeholder)

	assert.True(t, acc.HasContent())

	output := acc.BuildOutput()
	// Order: reasoning → message → function_calls
	require.Len(t, output, 3)
	assert.Equal(t, "reasoning", output[0].Type)
	assert.Equal(t, "message", output[1].Type)
	assert.Equal(t, "function_call", output[2].Type)
	assert.Equal(t, "The answer is 42.", output[1].Content[0].Text)
	assert.Equal(t, "verify", output[2].Name)
placeholder

func TestBufferedResponseAccumulator_SupplementEmptyOutput(t *testing.T) {
	acc := NewBufferedResponseAccumulator()
	acc.ProcessEvent(&ResponsesStreamEvent{Type: "response.output_text.delta", Delta: "Hello"placeholder)

	resp := &ResponsesResponse{
		ID:     "resp_1",
		Status: "completed",
		Output: nil, // empty output
		Usage:  &ResponsesUsage{InputTokens: 10, OutputTokens: 5placeholder,
placeholder

	acc.SupplementResponseOutput(resp)

	require.Len(t, resp.Output, 1)
	assert.Equal(t, "message", resp.Output[0].Type)
	assert.Equal(t, "Hello", resp.Output[0].Content[0].Text)
	// Usage should be untouched
	assert.Equal(t, 10, resp.Usage.InputTokens)
placeholder

func TestBufferedResponseAccumulator_NoSupplementWhenOutputExists(t *testing.T) {
	acc := NewBufferedResponseAccumulator()
	acc.ProcessEvent(&ResponsesStreamEvent{Type: "response.output_text.delta", Delta: "from deltas"placeholder)

	resp := &ResponsesResponse{
		ID:     "resp_2",
		Status: "completed",
		Output: []ResponsesOutput{
			{
				Type: "message",
				Content: []ResponsesContentPart{
					{Type: "output_text", Text: "from terminal event"placeholder,
			placeholder,
		placeholder,
	placeholder,
placeholder

	acc.SupplementResponseOutput(resp)

	// Output should NOT be overwritten
	require.Len(t, resp.Output, 1)
	assert.Equal(t, "from terminal event", resp.Output[0].Content[0].Text)
placeholder

func TestBufferedResponseAccumulator_EmptyDeltas(t *testing.T) {
	acc := NewBufferedResponseAccumulator()

	// Process events with empty delta — should not accumulate
	acc.ProcessEvent(&ResponsesStreamEvent{Type: "response.output_text.delta", Delta: ""placeholder)
	acc.ProcessEvent(&ResponsesStreamEvent{Type: "response.created"placeholder)

	assert.False(t, acc.HasContent())

	resp := &ResponsesResponse{ID: "resp_3", Status: "completed"placeholder
	acc.SupplementResponseOutput(resp)
	assert.Nil(t, resp.Output)
placeholder

func TestBufferedResponseAccumulator_IgnoresNonFunctionCallItems(t *testing.T) {
	acc := NewBufferedResponseAccumulator()

	// output_item.added with type "message" should be ignored
	acc.ProcessEvent(&ResponsesStreamEvent{
		Type:        "response.output_item.added",
		OutputIndex: 0,
		Item:        &ResponsesOutput{Type: "message"placeholder,
placeholder)

	assert.False(t, acc.HasContent())
placeholder
