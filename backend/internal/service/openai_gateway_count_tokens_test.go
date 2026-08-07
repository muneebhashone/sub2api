package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/apicompat"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

type countTokensRuntimeStateRepo struct {
	AccountRepository
	tempUnschedCalls int
	setErrorCalls    int
placeholder

func (r *countTokensRuntimeStateRepo) SetTempUnschedulable(_ context.Context, _ int64, _ time.Time, _ string) error {
	r.tempUnschedCalls++
	return nil
placeholder

func (r *countTokensRuntimeStateRepo) SetError(_ context.Context, _ int64, _ string) error {
	r.setErrorCalls++
	return nil
placeholder

func TestOpenAIGatewayService_ForwardCountTokensAsAnthropic_APIKeyUsesResponsesInputTokens(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	body := []byte(`{"model":"claude-sonnet-4-5","system":"You are helpful.","messages":[{"role":"user","content":"hello"placeholder],"tools":[{"name":"lookup","input_schema":{"type":"object"placeholderplaceholder]placeholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages/count_tokens", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"object":"response.input_tokens","input_tokens":42placeholder`)),
placeholderplaceholder

	svc := &OpenAIGatewayService{
		cfg: &config.Config{Security: config.SecurityConfig{URLAllowlist: config.URLAllowlistConfig{
			Enabled:           false,
			AllowInsecureHTTP: true,
	placeholderplaceholderplaceholder,
		httpUpstream: upstream,
placeholder
	account := &Account{
		ID:          101,
		Name:        "openai-apikey",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "http://upstream.example",
	placeholder,
		Status:      StatusActive,
		Schedulable: true,
placeholder

	err := svc.ForwardCountTokensAsAnthropic(context.Background(), c, account, body, "gpt-5.3-codex")
placeholder
	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `{"input_tokens":42placeholder`, rec.Body.String())
	require.NotNil(t, upstream.lastReq)
	require.Equal(t, "http://upstream.example/v1/responses/input_tokens", upstream.lastReq.URL.String())
	require.Equal(t, "Bearer sk-test", upstream.lastReq.Header.Get("authorization"))
	require.Equal(t, "gpt-5.3-codex", gjson.GetBytes(upstream.lastBody, "model").String())
	require.True(t, gjson.GetBytes(upstream.lastBody, "input").Exists())
	require.False(t, gjson.GetBytes(upstream.lastBody, "messages").Exists())
placeholder

func TestOpenAIGatewayService_ForwardCountTokensAsAnthropic_OAuthFallsBackWhenPlatformEndpointUnsupported(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"model":"claude-opus-4-1","messages":[{"role":"user","content":"hello"placeholder]placeholder`)
	account := &Account{
		ID:          202,
		Name:        "openai-oauth",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token":  "oauth-token",
			"refresh_token": "oauth-refresh-token",
	placeholder,
		Status:      StatusActive,
		Schedulable: true,
placeholder

	prepared, err := prepareOpenAIInputTokensCountRequest(body, account, "gpt-5.4")
placeholder
	expectedEstimate, err := estimateOpenAIInputTokens(prepared.Request)
placeholder

	cases := []struct {
		name       string
		statusCode int
		body       string
placeholder{
		{
			name:       "401_missing_responses_write_scope",
			statusCode: http.StatusUnauthorized,
			body:       `{"error":{"type":"invalid_request_error","code":"missing_scope","message":"You have insufficient permissions for this operation. Missing scopes: api.responses.write."placeholderplaceholder`,
	placeholder,
		{
			name:       "403_missing_responses_write_scope",
			statusCode: http.StatusForbidden,
			body:       `{"error":{"type":"invalid_request_error","code":"missing_scope","message":"Missing scopes: api.responses.write"placeholderplaceholder`,
	placeholder,
		{
			name:       "404_input_tokens_unsupported",
			statusCode: http.StatusNotFound,
			body:       `{"error":{"type":"invalid_request_error","message":"The /v1/responses/input_tokens endpoint was not found"placeholderplaceholder`,
	placeholder,
placeholder

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages/count_tokens", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Request.Header.Set("User-Agent", "Claude-Code/1.0")

			upstream := &httpUpstreamRecorder{resp: &http.Response{
				StatusCode: tt.statusCode,
				Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
				Body:       io.NopCloser(strings.NewReader(tt.body)),
		placeholderplaceholder
			repo := &countTokensRuntimeStateRepo{placeholder
			svc := &OpenAIGatewayService{
				cfg:              &config.Config{placeholder,
				httpUpstream:     upstream,
				rateLimitService: &RateLimitService{accountRepo: repo, cfg: &config.Config{placeholderplaceholder,
		placeholder

			err := svc.ForwardCountTokensAsAnthropic(context.Background(), c, account, body, "gpt-5.4")
		placeholder
			require.Equal(t, http.StatusOK, rec.Code)
			require.JSONEq(t, `{"input_tokens":`+strconv.Itoa(expectedEstimate)+`placeholder`, rec.Body.String())
			require.NotNil(t, upstream.lastReq)
			require.Equal(t, "https://api.openai.com/v1/responses/input_tokens", upstream.lastReq.URL.String())
			require.Equal(t, "Bearer oauth-token", upstream.lastReq.Header.Get("authorization"))
			require.Empty(t, upstream.lastReq.Header.Get("Chatgpt-Account-Id"))
			require.Zero(t, repo.tempUnschedCalls, "OAuth input_tokens unsupported errors must not temp-unschedule the account")
			require.Zero(t, repo.setErrorCalls, "OAuth input_tokens unsupported errors must not mark the account error")
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayService_OpenAIOAuthInputTokensFallbackUsesMinimumWhenEstimateFails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	prepared := &openAIInputTokensCountPrepared{
		Request: openAIInputTokensCountRequest{
			Model: "gpt-5",
			Input: json.RawMessage(`[`),
	placeholder,
		UpstreamModel: "gpt-5",
placeholder

	writeOpenAIOAuthInputTokensFallback(c, &Account{ID: 303placeholder, prepared, http.StatusUnauthorized)

	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `{"input_tokens":1placeholder`, rec.Body.String())
placeholder

func TestEstimateOpenAIInputTokens_RequestSamples(t *testing.T) {
	cases := []struct {
		name string
		req  openAIInputTokensCountRequest
		want int
placeholder{
		{
			name: "simple text input",
			req: openAIInputTokensCountRequest{
				Model: "gpt-5",
				Input: json.RawMessage(`[{"role":"user","content":"hello world"placeholder]`),
		placeholder,
			want: 6,
	placeholder,
		{
			name: "instructions plus tool schema",
			req: openAIInputTokensCountRequest{
				Model:        "gpt-5",
				Instructions: "You are helpful.",
				Input:        json.RawMessage(`[{"role":"user","content":"lookup weather in shanghai"placeholder]`),
				Tools: []apicompat.ResponsesTool{
					{
						Type:        "function",
						Name:        "lookup_weather",
						Description: "Look up current weather",
						Parameters:  json.RawMessage(`{"type":"object","properties":{"city":{"type":"string"placeholderplaceholder,"required":["city"]placeholder`),
				placeholder,
			placeholder,
		placeholder,
			want: 50,
	placeholder,
		{
			name: "input parts and tool output",
			req: openAIInputTokensCountRequest{
				Model: "gpt-4.1",
				Input: json.RawMessage(`[
					{"role":"user","content":[{"type":"input_text","text":"first line"placeholder,{"type":"input_text","text":"second line"placeholder]placeholder,
					{"type":"function_call_output","call_id":"call_123","output":"{\"ok\":trueplaceholder"placeholder
				]`),
		placeholder,
			want: 24,
	placeholder,
placeholder

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := estimateOpenAIInputTokens(tt.req)
		placeholder
			require.Equal(t, tt.want, got)
	placeholder)
placeholder
placeholder

func TestEstimateGrokCountTokens_AnthropicRequests(t *testing.T) {
	cases := []struct {
		name string
		body string
placeholder{
		{
			name: "simple message",
			body: `{"model":"grok-4","messages":[{"role":"user","content":"hello world"placeholder]placeholder`,
	placeholder,
		{
			name: "system blocks and tools",
			body: `{
				"model":"grok-4",
				"system":[{"type":"text","text":"You are helpful."placeholder],
				"messages":[{"role":"user","content":[{"type":"text","text":"look up the weather"placeholder]placeholder],
				"tools":[{"name":"lookup_weather","description":"Look up weather","input_schema":{"type":"object","properties":{"city":{"type":"string"placeholderplaceholderplaceholderplaceholder],
				"tool_choice":{"type":"auto"placeholder
		placeholder`,
	placeholder,
		{
			name: "empty conversation uses positive minimum",
			body: `{"model":"grok-4","messages":[]placeholder`,
	placeholder,
placeholder

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EstimateGrokCountTokens([]byte(tt.body))
		placeholder
			require.Positive(t, got)
	placeholder)
placeholder
placeholder

func TestEstimateGrokCountTokens_RejectsInvalidRequests(t *testing.T) {
	for _, body := range []string{
		`{`,
		`{"messages":[{"role":"user","content":"hello"placeholder]placeholder`,
		`{"model":"grok-4","messages":[{"role":"user","content":{"unexpected":trueplaceholderplaceholder]placeholder`,
placeholder {
		_, err := EstimateGrokCountTokens([]byte(body))
		require.Error(t, err, "body=%s", body)
placeholder
placeholder

func TestOpenAIInputTokensEncodingForModel(t *testing.T) {
	cases := []struct {
		model string
		want  string
placeholder{
		{model: "gpt-5", want: "o200k_base"placeholder,
		{model: "gpt-5.3-codex", want: "o200k_base"placeholder,
		{model: "gpt-4o-mini", want: "o200k_base"placeholder,
		{model: "gpt-4.1", want: "o200k_base"placeholder,
		{model: "gpt-4-turbo", want: "cl100k_base"placeholder,
		{model: "gpt-3.5-turbo", want: "cl100k_base"placeholder,
placeholder

	for _, tt := range cases {
		t.Run(tt.model, func(t *testing.T) {
			require.Equal(t, tt.want, string(openAIInputTokensEncodingForModel(tt.model)))
	placeholder)
placeholder
placeholder

func TestEstimateOpenAIInputTokens_CompareWithOpenAIAPI(t *testing.T) {
	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if apiKey == "" {
		t.Skip("OPENAI_API_KEY not set")
placeholder
	// Invalid/expired keys in local env must not fail the unit suite.
	if strings.HasPrefix(apiKey, "sk-") && len(apiKey) < 20 {
		t.Skip("OPENAI_API_KEY looks incomplete")
placeholder

	client := &http.Client{Timeout: 30 * time.Secondplaceholder
	cases := []struct {
		name               string
		anthropicBody      []byte
		defaultOpenAIModel string
placeholder{
		{
			name:               "simple user text",
			defaultOpenAIModel: "gpt-5",
			anthropicBody:      []byte(`{"model":"claude-sonnet-4-5","messages":[{"role":"user","content":"hello world from sub2api"placeholder]placeholder`),
	placeholder,
		{
			name:               "system plus tool",
			defaultOpenAIModel: "gpt-5",
			anthropicBody:      []byte(`{"model":"claude-sonnet-4-5","system":"You are helpful.","messages":[{"role":"user","content":"find weather in shanghai"placeholder],"tools":[{"name":"lookup_weather","description":"Look up current weather","input_schema":{"type":"object","properties":{"city":{"type":"string"placeholderplaceholder,"required":["city"]placeholderplaceholder]placeholder`),
	placeholder,
		{
			name:               "multi turn text",
			defaultOpenAIModel: "gpt-4.1",
			anthropicBody:      []byte(`{"model":"claude-opus-4-1","messages":[{"role":"user","content":"summarize this repo"placeholder,{"role":"assistant","content":"which repo?"placeholder,{"role":"user","content":"sub2api"placeholder]placeholder`),
	placeholder,
placeholder

	account := &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			prepared, err := prepareOpenAIInputTokensCountRequest(tc.anthropicBody, account, tc.defaultOpenAIModel)
		placeholder

			estimated, err := estimateOpenAIInputTokens(prepared.Request)
		placeholder

			actual, err := callOpenAIInputTokensAPIForTest(client, apiKey, prepared.Request)
			if err != nil {
				// Live-API comparison only; invalid/expired local keys should skip, not fail CI.
				if strings.Contains(err.Error(), "status=401") || strings.Contains(err.Error(), "invalid_api_key") {
					t.Skipf("OPENAI_API_KEY rejected by OpenAI: %v", err)
			placeholder
			placeholder
		placeholder

			diff := estimated - actual
			if diff < 0 {
				diff = -diff
		placeholder
			t.Logf("model=%s estimated=%d actual=%d diff=%d", prepared.Request.Model, estimated, actual, diff)
			require.LessOrEqual(t, diff, maxLocalInt(24, actual/4))
	placeholder)
placeholder
placeholder

func callOpenAIInputTokensAPIForTest(client *http.Client, apiKey string, reqBody openAIInputTokensCountRequest) (int, error) {
	body, err := marshalOpenAIUpstreamJSON(reqBody)
	if err != nil {
		return 0, err
placeholder
	req, err := http.NewRequest(http.MethodPost, openaiPlatformAPIInputTokensURL, bytes.NewReader(body))
	if err != nil {
		return 0, err
placeholder
	req.Header.Set("authorization", "Bearer "+apiKey)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
placeholder
	if resp.StatusCode >= 400 {
		return 0, fmt.Errorf("openai input_tokens api error: status=%d body=%s", resp.StatusCode, string(respBody))
placeholder

	value := gjson.GetBytes(respBody, "input_tokens")
	if !value.Exists() {
		return 0, fmt.Errorf("openai input_tokens api missing input_tokens: %s", string(respBody))
placeholder
	return int(value.Int()), nil
placeholder

func maxLocalInt(a, b int) int {
	if a > b {
		return a
placeholder
	return b
placeholder
