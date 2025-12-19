package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"sub2api/internal/model"
	"sub2api/internal/pkg/oauth"
	"sub2api/internal/service/ports"

	"github.com/imroc/req/v3"
)

// OAuthService handles OAuth authentication flows
type OAuthService struct {
	sessionStore *oauth.SessionStore
	proxyRepo    ports.ProxyRepository
placeholder

// NewOAuthService creates a new OAuth service
func NewOAuthService(proxyRepo ports.ProxyRepository) *OAuthService {
	return &OAuthService{
		sessionStore: oauth.NewSessionStore(),
		proxyRepo:    proxyRepo,
placeholder
placeholder

// GenerateAuthURLResult contains the authorization URL and session info
type GenerateAuthURLResult struct {
	AuthURL   string `json:"auth_url"`
	SessionID string `json:"session_id"`
placeholder

// GenerateAuthURL generates an OAuth authorization URL with full scope
func (s *OAuthService) GenerateAuthURL(ctx context.Context, proxyID *int64) (*GenerateAuthURLResult, error) {
	scope := fmt.Sprintf("%s %s", oauth.ScopeProfile, oauth.ScopeInference)
	return s.generateAuthURLWithScope(ctx, scope, proxyID)
placeholder

// GenerateSetupTokenURL generates an OAuth authorization URL for setup token (inference only)
func (s *OAuthService) GenerateSetupTokenURL(ctx context.Context, proxyID *int64) (*GenerateAuthURLResult, error) {
	scope := oauth.ScopeInference
	return s.generateAuthURLWithScope(ctx, scope, proxyID)
placeholder

func (s *OAuthService) generateAuthURLWithScope(ctx context.Context, scope string, proxyID *int64) (*GenerateAuthURLResult, error) {
	// Generate PKCE values
	state, err := oauth.GenerateState()
	if err != nil {
		return nil, fmt.Errorf("failed to generate state: %w", err)
placeholder

	codeVerifier, err := oauth.GenerateCodeVerifier()
	if err != nil {
		return nil, fmt.Errorf("failed to generate code verifier: %w", err)
placeholder

	codeChallenge := oauth.GenerateCodeChallenge(codeVerifier)

	// Generate session ID
	sessionID, err := oauth.GenerateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
placeholder

	// Get proxy URL if specified
	var proxyURL string
	if proxyID != nil {
		proxy, err := s.proxyRepo.GetByID(ctx, *proxyID)
		if err == nil && proxy != nil {
			proxyURL = proxy.URL()
	placeholder
placeholder

	// Store session
	session := &oauth.OAuthSession{
		State:        state,
		CodeVerifier: codeVerifier,
		Scope:        scope,
		ProxyURL:     proxyURL,
		CreatedAt:    time.Now(),
placeholder
	s.sessionStore.Set(sessionID, session)

	// Build authorization URL
	authURL := oauth.BuildAuthorizationURL(state, codeChallenge, scope)

	return &GenerateAuthURLResult{
		AuthURL:   authURL,
		SessionID: sessionID,
placeholder, nil
placeholder

// ExchangeCodeInput represents the input for code exchange
type ExchangeCodeInput struct {
	SessionID string
	Code      string
	ProxyID   *int64
placeholder

// TokenInfo represents the token information stored in credentials
type TokenInfo struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	ExpiresAt    int64  `json:"expires_at"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
	OrgUUID      string `json:"org_uuid,omitempty"`
	AccountUUID  string `json:"account_uuid,omitempty"`
placeholder

// ExchangeCode exchanges authorization code for tokens
func (s *OAuthService) ExchangeCode(ctx context.Context, input *ExchangeCodeInput) (*TokenInfo, error) {
	// Get session
	session, ok := s.sessionStore.Get(input.SessionID)
	if !ok {
		return nil, fmt.Errorf("session not found or expired")
placeholder

	// Get proxy URL
	proxyURL := session.ProxyURL
	if input.ProxyID != nil {
		proxy, err := s.proxyRepo.GetByID(ctx, *input.ProxyID)
		if err == nil && proxy != nil {
			proxyURL = proxy.URL()
	placeholder
placeholder

	// Exchange code for token
	tokenInfo, err := s.exchangeCodeForToken(ctx, input.Code, session.CodeVerifier, session.State, proxyURL)
	if err != nil {
		return nil, err
placeholder

	// Delete session after successful exchange
	s.sessionStore.Delete(input.SessionID)

	return tokenInfo, nil
placeholder

// CookieAuthInput represents the input for cookie-based authentication
type CookieAuthInput struct {
	SessionKey string
	ProxyID    *int64
	Scope      string // "full" or "inference"
placeholder

// CookieAuth performs OAuth using sessionKey (cookie-based auto-auth)
func (s *OAuthService) CookieAuth(ctx context.Context, input *CookieAuthInput) (*TokenInfo, error) {
	// Get proxy URL if specified
	var proxyURL string
	if input.ProxyID != nil {
		proxy, err := s.proxyRepo.GetByID(ctx, *input.ProxyID)
		if err == nil && proxy != nil {
			proxyURL = proxy.URL()
	placeholder
placeholder

	// Determine scope
	scope := fmt.Sprintf("%s %s", oauth.ScopeProfile, oauth.ScopeInference)
	if input.Scope == "inference" {
		scope = oauth.ScopeInference
placeholder

	// Step 1: Get organization info using sessionKey
	orgUUID, err := s.getOrganizationUUID(ctx, input.SessionKey, proxyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization info: %w", err)
placeholder

	// Step 2: Generate PKCE values
	codeVerifier, err := oauth.GenerateCodeVerifier()
	if err != nil {
		return nil, fmt.Errorf("failed to generate code verifier: %w", err)
placeholder
	codeChallenge := oauth.GenerateCodeChallenge(codeVerifier)

	state, err := oauth.GenerateState()
	if err != nil {
		return nil, fmt.Errorf("failed to generate state: %w", err)
placeholder

	// Step 3: Get authorization code using cookie
	authCode, err := s.getAuthorizationCode(ctx, input.SessionKey, orgUUID, scope, codeChallenge, state, proxyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get authorization code: %w", err)
placeholder

	// Step 4: Exchange code for token
	tokenInfo, err := s.exchangeCodeForToken(ctx, authCode, codeVerifier, state, proxyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
placeholder

	// Ensure org_uuid is set (from step 1 if not from token response)
	if tokenInfo.OrgUUID == "" && orgUUID != "" {
		tokenInfo.OrgUUID = orgUUID
		log.Printf("[OAuth] Set org_uuid from cookie auth: %s", orgUUID)
placeholder

	return tokenInfo, nil
placeholder

// getOrganizationUUID gets the organization UUID from claude.ai using sessionKey
func (s *OAuthService) getOrganizationUUID(ctx context.Context, sessionKey, proxyURL string) (string, error) {
	client := s.createReqClient(proxyURL)

	var orgs []struct {
		UUID string `json:"uuid"`
placeholder

	targetURL := "https://claude.ai/api/organizations"
	log.Printf("[OAuth] Step 1: Getting organization UUID from %s", targetURL)

	resp, err := client.R().
		SetContext(ctx).
		SetCookies(&http.Cookie{
			Name:  "sessionKey",
			Value: sessionKey,
	placeholder).
		SetSuccessResult(&orgs).
		Get(targetURL)

	if err != nil {
		log.Printf("[OAuth] Step 1 FAILED - Request error: %v", err)
		return "", fmt.Errorf("request failed: %w", err)
placeholder

	log.Printf("[OAuth] Step 1 Response - Status: %d, Body: %s", resp.StatusCode, resp.String())

	if !resp.IsSuccessState() {
		return "", fmt.Errorf("failed to get organizations: status %d, body: %s", resp.StatusCode, resp.String())
placeholder

	if len(orgs) == 0 {
		return "", fmt.Errorf("no organizations found")
placeholder

	log.Printf("[OAuth] Step 1 SUCCESS - Got org UUID: %s", orgs[0].UUID)
	return orgs[0].UUID, nil
placeholder

// getAuthorizationCode gets the authorization code using sessionKey
func (s *OAuthService) getAuthorizationCode(ctx context.Context, sessionKey, orgUUID, scope, codeChallenge, state, proxyURL string) (string, error) {
	client := s.createReqClient(proxyURL)

	authURL := fmt.Sprintf("https://claude.ai/v1/oauth/%s/authorize", orgUUID)

	// Build request body - must include organization_uuid as per CRS
	reqBody := map[string]interface{placeholder{
		"response_type":         "code",
		"client_id":             oauth.ClientID,
		"organization_uuid":     orgUUID, // Required field!
		"redirect_uri":          oauth.RedirectURI,
		"scope":                 scope,
		"state":                 state,
		"code_challenge":        codeChallenge,
		"code_challenge_method": "S256",
placeholder

	reqBodyJSON, _ := json.Marshal(reqBody)
	log.Printf("[OAuth] Step 2: Getting authorization code from %s", authURL)
	log.Printf("[OAuth] Step 2 Request Body: %s", string(reqBodyJSON))

	// Response contains redirect_uri with code, not direct code field
	var result struct {
		RedirectURI string `json:"redirect_uri"`
placeholder

	resp, err := client.R().
		SetContext(ctx).
		SetCookies(&http.Cookie{
			Name:  "sessionKey",
			Value: sessionKey,
	placeholder).
		SetHeader("Accept", "application/json").
		SetHeader("Accept-Language", "en-US,en;q=0.9").
		SetHeader("Cache-Control", "no-cache").
		SetHeader("Origin", "https://claude.ai").
		SetHeader("Referer", "https://claude.ai/new").
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetSuccessResult(&result).
		Post(authURL)

	if err != nil {
		log.Printf("[OAuth] Step 2 FAILED - Request error: %v", err)
		return "", fmt.Errorf("request failed: %w", err)
placeholder

	log.Printf("[OAuth] Step 2 Response - Status: %d, Body: %s", resp.StatusCode, resp.String())

	if !resp.IsSuccessState() {
		return "", fmt.Errorf("failed to get authorization code: status %d, body: %s", resp.StatusCode, resp.String())
placeholder

	if result.RedirectURI == "" {
		return "", fmt.Errorf("no redirect_uri in response")
placeholder

	// Parse redirect_uri to extract code and state
	parsedURL, err := url.Parse(result.RedirectURI)
	if err != nil {
		return "", fmt.Errorf("failed to parse redirect_uri: %w", err)
placeholder

	queryParams := parsedURL.Query()
	authCode := queryParams.Get("code")
	responseState := queryParams.Get("state")

	if authCode == "" {
		return "", fmt.Errorf("no authorization code in redirect_uri")
placeholder

	// Combine code with state if present (as CRS does)
	fullCode := authCode
	if responseState != "" {
		fullCode = authCode + "#" + responseState
placeholder

	log.Printf("[OAuth] Step 2 SUCCESS - Got authorization code: %s...", authCode[:20])
	return fullCode, nil
placeholder

// exchangeCodeForToken exchanges authorization code for tokens
func (s *OAuthService) exchangeCodeForToken(ctx context.Context, code, codeVerifier, state, proxyURL string) (*TokenInfo, error) {
	client := s.createReqClient(proxyURL)

	// Parse code#state format if present
	authCode := code
	codeState := ""
	if parts := strings.Split(code, "#"); len(parts) > 1 {
		authCode = parts[0]
		codeState = parts[1]
placeholder

	// Build JSON body as CRS does (not form data!)
	reqBody := map[string]interface{placeholder{
		"code":          authCode,
		"grant_type":    "authorization_code",
		"client_id":     oauth.ClientID,
		"redirect_uri":  oauth.RedirectURI,
		"code_verifier": codeVerifier,
placeholder

	// Add state if present
	if codeState != "" {
		reqBody["state"] = codeState
placeholder

	reqBodyJSON, _ := json.Marshal(reqBody)
	log.Printf("[OAuth] Step 3: Exchanging code for token at %s", oauth.TokenURL)
	log.Printf("[OAuth] Step 3 Request Body: %s", string(reqBodyJSON))

	var tokenResp oauth.TokenResponse

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetSuccessResult(&tokenResp).
		Post(oauth.TokenURL)

	if err != nil {
		log.Printf("[OAuth] Step 3 FAILED - Request error: %v", err)
		return nil, fmt.Errorf("request failed: %w", err)
placeholder

	log.Printf("[OAuth] Step 3 Response - Status: %d, Body: %s", resp.StatusCode, resp.String())

	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("token exchange failed: status %d, body: %s", resp.StatusCode, resp.String())
placeholder

	log.Printf("[OAuth] Step 3 SUCCESS - Got access token")

	tokenInfo := &TokenInfo{
		AccessToken:  tokenResp.AccessToken,
		TokenType:    tokenResp.TokenType,
		ExpiresIn:    tokenResp.ExpiresIn,
		ExpiresAt:    time.Now().Unix() + tokenResp.ExpiresIn,
		RefreshToken: tokenResp.RefreshToken,
		Scope:        tokenResp.Scope,
placeholder

	// Extract org_uuid and account_uuid from response
	if tokenResp.Organization != nil && tokenResp.Organization.UUID != "" {
		tokenInfo.OrgUUID = tokenResp.Organization.UUID
		log.Printf("[OAuth] Got org_uuid: %s", tokenInfo.OrgUUID)
placeholder
	if tokenResp.Account != nil && tokenResp.Account.UUID != "" {
		tokenInfo.AccountUUID = tokenResp.Account.UUID
		log.Printf("[OAuth] Got account_uuid: %s", tokenInfo.AccountUUID)
placeholder

	return tokenInfo, nil
placeholder

// RefreshToken refreshes an OAuth token
func (s *OAuthService) RefreshToken(ctx context.Context, refreshToken string, proxyURL string) (*TokenInfo, error) {
	client := s.createReqClient(proxyURL)

	formData := url.Values{placeholder
	formData.Set("grant_type", "refresh_token")
	formData.Set("refresh_token", refreshToken)
	formData.Set("client_id", oauth.ClientID)

	var tokenResp oauth.TokenResponse

	resp, err := client.R().
		SetContext(ctx).
		SetFormDataFromValues(formData).
		SetSuccessResult(&tokenResp).
		Post(oauth.TokenURL)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
placeholder

	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("token refresh failed: status %d, body: %s", resp.StatusCode, resp.String())
placeholder

	return &TokenInfo{
		AccessToken:  tokenResp.AccessToken,
		TokenType:    tokenResp.TokenType,
		ExpiresIn:    tokenResp.ExpiresIn,
		ExpiresAt:    time.Now().Unix() + tokenResp.ExpiresIn,
		RefreshToken: tokenResp.RefreshToken,
		Scope:        tokenResp.Scope,
placeholder, nil
placeholder

// RefreshAccountToken refreshes token for an account
func (s *OAuthService) RefreshAccountToken(ctx context.Context, account *model.Account) (*TokenInfo, error) {
	refreshToken := account.GetCredential("refresh_token")
	if refreshToken == "" {
		return nil, fmt.Errorf("no refresh token available")
placeholder

	var proxyURL string
	if account.ProxyID != nil {
		proxy, err := s.proxyRepo.GetByID(ctx, *account.ProxyID)
		if err == nil && proxy != nil {
			proxyURL = proxy.URL()
	placeholder
placeholder

	return s.RefreshToken(ctx, refreshToken, proxyURL)
placeholder

// createReqClient creates a req client with Chrome impersonation and optional proxy
func (s *OAuthService) createReqClient(proxyURL string) *req.Client {
	client := req.C().
		ImpersonateChrome(). // Impersonate Chrome browser to bypass Cloudflare
		SetTimeout(60 * time.Second)

	// Set proxy if specified
	if proxyURL != "" {
		client.SetProxyURL(proxyURL)
placeholder

	return client
placeholder
