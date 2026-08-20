package service

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/apicompat"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func openAIClientToolsRequest(stream bool) []byte {
	streamValue := "false"
	if stream {
		streamValue = "true"
placeholder
	return []byte(`{"model":"gpt-5.4","input":"fix it","stream":` + streamValue + `,"tools":[{"type":"custom","name":"exec"placeholder,{"type":"custom","name":"apply_patch"placeholder]placeholder`)
placeholder

func assertOpenAIClientToolsLowered(t *testing.T, body []byte) {
placeholder
	for index, name := range []string{"exec", "apply_patch"placeholder {
		tool := gjson.GetBytes(body, "tools."+string(rune('0'+index)))
		require.Equal(t, "function", tool.Get("type").String())
		require.Equal(t, name, tool.Get("name").String())
		require.Equal(t, "string", tool.Get("parameters.properties.input.type").String())
placeholder
placeholder

func openAIClientToolsTestService(upstream *httpUpstreamRecorder) *OpenAIGatewayService {
	return &OpenAIGatewayService{
		httpUpstream: upstream,
		cfg: &config.Config{Security: config.SecurityConfig{
			URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholder,
placeholder
placeholder
placeholder

func TestAdaptOpenAIResponsesClientToolsLeavesNamespaceOnlyBodyUnchanged(t *testing.T) {
	body := []byte(`{
		"model": "gpt-5.5",
		"tools": [{"type": "namespace", "name": "code_tools", "tools": [{"type": "function", "name": "run"placeholder]placeholder],
		"tool_choice": "auto"
placeholder`)

	adapted, mapping, err := adaptOpenAIResponsesClientTools(body)

placeholder
	require.Equal(t, body, adapted)
	require.Empty(t, mapping.CustomTools)
	require.Empty(t, mapping.NamespaceTools)
	require.False(t, mapping.ToolSearch)
placeholder

func TestAdaptOpenAIResponsesClientToolsRejectsTrailingData(t *testing.T) {
	tests := map[string][]byte{
		"trailing garbage":     append(openAIClientToolsRequest(false), []byte(` garbage`)...),
		"second JSON document": append(openAIClientToolsRequest(false), []byte(` {"model":"other"placeholder`)...),
placeholder

	for name, body := range tests {
		t.Run(name, func(t *testing.T) {
			adapted, mapping, err := adaptOpenAIResponsesClientTools(body)

			require.ErrorContains(t, err, "decode OpenAI Responses client tools trailing data")
			require.Equal(t, body, adapted)
			require.Empty(t, mapping)
	placeholder)
placeholder
placeholder

func TestResponsesFunctionUpstreamsLowerToolSearchDiscoveryOutput(t *testing.T) {
	body := []byte(`{"tools":[{"type":"tool_search"placeholder],"input":[{"type":"tool_search_output","call_id":"search_1","tools":[{"type":"namespace","name":"github"placeholder],"status":"completed","execution":"client"placeholder]placeholder`)
	adapters := map[string]func([]byte) ([]byte, apicompat.ResponsesClientToolMapping, error){
		"OpenAI API-key": adaptOpenAIResponsesClientTools,
		"Grok":           adaptGrokResponsesClientTools,
placeholder
	for name, adapt := range adapters {
		t.Run(name, func(t *testing.T) {
			adapted, mapping, err := adapt(body)
		placeholder
			require.True(t, mapping.ToolSearch)
			require.Equal(t, "function_call_output", gjson.GetBytes(adapted, "input.0.type").String())
			require.JSONEq(t, `[{"name":"github","type":"namespace"placeholder]`, gjson.GetBytes(adapted, "input.0.output").String())
			require.False(t, gjson.GetBytes(adapted, "input.0.tools").Exists())
			require.False(t, gjson.GetBytes(adapted, "input.0.status").Exists())
			require.False(t, gjson.GetBytes(adapted, "input.0.execution").Exists())
	placeholder)
placeholder
placeholder

func TestClearOpenAIResponsesClientToolMappingRemovesStaleContextState(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(openAIResponsesClientToolMappingContextKey, apicompat.ResponsesClientToolMapping{CustomTools: map[string]bool{"exec": trueplaceholderplaceholder)

	clearOpenAIResponsesClientToolMapping(c)

	_, ok := openAIResponsesClientToolMapping(c)
	require.False(t, ok)
placeholder

func TestOpenAIPassthroughAPIKeyRestoresClientToolsNonStreaming(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := openAIClientToolsRequest(false)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(`{"id":"resp_tools","status":"completed","output":[
			{"type":"function_call","id":"i1","call_id":"c1","name":"exec","arguments":"{\"input\":\"pwd\"placeholder"placeholder,
			{"type":"function_call","id":"i2","call_id":"c2","name":"apply_patch","arguments":"{\"input\":\"*** Begin Patch\"placeholder"placeholder],"usage":{placeholderplaceholder`)),
placeholderplaceholder
	svc := openAIClientToolsTestService(upstream)
	account := &Account{ID: 5659, Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Credentials: map[string]any{"api_key": "test-key"placeholderplaceholder

	result, err := svc.forwardOpenAIPassthrough(context.Background(), c, account, body, body, "gpt-5.4", false, nil, false, time.Now())

placeholder
	require.NotNil(t, result)
	assertOpenAIClientToolsLowered(t, upstream.lastBody)
	require.Equal(t, "custom_tool_call", gjson.Get(recorder.Body.String(), "output.0.type").String())
	require.Equal(t, "pwd", gjson.Get(recorder.Body.String(), "output.0.input").String())
	require.Equal(t, "custom_tool_call", gjson.Get(recorder.Body.String(), "output.1.type").String())
	require.Equal(t, "*** Begin Patch", gjson.Get(recorder.Body.String(), "output.1.input").String())
placeholder

func TestOpenAIPassthroughAPIKeyRestoresClientToolsStreaming(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := openAIClientToolsRequest(true)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))

	sse := strings.Join([]string{
		`data: {"type":"response.output_item.added","sequence_number":0,"output_index":0,"item":{"type":"function_call","id":"i1","call_id":"c1","name":"apply_patch","status":"in_progress"placeholderplaceholder`,
		`data: {"type":"response.function_call_arguments.done","sequence_number":1,"item_id":"i1","call_id":"c1","name":"apply_patch","arguments":"{\"input\":\"*** Begin Patch\"placeholder"placeholder`,
		`data: {"type":"response.output_item.done","sequence_number":2,"output_index":0,"item":{"type":"function_call","id":"i1","call_id":"c1","name":"apply_patch","arguments":"{\"input\":\"*** Begin Patch\"placeholder","status":"completed"placeholderplaceholder`,
		`data: {"type":"response.completed","sequence_number":3,"response":{"id":"resp_stream_tools","status":"completed","output":[{"type":"function_call","id":"i1","call_id":"c1","name":"apply_patch","arguments":"{\"input\":\"*** Begin Patch\"placeholder"placeholder],"usage":{"input_tokens":1,"output_tokens":1placeholderplaceholderplaceholder`,
placeholder, "\n\n") + "\n\n"
	upstream := &httpUpstreamRecorder{resp: &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder, Body: io.NopCloser(strings.NewReader(sse))placeholderplaceholder
	svc := openAIClientToolsTestService(upstream)
	account := &Account{ID: 5660, Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Credentials: map[string]any{"api_key": "test-key"placeholderplaceholder

	result, err := svc.forwardOpenAIPassthrough(context.Background(), c, account, body, body, "gpt-5.4", false, nil, true, time.Now())

placeholder
	require.NotNil(t, result)
	assertOpenAIClientToolsLowered(t, upstream.lastBody)
	output := recorder.Body.String()
	require.Contains(t, output, `"type":"custom_tool_call"`)
	require.Contains(t, output, `"type":"response.custom_tool_call_input.done"`)
	require.Contains(t, output, `"input":"*** Begin Patch"`)
	require.NotContains(t, output, `"input":{`)
placeholder
