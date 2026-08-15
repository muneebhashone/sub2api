package service

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// --- mock: 只记录临时不可调度写入，其余方法不应被调用 ---

type capacityShedAccountRepoStub struct {
	AccountRepository // 嵌入接口，未实现的方法会 panic（不应被调用）

	tempUnschedCalls int
placeholder

func (r *capacityShedAccountRepoStub) SetTempUnschedulable(_ context.Context, _ int64, _ time.Time, _ string) error {
	r.tempUnschedCalls++
	return nil
placeholder

// 上游容量降载是请求级信号：故障因素（客户端身份、模型容量）与账号无关，
// 同账号重试用尽后不得把账号临时摘掉——否则一个被降载的请求会顺着 failover
// 把整池账号逐个封禁，而每个账号都会以同一个错误失败。
func TestTempUnscheduleRetryableErrorSkipsRequestScopedTransient(t *testing.T) {
	t.Run("请求级瞬时故障不写账号状态", func(t *testing.T) {
		repo := &capacityShedAccountRepoStub{placeholder
		svc := &GatewayService{accountRepo: repoplaceholder

		svc.TempUnscheduleRetryableError(context.Background(), 1, &UpstreamFailoverError{
			StatusCode:             http.StatusBadGateway,
			RetryableOnSameAccount: true,
			RequestScopedTransient: true,
	placeholder)

		require.Zero(t, repo.tempUnschedCalls)
placeholder)

	// 对照组：同样的 502 在未标记请求级瞬时故障时仍按原有语义临时摘号，
	// 确认上面的断言来自新增守卫而非其他前置条件。
	t.Run("未标记时保持原有临时摘号语义", func(t *testing.T) {
		repo := &capacityShedAccountRepoStub{placeholder
		svc := &GatewayService{accountRepo: repoplaceholder

		svc.TempUnscheduleRetryableError(context.Background(), 1, &UpstreamFailoverError{
			StatusCode:             http.StatusBadGateway,
			RetryableOnSameAccount: true,
	placeholder)

		require.Equal(t, 1, repo.tempUnschedCalls)
placeholder)
placeholder

// 非池模式账号同样要先在同账号重试：换号不改变降载因素。
func TestStreamFailedEventCapacityShedRetriesOnSameAccount(t *testing.T) {
	nonPool := &Account{ID: 1, Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder

	for _, code := range []string{"server_is_overloaded", "slow_down"placeholder {
		payload := []byte(`{"type":"response.failed","response":{"error":{"code":"` + code + `"placeholderplaceholderplaceholder`)
		require.True(t, isOpenAIUpstreamCapacityShedEvent(payload), code)
		require.True(t, openAIStreamFailedEventRetryableOnSameAccount(nonPool, payload, "overloaded"), code)
placeholder

	// 非降载的 failed 事件在非池模式下仍不做同账号重试，避免放大改动面。
	other := []byte(`{"type":"response.failed","response":{"error":{"code":"server_error"placeholderplaceholderplaceholder`)
	require.False(t, isOpenAIUpstreamCapacityShedEvent(other))
	require.False(t, openAIStreamFailedEventRetryableOnSameAccount(nonPool, other, "boom"))
placeholder

// 上游降载的真实序列是「event: error → event: response.failed」。error 帧不算
// 客户端输出：若把它当首输出 flush，clientOutputStarted 被固化，随后的 failed
// 事件就进不了 pre-output failover 分支，只能把致命错误原样转发给客户端。
func TestOpenAIStreamErrorFrameDoesNotStartClientOutput(t *testing.T) {
	cases := []struct {
		data      string
		eventType string
		want      bool
placeholder{
		{`{"type":"error","error":{"code":"server_is_overloaded","message":"overloaded"placeholderplaceholder`, "error", falseplaceholder,
		{`{"type":"error","error":{"code":"slow_down","message":"slow down"placeholderplaceholder`, "error", falseplaceholder,
		{`{"type":"error","error":{"type":"rate_limit_error","code":"rate_limit_exceeded","message":"limited"placeholderplaceholder`, "error", falseplaceholder,
		// 不可重试类错误帧维持原样转发（不进 failover），保留上游错误细节。
		{`{"type":"error","error":{"type":"invalid_request_error","code":"content_policy_violation","message":"blocked"placeholderplaceholder`, "error", trueplaceholder,
		{`{"type":"response.failed","response":{"error":{"code":"server_is_overloaded"placeholderplaceholderplaceholder`, "response.failed", falseplaceholder,
		{`{"type":"response.created","response":{"id":"resp_1"placeholderplaceholder`, "response.created", falseplaceholder,
		{`{"type":"response.in_progress","response":{"id":"resp_1"placeholderplaceholder`, "response.in_progress", falseplaceholder,
		{`{"type":"response.output_item.added","item":{"type":"reasoning","summary":[]placeholderplaceholder`, "response.output_item.added", falseplaceholder,
		{`{"type":"response.output_item.added","item":{"type":"reasoning","encrypted_content":"ciphertext"placeholderplaceholder`, "response.output_item.added", trueplaceholder,
		{`{"type":"response.reasoning_summary_part.added","part":{"type":"summary_text","text":""placeholderplaceholder`, "response.reasoning_summary_part.added", falseplaceholder,
		{`{"type":"response.reasoning_summary_part.added","part":{"type":"summary_text","text":"thinking"placeholderplaceholder`, "response.reasoning_summary_part.added", trueplaceholder,
		{`{"type":"response.content_part.added","part":{"type":"output_text","text":""placeholderplaceholder`, "response.content_part.added", falseplaceholder,
		{`{"type":"response.output_text.delta","delta":"hi"placeholder`, "response.output_text.delta", trueplaceholder,
		{`[DONE]`, "", trueplaceholder,
placeholder
	for _, tc := range cases {
		require.Equal(t, tc.want, openAIStreamDataStartsClientOutput(tc.data, tc.eventType), "data=%s type=%s", tc.data, tc.eventType)
placeholder
placeholder

func TestOpenAIStreamMetadataPreambleAndMessageOnlyOverloadFailOver(t *testing.T) {
	gin.SetMode(gin.TestMode)
	largeMetadata := strings.Repeat("x", 16*1024)
	stream := strings.Join([]string{
		"event: response.created",
		`data: {"type":"response.created","response":{"id":"resp_1","metadata":{"padding":"` + largeMetadata + `"placeholderplaceholderplaceholder`,
		"",
		"event: response.output_item.added",
		`data: {"type":"response.output_item.added","item":{"type":"reasoning","summary":[]placeholderplaceholder`,
		"",
		"event: response.reasoning_summary_part.added",
		`data: {"type":"response.reasoning_summary_part.added","part":{"type":"summary_text","text":""placeholderplaceholder`,
		"",
		"event: error",
		`data: {"type":"error","error":{"type":"service_unavailable_error","message":"Our servers are currently overloaded. Please try again later."placeholderplaceholder`,
		"",
placeholder, "\n")

	tests := []struct {
		name string
		run  func(*OpenAIGatewayService, *gin.Context, *http.Response, *Account) error
placeholder{
		{
			name: "native",
			run: func(svc *OpenAIGatewayService, c *gin.Context, resp *http.Response, account *Account) error {
				_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, account, time.Now(), "model", "model")
				return err
		placeholder,
	placeholder,
		{
			name: "passthrough",
			run: func(svc *OpenAIGatewayService, c *gin.Context, resp *http.Response, account *Account) error {
				_, err := svc.handleStreamingResponsePassthrough(c.Request.Context(), resp, c, account, time.Now(), "model", "model")
				return err
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &OpenAIGatewayService{cfg: &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholderplaceholder
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(stream)),
				Header:     http.Header{"X-Request-Id": []string{"rid-message-only-overload"placeholderplaceholder,
		placeholder
			account := &Account{ID: 1, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Name: "acc"placeholder

			err := tt.run(svc, c, resp, account)
		placeholder
			var failoverErr *UpstreamFailoverError
			require.ErrorAs(t, err, &failoverErr)
			require.True(t, failoverErr.RetryableOnSameAccount)
			require.True(t, failoverErr.RequestScopedTransient)
			require.False(t, c.Writer.Written())
			require.Empty(t, rec.Body.String())
	placeholder)
placeholder
placeholder

// 回归用例（真实上游降载序列）：created → in_progress → error 帧 → response.failed。
// 期望仍然走 pre-output failover（同账号重试 + 请求级瞬时标记），且不向客户端写出任何字节。
func TestOpenAIStreamCapacityShedErrorFramePrecedingFailedStillFailsOver(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(strings.NewReader(strings.Join([]string{
			"event: response.created",
			`data: {"type":"response.created","response":{"id":"resp_1"placeholder,"sequence_number":0placeholder`,
			"",
			"event: response.in_progress",
			`data: {"type":"response.in_progress","response":{"id":"resp_1"placeholder,"sequence_number":1placeholder`,
			"",
			"event: error",
			`data: {"type":"error","error":{"type":"service_unavailable_error","code":"server_is_overloaded","message":"Our servers are currently overloaded. Please try again later."placeholder,"sequence_number":2placeholder`,
			"",
			"event: response.failed",
			`data: {"type":"response.failed","response":{"id":"resp_1","status":"failed","error":{"code":"server_is_overloaded","message":"Our servers are currently overloaded. Please try again later."placeholderplaceholder,"sequence_number":3placeholder`,
			"",
	placeholder, "\n"))),
		Header: http.Header{"X-Request-Id": []string{"rid-shed-error-then-failed"placeholderplaceholder,
placeholder

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Name: "acc"placeholder, time.Now(), "model", "model")
placeholder
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.True(t, failoverErr.RetryableOnSameAccount)
	require.True(t, failoverErr.RequestScopedTransient)
	require.False(t, c.Writer.Written())
	require.Empty(t, rec.Body.String())
placeholder

// 流中途（已有真实输出）降载时无法再 failover，此时必须把降载码改写为客户端
// 可重试的 server_error 再转发——Codex 对 server_is_overloaded/slow_down 判致命
// 并终止会话，对其余错误码执行内置退避重试。消息原样保留。
func TestOpenAIStreamCapacityShedAfterOutputRewritesCodeForClient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(strings.NewReader(strings.Join([]string{
			"event: response.created",
			`data: {"type":"response.created","response":{"id":"resp_1"placeholderplaceholder`,
			"",
			"event: response.output_text.delta",
			`data: {"type":"response.output_text.delta","delta":"partial"placeholder`,
			"",
			"event: error",
			`data: {"type":"error","error":{"type":"service_unavailable_error","code":"server_is_overloaded","message":"Our servers are currently overloaded. Please try again later."placeholder,"sequence_number":2placeholder`,
			"",
			"event: response.failed",
			`data: {"type":"response.failed","response":{"id":"resp_1","status":"failed","error":{"code":"server_is_overloaded","message":"Our servers are currently overloaded. Please try again later."placeholderplaceholder,"sequence_number":3placeholder`,
			"",
	placeholder, "\n"))),
		Header: http.Header{"X-Request-Id": []string{"rid-shed-after-output"placeholderplaceholder,
placeholder

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Name: "acc"placeholder, time.Now(), "model", "model")
placeholder
	var failoverErr *UpstreamFailoverError
	require.False(t, errors.As(err, &failoverErr))

	body := rec.Body.String()
	require.Contains(t, body, "partial")
	require.Contains(t, body, "event: response.failed")
	require.Contains(t, body, `"code":"server_error"`)
	require.NotContains(t, body, "server_is_overloaded")
	require.Contains(t, body, "Our servers are currently overloaded")
placeholder

// helper 单测：只有降载码被改写，其余错误码（尤其 rate_limit_exceeded，客户端
// 依赖其原码解析重试延时）必须原样保留。
func TestSanitizeOpenAICapacityShedErrorCodeForClient(t *testing.T) {
	cases := []struct {
		name        string
		payload     string
		wantChanged bool
		wantContain string
placeholder{
		{
			name:        "failed事件嵌套code改写",
			payload:     `{"type":"response.failed","response":{"error":{"code":"server_is_overloaded","message":"overloaded"placeholderplaceholderplaceholder`,
			wantChanged: true,
			wantContain: `"code":"server_error"`,
	placeholder,
		{
			name:        "error帧裸code改写",
			payload:     `{"type":"error","error":{"code":"slow_down","message":"slow down"placeholderplaceholder`,
			wantChanged: true,
			wantContain: `"code":"server_error"`,
	placeholder,
		{
			name:        "failed事件只有过载文案时补充code",
			payload:     `{"type":"response.failed","response":{"error":{"message":"Our servers are currently overloaded. Please try again later."placeholderplaceholderplaceholder`,
			wantChanged: true,
			wantContain: `"code":"server_error"`,
	placeholder,
		{
			name:        "error帧只有过载文案时补充code",
			payload:     `{"type":"error","error":{"message":"Server is overloaded. Please try again later."placeholderplaceholder`,
			wantChanged: true,
			wantContain: `"code":"server_error"`,
	placeholder,
		{
			name:        "rate_limit不改写",
			payload:     `{"type":"response.failed","response":{"error":{"code":"rate_limit_exceeded","message":"try again in 3s"placeholderplaceholderplaceholder`,
			wantChanged: false,
			wantContain: `"code":"rate_limit_exceeded"`,
	placeholder,
		{
			name:        "普通server_error不改写",
			payload:     `{"type":"response.failed","response":{"error":{"code":"server_error","message":"boom"placeholderplaceholderplaceholder`,
			wantChanged: false,
			wantContain: `"code":"server_error"`,
	placeholder,
		{
			name:        "非JSON不改写",
			payload:     `not-json`,
			wantChanged: false,
			wantContain: `not-json`,
	placeholder,
placeholder
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, changed := sanitizeOpenAICapacityShedErrorCodeForClient([]byte(tc.payload))
			require.Equal(t, tc.wantChanged, changed)
			require.Contains(t, string(out), tc.wantContain)
			if changed {
				require.NotContains(t, string(out), "server_is_overloaded")
				require.NotContains(t, string(out), "slow_down")
		placeholder
	placeholder)
placeholder
placeholder

// 出站身份的版本声明只能有一个来源：UA 的版本段、version 头、探针版本三处必须同源，
// 各自硬编码会漂移成互相矛盾的身份，而自相矛盾或陈旧的身份会被上游优先降载。
func TestCodexOutboundVersionHasSingleSource(t *testing.T) {
	require.True(t,
		strings.HasPrefix(codexCLIUserAgent, openai.CodexDefaultOriginator+"/"+codexCLIVersion+" "),
		"codexCLIUserAgent=%q 必须以 codexCLIVersion=%q 作为版本段", codexCLIUserAgent, codexCLIVersion,
	)
	require.Equal(t, codexCLIVersion, openAICodexProbeVersion)
	require.GreaterOrEqual(t, CompareVersions(codexCLIVersion, codexUpstreamMinVersion), 0,
		"codexCLIVersion=%q 不得低于上游最低门槛 %q", codexCLIVersion, codexUpstreamMinVersion,
	)
placeholder
