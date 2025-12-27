package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/oauth"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/imroc/req/v3"
)

func NewClaudeOAuthClient() service.ClaudeOAuthClient {
	return &claudeOAuthService{
		baseURL:       "https://claude.ai",
		tokenURL:      oauth.TokenURL,
		clientFactory: createReqClient,
placeholder
placeholder

type claudeOAuthService struct {
	baseURL       string
	tokenURL      string
	clientFactory func(proxyURL string) *req.Client
placeholder

func (s *claudeOAuthService) GetOrganizationUUID(ctx context.Context, sessionKey, proxyURL string) (string, error) {
	client := s.clientFactory(proxyURL)

	var orgs []struct {
		UUID string `json:"uuid"`
placeholder

	targetURL := s.baseURL + "/api/organizations"
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

func (s *claudeOAuthService) GetAuthorizationCode(ctx context.Context, sessionKey, orgUUID, scope, codeChallenge, state, proxyURL string) (string, error) {
	client := s.clientFactory(proxyURL)

	authURL := fmt.Sprintf("%s/v1/oauth/%s/authorize", s.baseURL, orgUUID)

	reqBody := map[string]any{
		"response_type":         "code",
		"client_id":             oauth.ClientID,
		"organization_uuid":     orgUUID,
		"redirect_uri":          oauth.RedirectURI,
		"scope":                 scope,
		"state":                 state,
		"code_challenge":        codeChallenge,
		"code_challenge_method": "S256",
placeholder

	reqBodyJSON, _ := json.Marshal(reqBody)
	log.Printf("[OAuth] Step 2: Getting authorization code from %s", authURL)
	log.Printf("[OAuth] Step 2 Request Body: %s", string(reqBodyJSON))

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

	fullCode := authCode
	if responseState != "" {
		fullCode = authCode + "#" + responseState
placeholder

	log.Printf("[OAuth] Step 2 SUCCESS - Got authorization code: %s...", prefix(authCode, 20))
	return fullCode, nil
placeholder

func (s *claudeOAuthService) ExchangeCodeForToken(ctx context.Context, code, codeVerifier, state, proxyURL string) (*oauth.TokenResponse, error) {
	client := s.clientFactory(proxyURL)

	// Parse code which may contain state in format "authCode#state"
	authCode := code
	codeState := ""
	if idx := strings.Index(code, "#"); idx != -1 {
		authCode = code[:idx]
		codeState = code[idx+1:]
placeholder

	reqBody := map[string]any{
		"code":          authCode,
		"grant_type":    "authorization_code",
		"client_id":     oauth.ClientID,
		"redirect_uri":  oauth.RedirectURI,
		"code_verifier": codeVerifier,
placeholder

	if codeState != "" {
		reqBody["state"] = codeState
placeholder

	reqBodyJSON, _ := json.Marshal(reqBody)
	log.Printf("[OAuth] Step 3: Exchanging code for token at %s", s.tokenURL)
	log.Printf("[OAuth] Step 3 Request Body: %s", string(reqBodyJSON))

	var tokenResp oauth.TokenResponse

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetSuccessResult(&tokenResp).
		Post(s.tokenURL)

	if err != nil {
		log.Printf("[OAuth] Step 3 FAILED - Request error: %v", err)
		return nil, fmt.Errorf("request failed: %w", err)
placeholder

	log.Printf("[OAuth] Step 3 Response - Status: %d, Body: %s", resp.StatusCode, resp.String())

	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("token exchange failed: status %d, body: %s", resp.StatusCode, resp.String())
placeholder

	log.Printf("[OAuth] Step 3 SUCCESS - Got access token")
	return &tokenResp, nil
placeholder

func (s *claudeOAuthService) RefreshToken(ctx context.Context, refreshToken, proxyURL string) (*oauth.TokenResponse, error) {
	client := s.clientFactory(proxyURL)

	// 使用 JSON 格式（与 ExchangeCodeForToken 保持一致）
	// Anthropic OAuth API 期望 JSON 格式的请求体
	reqBody := map[string]any{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"client_id":     oauth.ClientID,
placeholder

	var tokenResp oauth.TokenResponse

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetSuccessResult(&tokenResp).
		Post(s.tokenURL)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
placeholder

	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("token refresh failed: status %d, body: %s", resp.StatusCode, resp.String())
placeholder

	return &tokenResp, nil
placeholder

func createReqClient(proxyURL string) *req.Client {
	client := req.C().
		ImpersonateChrome().
		SetTimeout(60 * time.Second)

	if proxyURL != "" {
		client.SetProxyURL(proxyURL)
placeholder

	return client
placeholder

func prefix(s string, n int) string {
	if n <= 0 {
		return ""
placeholder
	if len(s) <= n {
		return s
placeholder
	return s[:n]
placeholder
