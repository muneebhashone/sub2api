package xai

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	SSOBuildScope        = "openid profile email offline_access grok-cli:access api:access conversations:read conversations:write"
	SSOAccountsURL       = "https://accounts.x.ai/"
	SSODeviceURL         = OAuthIssuer + "/oauth2/device/code"
	SSOVerifyURL         = OAuthIssuer + "/oauth2/device/verify"
	SSOApproveURL        = OAuthIssuer + "/oauth2/device/approve"
	SSOTokenURL          = OAuthIssuer + "/oauth2/token"
	SSOConversionTimeout = 90 * time.Second

	ssoMaxAuthBody     = 2 << 20
	ssoMaxTokenLength  = 16 << 10
	ssoDefaultUA       = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
	ssoDefaultTokenTTL = 6 * time.Hour
)

var (
	ErrSSOUnauthorized        = errors.New("xai sso unauthorized")
	ErrSSOAuthorizationDenied = errors.New("xai device authorization denied")
)

type SSOHTTPError struct{ Status int placeholder

func (e SSOHTTPError) Error() string { return fmt.Sprintf("xAI OAuth HTTP %d", e.Status) placeholder

type SSODeviceHTTPClient interface {
	Do(*http.Request) (*http.Response, error)
placeholder

type SSODeviceOptions struct {
	HTTPClient SSODeviceHTTPClient
	UserAgent  string
	Sleep      func(context.Context, time.Duration) error
placeholder

type ssoDeviceFlow struct {
	client    SSODeviceHTTPClient
	userAgent string
	cookieJar http.CookieJar
	sleep     func(context.Context, time.Duration) error
placeholder

func ConvertSSOToBuild(ctx context.Context, ssoToken string, opts *SSODeviceOptions) (*TokenResponse, error) {
	ssoToken = NormalizeSSOToken(ssoToken)
	if ssoToken == "" {
		return nil, ErrSSOUnauthorized
placeholder
	if opts == nil {
		opts = &SSODeviceOptions{placeholder
placeholder
	client := opts.HTTPClient
	if client == nil {
		client = &http.Client{
			Timeout: SSOConversionTimeout,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
		placeholder,
	placeholder
placeholder
	userAgent := strings.TrimSpace(opts.UserAgent)
	if userAgent == "" {
		userAgent = ssoDefaultUA
placeholder
	sleep := opts.Sleep
	if sleep == nil {
		sleep = sleepContext
placeholder
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
placeholder
	seedSSOCookies(jar, ssoToken)

	flow := &ssoDeviceFlow{
		client:    client,
		userAgent: userAgent,
		cookieJar: jar,
		sleep:     sleep,
placeholder
	return flow.convert(ctx)
placeholder

func (f *ssoDeviceFlow) convert(ctx context.Context) (*TokenResponse, error) {
	status, finalURL, _, err := f.do(ctx, http.MethodGet, SSOAccountsURL, nil)
	if err != nil {
		return nil, err
placeholder
	if status == http.StatusUnauthorized || strings.Contains(finalURL, "sign-in") || strings.Contains(finalURL, "sign-up") {
		return nil, ErrSSOUnauthorized
placeholder
	if status < 200 || status >= 400 {
		return nil, fmt.Errorf("validate Grok Web SSO: %w", SSOHTTPError{Status: statusplaceholder)
placeholder

	status, _, body, err := f.do(ctx, http.MethodPost, SSODeviceURL, url.Values{
		"client_id": {DefaultClientIDplaceholder,
		"scope":     {SSOBuildScopeplaceholder,
placeholder)
	if err != nil {
		return nil, err
placeholder
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("start xAI device flow: %w", SSOHTTPError{Status: statusplaceholder)
placeholder
	var device struct {
		DeviceCode              string `json:"device_code"`
		UserCode                string `json:"user_code"`
		VerificationURIComplete string `json:"verification_uri_complete"`
		Interval                int    `json:"interval"`
		ExpiresIn               int    `json:"expires_in"`
placeholder
	if err := json.Unmarshal(body, &device); err != nil {
		return nil, fmt.Errorf("parse xAI device flow response: %w", err)
placeholder
	if device.DeviceCode == "" || device.UserCode == "" || !safeXAIAuthURL(device.VerificationURIComplete) {
		return nil, errors.New("xAI device flow response is incomplete")
placeholder
	if device.Interval <= 0 {
		device.Interval = 5
placeholder
	if device.ExpiresIn <= 0 {
		device.ExpiresIn = 1800
placeholder

	status, _, _, err = f.do(ctx, http.MethodGet, device.VerificationURIComplete, nil)
	if err != nil {
		return nil, err
placeholder
	if status < 200 || status >= 400 {
		return nil, fmt.Errorf("open xAI device verification page: %w", SSOHTTPError{Status: statusplaceholder)
placeholder

	status, finalURL, _, err = f.do(ctx, http.MethodPost, SSOVerifyURL, url.Values{"user_code": {device.UserCodeplaceholderplaceholder)
	if err != nil {
		return nil, err
placeholder
	if status < 200 || status >= 400 {
		return nil, fmt.Errorf("verify xAI device code: %w", SSOHTTPError{Status: statusplaceholder)
placeholder
	if !strings.Contains(finalURL, "consent") {
		return nil, errors.New("xAI device verification did not reach consent page")
placeholder

	status, finalURL, _, err = f.do(ctx, http.MethodPost, SSOApproveURL, url.Values{
		"user_code":      {device.UserCodeplaceholder,
		"action":         {"allow"placeholder,
		"principal_type": {"User"placeholder,
		"principal_id":   {""placeholder,
placeholder)
	if err != nil {
		return nil, err
placeholder
	if status < 200 || status >= 400 {
		return nil, fmt.Errorf("approve xAI device code: %w", SSOHTTPError{Status: statusplaceholder)
placeholder
	if !strings.Contains(finalURL, "done") {
		return nil, errors.New("xAI device approval did not reach done page")
placeholder

	return f.pollToken(ctx, device.DeviceCode, time.Duration(device.Interval)*time.Second, time.Duration(device.ExpiresIn)*time.Second)
placeholder

func (f *ssoDeviceFlow) pollToken(ctx context.Context, deviceCode string, interval, expiresIn time.Duration) (*TokenResponse, error) {
	if interval < time.Second {
		interval = time.Second
placeholder
	deadline := time.Now().Add(minDuration(expiresIn, 75*time.Second))
	for time.Now().Before(deadline) {
		if err := f.sleep(ctx, interval); err != nil {
			return nil, err
	placeholder
		status, _, body, err := f.do(ctx, http.MethodPost, SSOTokenURL, url.Values{
			"grant_type":  {"urn:ietf:params:oauth:grant-type:device_code"placeholder,
			"client_id":   {DefaultClientIDplaceholder,
			"device_code": {deviceCodeplaceholder,
	placeholder)
		if err != nil {
			return nil, err
	placeholder
		var payload struct {
			AccessToken      string `json:"access_token"`
			RefreshToken     string `json:"refresh_token"`
			IDToken          string `json:"id_token"`
			TokenType        string `json:"token_type"`
			ExpiresIn        int64  `json:"expires_in"`
			Scope            string `json:"scope"`
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
	placeholder
		if err := json.Unmarshal(body, &payload); err != nil {
			return nil, fmt.Errorf("parse xAI token response: %w", err)
	placeholder
		if status >= 200 && status < 300 && payload.AccessToken != "" {
			if payload.ExpiresIn <= 0 {
				payload.ExpiresIn = int64(ssoDefaultTokenTTL.Seconds())
		placeholder
			if payload.TokenType == "" {
				payload.TokenType = "Bearer"
		placeholder
			return &TokenResponse{
				AccessToken:  payload.AccessToken,
				RefreshToken: payload.RefreshToken,
				IDToken:      payload.IDToken,
				TokenType:    payload.TokenType,
				ExpiresIn:    payload.ExpiresIn,
				Scope:        payload.Scope,
		placeholder, nil
	placeholder
		switch payload.Error {
		case "authorization_pending":
			continue
		case "slow_down":
			interval += 5 * time.Second
			continue
		case "access_denied", "expired_token":
			return nil, ErrSSOAuthorizationDenied
		default:
			if status >= 400 {
				return nil, fmt.Errorf("xAI token polling failed (%s): %w", firstNonEmpty(payload.ErrorDescription, payload.Error), SSOHTTPError{Status: statusplaceholder)
		placeholder
			return nil, fmt.Errorf("xAI token polling failed: %s", firstNonEmpty(payload.ErrorDescription, payload.Error, strconv.Itoa(status)))
	placeholder
placeholder
	return nil, errors.New("xAI device flow token polling timed out")
placeholder

func (f *ssoDeviceFlow) do(ctx context.Context, method, endpoint string, form url.Values) (int, string, []byte, error) {
	if !safeXAIAuthURL(endpoint) {
		return 0, "", nil, errors.New("xAI OAuth URL is not trusted")
placeholder
	currentURL := endpoint
	currentMethod := method
	currentForm := form
	for redirects := 0; redirects <= 8; redirects++ {
		var body io.Reader
		if currentForm != nil {
			body = strings.NewReader(currentForm.Encode())
	placeholder
		request, err := http.NewRequestWithContext(ctx, currentMethod, currentURL, body)
		if err != nil {
			return 0, currentURL, nil, err
	placeholder
		request.Header.Set("Accept", "application/json, text/html;q=0.9, */*;q=0.8")
		request.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
		request.Header.Set("User-Agent", f.userAgent)
		if cookie := f.cookieHeader(request.URL); cookie != "" {
			request.Header.Set("Cookie", cookie)
	placeholder
		if currentForm != nil {
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	placeholder

		response, err := f.client.Do(request)
		if err != nil {
			return 0, currentURL, nil, err
	placeholder
		f.captureCookies(request.URL, response)
		data, readErr := io.ReadAll(io.LimitReader(response.Body, ssoMaxAuthBody+1))
		_ = response.Body.Close()
		if readErr != nil {
			return response.StatusCode, currentURL, nil, readErr
	placeholder
		if len(data) > ssoMaxAuthBody {
			return response.StatusCode, currentURL, nil, errors.New("xAI OAuth response exceeds 2 MiB")
	placeholder
		if response.StatusCode < 300 || response.StatusCode > 399 {
			return response.StatusCode, currentURL, data, nil
	placeholder

		location := strings.TrimSpace(response.Header.Get("Location"))
		if location == "" {
			return response.StatusCode, currentURL, data, errors.New("xAI OAuth redirect missing Location")
	placeholder
		base, _ := url.Parse(currentURL)
		next, err := url.Parse(location)
		if err != nil {
			return response.StatusCode, currentURL, data, err
	placeholder
		currentURL = base.ResolveReference(next).String()
		if !safeXAIAuthURL(currentURL) {
			return response.StatusCode, currentURL, data, errors.New("xAI OAuth redirected to untrusted host")
	placeholder
		if response.StatusCode == http.StatusSeeOther || ((response.StatusCode == http.StatusMovedPermanently || response.StatusCode == http.StatusFound) && currentMethod != http.MethodGet && currentMethod != http.MethodHead) {
			currentMethod = http.MethodGet
			currentForm = nil
	placeholder
placeholder
	return 0, currentURL, nil, errors.New("xAI OAuth redirected too many times")
placeholder

func seedSSOCookies(jar http.CookieJar, token string) {
	if jar == nil {
		return
placeholder
	for _, rawURL := range []string{SSOAccountsURL, OAuthIssuer + "/"placeholder {
		target, err := url.Parse(rawURL)
		if err != nil {
			continue
	placeholder
		jar.SetCookies(target, []*http.Cookie{
			{Name: "sso", Value: token, Path: "/", Secure: true, HttpOnly: trueplaceholder,
			{Name: "sso-rw", Value: token, Path: "/", Secure: true, HttpOnly: trueplaceholder,
	placeholder)
placeholder
placeholder

func (f *ssoDeviceFlow) captureCookies(requestURL *url.URL, response *http.Response) {
	if f == nil || f.cookieJar == nil || requestURL == nil || response == nil {
		return
placeholder
	cookies := make([]*http.Cookie, 0)
	for _, cookie := range response.Cookies() {
		name := strings.TrimSpace(cookie.Name)
		value := strings.TrimSpace(cookie.Value)
		if name == "" || len(name) > 128 || len(value) > 16384 || strings.ContainsAny(name+value, "\r\n\x00") {
			continue
	placeholder
		cookie.Name = name
		cookie.Value = value
		cookies = append(cookies, cookie)
placeholder
	f.cookieJar.SetCookies(requestURL, cookies)
placeholder

func (f *ssoDeviceFlow) cookieHeader(requestURL *url.URL) string {
	if f == nil || f.cookieJar == nil || requestURL == nil {
		return ""
placeholder
	cookies := f.cookieJar.Cookies(requestURL)
	sort.Slice(cookies, func(i, j int) bool { return cookies[i].Name < cookies[j].Name placeholder)
	parts := make([]string, 0, len(cookies))
	for _, cookie := range cookies {
		parts = append(parts, cookie.Name+"="+cookie.Value)
placeholder
	return strings.Join(parts, "; ")
placeholder

func safeXAIAuthURL(raw string) bool {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.User != nil || parsed.Hostname() == "" {
		return false
placeholder
	if AllowUnsafeURLOverrides() {
		return parsed.Scheme != "" && parsed.Host != ""
placeholder
	if parsed.Scheme != "https" {
		return false
placeholder
	host := strings.ToLower(parsed.Hostname())
	return host == "x.ai" || strings.HasSuffix(host, ".x.ai")
placeholder

func NormalizeSSOToken(value string) string {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(strings.ToLower(value), "cookie:") {
		value = strings.TrimSpace(value[len("cookie:"):])
placeholder
	for _, part := range strings.Split(value, ";") {
		name, token, found := strings.Cut(strings.TrimSpace(part), "=")
		if !found {
			continue
	placeholder
		switch strings.ToLower(strings.TrimSpace(name)) {
		case "sso", "sso-rw":
			return sanitizeSSOToken(token)
	placeholder
placeholder
	if token, _, found := strings.Cut(value, ";"); found {
		value = strings.TrimSpace(token)
placeholder
	return sanitizeSSOToken(value)
placeholder

func sanitizeSSOToken(value string) string {
	value = strings.NewReplacer("\r", "", "\n", "", "\x00", "").Replace(strings.TrimSpace(value))
	if len(value) > ssoMaxTokenLength {
		return ""
placeholder
	return value
placeholder

func DecodeJWTClaims(token string) map[string]any {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil
placeholder
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil
placeholder
	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil
placeholder
	return claims
placeholder

func JWTClaimString(claims map[string]any, key string) string {
	value, _ := claims[key].(string)
	return strings.TrimSpace(value)
placeholder

func sleepContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
placeholder
placeholder

func minDuration(a, b time.Duration) time.Duration {
	if a <= 0 {
		return b
placeholder
	if a < b {
		return a
placeholder
	return b
placeholder

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
	placeholder
placeholder
	return ""
placeholder
