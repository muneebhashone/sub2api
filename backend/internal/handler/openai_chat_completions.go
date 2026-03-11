package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// ChatCompletions handles OpenAI Chat Completions API compatibility.
// POST /v1/chat/completions
func (h *OpenAIGatewayHandler) ChatCompletions(c *gin.Context) {
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

	// Preserve original chat-completions request for upstream passthrough when needed.
	c.Set(service.OpenAIChatCompletionsBodyKey, body)

	var chatReq map[string]any
	if err := json.Unmarshal(body, &chatReq); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Failed to parse request body")
		return
placeholder

	includeUsage := false
	if streamOptions, ok := chatReq["stream_options"].(map[string]any); ok {
		if v, ok := streamOptions["include_usage"].(bool); ok {
			includeUsage = v
	placeholder
placeholder
	c.Set(service.OpenAIChatCompletionsIncludeUsageKey, includeUsage)

	converted, err := service.ConvertChatCompletionsToResponses(chatReq)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", err.Error())
		return
placeholder

	convertedBody, err := json.Marshal(converted)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "api_error", "Failed to process request")
		return
placeholder

	stream, _ := converted["stream"].(bool)
	model, _ := converted["model"].(string)
	originalWriter := c.Writer
	writer := newChatCompletionsResponseWriter(c.Writer, stream, includeUsage, model)
	c.Writer = writer
	c.Request.Body = io.NopCloser(bytes.NewReader(convertedBody))
	c.Request.ContentLength = int64(len(convertedBody))

	h.Responses(c)
	writer.Finalize()
	c.Writer = originalWriter
placeholder

type chatCompletionsResponseWriter struct {
	gin.ResponseWriter
	stream       bool
	includeUsage bool
	buffer       bytes.Buffer
	streamBuf    bytes.Buffer
	state        *chatCompletionStreamState
	corrector    *service.CodexToolCorrector
	finalized    bool
	passthrough  bool
placeholder

type chatCompletionStreamState struct {
	id            string
	model         string
	created       int64
	sentRole      bool
	sawToolCall   bool
	sawText       bool
	toolCallIndex map[string]int
	usage         map[string]any
placeholder

func newChatCompletionsResponseWriter(w gin.ResponseWriter, stream bool, includeUsage bool, model string) *chatCompletionsResponseWriter {
	return &chatCompletionsResponseWriter{
		ResponseWriter: w,
		stream:         stream,
		includeUsage:   includeUsage,
		state: &chatCompletionStreamState{
			model:         strings.TrimSpace(model),
			toolCallIndex: make(map[string]int),
	placeholder,
		corrector: service.NewCodexToolCorrector(),
placeholder
placeholder

func (w *chatCompletionsResponseWriter) Write(data []byte) (int, error) {
	if w.passthrough {
		return w.ResponseWriter.Write(data)
placeholder
	if w.stream {
		n, err := w.streamBuf.Write(data)
		if err != nil {
			return n, err
	placeholder
		w.flushStreamBuffer()
		return n, nil
placeholder

	if w.finalized {
		return len(data), nil
placeholder
	return w.buffer.Write(data)
placeholder

func (w *chatCompletionsResponseWriter) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
placeholder

func (w *chatCompletionsResponseWriter) Finalize() {
	if w.finalized {
		return
placeholder
	w.finalized = true
	if w.passthrough {
		return
placeholder
	if w.stream {
		return
placeholder

	body := w.buffer.Bytes()
	if len(body) == 0 {
		return
placeholder

	w.ResponseWriter.Header().Del("Content-Length")

	converted, err := service.ConvertResponsesToChatCompletion(body)
	if err != nil {
		_, _ = w.ResponseWriter.Write(body)
		return
placeholder

	corrected := converted
	if correctedStr, ok := w.corrector.CorrectToolCallsInSSEData(string(converted)); ok {
		corrected = []byte(correctedStr)
placeholder

	_, _ = w.ResponseWriter.Write(corrected)
placeholder

func (w *chatCompletionsResponseWriter) SetPassthrough() {
	w.passthrough = true
placeholder

func (w *chatCompletionsResponseWriter) Status() int {
	if w.ResponseWriter == nil {
		return 0
placeholder
	return w.ResponseWriter.Status()
placeholder

func (w *chatCompletionsResponseWriter) Written() bool {
	if w.ResponseWriter == nil {
		return false
placeholder
	return w.ResponseWriter.Written()
placeholder

func (w *chatCompletionsResponseWriter) flushStreamBuffer() {
	for {
		buf := w.streamBuf.Bytes()
		idx := bytes.IndexByte(buf, '\n')
		if idx == -1 {
			return
	placeholder
		lineBytes := w.streamBuf.Next(idx + 1)
		line := strings.TrimRight(string(lineBytes), "\r\n")
		w.handleStreamLine(line)
placeholder
placeholder

func (w *chatCompletionsResponseWriter) handleStreamLine(line string) {
	if line == "" {
		return
placeholder
	if strings.HasPrefix(line, ":") {
		_, _ = w.ResponseWriter.Write([]byte(line + "\n\n"))
		return
placeholder
	if !strings.HasPrefix(line, "data:") {
		return
placeholder

	data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
	for _, chunk := range w.convertResponseDataToChatChunks(data) {
		if chunk == "" {
			continue
	placeholder
		if chunk == "[DONE]" {
			_, _ = w.ResponseWriter.Write([]byte("data: [DONE]\n\n"))
			continue
	placeholder
		_, _ = w.ResponseWriter.Write([]byte("data: " + chunk + "\n\n"))
placeholder
placeholder

func (w *chatCompletionsResponseWriter) convertResponseDataToChatChunks(data string) []string {
	if data == "" {
		return nil
placeholder
	if data == "[DONE]" {
		return []string{"[DONE]"placeholder
placeholder

	var payload map[string]any
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		return []string{dataplaceholder
placeholder

	if _, ok := payload["error"]; ok {
		return []string{dataplaceholder
placeholder

	eventType := strings.TrimSpace(getString(payload["type"]))
	if eventType == "" {
		return []string{dataplaceholder
placeholder

	w.state.applyMetadata(payload)

	switch eventType {
	case "response.created":
		return nil
	case "response.output_text.delta":
		delta := getString(payload["delta"])
		if delta == "" {
			return nil
	placeholder
		w.state.sawText = true
		return []string{w.buildTextDeltaChunk(delta)placeholder
	case "response.output_text.done":
		if w.state.sawText {
			return nil
	placeholder
		text := getString(payload["text"])
		if text == "" {
			return nil
	placeholder
		w.state.sawText = true
		return []string{w.buildTextDeltaChunk(text)placeholder
	case "response.output_item.added", "response.output_item.delta":
		if item, ok := payload["item"].(map[string]any); ok {
			if callID, name, args, ok := extractToolCallFromItem(item); ok {
				w.state.sawToolCall = true
				return []string{w.buildToolCallChunk(callID, name, args)placeholder
		placeholder
	placeholder
	case "response.completed", "response.done":
		if responseObj, ok := payload["response"].(map[string]any); ok {
			w.state.applyResponseUsage(responseObj)
	placeholder
		return []string{w.buildFinalChunk()placeholder
placeholder

	if strings.Contains(eventType, "tool_call") || strings.Contains(eventType, "function_call") {
		callID := strings.TrimSpace(getString(payload["call_id"]))
		if callID == "" {
			callID = strings.TrimSpace(getString(payload["tool_call_id"]))
	placeholder
		if callID == "" {
			callID = strings.TrimSpace(getString(payload["id"]))
	placeholder
		args := getString(payload["delta"])
		name := strings.TrimSpace(getString(payload["name"]))
		if callID != "" && (args != "" || name != "") {
			w.state.sawToolCall = true
			return []string{w.buildToolCallChunk(callID, name, args)placeholder
	placeholder
placeholder

	return nil
placeholder

func (w *chatCompletionsResponseWriter) buildTextDeltaChunk(delta string) string {
	w.state.ensureDefaults()
	payload := map[string]any{
		"content": delta,
placeholder
	if !w.state.sentRole {
		payload["role"] = "assistant"
		w.state.sentRole = true
placeholder
	return w.buildChunk(payload, nil, nil)
placeholder

func (w *chatCompletionsResponseWriter) buildToolCallChunk(callID, name, args string) string {
	w.state.ensureDefaults()
	index := w.state.toolCallIndexFor(callID)
	function := map[string]any{placeholder
	if name != "" {
		function["name"] = name
placeholder
	if args != "" {
		function["arguments"] = args
placeholder
	toolCall := map[string]any{
		"index":    index,
		"id":       callID,
		"type":     "function",
		"function": function,
placeholder

	delta := map[string]any{
		"tool_calls": []any{toolCallplaceholder,
placeholder
	if !w.state.sentRole {
		delta["role"] = "assistant"
		w.state.sentRole = true
placeholder

	return w.buildChunk(delta, nil, nil)
placeholder

func (w *chatCompletionsResponseWriter) buildFinalChunk() string {
	w.state.ensureDefaults()
	finishReason := "stop"
	if w.state.sawToolCall {
		finishReason = "tool_calls"
placeholder
	usage := map[string]any(nil)
	if w.includeUsage && w.state.usage != nil {
		usage = w.state.usage
placeholder
	return w.buildChunk(map[string]any{placeholder, finishReason, usage)
placeholder

func (w *chatCompletionsResponseWriter) buildChunk(delta map[string]any, finishReason any, usage map[string]any) string {
	w.state.ensureDefaults()
	chunk := map[string]any{
		"id":      w.state.id,
		"object":  "chat.completion.chunk",
		"created": w.state.created,
		"model":   w.state.model,
		"choices": []any{
			map[string]any{
				"index":         0,
				"delta":         delta,
				"finish_reason": finishReason,
		placeholder,
	placeholder,
placeholder
	if usage != nil {
		chunk["usage"] = usage
placeholder

	data, _ := json.Marshal(chunk)
	if corrected, ok := w.corrector.CorrectToolCallsInSSEData(string(data)); ok {
		return corrected
placeholder
	return string(data)
placeholder

func (s *chatCompletionStreamState) ensureDefaults() {
	if s.id == "" {
		s.id = "chatcmpl-" + randomHexUnsafe(12)
placeholder
	if s.model == "" {
		s.model = "unknown"
placeholder
	if s.created == 0 {
		s.created = time.Now().Unix()
placeholder
placeholder

func (s *chatCompletionStreamState) toolCallIndexFor(callID string) int {
	if idx, ok := s.toolCallIndex[callID]; ok {
		return idx
placeholder
	idx := len(s.toolCallIndex)
	s.toolCallIndex[callID] = idx
	return idx
placeholder

func (s *chatCompletionStreamState) applyMetadata(payload map[string]any) {
	if responseObj, ok := payload["response"].(map[string]any); ok {
		s.applyResponseMetadata(responseObj)
placeholder

	if s.id == "" {
		if id := strings.TrimSpace(getString(payload["response_id"])); id != "" {
			s.id = id
	placeholder else if id := strings.TrimSpace(getString(payload["id"])); id != "" {
			s.id = id
	placeholder
placeholder
	if s.model == "" {
		if model := strings.TrimSpace(getString(payload["model"])); model != "" {
			s.model = model
	placeholder
placeholder
	if s.created == 0 {
		if created := getInt64(payload["created_at"]); created != 0 {
			s.created = created
	placeholder else if created := getInt64(payload["created"]); created != 0 {
			s.created = created
	placeholder
placeholder
placeholder

func (s *chatCompletionStreamState) applyResponseMetadata(responseObj map[string]any) {
	if s.id == "" {
		if id := strings.TrimSpace(getString(responseObj["id"])); id != "" {
			s.id = id
	placeholder
placeholder
	if s.model == "" {
		if model := strings.TrimSpace(getString(responseObj["model"])); model != "" {
			s.model = model
	placeholder
placeholder
	if s.created == 0 {
		if created := getInt64(responseObj["created_at"]); created != 0 {
			s.created = created
	placeholder
placeholder
placeholder

func (s *chatCompletionStreamState) applyResponseUsage(responseObj map[string]any) {
	usage, ok := responseObj["usage"].(map[string]any)
	if !ok {
		return
placeholder
	promptTokens := int(getNumber(usage["input_tokens"]))
	completionTokens := int(getNumber(usage["output_tokens"]))
	if promptTokens == 0 && completionTokens == 0 {
		return
placeholder
	s.usage = map[string]any{
		"prompt_tokens":     promptTokens,
		"completion_tokens": completionTokens,
		"total_tokens":      promptTokens + completionTokens,
placeholder
placeholder

func extractToolCallFromItem(item map[string]any) (string, string, string, bool) {
	itemType := strings.TrimSpace(getString(item["type"]))
	if itemType != "tool_call" && itemType != "function_call" {
		return "", "", "", false
placeholder
	callID := strings.TrimSpace(getString(item["call_id"]))
	if callID == "" {
		callID = strings.TrimSpace(getString(item["id"]))
placeholder
	name := strings.TrimSpace(getString(item["name"]))
	args := getString(item["arguments"])
	if fn, ok := item["function"].(map[string]any); ok {
		if name == "" {
			name = strings.TrimSpace(getString(fn["name"]))
	placeholder
		if args == "" {
			args = getString(fn["arguments"])
	placeholder
placeholder
	if callID == "" && name == "" && args == "" {
		return "", "", "", false
placeholder
	if callID == "" {
		callID = "call_" + randomHexUnsafe(6)
placeholder
	return callID, name, args, true
placeholder

func getString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case json.Number:
		return v.String()
	default:
		return ""
placeholder
placeholder

func getNumber(value any) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case json.Number:
		f, _ := v.Float64()
		return f
	default:
		return 0
placeholder
placeholder

func getInt64(value any) int64 {
	switch v := value.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	case json.Number:
		i, _ := v.Int64()
		return i
	default:
		return 0
placeholder
placeholder

func randomHexUnsafe(byteLength int) string {
	if byteLength <= 0 {
		byteLength = 8
placeholder
	buf := make([]byte, byteLength)
	if _, err := rand.Read(buf); err != nil {
		return "000000"
placeholder
	return hex.EncodeToString(buf)
placeholder
