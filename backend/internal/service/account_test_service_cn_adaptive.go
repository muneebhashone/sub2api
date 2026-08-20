package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/claude"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/gin-gonic/gin"
)

const accountTestSuppressCompletionContextKey = "account_test_suppress_completion"

// testCNProviderAdaptiveConnection verifies every native endpoint used by an
// adaptive CN-provider account. Kimi and Zhipu use Chat Completions plus
// Anthropic; DeepSeek additionally uses its native Responses endpoint.
func (s *AccountTestService) testCNProviderAdaptiveConnection(c *gin.Context, account *Account, modelID string, prompt string) error {
	testModelID := strings.TrimSpace(modelID)
	if testModelID == "" {
		testModelID = openai.DefaultTestModel
placeholder
	testModelID = account.GetMappedModel(testModelID)

	authToken := strings.TrimSpace(account.GetOpenAIProtocolAPIKey())
	if authToken == "" {
		return s.sendErrorAndEnd(c, "No API key available")
placeholder

	// The existing Chat probe owns the SSE lifecycle. Suppress intermediate
	// completion events until every native adaptive endpoint has passed.
	c.Set(accountTestSuppressCompletionContextKey, true)
	defer c.Set(accountTestSuppressCompletionContextKey, false)
	if err := s.testCNProviderChatCompletionsConnection(c, account, modelID, prompt); err != nil {
		return err
placeholder

	if err := s.testCNProviderAdaptiveAnthropicConnection(c, account, testModelID, authToken); err != nil {
		return err
placeholder

	if account.Platform == PlatformDeepseek {
		if err := s.testCNProviderAdaptiveResponsesConnection(c, account, testModelID, authToken); err != nil {
			return err
	placeholder
placeholder

	c.Set(accountTestSuppressCompletionContextKey, false)
	s.sendEvent(c, TestEvent{Type: "test_complete", Success: trueplaceholder)
	return nil
placeholder

func (s *AccountTestService) testCNProviderAdaptiveAnthropicConnection(c *gin.Context, account *Account, testModelID string, authToken string) error {
	ctx := c.Request.Context()
	baseURL, err := s.validateUpstreamBaseURL(account.GetCNProtocolBaseURL(APIProtocolAnthropic))
	if err != nil {
		return s.sendErrorAndEnd(c, fmt.Sprintf("Invalid adaptive Anthropic base URL: %s", err.Error()))
placeholder
	apiURL := strings.TrimRight(baseURL, "/") + "/v1/messages"

	payload, err := createTestPayload(testModelID)
	if err != nil {
		return s.sendErrorAndEnd(c, "Failed to create adaptive Anthropic test payload")
placeholder
	payloadBytes, _ := json.Marshal(payload)

	s.sendEvent(c, TestEvent{Type: "status", Text: "正在通过原生 /v1/messages 测试自适应 Anthropic 端点"placeholder)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return s.sendErrorAndEnd(c, "Failed to create adaptive Anthropic request")
placeholder
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("anthropic-version", "2023-06-01")
	for key, value := range claude.DefaultHeaders {
		req.Header.Set(key, value)
placeholder
	req.Header.Set("anthropic-beta", claude.APIKeyBetaHeader)
	setAnthropicAPIKeyAuthHeader(req.Header, account, authToken)
	account.ApplyHeaderOverrides(req.Header)

	resp, err := s.doCNProviderAdaptiveRequest(req, account)
	if err != nil {
		return s.sendErrorAndEnd(c, fmt.Sprintf("Adaptive Anthropic endpoint request failed: %s", err.Error()))
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("Adaptive Anthropic endpoint returned %d: %s", resp.StatusCode, string(body))
		if resp.StatusCode == http.StatusUnauthorized && s.accountRepo != nil {
			_ = s.accountRepo.SetError(ctx, account.ID, errMsg)
	placeholder
		return s.sendErrorAndEnd(c, errMsg)
placeholder

	if err := s.processCNProviderAdaptiveAnthropicStream(c, resp.Body); err != nil {
		return err
placeholder
	s.sendEvent(c, TestEvent{Type: "status", Text: "已通过原生 /v1/messages 验证"placeholder)
	return nil
placeholder

func (s *AccountTestService) processCNProviderAdaptiveAnthropicStream(c *gin.Context, body io.Reader) error {
	reader := bufio.NewReader(body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return s.sendErrorAndEnd(c, "Adaptive Anthropic stream ended before message_stop")
		placeholder
			return s.sendErrorAndEnd(c, fmt.Sprintf("Adaptive Anthropic stream read error: %s", err.Error()))
	placeholder

		line = strings.TrimSpace(line)
		if line == "" || !sseDataPrefix.MatchString(line) {
			continue
	placeholder
		jsonStr := sseDataPrefix.ReplaceAllString(line, "")
		if jsonStr == "[DONE]" {
			return nil
	placeholder

		var data map[string]any
		if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
			continue
	placeholder
		switch eventType, _ := data["type"].(string); eventType {
		case "content_block_delta":
			if delta, ok := data["delta"].(map[string]any); ok {
				if text, ok := delta["text"].(string); ok && text != "" {
					s.sendEvent(c, TestEvent{Type: "content", Text: textplaceholder)
			placeholder
		placeholder
		case "message_stop":
			return nil
		case "error":
			errorMsg := "Unknown error"
			if errData, ok := data["error"].(map[string]any); ok {
				if message, ok := errData["message"].(string); ok && message != "" {
					errorMsg = message
			placeholder
		placeholder
			return s.sendErrorAndEnd(c, fmt.Sprintf("Adaptive Anthropic endpoint error: %s", errorMsg))
	placeholder
placeholder
placeholder

func (s *AccountTestService) testCNProviderAdaptiveResponsesConnection(c *gin.Context, account *Account, testModelID string, authToken string) error {
	ctx := c.Request.Context()
	baseURL, err := s.validateUpstreamBaseURL(account.GetCNProtocolBaseURL(APIProtocolResponses))
	if err != nil {
		return s.sendErrorAndEnd(c, fmt.Sprintf("Invalid adaptive Responses base URL: %s", err.Error()))
placeholder
	apiURL := buildOpenAIResponsesURLForPlatform(account.Platform, baseURL)

	payload := createOpenAITestPayload(testModelID, false)
	// DeepSeek's native Responses endpoint is stateless and does not need the
	// OpenAI probe's synthetic instructions.
	delete(payload, "instructions")
	payloadBytes, _ := json.Marshal(payload)
	payloadBytes = normalizeDeepSeekResponsesRequestBody(account, payloadBytes)

	s.sendEvent(c, TestEvent{Type: "status", Text: "正在通过原生 /responses 测试自适应 Responses 端点"placeholder)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return s.sendErrorAndEnd(c, "Failed to create adaptive Responses request")
placeholder
	req = req.WithContext(WithHTTPUpstreamProfile(req.Context(), HTTPUpstreamProfileOpenAI))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Authorization", "Bearer "+authToken)
	applyOpenAICodexProbeHeaders(req.Header)
	account.ApplyHeaderOverrides(req.Header)

	resp, err := s.doCNProviderAdaptiveRequest(req, account)
	if err != nil {
		return s.sendErrorAndEnd(c, fmt.Sprintf("Adaptive Responses endpoint request failed: %s", err.Error()))
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("Adaptive Responses endpoint returned %d: %s", resp.StatusCode, string(body))
		if resp.StatusCode == http.StatusUnauthorized && s.accountRepo != nil {
			_ = s.accountRepo.SetError(ctx, account.ID, errMsg)
	placeholder
		return s.sendErrorAndEnd(c, errMsg)
placeholder

	if err := s.processOpenAIStream(c, resp.Body); err != nil {
		return err
placeholder
	s.sendEvent(c, TestEvent{Type: "status", Text: "已通过原生 /responses 验证"placeholder)
	return nil
placeholder

func (s *AccountTestService) doCNProviderAdaptiveRequest(req *http.Request, account *Account) (*http.Response, error) {
	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
placeholder
	return s.httpUpstream.DoWithTLS(req, proxyURL, account.ID, account.Concurrency, s.tlsFPProfileService.ResolveTLSProfile(account))
placeholder
