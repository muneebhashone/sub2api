package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/mail"
	"strings"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/redeemcode"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

func normalizeOAuthSignupSource(signupSource string) string {
	signupSource = strings.TrimSpace(strings.ToLower(signupSource))
	switch signupSource {
	case "", "email":
		return "email"
	case "linuxdo", "wechat", "oidc", "github", "google", "dingtalk":
		return signupSource
	default:
		return "email"
placeholder
placeholder

// SendPendingOAuthVerifyCode sends a local verification code for pending OAuth
// account-creation flows without relying on the public registration gate.
func (s *AuthService) SendPendingOAuthVerifyCode(ctx context.Context, email string, locale ...string) (*SendVerifyCodeResult, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return nil, ErrEmailVerifyRequired
placeholder
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, ErrEmailVerifyRequired
placeholder
	if isReservedEmail(email) {
		return nil, ErrEmailReserved
placeholder
	if s == nil || s.emailService == nil {
		return nil, ErrServiceUnavailable
placeholder
	if err := s.validateRegistrationEmailQuota(ctx, email); err != nil {
		return nil, err
placeholder

	siteName := "Sub2API"
	if s.settingService != nil {
		siteName = s.settingService.GetSiteName(ctx)
placeholder
	if err := s.emailService.SendVerifyCode(ctx, email, siteName, firstEmailLocale(locale)); err != nil {
		return nil, err
placeholder
	return &SendVerifyCodeResult{
		Countdown: int(verifyCodeCooldown / time.Second),
placeholder, nil
placeholder

func (s *AuthService) validateOAuthRegistrationInvitation(ctx context.Context, invitationCode string) (*RedeemCode, error) {
	if s == nil || s.settingService == nil || !s.settingService.IsInvitationCodeEnabled(ctx) {
		return nil, nil
placeholder
	if s.redeemRepo == nil && s.oauthEmailFlowClient(ctx) == nil {
		return nil, ErrServiceUnavailable
placeholder

	invitationCode = strings.TrimSpace(invitationCode)
	if invitationCode == "" {
		return nil, ErrInvitationCodeRequired
placeholder

	redeemCode, err := s.loadOAuthRegistrationInvitation(ctx, invitationCode)
	if err != nil {
		return nil, ErrInvitationCodeInvalid
placeholder
	if redeemCode.Type != RedeemTypeInvitation || !redeemCode.CanUse() {
		return nil, ErrInvitationCodeInvalid
placeholder
	return redeemCode, nil
placeholder

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
	if s.settingService == nil || (!s.settingService.IsRegistrationEnabled(ctx) && !s.canBypassRegistrationDisabledForOAuth(ctx, signupSource)) {
		return nil, nil, ErrRegDisabled
placeholder

	email = strings.TrimSpace(strings.ToLower(email))
	if isReservedEmail(email) {
		return nil, nil, ErrEmailReserved
placeholder
	if err := s.VerifyOAuthEmailCode(ctx, email, verifyCode); err != nil {
		slog.Error("oauth email register: verify code failed", "email", email, "error", err.Error())
		return nil, nil, err
placeholder

	if _, err := s.validateOAuthRegistrationInvitation(ctx, invitationCode); err != nil {
		slog.Error("oauth email register: invitation failed", "email", email, "error", err.Error())
		return nil, nil, err
placeholder

	// 含 +别名 / Gmail 点号 / FQDN 根点变体归一化：该路径同样发放注册赠额，不能被单个收件箱刷号。
	existsEmail, err := s.existsByEmailOrAlias(ctx, email)
	if err != nil {
		slog.Error("oauth email register: ExistsByEmail failed", "email", email, "error", err.Error())
		return nil, nil, ErrServiceUnavailable
placeholder
	if existsEmail {
		return nil, nil, ErrEmailExists
placeholder
	if err := s.validateRegistrationEmailQuota(ctx, email); err != nil {
		slog.Error("oauth email register: policy rejected", "email", email, "error", err.Error())
		return nil, nil, err
placeholder

	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, nil, fmt.Errorf("hash password: %w", err)
placeholder

	signupSource = normalizeOAuthSignupSource(signupSource)
	grantPlan := s.resolveSignupGrantPlan(ctx, signupSource)

	user := &User{
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         RoleUser,
		Balance:      grantPlan.Balance,
		Concurrency:  grantPlan.Concurrency,
		Status:       StatusActive,
		SignupSource: signupSource,
placeholder

	if err := s.createUserWithRegistrationEmailGuard(ctx, user); err != nil {
		switch {
		case errors.Is(err, ErrEmailExists):
			return nil, nil, ErrEmailExists
		case errors.Is(err, ErrEmailDomainRegistrationLimit):
			return nil, nil, ErrEmailDomainRegistrationLimit
		default:
			slog.Error("oauth email register: userRepo.Create failed", "email", email, "signup_source", signupSource, "error", err.Error())
			return nil, nil, ErrServiceUnavailable
	placeholder
placeholder

	tokenPair, err := s.GenerateTokenPair(ctx, user, "")
	if err != nil {
		_ = s.RollbackOAuthEmailAccountCreation(ctx, user.ID, "")
		return nil, nil, fmt.Errorf("generate token pair: %w", err)
placeholder
	return tokenPair, user, nil
placeholder

// RegisterVerifiedOAuthEmailAccount creates a local account from an OAuth
// provider that has already returned a verified email address.
func (s *AuthService) RegisterVerifiedOAuthEmailAccount(
	ctx context.Context,
	email string,
	password string,
	invitationCode string,
	signupSource string,
) (*TokenPair, *User, error) {
	if s == nil {
		return nil, nil, ErrServiceUnavailable
placeholder
	if s.settingService == nil || (!s.settingService.IsRegistrationEnabled(ctx) && !s.canBypassRegistrationDisabledForOAuth(ctx, signupSource)) {
		return nil, nil, ErrRegDisabled
placeholder

	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || len(email) > 255 {
		return nil, nil, ErrEmailVerifyRequired
placeholder
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, nil, ErrEmailVerifyRequired
placeholder
	if isReservedEmail(email) {
		return nil, nil, ErrEmailReserved
placeholder
	if strings.TrimSpace(password) == "" {
		return nil, nil, infraerrors.BadRequest("PASSWORD_REQUIRED", "password is required")
placeholder
	if _, err := s.validateOAuthRegistrationInvitation(ctx, invitationCode); err != nil {
		return nil, nil, err
placeholder

	// 与本地注册同口径：同一收件箱的别名变体不能各自建号（该路径也发放注册赠额）。
	existsEmail, err := s.existsByEmailOrAlias(ctx, email)
	if err != nil {
		return nil, nil, ErrServiceUnavailable
placeholder
	if existsEmail {
		return nil, nil, ErrEmailExists
placeholder
	if err := s.validateRegistrationEmailQuota(ctx, email); err != nil {
		return nil, nil, err
placeholder

	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, nil, fmt.Errorf("hash password: %w", err)
placeholder

	signupSource = normalizeOAuthSignupSource(signupSource)
	grantPlan := s.resolveSignupGrantPlan(ctx, signupSource)
	var defaultRPMLimit int
	if s.settingService != nil {
		defaultRPMLimit = s.settingService.GetDefaultUserRPMLimit(ctx)
placeholder
	user := &User{
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         RoleUser,
		Balance:      grantPlan.Balance,
		Concurrency:  grantPlan.Concurrency,
		RPMLimit:     defaultRPMLimit,
		Status:       StatusActive,
		SignupSource: signupSource,
placeholder

	if err := s.createUserWithRegistrationEmailGuard(ctx, user); err != nil {
		switch {
		case errors.Is(err, ErrEmailExists):
			return nil, nil, ErrEmailExists
		case errors.Is(err, ErrEmailDomainRegistrationLimit):
			return nil, nil, ErrEmailDomainRegistrationLimit
		default:
			return nil, nil, ErrServiceUnavailable
	placeholder
placeholder

	tokenPair, err := s.GenerateTokenPair(ctx, user, "")
	if err != nil {
		_ = s.RollbackOAuthEmailAccountCreation(ctx, user.ID, "")
		return nil, nil, fmt.Errorf("generate token pair: %w", err)
placeholder
	return tokenPair, user, nil
placeholder

// FinalizeOAuthEmailAccount applies invitation usage and normal signup bootstrap
// only after the pending OAuth flow has fully reached its last reversible step.
func (s *AuthService) FinalizeOAuthEmailAccount(
	ctx context.Context,
	user *User,
	invitationCode string,
	signupSource string,
	affiliateCode string,
) error {
	if s == nil || user == nil || user.ID <= 0 {
		return ErrServiceUnavailable
placeholder

	signupSource = normalizeOAuthSignupSource(signupSource)
	invitationRedeemCode, err := s.validateOAuthRegistrationInvitation(ctx, invitationCode)
	if err != nil {
		return err
placeholder
	if invitationRedeemCode != nil {
		if err := s.useOAuthRegistrationInvitation(ctx, invitationRedeemCode.ID, user.ID); err != nil {
			return ErrInvitationCodeInvalid
	placeholder
placeholder

	s.updateOAuthSignupSource(ctx, user.ID, signupSource)
	grantPlan := s.resolveSignupGrantPlan(ctx, signupSource)
	s.assignSubscriptions(ctx, user.ID, grantPlan.Subscriptions, "auto assigned by signup defaults")
	// snapshot user × platform quota（fail-open）
	_ = s.snapshotPlatformQuotaDefaults(ctx, user.ID, &grantPlan)
	s.bindOAuthAffiliate(ctx, user.ID, affiliateCode)
	return nil
placeholder

// RollbackOAuthEmailAccountCreation removes a partially-created local account
// and restores any invitation code already consumed by that account.
func (s *AuthService) RollbackOAuthEmailAccountCreation(ctx context.Context, userID int64, invitationCode string) error {
	if s == nil || s.userRepo == nil || userID <= 0 {
		return ErrServiceUnavailable
placeholder
	if err := s.restoreOAuthRegistrationInvitation(ctx, invitationCode, userID); err != nil {
		return err
placeholder
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("delete created oauth user: %w", err)
placeholder
	return nil
placeholder

func (s *AuthService) restoreOAuthRegistrationInvitation(ctx context.Context, invitationCode string, userID int64) error {
	if s == nil || s.settingService == nil || !s.settingService.IsInvitationCodeEnabled(ctx) {
		return nil
placeholder
	if s.redeemRepo == nil && s.oauthEmailFlowClient(ctx) == nil {
		return ErrServiceUnavailable
placeholder

	invitationCode = strings.TrimSpace(invitationCode)
	if invitationCode == "" || userID <= 0 {
		return nil
placeholder

	redeemCode, err := s.loadOAuthRegistrationInvitation(ctx, invitationCode)
	if err != nil {
		if errors.Is(err, ErrRedeemCodeNotFound) {
			return nil
	placeholder
		return fmt.Errorf("load invitation code: %w", err)
placeholder
	if redeemCode.Type != RedeemTypeInvitation || redeemCode.Status != StatusUsed || redeemCode.UsedBy == nil || *redeemCode.UsedBy != userID {
		return nil
placeholder

	redeemCode.Status = StatusUnused
	redeemCode.UsedBy = nil
	redeemCode.UsedAt = nil
	if err := s.updateOAuthRegistrationInvitation(ctx, redeemCode); err != nil {
		return fmt.Errorf("restore invitation code: %w", err)
placeholder
	return nil
placeholder

func (s *AuthService) oauthEmailFlowClient(ctx context.Context) *dbent.Client {
	if s == nil || s.entClient == nil {
		return nil
placeholder
	if tx := dbent.TxFromContext(ctx); tx != nil {
		return tx.Client()
placeholder
	return s.entClient
placeholder

func (s *AuthService) loadOAuthRegistrationInvitation(ctx context.Context, invitationCode string) (*RedeemCode, error) {
	if client := s.oauthEmailFlowClient(ctx); client != nil {
		entity, err := client.RedeemCode.Query().Where(redeemcode.CodeEQ(invitationCode)).Only(ctx)
		if err != nil {
			if dbent.IsNotFound(err) {
				return nil, ErrRedeemCodeNotFound
		placeholder
			return nil, err
	placeholder
		return &RedeemCode{
			ID:           entity.ID,
			Code:         entity.Code,
			Type:         entity.Type,
			Value:        entity.Value,
			Status:       entity.Status,
			UsedBy:       entity.UsedBy,
			UsedAt:       entity.UsedAt,
			Notes:        oauthEmailFlowStringValue(entity.Notes),
			CreatedAt:    entity.CreatedAt,
			ExpiresAt:    entity.ExpiresAt,
			GroupID:      entity.GroupID,
			ValidityDays: entity.ValidityDays,
	placeholder, nil
placeholder
	return s.redeemRepo.GetByCode(ctx, invitationCode)
placeholder

func (s *AuthService) useOAuthRegistrationInvitation(ctx context.Context, invitationID, userID int64) error {
	if client := s.oauthEmailFlowClient(ctx); client != nil {
		affected, err := client.RedeemCode.Update().
			Where(
				redeemcode.IDEQ(invitationID),
				redeemcode.StatusEQ(StatusUnused),
				redeemcode.Or(redeemcode.ExpiresAtIsNil(), redeemcode.ExpiresAtGT(time.Now().UTC())),
			).
			SetStatus(StatusUsed).
			SetUsedBy(userID).
			SetUsedAt(time.Now().UTC()).
			Save(ctx)
		if err != nil {
			return err
	placeholder
		if affected == 0 {
			return ErrRedeemCodeUsed
	placeholder
		return nil
placeholder
	return s.redeemRepo.Use(ctx, invitationID, userID)
placeholder

func (s *AuthService) updateOAuthRegistrationInvitation(ctx context.Context, code *RedeemCode) error {
	if code == nil {
		return nil
placeholder
	if client := s.oauthEmailFlowClient(ctx); client != nil {
		update := client.RedeemCode.UpdateOneID(code.ID).
			SetCode(code.Code).
			SetType(code.Type).
			SetValue(code.Value).
			SetStatus(code.Status).
			SetNotes(code.Notes).
			SetValidityDays(code.ValidityDays)
		if code.ExpiresAt != nil {
			update = update.SetExpiresAt(*code.ExpiresAt)
	placeholder else {
			update = update.ClearExpiresAt()
	placeholder
		if code.UsedBy != nil {
			update = update.SetUsedBy(*code.UsedBy)
	placeholder else {
			update = update.ClearUsedBy()
	placeholder
		if code.UsedAt != nil {
			update = update.SetUsedAt(*code.UsedAt)
	placeholder else {
			update = update.ClearUsedAt()
	placeholder
		if code.GroupID != nil {
			update = update.SetGroupID(*code.GroupID)
	placeholder else {
			update = update.ClearGroupID()
	placeholder
		_, err := update.Save(ctx)
		return err
placeholder
	return s.redeemRepo.Update(ctx, code)
placeholder

func (s *AuthService) updateOAuthSignupSource(ctx context.Context, userID int64, signupSource string) {
	client := s.oauthEmailFlowClient(ctx)
	if client == nil || userID <= 0 || strings.TrimSpace(signupSource) == "" {
		return
placeholder
	_ = client.User.UpdateOneID(userID).SetSignupSource(signupSource).Exec(ctx)
placeholder

func oauthEmailFlowStringValue(value *string) string {
	if value == nil {
		return ""
placeholder
	return *value
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
		if err == nil && user != nil && !isReservedEmail(user.Email) {
			s.backfillEmailIdentityOnSuccessfulLogin(ctx, user)
	placeholder
placeholder
	s.touchUserLogin(ctx, userID)
placeholder
