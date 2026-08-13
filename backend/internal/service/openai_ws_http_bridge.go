package service

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

const (
	openAIWSClientReadLimitBytesDefault     int64 = 64 * 1024 * 1024
	openAIWSHTTPBridgeThresholdBytesDefault int64 = 15 * 1024 * 1024
	openAIWSHTTPBridgeErrorBodyLimitBytes         = 64 * 1024
)

// ResolveOpenAIWSClientFirstMessageTimeout returns the effective client ingress deadline.
func ResolveOpenAIWSClientFirstMessageTimeout(cfg *config.Config) time.Duration {
	seconds := config.DefaultOpenAIWSClientFirstMessageTimeoutSeconds
	if cfg != nil && cfg.Gateway.OpenAIWS.ClientFirstMessageTimeoutSeconds > 0 {
		seconds = cfg.Gateway.OpenAIWS.ClientFirstMessageTimeoutSeconds
placeholder
	return time.Duration(seconds) * time.Second
placeholder

func ResolveOpenAIWSClientReadLimitBytes(cfg *config.Config) int64 {
	if cfg == nil || cfg.Gateway.OpenAIWS.ClientReadLimitBytes <= 0 {
		return openAIWSClientReadLimitBytesDefault
placeholder
	return cfg.Gateway.OpenAIWS.ClientReadLimitBytes
placeholder

func (s *OpenAIGatewayService) openAIWSHTTPBridgeEnabled() bool {
	return s != nil && s.cfg != nil && s.cfg.Gateway.OpenAIWS.HTTPBridgeEnabled
placeholder

func (s *OpenAIGatewayService) openAIWSHTTPBridgeThresholdBytes() int64 {
	if s == nil || s.cfg == nil || s.cfg.Gateway.OpenAIWS.HTTPBridgeThresholdBytes <= 0 {
		return openAIWSHTTPBridgeThresholdBytesDefault
placeholder
	return s.cfg.Gateway.OpenAIWS.HTTPBridgeThresholdBytes
placeholder

func (s *OpenAIGatewayService) shouldBridgeOpenAIWSHTTP(account *Account, payloadBytes int, previousResponseID string) bool {
	if account != nil && account.Platform == PlatformGrok {
		return true
placeholder
	if !s.openAIWSHTTPBridgeEnabled() {
		return false
placeholder
	if strings.TrimSpace(previousResponseID) != "" {
		return false
placeholder
	threshold := s.openAIWSHTTPBridgeThresholdBytes()
	return threshold > 0 && int64(payloadBytes) >= threshold
placeholder

func prepareOpenAIWSHTTPBridgeBody(payload []byte) ([]byte, error) {
	var body map[string]any
	if err := json.Unmarshal(payload, &body); err != nil {
		return nil, err
placeholder
	if body == nil {
		return nil, errors.New("response.create payload must be a JSON object")
placeholder
	delete(body, "type")
	delete(body, "generate")
	delete(body, "previous_response_id")
	body["stream"] = true
	return json.Marshal(body)
placeholder

type openAIWSToolCallReplayCollector struct {
	items []json.RawMessage
	seen  map[string]struct{placeholder
placeholder

func (c *openAIWSToolCallReplayCollector) AddEvent(eventType string, message []byte) {
	switch strings.TrimSpace(eventType) {
	case "response.output_item.done":
		c.addItem(gjson.GetBytes(message, "item"))
	case "response.completed", "response.done":
		output := gjson.GetBytes(message, "response.output")
		if !output.IsArray() {
			return
	placeholder
		for _, item := range output.Array() {
			c.addItem(item)
	placeholder
placeholder
placeholder

func (c *openAIWSToolCallReplayCollector) Items() []json.RawMessage {
	return cloneOpenAIWSRawMessages(c.items)
placeholder

func (c *openAIWSToolCallReplayCollector) addItem(item gjson.Result) {
	if !item.Exists() || item.Type != gjson.JSON {
		return
placeholder
	raw := strings.TrimSpace(item.Raw)
	if raw == "" || !strings.HasPrefix(raw, "{") {
		return
placeholder
	if !isCodexToolCallContextItemType(item.Get("type").String()) {
		return
placeholder
	key := strings.TrimSpace(item.Get("id").String())
	if key == "" {
		key = strings.TrimSpace(item.Get("call_id").String())
placeholder
	if key == "" {
		key = raw
placeholder
	if c.seen == nil {
		c.seen = make(map[string]struct{placeholder)
placeholder
	if _, ok := c.seen[key]; ok {
		return
placeholder
	c.seen[key] = struct{placeholder{placeholder
	c.items = append(c.items, json.RawMessage(raw))
placeholder

func buildOpenAIWSHTTPBridgeErrorEvent(statusCode int, message string) []byte {
	message = strings.TrimSpace(message)
	if message == "" {
		message = http.StatusText(statusCode)
placeholder
	if message == "" {
		message = "upstream request failed"
placeholder
	event := map[string]any{
		"type":   "error",
		"status": statusCode,
		"error": map[string]any{
			"type":    "upstream_error",
			"message": message,
	placeholder,
placeholder
	body, err := json.Marshal(event)
	if err != nil {
		return []byte(`{"type":"error","error":{"type":"upstream_error","message":"upstream request failed"placeholderplaceholder`)
placeholder
	return body
placeholder

func (s *OpenAIGatewayService) proxyOpenAIWSHTTPBridgeTurn(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	token string,
	payload []byte,
	payloadBytes int,
	originalModel string,
	imageBillingModel string,
	imageSizeTier string,
	imageInputSize string,
	grokCacheIdentity string,
	turn int,
	writeClientMessage func([]byte) error,
) (*OpenAIForwardResult, error) {
	if s == nil {
		return nil, errors.New("service is nil")
placeholder
	if s.httpUpstream == nil {
		return nil, errors.New("openai http upstream is nil")
placeholder
	if account == nil {
		return nil, errors.New("account is nil")
placeholder
	if writeClientMessage == nil {
		return nil, errors.New("client websocket writer is nil")
placeholder
	responseModelObserver := &upstreamResponseModelObserver{placeholder

	body, err := prepareOpenAIWSHTTPBridgeBody(payload)
	if err != nil {
		return nil, fmt.Errorf("prepare http bridge body: %w", err)
placeholder

	upstreamCtx, releaseUpstreamCtx := detachUpstreamContext(ctx)
	var upstreamReq *http.Request
	if account.Platform == PlatformGrok {
		upstreamModel := resolveGrokWSUpstreamModel(account, body, originalModel)
		grokIntentSourceBody := body
		body, err = patchGrokResponsesBody(body, upstreamModel)
		if err != nil {
			releaseUpstreamCtx()
			return nil, err
	placeholder
		grokMixedCacheIntentBody := append([]byte(nil), body...)
		body, err = applyGrokResponsesCacheIdentity(body, grokIntentSourceBody, grokCacheIdentity, account.IsGrokOAuth())
		if err != nil {
			releaseUpstreamCtx()
			return nil, fmt.Errorf("apply grok prompt cache identity: %w", err)
	placeholder
		body, err = applyGrokFreeRequestToolCacheRoute(c, body, grokMixedCacheIntentBody, account, grokCacheIdentity)
		if err != nil {
			releaseUpstreamCtx()
			return nil, fmt.Errorf("apply grok Free function-tool cache route: %w", err)
	placeholder
		upstreamReq, err = buildGrokResponsesRequest(upstreamCtx, c, account, body, token, grokCacheIdentity, s.cfg, s.settingService)
placeholder else {
		upstreamReq, err = s.buildUpstreamRequestOpenAIPassthrough(upstreamCtx, c, account, body, token)
placeholder
	releaseUpstreamCtx()
	if err != nil {
		return nil, err
placeholder
	if account.Platform != PlatformGrok && isOpenAIResponsesLiteWebSocketPayload(payload) {
		upstreamReq.Header.Set(responsesLiteHeader, "true")
placeholder

	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
placeholder
	if c != nil {
		c.Set("openai_passthrough", true)
		c.Set("openai_ws_http_bridge", true)
placeholder

	turnStart := time.Now()
	resp, err := s.httpUpstream.Do(upstreamReq, proxyURL, account.ID, account.Concurrency)
	if err != nil {
		if turn == 1 {
			return nil, s.handleOpenAIUpstreamTransportError(ctx, c, account, err, true)
	placeholder
		safeErr := sanitizeUpstreamErrorMessage(err.Error())
		_ = writeClientMessage(buildOpenAIWSHTTPBridgeErrorEvent(http.StatusBadGateway, "Upstream request failed"))
		return nil, fmt.Errorf("upstream http bridge request failed: %s", safeErr)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, openAIWSHTTPBridgeErrorBodyLimitBytes))
		upstreamMsg := sanitizeUpstreamErrorMessage(strings.TrimSpace(extractUpstreamErrorMessage(respBody)))
		if upstreamMsg == "" {
			upstreamMsg = http.StatusText(resp.StatusCode)
	placeholder
		shouldFailover := s.shouldFailoverOpenAIUpstreamResponse(resp.StatusCode, upstreamMsg, respBody)
		if account.Platform == PlatformGrok {
			shouldFailover = s.shouldFailoverGrokUpstreamError(resp.StatusCode, respBody)
			s.handleGrokAccountUpstreamError(withGrokTeamRateLimitModel(ctx, resolveGrokWSUpstreamModel(account, body, originalModel)), account, resp.StatusCode, resp.Header, respBody)
			if turn == 1 && shouldFailover {
				return nil, newOpenAIUpstreamFailoverError(resp.StatusCode, resp.Header, respBody, upstreamMsg, false)
		placeholder
	placeholder else if turn == 1 && shouldFailover {
			return nil, s.handleFailoverErrorResponsePassthrough(ctx, resp, c, account, body, respBody)
	placeholder
		if account.Platform != PlatformGrok && (shouldFailover || shouldCooldownOpenAITransientUpstreamError(resp.StatusCode, respBody)) {
			canonicalModel := canonicalOpenAIAccountSchedulingModel(account, originalModel)
			s.handleOpenAIAccountUpstreamError(ctx, account, resp.StatusCode, resp.Header, respBody, canonicalModel)
	placeholder
		_ = writeClientMessage(buildOpenAIWSHTTPBridgeErrorEvent(resp.StatusCode, upstreamMsg))
		return nil, fmt.Errorf("upstream http bridge error: status=%d message=%s", resp.StatusCode, upstreamMsg)
placeholder
	if account.Platform == PlatformGrok {
		s.updateGrokUsageFromResponse(withGrokTeamRateLimitModel(ctx, resolveGrokWSUpstreamModel(account, body, originalModel)), account, resp.Header, resp.StatusCode)
placeholder

	responseID := ""
	usage := OpenAIUsage{placeholder
	imageCounter := newOpenAIImageOutputCounter()
	var firstTokenMs *int
	reqStream := openAIWSPayloadBoolFromRaw(body, "stream", true)
	eventCount := 0
	tokenEventCount := 0
	terminalEventCount := 0
	replayCollector := &openAIWSToolCallReplayCollector{placeholder
	firstEventType := ""
	lastEventType := ""
	upstreamTerminalEvent := ""
	sawDone := false
	wroteDownstream := false
	clientDisconnected := false
	mappedModel := ""
	needModelReplace := false
	var mappedModelBytes []byte
	if originalModel != "" {
		mappedModel = strings.TrimSpace(gjson.GetBytes(body, "model").String())
		if mappedModel == "" {
			mappedModel = normalizeOpenAIModelForUpstream(account, account.GetMappedModel(originalModel))
	placeholder
		needModelReplace = mappedModel != "" && mappedModel != originalModel
		if needModelReplace {
			mappedModelBytes = []byte(mappedModel)
	placeholder
placeholder

	resultWithUsage := func() *OpenAIForwardResult {
		imageCount := imageCounter.Count()
		result := &OpenAIForwardResult{
			RequestID:                     responseID,
			Usage:                         usage,
			Model:                         originalModel,
			UpstreamModel:                 mappedModel,
			UpstreamResponseModel:         responseModelObserver.Model(),
			UpstreamResponseModelConflict: responseModelObserver.Conflict(),
			ServiceTier:                   extractOpenAIServiceTierFromBody(body),
			ReasoningEffort:               ApplyThinkingEnabledFallback(extractOpenAIReasoningEffortFromBody(body, mappedModel, originalModel), body, mappedModel),
			Stream:                        reqStream,
			OpenAIWSMode:                  true,
			UpstreamTerminalEvent:         upstreamTerminalEvent,
			ResponseHeaders:               cloneHeader(resp.Header),
			Duration:                      time.Since(turnStart),
			FirstTokenMs:                  firstTokenMs,
	placeholder
		if replayInput := replayCollector.Items(); len(replayInput) > 0 {
			result.wsReplayInput = replayInput
			result.wsReplayInputExists = true
	placeholder
		if imageCount > 0 {
			result.ImageCount = imageCount
			result.ImageSize = imageSizeTier
			result.ImageInputSize = imageInputSize
			result.ImageOutputSizes = imageCounter.Sizes()
			result.BillingModel = imageBillingModel
	placeholder
		return result
placeholder

	scanner := bufio.NewScanner(resp.Body)
	maxLineSize := defaultMaxLineSize
	if s.cfg != nil && s.cfg.Gateway.MaxLineSize > 0 {
		maxLineSize = s.cfg.Gateway.MaxLineSize
placeholder
	scanBuf := getSSEScannerBuf64K()
	scanner.Buffer(scanBuf[:0], maxLineSize)
	defer putSSEScannerBuf64K(scanBuf)

	for scanner.Scan() {
		line := scanner.Text()
		data, ok := extractOpenAISSEDataLine(line)
		if !ok {
			continue
	placeholder
		trimmedData := strings.TrimSpace(data)
		if trimmedData == "" {
			continue
	placeholder
		if trimmedData == "[DONE]" {
			sawDone = true
			continue
	placeholder

		upstreamMessage := []byte(trimmedData)
		if normalized, changed := normalizeCompletedImageGenerationStatus(upstreamMessage); changed {
			upstreamMessage = normalized
	placeholder
		eventType, eventResponseID, _ := parseOpenAIWSEventEnvelope(upstreamMessage)
		responseModelObserver.ObserveOpenAI(upstreamMessage, eventType)
		if responseID == "" && eventResponseID != "" {
			responseID = eventResponseID
	placeholder
		if eventType != "" {
			eventCount++
			if firstEventType == "" {
				firstEventType = eventType
		placeholder
			lastEventType = eventType
	placeholder
		if isOpenAIWSTokenEvent(eventType) {
			tokenEventCount++
			if firstTokenMs == nil {
				ms := int(time.Since(turnStart).Milliseconds())
				firstTokenMs = &ms
		placeholder
	placeholder
		if openAIWSEventShouldParseUsage(eventType) {
			parseOpenAIWSResponseUsageFromCompletedEvent(upstreamMessage, &usage)
	placeholder
		imageCounter.AddSSEData(upstreamMessage)

		if needModelReplace && len(mappedModelBytes) > 0 && openAIWSEventMayContainModel(eventType) && strings.Contains(trimmedData, mappedModel) {
			upstreamMessage = replaceOpenAIWSMessageModel(upstreamMessage, mappedModel, originalModel)
	placeholder
		if s.toolCorrector != nil && openAIWSEventMayContainToolCalls(eventType) && openAIWSMessageLikelyContainsToolCalls(upstreamMessage) {
			if corrected, changed := s.toolCorrector.CorrectToolCallsInSSEBytes(upstreamMessage); changed {
				upstreamMessage = corrected
		placeholder
	placeholder
		replayCollector.AddEvent(eventType, upstreamMessage)

		var upstreamEventErr error
		if eventType == "error" {
			errCodeRaw, errTypeRaw, errMsgRaw := parseOpenAIWSErrorEventFields(upstreamMessage)
			errMessage := strings.TrimSpace(errMsgRaw)
			if errMessage == "" {
				errMessage = "upstream error event"
		placeholder
			statusCode := openAIWSErrorHTTPStatusFromRaw(errCodeRaw, errTypeRaw)
			shouldFailover := s.shouldFailoverOpenAIUpstreamResponse(statusCode, errMessage, upstreamMessage)
			if account.Platform == PlatformGrok {
				// SSE error events do not carry an HTTP status. The local status
				// mapper therefore defaults unknown xAI codes (for example
				// new_sensitive) to 502; classify the body as a request-scoped
				// 403 before applying status-based failover or account state.
				if isGrokContentPolicyRejection(http.StatusForbidden, upstreamMessage) {
					shouldFailover = false
			placeholder else {
					shouldFailover = s.shouldFailoverGrokUpstreamError(statusCode, upstreamMessage)
					s.handleGrokAccountUpstreamError(ctx, account, statusCode, resp.Header, upstreamMessage)
			placeholder
		placeholder else if shouldFailover {
				accountStatus := statusCode
				if transientStatus := openAIWSPayloadTransientStatus(upstreamMessage); transientStatus != 0 {
					accountStatus = transientStatus
			placeholder
				canonicalModel := canonicalOpenAIAccountSchedulingModel(account, originalModel)
				s.handleOpenAIAccountUpstreamError(ctx, account, accountStatus, resp.Header, upstreamMessage, canonicalModel)
		placeholder
			if turn == 1 && !wroteDownstream && shouldFailover {
				return nil, newOpenAIUpstreamFailoverError(statusCode, resp.Header, upstreamMessage, errMessage, false)
		placeholder
			upstreamEventErr = errors.New(errMessage)
	placeholder

		// 客户端写出副本改写容量降载码：Codex 对 error/response.failed 中的
		// server_is_overloaded / slow_down 判致命并终止会话，改写后走客户端内置
		// 重试。账号状态与终止事件判定（下方 handleOpenAIWSTerminalTransientFailure）
		// 仍使用未改写的 upstreamMessage。
		clientMessage := upstreamMessage
		if eventType == "error" || eventType == "response.failed" {
			if rewritten, changed := sanitizeOpenAICapacityShedErrorCodeForClient(clientMessage); changed {
				clientMessage = rewritten
		placeholder
	placeholder
		if !clientDisconnected {
			if err := writeClientMessage(clientMessage); err != nil {
				if isOpenAIWSClientDisconnectError(err) {
					clientDisconnected = true
					closeStatus, closeReason := summarizeOpenAIWSReadCloseError(err)
					logOpenAIWSModeInfo(
						"ingress_ws_http_bridge_client_disconnected_drain account_id=%d turn=%d close_status=%s close_reason=%s",
						account.ID,
						turn,
						closeStatus,
						truncateOpenAIWSLogValue(closeReason, openAIWSHeaderValueMaxLen),
					)
			placeholder else {
					return nil, wrapOpenAIWSIngressTurnError(
						"write_client",
						fmt.Errorf("write client websocket event: %w", err),
						wroteDownstream,
					)
			placeholder
		placeholder else {
				wroteDownstream = true
		placeholder
	placeholder

		if upstreamEventErr != nil {
			return resultWithUsage(), upstreamEventErr
	placeholder
		if isOpenAIWSTerminalEvent(eventType) {
			upstreamTerminalEvent = s.handleOpenAIWSTerminalTransientFailure(ctx, account, canonicalOpenAIAccountSchedulingModel(account, originalModel), resp.Header, upstreamMessage)
			terminalEventCount++
			firstTokenMsValue := -1
			if firstTokenMs != nil {
				firstTokenMsValue = *firstTokenMs
		placeholder
			logOpenAIWSModeInfo(
				"ingress_ws_http_bridge_turn_completed account_id=%d turn=%d response_id=%s payload_bytes=%d duration_ms=%d events=%d token_events=%d terminal_events=%d first_event=%s last_event=%s first_token_ms=%d client_disconnected=%v",
				account.ID,
				turn,
				truncateOpenAIWSLogValue(responseID, openAIWSIDValueMaxLen),
				payloadBytes,
				time.Since(turnStart).Milliseconds(),
				eventCount,
				tokenEventCount,
				terminalEventCount,
				truncateOpenAIWSLogValue(firstEventType, openAIWSLogValueMaxLen),
				truncateOpenAIWSLogValue(lastEventType, openAIWSLogValueMaxLen),
				firstTokenMsValue,
				clientDisconnected,
			)
			return resultWithUsage(), nil
	placeholder
placeholder
	if err := scanner.Err(); err != nil {
		streamErr := fmt.Errorf("read upstream http bridge stream: %w", err)
		if turn == 1 && !wroteDownstream {
			return nil, s.handleOpenAIUpstreamTransportError(ctx, c, account, streamErr, true)
	placeholder
		return resultWithUsage(), streamErr
placeholder
	terminalErr := errors.New("upstream http bridge stream ended before terminal event")
	if sawDone {
		terminalErr = errors.New("upstream http bridge stream sent [DONE] before terminal event")
placeholder
	if turn == 1 && !wroteDownstream {
		return nil, s.handleOpenAIUpstreamTransportError(ctx, c, account, terminalErr, true)
placeholder
	return resultWithUsage(), terminalErr
placeholder

func resolveGrokWSCacheIdentity(c *gin.Context, account *Account, seedPayload, currentPayload []byte, originalModel string) (string, error) {
	body, err := prepareOpenAIWSHTTPBridgeBody(seedPayload)
	if err != nil {
		return "", err
placeholder
	upstreamModel := resolveGrokWSUpstreamModel(account, currentPayload, originalModel)
	body, err = patchGrokResponsesBody(body, upstreamModel)
	if err != nil {
		return "", err
placeholder
	return resolveGrokCacheIdentity(c, body, "", upstreamModel), nil
placeholder

func resolveGrokWSUpstreamModel(account *Account, body []byte, originalModel string) string {
	upstreamModel := strings.TrimSpace(gjson.GetBytes(body, "model").String())
	originalModel = strings.TrimSpace(originalModel)
	// Shared ingress has already applied channel and account mappings when the
	// body model differs from the client-facing model. Only resolve from the
	// original model when the body still carries that original value.
	if account != nil && originalModel != "" && (upstreamModel == "" || upstreamModel == originalModel) {
		if mappedModel := normalizeOpenAIModelForUpstream(account, account.GetMappedModel(originalModel)); mappedModel != "" {
			upstreamModel = mappedModel
	placeholder
placeholder
	if upstreamModel == "" {
		upstreamModel = grokDefaultResponsesModel
placeholder
	return upstreamModel
placeholder
