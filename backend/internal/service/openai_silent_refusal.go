package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/apicompat"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

const (
	openAISilentRefusalMinRequestBodyBytes = 64 * 1024
	openAISilentRefusalErrorCode           = "openai_silent_refusal"
	openAISilentRefusalUpstreamMessage     = "OpenAI upstream returned an empty completion stream with finish_reason=stop and no usage"
	openAISilentRefusalClientMessage       = "Upstream returned an empty completion without usage; no fallback account was available"
	openAIResponsesEmptyCompletedMessage   = "OpenAI upstream returned an empty response.completed stream with no output and no usage"
)

type openAIChatSilentRefusalDetector struct {
	enabled         bool
	sawContent      bool
	sawToolCall     bool
	sawFunctionCall bool
	sawUsage        bool
	sawError        bool
	sawReasoning    bool
	sawFinish       bool
	finishReason    string
placeholder

func newOpenAIChatSilentRefusalDetector(requestBodyLen int) *openAIChatSilentRefusalDetector {
	return &openAIChatSilentRefusalDetector{
		enabled: requestBodyLen >= openAISilentRefusalMinRequestBodyBytes,
placeholder
placeholder

func (d *openAIChatSilentRefusalDetector) Enabled() bool {
	return d != nil && d.enabled
placeholder

func (d *openAIChatSilentRefusalDetector) ObserveSSELine(line string) {
	if d == nil || !d.enabled {
		return
placeholder
	if eventType, ok := extractOpenAISSEEventLine(line); ok {
		d.observeEventType(eventType)
		return
placeholder
	if payload, ok := extractOpenAISSEDataLine(line); ok {
		d.ObservePayload([]byte(payload))
placeholder
placeholder

func (d *openAIChatSilentRefusalDetector) ObservePayload(payload []byte) {
	if d == nil || !d.enabled {
		return
placeholder
	payload = bytes.TrimSpace(payload)
	if len(payload) == 0 || bytes.Equal(payload, []byte("[DONE]")) {
		return
placeholder
	if !gjson.ValidBytes(payload) {
		return
placeholder

	eventType := strings.TrimSpace(gjson.GetBytes(payload, "type").String())
	d.observeEventType(eventType)

	if gjson.GetBytes(payload, "error").Exists() {
		d.sawError = true
placeholder
	if usage := gjson.GetBytes(payload, "usage"); usage.Exists() && usage.IsObject() {
		d.sawUsage = true
placeholder
	if usage := gjson.GetBytes(payload, "response.usage"); usage.Exists() && usage.IsObject() {
		d.sawUsage = true
placeholder

	d.observeChatChoicesPayload(payload)
	d.observeResponsesPayload(payload, eventType)
placeholder

func (d *openAIChatSilentRefusalDetector) ObserveChatChunk(chunk apicompat.ChatCompletionsChunk) {
	if d == nil || !d.enabled {
		return
placeholder
	if chunk.Usage != nil {
		d.sawUsage = true
placeholder
	for _, choice := range chunk.Choices {
		if choice.FinishReason != nil {
			d.observeFinishReason(*choice.FinishReason)
	placeholder
		delta := choice.Delta
		if delta.Content != nil && *delta.Content != "" {
			d.sawContent = true
	placeholder
		if delta.ReasoningContent != nil {
			d.sawReasoning = true
	placeholder
		if len(delta.ToolCalls) > 0 {
			d.sawToolCall = true
	placeholder
placeholder
placeholder

func (d *openAIChatSilentRefusalDetector) ShouldReleaseClientOutput() bool {
	if d == nil || !d.enabled {
		return true
placeholder
	if d.sawContent || d.sawToolCall || d.sawFunctionCall || d.sawUsage || d.sawError || d.sawReasoning {
		return true
placeholder
	return d.sawFinish && d.finishReason != "" && d.finishReason != "stop"
placeholder

func (d *openAIChatSilentRefusalDetector) IsSilentRefusal() bool {
	if d == nil || !d.enabled {
		return false
placeholder
	return !d.sawContent &&
		!d.sawToolCall &&
		!d.sawFunctionCall &&
		!d.sawUsage &&
		!d.sawError &&
		!d.sawReasoning &&
		d.sawFinish &&
		d.finishReason == "stop"
placeholder

func (d *openAIChatSilentRefusalDetector) observeEventType(eventType string) {
	eventType = strings.TrimSpace(eventType)
	if eventType == "" {
		return
placeholder
	if eventType == "error" || eventType == "response.failed" {
		d.sawError = true
placeholder
	if strings.Contains(eventType, "reasoning") || strings.Contains(eventType, "reasoning_summary") {
		d.sawReasoning = true
placeholder
placeholder

func (d *openAIChatSilentRefusalDetector) observeFinishReason(reason string) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return
placeholder
	d.sawFinish = true
	d.finishReason = reason
placeholder

func (d *openAIChatSilentRefusalDetector) observeChatChoicesPayload(payload []byte) {
	choices := gjson.GetBytes(payload, "choices")
	if !choices.Exists() || !choices.IsArray() {
		return
placeholder
	for _, choice := range choices.Array() {
		if finish := choice.Get("finish_reason"); finish.Exists() {
			d.observeFinishReason(finish.String())
	placeholder
		delta := choice.Get("delta")
		if !delta.Exists() {
			continue
	placeholder
		if content := delta.Get("content"); content.Exists() && content.String() != "" {
			d.sawContent = true
	placeholder
		if delta.Get("tool_calls").Exists() {
			d.sawToolCall = true
	placeholder
		if delta.Get("function_call").Exists() {
			d.sawFunctionCall = true
	placeholder
		if delta.Get("reasoning").Exists() ||
			delta.Get("reasoning_content").Exists() ||
			delta.Get("reasoning_summary").Exists() {
			d.sawReasoning = true
	placeholder
placeholder
placeholder

func (d *openAIChatSilentRefusalDetector) observeResponsesPayload(payload []byte, eventType string) {
	switch eventType {
	case "response.output_text.delta":
		if gjson.GetBytes(payload, "delta").String() != "" {
			d.sawContent = true
	placeholder
	case "response.output_item.added":
		switch strings.TrimSpace(gjson.GetBytes(payload, "item.type").String()) {
		case "function_call":
			d.sawToolCall = true
		case "reasoning":
			d.sawReasoning = true
	placeholder
	case "response.function_call_arguments.delta":
		d.sawToolCall = true
	case "response.reasoning_summary_text.delta", "response.reasoning_summary_text.done":
		d.sawReasoning = true
	case "response.completed", "response.done":
		d.observeFinishReason("stop")
	case "response.incomplete":
		d.observeFinishReason("length")
	case "response.failed":
		d.sawError = true
placeholder

	if output := gjson.GetBytes(payload, "response.output"); output.Exists() && output.IsArray() {
		for _, item := range output.Array() {
			switch strings.TrimSpace(item.Get("type").String()) {
			case "function_call":
				d.sawToolCall = true
			case "reasoning":
				d.sawReasoning = true
			case "message":
				d.observeResponseMessageItem(item)
		placeholder
	placeholder
placeholder
placeholder

func (d *openAIChatSilentRefusalDetector) observeResponseMessageItem(item gjson.Result) {
	content := item.Get("content")
	if !content.Exists() || !content.IsArray() {
		return
placeholder
	for _, part := range content.Array() {
		if part.Get("text").String() != "" {
			d.sawContent = true
			return
	placeholder
placeholder
placeholder

func newOpenAISilentRefusalFailoverError(c *gin.Context, account *Account, upstreamRequestID string) *UpstreamFailoverError {
	accountID := int64(0)
	accountName := ""
	platform := PlatformOpenAI
	if account != nil {
		accountID = account.ID
		accountName = account.Name
		platform = account.Platform
placeholder

	setOpsUpstreamError(c, http.StatusBadGateway, openAISilentRefusalUpstreamMessage, "")
	appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
		Platform:           platform,
		AccountID:          accountID,
		AccountName:        accountName,
		UpstreamStatusCode: http.StatusBadGateway,
		UpstreamRequestID:  upstreamRequestID,
		Kind:               "failover",
		Message:            openAISilentRefusalUpstreamMessage,
placeholder)

	headers := http.Header{placeholder
	if strings.TrimSpace(upstreamRequestID) != "" {
		headers.Set("x-request-id", strings.TrimSpace(upstreamRequestID))
placeholder
	return &UpstreamFailoverError{
		StatusCode:      http.StatusBadGateway,
		ResponseBody:    openAISilentRefusalErrorBody(),
		ResponseHeaders: headers,
placeholder
placeholder

// newOpenAIResponsesEmptyCompletedFailoverError marks an empty
// response.completed terminal event as a retryable upstream anomaly. OpenAI
// Responses streams that deliver only response.created + response.completed
// with no output, no usage and no error are treated as silent upstream
// refusals rather than successful empty replies (issue #5009).
func newOpenAIResponsesEmptyCompletedFailoverError(c *gin.Context, account *Account, upstreamRequestID string) *UpstreamFailoverError {
	accountID := int64(0)
	accountName := ""
	platform := PlatformOpenAI
	if account != nil {
		accountID = account.ID
		accountName = account.Name
		platform = account.Platform
placeholder

	setOpsUpstreamError(c, http.StatusBadGateway, openAIResponsesEmptyCompletedMessage, "")
	appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
		Platform:           platform,
		AccountID:          accountID,
		AccountName:        accountName,
		UpstreamStatusCode: http.StatusBadGateway,
		UpstreamRequestID:  upstreamRequestID,
		Kind:               "failover",
		Message:            openAIResponsesEmptyCompletedMessage,
placeholder)

	headers := http.Header{placeholder
	if strings.TrimSpace(upstreamRequestID) != "" {
		headers.Set("x-request-id", strings.TrimSpace(upstreamRequestID))
placeholder
	return &UpstreamFailoverError{
		StatusCode:      http.StatusBadGateway,
		ResponseBody:    openAISilentRefusalErrorBody(),
		ResponseHeaders: headers,
placeholder
placeholder

func openAISilentRefusalErrorBody() []byte {
	body, err := json.Marshal(map[string]any{
		"error": map[string]any{
			"type":    "upstream_error",
			"code":    openAISilentRefusalErrorCode,
			"message": openAISilentRefusalUpstreamMessage,
	placeholder,
placeholder)
	if err != nil {
		return []byte(`{"error":{"type":"upstream_error","code":"openai_silent_refusal","message":"OpenAI upstream returned an empty completion stream with finish_reason=stop and no usage"placeholderplaceholder`)
placeholder
	return body
placeholder

// IsOpenAISilentRefusalErrorBody reports whether a failover body was produced
// by the OpenAI silent-refusal detector.
func IsOpenAISilentRefusalErrorBody(body []byte) bool {
	return strings.TrimSpace(gjson.GetBytes(body, "error.code").String()) == openAISilentRefusalErrorCode
placeholder

// OpenAISilentRefusalClientMessage returns the exhausted-failover client message
// for OpenAI silent refusals.
func OpenAISilentRefusalClientMessage() string {
	return openAISilentRefusalClientMessage
placeholder
