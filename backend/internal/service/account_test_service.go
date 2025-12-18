package service

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"sub2api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	testClaudeAPIURL = "https://api.anthropic.com/v1/messages"
	testModel        = "claude-sonnet-4-5-20250929"
)

// TestEvent represents a SSE event for account testing
type TestEvent struct {
	Type    string `json:"type"`
	Text    string `json:"text,omitempty"`
	Model   string `json:"model,omitempty"`
	Success bool   `json:"success,omitempty"`
	Error   string `json:"error,omitempty"`
placeholder

// AccountTestService handles account testing operations
type AccountTestService struct {
	repos        *repository.Repositories
	oauthService *OAuthService
	httpClient   *http.Client
placeholder

// NewAccountTestService creates a new AccountTestService
func NewAccountTestService(repos *repository.Repositories, oauthService *OAuthService) *AccountTestService {
	return &AccountTestService{
		repos:        repos,
		oauthService: oauthService,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
	placeholder,
placeholder
placeholder

// generateSessionString generates a Claude Code style session string
func generateSessionString() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	hex64 := hex.EncodeToString(bytes)
	sessionUUID := uuid.New().String()
	return fmt.Sprintf("user_%s_account__session_%s", hex64, sessionUUID)
placeholder

// createTestPayload creates a minimal test request payload for OAuth/Setup Token accounts
func createTestPayload() map[string]interface{placeholder {
	return map[string]interface{placeholder{
		"model": testModel,
		"messages": []map[string]interface{placeholder{
			{
				"role": "user",
				"content": []map[string]interface{placeholder{
					{
						"type": "text",
						"text": "hi",
						"cache_control": map[string]string{
							"type": "ephemeral",
					placeholder,
				placeholder,
			placeholder,
		placeholder,
	placeholder,
		"system": []map[string]interface{placeholder{
			{
				"type": "text",
				"text": "You are Claude Code, Anthropic's official CLI for Claude.",
				"cache_control": map[string]string{
					"type": "ephemeral",
			placeholder,
		placeholder,
	placeholder,
		"metadata": map[string]string{
			"user_id": generateSessionString(),
	placeholder,
		"max_tokens": 1024,
		"temperature": 1,
		"stream":     true,
placeholder
placeholder

// createApiKeyTestPayload creates a simpler test request payload for API Key accounts
func createApiKeyTestPayload(model string) map[string]interface{placeholder {
	return map[string]interface{placeholder{
		"model": model,
		"messages": []map[string]interface{placeholder{
			{
				"role":    "user",
				"content": "hi",
		placeholder,
	placeholder,
		"max_tokens": 1024,
		"stream":     true,
placeholder
placeholder

// TestAccountConnection tests an account's connection by sending a test request
func (s *AccountTestService) TestAccountConnection(c *gin.Context, accountID int64) error {
	ctx := c.Request.Context()

	// Get account
	account, err := s.repos.Account.GetByID(ctx, accountID)
	if err != nil {
		return s.sendErrorAndEnd(c, "Account not found")
placeholder

	// Determine authentication method based on account type
	var authToken string
	var authType string // "bearer" for OAuth, "apikey" for API Key
	var apiURL string

	if account.IsOAuth() {
		// OAuth or Setup Token account
		authType = "bearer"
		apiURL = testClaudeAPIURL
		authToken = account.GetCredential("access_token")
		if authToken == "" {
			return s.sendErrorAndEnd(c, "No access token available")
	placeholder

		// Check if token needs refresh
		needRefresh := false
		if expiresAtStr := account.GetCredential("expires_at"); expiresAtStr != "" {
			expiresAt, err := strconv.ParseInt(expiresAtStr, 10, 64)
			if err == nil && time.Now().Unix()+300 > expiresAt { // 5 minute buffer
				needRefresh = true
		placeholder
	placeholder

		if needRefresh && s.oauthService != nil {
			tokenInfo, err := s.oauthService.RefreshAccountToken(ctx, account)
			if err != nil {
				return s.sendErrorAndEnd(c, fmt.Sprintf("Failed to refresh token: %s", err.Error()))
		placeholder
			authToken = tokenInfo.AccessToken
	placeholder
placeholder else if account.Type == "apikey" {
		// API Key account
		authType = "apikey"
		authToken = account.GetCredential("api_key")
		if authToken == "" {
			return s.sendErrorAndEnd(c, "No API key available")
	placeholder

		// Get base URL (use default if not set)
		apiURL = account.GetBaseURL()
		if apiURL == "" {
			apiURL = "https://api.anthropic.com"
	placeholder
		// Append /v1/messages endpoint
		apiURL = strings.TrimSuffix(apiURL, "/") + "/v1/messages"
placeholder else {
		return s.sendErrorAndEnd(c, fmt.Sprintf("Unsupported account type: %s", account.Type))
placeholder

	// Set SSE headers
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.Flush()

	// Create test request payload
	var payload map[string]interface{placeholder
	var actualModel string
	if authType == "apikey" {
		// Use simpler payload for API Key (without Claude Code specific fields)
		// Apply model mapping if configured
		actualModel = account.GetMappedModel(testModel)
		payload = createApiKeyTestPayload(actualModel)
placeholder else {
		actualModel = testModel
		payload = createTestPayload()
placeholder
	payloadBytes, _ := json.Marshal(payload)

	// Send test_start event with model info
	s.sendEvent(c, TestEvent{Type: "test_start", Model: actualModelplaceholder)

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return s.sendErrorAndEnd(c, "Failed to create request")
placeholder

	// Set headers based on auth type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01")

	if authType == "bearer" {
		req.Header.Set("Authorization", "Bearer "+authToken)
		req.Header.Set("anthropic-beta", "prompt-caching-2024-07-31,placeholder,output-128k-2025-02-19")
placeholder else {
		// API Key uses x-api-key header
		req.Header.Set("x-api-key", authToken)
placeholder

	// Configure proxy if account has one
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL := account.Proxy.URL()
		if proxyURL != "" {
			if parsedURL, err := url.Parse(proxyURL); err == nil {
				transport.Proxy = http.ProxyURL(parsedURL)
		placeholder
	placeholder
placeholder

	client := &http.Client{
		Transport: transport,
		Timeout:   60 * time.Second,
placeholder

	resp, err := client.Do(req)
	if err != nil {
		return s.sendErrorAndEnd(c, fmt.Sprintf("Request failed: %s", err.Error()))
placeholder
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return s.sendErrorAndEnd(c, fmt.Sprintf("API returned %d: %s", resp.StatusCode, string(body)))
placeholder

	// Process SSE stream
	return s.processStream(c, resp.Body)
placeholder

// processStream processes the SSE stream from Claude API
func (s *AccountTestService) processStream(c *gin.Context, body io.Reader) error {
	reader := bufio.NewReader(body)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// Stream ended, send complete event
				s.sendEvent(c, TestEvent{Type: "test_complete", Success: trueplaceholder)
				return nil
		placeholder
			return s.sendErrorAndEnd(c, fmt.Sprintf("Stream read error: %s", err.Error()))
	placeholder

		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "data: ") {
			continue
	placeholder

		jsonStr := strings.TrimPrefix(line, "data: ")
		if jsonStr == "[DONE]" {
			s.sendEvent(c, TestEvent{Type: "test_complete", Success: trueplaceholder)
			return nil
	placeholder

		var data map[string]interface{placeholder
		if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
			continue
	placeholder

		eventType, _ := data["type"].(string)

		switch eventType {
		case "content_block_delta":
			if delta, ok := data["delta"].(map[string]interface{placeholder); ok {
				if text, ok := delta["text"].(string); ok {
					s.sendEvent(c, TestEvent{Type: "content", Text: textplaceholder)
			placeholder
		placeholder
		case "message_stop":
			s.sendEvent(c, TestEvent{Type: "test_complete", Success: trueplaceholder)
			return nil
		case "error":
			errorMsg := "Unknown error"
			if errData, ok := data["error"].(map[string]interface{placeholder); ok {
				if msg, ok := errData["message"].(string); ok {
					errorMsg = msg
			placeholder
		placeholder
			return s.sendErrorAndEnd(c, errorMsg)
	placeholder
placeholder
placeholder

// sendEvent sends a SSE event to the client
func (s *AccountTestService) sendEvent(c *gin.Context, event TestEvent) {
	eventJSON, _ := json.Marshal(event)
	fmt.Fprintf(c.Writer, "data: %s\n\n", eventJSON)
	c.Writer.Flush()
placeholder

// sendErrorAndEnd sends an error event and ends the stream
func (s *AccountTestService) sendErrorAndEnd(c *gin.Context, errorMsg string) error {
	log.Printf("Account test error: %s", errorMsg)
	s.sendEvent(c, TestEvent{Type: "error", Error: errorMsgplaceholder)
	return fmt.Errorf(errorMsg)
placeholder
