package geminicli

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"
)

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	Scopes       string
placeholder

type OAuthSession struct {
	State        string `json:"state"`
	CodeVerifier string `json:"code_verifier"`
	ProxyURL     string `json:"proxy_url,omitempty"`
	RedirectURI  string `json:"redirect_uri"`
	ProjectID    string `json:"project_id,omitempty"`
	// TierID is a user-selected fallback tier.
	// For oauth types that support auto detection (google_one/code_assist), the server will prefer
	// the detected tier and fall back to TierID when detection fails.
	TierID    string    `json:"tier_id,omitempty"`
	OAuthType string    `json:"oauth_type"` // "code_assist" 或 "ai_studio"
	CreatedAt time.Time `json:"created_at"`
placeholder

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*OAuthSession
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
	select {
	case <-s.stopCh:
		return
	default:
		close(s.stopCh)
placeholder
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

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
placeholder
	return b, nil
placeholder

func GenerateState() (string, error) {
	bytes, err := GenerateRandomBytes(32)
	if err != nil {
		return "", err
placeholder
	return base64URLEncode(bytes), nil
placeholder

func GenerateSessionID() (string, error) {
	bytes, err := GenerateRandomBytes(16)
	if err != nil {
		return "", err
placeholder
	return hex.EncodeToString(bytes), nil
placeholder

// GenerateCodeVerifier returns an RFC 7636 compatible code verifier (43+ chars).
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

// EffectiveOAuthConfig returns the effective OAuth configuration.
// oauthType: "code_assist" or "ai_studio" (defaults to "code_assist" if empty).
//
// If ClientID/ClientSecret is not provided, this falls back to the built-in Gemini CLI OAuth client.
//
// Note: The built-in Gemini CLI OAuth client is restricted and may reject some scopes (e.g.
// https://www.googleapis.com/auth/generative-language), which will surface as
// "restricted_client" / "Unregistered scope(s)" errors during browser authorization.
func EffectiveOAuthConfig(cfg OAuthConfig, oauthType string) (OAuthConfig, error) {
	effective := OAuthConfig{
		ClientID:     strings.TrimSpace(cfg.ClientID),
		ClientSecret: strings.TrimSpace(cfg.ClientSecret),
		Scopes:       strings.TrimSpace(cfg.Scopes),
placeholder

	// Normalize scopes: allow comma-separated input but send space-delimited scopes to Google.
	if effective.Scopes != "" {
		effective.Scopes = strings.Join(strings.Fields(strings.ReplaceAll(effective.Scopes, ",", " ")), " ")
placeholder

	// Fall back to built-in Gemini CLI OAuth client when not configured.
	if effective.ClientID == "" && effective.ClientSecret == "" {
		effective.ClientID = GeminiCLIOAuthClientID
		effective.ClientSecret = GeminiCLIOAuthClientSecret
placeholder else if effective.ClientID == "" || effective.ClientSecret == "" {
		return OAuthConfig{placeholder, fmt.Errorf("OAuth client not configured: please set both client_id and client_secret (or leave both empty to use the built-in Gemini CLI client)")
placeholder

	isBuiltinClient := effective.ClientID == GeminiCLIOAuthClientID &&
		effective.ClientSecret == GeminiCLIOAuthClientSecret

	if effective.Scopes == "" {
		// Use different default scopes based on OAuth type
		switch oauthType {
		case "ai_studio":
			// Built-in client can't request some AI Studio scopes (notably generative-language).
			if isBuiltinClient {
				effective.Scopes = DefaultCodeAssistScopes
		placeholder else {
				effective.Scopes = DefaultAIStudioScopes
		placeholder
		case "google_one":
			// Google One uses built-in Gemini CLI client (same as code_assist)
			// Built-in client can't request restricted scopes like generative-language.retriever
			if isBuiltinClient {
				effective.Scopes = DefaultCodeAssistScopes
		placeholder else {
				effective.Scopes = DefaultGoogleOneScopes
		placeholder
		default:
			// Default to Code Assist scopes
			effective.Scopes = DefaultCodeAssistScopes
	placeholder
placeholder else if (oauthType == "ai_studio" || oauthType == "google_one") && isBuiltinClient {
		// If user overrides scopes while still using the built-in client, strip restricted scopes.
		parts := strings.Fields(effective.Scopes)
		filtered := make([]string, 0, len(parts))
		for _, s := range parts {
			if hasRestrictedScope(s) {
				continue
		placeholder
			filtered = append(filtered, s)
	placeholder
		if len(filtered) == 0 {
			effective.Scopes = DefaultCodeAssistScopes
	placeholder else {
			effective.Scopes = strings.Join(filtered, " ")
	placeholder
placeholder

	// Backward compatibility: normalize older AI Studio scope to the currently documented one.
	if oauthType == "ai_studio" && effective.Scopes != "" {
		parts := strings.Fields(effective.Scopes)
		for i := range parts {
			if parts[i] == "https://www.googleapis.com/auth/generative-language" {
				parts[i] = "https://www.googleapis.com/auth/generative-language.retriever"
		placeholder
	placeholder
		effective.Scopes = strings.Join(parts, " ")
placeholder

	return effective, nil
placeholder

func hasRestrictedScope(scope string) bool {
	return strings.HasPrefix(scope, "https://www.googleapis.com/auth/generative-language") ||
		strings.HasPrefix(scope, "https://www.googleapis.com/auth/drive")
placeholder

func BuildAuthorizationURL(cfg OAuthConfig, state, codeChallenge, redirectURI, projectID, oauthType string) (string, error) {
	effectiveCfg, err := EffectiveOAuthConfig(cfg, oauthType)
	if err != nil {
		return "", err
placeholder
	redirectURI = strings.TrimSpace(redirectURI)
	if redirectURI == "" {
		return "", fmt.Errorf("redirect_uri is required")
placeholder

	params := url.Values{placeholder
	params.Set("response_type", "code")
	params.Set("client_id", effectiveCfg.ClientID)
	params.Set("redirect_uri", redirectURI)
	params.Set("scope", effectiveCfg.Scopes)
	params.Set("state", state)
	params.Set("code_challenge", codeChallenge)
	params.Set("code_challenge_method", "S256")
	params.Set("access_type", "offline")
	params.Set("prompt", "consent")
	params.Set("include_granted_scopes", "true")
	if strings.TrimSpace(projectID) != "" {
		params.Set("project_id", strings.TrimSpace(projectID))
placeholder

	return fmt.Sprintf("%s?%s", AuthorizeURL, params.Encode()), nil
placeholder
