package service

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	coderws "github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

type openAIWSRateLimitSignalRepo struct {
	stubOpenAIAccountRepo
	rateLimitCalls []time.Time
	updateExtra    []map[string]any
placeholder

type openAICodexSnapshotAsyncRepo struct {
	stubOpenAIAccountRepo
	updateExtraCh chan map[string]any
	rateLimitCh   chan time.Time
placeholder

type openAICodexExtraListRepo struct {
	stubOpenAIAccountRepo
	rateLimitCh chan time.Time
placeholder

func (r *openAIWSRateLimitSignalRepo) SetRateLimited(_ context.Context, _ int64, resetAt time.Time) error {
	r.rateLimitCalls = append(r.rateLimitCalls, resetAt)
	return nil
placeholder

func (r *openAIWSRateLimitSignalRepo) UpdateExtra(_ context.Context, _ int64, updates map[string]any) error {
	copied := make(map[string]any, len(updates))
	for k, v := range updates {
		copied[k] = v
placeholder
	r.updateExtra = append(r.updateExtra, copied)
	return nil
placeholder

func (r *openAICodexSnapshotAsyncRepo) SetRateLimited(_ context.Context, _ int64, resetAt time.Time) error {
	if r.rateLimitCh != nil {
		r.rateLimitCh <- resetAt
placeholder
	return nil
placeholder

func (r *openAICodexSnapshotAsyncRepo) UpdateExtra(_ context.Context, _ int64, updates map[string]any) error {
	if r.updateExtraCh != nil {
		copied := make(map[string]any, len(updates))
		for k, v := range updates {
			copied[k] = v
	placeholder
		r.updateExtraCh <- copied
placeholder
	return nil
placeholder

func (r *openAICodexExtraListRepo) SetRateLimited(_ context.Context, _ int64, resetAt time.Time) error {
	if r.rateLimitCh != nil {
		r.rateLimitCh <- resetAt
placeholder
	return nil
placeholder

func (r *openAICodexExtraListRepo) ListWithFilters(_ context.Context, params pagination.PaginationParams, platform, accountType, status, search string, groupID int64, privacyMode string) ([]Account, *pagination.PaginationResult, error) {
	_ = platform
	_ = accountType
	_ = status
	_ = search
	_ = groupID
	_ = privacyMode
	return r.accounts, &pagination.PaginationResult{Total: int64(len(r.accounts)), Page: params.Page, PageSize: params.PageSizeplaceholder, nil
placeholder

func TestOpenAIGatewayService_Forward_WSv2ErrorEventUsageLimitPersistsRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	resetAt := time.Now().Add(2 * time.Hour).Unix()
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true placeholderplaceholder
	wsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket failed: %v", err)
			return
	placeholder
		defer func() { _ = conn.Close() placeholder()

		var req map[string]any
		if err := conn.ReadJSON(&req); err != nil {
			t.Errorf("read ws request failed: %v", err)
			return
	placeholder
		_ = conn.WriteJSON(map[string]any{
			"type": "error",
			"error": map[string]any{
				"code":      "rate_limit_exceeded",
				"type":      "usage_limit_reached",
				"message":   "The usage limit has been reached",
				"resets_at": resetAt,
		placeholder,
	placeholder)
placeholder))
	defer wsServer.Close()

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	c.Request.Header.Set("User-Agent", "unit-test-agent/1.0")

	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"id":"resp_http_should_not_run"placeholder`)),
	placeholder,
placeholder

	cfg := newOpenAIWSV2TestConfig()
	cfg.Security.URLAllowlist.Enabled = false
	cfg.Security.URLAllowlist.AllowInsecureHTTP = true

	account := Account{
		ID:          501,
		Name:        "openai-ws-rate-limit-event",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": wsServer.URL,
	placeholder,
		Extra: map[string]any{
			"responses_websockets_v2_enabled": true,
	placeholder,
placeholder
	repo := &openAIWSRateLimitSignalRepo{stubOpenAIAccountRepo: stubOpenAIAccountRepo{accounts: []Account{accountplaceholderplaceholderplaceholder
	rateSvc := &RateLimitService{accountRepo: repoplaceholder
	svc := &OpenAIGatewayService{
		accountRepo:      repo,
		rateLimitService: rateSvc,
		httpUpstream:     upstream,
		cache:            &stubGatewayCache{placeholder,
		cfg:              cfg,
		openaiWSResolver: NewOpenAIWSProtocolResolver(cfg),
		toolCorrector:    NewCodexToolCorrector(),
placeholder

	body := []byte(`{"model":"gpt-5.1","stream":false,"input":[{"type":"input_text","text":"hello"placeholder]placeholder`)
	result, err := svc.Forward(context.Background(), c, &account, body)
placeholder
	require.Nil(t, result)
	require.Equal(t, http.StatusTooManyRequests, rec.Code)
	require.Nil(t, upstream.lastReq, "WS 限流 error event 不应回退到同账号 HTTP")
	require.Len(t, repo.rateLimitCalls, 1)
	require.WithinDuration(t, time.Unix(resetAt, 0), repo.rateLimitCalls[0], 2*time.Second)
placeholder

func TestOpenAIGatewayService_Forward_WSv2Handshake429PersistsRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-codex-primary-used-percent", "100")
		w.Header().Set("x-codex-primary-reset-after-seconds", "7200")
		w.Header().Set("x-codex-primary-window-minutes", "10080")
		w.Header().Set("x-codex-secondary-used-percent", "3")
		w.Header().Set("x-codex-secondary-reset-after-seconds", "1800")
		w.Header().Set("x-codex-secondary-window-minutes", "300")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":{"type":"rate_limit_exceeded","message":"rate limited"placeholderplaceholder`))
placeholder))
	defer server.Close()

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	c.Request.Header.Set("User-Agent", "unit-test-agent/1.0")

	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"id":"resp_http_should_not_run"placeholder`)),
	placeholder,
placeholder

	cfg := newOpenAIWSV2TestConfig()
	cfg.Security.URLAllowlist.Enabled = false
	cfg.Security.URLAllowlist.AllowInsecureHTTP = true

	account := Account{
		ID:          502,
		Name:        "openai-ws-rate-limit-handshake",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": server.URL,
	placeholder,
		Extra: map[string]any{
			"responses_websockets_v2_enabled": true,
	placeholder,
placeholder
	repo := &openAIWSRateLimitSignalRepo{stubOpenAIAccountRepo: stubOpenAIAccountRepo{accounts: []Account{accountplaceholderplaceholderplaceholder
	rateSvc := &RateLimitService{accountRepo: repoplaceholder
	svc := &OpenAIGatewayService{
		accountRepo:      repo,
		rateLimitService: rateSvc,
		httpUpstream:     upstream,
		cache:            &stubGatewayCache{placeholder,
		cfg:              cfg,
		openaiWSResolver: NewOpenAIWSProtocolResolver(cfg),
		toolCorrector:    NewCodexToolCorrector(),
placeholder

	body := []byte(`{"model":"gpt-5.1","stream":false,"input":[{"type":"input_text","text":"hello"placeholder]placeholder`)
	result, err := svc.Forward(context.Background(), c, &account, body)
placeholder
	require.Nil(t, result)
	require.Equal(t, http.StatusTooManyRequests, rec.Code)
	require.Nil(t, upstream.lastReq, "WS 握手 429 不应回退到同账号 HTTP")
	require.Len(t, repo.rateLimitCalls, 1)
	require.NotEmpty(t, repo.updateExtra, "握手 429 的 x-codex 头应立即落库")
	require.Contains(t, repo.updateExtra[0], "codex_usage_updated_at")
placeholder

func TestOpenAIGatewayService_Forward_WSv2Handshake502RecordsModelTransient(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-request-id", "req-ws-502")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"error":{"type":"server_error","message":"bad gateway"placeholderplaceholder`))
placeholder))
	defer server.Close()

	cfg := newOpenAIWSV2TestConfig()
	cfg.Security.URLAllowlist.Enabled = false
	cfg.Security.URLAllowlist.AllowInsecureHTTP = true
	account := Account{
		ID:          504,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
placeholder"api_key": "sk-test", "base_url": server.URLplaceholder,
		Extra:       map[string]any{"responses_websockets_v2_enabled": trueplaceholder,
placeholder
	svc := &OpenAIGatewayService{
		cfg:              cfg,
		rateLimitService: NewRateLimitService(transientCooldownAccountRepo{placeholder, nil, cfg, nil, nil),
		httpUpstream:     &httpUpstreamRecorder{placeholder,
		cache:            &stubGatewayCache{placeholder,
		openaiWSResolver: NewOpenAIWSProtocolResolver(cfg),
		toolCorrector:    NewCodexToolCorrector(),
placeholder
	body := []byte(`{"model":"gpt-5.5","stream":false,"input":"hello"placeholder`)

	for range 2 {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
		c.Request.Header.Set("User-Agent", "unit-test-agent/1.0")
		result, err := svc.Forward(context.Background(), c, &account, body)
	placeholder
		require.Nil(t, result)
placeholder

	require.True(t, svc.isOpenAIAccountModelRuntimeBlocked(&account, "gpt-5.5"))
placeholder

func TestOpenAIGatewayService_ProxyResponsesWebSocketFromClient_ErrorEventUsageLimitPersistsRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := newOpenAIWSV2TestConfig()
	cfg.Security.URLAllowlist.Enabled = false
	cfg.Security.URLAllowlist.AllowInsecureHTTP = true
	cfg.Gateway.OpenAIWS.MaxConnsPerAccount = 1
	cfg.Gateway.OpenAIWS.MinIdlePerAccount = 0
	cfg.Gateway.OpenAIWS.MaxIdlePerAccount = 1
	cfg.Gateway.OpenAIWS.QueueLimitPerConn = 8
	cfg.Gateway.OpenAIWS.DialTimeoutSeconds = 3
	cfg.Gateway.OpenAIWS.ReadTimeoutSeconds = 3
	cfg.Gateway.OpenAIWS.WriteTimeoutSeconds = 3

	resetAt := time.Now().Add(90 * time.Minute).Unix()
	captureConn := &openAIWSCaptureConn{
		events: [][]byte{
			[]byte(`{"type":"error","error":{"code":"rate_limit_exceeded","type":"usage_limit_reached","message":"The usage limit has been reached","resets_at":PLACEHOLDERplaceholderplaceholder`),
	placeholder,
placeholder
	captureConn.events[0] = []byte(strings.ReplaceAll(string(captureConn.events[0]), "PLACEHOLDER", strconv.FormatInt(resetAt, 10)))
	captureDialer := &openAIWSCaptureDialer{conn: captureConnplaceholder
	pool := newOpenAIWSConnPool(cfg)
	pool.setClientDialerForTest(captureDialer)

	account := Account{
		ID:          503,
		Name:        "openai-ingress-rate-limit",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
placeholder
			"api_key": "sk-test",
	placeholder,
		Extra: map[string]any{
			"responses_websockets_v2_enabled": true,
	placeholder,
placeholder
	repo := &openAIWSRateLimitSignalRepo{stubOpenAIAccountRepo: stubOpenAIAccountRepo{accounts: []Account{accountplaceholderplaceholderplaceholder
	rateSvc := &RateLimitService{accountRepo: repoplaceholder
	svc := &OpenAIGatewayService{
		accountRepo:      repo,
		rateLimitService: rateSvc,
		httpUpstream:     &httpUpstreamRecorder{placeholder,
		cache:            &stubGatewayCache{placeholder,
		cfg:              cfg,
		openaiWSResolver: NewOpenAIWSProtocolResolver(cfg),
		toolCorrector:    NewCodexToolCorrector(),
		openaiWSPool:     pool,
placeholder

	serverErrCh := make(chan error, 1)
	wsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := coderws.Accept(w, r, &coderws.AcceptOptions{CompressionMode: coderws.CompressionContextTakeoverplaceholder)
		if err != nil {
			serverErrCh <- err
			return
	placeholder
		defer func() { _ = conn.CloseNow() placeholder()

		rec := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(rec)
		req := r.Clone(r.Context())
		req.Header = req.Header.Clone()
		req.Header.Set("User-Agent", "unit-test-agent/1.0")
		ginCtx.Request = req

		readCtx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		msgType, firstMessage, readErr := conn.Read(readCtx)
		cancel()
		if readErr != nil {
			serverErrCh <- readErr
			return
	placeholder
		if msgType != coderws.MessageText && msgType != coderws.MessageBinary {
			serverErrCh <- io.ErrUnexpectedEOF
			return
	placeholder

		serverErrCh <- svc.ProxyResponsesWebSocketFromClient(r.Context(), ginCtx, conn, &account, "sk-test", firstMessage, nil)
placeholder))
	defer wsServer.Close()

	dialCtx, cancelDial := context.WithTimeout(context.Background(), 3*time.Second)
	clientConn, _, err := coderws.Dial(dialCtx, "ws"+strings.TrimPrefix(wsServer.URL, "http"), nil)
	cancelDial()
placeholder
	defer func() { _ = clientConn.CloseNow() placeholder()

	writeCtx, cancelWrite := context.WithTimeout(context.Background(), 3*time.Second)
	err = clientConn.Write(writeCtx, coderws.MessageText, []byte(`{"type":"response.create","model":"gpt-5.1","stream":falseplaceholder`))
	cancelWrite()
placeholder

	select {
	case serverErr := <-serverErrCh:
		require.Error(t, serverErr)
		var failoverErr *UpstreamFailoverError
		require.ErrorAs(t, serverErr, &failoverErr)
		require.Equal(t, http.StatusTooManyRequests, failoverErr.StatusCode)
		require.Len(t, repo.rateLimitCalls, 1)
		require.WithinDuration(t, time.Unix(resetAt, 0), repo.rateLimitCalls[0], 2*time.Second)
	case <-time.After(5 * time.Second):
		t.Fatal("等待 ingress websocket 结束超时")
placeholder
placeholder

func TestOpenAIGatewayService_UpdateCodexUsageSnapshot_ExhaustedSnapshotDoesNotSetRateLimit(t *testing.T) {
	repo := &openAICodexSnapshotAsyncRepo{
		updateExtraCh: make(chan map[string]any, 1),
		rateLimitCh:   make(chan time.Time, 1),
placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder
	snapshot := &OpenAICodexUsageSnapshot{
		PrimaryUsedPercent:         ptrFloat64WS(100),
		PrimaryResetAfterSeconds:   ptrIntWS(3600),
		PrimaryWindowMinutes:       ptrIntWS(10080),
		SecondaryUsedPercent:       ptrFloat64WS(12),
		SecondaryResetAfterSeconds: ptrIntWS(1200),
		SecondaryWindowMinutes:     ptrIntWS(300),
placeholder
	svc.updateCodexUsageSnapshot(context.Background(), 601, snapshot)

	select {
	case updates := <-repo.updateExtraCh:
		require.Equal(t, 100.0, updates["codex_7d_used_percent"])
	case <-time.After(2 * time.Second):
		t.Fatal("等待 codex 快照落库超时")
placeholder

	select {
	case resetAt := <-repo.rateLimitCh:
		t.Fatalf("不应因仅写入快照而生成运行时限流时间: %v", resetAt)
	case <-time.After(2 * time.Second):
placeholder
placeholder

func TestOpenAIGatewayService_UpdateCodexUsageSnapshot_NonExhaustedSnapshotDoesNotSetRateLimit(t *testing.T) {
	repo := &openAICodexSnapshotAsyncRepo{
		updateExtraCh: make(chan map[string]any, 1),
		rateLimitCh:   make(chan time.Time, 1),
placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder
	snapshot := &OpenAICodexUsageSnapshot{
		PrimaryUsedPercent:         ptrFloat64WS(94),
		PrimaryResetAfterSeconds:   ptrIntWS(3600),
		PrimaryWindowMinutes:       ptrIntWS(10080),
		SecondaryUsedPercent:       ptrFloat64WS(22),
		SecondaryResetAfterSeconds: ptrIntWS(1200),
		SecondaryWindowMinutes:     ptrIntWS(300),
placeholder
	svc.updateCodexUsageSnapshot(context.Background(), 602, snapshot)

	select {
	case <-repo.updateExtraCh:
	case <-time.After(2 * time.Second):
		t.Fatal("等待 codex 快照落库超时")
placeholder

	select {
	case resetAt := <-repo.rateLimitCh:
		t.Fatalf("不应写入运行时限流时间: %v", resetAt)
	case <-time.After(200 * time.Millisecond):
placeholder
placeholder

func TestOpenAIGatewayService_UpdateCodexUsageSnapshot_ThrottlesExtraWrites(t *testing.T) {
	repo := &openAICodexSnapshotAsyncRepo{
		updateExtraCh: make(chan map[string]any, 2),
placeholder
	svc := &OpenAIGatewayService{
		accountRepo:           repo,
		codexSnapshotThrottle: newAccountWriteThrottle(time.Hour),
placeholder
	snapshot := &OpenAICodexUsageSnapshot{
		PrimaryUsedPercent:         ptrFloat64WS(94),
		PrimaryResetAfterSeconds:   ptrIntWS(3600),
		PrimaryWindowMinutes:       ptrIntWS(10080),
		SecondaryUsedPercent:       ptrFloat64WS(22),
		SecondaryResetAfterSeconds: ptrIntWS(1200),
		SecondaryWindowMinutes:     ptrIntWS(300),
placeholder

	svc.updateCodexUsageSnapshot(context.Background(), 777, snapshot)
	svc.updateCodexUsageSnapshot(context.Background(), 777, snapshot)

	select {
	case <-repo.updateExtraCh:
	case <-time.After(2 * time.Second):
		t.Fatal("等待第一次 codex 快照落库超时")
placeholder

	select {
	case updates := <-repo.updateExtraCh:
		t.Fatalf("unexpected second codex snapshot write: %v", updates)
	case <-time.After(200 * time.Millisecond):
placeholder
placeholder

func ptrFloat64WS(v float64) *float64 { return &v placeholder
func ptrIntWS(v int) *int             { return &v placeholder

func TestOpenAIGatewayService_GetSchedulableAccount_ExhaustedCodexExtraDoesNotSetRateLimit(t *testing.T) {
	resetAt := time.Now().Add(6 * 24 * time.Hour)
	account := Account{
		ID:          701,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
		Extra: map[string]any{
			"codex_7d_used_percent": 100.0,
			"codex_7d_reset_at":     resetAt.UTC().Format(time.RFC3339),
	placeholder,
placeholder
	repo := &openAICodexExtraListRepo{stubOpenAIAccountRepo: stubOpenAIAccountRepo{accounts: []Account{accountplaceholderplaceholder, rateLimitCh: make(chan time.Time, 1)placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder

	fresh, err := svc.getSchedulableAccount(context.Background(), account.ID)
placeholder
	require.NotNil(t, fresh)
	require.Nil(t, fresh.RateLimitResetAt)
	select {
	case persisted := <-repo.rateLimitCh:
		t.Fatalf("不应将已耗尽的 codex extra 提升为运行时限流状态: %v", persisted)
	case <-time.After(2 * time.Second):
placeholder
placeholder

func TestAdminService_ListAccounts_ExhaustedCodexExtraDoesNotSetRateLimit(t *testing.T) {
	resetAt := time.Now().Add(4 * 24 * time.Hour)
	repo := &openAICodexExtraListRepo{
		stubOpenAIAccountRepo: stubOpenAIAccountRepo{accounts: []Account{{
			ID:          702,
			Platform:    PlatformOpenAI,
			Type:        AccountTypeOAuth,
			Status:      StatusActive,
			Schedulable: true,
			Concurrency: 1,
			Extra: map[string]any{
				"codex_7d_used_percent": 100.0,
				"codex_7d_reset_at":     resetAt.UTC().Format(time.RFC3339),
		placeholder,
	placeholderplaceholderplaceholder,
		rateLimitCh: make(chan time.Time, 1),
placeholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder

	accounts, total, err := svc.ListAccounts(context.Background(), 1, 20, PlatformOpenAI, AccountTypeOAuth, "", "", 0, "", "", "")
placeholder
	require.Equal(t, int64(1), total)
	require.Len(t, accounts, 1)
	require.Nil(t, accounts[0].RateLimitResetAt)
	select {
	case persisted := <-repo.rateLimitCh:
		t.Fatalf("不应在账号列表查询时将 codex extra 持久化为运行时限流状态: %v", persisted)
	case <-time.After(2 * time.Second):
placeholder
placeholder

func TestOpenAIWSErrorHTTPStatusFromRaw_UsageLimitReachedIs429(t *testing.T) {
	require.Equal(t, http.StatusTooManyRequests, openAIWSErrorHTTPStatusFromRaw("", "usage_limit_reached"))
	require.Equal(t, http.StatusTooManyRequests, openAIWSErrorHTTPStatusFromRaw("rate_limit_exceeded", ""))
placeholder

func TestOpenAIWSRateLimitFailoverError_OAuthKeepsSameAccountDeadline(t *testing.T) {
	svc := &OpenAIGatewayService{placeholder
	headers := http.Header{"Retry-After": []string{"30"placeholderplaceholder
	body := []byte(`{"error":{"type":"rate_limit_error","message":"limited"placeholderplaceholder`)

	oauthErr := svc.newOpenAIWSRateLimitFailoverError(&Account{
		ID:       904,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder, headers, body, "limited")
	require.True(t, oauthErr.RetryableOnSameAccount)
	require.False(t, oauthErr.SameAccountRetryDeadline.IsZero())
	require.Positive(t, oauthErr.SameAccountRetryDelay)
	require.LessOrEqual(t, oauthErr.SameAccountRetryDelay, openAIOAuth429MaxRetryDelay)
	require.Equal(t, body, oauthErr.ResponseBody)
	require.Equal(t, "30", oauthErr.ResponseHeaders.Get("Retry-After"))

	apiKeyErr := svc.newOpenAIWSRateLimitFailoverError(&Account{
		ID:       905,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
placeholder, headers, body, "limited")
	require.False(t, apiKeyErr.RetryableOnSameAccount)
	require.True(t, apiKeyErr.SameAccountRetryDeadline.IsZero())
	require.Zero(t, apiKeyErr.SameAccountRetryDelay)
placeholder
