package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/antigravity"
	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// antigravityFailingWriter 模拟客户端断开连接的 gin.ResponseWriter
type antigravityFailingWriter struct {
	gin.ResponseWriter
	failAfter int // 允许成功写入的次数，之后所有写入返回错误
	writes    int
placeholder

func (w *antigravityFailingWriter) Write(p []byte) (int, error) {
	if w.writes >= w.failAfter {
		return 0, errors.New("write failed: client disconnected")
placeholder
	w.writes++
	return w.ResponseWriter.Write(p)
placeholder

// newAntigravityTestService 创建用于流式测试的 AntigravityGatewayService
func newAntigravityTestService(cfg *config.Config) *AntigravityGatewayService {
	return &AntigravityGatewayService{
		settingService: &SettingService{cfg: cfgplaceholder,
placeholder
placeholder

func TestAntigravityUpstreamErrorBodyReadLimit_RespectsDiagnosticLimit(t *testing.T) {
	svc := newAntigravityTestService(&config.Config{Gateway: config.GatewayConfig{
		LogUpstreamErrorBody:         true,
		LogUpstreamErrorBodyMaxBytes: int(gatewayUpstreamErrorBodyReadLimit) + 1024,
placeholderplaceholder)

	require.Equal(t, int64(svc.settingService.cfg.Gateway.LogUpstreamErrorBodyMaxBytes), svc.upstreamErrorBodyReadLimit())
placeholder

func TestStripSignatureSensitiveBlocksFromClaudeRequest(t *testing.T) {
	req := &antigravity.ClaudeRequest{
		Model: "claude-sonnet-4-5",
		Thinking: &antigravity.ThinkingConfig{
			Type:         "enabled",
			BudgetTokens: 1024,
	placeholder,
		Messages: []antigravity.ClaudeMessage{
			{
				Role: "assistant",
				Content: json.RawMessage(`[
					{"type":"thinking","thinking":"secret plan","signature":""placeholder,
					{"type":"tool_use","id":"t1","name":"Bash","input":{"command":"ls"placeholderplaceholder
				]`),
		placeholder,
			{
				Role: "user",
				Content: json.RawMessage(`[
					{"type":"tool_result","tool_use_id":"t1","content":"ok","is_error":falseplaceholder,
					{"type":"redacted_thinking","data":"..."placeholder
				]`),
		placeholder,
	placeholder,
placeholder

	changed, err := stripSignatureSensitiveBlocksFromClaudeRequest(req)
placeholder
	require.True(t, changed)
	require.Nil(t, req.Thinking)

	require.Len(t, req.Messages, 2)

	var blocks0 []map[string]any
	require.NoError(t, json.Unmarshal(req.Messages[0].Content, &blocks0))
	require.Len(t, blocks0, 2)
	require.Equal(t, "text", blocks0[0]["type"])
	require.Equal(t, "secret plan", blocks0[0]["text"])
	require.Equal(t, "text", blocks0[1]["type"])

	var blocks1 []map[string]any
	require.NoError(t, json.Unmarshal(req.Messages[1].Content, &blocks1))
	require.Len(t, blocks1, 1)
	require.Equal(t, "text", blocks1[0]["type"])
	require.NotEmpty(t, blocks1[0]["text"])
placeholder

func TestStripThinkingFromClaudeRequest_DoesNotDowngradeTools(t *testing.T) {
	req := &antigravity.ClaudeRequest{
		Model: "claude-sonnet-4-5",
		Thinking: &antigravity.ThinkingConfig{
			Type:         "enabled",
			BudgetTokens: 1024,
	placeholder,
		Messages: []antigravity.ClaudeMessage{
			{
				Role:    "assistant",
				Content: json.RawMessage(`[{"type":"thinking","thinking":"secret plan"placeholder,{"type":"tool_use","id":"t1","name":"Bash","input":{"command":"ls"placeholderplaceholder]`),
		placeholder,
	placeholder,
placeholder

	changed, err := stripThinkingFromClaudeRequest(req)
placeholder
	require.True(t, changed)
	require.Nil(t, req.Thinking)

	var blocks []map[string]any
	require.NoError(t, json.Unmarshal(req.Messages[0].Content, &blocks))
	require.Len(t, blocks, 2)
	require.Equal(t, "text", blocks[0]["type"])
	require.Equal(t, "secret plan", blocks[0]["text"])
	require.Equal(t, "tool_use", blocks[1]["type"])
placeholder

func TestIsPromptTooLongError(t *testing.T) {
	require.True(t, isPromptTooLongError([]byte(`{"error":{"message":"Prompt is too long"placeholderplaceholder`)))
	require.True(t, isPromptTooLongError([]byte(`{"message":"Prompt is too long"placeholder`)))
	require.False(t, isPromptTooLongError([]byte(`{"error":{"message":"other"placeholderplaceholder`)))
placeholder

type httpUpstreamStub struct {
	resp *http.Response
	err  error
placeholder

func (s *httpUpstreamStub) Do(_ *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
	return s.resp, s.err
placeholder

func (s *httpUpstreamStub) DoWithTLS(_ *http.Request, _ string, _ int64, _ int, _ *tlsfingerprint.Profile) (*http.Response, error) {
	return s.resp, s.err
placeholder

type queuedHTTPUpstreamStub struct {
	responses     []*http.Response
	errors        []error
	requestBodies [][]byte
	callCount     int
	onCall        func(*http.Request, *queuedHTTPUpstreamStub)
placeholder

func (s *queuedHTTPUpstreamStub) Do(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
	if req != nil && req.Body != nil {
		body, _ := io.ReadAll(req.Body)
		s.requestBodies = append(s.requestBodies, body)
		req.Body = io.NopCloser(bytes.NewReader(body))
placeholder else {
		s.requestBodies = append(s.requestBodies, nil)
placeholder

	idx := s.callCount
	s.callCount++
	if s.onCall != nil {
		s.onCall(req, s)
placeholder

	var resp *http.Response
	if idx < len(s.responses) {
		resp = s.responses[idx]
placeholder
	var err error
	if idx < len(s.errors) {
		err = s.errors[idx]
placeholder
	if resp == nil && err == nil {
		return nil, errors.New("unexpected upstream call")
placeholder
	return resp, err
placeholder

func (s *queuedHTTPUpstreamStub) DoWithTLS(req *http.Request, proxyURL string, accountID int64, concurrency int, _ *tlsfingerprint.Profile) (*http.Response, error) {
	return s.Do(req, proxyURL, accountID, concurrency)
placeholder

type recordingInternal500CounterCache struct {
	incrementCalls []int64
	resetCalls     []int64
placeholder

func (c *recordingInternal500CounterCache) IncrementInternal500Count(_ context.Context, accountID int64) (int64, error) {
	c.incrementCalls = append(c.incrementCalls, accountID)
	return int64(len(c.incrementCalls)), nil
placeholder

func (c *recordingInternal500CounterCache) ResetInternal500Count(_ context.Context, accountID int64) error {
	c.resetCalls = append(c.resetCalls, accountID)
	return nil
placeholder

type antigravitySettingRepoStub struct {
	values map[string]string
placeholder

func (s *antigravitySettingRepoStub) Get(ctx context.Context, key string) (*Setting, error) {
	panic("unexpected Get call")
placeholder

func (s *antigravitySettingRepoStub) GetValue(ctx context.Context, key string) (string, error) {
	if value, ok := s.values[key]; ok {
		return value, nil
placeholder
	return "", ErrSettingNotFound
placeholder

func (s *antigravitySettingRepoStub) Set(ctx context.Context, key, value string) error {
	panic("unexpected Set call")
placeholder

func (s *antigravitySettingRepoStub) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	panic("unexpected GetMultiple call")
placeholder

func (s *antigravitySettingRepoStub) SetMultiple(ctx context.Context, settings map[string]string) error {
	panic("unexpected SetMultiple call")
placeholder

func (s *antigravitySettingRepoStub) GetAll(ctx context.Context) (map[string]string, error) {
	panic("unexpected GetAll call")
placeholder

func (s *antigravitySettingRepoStub) Delete(ctx context.Context, key string) error {
	panic("unexpected Delete call")
placeholder

func TestResolveAntigravityProjectID(t *testing.T) {
	tests := []struct {
		name    string
		account *Account
		want    string
		wantErr bool
placeholder{
		{
			name: "uses onboard project_id first",
			account: &Account{Credentials: map[string]any{
				"project_id": " onboard-project ",
				antigravityProjectIDFallbackCredentialKey: " configured-project ",
	placeholder
			want: "onboard-project",
	placeholder,
		{
			name: "uses configured credentials fallback",
			account: &Account{Credentials: map[string]any{
				antigravityProjectIDFallbackCredentialKey: " configured-project ",
	placeholder
			want: "configured-project",
	placeholder,
		{
			name: "uses configured extra fallback",
			account: &Account{Extra: map[string]any{
				antigravityProjectIDFallbackCredentialKey: " extra-project ",
	placeholder
			want: "extra-project",
	placeholder,
		{
			name:    "missing project",
			account: &Account{Credentials: map[string]any{placeholderplaceholder,
			wantErr: true,
	placeholder,
placeholder

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := resolveAntigravityProjectID(tc.account)
			if tc.wantErr {
				require.ErrorIs(t, err, errAntigravityProjectIDRequired)
				require.Empty(t, got)
				return
		placeholder
		placeholder
			require.Equal(t, tc.want, got)
	placeholder)
placeholder
placeholder

func TestAntigravityGatewayService_ForwardGemini_UsesConfiguredProjectFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)

	body, err := json.Marshal(map[string]any{
		"contents": []map[string]any{
			{"role": "user", "parts": []map[string]any{{"text": "hello"placeholderplaceholderplaceholder,
	placeholder,
placeholder)
placeholder
	c.Request = httptest.NewRequest(http.MethodPost, "/antigravity/v1beta/models/gemini-2.5-flash:streamGenerateContent", bytes.NewReader(body))

	upstreamBody := []byte("data: {\"response\":{\"candidates\":[{\"content\":{\"parts\":[{\"text\":\"ok\"placeholder]placeholder,\"finishReason\":\"STOP\"placeholder],\"usageMetadata\":{\"promptTokenCount\":1,\"candidatesTokenCount\":1placeholderplaceholderplaceholder\n\n")
	upstream := &queuedHTTPUpstreamStub{
		responses: []*http.Response{
			{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
				Body:       io.NopCloser(bytes.NewReader(upstreamBody)),
		placeholder,
	placeholder,
placeholder
	svc := &AntigravityGatewayService{
		settingService: NewSettingService(&antigravitySettingRepoStub{placeholder, &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder),
		tokenProvider:  &AntigravityTokenProvider{placeholder,
		httpUpstream:   upstream,
placeholder

	account := &Account{
		ID:          101,
		Name:        "acc-configured-project",
		Platform:    PlatformAntigravity,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"access_token": "token",
			antigravityProjectIDFallbackCredentialKey: "configured-project",
			"model_mapping": map[string]any{
				"gemini-2.5-flash": "gemini-2.5-flash",
		placeholder,
	placeholder,
placeholder

	result, err := svc.ForwardGemini(context.Background(), c, account, "gemini-2.5-flash", "streamGenerateContent", true, body, false)
placeholder
	require.NotNil(t, result)
	require.Len(t, upstream.requestBodies, 1)

	var wrapped map[string]any
	require.NoError(t, json.Unmarshal(upstream.requestBodies[0], &wrapped))
	require.Equal(t, "configured-project", wrapped["project"])
placeholder

func TestAntigravityGatewayService_ForwardGemini_PreservesServerSideToolInvocationConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"contents":[{"role":"user","parts":[{"text":"hello"placeholder]placeholder],"tools":[{"functionDeclarations":[{"name":"get_weather","parameters":{"type":"object","additionalProperties":falseplaceholderplaceholder]placeholder,{"googleSearch":{placeholderplaceholder],"toolConfig":{"includeServerSideToolInvocations":trueplaceholderplaceholder`)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)
	body = bytes.ReplaceAll(body, []byte{92placeholder, nil)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1beta/models/gemini-2.5-flash:generateContent", bytes.NewReader(body))

	upstream := &queuedHTTPUpstreamStub{responses: []*http.Response{{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader("data: {\"response\":{\"candidates\":[{\"content\":{\"parts\":[{\"text\":\"ok\"placeholder]placeholder,\"finishReason\":\"STOP\"placeholder],\"usageMetadata\":{placeholderplaceholderplaceholder\n\n")),
placeholderplaceholderplaceholder
	svc := &AntigravityGatewayService{
		settingService: NewSettingService(&antigravitySettingRepoStub{placeholder, &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder),
		tokenProvider:  &AntigravityTokenProvider{placeholder,
		httpUpstream:   upstream,
placeholder
	account := &Account{
		ID: 103, Name: "native-gemini", Platform: PlatformAntigravity, Type: AccountTypeOAuth, Status: StatusActive, Concurrency: 1,
placeholder"access_token": "token", "project_id": "project-103", "model_mapping": map[string]any{"gemini-2.5-flash": "gemini-2.5-flash"placeholderplaceholder,
placeholder

	result, err := svc.ForwardGemini(context.Background(), c, account, "gemini-2.5-flash", "generateContent", false, body, false)
placeholder
	require.NotNil(t, result)
	require.Len(t, upstream.requestBodies, 1)

	var wrapped map[string]any
	require.NoError(t, json.Unmarshal(upstream.requestBodies[0], &wrapped))
	request := wrapped["request"].(map[string]any)
	toolConfig := request["toolConfig"].(map[string]any)
	require.Equal(t, true, toolConfig["includeServerSideToolInvocations"])
	require.NotContains(t, toolConfig, "include_server_side_tool_invocations")
placeholder

func TestAntigravityGatewayService_ForwardGemini_MissingProjectReturnsLocalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)

	body, err := json.Marshal(map[string]any{
		"contents": []map[string]any{
			{"role": "user", "parts": []map[string]any{{"text": "hello"placeholderplaceholderplaceholder,
	placeholder,
placeholder)
placeholder
	c.Request = httptest.NewRequest(http.MethodPost, "/antigravity/v1beta/models/gemini-2.5-flash:streamGenerateContent", bytes.NewReader(body))

	upstream := &queuedHTTPUpstreamStub{placeholder
	internal500Cache := &recordingInternal500CounterCache{placeholder
	svc := &AntigravityGatewayService{
		tokenProvider:    &AntigravityTokenProvider{placeholder,
		httpUpstream:     upstream,
		internal500Cache: internal500Cache,
placeholder

	account := &Account{
		ID:          102,
		Name:        "acc-missing-project",
		Platform:    PlatformAntigravity,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"access_token": "token",
			"model_mapping": map[string]any{
				"gemini-2.5-flash": "gemini-2.5-flash",
		placeholder,
	placeholder,
placeholder

	result, err := svc.ForwardGemini(context.Background(), c, account, "gemini-2.5-flash", "streamGenerateContent", true, body, false)
	require.Nil(t, result)
	require.ErrorIs(t, err, errAntigravityProjectIDRequired)
	require.Equal(t, http.StatusBadRequest, writer.Code)
	require.Empty(t, upstream.requestBodies)
	require.Empty(t, internal500Cache.incrementCalls)
	require.Contains(t, writer.Body.String(), "project_id")
	require.NotContains(t, writer.Body.String(), `"project":""`)
placeholder

func TestAntigravityGatewayService_Forward_PromptTooLong(t *testing.T) {
	gin.SetMode(gin.TestMode)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)

	body, err := json.Marshal(map[string]any{
		"model": "claude-opus-4-6",
		"messages": []map[string]any{
			{"role": "user", "content": "hi"placeholder,
	placeholder,
		"max_tokens": 1,
		"stream":     false,
placeholder)
placeholder

	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	c.Request = req

	respBody := []byte(`{"error":{"message":"Prompt is too long"placeholderplaceholder`)
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     http.Header{"X-Request-Id": []string{"req-1"placeholderplaceholder,
		Body:       io.NopCloser(bytes.NewReader(respBody)),
placeholder

	svc := &AntigravityGatewayService{
		settingService: NewSettingService(&antigravitySettingRepoStub{placeholder, &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder),
		tokenProvider:  &AntigravityTokenProvider{placeholder,
		httpUpstream:   &httpUpstreamStub{resp: respplaceholder,
placeholder

	account := &Account{
		ID:          1,
		Name:        "acc-1",
		Platform:    PlatformAntigravity,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"access_token": "token",
	placeholder
	placeholder,
placeholder

	result, err := svc.Forward(context.Background(), c, account, body, false)
	require.Nil(t, result)

	var promptErr *PromptTooLongError
	require.ErrorAs(t, err, &promptErr)
	require.Equal(t, http.StatusBadRequest, promptErr.StatusCode)
	require.Equal(t, "req-1", promptErr.RequestID)
	require.NotEmpty(t, promptErr.Body)

	raw, ok := c.Get(OpsUpstreamErrorsKey)
	require.True(t, ok)
	events, ok := raw.([]*OpsUpstreamErrorEvent)
	require.True(t, ok)
	require.Len(t, events, 1)
	require.Equal(t, "prompt_too_long", events[0].Kind)
placeholder

// TestAntigravityGatewayService_Forward_ModelRateLimitTriggersFailover
// 验证：当账号存在模型限流且剩余时间 >= antigravityRateLimitThreshold 时，
// Forward 方法应返回 UpstreamFailoverError，触发 Handler 切换账号
func TestAntigravityGatewayService_Forward_ModelRateLimitTriggersFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)

	body, err := json.Marshal(map[string]any{
		"model": "claude-opus-4-6",
		"messages": []map[string]any{
			{"role": "user", "content": "hi"placeholder,
	placeholder,
		"max_tokens": 1,
		"stream":     false,
placeholder)
placeholder

	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	c.Request = req

	// 不需要真正调用上游，因为预检查会直接返回切换信号
	svc := &AntigravityGatewayService{
		tokenProvider: &AntigravityTokenProvider{placeholder,
		httpUpstream:  &httpUpstreamStub{resp: nil, err: nilplaceholder,
placeholder

	// 设置模型限流：剩余时间 30 秒（> antigravityRateLimitThreshold 7s）
	futureResetAt := time.Now().Add(30 * time.Second).Format(time.RFC3339)
	account := &Account{
		ID:          1,
		Name:        "acc-rate-limited",
		Platform:    PlatformAntigravity,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"access_token": "token",
	placeholder
	placeholder,
		Extra: map[string]any{
			modelRateLimitsKey: map[string]any{
				"claude-opus-4-6-thinking": map[string]any{
					"rate_limit_reset_at": futureResetAt,
			placeholder,
		placeholder,
	placeholder,
placeholder

	result, err := svc.Forward(context.Background(), c, account, body, false)
	require.Nil(t, result, "Forward should not return result when model rate limited")
	require.NotNil(t, err, "Forward should return error")

	// 核心验证：错误应该是 UpstreamFailoverError，而不是普通 502 错误
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr, "error should be UpstreamFailoverError to trigger account switch")
	require.Equal(t, http.StatusServiceUnavailable, failoverErr.StatusCode)
	// 非粘性会话请求，ForceCacheBilling 应为 false
	require.False(t, failoverErr.ForceCacheBilling, "ForceCacheBilling should be false for non-sticky session")
placeholder

// TestAntigravityGatewayService_ForwardGemini_ModelRateLimitTriggersFailover
// 验证：ForwardGemini 方法同样能正确将 AntigravityAccountSwitchError 转换为 UpstreamFailoverError
func TestAntigravityGatewayService_ForwardGemini_ModelRateLimitTriggersFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)

	body, err := json.Marshal(map[string]any{
		"contents": []map[string]any{
			{"role": "user", "parts": []map[string]any{{"text": "hi"placeholderplaceholderplaceholder,
	placeholder,
placeholder)
placeholder

	req := httptest.NewRequest(http.MethodPost, "/v1beta/models/gemini-2.5-flash:generateContent", bytes.NewReader(body))
	c.Request = req

	// 不需要真正调用上游，因为预检查会直接返回切换信号
	svc := &AntigravityGatewayService{
		tokenProvider: &AntigravityTokenProvider{placeholder,
		httpUpstream:  &httpUpstreamStub{resp: nil, err: nilplaceholder,
placeholder

	// 设置模型限流：剩余时间 30 秒（> antigravityRateLimitThreshold 7s）
	futureResetAt := time.Now().Add(30 * time.Second).Format(time.RFC3339)
	account := &Account{
		ID:          2,
		Name:        "acc-gemini-rate-limited",
		Platform:    PlatformAntigravity,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"access_token": "token",
	placeholder
	placeholder,
		Extra: map[string]any{
			modelRateLimitsKey: map[string]any{
				"gemini-2.5-flash": map[string]any{
					"rate_limit_reset_at": futureResetAt,
			placeholder,
		placeholder,
	placeholder,
placeholder

	result, err := svc.ForwardGemini(context.Background(), c, account, "gemini-2.5-flash", "generateContent", false, body, false)
	require.Nil(t, result, "ForwardGemini should not return result when model rate limited")
	require.NotNil(t, err, "ForwardGemini should return error")

	// 核心验证：错误应该是 UpstreamFailoverError，而不是普通 502 错误
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr, "error should be UpstreamFailoverError to trigger account switch")
	require.Equal(t, http.StatusServiceUnavailable, failoverErr.StatusCode)
	// 非粘性会话请求，ForceCacheBilling 应为 false
	require.False(t, failoverErr.ForceCacheBilling, "ForceCacheBilling should be false for non-sticky session")
placeholder

// TestAntigravityGatewayService_Forward_StickySessionForceCacheBilling
// 验证：粘性会话切换时，UpstreamFailoverError.ForceCacheBilling 应为 true
func TestAntigravityGatewayService_Forward_StickySessionForceCacheBilling(t *testing.T) {
	gin.SetMode(gin.TestMode)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)

	body, err := json.Marshal(map[string]any{
		"model":    "claude-opus-4-6",
		"messages": []map[string]string{{"role": "user", "content": "hello"placeholderplaceholder,
placeholder)
placeholder

	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	c.Request = req

	svc := &AntigravityGatewayService{
		tokenProvider: &AntigravityTokenProvider{placeholder,
		httpUpstream:  &httpUpstreamStub{resp: nil, err: nilplaceholder,
placeholder

	// 设置模型限流：剩余时间 30 秒（> antigravityRateLimitThreshold 7s）
	futureResetAt := time.Now().Add(30 * time.Second).Format(time.RFC3339)
	account := &Account{
		ID:          3,
		Name:        "acc-sticky-rate-limited",
		Platform:    PlatformAntigravity,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"access_token": "token",
	placeholder
	placeholder,
		Extra: map[string]any{
			modelRateLimitsKey: map[string]any{
				"claude-opus-4-6-thinking": map[string]any{
					"rate_limit_reset_at": futureResetAt,
			placeholder,
		placeholder,
	placeholder,
placeholder

	// 传入 isStickySession = true
	result, err := svc.Forward(context.Background(), c, account, body, true)
	require.Nil(t, result, "Forward should not return result when model rate limited")
	require.NotNil(t, err, "Forward should return error")

	// 核心验证：粘性会话切换时，ForceCacheBilling 应为 true
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr, "error should be UpstreamFailoverError to trigger account switch")
	require.Equal(t, http.StatusServiceUnavailable, failoverErr.StatusCode)
	require.True(t, failoverErr.ForceCacheBilling, "ForceCacheBilling should be true for sticky session switch")
placeholder

// TestAntigravityGatewayService_ForwardGemini_StickySessionForceCacheBilling verifies
// that ForwardGemini sets ForceCacheBilling=true for sticky session switch.
func TestAntigravityGatewayService_ForwardGemini_StickySessionForceCacheBilling(t *testing.T) {
	gin.SetMode(gin.TestMode)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)

	body, err := json.Marshal(map[string]any{
		"contents": []map[string]any{
			{"role": "user", "parts": []map[string]any{{"text": "hi"placeholderplaceholderplaceholder,
	placeholder,
placeholder)
placeholder

	req := httptest.NewRequest(http.MethodPost, "/v1beta/models/gemini-2.5-flash:generateContent", bytes.NewReader(body))
	c.Request = req

	svc := &AntigravityGatewayService{
		tokenProvider: &AntigravityTokenProvider{placeholder,
		httpUpstream:  &httpUpstreamStub{resp: nil, err: nilplaceholder,
placeholder

	// 设置模型限流：剩余时间 30 秒（> antigravityRateLimitThreshold 7s）
	futureResetAt := time.Now().Add(30 * time.Second).Format(time.RFC3339)
	account := &Account{
		ID:          4,
		Name:        "acc-gemini-sticky-rate-limited",
		Platform:    PlatformAntigravity,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"access_token": "token",
	placeholder
	placeholder,
		Extra: map[string]any{
			modelRateLimitsKey: map[string]any{
				"gemini-2.5-flash": map[string]any{
					"rate_limit_reset_at": futureResetAt,
			placeholder,
		placeholder,
	placeholder,
placeholder

	// 传入 isStickySession = true
	result, err := svc.ForwardGemini(context.Background(), c, account, "gemini-2.5-flash", "generateContent", false, body, true)
	require.Nil(t, result, "ForwardGemini should not return result when model rate limited")
	require.NotNil(t, err, "ForwardGemini should return error")

	// 核心验证：粘性会话切换时，ForceCacheBilling 应为 true
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr, "error should be UpstreamFailoverError to trigger account switch")
	require.Equal(t, http.StatusServiceUnavailable, failoverErr.StatusCode)
	require.True(t, failoverErr.ForceCacheBilling, "ForceCacheBilling should be true for sticky session switch")
placeholder

func TestAntigravityGatewayService_ForwardGemini_ClearsStickySessionOnGeminiRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)

	body, err := json.Marshal(map[string]any{
		"contents": []map[string]any{
			{"role": "user", "parts": []map[string]any{{"text": "hi"placeholderplaceholderplaceholder,
	placeholder,
placeholder)
placeholder

	req := httptest.NewRequest(http.MethodPost, "/v1beta/models/gemini-3-flash-preview:generateContent", bytes.NewReader(body))
	c.Request = req

	respBody := []byte(`{
		"error": {
			"status": "RESOURCE_EXHAUSTED",
			"details": [
				{"@type": "type.googleapis.com/google.rpc.ErrorInfo", "metadata": {"model": "gemini-3-flash"placeholder, "reason": "RATE_LIMIT_EXCEEDED"placeholder,
				{"@type": "type.googleapis.com/google.rpc.RetryInfo", "retryDelay": "15s"placeholder
			]
	placeholder
placeholder`)
	upstream := &httpUpstreamStub{resp: &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{placeholder,
		Body:       io.NopCloser(bytes.NewReader(respBody)),
placeholderplaceholder
	repo := &stubAntigravityAccountRepo{placeholder
	cache := &stubSmartRetryCache{placeholder
	svc := &AntigravityGatewayService{
		tokenProvider: &AntigravityTokenProvider{placeholder,
		httpUpstream:  upstream,
		accountRepo:   repo,
		cache:         cache,
placeholder

	account := &Account{
		ID:          44,
		Name:        "acc-gemini-runtime-rate-limited",
		Platform:    PlatformAntigravity,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
placeholder
			"access_token": "token",
			"expires_at":   time.Now().Add(time.Hour).Format(time.RFC3339),
	placeholder
	placeholder,
		Extra: map[string]any{
			"mixed_scheduling": true,
	placeholder,
placeholder

	result, err := svc.ForwardGemini(
		context.Background(),
		c,
		account,
		"gemini-3-flash-preview",
		"generateContent",
		false,
		body,
		true,
		WithForwardGeminiSession(77, "gemini:sticky-runtime"),
	)

	require.Nil(t, result)
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusServiceUnavailable, failoverErr.StatusCode)
	require.Len(t, repo.modelRateLimitCalls, 2)
	require.Equal(t, "gemini-3-flash", repo.modelRateLimitCalls[0].modelKey)
	require.Equal(t, antigravityGeminiModelRateLimitKey, repo.modelRateLimitCalls[1].modelKey)
	require.Len(t, cache.deleteCalls, 1)
	require.Equal(t, int64(77), cache.deleteCalls[0].groupID)
	require.Equal(t, "gemini:sticky-runtime", cache.deleteCalls[0].sessionHash)
placeholder

// TestAntigravityGatewayService_Forward_BillsWithMappedModel
// 验证：Antigravity Claude 转发返回的计费模型使用映射后的模型
func TestAntigravityGatewayService_Forward_BillsWithMappedModel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)

	body, err := json.Marshal(map[string]any{
		"model": "claude-sonnet-4-5",
		"messages": []map[string]any{
			{"role": "user", "content": "hello"placeholder,
	placeholder,
		"max_tokens": 16,
		"stream":     true,
placeholder)
placeholder

	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	c.Request = req

	upstreamBody := []byte("data: {\"response\":{\"candidates\":[{\"content\":{\"parts\":[{\"text\":\"ok\"placeholder]placeholder,\"finishReason\":\"STOP\"placeholder],\"usageMetadata\":{\"promptTokenCount\":8,\"candidatesTokenCount\":3placeholderplaceholderplaceholder\n\n")
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"X-Request-Id": []string{"req-bill-1"placeholderplaceholder,
		Body:       io.NopCloser(bytes.NewReader(upstreamBody)),
placeholder

	svc := &AntigravityGatewayService{
		settingService: NewSettingService(&antigravitySettingRepoStub{placeholder, &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder),
		tokenProvider:  &AntigravityTokenProvider{placeholder,
		httpUpstream:   &httpUpstreamStub{resp: respplaceholder,
placeholder

	const mappedModel = "gemini-3-pro-high"
	account := &Account{
		ID:          5,
		Name:        "acc-forward-billing",
		Platform:    PlatformAntigravity,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"access_token": "token",
	placeholder
			"model_mapping": map[string]any{
				"claude-sonnet-4-5": mappedModel,
		placeholder,
	placeholder,
placeholder

	result, err := svc.Forward(context.Background(), c, account, body, false)
placeholder
	require.NotNil(t, result)
	require.Equal(t, "claude-sonnet-4-5", result.Model)
	require.Equal(t, mappedModel, result.UpstreamModel)
placeholder

// TestAntigravityGatewayService_ForwardGemini_BillsWithMappedModel
// 验证：Antigravity Gemini 转发返回的计费模型使用映射后的模型
func TestAntigravityGatewayService_ForwardGemini_BillsWithMappedModel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)

	body, err := json.Marshal(map[string]any{
		"contents": []map[string]any{
			{"role": "user", "parts": []map[string]any{{"text": "hello"placeholderplaceholderplaceholder,
	placeholder,
placeholder)
placeholder

	req := httptest.NewRequest(http.MethodPost, "/v1beta/models/gemini-2.5-flash:generateContent", bytes.NewReader(body))
	c.Request = req

	upstreamBody := []byte("data: {\"response\":{\"candidates\":[{\"content\":{\"parts\":[{\"text\":\"ok\"placeholder]placeholder,\"finishReason\":\"STOP\"placeholder],\"usageMetadata\":{\"promptTokenCount\":8,\"candidatesTokenCount\":3placeholderplaceholderplaceholder\n\n")
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"X-Request-Id": []string{"req-bill-2"placeholderplaceholder,
		Body:       io.NopCloser(bytes.NewReader(upstreamBody)),
placeholder

	svc := &AntigravityGatewayService{
		settingService: NewSettingService(&antigravitySettingRepoStub{placeholder, &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder),
		tokenProvider:  &AntigravityTokenProvider{placeholder,
		httpUpstream:   &httpUpstreamStub{resp: respplaceholder,
placeholder

	const mappedModel = "gemini-3-pro-high"
	account := &Account{
		ID:          6,
		Name:        "acc-gemini-billing",
		Platform:    PlatformAntigravity,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"access_token": "token",
	placeholder
			"model_mapping": map[string]any{
				"gemini-2.5-flash": mappedModel,
		placeholder,
	placeholder,
placeholder

	result, err := svc.ForwardGemini(context.Background(), c, account, "gemini-2.5-flash", "generateContent", true, body, false)
placeholder
	require.NotNil(t, result)
	require.Equal(t, "gemini-2.5-flash", result.Model)
	require.Equal(t, mappedModel, result.UpstreamModel)
placeholder

func TestAntigravityGatewayService_ForwardGemini_FallbackReportsActualUpstreamModel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)

	body := []byte(`{"contents":[{"role":"user","parts":[{"text":"hello"placeholder]placeholder]placeholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1beta/models/gemini-primary:generateContent", bytes.NewReader(body))

	const (
		originalModel = "gemini-primary"
		mappedModel   = "gemini-primary-upstream"
		fallbackModel = "gemini-fallback-upstream"
	)
	upstream := &queuedHTTPUpstreamStub{responses: []*http.Response{
		{
			StatusCode: http.StatusNotFound,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"error":{"code":404,"message":"model not found"placeholderplaceholder`)),
	placeholder,
		{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
			Body: io.NopCloser(strings.NewReader(
				`data: {"response":{"modelVersion":"gemini-fallback-upstream","candidates":[{"content":{"parts":[{"text":"ok"placeholder]placeholder,"finishReason":"STOP"placeholder],"usageMetadata":{"promptTokenCount":8,"candidatesTokenCount":3placeholderplaceholderplaceholder` + "\n\n",
			)),
	placeholder,
placeholderplaceholder
	settings := &antigravitySettingRepoStub{values: map[string]string{
		SettingKeyEnableModelFallback:      "true",
		SettingKeyFallbackModelAntigravity: fallbackModel,
placeholderplaceholder
	svc := &AntigravityGatewayService{
		settingService: NewSettingService(settings, &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder),
		tokenProvider:  &AntigravityTokenProvider{placeholder,
		httpUpstream:   upstream,
placeholder
	account := &Account{
		ID:          9,
		Name:        "acc-gemini-fallback",
		Platform:    PlatformAntigravity,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"access_token": "token",
	placeholder
			"model_mapping": map[string]any{
				originalModel: mappedModel,
		placeholder,
	placeholder,
placeholder

	result, err := svc.ForwardGemini(context.Background(), c, account, originalModel, "generateContent", true, body, false)
placeholder
	require.NotNil(t, result)
	require.Equal(t, originalModel, result.Model)
	require.Equal(t, fallbackModel, result.UpstreamModel)
	require.Equal(t, fallbackModel, result.UpstreamResponseModel)
	require.False(t, result.UpstreamResponseModelConflict)
	mismatch := upstreamModelMismatch(result.UpstreamModel, result.UpstreamResponseModel)
	require.NotNil(t, mismatch)
	require.False(t, *mismatch)
	require.Len(t, upstream.requestBodies, 2)
	require.Contains(t, string(upstream.requestBodies[0]), `"model":"`+mappedModel+`"`)
	require.Contains(t, string(upstream.requestBodies[1]), `"model":"`+fallbackModel+`"`)
placeholder

func TestAntigravityGatewayService_ForwardGemini_RetriesCorruptedThoughtSignature(t *testing.T) {
	gin.SetMode(gin.TestMode)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)

	body, err := json.Marshal(map[string]any{
		"contents": []map[string]any{
			{"role": "user", "parts": []map[string]any{{"text": "hello"placeholderplaceholderplaceholder,
			{"role": "model", "parts": []map[string]any{{"text": "thinking", "thought": true, "thoughtSignature": "sig_bad_1"placeholderplaceholderplaceholder,
			{"role": "model", "parts": []map[string]any{{"functionCall": map[string]any{"name": "toolA", "args": map[string]any{"x": 1placeholderplaceholder, "thoughtSignature": "sig_bad_2"placeholderplaceholderplaceholder,
	placeholder,
placeholder)
placeholder

	req := httptest.NewRequest(http.MethodPost, "/antigravity/v1beta/models/gemini-3.1-pro-preview:streamGenerateContent", bytes.NewReader(body))
	c.Request = req

	firstRespBody := []byte(`{"response":{"error":{"code":400,"message":"Corrupted thought signature.","status":"INVALID_ARGUMENT"placeholderplaceholderplaceholder`)
	secondRespBody := []byte("data: {\"response\":{\"candidates\":[{\"content\":{\"parts\":[{\"text\":\"ok\"placeholder]placeholder,\"finishReason\":\"STOP\"placeholder],\"usageMetadata\":{\"promptTokenCount\":8,\"candidatesTokenCount\":3placeholderplaceholderplaceholder\n\n")

	upstream := &queuedHTTPUpstreamStub{
		responses: []*http.Response{
			{
				StatusCode: http.StatusBadRequest,
				Header: http.Header{
					"Content-Type": []string{"application/json"placeholder,
					"X-Request-Id": []string{"req-sig-1"placeholder,
			placeholder,
				Body: io.NopCloser(bytes.NewReader(firstRespBody)),
		placeholder,
			{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Type": []string{"text/event-stream"placeholder,
					"X-Request-Id": []string{"req-sig-2"placeholder,
			placeholder,
				Body: io.NopCloser(bytes.NewReader(secondRespBody)),
		placeholder,
	placeholder,
placeholder

	svc := &AntigravityGatewayService{
		settingService: NewSettingService(&antigravitySettingRepoStub{placeholder, &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder),
		tokenProvider:  &AntigravityTokenProvider{placeholder,
		httpUpstream:   upstream,
placeholder

	const originalModel = "gemini-3.1-pro-preview"
	const mappedModel = "gemini-3.1-pro-high"
	account := &Account{
		ID:          7,
		Name:        "acc-gemini-signature",
		Platform:    PlatformAntigravity,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"access_token": "token",
	placeholder
			"model_mapping": map[string]any{
				originalModel: mappedModel,
		placeholder,
	placeholder,
placeholder

	result, err := svc.ForwardGemini(context.Background(), c, account, originalModel, "streamGenerateContent", true, body, false)
placeholder
	require.NotNil(t, result)
	require.Equal(t, originalModel, result.Model)
	require.Equal(t, mappedModel, result.UpstreamModel)
	require.Len(t, upstream.requestBodies, 2, "signature error should trigger exactly one retry")

	firstReq := string(upstream.requestBodies[0])
	secondReq := string(upstream.requestBodies[1])
	require.Contains(t, firstReq, `"thoughtSignature":"sig_bad_1"`)
	require.Contains(t, firstReq, `"thoughtSignature":"sig_bad_2"`)
	require.Contains(t, secondReq, `"thoughtSignature":"skip_thought_signature_validator"`)
	require.NotContains(t, secondReq, `"thoughtSignature":"sig_bad_1"`)
	require.NotContains(t, secondReq, `"thoughtSignature":"sig_bad_2"`)

	raw, ok := c.Get(OpsUpstreamErrorsKey)
	require.True(t, ok)
	events, ok := raw.([]*OpsUpstreamErrorEvent)
	require.True(t, ok)
	require.NotEmpty(t, events)
	require.Equal(t, "signature_error", events[0].Kind)
placeholder

func TestAntigravityGatewayService_ForwardGemini_SignatureRetryPropagatesFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)

	body, err := json.Marshal(map[string]any{
		"contents": []map[string]any{
			{"role": "user", "parts": []map[string]any{{"text": "hello"placeholderplaceholderplaceholder,
			{"role": "model", "parts": []map[string]any{{"text": "thinking", "thought": true, "thoughtSignature": "sig_bad_1"placeholderplaceholderplaceholder,
	placeholder,
placeholder)
placeholder

	req := httptest.NewRequest(http.MethodPost, "/antigravity/v1beta/models/gemini-3.1-pro-preview:streamGenerateContent", bytes.NewReader(body))
	c.Request = req

	firstRespBody := []byte(`{"response":{"error":{"code":400,"message":"Corrupted thought signature.","status":"INVALID_ARGUMENT"placeholderplaceholderplaceholder`)

	const originalModel = "gemini-3.1-pro-preview"
	const mappedModel = "gemini-3.1-pro-high"
	account := &Account{
		ID:          8,
		Name:        "acc-gemini-signature-failover",
		Platform:    PlatformAntigravity,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Concurrency: 1,
placeholder
			"access_token": "token",
	placeholder
			"model_mapping": map[string]any{
				originalModel: mappedModel,
		placeholder,
	placeholder,
placeholder

	upstream := &queuedHTTPUpstreamStub{
		responses: []*http.Response{
			{
				StatusCode: http.StatusBadRequest,
				Header: http.Header{
					"Content-Type": []string{"application/json"placeholder,
					"X-Request-Id": []string{"req-sig-failover-1"placeholder,
			placeholder,
				Body: io.NopCloser(bytes.NewReader(firstRespBody)),
		placeholder,
	placeholder,
		onCall: func(_ *http.Request, stub *queuedHTTPUpstreamStub) {
			if stub.callCount != 1 {
				return
		placeholder
			futureResetAt := time.Now().Add(30 * time.Second).Format(time.RFC3339)
			account.Extra = map[string]any{
				modelRateLimitsKey: map[string]any{
					mappedModel: map[string]any{
						"rate_limit_reset_at": futureResetAt,
				placeholder,
			placeholder,
		placeholder
	placeholder,
placeholder

	svc := &AntigravityGatewayService{
		settingService: NewSettingService(&antigravitySettingRepoStub{placeholder, &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder),
		tokenProvider:  &AntigravityTokenProvider{placeholder,
		httpUpstream:   upstream,
placeholder

	result, err := svc.ForwardGemini(context.Background(), c, account, originalModel, "streamGenerateContent", true, body, true)
	require.Nil(t, result)

	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr, "signature retry should propagate failover instead of falling back to the original 400")
	require.Equal(t, http.StatusServiceUnavailable, failoverErr.StatusCode)
	require.True(t, failoverErr.ForceCacheBilling)
	require.Len(t, upstream.requestBodies, 1, "retry should stop at preflight failover and not issue a second upstream request")

	raw, ok := c.Get(OpsUpstreamErrorsKey)
	require.True(t, ok)
	events, ok := raw.([]*OpsUpstreamErrorEvent)
	require.True(t, ok)
	require.Len(t, events, 2)
	require.Equal(t, "signature_error", events[0].Kind)
	require.Equal(t, "failover", events[1].Kind)
placeholder

// TestStreamUpstreamResponse_UsageAndFirstToken
// 验证：usage 字段可被累积/覆盖更新，并且能记录首 token 时间
func TestStreamUpstreamResponse_UsageAndFirstToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Header: http.Header{placeholder, Body: prplaceholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		fmt.Fprintln(pw, `data: {"usage":{"input_tokens":1,"output_tokens":2,"cache_read_input_tokens":3,"cache_creation_input_tokens":4placeholderplaceholder`)
		fmt.Fprintln(pw, `data: {"usage":{"output_tokens":5placeholderplaceholder`)
placeholder()

	start := time.Now().Add(-10 * time.Millisecond)
	result := svc.streamUpstreamResponse(c, resp, start)
	_ = pr.Close()

	require.NotNil(t, result)
	require.NotNil(t, result.usage)
	require.Equal(t, 1, result.usage.InputTokens)
	// 第二次事件覆盖 output_tokens
	require.Equal(t, 5, result.usage.OutputTokens)
	require.Equal(t, 3, result.usage.CacheReadInputTokens)
	require.Equal(t, 4, result.usage.CacheCreationInputTokens)
	require.NotNil(t, result.firstTokenMs)

	// 确保有透传输出
	require.Contains(t, rec.Body.String(), "data:")
placeholder

// --- 流式 happy path 测试 ---

// TestStreamUpstreamResponse_NormalComplete
// 验证：正常流式转发完成时，数据正确透传、usage 正确收集、clientDisconnect=false
func TestStreamUpstreamResponse_NormalComplete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		fmt.Fprintln(pw, `event: message_start`)
		fmt.Fprintln(pw, `data: {"type":"message_start","message":{"usage":{"input_tokens":10placeholderplaceholderplaceholder`)
		fmt.Fprintln(pw, "")
		fmt.Fprintln(pw, `event: content_block_delta`)
		fmt.Fprintln(pw, `data: {"type":"content_block_delta","delta":{"text":"hello"placeholderplaceholder`)
		fmt.Fprintln(pw, "")
		fmt.Fprintln(pw, `event: message_delta`)
		fmt.Fprintln(pw, `data: {"type":"message_delta","usage":{"output_tokens":5placeholderplaceholder`)
		fmt.Fprintln(pw, "")
placeholder()

	result := svc.streamUpstreamResponse(c, resp, time.Now())
	_ = pr.Close()

	require.NotNil(t, result)
	require.False(t, result.clientDisconnect, "normal completion should not set clientDisconnect")
	require.NotNil(t, result.usage)
	require.Equal(t, 5, result.usage.OutputTokens, "should collect output_tokens from message_delta")
	require.NotNil(t, result.firstTokenMs, "should record first token time")

	// 验证数据被透传到客户端
	body := rec.Body.String()
	require.Contains(t, body, "event: message_start")
	require.Contains(t, body, "content_block_delta")
	require.Contains(t, body, "message_delta")
placeholder

// TestHandleGeminiStreamingResponse_NormalComplete
// 验证：正常 Gemini 流式转发，数据正确透传、usage 正确收集
func TestHandleGeminiStreamingResponse_NormalComplete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		// 第一个 chunk（部分内容）
		fmt.Fprintln(pw, `data: {"candidates":[{"content":{"parts":[{"text":"Hello"placeholder]placeholderplaceholder],"usageMetadata":{"promptTokenCount":10,"candidatesTokenCount":3placeholderplaceholder`)
		fmt.Fprintln(pw, "")
		// 第二个 chunk（最终内容+完整 usage）
		fmt.Fprintln(pw, `data: {"candidates":[{"content":{"parts":[{"text":" world"placeholder]placeholder,"finishReason":"STOP"placeholder],"usageMetadata":{"promptTokenCount":10,"candidatesTokenCount":8,"cachedContentTokenCount":2placeholderplaceholder`)
		fmt.Fprintln(pw, "")
placeholder()

	result, err := svc.handleGeminiStreamingResponse(c, resp, time.Now())
	_ = pr.Close()

placeholder
	require.NotNil(t, result)
	require.False(t, result.clientDisconnect, "normal completion should not set clientDisconnect")
	require.NotNil(t, result.usage)
	// Gemini usage: promptTokenCount=10, candidatesTokenCount=8, cachedContentTokenCount=2
	// → InputTokens=10-2=8, OutputTokens=8, CacheReadInputTokens=2
	require.Equal(t, 8, result.usage.InputTokens)
	require.Equal(t, 8, result.usage.OutputTokens)
	require.Equal(t, 2, result.usage.CacheReadInputTokens)
	require.NotNil(t, result.firstTokenMs, "should record first token time")

	// 验证数据被透传到客户端
	body := rec.Body.String()
	require.Contains(t, body, "Hello")
	require.Contains(t, body, "world")
	// 不应包含错误事件
	require.NotContains(t, body, "event: error")
placeholder

// TestHandleClaudeStreamingResponse_NormalComplete
// 验证：正常 Claude 流式转发（Gemini→Claude 转换），数据正确转换并输出
func TestHandleClaudeStreamingResponse_NormalComplete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		// v1internal 包装格式：Gemini 数据嵌套在 "response" 字段下
		// ProcessLine 先尝试反序列化为 V1InternalResponse，裸格式会导致 Response.UsageMetadata 为空
		fmt.Fprintln(pw, `data: {"response":{"candidates":[{"content":{"parts":[{"text":"Hi there"placeholder]placeholder,"finishReason":"STOP"placeholder],"usageMetadata":{"promptTokenCount":5,"candidatesTokenCount":3placeholderplaceholderplaceholder`)
		fmt.Fprintln(pw, "")
placeholder()

	result, err := svc.handleClaudeStreamingResponse(c, resp, time.Now(), "claude-sonnet-4-5")
	_ = pr.Close()

placeholder
	require.NotNil(t, result)
	require.False(t, result.clientDisconnect, "normal completion should not set clientDisconnect")
	require.NotNil(t, result.usage)
	// Gemini→Claude 转换的 usage：promptTokenCount=5→InputTokens=5, candidatesTokenCount=3→OutputTokens=3
	require.Equal(t, 5, result.usage.InputTokens)
	require.Equal(t, 3, result.usage.OutputTokens)
	require.NotNil(t, result.firstTokenMs, "should record first token time")

	// 验证输出是 Claude SSE 格式（processor 会转换）
	body := rec.Body.String()
	require.Contains(t, body, "event: message_start", "should contain Claude message_start event")
	require.Contains(t, body, "event: message_stop", "should contain Claude message_stop event")
	// 不应包含错误事件
	require.NotContains(t, body, "event: error")
placeholder

// TestHandleGeminiStreamingResponse_ThoughtsTokenCount
// 验证：Gemini 流式转发时 thoughtsTokenCount 被计入 OutputTokens
func TestHandleGeminiStreamingResponse_ThoughtsTokenCount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		fmt.Fprintln(pw, `data: {"candidates":[{"content":{"parts":[{"text":"Hello"placeholder]placeholderplaceholder],"usageMetadata":{"promptTokenCount":100,"candidatesTokenCount":20,"thoughtsTokenCount":50placeholderplaceholder`)
		fmt.Fprintln(pw, "")
		fmt.Fprintln(pw, `data: {"candidates":[{"content":{"parts":[{"text":" world"placeholder]placeholder,"finishReason":"STOP"placeholder],"usageMetadata":{"promptTokenCount":100,"candidatesTokenCount":30,"thoughtsTokenCount":80,"cachedContentTokenCount":10placeholderplaceholder`)
		fmt.Fprintln(pw, "")
placeholder()

	result, err := svc.handleGeminiStreamingResponse(c, resp, time.Now())
	_ = pr.Close()

placeholder
	require.NotNil(t, result)
	require.NotNil(t, result.usage)
	// promptTokenCount=100, cachedContentTokenCount=10 → InputTokens=90
	require.Equal(t, 90, result.usage.InputTokens)
	// candidatesTokenCount=30 + thoughtsTokenCount=80 → OutputTokens=110
	require.Equal(t, 110, result.usage.OutputTokens)
	require.Equal(t, 10, result.usage.CacheReadInputTokens)
placeholder

// TestHandleClaudeStreamingResponse_ThoughtsTokenCount
// 验证：Gemini→Claude 流式转换时 thoughtsTokenCount 被计入 OutputTokens
func TestHandleClaudeStreamingResponse_ThoughtsTokenCount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		fmt.Fprintln(pw, `data: {"response":{"candidates":[{"content":{"parts":[{"text":"Hi"placeholder]placeholder,"finishReason":"STOP"placeholder],"usageMetadata":{"promptTokenCount":50,"candidatesTokenCount":10,"thoughtsTokenCount":25placeholderplaceholderplaceholder`)
		fmt.Fprintln(pw, "")
placeholder()

	result, err := svc.handleClaudeStreamingResponse(c, resp, time.Now(), "gemini-2.5-pro")
	_ = pr.Close()

placeholder
	require.NotNil(t, result)
	require.NotNil(t, result.usage)
	// promptTokenCount=50 → InputTokens=50
	require.Equal(t, 50, result.usage.InputTokens)
	// candidatesTokenCount=10 + thoughtsTokenCount=25 → OutputTokens=35
	require.Equal(t, 35, result.usage.OutputTokens)
placeholder

// --- 流式客户端断开检测测试 ---

// TestStreamUpstreamResponse_ClientDisconnectDrainsUsage
// 验证：客户端写入失败后，streamUpstreamResponse 继续读取上游以收集 usage
func TestStreamUpstreamResponse_ClientDisconnectDrainsUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	c.Writer = &antigravityFailingWriter{ResponseWriter: c.Writer, failAfter: 0placeholder

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		fmt.Fprintln(pw, `event: message_start`)
		fmt.Fprintln(pw, `data: {"type":"message_start","message":{"usage":{"input_tokens":10placeholderplaceholderplaceholder`)
		fmt.Fprintln(pw, "")
		fmt.Fprintln(pw, `event: message_delta`)
		fmt.Fprintln(pw, `data: {"type":"message_delta","usage":{"output_tokens":20placeholderplaceholder`)
		fmt.Fprintln(pw, "")
placeholder()

	result := svc.streamUpstreamResponse(c, resp, time.Now())
	_ = pr.Close()

	require.NotNil(t, result)
	require.True(t, result.clientDisconnect)
	require.NotNil(t, result.usage)
	require.Equal(t, 20, result.usage.OutputTokens)
placeholder

// TestStreamUpstreamResponse_ContextCanceled
// 验证：context 取消时返回 usage 且标记 clientDisconnect
func TestStreamUpstreamResponse_ContextCanceled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil).WithContext(ctx)

	resp := &http.Response{StatusCode: http.StatusOK, Body: cancelReadCloser{placeholder, Header: http.Header{placeholderplaceholder

	result := svc.streamUpstreamResponse(c, resp, time.Now())

	require.NotNil(t, result)
	require.True(t, result.clientDisconnect)
	require.NotContains(t, rec.Body.String(), "event: error")
placeholder

// TestStreamUpstreamResponse_Timeout
// 验证：上游超时时返回已收集的 usage
func TestStreamUpstreamResponse_Timeout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{StreamDataIntervalTimeout: 1, MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	result := svc.streamUpstreamResponse(c, resp, time.Now())
	_ = pw.Close()
	_ = pr.Close()

	require.NotNil(t, result)
	require.False(t, result.clientDisconnect)
placeholder

// TestStreamUpstreamResponse_TimeoutAfterClientDisconnect
// 验证：客户端断开后上游超时，返回 usage 并标记 clientDisconnect
func TestStreamUpstreamResponse_TimeoutAfterClientDisconnect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{StreamDataIntervalTimeout: 1, MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	c.Writer = &antigravityFailingWriter{ResponseWriter: c.Writer, failAfter: 0placeholder

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	go func() {
		fmt.Fprintln(pw, `data: {"type":"message_start","message":{"usage":{"input_tokens":5placeholderplaceholderplaceholder`)
		fmt.Fprintln(pw, "")
		// 不关闭 pw → 等待超时
placeholder()

	result := svc.streamUpstreamResponse(c, resp, time.Now())
	_ = pw.Close()
	_ = pr.Close()

	require.NotNil(t, result)
	require.True(t, result.clientDisconnect)
placeholder

// TestHandleGeminiStreamingResponse_ClientDisconnect
// 验证：Gemini 流式转发中客户端断开后继续 drain 上游
func TestHandleGeminiStreamingResponse_ClientDisconnect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	c.Writer = &antigravityFailingWriter{ResponseWriter: c.Writer, failAfter: 0placeholder

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		fmt.Fprintln(pw, `data: {"candidates":[{"content":{"parts":[{"text":"hi"placeholder]placeholderplaceholder],"usageMetadata":{"promptTokenCount":5,"candidatesTokenCount":10placeholderplaceholder`)
		fmt.Fprintln(pw, "")
placeholder()

	result, err := svc.handleGeminiStreamingResponse(c, resp, time.Now())
	_ = pr.Close()

placeholder
	require.NotNil(t, result)
	require.True(t, result.clientDisconnect)
	require.NotContains(t, rec.Body.String(), "write_failed")
placeholder

// TestHandleGeminiStreamingResponse_ContextCanceled
// 验证：context 取消时不注入错误事件
func TestHandleGeminiStreamingResponse_ContextCanceled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil).WithContext(ctx)

	resp := &http.Response{StatusCode: http.StatusOK, Body: cancelReadCloser{placeholder, Header: http.Header{placeholderplaceholder

	result, err := svc.handleGeminiStreamingResponse(c, resp, time.Now())

placeholder
	require.NotNil(t, result)
	require.True(t, result.clientDisconnect)
	require.NotContains(t, rec.Body.String(), "event: error")
placeholder

// TestHandleClaudeStreamingResponse_ClientDisconnect
// 验证：Claude 流式转发中客户端断开后继续 drain 上游
func TestHandleClaudeStreamingResponse_ClientDisconnect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	c.Writer = &antigravityFailingWriter{ResponseWriter: c.Writer, failAfter: 0placeholder

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		// v1internal 包装格式
		fmt.Fprintln(pw, `data: {"response":{"candidates":[{"content":{"parts":[{"text":"hello"placeholder]placeholder,"finishReason":"STOP"placeholder],"usageMetadata":{"promptTokenCount":8,"candidatesTokenCount":15placeholderplaceholderplaceholder`)
		fmt.Fprintln(pw, "")
placeholder()

	result, err := svc.handleClaudeStreamingResponse(c, resp, time.Now(), "claude-sonnet-4-5")
	_ = pr.Close()

placeholder
	require.NotNil(t, result)
	require.True(t, result.clientDisconnect)
placeholder

// TestHandleClaudeStreamingResponse_EmptyStream
// 验证：上游只返回无法解析的 SSE 行时，触发 UpstreamFailoverError 而不是向客户端发出残缺流
func TestHandleClaudeStreamingResponse_EmptyStream(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{placeholderplaceholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		// 所有行均为无法 JSON 解析的内容，ProcessLine 全部返回 nil
		fmt.Fprintln(pw, "data: not-valid-json")
		fmt.Fprintln(pw, "")
		fmt.Fprintln(pw, "data: also-invalid")
		fmt.Fprintln(pw, "")
placeholder()

	_, err := svc.handleClaudeStreamingResponse(c, resp, time.Now(), "claude-sonnet-4-5")
	_ = pr.Close()

	// 应当返回 UpstreamFailoverError 而非 nil，以便上层触发 failover
placeholder
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.True(t, failoverErr.RetryableOnSameAccount)

	// 客户端不应收到任何 SSE 事件（既无 message_start 也无 message_stop）
	body := rec.Body.String()
	require.NotContains(t, body, "event: message_start")
	require.NotContains(t, body, "event: message_stop")
	require.NotContains(t, body, "event: message_delta")
placeholder

// TestHandleClaudeStreamingResponse_ContextCanceled
// 验证：context 取消时不注入错误事件
func TestHandleClaudeStreamingResponse_ContextCanceled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newAntigravityTestService(&config.Config{
		Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholder,
placeholder)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil).WithContext(ctx)

	resp := &http.Response{StatusCode: http.StatusOK, Body: cancelReadCloser{placeholder, Header: http.Header{placeholderplaceholder

	result, err := svc.handleClaudeStreamingResponse(c, resp, time.Now(), "claude-sonnet-4-5")

placeholder
	require.NotNil(t, result)
	require.True(t, result.clientDisconnect)
	require.NotContains(t, rec.Body.String(), "event: error")
placeholder

// TestExtractSSEUsage 验证 extractSSEUsage 从 SSE data 行正确提取 usage
func TestExtractSSEUsage(t *testing.T) {
	svc := &AntigravityGatewayService{placeholder
	tests := []struct {
		name     string
		line     string
		expected ClaudeUsage
placeholder{
		{
			name:     "message_delta with output_tokens",
			line:     `data: {"type":"message_delta","usage":{"output_tokens":42placeholderplaceholder`,
			expected: ClaudeUsage{OutputTokens: 42placeholder,
	placeholder,
		{
			name:     "non-data line ignored",
			line:     `event: message_start`,
			expected: ClaudeUsage{placeholder,
	placeholder,
		{
			name:     "top-level usage with all fields",
			line:     `data: {"usage":{"input_tokens":10,"output_tokens":20,"cache_read_input_tokens":5,"cache_creation_input_tokens":3placeholderplaceholder`,
			expected: ClaudeUsage{InputTokens: 10, OutputTokens: 20, CacheReadInputTokens: 5, CacheCreationInputTokens: 3placeholder,
	placeholder,
		{
			// Anthropic message_start 把 usage 嵌套在 message.usage 下，
			// 必须从这里提取输入侧字段（含 cache_read/cache_creation_input_tokens）。
			name:     "message_start nested usage with input/cache tokens",
			line:     `data: {"type":"message_start","message":{"id":"msg_01","usage":{"input_tokens":35576,"cache_creation_input_tokens":0,"cache_read_input_tokens":12000,"output_tokens":1placeholderplaceholderplaceholder`,
			expected: ClaudeUsage{InputTokens: 35576, OutputTokens: 1, CacheReadInputTokens: 12000placeholder,
	placeholder,
		{
			// message_start.message.usage.cache_creation 内的 5m/1h 明细也要解析。
			name:     "message_start nested usage with cache_creation breakdown",
			line:     `data: {"type":"message_start","message":{"usage":{"input_tokens":100,"cache_creation":{"ephemeral_5m_input_tokens":30,"ephemeral_1h_input_tokens":70placeholderplaceholderplaceholderplaceholder`,
			expected: ClaudeUsage{InputTokens: 100, CacheCreation5mTokens: 30, CacheCreation1hTokens: 70placeholder,
	placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usage := &ClaudeUsage{placeholder
			svc.extractSSEUsage(tt.line, usage)
			require.Equal(t, tt.expected, *usage)
	placeholder)
placeholder
placeholder

// TestExtractSSEUsage_StreamingSequence 复现 issue #2332：完整的 Anthropic streaming
// 序列（message_start → message_delta）必须把两类事件中的 usage 字段都汇入同一份累计值，
// 否则透传账号产出的 usage_logs 会出现 input_tokens=0、仅有 output_tokens 的"残缺"记录。
func TestExtractSSEUsage_StreamingSequence(t *testing.T) {
	svc := &AntigravityGatewayService{placeholder
	usage := &ClaudeUsage{placeholder

	// 1) message_start：携带完整输入侧 usage（input_tokens + cache_read）
	svc.extractSSEUsage(
		`data: {"type":"message_start","message":{"id":"msg_01","type":"message","role":"assistant","content":[],"model":"claude-opus-4-6","stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":35576,"cache_creation_input_tokens":0,"cache_read_input_tokens":12000,"output_tokens":1placeholderplaceholderplaceholder`,
		usage,
	)
	// 2) message_delta：流结束时只带 output_tokens（无 input_tokens 字段）
	svc.extractSSEUsage(
		`data: {"type":"message_delta","delta":{"stop_reason":"end_turn","stop_sequence":nullplaceholder,"usage":{"output_tokens":816placeholderplaceholder`,
		usage,
	)

	require.Equal(t, 35576, usage.InputTokens, "message_start 的 input_tokens 必须被记录，否则记账会缺失输入侧 token (#2332)")
	require.Equal(t, 12000, usage.CacheReadInputTokens, "message_start 的 cache_read_input_tokens 必须被记录")
	require.Equal(t, 816, usage.OutputTokens, "message_delta 的最终 output_tokens 必须被记录")
placeholder

// TestAntigravityClientWriter 验证 antigravityClientWriter 的断开检测
func TestAntigravityClientWriter(t *testing.T) {
	t.Run("normal write succeeds", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		flusher, _ := c.Writer.(http.Flusher)
		cw := newAntigravityClientWriter(c.Writer, flusher, "test")

		ok := cw.Write([]byte("hello"))
		require.True(t, ok)
		require.False(t, cw.Disconnected())
		require.Contains(t, rec.Body.String(), "hello")
placeholder)

	t.Run("write failure marks disconnected", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		fw := &antigravityFailingWriter{ResponseWriter: c.Writer, failAfter: 0placeholder
		flusher, _ := c.Writer.(http.Flusher)
		cw := newAntigravityClientWriter(fw, flusher, "test")

		ok := cw.Write([]byte("hello"))
		require.False(t, ok)
		require.True(t, cw.Disconnected())
placeholder)

	t.Run("subsequent writes are no-op", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		fw := &antigravityFailingWriter{ResponseWriter: c.Writer, failAfter: 0placeholder
		flusher, _ := c.Writer.(http.Flusher)
		cw := newAntigravityClientWriter(fw, flusher, "test")

		cw.Write([]byte("first"))
		ok := cw.Fprintf("second %d", 2)
		require.False(t, ok)
		require.True(t, cw.Disconnected())
placeholder)
placeholder

// TestUnwrapV1InternalResponse 测试 unwrapV1InternalResponse 的各种输入场景
func TestUnwrapV1InternalResponse(t *testing.T) {
	svc := &AntigravityGatewayService{placeholder

	// 构造 >50KB 的大型 JSON
	largePadding := strings.Repeat("x", 50*1024)
	largeInput := []byte(fmt.Sprintf(`{"response":{"id":"big","pad":"%s"placeholderplaceholder`, largePadding))
	largeExpected := fmt.Sprintf(`{"id":"big","pad":"%s"placeholder`, largePadding)

	tests := []struct {
		name     string
		input    []byte
		expected string
		wantErr  bool
placeholder{
		{
			name:     "正常 response 包装",
			input:    []byte(`{"response":{"id":"123","content":"hello"placeholderplaceholder`),
			expected: `{"id":"123","content":"hello"placeholder`,
	placeholder,
		{
			name:     "无 response 透传",
			input:    []byte(`{"id":"456"placeholder`),
			expected: `{"id":"456"placeholder`,
	placeholder,
		{
			name:     "空 JSON",
			input:    []byte(`{placeholder`),
			expected: `{placeholder`,
	placeholder,
		{
			name:     "response 为 null",
			input:    []byte(`{"response":nullplaceholder`),
			expected: `null`,
	placeholder,
		{
			name:     "response 为基础类型 string",
			input:    []byte(`{"response":"hello"placeholder`),
			expected: `"hello"`,
	placeholder,
		{
			name:     "非法 JSON",
			input:    []byte(`not json`),
			expected: `not json`,
	placeholder,
		{
			name:     "嵌套 response 只解一层",
			input:    []byte(`{"response":{"response":{"inner":trueplaceholderplaceholderplaceholder`),
			expected: `{"response":{"inner":trueplaceholderplaceholder`,
	placeholder,
		{
			name:     "大型 JSON >50KB",
			input:    largeInput,
			expected: largeExpected,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.unwrapV1InternalResponse(tt.input)
			if tt.wantErr {
			placeholder
				return
		placeholder
		placeholder
			require.Equal(t, tt.expected, strings.TrimSpace(string(got)))
	placeholder)
placeholder
placeholder

// --- unwrapV1InternalResponse benchmark 对照组 ---

// unwrapV1InternalResponseOld 旧实现：Unmarshal+Marshal 双重开销（仅用于 benchmark 对照）
func unwrapV1InternalResponseOld(body []byte) ([]byte, error) {
	var outer map[string]any
	if err := json.Unmarshal(body, &outer); err != nil {
		return nil, err
placeholder
	if resp, ok := outer["response"]; ok {
		return json.Marshal(resp)
placeholder
	return body, nil
placeholder

func BenchmarkUnwrapV1Internal_Old_Small(b *testing.B) {
	body := []byte(`{"response":{"candidates":[{"content":{"parts":[{"text":"hello world"placeholder]placeholderplaceholder],"usageMetadata":{"promptTokenCount":10,"candidatesTokenCount":5placeholderplaceholderplaceholder`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = unwrapV1InternalResponseOld(body)
placeholder
placeholder

func BenchmarkUnwrapV1Internal_New_Small(b *testing.B) {
	body := []byte(`{"response":{"candidates":[{"content":{"parts":[{"text":"hello world"placeholder]placeholderplaceholder],"usageMetadata":{"promptTokenCount":10,"candidatesTokenCount":5placeholderplaceholderplaceholder`)
	svc := &AntigravityGatewayService{placeholder
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.unwrapV1InternalResponse(body)
placeholder
placeholder

func BenchmarkUnwrapV1Internal_Old_Large(b *testing.B) {
	body := generateLargeUnwrapJSON(10 * 1024) // ~10KB
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = unwrapV1InternalResponseOld(body)
placeholder
placeholder

func BenchmarkUnwrapV1Internal_New_Large(b *testing.B) {
	body := generateLargeUnwrapJSON(10 * 1024) // ~10KB
	svc := &AntigravityGatewayService{placeholder
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.unwrapV1InternalResponse(body)
placeholder
placeholder

// generateLargeUnwrapJSON 生成指定最小大小的包含 response 包装的 JSON
func generateLargeUnwrapJSON(minSize int) []byte {
	parts := make([]map[string]string, 0)
	current := 0
	for current < minSize {
		text := fmt.Sprintf("这是第 %d 段内容，用于填充 JSON 到目标大小。", len(parts)+1)
		parts = append(parts, map[string]string{"text": textplaceholder)
		current += len(text) + 20 // 估算 JSON 编码开销
placeholder
	inner := map[string]any{
		"candidates": []map[string]any{
			{"content": map[string]any{"parts": partsplaceholderplaceholder,
	placeholder,
		"usageMetadata": map[string]any{
			"promptTokenCount":     100,
			"candidatesTokenCount": 50,
	placeholder,
placeholder
	outer := map[string]any{"response": innerplaceholder
	b, _ := json.Marshal(outer)
	return b
placeholder
