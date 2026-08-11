package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/apicompat"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/Wei-Shaw/sub2api/internal/util/responseheaders"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// openaiStreamingResult streaming response result
type openaiStreamingResult struct {
	usage            *OpenAIUsage
	firstTokenMs     *int
	responseID       string
	imageCount       int
	imageOutputSizes []string
placeholder

type openaiNonStreamingResult struct {
	*OpenAIUsage
	usage            *OpenAIUsage
	responseID       string
	imageCount       int
	imageOutputSizes []string
placeholder

func (s *OpenAIGatewayService) handleStreamingResponse(ctx context.Context, resp *http.Response, c *gin.Context, account *Account, startTime time.Time, originalModel, mappedModel string) (*openaiStreamingResult, error) {
	return s.handleStreamingResponseWithReasoning(ctx, resp, c, account, startTime, originalModel, mappedModel, "")
placeholder

func (s *OpenAIGatewayService) handleStreamingResponseWithReasoning(ctx context.Context, resp *http.Response, c *gin.Context, account *Account, startTime time.Time, originalModel, mappedModel, reasoningEffort string) (*openaiStreamingResult, error) {
	firstOutputTimeout := time.Duration(0)
	if account != nil && account.Platform == PlatformOpenAI {
		firstOutputTimeout = s.openAIFirstOutputTimeout(reasoningEffort)
placeholder
	guardFirstOutput := firstOutputTimeout > 0
	var attemptResponseHeaders http.Header
	if guardFirstOutput {
		if s.responseHeaderFilter != nil {
			attemptResponseHeaders = responseheaders.FilterHeaders(resp.Header, s.responseHeaderFilter)
	placeholder else if requestID := strings.TrimSpace(resp.Header.Get("x-request-id")); requestID != "" {
			attemptResponseHeaders = http.Header{"X-Request-Id": []string{requestIDplaceholderplaceholder
	placeholder
placeholder else if s.responseHeaderFilter != nil {
		responseheaders.WriteFilteredHeaders(c.Writer.Header(), resp.Header, s.responseHeaderFilter)
placeholder

	// Set SSE response headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	// Pass through other headers
	if !guardFirstOutput && resp.Header.Get("x-request-id") != "" {
		v := resp.Header.Get("x-request-id")
		c.Header("x-request-id", v)
placeholder
	applyAttemptResponseHeaders := func() {
		if !guardFirstOutput || len(attemptResponseHeaders) == 0 || c.Writer.Written() {
			return
	placeholder
		for key, values := range attemptResponseHeaders {
			for _, value := range values {
				c.Writer.Header().Add(key, value)
		placeholder
	placeholder
		// These headers describe this gateway's SSE stream and are stable across
		// account attempts. Keep them authoritative over upstream values.
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("X-Accel-Buffering", "no")
placeholder

	w := c.Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, errors.New("streaming not supported")
placeholder
	maxLineSize := defaultMaxLineSize
	if s.cfg != nil && s.cfg.Gateway.MaxLineSize > 0 {
		maxLineSize = s.cfg.Gateway.MaxLineSize
placeholder
	var firstTokenMs *int
	bufferedWriter := bufio.NewWriterSize(w, 4*1024)
	var firstOutputStage *openAIFirstOutputStage
	if guardFirstOutput {
		firstOutputStage = newDefaultOpenAIFirstOutputStage()
		defer func() {
			if err := firstOutputStage.Close(); err != nil {
				logger.LegacyPrintf("service.openai_gateway", "OpenAI first-output staging cleanup failed: account=%d model=%s error=%v", account.ID, originalModel, err)
		placeholder
	placeholder()
placeholder
	writePendingString := func(value string) (int, error) {
		if firstOutputStage != nil && firstTokenMs == nil && !firstOutputStage.closed {
			return firstOutputStage.WriteString(value)
	placeholder
		return bufferedWriter.WriteString(value)
placeholder
	pendingBytes := func() int64 {
		if firstOutputStage != nil && firstTokenMs == nil && !firstOutputStage.closed {
			return firstOutputStage.Buffered()
	placeholder
		return int64(bufferedWriter.Buffered())
placeholder
	flushBuffered := func() error {
		if firstOutputStage != nil && firstTokenMs == nil && !firstOutputStage.closed {
			if err := firstOutputStage.CommitTo(w); err != nil {
				return err
		placeholder
	placeholder else {
			if err := bufferedWriter.Flush(); err != nil {
				return err
		placeholder
	placeholder
		flusher.Flush()
		return nil
placeholder

	usage := &OpenAIUsage{placeholder
	imageCounter := newOpenAIImageOutputCounter()
	responseID := ""
	var firstOutputScanGuard atomic.Bool
	firstOutputScanGuard.Store(guardFirstOutput)
	scanner := bufio.NewScanner(resp.Body)
	scanBuf := getSSEScannerBuf64K()
	scanner.Buffer(scanBuf[:0], maxLineSize)
	if guardFirstOutput {
		scanner.Split(openAIFirstOutputDynamicScanLines(&firstOutputScanGuard))
placeholder
	documentScanner := newOpenAISSEJSONDocumentScanner(scanner)

	streamInterval := time.Duration(0)
	if s.cfg != nil && s.cfg.Gateway.StreamDataIntervalTimeout > 0 {
		streamInterval = time.Duration(s.cfg.Gateway.StreamDataIntervalTimeout) * time.Second
placeholder
	// 仅监控上游数据间隔超时，不被下游写入阻塞影响
	var intervalTicker *time.Ticker
	if streamInterval > 0 {
		intervalTicker = time.NewTicker(streamInterval)
		defer intervalTicker.Stop()
placeholder
	var intervalCh <-chan time.Time
	if intervalTicker != nil {
		intervalCh = intervalTicker.C
placeholder

	keepaliveInterval := time.Duration(0)
	if s.cfg != nil && s.cfg.Gateway.StreamKeepaliveInterval > 0 {
		keepaliveInterval = time.Duration(s.cfg.Gateway.StreamKeepaliveInterval) * time.Second
placeholder
	// 下游 keepalive 仅用于防止代理空闲断开
	var keepaliveTicker *time.Ticker
	if keepaliveInterval > 0 {
		keepaliveTicker = time.NewTicker(keepaliveInterval)
		defer keepaliveTicker.Stop()
placeholder
	var keepaliveCh <-chan time.Time
	if keepaliveTicker != nil {
		keepaliveCh = keepaliveTicker.C
placeholder

	var firstOutputTimer *time.Timer
	var firstOutputCh <-chan time.Time
	if firstOutputTimeout > 0 {
		remaining := time.Until(startTime.Add(firstOutputTimeout))
		if remaining <= 0 {
			remaining = time.Nanosecond
	placeholder
		firstOutputTimer = time.NewTimer(remaining)
		firstOutputCh = firstOutputTimer.C
		defer firstOutputTimer.Stop()
placeholder
	stopFirstOutputTimer := func() {
		if firstOutputTimer == nil {
			return
	placeholder
		if !firstOutputTimer.Stop() {
			select {
			case <-firstOutputTimer.C:
			default:
		placeholder
	placeholder
		firstOutputTimer = nil
		firstOutputCh = nil
placeholder
	// Track downstream writes separately from upstream reads: pre-output failover
	// can buffer response.created / response.in_progress, so keepalive must be
	// based on downstream idle time.
	lastDownstreamWriteAt := time.Now()

	// 仅发送一次错误事件，避免多次写入导致协议混乱。
	// 注意：OpenAI `/v1/responses` streaming 事件必须符合 OpenAI Responses schema；
	// 否则下游 SDK（例如 OpenCode）会因为类型校验失败而报错。
	errorEventSent := false
	clientDisconnected := false // 客户端断开后继续 drain 上游以收集 usage
	sawTerminalEvent := false
	sawFailedEvent := false
	failedMessage := ""
	clientOutputStarted := false
	upstreamRequestID := strings.TrimSpace(resp.Header.Get("x-request-id"))
	var streamEarlyErr error
	eventInProgress := false
	eventStartsClientOutput := false
	eventShouldFlush := false
	handlePendingWriteError := func(err error) {
		if firstOutputStage != nil && firstTokenMs == nil && !firstOutputStage.closed {
			message := "OpenAI first-output staging failed"
			if errors.Is(err, errOpenAIFirstOutputStageLimit) {
				message = "OpenAI first-output staging limit exceeded"
		placeholder
			logger.LegacyPrintf("service.openai_gateway", "%s: account=%d model=%s error=%v", message, account.ID, originalModel, err)
			failoverErr := s.newOpenAIStreamFailoverError(c, account, false, upstreamRequestID, nil, message)
			failoverErr.SafeToFailoverAfterWrite = true
			streamEarlyErr = failoverErr
			_ = resp.Body.Close()
			return
	placeholder
		clientDisconnected = true
		logger.LegacyPrintf("service.openai_gateway", "Client disconnected during streaming, continuing to drain upstream for billing")
placeholder
	completeGuardedEvent := func(queueDrained bool) {
		completedSemanticEvent := eventStartsClientOutput
		shouldFlush := eventShouldFlush || (queueDrained && clientOutputStarted)
		eventInProgress = false
		if !clientDisconnected {
			if completedSemanticEvent {
				applyAttemptResponseHeaders()
		placeholder
			if shouldFlush {
				if err := flushBuffered(); err != nil {
					clientDisconnected = true
					logger.LegacyPrintf("service.openai_gateway", "Client disconnected during streaming flush, continuing to drain upstream for billing")
			placeholder else {
					clientOutputStarted = true
					lastDownstreamWriteAt = time.Now()
			placeholder
		placeholder
	placeholder
		if completedSemanticEvent && firstTokenMs == nil {
			firstOutputScanGuard.Store(false)
			ms := int(time.Since(startTime).Milliseconds())
			firstTokenMs = &ms
			stopFirstOutputTimer()
	placeholder
		eventStartsClientOutput = false
		eventShouldFlush = false
placeholder
	sendErrorEvent := func(reason string) {
		if errorEventSent || clientDisconnected {
			return
	placeholder
		errorEventSent = true
		payload := `{"type":"error","sequence_number":0,"error":{"type":"upstream_error","message":` + strconv.Quote(reason) + `,"code":` + strconv.Quote(reason) + `placeholderplaceholder`
		if err := flushBuffered(); err != nil {
			clientDisconnected = true
			return
	placeholder
		if _, err := writePendingString("data: " + payload + "\n\n"); err != nil {
			clientDisconnected = true
			return
	placeholder
		if err := flushBuffered(); err != nil {
			clientDisconnected = true
			return
	placeholder
		clientOutputStarted = true
		lastDownstreamWriteAt = time.Now()
placeholder

	needModelReplace := originalModel != mappedModel
	streamOutputAccumulator := apicompat.NewBufferedResponseAccumulator()
	streamImageOutputs := make([]json.RawMessage, 0, 1)
	streamSeenImages := make(map[string]struct{placeholder)
	resultWithUsage := func() *openaiStreamingResult {
		return &openaiStreamingResult{
			usage:            usage,
			firstTokenMs:     firstTokenMs,
			responseID:       responseID,
			imageCount:       imageCounter.Count(),
			imageOutputSizes: imageCounter.Sizes(),
	placeholder
placeholder
	flushPending := func(disconnectMessage string) {
		if clientDisconnected || pendingBytes() == 0 {
			return
	placeholder
		if err := flushBuffered(); err != nil {
			clientDisconnected = true
			logger.LegacyPrintf("service.openai_gateway", "%s", disconnectMessage)
			return
	placeholder
		clientOutputStarted = true
		lastDownstreamWriteAt = time.Now()
placeholder
	finalizeStream := func() (*openaiStreamingResult, error) {
		if guardFirstOutput && eventInProgress {
			// EOF dispatches the final SSE event even without a trailing blank line.
			completeGuardedEvent(true)
	placeholder
		if sawTerminalEvent && !sawFailedEvent {
			s.clearOpenAIProxyStreamDisconnect(account)
	placeholder
		if !sawTerminalEvent && !openAIStreamClientOutputStarted(c, clientOutputStarted) && !eventShouldFlush {
			return resultWithUsage(), s.newOpenAIStreamFailoverError(
				c,
				account,
				false,
				upstreamRequestID,
				nil,
				"OpenAI stream ended before a terminal event",
			)
	placeholder
		flushPending("Client disconnected during final flush, returning collected usage")
		if !sawTerminalEvent {
			if openAIStreamClientOutputStarted(c, clientOutputStarted) && !clientDisconnected {
				s.recordOpenAIProxyStreamDisconnect(account, errors.New("stream ended before terminal event"), upstreamRequestID)
		placeholder
			return resultWithUsage(), fmt.Errorf("stream usage incomplete: missing terminal event")
	placeholder
		if sawFailedEvent {
			return resultWithUsage(), fmt.Errorf("upstream response failed: %s", failedMessage)
	placeholder
		return resultWithUsage(), nil
placeholder
	handleScanErr := func(scanErr error) (*openaiStreamingResult, error, bool) {
		if scanErr == nil {
			return nil, nil, false
	placeholder
		if errors.Is(scanErr, errOpenAIFirstOutputScannerLimit) && firstTokenMs == nil {
			logger.LegacyPrintf("service.openai_gateway", "SSE token exceeded guarded first-output limit: account=%d limit=%d error=%v", account.ID, openAIFirstOutputStageMaxBytes+openAIFirstOutputScannerFramingAllowance, scanErr)
			failoverErr := s.newOpenAIStreamFailoverError(
				c, account, false, upstreamRequestID, nil,
				"OpenAI SSE line exceeds guarded first-output limit",
			)
			failoverErr.SafeToFailoverAfterWrite = true
			return resultWithUsage(), failoverErr, true
	placeholder
		if errors.Is(scanErr, bufio.ErrTooLong) && guardFirstOutput && firstTokenMs == nil {
			logger.LegacyPrintf("service.openai_gateway", "SSE line too long before first output: account=%d max_size=%d error=%v", account.ID, maxLineSize, scanErr)
			failoverErr := s.newOpenAIStreamFailoverError(
				c, account, false, upstreamRequestID, nil,
				"OpenAI SSE line exceeds guarded first-output limit",
			)
			failoverErr.SafeToFailoverAfterWrite = true
			return resultWithUsage(), failoverErr, true
	placeholder
		if sawTerminalEvent {
			if !sawFailedEvent {
				s.clearOpenAIProxyStreamDisconnect(account)
				logger.LegacyPrintf("service.openai_gateway", "Upstream scan ended after terminal event: %v", scanErr)
		placeholder
			result, err := finalizeStream()
			return result, err, true
	placeholder
		// 客户端断开/取消请求时，上游读取往往会返回 context canceled。
		// /v1/responses 的 SSE 事件必须符合 OpenAI 协议；这里不注入自定义 error event，避免下游 SDK 解析失败。
		if errors.Is(scanErr, context.Canceled) || errors.Is(scanErr, context.DeadlineExceeded) {
			if eventShouldFlush {
				flushPending("Client disconnected during canceled stream flush, returning collected usage")
		placeholder
			return resultWithUsage(), fmt.Errorf("stream usage incomplete: %w", scanErr), true
	placeholder
		if errors.Is(scanErr, bufio.ErrTooLong) {
			logger.LegacyPrintf("service.openai_gateway", "SSE line too long: account=%d max_size=%d error=%v", account.ID, maxLineSize, scanErr)
			sendErrorEvent("response_too_large")
			return resultWithUsage(), scanErr, true
	placeholder
		if !openAIStreamClientOutputStarted(c, clientOutputStarted) && !eventShouldFlush {
			msg := "OpenAI stream disconnected before completion"
			if errText := strings.TrimSpace(scanErr.Error()); errText != "" {
				msg += ": " + errText
		placeholder
			return resultWithUsage(), s.newOpenAIStreamFailoverError(c, account, false, upstreamRequestID, nil, msg), true
	placeholder
		// 客户端已断开时，上游出错仅影响体验，不影响计费；返回已收集 usage
		if clientDisconnected {
			return resultWithUsage(), fmt.Errorf("stream usage incomplete after disconnect: %w", scanErr), true
	placeholder
		s.recordOpenAIProxyStreamDisconnect(account, scanErr, upstreamRequestID)
		sendErrorEvent("stream_read_error")
		return resultWithUsage(), fmt.Errorf("stream read error: %w", scanErr), true
placeholder
	processSSELine := func(line string, queueDrained bool) {
		if streamEarlyErr != nil {
			return
	placeholder
		// Extract data from SSE line (supports both "data: " and "data:" formats)
		if data, ok := extractOpenAISSEDataLine(line); ok {
			dataBytes := []byte(data)
			eventTypeRaw := gjson.GetBytes(dataBytes, "type").String()
			eventType := strings.TrimSpace(eventTypeRaw)
			// 初始上游 data 的 type 只解析一次：原始值保持终止事件的精确匹配，规范化值供后续分支复用。
			if openAIStreamEventIsTerminalWithType(data, eventTypeRaw) {
				sawTerminalEvent = true
		placeholder
			if responseID == "" {
				responseID = extractOpenAIResponseIDFromJSONBytes(dataBytes)
		placeholder
			forceFlushFailedEvent := false
			if eventType == "response.failed" {
				failedMessage = extractOpenAISSEErrorMessage(dataBytes)
				// response.failed 自带上游已消耗的 usage（input token 通常已扣）；必须先解析
				// 再打 cyber 标记，否则 mark 记到的是解析前的 0，导致流式 cyber 按 0 token 计费
				// 而漏记真实用量。对齐 WS V2 / Chat 流式路径（均先解析 usage 再 Mark）。
				s.parseSSEUsageBytes(dataBytes, usage)
				if hit, code, msg := detectOpenAICyberPolicy(dataBytes); hit {
					MarkOpsCyberPolicy(c, CyberPolicyMark{
						Code:           code,
						Message:        msg,
						Body:           truncateString(string(dataBytes), 4096),
						UpstreamStatus: http.StatusOK,
						UpstreamInTok:  usage.InputTokens,
						UpstreamOutTok: usage.OutputTokens,
				placeholder)
			placeholder
				if !openAIStreamClientOutputStarted(c, clientOutputStarted) {
					if status, errType, errMsg, matched := applyOpenAIStreamFailedErrorPassthroughRule(c, account.Platform, dataBytes, failedMessage); matched {
						sawFailedEvent = true
						// 命中透传规则也要记录 ops 上游错误事件（对齐 CC/Messages 与
						// antigravity 先例），否则透传命中的 failed 在监控中不可见。
						s.recordOpenAIStreamUpstreamError(c, account, false, upstreamRequestID, "http_error", dataBytes, failedMessage)
						MarkResponseCommitted(c)
						c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
						c.JSON(status, gin.H{
							"error": gin.H{
								"type":    errType,
								"message": errMsg,
						placeholder,
					placeholder)
						streamEarlyErr = fmt.Errorf("upstream response failed: passthrough rule matched message=%s", errMsg)
						return
				placeholder
					if openAIStreamFailedEventShouldFailover(dataBytes, failedMessage) {
						sawFailedEvent = true
						streamEarlyErr = s.newOpenAIStreamFailoverError(c, account, false, upstreamRequestID, dataBytes, failedMessage, resp.Header)
						return
				placeholder
			placeholder
				forceFlushFailedEvent = true
				sawFailedEvent = true
		placeholder
			if normalizedData, normalized := normalizeCompletedImageGenerationStatus(dataBytes); normalized {
				dataBytes = normalizedData
				data = string(normalizedData)
				line = "data: " + data
		placeholder
			imageCounter.AddSSEData(dataBytes)

			// Correct Codex tool calls if needed (apply_patch -> edit, etc.)
			if correctedData, corrected := s.toolCorrector.CorrectToolCallsInSSEBytes(dataBytes); corrected {
				dataBytes = correctedData
				data = string(correctedData)
				line = "data: " + data
				eventType = strings.TrimSpace(gjson.GetBytes(dataBytes, "type").String())
		placeholder
			if imageOutput, ok := extractImageGenerationOutputFromSSEData(dataBytes, streamSeenImages); ok {
				streamImageOutputs = append(streamImageOutputs, imageOutput)
		placeholder
			if responsesStreamEventMayContributeToOutput(eventType) {
				var streamEvent apicompat.ResponsesStreamEvent
				if err := json.Unmarshal(dataBytes, &streamEvent); err == nil {
					streamOutputAccumulator.ProcessEvent(&streamEvent)
			placeholder
		placeholder
			if normalizedData, normalized := normalizeResponsesStreamingTerminalOutput(dataBytes, streamOutputAccumulator, streamImageOutputs); normalized {
				dataBytes = normalizedData
				data = string(normalizedData)
				line = "data: " + data
				eventType = strings.TrimSpace(gjson.GetBytes(dataBytes, "type").String())
		placeholder
			restoredData, restoreErr := restoreGrokResponsesClientToolPayload(c, dataBytes)
			if restoreErr != nil {
				streamEarlyErr = fmt.Errorf("restore Grok Responses client tool response: %w", restoreErr)
				return
		placeholder
			restoredData, restoreErr = restoreOpenAIResponsesNamespacePayload(c, restoredData)
			if restoreErr != nil {
				streamEarlyErr = fmt.Errorf("restore OpenAI namespace response: %w", restoreErr)
				return
		placeholder
			if !bytes.Equal(restoredData, dataBytes) {
				dataBytes = restoredData
				data = string(restoredData)
				line = "data: " + data
				eventType = strings.TrimSpace(gjson.GetBytes(dataBytes, "type").String())
		placeholder
			if sanitizedData, sanitized := sanitizeOpenAIResponseFailedEventForClient(
				dataBytes,
				eventType,
				openAIStreamClientOutputStarted(c, clientOutputStarted),
			); sanitized {
				dataBytes = sanitizedData
				data = string(sanitizedData)
				line = "data: " + data
		placeholder
			// Replace model in response if needed.
			// Fast path: most events do not contain model field values.
			if needModelReplace && mappedModel != "" && strings.Contains(line, mappedModel) {
				line = s.replaceModelInSSELine(line, mappedModel, originalModel)
		placeholder
			startsClientOutput := forceFlushFailedEvent || openAIStreamDataStartsClientOutput(data, eventType)
			if guardFirstOutput {
				eventStartsClientOutput = eventStartsClientOutput || startsClientOutput
		placeholder

			// 写入客户端（客户端断开后继续 drain 上游）
			if !clientDisconnected {
				shouldFlush := queueDrained && (clientOutputStarted || startsClientOutput)
				if firstTokenMs == nil && startsClientOutput {
					// 保证首个 token 事件尽快出站，避免影响 TTFT。
					shouldFlush = true
			placeholder
				eventShouldFlush = eventShouldFlush || shouldFlush
				if _, err := writePendingString(line); err != nil {
					handlePendingWriteError(err)
			placeholder else if _, err := writePendingString("\n"); err != nil {
					handlePendingWriteError(err)
			placeholder else {
					eventInProgress = true
			placeholder
		placeholder

			// Record first token time
			if !guardFirstOutput && firstTokenMs == nil && startsClientOutput {
				ms := int(time.Since(startTime).Milliseconds())
				firstTokenMs = &ms
				stopFirstOutputTimer()
		placeholder
			s.parseSSEUsageBytes(dataBytes, usage)
			return
	placeholder

		// A blank line dispatches a guarded event from the attempt-local stage.
		if guardFirstOutput && line == "" {
			if !clientDisconnected {
				if _, err := writePendingString("\n"); err != nil {
					handlePendingWriteError(err)
			placeholder
		placeholder
			if streamEarlyErr == nil {
				completeGuardedEvent(queueDrained)
		placeholder
			return
	placeholder
		// Non-guarded streams retain upstream's event-boundary flushing: a keepalive
		// or queue-drain flush must never split an open SSE event.
		shouldFlush := false
		if line == "" {
			shouldFlush = eventShouldFlush || (queueDrained && clientOutputStarted)
			eventShouldFlush = false
	placeholder
		if !clientDisconnected {
			if _, err := writePendingString(line); err != nil {
				handlePendingWriteError(err)
		placeholder else if _, err := writePendingString("\n"); err != nil {
				handlePendingWriteError(err)
		placeholder else {
				eventInProgress = line != ""
				if shouldFlush {
					if err := flushBuffered(); err != nil {
						clientDisconnected = true
						logger.LegacyPrintf("service.openai_gateway", "Client disconnected during streaming flush, continuing to drain upstream for billing")
				placeholder else {
						clientOutputStarted = true
						lastDownstreamWriteAt = time.Now()
				placeholder
			placeholder
		placeholder
	placeholder
placeholder

	// 无超时/无 keepalive 的常见路径走同步扫描，减少 goroutine 与 channel 开销。
	if streamInterval <= 0 && keepaliveInterval <= 0 && firstOutputTimeout <= 0 {
		defer putSSEScannerBuf64K(scanBuf)
		for documentScanner.Scan() {
			processSSELine(documentScanner.Text(), true)
			if streamEarlyErr != nil {
				return resultWithUsage(), streamEarlyErr
		placeholder
	placeholder
		if result, err, done := handleScanErr(documentScanner.Err()); done {
			return result, err
	placeholder
		return finalizeStream()
placeholder

	type scanEvent struct {
		line      string
		err       error
		processed chan struct{placeholder
placeholder
	// 独立 goroutine 读取上游，避免读取阻塞影响 keepalive/超时处理
	// Guard mode permits one queued token plus the token being processed. With
	// the guarded scanner cap this bounds scanner/channel retention near 16 MiB;
	// the timeout-disabled path preserves the legacy depth of 16.
	events := make(chan scanEvent, openAIFirstOutputEventQueueSize(guardFirstOutput))
	done := make(chan struct{placeholder)
	sendEvent := func(ev scanEvent) bool {
		if guardFirstOutput {
			ev.processed = make(chan struct{placeholder)
	placeholder
		select {
		case events <- ev:
		case <-done:
			return false
	placeholder
		if ev.processed == nil {
			return true
	placeholder
		select {
		case <-ev.processed:
			return true
		case <-done:
			return false
	placeholder
placeholder
	markEventProcessed := func(ev scanEvent) {
		if ev.processed != nil {
			close(ev.processed)
	placeholder
placeholder
	var lastReadAt int64
	atomic.StoreInt64(&lastReadAt, time.Now().UnixNano())
	go func(scanBuf *sseScannerBuf64K) {
		defer putSSEScannerBuf64K(scanBuf)
		defer close(events)
		for documentScanner.Scan() {
			atomic.StoreInt64(&lastReadAt, time.Now().UnixNano())
			if !sendEvent(scanEvent{line: documentScanner.Text()placeholder) {
				return
		placeholder
	placeholder
		if err := documentScanner.Err(); err != nil {
			_ = sendEvent(scanEvent{err: errplaceholder)
	placeholder
placeholder(scanBuf)
	defer close(done)

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				if guardFirstOutput && eventInProgress {
					// EOF dispatches the final SSE event even without a trailing blank
					// line. Do not synthesize extra bytes on the downstream wire.
					completeGuardedEvent(true)
			placeholder
				return finalizeStream()
		placeholder
			if result, err, done := handleScanErr(ev.err); done {
				markEventProcessed(ev)
				return result, err
		placeholder
			processSSELine(ev.line, len(events) == 0)
			markEventProcessed(ev)
			if streamEarlyErr != nil {
				return resultWithUsage(), streamEarlyErr
		placeholder

		case <-intervalCh:
			lastRead := time.Unix(0, atomic.LoadInt64(&lastReadAt))
			if time.Since(lastRead) < streamInterval {
				continue
		placeholder
			if clientDisconnected {
				return resultWithUsage(), fmt.Errorf("stream usage incomplete after timeout")
		placeholder
			logger.LegacyPrintf("service.openai_gateway", "Stream data interval timeout: account=%d model=%s interval=%s", account.ID, originalModel, streamInterval)
			// 处理流超时，可能标记账户为临时不可调度或错误状态
			if s.rateLimitService != nil {
				s.rateLimitService.HandleStreamTimeout(ctx, account, originalModel)
		placeholder
			sendErrorEvent("stream_timeout")
			return resultWithUsage(), fmt.Errorf("stream data interval timeout")

		case <-firstOutputCh:
			if firstTokenMs != nil {
				stopFirstOutputTimer()
				continue
		placeholder
			_ = resp.Body.Close()
			for ev := range events {
				markEventProcessed(ev)
		placeholder
			return resultWithUsage(), s.newOpenAIFirstOutputTimeoutError(
				ctx, c, account, startTime, originalModel, reasoningEffort,
				firstOutputTimeout, "semantic_output", resp.Header,
			)

		case <-keepaliveCh:
			if clientDisconnected {
				continue
		placeholder
			if eventInProgress {
				continue
		placeholder
			if time.Since(lastDownstreamWriteAt) < keepaliveInterval {
				continue
		placeholder
			if guardFirstOutput {
				// Bypass attempt-local buffered frames. The stable SSE headers may be
				// committed here, but account headers remain private until semantic output.
				if _, err := w.Write([]byte(":\n\n")); err != nil {
					clientDisconnected = true
					logger.LegacyPrintf("service.openai_gateway", "Client disconnected during streaming, continuing to drain upstream for billing")
					continue
			placeholder
				flusher.Flush()
				lastDownstreamWriteAt = time.Now()
				continue
		placeholder
			if _, err := writePendingString(":\n\n"); err != nil {
				clientDisconnected = true
				logger.LegacyPrintf("service.openai_gateway", "Client disconnected during streaming, continuing to drain upstream for billing")
				continue
		placeholder
			if err := flushBuffered(); err != nil {
				clientDisconnected = true
				logger.LegacyPrintf("service.openai_gateway", "Client disconnected during keepalive flush, continuing to drain upstream for billing")
		placeholder else {
				lastDownstreamWriteAt = time.Now()
		placeholder
	placeholder
placeholder

placeholder

// extractOpenAISSEDataLine 低开销提取 SSE `data:` 行内容。
// 兼容 `data: xxx` 与 `data:xxx` 两种格式。
func extractOpenAISSEDataLine(line string) (string, bool) {
	if !strings.HasPrefix(line, "data:") {
		return "", false
placeholder
	start := len("data:")
	for start < len(line) {
		if line[start] != ' ' && line[start] != '	' {
			break
	placeholder
		start++
placeholder
	return line[start:], true
placeholder

func extractOpenAISSEEventLine(line string) (string, bool) {
	if !strings.HasPrefix(line, "event:") {
		return "", false
placeholder
	start := len("event:")
	for start < len(line) {
		if line[start] != ' ' && line[start] != '	' {
			break
	placeholder
		start++
placeholder
	return strings.TrimSpace(line[start:]), true
placeholder

type openAICompatSSEFrame struct {
	EventType string
	Data      string
placeholder

type openAICompatSSEFrameParser struct {
	eventType string
	dataLines []string
placeholder

func (p *openAICompatSSEFrameParser) AddLine(line string) (openAICompatSSEFrame, bool) {
	if line == "" {
		return p.dispatch()
placeholder
	if strings.HasPrefix(line, ":") {
		return openAICompatSSEFrame{placeholder, false
placeholder
	if eventType, ok := extractOpenAISSEEventLine(line); ok {
		p.eventType = eventType
		return openAICompatSSEFrame{placeholder, false
placeholder
	if data, ok := extractOpenAISSEDataLine(line); ok {
		p.dataLines = append(p.dataLines, data)
placeholder
	return openAICompatSSEFrame{placeholder, false
placeholder

func (p *openAICompatSSEFrameParser) Finish() (openAICompatSSEFrame, bool) {
	return p.dispatch()
placeholder

func (p *openAICompatSSEFrameParser) dispatch() (openAICompatSSEFrame, bool) {
	frame := openAICompatSSEFrame{
		EventType: p.eventType,
		Data:      strings.Join(p.dataLines, "\n"),
placeholder
	p.eventType = ""
	p.dataLines = nil
	return frame, frame.Data != ""
placeholder

func openAICompatPayloadWithEventType(payload, eventType string) string {
	eventType = strings.TrimSpace(eventType)
	if eventType == "" || strings.TrimSpace(payload) == "" || strings.TrimSpace(payload) == "[DONE]" {
		return payload
placeholder
	if gjson.Get(payload, "type").Exists() {
		return payload
placeholder
	patched, err := sjson.Set(payload, "type", eventType)
	if err != nil {
		return payload
placeholder
	return patched
placeholder

func (s *OpenAIGatewayService) replaceModelInSSELine(line, fromModel, toModel string) string {
	data, ok := extractOpenAISSEDataLine(line)
	if !ok {
		return line
placeholder
	if data == "" || data == "[DONE]" {
		return line
placeholder

	// 使用 gjson 精确检查 model 字段，避免全量 JSON 反序列化
	if m := gjson.Get(data, "model"); m.Exists() && m.Str == fromModel {
		newData, err := sjson.Set(data, "model", toModel)
		if err != nil {
			return line
	placeholder
		return "data: " + newData
placeholder

	// 检查嵌套的 response.model 字段
	if m := gjson.Get(data, "response.model"); m.Exists() && m.Str == fromModel {
		newData, err := sjson.Set(data, "response.model", toModel)
		if err != nil {
			return line
	placeholder
		return "data: " + newData
placeholder

	return line
placeholder

// correctToolCallsInResponseBody 修正响应体中的工具调用
func (s *OpenAIGatewayService) correctToolCallsInResponseBody(body []byte) []byte {
	if len(body) == 0 {
		return body
placeholder

	updated := body
	if s != nil && s.toolCorrector != nil {
		if corrected, changed := s.toolCorrector.CorrectToolCallsInSSEBytes(updated); changed {
			updated = corrected
	placeholder
placeholder
	if normalized, changed := normalizeOpenAIResponsesFunctionCallArguments(updated); changed {
		updated = normalized
placeholder
	return updated
placeholder

func normalizeOpenAIResponsesFunctionCallArguments(data []byte) ([]byte, bool) {
	if len(bytes.TrimSpace(data)) == 0 || !bytes.Contains(data, []byte(`"arguments"`)) {
		return data, false
placeholder
	if !gjson.ValidBytes(data) {
		return data, false
placeholder

	updated := data
	changed := false
	setDedupedArgument := func(path string) {
		arg := gjson.GetBytes(updated, path)
		if !arg.Exists() || arg.Type != gjson.String {
			return
	placeholder
		deduped, ok := dedupeRepeatedJSONArgumentString(arg.Str)
		if !ok {
			return
	placeholder
		next, err := sjson.SetBytes(updated, path, deduped)
		if err != nil {
			return
	placeholder
		updated = next
		changed = true
placeholder

	eventType := strings.TrimSpace(gjson.GetBytes(updated, "type").String())
	if eventType == "response.function_call_arguments.done" {
		setDedupedArgument("arguments")
placeholder
	if itemType := strings.TrimSpace(gjson.GetBytes(updated, "item.type").String()); isResponsesFunctionCallItemType(itemType) {
		setDedupedArgument("item.arguments")
placeholder
	dedupeResponsesFunctionCallOutputArguments(updated, "response.output", setDedupedArgument)
	dedupeResponsesFunctionCallOutputArguments(updated, "output", setDedupedArgument)

	return updated, changed
placeholder

func dedupeResponsesFunctionCallOutputArguments(data []byte, outputPath string, setDedupedArgument func(string)) {
	output := gjson.GetBytes(data, outputPath)
	if !output.Exists() || !output.IsArray() {
		return
placeholder
	for i, item := range output.Array() {
		if !isResponsesFunctionCallItemType(strings.TrimSpace(item.Get("type").String())) {
			continue
	placeholder
		setDedupedArgument(outputPath + "." + strconv.Itoa(i) + ".arguments")
placeholder
placeholder

func isResponsesFunctionCallItemType(itemType string) bool {
	return itemType == "function_call" || itemType == "custom_tool_call"
placeholder

func dedupeRepeatedJSONArgumentString(arguments string) (string, bool) {
	if len(arguments) == 0 || len(arguments)%2 != 0 {
		return "", false
placeholder
	halfLen := len(arguments) / 2
	first := arguments[:halfLen]
	if first != arguments[halfLen:] {
		return "", false
placeholder
	trimmed := strings.TrimSpace(first)
	if trimmed == "" || (!strings.HasPrefix(trimmed, "{") && !strings.HasPrefix(trimmed, "[")) {
		return "", false
placeholder
	if !json.Valid([]byte(first)) {
		return "", false
placeholder
	return first, true
placeholder

func (s *OpenAIGatewayService) parseSSEUsage(data string, usage *OpenAIUsage) {
	s.parseSSEUsageBytes([]byte(data), usage)
placeholder

func (s *OpenAIGatewayService) parseSSEUsageBytes(data []byte, usage *OpenAIUsage) {
	if usage == nil || len(data) == 0 || bytes.Equal(data, []byte("[DONE]")) {
		return
placeholder
	// 选择性解析：仅在数据中包含终止事件标识时才进入字段提取。
	if len(data) < 72 {
		return
placeholder
	eventType := gjson.GetBytes(data, "type").String()
	if eventType != "response.completed" && eventType != "response.done" && eventType != "response.failed" &&
		eventType != "response.incomplete" && eventType != "response.cancelled" && eventType != "response.canceled" {
		return
placeholder

	if parsedUsage, ok := extractOpenAIUsageFromJSONBytes(data); ok {
		*usage = parsedUsage
placeholder
placeholder

func extractOpenAIUsageFromJSONBytes(body []byte) (OpenAIUsage, bool) {
	if len(body) == 0 || !gjson.ValidBytes(body) {
		return OpenAIUsage{placeholder, false
placeholder
	// 部分 OpenAI 兼容上游（例如 Cline API）会将标准响应包在 data 字段中：
	// {"data":{"choices": [...], "usage": {...placeholderplaceholder, "success":trueplaceholder。
	// 按优先级尝试原生 OpenAI Responses、兼容层 data 包装和 Responses 包装，
	// 避免同步请求能正常返回但用量被静默记录为 0。
	candidates := []struct {
		usagePath      string
		imageUsagePath string
placeholder{
		{usagePath: "usage", imageUsagePath: "tool_usage.image_gen"placeholder,
		{usagePath: "data.usage", imageUsagePath: "data.tool_usage.image_gen"placeholder,
		{usagePath: "response.usage", imageUsagePath: "response.tool_usage.image_gen"placeholder,
		{usagePath: "data.response.usage", imageUsagePath: "data.response.tool_usage.image_gen"placeholder,
placeholder
	for _, candidate := range candidates {
		if usage, ok := openAIUsageFromGJSON(gjson.GetBytes(body, candidate.usagePath)); ok {
			mergeHostedImageGenToolUsage(gjson.GetBytes(body, candidate.imageUsagePath), &usage)
			return usage, true
	placeholder
placeholder
	return OpenAIUsage{placeholder, false
placeholder

func mergeHostedImageGenToolUsage(imageGen gjson.Result, usage *OpenAIUsage) {
	if !imageGen.Exists() || !imageGen.IsObject() {
		return
placeholder
	if usage.ImageOutputTokens == 0 {
		if v := imageGen.Get("output_tokens_details.image_tokens").Int(); v > 0 {
			usage.ImageOutputTokens = int(v)
	placeholder
placeholder
	if usage.ImageInputTokens == 0 {
		if v := imageGen.Get("input_tokens_details.image_tokens").Int(); v > 0 {
			usage.ImageInputTokens = int(v)
	placeholder
placeholder
placeholder

func extractOpenAIResponseIDFromJSONBytes(body []byte) string {
	if len(body) == 0 || !gjson.ValidBytes(body) {
		return ""
placeholder
	if id := strings.TrimSpace(gjson.GetBytes(body, "id").String()); id != "" {
		return id
placeholder
	return strings.TrimSpace(gjson.GetBytes(body, "response.id").String())
placeholder

func (s *OpenAIGatewayService) bindHTTPResponseAccount(ctx context.Context, c *gin.Context, account *Account, responseID string) {
	if s == nil || account == nil || account.ID <= 0 {
		return
placeholder
	responseID = strings.TrimSpace(responseID)
	if responseID == "" {
		return
placeholder
	store := s.getOpenAIWSStateStore()
	if store == nil {
		return
placeholder
	groupID := getOpenAIGroupIDFromContext(c)
	ttl := s.openAIWSResponseStickyTTL()
	logOpenAIWSBindResponseAccountWarn(groupID, account.ID, responseID, store.BindResponseAccount(ctx, groupID, responseID, account.ID, ttl))
placeholder

func openAIUsageFromGJSON(value gjson.Result) (OpenAIUsage, bool) {
	if !value.Exists() || !value.IsObject() {
		return OpenAIUsage{placeholder, false
placeholder
	inputTokens := value.Get("input_tokens").Int()
	if inputTokens == 0 {
		inputTokens = value.Get("prompt_tokens").Int()
placeholder
	outputTokens := value.Get("output_tokens").Int()
	if outputTokens == 0 {
		outputTokens = value.Get("completion_tokens").Int()
placeholder
	cacheReadTokens := openAICacheReadTokensFromUsage(value)
	cacheCreationTokens := openAICacheCreationTokensFromUsage(value)
	imageOutputTokens := value.Get("output_tokens_details.image_tokens").Int()
	if imageOutputTokens == 0 {
		imageOutputTokens = value.Get("completion_tokens_details.image_tokens").Int()
placeholder
	// 图片输入 token（如 gpt-image-2 的 /v1/images/edits 带图请求），
	// 上游在 input_tokens_details.image_tokens 单独回传，用于图/文输入分价计费。
	// 普通文本请求该字段为 0，走原路径行为不变。
	imageInputTokens := firstPositiveGJSONInt(
		value.Get("input_tokens_details.image_tokens"),
		value.Get("prompt_tokens_details.image_tokens"),
	)
	return OpenAIUsage{
		InputTokens:              int(inputTokens),
		ImageInputTokens:         imageInputTokens,
		OutputTokens:             int(outputTokens),
		CacheCreationInputTokens: cacheCreationTokens,
		CacheReadInputTokens:     cacheReadTokens,
		ImageOutputTokens:        int(imageOutputTokens),
placeholder, true
placeholder

func openAICacheReadTokensFromUsage(value gjson.Result) int {
	for _, nested := range []gjson.Result{
		value.Get("input_tokens_details.cached_tokens"),
		value.Get("prompt_tokens_details.cached_tokens"),
placeholder {
		if nested.Exists() {
			return max(int(nested.Int()), 0)
	placeholder
placeholder

	return firstPositiveGJSONInt(
		value.Get("cache_read_input_tokens"),
		value.Get("cache_read_tokens"),
		value.Get("cached_tokens"),
	)
placeholder

func openAICacheCreationTokensFromUsage(value gjson.Result) int {
	for _, nested := range []gjson.Result{
		value.Get("input_tokens_details.cache_write_tokens"),
		value.Get("prompt_tokens_details.cache_write_tokens"),
		value.Get("input_tokens_details.cache_creation_tokens"),
		value.Get("prompt_tokens_details.cache_creation_tokens"),
placeholder {
		if nested.Exists() {
			return max(int(nested.Int()), 0)
	placeholder
placeholder

	return firstPositiveGJSONInt(
		value.Get("cache_write_tokens"),
		value.Get("cache_creation_input_tokens"),
		value.Get("cache_write_input_tokens"),
		value.Get("cache_creation_tokens"),
	)
placeholder

func (s *OpenAIGatewayService) handleNonStreamingResponse(ctx context.Context, resp *http.Response, c *gin.Context, account *Account, originalModel, mappedModel string) (*openaiNonStreamingResult, error) {
	body, err := ReadUpstreamResponseBody(resp.Body, s.cfg, c, openAITooLargeError)
	if err != nil {
		return nil, err
placeholder

	// Detect SSE responses for ALL account types via Content-Type header.
	// Some OpenAI-compatible upstreams (including other sub2api instances)
	// may return SSE even when stream=false was requested.
	if isEventStreamResponse(resp.Header) {
		return s.handleSSEToJSON(resp, c, body, originalModel, mappedModel)
placeholder
	// bodyLooksLikeSSE is a line-level heuristic: real SSE framing requires
	// "data:"/"event:" field names at the very start of a physical line. A
	// plain bytes.Contains scan would also match ordinary JSON responses
	// whose string content merely echoes the literal text "data:" or
	// "event:" (e.g. compact tool output), causing those JSON bodies to be
	// misrouted into handleSSEToJSON and lose their usage accounting.
	bodyLooksLikeSSE := bodyHasSSEFraming(body)

	// For OAuth accounts, also fall back to a body-content heuristic because
	// the upstream may omit the Content-Type header while still sending SSE.
	// This heuristic is NOT applied to API-key accounts to avoid false
	// positives on JSON responses that coincidentally contain "data:" or
	// "event:" in their text content.
	if account.Type == AccountTypeOAuth && bodyLooksLikeSSE {
		return s.handleSSEToJSON(resp, c, body, originalModel, mappedModel)
placeholder
	if account != nil && account.IsGrok() && isOpenAIResponsesCompactPath(c) {
		body, err = convertGrokResponseToOpenAICompact(body)
		if err != nil {
			return nil, fmt.Errorf("convert Grok compact response: %w", err)
	placeholder
placeholder

	usageValue, usageOK := extractOpenAIUsageFromJSONBytes(body)
	if !usageOK {
		if bodyLooksLikeSSE {
			return s.handleSSEToJSON(resp, c, body, originalModel, mappedModel)
	placeholder
		return nil, fmt.Errorf("parse response: invalid json response")
placeholder
	usage := &usageValue

	// Replace model in response if needed
	if originalModel != mappedModel {
		body = s.replaceModelInResponseBody(body, mappedModel, originalModel)
placeholder
	body, err = restoreGrokResponsesClientToolPayload(c, body)
	if err != nil {
		return nil, fmt.Errorf("restore Grok Responses client tool response: %w", err)
placeholder
	body, err = restoreOpenAIResponsesNamespacePayload(c, body)
	if err != nil {
		return nil, fmt.Errorf("restore OpenAI namespace response: %w", err)
placeholder
	responseheaders.WriteFilteredHeaders(c.Writer.Header(), resp.Header, s.responseHeaderFilter)

	contentType := "application/json"
	if s.cfg != nil && !s.cfg.Security.ResponseHeaders.Enabled {
		if upstreamType := resp.Header.Get("Content-Type"); upstreamType != "" {
			contentType = upstreamType
	placeholder
placeholder

	if !writeOpenAICompactSSEBridge(c, resp.StatusCode, body) {
		c.Data(resp.StatusCode, contentType, body)
placeholder

	return &openaiNonStreamingResult{
		OpenAIUsage:      usage,
		usage:            usage,
		responseID:       extractOpenAIResponseIDFromJSONBytes(body),
		imageCount:       countOpenAIResponseImageOutputsFromJSONBytes(body),
		imageOutputSizes: collectOpenAIResponseImageOutputSizesFromJSONBytes(body),
placeholder, nil
placeholder

func isEventStreamResponse(header http.Header) bool {
	contentType := strings.ToLower(header.Get("Content-Type"))
	return strings.Contains(contentType, "text/event-stream")
placeholder

// bodyHasSSEFraming reports whether body contains genuine SSE framing by
// scanning for physical lines that begin with the "data:" or "event:"
// field names, per the SSE spec. Unlike a raw substring scan, this does not
// match when those strings only appear embedded inside JSON string values
// (e.g. "data: foo" quoted as part of an assistant text field), since such
// occurrences never start a physical line in a valid JSON encoding.
func bodyHasSSEFraming(body []byte) bool {
	for _, line := range bytes.Split(body, []byte("\n")) {
		line = bytes.TrimRight(line, "\r")
		if bytes.HasPrefix(line, []byte("data:")) || bytes.HasPrefix(line, []byte("event:")) {
			return true
	placeholder
placeholder
	return false
placeholder

func (s *OpenAIGatewayService) handleSSEToJSON(resp *http.Response, c *gin.Context, body []byte, originalModel, mappedModel string) (*openaiNonStreamingResult, error) {
	bodyText := string(body)
	finalResponse, ok := extractCodexFinalResponse(bodyText)

	usage := &OpenAIUsage{placeholder
	if ok {
		if parsedUsage, parsed := extractOpenAIUsageFromJSONBytes(finalResponse); parsed {
			*usage = parsedUsage
	placeholder
		// When the terminal event has an empty output array, reconstruct
		// output from accumulated delta events so the client gets full content.
		// gjson Array() returns empty slice for null, missing, or empty arrays.
		if len(gjson.GetBytes(finalResponse, "output").Array()) == 0 {
			if outputJSON, reconstructed := reconstructResponseOutputFromSSE(bodyText); reconstructed {
				if patched, err := sjson.SetRawBytes(finalResponse, "output", outputJSON); err == nil {
					finalResponse = patched
			placeholder
		placeholder
	placeholder
		finalResponse = supplementCompactionItemFromSSE(c, finalResponse, bodyText)
		body = finalResponse
		if originalModel != mappedModel {
			body = s.replaceModelInResponseBody(body, mappedModel, originalModel)
	placeholder
		// Correct tool calls in final response
		body = s.correctToolCallsInResponseBody(body)
		restoredBody, restoreErr := restoreGrokResponsesClientToolPayload(c, body)
		if restoreErr != nil {
			return nil, fmt.Errorf("restore Grok Responses client tool response: %w", restoreErr)
	placeholder
		restoredBody, restoreErr = restoreOpenAIResponsesNamespacePayload(c, restoredBody)
		if restoreErr != nil {
			return nil, fmt.Errorf("restore OpenAI namespace response: %w", restoreErr)
	placeholder
		body = restoredBody
placeholder else {
		terminalType, terminalPayload, terminalOK := extractOpenAISSETerminalEvent(bodyText)
		if terminalOK && terminalType == "response.failed" {
			msg := extractOpenAISSEErrorMessage(terminalPayload)
			if msg == "" {
				msg = "Upstream compact response failed"
		placeholder
			return nil, s.writeOpenAINonStreamingProtocolError(resp, c, msg)
	placeholder
		usage = s.parseSSEUsageFromBody(bodyText)
		if originalModel != mappedModel {
			bodyText = s.replaceModelInSSEBody(bodyText, mappedModel, originalModel)
	placeholder
		body = []byte(bodyText)
placeholder

	responseheaders.WriteFilteredHeaders(c.Writer.Header(), resp.Header, s.responseHeaderFilter)

	contentType := "application/json; charset=utf-8"
	if !ok {
		contentType = resp.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "text/event-stream"
	placeholder
placeholder
	if !writeOpenAICompactSSEBridge(c, resp.StatusCode, body) {
		c.Data(resp.StatusCode, contentType, body)
placeholder

	return &openaiNonStreamingResult{
		OpenAIUsage:      usage,
		usage:            usage,
		responseID:       extractOpenAIResponseIDFromJSONBytes(body),
		imageCount:       countOpenAIImageOutputsFromSSEBody(bodyText),
		imageOutputSizes: collectOpenAIImageOutputSizesFromSSEBody(bodyText),
placeholder, nil
placeholder

func extractOpenAISSETerminalEvent(body string) (string, []byte, bool) {
	var terminalType string
	var terminalPayload []byte
	forEachOpenAISSEDataPayload(body, func(data []byte) {
		if terminalPayload != nil {
			return
	placeholder
		eventType := strings.TrimSpace(gjson.GetBytes(data, "type").String())
		switch eventType {
		case "response.completed", "response.done", "response.failed", "response.incomplete", "response.cancelled", "response.canceled":
			terminalType = eventType
			terminalPayload = append([]byte(nil), data...)
	placeholder
placeholder)
	if terminalPayload != nil {
		return terminalType, terminalPayload, true
placeholder
	return "", nil, false
placeholder

func extractOpenAISSEErrorMessage(payload []byte) string {
	if len(payload) == 0 {
		return ""
placeholder
	for _, path := range []string{"response.error.message", "error.message", "message"placeholder {
		if msg := strings.TrimSpace(gjson.GetBytes(payload, path).String()); msg != "" {
			return sanitizeUpstreamErrorMessage(msg)
	placeholder
placeholder
	return sanitizeUpstreamErrorMessage(strings.TrimSpace(extractUpstreamErrorMessage(payload)))
placeholder

func sanitizeOpenAIResponseFailedEventForClient(payload []byte, eventType string, clientOutputStarted bool) ([]byte, bool) {
	if eventType != "response.failed" || len(payload) == 0 || !gjson.ValidBytes(payload) {
		return payload, false
placeholder
	updated := payload
	if clientOutputStarted && isOpenAIContextWindowError(extractOpenAISSEErrorMessage(payload), payload) {
		errorPath := ""
		switch {
		case gjson.GetBytes(updated, "response.error").Exists():
			errorPath = "response.error"
		case gjson.GetBytes(updated, "error").Exists():
			errorPath = "error"
	placeholder
		if errorPath != "" {
			next, err := sjson.SetBytes(updated, errorPath+".type", "invalid_request_error")
			if err != nil {
				return payload, false
		placeholder
			updated = next
			next, err = sjson.SetBytes(updated, errorPath+".code", "context_length_exceeded")
			if err != nil {
				return payload, false
		placeholder
			updated = next
	placeholder
placeholder
	if !gjson.GetBytes(updated, "response").Exists() {
		return updated, !bytes.Equal(updated, payload)
placeholder
	for _, path := range []string{
		"response.instructions",
		"response.output",
		"response.usage",
		"response.metadata",
		"response.reasoning",
		"response.tools",
		"response.tool_choice",
		"response.parallel_tool_calls",
		"response.text",
		"response.truncation",
		"response.max_output_tokens",
		"response.incomplete_details",
placeholder {
		next, err := sjson.DeleteBytes(updated, path)
		if err != nil {
			return payload, false
	placeholder
		updated = next
placeholder
	return updated, !bytes.Equal(updated, payload)
placeholder

func (s *OpenAIGatewayService) writeOpenAINonStreamingProtocolError(resp *http.Response, c *gin.Context, message string) error {
	message = sanitizeUpstreamErrorMessage(strings.TrimSpace(message))
	if message == "" {
		message = "Upstream returned an invalid non-streaming response"
placeholder
	setOpsUpstreamError(c, http.StatusBadGateway, message, "")
	// body-signal compact 心跳可能已把响应头提交为 200，此时只能以
	// response.failed 终止事件回传错误，不能再写 JSON+状态码。
	if openAICompactClientWantsStream(c) && StopOpenAICompactSSEKeepaliveCommitted(c) {
		writeOpenAICompactSSEFailureMessage(c, http.StatusBadGateway, "upstream_error", message)
		return fmt.Errorf("non-streaming openai protocol error: %s", message)
placeholder
	responseheaders.WriteFilteredHeaders(c.Writer.Header(), resp.Header, s.responseHeaderFilter)
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusBadGateway, gin.H{
		"error": gin.H{
			"type":    "upstream_error",
			"message": message,
	placeholder,
placeholder)
	return fmt.Errorf("non-streaming openai protocol error: %s", message)
placeholder

func extractCodexFinalResponse(body string) ([]byte, bool) {
	var finalResponse []byte
	forEachOpenAISSEDataPayload(body, func(data []byte) {
		if finalResponse != nil {
			return
	placeholder
		if normalized, changed := normalizeCompletedImageGenerationStatus(data); changed {
			data = normalized
	placeholder
		eventType := gjson.GetBytes(data, "type").String()
		if eventType == "response.done" || eventType == "response.completed" {
			if response := gjson.GetBytes(data, "response"); response.Exists() && response.Type == gjson.JSON && response.Raw != "" {
				finalResponse = []byte(response.Raw)
		placeholder
	placeholder
placeholder)
	if finalResponse != nil {
		return finalResponse, true
placeholder
	return nil, false
placeholder

func normalizeCompletedImageGenerationStatus(data []byte) ([]byte, bool) {
	if len(data) == 0 || !gjson.ValidBytes(data) {
		return data, false
placeholder

	shouldNormalize := func(item gjson.Result) bool {
		if !item.Exists() || !item.IsObject() ||
			strings.TrimSpace(item.Get("type").String()) != "image_generation_call" {
			return false
	placeholder
		switch strings.TrimSpace(item.Get("status").String()) {
		case "generating", "in_progress":
			return strings.TrimSpace(item.Get("result").String()) != ""
		default:
			return false
	placeholder
placeholder

	eventType := strings.TrimSpace(gjson.GetBytes(data, "type").String())
	switch eventType {
	case "response.output_item.done":
		if !shouldNormalize(gjson.GetBytes(data, "item")) {
			return data, false
	placeholder
		updated, err := sjson.SetBytes(data, "item.status", "completed")
		if err != nil {
			return data, false
	placeholder
		return updated, true
	case "response.completed", "response.done":
		output := gjson.GetBytes(data, "response.output")
		if !output.Exists() || !output.IsArray() {
			return data, false
	placeholder
		updated := data
		changed := false
		for i, item := range output.Array() {
			if !shouldNormalize(item) {
				continue
		placeholder
			next, err := sjson.SetBytes(updated, "response.output."+strconv.Itoa(i)+".status", "completed")
			if err != nil {
				return data, false
		placeholder
			updated = next
			changed = true
	placeholder
		return updated, changed
	default:
		return data, false
placeholder
placeholder

func normalizeResponsesStreamingTerminalOutput(data []byte, acc *apicompat.BufferedResponseAccumulator, imageOutputs []json.RawMessage) ([]byte, bool) {
	eventType := strings.TrimSpace(gjson.GetBytes(data, "type").String())
	switch eventType {
	case "response.completed", "response.done", "response.incomplete", "response.cancelled", "response.canceled":
	default:
		return data, false
placeholder

	output := gjson.GetBytes(data, "response.output")
	hasAccumulatedOutput := (acc != nil && acc.HasContent()) || len(imageOutputs) > 0
	if output.Exists() && output.IsArray() {
		if len(output.Array()) > 0 || !hasAccumulatedOutput {
			return data, false
	placeholder
placeholder

	outputJSON := []byte("[]")
	if reconstructed, ok := buildResponsesOutputJSON(acc, imageOutputs); ok {
		outputJSON = reconstructed
placeholder
	updated, err := sjson.SetRawBytes(data, "response.output", outputJSON)
	if err != nil {
		return data, false
placeholder
	return updated, true
placeholder

func responsesStreamEventMayContributeToOutput(eventType string) bool {
	switch eventType {
	case "response.output_text.delta",
		"response.output_item.added",
		"response.function_call_arguments.delta",
		"response.reasoning_summary_text.delta":
		return true
	default:
		return false
placeholder
placeholder

// collectRawResponsesOutputItemsFromSSE 按到达顺序收集 SSE 流中
// response.output_item.done 携带的原始 item。除已产生结果但仍停留在进行中
// 的图片状态外，item 以 raw JSON 逐字节保留，
// 避免经窄结构体重建时丢弃 encrypted_content/summary/opaque 等 compact
// 专属或未来新增字段（#3777 问题 2）。若整条流没有任何 done 事件，退回
// 收集 output_item.added 中的 compaction 类 item——compaction 结果没有
// delta 事件，部分上游只在 added 事件中携带完整 item。
func collectRawResponsesOutputItemsFromSSE(bodyText string) ([]byte, bool) {
	var items []json.RawMessage
	seen := make(map[string]struct{placeholder)
	hasCompactionItem := false
	appendItem := func(item gjson.Result) {
		if !item.Exists() || !item.IsObject() {
			return
	placeholder
		key := strings.TrimSpace(item.Get("id").String())
		if key == "" {
			key = item.Raw
	placeholder
		if _, dup := seen[key]; dup {
			return
	placeholder
		seen[key] = struct{placeholder{placeholder
		if isResponsesCompactionItemType(item.Get("type").String()) {
			hasCompactionItem = true
	placeholder
		items = append(items, json.RawMessage(item.Raw))
placeholder
	forEachOpenAISSEDataPayload(bodyText, func(data []byte) {
		if normalized, changed := normalizeCompletedImageGenerationStatus(data); changed {
			data = normalized
	placeholder
		if strings.TrimSpace(gjson.GetBytes(data, "type").String()) != "response.output_item.done" {
			return
	placeholder
		appendItem(gjson.GetBytes(data, "item"))
placeholder)
	// done 事件未携带 compaction item 时再看 added：覆盖"其他 item 有 done、
	// compaction 只在 added 中"的混合形态；done 已含 compaction 时跳过，
	// 避免同一 item 在无 id 可去重时被收集两份（Codex 要求恰好一个）。
	if !hasCompactionItem {
		forEachOpenAISSEDataPayload(bodyText, func(data []byte) {
			if strings.TrimSpace(gjson.GetBytes(data, "type").String()) != "response.output_item.added" {
				return
		placeholder
			item := gjson.GetBytes(data, "item")
			if !isResponsesCompactionItemType(item.Get("type").String()) {
				return
		placeholder
			appendItem(item)
	placeholder)
placeholder
	if len(items) == 0 {
		return nil, false
placeholder
	outputJSON, err := json.Marshal(items)
	if err != nil {
		return nil, false
placeholder
	return outputJSON, true
placeholder

// isResponsesCompactionItemType reports whether the item type is the Codex
// remote-compact result item ("compaction", upstream alias "compaction_summary").
func isResponsesCompactionItemType(itemType string) bool {
	switch strings.TrimSpace(itemType) {
	case "compaction", "compaction_summary":
		return true
	default:
		return false
placeholder
placeholder

// supplementCompactionItemFromSSE 保证 compact 请求的终态 output 携带
// compaction item：终态 output 非空但缺失 compaction、而原始事件流的
// output_item.done（或 added）中存在时（上游不一致形态），以 raw JSON 补入。
// Codex remote compact v2 只从 output_item.done 收集 item 且要求恰好一个
// compaction item——纯流式透传（v0.1.146）下客户端直接读事件流天然拿得到，
// SSE→JSON 提取链路必须给出等价结果。非 compact 请求原样返回。
func supplementCompactionItemFromSSE(c *gin.Context, finalResponse []byte, bodyText string) []byte {
	if !isOpenAIResponsesCompactPath(c) {
		return finalResponse
placeholder
	if len(gjson.GetBytes(finalResponse, "output").Array()) == 0 {
		// 空 output 由 reconstructResponseOutputFromSSE 整体修补，不在此处理。
		return finalResponse
placeholder
	if responsesOutputHasCompactionItem(finalResponse) {
		return finalResponse
placeholder
	item, found := findRawCompactionItemFromSSE(bodyText)
	if !found {
		return finalResponse
placeholder
	patched, err := sjson.SetRawBytes(finalResponse, "output.-1", item)
	if err != nil {
		return finalResponse
placeholder
	return patched
placeholder

// responsesOutputHasCompactionItem reports whether the response JSON already
// carries a compaction item in its output array.
func responsesOutputHasCompactionItem(response []byte) bool {
	for _, item := range gjson.GetBytes(response, "output").Array() {
		if isResponsesCompactionItemType(item.Get("type").String()) {
			return true
	placeholder
placeholder
	return false
placeholder

// findRawCompactionItemFromSSE 从原始 SSE 事件流中提取第一个 compaction 类
// item 的 raw JSON：output_item.done 优先，output_item.added 兜底。
func findRawCompactionItemFromSSE(bodyText string) (json.RawMessage, bool) {
	var found json.RawMessage
	pick := func(eventType string) {
		forEachOpenAISSEDataPayload(bodyText, func(data []byte) {
			if found != nil {
				return
		placeholder
			if strings.TrimSpace(gjson.GetBytes(data, "type").String()) != eventType {
				return
		placeholder
			item := gjson.GetBytes(data, "item")
			if !item.IsObject() || !isResponsesCompactionItemType(item.Get("type").String()) {
				return
		placeholder
			found = json.RawMessage(item.Raw)
	placeholder)
placeholder
	pick("response.output_item.done")
	if found == nil {
		pick("response.output_item.added")
placeholder
	return found, found != nil
placeholder

// reconstructResponseOutputFromSSE scans raw SSE body text and returns a
// JSON-encoded output array for a terminal event whose output is empty.
// Raw output_item.done items are preferred: per the Responses protocol they
// are the authoritative final form of each item. Delta accumulation only
// covers text/function_call/reasoning content and silently drops unknown
// item types such as compaction — Codex remote compact v2 then fails with
// "expected exactly one compaction output item, got 0" (#3887).
// Returns (nil, false) if nothing could be reconstructed.
func reconstructResponseOutputFromSSE(bodyText string) ([]byte, bool) {
	if outputJSON, ok := collectRawResponsesOutputItemsFromSSE(bodyText); ok {
		return outputJSON, true
placeholder
	acc := apicompat.NewBufferedResponseAccumulator()
	imageOutputs := make([]json.RawMessage, 0, 1)
	seenImages := make(map[string]struct{placeholder)
	forEachOpenAISSEDataPayload(bodyText, func(data []byte) {
		if imageOutput, ok := extractImageGenerationOutputFromSSEData(data, seenImages); ok {
			imageOutputs = append(imageOutputs, imageOutput)
	placeholder
		eventType := strings.TrimSpace(gjson.GetBytes(data, "type").String())
		if responsesStreamEventMayContributeToOutput(eventType) {
			var event apicompat.ResponsesStreamEvent
			if err := json.Unmarshal(data, &event); err == nil {
				acc.ProcessEvent(&event)
		placeholder
	placeholder
placeholder)
	return buildResponsesOutputJSON(acc, imageOutputs)
placeholder

func buildResponsesOutputJSON(acc *apicompat.BufferedResponseAccumulator, imageOutputs []json.RawMessage) ([]byte, bool) {
	if (acc == nil || !acc.HasContent()) && len(imageOutputs) == 0 {
		return nil, false
placeholder
	var output []json.RawMessage
	if acc != nil && acc.HasContent() {
		outputJSON, err := json.Marshal(acc.BuildOutput())
		if err == nil {
			_ = json.Unmarshal(outputJSON, &output)
	placeholder
placeholder
	output = append(output, imageOutputs...)
	if len(output) == 0 {
		return nil, false
placeholder

	outputJSON, err := json.Marshal(output)
	if err != nil {
		return nil, false
placeholder
	return outputJSON, true
placeholder

func extractImageGenerationOutputFromSSEData(data []byte, seen map[string]struct{placeholder) (json.RawMessage, bool) {
	if len(data) == 0 || !gjson.ValidBytes(data) {
		return nil, false
placeholder
	if gjson.GetBytes(data, "type").String() != "response.output_item.done" {
		return nil, false
placeholder
	item := gjson.GetBytes(data, "item")
	if !item.Exists() || !item.IsObject() || item.Get("type").String() != "image_generation_call" {
		return nil, false
placeholder
	if strings.TrimSpace(item.Get("result").String()) == "" {
		return nil, false
placeholder
	key := strings.TrimSpace(item.Get("id").String())
	if key == "" {
		key = strings.TrimSpace(item.Get("output_format").String()) + "|" + strings.TrimSpace(item.Get("result").String())
placeholder
	if key != "" && seen != nil {
		if _, exists := seen[key]; exists {
			return nil, false
	placeholder
		seen[key] = struct{placeholder{placeholder
placeholder
	return json.RawMessage(item.Raw), true
placeholder

func (s *OpenAIGatewayService) parseSSEUsageFromBody(body string) *OpenAIUsage {
	usage := &OpenAIUsage{placeholder
	forEachOpenAISSEDataPayload(body, func(data []byte) {
		s.parseSSEUsageBytes(data, usage)
placeholder)
	return usage
placeholder

func (s *OpenAIGatewayService) replaceModelInSSEBody(body, fromModel, toModel string) string {
	lines := strings.Split(body, "\n")
	for i, line := range lines {
		if _, ok := extractOpenAISSEDataLine(line); !ok {
			continue
	placeholder
		lines[i] = s.replaceModelInSSELine(line, fromModel, toModel)
placeholder
	return strings.Join(lines, "\n")
placeholder
