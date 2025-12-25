package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/model"
	"github.com/Wei-Shaw/sub2api/internal/pkg/geminicli"
	"github.com/Wei-Shaw/sub2api/internal/service/ports"
)

type GeminiOAuthService struct {
	sessionStore *geminicli.SessionStore
	proxyRepo    ports.ProxyRepository
	oauthClient  ports.GeminiOAuthClient
	codeAssist   ports.GeminiCliCodeAssistClient
	cfg          *config.Config
placeholder

func NewGeminiOAuthService(
	proxyRepo ports.ProxyRepository,
	oauthClient ports.GeminiOAuthClient,
	codeAssist ports.GeminiCliCodeAssistClient,
	cfg *config.Config,
) *GeminiOAuthService {
	return &GeminiOAuthService{
		sessionStore: geminicli.NewSessionStore(),
		proxyRepo:    proxyRepo,
		oauthClient:  oauthClient,
		codeAssist:   codeAssist,
		cfg:          cfg,
placeholder
placeholder

type GeminiAuthURLResult struct {
	AuthURL   string `json:"auth_url"`
	SessionID string `json:"session_id"`
	State     string `json:"state"`
placeholder

func (s *GeminiOAuthService) GenerateAuthURL(ctx context.Context, proxyID *int64, redirectURI string) (*GeminiAuthURLResult, error) {
	state, err := geminicli.GenerateState()
	if err != nil {
		return nil, fmt.Errorf("failed to generate state: %w", err)
placeholder
	codeVerifier, err := geminicli.GenerateCodeVerifier()
	if err != nil {
		return nil, fmt.Errorf("failed to generate code verifier: %w", err)
placeholder
	codeChallenge := geminicli.GenerateCodeChallenge(codeVerifier)
	sessionID, err := geminicli.GenerateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
placeholder

	var proxyURL string
	if proxyID != nil {
		proxy, err := s.proxyRepo.GetByID(ctx, *proxyID)
		if err == nil && proxy != nil {
			proxyURL = proxy.URL()
	placeholder
placeholder

	session := &geminicli.OAuthSession{
		State:        state,
		CodeVerifier: codeVerifier,
		ProxyURL:     proxyURL,
		RedirectURI:  redirectURI,
		CreatedAt:    time.Now(),
placeholder
	s.sessionStore.Set(sessionID, session)

	oauthCfg := geminicli.OAuthConfig{
		ClientID:     s.cfg.Gemini.OAuth.ClientID,
		ClientSecret: s.cfg.Gemini.OAuth.ClientSecret,
		Scopes:       s.cfg.Gemini.OAuth.Scopes,
placeholder

	authURL, err := geminicli.BuildAuthorizationURL(oauthCfg, state, codeChallenge, redirectURI)
	if err != nil {
		return nil, err
placeholder

	return &GeminiAuthURLResult{
		AuthURL:   authURL,
		SessionID: sessionID,
		State:     state,
placeholder, nil
placeholder

type GeminiExchangeCodeInput struct {
	SessionID   string
	State       string
	Code        string
	RedirectURI string
	ProxyID     *int64
placeholder

type GeminiTokenInfo struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	ExpiresAt    int64  `json:"expires_at"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope,omitempty"`
	ProjectID    string `json:"project_id,omitempty"`
placeholder

func (s *GeminiOAuthService) ExchangeCode(ctx context.Context, input *GeminiExchangeCodeInput) (*GeminiTokenInfo, error) {
	session, ok := s.sessionStore.Get(input.SessionID)
	if !ok {
		return nil, fmt.Errorf("session not found or expired")
placeholder
	if strings.TrimSpace(input.State) == "" || input.State != session.State {
		return nil, fmt.Errorf("invalid state")
placeholder

	proxyURL := session.ProxyURL
	if input.ProxyID != nil {
		proxy, err := s.proxyRepo.GetByID(ctx, *input.ProxyID)
		if err == nil && proxy != nil {
			proxyURL = proxy.URL()
	placeholder
placeholder

	redirectURI := session.RedirectURI
	if strings.TrimSpace(input.RedirectURI) != "" {
		redirectURI = input.RedirectURI
placeholder

	tokenResp, err := s.oauthClient.ExchangeCode(ctx, input.Code, session.CodeVerifier, redirectURI, proxyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
placeholder
	s.sessionStore.Delete(input.SessionID)

	expiresAt := time.Now().Unix() + tokenResp.ExpiresIn
	projectID, _ := s.fetchProjectID(ctx, tokenResp.AccessToken, proxyURL)

	return &GeminiTokenInfo{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		TokenType:    tokenResp.TokenType,
		ExpiresIn:    tokenResp.ExpiresIn,
		ExpiresAt:    expiresAt,
		Scope:        tokenResp.Scope,
		ProjectID:    projectID,
placeholder, nil
placeholder

func (s *GeminiOAuthService) RefreshToken(ctx context.Context, refreshToken, proxyURL string) (*GeminiTokenInfo, error) {
	var lastErr error

	for attempt := 0; attempt <= 3; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			if backoff > 30*time.Second {
				backoff = 30 * time.Second
		placeholder
			time.Sleep(backoff)
	placeholder

		tokenResp, err := s.oauthClient.RefreshToken(ctx, refreshToken, proxyURL)
		if err == nil {
			expiresAt := time.Now().Unix() + tokenResp.ExpiresIn
			return &GeminiTokenInfo{
				AccessToken:  tokenResp.AccessToken,
				RefreshToken: tokenResp.RefreshToken,
				TokenType:    tokenResp.TokenType,
				ExpiresIn:    tokenResp.ExpiresIn,
				ExpiresAt:    expiresAt,
				Scope:        tokenResp.Scope,
		placeholder, nil
	placeholder

		if isNonRetryableGeminiOAuthError(err) {
			return nil, err
	placeholder
		lastErr = err
placeholder

	return nil, fmt.Errorf("token refresh failed after retries: %w", lastErr)
placeholder

func isNonRetryableGeminiOAuthError(err error) bool {
	msg := err.Error()
	nonRetryable := []string{
		"invalid_grant",
		"invalid_client",
		"unauthorized_client",
		"access_denied",
placeholder
	for _, needle := range nonRetryable {
		if strings.Contains(msg, needle) {
			return true
	placeholder
placeholder
	return false
placeholder

func (s *GeminiOAuthService) RefreshAccountToken(ctx context.Context, account *model.Account) (*GeminiTokenInfo, error) {
	if account.Platform != model.PlatformGemini || account.Type != model.AccountTypeOAuth {
		return nil, fmt.Errorf("account is not a Gemini OAuth account")
placeholder

	refreshToken := account.GetCredential("refresh_token")
	if strings.TrimSpace(refreshToken) == "" {
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

func (s *GeminiOAuthService) BuildAccountCredentials(tokenInfo *GeminiTokenInfo) map[string]any {
	creds := map[string]any{
		"access_token": tokenInfo.AccessToken,
		"expires_at":   strconv.FormatInt(tokenInfo.ExpiresAt, 10),
placeholder
	if tokenInfo.RefreshToken != "" {
		creds["refresh_token"] = tokenInfo.RefreshToken
placeholder
	if tokenInfo.TokenType != "" {
		creds["token_type"] = tokenInfo.TokenType
placeholder
	if tokenInfo.Scope != "" {
		creds["scope"] = tokenInfo.Scope
placeholder
	if tokenInfo.ProjectID != "" {
		creds["project_id"] = tokenInfo.ProjectID
placeholder
	return creds
placeholder

func (s *GeminiOAuthService) Stop() {
	s.sessionStore.Stop()
placeholder

func (s *GeminiOAuthService) fetchProjectID(ctx context.Context, accessToken, proxyURL string) (string, error) {
	if s.codeAssist == nil {
		return "", errors.New("code assist client not configured")
placeholder

	loadResp, err := s.codeAssist.LoadCodeAssist(ctx, accessToken, proxyURL, nil)
	if err == nil && strings.TrimSpace(loadResp.CurrentTier) != "" && strings.TrimSpace(loadResp.CloudAICompanionProject) != "" {
		return strings.TrimSpace(loadResp.CloudAICompanionProject), nil
placeholder

	// pick default tier from allowedTiers, fallback to LEGACY.
	tierID := "LEGACY"
	if loadResp != nil {
		for _, tier := range loadResp.AllowedTiers {
			if tier.IsDefault && strings.TrimSpace(tier.ID) != "" {
				tierID = tier.ID
				break
		placeholder
	placeholder
placeholder

	req := &geminicli.OnboardUserRequest{
		TierID: tierID,
		Metadata: geminicli.LoadCodeAssistMetadata{
			IDEType:    "ANTIGRAVITY",
			Platform:   "PLATFORM_UNSPECIFIED",
			PluginType: "GEMINI",
	placeholder,
placeholder

	maxAttempts := 5
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resp, err := s.codeAssist.OnboardUser(ctx, accessToken, proxyURL, req)
		if err != nil {
			return "", err
	placeholder
		if resp.Done {
			if resp.Response == nil || resp.Response.CloudAICompanionProject == nil {
				return "", errors.New("onboardUser completed but no project_id returned")
		placeholder
			switch v := resp.Response.CloudAICompanionProject.(type) {
			case string:
				return strings.TrimSpace(v), nil
			case map[string]any:
				if id, ok := v["id"].(string); ok {
					return strings.TrimSpace(id), nil
			placeholder
		placeholder
			return "", errors.New("onboardUser returned unsupported project_id format")
	placeholder
		time.Sleep(2 * time.Second)
placeholder

	return "", fmt.Errorf("onboardUser timeout after %d attempts", maxAttempts)
placeholder
