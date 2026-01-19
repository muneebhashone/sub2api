package repository

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/oauth"
	"github.com/imroc/req/v3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ClaudeOAuthServiceSuite struct {
	suite.Suite
	client *claudeOAuthService
placeholder

// requestCapture holds captured request data for assertions in the main goroutine.
type requestCapture struct {
	path        string
	method      string
	cookies     []*http.Cookie
	body        []byte
	bodyJSON    map[string]any
	contentType string
placeholder

func newTestReqClient(rt http.RoundTripper) *req.Client {
	c := req.C()
	c.GetClient().Transport = rt
	return c
placeholder

func (s *ClaudeOAuthServiceSuite) TestGetOrganizationUUID() {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantErr    bool
		errContain string
		wantUUID   string
		validate   func(captured requestCapture)
placeholder{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`[{"uuid":"org-1"placeholder]`))
		placeholder,
			wantUUID: "org-1",
			validate: func(captured requestCapture) {
				require.Equal(s.T(), "/api/organizations", captured.path, "unexpected path")
				require.Len(s.T(), captured.cookies, 1, "expected 1 cookie")
				require.Equal(s.T(), "sessionKey", captured.cookies[0].Name)
				require.Equal(s.T(), "sess", captured.cookies[0].Value)
		placeholder,
	placeholder,
		{
			name: "non_200_returns_error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("unauthorized"))
		placeholder,
			wantErr:    true,
			errContain: "401",
	placeholder,
		{
			name: "invalid_json_returns_error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte("not-json"))
		placeholder,
			wantErr: true,
	placeholder,
placeholder

	for _, tt := range tests {
		s.Run(tt.name, func() {
			var captured requestCapture

			rt := newInProcessTransport(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				captured.path = r.URL.Path
				captured.cookies = r.Cookies()
				tt.handler(w, r)
		placeholder), nil)

			client, ok := NewClaudeOAuthClient().(*claudeOAuthService)
			require.True(s.T(), ok, "type assertion failed")
			s.client = client
			s.client.baseURL = "http://in-process"
			s.client.clientFactory = func(string) *req.Client { return newTestReqClient(rt) placeholder

			got, err := s.client.GetOrganizationUUID(context.Background(), "sess", "")

			if tt.wantErr {
				require.Error(s.T(), err)
				if tt.errContain != "" {
					require.ErrorContains(s.T(), err, tt.errContain)
			placeholder
				return
		placeholder

			require.NoError(s.T(), err)
			require.Equal(s.T(), tt.wantUUID, got)
			if tt.validate != nil {
				tt.validate(captured)
		placeholder
	placeholder)
placeholder
placeholder

func (s *ClaudeOAuthServiceSuite) TestGetAuthorizationCode() {
	tests := []struct {
		name     string
		handler  http.HandlerFunc
		wantErr  bool
		wantCode string
		validate func(captured requestCapture)
placeholder{
		{
			name: "parses_redirect_uri",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]string{
					"redirect_uri": oauth.RedirectURI + "?code=AUTH&state=STATE",
			placeholder)
		placeholder,
			wantCode: "AUTH#STATE",
			validate: func(captured requestCapture) {
				require.True(s.T(), strings.HasPrefix(captured.path, "/v1/oauth/") && strings.HasSuffix(captured.path, "/authorize"), "unexpected path: %s", captured.path)
				require.Equal(s.T(), http.MethodPost, captured.method, "expected POST")
				require.Len(s.T(), captured.cookies, 1, "expected 1 cookie")
				require.Equal(s.T(), "sess", captured.cookies[0].Value)
				require.Equal(s.T(), "org-1", captured.bodyJSON["organization_uuid"])
				require.Equal(s.T(), oauth.ClientID, captured.bodyJSON["client_id"])
				require.Equal(s.T(), oauth.RedirectURI, captured.bodyJSON["redirect_uri"])
				require.Equal(s.T(), "st", captured.bodyJSON["state"])
		placeholder,
	placeholder,
		{
			name: "missing_code_returns_error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]string{
					"redirect_uri": oauth.RedirectURI + "?state=STATE", // no code
			placeholder)
		placeholder,
			wantErr: true,
	placeholder,
placeholder

	for _, tt := range tests {
		s.Run(tt.name, func() {
			var captured requestCapture

			rt := newInProcessTransport(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				captured.path = r.URL.Path
				captured.method = r.Method
				captured.cookies = r.Cookies()
				captured.body, _ = io.ReadAll(r.Body)
				_ = json.Unmarshal(captured.body, &captured.bodyJSON)
				tt.handler(w, r)
		placeholder), nil)

			client, ok := NewClaudeOAuthClient().(*claudeOAuthService)
			require.True(s.T(), ok, "type assertion failed")
			s.client = client
			s.client.baseURL = "http://in-process"
			s.client.clientFactory = func(string) *req.Client { return newTestReqClient(rt) placeholder

			code, err := s.client.GetAuthorizationCode(context.Background(), "sess", "org-1", oauth.ScopeInference, "cc", "st", "")

			if tt.wantErr {
				require.Error(s.T(), err)
				return
		placeholder

			require.NoError(s.T(), err)
			require.Equal(s.T(), tt.wantCode, code)
			if tt.validate != nil {
				tt.validate(captured)
		placeholder
	placeholder)
placeholder
placeholder

func (s *ClaudeOAuthServiceSuite) TestExchangeCodeForToken() {
	tests := []struct {
		name         string
		handler      http.HandlerFunc
		code         string
		isSetupToken bool
		wantErr      bool
		wantResp     *oauth.TokenResponse
		validate     func(captured requestCapture)
placeholder{
		{
			name: "sends_state_when_embedded",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(oauth.TokenResponse{
					AccessToken:  "at",
					TokenType:    "bearer",
					ExpiresIn:    3600,
					RefreshToken: "rt",
					Scope:        "s",
			placeholder)
		placeholder,
			code:         "AUTH#STATE2",
			isSetupToken: false,
			wantResp: &oauth.TokenResponse{
				AccessToken:  "at",
				RefreshToken: "rt",
		placeholder,
			validate: func(captured requestCapture) {
				require.Equal(s.T(), http.MethodPost, captured.method, "expected POST")
				require.True(s.T(), strings.HasPrefix(captured.contentType, "application/json"), "unexpected content-type")
				require.Equal(s.T(), "AUTH", captured.bodyJSON["code"])
				require.Equal(s.T(), "STATE2", captured.bodyJSON["state"])
				require.Equal(s.T(), oauth.ClientID, captured.bodyJSON["client_id"])
				require.Equal(s.T(), oauth.RedirectURI, captured.bodyJSON["redirect_uri"])
				require.Equal(s.T(), "ver", captured.bodyJSON["code_verifier"])
				// Regular OAuth should not include expires_in
				require.Nil(s.T(), captured.bodyJSON["expires_in"], "regular OAuth should not include expires_in")
		placeholder,
	placeholder,
		{
			name: "setup_token_includes_expires_in",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(oauth.TokenResponse{
					AccessToken: "at",
					TokenType:   "bearer",
					ExpiresIn:   31536000,
			placeholder)
		placeholder,
			code:         "AUTH",
			isSetupToken: true,
			wantResp: &oauth.TokenResponse{
				AccessToken: "at",
		placeholder,
			validate: func(captured requestCapture) {
				// Setup token should include expires_in with 1 year value
				require.Equal(s.T(), float64(31536000), captured.bodyJSON["expires_in"],
					"setup token should include expires_in: 31536000")
		placeholder,
	placeholder,
		{
			name: "non_200_returns_error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("bad request"))
		placeholder,
			code:         "AUTH",
			isSetupToken: false,
			wantErr:      true,
	placeholder,
placeholder

	for _, tt := range tests {
		s.Run(tt.name, func() {
			var captured requestCapture

			rt := newInProcessTransport(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				captured.method = r.Method
				captured.contentType = r.Header.Get("Content-Type")
				captured.body, _ = io.ReadAll(r.Body)
				_ = json.Unmarshal(captured.body, &captured.bodyJSON)
				tt.handler(w, r)
		placeholder), nil)

			client, ok := NewClaudeOAuthClient().(*claudeOAuthService)
			require.True(s.T(), ok, "type assertion failed")
			s.client = client
			s.client.tokenURL = "http://in-process/token"
			s.client.clientFactory = func(string) *req.Client { return newTestReqClient(rt) placeholder

			resp, err := s.client.ExchangeCodeForToken(context.Background(), tt.code, "ver", "", "", tt.isSetupToken)

			if tt.wantErr {
				require.Error(s.T(), err)
				return
		placeholder

			require.NoError(s.T(), err)
			require.Equal(s.T(), tt.wantResp.AccessToken, resp.AccessToken)
			require.Equal(s.T(), tt.wantResp.RefreshToken, resp.RefreshToken)
			if tt.validate != nil {
				tt.validate(captured)
		placeholder
	placeholder)
placeholder
placeholder

func (s *ClaudeOAuthServiceSuite) TestRefreshToken() {
	tests := []struct {
		name     string
		handler  http.HandlerFunc
		wantErr  bool
		wantResp *oauth.TokenResponse
		validate func(captured requestCapture)
placeholder{
		{
			name: "sends_json_format",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(oauth.TokenResponse{
					AccessToken:  "new_access_token",
					TokenType:    "bearer",
					ExpiresIn:    28800,
					RefreshToken: "new_refresh_token",
					Scope:        "user:profile user:inference",
			placeholder)
		placeholder,
			wantResp: &oauth.TokenResponse{
				AccessToken:  "new_access_token",
				RefreshToken: "new_refresh_token",
		placeholder,
			validate: func(captured requestCapture) {
				require.Equal(s.T(), http.MethodPost, captured.method, "expected POST")
				// 验证使用 JSON 格式（不是 form 格式）
				require.True(s.T(), strings.HasPrefix(captured.contentType, "application/json"),
					"expected JSON content-type, got: %s", captured.contentType)
				// 验证 JSON body 内容
				require.Equal(s.T(), "refresh_token", captured.bodyJSON["grant_type"])
				require.Equal(s.T(), "rt", captured.bodyJSON["refresh_token"])
				require.Equal(s.T(), oauth.ClientID, captured.bodyJSON["client_id"])
		placeholder,
	placeholder,
		{
			name: "returns_new_refresh_token",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(oauth.TokenResponse{
					AccessToken:  "at",
					TokenType:    "bearer",
					ExpiresIn:    28800,
					RefreshToken: "rotated_rt", // Anthropic rotates refresh tokens
			placeholder)
		placeholder,
			wantResp: &oauth.TokenResponse{
				AccessToken:  "at",
				RefreshToken: "rotated_rt",
		placeholder,
	placeholder,
		{
			name: "non_200_returns_error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"invalid_grant"placeholder`))
		placeholder,
			wantErr: true,
	placeholder,
placeholder

	for _, tt := range tests {
		s.Run(tt.name, func() {
			var captured requestCapture

			rt := newInProcessTransport(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				captured.method = r.Method
				captured.contentType = r.Header.Get("Content-Type")
				captured.body, _ = io.ReadAll(r.Body)
				_ = json.Unmarshal(captured.body, &captured.bodyJSON)
				tt.handler(w, r)
		placeholder), nil)

			client, ok := NewClaudeOAuthClient().(*claudeOAuthService)
			require.True(s.T(), ok, "type assertion failed")
			s.client = client
			s.client.tokenURL = "http://in-process/token"
			s.client.clientFactory = func(string) *req.Client { return newTestReqClient(rt) placeholder

			resp, err := s.client.RefreshToken(context.Background(), "rt", "")

			if tt.wantErr {
				require.Error(s.T(), err)
				return
		placeholder

			require.NoError(s.T(), err)
			require.Equal(s.T(), tt.wantResp.AccessToken, resp.AccessToken)
			require.Equal(s.T(), tt.wantResp.RefreshToken, resp.RefreshToken)
			if tt.validate != nil {
				tt.validate(captured)
		placeholder
	placeholder)
placeholder
placeholder

func TestClaudeOAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(ClaudeOAuthServiceSuite))
placeholder
