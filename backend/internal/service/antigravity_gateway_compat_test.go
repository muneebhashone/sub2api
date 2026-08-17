package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

type antigravityCompatTokenCache struct {
	token string
placeholder

type antigravityCompatErrorReader struct {
	data []byte
	off  int
	err  error
placeholder

func (r *antigravityCompatErrorReader) Read(p []byte) (int, error) {
	if r.off < len(r.data) {
		n := copy(p, r.data[r.off:])
		r.off += n
		return n, nil
placeholder
	return 0, r.err
placeholder

func (r *antigravityCompatErrorReader) Close() error { return nil placeholder

func (c *antigravityCompatTokenCache) GetAccessToken(context.Context, string) (string, error) {
	return c.token, nil
placeholder

func (c *antigravityCompatTokenCache) SetAccessToken(context.Context, string, string, time.Duration) error {
	return nil
placeholder

func (c *antigravityCompatTokenCache) DeleteAccessToken(context.Context, string) error {
	return nil
placeholder

func (c *antigravityCompatTokenCache) AcquireRefreshLock(context.Context, string, time.Duration) (bool, error) {
	return true, nil
placeholder

func (c *antigravityCompatTokenCache) ReleaseRefreshLock(context.Context, string) error {
	return nil
placeholder

func newAntigravityCompatService(cfg config.GatewayConfig, upstream HTTPUpstream) *AntigravityGatewayService {
	tokenProvider := NewAntigravityTokenProvider(
		nil,
		&antigravityCompatTokenCache{token: "fresh-oauth-token"placeholder,
		nil,
	)
	return NewAntigravityGatewayService(
		nil,
		nil,
		nil,
		tokenProvider,
		nil,
		upstream,
		NewSettingService(&antigravitySettingRepoStub{placeholder, &config.Config{Gateway: cfgplaceholder),
		nil,
	)
placeholder

func newAntigravityCompatAccount(accountType string) *Account {
placeholder
		ID:          3757,
		Name:        "antigravity-compat",
		Platform:    PlatformAntigravity,
		Type:        accountType,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"access_token": "stale-account-token",
			"project_id":   "project-3757",
			"model_mapping": map[string]any{
				"gemini-3.1-pro-high": "gemini-3.1-pro-high",
				"claude-sonnet-4-5":   "claude-sonnet-4-5",
		placeholder,
	placeholder,
placeholder
placeholder

func newAntigravityCompatContext(method, path string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(method, path, bytes.NewReader(body))
	return c, recorder
placeholder

func antigravityCompatSuccessResponse() *http.Response {
	body := `data: {"response":{"responseId":"resp_3757","candidates":[{"content":{"parts":[{"text":"ok"placeholder]placeholder,"finishReason":"STOP"placeholder],"usageMetadata":{"promptTokenCount":8,"candidatesTokenCount":3placeholderplaceholderplaceholder` + "\n\n"
	return &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"text/event-stream"placeholder,
			"X-Request-Id": []string{"request-3757"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(body)),
placeholder
placeholder

func TestAntigravityCompatOAuthUsesNativeTokenAndRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string
		path string
		body []byte
		call func(*AntigravityGatewayService, context.Context, *gin.Context, *Account, []byte) (*ForwardResult, error)
placeholder{
		{
			name: "chat completions",
			path: "/v1/chat/completions",
			body: []byte(`{"model":"gemini-3.1-pro-high","messages":[{"role":"user","content":"Reply exactly: ok"placeholder]placeholder`),
			call: func(svc *AntigravityGatewayService, ctx context.Context, c *gin.Context, account *Account, body []byte) (*ForwardResult, error) {
				return svc.ForwardAsChatCompletions(ctx, c, account, body, nil)
		placeholder,
	placeholder,
		{
			name: "responses",
			path: "/v1/responses",
			body: []byte(`{"model":"gemini-3.1-pro-high","input":"Reply exactly: ok"placeholder`),
			call: func(svc *AntigravityGatewayService, ctx context.Context, c *gin.Context, account *Account, body []byte) (*ForwardResult, error) {
				return svc.ForwardAsResponses(ctx, c, account, body, nil)
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var authorization string
			var upstreamPath string
			var upstreamAlt string
			upstream := &queuedHTTPUpstreamStub{
				responses: []*http.Response{antigravityCompatSuccessResponse()placeholder,
				onCall: func(req *http.Request, _ *queuedHTTPUpstreamStub) {
					authorization = req.Header.Get("Authorization")
					upstreamPath = req.URL.Path
					upstreamAlt = req.URL.Query().Get("alt")
			placeholder,
		placeholder
			svc := newAntigravityCompatService(
				config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
				upstream,
			)
			c, recorder := newAntigravityCompatContext(http.MethodPost, tt.path, tt.body)

			result, err := tt.call(svc, context.Background(), c, newAntigravityCompatAccount(AccountTypeOAuth), tt.body)

		placeholder
			require.NotNil(t, result)
			require.Equal(t, "Bearer fresh-oauth-token", authorization)
			require.Equal(t, "/v1internal:streamGenerateContent", upstreamPath)
			require.Equal(t, "sse", upstreamAlt)
			require.Equal(t, "request-3757", result.RequestID)
			require.Equal(t, 8, result.Usage.InputTokens)
			require.Equal(t, 3, result.Usage.OutputTokens)
			require.Equal(t, http.StatusOK, recorder.Code)
			require.Contains(t, recorder.Body.String(), "ok")
			if tt.name == "chat completions" {
				require.Equal(t, "stop", gjson.Get(recorder.Body.String(), "choices.0.finish_reason").String())
				require.Equal(t, int64(8), gjson.Get(recorder.Body.String(), "usage.prompt_tokens").Int())
				require.Equal(t, int64(3), gjson.Get(recorder.Body.String(), "usage.completion_tokens").Int())
		placeholder
	placeholder)
placeholder
placeholder

func TestAntigravityCompatRejectsUnsupportedAccountType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		path        string
		accountType string
		call        func(*AntigravityGatewayService, context.Context, *gin.Context, *Account, []byte) (*ForwardResult, error)
placeholder{
		{
			name:        "chat completions upstream",
			path:        "/v1/chat/completions",
			accountType: AccountTypeUpstream,
			call: func(svc *AntigravityGatewayService, ctx context.Context, c *gin.Context, account *Account, body []byte) (*ForwardResult, error) {
				return svc.ForwardAsChatCompletions(ctx, c, account, body, nil)
		placeholder,
	placeholder,
		{
			name:        "responses setup token",
			path:        "/v1/responses",
			accountType: AccountTypeSetupToken,
			call: func(svc *AntigravityGatewayService, ctx context.Context, c *gin.Context, account *Account, body []byte) (*ForwardResult, error) {
				return svc.ForwardAsResponses(ctx, c, account, body, nil)
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := []byte(`{"model":"gemini-3.1-pro-high"placeholder`)
			c, recorder := newAntigravityCompatContext(http.MethodPost, tt.path, body)

			result, err := tt.call(&AntigravityGatewayService{placeholder, context.Background(), c, newAntigravityCompatAccount(tt.accountType), body)

		placeholder
			require.Nil(t, result)
			require.Equal(t, http.StatusBadRequest, recorder.Code)
			require.Contains(t, recorder.Body.String(), "native OAuth account required for antigravity compatibility mode")
	placeholder)
placeholder
placeholder

func TestBuildAntigravityCompatGeminiBody_ConfiguresMixedToolInvocations(t *testing.T) {
	svc := &AntigravityGatewayService{placeholder
	tests := []struct {
		name      string
		tools     string
		wantField bool
placeholder{
		{
			name:      "mixed server and client tools",
			tools:     `[{"name":"get_weather","input_schema":{"type":"object"placeholderplaceholder,{"type":"web_search_20250305","name":"web_search"placeholder]`,
			wantField: true,
	placeholder,
		{
			name:  "client tools only",
			tools: `[{"name":"get_weather","input_schema":{"type":"object"placeholderplaceholder]`,
	placeholder,
		{
			name:  "server tools only",
			tools: `[{"type":"web_search_20250305","name":"web_search"placeholder]`,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claudeBody := []byte(`{"messages":[{"role":"user","content":"hello"placeholder],"tools":` + tt.tools + `placeholder`)
			claudeBody = bytes.ReplaceAll(claudeBody, []byte{92placeholder, nil)
			body, err := svc.buildAntigravityCompatGeminiBody(context.Background(), claudeBody, nil, "project-1", "gemini-2.5-flash")
		placeholder

			var wrapped map[string]any
			require.NoError(t, json.Unmarshal(body, &wrapped))
			request := wrapped["request"].(map[string]any)
			toolConfig, exists := request["toolConfig"].(map[string]any)
			if !tt.wantField {
				require.False(t, exists)
				return
		placeholder
			require.True(t, exists)
			require.Equal(t, true, toolConfig["includeServerSideToolInvocations"])
			require.NotContains(t, toolConfig, "include_server_side_tool_invocations")
	placeholder)
placeholder
placeholder

func TestAntigravityCompatPreservesChatTokenLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name string
		body string
		want int64
placeholder{
		{
			name: "legacy max_tokens below bridge floor",
			body: `{"model":"gemini-3.1-pro-high","messages":[{"role":"user","content":"ok"placeholder],"max_tokens":8placeholder`,
			want: 8,
	placeholder,
		{
			name: "max_completion_tokens takes precedence",
			body: `{"model":"gemini-3.1-pro-high","messages":[{"role":"user","content":"ok"placeholder],"max_tokens":8,"max_completion_tokens":13placeholder`,
			want: 13,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upstream := &queuedHTTPUpstreamStub{responses: []*http.Response{antigravityCompatSuccessResponse()placeholderplaceholder
			svc := newAntigravityCompatService(config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder, upstream)
			body := []byte(tt.body)
			c, _ := newAntigravityCompatContext(http.MethodPost, "/v1/chat/completions", body)

			result, err := svc.ForwardAsChatCompletions(
				context.Background(),
				c,
				newAntigravityCompatAccount(AccountTypeOAuth),
				body,
				nil,
			)

		placeholder
			require.NotNil(t, result)
			require.Len(t, upstream.requestBodies, 1)
			require.Equal(t, tt.want, gjson.GetBytes(upstream.requestBodies[0], "request.generationConfig.maxOutputTokens").Int())
	placeholder)
placeholder
placeholder

func TestAntigravityCompatRoutesByMappedModelFamily(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		model         string
		wantSessionID bool
placeholder{
		{model: "gemini-3.1-pro-high", wantSessionID: falseplaceholder,
		{model: "claude-sonnet-4-5", wantSessionID: trueplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			upstream := &queuedHTTPUpstreamStub{responses: []*http.Response{antigravityCompatSuccessResponse()placeholderplaceholder
			svc := newAntigravityCompatService(config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder, upstream)
			body := []byte(`{"model":"` + tt.model + `","messages":[{"role":"user","content":"ok"placeholder],"max_tokens":8placeholder`)
			c, _ := newAntigravityCompatContext(http.MethodPost, "/v1/chat/completions", body)

			result, err := svc.ForwardAsChatCompletions(
				context.Background(),
				c,
				newAntigravityCompatAccount(AccountTypeOAuth),
				body,
				nil,
			)

		placeholder
			require.NotNil(t, result)
			require.Len(t, upstream.requestBodies, 1)
			require.Equal(t, tt.model, gjson.GetBytes(upstream.requestBodies[0], "model").String())
			require.Equal(t, tt.wantSessionID, gjson.GetBytes(upstream.requestBodies[0], "request.sessionId").Exists())
	placeholder)
placeholder
placeholder

func TestAntigravityCompatUnauthorizedIsCredentialFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstream := &queuedHTTPUpstreamStub{responses: []*http.Response{{
		StatusCode: http.StatusUnauthorized,
		Header:     http.Header{"X-Request-Id": []string{"auth-3757"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"Invalid bearer token"placeholderplaceholder`)),
placeholderplaceholderplaceholder
	svc := newAntigravityCompatService(config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder, upstream)
	body := []byte(`{"model":"gemini-3.1-pro-high","messages":[{"role":"user","content":"ok"placeholder]placeholder`)
	c, recorder := newAntigravityCompatContext(http.MethodPost, "/v1/chat/completions", body)

	result, err := svc.ForwardAsChatCompletions(
		context.Background(),
		c,
		newAntigravityCompatAccount(AccountTypeOAuth),
		body,
		nil,
	)

	require.Nil(t, result)
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, GatewayFailureStageAccountAuth, failoverErr.Stage)
	require.Equal(t, GatewayFailureScopeAccount, failoverErr.Scope)
	require.Equal(t, AntigravityCredentialRejectedReason, failoverErr.Reason)
	require.Equal(t, NextAccountRetry, failoverErr.NextAccountAction)
	require.Equal(t, http.StatusBadGateway, failoverErr.ClientStatusCode)
	require.Equal(t, AntigravityCredentialRejectedClientMessage, failoverErr.ClientMessage)
	require.Equal(t, "auth-3757", failoverErr.ResponseHeaders.Get("X-Request-Id"))
	require.Empty(t, recorder.Body.String())
placeholder

func TestAntigravityCompatEmptyStreamTriggersFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string
		run  func(*AntigravityGatewayService, *gin.Context, *http.Response) (*antigravityStreamResult, error)
placeholder{
		{
			name: "chat completions",
			run: func(svc *AntigravityGatewayService, c *gin.Context, resp *http.Response) (*antigravityStreamResult, error) {
				return svc.handleChatCompletionsStreamingFromAntigravity(c, resp, time.Now(), "gemini-3.1-pro-high", true)
		placeholder,
	placeholder,
		{
			name: "responses",
			run: func(svc *AntigravityGatewayService, c *gin.Context, resp *http.Response) (*antigravityStreamResult, error) {
				return svc.handleResponsesStreamingFromAntigravity(c, resp, time.Now(), "gemini-3.1-pro-high")
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newAntigravityCompatService(config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder, nil)
			c, recorder := newAntigravityCompatContext(http.MethodPost, "/", nil)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
				Body:       io.NopCloser(strings.NewReader("data: malformed\n\ndata: [DONE]\n\n")),
		placeholder

			result, err := tt.run(svc, c, resp)

			require.Nil(t, result)
			var failoverErr *UpstreamFailoverError
			require.ErrorAs(t, err, &failoverErr)
			require.True(t, failoverErr.RetryableOnSameAccount)
			require.Empty(t, recorder.Body.String())
			require.Empty(t, recorder.Header().Get("Content-Type"))
	placeholder)
placeholder
placeholder

func TestAntigravityCompatUsageOnlyStreamTriggersFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string
		run  func(*AntigravityGatewayService, *gin.Context, *http.Response) (*antigravityStreamResult, error)
placeholder{
		{
			name: "chat completions",
			run: func(svc *AntigravityGatewayService, c *gin.Context, resp *http.Response) (*antigravityStreamResult, error) {
				return svc.handleChatCompletionsStreamingFromAntigravity(c, resp, time.Now(), "gemini-3.1-pro-high", true)
		placeholder,
	placeholder,
		{
			name: "responses",
			run: func(svc *AntigravityGatewayService, c *gin.Context, resp *http.Response) (*antigravityStreamResult, error) {
				return svc.handleResponsesStreamingFromAntigravity(c, resp, time.Now(), "gemini-3.1-pro-high")
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newAntigravityCompatService(config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder, nil)
			c, recorder := newAntigravityCompatContext(http.MethodPost, "/", nil)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
				Body: io.NopCloser(strings.NewReader(
					`data: {"response":{"responseId":"resp_3757","usageMetadata":{"promptTokenCount":8placeholderplaceholderplaceholder` + "\n\n",
				)),
		placeholder

			result, err := tt.run(svc, c, resp)

			require.Nil(t, result)
			var failoverErr *UpstreamFailoverError
			require.ErrorAs(t, err, &failoverErr)
			require.True(t, failoverErr.RetryableOnSameAccount)
			require.Empty(t, recorder.Body.String())
			require.Empty(t, recorder.Header().Get("Content-Type"))
	placeholder)
placeholder
placeholder

func TestAntigravityCompatUsageOnlyNonStreamingTriggersFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string
		path string
		body []byte
		call func(*AntigravityGatewayService, context.Context, *gin.Context, *Account, []byte) (*ForwardResult, error)
placeholder{
		{
			name: "chat completions",
			path: "/v1/chat/completions",
			body: []byte(`{"model":"gemini-3.1-pro-high","messages":[{"role":"user","content":"ok"placeholder]placeholder`),
			call: func(svc *AntigravityGatewayService, ctx context.Context, c *gin.Context, account *Account, body []byte) (*ForwardResult, error) {
				return svc.ForwardAsChatCompletions(ctx, c, account, body, nil)
		placeholder,
	placeholder,
		{
			name: "responses",
			path: "/v1/responses",
			body: []byte(`{"model":"gemini-3.1-pro-high","input":"ok"placeholder`),
			call: func(svc *AntigravityGatewayService, ctx context.Context, c *gin.Context, account *Account, body []byte) (*ForwardResult, error) {
				return svc.ForwardAsResponses(ctx, c, account, body, nil)
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upstream := &queuedHTTPUpstreamStub{responses: []*http.Response{{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
				Body: io.NopCloser(strings.NewReader(
					`data: {"response":{"responseId":"resp_3757","usageMetadata":{"promptTokenCount":8placeholderplaceholderplaceholder` + "\n\n",
				)),
		placeholderplaceholderplaceholder
			svc := newAntigravityCompatService(config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder, upstream)
			c, recorder := newAntigravityCompatContext(http.MethodPost, tt.path, tt.body)

			result, err := tt.call(
				svc,
				context.Background(),
				c,
				newAntigravityCompatAccount(AccountTypeOAuth),
				tt.body,
			)

			require.Nil(t, result)
			var failoverErr *UpstreamFailoverError
			require.ErrorAs(t, err, &failoverErr)
			require.True(t, failoverErr.RetryableOnSameAccount)
			require.Empty(t, recorder.Body.String())
			require.Empty(t, recorder.Header().Get("Content-Type"))
	placeholder)
placeholder
placeholder

func TestAntigravityCompatChatStreamMapsToolCallAndUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityCompatService(config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder, nil)
	c, recorder := newAntigravityCompatContext(http.MethodPost, "/v1/chat/completions", nil)
	body := `data: {"response":{"responseId":"resp_3757","candidates":[{"content":{"parts":[{"functionCall":{"id":"call_3757","name":"get_weather","args":{"city":"Tokyo"placeholderplaceholderplaceholder]placeholder,"finishReason":"STOP"placeholder],"usageMetadata":{"promptTokenCount":8,"candidatesTokenCount":3placeholderplaceholderplaceholder` + "\n\n"
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(body)),
placeholder

	result, err := svc.handleChatCompletionsStreamingFromAntigravity(
		c,
		resp,
		time.Now(),
		"gemini-3.1-pro-high",
		true,
	)

placeholder
	require.NotNil(t, result)
	require.Equal(t, 8, result.usage.InputTokens)
	require.Equal(t, 3, result.usage.OutputTokens)
	require.Contains(t, recorder.Body.String(), `"tool_calls"`)
	require.Contains(t, recorder.Body.String(), `"get_weather"`)
	require.Contains(t, recorder.Body.String(), `"finish_reason":"tool_calls"`)
	require.Contains(t, recorder.Body.String(), `"prompt_tokens":8`)
	require.Equal(t, 1, strings.Count(recorder.Body.String(), "data: [DONE]"))
placeholder

func TestAntigravityCompatFirstEventTimeoutTriggersFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityCompatService(
		config.GatewayConfig{MaxLineSize: defaultMaxLineSize, StreamDataIntervalTimeout: 1placeholder,
		nil,
	)
	c, recorder := newAntigravityCompatContext(http.MethodPost, "/v1/chat/completions", nil)
	reader, writer := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Header: http.Header{placeholder, Body: readerplaceholder
	type outcome struct {
		result *antigravityStreamResult
		err    error
placeholder
	done := make(chan outcome, 1)

	go func() {
		result, err := svc.handleChatCompletionsStreamingFromAntigravity(c, resp, time.Now(), "gemini-3.1-pro-high", false)
		done <- outcome{result: result, err: errplaceholder
placeholder()

	select {
	case got := <-done:
		require.Nil(t, got.result)
		var failoverErr *UpstreamFailoverError
		require.ErrorAs(t, got.err, &failoverErr)
		require.True(t, failoverErr.RetryableOnSameAccount)
		require.Empty(t, recorder.Header().Get("Content-Type"))
	case <-time.After(2 * time.Second):
		_ = writer.Close()
		_ = reader.Close()
		t.Fatal("compat stream ignored StreamDataIntervalTimeout")
placeholder
	_ = writer.Close()
	_ = reader.Close()
placeholder

func TestAntigravityCompatClientDisconnectDrainsUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityCompatService(config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder, nil)
	c, _ := newAntigravityCompatContext(http.MethodPost, "/v1/chat/completions", nil)
	c.Writer = &antigravityFailingWriter{ResponseWriter: c.Writer, failAfter: 0placeholder
	body := strings.Join([]string{
		`data: {"response":{"responseId":"resp_3757","candidates":[{"content":{"parts":[{"text":"partial"placeholder]placeholderplaceholder],"usageMetadata":{"promptTokenCount":8,"candidatesTokenCount":1placeholderplaceholderplaceholder`,
		"",
		`data: {"response":{"responseId":"resp_3757","candidates":[{"finishReason":"STOP"placeholder],"usageMetadata":{"promptTokenCount":8,"candidatesTokenCount":15placeholderplaceholderplaceholder`,
		"",
placeholder, "\n")
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(body)),
placeholder

	result, err := svc.handleChatCompletionsStreamingFromAntigravity(c, resp, time.Now(), "gemini-3.1-pro-high", false)

placeholder
	require.NotNil(t, result)
	require.True(t, result.clientDisconnect)
	require.Equal(t, 15, result.usage.OutputTokens)
placeholder

func TestAntigravityCompatStreamErrorCommitsSingleTerminalFrame(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityCompatService(config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder, nil)
	c, recorder := newAntigravityCompatContext(http.MethodPost, "/v1/responses", nil)
	body := []byte(`data: {"response":{"responseId":"resp_3757","candidates":[{"content":{"parts":[{"text":"partial"placeholder]placeholderplaceholder],"usageMetadata":{"promptTokenCount":8,"candidatesTokenCount":1placeholderplaceholderplaceholder` + "\n\n")
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body: &antigravityCompatErrorReader{
			data: body,
			err:  io.ErrUnexpectedEOF,
	placeholder,
placeholder

	result, err := svc.handleResponsesStreamingFromAntigravity(c, resp, time.Now(), "gemini-3.1-pro-high")

placeholder
	require.Nil(t, result)
	require.True(t, IsResponseCommitted(c))
	require.Equal(t, 1, strings.Count(recorder.Body.String(), "event: error"))
placeholder

func TestAntigravityCompatKeepaliveAfterFirstEvent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityCompatService(
		config.GatewayConfig{MaxLineSize: defaultMaxLineSize, StreamKeepaliveInterval: 1placeholder,
		nil,
	)
	c, recorder := newAntigravityCompatContext(http.MethodPost, "/v1/responses", nil)
	reader, writer := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Header: http.Header{placeholder, Body: readerplaceholder
	done := make(chan error, 1)

	go func() {
		_, err := svc.handleResponsesStreamingFromAntigravity(c, resp, time.Now(), "gemini-3.1-pro-high")
		done <- err
placeholder()
	_, err := io.WriteString(
		writer,
		`data: {"response":{"responseId":"resp_3757","candidates":[{"content":{"parts":[{"text":"partial"placeholder]placeholderplaceholder],"usageMetadata":{"promptTokenCount":8,"candidatesTokenCount":1placeholderplaceholderplaceholder`+"\n\n",
	)
placeholder
	time.Sleep(1200 * time.Millisecond)
	require.NoError(t, writer.Close())
	require.NoError(t, <-done)
	require.Contains(t, recorder.Body.String(), ": ping\n\n")
	require.Contains(t, recorder.Header().Get("Content-Type"), "text/event-stream")
	require.NoError(t, reader.Close())
placeholder
