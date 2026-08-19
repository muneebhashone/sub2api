package apicompat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponsesInputToChatMessages_DeveloperRoleMapsToSystem(t *testing.T) {
	messages, err := responsesInputToChatMessages("", json.RawMessage(`[{"role":"developer","content":"follow project instructions"placeholder]`))
placeholder
	require.Len(t, messages, 1)

	assert.Equal(t, "system", messages[0].Role)
	assert.JSONEq(t, `"follow project instructions"`, string(messages[0].Content))
placeholder

func TestResponsesInputToChatMessages_SkipsInvalidHistoricalFunctionCall(t *testing.T) {
	input := json.RawMessage(`[
		{"type":"function_call","call_id":"call_bad","name":"exec_command","arguments":"{\"cmd\": \"ssh root@HOST"placeholder,
		{"type":"function_call_output","call_id":"call_bad","output":"failed to parse function arguments"placeholder,
		{"type":"function_call","call_id":"call_ok","name":"exec_command","arguments":"{placeholder"placeholder,
		{"type":"function_call_output","call_id":"call_ok","output":"ok"placeholder,
		{"role":"user","content":"continue"placeholder
	]`)

	messages, err := responsesInputToChatMessages("", input)
placeholder
	require.Len(t, messages, 3)
	require.Equal(t, "assistant", messages[0].Role)
	require.Len(t, messages[0].ToolCalls, 1)
	require.Equal(t, "call_ok", messages[0].ToolCalls[0].ID)
	require.Equal(t, "tool", messages[1].Role)
	require.Equal(t, "call_ok", messages[1].ToolCallID)
	require.Equal(t, "user", messages[2].Role)
placeholder

func TestResponsesInputToChatMessages_KeepsChatCompletionRoles(t *testing.T) {
	input := json.RawMessage(`[
		{"role":"system","content":"system message"placeholder,
		{"role":"user","content":"user message"placeholder,
		{"role":"assistant","content":"assistant message"placeholder,
		{"role":"tool","content":"tool message"placeholder
	]`)

	messages, err := responsesInputToChatMessages("", input)
placeholder
	require.Len(t, messages, 4)

	assert.Equal(t, []string{"system", "user", "assistant", "tool"placeholder, chatMessageRoles(messages))
placeholder

func TestResponsesInputToChatMessages_EmptyRoleFallsBackToUser(t *testing.T) {
	messages, err := responsesInputToChatMessages("", json.RawMessage(`[{"role":"","content":"hello"placeholder]`))
placeholder
	require.Len(t, messages, 1)

	assert.Equal(t, "user", messages[0].Role)
placeholder

func TestResponsesInputToChatMessages_DeveloperRoleTrimAndCaseInsensitive(t *testing.T) {
	input := json.RawMessage(`[
		{"role":" Developer ","content":"one"placeholder,
		{"role":"\tDEVELOPER\n","content":"two"placeholder
	]`)

	messages, err := responsesInputToChatMessages("", input)
placeholder
	require.Len(t, messages, 2)

	assert.Equal(t, []string{"system", "system"placeholder, chatMessageRoles(messages))
placeholder

func TestResponsesToChatCompletionsRequest_InstructionsAndInputDeveloperRole(t *testing.T) {
	req := &ResponsesRequest{
		Model:        "gpt-4o",
		Instructions: "Use concise answers.",
		Input: json.RawMessage(`[
			{"role":"developer","content":[{"type":"input_text","text":"Prefer JSON."placeholder]placeholder,
			{"role":"user","content":"Hello"placeholder
		]`),
placeholder

	out, err := ResponsesToChatCompletionsRequest(req)
placeholder
	require.Len(t, out.Messages, 3)

	assert.Equal(t, []string{"system", "system", "user"placeholder, chatMessageRoles(out.Messages))
	assert.JSONEq(t, `"Use concise answers."`, string(out.Messages[0].Content))
	assert.JSONEq(t, `"Prefer JSON."`, string(out.Messages[1].Content))
	assert.JSONEq(t, `"Hello"`, string(out.Messages[2].Content))
placeholder

func TestResponsesToChatCompletionsRequest_TextFormatJsonObject(t *testing.T) {
	req := &ResponsesRequest{
		Model: "gpt-4o",
		Input: json.RawMessage(`[
			{"role":"user","content":"Return JSON"placeholder
		]`),
		Text: &ResponsesText{
			Format: json.RawMessage(`{"type":"json_object"placeholder`),
	placeholder,
placeholder

	out, err := ResponsesToChatCompletionsRequest(req)
placeholder
	assert.JSONEq(t, `{"type":"json_object"placeholder`, string(out.ResponseFormat))
placeholder

func TestResponsesToChatCompletionsRequest_TextFormatJsonSchema(t *testing.T) {
	req := &ResponsesRequest{
		Model: "gpt-4o",
		Input: json.RawMessage(`[
			{"role":"user","content":"Return structured JSON"placeholder
		]`),
		Text: &ResponsesText{
			Format: json.RawMessage(`{
				"type":"json_schema",
				"name":"answer",
				"schema":{
					"type":"object",
					"properties":{"ok":{"type":"boolean"placeholderplaceholder,
					"required":["ok"],
					"additionalProperties":false
			placeholder,
				"strict":true
		placeholder`),
	placeholder,
placeholder

	out, err := ResponsesToChatCompletionsRequest(req)
placeholder
	assert.JSONEq(t, `{
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
placeholder`, string(out.ResponseFormat))
placeholder

func TestResponsesToChatCompletionsRequest_ParallelToolCalls(t *testing.T) {
	parallel := false
	req := &ResponsesRequest{
		Model: "gpt-4o",
		Input: json.RawMessage(`[
			{"role":"user","content":"Use tools"placeholder
		]`),
		ParallelToolCalls: &parallel,
placeholder

	out, err := ResponsesToChatCompletionsRequest(req)
placeholder
	require.NotNil(t, out.ParallelToolCalls)
	assert.False(t, *out.ParallelToolCalls)

	payload, err := json.Marshal(out)
placeholder
	assert.Contains(t, string(payload), `"parallel_tool_calls":false`)
placeholder

func chatMessageRoles(messages []ChatMessage) []string {
	roles := make([]string, 0, len(messages))
	for _, message := range messages {
		roles = append(roles, message.Role)
placeholder
	return roles
placeholder
