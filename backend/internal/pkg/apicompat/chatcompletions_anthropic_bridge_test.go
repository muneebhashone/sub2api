package apicompat

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func strPtr(s string) *string { return &s placeholder
func intPtr(v int) *int       { return &v placeholder

// collectAnthropicStreamEvents feeds CC chunks through the direct bridge and
// appends finalize events, returning the full Anthropic event sequence.
func collectAnthropicStreamEvents(t *testing.T, chunks []string) []AnthropicStreamEvent {
placeholder
	state := NewChatCompletionsToAnthropicStreamState("deepseek-v4-pro")
	var events []AnthropicStreamEvent
	for _, payload := range chunks {
		var chunk ChatCompletionsChunk
		require.NoError(t, json.Unmarshal([]byte(payload), &chunk))
		events = append(events, ChatCompletionsChunkToAnthropicEvents(&chunk, state)...)
placeholder
	events = append(events, FinalizeChatCompletionsAnthropicStream(state)...)
	return events
placeholder

// anthropicEventTypes extracts the sequence of event types for concise assertions.
func anthropicEventTypes(events []AnthropicStreamEvent) []string {
	out := make([]string, 0, len(events))
	for _, e := range events {
		out = append(out, e.Type)
placeholder
	return out
placeholder

// ---------------------------------------------------------------------------
// Request: AnthropicToChatCompletionsRequest
// ---------------------------------------------------------------------------

func TestAnthropicToChatCompletionsRequest_BasicText(t *testing.T) {
	req := &AnthropicRequest{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 1024,
		Messages: []AnthropicMessage{
			{Role: "user", Content: json.RawMessage(`"hello"`)placeholder,
	placeholder,
placeholder

	out, err := AnthropicToChatCompletionsRequest(req)
placeholder
	require.Equal(t, "claude-sonnet-4-20250514", out.Model)
	require.Len(t, out.Messages, 1)
	require.Equal(t, "user", out.Messages[0].Role)
	require.Equal(t, `"hello"`, string(out.Messages[0].Content))
	require.NotNil(t, out.MaxCompletionTokens)
	require.Equal(t, 1024, *out.MaxCompletionTokens)
placeholder

func TestAnthropicToChatCompletionsRequest_SystemPrompt(t *testing.T) {
	req := &AnthropicRequest{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 100,
		System:    json.RawMessage(`"You are helpful"`),
		Messages: []AnthropicMessage{
			{Role: "user", Content: json.RawMessage(`"hi"`)placeholder,
	placeholder,
placeholder

	out, err := AnthropicToChatCompletionsRequest(req)
placeholder
	require.Len(t, out.Messages, 2)
	require.Equal(t, "system", out.Messages[0].Role)
	require.Equal(t, `"You are helpful"`, string(out.Messages[0].Content))
placeholder

func TestAnthropicToChatCompletionsRequest_ToolUseInAssistant(t *testing.T) {
	req := &AnthropicRequest{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 100,
		Messages: []AnthropicMessage{
			{Role: "user", Content: json.RawMessage(`"check weather"`)placeholder,
			{Role: "assistant", Content: json.RawMessage(`[{"type":"text","text":"Let me check."placeholder,{"type":"tool_use","id":"toolu_1","name":"get_weather","input":{"city":"SF"placeholderplaceholder]`)placeholder,
			{Role: "user", Content: json.RawMessage(`[{"type":"tool_result","tool_use_id":"toolu_1","content":"sunny"placeholder]`)placeholder,
	placeholder,
placeholder

	out, err := AnthropicToChatCompletionsRequest(req)
placeholder
	// user + assistant(with tool_calls) + tool reply
	require.GreaterOrEqual(t, len(out.Messages), 2)
	// Find the assistant message with tool_calls
	var assistant *ChatMessage
	for i := range out.Messages {
		if out.Messages[i].Role == "assistant" && len(out.Messages[i].ToolCalls) > 0 {
			assistant = &out.Messages[i]
	placeholder
placeholder
	require.NotNil(t, assistant, "assistant message with tool_calls should survive normalization")
	require.Len(t, assistant.ToolCalls, 1)
	require.Equal(t, "toolu_1", assistant.ToolCalls[0].ID)
	require.Equal(t, "function", assistant.ToolCalls[0].Type)
	require.Equal(t, "get_weather", assistant.ToolCalls[0].Function.Name)
	require.Equal(t, `{"city":"SF"placeholder`, assistant.ToolCalls[0].Function.Arguments)
placeholder

func TestAnthropicToChatCompletionsRequest_ToolResultBecomesToolMessage(t *testing.T) {
	req := &AnthropicRequest{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 100,
		Messages: []AnthropicMessage{
			{Role: "user", Content: json.RawMessage(`"check weather"`)placeholder,
			{Role: "assistant", Content: json.RawMessage(`[{"type":"tool_use","id":"toolu_1","name":"get_weather","input":{"city":"SF"placeholderplaceholder]`)placeholder,
			{Role: "user", Content: json.RawMessage(`[{"type":"tool_result","tool_use_id":"toolu_1","content":"sunny, 72F"placeholder]`)placeholder,
	placeholder,
placeholder

	out, err := AnthropicToChatCompletionsRequest(req)
placeholder
	// Find the tool reply message
	var toolMsg *ChatMessage
	for i := range out.Messages {
		if out.Messages[i].Role == "tool" {
			toolMsg = &out.Messages[i]
	placeholder
placeholder
	require.NotNil(t, toolMsg, "tool_result should become a tool role message")
	require.Equal(t, "toolu_1", toolMsg.ToolCallID)
	require.Equal(t, `"sunny, 72F"`, string(toolMsg.Content))
placeholder

func TestAnthropicToChatCompletionsRequest_ThinkingDropped(t *testing.T) {
	req := &AnthropicRequest{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 100,
		Messages: []AnthropicMessage{
			{Role: "assistant", Content: json.RawMessage(`[{"type":"thinking","thinking":"secret thoughts"placeholder,{"type":"text","text":"answer"placeholder]`)placeholder,
	placeholder,
placeholder

	out, err := AnthropicToChatCompletionsRequest(req)
placeholder
	require.Len(t, out.Messages, 1)
	// Only text survives; thinking is dropped
	require.Equal(t, `"answer"`, string(out.Messages[0].Content))
placeholder

func TestAnthropicToChatCompletionsRequest_ToolChoiceAuto(t *testing.T) {
	req := &AnthropicRequest{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 100,
		Tools: []AnthropicTool{
			{Name: "get_weather", InputSchema: json.RawMessage(`{"type":"object","properties":{placeholderplaceholder`)placeholder,
	placeholder,
		ToolChoice: json.RawMessage(`{"type":"auto"placeholder`),
		Messages:   []AnthropicMessage{{Role: "user", Content: json.RawMessage(`"hi"`)placeholderplaceholder,
placeholder

	out, err := AnthropicToChatCompletionsRequest(req)
placeholder
	require.Len(t, out.Tools, 1)
	require.Equal(t, `"auto"`, string(out.ToolChoice))
placeholder

func TestAnthropicToChatCompletionsRequest_ToolChoiceAny(t *testing.T) {
	req := &AnthropicRequest{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 100,
		Tools: []AnthropicTool{
			{Name: "get_weather", InputSchema: json.RawMessage(`{"type":"object","properties":{placeholderplaceholder`)placeholder,
	placeholder,
		ToolChoice: json.RawMessage(`{"type":"any"placeholder`),
		Messages:   []AnthropicMessage{{Role: "user", Content: json.RawMessage(`"hi"`)placeholderplaceholder,
placeholder

	out, err := AnthropicToChatCompletionsRequest(req)
placeholder
	require.Equal(t, `"required"`, string(out.ToolChoice))
placeholder

func TestAnthropicToChatCompletionsRequest_ToolChoiceSpecificTool(t *testing.T) {
	req := &AnthropicRequest{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 100,
		Tools: []AnthropicTool{
			{Name: "get_weather", InputSchema: json.RawMessage(`{"type":"object","properties":{placeholderplaceholder`)placeholder,
	placeholder,
		ToolChoice: json.RawMessage(`{"type":"tool","name":"get_weather"placeholder`),
		Messages:   []AnthropicMessage{{Role: "user", Content: json.RawMessage(`"hi"`)placeholderplaceholder,
placeholder

	out, err := AnthropicToChatCompletionsRequest(req)
placeholder
	var tc map[string]any
	require.NoError(t, json.Unmarshal(out.ToolChoice, &tc))
	require.Equal(t, "function", tc["type"])
	fn := tc["function"].(map[string]any)
	require.Equal(t, "get_weather", fn["name"])
placeholder

func TestAnthropicToChatCompletionsRequest_TemperatureStrippedForReasoningModel(t *testing.T) {
	temp := 0.7
	topP := 0.9
	req := &AnthropicRequest{
		Model:       "gpt-5.4",
		MaxTokens:   100,
		Temperature: &temp,
		TopP:        &topP,
		Messages:    []AnthropicMessage{{Role: "user", Content: json.RawMessage(`"hi"`)placeholderplaceholder,
placeholder

	out, err := AnthropicToChatCompletionsRequest(req)
placeholder
	require.Nil(t, out.Temperature, "temperature should be stripped for reasoning models")
	require.Nil(t, out.TopP, "top_p should be stripped for reasoning models")
placeholder

func TestAnthropicToChatCompletionsRequest_TemperaturePreservedForNonReasoningModel(t *testing.T) {
	temp := 0.7
	topP := 0.9
	req := &AnthropicRequest{
		Model:       "deepseek-v4-pro",
		MaxTokens:   100,
		Temperature: &temp,
		TopP:        &topP,
		Messages:    []AnthropicMessage{{Role: "user", Content: json.RawMessage(`"hi"`)placeholderplaceholder,
placeholder

	out, err := AnthropicToChatCompletionsRequest(req)
placeholder
	require.NotNil(t, out.Temperature)
	require.Equal(t, 0.7, *out.Temperature)
	require.NotNil(t, out.TopP)
	require.Equal(t, 0.9, *out.TopP)
placeholder

func TestAnthropicToChatCompletionsRequest_MaxTokensFloor(t *testing.T) {
	req := &AnthropicRequest{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 10, // below minMaxOutputTokens (128)
		Messages:  []AnthropicMessage{{Role: "user", Content: json.RawMessage(`"hi"`)placeholderplaceholder,
placeholder

	out, err := AnthropicToChatCompletionsRequest(req)
placeholder
	require.NotNil(t, out.MaxCompletionTokens)
	require.Equal(t, minMaxOutputTokens, *out.MaxCompletionTokens)
placeholder

func TestAnthropicToChatCompletionsRequest_ReasoningEffortMapping(t *testing.T) {
	req := &AnthropicRequest{
		Model:         "gpt-5.4",
		MaxTokens:     100,
		OutputConfig:  &AnthropicOutputConfig{Effort: "max"placeholder,
		Messages:      []AnthropicMessage{{Role: "user", Content: json.RawMessage(`"hi"`)placeholderplaceholder,
placeholder

	out, err := AnthropicToChatCompletionsRequest(req)
placeholder
	require.Equal(t, "xhigh", out.ReasoningEffort)
placeholder

func TestAnthropicToChatCompletionsRequest_ReasoningEffortDefaultMedium(t *testing.T) {
	req := &AnthropicRequest{
		Model:     "gpt-5.4",
		MaxTokens: 100,
		Messages:  []AnthropicMessage{{Role: "user", Content: json.RawMessage(`"hi"`)placeholderplaceholder,
placeholder

	out, err := AnthropicToChatCompletionsRequest(req)
placeholder
	require.Equal(t, "medium", out.ReasoningEffort)
placeholder

func TestAnthropicToChatCompletionsRequest_ServerToolDropped(t *testing.T) {
	req := &AnthropicRequest{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 100,
		Tools: []AnthropicTool{
			{Type: "web_search_20250305", Name: "web_search"placeholder,
			{Name: "get_weather", InputSchema: json.RawMessage(`{"type":"object","properties":{placeholderplaceholder`)placeholder,
	placeholder,
		Messages: []AnthropicMessage{{Role: "user", Content: json.RawMessage(`"hi"`)placeholderplaceholder,
placeholder

	out, err := AnthropicToChatCompletionsRequest(req)
placeholder
	require.Len(t, out.Tools, 1, "web_search server tool should be dropped")
	require.Equal(t, "get_weather", out.Tools[0].Function.Name)
placeholder

// ---------------------------------------------------------------------------
// Non-streaming response: ChatCompletionsResponseToAnthropic
// ---------------------------------------------------------------------------

func TestChatCompletionsResponseToAnthropic_TextOnly(t *testing.T) {
	resp := &ChatCompletionsResponse{
		ID:    "chatcmpl-1",
		Model: "deepseek-v4-pro",
		Choices: []ChatChoice{{
			Index:        0,
			Message:      ChatMessage{Role: "assistant", Content: json.RawMessage(`"hello world"`)placeholder,
			FinishReason: "stop",
placeholder
		Usage: &ChatUsage{PromptTokens: 5, CompletionTokens: 2, TotalTokens: 7placeholder,
placeholder

	out := ChatCompletionsResponseToAnthropic(resp, "claude-sonnet-4-20250514")
	require.Equal(t, "chatcmpl-1", out.ID)
	require.Equal(t, "claude-sonnet-4-20250514", out.Model)
	require.Len(t, out.Content, 1)
	require.Equal(t, "text", out.Content[0].Type)
	require.Equal(t, "hello world", out.Content[0].Text)
	require.Equal(t, "end_turn", out.StopReason)
	require.Equal(t, 5, out.Usage.InputTokens)
	require.Equal(t, 2, out.Usage.OutputTokens)
placeholder

func TestChatCompletionsResponseToAnthropic_ToolUse(t *testing.T) {
	resp := &ChatCompletionsResponse{
		ID:    "chatcmpl-2",
		Model: "deepseek-v4-pro",
		Choices: []ChatChoice{{
			Index: 0,
			Message: ChatMessage{
				Role: "assistant",
				ToolCalls: []ChatToolCall{{
					ID:   "call_1",
					Type: "function",
					Function: ChatFunctionCall{
						Name:      "get_weather",
						Arguments: `{"city":"SF"placeholder`,
				placeholder,
		placeholder
		placeholder,
			FinishReason: "tool_calls",
placeholder
placeholder

	out := ChatCompletionsResponseToAnthropic(resp, "claude-sonnet-4-20250514")
	require.Len(t, out.Content, 1)
	require.Equal(t, "tool_use", out.Content[0].Type)
	require.Equal(t, "call_1", out.Content[0].ID)
	require.Equal(t, "get_weather", out.Content[0].Name)
	require.Equal(t, `{"city":"SF"placeholder`, string(out.Content[0].Input))
	require.Equal(t, "tool_use", out.StopReason)
placeholder

func TestChatCompletionsResponseToAnthropic_ReasoningOnlyFallback(t *testing.T) {
	resp := &ChatCompletionsResponse{
		ID:    "chatcmpl-3",
		Model: "deepseek-v4-pro",
		Choices: []ChatChoice{{
			Index: 0,
			Message: ChatMessage{
				Role:             "assistant",
				ReasoningContent: "I should think about this",
		placeholder,
			FinishReason: "stop",
placeholder
placeholder

	out := ChatCompletionsResponseToAnthropic(resp, "claude-sonnet-4-20250514")
	// thinking block + text block (fallback uses reasoning as visible text)
	require.Len(t, out.Content, 2)
	require.Equal(t, "thinking", out.Content[0].Type)
	require.Equal(t, "I should think about this", out.Content[0].Thinking)
	require.Equal(t, "text", out.Content[1].Type)
	require.Equal(t, "I should think about this", out.Content[1].Text)
placeholder

func TestChatCompletionsResponseToAnthropic_FinishReasonLength(t *testing.T) {
	resp := &ChatCompletionsResponse{
		ID:    "chatcmpl-4",
		Model: "deepseek-v4-pro",
		Choices: []ChatChoice{{
			Index:        0,
			Message:      ChatMessage{Role: "assistant", Content: json.RawMessage(`"truncated"`)placeholder,
			FinishReason: "length",
placeholder
placeholder

	out := ChatCompletionsResponseToAnthropic(resp, "claude-sonnet-4-20250514")
	require.Equal(t, "max_tokens", out.StopReason)
placeholder

func TestChatCompletionsResponseToAnthropic_EmptyChoices(t *testing.T) {
	resp := &ChatCompletionsResponse{
		ID:      "chatcmpl-5",
		Model:   "deepseek-v4-pro",
		Choices: []ChatChoice{placeholder,
placeholder

	out := ChatCompletionsResponseToAnthropic(resp, "claude-sonnet-4-20250514")
	require.Len(t, out.Content, 1)
	require.Equal(t, "text", out.Content[0].Type)
	require.Equal(t, "", out.Content[0].Text)
placeholder

func TestChatCompletionsResponseToAnthropic_CacheTokens(t *testing.T) {
	resp := &ChatCompletionsResponse{
		ID:    "chatcmpl-6",
		Model: "deepseek-v4-pro",
		Choices: []ChatChoice{{
			Index:        0,
			Message:      ChatMessage{Role: "assistant", Content: json.RawMessage(`"hi"`)placeholder,
			FinishReason: "stop",
placeholder
		Usage: &ChatUsage{
			PromptTokens:     100,
			CompletionTokens: 5,
			TotalTokens:      105,
			PromptTokensDetails: &ChatTokenDetails{
				CachedTokens:        30,
				CacheCreationTokens: 10,
		placeholder,
	placeholder,
placeholder

	out := ChatCompletionsResponseToAnthropic(resp, "claude-sonnet-4-20250514")
	// input = prompt(100) - cached(30) - cacheCreation(10) = 60
	require.Equal(t, 60, out.Usage.InputTokens)
	require.Equal(t, 5, out.Usage.OutputTokens)
	require.Equal(t, 30, out.Usage.CacheReadInputTokens)
	require.Equal(t, 10, out.Usage.CacheCreationInputTokens)
placeholder

func TestChatCompletionsResponseToAnthropic_NilResponse(t *testing.T) {
	out := ChatCompletionsResponseToAnthropic(nil, "claude-sonnet-4-20250514")
	require.Len(t, out.Content, 1)
	require.Equal(t, "text", out.Content[0].Type)
placeholder

// ---------------------------------------------------------------------------
// Streaming: ChatCompletionsChunkToAnthropicEvents
// ---------------------------------------------------------------------------

func TestChatCompletionsChunkToAnthropicEvents_TextOnly(t *testing.T) {
	events := collectAnthropicStreamEvents(t, []string{
		`{"choices":[{"index":0,"delta":{"role":"assistant","content":"hello"placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"content":" world"placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"content":""placeholder,"finish_reason":"stop"placeholder],"usage":{"prompt_tokens":5,"completion_tokens":2,"total_tokens":7placeholderplaceholder`,
placeholder)

	types := anthropicEventTypes(events)
	// message_start → content_block_start(text) → 2× content_block_delta → content_block_stop → message_delta → message_stop
	require.Equal(t, []string{
		"message_start",
		"content_block_start",
		"content_block_delta",
		"content_block_delta",
		"content_block_stop",
		"message_delta",
		"message_stop",
placeholder, types)

	// Verify deltas
	var texts []string
	for _, e := range events {
		if e.Type == "content_block_delta" && e.Delta != nil {
			texts = append(texts, e.Delta.Text)
	placeholder
placeholder
	require.Equal(t, []string{"hello", " world"placeholder, texts)

	// Verify stop reason
	for _, e := range events {
		if e.Type == "message_delta" {
			require.Equal(t, "end_turn", e.Delta.StopReason)
			require.Equal(t, 5, e.Usage.InputTokens)
			require.Equal(t, 2, e.Usage.OutputTokens)
	placeholder
placeholder
placeholder

func TestChatCompletionsChunkToAnthropicEvents_ReasoningThenContent(t *testing.T) {
	events := collectAnthropicStreamEvents(t, []string{
		`{"choices":[{"index":0,"delta":{"reasoning_content":"thinking..."placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"content":"answer"placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"content":""placeholder,"finish_reason":"stop"placeholder],"usage":{"prompt_tokens":1,"completion_tokens":2,"total_tokens":3placeholderplaceholder`,
placeholder)

	types := anthropicEventTypes(events)
	// message_start → thinking block start → thinking_delta → thinking block stop
	// → text block start → text_delta → text block stop → message_delta → message_stop
	require.Equal(t, []string{
		"message_start",
		"content_block_start",
		"content_block_delta",
		"content_block_stop",
		"content_block_start",
		"content_block_delta",
		"content_block_stop",
		"message_delta",
		"message_stop",
placeholder, types)

	// First content_block_start should be thinking, second text
	var blockTypes []string
	for _, e := range events {
		if e.Type == "content_block_start" && e.ContentBlock != nil {
			blockTypes = append(blockTypes, e.ContentBlock.Type)
	placeholder
placeholder
	require.Equal(t, []string{"thinking", "text"placeholder, blockTypes)
placeholder

func TestChatCompletionsChunkToAnthropicEvents_ToolCallAggregation(t *testing.T) {
	events := collectAnthropicStreamEvents(t, []string{
		`{"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_1","type":"function","function":{"name":"get_weather","arguments":""placeholderplaceholder]placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"{\"city\":"placeholderplaceholder]placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"\"SF\"placeholder"placeholderplaceholder]placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{placeholder,"finish_reason":"tool_calls"placeholder],"usage":{"prompt_tokens":10,"completion_tokens":5,"total_tokens":15placeholderplaceholder`,
placeholder)

	types := anthropicEventTypes(events)
	// message_start → content_block_start(tool_use) → 2× input_json_delta (empty first arg skipped) → content_block_stop → message_delta(tool_use) → message_stop
	require.Equal(t, []string{
		"message_start",
		"content_block_start",
		"content_block_delta",
		"content_block_delta",
		"content_block_stop",
		"message_delta",
		"message_stop",
placeholder, types)

	// Verify tool_use block
	for _, e := range events {
		if e.Type == "content_block_start" && e.ContentBlock != nil {
			require.Equal(t, "tool_use", e.ContentBlock.Type)
			require.Equal(t, "call_1", e.ContentBlock.ID)
			require.Equal(t, "get_weather", e.ContentBlock.Name)
	placeholder
		if e.Type == "message_delta" {
			require.Equal(t, "tool_use", e.Delta.StopReason)
	placeholder
placeholder

	// Verify arguments assembled (empty first fragment skipped)
	var partials []string
	for _, e := range events {
		if e.Type == "content_block_delta" && e.Delta != nil {
			partials = append(partials, e.Delta.PartialJSON)
	placeholder
placeholder
	require.Equal(t, []string{`{"city":`, `"SF"placeholder`placeholder, partials)
placeholder

func TestChatCompletionsChunkToAnthropicEvents_LengthMapsToMaxTokens(t *testing.T) {
	events := collectAnthropicStreamEvents(t, []string{
		`{"choices":[{"index":0,"delta":{"content":"partial"placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"content":""placeholder,"finish_reason":"length"placeholder],"usage":{"prompt_tokens":1,"completion_tokens":2,"total_tokens":3placeholderplaceholder`,
placeholder)

	for _, e := range events {
		if e.Type == "message_delta" {
			require.Equal(t, "max_tokens", e.Delta.StopReason)
	placeholder
placeholder
placeholder

func TestChatCompletionsChunkToAnthropicEvents_EmptyStream(t *testing.T) {
	events := collectAnthropicStreamEvents(t, []string{
		`{"choices":[{"index":0,"delta":{placeholder,"finish_reason":"stop"placeholder],"usage":{"prompt_tokens":1,"completion_tokens":0,"total_tokens":1placeholderplaceholder`,
placeholder)

	types := anthropicEventTypes(events)
	// Even with no content, message_start + message_delta + message_stop should fire.
	require.Contains(t, types, "message_start")
	require.Contains(t, types, "message_stop")
placeholder

func TestChatCompletionsChunkToAnthropicEvents_MessageStartEmittedOnce(t *testing.T) {
	events := collectAnthropicStreamEvents(t, []string{
		`{"choices":[{"index":0,"delta":{"content":"a"placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"content":"b"placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"content":""placeholder,"finish_reason":"stop"placeholder],"usage":{"prompt_tokens":1,"completion_tokens":2,"total_tokens":3placeholderplaceholder`,
placeholder)

	count := 0
	for _, e := range events {
		if e.Type == "message_start" {
			count++
	placeholder
placeholder
	require.Equal(t, 1, count, "message_start should only be emitted once")
placeholder

func TestChatCompletionsChunkToAnthropicEvents_ParallelToolCalls(t *testing.T) {
	events := collectAnthropicStreamEvents(t, []string{
		`{"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_1","type":"function","function":{"name":"tool_a","arguments":"{placeholder"placeholderplaceholder]placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"tool_calls":[{"index":1,"id":"call_2","type":"function","function":{"name":"tool_b","arguments":"{placeholder"placeholderplaceholder]placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{placeholder,"finish_reason":"tool_calls"placeholder],"usage":{"prompt_tokens":1,"completion_tokens":2,"total_tokens":3placeholderplaceholder`,
placeholder)

	// Two tool_use blocks should be opened
	var toolBlocks []string
	for _, e := range events {
		if e.Type == "content_block_start" && e.ContentBlock != nil && e.ContentBlock.Type == "tool_use" {
			toolBlocks = append(toolBlocks, e.ContentBlock.Name)
	placeholder
placeholder
	require.Equal(t, []string{"tool_a", "tool_b"placeholder, toolBlocks)

	for _, e := range events {
		if e.Type == "message_delta" {
			require.Equal(t, "tool_use", e.Delta.StopReason)
	placeholder
placeholder
placeholder

func TestFinalizeChatCompletionsAnthropicStream_NoOpAfterStop(t *testing.T) {
	state := NewChatCompletionsToAnthropicStreamState("test")
	state.MessageStopSent = true

	events := FinalizeChatCompletionsAnthropicStream(state)
	require.Nil(t, events, "finalize should be a no-op after message_stop")
placeholder

func TestFinalizeChatCompletionsAnthropicStream_EmitsMessageStartIfMissing(t *testing.T) {
	state := NewChatCompletionsToAnthropicStreamState("test")
	// Never fed any chunks — message_start not yet sent

	events := FinalizeChatCompletionsAnthropicStream(state)
	types := anthropicEventTypes(events)
	require.Contains(t, types, "message_start")
	require.Contains(t, types, "message_stop")
placeholder

// ---------------------------------------------------------------------------
// Equivalence: direct bridge matches the double-conversion bridge
// ---------------------------------------------------------------------------

// TestDirectBridge_NonStreamingMatchesDoubleConversion verifies that
// ChatCompletionsResponseToAnthropic produces the same Anthropic response as the
// existing ChatCompletionsResponseToResponses + ResponsesToAnthropic chain.
func TestDirectBridge_NonStreamingMatchesDoubleConversion(t *testing.T) {
	resp := &ChatCompletionsResponse{
		ID:    "chatcmpl-eq",
		Model: "deepseek-v4-pro",
		Choices: []ChatChoice{{
			Index: 0,
			Message: ChatMessage{
				Role:             "assistant",
				Content:          json.RawMessage(`"hello"`),
				ReasoningContent: "reasoning text",
				ToolCalls: []ChatToolCall{{
					ID:   "call_eq",
					Type: "function",
					Function: ChatFunctionCall{
						Name:      "search",
						Arguments: `{"q":"test"placeholder`,
				placeholder,
		placeholder
		placeholder,
			FinishReason: "tool_calls",
placeholder
		Usage: &ChatUsage{
			PromptTokens:        50,
			CompletionTokens:    10,
			TotalTokens:         60,
			PromptTokensDetails: &ChatTokenDetails{CachedTokens: 5placeholder,
	placeholder,
placeholder

	// Direct bridge
	direct := ChatCompletionsResponseToAnthropic(resp, "claude-sonnet-4-20250514")

	// Double-conversion bridge
	responsesResp := ChatCompletionsResponseToResponses(resp, "claude-sonnet-4-20250514", nil, false, nil)
	double := ResponsesToAnthropic(responsesResp, "claude-sonnet-4-20250514")

	// Compare key fields
	require.Equal(t, direct.StopReason, double.StopReason)
	require.Equal(t, direct.Model, double.Model)
	require.Len(t, direct.Content, len(double.Content))
	for i := range direct.Content {
		require.Equal(t, double.Content[i].Type, direct.Content[i].Type, "block %d type mismatch", i)
		require.Equal(t, double.Content[i].Text, direct.Content[i].Text, "block %d text mismatch", i)
		require.Equal(t, double.Content[i].Thinking, direct.Content[i].Thinking, "block %d thinking mismatch", i)
		require.Equal(t, double.Content[i].Name, direct.Content[i].Name, "block %d name mismatch", i)
		require.Equal(t, double.Content[i].ID, direct.Content[i].ID, "block %d id mismatch", i)
placeholder
	require.Equal(t, double.Usage.InputTokens, direct.Usage.InputTokens)
	require.Equal(t, double.Usage.OutputTokens, direct.Usage.OutputTokens)
	require.Equal(t, double.Usage.CacheReadInputTokens, direct.Usage.CacheReadInputTokens)
	require.Equal(t, double.Usage.CacheCreationInputTokens, direct.Usage.CacheCreationInputTokens)
placeholder

// TestDirectBridge_RequestMatchesDoubleConversion verifies that
// AnthropicToChatCompletionsRequest produces an equivalent Chat Completions
// request as the AnthropicToResponses + ResponsesToChatCompletionsRequest chain.
func TestDirectBridge_RequestMatchesDoubleConversion(t *testing.T) {
	temp := 0.5
	req := &AnthropicRequest{
		Model:       "deepseek-v4-pro",
		MaxTokens:   500,
		Temperature: &temp,
		System:      json.RawMessage(`"be helpful"`),
		Tools: []AnthropicTool{
			{Name: "get_weather", InputSchema: json.RawMessage(`{"type":"object","properties":{"city":{"type":"string"placeholderplaceholderplaceholder`)placeholder,
	placeholder,
		ToolChoice: json.RawMessage(`{"type":"auto"placeholder`),
		Messages: []AnthropicMessage{
			{Role: "user", Content: json.RawMessage(`"what's the weather?"`)placeholder,
			{Role: "assistant", Content: json.RawMessage(`[{"type":"text","text":"checking"placeholder,{"type":"tool_use","id":"toolu_1","name":"get_weather","input":{"city":"SF"placeholderplaceholder]`)placeholder,
			{Role: "user", Content: json.RawMessage(`[{"type":"tool_result","tool_use_id":"toolu_1","content":"sunny"placeholder]`)placeholder,
	placeholder,
placeholder

	// Direct bridge
	direct, err := AnthropicToChatCompletionsRequest(req)
placeholder

	// Double-conversion bridge
	responsesReq, err := AnthropicToResponses(req)
placeholder
	double, err := ResponsesToChatCompletionsRequest(responsesReq)
placeholder

	// Compare key fields
	require.Equal(t, double.Model, direct.Model)
	require.Equal(t, double.Temperature, direct.Temperature)
	require.Equal(t, double.MaxCompletionTokens, direct.MaxCompletionTokens)
	require.Equal(t, double.ReasoningEffort, direct.ReasoningEffort)
	require.Equal(t, string(double.ToolChoice), string(direct.ToolChoice))
	require.Len(t, direct.Tools, len(double.Tools))

	// Compare messages — same count, same roles, same content
	require.Len(t, direct.Messages, len(double.Messages), "message count mismatch")
	for i := range direct.Messages {
		require.Equal(t, double.Messages[i].Role, direct.Messages[i].Role, "msg %d role mismatch", i)
		// Normalize content for comparison (both should be valid JSON)
		var dContent, dblContent any
		_ = json.Unmarshal(double.Messages[i].Content, &dblContent)
		_ = json.Unmarshal(direct.Messages[i].Content, &dContent)
		require.Equal(t, dblContent, dContent, "msg %d content mismatch", i)
		require.Equal(t, double.Messages[i].ToolCallID, direct.Messages[i].ToolCallID, "msg %d tool_call_id mismatch", i)
		require.Len(t, direct.Messages[i].ToolCalls, len(double.Messages[i].ToolCalls), "msg %d tool_calls count mismatch", i)
		for j := range direct.Messages[i].ToolCalls {
			require.Equal(t, double.Messages[i].ToolCalls[j].ID, direct.Messages[i].ToolCalls[j].ID, "msg %d tool %d id mismatch", i, j)
			require.Equal(t, double.Messages[i].ToolCalls[j].Function.Name, direct.Messages[i].ToolCalls[j].Function.Name, "msg %d tool %d name mismatch", i, j)
	placeholder
placeholder
placeholder

// ---------------------------------------------------------------------------
// Edge cases
// ---------------------------------------------------------------------------

func TestChatCompletionsChunkToAnthropicEvents_ImageInToolResult(t *testing.T) {
	// A multi-turn conversation: assistant calls a tool, user replies with a
	// tool_result containing text + an image. The image should be lifted into
	// a follow-up user message as an image_url part.
	req := &AnthropicRequest{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 100,
		Messages: []AnthropicMessage{
			{Role: "user", Content: json.RawMessage(`"check this image"`)placeholder,
			{Role: "assistant", Content: json.RawMessage(`[{"type":"text","text":"let me look"placeholder,{"type":"tool_use","id":"toolu_1","name":"analyze","input":{"x":1placeholderplaceholder]`)placeholder,
			{Role: "user", Content: json.RawMessage(`[{"type":"tool_result","tool_use_id":"toolu_1","content":[{"type":"text","text":"result"placeholder,{"type":"image","source":{"type":"base64","media_type":"image/png","data":"aGVsbG8="placeholderplaceholder]placeholder]`)placeholder,
	placeholder,
placeholder

	out, err := AnthropicToChatCompletionsRequest(req)
placeholder
	// user + assistant(tool_use) + tool + user(image)
	require.GreaterOrEqual(t, len(out.Messages), 3)

	// Find the user message with image content
	var foundImage bool
	for _, m := range out.Messages {
		if m.Role != "user" {
			continue
	placeholder
		var parts []ChatContentPart
		if err := json.Unmarshal(m.Content, &parts); err == nil {
			for _, p := range parts {
				if p.Type == "image_url" && p.ImageURL != nil {
					foundImage = true
					require.True(t, strings.HasPrefix(p.ImageURL.URL, "data:image/png;base64,"))
			placeholder
		placeholder
	placeholder
placeholder
	require.True(t, foundImage, "image from tool_result should appear in user message")
placeholder

func TestChatCompletionsToAnthropicStreamState_ToolCallNameArrivesLate(t *testing.T) {
	// Some upstreams send the tool_call index + arguments before the name.
	events := collectAnthropicStreamEvents(t, []string{
		`{"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_late"placeholder]placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"name":"late_tool"placeholderplaceholder]placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"{placeholder"placeholderplaceholder]placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{placeholder,"finish_reason":"tool_calls"placeholder],"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2placeholderplaceholder`,
placeholder)

	// The tool_use block should still be opened with the correct name.
	var toolName string
	for _, e := range events {
		if e.Type == "content_block_start" && e.ContentBlock != nil && e.ContentBlock.Type == "tool_use" {
			toolName = e.ContentBlock.Name
	placeholder
placeholder
	require.Equal(t, "late_tool", toolName)
placeholder
