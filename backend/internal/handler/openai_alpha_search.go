package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	pkghttputil "github.com/Wei-Shaw/sub2api/internal/pkg/httputil"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

// AlphaSearch proxies the standalone search endpoint used by Codex Responses Lite.
func (h *OpenAIGatewayHandler) AlphaSearch(c *gin.Context) {
	streamStarted := false
	defer h.recoverResponsesPanic(c, &streamStarted)
	setOpenAIClientTransportHTTP(c)
	requestStart := time.Now()

	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok || apiKey.Group == nil {
		h.errorResponse(c, http.StatusUnauthorized, "authentication_error", "Invalid API key")
		return
placeholder
	if apiKey.Group.Platform != service.PlatformOpenAI {
		h.errorResponse(c, http.StatusNotFound, "not_found_error", "Codex alpha search is only available for OpenAI groups")
		return
placeholder
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		h.errorResponse(c, http.StatusInternalServerError, "api_error", "User context not found")
		return
placeholder
	reqLog := requestLogger(
		c,
		"handler.openai_gateway.alpha_search",
		zap.Int64("user_id", subject.UserID),
		zap.Int64("api_key_id", apiKey.ID),
		zap.Any("group_id", apiKey.GroupID),
	)
	if !h.ensureResponsesDependencies(c, reqLog) {
		return
placeholder

	body, err := pkghttputil.ReadRequestBodyWithPrealloc(c.Request)
	if err != nil {
		if maxErr, ok := extractMaxBytesError(err); ok {
			h.errorResponse(c, http.StatusRequestEntityTooLarge, "invalid_request_error", buildBodyTooLargeMessage(maxErr.Limit))
			return
	placeholder
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Failed to read request body")
		return
placeholder
	if len(body) == 0 {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Request body is empty")
		return
placeholder
	if !gjson.ValidBytes(body) {
		logRequestBodyParseFailure(reqLog, body, nil)
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Failed to parse request body")
		return
placeholder

	modelResult := gjson.GetBytes(body, "model")
	if !modelResult.Exists() || modelResult.Type != gjson.String || strings.TrimSpace(modelResult.String()) == "" {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "model is required")
		return
placeholder
	requestedModel := strings.TrimSpace(modelResult.String())
	reqLog = reqLog.With(zap.String("model", requestedModel))
	setOpsRequestContext(c, requestedModel, false)
	setOpsEndpointContext(c, "", int16(service.RequestTypeSync))

	channelMapping, _ := h.gatewayService.ResolveChannelMappingAndRestrict(c.Request.Context(), apiKey.GroupID, requestedModel)
	forwardBody := openAIModelMappedBody(body, channelMapping.Mapped, channelMapping.MappedModel, h.gatewayService.ReplaceModelInBody)
	subscription, _ := middleware2.GetSubscriptionFromContext(c)
	service.SetOpsLatencyMs(c, service.OpsAuthLatencyMsKey, time.Since(requestStart).Milliseconds())

	userRelease, acquired := h.acquireResponsesUserSlot(c, subject.UserID, subject.Concurrency, false, &streamStarted, reqLog)
	if !acquired {
		return
placeholder
	if userRelease != nil {
		defer userRelease()
placeholder

	if err := h.billingCacheService.CheckBillingEligibility(c.Request.Context(), apiKey.User, apiKey, apiKey.Group, subscription, service.QuotaPlatform(c.Request.Context(), apiKey)); err != nil {
		status, code, message, retryAfter := billingErrorDetails(err)
		if retryAfter > 0 {
			c.Header("Retry-After", strconv.Itoa(retryAfter))
	placeholder
		h.errorResponse(c, status, code, message)
		return
placeholder

	searchID := strings.TrimSpace(gjson.GetBytes(body, "id").String())
	sessionHash := h.gatewayService.GenerateSessionHashWithFallback(c, nil, searchID)
	failedAccountIDs := make(map[int64]struct{placeholder)
	var lastFailoverErr *service.UpstreamFailoverError
	switchCount := 0
	routingStart := time.Now()

	for {
		selection, _, err := h.gatewayService.SelectAccountWithSchedulerForCapability(
			c.Request.Context(),
			apiKey.GroupID,
			"",
			sessionHash,
			requestedModel,
			failedAccountIDs,
			service.OpenAIUpstreamTransportHTTPSSE,
			service.OpenAIEndpointCapabilityChatCompletions,
			false,
			false,
			service.PlatformOpenAI,
		)
		if err != nil || selection == nil || selection.Account == nil {
			if len(failedAccountIDs) == 0 {
				cls := classifyNoAccountErrorFromGin(c, h.gatewayService, apiKey, requestedModel, requestedModel, service.PlatformOpenAI)
				if !cls.ModelNotFound {
					markOpsRoutingCapacityLimitedIfNoAvailable(c, err)
			placeholder
				h.errorResponse(c, cls.Status, cls.ErrType, cls.Message)
				return
		placeholder
			if lastFailoverErr != nil {
				h.handleFailoverExhausted(c, lastFailoverErr, false)
		placeholder else {
				h.errorResponse(c, http.StatusBadGateway, "upstream_error", "Upstream request failed")
		placeholder
			return
	placeholder

		account := selection.Account
		setOpsSelectedAccount(c, account.ID, account.Platform)
		accountRelease, acquired := h.acquireResponsesAccountSlot(c, apiKey.GroupID, sessionHash, selection, false, &streamStarted, reqLog)
		if !acquired {
			return
	placeholder
		service.SetOpsLatencyMs(c, service.OpsRoutingLatencyMsKey, time.Since(routingStart).Milliseconds())
		writerSizeBeforeForward := c.Writer.Size()
		forwardStart := time.Now()
		err = func() error {
			if accountRelease != nil {
				defer accountRelease()
		placeholder
			return h.gatewayService.ForwardAlphaSearch(c.Request.Context(), c, account, forwardBody)
	placeholder()
		service.SetOpsLatencyMs(c, service.OpsResponseLatencyMsKey, time.Since(forwardStart).Milliseconds())

		if err == nil {
			h.gatewayService.ReportOpenAIAccountScheduleResult(account.ID, true, nil)
			return
	placeholder

		var failoverErr *service.UpstreamFailoverError
		if !errors.As(err, &failoverErr) {
			h.gatewayService.ReportOpenAIAccountScheduleResult(account.ID, false, nil)
			if c.Writer.Size() == writerSizeBeforeForward {
				h.errorResponse(c, http.StatusBadGateway, "upstream_error", "Upstream request failed")
		placeholder
			reqLog.Warn("openai_alpha_search.forward_failed", zap.Int64("account_id", account.ID), zap.Error(err))
			return
	placeholder

		h.gatewayService.ReportOpenAIAccountScheduleResult(account.ID, false, nil)
		if c.Writer.Size() != writerSizeBeforeForward {
			h.handleFailoverExhausted(c, failoverErr, true)
			return
	placeholder
		h.gatewayService.RecordOpenAIAccountSwitch()
		failedAccountIDs[account.ID] = struct{placeholder{placeholder
		lastFailoverErr = failoverErr
		if switchCount >= h.maxAccountSwitches {
			h.handleFailoverExhausted(c, failoverErr, false)
			return
	placeholder
		switchCount++
		if h.gatewayService.ShouldStopOpenAIOAuth429Failover(account, failoverErr.StatusCode, switchCount) {
			h.handleFailoverExhausted(c, failoverErr, false)
			return
	placeholder
		reqLog.Warn("openai_alpha_search.upstream_failover_switching",
			zap.Int64("account_id", account.ID),
			zap.Int("upstream_status", failoverErr.StatusCode),
			zap.Int("switch_count", switchCount),
		)
placeholder
placeholder
