//go:build unit

package service_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/authidentity"
	"github.com/Wei-Shaw/sub2api/ent/enttest"
	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/repository"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "modernc.org/sqlite"
)

type authIdentityDefaultSubAssignerStub struct {
	calls []*service.AssignSubscriptionInput
placeholder

func (s *authIdentityDefaultSubAssignerStub) AssignOrExtendSubscription(
	_ context.Context,
	input *service.AssignSubscriptionInput,
) (*service.UserSubscription, bool, error) {
	cloned := *input
	s.calls = append(s.calls, &cloned)
	return &service.UserSubscription{UserID: input.UserID, GroupID: input.GroupIDplaceholder, true, nil
placeholder

type flakyAuthIdentityDefaultSubAssignerStub struct {
	failuresRemaining int
	calls             []*service.AssignSubscriptionInput
placeholder

func (s *flakyAuthIdentityDefaultSubAssignerStub) AssignOrExtendSubscription(
	_ context.Context,
	input *service.AssignSubscriptionInput,
) (*service.UserSubscription, bool, error) {
	cloned := *input
	s.calls = append(s.calls, &cloned)
	if s.failuresRemaining > 0 {
		s.failuresRemaining--
		return nil, false, errors.New("temporary assign failure")
placeholder
	return &service.UserSubscription{UserID: input.UserID, GroupID: input.GroupIDplaceholder, true, nil
placeholder

type authIdentitySettingRepoStub struct {
	values map[string]string
placeholder

func (s *authIdentitySettingRepoStub) Get(context.Context, string) (*service.Setting, error) {
	panic("unexpected Get call")
placeholder

func (s *authIdentitySettingRepoStub) GetValue(_ context.Context, key string) (string, error) {
	if v, ok := s.values[key]; ok {
		return v, nil
placeholder
	return "", service.ErrSettingNotFound
placeholder

func (s *authIdentitySettingRepoStub) Set(context.Context, string, string) error {
	panic("unexpected Set call")
placeholder

func (s *authIdentitySettingRepoStub) GetMultiple(_ context.Context, keys []string) (map[string]string, error) {
	out := make(map[string]string, len(keys))
	for _, key := range keys {
		if v, ok := s.values[key]; ok {
			out[key] = v
	placeholder
placeholder
	return out, nil
placeholder

func (s *authIdentitySettingRepoStub) SetMultiple(context.Context, map[string]string) error {
	panic("unexpected SetMultiple call")
placeholder

func (s *authIdentitySettingRepoStub) GetAll(context.Context) (map[string]string, error) {
	panic("unexpected GetAll call")
placeholder

func (s *authIdentitySettingRepoStub) Delete(context.Context, string) error {
	panic("unexpected Delete call")
placeholder

func newAuthServiceWithEnt(
	t *testing.T,
	settings map[string]string,
	defaultSubAssigner service.DefaultSubscriptionAssigner,
) (*service.AuthService, service.UserRepository, *dbent.Client) {
placeholder

	db, err := sql.Open("sqlite", "file:auth_service_identity_sync?mode=memory&cache=shared")
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
			Secret:     "test-auth-identity-secret",
			ExpireHour: 1,
	placeholder,
		Default: config.DefaultConfig{
			UserBalance:     3.5,
			UserConcurrency: 2,
	placeholder,
placeholder
	settingSvc := service.NewSettingService(&authIdentitySettingRepoStub{
		values: settings,
placeholder, cfg)

	svc := service.NewAuthService(client, repo, nil, nil, cfg, settingSvc, nil, nil, nil, nil, defaultSubAssigner)
	return svc, repo, client
placeholder

func TestAuthServiceRegisterDualWritesEmailIdentity(t *testing.T) {
	svc, _, client := newAuthServiceWithEnt(t, map[string]string{
		service.SettingKeyRegistrationEnabled: "true",
placeholder, nil)
	ctx := context.Background()

	token, user, err := svc.Register(ctx, "user@example.com", "password")
placeholder
	require.NotEmpty(t, token)
	require.NotNil(t, user)

	storedUser, err := client.User.Get(ctx, user.ID)
placeholder
	require.Equal(t, "email", storedUser.SignupSource)
	require.NotNil(t, storedUser.LastLoginAt)
	require.NotNil(t, storedUser.LastActiveAt)

	identity, err := client.AuthIdentity.Query().
		Where(
			authidentity.ProviderTypeEQ("email"),
			authidentity.ProviderKeyEQ("email"),
			authidentity.ProviderSubjectEQ("user@example.com"),
		).
		Only(ctx)
placeholder
	require.Equal(t, user.ID, identity.UserID)
	require.NotNil(t, identity.VerifiedAt)
placeholder

func TestAuthServiceLoginDefersLastLoginTouchUntilRecordSuccessfulLogin(t *testing.T) {
	svc, _, client := newAuthServiceWithEnt(t, map[string]string{
		service.SettingKeyRegistrationEnabled: "true",
placeholder, nil)
	ctx := context.Background()

	passwordHash, err := svc.HashPassword("password")
placeholder
	user, err := client.User.Create().
		SetEmail("login@example.com").
		SetPasswordHash(passwordHash).
		SetRole(service.RoleUser).
		SetStatus(service.StatusActive).
		SetBalance(1).
		SetConcurrency(1).
		Save(ctx)
placeholder

	old := time.Now().Add(-2 * time.Hour).UTC().Round(time.Second)
	_, err = client.User.UpdateOneID(user.ID).
		SetLastLoginAt(old).
		SetLastActiveAt(old).
		Save(ctx)
placeholder

	token, gotUser, err := svc.Login(ctx, user.Email, "password")
placeholder
	require.NotEmpty(t, token)
	require.NotNil(t, gotUser)

	storedUser, err := client.User.Get(ctx, user.ID)
placeholder
	require.NotNil(t, storedUser.LastLoginAt)
	require.NotNil(t, storedUser.LastActiveAt)
	require.True(t, storedUser.LastLoginAt.Equal(old))
	require.True(t, storedUser.LastActiveAt.Equal(old))

	identityCount, err := client.AuthIdentity.Query().
		Where(
			authidentity.ProviderTypeEQ("email"),
			authidentity.ProviderKeyEQ("email"),
			authidentity.ProviderSubjectEQ("login@example.com"),
		).
		Count(ctx)
placeholder
	require.Zero(t, identityCount)

	svc.RecordSuccessfulLogin(ctx, user.ID)

	identity, err := client.AuthIdentity.Query().
		Where(
			authidentity.ProviderTypeEQ("email"),
			authidentity.ProviderKeyEQ("email"),
			authidentity.ProviderSubjectEQ("login@example.com"),
		).
		Only(ctx)
placeholder
	require.Equal(t, user.ID, identity.UserID)
placeholder

func TestAuthServiceRecordSuccessfulLoginBackfillsEmailIdentity(t *testing.T) {
	svc, repo, client := newAuthServiceWithEnt(t, map[string]string{
		service.SettingKeyRegistrationEnabled: "true",
placeholder, nil)
	ctx := context.Background()

	user := &service.User{
		Email:       "record@example.com",
		Role:        service.RoleUser,
		Status:      service.StatusActive,
		Balance:     1,
		Concurrency: 1,
placeholder
	require.NoError(t, user.SetPassword("password"))
	require.NoError(t, repo.Create(ctx, user))

	svc.RecordSuccessfulLogin(ctx, user.ID)

	identity, err := client.AuthIdentity.Query().
		Where(
			authidentity.ProviderTypeEQ("email"),
			authidentity.ProviderKeyEQ("email"),
			authidentity.ProviderSubjectEQ("record@example.com"),
		).
		Only(ctx)
placeholder
	require.Equal(t, user.ID, identity.UserID)
placeholder

func TestAuthServiceLogin_DoesNotApplyEmailFirstBindDefaultsWhenBackfillingLegacyEmailIdentity(t *testing.T) {
	assigner := &authIdentityDefaultSubAssignerStub{placeholder
	svc, _, client := newAuthServiceWithEnt(t, map[string]string{
		service.SettingKeyRegistrationEnabled:                    "true",
		service.SettingKeyAuthSourceDefaultEmailBalance:          "8.5",
		service.SettingKeyAuthSourceDefaultEmailConcurrency:      "4",
		service.SettingKeyAuthSourceDefaultEmailSubscriptions:    `[{"group_id":11,"validity_days":30placeholder]`,
		service.SettingKeyAuthSourceDefaultEmailGrantOnFirstBind: "true",
placeholder, assigner)
	ctx := context.Background()

	passwordHash, err := svc.HashPassword("password")
placeholder
	user, err := client.User.Create().
		SetEmail("legacy@example.com").
		SetUsername("legacy-user").
		SetPasswordHash(passwordHash).
		SetBalance(1.5).
		SetConcurrency(2).
		SetRole(service.RoleUser).
		SetStatus(service.StatusActive).
		Save(ctx)
placeholder

	token, gotUser, err := svc.Login(ctx, user.Email, "password")
placeholder
	require.NotEmpty(t, token)
	require.NotNil(t, gotUser)
	svc.RecordSuccessfulLogin(ctx, user.ID)

	storedUser, err := client.User.Get(ctx, user.ID)
placeholder
	require.Equal(t, 1.5, storedUser.Balance)
	require.Equal(t, 2, storedUser.Concurrency)
	require.Empty(t, assigner.calls)

	identityCount, err := client.AuthIdentity.Query().
		Where(
			authidentity.ProviderTypeEQ("email"),
			authidentity.ProviderKeyEQ("email"),
			authidentity.ProviderSubjectEQ("legacy@example.com"),
		).
		Count(ctx)
placeholder
	require.Equal(t, 1, identityCount)
	require.Equal(t, 0, countProviderGrantRecords(t, client, user.ID, "email", "first_bind"))

	token, gotUser, err = svc.Login(ctx, user.Email, "password")
placeholder
	require.NotEmpty(t, token)
	require.NotNil(t, gotUser)

	storedUser, err = client.User.Get(ctx, user.ID)
placeholder
	require.Equal(t, 1.5, storedUser.Balance)
	require.Equal(t, 2, storedUser.Concurrency)
	require.Empty(t, assigner.calls)
	require.Equal(t, 0, countProviderGrantRecords(t, client, user.ID, "email", "first_bind"))
placeholder

func TestAuthServiceLogin_DoesNotApplyMergedEmailFirstBindDefaultsWhenBackfillingLegacyEmailIdentity(t *testing.T) {
	assigner := &authIdentityDefaultSubAssignerStub{placeholder
	svc, _, client := newAuthServiceWithEnt(t, map[string]string{
		service.SettingKeyRegistrationEnabled:                    "true",
		service.SettingKeyDefaultSubscriptions:                   `[{"group_id":21,"validity_days":14placeholder]`,
		service.SettingKeyAuthSourceDefaultEmailBalance:          "8.5",
		service.SettingKeyAuthSourceDefaultEmailConcurrency:      "5",
		service.SettingKeyAuthSourceDefaultEmailSubscriptions:    `[]`,
		service.SettingKeyAuthSourceDefaultEmailGrantOnFirstBind: "true",
placeholder, assigner)
	ctx := context.Background()

	passwordHash, err := svc.HashPassword("password")
placeholder
	user, err := client.User.Create().
		SetEmail("merged-first-bind@example.com").
		SetUsername("merged-user").
		SetPasswordHash(passwordHash).
		SetBalance(1.5).
		SetConcurrency(2).
		SetRole(service.RoleUser).
		SetStatus(service.StatusActive).
		Save(ctx)
placeholder

	token, gotUser, err := svc.Login(ctx, user.Email, "password")
placeholder
	require.NotEmpty(t, token)
	require.NotNil(t, gotUser)
	svc.RecordSuccessfulLogin(ctx, user.ID)

	storedUser, err := client.User.Get(ctx, user.ID)
placeholder
	require.Equal(t, 1.5, storedUser.Balance)
	require.Equal(t, 2, storedUser.Concurrency)
	require.Empty(t, assigner.calls)
	require.Equal(t, 0, countProviderGrantRecords(t, client, user.ID, "email", "first_bind"))
placeholder

func TestAuthServiceLogin_DoesNotApplyEmailFirstBindDefaultsWhenIdentityAlreadyExists(t *testing.T) {
	assigner := &authIdentityDefaultSubAssignerStub{placeholder
	svc, _, client := newAuthServiceWithEnt(t, map[string]string{
		service.SettingKeyRegistrationEnabled:                    "true",
		service.SettingKeyAuthSourceDefaultEmailBalance:          "8.5",
		service.SettingKeyAuthSourceDefaultEmailConcurrency:      "4",
		service.SettingKeyAuthSourceDefaultEmailSubscriptions:    `[{"group_id":11,"validity_days":30placeholder]`,
		service.SettingKeyAuthSourceDefaultEmailGrantOnFirstBind: "true",
placeholder, assigner)
	ctx := context.Background()

	passwordHash, err := svc.HashPassword("password")
placeholder
	user, err := client.User.Create().
		SetEmail("bound@example.com").
		SetUsername("bound-user").
		SetPasswordHash(passwordHash).
		SetBalance(2).
		SetConcurrency(3).
		SetRole(service.RoleUser).
		SetStatus(service.StatusActive).
		Save(ctx)
placeholder
	_, err = client.AuthIdentity.Create().
		SetUserID(user.ID).
		SetProviderType("email").
		SetProviderKey("email").
		SetProviderSubject("bound@example.com").
		SetVerifiedAt(time.Now().UTC()).
		SetMetadata(map[string]any{"source": "preexisting"placeholder).
		Save(ctx)
placeholder

	token, gotUser, err := svc.Login(ctx, user.Email, "password")
placeholder
	require.NotEmpty(t, token)
	require.NotNil(t, gotUser)
	svc.RecordSuccessfulLogin(ctx, user.ID)

	storedUser, err := client.User.Get(ctx, user.ID)
placeholder
	require.Equal(t, 2.0, storedUser.Balance)
	require.Equal(t, 3, storedUser.Concurrency)
	require.Empty(t, assigner.calls)
	require.Equal(t, 0, countProviderGrantRecords(t, client, user.ID, "email", "first_bind"))
placeholder

func TestAuthServiceLogin_DoesNotRetryEmailFirstBindDefaultsForBackfilledEmailIdentity(t *testing.T) {
	assigner := &flakyAuthIdentityDefaultSubAssignerStub{failuresRemaining: 1placeholder
	svc, _, client := newAuthServiceWithEnt(t, map[string]string{
		service.SettingKeyRegistrationEnabled:                    "true",
		service.SettingKeyAuthSourceDefaultEmailBalance:          "8.5",
		service.SettingKeyAuthSourceDefaultEmailConcurrency:      "4",
		service.SettingKeyAuthSourceDefaultEmailSubscriptions:    `[{"group_id":11,"validity_days":30placeholder]`,
		service.SettingKeyAuthSourceDefaultEmailGrantOnFirstBind: "true",
placeholder, assigner)
	ctx := context.Background()

	passwordHash, err := svc.HashPassword("password")
placeholder
	user, err := client.User.Create().
		SetEmail("retry-first-bind@example.com").
		SetUsername("retry-user").
		SetPasswordHash(passwordHash).
		SetBalance(1.5).
		SetConcurrency(2).
		SetRole(service.RoleUser).
		SetStatus(service.StatusActive).
		Save(ctx)
placeholder

	token, gotUser, err := svc.Login(ctx, user.Email, "password")
placeholder
	require.NotEmpty(t, token)
	require.NotNil(t, gotUser)
	svc.RecordSuccessfulLogin(ctx, user.ID)

	storedUser, err := client.User.Get(ctx, user.ID)
placeholder
	require.Equal(t, 1.5, storedUser.Balance)
	require.Equal(t, 2, storedUser.Concurrency)
	require.Empty(t, assigner.calls)
	require.Equal(t, 0, countProviderGrantRecords(t, client, user.ID, "email", "first_bind"))

	token, gotUser, err = svc.Login(ctx, user.Email, "password")
placeholder
	require.NotEmpty(t, token)
	require.NotNil(t, gotUser)
	svc.RecordSuccessfulLogin(ctx, user.ID)

	storedUser, err = client.User.Get(ctx, user.ID)
placeholder
	require.Equal(t, 1.5, storedUser.Balance)
	require.Equal(t, 2, storedUser.Concurrency)
	require.Empty(t, assigner.calls)
	require.Equal(t, 0, countProviderGrantRecords(t, client, user.ID, "email", "first_bind"))
placeholder

func countProviderGrantRecords(
	t *testing.T,
	client *dbent.Client,
	userID int64,
	providerType string,
	grantReason string,
) int {
placeholder

	var count int
	rows, err := client.QueryContext(
		context.Background(),
		`SELECT COUNT(*) FROM user_provider_default_grants WHERE user_id = ? AND provider_type = ? AND grant_reason = ?`,
		userID,
		providerType,
		grantReason,
	)
placeholder
	defer rows.Close()
	require.True(t, rows.Next())
	require.NoError(t, rows.Scan(&count))
	require.NoError(t, rows.Err())
	return count
placeholder
