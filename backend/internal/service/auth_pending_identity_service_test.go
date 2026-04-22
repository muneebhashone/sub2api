//go:build unit

package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/enttest"
	"github.com/Wei-Shaw/sub2api/ent/identityadoptiondecision"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "modernc.org/sqlite"
)

func newAuthPendingIdentityServiceTestClient(t *testing.T) (*AuthPendingIdentityService, *dbent.Client) {
placeholder

	db, err := sql.Open("sqlite", "file:auth_pending_identity_service?mode=memory&cache=shared")
placeholder
	t.Cleanup(func() { _ = db.Close() placeholder)

	_, err = db.Exec("PRAGMA foreign_keys = ON")
placeholder

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := enttest.NewClient(t, enttest.WithOptions(dbent.Driver(drv)))
	t.Cleanup(func() { _ = client.Close() placeholder)

	return NewAuthPendingIdentityService(client), client
placeholder

func TestAuthPendingIdentityService_CreatePendingSessionStoresSeparatedState(t *testing.T) {
	svc, client := newAuthPendingIdentityServiceTestClient(t)
	ctx := context.Background()

	targetUser, err := client.User.Create().
		SetEmail("pending-target@example.com").
		SetPasswordHash("hash").
		SetRole(RoleUser).
		SetStatus(StatusActive).
		Save(ctx)
placeholder

	session, err := svc.CreatePendingSession(ctx, CreatePendingAuthSessionInput{
		Intent: "bind_current_user",
		Identity: PendingAuthIdentityKey{
			ProviderType:    "wechat",
			ProviderKey:     "wechat-open",
			ProviderSubject: "union-123",
	placeholder,
		TargetUserID:           &targetUser.ID,
		RedirectTo:             "/profile",
		ResolvedEmail:          "user@example.com",
		BrowserSessionKey:      "browser-1",
		UpstreamIdentityClaims: map[string]any{"nickname": "wx-user", "avatar_url": "https://cdn.example/avatar.png"placeholder,
		LocalFlowState:         map[string]any{"step": "email_required"placeholder,
placeholder)
placeholder
	require.NotEmpty(t, session.SessionToken)
	require.Equal(t, "bind_current_user", session.Intent)
	require.Equal(t, "wechat", session.ProviderType)
	require.NotNil(t, session.TargetUserID)
	require.Equal(t, targetUser.ID, *session.TargetUserID)
	require.Equal(t, "wx-user", session.UpstreamIdentityClaims["nickname"])
	require.Equal(t, "email_required", session.LocalFlowState["step"])
placeholder

func TestAuthPendingIdentityService_CompletionCodeIsBrowserBoundAndOneTime(t *testing.T) {
	svc, _ := newAuthPendingIdentityServiceTestClient(t)
	ctx := context.Background()

	session, err := svc.CreatePendingSession(ctx, CreatePendingAuthSessionInput{
		Intent: "login",
		Identity: PendingAuthIdentityKey{
			ProviderType:    "linuxdo",
			ProviderKey:     "linuxdo-main",
			ProviderSubject: "subject-1",
	placeholder,
		BrowserSessionKey:      "browser-expected",
		UpstreamIdentityClaims: map[string]any{"nickname": "linux-user"placeholder,
		LocalFlowState:         map[string]any{"step": "pending"placeholder,
placeholder)
placeholder

	issued, err := svc.IssueCompletionCode(ctx, IssuePendingAuthCompletionCodeInput{
		PendingAuthSessionID: session.ID,
		BrowserSessionKey:    "browser-expected",
placeholder)
placeholder
	require.NotEmpty(t, issued.Code)

	_, err = svc.ConsumeCompletionCode(ctx, issued.Code, "browser-other")
	require.ErrorIs(t, err, ErrPendingAuthBrowserMismatch)

	consumed, err := svc.ConsumeCompletionCode(ctx, issued.Code, "browser-expected")
placeholder
	require.NotNil(t, consumed.ConsumedAt)
	require.Empty(t, consumed.CompletionCodeHash)
	require.Nil(t, consumed.CompletionCodeExpiresAt)

	_, err = svc.ConsumeCompletionCode(ctx, issued.Code, "browser-expected")
	require.ErrorIs(t, err, ErrPendingAuthCodeInvalid)
placeholder

func TestAuthPendingIdentityService_CompletionCodeExpires(t *testing.T) {
	svc, client := newAuthPendingIdentityServiceTestClient(t)
	ctx := context.Background()

	session, err := svc.CreatePendingSession(ctx, CreatePendingAuthSessionInput{
		Intent: "login",
		Identity: PendingAuthIdentityKey{
			ProviderType:    "oidc",
			ProviderKey:     "https://issuer.example",
			ProviderSubject: "subject-1",
	placeholder,
		BrowserSessionKey: "browser-expired",
placeholder)
placeholder

	issued, err := svc.IssueCompletionCode(ctx, IssuePendingAuthCompletionCodeInput{
		PendingAuthSessionID: session.ID,
		BrowserSessionKey:    "browser-expired",
		TTL:                  time.Second,
placeholder)
placeholder

	_, err = client.PendingAuthSession.UpdateOneID(session.ID).
		SetCompletionCodeExpiresAt(time.Now().UTC().Add(-time.Minute)).
		Save(ctx)
placeholder

	_, err = svc.ConsumeCompletionCode(ctx, issued.Code, "browser-expired")
	require.ErrorIs(t, err, ErrPendingAuthCodeExpired)
placeholder

func TestAuthPendingIdentityService_UpsertAdoptionDecision(t *testing.T) {
	svc, client := newAuthPendingIdentityServiceTestClient(t)
	ctx := context.Background()

	user, err := client.User.Create().
		SetEmail("adoption@example.com").
		SetPasswordHash("hash").
		SetRole(RoleUser).
		SetStatus(StatusActive).
		Save(ctx)
placeholder

	identity, err := client.AuthIdentity.Create().
		SetUserID(user.ID).
		SetProviderType("wechat").
		SetProviderKey("wechat-open").
		SetProviderSubject("union-adoption").
		SetMetadata(map[string]any{placeholder).
		Save(ctx)
placeholder

	session, err := svc.CreatePendingSession(ctx, CreatePendingAuthSessionInput{
		Intent: "bind_current_user",
		Identity: PendingAuthIdentityKey{
			ProviderType:    "wechat",
			ProviderKey:     "wechat-open",
			ProviderSubject: "union-adoption",
	placeholder,
placeholder)
placeholder

	first, err := svc.UpsertAdoptionDecision(ctx, PendingIdentityAdoptionDecisionInput{
		PendingAuthSessionID: session.ID,
		AdoptDisplayName:     true,
		AdoptAvatar:          false,
placeholder)
placeholder
	require.True(t, first.AdoptDisplayName)
	require.False(t, first.AdoptAvatar)
	require.Nil(t, first.IdentityID)

	second, err := svc.UpsertAdoptionDecision(ctx, PendingIdentityAdoptionDecisionInput{
		PendingAuthSessionID: session.ID,
		IdentityID:           &identity.ID,
		AdoptDisplayName:     true,
		AdoptAvatar:          true,
placeholder)
placeholder
	require.Equal(t, first.ID, second.ID)
	require.NotNil(t, second.IdentityID)
	require.Equal(t, identity.ID, *second.IdentityID)
	require.True(t, second.AdoptAvatar)
placeholder

func TestAuthPendingIdentityService_UpsertAdoptionDecision_ReassignsExistingIdentityReference(t *testing.T) {
	svc, client := newAuthPendingIdentityServiceTestClient(t)
	ctx := context.Background()

	user, err := client.User.Create().
		SetEmail("adoption-reassign@example.com").
		SetPasswordHash("hash").
		SetRole(RoleUser).
		SetStatus(StatusActive).
		Save(ctx)
placeholder

	identity, err := client.AuthIdentity.Create().
		SetUserID(user.ID).
		SetProviderType("wechat").
		SetProviderKey("wechat-open").
		SetProviderSubject("union-reassign").
		SetMetadata(map[string]any{placeholder).
		Save(ctx)
placeholder

	firstSession, err := svc.CreatePendingSession(ctx, CreatePendingAuthSessionInput{
		Intent: "bind_current_user",
		Identity: PendingAuthIdentityKey{
			ProviderType:    "wechat",
			ProviderKey:     "wechat-open",
			ProviderSubject: "union-reassign",
	placeholder,
placeholder)
placeholder

	firstDecision, err := svc.UpsertAdoptionDecision(ctx, PendingIdentityAdoptionDecisionInput{
		PendingAuthSessionID: firstSession.ID,
		IdentityID:           &identity.ID,
		AdoptDisplayName:     true,
		AdoptAvatar:          false,
placeholder)
placeholder
	require.NotNil(t, firstDecision.IdentityID)
	require.Equal(t, identity.ID, *firstDecision.IdentityID)

	secondSession, err := svc.CreatePendingSession(ctx, CreatePendingAuthSessionInput{
		Intent: "bind_current_user",
		Identity: PendingAuthIdentityKey{
			ProviderType:    "wechat",
			ProviderKey:     "wechat-open",
			ProviderSubject: "union-reassign",
	placeholder,
placeholder)
placeholder

	secondDecision, err := svc.UpsertAdoptionDecision(ctx, PendingIdentityAdoptionDecisionInput{
		PendingAuthSessionID: secondSession.ID,
		IdentityID:           &identity.ID,
		AdoptDisplayName:     false,
		AdoptAvatar:          true,
placeholder)
placeholder
	require.NotNil(t, secondDecision.IdentityID)
	require.Equal(t, identity.ID, *secondDecision.IdentityID)

	reloadedFirst, err := client.IdentityAdoptionDecision.Get(ctx, firstDecision.ID)
placeholder
	require.Nil(t, reloadedFirst.IdentityID)
placeholder

func TestAuthPendingIdentityService_UpsertAdoptionDecision_ClearsLegacyNullSessionReference(t *testing.T) {
	t.Skip("legacy NULL pending_auth_session_id rows only exist in production PostgreSQL history; sqlite unit schema rejects NULL")

	svc, client := newAuthPendingIdentityServiceTestClient(t)
	ctx := context.Background()

	user, err := client.User.Create().
		SetEmail("legacy-null-session@example.com").
		SetPasswordHash("hash").
		SetRole(RoleUser).
		SetStatus(StatusActive).
		Save(ctx)
placeholder

	identity, err := client.AuthIdentity.Create().
		SetUserID(user.ID).
		SetProviderType("wechat").
		SetProviderKey("wechat-main").
		SetProviderSubject("legacy-null-session").
		SetMetadata(map[string]any{placeholder).
		Save(ctx)
placeholder

	_, err = client.ExecContext(
		ctx,
		`INSERT INTO identity_adoption_decisions
			(identity_id, adopt_display_name, adopt_avatar, decided_at, created_at, updated_at, pending_auth_session_id)
		VALUES (?, ?, ?, ?, ?, ?, NULL)`,
		identity.ID,
		true,
		false,
		time.Now().UTC(),
		time.Now().UTC(),
		time.Now().UTC(),
	)
placeholder
	legacyDecision, err := client.IdentityAdoptionDecision.Query().
		Where(identityadoptiondecision.IdentityIDEQ(identity.ID)).
		Only(ctx)
placeholder
	require.NotNil(t, legacyDecision.IdentityID)

	session, err := svc.CreatePendingSession(ctx, CreatePendingAuthSessionInput{
		Intent: "bind_current_user",
		Identity: PendingAuthIdentityKey{
			ProviderType:    "wechat",
			ProviderKey:     "wechat-main",
			ProviderSubject: "legacy-null-session",
	placeholder,
placeholder)
placeholder

	decision, err := svc.UpsertAdoptionDecision(ctx, PendingIdentityAdoptionDecisionInput{
		PendingAuthSessionID: session.ID,
		IdentityID:           &identity.ID,
		AdoptDisplayName:     false,
		AdoptAvatar:          true,
placeholder)
placeholder
	require.NotNil(t, decision.IdentityID)
	require.Equal(t, identity.ID, *decision.IdentityID)

	reloadedLegacy, err := client.IdentityAdoptionDecision.Get(ctx, legacyDecision.ID)
placeholder
	require.Nil(t, reloadedLegacy.IdentityID)
placeholder

func TestAuthPendingIdentityService_ConsumeBrowserSession(t *testing.T) {
	svc, _ := newAuthPendingIdentityServiceTestClient(t)
	ctx := context.Background()

	session, err := svc.CreatePendingSession(ctx, CreatePendingAuthSessionInput{
		Intent: "login",
		Identity: PendingAuthIdentityKey{
			ProviderType:    "linuxdo",
			ProviderKey:     "linuxdo",
			ProviderSubject: "subject-session-token",
	placeholder,
		BrowserSessionKey: "browser-session",
		LocalFlowState: map[string]any{
			"completion_response": map[string]any{
				"access_token": "token",
		placeholder,
	placeholder,
placeholder)
placeholder

	_, err = svc.ConsumeBrowserSession(ctx, session.SessionToken, "browser-other")
	require.ErrorIs(t, err, ErrPendingAuthBrowserMismatch)

	consumed, err := svc.ConsumeBrowserSession(ctx, session.SessionToken, "browser-session")
placeholder
	require.NotNil(t, consumed.ConsumedAt)

	_, err = svc.ConsumeBrowserSession(ctx, session.SessionToken, "browser-session")
	require.ErrorIs(t, err, ErrPendingAuthSessionConsumed)
placeholder

func TestAuthPendingIdentityService_ConsumeBrowserSessionRejectsStaleLoadedSessionReplay(t *testing.T) {
	svc, _ := newAuthPendingIdentityServiceTestClient(t)
	ctx := context.Background()

	session, err := svc.CreatePendingSession(ctx, CreatePendingAuthSessionInput{
		Intent: "login",
		Identity: PendingAuthIdentityKey{
			ProviderType:    "linuxdo",
			ProviderKey:     "linuxdo",
			ProviderSubject: "stale-replay-subject",
	placeholder,
		BrowserSessionKey: "browser-session",
placeholder)
placeholder

	loaded, err := svc.getBrowserSession(ctx, session.SessionToken)
placeholder

	consumed, err := svc.consumeSession(ctx, loaded, "browser-session", ErrPendingAuthSessionExpired, ErrPendingAuthSessionConsumed)
placeholder
	require.NotNil(t, consumed.ConsumedAt)

	_, err = svc.consumeSession(ctx, loaded, "browser-session", ErrPendingAuthSessionExpired, ErrPendingAuthSessionConsumed)
	require.ErrorIs(t, err, ErrPendingAuthSessionConsumed)
placeholder

func TestAuthPendingIdentityService_ConsumeBrowserSessionScrubsLegacyCompletionTokens(t *testing.T) {
	svc, client := newAuthPendingIdentityServiceTestClient(t)
	ctx := context.Background()

	session, err := svc.CreatePendingSession(ctx, CreatePendingAuthSessionInput{
		Intent: "login",
		Identity: PendingAuthIdentityKey{
			ProviderType:    "linuxdo",
			ProviderKey:     "linuxdo",
			ProviderSubject: "legacy-token-subject",
	placeholder,
		BrowserSessionKey: "browser-session",
		LocalFlowState: map[string]any{
			"completion_response": map[string]any{
				"access_token":  "legacy-access-token",
				"refresh_token": "legacy-refresh-token",
				"expires_in":    float64(3600),
				"token_type":    "Bearer",
				"redirect":      "/dashboard",
		placeholder,
	placeholder,
placeholder)
placeholder

	consumed, err := svc.ConsumeBrowserSession(ctx, session.SessionToken, "browser-session")
placeholder
	require.NotNil(t, consumed.ConsumedAt)

	stored, err := client.PendingAuthSession.Get(ctx, session.ID)
placeholder

	completion, ok := stored.LocalFlowState["completion_response"].(map[string]any)
	require.True(t, ok)
	require.NotContains(t, completion, "access_token")
	require.NotContains(t, completion, "refresh_token")
	require.NotContains(t, completion, "expires_in")
	require.NotContains(t, completion, "token_type")
	require.Equal(t, "/dashboard", completion["redirect"])
placeholder
