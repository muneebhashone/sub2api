package service

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	mathrand "math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/model"
	"github.com/Wei-Shaw/sub2api/internal/pkg/geminicli"
	"github.com/Wei-Shaw/sub2api/internal/service/ports"

	"github.com/gin-gonic/gin"
)

const geminiStickySessionTTL = time.Hour

const (
	geminiMaxRetries     = 5
	geminiRetryBaseDelay = 1 * time.Second
	geminiRetryMaxDelay  = 16 * time.Second
)

type GeminiMessagesCompatService struct {
	accountRepo      ports.AccountRepository
	cache            ports.GatewayCache
	tokenProvider    *GeminiTokenProvider
	rateLimitService *RateLimitService
	httpUpstream     ports.HTTPUpstream
placeholder

func NewGeminiMessagesCompatService(
	accountRepo ports.AccountRepository,
	cache ports.GatewayCache,
	tokenProvider *GeminiTokenProvider,
	rateLimitService *RateLimitService,
	httpUpstream ports.HTTPUpstream,
) *GeminiMessagesCompatService {
	return &GeminiMessagesCompatService{
		accountRepo:      accountRepo,
		cache:            cache,
		tokenProvider:    tokenProvider,
		rateLimitService: rateLimitService,
		httpUpstream:     httpUpstream,
placeholder
placeholder

func (s *GeminiMessagesCompatService) SelectAccountForModel(ctx context.Context, groupID *int64, sessionHash string, requestedModel string) (*model.Account, error) {
	cacheKey := "gemini:" + sessionHash
	if sessionHash != "" {
		accountID, err := s.cache.GetSessionAccountID(ctx, cacheKey)
		if err == nil && accountID > 0 {
			account, err := s.accountRepo.GetByID(ctx, accountID)
			if err == nil && account.IsSchedulable() && account.Platform == model.PlatformGemini && (requestedModel == "" || account.IsModelSupported(requestedModel)) {
				_ = s.cache.RefreshSessionTTL(ctx, cacheKey, geminiStickySessionTTL)
				return account, nil
		placeholder
	placeholder
placeholder

	var accounts []model.Account
	var err error
	if groupID != nil {
		accounts, err = s.accountRepo.ListSchedulableByGroupIDAndPlatform(ctx, *groupID, model.PlatformGemini)
placeholder else {
		accounts, err = s.accountRepo.ListSchedulableByPlatform(ctx, model.PlatformGemini)
placeholder
	if err != nil {
		return nil, fmt.Errorf("query accounts failed: %w", err)
placeholder

	var selected *model.Account
	for i := range accounts {
		acc := &accounts[i]
		if requestedModel != "" && !acc.IsModelSupported(requestedModel) {
			continue
	placeholder
		if selected == nil {
			selected = acc
			continue
	placeholder
		if acc.Priority < selected.Priority {
			selected = acc
	placeholder else if acc.Priority == selected.Priority {
			if acc.LastUsedAt == nil || (selected.LastUsedAt != nil && acc.LastUsedAt.Before(*selected.LastUsedAt)) {
				selected = acc
		placeholder
	placeholder
placeholder

	if selected == nil {
		if requestedModel != "" {
			return nil, fmt.Errorf("no available Gemini accounts supporting model: %s", requestedModel)
	placeholder
		return nil, errors.New("no available Gemini accounts")
placeholder

	if sessionHash != "" {
		_ = s.cache.SetSessionAccountID(ctx, cacheKey, selected.ID, geminiStickySessionTTL)
placeholder

	return selected, nil
placeholder

func (s *GeminiMessagesCompatService) Forward(ctx context.Context, c *gin.Context, account *model.Account, body []byte) (*ForwardResult, error) {
	startTime := time.Now()

	var req struct {
		Model  string `json:"model"`
		Stream bool   `json:"stream"`
placeholder
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("parse request: %w", err)
placeholder
	if strings.TrimSpace(req.Model) == "" {
		return nil, fmt.Errorf("missing model")
placeholder

	originalModel := req.Model
	mappedModel := req.Model
	if account.Type == model.AccountTypeApiKey {
		mappedModel = account.GetMappedModel(req.Model)
placeholder

	geminiReq, err := convertClaudeMessagesToGeminiGenerateContent(body)
	if err != nil {
		return nil, s.writeClaudeError(c, http.StatusBadRequest, "invalid_request_error", err.Error())
placeholder

	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
placeholder

	var requestIDHeader string
	var buildReq func(ctx context.Context) (*http.Request, string, error)

	switch account.Type {
	case model.AccountTypeApiKey:
		buildReq = func(ctx context.Context) (*http.Request, string, error) {
			apiKey := account.GetCredential("api_key")
			if strings.TrimSpace(apiKey) == "" {
				return nil, "", errors.New("Gemini api_key not configured")
		placeholder

			baseURL := strings.TrimRight(account.GetCredential("base_url"), "/")
			if baseURL == "" {
				baseURL = geminicli.AIStudioBaseURL
		placeholder

			action := "generateContent"
			if req.Stream {
				action = "streamGenerateContent"
		placeholder
			fullURL := fmt.Sprintf("%s/v1beta/models/%s:%s", strings.TrimRight(baseURL, "/"), mappedModel, action)
			if req.Stream {
				fullURL += "?alt=sse"
		placeholder

			upstreamReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(geminiReq))
			if err != nil {
				return nil, "", err
		placeholder
			upstreamReq.Header.Set("Content-Type", "application/json")
			upstreamReq.Header.Set("x-goog-api-key", apiKey)
			return upstreamReq, "x-request-id", nil
	placeholder
		requestIDHeader = "x-request-id"

	case model.AccountTypeOAuth:
		buildReq = func(ctx context.Context) (*http.Request, string, error) {
			if s.tokenProvider == nil {
				return nil, "", errors.New("Gemini token provider not configured")
		placeholder
			accessToken, err := s.tokenProvider.GetAccessToken(ctx, account)
			if err != nil {
				return nil, "", err
		placeholder

			projectID := strings.TrimSpace(account.GetCredential("project_id"))
			if projectID == "" {
				return nil, "", errors.New("missing project_id in account credentials")
		placeholder

			action := "generateContent"
			if req.Stream {
				action = "streamGenerateContent"
		placeholder
			fullURL := fmt.Sprintf("%s/v1internal:%s", geminicli.GeminiCliBaseURL, action)
			if req.Stream {
				fullURL += "?alt=sse"
		placeholder

			wrapped := map[string]any{
				"model":   mappedModel,
				"project": projectID,
		placeholder
			var inner any
			if err := json.Unmarshal(geminiReq, &inner); err != nil {
				return nil, "", fmt.Errorf("failed to parse gemini request: %w", err)
		placeholder
			wrapped["request"] = inner
			wrappedBytes, _ := json.Marshal(wrapped)

			upstreamReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(wrappedBytes))
			if err != nil {
				return nil, "", err
		placeholder
			upstreamReq.Header.Set("Content-Type", "application/json")
			upstreamReq.Header.Set("Authorization", "Bearer "+accessToken)
			upstreamReq.Header.Set("User-Agent", geminicli.GeminiCLIUserAgent)
			return upstreamReq, "x-request-id", nil
	placeholder
		requestIDHeader = "x-request-id"

	default:
		return nil, fmt.Errorf("unsupported account type: %s", account.Type)
placeholder

	if buildReq == nil {
		return nil, s.writeClaudeError(c, http.StatusBadGateway, "upstream_error", "Gemini upstream not configured")
placeholder

	var resp *http.Response
	for attempt := 1; attempt <= geminiMaxRetries; attempt++ {
		upstreamReq, idHeader, err := buildReq(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil, err
		placeholder
			// Local build error: don't retry.
			if strings.Contains(err.Error(), "missing project_id") {
				return nil, s.writeClaudeError(c, http.StatusBadRequest, "invalid_request_error", err.Error())
		placeholder
			return nil, s.writeClaudeError(c, http.StatusBadGateway, "upstream_error", err.Error())
	placeholder
		requestIDHeader = idHeader

		resp, err = s.httpUpstream.Do(upstreamReq, proxyURL)
		if err != nil {
			if attempt < geminiMaxRetries {
				log.Printf("Gemini account %d: upstream request failed, retry %d/%d: %v", account.ID, attempt, geminiMaxRetries, err)
				sleepGeminiBackoff(attempt)
				continue
		placeholder
			return nil, s.writeClaudeError(c, http.StatusBadGateway, "upstream_error", "Upstream request failed after retries")
	placeholder

		if resp.StatusCode >= 400 && s.shouldRetryGeminiUpstreamError(account, resp.StatusCode) {
			respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
			_ = resp.Body.Close()
			if resp.StatusCode == 429 {
				// Mark as rate-limited early so concurrent requests avoid this account.
				s.handleGeminiUpstreamError(ctx, account, resp.StatusCode, resp.Header, respBody)
		placeholder
			if attempt < geminiMaxRetries {
				log.Printf("Gemini account %d: upstream status %d, retry %d/%d", account.ID, resp.StatusCode, attempt, geminiMaxRetries)
				sleepGeminiBackoff(attempt)
				continue
		placeholder
			return nil, s.writeClaudeError(c, http.StatusBadGateway, "upstream_error", "Upstream request failed after retries")
	placeholder

		break
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
		s.handleGeminiUpstreamError(ctx, account, resp.StatusCode, resp.Header, respBody)
		return nil, s.writeGeminiMappedError(c, resp.StatusCode, respBody)
placeholder

	requestID := resp.Header.Get(requestIDHeader)
	if requestID == "" {
		requestID = resp.Header.Get("x-goog-request-id")
placeholder
	if requestID != "" {
		c.Header("x-request-id", requestID)
placeholder

	var usage *ClaudeUsage
	var firstTokenMs *int
	if req.Stream {
		streamRes, err := s.handleStreamingResponse(c, resp, startTime, originalModel)
		if err != nil {
			return nil, err
	placeholder
		usage = streamRes.usage
		firstTokenMs = streamRes.firstTokenMs
placeholder else {
		usage, err = s.handleNonStreamingResponse(c, resp, originalModel)
		if err != nil {
			return nil, err
	placeholder
placeholder

	return &ForwardResult{
		RequestID:    requestID,
		Usage:        *usage,
		Model:        originalModel,
		Stream:       req.Stream,
		Duration:     time.Since(startTime),
		FirstTokenMs: firstTokenMs,
placeholder, nil
placeholder

func (s *GeminiMessagesCompatService) shouldRetryGeminiUpstreamError(account *model.Account, statusCode int) bool {
	switch statusCode {
	case 429, 500, 502, 503, 504, 529:
		return true
	case 403:
		// GeminiCli OAuth occasionally returns 403 transiently (activation/quota propagation); allow retry.
		return account != nil && account.Type == model.AccountTypeOAuth
	default:
		return false
placeholder
placeholder

func sleepGeminiBackoff(attempt int) {
	delay := geminiRetryBaseDelay * time.Duration(1<<uint(attempt-1))
	if delay > geminiRetryMaxDelay {
		delay = geminiRetryMaxDelay
placeholder

	// +/- 20% jitter
	r := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
	jitter := time.Duration(float64(delay) * 0.2 * (r.Float64()*2 - 1))
	sleepFor := delay + jitter
	if sleepFor < 0 {
		sleepFor = 0
placeholder
	time.Sleep(sleepFor)
placeholder

func (s *GeminiMessagesCompatService) writeGeminiMappedError(c *gin.Context, upstreamStatus int, body []byte) error {
	var statusCode int
	var errType, errMsg string

	if mapped := mapGeminiErrorBodyToClaudeError(body); mapped != nil {
		errType = mapped.Type
		if mapped.Message != "" {
			errMsg = mapped.Message
	placeholder
		if mapped.StatusCode > 0 {
			statusCode = mapped.StatusCode
	placeholder
placeholder

	switch upstreamStatus {
	case 400:
		if statusCode == 0 {
			statusCode = http.StatusBadRequest
	placeholder
		if errType == "" {
			errType = "invalid_request_error"
	placeholder
		if errMsg == "" {
			errMsg = "Invalid request"
	placeholder
	case 401:
		if statusCode == 0 {
			statusCode = http.StatusBadGateway
	placeholder
		if errType == "" {
			errType = "authentication_error"
	placeholder
		if errMsg == "" {
			errMsg = "Upstream authentication failed, please contact administrator"
	placeholder
	case 403:
		if statusCode == 0 {
			statusCode = http.StatusBadGateway
	placeholder
		if errType == "" {
			errType = "permission_error"
	placeholder
		if errMsg == "" {
			errMsg = "Upstream access forbidden, please contact administrator"
	placeholder
	case 404:
		if statusCode == 0 {
			statusCode = http.StatusNotFound
	placeholder
		if errType == "" {
			errType = "not_found_error"
	placeholder
		if errMsg == "" {
			errMsg = "Resource not found"
	placeholder
	case 429:
		if statusCode == 0 {
			statusCode = http.StatusTooManyRequests
	placeholder
		if errType == "" {
			errType = "rate_limit_error"
	placeholder
		if errMsg == "" {
			errMsg = "Upstream rate limit exceeded, please retry later"
	placeholder
	case 529:
		if statusCode == 0 {
			statusCode = http.StatusServiceUnavailable
	placeholder
		if errType == "" {
			errType = "overloaded_error"
	placeholder
		if errMsg == "" {
			errMsg = "Upstream service overloaded, please retry later"
	placeholder
	case 500, 502, 503, 504:
		if statusCode == 0 {
			statusCode = http.StatusBadGateway
	placeholder
		if errType == "" {
			switch upstreamStatus {
			case 504:
				errType = "timeout_error"
			case 503:
				errType = "overloaded_error"
			default:
				errType = "api_error"
		placeholder
	placeholder
		if errMsg == "" {
			errMsg = "Upstream service temporarily unavailable"
	placeholder
	default:
		if statusCode == 0 {
			statusCode = http.StatusBadGateway
	placeholder
		if errType == "" {
			errType = "upstream_error"
	placeholder
		if errMsg == "" {
			errMsg = "Upstream request failed"
	placeholder
placeholder

	c.JSON(statusCode, gin.H{
		"type":  "error",
		"error": gin.H{"type": errType, "message": errMsgplaceholder,
placeholder)
	return fmt.Errorf("upstream error: %d", upstreamStatus)
placeholder

type claudeErrorMapping struct {
	Type       string
	Message    string
	StatusCode int
placeholder

func mapGeminiErrorBodyToClaudeError(body []byte) *claudeErrorMapping {
	if len(body) == 0 {
		return nil
placeholder

	var parsed struct {
		Error struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Status  string `json:"status"`
	placeholder `json:"error"`
placeholder
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil
placeholder
	if strings.TrimSpace(parsed.Error.Status) == "" && parsed.Error.Code == 0 && strings.TrimSpace(parsed.Error.Message) == "" {
		return nil
placeholder

	mapped := &claudeErrorMapping{
		Type:    mapGeminiStatusToClaudeErrorType(parsed.Error.Status),
		Message: "",
placeholder
	if mapped.Type == "" {
		mapped.Type = "upstream_error"
placeholder

	switch strings.ToUpper(strings.TrimSpace(parsed.Error.Status)) {
	case "INVALID_ARGUMENT":
		mapped.StatusCode = http.StatusBadRequest
	case "NOT_FOUND":
		mapped.StatusCode = http.StatusNotFound
	case "RESOURCE_EXHAUSTED":
		mapped.StatusCode = http.StatusTooManyRequests
	default:
		// Keep StatusCode unset and let HTTP status mapping decide.
placeholder

	// Keep messages generic by default; upstream error message can be long or include sensitive fragments.
	return mapped
placeholder

func mapGeminiStatusToClaudeErrorType(status string) string {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "INVALID_ARGUMENT":
		return "invalid_request_error"
	case "PERMISSION_DENIED":
		return "permission_error"
	case "NOT_FOUND":
		return "not_found_error"
	case "RESOURCE_EXHAUSTED":
		return "rate_limit_error"
	case "UNAUTHENTICATED":
		return "authentication_error"
	case "UNAVAILABLE":
		return "overloaded_error"
	case "INTERNAL":
		return "api_error"
	case "DEADLINE_EXCEEDED":
		return "timeout_error"
	default:
		return ""
placeholder
placeholder

type geminiStreamResult struct {
	usage        *ClaudeUsage
	firstTokenMs *int
placeholder

func (s *GeminiMessagesCompatService) handleNonStreamingResponse(c *gin.Context, resp *http.Response, originalModel string) (*ClaudeUsage, error) {
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, s.writeClaudeError(c, http.StatusBadGateway, "upstream_error", "Failed to read upstream response")
placeholder

	geminiResp, err := unwrapGeminiResponse(body)
	if err != nil {
		return nil, s.writeClaudeError(c, http.StatusBadGateway, "upstream_error", "Failed to parse upstream response")
placeholder

	claudeResp, usage := convertGeminiToClaudeMessage(geminiResp, originalModel)
	c.JSON(http.StatusOK, claudeResp)

	return usage, nil
placeholder

func (s *GeminiMessagesCompatService) handleStreamingResponse(c *gin.Context, resp *http.Response, startTime time.Time, originalModel string) (*geminiStreamResult, error) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return nil, errors.New("streaming not supported")
placeholder

	messageID := "msg_" + randomHex(12)
	messageStart := map[string]any{
		"type": "message_start",
		"message": map[string]any{
			"id":            messageID,
			"type":          "message",
			"role":          "assistant",
			"model":         originalModel,
			"content":       []any{placeholder,
			"stop_reason":   nil,
			"stop_sequence": nil,
			"usage": map[string]any{
				"input_tokens":  0,
				"output_tokens": 0,
		placeholder,
	placeholder,
placeholder
	writeSSE(c.Writer, "message_start", messageStart)
	flusher.Flush()

	var firstTokenMs *int
	var usage ClaudeUsage
	finishReason := ""
	sawToolUse := false

	nextBlockIndex := 0
	openBlockIndex := -1
	openBlockType := ""
	seenText := ""

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
		placeholder
			return nil, fmt.Errorf("stream read error: %w", err)
	placeholder

		if !strings.HasPrefix(line, "data:") {
			continue
	placeholder
		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "" || payload == "[DONE]" {
			continue
	placeholder

		geminiResp, err := unwrapGeminiResponse([]byte(payload))
		if err != nil {
			continue
	placeholder

		if fr := extractGeminiFinishReason(geminiResp); fr != "" {
			finishReason = fr
	placeholder

		parts := extractGeminiParts(geminiResp)
		for _, part := range parts {
			if text, ok := part["text"].(string); ok && text != "" {
				delta, newSeen := computeGeminiTextDelta(seenText, text)
				seenText = newSeen
				if delta == "" {
					continue
			placeholder

				if openBlockType != "text" {
					if openBlockIndex >= 0 {
						writeSSE(c.Writer, "content_block_stop", map[string]any{
							"type":  "content_block_stop",
							"index": openBlockIndex,
					placeholder)
				placeholder
					openBlockType = "text"
					openBlockIndex = nextBlockIndex
					nextBlockIndex++
					writeSSE(c.Writer, "content_block_start", map[string]any{
						"type":  "content_block_start",
						"index": openBlockIndex,
						"content_block": map[string]any{
							"type": "text",
							"text": "",
					placeholder,
				placeholder)
			placeholder

				if firstTokenMs == nil {
					ms := int(time.Since(startTime).Milliseconds())
					firstTokenMs = &ms
			placeholder
				writeSSE(c.Writer, "content_block_delta", map[string]any{
					"type":  "content_block_delta",
					"index": openBlockIndex,
					"delta": map[string]any{
						"type": "text_delta",
						"text": delta,
				placeholder,
			placeholder)
				flusher.Flush()
				continue
		placeholder

			if fc, ok := part["functionCall"].(map[string]any); ok && fc != nil {
				name, _ := fc["name"].(string)
				args := fc["args"]
				if strings.TrimSpace(name) == "" {
					name = "tool"
			placeholder

				// Close any open block before tool_use.
				if openBlockIndex >= 0 {
					writeSSE(c.Writer, "content_block_stop", map[string]any{
						"type":  "content_block_stop",
						"index": openBlockIndex,
				placeholder)
					openBlockIndex = -1
					openBlockType = ""
			placeholder

				toolID := "toolu_" + randomHex(8)
				toolIndex := nextBlockIndex
				nextBlockIndex++
				sawToolUse = true

				writeSSE(c.Writer, "content_block_start", map[string]any{
					"type":  "content_block_start",
					"index": toolIndex,
					"content_block": map[string]any{
						"type":  "tool_use",
						"id":    toolID,
						"name":  name,
						"input": map[string]any{placeholder,
				placeholder,
			placeholder)

				argsJSON := "{placeholder"
				if args != nil {
					if b, err := json.Marshal(args); err == nil {
						argsJSON = string(b)
				placeholder
			placeholder
				writeSSE(c.Writer, "content_block_delta", map[string]any{
					"type":  "content_block_delta",
					"index": toolIndex,
					"delta": map[string]any{
						"type":         "input_json_delta",
						"partial_json": argsJSON,
				placeholder,
			placeholder)
				writeSSE(c.Writer, "content_block_stop", map[string]any{
					"type":  "content_block_stop",
					"index": toolIndex,
			placeholder)
				flusher.Flush()
		placeholder
	placeholder

		if u := extractGeminiUsage(geminiResp); u != nil {
			usage = *u
	placeholder
placeholder

	if openBlockIndex >= 0 {
		writeSSE(c.Writer, "content_block_stop", map[string]any{
			"type":  "content_block_stop",
			"index": openBlockIndex,
	placeholder)
placeholder

	stopReason := mapGeminiFinishReasonToClaudeStopReason(finishReason)
	if sawToolUse {
		stopReason = "tool_use"
placeholder

	usageObj := map[string]any{
		"output_tokens": usage.OutputTokens,
placeholder
	if usage.InputTokens > 0 {
		usageObj["input_tokens"] = usage.InputTokens
placeholder
	writeSSE(c.Writer, "message_delta", map[string]any{
		"type": "message_delta",
		"delta": map[string]any{
			"stop_reason":   stopReason,
			"stop_sequence": nil,
	placeholder,
		"usage": usageObj,
placeholder)
	writeSSE(c.Writer, "message_stop", map[string]any{
		"type": "message_stop",
placeholder)
	flusher.Flush()

	return &geminiStreamResult{usage: &usage, firstTokenMs: firstTokenMsplaceholder, nil
placeholder

func writeSSE(w io.Writer, event string, data any) {
	if event != "" {
		_, _ = fmt.Fprintf(w, "event: %s\n", event)
placeholder
	b, _ := json.Marshal(data)
	_, _ = fmt.Fprintf(w, "data: %s\n\n", string(b))
placeholder

func randomHex(nBytes int) string {
	b := make([]byte, nBytes)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
placeholder

func (s *GeminiMessagesCompatService) writeClaudeError(c *gin.Context, status int, errType, message string) error {
	c.JSON(status, gin.H{
		"type":  "error",
		"error": gin.H{"type": errType, "message": messageplaceholder,
placeholder)
	return fmt.Errorf("%s", message)
placeholder

func unwrapGeminiResponse(raw []byte) (map[string]any, error) {
	var outer map[string]any
	if err := json.Unmarshal(raw, &outer); err != nil {
		return nil, err
placeholder
	if resp, ok := outer["response"].(map[string]any); ok && resp != nil {
		return resp, nil
placeholder
	return outer, nil
placeholder

func convertGeminiToClaudeMessage(geminiResp map[string]any, originalModel string) (map[string]any, *ClaudeUsage) {
	usage := extractGeminiUsage(geminiResp)
	if usage == nil {
		usage = &ClaudeUsage{placeholder
placeholder

	contentBlocks := make([]any, 0)
	sawToolUse := false
	if candidates, ok := geminiResp["candidates"].([]any); ok && len(candidates) > 0 {
		if cand, ok := candidates[0].(map[string]any); ok {
			if content, ok := cand["content"].(map[string]any); ok {
				if parts, ok := content["parts"].([]any); ok {
					for _, part := range parts {
						pm, ok := part.(map[string]any)
						if !ok {
							continue
					placeholder
						if text, ok := pm["text"].(string); ok && text != "" {
							contentBlocks = append(contentBlocks, map[string]any{
								"type": "text",
								"text": text,
						placeholder)
					placeholder
						if fc, ok := pm["functionCall"].(map[string]any); ok {
							name, _ := fc["name"].(string)
							if strings.TrimSpace(name) == "" {
								name = "tool"
						placeholder
							args := fc["args"]
							sawToolUse = true
							contentBlocks = append(contentBlocks, map[string]any{
								"type":  "tool_use",
								"id":    "toolu_" + randomHex(8),
								"name":  name,
								"input": args,
						placeholder)
					placeholder
				placeholder
			placeholder
		placeholder
	placeholder
placeholder

	stopReason := mapGeminiFinishReasonToClaudeStopReason(extractGeminiFinishReason(geminiResp))
	if sawToolUse {
		stopReason = "tool_use"
placeholder

	resp := map[string]any{
		"id":            "msg_" + randomHex(12),
		"type":          "message",
		"role":          "assistant",
		"model":         originalModel,
		"content":       contentBlocks,
		"stop_reason":   stopReason,
		"stop_sequence": nil,
		"usage": map[string]any{
			"input_tokens":  usage.InputTokens,
			"output_tokens": usage.OutputTokens,
	placeholder,
placeholder

	return resp, usage
placeholder

func extractGeminiUsage(geminiResp map[string]any) *ClaudeUsage {
	usageMeta, ok := geminiResp["usageMetadata"].(map[string]any)
	if !ok || usageMeta == nil {
		return nil
placeholder
	prompt, _ := asInt(usageMeta["promptTokenCount"])
	cand, _ := asInt(usageMeta["candidatesTokenCount"])
	return &ClaudeUsage{
		InputTokens:  prompt,
		OutputTokens: cand,
placeholder
placeholder

func asInt(v any) (int, bool) {
	switch t := v.(type) {
	case float64:
		return int(t), true
	case int:
		return t, true
	case int64:
		return int(t), true
	case json.Number:
		i, err := t.Int64()
		if err != nil {
			return 0, false
	placeholder
		return int(i), true
	default:
		return 0, false
placeholder
placeholder

func (s *GeminiMessagesCompatService) handleGeminiUpstreamError(ctx context.Context, account *model.Account, statusCode int, headers http.Header, body []byte) {
	if s.rateLimitService != nil && (statusCode == 401 || statusCode == 403 || statusCode == 529) {
		s.rateLimitService.HandleUpstreamError(ctx, account, statusCode, headers, body)
		return
placeholder
	if statusCode != 429 {
		return
placeholder
	resetAt := parseGeminiRateLimitResetTime(body)
	if resetAt == nil {
		ra := time.Now().Add(5 * time.Minute)
		_ = s.accountRepo.SetRateLimited(ctx, account.ID, ra)
		return
placeholder
	_ = s.accountRepo.SetRateLimited(ctx, account.ID, time.Unix(*resetAt, 0))
placeholder

func parseGeminiRateLimitResetTime(body []byte) *int64 {
	// Try to parse metadata.quotaResetDelay like "12.345s"
	var parsed map[string]any
	if err := json.Unmarshal(body, &parsed); err == nil {
		if errObj, ok := parsed["error"].(map[string]any); ok {
			if msg, ok := errObj["message"].(string); ok {
				if looksLikeGeminiDailyQuota(msg) {
					if ts := nextGeminiDailyResetUnix(); ts != nil {
						return ts
				placeholder
			placeholder
		placeholder
			if details, ok := errObj["details"].([]any); ok {
				for _, d := range details {
					dm, ok := d.(map[string]any)
					if !ok {
						continue
				placeholder
					if meta, ok := dm["metadata"].(map[string]any); ok {
						if v, ok := meta["quotaResetDelay"].(string); ok {
							if dur, err := time.ParseDuration(v); err == nil {
								ts := time.Now().Unix() + int64(dur.Seconds())
								return &ts
						placeholder
					placeholder
				placeholder
			placeholder
		placeholder
	placeholder
placeholder

	// Match "Please retry in Xs"
	retryInRegex := regexp.MustCompile(`Please retry in ([0-9.]+)s`)
	matches := retryInRegex.FindStringSubmatch(string(body))
	if len(matches) == 2 {
		if dur, err := time.ParseDuration(matches[1] + "s"); err == nil {
			ts := time.Now().Unix() + int64(math.Ceil(dur.Seconds()))
			return &ts
	placeholder
placeholder

	return nil
placeholder

func looksLikeGeminiDailyQuota(message string) bool {
	m := strings.ToLower(message)
	if strings.Contains(m, "per day") || strings.Contains(m, "requests per day") || strings.Contains(m, "quota") && strings.Contains(m, "per day") {
		return true
placeholder
	return false
placeholder

func nextGeminiDailyResetUnix() *int64 {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		// Fallback: PST without DST.
		loc = time.FixedZone("PST", -8*3600)
placeholder
	now := time.Now().In(loc)
	reset := time.Date(now.Year(), now.Month(), now.Day(), 0, 5, 0, 0, loc)
	if !reset.After(now) {
		reset = reset.Add(24 * time.Hour)
placeholder
	ts := reset.Unix()
	return &ts
placeholder

func extractGeminiFinishReason(geminiResp map[string]any) string {
	if candidates, ok := geminiResp["candidates"].([]any); ok && len(candidates) > 0 {
		if cand, ok := candidates[0].(map[string]any); ok {
			if fr, ok := cand["finishReason"].(string); ok {
				return fr
		placeholder
	placeholder
placeholder
	return ""
placeholder

func extractGeminiParts(geminiResp map[string]any) []map[string]any {
	if candidates, ok := geminiResp["candidates"].([]any); ok && len(candidates) > 0 {
		if cand, ok := candidates[0].(map[string]any); ok {
			if content, ok := cand["content"].(map[string]any); ok {
				if partsAny, ok := content["parts"].([]any); ok && len(partsAny) > 0 {
					out := make([]map[string]any, 0, len(partsAny))
					for _, p := range partsAny {
						pm, ok := p.(map[string]any)
						if !ok {
							continue
					placeholder
						out = append(out, pm)
				placeholder
					return out
			placeholder
		placeholder
	placeholder
placeholder
	return nil
placeholder

func computeGeminiTextDelta(seen, incoming string) (delta, newSeen string) {
	incoming = strings.TrimSuffix(incoming, "\u0000")
	if incoming == "" {
		return "", seen
placeholder

	// Cumulative mode: incoming contains full text so far.
	if strings.HasPrefix(incoming, seen) {
		return strings.TrimPrefix(incoming, seen), incoming
placeholder
	// Duplicate/rewind: ignore.
	if strings.HasPrefix(seen, incoming) {
		return "", seen
placeholder
	// Delta mode: treat incoming as incremental chunk.
	return incoming, seen + incoming
placeholder

func mapGeminiFinishReasonToClaudeStopReason(finishReason string) string {
	switch strings.ToUpper(strings.TrimSpace(finishReason)) {
	case "MAX_TOKENS":
		return "max_tokens"
	case "STOP":
		return "end_turn"
	default:
		return "end_turn"
placeholder
placeholder

func convertClaudeMessagesToGeminiGenerateContent(body []byte) ([]byte, error) {
	var req map[string]any
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
placeholder

	toolUseIDToName := make(map[string]string)

	systemText := extractClaudeSystemText(req["system"])
	contents, err := convertClaudeMessagesToGeminiContents(req["messages"], toolUseIDToName)
	if err != nil {
		return nil, err
placeholder

	out := make(map[string]any)
	if systemText != "" {
		out["systemInstruction"] = map[string]any{
			"parts": []any{map[string]any{"text": systemTextplaceholderplaceholder,
	placeholder
placeholder
	out["contents"] = contents

	if tools := convertClaudeToolsToGeminiTools(req["tools"]); tools != nil {
		out["tools"] = tools
placeholder

	generationConfig := convertClaudeGenerationConfig(req)
	if generationConfig != nil {
		out["generationConfig"] = generationConfig
placeholder

	stripGeminiFunctionIDs(out)
	return json.Marshal(out)
placeholder

func stripGeminiFunctionIDs(req map[string]any) {
	// Defensive cleanup: some upstreams reject unexpected `id` fields in functionCall/functionResponse.
	contents, ok := req["contents"].([]any)
	if !ok {
		return
placeholder
	for _, c := range contents {
		cm, ok := c.(map[string]any)
		if !ok {
			continue
	placeholder
		contentParts, ok := cm["parts"].([]any)
		if !ok {
			continue
	placeholder
		for _, p := range contentParts {
			pm, ok := p.(map[string]any)
			if !ok {
				continue
		placeholder
			if fc, ok := pm["functionCall"].(map[string]any); ok && fc != nil {
				delete(fc, "id")
		placeholder
			if fr, ok := pm["functionResponse"].(map[string]any); ok && fr != nil {
				delete(fr, "id")
		placeholder
	placeholder
placeholder
placeholder

func extractClaudeSystemText(system any) string {
	switch v := system.(type) {
	case string:
		return strings.TrimSpace(v)
	case []any:
		var parts []string
		for _, p := range v {
			pm, ok := p.(map[string]any)
			if !ok {
				continue
		placeholder
			if t, _ := pm["type"].(string); t != "text" {
				continue
		placeholder
			if text, ok := pm["text"].(string); ok && strings.TrimSpace(text) != "" {
				parts = append(parts, text)
		placeholder
	placeholder
		return strings.TrimSpace(strings.Join(parts, "\n"))
	default:
		return ""
placeholder
placeholder

func convertClaudeMessagesToGeminiContents(messages any, toolUseIDToName map[string]string) ([]any, error) {
	arr, ok := messages.([]any)
	if !ok {
		return nil, errors.New("messages must be an array")
placeholder

	out := make([]any, 0, len(arr))
	for _, m := range arr {
		mm, ok := m.(map[string]any)
		if !ok {
			continue
	placeholder
		role, _ := mm["role"].(string)
		role = strings.ToLower(strings.TrimSpace(role))
		gRole := "user"
		if role == "assistant" {
			gRole = "model"
	placeholder

		parts := make([]any, 0)
		switch content := mm["content"].(type) {
		case string:
			if strings.TrimSpace(content) != "" {
				parts = append(parts, map[string]any{"text": contentplaceholder)
		placeholder
		case []any:
			for _, block := range content {
				bm, ok := block.(map[string]any)
				if !ok {
					continue
			placeholder
				bt, _ := bm["type"].(string)
				switch bt {
				case "text":
					if text, ok := bm["text"].(string); ok && strings.TrimSpace(text) != "" {
						parts = append(parts, map[string]any{"text": textplaceholder)
				placeholder
				case "tool_use":
					id, _ := bm["id"].(string)
					name, _ := bm["name"].(string)
					if strings.TrimSpace(id) != "" && strings.TrimSpace(name) != "" {
						toolUseIDToName[id] = name
				placeholder
					parts = append(parts, map[string]any{
						"functionCall": map[string]any{
							"name": name,
							"args": bm["input"],
					placeholder,
				placeholder)
				case "tool_result":
					toolUseID, _ := bm["tool_use_id"].(string)
					name := toolUseIDToName[toolUseID]
					if name == "" {
						name = "tool"
				placeholder
					parts = append(parts, map[string]any{
						"functionResponse": map[string]any{
							"name": name,
							"response": map[string]any{
								"content": extractClaudeContentText(bm["content"]),
						placeholder,
					placeholder,
				placeholder)
				case "image":
					if src, ok := bm["source"].(map[string]any); ok {
						if srcType, _ := src["type"].(string); srcType == "base64" {
							mediaType, _ := src["media_type"].(string)
							data, _ := src["data"].(string)
							if mediaType != "" && data != "" {
								parts = append(parts, map[string]any{
									"inlineData": map[string]any{
										"mimeType": mediaType,
										"data":     data,
								placeholder,
							placeholder)
						placeholder
					placeholder
				placeholder
				default:
					// best-effort: preserve unknown blocks as text
					if b, err := json.Marshal(bm); err == nil {
						parts = append(parts, map[string]any{"text": string(b)placeholder)
				placeholder
			placeholder
		placeholder
		default:
			// ignore
	placeholder

		out = append(out, map[string]any{
			"role":  gRole,
			"parts": parts,
	placeholder)
placeholder
	return out, nil
placeholder

func extractClaudeContentText(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case []any:
		var sb strings.Builder
		for _, part := range t {
			pm, ok := part.(map[string]any)
			if !ok {
				continue
		placeholder
			if pm["type"] == "text" {
				if text, ok := pm["text"].(string); ok {
					sb.WriteString(text)
			placeholder
		placeholder
	placeholder
		return sb.String()
	default:
		b, _ := json.Marshal(t)
		return string(b)
placeholder
placeholder

func convertClaudeToolsToGeminiTools(tools any) []any {
	arr, ok := tools.([]any)
	if !ok || len(arr) == 0 {
		return nil
placeholder

	funcDecls := make([]any, 0, len(arr))
	for _, t := range arr {
		tm, ok := t.(map[string]any)
		if !ok {
			continue
	placeholder
		name, _ := tm["name"].(string)
		desc, _ := tm["description"].(string)
		params := tm["input_schema"]
		if name == "" {
			continue
	placeholder
		funcDecls = append(funcDecls, map[string]any{
			"name":        name,
			"description": desc,
			"parameters":  params,
	placeholder)
placeholder

	if len(funcDecls) == 0 {
		return nil
placeholder
	return []any{
		map[string]any{
			"functionDeclarations": funcDecls,
	placeholder,
placeholder
placeholder

func convertClaudeGenerationConfig(req map[string]any) map[string]any {
	out := make(map[string]any)
	if mt, ok := asInt(req["max_tokens"]); ok && mt > 0 {
		out["maxOutputTokens"] = mt
placeholder
	if temp, ok := req["temperature"].(float64); ok {
		out["temperature"] = temp
placeholder
	if topP, ok := req["top_p"].(float64); ok {
		out["topP"] = topP
placeholder
	if stopSeq, ok := req["stop_sequences"].([]any); ok && len(stopSeq) > 0 {
		out["stopSequences"] = stopSeq
placeholder
	if len(out) == 0 {
		return nil
placeholder
	return out
placeholder
