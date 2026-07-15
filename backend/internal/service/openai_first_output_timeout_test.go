package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type blockingOpenAIResponseHeaderUpstream struct {
	canceled chan struct{placeholder
	once     sync.Once
placeholder

type firstOutputCloseTrackingBody struct {
	io.ReadCloser
	closed chan struct{placeholder
	once   sync.Once
placeholder

func (b *firstOutputCloseTrackingBody) Close() error {
	b.once.Do(func() { close(b.closed) placeholder)
	return b.ReadCloser.Close()
placeholder

func (u *blockingOpenAIResponseHeaderUpstream) Do(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
	select {
	case <-req.Context().Done():
		u.once.Do(func() { close(u.canceled) placeholder)
		return nil, req.Context().Err()
	case <-time.After(1500 * time.Millisecond):
		return nil, errors.New("test upstream was not canceled before response headers")
placeholder
placeholder

func (u *blockingOpenAIResponseHeaderUpstream) DoWithTLS(req *http.Request, _ string, _ int64, _ int, _ *tlsfingerprint.Profile) (*http.Response, error) {
	return u.Do(req, "", 0, 0)
placeholder

func TestOpenAIForwardFirstOutputTimeoutIncludesResponseHeaderWait(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstream := &blockingOpenAIResponseHeaderUpstream{canceled: make(chan struct{placeholder)placeholder
	svc := &OpenAIGatewayService{
		cfg: &config.Config{Gateway: config.GatewayConfig{
			OpenAIFirstOutputTimeoutSeconds: 1,
			MaxLineSize:                     defaultMaxLineSize,
placeholder
		httpUpstream: upstream,
placeholder
	body := []byte(`{"model":"gpt-5.5","stream":true,"reasoning":{"effort":"low"placeholder,"input":"hello"placeholder`)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))
	account := &Account{
		ID: 1, Name: "oauth-test", Platform: PlatformOpenAI, Type: AccountTypeOAuth,
		Status: StatusActive, Schedulable: true, Concurrency: 1,
placeholder"access_token": "test-token", "chatgpt_account_id": "test-account"placeholder,
placeholder

	started := time.Now()
	_, err := svc.Forward(context.Background(), c, account, body)

placeholder
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusGatewayTimeout, failoverErr.StatusCode)
	require.Contains(t, string(failoverErr.ResponseBody), "first_output_timeout")
	require.True(t, failoverErr.SafeToFailoverAfterWrite)
	require.Less(t, time.Since(started), 1300*time.Millisecond)
	require.Empty(t, rec.Body.String())
	select {
	case <-upstream.canceled:
	default:
		t.Fatal("response-header timeout did not cancel the upstream request context")
placeholder
placeholder

func TestOpenAINativeFirstOutputTimeoutDisabledPreservesSynchronousStream(t *testing.T) {
	svc := &OpenAIGatewayService{cfg: &config.Config{Gateway: config.GatewayConfig{
		OpenAIFirstOutputTimeoutSeconds: 0,
		MaxLineSize:                     defaultMaxLineSize,
placeholderplaceholderplaceholder
	resp := &http.Response{StatusCode: http.StatusOK, Header: http.Header{placeholder, Body: io.NopCloser(strings.NewReader(strings.Join([]string{
		`data: {"type":"response.created","response":{"id":"resp_disabled"placeholderplaceholder`,
		"",
		`data: {"type":"response.completed","response":{"id":"resp_disabled","usage":{"input_tokens":1,"output_tokens":1placeholderplaceholderplaceholder`,
		"",
placeholder, "\n")))placeholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)

	result, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAIplaceholder, time.Now(), "model", "model")

placeholder
	require.NotNil(t, result)
	require.Contains(t, rec.Body.String(), "response.completed")
placeholder

func TestOpenAINativeFirstOutputTimeoutIgnoresPreambleAndCleansReader(t *testing.T) {
	svc := &OpenAIGatewayService{cfg: &config.Config{Gateway: config.GatewayConfig{
		OpenAIFirstOutputTimeoutSeconds: 1,
		MaxLineSize:                     defaultMaxLineSize,
placeholderplaceholderplaceholder
	pr, pw := io.Pipe()
	writerDone := make(chan struct{placeholder)
	go func() {
		defer close(writerDone)
		defer func() { _ = pw.Close() placeholder()
		_, _ = pw.Write([]byte("data: {\"type\":\"response.created\",\"response\":{\"id\":\"resp_slow\"placeholderplaceholder\n\n"))
		_, _ = pw.Write([]byte("data: {\"type\":\"response.in_progress\",\"response\":{\"id\":\"resp_slow\"placeholderplaceholder\n\n"))
		time.Sleep(200 * time.Millisecond)
placeholder()
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	body := &firstOutputCloseTrackingBody{ReadCloser: pr, closed: make(chan struct{placeholder)placeholder
	resp := &http.Response{StatusCode: http.StatusOK, Header: http.Header{placeholder, Body: bodyplaceholder

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAIplaceholder, time.Now().Add(-2*time.Second), "model", "model")

placeholder
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusGatewayTimeout, failoverErr.StatusCode)
	require.Contains(t, string(failoverErr.ResponseBody), "first_output_timeout")
	require.True(t, failoverErr.SafeToFailoverAfterWrite)
	require.Empty(t, rec.Body.String())
	select {
	case <-body.closed:
	default:
		t.Fatal("first-output timeout did not close the upstream response body")
placeholder
	select {
	case <-writerDone:
	case <-time.After(time.Second):
		t.Fatal("stream reader/writer goroutine did not exit after first-output timeout")
placeholder
placeholder

func TestOpenAIFirstOutputTimeoutForReasoningEffort(t *testing.T) {
	svc := &OpenAIGatewayService{cfg: &config.Config{Gateway: config.GatewayConfig{
		OpenAIFirstOutputTimeoutSeconds:           120,
		OpenAIHighEffortFirstOutputTimeoutSeconds: 300,
placeholderplaceholderplaceholder

	require.Equal(t, 120*time.Second, svc.openAIFirstOutputTimeout("low"))
	require.Equal(t, 300*time.Second, svc.openAIFirstOutputTimeout("high"))
	require.Equal(t, 300*time.Second, svc.openAIFirstOutputTimeout("xhigh"))
	require.Equal(t, 300*time.Second, svc.openAIFirstOutputTimeout("max"))
placeholder

func TestOpenAIFirstOutputStageDefaultLimitIsIndependentFromScannerLimit(t *testing.T) {
	stage := newDefaultOpenAIFirstOutputStage()
	defer func() { require.NoError(t, stage.Close()) placeholder()

	require.EqualValues(t, 8*1024*1024, stage.limit)
	require.Greater(t, stage.limit, int64(68106))
	require.Less(t, stage.limit, int64(defaultMaxLineSize))
placeholder

func TestOpenAIFirstOutputEventQueueSizeBackpressuresGuardedStreams(t *testing.T) {
	require.Equal(t, 1, openAIFirstOutputEventQueueSize(true))
	require.Equal(t, 16, openAIFirstOutputEventQueueSize(false))
placeholder

func TestOpenAIFirstOutputDynamicScannerLimitsOnlyWhileGuardIsActive(t *testing.T) {
	var guardActive atomic.Bool
	guardActive.Store(true)
	split := openAIFirstOutputDynamicScanLines(&guardActive)
	guardLimit := openAIFirstOutputStageMaxBytes + openAIFirstOutputScannerFramingAllowance
	undelimited := bytes.Repeat([]byte("x"), guardLimit)

	_, _, err := split(undelimited, false)
	require.ErrorIs(t, err, errOpenAIFirstOutputScannerLimit)

	guardActive.Store(false)
	advance, token, err := split(undelimited, false)
placeholder
	require.Zero(t, advance)
	require.Nil(t, token)
placeholder

func TestOpenAIFirstOutputStageOverflowIsAtomicAndCleanupRemovesSpool(t *testing.T) {
	stage := newOpenAIFirstOutputStage(70 * 1024)
	payload := bytes.Repeat([]byte("x"), 68*1024)
	n, err := stage.Write(payload)
placeholder
	require.Equal(t, len(payload), n)
	if runtime.GOOS == "windows" {
		require.Nil(t, stage.tempFile)
		require.Empty(t, stage.tempPath)
placeholder else {
		require.NotNil(t, stage.tempFile)
		require.NotEmpty(t, stage.tempPath)
		_, err = os.Stat(stage.tempPath)
		require.ErrorIs(t, err, os.ErrNotExist)
		stat, statErr := stage.tempFile.Stat()
		require.NoError(t, statErr)
		require.Equal(t, os.FileMode(0o600), stat.Mode().Perm())
placeholder

	n, err = stage.Write(bytes.Repeat([]byte("y"), 3*1024))
	require.Zero(t, n)
	require.ErrorIs(t, err, errOpenAIFirstOutputStageLimit)
	require.EqualValues(t, len(payload), stage.Buffered())
	path := stage.tempPath
	require.NoError(t, stage.Close())
	require.True(t, stage.closed)
	require.Nil(t, stage.tempFile)
	require.Empty(t, stage.tempPath)
	if path != "" {
		_, err = os.Stat(path)
		require.ErrorIs(t, err, os.ErrNotExist)
placeholder
placeholder

func TestOpenAIFirstOutputStageCommitCopiesSpoolAndRemovesTemp(t *testing.T) {
	stage := newOpenAIFirstOutputStage(80 * 1024)
	payload := bytes.Repeat([]byte("z"), 68*1024)
	_, err := stage.Write(payload)
placeholder
	path := stage.tempPath
	if runtime.GOOS == "windows" {
		require.Empty(t, path)
		require.Nil(t, stage.tempFile)
placeholder else {
		require.NotEmpty(t, path)
		require.NotNil(t, stage.tempFile)
		_, statErr := os.Stat(path)
		require.ErrorIs(t, statErr, os.ErrNotExist)
placeholder

	var downstream bytes.Buffer
	require.NoError(t, stage.CommitTo(&downstream))
	require.Equal(t, payload, downstream.Bytes())
	require.Zero(t, stage.Buffered())
	if path != "" {
		_, err = os.Stat(path)
		require.ErrorIs(t, err, os.ErrNotExist)
placeholder
	require.NoError(t, stage.Close())
placeholder

func TestOpenAIFirstOutputStageUnlinkFailurePermanentlyFallsBackToMemoryAndRetriesCleanup(t *testing.T) {
	stage := newDefaultOpenAIFirstOutputStage()
	stage.memoryOnly = false
	t.Cleanup(func() {
		stage.removeFile = os.Remove
		_ = stage.Close()
placeholder)
	createCalls := 0
	stage.createTemp = func() (*os.File, error) {
		createCalls++
		return os.CreateTemp("", "sub2api-openai-first-output-fallback-*")
placeholder
	removeCalls := 0
	stage.removeFile = func(path string) error {
		removeCalls++
		if removeCalls <= 2 {
			return errors.New("forced remove failure")
	placeholder
		return os.Remove(path)
placeholder

	payload := bytes.Repeat([]byte("m"), 68*1024)
	_, err := stage.Write(payload)
placeholder
	require.True(t, stage.memoryOnly)
	require.Nil(t, stage.tempFile)
	require.NotEmpty(t, stage.tempPath)
	require.Equal(t, 1, createCalls)
	stat, statErr := os.Stat(stage.tempPath)
	require.NoError(t, statErr)
	require.Zero(t, stat.Size(), "failed-unlink fallback must never write plaintext to the named file")

	_, err = stage.WriteString("more")
placeholder
	require.Equal(t, 1, createCalls, "memory-only fallback must not retry CreateTemp")
	path := stage.tempPath
	cleanupErr := stage.Close()
	require.ErrorContains(t, cleanupErr, "forced remove failure")
	require.Empty(t, stage.tempPath)
	_, err = os.Stat(path)
	require.ErrorIs(t, err, os.ErrNotExist)
	require.NoError(t, stage.Close())
placeholder

func TestOpenAINativeFirstOutputTimeoutDisarmsAfterSemanticOutput(t *testing.T) {
	cfg := &config.Config{Gateway: config.GatewayConfig{
		OpenAIFirstOutputTimeoutSeconds: 1,
		MaxLineSize:                     defaultMaxLineSize,
placeholderplaceholder
	svc := &OpenAIGatewayService{cfg: cfg, responseHeaderFilter: compileResponseHeaderFilter(cfg)placeholder
	pr, pw := io.Pipe()
	go func() {
		defer func() { _ = pw.Close() placeholder()
		_, _ = pw.Write([]byte("data: {\"type\":\"response.created\",\"response\":{\"id\":\"resp_ok\"placeholderplaceholder\n\n"))
		_, _ = pw.Write([]byte("data: {\"type\":\"response.output_text.delta\",\"delta\":\"hello\"placeholder\n\n"))
		time.Sleep(1100 * time.Millisecond)
		_, _ = pw.Write([]byte("data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_ok\",\"usage\":{\"input_tokens\":1,\"output_tokens\":1placeholderplaceholderplaceholder\n\n"))
placeholder()
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	resp := &http.Response{StatusCode: http.StatusOK, Header: http.Header{
		"X-Request-Id":                   []string{"request-winning"placeholder,
		"X-Ratelimit-Remaining-Requests": []string{"42"placeholder,
placeholder, Body: prplaceholder

	result, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAIplaceholder, time.Now(), "model", "model")

placeholder
	require.NotNil(t, result)
	require.NotNil(t, result.firstTokenMs)
	require.Contains(t, rec.Body.String(), "response.output_text.delta")
	require.Contains(t, rec.Body.String(), "response.completed")
	require.Equal(t, "request-winning", rec.Result().Header.Get("X-Request-Id"))
	require.Equal(t, "42", rec.Result().Header.Get("X-Ratelimit-Remaining-Requests"))
placeholder

func TestOpenAINativeFirstOutputTimeoutWaitsForCompleteSemanticEvent(t *testing.T) {
	const lineSize = 68106
	prefix := `data: {"type":"response.output_text.delta","delta":"`
	suffix := `"placeholder`
	line := prefix + strings.Repeat("x", lineSize-len(prefix)-len(suffix)) + suffix
	require.Len(t, line, lineSize)
	assertOpenAINativeLargeOpenEventTimesOutWithoutLeak(t, line)
placeholder

func TestOpenAINativeFirstOutputTimeoutDoesNotLeakLargePreambleEvent(t *testing.T) {
	const lineSize = 68106
	prefix := `data: {"type":"response.created","response":{"id":"resp_partial","padding":"`
	suffix := `"placeholderplaceholder`
	line := prefix + strings.Repeat("x", lineSize-len(prefix)-len(suffix)) + suffix
	require.Len(t, line, lineSize)
	assertOpenAINativeLargeOpenEventTimesOutWithoutLeak(t, line)
placeholder

func assertOpenAINativeLargeOpenEventTimesOutWithoutLeak(t *testing.T, line string) {
placeholder
	cfg := &config.Config{Gateway: config.GatewayConfig{
		OpenAIFirstOutputTimeoutSeconds: 1,
		StreamKeepaliveInterval:         1,
		MaxLineSize:                     defaultMaxLineSize,
placeholderplaceholder
	svc := &OpenAIGatewayService{cfg: cfg, responseHeaderFilter: compileResponseHeaderFilter(cfg)placeholder
	pr, pw := io.Pipe()
	body := &firstOutputCloseTrackingBody{ReadCloser: pr, closed: make(chan struct{placeholder)placeholder
	writerDone := make(chan struct{placeholder)
	go func() {
		defer close(writerDone)
		defer func() { _ = pw.Close() placeholder()
		_, _ = pw.Write([]byte(line + "\n"))
		select {
		case <-body.closed:
		case <-time.After(2 * time.Second):
	placeholder
placeholder()
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	resp := &http.Response{StatusCode: http.StatusOK, Header: http.Header{
		"X-Request-Id":                   []string{"request-partial"placeholder,
		"X-Ratelimit-Remaining-Requests": []string{"1"placeholder,
placeholder, Body: bodyplaceholder

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAIplaceholder, time.Now(), "model", "model")

	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusGatewayTimeout, failoverErr.StatusCode)
	require.Contains(t, string(failoverErr.ResponseBody), "first_output_timeout")
	require.True(t, failoverErr.SafeToFailoverAfterWrite)
	require.NotContains(t, rec.Body.String(), "data:", "attempt JSON must remain private before the SSE boundary")
	require.NotContains(t, rec.Body.String(), `"type"`, "attempt JSON must remain private before the SSE boundary")
	for _, outputLine := range strings.Split(strings.TrimSpace(rec.Body.String()), "\n") {
		if outputLine != "" {
			require.True(t, strings.HasPrefix(outputLine, ":"), "only keepalive comments may precede failover: %q", outputLine)
	placeholder
placeholder
	require.Empty(t, rec.Header().Values("X-Request-Id"))
	require.Empty(t, rec.Header().Values("X-Ratelimit-Remaining-Requests"))
	select {
	case <-writerDone:
	case <-time.After(time.Second):
		t.Fatal("partial-event writer did not exit after timeout closed the body")
placeholder
placeholder

func TestOpenAINativeFirstOutputEOFDispatchesTerminalEventWithoutBlankLine(t *testing.T) {
	cfg := &config.Config{Gateway: config.GatewayConfig{
		OpenAIFirstOutputTimeoutSeconds: 1,
		MaxLineSize:                     defaultMaxLineSize,
placeholderplaceholder
	svc := &OpenAIGatewayService{cfg: cfg, responseHeaderFilter: compileResponseHeaderFilter(cfg)placeholder
	payload := `data: {"type":"response.completed","response":{"id":"resp_eof","usage":{"input_tokens":3,"output_tokens":2placeholderplaceholderplaceholder`
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"X-Request-Id":                   []string{"request-eof"placeholder,
			"X-Ratelimit-Remaining-Requests": []string{"17"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(payload)),
placeholder

	result, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAIplaceholder, time.Now(), "model", "model")

placeholder
	require.NotNil(t, result)
	require.NotNil(t, result.firstTokenMs)
	require.Equal(t, "resp_eof", result.responseID)
	require.Equal(t, 3, result.usage.InputTokens)
	require.Equal(t, 2, result.usage.OutputTokens)
	require.Contains(t, rec.Body.String(), `"type":"response.completed"`)
	require.Contains(t, rec.Body.String(), `"id":"resp_eof"`)
	require.True(t, strings.HasSuffix(rec.Body.String(), "\n"))
	require.False(t, strings.HasSuffix(rec.Body.String(), "\n\n"), "EOF dispatch must not synthesize a blank line")
	require.Equal(t, "request-eof", rec.Result().Header.Get("X-Request-Id"))
	require.Equal(t, "17", rec.Result().Header.Get("X-Ratelimit-Remaining-Requests"))
placeholder

func TestOpenAINativeFirstOutputStageOverflowFailsOverWithoutAttemptBytes(t *testing.T) {
	cfg := &config.Config{Gateway: config.GatewayConfig{
		OpenAIFirstOutputTimeoutSeconds: 30,
		MaxLineSize:                     2 * 1024 * 1024,
placeholderplaceholder
	svc := &OpenAIGatewayService{cfg: cfg, responseHeaderFilter: compileResponseHeaderFilter(cfg)placeholder
	const lineSize = 1024*1024 - 256
	prefix := `data: {"type":"response.output_text.delta","delta":"`
	suffix := `"placeholder`
	line := prefix + strings.Repeat("x", lineSize-len(prefix)-len(suffix)) + suffix
	body := strings.Repeat(line+"\n", 9)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"X-Request-Id":                   []string{"request-overflow"placeholder,
			"X-Ratelimit-Remaining-Requests": []string{"1"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(body)),
placeholder

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAIplaceholder, time.Now(), "model", "model")

	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusBadGateway, failoverErr.StatusCode)
	require.True(t, failoverErr.SafeToFailoverAfterWrite)
	require.Contains(t, string(failoverErr.ResponseBody), "staging limit exceeded")
	require.Empty(t, rec.Body.String())
	require.Empty(t, rec.Header().Values("X-Request-Id"))
	require.Empty(t, rec.Header().Values("X-Ratelimit-Remaining-Requests"))
placeholder

func TestOpenAINativeFirstOutputScannerRejectsOversizedLineWithoutLeak(t *testing.T) {
	cfg := &config.Config{Gateway: config.GatewayConfig{
		OpenAIFirstOutputTimeoutSeconds: 30,
		MaxLineSize:                     defaultMaxLineSize,
placeholderplaceholder
	svc := &OpenAIGatewayService{cfg: cfg, responseHeaderFilter: compileResponseHeaderFilter(cfg)placeholder
	oversizedLine := "data: " + strings.Repeat("x", openAIFirstOutputStageMaxBytes+openAIFirstOutputScannerFramingAllowance+1024)
	body := "data: {\"type\":\"response.created\",\"response\":{\"id\":\"resp_private\"placeholderplaceholder\n\n" + oversizedLine + "\n"
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"X-Request-Id":                   []string{"request-too-large"placeholder,
			"X-Ratelimit-Remaining-Requests": []string{"1"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(body)),
placeholder

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAIplaceholder, time.Now(), "model", "model")

	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusBadGateway, failoverErr.StatusCode)
	require.True(t, failoverErr.SafeToFailoverAfterWrite)
	require.Contains(t, string(failoverErr.ResponseBody), "line exceeds guarded first-output limit")
	require.Empty(t, rec.Body.String())
	require.Empty(t, rec.Header().Values("X-Request-Id"))
	require.Empty(t, rec.Header().Values("X-Ratelimit-Remaining-Requests"))
placeholder

func TestOpenAINativeFirstOutputScannerAllowsLargeEventAfterSemanticBoundary(t *testing.T) {
	cfg := &config.Config{Gateway: config.GatewayConfig{
		OpenAIFirstOutputTimeoutSeconds: 30,
		MaxLineSize:                     defaultMaxLineSize,
placeholderplaceholder
	svc := &OpenAIGatewayService{cfg: cfg, responseHeaderFilter: compileResponseHeaderFilter(cfg)placeholder
	largeDelta := strings.Repeat("i", openAIFirstOutputStageMaxBytes+openAIFirstOutputScannerFramingAllowance+1024)
	body := strings.Join([]string{
		`data: {"type":"response.output_text.delta","delta":"ready"placeholder`,
		"",
		`data: {"type":"response.output_text.delta","delta":"` + largeDelta + `"placeholder`,
		"",
		`data: {"type":"response.completed","response":{"id":"resp_large_image","usage":{"input_tokens":4,"output_tokens":3placeholderplaceholderplaceholder`,
		"",
placeholder, "\n")
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"X-Request-Id": []string{"request-large-image"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(body)),
placeholder

	result, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAIplaceholder, time.Now(), "model", "model")

placeholder
	require.NotNil(t, result)
	require.NotNil(t, result.firstTokenMs)
	require.Equal(t, "resp_large_image", result.responseID)
	require.Equal(t, 4, result.usage.InputTokens)
	require.Equal(t, 3, result.usage.OutputTokens)
	require.Contains(t, rec.Body.String(), `"delta":"ready"`)
	require.Contains(t, rec.Body.String(), `"id":"resp_large_image"`)
	require.Contains(t, rec.Body.String(), strings.Repeat("i", 1024))
	require.Equal(t, "request-large-image", rec.Result().Header.Get("X-Request-Id"))
placeholder

func TestOpenAINativeFirstOutputTimeoutDisabledPreservesKeepaliveFlush(t *testing.T) {
	svc := &OpenAIGatewayService{cfg: &config.Config{Gateway: config.GatewayConfig{
		StreamKeepaliveInterval: 1,
		MaxLineSize:             defaultMaxLineSize,
placeholderplaceholderplaceholder
	pr, pw := io.Pipe()
	go func() {
		defer func() { _ = pw.Close() placeholder()
		_, _ = pw.Write([]byte("data: {\"type\":\"response.created\",\"response\":{\"id\":\"resp_stalled\"placeholderplaceholder\n\n"))
		_, _ = pw.Write([]byte("data: {\"type\":\"response.in_progress\",\"response\":{\"id\":\"resp_stalled\"placeholderplaceholder\n\n"))
		time.Sleep(1100 * time.Millisecond)
placeholder()
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	resp := &http.Response{StatusCode: http.StatusOK, Header: http.Header{placeholder, Body: prplaceholder

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAIplaceholder, time.Now(), "model", "model")

placeholder
	require.Contains(t, rec.Body.String(), ":\n\n")
	require.Contains(t, rec.Body.String(), "response.created")
	require.Contains(t, rec.Body.String(), "response.in_progress")
placeholder

func TestOpenAINativeFirstOutputFailoverKeepsAttemptHeadersPrivateAfterKeepaliveCommit(t *testing.T) {
	cfg := &config.Config{Gateway: config.GatewayConfig{
		OpenAIFirstOutputTimeoutSeconds: 2,
		StreamKeepaliveInterval:         1,
		MaxLineSize:                     defaultMaxLineSize,
placeholderplaceholder
	svc := &OpenAIGatewayService{
		cfg:                  cfg,
		responseHeaderFilter: compileResponseHeaderFilter(cfg),
placeholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	firstBody, firstWriter := io.Pipe()
	trackedFirstBody := &firstOutputCloseTrackingBody{ReadCloser: firstBody, closed: make(chan struct{placeholder)placeholder
	firstWriterDone := make(chan struct{placeholder)
	go func() {
		defer close(firstWriterDone)
		defer func() { _ = firstWriter.Close() placeholder()
		_, _ = firstWriter.Write([]byte("data: {\"type\":\"response.created\",\"response\":{\"id\":\"resp_first\"placeholderplaceholder\n\n"))
		select {
		case <-trackedFirstBody.closed:
		case <-time.After(4 * time.Second):
	placeholder
placeholder()
	firstResp := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type":                   []string{"text/event-stream"placeholder,
			"X-Request-Id":                   []string{"request-first"placeholder,
			"X-Ratelimit-Remaining-Requests": []string{"1"placeholder,
	placeholder,
		Body: trackedFirstBody,
placeholder

	_, firstErr := svc.handleStreamingResponse(c.Request.Context(), firstResp, c, &Account{ID: 1, Platform: PlatformOpenAIplaceholder, time.Now(), "model", "model")
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, firstErr, &failoverErr)
	require.Contains(t, rec.Body.String(), ":\n\n", "first attempt should have committed only a stable keepalive")
	require.NotContains(t, rec.Body.String(), "resp_first")

	secondResp := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type":                   []string{"text/event-stream"placeholder,
			"X-Request-Id":                   []string{"request-second"placeholder,
			"X-Ratelimit-Remaining-Requests": []string{"99"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(strings.Join([]string{
			`data: {"type":"response.output_text.delta","delta":"hello"placeholder`,
			"",
			`data: {"type":"response.completed","response":{"id":"resp_second","usage":{"input_tokens":1,"output_tokens":1placeholderplaceholderplaceholder`,
			"",
	placeholder, "\n"))),
placeholder
	result, secondErr := svc.handleStreamingResponse(c.Request.Context(), secondResp, c, &Account{ID: 2, Platform: PlatformOpenAIplaceholder, time.Now(), "model", "model")

	require.NoError(t, secondErr)
	require.NotNil(t, result)
	require.Contains(t, rec.Body.String(), "resp_second")
	wireHeaders := rec.Result().Header
	require.Empty(t, wireHeaders.Values("X-Request-Id"))
	require.Empty(t, wireHeaders.Values("X-Ratelimit-Remaining-Requests"))
	require.Empty(t, rec.Header().Values("X-Request-Id"))
	require.Empty(t, rec.Header().Values("X-Ratelimit-Remaining-Requests"))
	select {
	case <-firstWriterDone:
	case <-time.After(time.Second):
		t.Fatal("first account writer did not exit after timeout")
placeholder
placeholder
