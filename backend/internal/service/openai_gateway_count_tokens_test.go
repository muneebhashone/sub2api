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

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	body := []byte(`{"model":"claude-opus-4-1","messages":[{"role":"user","content":"hello"placeholder]placeholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages/count_tokens", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("User-Agent", "Claude-Code/1.0")

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusUnauthorized,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"error":{"type":"invalid_request_error","message":"unauthorized"placeholderplaceholder`)),
placeholderplaceholder

	svc := &OpenAIGatewayService{
		cfg:          &config.Config{placeholder,
		httpUpstream: upstream,
placeholder
	account := &Account{
		ID:          202,
		Name:        "openai-oauth",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token": "oauth-token",
	placeholder,
		Status:      StatusActive,
		Schedulable: true,
placeholder

	err := svc.ForwardCountTokensAsAnthropic(context.Background(), c, account, body, "gpt-5.4")
placeholder
	require.Equal(t, http.StatusNotFound, rec.Code)
	require.Contains(t, rec.Body.String(), "Token counting is not supported for this OpenAI account type")
	require.NotNil(t, upstream.lastReq)
	require.Equal(t, "https://api.openai.com/v1/responses/input_tokens", upstream.lastReq.URL.String())
	require.Equal(t, "Bearer oauth-token", upstream.lastReq.Header.Get("authorization"))
	require.Empty(t, upstream.lastReq.Header.Get("Chatgpt-Account-Id"))
placeholder
