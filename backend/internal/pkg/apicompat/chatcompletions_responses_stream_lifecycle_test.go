package apicompat

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func collectStreamEvents(t *testing.T, chunks []string) []ResponsesStreamEvent {
placeholder
	state := NewChatCompletionsToResponsesStreamState("deepseek-v4-pro")
	var events []ResponsesStreamEvent
	for _, payload := range chunks {
		var chunk ChatCompletionsChunk
		require.NoError(t, json.Unmarshal([]byte(payload), &chunk))
		events = append(events, ChatCompletionsChunkToResponsesEvents(&chunk, state)...)
placeholder
	events = append(events, FinalizeChatCompletionsResponsesStream(state)...)
	return events
placeholder

// TestStream_ReasoningOpensItemBeforeDelta guards the bug where a strict client
// (Codex) drops reasoning deltas that reference an item not yet opened.
func TestStream_ReasoningOpensItemBeforeDelta(t *testing.T) {
	events := collectStreamEvents(t, []string{
		`{"choices":[{"index":0,"delta":{"role":"assistant","content":null,"reasoning_content":""placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"reasoning_content":"think"placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"content":"hello"placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"content":""placeholder,"finish_reason":"stop"placeholder],"usage":{"prompt_tokens":1,"completion_tokens":2,"total_tokens":3placeholderplaceholder`,
placeholder)

	open := map[int]string{placeholder // output_index -> item type
	for _, e := range events {
		switch e.Type {
		case "response.output_item.added":
			require.NotNil(t, e.Item)
			open[e.OutputIndex] = e.Item.Type
		case "response.reasoning_summary_text.delta":
			require.Equalf(t, "reasoning", open[e.OutputIndex], "reasoning delta before its item was opened")
		case "response.output_text.delta":
			require.Equalf(t, "message", open[e.OutputIndex], "text delta before its item was opened")
	placeholder
placeholder
placeholder

func TestStream_ReasoningOnlySynthesizesVisibleText(t *testing.T) {
	events := collectStreamEvents(t, []string{
		`{"choices":[{"index":0,"delta":{"role":"assistant","content":null,"reasoning_content":""placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"reasoning_content":"thinking before final"placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"content":""placeholder,"finish_reason":"length"placeholder],"usage":{"prompt_tokens":1,"completion_tokens":2,"total_tokens":3placeholderplaceholder`,
placeholder)

	open := map[int]string{placeholder
	var sawTextDelta, sawTextDone, sawMessageDone bool
	for _, e := range events {
		switch e.Type {
		case "response.output_item.added":
			require.NotNil(t, e.Item)
			open[e.OutputIndex] = e.Item.Type
		case "response.output_text.delta":
			sawTextDelta = true
			require.Equalf(t, "message", open[e.OutputIndex], "fallback text delta before its item was opened")
			require.Equal(t, "thinking before final", e.Delta)
		case "response.output_text.done":
			sawTextDone = true
			require.Equal(t, "thinking before final", e.Text)
		case "response.output_item.done":
			if e.Item != nil && e.Item.Type == "message" {
				sawMessageDone = true
				require.Equal(t, "thinking before final", e.Item.Content[0].Text)
		placeholder
		case "response.completed":
			require.NotNil(t, e.Response)
			require.Equal(t, "incomplete", e.Response.Status)
			require.NotNil(t, e.Response.IncompleteDetails)
			require.Equal(t, "max_output_tokens", e.Response.IncompleteDetails.Reason)
			require.Len(t, e.Response.Output, 2)
			require.Equal(t, "reasoning", e.Response.Output[0].Type)
			require.Equal(t, "message", e.Response.Output[1].Type)
			require.Equal(t, "thinking before final", e.Response.Output[1].Content[0].Text)
	placeholder
placeholder
	require.True(t, sawTextDelta, "reasoning-only stream must produce visible text delta")
	require.True(t, sawTextDone, "reasoning-only stream must close visible text part")
	require.True(t, sawMessageDone, "reasoning-only stream must close synthesized message item")
placeholder

func TestStream_ReasoningOnlyBlankDoesNotSynthesizeVisibleText(t *testing.T) {
	events := collectStreamEvents(t, []string{
		`{"choices":[{"index":0,"delta":{"reasoning_content":"   "placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{placeholder,"finish_reason":"stop"placeholder]placeholder`,
placeholder)

	for _, e := range events {
		require.NotEqual(t, "response.output_text.delta", e.Type)
		if e.Type == "response.completed" {
			require.NotNil(t, e.Response)
			require.Len(t, e.Response.Output, 2)
			require.Equal(t, "reasoning", e.Response.Output[0].Type)
			require.Equal(t, "message", e.Response.Output[1].Type)
			require.Equal(t, "", e.Response.Output[1].Content[0].Text)
	placeholder
placeholder
placeholder

func TestStream_ReasoningThenContentDoesNotDuplicateFallbackText(t *testing.T) {
	events := collectStreamEvents(t, []string{
		`{"choices":[{"index":0,"delta":{"reasoning_content":"private plan"placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"content":"final answer"placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{placeholder,"finish_reason":"stop"placeholder]placeholder`,
placeholder)

	var textDeltas []string
	for _, e := range events {
		switch e.Type {
		case "response.output_text.delta":
			textDeltas = append(textDeltas, e.Delta)
		case "response.completed":
			require.NotNil(t, e.Response)
			require.Len(t, e.Response.Output, 2)
			require.Equal(t, "private plan", e.Response.Output[0].Summary[0].Text)
			require.Equal(t, "final answer", e.Response.Output[1].Content[0].Text)
	placeholder
placeholder
	require.Equal(t, []string{"final answer"placeholder, textDeltas)
placeholder

func TestStream_ReasoningThenToolCallDoesNotSynthesizeVisibleText(t *testing.T) {
	events := collectStreamEvents(t, []string{
		`{"choices":[{"index":0,"delta":{"reasoning_content":"call a tool"placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_a","type":"function","function":{"name":"exec","arguments":"{placeholder"placeholderplaceholder]placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{placeholder,"finish_reason":"tool_calls"placeholder]placeholder`,
placeholder)

	for _, e := range events {
		require.NotEqual(t, "response.output_text.delta", e.Type)
		if e.Type == "response.completed" {
			require.NotNil(t, e.Response)
			require.Len(t, e.Response.Output, 2)
			require.Equal(t, "reasoning", e.Response.Output[0].Type)
			require.Equal(t, "function_call", e.Response.Output[1].Type)
	placeholder
placeholder
placeholder

// TestStream_ToolCallLifecycleComplete guards that a tool call is fully closed
// (function_call_arguments.done + output_item.done with full arguments), which
// codex needs to execute the call.
func TestStream_ToolCallLifecycleComplete(t *testing.T) {
	events := collectStreamEvents(t, []string{
		`{"choices":[{"index":0,"delta":{"role":"assistant","reasoning_content":"plan"placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_a","type":"function","function":{"name":"exec","arguments":""placeholderplaceholder]placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"{\"cmd\":\"ls\"placeholder"placeholderplaceholder]placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{placeholder,"finish_reason":"tool_calls"placeholder],"usage":{"prompt_tokens":1,"completion_tokens":2,"total_tokens":3placeholderplaceholder`,
placeholder)

	var sawAdded, sawArgsDone, sawItemDone bool
	for _, e := range events {
		switch e.Type {
		case "response.output_item.added":
			if e.Item != nil && e.Item.Type == "function_call" {
				sawAdded = true
		placeholder
		case "response.function_call_arguments.done":
			sawArgsDone = true
			require.Equal(t, `{"cmd":"ls"placeholder`, e.Arguments)
		case "response.output_item.done":
			if e.Item != nil && e.Item.Type == "function_call" {
				sawItemDone = true
				require.Equal(t, `{"cmd":"ls"placeholder`, e.Item.Arguments)
				require.Equal(t, "call_a", e.Item.CallID)
		placeholder
	placeholder
placeholder
	require.True(t, sawAdded, "function_call output_item.added missing")
	require.True(t, sawArgsDone, "function_call_arguments.done missing")
	require.True(t, sawItemDone, "function_call output_item.done missing")
placeholder

// TestStream_ToolCallArgumentsInFirstChunkNotDoubled guards the GLM/Zhipu shape
// where a single tool_call delta chunk carries id+name+arguments together.
// Earlier code copied the whole tool_call (including arguments) into state and
// then accumulated the same chunk's arguments again, producing a doubled,
// invalid JSON like {"cmd":"ls"placeholder{"cmd":"ls"placeholder that breaks Codex tool parsing
// ("trailing characters").
func TestStream_ToolCallArgumentsInFirstChunkNotDoubled(t *testing.T) {
	events := collectStreamEvents(t, []string{
		`{"choices":[{"index":0,"delta":{"role":"assistant"placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_a","type":"function","function":{"name":"exec","arguments":"{\"cmd\":\"ls\"placeholder"placeholderplaceholder]placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{placeholder,"finish_reason":"tool_calls"placeholder]placeholder`,
placeholder)

	var argsDelta strings.Builder
	var sawArgsDone, sawItemDone bool
	for _, e := range events {
		switch e.Type {
		case "response.function_call_arguments.delta":
			_, _ = argsDelta.WriteString(e.Delta)
		case "response.function_call_arguments.done":
			sawArgsDone = true
			require.Equal(t, `{"cmd":"ls"placeholder`, e.Arguments)
		case "response.output_item.done":
			if e.Item != nil && e.Item.Type == "function_call" {
				sawItemDone = true
				require.Equal(t, `{"cmd":"ls"placeholder`, e.Item.Arguments)
		placeholder
	placeholder
placeholder
	require.True(t, sawArgsDone, "function_call_arguments.done missing")
	require.True(t, sawItemDone, "function_call output_item.done missing")
	// Accumulated deltas must equal the final arguments exactly (no duplication).
	require.Equal(t, `{"cmd":"ls"placeholder`, argsDelta.String())
placeholder

func TestStream_InvalidToolArgumentsAreRejectedBeforeFinalize(t *testing.T) {
	idx := 0
	state := NewChatCompletionsToResponsesStreamState("deepseek-v4-flash")
	chunk := &ChatCompletionsChunk{
		Choices: []ChatChunkChoice{
			{
				Index: 0,
				Delta: ChatDelta{
					ToolCalls: []ChatToolCall{
						{
							Index: &idx,
							ID:    "call_bad",
							Type:  "function",
							Function: ChatFunctionCall{
								Name:      "exec_command",
								Arguments: `{"cmd": "ssh root@HOST`,
						placeholder,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
	placeholder,
placeholder
	ChatCompletionsChunkToResponsesEvents(chunk, state)

	err := state.ValidateToolCallArguments()
	require.ErrorContains(t, err, "invalid JSON")
placeholder

func TestStream_ValidToolCallAtOutputLimitKeepsIncompleteResponse(t *testing.T) {
	idx := 0
	state := NewChatCompletionsToResponsesStreamState("deepseek-v4-flash")
	chunk := &ChatCompletionsChunk{
		Choices: []ChatChunkChoice{
			{
				Index: 0,
				Delta: ChatDelta{
					ToolCalls: []ChatToolCall{
						{
							Index: &idx,
							ID:    "call_at_limit",
							Type:  "function",
							Function: ChatFunctionCall{
								Name:      "exec_command",
								Arguments: `{placeholder`,
						placeholder,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
	placeholder,
placeholder
	ChatCompletionsChunkToResponsesEvents(chunk, state)
	state.FinishReason = "length"

	require.NoError(t, state.ValidateToolCallArguments())
	events := FinalizeChatCompletionsResponsesStream(state)
	var sawArgsDone, sawIncomplete bool
	for _, event := range events {
		switch event.Type {
		case "response.function_call_arguments.done":
			sawArgsDone = true
			require.Equal(t, `{placeholder`, event.Arguments)
		case "response.completed":
			require.NotNil(t, event.Response)
			sawIncomplete = event.Response.Status == "incomplete"
	placeholder
placeholder
	require.True(t, sawArgsDone)
	require.True(t, sawIncomplete)
placeholder

// TestStream_SSEWireComplete drives the full stream through SSE encoding and
// asserts the function_call events carry complete fields on the wire.
func TestStream_SSEWireComplete(t *testing.T) {
	events := collectStreamEvents(t, []string{
		`{"choices":[{"index":0,"delta":{"role":"assistant","reasoning_content":"plan"placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_a","type":"function","function":{"name":"exec","arguments":"{placeholder"placeholderplaceholder]placeholderplaceholder]placeholder`,
		`{"choices":[{"index":0,"delta":{placeholder,"finish_reason":"tool_calls"placeholder]placeholder`,
placeholder)

	var addedLine string
	for _, e := range events {
		sse, err := ResponsesEventToSSE(e)
	placeholder
		if e.Type == "response.output_item.added" && e.Item != nil && e.Item.Type == "function_call" {
			addedLine = sse
	placeholder
placeholder
	require.NotEmpty(t, addedLine)
	// The function_call added event must carry arguments:"" on the wire.
	require.True(t, strings.Contains(addedLine, `"arguments":""`), "added line missing arguments: %s", addedLine)
	require.Contains(t, addedLine, `"call_id":"call_a"`)
placeholder
