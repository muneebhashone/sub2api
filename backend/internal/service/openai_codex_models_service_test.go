package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
)

type codexModelsHTTPUpstreamStub struct {
	do func(req *http.Request, proxyURL string, accountID int64, accountConcurrency int) (*http.Response, error)
placeholder

type codexModelsBlockingBody struct {
	ctx         context.Context
	readStarted chan struct{placeholder
	startedOnce *sync.Once
	release     <-chan struct{placeholder
	body        *strings.Reader
placeholder

func (b *codexModelsBlockingBody) Read(p []byte) (int, error) {
	b.startedOnce.Do(func() { close(b.readStarted) placeholder)
	select {
	case <-b.release:
		return b.body.Read(p)
	case <-b.ctx.Done():
		return 0, b.ctx.Err()
placeholder
placeholder

func (b *codexModelsBlockingBody) Close() error { return nil placeholder

func (s *codexModelsHTTPUpstreamStub) Do(req *http.Request, proxyURL string, accountID int64, accountConcurrency int) (*http.Response, error) {
	return s.do(req, proxyURL, accountID, accountConcurrency)
placeholder

func (s *codexModelsHTTPUpstreamStub) DoWithTLS(req *http.Request, proxyURL string, accountID int64, accountConcurrency int, _ *tlsfingerprint.Profile) (*http.Response, error) {
	return s.Do(req, proxyURL, accountID, accountConcurrency)
placeholder

func TestIsRetryableCodexModelsManifestTransportError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		retryable bool
placeholder{
		{name: "nil", err: nilplaceholder,
		{name: "configuration error", err: errors.New("invalid proxy URL")placeholder,
		{name: "upstream configuration error", err: errors.New("upstream error: invalid proxy")placeholder,
		{name: "proxy connection configuration error", err: errors.New("proxy connection error: invalid configuration")placeholder,
		{name: "canceled request", err: context.Canceledplaceholder,
		{
			name: "redirect policy error",
			err: &url.Error{
				Op:  "Get",
				URL: "https://upstream.example/v1/models",
				Err: errors.New("stopped after 10 redirects"),
		placeholder,
	placeholder,
		{name: "deadline exceeded", err: context.DeadlineExceeded, retryable: trueplaceholder,
		{name: "unexpected EOF", err: io.ErrUnexpectedEOF, retryable: trueplaceholder,
		{name: "closed connection", err: net.ErrClosed, retryable: trueplaceholder,
		{
			name: "network operation",
			err: &net.OpError{
				Op:  "read",
				Net: "tcp",
				Err: errors.New("connection reset"),
		placeholder,
			retryable: true,
	placeholder,
		{
			name:      "DNS error",
			err:       &net.DNSError{Err: "temporary failure", Name: "upstream.example"placeholder,
			retryable: true,
	placeholder,
		{
			name:      "typed HTTP2 GOAWAY",
			err:       http2.GoAwayError{ErrCode: http2.ErrCodeNoplaceholder,
			retryable: true,
	placeholder,
		{
			name:      "stdlib HTTP2 GOAWAY",
			err:       errors.New("http2: server sent GOAWAY and closed the connection; LastStreamID=1, ErrCode=NO_ERROR"),
			retryable: true,
	placeholder,
		{
			name:      "stdlib HTTP2 refused stream",
			err:       errors.New("stream error: stream ID 3; REFUSED_STREAM"),
			retryable: true,
	placeholder,
		{
			name:      "stdlib HTTP2 connection error",
			err:       errors.New(`Get "https://upstream.example/v1/models": connection error: PROTOCOL_ERROR`),
			retryable: true,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRetryableCodexModelsManifestTransportError(tt.err); got != tt.retryable {
				t.Fatalf("retryable = %v, want %v", got, tt.retryable)
		placeholder
	placeholder)
placeholder
placeholder

func newCodexModelsAPIKeyTestService(upstream HTTPUpstream) *OpenAIGatewayService {
	return &OpenAIGatewayService{
		cfg: &config.Config{Security: config.SecurityConfig{URLAllowlist: config.URLAllowlistConfig{
			Enabled: false,
	placeholderplaceholderplaceholder,
		httpUpstream: upstream,
placeholder
placeholder

func newCodexModelsAPIKeyTestAccount(baseURL string) *Account {
	credentials := map[string]any{"api_key": "sk-upstream"placeholder
	if baseURL != "" {
		credentials["base_url"] = baseURL
placeholder
placeholder
		ID:          2,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Credentials: credentials,
		Concurrency: 3,
placeholder
placeholder

func newCodexModelsTestAccount() *Account {
placeholder
		ID:       1,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token":       "test-access-token",
			"chatgpt_account_id": "acc-123",
	placeholder,
placeholder
placeholder

func TestFetchCodexModelsManifestPassthrough(t *testing.T) {
	manifestBody := `{"models":[{"slug":"gpt-5.5","display_name":"GPT-5.5"placeholder]placeholder`

	var gotAuth, gotAccountID, gotOriginator, gotClientVersion string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotAccountID = r.Header.Get("chatgpt-account-id")
		gotOriginator = r.Header.Get("Originator")
		gotClientVersion = r.URL.Query().Get("client_version")
		w.Header().Set("ETag", `W/"abc123"`)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(manifestBody))
placeholder))
	defer server.Close()

	original := chatgptCodexModelsURL
	chatgptCodexModelsURL = server.URL
	defer func() { chatgptCodexModelsURL = original placeholder()

	s := &OpenAIGatewayService{placeholder
	manifest, err := s.FetchCodexModelsManifest(context.Background(), newCodexModelsTestAccount(), "0.137.0", "")
	if err != nil {
		t.Fatalf("FetchCodexModelsManifest returned error: %v", err)
placeholder

	if string(manifest.Body) != manifestBody {
		t.Errorf("body not passed through verbatim: got %q", manifest.Body)
placeholder
	if manifest.ETag != `W/"abc123"` {
		t.Errorf("etag not passed through: got %q", manifest.ETag)
placeholder
	if gotAuth != "Bearer test-access-token" {
		t.Errorf("authorization header: got %q", gotAuth)
placeholder
	if gotAccountID != "acc-123" {
		t.Errorf("chatgpt-account-id header: got %q", gotAccountID)
placeholder
	if gotOriginator != openai.CodexDefaultOriginator {
		t.Errorf("originator header: got %q", gotOriginator)
placeholder
	if gotClientVersion != "0.137.0" {
		t.Errorf("client_version query: got %q", gotClientVersion)
placeholder
placeholder

func TestFetchCodexModelsManifestAgentIdentityUsesAssertionWithoutOAuthToken(t *testing.T) {
	key, privateKey := newTestAgentIdentityKey(t)
	account := &Account{
		ID:       3,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"auth_mode":          OpenAIAuthModeAgentIdentity,
			"agent_runtime_id":   key.runtimeID,
			"agent_private_key":  privateKey,
			"task_id":            key.taskID,
			"chatgpt_account_id": "acc-agent",
	placeholder,
placeholder

	var gotAuth, gotAccountID string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotAccountID = r.Header.Get("chatgpt-account-id")
		_, _ = w.Write([]byte(`{"models":[]placeholder`))
placeholder))
	defer server.Close()

	original := chatgptCodexModelsURL
	chatgptCodexModelsURL = server.URL
	defer func() { chatgptCodexModelsURL = original placeholder()

	s := &OpenAIGatewayService{placeholder
	manifest, err := s.FetchCodexModelsManifest(context.Background(), account, "0.137.0", "")
	if err != nil {
		t.Fatalf("FetchCodexModelsManifest returned error: %v", err)
placeholder
	if string(manifest.Body) != `{"models":[]placeholder` {
		t.Fatalf("unexpected manifest body: %q", manifest.Body)
placeholder
	if !strings.HasPrefix(gotAuth, "AgentAssertion ") {
		t.Fatalf("authorization scheme: got %q", strings.SplitN(gotAuth, " ", 2)[0])
placeholder
	if gotAccountID != "acc-agent" {
		t.Fatalf("chatgpt-account-id header: got %q", gotAccountID)
placeholder
placeholder

func TestFetchCodexModelsManifestAgentIdentityRecoversInvalidTaskOnce(t *testing.T) {
	key, privateKey := newTestAgentIdentityKey(t)
	account := &Account{
		ID:       4,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"auth_mode":          OpenAIAuthModeAgentIdentity,
			"agent_runtime_id":   key.runtimeID,
			"agent_private_key":  privateKey,
			"task_id":            "task-models-old",
			"chatgpt_account_id": "acc-agent-recovery",
	placeholder,
placeholder
	repo := &stubQuotaAccountRepo{accounts: map[int64]*Account{account.ID: accountplaceholderplaceholder
	modelsCalls := 0
	registerCalls := 0
	var assertions []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		if strings.Contains(r.URL.Path, "/task/register") {
			registerCalls++
			_, _ = w.Write([]byte(`{"task_id":"task-models-new"placeholder`))
			return
	placeholder
		modelsCalls++
		assertions = append(assertions, r.Header.Get("Authorization"))
		if modelsCalls == 1 {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":{"code":"invalid_task_id"placeholderplaceholder`))
			return
	placeholder
		_, _ = w.Write([]byte(`{"models":[]placeholder`))
placeholder))
	defer server.Close()

	originalModelsURL := chatgptCodexModelsURL
	chatgptCodexModelsURL = server.URL
	t.Cleanup(func() { chatgptCodexModelsURL = originalModelsURL placeholder)
	originalAuthBase := openAIAgentIdentityAuthAPIBaseURL
	openAIAgentIdentityAuthAPIBaseURL = server.URL
	t.Cleanup(func() { openAIAgentIdentityAuthAPIBaseURL = originalAuthBase placeholder)

	s := &OpenAIGatewayService{accountRepo: repoplaceholder
	manifest, err := s.FetchCodexModelsManifest(context.Background(), account, "0.137.0", "")
placeholder
	require.Equal(t, `{"models":[]placeholder`, string(manifest.Body))
	require.Equal(t, 2, modelsCalls)
	require.Equal(t, 1, registerCalls)
	require.Len(t, assertions, 2)
	require.Equal(t, "task-models-old", decodeAgentAssertionTask(t, assertions[0]))
	require.Equal(t, "task-models-new", decodeAgentAssertionTask(t, assertions[1]))
placeholder

func TestFetchCodexModelsManifestAgentIdentityRedactsUpstreamErrors(t *testing.T) {
	key, privateKey := newTestAgentIdentityKey(t)
	account := &Account{
		ID:       5,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"auth_mode":          OpenAIAuthModeAgentIdentity,
			"agent_runtime_id":   key.runtimeID,
			"agent_private_key":  privateKey,
			"task_id":            key.taskID,
			"chatgpt_account_id": "acc-agent-redaction",
	placeholder,
placeholder
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"error":"%s %s %s AgentAssertion leaked"placeholder`, key.runtimeID, key.taskID, privateKey)
placeholder))
	defer server.Close()
	original := chatgptCodexModelsURL
	chatgptCodexModelsURL = server.URL
	t.Cleanup(func() { chatgptCodexModelsURL = original placeholder)

	s := &OpenAIGatewayService{placeholder
	_, err := s.FetchCodexModelsManifest(context.Background(), account, "0.137.0", "")
placeholder
	require.NotContains(t, err.Error(), key.runtimeID)
	require.NotContains(t, err.Error(), key.taskID)
placeholder
	require.NotContains(t, err.Error(), "AgentAssertion leaked")
	require.Contains(t, err.Error(), "[redacted]")
placeholder

func TestFetchCodexModelsManifestDefaultClientVersion(t *testing.T) {
	var gotClientVersion string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotClientVersion = r.URL.Query().Get("client_version")
		_, _ = w.Write([]byte(`{"models":[]placeholder`))
placeholder))
	defer server.Close()

	original := chatgptCodexModelsURL
	chatgptCodexModelsURL = server.URL
	defer func() { chatgptCodexModelsURL = original placeholder()

	s := &OpenAIGatewayService{placeholder
	if _, err := s.FetchCodexModelsManifest(context.Background(), newCodexModelsTestAccount(), "", ""); err != nil {
		t.Fatalf("FetchCodexModelsManifest returned error: %v", err)
placeholder
	if gotClientVersion != openAICodexProbeVersion {
		t.Errorf("default client_version: got %q, want %q", gotClientVersion, openAICodexProbeVersion)
placeholder
placeholder

func TestFetchCodexModelsManifestNotModified(t *testing.T) {
	var gotIfNoneMatch string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotIfNoneMatch = r.Header.Get("If-None-Match")
		w.Header().Set("ETag", `W/"abc123"`)
		w.WriteHeader(http.StatusNotModified)
placeholder))
	defer server.Close()

	original := chatgptCodexModelsURL
	chatgptCodexModelsURL = server.URL
	defer func() { chatgptCodexModelsURL = original placeholder()

	s := &OpenAIGatewayService{placeholder
	manifest, err := s.FetchCodexModelsManifest(context.Background(), newCodexModelsTestAccount(), "0.137.0", `W/"abc123"`)
	if err != nil {
		t.Fatalf("FetchCodexModelsManifest returned error: %v", err)
placeholder
	if !manifest.NotModified {
		t.Error("expected NotModified to be true")
placeholder
	if gotIfNoneMatch != `W/"abc123"` {
		t.Errorf("if-none-match header: got %q", gotIfNoneMatch)
placeholder
placeholder

func TestFetchCodexModelsManifestUpstreamError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"detail":"boom"placeholder`, http.StatusInternalServerError)
placeholder))
	defer server.Close()

	original := chatgptCodexModelsURL
	chatgptCodexModelsURL = server.URL
	defer func() { chatgptCodexModelsURL = original placeholder()

	s := &OpenAIGatewayService{placeholder
	if _, err := s.FetchCodexModelsManifest(context.Background(), newCodexModelsTestAccount(), "0.137.0", ""); err == nil {
		t.Fatal("expected error for upstream 500, got nil")
placeholder
placeholder

func TestFetchCodexModelsManifestMissingToken(t *testing.T) {
	account := newCodexModelsTestAccount()
	delete(account.Credentials, "access_token")

	s := &OpenAIGatewayService{placeholder
	if _, err := s.FetchCodexModelsManifest(context.Background(), account, "0.137.0", ""); err == nil {
		t.Fatal("expected error for missing access token, got nil")
placeholder
placeholder

func TestFetchCodexModelsManifestAPIKeyCustomUpstream(t *testing.T) {
	manifestBody := `{"models":[{"slug":"gpt-5.6"placeholder]placeholder`
	var gotRequest *http.Request
	var gotProxyURL string
	var gotAccountID int64
	var gotConcurrency int
	upstream := &codexModelsHTTPUpstreamStub{do: func(req *http.Request, proxyURL string, accountID int64, accountConcurrency int) (*http.Response, error) {
		gotRequest = req
		gotProxyURL = proxyURL
		gotAccountID = accountID
		gotConcurrency = accountConcurrency
		header := make(http.Header)
		header.Set("ETag", `W/"api-key-manifest"`)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       io.NopCloser(strings.NewReader(manifestBody)),
	placeholder, nil
placeholderplaceholder

	s := newCodexModelsAPIKeyTestService(upstream)
	manifest, err := s.FetchCodexModelsManifest(
		context.Background(),
		newCodexModelsAPIKeyTestAccount("https://upstream.example/v1"),
		"0.144.0",
		"",
	)
	if err != nil {
		t.Fatalf("FetchCodexModelsManifest returned error: %v", err)
placeholder

	if gotRequest == nil {
		t.Fatal("expected request to custom API key upstream")
placeholder
	if gotRequest.Method != http.MethodGet {
		t.Errorf("method: got %q", gotRequest.Method)
placeholder
	if gotRequest.URL.String() != "https://upstream.example/v1/models?client_version=0.144.0" {
		t.Errorf("request URL: got %q", gotRequest.URL.String())
placeholder
	if gotRequest.Header.Get("Authorization") != "Bearer sk-upstream" {
		t.Errorf("authorization header: got %q", gotRequest.Header.Get("Authorization"))
placeholder
	if gotRequest.Header.Get("Originator") != openai.CodexDefaultOriginator {
		t.Errorf("originator header: got %q", gotRequest.Header.Get("Originator"))
placeholder
	if gotRequest.Header.Get("Version") != "0.144.0" {
		t.Errorf("version header: got %q", gotRequest.Header.Get("Version"))
placeholder
	if gotRequest.Header.Get("User-Agent") != codexCLIUserAgent {
		t.Errorf("user-agent header: got %q", gotRequest.Header.Get("User-Agent"))
placeholder
	if gotRequest.Header.Get("chatgpt-account-id") != "" {
		t.Errorf("chatgpt-account-id must not be sent to API key upstream: got %q", gotRequest.Header.Get("chatgpt-account-id"))
placeholder
	if gotProxyURL != "" || gotAccountID != 2 || gotConcurrency != 3 {
		t.Errorf("upstream routing metadata: proxy=%q account_id=%d concurrency=%d", gotProxyURL, gotAccountID, gotConcurrency)
placeholder
	if string(manifest.Body) != manifestBody {
		t.Errorf("body not passed through verbatim: got %q", manifest.Body)
placeholder
	if manifest.ETag != `W/"api-key-manifest"` {
		t.Errorf("etag not passed through: got %q", manifest.ETag)
placeholder
placeholder

func TestFetchCodexModelsManifestAPIKeyConvertsStandardOpenAIModelList(t *testing.T) {
	upstreamBody := `{"object":"list","data":[{"id":"gpt-5.6","object":"model"placeholder,{"id":"  ","object":"model"placeholder,{"id":"gpt-5.6-codex","object":"model"placeholder]placeholder`
	upstream := &codexModelsHTTPUpstreamStub{do: func(_ *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
		header := make(http.Header)
		header.Set("ETag", `W/"openai-list"`)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       io.NopCloser(strings.NewReader(upstreamBody)),
	placeholder, nil
placeholderplaceholder

	s := newCodexModelsAPIKeyTestService(upstream)
	manifest, err := s.FetchCodexModelsManifest(
		context.Background(),
		newCodexModelsAPIKeyTestAccount("https://upstream.example/v1"),
		"0.144.0",
		"",
	)
	if err != nil {
		t.Fatalf("FetchCodexModelsManifest returned error: %v", err)
placeholder
	if got, want := string(manifest.Body), `{"models":[{"slug":"gpt-5.6"placeholder,{"slug":"gpt-5.6-codex"placeholder]placeholder`; got != want {
		t.Errorf("converted body: got %q, want %q", got, want)
placeholder
	require.Equal(t, codexModelsManifestBodyETag(manifest.Body), manifest.ETag)
	require.Equal(t, `W/"openai-list"`, manifest.upstreamETag)
placeholder

func TestAdjustAPIKeyCodexModelsManifest(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
placeholder{
		{
			name: "affected models disable responses lite and preserve unknown fields",
			body: `{"models":[{"slug":"gpt-5.6-sol","use_responses_lite":true,"unknown_model":{"enabled":trueplaceholderplaceholder,{"slug":"gpt-5.6-terra","use_responses_lite":trueplaceholder,{"slug":"gpt-5.6-luna","use_responses_lite":trueplaceholder],"unknown_top":{"version":1placeholderplaceholder`,
			want: `{"models":[{"slug":"gpt-5.6-sol","unknown_model":{"enabled":trueplaceholder,"use_responses_lite":falseplaceholder,{"slug":"gpt-5.6-terra","use_responses_lite":falseplaceholder,{"slug":"gpt-5.6-luna","use_responses_lite":falseplaceholder],"unknown_top":{"version":1placeholderplaceholder`,
	placeholder,
		{
			name: "unaffected model unchanged",
			body: ` {"models":[{"slug":"gpt-5.6-codex","use_responses_lite":trueplaceholder]placeholder `,
			want: ` {"models":[{"slug":"gpt-5.6-codex","use_responses_lite":trueplaceholder]placeholder `,
	placeholder,
		{
			name: "false missing and alternate entries unchanged",
			body: `{"models":[{"slug":"gpt-5.6-sol","use_responses_lite":falseplaceholder,{"slug":"gpt-5.6-terra"placeholder,null,"gpt-5.6-luna",{"slug":17,"use_responses_lite":trueplaceholder]placeholder`,
			want: `{"models":[{"slug":"gpt-5.6-sol","use_responses_lite":falseplaceholder,{"slug":"gpt-5.6-terra"placeholder,null,"gpt-5.6-luna",{"slug":17,"use_responses_lite":trueplaceholder]placeholder`,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := adjustAPIKeyCodexModelsManifest([]byte(tt.body))
		placeholder
			require.Equal(t, tt.want, string(got))
	placeholder)
placeholder
placeholder

func TestFetchCodexModelsManifestAPIKeyDisablesResponsesLiteForAffectedModels(t *testing.T) {
	const upstreamBody = `{"models":[{"slug":"gpt-5.6-sol","use_responses_lite":trueplaceholder,{"slug":"gpt-5.6-codex","use_responses_lite":trueplaceholder],"metadata":{"version":1placeholderplaceholder`
	upstream := &codexModelsHTTPUpstreamStub{do: func(_ *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Etag": []string{`"upstream-strong"`placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(upstreamBody)),
	placeholder, nil
placeholderplaceholder

	s := newCodexModelsAPIKeyTestService(upstream)
	manifest, err := s.FetchCodexModelsManifest(context.Background(), newCodexModelsAPIKeyTestAccount("https://upstream.example"), "0.145.0", "")
placeholder
	require.JSONEq(t, `{"models":[{"slug":"gpt-5.6-sol","use_responses_lite":falseplaceholder,{"slug":"gpt-5.6-codex","use_responses_lite":trueplaceholder],"metadata":{"version":1placeholderplaceholder`, string(manifest.Body))
	require.Equal(t, codexModelsManifestBodyETag(manifest.Body), manifest.ETag)
	require.Equal(t, `"upstream-strong"`, manifest.upstreamETag)

	notModified, err := s.FetchCodexModelsManifest(context.Background(), newCodexModelsAPIKeyTestAccount("https://upstream.example"), "0.145.0", manifest.ETag)
placeholder
	require.True(t, notModified.NotModified)
	require.Equal(t, manifest.ETag, notModified.ETag)
placeholder

func TestFetchCodexModelsManifestOAuthPreservesResponsesLite(t *testing.T) {
	const manifestBody = ` {"models":[{"slug":"gpt-5.6-sol","use_responses_lite":trueplaceholder]placeholder `
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(manifestBody))
placeholder))
	defer server.Close()
	original := chatgptCodexModelsURL
	chatgptCodexModelsURL = server.URL
	defer func() { chatgptCodexModelsURL = original placeholder()

	s := &OpenAIGatewayService{placeholder
	manifest, err := s.FetchCodexModelsManifest(context.Background(), newCodexModelsTestAccount(), "0.145.0", "")
placeholder
	require.Equal(t, manifestBody, string(manifest.Body))
placeholder

func TestConvertOpenAIModelListToCodexManifest(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
placeholder{
		{
			name: "standard list",
			body: `{"object":"list","data":[{"id":"m-1"placeholder,{"id":"m-2"placeholder]placeholder`,
			want: `{"models":[{"slug":"m-1"placeholder,{"slug":"m-2"placeholder]placeholder`,
	placeholder,
		{
			name: "codex manifest unchanged",
			body: `{"models":[{"slug":"m-1"placeholder]placeholder`,
			want: `{"models":[{"slug":"m-1"placeholder]placeholder`,
	placeholder,
		{
			name: "empty data unchanged",
			body: `{"object":"list","data":[]placeholder`,
			want: `{"object":"list","data":[]placeholder`,
	placeholder,
		{
			name: "data not an array unchanged",
			body: `{"object":"list","data":{"id":"m-1"placeholderplaceholder`,
			want: `{"object":"list","data":{"id":"m-1"placeholderplaceholder`,
	placeholder,
		{
			name: "entries without usable IDs unchanged",
			body: `{"object":"list","data":[{"id":""placeholder,{"object":"model"placeholder]placeholder`,
			want: `{"object":"list","data":[{"id":""placeholder,{"object":"model"placeholder]placeholder`,
	placeholder,
		{
			name: "invalid JSON unchanged",
			body: `{"data":`,
			want: `{"data":`,
	placeholder,
		{
			name: "non-object unchanged",
			body: `[]`,
			want: `[]`,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(convertOpenAIModelListToCodexManifest([]byte(tt.body))); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestFetchCodexModelsManifestRejectsInvalidEnvelope(t *testing.T) {
	tests := []struct {
		name string
		body string
placeholder{
		{name: "OpenAI models list", body: `{"object":"list","data":[]placeholder`placeholder,
		{name: "invalid JSON", body: `{"models":`placeholder,
		{name: "non-object", body: `[]`placeholder,
		{name: "null object", body: `null`placeholder,
		{name: "missing models", body: `{placeholder`placeholder,
		{name: "models object", body: `{"models":{placeholderplaceholder`placeholder,
		{name: "models null", body: `{"models":nullplaceholder`placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upstream := &codexModelsHTTPUpstreamStub{do: func(_ *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     make(http.Header),
					Body:       io.NopCloser(strings.NewReader(tt.body)),
			placeholder, nil
		placeholderplaceholder

			s := newCodexModelsAPIKeyTestService(upstream)
			_, err := s.FetchCodexModelsManifest(
				context.Background(),
				newCodexModelsAPIKeyTestAccount("https://upstream.example"),
				"0.144.0",
				"",
			)
			if err == nil {
				t.Fatal("expected invalid manifest error, got nil")
		placeholder
			if infraerrors.Reason(err) != "OPENAI_CODEX_MODELS_UPSTREAM_INVALID_MANIFEST" {
				t.Errorf("error reason: got %q", infraerrors.Reason(err))
		placeholder
			if !IsRetryableCodexModelsManifestError(err) {
				t.Error("invalid upstream manifest must be retryable")
		placeholder
	placeholder)
placeholder
placeholder

func TestFetchCodexModelsManifestAPIKeyDoesNotCacheInvalidEnvelope(t *testing.T) {
	var calls atomic.Int32
	upstream := &codexModelsHTTPUpstreamStub{do: func(_ *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
		body := `{"object":"list","data":[]placeholder`
		if calls.Add(1) > 1 {
			body = `{"models":[{"slug":"gpt-5.6"placeholder]placeholder`
	placeholder
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(body)),
	placeholder, nil
placeholderplaceholder

	s := newCodexModelsAPIKeyTestService(upstream)
	account := newCodexModelsAPIKeyTestAccount("https://upstream.example")
	if _, err := s.FetchCodexModelsManifest(context.Background(), account, "0.144.0", ""); err == nil {
		t.Fatal("expected invalid manifest error on first fetch")
placeholder
	manifest, err := s.FetchCodexModelsManifest(context.Background(), account, "0.144.0", "")
	if err != nil {
		t.Fatalf("second fetch returned error: %v", err)
placeholder
	if got, want := string(manifest.Body), `{"models":[{"slug":"gpt-5.6"placeholder]placeholder`; got != want {
		t.Errorf("body: got %q, want %q", got, want)
placeholder
	if got := calls.Load(); got != 2 {
		t.Errorf("upstream calls: got %d, want 2", got)
placeholder
placeholder

func TestFetchCodexModelsManifestAPIKeySharedRefreshSurvivesCallerCancellation(t *testing.T) {
	const manifestBody = `{"models":[{"slug":"gpt-5.6"placeholder]placeholder`
	var calls atomic.Int32
	var readStartedOnce sync.Once
	readStarted := make(chan struct{placeholder)
	deadlineRemaining := make(chan time.Duration, 1)
	release := make(chan struct{placeholder)
	upstream := &codexModelsHTTPUpstreamStub{do: func(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
		calls.Add(1)
		deadline, ok := req.Context().Deadline()
		if !ok {
			deadlineRemaining <- 0
	placeholder else {
			deadlineRemaining <- time.Until(deadline)
	placeholder
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Etag": []string{`W/"shared"`placeholderplaceholder,
			Body: &codexModelsBlockingBody{
				ctx:         req.Context(),
				readStarted: readStarted,
				startedOnce: &readStartedOnce,
				release:     release,
				body:        strings.NewReader(manifestBody),
		placeholder,
	placeholder, nil
placeholderplaceholder

	s := newCodexModelsAPIKeyTestService(upstream)
	account := newCodexModelsAPIKeyTestAccount("https://upstream.example")
	firstCtx, cancelFirst := context.WithCancel(context.Background())
	firstErr := make(chan error, 1)
	go func() {
		_, err := s.FetchCodexModelsManifest(firstCtx, account, "0.144.0", "")
		firstErr <- err
placeholder()

	select {
	case <-readStarted:
	case <-time.After(time.Second):
		t.Fatal("upstream body read did not start")
placeholder
	remaining := <-deadlineRemaining
	if remaining < 14*time.Second || remaining > codexModelsManifestRequestTimeout {
		t.Errorf("detached refresh deadline: got %s, want approximately %s", remaining, codexModelsManifestRequestTimeout)
placeholder
	cancelFirst()
	select {
	case err := <-firstErr:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("first caller error: got %v, want context.Canceled", err)
	placeholder
	case <-time.After(time.Second):
		t.Fatal("canceled caller did not return promptly")
placeholder

	secondResult := make(chan struct {
		manifest *CodexModelsManifest
		err      error
placeholder, 1)
	go func() {
		manifest, err := s.FetchCodexModelsManifest(context.Background(), account, "0.144.0", "")
		secondResult <- struct {
			manifest *CodexModelsManifest
			err      error
	placeholder{manifest: manifest, err: errplaceholder
placeholder()

	time.Sleep(50 * time.Millisecond)
	if got := calls.Load(); got != 1 {
		t.Errorf("upstream calls before shared refresh completed: got %d, want 1", got)
placeholder
	close(release)
	select {
	case result := <-secondResult:
		if result.err != nil {
			t.Fatalf("second caller returned error: %v", result.err)
	placeholder
		if string(result.manifest.Body) != manifestBody {
			t.Errorf("second caller body: got %q", result.manifest.Body)
	placeholder
	case <-time.After(time.Second):
		t.Fatal("second caller did not receive shared refresh result")
placeholder
	if got := calls.Load(); got != 1 {
		t.Errorf("total upstream calls: got %d, want 1", got)
placeholder
placeholder

func TestFetchCodexModelsManifestAPIKeyConcurrentRequestsShareRefresh(t *testing.T) {
	const callers = 8
	var calls atomic.Int32
	started := make(chan struct{placeholder)
	var startedOnce sync.Once
	release := make(chan struct{placeholder)
	upstream := &codexModelsHTTPUpstreamStub{do: func(_ *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
		calls.Add(1)
		startedOnce.Do(func() { close(started) placeholder)
		<-release
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"models":[]placeholder`)),
	placeholder, nil
placeholderplaceholder

	s := newCodexModelsAPIKeyTestService(upstream)
	account := newCodexModelsAPIKeyTestAccount("https://upstream.example")
	begin := make(chan struct{placeholder)
	errs := make(chan error, callers)
	for i := 0; i < callers; i++ {
		go func() {
			<-begin
			_, err := s.FetchCodexModelsManifest(context.Background(), account, "0.144.0", "")
			errs <- err
	placeholder()
placeholder
	close(begin)
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("upstream request did not start")
placeholder
	time.Sleep(50 * time.Millisecond)
	if got := calls.Load(); got != 1 {
		t.Errorf("concurrent upstream calls: got %d, want 1", got)
placeholder
	close(release)
	for i := 0; i < callers; i++ {
		if err := <-errs; err != nil {
			t.Errorf("caller %d returned error: %v", i, err)
	placeholder
placeholder
placeholder

func TestFetchCodexModelsManifestAPIKeyFreshCacheHandlesETagLocally(t *testing.T) {
	var calls atomic.Int32
	upstream := &codexModelsHTTPUpstreamStub{do: func(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
		calls.Add(1)
		if got := req.Header.Get("If-None-Match"); got != "" {
			t.Errorf("cache refresh must not inherit a caller's If-None-Match: got %q", got)
	placeholder
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Etag": []string{`W/"cached"`placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"models":[]placeholder`)),
	placeholder, nil
placeholderplaceholder

	s := newCodexModelsAPIKeyTestService(upstream)
	account := newCodexModelsAPIKeyTestAccount("https://upstream.example")
	if _, err := s.FetchCodexModelsManifest(context.Background(), account, "0.144.0", ""); err != nil {
		t.Fatalf("initial fetch returned error: %v", err)
placeholder
	manifest, err := s.FetchCodexModelsManifest(context.Background(), account, "0.144.0", `W/"cached"`)
	if err != nil {
		t.Fatalf("cached fetch returned error: %v", err)
placeholder
	if !manifest.NotModified {
		t.Fatal("matching cached ETag must return NotModified")
placeholder
	if got := calls.Load(); got != 1 {
		t.Errorf("upstream calls: got %d, want 1", got)
placeholder
placeholder

func TestFetchCodexModelsManifestAPIKeyCacheKeyIsolatesRequestIdentity(t *testing.T) {
	var calls atomic.Int32
	upstream := &codexModelsHTTPUpstreamStub{do: func(_ *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
		calls.Add(1)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"models":[]placeholder`)),
	placeholder, nil
placeholderplaceholder
	s := newCodexModelsAPIKeyTestService(upstream)

	base := newCodexModelsAPIKeyTestAccount("https://upstream.example")
	fetch := func(account *Account, version string) {
	placeholder
		if _, err := s.FetchCodexModelsManifest(context.Background(), account, version, ""); err != nil {
			t.Fatalf("fetch returned error: %v", err)
	placeholder
placeholder
	fetch(base, "0.144.0")
	fetch(base, "0.144.0")

	differentAccount := newCodexModelsAPIKeyTestAccount("https://upstream.example")
	differentAccount.ID = 3
	fetch(differentAccount, "0.144.0")

	differentToken := newCodexModelsAPIKeyTestAccount("https://upstream.example")
	differentToken.Credentials["api_key"] = "sk-other"
	fetch(differentToken, "0.144.0")

	differentUpstream := newCodexModelsAPIKeyTestAccount("https://other-upstream.example")
	fetch(differentUpstream, "0.144.0")
	fetch(base, "0.145.0")

	differentHeaders := newCodexModelsAPIKeyTestAccount("https://upstream.example")
	differentHeaders.Credentials[credKeyHeaderOverrideEnabled] = true
	differentHeaders.Credentials[credKeyHeaderOverrides] = map[string]any{"x-tenant": "other"placeholder
	fetch(differentHeaders, "0.144.0")

	proxyID := int64(9)
	differentProxy := newCodexModelsAPIKeyTestAccount("https://upstream.example")
	differentProxy.ProxyID = &proxyID
	differentProxy.Proxy = &Proxy{Protocol: "http", Host: "127.0.0.1", Port: 8080placeholder
	fetch(differentProxy, "0.144.0")
	fetch(differentProxy, "0.144.0")

	if got := calls.Load(); got != 7 {
		t.Errorf("isolated upstream calls: got %d, want 7", got)
placeholder
placeholder

func TestFetchCodexModelsManifestAPIKeyCacheBoundsEntriesAndBodySize(t *testing.T) {
	var calls atomic.Int32
	upstream := &codexModelsHTTPUpstreamStub{do: func(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
		calls.Add(1)
		body := `{"models":[]placeholder`
		if strings.Contains(req.URL.Host, "large") {
			body = `{"models":[],"padding":"` + strings.Repeat("x", (1<<20)+1) + `"placeholder`
	placeholder
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(body)),
	placeholder, nil
placeholderplaceholder
	s := newCodexModelsAPIKeyTestService(upstream)
	fetch := func(account *Account) {
	placeholder
		if _, err := s.FetchCodexModelsManifest(context.Background(), account, "0.144.0", ""); err != nil {
			t.Fatalf("fetch returned error: %v", err)
	placeholder
placeholder

	small := newCodexModelsAPIKeyTestAccount("https://small.example")
	fetch(small)
	fetch(small)
	large := newCodexModelsAPIKeyTestAccount("https://large.example")
	large.ID = 3
	fetch(large)
	fetch(large)
	if got := calls.Load(); got != 3 {
		t.Fatalf("body-size bounded cache calls: got %d, want 3", got)
placeholder

	for i := int64(10); i < 75; i++ {
		account := newCodexModelsAPIKeyTestAccount("https://bounded.example")
		account.ID = i
		fetch(account)
placeholder
	last := newCodexModelsAPIKeyTestAccount("https://bounded.example")
	last.ID = 74
	fetch(last)
	if got := calls.Load(); got != 68 {
		t.Fatalf("most recent cache entry was not retained: calls=%d, want 68", got)
placeholder
	first := newCodexModelsAPIKeyTestAccount("https://bounded.example")
	first.ID = 10
	fetch(first)
	if got := calls.Load(); got != 69 {
		t.Errorf("oldest cache entry was not evicted: calls=%d, want 69", got)
placeholder
placeholder

func TestFetchCodexModelsManifestAPIKeyServesStaleWhileRefreshing(t *testing.T) {
	var calls atomic.Int32
	refreshStarted := make(chan struct{placeholder)
	releaseRefresh := make(chan struct{placeholder)
	upstream := &codexModelsHTTPUpstreamStub{do: func(_ *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
		call := calls.Add(1)
		body := `{"models":[{"slug":"old"placeholder]placeholder`
		if call > 1 {
			if call == 2 {
				close(refreshStarted)
		placeholder
			<-releaseRefresh
			body = `{"models":[{"slug":"new"placeholder]placeholder`
	placeholder
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(body)),
	placeholder, nil
placeholderplaceholder
	s := newCodexModelsAPIKeyTestService(upstream)
	account := newCodexModelsAPIKeyTestAccount("https://upstream.example")
	if _, err := s.FetchCodexModelsManifest(context.Background(), account, "0.144.0", ""); err != nil {
		t.Fatalf("initial fetch returned error: %v", err)
placeholder

	s.codexModelsManifestCache.mu.Lock()
	for key, entry := range s.codexModelsManifestCache.entries {
		entry.expiresAt = time.Now().Add(-time.Second)
		s.codexModelsManifestCache.entries[key] = entry
placeholder
	s.codexModelsManifestCache.mu.Unlock()

	resultCh := make(chan struct {
		manifest *CodexModelsManifest
		err      error
placeholder, 1)
	go func() {
		manifest, err := s.FetchCodexModelsManifest(context.Background(), account, "0.144.0", "")
		resultCh <- struct {
			manifest *CodexModelsManifest
			err      error
	placeholder{manifest: manifest, err: errplaceholder
placeholder()
	select {
	case <-refreshStarted:
	case <-time.After(time.Second):
		t.Fatal("background refresh did not start")
placeholder

	var staleResult struct {
		manifest *CodexModelsManifest
		err      error
placeholder
	select {
	case staleResult = <-resultCh:
	case <-time.After(100 * time.Millisecond):
		t.Error("stale manifest was not returned while refresh was blocked")
		close(releaseRefresh)
		staleResult = <-resultCh
placeholder
	if staleResult.err != nil {
		t.Fatalf("stale fetch returned error: %v", staleResult.err)
placeholder
	if got := string(staleResult.manifest.Body); got != `{"models":[{"slug":"old"placeholder]placeholder` {
		t.Errorf("stale body: got %q", got)
placeholder
	if got := calls.Load(); got != 2 {
		t.Errorf("upstream calls during stale refresh: got %d, want 2", got)
placeholder

	select {
	case <-releaseRefresh:
	default:
		close(releaseRefresh)
placeholder
	deadline := time.Now().Add(time.Second)
	for {
		manifest, err := s.FetchCodexModelsManifest(context.Background(), account, "0.144.0", "")
		if err == nil && string(manifest.Body) == `{"models":[{"slug":"new"placeholder]placeholder` {
			break
	placeholder
		if time.Now().After(deadline) {
			t.Fatalf("refreshed manifest was not cached: manifest=%v err=%v", manifest, err)
	placeholder
		time.Sleep(10 * time.Millisecond)
placeholder
	if got := calls.Load(); got != 2 {
		t.Errorf("stale refresh was not deduplicated: calls=%d, want 2", got)
placeholder
placeholder

func TestFetchCodexModelsManifestAPIKeyRevalidatesStaleETag(t *testing.T) {
	var calls atomic.Int32
	refreshDone := make(chan struct{placeholder)
	upstream := &codexModelsHTTPUpstreamStub{do: func(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
		call := calls.Add(1)
		if call == 1 {
			header := make(http.Header)
			header.Set("ETag", `"upstream-cached"`)
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     header,
				Body:       io.NopCloser(strings.NewReader(`{"models":[{"slug":"gpt-5.6-sol","use_responses_lite":trueplaceholder]placeholder`)),
		placeholder, nil
	placeholder
		if got := req.Header.Get("If-None-Match"); got != `"upstream-cached"` {
			t.Errorf("background revalidation If-None-Match: got %q", got)
	placeholder
		close(refreshDone)
		header := make(http.Header)
		header.Set("ETag", `"upstream-cached"`)
		return &http.Response{StatusCode: http.StatusNotModified, Header: header, Body: http.NoBodyplaceholder, nil
placeholderplaceholder
	s := newCodexModelsAPIKeyTestService(upstream)
	account := newCodexModelsAPIKeyTestAccount("https://upstream.example")
	if _, err := s.FetchCodexModelsManifest(context.Background(), account, "0.144.0", ""); err != nil {
		t.Fatalf("initial fetch returned error: %v", err)
placeholder
	s.codexModelsManifestCache.mu.Lock()
	for key, entry := range s.codexModelsManifestCache.entries {
		entry.expiresAt = time.Now().Add(-time.Second)
		s.codexModelsManifestCache.entries[key] = entry
placeholder
	s.codexModelsManifestCache.mu.Unlock()

	manifest, err := s.FetchCodexModelsManifest(context.Background(), account, "0.144.0", "")
	if err != nil {
		t.Fatalf("stale fetch returned error: %v", err)
placeholder
	if got := string(manifest.Body); got != `{"models":[{"slug":"gpt-5.6-sol","use_responses_lite":falseplaceholder]placeholder` {
		t.Fatalf("stale body: got %q", got)
placeholder
	select {
	case <-refreshDone:
	case <-time.After(time.Second):
		t.Fatal("ETag revalidation did not complete")
placeholder

	deadline := time.Now().Add(time.Second)
	for {
		s.codexModelsManifestCache.mu.Lock()
		fresh := false
		for _, entry := range s.codexModelsManifestCache.entries {
			fresh = time.Now().Before(entry.expiresAt)
	placeholder
		s.codexModelsManifestCache.mu.Unlock()
		if fresh {
			break
	placeholder
		if time.Now().After(deadline) {
			t.Fatal("304 revalidation did not renew the cached manifest")
	placeholder
		time.Sleep(10 * time.Millisecond)
placeholder
	manifest, err = s.FetchCodexModelsManifest(context.Background(), account, "0.144.0", "")
	if err != nil || string(manifest.Body) != `{"models":[{"slug":"gpt-5.6-sol","use_responses_lite":falseplaceholder]placeholder` {
		t.Fatalf("renewed cached manifest: body=%q err=%v", manifest.Body, err)
placeholder
	if got := calls.Load(); got != 2 {
		t.Errorf("upstream calls: got %d, want 2", got)
placeholder
placeholder

func TestFetchCodexModelsManifestAPIKeyColdCacheHandlesNotModifiedLocally(t *testing.T) {
	var gotIfNoneMatch string
	upstream := &codexModelsHTTPUpstreamStub{do: func(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
		gotIfNoneMatch = req.Header.Get("If-None-Match")
		header := make(http.Header)
		header.Set("ETag", `W/"api-key-manifest"`)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       io.NopCloser(strings.NewReader(`{"models":[]placeholder`)),
	placeholder, nil
placeholderplaceholder

	s := newCodexModelsAPIKeyTestService(upstream)
	manifest, err := s.FetchCodexModelsManifest(
		context.Background(),
		newCodexModelsAPIKeyTestAccount("https://upstream.example"),
		"0.144.0",
		`W/"api-key-manifest"`,
	)
	if err != nil {
		t.Fatalf("FetchCodexModelsManifest returned error: %v", err)
placeholder
	if !manifest.NotModified {
		t.Error("expected NotModified to be true")
placeholder
	if manifest.ETag != `W/"api-key-manifest"` {
		t.Errorf("etag not passed through: got %q", manifest.ETag)
placeholder
	if gotIfNoneMatch != "" {
		t.Errorf("cold shared refresh must not inherit caller if-none-match: got %q", gotIfNoneMatch)
placeholder
placeholder

func TestFetchCodexModelsManifestAPIKeyDoesNotCacheUnexpectedColdNotModified(t *testing.T) {
	var calls atomic.Int32
	upstream := &codexModelsHTTPUpstreamStub{do: func(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
		calls.Add(1)
		if got := req.Header.Get("If-None-Match"); got != "" {
			t.Errorf("cold shared refresh If-None-Match: got %q", got)
	placeholder
		header := make(http.Header)
		header.Set("ETag", `W/"unexpected"`)
		return &http.Response{StatusCode: http.StatusNotModified, Header: header, Body: http.NoBodyplaceholder, nil
placeholderplaceholder
	s := newCodexModelsAPIKeyTestService(upstream)
	account := newCodexModelsAPIKeyTestAccount("https://upstream.example")
	for i := 0; i < 2; i++ {
		manifest, err := s.FetchCodexModelsManifest(context.Background(), account, "0.144.0", "")
		if err != nil {
			t.Fatalf("fetch %d returned error: %v", i, err)
	placeholder
		if !manifest.NotModified {
			t.Fatalf("fetch %d: expected upstream NotModified response", i)
	placeholder
placeholder
	if got := calls.Load(); got != 2 {
		t.Errorf("unexpected cold 304 was cached: upstream calls=%d, want 2", got)
placeholder
placeholder

func TestFetchCodexModelsManifestAPIKeyPreservesBaseURLQuery(t *testing.T) {
	var gotURL string
	upstream := &codexModelsHTTPUpstreamStub{do: func(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
		gotURL = req.URL.String()
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"models":[]placeholder`)),
	placeholder, nil
placeholderplaceholder

	s := newCodexModelsAPIKeyTestService(upstream)
	_, err := s.FetchCodexModelsManifest(
		context.Background(),
		newCodexModelsAPIKeyTestAccount("https://upstream.example/v1?tenant=acme"),
		"0.144.0",
		"",
	)
	if err != nil {
		t.Fatalf("FetchCodexModelsManifest returned error: %v", err)
placeholder
	if gotURL != "https://upstream.example/v1/models?client_version=0.144.0&tenant=acme" {
		t.Errorf("request URL: got %q", gotURL)
placeholder
placeholder

func TestFetchCodexModelsManifestAPIKeyRejectsBaseURLFragment(t *testing.T) {
	called := false
	upstream := &codexModelsHTTPUpstreamStub{do: func(_ *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
		called = true
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"models":[]placeholder`)),
	placeholder, nil
placeholderplaceholder

	s := newCodexModelsAPIKeyTestService(upstream)
	_, err := s.FetchCodexModelsManifest(
		context.Background(),
		newCodexModelsAPIKeyTestAccount("https://upstream.example/v1#models"),
		"0.144.0",
		"",
	)
	if err == nil {
		t.Fatal("expected invalid upstream base URL error, got nil")
placeholder
	if infraerrors.Reason(err) != "OPENAI_CODEX_MODELS_API_KEY_UPSTREAM_INVALID" {
		t.Errorf("error reason: got %q", infraerrors.Reason(err))
placeholder
	if called {
		t.Fatal("fragment-bearing base URL must be rejected before the upstream request")
placeholder
placeholder

// codexModelsAccountStateRepo records account state transitions triggered by
// manifest upstream errors (#4544).
type codexModelsAccountStateRepo struct {
	AccountRepository
	mu                  sync.Mutex
	setErrorCalls       int
	lastErrorMsg        string
	setTempUnschedCalls int
	lastTempReason      string
placeholder

func (r *codexModelsAccountStateRepo) SetError(_ context.Context, _ int64, errorMsg string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.setErrorCalls++
	r.lastErrorMsg = errorMsg
	return nil
placeholder

func (r *codexModelsAccountStateRepo) SetTempUnschedulable(_ context.Context, _ int64, _ time.Time, reason string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.setTempUnschedCalls++
	r.lastTempReason = reason
	return nil
placeholder

func newCodexModels401TestService(repo AccountRepository) *OpenAIGatewayService {
	rateLimitService := NewRateLimitService(repo, nil, &config.Config{placeholder, nil, nil)
	s := &OpenAIGatewayService{rateLimitService: rateLimitServiceplaceholder
	rateLimitService.SetAccountRuntimeBlocker(s)
	return s
placeholder

func TestFetchCodexModelsManifestOAuth401MarksAccountUnschedulable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"detail":{"message":"invalid token"placeholderplaceholder`))
placeholder))
	defer server.Close()

	original := chatgptCodexModelsURL
	chatgptCodexModelsURL = server.URL
	defer func() { chatgptCodexModelsURL = original placeholder()

	repo := &codexModelsAccountStateRepo{placeholder
	s := newCodexModels401TestService(repo)
	account := newCodexModelsTestAccount()
	account.Credentials["refresh_token"] = "test-refresh-token"

	_, err := s.FetchCodexModelsManifest(context.Background(), account, "0.137.0", "")
placeholder
	require.True(t, IsRetryableCodexModelsManifestError(err), "manifest 401 should allow account failover")
	require.Equal(t, 1, repo.setTempUnschedCalls, "OAuth 401 should temp-unschedule the account")
	require.Equal(t, 0, repo.setErrorCalls)
	require.True(t, s.isOpenAIAccountRuntimeBlocked(account), "account should be runtime-blocked after manifest 401")
placeholder

func TestFetchCodexModelsManifestOAuth401TokenRevokedDisablesAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"code":"token_revoked","message":"token has been revoked"placeholderplaceholder`))
placeholder))
	defer server.Close()

	original := chatgptCodexModelsURL
	chatgptCodexModelsURL = server.URL
	defer func() { chatgptCodexModelsURL = original placeholder()

	repo := &codexModelsAccountStateRepo{placeholder
	s := newCodexModels401TestService(repo)
	account := newCodexModelsTestAccount()
	account.Credentials["refresh_token"] = "test-refresh-token"

	_, err := s.FetchCodexModelsManifest(context.Background(), account, "0.137.0", "")
placeholder
	require.True(t, IsRetryableCodexModelsManifestError(err))
	require.Equal(t, 1, repo.setErrorCalls, "revoked token should permanently disable the account")
	require.Contains(t, repo.lastErrorMsg, "Token revoked")
	require.Equal(t, 0, repo.setTempUnschedCalls)
placeholder

func TestFetchCodexModelsManifestAgentIdentity401DoesNotDisableAccount(t *testing.T) {
	key, privateKey := newTestAgentIdentityKey(t)
	account := &Account{
		ID:       6,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"auth_mode":          OpenAIAuthModeAgentIdentity,
			"agent_runtime_id":   key.runtimeID,
			"agent_private_key":  privateKey,
			"task_id":            key.taskID,
			"chatgpt_account_id": "acc-agent-401",
	placeholder,
placeholder
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"detail":"some non-task 401"placeholder`))
placeholder))
	defer server.Close()

	original := chatgptCodexModelsURL
	chatgptCodexModelsURL = server.URL
	defer func() { chatgptCodexModelsURL = original placeholder()

	repo := &codexModelsAccountStateRepo{placeholder
	s := newCodexModels401TestService(repo)

	_, err := s.FetchCodexModelsManifest(context.Background(), account, "0.137.0", "")
placeholder
	require.Equal(t, 0, repo.setErrorCalls, "agent identity 401s must not disable the account")
	require.Equal(t, 0, repo.setTempUnschedCalls)
placeholder

func TestFetchCodexModelsManifestAPIKey401KeepsNoFailoverAndNoDisable(t *testing.T) {
	upstream := &codexModelsHTTPUpstreamStub{do: func(_ *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Status:     "401 Unauthorized",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"error":"invalid api key"placeholder`)),
	placeholder, nil
placeholderplaceholder

	repo := &codexModelsAccountStateRepo{placeholder
	s := newCodexModelsAPIKeyTestService(upstream)
	s.rateLimitService = NewRateLimitService(repo, nil, &config.Config{placeholder, nil, nil)

	_, err := s.FetchCodexModelsManifest(
		context.Background(),
		newCodexModelsAPIKeyTestAccount("https://upstream.example"),
		"0.144.0",
		"",
	)
placeholder
	require.False(t, IsRetryableCodexModelsManifestError(err), "custom upstream manifest 401 keeps the no-failover behavior")
	require.Equal(t, 0, repo.setErrorCalls, "custom upstream manifest 401 must not disable the account")
	require.Equal(t, 0, repo.setTempUnschedCalls)
placeholder

func TestFetchCodexModelsManifestAPIKeyUpstreamError(t *testing.T) {
	upstream := &codexModelsHTTPUpstreamStub{do: func(_ *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Status:     "429 Too Many Requests",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"error":"rate limited"placeholder`)),
	placeholder, nil
placeholderplaceholder

	s := newCodexModelsAPIKeyTestService(upstream)
	_, err := s.FetchCodexModelsManifest(
		context.Background(),
		newCodexModelsAPIKeyTestAccount("https://upstream.example"),
		"0.144.0",
		"",
	)
	if err == nil {
		t.Fatal("expected error for upstream 429, got nil")
placeholder
	if infraerrors.Code(err) != http.StatusBadGateway {
		t.Errorf("error status: got %d, want %d", infraerrors.Code(err), http.StatusBadGateway)
placeholder
	if infraerrors.Reason(err) != "OPENAI_CODEX_MODELS_UPSTREAM_FAILED" {
		t.Errorf("error reason: got %q", infraerrors.Reason(err))
placeholder
placeholder

func TestFetchCodexModelsManifestAPIKeyRejectsOfficialOpenAIBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
placeholder{
		{name: "missing base URL"placeholder,
		{name: "official host", baseURL: "https://api.openai.com"placeholder,
		{name: "official versioned URL", baseURL: "https://API.OPENAI.COM:443/v1/"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newCodexModelsAPIKeyTestService(&codexModelsHTTPUpstreamStub{do: func(_ *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
				t.Fatal("official OpenAI API key must not be used as a Codex manifest upstream")
				return nil, nil
		placeholderplaceholder)

			_, err := s.FetchCodexModelsManifest(
				context.Background(),
				newCodexModelsAPIKeyTestAccount(tt.baseURL),
				"0.144.0",
				"",
			)
			if err == nil {
				t.Fatal("expected unsupported API key upstream error, got nil")
		placeholder
			if infraerrors.Reason(err) != "OPENAI_CODEX_MODELS_API_KEY_UPSTREAM_UNSUPPORTED" {
				t.Errorf("error reason: got %q", infraerrors.Reason(err))
		placeholder
	placeholder)
placeholder
placeholder
