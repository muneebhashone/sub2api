package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/oauth"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
)

const (
	emailOAuthCookiePath      = "/api/v1/auth/oauth"
	emailOAuthStateCookieName = "email_oauth_state"
	emailOAuthRedirectCookie  = "email_oauth_redirect"
	emailOAuthProviderCookie  = "email_oauth_provider"
	emailOAuthAffiliateCookie = "email_oauth_affiliate"
	emailOAuthCookieMaxAgeSec = 10 * 60
	emailOAuthDefaultRedirect = "/dashboard"
)

type emailOAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope,omitempty"`
placeholder

type emailOAuthProfile struct {
	Subject       string
	Email         string
	EmailVerified bool
	Username      string
	DisplayName   string
	AvatarURL     string
	Metadata      map[string]any
placeholder

func (h *AuthHandler) GitHubOAuthStart(c *gin.Context) { h.emailOAuthStart(c, "github") placeholder
func (h *AuthHandler) GoogleOAuthStart(c *gin.Context) { h.emailOAuthStart(c, "google") placeholder

func (h *AuthHandler) GitHubOAuthCallback(c *gin.Context) { h.emailOAuthCallback(c, "github") placeholder
func (h *AuthHandler) GoogleOAuthCallback(c *gin.Context) { h.emailOAuthCallback(c, "google") placeholder
func (h *AuthHandler) CompleteGitHubOAuthRegistration(c *gin.Context) {
	h.completeEmailOAuthRegistration(c, "github")
placeholder
func (h *AuthHandler) CompleteGoogleOAuthRegistration(c *gin.Context) {
	h.completeEmailOAuthRegistration(c, "google")
placeholder

func (h *AuthHandler) emailOAuthStart(c *gin.Context, provider string) {
	cfg, err := h.getEmailOAuthConfig(c.Request.Context(), provider)
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
		redirectTo = emailOAuthDefaultRedirect
placeholder

	secureCookie := isRequestHTTPS(c)
	emailOAuthSetCookie(c, emailOAuthStateCookieName, encodeCookieValue(state), secureCookie)
	emailOAuthSetCookie(c, emailOAuthRedirectCookie, encodeCookieValue(redirectTo), secureCookie)
	emailOAuthSetCookie(c, emailOAuthProviderCookie, encodeCookieValue(provider), secureCookie)
	if affCode := strings.TrimSpace(firstNonEmpty(c.Query("aff_code"), c.Query("aff"))); affCode != "" {
		emailOAuthSetCookie(c, emailOAuthAffiliateCookie, encodeCookieValue(affCode), secureCookie)
placeholder else {
		emailOAuthClearCookie(c, emailOAuthAffiliateCookie, secureCookie)
placeholder

	authURL, err := buildEmailOAuthAuthorizeURL(cfg, state)
	if err != nil {
		response.ErrorFrom(c, infraerrors.InternalServer("OAUTH_BUILD_URL_FAILED", "failed to build oauth authorization url").WithCause(err))
		return
placeholder
	c.Redirect(http.StatusFound, authURL)
placeholder

func (h *AuthHandler) emailOAuthCallback(c *gin.Context, provider string) {
	cfg, cfgErr := h.getEmailOAuthConfig(c.Request.Context(), provider)
	if cfgErr != nil {
		response.ErrorFrom(c, cfgErr)
		return
placeholder
	frontendCallback := strings.TrimSpace(cfg.FrontendRedirectURL)
	if frontendCallback == "" {
		frontendCallback = "/auth/oauth/callback"
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
		emailOAuthClearCookie(c, emailOAuthStateCookieName, secureCookie)
		emailOAuthClearCookie(c, emailOAuthRedirectCookie, secureCookie)
		emailOAuthClearCookie(c, emailOAuthProviderCookie, secureCookie)
		emailOAuthClearCookie(c, emailOAuthAffiliateCookie, secureCookie)
placeholder()
	expectedState, err := readCookieDecoded(c, emailOAuthStateCookieName)
	if err != nil || expectedState == "" || expectedState != state {
		redirectOAuthError(c, frontendCallback, "invalid_state", "invalid oauth state", "")
		return
placeholder
	expectedProvider, _ := readCookieDecoded(c, emailOAuthProviderCookie)
	if !strings.EqualFold(strings.TrimSpace(expectedProvider), provider) {
		redirectOAuthError(c, frontendCallback, "invalid_state", "invalid oauth provider", "")
		return
placeholder
	redirectTo, _ := readCookieDecoded(c, emailOAuthRedirectCookie)
	redirectTo = sanitizeFrontendRedirectPath(redirectTo)
	if redirectTo == "" {
		redirectTo = emailOAuthDefaultRedirect
placeholder

	tokenResp, err := exchangeEmailOAuthCode(c.Request.Context(), cfg, code)
	if err != nil {
		redirectOAuthError(c, frontendCallback, "token_exchange_failed", "failed to exchange oauth code", singleLine(err.Error()))
		return
placeholder
	profile, err := fetchEmailOAuthProfile(c.Request.Context(), provider, cfg, tokenResp)
	if err != nil {
		redirectOAuthError(c, frontendCallback, "userinfo_failed", "failed to fetch verified email", singleLine(err.Error()))
		return
placeholder
	h.emailOAuthCallbackWithProfile(c, provider, cfg, frontendCallback, redirectTo, profile)
placeholder

func (h *AuthHandler) emailOAuthCallbackWithProfile(
	c *gin.Context,
	provider string,
	cfg config.EmailOAuthProviderConfig,
	frontendCallback string,
	redirectTo string,
	profile *emailOAuthProfile,
) {
	input := service.EmailOAuthIdentityInput{
		ProviderType:     provider,
		ProviderKey:      provider,
		ProviderSubject:  profile.Subject,
		Email:            profile.Email,
		EmailVerified:    profile.EmailVerified,
		Username:         profile.Username,
		DisplayName:      profile.DisplayName,
		AvatarURL:        profile.AvatarURL,
		UpstreamMetadata: profile.Metadata,
placeholder
	affiliateCode := h.emailOAuthAffiliateCode(c)
	tokenPair, user, err := h.authService.LoginOrRegisterVerifiedEmailOAuthWithInvitation(c.Request.Context(), input, "", affiliateCode)
	if err != nil {
		if errors.Is(err, service.ErrOAuthInvitationRequired) {
			if pendingErr := h.createEmailOAuthInvitationPendingSession(c, provider, frontendCallback, redirectTo, profile); pendingErr != nil {
				redirectOAuthError(c, frontendCallback, infraerrors.Reason(pendingErr), infraerrors.Message(pendingErr), "")
				return
		placeholder
			redirectToFrontendCallback(c, frontendCallback)
			return
	placeholder
		redirectOAuthError(c, frontendCallback, infraerrors.Reason(err), infraerrors.Message(err), "")
		return
placeholder
	if err := h.ensureBackendModeAllowsUser(c.Request.Context(), user); err != nil {
		redirectOAuthError(c, frontendCallback, "login_blocked", infraerrors.Reason(err), infraerrors.Message(err))
		return
placeholder

	fragment := url.Values{placeholder
	fragment.Set("access_token", tokenPair.AccessToken)
	fragment.Set("refresh_token", tokenPair.RefreshToken)
	fragment.Set("expires_in", fmt.Sprintf("%d", tokenPair.ExpiresIn))
	fragment.Set("token_type", "Bearer")
	fragment.Set("redirect", redirectTo)
	redirectWithFragment(c, frontendCallback, fragment)
placeholder

func (h *AuthHandler) emailOAuthAffiliateCode(c *gin.Context) string {
	if c == nil {
		return ""
placeholder
	if code, err := readCookieDecoded(c, emailOAuthAffiliateCookie); err == nil {
		return strings.TrimSpace(code)
placeholder
	return ""
placeholder

func (h *AuthHandler) createEmailOAuthInvitationPendingSession(
	c *gin.Context,
	provider string,
	frontendCallback string,
	redirectTo string,
	profile *emailOAuthProfile,
) error {
	if h == nil || profile == nil {
		return infraerrors.ServiceUnavailable("PENDING_AUTH_NOT_READY", "pending auth service is not ready")
placeholder
	browserSessionKey, err := generateOAuthPendingBrowserSession()
	if err != nil {
		return infraerrors.InternalServer("PENDING_AUTH_SESSION_CREATE_FAILED", "failed to create pending auth session").WithCause(err)
placeholder
	setOAuthPendingBrowserCookie(c, browserSessionKey, isRequestHTTPS(c))

	email := strings.TrimSpace(strings.ToLower(profile.Email))
	username := strings.TrimSpace(profile.Username)
	affiliateCode := h.emailOAuthAffiliateCode(c)
	upstreamClaims := map[string]any{
		"email":            email,
		"email_verified":   profile.EmailVerified,
		"username":         username,
		"provider":         provider,
		"provider_key":     provider,
		"provider_subject": strings.TrimSpace(profile.Subject),
placeholder
	if strings.TrimSpace(profile.DisplayName) != "" {
		upstreamClaims["suggested_display_name"] = strings.TrimSpace(profile.DisplayName)
placeholder
	if strings.TrimSpace(profile.AvatarURL) != "" {
		upstreamClaims["suggested_avatar_url"] = strings.TrimSpace(profile.AvatarURL)
placeholder
	if affiliateCode != "" {
		upstreamClaims["aff_code"] = affiliateCode
placeholder
	for key, value := range profile.Metadata {
		if _, exists := upstreamClaims[key]; !exists {
			upstreamClaims[key] = value
	placeholder
placeholder

	completionResponse := map[string]any{
		"step":                      oauthPendingChoiceStep,
		"error":                     "invitation_required",
		"choice_reason":             "invitation_required",
		"adoption_required":         false,
		"create_account_allowed":    true,
		"existing_account_bindable": false,
		"force_email_on_signup":     true,
		"email":                     email,
		"resolved_email":            email,
		"provider":                  provider,
		"redirect":                  redirectTo,
placeholder
	if strings.TrimSpace(frontendCallback) != "" {
		completionResponse["frontend_callback"] = strings.TrimSpace(frontendCallback)
placeholder

	return h.createOAuthPendingSession(c, oauthPendingSessionPayload{
		Intent:                 oauthIntentLogin,
		Identity:               service.PendingAuthIdentityKey{ProviderType: provider, ProviderKey: provider, ProviderSubject: strings.TrimSpace(profile.Subject)placeholder,
		ResolvedEmail:          email,
		RedirectTo:             redirectTo,
		BrowserSessionKey:      browserSessionKey,
		UpstreamIdentityClaims: upstreamClaims,
		CompletionResponse:     completionResponse,
placeholder)
placeholder

type completeEmailOAuthRequest struct {
	InvitationCode string `json:"invitation_code" binding:"required"`
	AffCode        string `json:"aff_code,omitempty"`
placeholder

func (h *AuthHandler) completeEmailOAuthRegistration(c *gin.Context, provider string) {
	var req completeEmailOAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	_, session, clearCookies, err := readPendingOAuthBrowserSession(c, h)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	if err := ensurePendingOAuthCompleteRegistrationSession(session); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	if !strings.EqualFold(strings.TrimSpace(session.ProviderType), provider) {
		response.BadRequest(c, "Pending oauth session provider mismatch")
		return
placeholder
	if err := h.ensureBackendModeAllowsNewUserLogin(c.Request.Context()); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	affiliateCode := strings.TrimSpace(req.AffCode)
	if affiliateCode == "" {
		affiliateCode = pendingSessionStringValue(session.UpstreamIdentityClaims, "aff_code")
placeholder

	tokenPair, user, err := h.authService.LoginOrRegisterVerifiedEmailOAuthWithInvitation(
		c.Request.Context(),
		service.EmailOAuthIdentityInput{
			ProviderType:     strings.TrimSpace(session.ProviderType),
			ProviderKey:      strings.TrimSpace(session.ProviderKey),
			ProviderSubject:  strings.TrimSpace(session.ProviderSubject),
			Email:            strings.TrimSpace(session.ResolvedEmail),
			EmailVerified:    true,
			Username:         pendingSessionStringValue(session.UpstreamIdentityClaims, "username"),
			DisplayName:      pendingSessionStringValue(session.UpstreamIdentityClaims, "suggested_display_name"),
			AvatarURL:        pendingSessionStringValue(session.UpstreamIdentityClaims, "suggested_avatar_url"),
			UpstreamMetadata: clonePendingMap(session.UpstreamIdentityClaims),
	placeholder,
		strings.TrimSpace(req.InvitationCode),
		affiliateCode,
	)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	client := h.entClient()
	if client == nil {
		response.ErrorFrom(c, infraerrors.ServiceUnavailable("PENDING_AUTH_NOT_READY", "pending auth service is not ready"))
		return
placeholder
	tx, err := client.Tx(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, infraerrors.InternalServer("PENDING_AUTH_BIND_APPLY_FAILED", "failed to consume pending oauth session").WithCause(err))
		return
placeholder
	defer func() { _ = tx.Rollback() placeholder()
	if err := consumePendingOAuthBrowserSessionTx(c.Request.Context(), tx, session); err != nil {
		_ = tx.Rollback()
		clearCookies()
		response.ErrorFrom(c, err)
		return
placeholder
	if err := tx.Commit(); err != nil {
		response.ErrorFrom(c, infraerrors.InternalServer("PENDING_AUTH_BIND_APPLY_FAILED", "failed to consume pending oauth session").WithCause(err))
		return
placeholder
	h.authService.RecordSuccessfulLogin(c.Request.Context(), user.ID)
	clearCookies()
	writeOAuthTokenPairResponse(c, tokenPair)
placeholder

func (h *AuthHandler) getEmailOAuthConfig(ctx context.Context, provider string) (config.EmailOAuthProviderConfig, error) {
	if h != nil && h.settingSvc != nil {
		return h.settingSvc.GetEmailOAuthProviderConfig(ctx, provider)
placeholder
	return config.EmailOAuthProviderConfig{placeholder, infraerrors.ServiceUnavailable("CONFIG_NOT_READY", "config not loaded")
placeholder

func buildEmailOAuthAuthorizeURL(cfg config.EmailOAuthProviderConfig, state string) (string, error) {
	u, err := url.Parse(cfg.AuthorizeURL)
	if err != nil {
		return "", fmt.Errorf("parse authorize_url: %w", err)
placeholder
	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", cfg.ClientID)
	q.Set("redirect_uri", cfg.RedirectURL)
	q.Set("state", state)
	if strings.TrimSpace(cfg.Scopes) != "" {
		q.Set("scope", cfg.Scopes)
placeholder
	u.RawQuery = q.Encode()
	return u.String(), nil
placeholder

func exchangeEmailOAuthCode(ctx context.Context, cfg config.EmailOAuthProviderConfig, code string) (*emailOAuthTokenResponse, error) {
	resp, err := req.C().
		R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetFormData(map[string]string{
			"grant_type":    "authorization_code",
			"client_id":     cfg.ClientID,
			"client_secret": cfg.ClientSecret,
			"code":          code,
			"redirect_uri":  cfg.RedirectURL,
	placeholder).
		Post(cfg.TokenURL)
	if err != nil {
		return nil, err
placeholder
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("token endpoint status %d: %s", resp.StatusCode, truncateLogValue(resp.String(), 1024))
placeholder
	var tokenResp emailOAuthTokenResponse
	if err := json.Unmarshal(resp.Bytes(), &tokenResp); err != nil {
		return nil, err
placeholder
	if strings.TrimSpace(tokenResp.AccessToken) == "" {
		return nil, errors.New("missing access_token")
placeholder
	return &tokenResp, nil
placeholder

func fetchEmailOAuthProfile(ctx context.Context, provider string, cfg config.EmailOAuthProviderConfig, token *emailOAuthTokenResponse) (*emailOAuthProfile, error) {
	resp, err := req.C().
		R().
		SetContext(ctx).
		SetBearerAuthToken(token.AccessToken).
		SetHeader("Accept", "application/json").
		Get(cfg.UserInfoURL)
	if err != nil {
		return nil, err
placeholder
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("userinfo endpoint status %d: %s", resp.StatusCode, truncateLogValue(resp.String(), 1024))
placeholder
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "github":
		return parseGitHubOAuthProfile(ctx, cfg, token, resp.String())
	case "google":
		return parseGoogleOAuthProfile(resp.String())
	default:
		return nil, errors.New("unsupported oauth provider")
placeholder
placeholder

func parseGitHubOAuthProfile(ctx context.Context, cfg config.EmailOAuthProviderConfig, token *emailOAuthTokenResponse, body string) (*emailOAuthProfile, error) {
	subject := strings.TrimSpace(gjson.Get(body, "id").String())
	if subject == "" {
		return nil, errors.New("github user id is missing")
placeholder
	email := strings.TrimSpace(gjson.Get(body, "email").String())
	emailVerified := false
	if email != "" {
		emailVerified = true
placeholder
	if strings.TrimSpace(cfg.EmailsURL) != "" {
		if verifiedEmail, err := fetchGitHubPrimaryVerifiedEmail(ctx, cfg.EmailsURL, token.AccessToken); err == nil && verifiedEmail != "" {
			email = verifiedEmail
			emailVerified = true
	placeholder else if email == "" && err != nil {
			return nil, err
	placeholder
placeholder
	if email == "" || !emailVerified {
		return nil, errors.New("github verified email is missing")
placeholder
	login := strings.TrimSpace(gjson.Get(body, "login").String())
	name := strings.TrimSpace(gjson.Get(body, "name").String())
	return &emailOAuthProfile{
		Subject:       subject,
		Email:         email,
		EmailVerified: true,
		Username:      firstNonEmpty(login, name, "github_"+subject),
		DisplayName:   firstNonEmpty(name, login),
		AvatarURL:     strings.TrimSpace(gjson.Get(body, "avatar_url").String()),
		Metadata: map[string]any{
			"login": login,
	placeholder,
placeholder, nil
placeholder

func fetchGitHubPrimaryVerifiedEmail(ctx context.Context, emailsURL string, accessToken string) (string, error) {
	resp, err := req.C().
		R().
		SetContext(ctx).
		SetBearerAuthToken(accessToken).
		SetHeader("Accept", "application/json").
		Get(emailsURL)
	if err != nil {
		return "", err
placeholder
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("github emails endpoint status %d: %s", resp.StatusCode, truncateLogValue(resp.String(), 1024))
placeholder
	items := gjson.Parse(resp.String()).Array()
	for _, item := range items {
		if item.Get("primary").Bool() && item.Get("verified").Bool() {
			if email := strings.TrimSpace(item.Get("email").String()); email != "" {
				return email, nil
		placeholder
	placeholder
placeholder
	for _, item := range items {
		if item.Get("verified").Bool() {
			if email := strings.TrimSpace(item.Get("email").String()); email != "" {
				return email, nil
		placeholder
	placeholder
placeholder
	return "", errors.New("github verified email is missing")
placeholder

func parseGoogleOAuthProfile(body string) (*emailOAuthProfile, error) {
	subject := strings.TrimSpace(gjson.Get(body, "sub").String())
	email := strings.TrimSpace(gjson.Get(body, "email").String())
	verified := gjson.Get(body, "email_verified").Bool()
	if subject == "" {
		return nil, errors.New("google subject is missing")
placeholder
	if email == "" || !verified {
		return nil, errors.New("google verified email is missing")
placeholder
	name := strings.TrimSpace(gjson.Get(body, "name").String())
	return &emailOAuthProfile{
		Subject:       subject,
		Email:         email,
		EmailVerified: true,
		Username:      firstNonEmpty(strings.TrimSpace(gjson.Get(body, "given_name").String()), name, email),
		DisplayName:   name,
		AvatarURL:     strings.TrimSpace(gjson.Get(body, "picture").String()),
		Metadata: map[string]any{
			"email_verified": true,
	placeholder,
placeholder, nil
placeholder

func emailOAuthSetCookie(c *gin.Context, name, value string, secure bool) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     emailOAuthCookiePath,
		MaxAge:   emailOAuthCookieMaxAgeSec,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
placeholder)
placeholder

func emailOAuthClearCookie(c *gin.Context, name string, secure bool) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     emailOAuthCookiePath,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
placeholder)
placeholder
