package service

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/antigravity"
	"github.com/Wei-Shaw/sub2api/internal/pkg/apicompat"
	"github.com/gin-gonic/gin"
)

type antigravityCompatProtocol uint8

const (
	antigravityCompatChatCompletions antigravityCompatProtocol = iota
	antigravityCompatResponses
)

const (
	// AntigravityCredentialRejectedClientMessage 是可安全返回给客户端的认证修复提示。
	AntigravityCredentialRejectedClientMessage = "Antigravity rejected the OAuth credential after refresh; reauthorize the account and verify project_id"
	// AntigravityCredentialRejectedReason 标识上游拒绝已刷新 OAuth 凭据。
	AntigravityCredentialRejectedReason GatewayFailureReason = "antigravity_oauth_credential_rejected"
)

type antigravityCompatRequest struct {
	protocol        antigravityCompatProtocol
	originalBody    []byte
	claudeBody      []byte
	originalModel   string
	clientStream    bool
	includeUsage    bool
	startTime       time.Time
	reasoningEffort *string
placeholder

type antigravityCompatUpstreamCall struct {
	request      antigravityCompatRequest
	billingModel string
	prefix       string
	proxyURL     string
	accessToken  string
	geminiBody   []byte
placeholder

// ForwardAsChatCompletions 使用 Antigravity 原生 OAuth 账号转发 Chat Completions 请求。
func (s *AntigravityGatewayService) ForwardAsChatCompletions(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	body []byte,
	_ *ParsedRequest,
) (*ForwardResult, error) {
	if err := s.validateAntigravityCompatAccount(c, account); err != nil {
		return nil, err
placeholder

	var request apicompat.ChatCompletionsRequest
	if json.Unmarshal(body, &request) != nil {
		return nil, s.writeAntigravityCompatError(c, http.StatusBadRequest, "invalid_request_error", "Failed to parse request body")
placeholder
	if strings.TrimSpace(request.Model) == "" {
		return nil, s.writeAntigravityCompatError(c, http.StatusBadRequest, "invalid_request_error", "model is required")
placeholder

	responsesRequest, err := apicompat.ChatCompletionsToResponses(&request)
	if err != nil {
		return nil, s.writeAntigravityCompatError(c, http.StatusBadRequest, "invalid_request_error", err.Error())
placeholder
	claudeRequest, err := apicompat.ResponsesToAnthropicRequest(responsesRequest)
	if err != nil {
		return nil, s.writeAntigravityCompatError(c, http.StatusBadRequest, "invalid_request_error", err.Error())
placeholder
	preserveChatCompletionTokenLimit(&request, claudeRequest)
	claudeRequest.Stream = request.Stream
	claudeBody, err := json.Marshal(claudeRequest)
	if err != nil {
		return nil, fmt.Errorf("marshal anthropic request: %w", err)
placeholder

	return s.forwardAntigravityCompat(ctx, c, account, antigravityCompatRequest{
		protocol:        antigravityCompatChatCompletions,
		originalBody:    body,
		claudeBody:      claudeBody,
		originalModel:   request.Model,
		clientStream:    request.Stream,
		includeUsage:    request.StreamOptions != nil && request.StreamOptions.IncludeUsage,
		startTime:       time.Now(),
		reasoningEffort: extractCCReasoningEffortFromBody(body),
placeholder)
placeholder

// ForwardAsResponses 使用 Antigravity 原生 OAuth 账号转发 Responses 请求。
func (s *AntigravityGatewayService) ForwardAsResponses(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	body []byte,
	_ *ParsedRequest,
) (*ForwardResult, error) {
	if err := s.validateAntigravityCompatAccount(c, account); err != nil {
		return nil, err
placeholder

	var request apicompat.ResponsesRequest
	if json.Unmarshal(body, &request) != nil {
		return nil, s.writeAntigravityCompatError(c, http.StatusBadRequest, "invalid_request_error", "Failed to parse request body")
placeholder
	if strings.TrimSpace(request.Model) == "" {
		return nil, s.writeAntigravityCompatError(c, http.StatusBadRequest, "invalid_request_error", "model is required")
placeholder

	claudeRequest, err := apicompat.ResponsesToAnthropicRequest(&request)
	if err != nil {
		return nil, s.writeAntigravityCompatError(c, http.StatusBadRequest, "invalid_request_error", err.Error())
placeholder
	claudeRequest.Stream = request.Stream
	claudeBody, err := json.Marshal(claudeRequest)
	if err != nil {
		return nil, fmt.Errorf("marshal anthropic request: %w", err)
placeholder

	return s.forwardAntigravityCompat(ctx, c, account, antigravityCompatRequest{
		protocol:        antigravityCompatResponses,
		originalBody:    body,
		claudeBody:      claudeBody,
		originalModel:   request.Model,
		clientStream:    request.Stream,
		startTime:       time.Now(),
		reasoningEffort: ExtractResponsesReasoningEffortFromBody(body),
placeholder)
placeholder

func (s *AntigravityGatewayService) validateAntigravityCompatAccount(c *gin.Context, account *Account) error {
	if account != nil && account.Platform == PlatformAntigravity && account.Type == AccountTypeOAuth {
		return nil
placeholder
	return s.writeAntigravityCompatError(
		c,
		http.StatusBadRequest,
		"invalid_request_error",
		"native OAuth account required for antigravity compatibility mode",
	)
placeholder

func preserveChatCompletionTokenLimit(request *apicompat.ChatCompletionsRequest, claudeRequest *apicompat.AnthropicRequest) {
	if request == nil || claudeRequest == nil {
		return
placeholder
	limit := request.MaxTokens
	if request.MaxCompletionTokens != nil {
		limit = request.MaxCompletionTokens
placeholder
	if limit != nil && *limit > 0 {
		claudeRequest.MaxTokens = *limit
placeholder
placeholder

func (s *AntigravityGatewayService) forwardAntigravityCompat(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	request antigravityCompatRequest,
) (*ForwardResult, error) {
	call, err := s.prepareAntigravityCompatCall(ctx, c, account, request)
	if err != nil {
		return nil, err
placeholder

	result, err := s.antigravityRetryLoop(antigravityRetryLoopParams{
		ctx:             ctx,
		prefix:          call.prefix,
		account:         account,
		proxyURL:        call.proxyURL,
		accessToken:     call.accessToken,
		action:          "streamGenerateContent",
		body:            call.geminiBody,
		c:               c,
		httpUpstream:    s.httpUpstream,
		settingService:  s.settingService,
		accountRepo:     s.accountRepo,
		handleError:     s.handleUpstreamError,
		requestedModel:  request.originalModel,
		isStickySession: false,
		groupID:         0,
		sessionHash:     "",
placeholder)
	if err != nil {
		return nil, s.handleAntigravityCompatTransportError(c, err)
placeholder

	return s.consumeAntigravityCompatResponse(ctx, c, account, call, result.resp)
placeholder

func (s *AntigravityGatewayService) prepareAntigravityCompatCall(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	request antigravityCompatRequest,
) (*antigravityCompatUpstreamCall, error) {
	var claudeRequest antigravity.ClaudeRequest
	if json.Unmarshal(request.claudeBody, &claudeRequest) != nil {
		return nil, s.writeAntigravityCompatError(c, http.StatusBadRequest, "invalid_request_error", "Invalid request body")
placeholder

	mappedModel := s.getMappedModel(account, request.originalModel)
	if mappedModel == "" {
		MarkOpsClientBusinessLimited(c, OpsClientBusinessLimitedReasonLocalFeatureGate)
		message := fmt.Sprintf("model %s not in whitelist", request.originalModel)
		return nil, s.writeAntigravityCompatError(c, http.StatusForbidden, "permission_error", message)
placeholder
	thinkingEnabled := claudeRequest.Thinking != nil &&
		(claudeRequest.Thinking.Type == "enabled" || claudeRequest.Thinking.Type == "adaptive")
	mappedModel = applyThinkingModelSuffix(mappedModel, thinkingEnabled)

	if s.tokenProvider == nil {
		return nil, s.writeAntigravityCompatError(c, http.StatusBadGateway, "api_error", "Antigravity token provider not configured")
placeholder
	accessToken, err := s.tokenProvider.GetAccessToken(ctx, account)
	if err != nil {
		return nil, &UpstreamFailoverError{
			StatusCode:   http.StatusBadGateway,
			ResponseBody: []byte(`{"error":{"type":"authentication_error","message":"Failed to get upstream access token"placeholder,"type":"error"placeholder`),
	placeholder
placeholder

	projectID, err := resolveAntigravityProjectID(account)
	if err != nil {
		_ = s.writeAntigravityCompatError(c, http.StatusBadRequest, "invalid_request_error", err.Error())
		return nil, err
placeholder
	geminiBody, err := s.buildAntigravityCompatGeminiBody(ctx, request.claudeBody, &claudeRequest, projectID, mappedModel)
	if err != nil {
		return nil, s.writeAntigravityCompatError(c, http.StatusBadRequest, "invalid_request_error", "Invalid request")
placeholder

	request.reasoningEffort = ApplyThinkingEnabledFallback(request.reasoningEffort, request.originalBody, mappedModel)
	return &antigravityCompatUpstreamCall{
		request:      request,
		billingModel: mappedModel,
		prefix:       logPrefix(getSessionID(c), account.Name),
		proxyURL:     antigravityCompatProxyURL(account),
		accessToken:  accessToken,
		geminiBody:   geminiBody,
placeholder, nil
placeholder

func (s *AntigravityGatewayService) buildAntigravityCompatGeminiBody(
	ctx context.Context,
	claudeBody []byte,
	claudeRequest *antigravity.ClaudeRequest,
	projectID string,
	mappedModel string,
) ([]byte, error) {
	if strings.HasPrefix(strings.ToLower(mappedModel), "gemini-") {
		body, err := convertClaudeMessagesToGeminiGenerateContent(claudeBody)
		if err != nil {
			return nil, err
	placeholder
		body = ensureGeminiFunctionCallThoughtSignatures(body)
		body, err = injectIdentityPatchToGeminiRequest(body)
		if err != nil {
			return nil, err
	placeholder
		if cleaned, cleanErr := cleanGeminiRequest(body); cleanErr == nil {
			body = cleaned
	placeholder
		return s.wrapV1InternalRequest(projectID, mappedModel, body)
placeholder

	options := s.getClaudeTransformOptions(ctx)
	options.EnableIdentityPatch = true
	return antigravity.TransformClaudeToGeminiWithOptions(claudeRequest, projectID, mappedModel, options)
placeholder

func antigravityCompatProxyURL(account *Account) string {
	if account.ProxyID == nil || account.Proxy == nil {
		return ""
placeholder
	return account.Proxy.URL()
placeholder

func (s *AntigravityGatewayService) handleAntigravityCompatTransportError(c *gin.Context, err error) error {
	if switchErr, ok := IsAntigravityAccountSwitchError(err); ok {
		return &UpstreamFailoverError{
			StatusCode:        http.StatusServiceUnavailable,
			ForceCacheBilling: switchErr.IsStickySession,
	placeholder
placeholder
	if c.Request.Context().Err() != nil {
		return s.writeAntigravityCompatError(c, http.StatusBadGateway, "client_disconnected", "Client disconnected before upstream response")
placeholder
	return s.writeAntigravityCompatError(c, http.StatusBadGateway, "upstream_error", "Upstream request failed after retries")
placeholder

func (s *AntigravityGatewayService) consumeAntigravityCompatResponse(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	call *antigravityCompatUpstreamCall,
	resp *http.Response,
) (*ForwardResult, error) {
	defer func() { _ = resp.Body.Close() placeholder()
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, s.handleAntigravityCompatHTTPError(ctx, c, account, call, resp)
placeholder

	requestID := resp.Header.Get("x-request-id")
	if requestID != "" {
		c.Header("x-request-id", requestID)
placeholder
	streamResult, err := s.consumeAntigravityCompatSuccess(c, call, resp)
	if err != nil {
		return nil, err
placeholder
	if streamResult.usage == nil {
		streamResult.usage = &ClaudeUsage{placeholder
placeholder

	return &ForwardResult{
		RequestID:        requestID,
		Usage:            *streamResult.usage,
		Model:            call.request.originalModel,
		UpstreamModel:    call.billingModel,
		Stream:           call.request.clientStream,
		Duration:         time.Since(call.request.startTime),
		FirstTokenMs:     streamResult.firstTokenMs,
		ReasoningEffort:  call.request.reasoningEffort,
		ClientDisconnect: streamResult.clientDisconnect,
placeholder, nil
placeholder

func (s *AntigravityGatewayService) consumeAntigravityCompatSuccess(
	c *gin.Context,
	call *antigravityCompatUpstreamCall,
	resp *http.Response,
) (*antigravityStreamResult, error) {
	if call.request.clientStream {
		if call.request.protocol == antigravityCompatChatCompletions {
			return s.handleChatCompletionsStreamingFromAntigravity(
				c,
				resp,
				call.request.startTime,
				call.request.originalModel,
				call.request.includeUsage,
			)
	placeholder
		return s.handleResponsesStreamingFromAntigravity(c, resp, call.request.startTime, call.request.originalModel)
placeholder

	if call.request.protocol == antigravityCompatChatCompletions {
		return s.handleChatCompletionsNonStreamingFromAntigravity(c, resp, call.request.startTime, call.request.originalModel)
placeholder
	return s.handleResponsesNonStreamingFromAntigravity(c, resp, call.request.startTime, call.request.originalModel)
placeholder

func (s *AntigravityGatewayService) handleAntigravityCompatHTTPError(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	call *antigravityCompatUpstreamCall,
	resp *http.Response,
) error {
	body := s.readUpstreamErrorBody(resp)
	s.handleUpstreamError(
		ctx,
		call.prefix,
		account,
		resp.StatusCode,
		resp.Header,
		body,
		call.request.originalModel,
		0,
		"",
		false,
	)
	if s.shouldFailoverUpstreamError(resp.StatusCode) {
		message := sanitizeUpstreamErrorMessage(strings.TrimSpace(extractAntigravityErrorMessage(body)))
		event := OpsUpstreamErrorEvent{
			Platform:           account.Platform,
			AccountID:          account.ID,
			AccountName:        account.Name,
			UpstreamStatusCode: resp.StatusCode,
			UpstreamRequestID:  resp.Header.Get("x-request-id"),
			Kind:               "failover",
			Message:            message,
			Detail:             s.getUpstreamErrorDetail(body),
	placeholder
		if resp.StatusCode == http.StatusUnauthorized {
			event.Stage = string(GatewayFailureStageAccountAuth)
			event.Scope = string(GatewayFailureScopeAccount)
			event.Reason = string(AntigravityCredentialRejectedReason)
			appendOpsUpstreamError(c, event)
			return antigravityCredentialRejectedError(resp, body)
	placeholder
		appendOpsUpstreamError(c, event)
		return &UpstreamFailoverError{
			StatusCode:      resp.StatusCode,
			ResponseBody:    body,
			ResponseHeaders: resp.Header.Clone(),
	placeholder
placeholder
	return s.writeMappedAntigravityCompatError(c, account, resp.StatusCode, resp.Header.Get("x-request-id"), body)
placeholder

func antigravityCredentialRejectedError(resp *http.Response, body []byte) *UpstreamFailoverError {
	return &UpstreamFailoverError{
		StatusCode:        resp.StatusCode,
		ResponseBody:      body,
		ResponseHeaders:   resp.Header.Clone(),
		Stage:             GatewayFailureStageAccountAuth,
		Scope:             GatewayFailureScopeAccount,
		Reason:            AntigravityCredentialRejectedReason,
		NextAccountAction: NextAccountRetry,
		ClientStatusCode:  http.StatusBadGateway,
		ClientMessage:     AntigravityCredentialRejectedClientMessage,
placeholder
placeholder

func (s *AntigravityGatewayService) writeAntigravityCompatError(
	c *gin.Context,
	status int,
	errType string,
	message string,
) error {
	MarkResponseCommitted(c)
	c.JSON(status, gin.H{
		"error": gin.H{
			"message": message,
			"type":    errType,
			"param":   nil,
			"code":    nil,
	placeholder,
placeholder)
	return errors.New(message)
placeholder

func (s *AntigravityGatewayService) writeMappedAntigravityCompatError(
	c *gin.Context,
	account *Account,
	upstreamStatus int,
	upstreamRequestID string,
	body []byte,
) error {
	MarkResponseCommitted(c)
	message := sanitizeUpstreamErrorMessage(strings.TrimSpace(extractAntigravityErrorMessage(body)))
	setOpsUpstreamError(c, upstreamStatus, message, s.getUpstreamErrorDetail(body))
	appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
		Platform:           account.Platform,
		AccountID:          account.ID,
		AccountName:        account.Name,
		UpstreamStatusCode: upstreamStatus,
		UpstreamRequestID:  upstreamRequestID,
		Kind:               "http_error",
		Message:            message,
placeholder)
	c.JSON(mapUpstreamStatusCode(upstreamStatus), gin.H{
		"error": gin.H{
			"message": getPassthroughOrDefault(message, "Upstream request failed"),
			"type":    "upstream_error",
			"param":   nil,
			"code":    nil,
	placeholder,
placeholder)
	return fmt.Errorf("upstream error: %d %s", upstreamStatus, message)
placeholder

func (s *AntigravityGatewayService) handleChatCompletionsNonStreamingFromAntigravity(
	c *gin.Context,
	resp *http.Response,
	startTime time.Time,
	originalModel string,
) (*antigravityStreamResult, error) {
	claudeResponse, result, err := s.collectClaudeStreamResponse(resp, startTime, originalModel)
	if err != nil {
		return nil, s.mapAntigravityCompatCollectionError(c, err)
placeholder
	var anthropicResponse apicompat.AnthropicResponse
	if json.Unmarshal(claudeResponse, &anthropicResponse) != nil {
		return nil, s.writeAntigravityCompatError(c, http.StatusBadGateway, "upstream_error", "Failed to parse upstream response")
placeholder
	responsesResponse := apicompat.AnthropicToResponsesResponse(&anthropicResponse)
	c.JSON(http.StatusOK, apicompat.ResponsesToChatCompletions(responsesResponse, originalModel))
	return result, nil
placeholder

func (s *AntigravityGatewayService) handleResponsesNonStreamingFromAntigravity(
	c *gin.Context,
	resp *http.Response,
	startTime time.Time,
	originalModel string,
) (*antigravityStreamResult, error) {
	claudeResponse, result, err := s.collectClaudeStreamResponse(resp, startTime, originalModel)
	if err != nil {
		return nil, s.mapAntigravityCompatCollectionError(c, err)
placeholder
	var anthropicResponse apicompat.AnthropicResponse
	if json.Unmarshal(claudeResponse, &anthropicResponse) != nil {
		return nil, s.writeAntigravityCompatError(c, http.StatusBadGateway, "upstream_error", "Failed to parse upstream response")
placeholder
	c.JSON(http.StatusOK, apicompat.AnthropicToResponsesResponse(&anthropicResponse))
	return result, nil
placeholder

func (s *AntigravityGatewayService) mapAntigravityCompatCollectionError(c *gin.Context, err error) error {
	var failoverError *UpstreamFailoverError
	if errors.As(err, &failoverError) {
		return err
placeholder
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
placeholder
	if strings.Contains(err.Error(), "stream data interval timeout") {
		return s.writeAntigravityCompatError(c, http.StatusBadGateway, "upstream_timeout", "Upstream stream data interval timeout")
placeholder
	if errors.Is(err, bufio.ErrTooLong) {
		return s.writeAntigravityCompatError(c, http.StatusBadGateway, "response_too_large", "Upstream response line too long")
placeholder
	return s.writeAntigravityCompatError(c, http.StatusBadGateway, "upstream_error", "Failed to parse upstream response")
placeholder
