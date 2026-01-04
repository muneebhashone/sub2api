package service

import (
	"context"
	"net/url"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/geminicli"
)

func TestGeminiOAuthService_GenerateAuthURL_RedirectURIStrategy(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name          string
		cfg           *config.Config
		oauthType     string
		projectID     string
		wantClientID  string
		wantRedirect  string
		wantScope     string
		wantProjectID string
		wantErrSubstr string
placeholder

	tests := []testCase{
		{
			name: "google_one uses built-in client when not configured and redirects to upstream",
			cfg: &config.Config{
				Gemini: config.GeminiConfig{
					OAuth: config.GeminiOAuthConfig{placeholder,
			placeholder,
		placeholder,
			oauthType:     "google_one",
			wantClientID:  geminicli.GeminiCLIOAuthClientID,
			wantRedirect:  geminicli.GeminiCLIRedirectURI,
			wantScope:     geminicli.DefaultCodeAssistScopes,
			wantProjectID: "",
	placeholder,
		{
			name: "google_one uses custom client when configured and redirects to localhost",
			cfg: &config.Config{
				Gemini: config.GeminiConfig{
					OAuth: config.GeminiOAuthConfig{
						ClientID:     "custom-client-id",
						ClientSecret: "custom-client-secret",
				placeholder,
			placeholder,
		placeholder,
			oauthType:     "google_one",
			wantClientID:  "custom-client-id",
			wantRedirect:  geminicli.AIStudioOAuthRedirectURI,
			wantScope:     geminicli.DefaultGoogleOneScopes,
			wantProjectID: "",
	placeholder,
		{
			name: "code_assist always forces built-in client even when custom client configured",
			cfg: &config.Config{
				Gemini: config.GeminiConfig{
					OAuth: config.GeminiOAuthConfig{
						ClientID:     "custom-client-id",
						ClientSecret: "custom-client-secret",
				placeholder,
			placeholder,
		placeholder,
			oauthType:     "code_assist",
			projectID:     "my-gcp-project",
			wantClientID:  geminicli.GeminiCLIOAuthClientID,
			wantRedirect:  geminicli.GeminiCLIRedirectURI,
			wantScope:     geminicli.DefaultCodeAssistScopes,
			wantProjectID: "my-gcp-project",
	placeholder,
		{
			name: "ai_studio requires custom client",
			cfg: &config.Config{
				Gemini: config.GeminiConfig{
					OAuth: config.GeminiOAuthConfig{placeholder,
			placeholder,
		placeholder,
			oauthType:     "ai_studio",
			wantErrSubstr: "AI Studio OAuth requires a custom OAuth Client",
	placeholder,
placeholder

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := NewGeminiOAuthService(nil, nil, nil, tt.cfg)
			got, err := svc.GenerateAuthURL(context.Background(), nil, "https://example.com/auth/callback", tt.projectID, tt.oauthType, "")
			if tt.wantErrSubstr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErrSubstr)
			placeholder
				if !strings.Contains(err.Error(), tt.wantErrSubstr) {
					t.Fatalf("expected error containing %q, got: %v", tt.wantErrSubstr, err)
			placeholder
				return
		placeholder
			if err != nil {
				t.Fatalf("GenerateAuthURL returned error: %v", err)
		placeholder

			parsed, err := url.Parse(got.AuthURL)
			if err != nil {
				t.Fatalf("failed to parse auth_url: %v", err)
		placeholder
			q := parsed.Query()

			if gotState := q.Get("state"); gotState != got.State {
				t.Fatalf("state mismatch: query=%q result=%q", gotState, got.State)
		placeholder
			if gotClientID := q.Get("client_id"); gotClientID != tt.wantClientID {
				t.Fatalf("client_id mismatch: got=%q want=%q", gotClientID, tt.wantClientID)
		placeholder
			if gotRedirect := q.Get("redirect_uri"); gotRedirect != tt.wantRedirect {
				t.Fatalf("redirect_uri mismatch: got=%q want=%q", gotRedirect, tt.wantRedirect)
		placeholder
			if gotScope := q.Get("scope"); gotScope != tt.wantScope {
				t.Fatalf("scope mismatch: got=%q want=%q", gotScope, tt.wantScope)
		placeholder
			if gotProjectID := q.Get("project_id"); gotProjectID != tt.wantProjectID {
				t.Fatalf("project_id mismatch: got=%q want=%q", gotProjectID, tt.wantProjectID)
		placeholder
	placeholder)
placeholder
placeholder
