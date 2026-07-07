package apicompat

// custom/freeform 工具（如 Codex 0.14x 的 exec）在 responses→chat 桥上的双向转换。
// 背景：Codex 的核心命令执行工具 exec 是 type=custom（输入为自由文本），此前被
// responsesToolsToChatTools 丢弃，导致模型工具列表中没有 exec、无法执行任何命令。

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponsesToChatCompletionsRequest_CustomToolBecomesFunctionTool(t *testing.T) {
	req := &ResponsesRequest{
		Model: "glm-5.2",
		Input: json.RawMessage(`"run dir"`),
		Tools: []ResponsesTool{
			{Type: "custom", Name: "exec", Description: "Run JavaScript code"placeholder,
			{Type: "function", Name: "wait", Parameters: json.RawMessage(`{"type":"object","properties":{placeholderplaceholder`)placeholder,
	placeholder,
placeholder

	out, err := ResponsesToChatCompletionsRequest(req)
placeholder
	require.Len(t, out.Tools, 2)

	assert.Equal(t, "function", out.Tools[0].Type)
	assert.Equal(t, "exec", out.Tools[0].Function.Name)
	assert.Equal(t, "Run JavaScript code", out.Tools[0].Function.Description)
	assert.JSONEq(t, customToolInputSchema, string(out.Tools[0].Function.Parameters))

	assert.Equal(t, "wait", out.Tools[1].Function.Name)
placeholder

func TestResponsesToChatCompletionsRequest_DropsToolChoiceWhenNoConvertibleTools(t *testing.T) {
	req := &ResponsesRequest{
		Model: "glm-5.2",
		Input: json.RawMessage(`"hi"`),
		Tools: []ResponsesTool{
			{Type: "web_search"placeholder,
			{Type: "image_generation"placeholder,
	placeholder,
		ToolChoice: json.RawMessage(`"auto"`),
placeholder

	out, err := ResponsesToChatCompletionsRequest(req)
placeholder

	assert.Empty(t, out.Tools)
	assert.Empty(t, out.ToolChoice, "tools 为空时转发 tool_choice 会被上游 400 拒绝")
placeholder

func TestResponsesToChatCompletionsRequest_CustomToolChoiceMapsToFunctionChoice(t *testing.T) {
	req := &ResponsesRequest{
		Model:      "glm-5.2",
		Input:      json.RawMessage(`"run dir"`),
		Tools:      []ResponsesTool{{Type: "custom", Name: "exec"placeholderplaceholder,
		ToolChoice: json.RawMessage(`{"type":"custom","name":"exec"placeholder`),
placeholder

	out, err := ResponsesToChatCompletionsRequest(req)
placeholder

	assert.JSONEq(t, `{"type":"function","function":{"name":"exec"placeholderplaceholder`, string(out.ToolChoice))
placeholder

func TestResponsesInputToChatMessages_CustomToolCallHistory(t *testing.T) {
	input := json.RawMessage(`[
		{"role":"user","content":"list files"placeholder,
		{"type":"custom_tool_call","call_id":"call_1","name":"exec","input":"dir"placeholder,
		{"type":"custom_tool_call_output","call_id":"call_1","output":"main.go"placeholder
	]`)

	messages, err := responsesInputToChatMessages("", input)
placeholder
	require.Len(t, messages, 3)

	assert.Equal(t, []string{"user", "assistant", "tool"placeholder, chatMessageRoles(messages))

	require.Len(t, messages[1].ToolCalls, 1)
	toolCall := messages[1].ToolCalls[0]
	assert.Equal(t, "call_1", toolCall.ID)
	assert.Equal(t, "exec", toolCall.Function.Name)
	assert.JSONEq(t, `{"input":"dir"placeholder`, toolCall.Function.Arguments)

	assert.Equal(t, "call_1", messages[2].ToolCallID)
	assert.JSONEq(t, `"main.go"`, string(messages[2].Content))
placeholder

func TestChatCompletionsResponseToResponses_CustomToolCallOutputItem(t *testing.T) {
	resp := &ChatCompletionsResponse{
		ID: "cc-1",
		Choices: []ChatChoice{{
			Message: ChatMessage{
				Role: "assistant",
				ToolCalls: []ChatToolCall{
					{ID: "call_1", Function: ChatFunctionCall{Name: "exec", Arguments: `{"input": "dir"placeholder`placeholderplaceholder,
					{ID: "call_2", Function: ChatFunctionCall{Name: "wait", Arguments: `{"cell_id": 3placeholder`placeholderplaceholder,
			placeholder,
		placeholder,
placeholder
placeholder

	out := ChatCompletionsResponseToResponses(resp, "glm-5.2", map[string]bool{"exec": trueplaceholder)
	require.Len(t, out.Output, 2)

	assert.Equal(t, "custom_tool_call", out.Output[0].Type)
	assert.Equal(t, "call_1", out.Output[0].CallID)
	assert.Equal(t, "exec", out.Output[0].Name)
	assert.Equal(t, "dir", out.Output[0].Input)
	assert.Empty(t, out.Output[0].Arguments)

	assert.Equal(t, "function_call", out.Output[1].Type)
	assert.Equal(t, "wait", out.Output[1].Name)
	assert.Equal(t, `{"cell_id": 3placeholder`, out.Output[1].Arguments)
placeholder

func TestExtractCustomToolCallInput_FallsBackToRawArguments(t *testing.T) {
	assert.Equal(t, "dir", extractCustomToolCallInput(`{"input": "dir"placeholder`))
	assert.Equal(t, "console.log(1)", extractCustomToolCallInput(`console.log(1)`))
	assert.Equal(t, `{"other": "x"placeholder`, extractCustomToolCallInput(`{"other": "x"placeholder`))
	assert.Equal(t, "", extractCustomToolCallInput(`{placeholder`))
	assert.Equal(t, "", extractCustomToolCallInput(""))
placeholder

func TestChatCompletionsChunkToResponsesEvents_CustomToolCallStream(t *testing.T) {
	state := NewChatCompletionsToResponsesStreamState("glm-5.2")
	state.CustomTools = map[string]bool{"exec": trueplaceholder

	idx := 0
	chunk := &ChatCompletionsChunk{
		ID: "cc-1",
		Choices: []ChatChunkChoice{{
			Delta: ChatDelta{
				ToolCalls: []ChatToolCall{{
					Index:    &idx,
					ID:       "call_1",
					Function: ChatFunctionCall{Name: "exec", Arguments: `{"input": "dir"placeholder`placeholder,
		placeholder
		placeholder,
placeholder
placeholder

	events := ChatCompletionsChunkToResponsesEvents(chunk, state)
	events = append(events, FinalizeChatCompletionsResponsesStream(state)...)

	var added, inputDone, itemDone *ResponsesStreamEvent
	for i := range events {
		evt := &events[i]
		switch evt.Type {
		case "response.output_item.added":
			if evt.Item != nil && evt.Item.Type != "message" && evt.Item.Type != "reasoning" {
				added = evt
		placeholder
		case "response.custom_tool_call_input.done":
			inputDone = evt
		case "response.output_item.done":
			if evt.Item != nil && evt.Item.Type == "custom_tool_call" {
				itemDone = evt
		placeholder
		case "response.function_call_arguments.delta", "response.function_call_arguments.done":
			t.Fatalf("custom 工具调用不应产出 function_call 参数事件: %s", evt.Type)
	placeholder
placeholder

	require.NotNil(t, added, "缺少 custom_tool_call 的 output_item.added")
	assert.Equal(t, "custom_tool_call", added.Item.Type)
	assert.Equal(t, "exec", added.Item.Name)

	require.NotNil(t, inputDone, "缺少 response.custom_tool_call_input.done")
	assert.Equal(t, "dir", inputDone.Input)
	assert.Equal(t, "call_1", inputDone.CallID)

	require.NotNil(t, itemDone, "缺少 custom_tool_call 的 output_item.done")
	assert.Equal(t, "call_1", itemDone.Item.CallID)
	assert.Equal(t, "exec", itemDone.Item.Name)
	assert.Equal(t, "dir", itemDone.Item.Input)
	assert.Empty(t, itemDone.Item.Arguments)

	// response.completed 的 output 数组同样携带 custom_tool_call 项。
	final := events[len(events)-1]
	require.Equal(t, "response.completed", final.Type)
	require.NotNil(t, final.Response)
	foundCustom := false
	for _, item := range final.Response.Output {
		if item.Type == "custom_tool_call" {
			foundCustom = true
			assert.Equal(t, "exec", item.Name)
			assert.Equal(t, "dir", item.Input)
	placeholder
placeholder
	assert.True(t, foundCustom, "response.completed 缺少 custom_tool_call 输出项")
placeholder

func TestResponsesToChatCompletionsRequest_ToolSearchToolBecomesProxyFunction(t *testing.T) {
	req := &ResponsesRequest{
		Model: "glm-5.2",
		Input: json.RawMessage(`"hi"`),
		Tools: []ResponsesTool{{Type: "tool_search"placeholderplaceholder,
placeholder

	out, err := ResponsesToChatCompletionsRequest(req)
placeholder
	require.Len(t, out.Tools, 1)

	assert.Equal(t, "function", out.Tools[0].Type)
	assert.Equal(t, "tool_search", out.Tools[0].Function.Name)
	assert.Contains(t, string(out.Tools[0].Function.Parameters), `"query"`)
placeholder

func TestResponsesToChatCompletionsRequest_NamespaceToolFlattensChildren(t *testing.T) {
	req := &ResponsesRequest{
		Model: "glm-5.2",
		Input: json.RawMessage(`"hi"`),
		Tools: []ResponsesTool{{
			Type: "namespace",
			Name: "gmail",
			Tools: []ResponsesTool{
				{Type: "function", Name: "send", Description: "Send mail", Parameters: json.RawMessage(`{"type":"object","properties":{placeholderplaceholder`)placeholder,
				{Type: "custom", Name: "ignored_child"placeholder,
		placeholder,
placeholder
placeholder

	out, err := ResponsesToChatCompletionsRequest(req)
placeholder
	require.Len(t, out.Tools, 1, "namespace 子工具中仅 function 类型被摊平")

	assert.Equal(t, "gmail__send", out.Tools[0].Function.Name)
	assert.Equal(t, "Send mail", out.Tools[0].Function.Description)
placeholder

func TestResponsesToolsParsing_StringToolBecomesCustom(t *testing.T) {
	var req ResponsesRequest
	require.NoError(t, json.Unmarshal([]byte(`{"model":"glm-5.2","input":"hi","tools":["exec",{"type":"function","name":"wait"placeholder]placeholder`), &req))

	require.Len(t, req.Tools, 2)
	assert.Equal(t, "custom", req.Tools[0].Type)
	assert.Equal(t, "exec", req.Tools[0].Name)
	assert.Equal(t, "function", req.Tools[1].Type)

	assert.True(t, CustomToolNames(req.Tools)["exec"])
placeholder

func TestFlattenNamespaceToolName_CapsAt64WithHashSuffix(t *testing.T) {
	assert.Equal(t, "gmail__send", flattenNamespaceToolName("gmail", "send"))

	long := flattenNamespaceToolName("very_long_namespace_prefix_for_testing_purposes", "and_a_rather_long_tool_name_too")
	assert.LessOrEqual(t, len(long), 64)
	assert.Contains(t, long, "__")
	// 同输入结果稳定
	assert.Equal(t, long, flattenNamespaceToolName("very_long_namespace_prefix_for_testing_purposes", "and_a_rather_long_tool_name_too"))
placeholder

func TestResponsesInputToChatMessages_ToolSearchCallHistory(t *testing.T) {
	input := json.RawMessage(`[
		{"role":"user","content":"find tools"placeholder,
		{"type":"tool_search_call","call_id":"call_s","arguments":{"query":"gmail"placeholderplaceholder,
		{"type":"tool_search_output","call_id":"call_s","output":{"groups":["gmail"]placeholderplaceholder
	]`)

	messages, err := responsesInputToChatMessages("", input)
placeholder
	require.Len(t, messages, 3)

	require.Len(t, messages[1].ToolCalls, 1)
	assert.Equal(t, "tool_search", messages[1].ToolCalls[0].Function.Name)
	assert.JSONEq(t, `{"query":"gmail"placeholder`, messages[1].ToolCalls[0].Function.Arguments)

	assert.Equal(t, "tool", messages[2].Role)
	assert.Equal(t, "call_s", messages[2].ToolCallID)
	assert.JSONEq(t, `"{\"groups\":[\"gmail\"]placeholder"`, string(messages[2].Content))
placeholder

func TestResponsesInputToChatMessages_NamespacedFunctionCallHistory(t *testing.T) {
	input := json.RawMessage(`[
		{"type":"function_call","call_id":"call_n","name":"send","namespace":"gmail","arguments":"{\"to\":\"a\"placeholder"placeholder,
		{"type":"function_call_output","call_id":"call_n","output":"ok"placeholder
	]`)

	messages, err := responsesInputToChatMessages("", input)
placeholder
	require.Len(t, messages, 2)

	require.Len(t, messages[0].ToolCalls, 1)
	assert.Equal(t, "gmail__send", messages[0].ToolCalls[0].Function.Name)
placeholder

func TestChatCompletionsChunkToResponsesEvents_CustomToolNameArrivesLate(t *testing.T) {
	state := NewChatCompletionsToResponsesStreamState("glm-5.2")
	state.CustomTools = map[string]bool{"exec": trueplaceholder

	idx := 0
	chunk1 := &ChatCompletionsChunk{Choices: []ChatChunkChoice{{Delta: ChatDelta{
		ToolCalls: []ChatToolCall{{Index: &idx, ID: "call_1", Function: ChatFunctionCall{Arguments: `{"inp`placeholderplaceholderplaceholder,
placeholderplaceholderplaceholderplaceholder
	chunk2 := &ChatCompletionsChunk{Choices: []ChatChunkChoice{{Delta: ChatDelta{
		ToolCalls: []ChatToolCall{{Index: &idx, Function: ChatFunctionCall{Name: "exec", Arguments: `ut": "dir"placeholder`placeholderplaceholderplaceholder,
placeholderplaceholderplaceholderplaceholder

	var events []ResponsesStreamEvent
	events = append(events, ChatCompletionsChunkToResponsesEvents(chunk1, state)...)
	events = append(events, ChatCompletionsChunkToResponsesEvents(chunk2, state)...)
	events = append(events, FinalizeChatCompletionsResponsesStream(state)...)

	addedCount := 0
	for _, evt := range events {
		switch evt.Type {
		case "response.output_item.added":
			if evt.Item != nil && evt.Item.Type != "reasoning" && evt.Item.Type != "message" {
				addedCount++
				assert.Equal(t, "custom_tool_call", evt.Item.Type, "迟到的名字命中 custom 工具时按 custom_tool_call 宣告")
				assert.Equal(t, "exec", evt.Item.Name)
		placeholder
		case "response.function_call_arguments.delta", "response.function_call_arguments.done":
			t.Fatalf("custom 调用不应产出 function 参数事件: %s", evt.Type)
		case "response.custom_tool_call_input.done":
			assert.Equal(t, "dir", evt.Input)
	placeholder
placeholder
	assert.Equal(t, 1, addedCount, "工具调用只宣告一次")
placeholder

func TestChatCompletionsChunkToResponsesEvents_FunctionToolNameArrivesLate(t *testing.T) {
	state := NewChatCompletionsToResponsesStreamState("glm-5.2")
	state.CustomTools = map[string]bool{"exec": trueplaceholder

	idx := 0
	chunk1 := &ChatCompletionsChunk{Choices: []ChatChunkChoice{{Delta: ChatDelta{
		ToolCalls: []ChatToolCall{{Index: &idx, ID: "call_9", Function: ChatFunctionCall{Arguments: `{"cell`placeholderplaceholderplaceholder,
placeholderplaceholderplaceholderplaceholder
	chunk2 := &ChatCompletionsChunk{Choices: []ChatChunkChoice{{Delta: ChatDelta{
		ToolCalls: []ChatToolCall{{Index: &idx, Function: ChatFunctionCall{Name: "wait", Arguments: `_id": 3placeholder`placeholderplaceholderplaceholder,
placeholderplaceholderplaceholderplaceholder

	var events []ResponsesStreamEvent
	events = append(events, ChatCompletionsChunkToResponsesEvents(chunk1, state)...)
	events = append(events, ChatCompletionsChunkToResponsesEvents(chunk2, state)...)
	events = append(events, FinalizeChatCompletionsResponsesStream(state)...)

	deltas := ""
	argsDone := ""
	for _, evt := range events {
		switch evt.Type {
		case "response.function_call_arguments.delta":
			deltas += evt.Delta
		case "response.function_call_arguments.done":
			argsDone = evt.Arguments
		case "response.custom_tool_call_input.done":
			t.Fatal("function 调用不应产出 custom 事件")
	placeholder
placeholder
	assert.Equal(t, `{"cell_id": 3placeholder`, deltas, "宣告前累积的参数需在宣告时补发")
	assert.Equal(t, `{"cell_id": 3placeholder`, argsDone)
placeholder

// 序列化层（MarshalJSON → responsesItemWire）单独走白名单重组，事件结构体上的字段
// 齐全不代表落到 SSE 线上的 JSON 齐全，必须在 wire 层再断言一次。
func TestResponsesEventToSSE_CustomToolCallItemCarriesAllFields(t *testing.T) {
	evt := ResponsesStreamEvent{
		Type:        "response.output_item.done",
		OutputIndex: 1,
		Item: &ResponsesOutput{
			Type:   "custom_tool_call",
			ID:     "item_1",
			CallID: "call_1",
			Name:   "exec",
			Input:  "dir",
			Status: "completed",
	placeholder,
placeholder

	sse, err := ResponsesEventToSSE(evt)
placeholder

	assert.Contains(t, sse, `"call_id":"call_1"`)
	assert.Contains(t, sse, `"name":"exec"`)
	assert.Contains(t, sse, `"input":"dir"`)
	assert.Contains(t, sse, `"type":"custom_tool_call"`)
placeholder

func TestChatCompletionsChunkToResponsesEvents_FunctionToolStreamUnaffected(t *testing.T) {
	state := NewChatCompletionsToResponsesStreamState("glm-5.2")
	state.CustomTools = map[string]bool{"exec": trueplaceholder

	idx := 0
	chunk := &ChatCompletionsChunk{
		Choices: []ChatChunkChoice{{
			Delta: ChatDelta{
				ToolCalls: []ChatToolCall{{
					Index:    &idx,
					ID:       "call_9",
					Function: ChatFunctionCall{Name: "wait", Arguments: `{"cell_id": 3placeholder`placeholder,
		placeholder
		placeholder,
placeholder
placeholder

	events := ChatCompletionsChunkToResponsesEvents(chunk, state)
	events = append(events, FinalizeChatCompletionsResponsesStream(state)...)

	sawArgsDelta := false
	for _, evt := range events {
		if evt.Type == "response.function_call_arguments.delta" {
			sawArgsDelta = true
	placeholder
		if evt.Type == "response.custom_tool_call_input.done" {
			t.Fatal("function 工具不应产出 custom_tool_call 事件")
	placeholder
placeholder
	assert.True(t, sawArgsDelta, "function 工具应保持原有参数增量事件")
placeholder
