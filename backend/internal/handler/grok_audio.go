package handler

import (
	"errors"
	"io"
	"net/http"
	"strings"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	coderws "github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GrokRealtime exposes xAI's native Voice Realtime WebSocket.
// Only Grok-platform API keys may use this endpoint.
func (h *OpenAIGatewayHandler) GrokRealtime(c *gin.Context) {
	if c == nil || c.Request == nil || !isOpenAIWSUpgradeRequest(c.Request) {
		h.errorResponse(c, http.StatusUpgradeRequired, "invalid_request_error", "WebSocket upgrade required (Upgrade: websocket)")
		return
placeholder
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok || apiKey.Group == nil || apiKey.Group.Platform != service.PlatformGrok {
		h.errorResponse(c, http.StatusNotFound, "not_found_error", "Realtime API is not supported for this platform")
		return
placeholder
	if !h.ensureResponsesDependencies(c, nil) {
		return
placeholder

	selection, _, err := h.gatewayService.SelectAccountWithSchedulerForCapability(
		c.Request.Context(),
		apiKey.GroupID,
		"",
		"",
		"grok-4.5",
		nil,
		service.OpenAIUpstreamTransportHTTPSSE,
		// Grok only advertises chat_completions + media capabilities on HEAD.
		service.OpenAIEndpointCapabilityChatCompletions,
		false,
		false,
		false,
		service.PlatformGrok,
	)
	if err != nil || selection == nil || selection.Account == nil {
		h.errorResponse(c, http.StatusServiceUnavailable, "api_error", "No available Grok accounts")
		return
placeholder

	var streamStarted bool
	reqLog := requestLogger(c, "handler.openai_gateway.grok_realtime")
	release, slotStatus := h.acquireResponsesAccountSlot(c, apiKey.GroupID, "", selection, true, &streamStarted, reqLog)
	if slotStatus != openAISlotAcquireOK {
		return
placeholder
	defer release()

	token, _, err := h.gatewayService.GetRequestCredential(c.Request.Context(), c, selection.Account)
	if err != nil {
		h.errorResponse(c, http.StatusBadGateway, "upstream_error", "Grok credential unavailable")
		return
placeholder

	conn, err := coderws.Accept(c.Writer, c.Request, &coderws.AcceptOptions{CompressionMode: coderws.CompressionContextTakeoverplaceholder)
	if err != nil {
		return
placeholder
	defer conn.CloseNow()

	model := c.Query("model")
	if err := h.gatewayService.ProxyGrokRealtime(c.Request.Context(), c, conn, selection.Account, token, model); err != nil {
		reqLog.Info("grok_realtime.proxy_failed", zap.Error(err))
		_ = conn.Close(coderws.StatusInternalError, "upstream realtime websocket failed")
placeholder
placeholder

// GrokVoice handles xAI Voice HTTP endpoints. endpoint is "tts", "stt", or "custom-voices".
func (h *OpenAIGatewayHandler) GrokVoice(c *gin.Context, endpoint string) {
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok || apiKey.Group == nil || apiKey.Group.Platform != service.PlatformGrok {
		h.errorResponse(c, http.StatusNotFound, "not_found_error", "Voice API is not supported for this platform")
		return
placeholder
	if !h.ensureResponsesDependencies(c, nil) {
		return
placeholder

	body, err := readGrokVoiceGatewayBody(c)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", err.Error())
		return
placeholder
	contentType := c.GetHeader("Content-Type")
	if strings.TrimSpace(contentType) == "" {
		contentType = "application/json"
placeholder

	failed := map[int64]struct{placeholder{placeholder
	var last *service.UpstreamFailoverError
	reqLog := requestLogger(c, "handler.openai_gateway.grok_voice", zap.String("endpoint", endpoint))
	selectionModel := "grok-4.5"

	for attempts := 0; attempts < 4; attempts++ {
		selection, _, selectErr := h.gatewayService.SelectAccountWithSchedulerForCapability(
			c.Request.Context(),
			apiKey.GroupID,
			"",
			"",
			selectionModel,
			failed,
			service.OpenAIUpstreamTransportHTTPSSE,
			service.OpenAIEndpointCapabilityChatCompletions,
			false,
			false,
			false,
			service.PlatformGrok,
		)
		if selectErr != nil || selection == nil || selection.Account == nil {
			if last != nil {
				h.handleFailoverExhausted(c, last, false)
		placeholder else {
				h.errorResponse(c, http.StatusServiceUnavailable, "api_error", "No available Grok accounts")
		placeholder
			return
	placeholder
		account := selection.Account
		var started bool
		release, status := h.acquireResponsesAccountSlot(c, apiKey.GroupID, "", selection, false, &started, reqLog)
		if status == openAISlotAcquireProfitVetoed {
			failed[account.ID] = struct{placeholder{placeholder
			continue
	placeholder
		if status != openAISlotAcquireOK {
			// Failed already wrote error response (or transient reject).
			if status == openAISlotAcquireFailed && len(failed) == 0 {
				// Slot path wrote the response; stop.
				return
		placeholder
			failed[account.ID] = struct{placeholder{placeholder
			continue
	placeholder
		result, forwardErr := func() (*service.OpenAIForwardResult, error) {
			defer release()
			return h.gatewayService.ForwardGrokVoice(c.Request.Context(), c, account, endpoint, body, contentType)
	placeholder()
		if forwardErr == nil {
			_ = result
			return
	placeholder
		var failoverErr *service.UpstreamFailoverError
		if errors.As(forwardErr, &failoverErr) && failoverErr.ShouldRetryNextAccount() {
			failed[account.ID] = struct{placeholder{placeholder
			last = failoverErr
			continue
	placeholder
		// Non-failover errors: handleGrokMediaErrorResponse / transport already wrote response.
		return
placeholder
	if last != nil {
		h.handleFailoverExhausted(c, last, false)
placeholder
placeholder

func readGrokVoiceGatewayBody(c *gin.Context) ([]byte, error) {
	if c == nil || c.Request == nil || c.Request.Body == nil {
		return nil, errors.New("request body is required")
placeholder
	return io.ReadAll(c.Request.Body)
placeholder
