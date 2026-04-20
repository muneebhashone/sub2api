package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/ent/authidentity"
	"github.com/Wei-Shaw/sub2api/ent/identityadoptiondecision"
	"github.com/Wei-Shaw/sub2api/ent/pendingauthsession"
	dbuser "github.com/Wei-Shaw/sub2api/ent/user"
	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
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

	email, username, subject, displayName, avatarURL, err := linuxDoParseUserInfo(`{"id":123,"username":"alice","name":"Alice","avatar_url":"https://cdn.example/avatar.png"placeholder`, cfg)
placeholder
	require.Equal(t, "123", subject)
	require.Equal(t, "alice", username)
	require.Equal(t, "linuxdo-123@linuxdo-connect.invalid", email)
	require.Equal(t, "Alice", displayName)
	require.Equal(t, "https://cdn.example/avatar.png", avatarURL)
placeholder

func TestLinuxDoParseUserInfoDefaultsUsername(t *testing.T) {
	cfg := config.LinuxDoConnectConfig{
		UserInfoURL: "https://connect.linux.do/api/user",
placeholder

	email, username, subject, displayName, avatarURL, err := linuxDoParseUserInfo(`{"id":"123"placeholder`, cfg)
placeholder
	require.Equal(t, "123", subject)
	require.Equal(t, "linuxdo_123", username)
	require.Equal(t, "linuxdo-123@linuxdo-connect.invalid", email)
	require.Equal(t, "linuxdo_123", displayName)
	require.Equal(t, "", avatarURL)
placeholder

func TestLinuxDoParseUserInfoRejectsUnsafeSubject(t *testing.T) {
	cfg := config.LinuxDoConnectConfig{
		UserInfoURL: "https://connect.linux.do/api/user",
placeholder

	_, _, _, _, _, err := linuxDoParseUserInfo(`{"id":"123@456"placeholder`, cfg)
placeholder

	tooLong := strings.Repeat("a", linuxDoOAuthMaxSubjectLen+1)
	_, _, _, _, _, err = linuxDoParseUserInfo(`{"id":"`+tooLong+`"placeholder`, cfg)
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

func TestCompleteLinuxDoOAuthRegistrationAppliesPendingAdoptionDecision(t *testing.T) {
	handler, client := newOAuthPendingFlowTestHandler(t, false)
	ctx := context.Background()

	session, err := client.PendingAuthSession.Create().
		SetSessionToken("linuxdo-complete-session").
		SetIntent("login").
		SetProviderType("linuxdo").
		SetProviderKey("linuxdo").
		SetProviderSubject("linuxdo-subject-1").
		SetResolvedEmail("linuxdo-subject-1@linuxdo-connect.invalid").
		SetBrowserSessionKey("linuxdo-browser").
		SetUpstreamIdentityClaims(map[string]any{
			"username":               "linuxdo_user",
			"suggested_display_name": "LinuxDo Display",
			"suggested_avatar_url":   "https://cdn.example/linuxdo.png",
	placeholder).
		SetExpiresAt(time.Now().UTC().Add(10 * time.Minute)).
		Save(ctx)
placeholder

	_, err = service.NewAuthPendingIdentityService(client).UpsertAdoptionDecision(ctx, service.PendingIdentityAdoptionDecisionInput{
		PendingAuthSessionID: session.ID,
		AdoptAvatar:          true,
placeholder)
placeholder

	body := bytes.NewBufferString(`{"invitation_code":"invite-1","adopt_display_name":trueplaceholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/oauth/linuxdo/complete-registration", body)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: oauthPendingSessionCookieName, Value: encodeCookieValue(session.SessionToken)placeholder)
	req.AddCookie(&http.Cookie{Name: oauthPendingBrowserCookieName, Value: encodeCookieValue("linuxdo-browser")placeholder)
	c.Request = req

	handler.CompleteLinuxDoOAuthRegistration(c)

	require.Equal(t, http.StatusOK, recorder.Code)
	responseData := decodeJSONBody(t, recorder)
	require.NotEmpty(t, responseData["access_token"])

	userEntity, err := client.User.Query().
		Where(dbuser.EmailEQ(session.ResolvedEmail)).
		Only(ctx)
placeholder
	require.Equal(t, "LinuxDo Display", userEntity.Username)

	identity, err := client.AuthIdentity.Query().
		Where(
			authidentity.ProviderTypeEQ("linuxdo"),
			authidentity.ProviderKeyEQ("linuxdo"),
			authidentity.ProviderSubjectEQ("linuxdo-subject-1"),
		).
		Only(ctx)
placeholder
	require.Equal(t, userEntity.ID, identity.UserID)
	require.Equal(t, "LinuxDo Display", identity.Metadata["display_name"])
	require.Equal(t, "https://cdn.example/linuxdo.png", identity.Metadata["avatar_url"])

	decision, err := client.IdentityAdoptionDecision.Query().
		Where(identityadoptiondecision.PendingAuthSessionIDEQ(session.ID)).
		Only(ctx)
placeholder
	require.NotNil(t, decision.IdentityID)
	require.Equal(t, identity.ID, *decision.IdentityID)
	require.True(t, decision.AdoptDisplayName)
	require.True(t, decision.AdoptAvatar)

	consumed, err := client.PendingAuthSession.Query().
		Where(pendingauthsession.IDEQ(session.ID)).
		Only(ctx)
placeholder
	require.NotNil(t, consumed.ConsumedAt)
placeholder
