package service

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// response.incomplete（生成超时/截断）应被识别为可重试的 502 上游错误，触发 failover。
func TestExtractImagesUpstreamError_IncompleteIsRetryable(t *testing.T) {
	body := "data: {\"type\":\"response.created\",\"response\":{\"id\":\"resp_1\"placeholderplaceholder\n\n" +
		"data: {\"type\":\"response.incomplete\",\"response\":{\"id\":\"resp_1\",\"status\":\"incomplete\",\"incomplete_details\":{\"reason\":\"max_output_tokens\"placeholderplaceholderplaceholder\n\n"
	got := extractOpenAIImagesUpstreamError([]byte(body))
	if got == nil {
		t.Fatal("incomplete event should produce an upstream error, got nil")
placeholder
	if got.StatusCode != http.StatusBadGateway {
		t.Fatalf("incomplete(max_output_tokens) should be 502 retryable, got %d", got.StatusCode)
placeholder
	if !IsOpenAIImagesRetryableUpstreamError(got) {
		t.Fatal("incomplete(max_output_tokens) should be retryable for failover")
placeholder
	if got.Code != "response_incomplete" {
		t.Fatalf("unexpected code %q", got.Code)
placeholder
	if !strings.Contains(got.Message, "max_output_tokens") {
		t.Fatalf("message should carry reason, got %q", got.Message)
placeholder
placeholder

// incomplete 因 content_filter → 400，重试无意义，不应触发 failover。
func TestExtractImagesUpstreamError_IncompleteContentFilterNotRetryable(t *testing.T) {
	body := "data: {\"type\":\"response.incomplete\",\"response\":{\"id\":\"r\",\"status\":\"incomplete\",\"incomplete_details\":{\"reason\":\"content_filter\"placeholderplaceholderplaceholder\n\n"
	got := extractOpenAIImagesUpstreamError([]byte(body))
	if got == nil {
		t.Fatal("content_filter incomplete should produce error")
placeholder
	if got.StatusCode != http.StatusBadRequest {
		t.Fatalf("content_filter should be 400 (non-retryable), got %d", got.StatusCode)
placeholder
	if IsOpenAIImagesRetryableUpstreamError(got) {
		t.Fatal("content_filter must NOT be retryable")
placeholder
placeholder

// 旧行为不变：error / response.failed 仍按原逻辑识别。
func TestExtractImagesUpstreamError_ErrorAndFailedUnchanged(t *testing.T) {
	errBody := "data: {\"type\":\"error\",\"error\":{\"type\":\"image_generation_user_error\",\"code\":\"moderation_blocked\",\"message\":\"rejected\"placeholderplaceholder\n\n"
	if got := extractOpenAIImagesUpstreamError([]byte(errBody)); got == nil || got.StatusCode != http.StatusBadRequest {
		t.Fatalf("moderation_blocked should still be 400, got %+v", got)
placeholder
placeholder

// 上游既无图、又无任何可识别事件时，摘要函数应提取诊断信息。
func TestSummarizeNoOutputBody_ExtractsDiagnostics(t *testing.T) {
	body := "data: {\"type\":\"response.created\",\"response\":{\"id\":\"r\"placeholderplaceholder\n\n" +
		"data: {\"type\":\"response.in_progress\",\"response\":{\"id\":\"r\",\"status\":\"in_progress\"placeholderplaceholder\n\n"
	summary := summarizeOpenAIImagesNoOutputBody([]byte(body))
	if !strings.HasPrefix(summary, "no_image_output") {
		t.Fatalf("summary should start with marker, got %q", summary)
placeholder
	if !strings.Contains(summary, "last_event=response.in_progress") {
		t.Fatalf("summary should capture last event type, got %q", summary)
placeholder
	if !strings.Contains(summary, "status=in_progress") {
		t.Fatalf("summary should capture response status, got %q", summary)
placeholder
placeholder

// 摘要应能抓到 incomplete_reason 并对超长 body 截断。
func TestSummarizeNoOutputBody_IncompleteReasonAndTruncation(t *testing.T) {
	long := strings.Repeat("x", 2000)
	body := "data: {\"type\":\"response.incomplete\",\"response\":{\"status\":\"incomplete\",\"incomplete_details\":{\"reason\":\"max_output_tokens\"placeholder,\"junk\":\"" + long + "\"placeholderplaceholder\n\n"
	summary := summarizeOpenAIImagesNoOutputBody([]byte(body))
	if !strings.Contains(summary, "incomplete_reason=max_output_tokens") {
		t.Fatalf("should capture incomplete reason, got %q", summary[:120])
placeholder
	if !strings.Contains(summary, "truncated") {
		t.Fatalf("oversized body should be truncated, len=%d", len(summary))
placeholder
placeholder

// 软失败（上游 completed 但无图，如偶发路由到 mini 模型）应返回可重试的
// UpstreamFailoverError 且优先同账号重试，而非一次性失败。
func TestImagesOAuthNonStreaming_CompletedNoImageTriggersSameAccountRetry(t *testing.T) {
	// 上游 SSE：response.completed 但 output 为空（实测的真实失败形态）。
	upstreamSSE := "event: response.created\n" +
		"data: {\"type\":\"response.created\",\"response\":{\"id\":\"resp_x\",\"status\":\"in_progress\",\"model\":\"gpt-5.4-mini-2026-03-17\",\"output\":[]placeholderplaceholder\n\n" +
		"event: response.completed\n" +
		"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_x\",\"status\":\"completed\",\"model\":\"gpt-5.4-mini-2026-03-17\",\"output\":[],\"tool_usage\":{\"image_gen\":{\"output_tokens\":0placeholderplaceholderplaceholderplaceholder\n\n"

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", nil)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{placeholder,
		Body:       io.NopCloser(strings.NewReader(upstreamSSE)),
placeholder

	svc := &OpenAIGatewayService{placeholder
	_, _, _, err := svc.handleOpenAIImagesOAuthNonStreamingResponse(resp, c, "b64_json", "gpt-image-2")

	if err == nil {
		t.Fatal("completed-but-no-image should return an error")
placeholder
	var failoverErr *UpstreamFailoverError
	if !errors.As(err, &failoverErr) {
		t.Fatalf("expected *UpstreamFailoverError to trigger retry, got %T: %v", err, err)
placeholder
	if failoverErr.StatusCode != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", failoverErr.StatusCode)
placeholder
	if !failoverErr.RetryableOnSameAccount {
		t.Fatal("soft-failure should prefer same-account retry (probabilistic upstream failure)")
placeholder
placeholder

// 内容审核拒绝（模型未出图但输出文字拒绝）应返回 400 content_policy 错误且不重试，
// 而非可重试的 UpstreamFailoverError。
func TestImagesOAuthNonStreaming_ContentRefusalReturns400NoRetry(t *testing.T) {
	// 上游：completed 无图，但模型输出了文字拒绝（内容审核场景）。
	upstreamSSE := "event: response.created\n" +
		"data: {\"type\":\"response.created\",\"response\":{\"id\":\"r\",\"status\":\"in_progress\",\"model\":\"gpt-5.4-mini\",\"output\":[]placeholderplaceholder\n\n" +
		"event: response.output_text.delta\n" +
		"data: {\"type\":\"response.output_text.delta\",\"delta\":\"抱歉，这个请求因涉及违规内容被安全系统判定为不适合生成。\"placeholder\n\n" +
		"event: response.completed\n" +
		"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"r\",\"status\":\"completed\",\"model\":\"gpt-5.4-mini\",\"output\":[{\"type\":\"message\",\"content\":[{\"type\":\"output_text\",\"text\":\"抱歉，这个请求因涉及违规内容被安全系统判定为不适合生成。\"placeholder]placeholder],\"tool_usage\":{\"image_gen\":{\"output_tokens\":0placeholderplaceholderplaceholderplaceholder\n\n"

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", nil)
	resp := &http.Response{StatusCode: http.StatusOK, Header: http.Header{placeholder, Body: io.NopCloser(strings.NewReader(upstreamSSE))placeholder

	svc := &OpenAIGatewayService{placeholder
	_, _, _, err := svc.handleOpenAIImagesOAuthNonStreamingResponse(resp, c, "b64_json", "gpt-image-2")

	if err == nil {
		t.Fatal("content refusal should return an error")
placeholder
	// 应是不可重试的内容策略错误（400），而非 UpstreamFailoverError。
	var failoverErr *UpstreamFailoverError
	if errors.As(err, &failoverErr) {
		t.Fatalf("content refusal must NOT be a retryable failover error, got %v", failoverErr)
placeholder
	var imgErr *OpenAIImagesUpstreamError
	if !errors.As(err, &imgErr) {
		t.Fatalf("expected *OpenAIImagesUpstreamError, got %T: %v", err, err)
placeholder
	if imgErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("content refusal should be 400, got %d", imgErr.StatusCode)
placeholder
	if !strings.Contains(imgErr.Message, "安全系统") && !strings.Contains(imgErr.Message, "违规") {
		t.Fatalf("refusal message should carry model's reason, got %q", imgErr.Message)
placeholder
placeholder

// extractOpenAIImagesModelRefusal：真空响应（无文字）返回空串。
func TestExtractModelRefusal_EmptyWhenNoText(t *testing.T) {
	body := "data: {\"type\":\"response.completed\",\"response\":{\"output\":[],\"tool_usage\":{\"image_gen\":{\"output_tokens\":0placeholderplaceholderplaceholderplaceholder\n\n"
	if r := extractOpenAIImagesModelRefusal([]byte(body)); r != "" {
		t.Fatalf("empty response should yield no refusal, got %q", r)
placeholder
placeholder
