package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai_compat"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestOpenAIResponsesRejectedFieldRetryStateRejectsDuplicateBodyAndCap(t *testing.T) {
	initialBody := []byte(`{"model":"gpt-5.5"placeholder`)
	state := newOpenAIResponsesRejectedFieldRetryState(initialBody)

	require.False(t, state.Allow(initialBody))
	for attempt := 0; attempt < maxOpenAIResponsesRejectedFieldRetries; attempt++ {
		nextBody := []byte(fmt.Sprintf(`{"model":"gpt-5.5","variant":%dplaceholder`, attempt))
		require.True(t, state.Allow(nextBody))
		require.False(t, state.Allow(nextBody))
placeholder
	require.False(t, state.Allow([]byte(`{"model":"gpt-5.5","variant":"overflow"placeholder`)))
placeholder

func TestNormalizeOpenAIResponsesRejectedFieldRetryBodyRejectsAmbiguousErrors(t *testing.T) {
	tests := []struct {
		name         string
		body         []byte
		responseBody []byte
placeholder{
		{
			name:         "namespace belongs to message",
			body:         []byte(`{"input":[{"type":"message","namespace":"keep"placeholder]placeholder`),
			responseBody: []byte(`{"error":{"code":"unknown_parameter","message":"Unknown parameter: 'input[0].namespace'.","param":"input[0].namespace"placeholderplaceholder`),
	placeholder,
		{
			name:         "max output tokens only mentioned",
			body:         []byte(`{"max_output_tokens":placeholder`),
			responseBody: []byte(`{"error":{"code":"invalid_request_error","message":"max_output_tokens must be positive","param":"max_output_tokens"placeholderplaceholder`),
	placeholder,
		{
			name:         "structured param overrides namespace mention",
			body:         []byte(`{"input":[{"type":"function_call","namespace":"keep","arguments":"{placeholder"placeholder]placeholder`),
			responseBody: []byte(`{"error":{"code":"unknown_parameter","message":"Unknown parameter: 'input[0].namespace'.","param":"tools"placeholderplaceholder`),
	placeholder,
		{
			name:         "nested max output tokens param is not top level",
			body:         []byte(`{"max_output_tokens":4096,"input":[{"type":"message","content":{"max_output_tokens":"keep"placeholderplaceholder]placeholder`),
			responseBody: []byte(`{"error":{"code":"unknown_parameter","message":"Unknown parameter: input[0].content.max_output_tokens","param":"input[0].content.max_output_tokens"placeholderplaceholder`),
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryBody, _, changed, err := normalizeOpenAIResponsesRejectedFieldRetryBody(http.StatusBadRequest, tt.body, tt.responseBody)
		placeholder
			require.False(t, changed)
			require.Nil(t, retryBody)
	placeholder)
placeholder
placeholder

func TestNormalizeOpenAIResponsesRejectedFieldRetryBodyFindsNamespacePathInMessage(t *testing.T) {
	body := []byte(`{"input":[{"type":"function_call","namespace":"keep","arguments":"{placeholder"placeholder,{"type":"function_call","namespace":"remove","arguments":"{placeholder"placeholder]placeholder`)
	responseBody := []byte(`{"error":{"code":"unknown_parameter","message":"input[0] was accepted; Unknown parameter: 'input[1].namespace'."placeholderplaceholder`)

	retryBody, _, changed, err := normalizeOpenAIResponsesRejectedFieldRetryBody(http.StatusBadRequest, body, responseBody)

placeholder
	require.True(t, changed)
	require.Equal(t, "keep", gjson.GetBytes(retryBody, "input.0.namespace").String())
	require.False(t, gjson.GetBytes(retryBody, "input.1.namespace").Exists())
placeholder

func TestNormalizeOpenAIResponsesRejectedFieldRetryBodyBindsNamespacePathToRejectionPhrase(t *testing.T) {
	body := []byte(`{"input":[{"type":"function_call","namespace":"keep","arguments":"{placeholder"placeholder,{"type":"function_call","namespace":"remove","arguments":"{placeholder"placeholder]placeholder`)
	responseBody := []byte(`{"error":{"code":"unknown_parameter","message":"input[0].namespace is supported; Unknown parameter: input[1].namespace."placeholderplaceholder`)

	retryBody, _, changed, err := normalizeOpenAIResponsesRejectedFieldRetryBody(http.StatusBadRequest, body, responseBody)

placeholder
	require.True(t, changed)
	require.Equal(t, "keep", gjson.GetBytes(retryBody, "input.0.namespace").String())
	require.False(t, gjson.GetBytes(retryBody, "input.1.namespace").Exists())
placeholder

func TestNormalizeOpenAIResponsesRejectedFieldRetryBodyDoesNotTreatMaxOutputTokensSuggestionAsRejection(t *testing.T) {
	body := []byte(`{"max_tokens":4096,"max_output_tokens":placeholder`)
	responseBody := []byte(`{"error":{"code":"unknown_parameter","message":"Unknown parameter: max_tokens. Use max_output_tokens instead."placeholderplaceholder`)

	retryBody, _, changed, err := normalizeOpenAIResponsesRejectedFieldRetryBody(http.StatusBadRequest, body, responseBody)

placeholder
	require.False(t, changed)
	require.Nil(t, retryBody)
placeholder

func TestNormalizeOpenAIResponsesRejectedFieldRetryBodyBindsMaxOutputTokensToRejectionPhrase(t *testing.T) {
	body := []byte(`{"max_output_tokens":placeholder`)
	responseBody := []byte(`{"error":{"code":"unsupported_parameter","message":"Unsupported parameter: max_output_tokens."placeholderplaceholder`)

	retryBody, _, changed, err := normalizeOpenAIResponsesRejectedFieldRetryBody(http.StatusBadRequest, body, responseBody)

placeholder
	require.True(t, changed)
	require.False(t, gjson.GetBytes(retryBody, "max_output_tokens").Exists())
placeholder

func TestNormalizeOpenAIResponsesRejectedFieldRetryBodyRemovesExactIndexedStatus(t *testing.T) {
	body := []byte(`{"input":[{"type":"message","status":"keep","content":"one"placeholder,{"type":"reasoning","status":"remove","summary":[]placeholder]placeholder`)
	responses := []struct {
		name string
		body []byte
placeholder{
		{
			name: "structured param",
			body: []byte(`{"error":{"code":"unknown_parameter","message":"Unknown parameter: 'input[1].status'.","param":"input[1].status"placeholderplaceholder`),
	placeholder,
		{
			name: "message param",
			body: []byte(`{"error":{"code":"unsupported_parameter","message":"Unsupported parameter: input[1].status."placeholderplaceholder`),
	placeholder,
placeholder
	for _, tt := range responses {
		t.Run(tt.name, func(t *testing.T) {
			retryBody, _, changed, err := normalizeOpenAIResponsesRejectedFieldRetryBody(http.StatusBadRequest, body, tt.body)
		placeholder
			require.True(t, changed)
			require.Equal(t, "keep", gjson.GetBytes(retryBody, "input.0.status").String())
			require.False(t, gjson.GetBytes(retryBody, "input.1.status").Exists())
	placeholder)
placeholder
placeholder

func TestNormalizeOpenAIResponsesRejectedFieldRetryBodyNormalizesExactNullContent(t *testing.T) {
	tests := []struct {
		name       string
		body       []byte
		wantChange bool
		wantValue  string
		wantExists bool
placeholder{
		{
			name:       "message becomes empty string",
			body:       []byte(`{"input":[{"type":"message","role":"assistant","content":nullplaceholder]placeholder`),
			wantChange: true,
			wantValue:  "",
			wantExists: true,
	placeholder,
		{
			name:       "reasoning content is removed",
			body:       []byte(`{"input":[{"type":"reasoning","content":null,"summary":[]placeholder]placeholder`),
			wantChange: true,
			wantExists: false,
	placeholder,
		{
			name:       "unknown item is unchanged",
			body:       []byte(`{"input":[{"type":"future_item","content":nullplaceholder]placeholder`),
			wantChange: false,
	placeholder,
		{
			name:       "non null content is unchanged",
			body:       []byte(`{"input":[{"type":"message","content":"keep"placeholder]placeholder`),
			wantChange: false,
	placeholder,
placeholder
	responseBody := []byte(`{"error":{"code":"invalid_type","message":"Invalid type for 'input[0].content': expected one of a string or a list of input items, but got null instead.","param":"input[0].content"placeholderplaceholder`)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryBody, _, changed, err := normalizeOpenAIResponsesRejectedFieldRetryBody(http.StatusBadRequest, tt.body, responseBody)
		placeholder
			require.Equal(t, tt.wantChange, changed)
			if !tt.wantChange {
				require.Nil(t, retryBody)
				return
		placeholder
			content := gjson.GetBytes(retryBody, "input.0.content")
			require.Equal(t, tt.wantExists, content.Exists())
			if tt.wantExists {
				require.Equal(t, tt.wantValue, content.String())
				require.Equal(t, gjson.String, content.Type)
		placeholder
	placeholder)
placeholder
placeholder

func TestNormalizeOpenAIResponsesRejectedFieldRetryBodyRejectsUnsafeIndexedMutations(t *testing.T) {
	tests := []struct {
		name         string
		body         []byte
		responseBody []byte
placeholder{
		{
			name:         "nested status path",
			body:         []byte(`{"input":[{"type":"message","content":{"status":"keep"placeholderplaceholder]placeholder`),
			responseBody: []byte(`{"error":{"code":"unknown_parameter","message":"Unknown parameter: input[0].content.status.","param":"input[0].content.status"placeholderplaceholder`),
	placeholder,
		{
			name:         "status index out of bounds",
			body:         []byte(`{"input":[{"type":"message","status":"keep"placeholder]placeholder`),
			responseBody: []byte(`{"error":{"code":"unknown_parameter","message":"Unknown parameter: input[4].status.","param":"input[4].status"placeholderplaceholder`),
	placeholder,
		{
			name:         "status path only mentioned",
			body:         []byte(`{"input":[{"type":"message","status":"keep"placeholder]placeholder`),
			responseBody: []byte(`{"error":{"code":"invalid_request_error","message":"input[0].status must be completed","param":"input[0].status"placeholderplaceholder`),
	placeholder,
		{
			name:         "content param and message disagree",
			body:         []byte(`{"input":[{"type":"message","content":nullplaceholder,{"type":"message","content":nullplaceholder]placeholder`),
			responseBody: []byte(`{"error":{"code":"invalid_type","message":"Invalid type for input[1].content: got null instead.","param":"input[0].content"placeholderplaceholder`),
	placeholder,
		{
			name:         "content error only mentions null",
			body:         []byte(`{"input":[{"type":"message","content":nullplaceholder]placeholder`),
			responseBody: []byte(`{"error":{"code":"invalid_type","message":"content cannot be null","param":"input[0].content"placeholderplaceholder`),
	placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryBody, _, changed, err := normalizeOpenAIResponsesRejectedFieldRetryBody(http.StatusBadRequest, tt.body, tt.responseBody)
		placeholder
			require.False(t, changed)
			require.Nil(t, retryBody)
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayService_OAuthRetriesExactRejectedStatus(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5","stream":true,"instructions":"test","input":[{"type":"message","role":"user","status":"completed","content":"hello"placeholder]placeholder`)
	upstream := &httpUpstreamRecorder{responses: []*http.Response{
		newOpenAIRejectedFieldTestResponse(http.StatusBadRequest, `{"error":{"code":"unknown_parameter","message":"Unknown parameter: 'input[0].status'.","param":"input[0].status"placeholderplaceholder`),
		newOpenAIRejectedFieldTestResponse(http.StatusOK, "data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_ok\",\"output\":[],\"usage\":{\"input_tokens\":1,\"output_tokens\":1,\"total_tokens\":2placeholderplaceholderplaceholder\n\ndata: [DONE]\n\n"),
placeholderplaceholder
	upstream.responses[1].Header.Set("Content-Type", "text/event-stream")

	result, err := newOpenAIRejectedFieldTestService(upstream).Forward(
		context.Background(), newOpenAIRejectedFieldTestContext(body), newOpenAIOAuthNamespaceTestAccount(), body,
	)

placeholder
	require.NotNil(t, result)
	require.Len(t, upstream.bodies, 2)
	require.Equal(t, "completed", gjson.GetBytes(upstream.bodies[0], "input.0.status").String())
	require.False(t, gjson.GetBytes(upstream.bodies[1], "input.0.status").Exists())
placeholder

func TestOpenAIGatewayService_APIKeyRetriesExactRejectedNullMessageContent(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5","stream":false,"input":[{"type":"message","role":"assistant","content":nullplaceholder,{"type":"message","role":"user","content":"continue"placeholder]placeholder`)
	upstream := &httpUpstreamRecorder{responses: []*http.Response{
		newOpenAIRejectedFieldTestResponse(http.StatusBadRequest, `{"error":{"code":"invalid_type","message":"Invalid type for 'input[0].content': expected one of a string or a list of input items, but got null instead.","param":"input[0].content"placeholderplaceholder`),
		newOpenAIRejectedFieldTestResponse(http.StatusOK, `{"output":[],"usage":{"input_tokens":1,"output_tokens":1,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholder`),
placeholderplaceholder

	result, err := newOpenAIRejectedFieldTestService(upstream).Forward(
		context.Background(), newOpenAIRejectedFieldTestContext(body), newOpenAIRejectedFieldTestAccount(), body,
	)

placeholder
	require.NotNil(t, result)
	require.Len(t, upstream.bodies, 2)
	require.Equal(t, gjson.Null, gjson.GetBytes(upstream.bodies[0], "input.0.content").Type)
	require.Equal(t, gjson.String, gjson.GetBytes(upstream.bodies[1], "input.0.content").Type)
	require.Equal(t, "continue", gjson.GetBytes(upstream.bodies[1], "input.1.content").String())
placeholder

func TestNormalizeOpenAIResponsesRejectedFieldRetryBodyRemovesModelRejectedPromptCacheBreakpoint(t *testing.T) {
	tests := []struct {
		name         string
		body         []byte
		responseBody []byte
		removedPath  string
		preserved    string
		reason       string
placeholder{
		{
			name:         "top level",
			body:         []byte(`{"model":"gpt-5.6-sol","prompt_cache_breakpoint":{"type":"message_start"placeholder,"input":"hello"placeholder`),
			responseBody: []byte(`{"error":{"code":"invalid_parameter","message":"prompt_cache_breakpoint is not supported on this model","param":"prompt_cache_breakpoint"placeholderplaceholder`),
			removedPath:  "prompt_cache_breakpoint",
			preserved:    "input",
			reason:       "prompt_cache_breakpoint parameter rejection",
	placeholder,
		{
			name:         "indexed path from message",
			body:         []byte(`{"input":[{"type":"message","prompt_cache_breakpoint":{"type":"message_start"placeholderplaceholder,{"type":"message","prompt_cache_breakpoint":{"type":"message_end"placeholderplaceholder]placeholder`),
			responseBody: []byte(`{"error":{"code":"invalid_parameter","message":"input[1].prompt_cache_breakpoint is not supported on this model"placeholderplaceholder`),
			removedPath:  "input.1.prompt_cache_breakpoint",
			preserved:    "input.0.prompt_cache_breakpoint",
			reason:       "indexed prompt_cache_breakpoint parameter rejection",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryBody, reason, changed, err := normalizeOpenAIResponsesRejectedFieldRetryBody(http.StatusBadRequest, tt.body, tt.responseBody)

		placeholder
			require.True(t, changed)
			require.Equal(t, tt.reason, reason)
			require.False(t, gjson.GetBytes(retryBody, tt.removedPath).Exists())
			require.True(t, gjson.GetBytes(retryBody, tt.preserved).Exists())
	placeholder)
placeholder
placeholder

func TestNormalizeOpenAIResponsesRejectedFieldRetryBodyRejectsAmbiguousPromptCacheBreakpointErrors(t *testing.T) {
	body := []byte(`{"prompt_cache_breakpoint":{"type":"message_start"placeholder,"input":[{"type":"message","prompt_cache_breakpoint":{"type":"message_end"placeholderplaceholder]placeholder`)
	tests := []struct {
		name         string
		responseBody []byte
placeholder{
		{
			name:         "structured param disagrees",
			responseBody: []byte(`{"error":{"code":"invalid_parameter","message":"input[0].prompt_cache_breakpoint is not supported on this model","param":"prompt_cache_breakpoint"placeholderplaceholder`),
	placeholder,
		{
			name:         "index out of bounds",
			responseBody: []byte(`{"error":{"code":"invalid_parameter","message":"input[4].prompt_cache_breakpoint is not supported on this model","param":"input[4].prompt_cache_breakpoint"placeholderplaceholder`),
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryBody, _, changed, err := normalizeOpenAIResponsesRejectedFieldRetryBody(http.StatusBadRequest, body, tt.responseBody)

		placeholder
			require.False(t, changed)
			require.Nil(t, retryBody)
	placeholder)
placeholder
placeholder

func TestNormalizeOpenAIResponsesRejectedFieldRetryBodyAcceptsEitherCacheModelRejectionSignal(t *testing.T) {
	tests := []struct {
		name         string
		responseBody []byte
placeholder{
		{
			name:         "invalid parameter code",
			responseBody: []byte(`{"error":{"code":"invalid_parameter","message":"This optional cache hint cannot be used here","param":"prompt_cache_breakpoint"placeholderplaceholder`),
	placeholder,
		{
			name:         "model rejection message",
			responseBody: []byte(`{"error":{"code":"invalid_request_error","message":"prompt_cache_breakpoint is not supported on this model","param":"prompt_cache_breakpoint"placeholderplaceholder`),
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryBody, _, changed, err := normalizeOpenAIResponsesRejectedFieldRetryBody(http.StatusBadRequest, []byte(`{"prompt_cache_breakpoint":true,"input":"keep"placeholder`), tt.responseBody)

		placeholder
			require.True(t, changed)
			require.False(t, gjson.GetBytes(retryBody, "prompt_cache_breakpoint").Exists())
			require.Equal(t, "keep", gjson.GetBytes(retryBody, "input").String())
	placeholder)
placeholder
placeholder

func TestOpenAIResponsesRejectedFieldRetryStateAllowsPromptCacheBreakpointVariantOnce(t *testing.T) {
	body := []byte(`{"input":[{"prompt_cache_breakpoint":{"type":"message_start"placeholderplaceholder]placeholder`)
	responseBody := []byte(`{"error":{"code":"invalid_parameter","message":"input[0].prompt_cache_breakpoint is not supported on this model","param":"input[0].prompt_cache_breakpoint"placeholderplaceholder`)
	state := newOpenAIResponsesRejectedFieldRetryState(body)

	retryBody, _, changed, err := normalizeOpenAIResponsesRejectedFieldRetryBody(http.StatusBadRequest, body, responseBody)
placeholder
	require.True(t, changed)
	require.True(t, state.Allow(retryBody))
	require.False(t, state.Allow(retryBody))
placeholder

func TestOpenAIGatewayService_APIKeyStripsAllIndexedNamespacesBeforeFirstForward(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5","stream":false,"input":[{"type":"function_call","name":"first","namespace":"remove-first","arguments":"{placeholder"placeholder,{"type":"custom_tool_call","name":"second","namespace":"remove-second","input":"{placeholder"placeholder]placeholder`)
	upstream := &httpUpstreamRecorder{responses: []*http.Response{
		newOpenAIRejectedFieldTestResponse(http.StatusOK, `{"output":[],"usage":{"input_tokens":1,"output_tokens":1,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholder`),
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
	require.False(t, gjson.GetBytes(upstream.bodies[0], "input.0.namespace").Exists())
	require.False(t, gjson.GetBytes(upstream.bodies[0], "input.1.namespace").Exists())
placeholder

func TestOpenAIGatewayService_OpenAIHTTPStripsInputNamespacesBeforeFirstForward(t *testing.T) {
	accounts := []struct {
		name    string
		account *Account
placeholder{
		{name: "oauth", account: newOpenAIOAuthNamespaceTestAccount()placeholder,
		{name: "apikey", account: newOpenAIRejectedFieldTestAccount()placeholder,
placeholder
	for _, tt := range accounts {
		for _, path := range []string{"/v1/responses", "/v1/responses/compact"placeholder {
			t.Run(tt.name+path, func(t *testing.T) {
				body := []byte(`{"model":"gpt-5.5","stream":false,"instructions":"test","input":[{"type":"message","role":"user","namespace":"remove","content":[{"type":"input_text","text":"hello","namespace":"nested-keep"placeholder]placeholder]placeholder`)
				upstream := &httpUpstreamRecorder{responses: []*http.Response{
					newOpenAIRejectedFieldTestResponse(http.StatusOK, `{"id":"resp_namespace_ok","output":[],"usage":{"input_tokens":1,"output_tokens":1,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholder`),
			placeholderplaceholder
				c := newOpenAIRejectedFieldTestContext(body)
				c.Request.URL.Path = path

				result, err := newOpenAIRejectedFieldTestService(upstream).Forward(
					context.Background(),
					c,
					tt.account,
					body,
				)

			placeholder
				require.NotNil(t, result)
				require.Len(t, upstream.bodies, 1, "namespace must be removed before the first upstream request")
				require.False(t, gjson.GetBytes(upstream.bodies[0], "input.0.namespace").Exists())
				require.Equal(t, "nested-keep", gjson.GetBytes(upstream.bodies[0], "input.0.content.0.namespace").String())
		placeholder)
	placeholder
placeholder
placeholder

func TestOpenAIGatewayService_RetriesExplicitMaxOutputTokensRejection(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5","stream":false,"max_output_tokens":4096,"input":[{"type":"message","role":"user","content":{"max_output_tokens":"keep"placeholderplaceholder]placeholder`)
	upstream := &httpUpstreamRecorder{responses: []*http.Response{
		newOpenAIRejectedFieldTestResponse(http.StatusBadRequest, `{"error":{"code":"unsupported_parameter","message":"Unsupported parameter: max_output_tokens","param":"max_output_tokens","type":"invalid_request_error"placeholderplaceholder`),
		newOpenAIRejectedFieldTestResponse(http.StatusOK, `{"output":[],"usage":{"input_tokens":1,"output_tokens":1,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholder`),
placeholderplaceholder

	result, err := newOpenAIRejectedFieldTestService(upstream).Forward(
		context.Background(),
		newOpenAIRejectedFieldTestContext(body),
		newOpenAIRejectedFieldTestAccount(),
		body,
	)

placeholder
	require.NotNil(t, result)
	require.Len(t, upstream.bodies, 2)
	require.Equal(t, int64(4096), gjson.GetBytes(upstream.bodies[0], "max_output_tokens").Int())
	require.False(t, gjson.GetBytes(upstream.bodies[1], "max_output_tokens").Exists())
	require.Equal(t, "keep", gjson.GetBytes(upstream.bodies[1], "input.0.content.max_output_tokens").String())
placeholder

func TestOpenAIGatewayService_ComposesProactiveNamespaceStripWithRejectedFieldRetry(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5","stream":false,"max_output_tokens":2048,"input":[{"type":"function_call","name":"first","namespace":"remove-first","arguments":"{placeholder"placeholder,{"type":"custom_tool_call","name":"second","namespace":"remove-second","input":"{placeholder"placeholder]placeholder`)
	upstream := &httpUpstreamRecorder{responses: []*http.Response{
		newOpenAIRejectedFieldTestResponse(http.StatusBadRequest, `{"error":{"code":"unsupported_parameter","message":"Unsupported parameter: max_output_tokens","param":"max_output_tokens"placeholderplaceholder`),
		newOpenAIRejectedFieldTestResponse(http.StatusOK, `{"output":[],"usage":{"input_tokens":1,"output_tokens":1,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholder`),
placeholderplaceholder

	result, err := newOpenAIRejectedFieldTestService(upstream).Forward(
		context.Background(),
		newOpenAIRejectedFieldTestContext(body),
		newOpenAIRejectedFieldTestAccount(),
		body,
	)

placeholder
	require.NotNil(t, result)
	require.Len(t, upstream.bodies, 2)
	for _, forwardedBody := range upstream.bodies {
		require.False(t, gjson.GetBytes(forwardedBody, "input.0.namespace").Exists())
		require.False(t, gjson.GetBytes(forwardedBody, "input.1.namespace").Exists())
placeholder
	require.Equal(t, int64(2048), gjson.GetBytes(upstream.bodies[0], "max_output_tokens").Int())
	require.False(t, gjson.GetBytes(upstream.bodies[1], "max_output_tokens").Exists())
placeholder

func newOpenAIRejectedFieldTestService(upstream *httpUpstreamRecorder) *OpenAIGatewayService {
	return &OpenAIGatewayService{
		cfg: &config.Config{Security: config.SecurityConfig{
			URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholder,
placeholder
		httpUpstream: upstream,
placeholder
placeholder

func newOpenAIRejectedFieldTestContext(body []byte) *gin.Context {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("User-Agent", "curl/8.0")
	return c
placeholder

func newOpenAIRejectedFieldTestAccount() *Account {
placeholder
		ID:          5107,
		Name:        "responses-compatible",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "https://compat.example",
	placeholder,
		Extra: map[string]any{
			openai_compat.ExtraKeyResponsesMode:      string(openai_compat.ResponsesSupportModeAuto),
			openai_compat.ExtraKeyResponsesSupported: true,
	placeholder,
		Status:      StatusActive,
		Schedulable: true,
placeholder
placeholder

func newOpenAIOAuthNamespaceTestAccount() *Account {
placeholder
		ID:          5108,
		Name:        "openai-oauth-namespace",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token":       "oauth-token",
			"chatgpt_account_id": "chatgpt-account",
	placeholder,
		Status:      StatusActive,
		Schedulable: true,
placeholder
placeholder

func newOpenAIRejectedFieldTestResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(body)),
placeholder
placeholder
