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
	State        string    `json:"state"`
	CodeVerifier string    `json:"code_verifier"`
	ProxyURL     string    `json:"proxy_url,omitempty"`
	RedirectURI  string    `json:"redirect_uri"`
	ProjectID    string    `json:"project_id,omitempty"`
	OAuthType    string    `json:"oauth_type"` // "code_assist" 或 "ai_studio"
	CreatedAt    time.Time `json:"created_at"`
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
// oauthType: "code_assist" or "ai_studio" (defaults to "code_assist" if empty)
// Returns error if ClientID or ClientSecret is not configured.
// Configure via GEMINI_OAUTH_CLIENT_ID and GEMINI_OAUTH_CLIENT_SECRET environment variables.
func EffectiveOAuthConfig(cfg OAuthConfig, oauthType string) (OAuthConfig, error) {
	effective := OAuthConfig{
		ClientID:     strings.TrimSpace(cfg.ClientID),
		ClientSecret: strings.TrimSpace(cfg.ClientSecret),
		Scopes:       strings.TrimSpace(cfg.Scopes),
placeholder

	// Require OAuth credentials to be configured
	if effective.ClientID == "" || effective.ClientSecret == "" {
		return OAuthConfig{placeholder, fmt.Errorf("gemini OAuth credentials not configured, set GEMINI_OAUTH_CLIENT_ID and GEMINI_OAUTH_CLIENT_SECRET environment variables")
placeholder

	if effective.Scopes == "" {
		// Use different default scopes based on OAuth type
		if oauthType == "ai_studio" {
			effective.Scopes = DefaultAIStudioScopes
	placeholder else {
			// Default to Code Assist scopes
			effective.Scopes = DefaultCodeAssistScopes
	placeholder
placeholder

	return effective, nil
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
