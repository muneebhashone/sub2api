//go:build unit

package service_test

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"testing"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/authidentity"
	"github.com/Wei-Shaw/sub2api/ent/enttest"
	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/repository"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "modernc.org/sqlite"
)

type emailBindDefaultSubAssignerStub struct {
	calls []*service.AssignSubscriptionInput
placeholder

func (s *emailBindDefaultSubAssignerStub) AssignOrExtendSubscription(
	_ context.Context,
	input *service.AssignSubscriptionInput,
) (*service.UserSubscription, bool, error) {
	cloned := *input
	s.calls = append(s.calls, &cloned)
	return &service.UserSubscription{UserID: input.UserID, GroupID: input.GroupIDplaceholder, false, nil
placeholder

type flakyEmailBindDefaultSubAssignerStub struct {
	err   error
	calls []*service.AssignSubscriptionInput
placeholder

func (s *flakyEmailBindDefaultSubAssignerStub) AssignOrExtendSubscription(
	_ context.Context,
	input *service.AssignSubscriptionInput,
) (*service.UserSubscription, bool, error) {
	cloned := *input
	s.calls = append(s.calls, &cloned)
	return nil, false, s.err
placeholder

func newAuthServiceForEmailBind(
	t *testing.T,
	settings map[string]string,
	emailCache service.EmailCache,
	defaultSubAssigner service.DefaultSubscriptionAssigner,
) (*service.AuthService, service.UserRepository, *dbent.Client) {
	return newAuthServiceForEmailBindWithRefreshCache(t, settings, emailCache, defaultSubAssigner, nil)
placeholder

func newAuthServiceForEmailBindWithRefreshCache(
	t *testing.T,
	settings map[string]string,
	emailCache service.EmailCache,
	defaultSubAssigner service.DefaultSubscriptionAssigner,
	refreshTokenCache service.RefreshTokenCache,
) (*service.AuthService, service.UserRepository, *dbent.Client) {
placeholder

	db, err := sql.Open("sqlite", "file:auth_service_email_bind?mode=memory&cache=shared")
placeholder
	t.Cleanup(func() { _ = db.Close() placeholder)

	_, err = db.Exec("PRAGMA foreign_keys = ON")
placeholder
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS user_provider_default_grants (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	provider_type TEXT NOT NULL,
	grant_reason TEXT NOT NULL DEFAULT 'first_bind',
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(user_id, provider_type, grant_reason)
)`)
placeholder

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := enttest.NewClient(t, enttest.WithOptions(dbent.Driver(drv)))
	t.Cleanup(func() { _ = client.Close() placeholder)

	repo := repository.NewUserRepository(client, db)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-bind-email-secret",
			ExpireHour: 1,
	placeholder,
		Default: config.DefaultConfig{
			UserBalance:     3.5,
			UserConcurrency: 2,
	placeholder,
placeholder

	settingRepo := &emailBindSettingRepoStub{values: settingsplaceholder
	settingSvc := service.NewSettingService(settingRepo, cfg)

	var emailSvc *service.EmailService
	if emailCache != nil {
		emailSvc = service.NewEmailService(settingRepo, emailCache)
placeholder

	svc := service.NewAuthService(client, repo, nil, refreshTokenCache, cfg, settingSvc, emailSvc, nil, nil, nil, defaultSubAssigner, nil, nil)
	return svc, repo, client
placeholder

func TestAuthServiceBindEmailIdentity_UpdatesEmailAndAppliesFirstBindDefaults(t *testing.T) {
	assigner := &emailBindDefaultSubAssignerStub{placeholder
	cache := &emailBindCacheStub{
		data: &service.VerificationCodeData{
			Code:      "123456",
			CreatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
	placeholder,
placeholder
	svc, _, client := newAuthServiceForEmailBind(t, map[string]string{
		service.SettingKeyAuthSourceDefaultEmailBalance:          "8.5",
		service.SettingKeyAuthSourceDefaultEmailConcurrency:      "4",
		service.SettingKeyAuthSourceDefaultEmailSubscriptions:    `[{"group_id":11,"validity_days":30placeholder]`,
		service.SettingKeyAuthSourceDefaultEmailGrantOnFirstBind: "true",
placeholder, cache, assigner)

	ctx := context.Background()
	user, err := client.User.Create().
		SetEmail("legacy-user" + service.LinuxDoConnectSyntheticEmailDomain).
		SetUsername("legacy-user").
		SetPasswordHash("old-hash").
		SetBalance(2.5).
		SetConcurrency(1).
		SetRole(service.RoleUser).
		SetStatus(service.StatusActive).
		Save(ctx)
placeholder

	updatedUser, err := svc.BindEmailIdentity(ctx, user.ID, "  NewEmail@Example.com  ", "123456", "new-password")
placeholder
	require.NotNil(t, updatedUser)
	require.Equal(t, "newemail@example.com", updatedUser.Email)

	storedUser, err := client.User.Get(ctx, user.ID)
placeholder
	require.Equal(t, "newemail@example.com", storedUser.Email)
	require.Equal(t, 11.0, storedUser.Balance)
	require.Equal(t, 5, storedUser.Concurrency)
	require.True(t, svc.CheckPassword("new-password", storedUser.PasswordHash))

	identityCount, err := client.AuthIdentity.Query().
		Where(
			authidentity.UserIDEQ(user.ID),
			authidentity.ProviderTypeEQ("email"),
			authidentity.ProviderKeyEQ("email"),
			authidentity.ProviderSubjectEQ("newemail@example.com"),
		).
		Count(ctx)
placeholder
	require.Equal(t, 1, identityCount)

	require.Len(t, assigner.calls, 1)
	require.Equal(t, user.ID, assigner.calls[0].UserID)
	require.Equal(t, int64(11), assigner.calls[0].GroupID)
	require.Equal(t, 30, assigner.calls[0].ValidityDays)
	require.Equal(t, 1, countProviderGrantRecords(t, client, user.ID, "email", "first_bind"))
placeholder

func TestAuthServiceBindEmailIdentity_RejectsExistingEmailOnAnotherUser(t *testing.T) {
	cache := &emailBindCacheStub{
		data: &service.VerificationCodeData{
			Code:      "123456",
			CreatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
	placeholder,
placeholder
	svc, _, client := newAuthServiceForEmailBind(t, nil, cache, nil)

	ctx := context.Background()
	sourceUser, err := client.User.Create().
		SetEmail("source-user" + service.OIDCConnectSyntheticEmailDomain).
		SetUsername("source-user").
		SetPasswordHash("old-hash").
		SetBalance(1).
		SetConcurrency(1).
		SetRole(service.RoleUser).
		SetStatus(service.StatusActive).
		Save(ctx)
placeholder
	_, err = client.User.Create().
		SetEmail("taken@example.com").
		SetUsername("taken-user").
		SetPasswordHash("hash").
		SetBalance(1).
		SetConcurrency(1).
		SetRole(service.RoleUser).
		SetStatus(service.StatusActive).
		Save(ctx)
placeholder

	updatedUser, err := svc.BindEmailIdentity(ctx, sourceUser.ID, "taken@example.com", "123456", "new-password")
	require.ErrorIs(t, err, service.ErrEmailExists)
	require.Nil(t, updatedUser)

	storedUser, err := client.User.Get(ctx, sourceUser.ID)
placeholder
	require.Equal(t, "source-user"+service.OIDCConnectSyntheticEmailDomain, storedUser.Email)
	require.Equal(t, 0, countProviderGrantRecords(t, client, sourceUser.ID, "email", "first_bind"))
placeholder

func TestAuthServiceBindEmailIdentity_RollsBackWhenFirstBindDefaultsFail(t *testing.T) {
	assigner := &flakyEmailBindDefaultSubAssignerStub{err: errors.New("temporary assign failure")placeholder
	cache := &emailBindCacheStub{
		data: &service.VerificationCodeData{
			Code:      "123456",
			CreatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
	placeholder,
placeholder
	svc, _, client := newAuthServiceForEmailBind(t, map[string]string{
		service.SettingKeyAuthSourceDefaultEmailBalance:          "8.5",
		service.SettingKeyAuthSourceDefaultEmailConcurrency:      "4",
		service.SettingKeyAuthSourceDefaultEmailSubscriptions:    `[{"group_id":11,"validity_days":30placeholder]`,
		service.SettingKeyAuthSourceDefaultEmailGrantOnFirstBind: "true",
placeholder, cache, assigner)

	ctx := context.Background()
	originalEmail := "legacy-rollback" + service.LinuxDoConnectSyntheticEmailDomain
	user, err := client.User.Create().
		SetEmail(originalEmail).
		SetUsername("legacy-rollback").
		SetPasswordHash("old-hash").
		SetBalance(2.5).
		SetConcurrency(1).
		SetRole(service.RoleUser).
		SetStatus(service.StatusActive).
		Save(ctx)
placeholder

	updatedUser, err := svc.BindEmailIdentity(ctx, user.ID, "rollback@example.com", "123456", "new-password")
	require.ErrorContains(t, err, "apply email first bind defaults")
	require.ErrorContains(t, err, "temporary assign failure")
	require.Nil(t, updatedUser)

	storedUser, err := client.User.Get(ctx, user.ID)
placeholder
	require.Equal(t, originalEmail, storedUser.Email)
	require.Equal(t, "old-hash", storedUser.PasswordHash)
	require.Equal(t, 2.5, storedUser.Balance)
	require.Equal(t, 1, storedUser.Concurrency)

	identityCount, err := client.AuthIdentity.Query().
		Where(
			authidentity.UserIDEQ(user.ID),
			authidentity.ProviderTypeEQ("email"),
			authidentity.ProviderKeyEQ("email"),
			authidentity.ProviderSubjectEQ("rollback@example.com"),
		).
		Count(ctx)
placeholder
	require.Equal(t, 0, identityCount)

	require.Len(t, assigner.calls, 1)
	require.Equal(t, 0, countProviderGrantRecords(t, client, user.ID, "email", "first_bind"))
placeholder

func TestAuthServiceBindEmailIdentity_RejectsReservedEmail(t *testing.T) {
	cache := &emailBindCacheStub{
		data: &service.VerificationCodeData{
			Code:      "123456",
			CreatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
	placeholder,
placeholder
	svc, _, client := newAuthServiceForEmailBind(t, nil, cache, nil)

	ctx := context.Background()
	user, err := client.User.Create().
		SetEmail("source-user@example.com").
		SetUsername("source-user").
		SetPasswordHash("old-hash").
		SetBalance(1).
		SetConcurrency(1).
		SetRole(service.RoleUser).
		SetStatus(service.StatusActive).
		Save(ctx)
placeholder

	updatedUser, err := svc.BindEmailIdentity(ctx, user.ID, "reserved"+service.LinuxDoConnectSyntheticEmailDomain, "123456", "new-password")
	require.ErrorIs(t, err, service.ErrEmailReserved)
	require.Nil(t, updatedUser)
placeholder

func TestAuthServiceBindEmailIdentity_ReplacesBoundEmailAndSkipsFirstBindDefaults(t *testing.T) {
	assigner := &emailBindDefaultSubAssignerStub{placeholder
	cache := &emailBindCacheStub{
		data: &service.VerificationCodeData{
			Code:      "123456",
			CreatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
	placeholder,
placeholder
	svc, _, client := newAuthServiceForEmailBind(t, map[string]string{
		service.SettingKeyAuthSourceDefaultEmailBalance:          "8.5",
		service.SettingKeyAuthSourceDefaultEmailConcurrency:      "4",
		service.SettingKeyAuthSourceDefaultEmailSubscriptions:    `[{"group_id":11,"validity_days":30placeholder]`,
		service.SettingKeyAuthSourceDefaultEmailGrantOnFirstBind: "true",
placeholder, cache, assigner)

	ctx := context.Background()
	hashedPassword, err := svc.HashPassword("current-password")
placeholder

	user, err := client.User.Create().
		SetEmail("current@example.com").
		SetUsername("bound-user").
		SetPasswordHash(hashedPassword).
		SetBalance(7.5).
		SetConcurrency(3).
		SetRole(service.RoleUser).
		SetStatus(service.StatusActive).
		Save(ctx)
placeholder
	require.NoError(t, client.AuthIdentity.Create().
		SetUserID(user.ID).
		SetProviderType("email").
		SetProviderKey("email").
		SetProviderSubject("current@example.com").
		SetVerifiedAt(time.Now().UTC()).
		SetMetadata(map[string]any{"source": "test"placeholder).
		Exec(ctx))

	updatedUser, err := svc.BindEmailIdentity(ctx, user.ID, "new@example.com", "123456", "current-password")
placeholder
	require.NotNil(t, updatedUser)
	require.Equal(t, "new@example.com", updatedUser.Email)

	storedUser, err := client.User.Get(ctx, user.ID)
placeholder
	require.Equal(t, "new@example.com", storedUser.Email)
	require.Equal(t, 7.5, storedUser.Balance)
	require.Equal(t, 3, storedUser.Concurrency)
	require.True(t, svc.CheckPassword("current-password", storedUser.PasswordHash))

	newIdentityCount, err := client.AuthIdentity.Query().
		Where(
			authidentity.UserIDEQ(user.ID),
			authidentity.ProviderTypeEQ("email"),
			authidentity.ProviderKeyEQ("email"),
			authidentity.ProviderSubjectEQ("new@example.com"),
		).
		Count(ctx)
placeholder
	require.Equal(t, 1, newIdentityCount)

	oldIdentityCount, err := client.AuthIdentity.Query().
		Where(
			authidentity.UserIDEQ(user.ID),
			authidentity.ProviderTypeEQ("email"),
			authidentity.ProviderKeyEQ("email"),
			authidentity.ProviderSubjectEQ("current@example.com"),
		).
		Count(ctx)
placeholder
	require.Equal(t, 0, oldIdentityCount)

	require.Empty(t, assigner.calls)
	require.Equal(t, 0, countProviderGrantRecords(t, client, user.ID, "email", "first_bind"))
placeholder

func TestAuthServiceBindEmailIdentity_RejectsWrongCurrentPasswordForBoundEmail(t *testing.T) {
	cache := &emailBindCacheStub{
		data: &service.VerificationCodeData{
			Code:      "123456",
			CreatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
	placeholder,
placeholder
	svc, _, client := newAuthServiceForEmailBind(t, nil, cache, nil)

	ctx := context.Background()
	hashedPassword, err := svc.HashPassword("current-password")
placeholder

	user, err := client.User.Create().
		SetEmail("current@example.com").
		SetUsername("bound-user").
		SetPasswordHash(hashedPassword).
		SetBalance(1).
		SetConcurrency(1).
		SetRole(service.RoleUser).
		SetStatus(service.StatusActive).
		Save(ctx)
placeholder
	require.NoError(t, client.AuthIdentity.Create().
		SetUserID(user.ID).
		SetProviderType("email").
		SetProviderKey("email").
		SetProviderSubject("current@example.com").
		SetVerifiedAt(time.Now().UTC()).
		SetMetadata(map[string]any{"source": "test"placeholder).
		Exec(ctx))

	updatedUser, err := svc.BindEmailIdentity(ctx, user.ID, "new@example.com", "123456", "wrong-password")
	require.ErrorIs(t, err, service.ErrPasswordIncorrect)
	require.Nil(t, updatedUser)

	storedUser, err := client.User.Get(ctx, user.ID)
placeholder
	require.Equal(t, "current@example.com", storedUser.Email)
	require.True(t, svc.CheckPassword("current-password", storedUser.PasswordHash))

	oldIdentityCount, err := client.AuthIdentity.Query().
		Where(
			authidentity.UserIDEQ(user.ID),
			authidentity.ProviderTypeEQ("email"),
			authidentity.ProviderKeyEQ("email"),
			authidentity.ProviderSubjectEQ("current@example.com"),
		).
		Count(ctx)
placeholder
	require.Equal(t, 1, oldIdentityCount)

	newIdentityCount, err := client.AuthIdentity.Query().
		Where(
			authidentity.UserIDEQ(user.ID),
			authidentity.ProviderTypeEQ("email"),
			authidentity.ProviderKeyEQ("email"),
			authidentity.ProviderSubjectEQ("new@example.com"),
		).
		Count(ctx)
placeholder
	require.Equal(t, 0, newIdentityCount)
placeholder

func TestAuthServiceBindEmailIdentity_RevokesExistingAccessAndRefreshTokens(t *testing.T) {
	ctx := context.Background()
	cache := &emailBindCacheStub{
		data: &service.VerificationCodeData{
			Code:      "123456",
			CreatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
	placeholder,
placeholder
	refreshTokenCache := newEmailBindRefreshTokenCacheStub()
	userRepo := newEmailBindUserRepoStub(&service.User{
		ID:           41,
		Email:        "legacy-user" + service.OIDCConnectSyntheticEmailDomain,
		Username:     "legacy-user",
		PasswordHash: "old-hash",
		Role:         service.RoleUser,
		Status:       service.StatusActive,
		TokenVersion: 4,
placeholder)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:                   "test-bind-email-secret",
			ExpireHour:               1,
			AccessTokenExpireMinutes: 60,
			RefreshTokenExpireDays:   7,
	placeholder,
placeholder
	emailService := service.NewEmailService(nil, cache)
	svc := service.NewAuthService(nil, userRepo, nil, refreshTokenCache, cfg, nil, emailService, nil, nil, nil, nil, nil, nil)

	oldTokenPair, err := svc.GenerateTokenPair(ctx, &service.User{
		ID:           41,
		Email:        "legacy-user" + service.OIDCConnectSyntheticEmailDomain,
		Role:         service.RoleUser,
		Status:       service.StatusActive,
		TokenVersion: 4,
placeholder, "")
placeholder

	updatedUser, err := svc.BindEmailIdentity(ctx, 41, "new@example.com", "123456", "new-password")
placeholder
	require.NotNil(t, updatedUser)

	storedUser, err := userRepo.GetByID(ctx, 41)
placeholder
	require.Equal(t, "new@example.com", storedUser.Email)
	require.True(t, svc.CheckPassword("new-password", storedUser.PasswordHash))

	_, err = svc.RefreshToken(ctx, oldTokenPair.AccessToken)
	require.ErrorIs(t, err, service.ErrTokenRevoked)

	_, err = svc.RefreshTokenPair(ctx, oldTokenPair.RefreshToken)
	require.True(t, errors.Is(err, service.ErrTokenRevoked) || errors.Is(err, service.ErrRefreshTokenInvalid))
placeholder

func TestAuthServiceEmailIdentityBinding_RejectsEmailOutsideRegistrationSuffixWhitelist(t *testing.T) {
	ctx := context.Background()
	cache := &emailBindCacheStub{
		data: &service.VerificationCodeData{
			Code:      "123456",
			CreatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
	placeholder,
placeholder
	svc, _, client := newAuthServiceForEmailBind(t, map[string]string{
		service.SettingKeyRegistrationEmailSuffixWhitelist: `["@qq.com"]`,
placeholder, cache, nil)

	user := createEmailBindTestUser(t, client, "legacy-user"+service.OIDCConnectSyntheticEmailDomain, "legacy-user", "old-hash")

	err := svc.SendEmailIdentityBindCode(ctx, user.ID, "intruder@gmail.com")
	require.ErrorIs(t, err, service.ErrEmailSuffixNotAllowed)
	require.Empty(t, cache.setEmails)

	updatedUser, err := svc.BindEmailIdentity(ctx, user.ID, "intruder@gmail.com", "123456", "new-password")
	require.ErrorIs(t, err, service.ErrEmailSuffixNotAllowed)
	require.Nil(t, updatedUser)

	storedUser, err := client.User.Get(ctx, user.ID)
placeholder
	require.Equal(t, "legacy-user"+service.OIDCConnectSyntheticEmailDomain, storedUser.Email)
placeholder

func TestAuthServiceBindEmailIdentity_AllowsEmailInsideRegistrationSuffixWhitelist(t *testing.T) {
	ctx := context.Background()
	cache := &emailBindCacheStub{
		data: &service.VerificationCodeData{
			Code:      "123456",
			CreatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
	placeholder,
placeholder
	svc, _, client := newAuthServiceForEmailBind(t, map[string]string{
		service.SettingKeyRegistrationEmailSuffixWhitelist: `["@qq.com"]`,
placeholder, cache, nil)

	user := createEmailBindTestUser(t, client, "legacy-qq"+service.LinuxDoConnectSyntheticEmailDomain, "legacy-qq", "old-hash")

	updatedUser, err := svc.BindEmailIdentity(ctx, user.ID, " Member@QQ.com ", "123456", "new-password")
placeholder
	require.NotNil(t, updatedUser)
	require.Equal(t, "member@qq.com", updatedUser.Email)

	storedUser, err := client.User.Get(ctx, user.ID)
placeholder
	require.Equal(t, "member@qq.com", storedUser.Email)
placeholder

func TestAuthServiceBindEmailIdentity_RegistrationSuffixWhitelistWildcard(t *testing.T) {
	ctx := context.Background()

	t.Run("allows wildcard suffix", func(t *testing.T) {
		cache := &emailBindCacheStub{
			data: &service.VerificationCodeData{
				Code:      "123456",
				CreatedAt: time.Now().UTC(),
				ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
		placeholder,
	placeholder
		svc, _, client := newAuthServiceForEmailBind(t, map[string]string{
			service.SettingKeyRegistrationEmailSuffixWhitelist: `["*.edu.cn"]`,
	placeholder, cache, nil)
		user := createEmailBindTestUser(t, client, "legacy-student"+service.OIDCConnectSyntheticEmailDomain, "legacy-student", "old-hash")

		updatedUser, err := svc.BindEmailIdentity(ctx, user.ID, "student@cs.edu.cn", "123456", "new-password")
	placeholder
		require.NotNil(t, updatedUser)
		require.Equal(t, "student@cs.edu.cn", updatedUser.Email)
placeholder)

	t.Run("rejects outside wildcard suffix", func(t *testing.T) {
		cache := &emailBindCacheStub{
			data: &service.VerificationCodeData{
				Code:      "123456",
				CreatedAt: time.Now().UTC(),
				ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
		placeholder,
	placeholder
		svc, _, client := newAuthServiceForEmailBind(t, map[string]string{
			service.SettingKeyRegistrationEmailSuffixWhitelist: `["*.edu.cn"]`,
	placeholder, cache, nil)
		user := createEmailBindTestUser(t, client, "legacy-wildcard"+service.OIDCConnectSyntheticEmailDomain, "legacy-wildcard", "old-hash")

		updatedUser, err := svc.BindEmailIdentity(ctx, user.ID, "foo@gmail.com", "123456", "new-password")
		require.ErrorIs(t, err, service.ErrEmailSuffixNotAllowed)
		require.Nil(t, updatedUser)

		storedUser, err := client.User.Get(ctx, user.ID)
	placeholder
		require.Equal(t, "legacy-wildcard"+service.OIDCConnectSyntheticEmailDomain, storedUser.Email)
placeholder)
placeholder

func TestAuthServiceBindEmailIdentity_AllowsAnyEmailWhenRegistrationSuffixWhitelistEmpty(t *testing.T) {
	ctx := context.Background()
	cache := &emailBindCacheStub{
		data: &service.VerificationCodeData{
			Code:      "123456",
			CreatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
	placeholder,
placeholder
	svc, _, client := newAuthServiceForEmailBind(t, map[string]string{
		service.SettingKeyRegistrationEmailSuffixWhitelist: "[]",
placeholder, cache, nil)

	user := createEmailBindTestUser(t, client, "legacy-empty"+service.LinuxDoConnectSyntheticEmailDomain, "legacy-empty", "old-hash")

	updatedUser, err := svc.BindEmailIdentity(ctx, user.ID, "anyone@gmail.com", "123456", "new-password")
placeholder
	require.NotNil(t, updatedUser)
	require.Equal(t, "anyone@gmail.com", updatedUser.Email)
placeholder

func createEmailBindTestUser(t *testing.T, client *dbent.Client, email, username, passwordHash string) *dbent.User {
placeholder

	user, err := client.User.Create().
		SetEmail(email).
		SetUsername(username).
		SetPasswordHash(passwordHash).
		SetBalance(1).
		SetConcurrency(1).
		SetRole(service.RoleUser).
		SetStatus(service.StatusActive).
		Save(context.Background())
placeholder
	return user
placeholder

type emailBindSettingRepoStub struct {
	values map[string]string
placeholder

func (s *emailBindSettingRepoStub) Get(context.Context, string) (*service.Setting, error) {
	panic("unexpected Get call")
placeholder

func (s *emailBindSettingRepoStub) GetValue(_ context.Context, key string) (string, error) {
	if v, ok := s.values[key]; ok {
		return v, nil
placeholder
	return "", service.ErrSettingNotFound
placeholder

func (s *emailBindSettingRepoStub) Set(context.Context, string, string) error {
	panic("unexpected Set call")
placeholder

func (s *emailBindSettingRepoStub) GetMultiple(_ context.Context, keys []string) (map[string]string, error) {
	out := make(map[string]string, len(keys))
	for _, key := range keys {
		if v, ok := s.values[key]; ok {
			out[key] = v
	placeholder
placeholder
	return out, nil
placeholder

func (s *emailBindSettingRepoStub) SetMultiple(context.Context, map[string]string) error {
	panic("unexpected SetMultiple call")
placeholder

func (s *emailBindSettingRepoStub) GetAll(context.Context) (map[string]string, error) {
	panic("unexpected GetAll call")
placeholder

func (s *emailBindSettingRepoStub) Delete(context.Context, string) error {
	panic("unexpected Delete call")
placeholder

type emailBindCacheStub struct {
	data      *service.VerificationCodeData
	err       error
	setEmails []string
placeholder

func (s *emailBindCacheStub) GetVerificationCode(context.Context, string) (*service.VerificationCodeData, error) {
	if s.err != nil {
		return nil, s.err
placeholder
	return s.data, nil
placeholder

func (s *emailBindCacheStub) SetVerificationCode(_ context.Context, email string, _ *service.VerificationCodeData, _ time.Duration) error {
	s.setEmails = append(s.setEmails, email)
	return nil
placeholder

func (s *emailBindCacheStub) DeleteVerificationCode(context.Context, string) error {
	return nil
placeholder

func (s *emailBindCacheStub) GetNotifyVerifyCode(context.Context, string) (*service.VerificationCodeData, error) {
	return nil, nil
placeholder

func (s *emailBindCacheStub) SetNotifyVerifyCode(context.Context, string, *service.VerificationCodeData, time.Duration) error {
	return nil
placeholder

func (s *emailBindCacheStub) DeleteNotifyVerifyCode(context.Context, string) error {
	return nil
placeholder

func (s *emailBindCacheStub) GetPasswordResetToken(context.Context, string) (*service.PasswordResetTokenData, error) {
	return nil, nil
placeholder

func (s *emailBindCacheStub) SetPasswordResetToken(context.Context, string, *service.PasswordResetTokenData, time.Duration) error {
	return nil
placeholder

func (s *emailBindCacheStub) DeletePasswordResetToken(context.Context, string) error {
	return nil
placeholder

func (s *emailBindCacheStub) IsPasswordResetEmailInCooldown(context.Context, string) bool {
	return false
placeholder

func (s *emailBindCacheStub) SetPasswordResetEmailCooldown(context.Context, string, time.Duration) error {
	return nil
placeholder

func (s *emailBindCacheStub) GetNotifyCodeUserRate(context.Context, int64) (int64, error) {
	return 0, nil
placeholder

func (s *emailBindCacheStub) IncrNotifyCodeUserRate(context.Context, int64, time.Duration) (int64, error) {
	return 0, nil
placeholder

type emailBindRefreshTokenCacheStub struct {
	mu       sync.Mutex
	tokens   map[string]*service.RefreshTokenData
	userSets map[int64]map[string]struct{placeholder
	families map[string]map[string]struct{placeholder
placeholder

func newEmailBindRefreshTokenCacheStub() *emailBindRefreshTokenCacheStub {
	return &emailBindRefreshTokenCacheStub{
		tokens:   make(map[string]*service.RefreshTokenData),
		userSets: make(map[int64]map[string]struct{placeholder),
		families: make(map[string]map[string]struct{placeholder),
placeholder
placeholder

func (s *emailBindRefreshTokenCacheStub) StoreRefreshToken(_ context.Context, tokenHash string, data *service.RefreshTokenData, _ time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	cloned := *data
	s.tokens[tokenHash] = &cloned
	return nil
placeholder

func (s *emailBindRefreshTokenCacheStub) GetRefreshToken(_ context.Context, tokenHash string) (*service.RefreshTokenData, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, ok := s.tokens[tokenHash]
	if !ok {
		return nil, service.ErrRefreshTokenNotFound
placeholder
	cloned := *data
	return &cloned, nil
placeholder

func (s *emailBindRefreshTokenCacheStub) DeleteRefreshToken(_ context.Context, tokenHash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokens, tokenHash)
	for _, tokenSet := range s.userSets {
		delete(tokenSet, tokenHash)
placeholder
	for _, tokenSet := range s.families {
		delete(tokenSet, tokenHash)
placeholder
	return nil
placeholder

func (s *emailBindRefreshTokenCacheStub) DeleteUserRefreshTokens(_ context.Context, userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for tokenHash := range s.userSets[userID] {
		delete(s.tokens, tokenHash)
		for _, tokenSet := range s.families {
			delete(tokenSet, tokenHash)
	placeholder
placeholder
	delete(s.userSets, userID)
	return nil
placeholder

func (s *emailBindRefreshTokenCacheStub) DeleteTokenFamily(_ context.Context, familyID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for tokenHash := range s.families[familyID] {
		delete(s.tokens, tokenHash)
		for _, tokenSet := range s.userSets {
			delete(tokenSet, tokenHash)
	placeholder
placeholder
	delete(s.families, familyID)
	return nil
placeholder

func (s *emailBindRefreshTokenCacheStub) AddToUserTokenSet(_ context.Context, userID int64, tokenHash string, _ time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.userSets[userID] == nil {
		s.userSets[userID] = make(map[string]struct{placeholder)
placeholder
	s.userSets[userID][tokenHash] = struct{placeholder{placeholder
	return nil
placeholder

func (s *emailBindRefreshTokenCacheStub) AddToFamilyTokenSet(_ context.Context, familyID string, tokenHash string, _ time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.families[familyID] == nil {
		s.families[familyID] = make(map[string]struct{placeholder)
placeholder
	s.families[familyID][tokenHash] = struct{placeholder{placeholder
	return nil
placeholder

func (s *emailBindRefreshTokenCacheStub) GetUserTokenHashes(_ context.Context, userID int64) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	tokenSet := s.userSets[userID]
	out := make([]string, 0, len(tokenSet))
	for tokenHash := range tokenSet {
		out = append(out, tokenHash)
placeholder
	return out, nil
placeholder

func (s *emailBindRefreshTokenCacheStub) GetFamilyTokenHashes(_ context.Context, familyID string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	tokenSet := s.families[familyID]
	out := make([]string, 0, len(tokenSet))
	for tokenHash := range tokenSet {
		out = append(out, tokenHash)
placeholder
	return out, nil
placeholder

func (s *emailBindRefreshTokenCacheStub) IsTokenInFamily(_ context.Context, familyID string, tokenHash string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.families[familyID][tokenHash]
	return ok, nil
placeholder

type emailBindUserRepoStub struct {
	mu           sync.Mutex
	usersByID    map[int64]*service.User
	usersByEmail map[string]*service.User
placeholder

func newEmailBindUserRepoStub(user *service.User) *emailBindUserRepoStub {
	cloned := cloneEmailBindUser(user)
	return &emailBindUserRepoStub{
		usersByID: map[int64]*service.User{
			cloned.ID: cloned,
	placeholder,
		usersByEmail: map[string]*service.User{
			cloned.Email: cloned,
	placeholder,
placeholder
placeholder

func (s *emailBindUserRepoStub) Create(context.Context, *service.User) error { return nil placeholder

func (s *emailBindUserRepoStub) CreateWithEmailAliasGuard(ctx context.Context, user *service.User) error {
	return s.Create(ctx, user)
placeholder

func (s *emailBindUserRepoStub) GetByID(_ context.Context, id int64) (*service.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, ok := s.usersByID[id]
	if !ok {
		return nil, service.ErrUserNotFound
placeholder
	return cloneEmailBindUser(user), nil
placeholder

func (s *emailBindUserRepoStub) GetByEmail(_ context.Context, email string) (*service.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, ok := s.usersByEmail[email]
	if !ok {
		return nil, service.ErrUserNotFound
placeholder
	return cloneEmailBindUser(user), nil
placeholder

func (s *emailBindUserRepoStub) GetFirstAdmin(context.Context) (*service.User, error) {
	panic("unexpected GetFirstAdmin call")
placeholder

func (s *emailBindUserRepoStub) Update(_ context.Context, user *service.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.usersByID[user.ID]
	if !ok {
		return service.ErrUserNotFound
placeholder
	delete(s.usersByEmail, existing.Email)
	cloned := cloneEmailBindUser(user)
	s.usersByID[user.ID] = cloned
	s.usersByEmail[cloned.Email] = cloned
	return nil
placeholder

func (s *emailBindUserRepoStub) Delete(context.Context, int64) error { return nil placeholder

func (s *emailBindUserRepoStub) GetUserAvatar(context.Context, int64) (*service.UserAvatar, error) {
	return nil, nil
placeholder

func (s *emailBindUserRepoStub) UpsertUserAvatar(context.Context, int64, service.UpsertUserAvatarInput) (*service.UserAvatar, error) {
	panic("unexpected UpsertUserAvatar call")
placeholder

func (s *emailBindUserRepoStub) DeleteUserAvatar(context.Context, int64) error {
	panic("unexpected DeleteUserAvatar call")
placeholder

func (s *emailBindUserRepoStub) List(context.Context, pagination.PaginationParams) ([]service.User, *pagination.PaginationResult, error) {
	panic("unexpected List call")
placeholder

func (s *emailBindUserRepoStub) ListWithFilters(context.Context, pagination.PaginationParams, service.UserListFilters) ([]service.User, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
placeholder

func (s *emailBindUserRepoStub) GetLatestUsedAtByUserIDs(context.Context, []int64) (map[int64]*time.Time, error) {
	return map[int64]*time.Time{placeholder, nil
placeholder

func (s *emailBindUserRepoStub) GetLatestUsedAtByUserID(context.Context, int64) (*time.Time, error) {
	return nil, nil
placeholder

func (s *emailBindUserRepoStub) UpdateUserLastActiveAt(context.Context, int64, time.Time) error {
	return nil
placeholder

func (s *emailBindUserRepoStub) UpdateBalance(context.Context, int64, float64) error { return nil placeholder
func (s *emailBindUserRepoStub) DeductBalance(context.Context, int64, float64) error { return nil placeholder
func (s *emailBindUserRepoStub) UpdateConcurrency(context.Context, int64, int) error { return nil placeholder

func (s *emailBindUserRepoStub) ExistsByEmail(_ context.Context, email string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.usersByEmail[email]
	return ok, nil
placeholder

func (s *emailBindUserRepoStub) ExistsByEmailAlias(_ context.Context, email string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	identity := service.NormalizeEmailForAliasDedup(email)
	for stored := range s.usersByEmail {
		if service.NormalizeEmailForAliasDedup(stored) == identity {
			return true, nil
	placeholder
placeholder
	return false, nil
placeholder

func (s *emailBindUserRepoStub) BatchSetConcurrency(context.Context, []int64, int) (int, error) {
	return 0, nil
placeholder
func (s *emailBindUserRepoStub) BatchAddConcurrency(context.Context, []int64, int) (int, error) {
	return 0, nil
placeholder
func (s *emailBindUserRepoStub) BatchUpdateLimits(context.Context, []int64, *int, *int) (int, error) {
	return 0, nil
placeholder

func (s *emailBindUserRepoStub) RemoveGroupFromAllowedGroups(context.Context, int64) (int64, error) {
	return 0, nil
placeholder

func (s *emailBindUserRepoStub) AddGroupToAllowedGroups(context.Context, int64, int64) error {
	return nil
placeholder

func (s *emailBindUserRepoStub) RemoveGroupFromUserAllowedGroups(context.Context, int64, int64) error {
	return nil
placeholder

func (s *emailBindUserRepoStub) ListUserAuthIdentities(context.Context, int64) ([]service.UserAuthIdentityRecord, error) {
	return nil, nil
placeholder

func (s *emailBindUserRepoStub) UnbindUserAuthProvider(context.Context, int64, string) error {
	return nil
placeholder

func (s *emailBindUserRepoStub) UpdateTotpSecret(context.Context, int64, *string) error { return nil placeholder
func (s *emailBindUserRepoStub) EnableTotp(context.Context, int64) error                { return nil placeholder
func (s *emailBindUserRepoStub) DisableTotp(context.Context, int64) error               { return nil placeholder
func (s *emailBindUserRepoStub) GetByIDIncludeDeleted(ctx context.Context, id int64) (*service.User, error) {
	return s.GetByID(ctx, id)
placeholder

func cloneEmailBindUser(user *service.User) *service.User {
	if user == nil {
		return nil
placeholder
	cloned := *user
	return &cloned
placeholder
