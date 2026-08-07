package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	sharedhttp "github.com/Wei-Shaw/sub2api/internal/pkg/httpclient"
	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/Wei-Shaw/sub2api/internal/util/logredact"
	"github.com/imroc/req/v3"
)

type grokOAuthClient struct {
	tokenURL string
placeholder

const (
	accountsBaseURL      = "https://accounts.x.ai"
	loginRPCEndpoint     = accountsBaseURL + "/api/rpc"
	turnstileWebsiteURL  = accountsBaseURL
	turnstileWebsiteKey  = "0x4AAAAAAAhr9JGVDZbrZOo0"
	yesCaptchaCreateTask = "https://api.yescaptcha.com/createTask"
	yesCaptchaGetResult  = "https://api.yescaptcha.com/getTaskResult"
)

func NewGrokOAuthClient() service.GrokOAuthClient {
	// Fail closed: never fall back to an unvalidated EffectiveTokenURL (env can
	// point at an attacker host and steal code/refresh tokens).
	tokenURL, err := xai.ValidatedTokenURL()
	if err != nil || strings.TrimSpace(tokenURL) == "" {
		// Official allowlisted endpoint only — never EffectiveTokenURL() (raw env).
		tokenURL = xai.DefaultTokenURL
placeholder
	return &grokOAuthClient{tokenURL: tokenURLplaceholder
placeholder

func (c *grokOAuthClient) ExchangeCode(ctx context.Context, code, codeVerifier, redirectURI, proxyURL, clientID string) (*xai.TokenResponse, error) {
	client, err := createGrokReqClient(proxyURL)
	if err != nil {
		return nil, infraerrors.Newf(http.StatusBadGateway, "GROK_OAUTH_CLIENT_INIT_FAILED", "create HTTP client: %v", err)
placeholder

	clientID = strings.TrimSpace(clientID)
	if clientID == "" {
		clientID = xai.EffectiveClientID()
placeholder

	formData := url.Values{placeholder
	formData.Set("grant_type", "authorization_code")
	formData.Set("client_id", clientID)
	formData.Set("code", code)
	formData.Set("redirect_uri", xai.EffectiveRedirectURI(redirectURI))
	formData.Set("code_verifier", codeVerifier)

	var tokenResp xai.TokenResponse
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("User-Agent", "sub2api-grok-oauth/1.0").
		SetFormDataFromValues(formData).
		SetSuccessResult(&tokenResp).
		Post(c.tokenURL)
	if err != nil {
		return nil, infraerrors.Newf(http.StatusBadGateway, "GROK_OAUTH_REQUEST_FAILED", "request failed: %v", err)
placeholder
	if !resp.IsSuccessState() {
		return nil, grokOAuthStatusError("GROK_OAUTH_TOKEN_EXCHANGE_FAILED", "token exchange failed", resp)
placeholder
	return &tokenResp, nil
placeholder

func (c *grokOAuthClient) RefreshToken(ctx context.Context, refreshToken, proxyURL, clientID string) (*xai.TokenResponse, error) {
	client, err := createGrokReqClient(proxyURL)
	if err != nil {
		return nil, infraerrors.Newf(http.StatusBadGateway, "GROK_OAUTH_CLIENT_INIT_FAILED", "create HTTP client: %v", err)
placeholder

	clientID = strings.TrimSpace(clientID)
	if clientID == "" {
		clientID = xai.EffectiveClientID()
placeholder

	formData := url.Values{placeholder
	formData.Set("grant_type", "refresh_token")
	formData.Set("client_id", clientID)
	formData.Set("refresh_token", refreshToken)

	var tokenResp xai.TokenResponse
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("User-Agent", "sub2api-grok-oauth/1.0").
		SetFormDataFromValues(formData).
		SetSuccessResult(&tokenResp).
		Post(c.tokenURL)
	if err != nil {
		return nil, infraerrors.Newf(http.StatusBadGateway, "GROK_OAUTH_REQUEST_FAILED", "request failed: %v", err)
placeholder
	if !resp.IsSuccessState() {
		return nil, grokOAuthStatusError("GROK_OAUTH_TOKEN_REFRESH_FAILED", "token refresh failed", resp)
placeholder
	return &tokenResp, nil
placeholder

// LoginWithPassword authenticates against accounts.x.ai and returns an ephemeral SSO cookie.
// Password and SSO must never be written to account credentials or logs.
func (c *grokOAuthClient) LoginWithPassword(ctx context.Context, email, password, proxyURL string) (*service.GrokPasswordLoginResult, error) {
	turnstileToken, err := solveTurnstile(ctx)
	if err != nil {
		return nil, err
placeholder
	httpClient, err := createGrokHTTPClient(proxyURL, true)
	if err != nil {
		return nil, infraerrors.Newf(http.StatusBadGateway, "GROK_OAUTH_CLIENT_INIT_FAILED", "create HTTP client: %v", err)
placeholder
	cookieSetterURL, err := createGrokPasswordSession(ctx, httpClient, strings.TrimSpace(email), password, turnstileToken)
	if err != nil {
		return nil, err
placeholder
	ssoToken, err := extractGrokSSOToken(ctx, httpClient, cookieSetterURL)
	if err != nil {
		return nil, err
placeholder
	return &service.GrokPasswordLoginResult{
		Email:    strings.TrimSpace(email),
		SSOToken: ssoToken,
placeholder, nil
placeholder

func (c *grokOAuthClient) ConvertSSOToBuild(ctx context.Context, ssoToken, proxyURL string) (*xai.TokenResponse, error) {
	client, err := createGrokSSOHTTPClient(proxyURL)
	if err != nil {
		return nil, infraerrors.Newf(http.StatusBadGateway, "GROK_SSO_CLIENT_INIT_FAILED", "create HTTP client: %v", err)
placeholder

	requestCtx, cancel := context.WithTimeout(ctx, xai.SSOConversionTimeout)
	defer cancel()
	tokenResp, err := xai.ConvertSSOToBuild(requestCtx, ssoToken, &xai.SSODeviceOptions{HTTPClient: clientplaceholder)
	if err != nil {
		return nil, grokSSOConversionError(err)
placeholder
	return tokenResp, nil
placeholder

func createGrokReqClient(proxyURL string) (*req.Client, error) {
	return getSharedReqClient(reqClientOptions{
		ProxyURL: proxyURL,
		Timeout:  60 * time.Second,
placeholder)
placeholder

func createGrokSSOHTTPClient(proxyURL string) (*http.Client, error) {
	client, err := sharedhttp.GetClient(sharedhttp.Options{
		ProxyURL:              proxyURL,
		Timeout:               xai.SSOConversionTimeout,
		ResponseHeaderTimeout: 30 * time.Second,
placeholder)
	if err != nil {
		return nil, err
placeholder
	clone := *client
	clone.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
placeholder
	return &clone, nil
placeholder

func grokSSOConversionError(err error) error {
	if errors.Is(err, xai.ErrSSOUnauthorized) {
		return infraerrors.New(http.StatusUnauthorized, "GROK_SSO_UNAUTHORIZED", "Grok Web SSO cookie is invalid or expired")
placeholder
	if errors.Is(err, xai.ErrSSOAuthorizationDenied) {
		return infraerrors.New(http.StatusForbidden, "GROK_SSO_AUTHORIZATION_DENIED", "xAI device authorization was denied or expired")
placeholder
	var statusErr xai.SSOHTTPError
	if errors.As(err, &statusErr) {
		statusCode := http.StatusBadGateway
		if statusErr.Status == http.StatusForbidden {
			statusCode = http.StatusForbidden
	placeholder
		return infraerrors.Newf(statusCode, "GROK_SSO_UPSTREAM_FAILED", "xAI SSO conversion failed: %v", err)
placeholder
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return infraerrors.Newf(http.StatusGatewayTimeout, "GROK_SSO_TIMEOUT", "xAI SSO conversion timed out: %v", err)
placeholder
	return infraerrors.Newf(http.StatusBadGateway, "GROK_SSO_CONVERSION_FAILED", "xAI SSO conversion failed: %v", err)
placeholder

func grokOAuthStatusError(code, message string, resp *req.Response) error {
	statusCode := http.StatusBadGateway
	errorCode := code
	upstreamStatus := 0
	body := ""
	if resp != nil {
		upstreamStatus = resp.StatusCode
		body = logredact.RedactText(resp.String())
		if resp.StatusCode == http.StatusForbidden && grokOAuthHasExplicitEntitlementDenial(body) {
			statusCode = http.StatusForbidden
			errorCode = "GROK_OAUTH_ENTITLEMENT_DENIED"
	placeholder
placeholder
	return infraerrors.Newf(statusCode, errorCode, "%s: status %d, body: %s", message, upstreamStatus, body)
placeholder

func grokOAuthHasExplicitEntitlementDenial(body string) bool {
	lower := strings.ToLower(body)
	compact := strings.NewReplacer(" ", "", "\n", "", "\r", "", "\t", "").Replace(lower)
	for _, field := range []string{"error", "code", "reason"placeholder {
		for _, value := range []string{"access_denied", "entitlement_denied", "subscription_required", "no_active_subscription"placeholder {
			if strings.Contains(compact, `"`+field+`":"`+value+`"`) {
				return true
		placeholder
	placeholder
placeholder
	return strings.Contains(lower, "entitlement denied") ||
		strings.Contains(lower, "subscription required") ||
		strings.Contains(lower, "no active grok subscription")
placeholder

func createGrokHTTPClient(proxyURL string, noRedirect bool) (*http.Client, error) {
	transport := &http.Transport{placeholder
	if strings.TrimSpace(proxyURL) != "" {
		parsed, err := url.Parse(proxyURL)
		if err != nil {
			return nil, err
	placeholder
		transport.Proxy = http.ProxyURL(parsed)
placeholder
	client := &http.Client{Timeout: 120 * time.Second, Transport: transportplaceholder
	if noRedirect {
		client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
	placeholder
placeholder
	return client, nil
placeholder

func solveTurnstile(ctx context.Context) (string, error) {
	clientKey := strings.TrimSpace(os.Getenv("YESCAPTCHA_CLIENT_KEY"))
	if clientKey == "" {
		clientKey = strings.TrimSpace(os.Getenv("YESCAPTCHA_API_KEY"))
placeholder
	if clientKey == "" {
		return "", infraerrors.New(http.StatusBadRequest, "GROK_OAUTH_CAPTCHA_KEY_REQUIRED", "yescaptcha client key is required for Grok password authorization")
placeholder
	createBody, err := json.Marshal(map[string]any{
		"clientKey": clientKey,
		"task": map[string]any{
			"type":       "TurnstileTaskProxyless",
			"websiteURL": turnstileWebsiteURL,
			"websiteKey": turnstileWebsiteKey,
	placeholder,
placeholder)
	if err != nil {
		return "", infraerrors.Newf(http.StatusInternalServerError, "GROK_OAUTH_CAPTCHA_FAILED", "encode captcha create request failed: %v", err)
placeholder
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, yesCaptchaCreateTask, bytes.NewReader(createBody))
	if err != nil {
		return "", infraerrors.Newf(http.StatusBadGateway, "GROK_OAUTH_CAPTCHA_FAILED", "build captcha create request failed: %v", err)
placeholder
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", infraerrors.Newf(http.StatusBadGateway, "GROK_OAUTH_CAPTCHA_FAILED", "create captcha task failed: %v", err)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	var createResp struct {
		ErrorID          int    `json:"errorId"`
		TaskID           string `json:"taskId"`
		ErrorDescription string `json:"errorDescription"`
placeholder
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return "", infraerrors.Newf(http.StatusBadGateway, "GROK_OAUTH_CAPTCHA_FAILED", "decode captcha create response failed: %v", err)
placeholder
	if createResp.ErrorID != 0 || strings.TrimSpace(createResp.TaskID) == "" {
		return "", infraerrors.Newf(http.StatusBadGateway, "GROK_OAUTH_CAPTCHA_FAILED", "captcha create failed: %s", createResp.ErrorDescription)
placeholder
	deadline := time.Now().Add(90 * time.Second)
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(5 * time.Second):
	placeholder
		body, err := json.Marshal(map[string]any{"clientKey": clientKey, "taskId": createResp.TaskIDplaceholder)
		if err != nil {
			return "", infraerrors.Newf(http.StatusInternalServerError, "GROK_OAUTH_CAPTCHA_FAILED", "encode captcha poll request failed: %v", err)
	placeholder
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, yesCaptchaGetResult, bytes.NewReader(body))
		if err != nil {
			continue
	placeholder
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
	placeholder
		var pollResp struct {
			ErrorID          int    `json:"errorId"`
			Status           string `json:"status"`
			ErrorDescription string `json:"errorDescription"`
			Solution         struct {
				Token string `json:"token"`
		placeholder `json:"solution"`
	placeholder
		err = json.NewDecoder(resp.Body).Decode(&pollResp)
		_ = resp.Body.Close()
		if err != nil {
			continue
	placeholder
		if pollResp.ErrorID != 0 {
			return "", infraerrors.Newf(http.StatusBadGateway, "GROK_OAUTH_CAPTCHA_FAILED", "captcha poll failed: %s", pollResp.ErrorDescription)
	placeholder
		if pollResp.Status == "ready" && strings.TrimSpace(pollResp.Solution.Token) != "" {
			return pollResp.Solution.Token, nil
	placeholder
placeholder
	return "", infraerrors.New(http.StatusGatewayTimeout, "GROK_OAUTH_CAPTCHA_TIMEOUT", "captcha solve timed out")
placeholder

func createGrokPasswordSession(ctx context.Context, client *http.Client, email, password, turnstileToken string) (string, error) {
	payload, err := json.Marshal(map[string]any{
		"rpc": "createSession",
		"req": map[string]any{
			"createSessionRequest": map[string]any{
				"credentials": map[string]any{
					"case": "emailAndPassword",
					"value": map[string]any{
						"email":             email,
						"clearTextPassword": password,
				placeholder,
			placeholder,
		placeholder,
			"turnstileToken": turnstileToken,
	placeholder,
placeholder)
	if err != nil {
		return "", infraerrors.Newf(http.StatusInternalServerError, "GROK_OAUTH_PASSWORD_LOGIN_FAILED", "encode password login request failed: %v", err)
placeholder
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, loginRPCEndpoint, bytes.NewReader(payload))
	if err != nil {
		return "", infraerrors.Newf(http.StatusBadGateway, "GROK_OAUTH_PASSWORD_LOGIN_FAILED", "build password login request failed: %v", err)
placeholder
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", accountsBaseURL)
	req.Header.Set("Referer", accountsBaseURL+"/sign-in?redirect=grok-com&email=true")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "*/*")
	resp, err := client.Do(req)
	if err != nil {
		return "", infraerrors.Newf(http.StatusBadGateway, "GROK_OAUTH_PASSWORD_LOGIN_FAILED", "password login request failed: %v", err)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", infraerrors.Newf(http.StatusBadGateway, "GROK_OAUTH_PASSWORD_LOGIN_FAILED", "password login returned status %d: %s", resp.StatusCode, logredact.RedactText(string(body)))
placeholder
	var loginResp struct {
		CookieSetterURL string `json:"cookieSetterUrl"`
		Error           string `json:"error"`
placeholder
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return "", infraerrors.Newf(http.StatusBadGateway, "GROK_OAUTH_PASSWORD_LOGIN_FAILED", "decode password login response failed: %v", err)
placeholder
	if strings.TrimSpace(loginResp.Error) != "" {
		return "", infraerrors.Newf(http.StatusBadGateway, "GROK_OAUTH_PASSWORD_LOGIN_FAILED", "password login error: %s", logredact.RedactText(loginResp.Error))
placeholder
	if strings.TrimSpace(loginResp.CookieSetterURL) == "" {
		return "", infraerrors.New(http.StatusBadGateway, "GROK_OAUTH_PASSWORD_LOGIN_FAILED", "password login did not return cookieSetterUrl")
placeholder
	return loginResp.CookieSetterURL, nil
placeholder

func extractGrokSSOToken(ctx context.Context, client *http.Client, cookieSetterURL string) (string, error) {
	safeURL, err := validateGrokCookieSetterURL(cookieSetterURL)
	if err != nil {
		return "", infraerrors.Newf(http.StatusBadGateway, "GROK_OAUTH_PASSWORD_LOGIN_FAILED", "invalid cookie setter url: %v", err)
placeholder
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, safeURL.String(), nil)
	if err != nil {
		return "", infraerrors.Newf(http.StatusBadGateway, "GROK_OAUTH_PASSWORD_LOGIN_FAILED", "build cookie setter request: %v", err)
placeholder
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Referer", accountsBaseURL+"/")
	resp, err := client.Do(req)
	if err != nil {
		return "", infraerrors.Newf(http.StatusBadGateway, "GROK_OAUTH_PASSWORD_LOGIN_FAILED", "follow cookie setter url failed: %v", err)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	for _, cookie := range resp.Header.Values("Set-Cookie") {
		if token, ok := strings.CutPrefix(cookie, "sso="); ok {
			if idx := strings.Index(token, ";"); idx > 0 {
				token = token[:idx]
		placeholder
			return strings.TrimSpace(token), nil
	placeholder
placeholder
	return "", infraerrors.Newf(http.StatusBadGateway, "GROK_OAUTH_PASSWORD_LOGIN_FAILED", "no sso cookie found in response (status=%d)", resp.StatusCode)
placeholder

func validateGrokCookieSetterURL(rawURL string) (*url.URL, error) {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return nil, err
placeholder
	if parsed.Scheme != "https" || !strings.EqualFold(parsed.Hostname(), "accounts.x.ai") {
		return nil, fmt.Errorf("url must use https://accounts.x.ai")
placeholder
	if parsed.User != nil || parsed.Port() != "" || parsed.Fragment != "" || parsed.Opaque != "" {
		return nil, fmt.Errorf("url contains disallowed authority or fragment components")
placeholder
	return parsed, nil
placeholder
