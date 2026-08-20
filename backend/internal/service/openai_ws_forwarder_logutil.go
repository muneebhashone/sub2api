package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	coderws "github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

func normalizeOpenAIWSLogValue(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "-"
placeholder
	return openAIWSLogValueReplacer.Replace(trimmed)
placeholder

func truncateOpenAIWSLogValue(value string, maxLen int) string {
	normalized := normalizeOpenAIWSLogValue(value)
	if normalized == "-" || maxLen <= 0 {
		return normalized
placeholder
	if len(normalized) <= maxLen {
		return normalized
placeholder
	return normalized[:maxLen] + "..."
placeholder

func openAIWSHeaderValueForLog(headers http.Header, key string) string {
	if headers == nil {
		return "-"
placeholder
	return truncateOpenAIWSLogValue(headers.Get(key), openAIWSHeaderValueMaxLen)
placeholder

func hasOpenAIWSHeader(headers http.Header, key string) bool {
	if headers == nil {
		return false
placeholder
	return strings.TrimSpace(headers.Get(key)) != ""
placeholder

type openAIWSSessionHeaderResolution struct {
	SessionID          string
	ConversationID     string
	SessionSource      string
	ConversationSource string
placeholder

func resolveOpenAIWSSessionHeaders(c *gin.Context, promptCacheKey string) openAIWSSessionHeaderResolution {
	resolution := openAIWSSessionHeaderResolution{
		SessionSource:      "none",
		ConversationSource: "none",
placeholder
	if c != nil && c.Request != nil {
		if sessionID := strings.TrimSpace(c.Request.Header.Get("session_id")); sessionID != "" {
			resolution.SessionID = sessionID
			resolution.SessionSource = "header_session_id"
	placeholder
		if conversationID := strings.TrimSpace(c.Request.Header.Get("conversation_id")); conversationID != "" {
			resolution.ConversationID = conversationID
			resolution.ConversationSource = "header_conversation_id"
			if resolution.SessionID == "" {
				resolution.SessionID = conversationID
				resolution.SessionSource = "header_conversation_id"
		placeholder
	placeholder
placeholder

	cacheKey := strings.TrimSpace(promptCacheKey)
	if cacheKey != "" {
		if resolution.SessionID == "" {
			resolution.SessionID = cacheKey
			resolution.SessionSource = "prompt_cache_key"
	placeholder
placeholder
	return resolution
placeholder

func shouldLogOpenAIWSEvent(idx int, eventType string) bool {
	if idx <= openAIWSEventLogHeadLimit {
		return true
placeholder
	if openAIWSEventLogEveryN > 0 && idx%openAIWSEventLogEveryN == 0 {
		return true
placeholder
	if eventType == "error" || isOpenAIWSTerminalEvent(eventType) {
		return true
placeholder
	return false
placeholder

func shouldLogOpenAIWSBufferedEvent(idx int) bool {
	if idx <= openAIWSBufferLogHeadLimit {
		return true
placeholder
	if openAIWSBufferLogEveryN > 0 && idx%openAIWSBufferLogEveryN == 0 {
		return true
placeholder
	return false
placeholder

func openAIWSEventMayContainModel(eventType string) bool {
	switch eventType {
	case "response.created",
		"response.in_progress",
		"response.completed",
		"response.done",
		"response.failed",
		"response.incomplete",
		"response.cancelled",
		"response.canceled":
		return true
	default:
		trimmed := strings.TrimSpace(eventType)
		if trimmed == eventType {
			return false
	placeholder
		switch trimmed {
		case "response.created",
			"response.in_progress",
			"response.completed",
			"response.done",
			"response.failed",
			"response.incomplete",
			"response.cancelled",
			"response.canceled":
			return true
		default:
			return false
	placeholder
placeholder
placeholder

func openAIWSEventMayContainToolCalls(eventType string) bool {
	eventType = strings.TrimSpace(eventType)
	if eventType == "" {
		return false
placeholder
	if strings.Contains(eventType, "function_call") || strings.Contains(eventType, "tool_call") {
		return true
placeholder
	switch eventType {
	case "response.output_item.added", "response.output_item.done", "response.completed", "response.done":
		return true
	default:
		return false
placeholder
placeholder

func openAIWSEventShouldParseUsage(eventType string) bool {
	eventType = strings.TrimSpace(eventType)
	if eventType == "error" || openAIStreamEventTypeIsTerminal(eventType) {
		return true
placeholder
	return strings.HasPrefix(eventType, "response.") && !strings.HasSuffix(eventType, ".delta")
placeholder

func openAIWSMessageShouldParseUsage(eventType string, message []byte) bool {
	return openAIWSEventShouldParseUsage(eventType) && bytes.Contains(message, []byte(`"usage"`))
placeholder

func parseOpenAIWSEventEnvelope(message []byte) (eventType string, responseID string, response gjson.Result) {
	if len(message) == 0 {
		return "", "", gjson.Result{placeholder
placeholder
	values := gjson.GetManyBytes(message, "type", "response.id", "id", "response")
	eventType = strings.TrimSpace(values[0].String())
	if id := strings.TrimSpace(values[1].String()); id != "" {
		responseID = id
placeholder else {
		responseID = strings.TrimSpace(values[2].String())
placeholder
	return eventType, responseID, values[3]
placeholder

func openAIWSMessageLikelyContainsToolCalls(message []byte) bool {
	if len(message) == 0 {
		return false
placeholder
	return bytes.Contains(message, []byte(`"tool_calls"`)) ||
		bytes.Contains(message, []byte(`"tool_call"`)) ||
		bytes.Contains(message, []byte(`"function_call"`))
placeholder

func parseOpenAIWSResponseUsageFromCompletedEvent(message []byte, usage *OpenAIUsage) {
	if usage == nil || len(message) == 0 || !bytes.Contains(message, []byte(`"usage"`)) {
		return
placeholder
	if parsedUsage, ok := extractOpenAIUsageFromJSONBytes(message); ok {
		if openAIStreamEventTypeIsTerminal(effectiveOpenAISSEEventType(message, "")) {
			if !openAIUsageHasTokens(&parsedUsage) && openAIUsageHasTokens(usage) {
				return
		placeholder
			*usage = parsedUsage
	placeholder else {
			mergeOpenAIUsageNonZero(usage, parsedUsage)
	placeholder
placeholder
placeholder

func parseOpenAIWSErrorEventFields(message []byte) (code string, errType string, errMessage string) {
	if len(message) == 0 {
		return "", "", ""
placeholder
	values := gjson.GetManyBytes(message, "error.code", "error.type", "error.message")
	return strings.TrimSpace(values[0].String()), strings.TrimSpace(values[1].String()), strings.TrimSpace(values[2].String())
placeholder

func summarizeOpenAIWSErrorEventFieldsFromRaw(codeRaw, errTypeRaw, errMessageRaw string) (code string, errType string, errMessage string) {
	code = truncateOpenAIWSLogValue(codeRaw, openAIWSLogValueMaxLen)
	errType = truncateOpenAIWSLogValue(errTypeRaw, openAIWSLogValueMaxLen)
	errMessage = truncateOpenAIWSLogValue(errMessageRaw, openAIWSLogValueMaxLen)
	return code, errType, errMessage
placeholder

func summarizeOpenAIWSErrorEventFields(message []byte) (code string, errType string, errMessage string) {
	if len(message) == 0 {
		return "-", "-", "-"
placeholder
	return summarizeOpenAIWSErrorEventFieldsFromRaw(parseOpenAIWSErrorEventFields(message))
placeholder

func summarizeOpenAIWSPayloadKeySizes(payload map[string]any, topN int) string {
	if len(payload) == 0 {
		return "-"
placeholder
	type keySize struct {
		Key  string
		Size int
placeholder
	sizes := make([]keySize, 0, len(payload))
	for key, value := range payload {
		size := estimateOpenAIWSPayloadValueSize(value, openAIWSPayloadSizeEstimateDepth)
		sizes = append(sizes, keySize{Key: key, Size: sizeplaceholder)
placeholder
	sort.Slice(sizes, func(i, j int) bool {
		if sizes[i].Size == sizes[j].Size {
			return sizes[i].Key < sizes[j].Key
	placeholder
		return sizes[i].Size > sizes[j].Size
placeholder)

	if topN <= 0 || topN > len(sizes) {
		topN = len(sizes)
placeholder
	parts := make([]string, 0, topN)
	for idx := 0; idx < topN; idx++ {
		item := sizes[idx]
		parts = append(parts, fmt.Sprintf("%s:%d", item.Key, item.Size))
placeholder
	return strings.Join(parts, ",")
placeholder

func estimateOpenAIWSPayloadValueSize(value any, depth int) int {
	if depth <= 0 {
		return -1
placeholder
	switch v := value.(type) {
	case nil:
		return 0
	case string:
		return len(v)
	case []byte:
		return len(v)
	case bool:
		return 1
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return 8
	case float32, float64:
		return 8
	case map[string]any:
		if len(v) == 0 {
			return 2
	placeholder
		total := 2
		count := 0
		for key, item := range v {
			count++
			if count > openAIWSPayloadSizeEstimateMaxItems {
				return -1
		placeholder
			itemSize := estimateOpenAIWSPayloadValueSize(item, depth-1)
			if itemSize < 0 {
				return -1
		placeholder
			total += len(key) + itemSize + 3
			if total > openAIWSPayloadSizeEstimateMaxBytes {
				return -1
		placeholder
	placeholder
		return total
	case []any:
		if len(v) == 0 {
			return 2
	placeholder
		total := 2
		limit := len(v)
		if limit > openAIWSPayloadSizeEstimateMaxItems {
			return -1
	placeholder
		for i := 0; i < limit; i++ {
			itemSize := estimateOpenAIWSPayloadValueSize(v[i], depth-1)
			if itemSize < 0 {
				return -1
		placeholder
			total += itemSize + 1
			if total > openAIWSPayloadSizeEstimateMaxBytes {
				return -1
		placeholder
	placeholder
		return total
	default:
		raw, err := json.Marshal(v)
		if err != nil {
			return -1
	placeholder
		if len(raw) > openAIWSPayloadSizeEstimateMaxBytes {
			return -1
	placeholder
		return len(raw)
placeholder
placeholder

func openAIWSPayloadString(payload map[string]any, key string) string {
	if len(payload) == 0 {
		return ""
placeholder
	raw, ok := payload[key]
	if !ok {
		return ""
placeholder
	switch v := raw.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(v)
	case []byte:
		return strings.TrimSpace(string(v))
	default:
		return ""
placeholder
placeholder

func openAIWSPayloadStringFromRaw(payload []byte, key string) string {
	if len(payload) == 0 || strings.TrimSpace(key) == "" {
		return ""
placeholder
	return strings.TrimSpace(gjson.GetBytes(payload, key).String())
placeholder

func openAIWSPayloadBoolFromRaw(payload []byte, key string, defaultValue bool) bool {
	if len(payload) == 0 || strings.TrimSpace(key) == "" {
		return defaultValue
placeholder
	value := gjson.GetBytes(payload, key)
	if !value.Exists() {
		return defaultValue
placeholder
	if value.Type != gjson.True && value.Type != gjson.False {
		return defaultValue
placeholder
	return value.Bool()
placeholder

func openAIWSSessionHashesFromID(sessionID string) (string, string) {
	return deriveOpenAISessionHashes(sessionID)
placeholder

func extractOpenAIWSImageURL(value any) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case map[string]any:
		if raw, ok := v["url"].(string); ok {
			return strings.TrimSpace(raw)
	placeholder
placeholder
	return ""
placeholder

func summarizeOpenAIWSInput(input any) string {
	items, ok := input.([]any)
	if !ok || len(items) == 0 {
		return "-"
placeholder

	itemCount := len(items)
	textChars := 0
	imageDataURLs := 0
	imageDataURLChars := 0
	imageRemoteURLs := 0

	handleContentItem := func(contentItem map[string]any) {
		contentType, _ := contentItem["type"].(string)
		switch strings.TrimSpace(contentType) {
		case "input_text", "output_text", "text":
			if text, ok := contentItem["text"].(string); ok {
				textChars += len(text)
		placeholder
		case "input_image":
			imageURL := extractOpenAIWSImageURL(contentItem["image_url"])
			if imageURL == "" {
				return
		placeholder
			if strings.HasPrefix(strings.ToLower(imageURL), "data:image/") {
				imageDataURLs++
				imageDataURLChars += len(imageURL)
				return
		placeholder
			imageRemoteURLs++
	placeholder
placeholder

	handleInputItem := func(inputItem map[string]any) {
		if content, ok := inputItem["content"].([]any); ok {
			for _, rawContent := range content {
				contentItem, ok := rawContent.(map[string]any)
				if !ok {
					continue
			placeholder
				handleContentItem(contentItem)
		placeholder
			return
	placeholder

		itemType, _ := inputItem["type"].(string)
		switch strings.TrimSpace(itemType) {
		case "input_text", "output_text", "text":
			if text, ok := inputItem["text"].(string); ok {
				textChars += len(text)
		placeholder
		case "input_image":
			imageURL := extractOpenAIWSImageURL(inputItem["image_url"])
			if imageURL == "" {
				return
		placeholder
			if strings.HasPrefix(strings.ToLower(imageURL), "data:image/") {
				imageDataURLs++
				imageDataURLChars += len(imageURL)
				return
		placeholder
			imageRemoteURLs++
	placeholder
placeholder

	for _, rawItem := range items {
		inputItem, ok := rawItem.(map[string]any)
		if !ok {
			continue
	placeholder
		handleInputItem(inputItem)
placeholder

	return fmt.Sprintf(
		"items=%d,text_chars=%d,image_data_urls=%d,image_data_url_chars=%d,image_remote_urls=%d",
		itemCount,
		textChars,
		imageDataURLs,
		imageDataURLChars,
		imageRemoteURLs,
	)
placeholder

func dropOpenAIWSPayloadKey(payload map[string]any, key string, removed *[]string) {
	if len(payload) == 0 || strings.TrimSpace(key) == "" {
		return
placeholder
	if _, exists := payload[key]; !exists {
		return
placeholder
	delete(payload, key)
	*removed = append(*removed, key)
placeholder

// applyOpenAIWSRetryPayloadStrategy 在 WS 连续失败时仅移除无语义字段，
// 避免重试成功却改变原始请求语义。
// 注意：prompt_cache_key 不应在重试中移除；它常用于会话稳定标识（session_id 兜底）。
func applyOpenAIWSRetryPayloadStrategy(payload map[string]any, attempt int) (strategy string, removedKeys []string) {
	if len(payload) == 0 {
		return "empty", nil
placeholder
	if attempt <= 1 {
		return "full", nil
placeholder

	removed := make([]string, 0, 2)
	if attempt >= 2 {
		dropOpenAIWSPayloadKey(payload, "include", &removed)
placeholder

	if len(removed) == 0 {
		return "full", nil
placeholder
	sort.Strings(removed)
	return "trim_optional_fields", removed
placeholder

func logOpenAIWSModeInfo(format string, args ...any) {
	logger.LegacyPrintf("service.openai_gateway", "[OpenAI WS Mode][openai_ws_mode=true] "+format, args...)
placeholder

func isOpenAIWSModeDebugEnabled() bool {
	return logger.L().Core().Enabled(zap.DebugLevel)
placeholder

func logOpenAIWSModeDebug(format string, args ...any) {
	if !isOpenAIWSModeDebugEnabled() {
		return
placeholder
	logger.LegacyPrintf("service.openai_gateway", "[debug] [OpenAI WS Mode][openai_ws_mode=true] "+format, args...)
placeholder

func logOpenAIWSBindResponseAccountWarn(groupID, accountID int64, responseID string, err error) {
	if err == nil {
		return
placeholder
	logger.L().Warn(
		"openai.ws_bind_response_account_failed",
		zap.Int64("group_id", groupID),
		zap.Int64("account_id", accountID),
		zap.String("response_id", truncateOpenAIWSLogValue(responseID, openAIWSIDValueMaxLen)),
		zap.Error(err),
	)
placeholder

func summarizeOpenAIWSReadCloseError(err error) (status string, reason string) {
	if err == nil {
		return "-", "-"
placeholder
	statusCode := coderws.CloseStatus(err)
	if statusCode == -1 {
		return "-", "-"
placeholder
	closeStatus := fmt.Sprintf("%d(%s)", int(statusCode), statusCode.String())
	closeReason := "-"
	var closeErr coderws.CloseError
	if errors.As(err, &closeErr) {
		reasonText := strings.TrimSpace(closeErr.Reason)
		if reasonText != "" {
			closeReason = normalizeOpenAIWSLogValue(reasonText)
	placeholder
placeholder
	return normalizeOpenAIWSLogValue(closeStatus), closeReason
placeholder

func unwrapOpenAIWSDialBaseError(err error) error {
	if err == nil {
		return nil
placeholder
	var dialErr *openAIWSDialError
	if errors.As(err, &dialErr) && dialErr != nil && dialErr.Err != nil {
		return dialErr.Err
placeholder
	return err
placeholder

func openAIWSDialRespHeaderForLog(err error, key string) string {
	var dialErr *openAIWSDialError
	if !errors.As(err, &dialErr) || dialErr == nil || dialErr.ResponseHeaders == nil {
		return "-"
placeholder
	return truncateOpenAIWSLogValue(dialErr.ResponseHeaders.Get(key), openAIWSHeaderValueMaxLen)
placeholder

func classifyOpenAIWSDialError(err error) string {
	if err == nil {
		return "-"
placeholder
	baseErr := unwrapOpenAIWSDialBaseError(err)
	if baseErr == nil {
		return "-"
placeholder
	if errors.Is(baseErr, context.DeadlineExceeded) {
		return "ctx_deadline_exceeded"
placeholder
	if errors.Is(baseErr, context.Canceled) {
		return "ctx_canceled"
placeholder
	var netErr net.Error
	if errors.As(baseErr, &netErr) && netErr.Timeout() {
		return "net_timeout"
placeholder
	if status := coderws.CloseStatus(baseErr); status != -1 {
		return normalizeOpenAIWSLogValue(fmt.Sprintf("ws_close_%d", int(status)))
placeholder
	message := strings.ToLower(strings.TrimSpace(baseErr.Error()))
	switch {
	case strings.Contains(message, "handshake not finished"):
		return "handshake_not_finished"
	case strings.Contains(message, "bad handshake"):
		return "bad_handshake"
	case strings.Contains(message, "connection refused"):
		return "connection_refused"
	case strings.Contains(message, "no such host"):
		return "dns_not_found"
	case strings.Contains(message, "tls"):
		return "tls_error"
	case strings.Contains(message, "i/o timeout"):
		return "io_timeout"
	case strings.Contains(message, "context deadline exceeded"):
		return "ctx_deadline_exceeded"
	default:
		return "dial_error"
placeholder
placeholder

func summarizeOpenAIWSDialError(err error) (
	statusCode int,
	dialClass string,
	closeStatus string,
	closeReason string,
	respServer string,
	respVia string,
	respCFRay string,
	respRequestID string,
) {
	dialClass = "-"
	closeStatus = "-"
	closeReason = "-"
	respServer = "-"
	respVia = "-"
	respCFRay = "-"
	respRequestID = "-"
	if err == nil {
		return
placeholder
	var dialErr *openAIWSDialError
	if errors.As(err, &dialErr) && dialErr != nil {
		statusCode = dialErr.StatusCode
		respServer = openAIWSDialRespHeaderForLog(err, "server")
		respVia = openAIWSDialRespHeaderForLog(err, "via")
		respCFRay = openAIWSDialRespHeaderForLog(err, "cf-ray")
		respRequestID = openAIWSDialRespHeaderForLog(err, "x-request-id")
placeholder
	dialClass = normalizeOpenAIWSLogValue(classifyOpenAIWSDialError(err))
	closeStatus, closeReason = summarizeOpenAIWSReadCloseError(unwrapOpenAIWSDialBaseError(err))
	return
placeholder

func isOpenAIWSClientDisconnectError(err error) bool {
	if err == nil {
		return false
placeholder
	if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) || errors.Is(err, context.Canceled) {
		return true
placeholder
	switch coderws.CloseStatus(err) {
	case coderws.StatusNormalClosure, coderws.StatusGoingAway, coderws.StatusNoStatusRcvd, coderws.StatusAbnormalClosure:
		return true
placeholder
	message := strings.ToLower(strings.TrimSpace(err.Error()))
	if message == "" {
		return false
placeholder
	return strings.Contains(message, "failed to read frame header: eof") ||
		strings.Contains(message, "unexpected eof") ||
		strings.Contains(message, "use of closed network connection") ||
		strings.Contains(message, "connection reset by peer") ||
		strings.Contains(message, "broken pipe") ||
		strings.Contains(message, "an existing connection was forcibly closed by the remote host") ||
		strings.Contains(message, "an established connection was aborted")
placeholder

func classifyOpenAIWSReadFallbackReason(err error) string {
	if err == nil {
		return "read_event"
placeholder
	switch coderws.CloseStatus(err) {
	case coderws.StatusPolicyViolation:
		return "policy_violation"
	case coderws.StatusMessageTooBig:
		return "message_too_big"
	default:
		return "read_event"
placeholder
placeholder

func sortedKeys(m map[string]any) []string {
	if len(m) == 0 {
		return nil
placeholder
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
placeholder
	sort.Strings(keys)
	return keys
placeholder
