//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

type redeemCodeRepoStub struct {
	codesByCode map[string]*RedeemCode
	useCalls    []struct {
		id     int64
		userID int64
placeholder
	updateCalls []*RedeemCode
placeholder

func (s *redeemCodeRepoStub) Create(context.Context, *RedeemCode) error {
	panic("unexpected Create call")
placeholder

func (s *redeemCodeRepoStub) CreateBatch(context.Context, []RedeemCode) error {
	panic("unexpected CreateBatch call")
placeholder

func (s *redeemCodeRepoStub) GetByID(context.Context, int64) (*RedeemCode, error) {
	panic("unexpected GetByID call")
placeholder

func (s *redeemCodeRepoStub) GetByCode(_ context.Context, code string) (*RedeemCode, error) {
	if s.codesByCode == nil {
		return nil, ErrRedeemCodeNotFound
placeholder
	redeemCode, ok := s.codesByCode[code]
	if !ok {
		return nil, ErrRedeemCodeNotFound
placeholder
	cloned := *redeemCode
	return &cloned, nil
placeholder

func (s *redeemCodeRepoStub) Update(_ context.Context, code *RedeemCode) error {
	if code == nil {
		return nil
placeholder
	cloned := *code
	s.updateCalls = append(s.updateCalls, &cloned)
	if s.codesByCode == nil {
		s.codesByCode = make(map[string]*RedeemCode)
placeholder
	s.codesByCode[cloned.Code] = &cloned
	return nil
placeholder

func (s *redeemCodeRepoStub) Delete(context.Context, int64) error {
	panic("unexpected Delete call")
placeholder

func (s *redeemCodeRepoStub) Use(_ context.Context, id, userID int64) error {
	for code, redeemCode := range s.codesByCode {
		if redeemCode.ID != id {
			continue
	placeholder
		now := time.Now().UTC()
		redeemCode.Status = StatusUsed
		redeemCode.UsedBy = &userID
		redeemCode.UsedAt = &now
		s.codesByCode[code] = redeemCode
		s.useCalls = append(s.useCalls, struct {
			id     int64
			userID int64
	placeholder{id: id, userID: userIDplaceholder)
		return nil
placeholder
	return ErrRedeemCodeNotFound
placeholder

func (s *redeemCodeRepoStub) List(context.Context, pagination.PaginationParams) ([]RedeemCode, *pagination.PaginationResult, error) {
	panic("unexpected List call")
placeholder

func (s *redeemCodeRepoStub) ListWithFilters(context.Context, pagination.PaginationParams, string, string, string) ([]RedeemCode, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
placeholder

func (s *redeemCodeRepoStub) ListByUser(context.Context, int64, int) ([]RedeemCode, error) {
	panic("unexpected ListByUser call")
placeholder

func (s *redeemCodeRepoStub) ListByUserPaginated(context.Context, int64, pagination.PaginationParams, string) ([]RedeemCode, *pagination.PaginationResult, error) {
	panic("unexpected ListByUserPaginated call")
placeholder

func (s *redeemCodeRepoStub) SumPositiveBalanceByUser(context.Context, int64) (float64, error) {
	panic("unexpected SumPositiveBalanceByUser call")
placeholder

func newOAuthEmailFlowAuthService(
	userRepo UserRepository,
	redeemRepo RedeemCodeRepository,
	refreshTokenCache RefreshTokenCache,
	settings map[string]string,
	emailCache EmailCache,
) *AuthService {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:                   "test-secret",
			ExpireHour:               1,
			AccessTokenExpireMinutes: 60,
			RefreshTokenExpireDays:   7,
	placeholder,
		Default: config.DefaultConfig{
			UserBalance:     3.5,
			UserConcurrency: 2,
	placeholder,
placeholder

	settingService := NewSettingService(&settingRepoStub{values: settingsplaceholder, cfg)
	emailService := NewEmailService(&settingRepoStub{values: settingsplaceholder, emailCache)

	return NewAuthService(
		nil,
		userRepo,
		redeemRepo,
		refreshTokenCache,
		cfg,
		settingService,
		emailService,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
placeholder

func TestRegisterOAuthEmailAccountRollsBackCreatedUserWhenTokenPairGenerationFails(t *testing.T) {
	userRepo := &userRepoStub{nextID: 42placeholder
	redeemRepo := &redeemCodeRepoStub{
		codesByCode: map[string]*RedeemCode{
			"INVITE123": {
				ID:     7,
				Code:   "INVITE123",
				Type:   RedeemTypeInvitation,
				Status: StatusUnused,
		placeholder,
	placeholder,
placeholder
	emailCache := &emailCacheStub{
		data: &VerificationCodeData{
			Code:      "246810",
			Attempts:  0,
			CreatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(15 * time.Minute),
	placeholder,
placeholder
	authService := newOAuthEmailFlowAuthService(
		userRepo,
		redeemRepo,
		nil,
		map[string]string{
			SettingKeyRegistrationEnabled:   "true",
			SettingKeyInvitationCodeEnabled: "true",
			SettingKeyEmailVerifyEnabled:    "true",
	placeholder,
		emailCache,
	)

	tokenPair, user, err := authService.RegisterOAuthEmailAccount(
		context.Background(),
		"fresh@example.com",
		"secret-123",
		"246810",
		"INVITE123",
		"oidc",
	)

	require.Nil(t, tokenPair)
	require.Nil(t, user)
placeholder
	require.Contains(t, err.Error(), "generate token pair")
	require.Equal(t, []int64{42placeholder, userRepo.deletedIDs)
	require.Len(t, userRepo.created, 1)
	require.Empty(t, redeemRepo.useCalls)
	require.Empty(t, redeemRepo.updateCalls)
placeholder

func TestRegisterOAuthEmailAccountSetsNormalizedSignupSourceOnCreatedUser(t *testing.T) {
	userRepo := &userRepoStub{nextID: 42placeholder
	emailCache := &emailCacheStub{
		data: &VerificationCodeData{
			Code:      "246810",
			Attempts:  0,
			CreatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(15 * time.Minute),
	placeholder,
placeholder
	authService := newOAuthEmailFlowAuthService(
		userRepo,
		&redeemCodeRepoStub{placeholder,
		&refreshTokenCacheStub{placeholder,
		map[string]string{
			SettingKeyRegistrationEnabled: "true",
			SettingKeyEmailVerifyEnabled:  "true",
	placeholder,
		emailCache,
	)

	tokenPair, user, err := authService.RegisterOAuthEmailAccount(
		context.Background(),
		"fresh@example.com",
		"secret-123",
		"246810",
		"",
		" OIDC ",
	)

placeholder
	require.NotNil(t, tokenPair)
	require.NotNil(t, user)
	require.Len(t, userRepo.created, 1)
	require.Equal(t, "oidc", userRepo.created[0].SignupSource)
placeholder

func TestRegisterOAuthEmailAccountKeepsGitHubAndGoogleSignupSource(t *testing.T) {
	tests := []struct {
		name         string
		email        string
		signupSource string
		want         string
placeholder{
		{
			name:         "github",
			email:        "github@example.com",
			signupSource: " GitHub ",
			want:         "github",
	placeholder,
		{
			name:         "google",
			email:        "google@example.com",
			signupSource: " Google ",
			want:         "google",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &userRepoStub{nextID: 43placeholder
			emailCache := &emailCacheStub{
				data: &VerificationCodeData{
					Code:      "246810",
					Attempts:  0,
					CreatedAt: time.Now().UTC(),
					ExpiresAt: time.Now().UTC().Add(15 * time.Minute),
			placeholder,
		placeholder
			authService := newOAuthEmailFlowAuthService(
				userRepo,
				&redeemCodeRepoStub{placeholder,
				&refreshTokenCacheStub{placeholder,
				map[string]string{
					SettingKeyRegistrationEnabled: "true",
					SettingKeyEmailVerifyEnabled:  "true",
			placeholder,
				emailCache,
			)

			tokenPair, user, err := authService.RegisterOAuthEmailAccount(
				context.Background(),
				tt.email,
				"secret-123",
				"246810",
				"",
				tt.signupSource,
			)

		placeholder
			require.NotNil(t, tokenPair)
			require.NotNil(t, user)
			require.Len(t, userRepo.created, 1)
			require.Equal(t, tt.want, userRepo.created[0].SignupSource)
	placeholder)
placeholder
placeholder

func TestRegisterOAuthEmailAccountFallsBackUnknownSignupSourceToEmail(t *testing.T) {
	userRepo := &userRepoStub{nextID: 43placeholder
	emailCache := &emailCacheStub{
		data: &VerificationCodeData{
			Code:      "246810",
			Attempts:  0,
			CreatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(15 * time.Minute),
	placeholder,
placeholder
	authService := newOAuthEmailFlowAuthService(
		userRepo,
		&redeemCodeRepoStub{placeholder,
		&refreshTokenCacheStub{placeholder,
		map[string]string{
			SettingKeyRegistrationEnabled: "true",
			SettingKeyEmailVerifyEnabled:  "true",
	placeholder,
		emailCache,
	)

	tokenPair, user, err := authService.RegisterOAuthEmailAccount(
		context.Background(),
		"fallback@example.com",
		"secret-123",
		"246810",
		"",
		"unknown-provider",
	)

placeholder
	require.NotNil(t, tokenPair)
	require.NotNil(t, user)
	require.Len(t, userRepo.created, 1)
	require.Equal(t, "email", userRepo.created[0].SignupSource)
placeholder

func TestRollbackOAuthEmailAccountCreationRestoresInvitationUsage(t *testing.T) {
	userRepo := &userRepoStub{placeholder
	redeemRepo := &redeemCodeRepoStub{
		codesByCode: map[string]*RedeemCode{
			"INVITE123": {
				ID:     7,
				Code:   "INVITE123",
				Type:   RedeemTypeInvitation,
				Status: StatusUsed,
				UsedBy: func() *int64 {
					v := int64(42)
					return &v
			placeholder(),
				UsedAt: func() *time.Time {
					v := time.Now().UTC()
					return &v
			placeholder(),
		placeholder,
	placeholder,
placeholder
	authService := newOAuthEmailFlowAuthService(
		userRepo,
		redeemRepo,
		&refreshTokenCacheStub{placeholder,
		map[string]string{
			SettingKeyRegistrationEnabled:   "true",
			SettingKeyInvitationCodeEnabled: "true",
	placeholder,
		&emailCacheStub{placeholder,
	)

	err := authService.RollbackOAuthEmailAccountCreation(context.Background(), 42, "INVITE123")

placeholder
	require.Equal(t, []int64{42placeholder, userRepo.deletedIDs)
	require.Len(t, redeemRepo.updateCalls, 1)
	require.Equal(t, StatusUnused, redeemRepo.updateCalls[0].Status)
	require.Nil(t, redeemRepo.updateCalls[0].UsedBy)
	require.Nil(t, redeemRepo.updateCalls[0].UsedAt)
placeholder

func TestRollbackOAuthEmailAccountCreationPropagatesDeleteError(t *testing.T) {
	userRepo := &userRepoStub{deleteErr: errors.New("delete failed")placeholder
	authService := newOAuthEmailFlowAuthService(
		userRepo,
		&redeemCodeRepoStub{placeholder,
		&refreshTokenCacheStub{placeholder,
		map[string]string{
			SettingKeyRegistrationEnabled: "true",
	placeholder,
		&emailCacheStub{placeholder,
	)

	err := authService.RollbackOAuthEmailAccountCreation(context.Background(), 42, "")

placeholder
	require.Contains(t, err.Error(), "delete created oauth user")
placeholder
