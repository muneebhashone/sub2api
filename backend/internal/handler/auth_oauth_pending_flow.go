package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/authidentity"
	"github.com/Wei-Shaw/sub2api/ent/authidentitychannel"
	"github.com/Wei-Shaw/sub2api/ent/identityadoptiondecision"
	"github.com/Wei-Shaw/sub2api/ent/predicate"
	dbuser "github.com/Wei-Shaw/sub2api/ent/user"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/oauth"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"
)

const (
	oauthPendingBrowserCookiePath = "/api/v1/auth/oauth"
	oauthPendingBrowserCookieName = "oauth_pending_browser_session"
	oauthPendingSessionCookiePath = "/api/v1/auth/oauth/pending"
	oauthPendingSessionCookieName = "oauth_pending_session"
	oauthPendingCookieMaxAgeSec   = 10 * 60

	oauthCompletionResponseKey = "completion_response"
)

type oauthPendingSessionPayload struct {
	Intent                 string
	Identity               service.PendingAuthIdentityKey
	TargetUserID           *int64
	ResolvedEmail          string
	RedirectTo             string
	BrowserSessionKey      string
	UpstreamIdentityClaims map[string]any
	CompletionResponse     map[string]any
placeholder

type oauthAdoptionDecisionRequest struct {
	AdoptDisplayName *bool `json:"adopt_display_name,omitempty"`
	AdoptAvatar      *bool `json:"adopt_avatar,omitempty"`
placeholder

type bindPendingOAuthLoginRequest struct {
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required"`
	AdoptDisplayName *bool  `json:"adopt_display_name,omitempty"`
	AdoptAvatar      *bool  `json:"adopt_avatar,omitempty"`
placeholder

type createPendingOAuthAccountRequest struct {
	Email            string `json:"email" binding:"required,email"`
	VerifyCode       string `json:"verify_code,omitempty"`
	Password         string `json:"password" binding:"required,min=6"`
	InvitationCode   string `json:"invitation_code,omitempty"`
	AdoptDisplayName *bool  `json:"adopt_display_name,omitempty"`
	AdoptAvatar      *bool  `json:"adopt_avatar,omitempty"`
placeholder

func (r bindPendingOAuthLoginRequest) adoptionDecision() oauthAdoptionDecisionRequest {
	return oauthAdoptionDecisionRequest{
		AdoptDisplayName: r.AdoptDisplayName,
		AdoptAvatar:      r.AdoptAvatar,
placeholder
placeholder

func (r createPendingOAuthAccountRequest) adoptionDecision() oauthAdoptionDecisionRequest {
	return oauthAdoptionDecisionRequest{
		AdoptDisplayName: r.AdoptDisplayName,
		AdoptAvatar:      r.AdoptAvatar,
placeholder
placeholder

func (h *AuthHandler) pendingIdentityService() (*service.AuthPendingIdentityService, error) {
	if h == nil || h.authService == nil || h.authService.EntClient() == nil {
		return nil, infraerrors.ServiceUnavailable("PENDING_AUTH_NOT_READY", "pending auth service is not ready")
placeholder
	return service.NewAuthPendingIdentityService(h.authService.EntClient()), nil
placeholder

func generateOAuthPendingBrowserSession() (string, error) {
	return oauth.GenerateState()
placeholder

func setOAuthPendingBrowserCookie(c *gin.Context, sessionKey string, secure bool) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     oauthPendingBrowserCookieName,
		Value:    encodeCookieValue(sessionKey),
		Path:     oauthPendingBrowserCookiePath,
		MaxAge:   oauthPendingCookieMaxAgeSec,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
placeholder)
placeholder

func clearOAuthPendingBrowserCookie(c *gin.Context, secure bool) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     oauthPendingBrowserCookieName,
		Value:    "",
		Path:     oauthPendingBrowserCookiePath,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
placeholder)
placeholder

func readOAuthPendingBrowserCookie(c *gin.Context) (string, error) {
	return readCookieDecoded(c, oauthPendingBrowserCookieName)
placeholder

func setOAuthPendingSessionCookie(c *gin.Context, sessionToken string, secure bool) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     oauthPendingSessionCookieName,
		Value:    encodeCookieValue(sessionToken),
		Path:     oauthPendingSessionCookiePath,
		MaxAge:   oauthPendingCookieMaxAgeSec,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
placeholder)
placeholder

func clearOAuthPendingSessionCookie(c *gin.Context, secure bool) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     oauthPendingSessionCookieName,
		Value:    "",
		Path:     oauthPendingSessionCookiePath,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
placeholder)
placeholder

func readOAuthPendingSessionCookie(c *gin.Context) (string, error) {
	return readCookieDecoded(c, oauthPendingSessionCookieName)
placeholder

func redirectToFrontendCallback(c *gin.Context, frontendCallback string) {
	u, err := url.Parse(frontendCallback)
	if err != nil {
		c.Redirect(http.StatusFound, linuxDoOAuthDefaultRedirectTo)
		return
placeholder
	if u.Scheme != "" && !strings.EqualFold(u.Scheme, "http") && !strings.EqualFold(u.Scheme, "https") {
		c.Redirect(http.StatusFound, linuxDoOAuthDefaultRedirectTo)
		return
placeholder
	u.Fragment = ""
	c.Header("Cache-Control", "no-store")
	c.Header("Pragma", "no-cache")
	c.Redirect(http.StatusFound, u.String())
placeholder

func (h *AuthHandler) createOAuthPendingSession(c *gin.Context, payload oauthPendingSessionPayload) error {
	svc, err := h.pendingIdentityService()
	if err != nil {
		return err
placeholder

	session, err := svc.CreatePendingSession(c.Request.Context(), service.CreatePendingAuthSessionInput{
		Intent:                 strings.TrimSpace(payload.Intent),
		Identity:               payload.Identity,
		TargetUserID:           payload.TargetUserID,
		ResolvedEmail:          strings.TrimSpace(payload.ResolvedEmail),
		RedirectTo:             strings.TrimSpace(payload.RedirectTo),
		BrowserSessionKey:      strings.TrimSpace(payload.BrowserSessionKey),
		UpstreamIdentityClaims: payload.UpstreamIdentityClaims,
		LocalFlowState: map[string]any{
			oauthCompletionResponseKey: payload.CompletionResponse,
	placeholder,
placeholder)
	if err != nil {
		return infraerrors.InternalServer("PENDING_AUTH_SESSION_CREATE_FAILED", "failed to create pending auth session").WithCause(err)
placeholder

	setOAuthPendingSessionCookie(c, session.SessionToken, isRequestHTTPS(c))
	return nil
placeholder

func readCompletionResponse(session map[string]any) (map[string]any, bool) {
	if len(session) == 0 {
		return nil, false
placeholder
	value, ok := session[oauthCompletionResponseKey]
	if !ok {
		return nil, false
placeholder
	result, ok := value.(map[string]any)
	if !ok {
		return nil, false
placeholder
	return result, true
placeholder

func clonePendingMap(values map[string]any) map[string]any {
	if len(values) == 0 {
		return map[string]any{placeholder
placeholder
	cloned := make(map[string]any, len(values))
	for key, value := range values {
		cloned[key] = value
placeholder
	return cloned
placeholder

func mergePendingCompletionResponse(session *dbent.PendingAuthSession, overrides map[string]any) map[string]any {
	payload, _ := readCompletionResponse(session.LocalFlowState)
	merged := clonePendingMap(payload)
	if strings.TrimSpace(session.RedirectTo) != "" {
		if _, exists := merged["redirect"]; !exists {
			merged["redirect"] = session.RedirectTo
	placeholder
placeholder
	for key, value := range overrides {
		if value == nil {
			delete(merged, key)
			continue
	placeholder
		merged[key] = value
placeholder
	applySuggestedProfileToCompletionResponse(merged, session.UpstreamIdentityClaims)
	return merged
placeholder

func pendingSessionStringValue(values map[string]any, key string) string {
	if len(values) == 0 {
		return ""
placeholder
	raw, ok := values[key]
	if !ok {
		return ""
placeholder
	value, ok := raw.(string)
	if !ok {
		return ""
placeholder
	return strings.TrimSpace(value)
placeholder

func pendingSessionWantsInvitation(payload map[string]any) bool {
	return strings.EqualFold(strings.TrimSpace(pendingSessionStringValue(payload, "error")), "invitation_required")
placeholder

func (r oauthAdoptionDecisionRequest) hasDecision() bool {
	return r.AdoptDisplayName != nil || r.AdoptAvatar != nil
placeholder

func (r oauthAdoptionDecisionRequest) toServiceInput(sessionID int64) service.PendingIdentityAdoptionDecisionInput {
	input := service.PendingIdentityAdoptionDecisionInput{
		PendingAuthSessionID: sessionID,
placeholder
	if r.AdoptDisplayName != nil {
		input.AdoptDisplayName = *r.AdoptDisplayName
placeholder
	if r.AdoptAvatar != nil {
		input.AdoptAvatar = *r.AdoptAvatar
placeholder
	return input
placeholder

func bindOptionalOAuthAdoptionDecision(c *gin.Context) (oauthAdoptionDecisionRequest, error) {
	var req oauthAdoptionDecisionRequest
	if c == nil || c.Request == nil || c.Request.Body == nil {
		return req, nil
placeholder
	if err := c.ShouldBindJSON(&req); err != nil {
		if errors.Is(err, io.EOF) {
			return req, nil
	placeholder
		return req, err
placeholder
	return req, nil
placeholder

func persistPendingOAuthAdoptionDecision(
	c *gin.Context,
	svc *service.AuthPendingIdentityService,
	sessionID int64,
	req oauthAdoptionDecisionRequest,
) error {
	if !req.hasDecision() {
		return nil
placeholder
	if svc == nil {
		return infraerrors.ServiceUnavailable("PENDING_AUTH_NOT_READY", "pending auth service is not ready")
placeholder
	if _, err := svc.UpsertAdoptionDecision(c.Request.Context(), req.toServiceInput(sessionID)); err != nil {
		return infraerrors.InternalServer("PENDING_AUTH_ADOPTION_SAVE_FAILED", "failed to save oauth profile adoption decision").WithCause(err)
placeholder
	return nil
placeholder

func cloneOAuthMetadata(values map[string]any) map[string]any {
	if len(values) == 0 {
		return map[string]any{placeholder
placeholder
	cloned := make(map[string]any, len(values))
	for key, value := range values {
		cloned[key] = value
placeholder
	return cloned
placeholder

func mergeOAuthMetadata(base map[string]any, overlay map[string]any) map[string]any {
	merged := cloneOAuthMetadata(base)
	for key, value := range overlay {
		merged[key] = value
placeholder
	return merged
placeholder

func normalizeAdoptedOAuthDisplayName(value string) string {
	value = strings.TrimSpace(value)
	if len([]rune(value)) > 100 {
		value = string([]rune(value)[:100])
placeholder
	return value
placeholder

func (h *AuthHandler) entClient() *dbent.Client {
	if h == nil || h.authService == nil {
		return nil
placeholder
	return h.authService.EntClient()
placeholder

func (h *AuthHandler) isForceEmailOnThirdPartySignup(ctx context.Context) bool {
	if h == nil || h.settingSvc == nil {
		return false
placeholder
	defaults, err := h.settingSvc.GetAuthSourceDefaultSettings(ctx)
	if err != nil || defaults == nil {
		return false
placeholder
	return defaults.ForceEmailOnThirdPartySignup
placeholder

func (h *AuthHandler) findOAuthIdentityUser(ctx context.Context, identity service.PendingAuthIdentityKey) (*dbent.User, error) {
	client := h.entClient()
	if client == nil {
		return nil, infraerrors.ServiceUnavailable("PENDING_AUTH_NOT_READY", "pending auth service is not ready")
placeholder

	record, err := client.AuthIdentity.Query().
		Where(
			authidentity.ProviderTypeEQ(strings.TrimSpace(identity.ProviderType)),
			authidentity.ProviderKeyEQ(strings.TrimSpace(identity.ProviderKey)),
			authidentity.ProviderSubjectEQ(strings.TrimSpace(identity.ProviderSubject)),
		).
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, nil
	placeholder
		return nil, infraerrors.InternalServer("AUTH_IDENTITY_LOOKUP_FAILED", "failed to inspect auth identity ownership").WithCause(err)
placeholder

	userEntity, err := client.User.Get(ctx, record.UserID)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, nil
	placeholder
		return nil, infraerrors.InternalServer("AUTH_IDENTITY_USER_LOOKUP_FAILED", "failed to load auth identity user").WithCause(err)
placeholder
	return userEntity, nil
placeholder

func (h *AuthHandler) createOAuthEmailRequiredPendingSession(
	c *gin.Context,
	identity service.PendingAuthIdentityKey,
	redirectTo string,
	browserSessionKey string,
	upstreamClaims map[string]any,
) error {
	return h.createOAuthPendingSession(c, oauthPendingSessionPayload{
		Intent:                 oauthIntentLogin,
		Identity:               identity,
		RedirectTo:             redirectTo,
		BrowserSessionKey:      browserSessionKey,
		UpstreamIdentityClaims: upstreamClaims,
		CompletionResponse: map[string]any{
			"redirect":                  redirectTo,
			"step":                      "email_required",
			"force_email_on_signup":     true,
			"email_binding_required":    true,
			"existing_account_bindable": true,
	placeholder,
placeholder)
placeholder

func (h *AuthHandler) BindLinuxDoOAuthLogin(c *gin.Context) { h.bindPendingOAuthLogin(c, "linuxdo") placeholder
func (h *AuthHandler) BindOIDCOAuthLogin(c *gin.Context)    { h.bindPendingOAuthLogin(c, "oidc") placeholder
func (h *AuthHandler) BindWeChatOAuthLogin(c *gin.Context)  { h.bindPendingOAuthLogin(c, "wechat") placeholder
func (h *AuthHandler) BindPendingOAuthLogin(c *gin.Context) { h.bindPendingOAuthLogin(c, "") placeholder

func (h *AuthHandler) CreateLinuxDoOAuthAccount(c *gin.Context) {
	h.createPendingOAuthAccount(c, "linuxdo")
placeholder

func (h *AuthHandler) CreateOIDCOAuthAccount(c *gin.Context) { h.createPendingOAuthAccount(c, "oidc") placeholder

func (h *AuthHandler) CreateWeChatOAuthAccount(c *gin.Context) {
	h.createPendingOAuthAccount(c, "wechat")
placeholder

func (h *AuthHandler) CreatePendingOAuthAccount(c *gin.Context) {
	h.createPendingOAuthAccount(c, "")
placeholder

func (h *AuthHandler) upsertPendingOAuthAdoptionDecision(
	c *gin.Context,
	sessionID int64,
	req oauthAdoptionDecisionRequest,
) (*dbent.IdentityAdoptionDecision, error) {
	client := h.entClient()
	if client == nil {
		return nil, infraerrors.ServiceUnavailable("PENDING_AUTH_NOT_READY", "pending auth service is not ready")
placeholder

	existing, err := client.IdentityAdoptionDecision.Query().
		Where(identityadoptiondecision.PendingAuthSessionIDEQ(sessionID)).
		Only(c.Request.Context())
	if err != nil && !dbent.IsNotFound(err) {
		return nil, infraerrors.InternalServer("PENDING_AUTH_ADOPTION_LOAD_FAILED", "failed to load oauth profile adoption decision").WithCause(err)
placeholder
	if existing != nil && !req.hasDecision() {
		return existing, nil
placeholder
	if existing == nil && !req.hasDecision() {
		return nil, nil
placeholder

	input := service.PendingIdentityAdoptionDecisionInput{
		PendingAuthSessionID: sessionID,
placeholder
	if existing != nil {
		input.AdoptDisplayName = existing.AdoptDisplayName
		input.AdoptAvatar = existing.AdoptAvatar
		input.IdentityID = existing.IdentityID
placeholder
	if req.AdoptDisplayName != nil {
		input.AdoptDisplayName = *req.AdoptDisplayName
placeholder
	if req.AdoptAvatar != nil {
		input.AdoptAvatar = *req.AdoptAvatar
placeholder

	svc, err := h.pendingIdentityService()
	if err != nil {
		return nil, err
placeholder
	decision, err := svc.UpsertAdoptionDecision(c.Request.Context(), input)
	if err != nil {
		return nil, infraerrors.InternalServer("PENDING_AUTH_ADOPTION_SAVE_FAILED", "failed to save oauth profile adoption decision").WithCause(err)
placeholder
	return decision, nil
placeholder

func (h *AuthHandler) ensurePendingOAuthAdoptionDecision(
	c *gin.Context,
	sessionID int64,
	req oauthAdoptionDecisionRequest,
) (*dbent.IdentityAdoptionDecision, error) {
	decision, err := h.upsertPendingOAuthAdoptionDecision(c, sessionID, req)
	if err != nil {
		return nil, err
placeholder
	if decision != nil {
		return decision, nil
placeholder

	svc, err := h.pendingIdentityService()
	if err != nil {
		return nil, err
placeholder
	decision, err = svc.UpsertAdoptionDecision(c.Request.Context(), service.PendingIdentityAdoptionDecisionInput{
		PendingAuthSessionID: sessionID,
placeholder)
	if err != nil {
		return nil, infraerrors.InternalServer("PENDING_AUTH_ADOPTION_SAVE_FAILED", "failed to save oauth profile adoption decision").WithCause(err)
placeholder
	return decision, nil
placeholder

func updatePendingOAuthSessionProgress(
	ctx context.Context,
	client *dbent.Client,
	session *dbent.PendingAuthSession,
	intent string,
	resolvedEmail string,
	targetUserID *int64,
	completionResponse map[string]any,
) (*dbent.PendingAuthSession, error) {
	if client == nil || session == nil {
		return nil, infraerrors.BadRequest("PENDING_AUTH_SESSION_INVALID", "pending auth session is invalid")
placeholder

	localFlowState := clonePendingMap(session.LocalFlowState)
	localFlowState[oauthCompletionResponseKey] = clonePendingMap(completionResponse)

	update := client.PendingAuthSession.UpdateOneID(session.ID).
		SetIntent(strings.TrimSpace(intent)).
		SetResolvedEmail(strings.TrimSpace(resolvedEmail)).
		SetLocalFlowState(localFlowState)
	if targetUserID != nil && *targetUserID > 0 {
		update = update.SetTargetUserID(*targetUserID)
placeholder else {
		update = update.ClearTargetUserID()
placeholder
	return update.Save(ctx)
placeholder

func resolvePendingOAuthTargetUserID(ctx context.Context, client *dbent.Client, session *dbent.PendingAuthSession) (int64, error) {
	if session == nil {
		return 0, infraerrors.BadRequest("PENDING_AUTH_SESSION_INVALID", "pending auth session is invalid")
placeholder
	if session.TargetUserID != nil && *session.TargetUserID > 0 {
		return *session.TargetUserID, nil
placeholder
	email := strings.TrimSpace(session.ResolvedEmail)
	if email == "" {
		return 0, infraerrors.BadRequest("PENDING_AUTH_TARGET_USER_MISSING", "pending auth target user is missing")
placeholder

	userEntity, err := findUserByNormalizedEmail(ctx, client, email)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return 0, infraerrors.InternalServer("PENDING_AUTH_TARGET_USER_NOT_FOUND", "pending auth target user was not found")
	placeholder
		return 0, err
placeholder
	return userEntity.ID, nil
placeholder

func userNormalizedEmailPredicate(email string) predicate.User {
	normalized := strings.TrimSpace(email)
	if normalized == "" {
		return dbuser.EmailEQ(email)
placeholder
	return predicate.User(func(s *entsql.Selector) {
		s.Where(entsql.ExprP(
			fmt.Sprintf("LOWER(TRIM(%s)) = LOWER(TRIM(?))", s.C(dbuser.FieldEmail)),
			normalized,
		))
placeholder)
placeholder

func findUserByNormalizedEmail(ctx context.Context, client *dbent.Client, email string) (*dbent.User, error) {
	if client == nil {
		return nil, infraerrors.ServiceUnavailable("PENDING_AUTH_NOT_READY", "pending auth service is not ready")
placeholder

	matches, err := client.User.Query().
		Where(userNormalizedEmailPredicate(email)).
		Order(dbent.Asc(dbuser.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
placeholder
	if len(matches) == 0 {
		return nil, service.ErrUserNotFound
placeholder
	if len(matches) > 1 {
		return nil, infraerrors.Conflict("USER_EMAIL_CONFLICT", "normalized email matched multiple users")
placeholder
	return matches[0], nil
placeholder

func oauthIdentityIssuer(session *dbent.PendingAuthSession) *string {
	if session == nil {
		return nil
placeholder
	switch strings.TrimSpace(session.ProviderType) {
	case "oidc":
		issuer := strings.TrimSpace(session.ProviderKey)
		if issuer == "" {
			issuer = pendingSessionStringValue(session.UpstreamIdentityClaims, "issuer")
	placeholder
		if issuer == "" {
			return nil
	placeholder
		return &issuer
	default:
		issuer := pendingSessionStringValue(session.UpstreamIdentityClaims, "issuer")
		if issuer == "" {
			return nil
	placeholder
		return &issuer
placeholder
placeholder

func ensurePendingOAuthIdentityForUser(ctx context.Context, tx *dbent.Tx, session *dbent.PendingAuthSession, userID int64) (*dbent.AuthIdentity, error) {
	if session != nil && strings.EqualFold(strings.TrimSpace(session.ProviderType), "wechat") {
		return ensurePendingWeChatOAuthIdentityForUser(ctx, tx, session, userID)
placeholder

	client := tx.Client()
	identity, err := client.AuthIdentity.Query().
		Where(
			authidentity.ProviderTypeEQ(strings.TrimSpace(session.ProviderType)),
			authidentity.ProviderKeyEQ(strings.TrimSpace(session.ProviderKey)),
			authidentity.ProviderSubjectEQ(strings.TrimSpace(session.ProviderSubject)),
		).
		Only(ctx)
	if err != nil && !dbent.IsNotFound(err) {
		return nil, err
placeholder
	if identity != nil {
		if identity.UserID != userID {
			return nil, infraerrors.Conflict("AUTH_IDENTITY_OWNERSHIP_CONFLICT", "auth identity already belongs to another user")
	placeholder
		return identity, nil
placeholder

	create := client.AuthIdentity.Create().
		SetUserID(userID).
		SetProviderType(strings.TrimSpace(session.ProviderType)).
		SetProviderKey(strings.TrimSpace(session.ProviderKey)).
		SetProviderSubject(strings.TrimSpace(session.ProviderSubject)).
		SetMetadata(cloneOAuthMetadata(session.UpstreamIdentityClaims))
	if issuer := oauthIdentityIssuer(session); issuer != nil {
		create = create.SetIssuer(strings.TrimSpace(*issuer))
placeholder
	return create.Save(ctx)
placeholder

func ensurePendingWeChatOAuthIdentityForUser(ctx context.Context, tx *dbent.Tx, session *dbent.PendingAuthSession, userID int64) (*dbent.AuthIdentity, error) {
	client := tx.Client()
	providerType := strings.TrimSpace(session.ProviderType)
	providerKey := strings.TrimSpace(session.ProviderKey)
	providerSubject := strings.TrimSpace(session.ProviderSubject)
	providerKeys := wechatCompatibleProviderKeys(providerKey)
	channel := strings.TrimSpace(pendingSessionStringValue(session.UpstreamIdentityClaims, "channel"))
	channelAppID := strings.TrimSpace(pendingSessionStringValue(session.UpstreamIdentityClaims, "channel_app_id"))
	channelSubject := strings.TrimSpace(pendingSessionStringValue(session.UpstreamIdentityClaims, "channel_subject"))
	metadata := cloneOAuthMetadata(session.UpstreamIdentityClaims)

	identityRecords, err := client.AuthIdentity.Query().
		Where(
			authidentity.ProviderTypeEQ(providerType),
			authidentity.ProviderKeyIn(providerKeys...),
			authidentity.ProviderSubjectEQ(providerSubject),
		).
		All(ctx)
	if err != nil {
		return nil, err
placeholder
	identity, hasCanonicalKey, err := chooseWeChatIdentityForUser(identityRecords, userID, providerKey)
	if err != nil {
		return nil, err
placeholder

	var legacyOpenIDIdentity *dbent.AuthIdentity
	if channelSubject != "" && channelSubject != providerSubject {
		legacyOpenIDRecords, err := client.AuthIdentity.Query().
			Where(
				authidentity.ProviderTypeEQ(providerType),
				authidentity.ProviderKeyIn(providerKeys...),
				authidentity.ProviderSubjectEQ(channelSubject),
			).
			All(ctx)
		if err != nil {
			return nil, err
	placeholder
		legacyOpenIDIdentity, _, err = chooseWeChatIdentityForUser(legacyOpenIDRecords, userID, providerKey)
		if err != nil {
			return nil, err
	placeholder
placeholder

	switch {
	case identity != nil:
		update := client.AuthIdentity.UpdateOneID(identity.ID).
			SetMetadata(mergeOAuthMetadata(identity.Metadata, metadata))
		if !strings.EqualFold(strings.TrimSpace(identity.ProviderKey), providerKey) && !hasCanonicalKey {
			update = update.SetProviderKey(providerKey)
	placeholder
		if issuer := oauthIdentityIssuer(session); issuer != nil {
			update = update.SetIssuer(strings.TrimSpace(*issuer))
	placeholder
		identity, err = update.Save(ctx)
		if err != nil {
			return nil, err
	placeholder
	case legacyOpenIDIdentity != nil:
		update := client.AuthIdentity.UpdateOneID(legacyOpenIDIdentity.ID).
			SetProviderKey(providerKey).
			SetProviderSubject(providerSubject).
			SetMetadata(mergeOAuthMetadata(legacyOpenIDIdentity.Metadata, metadata))
		if issuer := oauthIdentityIssuer(session); issuer != nil {
			update = update.SetIssuer(strings.TrimSpace(*issuer))
	placeholder
		identity, err = update.Save(ctx)
		if err != nil {
			return nil, err
	placeholder
	default:
		create := client.AuthIdentity.Create().
			SetUserID(userID).
			SetProviderType(providerType).
			SetProviderKey(providerKey).
			SetProviderSubject(providerSubject).
			SetMetadata(metadata)
		if issuer := oauthIdentityIssuer(session); issuer != nil {
			create = create.SetIssuer(strings.TrimSpace(*issuer))
	placeholder
		identity, err = create.Save(ctx)
		if err != nil {
			return nil, err
	placeholder
placeholder

	if channel == "" || channelAppID == "" || channelSubject == "" {
		return identity, nil
placeholder

	channelRecords, err := client.AuthIdentityChannel.Query().
		Where(
			authidentitychannel.ProviderTypeEQ(providerType),
			authidentitychannel.ProviderKeyIn(providerKeys...),
			authidentitychannel.ChannelEQ(channel),
			authidentitychannel.ChannelAppIDEQ(channelAppID),
			authidentitychannel.ChannelSubjectEQ(channelSubject),
		).
		WithIdentity().
		All(ctx)
	if err != nil {
		return nil, err
placeholder
	channelRecord, hasCanonicalChannelKey, err := chooseWeChatChannelForUser(channelRecords, userID, providerKey)
	if err != nil {
		return nil, err
placeholder

	channelMetadata := mergeOAuthMetadata(channelRecordMetadata(channelRecord), metadata)
	if channelRecord == nil {
		if _, err := client.AuthIdentityChannel.Create().
			SetIdentityID(identity.ID).
			SetProviderType(providerType).
			SetProviderKey(providerKey).
			SetChannel(channel).
			SetChannelAppID(channelAppID).
			SetChannelSubject(channelSubject).
			SetMetadata(channelMetadata).
			Save(ctx); err != nil {
			return nil, err
	placeholder
		return identity, nil
placeholder

	updateChannel := client.AuthIdentityChannel.UpdateOneID(channelRecord.ID).
		SetIdentityID(identity.ID).
		SetMetadata(channelMetadata)
	if !strings.EqualFold(strings.TrimSpace(channelRecord.ProviderKey), providerKey) && !hasCanonicalChannelKey {
		updateChannel = updateChannel.SetProviderKey(providerKey)
placeholder
	_, err = updateChannel.Save(ctx)
	if err != nil {
		return nil, err
placeholder
	return identity, nil
placeholder

func chooseWeChatIdentityForUser(records []*dbent.AuthIdentity, userID int64, preferredProviderKey string) (*dbent.AuthIdentity, bool, error) {
	var preferred *dbent.AuthIdentity
	var fallback *dbent.AuthIdentity
	hasCanonicalKey := false
	for _, record := range records {
		if record == nil {
			continue
	placeholder
		if record.UserID != userID {
			return nil, false, infraerrors.Conflict("AUTH_IDENTITY_OWNERSHIP_CONFLICT", "auth identity already belongs to another user")
	placeholder
		if strings.EqualFold(strings.TrimSpace(record.ProviderKey), preferredProviderKey) {
			hasCanonicalKey = true
			if preferred == nil {
				preferred = record
		placeholder
			continue
	placeholder
		if fallback == nil {
			fallback = record
	placeholder
placeholder
	if preferred != nil {
		return preferred, hasCanonicalKey, nil
placeholder
	return fallback, hasCanonicalKey, nil
placeholder

func chooseWeChatChannelForUser(records []*dbent.AuthIdentityChannel, userID int64, preferredProviderKey string) (*dbent.AuthIdentityChannel, bool, error) {
	var preferred *dbent.AuthIdentityChannel
	var fallback *dbent.AuthIdentityChannel
	hasCanonicalKey := false
	for _, record := range records {
		if record == nil {
			continue
	placeholder
		if record.Edges.Identity != nil && record.Edges.Identity.UserID != userID {
			return nil, false, infraerrors.Conflict("AUTH_IDENTITY_CHANNEL_OWNERSHIP_CONFLICT", "auth identity channel already belongs to another user")
	placeholder
		if strings.EqualFold(strings.TrimSpace(record.ProviderKey), preferredProviderKey) {
			hasCanonicalKey = true
			if preferred == nil {
				preferred = record
		placeholder
			continue
	placeholder
		if fallback == nil {
			fallback = record
	placeholder
placeholder
	if preferred != nil {
		return preferred, hasCanonicalKey, nil
placeholder
	return fallback, hasCanonicalKey, nil
placeholder

func channelRecordMetadata(channel *dbent.AuthIdentityChannel) map[string]any {
	if channel == nil {
		return map[string]any{placeholder
placeholder
	return cloneOAuthMetadata(channel.Metadata)
placeholder

func shouldBindPendingOAuthIdentity(session *dbent.PendingAuthSession, decision *dbent.IdentityAdoptionDecision) bool {
	if session == nil || decision == nil {
		return false
placeholder
	switch strings.ToLower(strings.TrimSpace(session.Intent)) {
	case "bind_current_user", "login", "adopt_existing_user_by_email":
		return true
	default:
		return decision.AdoptDisplayName || decision.AdoptAvatar
placeholder
placeholder

func applyPendingOAuthBinding(
	ctx context.Context,
	client *dbent.Client,
	authService *service.AuthService,
	session *dbent.PendingAuthSession,
	decision *dbent.IdentityAdoptionDecision,
	overrideUserID *int64,
	forceBind bool,
	applyFirstBindDefaults bool,
) error {
	if client == nil || session == nil {
		return nil
placeholder
	if !forceBind && !shouldBindPendingOAuthIdentity(session, decision) {
		return nil
placeholder

	targetUserID := int64(0)
	if overrideUserID != nil && *overrideUserID > 0 {
		targetUserID = *overrideUserID
placeholder else {
		resolvedUserID, err := resolvePendingOAuthTargetUserID(ctx, client, session)
		if err != nil {
			return err
	placeholder
		targetUserID = resolvedUserID
placeholder

	adoptedDisplayName := ""
	if decision != nil && decision.AdoptDisplayName {
		adoptedDisplayName = normalizeAdoptedOAuthDisplayName(pendingSessionStringValue(session.UpstreamIdentityClaims, "suggested_display_name"))
placeholder
	adoptedAvatarURL := ""
	if decision != nil && decision.AdoptAvatar {
		adoptedAvatarURL = pendingSessionStringValue(session.UpstreamIdentityClaims, "suggested_avatar_url")
placeholder

	tx, err := client.Tx(ctx)
	if err != nil {
		return err
placeholder
	defer func() { _ = tx.Rollback() placeholder()
	txCtx := dbent.NewTxContext(ctx, tx)

	if decision != nil && decision.AdoptDisplayName && adoptedDisplayName != "" {
		if err := tx.Client().User.UpdateOneID(targetUserID).
			SetUsername(adoptedDisplayName).
			Exec(txCtx); err != nil {
			return err
	placeholder
placeholder

	identity, err := ensurePendingOAuthIdentityForUser(txCtx, tx, session, targetUserID)
	if err != nil {
		return err
placeholder

	metadata := cloneOAuthMetadata(identity.Metadata)
	for key, value := range session.UpstreamIdentityClaims {
		metadata[key] = value
placeholder
	if decision != nil && decision.AdoptDisplayName && adoptedDisplayName != "" {
		metadata["display_name"] = adoptedDisplayName
placeholder
	if decision != nil && decision.AdoptAvatar && adoptedAvatarURL != "" {
		metadata["avatar_url"] = adoptedAvatarURL
placeholder

	updateIdentity := tx.Client().AuthIdentity.UpdateOneID(identity.ID).SetMetadata(metadata)
	if issuer := oauthIdentityIssuer(session); issuer != nil {
		updateIdentity = updateIdentity.SetIssuer(strings.TrimSpace(*issuer))
placeholder
	if _, err := updateIdentity.Save(txCtx); err != nil {
		return err
placeholder

	if decision != nil && (decision.IdentityID == nil || *decision.IdentityID != identity.ID) {
		if _, err := tx.Client().IdentityAdoptionDecision.UpdateOneID(decision.ID).
			SetIdentityID(identity.ID).
			Save(txCtx); err != nil {
			return err
	placeholder
placeholder

	if applyFirstBindDefaults && authService != nil {
		if err := authService.ApplyProviderDefaultSettingsOnFirstBind(txCtx, targetUserID, session.ProviderType); err != nil {
			return err
	placeholder
placeholder

	return tx.Commit()
placeholder

func applyPendingOAuthAdoption(
	ctx context.Context,
	client *dbent.Client,
	authService *service.AuthService,
	session *dbent.PendingAuthSession,
	decision *dbent.IdentityAdoptionDecision,
	overrideUserID *int64,
) error {
	return applyPendingOAuthBinding(
		ctx,
		client,
		authService,
		session,
		decision,
		overrideUserID,
		false,
		strings.EqualFold(strings.TrimSpace(session.Intent), "bind_current_user"),
	)
placeholder

func applySuggestedProfileToCompletionResponse(payload map[string]any, upstream map[string]any) {
	if len(payload) == 0 || len(upstream) == 0 {
		return
placeholder

	displayName := pendingSessionStringValue(upstream, "suggested_display_name")
	avatarURL := pendingSessionStringValue(upstream, "suggested_avatar_url")

	if displayName != "" {
		if _, exists := payload["suggested_display_name"]; !exists {
			payload["suggested_display_name"] = displayName
	placeholder
placeholder
	if avatarURL != "" {
		if _, exists := payload["suggested_avatar_url"]; !exists {
			payload["suggested_avatar_url"] = avatarURL
	placeholder
placeholder
	if displayName != "" || avatarURL != "" {
		payload["adoption_required"] = true
placeholder
placeholder

func readPendingOAuthBrowserSession(c *gin.Context, h *AuthHandler) (*service.AuthPendingIdentityService, *dbent.PendingAuthSession, func(), error) {
	secureCookie := isRequestHTTPS(c)
	clearCookies := func() {
		clearOAuthPendingSessionCookie(c, secureCookie)
		clearOAuthPendingBrowserCookie(c, secureCookie)
placeholder

	sessionToken, err := readOAuthPendingSessionCookie(c)
	if err != nil || strings.TrimSpace(sessionToken) == "" {
		clearCookies()
		return nil, nil, clearCookies, service.ErrPendingAuthSessionNotFound
placeholder
	browserSessionKey, err := readOAuthPendingBrowserCookie(c)
	if err != nil || strings.TrimSpace(browserSessionKey) == "" {
		clearCookies()
		return nil, nil, clearCookies, service.ErrPendingAuthBrowserMismatch
placeholder

	svc, err := h.pendingIdentityService()
	if err != nil {
		clearCookies()
		return nil, nil, clearCookies, err
placeholder

	session, err := svc.GetBrowserSession(c.Request.Context(), sessionToken, browserSessionKey)
	if err != nil {
		clearCookies()
		return nil, nil, clearCookies, err
placeholder

	return svc, session, clearCookies, nil
placeholder

func buildPendingOAuthSessionStatusPayload(session *dbent.PendingAuthSession) gin.H {
	payload := gin.H{
		"auth_result": "pending_session",
		"provider":    strings.TrimSpace(session.ProviderType),
		"intent":      strings.TrimSpace(session.Intent),
placeholder
	for key, value := range mergePendingCompletionResponse(session, nil) {
		payload[key] = value
placeholder
	if email := strings.TrimSpace(session.ResolvedEmail); email != "" {
		payload["email"] = email
placeholder
	return payload
placeholder

func writeOAuthTokenPairResponse(c *gin.Context, tokenPair *service.TokenPair) {
	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"expires_in":    tokenPair.ExpiresIn,
		"token_type":    "Bearer",
placeholder)
placeholder

func (h *AuthHandler) bindPendingOAuthLogin(c *gin.Context, provider string) {
	var req bindPendingOAuthLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	pendingSvc, session, clearCookies, err := readPendingOAuthBrowserSession(c, h)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	if strings.TrimSpace(provider) != "" && !strings.EqualFold(strings.TrimSpace(session.ProviderType), provider) {
		response.BadRequest(c, "Pending oauth session provider mismatch")
		return
placeholder

	user, err := h.authService.ValidatePasswordCredentials(c.Request.Context(), strings.TrimSpace(req.Email), req.Password)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	if session.TargetUserID != nil && *session.TargetUserID > 0 && user.ID != *session.TargetUserID {
		response.ErrorFrom(c, infraerrors.Conflict("PENDING_AUTH_TARGET_USER_MISMATCH", "pending oauth session must be completed by the targeted user"))
		return
placeholder

	decision, err := h.ensurePendingOAuthAdoptionDecision(c, session.ID, req.adoptionDecision())
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	if h.totpService != nil && h.settingSvc.IsTotpEnabled(c.Request.Context()) && user.TotpEnabled {
		tempToken, err := h.totpService.CreatePendingOAuthBindLoginSession(
			c.Request.Context(),
			user.ID,
			user.Email,
			session.SessionToken,
			session.BrowserSessionKey,
		)
		if err != nil {
			response.InternalError(c, "Failed to create 2FA session")
			return
	placeholder
		response.Success(c, TotpLoginResponse{
			Requires2FA:     true,
			TempToken:       tempToken,
			UserEmailMasked: service.MaskEmail(user.Email),
	placeholder)
		return
placeholder
	if err := applyPendingOAuthBinding(c.Request.Context(), h.entClient(), h.authService, session, decision, &user.ID, true, true); err != nil {
		response.ErrorFrom(c, infraerrors.InternalServer("PENDING_AUTH_BIND_APPLY_FAILED", "failed to bind pending oauth identity").WithCause(err))
		return
placeholder

	h.authService.RecordSuccessfulLogin(c.Request.Context(), user.ID)
	tokenPair, err := h.authService.GenerateTokenPair(c.Request.Context(), user, "")
	if err != nil {
		response.InternalError(c, "Failed to generate token pair")
		return
placeholder
	if _, err := pendingSvc.ConsumeBrowserSession(c.Request.Context(), session.SessionToken, session.BrowserSessionKey); err != nil {
		clearCookies()
		response.ErrorFrom(c, err)
		return
placeholder

	clearCookies()
	writeOAuthTokenPairResponse(c, tokenPair)
placeholder

func (h *AuthHandler) createPendingOAuthAccount(c *gin.Context, provider string) {
	var req createPendingOAuthAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	pendingSvc, session, clearCookies, err := readPendingOAuthBrowserSession(c, h)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	if strings.TrimSpace(provider) != "" && !strings.EqualFold(strings.TrimSpace(session.ProviderType), provider) {
		response.BadRequest(c, "Pending oauth session provider mismatch")
		return
placeholder

	client := h.entClient()
	if client == nil {
		response.ErrorFrom(c, infraerrors.ServiceUnavailable("PENDING_AUTH_NOT_READY", "pending auth service is not ready"))
		return
placeholder

	email := strings.TrimSpace(strings.ToLower(req.Email))
	existingUser, err := findUserByNormalizedEmail(c.Request.Context(), client, email)
	if err != nil && !errors.Is(err, service.ErrUserNotFound) {
		response.ErrorFrom(c, infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "service temporarily unavailable"))
		return
placeholder
	if existingUser != nil {
		completionResponse := mergePendingCompletionResponse(session, map[string]any{
			"step":  "bind_login_required",
			"email": email,
	placeholder)
		session, err = updatePendingOAuthSessionProgress(
			c.Request.Context(),
			client,
			session,
			"adopt_existing_user_by_email",
			email,
			&existingUser.ID,
			completionResponse,
		)
		if err != nil {
			response.ErrorFrom(c, infraerrors.InternalServer("PENDING_AUTH_SESSION_UPDATE_FAILED", "failed to update pending oauth session").WithCause(err))
			return
	placeholder

		if _, err := h.ensurePendingOAuthAdoptionDecision(c, session.ID, req.adoptionDecision()); err != nil {
			response.ErrorFrom(c, err)
			return
	placeholder

		c.JSON(http.StatusOK, buildPendingOAuthSessionStatusPayload(session))
		return
placeholder

	tokenPair, user, err := h.authService.RegisterOAuthEmailAccount(
		c.Request.Context(),
		email,
		req.Password,
		strings.TrimSpace(req.VerifyCode),
		strings.TrimSpace(req.InvitationCode),
		strings.TrimSpace(session.ProviderType),
	)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	decision, err := h.ensurePendingOAuthAdoptionDecision(c, session.ID, req.adoptionDecision())
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	if err := applyPendingOAuthBinding(c.Request.Context(), client, h.authService, session, decision, &user.ID, true, false); err != nil {
		response.ErrorFrom(c, infraerrors.InternalServer("PENDING_AUTH_BIND_APPLY_FAILED", "failed to bind pending oauth identity").WithCause(err))
		return
placeholder

	if _, err := pendingSvc.ConsumeBrowserSession(c.Request.Context(), session.SessionToken, session.BrowserSessionKey); err != nil {
		clearCookies()
		response.ErrorFrom(c, err)
		return
placeholder

	clearCookies()
	writeOAuthTokenPairResponse(c, tokenPair)
placeholder

// ExchangePendingOAuthCompletion redeems a pending OAuth browser session into a frontend-safe payload.
// POST /api/v1/auth/oauth/pending/exchange
func (h *AuthHandler) ExchangePendingOAuthCompletion(c *gin.Context) {
	secureCookie := isRequestHTTPS(c)
	clearCookies := func() {
		clearOAuthPendingSessionCookie(c, secureCookie)
		clearOAuthPendingBrowserCookie(c, secureCookie)
placeholder
	adoptionDecision, err := bindOptionalOAuthAdoptionDecision(c)
	if err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	sessionToken, err := readOAuthPendingSessionCookie(c)
	if err != nil || strings.TrimSpace(sessionToken) == "" {
		clearCookies()
		response.ErrorFrom(c, service.ErrPendingAuthSessionNotFound)
		return
placeholder
	browserSessionKey, err := readOAuthPendingBrowserCookie(c)
	if err != nil || strings.TrimSpace(browserSessionKey) == "" {
		clearCookies()
		response.ErrorFrom(c, service.ErrPendingAuthBrowserMismatch)
		return
placeholder

	svc, err := h.pendingIdentityService()
	if err != nil {
		clearCookies()
		response.ErrorFrom(c, err)
		return
placeholder

	session, err := svc.GetBrowserSession(c.Request.Context(), sessionToken, browserSessionKey)
	if err != nil {
		clearCookies()
		response.ErrorFrom(c, err)
		return
placeholder

	payload, ok := readCompletionResponse(session.LocalFlowState)
	if !ok {
		clearCookies()
		response.ErrorFrom(c, infraerrors.InternalServer("PENDING_AUTH_COMPLETION_INVALID", "pending auth completion payload is invalid"))
		return
placeholder
	if strings.TrimSpace(session.RedirectTo) != "" {
		if _, exists := payload["redirect"]; !exists {
			payload["redirect"] = session.RedirectTo
	placeholder
placeholder
	applySuggestedProfileToCompletionResponse(payload, session.UpstreamIdentityClaims)

	if pendingSessionWantsInvitation(payload) {
		if adoptionDecision.hasDecision() {
			decision, err := h.upsertPendingOAuthAdoptionDecision(c, session.ID, adoptionDecision)
			if err != nil {
				response.ErrorFrom(c, err)
				return
		placeholder
			_ = decision
	placeholder
		response.Success(c, payload)
		return
placeholder
	if !adoptionDecision.hasDecision() {
		response.Success(c, payload)
		return
placeholder
	decision, err := h.upsertPendingOAuthAdoptionDecision(c, session.ID, adoptionDecision)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	if err := applyPendingOAuthAdoption(c.Request.Context(), h.entClient(), h.authService, session, decision, session.TargetUserID); err != nil {
		response.ErrorFrom(c, infraerrors.InternalServer("PENDING_AUTH_ADOPTION_APPLY_FAILED", "failed to apply oauth profile adoption").WithCause(err))
		return
placeholder

	if _, err := svc.ConsumeBrowserSession(c.Request.Context(), sessionToken, browserSessionKey); err != nil {
		clearCookies()
		response.ErrorFrom(c, err)
		return
placeholder

	clearCookies()
	response.Success(c, payload)
placeholder
