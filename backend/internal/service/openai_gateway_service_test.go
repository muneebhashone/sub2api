package service

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/model"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/cespare/xxhash/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

// 编译期接口断言
var _ AccountRepository = (*stubOpenAIAccountRepo)(nil)
var _ GatewayCache = (*stubGatewayCache)(nil)

type stubOpenAIAccountRepo struct {
	AccountRepository
	accounts []Account
placeholder

type tempUnschedulableOpenAIAccountRepo struct {
	stubOpenAIAccountRepo
	modelRateLimitAccountID int64
	modelRateLimitKey       string
placeholder

func (r *tempUnschedulableOpenAIAccountRepo) SetModelRateLimit(_ context.Context, accountID int64, modelKey string, _ time.Time, _ ...string) error {
	r.modelRateLimitAccountID = accountID
	r.modelRateLimitKey = modelKey
	return nil
placeholder

type snapshotUpdateAccountRepo struct {
	stubOpenAIAccountRepo
	updateExtraCalls chan map[string]any
placeholder

func (r *snapshotUpdateAccountRepo) UpdateExtra(ctx context.Context, id int64, updates map[string]any) error {
	if r.updateExtraCalls != nil {
		copied := make(map[string]any, len(updates))
		for k, v := range updates {
			copied[k] = v
	placeholder
		r.updateExtraCalls <- copied
placeholder
	return nil
placeholder

func (r stubOpenAIAccountRepo) GetByID(ctx context.Context, id int64) (*Account, error) {
	for i := range r.accounts {
		if r.accounts[i].ID == id {
			return &r.accounts[i], nil
	placeholder
placeholder
	return nil, errors.New("account not found")
placeholder

func (r stubOpenAIAccountRepo) ListSchedulableByGroupIDAndPlatform(ctx context.Context, groupID int64, platform string) ([]Account, error) {
	var result []Account
	for _, acc := range r.accounts {
		if acc.Platform == platform {
			result = append(result, acc)
	placeholder
placeholder
	return result, nil
placeholder

func (r stubOpenAIAccountRepo) ListSchedulableByPlatform(ctx context.Context, platform string) ([]Account, error) {
	var result []Account
	for _, acc := range r.accounts {
		if acc.Platform == platform {
			result = append(result, acc)
	placeholder
placeholder
	return result, nil
placeholder

func (r stubOpenAIAccountRepo) ListSchedulableUngroupedByPlatform(ctx context.Context, platform string) ([]Account, error) {
	return r.ListSchedulableByPlatform(ctx, platform)
placeholder

func TestOpenAIGatewayService_ForwardAsAnthropic_TempUnschedulableReturnsFailoverWithoutCommit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"model":"gpt-5.4","max_tokens":32,"messages":[{"role":"user","content":"hello"placeholder],"stream":falseplaceholder`)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	upstreamBody := []byte(`{"error":{"message":"Our servers are currently overloaded. Please try again later."placeholderplaceholder`)
	upstream := &httpUpstreamRecorder{responses: []*http.Response{
		{
			StatusCode: http.StatusBadRequest,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(bytes.NewReader(upstreamBody)),
	placeholder,
		{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
			Body: io.NopCloser(strings.NewReader(strings.Join([]string{
				`data: {"type":"response.output_text.delta","delta":"ok"placeholder`,
				"",
				`data: {"type":"response.completed","response":{"id":"resp_second","status":"completed","output":[],"usage":{"input_tokens":1,"output_tokens":1placeholderplaceholderplaceholder`,
				"",
		placeholder, "\n"))),
	placeholder,
placeholderplaceholder
	repo := &tempUnschedulableOpenAIAccountRepo{placeholder
	rateLimitService := NewRateLimitService(repo, nil, &config.Config{placeholder, nil, nil)
	svc := &OpenAIGatewayService{
		cfg: &config.Config{Security: config.SecurityConfig{URLAllowlist: config.URLAllowlistConfig{
			Enabled: false, AllowInsecureHTTP: true,
	placeholderplaceholderplaceholder,
		httpUpstream:     upstream,
		rateLimitService: rateLimitService,
placeholder
	account := &Account{
		ID: 5099, Name: "temporary-unschedulable", Platform: PlatformOpenAI,
		Type: AccountTypeAPIKey, Concurrency: 1,
placeholder
			"api_key":                    "sk-test",
			"base_url":                   "http://upstream.example",
			"model_provider":             "env-openai",
			"temp_unschedulable_enabled": true,
			"temp_unschedulable_rules": []any{map[string]any{
				"error_code":       float64(http.StatusBadRequest),
				"keywords":         []any{"our servers are currently overloaded", "please try again later"placeholder,
				"duration_minutes": float64(1),
	placeholder
	placeholder,
placeholder

	_, err := svc.ForwardAsAnthropic(context.Background(), c, account, body, "", "")

	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusBadRequest, failoverErr.StatusCode)
	require.True(t, failoverErr.ShouldRetryNextAccount())
	require.Equal(t, account.ID, repo.modelRateLimitAccountID, "temporary unschedulability should exclude this account from reselection")
	require.Equal(t, "gpt-5.4", repo.modelRateLimitKey)
	require.False(t, IsResponseCommitted(c))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Empty(t, rec.Body.String())

	secondRec := httptest.NewRecorder()
	secondContext, _ := gin.CreateTestContext(secondRec)
	secondContext.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	secondContext.Request.Header.Set("Content-Type", "application/json")
	secondAccount := *account
	secondAccount.ID = 5100
	secondAccount.Name = "healthy-failover-account"
	result, secondErr := svc.ForwardAsAnthropic(context.Background(), secondContext, &secondAccount, body, "", "")
	require.NoError(t, secondErr)
	require.NotNil(t, result)
	require.Equal(t, "resp_second", result.ResponseID)
	require.NotEmpty(t, secondRec.Body.String())
placeholder

func TestFailoverOpenAIUpstreamHTTPError_NilContextSkipsTempUnschedulablePolicy(t *testing.T) {
	repo := &tempUnschedulableOpenAIAccountRepo{placeholder
	svc := &OpenAIGatewayService{
		rateLimitService: NewRateLimitService(repo, nil, &config.Config{placeholder, nil, nil),
placeholder
	account := &Account{
		ID: 5099, Platform: PlatformOpenAI, Type: AccountTypeAPIKey,
placeholder
			"temp_unschedulable_enabled": true,
			"temp_unschedulable_rules": []any{map[string]any{
				"error_code":       float64(http.StatusBadRequest),
				"keywords":         []any{"servers are currently overloaded"placeholder,
				"duration_minutes": float64(1),
	placeholder
	placeholder,
placeholder
	body := []byte(`{"error":{"message":"Our servers are currently overloaded."placeholderplaceholder`)
	resp := &http.Response{StatusCode: http.StatusBadRequest, Header: http.Header{placeholderplaceholder

	got := svc.failoverOpenAIUpstreamHTTPError(
		context.Background(), nil, account, resp, body,
		"Our servers are currently overloaded.", "gpt-5.4",
	)

	require.Nil(t, got)
	require.Zero(t, repo.modelRateLimitAccountID)
	require.Empty(t, repo.modelRateLimitKey)
placeholder

type groupAwareStubOpenAIAccountRepo struct {
	stubOpenAIAccountRepo
placeholder

func (r groupAwareStubOpenAIAccountRepo) ListSchedulableByGroupIDAndPlatform(ctx context.Context, groupID int64, platform string) ([]Account, error) {
	var result []Account
	for _, acc := range r.accounts {
		if acc.Platform == platform && openAIStickyAccountMatchesGroup(&acc, &groupID) {
			result = append(result, acc)
	placeholder
placeholder
	return result, nil
placeholder

func (r groupAwareStubOpenAIAccountRepo) ListSchedulableUngroupedByPlatform(ctx context.Context, platform string) ([]Account, error) {
	var result []Account
	for _, acc := range r.accounts {
		if acc.Platform == platform && openAIStickyAccountMatchesGroup(&acc, nil) {
			result = append(result, acc)
	placeholder
placeholder
	return result, nil
placeholder

type stubConcurrencyCache struct {
	ConcurrencyCache
	loadBatchErr    error
	loadMap         map[int64]*AccountLoadInfo
	acquireResults  map[int64]bool
	waitCounts      map[int64]int
	skipDefaultLoad bool
placeholder

type cancelReadCloser struct{placeholder

func (c cancelReadCloser) Read(p []byte) (int, error) { return 0, context.Canceled placeholder
func (c cancelReadCloser) Close() error               { return nil placeholder

type errReadCloser struct {
	err error
placeholder

func (r errReadCloser) Read([]byte) (int, error) { return 0, r.err placeholder
func (r errReadCloser) Close() error             { return nil placeholder

type openAIStreamReadThenErrorCloser struct {
	reader *strings.Reader
	err    error
placeholder

func (r *openAIStreamReadThenErrorCloser) Read(p []byte) (int, error) {
	if r.reader != nil && r.reader.Len() > 0 {
		return r.reader.Read(p)
placeholder
	return 0, r.err
placeholder

func (r *openAIStreamReadThenErrorCloser) Close() error { return nil placeholder

type failingGinWriter struct {
	gin.ResponseWriter
	failAfter int
	writes    int
placeholder

func (w *failingGinWriter) Write(p []byte) (int, error) {
	if w.writes >= w.failAfter {
		return 0, errors.New("write failed")
placeholder
	w.writes++
	return w.ResponseWriter.Write(p)
placeholder

func (c stubConcurrencyCache) AcquireAccountSlot(ctx context.Context, accountID int64, maxConcurrency int, requestID string) (bool, error) {
	if c.acquireResults != nil {
		if result, ok := c.acquireResults[accountID]; ok {
			return result, nil
	placeholder
placeholder
	return true, nil
placeholder

func (c stubConcurrencyCache) ReleaseAccountSlot(ctx context.Context, accountID int64, requestID string) error {
	return nil
placeholder

func (c stubConcurrencyCache) GetAccountsLoadBatch(ctx context.Context, accounts []AccountWithConcurrency) (map[int64]*AccountLoadInfo, error) {
	if c.loadBatchErr != nil {
		return nil, c.loadBatchErr
placeholder
	out := make(map[int64]*AccountLoadInfo, len(accounts))
	if c.skipDefaultLoad && c.loadMap != nil {
		for _, acc := range accounts {
			if load, ok := c.loadMap[acc.ID]; ok {
				out[acc.ID] = load
		placeholder
	placeholder
		return out, nil
placeholder
	for _, acc := range accounts {
		if c.loadMap != nil {
			if load, ok := c.loadMap[acc.ID]; ok {
				out[acc.ID] = load
				continue
		placeholder
	placeholder
		out[acc.ID] = &AccountLoadInfo{AccountID: acc.ID, LoadRate: 0placeholder
placeholder
	return out, nil
placeholder

func TestOpenAIGatewayService_GenerateSessionHash_Priority(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)

	svc := &OpenAIGatewayService{placeholder

	bodyWithKey := []byte(`{"prompt_cache_key":"ses_aaa"placeholder`)

	// 1) session_id header wins
	c.Request.Header.Set("session_id", "sess-123")
	c.Request.Header.Set("conversation_id", "conv-456")
	h1 := svc.GenerateSessionHash(c, bodyWithKey)
	if h1 == "" {
		t.Fatalf("expected non-empty hash")
placeholder

	// 2) conversation_id used when session_id absent
	c.Request.Header.Del("session_id")
	h2 := svc.GenerateSessionHash(c, bodyWithKey)
	if h2 == "" {
		t.Fatalf("expected non-empty hash")
placeholder
	if h1 == h2 {
		t.Fatalf("expected different hashes for different keys")
placeholder

	// 3) prompt_cache_key used when both headers absent
	c.Request.Header.Del("conversation_id")
	h3 := svc.GenerateSessionHash(c, bodyWithKey)
	if h3 == "" {
		t.Fatalf("expected non-empty hash")
placeholder
	if h2 == h3 {
		t.Fatalf("expected different hashes for different keys")
placeholder

	// 4) empty when no signals
	h4 := svc.GenerateSessionHash(c, []byte(`{placeholder`))
	if h4 != "" {
		t.Fatalf("expected empty hash when no signals")
placeholder
placeholder

func TestOpenAIGatewayService_ClientSessionHeaderPriority(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	c.Set("api_key", &APIKey{ID: 901, Group: &Group{Platform: PlatformGrokplaceholderplaceholder)

	headers := []struct {
		name  string
		value string
placeholder{
		{name: "session_id", value: "generic-session"placeholder,
		{name: "conversation_id", value: "generic-conversation"placeholder,
		{name: openCodeSessionAffinityHeader, value: "opencode-affinity"placeholder,
		{name: openCodeSessionIDHeader, value: "opencode-session-id"placeholder,
		{name: openCodeNativeSessionHeader, value: "opencode-native-session"placeholder,
		{name: codeBuddyConversationHeader, value: "codebuddy-conversation"placeholder,
		{name: grokConversationIDHeader, value: "grok-conversation"placeholder,
placeholder
	for _, header := range headers {
		c.Request.Header.Set(header.name, header.value)
placeholder

	svc := &OpenAIGatewayService{placeholder
	body := []byte(`{"prompt_cache_key":"body-session"placeholder`)
	for _, header := range headers {
		require.Equal(t, header.value, svc.ExtractSessionID(c, body), header.name)
		require.Equal(t, fmt.Sprintf("%016x", xxhash.Sum64String(header.value)), svc.GenerateExplicitSessionHash(c, body), header.name)
		if header.name != grokConversationIDHeader {
			require.Equal(t, header.value, explicitOpenAISessionID(c, body), header.name)
	placeholder
		c.Request.Header.Del(header.name)
placeholder
	require.Equal(t, "body-session", svc.ExtractSessionID(c, body))
placeholder

func TestOpenAIGatewayService_ClientSessionHeadersIgnorePerRequestIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	for name, value := range map[string]string{
		"X-Conversation-Request-ID": "request-rotates-every-turn",
		"X-Conversation-Message-ID": "message-rotates-every-turn",
		"X-Request-ID":              "generic-request-id",
placeholder {
		c.Request.Header.Set(name, value)
placeholder

	svc := &OpenAIGatewayService{placeholder
	require.Empty(t, explicitOpenAIHeaderSessionID(c))
	require.Empty(t, svc.ExtractSessionID(c, nil))
	require.Empty(t, svc.GenerateExplicitSessionHash(c, nil))
placeholder

func TestOpenAIGatewayService_GenerateSessionHash_UsesXXHash64(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)

	c.Request.Header.Set("session_id", "sess-fixed-value")
	svc := &OpenAIGatewayService{placeholder

	got := svc.GenerateSessionHash(c, nil)
	want := fmt.Sprintf("%016x", xxhash.Sum64String("sess-fixed-value"))
	require.Equal(t, want, got)
placeholder

func TestOpenAIGatewayService_GenerateSessionHash_AttachesLegacyHashToContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)

	c.Request.Header.Set("session_id", "sess-legacy-check")
	svc := &OpenAIGatewayService{placeholder

	sessionHash := svc.GenerateSessionHash(c, nil)
	require.NotEmpty(t, sessionHash)
	require.NotNil(t, c.Request)
	require.NotNil(t, c.Request.Context())
	require.NotEmpty(t, openAILegacySessionHashFromContext(c.Request.Context()))
placeholder

func TestExtractOpenAIResponseIDFromJSONBytes(t *testing.T) {
	require.Equal(t, "resp_json", extractOpenAIResponseIDFromJSONBytes([]byte(`{"id":"resp_json"placeholder`)))
	require.Equal(t, "resp_sse", extractOpenAIResponseIDFromJSONBytes([]byte(`{"type":"response.completed","response":{"id":"resp_sse"placeholderplaceholder`)))
	require.Empty(t, extractOpenAIResponseIDFromJSONBytes([]byte(`{"response":{placeholderplaceholder`)))
	require.Empty(t, extractOpenAIResponseIDFromJSONBytes([]byte(`not-json`)))
placeholder

// 复现 #4386：gpt-image-2 /v1/images/edits 的 usage 携带 input_tokens_details.image_tokens，
// 提取器须将图片输入 token 单独填入 ImageInputTokens（此前被丢弃并入 InputTokens 按文本价计费）。
func TestExtractOpenAIUsage_CapturesImageInputTokens(t *testing.T) {
	body := []byte(`{"usage":{"input_tokens":371,"input_tokens_details":{"image_tokens":352,"text_tokens":19placeholder,"output_tokens":439,"output_tokens_details":{"image_tokens":439,"text_tokens":0placeholder,"total_tokens":810placeholderplaceholder`)
	usage, ok := extractOpenAIUsageFromJSONBytes(body)
	require.True(t, ok)
	require.Equal(t, 371, usage.InputTokens)
	require.Equal(t, 352, usage.ImageInputTokens)
	require.Equal(t, 439, usage.OutputTokens)
	require.Equal(t, 439, usage.ImageOutputTokens)

	// prompt_tokens_details 回退路径（部分上游用 prompt_tokens 口径）。
	promptStyle := []byte(`{"usage":{"prompt_tokens":100,"prompt_tokens_details":{"image_tokens":80placeholderplaceholderplaceholder`)
	pu, ok := extractOpenAIUsageFromJSONBytes(promptStyle)
	require.True(t, ok)
	require.Equal(t, 100, pu.InputTokens)
	require.Equal(t, 80, pu.ImageInputTokens)

	// 纯文本请求：无 image_tokens 时 ImageInputTokens 为 0，行为不变。
	textOnly := []byte(`{"usage":{"input_tokens":50,"output_tokens":10placeholderplaceholder`)
	tu, ok := extractOpenAIUsageFromJSONBytes(textOnly)
	require.True(t, ok)
	require.Zero(t, tu.ImageInputTokens)
placeholder

func TestOpenAIGatewayService_BindHTTPResponseAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	groupID := int64(4201)
	c.Set("api_key", &APIKey{ID: 501, GroupID: &groupIDplaceholder)

	svc := &OpenAIGatewayService{placeholder
	account := &Account{ID: 37001, Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder
	svc.bindHTTPResponseAccount(context.Background(), c, account, "resp_http_001")

	got, err := svc.getOpenAIWSStateStore().GetResponseAccount(context.Background(), groupID, "resp_http_001")
placeholder
	require.Equal(t, account.ID, got)
placeholder

func TestOpenAIGatewayService_GenerateExplicitSessionHash_SkipsContentFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &OpenAIGatewayService{placeholder
	body := []byte(`{"model":"gpt-image-2","prompt":"draw a cat"placeholder`)

	t.Run("stateless image body stays unstuck", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", nil)

		require.Empty(t, svc.GenerateExplicitSessionHash(c, body))
		require.Empty(t, openAILegacySessionHashFromContext(c.Request.Context()))
placeholder)

	t.Run("prompt_cache_key is explicit", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", nil)

		got := svc.GenerateExplicitSessionHash(c, []byte(`{"model":"gpt-image-2","prompt_cache_key":"image-session"placeholder`))
		require.Equal(t, fmt.Sprintf("%016x", xxhash.Sum64String("image-session")), got)
		require.NotEmpty(t, openAILegacySessionHashFromContext(c.Request.Context()))
placeholder)

	t.Run("header overrides body", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", nil)
		c.Request.Header.Set("session_id", "header-session")

		got := svc.GenerateExplicitSessionHash(c, []byte(`{"prompt_cache_key":"body-session"placeholder`))
		require.Equal(t, fmt.Sprintf("%016x", xxhash.Sum64String("header-session")), got)
placeholder)
placeholder

func TestOpenAIGatewayService_GenerateSessionHashWithFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)

	svc := &OpenAIGatewayService{placeholder
	seed := "openai_ws_ingress:9:100:200"

	got := svc.GenerateSessionHashWithFallback(c, []byte(`{placeholder`), seed)
	want := fmt.Sprintf("%016x", xxhash.Sum64String(seed))
	require.Equal(t, want, got)
	require.NotEmpty(t, openAILegacySessionHashFromContext(c.Request.Context()))

	empty := svc.GenerateSessionHashWithFallback(c, []byte(`{placeholder`), "   ")
	require.Equal(t, "", empty)
placeholder

func TestOpenAIGatewayService_GenerateSessionHash_ContentFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/chat/completions", nil)

	svc := &OpenAIGatewayService{placeholder

	body := []byte(`{"model":"gpt-5.4","messages":[{"role":"system","content":"You are helpful."placeholder,{"role":"user","content":"Hello"placeholder]placeholder`)

	hash := svc.GenerateSessionHash(c, body)
	require.NotEmpty(t, hash, "content-based fallback should produce a hash")

	hash2 := svc.GenerateSessionHash(c, body)
	require.Equal(t, hash, hash2, "same content should produce same hash")

	bodyExtended := []byte(`{"model":"gpt-5.4","messages":[{"role":"system","content":"You are helpful."placeholder,{"role":"user","content":"Hello"placeholder,{"role":"assistant","content":"Hi!"placeholder,{"role":"user","content":"How are you?"placeholder]placeholder`)
	hashExtended := svc.GenerateSessionHash(c, bodyExtended)
	require.Equal(t, hash, hashExtended, "hash should be stable across later turns")

	bodyDifferent := []byte(`{"model":"gpt-5.4","messages":[{"role":"user","content":"Different question"placeholder]placeholder`)
	hashDifferent := svc.GenerateSessionHash(c, bodyDifferent)
	require.NotEqual(t, hash, hashDifferent, "different content should produce different hash")
placeholder

func TestOpenAIGatewayService_GenerateSessionHash_ExplicitSignalWinsOverContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/chat/completions", nil)

	svc := &OpenAIGatewayService{placeholder
	body := []byte(`{"model":"gpt-5.4","messages":[{"role":"user","content":"Hello"placeholder]placeholder`)

	contentHash := svc.GenerateSessionHash(c, body)
	require.NotEmpty(t, contentHash)

	c.Request.Header.Set("session_id", "explicit-session")
	explicitHash := svc.GenerateSessionHash(c, body)
	require.NotEmpty(t, explicitHash)
	require.NotEqual(t, contentHash, explicitHash, "explicit session_id should override content fallback")
placeholder

func TestOpenAIGatewayService_GenerateSessionHash_EmptyBodyStillEmpty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/chat/completions", nil)

	svc := &OpenAIGatewayService{placeholder
	require.Empty(t, svc.GenerateSessionHash(c, []byte(`{placeholder`)))
	require.Empty(t, svc.GenerateSessionHash(c, nil))
placeholder

func (c stubConcurrencyCache) GetAccountWaitingCount(ctx context.Context, accountID int64) (int, error) {
	if c.waitCounts != nil {
		if count, ok := c.waitCounts[accountID]; ok {
			return count, nil
	placeholder
placeholder
	return 0, nil
placeholder

type stubGatewayCache struct {
	sessionBindings map[string]int64
	deletedSessions map[string]int
placeholder

func (c *stubGatewayCache) GetSessionAccountID(ctx context.Context, groupID int64, sessionHash string) (int64, error) {
	if id, ok := c.sessionBindings[sessionHash]; ok {
		return id, nil
placeholder
	return 0, errors.New("not found")
placeholder

func (c *stubGatewayCache) SetSessionAccountID(ctx context.Context, groupID int64, sessionHash string, accountID int64, ttl time.Duration) error {
	if c.sessionBindings == nil {
		c.sessionBindings = make(map[string]int64)
placeholder
	c.sessionBindings[sessionHash] = accountID
	return nil
placeholder

func (c *stubGatewayCache) RefreshSessionTTL(ctx context.Context, groupID int64, sessionHash string, ttl time.Duration) error {
	return nil
placeholder

func (c *stubGatewayCache) DeleteSessionAccountID(ctx context.Context, groupID int64, sessionHash string) error {
	if c.sessionBindings == nil {
		return nil
placeholder
	if c.deletedSessions == nil {
		c.deletedSessions = make(map[string]int)
placeholder
	c.deletedSessions[sessionHash]++
	delete(c.sessionBindings, sessionHash)
	return nil
placeholder

func TestOpenAISelectAccountWithLoadAwareness_FiltersUnschedulable(t *testing.T) {
	now := time.Now()
	resetAt := now.Add(10 * time.Minute)
	groupID := int64(1)

	rateLimited := Account{
		ID:               1,
		Platform:         PlatformOpenAI,
		Type:             AccountTypeAPIKey,
		Status:           StatusActive,
		Schedulable:      true,
		Concurrency:      1,
		Priority:         0,
		RateLimitResetAt: &resetAt,
placeholder
	available := Account{
		ID:          2,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
		Priority:    1,
placeholder

	svc := &OpenAIGatewayService{
		accountRepo:        stubOpenAIAccountRepo{accounts: []Account{rateLimited, availableplaceholderplaceholder,
		concurrencyService: NewConcurrencyService(stubConcurrencyCache{placeholder),
placeholder

	selection, err := svc.SelectAccountWithLoadAwareness(context.Background(), &groupID, "", "gpt-5.2", nil)
	if err != nil {
		t.Fatalf("SelectAccountWithLoadAwareness error: %v", err)
placeholder
	if selection == nil || selection.Account == nil {
		t.Fatalf("expected selection with account")
placeholder
	if selection.Account.ID != available.ID {
		t.Fatalf("expected account %d, got %d", available.ID, selection.Account.ID)
placeholder
	if selection.ReleaseFunc != nil {
		selection.ReleaseFunc()
placeholder
placeholder

func TestOpenAISelectAccountWithLoadAwareness_ImageRateLimitSkipsOnlyImageRequests(t *testing.T) {
	future := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
	groupID := int64(1)

	imageLimited := Account{
		ID:          1,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
		Priority:    0,
		Extra: map[string]any{
			modelRateLimitsKey: map[string]any{
				openAIImageGenerationRateLimitKey: map[string]any{
					"rate_limit_reset_at": future,
			placeholder,
		placeholder,
	placeholder,
placeholder
	available := Account{
		ID:          2,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
		Priority:    1,
placeholder
	svc := &OpenAIGatewayService{
		accountRepo:        stubOpenAIAccountRepo{accounts: []Account{imageLimited, availableplaceholderplaceholder,
		concurrencyService: NewConcurrencyService(stubConcurrencyCache{placeholder),
placeholder

	imageSelection, err := svc.SelectAccountWithLoadAwareness(WithOpenAIImageGenerationIntent(context.Background()), &groupID, "", "gpt-5.4", nil)
placeholder
	require.NotNil(t, imageSelection)
	require.Equal(t, available.ID, imageSelection.Account.ID)
	if imageSelection.ReleaseFunc != nil {
		imageSelection.ReleaseFunc()
placeholder

	textSelection, err := svc.SelectAccountWithLoadAwareness(context.Background(), &groupID, "", "gpt-5.4", nil)
placeholder
	require.NotNil(t, textSelection)
	require.Equal(t, imageLimited.ID, textSelection.Account.ID)
	if textSelection.ReleaseFunc != nil {
		textSelection.ReleaseFunc()
placeholder
placeholder

func TestOpenAISelectAccountWithLoadAwareness_FiltersUnschedulableWhenNoConcurrencyService(t *testing.T) {
	now := time.Now()
	resetAt := now.Add(10 * time.Minute)
	groupID := int64(1)

	rateLimited := Account{
		ID:               1,
		Platform:         PlatformOpenAI,
		Type:             AccountTypeAPIKey,
		Status:           StatusActive,
		Schedulable:      true,
		Concurrency:      1,
		Priority:         0,
		RateLimitResetAt: &resetAt,
placeholder
	available := Account{
		ID:          2,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
		Priority:    1,
placeholder

	svc := &OpenAIGatewayService{
		accountRepo: stubOpenAIAccountRepo{accounts: []Account{rateLimited, availableplaceholderplaceholder,
		// concurrencyService is nil, forcing the non-load-batch selection path.
placeholder

	selection, err := svc.SelectAccountWithLoadAwareness(context.Background(), &groupID, "", "gpt-5.2", nil)
	if err != nil {
		t.Fatalf("SelectAccountWithLoadAwareness error: %v", err)
placeholder
	if selection == nil || selection.Account == nil {
		t.Fatalf("expected selection with account")
placeholder
	if selection.Account.ID != available.ID {
		t.Fatalf("expected account %d, got %d", available.ID, selection.Account.ID)
placeholder
	if selection.ReleaseFunc != nil {
		selection.ReleaseFunc()
placeholder
placeholder

func TestOpenAISelectAccountForModelWithExclusions_StickyUnschedulableClearsSession(t *testing.T) {
	sessionHash := "session-1"
	repo := stubOpenAIAccountRepo{
		accounts: []Account{
			{ID: 1, Platform: PlatformOpenAI, Status: StatusDisabled, Schedulable: true, Concurrency: 1placeholder,
			{ID: 2, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1placeholder,
	placeholder,
placeholder
	cache := &stubGatewayCache{
		sessionBindings: map[string]int64{"openai:" + sessionHash: 1placeholder,
placeholder

	svc := &OpenAIGatewayService{
		accountRepo: repo,
		cache:       cache,
placeholder

	acc, err := svc.SelectAccountForModelWithExclusions(context.Background(), nil, sessionHash, "gpt-4", nil)
	if err != nil {
		t.Fatalf("SelectAccountForModelWithExclusions error: %v", err)
placeholder
	if acc == nil || acc.ID != 2 {
		t.Fatalf("expected account 2, got %+v", acc)
placeholder
	if cache.deletedSessions["openai:"+sessionHash] != 1 {
		t.Fatalf("expected sticky session to be deleted")
placeholder
	if cache.sessionBindings["openai:"+sessionHash] != 2 {
		t.Fatalf("expected sticky session to bind to account 2")
placeholder
placeholder

func TestOpenAISelectAccountForModelWithExclusions_StickyOutsideGroupClearsSession(t *testing.T) {
	sessionHash := "session-outside-group"
	groupID := int64(1001)
	repo := groupAwareStubOpenAIAccountRepo{
		stubOpenAIAccountRepo{
			accounts: []Account{
				{ID: 1, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1placeholder,
				{ID: 2, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1, AccountGroups: []AccountGroup{{GroupID: groupIDplaceholderplaceholderplaceholder,
		placeholder,
	placeholder,
placeholder
	cache := &stubGatewayCache{
		sessionBindings: map[string]int64{"openai:" + sessionHash: 1placeholder,
placeholder

	svc := &OpenAIGatewayService{
		accountRepo: repo,
		cache:       cache,
placeholder

	acc, err := svc.SelectAccountForModelWithExclusions(context.Background(), &groupID, sessionHash, "gpt-4", nil)
	if err != nil {
		t.Fatalf("SelectAccountForModelWithExclusions error: %v", err)
placeholder
	if acc == nil || acc.ID != 2 {
		t.Fatalf("expected account 2, got %+v", acc)
placeholder
	if cache.deletedSessions["openai:"+sessionHash] != 1 {
		t.Fatalf("expected sticky session to be deleted")
placeholder
	if cache.sessionBindings["openai:"+sessionHash] != 2 {
		t.Fatalf("expected sticky session to bind to account 2")
placeholder
placeholder

func TestOpenAISelectAccountWithLoadAwareness_StickyUnschedulableClearsSession(t *testing.T) {
	sessionHash := "session-2"
	groupID := int64(1)
	repo := stubOpenAIAccountRepo{
		accounts: []Account{
			{ID: 1, Platform: PlatformOpenAI, Status: StatusDisabled, Schedulable: true, Concurrency: 1placeholder,
			{ID: 2, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1placeholder,
	placeholder,
placeholder
	cache := &stubGatewayCache{
		sessionBindings: map[string]int64{"openai:" + sessionHash: 1placeholder,
placeholder

	svc := &OpenAIGatewayService{
		accountRepo:        repo,
		cache:              cache,
		concurrencyService: NewConcurrencyService(stubConcurrencyCache{placeholder),
placeholder

	selection, err := svc.SelectAccountWithLoadAwareness(context.Background(), &groupID, sessionHash, "gpt-4", nil)
	if err != nil {
		t.Fatalf("SelectAccountWithLoadAwareness error: %v", err)
placeholder
	if selection == nil || selection.Account == nil || selection.Account.ID != 2 {
		t.Fatalf("expected account 2, got %+v", selection)
placeholder
	if cache.deletedSessions["openai:"+sessionHash] != 1 {
		t.Fatalf("expected sticky session to be deleted")
placeholder
	if cache.sessionBindings["openai:"+sessionHash] != 2 {
		t.Fatalf("expected sticky session to bind to account 2")
placeholder
	if selection.ReleaseFunc != nil {
		selection.ReleaseFunc()
placeholder
placeholder

func TestOpenAISelectAccountWithLoadAwareness_StickyOutsideGroupClearsSession(t *testing.T) {
	sessionHash := "session-load-outside-group"
	groupID := int64(1002)
	repo := groupAwareStubOpenAIAccountRepo{
		stubOpenAIAccountRepo{
			accounts: []Account{
				{ID: 1, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1placeholder,
				{ID: 2, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1, AccountGroups: []AccountGroup{{GroupID: groupIDplaceholderplaceholderplaceholder,
		placeholder,
	placeholder,
placeholder
	cache := &stubGatewayCache{
		sessionBindings: map[string]int64{"openai:" + sessionHash: 1placeholder,
placeholder

	svc := &OpenAIGatewayService{
		accountRepo:        repo,
		cache:              cache,
		concurrencyService: NewConcurrencyService(stubConcurrencyCache{placeholder),
placeholder

	selection, err := svc.SelectAccountWithLoadAwareness(context.Background(), &groupID, sessionHash, "gpt-4", nil)
	if err != nil {
		t.Fatalf("SelectAccountWithLoadAwareness error: %v", err)
placeholder
	if selection == nil || selection.Account == nil || selection.Account.ID != 2 {
		t.Fatalf("expected account 2, got %+v", selection)
placeholder
	if cache.deletedSessions["openai:"+sessionHash] != 1 {
		t.Fatalf("expected sticky session to be deleted")
placeholder
	if cache.sessionBindings["openai:"+sessionHash] != 2 {
		t.Fatalf("expected sticky session to bind to account 2")
placeholder
	if selection.ReleaseFunc != nil {
		selection.ReleaseFunc()
placeholder
placeholder

func TestOpenAISelectAccountForModelWithExclusions_NoModelSupport(t *testing.T) {
	repo := stubOpenAIAccountRepo{
		accounts: []Account{
			{
				ID:          1,
				Platform:    PlatformOpenAI,
				Status:      StatusActive,
				Schedulable: true,
		placeholder"model_mapping": map[string]any{"gpt-3.5-turbo": "gpt-3.5-turbo"placeholderplaceholder,
		placeholder,
	placeholder,
placeholder
	cache := &stubGatewayCache{placeholder

	svc := &OpenAIGatewayService{
		accountRepo: repo,
		cache:       cache,
placeholder

	acc, err := svc.SelectAccountForModelWithExclusions(context.Background(), nil, "", "gpt-4", nil)
	if err == nil {
		t.Fatalf("expected error for unsupported model")
placeholder
	if acc != nil {
		t.Fatalf("expected nil account for unsupported model")
placeholder
	if !strings.Contains(err.Error(), "supporting model") {
		t.Fatalf("unexpected error: %v", err)
placeholder
placeholder

func TestOpenAISelectAccountWithLoadAwareness_LoadBatchErrorFallback(t *testing.T) {
	groupID := int64(1)
	repo := stubOpenAIAccountRepo{
		accounts: []Account{
			{ID: 1, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 2placeholder,
			{ID: 2, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 1placeholder,
	placeholder,
placeholder
	cache := &stubGatewayCache{placeholder
	concurrencyCache := stubConcurrencyCache{
		loadBatchErr: errors.New("load batch failed"),
placeholder

	svc := &OpenAIGatewayService{
		accountRepo:        repo,
		cache:              cache,
		concurrencyService: NewConcurrencyService(concurrencyCache),
placeholder

	selection, err := svc.SelectAccountWithLoadAwareness(context.Background(), &groupID, "fallback", "gpt-4", nil)
	if err != nil {
		t.Fatalf("SelectAccountWithLoadAwareness error: %v", err)
placeholder
	if selection == nil || selection.Account == nil {
		t.Fatalf("expected selection")
placeholder
	if selection.Account.ID != 2 {
		t.Fatalf("expected account 2, got %d", selection.Account.ID)
placeholder
	if cache.sessionBindings["openai:fallback"] != 2 {
		t.Fatalf("expected sticky session updated")
placeholder
	if selection.ReleaseFunc != nil {
		selection.ReleaseFunc()
placeholder
placeholder

func TestOpenAISelectAccountWithLoadAwareness_NoSlotFallbackWait(t *testing.T) {
	groupID := int64(1)
	repo := stubOpenAIAccountRepo{
		accounts: []Account{
			{ID: 1, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 1placeholder,
	placeholder,
placeholder
	cache := &stubGatewayCache{placeholder
	concurrencyCache := stubConcurrencyCache{
		acquireResults: map[int64]bool{1: falseplaceholder,
		loadMap: map[int64]*AccountLoadInfo{
			1: {AccountID: 1, LoadRate: 10placeholder,
	placeholder,
placeholder

	svc := &OpenAIGatewayService{
		accountRepo:        repo,
		cache:              cache,
		concurrencyService: NewConcurrencyService(concurrencyCache),
placeholder

	selection, err := svc.SelectAccountWithLoadAwareness(context.Background(), &groupID, "", "gpt-4", nil)
	if err != nil {
		t.Fatalf("SelectAccountWithLoadAwareness error: %v", err)
placeholder
	if selection == nil || selection.WaitPlan == nil {
		t.Fatalf("expected wait plan fallback")
placeholder
	if selection.Account == nil || selection.Account.ID != 1 {
		t.Fatalf("expected account 1")
placeholder
placeholder

func TestOpenAISelectAccountForModelWithExclusions_SetsStickyBinding(t *testing.T) {
	sessionHash := "bind"
	repo := stubOpenAIAccountRepo{
		accounts: []Account{
			{ID: 1, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 1placeholder,
	placeholder,
placeholder
	cache := &stubGatewayCache{placeholder

	svc := &OpenAIGatewayService{
		accountRepo: repo,
		cache:       cache,
placeholder

	acc, err := svc.SelectAccountForModelWithExclusions(context.Background(), nil, sessionHash, "gpt-4", nil)
	if err != nil {
		t.Fatalf("SelectAccountForModelWithExclusions error: %v", err)
placeholder
	if acc == nil || acc.ID != 1 {
		t.Fatalf("expected account 1")
placeholder
	if cache.sessionBindings["openai:"+sessionHash] != 1 {
		t.Fatalf("expected sticky session binding")
placeholder
placeholder

func TestOpenAISelectAccountWithLoadAwareness_StickyWaitPlan(t *testing.T) {
	sessionHash := "sticky-wait"
	groupID := int64(1)
	repo := stubOpenAIAccountRepo{
		accounts: []Account{
			{ID: 1, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 1placeholder,
	placeholder,
placeholder
	cache := &stubGatewayCache{
		sessionBindings: map[string]int64{"openai:" + sessionHash: 1placeholder,
placeholder
	concurrencyCache := stubConcurrencyCache{
		acquireResults: map[int64]bool{1: falseplaceholder,
		waitCounts:     map[int64]int{1: 0placeholder,
placeholder

	svc := &OpenAIGatewayService{
		accountRepo:        repo,
		cache:              cache,
		concurrencyService: NewConcurrencyService(concurrencyCache),
placeholder

	selection, err := svc.SelectAccountWithLoadAwareness(context.Background(), &groupID, sessionHash, "gpt-4", nil)
	if err != nil {
		t.Fatalf("SelectAccountWithLoadAwareness error: %v", err)
placeholder
	if selection == nil || selection.WaitPlan == nil {
		t.Fatalf("expected sticky wait plan")
placeholder
	if selection.Account == nil || selection.Account.ID != 1 {
		t.Fatalf("expected account 1")
placeholder
placeholder

func TestOpenAISelectAccountWithLoadAwareness_PrefersLowerLoad(t *testing.T) {
	groupID := int64(1)
	repo := stubOpenAIAccountRepo{
		accounts: []Account{
			{ID: 1, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 1placeholder,
			{ID: 2, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 1placeholder,
	placeholder,
placeholder
	cache := &stubGatewayCache{placeholder
	concurrencyCache := stubConcurrencyCache{
		loadMap: map[int64]*AccountLoadInfo{
			1: {AccountID: 1, LoadRate: 80placeholder,
			2: {AccountID: 2, LoadRate: 10placeholder,
	placeholder,
placeholder

	svc := &OpenAIGatewayService{
		accountRepo:        repo,
		cache:              cache,
		concurrencyService: NewConcurrencyService(concurrencyCache),
placeholder

	selection, err := svc.SelectAccountWithLoadAwareness(context.Background(), &groupID, "load", "gpt-4", nil)
	if err != nil {
		t.Fatalf("SelectAccountWithLoadAwareness error: %v", err)
placeholder
	if selection == nil || selection.Account == nil || selection.Account.ID != 2 {
		t.Fatalf("expected account 2")
placeholder
	if cache.sessionBindings["openai:load"] != 2 {
		t.Fatalf("expected sticky session updated")
placeholder
placeholder

func TestOpenAISelectAccountForModelWithExclusions_StickyExcludedFallback(t *testing.T) {
	sessionHash := "excluded"
	repo := stubOpenAIAccountRepo{
		accounts: []Account{
			{ID: 1, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 1placeholder,
			{ID: 2, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 2placeholder,
	placeholder,
placeholder
	cache := &stubGatewayCache{
		sessionBindings: map[string]int64{"openai:" + sessionHash: 1placeholder,
placeholder

	svc := &OpenAIGatewayService{
		accountRepo: repo,
		cache:       cache,
placeholder

	excluded := map[int64]struct{placeholder{1: {placeholderplaceholder
	acc, err := svc.SelectAccountForModelWithExclusions(context.Background(), nil, sessionHash, "gpt-4", excluded)
	if err != nil {
		t.Fatalf("SelectAccountForModelWithExclusions error: %v", err)
placeholder
	if acc == nil || acc.ID != 2 {
		t.Fatalf("expected account 2")
placeholder
placeholder

func TestOpenAISelectAccountForModelWithExclusions_StickyNonOpenAI(t *testing.T) {
	sessionHash := "non-openai"
	repo := stubOpenAIAccountRepo{
		accounts: []Account{
			{ID: 1, Platform: PlatformAnthropic, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 1placeholder,
			{ID: 2, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 2placeholder,
	placeholder,
placeholder
	cache := &stubGatewayCache{
		sessionBindings: map[string]int64{"openai:" + sessionHash: 1placeholder,
placeholder

	svc := &OpenAIGatewayService{
		accountRepo: repo,
		cache:       cache,
placeholder

	acc, err := svc.SelectAccountForModelWithExclusions(context.Background(), nil, sessionHash, "gpt-4", nil)
	if err != nil {
		t.Fatalf("SelectAccountForModelWithExclusions error: %v", err)
placeholder
	if acc == nil || acc.ID != 2 {
		t.Fatalf("expected account 2")
placeholder
placeholder

func TestOpenAISelectAccountForModelWithExclusions_NoAccounts(t *testing.T) {
	repo := stubOpenAIAccountRepo{accounts: []Account{placeholderplaceholder
	cache := &stubGatewayCache{placeholder

	svc := &OpenAIGatewayService{
		accountRepo: repo,
		cache:       cache,
placeholder

	acc, err := svc.SelectAccountForModelWithExclusions(context.Background(), nil, "", "", nil)
	if err == nil {
		t.Fatalf("expected error for no accounts")
placeholder
	if acc != nil {
		t.Fatalf("expected nil account")
placeholder
	if !strings.Contains(err.Error(), "no available OpenAI accounts") {
		t.Fatalf("unexpected error: %v", err)
placeholder
placeholder

func TestOpenAISelectAccountWithLoadAwareness_NoCandidates(t *testing.T) {
	groupID := int64(1)
	resetAt := time.Now().Add(1 * time.Hour)
	repo := stubOpenAIAccountRepo{
		accounts: []Account{
			{ID: 1, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 1, RateLimitResetAt: &resetAtplaceholder,
	placeholder,
placeholder
	cache := &stubGatewayCache{placeholder
	concurrencyCache := stubConcurrencyCache{placeholder

	svc := &OpenAIGatewayService{
		accountRepo:        repo,
		cache:              cache,
		concurrencyService: NewConcurrencyService(concurrencyCache),
placeholder

	selection, err := svc.SelectAccountWithLoadAwareness(context.Background(), &groupID, "", "gpt-4", nil)
	if err == nil {
		t.Fatalf("expected error for no candidates")
placeholder
	if selection != nil {
		t.Fatalf("expected nil selection")
placeholder
placeholder

func TestOpenAISelectAccountWithLoadAwareness_AllFullWaitPlan(t *testing.T) {
	groupID := int64(1)
	repo := stubOpenAIAccountRepo{
		accounts: []Account{
			{ID: 1, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 1placeholder,
	placeholder,
placeholder
	cache := &stubGatewayCache{placeholder
	concurrencyCache := stubConcurrencyCache{
		loadMap: map[int64]*AccountLoadInfo{
			1: {AccountID: 1, LoadRate: 100placeholder,
	placeholder,
placeholder

	svc := &OpenAIGatewayService{
		accountRepo:        repo,
		cache:              cache,
		concurrencyService: NewConcurrencyService(concurrencyCache),
placeholder

	selection, err := svc.SelectAccountWithLoadAwareness(context.Background(), &groupID, "", "gpt-4", nil)
	if err != nil {
		t.Fatalf("SelectAccountWithLoadAwareness error: %v", err)
placeholder
	if selection == nil || selection.WaitPlan == nil {
		t.Fatalf("expected wait plan")
placeholder
placeholder

func TestOpenAISelectAccountWithLoadAwareness_LoadBatchErrorNoAcquire(t *testing.T) {
	groupID := int64(1)
	repo := stubOpenAIAccountRepo{
		accounts: []Account{
			{ID: 1, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 1placeholder,
	placeholder,
placeholder
	cache := &stubGatewayCache{placeholder
	concurrencyCache := stubConcurrencyCache{
		loadBatchErr:   errors.New("load batch failed"),
		acquireResults: map[int64]bool{1: falseplaceholder,
placeholder

	svc := &OpenAIGatewayService{
		accountRepo:        repo,
		cache:              cache,
		concurrencyService: NewConcurrencyService(concurrencyCache),
placeholder

	selection, err := svc.SelectAccountWithLoadAwareness(context.Background(), &groupID, "", "gpt-4", nil)
	if err != nil {
		t.Fatalf("SelectAccountWithLoadAwareness error: %v", err)
placeholder
	if selection == nil || selection.WaitPlan == nil {
		t.Fatalf("expected wait plan")
placeholder
placeholder

func TestOpenAISelectAccountWithLoadAwareness_MissingLoadInfo(t *testing.T) {
	groupID := int64(1)
	repo := stubOpenAIAccountRepo{
		accounts: []Account{
			{ID: 1, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 1placeholder,
			{ID: 2, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 1placeholder,
	placeholder,
placeholder
	cache := &stubGatewayCache{placeholder
	concurrencyCache := stubConcurrencyCache{
		loadMap: map[int64]*AccountLoadInfo{
			1: {AccountID: 1, LoadRate: 50placeholder,
	placeholder,
		skipDefaultLoad: true,
placeholder

	svc := &OpenAIGatewayService{
		accountRepo:        repo,
		cache:              cache,
		concurrencyService: NewConcurrencyService(concurrencyCache),
placeholder

	selection, err := svc.SelectAccountWithLoadAwareness(context.Background(), &groupID, "", "gpt-4", nil)
	if err != nil {
		t.Fatalf("SelectAccountWithLoadAwareness error: %v", err)
placeholder
	if selection == nil || selection.Account == nil || selection.Account.ID != 2 {
		t.Fatalf("expected account 2")
placeholder
placeholder

func TestOpenAISelectAccountForModelWithExclusions_LeastRecentlyUsed(t *testing.T) {
	oldTime := time.Now().Add(-2 * time.Hour)
	newTime := time.Now().Add(-1 * time.Hour)
	repo := stubOpenAIAccountRepo{
		accounts: []Account{
			{ID: 1, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Priority: 1, LastUsedAt: &newTimeplaceholder,
			{ID: 2, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Priority: 1, LastUsedAt: &oldTimeplaceholder,
	placeholder,
placeholder
	cache := &stubGatewayCache{placeholder

	svc := &OpenAIGatewayService{
		accountRepo: repo,
		cache:       cache,
placeholder

	acc, err := svc.SelectAccountForModelWithExclusions(context.Background(), nil, "", "gpt-4", nil)
	if err != nil {
		t.Fatalf("SelectAccountForModelWithExclusions error: %v", err)
placeholder
	if acc == nil || acc.ID != 2 {
		t.Fatalf("expected account 2")
placeholder
placeholder

func TestOpenAISelectAccountWithLoadAwareness_PreferNeverUsed(t *testing.T) {
	groupID := int64(1)
	lastUsed := time.Now().Add(-1 * time.Hour)
	repo := stubOpenAIAccountRepo{
		accounts: []Account{
			{ID: 1, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 1, LastUsedAt: &lastUsedplaceholder,
			{ID: 2, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 1placeholder,
	placeholder,
placeholder
	cache := &stubGatewayCache{placeholder
	concurrencyCache := stubConcurrencyCache{
		loadMap: map[int64]*AccountLoadInfo{
			1: {AccountID: 1, LoadRate: 10placeholder,
			2: {AccountID: 2, LoadRate: 10placeholder,
	placeholder,
placeholder

	svc := &OpenAIGatewayService{
		accountRepo:        repo,
		cache:              cache,
		concurrencyService: NewConcurrencyService(concurrencyCache),
placeholder

	selection, err := svc.SelectAccountWithLoadAwareness(context.Background(), &groupID, "", "gpt-4", nil)
	if err != nil {
		t.Fatalf("SelectAccountWithLoadAwareness error: %v", err)
placeholder
	if selection == nil || selection.Account == nil || selection.Account.ID != 2 {
		t.Fatalf("expected account 2")
placeholder
placeholder

func TestOpenAIStreamingTimeout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 1,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       pr,
		Header:     http.Header{placeholder,
placeholder

	start := time.Now()
	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1placeholder, start, "model", "model")
	_ = pw.Close()
	_ = pr.Close()

	if err == nil || !strings.Contains(err.Error(), "stream data interval timeout") {
		t.Fatalf("expected stream timeout error, got %v", err)
placeholder
	if !strings.Contains(rec.Body.String(), "\"type\":\"error\"") || !strings.Contains(rec.Body.String(), "stream_timeout") {
		t.Fatalf("expected OpenAI-compatible error SSE event, got %q", rec.Body.String())
placeholder
placeholder

func TestOpenAIStreamingContextCanceledReturnsIncompleteErrorWithoutInjectingErrorEvent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil).WithContext(ctx)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       cancelReadCloser{placeholder,
		Header:     http.Header{placeholder,
placeholder

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1placeholder, time.Now(), "model", "model")
	if err == nil || !strings.Contains(err.Error(), "stream usage incomplete") {
		t.Fatalf("expected incomplete stream error, got %v", err)
placeholder
	if strings.Contains(rec.Body.String(), "event: error") || strings.Contains(rec.Body.String(), "stream_read_error") {
		t.Fatalf("expected no injected SSE error event, got %q", rec.Body.String())
placeholder
placeholder

func TestOpenAIStreamingReadErrorBeforeOutputReturnsFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       errReadCloser{err: io.ErrUnexpectedEOFplaceholder,
		Header:     http.Header{"X-Request-Id": []string{"rid-disconnect"placeholderplaceholder,
placeholder

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Name: "acc"placeholder, time.Now(), "model", "model")
placeholder
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusBadGateway, failoverErr.StatusCode)
	require.False(t, c.Writer.Written())
	require.Empty(t, rec.Body.String())
placeholder

func TestOpenAIStreamingPostOutputDisconnectQuarantinesSharedProxyWithoutSameStreamFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)
	proxyID := int64(4698)
	account := &Account{
		ID:       469801,
		Name:     "oauth-on-shared-proxy",
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
		ProxyID:  &proxyID,
placeholder
	svc := &OpenAIGatewayService{cfg: &config.Config{Gateway: config.GatewayConfig{
		MaxLineSize: defaultMaxLineSize,
placeholderplaceholderplaceholder
	// collapseInterval 0: the two loop iterations below record within the
	// production collapse window and must count as distinct failure events here.
	svc.openaiProxyStreamCircuit = newOpenAIProxyStreamCircuit(openAIProxyStreamCircuitSettings{
		failureThreshold: 2,
		failureWindow:    time.Minute,
		quarantineTTL:    10 * time.Minute,
		maxEntries:       16,
placeholder)

	for _, readErr := range []error{
		io.ErrUnexpectedEOF,
		errors.New("http2: client connection lost"),
placeholder {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body: &openAIStreamReadThenErrorCloser{
				reader: strings.NewReader(strings.Join([]string{
					"event: response.output_text.delta",
					`data: {"type":"response.output_text.delta","delta":"partial"placeholder`,
					"",
			placeholder, "\n")),
				err: readErr,
		placeholder,
			Header: http.Header{"X-Request-Id": []string{"rid-proxy-disconnect"placeholderplaceholder,
	placeholder

		_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, account, time.Now(), "gpt-5.6-sol", "gpt-5.6-sol")
	placeholder
		var failoverErr *UpstreamFailoverError
		require.False(t, errors.As(err, &failoverErr), "post-output disconnect must not fail over inside the same stream")
		require.Contains(t, rec.Body.String(), "partial")
placeholder

	scheduler := &defaultOpenAIAccountScheduler{service: svcplaceholder
	compatible, reason := scheduler.isAccountRequestCompatibleReason(context.Background(), account, OpenAIAccountScheduleRequest{placeholder)
	require.False(t, compatible, "the next request must exclude accounts sharing the quarantined proxy")
	require.Equal(t, "proxy_stream_quarantined", reason)
placeholder

func TestOpenAIStreamingTerminalAndClientCancellationDoNotQuarantineProxy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	proxyID := int64(4699)
	account := &Account{ID: 469901, Name: "oauth", Platform: PlatformOpenAI, Type: AccountTypeOAuth, ProxyID: &proxyIDplaceholder
	svc := &OpenAIGatewayService{cfg: &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholderplaceholder

	terminalRecorder := httptest.NewRecorder()
	terminalCtx, _ := gin.CreateTestContext(terminalRecorder)
	terminalCtx.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	terminalResp := &http.Response{
		StatusCode: http.StatusOK,
		Body: &openAIStreamReadThenErrorCloser{
			reader: strings.NewReader(strings.Join([]string{
				"event: response.completed",
				`data: {"type":"response.completed","response":{"status":"completed","output":[]placeholderplaceholder`,
				"",
		placeholder, "\n")),
			err: io.ErrUnexpectedEOF,
	placeholder,
		Header: http.Header{placeholder,
placeholder
	_, err := svc.handleStreamingResponse(terminalCtx.Request.Context(), terminalResp, terminalCtx, account, time.Now(), "model", "model")
placeholder

	for range 2 {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil).WithContext(ctx)
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body: &openAIStreamReadThenErrorCloser{
				reader: strings.NewReader("data: {\"type\":\"response.output_text.delta\",\"delta\":\"partial\"placeholder\n\n"),
				err:    context.Canceled,
		placeholder,
			Header: http.Header{placeholder,
	placeholder
		_, err = svc.handleStreamingResponse(c.Request.Context(), resp, c, account, time.Now(), "model", "model")
	placeholder
placeholder

	scheduler := &defaultOpenAIAccountScheduler{service: svcplaceholder
	compatible, reason := scheduler.isAccountRequestCompatibleReason(context.Background(), account, OpenAIAccountScheduleRequest{placeholder)
	require.True(t, compatible)
	require.Empty(t, reason)
placeholder

func TestOpenAIStreamingResponseFailedBeforeOutputReturnsFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
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
			"event: response.in_progress",
			`data: {"type":"response.in_progress","response":{"id":"resp_1"placeholderplaceholder`,
			"",
			"event: response.failed",
			`data: {"type":"response.failed","error":{"message":"An error occurred while processing your request."placeholderplaceholder`,
			"",
	placeholder, "\n"))),
		Header: http.Header{"X-Request-Id": []string{"rid-failed"placeholderplaceholder,
placeholder

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Name: "acc"placeholder, time.Now(), "model", "model")
placeholder
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusBadGateway, failoverErr.StatusCode)
	require.False(t, failoverErr.RetryableOnSameAccount)
	require.Contains(t, string(failoverErr.ResponseBody), "An error occurred while processing your request")
	require.False(t, c.Writer.Written())
	require.Empty(t, rec.Body.String())
placeholder

func TestOpenAIStreamingResponseFailedBeforeOutputCapacityErrorReturnsFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
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
			"event: response.in_progress",
			`data: {"type":"response.in_progress","response":{"id":"resp_1"placeholderplaceholder`,
			"",
			"event: response.failed",
			`data: {"type":"response.failed","error":{"message":"Selected model is at capacity. Please try a different model.","type":"invalid_request_error"placeholderplaceholder`,
			"",
	placeholder, "\n"))),
		Header: http.Header{"X-Request-Id": []string{"rid-capacity-failed"placeholderplaceholder,
placeholder

	account := &Account{
		ID:       1,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Name:     "pool-account",
placeholder
			"pool_mode": true,
	placeholder,
placeholder
	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, account, time.Now(), "model", "model")
placeholder
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusBadGateway, failoverErr.StatusCode)
	require.True(t, failoverErr.RetryableOnSameAccount)
	require.Contains(t, string(failoverErr.ResponseBody), "Selected model is at capacity")
	require.False(t, c.Writer.Written())
	require.Empty(t, rec.Body.String())
placeholder

func TestOpenAIStreamingResponseFailedBeforeOutputServerOverloadedCodeReturnsFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
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
			"event: response.failed",
			`data: {"type":"response.failed","response":{"id":"resp_1","error":{"code":"server_is_overloaded","message":"Please retry later."placeholderplaceholderplaceholder`,
			"",
	placeholder, "\n"))),
		Header: http.Header{"X-Request-Id": []string{"rid-overloaded-failed"placeholderplaceholder,
placeholder

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Name: "acc"placeholder, time.Now(), "model", "model")
placeholder
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusBadGateway, failoverErr.StatusCode)
	require.Contains(t, string(failoverErr.ResponseBody), "Please retry later")
	// 容量降载是请求级信号：非池模式账号也要先在同账号重试，且不得据此临时封禁账号。
	// 否则单个被降载的请求会把整池账号逐个消耗掉，而降载因素在每个账号上都相同。
	require.True(t, failoverErr.RetryableOnSameAccount)
	require.True(t, failoverErr.RequestScopedTransient)
	require.False(t, c.Writer.Written())
	require.Empty(t, rec.Body.String())
placeholder

func TestOpenAIStreamingResponseFailedBeforeOutputRateLimitUsesPoolRetryPolicy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
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
			"event: response.failed",
			`data: {"type":"response.failed","response":{"id":"resp_1","status":"failed","error":{"code":"rate_limit_exceeded","message":"Concurrency limit exceeded for account, please retry later"placeholderplaceholderplaceholder`,
			"",
	placeholder, "\n"))),
		Header: http.Header{
			"X-Request-Id": []string{"rid-rate-limit-failed"placeholder,
			"Retry-After":  []string{"1"placeholder,
	placeholder,
placeholder
	account := &Account{
		ID:       1,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Name:     "pool-account",
placeholder
			"pool_mode":                    true,
			"pool_mode_retry_count":        float64(1),
			"pool_mode_retry_status_codes": []any{float64(http.StatusTooManyRequests)placeholder,
	placeholder,
placeholder

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, account, time.Now(), "model", "model")
placeholder
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusTooManyRequests, failoverErr.StatusCode)
	require.True(t, failoverErr.RetryableOnSameAccount)
	require.Equal(t, "1", failoverErr.ResponseHeaders.Get("Retry-After"))
	require.Equal(t, "rate_limit_error", gjson.GetBytes(failoverErr.ResponseBody, "error.type").String())
	require.Contains(t, string(failoverErr.ResponseBody), "Concurrency limit exceeded")
	require.False(t, c.Writer.Written())
	require.Empty(t, rec.Body.String())

	opsVal, ok := c.Get(OpsUpstreamErrorsKey)
	require.True(t, ok)
	opsEvents, ok := opsVal.([]*OpsUpstreamErrorEvent)
	require.True(t, ok)
	require.NotEmpty(t, opsEvents)
	require.Equal(t, http.StatusTooManyRequests, opsEvents[len(opsEvents)-1].UpstreamStatusCode)
placeholder

// 流内 rate limit 只产生 failover 错误，不写账号级限流/封禁状态：
// HTTP 200 流的 x-codex-* 头是正常配额快照，不能按 429 头驱动账号冷却。
func TestOpenAIStreamingResponseFailedRateLimitDoesNotBlockAccountScheduling(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
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
			"event: response.failed",
			`data: {"type":"response.failed","response":{"id":"resp_1","status":"failed","error":{"code":"rate_limit_exceeded","message":"Concurrency limit exceeded for account, please retry later"placeholderplaceholderplaceholder`,
			"",
	placeholder, "\n"))),
		Header: http.Header{
			"X-Codex-Primary-Used-Percent":        []string{"12"placeholder,
			"X-Codex-Primary-Reset-After-Seconds": []string{"604800"placeholder,
			"Retry-After":                         []string{"1"placeholder,
	placeholder,
placeholder
	account := &Account{
		ID:       11,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
		Name:     "oauth-account",
placeholder

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, account, time.Now(), "model", "model")
placeholder
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusTooManyRequests, failoverErr.StatusCode)
	require.False(t, failoverErr.RetryableOnSameAccount)
	require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestOpenAIStreamingResponseFailedAfterOutputSanitizesVerboseResponseForClient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	longInstructions := strings.Repeat("You are GPT-5.1 running in the Codex CLI. ", 20)
	failedPayload := fmt.Sprintf(
		`{"type":"response.failed","response":{"id":"resp_failed","object":"response","created_at":1782446336,"status":"failed","instructions":%q,"output":[{"type":"message","content":[{"type":"output_text","text":"large"placeholder]placeholder],"usage":{"input_tokens":123,"output_tokens":0placeholder,"error":{"code":"context_length_exceeded","message":"Your input exceeds the context window of this model. Please adjust your input and try again."placeholderplaceholderplaceholder`,
		longInstructions,
	)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(strings.NewReader(strings.Join([]string{
			"event: response.created",
			`data: {"type":"response.created","response":{"id":"resp_failed"placeholderplaceholder`,
			"",
			"event: response.output_text.delta",
			`data: {"type":"response.output_text.delta","delta":"partial"placeholder`,
			"",
			"event: response.failed",
			"data: " + failedPayload,
			"",
	placeholder, "\n"))),
		Header: http.Header{"X-Request-Id": []string{"rid-failed-after-output"placeholderplaceholder,
placeholder

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Name: "acc"placeholder, time.Now(), "model", "model")
placeholder

	body := rec.Body.String()
	require.Contains(t, body, "event: response.failed")
	require.Contains(t, body, "context_length_exceeded")
	require.Contains(t, body, `"type":"invalid_request_error"`)
	require.Contains(t, body, "Your input exceeds the context window")
	require.NotContains(t, body, "You are GPT-5.1 running in the Codex CLI")
	require.NotContains(t, body, `"instructions"`)
	require.NotContains(t, body, `"output"`)
	require.NotContains(t, body, `"usage"`)
placeholder

func TestOpenAIStreamingContextWindowResponseFailedBeforeOutputPassesThrough(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
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
			"event: response.failed",
			`data: {"type":"response.failed","error":{"type":"upstream_error","message":"Your input exceeds the context window of this model. Please adjust your input and try again.","code":nullplaceholderplaceholder`,
			"",
	placeholder, "\n"))),
		Header: http.Header{"X-Request-Id": []string{"rid-context-window-failed"placeholderplaceholder,
placeholder

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Name: "acc"placeholder, time.Now(), "model", "model")
placeholder
	var failoverErr *UpstreamFailoverError
	require.False(t, errors.As(err, &failoverErr))
	require.True(t, c.Writer.Written())
	require.Contains(t, rec.Body.String(), "response.failed")
	require.Contains(t, rec.Body.String(), `"type":"upstream_error"`)
	require.Contains(t, rec.Body.String(), "Your input exceeds the context window")
placeholder

func TestOpenAIStreamingContextWindowResponseFailedBeforeOutputAppliesPassthroughRule(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	rule := newNonFailoverPassthroughRule(http.StatusBadRequest, "context_length_exceeded", http.StatusBadRequest, "")
	rule.Platforms = []string{PlatformOpenAIplaceholder
	rule.PassthroughBody = true
	rule.CustomMessage = nil
	ruleSvc := &ErrorPassthroughService{placeholder
	ruleSvc.setLocalCache([]*model.ErrorPassthroughRule{ruleplaceholder)
	BindErrorPassthroughService(c, ruleSvc)

	upstreamMessage := "Your input exceeds the context window of this model. Please adjust your input and try again."
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(strings.NewReader(strings.Join([]string{
			"event: response.created",
			`data: {"type":"response.created","response":{"id":"resp_1"placeholderplaceholder`,
			"",
			"event: response.failed",
			`data: {"type":"response.failed","response":{"id":"resp_1","error":{"type":"invalid_request_error","code":"context_length_exceeded","message":"` + upstreamMessage + `"placeholderplaceholderplaceholder`,
			"",
	placeholder, "\n"))),
		Header: http.Header{"X-Request-Id": []string{"rid-context-window-passthrough-rule"placeholderplaceholder,
placeholder

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Name: "acc"placeholder, time.Now(), "model", "model")
placeholder
	var failoverErr *UpstreamFailoverError
	require.False(t, errors.As(err, &failoverErr))
	require.True(t, IsResponseCommitted(c))
	require.Equal(t, http.StatusBadRequest, rec.Code)
	body := rec.Body.String()
	require.Equal(t, "upstream_error", gjson.Get(body, "error.type").String())
	require.Equal(t, upstreamMessage, gjson.Get(body, "error.message").String())
	require.NotContains(t, body, "response.failed")
	require.NotContains(t, body, "Upstream request failed")
	// 命中透传规则也应记录 ops 上游错误事件（对齐 CC/Messages 与 antigravity 先例）。
	opsVal, opsRecorded := c.Get(OpsUpstreamErrorsKey)
	require.True(t, opsRecorded, "passthrough hit should record an ops upstream error event")
	opsEvents, _ := opsVal.([]*OpsUpstreamErrorEvent)
	require.NotEmpty(t, opsEvents)
placeholder

func TestOpenAIStreamingPreambleOnlyMissingTerminalReturnsFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
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
			"event: response.in_progress",
			`data: {"type":"response.in_progress","response":{"id":"resp_1"placeholderplaceholder`,
			"",
	placeholder, "\n"))),
		Header: http.Header{"X-Request-Id": []string{"rid-missing-terminal"placeholderplaceholder,
placeholder

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Name: "acc"placeholder, time.Now(), "model", "model")
placeholder
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.False(t, c.Writer.Written())
	require.Empty(t, rec.Body.String())
placeholder

func TestOpenAIStreamingPreambleKeepaliveUsesDownstreamIdle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   1,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       pr,
		Header:     http.Header{placeholder,
placeholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		_, _ = pw.Write([]byte("data: {\"type\":\"response.created\",\"response\":{\"id\":\"resp_1\"placeholderplaceholder\n\n"))
		for i := 0; i < 6; i++ {
			time.Sleep(250 * time.Millisecond)
			_, _ = pw.Write([]byte("data: {\"type\":\"response.in_progress\",\"response\":{\"id\":\"resp_1\"placeholderplaceholder\n\n"))
	placeholder
		_, _ = pw.Write([]byte("data: {\"type\":\"response.completed\",\"response\":{\"usage\":{\"input_tokens\":1,\"output_tokens\":2placeholderplaceholderplaceholder\n\n"))
placeholder()

	result, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Name: "acc"placeholder, time.Now(), "model", "model")
	_ = pr.Close()
placeholder
	require.NotNil(t, result)
	require.Contains(t, rec.Body.String(), ":\n\n")
	require.Contains(t, rec.Body.String(), "response.completed")
placeholder

func TestOpenAIStreamingNormalizesTerminalOutputFromDeltas(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(strings.NewReader(strings.Join([]string{
			`data: {"type":"response.created","response":{"id":"resp_sdk_parse"placeholderplaceholder`,
			"",
			`data: {"type":"response.output_text.delta","delta":"pon"placeholder`,
			"",
			`data: {"type":"response.output_text.delta","delta":"g"placeholder`,
			"",
			`data: {"type":"response.completed","response":{"id":"resp_sdk_parse","status":"completed","output":null,"usage":{"input_tokens":1,"output_tokens":1placeholderplaceholderplaceholder`,
			"",
	placeholder, "\n"))),
		Header: http.Header{"X-Request-Id": []string{"rid-sdk-parse"placeholderplaceholder,
placeholder

	result, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Name: "acc"placeholder, time.Now(), "model", "model")
placeholder
	require.NotNil(t, result)

	terminalType, terminalPayload, ok := extractOpenAISSETerminalEvent(rec.Body.String())
	require.True(t, ok)
	require.Equal(t, "response.completed", terminalType)
	output := gjson.GetBytes(terminalPayload, "response.output")
	require.True(t, output.IsArray())
	require.Len(t, output.Array(), 1)
	require.Equal(t, "pong", gjson.GetBytes(terminalPayload, "response.output.0.content.0.text").String())
placeholder

func TestOpenAIStreamingNormalizesTerminalOutputToEmptyArray(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(strings.NewReader(strings.Join([]string{
			`data: {"type":"response.completed","response":{"id":"resp_empty","status":"completed","output":null,"usage":{"input_tokens":1,"output_tokens":0placeholderplaceholderplaceholder`,
			"",
	placeholder, "\n"))),
		Header: http.Header{"X-Request-Id": []string{"rid-empty-output"placeholderplaceholder,
placeholder

	result, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Name: "acc"placeholder, time.Now(), "model", "model")
placeholder
	require.NotNil(t, result)

	terminalType, terminalPayload, ok := extractOpenAISSETerminalEvent(rec.Body.String())
	require.True(t, ok)
	require.Equal(t, "response.completed", terminalType)
	output := gjson.GetBytes(terminalPayload, "response.output")
	require.True(t, output.IsArray())
	require.Len(t, output.Array(), 0)
placeholder

func TestOpenAIStreamingPolicyResponseFailedBeforeOutputPassesThrough(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
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
			"event: response.failed",
			`data: {"type":"response.failed","error":{"type":"safety_error","message":"This request has been flagged for potentially high-risk cyber activity."placeholderplaceholder`,
			"",
	placeholder, "\n"))),
		Header: http.Header{"X-Request-Id": []string{"rid-policy-failed"placeholderplaceholder,
placeholder

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Name: "acc"placeholder, time.Now(), "model", "model")
placeholder
	var failoverErr *UpstreamFailoverError
	require.False(t, errors.As(err, &failoverErr))
	require.True(t, c.Writer.Written())
	require.Contains(t, rec.Body.String(), "response.failed")
	require.Contains(t, rec.Body.String(), "high-risk cyber activity")
placeholder

func TestOpenAIStreamingClientDisconnectDrainsUpstreamUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	c.Writer = &failingGinWriter{ResponseWriter: c.Writer, failAfter: 0placeholder

	pr, pw := io.Pipe()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       pr,
		Header:     http.Header{placeholder,
placeholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		_, _ = pw.Write([]byte("data: {\"type\":\"response.in_progress\",\"response\":{placeholderplaceholder\n\n"))
		_, _ = pw.Write([]byte("data: {\"type\":\"response.completed\",\"response\":{\"usage\":{\"input_tokens\":3,\"output_tokens\":5,\"input_tokens_details\":{\"cached_tokens\":1placeholderplaceholderplaceholderplaceholder\n\n"))
placeholder()

	result, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1placeholder, time.Now(), "model", "model")
	_ = pr.Close()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
placeholder
	if result == nil || result.usage == nil {
		t.Fatalf("expected usage result")
placeholder
	if result.usage.InputTokens != 3 || result.usage.OutputTokens != 5 || result.usage.CacheReadInputTokens != 1 {
		t.Fatalf("unexpected usage: %+v", *result.usage)
placeholder
	if strings.Contains(rec.Body.String(), "event: error") || strings.Contains(rec.Body.String(), "write_failed") {
		t.Fatalf("expected no injected SSE error event, got %q", rec.Body.String())
placeholder
placeholder

func TestOpenAIStreamingMissingTerminalEventReturnsIncompleteError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       pr,
		Header:     http.Header{placeholder,
placeholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		_, _ = pw.Write([]byte("data: {\"type\":\"response.output_item.added\",\"item\":{\"type\":\"message\"placeholder,\"output_index\":0placeholder\n\n"))
placeholder()

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1placeholder, time.Now(), "model", "model")
	_ = pr.Close()
	if err == nil || !strings.Contains(err.Error(), "missing terminal event") {
		t.Fatalf("expected missing terminal event error, got %v", err)
placeholder
placeholder

func TestOpenAIStreamingPassthroughMissingTerminalEventReturnsIncompleteError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			MaxLineSize: defaultMaxLineSize,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       pr,
		Header:     http.Header{placeholder,
placeholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		_, _ = pw.Write([]byte("data: {\"type\":\"response.output_item.added\",\"item\":{\"type\":\"message\"placeholder,\"output_index\":0placeholder\n\n"))
placeholder()

	_, err := svc.handleStreamingResponsePassthrough(c.Request.Context(), resp, c, &Account{ID: 1placeholder, time.Now(), "", "")
	_ = pr.Close()
	if err == nil || !strings.Contains(err.Error(), "missing terminal event") {
		t.Fatalf("expected missing terminal event error, got %v", err)
placeholder
placeholder

func TestOpenAIStreamingPassthroughPostOutputDisconnectQuarantinesSharedProxy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	proxyID := int64(4698)
	account := &Account{ID: 469804, Platform: PlatformOpenAI, Type: AccountTypeAPIKey, ProxyID: &proxyIDplaceholder
	svc := &OpenAIGatewayService{cfg: &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholderplaceholder
	// collapseInterval 0: the loop below records within the production collapse
	// window and must count as distinct failure events here.
	svc.openaiProxyStreamCircuit = newOpenAIProxyStreamCircuit(openAIProxyStreamCircuitSettings{
		failureThreshold: 2,
		failureWindow:    time.Minute,
		quarantineTTL:    10 * time.Minute,
		maxEntries:       16,
placeholder)

	for _, readErr := range []error{io.ErrUnexpectedEOF, errors.New("http2: client connection lost")placeholder {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body: &openAIStreamReadThenErrorCloser{
				reader: strings.NewReader("data: {\"type\":\"response.output_text.delta\",\"delta\":\"partial\"placeholder\n\n"),
				err:    readErr,
		placeholder,
			Header: http.Header{"X-Request-Id": []string{"rid-passthrough-proxy-disconnect"placeholderplaceholder,
	placeholder

		_, err := svc.handleStreamingResponsePassthrough(c.Request.Context(), resp, c, account, time.Now(), "model", "model")
	placeholder
		var failoverErr *UpstreamFailoverError
		require.False(t, errors.As(err, &failoverErr), "post-output disconnect must not fail over inside the same stream")
		require.Contains(t, rec.Body.String(), "partial")
placeholder

	require.True(t, svc.isOpenAIProxyStreamQuarantined(context.Background(), account))
placeholder

func TestOpenAIStreamingPassthroughResponseFailedBeforeOutputReturnsFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			MaxLineSize: defaultMaxLineSize,
	placeholder,
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
			"event: response.failed",
			`data: {"type":"response.failed","error":{"message":"upstream processing failed"placeholderplaceholder`,
			"",
	placeholder, "\n"))),
		Header: http.Header{"X-Request-Id": []string{"rid-passthrough-failed"placeholderplaceholder,
placeholder

	_, err := svc.handleStreamingResponsePassthrough(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Name: "acc"placeholder, time.Now(), "", "")
placeholder
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusBadGateway, failoverErr.StatusCode)
	require.Contains(t, string(failoverErr.ResponseBody), "upstream processing failed")
	require.False(t, c.Writer.Written())
	require.Empty(t, rec.Body.String())
placeholder

func TestOpenAIStreamingPassthroughContextWindowResponseFailedBeforeOutputAppliesPassthroughRule(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			MaxLineSize: defaultMaxLineSize,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	rule := newNonFailoverPassthroughRule(http.StatusBadRequest, "input exceeds the context window", http.StatusBadRequest, "")
	rule.Platforms = []string{PlatformOpenAIplaceholder
	rule.PassthroughBody = true
	rule.CustomMessage = nil
	ruleSvc := &ErrorPassthroughService{placeholder
	ruleSvc.setLocalCache([]*model.ErrorPassthroughRule{ruleplaceholder)
	BindErrorPassthroughService(c, ruleSvc)

	upstreamMessage := "Your input exceeds the context window of this model. Please adjust your input and try again."
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(strings.NewReader(strings.Join([]string{
			"event: response.created",
			`data: {"type":"response.created","response":{"id":"resp_1"placeholderplaceholder`,
			"",
			"event: response.failed",
			`data: {"type":"response.failed","response":{"id":"resp_1","error":{"type":"invalid_request_error","code":"context_length_exceeded","message":"` + upstreamMessage + `"placeholderplaceholderplaceholder`,
			"",
	placeholder, "\n"))),
		Header: http.Header{"X-Request-Id": []string{"rid-pass-context-window-passthrough-rule"placeholderplaceholder,
placeholder

	_, err := svc.handleStreamingResponsePassthrough(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Name: "acc"placeholder, time.Now(), "", "")
placeholder
	var failoverErr *UpstreamFailoverError
	require.False(t, errors.As(err, &failoverErr))
	require.True(t, IsResponseCommitted(c))
	require.Equal(t, http.StatusBadRequest, rec.Code)
	body := rec.Body.String()
	require.Equal(t, "upstream_error", gjson.Get(body, "error.type").String())
	require.Equal(t, upstreamMessage, gjson.Get(body, "error.message").String())
	require.NotContains(t, body, "response.failed")
	require.NotContains(t, body, "Upstream request failed")
	// 命中透传规则也应记录 ops 上游错误事件（对齐 CC/Messages 与 antigravity 先例）。
	opsVal, opsRecorded := c.Get(OpsUpstreamErrorsKey)
	require.True(t, opsRecorded, "passthrough hit should record an ops upstream error event")
	opsEvents, _ := opsVal.([]*OpsUpstreamErrorEvent)
	require.NotEmpty(t, opsEvents)
placeholder

func TestOpenAIStreamingPassthroughContextWindowResponseFailedBeforeOutputWithoutRulePassesThrough(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			MaxLineSize: defaultMaxLineSize,
	placeholder,
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
			"event: response.failed",
			`data: {"type":"response.failed","response":{"id":"resp_1","error":{"type":"invalid_request_error","code":"context_length_exceeded","message":"Your input exceeds the context window of this model. Please adjust your input and try again."placeholderplaceholderplaceholder`,
			"",
	placeholder, "\n"))),
		Header: http.Header{"X-Request-Id": []string{"rid-pass-context-window-no-rule"placeholderplaceholder,
placeholder

	_, err := svc.handleStreamingResponsePassthrough(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Name: "acc"placeholder, time.Now(), "", "")
placeholder
	var failoverErr *UpstreamFailoverError
	require.False(t, errors.As(err, &failoverErr))
	body := rec.Body.String()
	require.Contains(t, body, "event: response.failed")
	require.Contains(t, body, "context_length_exceeded")
	require.Contains(t, body, "Your input exceeds the context window")
placeholder

func TestOpenAIStreamingPassthroughResponseFailedAfterOutputSanitizesVerboseResponseForClient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			MaxLineSize: defaultMaxLineSize,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	longInstructions := strings.Repeat("You are GPT-5.1 running in the Codex CLI. ", 20)
	failedPayload := fmt.Sprintf(
		`{"type":"response.failed","response":{"id":"resp_pass_failed","object":"response","created_at":1782446336,"status":"failed","instructions":%q,"output":[{"type":"message","content":[{"type":"output_text","text":"large"placeholder]placeholder],"usage":{"input_tokens":123,"output_tokens":0placeholder,"error":{"code":"context_length_exceeded","message":"Your input exceeds the context window of this model. Please adjust your input and try again."placeholderplaceholderplaceholder`,
		longInstructions,
	)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(strings.NewReader(strings.Join([]string{
			"event: response.created",
			`data: {"type":"response.created","response":{"id":"resp_pass_failed"placeholderplaceholder`,
			"",
			"event: response.output_text.delta",
			`data: {"type":"response.output_text.delta","delta":"partial"placeholder`,
			"",
			"event: response.failed",
			"data: " + failedPayload,
			"",
	placeholder, "\n"))),
		Header: http.Header{"X-Request-Id": []string{"rid-pass-failed-after-output"placeholderplaceholder,
placeholder

	_, err := svc.handleStreamingResponsePassthrough(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Name: "acc"placeholder, time.Now(), "", "")
placeholder

	body := rec.Body.String()
	require.Contains(t, body, "event: response.failed")
	require.Contains(t, body, "context_length_exceeded")
	require.Contains(t, body, `"type":"invalid_request_error"`)
	require.Contains(t, body, "Your input exceeds the context window")
	require.NotContains(t, body, "You are GPT-5.1 running in the Codex CLI")
	require.NotContains(t, body, `"instructions"`)
	require.NotContains(t, body, `"output"`)
	require.NotContains(t, body, `"usage"`)
placeholder

func TestOpenAIStreamingPassthroughResponseDoneWithoutDoneMarkerStillSucceeds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			MaxLineSize: defaultMaxLineSize,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       pr,
		Header:     http.Header{placeholder,
placeholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		_, _ = pw.Write([]byte("data: {\"type\":\"response.done\",\"response\":{\"usage\":{\"input_tokens\":2,\"output_tokens\":3,\"input_tokens_details\":{\"cached_tokens\":1placeholderplaceholderplaceholderplaceholder\n\n"))
placeholder()

	result, err := svc.handleStreamingResponsePassthrough(c.Request.Context(), resp, c, &Account{ID: 1placeholder, time.Now(), "", "")
	_ = pr.Close()
placeholder
	require.NotNil(t, result)
	require.NotNil(t, result.usage)
	require.Equal(t, 2, result.usage.InputTokens)
	require.Equal(t, 3, result.usage.OutputTokens)
	require.Equal(t, 1, result.usage.CacheReadInputTokens)
placeholder

func TestOpenAIStreamingPassthroughResponseIncompleteWithoutDoneMarkerStillSucceeds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			MaxLineSize: defaultMaxLineSize,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       pr,
		Header:     http.Header{placeholder,
placeholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		_, _ = pw.Write([]byte("data: {\"type\":\"response.incomplete\",\"response\":{\"usage\":{\"input_tokens\":2,\"output_tokens\":3,\"input_tokens_details\":{\"cached_tokens\":1placeholderplaceholderplaceholderplaceholder\n\n"))
placeholder()

	result, err := svc.handleStreamingResponsePassthrough(c.Request.Context(), resp, c, &Account{ID: 1placeholder, time.Now(), "", "")
	_ = pr.Close()
placeholder
	require.NotNil(t, result)
	require.NotNil(t, result.usage)
	require.Equal(t, 2, result.usage.InputTokens)
	require.Equal(t, 3, result.usage.OutputTokens)
	require.Equal(t, 1, result.usage.CacheReadInputTokens)
placeholder

func TestOpenAIStreamingTooLong(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               64 * 1024,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       pr,
		Header:     http.Header{placeholder,
placeholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		// 写入超过 MaxLineSize 的单行数据，触发 ErrTooLong
		payload := "data: " + strings.Repeat("a", 128*1024) + "\n"
		_, _ = pw.Write([]byte(payload))
placeholder()

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 2placeholder, time.Now(), "model", "model")
	_ = pr.Close()

	if !errors.Is(err, bufio.ErrTooLong) {
		t.Fatalf("expected ErrTooLong, got %v", err)
placeholder
	if !strings.Contains(rec.Body.String(), "\"type\":\"error\"") || !strings.Contains(rec.Body.String(), "response_too_large") {
		t.Fatalf("expected OpenAI-compatible error SSE event, got %q", rec.Body.String())
placeholder
placeholder

func TestOpenAINonStreamingContentTypePassThrough(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Security: config.SecurityConfig{
			ResponseHeaders: config.ResponseHeaderConfig{Enabled: falseplaceholder,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	body := []byte(`{"usage":{"input_tokens":1,"output_tokens":2,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholder`)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(body)),
		Header:     http.Header{"Content-Type": []string{"application/vnd.test+json"placeholderplaceholder,
placeholder

	_, err := svc.handleNonStreamingResponse(c.Request.Context(), resp, c, &Account{placeholder, "model", "model")
	if err != nil {
		t.Fatalf("handleNonStreamingResponse error: %v", err)
placeholder

	if !strings.Contains(rec.Header().Get("Content-Type"), "application/vnd.test+json") {
		t.Fatalf("expected Content-Type passthrough, got %q", rec.Header().Get("Content-Type"))
placeholder
placeholder

func TestOpenAINonStreamingContentTypeDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Security: config.SecurityConfig{
			ResponseHeaders: config.ResponseHeaderConfig{Enabled: falseplaceholder,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	body := []byte(`{"usage":{"input_tokens":1,"output_tokens":2,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholder`)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(body)),
		Header:     http.Header{placeholder,
placeholder

	_, err := svc.handleNonStreamingResponse(c.Request.Context(), resp, c, &Account{placeholder, "model", "model")
	if err != nil {
		t.Fatalf("handleNonStreamingResponse error: %v", err)
placeholder

	if !strings.Contains(rec.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("expected default Content-Type, got %q", rec.Header().Get("Content-Type"))
placeholder
placeholder

func TestOpenAIStreamingHeadersOverride(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Security: config.SecurityConfig{
			ResponseHeaders: config.ResponseHeaderConfig{Enabled: falseplaceholder,
	placeholder,
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       pr,
		Header: http.Header{
			"Cache-Control": []string{"upstream"placeholder,
			"X-Request-Id":  []string{"req-123"placeholder,
			"Content-Type":  []string{"application/custom"placeholder,
	placeholder,
placeholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		_, _ = pw.Write([]byte("data: {\"type\":\"response.completed\",\"response\":{placeholderplaceholder\n\n"))
placeholder()

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1placeholder, time.Now(), "model", "model")
	_ = pr.Close()
	if err != nil {
		t.Fatalf("handleStreamingResponse error: %v", err)
placeholder

	if rec.Header().Get("Cache-Control") != "no-cache" {
		t.Fatalf("expected Cache-Control override, got %q", rec.Header().Get("Cache-Control"))
placeholder
	if rec.Header().Get("Content-Type") != "text/event-stream" {
		t.Fatalf("expected Content-Type override, got %q", rec.Header().Get("Content-Type"))
placeholder
	if rec.Header().Get("X-Request-Id") != "req-123" {
		t.Fatalf("expected X-Request-Id passthrough, got %q", rec.Header().Get("X-Request-Id"))
placeholder
placeholder

func TestOpenAIStreamingReuseScannerBufferAndStillWorks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       pr,
		Header:     http.Header{placeholder,
placeholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		_, _ = pw.Write([]byte("data: {\"type\":\"response.completed\",\"response\":{\"usage\":{\"input_tokens\":1,\"output_tokens\":2,\"input_tokens_details\":{\"cached_tokens\":3placeholderplaceholderplaceholderplaceholder\n\n"))
placeholder()

	result, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1placeholder, time.Now(), "model", "model")
	_ = pr.Close()
placeholder
	require.NotNil(t, result)
	require.NotNil(t, result.usage)
	require.Equal(t, 1, result.usage.InputTokens)
	require.Equal(t, 2, result.usage.OutputTokens)
	require.Equal(t, 3, result.usage.CacheReadInputTokens)
placeholder

func TestOpenAIInvalidBaseURLWhenAllowlistDisabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Security: config.SecurityConfig{
			URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholder,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	account := &Account{
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
placeholder"base_url": "://invalid-url"placeholder,
placeholder

	_, err := svc.buildUpstreamRequest(c.Request.Context(), c, account, []byte("{placeholder"), "token", false, "", false)
	if err == nil {
		t.Fatalf("expected error for invalid base_url when allowlist disabled")
placeholder
placeholder

func TestOpenAIValidateUpstreamBaseURLDisabledRequiresHTTPS(t *testing.T) {
	cfg := &config.Config{
		Security: config.SecurityConfig{
			URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholder,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	if _, err := svc.validateUpstreamBaseURL("http://not-https.example.com"); err == nil {
		t.Fatalf("expected http to be rejected when allow_insecure_http is false")
placeholder
	normalized, err := svc.validateUpstreamBaseURL("https://example.com")
	if err != nil {
		t.Fatalf("expected https to be allowed when allowlist disabled, got %v", err)
placeholder
	if normalized != "https://example.com" {
		t.Fatalf("expected raw url passthrough, got %q", normalized)
placeholder
placeholder

func TestOpenAIValidateUpstreamBaseURLDisabledAllowsHTTP(t *testing.T) {
	cfg := &config.Config{
		Security: config.SecurityConfig{
			URLAllowlist: config.URLAllowlistConfig{
				Enabled:           false,
				AllowInsecureHTTP: true,
		placeholder,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	normalized, err := svc.validateUpstreamBaseURL("http://not-https.example.com")
	if err != nil {
		t.Fatalf("expected http allowed when allow_insecure_http is true, got %v", err)
placeholder
	if normalized != "http://not-https.example.com" {
		t.Fatalf("expected raw url passthrough, got %q", normalized)
placeholder
placeholder

func TestOpenAIValidateUpstreamBaseURLEnabledEnforcesAllowlist(t *testing.T) {
	cfg := &config.Config{
		Security: config.SecurityConfig{
			URLAllowlist: config.URLAllowlistConfig{
				Enabled:       true,
				UpstreamHosts: []string{"example.com"placeholder,
		placeholder,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	if _, err := svc.validateUpstreamBaseURL("https://example.com"); err != nil {
		t.Fatalf("expected allowlisted host to pass, got %v", err)
placeholder
	if _, err := svc.validateUpstreamBaseURL("https://evil.com"); err == nil {
		t.Fatalf("expected non-allowlisted host to fail")
placeholder
placeholder

func TestOpenAIUpdateCodexUsageSnapshotFromHeaders(t *testing.T) {
	repo := &snapshotUpdateAccountRepo{updateExtraCalls: make(chan map[string]any, 1)placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder
	headers := http.Header{placeholder
	headers.Set("x-codex-primary-used-percent", "12")
	headers.Set("x-codex-secondary-used-percent", "34")
	headers.Set("x-codex-primary-window-minutes", "300")
	headers.Set("x-codex-secondary-window-minutes", "10080")
	headers.Set("x-codex-primary-reset-after-seconds", "600")
	headers.Set("x-codex-secondary-reset-after-seconds", "86400")

	svc.UpdateCodexUsageSnapshotFromHeaders(context.Background(), 123, headers)

	select {
	case updates := <-repo.updateExtraCalls:
		require.Equal(t, 12.0, updates["codex_5h_used_percent"])
		require.Equal(t, 34.0, updates["codex_7d_used_percent"])
		require.Equal(t, 600, updates["codex_5h_reset_after_seconds"])
		require.Equal(t, 86400, updates["codex_7d_reset_after_seconds"])
	case <-time.After(2 * time.Second):
		t.Fatal("expected UpdateExtra to be called")
placeholder
placeholder

func TestOpenAIResponsesRequestPathSuffix(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	tests := []struct {
		name string
		path string
		want string
placeholder{
		{name: "exact v1 responses", path: "/v1/responses", want: ""placeholder,
		{name: "compact v1 responses", path: "/v1/responses/compact", want: "/compact"placeholder,
		{name: "compact alias responses", path: "/responses/compact/", want: "/compact"placeholder,
		{name: "nested suffix", path: "/openai/v1/responses/compact/detail", want: "/compact/detail"placeholder,
		{name: "unrelated path", path: "/v1/chat/completions", want: ""placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.Request = httptest.NewRequest(http.MethodPost, tt.path, nil)
			require.Equal(t, tt.want, openAIResponsesRequestPathSuffix(c))
	placeholder)
placeholder
placeholder

func TestNormalizeOpenAICompactRequestBodyPreservesCurrentCodexPayloadFields(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5","input":[{"type":"message","role":"user","content":"compact me"placeholder],"instructions":"compact-test","tools":[{"type":"function","name":"shell"placeholder],"parallel_tool_calls":true,"reasoning":{"effort":"high"placeholder,"text":{"verbosity":"low"placeholder,"previous_response_id":"resp_123","store":true,"stream":true,"prompt_cache_key":"cache_123"placeholder`)

	normalized, changed, err := normalizeOpenAICompactRequestBody(body)

placeholder
	require.True(t, changed)
	require.Equal(t, "gpt-5.5", gjson.GetBytes(normalized, "model").String())
	require.True(t, gjson.GetBytes(normalized, "tools").Exists())
	require.True(t, gjson.GetBytes(normalized, "parallel_tool_calls").Bool())
	require.Equal(t, "high", gjson.GetBytes(normalized, "reasoning.effort").String())
	require.Equal(t, "low", gjson.GetBytes(normalized, "text.verbosity").String())
	require.Equal(t, "resp_123", gjson.GetBytes(normalized, "previous_response_id").String())
	require.False(t, gjson.GetBytes(normalized, "store").Exists())
	require.False(t, gjson.GetBytes(normalized, "stream").Exists())
	require.False(t, gjson.GetBytes(normalized, "prompt_cache_key").Exists())
placeholder

func TestOpenAIBuildUpstreamRequestOpenAIPassthroughPreservesCompactPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses/compact", bytes.NewReader([]byte(`{"model":"gpt-5"placeholder`)))

	svc := &OpenAIGatewayService{placeholder
	account := &Account{Type: AccountTypeOAuthplaceholder

	req, err := svc.buildUpstreamRequestOpenAIPassthrough(c.Request.Context(), c, account, []byte(`{"model":"gpt-5"placeholder`), "token")
placeholder
	require.Equal(t, chatgptCodexURL+"/compact", req.URL.String())
	require.Equal(t, "application/json", req.Header.Get("Accept"))
	require.Equal(t, codexCLIVersion, req.Header.Get("Version"))
	require.NotEmpty(t, req.Header.Get("Session_Id"))
	require.Equal(t, HTTPUpstreamProfileOpenAI, HTTPUpstreamProfileFromContext(req.Context()))
placeholder

func TestOpenAIBuildUpstreamRequestCompactForcesJSONAcceptForOAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses/compact", bytes.NewReader([]byte(`{"model":"gpt-5"placeholder`)))

	svc := &OpenAIGatewayService{placeholder
	account := &Account{
		Type:        AccountTypeOAuth,
placeholder"chatgpt_account_id": "chatgpt-acc"placeholder,
placeholder

	req, err := svc.buildUpstreamRequest(c.Request.Context(), c, account, []byte(`{"model":"gpt-5"placeholder`), "token", false, "", true)
placeholder
	require.Equal(t, chatgptCodexURL+"/compact", req.URL.String())
	require.Equal(t, "application/json", req.Header.Get("Accept"))
	require.Equal(t, codexCLIVersion, req.Header.Get("Version"))
	require.NotEmpty(t, req.Header.Get("Session_Id"))
	require.Equal(t, HTTPUpstreamProfileOpenAI, HTTPUpstreamProfileFromContext(req.Context()))
placeholder

func TestOpenAIBuildUpstreamRequestOAuthMessagesBridgeUsesSessionOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	body := []byte(`{"model":"gpt-5.5","prompt_cache_key":"anthropic-metadata-session-1","input":[{"type":"message","role":"developer","content":[{"type":"input_text","text":"<sub2api-claude-code-todo-guard>"placeholder]placeholder,{"type":"message","role":"user","content":"hello"placeholder]placeholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))
	c.Request.Header.Set("OpenAI-Beta", "responses=experimental")
	c.Request.Header.Set("originator", "codex_cli_rs")

	svc := &OpenAIGatewayService{placeholder
	account := &Account{
		Type:        AccountTypeOAuth,
placeholder"chatgpt_account_id": "chatgpt-acc"placeholder,
placeholder

	req, err := svc.buildUpstreamRequest(c.Request.Context(), c, account, body, "token", true, "anthropic-metadata-session-1", false)
placeholder
	require.NotEmpty(t, req.Header.Get("Session_Id"))
	require.Empty(t, req.Header.Get("Conversation_Id"))
	require.Empty(t, req.Header.Get("OpenAI-Beta"))
	require.Empty(t, req.Header.Get("originator"))
placeholder

func TestOpenAIBuildUpstreamRequestPreservesCompactPathForAPIKeyBaseURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/responses/compact", bytes.NewReader([]byte(`{"model":"gpt-5"placeholder`)))

	svc := &OpenAIGatewayService{cfg: &config.Config{
		Security: config.SecurityConfig{
			URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholder,
	placeholder,
placeholderplaceholder
	account := &Account{
		Type:        AccountTypeAPIKey,
		Platform:    PlatformOpenAI,
placeholder"base_url": "https://example.com/v1"placeholder,
placeholder

	req, err := svc.buildUpstreamRequest(c.Request.Context(), c, account, []byte(`{"model":"gpt-5"placeholder`), "token", false, "", false)
placeholder
	require.Equal(t, "https://example.com/v1/responses/compact", req.URL.String())
placeholder

func TestOpenAIBuildUpstreamRequestPreservesCodexIdentityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	c.Request.Header.Set("User-Agent", "codex_cli_rs/0.144.1")
	c.Request.Header.Set("X-Codex-Window-ID", "window-http")
	c.Request.Header.Set("X-Codex-Installation-ID", "installation-http")
	c.Request.Header.Set("X-Test", "blocked")

	body := []byte(`{"model":"gpt-5","input":"hello"placeholder`)
	svc := &OpenAIGatewayService{cfg: &config.Config{
		Security: config.SecurityConfig{
			URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholder,
	placeholder,
placeholderplaceholder
	account := &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder

	req, err := svc.buildUpstreamRequest(c.Request.Context(), c, account, body, "token", false, "", true)
placeholder
	require.Equal(t, "window-http", req.Header.Get("X-Codex-Window-ID"))
	require.Equal(t, "installation-http", req.Header.Get("X-Codex-Installation-ID"))
	require.Empty(t, req.Header.Get("X-Test"))
	require.True(t, openai.EvaluateEngineFingerprint(req.Header, body, openai.DefaultEngineFingerprintSignals))
placeholder

func TestOpenAIBuildUpstreamRequestOAuthOfficialClientOriginatorCompatibility(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 强制统一出口：客户端自报的 originator / User-Agent 都不参与上游身份构造，
	// 一律改写为网关规范身份，天然满足 originator 与 UA 首段配套的上游校验（issue #3901）。
	tests := []struct {
		name       string
		userAgent  string
		originator string
placeholder{
		{name: "official desktop ua", userAgent: "Codex Desktop/1.2.3"placeholder,
		{
			name:       "mismatched originator",
			userAgent:  "codex_vscode/0.140.2 (Mac OS X 14.0; arm64) vscode (codex_vscode; 0.140.2)",
			originator: "codex_cli_rs",
	placeholder,
		{
			name:       "tui identity",
			userAgent:  "codex-tui/0.140.2 (Mac OS X 14.0; arm64) iTerm (codex-tui; 0.140.2)",
			originator: "codex-tui",
	placeholder,
		{name: "official originator without ua", originator: "codex_vscode"placeholder,
		{name: "third-party ua", userAgent: "luna/1.2.0"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader([]byte(`{"model":"gpt-5"placeholder`)))
			if tt.userAgent != "" {
				c.Request.Header.Set("User-Agent", tt.userAgent)
		placeholder
			if tt.originator != "" {
				c.Request.Header.Set("originator", tt.originator)
		placeholder

			svc := &OpenAIGatewayService{placeholder
			account := &Account{
				Type:        AccountTypeOAuth,
		placeholder"chatgpt_account_id": "chatgpt-acc"placeholder,
		placeholder

			isCodexCLI := openai.IsCodexOfficialClientByHeaders(c.GetHeader("User-Agent"), c.GetHeader("originator"))
			req, err := svc.buildUpstreamRequest(c.Request.Context(), c, account, []byte(`{"model":"gpt-5"placeholder`), "token", false, "", isCodexCLI)
		placeholder
			require.Equal(t, openai.CodexDefaultOriginator, req.Header.Get("originator"))
			require.Equal(t, codexCLIUserAgent, req.Header.Get("User-Agent"))
			require.Equal(t, codexCLIVersion, req.Header.Get("version"))
	placeholder)
placeholder
placeholder

// ==================== P1-08 修复：model 替换性能优化测试 ====================

// ==================== P1-08 修复：model 替换性能优化测试 =============
func TestReplaceModelInSSELine(t *testing.T) {
	svc := &OpenAIGatewayService{placeholder

	tests := []struct {
		name     string
		line     string
		from     string
		to       string
		expected string
placeholder{
		{
			name:     "顶层 model 字段替换",
			line:     `data: {"id":"chatcmpl-123","model":"gpt-4o","choices":[]placeholder`,
			from:     "gpt-4o",
			to:       "my-custom-model",
			expected: `data: {"id":"chatcmpl-123","model":"my-custom-model","choices":[]placeholder`,
	placeholder,
		{
			name:     "嵌套 response.model 替换",
			line:     `data: {"type":"response","response":{"id":"resp-1","model":"gpt-4o","output":[]placeholderplaceholder`,
			from:     "gpt-4o",
			to:       "my-model",
			expected: `data: {"type":"response","response":{"id":"resp-1","model":"my-model","output":[]placeholderplaceholder`,
	placeholder,
		{
			name:     "model 不匹配时不替换",
			line:     `data: {"id":"chatcmpl-123","model":"gpt-3.5-turbo","choices":[]placeholder`,
			from:     "gpt-4o",
			to:       "my-model",
			expected: `data: {"id":"chatcmpl-123","model":"gpt-3.5-turbo","choices":[]placeholder`,
	placeholder,
		{
			name:     "无 model 字段时不替换",
			line:     `data: {"id":"chatcmpl-123","choices":[]placeholder`,
			from:     "gpt-4o",
			to:       "my-model",
			expected: `data: {"id":"chatcmpl-123","choices":[]placeholder`,
	placeholder,
		{
			name:     "空 data 行",
			line:     `data: `,
			from:     "gpt-4o",
			to:       "my-model",
			expected: `data: `,
	placeholder,
		{
			name:     "[DONE] 行",
			line:     `data: [DONE]`,
			from:     "gpt-4o",
			to:       "my-model",
			expected: `data: [DONE]`,
	placeholder,
		{
			name:     "非 data: 前缀行",
			line:     `event: message`,
			from:     "gpt-4o",
			to:       "my-model",
			expected: `event: message`,
	placeholder,
		{
			name:     "非法 JSON 不替换",
			line:     `data: {invalid jsonplaceholder`,
			from:     "gpt-4o",
			to:       "my-model",
			expected: `data: {invalid jsonplaceholder`,
	placeholder,
		{
			name:     "无空格 data: 格式",
			line:     `data:{"id":"x","model":"gpt-4o"placeholder`,
			from:     "gpt-4o",
			to:       "my-model",
			expected: `data: {"id":"x","model":"my-model"placeholder`,
	placeholder,
		{
			name:     "model 名含特殊字符",
			line:     `data: {"model":"org/model-v2.1-beta"placeholder`,
			from:     "org/model-v2.1-beta",
			to:       "custom/alias",
			expected: `data: {"model":"custom/alias"placeholder`,
	placeholder,
		{
			name:     "空行",
			line:     "",
			from:     "gpt-4o",
			to:       "my-model",
			expected: "",
	placeholder,
		{
			name:     "保持其他字段不变",
			line:     `data: {"id":"abc","object":"chat.completion.chunk","model":"gpt-4o","created":1234567890,"choices":[{"index":0,"delta":{"content":"hi"placeholderplaceholder]placeholder`,
			from:     "gpt-4o",
			to:       "alias",
			expected: `data: {"id":"abc","object":"chat.completion.chunk","model":"alias","created":1234567890,"choices":[{"index":0,"delta":{"content":"hi"placeholderplaceholder]placeholder`,
	placeholder,
		{
			name:     "顶层优先于嵌套：同时存在两个 model",
			line:     `data: {"model":"gpt-4o","response":{"model":"gpt-4o"placeholderplaceholder`,
			from:     "gpt-4o",
			to:       "replaced",
			expected: `data: {"model":"replaced","response":{"model":"gpt-4o"placeholderplaceholder`,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.replaceModelInSSELine(tt.line, tt.from, tt.to)
			require.Equal(t, tt.expected, got)
	placeholder)
placeholder
placeholder

func TestReplaceModelInSSEBody(t *testing.T) {
	svc := &OpenAIGatewayService{placeholder

	tests := []struct {
		name     string
		body     string
		from     string
		to       string
		expected string
placeholder{
		{
			name:     "多行 SSE body 替换",
			body:     "data: {\"model\":\"gpt-4o\",\"choices\":[]placeholder\n\ndata: {\"model\":\"gpt-4o\",\"choices\":[{\"delta\":{\"content\":\"hi\"placeholderplaceholder]placeholder\n\ndata: [DONE]\n",
			from:     "gpt-4o",
			to:       "alias",
			expected: "data: {\"model\":\"alias\",\"choices\":[]placeholder\n\ndata: {\"model\":\"alias\",\"choices\":[{\"delta\":{\"content\":\"hi\"placeholderplaceholder]placeholder\n\ndata: [DONE]\n",
	placeholder,
		{
			name:     "无需替换的 body",
			body:     "data: {\"model\":\"gpt-3.5-turbo\"placeholder\n\ndata: [DONE]\n",
			from:     "gpt-4o",
			to:       "alias",
			expected: "data: {\"model\":\"gpt-3.5-turbo\"placeholder\n\ndata: [DONE]\n",
	placeholder,
		{
			name:     "混合 event 和 data 行",
			body:     "event: message\ndata: {\"model\":\"gpt-4o\"placeholder\n\n",
			from:     "gpt-4o",
			to:       "alias",
			expected: "event: message\ndata: {\"model\":\"alias\"placeholder\n\n",
	placeholder,
		{
			name:     "空 body",
			body:     "",
			from:     "gpt-4o",
			to:       "alias",
			expected: "",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.replaceModelInSSEBody(tt.body, tt.from, tt.to)
			require.Equal(t, tt.expected, got)
	placeholder)
placeholder
placeholder

func TestReplaceModelInResponseBody(t *testing.T) {
	svc := &OpenAIGatewayService{placeholder

	tests := []struct {
		name     string
		body     string
		from     string
		to       string
		expected string
placeholder{
		{
			name:     "替换顶层 model",
			body:     `{"id":"chatcmpl-123","model":"gpt-4o","choices":[]placeholder`,
			from:     "gpt-4o",
			to:       "alias",
			expected: `{"id":"chatcmpl-123","model":"alias","choices":[]placeholder`,
	placeholder,
		{
			name:     "model 不匹配不替换",
			body:     `{"id":"chatcmpl-123","model":"gpt-3.5-turbo","choices":[]placeholder`,
			from:     "gpt-4o",
			to:       "alias",
			expected: `{"id":"chatcmpl-123","model":"gpt-3.5-turbo","choices":[]placeholder`,
	placeholder,
		{
			name:     "无 model 字段不替换",
			body:     `{"id":"chatcmpl-123","choices":[]placeholder`,
			from:     "gpt-4o",
			to:       "alias",
			expected: `{"id":"chatcmpl-123","choices":[]placeholder`,
	placeholder,
		{
			name:     "非法 JSON 返回原值",
			body:     `not json`,
			from:     "gpt-4o",
			to:       "alias",
			expected: `not json`,
	placeholder,
		{
			name:     "空 body 返回原值",
			body:     ``,
			from:     "gpt-4o",
			to:       "alias",
			expected: ``,
	placeholder,
		{
			name:     "保持嵌套结构不变",
			body:     `{"model":"gpt-4o","usage":{"prompt_tokens":10,"completion_tokens":20placeholder,"choices":[{"message":{"role":"assistant","content":"hello"placeholderplaceholder]placeholder`,
			from:     "gpt-4o",
			to:       "alias",
			expected: `{"model":"alias","usage":{"prompt_tokens":10,"completion_tokens":20placeholder,"choices":[{"message":{"role":"assistant","content":"hello"placeholderplaceholder]placeholder`,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.replaceModelInResponseBody([]byte(tt.body), tt.from, tt.to)
			require.Equal(t, tt.expected, string(got))
	placeholder)
placeholder
placeholder

func TestExtractOpenAISSEDataLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantData string
		wantOK   bool
placeholder{
		{name: "标准格式", line: `data: {"type":"x"placeholder`, wantData: `{"type":"x"placeholder`, wantOK: trueplaceholder,
		{name: "无空格格式", line: `data:{"type":"x"placeholder`, wantData: `{"type":"x"placeholder`, wantOK: trueplaceholder,
		{name: "纯空数据", line: `data:   `, wantData: ``, wantOK: trueplaceholder,
		{name: "非 data 行", line: `event: message`, wantData: ``, wantOK: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := extractOpenAISSEDataLine(tt.line)
			require.Equal(t, tt.wantOK, ok)
			require.Equal(t, tt.wantData, got)
	placeholder)
placeholder
placeholder

func TestParseSSEUsage_SelectiveParsing(t *testing.T) {
	svc := &OpenAIGatewayService{placeholder
	usage := &OpenAIUsage{InputTokens: 9, OutputTokens: 8, CacheReadInputTokens: 7placeholder

	// 非 completed 事件，不应覆盖 usage
	svc.parseSSEUsage(`{"type":"response.in_progress","response":{"usage":{"input_tokens":1,"output_tokens":2placeholderplaceholderplaceholder`, usage)
	require.Equal(t, 9, usage.InputTokens)
	require.Equal(t, 8, usage.OutputTokens)
	require.Equal(t, 7, usage.CacheReadInputTokens)

	// completed 事件，应提取 usage
	svc.parseSSEUsage(`{"type":"response.completed","response":{"usage":{"input_tokens":3,"output_tokens":5,"input_tokens_details":{"cached_tokens":2placeholderplaceholderplaceholderplaceholder`, usage)
	require.Equal(t, 3, usage.InputTokens)
	require.Equal(t, 5, usage.OutputTokens)
	require.Equal(t, 2, usage.CacheReadInputTokens)

	// done 事件同样可能携带最终 usage
	svc.parseSSEUsage(`{"type":"response.done","response":{"usage":{"input_tokens":13,"output_tokens":15,"input_tokens_details":{"cached_tokens":4placeholderplaceholderplaceholderplaceholder`, usage)
	require.Equal(t, 13, usage.InputTokens)
	require.Equal(t, 15, usage.OutputTokens)
	require.Equal(t, 4, usage.CacheReadInputTokens)

	// failed 事件在部分上游路径也会携带已消耗 usage，应与 WS/passthrough 保持一致
	svc.parseSSEUsage(`{"type":"response.failed","response":{"usage":{"input_tokens":17,"output_tokens":19,"input_tokens_details":{"cached_tokens":6placeholderplaceholderplaceholderplaceholder`, usage)
	require.Equal(t, 17, usage.InputTokens)
	require.Equal(t, 19, usage.OutputTokens)
	require.Equal(t, 6, usage.CacheReadInputTokens)

	svc.parseSSEUsage(`{"type":"response.completed","response":{"usage":{"prompt_tokens":21,"completion_tokens":8,"prompt_tokens_details":{"cached_tokens":6placeholderplaceholderplaceholderplaceholder`, usage)
	require.Equal(t, 21, usage.InputTokens)
	require.Equal(t, 8, usage.OutputTokens)
	require.Equal(t, 6, usage.CacheReadInputTokens)
placeholder

func TestExtractOpenAIUsageFromJSONBytes_AcceptsResponseAndChatUsageShapes(t *testing.T) {
	usage, ok := extractOpenAIUsageFromJSONBytes([]byte(`{"id":"resp_1","usage":{"input_tokens":9,"output_tokens":5,"input_tokens_details":{"cached_tokens":2,"cache_write_tokens":4placeholderplaceholderplaceholder`))
	require.True(t, ok)
	require.Equal(t, 9, usage.InputTokens)
	require.Equal(t, 5, usage.OutputTokens)
	require.Equal(t, 2, usage.CacheReadInputTokens)
	require.Equal(t, 4, usage.CacheCreationInputTokens)

	usage, ok = extractOpenAIUsageFromJSONBytes([]byte(`{"type":"response.completed","response":{"usage":{"prompt_tokens":13,"completion_tokens":7,"prompt_tokens_details":{"cached_tokens":4,"cache_creation_tokens":3placeholderplaceholderplaceholderplaceholder`))
	require.True(t, ok)
	require.Equal(t, 13, usage.InputTokens)
	require.Equal(t, 7, usage.OutputTokens)
	require.Equal(t, 4, usage.CacheReadInputTokens)
	require.Equal(t, 3, usage.CacheCreationInputTokens)

	usage, ok = extractOpenAIUsageFromJSONBytes([]byte(`{"usage":{"input_tokens":11,"output_tokens":2,"cache_write_input_tokens":6placeholderplaceholder`))
	require.True(t, ok)
	require.Equal(t, 6, usage.CacheCreationInputTokens)

	usage, ok = extractOpenAIUsageFromJSONBytes([]byte(`{"usage":{"input_tokens":20,"output_tokens":2,"cache_creation_input_tokens":19,"input_tokens_details":{"cache_write_tokens":7placeholderplaceholderplaceholder`))
	require.True(t, ok)
	require.Equal(t, 7, usage.CacheCreationInputTokens, "官方嵌套字段应优先于兼容顶层别名")

	usage, ok = extractOpenAIUsageFromJSONBytes([]byte(`{"usage":{"input_tokens":20,"output_tokens":2,"cache_creation_input_tokens":19,"input_tokens_details":{"cache_write_tokens":0placeholderplaceholderplaceholder`))
	require.True(t, ok)
	require.Zero(t, usage.CacheCreationInputTokens, "官方嵌套字段显式为零时仍应优先于兼容顶层别名")

	usage, ok = extractOpenAIUsageFromJSONBytes([]byte(`{"usage":{"input_tokens":20,"output_tokens":2,"cache_read_input_tokens":19,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholder`))
	require.True(t, ok)
	require.Zero(t, usage.CacheReadInputTokens, "官方嵌套缓存读取字段显式为零时仍应优先于兼容顶层别名")
placeholder

func TestExtractCodexFinalResponse_SampleReplay(t *testing.T) {
	body := strings.Join([]string{
		`event: message`,
		`data: {"type":"response.in_progress","response":{"id":"resp_1"placeholderplaceholder`,
		`data: {"type":"response.completed","response":{"id":"resp_1","model":"gpt-4o","usage":{"input_tokens":11,"output_tokens":22,"input_tokens_details":{"cached_tokens":3placeholderplaceholderplaceholderplaceholder`,
		`data: [DONE]`,
placeholder, "\n")

	finalResp, ok := extractCodexFinalResponse(body)
	require.True(t, ok)
	require.Contains(t, string(finalResp), `"id":"resp_1"`)
	require.Contains(t, string(finalResp), `"input_tokens":11`)
placeholder

func TestHandleSSEToJSON_CompletedEventReturnsJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	svc := &OpenAIGatewayService{cfg: &config.Config{placeholderplaceholder
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
placeholder
	body := []byte(strings.Join([]string{
		`data: {"type":"response.in_progress","response":{"id":"resp_2"placeholderplaceholder`,
		`data: {"type":"response.completed","response":{"id":"resp_2","model":"gpt-4o","usage":{"input_tokens":7,"output_tokens":9,"input_tokens_details":{"cached_tokens":1placeholderplaceholderplaceholderplaceholder`,
		`data: [DONE]`,
placeholder, "\n"))

	usage, err := svc.handleSSEToJSON(resp, c, body, "gpt-4o", "gpt-4o")
placeholder
	require.NotNil(t, usage)
	require.Equal(t, 7, usage.InputTokens)
	require.Equal(t, 9, usage.OutputTokens)
	require.Equal(t, 1, usage.CacheReadInputTokens)
	// Header 可能由上游 Content-Type 透传；关键是 body 已转换为最终 JSON 响应。
	require.NotContains(t, rec.Body.String(), "event:")
	require.Contains(t, rec.Body.String(), `"id":"resp_2"`)
	require.NotContains(t, rec.Body.String(), "data:")
placeholder

func TestHandleNonStreamingResponse_APIKeyFallsBackToSSEBodyWhenContentTypeIsWrong(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)

	svc := &OpenAIGatewayService{cfg: &config.Config{placeholderplaceholder
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(strings.Join([]string{
			`data: {"type":"response.output_text.delta","delta":"hel"placeholder`,
			`data: {"type":"response.output_text.delta","delta":"lo"placeholder`,
			`data: {"type":"response.completed","response":{"id":"resp_api_key_sse","object":"response","model":"gpt-5.4","status":"completed","output":[],"usage":{"input_tokens":3,"output_tokens":2,"total_tokens":5placeholderplaceholderplaceholder`,
			`data: [DONE]`,
	placeholder, "\n"))),
placeholder
	account := &Account{ID: 1, Type: AccountTypeAPIKeyplaceholder

	result, err := svc.handleNonStreamingResponse(context.Background(), resp, c, account, "gpt-5.4", "gpt-5.4")
placeholder
	require.NotNil(t, result)
	require.Equal(t, 3, result.InputTokens)
	require.Equal(t, 2, result.OutputTokens)
	require.NotContains(t, rec.Body.String(), "data:")
	require.Equal(t, "resp_api_key_sse", gjson.Get(rec.Body.String(), "id").String())
	require.Equal(t, "hello", gjson.Get(rec.Body.String(), "output.0.content.0.text").String())
placeholder

func TestHandleNonStreamingResponse_OAuthJSONBodyWithDataEventTextKeepsJSONUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses/compact", nil)

	svc := &OpenAIGatewayService{cfg: &config.Config{placeholderplaceholder
	// Plain JSON compact response whose output text happens to contain the
	// literal substrings "data:" and "event:" (e.g. echoing shell/log output).
	// This must NOT be misdetected as SSE framing: it has a top-level usage
	// object and no upstream text/event-stream Content-Type.
	jsonBody := `{"id":"resp_oauth_compact","object":"response","model":"gpt-5.4","status":"completed",` +
		`"output":[{"type":"message","content":[{"type":"output_text",` +
		`"text":"processing data: 1,2,3 then event: click finished"placeholder]placeholder],` +
		`"usage":{"input_tokens":11,"output_tokens":22,"total_tokens":33placeholderplaceholder`
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(jsonBody)),
placeholder
	account := &Account{ID: 146, Type: AccountTypeOAuthplaceholder

	result, err := svc.handleNonStreamingResponse(context.Background(), resp, c, account, "gpt-5.4", "gpt-5.4")
placeholder
	require.NotNil(t, result)
	require.Equal(t, 11, result.InputTokens)
	require.Equal(t, 22, result.OutputTokens)
	// Response must remain the original JSON body (not routed through the SSE
	// path, which would rewrite/lose the body or usage).
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	require.Equal(t, "resp_oauth_compact", gjson.Get(rec.Body.String(), "id").String())
	require.Equal(t, int64(33), gjson.Get(rec.Body.String(), "usage.total_tokens").Int())
	require.Contains(t, rec.Body.String(), "processing data: 1,2,3 then event: click finished")
placeholder

func TestHandleSSEToJSON_ReconstructsImageGenerationOutputItemDone(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	svc := &OpenAIGatewayService{cfg: &config.Config{placeholderplaceholder
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
placeholder
	body := []byte(strings.Join([]string{
		`data: {"type":"response.output_item.done","item":{"id":"ig_123","type":"image_generation_call","status":"generating","result":"aGVsbG8=","revised_prompt":"draw a cat","output_format":"png"placeholderplaceholder`,
		`data: {"type":"response.completed","response":{"id":"resp_img","model":"gpt-5.4","output":[],"usage":{"input_tokens":7,"output_tokens":9,"output_tokens_details":{"image_tokens":4placeholderplaceholderplaceholderplaceholder`,
		`data: [DONE]`,
placeholder, "\n"))

	usage, err := svc.handleSSEToJSON(resp, c, body, "gpt-5.4", "gpt-5.4")
placeholder
	require.NotNil(t, usage)
	require.Equal(t, 4, usage.ImageOutputTokens)
	require.NotContains(t, rec.Body.String(), "data:")
	require.Equal(t, "image_generation_call", gjson.Get(rec.Body.String(), "output.0.type").String())
	require.Equal(t, "completed", gjson.Get(rec.Body.String(), "output.0.status").String())
	require.Equal(t, "aGVsbG8=", gjson.Get(rec.Body.String(), "output.0.result").String())
	require.Equal(t, "draw a cat", gjson.Get(rec.Body.String(), "output.0.revised_prompt").String())
placeholder

func TestHandleSSEToJSON_NoFinalResponseKeepsSSEBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	svc := &OpenAIGatewayService{cfg: &config.Config{placeholderplaceholder
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
placeholder
	body := []byte(strings.Join([]string{
		`data: {"type":"response.in_progress","response":{"id":"resp_3"placeholderplaceholder`,
		`data: [DONE]`,
placeholder, "\n"))

	usage, err := svc.handleSSEToJSON(resp, c, body, "gpt-4o", "gpt-4o")
placeholder
	require.NotNil(t, usage)
	require.Equal(t, 0, usage.InputTokens)
	require.Contains(t, rec.Header().Get("Content-Type"), "text/event-stream")
	require.Contains(t, rec.Body.String(), `data: {"type":"response.in_progress"`)
placeholder

func TestHandleSSEToJSON_ResponseFailedReturnsProtocolError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	svc := &OpenAIGatewayService{cfg: &config.Config{placeholderplaceholder
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
placeholder
	body := []byte(strings.Join([]string{
		`data: {"type":"response.failed","error":{"message":"upstream rejected request"placeholderplaceholder`,
		`data: [DONE]`,
placeholder, "\n"))

	usage, err := svc.handleSSEToJSON(resp, c, body, "gpt-4o", "gpt-4o")
	require.Nil(t, usage)
placeholder
	require.Equal(t, http.StatusBadGateway, rec.Code)
	require.Contains(t, rec.Body.String(), "upstream rejected request")
	require.Contains(t, rec.Header().Get("Content-Type"), "application/json")
placeholder

func TestOpenAICompatSSEFrameParserResetsEventTypeAtFrameBoundary(t *testing.T) {
	var parser openAICompatSSEFrameParser

	frame, ok := parser.AddLine("event: response.created")
	require.False(t, ok)
	require.Empty(t, frame)

	frame, ok = parser.AddLine(`data: {"response":{"id":"resp_1"placeholderplaceholder`)
	require.False(t, ok)
	require.Empty(t, frame)

	frame, ok = parser.AddLine("")
	require.True(t, ok)
	require.Equal(t, "response.created", frame.EventType)
	require.JSONEq(t, `{"response":{"id":"resp_1"placeholderplaceholder`, frame.Data)

	frame, ok = parser.AddLine(`data: {"delta":"ok"placeholder`)
	require.False(t, ok)
	require.Empty(t, frame.EventType)

	frame, ok = parser.AddLine("")
	require.True(t, ok)
	require.Empty(t, frame.EventType)
	require.JSONEq(t, `{"delta":"ok"placeholder`, frame.Data)
placeholder

func TestStreamingPassthroughCyberPolicyMarksAndPassesThrough(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(strings.NewReader(strings.Join([]string{
			"event: response.created",
			`data: {"type":"response.created","response":{"id":"r1"placeholderplaceholder`,
			"",
			"event: response.failed",
			`data: {"type":"response.failed","response":{"error":{"code":"cyber_policy","message":"flagged for cyber policy"placeholderplaceholderplaceholder`,
			"",
	placeholder, "\n"))),
		Header: http.Header{"X-Request-Id": []string{"rid-cyber"placeholderplaceholder,
placeholder

	_, err := svc.handleStreamingResponsePassthrough(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Name: "a"placeholder, time.Now(), "m", "m")
placeholder
	var failoverErr *UpstreamFailoverError
	require.False(t, errors.As(err, &failoverErr), "cyber must NOT failover")
	require.Contains(t, rec.Body.String(), "cyber_policy", "response.failed passed through to client")
	mark := GetOpsCyberPolicy(c)
	require.NotNil(t, mark)
	require.Equal(t, "flagged for cyber policy", mark.Message)
placeholder

func TestHandleStreamingResponseCyberPolicyMarks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{Gateway: config.GatewayConfig{MaxLineSize: defaultMaxLineSizeplaceholderplaceholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(strings.NewReader(strings.Join([]string{
			"event: response.created",
			`data: {"type":"response.created","response":{"id":"r1"placeholderplaceholder`,
			"",
			"event: response.failed",
			`data: {"type":"response.failed","error":{"code":"cyber_policy","message":"flagged"placeholderplaceholder`,
			"",
	placeholder, "\n"))),
		Header: http.Header{"X-Request-Id": []string{"rid"placeholderplaceholder,
placeholder
	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Name: "a"placeholder, time.Now(), "m", "m")
placeholder
	var fo *UpstreamFailoverError
	require.False(t, errors.As(err, &fo))
	require.NotNil(t, GetOpsCyberPolicy(c))
placeholder

func TestHandleErrorResponseCyberPolicyPassthrough(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &OpenAIGatewayService{cfg: &config.Config{placeholderplaceholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	cyberBody := `{"error":{"code":"cyber_policy","message":"flagged for cyber policy"placeholderplaceholder`
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholder, "X-Request-Id": []string{"rid"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(cyberBody)),
placeholder
	_, err := svc.handleErrorResponse(context.Background(), resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Name: "a"placeholder, nil)
placeholder
	require.Equal(t, http.StatusBadRequest, rec.Code, "passthrough upstream 400, not rewrapped 502")
	require.Contains(t, rec.Body.String(), "cyber_policy", "client sees original cyber body")
	require.NotContains(t, rec.Body.String(), "Upstream request failed", "must not 502-rewrap")
	mark := GetOpsCyberPolicy(c)
	require.NotNil(t, mark)
	require.Equal(t, http.StatusBadRequest, mark.UpstreamStatus)
placeholder

func TestHandleCompatErrorResponseCyberPolicyEarlyReturn(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &OpenAIGatewayService{cfg: &config.Config{placeholderplaceholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	cyberBody := `{"error":{"code":"cyber_policy","message":"flagged for cyber policy"placeholderplaceholder`
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(cyberBody)),
placeholder
	var gotStatus int
	var gotType, gotMsg string
	writeError := func(_ *gin.Context, statusCode int, errType, message string) {
		gotStatus, gotType, gotMsg = statusCode, errType, message
placeholder
	// cyber 命中应早返回(写兼容错误 + 不冷却账号)，而非落到通用 "Upstream request failed"。
	_, err := svc.handleCompatErrorResponse(resp, c, &Account{ID: 1, Platform: PlatformOpenAI, Name: "a"placeholder, writeError)
placeholder
	require.Equal(t, http.StatusBadRequest, gotStatus)
	require.Equal(t, "invalid_request_error", gotType)
	require.Contains(t, gotMsg, "flagged for cyber policy")
	require.NotContains(t, gotMsg, "Upstream request failed")
	require.NotNil(t, GetOpsCyberPolicy(c))
placeholder
