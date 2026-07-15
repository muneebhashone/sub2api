package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/apicompat"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

const (
	grokChatResponsesEndpoint = "/v1/responses"
	grokChatRawEndpoint       = "/v1/chat/completions"
)

var grokChatResponsesBridgeTopLevelFields = map[string]struct{placeholder{
	"model":                 {placeholder,
	"messages":              {placeholder,
	"stream":                {placeholder,
	"stream_options":        {placeholder,
	"max_tokens":            {placeholder,
	"max_completion_tokens": {placeholder,
	"temperature":           {placeholder,
	"top_p":                 {placeholder,
	"prompt_cache_key":      {placeholder,
	"tools":                 {placeholder,
	"tool_choice":           {placeholder,
	"functions":             {placeholder,
	"function_call":         {placeholder,
placeholder

// grokChatResponsesBridgeEligibility deliberately accepts only request shapes
// whose Chat Completions semantics are preserved by the Responses bridge.
// Everything else stays on raw Chat Completions rather than being silently
// dropped or rewritten.
func grokChatResponsesBridgeEligibility(body []byte) (bool, string) {
	var root map[string]json.RawMessage
	if err := json.Unmarshal(body, &root); err != nil || root == nil {
		return false, "invalid_json"
placeholder

	for _, field := range []string{"stop", "reasoning_effort"placeholder {
		if _, exists := root[field]; exists {
			return false, "unsupported_" + field
	placeholder
placeholder
	for _, field := range []string{"tools", "functions"placeholder {
		if raw, exists := root[field]; exists && !grokChatNullOrEmptyArray(raw) {
			return false, "unsupported_" + field
	placeholder
placeholder
	if raw, exists := root["tool_choice"]; exists && !grokChatNullOrNone(raw) {
		return false, "unsupported_tool_choice"
placeholder
	if raw, exists := root["function_call"]; exists && !grokChatNullOrNone(raw) {
		return false, "unsupported_function_call"
placeholder
	for field := range root {
		if _, supported := grokChatResponsesBridgeTopLevelFields[field]; !supported {
			return false, "unknown_field_" + field
	placeholder
placeholder

	var model string
	if raw, ok := root["model"]; !ok || json.Unmarshal(raw, &model) != nil || strings.TrimSpace(model) == "" {
		return false, "invalid_model"
placeholder

	if raw, ok := root["stream"]; ok {
		var stream *bool
		if json.Unmarshal(raw, &stream) != nil || stream == nil {
			return false, "invalid_stream"
	placeholder
placeholder
	if raw, ok := root["stream_options"]; ok {
		var options map[string]json.RawMessage
		if json.Unmarshal(raw, &options) != nil || options == nil {
			return false, "invalid_stream_options"
	placeholder
		for field, value := range options {
			if field != "include_usage" {
				return false, "unknown_stream_option_" + field
		placeholder
			var includeUsage *bool
			if json.Unmarshal(value, &includeUsage) != nil || includeUsage == nil {
				return false, "invalid_stream_include_usage"
		placeholder
	placeholder
placeholder

	for _, field := range []string{"max_tokens", "max_completion_tokens"placeholder {
		if raw, ok := root[field]; ok {
			var value *int
			if json.Unmarshal(raw, &value) != nil || value == nil || *value < 128 {
				return false, "unsafe_" + field
		placeholder
	placeholder
placeholder
	if _, hasMaxTokens := root["max_tokens"]; hasMaxTokens {
		if _, hasMaxCompletionTokens := root["max_completion_tokens"]; hasMaxCompletionTokens {
			return false, "conflicting_max_tokens"
	placeholder
placeholder
	for _, field := range []string{"temperature", "top_p"placeholder {
		if raw, ok := root[field]; ok {
			var value *float64
			if json.Unmarshal(raw, &value) != nil || value == nil {
				return false, "invalid_" + field
		placeholder
	placeholder
placeholder
	if raw, ok := root["prompt_cache_key"]; ok {
		var key string
		if json.Unmarshal(raw, &key) != nil {
			return false, "invalid_prompt_cache_key"
	placeholder
placeholder

	var messages []map[string]json.RawMessage
	rawMessages, ok := root["messages"]
	if !ok || json.Unmarshal(rawMessages, &messages) != nil || len(messages) == 0 {
		return false, "invalid_messages"
placeholder
	for _, message := range messages {
		for field := range message {
			if field != "role" && field != "content" {
				return false, "unsafe_message_field_" + field
		placeholder
	placeholder
		var role string
		if raw, exists := message["role"]; !exists || json.Unmarshal(raw, &role) != nil {
			return false, "invalid_message_role"
	placeholder
		switch role {
		case "system", "user", "assistant":
		default:
			return false, "unsupported_message_role_" + role
	placeholder
		raw, exists := message["content"]
		if !exists {
			return false, "non_text_message_content"
	placeholder
		var content string
		if json.Unmarshal(raw, &content) == nil {
			if strings.TrimSpace(content) == "" {
				return false, "empty_message_content"
		placeholder
			continue
	placeholder
		// Structured content: only allow arrays whose parts are text or
		// image_url. These are losslessly convertible to Responses input_text/
		// input_image parts, so the bridge preserves Chat Completions semantics.
		if ok, reason := grokChatStructuredContentBridgeable(raw); !ok {
			return false, reason
	placeholder
placeholder

	return true, ""
placeholder

func grokChatStructuredContentBridgeable(raw json.RawMessage) (bool, string) {
	var parts []map[string]json.RawMessage
	if err := json.Unmarshal(raw, &parts); err != nil {
		return false, "non_text_message_content"
placeholder
	if len(parts) == 0 {
		return false, "empty_message_content"
placeholder
	hasContent := false
	for _, part := range parts {
		var partType string
		rawType, ok := part["type"]
		if !ok || json.Unmarshal(rawType, &partType) != nil {
			return false, "non_text_message_content"
	placeholder
		switch strings.TrimSpace(partType) {
		case "text":
			var text string
			if raw, ok := part["text"]; ok && json.Unmarshal(raw, &text) == nil {
				if strings.TrimSpace(text) != "" {
					hasContent = true
			placeholder
		placeholder
		case "image_url", "input_image":
			hasContent = true
		default:
			return false, "unsupported_content_part_" + strings.TrimSpace(partType)
	placeholder
placeholder
	if !hasContent {
		return false, "empty_message_content"
placeholder
	return true, ""
placeholder

func grokChatNullOrEmptyArray(raw json.RawMessage) bool {
	if strings.TrimSpace(string(raw)) == "null" {
		return true
placeholder
	var values []json.RawMessage
	return json.Unmarshal(raw, &values) == nil && len(values) == 0
placeholder

func grokChatNullOrNone(raw json.RawMessage) bool {
	if strings.TrimSpace(string(raw)) == "null" {
		return true
placeholder
	var value string
	return json.Unmarshal(raw, &value) == nil && strings.EqualFold(strings.TrimSpace(value), "none")
placeholder

func grokChatCacheIntentBody(body []byte) ([]byte, error) {
	var root map[string]json.RawMessage
	if err := json.Unmarshal(body, &root); err != nil {
		return nil, err
placeholder
	for _, field := range []string{"tools", "tool_choice", "functions", "function_call"placeholder {
		delete(root, field)
placeholder
	return json.Marshal(root)
placeholder

func grokChatResponsesRuntimeEligible(upstreamModel, cacheIdentity string) bool {
	return strings.TrimSpace(upstreamModel) == "grok-4.5" && strings.TrimSpace(cacheIdentity) != ""
placeholder

// forwardGrokChatCompletionsViaResponses converts a strictly compatible Chat
// request into xAI Responses format and reuses the established Responses-to-
// Chat response translators. It intentionally does not run the Codex OAuth
// transform because Grok CLI is a separate upstream protocol.
func (s *OpenAIGatewayService) forwardGrokChatCompletionsViaResponses(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	body []byte,
	promptCacheKey string,
	defaultMappedModel string,
) (*OpenAIForwardResult, error) {
	startTime := time.Now()

	var chatReq apicompat.ChatCompletionsRequest
	if err := json.Unmarshal(body, &chatReq); err != nil {
		return nil, fmt.Errorf("parse grok chat completions request: %w", err)
placeholder
	originalModel := chatReq.Model
	clientStream := chatReq.Stream
	billingModel := resolveOpenAIForwardModel(account, originalModel, defaultMappedModel)
	upstreamModel := normalizeOpenAIModelForUpstream(account, billingModel)
	cacheIdentity := resolveGrokCacheIdentity(c, body, promptCacheKey, upstreamModel)
	// Image inputs must go through the Responses bridge: the raw Chat
	// Completions path cannot forward image_url parts to Grok's native vision
	// for non-composer models, so they would be silently dropped. Route them to
	// Responses even when no prompt-cache identity is available.
	hasImageInput := openAIJSONValueMayContainImageInput(gjson.GetBytes(body, "messages"))
	if !grokChatResponsesRuntimeEligible(upstreamModel, cacheIdentity) && (!hasImageInput || strings.TrimSpace(upstreamModel) != "grok-4.5") {
		return s.forwardAsRawChatCompletions(ctx, c, account, body, defaultMappedModel)
placeholder

	responsesReq, err := apicompat.ChatCompletionsToResponses(&chatReq)
	if err != nil {
		return nil, fmt.Errorf("convert grok chat completions to responses: %w", err)
placeholder
	responsesReq.Model = upstreamModel
	responsesReq.Stream = true
	// These fields are useful to Codex but are not needed by the Grok CLI
	// protocol. Keep the bridge request as close as possible to native Grok.
	responsesReq.Include = nil
	responsesReq.Store = nil

	responsesBody, err := json.Marshal(responsesReq)
	if err != nil {
		return nil, fmt.Errorf("marshal grok responses bridge request: %w", err)
placeholder
	responsesBody, err = patchGrokResponsesBody(responsesBody, upstreamModel)
	if err != nil {
		return nil, fmt.Errorf("patch grok responses bridge request: %w", err)
placeholder
	intentBody, err := grokChatCacheIntentBody(body)
	if err != nil {
		return nil, fmt.Errorf("normalize grok responses bridge tool intent: %w", err)
placeholder
	responsesBody, err = applyGrokResponsesCacheIdentity(responsesBody, intentBody, cacheIdentity, true)
	if err != nil {
		return nil, fmt.Errorf("apply grok responses bridge cache identity: %w", err)
placeholder

	updatedBody, policyErr := s.applyOpenAIFastPolicyToBody(ctx, account, upstreamModel, responsesBody)
	if policyErr != nil {
		var blocked *OpenAIFastBlockedError
		if errors.As(policyErr, &blocked) {
			MarkOpsClientBusinessLimited(c, OpsClientBusinessLimitedReasonLocalPolicyDenied)
			writeChatCompletionsError(c, http.StatusForbidden, "permission_error", blocked.Message)
	placeholder
		return nil, policyErr
placeholder
	responsesBody = updatedBody

	token, _, err := s.getRequestCredential(ctx, c, account)
	if err != nil {
		return nil, fmt.Errorf("get grok access token: %w", err)
placeholder
	upstreamCtx, releaseUpstreamCtx := detachUpstreamContext(ctx)
	upstreamReq, err := buildGrokResponsesRequest(upstreamCtx, c, account, responsesBody, token, cacheIdentity, s.cfg)
	releaseUpstreamCtx()
	if err != nil {
		return nil, fmt.Errorf("build grok responses bridge request: %w", err)
placeholder
	SetActualOpenAIUpstreamEndpoint(c, grokChatResponsesEndpoint)

	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
placeholder
	resp, err := s.httpUpstream.Do(upstreamReq, proxyURL, account.ID, account.Concurrency)
	if err != nil {
		return nil, s.handleOpenAIUpstreamTransportError(ctx, c, account, err, false)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	if resp.StatusCode >= http.StatusBadRequest {
		respBody, upstreamMsg := s.readOpenAIUpstreamError(resp)
		if upstreamMsg == "" {
			upstreamMsg = fmt.Sprintf("xAI upstream returned status %d", resp.StatusCode)
	placeholder
		appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
			Platform:           account.Platform,
			AccountID:          account.ID,
			AccountName:        account.Name,
			UpstreamStatusCode: resp.StatusCode,
			UpstreamRequestID:  firstNonEmpty(resp.Header.Get("x-request-id"), resp.Header.Get("xai-request-id")),
			Kind:               "failover",
			Message:            upstreamMsg,
	placeholder)
		s.handleGrokAccountUpstreamError(ctx, account, resp.StatusCode, resp.Header, respBody)
		if s.shouldFailoverUpstreamError(resp.StatusCode) {
			return nil, &UpstreamFailoverError{
				StatusCode:             resp.StatusCode,
				ResponseBody:           respBody,
				ResponseHeaders:        resp.Header.Clone(),
				RetryableOnSameAccount: account.IsPoolMode() && account.IsPoolModeRetryableStatus(resp.StatusCode),
		placeholder
	placeholder
		return s.handleChatCompletionsErrorResponse(resp, c, account, billingModel)
placeholder

	s.updateGrokUsageFromResponse(ctx, account, resp.Header, resp.StatusCode)

	var result *OpenAIForwardResult
	if clientStream {
		result, err = s.handleChatStreamingResponse(resp, c, account, originalModel, billingModel, upstreamModel, startTime, len(body))
placeholder else {
		result, err = s.handleChatBufferedStreamingResponse(resp, c, account, originalModel, billingModel, upstreamModel, startTime)
placeholder
	if result != nil {
		result.UpstreamEndpoint = grokChatResponsesEndpoint
		result.ResponseHeaders = resp.Header.Clone()
		if result.RequestID == "" {
			result.RequestID = firstNonEmpty(resp.Header.Get("x-request-id"), resp.Header.Get("xai-request-id"))
	placeholder
		result.ReasoningEffort = extractOpenAIReasoningEffortFromBody(body, upstreamModel, billingModel, originalModel)
placeholder
	return result, err
placeholder
