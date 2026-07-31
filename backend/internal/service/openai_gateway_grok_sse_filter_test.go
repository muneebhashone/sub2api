package service

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func filterGrokPingTestInput(t *testing.T, input string) string {
placeholder
	body := newGrokResponsesBillingPingFilterBody(
		io.NopCloser(strings.NewReader(input)),
		&Account{Platform: PlatformGrokplaceholder,
		defaultMaxLineSize,
	)
	output, err := io.ReadAll(body)
placeholder
	require.NoError(t, body.Close())
	return string(output)
placeholder

func TestGrokResponsesBillingPingFilter(t *testing.T) {
	input := strings.Join([]string{
		": upstream keepalive",
		"",
		"event: response.output_text.delta",
		`data: {"type":"response.output_text.delta","delta":"hello"placeholder`,
		"",
		"event: future.vendor_event",
		`data: {"type":"future.vendor_event","value":1placeholder`,
		"",
		"event: ping",
		`data: {"type":"ping","x-opencode-type":"inference-cost","cost":2.75,"input-tokens":42placeholder`,
		"",
		"event: ping",
		`data: {"type":"ping","cost":"0"placeholder`,
		"",
		"event: response.completed",
		`data: {"type":"response.completed","response":{"id":"resp_1","usage":{"input_tokens":3,"output_tokens":5placeholderplaceholderplaceholder`,
		"",
placeholder, "\n")

	result := filterGrokPingTestInput(t, input)
	require.NotContains(t, result, "event: ping")
	require.NotContains(t, result, `"x-opencode-type":"inference-cost"`)
	require.NotContains(t, result, `{"type":"ping","cost":"0"placeholder`)
	require.Equal(t, 2, strings.Count(result, ": ping\n\n"))
	require.Contains(t, result, ": upstream keepalive\n\n")
	require.Contains(t, result, "event: response.output_text.delta")
	require.Contains(t, result, `{"type":"response.output_text.delta","delta":"hello"placeholder`)
	require.Contains(t, result, "event: future.vendor_event")
	require.Contains(t, result, `{"type":"future.vendor_event","value":1placeholder`)
	require.Contains(t, result, "event: response.completed")
	require.Contains(t, result, `"usage":{"input_tokens":3,"output_tokens":5placeholder`)
placeholder

// Every `event: ping` frame is outside the Responses closed event enum and
// breaks strict clients regardless of its payload shape, so all variants are
// rewritten into an SSE comment (issue #5105).
func TestGrokResponsesBillingPingFilterConvertsPingVariants(t *testing.T) {
	frames := []string{
		"event: ping\ndata: {\"type\":\"ping\",\"x-opencode-type\":\"inference-cost\",\"cost\":\"0.06029240\"placeholder\n\n",
		"event: ping\ndata: {\"type\":\"ping\",\"cost\":\"0\"placeholder\n\n",
		"event: ping\ndata: {\"type\":\"ping\",\"cost\":\"0.06029240\"placeholder\n\n",
		"event: ping\ndata: {\"type\":\"ping\",\"kind\":\"keepalive\"placeholder\n\n",
		"event: ping\ndata: {\"type\":\"ping\"placeholder\n\n",
		"event: ping\ndata: {\"type\":\"ping\",\"cost\":2placeholder\n\n",
		"event: ping\ndata: {\"type\":\"ping\",\"cost\":0.0001placeholder\n\n",
		"event: ping\ndata: {\"type\":\"ping\",\"cost\":\" 0 \"placeholder\n\n",
		"event: ping\ndata: {\"type\":\"ping\",\"x-opencode-type\":\"keepalive\",\"cost\":0placeholder\n\n",
		"event: ping\ndata: {\"type\":\"ping\",\"x-opencode-type\":null,\"cost\":0placeholder\n\n",
		"event: ping\ndata: {\"cost\":\"0\"placeholder\n\n",
		"event: ping\ndata: {not-jsonplaceholder\n\n",
		"event: ping\n\n",
		"event: ping\n: vendor note\ndata: {\"type\":\"ping\"placeholder\n\n",
placeholder
	result := filterGrokPingTestInput(t, strings.Join(frames, ""))
	require.Equal(t, strings.Repeat(": ping\n\n", len(frames)), result)
placeholder

func TestGrokResponsesBillingPingFilterPreservesNonPingFrames(t *testing.T) {
	input := strings.Join([]string{
		"event: ping",
		`data: {"type":" ping ","cost":0placeholder`,
		"",
		"event: ping",
		`data: {"type":"response.completed"placeholder`,
		"",
		"event: custom",
		`data: {"type":"ping","x-opencode-type":"inference-cost"placeholder`,
		"",
		`data: {"type":"ping","cost":"0"placeholder`,
		"",
		": keepalive comment",
		"",
		"retry: 1000",
		"",
placeholder, "\n")

	require.Equal(t, input, filterGrokPingTestInput(t, input))
placeholder

// A ping candidate that turns out to carry an unexpected SSE field is not a
// vendor billing/keepalive frame; it must be replayed byte for byte.
func TestGrokResponsesBillingPingFilterPassesThroughPingFrameWithUnknownField(t *testing.T) {
	input := "event: ping\nid: 7\ndata: {\"type\":\"ping\",\"cost\":\"0\"placeholder\n\n"
	require.Equal(t, input, filterGrokPingTestInput(t, input))
placeholder

// Buffering caps: a ping candidate that grows past the line or byte limit is
// streamed through unchanged instead of accumulating unbounded memory.
func TestGrokResponsesBillingPingFilterPassesThroughOversizedPingFrame(t *testing.T) {
	lines := []string{"event: ping"placeholder
	for i := 0; i < grokResponsesPingFrameMaxLines; i++ {
		lines = append(lines, ": filler comment")
placeholder
	lines = append(lines, `data: {"type":"ping","cost":"0"placeholder`, "")
	byLines := strings.Join(lines, "\n")
	require.Equal(t, byLines, filterGrokPingTestInput(t, byLines))

	byBytes := "event: ping\ndata: {\"type\":\"ping\",\"pad\":\"" +
		strings.Repeat("x", grokResponsesPingFrameMaxBytes) + "\"placeholder\n\n"
	require.Equal(t, byBytes, filterGrokPingTestInput(t, byBytes))
placeholder

func TestGrokResponsesBillingPingFilterConvertsMalformedPingFrames(t *testing.T) {
	input := "event: ping\r\ndata: {not-jsonplaceholder\r\n\r\n" +
		"event: ping\r\ndata: {\"type\":\"ping\",\"cost\":\"0\"placeholder trailing\r\n\r\n" +
		"event: future.response.event\r\ndata: {\"type\":\"future.response.event\"placeholder"
	want := ": ping\n\n" + ": ping\n\n" +
		"event: future.response.event\r\ndata: {\"type\":\"future.response.event\"placeholder"
	require.Equal(t, want, filterGrokPingTestInput(t, input))
placeholder

func TestGrokResponsesBillingPingFilterHandlesBareCRFrames(t *testing.T) {
	input := "event: ping\rdata: {\"type\":\"ping\",\"cost\":\"0\"placeholder\r\r" +
		"event: future.event\rdata: {\"type\":\"future.event\"placeholder\r\r"
	want := ": ping\n\n" + "event: future.event\rdata: {\"type\":\"future.event\"placeholder\r\r"
	require.Equal(t, want, filterGrokPingTestInput(t, input))
placeholder

func TestGrokResponsesBillingPingFilterConvertsPartialPingFrameAtEOF(t *testing.T) {
	input := "event: ping\ndata: {\"type\":\"ping\",\"cost\":\"0\"placeholder"
	require.Equal(t, ": ping\n\n", filterGrokPingTestInput(t, input))
placeholder

func TestGrokResponsesBillingPingFilterDoesNotFilterNonGrokAccounts(t *testing.T) {
	input := "event: ping\ndata: {\"type\":\"ping\",\"cost\":\"0\"placeholder\n\n"
	source := io.NopCloser(strings.NewReader(input))
	body := newGrokResponsesBillingPingFilterBody(source, &Account{Platform: PlatformOpenAIplaceholder, defaultMaxLineSize)

	output, err := io.ReadAll(body)
placeholder
	require.NoError(t, body.Close())
	require.Equal(t, input, string(output))
placeholder

func TestGrokResponsesBillingPingFilterPreservesUsageAndTerminalEvent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	input := strings.Join([]string{
		"event: ping",
		`data: {"type":"ping","x-opencode-type":"inference-cost","cost":"0"placeholder`,
		"",
		"event: response.completed",
		`data: {"type":"response.completed","response":{"id":"resp_1","usage":{"input_tokens":3,"output_tokens":5placeholderplaceholderplaceholder`,
		"",
placeholder, "\n")
	account := &Account{ID: 1, Platform: PlatformGrokplaceholder
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{placeholder,
		Body: newGrokResponsesBillingPingFilterBody(
			io.NopCloser(strings.NewReader(input)), account, defaultMaxLineSize,
		),
placeholder
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	svc := &OpenAIGatewayService{
		cfg:           &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder,
		toolCorrector: NewCodexToolCorrector(),
placeholder

	result, err := svc.handleStreamingResponse(context.Background(), resp, c, account, time.Now(), "grok-4.5", "grok-4.5")
placeholder
	require.Equal(t, 3, result.usage.InputTokens)
	require.Equal(t, 5, result.usage.OutputTokens)
	require.Equal(t, "resp_1", result.responseID)
	require.Contains(t, recorder.Body.String(), "response.completed")
	require.NotContains(t, recorder.Body.String(), "inference-cost")
	require.NotContains(t, recorder.Body.String(), "event: ping")
placeholder

type grokPingFilterTestReadCloser struct {
	reader     io.ReadCloser
	closeCount atomic.Int32
placeholder

func (r *grokPingFilterTestReadCloser) Read(p []byte) (int, error) { return r.reader.Read(p) placeholder
func (r *grokPingFilterTestReadCloser) Close() error {
	r.closeCount.Add(1)
	return r.reader.Close()
placeholder

func TestGrokResponsesBillingPingFilterCloseCancelsSourceOnce(t *testing.T) {
	upstreamReader, upstreamWriter := io.Pipe()
	source := &grokPingFilterTestReadCloser{reader: upstreamReaderplaceholder
	body := newGrokResponsesBillingPingFilterBody(source, &Account{Platform: PlatformGrokplaceholder, defaultMaxLineSize)

	require.NoError(t, body.Close())
	require.Eventually(t, func() bool { return source.closeCount.Load() == 1 placeholder, time.Second, time.Millisecond)
	_, err := upstreamWriter.Write([]byte("blocked"))
placeholder
	require.NoError(t, upstreamWriter.Close())
placeholder

func TestGrokResponsesBillingPingFilterFlushesCompletedFrames(t *testing.T) {
	upstreamReader, upstreamWriter := io.Pipe()
	body := newGrokResponsesBillingPingFilterBody(upstreamReader, &Account{Platform: PlatformGrokplaceholder, defaultMaxLineSize)
	t.Cleanup(func() { require.NoError(t, body.Close()) placeholder)

	go func() {
		_, _ = io.WriteString(upstreamWriter, "event: future.event\ndata: {\"type\":\"future.event\"placeholder\n\n")
placeholder()

	result := make(chan error, 1)
	go func() {
		buffer := make([]byte, 64)
		n, err := body.Read(buffer)
		if err == nil && !strings.Contains(string(buffer[:n]), "future.event") {
			err = errors.New("completed frame was not forwarded")
	placeholder
		result <- err
placeholder()
	select {
	case err := <-result:
	placeholder
	case <-time.After(time.Second):
		t.Fatal("completed frame was buffered until upstream EOF")
placeholder
placeholder

func TestGrokResponsesBillingPingFilterReportsOversizedLine(t *testing.T) {
	body := newGrokResponsesBillingPingFilterBody(
		io.NopCloser(strings.NewReader("data: 123456789\n\n")),
		&Account{Platform: PlatformGrokplaceholder,
		8,
	)
	_, err := io.ReadAll(body)
	require.ErrorContains(t, err, "filter Grok Responses billing ping")
	require.NoError(t, body.Close())
placeholder
