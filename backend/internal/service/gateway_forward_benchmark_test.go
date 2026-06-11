//go:build unit

// Phase-0 TASK-005 热路径基准：覆盖 anthropic API Key 透传账号的完整
// Forward 路径（非流式 + 流式 SSE）。与 scripts/bench-baseline.sh 配套使用，
// 基线存于 testdata/bench/baseline.txt；对比策略：allocs/op 严格（不允许增加）、
// ns/op 宽松（容忍 15% 抖动）。复用 TASK-002 特征化测试的夹具（passChar* 前缀）。
package service

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
)

func benchFwdRun(b *testing.B, header http.Header, upstreamBody string, reqBody string, stream bool) {
	b.Helper()
	gin.SetMode(gin.TestMode)
	// 透传分支每请求有一条 stdlib log，静音以降低基准 I/O 噪声
	// （日志成本在 baseline 与 compare 两侧同等消除，不影响对比有效性）。
	prevLogOut := log.Writer()
	log.SetOutput(io.Discard)
	b.Cleanup(func() { log.SetOutput(prevLogOut) placeholder)
	cfg := &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder
	account := passCharAnthropicAccount(9100)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		upstream := &anthropicHTTPUpstreamRecorder{
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Header:     header.Clone(),
				Body:       io.NopCloser(strings.NewReader(upstreamBody)),
		placeholder,
	placeholder
		svc := passCharGatewayService(cfg, upstream)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
		parsed := &ParsedRequest{
			Body:   NewRequestBodyRef([]byte(reqBody)),
			Model:  "claude-sonnet-4-20250514",
			Stream: stream,
	placeholder
		if _, err := svc.Forward(context.Background(), c, account, parsed); err != nil {
			b.Fatal(err)
	placeholder
placeholder
placeholder

func BenchmarkGatewayForward_AnthropicNonStreamPassthrough(b *testing.B) {
	upstreamJSON := `{"id":"msg_bench_01","type":"message","content":[{"type":"text","text":"hello"placeholder],` +
		`"usage":{"input_tokens":12,"output_tokens":7,"cache_read_input_tokens":3placeholderplaceholder`
	benchFwdRun(b,
		http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		upstreamJSON,
		`{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"placeholder]placeholder`,
		false,
	)
placeholder

func BenchmarkGatewayForward_AnthropicStreamPassthrough(b *testing.B) {
	events := []string{
		"event: message_start\ndata: {\"type\":\"message_start\",\"message\":{\"id\":\"msg_bench_02\",\"usage\":{\"input_tokens\":9,\"cached_tokens\":3placeholderplaceholderplaceholder",
		"event: content_block_delta\ndata: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"text_delta\",\"text\":\"Hel\"placeholderplaceholder",
		"event: content_block_delta\ndata: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"text_delta\",\"text\":\"lo\"placeholderplaceholder",
		"event: message_delta\ndata: {\"type\":\"message_delta\",\"delta\":{\"stop_reason\":\"end_turn\"placeholder,\"usage\":{\"output_tokens\":5placeholderplaceholder",
		"event: message_stop\ndata: {\"type\":\"message_stop\"placeholder",
placeholder
	benchFwdRun(b,
		http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		strings.Join(events, "\n\n")+"\n\n",
		`{"model":"claude-sonnet-4-20250514","stream":true,"messages":[{"role":"user","content":"hi"placeholder]placeholder`,
		true,
	)
placeholder
