//go:build unit

package service

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func adaptiveCNAccountTestAccount(id int64, platform string) *Account {
placeholder
		ID:          id,
		Name:        "adaptive-cn-test",
		Platform:    platform,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"api_key":      "sk-adaptive-test",
			"api_protocol": APIProtocolAdaptive,
			"api_base_urls": map[string]any{
				APIProtocolChatCompletions: "http://chat.example/v1",
				APIProtocolAnthropic:       "http://anthropic.example",
				APIProtocolResponses:       "http://responses.example",
		placeholder,
	placeholder,
placeholder
placeholder

func adaptiveCNAccountTestService(account *Account, responses ...*http.Response) (*AccountTestService, *httpUpstreamRecorder) {
	repo := &openAIAccountTestRepo{
		mockAccountRepoForGemini: mockAccountRepoForGemini{
			accountsByID: map[int64]*Account{account.ID: accountplaceholder,
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{responses: responsesplaceholder
	return &AccountTestService{
		accountRepo:  repo,
		httpUpstream: upstream,
		cfg:          rawChatCompletionsTestConfig(),
placeholder, upstream
placeholder

func adaptiveCNChatTestResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(`data: {"choices":[{"delta":{"content":"chat ok"placeholder,"finish_reason":"stop"placeholder]placeholder

data: [DONE]

`)),
placeholder
placeholder

func adaptiveCNAnthropicTestResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(`data: {"type":"content_block_delta","delta":{"text":"anthropic ok"placeholderplaceholder

data: {"type":"message_stop"placeholder

`)),
placeholder
placeholder

func adaptiveCNResponsesTestResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(`data: {"type":"response.output_text.delta","delta":"responses ok"placeholder

data: {"type":"response.completed"placeholder

`)),
placeholder
placeholder

func TestAccountTestService_AdaptiveChatOnlyProvidersTestChatAndAnthropicEndpoints(t *testing.T) {
	for index, testCase := range []struct {
		name     string
		platform string
		model    string
placeholder{
		{name: "Kimi", platform: PlatformKimi, model: "kimi-k2.5"placeholder,
		{name: "Zhipu", platform: PlatformZhipu, model: "glm-4.7"placeholder,
placeholder {
		t.Run(testCase.name, func(t *testing.T) {
			account := adaptiveCNAccountTestAccount(int64(301+index), testCase.platform)
			svc, upstream := adaptiveCNAccountTestService(
				account,
				adaptiveCNChatTestResponse(),
				adaptiveCNAnthropicTestResponse(),
			)
			c, recorder := newTestContext()

			err := svc.TestAccountConnection(c, account.ID, testCase.model, "hello", AccountTestModeDefault)

		placeholder
			require.Len(t, upstream.requests, 2)
			require.Equal(t, "http://chat.example/v1/chat/completions", upstream.requests[0].URL.String())
			require.Equal(t, "http://anthropic.example/v1/messages", upstream.requests[1].URL.String())
			require.Equal(t, "Bearer sk-adaptive-test", upstream.requests[0].Header.Get("Authorization"))
			require.Equal(t, "sk-adaptive-test", upstream.requests[1].Header.Get("x-api-key"))
			require.Equal(t, 1, strings.Count(recorder.Body.String(), `"type":"test_start"`))
			require.Equal(t, 1, strings.Count(recorder.Body.String(), `"type":"test_complete"`))
			require.Contains(t, recorder.Body.String(), "已通过原生 /v1/messages 验证")
	placeholder)
placeholder
placeholder

func TestAccountTestService_AdaptiveDeepSeekAlsoTestsResponsesEndpoint(t *testing.T) {
	account := adaptiveCNAccountTestAccount(302, PlatformDeepseek)
	svc, upstream := adaptiveCNAccountTestService(
		account,
		adaptiveCNChatTestResponse(),
		adaptiveCNAnthropicTestResponse(),
		adaptiveCNResponsesTestResponse(),
	)
	c, recorder := newTestContext()

	err := svc.TestAccountConnection(c, account.ID, "deepseek-chat", "", AccountTestModeDefault)

placeholder
	require.Len(t, upstream.requests, 3)
	require.Equal(t, "http://responses.example/responses", upstream.requests[2].URL.String())
	require.Equal(t, HTTPUpstreamProfileOpenAI, HTTPUpstreamProfileFromContext(upstream.requests[2].Context()))
	require.Equal(t, "Bearer sk-adaptive-test", upstream.requests[2].Header.Get("Authorization"))
	require.True(t, gjson.GetBytes(upstream.bodies[2], "stream").Bool())
	require.False(t, gjson.GetBytes(upstream.bodies[2], "store").Bool())
	require.False(t, gjson.GetBytes(upstream.bodies[2], "instructions").Exists())
	require.Equal(t, 1, strings.Count(recorder.Body.String(), `"type":"test_complete"`))
	require.Contains(t, recorder.Body.String(), "已通过原生 /responses 验证")
placeholder

func TestAccountTestService_AdaptiveStopsAndNamesFailingEndpoint(t *testing.T) {
	account := adaptiveCNAccountTestAccount(303, PlatformDeepseek)
	svc, upstream := adaptiveCNAccountTestService(
		account,
		adaptiveCNChatTestResponse(),
		newJSONResponse(http.StatusNotFound, `{"error":{"message":"missing messages route"placeholderplaceholder`),
	)
	c, recorder := newTestContext()

	err := svc.TestAccountConnection(c, account.ID, "deepseek-chat", "", AccountTestModeDefault)

placeholder
	require.Contains(t, err.Error(), "Adaptive Anthropic endpoint returned 404")
	require.Len(t, upstream.requests, 2)
	require.Contains(t, recorder.Body.String(), `"type":"error"`)
	require.NotContains(t, recorder.Body.String(), `"type":"test_complete"`)
placeholder

func TestAccountTestService_AdaptiveRejectsInvalidAnthropicSuccessBody(t *testing.T) {
	account := adaptiveCNAccountTestAccount(305, PlatformKimi)
	svc, upstream := adaptiveCNAccountTestService(
		account,
		adaptiveCNChatTestResponse(),
		newJSONResponse(http.StatusOK, `<html>not an Anthropic stream</html>`),
	)
	c, recorder := newTestContext()

	err := svc.TestAccountConnection(c, account.ID, "kimi-k2.5", "", AccountTestModeDefault)

placeholder
	require.Contains(t, err.Error(), "Adaptive Anthropic stream ended before message_stop")
	require.Len(t, upstream.requests, 2)
	require.Contains(t, recorder.Body.String(), `"type":"error"`)
	require.NotContains(t, recorder.Body.String(), `"type":"test_complete"`)
placeholder

func TestAccountTestService_FixedCNChatProtocolStillTestsOnlyChatEndpoint(t *testing.T) {
	account := adaptiveCNAccountTestAccount(304, PlatformZhipu)
	account.Credentials["api_protocol"] = APIProtocolChatCompletions
	account.Credentials["base_url"] = "http://fixed-chat.example/v1"
	svc, upstream := adaptiveCNAccountTestService(account, adaptiveCNChatTestResponse())
	c, recorder := newTestContext()

	err := svc.TestAccountConnection(c, account.ID, "glm-4.7", "", AccountTestModeDefault)

placeholder
	require.Len(t, upstream.requests, 1)
	require.Equal(t, "http://fixed-chat.example/v1/chat/completions", upstream.requests[0].URL.String())
	require.Equal(t, 1, strings.Count(recorder.Body.String(), `"type":"test_complete"`))
placeholder

func anthropicProtocolCNAccount(id int64, platform string, credentials map[string]any) *Account {
	base := map[string]any{
		"api_key":      "sk-anthropic-test",
		"api_protocol": APIProtocolAnthropic,
placeholder
	for key, value := range credentials {
		base[key] = value
placeholder
placeholder
		ID:          id,
		Name:        "anthropic-protocol-cn-test",
		Platform:    platform,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Concurrency: 1,
		Credentials: base,
placeholder
placeholder

func TestAccountTestService_AnthropicProtocolProbesNativeEndpointWithoutBetaQuery(t *testing.T) {
	account := anthropicProtocolCNAccount(311, PlatformZhipu, map[string]any{
		"base_url": "https://open.bigmodel.cn/api/anthropic",
placeholder)
	svc, upstream := adaptiveCNAccountTestService(account, adaptiveCNAnthropicTestResponse())
	c, recorder := newTestContext()

	err := svc.TestAccountConnection(c, account.ID, "glm-4.7", "", AccountTestModeDefault)

placeholder
	require.Len(t, upstream.requests, 1)
	req := upstream.requests[0]
	// Native Anthropic path without the ?beta=true suffix the generic Claude tester appends.
	require.Equal(t, "https://open.bigmodel.cn/api/anthropic/v1/messages", req.URL.String())
	require.Empty(t, req.URL.RawQuery)
	require.Equal(t, "sk-anthropic-test", req.Header.Get("x-api-key"))
	require.Equal(t, "2023-06-01", req.Header.Get("anthropic-version"))
	require.Contains(t, recorder.Body.String(), `"type":"test_complete"`)
placeholder

func TestAccountTestService_AnthropicProtocolFallsBackToProviderDefaultNotAnthropicDotCom(t *testing.T) {
	// base_url intentionally absent: forwarding resolves the per-platform default
	// Anthropic endpoint. The old fall-through probed https://api.anthropic.com
	// with the provider's API key.
	account := anthropicProtocolCNAccount(312, PlatformZhipu, nil)
	svc, upstream := adaptiveCNAccountTestService(account, adaptiveCNAnthropicTestResponse())
	c, _ := newTestContext()

	err := svc.TestAccountConnection(c, account.ID, "glm-4.7", "", AccountTestModeDefault)

placeholder
	require.Len(t, upstream.requests, 1)
	require.Equal(t, "https://open.bigmodel.cn/api/anthropic/v1/messages", upstream.requests[0].URL.String())
placeholder

func TestAccountTestService_AnthropicProtocolRejectsOpenAICompatBaseURL(t *testing.T) {
	account := anthropicProtocolCNAccount(313, PlatformZhipu, map[string]any{
		"base_url": "https://open.bigmodel.cn/api/paas/v4",
placeholder)
	svc, upstream := adaptiveCNAccountTestService(account)
	c, recorder := newTestContext()

	err := svc.TestAccountConnection(c, account.ID, "glm-4.7", "", AccountTestModeDefault)

placeholder
	// Fails fast locally: no upstream request with the wrong endpoint shape.
	require.Empty(t, upstream.requests)
	require.Contains(t, recorder.Body.String(), "looks like an OpenAI-compatible endpoint")
	require.Contains(t, recorder.Body.String(), "https://open.bigmodel.cn/api/anthropic")
placeholder

func TestAccountTestService_AnthropicProtocol401MarksAccountError(t *testing.T) {
	account := anthropicProtocolCNAccount(314, PlatformKimi, map[string]any{
		"base_url": "https://api.moonshot.cn/anthropic",
placeholder)
	svc, _ := adaptiveCNAccountTestService(
		account,
		newJSONResponse(http.StatusUnauthorized, `{"error":{"message":"invalid key"placeholderplaceholder`),
	)
	c, _ := newTestContext()

	err := svc.TestAccountConnection(c, account.ID, "kimi-k2.5", "", AccountTestModeDefault)

placeholder
	require.Contains(t, err.Error(), "Anthropic endpoint returned 401")
	repo := svc.accountRepo.(*openAIAccountTestRepo)
	require.Equal(t, account.ID, repo.setErrorID)
placeholder
