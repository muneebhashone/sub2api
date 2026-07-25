package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	coderws "github.com/coder/websocket"
	"github.com/stretchr/testify/require"
)

type liveHTTPUpstreamStub struct {
	request *http.Request
	body    []byte
placeholder

func (s *liveHTTPUpstreamStub) Do(
	request *http.Request,
	_ string,
	_ int64,
	_ int,
) (*http.Response, error) {
	s.request = request
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
placeholder
	s.body = body
	return &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Location": {"/backend-api/codex/call_test"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader("v=0\r\n")),
placeholder, nil
placeholder

func (s *liveHTTPUpstreamStub) DoWithTLS(
	request *http.Request,
	proxyURL string,
	accountID int64,
	accountConcurrency int,
	_ *tlsfingerprint.Profile,
) (*http.Response, error) {
	return s.Do(request, proxyURL, accountID, accountConcurrency)
placeholder

func TestLiveCapabilityOnlyAllowsOpenAIOAuth(t *testing.T) {
	require.True(t, (&Account{Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder).SupportsOpenAIEndpointCapability(OpenAIEndpointCapabilityLive))
	require.False(t, (&Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder).SupportsOpenAIEndpointCapability(OpenAIEndpointCapabilityLive))
	require.False(t, (&Account{Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder).SupportsOpenAIEndpointCapability(OpenAIEndpointCapabilityLive))
	require.False(t, (&Account{
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			openAIAuthModeCredentialKey: OpenAIAuthModePersonalAccessToken,
	placeholder,
placeholder).SupportsOpenAIEndpointCapability(OpenAIEndpointCapabilityLive))
	require.False(t, (&Account{
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			openAIAuthModeCredentialKey: OpenAIAuthModeAgentIdentity,
	placeholder,
placeholder).SupportsOpenAIEndpointCapability(OpenAIEndpointCapabilityLive))
placeholder

func TestValidateLiveCallRequestDoesNotRequireDelegation(t *testing.T) {
	request := &LiveCallRequest{
		SDP:     "v=0\r\n",
		Session: json.RawMessage(`{"model":"gpt-live-test","instructions":"hello"placeholder`),
placeholder
	require.NoError(t, ValidateLiveCallRequest(request))
	require.NotContains(t, string(request.Session), "delegation")
placeholder

func TestCreateUpstreamLiveCallPreservesSession(t *testing.T) {
	upstream := &liveHTTPUpstreamStub{placeholder
	service := &OpenAIGatewayService{
		cfg:          &config.Config{placeholder,
		httpUpstream: upstream,
placeholder
	account := &Account{
		ID:          7,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 2,
placeholder
			"access_token":       "test-access-token",
			"chatgpt_account_id": "acct_test",
	placeholder,
placeholder
	session := json.RawMessage(`{
		"model":"gpt-live-test",
		"delegation":{"type":"client"placeholder,
		"custom":{"keep":trueplaceholder
placeholder`)

	created, err := service.createUpstreamLiveCall(context.Background(), account, &LiveCallRequest{
		SDP:     "v=offer\r\n",
		Session: session,
placeholder)
placeholder
	require.Equal(t, "call_test", created.CallID)
	require.Equal(t, []byte("v=0\r\n"), created.SDP)

	var forwarded struct {
		SDP     string          `json:"sdp"`
		Session json.RawMessage `json:"session"`
placeholder
	require.NoError(t, json.Unmarshal(upstream.body, &forwarded))
	require.Equal(t, "v=offer\r\n", forwarded.SDP)
	require.JSONEq(t, string(session), string(forwarded.Session))
	require.Equal(t, "Bearer test-access-token", upstream.request.Header.Get("Authorization"))
	require.Equal(t, "acct_test", upstream.request.Header.Get("Chatgpt-Account-Id"))
	require.Equal(t, "quicksilver=v2", upstream.request.Header.Get("OpenAI-Alpha"))
	require.NotEmpty(t, upstream.request.Header.Get("Session-Id"))
	require.NotEmpty(t, upstream.request.Header.Get("Thread-Id"))
	require.Empty(t, upstream.request.Header.Get("OpenAI-Beta"))
	require.Equal(t, HTTPUpstreamProfileOpenAI, HTTPUpstreamProfileFromContext(upstream.request.Context()))
	require.True(t, HTTPUpstreamRedirectsDisabled(upstream.request.Context()))
placeholder

func TestLiveMaxSessionDurationDefaultsAndOverrides(t *testing.T) {
	require.Equal(t, defaultLiveMaxSessionDuration, (&OpenAIGatewayService{placeholder).liveMaxSessionDuration())
	require.Equal(
		t,
		90*time.Second,
		(&OpenAIGatewayService{cfg: &config.Config{
			Gateway: config.GatewayConfig{
				Live: config.GatewayLiveConfig{MaxSessionDurationSeconds: 90placeholder,
		placeholder,
	placeholderplaceholder).liveMaxSessionDuration(),
	)
placeholder

func TestLiveSidebandNormalCloseEndsCall(t *testing.T) {
	normalClose := coderws.CloseError{Code: coderws.StatusNormalClosureplaceholder
	require.ErrorIs(t, liveSidebandReadError(normalClose), ErrLiveCallNotFound)

	abnormalClose := coderws.CloseError{Code: coderws.StatusInternalErrorplaceholder
	require.Equal(t, abnormalClose, liveSidebandReadError(abnormalClose))
placeholder

func TestLiveCreateFailoverUsesExistingOpenAIPolicy(t *testing.T) {
	service := &OpenAIGatewayService{placeholder
	require.False(t, service.shouldFailoverLiveCreateError(&UpstreamFailoverError{
		StatusCode:   http.StatusBadRequest,
		ResponseBody: []byte(`{"error":{"message":"invalid session"placeholderplaceholder`),
placeholder))
	require.True(t, service.shouldFailoverLiveCreateError(&UpstreamFailoverError{
		StatusCode: http.StatusForbidden,
placeholder))
	require.True(t, service.shouldFailoverLiveCreateError(&UpstreamFailoverError{
		StatusCode: http.StatusBadGateway,
placeholder))
	require.True(t, service.shouldFailoverLiveCreateError(errors.New("transport failed")))
placeholder

func TestLiveCallIDFromLocation(t *testing.T) {
	callID, err := liveCallIDFromLocation("https://chatgpt.com/backend-api/codex/call_123?intent=quicksilver")
placeholder
	require.Equal(t, "call_123", callID)

	callID, err = liveCallIDFromLocation("/backend-api/codex/call_456")
placeholder
	require.Equal(t, "call_456", callID)
placeholder

func TestRequestTypeLive(t *testing.T) {
	require.True(t, RequestTypeLive.IsValid())
	require.Equal(t, "live", RequestTypeLive.String())
	parsed, err := ParseUsageRequestType("live")
placeholder
	require.Equal(t, RequestTypeLive, parsed)
placeholder
