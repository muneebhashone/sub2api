package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/fnv"
	"sort"
	"strings"
	"sync"
	"time"

	"entgo.io/ent/dialect"
	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/identityadoptiondecision"
	"github.com/Wei-Shaw/sub2api/ent/pendingauthsession"
	dbpredicate "github.com/Wei-Shaw/sub2api/ent/predicate"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"

	entsql "entgo.io/ent/dialect/sql"
)

var (
	ErrPendingAuthSessionNotFound = infraerrors.NotFound("PENDING_AUTH_SESSION_NOT_FOUND", "pending auth session not found")
	ErrPendingAuthSessionExpired  = infraerrors.Unauthorized("PENDING_AUTH_SESSION_EXPIRED", "pending auth session has expired")
	ErrPendingAuthSessionConsumed = infraerrors.Unauthorized("PENDING_AUTH_SESSION_CONSUMED", "pending auth session has already been used")
	ErrPendingAuthCodeInvalid     = infraerrors.Unauthorized("PENDING_AUTH_CODE_INVALID", "pending auth completion code is invalid")
	ErrPendingAuthCodeExpired     = infraerrors.Unauthorized("PENDING_AUTH_CODE_EXPIRED", "pending auth completion code has expired")
	ErrPendingAuthCodeConsumed    = infraerrors.Unauthorized("PENDING_AUTH_CODE_CONSUMED", "pending auth completion code has already been used")
	ErrPendingAuthBrowserMismatch = infraerrors.Unauthorized("PENDING_AUTH_BROWSER_MISMATCH", "pending auth completion code does not match this browser session")
)

const (
	defaultPendingAuthTTL           = 15 * time.Minute
	defaultPendingAuthCompletionTTL = 5 * time.Minute
)

type PendingAuthIdentityKey struct {
	ProviderType    string
	ProviderKey     string
	ProviderSubject string
placeholder

type CreatePendingAuthSessionInput struct {
	SessionToken             string
	Intent                   string
	Identity                 PendingAuthIdentityKey
	TargetUserID             *int64
	RedirectTo               string
	ResolvedEmail            string
	RegistrationPasswordHash string
	BrowserSessionKey        string
	UpstreamIdentityClaims   map[string]any
	LocalFlowState           map[string]any
	ExpiresAt                time.Time
placeholder

type IssuePendingAuthCompletionCodeInput struct {
	PendingAuthSessionID int64
	BrowserSessionKey    string
	TTL                  time.Duration
placeholder

type IssuePendingAuthCompletionCodeResult struct {
	Code      string
	ExpiresAt time.Time
placeholder

type PendingIdentityAdoptionDecisionInput struct {
	PendingAuthSessionID int64
	IdentityID           *int64
	AdoptDisplayName     bool
	AdoptAvatar          bool
placeholder

type AuthPendingIdentityService struct {
	entClient *dbent.Client
placeholder

var authPendingIdentityScopedKeyLocks = newAuthPendingIdentityScopedKeyLockRegistry()

type authPendingIdentityScopedKeyLockRegistry struct {
	mu    sync.Mutex
	locks map[string]*authPendingIdentityScopedKeyLockEntry
placeholder

type authPendingIdentityScopedKeyLockEntry struct {
	mu   sync.Mutex
	refs int
placeholder

func newAuthPendingIdentityScopedKeyLockRegistry() *authPendingIdentityScopedKeyLockRegistry {
	return &authPendingIdentityScopedKeyLockRegistry{
		locks: make(map[string]*authPendingIdentityScopedKeyLockEntry),
placeholder
placeholder

func (r *authPendingIdentityScopedKeyLockRegistry) lock(keys ...string) func() {
	normalized := normalizeAuthPendingIdentityLockKeys(keys...)
	if len(normalized) == 0 {
		return func() {placeholder
placeholder

	entries := make([]*authPendingIdentityScopedKeyLockEntry, 0, len(normalized))
	r.mu.Lock()
	for _, key := range normalized {
		entry := r.locks[key]
		if entry == nil {
			entry = &authPendingIdentityScopedKeyLockEntry{placeholder
			r.locks[key] = entry
	placeholder
		entry.refs++
		entries = append(entries, entry)
placeholder
	r.mu.Unlock()

	for _, entry := range entries {
		entry.mu.Lock()
placeholder

	return func() {
		for i := len(entries) - 1; i >= 0; i-- {
			entries[i].mu.Unlock()
	placeholder

		r.mu.Lock()
		defer r.mu.Unlock()
		for idx, key := range normalized {
			entry := entries[idx]
			entry.refs--
			if entry.refs == 0 {
				delete(r.locks, key)
		placeholder
	placeholder
placeholder
placeholder

func normalizeAuthPendingIdentityLockKeys(keys ...string) []string {
	if len(keys) == 0 {
		return nil
placeholder

	deduped := make(map[string]struct{placeholder, len(keys))
	for _, key := range keys {
		trimmed := strings.TrimSpace(key)
		if trimmed == "" {
			continue
	placeholder
		deduped[trimmed] = struct{placeholder{placeholder
placeholder
	if len(deduped) == 0 {
		return nil
placeholder

	normalized := make([]string, 0, len(deduped))
	for key := range deduped {
		normalized = append(normalized, key)
placeholder
	sort.Strings(normalized)
	return normalized
placeholder

func authPendingIdentityAdvisoryLockHash(key string) int64 {
	hasher := fnv.New64a()
	_, _ = hasher.Write([]byte(key))
	return int64(hasher.Sum64())
placeholder

func lockAuthPendingIdentityKeys(ctx context.Context, client *dbent.Client, keys ...string) (func(), error) {
	release := authPendingIdentityScopedKeyLocks.lock(keys...)
	normalized := normalizeAuthPendingIdentityLockKeys(keys...)
	if len(normalized) == 0 || client == nil || client.Driver().Dialect() != dialect.Postgres {
		return release, nil
placeholder

	for _, key := range normalized {
		var rows entsql.Rows
		if err := client.Driver().Query(ctx, "SELECT pg_advisory_xact_lock($1)", []any{authPendingIdentityAdvisoryLockHash(key)placeholder, &rows); err != nil {
			release()
			return nil, err
	placeholder
		_ = rows.Close()
placeholder

	return release, nil
placeholder

func pendingIdentityAdoptionLockKeys(pendingAuthSessionID int64, identityID *int64) []string {
	keys := []string{fmt.Sprintf("pending-auth-adoption:pending:%d", pendingAuthSessionID)placeholder
	if identityID != nil && *identityID > 0 {
		keys = append(keys, fmt.Sprintf("pending-auth-adoption:identity:%d", *identityID))
placeholder
	return keys
placeholder

func NewAuthPendingIdentityService(entClient *dbent.Client) *AuthPendingIdentityService {
	return &AuthPendingIdentityService{entClient: entClientplaceholder
placeholder

func (s *AuthPendingIdentityService) CreatePendingSession(ctx context.Context, input CreatePendingAuthSessionInput) (*dbent.PendingAuthSession, error) {
	if s == nil || s.entClient == nil {
		return nil, fmt.Errorf("pending auth ent client is not configured")
placeholder

	sessionToken := strings.TrimSpace(input.SessionToken)
	if sessionToken == "" {
		var err error
		sessionToken, err = randomOpaqueToken(24)
		if err != nil {
			return nil, err
	placeholder
placeholder

	expiresAt := input.ExpiresAt.UTC()
	if expiresAt.IsZero() {
		expiresAt = time.Now().UTC().Add(defaultPendingAuthTTL)
placeholder

	create := s.entClient.PendingAuthSession.Create().
		SetSessionToken(sessionToken).
		SetIntent(strings.TrimSpace(input.Intent)).
		SetProviderType(strings.TrimSpace(input.Identity.ProviderType)).
		SetProviderKey(strings.TrimSpace(input.Identity.ProviderKey)).
		SetProviderSubject(strings.TrimSpace(input.Identity.ProviderSubject)).
		SetRedirectTo(strings.TrimSpace(input.RedirectTo)).
		SetResolvedEmail(strings.TrimSpace(input.ResolvedEmail)).
		SetRegistrationPasswordHash(strings.TrimSpace(input.RegistrationPasswordHash)).
		SetBrowserSessionKey(strings.TrimSpace(input.BrowserSessionKey)).
		SetUpstreamIdentityClaims(copyPendingMap(input.UpstreamIdentityClaims)).
		SetLocalFlowState(copyPendingMap(input.LocalFlowState)).
		SetExpiresAt(expiresAt)
	if input.TargetUserID != nil {
		create = create.SetTargetUserID(*input.TargetUserID)
placeholder
	return create.Save(ctx)
placeholder

func (s *AuthPendingIdentityService) IssueCompletionCode(ctx context.Context, input IssuePendingAuthCompletionCodeInput) (*IssuePendingAuthCompletionCodeResult, error) {
	if s == nil || s.entClient == nil {
		return nil, fmt.Errorf("pending auth ent client is not configured")
placeholder

	session, err := s.entClient.PendingAuthSession.Get(ctx, input.PendingAuthSessionID)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, ErrPendingAuthSessionNotFound
	placeholder
		return nil, err
placeholder

	code, err := randomOpaqueToken(24)
	if err != nil {
		return nil, err
placeholder
	ttl := input.TTL
	if ttl <= 0 {
		ttl = defaultPendingAuthCompletionTTL
placeholder
	expiresAt := time.Now().UTC().Add(ttl)

	update := s.entClient.PendingAuthSession.UpdateOneID(session.ID).
		SetCompletionCodeHash(hashPendingAuthCode(code)).
		SetCompletionCodeExpiresAt(expiresAt)
	if strings.TrimSpace(input.BrowserSessionKey) != "" {
		update = update.SetBrowserSessionKey(strings.TrimSpace(input.BrowserSessionKey))
placeholder
	if _, err := update.Save(ctx); err != nil {
		return nil, err
placeholder

	return &IssuePendingAuthCompletionCodeResult{
		Code:      code,
		ExpiresAt: expiresAt,
placeholder, nil
placeholder

func (s *AuthPendingIdentityService) ConsumeCompletionCode(ctx context.Context, rawCode, browserSessionKey string) (*dbent.PendingAuthSession, error) {
	if s == nil || s.entClient == nil {
		return nil, fmt.Errorf("pending auth ent client is not configured")
placeholder

	codeHash := hashPendingAuthCode(strings.TrimSpace(rawCode))
	session, err := s.entClient.PendingAuthSession.Query().
		Where(pendingauthsession.CompletionCodeHashEQ(codeHash)).
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, ErrPendingAuthCodeInvalid
	placeholder
		return nil, err
placeholder

	return s.consumeSession(ctx, session, browserSessionKey, ErrPendingAuthCodeExpired, ErrPendingAuthCodeConsumed)
placeholder

func (s *AuthPendingIdentityService) ConsumeBrowserSession(ctx context.Context, sessionToken, browserSessionKey string) (*dbent.PendingAuthSession, error) {
	if s == nil || s.entClient == nil {
		return nil, fmt.Errorf("pending auth ent client is not configured")
placeholder

	session, err := s.getBrowserSession(ctx, sessionToken)
	if err != nil {
		return nil, err
placeholder

	return s.consumeSession(ctx, session, browserSessionKey, ErrPendingAuthSessionExpired, ErrPendingAuthSessionConsumed)
placeholder

func (s *AuthPendingIdentityService) GetBrowserSession(ctx context.Context, sessionToken, browserSessionKey string) (*dbent.PendingAuthSession, error) {
	if s == nil || s.entClient == nil {
		return nil, fmt.Errorf("pending auth ent client is not configured")
placeholder

	session, err := s.getBrowserSession(ctx, sessionToken)
	if err != nil {
		return nil, err
placeholder
	if err := validatePendingSessionState(session, browserSessionKey, ErrPendingAuthSessionExpired, ErrPendingAuthSessionConsumed); err != nil {
		return nil, err
placeholder
	return session, nil
placeholder

func (s *AuthPendingIdentityService) getBrowserSession(ctx context.Context, sessionToken string) (*dbent.PendingAuthSession, error) {
	if s == nil || s.entClient == nil {
		return nil, fmt.Errorf("pending auth ent client is not configured")
placeholder

	sessionToken = strings.TrimSpace(sessionToken)
	if sessionToken == "" {
		return nil, ErrPendingAuthSessionNotFound
placeholder

	session, err := s.entClient.PendingAuthSession.Query().
		Where(pendingauthsession.SessionTokenEQ(sessionToken)).
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, ErrPendingAuthSessionNotFound
	placeholder
		return nil, err
placeholder
	return session, nil
placeholder

func (s *AuthPendingIdentityService) consumeSession(
	ctx context.Context,
	session *dbent.PendingAuthSession,
	browserSessionKey string,
	expiredErr error,
	consumedErr error,
) (*dbent.PendingAuthSession, error) {
	if err := validatePendingSessionState(session, browserSessionKey, expiredErr, consumedErr); err != nil {
		return nil, err
placeholder

	sanitizedLocalFlowState := sanitizePendingAuthLocalFlowState(session.LocalFlowState)
	now := time.Now().UTC()
	update := s.entClient.PendingAuthSession.UpdateOneID(session.ID).
		Where(
			pendingauthsession.ConsumedAtIsNil(),
			pendingauthsession.ExpiresAtGTE(now),
			pendingauthsession.Or(
				pendingauthsession.CompletionCodeExpiresAtIsNil(),
				pendingauthsession.CompletionCodeExpiresAtGTE(now),
			),
		).
		SetConsumedAt(now).
		SetLocalFlowState(sanitizedLocalFlowState).
		SetCompletionCodeHash("").
		ClearCompletionCodeExpiresAt()
	if expectedBrowserSessionKey := strings.TrimSpace(session.BrowserSessionKey); expectedBrowserSessionKey != "" {
		update = update.Where(pendingauthsession.BrowserSessionKeyEQ(expectedBrowserSessionKey))
placeholder
	updated, err := update.Save(ctx)
	if err == nil {
		return updated, nil
placeholder
	if !dbent.IsNotFound(err) {
		return nil, err
placeholder

	current, currentErr := s.entClient.PendingAuthSession.Get(ctx, session.ID)
	if currentErr != nil {
		if dbent.IsNotFound(currentErr) {
			return nil, ErrPendingAuthSessionNotFound
	placeholder
		return nil, currentErr
placeholder
	if err := validatePendingSessionState(current, browserSessionKey, expiredErr, consumedErr); err != nil {
		return nil, err
placeholder
	return nil, consumedErr
placeholder

func sanitizePendingAuthLocalFlowState(localFlowState map[string]any) map[string]any {
	sanitized := copyPendingMap(localFlowState)
	if len(sanitized) == 0 {
		return sanitized
placeholder

	rawCompletion, ok := sanitized["completion_response"]
	if !ok {
		return sanitized
placeholder
	completion, ok := rawCompletion.(map[string]any)
	if !ok {
		return sanitized
placeholder

	cleanedCompletion := copyPendingMap(completion)
	for _, key := range []string{"access_token", "refresh_token", "expires_in", "token_type"placeholder {
		delete(cleanedCompletion, key)
placeholder
	sanitized["completion_response"] = cleanedCompletion
	return sanitized
placeholder

func validatePendingSessionState(session *dbent.PendingAuthSession, browserSessionKey string, expiredErr error, consumedErr error) error {
	if session == nil {
		return ErrPendingAuthSessionNotFound
placeholder

	now := time.Now().UTC()
	if session.ConsumedAt != nil {
		return consumedErr
placeholder
	if !session.ExpiresAt.IsZero() && now.After(session.ExpiresAt) {
		return expiredErr
placeholder
	if session.CompletionCodeExpiresAt != nil && now.After(*session.CompletionCodeExpiresAt) {
		return expiredErr
placeholder
	if strings.TrimSpace(session.BrowserSessionKey) != "" && strings.TrimSpace(browserSessionKey) != strings.TrimSpace(session.BrowserSessionKey) {
		return ErrPendingAuthBrowserMismatch
placeholder
	return nil
placeholder

func (s *AuthPendingIdentityService) UpsertAdoptionDecision(ctx context.Context, input PendingIdentityAdoptionDecisionInput) (*dbent.IdentityAdoptionDecision, error) {
	if s == nil || s.entClient == nil {
		return nil, fmt.Errorf("pending auth ent client is not configured")
placeholder

	tx, err := s.entClient.Tx(ctx)
	if err != nil && !errors.Is(err, dbent.ErrTxStarted) {
		return nil, err
placeholder

	client := s.entClient
	txCtx := ctx
	if err == nil {
		defer func() { _ = tx.Rollback() placeholder()
		client = tx.Client()
		txCtx = dbent.NewTxContext(ctx, tx)
placeholder else if existingTx := dbent.TxFromContext(ctx); existingTx != nil {
		client = existingTx.Client()
placeholder

	releaseLocks, err := lockAuthPendingIdentityKeys(txCtx, client, pendingIdentityAdoptionLockKeys(input.PendingAuthSessionID, input.IdentityID)...)
	if err != nil {
		return nil, err
placeholder
	defer releaseLocks()

	if input.IdentityID != nil && *input.IdentityID > 0 {
		if _, err := client.IdentityAdoptionDecision.Update().
			Where(
				identityadoptiondecision.IdentityIDEQ(*input.IdentityID),
				dbpredicate.IdentityAdoptionDecision(func(s *entsql.Selector) {
					col := s.C(identityadoptiondecision.FieldPendingAuthSessionID)
					s.Where(entsql.Or(
						entsql.IsNull(col),
						entsql.NEQ(col, input.PendingAuthSessionID),
					))
			placeholder),
			).
			ClearIdentityID().
			Save(txCtx); err != nil {
			return nil, err
	placeholder
placeholder

	create := client.IdentityAdoptionDecision.Create().
		SetPendingAuthSessionID(input.PendingAuthSessionID).
		SetAdoptDisplayName(input.AdoptDisplayName).
		SetAdoptAvatar(input.AdoptAvatar).
		SetDecidedAt(time.Now().UTC())
	if input.IdentityID != nil && *input.IdentityID > 0 {
		create = create.SetIdentityID(*input.IdentityID)
placeholder

	decisionID, err := create.
		OnConflictColumns(identityadoptiondecision.FieldPendingAuthSessionID).
		UpdateNewValues().
		ID(txCtx)
	if err != nil {
		return nil, err
placeholder

	decision, err := client.IdentityAdoptionDecision.Get(txCtx, decisionID)
	if err != nil {
		return nil, err
placeholder

	if tx != nil {
		if err := tx.Commit(); err != nil {
			return nil, err
	placeholder
placeholder

	return decision, nil
placeholder

func copyPendingMap(in map[string]any) map[string]any {
	if len(in) == 0 {
		return map[string]any{placeholder
placeholder
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
placeholder
	return out
placeholder

func randomOpaqueToken(byteLen int) (string, error) {
	if byteLen <= 0 {
		byteLen = 16
placeholder
	buf := make([]byte, byteLen)
	if _, err := rand.Read(buf); err != nil {
		return "", err
placeholder
	return hex.EncodeToString(buf), nil
placeholder

func hashPendingAuthCode(code string) string {
	sum := sha256.Sum256([]byte(code))
	return hex.EncodeToString(sum[:])
placeholder
