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
		body, err = applyGrokResponsesCacheIdentity(body, grokIntentSourceBody, grokCacheIdentity, account.IsGrokOAuth())
		if err != nil {
			releaseUpstreamCtx()
			return nil, fmt.Errorf("apply grok prompt cache identity: %w", err)
	placeholder
		upstreamReq, err = buildGrokResponsesRequest(upstreamCtx, c, account, body, token, grokCacheIdentity, s.cfg)
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
		safeErr := sanitizeUpstreamErrorMessage(err.Error())
		_ = writeClientMessage(buildOpenAIWSHTTPBridgeErrorEvent(http.StatusBadGateway, "Upstream request failed"))
		return nil, fmt.Errorf("upstream http bridge request failed: %s", safeErr)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, openAIWSHTTPBridgeErrorBodyLimitBytes))
		if account.Platform == PlatformGrok {
			s.handleGrokAccountUpstreamError(ctx, account, resp.StatusCode, resp.Header, respBody)
	placeholder
		upstreamMsg := sanitizeUpstreamErrorMessage(strings.TrimSpace(extractUpstreamErrorMessage(respBody)))
		if upstreamMsg == "" {
			upstreamMsg = http.StatusText(resp.StatusCode)
	placeholder
		_ = writeClientMessage(buildOpenAIWSHTTPBridgeErrorEvent(resp.StatusCode, upstreamMsg))
		return nil, fmt.Errorf("upstream http bridge error: status=%d message=%s", resp.StatusCode, upstreamMsg)
placeholder
	if account.Platform == PlatformGrok {
		s.updateGrokUsageFromResponse(ctx, account, resp.Header, resp.StatusCode)
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
	sawDone := false
	wroteDownstream := false
	clientDisconnected := false
	mappedModel := ""
	needModelReplace := false
	var mappedModelBytes []byte
	if originalModel != "" {
		mappedModel = normalizeOpenAIModelForUpstream(account, account.GetMappedModel(originalModel))
		needModelReplace = mappedModel != "" && mappedModel != originalModel
		if needModelReplace {
			mappedModelBytes = []byte(mappedModel)
	placeholder
placeholder

	resultWithUsage := func() *OpenAIForwardResult {
		imageCount := imageCounter.Count()
		result := &OpenAIForwardResult{
			RequestID:       responseID,
			Usage:           usage,
			Model:           originalModel,
			UpstreamModel:   mappedModel,
			ServiceTier:     extractOpenAIServiceTierFromBody(body),
			ReasoningEffort: ApplyThinkingEnabledFallback(extractOpenAIReasoningEffortFromBody(body, mappedModel, originalModel), body, mappedModel),
			Stream:          reqStream,
			OpenAIWSMode:    true,
			ResponseHeaders: cloneHeader(resp.Header),
			Duration:        time.Since(turnStart),
			FirstTokenMs:    firstTokenMs,
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
		eventType, eventResponseID, _ := parseOpenAIWSEventEnvelope(upstreamMessage)
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

		if !clientDisconnected {
			if err := writeClientMessage(upstreamMessage); err != nil {
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

		if eventType == "error" {
			errCodeRaw, errTypeRaw, errMsgRaw := parseOpenAIWSErrorEventFields(upstreamMessage)
			s.persistOpenAIWSRateLimitSignal(ctx, account, resp.Header, upstreamMessage, errCodeRaw, errTypeRaw, errMsgRaw)
			errMessage := strings.TrimSpace(errMsgRaw)
			if errMessage == "" {
				errMessage = "upstream error event"
		placeholder
			return resultWithUsage(), errors.New(errMessage)
	placeholder
		if isOpenAIWSTerminalEvent(eventType) {
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
		return resultWithUsage(), fmt.Errorf("read upstream http bridge stream: %w", err)
placeholder
	if sawDone && eventCount > 0 {
		return resultWithUsage(), nil
placeholder
	return resultWithUsage(), errors.New("upstream http bridge stream ended before terminal event")
placeholder

func resolveGrokWSCacheIdentity(c *gin.Context, account *Account, payload []byte, originalModel string) (string, error) {
	body, err := prepareOpenAIWSHTTPBridgeBody(payload)
	if err != nil {
		return "", err
placeholder
	upstreamModel := resolveGrokWSUpstreamModel(account, body, originalModel)
	return resolveGrokCacheIdentity(c, body, "", upstreamModel), nil
placeholder

func resolveGrokWSUpstreamModel(account *Account, body []byte, originalModel string) string {
	upstreamModel := strings.TrimSpace(gjson.GetBytes(body, "model").String())
	if account != nil && originalModel != "" {
		if mappedModel := normalizeOpenAIModelForUpstream(account, account.GetMappedModel(originalModel)); mappedModel != "" {
			upstreamModel = mappedModel
	placeholder
placeholder
	if upstreamModel == "" {
		upstreamModel = grokDefaultResponsesModel
placeholder
	return upstreamModel
placeholder
