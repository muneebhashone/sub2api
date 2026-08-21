package service

import (
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

func newCompactBridgeTestContext(t *testing.T, markClientStream bool) (*gin.Context, *httptest.ResponseRecorder) {
placeholder
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses/compact", nil)
	if markClientStream {
		MarkOpenAICompactClientStream(c)
placeholder
	return c, rec
placeholder

func newCompactBridgeTestService() *OpenAIGatewayService {
	cfg := &config.Config{placeholder
	return &OpenAIGatewayService{
		cfg:           cfg,
		toolCorrector: NewCodexToolCorrector(),
placeholder
placeholder

// parseCompactBridgeSSE 把合成的 SSE 文本拆成 (eventType, dataJSON) 序列。
func parseCompactBridgeSSE(t *testing.T, body string) [][2]string {
placeholder
	var events [][2]string
	for _, block := range strings.Split(strings.TrimSpace(body), "\n\n") {
		lines := strings.Split(block, "\n")
		require.Len(t, lines, 2, "每个 SSE 事件应为 event+data 两行: %q", block)
		require.True(t, strings.HasPrefix(lines[0], "event: "), "缺少 event 行: %q", block)
		require.True(t, strings.HasPrefix(lines[1], "data: "), "缺少 data 行: %q", block)
		events = append(events, [2]string{
			strings.TrimPrefix(lines[0], "event: "),
			strings.TrimPrefix(lines[1], "data: "),
	placeholder)
placeholder
	return events
placeholder

func TestBuildOpenAICompactSSEPayload_EmitsItemsAndCompleted(t *testing.T) {
	finalResponse := []byte(`{
		"id":"resp_compact_1",
		"object":"response",
		"model":"gpt-5.1-codex",
		"status":"completed",
		"output":[
			{"id":"cmp_1","type":"compaction","status":"completed","encrypted_content":"compact-payload","summary":[{"type":"summary_text","text":"compact summary"placeholder],"opaque":{"kept":trueplaceholderplaceholder,
			{"id":"msg_1","type":"message","role":"assistant","content":[{"type":"output_text","text":"done"placeholder]placeholder
		],
		"usage":{"input_tokens":9,"output_tokens":4,"total_tokens":13placeholder
placeholder`)

	payload, ok := buildOpenAICompactSSEPayload(finalResponse)
	require.True(t, ok)

	events := parseCompactBridgeSSE(t, string(payload))
	require.Len(t, events, 3)

	require.Equal(t, "response.output_item.done", events[0][0])
	first := events[0][1]
	require.Equal(t, "response.output_item.done", gjson.Get(first, "type").String())
	require.Equal(t, int64(0), gjson.Get(first, "output_index").Int())
	require.Equal(t, "compaction", gjson.Get(first, "item.type").String())
	require.Equal(t, "cmp_1", gjson.Get(first, "item.id").String())
	require.Equal(t, "compact-payload", gjson.Get(first, "item.encrypted_content").String())
	require.Equal(t, "compact summary", gjson.Get(first, "item.summary.0.text").String())
	require.True(t, gjson.Get(first, "item.opaque.kept").Bool(), "item 原始字段必须逐字节保留")

	require.Equal(t, "response.output_item.done", events[1][0])
	require.Equal(t, int64(1), gjson.Get(events[1][1], "output_index").Int())
	require.Equal(t, "message", gjson.Get(events[1][1], "item.type").String())

	require.Equal(t, "response.completed", events[2][0])
	completed := events[2][1]
	require.Equal(t, "response.completed", gjson.Get(completed, "type").String())
	require.Equal(t, "resp_compact_1", gjson.Get(completed, "response.id").String())
	require.Equal(t, int64(13), gjson.Get(completed, "response.usage.total_tokens").Int())
	require.Len(t, gjson.Get(completed, "response.output").Array(), 2)
placeholder

func TestBuildOpenAICompactSSEPayload_InjectsMissingResponseID(t *testing.T) {
	payload, ok := buildOpenAICompactSSEPayload([]byte(`{"output":[{"type":"compaction","encrypted_content":"x"placeholder]placeholder`))
	require.True(t, ok)

	events := parseCompactBridgeSSE(t, string(payload))
	require.Len(t, events, 2)
	completed := events[1][1]
	// Codex 的 ResponseCompleted 解析要求 response.id 为非空 string，缺失时必须注入。
	id := gjson.Get(completed, "response.id").String()
	require.True(t, strings.HasPrefix(id, "resp_"), "缺失 id 必须注入 resp_* 兜底: %q", id)
	require.NotEqual(t, "resp_", id)
placeholder

func TestBuildOpenAICompactSSEPayload_DropsMalformedUsage(t *testing.T) {
	payload, ok := buildOpenAICompactSSEPayload([]byte(`{
		"id":"resp_1",
		"output":[{"type":"compaction","encrypted_content":"x"placeholder],
		"usage":{"prompt_tokens":9,"completion_tokens":4placeholder
placeholder`))
	require.True(t, ok)

	events := parseCompactBridgeSSE(t, string(payload))
	completed := events[len(events)-1][1]
	// usage 缺少 Codex 必需的整数字段时必须整体删除，否则 completed 事件解析失败。
	require.False(t, gjson.Get(completed, "response.usage").Exists())
placeholder

func TestBuildOpenAICompactSSEPayload_KeepsWellFormedUsage(t *testing.T) {
	payload, ok := buildOpenAICompactSSEPayload([]byte(`{
		"id":"resp_1",
		"output":[{"type":"compaction","encrypted_content":"x"placeholder],
		"usage":{"input_tokens":9,"output_tokens":4,"total_tokens":13,"input_tokens_details":{"cached_tokens":2placeholderplaceholder
placeholder`))
	require.True(t, ok)

	events := parseCompactBridgeSSE(t, string(payload))
	completed := events[len(events)-1][1]
	require.Equal(t, int64(9), gjson.Get(completed, "response.usage.input_tokens").Int())
	require.Equal(t, int64(2), gjson.Get(completed, "response.usage.input_tokens_details.cached_tokens").Int())
placeholder

func TestBuildOpenAICompactSSEPayload_RejectsNonJSONObject(t *testing.T) {
	for name, body := range map[string][]byte{
		"empty":     nil,
		"sse_text":  []byte("data: {\"type\":\"response.completed\"placeholder\n\n"),
		"array":     []byte(`[{"id":"resp_1"placeholder]`),
		"non_json":  []byte("upstream said no"),
		"bare_true": []byte("true"),
placeholder {
		_, ok := buildOpenAICompactSSEPayload(body)
		require.False(t, ok, "case %s 不应被合成为 SSE", name)
placeholder
placeholder

func TestWriteOpenAICompactSSEBridge_RequiresMarkAndSuccessStatus(t *testing.T) {
	finalResponse := []byte(`{"id":"resp_1","output":[{"type":"compaction","encrypted_content":"x"placeholder]placeholder`)

	// 未标记 client stream：不写出，走原 JSON 路径。
	c, rec := newCompactBridgeTestContext(t, false)
	require.False(t, writeOpenAICompactSSEBridge(c, http.StatusOK, finalResponse))
	require.Zero(t, rec.Body.Len())

	// 标记但上游非 2xx：错误响应保持 JSON 原样（Codex 依赖 HTTP 状态码走重试）。
	c, rec = newCompactBridgeTestContext(t, true)
	require.False(t, writeOpenAICompactSSEBridge(c, http.StatusBadGateway, finalResponse))
	require.Zero(t, rec.Body.Len())

	// 标记且 2xx：合成 SSE。
	c, rec = newCompactBridgeTestContext(t, true)
	require.True(t, writeOpenAICompactSSEBridge(c, http.StatusOK, finalResponse))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "text/event-stream", rec.Header().Get("Content-Type"))
	require.Contains(t, rec.Body.String(), "event: response.completed")
placeholder

// 回归 #3875：body-signal 提升后的 compact 请求，上游返回 unary JSON，
// 客户端（Codex remote compact v2）必须收到 SSE 事件流而非 JSON 文档，
// 否则报 "stream closed before response.completed" 并无限重连。
func TestHandleNonStreamingResponse_CompactClientStreamBridgesToSSE(t *testing.T) {
	svc := newCompactBridgeTestService()
	c, rec := newCompactBridgeTestContext(t, true)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(`{
			"id":"resp_compact_json",
			"object":"response",
			"model":"gpt-5.1-codex",
			"status":"completed",
			"output":[{"id":"cmp_1","type":"compaction","status":"completed","encrypted_content":"compact-payload"placeholder],
			"usage":{"input_tokens":9,"output_tokens":4,"total_tokens":13placeholder
	placeholder`)),
placeholder

	result, err := svc.handleNonStreamingResponse(context.Background(), resp, c, &Account{ID: 1, Type: AccountTypeOAuthplaceholder, "gpt-5.5", "gpt-5.5")
placeholder
	require.NotNil(t, result)

	require.Equal(t, "text/event-stream", rec.Header().Get("Content-Type"))
	events := parseCompactBridgeSSE(t, rec.Body.String())
	require.Len(t, events, 2)
	require.Equal(t, "response.output_item.done", events[0][0])
	require.Equal(t, "compaction", gjson.Get(events[0][1], "item.type").String())
	require.Equal(t, "response.completed", events[1][0])
	require.Equal(t, "resp_compact_json", gjson.Get(events[1][1], "response.id").String())

	// 计费与响应元数据不受写回形态影响。
	require.NotNil(t, result.usage)
	require.Equal(t, 9, result.usage.InputTokens)
	require.Equal(t, 4, result.usage.OutputTokens)
	require.Equal(t, "resp_compact_json", result.responseID)
placeholder

// 回归防护：path-based compact（Codex v1 unary 协议、链式 sub2api）未标记
// client stream，必须保持 v0.1.146 以来的 JSON 写回行为。
func TestHandleNonStreamingResponse_PathBasedCompactStaysJSON(t *testing.T) {
	svc := newCompactBridgeTestService()
	c, rec := newCompactBridgeTestContext(t, false)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(`{
			"id":"resp_compact_json",
			"output":[{"id":"cmp_1","type":"compaction","encrypted_content":"compact-payload"placeholder],
			"usage":{"input_tokens":9,"output_tokens":4,"total_tokens":13placeholder
	placeholder`)),
placeholder

	result, err := svc.handleNonStreamingResponse(context.Background(), resp, c, &Account{ID: 1, Type: AccountTypeOAuthplaceholder, "gpt-5.5", "gpt-5.5")
placeholder
	require.NotNil(t, result)

	require.NotContains(t, rec.Header().Get("Content-Type"), "text/event-stream")
	body := rec.Body.String()
	require.Equal(t, "resp_compact_json", gjson.Get(body, "id").String())
	require.Equal(t, "compaction", gjson.Get(body, "output.0.type").String())
placeholder

// 上游对 compact 返回 SSE（如链式网关）时，最终响应经 SSE→JSON 提取后，
// 对 client-stream 请求同样必须再合成回 SSE。
func TestHandleSSEToJSON_CompactClientStreamBridgesToSSE(t *testing.T) {
	svc := newCompactBridgeTestService()
	c, rec := newCompactBridgeTestContext(t, true)
	upstreamSSE := strings.Join([]string{
		`data: {"type":"response.completed","response":{"id":"resp_compact_sse","object":"response","model":"gpt-5.1-codex","status":"completed","output":[{"id":"cmp_sse_1","type":"compaction","status":"completed","encrypted_content":"compact-sse-payload"placeholder],"usage":{"input_tokens":3,"output_tokens":2,"total_tokens":5placeholderplaceholderplaceholder`,
		"",
placeholder, "\n")
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamSSE)),
placeholder

	result, err := svc.handleNonStreamingResponse(context.Background(), resp, c, &Account{ID: 1, Type: AccountTypeOAuthplaceholder, "gpt-5.5", "gpt-5.5")
placeholder
	require.NotNil(t, result)

	require.Equal(t, "text/event-stream", rec.Header().Get("Content-Type"))
	events := parseCompactBridgeSSE(t, rec.Body.String())
	require.Len(t, events, 2)
	require.Equal(t, "response.output_item.done", events[0][0])
	require.Equal(t, "compact-sse-payload", gjson.Get(events[0][1], "item.encrypted_content").String())
	require.Equal(t, "response.completed", events[1][0])
	require.Equal(t, "resp_compact_sse", gjson.Get(events[1][1], "response.id").String())
placeholder

// 回归 #3887（#3777 问题 2）：上游对 compact 返回 SSE，compaction item 只在
// raw output_item.done 中、终态 response.completed 的 output 为空。SSE→JSON
// 提取必须保留 raw item 修补终态 output，否则桥接合成 0 个 output_item.done，
// Codex 报 "expected exactly one compaction output item, got 0" 并盲目重试，
// 每次重试都重新计费。fixture 取自 #3777 的上游实录形态。
func TestHandleSSEToJSON_CompactRawOutputItemDoneRepairsEmptyTerminalOutput(t *testing.T) {
	svc := newCompactBridgeTestService()
	c, rec := newCompactBridgeTestContext(t, true)
	upstreamSSE := strings.Join([]string{
		`data: {"type":"response.output_item.done","output_index":0,"item":{"id":"cmp_1","type":"compaction_summary","status":"completed","summary":[{"type":"summary_text","text":"compact summary"placeholder],"encrypted_content":"compact-payload","opaque":{"kept":trueplaceholderplaceholderplaceholder`,
		``,
		`data: {"type":"response.completed","response":{"id":"resp_compact","object":"response","model":"gpt-5.1-codex","status":"completed","output":[],"usage":{"input_tokens":9,"output_tokens":4,"total_tokens":13placeholderplaceholderplaceholder`,
		``,
placeholder, "\n")
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamSSE)),
placeholder

	result, err := svc.handleNonStreamingResponse(context.Background(), resp, c, &Account{ID: 1, Type: AccountTypeOAuthplaceholder, "gpt-5.5", "gpt-5.5")
placeholder
	require.NotNil(t, result)

	require.Equal(t, "text/event-stream", rec.Header().Get("Content-Type"))
	events := parseCompactBridgeSSE(t, rec.Body.String())
	require.Len(t, events, 2)
	require.Equal(t, "response.output_item.done", events[0][0])
	item := gjson.Get(events[0][1], "item")
	require.Equal(t, "compaction_summary", item.Get("type").String())
	require.Equal(t, "cmp_1", item.Get("id").String())
	require.Equal(t, "compact-payload", item.Get("encrypted_content").String())
	require.Equal(t, "compact summary", item.Get("summary.0.text").String())
	require.True(t, item.Get("opaque.kept").Bool(), "raw item 字段必须逐字节保留")
	require.Equal(t, "response.completed", events[1][0])
	require.Equal(t, "resp_compact", gjson.Get(events[1][1], "response.id").String())
	require.Len(t, gjson.Get(events[1][1], "response.output").Array(), 1)
	require.Equal(t, int64(13), gjson.Get(events[1][1], "response.usage.total_tokens").Int())

	require.NotNil(t, result.usage)
	require.Equal(t, 9, result.usage.InputTokens)
	require.Equal(t, 4, result.usage.OutputTokens)
placeholder

// 同一形态经透传分支（handlePassthroughSSEToJSON）也必须修补。
func TestHandlePassthroughSSEToJSON_CompactRawOutputItemDoneRepairsEmptyTerminalOutput(t *testing.T) {
	svc := newCompactBridgeTestService()
	c, rec := newCompactBridgeTestContext(t, true)
	upstreamSSE := strings.Join([]string{
		`data: {"type":"response.output_item.done","output_index":0,"item":{"id":"cmp_pt_1","type":"compaction","status":"completed","encrypted_content":"compact-pt-raw"placeholderplaceholder`,
		``,
		`data: {"type":"response.completed","response":{"id":"resp_compact_pt_raw","object":"response","status":"completed","output":[],"usage":{"input_tokens":6,"output_tokens":2,"total_tokens":8placeholderplaceholderplaceholder`,
		``,
placeholder, "\n")
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamSSE)),
placeholder

	result, err := svc.handleNonStreamingResponsePassthrough(context.Background(), resp, c, nil, "gpt-5.5", "")
placeholder
	require.NotNil(t, result)

	require.Equal(t, "text/event-stream", rec.Header().Get("Content-Type"))
	events := parseCompactBridgeSSE(t, rec.Body.String())
	require.Len(t, events, 2)
	require.Equal(t, "compaction", gjson.Get(events[0][1], "item.type").String())
	require.Equal(t, "compact-pt-raw", gjson.Get(events[0][1], "item.encrypted_content").String())
	require.Len(t, gjson.Get(events[1][1], "response.output").Array(), 1)
placeholder

// path-based（Codex v1 unary、链式 sub2api）未标记 client stream：同一上游
// 形态修补后仍按 JSON 写回，output 中必须包含 compaction item。
func TestHandleSSEToJSON_PathBasedCompactRawOutputItemDoneRepairsJSON(t *testing.T) {
	svc := newCompactBridgeTestService()
	c, rec := newCompactBridgeTestContext(t, false)
	upstreamSSE := strings.Join([]string{
		`data: {"type":"response.output_item.done","output_index":0,"item":{"id":"cmp_v1","type":"compaction_summary","encrypted_content":"compact-v1-raw"placeholderplaceholder`,
		``,
		`data: {"type":"response.completed","response":{"id":"resp_compact_v1","object":"response","status":"completed","output":[],"usage":{"input_tokens":5,"output_tokens":1,"total_tokens":6placeholderplaceholderplaceholder`,
		``,
placeholder, "\n")
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamSSE)),
placeholder

	result, err := svc.handleNonStreamingResponse(context.Background(), resp, c, &Account{ID: 1, Type: AccountTypeOAuthplaceholder, "gpt-5.5", "gpt-5.5")
placeholder
	require.NotNil(t, result)

	// 写回 body 必须是修补后的 JSON 文档（非 SSE 事件流）。
	body := rec.Body.String()
	require.NotContains(t, body, "event:")
	require.NotContains(t, body, "data:")
	require.Equal(t, "resp_compact_v1", gjson.Get(body, "id").String())
	require.Equal(t, "compaction_summary", gjson.Get(body, "output.0.type").String())
	require.Equal(t, "compact-v1-raw", gjson.Get(body, "output.0.encrypted_content").String())
placeholder

// raw done item 是协议上的最终完整形态，优先于 delta 重建且不得重复计入。
func TestReconstructResponseOutputFromSSE_PrefersRawDoneItems(t *testing.T) {
	bodyText := strings.Join([]string{
		`data: {"type":"response.output_text.delta","delta":"hel"placeholder`,
		`data: {"type":"response.output_text.delta","delta":"lo"placeholder`,
		`data: {"type":"response.output_item.done","output_index":0,"item":{"id":"msg_1","type":"message","role":"assistant","status":"completed","content":[{"type":"output_text","text":"hello"placeholder]placeholderplaceholder`,
		`data: {"type":"response.completed","response":{"id":"resp_1","output":[]placeholderplaceholder`,
placeholder, "\n")

	outputJSON, ok := reconstructResponseOutputFromSSE(bodyText)
	require.True(t, ok)
	items := gjson.ParseBytes(outputJSON).Array()
	require.Len(t, items, 1, "raw done item 与 delta 重建不得重复")
	require.Equal(t, "msg_1", items[0].Get("id").String())
	require.Equal(t, "hello", items[0].Get("content.0.text").String())
placeholder

// 无任何 done 事件时，退回收集 output_item.added 中的 compaction 类 item。
func TestReconstructResponseOutputFromSSE_CompactionAddedFallback(t *testing.T) {
	bodyText := strings.Join([]string{
		`data: {"type":"response.output_item.added","output_index":0,"item":{"id":"cmp_add","type":"compaction","encrypted_content":"added-only"placeholderplaceholder`,
		`data: {"type":"response.completed","response":{"id":"resp_1","output":[]placeholderplaceholder`,
placeholder, "\n")

	outputJSON, ok := reconstructResponseOutputFromSSE(bodyText)
	require.True(t, ok)
	items := gjson.ParseBytes(outputJSON).Array()
	require.Len(t, items, 1)
	require.Equal(t, "compaction", items[0].Get("type").String())
	require.Equal(t, "added-only", items[0].Get("encrypted_content").String())
placeholder

// 混合形态：其他 item 有 done、compaction 只在 added 中——compaction 必须
// 被补入；done 已含 compaction 时 added 不得重复计入。
func TestReconstructResponseOutputFromSSE_MixedDoneAndCompactionAdded(t *testing.T) {
	bodyText := strings.Join([]string{
		`data: {"type":"response.output_item.added","output_index":0,"item":{"id":"cmp_mixed","type":"compaction","encrypted_content":"mixed"placeholderplaceholder`,
		`data: {"type":"response.output_item.done","output_index":1,"item":{"id":"msg_1","type":"message","content":[{"type":"output_text","text":"hi"placeholder]placeholderplaceholder`,
		`data: {"type":"response.completed","response":{"id":"resp_1","output":[]placeholderplaceholder`,
placeholder, "\n")

	outputJSON, ok := reconstructResponseOutputFromSSE(bodyText)
	require.True(t, ok)
	items := gjson.ParseBytes(outputJSON).Array()
	require.Len(t, items, 2)
	require.Equal(t, "msg_1", items[0].Get("id").String())
	require.Equal(t, "cmp_mixed", items[1].Get("id").String())

	// done 已含 compaction：added 中的同一 item（无 id 可去重的最坏情况用
	// 不同 raw 表达）不得再收集，Codex 要求恰好一个 compaction item。
	bodyText = strings.Join([]string{
		`data: {"type":"response.output_item.added","output_index":0,"item":{"type":"compaction","status":"in_progress"placeholderplaceholder`,
		`data: {"type":"response.output_item.done","output_index":0,"item":{"type":"compaction","status":"completed","encrypted_content":"final"placeholderplaceholder`,
		`data: {"type":"response.completed","response":{"id":"resp_1","output":[]placeholderplaceholder`,
placeholder, "\n")
	outputJSON, ok = reconstructResponseOutputFromSSE(bodyText)
	require.True(t, ok)
	items = gjson.ParseBytes(outputJSON).Array()
	require.Len(t, items, 1)
	require.Equal(t, "final", items[0].Get("encrypted_content").String())
placeholder

// 上游不一致形态：终态 output 非空（含 message）但 compaction 只在 raw
// output_item.done 中。146 纯流式透传下 Codex 直接读事件流能拿到 compaction，
// SSE→JSON 提取必须补入等价结果。
func TestHandleSSEToJSON_CompactSupplementsMissingCompactionIntoNonEmptyOutput(t *testing.T) {
	svc := newCompactBridgeTestService()
	c, rec := newCompactBridgeTestContext(t, true)
	upstreamSSE := strings.Join([]string{
		`data: {"type":"response.output_item.done","output_index":0,"item":{"id":"cmp_sup","type":"compaction","encrypted_content":"supplement"placeholderplaceholder`,
		``,
		`data: {"type":"response.completed","response":{"id":"resp_sup","object":"response","status":"completed","output":[{"id":"msg_sup","type":"message","role":"assistant","content":[{"type":"output_text","text":"note"placeholder]placeholder],"usage":{"input_tokens":2,"output_tokens":1,"total_tokens":3placeholderplaceholderplaceholder`,
		``,
placeholder, "\n")
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamSSE)),
placeholder

	result, err := svc.handleNonStreamingResponse(context.Background(), resp, c, &Account{ID: 1, Type: AccountTypeOAuthplaceholder, "gpt-5.5", "gpt-5.5")
placeholder
	require.NotNil(t, result)

	events := parseCompactBridgeSSE(t, rec.Body.String())
	require.Len(t, events, 3)
	itemTypes := []string{
		gjson.Get(events[0][1], "item.type").String(),
		gjson.Get(events[1][1], "item.type").String(),
placeholder
	require.Contains(t, itemTypes, "compaction")
	require.Contains(t, itemTypes, "message")
	require.Equal(t, "response.completed", events[2][0])
	require.Len(t, gjson.Get(events[2][1], "response.output").Array(), 2)
placeholder

// 补全逻辑的门控：非 compact 请求原样返回；终态已含 compaction 不重复补入。
func TestSupplementCompactionItemFromSSE_Gating(t *testing.T) {
	bodyText := `data: {"type":"response.output_item.done","item":{"id":"cmp_g","type":"compaction","encrypted_content":"g"placeholderplaceholder` + "\n"

	// 非 compact 路径：不补入。
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	finalResponse := []byte(`{"id":"r1","output":[{"type":"message"placeholder]placeholder`)
	require.Equal(t, string(finalResponse), string(supplementCompactionItemFromSSE(c, finalResponse, bodyText)))

	// compact 路径 + 终态已含 compaction：不重复补入。
	c2, _ := newCompactBridgeTestContext(t, false)
	already := []byte(`{"id":"r2","output":[{"type":"compaction","encrypted_content":"x"placeholder]placeholder`)
	require.Equal(t, string(already), string(supplementCompactionItemFromSSE(c2, already, bodyText)))

	// compact 路径 + 终态非空缺 compaction：补入到末尾。
	missing := []byte(`{"id":"r3","output":[{"type":"message"placeholder]placeholder`)
	patched := supplementCompactionItemFromSSE(c2, missing, bodyText)
	items := gjson.GetBytes(patched, "output").Array()
	require.Len(t, items, 2)
	require.Equal(t, "compaction", items[1].Get("type").String())
	require.Equal(t, "g", items[1].Get("encrypted_content").String())
placeholder

// 非 compaction 的 output_item.added 不参与回退收集（added 阶段的 message
// 通常是空壳），仍走 delta 重建。
func TestReconstructResponseOutputFromSSE_NonCompactionAddedStillUsesDeltas(t *testing.T) {
	bodyText := strings.Join([]string{
		`data: {"type":"response.output_item.added","output_index":0,"item":{"id":"msg_1","type":"message","content":[]placeholderplaceholder`,
		`data: {"type":"response.output_text.delta","delta":"hi"placeholder`,
		`data: {"type":"response.completed","response":{"id":"resp_1","output":[]placeholderplaceholder`,
placeholder, "\n")

	outputJSON, ok := reconstructResponseOutputFromSSE(bodyText)
	require.True(t, ok)
	items := gjson.ParseBytes(outputJSON).Array()
	require.Len(t, items, 1)
	require.Equal(t, "hi", items[0].Get("content.0.text").String())
placeholder

// 透传分支（OAuth passthrough）同样命中桥接。
func TestHandleNonStreamingResponsePassthrough_CompactClientStreamBridgesToSSE(t *testing.T) {
	svc := newCompactBridgeTestService()
	c, rec := newCompactBridgeTestContext(t, true)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(`{
			"id":"resp_compact_pt",
			"output":[{"id":"cmp_pt_1","type":"compaction","encrypted_content":"compact-pt-payload"placeholder],
			"usage":{"input_tokens":7,"output_tokens":3,"total_tokens":10placeholder
	placeholder`)),
placeholder

	result, err := svc.handleNonStreamingResponsePassthrough(context.Background(), resp, c, nil, "gpt-5.5", "")
placeholder
	require.NotNil(t, result)

	require.Equal(t, "text/event-stream", rec.Header().Get("Content-Type"))
	events := parseCompactBridgeSSE(t, rec.Body.String())
	require.Len(t, events, 2)
	require.Equal(t, "compaction", gjson.Get(events[0][1], "item.type").String())
	require.Equal(t, "resp_compact_pt", gjson.Get(events[1][1], "response.id").String())
	require.NotNil(t, result.usage)
	require.Equal(t, 7, result.usage.InputTokens)
placeholder
