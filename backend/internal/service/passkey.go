package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

const (
	passkeySessionRegistration = "registration"
	passkeySessionLogin        = "login"
	passkeySessionTTL          = 5 * time.Minute
	defaultPasskeyName         = "Passkey"
	maxPasskeyNameLength       = 100
)

var (
	ErrPasskeysDisabled = infraerrors.Forbidden("PASSKEY_DISABLED", "passkey authentication is not enabled")
	ErrPasskeyNotFound  = infraerrors.NotFound("PASSKEY_NOT_FOUND", "passkey not found")
	ErrPasskeyExists    = infraerrors.Conflict("PASSKEY_ALREADY_EXISTS", "this passkey is already registered")
	ErrPasskeySession   = infraerrors.BadRequest("PASSKEY_SESSION_INVALID", "passkey session is invalid or expired")
	ErrPasskeyVerify    = infraerrors.Unauthorized("PASSKEY_VERIFICATION_FAILED", "passkey verification failed")
)

// PasskeyCredentialRecord is the persistence representation used by the
// WebAuthn service. Credential contains the complete WebAuthn credential record
// so future library versions can continue to validate and update it.
type PasskeyCredentialRecord struct {
	ID         int64
	UserID     int64
	UserHandle []byte
	Name       string
	Credential webauthn.Credential
	LastUsedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
placeholder

type PasskeyRepository interface {
	EnsureUserHandle(ctx context.Context, userID int64, candidate []byte) ([]byte, error)
	GetUserHandle(ctx context.Context, userID int64) ([]byte, error)
	GetByCredentialID(ctx context.Context, credentialID []byte) (*PasskeyCredentialRecord, error)
	ListByUserID(ctx context.Context, userID int64) ([]PasskeyCredentialRecord, error)
	Create(ctx context.Context, record *PasskeyCredentialRecord) (*PasskeyCredentialRecord, error)
	UpdateCredential(ctx context.Context, userID int64, credential *webauthn.Credential, usedAt time.Time) error
	Rename(ctx context.Context, userID, credentialID int64, name string) error
	Delete(ctx context.Context, userID, credentialID int64) error
placeholder

type PasskeySession struct {
	Kind     string               `json:"kind"`
	UserID   int64                `json:"user_id,omitempty"`
	WebAuthn webauthn.SessionData `json:"webauthn"`
placeholder

type PasskeySessionStore interface {
	Store(ctx context.Context, session *PasskeySession, ttl time.Duration) (string, error)
	Consume(ctx context.Context, token string) (*PasskeySession, error)
placeholder

type PasskeyCredentialSummary struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	Backup     bool       `json:"backup"`
placeholder

type passkeyUser struct {
	account     *User
	handle      []byte
	credentials []webauthn.Credential
placeholder

func (u *passkeyUser) WebAuthnID() []byte {
	return u.handle
placeholder

func (u *passkeyUser) WebAuthnName() string {
	return u.account.Email
placeholder

func (u *passkeyUser) WebAuthnDisplayName() string {
	if name := strings.TrimSpace(u.account.Username); name != "" {
		return name
placeholder
	return u.account.Email
placeholder

func (u *passkeyUser) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
placeholder

type PasskeyService struct {
	enabled  bool
	webAuthn *webauthn.WebAuthn
	repo     PasskeyRepository
	sessions PasskeySessionStore
	userRepo UserRepository
placeholder

func NewPasskeyService(
	cfg *config.Config,
	repo PasskeyRepository,
	sessions PasskeySessionStore,
	userRepo UserRepository,
) (*PasskeyService, error) {
	s := &PasskeyService{
		repo:     repo,
		sessions: sessions,
		userRepo: userRepo,
placeholder
	if cfg == nil || !cfg.WebAuthn.Enabled {
		return s, nil
placeholder

	instance, err := webauthn.New(&webauthn.Config{
		RPDisplayName: cfg.WebAuthn.RPDisplayName,
		RPID:          cfg.WebAuthn.RPID,
		RPOrigins:     cfg.WebAuthn.RPOrigins,
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			ResidentKey:      protocol.ResidentKeyRequirementRequired,
			UserVerification: protocol.VerificationRequired,
	placeholder,
placeholder)
	if err != nil {
		return nil, fmt.Errorf("initialize WebAuthn: %w", err)
placeholder
	s.enabled = true
	s.webAuthn = instance
	return s, nil
placeholder

func (s *PasskeyService) Enabled() bool {
	return s != nil && s.enabled && s.webAuthn != nil
placeholder

func (s *PasskeyService) requireEnabled() error {
	if !s.Enabled() {
		return ErrPasskeysDisabled
placeholder
	return nil
placeholder

func (s *PasskeyService) BeginRegistration(
	ctx context.Context,
	userID int64,
) (creation *protocol.CredentialCreation, sessionToken string, err error) {
	if err = s.requireEnabled(); err != nil {
		return nil, "", err
placeholder
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, "", err
placeholder
	if !user.IsActive() {
		return nil, "", ErrUserNotActive
placeholder

	candidate := make([]byte, 32)
	if _, err = rand.Read(candidate); err != nil {
		return nil, "", fmt.Errorf("generate passkey user handle: %w", err)
placeholder
	handle, err := s.repo.EnsureUserHandle(ctx, userID, candidate)
	if err != nil {
		return nil, "", err
placeholder
	waUser, err := s.loadWebAuthnUser(ctx, user, handle)
	if err != nil {
		return nil, "", err
placeholder

	creation, session, err := s.webAuthn.BeginRegistration(
		waUser,
		webauthn.WithResidentKeyRequirement(protocol.ResidentKeyRequirementRequired),
		webauthn.WithExclusions(webauthn.Credentials(waUser.credentials).CredentialDescriptors()),
		webauthn.WithExtensions(protocol.AuthenticationExtensions{"credProps": trueplaceholder),
	)
	if err != nil {
		return nil, "", fmt.Errorf("begin passkey registration: %w", err)
placeholder
	sessionToken, err = s.sessions.Store(ctx, &PasskeySession{
		Kind:     passkeySessionRegistration,
		UserID:   userID,
		WebAuthn: *session,
placeholder, passkeySessionTTL)
	if err != nil {
		return nil, "", err
placeholder
	return creation, sessionToken, nil
placeholder

func (s *PasskeyService) FinishRegistration(
	ctx context.Context,
	userID int64,
	sessionToken, name string,
	request *http.Request,
) (*PasskeyCredentialSummary, error) {
	if err := s.requireEnabled(); err != nil {
		return nil, err
placeholder
	session, err := s.sessions.Consume(ctx, sessionToken)
	if err != nil {
		return nil, err
placeholder
	if session == nil || session.Kind != passkeySessionRegistration || session.UserID != userID {
		return nil, ErrPasskeySession
placeholder
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
placeholder
	handle, err := s.repo.GetUserHandle(ctx, userID)
	if err != nil {
		return nil, err
placeholder
	waUser, err := s.loadWebAuthnUser(ctx, user, handle)
	if err != nil {
		return nil, err
placeholder
	credential, err := s.webAuthn.FinishRegistration(waUser, session.WebAuthn, request)
	if err != nil {
		return nil, ErrPasskeyVerify
placeholder

	record, err := s.repo.Create(ctx, &PasskeyCredentialRecord{
		UserID:     userID,
		UserHandle: handle,
		Name:       normalizePasskeyName(name),
		Credential: *credential,
placeholder)
	if err != nil {
		return nil, err
placeholder
	return passkeySummary(record), nil
placeholder

func (s *PasskeyService) BeginLogin(
	ctx context.Context,
) (assertion *protocol.CredentialAssertion, sessionToken string, err error) {
	if err = s.requireEnabled(); err != nil {
		return nil, "", err
placeholder
	assertion, session, err := s.webAuthn.BeginDiscoverableLogin(
		webauthn.WithUserVerification(protocol.VerificationRequired),
	)
	if err != nil {
		return nil, "", fmt.Errorf("begin passkey login: %w", err)
placeholder
	sessionToken, err = s.sessions.Store(ctx, &PasskeySession{
		Kind:     passkeySessionLogin,
		WebAuthn: *session,
placeholder, passkeySessionTTL)
	if err != nil {
		return nil, "", err
placeholder
	return assertion, sessionToken, nil
placeholder

func (s *PasskeyService) FinishLogin(
	ctx context.Context,
	sessionToken string,
	request *http.Request,
) (*User, error) {
	if err := s.requireEnabled(); err != nil {
		return nil, err
placeholder
	session, err := s.sessions.Consume(ctx, sessionToken)
	if err != nil {
		return nil, err
placeholder
	if session == nil || session.Kind != passkeySessionLogin {
		return nil, ErrPasskeySession
placeholder

	handler := func(rawID, userHandle []byte) (webauthn.User, error) {
		record, lookupErr := s.repo.GetByCredentialID(ctx, rawID)
		if lookupErr != nil || record == nil || !bytes.Equal(record.UserHandle, userHandle) {
			return nil, ErrPasskeyVerify
	placeholder
		account, lookupErr := s.userRepo.GetByID(ctx, record.UserID)
		if lookupErr != nil || account == nil || !account.IsActive() {
			return nil, ErrPasskeyVerify
	placeholder
		return s.loadWebAuthnUser(ctx, account, record.UserHandle)
placeholder

	validatedUser, credential, err := s.webAuthn.FinishPasskeyLogin(handler, session.WebAuthn, request)
	if err != nil {
		return nil, ErrPasskeyVerify
placeholder
	waUser, ok := validatedUser.(*passkeyUser)
	if !ok || waUser.account == nil {
		return nil, ErrPasskeyVerify
placeholder
	if err = s.repo.UpdateCredential(ctx, waUser.account.ID, credential, time.Now().UTC()); err != nil {
		return nil, err
placeholder
	return waUser.account, nil
placeholder

func (s *PasskeyService) List(ctx context.Context, userID int64) ([]PasskeyCredentialSummary, error) {
	if err := s.requireEnabled(); err != nil {
		return nil, err
placeholder
	records, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
placeholder
	result := make([]PasskeyCredentialSummary, 0, len(records))
	for i := range records {
		result = append(result, *passkeySummary(&records[i]))
placeholder
	return result, nil
placeholder

func (s *PasskeyService) Rename(ctx context.Context, userID, credentialID int64, name string) error {
	if err := s.requireEnabled(); err != nil {
		return err
placeholder
	return s.repo.Rename(ctx, userID, credentialID, normalizePasskeyName(name))
placeholder

func (s *PasskeyService) Delete(ctx context.Context, userID, credentialID int64) error {
	if err := s.requireEnabled(); err != nil {
		return err
placeholder
	return s.repo.Delete(ctx, userID, credentialID)
placeholder

func (s *PasskeyService) loadWebAuthnUser(
	ctx context.Context,
	user *User,
	handle []byte,
) (*passkeyUser, error) {
	records, err := s.repo.ListByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
placeholder
	credentials := make([]webauthn.Credential, 0, len(records))
	for i := range records {
		credentials = append(credentials, records[i].Credential)
placeholder
	return &passkeyUser{account: user, handle: handle, credentials: credentialsplaceholder, nil
placeholder

func normalizePasskeyName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return defaultPasskeyName
placeholder
	runes := []rune(name)
	if len(runes) > maxPasskeyNameLength {
		name = string(runes[:maxPasskeyNameLength])
placeholder
	return name
placeholder

func passkeySummary(record *PasskeyCredentialRecord) *PasskeyCredentialSummary {
	return &PasskeyCredentialSummary{
		ID:         record.ID,
		Name:       record.Name,
		CreatedAt:  record.CreatedAt,
		LastUsedAt: record.LastUsedAt,
		Backup:     record.Credential.Flags.BackupState,
placeholder
placeholder

var _ webauthn.User = (*passkeyUser)(nil)
