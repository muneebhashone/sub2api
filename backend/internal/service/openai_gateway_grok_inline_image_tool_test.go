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

	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestForwardGrokChatViaResponsesDropsRedundantViewImage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := grokChatInlineImageRequest()
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	c.Set("api_key", &APIKey{ID: 7991placeholder)

	account := grokChatBridgeTestAccount(799)
	repo := &grokQuotaAccountRepo{mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
		accountsByID: map[int64]*Account{account.ID: accountplaceholder,
placeholderplaceholder
	upstream := &httpUpstreamRecorder{resp: grokChatBridgeCompletedResponse("resp_chat_image", 0)placeholder
	svc := &OpenAIGatewayService{
		httpUpstream:      upstream,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		accountRepo:       repo,
placeholder

	result, err := svc.ForwardAsChatCompletions(context.Background(), c, account, body, "", "")
placeholder
	require.NotNil(t, result)
	require.Equal(t, xai.DefaultCLIBaseURL+"/responses", upstream.lastReq.URL.String())
	require.Equal(t, "input_image", gjson.GetBytes(upstream.lastBody, "input.0.content.1.type").String())
	assertGrokUpstreamKeepsOtherToolAndDropsViewImage(t, upstream.lastBody, "tools.#(name==\"%s\")")
placeholder

func TestForwardGrokRawChatDropsRedundantViewImage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := grokChatInlineImageRequest()
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))

	account := &Account{
		ID: 800, Platform: PlatformGrok, Type: AccountTypeAPIKey, Concurrency: 1,
placeholder"api_key": "test-key", "base_url": "https://grok.example.test/v1"placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(
			`{"id":"chatcmpl","object":"chat.completion","model":"grok-4.6","choices":[],"usage":{"prompt_tokens":1,"completion_tokens":1placeholderplaceholder`,
		)),
placeholderplaceholder
	svc := &OpenAIGatewayService{cfg: rawChatCompletionsTestConfig(), httpUpstream: upstreamplaceholder

	result, err := svc.ForwardAsChatCompletions(context.Background(), c, account, body, "", "")
placeholder
	require.NotNil(t, result)
	require.Equal(t, "https://grok.example.test/v1/chat/completions", upstream.lastReq.URL.String())
	require.Equal(t, "image_url", gjson.GetBytes(upstream.lastBody, "messages.0.content.1.type").String())
	assertGrokUpstreamKeepsOtherToolAndDropsViewImage(t, upstream.lastBody, "tools.#(function.name==\"%s\")")
placeholder

func TestForwardGrokMessagesDropsRedundantViewImage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{
		"model":"grok-4.6","max_tokens":32,"stream":false,
		"messages":[{"role":"user","content":[
			{"type":"text","text":"What text is in this image?"placeholder,
			{"type":"image","source":{"type":"base64","media_type":"image/png","data":"AA=="placeholderplaceholder
		]placeholder],
		"tools":[
			{"name":"view_image","input_schema":{"type":"object","properties":{"path":{"type":"string"placeholderplaceholderplaceholderplaceholder,
			{"name":"shell_command","input_schema":{"type":"object","properties":{"cmd":{"type":"string"placeholderplaceholderplaceholderplaceholder
		],
		"tool_choice":{"type":"auto"placeholder
placeholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	c.Set("api_key", &APIKey{ID: 7992placeholder)

	account := healthyGrokOAuthGatewayTestAccount(801, "access-token")
	repo := &grokQuotaAccountRepo{mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
		accountsByID: map[int64]*Account{account.ID: accountplaceholder,
placeholderplaceholder
	upstream := &httpUpstreamRecorder{resp: grokMessagesSSECompletedResponse("resp_messages_image", 0)placeholder
	svc := &OpenAIGatewayService{
		httpUpstream:      upstream,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		accountRepo:       repo,
placeholder

	result, err := svc.ForwardAsAnthropic(context.Background(), c, account, body, "", "")
placeholder
	require.NotNil(t, result)
	require.Equal(t, xai.DefaultCLIBaseURL+"/responses", upstream.lastReq.URL.String())
	require.Equal(t, "input_image", gjson.GetBytes(upstream.lastBody, "input.0.content.1.type").String())
	assertGrokUpstreamKeepsOtherToolAndDropsViewImage(t, upstream.lastBody, "tools.#(name==\"%s\")")
placeholder

func TestStripRedundantGrokChatViewImageToolLeavesNonTargetRequestsByteExact(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		body string
placeholder{
		{
			name: "current turn has no inline image",
			body: `{"messages":[{"role":"user","content":"Inspect a local image"placeholder],"tools":[{"type":"function","function":{"name":"view_image"placeholderplaceholder]placeholder`,
	placeholder,
		{
			name: "inline image is only historical",
			body: `{"messages":[{"role":"user","content":[{"type":"image_url","image_url":{"url":"data:image/png;base64,AA=="placeholderplaceholder]placeholder,{"role":"assistant","content":"Done"placeholder,{"role":"user","content":"Inspect another local image"placeholder],"tools":[{"type":"function","function":{"name":"view_image"placeholderplaceholder]placeholder`,
	placeholder,
		{
			name: "view image is explicitly selected",
			body: `{"messages":[{"role":"user","content":[{"type":"image_url","image_url":{"url":"data:image/png;base64,AA=="placeholderplaceholder]placeholder],"tools":[{"type":"function","function":{"name":"view_image"placeholderplaceholder],"tool_choice":{"type":"function","function":{"name":"view_image"placeholderplaceholderplaceholder`,
	placeholder,
		{
			name: "required with view image as the only tool",
			body: `{"messages":[{"role":"user","content":[{"type":"image_url","image_url":{"url":"data:image/png;base64,AA=="placeholderplaceholder]placeholder],"tools":[{"type":"function","function":{"name":"view_image"placeholderplaceholder],"tool_choice":"required"placeholder`,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			body := []byte(tt.body)
			patched, err := stripRedundantGrokChatViewImageTool(body)
		placeholder
			require.Equal(t, body, patched)
	placeholder)
placeholder
placeholder

func TestStripRedundantGrokChatViewImageToolDropsOnlyToolMetadata(t *testing.T) {
	t.Parallel()
	body := []byte(`{
		"messages":[{"role":"user","content":[{"type":"image_url","image_url":{"url":"data:image/png;base64,AA=="placeholderplaceholder]placeholder],
		"tools":[{"type":"function","function":{"name":"view_image"placeholderplaceholder],
		"tool_choice":"auto",
		"parallel_tool_calls":true
placeholder`)

	patched, err := stripRedundantGrokChatViewImageTool(body)
placeholder
	require.False(t, gjson.GetBytes(patched, "tools").Exists())
	require.False(t, gjson.GetBytes(patched, "tool_choice").Exists())
	require.False(t, gjson.GetBytes(patched, "parallel_tool_calls").Exists())
placeholder

func grokChatInlineImageRequest() []byte {
	return []byte(`{
		"model":"grok-4.6",
		"messages":[{"role":"user","content":[
			{"type":"text","text":"What text is in this image?"placeholder,
			{"type":"image_url","image_url":{"url":"data:image/png;base64,AA=="placeholderplaceholder
		]placeholder],
		"stream":false,
		"tools":[
			{"type":"function","function":{"name":"view_image","parameters":{"type":"object","properties":{"path":{"type":"string"placeholderplaceholderplaceholderplaceholderplaceholder,
			{"type":"function","function":{"name":"shell_command","parameters":{"type":"object","properties":{"cmd":{"type":"string"placeholderplaceholderplaceholderplaceholderplaceholder
		],
		"tool_choice":"auto"
placeholder`)
placeholder

func assertGrokUpstreamKeepsOtherToolAndDropsViewImage(t *testing.T, body []byte, pathTemplate string) {
placeholder
	require.False(t, gjson.GetBytes(body, strings.Replace(pathTemplate, "%s", "view_image", 1)).Exists(), string(body))
	require.True(t, gjson.GetBytes(body, strings.Replace(pathTemplate, "%s", "shell_command", 1)).Exists(), string(body))
placeholder
