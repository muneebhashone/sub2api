//go:build unit

package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestGrokChatResponsesBridgeEligibility(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		body   string
		want   bool
		reason string
placeholder{
		{
			name: "plain text chat",
			body: `{"model":"grok","messages":[{"role":"system","content":"concise"placeholder,{"role":"user","content":"hi"placeholder],"stream":falseplaceholder`,
			want: true,
	placeholder,
		{
			name: "safe generation options",
			body: `{"model":"grok","messages":[{"role":"user","content":"hi"placeholder],"stream":true,"stream_options":{"include_usage":trueplaceholder,"max_completion_tokens":256,"temperature":0.2,"top_p":0.9,"prompt_cache_key":"session","tools":[],"functions":null,"tool_choice":"none"placeholder`,
			want: true,
	placeholder,
		{
			name:   "stop falls back",
			body:   `{"model":"grok","messages":[{"role":"user","content":"hi"placeholder],"stop":"done"placeholder`,
			reason: "unsupported_stop",
	placeholder,
		{
			name:   "developer role falls back",
			body:   `{"model":"grok","messages":[{"role":"developer","content":"rules"placeholder,{"role":"user","content":"hi"placeholder]placeholder`,
			reason: "unsupported_message_role_developer",
	placeholder,
		{
			name:   "image content falls back",
			body:   `{"model":"grok","messages":[{"role":"user","content":[{"type":"image_url","image_url":{"url":"data:image/png;base64,QQ=="placeholderplaceholder]placeholder]placeholder`,
			reason: "non_text_message_content",
	placeholder,
		{
			name:   "function tools fall back",
			body:   `{"model":"grok","messages":[{"role":"user","content":"hi"placeholder],"tools":[{"type":"function","function":{"name":"lookup"placeholderplaceholder]placeholder`,
			reason: "unsupported_tools",
	placeholder,
		{
			name:   "automatic tool choice falls back",
			body:   `{"model":"grok","messages":[{"role":"user","content":"hi"placeholder],"tools":[],"tool_choice":"auto"placeholder`,
			reason: "unsupported_tool_choice",
	placeholder,
		{
			name:   "reasoning effort falls back because conversion adds summary",
			body:   `{"model":"grok","messages":[{"role":"user","content":"hi"placeholder],"reasoning_effort":"high"placeholder`,
			reason: "unsupported_reasoning_effort",
	placeholder,
		{
			name:   "both token limits fall back",
			body:   `{"model":"grok","messages":[{"role":"user","content":"hi"placeholder],"max_tokens":256,"max_completion_tokens":placeholder`,
			reason: "conflicting_max_tokens",
	placeholder,
		{
			name:   "empty message falls back",
			body:   `{"model":"grok","messages":[{"role":"assistant","content":""placeholder,{"role":"user","content":"hi"placeholder]placeholder`,
			reason: "empty_message_content",
	placeholder,
		{
			name:   "tool history falls back",
			body:   `{"model":"grok","messages":[{"role":"assistant","content":"","tool_calls":[]placeholder]placeholder`,
			reason: "unsafe_message_field_tool_calls",
	placeholder,
		{
			name:   "unknown field falls back",
			body:   `{"model":"grok","messages":[{"role":"user","content":"hi"placeholder],"seed":7placeholder`,
			reason: "unknown_field_seed",
	placeholder,
		{
			name:   "small max tokens falls back because conversion clamps it",
			body:   `{"model":"grok","messages":[{"role":"user","content":"hi"placeholder],"max_tokens":32placeholder`,
			reason: "unsafe_max_tokens",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, reason := grokChatResponsesBridgeEligibility([]byte(tt.body))
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.reason, reason)
	placeholder)
placeholder
placeholder

func TestGrokChatResponsesRuntimeEligibility(t *testing.T) {
	t.Parallel()
	require.True(t, grokChatResponsesRuntimeEligible("grok-4.5", "isolated-id"))
	require.False(t, grokChatResponsesRuntimeEligible("grok-4.3", "isolated-id"))
	require.False(t, grokChatResponsesRuntimeEligible("grok-4.5-build-free", "isolated-id"))
	require.False(t, grokChatResponsesRuntimeEligible("grok-4.5", ""))
placeholder

func TestForwardGrokChatViaResponsesNonStreamingCachesAndReturnsChat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"model":"grok","messages":[{"role":"system","content":"be concise"placeholder,{"role":"user","content":"hi"placeholder],"stream":false,"prompt_cache_key":"stable-session","tools":[],"functions":null,"tool_choice":"none"placeholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, grokChatRawEndpoint, bytes.NewReader(body))
	c.Set("api_key", &APIKey{ID: 7101placeholder)

	account := grokChatBridgeTestAccount(71)
	repo := &grokQuotaAccountRepo{mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
		accountsByID: map[int64]*Account{account.ID: accountplaceholder,
placeholderplaceholder
	upstream := &httpUpstreamRecorder{resp: grokChatBridgeCompletedResponse("resp_grok_chat_cache", 9856)placeholder
	svc := &OpenAIGatewayService{
		httpUpstream:      upstream,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		accountRepo:       repo,
placeholder

	result, err := svc.ForwardAsChatCompletions(context.Background(), c, account, body, "", "")
placeholder
	require.NotNil(t, result)
	require.Equal(t, xai.DefaultCLIBaseURL+"/responses", upstream.lastReq.URL.String())
	require.Equal(t, grokChatResponsesEndpoint, result.UpstreamEndpoint)
	require.Equal(t, "grok-4.5", result.UpstreamModel)
	require.Equal(t, 9908, result.Usage.InputTokens)
	require.Equal(t, 12, result.Usage.OutputTokens)
	require.Equal(t, 9856, result.Usage.CacheReadInputTokens)

	identity := gjson.GetBytes(upstream.lastBody, "prompt_cache_key").String()
	require.NotEmpty(t, identity)
	require.NotEqual(t, "stable-session", identity)
	require.Equal(t, identity, upstream.lastReq.Header.Get(grokConversationIDHeader))
	require.Equal(t, "web_search", gjson.GetBytes(upstream.lastBody, "tools.0.type").String())
	require.Equal(t, "x_search", gjson.GetBytes(upstream.lastBody, "tools.1.type").String())
	require.Equal(t, grokFreeCacheDisabledToolChoice, gjson.GetBytes(upstream.lastBody, "tool_choice").String())
	require.True(t, gjson.GetBytes(upstream.lastBody, "stream").Bool())
	require.Equal(t, "system", gjson.GetBytes(upstream.lastBody, "input.0.role").String())
	require.Equal(t, "user", gjson.GetBytes(upstream.lastBody, "input.1.role").String())
	require.False(t, gjson.GetBytes(upstream.lastBody, "instructions").Exists())
	require.False(t, gjson.GetBytes(upstream.lastBody, "include").Exists())
	require.False(t, gjson.GetBytes(upstream.lastBody, "store").Exists())

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "cached ok", gjson.Get(recorder.Body.String(), "choices.0.message.content").String())
	require.Equal(t, int64(9856), gjson.Get(recorder.Body.String(), "usage.prompt_tokens_details.cached_tokens").Int())
	require.NotNil(t, repo.updates[account.ID][grokQuotaSnapshotExtraKey])
placeholder

func TestForwardGrokChatViaResponsesStreamingPropagatesCachedUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"model":"grok","messages":[{"role":"user","content":"hi"placeholder],"stream":trueplaceholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, grokChatRawEndpoint, bytes.NewReader(body))
	c.Set("api_key", &APIKey{ID: 7201placeholder)

	account := grokChatBridgeTestAccount(72)
	repo := &grokQuotaAccountRepo{mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
		accountsByID: map[int64]*Account{account.ID: accountplaceholder,
placeholderplaceholder
	upstream := &httpUpstreamRecorder{resp: grokChatBridgeCompletedResponse("resp_grok_chat_stream", 4096)placeholder
	svc := &OpenAIGatewayService{
		httpUpstream:      upstream,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		accountRepo:       repo,
placeholder

	result, err := svc.ForwardAsChatCompletions(context.Background(), c, account, body, "", "")
placeholder
	require.NotNil(t, result)
	require.True(t, result.Stream)
	require.Equal(t, grokChatResponsesEndpoint, result.UpstreamEndpoint)
	require.Equal(t, 4096, result.Usage.CacheReadInputTokens)
	require.Contains(t, recorder.Header().Get("Content-Type"), "text/event-stream")
	require.Contains(t, recorder.Body.String(), `"content":"cached ok"`)
	require.Contains(t, recorder.Body.String(), `"cached_tokens":4096`)
	require.Contains(t, recorder.Body.String(), "data: [DONE]")
placeholder

func TestForwardGrokChatRuntimeGateFallsBackToRaw(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		setAPIKey    bool
		mappedModel  string
		wantUpstream string
placeholder{
		{name: "missing cache identity", wantUpstream: "grok-4.5"placeholder,
		{name: "non cache capable mapped model", setAPIKey: true, mappedModel: "grok-4.3", wantUpstream: "grok-4.3"placeholder,
placeholder

	for index, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := []byte(`{"model":"grok","messages":[{"role":"user","content":"hi"placeholder],"stream":falseplaceholder`)
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request = httptest.NewRequest(http.MethodPost, grokChatRawEndpoint, bytes.NewReader(body))
			if tt.setAPIKey {
				c.Set("api_key", &APIKey{ID: int64(7301 + index)placeholder)
		placeholder

			account := grokChatBridgeTestAccount(int64(73 + index))
			if tt.mappedModel != "" {
				account.Credentials["model_mapping"] = map[string]any{"grok": tt.mappedModelplaceholder
		placeholder
			repo := &grokQuotaAccountRepo{mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
				accountsByID: map[int64]*Account{account.ID: accountplaceholder,
		placeholderplaceholder
			upstream := &httpUpstreamRecorder{resp: &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
				Body: io.NopCloser(strings.NewReader(
					`{"id":"chat_raw","object":"chat.completion","model":"` + tt.wantUpstream + `","choices":[{"index":0,"message":{"role":"assistant","content":"raw ok"placeholder,"finish_reason":"stop"placeholder],"usage":{"prompt_tokens":2,"completion_tokens":1,"total_tokens":3placeholderplaceholder`,
				)),
		placeholderplaceholder
			svc := &OpenAIGatewayService{
				httpUpstream:      upstream,
				grokTokenProvider: NewGrokTokenProvider(repo, nil),
				accountRepo:       repo,
		placeholder

			result, err := svc.ForwardAsChatCompletions(context.Background(), c, account, body, "", "")
		placeholder
			require.NotNil(t, result)
			require.Equal(t, xai.DefaultCLIBaseURL+"/chat/completions", upstream.lastReq.URL.String())
			require.Equal(t, grokChatRawEndpoint, result.UpstreamEndpoint)
			require.Equal(t, tt.wantUpstream, result.UpstreamModel)
			require.False(t, gjson.GetBytes(upstream.lastBody, "tools").Exists())
			require.Equal(t, "raw ok", gjson.Get(recorder.Body.String(), "choices.0.message.content").String())
	placeholder)
placeholder
placeholder

func TestForwardGrokChatViaResponses429UsesGrokRateLimitPolicy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"model":"grok","messages":[{"role":"user","content":"hi"placeholder],"stream":falseplaceholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, grokChatRawEndpoint, bytes.NewReader(body))
	c.Set("api_key", &APIKey{ID: 7501placeholder)

	account := grokChatBridgeTestAccount(75)
	repo := &grokQuotaAccountRepo{mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
		accountsByID: map[int64]*Account{account.ID: accountplaceholder,
placeholderplaceholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header: http.Header{
			"Content-Type": []string{"application/json"placeholder,
			"Retry-After":  []string{"45"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(`{"error":{"message":"rate limited"placeholderplaceholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{
		httpUpstream:      upstream,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		accountRepo:       repo,
placeholder
	before := time.Now()

	result, err := svc.ForwardAsChatCompletions(context.Background(), c, account, body, "", "")
placeholder
	require.Nil(t, result)
	var failoverErr *UpstreamFailoverError
	require.True(t, errors.As(err, &failoverErr))
	require.Equal(t, http.StatusTooManyRequests, failoverErr.StatusCode)
	require.Equal(t, "45", failoverErr.ResponseHeaders.Get("Retry-After"))
	require.Equal(t, xai.DefaultCLIBaseURL+"/responses", upstream.lastReq.URL.String())
	require.Equal(t, grokChatResponsesEndpoint, GetActualOpenAIUpstreamEndpoint(c))
	require.Equal(t, 1, repo.rateLimitedCalls)
	require.Zero(t, repo.tempUnschedCalls)
	require.WithinDuration(t, before.Add(45*time.Second), repo.lastRateLimitResetAt, time.Second)
	require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestForwardGrokRawChat429PreservesRetryAfter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"model":"grok","messages":[{"role":"user","content":"hi"placeholder],"stream":false,"stop":"done"placeholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, grokChatRawEndpoint, bytes.NewReader(body))
	c.Set("api_key", &APIKey{ID: 7551placeholder)

	account := grokChatBridgeTestAccount(755)
	account.Credentials["expires_at"] = time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339)
	repo := &grokQuotaAccountRepo{mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
		accountsByID: map[int64]*Account{account.ID: accountplaceholder,
placeholderplaceholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header: http.Header{
			"Content-Type": []string{"application/json"placeholder,
			"Retry-After":  []string{"45"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(`{"error":{"message":"rate limited"placeholderplaceholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{
		httpUpstream:      upstream,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		accountRepo:       repo,
placeholder

	result, err := svc.ForwardAsChatCompletions(context.Background(), c, account, body, "", "")

placeholder
	require.Nil(t, result)
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusTooManyRequests, failoverErr.StatusCode)
	require.Equal(t, "45", failoverErr.ResponseHeaders.Get("Retry-After"))
	require.Equal(t, xai.DefaultCLIBaseURL+"/chat/completions", upstream.lastReq.URL.String())
placeholder

func TestForwardGrokRawChatErrorRecordsActualEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"model":"grok","messages":[{"role":"user","content":"hi"placeholder],"stream":false,"stop":"done"placeholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, grokChatRawEndpoint, bytes.NewReader(body))
	c.Set("api_key", &APIKey{ID: 7601placeholder)

	account := grokChatBridgeTestAccount(76)
	repo := &grokQuotaAccountRepo{mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
		accountsByID: map[int64]*Account{account.ID: accountplaceholder,
placeholderplaceholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"bad request"placeholderplaceholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{
		httpUpstream:      upstream,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		accountRepo:       repo,
placeholder

	result, err := svc.ForwardAsChatCompletions(context.Background(), c, account, body, "", "")
placeholder
	require.Nil(t, result)
	require.Equal(t, xai.DefaultCLIBaseURL+"/chat/completions", upstream.lastReq.URL.String())
	require.Equal(t, grokChatRawEndpoint, GetActualOpenAIUpstreamEndpoint(c))
placeholder

func grokChatBridgeTestAccount(id int64) *Account {
placeholder
		ID:          id,
		Name:        "grok-cache-bridge",
		Platform:    PlatformGrok,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
placeholder
			"access_token":  "access-token",
			"refresh_token": "refresh-token",
			"expires_at":    time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339),
			"base_url":      xai.DefaultCLIBaseURL,
	placeholder,
placeholder
placeholder

func grokChatBridgeCompletedResponse(responseID string, cachedTokens int) *http.Response {
	body := strings.Join([]string{
		`data: {"type":"response.output_text.delta","sequence_number":0,"delta":"cached ok"placeholder`,
		"",
		`data: {"type":"response.completed","sequence_number":1,"response":{"id":"` + responseID + `","object":"response","model":"grok-4.5","status":"completed","output":[{"type":"message","id":"msg_1","role":"assistant","status":"completed","content":[{"type":"output_text","text":"cached ok"placeholder]placeholder],"usage":{"input_tokens":9908,"output_tokens":12,"total_tokens":9920,"input_tokens_details":{"cached_tokens":` + strconv.Itoa(cachedTokens) + `placeholderplaceholderplaceholderplaceholder`,
		"",
placeholder, "\n")
	return &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type":                   []string{"text/event-stream"placeholder,
			"Xai-Request-Id":                 []string{responseID + "-request"placeholder,
			"X-Ratelimit-Limit-Requests":     []string{"10"placeholder,
			"X-Ratelimit-Remaining-Requests": []string{"9"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(body)),
placeholder
placeholder
