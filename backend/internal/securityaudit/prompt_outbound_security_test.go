package securityaudit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type staticResolver struct{ addresses []netip.Addr placeholder

func (r staticResolver) LookupNetIP(context.Context, string, string) ([]netip.Addr, error) {
	return r.addresses, nil
placeholder

func TestNormalizeBaseURLSecurity(t *testing.T) {
	allowed := []string{"https://guard.example.com", "https://guard.example.com/v1", "http://127.0.0.1:8080", "http://10.0.0.8:8080"placeholder
	for _, raw := range allowed {
		_, err := NormalizeBaseURL(raw)
		require.NoError(t, err, raw)
placeholder
	blocked := []string{
		"ftp://guard.example.com", "http://guard.example.com", "https://user:pass@guard.example.com",
		"https://guard.example.com?q=secret", "https://guard.example.com/#fragment", "http://169.254.169.254",
		"https://metadata.google.internal", "https://0.0.0.0", "https://224.0.0.1", "https://192.0.2.1",
		"https://[::]", "https://[fe80::1]", "https://[ff02::1]", "https://[2001:db8::1]",
placeholder
	for _, raw := range blocked {
		_, err := NormalizeBaseURL(raw)
		require.Error(t, err, raw)
placeholder
	url, err := ChatCompletionsURL("https://guard.example.com/v1")
placeholder
	require.Equal(t, "https://guard.example.com/v1/chat/completions", url)
placeholder

func TestSecureDialRejectsDNSRebindingToPrivateAddress(t *testing.T) {
	dial := secureDialContext(nil, staticResolver{addresses: []netip.Addr{netip.MustParseAddr("127.0.0.1")placeholderplaceholder, false)
	_, err := dial(context.Background(), "tcp", "guard.example.com:443")
placeholder
placeholder

func TestSecureHTTPClientDoesNotBypassDestinationValidationThroughEnvironmentProxy(t *testing.T) {
	client, err := NewSecureHTTPClient(ActiveEndpoint{BaseURL: "https://guard.example.com", TimeoutMS: 1000placeholder)
placeholder
	transport, ok := client.Transport.(*http.Transport)
	require.True(t, ok)
	require.Nil(t, transport.Proxy)
placeholder

func TestOpenAICompatibleScannerRequestContract(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v1/chat/completions", r.URL.Path)
		require.Equal(t, "Bearer token", r.Header.Get("Authorization"))
		var payload map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, DefaultGuardModel, payload["model"])
		require.Equal(t, float64(0), payload["temperature"])
		require.Equal(t, float64(64), payload["max_tokens"])
		require.Equal(t, float64(42), payload["seed"])
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"Safety: Safe\nCategories: None"placeholderplaceholder]placeholder`))
placeholder))
	defer server.Close()
	scanner := NewOpenAICompatibleScanner()
	result, err := scanner.Scan(context.Background(), ActiveEndpoint{ID: "one", BaseURL: server.URL, Model: DefaultGuardModel, Token: "token", TimeoutMS: 1000placeholder, "hello", AllScannerIDs)
placeholder
	require.Equal(t, EventPass, result.Decision)
placeholder

func TestOpenAICompatibleScannerRejectsRedirectAndOversize(t *testing.T) {
	redirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://127.0.0.1/other", http.StatusFound)
placeholder))
	defer redirect.Close()
	_, err := NewOpenAICompatibleScanner().Scan(context.Background(), ActiveEndpoint{ID: "redirect", BaseURL: redirect.URL, Model: DefaultGuardModel, TimeoutMS: 1000placeholder, "hello", AllScannerIDs)
placeholder
	oversize := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(strings.Repeat("x", int(maxGuardResponseBytes)+1)))
placeholder))
	defer oversize.Close()
	_, err = NewOpenAICompatibleScanner().Scan(context.Background(), ActiveEndpoint{ID: "large", BaseURL: oversize.URL, Model: DefaultGuardModel, TimeoutMS: 1000placeholder, "hello", AllScannerIDs)
placeholder
placeholder

func TestOpenAICompatibleScannerClassifiesHTTPConnectionAndTimeoutFailures(t *testing.T) {
	tests := []struct {
		name      string
		status    int
		retryable bool
placeholder{
		{name: "authentication", status: http.StatusUnauthorized, retryable: falseplaceholder,
		{name: "forbidden", status: http.StatusForbidden, retryable: falseplaceholder,
		{name: "rate limited", status: http.StatusTooManyRequests, retryable: trueplaceholder,
		{name: "server failure", status: http.StatusBadGateway, retryable: trueplaceholder,
		{name: "other client error", status: http.StatusBadRequest, retryable: falseplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
		placeholder))
			defer server.Close()
			_, err := NewOpenAICompatibleScanner().Scan(context.Background(), ActiveEndpoint{ID: "status", BaseURL: server.URL, Model: DefaultGuardModel, TimeoutMS: 1000placeholder, "hello", AllScannerIDs)
			var guardErr *GuardError
			require.ErrorAs(t, err, &guardErr)
			require.Equal(t, ErrorCodeUnavailable, guardErr.Code)
			require.Equal(t, tt.status, guardErr.HTTPStatus)
			require.Equal(t, tt.retryable, guardErr.Retryable)
			require.NotContains(t, err.Error(), server.URL)
	placeholder)
placeholder

	closed := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {placeholder))
	closedURL := closed.URL
	closed.Close()
	_, err := NewOpenAICompatibleScanner().Scan(context.Background(), ActiveEndpoint{ID: "closed", BaseURL: closedURL, Model: DefaultGuardModel, TimeoutMS: 100placeholder, "hello", AllScannerIDs)
	var connectionErr *GuardError
	require.ErrorAs(t, err, &connectionErr)
	require.True(t, connectionErr.Retryable)

	timeout := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
placeholder))
	defer timeout.Close()
	_, err = NewOpenAICompatibleScanner().Scan(context.Background(), ActiveEndpoint{ID: "timeout", BaseURL: timeout.URL, Model: DefaultGuardModel, TimeoutMS: 20placeholder, "hello", AllScannerIDs)
	var timeoutErr *GuardError
	require.ErrorAs(t, err, &timeoutErr)
	require.True(t, timeoutErr.Retryable)
	require.True(t, timeoutErr.Timeout)
placeholder

func TestPromptAuditProbeModelsFallbackAndResponseSafety(t *testing.T) {
	t.Run("models contains configured model", func(t *testing.T) {
		var chatCalls atomic.Int64
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "Bearer temporary-token", r.Header.Get("Authorization"))
			if r.URL.Path == "/v1/models" {
				_, _ = w.Write([]byte(`{"data":[{"id":"` + DefaultGuardModel + `"placeholder]placeholder`))
				return
		placeholder
			chatCalls.Add(1)
	placeholder))
		defer server.Close()
		result := newProbeTestService().Probe(context.Background(), ProbeRequest{Endpoint: probeEndpoint(server.URL, "temporary-token")placeholder)
		require.True(t, result.OK)
		require.True(t, result.TokenApplied)
		require.Equal(t, http.StatusOK, result.HTTPStatus)
		require.Zero(t, chatCalls.Load())
placeholder)

	t.Run("invalid models response performs real guard fallback", func(t *testing.T) {
		var chatCalls atomic.Int64
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/v1/models" {
				_, _ = w.Write([]byte(`{"unexpected":trueplaceholder`))
				return
		placeholder
			chatCalls.Add(1)
			_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"Safety: Safe\nCategories: None"placeholderplaceholder]placeholder`))
	placeholder))
		defer server.Close()
		result := newProbeTestService().Probe(context.Background(), ProbeRequest{Endpoint: probeEndpoint(server.URL, "temporary-token")placeholder)
		require.True(t, result.OK)
		require.Equal(t, int64(1), chatCalls.Load())
placeholder)

	t.Run("fallback authentication failure is stable", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/v1/models" {
				w.WriteHeader(http.StatusNotFound)
				return
		placeholder
			w.WriteHeader(http.StatusUnauthorized)
	placeholder))
		defer server.Close()
		result := newProbeTestService().Probe(context.Background(), ProbeRequest{Endpoint: probeEndpoint(server.URL, "temporary-token")placeholder)
		require.False(t, result.OK)
		require.Equal(t, ErrorCodeUnavailable, result.ErrorCode)
		require.Equal(t, http.StatusUnauthorized, result.HTTPStatus)
		require.False(t, result.Retryable)
placeholder)

	t.Run("oversized models response is rejected without fallback", func(t *testing.T) {
		var chatCalls atomic.Int64
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/v1/models" {
				chatCalls.Add(1)
		placeholder
			_, _ = w.Write([]byte(strings.Repeat("x", int(maxGuardResponseBytes)+1)))
	placeholder))
		defer server.Close()
		result := newProbeTestService().Probe(context.Background(), ProbeRequest{Endpoint: probeEndpoint(server.URL, "temporary-token")placeholder)
		require.False(t, result.OK)
		require.Equal(t, "response_too_large", result.ErrorCode)
		require.Zero(t, chatCalls.Load())
placeholder)
placeholder

func newProbeTestService() *PromptService {
	return &PromptService{
		config: &ConfigManager{placeholder, scanner: NewOpenAICompatibleScanner(), clock: realClock{placeholder,
		probes: map[string]ProbeResult{placeholder,
placeholder
placeholder

func probeEndpoint(baseURL, token string) UpdateEndpoint {
	return UpdateEndpoint{
		ID: "probe-one", Name: "Probe One", Protocol: "openai_compatible", BaseURL: baseURL,
		Model: DefaultGuardModel, Token: token, TimeoutMS: 1000, InputLimit: 1024, Enabled: true,
placeholder
placeholder
