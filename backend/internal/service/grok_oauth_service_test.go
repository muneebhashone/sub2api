//go:build unit

package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/stretchr/testify/require"
)

type grokOAuthClientStub struct {
	refreshResponse     *xai.TokenResponse
	ssoResponse         *xai.TokenResponse
	loginResult         *GrokPasswordLoginResult
	loginEmail          string
	loginPassword       string
	exchangeCalls       int
	exchangeRedirectURI string
placeholder

func (s *grokOAuthClientStub) ExchangeCode(_ context.Context, _, _, redirectURI, _, _ string) (*xai.TokenResponse, error) {
	s.exchangeCalls++
	s.exchangeRedirectURI = redirectURI
	return &xai.TokenResponse{AccessToken: "access-token"placeholder, nil
placeholder

func (s *grokOAuthClientStub) RefreshToken(context.Context, string, string, string) (*xai.TokenResponse, error) {
	return s.refreshResponse, nil
placeholder

func (s *grokOAuthClientStub) LoginWithPassword(_ context.Context, email, password, _ string) (*GrokPasswordLoginResult, error) {
	s.loginEmail = email
	s.loginPassword = password
	return s.loginResult, nil
placeholder

func (s *grokOAuthClientStub) ConvertSSOToBuild(context.Context, string, string) (*xai.TokenResponse, error) {
	return s.ssoResponse, nil
placeholder

func TestGrokOAuthServiceRefreshTokenPreservesOriginalRefreshTokenWhenNotRotated(t *testing.T) {
	svc := NewGrokOAuthService(nil, &grokOAuthClientStub{
		refreshResponse: &xai.TokenResponse{
			AccessToken: "new-access-token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
	placeholder,
placeholder)
	defer svc.Stop()

	info, err := svc.RefreshToken(context.Background(), "original-refresh-token", "", "client-id")
placeholder
	require.Equal(t, "new-access-token", info.AccessToken)
	require.Equal(t, "original-refresh-token", info.RefreshToken)
	require.Equal(t, "client-id", info.ClientID)
placeholder

func TestGrokOAuthServiceRefreshTokenRejectsEmptyUpstreamResponse(t *testing.T) {
	svc := NewGrokOAuthService(nil, &grokOAuthClientStub{placeholder)
	defer svc.Stop()

	require.NotPanics(t, func() {
		info, err := svc.RefreshToken(context.Background(), "refresh-token", "", "client-id")
		require.Nil(t, info)
	placeholder
		require.Contains(t, err.Error(), "GROK_OAUTH_INVALID_TOKEN_RESPONSE")
placeholder)
placeholder

func TestGrokOAuthServiceExchangeCodeConsumesOnlyAfterValidation(t *testing.T) {
	client := &grokOAuthClientStub{placeholder
	svc := NewGrokOAuthService(nil, client)
	defer svc.Stop()

	auth, err := svc.GenerateAuthURL(context.Background(), nil, "")
placeholder

	_, err = svc.ExchangeCode(context.Background(), &GrokExchangeCodeInput{
		SessionID: auth.SessionID,
		Code:      "http://127.0.0.1:56121/callback?code=code-without-state",
placeholder)
placeholder
	require.Contains(t, err.Error(), "GROK_OAUTH_STATE_REQUIRED")
	require.Zero(t, client.exchangeCalls)

	_, err = svc.ExchangeCode(context.Background(), &GrokExchangeCodeInput{
		SessionID: auth.SessionID,
		Code:      "code-with-state",
		State:     auth.State,
placeholder)
placeholder
	require.Equal(t, 1, client.exchangeCalls)

	_, err = svc.ExchangeCode(context.Background(), &GrokExchangeCodeInput{
		SessionID: auth.SessionID,
		Code:      "replayed-code",
		State:     auth.State,
placeholder)
placeholder
	require.Contains(t, err.Error(), "GROK_OAUTH_SESSION_NOT_FOUND")
	require.Equal(t, 1, client.exchangeCalls)
placeholder

func TestGrokOAuthServiceExchangeCodeRejectsMissingClientWithoutConsumingSession(t *testing.T) {
	svc := NewGrokOAuthService(nil, nil)
	defer svc.Stop()
	auth, err := svc.GenerateAuthURL(context.Background(), nil, "")
placeholder

	_, err = svc.ExchangeCode(context.Background(), &GrokExchangeCodeInput{
		SessionID: auth.SessionID,
		Code:      "code",
		State:     auth.State,
placeholder)
placeholder
	require.Contains(t, err.Error(), "GROK_OAUTH_CLIENT_NOT_CONFIGURED")
	_, ok := svc.sessionStore.Get(auth.SessionID)
	require.True(t, ok)
placeholder

func TestGrokOAuthServiceExchangeCodeRequiresStateForBareCode(t *testing.T) {
	client := &grokOAuthClientStub{placeholder
	svc := NewGrokOAuthService(nil, client)
	defer svc.Stop()
	auth, err := svc.GenerateAuthURL(context.Background(), nil, "")
placeholder

	_, err = svc.ExchangeCode(context.Background(), &GrokExchangeCodeInput{
		SessionID: auth.SessionID,
		Code:      "bare-authorization-code",
placeholder)
placeholder
	require.Contains(t, err.Error(), "GROK_OAUTH_STATE_REQUIRED")
	require.Zero(t, client.exchangeCalls)
	_, ok := svc.sessionStore.Get(auth.SessionID)
	require.True(t, ok)
placeholder

func TestGrokOAuthServiceExchangeCodeRejectsRedirectURIOverride(t *testing.T) {
	client := &grokOAuthClientStub{placeholder
	svc := NewGrokOAuthService(nil, client)
	defer svc.Stop()
	auth, err := svc.GenerateAuthURL(context.Background(), nil, "")
placeholder

	_, err = svc.ExchangeCode(context.Background(), &GrokExchangeCodeInput{
		SessionID:   auth.SessionID,
		Code:        "authorization-code",
		State:       auth.State,
		RedirectURI: "http://127.0.0.1:9999/callback",
placeholder)
placeholder
	require.Contains(t, err.Error(), "GROK_OAUTH_REDIRECT_URI_MISMATCH")
	require.Zero(t, client.exchangeCalls)

	_, err = svc.ExchangeCode(context.Background(), &GrokExchangeCodeInput{
		SessionID:   auth.SessionID,
		Code:        "authorization-code",
		State:       auth.State,
		RedirectURI: xai.DefaultRedirectURI,
placeholder)
placeholder
	require.Equal(t, xai.DefaultRedirectURI, client.exchangeRedirectURI)
placeholder

func TestGrokOAuthServiceExternalFlowsRejectMissingClient(t *testing.T) {
	svc := NewGrokOAuthService(nil, nil)
	defer svc.Stop()

	_, err := svc.RefreshToken(context.Background(), "refresh-token", "", "")
placeholder
	require.Contains(t, err.Error(), "GROK_OAUTH_CLIENT_NOT_CONFIGURED")

	_, err = svc.ValidateSSOToken(context.Background(), "sso-token", nil)
placeholder
	require.Contains(t, err.Error(), "GROK_OAUTH_CLIENT_NOT_CONFIGURED")
placeholder

func TestGrokOAuthServiceBuildAccountCredentialsDefaultsToSubscriptionProxy(t *testing.T) {
	svc := NewGrokOAuthService(nil, &grokOAuthClientStub{placeholder)
	defer svc.Stop()

	credentials := svc.BuildAccountCredentials(&GrokTokenInfo{
		AccessToken: "access-token",
		ExpiresAt:   time.Now().Add(time.Hour).Unix(),
placeholder)

	require.Equal(t, xai.DefaultCLIBaseURL, credentials["base_url"])
placeholder

func TestGrokOAuthServiceConvertFromSSOExtractsBuildClaims(t *testing.T) {
	svc := NewGrokOAuthService(nil, &grokOAuthClientStub{
		ssoResponse: &xai.TokenResponse{
			AccessToken:  makeGrokOAuthJWT(map[string]any{"sub": "user-sub", "team_id": "team-1"placeholder),
			RefreshToken: "refresh-token",
			IDToken:      makeGrokOAuthJWT(map[string]any{"email": "user@example.com"placeholder),
			ExpiresIn:    3600,
	placeholder,
placeholder)
	defer svc.Stop()

	info, err := svc.ConvertFromSSO(context.Background(), "sso-token", nil)
placeholder
	require.Equal(t, "user@example.com", info.Email)
	require.Equal(t, "user-sub", info.Subject)
	require.Equal(t, "team-1", info.TeamID)

	credentials := svc.BuildAccountCredentials(info)
	require.Equal(t, "user@example.com", credentials["email"])
	require.Equal(t, "user-sub", credentials["sub"])
	require.Equal(t, "team-1", credentials["team_id"])
	require.NotContains(t, credentials, "sso_token")
placeholder

func TestGrokOAuthServiceValidateSSOTokenReturnsOAuthTokensWithoutPersistingSSO(t *testing.T) {
	svc := NewGrokOAuthService(nil, &grokOAuthClientStub{
		ssoResponse: &xai.TokenResponse{
			AccessToken:  "access-from-sso",
			RefreshToken: "refresh-from-sso",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
	placeholder,
placeholder)
	defer svc.Stop()

	info, err := svc.ValidateSSOToken(context.Background(), "sso-token", nil)
placeholder
	require.Equal(t, "access-from-sso", info.AccessToken)
	require.Equal(t, "refresh-from-sso", info.RefreshToken)

	creds := svc.BuildAccountCredentials(info)
	require.NotContains(t, creds, "sso_token")
	require.NotContains(t, creds, "password")
placeholder

func TestGrokOAuthServiceAuthorizePasswordUsesLoginThenSSOAuthorize(t *testing.T) {
	client := &grokOAuthClientStub{
		loginResult: &GrokPasswordLoginResult{
			Email:    "user@example.com",
			SSOToken: "password-derived-sso",
	placeholder,
		ssoResponse: &xai.TokenResponse{
			AccessToken:  "access-from-password",
			RefreshToken: "refresh-from-password",
			ExpiresIn:    3600,
	placeholder,
placeholder
	cfg := &config.Config{placeholder
	cfg.Gateway.Grok.PasswordAuthEnabled = true
	svc := NewGrokOAuthService(nil, client, cfg)
	defer svc.Stop()

	require.True(t, svc.GetCapabilities().PasswordAuthEnabled)
	info, err := svc.AuthorizePassword(context.Background(), " user@example.com ", "  super-secret  ", nil)
placeholder
	require.Equal(t, "user@example.com", info.Email)
	require.Equal(t, "access-from-password", info.AccessToken)
	creds := svc.BuildAccountCredentials(info)
	require.NotContains(t, creds, "password")
	require.NotContains(t, creds, "sso_token")
	require.Equal(t, "user@example.com", client.loginEmail)
	require.Equal(t, "  super-secret  ", client.loginPassword)
placeholder

func TestGrokOAuthServiceAuthorizePasswordDisabledByDefault(t *testing.T) {
	client := &grokOAuthClientStub{placeholder
	svc := NewGrokOAuthService(nil, client)
	defer svc.Stop()

	require.False(t, svc.GetCapabilities().PasswordAuthEnabled)
	_, err := svc.AuthorizePassword(context.Background(), "user@example.com", "secret", nil)
placeholder
	require.Contains(t, err.Error(), "GROK_OAUTH_PASSWORD_AUTH_DISABLED")
	require.Empty(t, client.loginEmail)
placeholder

func makeGrokOAuthJWT(claims map[string]any) string {
	payload, _ := json.Marshal(claims)
	return "header." + base64.RawURLEncoding.EncodeToString(payload) + ".signature"
placeholder
