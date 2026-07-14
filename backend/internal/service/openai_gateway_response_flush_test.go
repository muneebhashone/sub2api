package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type openAIResponseFlushRecorder struct {
	header          http.Header
	mu              sync.Mutex
	body            bytes.Buffer
	status          int
	writes          int
	failAfterWrites int
	flushSnapshots  []string
	flushEvents     chan int
	blockFlush      int
	flushBlocked    chan struct{placeholder
	releaseFlush    <-chan struct{placeholder
placeholder

func newOpenAIResponseFlushRecorder() *openAIResponseFlushRecorder {
	return &openAIResponseFlushRecorder{
		header:          make(http.Header),
		failAfterWrites: -1,
		flushEvents:     make(chan int, 16),
placeholder
placeholder

func (w *openAIResponseFlushRecorder) Header() http.Header {
	return w.header
placeholder

func (w *openAIResponseFlushRecorder) WriteHeader(statusCode int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.status == 0 {
		w.status = statusCode
placeholder
placeholder

func (w *openAIResponseFlushRecorder) Write(data []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.failAfterWrites >= 0 && w.writes >= w.failAfterWrites {
		return 0, errors.New("client disconnected")
placeholder
	w.writes++
	if w.status == 0 {
		w.status = http.StatusOK
placeholder
	return w.body.Write(data)
placeholder

func (w *openAIResponseFlushRecorder) Flush() {
	w.mu.Lock()
	w.flushSnapshots = append(w.flushSnapshots, w.body.String())
	count := len(w.flushSnapshots)
	w.mu.Unlock()
	w.flushEvents <- count
	if count == w.blockFlush {
		close(w.flushBlocked)
		<-w.releaseFlush
placeholder
placeholder

func (w *openAIResponseFlushRecorder) snapshot() (string, []string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.body.String(), append([]string(nil), w.flushSnapshots...)
placeholder

type stagedOpenAISSEReadCloser struct {
	segments   [][]byte
	gates      []<-chan struct{placeholder
	waiting    []chan struct{placeholder
	eofReached chan struct{placeholder
	current    []byte
	index      int
placeholder

func (r *stagedOpenAISSEReadCloser) Read(data []byte) (int, error) {
	if len(r.current) == 0 {
		if r.index >= len(r.segments) {
			if r.eofReached != nil {
				close(r.eofReached)
				r.eofReached = nil
		placeholder
			return 0, io.EOF
	placeholder
		index := r.index
		r.index++
		if index < len(r.waiting) && r.waiting[index] != nil {
			close(r.waiting[index])
	placeholder
		if index < len(r.gates) && r.gates[index] != nil {
			<-r.gates[index]
	placeholder
		r.current = r.segments[index]
placeholder
	n := copy(data, r.current)
	r.current = r.current[n:]
	return n, nil
placeholder

func (r *stagedOpenAISSEReadCloser) Close() error { return nil placeholder

type openAIResponseFlushReadError struct {
	payload []byte
	err     error
	sent    bool
placeholder

func (r *openAIResponseFlushReadError) Read(data []byte) (int, error) {
	if !r.sent {
		r.sent = true
		return copy(data, r.payload), nil
placeholder
	if r.err != nil {
		return 0, r.err
placeholder
	return 0, io.ErrUnexpectedEOF
placeholder

func (r *openAIResponseFlushReadError) Close() error { return nil placeholder

func TestOpenAIResponseFlush_SlowEventsFlushOnceAtBoundaries(t *testing.T) {
	events := []string{
		`data: {"type":"response.output_text.delta","delta":"a"placeholder`,
		`data: {"type":"response.output_text.delta","delta":"b"placeholder`,
		`data: {"type":"response.output_text.delta","delta":"c"placeholder`,
		`data: [DONE]`,
placeholder
	body := strings.Join(events, "\n\n") + "\n\n"
	recorder := newOpenAIResponseFlushRecorder()

	result, err := runOpenAIResponseFlushTest(recorder, io.NopCloser(strings.NewReader(body)), config.GatewayConfig{placeholder)

placeholder
	require.NotNil(t, result)
	gotBody, flushes := recorder.snapshot()
	require.Equal(t, body, gotBody)
	require.Len(t, flushes, len(events))
	for _, flushed := range flushes {
		require.True(t, strings.HasSuffix(flushed, "\n\n"), "flush must occur after a complete SSE event")
placeholder
placeholder

func TestOpenAIResponseFlush_DataQueuedButBlankDrainsFlushesOnce(t *testing.T) {
	first := "data: {\"type\":\"response.output_text.delta\",\"delta\":\"first\"placeholder\n\n"
	second := "data: {\"type\":\"response.output_text.delta\",\"delta\":\"second\"placeholder\n\n"
	terminal := "data: [DONE]\n\n"
	allowSecond := make(chan struct{placeholder)
	allowTerminal := make(chan struct{placeholder)
	terminalWaiting := make(chan struct{placeholder)
	reader := &stagedOpenAISSEReadCloser{
		segments: [][]byte{[]byte(first), []byte(second), []byte(terminal)placeholder,
		gates:    []<-chan struct{placeholder{nil, allowSecond, allowTerminalplaceholder,
		waiting:  []chan struct{placeholder{nil, nil, terminalWaitingplaceholder,
placeholder
	releaseFirstFlush := make(chan struct{placeholder)
	recorder := newOpenAIResponseFlushRecorder()
	recorder.blockFlush = 1
	recorder.flushBlocked = make(chan struct{placeholder)
	recorder.releaseFlush = releaseFirstFlush
	resultCh, errCh := runOpenAIResponseFlushTestAsync(recorder, reader, config.GatewayConfig{StreamDataIntervalTimeout: 30placeholder)

	waitOpenAIResponseFlushSignal(t, recorder.flushBlocked)
	close(allowSecond)
	waitOpenAIResponseFlushSignal(t, terminalWaiting)
	close(releaseFirstFlush)
	waitOpenAIResponseFlushCount(t, recorder, 2)
	close(allowTerminal)

	require.NoError(t, <-errCh)
	require.NotNil(t, <-resultCh)
	gotBody, flushes := recorder.snapshot()
	require.Equal(t, first+second+terminal, gotBody)
	require.Len(t, flushes, 3)
	require.Equal(t, first, flushes[0])
	require.Equal(t, first+second, flushes[1], "blank line that drains the queue must flush the complete event exactly once")
placeholder

func TestOpenAIResponseFlush_BurstDoesNotIncreaseFlushes(t *testing.T) {
	first := "data: {\"type\":\"response.output_text.delta\",\"delta\":\"first\"placeholder\n\n"
	burst := strings.Join([]string{
		`data: {"type":"response.output_text.delta","delta":"second"placeholder`,
		`data: {"type":"response.output_text.delta","delta":"third"placeholder`,
		`data: [DONE]`,
placeholder, "\n\n") + "\n\n"
	allowBurst := make(chan struct{placeholder)
	eofReached := make(chan struct{placeholder)
	reader := &stagedOpenAISSEReadCloser{
		segments:   [][]byte{[]byte(first), []byte(burst)placeholder,
		gates:      []<-chan struct{placeholder{nil, allowBurstplaceholder,
		eofReached: eofReached,
placeholder
	releaseFirstFlush := make(chan struct{placeholder)
	recorder := newOpenAIResponseFlushRecorder()
	recorder.blockFlush = 1
	recorder.flushBlocked = make(chan struct{placeholder)
	recorder.releaseFlush = releaseFirstFlush
	resultCh, errCh := runOpenAIResponseFlushTestAsync(recorder, reader, config.GatewayConfig{StreamDataIntervalTimeout: 30placeholder)

	waitOpenAIResponseFlushSignal(t, recorder.flushBlocked)
	close(allowBurst)
	waitOpenAIResponseFlushSignal(t, eofReached)
	close(releaseFirstFlush)

	require.NoError(t, <-errCh)
	require.NotNil(t, <-resultCh)
	gotBody, flushes := recorder.snapshot()
	require.Equal(t, first+burst, gotBody)
	require.Len(t, flushes, 2, "queued burst must remain batched until its drained event boundary")
	require.Equal(t, first, flushes[0])
	require.Equal(t, first+burst, flushes[1])
placeholder

func TestOpenAIResponseFlush_CommentAndEOFOnlyFlushCompleteResidual(t *testing.T) {
	body := "data: {\"type\":\"response.output_text.delta\",\"delta\":\"a\"placeholder\n\n" +
		": upstream-comment\n\n" +
		"data: [DONE]\n"
	recorder := newOpenAIResponseFlushRecorder()

	result, err := runOpenAIResponseFlushTest(recorder, io.NopCloser(strings.NewReader(body)), config.GatewayConfig{placeholder)

placeholder
	require.NotNil(t, result)
	gotBody, flushes := recorder.snapshot()
	require.Equal(t, body, gotBody)
	require.Len(t, flushes, 3)
	require.True(t, strings.HasSuffix(flushes[0], "\n\n"))
	require.True(t, strings.HasSuffix(flushes[1], "\n\n"))
	require.True(t, strings.HasSuffix(flushes[2], "data: [DONE]\n"), "EOF must flush only the remaining bytes")
placeholder

func TestOpenAIResponseFlush_TerminalReadErrorFlushesResidual(t *testing.T) {
	body := "data: [DONE]\n"
	recorder := newOpenAIResponseFlushRecorder()

	result, err := runOpenAIResponseFlushTest(recorder, &openAIResponseFlushReadError{payload: []byte(body)placeholder, config.GatewayConfig{placeholder)

placeholder
	require.NotNil(t, result)
	gotBody, flushes := recorder.snapshot()
	require.Equal(t, body, gotBody)
	require.Equal(t, []string{bodyplaceholder, flushes)
placeholder

func TestOpenAIResponseFlush_OutputWithoutTerminalFlushesResidualWithoutFailover(t *testing.T) {
	body := "data: {\"type\":\"response.output_text.delta\",\"delta\":\"partial\"placeholder\n"
	recorder := newOpenAIResponseFlushRecorder()

	result, err := runOpenAIResponseFlushTest(recorder, io.NopCloser(strings.NewReader(body)), config.GatewayConfig{placeholder)

	require.ErrorContains(t, err, "missing terminal event")
	var failoverErr *UpstreamFailoverError
	require.False(t, errors.As(err, &failoverErr))
	require.NotNil(t, result)
	gotBody, flushes := recorder.snapshot()
	require.Equal(t, body, gotBody)
	require.Equal(t, []string{bodyplaceholder, flushes)
placeholder

func TestOpenAIResponseFlush_PreambleWithoutTerminalRemainsBufferedForFailover(t *testing.T) {
	body := "data: {\"type\":\"response.created\",\"response\":{\"id\":\"resp_1\"placeholderplaceholder\n"
	recorder := newOpenAIResponseFlushRecorder()

	result, err := runOpenAIResponseFlushTest(recorder, io.NopCloser(strings.NewReader(body)), config.GatewayConfig{placeholder)

	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.NotNil(t, result)
	gotBody, flushes := recorder.snapshot()
	require.Empty(t, gotBody)
	require.Empty(t, flushes)
placeholder

func TestOpenAIResponseFlush_CanceledAfterOutputFlushesResidualWithoutErrorEvent(t *testing.T) {
	body := "data: {\"type\":\"response.output_text.delta\",\"delta\":\"partial\"placeholder\n"
	recorder := newOpenAIResponseFlushRecorder()

	result, err := runOpenAIResponseFlushTest(recorder, &openAIResponseFlushReadError{payload: []byte(body), err: context.Canceledplaceholder, config.GatewayConfig{placeholder)

	require.ErrorIs(t, err, context.Canceled)
	require.NotNil(t, result)
	gotBody, flushes := recorder.snapshot()
	require.Equal(t, body, gotBody)
	require.Equal(t, []string{bodyplaceholder, flushes)
	require.NotContains(t, gotBody, "stream_read_error")
placeholder

func TestOpenAIResponseFlush_KeepaliveFlushesImmediately(t *testing.T) {
	recorder := newOpenAIResponseFlushRecorder()
	reader, writer := io.Pipe()
	resultCh, errCh := runOpenAIResponseFlushTestAsync(recorder, reader, config.GatewayConfig{StreamKeepaliveInterval: 1placeholder)

	waitOpenAIResponseFlushCount(t, recorder, 1)
	_, flushes := recorder.snapshot()
	require.Equal(t, ":\n\n", flushes[0])
	_, err := writer.Write([]byte("data: [DONE]\n\n"))
placeholder
	require.NoError(t, writer.Close())

	require.NoError(t, <-errCh)
	require.NotNil(t, <-resultCh)
	gotBody, flushes := recorder.snapshot()
	require.Equal(t, ":\n\ndata: [DONE]\n\n", gotBody)
	require.Len(t, flushes, 2)
placeholder

func TestOpenAIResponseFlush_KeepaliveDoesNotSplitOpenEvent(t *testing.T) {
	const dataLine = `data: {"type":"response.output_text.delta","delta":"a"placeholder`
	// Filling the 16-slot scan queue proves the main loop processed data before the reader reaches the gated blank.
	dataLines := make([]string, 17)
	for i := range dataLines {
		dataLines[i] = dataLine
placeholder
	partialEvent := strings.Join(dataLines, "\n") + "\n"
	completeEvent := partialEvent + "\n"
	terminal := "data: [DONE]\n\n"
	allowBlank := make(chan struct{placeholder)
	allowTerminal := make(chan struct{placeholder)
	blankWaiting := make(chan struct{placeholder)
	terminalWaiting := make(chan struct{placeholder)
	reader := &stagedOpenAISSEReadCloser{
		segments: [][]byte{[]byte(partialEvent), []byte("\n"), []byte(terminal)placeholder,
		gates:    []<-chan struct{placeholder{nil, allowBlank, allowTerminalplaceholder,
		waiting:  []chan struct{placeholder{nil, blankWaiting, terminalWaitingplaceholder,
placeholder
	recorder := newOpenAIResponseFlushRecorder()
	resultCh, errCh := runOpenAIResponseFlushTestAsync(recorder, reader, config.GatewayConfig{StreamKeepaliveInterval: 1placeholder)

	waitOpenAIResponseFlushSignal(t, blankWaiting)
	timer := time.NewTimer(1250 * time.Millisecond)
	select {
	case count := <-recorder.flushEvents:
		timer.Stop()
		t.Fatalf("keepalive flushed open event before its blank boundary: flush %d", count)
	case <-timer.C:
placeholder

	close(allowBlank)
	waitOpenAIResponseFlushSignal(t, terminalWaiting)
	waitOpenAIResponseFlushCount(t, recorder, 1)
	gotBody, flushes := recorder.snapshot()
	require.Equal(t, completeEvent, gotBody)
	require.Equal(t, []string{completeEventplaceholder, flushes)

	close(allowTerminal)
	require.NoError(t, <-errCh)
	require.NotNil(t, <-resultCh)
	gotBody, flushes = recorder.snapshot()
	require.Equal(t, completeEvent+terminal, gotBody)
	require.Len(t, flushes, 2)
	require.Equal(t, completeEvent+terminal, flushes[1])
placeholder

func TestOpenAIResponseFlush_FailedAndErrorEventsFlushAtBoundaries(t *testing.T) {
	t.Run("failed at EOF", func(t *testing.T) {
		body := "data: {\"type\":\"response.output_text.delta\",\"delta\":\"a\"placeholder\n\n" +
			"data: {\"type\":\"response.failed\",\"response\":{\"error\":{\"code\":\"safety_error\",\"message\":\"blocked\"placeholder,\"usage\":{\"input_tokens\":3,\"output_tokens\":1placeholderplaceholderplaceholder\n"
		recorder := newOpenAIResponseFlushRecorder()

		result, err := runOpenAIResponseFlushTest(recorder, io.NopCloser(strings.NewReader(body)), config.GatewayConfig{placeholder)

	placeholder
		require.NotNil(t, result)
		require.Equal(t, 3, result.usage.InputTokens)
		gotBody, flushes := recorder.snapshot()
		expectedBody := "data: {\"type\":\"response.output_text.delta\",\"delta\":\"a\"placeholder\n\n" +
			"data: {\"type\":\"response.failed\",\"response\":{\"error\":{\"code\":\"safety_error\",\"message\":\"blocked\"placeholderplaceholderplaceholder\n"
		require.Equal(t, expectedBody, gotBody)
		require.Len(t, flushes, 2)
		require.Contains(t, flushes[1], "response.failed")
placeholder)

	t.Run("error event", func(t *testing.T) {
		body := "data: {\"type\":\"error\",\"error\":{\"message\":\"failed\"placeholderplaceholder\n\n" +
			"data: [DONE]\n\n"
		recorder := newOpenAIResponseFlushRecorder()

		result, err := runOpenAIResponseFlushTest(recorder, io.NopCloser(strings.NewReader(body)), config.GatewayConfig{placeholder)

	placeholder
		require.NotNil(t, result)
		gotBody, flushes := recorder.snapshot()
		require.Equal(t, body, gotBody)
		require.Len(t, flushes, 2)
placeholder)
placeholder

func TestOpenAIResponseFlush_ClientDisconnectStillDrainsUsage(t *testing.T) {
	first := "data: {\"type\":\"response.output_text.delta\",\"delta\":\"a\"placeholder\n\n"
	terminal := "data: {\"type\":\"response.completed\",\"response\":{\"status\":\"completed\",\"output\":[],\"usage\":{\"input_tokens\":7,\"output_tokens\":5,\"input_tokens_details\":{\"cached_tokens\":2placeholderplaceholderplaceholderplaceholder\n\n"
	recorder := newOpenAIResponseFlushRecorder()
	recorder.failAfterWrites = 1

	result, err := runOpenAIResponseFlushTest(recorder, io.NopCloser(strings.NewReader(first+terminal)), config.GatewayConfig{placeholder)

placeholder
	require.NotNil(t, result)
	require.Equal(t, 7, result.usage.InputTokens)
	require.Equal(t, 5, result.usage.OutputTokens)
	require.Equal(t, 2, result.usage.CacheReadInputTokens)
	gotBody, flushes := recorder.snapshot()
	require.Equal(t, first, gotBody)
	require.Len(t, flushes, 1)
placeholder

func runOpenAIResponseFlushTest(recorder *openAIResponseFlushRecorder, body io.ReadCloser, gatewayCfg config.GatewayConfig) (*openaiStreamingResult, error) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	svc := &OpenAIGatewayService{
		cfg:           &config.Config{Gateway: gatewayCfgplaceholder,
		toolCorrector: NewCodexToolCorrector(),
placeholder
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       body,
placeholder
	return svc.handleStreamingResponse(context.Background(), resp, c, &Account{ID: 1, Platform: PlatformOpenAIplaceholder, time.Now(), "gpt-5", "gpt-5")
placeholder

func runOpenAIResponseFlushTestAsync(recorder *openAIResponseFlushRecorder, body io.ReadCloser, gatewayCfg config.GatewayConfig) (<-chan *openaiStreamingResult, <-chan error) {
	resultCh := make(chan *openaiStreamingResult, 1)
	errCh := make(chan error, 1)
	go func() {
		result, err := runOpenAIResponseFlushTest(recorder, body, gatewayCfg)
		resultCh <- result
		errCh <- err
placeholder()
	return resultCh, errCh
placeholder

func waitOpenAIResponseFlushCount(t *testing.T, recorder *openAIResponseFlushRecorder, want int) {
placeholder
	timer := time.NewTimer(3 * time.Second)
	defer timer.Stop()
	for {
		select {
		case count := <-recorder.flushEvents:
			if count >= want {
				return
		placeholder
		case <-timer.C:
			t.Fatalf("timed out waiting for flush %d", want)
	placeholder
placeholder
placeholder

func waitOpenAIResponseFlushSignal(t *testing.T, signal <-chan struct{placeholder) {
placeholder
	select {
	case <-signal:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for stream signal")
placeholder
placeholder
