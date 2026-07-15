package xai

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/util/logredact"
	"github.com/Wei-Shaw/sub2api/internal/util/urlvalidator"
)

const (
	OAuthIssuer         = "https://auth.x.ai"
	DiscoveryURL        = OAuthIssuer + "/.well-known/openid-configuration"
	DefaultAuthorizeURL = OAuthIssuer + "/oauth2/authorize"
	DefaultTokenURL     = OAuthIssuer + "/oauth2/token"
	DefaultBaseURL      = "https://api.x.ai/v1"
	DefaultCLIBaseURL   = "https://cli-chat-proxy.grok.com/v1"
	DefaultClientID     = "b1a00492-073a-47ea-816f-4c329264a828"
	DefaultScope        = "openid profile email offline_access grok-cli:access api:access"
	DefaultRedirectURI  = "http://127.0.0.1:56121/callback"
	SessionTTL          = 30 * time.Minute

	EnvAuthorizeURL               = "XAI_OAUTH_AUTHORIZE_URL"
	EnvTokenURL                   = "XAI_OAUTH_TOKEN_URL"
	EnvClientID                   = "XAI_OAUTH_CLIENT_ID"
	EnvScope                      = "XAI_OAUTH_SCOPE"
	EnvRedirectURI                = "XAI_OAUTH_REDIRECT_URI"
	EnvBaseURL                    = "XAI_BASE_URL"
	EnvAllowUnsafeURLOverrides    = "XAI_ALLOW_UNSAFE_URL_OVERRIDES"
	EnvUnsafeAllowHighConcurrency = "XAI_GROK_UNSAFE_ALLOW_CONCURRENCY_GT_ONE"
)

var (
	oauthEndpointAllowedHosts = []string{"x.ai", "*.x.ai"placeholder
	baseURLAllowedHosts       = []string{"api.x.ai", "cli-chat-proxy.grok.com"placeholder
)

// OAuthSession stores one PKCE OAuth flow.
type OAuthSession struct {
	State         string    `json:"state"`
	CodeVerifier  string    `json:"code_verifier"`
	CodeChallenge string    `json:"code_challenge"`
	ClientID      string    `json:"client_id,omitempty"`
	Scope         string    `json:"scope,omitempty"`
	ProxyURL      string    `json:"proxy_url,omitempty"`
	RedirectURI   string    `json:"redirect_uri"`
	CreatedAt     time.Time `json:"created_at"`
placeholder

// SessionStore manages xAI OAuth sessions in memory.
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*OAuthSession
	stopOnce sync.Once
	stopCh   chan struct{placeholder
placeholder

func NewSessionStore() *SessionStore {
	store := &SessionStore{
		sessions: make(map[string]*OAuthSession),
		stopCh:   make(chan struct{placeholder),
placeholder
	go store.cleanup()
	return store
placeholder

func (s *SessionStore) Set(sessionID string, session *OAuthSession) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[sessionID] = session
placeholder

func (s *SessionStore) Get(sessionID string) (*OAuthSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessions[sessionID]
	if !ok {
		return nil, false
placeholder
	if time.Since(session.CreatedAt) > SessionTTL {
		return nil, false
placeholder
	return session, true
placeholder

func (s *SessionStore) Delete(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, sessionID)
placeholder

func (s *SessionStore) Stop() {
	s.stopOnce.Do(func() {
		close(s.stopCh)
placeholder)
placeholder

func (s *SessionStore) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.mu.Lock()
			for id, session := range s.sessions {
				if time.Since(session.CreatedAt) > SessionTTL {
					delete(s.sessions, id)
			placeholder
		placeholder
			s.mu.Unlock()
	placeholder
placeholder
placeholder

func EffectiveAuthorizeURL() string {
	return envOrDefault(EnvAuthorizeURL, DefaultAuthorizeURL)
placeholder

func ValidatedAuthorizeURL() (string, error) {
	return ValidateOAuthEndpointURL(EffectiveAuthorizeURL())
placeholder

func EffectiveTokenURL() string {
	return envOrDefault(EnvTokenURL, DefaultTokenURL)
placeholder

func ValidatedTokenURL() (string, error) {
	return ValidateOAuthEndpointURL(EffectiveTokenURL())
placeholder

func EffectiveClientID() string {
	return envOrDefault(EnvClientID, DefaultClientID)
placeholder

func EffectiveScope() string {
	return envOrDefault(EnvScope, DefaultScope)
placeholder

func EffectiveRedirectURI(override string) string {
	if trimmed := strings.TrimSpace(override); trimmed != "" {
		return trimmed
placeholder
	return envOrDefault(EnvRedirectURI, DefaultRedirectURI)
placeholder

func EffectiveBaseURL(override string) string {
	if trimmed := strings.TrimSpace(override); trimmed != "" {
		return strings.TrimRight(trimmed, "/")
placeholder
	return strings.TrimRight(envOrDefault(EnvBaseURL, DefaultBaseURL), "/")
placeholder

func ValidatedBaseURL(override string) (string, error) {
	return ValidateBaseURL(EffectiveBaseURL(override))
placeholder

// BaseURLValidator applies the caller's outbound URL trust policy before xAI
// endpoint paths are appended. The service layer uses this for API-key accounts
// so the global security.url_allowlist policy remains the single source of
// truth; OAuth callers keep using the strict trusted-host validator.
type BaseURLValidator func(string) (string, error)

func validatedBaseURLWithValidator(override string, validator BaseURLValidator) (string, error) {
	if validator == nil {
		return ValidatedBaseURL(override)
placeholder
	raw := EffectiveBaseURL(override)
	validated, err := validator(raw)
	if err != nil {
		return "", err
placeholder
	return normalizeKnownBaseURLPath(validated)
placeholder

type RuntimeSanityCheck struct {
	Value     string `json:"value"`
	Valid     bool   `json:"valid"`
	Error     string `json:"error,omitempty"`
	IsDefault bool   `json:"is_default,omitempty"`
placeholder

type RuntimeSanityReport struct {
	BaseURL               RuntimeSanityCheck `json:"base_url"`
	OAuthAuthorizeURL     RuntimeSanityCheck `json:"oauth_authorize_url"`
	OAuthTokenURL         RuntimeSanityCheck `json:"oauth_token_url"`
	OAuthRedirectURI      RuntimeSanityCheck `json:"oauth_redirect_uri"`
	UnsafeURLOverrides    bool               `json:"unsafe_url_overrides"`
	UnsafeHighConcurrency bool               `json:"unsafe_high_concurrency"`
	PublicGatewayScope    string             `json:"public_gateway_scope"`
	ProxyPolicy           string             `json:"proxy_policy"`
placeholder

func RuntimeSanity() RuntimeSanityReport {
	return RuntimeSanityReport{
		BaseURL:               runtimeSanityCheck(EffectiveBaseURL(""), EnvBaseURL, ValidatedBaseURL),
		OAuthAuthorizeURL:     runtimeSanityCheck(EffectiveAuthorizeURL(), EnvAuthorizeURL, func(string) (string, error) { return ValidatedAuthorizeURL() placeholder),
		OAuthTokenURL:         runtimeSanityCheck(EffectiveTokenURL(), EnvTokenURL, func(string) (string, error) { return ValidatedTokenURL() placeholder),
		OAuthRedirectURI:      runtimeSanityCheck(EffectiveRedirectURI(""), EnvRedirectURI, validateRedirectURI),
		UnsafeURLOverrides:    AllowUnsafeURLOverrides(),
		UnsafeHighConcurrency: AllowUnsafeHighConcurrency(),
		PublicGatewayScope:    "responses_only",
		ProxyPolicy:           "account_proxy_optional; OAuth URLs use trusted-host allowlists; API-key base URLs require public HTTPS unless unsafe overrides are enabled",
placeholder
placeholder

func runtimeSanityCheck(value string, envKey string, validate func(string) (string, error)) RuntimeSanityCheck {
	normalized, err := validate(value)
	check := RuntimeSanityCheck{
		Value:     sanitizeRuntimeURLValue(normalized),
		Valid:     err == nil,
		IsDefault: strings.TrimSpace(os.Getenv(envKey)) == "",
placeholder
	if err != nil {
		check.Value = sanitizeRuntimeURLValue(value)
		check.Error = sanitizeRuntimeError(err.Error(), value)
placeholder
	return check
placeholder

func validateRedirectURI(raw string) (string, error) {
	return urlvalidator.ValidateURLFormat(raw, true)
placeholder

func sanitizeRuntimeURLValue(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
placeholder
	parsed, err := url.Parse(trimmed)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return trimmed
placeholder
	parsed.User = nil
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return strings.TrimRight(parsed.String(), "/")
placeholder

func sanitizeRuntimeError(rawErr string, rawValue string) string {
	redacted := logredact.RedactText(rawErr)
	trimmedValue := strings.TrimSpace(rawValue)
	if trimmedValue == "" {
		return redacted
placeholder
	sanitizedValue := sanitizeRuntimeURLValue(trimmedValue)
	redacted = strings.ReplaceAll(redacted, trimmedValue, sanitizedValue)
	redacted = strings.ReplaceAll(redacted, logredact.RedactText(trimmedValue), sanitizedValue)
	return redacted
placeholder

func ValidateOAuthEndpointURL(raw string) (string, error) {
	if AllowUnsafeURLOverrides() {
		return urlvalidator.ValidateURLFormat(raw, true)
placeholder
	return urlvalidator.ValidateHTTPSURL(raw, urlvalidator.ValidationOptions{
		AllowedHosts:     oauthEndpointAllowedHosts,
		RequireAllowlist: true,
		AllowPrivate:     false,
placeholder)
placeholder

func ValidateBaseURL(raw string) (string, error) {
	if AllowUnsafeURLOverrides() {
		return urlvalidator.ValidateURLFormat(raw, true)
placeholder
	normalized, err := urlvalidator.ValidateHTTPSURL(raw, urlvalidator.ValidationOptions{
		AllowPrivate: false,
placeholder)
	if err != nil {
		return "", err
placeholder
	return normalizeKnownBaseURLPath(normalized)
placeholder

func ValidateTrustedBaseURL(raw string) (string, error) {
	if AllowUnsafeURLOverrides() {
		return urlvalidator.ValidateURLFormat(raw, true)
placeholder
	normalized, err := urlvalidator.ValidateHTTPSURL(raw, urlvalidator.ValidationOptions{
		AllowedHosts:     baseURLAllowedHosts,
		RequireAllowlist: true,
		AllowPrivate:     false,
placeholder)
	if err != nil {
		return "", err
placeholder
	return normalizeKnownBaseURLPath(normalized)
placeholder

func normalizeKnownBaseURLPath(raw string) (string, error) {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", errors.New("invalid base URL")
placeholder
	if parsed.User != nil {
		return "", errors.New("base URL must not include userinfo")
placeholder
	if parsed.RawQuery != "" {
		return "", errors.New("base URL must not include a query")
placeholder
	if parsed.Fragment != "" {
		return "", errors.New("base URL must not include a fragment")
placeholder
	path := strings.TrimRight(parsed.Path, "/")
	if path == "" {
		parsed.Path = "/v1"
		parsed.RawPath = ""
		return strings.TrimRight(parsed.String(), "/"), nil
placeholder
	if path != "/v1" {
		return "", fmt.Errorf("base URL path must be /v1")
placeholder
	parsed.Path = path
	parsed.RawPath = ""
	return strings.TrimRight(parsed.String(), "/"), nil
placeholder

func AllowUnsafeURLOverrides() bool {
	return envBool(EnvAllowUnsafeURLOverrides)
placeholder

func AllowUnsafeHighConcurrency() bool {
	return envBool(EnvUnsafeAllowHighConcurrency)
placeholder

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
placeholder
	return fallback
placeholder

func envBool(key string) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(key))) {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
placeholder
placeholder

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
placeholder
	return b, nil
placeholder

func GenerateState() (string, error) {
	bytes, err := GenerateRandomBytes(32)
	if err != nil {
		return "", err
placeholder
	return hex.EncodeToString(bytes), nil
placeholder

func GenerateNonce() (string, error) {
	bytes, err := GenerateRandomBytes(16)
	if err != nil {
		return "", err
placeholder
	return hex.EncodeToString(bytes), nil
placeholder

func GenerateSessionID() (string, error) {
	bytes, err := GenerateRandomBytes(16)
	if err != nil {
		return "", err
placeholder
	return hex.EncodeToString(bytes), nil
placeholder

func GenerateCodeVerifier() (string, error) {
	bytes, err := GenerateRandomBytes(32)
	if err != nil {
		return "", err
placeholder
	return base64URLEncode(bytes), nil
placeholder

func GenerateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64URLEncode(hash[:])
placeholder

func base64URLEncode(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
placeholder

func BuildAuthorizationURL(state, codeChallenge, redirectURI, nonce string) (string, error) {
	redirectURI = EffectiveRedirectURI(redirectURI)
	authorizeURL, err := ValidatedAuthorizeURL()
	if err != nil {
		return "", fmt.Errorf("invalid authorize url: %w", err)
placeholder

	params := url.Values{placeholder
	params.Set("response_type", "code")
	params.Set("client_id", EffectiveClientID())
	params.Set("redirect_uri", redirectURI)
	params.Set("scope", EffectiveScope())
	params.Set("state", state)
	params.Set("nonce", nonce)
	params.Set("code_challenge", codeChallenge)
	params.Set("code_challenge_method", "S256")
	params.Set("plan", "generic")
	params.Set("referrer", "sub2api")

	return fmt.Sprintf("%s?%s", authorizeURL, params.Encode()), nil
placeholder

// AuthorizationInput is a parsed manual OAuth callback input.
type AuthorizationInput struct {
	Code          string
	State         string
	RequiresState bool
placeholder

// ParseAuthorizationInput accepts a full callback URL, query string, or bare code.
func ParseAuthorizationInput(raw string) AuthorizationInput {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return AuthorizationInput{placeholder
placeholder

	if parsed, err := url.Parse(trimmed); err == nil && parsed != nil {
		values := parsed.Query()
		if code := strings.TrimSpace(values.Get("code")); code != "" {
			return AuthorizationInput{
				Code:          code,
				State:         strings.TrimSpace(values.Get("state")),
				RequiresState: true,
		placeholder
	placeholder
placeholder

	queryCandidate := strings.TrimPrefix(trimmed, "?")
	if strings.Contains(queryCandidate, "=") {
		if values, err := url.ParseQuery(queryCandidate); err == nil {
			if code := strings.TrimSpace(values.Get("code")); code != "" {
				return AuthorizationInput{
					Code:          code,
					State:         strings.TrimSpace(values.Get("state")),
					RequiresState: true,
			placeholder
		placeholder
	placeholder
placeholder

	return AuthorizationInput{Code: trimmedplaceholder
placeholder

func BuildResponsesURL(baseURL string) (string, error) {
	return BuildResponsesURLWithValidator(baseURL, nil)
placeholder

func BuildResponsesURLWithValidator(baseURL string, validator BaseURLValidator) (string, error) {
	validatedBaseURL, err := validatedBaseURLWithValidator(baseURL, validator)
	if err != nil {
		return "", fmt.Errorf("invalid base url: %w", err)
placeholder
	return validatedBaseURL + "/responses", nil
placeholder

func BuildChatCompletionsURL(baseURL string) (string, error) {
	return BuildChatCompletionsURLWithValidator(baseURL, nil)
placeholder

func BuildChatCompletionsURLWithValidator(baseURL string, validator BaseURLValidator) (string, error) {
	validatedBaseURL, err := validatedBaseURLWithValidator(baseURL, validator)
	if err != nil {
		return "", fmt.Errorf("invalid base url: %w", err)
placeholder
	return validatedBaseURL + "/chat/completions", nil
placeholder

func BuildImagesGenerationsURL(baseURL string) (string, error) {
	return BuildImagesGenerationsURLWithValidator(baseURL, nil)
placeholder

func BuildImagesGenerationsURLWithValidator(baseURL string, validator BaseURLValidator) (string, error) {
	validatedBaseURL, err := validatedBaseURLWithValidator(baseURL, validator)
	if err != nil {
		return "", fmt.Errorf("invalid base url: %w", err)
placeholder
	return validatedBaseURL + "/images/generations", nil
placeholder

func BuildImagesEditsURL(baseURL string) (string, error) {
	return BuildImagesEditsURLWithValidator(baseURL, nil)
placeholder

func BuildImagesEditsURLWithValidator(baseURL string, validator BaseURLValidator) (string, error) {
	validatedBaseURL, err := validatedBaseURLWithValidator(baseURL, validator)
	if err != nil {
		return "", fmt.Errorf("invalid base url: %w", err)
placeholder
	return validatedBaseURL + "/images/edits", nil
placeholder

func BuildVideosGenerationsURL(baseURL string) (string, error) {
	return BuildVideosGenerationsURLWithValidator(baseURL, nil)
placeholder

func BuildVideosGenerationsURLWithValidator(baseURL string, validator BaseURLValidator) (string, error) {
	validatedBaseURL, err := validatedBaseURLWithValidator(baseURL, validator)
	if err != nil {
		return "", fmt.Errorf("invalid base url: %w", err)
placeholder
	return validatedBaseURL + "/videos/generations", nil
placeholder

func BuildVideosEditsURL(baseURL string) (string, error) {
	return BuildVideosEditsURLWithValidator(baseURL, nil)
placeholder

func BuildVideosEditsURLWithValidator(baseURL string, validator BaseURLValidator) (string, error) {
	validatedBaseURL, err := validatedBaseURLWithValidator(baseURL, validator)
	if err != nil {
		return "", fmt.Errorf("invalid base url: %w", err)
placeholder
	return validatedBaseURL + "/videos/edits", nil
placeholder

func BuildVideosExtensionsURL(baseURL string) (string, error) {
	return BuildVideosExtensionsURLWithValidator(baseURL, nil)
placeholder

func BuildVideosExtensionsURLWithValidator(baseURL string, validator BaseURLValidator) (string, error) {
	validatedBaseURL, err := validatedBaseURLWithValidator(baseURL, validator)
	if err != nil {
		return "", fmt.Errorf("invalid base url: %w", err)
placeholder
	return validatedBaseURL + "/videos/extensions", nil
placeholder

func BuildVideoURL(baseURL, requestID string) (string, error) {
	return BuildVideoURLWithValidator(baseURL, requestID, nil)
placeholder

func BuildVideoURLWithValidator(baseURL, requestID string, validator BaseURLValidator) (string, error) {
	validatedBaseURL, err := validatedBaseURLWithValidator(baseURL, validator)
	if err != nil {
		return "", fmt.Errorf("invalid base url: %w", err)
placeholder
	requestID = strings.TrimSpace(requestID)
	if requestID == "" {
		return "", fmt.Errorf("request id is required")
placeholder
	return validatedBaseURL + "/videos/" + url.PathEscape(requestID), nil
placeholder

// TokenResponse represents xAI OAuth token responses.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	Scope        string `json:"scope,omitempty"`
placeholder
