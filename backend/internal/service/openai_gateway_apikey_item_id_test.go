//go:build unit

package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestOpenAIGatewayService_APIKeyPassthrough_StripsInvalidInputItemIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(
			`{"id":"resp_test","model":"gpt-5.6-sol","output":[],"usage":{"input_tokens":1,"output_tokens":1,"total_tokens":2placeholderplaceholder`,
		)),
placeholderplaceholder
	svc := newOpenAIImageGenerationControlTestService(upstream)
	c, _ := newOpenAIImageGenerationControlTestContext(true, "codex_cli_rs/0.144.1")
	account := newOpenAIImageGenerationControlTestAccount()
	account.Extra = map[string]any{"openai_passthrough": trueplaceholder

	body := []byte(`{
		"model":"gpt-5.6-sol",
		"stream":false,
		"input":[
			{"type":"message","id":"item_bad_message","role":"assistant","content":[{"type":"output_text","text":"hello"placeholder]placeholder,
			{"type":"function_call","id":"item_bad_call","call_id":"call_123","name":"exec_command","arguments":"{placeholder"placeholder,
			{"type":"message","id":"msg_valid","role":"user","content":[{"type":"input_text","text":"continue"placeholder]placeholder,
			{"type":"function_call","id":"fc_valid","call_id":"call_456","name":"apply_patch","arguments":"{placeholder"placeholder,
			{"type":"custom_tool_call","id":"fc_wrong_custom","call_id":"call_custom_1","name":"apply_patch","input":"patch"placeholder,
			{"type":"custom_tool_call","id":"ctc_valid","call_id":"call_custom_2","name":"apply_patch","input":"patch"placeholder,
			{"type":"tool_search_call","id":"fc_wrong_search","call_id":"call_search_1","arguments":{"query":"docs"placeholderplaceholder,
			{"type":"tool_search_call","id":"tsc_valid","call_id":"call_search_2","arguments":{"query":"docs"placeholderplaceholder,
			{"type":"function_call_output","id":"item_output","call_id":"call_123","output":"done"placeholder,
			{"type":"web_search_call","id":"item_wrong_web"placeholder
		]
placeholder`)

	result, err := svc.Forward(context.Background(), c, account, body)
placeholder
	require.NotNil(t, result)
	require.NotNil(t, upstream.lastReq)

	forwarded := upstream.lastBody
	require.False(t, gjson.GetBytes(forwarded, "input.0.id").Exists())
	require.Equal(t, "hello", gjson.GetBytes(forwarded, "input.0.content.0.text").String())
	require.False(t, gjson.GetBytes(forwarded, "input.1.id").Exists())
	require.Equal(t, "call_123", gjson.GetBytes(forwarded, "input.1.call_id").String())
	require.Equal(t, "exec_command", gjson.GetBytes(forwarded, "input.1.name").String())
	require.Equal(t, "{placeholder", gjson.GetBytes(forwarded, "input.1.arguments").String())
	require.Equal(t, "msg_valid", gjson.GetBytes(forwarded, "input.2.id").String())
	require.Equal(t, "fc_valid", gjson.GetBytes(forwarded, "input.3.id").String())
	require.False(t, gjson.GetBytes(forwarded, "input.4.id").Exists())
	require.Equal(t, "ctc_valid", gjson.GetBytes(forwarded, "input.5.id").String())
	require.False(t, gjson.GetBytes(forwarded, "input.6.id").Exists())
	require.Equal(t, "tsc_valid", gjson.GetBytes(forwarded, "input.7.id").String())
	require.Equal(t, "item_output", gjson.GetBytes(forwarded, "input.8.id").String())
	require.Equal(t, "call_123", gjson.GetBytes(forwarded, "input.8.call_id").String())
	require.False(t, gjson.GetBytes(forwarded, "input.9.id").Exists())
placeholder

func TestOpenAIGatewayService_OAuthPassthrough_SanitizesNativeToolItemIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	for _, accountType := range []string{AccountTypeOAuth, AccountTypeSetupTokenplaceholder {
		t.Run(accountType, func(t *testing.T) {
			upstreamSSE := "data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_test\",\"model\":\"gpt-5.6-sol\",\"output\":[],\"usage\":{\"input_tokens\":1,\"output_tokens\":1,\"total_tokens\":2placeholderplaceholderplaceholder\n\ndata: [DONE]\n\n"
			upstream := &httpUpstreamRecorder{resp: &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
				Body:       io.NopCloser(strings.NewReader(upstreamSSE)),
		placeholderplaceholder
			svc := newOpenAIImageGenerationControlTestService(upstream)
			c, _ := newOpenAIImageGenerationControlTestContext(true, "codex_cli_rs/0.144.1")
			account := newOpenAIImageGenerationControlTestAccount()
			account.Type = accountType
			account.Credentials = map[string]any{
				"access_token":       "oauth-token",
				"chatgpt_account_id": "chatgpt-account",
		placeholder
			account.Extra = map[string]any{"openai_passthrough": trueplaceholder

			body := []byte(`{
		"model":"gpt-5.6-sol",
		"stream":true,
		"instructions":"test",
		"input":[
			{"type":"custom_tool_call","id":"fc_wrong_custom","call_id":"call_custom_1","name":"apply_patch","input":"patch"placeholder,
			{"type":"custom_tool_call","id":"ctc_valid","call_id":"call_custom_2","name":"apply_patch","input":"patch"placeholder,
			{"type":"tool_search_call","id":"fc_wrong_search","call_id":"call_search_1","arguments":{"query":"docs"placeholderplaceholder,
			{"type":"tool_search_call","id":"tsc_valid","call_id":"call_search_2","arguments":{"query":"docs"placeholderplaceholder
		]
placeholder`)

			result, err := svc.Forward(context.Background(), c, account, body)
		placeholder
			require.NotNil(t, result)
			require.NotNil(t, upstream.lastReq)
			require.Equal(t, "https://chatgpt.com/backend-api/codex/responses", upstream.lastReq.URL.String())
			require.False(t, gjson.GetBytes(upstream.lastBody, "input.0.id").Exists())
			require.Equal(t, "ctc_valid", gjson.GetBytes(upstream.lastBody, "input.1.id").String())
			require.False(t, gjson.GetBytes(upstream.lastBody, "input.2.id").Exists())
			require.Equal(t, "tsc_valid", gjson.GetBytes(upstream.lastBody, "input.3.id").String())
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayService_SetupTokenLegacy_SanitizesAndTransforms(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstreamSSE := "data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_test\",\"model\":\"gpt-5.6-sol\",\"output\":[],\"usage\":{\"input_tokens\":1,\"output_tokens\":1,\"total_tokens\":2placeholderplaceholderplaceholder\n\ndata: [DONE]\n\n"
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamSSE)),
placeholderplaceholder
	svc := newOpenAIImageGenerationControlTestService(upstream)
	c, _ := newOpenAIImageGenerationControlTestContext(true, "codex_cli_rs/0.144.1")
	account := newOpenAIImageGenerationControlTestAccount()
	account.Type = AccountTypeSetupToken
	account.Credentials = map[string]any{
		"access_token":       "setup-token",
		"chatgpt_account_id": "chatgpt-account",
placeholder
	account.Extra = map[string]any{"openai_passthrough": falseplaceholder
	body := []byte(`{
		"model":"gpt-5.6-sol",
		"stream":true,
		"store":true,
		"reasoning":{"mode":"pro"placeholder,
		"instructions":"test",
		"input":[
			{"type":"custom_tool_call","id":"fc_wrong_custom","call_id":"call_custom_1","name":"apply_patch","input":"patch"placeholder,
			{"type":"function_call_output","call_id":"call_orphan","output":"orphan"placeholder
		]
placeholder`)

	result, err := svc.Forward(context.Background(), c, account, body)
placeholder
	require.NotNil(t, result)
	require.NotNil(t, upstream.lastReq)
	require.Equal(t, "https://chatgpt.com/backend-api/codex/responses", upstream.lastReq.URL.String())
	require.Equal(t, false, gjson.GetBytes(upstream.lastBody, "store").Bool())
	require.Equal(t, "max", gjson.GetBytes(upstream.lastBody, "reasoning.effort").String())
	require.False(t, gjson.GetBytes(upstream.lastBody, "reasoning.mode").Exists())
	require.False(t, gjson.GetBytes(upstream.lastBody, "input.0.id").Exists())
	require.Len(t, gjson.GetBytes(upstream.lastBody, "input").Array(), 1)
placeholder

// TestOpenAIGatewayService_APIKeyPassthrough_StripsInvalidReasoningItemIDs
// verifies that reasoning items with a non-rs id (e.g. item_*) are stripped
// before forwarding. OpenAI upstream requires reasoning ids to begin with
// "rs" and rejects item_* with 400:
// "Expected an ID that begins with 'rs'." (#5410)
func TestOpenAIGatewayService_APIKeyPassthrough_StripsInvalidReasoningItemIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(
			`{"id":"resp_test","model":"gpt-5.6-sol","output":[],"usage":{"input_tokens":1,"output_tokens":1,"total_tokens":2placeholderplaceholder`,
		)),
placeholderplaceholder
	svc := newOpenAIImageGenerationControlTestService(upstream)
	c, _ := newOpenAIImageGenerationControlTestContext(true, "codex_cli_rs/0.144.1")
	account := newOpenAIImageGenerationControlTestAccount()
	account.Extra = map[string]any{"openai_passthrough": trueplaceholder

	body := []byte(`{
		"model":"gpt-5.6-sol",
		"stream":false,
		"input":[
			{"type":"reasoning","id":"item_bad_reasoning","summary":[]placeholder,
			{"type":"reasoning","id":"rs_valid","summary":[]placeholder,
			{"type":"message","id":"msg_valid","role":"user","content":[{"type":"input_text","text":"continue"placeholder]placeholder
		]
placeholder`)

	result, err := svc.Forward(context.Background(), c, account, body)
placeholder
	require.NotNil(t, result)
	require.NotNil(t, upstream.lastReq)

	forwarded := upstream.lastBody
	require.False(t, gjson.GetBytes(forwarded, "input.0.id").Exists(),
		"item_* id should be stripped from reasoning")
	require.Equal(t, "rs_valid", gjson.GetBytes(forwarded, "input.1.id").String(),
		"valid rs* id must be preserved")
	require.Equal(t, "msg_valid", gjson.GetBytes(forwarded, "input.2.id").String())
placeholder

func TestShouldStripOpenAIResponsesInputItemID_Reasoning(t *testing.T) {
	cases := []struct {
		name     string
		itemType string
		id       string
		want     bool
placeholder{
		{"reasoning item_* id", "reasoning", "item_bad_reasoning", trueplaceholder,
		{"reasoning rs id", "reasoning", "rs_abc123", falseplaceholder,
		{"reasoning empty id", "reasoning", "", trueplaceholder,
		{"message msg id", "message", "msg_abc", falseplaceholder,
		{"message item id", "message", "item_x", trueplaceholder,
		{"function_call fc id", "function_call", "fc_abc", falseplaceholder,
		{"function_call ctc id", "function_call", "ctc_abc", trueplaceholder,
		{"function_call item id", "function_call", "item_x", trueplaceholder,
		{"custom tool ctc id", "custom_tool_call", "ctc_abc", falseplaceholder,
		{"custom tool fc id", "custom_tool_call", "fc_abc", trueplaceholder,
		{"tool search tsc id", "tool_search_call", "tsc_abc", falseplaceholder,
		{"tool search fc id", "tool_search_call", "fc_abc", trueplaceholder,
		{"web search ws id", "web_search_call", "ws_001", falseplaceholder,
		{"web search item id", "web_search_call", "item_001", trueplaceholder,
		{"custom output fc id", "custom_tool_call_output", "fc_001", falseplaceholder,
		{"custom output ctco id", "custom_tool_call_output", "ctco_001", trueplaceholder,
placeholder
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, shouldStripOpenAIResponsesInputItemID(tc.itemType, tc.id))
	placeholder)
placeholder
placeholder

func TestSanitizeOpenAIResponsesInputItemIDs_AllocationGrowthIsLinear(t *testing.T) {
	makeBody := func(itemCount int) []byte {
		items := make([]string, itemCount)
		for i := range items {
			items[i] = fmt.Sprintf(`{"type":"message","id":"item_%d","role":"user","content":[{"type":"input_text","text":"hello"placeholder]placeholder`, i)
	placeholder
		return []byte(`{"model":"gpt-5.6-sol","input":[` + strings.Join(items, ",") + `]placeholder`)
placeholder
	allocatedBytes := func(body []byte) uint64 {
		runtime.GC()
		var before, after runtime.MemStats
		runtime.ReadMemStats(&before)
		sanitized, changed, err := sanitizeOpenAIResponsesInputItemIDs(body)
		runtime.ReadMemStats(&after)
	placeholder
		require.True(t, changed)
		require.NotEmpty(t, sanitized)
		return after.TotalAlloc - before.TotalAlloc
placeholder

	smallAllocated := allocatedBytes(makeBody(20))
	largeAllocated := allocatedBytes(makeBody(200))
	require.Less(t, largeAllocated, smallAllocated*30,
		"10x more input items must not cause quadratic whole-body allocation growth")
placeholder

func TestNormalizeOpenAIResponsesWebSocketCompatibilityBodyClosesInvalidIDReferences(t *testing.T) {
	body := []byte(`{"type":"response.create","input":[
		{"type":"custom_tool_call","id":"ctc_call","call_id":"call_custom","name":"apply_patch","input":"patch"placeholder,
		{"type":"custom_tool_call_output","id":"ctco_bad","call_id":"call_custom","output":"done"placeholder,
		{"type":"item_reference","id":"ctco_bad"placeholder,
		{"type":"future_item","id":"item_future","payload":"keep"placeholder
	]placeholder`)

	for _, accountType := range []string{AccountTypeAPIKey, AccountTypeOAuthplaceholder {
		t.Run(accountType, func(t *testing.T) {
			normalized, changed, err := normalizeOpenAIResponsesWebSocketCompatibilityBody(body, &Account{
				Platform: PlatformOpenAI,
				Type:     accountType,
		placeholder)

		placeholder
			require.True(t, changed)
			require.Len(t, gjson.GetBytes(normalized, "input").Array(), 3)
			require.Equal(t, "ctc_call", gjson.GetBytes(normalized, "input.0.id").String())
			require.Equal(t, "call_custom", gjson.GetBytes(normalized, "input.1.call_id").String())
			require.False(t, gjson.GetBytes(normalized, "input.1.id").Exists())
			require.Equal(t, "item_future", gjson.GetBytes(normalized, "input.2.id").String())

			second, changedAgain, err := normalizeOpenAIResponsesWebSocketCompatibilityBody(normalized, &Account{
				Platform: PlatformOpenAI,
				Type:     accountType,
		placeholder)
		placeholder
			require.False(t, changedAgain)
			require.JSONEq(t, string(normalized), string(second))
	placeholder)
placeholder
placeholder
