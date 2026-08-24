//go:build unit

package service

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestNormalizeOpenAIResponsesLiteTools_MovesNamespacesAndKeepsSupportedTools(t *testing.T) {
	reqBody := map[string]any{
		"model": "gpt-5.6-terra",
		"tools": []any{
			map[string]any{"type": "function", "name": "shell"placeholder,
			map[string]any{"type": "custom", "name": "exec"placeholder,
			map[string]any{"type": "tool_search"placeholder,
			map[string]any{"type": "namespace", "name": "collaboration", "tools": []any{
				map[string]any{"type": "function", "name": "spawn_agent"placeholder,
	placeholder
	placeholder,
		"input": []any{
			map[string]any{"type": "message", "role": "user", "content": "hello"placeholder,
			map[string]any{"type": "additional_tools", "role": "developer", "tools": []any{
				map[string]any{"type": "namespace", "name": "image_gen"placeholder,
				map[string]any{"type": "namespace", "name": "collaboration", "tools": []any{
					map[string]any{"type": "function", "name": "spawn_agent"placeholder,
		placeholder
	placeholder
	placeholder,
		"tool_choice": map[string]any{"type": "namespace", "name": "collaboration"placeholder,
placeholder

	changed, err := normalizeOpenAIResponsesLiteTools(reqBody)

placeholder
	require.True(t, changed)
	tools := reqBody["tools"].([]any)
	require.Len(t, tools, 3)
	require.Equal(t, "function", tools[0].(map[string]any)["type"])
	require.Equal(t, "custom", tools[1].(map[string]any)["type"])
	require.Equal(t, "tool_search", tools[2].(map[string]any)["type"])
	input := reqBody["input"].([]any)
	require.Len(t, input, 2)
	additional := input[1].(map[string]any)["tools"].([]any)
	require.Len(t, additional, 2)
	require.Equal(t, "image_gen", additional[0].(map[string]any)["name"])
	require.Equal(t, "collaboration", additional[1].(map[string]any)["name"], "existing namespace must not be duplicated")
	require.Equal(t, map[string]any{"type": "namespace", "name": "collaboration"placeholder, reqBody["tool_choice"])
placeholder

func TestNormalizeOpenAIResponsesLiteTools_PreservesDeferredFlagsWithToolSearch(t *testing.T) {
	reqBody := map[string]any{
		"tools": []any{
			map[string]any{"type": "tool_search"placeholder,
			map[string]any{"type": "function", "name": "shell", "defer_loading": trueplaceholder,
	placeholder,
placeholder

	_, err := normalizeOpenAIResponsesLiteTools(reqBody)
placeholder
	tools := reqBody["tools"].([]any)
	require.Equal(t, "tool_search", tools[0].(map[string]any)["type"])
	require.Equal(t, true, tools[1].(map[string]any)["defer_loading"])
placeholder

func TestNormalizeOpenAIResponsesLiteTools_RejectsConflictingAdditionalTool(t *testing.T) {
	reqBody := map[string]any{
		"tools": []any{map[string]any{
			"type":  "namespace",
			"name":  "collaboration",
			"tools": []any{map[string]any{"type": "function", "name": "spawn_agent"placeholderplaceholder,
placeholder
		"input": []any{map[string]any{
			"type": "additional_tools",
			"tools": []any{map[string]any{
				"type":  "namespace",
				"name":  "collaboration",
				"tools": []any{map[string]any{"type": "function", "name": "send_message"placeholderplaceholder,
	placeholder
placeholder
placeholder

	changed, err := normalizeOpenAIResponsesLiteTools(reqBody)

	require.ErrorContains(t, err, `conflicts with migrated tool type "namespace" name "collaboration"`)
	require.False(t, changed)
	require.Len(t, reqBody["tools"], 1, "conflicts must not partially remove top-level tools")
placeholder

func TestNormalizeOpenAIResponsesLiteTools_DeduplicatesAcrossAdditionalToolItems(t *testing.T) {
	namespace := map[string]any{
		"type":  "namespace",
		"name":  "collaboration",
		"tools": []any{map[string]any{"type": "function", "name": "spawn_agent"placeholderplaceholder,
placeholder
	reqBody := map[string]any{
		"tools": []any{namespaceplaceholder,
		"input": []any{
			map[string]any{
				"type":  "additional_tools",
				"tools": []any{map[string]any{"type": "custom", "name": "exec"placeholderplaceholder,
		placeholder,
			map[string]any{
				"type":  "additional_tools",
				"tools": []any{namespaceplaceholder,
		placeholder,
	placeholder,
placeholder

	changed, err := normalizeOpenAIResponsesLiteTools(reqBody)

placeholder
	require.True(t, changed)
	require.NotContains(t, reqBody, "tools")
	input := reqBody["input"].([]any)
	require.Len(t, input[0].(map[string]any)["tools"], 1)
	require.Len(t, input[1].(map[string]any)["tools"], 1)
placeholder

func TestNormalizeOpenAIResponsesLiteTools_ConvertsStringInput(t *testing.T) {
	reqBody := map[string]any{
		"input": "hello",
		"tools": []any{map[string]any{
			"type": "namespace",
			"name": "collaboration",
placeholder
placeholder

	changed, err := normalizeOpenAIResponsesLiteTools(reqBody)

placeholder
	require.True(t, changed)
	require.NotContains(t, reqBody, "tools")
	input := reqBody["input"].([]any)
	require.Len(t, input, 2)
	require.Equal(t, "message", input[0].(map[string]any)["type"])
	require.Equal(t, "hello", input[0].(map[string]any)["content"])
	require.Equal(t, "additional_tools", input[1].(map[string]any)["type"])
placeholder

func TestNormalizeOpenAIResponsesLiteTools_KeepsSupportedTopLevelTools(t *testing.T) {
	reqBody := map[string]any{
		"reasoning": map[string]any{"context": "all_turns"placeholder,
		"tools": []any{
			map[string]any{"type": "function", "name": "shell"placeholder,
			map[string]any{"type": "custom", "name": "exec"placeholder,
			map[string]any{"type": "tool_search"placeholder,
			"custom shorthand",
	placeholder,
placeholder

	changed, err := normalizeOpenAIResponsesLiteTools(reqBody)

placeholder
	require.False(t, changed)
	require.Len(t, reqBody["tools"], 4)
placeholder

func TestNormalizeOpenAIResponsesLiteTools_EnsuresReasoningContext(t *testing.T) {
	tests := []struct {
		name      string
		reasoning any
placeholder{
		{name: "missing"placeholder,
		{name: "missing context", reasoning: map[string]any{"effort": "high"placeholderplaceholder,
		{name: "wrong context", reasoning: map[string]any{"effort": "medium", "context": "current_turn"placeholderplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := map[string]any{"input": "hello"placeholder
			if tt.reasoning != nil {
				reqBody["reasoning"] = tt.reasoning
		placeholder

			changed, err := normalizeOpenAIResponsesLiteTools(reqBody)

		placeholder
			require.True(t, changed)
			reasoning := reqBody["reasoning"].(map[string]any)
			require.Equal(t, "all_turns", reasoning["context"])
			if tt.name != "missing" {
				require.Equal(t, tt.reasoning.(map[string]any)["effort"], reasoning["effort"])
		placeholder
	placeholder)
placeholder
placeholder

func TestNormalizeOpenAIResponsesLiteTools_RejectsNonObjectReasoning(t *testing.T) {
	reqBody := map[string]any{"reasoning": "high"placeholder

	changed, err := normalizeOpenAIResponsesLiteTools(reqBody)

	require.ErrorContains(t, err, "reasoning to be an object")
	require.False(t, changed)
	require.Equal(t, "high", reqBody["reasoning"])
placeholder

func TestNormalizeOpenAIResponsesLiteTools_RejectsUnsupportedTools(t *testing.T) {
	tests := []struct {
		name string
		tool map[string]any
		want string
placeholder{
		{name: "hosted web search", tool: map[string]any{"type": "web_search"placeholder, want: `top-level tool type "web_search"`placeholder,
		{name: "hosted image generation", tool: map[string]any{"type": "image_generation"placeholder, want: `top-level tool type "image_generation"`placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := map[string]any{"tools": []any{tt.toolplaceholderplaceholder
			changed, err := normalizeOpenAIResponsesLiteTools(reqBody)
			require.ErrorContains(t, err, tt.want)
			require.False(t, changed)
			require.Len(t, reqBody["tools"], 1, "validation errors must not partially mutate tools")
	placeholder)
placeholder
placeholder

func TestNormalizeOpenAIResponsesLiteToolsPayload_PreservesResponseCreateShape(t *testing.T) {
	body := []byte(`{
		"type":"response.create",
		"model":"gpt-5.6-terra",
		"client_metadata":{"ws_request_header_x_openai_internal_codex_responses_lite":"true"placeholder,
		"input":[{"type":"message","role":"user","content":"hello"placeholder],
		"tools":[{"type":"namespace","name":"collaboration","tools":[{"type":"function","name":"spawn_agent"placeholder]placeholder],
		"tool_choice":{"type":"namespace","name":"collaboration"placeholder
placeholder`)

	updated, changed, err := normalizeOpenAIResponsesLiteToolsPayload(body)

placeholder
	require.True(t, changed)
	require.Equal(t, "response.create", gjson.GetBytes(updated, "type").String())
	require.False(t, gjson.GetBytes(updated, "tools").Exists())
	require.Equal(t, "collaboration", gjson.GetBytes(updated, `input.#(type=="additional_tools").tools.0.name`).String())
	require.Equal(t, "namespace", gjson.GetBytes(updated, "tool_choice.type").String())
placeholder

func TestApplyCodexOAuthTransform_PreservesLiteNamespaceToolChoice(t *testing.T) {
	reqBody := map[string]any{
		"model": "gpt-5.6-terra",
		"input": []any{map[string]any{
			"type": "additional_tools",
			"tools": []any{map[string]any{
				"type": "namespace",
				"name": "collaboration",
	placeholder
placeholder
		"tool_choice": map[string]any{"type": "namespace", "name": "collaboration"placeholder,
placeholder

	applyCodexOAuthTransform(reqBody, true, false)

	require.Equal(t, map[string]any{"type": "namespace", "name": "collaboration"placeholder, reqBody["tool_choice"])
placeholder

func TestOpenAIGatewayServiceForward_NormalizesResponsesLiteToolsForOAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	for _, passthrough := range []bool{false, trueplaceholder {
		name := "managed"
		if passthrough {
			name = "passthrough"
	placeholder
		t.Run(name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(nil))
			c.Request.Header.Set("User-Agent", "codex_cli_rs/0.144.1")
			c.Request.Header.Set(responsesLiteHeader, "true")
			upstream := &httpUpstreamRecorder{resp: &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
				Body: io.NopCloser(strings.NewReader(
					"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_lite\",\"usage\":{\"input_tokens\":1,\"output_tokens\":1placeholderplaceholderplaceholder\n\n" +
						"data: [DONE]\n\n",
				)),
		placeholderplaceholder
			svc := &OpenAIGatewayService{cfg: &config.Config{placeholder, httpUpstream: upstreamplaceholder
			account := &Account{
				ID: 501, Name: "responses-lite", Platform: PlatformOpenAI, Type: AccountTypeOAuth,
				Concurrency: 1, Status: StatusActive, Schedulable: true, RateMultiplier: f64p(1),
		placeholder"access_token": "oauth-token", "chatgpt_account_id": "chatgpt-account"placeholder,
				Extra:       map[string]any{"openai_passthrough": passthroughplaceholder,
		placeholder
			body := []byte(`{
				"model":"gpt-5.6-terra","stream":true,"instructions":"test",
				"reasoning":{"effort":"high","context":"current_turn"placeholder,
				"tools":[
					{"type":"function","name":"shell","parameters":{"type":"object"placeholderplaceholder,
					{"type":"custom","name":"exec"placeholder,
					{"type":"tool_search"placeholder,
					{"type":"namespace","name":"collaboration","tools":[{"type":"function","name":"spawn_agent","parameters":{"type":"object"placeholderplaceholder]placeholder
				],
				"input":[{"type":"message","role":"user","content":"hello"placeholder],
				"tool_choice":{"type":"namespace","name":"collaboration"placeholder
		placeholder`)

			result, err := svc.Forward(context.Background(), c, account, body)

		placeholder
			require.NotNil(t, result)
			require.Equal(t, "true", upstream.lastReq.Header.Get(responsesLiteHeader))
			require.Equal(t, "high", gjson.GetBytes(upstream.lastBody, "reasoning.effort").String())
			require.Equal(t, "all_turns", gjson.GetBytes(upstream.lastBody, "reasoning.context").String())
			require.False(t, gjson.GetBytes(upstream.lastBody, `tools.#(type=="namespace")`).Exists())
			require.Equal(t, "shell", gjson.GetBytes(upstream.lastBody, `tools.#(type=="function").name`).String())
			require.Equal(t, "exec", gjson.GetBytes(upstream.lastBody, `tools.#(type=="custom").name`).String())
			require.True(t, gjson.GetBytes(upstream.lastBody, `tools.#(type=="tool_search")`).Exists())
			require.Equal(t, "collaboration", gjson.GetBytes(upstream.lastBody, `input.#(type=="additional_tools").tools.0.name`).String())
			require.Equal(t, "namespace", gjson.GetBytes(upstream.lastBody, "tool_choice.type").String())
			require.Equal(t, "collaboration", gjson.GetBytes(upstream.lastBody, "tool_choice.name").String())
	placeholder)
placeholder
placeholder
