package service

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/sha3"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/Wei-Shaw/sub2api/internal/pkg/proxyurl"
	"github.com/Wei-Shaw/sub2api/internal/util/responseheaders"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	openAIImagesGenerationsEndpoint = "/v1/images/generations"
	openAIImagesEditsEndpoint       = "/v1/images/edits"

	openAIImagesGenerationsURL = "https://api.openai.com/v1/images/generations"
	openAIImagesEditsURL       = "https://api.openai.com/v1/images/edits"

	openAIChatGPTStartURL               = "https://chatgpt.com/"
	openAIChatGPTFilesURL               = "https://chatgpt.com/backend-api/files"
	openAIChatGPTConversationInitURL    = "https://chatgpt.com/backend-api/conversation/init"
	openAIChatGPTConversationURL        = "https://chatgpt.com/backend-api/f/conversation"
	openAIChatGPTConversationPrepareURL = "https://chatgpt.com/backend-api/f/conversation/prepare"
	openAIChatGPTChatRequirementsURL    = "https://chatgpt.com/backend-api/sentinel/chat-requirements"

	openAIImageBackendUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
	openAIImageRequirementsDiff = "0fffff"
)

type OpenAIImagesCapability string

const (
	OpenAIImagesCapabilityBasic  OpenAIImagesCapability = "images-basic"
	OpenAIImagesCapabilityNative OpenAIImagesCapability = "images-native"
)

type OpenAIImagesUpload struct {
	FieldName   string
	FileName    string
	ContentType string
	Data        []byte
	Width       int
	Height      int
placeholder

type OpenAIImagesRequest struct {
	Endpoint           string
	ContentType        string
	Multipart          bool
	Model              string
	ExplicitModel      bool
	Prompt             string
	Stream             bool
	N                  int
	Size               string
	ExplicitSize       bool
	SizeTier           string
	ResponseFormat     string
	HasMask            bool
	HasNativeOptions   bool
	RequiredCapability OpenAIImagesCapability
	Uploads            []OpenAIImagesUpload
	Body               []byte
	bodyHash           string
placeholder

func (r *OpenAIImagesRequest) IsEdits() bool {
	return r != nil && r.Endpoint == openAIImagesEditsEndpoint
placeholder

func (r *OpenAIImagesRequest) StickySessionSeed() string {
	if r == nil {
		return ""
placeholder
	parts := []string{
		"openai-images",
		strings.TrimSpace(r.Endpoint),
		strings.TrimSpace(r.Model),
		strings.TrimSpace(r.Size),
		strings.TrimSpace(r.Prompt),
placeholder
	seed := strings.Join(parts, "|")
	if strings.TrimSpace(r.Prompt) == "" && r.bodyHash != "" {
		seed += "|body=" + r.bodyHash
placeholder
	return seed
placeholder

func (s *OpenAIGatewayService) ParseOpenAIImagesRequest(c *gin.Context, body []byte) (*OpenAIImagesRequest, error) {
	if c == nil || c.Request == nil {
		return nil, fmt.Errorf("missing request context")
placeholder
	endpoint := normalizeOpenAIImagesEndpointPath(c.Request.URL.Path)
	if endpoint == "" {
		return nil, fmt.Errorf("unsupported images endpoint")
placeholder

	contentType := strings.TrimSpace(c.GetHeader("Content-Type"))
	req := &OpenAIImagesRequest{
		Endpoint:    endpoint,
		ContentType: contentType,
		N:           1,
		Body:        body,
placeholder
	if len(body) > 0 {
		sum := sha256.Sum256(body)
		req.bodyHash = hex.EncodeToString(sum[:8])
placeholder

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err == nil && strings.EqualFold(mediaType, "multipart/form-data") {
		req.Multipart = true
		if parseErr := parseOpenAIImagesMultipartRequest(body, contentType, req); parseErr != nil {
			return nil, parseErr
	placeholder
placeholder else {
		if len(body) == 0 {
			return nil, fmt.Errorf("request body is empty")
	placeholder
		if !gjson.ValidBytes(body) {
			return nil, fmt.Errorf("failed to parse request body")
	placeholder
		if parseErr := parseOpenAIImagesJSONRequest(body, req); parseErr != nil {
			return nil, parseErr
	placeholder
placeholder

	applyOpenAIImagesDefaults(req)
	req.SizeTier = normalizeOpenAIImageSizeTier(req.Size)
	req.RequiredCapability = classifyOpenAIImagesCapability(req)
	return req, nil
placeholder

func parseOpenAIImagesJSONRequest(body []byte, req *OpenAIImagesRequest) error {
	if modelResult := gjson.GetBytes(body, "model"); modelResult.Exists() {
		req.Model = strings.TrimSpace(modelResult.String())
		req.ExplicitModel = req.Model != ""
placeholder
	req.Prompt = strings.TrimSpace(gjson.GetBytes(body, "prompt").String())

	if streamResult := gjson.GetBytes(body, "stream"); streamResult.Exists() {
		if streamResult.Type != gjson.True && streamResult.Type != gjson.False {
			return fmt.Errorf("invalid stream field type")
	placeholder
		req.Stream = streamResult.Bool()
placeholder

	if nResult := gjson.GetBytes(body, "n"); nResult.Exists() {
		if nResult.Type != gjson.Number {
			return fmt.Errorf("invalid n field type")
	placeholder
		req.N = int(nResult.Int())
		if req.N <= 0 {
			return fmt.Errorf("n must be greater than 0")
	placeholder
placeholder

	if sizeResult := gjson.GetBytes(body, "size"); sizeResult.Exists() {
		req.Size = strings.TrimSpace(sizeResult.String())
		req.ExplicitSize = req.Size != ""
placeholder
	req.ResponseFormat = strings.ToLower(strings.TrimSpace(gjson.GetBytes(body, "response_format").String()))
	req.HasMask = gjson.GetBytes(body, "mask").Exists()
	req.HasNativeOptions = hasOpenAINativeImageOptions(func(path string) bool {
		return gjson.GetBytes(body, path).Exists()
placeholder)
	return nil
placeholder

func parseOpenAIImagesMultipartRequest(body []byte, contentType string, req *OpenAIImagesRequest) error {
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return fmt.Errorf("invalid multipart content-type: %w", err)
placeholder
	boundary := strings.TrimSpace(params["boundary"])
	if boundary == "" {
		return fmt.Errorf("multipart boundary is required")
placeholder

	reader := multipart.NewReader(bytes.NewReader(body), boundary)
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
	placeholder
		if err != nil {
			return fmt.Errorf("read multipart body: %w", err)
	placeholder
		name := strings.TrimSpace(part.FormName())
		if name == "" {
			_ = part.Close()
			continue
	placeholder

		data, err := io.ReadAll(part)
		_ = part.Close()
		if err != nil {
			return fmt.Errorf("read multipart field %s: %w", name, err)
	placeholder

		fileName := strings.TrimSpace(part.FileName())
		if fileName != "" {
			partContentType := strings.TrimSpace(part.Header.Get("Content-Type"))
			if name == "mask" && len(data) > 0 {
				req.HasMask = true
		placeholder
			if name == "image" || strings.HasPrefix(name, "image[") {
				width, height := parseOpenAIImageDimensions(part.Header)
				req.Uploads = append(req.Uploads, OpenAIImagesUpload{
					FieldName:   name,
					FileName:    fileName,
					ContentType: partContentType,
					Data:        data,
					Width:       width,
					Height:      height,
			placeholder)
		placeholder
			continue
	placeholder

		value := strings.TrimSpace(string(data))
		switch name {
		case "model":
			req.Model = value
			req.ExplicitModel = value != ""
		case "prompt":
			req.Prompt = value
		case "size":
			req.Size = value
			req.ExplicitSize = value != ""
		case "response_format":
			req.ResponseFormat = strings.ToLower(value)
		case "stream":
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("invalid stream field value")
		placeholder
			req.Stream = parsed
		case "n":
			n, err := strconv.Atoi(value)
			if err != nil || n <= 0 {
				return fmt.Errorf("n must be a positive integer")
		placeholder
			req.N = n
		default:
			if isOpenAINativeImageOption(name) && value != "" {
				req.HasNativeOptions = true
		placeholder
	placeholder
placeholder

	if len(req.Uploads) == 0 && req.IsEdits() {
		return fmt.Errorf("image file is required")
placeholder
	return nil
placeholder

func parseOpenAIImageDimensions(_ textproto.MIMEHeader) (int, int) {
	return 0, 0
placeholder

func applyOpenAIImagesDefaults(req *OpenAIImagesRequest) {
	if req == nil {
		return
placeholder
	if req.N <= 0 {
		req.N = 1
placeholder
	if strings.TrimSpace(req.Model) != "" {
		req.Model = strings.TrimSpace(req.Model)
		return
placeholder
	req.Model = "gpt-image-2"
placeholder

func normalizeOpenAIImagesEndpointPath(path string) string {
	trimmed := strings.TrimSpace(path)
	switch {
	case strings.Contains(trimmed, "/images/generations"):
		return openAIImagesGenerationsEndpoint
	case strings.Contains(trimmed, "/images/edits"):
		return openAIImagesEditsEndpoint
	default:
		return ""
placeholder
placeholder

func classifyOpenAIImagesCapability(req *OpenAIImagesRequest) OpenAIImagesCapability {
	if req == nil {
		return OpenAIImagesCapabilityNative
placeholder
	if req.ExplicitModel || req.ExplicitSize {
		return OpenAIImagesCapabilityNative
placeholder
	model := strings.ToLower(strings.TrimSpace(req.Model))
	if !strings.HasPrefix(model, "gpt-image-") {
		return OpenAIImagesCapabilityNative
placeholder
	if req.Stream || req.N != 1 || req.HasMask || req.HasNativeOptions {
		return OpenAIImagesCapabilityNative
placeholder
	if req.IsEdits() && !req.Multipart {
		return OpenAIImagesCapabilityNative
placeholder
	if req.ResponseFormat != "" && req.ResponseFormat != "b64_json" {
		return OpenAIImagesCapabilityNative
placeholder
	return OpenAIImagesCapabilityBasic
placeholder

func hasOpenAINativeImageOptions(exists func(path string) bool) bool {
	for _, path := range []string{
		"background",
		"quality",
		"style",
		"output_format",
		"output_compression",
		"moderation",
placeholder {
		if exists(path) {
			return true
	placeholder
placeholder
	return false
placeholder

func isOpenAINativeImageOption(name string) bool {
	switch strings.TrimSpace(strings.ToLower(name)) {
	case "background", "quality", "style", "output_format", "output_compression", "moderation":
		return true
	default:
		return false
placeholder
placeholder

func normalizeOpenAIImageSizeTier(size string) string {
	switch strings.ToLower(strings.TrimSpace(size)) {
	case "1024x1024":
		return "1K"
	case "1536x1024", "1024x1536", "1792x1024", "1024x1792", "", "auto":
		return "2K"
	default:
		return "2K"
placeholder
placeholder

func (s *OpenAIGatewayService) ForwardImages(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	body []byte,
	parsed *OpenAIImagesRequest,
	channelMappedModel string,
) (*OpenAIForwardResult, error) {
	if parsed == nil {
		return nil, fmt.Errorf("parsed images request is required")
placeholder
	switch account.Type {
	case AccountTypeAPIKey:
		return s.forwardOpenAIImagesAPIKey(ctx, c, account, body, parsed, channelMappedModel)
	case AccountTypeOAuth:
		return s.forwardOpenAIImagesOAuth(ctx, c, account, parsed, channelMappedModel)
	default:
		return nil, fmt.Errorf("unsupported account type: %s", account.Type)
placeholder
placeholder

func (s *OpenAIGatewayService) forwardOpenAIImagesAPIKey(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	body []byte,
	parsed *OpenAIImagesRequest,
	channelMappedModel string,
) (*OpenAIForwardResult, error) {
	startTime := time.Now()
	requestModel := strings.TrimSpace(parsed.Model)
	if mapped := strings.TrimSpace(channelMappedModel); mapped != "" {
		requestModel = mapped
placeholder
	upstreamModel := account.GetMappedModel(requestModel)
	forwardBody, forwardContentType, err := rewriteOpenAIImagesModel(body, parsed.ContentType, upstreamModel)
	if err != nil {
		return nil, err
placeholder
	if !parsed.Multipart {
		setOpsUpstreamRequestBody(c, forwardBody)
placeholder

	token, _, err := s.GetAccessToken(ctx, account)
	if err != nil {
		return nil, err
placeholder
	upstreamReq, err := s.buildOpenAIImagesRequest(ctx, c, account, forwardBody, forwardContentType, token, parsed.Endpoint)
	if err != nil {
		return nil, err
placeholder

	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
placeholder
	upstreamStart := time.Now()
	resp, err := s.httpUpstream.Do(upstreamReq, proxyURL, account.ID, account.Concurrency)
	SetOpsLatencyMs(c, OpsUpstreamLatencyMsKey, time.Since(upstreamStart).Milliseconds())
	if err != nil {
		safeErr := sanitizeUpstreamErrorMessage(err.Error())
		setOpsUpstreamError(c, 0, safeErr, "")
		appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
			Platform:           account.Platform,
			AccountID:          account.ID,
			AccountName:        account.Name,
			UpstreamStatusCode: 0,
			UpstreamURL:        safeUpstreamURL(upstreamReq.URL.String()),
			Kind:               "request_error",
			Message:            safeErr,
	placeholder)
		return nil, fmt.Errorf("upstream request failed: %s", safeErr)
placeholder
	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
		_ = resp.Body.Close()
		resp.Body = io.NopCloser(bytes.NewReader(respBody))
		upstreamMsg := strings.TrimSpace(extractUpstreamErrorMessage(respBody))
		upstreamMsg = sanitizeUpstreamErrorMessage(upstreamMsg)
		if s.shouldFailoverOpenAIUpstreamResponse(resp.StatusCode, upstreamMsg, respBody) {
			appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
				Platform:           account.Platform,
				AccountID:          account.ID,
				AccountName:        account.Name,
				UpstreamStatusCode: resp.StatusCode,
				UpstreamRequestID:  resp.Header.Get("x-request-id"),
				UpstreamURL:        safeUpstreamURL(upstreamReq.URL.String()),
				Kind:               "failover",
				Message:            upstreamMsg,
		placeholder)
			s.handleFailoverSideEffects(ctx, resp, account)
			return nil, &UpstreamFailoverError{
				StatusCode:             resp.StatusCode,
				ResponseBody:           respBody,
				RetryableOnSameAccount: account.IsPoolMode() && isPoolModeRetryableStatus(resp.StatusCode),
		placeholder
	placeholder
		return s.handleErrorResponse(ctx, resp, c, account, forwardBody)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	var usage OpenAIUsage
	imageCount := parsed.N
	var firstTokenMs *int
	if parsed.Stream {
		streamUsage, streamCount, ttft, err := s.handleOpenAIImagesStreamingResponse(resp, c, startTime)
		if err != nil {
			return nil, err
	placeholder
		usage = streamUsage
		imageCount = streamCount
		firstTokenMs = ttft
placeholder else {
		nonStreamUsage, nonStreamCount, err := s.handleOpenAIImagesNonStreamingResponse(resp, c)
		if err != nil {
			return nil, err
	placeholder
		usage = nonStreamUsage
		if nonStreamCount > 0 {
			imageCount = nonStreamCount
	placeholder
placeholder
	return &OpenAIForwardResult{
		RequestID:       resp.Header.Get("x-request-id"),
		Usage:           usage,
		Model:           requestModel,
		UpstreamModel:   upstreamModel,
		Stream:          parsed.Stream,
		ResponseHeaders: resp.Header.Clone(),
		Duration:        time.Since(startTime),
		FirstTokenMs:    firstTokenMs,
		ImageCount:      imageCount,
		ImageSize:       parsed.SizeTier,
placeholder, nil
placeholder

func (s *OpenAIGatewayService) buildOpenAIImagesRequest(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	body []byte,
	contentType string,
	token string,
	endpoint string,
) (*http.Request, error) {
	targetURL := openAIImagesGenerationsURL
	if endpoint == openAIImagesEditsEndpoint {
		targetURL = openAIImagesEditsURL
placeholder
	baseURL := account.GetOpenAIBaseURL()
	if baseURL != "" {
		validatedURL, err := s.validateUpstreamBaseURL(baseURL)
		if err != nil {
			return nil, err
	placeholder
		targetURL = buildOpenAIImagesURL(validatedURL, endpoint)
placeholder

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
placeholder
	req.Header.Set("Authorization", "Bearer "+token)
	for key, values := range c.Request.Header {
		if !openaiPassthroughAllowedHeaders[strings.ToLower(key)] {
			continue
	placeholder
		for _, value := range values {
			req.Header.Add(key, value)
	placeholder
placeholder
	customUA := account.GetOpenAIUserAgent()
	if customUA != "" {
		req.Header.Set("User-Agent", customUA)
placeholder
	if strings.TrimSpace(contentType) != "" {
		req.Header.Set("Content-Type", contentType)
placeholder
	return req, nil
placeholder

func buildOpenAIImagesURL(base string, endpoint string) string {
	normalized := strings.TrimRight(strings.TrimSpace(base), "/")
	relative := strings.TrimPrefix(strings.TrimSpace(endpoint), "/v1")
	if strings.HasSuffix(normalized, endpoint) || strings.HasSuffix(normalized, relative) {
		return normalized
placeholder
	if strings.HasSuffix(normalized, "/v1") {
		return normalized + relative
placeholder
	return normalized + endpoint
placeholder

func rewriteOpenAIImagesModel(body []byte, contentType string, model string) ([]byte, string, error) {
	model = strings.TrimSpace(model)
	if model == "" {
		return body, contentType, nil
placeholder
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err == nil && strings.EqualFold(mediaType, "multipart/form-data") {
		rewrittenBody, rewrittenType, rewriteErr := rewriteOpenAIImagesMultipartModel(body, contentType, model)
		return rewrittenBody, rewrittenType, rewriteErr
placeholder
	rewritten, err := sjson.SetBytes(body, "model", model)
	if err != nil {
		return nil, "", fmt.Errorf("rewrite image request model: %w", err)
placeholder
	return rewritten, contentType, nil
placeholder

func rewriteOpenAIImagesMultipartModel(body []byte, contentType string, model string) ([]byte, string, error) {
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, "", fmt.Errorf("parse multipart content-type: %w", err)
placeholder
	boundary := strings.TrimSpace(params["boundary"])
	if boundary == "" {
		return nil, "", fmt.Errorf("multipart boundary is required")
placeholder

	reader := multipart.NewReader(bytes.NewReader(body), boundary)
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	modelWritten := false

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
	placeholder
		if err != nil {
			return nil, "", fmt.Errorf("read multipart body: %w", err)
	placeholder

		formName := strings.TrimSpace(part.FormName())
		partHeader := cloneMultipartHeader(part.Header)
		target, err := writer.CreatePart(partHeader)
		if err != nil {
			_ = part.Close()
			return nil, "", fmt.Errorf("create multipart part: %w", err)
	placeholder

		if formName == "model" && part.FileName() == "" {
			if _, err := target.Write([]byte(model)); err != nil {
				_ = part.Close()
				return nil, "", fmt.Errorf("rewrite multipart model: %w", err)
		placeholder
			modelWritten = true
			_ = part.Close()
			continue
	placeholder
		if _, err := io.Copy(target, part); err != nil {
			_ = part.Close()
			return nil, "", fmt.Errorf("copy multipart part: %w", err)
	placeholder
		_ = part.Close()
placeholder

	if !modelWritten {
		if err := writer.WriteField("model", model); err != nil {
			return nil, "", fmt.Errorf("append multipart model field: %w", err)
	placeholder
placeholder
	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("finalize multipart body: %w", err)
placeholder
	return buffer.Bytes(), writer.FormDataContentType(), nil
placeholder

func cloneMultipartHeader(src textproto.MIMEHeader) textproto.MIMEHeader {
	dst := make(textproto.MIMEHeader, len(src))
	for key, values := range src {
		copied := make([]string, len(values))
		copy(copied, values)
		dst[key] = copied
placeholder
	return dst
placeholder

func (s *OpenAIGatewayService) handleOpenAIImagesNonStreamingResponse(resp *http.Response, c *gin.Context) (OpenAIUsage, int, error) {
	body, err := ReadUpstreamResponseBody(resp.Body, s.cfg, c, openAITooLargeError)
	if err != nil {
		return OpenAIUsage{placeholder, 0, err
placeholder
	responseheaders.WriteFilteredHeaders(c.Writer.Header(), resp.Header, s.responseHeaderFilter)
	contentType := "application/json"
	if s.cfg != nil && !s.cfg.Security.ResponseHeaders.Enabled {
		if upstreamType := resp.Header.Get("Content-Type"); upstreamType != "" {
			contentType = upstreamType
	placeholder
placeholder
	c.Data(resp.StatusCode, contentType, body)

	usage, _ := extractOpenAIUsageFromJSONBytes(body)
	return usage, extractOpenAIImageCountFromJSONBytes(body), nil
placeholder

func (s *OpenAIGatewayService) handleOpenAIImagesStreamingResponse(
	resp *http.Response,
	c *gin.Context,
	startTime time.Time,
) (OpenAIUsage, int, *int, error) {
	responseheaders.WriteFilteredHeaders(c.Writer.Header(), resp.Header, s.responseHeaderFilter)
	contentType := strings.TrimSpace(resp.Header.Get("Content-Type"))
	if contentType == "" {
		contentType = "text/event-stream"
placeholder
	c.Status(resp.StatusCode)
	c.Header("Content-Type", contentType)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return OpenAIUsage{placeholder, 0, nil, fmt.Errorf("streaming is not supported by response writer")
placeholder

	reader := bufio.NewReader(resp.Body)
	usage := OpenAIUsage{placeholder
	imageCount := 0
	var firstTokenMs *int

	for {
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			if firstTokenMs == nil {
				ms := int(time.Since(startTime).Milliseconds())
				firstTokenMs = &ms
		placeholder
			if _, writeErr := c.Writer.Write(line); writeErr != nil {
				return OpenAIUsage{placeholder, 0, firstTokenMs, writeErr
		placeholder
			flusher.Flush()

			if data, ok := extractOpenAISSEDataLine(strings.TrimRight(string(line), "\r\n")); ok && data != "" && data != "[DONE]" {
				dataBytes := []byte(data)
				mergeOpenAIUsage(&usage, dataBytes)
				if count := extractOpenAIImageCountFromJSONBytes(dataBytes); count > imageCount {
					imageCount = count
			placeholder
		placeholder
	placeholder
		if err == io.EOF {
			break
	placeholder
		if err != nil {
			return OpenAIUsage{placeholder, 0, firstTokenMs, err
	placeholder
placeholder
	return usage, imageCount, firstTokenMs, nil
placeholder

func mergeOpenAIUsage(dst *OpenAIUsage, body []byte) {
	if dst == nil {
		return
placeholder
	if parsed, ok := extractOpenAIUsageFromJSONBytes(body); ok {
		if parsed.InputTokens > 0 {
			dst.InputTokens = parsed.InputTokens
	placeholder
		if parsed.OutputTokens > 0 {
			dst.OutputTokens = parsed.OutputTokens
	placeholder
		if parsed.CacheReadInputTokens > 0 {
			dst.CacheReadInputTokens = parsed.CacheReadInputTokens
	placeholder
		if parsed.ImageOutputTokens > 0 {
			dst.ImageOutputTokens = parsed.ImageOutputTokens
	placeholder
placeholder
placeholder

func extractOpenAIImageCountFromJSONBytes(body []byte) int {
	if len(body) == 0 || !gjson.ValidBytes(body) {
		return 0
placeholder
	data := gjson.GetBytes(body, "data")
	if data.Exists() && data.IsArray() {
		return len(data.Array())
placeholder
	return 0
placeholder

func (s *OpenAIGatewayService) forwardOpenAIImagesOAuth(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	parsed *OpenAIImagesRequest,
	channelMappedModel string,
) (*OpenAIForwardResult, error) {
	startTime := time.Now()
	requestModel := strings.TrimSpace(parsed.Model)
	if mapped := strings.TrimSpace(channelMappedModel); mapped != "" {
		requestModel = mapped
placeholder

	token, _, err := s.GetAccessToken(ctx, account)
	if err != nil {
		return nil, err
placeholder
	client, err := newOpenAIBackendAPIClient(resolveOpenAIProxyURL(account))
	if err != nil {
		return nil, err
placeholder
	headers, err := s.buildOpenAIBackendAPIHeaders(account, token)
	if err != nil {
		return nil, err
placeholder
	if bootstrapErr := bootstrapOpenAIBackendAPI(ctx, client, headers); bootstrapErr != nil {
		logger.LegacyPrintf("service.openai_gateway", "OpenAI image bootstrap failed: %v", bootstrapErr)
placeholder

	chatReqs, err := fetchOpenAIChatRequirements(ctx, client, headers)
	if err != nil {
		return nil, s.wrapOpenAIImageBackendError(ctx, c, account, err)
placeholder
	if chatReqs.Arkose.Required {
		return nil, s.wrapOpenAIImageBackendError(
			ctx,
			c,
			account,
			newOpenAIImageSyntheticStatusError(
				http.StatusForbidden,
				"chat-requirements requires unsupported challenge (arkose)",
				openAIChatGPTChatRequirementsURL,
			),
		)
placeholder

	parentMessageID := uuid.NewString()
	proofToken := generateOpenAIProofToken(chatReqs.ProofOfWork.Required, chatReqs.ProofOfWork.Seed, chatReqs.ProofOfWork.Difficulty, headers.Get("User-Agent"))
	_ = initializeOpenAIImageConversation(ctx, client, headers)
	conduitToken, err := prepareOpenAIImageConversation(ctx, client, headers, parsed.Prompt, parentMessageID, chatReqs.Token, proofToken)
	if err != nil {
		return nil, s.wrapOpenAIImageBackendError(ctx, c, account, err)
placeholder

	uploads, err := uploadOpenAIImageFiles(ctx, client, headers, parsed.Uploads)
	if err != nil {
		return nil, s.wrapOpenAIImageBackendError(ctx, c, account, err)
placeholder

	convReq := buildOpenAIImageConversationRequest(parsed, parentMessageID, uploads)
	if parsedContent, err := json.Marshal(convReq); err == nil {
		setOpsUpstreamRequestBody(c, parsedContent)
placeholder
	convHeaders := cloneHTTPHeader(headers)
	convHeaders.Set("Accept", "text/event-stream")
	convHeaders.Set("Content-Type", "application/json")
	convHeaders.Set("openai-sentinel-chat-requirements-token", chatReqs.Token)
	if conduitToken != "" {
		convHeaders.Set("x-conduit-token", conduitToken)
placeholder
	if proofToken != "" {
		convHeaders.Set("openai-sentinel-proof-token", proofToken)
placeholder

	resp, err := client.R().
		SetContext(ctx).
		DisableAutoReadResponse().
		SetHeaders(headerToMap(convHeaders)).
		SetBodyJsonMarshal(convReq).
		Post(openAIChatGPTConversationURL)
	if err != nil {
		return nil, fmt.Errorf("openai image conversation request failed: %w", err)
placeholder
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
	placeholder
placeholder()
	if resp.StatusCode >= 400 {
		return nil, s.wrapOpenAIImageBackendError(ctx, c, account, handleOpenAIImageBackendError(resp))
placeholder

	conversationID, pointerInfos, usage, firstTokenMs, err := readOpenAIImageConversationStream(resp, startTime)
	if err != nil {
		return nil, err
placeholder
	pointerInfos = mergeOpenAIImagePointerInfos(pointerInfos, nil)
	if conversationID != "" && !hasOpenAIFileServicePointerInfos(pointerInfos) {
		polledPointers, pollErr := pollOpenAIImageConversation(ctx, client, headers, conversationID)
		if pollErr != nil {
			return nil, s.wrapOpenAIImageBackendError(ctx, c, account, pollErr)
	placeholder
		pointerInfos = mergeOpenAIImagePointerInfos(pointerInfos, polledPointers)
placeholder
	pointerInfos = preferOpenAIFileServicePointerInfos(pointerInfos)
	if len(pointerInfos) == 0 {
		return nil, fmt.Errorf("openai image conversation returned no downloadable images")
placeholder

	responseBody, imageCount, err := buildOpenAIImageResponse(ctx, client, headers, conversationID, pointerInfos)
	if err != nil {
		return nil, s.wrapOpenAIImageBackendError(ctx, c, account, err)
placeholder

	c.Data(http.StatusOK, "application/json; charset=utf-8", responseBody)
	return &OpenAIForwardResult{
		RequestID:     resp.Header.Get("x-request-id"),
		Usage:         usage,
		Model:         requestModel,
		UpstreamModel: requestModel,
		Stream:        false,
		Duration:      time.Since(startTime),
		FirstTokenMs:  firstTokenMs,
		ImageCount:    imageCount,
		ImageSize:     parsed.SizeTier,
placeholder, nil
placeholder

func resolveOpenAIProxyURL(account *Account) string {
	if account != nil && account.ProxyID != nil && account.Proxy != nil {
		return account.Proxy.URL()
placeholder
	return ""
placeholder

func newOpenAIBackendAPIClient(proxyURL string) (*req.Client, error) {
	client := req.C().
		SetTimeout(180 * time.Second).
		ImpersonateChrome()
	trimmed, _, err := proxyurl.Parse(proxyURL)
	if err != nil {
		return nil, err
placeholder
	if trimmed != "" {
		client.SetProxyURL(trimmed)
placeholder
	return client, nil
placeholder

func (s *OpenAIGatewayService) buildOpenAIBackendAPIHeaders(account *Account, token string) (http.Header, error) {
	deviceID, sessionID := s.ensureOpenAIImageSessionCredentials(context.Background(), account)
	headers := make(http.Header)
	headers.Set("Authorization", "Bearer "+token)
	headers.Set("Accept", "application/json")
	headers.Set("Origin", "https://chatgpt.com")
	headers.Set("Referer", "https://chatgpt.com/")
	headers.Set("Sec-Fetch-Dest", "empty")
	headers.Set("Sec-Fetch-Mode", "cors")
	headers.Set("Sec-Fetch-Site", "same-origin")
	headers.Set("User-Agent", openAIImageBackendUserAgent)
	if customUA := strings.TrimSpace(account.GetOpenAIUserAgent()); customUA != "" {
		headers.Set("User-Agent", customUA)
placeholder
	if chatgptAccountID := strings.TrimSpace(account.GetChatGPTAccountID()); chatgptAccountID != "" {
		headers.Set("chatgpt-account-id", chatgptAccountID)
placeholder
	if deviceID != "" {
		headers.Set("oai-device-id", deviceID)
		headers.Set("Cookie", "oai-did="+deviceID)
placeholder
	if sessionID != "" {
		headers.Set("oai-session-id", sessionID)
placeholder
	return headers, nil
placeholder

func (s *OpenAIGatewayService) ensureOpenAIImageSessionCredentials(ctx context.Context, account *Account) (string, string) {
	if account == nil {
		return "", ""
placeholder
	deviceID := account.GetOpenAIDeviceID()
	sessionID := account.GetOpenAISessionID()
	if deviceID != "" && sessionID != "" {
		return deviceID, sessionID
placeholder

	updates := map[string]any{placeholder
	if deviceID == "" {
		deviceID = uuid.NewString()
		updates["openai_device_id"] = deviceID
placeholder
	if sessionID == "" {
		sessionID = uuid.NewString()
		updates["openai_session_id"] = sessionID
placeholder
	if account.Extra == nil {
		account.Extra = map[string]any{placeholder
placeholder
	for key, value := range updates {
		account.Extra[key] = value
placeholder
	if len(updates) == 0 || s == nil || s.accountRepo == nil {
		return deviceID, sessionID
placeholder

	updateCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := s.accountRepo.UpdateExtra(updateCtx, account.ID, updates); err != nil {
		logger.LegacyPrintf("service.openai_gateway", "persist openai image session creds failed: account=%d err=%v", account.ID, err)
placeholder
	return deviceID, sessionID
placeholder

func bootstrapOpenAIBackendAPI(ctx context.Context, client *req.Client, headers http.Header) error {
	resp, err := client.R().
		SetContext(ctx).
		DisableAutoReadResponse().
		SetHeaders(headerToMap(headers)).
		Get(openAIChatGPTStartURL)
	if err != nil {
		return err
placeholder
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
placeholder
	return nil
placeholder

func initializeOpenAIImageConversation(ctx context.Context, client *req.Client, headers http.Header) error {
	payload := map[string]any{
		"gizmo_id":                nil,
		"requested_default_model": nil,
		"conversation_id":         nil,
		"timezone_offset_min":     openAITimezoneOffsetMinutes(),
		"system_hints":            []string{"picture_v2"placeholder,
placeholder
	resp, err := client.R().
		SetContext(ctx).
		SetHeaders(headerToMap(headers)).
		SetBodyJsonMarshal(payload).
		Post(openAIChatGPTConversationInitURL)
	if err != nil {
		return err
placeholder
	if !resp.IsSuccessState() {
		return newOpenAIImageStatusError(resp, "conversation init failed")
placeholder
	return nil
placeholder

type openAIChatRequirements struct {
	Token     string `json:"token"`
	Turnstile struct {
		Required bool `json:"required"`
placeholder `json:"turnstile"`
	Arkose struct {
		Required bool `json:"required"`
placeholder `json:"arkose"`
	ProofOfWork struct {
		Required   bool   `json:"required"`
		Seed       string `json:"seed"`
		Difficulty string `json:"difficulty"`
placeholder `json:"proofofwork"`
placeholder

func fetchOpenAIChatRequirements(ctx context.Context, client *req.Client, headers http.Header) (*openAIChatRequirements, error) {
	var lastErr error
	for _, payload := range []map[string]any{
		{"p": nilplaceholder,
		{"p": generateOpenAIRequirementsToken(headers.Get("User-Agent"))placeholder,
placeholder {
		var result openAIChatRequirements
		resp, err := client.R().
			SetContext(ctx).
			SetHeaders(headerToMap(headers)).
			SetBodyJsonMarshal(payload).
			SetSuccessResult(&result).
			Post(openAIChatGPTChatRequirementsURL)
		if err != nil {
			lastErr = err
			continue
	placeholder
		if resp.IsSuccessState() && strings.TrimSpace(result.Token) != "" {
			return &result, nil
	placeholder
		lastErr = newOpenAIImageStatusError(resp, "chat-requirements failed")
placeholder
	if lastErr == nil {
		lastErr = fmt.Errorf("chat-requirements failed")
placeholder
	return nil, lastErr
placeholder

func prepareOpenAIImageConversation(
	ctx context.Context,
	client *req.Client,
	headers http.Header,
	prompt string,
	parentMessageID string,
	chatToken string,
	proofToken string,
) (string, error) {
	messageID := uuid.NewString()
	payload := map[string]any{
		"action":                "next",
		"client_prepare_state":  "success",
		"fork_from_shared_post": false,
		"parent_message_id":     parentMessageID,
		"model":                 "auto",
		"timezone_offset_min":   openAITimezoneOffsetMinutes(),
		"timezone":              openAITimezoneName(),
		"conversation_mode":     map[string]any{"kind": "primary_assistant"placeholder,
		"system_hints":          []string{"picture_v2"placeholder,
		"supports_buffering":    true,
		"supported_encodings":   []string{"v1"placeholder,
		"partial_query": map[string]any{
			"id":     messageID,
			"author": map[string]any{"role": "user"placeholder,
			"content": map[string]any{
				"content_type": "text",
				"parts":        []string{coalesceOpenAIFileName(prompt, "Generate an image.")placeholder,
		placeholder,
	placeholder,
		"client_contextual_info": map[string]any{
			"app_name": "chatgpt.com",
	placeholder,
placeholder
	prepareHeaders := cloneHTTPHeader(headers)
	prepareHeaders.Set("Accept", "*/*")
	prepareHeaders.Set("Content-Type", "application/json")
	if strings.TrimSpace(chatToken) != "" {
		prepareHeaders.Set("openai-sentinel-chat-requirements-token", strings.TrimSpace(chatToken))
placeholder
	if strings.TrimSpace(proofToken) != "" {
		prepareHeaders.Set("openai-sentinel-proof-token", strings.TrimSpace(proofToken))
placeholder
	var result struct {
		ConduitToken string `json:"conduit_token"`
placeholder
	resp, err := client.R().
		SetContext(ctx).
		SetHeaders(headerToMap(prepareHeaders)).
		SetBodyJsonMarshal(payload).
		SetSuccessResult(&result).
		Post(openAIChatGPTConversationPrepareURL)
	if err != nil {
		return "", err
placeholder
	if !resp.IsSuccessState() {
		return "", newOpenAIImageStatusError(resp, "conversation prepare failed")
placeholder
	return strings.TrimSpace(result.ConduitToken), nil
placeholder

type openAIUploadedImage struct {
	FileID   string
	FileName string
	FileSize int
	MimeType string
	Width    int
	Height   int
placeholder

func uploadOpenAIImageFiles(ctx context.Context, client *req.Client, headers http.Header, uploads []OpenAIImagesUpload) ([]openAIUploadedImage, error) {
	if len(uploads) == 0 {
		return nil, nil
placeholder
	results := make([]openAIUploadedImage, 0, len(uploads))
	for i := range uploads {
		item := uploads[i]
		fileName := coalesceOpenAIFileName(item.FileName, "image.png")
		payload := map[string]any{
			"file_name": fileName,
			"file_size": len(item.Data),
			"use_case":  "multimodal",
	placeholder
		var created struct {
			FileID    string `json:"file_id"`
			UploadURL string `json:"upload_url"`
	placeholder
		resp, err := client.R().
			SetContext(ctx).
			SetHeaders(headerToMap(headers)).
			SetBodyJsonMarshal(payload).
			SetSuccessResult(&created).
			Post(openAIChatGPTFilesURL)
		if err != nil {
			return nil, err
	placeholder
		if !resp.IsSuccessState() || strings.TrimSpace(created.FileID) == "" || strings.TrimSpace(created.UploadURL) == "" {
			return nil, newOpenAIImageStatusError(resp, "create upload slot failed")
	placeholder

		uploadHeaders := map[string]string{
			"Content-Type":   coalesceOpenAIFileName(item.ContentType, "application/octet-stream"),
			"Origin":         "https://chatgpt.com",
			"x-ms-blob-type": "BlockBlob",
			"x-ms-version":   "2020-04-08",
			"User-Agent":     headers.Get("User-Agent"),
	placeholder
		putResp, err := client.R().
			SetContext(ctx).
			SetHeaders(uploadHeaders).
			SetBody(item.Data).
			DisableAutoReadResponse().
			Put(created.UploadURL)
		if err != nil {
			return nil, err
	placeholder
		if putResp.Response != nil && putResp.Body != nil {
			_, _ = io.Copy(io.Discard, putResp.Body)
			_ = putResp.Body.Close()
	placeholder
		if putResp.StatusCode < 200 || putResp.StatusCode >= 300 {
			return nil, newOpenAIImageStatusError(putResp, "upload image bytes failed")
	placeholder

		uploadedResp, err := client.R().
			SetContext(ctx).
			SetHeaders(headerToMap(headers)).
			SetBodyJsonMarshal(map[string]any{placeholder).
			Post(fmt.Sprintf("%s/%s/uploaded", openAIChatGPTFilesURL, created.FileID))
		if err != nil {
			return nil, err
	placeholder
		if !uploadedResp.IsSuccessState() {
			return nil, newOpenAIImageStatusError(uploadedResp, "mark upload complete failed")
	placeholder

		results = append(results, openAIUploadedImage{
			FileID:   created.FileID,
			FileName: fileName,
			FileSize: len(item.Data),
			MimeType: coalesceOpenAIFileName(item.ContentType, "application/octet-stream"),
			Width:    item.Width,
			Height:   item.Height,
	placeholder)
placeholder
	return results, nil
placeholder

func coalesceOpenAIFileName(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
placeholder
	return value
placeholder

func buildOpenAIImageConversationRequest(parsed *OpenAIImagesRequest, parentMessageID string, uploads []openAIUploadedImage) map[string]any {
	parts := []any{coalesceOpenAIFileName(parsed.Prompt, "Generate an image.")placeholder
	attachments := make([]map[string]any, 0, len(uploads))
	if len(uploads) > 0 {
		parts = make([]any, 0, len(uploads)+1)
		for _, upload := range uploads {
			parts = append(parts, map[string]any{
				"content_type":  "image_asset_pointer",
				"asset_pointer": "file-service://" + upload.FileID,
				"size_bytes":    upload.FileSize,
				"width":         upload.Width,
				"height":        upload.Height,
		placeholder)
			attachment := map[string]any{
				"id":       upload.FileID,
				"mimeType": upload.MimeType,
				"name":     upload.FileName,
				"size":     upload.FileSize,
		placeholder
			if upload.Width > 0 {
				attachment["width"] = upload.Width
		placeholder
			if upload.Height > 0 {
				attachment["height"] = upload.Height
		placeholder
			attachments = append(attachments, attachment)
	placeholder
		parts = append(parts, coalesceOpenAIFileName(parsed.Prompt, "Edit this image."))
placeholder

	contentType := "text"
	if len(uploads) > 0 {
		contentType = "multimodal_text"
placeholder
	metadata := map[string]any{
		"developer_mode_connector_ids": []any{placeholder,
		"selected_github_repos":        []any{placeholder,
		"selected_all_github_repos":    false,
		"system_hints":                 []string{"picture_v2"placeholder,
		"serialization_metadata": map[string]any{
			"custom_symbol_offsets": []any{placeholder,
	placeholder,
placeholder
	message := map[string]any{
		"id":     uuid.NewString(),
		"author": map[string]any{"role": "user"placeholder,
		"content": map[string]any{
			"content_type": contentType,
			"parts":        parts,
	placeholder,
		"metadata":    metadata,
		"create_time": float64(time.Now().UnixMilli()) / 1000,
placeholder
	if len(attachments) > 0 {
		metadata["attachments"] = attachments
placeholder

	return map[string]any{
		"action":                               "next",
		"client_prepare_state":                 "sent",
		"parent_message_id":                    parentMessageID,
		"model":                                "auto",
		"timezone_offset_min":                  openAITimezoneOffsetMinutes(),
		"timezone":                             openAITimezoneName(),
		"conversation_mode":                    map[string]any{"kind": "primary_assistant"placeholder,
		"enable_message_followups":             true,
		"system_hints":                         []string{"picture_v2"placeholder,
		"supports_buffering":                   true,
		"supported_encodings":                  []string{"v1"placeholder,
		"paragen_cot_summary_display_override": "allow",
		"force_parallel_switch":                "auto",
		"client_contextual_info": map[string]any{
			"is_dark_mode":      false,
			"time_since_loaded": 200,
			"page_height":       900,
			"page_width":        1440,
			"pixel_ratio":       1,
			"screen_height":     1080,
			"screen_width":      1920,
			"app_name":          "chatgpt.com",
	placeholder,
		"messages": []any{messageplaceholder,
placeholder
placeholder

type openAIImagePointerInfo struct {
	Pointer string
	Prompt  string
placeholder

type openAIImageToolMessage struct {
	MessageID    string
	CreateTime   float64
	PointerInfos []openAIImagePointerInfo
placeholder

func readOpenAIImageConversationStream(resp *req.Response, startTime time.Time) (string, []openAIImagePointerInfo, OpenAIUsage, *int, error) {
	if resp == nil || resp.Response == nil || resp.Body == nil {
		return "", nil, OpenAIUsage{placeholder, nil, fmt.Errorf("empty conversation response")
placeholder
	reader := bufio.NewReader(resp.Body)
	var (
		conversationID string
		firstTokenMs   *int
		usage          OpenAIUsage
		pointers       []openAIImagePointerInfo
	)

	for {
		line, err := reader.ReadString('\n')
		if strings.TrimSpace(line) != "" && firstTokenMs == nil {
			ms := int(time.Since(startTime).Milliseconds())
			firstTokenMs = &ms
	placeholder
		if data, ok := extractOpenAISSEDataLine(strings.TrimRight(line, "\r\n")); ok && data != "" && data != "[DONE]" {
			dataBytes := []byte(data)
			if conversationID == "" {
				conversationID = strings.TrimSpace(gjson.GetBytes(dataBytes, "v.conversation_id").String())
				if conversationID == "" {
					conversationID = strings.TrimSpace(gjson.GetBytes(dataBytes, "conversation_id").String())
			placeholder
		placeholder
			mergeOpenAIUsage(&usage, dataBytes)
			pointers = mergeOpenAIImagePointerInfos(pointers, collectOpenAIImagePointers(dataBytes))
	placeholder
		if err == io.EOF {
			break
	placeholder
		if err != nil {
			return "", nil, OpenAIUsage{placeholder, firstTokenMs, err
	placeholder
placeholder
	return conversationID, pointers, usage, firstTokenMs, nil
placeholder

func collectOpenAIImagePointers(body []byte) []openAIImagePointerInfo {
	if len(body) == 0 {
		return nil
placeholder
	matches := openAIImagePointerMatches(body)
	if len(matches) == 0 {
		return nil
placeholder
	prompt := ""
	for _, path := range []string{
		"message.metadata.dalle.prompt",
		"metadata.dalle.prompt",
		"revised_prompt",
placeholder {
		if value := strings.TrimSpace(gjson.GetBytes(body, path).String()); value != "" {
			prompt = value
			break
	placeholder
placeholder
	out := make([]openAIImagePointerInfo, 0, len(matches))
	for _, pointer := range matches {
		out = append(out, openAIImagePointerInfo{Pointer: pointer, Prompt: promptplaceholder)
placeholder
	return out
placeholder

func openAIImagePointerMatches(body []byte) []string {
	raw := string(body)
	matches := make([]string, 0, 4)
	for _, prefix := range []string{"file-service://", "sediment://"placeholder {
		start := 0
		for {
			idx := strings.Index(raw[start:], prefix)
			if idx < 0 {
				break
		placeholder
			idx += start
			end := idx + len(prefix)
			for end < len(raw) {
				ch := raw[end]
				if ch != '-' && ch != '_' &&
					(ch < '0' || ch > '9') &&
					(ch < 'a' || ch > 'z') &&
					(ch < 'A' || ch > 'Z') {
					break
			placeholder
				end++
		placeholder
			matches = append(matches, raw[idx:end])
			start = end
	placeholder
placeholder
	return dedupeStrings(matches)
placeholder

func mergeOpenAIImagePointerInfos(existing []openAIImagePointerInfo, next []openAIImagePointerInfo) []openAIImagePointerInfo {
	if len(next) == 0 {
		return existing
placeholder
	seen := make(map[string]openAIImagePointerInfo, len(existing)+len(next))
	out := make([]openAIImagePointerInfo, 0, len(existing)+len(next))
	for _, item := range existing {
		seen[item.Pointer] = item
		out = append(out, item)
placeholder
	for _, item := range next {
		if existingItem, ok := seen[item.Pointer]; ok {
			if existingItem.Prompt == "" && item.Prompt != "" {
				for i := range out {
					if out[i].Pointer == item.Pointer {
						out[i].Prompt = item.Prompt
						break
				placeholder
			placeholder
		placeholder
			continue
	placeholder
		seen[item.Pointer] = item
		out = append(out, item)
placeholder
	return out
placeholder

func hasOpenAIFileServicePointerInfos(items []openAIImagePointerInfo) bool {
	for _, item := range items {
		if strings.HasPrefix(item.Pointer, "file-service://") {
			return true
	placeholder
placeholder
	return false
placeholder

func preferOpenAIFileServicePointerInfos(items []openAIImagePointerInfo) []openAIImagePointerInfo {
	if !hasOpenAIFileServicePointerInfos(items) {
		return items
placeholder
	out := make([]openAIImagePointerInfo, 0, len(items))
	for _, item := range items {
		if strings.HasPrefix(item.Pointer, "file-service://") {
			out = append(out, item)
	placeholder
placeholder
	return out
placeholder

func extractOpenAIImageToolMessages(mapping map[string]any) []openAIImageToolMessage {
	if len(mapping) == 0 {
		return nil
placeholder
	out := make([]openAIImageToolMessage, 0, 4)
	for messageID, raw := range mapping {
		node, _ := raw.(map[string]any)
		if node == nil {
			continue
	placeholder
		message, _ := node["message"].(map[string]any)
		if message == nil {
			continue
	placeholder
		author, _ := message["author"].(map[string]any)
		metadata, _ := message["metadata"].(map[string]any)
		content, _ := message["content"].(map[string]any)
		if author == nil || metadata == nil || content == nil {
			continue
	placeholder
		if role, _ := author["role"].(string); role != "tool" {
			continue
	placeholder
		if asyncTaskType, _ := metadata["async_task_type"].(string); asyncTaskType != "image_gen" {
			continue
	placeholder
		if contentType, _ := content["content_type"].(string); contentType != "multimodal_text" {
			continue
	placeholder
		prompt := ""
		if title, _ := metadata["image_gen_title"].(string); strings.TrimSpace(title) != "" {
			prompt = strings.TrimSpace(title)
	placeholder
		item := openAIImageToolMessage{MessageID: messageIDplaceholder
		if createTime, ok := message["create_time"].(float64); ok {
			item.CreateTime = createTime
	placeholder
		parts, _ := content["parts"].([]any)
		for _, part := range parts {
			switch value := part.(type) {
			case map[string]any:
				if assetPointer, _ := value["asset_pointer"].(string); strings.TrimSpace(assetPointer) != "" {
					for _, pointer := range openAIImagePointerMatches([]byte(assetPointer)) {
						item.PointerInfos = append(item.PointerInfos, openAIImagePointerInfo{
							Pointer: pointer,
							Prompt:  prompt,
					placeholder)
				placeholder
			placeholder
			case string:
				for _, pointer := range openAIImagePointerMatches([]byte(value)) {
					item.PointerInfos = append(item.PointerInfos, openAIImagePointerInfo{
						Pointer: pointer,
						Prompt:  prompt,
				placeholder)
			placeholder
		placeholder
	placeholder
		if len(item.PointerInfos) == 0 {
			continue
	placeholder
		item.PointerInfos = mergeOpenAIImagePointerInfos(nil, item.PointerInfos)
		out = append(out, item)
placeholder
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreateTime < out[j].CreateTime
placeholder)
	return out
placeholder

func pollOpenAIImageConversation(ctx context.Context, client *req.Client, headers http.Header, conversationID string) ([]openAIImagePointerInfo, error) {
	conversationID = strings.TrimSpace(conversationID)
	if conversationID == "" {
		return nil, nil
placeholder
	deadline := time.Now().Add(90 * time.Second)
	interval := 3 * time.Second
	previewWait := 15 * time.Second
	var (
		lastErr     error
		firstToolAt time.Time
	)
	for time.Now().Before(deadline) {
		resp, err := client.R().
			SetContext(ctx).
			SetHeaders(headerToMap(headers)).
			DisableAutoReadResponse().
			Get(fmt.Sprintf("https://chatgpt.com/backend-api/conversation/%s", conversationID))
		if err != nil {
			lastErr = err
	placeholder else {
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				body, readErr := io.ReadAll(resp.Body)
				_ = resp.Body.Close()
				if readErr != nil {
					lastErr = readErr
					goto waitNextPoll
			placeholder
				pointers := mergeOpenAIImagePointerInfos(nil, collectOpenAIImagePointers(body))
				var decoded map[string]any
				if err := json.Unmarshal(body, &decoded); err == nil {
					if mapping, _ := decoded["mapping"].(map[string]any); len(mapping) > 0 {
						toolMessages := extractOpenAIImageToolMessages(mapping)
						if len(toolMessages) > 0 && firstToolAt.IsZero() {
							firstToolAt = time.Now()
					placeholder
						for _, msg := range toolMessages {
							pointers = mergeOpenAIImagePointerInfos(pointers, msg.PointerInfos)
					placeholder
				placeholder
			placeholder
				if hasOpenAIFileServicePointerInfos(pointers) {
					return preferOpenAIFileServicePointerInfos(pointers), nil
			placeholder
				if len(pointers) > 0 && !firstToolAt.IsZero() && time.Since(firstToolAt) >= previewWait {
					return pointers, nil
			placeholder
		placeholder else {
				statusErr := newOpenAIImageStatusError(resp, "conversation poll failed")
				if isOpenAIImageTransientConversationNotFoundError(statusErr) {
					lastErr = statusErr
					goto waitNextPoll
			placeholder
				return nil, statusErr
		placeholder
	placeholder

	waitNextPoll:
		timer := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
		placeholder
			return nil, ctx.Err()
		case <-timer.C:
	placeholder
placeholder
	return nil, lastErr
placeholder

func buildOpenAIImageResponse(
	ctx context.Context,
	client *req.Client,
	headers http.Header,
	conversationID string,
	pointers []openAIImagePointerInfo,
) ([]byte, int, error) {
	type responseItem struct {
		B64JSON       string `json:"b64_json"`
		RevisedPrompt string `json:"revised_prompt,omitempty"`
placeholder
	items := make([]responseItem, 0, len(pointers))
	for _, pointer := range pointers {
		downloadURL, err := fetchOpenAIImageDownloadURL(ctx, client, headers, conversationID, pointer.Pointer)
		if err != nil {
			return nil, 0, err
	placeholder
		data, err := downloadOpenAIImageBytes(ctx, client, headers, downloadURL)
		if err != nil {
			return nil, 0, err
	placeholder
		items = append(items, responseItem{
			B64JSON:       base64.StdEncoding.EncodeToString(data),
			RevisedPrompt: pointer.Prompt,
	placeholder)
placeholder
	payload := map[string]any{
		"created": time.Now().Unix(),
		"data":    items,
placeholder
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, err
placeholder
	return body, len(items), nil
placeholder

func fetchOpenAIImageDownloadURL(
	ctx context.Context,
	client *req.Client,
	headers http.Header,
	conversationID string,
	pointer string,
) (string, error) {
	url := ""
	allowConversationRetry := false
	switch {
	case strings.HasPrefix(pointer, "file-service://"):
		fileID := strings.TrimPrefix(pointer, "file-service://")
		url = fmt.Sprintf("%s/%s/download", openAIChatGPTFilesURL, fileID)
	case strings.HasPrefix(pointer, "sediment://"):
		attachmentID := strings.TrimPrefix(pointer, "sediment://")
		url = fmt.Sprintf("https://chatgpt.com/backend-api/conversation/%s/attachment/%s/download", conversationID, attachmentID)
		allowConversationRetry = true
	default:
		return "", fmt.Errorf("unsupported image pointer: %s", pointer)
placeholder

	var lastErr error
	for attempt := 0; attempt < 8; attempt++ {
		var result struct {
			DownloadURL string `json:"download_url"`
	placeholder
		resp, err := client.R().
			SetContext(ctx).
			SetHeaders(headerToMap(headers)).
			SetSuccessResult(&result).
			Get(url)
		if err != nil {
			lastErr = err
	placeholder else if resp.IsSuccessState() && strings.TrimSpace(result.DownloadURL) != "" {
			return strings.TrimSpace(result.DownloadURL), nil
	placeholder else {
			statusErr := newOpenAIImageStatusError(resp, "fetch image download url failed")
			if !allowConversationRetry || !isOpenAIImageTransientConversationNotFoundError(statusErr) {
				return "", statusErr
		placeholder
			lastErr = statusErr
	placeholder
		if attempt == 7 {
			break
	placeholder
		timer := time.NewTimer(750 * time.Millisecond)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
		placeholder
			return "", ctx.Err()
		case <-timer.C:
	placeholder
placeholder
	if lastErr == nil {
		lastErr = fmt.Errorf("fetch image download url failed")
placeholder
	return "", lastErr
placeholder

func downloadOpenAIImageBytes(ctx context.Context, client *req.Client, headers http.Header, downloadURL string) ([]byte, error) {
	request := client.R().
		SetContext(ctx).
		DisableAutoReadResponse()

	if strings.HasPrefix(downloadURL, openAIChatGPTStartURL) {
		downloadHeaders := cloneHTTPHeader(headers)
		downloadHeaders.Set("Accept", "image/*,*/*;q=0.8")
		downloadHeaders.Del("Content-Type")
		request.SetHeaders(headerToMap(downloadHeaders))
placeholder else {
		userAgent := strings.TrimSpace(headers.Get("User-Agent"))
		if userAgent == "" {
			userAgent = openAIImageBackendUserAgent
	placeholder
		request.SetHeader("User-Agent", userAgent)
placeholder

	resp, err := request.Get(downloadURL)
	if err != nil {
		return nil, err
placeholder
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
	placeholder
placeholder()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, newOpenAIImageStatusError(resp, "download image bytes failed")
placeholder
	return io.ReadAll(resp.Body)
placeholder

func handleOpenAIImageBackendError(resp *req.Response) error {
	return newOpenAIImageStatusError(resp, "backend-api request failed")
placeholder

type openAIImageStatusError struct {
	StatusCode      int
	Message         string
	ResponseBody    []byte
	ResponseHeaders http.Header
	RequestID       string
	URL             string
placeholder

func (e *openAIImageStatusError) Error() string {
	if e == nil {
		return "openai image backend request failed"
placeholder
	if e.Message != "" {
		return e.Message
placeholder
	if e.StatusCode > 0 {
		return fmt.Sprintf("openai image backend request failed: status %d", e.StatusCode)
placeholder
	return "openai image backend request failed"
placeholder

func newOpenAIImageStatusError(resp *req.Response, fallback string) error {
	if resp == nil {
		if strings.TrimSpace(fallback) == "" {
			fallback = "openai image backend request failed"
	placeholder
		return fmt.Errorf("%s", fallback)
placeholder

	statusCode := resp.StatusCode
	headers := http.Header(nil)
	requestID := ""
	requestURL := ""
	body := []byte(nil)

	if resp.Response != nil {
		headers = resp.Header.Clone()
		requestID = strings.TrimSpace(resp.Header.Get("x-request-id"))
		if resp.Response.Request != nil && resp.Response.Request.URL != nil {
			requestURL = resp.Response.Request.URL.String()
	placeholder
		if resp.Body != nil {
			body, _ = io.ReadAll(io.LimitReader(resp.Body, 2<<20))
			_ = resp.Body.Close()
	placeholder
placeholder

	message := sanitizeUpstreamErrorMessage(extractUpstreamErrorMessage(body))
	if message == "" {
		prefix := strings.TrimSpace(fallback)
		if prefix == "" {
			prefix = "openai image backend request failed"
	placeholder
		message = fmt.Sprintf("%s: status %d", prefix, statusCode)
placeholder

	return &openAIImageStatusError{
		StatusCode:      statusCode,
		Message:         message,
		ResponseBody:    body,
		ResponseHeaders: headers,
		RequestID:       requestID,
		URL:             requestURL,
placeholder
placeholder

func newOpenAIImageSyntheticStatusError(statusCode int, message string, requestURL string) *openAIImageStatusError {
	message = sanitizeUpstreamErrorMessage(strings.TrimSpace(message))
	if message == "" {
		message = "openai image backend request failed"
placeholder
	var body []byte
	if payload, err := json.Marshal(map[string]string{"detail": messageplaceholder); err == nil {
		body = payload
placeholder
	return &openAIImageStatusError{
		StatusCode:   statusCode,
		Message:      message,
		ResponseBody: body,
		URL:          strings.TrimSpace(requestURL),
placeholder
placeholder

func isOpenAIImageTransientConversationNotFoundError(err error) bool {
	statusErr, ok := err.(*openAIImageStatusError)
	if !ok || statusErr == nil || statusErr.StatusCode != http.StatusNotFound {
		return false
placeholder
	msg := strings.ToLower(strings.TrimSpace(statusErr.Message))
	if strings.Contains(msg, "conversation_not_found") {
		return true
placeholder
	if strings.Contains(msg, "conversation") && strings.Contains(msg, "not found") {
		return true
placeholder
	bodyMsg := strings.ToLower(strings.TrimSpace(extractUpstreamErrorMessage(statusErr.ResponseBody)))
	if strings.Contains(bodyMsg, "conversation_not_found") {
		return true
placeholder
	return strings.Contains(bodyMsg, "conversation") && strings.Contains(bodyMsg, "not found")
placeholder

func (s *OpenAIGatewayService) wrapOpenAIImageBackendError(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	err error,
) error {
	var statusErr *openAIImageStatusError
	if !errors.As(err, &statusErr) || statusErr == nil {
		return err
placeholder

	upstreamMsg := sanitizeUpstreamErrorMessage(statusErr.Message)
	appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
		Platform:           account.Platform,
		AccountID:          account.ID,
		AccountName:        account.Name,
		UpstreamStatusCode: statusErr.StatusCode,
		UpstreamRequestID:  statusErr.RequestID,
		UpstreamURL:        safeUpstreamURL(statusErr.URL),
		Kind:               "request_error",
		Message:            upstreamMsg,
placeholder)
	setOpsUpstreamError(c, statusErr.StatusCode, upstreamMsg, "")

	if s.shouldFailoverOpenAIUpstreamResponse(statusErr.StatusCode, upstreamMsg, statusErr.ResponseBody) {
		if s.rateLimitService != nil {
			s.rateLimitService.HandleUpstreamError(ctx, account, statusErr.StatusCode, statusErr.ResponseHeaders, statusErr.ResponseBody)
	placeholder
		appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
			Platform:           account.Platform,
			AccountID:          account.ID,
			AccountName:        account.Name,
			UpstreamStatusCode: statusErr.StatusCode,
			UpstreamRequestID:  statusErr.RequestID,
			UpstreamURL:        safeUpstreamURL(statusErr.URL),
			Kind:               "failover",
			Message:            upstreamMsg,
	placeholder)
		retryableOnSameAccount := account.IsPoolMode() && isPoolModeRetryableStatus(statusErr.StatusCode)
		if strings.Contains(strings.ToLower(statusErr.Message), "unsupported challenge") {
			retryableOnSameAccount = false
	placeholder
		return &UpstreamFailoverError{
			StatusCode:             statusErr.StatusCode,
			ResponseBody:           statusErr.ResponseBody,
			RetryableOnSameAccount: retryableOnSameAccount,
	placeholder
placeholder

	return statusErr
placeholder

func cloneHTTPHeader(src http.Header) http.Header {
	dst := make(http.Header, len(src))
	for key, values := range src {
		copied := make([]string, len(values))
		copy(copied, values)
		dst[key] = copied
placeholder
	return dst
placeholder

func headerToMap(header http.Header) map[string]string {
	if len(header) == 0 {
		return nil
placeholder
	result := make(map[string]string, len(header))
	for key, values := range header {
		if len(values) == 0 {
			continue
	placeholder
		result[key] = values[0]
placeholder
	return result
placeholder

func openAITimezoneOffsetMinutes() int {
	_, offset := time.Now().Zone()
	return offset / 60
placeholder

func openAITimezoneName() string {
	return time.Now().Location().String()
placeholder

func generateOpenAIRequirementsToken(userAgent string) string {
	config := []any{
		"core" + strconv.Itoa(3008),
		time.Now().UTC().Format(time.RFC1123),
		nil,
		0.123456,
		coalesceOpenAIFileName(strings.TrimSpace(userAgent), openAIImageBackendUserAgent),
		nil,
		"prod-openai-images",
		"en-US",
		"en-US,en",
		0,
		"navigator.webdriver",
		"location",
		"document.body",
		float64(time.Now().UnixMilli()) / 1000,
		uuid.NewString(),
		"",
		8,
		time.Now().Unix(),
placeholder
	answer, solved := generateOpenAIChallengeAnswer(strconv.FormatInt(time.Now().UnixNano(), 10), openAIImageRequirementsDiff, config)
	if solved {
		return "gAAAAAC" + answer
placeholder
	return ""
placeholder

func generateOpenAIChallengeAnswer(seed string, difficulty string, config []any) (string, bool) {
	diffBytes, err := hex.DecodeString(difficulty)
	if err != nil {
		return "", false
placeholder
	p1 := []byte(jsonCompactSlice(config[:3], true))
	p2 := []byte(jsonCompactSlice(config[4:9], false))
	p3 := []byte(jsonCompactSlice(config[10:], false))
	seedBytes := []byte(seed)

	for i := 0; i < 100000; i++ {
		payload := fmt.Sprintf("%s%d,%s,%d,%s", p1, i, p2, i>>1, p3)
		encoded := base64.StdEncoding.EncodeToString([]byte(payload))
		sum := sha3.Sum512(append(seedBytes, []byte(encoded)...))
		if bytes.Compare(sum[:len(diffBytes)], diffBytes) <= 0 {
			return encoded, true
	placeholder
placeholder
	return "", false
placeholder

func jsonCompactSlice(values []any, trimSuffixComma bool) string {
	raw, _ := json.Marshal(values)
	text := string(raw)
	if trimSuffixComma {
		return strings.TrimSuffix(text, "]")
placeholder
	return strings.TrimPrefix(text, "[")
placeholder

func generateOpenAIProofToken(required bool, seed string, difficulty string, userAgent string) string {
	if !required || strings.TrimSpace(seed) == "" || strings.TrimSpace(difficulty) == "" {
		return ""
placeholder
	screen := 3008
	if len(seed)%2 == 0 {
		screen = 4010
placeholder
	proofToken := []any{
		screen,
		time.Now().UTC().Format(time.RFC1123),
		nil,
		0,
		coalesceOpenAIFileName(strings.TrimSpace(userAgent), openAIImageBackendUserAgent),
		"https://chatgpt.com/",
		"dpl=openai-images",
		"en",
		"en-US",
		nil,
		"plugins[object PluginArray]",
		"_reactListening",
		"alert",
placeholder
	diffLen := len(difficulty)
	for i := 0; i < 100000; i++ {
		proofToken[3] = i
		raw, _ := json.Marshal(proofToken)
		encoded := base64.StdEncoding.EncodeToString(raw)
		sum := sha3.Sum512([]byte(seed + encoded))
		if strings.Compare(hex.EncodeToString(sum[:])[:diffLen], difficulty) <= 0 {
			return "gAAAAAB" + encoded
	placeholder
placeholder
	fallbackBase := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%q", seed)))
	return "gAAAAABwQ8Lk5FbGpA2NcR9dShT6gYjU7VxZ4D" + fallbackBase
placeholder

func dedupeStrings(values []string) []string {
	if len(values) == 0 {
		return nil
placeholder
	seen := make(map[string]struct{placeholder, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
	placeholder
		seen[value] = struct{placeholder{placeholder
		out = append(out, value)
placeholder
	return out
placeholder
