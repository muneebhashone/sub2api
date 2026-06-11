//go:build unit

// Phase-3 TASK-001 前置特征化测试（T2 补充）：ForwardGemini 的 action 参数契约。
//
// handler :444 调用点（以及未来 gemini adapter）以硬编码 "generateContent"
// 传入 action；本文件锁定该参数的全部行为面，使 adapter 改写 action 字面量
// 时测试必红：
//   - "generateContent" / "streamGenerateContent" 被接受，且上游 action 恒为
//     streamGenerateContent（客户端流式与否由独立的 stream 参数决定）；
//   - "countTokens" 短路返回 {"totalTokens":0placeholder，不触上游；
//   - 其他 action → 404 "Unsupported action"，不触上游。
//
// 夹具复用本包既有 p3 系列/passChar 系列（antigravitySettingRepoStub、
// queuedHTTPUpstreamStub、passCharAntigravityService/Account）。
package service

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// p3CharGeminiNativeBody 返回 Gemini 原生请求体。
func p3CharGeminiNativeBody(t *testing.T) []byte {
placeholder
	return []byte(`{"contents":[{"role":"user","parts":[{"text":"action contract probe"placeholder]placeholder]placeholder`)
placeholder

// p3CharGeminiSSESuccess 返回 v1internal 包裹的上游流式成功响应。
func p3CharGeminiSSESuccess() *http.Response {
	sse := "data: {\"response\":{\"candidates\":[{\"content\":{\"role\":\"model\",\"parts\":[{\"text\":\"ok\"placeholder]placeholder,\"finishReason\":\"STOP\"placeholder],\"usageMetadata\":{\"promptTokenCount\":4,\"candidatesTokenCount\":2placeholderplaceholderplaceholder\n\n"
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(bytes.NewReader([]byte(sse))),
placeholder
placeholder

// TestGatewayCharacterization_ForwardGeminiActionContract 锁定 action 参数契约。
func TestGatewayCharacterization_ForwardGeminiActionContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	newCtx := func(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1beta/models/gemini-2.5-flash:generateContent", bytes.NewReader(p3CharGeminiNativeBody(t)))
		return c, rec
placeholder

	t.Run("generateContent被接受_上游action恒为streamGenerateContent", func(t *testing.T) {
		upstream := &queuedHTTPUpstreamStub{responses: []*http.Response{p3CharGeminiSSESuccess()placeholderplaceholder
		svc := passCharAntigravityService(upstream)
		c, _ := newCtx(t)
		account := passCharAntigravityAccount(9701, map[string]any{"gemini-2.5-flash": "gemini-3-pro-high"placeholder)

		result, err := svc.ForwardGemini(context.Background(), c, account, "gemini-2.5-flash", "generateContent", false, p3CharGeminiNativeBody(t), false)
	placeholder
		require.NotNil(t, result)
		require.Equal(t, 1, upstream.callCount)
placeholder)

	t.Run("streamGenerateContent同样被接受_与generateContent上游行为一致", func(t *testing.T) {
		upstream := &queuedHTTPUpstreamStub{responses: []*http.Response{p3CharGeminiSSESuccess()placeholderplaceholder
		svc := passCharAntigravityService(upstream)
		c, _ := newCtx(t)
		account := passCharAntigravityAccount(9702, map[string]any{"gemini-2.5-flash": "gemini-3-pro-high"placeholder)

		result, err := svc.ForwardGemini(context.Background(), c, account, "gemini-2.5-flash", "streamGenerateContent", false, p3CharGeminiNativeBody(t), false)
	placeholder
		require.NotNil(t, result)
		require.Equal(t, 1, upstream.callCount)
placeholder)

	t.Run("countTokens短路返回零值不触上游", func(t *testing.T) {
		upstream := &queuedHTTPUpstreamStub{placeholder // 任何上游调用都会失败
		svc := passCharAntigravityService(upstream)
		c, rec := newCtx(t)
		account := passCharAntigravityAccount(9703, map[string]any{"gemini-2.5-flash": "gemini-3-pro-high"placeholder)

		result, err := svc.ForwardGemini(context.Background(), c, account, "gemini-2.5-flash", "countTokens", false, p3CharGeminiNativeBody(t), false)
	placeholder
		require.NotNil(t, result)
		require.Equal(t, 0, upstream.callCount, "countTokens 不触上游")
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"totalTokens":0placeholder`, rec.Body.String())
placeholder)

	t.Run("未知action返回404不触上游", func(t *testing.T) {
		upstream := &queuedHTTPUpstreamStub{placeholder
		svc := passCharAntigravityService(upstream)
		c, rec := newCtx(t)
		account := passCharAntigravityAccount(9704, map[string]any{"gemini-2.5-flash": "gemini-3-pro-high"placeholder)

		result, err := svc.ForwardGemini(context.Background(), c, account, "gemini-2.5-flash", "bogusAction", false, p3CharGeminiNativeBody(t), false)
	placeholder
		require.Nil(t, result)
		require.Equal(t, 0, upstream.callCount)
		require.Equal(t, http.StatusNotFound, rec.Code)
		require.Contains(t, rec.Body.String(), "Unsupported action")
placeholder)
placeholder
