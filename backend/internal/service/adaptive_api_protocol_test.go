//go:build unit

package service

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/openai_compat"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func adaptiveProtocolTestAccount(platform string, baseURLs map[string]any) *Account {
placeholder
		ID:          701,
		Name:        "adaptive-cn",
		Platform:    platform,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":       "sk-test",
			"api_protocol":  APIProtocolAdaptive,
			"account_mode":  AccountModePayG,
			"api_base_urls": baseURLs,
	placeholder,
placeholder
placeholder

func adaptiveProtocolTestContext(path string, body []byte) *gin.Context {
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	return c
placeholder

type cnProtocolIngressCase struct {
	name    string
	path    string
	body    []byte
	forward func(*OpenAIGatewayService, *gin.Context, *Account, []byte) error
placeholder

func cnProtocolIngressCases() []cnProtocolIngressCase {
	return []cnProtocolIngressCase{
		{
			name: "chat completions",
			path: "/v1/chat/completions",
			body: []byte(`{"model":"deepseek-chat","messages":[{"role":"user","content":"hello"placeholder],"stream":falseplaceholder`),
			forward: func(svc *OpenAIGatewayService, c *gin.Context, account *Account, body []byte) error {
				_, err := svc.ForwardAsChatCompletions(context.Background(), c, account, body, "", "")
				return err
		placeholder,
	placeholder,
		{
			name: "messages",
			path: "/v1/messages",
			body: []byte(`{"model":"deepseek-chat","max_tokens":32,"messages":[{"role":"user","content":"hello"placeholder],"stream":falseplaceholder`),
			forward: func(svc *OpenAIGatewayService, c *gin.Context, account *Account, body []byte) error {
				_, err := svc.ForwardAsAnthropic(context.Background(), c, account, body, "", "")
				return err
		placeholder,
	placeholder,
		{
			name: "responses",
			path: "/v1/responses",
			body: []byte(`{"model":"deepseek-chat","input":"hello","stream":falseplaceholder`),
			forward: func(svc *OpenAIGatewayService, c *gin.Context, account *Account, body []byte) error {
				_, err := svc.Forward(context.Background(), c, account, body)
				return err
		placeholder,
	placeholder,
placeholder
placeholder

func TestAdaptiveProtocolRoutesChatCompletionsToNativeChat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"glm-4.7","messages":[{"role":"user","content":"hello"placeholder],"stream":falseplaceholder`)
	upstream := &httpUpstreamRecorder{err: errors.New("stop after capture")placeholder
	svc := &OpenAIGatewayService{cfg: rawChatCompletionsTestConfig(), httpUpstream: upstreamplaceholder
	account := adaptiveProtocolTestAccount(PlatformZhipu, map[string]any{
		APIProtocolChatCompletions: "http://chat.example",
		APIProtocolAnthropic:       "http://anthropic.example",
placeholder)

	_, err := svc.ForwardAsChatCompletions(context.Background(), adaptiveProtocolTestContext("/v1/chat/completions", body), account, body, "", "")
placeholder
	require.Equal(t, "http://chat.example/v1/chat/completions", upstream.lastReq.URL.String())
	require.True(t, gjson.GetBytes(upstream.lastBody, "messages").IsArray())
	require.False(t, gjson.GetBytes(upstream.lastBody, "input").Exists())
placeholder

func TestAdaptiveProtocolRoutesResponsesShapedChatToNativeResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"deepseek-v4","input":"hello","max_output_tokens":32,"stream":falseplaceholder`)
	upstream := &httpUpstreamRecorder{err: errors.New("stop after capture")placeholder
	svc := &OpenAIGatewayService{cfg: rawChatCompletionsTestConfig(), httpUpstream: upstreamplaceholder
	account := adaptiveProtocolTestAccount(PlatformDeepseek, map[string]any{
		APIProtocolChatCompletions: "http://chat.example",
		APIProtocolAnthropic:       "http://anthropic.example",
		APIProtocolResponses:       "http://responses.example",
placeholder)

	_, err := svc.ForwardAsChatCompletions(context.Background(), adaptiveProtocolTestContext("/v1/chat/completions", body), account, body, "", "")
placeholder
	require.Equal(t, "http://responses.example/responses", upstream.lastReq.URL.String())
	require.True(t, gjson.GetBytes(upstream.lastBody, "input").Exists())
	require.False(t, gjson.GetBytes(upstream.lastBody, "messages").Exists())
placeholder

func TestAdaptiveProtocolConvertsResponsesShapedChatForChatOnlyProvider(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"kimi-k2.5","input":"hello","max_output_tokens":32,"stream":falseplaceholder`)
	upstream := &httpUpstreamRecorder{err: errors.New("stop after capture")placeholder
	svc := &OpenAIGatewayService{cfg: rawChatCompletionsTestConfig(), httpUpstream: upstreamplaceholder
	account := adaptiveProtocolTestAccount(PlatformKimi, map[string]any{
		APIProtocolChatCompletions: "http://chat.example",
		APIProtocolAnthropic:       "http://anthropic.example",
placeholder)

	_, err := svc.ForwardAsChatCompletions(context.Background(), adaptiveProtocolTestContext("/v1/chat/completions", body), account, body, "", "")
placeholder
	require.Equal(t, "http://chat.example/v1/chat/completions", upstream.lastReq.URL.String())
	require.True(t, gjson.GetBytes(upstream.lastBody, "messages").IsArray())
	require.False(t, gjson.GetBytes(upstream.lastBody, "input").Exists())
placeholder

func TestAdaptiveProtocolRoutesMessagesToNativeAnthropic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"glm-4.7","max_tokens":32,"messages":[{"role":"user","content":"hello"placeholder],"stream":falseplaceholder`)
	upstream := &httpUpstreamRecorder{err: errors.New("stop after capture")placeholder
	svc := &OpenAIGatewayService{cfg: rawChatCompletionsTestConfig(), httpUpstream: upstreamplaceholder
	account := adaptiveProtocolTestAccount(PlatformZhipu, map[string]any{
		APIProtocolChatCompletions: "http://chat.example",
		APIProtocolAnthropic:       "http://anthropic.example",
placeholder)

	_, err := svc.ForwardAsAnthropic(context.Background(), adaptiveProtocolTestContext("/v1/messages", body), account, body, "", "")
placeholder
	require.Equal(t, "http://anthropic.example/v1/messages", upstream.lastReq.URL.String())
	require.Equal(t, "glm-4.7", gjson.GetBytes(upstream.lastBody, "model").String())
placeholder

func TestAdaptiveProtocolConvertsKimiResponsesToChatCompletions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"kimi-k2.5","input":"hello","stream":falseplaceholder`)
	upstream := &httpUpstreamRecorder{err: errors.New("stop after capture")placeholder
	svc := &OpenAIGatewayService{cfg: rawChatCompletionsTestConfig(), httpUpstream: upstreamplaceholder
	account := adaptiveProtocolTestAccount(PlatformKimi, map[string]any{
		APIProtocolChatCompletions: "http://chat.example",
		APIProtocolAnthropic:       "http://anthropic.example",
placeholder)

	_, err := svc.Forward(context.Background(), adaptiveProtocolTestContext("/v1/responses", body), account, body)
placeholder
	require.Equal(t, "http://chat.example/v1/chat/completions", upstream.lastReq.URL.String())
	require.True(t, gjson.GetBytes(upstream.lastBody, "messages").IsArray())
	require.False(t, gjson.GetBytes(upstream.lastBody, "input").Exists())
placeholder

func TestAdaptiveProtocolRoutesDeepSeekResponsesToNativeResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"deepseek-v4","input":"hello","max_output_tokens":32,"store":true,"previous_response_id":"resp_old","stream":falseplaceholder`)
	upstream := &httpUpstreamRecorder{err: errors.New("stop after capture")placeholder
	svc := &OpenAIGatewayService{cfg: rawChatCompletionsTestConfig(), httpUpstream: upstreamplaceholder
	account := adaptiveProtocolTestAccount(PlatformDeepseek, map[string]any{
		APIProtocolChatCompletions: "http://chat.example",
		APIProtocolAnthropic:       "http://anthropic.example",
		APIProtocolResponses:       "http://responses.example",
placeholder)

	_, err := svc.Forward(context.Background(), adaptiveProtocolTestContext("/v1/responses", body), account, body)
placeholder
	require.Equal(t, "http://responses.example/responses", upstream.lastReq.URL.String())
	require.False(t, gjson.GetBytes(upstream.lastBody, "store").Bool())
	require.False(t, gjson.GetBytes(upstream.lastBody, "previous_response_id").Exists())
	require.Equal(t, int64(32), gjson.GetBytes(upstream.lastBody, "max_output_tokens").Int())
	require.False(t, gjson.GetBytes(upstream.lastBody, "instructions").Exists())
placeholder

func TestFixedCNChatProtocolOverridesStaleResponsesMode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range cnProtocolIngressCases() {
		t.Run(tc.name, func(t *testing.T) {
			upstream := &httpUpstreamRecorder{err: errors.New("stop after capture")placeholder
			svc := &OpenAIGatewayService{cfg: rawChatCompletionsTestConfig(), httpUpstream: upstreamplaceholder
			account := adaptiveProtocolTestAccount(PlatformDeepseek, nil)
			account.Credentials["api_protocol"] = APIProtocolChatCompletions
			account.Credentials["base_url"] = "http://chat.example"
			account.Extra = map[string]any{
				openai_compat.ExtraKeyResponsesMode: string(openai_compat.ResponsesSupportModeForceResponses),
		placeholder

			err := tc.forward(svc, adaptiveProtocolTestContext(tc.path, tc.body), account, tc.body)

		placeholder
			require.Equal(t, "http://chat.example/v1/chat/completions", upstream.lastReq.URL.String())
	placeholder)
placeholder
placeholder

func TestFixedCNResponsesProtocolOverridesStaleChatMode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range cnProtocolIngressCases() {
		t.Run(tc.name, func(t *testing.T) {
			upstream := &httpUpstreamRecorder{err: errors.New("stop after capture")placeholder
			svc := &OpenAIGatewayService{cfg: rawChatCompletionsTestConfig(), httpUpstream: upstreamplaceholder
			account := adaptiveProtocolTestAccount(PlatformDeepseek, nil)
			account.Credentials["api_protocol"] = APIProtocolResponses
			account.Credentials["base_url"] = "http://responses.example"
			account.Extra = map[string]any{
				openai_compat.ExtraKeyResponsesMode: string(openai_compat.ResponsesSupportModeForceChatCompletions),
		placeholder

			err := tc.forward(svc, adaptiveProtocolTestContext(tc.path, tc.body), account, tc.body)

		placeholder
			require.Equal(t, "http://responses.example/responses", upstream.lastReq.URL.String())
	placeholder)
placeholder
placeholder
