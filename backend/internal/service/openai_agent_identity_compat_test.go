package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestAccountTestServiceOpenAICompactAgentIdentityUsesFreshAssertion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	key, privateKey := newTestAgentIdentityKey(t)
	account := Account{
		ID:          21,
		Name:        "agent-identity",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
placeholder
			"auth_mode":                  OpenAIAuthModeAgentIdentity,
			"agent_runtime_id":           key.runtimeID,
			"agent_private_key":          privateKey,
			"task_id":                    key.taskID,
			"chatgpt_account_id":         "account-agent-test",
			"chatgpt_account_is_fedramp": true,
	placeholder,
placeholder
	repo := &snapshotUpdateAccountRepo{stubOpenAIAccountRepo: stubOpenAIAccountRepo{accounts: []Account{accountplaceholderplaceholderplaceholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"id":"compact-agent","status":"completed"placeholder`)),
placeholderplaceholder
	svc := &AccountTestService{accountRepo: repo, httpUpstream: upstreamplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/accounts/21/test", bytes.NewReader(nil))

	require.NoError(t, svc.TestAccountConnection(c, account.ID, "gpt-5.4", "", AccountTestModeCompact))
	require.Equal(t, "AgentAssertion", strings.SplitN(upstream.lastReq.Header.Get("Authorization"), " ", 2)[0])
	require.Equal(t, "account-agent-test", upstream.lastReq.Header.Get("chatgpt-account-id"))
	require.Equal(t, "true", upstream.lastReq.Header.Get("x-openai-fedramp"))
	require.NotContains(t, upstream.lastReq.Header.Get("Authorization"), privateKey)
placeholder

func TestAccountTestServiceOpenAICompactAgentIdentityRecoversInvalidTaskOnce(t *testing.T) {
	gin.SetMode(gin.TestMode)
	key, privateKey := newTestAgentIdentityKey(t)
	account := &Account{
		ID:          22,
		Name:        "agent-identity-recovery",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
placeholder
			"auth_mode":          OpenAIAuthModeAgentIdentity,
			"agent_runtime_id":   key.runtimeID,
			"agent_private_key":  privateKey,
			"task_id":            "task-compact-old",
			"chatgpt_account_id": "account-agent-compact-recovery",
	placeholder,
placeholder
	repo := &accountTestAgentIdentityRepo{account: accountplaceholder
	registerCalls := 0
	registerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		registerCalls++
		_, _ = io.WriteString(w, `{"task_id":"task-compact-new"placeholder`)
placeholder))
	defer registerServer.Close()
	oldBase := openAIAgentIdentityAuthAPIBaseURL
	openAIAgentIdentityAuthAPIBaseURL = registerServer.URL
	t.Cleanup(func() { openAIAgentIdentityAuthAPIBaseURL = oldBase placeholder)

	upstream := &httpUpstreamRecorder{responses: []*http.Response{
		{StatusCode: http.StatusUnauthorized, Header: http.Header{"Content-Type": []string{"application/json"placeholderplaceholder, Body: io.NopCloser(strings.NewReader(`{"error":{"code":"invalid_task_id"placeholderplaceholder`))placeholder,
		{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": []string{"application/json"placeholderplaceholder, Body: io.NopCloser(strings.NewReader(`{"id":"compact-agent","status":"completed"placeholder`))placeholder,
placeholderplaceholder
	invalidator := &agentIdentityWSInvalidationRecorder{placeholder
	svc := &AccountTestService{accountRepo: repo, httpUpstream: upstream, agentIdentityWS: invalidatorplaceholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/accounts/22/test", bytes.NewReader(nil))

	require.NoError(t, svc.TestAccountConnection(c, account.ID, "gpt-5.4", "", AccountTestModeCompact))
	require.Equal(t, 1, registerCalls)
	require.Len(t, upstream.requests, 2)
	require.Equal(t, "task-compact-new", account.GetCredential("task_id"))
	require.Equal(t, 0, repo.setErrorCalls)
	require.Equal(t, []int64{account.IDplaceholder, invalidator.accountIDs)
placeholder

func TestOpenAIAgentIdentityPassthroughKeepsSessionAndPromptCacheHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	key, privateKey := newTestAgentIdentityKey(t)
	account := &Account{
		ID:       24,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"auth_mode":          OpenAIAuthModeAgentIdentity,
			"agent_runtime_id":   key.runtimeID,
			"agent_private_key":  privateKey,
			"task_id":            key.taskID,
			"chatgpt_account_id": "account-agent-passthrough",
	placeholder,
placeholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	body := []byte(`{"model":"gpt-5.4","instructions":"Reply OK","input":[],"stream":true,"prompt_cache_key":"cache-agent"placeholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))
	c.Request.Header.Set("session_id", "client-session")
	c.Request.Header.Set("conversation_id", "client-conversation")
	c.Request.Header.Set("Authorization", "Bearer inbound-must-not-forward")

	svc := &OpenAIGatewayService{placeholder
	req, err := svc.buildUpstreamRequestOpenAIPassthrough(context.Background(), c, account, body, "")
placeholder
	require.Equal(t, "AgentAssertion", strings.SplitN(req.Header.Get("Authorization"), " ", 2)[0])
	require.Equal(t, "account-agent-passthrough", req.Header.Get("chatgpt-account-id"))
	require.NotEqual(t, "client-session", req.Header.Get("session_id"))
	require.NotEqual(t, "client-conversation", req.Header.Get("conversation_id"))
	require.Equal(t, isolateOpenAISessionID(0, "client-session"), req.Header.Get("session_id"))
	require.Equal(t, isolateOpenAISessionID(0, "client-conversation"), req.Header.Get("conversation_id"))
	requestBody, err := io.ReadAll(req.Body)
placeholder
	require.Contains(t, string(requestBody), `"prompt_cache_key":"cache-agent"`)

	// Authentication mode must not affect session isolation or prompt-cache
	// behavior. Compare the same request with the existing OAuth path instead
	// of pinning this test to an implementation-specific hash.
	oauthAccount := &Account{
		ID:       26,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"chatgpt_account_id": "account-oauth-passthrough",
	placeholder,
placeholder
	oauthRecorder := httptest.NewRecorder()
	oauthContext, _ := gin.CreateTestContext(oauthRecorder)
	oauthContext.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))
	oauthContext.Request.Header.Set("session_id", "client-session")
	oauthContext.Request.Header.Set("conversation_id", "client-conversation")
	oauthReq, err := svc.buildUpstreamRequestOpenAIPassthrough(context.Background(), oauthContext, oauthAccount, body, "oauth-token")
placeholder
	require.Equal(t, oauthReq.Header.Get("session_id"), req.Header.Get("session_id"))
	require.Equal(t, oauthReq.Header.Get("conversation_id"), req.Header.Get("conversation_id"))
placeholder

func TestOpenAIAgentIdentityErrorRedactionDoesNotLeakCredentialValues(t *testing.T) {
	key, privateKey := newTestAgentIdentityKey(t)
	account := &Account{
		ID:       25,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"auth_mode":         OpenAIAuthModeAgentIdentity,
			"agent_runtime_id":  key.runtimeID,
			"agent_private_key": privateKey,
			"task_id":           key.taskID,
			"access_token":      key.runtimeID + "-oauth-value",
	placeholder,
placeholder
	svc := &OpenAIGatewayService{placeholder
	oauthValue := account.GetCredential("access_token")
	redacted := svc.redactAgentIdentitySensitiveBody(context.Background(), account, []byte(`{"message":"runtime-test task-test `+oauthValue+` AgentAssertion abc123"placeholder`))
	require.NotContains(t, string(redacted), key.runtimeID)
	require.NotContains(t, string(redacted), key.taskID)
	require.NotContains(t, string(redacted), oauthValue)
	require.NotContains(t, string(redacted), "AgentAssertion abc123")
	require.Contains(t, string(redacted), "[redacted]")
placeholder

func TestOpenAIAuthenticationHeadersPreserveOAuthPATAndAPIKeyBearerModes(t *testing.T) {
	svc := &OpenAIGatewayService{placeholder
	tests := []struct {
		name    string
		account *Account
		token   string
placeholder{
		{name: "oauth", account: &Account{Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder, token: "oauth-runtime-token"placeholder,
		{name: "personal access token", account: &Account{Platform: PlatformOpenAI, Type: AccountTypeOAuth, Credentials: map[string]any{"auth_mode": OpenAIAuthModePersonalAccessTokenplaceholderplaceholder, token: "pat-runtime-token"placeholder,
		{name: "api key", account: &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder, token: "api-key-runtime-token"placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers, err := svc.buildOpenAIAuthenticationHeaders(context.Background(), tt.account, tt.token)
		placeholder
			require.Equal(t, "Bearer "+tt.token, headers.Get("Authorization"))
	placeholder)
placeholder
placeholder

func TestOpenAIWSAgentIdentityRecoveryRequiresTaskInvalidBody(t *testing.T) {
	require.False(t, isAgentIdentityTaskInvalidWSDialError(&openAIWSDialError{
		StatusCode:   http.StatusUnauthorized,
		ResponseBody: []byte(`{"error":{"code":"invalid_signature"placeholderplaceholder`),
placeholder))
	require.True(t, isAgentIdentityTaskInvalidWSDialError(&openAIWSDialError{
		StatusCode:   http.StatusUnauthorized,
		ResponseBody: []byte(`{"error":{"code":"invalid_task_id"placeholderplaceholder`),
placeholder))
placeholder

func TestValidateOpenAIWSBearerTokenAllowsAgentIdentityWithoutStoredToken(t *testing.T) {
	t.Run("Given Agent Identity When a WS path receives no bearer token Then dial-time assertion auth is allowed", func(t *testing.T) {
		account := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeOAuth,
	placeholder
				"auth_mode": OpenAIAuthModeAgentIdentity,
		placeholder,
	placeholder

		require.NoError(t, validateOpenAIWSBearerToken(account, ""))
placeholder)

	t.Run("Given bearer credentials When a WS path receives no token Then the request is rejected", func(t *testing.T) {
		accounts := []*Account{
			{Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder,
			{Platform: PlatformOpenAI, Type: AccountTypeOAuth, Credentials: map[string]any{"auth_mode": OpenAIAuthModePersonalAccessTokenplaceholderplaceholder,
			{Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder,
	placeholder

		for _, account := range accounts {
			require.EqualError(t, validateOpenAIWSBearerToken(account, ""), "token is empty")
	placeholder
placeholder)
placeholder

func TestOpenAIWSConnPoolHeadersFactoryRunsAtDialAndStalePrewarmIsDiscarded(t *testing.T) {
	cfg := &config.Config{placeholder
	cfg.Gateway.OpenAIWS.MaxConnsPerAccount = 1
	pool := newOpenAIWSConnPool(cfg)
	defer pool.Close()
	pool.setClientDialerForTest(&openAIWSFakeDialer{placeholder)

	accountID := int64(22)
	ap := pool.getOrCreateAccountPool(accountID)
	factoryCalls := 0
	latestHeader := ""
	req := openAIWSAcquireRequest{
		Account: &Account{ID: accountID, Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder,
		WSURL:   "wss://example.com/v1/responses",
		HeadersFactory: func(_ context.Context, headers http.Header) (http.Header, error) {
			factoryCalls++
			latestHeader = "AgentAssertion dial-" + string(rune('0'+factoryCalls))
			if headers == nil {
				headers = make(http.Header)
		placeholder
			headers.Set("Authorization", latestHeader)
			return headers, nil
	placeholder,
placeholder
	ap.mu.Lock()
	ap.lastAcquire = &req
	generation := ap.generation
	ap.mu.Unlock()

	pool.prewarmConns(accountID, req, 1, generation)
	require.Equal(t, 1, factoryCalls, "prewarm must generate authorization inside the actual dial")
	require.Equal(t, "AgentAssertion dial-1", latestHeader)

	pool.ClearAccount(accountID)
	ap.mu.Lock()
	require.Empty(t, ap.conns, "credential recovery must remove pooled connections")
	require.Nil(t, ap.lastAcquire, "credential recovery must discard delayed acquire state")
	require.Equal(t, generation+1, ap.generation)
	ap.mu.Unlock()

	// A prewarm captured before ClearAccount must not be admitted after recovery.
	pool.prewarmConns(accountID, req, 1, generation)
	ap.mu.Lock()
	require.Empty(t, ap.conns)
	ap.mu.Unlock()
placeholder

func TestOpenAIAgentIdentityTaskInvalidRetriesExactlyOnce(t *testing.T) {
	gin.SetMode(gin.TestMode)
	key, privateKey := newTestAgentIdentityKey(t)
	account := &Account{
		ID:          23,
		Name:        "agent-identity",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
placeholder
			"auth_mode":          OpenAIAuthModeAgentIdentity,
			"agent_runtime_id":   key.runtimeID,
			"agent_private_key":  privateKey,
			"task_id":            "task-old",
			"chatgpt_account_id": "account-agent-retry",
	placeholder,
placeholder
	repo := &agentIdentityForwardRepo{account: accountplaceholder
	registerCalls := 0
	registerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		registerCalls++
		_, _ = io.WriteString(w, `{"task_id":"task-new"placeholder`)
placeholder))
	defer registerServer.Close()
	oldBase := openAIAgentIdentityAuthAPIBaseURL
	openAIAgentIdentityAuthAPIBaseURL = registerServer.URL
	t.Cleanup(func() { openAIAgentIdentityAuthAPIBaseURL = oldBase placeholder)

	successBody := `{"id":"resp-agent-retry","object":"response","model":"gpt-5.4","output":[],"usage":{"input_tokens":1,"output_tokens":1,"total_tokens":2placeholderplaceholder`
	upstream := &httpUpstreamRecorder{responses: []*http.Response{
		{StatusCode: http.StatusUnauthorized, Header: http.Header{"Content-Type": []string{"application/json"placeholderplaceholder, Body: io.NopCloser(strings.NewReader(`{"error":{"code":"invalid_task_id"placeholderplaceholder`))placeholder,
		{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": []string{"application/json"placeholderplaceholder, Body: io.NopCloser(strings.NewReader(successBody))placeholder,
placeholderplaceholder
	require.True(t, isAgentIdentityTaskInvalidHTTPResponse(http.StatusUnauthorized, []byte(`{"error":{"code":"invalid_task_id"placeholderplaceholder`)))
	svc := &OpenAIGatewayService{cfg: &config.Config{placeholder, accountRepo: repo, httpUpstream: upstreamplaceholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", strings.NewReader(`{"model":"gpt-5.4","instructions":"Reply OK","input":[],"stream":falseplaceholder`))

	_, err := svc.Forward(context.Background(), c, account, []byte(`{"model":"gpt-5.4","instructions":"Reply OK","input":[],"stream":falseplaceholder`))
placeholder
	require.Equal(t, 1, registerCalls)
	require.Len(t, upstream.requests, 2)
	require.NotEqual(t, upstream.requests[0].Header.Get("Authorization"), upstream.requests[1].Header.Get("Authorization"))
	require.Equal(t, "task-new", decodeAgentAssertionTask(t, upstream.requests[1].Header.Get("Authorization")))

	// Two consecutive invalid responses still produce only one retry for this
	// request; the recovery path must not loop indefinitely.
	upstream.responses = []*http.Response{
		{StatusCode: http.StatusUnauthorized, Header: http.Header{"Content-Type": []string{"application/json"placeholderplaceholder, Body: io.NopCloser(strings.NewReader(`{"error":{"code":"invalid_task_id"placeholderplaceholder`))placeholder,
		{StatusCode: http.StatusUnauthorized, Header: http.Header{"Content-Type": []string{"application/json"placeholderplaceholder, Body: io.NopCloser(strings.NewReader(`{"error":{"code":"invalid_task_id"placeholderplaceholder`))placeholder,
placeholder
	rec2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(rec2)
	c2.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", strings.NewReader(`{"model":"gpt-5.4","instructions":"Reply OK","input":[],"stream":falseplaceholder`))
	_, err = svc.Forward(context.Background(), c2, account, []byte(`{"model":"gpt-5.4","instructions":"Reply OK","input":[],"stream":falseplaceholder`))
placeholder
	require.Equal(t, 2, registerCalls)
	require.Len(t, upstream.requests, 4)

	// Passthrough uses the same one-shot task recovery contract.
	account.Extra = map[string]any{"openai_passthrough": trueplaceholder
	account.Credentials["task_id"] = "task-old-passthrough"
	upstream.responses = []*http.Response{
		{StatusCode: http.StatusUnauthorized, Header: http.Header{"Content-Type": []string{"application/json"placeholderplaceholder, Body: io.NopCloser(strings.NewReader(`{"error":{"code":"invalid_task_id"placeholderplaceholder`))placeholder,
		{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder, Body: io.NopCloser(strings.NewReader("data: {\"type\":\"response.completed\",\"response\":{\"usage\":{\"input_tokens\":1,\"output_tokens\":1,\"total_tokens\":2placeholderplaceholderplaceholder\n\ndata: [DONE]\n\n"))placeholder,
placeholder
	rec3 := httptest.NewRecorder()
	c3, _ := gin.CreateTestContext(rec3)
	c3.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", strings.NewReader(`{"model":"gpt-5.4","instructions":"Reply OK","input":[],"stream":falseplaceholder`))
	_, err = svc.Forward(context.Background(), c3, account, []byte(`{"model":"gpt-5.4","instructions":"Reply OK","input":[],"stream":falseplaceholder`))
placeholder
	require.Equal(t, 3, registerCalls)
	require.Len(t, upstream.requests, 6)
placeholder

func TestOpenAIAgentIdentityCompatRoutesRecoverInvalidTaskOnce(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name string
		path string
		body []byte
		call func(*OpenAIGatewayService, context.Context, *gin.Context, *Account, []byte) (*OpenAIForwardResult, error)
placeholder{
		{
			name: "chat completions",
			path: "/v1/chat/completions",
			body: []byte(`{"model":"gpt-5.4","stream":false,"messages":[{"role":"user","content":"hi"placeholder]placeholder`),
			call: func(s *OpenAIGatewayService, ctx context.Context, c *gin.Context, account *Account, body []byte) (*OpenAIForwardResult, error) {
				return s.ForwardAsChatCompletions(ctx, c, account, body, "", "gpt-5.4")
		placeholder,
	placeholder,
		{
			name: "anthropic messages",
			path: "/v1/messages",
			body: []byte(`{"model":"gpt-5.4","stream":false,"max_tokens":32,"messages":[{"role":"user","content":"hi"placeholder]placeholder`),
			call: func(s *OpenAIGatewayService, ctx context.Context, c *gin.Context, account *Account, body []byte) (*OpenAIForwardResult, error) {
				return s.ForwardAsAnthropic(ctx, c, account, body, "", "gpt-5.4")
		placeholder,
	placeholder,
placeholder

	for index, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, privateKey := newTestAgentIdentityKey(t)
			account := &Account{
				ID:          int64(40 + index),
				Name:        "agent-identity-compat",
				Platform:    PlatformOpenAI,
				Type:        AccountTypeOAuth,
				Status:      StatusActive,
				Schedulable: true,
				Concurrency: 1,
		placeholder
					"auth_mode":          OpenAIAuthModeAgentIdentity,
					"agent_runtime_id":   key.runtimeID,
					"agent_private_key":  privateKey,
					"task_id":            "task-compat-old",
					"chatgpt_account_id": "account-compat-recovery",
			placeholder,
		placeholder
			repo := &agentIdentityForwardRepo{account: accountplaceholder
			registerCalls := 0
			registerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				registerCalls++
				_, _ = io.WriteString(w, `{"task_id":"task-compat-new"placeholder`)
		placeholder))
			defer registerServer.Close()
			oldBase := openAIAgentIdentityAuthAPIBaseURL
			openAIAgentIdentityAuthAPIBaseURL = registerServer.URL
			t.Cleanup(func() { openAIAgentIdentityAuthAPIBaseURL = oldBase placeholder)

			upstream := &httpUpstreamRecorder{responses: []*http.Response{
				{StatusCode: http.StatusUnauthorized, Header: http.Header{"Content-Type": []string{"application/json"placeholderplaceholder, Body: io.NopCloser(strings.NewReader(`{"error":{"code":"invalid_task_id"placeholderplaceholder`))placeholder,
				{StatusCode: http.StatusUnauthorized, Header: http.Header{"Content-Type": []string{"application/json"placeholderplaceholder, Body: io.NopCloser(strings.NewReader(`{"error":{"code":"invalid_task_id"placeholderplaceholder`))placeholder,
		placeholderplaceholder
			svc := &OpenAIGatewayService{cfg: &config.Config{placeholder, accountRepo: repo, httpUpstream: upstreamplaceholder
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodPost, tt.path, bytes.NewReader(tt.body))

			_, err := tt.call(svc, context.Background(), c, account, tt.body)
		placeholder
			require.Equal(t, 1, registerCalls)
			require.Len(t, upstream.requests, 2)
			require.Equal(t, "task-compat-new", account.GetCredential("task_id"))
	placeholder)
placeholder
placeholder

func decodeAgentAssertionTask(t *testing.T, header string) string {
placeholder
	encoded := strings.TrimPrefix(header, "AgentAssertion ")
	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
placeholder
	var envelope struct {
		TaskID string `json:"task_id"`
placeholder
	require.NoError(t, json.Unmarshal(decoded, &envelope))
	return envelope.TaskID
placeholder

type agentIdentityForwardRepo struct {
	AccountRepository
	account *Account
placeholder

type agentIdentityWSInvalidationRecorder struct {
	accountIDs []int64
placeholder

func (r *agentIdentityWSInvalidationRecorder) InvalidateAgentIdentityWSConnections(accountID int64) {
	r.accountIDs = append(r.accountIDs, accountID)
placeholder

type accountTestAgentIdentityRepo struct {
	AccountRepository
	account       *Account
	setErrorCalls int
placeholder

func (r *accountTestAgentIdentityRepo) GetByID(_ context.Context, _ int64) (*Account, error) {
	return r.account, nil
placeholder

func (r *accountTestAgentIdentityRepo) UpdateCredentials(_ context.Context, _ int64, credentials map[string]any) error {
	r.account.Credentials = credentials
	return nil
placeholder

func (r *accountTestAgentIdentityRepo) UpdateExtra(_ context.Context, _ int64, _ map[string]any) error {
	return nil
placeholder

func (r *accountTestAgentIdentityRepo) SetError(_ context.Context, _ int64, _ string) error {
	r.setErrorCalls++
	return nil
placeholder

func (r *agentIdentityForwardRepo) GetByID(_ context.Context, _ int64) (*Account, error) {
	return r.account, nil
placeholder

func (r *agentIdentityForwardRepo) UpdateCredentials(_ context.Context, _ int64, credentials map[string]any) error {
	r.account.Credentials = credentials
	return nil
placeholder
