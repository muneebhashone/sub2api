package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/apicompat"
	"github.com/Wei-Shaw/sub2api/internal/pkg/claude"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/Wei-Shaw/sub2api/internal/util/responseheaders"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ForwardAsAnthropic accepts an Anthropic Messages request body, converts it
// to OpenAI Responses API format, forwards to the OpenAI upstream, and converts
// the response back to Anthropic Messages format. This enables Claude Code
// clients to access OpenAI models through the standard /v1/messages endpoint.
func (s *OpenAIGatewayService) ForwardAsAnthropic(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	body []byte,
	promptCacheKey string,
	defaultMappedModel string,
) (*OpenAIForwardResult, error) {
	startTime := time.Now()

	// 1. Parse Anthropic request
	var anthropicReq apicompat.AnthropicRequest
	if err := json.Unmarshal(body, &anthropicReq); err != nil {
		return nil, fmt.Errorf("parse anthropic request: %w", err)
placeholder
	originalModel := anthropicReq.Model
	applyOpenAICompatModelNormalization(&anthropicReq)
	normalizedModel := anthropicReq.Model
	clientStream := anthropicReq.Stream // client's original stream preference

	// 2. Convert Anthropic → Responses
	responsesReq, err := apicompat.AnthropicToResponses(&anthropicReq)
	if err != nil {
		return nil, fmt.Errorf("convert anthropic to responses: %w", err)
placeholder

	// Upstream always uses streaming (upstream may not support sync mode).
	// The client's original preference determines the response format.
	responsesReq.Stream = true
	isStream := true

	// 2b. Handle BetaFastMode → service_tier: "priority"
	if containsBetaToken(c.GetHeader("anthropic-beta"), claude.BetaFastMode) {
		responsesReq.ServiceTier = "priority"
placeholder

	// 3. Model mapping
	billingModel := resolveOpenAIForwardModel(account, normalizedModel, defaultMappedModel)
	upstreamModel := normalizeOpenAIModelForUpstream(account, billingModel)
	responsesReq.Model = upstreamModel

	logger.L().Debug("openai messages: model mapping applied",
		zap.Int64("account_id", account.ID),
		zap.String("original_model", originalModel),
		zap.String("normalized_model", normalizedModel),
		zap.String("billing_model", billingModel),
		zap.String("upstream_model", upstreamModel),
		zap.Bool("stream", isStream),
	)

	// 4. Marshal Responses request body, then apply OAuth codex transform
	responsesBody, err := json.Marshal(responsesReq)
	if err != nil {
		return nil, fmt.Errorf("marshal responses request: %w", err)
placeholder

	if account.Type == AccountTypeOAuth {
		var reqBody map[string]any
		if err := json.Unmarshal(responsesBody, &reqBody); err != nil {
			return nil, fmt.Errorf("unmarshal for codex transform: %w", err)
	placeholder
		codexResult := applyCodexOAuthTransform(reqBody, false, false)
		forcedTemplateText := ""
		if s.cfg != nil {
			forcedTemplateText = s.cfg.Gateway.ForcedCodexInstructionsTemplate
	placeholder
		templateUpstreamModel := upstreamModel
		if codexResult.NormalizedModel != "" {
			templateUpstreamModel = codexResult.NormalizedModel
	placeholder
		existingInstructions, _ := reqBody["instructions"].(string)
		if _, err := applyForcedCodexInstructionsTemplate(reqBody, forcedTemplateText, forcedCodexInstructionsTemplateData{
			ExistingInstructions: strings.TrimSpace(existingInstructions),
			OriginalModel:        originalModel,
			NormalizedModel:      normalizedModel,
			BillingModel:         billingModel,
			UpstreamModel:        templateUpstreamModel,
	placeholder); err != nil {
			return nil, err
	placeholder
		if codexResult.NormalizedModel != "" {
			upstreamModel = codexResult.NormalizedModel
	placeholder
		if codexResult.PromptCacheKey != "" {
			promptCacheKey = codexResult.PromptCacheKey
	placeholder else if promptCacheKey != "" {
			reqBody["prompt_cache_key"] = promptCacheKey
	placeholder
		// OAuth codex transform forces stream=true upstream, so always use
		// the streaming response handler regardless of what the client asked.
		isStream = true
		responsesBody, err = json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("remarshal after codex transform: %w", err)
	placeholder
placeholder

	// For API key accounts (including OpenAI-compatible upstream gateways),
	// ensure promptCacheKey is also propagated via the request body so that
	// upstreams using the Responses API can derive a stable session identifier
	// from prompt_cache_key. This makes our Anthropic /v1/messages compatibility
	// path behave more like a native Responses client.
	if account.Type == AccountTypeAPIKey {
		if trimmedKey := strings.TrimSpace(promptCacheKey); trimmedKey != "" {
			var reqBody map[string]any
			if err := json.Unmarshal(responsesBody, &reqBody); err != nil {
				return nil, fmt.Errorf("unmarshal for prompt cache key injection: %w", err)
		placeholder
			if existing, ok := reqBody["prompt_cache_key"].(string); !ok || strings.TrimSpace(existing) == "" {
				reqBody["prompt_cache_key"] = trimmedKey
				updated, err := json.Marshal(reqBody)
				if err != nil {
					return nil, fmt.Errorf("remarshal after prompt cache key injection: %w", err)
			placeholder
				responsesBody = updated
		placeholder
	placeholder
placeholder

	// 4c. Apply OpenAI fast policy (may filter service_tier or block the request).
	// Mirrors the Claude anthropic-beta "fast-mode-2026-02-01" filter, but keyed
	// on the body-level service_tier field (priority/flex).
	updatedBody, policyErr := s.applyOpenAIFastPolicyToBody(ctx, account, upstreamModel, responsesBody)
	if policyErr != nil {
		var blocked *OpenAIFastBlockedError
		if errors.As(policyErr, &blocked) {
			writeAnthropicError(c, http.StatusForbidden, "forbidden_error", blocked.Message)
	placeholder
		return nil, policyErr
placeholder
	responsesBody = updatedBody

	// 5. Get access token
	token, _, err := s.GetAccessToken(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
placeholder

	// 6. Build upstream request
	upstreamCtx, releaseUpstreamCtx := detachUpstreamContext(ctx)
	upstreamReq, err := s.buildUpstreamRequest(upstreamCtx, c, account, responsesBody, token, isStream, promptCacheKey, false)
	releaseUpstreamCtx()
	if err != nil {
		return nil, fmt.Errorf("build upstream request: %w", err)
placeholder

	// Override session_id with a deterministic UUID derived from the isolated
	// session key, ensuring different API keys produce different upstream sessions.
	if promptCacheKey != "" {
		apiKeyID := getAPIKeyIDFromContext(c)
		upstreamReq.Header.Set("session_id", generateSessionUUID(isolateOpenAISessionID(apiKeyID, promptCacheKey)))
placeholder

	// 7. Send request
	proxyURL := ""
	if account.Proxy != nil {
		proxyURL = account.Proxy.URL()
placeholder
	resp, err := s.httpUpstream.Do(upstreamReq, proxyURL, account.ID, account.Concurrency)
	if err != nil {
		safeErr := sanitizeUpstreamErrorMessage(err.Error())
		setOpsUpstreamError(c, 0, safeErr, "")
		appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
			Platform:           account.Platform,
			AccountID:          account.ID,
			AccountName:        account.Name,
			UpstreamStatusCode: 0,
			Kind:               "request_error",
			Message:            safeErr,
	placeholder)
		writeAnthropicError(c, http.StatusBadGateway, "api_error", "Upstream request failed")
		return nil, fmt.Errorf("upstream request failed: %s", safeErr)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	// 8. Handle error response with failover
	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
		_ = resp.Body.Close()
		resp.Body = io.NopCloser(bytes.NewReader(respBody))

		upstreamMsg := strings.TrimSpace(extractUpstreamErrorMessage(respBody))
		upstreamMsg = sanitizeUpstreamErrorMessage(upstreamMsg)
		if s.shouldFailoverOpenAIUpstreamResponse(resp.StatusCode, upstreamMsg, respBody) {
			upstreamDetail := ""
			if s.cfg != nil && s.cfg.Gateway.LogUpstreamErrorBody {
				maxBytes := s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes
				if maxBytes <= 0 {
					maxBytes = 2048
			placeholder
				upstreamDetail = truncateString(string(respBody), maxBytes)
		placeholder
			appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
				Platform:           account.Platform,
				AccountID:          account.ID,
				AccountName:        account.Name,
				UpstreamStatusCode: resp.StatusCode,
				UpstreamRequestID:  resp.Header.Get("x-request-id"),
				Kind:               "failover",
				Message:            upstreamMsg,
				Detail:             upstreamDetail,
		placeholder)
			if s.rateLimitService != nil {
				s.rateLimitService.HandleUpstreamError(ctx, account, resp.StatusCode, resp.Header, respBody)
		placeholder
			return nil, &UpstreamFailoverError{
				StatusCode:             resp.StatusCode,
				ResponseBody:           respBody,
				RetryableOnSameAccount: account.IsPoolMode() && (isPoolModeRetryableStatus(resp.StatusCode) || isOpenAITransientProcessingError(resp.StatusCode, upstreamMsg, respBody)),
		placeholder
	placeholder
		// Non-failover error: return Anthropic-formatted error to client
		return s.handleAnthropicErrorResponse(resp, c, account)
placeholder

	// 9. Handle normal response
	// Upstream is always streaming; choose response format based on client preference.
	var result *OpenAIForwardResult
	var handleErr error
	if clientStream {
		result, handleErr = s.handleAnthropicStreamingResponse(resp, c, originalModel, billingModel, upstreamModel, startTime)
placeholder else {
		// Client wants JSON: buffer the streaming response and assemble a JSON reply.
		result, handleErr = s.handleAnthropicBufferedStreamingResponse(resp, c, originalModel, billingModel, upstreamModel, startTime)
placeholder

	// Propagate ServiceTier and ReasoningEffort to result for billing
	if handleErr == nil && result != nil {
		if responsesReq.ServiceTier != "" {
			st := responsesReq.ServiceTier
			result.ServiceTier = &st
	placeholder
		if responsesReq.Reasoning != nil && responsesReq.Reasoning.Effort != "" {
			re := responsesReq.Reasoning.Effort
			result.ReasoningEffort = &re
	placeholder
placeholder

	// Extract and save Codex usage snapshot from response headers (for OAuth accounts)
	if handleErr == nil && account.Type == AccountTypeOAuth {
		if snapshot := ParseCodexRateLimitHeaders(resp.Header); snapshot != nil {
			s.updateCodexUsageSnapshot(ctx, account.ID, snapshot)
	placeholder
placeholder

	return result, handleErr
placeholder

// handleAnthropicErrorResponse reads an upstream error and returns it in
// Anthropic error format.
func (s *OpenAIGatewayService) handleAnthropicErrorResponse(
	resp *http.Response,
	c *gin.Context,
	account *Account,
) (*OpenAIForwardResult, error) {
	return s.handleCompatErrorResponse(resp, c, account, writeAnthropicError)
placeholder

// handleAnthropicBufferedStreamingResponse reads all Responses SSE events from
// the upstream streaming response, finds the terminal event (response.completed
// / response.incomplete / response.failed), converts the complete response to
// Anthropic Messages JSON format, and writes it to the client.
// This is used when the client requested stream=false but the upstream is always
// streaming.
func (s *OpenAIGatewayService) handleAnthropicBufferedStreamingResponse(
	resp *http.Response,
	c *gin.Context,
	originalModel string,
	billingModel string,
	upstreamModel string,
	startTime time.Time,
) (*OpenAIForwardResult, error) {
	requestID := resp.Header.Get("x-request-id")

	finalResponse, usage, acc, err := s.readOpenAICompatBufferedTerminal(resp, "openai messages buffered", requestID)
	if err != nil {
		return nil, err
placeholder

	if finalResponse == nil {
		writeAnthropicError(c, http.StatusBadGateway, "api_error", "Upstream stream ended without a terminal response event")
		return nil, fmt.Errorf("upstream stream ended without terminal event")
placeholder

	// When the terminal event has an empty output array, reconstruct from
	// accumulated delta events so the client receives the full content.
	acc.SupplementResponseOutput(finalResponse)

	anthropicResp := apicompat.ResponsesToAnthropic(finalResponse, originalModel)

	if s.responseHeaderFilter != nil {
		responseheaders.WriteFilteredHeaders(c.Writer.Header(), resp.Header, s.responseHeaderFilter)
placeholder
	c.JSON(http.StatusOK, anthropicResp)

	return &OpenAIForwardResult{
		RequestID:     requestID,
		Usage:         usage,
		Model:         originalModel,
		BillingModel:  billingModel,
		UpstreamModel: upstreamModel,
		Stream:        false,
		Duration:      time.Since(startTime),
placeholder, nil
placeholder

func isOpenAICompatResponsesTerminalEvent(eventType string) bool {
	switch strings.TrimSpace(eventType) {
	case "response.completed", "response.done", "response.incomplete", "response.failed":
		return true
	default:
		return false
placeholder
placeholder

func isOpenAICompatDoneSentinelLine(line string) bool {
	payload, ok := extractOpenAISSEDataLine(line)
	return ok && strings.TrimSpace(payload) == "[DONE]"
placeholder

func (s *OpenAIGatewayService) readOpenAICompatBufferedTerminal(
	resp *http.Response,
	logPrefix string,
	requestID string,
) (*apicompat.ResponsesResponse, OpenAIUsage, *apicompat.BufferedResponseAccumulator, error) {
	acc := apicompat.NewBufferedResponseAccumulator()
	var usage OpenAIUsage
	if resp == nil || resp.Body == nil {
		return nil, usage, acc, errors.New("upstream response body is nil")
placeholder

	scanner := bufio.NewScanner(resp.Body)
	maxLineSize := defaultMaxLineSize
	if s.cfg != nil && s.cfg.Gateway.MaxLineSize > 0 {
		maxLineSize = s.cfg.Gateway.MaxLineSize
placeholder
	scanner.Buffer(make([]byte, 0, 64*1024), maxLineSize)

	streamInterval := time.Duration(0)
	if s.cfg != nil && s.cfg.Gateway.StreamDataIntervalTimeout > 0 {
		streamInterval = time.Duration(s.cfg.Gateway.StreamDataIntervalTimeout) * time.Second
placeholder
	var timeoutCh <-chan time.Time
	var timeoutTimer *time.Timer
	resetTimeout := func() {
		if streamInterval <= 0 {
			return
	placeholder
		if timeoutTimer == nil {
			timeoutTimer = time.NewTimer(streamInterval)
			timeoutCh = timeoutTimer.C
			return
	placeholder
		if !timeoutTimer.Stop() {
			select {
			case <-timeoutTimer.C:
			default:
		placeholder
	placeholder
		timeoutTimer.Reset(streamInterval)
placeholder
	stopTimeout := func() {
		if timeoutTimer == nil {
			return
	placeholder
		if !timeoutTimer.Stop() {
			select {
			case <-timeoutTimer.C:
			default:
		placeholder
	placeholder
placeholder
	resetTimeout()
	defer stopTimeout()

	type scanEvent struct {
		line string
		err  error
placeholder
	events := make(chan scanEvent, 16)
	done := make(chan struct{placeholder)
	go func() {
		defer close(events)
		for scanner.Scan() {
			select {
			case events <- scanEvent{line: scanner.Text()placeholder:
			case <-done:
				return
		placeholder
	placeholder
		if err := scanner.Err(); err != nil {
			select {
			case events <- scanEvent{err: errplaceholder:
			case <-done:
		placeholder
	placeholder
placeholder()
	defer close(done)

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				return nil, usage, acc, nil
		placeholder
			resetTimeout()
			if ev.err != nil {
				if !errors.Is(ev.err, context.Canceled) && !errors.Is(ev.err, context.DeadlineExceeded) {
					logger.L().Warn(logPrefix+": read error",
						zap.Error(ev.err),
						zap.String("request_id", requestID),
					)
			placeholder
				return nil, usage, acc, ev.err
		placeholder

			payload, ok := extractOpenAISSEDataLine(ev.line)
			if !ok || payload == "" {
				continue
		placeholder
			if strings.TrimSpace(payload) == "[DONE]" {
				return nil, usage, acc, nil
		placeholder

			var event apicompat.ResponsesStreamEvent
			if err := json.Unmarshal([]byte(payload), &event); err != nil {
				logger.L().Warn(logPrefix+": failed to parse event",
					zap.Error(err),
					zap.String("request_id", requestID),
				)
				continue
		placeholder

			acc.ProcessEvent(&event)

			if isOpenAICompatResponsesTerminalEvent(event.Type) && event.Response != nil {
				if event.Response.Usage != nil {
					usage = copyOpenAIUsageFromResponsesUsage(event.Response.Usage)
			placeholder
				return event.Response, usage, acc, nil
		placeholder

		case <-timeoutCh:
			_ = resp.Body.Close()
			logger.L().Warn(logPrefix+": data interval timeout",
				zap.String("request_id", requestID),
				zap.Duration("interval", streamInterval),
			)
			return nil, usage, acc, fmt.Errorf("stream data interval timeout")
	placeholder
placeholder
placeholder

// handleAnthropicStreamingResponse reads Responses SSE events from upstream,
// converts each to Anthropic SSE events, and writes them to the client.
// When StreamKeepaliveInterval is configured, it uses a goroutine + channel
// pattern to send Anthropic ping events during periods of upstream silence,
// preventing proxy/client timeout disconnections.
func (s *OpenAIGatewayService) handleAnthropicStreamingResponse(
	resp *http.Response,
	c *gin.Context,
	originalModel string,
	billingModel string,
	upstreamModel string,
	startTime time.Time,
) (*OpenAIForwardResult, error) {
	requestID := resp.Header.Get("x-request-id")

	if s.responseHeaderFilter != nil {
		responseheaders.WriteFilteredHeaders(c.Writer.Header(), resp.Header, s.responseHeaderFilter)
placeholder
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.WriteHeader(http.StatusOK)

	state := apicompat.NewResponsesEventToAnthropicState()
	state.Model = originalModel
	var usage OpenAIUsage
	var firstTokenMs *int
	firstChunk := true
	clientDisconnected := false

	scanner := bufio.NewScanner(resp.Body)
	maxLineSize := defaultMaxLineSize
	if s.cfg != nil && s.cfg.Gateway.MaxLineSize > 0 {
		maxLineSize = s.cfg.Gateway.MaxLineSize
placeholder
	scanner.Buffer(make([]byte, 0, 64*1024), maxLineSize)

	streamInterval := time.Duration(0)
	if s.cfg != nil && s.cfg.Gateway.StreamDataIntervalTimeout > 0 {
		streamInterval = time.Duration(s.cfg.Gateway.StreamDataIntervalTimeout) * time.Second
placeholder
	var intervalTicker *time.Ticker
	if streamInterval > 0 {
		intervalTicker = time.NewTicker(streamInterval)
		defer intervalTicker.Stop()
placeholder
	var intervalCh <-chan time.Time
	if intervalTicker != nil {
		intervalCh = intervalTicker.C
placeholder

	// resultWithUsage builds the final result snapshot.
	resultWithUsage := func() *OpenAIForwardResult {
		return &OpenAIForwardResult{
			RequestID:     requestID,
			Usage:         usage,
			Model:         originalModel,
			BillingModel:  billingModel,
			UpstreamModel: upstreamModel,
			Stream:        true,
			Duration:      time.Since(startTime),
			FirstTokenMs:  firstTokenMs,
	placeholder
placeholder

	// processDataLine handles a single "data: ..." SSE line from upstream.
	processDataLine := func(payload string) bool {
		if firstChunk {
			firstChunk = false
			ms := int(time.Since(startTime).Milliseconds())
			firstTokenMs = &ms
	placeholder

		var event apicompat.ResponsesStreamEvent
		if err := json.Unmarshal([]byte(payload), &event); err != nil {
			logger.L().Warn("openai messages stream: failed to parse event",
				zap.Error(err),
				zap.String("request_id", requestID),
			)
			return false
	placeholder

		// 仅按兼容转换器支持的终止事件提取 usage，避免无意扩大事件语义。
		isTerminalEvent := isOpenAICompatResponsesTerminalEvent(event.Type)
		if isTerminalEvent && event.Response != nil && event.Response.Usage != nil {
			usage = copyOpenAIUsageFromResponsesUsage(event.Response.Usage)
	placeholder

		// Convert to Anthropic events
		events := apicompat.ResponsesEventToAnthropicEvents(&event, state)
		if !clientDisconnected {
			for _, evt := range events {
				sse, err := apicompat.ResponsesAnthropicEventToSSE(evt)
				if err != nil {
					logger.L().Warn("openai messages stream: failed to marshal event",
						zap.Error(err),
						zap.String("request_id", requestID),
					)
					continue
			placeholder
				if _, err := fmt.Fprint(c.Writer, sse); err != nil {
					clientDisconnected = true
					logger.L().Info("openai messages stream: client disconnected, continuing to drain upstream for billing",
						zap.String("request_id", requestID),
					)
					break
			placeholder
		placeholder
	placeholder
		if len(events) > 0 && !clientDisconnected {
			c.Writer.Flush()
	placeholder
		return isTerminalEvent
placeholder

	// finalizeStream sends any remaining Anthropic events and returns the result.
	finalizeStream := func() (*OpenAIForwardResult, error) {
		if finalEvents := apicompat.FinalizeResponsesAnthropicStream(state); len(finalEvents) > 0 && !clientDisconnected {
			for _, evt := range finalEvents {
				sse, err := apicompat.ResponsesAnthropicEventToSSE(evt)
				if err != nil {
					continue
			placeholder
				if _, err := fmt.Fprint(c.Writer, sse); err != nil {
					clientDisconnected = true
					logger.L().Info("openai messages stream: client disconnected during final flush",
						zap.String("request_id", requestID),
					)
					break
			placeholder
		placeholder
			if !clientDisconnected {
				c.Writer.Flush()
		placeholder
	placeholder
		return resultWithUsage(), nil
placeholder

	// handleScanErr logs scanner errors if meaningful.
	handleScanErr := func(err error) {
		if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			logger.L().Warn("openai messages stream: read error",
				zap.Error(err),
				zap.String("request_id", requestID),
			)
	placeholder
placeholder
	missingTerminalErr := func() (*OpenAIForwardResult, error) {
		return resultWithUsage(), fmt.Errorf("stream usage incomplete: missing terminal event")
placeholder

	// ── Determine keepalive interval ──
	keepaliveInterval := time.Duration(0)
	if s.cfg != nil && s.cfg.Gateway.StreamKeepaliveInterval > 0 {
		keepaliveInterval = time.Duration(s.cfg.Gateway.StreamKeepaliveInterval) * time.Second
placeholder

	// ── No keepalive: fast synchronous path (no goroutine overhead) ──
	if streamInterval <= 0 && keepaliveInterval <= 0 {
		for scanner.Scan() {
			line := scanner.Text()
			payload, ok := extractOpenAISSEDataLine(line)
			if !ok {
				continue
		placeholder
			if strings.TrimSpace(payload) == "[DONE]" {
				return missingTerminalErr()
		placeholder
			if processDataLine(payload) {
				return finalizeStream()
		placeholder
	placeholder
		if err := scanner.Err(); err != nil {
			handleScanErr(err)
			return resultWithUsage(), fmt.Errorf("stream usage incomplete: %w", err)
	placeholder
		return missingTerminalErr()
placeholder

	// ── With keepalive: goroutine + channel + select ──
	type scanEvent struct {
		line string
		err  error
placeholder
	events := make(chan scanEvent, 16)
	done := make(chan struct{placeholder)
	var lastReadAt int64
	atomic.StoreInt64(&lastReadAt, time.Now().UnixNano())
	sendEvent := func(ev scanEvent) bool {
		select {
		case events <- ev:
			return true
		case <-done:
			return false
	placeholder
placeholder
	go func() {
		defer close(events)
		for scanner.Scan() {
			atomic.StoreInt64(&lastReadAt, time.Now().UnixNano())
			if !sendEvent(scanEvent{line: scanner.Text()placeholder) {
				return
		placeholder
	placeholder
		if err := scanner.Err(); err != nil {
			_ = sendEvent(scanEvent{err: errplaceholder)
	placeholder
placeholder()
	defer close(done)

	var keepaliveTicker *time.Ticker
	if keepaliveInterval > 0 {
		keepaliveTicker = time.NewTicker(keepaliveInterval)
		defer keepaliveTicker.Stop()
placeholder
	var keepaliveCh <-chan time.Time
	if keepaliveTicker != nil {
		keepaliveCh = keepaliveTicker.C
placeholder
	lastDataAt := time.Now()

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				// Upstream closed
				return missingTerminalErr()
		placeholder
			if ev.err != nil {
				handleScanErr(ev.err)
				return resultWithUsage(), fmt.Errorf("stream usage incomplete: %w", ev.err)
		placeholder
			lastDataAt = time.Now()
			line := ev.line
			payload, ok := extractOpenAISSEDataLine(line)
			if !ok {
				continue
		placeholder
			if strings.TrimSpace(payload) == "[DONE]" {
				return missingTerminalErr()
		placeholder
			if processDataLine(payload) {
				return finalizeStream()
		placeholder

		case <-intervalCh:
			lastRead := time.Unix(0, atomic.LoadInt64(&lastReadAt))
			if time.Since(lastRead) < streamInterval {
				continue
		placeholder
			if clientDisconnected {
				return resultWithUsage(), fmt.Errorf("stream usage incomplete after timeout")
		placeholder
			logger.L().Warn("openai messages stream: data interval timeout",
				zap.String("request_id", requestID),
				zap.String("model", originalModel),
				zap.Duration("interval", streamInterval),
			)
			return resultWithUsage(), fmt.Errorf("stream data interval timeout")

		case <-keepaliveCh:
			if clientDisconnected {
				continue
		placeholder
			if time.Since(lastDataAt) < keepaliveInterval {
				continue
		placeholder
			// Send Anthropic-format ping event
			if _, err := fmt.Fprint(c.Writer, "event: ping\ndata: {\"type\":\"ping\"placeholder\n\n"); err != nil {
				// Client disconnected
				logger.L().Info("openai messages stream: client disconnected during keepalive",
					zap.String("request_id", requestID),
				)
				clientDisconnected = true
				continue
		placeholder
			c.Writer.Flush()
	placeholder
placeholder
placeholder

// writeAnthropicError writes an error response in Anthropic Messages API format.
func writeAnthropicError(c *gin.Context, statusCode int, errType, message string) {
	c.JSON(statusCode, gin.H{
		"type": "error",
		"error": gin.H{
			"type":    errType,
			"message": message,
	placeholder,
placeholder)
placeholder

func copyOpenAIUsageFromResponsesUsage(usage *apicompat.ResponsesUsage) OpenAIUsage {
	if usage == nil {
		return OpenAIUsage{placeholder
placeholder
	result := OpenAIUsage{
		InputTokens:  usage.InputTokens,
		OutputTokens: usage.OutputTokens,
placeholder
	if usage.InputTokensDetails != nil {
		result.CacheReadInputTokens = usage.InputTokensDetails.CachedTokens
placeholder
	return result
placeholder
