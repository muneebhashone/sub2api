//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type settingRepoStub struct {
	values map[string]string
	err    error
placeholder

func (s *settingRepoStub) Get(ctx context.Context, key string) (*Setting, error) {
	panic("unexpected Get call")
placeholder

func (s *settingRepoStub) GetValue(ctx context.Context, key string) (string, error) {
	if s.err != nil {
		return "", s.err
placeholder
	if v, ok := s.values[key]; ok {
		return v, nil
placeholder
	return "", ErrSettingNotFound
placeholder

func (s *settingRepoStub) Set(ctx context.Context, key, value string) error {
	panic("unexpected Set call")
placeholder

func (s *settingRepoStub) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	panic("unexpected GetMultiple call")
placeholder

func (s *settingRepoStub) SetMultiple(ctx context.Context, settings map[string]string) error {
	panic("unexpected SetMultiple call")
placeholder

func (s *settingRepoStub) GetAll(ctx context.Context) (map[string]string, error) {
	panic("unexpected GetAll call")
placeholder

func (s *settingRepoStub) Delete(ctx context.Context, key string) error {
	panic("unexpected Delete call")
placeholder

type emailCacheStub struct {
	data *VerificationCodeData
	err  error
placeholder

func (s *emailCacheStub) GetVerificationCode(ctx context.Context, email string) (*VerificationCodeData, error) {
	if s.err != nil {
		return nil, s.err
placeholder
	return s.data, nil
placeholder

func (s *emailCacheStub) SetVerificationCode(ctx context.Context, email string, data *VerificationCodeData, ttl time.Duration) error {
	return nil
placeholder

func (s *emailCacheStub) DeleteVerificationCode(ctx context.Context, email string) error {
	return nil
placeholder

func newAuthService(repo *userRepoStub, settings map[string]string, emailCache EmailCache) *AuthService {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret",
			ExpireHour: 1,
	placeholder,
		Default: config.DefaultConfig{
			UserBalance:     3.5,
			UserConcurrency: 2,
	placeholder,
placeholder

	var settingService *SettingService
	if settings != nil {
		settingService = NewSettingService(&settingRepoStub{values: settingsplaceholder, cfg)
placeholder

	var emailService *EmailService
	if emailCache != nil {
		emailService = NewEmailService(&settingRepoStub{values: settingsplaceholder, emailCache)
placeholder

	return NewAuthService(
		repo,
		cfg,
		settingService,
		emailService,
		nil,
		nil,
	)
placeholder

func TestAuthService_Register_Disabled(t *testing.T) {
	repo := &userRepoStub{placeholder
	service := newAuthService(repo, map[string]string{
		SettingKeyRegistrationEnabled: "false",
placeholder, nil)

	_, _, err := service.Register(context.Background(), "user@test.com", "password")
	require.ErrorIs(t, err, ErrRegDisabled)
placeholder

func TestAuthService_Register_DisabledByDefault(t *testing.T) {
	// 当 settings 为 nil（设置项不存在）时，注册应该默认关闭
	repo := &userRepoStub{placeholder
	service := newAuthService(repo, nil, nil)

	_, _, err := service.Register(context.Background(), "user@test.com", "password")
	require.ErrorIs(t, err, ErrRegDisabled)
placeholder

func TestAuthService_Register_EmailVerifyEnabledButServiceNotConfigured(t *testing.T) {
	repo := &userRepoStub{placeholder
	// 邮件验证开启但 emailCache 为 nil（emailService 未配置）
	service := newAuthService(repo, map[string]string{
		SettingKeyRegistrationEnabled: "true",
		SettingKeyEmailVerifyEnabled:  "true",
placeholder, nil)

	// 应返回服务不可用错误，而不是允许绕过验证
	_, _, err := service.RegisterWithVerification(context.Background(), "user@test.com", "password", "any-code")
	require.ErrorIs(t, err, ErrServiceUnavailable)
placeholder

func TestAuthService_Register_EmailVerifyRequired(t *testing.T) {
	repo := &userRepoStub{placeholder
	cache := &emailCacheStub{placeholder // 配置 emailService
	service := newAuthService(repo, map[string]string{
		SettingKeyRegistrationEnabled: "true",
		SettingKeyEmailVerifyEnabled:  "true",
placeholder, cache)

	_, _, err := service.RegisterWithVerification(context.Background(), "user@test.com", "password", "")
	require.ErrorIs(t, err, ErrEmailVerifyRequired)
placeholder

func TestAuthService_Register_EmailVerifyInvalid(t *testing.T) {
	repo := &userRepoStub{placeholder
	cache := &emailCacheStub{
		data: &VerificationCodeData{Code: "expected", Attempts: 0placeholder,
placeholder
	service := newAuthService(repo, map[string]string{
		SettingKeyRegistrationEnabled: "true",
		SettingKeyEmailVerifyEnabled:  "true",
placeholder, cache)

	_, _, err := service.RegisterWithVerification(context.Background(), "user@test.com", "password", "wrong")
	require.ErrorIs(t, err, ErrInvalidVerifyCode)
	require.ErrorContains(t, err, "verify code")
placeholder

func TestAuthService_Register_EmailExists(t *testing.T) {
	repo := &userRepoStub{exists: trueplaceholder
	service := newAuthService(repo, map[string]string{
		SettingKeyRegistrationEnabled: "true",
placeholder, nil)

	_, _, err := service.Register(context.Background(), "user@test.com", "password")
	require.ErrorIs(t, err, ErrEmailExists)
placeholder

func TestAuthService_Register_CheckEmailError(t *testing.T) {
	repo := &userRepoStub{existsErr: errors.New("db down")placeholder
	service := newAuthService(repo, map[string]string{
		SettingKeyRegistrationEnabled: "true",
placeholder, nil)

	_, _, err := service.Register(context.Background(), "user@test.com", "password")
	require.ErrorIs(t, err, ErrServiceUnavailable)
placeholder

func TestAuthService_Register_ReservedEmail(t *testing.T) {
	repo := &userRepoStub{placeholder
	service := newAuthService(repo, nil, nil)

	_, _, err := service.Register(context.Background(), "linuxdo-123@linuxdo-connect.invalid", "password")
	require.ErrorIs(t, err, ErrEmailReserved)
placeholder

func TestAuthService_Register_CreateError(t *testing.T) {
	repo := &userRepoStub{createErr: errors.New("create failed")placeholder
	service := newAuthService(repo, map[string]string{
		SettingKeyRegistrationEnabled: "true",
placeholder, nil)

	_, _, err := service.Register(context.Background(), "user@test.com", "password")
	require.ErrorIs(t, err, ErrServiceUnavailable)
placeholder

func TestAuthService_Register_CreateEmailExistsRace(t *testing.T) {
	// 模拟竞态条件：ExistsByEmail 返回 false，但 Create 时因唯一约束失败
	repo := &userRepoStub{createErr: ErrEmailExistsplaceholder
	service := newAuthService(repo, map[string]string{
		SettingKeyRegistrationEnabled: "true",
placeholder, nil)

	_, _, err := service.Register(context.Background(), "user@test.com", "password")
	require.ErrorIs(t, err, ErrEmailExists)
placeholder

func TestAuthService_Register_Success(t *testing.T) {
	repo := &userRepoStub{nextID: 5placeholder
	service := newAuthService(repo, map[string]string{
		SettingKeyRegistrationEnabled: "true",
placeholder, nil)

	token, user, err := service.Register(context.Background(), "user@test.com", "password")
placeholder
	require.NotEmpty(t, token)
	require.NotNil(t, user)
	require.Equal(t, int64(5), user.ID)
	require.Equal(t, "user@test.com", user.Email)
	require.Equal(t, RoleUser, user.Role)
	require.Equal(t, StatusActive, user.Status)
	require.Equal(t, 3.5, user.Balance)
	require.Equal(t, 2, user.Concurrency)
	require.Len(t, repo.created, 1)
	require.True(t, user.CheckPassword("password"))
placeholder

func TestAuthService_ValidateToken_ExpiredReturnsClaimsWithError(t *testing.T) {
	repo := &userRepoStub{placeholder
	service := newAuthService(repo, nil, nil)

	// 创建用户并生成 token
	user := &User{
		ID:           1,
		Email:        "test@test.com",
		Role:         RoleUser,
		Status:       StatusActive,
		TokenVersion: 1,
placeholder
	token, err := service.GenerateToken(user)
placeholder

	// 验证有效 token
	claims, err := service.ValidateToken(token)
placeholder
	require.NotNil(t, claims)
	require.Equal(t, int64(1), claims.UserID)

	// 模拟过期 token（通过创建一个过期很久的 token）
	service.cfg.JWT.ExpireHour = -1 // 设置为负数使 token 立即过期
	expiredToken, err := service.GenerateToken(user)
placeholder
	service.cfg.JWT.ExpireHour = 1 // 恢复

	// 验证过期 token 应返回 claims 和 ErrTokenExpired
	claims, err = service.ValidateToken(expiredToken)
	require.ErrorIs(t, err, ErrTokenExpired)
	require.NotNil(t, claims, "claims should not be nil when token is expired")
	require.Equal(t, int64(1), claims.UserID)
	require.Equal(t, "test@test.com", claims.Email)
placeholder

func TestAuthService_RefreshToken_ExpiredTokenNoPanic(t *testing.T) {
	user := &User{
		ID:           1,
		Email:        "test@test.com",
		Role:         RoleUser,
		Status:       StatusActive,
		TokenVersion: 1,
placeholder
	repo := &userRepoStub{user: userplaceholder
	service := newAuthService(repo, nil, nil)

	// 创建过期 token
	service.cfg.JWT.ExpireHour = -1
	expiredToken, err := service.GenerateToken(user)
placeholder
	service.cfg.JWT.ExpireHour = 1

	// RefreshToken 使用过期 token 不应 panic
	require.NotPanics(t, func() {
		newToken, err := service.RefreshToken(context.Background(), expiredToken)
	placeholder
		require.NotEmpty(t, newToken)
placeholder)
placeholder
