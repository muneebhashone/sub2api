package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/apicompat"
	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	grokComposerImageBridgeVisionModel     = "grok-build-0.1"
	grokComposerImageBridgeMaxOutputTokens = 512
	// grokUpstreamUserAgent lives in grok_upstream_headers.go (shared with TLS header helpers).
	grokCLIVersion                   = xai.CLIClientVersion
	grokDefaultResponsesModel        = "grok-4.5"
	grokRateLimitFallbackCooldown    = 2 * time.Minute
	grokRateLimitRepeatCooldown      = 10 * time.Minute
	grokRateLimitSustainedCooldown   = 30 * time.Minute
	grokRateLimitMaxAdaptiveCooldown = time.Hour
	grokRateLimitBackoffQuietPeriod  = time.Hour
)

func (s *OpenAIGatewayService) forwardGrokResponses(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	body []byte,
	originalModel string,
	reqStream bool,
	startTime time.Time,
) (*OpenAIForwardResult, error) {
	if account.Type != AccountTypeOAuth && account.Type != AccountTypeAPIKey {
		return nil, fmt.Errorf("grok account type %s is not supported by Responses forwarding", account.Type)
placeholder

	upstreamModel := account.GetMappedModel(originalModel)
	if strings.TrimSpace(upstreamModel) == "" {
		upstreamModel = grokDefaultResponsesModel
placeholder
	// Account mappings are optional. Canonicalize client aliases even when the
	// account has no model_mapping, matching the Chat Completions path and xAI's
	// actual Responses model IDs.
	upstreamModel = xai.ResolveGrokTextResponsesModelID(upstreamModel, grokDefaultResponsesModel)
	if isGrokImageGenerationModel(upstreamModel) {
		return nil, fmt.Errorf("model %s is an image model and is not available on the Responses endpoint; use /v1/images/generations instead", upstreamModel)
placeholder
	patchedBody, clientToolMapping, err := patchGrokResponsesBodyWithClientTools(body, upstreamModel)
	if err != nil {
		setOpsUpstreamError(c, http.StatusBadRequest, err.Error(), "")
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"type": "invalid_request_error", "message": err.Error(), "param": "tools",
	placeholderplaceholder)
		return nil, err
placeholder
	setGrokResponsesClientToolMapping(c, clientToolMapping)
	// OpenAI /responses/compact is not a native xAI endpoint. Convert it into a
	// normal Grok Responses turn that asks for a structured summary, then map the
	// reply back to an OpenAI compaction item on the way out.
	if isOpenAIResponsesCompactPath(c) {
		patchedBody, err = buildGrokCompactRequestBody(patchedBody)
		if err != nil {
			return nil, err
	placeholder
placeholder
	// Derive the identity from the request xAI will actually see. This makes
	// Codex Responses Lite additional_tools part of the stable tool prefix.
	cacheIdentity := resolveGrokCacheIdentity(c, patchedBody, "", upstreamModel)
	mixedCacheIntentBody := append([]byte(nil), patchedBody...)
	patchedBody, err = applyGrokResponsesCacheIdentity(patchedBody, body, cacheIdentity, account.IsGrokOAuth())
	if err != nil {
		return nil, fmt.Errorf("apply grok prompt cache identity: %w", err)
placeholder
	// Free OAuth + client function tools: reuse Messages mixed-tools cache route
	// (append web_search/x_search so xAI does not force non-cacheable build-free).
	patchedBody, err = applyGrokFreeRequestToolCacheRoute(c, patchedBody, mixedCacheIntentBody, account, cacheIdentity)
	if err != nil {
		return nil, fmt.Errorf("apply grok Free function-tool cache route: %w", err)
placeholder

	token, _, err := s.getRequestCredential(ctx, c, account)
	if err != nil {
		return nil, err
placeholder

	upstreamCtx, releaseUpstreamCtx := detachUpstreamContext(ctx)
	defer releaseUpstreamCtx()

	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
placeholder

	upstreamStart := time.Now()
	var resp *http.Response
	for attempt := 0; ; attempt++ {
		upstreamReq, buildErr := buildGrokResponsesRequest(upstreamCtx, c, account, patchedBody, token, cacheIdentity, s.cfg, s.settingService)
		if buildErr != nil {
			return nil, buildErr
	placeholder

		resp, err = s.httpUpstream.Do(upstreamReq, proxyURL, account.ID, account.Concurrency)
		SetOpsLatencyMs(c, OpsUpstreamLatencyMsKey, time.Since(upstreamStart).Milliseconds())
		if err != nil {
			return nil, s.handleOpenAIUpstreamTransportError(ctx, c, account, err, false)
	placeholder

		// xAI can reject encrypted reasoning copied from a response produced under
		// another account or cache identity. Retry once with the same routing and
		// credential after removing only the rejected encrypted reasoning payload.
		if attempt > 0 || resp.StatusCode != http.StatusBadRequest {
			break
	placeholder
		respBody := s.readUpstreamErrorBody(resp)
		if resp.Body != nil {
			_ = resp.Body.Close()
	placeholder
		if !isGrokInvalidEncryptedContentResponse(resp.StatusCode, respBody) {
			resp.Body = io.NopCloser(bytes.NewReader(respBody))
			break
	placeholder

		retryBody, changed, trimErr := trimGrokInvalidEncryptedContentRetryBody(patchedBody)
		if trimErr != nil {
			return nil, fmt.Errorf("prepare Grok invalid encrypted_content retry: %w", trimErr)
	placeholder
		if !changed {
			resp.Body = io.NopCloser(bytes.NewReader(respBody))
			break
	placeholder

		patchedBody = retryBody
		slog.Info("grok_invalid_encrypted_content_retry", "account_id", account.ID, "cache_identity_present", cacheIdentity != "")
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	if resp.StatusCode >= 400 {
		respBody := s.readUpstreamErrorBody(resp)
		resp.Body = io.NopCloser(bytes.NewReader(respBody))
		upstreamMsg := sanitizeUpstreamErrorMessage(extractUpstreamErrorMessage(respBody))
		if upstreamMsg == "" {
			upstreamMsg = fmt.Sprintf("xAI upstream returned status %d", resp.StatusCode)
	placeholder
		kind := "http_error"
		if s.shouldFailoverGrokUpstreamError(resp.StatusCode, respBody) {
			kind = "failover"
	placeholder
		appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
			Platform:           account.Platform,
			AccountID:          account.ID,
			AccountName:        account.Name,
			UpstreamStatusCode: resp.StatusCode,
			UpstreamRequestID:  firstNonEmpty(resp.Header.Get("x-request-id"), resp.Header.Get("xai-request-id")),
			Kind:               kind,
			Message:            upstreamMsg,
	placeholder)
		errCtx := withGrokTeamRateLimitModel(ctx, upstreamModel)
		s.handleGrokAccountUpstreamError(errCtx, account, resp.StatusCode, resp.Header, respBody)
		// 429 / free-usage: stamp team+model cool so sibling accounts skip this model.
		if resp.StatusCode == http.StatusTooManyRequests ||
			classifyGrokUpstreamFailure(resp.StatusCode, respBody, upstreamModel).Class == GrokFailureFreeUsage {
			markGrokTeamModelRateLimit(account, upstreamModel, resolveGrokTeamRateLimitUntil(time.Now().Add(grokTeamRateLimitDefaultTTL), time.Now()))
	placeholder
		if s.shouldFailoverGrokUpstreamError(resp.StatusCode, respBody) {
			return nil, &UpstreamFailoverError{
				StatusCode:             resp.StatusCode,
				ResponseBody:           respBody,
				ResponseHeaders:        resp.Header.Clone(),
				RetryableOnSameAccount: account.IsPoolMode() && account.IsPoolModeRetryableStatus(resp.StatusCode),
		placeholder
	placeholder
		return s.handleErrorResponse(ctx, resp, c, account, patchedBody, upstreamModel)
placeholder

	// Attach model so rate-limit snapshots can fan out a team+model cool.
	stateCtx := withGrokTeamRateLimitModel(ctx, upstreamModel)
	s.updateGrokUsageFromResponse(stateCtx, account, resp.Header, resp.StatusCode)

	var usage *OpenAIUsage
	var firstTokenMs *int
	responseID := ""
	searchCount := 0
	imageCount := 0
	var imageOutputSizes []string
	if reqStream {
		maxLineSize := defaultMaxLineSize
		if s.cfg != nil && s.cfg.Gateway.MaxLineSize > 0 {
			maxLineSize = s.cfg.Gateway.MaxLineSize
	placeholder
		resp.Body = newGrokResponsesBillingPingFilterBody(resp.Body, account, maxLineSize)
		if hasGrokResponsesClientToolMapping(clientToolMapping) {
			resp.Body = newGrokResponsesClientToolStreamBody(resp.Body, clientToolMapping, maxLineSize)
	placeholder
		streamResult, err := s.handleStreamingResponse(ctx, resp, c, account, startTime, originalModel, upstreamModel)
		if err != nil {
			return nil, err
	placeholder
		usage = streamResult.usage
		firstTokenMs = streamResult.firstTokenMs
		responseID = strings.TrimSpace(streamResult.responseID)
		searchCount = streamResult.searchCount
		imageCount = streamResult.imageCount
		imageOutputSizes = streamResult.imageOutputSizes
placeholder else {
		nonStreamResult, err := s.handleNonStreamingResponse(ctx, resp, c, account, originalModel, upstreamModel)
		if err != nil {
			return nil, err
	placeholder
		usage = nonStreamResult.usage
		responseID = strings.TrimSpace(nonStreamResult.responseID)
		searchCount = nonStreamResult.searchCount
		imageCount = nonStreamResult.imageCount
		imageOutputSizes = nonStreamResult.imageOutputSizes
placeholder

	if usage == nil {
		usage = &OpenAIUsage{placeholder
placeholder
	reasoningEffort := extractOpenAIReasoningEffortFromBody(patchedBody, originalModel)
	result := &OpenAIForwardResult{
		RequestID:       firstNonEmpty(resp.Header.Get("x-request-id"), resp.Header.Get("xai-request-id")),
		ResponseID:      responseID,
		Usage:           *usage,
		Model:           originalModel,
		UpstreamModel:   upstreamModel,
		ReasoningEffort: reasoningEffort,
		Stream:          reqStream,
		OpenAIWSMode:    false,
		ResponseHeaders: resp.Header.Clone(),
		Duration:        time.Since(startTime),
		FirstTokenMs:    firstTokenMs,
placeholder
	// Propagate search/image counters from the shared Responses handler — without
	// this, stream/JSON counting runs but search_price_per_1k / image bills never apply.
	if searchCount > 0 {
		result.SearchCount = searchCount
placeholder
	if imageCount > 0 {
		result.ImageCount = imageCount
		result.ImageOutputSizes = imageOutputSizes
placeholder
	return result, nil
placeholder

func isGrokInvalidEncryptedContentResponse(statusCode int, body []byte) bool {
	if statusCode != http.StatusBadRequest {
		return false
placeholder

	// xAI has used both flat and nested error envelopes:
	//   {"code":"invalid-argument","error":"Could not decrypt the provided encrypted_content."placeholder
	//   {"error":{"message":"Could not decrypt the provided encrypted_content."placeholderplaceholder
	code := strings.TrimSpace(gjson.GetBytes(body, "code").String())
	message := ""
	errNode := gjson.GetBytes(body, "error")
	switch {
	case errNode.Type == gjson.String:
		message = errNode.String()
	case errNode.IsObject():
		message = firstNonEmpty(errNode.Get("message").String(), errNode.Get("error").String())
		if code == "" {
			code = strings.TrimSpace(errNode.Get("code").String())
	placeholder
	default:
		message = gjson.GetBytes(body, "message").String()
placeholder
	normalizedMessage := strings.ToLower(strings.TrimSpace(message))
	if normalizedMessage == "" {
		return false
placeholder

	if strings.EqualFold(code, "invalid_encrypted_content") {
		return true
placeholder
	// Keep the official xAI flat-code gate so unrelated 400s are not retried.
	if !strings.EqualFold(code, "invalid-argument") && code != "" {
		return false
placeholder
	// Nested OpenAI-style envelopes may omit top-level code; require decrypt text.
	if code == "" && !strings.Contains(normalizedMessage, "decrypt") {
		return false
placeholder
	return strings.Contains(normalizedMessage, "encrypted_content") &&
		(strings.Contains(normalizedMessage, "decrypt") ||
			strings.Contains(normalizedMessage, "unmodified"))
placeholder

// requestHasGrokEncryptedReasoning reports whether the outbound Responses body
// still carries reasoning.encrypted_content that can be stripped for retry.
func requestHasGrokEncryptedReasoning(body []byte) bool {
	input := gjson.GetBytes(body, "input")
	if !input.Exists() {
		return false
placeholder
	items := input.Array()
	if input.IsObject() {
		items = []gjson.Result{inputplaceholder
placeholder
	for _, item := range items {
		if strings.TrimSpace(item.Get("type").String()) != "reasoning" {
			continue
	placeholder
		enc := item.Get("encrypted_content")
		if enc.Exists() && enc.Type != gjson.Null && strings.TrimSpace(enc.String()) != "" {
			return true
	placeholder
placeholder
	return false
placeholder

type grokEncryptedContentStripRetriedKey struct{placeholder

func markGrokEncryptedContentStripRetried(ctx context.Context) context.Context {
	return context.WithValue(ctx, grokEncryptedContentStripRetriedKey{placeholder, true)
placeholder

func grokEncryptedContentStripRetried(ctx context.Context) bool {
	v, _ := ctx.Value(grokEncryptedContentStripRetriedKey{placeholder).(bool)
	return v
placeholder

// stripAnthropicThinkingSignatures removes thinking.signature from Claude
// history so a different Grok OAuth account can accept multi-turn tool
// continuations after decrypt failures. Returns ok=false when nothing changed.
func stripAnthropicThinkingSignatures(body []byte) ([]byte, bool) {
	if len(body) == 0 || !bytes.Contains(body, []byte(`"signature"`)) {
		return body, false
placeholder
	var req map[string]any
	if err := json.Unmarshal(body, &req); err != nil {
		return body, false
placeholder
	messages, ok := req["messages"].([]any)
	if !ok || len(messages) == 0 {
		return body, false
placeholder
	changed := false
	for _, rawMsg := range messages {
		msg, ok := rawMsg.(map[string]any)
		if !ok {
			continue
	placeholder
		content, ok := msg["content"].([]any)
		if !ok {
			continue
	placeholder
		for _, rawBlock := range content {
			block, ok := rawBlock.(map[string]any)
			if !ok {
				continue
		placeholder
			if typ, _ := block["type"].(string); typ != "thinking" {
				continue
		placeholder
			if _, has := block["signature"]; has {
				delete(block, "signature")
				changed = true
		placeholder
	placeholder
placeholder
	if !changed {
		return body, false
placeholder
	out, err := json.Marshal(req)
	if err != nil {
		return body, false
placeholder
	return out, true
placeholder

func trimGrokInvalidEncryptedContentRetryBody(body []byte) ([]byte, bool, error) {
	input := gjson.GetBytes(body, "input")
	items := input.Array()
	if input.IsObject() {
		items = []gjson.Result{inputplaceholder
placeholder

	hasEncryptedReasoning := false
	for _, item := range items {
		if strings.TrimSpace(item.Get("type").String()) == "reasoning" && item.Get("encrypted_content").Exists() {
			hasEncryptedReasoning = true
			break
	placeholder
placeholder
	if !hasEncryptedReasoning {
		return body, false, nil
placeholder

	var requestBody map[string]any
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	if err := decoder.Decode(&requestBody); err != nil {
		return nil, false, err
placeholder
	if !trimOpenAIEncryptedReasoningItems(requestBody) {
		return body, false, nil
placeholder

	retryBody, err := marshalOpenAIUpstreamJSON(requestBody)
	if err != nil {
		return nil, false, err
placeholder
	return retryBody, true, nil
placeholder

func patchGrokResponsesBody(body []byte, upstreamModel string) ([]byte, error) {
	return patchGrokResponsesBodyBase(body, upstreamModel)
placeholder

func patchGrokResponsesBodyWithClientTools(body []byte, upstreamModel string) ([]byte, apicompat.ResponsesClientToolMapping, error) {
	if !json.Valid(body) {
		return nil, apicompat.ResponsesClientToolMapping{placeholder, fmt.Errorf("invalid json request body")
placeholder
	promoted, err := sanitizeGrokResponsesInput(body)
	if err != nil {
		return nil, apicompat.ResponsesClientToolMapping{placeholder, err
placeholder
	adapted, mapping, err := adaptGrokResponsesClientTools(promoted)
	if err != nil {
		return nil, apicompat.ResponsesClientToolMapping{placeholder, err
placeholder
	patched, err := patchGrokResponsesBodyBase(adapted, upstreamModel)
	if err != nil {
		return nil, apicompat.ResponsesClientToolMapping{placeholder, err
placeholder
	return patched, mapping, nil
placeholder

func patchGrokResponsesBodyBase(body []byte, upstreamModel string) ([]byte, error) {
	if !json.Valid(body) {
		return nil, fmt.Errorf("invalid json request body")
placeholder
	// sjson may reuse the input backing array; keep the caller's request bytes
	// unchanged because the same body can be inspected for billing/retry paths.
	out, err := sjson.SetBytes(append([]byte(nil), body...), "model", upstreamModel)
	if err != nil {
		return nil, err
placeholder
	out, err = normalizeGrokResponsesReasoningEffort(out, upstreamModel)
	if err != nil {
		return nil, err
placeholder
	out, err = sanitizeGrokResponsesModelCapabilities(out, upstreamModel)
	if err != nil {
		return nil, err
placeholder
	for _, unsupportedField := range []string{"prompt_cache_retention", "safety_identifier"placeholder {
		if gjson.GetBytes(out, unsupportedField).Exists() {
			out, err = sjson.DeleteBytes(out, unsupportedField)
			if err != nil {
				return nil, err
		placeholder
	placeholder
placeholder
	if strings.EqualFold(upstreamModel, "grok-4.5") {
		for _, unsupportedField := range []string{"presence_penalty", "presencePenalty", "frequency_penalty", "frequencyPenalty", "stop"placeholder {
			if gjson.GetBytes(out, unsupportedField).Exists() {
				out, err = sjson.DeleteBytes(out, unsupportedField)
				if err != nil {
					return nil, err
			placeholder
		placeholder
	placeholder
placeholder
	if grokModelRejectsLogprobs(upstreamModel) {
		for _, unsupportedField := range []string{"logprobs", "top_logprobs"placeholder {
			if gjson.GetBytes(out, unsupportedField).Exists() {
				out, err = sjson.DeleteBytes(out, unsupportedField)
				if err != nil {
					return nil, err
			placeholder
		placeholder
	placeholder
placeholder
	out, err = sanitizeGrokResponsesUnsupportedFields(out)
	if err != nil {
		return nil, err
placeholder
	out, err = convertOpenAICompactInputsForGrok(out)
	if err != nil {
		return nil, err
placeholder
	out, err = sanitizeGrokResponsesInput(out)
	if err != nil {
		return nil, err
placeholder
	out, err = sanitizeGrokReasoningNullContent(out)
	if err != nil {
		return nil, err
placeholder
	out, err = sanitizeGrokResponsesTools(out)
	if err != nil {
		return nil, err
placeholder
	return out, nil
placeholder

// xAI's Grok 4.20 family and newer models do not support OpenAI's logprobs
// fields. Remove them before egress instead of forwarding a request the
// upstream rejects. Older Grok models retain the fields for compatibility.
func grokModelRejectsLogprobs(model string) bool {
	model = strings.ToLower(strings.TrimSpace(model))
	if slash := strings.LastIndex(model, "/"); slash >= 0 {
		model = strings.TrimSpace(model[slash+1:])
placeholder
	return strings.HasPrefix(model, "grok-4.20")
placeholder

func sanitizeGrokResponsesModelCapabilities(body []byte, upstreamModel string) ([]byte, error) {
	if !grokModelRejectsReasoningEffort(upstreamModel) {
		return body, nil
placeholder

	out := body
	for _, field := range []string{"reasoning", "reasoning_effort", "reasoningEffort"placeholder {
		if !gjson.GetBytes(out, field).Exists() {
			continue
	placeholder
		var err error
		out, err = sjson.DeleteBytes(out, field)
		if err != nil {
			return nil, fmt.Errorf("remove unsupported Grok Composer %s: %w", field, err)
	placeholder
placeholder
	return out, nil
placeholder

func grokModelRejectsReasoningEffort(model string) bool {
	model = strings.TrimSpace(strings.ToLower(model))
	if slash := strings.LastIndex(model, "/"); slash >= 0 {
		model = strings.TrimSpace(model[slash+1:])
placeholder
	switch model {
	case "grok-composer", "grok-composer-2.5-fast", "composer-2.5":
		return true
	default:
		return false
placeholder
placeholder

func normalizeGrokResponsesReasoningEffort(body []byte, upstreamModel string) ([]byte, error) {
	supportsEffort := grokSupportsReasoningEffort(upstreamModel)
	out := body
	var err error
	for _, field := range []string{"reasoning.effort", "reasoning_effort"placeholder {
		value := gjson.GetBytes(out, field)
		if !value.Exists() {
			continue
	placeholder
		normalized, keep := normalizeGrokReasoningEffortValue(value.String())
		if !supportsEffort || !keep {
			out, err = sjson.DeleteBytes(out, field)
	placeholder else {
			out, err = sjson.SetBytes(out, field, normalized)
	placeholder
		if err != nil {
			return nil, fmt.Errorf("normalize Grok reasoning field %s: %w", field, err)
	placeholder
placeholder
	if camel := gjson.GetBytes(out, "reasoningEffort"); camel.Exists() {
		normalized, keep := normalizeGrokReasoningEffortValue(camel.String())
		out, err = sjson.DeleteBytes(out, "reasoningEffort")
		if err != nil {
			return nil, fmt.Errorf("remove Grok reasoningEffort: %w", err)
	placeholder
		if supportsEffort && keep && !gjson.GetBytes(out, "reasoning_effort").Exists() {
			out, err = sjson.SetBytes(out, "reasoning_effort", normalized)
			if err != nil {
				return nil, fmt.Errorf("set Grok reasoning_effort: %w", err)
		placeholder
	placeholder
placeholder
	if reasoning := gjson.GetBytes(out, "reasoning"); reasoning.Exists() && reasoning.IsObject() && len(reasoning.Map()) == 0 {
		out, err = sjson.DeleteBytes(out, "reasoning")
		if err != nil {
			return nil, fmt.Errorf("remove empty Grok reasoning: %w", err)
	placeholder
placeholder
	return out, nil
placeholder

func normalizeGrokChatReasoningEffort(body []byte, upstreamModel string) ([]byte, error) {
	raw := strings.TrimSpace(gjson.GetBytes(body, "reasoning_effort").String())
	if raw == "" {
		raw = strings.TrimSpace(gjson.GetBytes(body, "reasoningEffort").String())
placeholder
	normalized, keep := normalizeGrokReasoningEffortValue(raw)
	keep = keep && grokSupportsReasoningEffort(upstreamModel)
	out := body
	var err error
	if gjson.GetBytes(out, "reasoningEffort").Exists() {
		out, err = sjson.DeleteBytes(out, "reasoningEffort")
		if err != nil {
			return nil, err
	placeholder
placeholder
	if !keep {
		if gjson.GetBytes(out, "reasoning_effort").Exists() {
			out, err = sjson.DeleteBytes(out, "reasoning_effort")
	placeholder
		return out, err
placeholder
	out, err = sjson.SetBytes(out, "reasoning_effort", normalized)
	return out, err
placeholder

func normalizeGrokReasoningEffortValue(raw string) (string, bool) {
	value := strings.NewReplacer("-", "", "_", "", " ", "").Replace(strings.ToLower(strings.TrimSpace(raw)))
	switch value {
	case "none", "low", "medium", "high":
		return value, true
	case "minimal":
		return "low", true
	case "xhigh", "extrahigh", "max", "ultra":
		return "high", true
	default:
		return "", false
placeholder
placeholder

func grokSupportsReasoningEffort(model string) bool {
	model = strings.ToLower(xai.StripGrokProviderPrefix(strings.TrimSpace(model)))
	switch model {
	case xai.DefaultTextModel, "grok-4.5-latest", "grok-4.6", "grok-4.6-latest",
		"grok-4.3", "grok-4.3-latest",
		"grok-3-mini", "grok-3-mini-fast", "grok-4.20-0309-reasoning",
		"grok-4.20-reasoning", "grok-4.20-multi-agent-0309":
		return true
	default:
		return false
placeholder
placeholder

var grokResponsesUnsupportedRecursiveFields = map[string]struct{placeholder{
	"external_web_access": {placeholder,
placeholder

func sanitizeGrokResponsesUnsupportedFields(body []byte) ([]byte, error) {
	if !bytes.Contains(body, []byte(`"external_web_access"`)) {
		return body, nil
placeholder

	var payload any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
placeholder
	if !deleteJSONFields(payload, grokResponsesUnsupportedRecursiveFields) {
		return body, nil
placeholder
	return json.Marshal(payload)
placeholder

func deleteJSONFields(value any, fields map[string]struct{placeholder) bool {
	switch typed := value.(type) {
	case map[string]any:
		changed := false
		for field := range fields {
			if _, ok := typed[field]; ok {
				delete(typed, field)
				changed = true
		placeholder
	placeholder
		for _, child := range typed {
			if deleteJSONFields(child, fields) {
				changed = true
		placeholder
	placeholder
		return changed
	case []any:
		changed := false
		for _, child := range typed {
			if deleteJSONFields(child, fields) {
				changed = true
		placeholder
	placeholder
		return changed
	default:
		return false
placeholder
placeholder

// additional_tools is a Codex/Responses Lite private input carrier. xAI's
// Responses schema rejects the carrier itself, but accepts supported tools at
// the top level. Preserve top-level order, append newly discovered tools in
// carrier order, then let sanitizeGrokResponsesTools filter unsupported types.
func sanitizeGrokResponsesInput(body []byte) ([]byte, error) {
	if !bytes.Contains(body, []byte(`"additional_tools"`)) {
		return body, nil
placeholder
	input := gjson.GetBytes(body, "input")
	if !input.Exists() || !input.IsArray() {
		return body, nil
placeholder

	rawItems := input.Array()
	filtered := make([]json.RawMessage, 0, len(rawItems))
	topLevelTools := gjson.GetBytes(body, "tools")
	mergedTools := make([]json.RawMessage, 0)
	seenTools := make(map[string]struct{placeholder)
	appendTool := func(tool gjson.Result) bool {
		key := grokResponsesToolDedupKey(tool)
		if _, exists := seenTools[key]; exists {
			return false
	placeholder
		seenTools[key] = struct{placeholder{placeholder
		mergedTools = append(mergedTools, json.RawMessage(tool.Raw))
		return true
placeholder
	if topLevelTools.IsArray() {
		for _, tool := range topLevelTools.Array() {
			seenTools[grokResponsesToolDedupKey(tool)] = struct{placeholder{placeholder
			mergedTools = append(mergedTools, json.RawMessage(tool.Raw))
	placeholder
placeholder

	promoted := false
	for _, item := range rawItems {
		if strings.TrimSpace(item.Get("type").String()) == "additional_tools" {
			tools := item.Get("tools")
			if tools.IsArray() {
				for _, tool := range tools.Array() {
					if appendTool(tool) {
						promoted = true
				placeholder
			placeholder
		placeholder
			continue
	placeholder
		filtered = append(filtered, json.RawMessage(item.Raw))
placeholder
	if len(filtered) == len(rawItems) {
		return body, nil
placeholder
	encoded, err := json.Marshal(filtered)
	if err != nil {
		return nil, err
placeholder
	body, err = sjson.SetRawBytes(body, "input", encoded)
	if err != nil || !promoted {
		return body, err
placeholder
	encodedTools, err := json.Marshal(mergedTools)
	if err != nil {
		return nil, err
placeholder
	return sjson.SetRawBytes(body, "tools", encodedTools)
placeholder

func grokResponsesToolDedupKey(tool gjson.Result) string {
	toolType := strings.TrimSpace(tool.Get("type").String())
	if toolType != "" {
		if name := strings.TrimSpace(tool.Get("name").String()); name != "" {
			return "type:" + toolType + "\x00name:" + name
	placeholder
		if toolType == "mcp" {
			if label := strings.TrimSpace(tool.Get("server_label").String()); label != "" {
				return "type:mcp\x00server_label:" + label
		placeholder
	placeholder
placeholder
	return "json:" + normalizeCompatSeedJSON(json.RawMessage(tool.Raw))
placeholder

// sanitizeGrokReasoningNullContent 删除 reasoning 项中的 "content": null。
// xAI 的 untagged enum 反序列化器拒收该字段，返回 422。
func sanitizeGrokReasoningNullContent(body []byte) ([]byte, error) {
	input := gjson.GetBytes(body, "input")
	if !input.Exists() || !input.IsArray() {
		return body, nil
placeholder

	items := input.Array()
	changed := false
	for i := len(items) - 1; i >= 0; i-- {
		item := items[i]
		if strings.TrimSpace(item.Get("type").String()) != "reasoning" {
			continue
	placeholder
		contentResult := item.Get("content")
		if contentResult.Exists() && contentResult.Type == gjson.Null {
			var err error
			body, err = sjson.DeleteBytes(body, fmt.Sprintf("input.%d.content", i))
			if err != nil {
				return nil, err
		placeholder
			changed = true
	placeholder
placeholder
	_ = changed
	return body, nil
placeholder

var grokResponsesSupportedToolTypes = map[string]struct{placeholder{
	"code_execution":     {placeholder,
	"code_interpreter":   {placeholder,
	"collections_search": {placeholder,
	"file_search":        {placeholder,
	"function":           {placeholder,
	"mcp":                {placeholder,
	"shell":              {placeholder,
	"web_search":         {placeholder,
	"x_search":           {placeholder,
placeholder

func sanitizeGrokResponsesTools(body []byte) ([]byte, error) {
	tools := gjson.GetBytes(body, "tools")
	if !tools.Exists() {
		if gjson.GetBytes(body, "tool_choice").Exists() {
			return sjson.DeleteBytes(body, "tool_choice")
	placeholder
		return body, nil
placeholder
	if !tools.IsArray() {
		return body, nil
placeholder

	rawTools := tools.Array()
	filteredTools := make([]json.RawMessage, 0, len(rawTools))
	toolsChanged := false
	for _, tool := range rawTools {
		toolType := strings.TrimSpace(tool.Get("type").String())
		if _, ok := grokResponsesSupportedToolTypes[toolType]; ok {
			raw := json.RawMessage(tool.Raw)
			if toolType == "function" && (!tool.Get("parameters").Exists() || tool.Get("parameters").Type == gjson.Null) {
				var payload map[string]any
				if err := json.Unmarshal(raw, &payload); err != nil {
					return nil, err
			placeholder
				payload["parameters"] = map[string]any{"type": "object", "properties": map[string]any{placeholderplaceholder
				encoded, err := json.Marshal(payload)
				if err != nil {
					return nil, err
			placeholder
				raw = encoded
				toolsChanged = true
		placeholder
			filteredTools = append(filteredTools, raw)
	placeholder
placeholder

	var err error
	if len(filteredTools) != len(rawTools) || toolsChanged {
		if len(filteredTools) == 0 {
			body, err = sjson.DeleteBytes(body, "tools")
	placeholder else {
			var encoded []byte
			encoded, err = json.Marshal(filteredTools)
			if err != nil {
				return nil, err
		placeholder
			body, err = sjson.SetRawBytes(body, "tools", encoded)
	placeholder
		if err != nil {
			return nil, err
	placeholder
placeholder

	toolChoice := gjson.GetBytes(body, "tool_choice")
	if !toolChoice.Exists() {
		return body, nil
placeholder
	if shouldDropGrokToolChoice(toolChoice, filteredTools) {
		body, err = sjson.DeleteBytes(body, "tool_choice")
		if err != nil {
			return nil, err
	placeholder
placeholder
	return body, nil
placeholder

func shouldDropGrokToolChoice(toolChoice gjson.Result, tools []json.RawMessage) bool {
	if len(tools) == 0 {
		return true
placeholder
	if !toolChoice.IsObject() {
		return false
placeholder
	choiceType := strings.TrimSpace(toolChoice.Get("type").String())
	if choiceType == "" {
		return false
placeholder
	if _, ok := grokResponsesSupportedToolTypes[choiceType]; !ok {
		return true
placeholder
	if choiceType == "function" {
		choiceName := strings.TrimSpace(toolChoice.Get("name").String())
		if choiceName == "" {
			choiceName = strings.TrimSpace(toolChoice.Get("function.name").String())
	placeholder
		if choiceName == "" {
			return false
	placeholder
		for _, tool := range tools {
			var item struct {
				Type     string `json:"type"`
				Name     string `json:"name"`
				Function struct {
					Name string `json:"name"`
			placeholder `json:"function"`
		placeholder
			if err := json.Unmarshal(tool, &item); err != nil {
				continue
		placeholder
			name := strings.TrimSpace(item.Name)
			if name == "" {
				name = strings.TrimSpace(item.Function.Name)
		placeholder
			if strings.TrimSpace(item.Type) == "function" && name == choiceName {
				return false
		placeholder
	placeholder
		return true
placeholder
	return false
placeholder

func (s *OpenAIGatewayService) bridgeGrokComposerImageInputs(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	body []byte,
	token string,
) ([]byte, OpenAIUsage, bool, error) {
	if !shouldBridgeGrokComposerImageInputs(body) {
		return body, OpenAIUsage{placeholder, false, nil
placeholder

	var reqBody map[string]any
	if err := json.Unmarshal(body, &reqBody); err != nil {
		return body, OpenAIUsage{placeholder, false, fmt.Errorf("parse grok composer image bridge request: %w", err)
placeholder

	imageURLs := collectGrokComposerImageURLs(reqBody)
	if len(imageURLs) == 0 {
		return body, OpenAIUsage{placeholder, false, nil
placeholder

	descriptions := make([]string, 0, len(imageURLs))
	var bridgeUsage OpenAIUsage
	for index, imageURL := range imageURLs {
		description, usage, err := s.describeGrokComposerImage(ctx, c, account, token, imageURL, index+1)
		if err != nil {
			return body, bridgeUsage, false, err
	placeholder
		descriptions = append(descriptions, description)
		addOpenAIUsage(&bridgeUsage, usage)
placeholder

	if !rewriteGrokComposerImagesAsText(reqBody, descriptions) {
		return body, bridgeUsage, false, nil
placeholder
	bridgedBody, err := marshalOpenAIUpstreamJSON(reqBody)
	if err != nil {
		return body, bridgeUsage, false, fmt.Errorf("serialize grok composer image bridge request: %w", err)
placeholder
	return bridgedBody, bridgeUsage, true, nil
placeholder

func shouldBridgeGrokComposerImageInputs(body []byte) bool {
	if len(body) == 0 || !isGrokComposerModel(gjson.GetBytes(body, "model").String()) {
		return false
placeholder
	messages := gjson.GetBytes(body, "messages")
	if !messages.Exists() {
		return false
placeholder
	return openAIJSONValueMayContainImageInput(messages)
placeholder

func isGrokComposerModel(model string) bool {
	model = strings.TrimSpace(strings.ToLower(model))
	if model == "" {
		return false
placeholder
	if strings.Contains(model, "/") {
		parts := strings.Split(model, "/")
		model = strings.TrimSpace(parts[len(parts)-1])
placeholder
	return strings.Contains(model, "composer")
placeholder

func collectGrokComposerImageURLs(reqBody map[string]any) []string {
	messages, ok := reqBody["messages"].([]any)
	if !ok {
		return nil
placeholder

	var imageURLs []string
	for _, msg := range messages {
		msgMap, ok := msg.(map[string]any)
		if !ok {
			continue
	placeholder
		parts, ok := msgMap["content"].([]any)
		if !ok {
			continue
	placeholder
		for _, part := range parts {
			if imageURL := grokComposerImageURLFromPart(part); imageURL != "" {
				imageURLs = append(imageURLs, imageURL)
		placeholder
	placeholder
placeholder
	return imageURLs
placeholder

func grokComposerImageURLFromPart(part any) string {
	partMap, ok := part.(map[string]any)
	if !ok {
		return ""
placeholder
	if strings.TrimSpace(strings.ToLower(fmt.Sprint(partMap["type"]))) != "image_url" {
		return ""
placeholder
	switch imageURL := partMap["image_url"].(type) {
	case string:
		return normalizeGrokComposerImageURL(imageURL)
	case map[string]any:
		raw, _ := imageURL["url"].(string)
		return normalizeGrokComposerImageURL(raw)
	default:
		return ""
placeholder
placeholder

func normalizeGrokComposerImageURL(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || isEmptyBase64DataURI(trimmed) {
		return ""
placeholder
	return trimmed
placeholder

func (s *OpenAIGatewayService) describeGrokComposerImage(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	token string,
	imageURL string,
	index int,
) (string, OpenAIUsage, error) {
	body, err := buildGrokComposerImageDescriptionBody(imageURL, index)
	if err != nil {
		return "", OpenAIUsage{placeholder, err
placeholder

	upstreamCtx, releaseUpstreamCtx := detachUpstreamContext(ctx)
	// Image-description probes are auxiliary requests, not conversation turns.
	// Do not bind them to the caller's Grok prompt-cache identity.
	upstreamReq, err := buildGrokResponsesRequest(upstreamCtx, c, account, body, token, "", s.cfg, s.settingService)
	releaseUpstreamCtx()
	if err != nil {
		return "", OpenAIUsage{placeholder, fmt.Errorf("build grok composer image bridge request: %w", err)
placeholder

	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
placeholder

	resp, err := s.httpUpstream.Do(upstreamReq, proxyURL, account.ID, account.Concurrency)
	if err != nil {
		return "", OpenAIUsage{placeholder, s.handleOpenAIUpstreamTransportError(ctx, c, account, err, false)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	if resp.StatusCode >= 400 {
		respBody := s.readUpstreamErrorBody(resp)
		upstreamMsg := sanitizeUpstreamErrorMessage(extractUpstreamErrorMessage(respBody))
		if upstreamMsg == "" {
			upstreamMsg = fmt.Sprintf("xAI image bridge upstream returned status %d", resp.StatusCode)
	placeholder
		kind := "http_error"
		if s.shouldFailoverGrokUpstreamError(resp.StatusCode, respBody) {
			kind = "failover"
	placeholder
		appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
			Platform:           account.Platform,
			AccountID:          account.ID,
			AccountName:        account.Name,
			UpstreamStatusCode: resp.StatusCode,
			UpstreamRequestID:  firstNonEmpty(resp.Header.Get("x-request-id"), resp.Header.Get("xai-request-id")),
			Kind:               kind,
			Message:            upstreamMsg,
	placeholder)
		s.handleGrokAccountUpstreamError(withGrokTeamRateLimitModel(ctx, grokComposerImageBridgeVisionModel), account, resp.StatusCode, resp.Header, respBody)
		if s.shouldFailoverGrokUpstreamError(resp.StatusCode, respBody) {
			return "", OpenAIUsage{placeholder, &UpstreamFailoverError{
				StatusCode:             resp.StatusCode,
				ResponseBody:           respBody,
				ResponseHeaders:        resp.Header.Clone(),
				RetryableOnSameAccount: account.IsPoolMode() && account.IsPoolModeRetryableStatus(resp.StatusCode),
		placeholder
	placeholder
		return "", OpenAIUsage{placeholder, fmt.Errorf("grok composer image bridge upstream error: %s", upstreamMsg)
placeholder

	s.updateGrokUsageFromResponse(withGrokTeamRateLimitModel(ctx, grokComposerImageBridgeVisionModel), account, resp.Header, resp.StatusCode)
	respBody, err := ReadUpstreamResponseBody(resp.Body, s.cfg, c, nil)
	if err != nil {
		return "", OpenAIUsage{placeholder, fmt.Errorf("read grok composer image bridge response: %w", err)
placeholder

	var parsed apicompat.ResponsesResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", OpenAIUsage{placeholder, fmt.Errorf("parse grok composer image bridge response: %w", err)
placeholder
	description := strings.TrimSpace(grokResponsesOutputText(&parsed))
	if description == "" {
		return "", copyOpenAIUsageFromResponsesUsage(parsed.Usage), fmt.Errorf("grok composer image bridge returned empty description")
placeholder
	return description, copyOpenAIUsageFromResponsesUsage(parsed.Usage), nil
placeholder

func buildGrokComposerImageDescriptionBody(imageURL string, index int) ([]byte, error) {
	prompt := fmt.Sprintf("Describe image %d in concise, factual text for a downstream coding/composer model. Include visible text, UI elements, diagrams, errors, and spatial relationships. Do not mention that you are an image analysis bridge.", index)
	req := map[string]any{
		"model":             grokComposerImageBridgeVisionModel,
		"stream":            false,
		"store":             false,
		"max_output_tokens": grokComposerImageBridgeMaxOutputTokens,
		"input": []any{
			map[string]any{
				"type": "message",
				"role": "user",
				"content": []any{
					map[string]any{"type": "input_text", "text": promptplaceholder,
					map[string]any{"type": "input_image", "image_url": imageURLplaceholder,
			placeholder,
		placeholder,
	placeholder,
placeholder
	return marshalOpenAIUpstreamJSON(req)
placeholder

func grokResponsesOutputText(resp *apicompat.ResponsesResponse) string {
	if resp == nil {
		return ""
placeholder
	var parts []string
	for _, output := range resp.Output {
		for _, content := range output.Content {
			if content.Type == "output_text" || content.Type == "text" || content.Type == "input_text" {
				if text := strings.TrimSpace(content.Text); text != "" {
					parts = append(parts, text)
			placeholder
		placeholder
	placeholder
placeholder
	return strings.Join(parts, "\n\n")
placeholder

func rewriteGrokComposerImagesAsText(reqBody map[string]any, descriptions []string) bool {
	messages, ok := reqBody["messages"].([]any)
	if !ok {
		return false
placeholder

	imageIndex := 0
	changed := false
	for _, msg := range messages {
		msgMap, ok := msg.(map[string]any)
		if !ok {
			continue
	placeholder
		parts, ok := msgMap["content"].([]any)
		if !ok {
			continue
	placeholder
		var textParts []string
		messageChanged := false
		for _, part := range parts {
			if imageURL := grokComposerImageURLFromPart(part); imageURL != "" {
				if imageIndex < len(descriptions) {
					textParts = append(textParts, fmt.Sprintf("Image %d description: %s", imageIndex+1, strings.TrimSpace(descriptions[imageIndex])))
			placeholder
				imageIndex++
				messageChanged = true
				continue
		placeholder
			if text := grokComposerTextFromPart(part); text != "" {
				textParts = append(textParts, text)
		placeholder
	placeholder
		if messageChanged {
			msgMap["content"] = strings.Join(textParts, "\n\n")
			changed = true
	placeholder
placeholder
	return changed
placeholder

func grokComposerTextFromPart(part any) string {
	partMap, ok := part.(map[string]any)
	if !ok {
		return ""
placeholder
	partType := strings.TrimSpace(strings.ToLower(fmt.Sprint(partMap["type"])))
	switch partType {
	case "text", "input_text":
		text, _ := partMap["text"].(string)
		return strings.TrimSpace(text)
	default:
		return ""
placeholder
placeholder

func addOpenAIUsage(dst *OpenAIUsage, usage OpenAIUsage) {
	if dst == nil {
		return
placeholder
	dst.InputTokens += usage.InputTokens
	dst.ImageInputTokens += usage.ImageInputTokens
	dst.OutputTokens += usage.OutputTokens
	dst.CacheCreationInputTokens += usage.CacheCreationInputTokens
	dst.CacheReadInputTokens += usage.CacheReadInputTokens
	dst.ImageOutputTokens += usage.ImageOutputTokens
placeholder

func buildGrokResponsesRequest(ctx context.Context, c *gin.Context, account *Account, body []byte, token, cacheIdentity string, cfg *config.Config, settings ...*SettingService) (*http.Request, error) {
	targetURL, err := buildGrokResponsesURL(account, cfg, settings...)
	if err != nil {
		return nil, err
placeholder
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
placeholder
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	if account.IsGrokOAuth() {
		applyGrokCLIHeaders(req.Header)
placeholder
	applyGrokCacheHeaders(req.Header, cacheIdentity)
	if c != nil {
		if v := c.GetHeader("OpenAI-Beta"); strings.TrimSpace(v) != "" {
			req.Header.Set("OpenAI-Beta", v)
	placeholder
placeholder
	// 账号级请求头覆写最后应用，使配置值优先于上面的内置默认头；
	// 打到官方 CLI 网关时身份头仍由共享传输层最终强制。
	account.ApplyHeaderOverrides(req.Header)
	return req, nil
placeholder

// applyGrokCLIHeaders identifies subscription traffic as a supported Grok CLI
// version. The CLI gateway rejects otherwise valid OAuth requests without it.
// Identity pins come from package xai so service-layer headers match the final
// transport rewrite on cli-chat-proxy.grok.com.
func applyGrokCLIHeaders(headers http.Header) {
	if headers == nil {
		return
placeholder
	version := xai.ResolveCLIVersion()
	headers.Set("User-Agent", xai.CLIUserAgent(version))
	headers.Set("X-Grok-Client-Version", version)
	headers.Set("x-grok-client-version", version)
	headers.Set("x-grok-client-identifier", xai.CLIClientIdentifier)
	// Historical mode value expected by some unit tests / older CLI probes.
	headers.Set("X-Grok-Client-Mode", "interactive")
placeholder

func (s *OpenAIGatewayService) updateGrokUsageSnapshot(ctx context.Context, account *Account, snapshot *xai.QuotaSnapshot) {
	if s == nil || account == nil || account.ID <= 0 || snapshot == nil {
		return
placeholder
	accountID := account.ID
	now := time.Now()
	resetAt, hasActiveLimit := grokRateLimitResetAtForAccount(account, snapshot, now)
	if hasActiveLimit {
		normalizeGrokExhaustedWindowResets(snapshot, resetAt, now)
placeholder
	recovery := isSuccessfulGrokRateLimitRecovery(account, snapshot)
	critical := snapshot.StatusCode == http.StatusTooManyRequests || hasActiveLimit || recovery
	if s.codexSnapshotThrottle != nil {
		allowed := s.codexSnapshotThrottle.Allow(accountID, now)
		if !critical && !allowed {
			return
	placeholder
placeholder

	updates := map[string]any{
		grokQuotaSnapshotExtraKey: snapshot,
placeholder
	// Also derive the scheduling-threshold extras (grok_sched_*) the evaluator
	// reads in grokThresholdCandidates. Without this writer the admin-configured
	// Grok auto-pause threshold could never fire (the read side was dead config).
	for k, v := range buildGrokSchedulerExtraUpdates(snapshot) {
		updates[k] = v
placeholder
	stateCtx := ctx
	if hasActiveLimit {
		var cancel context.CancelFunc
		stateCtx, cancel = openAIAccountStateContext(ctx)
		defer cancel()
placeholder
	// Account pointers on the request path are per-request copies (Redis/DB decode),
	// not a shared in-process cache. Mutating Extra here matches token refresh /
	// rate-limit writers; do not reuse the same *Account across goroutines.
	if account.Extra == nil {
		account.Extra = map[string]any{placeholder
placeholder
	account.Extra[grokQuotaSnapshotExtraKey] = snapshot
	if s.accountRepo != nil {
		_ = s.accountRepo.UpdateExtra(stateCtx, accountID, updates)
placeholder
	// Error responses are reconciled by handleGrokAccountUpstreamError. Pool-mode
	// API keys retain the snapshot for observability but leave account health to
	// the upstream pool. Other accounts install the immediate runtime and durable
	// rate-limit state when the observed window is exhausted.
	if hasActiveLimit && !account.IsPoolMode() {
		s.rateLimitGrok(stateCtx, account, resetAt)
placeholder else if recovery {
		clearGrokRateLimitAfterRecovery(stateCtx, s.accountRepo, account)
placeholder
placeholder

func (s *OpenAIGatewayService) updateGrokUsageFromResponse(ctx context.Context, account *Account, headers http.Header, statusCode int) {
	snapshot := parseGrokQuotaSnapshot(headers, statusCode, time.Now())
	if snapshot != nil {
		stampGrokQuotaSnapshotForPlan(account, snapshot, grokRequestedModelFromCtx(ctx))
		s.updateGrokUsageSnapshot(ctx, account, snapshot)
		return
placeholder
	// Successful responses are recovery evidence even when the upstream omits
	// optional quota headers. Do not replace an informative stored snapshot with
	// an empty one; only clear the exact observed cooldown generation.
	recoverySnapshot := &xai.QuotaSnapshot{StatusCode: statusCodeplaceholder
	if isSuccessfulGrokRateLimitRecovery(account, recoverySnapshot) {
		clearGrokRateLimitAfterRecovery(ctx, s.accountRepo, account)
placeholder
placeholder

func parseGrokQuotaSnapshot(headers http.Header, statusCode int, now time.Time) *xai.QuotaSnapshot {
	snapshot := xai.ParseQuotaHeaders(headers, statusCode)
	if snapshot == nil && statusCode == http.StatusTooManyRequests {
		return &xai.QuotaSnapshot{
			StatusCode: statusCode,
			UpdatedAt:  now.UTC().Format(time.RFC3339),
	placeholder
placeholder
	return snapshot
placeholder

func normalizeGrokExhaustedWindowResets(snapshot *xai.QuotaSnapshot, resetAt, now time.Time) {
	if snapshot == nil || !resetAt.After(now) {
		return
placeholder
	for _, window := range []*xai.QuotaWindow{snapshot.Requests, snapshot.Tokensplaceholder {
		if window == nil || window.Remaining == nil || *window.Remaining > 0 {
			continue
	placeholder
		candidate := time.Time{placeholder
		if window.ResetUnix != nil && *window.ResetUnix > 0 {
			candidate = time.Unix(*window.ResetUnix, 0)
	placeholder else if parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(window.ResetAt)); err == nil {
			candidate = parsed
	placeholder
		if !candidate.After(now) {
			candidate = resetAt
	placeholder
		resetUnix := candidate.Unix()
		window.ResetUnix = &resetUnix
		window.ResetAt = candidate.UTC().Format(time.RFC3339)
placeholder
placeholder

func grokRateLimitResetAt(snapshot *xai.QuotaSnapshot, now time.Time) (time.Time, bool) {
	if snapshot == nil {
		return time.Time{placeholder, false
placeholder

	// Retry-After is xAI's explicit retry boundary. Use the observation time so
	// a persisted snapshot does not start a fresh cooldown every time it is read.
	retryAfterExpired := false
	var resetAt time.Time
	if snapshot.RetryAfterSeconds != nil && *snapshot.RetryAfterSeconds > 0 {
		observedAt := now
		if parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(snapshot.UpdatedAt)); err == nil {
			observedAt = parsed
	placeholder
		retryAfterResetAt := observedAt.Add(time.Duration(*snapshot.RetryAfterSeconds) * time.Second)
		if retryAfterResetAt.After(now) {
			resetAt = retryAfterResetAt
	placeholder else {
			retryAfterExpired = true
	placeholder
placeholder

	exhausted := false
	for _, window := range []*xai.QuotaWindow{snapshot.Requests, snapshot.Tokensplaceholder {
		if window == nil || window.Remaining == nil || *window.Remaining > 0 {
			continue
	placeholder
		exhausted = true
		candidate := time.Time{placeholder
		if window.ResetUnix != nil && *window.ResetUnix > 0 {
			candidate = time.Unix(*window.ResetUnix, 0)
	placeholder else if parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(window.ResetAt)); err == nil {
			candidate = parsed
	placeholder
		if candidate.After(now) && candidate.After(resetAt) {
			resetAt = candidate
	placeholder
placeholder
	if !resetAt.IsZero() {
		return resetAt, true
placeholder
	// An observed Retry-After is an absolute boundary once combined with the
	// snapshot timestamp. Do not turn an expired persisted snapshot into a new
	// rolling fallback cooldown, but still allow a later explicit window reset.
	if retryAfterExpired {
		return time.Time{placeholder, false
placeholder
	if exhausted || snapshot.StatusCode == http.StatusTooManyRequests {
		return now.Add(grokRateLimitFallbackCooldown), true
placeholder
	return time.Time{placeholder, false
placeholder

func grokRateLimitResetAtForAccount(account *Account, snapshot *xai.QuotaSnapshot, now time.Time) (time.Time, bool) {
	resetAt, limited := grokRateLimitResetAt(snapshot, now)
	if !limited || !isGrokOAuthAccount(account) || snapshot == nil || snapshot.StatusCode != http.StatusTooManyRequests {
		return resetAt, limited
placeholder
	if account.RateLimitedAt == nil || account.RateLimitResetAt == nil {
		return resetAt, true
placeholder
	previousResetAt := *account.RateLimitResetAt
	if previousResetAt.After(now) || now.Sub(previousResetAt) > grokRateLimitBackoffQuietPeriod {
		return resetAt, true
placeholder
	previousCooldown := previousResetAt.Sub(*account.RateLimitedAt)
	if previousCooldown <= 0 {
		return resetAt, true
placeholder

	adaptiveCooldown := grokRateLimitRepeatCooldown
	switch {
	case previousCooldown >= grokRateLimitSustainedCooldown:
		adaptiveCooldown = grokRateLimitMaxAdaptiveCooldown
	case previousCooldown >= grokRateLimitRepeatCooldown:
		adaptiveCooldown = grokRateLimitSustainedCooldown
placeholder
	adaptiveResetAt := now.Add(adaptiveCooldown)
	if adaptiveResetAt.After(resetAt) {
		resetAt = adaptiveResetAt
placeholder
	return resetAt, true
placeholder

func normalizeGrokRateLimitResetAt(account *Account, resetAt, now time.Time) time.Time {
	if !resetAt.After(now) {
		resetAt = now.Add(grokRateLimitFallbackCooldown)
placeholder
	if account != nil && account.RateLimitResetAt != nil && account.RateLimitResetAt.After(resetAt) {
		resetAt = *account.RateLimitResetAt
placeholder
	return resetAt
placeholder

type grokRateLimitExtendingRepository interface {
	SetRateLimitedIfLater(ctx context.Context, id int64, resetAt time.Time) error
placeholder

type grokRateLimitRecoveryRepository interface {
	ClearRateLimitIfObserved(ctx context.Context, id int64, observedLimitedAt, observedResetAt time.Time) (bool, error)
placeholder

func isSuccessfulGrokRateLimitRecovery(account *Account, snapshot *xai.QuotaSnapshot) bool {
	return isGrokOAuthAccount(account) &&
		account.RateLimitedAt != nil &&
		account.RateLimitResetAt != nil &&
		snapshot != nil &&
		snapshot.StatusCode >= http.StatusOK &&
		snapshot.StatusCode < http.StatusMultipleChoices
placeholder

func clearGrokRateLimitAfterRecovery(ctx context.Context, repo AccountRepository, account *Account) {
	if repo == nil || account == nil || account.RateLimitedAt == nil || account.RateLimitResetAt == nil || ctx.Err() != nil {
		return
placeholder
	recoveryRepo, ok := repo.(grokRateLimitRecoveryRepository)
	if !ok {
		return
placeholder
	_, err := recoveryRepo.ClearRateLimitIfObserved(ctx, account.ID, *account.RateLimitedAt, *account.RateLimitResetAt)
	if err != nil {
		slog.Warn("grok_rate_limit_recovery_clear_failed", "account_id", account.ID, "error", err)
placeholder
placeholder

func persistGrokRateLimit(ctx context.Context, repo AccountRepository, account *Account, resetAt time.Time) {
	if repo == nil || account == nil || account.ID <= 0 {
		return
placeholder
	resetAt = normalizeGrokRateLimitResetAt(account, resetAt, time.Now())
	stateCtx, cancel := openAIAccountStateContext(ctx)
	defer cancel()
	var err error
	if extendingRepo, ok := repo.(grokRateLimitExtendingRepository); ok {
		err = extendingRepo.SetRateLimitedIfLater(stateCtx, account.ID, resetAt)
placeholder else {
		err = repo.SetRateLimited(stateCtx, account.ID, resetAt)
placeholder
	if err != nil {
		slog.Warn("persist_grok_rate_limit_failed", "account_id", account.ID, "reset_at", resetAt.UTC(), "error", err)
placeholder
placeholder

func (s *OpenAIGatewayService) rateLimitGrok(ctx context.Context, account *Account, resetAt time.Time) {
	if s == nil || account == nil {
		return
placeholder
	now := time.Now()
	resetAt = normalizeGrokRateLimitResetAt(account, resetAt, now)

	runtimeUntil := resetAt
	if account.TempUnschedulableUntil != nil && account.TempUnschedulableUntil.After(runtimeUntil) {
		runtimeUntil = *account.TempUnschedulableUntil
placeholder
	s.BlockAccountScheduling(account, runtimeUntil, "429")
	persistGrokRateLimit(ctx, s.accountRepo, account, resetAt)

	// Propagate a short team+model cool so sibling OAuth accounts on the same
	// xAI team skip the hot model without waiting for each to hit 429 alone.
	// Model is taken from the latest request context when available; empty is a
	// no-op inside markGrokTeamModelRateLimit.
	if model, _ := ctx.Value(grokTeamRateLimitModelContextKey{placeholder).(string); model != "" {
		markGrokTeamModelRateLimit(account, model, resolveGrokTeamRateLimitUntil(resetAt, now))
placeholder
placeholder

// buildGrokSchedulerExtraUpdates derives the grok_sched_* scheduling snapshot
// (utilization percent + reset time) consumed by EvaluateAccountSchedulingThreshold.
// Utilization is the most-constrained of the requests/tokens windows.
func buildGrokSchedulerExtraUpdates(snapshot *xai.QuotaSnapshot) map[string]any {
	if snapshot == nil {
		return nil
placeholder
	util, reset, ok := grokSnapshotUtilization(snapshot)
	if !ok {
		return nil
placeholder
	updates := map[string]any{
		"grok_sched_utilization":      util,
		"grok_sched_usage_updated_at": time.Now().UTC().Format(time.RFC3339),
placeholder
	if reset != nil {
		// 防御：调度阈值暂停时长由 grok_sched_reset_at 决定。若上游返回脏的
		// reset 头（例如把相对毫秒 "6000" 误当相对秒解析出 ~33h 的未来时刻），
		// 不设上限会把耗尽账号长时间锁死。xAI 配额窗口不会超过一天，因此对
		// 未来时刻做 grokMaxSchedulingResetHorizon 钳制；过去/无效值直接不写。
		now := time.Now()
		if reset.After(now) {
			capped := *reset
			if horizon := now.Add(grokMaxSchedulingResetHorizon); capped.After(horizon) {
				capped = horizon
		placeholder
			updates["grok_sched_reset_at"] = capped.UTC().Format(time.RFC3339)
	placeholder
placeholder
	return updates
placeholder

// grokSnapshotUtilization returns the highest window utilization (0-100) across
// the requests/tokens quota windows and the reset time of that window.
func grokSnapshotUtilization(snapshot *xai.QuotaSnapshot) (float64, *time.Time, bool) {
	if snapshot == nil {
		return 0, nil, false
placeholder
	best := -1.0
	var bestReset *time.Time
	consider := func(window *xai.QuotaWindow) {
		if window == nil || window.Limit == nil || *window.Limit <= 0 || window.Remaining == nil {
			return
	placeholder
		remaining := *window.Remaining
		if remaining < 0 {
			remaining = 0
	placeholder
		util := (1 - float64(remaining)/float64(*window.Limit)) * 100
		if util < 0 {
			util = 0
	placeholder
		if util > 100 {
			util = 100
	placeholder
		if util > best {
			best = util
			if window.ResetUnix != nil {
				t := time.Unix(*window.ResetUnix, 0).UTC()
				bestReset = &t
		placeholder else {
				bestReset = nil
		placeholder
	placeholder
placeholder
	consider(snapshot.Requests)
	consider(snapshot.Tokens)
	if best < 0 {
		return 0, nil, false
placeholder
	return best, bestReset, true
placeholder

// grokMaxSchedulingResetHorizon bounds how far into the future a Grok
// scheduling-threshold pause (grok_sched_reset_at) may be set, so a malformed
// upstream reset header can't park an over-threshold account for days. xAI quota
// windows do not exceed ~a day.
const grokMaxSchedulingResetHorizon = 25 * time.Hour

// grokTeamRateLimitModelContextKey carries the upstream model for team cools.
type grokTeamRateLimitModelContextKey struct{placeholder

// withGrokTeamRateLimitModel attaches the upstream model name for rate-limit
// side effects (team+model cool). Safe when model is empty.
func withGrokTeamRateLimitModel(ctx context.Context, model string) context.Context {
	model = strings.TrimSpace(model)
	if model == "" || ctx == nil {
		return ctx
placeholder
	return context.WithValue(ctx, grokTeamRateLimitModelContextKey{placeholder, model)
placeholder

func grokRequestedModelFromCtx(ctx context.Context) string {
	if ctx == nil {
		return ""
placeholder
	model, _ := ctx.Value(grokTeamRateLimitModelContextKey{placeholder).(string)
	return strings.TrimSpace(model)
placeholder

func isGrokHeavyTransientModel(requestedModel string) bool {
	model := strings.ToLower(strings.TrimSpace(xai.ResolveGrokTextResponsesModelID(requestedModel)))
	return strings.Contains(model, "multi-agent")
placeholder

func persistGrokTransientModelCooldown(account *Account, decision GrokUpstreamFailureDecision) bool {
	if account == nil {
		return false
placeholder
	model := strings.TrimSpace(decision.Model)
	if model == "" || !isGrokHeavyTransientModel(model) {
		return false
placeholder
	cooldown := decision.Cooldown
	if cooldown <= 0 {
		cooldown = 3 * time.Minute
placeholder
	markGrokModelTransientBlock(account.ID, model, time.Now().Add(cooldown))
	return true
placeholder

func (s *OpenAIGatewayService) handleGrokAccountUpstreamError(ctx context.Context, account *Account, statusCode int, headers http.Header, responseBody []byte) {
	if s == nil || account == nil {
		return
placeholder
	if isGrokContentPolicyRejection(statusCode, responseBody) {
		return
placeholder
	now := time.Now()
	snapshot := parseGrokQuotaSnapshot(headers, statusCode, now)
	stampGrokQuotaSnapshotForPlan(account, snapshot, grokRequestedModelFromCtx(ctx))
	s.updateGrokUsageSnapshot(ctx, account, snapshot)

	// Body-first free-usage / empty / billing / capacity must run before the
	// status switch so non-429 free-usage bodies still cool the account.
	// Pool-mode still skips durable mutation unless an explicit temp rule matches.
	decision := classifyGrokUpstreamFailure(statusCode, responseBody, grokRequestedModelFromCtx(ctx))
	if decision.ShouldCooldown && decision.Class != GrokFailureNone && decision.Class != GrokFailureRateLimit {
		if account.IsPoolMode() {
			// Allow configured temp rules (403) below; skip default body cools.
	placeholder else {
			// A free-tier exhaustion message describes a rolling usage window. Use
			// an upstream absolute reset (or Retry-After) when available; otherwise
			// apply only a short probe cooldown. Never start a fabricated 24h window
			// at the instant this error was received.
			if decision.Class == GrokFailureFreeUsage {
				if resetAt, limited := grokRateLimitResetAtForAccount(account, parseGrokQuotaSnapshot(headers, statusCode, now), now); limited && resetAt.After(now) {
					if decision.Model != "" && isGrokModelSpecificFreeUsage(strings.ToLower(decision.Reason), decision.Model) {
						markGrokModelQuotaBlock(account.ID, decision.Model, resetAt)
						return
				placeholder
					s.rateLimitGrok(ctx, account, resetAt)
					return
			placeholder
		placeholder
			if s.applyGrokUpstreamFailureDecision(ctx, account, decision) {
				return
		placeholder
	placeholder
placeholder

	if statusCode == http.StatusForbidden && s.applyGrokForbiddenPolicy(ctx, account, responseBody) {
		return
placeholder
	if account.IsPoolMode() {
		slog.Info("grok_pool_mode_error_state_skipped", "account_id", account.ID, "status_code", statusCode)
		return
placeholder
	switch statusCode {
	case http.StatusUnauthorized:
		s.tempUnscheduleGrok(ctx, account, 10*time.Minute, "grok credentials unauthorized")
	case http.StatusPaymentRequired:
		// 402 without a body-classified billing decision: keep the legacy 30m cool.
		s.tempUnscheduleGrok(ctx, account, 30*time.Minute, "grok payment required")
	case http.StatusForbidden:
		// Spending-limit already handled by body classifier when phrasing matches.
		if isGrokSpendingLimitError(responseBody) {
			s.rateLimitGrok(ctx, account, grokSpendingLimitResetAt(account, time.Now()))
			return
	placeholder
		s.tempUnscheduleGrok(ctx, account, 30*time.Minute, "grok access or entitlement denied")
	case http.StatusTooManyRequests:
		// updateGrokUsageSnapshot installs rate-limit state for non-pool accounts.
		// Free-usage 429 was already cooled above via body classification.
	default:
		if statusCode >= 500 {
			s.tempUnscheduleGrok(ctx, account, 2*time.Minute, "grok upstream temporary error")
	placeholder
placeholder
placeholder

// isGrokSpendingLimitError detects xAI billing exhaustion bodies (often 403, sometimes 402).
func isGrokSpendingLimitError(responseBody []byte) bool {
	if len(responseBody) == 0 {
		return false
placeholder
	code := strings.ToLower(strings.TrimSpace(firstNonEmpty(
		gjson.GetBytes(responseBody, "code").String(),
		gjson.GetBytes(responseBody, "error.code").String(),
	)))
	if code == "personal-team-blocked:spending-limit" {
		return true
placeholder
	message := strings.ToLower(strings.TrimSpace(firstNonEmpty(
		gjson.GetBytes(responseBody, "error").String(),
		gjson.GetBytes(responseBody, "error.message").String(),
		gjson.GetBytes(responseBody, "message").String(),
	)))
	return strings.Contains(message, "spending limit") ||
		strings.Contains(message, "run out of credits")
placeholder

func (s *OpenAIGatewayService) tempUnscheduleGrok(ctx context.Context, account *Account, cooldown time.Duration, reason string) {
	if s == nil || account == nil {
		return
placeholder
	until := time.Now().Add(cooldown)
	if account.TempUnschedulableUntil != nil && account.TempUnschedulableUntil.After(until) {
		until = *account.TempUnschedulableUntil
placeholder
	s.BlockAccountScheduling(account, until, reason)
	if s.accountRepo != nil {
		stateCtx, cancel := openAIAccountStateContext(ctx)
		defer cancel()
		_ = s.accountRepo.SetTempUnschedulable(stateCtx, account.ID, until, reason)
placeholder
placeholder
