package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// VerifyOAuthEmailCode verifies the locally entered email verification code for
// third-party signup and binding flows. This is intentionally independent from
// the global registration email verification toggle.
func (s *AuthService) VerifyOAuthEmailCode(ctx context.Context, email, verifyCode string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	verifyCode = strings.TrimSpace(verifyCode)

	if email == "" {
		return ErrEmailVerifyRequired
placeholder
	if verifyCode == "" {
		return ErrEmailVerifyRequired
placeholder
	if s == nil || s.emailService == nil {
		return ErrServiceUnavailable
placeholder
	return s.emailService.VerifyCode(ctx, email, verifyCode)
placeholder

// RegisterOAuthEmailAccount creates a local account from a third-party first
// login after the user has verified a local email address.
func (s *AuthService) RegisterOAuthEmailAccount(
	ctx context.Context,
	email string,
	password string,
	verifyCode string,
	invitationCode string,
	signupSource string,
) (*TokenPair, *User, error) {
	if s == nil {
		return nil, nil, ErrServiceUnavailable
placeholder
	if s.settingService == nil || !s.settingService.IsRegistrationEnabled(ctx) {
		return nil, nil, ErrRegDisabled
placeholder

	email = strings.TrimSpace(strings.ToLower(email))
	if isReservedEmail(email) {
		return nil, nil, ErrEmailReserved
placeholder
	if err := s.validateRegistrationEmailPolicy(ctx, email); err != nil {
		return nil, nil, err
placeholder
	if err := s.VerifyOAuthEmailCode(ctx, email, verifyCode); err != nil {
		return nil, nil, err
placeholder

	var invitationRedeemCode *RedeemCode
	if s.settingService.IsInvitationCodeEnabled(ctx) {
		if invitationCode == "" {
			return nil, nil, ErrInvitationCodeRequired
	placeholder
		redeemCode, err := s.redeemRepo.GetByCode(ctx, invitationCode)
		if err != nil {
			return nil, nil, ErrInvitationCodeInvalid
	placeholder
		if redeemCode.Type != RedeemTypeInvitation || redeemCode.Status != StatusUnused {
			return nil, nil, ErrInvitationCodeInvalid
	placeholder
		invitationRedeemCode = redeemCode
placeholder

	existsEmail, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, nil, ErrServiceUnavailable
placeholder
	if existsEmail {
		return nil, nil, ErrEmailExists
placeholder

	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, nil, fmt.Errorf("hash password: %w", err)
placeholder

	signupSource = strings.TrimSpace(strings.ToLower(signupSource))
	if signupSource == "" {
		signupSource = "email"
placeholder
	grantPlan := s.resolveSignupGrantPlan(ctx, signupSource)

	user := &User{
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         RoleUser,
		Balance:      grantPlan.Balance,
		Concurrency:  grantPlan.Concurrency,
		Status:       StatusActive,
placeholder

	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, ErrEmailExists) {
			return nil, nil, ErrEmailExists
	placeholder
		return nil, nil, ErrServiceUnavailable
placeholder

	s.postAuthUserBootstrap(ctx, user, signupSource, false)
	s.assignSubscriptions(ctx, user.ID, grantPlan.Subscriptions, "auto assigned by signup defaults")

	if invitationRedeemCode != nil {
		if err := s.redeemRepo.Use(ctx, invitationRedeemCode.ID, user.ID); err != nil {
			return nil, nil, ErrInvitationCodeInvalid
	placeholder
placeholder

	tokenPair, err := s.GenerateTokenPair(ctx, user, "")
	if err != nil {
		return nil, nil, fmt.Errorf("generate token pair: %w", err)
placeholder
	return tokenPair, user, nil
placeholder

// ValidatePasswordCredentials checks the local password without completing the
// login flow. This is used by pending third-party account adoption flows before
// the external identity has been bound.
func (s *AuthService) ValidatePasswordCredentials(ctx context.Context, email, password string) (*User, error) {
	if s == nil {
		return nil, ErrServiceUnavailable
placeholder

	user, err := s.userRepo.GetByEmail(ctx, strings.TrimSpace(strings.ToLower(email)))
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
	placeholder
		return nil, ErrServiceUnavailable
placeholder
	if !user.IsActive() {
		return nil, ErrUserNotActive
placeholder
	if !s.CheckPassword(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
placeholder
	return user, nil
placeholder

// RecordSuccessfulLogin updates last-login activity after a non-standard login
// flow finishes with a real session.
func (s *AuthService) RecordSuccessfulLogin(ctx context.Context, userID int64) {
	if s != nil && s.userRepo != nil && userID > 0 {
		user, err := s.userRepo.GetByID(ctx, userID)
		if err == nil {
			s.backfillEmailIdentityOnSuccessfulLogin(ctx, user)
	placeholder
placeholder
	s.touchUserLogin(ctx, userID)
placeholder
