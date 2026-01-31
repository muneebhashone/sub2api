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
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
)

var soraSSEDataRe = regexp.MustCompile(`^data:\s*`)
var soraImageMarkdownRe = regexp.MustCompile(`!\[[^\]]*\]\(([^)]+)\)`)
var soraVideoHTMLRe = regexp.MustCompile(`(?i)<video[^>]+src=['"]([^'"]+)['"]`)

var soraImageSizeMap = map[string]string{
	"gpt-image":           "360",
	"gpt-image-landscape": "540",
	"gpt-image-portrait":  "540",
placeholder

type soraStreamingResult struct {
	content      string
	mediaType    string
	mediaURLs    []string
	imageCount   int
	imageSize    string
	firstTokenMs *int
placeholder

// SoraGatewayService handles forwarding requests to sora2api.
type SoraGatewayService struct {
	sora2api         *Sora2APIService
	httpUpstream     HTTPUpstream
	rateLimitService *RateLimitService
	cfg              *config.Config
placeholder

func NewSoraGatewayService(
	sora2api *Sora2APIService,
	httpUpstream HTTPUpstream,
	rateLimitService *RateLimitService,
	cfg *config.Config,
) *SoraGatewayService {
	return &SoraGatewayService{
		sora2api:         sora2api,
		httpUpstream:     httpUpstream,
		rateLimitService: rateLimitService,
		cfg:              cfg,
placeholder
placeholder

func (s *SoraGatewayService) Forward(ctx context.Context, c *gin.Context, account *Account, body []byte, clientStream bool) (*ForwardResult, error) {
	startTime := time.Now()

	if s.sora2api == nil || !s.sora2api.Enabled() {
		if c != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": gin.H{
					"type":    "api_error",
					"message": "sora2api 未配置",
			placeholder,
		placeholder)
	placeholder
		return nil, errors.New("sora2api not configured")
placeholder

	var reqBody map[string]any
	if err := json.Unmarshal(body, &reqBody); err != nil {
		return nil, fmt.Errorf("parse request: %w", err)
placeholder
	reqModel, _ := reqBody["model"].(string)
	reqStream, _ := reqBody["stream"].(bool)

	mappedModel := account.GetMappedModel(reqModel)
	if mappedModel != reqModel && mappedModel != "" {
		reqBody["model"] = mappedModel
		if updated, err := json.Marshal(reqBody); err == nil {
			body = updated
	placeholder
placeholder

	reqCtx, cancel := s.withSoraTimeout(ctx, reqStream)
	if cancel != nil {
		defer cancel()
placeholder

	upstreamReq, err := s.sora2api.NewAPIRequest(reqCtx, http.MethodPost, "/v1/chat/completions", body)
	if err != nil {
		return nil, err
placeholder
	if c != nil {
		if ua := strings.TrimSpace(c.GetHeader("User-Agent")); ua != "" {
			upstreamReq.Header.Set("User-Agent", ua)
	placeholder
placeholder
	if reqStream {
		upstreamReq.Header.Set("Accept", "text/event-stream")
placeholder

	if c != nil {
		c.Set(OpsUpstreamRequestBodyKey, string(body))
placeholder

	proxyURL := ""
	if account != nil && account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
placeholder

	var resp *http.Response
	if s.httpUpstream != nil {
		resp, err = s.httpUpstream.Do(upstreamReq, proxyURL, account.ID, account.Concurrency)
placeholder else {
		resp, err = http.DefaultClient.Do(upstreamReq)
placeholder
	if err != nil {
		s.setUpstreamRequestError(c, account, err)
		return nil, fmt.Errorf("upstream request failed: %w", err)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	if resp.StatusCode >= 400 {
		if s.shouldFailoverUpstreamError(resp.StatusCode) {
			respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
			_ = resp.Body.Close()
			resp.Body = io.NopCloser(bytes.NewReader(respBody))
			upstreamMsg := strings.TrimSpace(extractUpstreamErrorMessage(respBody))
			upstreamMsg = sanitizeUpstreamErrorMessage(upstreamMsg)
			appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
				Platform:           account.Platform,
				AccountID:          account.ID,
				AccountName:        account.Name,
				UpstreamStatusCode: resp.StatusCode,
				UpstreamRequestID:  resp.Header.Get("x-request-id"),
				Kind:               "failover",
				Message:            upstreamMsg,
		placeholder)
			s.handleFailoverSideEffects(ctx, resp, account)
			return nil, &UpstreamFailoverError{StatusCode: resp.StatusCodeplaceholder
	placeholder
		return s.handleErrorResponse(ctx, resp, c, account, reqModel)
placeholder

	streamResult, err := s.handleStreamingResponse(ctx, resp, c, account, startTime, reqModel, clientStream)
	if err != nil {
		return nil, err
placeholder

	result := &ForwardResult{
		RequestID:    resp.Header.Get("x-request-id"),
		Model:        reqModel,
		Stream:       clientStream,
		Duration:     time.Since(startTime),
		FirstTokenMs: streamResult.firstTokenMs,
		Usage:        ClaudeUsage{placeholder,
		MediaType:    streamResult.mediaType,
		MediaURL:     firstMediaURL(streamResult.mediaURLs),
		ImageCount:   streamResult.imageCount,
		ImageSize:    streamResult.imageSize,
placeholder

	return result, nil
placeholder

func (s *SoraGatewayService) withSoraTimeout(ctx context.Context, stream bool) (context.Context, context.CancelFunc) {
	if s == nil || s.cfg == nil {
		return ctx, nil
placeholder
	timeoutSeconds := s.cfg.Gateway.SoraRequestTimeoutSeconds
	if stream {
		timeoutSeconds = s.cfg.Gateway.SoraStreamTimeoutSeconds
placeholder
	if timeoutSeconds <= 0 {
		return ctx, nil
placeholder
	return context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
placeholder

func (s *SoraGatewayService) setUpstreamRequestError(c *gin.Context, account *Account, err error) {
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
	if c != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": gin.H{
				"type":    "upstream_error",
				"message": "Upstream request failed",
		placeholder,
	placeholder)
placeholder
placeholder

func (s *SoraGatewayService) shouldFailoverUpstreamError(statusCode int) bool {
	switch statusCode {
	case 401, 402, 403, 429, 529:
		return true
	default:
		return statusCode >= 500
placeholder
placeholder

func (s *SoraGatewayService) handleFailoverSideEffects(ctx context.Context, resp *http.Response, account *Account) {
	if s.rateLimitService == nil || account == nil || resp == nil {
		return
placeholder
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	s.rateLimitService.HandleUpstreamError(ctx, account, resp.StatusCode, resp.Header, body)
placeholder

func (s *SoraGatewayService) handleErrorResponse(ctx context.Context, resp *http.Response, c *gin.Context, account *Account, reqModel string) (*ForwardResult, error) {
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	_ = resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(respBody))

	upstreamMsg := strings.TrimSpace(extractUpstreamErrorMessage(respBody))
	upstreamMsg = sanitizeUpstreamErrorMessage(upstreamMsg)
	if msg := soraProErrorMessage(reqModel, upstreamMsg); msg != "" {
		upstreamMsg = msg
placeholder

	upstreamDetail := ""
	if s.cfg != nil && s.cfg.Gateway.LogUpstreamErrorBody {
		maxBytes := s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes
		if maxBytes <= 0 {
			maxBytes = 2048
	placeholder
		upstreamDetail = truncateString(string(respBody), maxBytes)
placeholder
	setOpsUpstreamError(c, resp.StatusCode, upstreamMsg, upstreamDetail)
	appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
		Platform:           account.Platform,
		AccountID:          account.ID,
		AccountName:        account.Name,
		UpstreamStatusCode: resp.StatusCode,
		UpstreamRequestID:  resp.Header.Get("x-request-id"),
		Kind:               "http_error",
		Message:            upstreamMsg,
		Detail:             upstreamDetail,
placeholder)

	if c != nil {
		responsePayload := s.buildErrorPayload(respBody, upstreamMsg)
		c.JSON(resp.StatusCode, responsePayload)
placeholder
	if upstreamMsg == "" {
		return nil, fmt.Errorf("upstream error: %d", resp.StatusCode)
placeholder
	return nil, fmt.Errorf("upstream error: %d message=%s", resp.StatusCode, upstreamMsg)
placeholder

func (s *SoraGatewayService) buildErrorPayload(respBody []byte, overrideMessage string) map[string]any {
	if len(respBody) > 0 {
		var payload map[string]any
		if err := json.Unmarshal(respBody, &payload); err == nil {
			if errObj, ok := payload["error"].(map[string]any); ok {
				if overrideMessage != "" {
					errObj["message"] = overrideMessage
			placeholder
				payload["error"] = errObj
				return payload
		placeholder
	placeholder
placeholder
	return map[string]any{
		"error": map[string]any{
			"type":    "upstream_error",
			"message": overrideMessage,
	placeholder,
placeholder
placeholder

func (s *SoraGatewayService) handleStreamingResponse(ctx context.Context, resp *http.Response, c *gin.Context, account *Account, startTime time.Time, originalModel string, clientStream bool) (*soraStreamingResult, error) {
	if resp == nil {
		return nil, errors.New("empty response")
placeholder

	if clientStream {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("X-Accel-Buffering", "no")
		if v := resp.Header.Get("x-request-id"); v != "" {
			c.Header("x-request-id", v)
	placeholder
placeholder

	w := c.Writer
	flusher, _ := w.(http.Flusher)

	contentBuilder := strings.Builder{placeholder
	var firstTokenMs *int
	var upstreamError error

	scanner := bufio.NewScanner(resp.Body)
	maxLineSize := defaultMaxLineSize
	if s.cfg != nil && s.cfg.Gateway.MaxLineSize > 0 {
		maxLineSize = s.cfg.Gateway.MaxLineSize
placeholder
	scanner.Buffer(make([]byte, 64*1024), maxLineSize)

	sendLine := func(line string) error {
		if !clientStream {
			return nil
	placeholder
		if _, err := fmt.Fprintf(w, "%s\n", line); err != nil {
			return err
	placeholder
		if flusher != nil {
			flusher.Flush()
	placeholder
		return nil
placeholder

	for scanner.Scan() {
		line := scanner.Text()
		if soraSSEDataRe.MatchString(line) {
			data := soraSSEDataRe.ReplaceAllString(line, "")
			if data == "[DONE]" {
				if err := sendLine("data: [DONE]"); err != nil {
					return nil, err
			placeholder
				break
		placeholder
			updatedLine, contentDelta, errEvent := s.processSoraSSEData(data, originalModel)
			if errEvent != nil && upstreamError == nil {
				upstreamError = errEvent
		placeholder
			if contentDelta != "" {
				if firstTokenMs == nil {
					ms := int(time.Since(startTime).Milliseconds())
					firstTokenMs = &ms
			placeholder
				contentBuilder.WriteString(contentDelta)
		placeholder
			if err := sendLine(updatedLine); err != nil {
				return nil, err
		placeholder
			continue
	placeholder
		if err := sendLine(line); err != nil {
			return nil, err
	placeholder
placeholder

	if err := scanner.Err(); err != nil {
		if errors.Is(err, bufio.ErrTooLong) {
			if clientStream {
				_, _ = fmt.Fprintf(w, "event: error\ndata: {\"error\":\"response_too_large\"placeholder\n\n")
				if flusher != nil {
					flusher.Flush()
			placeholder
		placeholder
			return nil, err
	placeholder
		if ctx.Err() == context.DeadlineExceeded && s.rateLimitService != nil && account != nil {
			s.rateLimitService.HandleStreamTimeout(ctx, account, originalModel)
	placeholder
		if clientStream {
			_, _ = fmt.Fprintf(w, "event: error\ndata: {\"error\":\"stream_read_error\"placeholder\n\n")
			if flusher != nil {
				flusher.Flush()
		placeholder
	placeholder
		return nil, err
placeholder

	content := contentBuilder.String()
	mediaType, mediaURLs := s.extractSoraMedia(content)
	if mediaType == "" && isSoraPromptEnhanceModel(originalModel) {
		mediaType = "prompt"
placeholder
	imageSize := ""
	imageCount := 0
	if mediaType == "image" {
		imageSize = soraImageSizeFromModel(originalModel)
		imageCount = len(mediaURLs)
placeholder

	if upstreamError != nil && !clientStream {
		if c != nil {
			c.JSON(http.StatusBadGateway, map[string]any{
				"error": map[string]any{
					"type":    "upstream_error",
					"message": upstreamError.Error(),
			placeholder,
		placeholder)
	placeholder
		return nil, upstreamError
placeholder

	if !clientStream {
		response := buildSoraNonStreamResponse(content, originalModel)
		if len(mediaURLs) > 0 {
			response["media_url"] = mediaURLs[0]
			if len(mediaURLs) > 1 {
				response["media_urls"] = mediaURLs
		placeholder
	placeholder
		c.JSON(http.StatusOK, response)
placeholder

	return &soraStreamingResult{
		content:      content,
		mediaType:    mediaType,
		mediaURLs:    mediaURLs,
		imageCount:   imageCount,
		imageSize:    imageSize,
		firstTokenMs: firstTokenMs,
placeholder, nil
placeholder

func (s *SoraGatewayService) processSoraSSEData(data string, originalModel string) (string, string, error) {
	if strings.TrimSpace(data) == "" {
		return "data: ", "", nil
placeholder

	var payload map[string]any
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		return "data: " + data, "", nil
placeholder

	if errObj, ok := payload["error"].(map[string]any); ok {
		if msg, ok := errObj["message"].(string); ok && strings.TrimSpace(msg) != "" {
			return "data: " + data, "", errors.New(msg)
	placeholder
placeholder

	if model, ok := payload["model"].(string); ok && model != "" && originalModel != "" {
		payload["model"] = originalModel
placeholder

	contentDelta, updated := extractSoraContent(payload)
	if updated {
		rewritten := s.rewriteSoraContent(contentDelta)
		if rewritten != contentDelta {
			applySoraContent(payload, rewritten)
			contentDelta = rewritten
	placeholder
placeholder

	updatedData, err := json.Marshal(payload)
	if err != nil {
		return "data: " + data, contentDelta, nil
placeholder
	return "data: " + string(updatedData), contentDelta, nil
placeholder

func extractSoraContent(payload map[string]any) (string, bool) {
	choices, ok := payload["choices"].([]any)
	if !ok || len(choices) == 0 {
		return "", false
placeholder
	choice, ok := choices[0].(map[string]any)
	if !ok {
		return "", false
placeholder
	if delta, ok := choice["delta"].(map[string]any); ok {
		if content, ok := delta["content"].(string); ok {
			return content, true
	placeholder
placeholder
	if message, ok := choice["message"].(map[string]any); ok {
		if content, ok := message["content"].(string); ok {
			return content, true
	placeholder
placeholder
	return "", false
placeholder

func applySoraContent(payload map[string]any, content string) {
	choices, ok := payload["choices"].([]any)
	if !ok || len(choices) == 0 {
		return
placeholder
	choice, ok := choices[0].(map[string]any)
	if !ok {
		return
placeholder
	if delta, ok := choice["delta"].(map[string]any); ok {
		delta["content"] = content
		choice["delta"] = delta
		return
placeholder
	if message, ok := choice["message"].(map[string]any); ok {
		message["content"] = content
		choice["message"] = message
placeholder
placeholder

func (s *SoraGatewayService) rewriteSoraContent(content string) string {
	if content == "" {
		return content
placeholder
	content = soraImageMarkdownRe.ReplaceAllStringFunc(content, func(match string) string {
		sub := soraImageMarkdownRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
	placeholder
		rewritten := s.rewriteSoraURL(sub[1])
		if rewritten == sub[1] {
			return match
	placeholder
		return strings.Replace(match, sub[1], rewritten, 1)
placeholder)
	content = soraVideoHTMLRe.ReplaceAllStringFunc(content, func(match string) string {
		sub := soraVideoHTMLRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
	placeholder
		rewritten := s.rewriteSoraURL(sub[1])
		if rewritten == sub[1] {
			return match
	placeholder
		return strings.Replace(match, sub[1], rewritten, 1)
placeholder)
	return content
placeholder

func (s *SoraGatewayService) rewriteSoraURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw
placeholder
	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
placeholder
	path := parsed.Path
	if !strings.HasPrefix(path, "/tmp/") && !strings.HasPrefix(path, "/static/") {
		return raw
placeholder
	return s.buildSoraMediaURL(path, parsed.RawQuery)
placeholder

func (s *SoraGatewayService) extractSoraMedia(content string) (string, []string) {
	if content == "" {
		return "", nil
placeholder
	if match := soraVideoHTMLRe.FindStringSubmatch(content); len(match) > 1 {
		return "video", []string{match[1]placeholder
placeholder
	imageMatches := soraImageMarkdownRe.FindAllStringSubmatch(content, -1)
	if len(imageMatches) == 0 {
		return "", nil
placeholder
	urls := make([]string, 0, len(imageMatches))
	for _, match := range imageMatches {
		if len(match) > 1 {
			urls = append(urls, match[1])
	placeholder
placeholder
	return "image", urls
placeholder

func buildSoraNonStreamResponse(content, model string) map[string]any {
	return map[string]any{
		"id":      fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano()),
		"object":  "chat.completion",
		"created": time.Now().Unix(),
		"model":   model,
		"choices": []any{
			map[string]any{
				"index": 0,
				"message": map[string]any{
					"role":    "assistant",
					"content": content,
			placeholder,
				"finish_reason": "stop",
		placeholder,
	placeholder,
placeholder
placeholder

func soraImageSizeFromModel(model string) string {
	modelLower := strings.ToLower(model)
	if size, ok := soraImageSizeMap[modelLower]; ok {
		return size
placeholder
	if strings.Contains(modelLower, "landscape") || strings.Contains(modelLower, "portrait") {
		return "540"
placeholder
	return "360"
placeholder

func isSoraPromptEnhanceModel(model string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(model)), "prompt-enhance")
placeholder

func soraProErrorMessage(model, upstreamMsg string) string {
	modelLower := strings.ToLower(model)
	if strings.Contains(modelLower, "sora2pro-hd") {
		return "当前账号无法使用 Sora Pro-HD 模型，请更换模型或账号"
placeholder
	if strings.Contains(modelLower, "sora2pro") {
		return "当前账号无法使用 Sora Pro 模型，请更换模型或账号"
placeholder
	return ""
placeholder

func firstMediaURL(urls []string) string {
	if len(urls) == 0 {
		return ""
placeholder
	return urls[0]
placeholder

func (s *SoraGatewayService) buildSoraMediaURL(path string, rawQuery string) string {
	if path == "" {
		return path
placeholder
	prefix := "/sora/media"
	values := url.Values{placeholder
	if rawQuery != "" {
		if parsed, err := url.ParseQuery(rawQuery); err == nil {
			values = parsed
	placeholder
placeholder

	signKey := ""
	ttlSeconds := 0
	if s != nil && s.cfg != nil {
		signKey = strings.TrimSpace(s.cfg.Gateway.SoraMediaSigningKey)
		ttlSeconds = s.cfg.Gateway.SoraMediaSignedURLTTLSeconds
placeholder
	values.Del("sig")
	values.Del("expires")
	signingQuery := values.Encode()
	if signKey != "" && ttlSeconds > 0 {
		expires := time.Now().Add(time.Duration(ttlSeconds) * time.Second).Unix()
		signature := SignSoraMediaURL(path, signingQuery, expires, signKey)
		if signature != "" {
			values.Set("expires", strconv.FormatInt(expires, 10))
			values.Set("sig", signature)
			prefix = "/sora/media-signed"
	placeholder
placeholder

	encoded := values.Encode()
	if encoded == "" {
		return prefix + path
placeholder
	return prefix + path + "?" + encoded
placeholder
