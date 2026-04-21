package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

// BindEmailIdentity verifies and binds a local email/password identity to the current user.
func (s *AuthService) BindEmailIdentity(
	ctx context.Context,
	userID int64,
	email string,
	verifyCode string,
	password string,
) (*User, error) {
	if s == nil {
		return nil, ErrServiceUnavailable
placeholder

	normalizedEmail, err := normalizeEmailForIdentityBinding(email)
	if err != nil {
		return nil, err
placeholder
	if isReservedEmail(normalizedEmail) {
		return nil, ErrEmailReserved
placeholder
	if strings.TrimSpace(password) == "" {
		return nil, ErrPasswordRequired
placeholder
	if err := s.VerifyOAuthEmailCode(ctx, normalizedEmail, verifyCode); err != nil {
		return nil, err
placeholder

	currentUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
placeholder

	existingUser, err := s.userRepo.GetByEmail(ctx, normalizedEmail)
	switch {
	case err == nil && existingUser != nil && existingUser.ID != userID:
		return nil, ErrEmailExists
	case err != nil && !errors.Is(err, ErrUserNotFound):
		return nil, ErrServiceUnavailable
placeholder

	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
placeholder

	firstRealEmailBind := !hasBindableEmailIdentitySubject(currentUser.Email)
	currentUser.Email = normalizedEmail
	currentUser.PasswordHash = hashedPassword
	if err := s.userRepo.Update(ctx, currentUser); err != nil {
		if errors.Is(err, ErrEmailExists) {
			return nil, ErrEmailExists
	placeholder
		return nil, ErrServiceUnavailable
placeholder

	if firstRealEmailBind {
		if err := s.ApplyProviderDefaultSettingsOnFirstBind(ctx, userID, "email"); err != nil {
			return nil, fmt.Errorf("apply email first bind defaults: %w", err)
	placeholder
placeholder

	return currentUser, nil
placeholder

// SendEmailIdentityBindCode sends a verification code for authenticated email binding flows.
func (s *AuthService) SendEmailIdentityBindCode(ctx context.Context, userID int64, email string) error {
	if s == nil {
		return ErrServiceUnavailable
placeholder

	normalizedEmail, err := normalizeEmailForIdentityBinding(email)
	if err != nil {
		return err
placeholder
	if isReservedEmail(normalizedEmail) {
		return ErrEmailReserved
placeholder
	if s.emailService == nil {
		return ErrServiceUnavailable
placeholder
	if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return ErrUserNotFound
	placeholder
		return ErrServiceUnavailable
placeholder

	existingUser, err := s.userRepo.GetByEmail(ctx, normalizedEmail)
	switch {
	case err == nil && existingUser != nil && existingUser.ID != userID:
		return ErrEmailExists
	case err != nil && !errors.Is(err, ErrUserNotFound):
		return ErrServiceUnavailable
placeholder

	siteName := "Sub2API"
	if s.settingService != nil {
		siteName = s.settingService.GetSiteName(ctx)
placeholder
	return s.emailService.SendVerifyCode(ctx, normalizedEmail, siteName)
placeholder

func normalizeEmailForIdentityBinding(email string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(email))
	if normalized == "" || len(normalized) > 255 {
		return "", infraerrors.BadRequest("INVALID_EMAIL", "invalid email")
placeholder
	if _, err := mail.ParseAddress(normalized); err != nil {
		return "", infraerrors.BadRequest("INVALID_EMAIL", "invalid email")
placeholder
	return normalized, nil
placeholder

func hasBindableEmailIdentitySubject(email string) bool {
	normalized := strings.ToLower(strings.TrimSpace(email))
	return normalized != "" && !isReservedEmail(normalized)
placeholder
