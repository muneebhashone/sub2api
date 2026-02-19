package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
)

const soraImageInputMaxBytes = 20 << 20
const soraImageInputMaxRedirects = 3
const soraImageInputTimeout = 20 * time.Second

var soraImageSizeMap = map[string]string{
	"gpt-image":           "360",
	"gpt-image-landscape": "540",
	"gpt-image-portrait":  "540",
placeholder

var soraBlockedHostnames = map[string]struct{placeholder{
	"localhost":                 {placeholder,
	"localhost.localdomain":     {placeholder,
	"metadata.google.internal":  {placeholder,
	"metadata.google.internal.": {placeholder,
placeholder

var soraBlockedCIDRs = mustParseCIDRs([]string{
	"0.0.0.0/8",
	"10.0.0.0/8",
	"100.64.0.0/10",
	"127.0.0.0/8",
	"169.254.0.0/16",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"224.0.0.0/4",
	"240.0.0.0/4",
	"::/128",
	"::1/128",
	"fc00::/7",
	"fe80::/10",
placeholder)

// SoraGatewayService handles forwarding requests to Sora upstream.
type SoraGatewayService struct {
	soraClient       SoraClient
	mediaStorage     *SoraMediaStorage
	rateLimitService *RateLimitService
	cfg              *config.Config
placeholder

type soraPreflightChecker interface {
	PreflightCheck(ctx context.Context, account *Account, requestedModel string, modelCfg SoraModelConfig) error
placeholder

func NewSoraGatewayService(
	soraClient SoraClient,
	mediaStorage *SoraMediaStorage,
	rateLimitService *RateLimitService,
	cfg *config.Config,
) *SoraGatewayService {
	return &SoraGatewayService{
		soraClient:       soraClient,
		mediaStorage:     mediaStorage,
		rateLimitService: rateLimitService,
		cfg:              cfg,
placeholder
placeholder

func (s *SoraGatewayService) Forward(ctx context.Context, c *gin.Context, account *Account, body []byte, clientStream bool) (*ForwardResult, error) {
	startTime := time.Now()

	if s.soraClient == nil || !s.soraClient.Enabled() {
		if c != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": gin.H{
					"type":    "api_error",
					"message": "Sora 上游未配置",
			placeholder,
		placeholder)
	placeholder
		return nil, errors.New("sora upstream not configured")
placeholder

	var reqBody map[string]any
	if err := json.Unmarshal(body, &reqBody); err != nil {
		s.writeSoraError(c, http.StatusBadRequest, "invalid_request_error", "Failed to parse request body", clientStream)
		return nil, fmt.Errorf("parse request: %w", err)
placeholder
	reqModel, _ := reqBody["model"].(string)
	reqStream, _ := reqBody["stream"].(bool)
	if strings.TrimSpace(reqModel) == "" {
		s.writeSoraError(c, http.StatusBadRequest, "invalid_request_error", "model is required", clientStream)
		return nil, errors.New("model is required")
placeholder

	mappedModel := account.GetMappedModel(reqModel)
	if mappedModel != "" && mappedModel != reqModel {
		reqModel = mappedModel
placeholder

	modelCfg, ok := GetSoraModelConfig(reqModel)
	if !ok {
		s.writeSoraError(c, http.StatusBadRequest, "invalid_request_error", "Unsupported Sora model", clientStream)
		return nil, fmt.Errorf("unsupported model: %s", reqModel)
placeholder
	prompt, imageInput, videoInput, remixTargetID := extractSoraInput(reqBody)
	if strings.TrimSpace(prompt) == "" {
		s.writeSoraError(c, http.StatusBadRequest, "invalid_request_error", "prompt is required", clientStream)
		return nil, errors.New("prompt is required")
placeholder
	if strings.TrimSpace(videoInput) != "" {
		s.writeSoraError(c, http.StatusBadRequest, "invalid_request_error", "Video input is not supported yet", clientStream)
		return nil, errors.New("video input not supported")
placeholder

	reqCtx, cancel := s.withSoraTimeout(ctx, reqStream)
	if cancel != nil {
		defer cancel()
placeholder
	if checker, ok := s.soraClient.(soraPreflightChecker); ok {
		if err := checker.PreflightCheck(reqCtx, account, reqModel, modelCfg); err != nil {
			return nil, s.handleSoraRequestError(ctx, account, err, reqModel, c, clientStream)
	placeholder
placeholder

	if modelCfg.Type == "prompt_enhance" {
		enhancedPrompt, err := s.soraClient.EnhancePrompt(reqCtx, account, prompt, modelCfg.ExpansionLevel, modelCfg.DurationS)
		if err != nil {
			return nil, s.handleSoraRequestError(ctx, account, err, reqModel, c, clientStream)
	placeholder
		content := strings.TrimSpace(enhancedPrompt)
		if content == "" {
			content = prompt
	placeholder
		var firstTokenMs *int
		if clientStream {
			ms, streamErr := s.writeSoraStream(c, reqModel, content, startTime)
			if streamErr != nil {
				return nil, streamErr
		placeholder
			firstTokenMs = ms
	placeholder else if c != nil {
			c.JSON(http.StatusOK, buildSoraNonStreamResponse(content, reqModel))
	placeholder
		return &ForwardResult{
			RequestID:    "",
			Model:        reqModel,
			Stream:       clientStream,
			Duration:     time.Since(startTime),
			FirstTokenMs: firstTokenMs,
			Usage:        ClaudeUsage{placeholder,
			MediaType:    "prompt",
	placeholder, nil
placeholder

	var imageData []byte
	imageFilename := ""
	if strings.TrimSpace(imageInput) != "" {
		decoded, filename, err := decodeSoraImageInput(reqCtx, imageInput)
		if err != nil {
			s.writeSoraError(c, http.StatusBadRequest, "invalid_request_error", err.Error(), clientStream)
			return nil, err
	placeholder
		imageData = decoded
		imageFilename = filename
placeholder

	mediaID := ""
	if len(imageData) > 0 {
		uploadID, err := s.soraClient.UploadImage(reqCtx, account, imageData, imageFilename)
		if err != nil {
			return nil, s.handleSoraRequestError(ctx, account, err, reqModel, c, clientStream)
	placeholder
		mediaID = uploadID
placeholder

	taskID := ""
	var err error
	switch modelCfg.Type {
	case "image":
		taskID, err = s.soraClient.CreateImageTask(reqCtx, account, SoraImageRequest{
			Prompt:  prompt,
			Width:   modelCfg.Width,
			Height:  modelCfg.Height,
			MediaID: mediaID,
	placeholder)
	case "video":
		taskID, err = s.soraClient.CreateVideoTask(reqCtx, account, SoraVideoRequest{
			Prompt:        prompt,
			Orientation:   modelCfg.Orientation,
			Frames:        modelCfg.Frames,
			Model:         modelCfg.Model,
			Size:          modelCfg.Size,
			MediaID:       mediaID,
			RemixTargetID: remixTargetID,
	placeholder)
	default:
		err = fmt.Errorf("unsupported model type: %s", modelCfg.Type)
placeholder
	if err != nil {
		return nil, s.handleSoraRequestError(ctx, account, err, reqModel, c, clientStream)
placeholder

	if clientStream && c != nil {
		s.prepareSoraStream(c, taskID)
placeholder

	var mediaURLs []string
	mediaType := modelCfg.Type
	imageCount := 0
	imageSize := ""
	switch modelCfg.Type {
	case "image":
		urls, pollErr := s.pollImageTask(reqCtx, c, account, taskID, clientStream)
		if pollErr != nil {
			return nil, s.handleSoraRequestError(ctx, account, pollErr, reqModel, c, clientStream)
	placeholder
		mediaURLs = urls
		imageCount = len(urls)
		imageSize = soraImageSizeFromModel(reqModel)
	case "video":
		urls, pollErr := s.pollVideoTask(reqCtx, c, account, taskID, clientStream)
		if pollErr != nil {
			return nil, s.handleSoraRequestError(ctx, account, pollErr, reqModel, c, clientStream)
	placeholder
		mediaURLs = urls
	default:
		mediaType = "prompt"
placeholder

	finalURLs := s.normalizeSoraMediaURLs(mediaURLs)
	if len(mediaURLs) > 0 && s.mediaStorage != nil && s.mediaStorage.Enabled() {
		stored, storeErr := s.mediaStorage.StoreFromURLs(reqCtx, mediaType, mediaURLs)
		if storeErr != nil {
			// 存储失败时降级使用原始 URL，不中断用户请求
			log.Printf("[Sora] StoreFromURLs failed, falling back to original URLs: %v", storeErr)
	placeholder else {
			finalURLs = s.normalizeSoraMediaURLs(stored)
	placeholder
placeholder

	content := buildSoraContent(mediaType, finalURLs)
	var firstTokenMs *int
	if clientStream {
		ms, streamErr := s.writeSoraStream(c, reqModel, content, startTime)
		if streamErr != nil {
			return nil, streamErr
	placeholder
		firstTokenMs = ms
placeholder else if c != nil {
		response := buildSoraNonStreamResponse(content, reqModel)
		if len(finalURLs) > 0 {
			response["media_url"] = finalURLs[0]
			if len(finalURLs) > 1 {
				response["media_urls"] = finalURLs
		placeholder
	placeholder
		c.JSON(http.StatusOK, response)
placeholder

	return &ForwardResult{
		RequestID:    taskID,
		Model:        reqModel,
		Stream:       clientStream,
		Duration:     time.Since(startTime),
		FirstTokenMs: firstTokenMs,
		Usage:        ClaudeUsage{placeholder,
		MediaType:    mediaType,
		MediaURL:     firstMediaURL(finalURLs),
		ImageCount:   imageCount,
		ImageSize:    imageSize,
placeholder, nil
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

func (s *SoraGatewayService) shouldFailoverUpstreamError(statusCode int) bool {
	switch statusCode {
	case 401, 402, 403, 404, 429, 529:
		return true
	default:
		return statusCode >= 500
placeholder
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

func (s *SoraGatewayService) prepareSoraStream(c *gin.Context, requestID string) {
	if c == nil {
		return
placeholder
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	if strings.TrimSpace(requestID) != "" {
		c.Header("x-request-id", requestID)
placeholder
placeholder

func (s *SoraGatewayService) writeSoraStream(c *gin.Context, model, content string, startTime time.Time) (*int, error) {
	if c == nil {
		return nil, nil
placeholder
	writer := c.Writer
	flusher, _ := writer.(http.Flusher)

	chunk := map[string]any{
		"id":      fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano()),
		"object":  "chat.completion.chunk",
		"created": time.Now().Unix(),
		"model":   model,
		"choices": []any{
			map[string]any{
				"index": 0,
				"delta": map[string]any{
					"content": content,
			placeholder,
		placeholder,
	placeholder,
placeholder
	encoded, _ := json.Marshal(chunk)
	if _, err := fmt.Fprintf(writer, "data: %s\n\n", encoded); err != nil {
		return nil, err
placeholder
	if flusher != nil {
		flusher.Flush()
placeholder
	ms := int(time.Since(startTime).Milliseconds())
	finalChunk := map[string]any{
		"id":      chunk["id"],
		"object":  "chat.completion.chunk",
		"created": time.Now().Unix(),
		"model":   model,
		"choices": []any{
			map[string]any{
				"index":         0,
				"delta":         map[string]any{placeholder,
				"finish_reason": "stop",
		placeholder,
	placeholder,
placeholder
	finalEncoded, _ := json.Marshal(finalChunk)
	if _, err := fmt.Fprintf(writer, "data: %s\n\n", finalEncoded); err != nil {
		return &ms, err
placeholder
	if _, err := fmt.Fprint(writer, "data: [DONE]\n\n"); err != nil {
		return &ms, err
placeholder
	if flusher != nil {
		flusher.Flush()
placeholder
	return &ms, nil
placeholder

func (s *SoraGatewayService) writeSoraError(c *gin.Context, status int, errType, message string, stream bool) {
	if c == nil {
		return
placeholder
	if stream {
		flusher, _ := c.Writer.(http.Flusher)
		errorEvent := fmt.Sprintf(`event: error`+"\n"+`data: {"error": {"type": "%s", "message": "%s"placeholderplaceholder`+"\n\n", errType, message)
		_, _ = fmt.Fprint(c.Writer, errorEvent)
		_, _ = fmt.Fprint(c.Writer, "data: [DONE]\n\n")
		if flusher != nil {
			flusher.Flush()
	placeholder
		return
placeholder
	c.JSON(status, gin.H{
		"error": gin.H{
			"type":    errType,
			"message": message,
	placeholder,
placeholder)
placeholder

func (s *SoraGatewayService) handleSoraRequestError(ctx context.Context, account *Account, err error, model string, c *gin.Context, stream bool) error {
	if err == nil {
		return nil
placeholder
	var upstreamErr *SoraUpstreamError
	if errors.As(err, &upstreamErr) {
		if s.rateLimitService != nil && account != nil {
			s.rateLimitService.HandleUpstreamError(ctx, account, upstreamErr.StatusCode, upstreamErr.Headers, upstreamErr.Body)
	placeholder
		if s.shouldFailoverUpstreamError(upstreamErr.StatusCode) {
			return &UpstreamFailoverError{StatusCode: upstreamErr.StatusCode, ResponseBody: upstreamErr.Bodyplaceholder
	placeholder
		msg := upstreamErr.Message
		if override := soraProErrorMessage(model, msg); override != "" {
			msg = override
	placeholder
		s.writeSoraError(c, upstreamErr.StatusCode, "upstream_error", msg, stream)
		return err
placeholder
	if errors.Is(err, context.DeadlineExceeded) {
		s.writeSoraError(c, http.StatusGatewayTimeout, "timeout_error", "Sora generation timeout", stream)
		return err
placeholder
	s.writeSoraError(c, http.StatusBadGateway, "api_error", err.Error(), stream)
	return err
placeholder

func (s *SoraGatewayService) pollImageTask(ctx context.Context, c *gin.Context, account *Account, taskID string, stream bool) ([]string, error) {
	interval := s.pollInterval()
	maxAttempts := s.pollMaxAttempts()
	lastPing := time.Now()
	for attempt := 0; attempt < maxAttempts; attempt++ {
		status, err := s.soraClient.GetImageTask(ctx, account, taskID)
		if err != nil {
			return nil, err
	placeholder
		switch strings.ToLower(status.Status) {
		case "succeeded", "completed":
			return status.URLs, nil
		case "failed":
			if status.ErrorMsg != "" {
				return nil, errors.New(status.ErrorMsg)
		placeholder
			return nil, errors.New("sora image generation failed")
	placeholder
		if stream {
			s.maybeSendPing(c, &lastPing)
	placeholder
		if err := sleepWithContext(ctx, interval); err != nil {
			return nil, err
	placeholder
placeholder
	return nil, errors.New("sora image generation timeout")
placeholder

func (s *SoraGatewayService) pollVideoTask(ctx context.Context, c *gin.Context, account *Account, taskID string, stream bool) ([]string, error) {
	interval := s.pollInterval()
	maxAttempts := s.pollMaxAttempts()
	lastPing := time.Now()
	for attempt := 0; attempt < maxAttempts; attempt++ {
		status, err := s.soraClient.GetVideoTask(ctx, account, taskID)
		if err != nil {
			return nil, err
	placeholder
		switch strings.ToLower(status.Status) {
		case "completed", "succeeded":
			return status.URLs, nil
		case "failed":
			if status.ErrorMsg != "" {
				return nil, errors.New(status.ErrorMsg)
		placeholder
			return nil, errors.New("sora video generation failed")
	placeholder
		if stream {
			s.maybeSendPing(c, &lastPing)
	placeholder
		if err := sleepWithContext(ctx, interval); err != nil {
			return nil, err
	placeholder
placeholder
	return nil, errors.New("sora video generation timeout")
placeholder

func (s *SoraGatewayService) pollInterval() time.Duration {
	if s == nil || s.cfg == nil {
		return 2 * time.Second
placeholder
	interval := s.cfg.Sora.Client.PollIntervalSeconds
	if interval <= 0 {
		interval = 2
placeholder
	return time.Duration(interval) * time.Second
placeholder

func (s *SoraGatewayService) pollMaxAttempts() int {
	if s == nil || s.cfg == nil {
		return 600
placeholder
	maxAttempts := s.cfg.Sora.Client.MaxPollAttempts
	if maxAttempts <= 0 {
		maxAttempts = 600
placeholder
	return maxAttempts
placeholder

func (s *SoraGatewayService) maybeSendPing(c *gin.Context, lastPing *time.Time) {
	if c == nil {
		return
placeholder
	interval := 10 * time.Second
	if s != nil && s.cfg != nil && s.cfg.Concurrency.PingInterval > 0 {
		interval = time.Duration(s.cfg.Concurrency.PingInterval) * time.Second
placeholder
	if time.Since(*lastPing) < interval {
		return
placeholder
	if _, err := fmt.Fprint(c.Writer, ":\n\n"); err == nil {
		if flusher, ok := c.Writer.(http.Flusher); ok {
			flusher.Flush()
	placeholder
		*lastPing = time.Now()
placeholder
placeholder

func (s *SoraGatewayService) normalizeSoraMediaURLs(urls []string) []string {
	if len(urls) == 0 {
		return urls
placeholder
	output := make([]string, 0, len(urls))
	for _, raw := range urls {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
	placeholder
		if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
			output = append(output, raw)
			continue
	placeholder
		pathVal := raw
		if !strings.HasPrefix(pathVal, "/") {
			pathVal = "/" + pathVal
	placeholder
		output = append(output, s.buildSoraMediaURL(pathVal, ""))
placeholder
	return output
placeholder

func buildSoraContent(mediaType string, urls []string) string {
	switch mediaType {
	case "image":
		parts := make([]string, 0, len(urls))
		for _, u := range urls {
			parts = append(parts, fmt.Sprintf("![image](%s)", u))
	placeholder
		return strings.Join(parts, "\n")
	case "video":
		if len(urls) == 0 {
			return ""
	placeholder
		return fmt.Sprintf("```html\n<video src='%s' controls></video>\n```", urls[0])
	default:
		return ""
placeholder
placeholder

func extractSoraInput(body map[string]any) (prompt, imageInput, videoInput, remixTargetID string) {
	if body == nil {
		return "", "", "", ""
placeholder
	if v, ok := body["remix_target_id"].(string); ok {
		remixTargetID = v
placeholder
	if v, ok := body["image"].(string); ok {
		imageInput = v
placeholder
	if v, ok := body["video"].(string); ok {
		videoInput = v
placeholder
	if v, ok := body["prompt"].(string); ok && strings.TrimSpace(v) != "" {
		prompt = v
placeholder
	if messages, ok := body["messages"].([]any); ok {
		builder := strings.Builder{placeholder
		for _, raw := range messages {
			msg, ok := raw.(map[string]any)
			if !ok {
				continue
		placeholder
			role, _ := msg["role"].(string)
			if role != "" && role != "user" {
				continue
		placeholder
			content := msg["content"]
			text, img, vid := parseSoraMessageContent(content)
			if text != "" {
				if builder.Len() > 0 {
					_, _ = builder.WriteString("\n")
			placeholder
				_, _ = builder.WriteString(text)
		placeholder
			if imageInput == "" && img != "" {
				imageInput = img
		placeholder
			if videoInput == "" && vid != "" {
				videoInput = vid
		placeholder
	placeholder
		if prompt == "" {
			prompt = builder.String()
	placeholder
placeholder
	return prompt, imageInput, videoInput, remixTargetID
placeholder

func parseSoraMessageContent(content any) (text, imageInput, videoInput string) {
	switch val := content.(type) {
	case string:
		return val, "", ""
	case []any:
		builder := strings.Builder{placeholder
		for _, item := range val {
			itemMap, ok := item.(map[string]any)
			if !ok {
				continue
		placeholder
			t, _ := itemMap["type"].(string)
			switch t {
			case "text":
				if txt, ok := itemMap["text"].(string); ok && strings.TrimSpace(txt) != "" {
					if builder.Len() > 0 {
						_, _ = builder.WriteString("\n")
				placeholder
					_, _ = builder.WriteString(txt)
			placeholder
			case "image_url":
				if imageInput == "" {
					if urlVal, ok := itemMap["image_url"].(map[string]any); ok {
						imageInput = fmt.Sprintf("%v", urlVal["url"])
				placeholder else if urlStr, ok := itemMap["image_url"].(string); ok {
						imageInput = urlStr
				placeholder
			placeholder
			case "video_url":
				if videoInput == "" {
					if urlVal, ok := itemMap["video_url"].(map[string]any); ok {
						videoInput = fmt.Sprintf("%v", urlVal["url"])
				placeholder else if urlStr, ok := itemMap["video_url"].(string); ok {
						videoInput = urlStr
				placeholder
			placeholder
		placeholder
	placeholder
		return builder.String(), imageInput, videoInput
	default:
		return "", "", ""
placeholder
placeholder

func decodeSoraImageInput(ctx context.Context, input string) ([]byte, string, error) {
	raw := strings.TrimSpace(input)
	if raw == "" {
		return nil, "", errors.New("empty image input")
placeholder
	if strings.HasPrefix(raw, "data:") {
		parts := strings.SplitN(raw, ",", 2)
		if len(parts) != 2 {
			return nil, "", errors.New("invalid data url")
	placeholder
		meta := parts[0]
		payload := parts[1]
		decoded, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			return nil, "", err
	placeholder
		ext := ""
		if strings.HasPrefix(meta, "data:") {
			metaParts := strings.SplitN(meta[5:], ";", 2)
			if len(metaParts) > 0 {
				if exts, err := mime.ExtensionsByType(metaParts[0]); err == nil && len(exts) > 0 {
					ext = exts[0]
			placeholder
		placeholder
	placeholder
		filename := "image" + ext
		return decoded, filename, nil
placeholder
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return downloadSoraImageInput(ctx, raw)
placeholder
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, "", errors.New("invalid base64 image")
placeholder
	return decoded, "image.png", nil
placeholder

func downloadSoraImageInput(ctx context.Context, rawURL string) ([]byte, string, error) {
	parsed, err := validateSoraImageURL(rawURL)
	if err != nil {
		return nil, "", err
placeholder
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsed.String(), nil)
	if err != nil {
		return nil, "", err
placeholder
	client := &http.Client{
		Timeout: soraImageInputTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= soraImageInputMaxRedirects {
				return errors.New("too many redirects")
		placeholder
			return validateSoraImageURLValue(req.URL)
	placeholder,
placeholder
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("download image failed: %d", resp.StatusCode)
placeholder
	data, err := io.ReadAll(io.LimitReader(resp.Body, soraImageInputMaxBytes))
	if err != nil {
		return nil, "", err
placeholder
	ext := fileExtFromURL(parsed.String())
	if ext == "" {
		ext = fileExtFromContentType(resp.Header.Get("Content-Type"))
placeholder
	filename := "image" + ext
	return data, filename, nil
placeholder

func validateSoraImageURL(raw string) (*url.URL, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, errors.New("empty image url")
placeholder
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid image url: %w", err)
placeholder
	if err := validateSoraImageURLValue(parsed); err != nil {
		return nil, err
placeholder
	return parsed, nil
placeholder

func validateSoraImageURLValue(parsed *url.URL) error {
	if parsed == nil {
		return errors.New("invalid image url")
placeholder
	scheme := strings.ToLower(strings.TrimSpace(parsed.Scheme))
	if scheme != "http" && scheme != "https" {
		return errors.New("only http/https image url is allowed")
placeholder
	if parsed.User != nil {
		return errors.New("image url cannot contain userinfo")
placeholder
	host := strings.ToLower(strings.TrimSpace(parsed.Hostname()))
	if host == "" {
		return errors.New("image url missing host")
placeholder
	if _, blocked := soraBlockedHostnames[host]; blocked {
		return errors.New("image url is not allowed")
placeholder
	if ip := net.ParseIP(host); ip != nil {
		if isSoraBlockedIP(ip) {
			return errors.New("image url is not allowed")
	placeholder
		return nil
placeholder
	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("resolve image url failed: %w", err)
placeholder
	for _, ip := range ips {
		if isSoraBlockedIP(ip) {
			return errors.New("image url is not allowed")
	placeholder
placeholder
	return nil
placeholder

func isSoraBlockedIP(ip net.IP) bool {
	if ip == nil {
		return true
placeholder
	for _, cidr := range soraBlockedCIDRs {
		if cidr.Contains(ip) {
			return true
	placeholder
placeholder
	return false
placeholder

func mustParseCIDRs(values []string) []*net.IPNet {
	out := make([]*net.IPNet, 0, len(values))
	for _, val := range values {
		_, cidr, err := net.ParseCIDR(val)
		if err != nil {
			continue
	placeholder
		out = append(out, cidr)
placeholder
	return out
placeholder
