package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/authidentity"
	"github.com/Wei-Shaw/sub2api/ent/enttest"
	"github.com/Wei-Shaw/sub2api/ent/identityadoptiondecision"
	"github.com/Wei-Shaw/sub2api/ent/pendingauthsession"
	dbuser "github.com/Wei-Shaw/sub2api/ent/user"
	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "modernc.org/sqlite"
)

func TestApplySuggestedProfileToCompletionResponse(t *testing.T) {
	payload := map[string]any{
		"access_token": "token",
placeholder
	upstream := map[string]any{
		"suggested_display_name": "Alice",
		"suggested_avatar_url":   "https://cdn.example/avatar.png",
placeholder

	applySuggestedProfileToCompletionResponse(payload, upstream)

	require.Equal(t, "Alice", payload["suggested_display_name"])
	require.Equal(t, "https://cdn.example/avatar.png", payload["suggested_avatar_url"])
	require.Equal(t, true, payload["adoption_required"])
placeholder

func TestApplySuggestedProfileToCompletionResponseKeepsExistingPayloadValues(t *testing.T) {
	payload := map[string]any{
		"suggested_display_name": "Existing",
		"adoption_required":      false,
placeholder
	upstream := map[string]any{
		"suggested_display_name": "Alice",
		"suggested_avatar_url":   "https://cdn.example/avatar.png",
placeholder

	applySuggestedProfileToCompletionResponse(payload, upstream)

	require.Equal(t, "Existing", payload["suggested_display_name"])
	require.Equal(t, "https://cdn.example/avatar.png", payload["suggested_avatar_url"])
	require.Equal(t, true, payload["adoption_required"])
placeholder

func TestExchangePendingOAuthCompletionPreviewThenFinalizeAppliesAdoptionDecision(t *testing.T) {
	handler, client := newOAuthPendingFlowTestHandler(t, false)
	ctx := context.Background()

	userEntity, err := client.User.Create().
		SetEmail("linuxdo-123@linuxdo-connect.invalid").
		SetUsername("legacy-name").
		SetPasswordHash("hash").
		SetRole(service.RoleUser).
		SetStatus(service.StatusActive).
		Save(ctx)
placeholder

	session, err := client.PendingAuthSession.Create().
		SetSessionToken("pending-session-token").
		SetIntent("login").
		SetProviderType("linuxdo").
		SetProviderKey("linuxdo").
		SetProviderSubject("123").
		SetTargetUserID(userEntity.ID).
		SetResolvedEmail(userEntity.Email).
		SetBrowserSessionKey("browser-session-key").
		SetUpstreamIdentityClaims(map[string]any{
			"username":               "linuxdo_user",
			"suggested_display_name": "Alice Example",
			"suggested_avatar_url":   "https://cdn.example/alice.png",
	placeholder).
		SetLocalFlowState(map[string]any{
			oauthCompletionResponseKey: map[string]any{
				"access_token": "access-token",
				"redirect":     "/dashboard",
		placeholder,
	placeholder).
		SetExpiresAt(time.Now().UTC().Add(10 * time.Minute)).
		Save(ctx)
placeholder

	previewRecorder := httptest.NewRecorder()
	previewCtx, _ := gin.CreateTestContext(previewRecorder)
	previewReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/oauth/pending/exchange", nil)
	previewReq.AddCookie(&http.Cookie{Name: oauthPendingSessionCookieName, Value: encodeCookieValue(session.SessionToken)placeholder)
	previewReq.AddCookie(&http.Cookie{Name: oauthPendingBrowserCookieName, Value: encodeCookieValue("browser-session-key")placeholder)
	previewCtx.Request = previewReq

	handler.ExchangePendingOAuthCompletion(previewCtx)

	require.Equal(t, http.StatusOK, previewRecorder.Code)
	previewData := decodeJSONResponseData(t, previewRecorder)
	require.Equal(t, "Alice Example", previewData["suggested_display_name"])
	require.Equal(t, "https://cdn.example/alice.png", previewData["suggested_avatar_url"])
	require.Equal(t, true, previewData["adoption_required"])

	storedUser, err := client.User.Get(ctx, userEntity.ID)
placeholder
	require.Equal(t, "legacy-name", storedUser.Username)

	previewSession, err := client.PendingAuthSession.Query().
		Where(pendingauthsession.IDEQ(session.ID)).
		Only(ctx)
placeholder
	require.Nil(t, previewSession.ConsumedAt)

	body := bytes.NewBufferString(`{"adopt_display_name":true,"adopt_avatar":trueplaceholder`)
	finalizeRecorder := httptest.NewRecorder()
	finalizeCtx, _ := gin.CreateTestContext(finalizeRecorder)
	finalizeReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/oauth/pending/exchange", body)
	finalizeReq.Header.Set("Content-Type", "application/json")
	finalizeReq.AddCookie(&http.Cookie{Name: oauthPendingSessionCookieName, Value: encodeCookieValue(session.SessionToken)placeholder)
	finalizeReq.AddCookie(&http.Cookie{Name: oauthPendingBrowserCookieName, Value: encodeCookieValue("browser-session-key")placeholder)
	finalizeCtx.Request = finalizeReq

	handler.ExchangePendingOAuthCompletion(finalizeCtx)

	require.Equal(t, http.StatusOK, finalizeRecorder.Code)

	storedUser, err = client.User.Get(ctx, userEntity.ID)
placeholder
	require.Equal(t, "Alice Example", storedUser.Username)

	identity, err := client.AuthIdentity.Query().
		Where(
			authidentity.ProviderTypeEQ("linuxdo"),
			authidentity.ProviderKeyEQ("linuxdo"),
			authidentity.ProviderSubjectEQ("123"),
		).
		Only(ctx)
placeholder
	require.Equal(t, userEntity.ID, identity.UserID)
	require.Equal(t, "Alice Example", identity.Metadata["display_name"])
	require.Equal(t, "https://cdn.example/alice.png", identity.Metadata["avatar_url"])

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

func newOAuthPendingFlowTestHandler(t *testing.T, invitationEnabled bool) (*AuthHandler, *dbent.Client) {
placeholder

	db, err := sql.Open("sqlite", "file:auth_oauth_pending_flow_handler?mode=memory&cache=shared")
placeholder
	t.Cleanup(func() { _ = db.Close() placeholder)

	_, err = db.Exec("PRAGMA foreign_keys = ON")
placeholder

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := enttest.NewClient(t, enttest.WithOptions(dbent.Driver(drv)))

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:                   "test-secret",
			ExpireHour:               1,
			AccessTokenExpireMinutes: 60,
			RefreshTokenExpireDays:   7,
	placeholder,
		Default: config.DefaultConfig{
			UserBalance:     0,
			UserConcurrency: 1,
	placeholder,
placeholder
	settingSvc := service.NewSettingService(&oauthPendingFlowSettingRepoStub{
		values: map[string]string{
			service.SettingKeyRegistrationEnabled:   "true",
			service.SettingKeyInvitationCodeEnabled: boolSettingValue(invitationEnabled),
	placeholder,
placeholder, cfg)
	authSvc := service.NewAuthService(
		client,
		&oauthPendingFlowUserRepo{client: clientplaceholder,
		nil,
		&oauthPendingFlowRefreshTokenCacheStub{placeholder,
		cfg,
		settingSvc,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	return &AuthHandler{
		authService: authSvc,
		settingSvc:  settingSvc,
placeholder, client
placeholder

func boolSettingValue(v bool) string {
	if v {
		return "true"
placeholder
	return "false"
placeholder

func boolPtr(v bool) *bool {
	return &v
placeholder

type oauthPendingFlowSettingRepoStub struct {
	values map[string]string
placeholder

func (s *oauthPendingFlowSettingRepoStub) Get(context.Context, string) (*service.Setting, error) {
	return nil, service.ErrSettingNotFound
placeholder

func (s *oauthPendingFlowSettingRepoStub) GetValue(_ context.Context, key string) (string, error) {
	value, ok := s.values[key]
	if !ok {
		return "", service.ErrSettingNotFound
placeholder
	return value, nil
placeholder

func (s *oauthPendingFlowSettingRepoStub) Set(context.Context, string, string) error {
	return nil
placeholder

func (s *oauthPendingFlowSettingRepoStub) GetMultiple(_ context.Context, keys []string) (map[string]string, error) {
	result := make(map[string]string, len(keys))
	for _, key := range keys {
		if value, ok := s.values[key]; ok {
			result[key] = value
	placeholder
placeholder
	return result, nil
placeholder

func (s *oauthPendingFlowSettingRepoStub) SetMultiple(context.Context, map[string]string) error {
	return nil
placeholder

func (s *oauthPendingFlowSettingRepoStub) GetAll(context.Context) (map[string]string, error) {
	result := make(map[string]string, len(s.values))
	for key, value := range s.values {
		result[key] = value
placeholder
	return result, nil
placeholder

func (s *oauthPendingFlowSettingRepoStub) Delete(context.Context, string) error {
	return nil
placeholder

type oauthPendingFlowRefreshTokenCacheStub struct{placeholder

func (s *oauthPendingFlowRefreshTokenCacheStub) StoreRefreshToken(context.Context, string, *service.RefreshTokenData, time.Duration) error {
	return nil
placeholder

func (s *oauthPendingFlowRefreshTokenCacheStub) GetRefreshToken(context.Context, string) (*service.RefreshTokenData, error) {
	return nil, service.ErrRefreshTokenNotFound
placeholder

func (s *oauthPendingFlowRefreshTokenCacheStub) DeleteRefreshToken(context.Context, string) error {
	return nil
placeholder

func (s *oauthPendingFlowRefreshTokenCacheStub) DeleteUserRefreshTokens(context.Context, int64) error {
	return nil
placeholder

func (s *oauthPendingFlowRefreshTokenCacheStub) DeleteTokenFamily(context.Context, string) error {
	return nil
placeholder

func (s *oauthPendingFlowRefreshTokenCacheStub) AddToUserTokenSet(context.Context, int64, string, time.Duration) error {
	return nil
placeholder

func (s *oauthPendingFlowRefreshTokenCacheStub) AddToFamilyTokenSet(context.Context, string, string, time.Duration) error {
	return nil
placeholder

func (s *oauthPendingFlowRefreshTokenCacheStub) GetUserTokenHashes(context.Context, int64) ([]string, error) {
	return nil, nil
placeholder

func (s *oauthPendingFlowRefreshTokenCacheStub) GetFamilyTokenHashes(context.Context, string) ([]string, error) {
	return nil, nil
placeholder

func (s *oauthPendingFlowRefreshTokenCacheStub) IsTokenInFamily(context.Context, string, string) (bool, error) {
	return false, nil
placeholder

func decodeJSONResponseData(t *testing.T, recorder *httptest.ResponseRecorder) map[string]any {
placeholder

	var envelope struct {
		Data map[string]any `json:"data"`
placeholder
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &envelope))
	return envelope.Data
placeholder

func decodeJSONBody(t *testing.T, recorder *httptest.ResponseRecorder) map[string]any {
placeholder

	var payload map[string]any
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &payload))
	return payload
placeholder

type oauthPendingFlowUserRepo struct {
	client *dbent.Client
placeholder

func (r *oauthPendingFlowUserRepo) Create(ctx context.Context, user *service.User) error {
	entity, err := r.client.User.Create().
		SetEmail(user.Email).
		SetUsername(user.Username).
		SetNotes(user.Notes).
		SetPasswordHash(user.PasswordHash).
		SetRole(user.Role).
		SetBalance(user.Balance).
		SetConcurrency(user.Concurrency).
		SetStatus(user.Status).
		SetSignupSource(user.SignupSource).
		SetNillableLastLoginAt(user.LastLoginAt).
		SetNillableLastActiveAt(user.LastActiveAt).
		Save(ctx)
	if err != nil {
		return err
placeholder
	user.ID = entity.ID
	user.CreatedAt = entity.CreatedAt
	user.UpdatedAt = entity.UpdatedAt
	return nil
placeholder

func (r *oauthPendingFlowUserRepo) GetByID(ctx context.Context, id int64) (*service.User, error) {
	entity, err := r.client.User.Get(ctx, id)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrUserNotFound
	placeholder
		return nil, err
placeholder
	return oauthPendingFlowServiceUser(entity), nil
placeholder

func (r *oauthPendingFlowUserRepo) GetByEmail(ctx context.Context, email string) (*service.User, error) {
	entity, err := r.client.User.Query().Where(dbuser.EmailEQ(email)).Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrUserNotFound
	placeholder
		return nil, err
placeholder
	return oauthPendingFlowServiceUser(entity), nil
placeholder

func (r *oauthPendingFlowUserRepo) GetFirstAdmin(context.Context) (*service.User, error) {
	panic("unexpected GetFirstAdmin call")
placeholder

func (r *oauthPendingFlowUserRepo) Update(ctx context.Context, user *service.User) error {
	entity, err := r.client.User.UpdateOneID(user.ID).
		SetEmail(user.Email).
		SetUsername(user.Username).
		SetNotes(user.Notes).
		SetPasswordHash(user.PasswordHash).
		SetRole(user.Role).
		SetBalance(user.Balance).
		SetConcurrency(user.Concurrency).
		SetStatus(user.Status).
		SetSignupSource(user.SignupSource).
		SetNillableLastLoginAt(user.LastLoginAt).
		SetNillableLastActiveAt(user.LastActiveAt).
		Save(ctx)
	if err != nil {
		return err
placeholder
	user.UpdatedAt = entity.UpdatedAt
	return nil
placeholder

func (r *oauthPendingFlowUserRepo) Delete(ctx context.Context, id int64) error {
	return r.client.User.DeleteOneID(id).Exec(ctx)
placeholder

func (r *oauthPendingFlowUserRepo) GetUserAvatar(context.Context, int64) (*service.UserAvatar, error) {
	return nil, service.ErrUserNotFound
placeholder

func (r *oauthPendingFlowUserRepo) UpsertUserAvatar(context.Context, int64, service.UpsertUserAvatarInput) (*service.UserAvatar, error) {
	panic("unexpected UpsertUserAvatar call")
placeholder

func (r *oauthPendingFlowUserRepo) DeleteUserAvatar(context.Context, int64) error {
	return nil
placeholder

func (r *oauthPendingFlowUserRepo) List(context.Context, pagination.PaginationParams) ([]service.User, *pagination.PaginationResult, error) {
	panic("unexpected List call")
placeholder

func (r *oauthPendingFlowUserRepo) ListWithFilters(context.Context, pagination.PaginationParams, service.UserListFilters) ([]service.User, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
placeholder

func (r *oauthPendingFlowUserRepo) UpdateBalance(context.Context, int64, float64) error {
	panic("unexpected UpdateBalance call")
placeholder

func (r *oauthPendingFlowUserRepo) DeductBalance(context.Context, int64, float64) error {
	panic("unexpected DeductBalance call")
placeholder

func (r *oauthPendingFlowUserRepo) UpdateConcurrency(context.Context, int64, int) error {
	panic("unexpected UpdateConcurrency call")
placeholder

func (r *oauthPendingFlowUserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	count, err := r.client.User.Query().Where(dbuser.EmailEQ(email)).Count(ctx)
	return count > 0, err
placeholder

func (r *oauthPendingFlowUserRepo) RemoveGroupFromAllowedGroups(context.Context, int64) (int64, error) {
	panic("unexpected RemoveGroupFromAllowedGroups call")
placeholder

func (r *oauthPendingFlowUserRepo) AddGroupToAllowedGroups(context.Context, int64, int64) error {
	panic("unexpected AddGroupToAllowedGroups call")
placeholder

func (r *oauthPendingFlowUserRepo) RemoveGroupFromUserAllowedGroups(context.Context, int64, int64) error {
	panic("unexpected RemoveGroupFromUserAllowedGroups call")
placeholder

func (r *oauthPendingFlowUserRepo) UpdateTotpSecret(context.Context, int64, *string) error {
	panic("unexpected UpdateTotpSecret call")
placeholder

func (r *oauthPendingFlowUserRepo) EnableTotp(context.Context, int64) error {
	panic("unexpected EnableTotp call")
placeholder

func (r *oauthPendingFlowUserRepo) DisableTotp(context.Context, int64) error {
	panic("unexpected DisableTotp call")
placeholder

func oauthPendingFlowServiceUser(entity *dbent.User) *service.User {
	if entity == nil {
		return nil
placeholder
	return &service.User{
		ID:           entity.ID,
		Email:        entity.Email,
		Username:     entity.Username,
		Notes:        entity.Notes,
		PasswordHash: entity.PasswordHash,
		Role:         entity.Role,
		Balance:      entity.Balance,
		Concurrency:  entity.Concurrency,
		Status:       entity.Status,
		SignupSource: entity.SignupSource,
		LastLoginAt:  entity.LastLoginAt,
		LastActiveAt: entity.LastActiveAt,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
placeholder
placeholder
