package handler

import (
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func TestSanitizeFrontendRedirectPath(t *testing.T) {
	require.Equal(t, "/dashboard", sanitizeFrontendRedirectPath("/dashboard"))
	require.Equal(t, "/dashboard", sanitizeFrontendRedirectPath(" /dashboard "))
	require.Equal(t, "", sanitizeFrontendRedirectPath("dashboard"))
	require.Equal(t, "", sanitizeFrontendRedirectPath("//evil.com"))
	require.Equal(t, "", sanitizeFrontendRedirectPath("https://evil.com"))
	require.Equal(t, "", sanitizeFrontendRedirectPath("/\nfoo"))

	long := "/" + strings.Repeat("a", linuxDoOAuthMaxRedirectLen)
	require.Equal(t, "", sanitizeFrontendRedirectPath(long))
placeholder

func TestBuildBearerAuthorization(t *testing.T) {
	auth, err := buildBearerAuthorization("", "token123")
placeholder
	require.Equal(t, "Bearer token123", auth)

	auth, err = buildBearerAuthorization("bearer", "token123")
placeholder
	require.Equal(t, "Bearer token123", auth)

	_, err = buildBearerAuthorization("MAC", "token123")
placeholder

	_, err = buildBearerAuthorization("Bearer", "token 123")
placeholder
placeholder

func TestLinuxDoParseUserInfoParsesIDAndUsername(t *testing.T) {
	cfg := config.LinuxDoConnectConfig{
		UserInfoURL: "https://connect.linux.do/api/user",
placeholder

	email, username, subject, err := linuxDoParseUserInfo(`{"id":123,"username":"alice"placeholder`, cfg)
placeholder
	require.Equal(t, "123", subject)
	require.Equal(t, "alice", username)
	require.Equal(t, "linuxdo-123@linuxdo-connect.invalid", email)
placeholder

func TestLinuxDoParseUserInfoDefaultsUsername(t *testing.T) {
	cfg := config.LinuxDoConnectConfig{
		UserInfoURL: "https://connect.linux.do/api/user",
placeholder

	email, username, subject, err := linuxDoParseUserInfo(`{"id":"123"placeholder`, cfg)
placeholder
	require.Equal(t, "123", subject)
	require.Equal(t, "linuxdo_123", username)
	require.Equal(t, "linuxdo-123@linuxdo-connect.invalid", email)
placeholder

func TestLinuxDoParseUserInfoRejectsUnsafeSubject(t *testing.T) {
	cfg := config.LinuxDoConnectConfig{
		UserInfoURL: "https://connect.linux.do/api/user",
placeholder

	_, _, _, err := linuxDoParseUserInfo(`{"id":"123@456"placeholder`, cfg)
placeholder

	tooLong := strings.Repeat("a", linuxDoOAuthMaxSubjectLen+1)
	_, _, _, err = linuxDoParseUserInfo(`{"id":"`+tooLong+`"placeholder`, cfg)
placeholder
placeholder

func TestParseOAuthProviderErrorJSON(t *testing.T) {
	code, desc := parseOAuthProviderError(`{"error":"invalid_client","error_description":"bad secret"placeholder`)
	require.Equal(t, "invalid_client", code)
	require.Equal(t, "bad secret", desc)
placeholder

func TestParseOAuthProviderErrorForm(t *testing.T) {
	code, desc := parseOAuthProviderError("error=invalid_request&error_description=Missing+code_verifier")
	require.Equal(t, "invalid_request", code)
	require.Equal(t, "Missing code_verifier", desc)
placeholder

func TestParseLinuxDoTokenResponseJSON(t *testing.T) {
	token, ok := parseLinuxDoTokenResponse(`{"access_token":"t1","token_type":"Bearer","expires_in":3600,"scope":"user"placeholder`)
	require.True(t, ok)
	require.Equal(t, "t1", token.AccessToken)
	require.Equal(t, "Bearer", token.TokenType)
	require.Equal(t, int64(3600), token.ExpiresIn)
	require.Equal(t, "user", token.Scope)
placeholder

func TestParseLinuxDoTokenResponseForm(t *testing.T) {
	token, ok := parseLinuxDoTokenResponse("access_token=t2&token_type=bearer&expires_in=60")
	require.True(t, ok)
	require.Equal(t, "t2", token.AccessToken)
	require.Equal(t, "bearer", token.TokenType)
	require.Equal(t, int64(60), token.ExpiresIn)
placeholder

func TestSingleLineStripsWhitespace(t *testing.T) {
	require.Equal(t, "hello world", singleLine("hello\r\nworld"))
	require.Equal(t, "", singleLine("\n\t\r"))
placeholder
