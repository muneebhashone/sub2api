package service

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func upstreamModelSyncTestConfig() *config.Config {
	return &config.Config{
		Security: config.SecurityConfig{
			URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholder,
	placeholder,
placeholder
placeholder

func grokOAuthModelSyncTestAccount(baseURL string) *Account {
	credentials := map[string]any{
		"access_token":  "oauth-access-token",
		"refresh_token": "oauth-refresh-token",
		"expires_at":    time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
		"sub":           "grok-user-id",
		"email":         "grok-user@example.com",
placeholder
	if strings.TrimSpace(baseURL) != "" {
		credentials["base_url"] = baseURL
placeholder
placeholder
		ID:          10,
		Platform:    PlatformGrok,
		Type:        AccountTypeOAuth,
		Credentials: credentials,
placeholder
placeholder

func TestBuildV1ModelsURL(t *testing.T) {
	t.Parallel()

	require.Equal(t, "https://api.anthropic.com/v1/models", buildV1ModelsURL("https://api.anthropic.com"))
	require.Equal(t, "https://api.anthropic.com/v1/models", buildV1ModelsURL("https://api.anthropic.com/v1"))
	require.Equal(t, "https://api.anthropic.com/v1/models", buildV1ModelsURL("https://api.anthropic.com/v1/models"))
	require.Equal(t, "https://gateway.example.com/antigravity/v1/models", buildV1ModelsURL("https://gateway.example.com/antigravity/"))
placeholder

func TestBuildOpenAIModelsURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		base string
		want string
placeholder{
		{
			name: "zhipu v4 coding base url",
			base: "https://open.bigmodel.cn/api/coding/paas/v4",
			want: "https://open.bigmodel.cn/api/coding/paas/v4/models",
	placeholder,
		{
			name: "openai v1 base url",
			base: "https://api.openai.com/v1",
			want: "https://api.openai.com/v1/models",
	placeholder,
		{
			name: "models url unchanged",
			base: "https://api.openai.com/v1/models",
			want: "https://api.openai.com/v1/models",
	placeholder,
		{
			name: "host fallback uses v1",
			base: "https://api.openai.com",
			want: "https://api.openai.com/v1/models",
	placeholder,
		{
			name: "trailing slash on v4",
			base: "https://open.bigmodel.cn/api/coding/paas/v4/",
			want: "https://open.bigmodel.cn/api/coding/paas/v4/models",
	placeholder,
		{
			name: "v2 base url",
			base: "https://gateway.example.com/openai/v2",
			want: "https://gateway.example.com/openai/v2/models",
	placeholder,
		{
			name: "v3 base url",
			base: "https://gateway.example.com/openai/v3",
			want: "https://gateway.example.com/openai/v3/models",
	placeholder,
placeholder

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, buildOpenAIModelsURL(tt.base))
	placeholder)
placeholder
placeholder

func TestBuildGeminiModelsURL(t *testing.T) {
	t.Parallel()

	require.Equal(t, "https://generativelanguage.googleapis.com/v1beta/models", buildGeminiModelsURL("https://generativelanguage.googleapis.com"))
	require.Equal(t, "https://generativelanguage.googleapis.com/v1beta/models", buildGeminiModelsURL("https://generativelanguage.googleapis.com/v1beta"))
	require.Equal(t, "https://generativelanguage.googleapis.com/v1beta/models", buildGeminiModelsURL("https://generativelanguage.googleapis.com/v1beta/models"))
placeholder

func TestExtractUpstreamModelIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
		want []string
placeholder{
		{
			name: "openai and anthropic data array",
			body: `{"data":[{"id":"claude-sonnet-4-5"placeholder,{"id":"gpt-5"placeholder,{"id":"gpt-5"placeholder,{"id":""placeholder]placeholder`,
			want: []string{"claude-sonnet-4-5", "gpt-5"placeholder,
	placeholder,
		{
			name: "gemini models array strips prefix",
			body: `{"models":[{"name":"models/gemini-2.5-pro"placeholder,{"name":"gemini-2.5-flash"placeholder]placeholder`,
			want: []string{"gemini-2.5-flash", "gemini-2.5-pro"placeholder,
	placeholder,
		{
			name: "top level array",
			body: `[{"id":"z-model"placeholder,{"name":"models/a-model"placeholder]`,
			want: []string{"a-model", "z-model"placeholder,
	placeholder,
		{
			name: "standard id wins over provider-specific model field",
			body: `{"data":[{"id":"canonical-id","model":"display-model"placeholder]placeholder`,
			want: []string{"canonical-id"placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := extractUpstreamModelIDs([]byte(tt.body))
		placeholder
			require.Equal(t, tt.want, got)
	placeholder)
placeholder
placeholder

func TestExtractGrokUpstreamModelIDs(t *testing.T) {
	t.Parallel()

	models, err := extractGrokUpstreamModelIDs([]byte(`{"data":[{"id":"display-id","model":"grok-4.5"placeholder,{"modelId":"grok-build-0.1"placeholder,{"model_id":"grok-composer-2.5-fast"placeholder,{"name":"Grok Meta Display Name","_meta":{"model":"grok-meta"placeholderplaceholder,{"name":"grok-name"placeholder,{"id":"grok-safe","_meta":"not-an-object"placeholder]placeholder`))
placeholder
	require.Equal(t, []string{"grok-4.5", "grok-build-0.1", "grok-composer-2.5-fast", "grok-meta", "grok-name", "grok-safe"placeholder, models)
placeholder

func TestBuildUpstreamModelsRequestsForAPIKeyAccounts(t *testing.T) {
	t.Parallel()

	svc := &AccountTestService{cfg: upstreamModelSyncTestConfig()placeholder
	ctx := context.Background()

	anthropicReq, err := svc.buildAnthropicUpstreamModelsRequest(ctx, &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeAPIKey,
placeholder
			"api_key":  "anthropic-key",
			"base_url": "https://anthropic.example.com/v1",
	placeholder,
placeholder)
placeholder
	require.Equal(t, "https://anthropic.example.com/v1/models", anthropicReq.URL.String())
	require.Equal(t, "anthropic-key", anthropicReq.Header.Get("x-api-key"))
	require.Equal(t, "2023-06-01", anthropicReq.Header.Get("anthropic-version"))

	anthropicBearerReq, err := svc.buildAnthropicUpstreamModelsRequest(ctx, &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeAPIKey,
placeholder
			"api_key":  "ollama-key",
			"base_url": "https://ollama.com",
	placeholder,
		Extra: map[string]any{
			"anthropic_apikey_auth_scheme": AnthropicAPIKeyAuthSchemeAuthorizationBearer,
	placeholder,
placeholder)
placeholder
	require.Equal(t, "https://ollama.com/v1/models", anthropicBearerReq.URL.String())
	require.Equal(t, "Bearer ollama-key", anthropicBearerReq.Header.Get("Authorization"))
	require.Empty(t, anthropicBearerReq.Header.Get("x-api-key"))
	require.Equal(t, "2023-06-01", anthropicBearerReq.Header.Get("anthropic-version"))

	openAIReq, err := svc.buildOpenAIUpstreamModelsRequest(ctx, &Account{
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
placeholder
			"api_key":  "openai-key",
			"base_url": "https://openai.example.com",
	placeholder,
placeholder)
placeholder
	require.Equal(t, "https://openai.example.com/v1/models", openAIReq.URL.String())
	require.Equal(t, "Bearer openai-key", openAIReq.Header.Get("Authorization"))

	grokReq, err := svc.buildUpstreamModelsRequest(ctx, &Account{
		Platform: PlatformGrok,
		Type:     AccountTypeAPIKey,
placeholder
			"api_key":  "xai-key",
			"base_url": "https://xai.example.com/v1",
	placeholder,
placeholder)
placeholder
	require.Equal(t, "https://xai.example.com/v1/models", grokReq.URL.String())
	require.Equal(t, "Bearer xai-key", grokReq.Header.Get("Authorization"))

	geminiReq, err := svc.buildGeminiUpstreamModelsRequest(ctx, &Account{
placeholder
		Type:     AccountTypeAPIKey,
placeholder
			"api_key":  "gemini-key",
			"base_url": "https://generativelanguage.googleapis.com/v1beta",
	placeholder,
placeholder)
placeholder
	require.Equal(t, "https://generativelanguage.googleapis.com/v1beta/models", geminiReq.URL.String())
	require.Equal(t, "gemini-key", geminiReq.Header.Get("x-goog-api-key"))

	antigravityReq, err := svc.buildAntigravityAPIKeyModelsRequest(ctx, &Account{
		Platform: PlatformAntigravity,
		Type:     AccountTypeAPIKey,
placeholder
			"api_key":  "antigravity-key",
			"base_url": "https://gateway.example.com/antigravity",
	placeholder,
placeholder)
placeholder
	require.Equal(t, "https://gateway.example.com/antigravity/v1/models", antigravityReq.URL.String())
	require.Equal(t, "antigravity-key", antigravityReq.Header.Get("x-api-key"))
placeholder

func TestBuildUpstreamModelsRequestSupportsGrokOAuth(t *testing.T) {
	t.Parallel()

	svc := &AccountTestService{
		cfg:               upstreamModelSyncTestConfig(),
		grokTokenProvider: NewGrokTokenProvider(nil, nil),
placeholder
	req, err := svc.buildUpstreamModelsRequest(context.Background(), grokOAuthModelSyncTestAccount(""))
placeholder
	require.Equal(t, "https://cli-chat-proxy.grok.com/v1/models", req.URL.String())
	require.Equal(t, "Bearer oauth-access-token", req.Header.Get("Authorization"))
	require.Equal(t, grokCLIVersion, req.Header.Get("X-Grok-Client-Version"))
	require.Equal(t, "interactive", req.Header.Get("X-Grok-Client-Mode"))
	require.Equal(t, grokUpstreamUserAgent, req.Header.Get("User-Agent"))
	require.Equal(t, "grok-user-id", req.Header.Get("X-UserID"))
	require.Equal(t, "grok-user@example.com", req.Header.Get("X-Email"))
	require.NotContains(t, req.Header.Get("Authorization"), "oauth-refresh-token")
placeholder

func TestBuildUpstreamModelsRequestGrokOAuthRequiresTokenProvider(t *testing.T) {
	t.Parallel()

	svc := &AccountTestService{cfg: upstreamModelSyncTestConfig()placeholder
	_, err := svc.buildUpstreamModelsRequest(context.Background(), grokOAuthModelSyncTestAccount(""))
placeholder

	var syncErr *UpstreamModelSyncError
	require.True(t, errors.As(err, &syncErr))
	require.Equal(t, UpstreamModelSyncErrorConfiguration, syncErr.Kind)
	require.Contains(t, syncErr.SafeMessage(), "token provider")
placeholder

func TestBuildAntigravityAPIKeyModelsRequestRejectsOfficialCloudCodeBase(t *testing.T) {
	t.Parallel()

	svc := &AccountTestService{cfg: upstreamModelSyncTestConfig()placeholder
	_, err := svc.buildAntigravityAPIKeyModelsRequest(context.Background(), &Account{
		Platform: PlatformAntigravity,
		Type:     AccountTypeAPIKey,
placeholder
			"api_key":  "antigravity-key",
			"base_url": "https://cloudcode-pa.googleapis.com",
	placeholder,
placeholder)
placeholder

	var syncErr *UpstreamModelSyncError
	require.True(t, errors.As(err, &syncErr))
	require.Equal(t, UpstreamModelSyncErrorUnsupported, syncErr.Kind)
	require.Contains(t, syncErr.SafeMessage(), "compatible gateway")
placeholder

func TestBuildAnthropicUpstreamModelsRequestRejectsBedrock(t *testing.T) {
	t.Parallel()

	svc := &AccountTestService{cfg: upstreamModelSyncTestConfig()placeholder
	_, err := svc.buildAnthropicUpstreamModelsRequest(context.Background(), &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeBedrock,
placeholder)
placeholder

	var syncErr *UpstreamModelSyncError
	require.True(t, errors.As(err, &syncErr))
	require.Equal(t, UpstreamModelSyncErrorUnsupported, syncErr.Kind)
placeholder

func TestFetchUpstreamSupportedModelsParsesOpenAIResponse(t *testing.T) {
	t.Parallel()

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"data":[{"id":"gpt-5"placeholder,{"id":"gpt-5"placeholder,{"name":"o3"placeholder]placeholder`)),
placeholderplaceholder
	svc := &AccountTestService{
		httpUpstream: upstream,
		cfg:          upstreamModelSyncTestConfig(),
placeholder

	models, err := svc.FetchUpstreamSupportedModels(context.Background(), &Account{
		ID:       7,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
placeholder
			"api_key":  "openai-key",
			"base_url": "https://openai.example.com/v1",
	placeholder,
placeholder)
placeholder
	require.Equal(t, []string{"gpt-5", "o3"placeholder, models)
	require.Equal(t, "https://openai.example.com/v1/models", upstream.lastReq.URL.String())
	require.Equal(t, "Bearer openai-key", upstream.lastReq.Header.Get("Authorization"))
placeholder

func TestFetchUpstreamSupportedModelsParsesGrokAPIKeyResponse(t *testing.T) {
	t.Parallel()

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"data":[{"id":"grok-4.5"placeholder,{"id":"grok-4.5"placeholder,{"id":"grok-imagine"placeholder]placeholder`)),
placeholderplaceholder
	svc := &AccountTestService{
		httpUpstream: upstream,
		cfg:          upstreamModelSyncTestConfig(),
placeholder

	models, err := svc.FetchUpstreamSupportedModels(context.Background(), &Account{
		ID:       9,
		Platform: PlatformGrok,
		Type:     AccountTypeAPIKey,
placeholder
			"api_key":  "xai-key",
			"base_url": "https://xai.example.com/v1",
	placeholder,
placeholder)
placeholder
	require.Equal(t, []string{"grok-4.5", "grok-imagine"placeholder, models)
	require.Equal(t, "https://xai.example.com/v1/models", upstream.lastReq.URL.String())
	require.Equal(t, "Bearer xai-key", upstream.lastReq.Header.Get("Authorization"))
placeholder

func TestFetchUpstreamSupportedModelsParsesGrokOAuthResponse(t *testing.T) {
	t.Parallel()

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"data":[{"model":"grok-4.5"placeholder,{"model":"grok-4.5"placeholder,{"modelId":"grok-build-0.1"placeholder]placeholder`)),
placeholderplaceholder
	svc := &AccountTestService{
		httpUpstream:      upstream,
		cfg:               upstreamModelSyncTestConfig(),
		grokTokenProvider: NewGrokTokenProvider(nil, nil),
placeholder

	models, err := svc.FetchUpstreamSupportedModels(context.Background(), grokOAuthModelSyncTestAccount(""))
placeholder
	require.Equal(t, []string{"grok-4.5", "grok-build-0.1"placeholder, models)
	require.Equal(t, "https://cli-chat-proxy.grok.com/v1/models", upstream.lastReq.URL.String())
	require.Equal(t, "Bearer oauth-access-token", upstream.lastReq.Header.Get("Authorization"))
	require.Equal(t, grokCLIVersion, upstream.lastReq.Header.Get("X-Grok-Client-Version"))
	require.Equal(t, "interactive", upstream.lastReq.Header.Get("X-Grok-Client-Mode"))
	require.Equal(t, "grok-user-id", upstream.lastReq.Header.Get("X-UserID"))
	require.Equal(t, "grok-user@example.com", upstream.lastReq.Header.Get("X-Email"))
placeholder

func TestBuildUpstreamModelsRequestGrokOAuthDoesNotSendIdentityToCustomBase(t *testing.T) {
	t.Parallel()

	svc := &AccountTestService{
		cfg:               upstreamModelSyncTestConfig(),
		grokTokenProvider: NewGrokTokenProvider(nil, nil),
placeholder
	req, err := svc.buildUpstreamModelsRequest(context.Background(), grokOAuthModelSyncTestAccount("https://relay.example/v1"))
placeholder
	require.Equal(t, "https://relay.example/v1/models", req.URL.String())
	require.Empty(t, req.Header.Get("X-UserID"))
	require.Empty(t, req.Header.Get("X-Email"))
placeholder

func TestFetchUpstreamSupportedModelsDoesNotExposeUpstreamBody(t *testing.T) {
	t.Parallel()

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusBadGateway,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"error":"SECRET_TOKEN should not be exposed"placeholder`)),
placeholderplaceholder
	svc := &AccountTestService{
		httpUpstream: upstream,
		cfg:          upstreamModelSyncTestConfig(),
placeholder

	_, err := svc.FetchUpstreamSupportedModels(context.Background(), &Account{
		ID:       8,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
placeholder
			"api_key":  "openai-key",
			"base_url": "https://openai.example.com/v1",
	placeholder,
placeholder)
placeholder
	require.NotContains(t, err.Error(), "SECRET_TOKEN")

	var syncErr *UpstreamModelSyncError
	require.True(t, errors.As(err, &syncErr))
	require.Equal(t, UpstreamModelSyncErrorUpstream, syncErr.Kind)
	require.NotContains(t, syncErr.SafeMessage(), "SECRET_TOKEN")
	require.Contains(t, syncErr.SafeMessage(), "HTTP 502")
placeholder
