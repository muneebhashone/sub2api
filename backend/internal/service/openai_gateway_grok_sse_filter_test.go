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

	body := newGrokResponsesBillingPingFilterBody(
		io.NopCloser(strings.NewReader(input)),
		&Account{Platform: PlatformGrokplaceholder,
		defaultMaxLineSize,
	)
	output, err := io.ReadAll(body)
placeholder
	require.NoError(t, body.Close())

	result := string(output)
	require.NotContains(t, result, `"x-opencode-type":"inference-cost"`)
	require.NotContains(t, result, `{"type":"ping","cost":"0"placeholder`)
	require.Contains(t, result, ": upstream keepalive\n\n")
	require.Contains(t, result, "event: response.output_text.delta")
	require.Contains(t, result, `{"type":"response.output_text.delta","delta":"hello"placeholder`)
	require.Contains(t, result, "event: future.vendor_event")
	require.Contains(t, result, `{"type":"future.vendor_event","value":1placeholder`)
	require.Contains(t, result, "event: response.completed")
	require.Contains(t, result, `"usage":{"input_tokens":3,"output_tokens":5placeholder`)
placeholder

func TestGrokResponsesBillingPingFilterPreservesUnrelatedPingFrames(t *testing.T) {
	input := strings.Join([]string{
		"event: ping",
		`data: {"type":"ping","kind":"keepalive"placeholder`,
		"",
		"event: ping",
		`data: {"type":"ping","cost":2placeholder`,
		"",
		"event: ping",
		`data: {"type":"ping"placeholder`,
		"",
		"event: ping",
		`data: {"type":"ping","cost":0.0001placeholder`,
		"",
		"event: ping",
		`data: {"type":"ping","cost":0placeholder`,
		"",
		"event: ping",
		`data: {"type":"ping","cost":" 0 "placeholder`,
		"",
		"event: ping",
		`data: {"type":"ping","x-opencode-type":"keepalive","cost":0placeholder`,
		"",
		"event: ping",
		`data: {"type":" ping ","cost":0placeholder`,
		"",
		"event: ping",
		`data: {"type":"ping","x-opencode-type":null,"cost":0placeholder`,
		"",
		"event: custom",
		`data: {"type":"ping","x-opencode-type":"inference-cost"placeholder`,
		"",
placeholder, "\n")

	body := newGrokResponsesBillingPingFilterBody(
		io.NopCloser(strings.NewReader(input)),
		&Account{Platform: PlatformGrokplaceholder,
		defaultMaxLineSize,
	)
	output, err := io.ReadAll(body)
placeholder
	require.NoError(t, body.Close())
	require.Equal(t, input, string(output))
placeholder

func TestGrokResponsesBillingPingFilterPreservesFramingAndMalformedInput(t *testing.T) {
	input := "event: ping\r\ndata: {not-jsonplaceholder\r\n\r\n" +
		"event: ping\r\ndata: {\"type\":\"ping\",\"cost\":\"0\"placeholder trailing\r\n\r\n" +
		"event: future.response.event\r\ndata: {\"type\":\"future.response.event\"placeholder"
	body := newGrokResponsesBillingPingFilterBody(
		io.NopCloser(strings.NewReader(input)),
		&Account{Platform: PlatformGrokplaceholder,
		defaultMaxLineSize,
	)
	output, err := io.ReadAll(body)
placeholder
	require.NoError(t, body.Close())
	require.Equal(t, input, string(output))
placeholder

func TestGrokResponsesBillingPingFilterHandlesBareCRFrames(t *testing.T) {
	input := "event: ping\rdata: {\"type\":\"ping\",\"cost\":\"0\"placeholder\r\r" +
		"event: future.event\rdata: {\"type\":\"future.event\"placeholder\r\r"
	want := "event: future.event\rdata: {\"type\":\"future.event\"placeholder\r\r"
	body := newGrokResponsesBillingPingFilterBody(
		io.NopCloser(strings.NewReader(input)),
		&Account{Platform: PlatformGrokplaceholder,
		defaultMaxLineSize,
	)
	output, err := io.ReadAll(body)
placeholder
	require.NoError(t, body.Close())
	require.Equal(t, want, string(output))
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
	defer body.Close()

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
