package service

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type gatewayForwardErrorPolicyRepoStub struct {
	AccountRepository
	tempCalls           int
	modelRateLimitCalls []gatewayForwardModelRateLimitCall
placeholder

type gatewayForwardModelRateLimitCall struct {
	accountID int64
	scope     string
placeholder

func (r *gatewayForwardErrorPolicyRepoStub) SetTempUnschedulable(context.Context, int64, time.Time, string) error {
	r.tempCalls++
	return nil
placeholder

func (r *gatewayForwardErrorPolicyRepoStub) SetModelRateLimit(_ context.Context, id int64, scope string, _ time.Time, _ ...string) error {
	r.modelRateLimitCalls = append(r.modelRateLimitCalls, gatewayForwardModelRateLimitCall{
		accountID: id,
		scope:     scope,
placeholder)
	return nil
placeholder

// 本文件覆盖 issue #5148：流式转发中途出错（缺失 terminal 事件、读错误等）时，
// 已观测到的上游 usage 不得随错误一起被丢弃，Forward 必须把部分结果与错误一同
// 返回，供 handler 照常提交 usage 记录。

func newForwardPartialUsageServiceForTest(upstream *anthropicHTTPUpstreamRecorder) *GatewayService {
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			MaxLineSize: defaultMaxLineSize,
	placeholder,
placeholder
	return &GatewayService{
		cfg:                  cfg,
		responseHeaderFilter: compileResponseHeaderFilter(cfg),
		httpUpstream:         upstream,
		rateLimitService:     &RateLimitService{placeholder,
		deferredService:      &DeferredService{placeholder,
placeholder
placeholder

func newAnthropicOAuthAccountForPartialUsageTest() *Account {
placeholder
		ID:          501,
		Name:        "anthropic-oauth-partial-usage",
		Platform:    PlatformAnthropic,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token": "oauth-token",
	placeholder,
		Status:      StatusActive,
		Schedulable: true,
placeholder
placeholder

func TestGatewayService_Forward_StreamMissingTerminalPreservesPartialUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	body := []byte(`{"model":"claude-3-5-sonnet-latest","stream":true,"messages":[{"role":"user","content":[{"type":"text","text":"hello"placeholder]placeholder]placeholder`)
	parsed, err := ParseGatewayRequest(NewRequestBodyRef(body), PlatformAnthropic)
placeholder

	// newapi 等聚合上游的典型失败形态：message_start/message_delta 携带 usage，
	// 但流在 message_stop 前直接结束。
	upstreamSSE := strings.Join([]string{
		`event: message_start`,
		`data: {"type":"message_start","message":{"id":"msg_1","type":"message","role":"assistant","model":"claude-3-5-sonnet-latest","content":[],"usage":{"input_tokens":11,"cache_read_input_tokens":7placeholderplaceholderplaceholder`,
		"",
		`event: content_block_delta`,
		`data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"hi"placeholderplaceholder`,
		"",
		`event: message_delta`,
		`data: {"type":"message_delta","delta":{"stop_reason":nullplaceholder,"usage":{"output_tokens":5placeholderplaceholder`,
		"",
		"",
placeholder, "\n")
	upstream := &anthropicHTTPUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"text/event-stream"placeholder,
			"X-Request-Id": []string{"rid-partial"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(upstreamSSE)),
placeholderplaceholder
	svc := newForwardPartialUsageServiceForTest(upstream)
	account := newAnthropicOAuthAccountForPartialUsageTest()

	result, err := svc.Forward(context.Background(), c, account, parsed)
placeholder
	require.Contains(t, err.Error(), "missing terminal event")
	require.NotNil(t, result, "流中断但已观测到 usage 时必须返回部分结果用于计费")
	require.True(t, result.Stream)
	require.Equal(t, 11, result.Usage.InputTokens)
	require.Equal(t, 7, result.Usage.CacheReadInputTokens)
	require.Equal(t, 5, result.Usage.OutputTokens)
	require.Equal(t, "rid-partial", result.RequestID)
	require.NotNil(t, result.FirstTokenMs)
placeholder

func TestGatewayService_Forward_StreamReadErrorAfterOutputPreservesPartialUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	body := []byte(`{"model":"claude-3-5-sonnet-latest","stream":true,"messages":[{"role":"user","content":[{"type":"text","text":"hello"placeholder]placeholder]placeholder`)
	parsed, err := ParseGatewayRequest(NewRequestBodyRef(body), PlatformAnthropic)
placeholder

	// message_start 已写出（含 usage），随后上游连接异常中断。
	upstream := &anthropicHTTPUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body: &streamReadCloser{
			payload: []byte("data: {\"type\":\"message_start\",\"message\":{\"usage\":{\"input_tokens\":9,\"cache_creation_input_tokens\":4placeholderplaceholderplaceholder\n\n"),
			err:     io.ErrUnexpectedEOF,
	placeholder,
placeholderplaceholder
	svc := newForwardPartialUsageServiceForTest(upstream)
	account := newAnthropicOAuthAccountForPartialUsageTest()

	result, err := svc.Forward(context.Background(), c, account, parsed)
placeholder
	require.Contains(t, err.Error(), "stream read error")
	require.NotNil(t, result, "已写出内容后的读错误必须保留部分 usage")
	require.Equal(t, 9, result.Usage.InputTokens)
	require.Equal(t, 4, result.Usage.CacheCreationInputTokens)
placeholder

func TestGatewayService_Forward_StreamErrorWithoutUsageReturnsNilResult(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	body := []byte(`{"model":"claude-3-5-sonnet-latest","stream":true,"messages":[{"role":"user","content":[{"type":"text","text":"hello"placeholder]placeholder]placeholder`)
	parsed, err := ParseGatewayRequest(NewRequestBodyRef(body), PlatformAnthropic)
placeholder

	// 只有 ping、没有任何 usage 的流中断：不应产生零 usage 的幽灵账单记录。
	upstream := &anthropicHTTPUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader("event: ping\ndata: {\"type\": \"ping\"placeholder\n\n")),
placeholderplaceholder
	svc := newForwardPartialUsageServiceForTest(upstream)
	account := newAnthropicOAuthAccountForPartialUsageTest()

	result, err := svc.Forward(context.Background(), c, account, parsed)
placeholder
	require.Contains(t, err.Error(), "missing terminal event")
	require.Nil(t, result, "无已观测 usage 时不应返回部分结果")
placeholder

func TestGatewayService_Forward_FailoverErrorKeepsNilResult(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	body := []byte(`{"model":"claude-3-5-sonnet-latest","stream":true,"messages":[{"role":"user","content":[{"type":"text","text":"hello"placeholder]placeholder]placeholder`)
	parsed, err := ParseGatewayRequest(NewRequestBodyRef(body), PlatformAnthropic)
placeholder

	// 未向客户端写出任何字节前的读错误会包成 UpstreamFailoverError 走换号重试。
	// 该路径必须保持 result=nil：failover 成功后按成功请求计费，双份结果会重复计费。
	upstream := &anthropicHTTPUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body: &streamReadCloser{
			err: errors.New("connection reset by peer"),
	placeholder,
placeholderplaceholder
	svc := newForwardPartialUsageServiceForTest(upstream)
	account := newAnthropicOAuthAccountForPartialUsageTest()

	result, err := svc.Forward(context.Background(), c, account, parsed)
placeholder
	var failoverErr *UpstreamFailoverError
	require.True(t, errors.As(err, &failoverErr))
	require.Nil(t, result, "failover 错误必须保持 result=nil，防止重试成功后双重计费")
placeholder

func TestGatewayService_Forward_PreOutputSSEOverloadedErrorUsesSemantic529(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	body := []byte(`{"model":"claude-3-5-sonnet-latest","stream":true,"messages":[{"role":"user","content":"hello"placeholder]placeholder`)
	parsed, err := ParseGatewayRequest(NewRequestBodyRef(body), PlatformAnthropic)
placeholder

	const errorJSON = `{"type":"error","error":{"details":null,"type":"overloaded_error","message":"Overloaded"placeholder,"request_id":"req_01"placeholder`
	upstream := &anthropicHTTPUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader("event: error\ndata: " + errorJSON + "\n\n")),
placeholderplaceholder
	repo := &gatewayForwardErrorPolicyRepoStub{placeholder
	cfg := &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder
	svc := &GatewayService{
		cfg:                  cfg,
		responseHeaderFilter: compileResponseHeaderFilter(cfg),
		httpUpstream:         upstream,
		rateLimitService:     NewRateLimitService(repo, nil, cfg, nil, nil),
		deferredService:      &DeferredService{placeholder,
placeholder
	account := newAnthropicOAuthAccountForPartialUsageTest()
	account.Credentials["temp_unschedulable_enabled"] = true
	account.Credentials["temp_unschedulable_rules"] = []any{map[string]any{
		"error_code":       float64(529),
		"keywords":         []any{"Overloaded"placeholder,
		"duration_minutes": float64(10),
placeholderplaceholder

	result, err := svc.Forward(context.Background(), c, account, parsed)
placeholder
	require.Nil(t, result)

	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, 529, failoverErr.StatusCode)
	require.JSONEq(t, errorJSON, string(failoverErr.ResponseBody))
	require.Len(t, repo.modelRateLimitCalls, 1, "synthetic 529 must participate in temp-unschedulable rules")
	require.Equal(t, account.ID, repo.modelRateLimitCalls[0].accountID)
	require.Equal(t, parsed.Model, repo.modelRateLimitCalls[0].scope)
	require.Empty(t, rec.Body.String(), "pre-output overload must remain eligible for account failover")
placeholder

func TestGatewayService_Forward_PostOutputSSEOverloadedErrorKeepsExistingStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	body := []byte(`{"model":"claude-3-5-sonnet-latest","stream":true,"messages":[{"role":"user","content":"hello"placeholder]placeholder`)
	parsed, err := ParseGatewayRequest(NewRequestBodyRef(body), PlatformAnthropic)
placeholder

	const errorJSON = `{"type":"error","error":{"type":"overloaded_error","message":"Overloaded"placeholderplaceholder`
	fixture := "event: message_start\ndata: {\"type\":\"message_start\",\"message\":{\"usage\":{\"input_tokens\":1placeholderplaceholderplaceholder\n\n" +
		"event: error\ndata: " + errorJSON + "\n\n"
	upstream := &anthropicHTTPUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(fixture)),
placeholderplaceholder
	repo := &gatewayForwardErrorPolicyRepoStub{placeholder
	cfg := &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder
	svc := &GatewayService{
		cfg:                  cfg,
		responseHeaderFilter: compileResponseHeaderFilter(cfg),
		httpUpstream:         upstream,
		rateLimitService:     NewRateLimitService(repo, nil, cfg, nil, nil),
		deferredService:      &DeferredService{placeholder,
placeholder

	result, err := svc.Forward(context.Background(), c, newAnthropicOAuthAccountForPartialUsageTest(), parsed)
placeholder
	require.Nil(t, result)

	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusForbidden, failoverErr.StatusCode)
	require.JSONEq(t, errorJSON, string(failoverErr.ResponseBody))
	require.Zero(t, repo.tempCalls)
	require.Contains(t, rec.Body.String(), "message_start")
placeholder

func TestGatewayService_AnthropicAPIKeyPassthrough_ForwardStreamMissingTerminalPreservesPartialUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	body := []byte(`{"model":"claude-3-7-sonnet-20250219","stream":true,"messages":[{"role":"user","content":[{"type":"text","text":"hello"placeholder]placeholder]placeholder`)
	parsed := &ParsedRequest{
		Body:   NewRequestBodyRef(body),
		Model:  "claude-3-7-sonnet-20250219",
		Stream: true,
placeholder

	upstreamSSE := strings.Join([]string{
		`data: {"type":"message_start","message":{"usage":{"input_tokens":9,"cache_read_input_tokens":2placeholderplaceholderplaceholder`,
		"",
		`data: {"type":"message_delta","usage":{"output_tokens":3placeholderplaceholder`,
		"",
placeholder, "\n")
	upstream := &anthropicHTTPUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"text/event-stream"placeholder,
			"X-Request-Id": []string{"rid-pass-partial"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(upstreamSSE)),
placeholderplaceholder
	svc := newForwardPartialUsageServiceForTest(upstream)
	account := newAnthropicAPIKeyAccountForTest()

	result, err := svc.Forward(context.Background(), c, account, parsed)
placeholder
	require.Contains(t, err.Error(), "missing terminal event")
	require.NotNil(t, result, "透传流中断但已观测到 usage 时必须返回部分结果用于计费")
	require.True(t, result.Stream)
	require.Equal(t, 9, result.Usage.InputTokens)
	require.Equal(t, 2, result.Usage.CacheReadInputTokens)
	require.Equal(t, 3, result.Usage.OutputTokens)
	require.Equal(t, "claude-3-7-sonnet-20250219", result.Model)
placeholder
