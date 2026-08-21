package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

type panicOnReadCloser struct{placeholder

func (panicOnReadCloser) Read(_ []byte) (int, error) {
	panic("response body should not be reread")
placeholder

func (panicOnReadCloser) Close() error { return nil placeholder

func TestOpenAIGatewayService_Forward_FailoverReparsesCachedBodyForNextAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		requestModel  string
		firstMapping  map[string]any
		secondMapping map[string]any
		wantFirst     string
		wantSecond    string
placeholder{
		{
			name:          "both accounts have mapping",
			firstMapping:  map[string]any{"alias-model": "base-model-a"placeholder,
			secondMapping: map[string]any{"alias-model": "base-model-b"placeholder,
			wantFirst:     "base-model-a",
			wantSecond:    "base-model-b",
	placeholder,
		{
			name:         "first account has mapping second account has none",
			requestModel: "gpt-5.4-high",
			firstMapping: map[string]any{"gpt-5.4-high": "gpt-5.4"placeholder,
			wantFirst:    "gpt-5.4",
			wantSecond:   "gpt-5.4",
	placeholder,
		{
			name:          "first account has no mapping second account has mapping",
			secondMapping: map[string]any{"alias-model": "base-model-b"placeholder,
			wantFirst:     "alias-model",
			wantSecond:    "base-model-b",
	placeholder,
		{
			name:          "legacy context cache is ignored when mappings differ",
			firstMapping:  map[string]any{"alias-model": "base-model-a"placeholder,
			secondMapping: map[string]any{"alias-model": "base-model-b"placeholder,
			wantFirst:     "base-model-a",
			wantSecond:    "base-model-b",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestModel := tt.requestModel
			if requestModel == "" {
				requestModel = "alias-model"
		placeholder
			body := []byte(`{"model":"` + requestModel + `","stream":false,"instructions":"cache-test","input":"hello"placeholder`)

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			upstream := &httpUpstreamRecorder{responses: []*http.Response{
				{
					StatusCode: http.StatusTooManyRequests,
					Header:     http.Header{"Content-Type": []string{"application/json"placeholder, "x-request-id": []string{"rid-failover-a"placeholderplaceholder,
					Body:       io.NopCloser(strings.NewReader(`{"error":{"type":"rate_limit_error","message":"rate limited"placeholderplaceholder`)),
			placeholder,
				{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"placeholder, "x-request-id": []string{"rid-ok-b"placeholderplaceholder,
					Body:       io.NopCloser(strings.NewReader(`{"id":"resp_123","status":"completed","model":"ok","output":[],"usage":{"input_tokens":1,"output_tokens":1placeholderplaceholder`)),
			placeholder,
		placeholderplaceholder
			svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder

			firstAccount := openAIFailoverCachedBodyTestAccount(1, "account-a", tt.firstMapping)
			secondAccount := openAIFailoverCachedBodyTestAccount(2, "account-b", tt.secondMapping)

			_, err := svc.Forward(context.Background(), c, firstAccount, body)
		placeholder
			var failoverErr *UpstreamFailoverError
			require.True(t, errors.As(err, &failoverErr))
			require.Len(t, upstream.bodies, 1)
			require.Equal(t, tt.wantFirst, gjson.GetBytes(upstream.bodies[0], "model").String())

			c.Set("openai_parsed_request_body", map[string]any{"model": tt.wantFirst, "stream": trueplaceholder)
			result, err := svc.Forward(context.Background(), c, secondAccount, body)
		placeholder
			require.NotNil(t, result)
			require.Len(t, upstream.bodies, 2)
			require.Equal(t, tt.wantSecond, gjson.GetBytes(upstream.bodies[1], "model").String())
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayService_HandleFailoverSideEffects_DoesNotRereadResponseBody(t *testing.T) {
	svc := &OpenAIGatewayService{placeholder
	account := &Account{
		ID:       88,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{placeholder,
		Body:       panicOnReadCloser{placeholder,
placeholder

	require.NotPanics(t, func() {
		svc.handleFailoverSideEffects(context.Background(), resp, account, []byte(`{"error":{"type":"rate_limit_error","message":"rate limited"placeholderplaceholder`))
placeholder)

	require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
	require.True(t, svc.shouldRetryOpenAIOAuth429OnSameAccount(account, http.StatusTooManyRequests, false))
placeholder

func TestGetOpenAIRequestBodyMap_IgnoresLegacyContextCache(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Set("openai_parsed_request_body", map[string]any{"model": "base-model-a", "stream": trueplaceholder)

	got, err := getOpenAIRequestBodyMap(c, []byte(`{"model":"alias-model","stream":falseplaceholder`))
placeholder
	require.Equal(t, "alias-model", got["model"])
	require.Equal(t, false, got["stream"])
placeholder

func openAIFailoverCachedBodyTestAccount(id int64, name string, mapping map[string]any) *Account {
	credentials := map[string]any{"access_token": "oauth-token", "chatgpt_account_id": "chatgpt-account"placeholder
	if mapping != nil {
		credentials["model_mapping"] = mapping
placeholder
placeholder
		ID:             id,
		Name:           name,
		Platform:       PlatformOpenAI,
		Type:           AccountTypeOAuth,
		Concurrency:    1,
		Credentials:    credentials,
		Status:         StatusActive,
		Schedulable:    true,
		RateMultiplier: f64p(1),
placeholder
placeholder
