package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/apicompat"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

type openAICompatFailingWriter struct {
	gin.ResponseWriter
	failAfter int
	writes    int
placeholder

func (w *openAICompatFailingWriter) Write(p []byte) (int, error) {
	if w.writes >= w.failAfter {
		return 0, errors.New("write failed: client disconnected")
placeholder
	w.writes++
	return w.ResponseWriter.Write(p)
placeholder

type openAICompatBlockingReadCloser struct {
	data      []byte
	offset    int
	closed    chan struct{placeholder
	closeOnce sync.Once
placeholder

func newOpenAICompatBlockingReadCloser(data []byte) *openAICompatBlockingReadCloser {
	return &openAICompatBlockingReadCloser{
		data:   data,
		closed: make(chan struct{placeholder),
placeholder
placeholder

func (r *openAICompatBlockingReadCloser) Read(p []byte) (int, error) {
	if r.offset < len(r.data) {
		n := copy(p, r.data[r.offset:])
		r.offset += n
		return n, nil
placeholder
	<-r.closed
	return 0, io.EOF
placeholder

func (r *openAICompatBlockingReadCloser) Close() error {
	r.closeOnce.Do(func() {
		close(r.closed)
placeholder)
	return nil
placeholder

func TestNormalizeOpenAICompatRequestedModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
placeholder{
		{name: "gpt reasoning alias strips xhigh", input: "gpt-5.4-xhigh", want: "gpt-5.4"placeholder,
		{name: "gpt reasoning alias strips none", input: "gpt-5.4-none", want: "gpt-5.4"placeholder,
		{name: "codex max model stays intact", input: "gpt-5.1-codex-max", want: "gpt-5.1-codex-max"placeholder,
		{name: "non openai model unchanged", input: "claude-opus-4-6", want: "claude-opus-4-6"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, NormalizeOpenAICompatRequestedModel(tt.input))
	placeholder)
placeholder
placeholder

func TestApplyOpenAICompatModelNormalization(t *testing.T) {
	t.Parallel()

	t.Run("derives xhigh from model suffix when output config missing", func(t *testing.T) {
		req := &apicompat.AnthropicRequest{Model: "gpt-5.4-xhigh"placeholder

		applyOpenAICompatModelNormalization(req)

		require.Equal(t, "gpt-5.4", req.Model)
		require.NotNil(t, req.OutputConfig)
		require.Equal(t, "max", req.OutputConfig.Effort)
placeholder)

	t.Run("explicit output config wins over model suffix", func(t *testing.T) {
		req := &apicompat.AnthropicRequest{
			Model:        "gpt-5.4-xhigh",
			OutputConfig: &apicompat.AnthropicOutputConfig{Effort: "low"placeholder,
	placeholder

		applyOpenAICompatModelNormalization(req)

		require.Equal(t, "gpt-5.4", req.Model)
		require.NotNil(t, req.OutputConfig)
		require.Equal(t, "low", req.OutputConfig.Effort)
placeholder)

	t.Run("non openai model is untouched", func(t *testing.T) {
		req := &apicompat.AnthropicRequest{Model: "claude-opus-4-6"placeholder

		applyOpenAICompatModelNormalization(req)

		require.Equal(t, "claude-opus-4-6", req.Model)
		require.Nil(t, req.OutputConfig)
placeholder)
placeholder

func TestForwardAsAnthropic_NormalizesRoutingAndEffortForGpt54XHigh(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	body := []byte(`{"model":"gpt-5.4-xhigh","max_tokens":16,"messages":[{"role":"user","content":"hello"placeholder],"stream":falseplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstreamBody := strings.Join([]string{
		`data: {"type":"response.completed","response":{"id":"resp_1","object":"response","model":"gpt-5.4","status":"completed","output":[{"type":"message","id":"msg_1","role":"assistant","status":"completed","content":[{"type":"output_text","text":"ok"placeholder]placeholder],"usage":{"input_tokens":5,"output_tokens":2,"total_tokens":7placeholderplaceholderplaceholder`,
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid_compat"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamBody)),
placeholderplaceholder

	svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          1,
		Name:        "openai-oauth",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token":       "oauth-token",
			"chatgpt_account_id": "chatgpt-acc",
			"model_mapping": map[string]any{
				"gpt-5.4": "gpt-5.4",
		placeholder,
	placeholder,
placeholder

	result, err := svc.ForwardAsAnthropic(context.Background(), c, account, body, "", "gpt-5.1")
placeholder
	require.NotNil(t, result)
	require.Equal(t, "gpt-5.4-xhigh", result.Model)
	require.Equal(t, "gpt-5.4", result.UpstreamModel)
	require.Equal(t, "gpt-5.4", result.BillingModel)
	require.NotNil(t, result.ReasoningEffort)
	require.Equal(t, "xhigh", *result.ReasoningEffort)

	require.Equal(t, "gpt-5.4", gjson.GetBytes(upstream.lastBody, "model").String())
	require.Equal(t, "xhigh", gjson.GetBytes(upstream.lastBody, "reasoning.effort").String())
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "gpt-5.4-xhigh", gjson.GetBytes(rec.Body.Bytes(), "model").String())
	require.Equal(t, "ok", gjson.GetBytes(rec.Body.Bytes(), "content.0.text").String())
	t.Logf("upstream body: %s", string(upstream.lastBody))
	t.Logf("response body: %s", rec.Body.String())
placeholder

func TestForwardAsAnthropic_ForcedCodexInstructionsTemplatePrependsRenderedInstructions(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	templateDir := t.TempDir()
	templatePath := filepath.Join(templateDir, "codex-instructions.md.tmpl")
	require.NoError(t, os.WriteFile(templatePath, []byte("server-prefix\n\n{{ .ExistingInstructions placeholderplaceholder"), 0o644))

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	body := []byte(`{"model":"gpt-5.4","max_tokens":16,"system":"client-system","messages":[{"role":"user","content":"hello"placeholder],"stream":falseplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstreamBody := strings.Join([]string{
		`data: {"type":"response.completed","response":{"id":"resp_1","object":"response","model":"gpt-5.4","status":"completed","output":[{"type":"message","id":"msg_1","role":"assistant","status":"completed","content":[{"type":"output_text","text":"ok"placeholder]placeholder],"usage":{"input_tokens":5,"output_tokens":2,"total_tokens":7placeholderplaceholderplaceholder`,
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid_forced"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamBody)),
placeholderplaceholder

	svc := &OpenAIGatewayService{
		cfg: &config.Config{Gateway: config.GatewayConfig{
			ForcedCodexInstructionsTemplateFile: templatePath,
			ForcedCodexInstructionsTemplate:     "server-prefix\n\n{{ .ExistingInstructions placeholderplaceholder",
placeholder
		httpUpstream: upstream,
placeholder
	account := &Account{
		ID:          1,
		Name:        "openai-oauth",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token":       "oauth-token",
			"chatgpt_account_id": "chatgpt-acc",
	placeholder,
placeholder

	result, err := svc.ForwardAsAnthropic(context.Background(), c, account, body, "", "gpt-5.1")
placeholder
	require.NotNil(t, result)
	require.Equal(t, "server-prefix\n\nclient-system", gjson.GetBytes(upstream.lastBody, "instructions").String())
placeholder

func TestForwardAsAnthropic_ForcedCodexInstructionsTemplateUsesCachedTemplateContent(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	body := []byte(`{"model":"gpt-5.4","max_tokens":16,"system":"client-system","messages":[{"role":"user","content":"hello"placeholder],"stream":falseplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstreamBody := strings.Join([]string{
		`data: {"type":"response.completed","response":{"id":"resp_1","object":"response","model":"gpt-5.4","status":"completed","output":[{"type":"message","id":"msg_1","role":"assistant","status":"completed","content":[{"type":"output_text","text":"ok"placeholder]placeholder],"usage":{"input_tokens":5,"output_tokens":2,"total_tokens":7placeholderplaceholderplaceholder`,
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid_forced_cached"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamBody)),
placeholderplaceholder

	svc := &OpenAIGatewayService{
		cfg: &config.Config{Gateway: config.GatewayConfig{
			ForcedCodexInstructionsTemplateFile: "/path/that/should/not/be/read.tmpl",
			ForcedCodexInstructionsTemplate:     "cached-prefix\n\n{{ .ExistingInstructions placeholderplaceholder",
placeholder
		httpUpstream: upstream,
placeholder
	account := &Account{
		ID:          1,
		Name:        "openai-oauth",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token":       "oauth-token",
			"chatgpt_account_id": "chatgpt-acc",
	placeholder,
placeholder

	result, err := svc.ForwardAsAnthropic(context.Background(), c, account, body, "", "gpt-5.1")
placeholder
	require.NotNil(t, result)
	require.Equal(t, "cached-prefix\n\nclient-system", gjson.GetBytes(upstream.lastBody, "instructions").String())
placeholder

func TestForwardAsAnthropic_ClientDisconnectDrainsUpstreamUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Writer = &openAICompatFailingWriter{ResponseWriter: c.Writer, failAfter: 0placeholder
	body := []byte(`{"model":"gpt-5.4","max_tokens":16,"messages":[{"role":"user","content":"hello"placeholder],"stream":trueplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstreamBody := strings.Join([]string{
		`data: {"type":"response.created","response":{"id":"resp_1","model":"gpt-5.4","status":"in_progress","output":[]placeholderplaceholder`,
		"",
		`data: {"type":"response.output_text.delta","delta":"ok"placeholder`,
		"",
		`data: {"type":"response.completed","response":{"id":"resp_1","object":"response","model":"gpt-5.4","status":"completed","output":[{"type":"message","id":"msg_1","role":"assistant","status":"completed","content":[{"type":"output_text","text":"ok"placeholder]placeholder],"usage":{"input_tokens":9,"output_tokens":4,"total_tokens":13,"input_tokens_details":{"cached_tokens":3placeholderplaceholderplaceholderplaceholder`,
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid_disconnect"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamBody)),
placeholderplaceholder

	svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          1,
		Name:        "openai-oauth",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token":       "oauth-token",
			"chatgpt_account_id": "chatgpt-acc",
	placeholder,
placeholder

	result, err := svc.ForwardAsAnthropic(context.Background(), c, account, body, "", "gpt-5.1")
placeholder
	require.NotNil(t, result)
	require.Equal(t, 9, result.Usage.InputTokens)
	require.Equal(t, 4, result.Usage.OutputTokens)
	require.Equal(t, 3, result.Usage.CacheReadInputTokens)
placeholder

func TestForwardAsAnthropic_TerminalUsageWithoutUpstreamCloseReturns(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Writer = &openAICompatFailingWriter{ResponseWriter: c.Writer, failAfter: 0placeholder
	body := []byte(`{"model":"gpt-5.4","max_tokens":16,"messages":[{"role":"user","content":"hello"placeholder],"stream":trueplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstreamBody := []byte(`data: {"type":"response.completed","response":{"id":"resp_1","object":"response","model":"gpt-5.4","status":"completed","output":[{"type":"message","id":"msg_1","role":"assistant","status":"completed","content":[{"type":"output_text","text":"ok"placeholder]placeholder],"usage":{"input_tokens":15,"output_tokens":6,"total_tokens":21,"input_tokens_details":{"cached_tokens":5placeholderplaceholderplaceholderplaceholder` + "\n\n")
	upstreamStream := newOpenAICompatBlockingReadCloser(upstreamBody)
	defer upstreamStream.Close()
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid_terminal_no_close"placeholderplaceholder,
		Body:       upstreamStream,
placeholderplaceholder

	svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          1,
		Name:        "openai-oauth",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token":       "oauth-token",
			"chatgpt_account_id": "chatgpt-acc",
	placeholder,
placeholder

	type forwardResult struct {
		result *OpenAIForwardResult
		err    error
placeholder
	resultCh := make(chan forwardResult, 1)
	go func() {
		result, err := svc.ForwardAsAnthropic(context.Background(), c, account, body, "", "gpt-5.1")
		resultCh <- forwardResult{result: result, err: errplaceholder
placeholder()

	select {
	case got := <-resultCh:
		require.NoError(t, got.err)
		require.NotNil(t, got.result)
		require.Equal(t, 15, got.result.Usage.InputTokens)
		require.Equal(t, 6, got.result.Usage.OutputTokens)
		require.Equal(t, 5, got.result.Usage.CacheReadInputTokens)
	case <-time.After(time.Second):
		require.Fail(t, "ForwardAsAnthropic should return after terminal usage event even if upstream keeps the connection open")
placeholder
placeholder

func TestForwardAsAnthropic_BufferedTerminalWithoutUpstreamCloseReturns(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	body := []byte(`{"model":"gpt-5.4","max_tokens":16,"messages":[{"role":"user","content":"hello"placeholder],"stream":falseplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstreamBody := []byte(`data: {"type":"response.completed","response":{"id":"resp_1","object":"response","model":"gpt-5.4","status":"completed","output":[{"type":"message","id":"msg_1","role":"assistant","status":"completed","content":[{"type":"output_text","text":"ok"placeholder]placeholder],"usage":{"input_tokens":15,"output_tokens":6,"total_tokens":21,"input_tokens_details":{"cached_tokens":5placeholderplaceholderplaceholderplaceholder` + "\n\n")
	upstreamStream := newOpenAICompatBlockingReadCloser(upstreamBody)
	defer upstreamStream.Close()
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid_buffered_terminal_no_close"placeholderplaceholder,
		Body:       upstreamStream,
placeholderplaceholder

	svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          1,
		Name:        "openai-oauth",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token":       "oauth-token",
			"chatgpt_account_id": "chatgpt-acc",
	placeholder,
placeholder

	type forwardResult struct {
		result *OpenAIForwardResult
		err    error
placeholder
	resultCh := make(chan forwardResult, 1)
	go func() {
		result, err := svc.ForwardAsAnthropic(context.Background(), c, account, body, "", "gpt-5.1")
		resultCh <- forwardResult{result: result, err: errplaceholder
placeholder()

	select {
	case got := <-resultCh:
		require.NoError(t, got.err)
		require.NotNil(t, got.result)
		require.Equal(t, 15, got.result.Usage.InputTokens)
		require.Equal(t, 6, got.result.Usage.OutputTokens)
		require.Equal(t, 5, got.result.Usage.CacheReadInputTokens)
		require.Contains(t, rec.Body.String(), `"stop_reason":"end_turn"`)
	case <-time.After(time.Second):
		require.Fail(t, "ForwardAsAnthropic buffered response should return after terminal usage event even if upstream keeps the connection open")
placeholder
placeholder

func TestForwardAsAnthropic_DoneSentinelWithoutTerminalReturnsError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	body := []byte(`{"model":"gpt-5.4","max_tokens":16,"messages":[{"role":"user","content":"hello"placeholder],"stream":trueplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstreamBody := "data: [DONE]\n\n"
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid_missing_terminal"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamBody)),
placeholderplaceholder

	svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          1,
		Name:        "openai-oauth",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token":       "oauth-token",
			"chatgpt_account_id": "chatgpt-acc",
	placeholder,
placeholder

	result, err := svc.ForwardAsAnthropic(context.Background(), c, account, body, "", "gpt-5.1")
placeholder
	require.Contains(t, err.Error(), "missing terminal event")
	require.NotNil(t, result)
	require.Zero(t, result.Usage.InputTokens)
	require.Zero(t, result.Usage.OutputTokens)
placeholder

func TestForwardAsAnthropic_UpstreamRequestIgnoresClientCancel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	reqCtx, cancel := context.WithCancel(context.Background())
	body := []byte(`{"model":"gpt-5.4","max_tokens":16,"messages":[{"role":"user","content":"hello"placeholder],"stream":falseplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body)).WithContext(reqCtx)
	c.Request.Header.Set("Content-Type", "application/json")
	cancel()

	upstreamBody := strings.Join([]string{
		`data: {"type":"response.completed","response":{"id":"resp_1","object":"response","model":"gpt-5.4","status":"completed","output":[{"type":"message","id":"msg_1","role":"assistant","status":"completed","content":[{"type":"output_text","text":"ok"placeholder]placeholder],"usage":{"input_tokens":5,"output_tokens":2,"total_tokens":7placeholderplaceholderplaceholder`,
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid_ctx"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamBody)),
placeholderplaceholder

	svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          1,
		Name:        "openai-oauth",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token":       "oauth-token",
			"chatgpt_account_id": "chatgpt-acc",
	placeholder,
placeholder

	result, err := svc.ForwardAsAnthropic(reqCtx, c, account, body, "", "gpt-5.1")
placeholder
	require.NotNil(t, result)
	require.NotNil(t, upstream.lastReq)
	require.NoError(t, upstream.lastReq.Context().Err())
placeholder
