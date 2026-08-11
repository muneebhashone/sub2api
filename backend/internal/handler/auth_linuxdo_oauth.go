package handler

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	dbuser "github.com/Wei-Shaw/sub2api/ent/user"
	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/oauth"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	servermiddleware "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
)

const (
	linuxDoOAuthCookiePath         = "/api/v1/auth/oauth/linuxdo"
	oauthBindAccessTokenCookiePath = "/api/v1/auth/oauth"
	linuxDoOAuthStateCookieName    = "linuxdo_oauth_state"
	linuxDoOAuthVerifierCookie     = "linuxdo_oauth_verifier"
	linuxDoOAuthRedirectCookie     = "linuxdo_oauth_redirect"
	linuxDoOAuthIntentCookieName   = "linuxdo_oauth_intent"
	linuxDoOAuthBindUserCookieName = "linuxdo_oauth_bind_user"
	oauthBindAccessTokenCookieName = "oauth_bind_access_token"
	linuxDoOAuthCookieMaxAgeSec    = 10 * 60 // 10 minutes
	linuxDoOAuthDefaultRedirectTo  = "/dashboard"
	linuxDoOAuthDefaultFrontendCB  = "/auth/linuxdo/callback"

	linuxDoOAuthMaxRedirectLen      = 2048
	linuxDoOAuthMaxFragmentValueLen = 512
	linuxDoOAuthMaxSubjectLen       = 64 - len("linuxdo-")

	oauthIntentLogin           = "login"
	oauthIntentBindCurrentUser = "bind_current_user"
)

type linuxDoTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
placeholder

type linuxDoTokenExchangeError struct {
	StatusCode          int
	ProviderError       string
	ProviderDescription string
	Body                string
placeholder

func (e *linuxDoTokenExchangeError) Error() string {
	if e == nil {
		return ""
placeholder
	parts := []string{fmt.Sprintf("token exchange status=%d", e.StatusCode)placeholder
	if strings.TrimSpace(e.ProviderError) != "" {
		parts = append(parts, "error="+strings.TrimSpace(e.ProviderError))
placeholder
	if strings.TrimSpace(e.ProviderDescription) != "" {
		parts = append(parts, "error_description="+strings.TrimSpace(e.ProviderDescription))
placeholder
	return strings.Join(parts, " ")
placeholder

// LinuxDoOAuthStart 启动 LinuxDo Connect OAuth 登录流程。
// GET /api/v1/auth/oauth/linuxdo/start?redirect=/dashboard
func (h *AuthHandler) LinuxDoOAuthStart(c *gin.Context) {
	if !h.requireActionCaptchaForOAuthLoginStart(c) {
		return
placeholder
	cfg, err := h.getLinuxDoOAuthConfig(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	state, err := oauth.GenerateState()
	if err != nil {
		response.ErrorFrom(c, infraerrors.InternalServer("OAUTH_STATE_GEN_FAILED", "failed to generate oauth state").WithCause(err))
		return
placeholder

	redirectTo := sanitizeFrontendRedirectPath(c.Query("redirect"))
	if redirectTo == "" {
		redirectTo = linuxDoOAuthDefaultRedirectTo
placeholder

	browserSessionKey, err := generateOAuthPendingBrowserSession()
	if err != nil {
		response.ErrorFrom(c, infraerrors.InternalServer("OAUTH_BROWSER_SESSION_GEN_FAILED", "failed to generate oauth browser session").WithCause(err))
		return
placeholder

	secureCookie := isRequestHTTPS(c)
	setCookie(c, linuxDoOAuthStateCookieName, encodeCookieValue(state), linuxDoOAuthCookieMaxAgeSec, secureCookie)
	setCookie(c, linuxDoOAuthRedirectCookie, encodeCookieValue(redirectTo), linuxDoOAuthCookieMaxAgeSec, secureCookie)
	intent := normalizeOAuthIntent(c.Query("intent"))
	setCookie(c, linuxDoOAuthIntentCookieName, encodeCookieValue(intent), linuxDoOAuthCookieMaxAgeSec, secureCookie)
	captureOAuthPromoCode(c, secureCookie)
	setOAuthPendingBrowserCookie(c, browserSessionKey, secureCookie)
	clearOAuthPendingSessionCookie(c, secureCookie)
	if intent == oauthIntentBindCurrentUser {
		bindCookieValue, err := h.buildOAuthBindUserCookieFromContext(c)
		if err != nil {
			response.ErrorFrom(c, err)
			return
	placeholder
		setCookie(c, linuxDoOAuthBindUserCookieName, encodeCookieValue(bindCookieValue), linuxDoOAuthCookieMaxAgeSec, secureCookie)
placeholder else {
		clearCookie(c, linuxDoOAuthBindUserCookieName, secureCookie)
placeholder

	codeChallenge := ""
	if cfg.UsePKCE {
		verifier, err := oauth.GenerateCodeVerifier()
		if err != nil {
			response.ErrorFrom(c, infraerrors.InternalServer("OAUTH_PKCE_GEN_FAILED", "failed to generate pkce verifier").WithCause(err))
			return
	placeholder
		codeChallenge = oauth.GenerateCodeChallenge(verifier)
		setCookie(c, linuxDoOAuthVerifierCookie, encodeCookieValue(verifier), linuxDoOAuthCookieMaxAgeSec, secureCookie)
placeholder

	redirectURI := strings.TrimSpace(cfg.RedirectURL)
	if redirectURI == "" {
		response.ErrorFrom(c, infraerrors.InternalServer("OAUTH_CONFIG_INVALID", "oauth redirect url not configured"))
		return
placeholder

	authURL, err := buildLinuxDoAuthorizeURL(cfg, state, codeChallenge, redirectURI)
	if err != nil {
		response.ErrorFrom(c, infraerrors.InternalServer("OAUTH_BUILD_URL_FAILED", "failed to build oauth authorization url").WithCause(err))
		return
placeholder

	respondOAuthStart(c, authURL)
placeholder

// LinuxDoOAuthCallback 处理 OAuth 回调：创建/登录用户，然后重定向到前端。
// GET /api/v1/auth/oauth/linuxdo/callback?code=...&state=...
func (h *AuthHandler) LinuxDoOAuthCallback(c *gin.Context) {
	cfg, cfgErr := h.getLinuxDoOAuthConfig(c.Request.Context())
	if cfgErr != nil {
		response.ErrorFrom(c, cfgErr)
		return
placeholder

	frontendCallback := strings.TrimSpace(cfg.FrontendRedirectURL)
	if frontendCallback == "" {
		frontendCallback = linuxDoOAuthDefaultFrontendCB
placeholder

	if providerErr := strings.TrimSpace(c.Query("error")); providerErr != "" {
		redirectOAuthError(c, frontendCallback, "provider_error", providerErr, c.Query("error_description"))
		return
placeholder

	code := strings.TrimSpace(c.Query("code"))
	state := strings.TrimSpace(c.Query("state"))
	if code == "" || state == "" {
		redirectOAuthError(c, frontendCallback, "missing_params", "missing code/state", "")
		return
placeholder

	secureCookie := isRequestHTTPS(c)
	defer func() {
		clearCookie(c, linuxDoOAuthStateCookieName, secureCookie)
		clearCookie(c, linuxDoOAuthVerifierCookie, secureCookie)
		clearCookie(c, linuxDoOAuthRedirectCookie, secureCookie)
		clearCookie(c, linuxDoOAuthIntentCookieName, secureCookie)
		clearCookie(c, linuxDoOAuthBindUserCookieName, secureCookie)
		clearOAuthPromoCodeCookie(c, secureCookie)
placeholder()

	expectedState, err := readCookieDecoded(c, linuxDoOAuthStateCookieName)
	if err != nil || expectedState == "" || state != expectedState {
		redirectOAuthError(c, frontendCallback, "invalid_state", "invalid oauth state", "")
		return
placeholder

	redirectTo, _ := readCookieDecoded(c, linuxDoOAuthRedirectCookie)
	redirectTo = sanitizeFrontendRedirectPath(redirectTo)
	if redirectTo == "" {
		redirectTo = linuxDoOAuthDefaultRedirectTo
placeholder
	browserSessionKey, _ := readOAuthPendingBrowserCookie(c)
	if strings.TrimSpace(browserSessionKey) == "" {
		redirectOAuthError(c, frontendCallback, "missing_browser_session", "missing oauth browser session", "")
		return
placeholder
	intent, _ := readCookieDecoded(c, linuxDoOAuthIntentCookieName)
	intent = normalizeOAuthIntent(intent)

	codeVerifier := ""
	if cfg.UsePKCE {
		codeVerifier, _ = readCookieDecoded(c, linuxDoOAuthVerifierCookie)
		if codeVerifier == "" {
			redirectOAuthError(c, frontendCallback, "missing_verifier", "missing pkce verifier", "")
			return
	placeholder
placeholder

	redirectURI := strings.TrimSpace(cfg.RedirectURL)
	if redirectURI == "" {
		redirectOAuthError(c, frontendCallback, "config_error", "oauth redirect url not configured", "")
		return
placeholder

	tokenResp, err := linuxDoExchangeCode(c.Request.Context(), cfg, code, redirectURI, codeVerifier)
	if err != nil {
		description := ""
		var exchangeErr *linuxDoTokenExchangeError
		if errors.As(err, &exchangeErr) && exchangeErr != nil {
			log.Printf(
				"[LinuxDo OAuth] token exchange failed: status=%d provider_error=%q provider_description=%q body=%s",
				exchangeErr.StatusCode,
				exchangeErr.ProviderError,
				exchangeErr.ProviderDescription,
				truncateLogValue(exchangeErr.Body, 2048),
			)
			description = exchangeErr.Error()
	placeholder else {
			log.Printf("[LinuxDo OAuth] token exchange failed: %v", err)
			description = err.Error()
	placeholder
		redirectOAuthError(c, frontendCallback, "token_exchange_failed", "failed to exchange oauth code", singleLine(description))
		return
placeholder

	email, username, subject, displayName, avatarURL, err := linuxDoFetchUserInfo(c.Request.Context(), cfg, tokenResp)
	if err != nil {
		log.Printf("[LinuxDo OAuth] userinfo fetch failed: %v", err)
		redirectOAuthError(c, frontendCallback, "userinfo_failed", "failed to fetch user info", "")
		return
placeholder
	compatEmail := strings.TrimSpace(email)

	// 安全考虑：不要把第三方返回的 email 直接映射到本地账号（可能与本地邮箱用户冲突导致账号被接管）。
	// 统一使用基于 subject 的稳定合成邮箱来做账号绑定。
	if subject != "" {
		email = linuxDoSyntheticEmail(subject)
placeholder
	identityKey := service.PendingAuthIdentityKey{
		ProviderType:    "linuxdo",
		ProviderKey:     "linuxdo",
		ProviderSubject: subject,
placeholder
	upstreamClaims := map[string]any{
		"email":                  email,
		"username":               username,
		"subject":                subject,
		"suggested_display_name": displayName,
		"suggested_avatar_url":   avatarURL,
placeholder
	if compatEmail != "" && !strings.EqualFold(strings.TrimSpace(compatEmail), strings.TrimSpace(email)) {
		upstreamClaims["compat_email"] = compatEmail
placeholder
	if intent == oauthIntentBindCurrentUser {
		targetUserID, err := h.readOAuthBindUserIDFromCookie(c, linuxDoOAuthBindUserCookieName)
		if err != nil {
			redirectOAuthError(c, frontendCallback, "invalid_state", "invalid oauth bind target", "")
			return
	placeholder
		if err := h.createOAuthPendingSession(c, oauthPendingSessionPayload{
			Intent:                 oauthIntentBindCurrentUser,
			Identity:               identityKey,
			TargetUserID:           &targetUserID,
			ResolvedEmail:          email,
			RedirectTo:             redirectTo,
			BrowserSessionKey:      browserSessionKey,
			UpstreamIdentityClaims: upstreamClaims,
			CompletionResponse: map[string]any{
				"redirect": redirectTo,
		placeholder,
	placeholder); err != nil {
			redirectOAuthError(c, frontendCallback, "session_error", "failed to continue oauth bind", "")
			return
	placeholder
		redirectToFrontendCallback(c, frontendCallback)
		return
placeholder

	existingIdentityUser, err := h.findOAuthIdentityUser(c.Request.Context(), identityKey)
	if err != nil {
		redirectOAuthError(c, frontendCallback, "session_error", infraerrors.Reason(err), infraerrors.Message(err))
		return
placeholder
	if existingIdentityUser != nil {
		if err := h.createOAuthPendingSession(c, oauthPendingSessionPayload{
			Intent:                 oauthIntentLogin,
			Identity:               identityKey,
			TargetUserID:           &existingIdentityUser.ID,
			ResolvedEmail:          existingIdentityUser.Email,
			RedirectTo:             redirectTo,
			BrowserSessionKey:      browserSessionKey,
			UpstreamIdentityClaims: upstreamClaims,
			CompletionResponse: map[string]any{
				"redirect": redirectTo,
		placeholder,
	placeholder); err != nil {
			redirectOAuthError(c, frontendCallback, "session_error", "failed to continue oauth login", "")
			return
	placeholder
		redirectToFrontendCallback(c, frontendCallback)
		return
placeholder

	compatEmailUser, err := h.findLinuxDoCompatEmailUser(c.Request.Context(), compatEmail)
	if err != nil {
		redirectOAuthError(c, frontendCallback, "session_error", infraerrors.Reason(err), infraerrors.Message(err))
		return
placeholder
	emailVerificationRequired := h != nil && h.authService != nil && h.authService.IsEmailVerifyEnabled(c.Request.Context())
	forceEmailOnSignup := h.isForceEmailOnThirdPartySignup(c.Request.Context())
	if compatEmailUser == nil && !emailVerificationRequired && !forceEmailOnSignup {
		if err := h.ensureBackendModeAllowsNewUserLogin(c.Request.Context()); err != nil {
			redirectOAuthError(c, frontendCallback, "session_error", infraerrors.Reason(err), infraerrors.Message(err))
			return
	placeholder
		tokenPair, user, err := h.authService.LoginOrRegisterOAuthWithTokenPairAndPromoCode(
			c.Request.Context(),
			email,
			username,
			"",
			"",
			readOAuthPromoCode(c),
			"linuxdo",
		)
		if err == nil {
			if err := applyPendingOAuthBinding(
				c.Request.Context(),
				h.entClient(),
				h.authService,
				h.userService,
				&dbent.PendingAuthSession{
					Intent:                 oauthIntentLogin,
					ProviderType:           identityKey.ProviderType,
					ProviderKey:            identityKey.ProviderKey,
					ProviderSubject:        identityKey.ProviderSubject,
					ResolvedEmail:          email,
					UpstreamIdentityClaims: upstreamClaims,
			placeholder,
				nil,
				&user.ID,
				true,
				false,
			); err != nil {
				redirectOAuthError(c, frontendCallback, "session_error", "failed to bind oauth identity", "")
				return
		placeholder
			h.authService.RecordSuccessfulLogin(c.Request.Context(), user.ID)
			clearOAuthPendingSessionCookie(c, secureCookie)
			clearOAuthPendingBrowserCookie(c, secureCookie)
			redirectOAuthTokenPair(c, frontendCallback, tokenPair, redirectTo)
			return
	placeholder
		if !errors.Is(err, service.ErrOAuthInvitationRequired) {
			redirectOAuthError(c, frontendCallback, "session_error", infraerrors.Reason(err), infraerrors.Message(err))
			return
	placeholder
placeholder
	if err := h.createLinuxDoOAuthChoicePendingSession(
		c,
		identityKey,
		email,
		email,
		redirectTo,
		browserSessionKey,
		upstreamClaims,
		compatEmail,
		compatEmailUser,
		emailVerificationRequired,
		forceEmailOnSignup,
	); err != nil {
		redirectOAuthError(c, frontendCallback, "session_error", "failed to continue oauth login", "")
		return
placeholder
	redirectToFrontendCallback(c, frontendCallback)
placeholder

func (h *AuthHandler) findLinuxDoCompatEmailUser(ctx context.Context, email string) (*dbent.User, error) {
	client := h.entClient()
	if client == nil {
		return nil, infraerrors.ServiceUnavailable("PENDING_AUTH_NOT_READY", "pending auth service is not ready")
placeholder

	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" ||
		strings.HasSuffix(email, service.LinuxDoConnectSyntheticEmailDomain) ||
		strings.HasSuffix(email, service.OIDCConnectSyntheticEmailDomain) ||
		strings.HasSuffix(email, service.WeChatConnectSyntheticEmailDomain) ||
		strings.HasSuffix(email, service.DingTalkConnectSyntheticEmailDomain) {
		return nil, nil
placeholder

	userEntity, err := client.User.Query().
		Where(userNormalizedEmailPredicate(email)).
		Order(dbent.Asc(dbuser.FieldID)).
		All(ctx)
	if err != nil {
		return nil, infraerrors.InternalServer("COMPAT_EMAIL_LOOKUP_FAILED", "failed to look up compat email user").WithCause(err)
placeholder
	switch len(userEntity) {
	case 0:
		return nil, nil
	case 1:
		return userEntity[0], nil
	default:
		return nil, infraerrors.Conflict("USER_EMAIL_CONFLICT", "normalized email matched multiple users")
placeholder
placeholder

func (h *AuthHandler) createLinuxDoOAuthChoicePendingSession(
	c *gin.Context,
	identity service.PendingAuthIdentityKey,
	suggestedEmail string,
	resolvedEmail string,
	redirectTo string,
	browserSessionKey string,
	upstreamClaims map[string]any,
	compatEmail string,
	compatEmailUser *dbent.User,
	emailVerificationRequired bool,
	forceEmailOnSignup bool,
) error {
	suggestionEmail := strings.TrimSpace(suggestedEmail)
	canonicalEmail := strings.TrimSpace(resolvedEmail)
	if suggestionEmail == "" {
		suggestionEmail = canonicalEmail
placeholder

	completionResponse := map[string]any{
		"step":                      oauthPendingChoiceStep,
		"adoption_required":         true,
		"redirect":                  strings.TrimSpace(redirectTo),
		"email":                     suggestionEmail,
		"resolved_email":            canonicalEmail,
		"existing_account_email":    "",
		"existing_account_bindable": false,
		"create_account_allowed":    true,
		"force_email_on_signup":     forceEmailOnSignup,
		"choice_reason":             "third_party_signup",
placeholder
	if strings.TrimSpace(compatEmail) != "" {
		completionResponse["compat_email"] = strings.TrimSpace(compatEmail)
placeholder
	resolvedChoiceEmail := suggestionEmail
	if compatEmailUser != nil {
		completionResponse["email"] = strings.TrimSpace(compatEmailUser.Email)
		completionResponse["existing_account_email"] = strings.TrimSpace(compatEmailUser.Email)
		completionResponse["existing_account_bindable"] = true
		completionResponse["choice_reason"] = "compat_email_match"
		resolvedChoiceEmail = strings.TrimSpace(compatEmailUser.Email)
placeholder
	if forceEmailOnSignup && compatEmailUser == nil {
		completionResponse["choice_reason"] = "force_email_on_signup"
placeholder
	if (emailVerificationRequired || forceEmailOnSignup) && compatEmailUser == nil {
		completionResponse["step"] = "create_account_required"
		completionResponse["email_binding_required"] = true
		completionResponse["force_email_on_signup"] = true
		if emailVerificationRequired {
			completionResponse["choice_reason"] = "email_verification_required"
	placeholder
		delete(completionResponse, "email")
		delete(completionResponse, "resolved_email")
		resolvedChoiceEmail = ""
placeholder

	var targetUserID *int64
	if compatEmailUser != nil && compatEmailUser.ID > 0 {
		targetUserID = &compatEmailUser.ID
placeholder

	return h.createOAuthPendingSession(c, oauthPendingSessionPayload{
		Intent:                 oauthIntentLogin,
		Identity:               identity,
		TargetUserID:           targetUserID,
		ResolvedEmail:          resolvedChoiceEmail,
		RedirectTo:             redirectTo,
		BrowserSessionKey:      browserSessionKey,
		UpstreamIdentityClaims: upstreamClaims,
		CompletionResponse:     completionResponse,
placeholder)
placeholder

type completeLinuxDoOAuthRequest struct {
	InvitationCode   string `json:"invitation_code" binding:"required"`
	AffCode          string `json:"aff_code,omitempty"`
	AdoptDisplayName *bool  `json:"adopt_display_name,omitempty"`
	AdoptAvatar      *bool  `json:"adopt_avatar,omitempty"`
placeholder

// CompleteLinuxDoOAuthRegistration completes a pending OAuth registration by validating
// the invitation code and creating the user account.
// POST /api/v1/auth/oauth/linuxdo/complete-registration
func (h *AuthHandler) CompleteLinuxDoOAuthRegistration(c *gin.Context) {
	var req completeLinuxDoOAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_REQUEST", "message": err.Error()placeholder)
		return
placeholder

	secureCookie := isRequestHTTPS(c)
	sessionToken, err := readOAuthPendingSessionCookie(c)
	if err != nil {
		clearOAuthPendingSessionCookie(c, secureCookie)
		clearOAuthPendingBrowserCookie(c, secureCookie)
		response.ErrorFrom(c, service.ErrPendingAuthSessionNotFound)
		return
placeholder
	browserSessionKey, err := readOAuthPendingBrowserCookie(c)
	if err != nil {
		clearOAuthPendingSessionCookie(c, secureCookie)
		clearOAuthPendingBrowserCookie(c, secureCookie)
		response.ErrorFrom(c, service.ErrPendingAuthBrowserMismatch)
		return
placeholder
	pendingSvc, err := h.pendingIdentityService()
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	session, err := pendingSvc.GetBrowserSession(c.Request.Context(), sessionToken, browserSessionKey)
	if err != nil {
		clearOAuthPendingSessionCookie(c, secureCookie)
		clearOAuthPendingBrowserCookie(c, secureCookie)
		response.ErrorFrom(c, err)
		return
placeholder
	if err := ensurePendingOAuthCompleteRegistrationSession(session); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	if updatedSession, handled, err := h.legacyCompleteRegistrationSessionStatus(c, session); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder else if handled {
		c.JSON(http.StatusOK, buildPendingOAuthSessionStatusPayload(updatedSession))
		return
placeholder else {
		session = updatedSession
placeholder
	if err := h.ensureBackendModeAllowsNewUserLogin(c.Request.Context()); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	email := strings.TrimSpace(session.ResolvedEmail)
	username := pendingSessionStringValue(session.UpstreamIdentityClaims, "username")
	if email == "" || username == "" {
		response.ErrorFrom(c, infraerrors.BadRequest("PENDING_AUTH_SESSION_INVALID", "pending auth registration context is invalid"))
		return
placeholder

	client := h.entClient()
	if client == nil {
		response.ErrorFrom(c, infraerrors.ServiceUnavailable("PENDING_AUTH_NOT_READY", "pending auth service is not ready"))
		return
placeholder
	if err := ensurePendingOAuthRegistrationIdentityAvailable(c.Request.Context(), client, session); err != nil {
		respondPendingOAuthBindingApplyError(c, err)
		return
placeholder
	decision, err := h.ensurePendingOAuthAdoptionDecision(c, session.ID, oauthAdoptionDecisionRequest{
		AdoptDisplayName: req.AdoptDisplayName,
		AdoptAvatar:      req.AdoptAvatar,
placeholder)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	tokenPair, user, err := h.authService.LoginOrRegisterOAuthWithTokenPairAndPromoCode(
		c.Request.Context(),
		email,
		username,
		req.InvitationCode,
		req.AffCode,
		pendingOAuthPromoCode(session),
		"linuxdo",
	)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	if err := applyPendingOAuthAdoptionAndConsumeSession(c.Request.Context(), client, h.authService, h.userService, session, decision, user.ID); err != nil {
		respondPendingOAuthBindingApplyError(c, err)
		return
placeholder
	h.authService.RecordSuccessfulLogin(c.Request.Context(), user.ID)
	clearOAuthPendingSessionCookie(c, secureCookie)
	clearOAuthPendingBrowserCookie(c, secureCookie)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"expires_in":    tokenPair.ExpiresIn,
		"token_type":    "Bearer",
placeholder)
placeholder

func (h *AuthHandler) getLinuxDoOAuthConfig(ctx context.Context) (config.LinuxDoConnectConfig, error) {
	if h != nil && h.settingSvc != nil {
		return h.settingSvc.GetLinuxDoConnectOAuthConfig(ctx)
placeholder
	if h == nil || h.cfg == nil {
		return config.LinuxDoConnectConfig{placeholder, infraerrors.ServiceUnavailable("CONFIG_NOT_READY", "config not loaded")
placeholder
	if !h.cfg.LinuxDo.Enabled {
		return config.LinuxDoConnectConfig{placeholder, infraerrors.NotFound("OAUTH_DISABLED", "oauth login is disabled")
placeholder
	return h.cfg.LinuxDo, nil
placeholder

func linuxDoExchangeCode(
	ctx context.Context,
	cfg config.LinuxDoConnectConfig,
	code string,
	redirectURI string,
	codeVerifier string,
) (*linuxDoTokenResponse, error) {
	client := req.C().SetTimeout(30 * time.Second)

	form := url.Values{placeholder
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", cfg.ClientID)
	form.Set("code", code)
	form.Set("redirect_uri", redirectURI)
	if strings.TrimSpace(codeVerifier) != "" {
		form.Set("code_verifier", codeVerifier)
placeholder

	r := client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json")

	switch strings.ToLower(strings.TrimSpace(cfg.TokenAuthMethod)) {
	case "", "client_secret_post":
		form.Set("client_secret", cfg.ClientSecret)
	case "client_secret_basic":
		r.SetBasicAuth(cfg.ClientID, cfg.ClientSecret)
	case "none":
	default:
		return nil, fmt.Errorf("unsupported token_auth_method: %s", cfg.TokenAuthMethod)
placeholder

	resp, err := r.SetFormDataFromValues(form).Post(cfg.TokenURL)
	if err != nil {
		return nil, fmt.Errorf("request token: %w", err)
placeholder
	body := strings.TrimSpace(resp.String())
	if !resp.IsSuccessState() {
		providerErr, providerDesc := parseOAuthProviderError(body)
		return nil, &linuxDoTokenExchangeError{
			StatusCode:          resp.StatusCode,
			ProviderError:       providerErr,
			ProviderDescription: providerDesc,
			Body:                body,
	placeholder
placeholder

	tokenResp, ok := parseLinuxDoTokenResponse(body)
	if !ok || strings.TrimSpace(tokenResp.AccessToken) == "" {
		return nil, &linuxDoTokenExchangeError{
			StatusCode: resp.StatusCode,
			Body:       body,
	placeholder
placeholder
	if strings.TrimSpace(tokenResp.TokenType) == "" {
		tokenResp.TokenType = "Bearer"
placeholder
	return tokenResp, nil
placeholder

func linuxDoFetchUserInfo(
	ctx context.Context,
	cfg config.LinuxDoConnectConfig,
	token *linuxDoTokenResponse,
) (email string, username string, subject string, displayName string, avatarURL string, err error) {
	client := req.C().SetTimeout(30 * time.Second)
	authorization, err := buildBearerAuthorization(token.TokenType, token.AccessToken)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("invalid token for userinfo request: %w", err)
placeholder

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", authorization).
		Get(cfg.UserInfoURL)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("request userinfo: %w", err)
placeholder
	if !resp.IsSuccessState() {
		return "", "", "", "", "", fmt.Errorf("userinfo status=%d", resp.StatusCode)
placeholder

	return linuxDoParseUserInfo(resp.String(), cfg)
placeholder

func linuxDoParseUserInfo(body string, cfg config.LinuxDoConnectConfig) (email string, username string, subject string, displayName string, avatarURL string, err error) {
	email = firstNonEmpty(
		getGJSON(body, cfg.UserInfoEmailPath),
		getGJSON(body, "email"),
		getGJSON(body, "user.email"),
		getGJSON(body, "data.email"),
		getGJSON(body, "attributes.email"),
	)
	username = firstNonEmpty(
		getGJSON(body, cfg.UserInfoUsernamePath),
		getGJSON(body, "username"),
		getGJSON(body, "preferred_username"),
		getGJSON(body, "name"),
		getGJSON(body, "user.username"),
		getGJSON(body, "user.name"),
	)
	subject = firstNonEmpty(
		getGJSON(body, cfg.UserInfoIDPath),
		getGJSON(body, "sub"),
		getGJSON(body, "id"),
		getGJSON(body, "user_id"),
		getGJSON(body, "uid"),
		getGJSON(body, "user.id"),
	)

	displayName = firstNonEmpty(
		getGJSON(body, "name"),
		getGJSON(body, "nickname"),
		getGJSON(body, "display_name"),
		getGJSON(body, "user.name"),
		getGJSON(body, "user.username"),
		username,
	)
	avatarURL = firstNonEmpty(
		getGJSON(body, "avatar_url"),
		getGJSON(body, "avatar"),
		getGJSON(body, "picture"),
		getGJSON(body, "profile_image_url"),
		getGJSON(body, "user.avatar"),
		getGJSON(body, "user.avatar_url"),
	)

	subject = strings.TrimSpace(subject)
	if subject == "" {
		return "", "", "", "", "", errors.New("userinfo missing id field")
placeholder
	if !isSafeLinuxDoSubject(subject) {
		return "", "", "", "", "", errors.New("userinfo returned invalid id field")
placeholder

	email = strings.TrimSpace(email)
	if email == "" {
		// LinuxDo Connect 的 userinfo 可能不提供 email。为兼容现有用户模型（email 必填且唯一），使用稳定的合成邮箱。
		email = linuxDoSyntheticEmail(subject)
placeholder

	username = strings.TrimSpace(username)
	if username == "" {
		username = "linuxdo_" + subject
placeholder
	displayName = strings.TrimSpace(displayName)
	if displayName == "" {
		displayName = username
placeholder
	avatarURL = strings.TrimSpace(avatarURL)

	return email, username, subject, displayName, avatarURL, nil
placeholder

func buildLinuxDoAuthorizeURL(cfg config.LinuxDoConnectConfig, state string, codeChallenge string, redirectURI string) (string, error) {
	u, err := url.Parse(cfg.AuthorizeURL)
	if err != nil {
		return "", fmt.Errorf("parse authorize_url: %w", err)
placeholder

	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", cfg.ClientID)
	q.Set("redirect_uri", redirectURI)
	if strings.TrimSpace(cfg.Scopes) != "" {
		q.Set("scope", cfg.Scopes)
placeholder
	q.Set("state", state)
	if strings.TrimSpace(codeChallenge) != "" {
		q.Set("code_challenge", codeChallenge)
		q.Set("code_challenge_method", "S256")
placeholder

	u.RawQuery = q.Encode()
	return u.String(), nil
placeholder

func redirectOAuthError(c *gin.Context, frontendCallback string, code string, message string, description string) {
	fragment := url.Values{placeholder
	fragment.Set("error", truncateFragmentValue(code))
	if strings.TrimSpace(message) != "" {
		fragment.Set("error_message", truncateFragmentValue(message))
placeholder
	if strings.TrimSpace(description) != "" {
		fragment.Set("error_description", truncateFragmentValue(description))
placeholder
	redirectWithFragment(c, frontendCallback, fragment)
placeholder

func redirectOAuthTokenPair(c *gin.Context, frontendCallback string, tokenPair *service.TokenPair, redirectTo string) {
	fragment := url.Values{placeholder
	if tokenPair != nil {
		fragment.Set("access_token", truncateFragmentValue(tokenPair.AccessToken))
		fragment.Set("refresh_token", truncateFragmentValue(tokenPair.RefreshToken))
		fragment.Set("expires_in", strconv.Itoa(tokenPair.ExpiresIn))
		fragment.Set("token_type", "Bearer")
placeholder
	if redirect := strings.TrimSpace(redirectTo); redirect != "" {
		originalRedirect := redirect
		for range 2 {
			decoded, err := url.QueryUnescape(redirect)
			if err != nil || decoded == redirect {
				break
		placeholder
			redirect = decoded
	placeholder
		if redirect != originalRedirect {
			if sanitized := sanitizeFrontendRedirectPath(redirect); sanitized != "" {
				redirect = sanitized
		placeholder else {
				redirect = originalRedirect
		placeholder
	placeholder
		fragment.Set("redirect", truncateFragmentValue(redirect))
placeholder
	redirectWithFragment(c, frontendCallback, fragment)
placeholder

func redirectWithFragment(c *gin.Context, frontendCallback string, fragment url.Values) {
	u, err := url.Parse(frontendCallback)
	if err != nil {
		// 兜底：尽力跳转到默认页面，避免卡死在回调页。
		c.Redirect(http.StatusFound, linuxDoOAuthDefaultRedirectTo)
		return
placeholder
	if u.Scheme != "" && !strings.EqualFold(u.Scheme, "http") && !strings.EqualFold(u.Scheme, "https") {
		c.Redirect(http.StatusFound, linuxDoOAuthDefaultRedirectTo)
		return
placeholder
	u.Fragment = fragment.Encode()
	c.Header("Cache-Control", "no-store")
	c.Header("Pragma", "no-cache")
	c.Redirect(http.StatusFound, u.String())
placeholder

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v != "" {
			return v
	placeholder
placeholder
	return ""
placeholder

func parseOAuthProviderError(body string) (providerErr string, providerDesc string) {
	body = strings.TrimSpace(body)
	if body == "" {
		return "", ""
placeholder

	providerErr = firstNonEmpty(
		getGJSON(body, "error"),
		getGJSON(body, "code"),
		getGJSON(body, "error.code"),
	)
	providerDesc = firstNonEmpty(
		getGJSON(body, "error_description"),
		getGJSON(body, "error.message"),
		getGJSON(body, "message"),
		getGJSON(body, "detail"),
	)

	if providerErr != "" || providerDesc != "" {
		return providerErr, providerDesc
placeholder

	values, err := url.ParseQuery(body)
	if err != nil {
		return "", ""
placeholder
	providerErr = firstNonEmpty(values.Get("error"), values.Get("code"))
	providerDesc = firstNonEmpty(values.Get("error_description"), values.Get("error_message"), values.Get("message"))
	return providerErr, providerDesc
placeholder

func parseLinuxDoTokenResponse(body string) (*linuxDoTokenResponse, bool) {
	body = strings.TrimSpace(body)
	if body == "" {
		return nil, false
placeholder

	accessToken := strings.TrimSpace(getGJSON(body, "access_token"))
	if accessToken != "" {
		tokenType := strings.TrimSpace(getGJSON(body, "token_type"))
		refreshToken := strings.TrimSpace(getGJSON(body, "refresh_token"))
		scope := strings.TrimSpace(getGJSON(body, "scope"))
		expiresIn := gjson.Get(body, "expires_in").Int()
		return &linuxDoTokenResponse{
			AccessToken:  accessToken,
			TokenType:    tokenType,
			ExpiresIn:    expiresIn,
			RefreshToken: refreshToken,
			Scope:        scope,
	placeholder, true
placeholder

	values, err := url.ParseQuery(body)
	if err != nil {
		return nil, false
placeholder
	accessToken = strings.TrimSpace(values.Get("access_token"))
	if accessToken == "" {
		return nil, false
placeholder
	expiresIn := int64(0)
	if raw := strings.TrimSpace(values.Get("expires_in")); raw != "" {
		if v, err := strconv.ParseInt(raw, 10, 64); err == nil {
			expiresIn = v
	placeholder
placeholder
	return &linuxDoTokenResponse{
		AccessToken:  accessToken,
		TokenType:    strings.TrimSpace(values.Get("token_type")),
		ExpiresIn:    expiresIn,
		RefreshToken: strings.TrimSpace(values.Get("refresh_token")),
		Scope:        strings.TrimSpace(values.Get("scope")),
placeholder, true
placeholder

func getGJSON(body string, path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
placeholder
	res := gjson.Get(body, path)
	if !res.Exists() {
		return ""
placeholder
	return res.String()
placeholder

func truncateLogValue(value string, maxLen int) string {
	value = strings.TrimSpace(value)
	if value == "" || maxLen <= 0 {
		return ""
placeholder
	if len(value) <= maxLen {
		return value
placeholder
	value = value[:maxLen]
	for !utf8.ValidString(value) {
		value = value[:len(value)-1]
placeholder
	return value
placeholder

func singleLine(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
placeholder
	return strings.Join(strings.Fields(value), " ")
placeholder

func sanitizeFrontendRedirectPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
placeholder
	if len(path) > linuxDoOAuthMaxRedirectLen {
		return ""
placeholder
	// 只允许同源相对路径（避免开放重定向）。
	if !strings.HasPrefix(path, "/") {
		return ""
placeholder
	if strings.HasPrefix(path, "//") {
		return ""
placeholder
	if strings.Contains(path, "://") {
		return ""
placeholder
	if strings.ContainsAny(path, "\r\n") {
		return ""
placeholder
	return path
placeholder

func isRequestHTTPS(c *gin.Context) bool {
	if c.Request.TLS != nil {
		return true
placeholder
	proto := strings.ToLower(strings.TrimSpace(c.GetHeader("X-Forwarded-Proto")))
	return proto == "https"
placeholder

func encodeCookieValue(value string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(value))
placeholder

func decodeCookieValue(value string) (string, error) {
	raw, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return "", err
placeholder
	return string(raw), nil
placeholder

func readCookieDecoded(c *gin.Context, name string) (string, error) {
	ck, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
placeholder
	return decodeCookieValue(ck.Value)
placeholder

func setCookie(c *gin.Context, name string, value string, maxAgeSec int, secure bool) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     linuxDoOAuthCookiePath,
		MaxAge:   maxAgeSec,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
placeholder)
placeholder

func clearCookie(c *gin.Context, name string, secure bool) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     linuxDoOAuthCookiePath,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
placeholder)
placeholder

func clearOAuthBindAccessTokenCookie(c *gin.Context, secure bool) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     oauthBindAccessTokenCookieName,
		Value:    "",
		Path:     oauthBindAccessTokenCookiePath,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
placeholder)
placeholder

func setOAuthBindAccessTokenCookie(c *gin.Context, token string, secure bool) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     oauthBindAccessTokenCookieName,
		Value:    url.QueryEscape(strings.TrimSpace(token)),
		Path:     oauthBindAccessTokenCookiePath,
		MaxAge:   linuxDoOAuthCookieMaxAgeSec,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
placeholder)
placeholder

func truncateFragmentValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
placeholder
	if len(value) > linuxDoOAuthMaxFragmentValueLen {
		value = value[:linuxDoOAuthMaxFragmentValueLen]
		for !utf8.ValidString(value) {
			value = value[:len(value)-1]
	placeholder
placeholder
	return value
placeholder

func buildBearerAuthorization(tokenType, accessToken string) (string, error) {
	tokenType = strings.TrimSpace(tokenType)
	if tokenType == "" {
		tokenType = "Bearer"
placeholder
	if !strings.EqualFold(tokenType, "Bearer") {
		return "", fmt.Errorf("unsupported token_type: %s", tokenType)
placeholder

	accessToken = strings.TrimSpace(accessToken)
	if accessToken == "" {
		return "", errors.New("missing access_token")
placeholder
	if strings.ContainsAny(accessToken, " \t\r\n") {
		return "", errors.New("access_token contains whitespace")
placeholder
	return "Bearer " + accessToken, nil
placeholder

func isSafeLinuxDoSubject(subject string) bool {
	subject = strings.TrimSpace(subject)
	if subject == "" || len(subject) > linuxDoOAuthMaxSubjectLen {
		return false
placeholder
	for _, r := range subject {
		switch {
		case r >= '0' && r <= '9':
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r == '_' || r == '-':
		default:
			return false
	placeholder
placeholder
	return true
placeholder

func linuxDoSyntheticEmail(subject string) string {
	subject = strings.TrimSpace(subject)
	if subject == "" {
		return ""
placeholder
	return "linuxdo-" + subject + service.LinuxDoConnectSyntheticEmailDomain
placeholder

func normalizeOAuthIntent(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", oauthIntentLogin:
		return oauthIntentLogin
	case "bind", oauthIntentBindCurrentUser:
		return oauthIntentBindCurrentUser
	default:
		return oauthIntentLogin
placeholder
placeholder

func (h *AuthHandler) buildOAuthBindUserCookieFromContext(c *gin.Context) (string, error) {
	userID, err := h.resolveOAuthBindTargetUserID(c)
	if err != nil || userID == nil || *userID <= 0 {
		return "", infraerrors.Unauthorized("UNAUTHORIZED", "authentication required")
placeholder
	return buildOAuthBindUserCookieValue(*userID, h.oauthBindCookieSecret())
placeholder

func (h *AuthHandler) PrepareOAuthBindAccessTokenCookie(c *gin.Context) {
	const bearerPrefix = "Bearer "

	authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
	if !strings.HasPrefix(strings.ToLower(authHeader), strings.ToLower(bearerPrefix)) {
		response.ErrorFrom(c, infraerrors.Unauthorized("UNAUTHORIZED", "authentication required"))
		return
placeholder

	token := strings.TrimSpace(authHeader[len(bearerPrefix):])
	if token == "" {
		response.ErrorFrom(c, infraerrors.Unauthorized("UNAUTHORIZED", "authentication required"))
		return
placeholder

	setOAuthBindAccessTokenCookie(c, token, isRequestHTTPS(c))
	c.Status(http.StatusNoContent)
	c.Writer.WriteHeaderNow()
placeholder

func (h *AuthHandler) resolveOAuthBindTargetUserID(c *gin.Context) (*int64, error) {
	if subject, ok := servermiddleware.GetAuthSubjectFromContext(c); ok && subject.UserID > 0 {
		return &subject.UserID, nil
placeholder
	if h == nil || h.authService == nil || h.userService == nil {
		return nil, service.ErrInvalidToken
placeholder

	ck, err := c.Request.Cookie(oauthBindAccessTokenCookieName)
	clearOAuthBindAccessTokenCookie(c, isRequestHTTPS(c))
	if err != nil {
		return nil, err
placeholder

	tokenString, err := url.QueryUnescape(strings.TrimSpace(ck.Value))
	if err != nil {
		return nil, err
placeholder
	if tokenString == "" {
		return nil, service.ErrInvalidToken
placeholder

	claims, err := h.authService.ValidateToken(tokenString)
	if err != nil {
		return nil, err
placeholder
	user, err := h.userService.GetByID(c.Request.Context(), claims.UserID)
	if err != nil {
		return nil, err
placeholder
	if user == nil || !user.IsActive() || claims.TokenVersion != user.TokenVersion {
		return nil, service.ErrInvalidToken
placeholder
	return &user.ID, nil
placeholder

func (h *AuthHandler) readOAuthBindUserIDFromCookie(c *gin.Context, cookieName string) (int64, error) {
	value, err := readCookieDecoded(c, cookieName)
	if err != nil {
		return 0, err
placeholder
	return parseOAuthBindUserCookieValue(value, h.oauthBindCookieSecret())
placeholder

func (h *AuthHandler) oauthBindCookieSecret() string {
	if h == nil || h.cfg == nil {
		return ""
placeholder
	return strings.TrimSpace(h.cfg.JWT.Secret)
placeholder

func buildOAuthBindUserCookieValue(userID int64, secret string) (string, error) {
	secret = strings.TrimSpace(secret)
	if userID <= 0 || secret == "" {
		return "", errors.New("invalid oauth bind cookie input")
placeholder
	payload := strconv.FormatInt(userID, 10)
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return payload + "." + signature, nil
placeholder

func parseOAuthBindUserCookieValue(value string, secret string) (int64, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return 0, errors.New("missing oauth bind cookie secret")
placeholder
	payload, signature, ok := strings.Cut(strings.TrimSpace(value), ".")
	if !ok || payload == "" || signature == "" {
		return 0, errors.New("invalid oauth bind cookie")
placeholder
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	expectedSignature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return 0, errors.New("invalid oauth bind cookie signature")
placeholder
	userID, err := strconv.ParseInt(payload, 10, 64)
	if err != nil || userID <= 0 {
		return 0, errors.New("invalid oauth bind cookie user")
placeholder
	return userID, nil
placeholder
