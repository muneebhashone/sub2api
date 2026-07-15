package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestOpenAIRequestBodyLimitFailover_HTTP413SwitchesAccountsBeforeWrite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	requestBody := []byte(`{"model":"gpt-5.2","stream":false,"input":"hello"placeholder`)

	for _, passthrough := range []bool{false, trueplaceholder {
		name := "native_responses"
		if passthrough {
			name = "api_key_passthrough"
	placeholder
		t.Run(name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(nil))

			const upstreamBody = `{"error":{"message":"request body exceeds this account's 16MB proxy limit; secret=must-not-leak","type":"invalid_request_error"placeholderplaceholder`
			body := &passthroughCloseTrackingReadCloser{Reader: strings.NewReader(upstreamBody)placeholder
			upstream := &httpUpstreamRecorder{resp: &http.Response{
				StatusCode: http.StatusRequestEntityTooLarge,
				Header: http.Header{
					"Content-Type": []string{"application/json"placeholder,
					"X-Request-Id": []string{"rid-body-limit"placeholder,
			placeholder,
				Body: body,
		placeholderplaceholder
			svc := &OpenAIGatewayService{
				cfg:          &config.Config{Gateway: config.GatewayConfig{ForceCodexCLI: falseplaceholderplaceholder,
				httpUpstream: upstream,
		placeholder
			account := &Account{
				ID:          161,
				Name:        name,
				Platform:    PlatformOpenAI,
				Type:        AccountTypeAPIKey,
				Concurrency: 1,
		placeholder
					"api_key":   "sk-test",
					"base_url":  "https://api.example.test",
					"pool_mode": true,
					"pool_mode_retry_status_codes": []any{
						float64(http.StatusRequestEntityTooLarge),
				placeholder,
			placeholder,
				Extra: map[string]any{
					"openai_passthrough":         passthrough,
					"openai_responses_supported": true,
			placeholder,
				Status:      StatusActive,
				Schedulable: true,
		placeholder

			result, err := svc.Forward(context.Background(), c, account, requestBody)

			require.Nil(t, result)
			var failoverErr *UpstreamFailoverError
			require.ErrorAs(t, err, &failoverErr)
			require.Equal(t, http.StatusRequestEntityTooLarge, failoverErr.StatusCode)
			require.Equal(t, GatewayFailureScopeAccount, failoverErr.Scope)
			require.Equal(t, GatewayFailureReason("openai_request_body_too_large"), failoverErr.Reason)
			require.Equal(t, NextAccountRetry, failoverErr.NextAccountAction)
			require.Equal(t, http.StatusRequestEntityTooLarge, failoverErr.ClientStatusCode)
			require.Equal(t, "Request payload is too large", failoverErr.ClientMessage)
			require.False(t, failoverErr.RetryableOnSameAccount, "a body limit requires another account, not another attempt on the same account")
			require.False(t, c.Writer.Written(), "account failover must happen before downstream output is committed")
			require.Empty(t, rec.Body.String())
			require.True(t, body.closed)
			if passthrough {
				require.Equal(t, requestBody, upstream.lastBody)
		placeholder else {
				require.Equal(t, "gpt-5.2", gjson.GetBytes(upstream.lastBody, "model").String())
				require.Equal(t, "hello", gjson.GetBytes(upstream.lastBody, "input").String())
		placeholder
	placeholder)
placeholder
placeholder

func TestOpenAIRequestBodyLimitFailover_ContextWindow413DoesNotSwitchAccounts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	requestBody := []byte(`{"model":"gpt-5.2","stream":false,"input":"hello"placeholder`)

	for _, passthrough := range []bool{false, trueplaceholder {
		t.Run(fmt.Sprintf("passthrough_%t", passthrough), func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(nil))

			const upstreamBody = `{"error":{"message":"Your input exceeds the context window of this model. Please adjust your input and try again.","type":"invalid_request_error"placeholderplaceholder`
			body := &passthroughCloseTrackingReadCloser{Reader: strings.NewReader(upstreamBody)placeholder
			svc := &OpenAIGatewayService{
				cfg: &config.Config{Gateway: config.GatewayConfig{ForceCodexCLI: falseplaceholderplaceholder,
				httpUpstream: &httpUpstreamRecorder{resp: &http.Response{
					StatusCode: http.StatusRequestEntityTooLarge,
					Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
					Body:       body,
		placeholder
		placeholder
			account := &Account{
				ID: 162, Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Concurrency: 1,
		placeholder"api_key": "sk-test", "base_url": "https://api.example.test"placeholder,
				Extra: map[string]any{
					"openai_passthrough":         passthrough,
					"openai_responses_supported": true,
			placeholder,
				Status: StatusActive, Schedulable: true,
		placeholder

			result, err := svc.Forward(context.Background(), c, account, requestBody)

			require.Nil(t, result)
		placeholder
			var failoverErr *UpstreamFailoverError
			require.False(t, errors.As(err, &failoverErr), "context-window failures are deterministic request errors")
			require.True(t, c.Writer.Written())
			require.Contains(t, rec.Body.String(), "exceeds the context window")
			require.True(t, body.closed)
	placeholder)
placeholder
placeholder
