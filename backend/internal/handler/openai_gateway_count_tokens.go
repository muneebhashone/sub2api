package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/domain"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

// ResponsesInputTokens handles native OpenAI POST
// /v1/responses/input_tokens requests without routing them through the normal
// Responses generation and usage-recording pipeline.
func (h *OpenAIGatewayHandler) ResponsesInputTokens(c *gin.Context) {
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok {
		h.errorResponse(c, http.StatusUnauthorized, "authentication_error", "Invalid API key")
		return
placeholder
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		h.errorResponse(c, http.StatusInternalServerError, "api_error", "User context not found")
		return
placeholder
	reqLog := requestLogger(
		c,
		"handler.openai_gateway.responses_input_tokens",
		zap.Int64("user_id", subject.UserID),
		zap.Int64("api_key_id", apiKey.ID),
		zap.Any("group_id", apiKey.GroupID),
	)
	if !h.ensureResponsesDependencies(c, reqLog) {
		return
placeholder

	body, err := readLenientJSONRequestBodyWithPrealloc(c.Request, h.cfg)
	if err != nil {
		if maxErr, ok := extractMaxBytesError(err); ok {
			h.errorResponse(c, http.StatusRequestEntityTooLarge, "invalid_request_error", buildBodyTooLargeMessage(maxErr.Limit))
			return
	placeholder
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Failed to read request body")
		return
placeholder
	if len(body) == 0 || !gjson.ValidBytes(body) {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Failed to parse request body")
		return
placeholder
	modelResult := gjson.GetBytes(body, "model")
	if !modelResult.Exists() || modelResult.Type != gjson.String || strings.TrimSpace(modelResult.String()) == "" {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "model is required")
		return
placeholder
	reqModel := strings.TrimSpace(modelResult.String())
	ensureCompositeTargetPlatform(c, apiKey, reqModel)
	if !openAICompatibleTextTargetAllowed(c, apiKey, reqModel) {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Model is not supported by this OpenAI-compatible endpoint for composite groups")
		return
placeholder

	setOpsRequestContext(c, reqModel, false)
	setOpsEndpointContext(c, "", int16(service.RequestTypeFromLegacy(false, false)))
	if decision := h.checkSecurityAudit(c, reqLog, apiKey, subject, service.ContentModerationProtocolOpenAIResponses, reqModel, body); decision != nil && !decision.AllowNextStage {
		h.openAISecurityAuditError(c, decision)
		return
placeholder

	subscription, _ := middleware2.GetSubscriptionFromContext(c)
	if err := h.billingCacheService.CheckBillingEligibility(c.Request.Context(), apiKey.User, apiKey, apiKey.Group, subscription, service.QuotaPlatform(c.Request.Context(), apiKey)); err != nil {
		reqLog.Info("openai_input_tokens.billing_eligibility_check_failed", zap.Error(err))
		status, code, message, retryAfter := billingErrorDetails(err)
		if retryAfter > 0 {
			c.Header("Retry-After", strconv.Itoa(retryAfter))
	placeholder
		h.errorResponse(c, status, code, message)
		return
placeholder

	channelMapping, _ := h.gatewayService.ResolveChannelMappingAndRestrict(c.Request.Context(), apiKey.GroupID, reqModel)
	routingModel := reqModel
	forwardBody := body
	if channelMapping.Mapped {
		routingModel = channelMapping.MappedModel
		forwardBody = h.gatewayService.ReplaceModelInBody(body, routingModel)
placeholder

	// Token counting is not billed, so it must not be excluded by the profit gate.
	c.Request = c.Request.WithContext(service.WithOpenAIProfitControlSuppressed(c.Request.Context()))
	requestPlatform := openAICompatibleRequestPlatform(c.Request.Context(), apiKey)
	sessionHash := h.gatewayService.GenerateSessionHash(c, body)
	requestStart := time.Now()
	selection, _, err := h.gatewayService.SelectAccountWithSchedulerForCapability(
		c.Request.Context(),
		apiKey.GroupID,
		"",
		sessionHash,
		routingModel,
		nil,
		service.OpenAIUpstreamTransportAny,
		service.OpenAIEndpointCapabilityChatCompletions,
		false,
		false,
		false,
		requestPlatform,
	)
	service.SetOpsLatencyMs(c, service.OpsAuthLatencyMsKey, time.Since(requestStart).Milliseconds())
	if err != nil {
		reqLog.Warn("openai_input_tokens.account_select_failed", zap.Error(openAICompatibleSelectionErrorForLog(err, requestPlatform)))
		cls := classifyOpenAICompatibleNoAccountErrorFromGin(c, h.gatewayService, apiKey, routingModel, reqModel)
		if !cls.ModelNotFound {
			markOpsRoutingCapacityLimitedIfNoAvailable(c, err)
	placeholder
		h.errorResponse(c, cls.Status, cls.ErrType, cls.Message)
		return
placeholder
	if selection == nil || selection.Account == nil {
		cls := classifyOpenAICompatibleNoAccountErrorFromGin(c, h.gatewayService, apiKey, routingModel, reqModel)
		if !cls.ModelNotFound {
			markOpsRoutingCapacityLimited(c)
	placeholder
		h.errorResponse(c, cls.Status, cls.ErrType, cls.Message)
		return
placeholder

	account := selection.Account
	setOpsSelectedAccount(c, account.ID, account.Platform)
	if selection.Acquired && selection.ReleaseFunc != nil {
		defer selection.ReleaseFunc()
placeholder
	if err := h.gatewayService.ForwardResponsesInputTokens(c.Request.Context(), c, account, forwardBody); err != nil {
		reqLog.Error("openai_input_tokens.forward_failed", zap.Int64("account_id", account.ID), zap.Error(err))
placeholder
placeholder

// GrokCountTokens handles Anthropic-compatible count_tokens requests locally.
// The route middleware already authenticates the API key and resolves the
// group; this handler intentionally does not select an account or check billing.
func (h *OpenAIGatewayHandler) GrokCountTokens(c *gin.Context) {
	body, err := readLenientJSONRequestBodyWithPrealloc(c.Request, h.cfg)
	if err != nil {
		if maxErr, ok := extractMaxBytesError(err); ok {
			h.anthropicErrorResponse(c, http.StatusRequestEntityTooLarge, "invalid_request_error", buildBodyTooLargeMessage(maxErr.Limit))
			return
	placeholder
		h.anthropicErrorResponse(c, http.StatusBadRequest, "invalid_request_error", "Failed to read request body")
		return
placeholder
	if len(body) == 0 {
		h.anthropicErrorResponse(c, http.StatusBadRequest, "invalid_request_error", "Request body is empty")
		return
placeholder

	bodyRef := service.NewRequestBodyRef(body)
	parsedReq, err := service.ParseGatewayRequest(bodyRef, domain.PlatformAnthropic)
	if err != nil {
		logRequestBodyParseFailure(requestLogger(c, "handler.openai_gateway.grok_count_tokens"), body, err)
		h.anthropicErrorResponse(c, http.StatusBadRequest, "invalid_request_error", "Failed to parse request body")
		return
placeholder
	if parsedReq.Model == "" {
		h.anthropicErrorResponse(c, http.StatusBadRequest, "invalid_request_error", "model is required")
		return
placeholder

	estimated, err := service.EstimateGrokCountTokens(parsedReq.Body.Bytes())
	if err != nil {
		requestLogger(c, "handler.openai_gateway.grok_count_tokens").Warn("grok_count_tokens.local_estimate_failed", zap.Error(err))
		h.anthropicErrorResponse(c, http.StatusBadRequest, "invalid_request_error", "Failed to parse request body")
		return
placeholder

	setOpsRequestContext(c, parsedReq.Model, false)
	setOpsEndpointContext(c, "", int16(service.RequestTypeFromLegacy(false, false)))
	c.JSON(http.StatusOK, gin.H{"input_tokens": estimatedplaceholder)
placeholder

// CountTokens handles Anthropic-compatible POST /v1/messages/count_tokens for OpenAI groups.
// It validates billing and routes to an OpenAI token-count bridge without taking concurrency slots
// or recording usage.
func (h *OpenAIGatewayHandler) CountTokens(c *gin.Context) {
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok {
		h.anthropicErrorResponse(c, http.StatusUnauthorized, "authentication_error", "Invalid API key")
		return
placeholder

	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		h.anthropicErrorResponse(c, http.StatusInternalServerError, "api_error", "User context not found")
		return
placeholder
	reqLog := requestLogger(
		c,
		"handler.openai_gateway.count_tokens",
		zap.Int64("user_id", subject.UserID),
		zap.Int64("api_key_id", apiKey.ID),
		zap.Any("group_id", apiKey.GroupID),
	)

	if !allowOpenAICompatibleMessagesDispatch(c, apiKey) {
		h.anthropicErrorResponse(c, http.StatusForbidden, "permission_error",
			"This group does not allow /v1/messages dispatch")
		return
placeholder

	if !h.ensureResponsesDependencies(c, reqLog) {
		return
placeholder

	body, err := readLenientJSONRequestBodyWithPrealloc(c.Request, h.cfg)
	if err != nil {
		if maxErr, ok := extractMaxBytesError(err); ok {
			h.anthropicErrorResponse(c, http.StatusRequestEntityTooLarge, "invalid_request_error", buildBodyTooLargeMessage(maxErr.Limit))
			return
	placeholder
		h.anthropicErrorResponse(c, http.StatusBadRequest, "invalid_request_error", "Failed to read request body")
		return
placeholder
	if len(body) == 0 {
		h.anthropicErrorResponse(c, http.StatusBadRequest, "invalid_request_error", "Request body is empty")
		return
placeholder

	bodyRef := service.NewRequestBodyRef(body)
	parsedReq, err := service.ParseGatewayRequest(bodyRef, domain.PlatformAnthropic)
	if err != nil {
		logRequestBodyParseFailure(reqLog, body, err)
		h.anthropicErrorResponse(c, http.StatusBadRequest, "invalid_request_error", "Failed to parse request body")
		return
placeholder
	body = parsedReq.Body.Bytes()
	if parsedReq.Model == "" {
		h.anthropicErrorResponse(c, http.StatusBadRequest, "invalid_request_error", "model is required")
		return
placeholder

	reqModel := parsedReq.Model
	ensureCompositeTargetPlatform(c, apiKey, reqModel)
	// composite+grok 在路由层已分流到 GrokCountTokens，这里可达的目标平台是
	// openai 与 CN 供应商；CN 账号由 ForwardCountTokensAsAnthropic 本地估算。
	if !openAICompatibleTextTargetAllowed(c, apiKey, reqModel) {
		h.anthropicErrorResponse(c, http.StatusBadRequest, "invalid_request_error", "Model is not supported by this OpenAI-compatible endpoint for composite groups")
		return
placeholder
	routingModel := service.NormalizeOpenAICompatRequestedModel(reqModel)
	preferredMappedModel := resolveOpenAIMessagesDispatchMappedModel(c, apiKey, reqModel)
	reqLog = reqLog.With(zap.String("model", reqModel), zap.Bool("stream", parsedReq.Stream))

	setOpsRequestContext(c, reqModel, false)
	setOpsEndpointContext(c, "", int16(service.RequestTypeFromLegacy(false, false)))

	channelMapping, _ := h.gatewayService.ResolveChannelMappingAndRestrict(c.Request.Context(), apiKey.GroupID, reqModel)
	mappedBodyForMessages := newOpenAIModelMappedBodyCache(body, h.gatewayService.ReplaceModelInBody)

	subscription, _ := middleware2.GetSubscriptionFromContext(c)
	if err := h.billingCacheService.CheckBillingEligibility(c.Request.Context(), apiKey.User, apiKey, apiKey.Group, subscription, service.QuotaPlatform(c.Request.Context(), apiKey)); err != nil {
		reqLog.Info("openai_count_tokens.billing_eligibility_check_failed", zap.Error(err))
		status, code, message, retryAfter := billingErrorDetails(err)
		if retryAfter > 0 {
			c.Header("Retry-After", strconv.Itoa(retryAfter))
	placeholder
		h.anthropicErrorResponse(c, status, code, message)
		return
placeholder

	requestStart := time.Now()
	// count_tokens 不计费：显式豁免利润门，避免高倍率账号池被门排除后连
	// token 计数都返回 no available accounts。
	c.Request = c.Request.WithContext(service.WithOpenAIProfitControlSuppressed(c.Request.Context()))
	sessionHash := h.gatewayService.GenerateSessionHash(c, body)
	currentRoutingModel := routingModel
	if preferredMappedModel != "" {
		currentRoutingModel = preferredMappedModel
placeholder
	selection, _, err := h.gatewayService.SelectAccountWithSchedulerForCapability(
		c.Request.Context(),
		apiKey.GroupID,
		"",
		sessionHash,
		currentRoutingModel,
		nil,
		service.OpenAIUpstreamTransportAny,
		service.OpenAIEndpointCapabilityChatCompletions,
		false,
		false,
		false,
		openAICompatibleRequestPlatform(c.Request.Context(), apiKey),
	)
	service.SetOpsLatencyMs(c, service.OpsAuthLatencyMsKey, time.Since(requestStart).Milliseconds())
	if err != nil {
		requestPlatform := openAICompatibleRequestPlatform(c.Request.Context(), apiKey)
		reqLog.Warn("openai_count_tokens.account_select_failed", zap.Error(openAICompatibleSelectionErrorForLog(err, requestPlatform)))
		cls := classifyOpenAICompatibleNoAccountErrorFromGin(c, h.gatewayService, apiKey, currentRoutingModel, reqModel)
		if !cls.ModelNotFound {
			markOpsRoutingCapacityLimitedIfNoAvailable(c, err)
	placeholder
		h.anthropicErrorResponse(c, cls.Status, cls.ErrType, cls.Message)
		return
placeholder
	if selection == nil || selection.Account == nil {
		cls := classifyOpenAICompatibleNoAccountErrorFromGin(c, h.gatewayService, apiKey, currentRoutingModel, reqModel)
		if !cls.ModelNotFound {
			markOpsRoutingCapacityLimited(c)
	placeholder
		h.anthropicErrorResponse(c, cls.Status, cls.ErrType, cls.Message)
		return
placeholder

	account := selection.Account
	setOpsSelectedAccount(c, account.ID, account.Platform)
	if selection.Acquired && selection.ReleaseFunc != nil {
		defer selection.ReleaseFunc()
placeholder
	forwardBody := mappedBodyForMessages(channelMapping.Mapped, channelMapping.MappedModel)
	defaultMappedModel := preferredMappedModel

	if err := h.gatewayService.ForwardCountTokensAsAnthropic(c.Request.Context(), c, account, forwardBody, defaultMappedModel); err != nil {
		reqLog.Error("openai_count_tokens.forward_failed", zap.Int64("account_id", account.ID), zap.Error(err))
placeholder
placeholder
