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
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestIsOpenAIOAuthLike(t *testing.T) {
	tests := []struct {
		name    string
		account *Account
		want    bool
		codex   bool
placeholder{
		{name: "openai_oauth", account: &Account{Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder, want: true, codex: trueplaceholder,
		{name: "openai_setup_token", account: &Account{Platform: PlatformOpenAI, Type: AccountTypeSetupTokenplaceholder, want: true, codex: trueplaceholder,
		{name: "openai_api_key", account: &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder, want: false, codex: falseplaceholder,
		{name: "anthropic_setup_token", account: &Account{Platform: PlatformAnthropic, Type: AccountTypeSetupTokenplaceholder, want: false, codex: falseplaceholder,
		{name: "grok_setup_token", account: &Account{Platform: PlatformGrok, Type: AccountTypeSetupTokenplaceholder, want: false, codex: falseplaceholder,
		{name: "nil", account: nil, want: false, codex: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.account.IsOpenAIOAuthLike())
			require.Equal(t, tt.codex, tt.account.UsesOpenAICodexProtocol())
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayServiceGetAccessTokenSetupToken(t *testing.T) {
	svc := &OpenAIGatewayService{openAITokenProvider: &OpenAITokenProvider{placeholderplaceholder
	account := &Account{
		Platform:    PlatformOpenAI,
		Type:        AccountTypeSetupToken,
placeholder"access_token": "setup-token-value"placeholder,
placeholder

	token, tokenType, err := svc.GetAccessToken(context.Background(), account)
placeholder
	require.Equal(t, "setup-token-value", token)
	require.Equal(t, "oauth", tokenType)

	delete(account.Credentials, "access_token")
	_, _, err = svc.GetAccessToken(context.Background(), account)
	require.EqualError(t, err, "access_token not found in credentials")

	for _, platform := range []string{PlatformAnthropic, PlatformGrokplaceholder {
		foreign := &Account{
			Platform:    platform,
			Type:        AccountTypeSetupToken,
	placeholder"access_token": "foreign-token"placeholder,
	placeholder
		_, _, err = svc.GetAccessToken(context.Background(), foreign)
		require.EqualError(t, err, "unsupported account type: setup-token")
placeholder
placeholder

func TestOpenAISetupTokenImagesUsesOAuthResponsesPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", nil)

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"rate limited"placeholderplaceholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{cfg: &config.Config{placeholder, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          73,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeSetupToken,
placeholder"access_token": "setup-token"placeholder,
placeholder
	parsed := &OpenAIImagesRequest{
		Endpoint:       openAIImagesGenerationsEndpoint,
		Model:          "gpt-image-2",
		Prompt:         "draw a square",
		N:              1,
		ResponseFormat: "b64_json",
placeholder

	result, err := svc.ForwardImages(context.Background(), c, account, nil, parsed, "")

	require.Nil(t, result)
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusTooManyRequests, failoverErr.StatusCode)
	require.True(t, failoverErr.RetryableOnSameAccount)
	require.False(t, failoverErr.SameAccountRetryDeadline.IsZero())
	require.Contains(t, upstream.lastReq.URL.String(), "/backend-api/codex/responses")
placeholder

func TestOpenAISetupTokenWSCompatibility(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/v1/responses", strings.NewReader(`{placeholder`))
	c.Request.Header.Set("session_id", "session-one")
	c.Set("api_key", &APIKey{ID: 17placeholder)

	account := &Account{
		Platform: PlatformOpenAI,
		Type:     AccountTypeSetupToken,
placeholder
			"access_token":       "setup-token-value",
			"chatgpt_account_id": "chatgpt-setup",
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: &config.Config{placeholderplaceholder

	wsURL, err := svc.buildOpenAIResponsesWSURL(account)
placeholder
	require.Equal(t, "wss://chatgpt.com/backend-api/codex/responses", wsURL)
	foreignURL, err := svc.buildOpenAIResponsesWSURL(&Account{Platform: PlatformGrok, Type: AccountTypeSetupTokenplaceholder)
placeholder
	require.Equal(t, "wss://api.openai.com/v1/responses", foreignURL)

	headers, session, err := svc.buildOpenAIWSHeaders(
		context.Background(), c, account, "setup-token-value",
		OpenAIWSProtocolDecision{Transport: OpenAIUpstreamTransportResponsesWebsocketV2placeholder,
		true, "", "", "", "gpt-5.1-codex", "",
	)
placeholder
	require.Equal(t, "Bearer setup-token-value", headers.Get("authorization"))
	require.Equal(t, "chatgpt-setup", headers.Get("chatgpt-account-id"))
	require.NotEmpty(t, headers.Get("originator"))
	require.Equal(t, "session-one", session.SessionID)
	require.NotEqual(t, session.SessionID, headers.Get("session_id"))

	payload := svc.buildOpenAIWSCreatePayload(map[string]any{"store": trueplaceholder, account)
	require.Equal(t, false, payload["store"])
placeholder

func TestOpenAISetupTokenChatCompletionsUsesCodexTransform(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"gpt-5.4","messages":[{"role":"system","content":"setup instructions"placeholder,{"role":"user","content":"hello"placeholder],"stream":falseplaceholder`)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"error":{"type":"invalid_request_error","message":"stop after request capture"placeholderplaceholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{cfg: &config.Config{placeholder, httpUpstream: upstreamplaceholder
	account := openAISetupTokenCompatAccount(71)

	result, err := svc.ForwardAsChatCompletions(context.Background(), c, account, body, "", "gpt-5.4")

placeholder
	require.Nil(t, result)
	require.NotNil(t, upstream.lastReq)
	require.Equal(t, chatgptCodexURL, upstream.lastReq.URL.String())
	require.Equal(t, "Bearer setup-token-value", upstream.lastReq.Header.Get("Authorization"))
	require.Equal(t, "chatgpt-setup", upstream.lastReq.Header.Get("chatgpt-account-id"))
	require.NotEmpty(t, upstream.lastReq.Header.Get("originator"))
	require.Equal(t, "setup instructions", gjson.GetBytes(upstream.lastBody, "instructions").String())
	require.Equal(t, int64(1), gjson.GetBytes(upstream.lastBody, "input.#").Int())
	require.Equal(t, "user", gjson.GetBytes(upstream.lastBody, "input.0.role").String())
	require.NotEmpty(t, gjson.GetBytes(upstream.lastBody, "prompt_cache_key").String())
placeholder

func TestOpenAISetupTokenMessagesUsesCodexBridgeAndTurnState(t *testing.T) {
	gin.SetMode(gin.TestMode)

	firstResp := openAICompatSSECompletedResponse("resp_setup_first", "gpt-5.4")
	firstResp.Header.Set("x-codex-turn-state", "turn_state_setup")
	upstream := &httpUpstreamRecorder{responses: []*http.Response{
		firstResp,
		openAICompatSSECompletedResponse("resp_setup_second", "gpt-5.4"),
placeholderplaceholder
	svc := &OpenAIGatewayService{
		cfg:          &config.Config{Security: config.SecurityConfig{URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholderplaceholderplaceholder,
		httpUpstream: upstream,
placeholder
	account := openAISetupTokenCompatAccount(72)

	messages := make([]string, 0, openAICompatAnthropicReplayMaxTailMessages+3)
	for i := 0; i < openAICompatAnthropicReplayMaxTailMessages+3; i++ {
		messages = append(messages, `{"role":"user","content":"message-`+fmt.Sprintf("%02d", i)+`"placeholder`)
placeholder
	firstBody := []byte(`{"model":"claude-sonnet-4-5","max_tokens":16,"messages":[` + strings.Join(messages, ",") + `],"stream":falseplaceholder`)
	firstRec := httptest.NewRecorder()
	firstCtx, _ := gin.CreateTestContext(firstRec)
	firstCtx.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(firstBody))
	firstCtx.Request.Header.Set("Content-Type", "application/json")

	firstResult, err := svc.ForwardAsAnthropic(context.Background(), firstCtx, account, firstBody, "stable-cache-key", "gpt-5.4")

placeholder
	require.NotNil(t, firstResult)
	require.True(t, isOpenAICompatMessagesBridgeContext(firstCtx))
	require.Equal(t, int64(openAICompatAnthropicReplayMaxTailMessages+4), gjson.GetBytes(upstream.bodies[0], "input.#").Int())
	require.Equal(t, "developer", gjson.GetBytes(upstream.bodies[0], "input.0.role").String())
	require.Contains(t, gjson.GetBytes(upstream.bodies[0], "input.0.content.0.text").String(), openAICompatClaudeCodeTodoGuardMarker)
	require.Equal(t, "message-00", gjson.GetBytes(upstream.bodies[0], "input.1.content.0.text").String())
	require.False(t, gjson.GetBytes(upstream.bodies[0], "prompt_cache_key").Exists())
	require.Equal(t, chatgptCodexURL, upstream.requests[0].URL.String())
	require.Equal(t, "Bearer setup-token-value", upstream.requests[0].Header.Get("Authorization"))
	require.Equal(t, "chatgpt-setup", upstream.requests[0].Header.Get("chatgpt-account-id"))
	requireOpenAIMessagesCodexIdentity(t, upstream.requests[0], codexCLIUserAgent, "codex-tui")
	require.Empty(t, upstream.requests[0].Header.Get("x-codex-turn-state"))

	secondBody := []byte(`{"model":"claude-sonnet-4-5","max_tokens":16,"messages":[{"role":"user","content":"next"placeholder],"stream":falseplaceholder`)
	secondRec := httptest.NewRecorder()
	secondCtx, _ := gin.CreateTestContext(secondRec)
	secondCtx.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(secondBody))
	secondCtx.Request.Header.Set("Content-Type", "application/json")

	secondResult, err := svc.ForwardAsAnthropic(context.Background(), secondCtx, account, secondBody, "stable-cache-key", "gpt-5.4")

placeholder
	require.NotNil(t, secondResult)
	require.True(t, isOpenAICompatMessagesBridgeContext(secondCtx))
	require.Equal(t, "turn_state_setup", upstream.requests[1].Header.Get("x-codex-turn-state"))
	require.Equal(t, generateSessionUUID(isolateOpenAISessionID(0, "stable-cache-key")), upstream.requests[1].Header.Get("session_id"))
	require.Empty(t, upstream.requests[1].Header.Get("conversation_id"))
	requireOpenAIMessagesCodexIdentity(t, upstream.requests[1], codexCLIUserAgent, "codex-tui")
placeholder

func openAISetupTokenCompatAccount(id int64) *Account {
placeholder
		ID:          id,
		Name:        "openai-setup-token",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeSetupToken,
		Concurrency: 1,
placeholder
			"access_token":       "setup-token-value",
			"chatgpt_account_id": "chatgpt-setup",
	placeholder,
placeholder
placeholder
