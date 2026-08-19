package service

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestSanitizeOpenAIResponsesOrphanToolOutputs(t *testing.T) {
	t.Run("preserves matches regardless of item order", func(t *testing.T) {
		input := []any{
			map[string]any{"type": "tool_search_output", "call_id": "search_1", "output": "first"placeholder,
			map[string]any{"type": "tool_search_call", "id": "search_1", "query": "docs"placeholder,
			map[string]any{"type": "custom_tool_call_output", "call_id": "custom_1", "output": "second"placeholder,
			map[string]any{"type": "item_reference", "id": "custom_1"placeholder,
	placeholder
		reqBody := map[string]any{"input": inputplaceholder

		require.False(t, sanitizeOpenAIResponsesOrphanToolOutputs(reqBody, input, false))
		require.Equal(t, input, reqBody["input"])
placeholder)

	t.Run("outputs do not legitimize each other", func(t *testing.T) {
		input := []any{
			map[string]any{"type": "function_call_output", "call_id": "missing", "output": "one"placeholder,
			map[string]any{"type": "tool_search_output", "call_id": "missing", "output": "two"placeholder,
			map[string]any{"type": "custom_tool_call_output", "call_id": "missing", "output": "three"placeholder,
			map[string]any{"type": "mcp_tool_call_output", "call_id": "missing", "output": "four"placeholder,
	placeholder
		reqBody := map[string]any{"input": inputplaceholder

		require.True(t, sanitizeOpenAIResponsesOrphanToolOutputs(reqBody, input, false))
		got := reqBody["input"].([]any)
		require.Empty(t, got)
placeholder)

	t.Run("preserves all output variants with matching calls", func(t *testing.T) {
		pairs := []struct {
			callType   string
			outputType string
	placeholder{
			{callType: "function_call", outputType: "function_call_output"placeholder,
			{callType: "tool_search_call", outputType: "tool_search_output"placeholder,
			{callType: "custom_tool_call", outputType: "custom_tool_call_output"placeholder,
			{callType: "mcp_tool_call", outputType: "mcp_tool_call_output"placeholder,
	placeholder
		input := make([]any, 0, len(pairs)*2)
		for index, pair := range pairs {
			callID := string(rune('a' + index))
			input = append(input,
				map[string]any{"type": pair.callType, "call_id": callIDplaceholder,
				map[string]any{"type": pair.outputType, "call_id": callID, "output": "ok"placeholder,
			)
	placeholder
		reqBody := map[string]any{"input": inputplaceholder

		require.False(t, sanitizeOpenAIResponsesOrphanToolOutputs(reqBody, input, false))
placeholder)

	t.Run("previous response may contain the missing call", func(t *testing.T) {
		input := []any{map[string]any{"type": "function_call_output", "call_id": "remote", "output": "ok"placeholderplaceholder
		reqBody := map[string]any{"input": input, "previous_response_id": "resp_1"placeholder

		require.False(t, sanitizeOpenAIResponsesOrphanToolOutputs(reqBody, input, true))
		require.Equal(t, input, reqBody["input"])
placeholder)
placeholder

func TestTruncateOpenAIResponsesInputText(t *testing.T) {
	oversized := strings.Repeat("a", openAIResponsesInputTextMaxChars+1)
	input := []any{
		map[string]any{"type": "function_call_output", "call_id": "a", "output": oversizedplaceholder,
		map[string]any{"type": "tool_search_output", "call_id": "b", "output": oversizedplaceholder,
		map[string]any{"type": "custom_tool_call_output", "call_id": "c", "output": oversizedplaceholder,
		map[string]any{"type": "mcp_tool_call_output", "call_id": "d", "output": oversizedplaceholder,
		map[string]any{
			"role": "user",
			"content": []any{
				map[string]any{"type": "input_text", "text": "short"placeholder,
				map[string]any{"type": "input_text", "text": oversizedplaceholder,
		placeholder,
	placeholder,
placeholder
	reqBody := map[string]any{"input": inputplaceholder

	require.True(t, truncateOpenAIResponsesInputText(reqBody))
	for _, rawItem := range input[:4] {
		item := rawItem.(map[string]any)
		require.Len(t, item["output"].(string), openAIResponsesInputTextMaxChars)
placeholder
	content := input[4].(map[string]any)["content"].([]any)
	require.Equal(t, "short", content[0].(map[string]any)["text"])
	require.Len(t, content[1].(map[string]any)["text"].(string), openAIResponsesInputTextMaxChars)
placeholder

func TestTruncateOpenAIResponsesInputStringPreservesUTF8Boundary(t *testing.T) {
	value := strings.Repeat("a", openAIResponsesInputTextMaxChars-1) + "中中"

	got, changed := truncateOpenAIResponsesInputString(value)

	require.True(t, changed)
	require.True(t, utf8.ValidString(got))
	require.Equal(t, openAIResponsesInputTextMaxChars, utf8.RuneCountInString(got))
	require.Equal(t, "中", got[len(got)-len("中"):])
placeholder

func TestOpenAIResponsesInputMayNeedTruncation(t *testing.T) {
	short := []byte(`{"input":[{"type":"function_call_output","output":"ok"placeholder]placeholder`)
	largeUnrelated := []byte(`{"input":"` + strings.Repeat("x", openAIResponsesInputTextMaxChars+1) + `"placeholder`)
	largeOutput := []byte(`{"input":[{"type":"function_call_output","output":"` + strings.Repeat("x", openAIResponsesInputTextMaxChars+1) + `"placeholder]placeholder`)

	require.False(t, openAIResponsesInputMayNeedTruncation(short))
	require.False(t, openAIResponsesInputMayNeedTruncation(largeUnrelated))
	require.True(t, openAIResponsesInputMayNeedTruncation(largeOutput))
placeholder

func TestOpenAIGatewayService_OAuthDropsOrphanAfterDroppingPreviousResponse(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5","stream":false,"previous_response_id":"resp_missing","input":[{"type":"function_call_output","call_id":"call_missing","output":"keep this result"placeholder]placeholder`)
	upstream := &httpUpstreamRecorder{responses: []*http.Response{
		newOpenAIRejectedFieldTestResponse(http.StatusOK, `{"id":"resp_ok","output":[],"usage":{"input_tokens":1,"output_tokens":1,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholder`),
placeholderplaceholder

	result, err := newOpenAIRejectedFieldTestService(upstream).Forward(
		context.Background(),
		newOpenAIRejectedFieldTestContext(body),
		newOpenAIOAuthNamespaceTestAccount(),
		body,
	)

placeholder
	require.NotNil(t, result)
	require.Len(t, upstream.bodies, 1)
	require.False(t, gjson.GetBytes(upstream.bodies[0], "previous_response_id").Exists())
	require.Empty(t, gjson.GetBytes(upstream.bodies[0], "input").Array())
placeholder

func TestOpenAIGatewayService_TruncatesOversizedToolOutputBeforeForward(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5","stream":false,"input":[{"type":"function_call","call_id":"call_1","name":"lookup","arguments":"{placeholder"placeholder,{"type":"function_call_output","call_id":"call_1","output":"` + strings.Repeat("x", openAIResponsesInputTextMaxChars+1) + `"placeholder]placeholder`)
	upstream := &httpUpstreamRecorder{responses: []*http.Response{
		newOpenAIRejectedFieldTestResponse(http.StatusOK, `{"id":"resp_ok","output":[],"usage":{"input_tokens":1,"output_tokens":1,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholder`),
placeholderplaceholder

	result, err := newOpenAIRejectedFieldTestService(upstream).Forward(
		context.Background(),
		newOpenAIRejectedFieldTestContext(body),
		newOpenAIRejectedFieldTestAccount(),
		body,
	)

placeholder
	require.NotNil(t, result)
	require.Len(t, upstream.bodies, 1)
	require.Equal(t, openAIResponsesInputTextMaxChars, len(gjson.GetBytes(upstream.bodies[0], "input.1.output").String()))
placeholder
