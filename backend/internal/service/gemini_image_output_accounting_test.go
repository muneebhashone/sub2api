//go:build unit

package service

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

const geminiTestPNG = "iVBORw0KGgoAAAANSUhEUg=="

func newGeminiImageTestContext(t *testing.T) *gin.Context {
placeholder
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost,
		"/v1beta/models/nana-banana-2:generateContent", strings.NewReader("{placeholder"))
	return c
placeholder

func geminiImageResponse(parts string) string {
	return `{"candidates":[{"content":{"role":"model","parts":[` + parts + `]placeholder,"finishReason":"STOP"placeholder]placeholder`
placeholder

func TestCountGeminiInlineImageOutputs(t *testing.T) {
	cases := []struct {
		name    string
		payload string
		want    int
placeholder{
		{
			name:    "camelCase inlineData",
			payload: geminiImageResponse(`{"inlineData":{"mimeType":"image/png","data":"` + geminiTestPNG + `"placeholderplaceholder`),
			want:    1,
	placeholder,
		{
			// 官方 SDK 与部分中转把字段回成 snake_case。
			name:    "snake_case inline_data",
			payload: geminiImageResponse(`{"inline_data":{"mime_type":"image/png","data":"` + geminiTestPNG + `"placeholderplaceholder`),
			want:    1,
	placeholder,
		{
			name: "text and image mixed",
			payload: geminiImageResponse(`{"text":"here you go"placeholder,` +
				`{"inlineData":{"mimeType":"image/jpeg","data":"` + geminiTestPNG + `"placeholderplaceholder`),
			want: 1,
	placeholder,
		{
			name: "multiple images",
			payload: geminiImageResponse(`{"inlineData":{"mimeType":"image/png","data":"` + geminiTestPNG + `"placeholderplaceholder,` +
				`{"inlineData":{"mimeType":"image/webp","data":"` + geminiTestPNG + `"placeholderplaceholder`),
			want: 2,
	placeholder,
		{
			name:    "uppercase mime type",
			payload: geminiImageResponse(`{"inlineData":{"mimeType":"IMAGE/PNG","data":"` + geminiTestPNG + `"placeholderplaceholder`),
			want:    1,
	placeholder,
		{
			name:    "text only",
			payload: geminiImageResponse(`{"text":"no image here"placeholder`),
			want:    0,
	placeholder,
		{
			// 非图片的内联附件（例如音频）不能按图片计费。
			name:    "non image mime type",
			payload: geminiImageResponse(`{"inlineData":{"mimeType":"audio/mpeg","data":"` + geminiTestPNG + `"placeholderplaceholder`),
			want:    0,
	placeholder,
		{
			name:    "empty data is not billable",
			payload: geminiImageResponse(`{"inlineData":{"mimeType":"image/png","data":""placeholderplaceholder`),
			want:    0,
	placeholder,
		{name: "empty payload", payload: "", want: 0placeholder,
		{name: "invalid json", payload: "not-json", want: 0placeholder,
		{name: "error response", payload: `{"error":{"code":429,"message":"quota"placeholderplaceholder`, want: 0placeholder,
placeholder

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, countGeminiInlineImageOutputs([]byte(tc.payload)))
	placeholder)
placeholder
placeholder

// 累积式 SSE 会把同一张图在后续 chunk 里整段重发，逐 chunk 累加会重复计费。
// 计数器取单个 payload 内的最大值，正是为了挡住这一点。
func TestObserveGeminiImageOutputs_CumulativeChunksDoNotDoubleCount(t *testing.T) {
	c := newGeminiImageTestContext(t)
	beginGeminiImageOutputObservation(c)

	oneImage := geminiImageResponse(`{"inlineData":{"mimeType":"image/png","data":"` + geminiTestPNG + `"placeholderplaceholder`)
	for range 4 {
		observeGeminiImageOutputs(c, []byte(oneImage))
placeholder

	require.Equal(t, 1, observedGeminiImageOutputs(c))
placeholder

func TestObserveGeminiImageOutputs_KeepsLargestChunk(t *testing.T) {
	c := newGeminiImageTestContext(t)
	beginGeminiImageOutputObservation(c)

	observeGeminiImageOutputs(c, []byte(geminiImageResponse(`{"text":"working"placeholder`)))
	observeGeminiImageOutputs(c, []byte(geminiImageResponse(
		`{"inlineData":{"mimeType":"image/png","data":"`+geminiTestPNG+`"placeholderplaceholder,`+
			`{"inlineData":{"mimeType":"image/png","data":"`+geminiTestPNG+`"placeholderplaceholder`)))
	// 收尾 chunk 只带 usageMetadata，不能把已数到的张数抹掉。
	observeGeminiImageOutputs(c, []byte(`{"usageMetadata":{"promptTokenCount":9placeholderplaceholder`))

	require.Equal(t, 2, observedGeminiImageOutputs(c))
placeholder

// failover 会拿同一个 gin.Context 重跑 Forward，计数器必须按次重置，
// 否则失败账号已经回吐的图会被叠加到成功账号的账单上。
func TestBeginGeminiImageOutputObservation_ResetsPerForward(t *testing.T) {
	c := newGeminiImageTestContext(t)
	oneImage := []byte(geminiImageResponse(`{"inlineData":{"mimeType":"image/png","data":"` + geminiTestPNG + `"placeholderplaceholder`))

	beginGeminiImageOutputObservation(c)
	observeGeminiImageOutputs(c, oneImage)
	require.Equal(t, 1, observedGeminiImageOutputs(c))

	beginGeminiImageOutputObservation(c)
	require.Equal(t, 0, observedGeminiImageOutputs(c))
	observeGeminiImageOutputs(c, oneImage)
	require.Equal(t, 1, observedGeminiImageOutputs(c))
placeholder

// issue #5358：自定义模型名（客户端名与上游映射名都不在白名单里）走 Gemini 原生
// generateContent 生图，改动前 ImageCount 恒为 0，calculateRecordUsageCost 的按次
// 计费分支整条不触发，四次生图全部记 $0。
func TestResolveGeminiImageCount(t *testing.T) {
	oneImage := []byte(geminiImageResponse(`{"inlineData":{"mimeType":"image/png","data":"` + geminiTestPNG + `"placeholderplaceholder`))
	textOnly := []byte(geminiImageResponse(`{"text":"hello"placeholder`))

	t.Run("custom model name bills by observed images", func(t *testing.T) {
		c := newGeminiImageTestContext(t)
		beginGeminiImageOutputObservation(c)
		observeGeminiImageOutputs(c, oneImage)

		require.False(t, isImageGenerationModel("nana-banana-2"), "前置条件：白名单判不出自定义名")
		require.Equal(t, 1, resolveGeminiImageCount(c, "nana-banana-2", "nana-banana-2"))
placeholder)

	t.Run("falls back to requested model name", func(t *testing.T) {
		c := newGeminiImageTestContext(t)
		beginGeminiImageOutputObservation(c)
		observeGeminiImageOutputs(c, textOnly)

		require.Equal(t, 1, resolveGeminiImageCount(c, "gemini-3-pro-image-preview", "gemini-3-pro-image-preview"))
placeholder)

	t.Run("falls back to mapped upstream model name", func(t *testing.T) {
		c := newGeminiImageTestContext(t)
		beginGeminiImageOutputObservation(c)
		observeGeminiImageOutputs(c, textOnly)

		require.Equal(t, 1, resolveGeminiImageCount(c, "my-image-alias", "gemini-2.5-flash-image"))
placeholder)

	t.Run("text model stays unbilled", func(t *testing.T) {
		c := newGeminiImageTestContext(t)
		beginGeminiImageOutputObservation(c)
		observeGeminiImageOutputs(c, textOnly)

		require.Equal(t, 0, resolveGeminiImageCount(c, "gemini-2.5-pro", "gemini-2.5-pro"))
placeholder)

	t.Run("no counter on context degrades to name heuristic", func(t *testing.T) {
		c := newGeminiImageTestContext(t)
		require.Equal(t, 0, resolveGeminiImageCount(c, "nana-banana-2", "nana-banana-2"))
		require.Equal(t, 1, resolveGeminiImageCount(c, "gemini-3-pro-image", "gemini-3-pro-image"))
placeholder)
placeholder

// 端到端守住接线：/v1beta/models/{modelplaceholder:generateContent 的非流式响应体
// 必须真的喂进计数器，否则上面的单测全绿而线上依然记 $0。
func TestHandleNativeNonStreamingResponse_FeedsImageCounter(t *testing.T) {
	c := newGeminiImageTestContext(t)
	beginGeminiImageOutputObservation(c)

	body := geminiImageResponse(`{"inlineData":{"mimeType":"image/png","data":"` + geminiTestPNG + `"placeholderplaceholder`)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(body)),
placeholder

	svc := &GeminiMessagesCompatService{placeholder
	usage, err := svc.handleNativeNonStreamingResponse(c, resp, false)
placeholder
	require.NotNil(t, usage)

	require.Equal(t, 1, observedGeminiImageOutputs(c))
	require.Equal(t, 1, resolveGeminiImageCount(c, "nana-banana-2", "nana-banana-2"))
placeholder
