//go:build unit

package service_test

import (
	"context"
	"database/sql"
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

func (s *authIdentitySettingRepoStub) GetMultiple(context.Context, []string) (map[string]string, error) {
	panic("unexpected GetMultiple call")
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

func newAuthServiceWithEnt(t *testing.T) (*service.AuthService, service.UserRepository, *dbent.Client) {
placeholder

	db, err := sql.Open("sqlite", "file:auth_service_identity_sync?mode=memory&cache=shared")
placeholder
	t.Cleanup(func() { _ = db.Close() placeholder)

	_, err = db.Exec("PRAGMA foreign_keys = ON")
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
		values: map[string]string{
			service.SettingKeyRegistrationEnabled: "true",
	placeholder,
placeholder, cfg)

	svc := service.NewAuthService(client, repo, nil, nil, cfg, settingSvc, nil, nil, nil, nil, nil)
	return svc, repo, client
placeholder

func TestAuthServiceRegisterDualWritesEmailIdentity(t *testing.T) {
	svc, _, client := newAuthServiceWithEnt(t)
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

func TestAuthServiceLoginTouchesLastLoginAt(t *testing.T) {
	svc, repo, client := newAuthServiceWithEnt(t)
	ctx := context.Background()

	user := &service.User{
		Email:       "login@example.com",
		Role:        service.RoleUser,
		Status:      service.StatusActive,
		Balance:     1,
		Concurrency: 1,
placeholder
	require.NoError(t, user.SetPassword("password"))
	require.NoError(t, repo.Create(ctx, user))

	old := time.Now().Add(-2 * time.Hour).UTC().Round(time.Second)
	_, err := client.User.UpdateOneID(user.ID).
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
	require.True(t, storedUser.LastLoginAt.After(old))
	require.True(t, storedUser.LastActiveAt.After(old))
placeholder
