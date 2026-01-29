package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/sora"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// SoraGatewayHandler handles Sora OpenAI compatible endpoints.
type SoraGatewayHandler struct {
	gatewayService      *service.GatewayService
	soraGatewayService  *service.SoraGatewayService
	billingCacheService *service.BillingCacheService
	concurrencyHelper   *ConcurrencyHelper
	maxAccountSwitches  int
placeholder

// NewSoraGatewayHandler creates a new SoraGatewayHandler.
func NewSoraGatewayHandler(
	gatewayService *service.GatewayService,
	soraGatewayService *service.SoraGatewayService,
	concurrencyService *service.ConcurrencyService,
	billingCacheService *service.BillingCacheService,
	cfg *config.Config,
) *SoraGatewayHandler {
	pingInterval := time.Duration(0)
	maxAccountSwitches := 3
	if cfg != nil {
		pingInterval = time.Duration(cfg.Concurrency.PingInterval) * time.Second
		if cfg.Gateway.MaxAccountSwitches > 0 {
			maxAccountSwitches = cfg.Gateway.MaxAccountSwitches
	placeholder
placeholder
	return &SoraGatewayHandler{
		gatewayService:      gatewayService,
		soraGatewayService:  soraGatewayService,
		billingCacheService: billingCacheService,
		concurrencyHelper:   NewConcurrencyHelper(concurrencyService, SSEPingFormatComment, pingInterval),
		maxAccountSwitches:  maxAccountSwitches,
placeholder
placeholder

// ChatCompletions handles Sora OpenAI-compatible chat completions endpoint.
// POST /v1/chat/completions
func (h *SoraGatewayHandler) ChatCompletions(c *gin.Context) {
	apiKey, ok := middleware.GetAPIKeyFromContext(c)
	if !ok {
		h.errorResponse(c, http.StatusUnauthorized, "authentication_error", "Invalid API key")
		return
placeholder
	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok {
		h.errorResponse(c, http.StatusInternalServerError, "api_error", "User context not found")
		return
placeholder

	body, err := io.ReadAll(c.Request.Body)
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

	var reqBody map[string]any
	if err := json.Unmarshal(body, &reqBody); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Failed to parse request body")
		return
placeholder

	model, _ := reqBody["model"].(string)
	if strings.TrimSpace(model) == "" {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "model is required")
		return
placeholder
	stream, _ := reqBody["stream"].(bool)

	prompt, imageData, videoData, remixID, err := parseSoraPrompt(reqBody)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", err.Error())
		return
placeholder
	if remixID == "" {
		remixID = sora.ExtractRemixID(prompt)
placeholder
	if remixID != "" {
		prompt = strings.ReplaceAll(prompt, remixID, "")
placeholder

	if apiKey.Group != nil && apiKey.Group.Platform != service.PlatformSora {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "当前分组不支持 Sora 平台")
		return
placeholder

	streamStarted := false
	maxWait := service.CalculateMaxWait(subject.Concurrency)
	canWait, err := h.concurrencyHelper.IncrementWaitCount(c.Request.Context(), subject.UserID, maxWait)
	waitCounted := false
	if err == nil && canWait {
		waitCounted = true
placeholder
	if err == nil && !canWait {
		h.errorResponse(c, http.StatusTooManyRequests, "rate_limit_error", "Too many pending requests, please retry later")
		return
placeholder
	defer func() {
		if waitCounted {
			h.concurrencyHelper.DecrementWaitCount(c.Request.Context(), subject.UserID)
	placeholder
placeholder()

	userReleaseFunc, err := h.concurrencyHelper.AcquireUserSlotWithWait(c, subject.UserID, subject.Concurrency, stream, &streamStarted)
	if err != nil {
		h.handleConcurrencyError(c, err, "user", streamStarted)
		return
placeholder
	if waitCounted {
		h.concurrencyHelper.DecrementWaitCount(c.Request.Context(), subject.UserID)
		waitCounted = false
placeholder
	userReleaseFunc = wrapReleaseOnDone(c.Request.Context(), userReleaseFunc)
	if userReleaseFunc != nil {
		defer userReleaseFunc()
placeholder

	failedAccountIDs := make(map[int64]struct{placeholder)
	maxSwitches := h.maxAccountSwitches
	if mode := h.soraGatewayService.CallLogicMode(c.Request.Context()); strings.EqualFold(mode, "native") {
		maxSwitches = 1
placeholder

	for switchCount := 0; switchCount < maxSwitches; switchCount++ {
		selection, err := h.gatewayService.SelectAccountWithLoadAwareness(c.Request.Context(), apiKey.GroupID, "", model, failedAccountIDs, "")
		if err != nil {
			h.errorResponse(c, http.StatusServiceUnavailable, "server_error", err.Error())
			return
	placeholder
		account := selection.Account
		releaseFunc := selection.ReleaseFunc

		result, err := h.soraGatewayService.Generate(c.Request.Context(), account, service.SoraGenerationRequest{
			Model:         model,
			Prompt:        prompt,
			Image:         imageData,
			Video:         videoData,
			RemixTargetID: remixID,
			Stream:        stream,
			UserID:        subject.UserID,
	placeholder)
		if err != nil {
			// 失败路径：立即释放槽位，而非 defer
			if releaseFunc != nil {
				releaseFunc()
		placeholder

			if errors.Is(err, service.ErrSoraAccountMissingToken) || errors.Is(err, service.ErrSoraAccountNotEligible) {
				failedAccountIDs[account.ID] = struct{placeholder{placeholder
				continue
		placeholder
			h.handleStreamingAwareError(c, http.StatusBadGateway, "server_error", err.Error(), streamStarted)
			return
	placeholder

		// 成功路径：使用 defer 在函数退出时释放
		if releaseFunc != nil {
			defer releaseFunc()
	placeholder

		h.respondCompletion(c, model, result, stream)
		return
placeholder

	h.handleFailoverExhausted(c, http.StatusServiceUnavailable, streamStarted)
placeholder

func (h *SoraGatewayHandler) respondCompletion(c *gin.Context, model string, result *service.SoraGenerationResult, stream bool) {
	if result == nil {
		h.errorResponse(c, http.StatusInternalServerError, "api_error", "Empty response")
		return
placeholder
	if stream {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		first := buildSoraStreamChunk(model, "", true, "")
		if _, err := c.Writer.WriteString(first); err != nil {
			return
	placeholder
		final := buildSoraStreamChunk(model, result.Content, false, "stop")
		if _, err := c.Writer.WriteString(final); err != nil {
			return
	placeholder
		_, _ = c.Writer.WriteString("data: [DONE]\n\n")
		return
placeholder

	c.JSON(http.StatusOK, buildSoraNonStreamResponse(model, result.Content))
placeholder

func buildSoraStreamChunk(model, content string, isFirst bool, finishReason string) string {
	chunkID := fmt.Sprintf("chatcmpl-%d", time.Now().UnixMilli())
	delta := map[string]any{placeholder
	if isFirst {
		delta["role"] = "assistant"
placeholder
	if content != "" {
		delta["content"] = content
placeholder else {
		delta["content"] = nil
placeholder
	response := map[string]any{
		"id":      chunkID,
		"object":  "chat.completion.chunk",
		"created": time.Now().Unix(),
		"model":   model,
		"choices": []any{
			map[string]any{
				"index":         0,
				"delta":         delta,
				"finish_reason": finishReason,
		placeholder,
	placeholder,
placeholder
	payload, _ := json.Marshal(response)
	return "data: " + string(payload) + "\n\n"
placeholder

func buildSoraNonStreamResponse(model, content string) map[string]any {
	return map[string]any{
		"id":      fmt.Sprintf("chatcmpl-%d", time.Now().UnixMilli()),
		"object":  "chat.completion",
		"created": time.Now().Unix(),
		"model":   model,
		"choices": []any{
			map[string]any{
				"index": 0,
				"message": map[string]any{
					"role":    "assistant",
					"content": content,
			placeholder,
				"finish_reason": "stop",
		placeholder,
	placeholder,
placeholder
placeholder

func parseSoraPrompt(req map[string]any) (prompt, imageData, videoData, remixID string, err error) {
	messages, ok := req["messages"].([]any)
	if !ok || len(messages) == 0 {
		return "", "", "", "", fmt.Errorf("messages is required")
placeholder
	last := messages[len(messages)-1]
	msg, ok := last.(map[string]any)
	if !ok {
		return "", "", "", "", fmt.Errorf("invalid message format")
placeholder
	content, ok := msg["content"]
	if !ok {
		return "", "", "", "", fmt.Errorf("content is required")
placeholder

	if v, ok := req["image"].(string); ok && v != "" {
		imageData = v
placeholder
	if v, ok := req["video"].(string); ok && v != "" {
		videoData = v
placeholder
	if v, ok := req["remix_target_id"].(string); ok {
		remixID = v
placeholder

	switch value := content.(type) {
	case string:
		prompt = value
	case []any:
		for _, item := range value {
			part, ok := item.(map[string]any)
			if !ok {
				continue
		placeholder
			switch part["type"] {
			case "text":
				if text, ok := part["text"].(string); ok {
					prompt = text
			placeholder
			case "image_url":
				if image, ok := part["image_url"].(map[string]any); ok {
					if url, ok := image["url"].(string); ok {
						imageData = url
				placeholder
			placeholder
			case "video_url":
				if video, ok := part["video_url"].(map[string]any); ok {
					if url, ok := video["url"].(string); ok {
						videoData = url
				placeholder
			placeholder
		placeholder
	placeholder
	default:
		return "", "", "", "", fmt.Errorf("invalid content format")
placeholder
	if strings.TrimSpace(prompt) == "" && strings.TrimSpace(videoData) == "" {
		return "", "", "", "", fmt.Errorf("prompt is required")
placeholder
	return prompt, imageData, videoData, remixID, nil
placeholder

func looksLikeURL(value string) bool {
	trimmed := strings.ToLower(strings.TrimSpace(value))
	return strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://")
placeholder

func (h *SoraGatewayHandler) handleConcurrencyError(c *gin.Context, err error, slotType string, streamStarted bool) {
	if streamStarted {
		h.handleStreamingAwareError(c, http.StatusTooManyRequests, "rate_limit_error", err.Error(), true)
		return
placeholder
	c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()placeholder)
placeholder

func (h *SoraGatewayHandler) handleFailoverExhausted(c *gin.Context, statusCode int, streamStarted bool) {
	message := "No available Sora accounts"
	h.handleStreamingAwareError(c, statusCode, "server_error", message, streamStarted)
placeholder

func (h *SoraGatewayHandler) handleStreamingAwareError(c *gin.Context, status int, errType, message string, streamStarted bool) {
	if streamStarted {
		payload := map[string]any{"error": map[string]any{"message": message, "type": errType, "param": nil, "code": nilplaceholderplaceholder
		data, _ := json.Marshal(payload)
		_, _ = c.Writer.WriteString("data: " + string(data) + "\n\n")
		_, _ = c.Writer.WriteString("data: [DONE]\n\n")
		return
placeholder
	h.errorResponse(c, status, errType, message)
placeholder

func (h *SoraGatewayHandler) errorResponse(c *gin.Context, status int, errType, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"message": message,
			"type":    errType,
			"param":   nil,
			"code":    nil,
	placeholder,
placeholder)
placeholder
