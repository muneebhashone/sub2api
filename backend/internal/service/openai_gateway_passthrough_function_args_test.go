package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai_compat"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestHandleStreamingResponsePassthroughDeduplicatesFunctionCallArguments(t *testing.T) {
	gin.SetMode(gin.TestMode)

	argsA := `{"cmd":"echo hi","meta":{"nested":[1,{"ok":trueplaceholder],"quote":"aplaceholderb"placeholderplaceholder`
	argsB := `{"path":"/tmp/file","patch":{"ops":[{"op":"replace","value":{"lines":["x","y"]placeholderplaceholder]placeholderplaceholder`
	upstreamBody := strings.Join([]string{
		passthroughSSEData(`{"type":"response.created","response":{"id":"resp_passthrough_args","model":"gpt-5.4"placeholderplaceholder`),
		passthroughSSEData(`{"type":"response.output_item.added","output_index":0,"item":{"type":"function_call","id":"fc_a","call_id":"call_a","name":"exec_command","arguments":"","status":"in_progress"placeholderplaceholder`),
		passthroughSSEData(functionArgsDeltaJSON(0, "fc_a", "call_a", "exec_command", `{"cmd":`)),
		passthroughSSEData(functionArgsDeltaJSON(0, "fc_a", "call_a", "exec_command", `"echo hi","meta":{"nested":[1,{"ok":trueplaceholder],"quote":"aplaceholderb"placeholderplaceholder`)),
		passthroughSSEData(functionArgsDoneJSON(0, "fc_a", "call_a", "exec_command", argsA+argsA)),
		passthroughSSEData(outputItemDoneJSON(0, "fc_a", "call_a", "exec_command", argsA+argsA)),
		passthroughSSEData(`{"type":"response.output_item.added","output_index":1,"item":{"type":"function_call","id":"fc_b","call_id":"call_b","name":"apply_patch","arguments":"","status":"in_progress"placeholderplaceholder`),
		passthroughSSEData(functionArgsDeltaJSON(1, "fc_b", "call_b", "apply_patch", `{"path":"/tmp/file",`)),
		passthroughSSEData(functionArgsDeltaJSON(1, "fc_b", "call_b", "apply_patch", `"patch":{"ops":[{"op":"replace","value":{"lines":["x","y"]placeholderplaceholder]placeholderplaceholder`)),
		passthroughSSEData(functionArgsDoneJSON(1, "fc_b", "call_b", "apply_patch", argsB+argsB)),
		passthroughSSEData(outputItemDoneJSON(1, "fc_b", "call_b", "apply_patch", argsB+argsB)),
		passthroughSSEData(completedWithFunctionCallsJSON(argsA+argsA, argsB+argsB)),
		"data: [DONE]\n\n",
placeholder, "")

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamBody)),
placeholder

	svc := &OpenAIGatewayService{placeholder
	result, err := svc.handleStreamingResponsePassthrough(context.Background(), resp, c, &Account{ID: 1placeholder, time.Now(), "gpt-5.4", "gpt-5.4")
placeholder
	require.NotNil(t, result)

	events := collectSSEDataPayloads(t, rec.Body.String())
	require.Equal(t, argsA, accumulateFunctionArgumentDeltas(events, "call_a"))
	require.Equal(t, argsB, accumulateFunctionArgumentDeltas(events, "call_b"))

	require.Equal(t, argsA, gjson.Get(findSSEEvent(t, events, "response.function_call_arguments.done", "call_a"), "arguments").String())
	require.Equal(t, argsB, gjson.Get(findSSEEvent(t, events, "response.function_call_arguments.done", "call_b"), "arguments").String())
	require.Equal(t, argsA, gjson.Get(findSSEEvent(t, events, "response.output_item.done", "call_a"), "item.arguments").String())
	require.Equal(t, argsB, gjson.Get(findSSEEvent(t, events, "response.output_item.done", "call_b"), "item.arguments").String())

	completed := findSSEEvent(t, events, "response.completed", "")
	require.Equal(t, argsA, gjson.Get(completed, "response.output.0.arguments").String())
	require.Equal(t, argsB, gjson.Get(completed, "response.output.1.arguments").String())
	requireJSONArgument(t, gjson.Get(completed, "response.output.0.arguments").String())
	requireJSONArgument(t, gjson.Get(completed, "response.output.1.arguments").String())
placeholder

func TestForwardResponsesChatCompletionsFallbackKeepsFunctionArgumentsSingle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"model":"gpt-5.4","input":"run a command","stream":trueplaceholder`)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", strings.NewReader(string(body)))
	c.Request.Header.Set("Content-Type", "application/json")

	upstreamBody := strings.Join([]string{
		passthroughSSEData(chatToolCallChunkJSON(true, "")),
		"",
		passthroughSSEData(chatToolCallChunkJSON(false, `{"cmd":"echo hi"placeholder`)),
		"",
		`data: {"id":"chatcmpl_tool","object":"chat.completion.chunk","model":"gpt-5.4","choices":[{"index":0,"delta":{placeholder,"finish_reason":"tool_calls"placeholder],"usage":{"prompt_tokens":3,"completion_tokens":2,"total_tokens":5placeholderplaceholder`,
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid_fallback_tool_args"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamBody)),
placeholderplaceholder
	account := passthroughArgsFallbackAccount()
	account.Extra = map[string]any{
		openai_compat.ExtraKeyResponsesMode: string(openai_compat.ResponsesSupportModeForceChatCompletions),
placeholder
	svc := &OpenAIGatewayService{
		cfg:          passthroughArgsTestConfig(),
		httpUpstream: upstream,
placeholder

	result, err := svc.Forward(context.Background(), c, account, body)
placeholder
	require.NotNil(t, result)

	const wantArgs = `{"cmd":"echo hi"placeholder`
	events := collectSSEDataPayloads(t, rec.Body.String())
	require.Equal(t, wantArgs, accumulateFunctionArgumentDeltas(events, "chatcmpl-tool-a"))
	require.Equal(t, wantArgs, gjson.Get(findSSEEvent(t, events, "response.function_call_arguments.done", "chatcmpl-tool-a"), "arguments").String())
	require.Equal(t, wantArgs, gjson.Get(findSSEEvent(t, events, "response.output_item.done", "chatcmpl-tool-a"), "item.arguments").String())
placeholder

func passthroughSSEData(payload string) string {
	return "data: " + payload + "\n\n"
placeholder

func functionArgsDeltaJSON(outputIndex int, itemID, callID, name, delta string) string {
	return fmt.Sprintf(
		`{"type":"response.function_call_arguments.delta","output_index":%d,"item_id":%s,"call_id":%s,"name":%s,"delta":%splaceholder`,
		outputIndex,
		strconv.Quote(itemID),
		strconv.Quote(callID),
		strconv.Quote(name),
		strconv.Quote(delta),
	)
placeholder

func functionArgsDoneJSON(outputIndex int, itemID, callID, name, arguments string) string {
	return fmt.Sprintf(
		`{"type":"response.function_call_arguments.done","output_index":%d,"item_id":%s,"call_id":%s,"name":%s,"arguments":%splaceholder`,
		outputIndex,
		strconv.Quote(itemID),
		strconv.Quote(callID),
		strconv.Quote(name),
		strconv.Quote(arguments),
	)
placeholder

func outputItemDoneJSON(outputIndex int, itemID, callID, name, arguments string) string {
	return fmt.Sprintf(
		`{"type":"response.output_item.done","output_index":%d,"item":{"type":"function_call","id":%s,"call_id":%s,"name":%s,"arguments":%s,"status":"completed"placeholderplaceholder`,
		outputIndex,
		strconv.Quote(itemID),
		strconv.Quote(callID),
		strconv.Quote(name),
		strconv.Quote(arguments),
	)
placeholder

func completedWithFunctionCallsJSON(argsA, argsB string) string {
	return fmt.Sprintf(
		`{"type":"response.completed","response":{"id":"resp_passthrough_args","status":"completed","output":[{"type":"function_call","id":"fc_a","call_id":"call_a","name":"exec_command","arguments":%s,"status":"completed"placeholder,{"type":"function_call","id":"fc_b","call_id":"call_b","name":"apply_patch","arguments":%s,"status":"completed"placeholder],"usage":{"input_tokens":2,"output_tokens":3,"total_tokens":5placeholderplaceholderplaceholder`,
		strconv.Quote(argsA),
		strconv.Quote(argsB),
	)
placeholder

func chatToolCallChunkJSON(includeIdentity bool, arguments string) string {
	identity := ""
	functionFields := make([]string, 0, 2)
	if includeIdentity {
		identity = `"id":"chatcmpl-tool-a","type":"function",`
		functionFields = append(functionFields, `"name":"exec_command"`)
placeholder
	if includeIdentity || arguments != "" {
		functionFields = append(functionFields, `"arguments":`+strconv.Quote(arguments))
placeholder
	return fmt.Sprintf(
		`{"id":"chatcmpl_tool","object":"chat.completion.chunk","model":"gpt-5.4","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,%s"function":{%splaceholderplaceholder]placeholder,"finish_reason":nullplaceholder]placeholder`,
		identity,
		strings.Join(functionFields, ","),
	)
placeholder

func passthroughArgsTestConfig() *config.Config {
	return &config.Config{
		Security: config.SecurityConfig{
			URLAllowlist: config.URLAllowlistConfig{
				Enabled:           false,
				AllowInsecureHTTP: true,
		placeholder,
	placeholder,
placeholder
placeholder

func passthroughArgsFallbackAccount() *Account {
placeholder
		ID:          102,
		Name:        "passthrough-args-openai-apikey",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "http://upstream.example",
	placeholder,
placeholder
placeholder

func collectSSEDataPayloads(t *testing.T, body string) []string {
placeholder
	scanner := bufio.NewScanner(strings.NewReader(body))
	var events []string
	for scanner.Scan() {
		data, ok := extractOpenAISSEDataLine(scanner.Text())
		if !ok {
			continue
	placeholder
		if strings.TrimSpace(data) == "[DONE]" {
			continue
	placeholder
		require.True(t, gjson.Valid(data), "invalid SSE data payload: %s", data)
		events = append(events, data)
placeholder
	require.NoError(t, scanner.Err())
	return events
placeholder

func findSSEEvent(t *testing.T, events []string, eventType, callID string) string {
placeholder
	for _, event := range events {
		if gjson.Get(event, "type").String() != eventType {
			continue
	placeholder
		if callID == "" ||
			gjson.Get(event, "call_id").String() == callID ||
			gjson.Get(event, "item.call_id").String() == callID {
			return event
	placeholder
placeholder
	t.Fatalf("missing event type=%s call_id=%s in %d events", eventType, callID, len(events))
	return ""
placeholder

func accumulateFunctionArgumentDeltas(events []string, callID string) string {
	var b strings.Builder
	for _, event := range events {
		if gjson.Get(event, "type").String() != "response.function_call_arguments.delta" {
			continue
	placeholder
		if gjson.Get(event, "call_id").String() != callID {
			continue
	placeholder
		_, _ = b.WriteString(gjson.Get(event, "delta").String())
placeholder
	return b.String()
placeholder

func requireJSONArgument(t *testing.T, arguments string) {
placeholder
	var decoded any
	require.NoError(t, json.Unmarshal([]byte(arguments), &decoded))
placeholder
