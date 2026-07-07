package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/Wei-Shaw/sub2api/internal/util/responseheaders"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type GrokMediaEndpoint string

const (
	GrokMediaEndpointImagesGenerations GrokMediaEndpoint = "images_generations"
	GrokMediaEndpointImagesEdits       GrokMediaEndpoint = "images_edits"
	GrokMediaEndpointVideosGenerations GrokMediaEndpoint = "videos_generations"
	GrokMediaEndpointVideoStatus       GrokMediaEndpoint = "video_status"
)

func (e GrokMediaEndpoint) RequiresRequestBody() bool {
	return e != GrokMediaEndpointVideoStatus
placeholder

func (e GrokMediaEndpoint) IsGenerationRequest() bool {
	switch e {
	case GrokMediaEndpointImagesGenerations, GrokMediaEndpointImagesEdits, GrokMediaEndpointVideosGenerations:
		return true
	default:
		return false
placeholder
placeholder

type GrokMediaRequestInfo struct {
	Model          string
	Prompt         string
	N              int
	Size           string
	SizeTier       string
	InputImageURLs []string
	MaskImageURL   string
	Uploads        []OpenAIImagesUpload
	MaskUpload     *OpenAIImagesUpload
placeholder

func (r GrokMediaRequestInfo) ModerationBody() []byte {
	payload := map[string]any{placeholder
	if prompt := strings.TrimSpace(r.Prompt); prompt != "" {
		payload["prompt"] = prompt
placeholder

	images := make([]map[string]string, 0, len(r.InputImageURLs)+len(r.Uploads)+1)
	for _, imageURL := range r.InputImageURLs {
		if imageURL = strings.TrimSpace(imageURL); imageURL != "" {
			images = append(images, map[string]string{"image_url": imageURLplaceholder)
	placeholder
placeholder
	for _, upload := range r.Uploads {
		if dataURL := upload.ModerationDataURL(); dataURL != "" {
			images = append(images, map[string]string{"image_url": dataURLplaceholder)
	placeholder
placeholder
	if maskURL := strings.TrimSpace(r.MaskImageURL); maskURL != "" {
		images = append(images, map[string]string{"image_url": maskURLplaceholder)
placeholder
	if r.MaskUpload != nil {
		if dataURL := r.MaskUpload.ModerationDataURL(); dataURL != "" {
			images = append(images, map[string]string{"image_url": dataURLplaceholder)
	placeholder
placeholder
	if len(images) > 0 {
		payload["images"] = images
placeholder
	if len(payload) == 0 {
		return nil
placeholder
	body, err := json.Marshal(payload)
	if err != nil {
		return nil
placeholder
	return body
placeholder

func (e GrokMediaEndpoint) httpMethod() string {
	if e == GrokMediaEndpointVideoStatus {
		return http.MethodGet
placeholder
	return http.MethodPost
placeholder

func ExtractGrokMediaModel(contentType string, body []byte) string {
	return ParseGrokMediaRequest(contentType, body).Model
placeholder

func ParseGrokMediaRequest(contentType string, body []byte) GrokMediaRequestInfo {
	info := GrokMediaRequestInfo{N: 1placeholder
	if gjson.ValidBytes(body) {
		parseGrokMediaJSONRequest(body, &info)
placeholder else {
		parseGrokMediaMultipartRequest(contentType, body, &info)
placeholder
	info.Model = strings.TrimSpace(info.Model)
	info.Prompt = strings.TrimSpace(info.Prompt)
	info.Size = strings.TrimSpace(info.Size)
	info.SizeTier = NormalizeImageBillingTierOrDefault(info.Size)
	if info.N <= 0 {
		info.N = 1
placeholder
	return info
placeholder

func parseGrokMediaJSONRequest(body []byte, info *GrokMediaRequestInfo) {
	if info == nil {
		return
placeholder
	info.Model = strings.TrimSpace(gjson.GetBytes(body, "model").String())
	info.Prompt = strings.TrimSpace(gjson.GetBytes(body, "prompt").String())
	info.Size = strings.TrimSpace(gjson.GetBytes(body, "size").String())
	if n := gjson.GetBytes(body, "n"); n.Exists() && n.Type == gjson.Number {
		info.N = int(n.Int())
placeholder
	appendJSONImageURLs := func(value gjson.Result) {
		if !value.Exists() {
			return
	placeholder
		switch {
		case value.IsArray():
			for _, item := range value.Array() {
				if imageURL := strings.TrimSpace(item.Get("image_url").String()); imageURL != "" {
					info.InputImageURLs = append(info.InputImageURLs, imageURL)
					continue
			placeholder
				if item.Type == gjson.String {
					imageURL := strings.TrimSpace(item.String())
					if imageURL == "" {
						continue
				placeholder
					info.InputImageURLs = append(info.InputImageURLs, imageURL)
			placeholder
		placeholder
		default:
			if imageURL := strings.TrimSpace(value.Get("image_url").String()); imageURL != "" {
				info.InputImageURLs = append(info.InputImageURLs, imageURL)
				return
		placeholder
			if value.Type == gjson.String {
				imageURL := strings.TrimSpace(value.String())
				if imageURL == "" {
					return
			placeholder
				info.InputImageURLs = append(info.InputImageURLs, imageURL)
		placeholder
	placeholder
placeholder
	appendJSONImageURLs(gjson.GetBytes(body, "image"))
	appendJSONImageURLs(gjson.GetBytes(body, "images"))
	info.MaskImageURL = strings.TrimSpace(gjson.GetBytes(body, "mask.image_url").String())
placeholder

func parseGrokMediaMultipartRequest(contentType string, body []byte, info *GrokMediaRequestInfo) {
	if info == nil {
		return
placeholder
	mediaType, params, err := mime.ParseMediaType(strings.TrimSpace(contentType))
	if err != nil || !strings.EqualFold(mediaType, "multipart/form-data") {
		return
placeholder
	boundary := strings.TrimSpace(params["boundary"])
	if boundary == "" {
		return
placeholder
	reader := multipart.NewReader(bytes.NewReader(body), boundary)
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			return
	placeholder
		if err != nil {
			return
	placeholder
		name := strings.TrimSpace(part.FormName())
		if name == "" {
			_ = part.Close()
			continue
	placeholder
		data, err := io.ReadAll(io.LimitReader(part, openAIImageMaxUploadPartSize))
		_ = part.Close()
		if err != nil {
			return
	placeholder
		fileName := strings.TrimSpace(part.FileName())
		partContentType := strings.TrimSpace(part.Header.Get("Content-Type"))
		if fileName != "" {
			upload := OpenAIImagesUpload{
				FieldName:   name,
				FileName:    fileName,
				ContentType: partContentType,
				Data:        data,
		placeholder
			if name == "mask" {
				info.MaskUpload = &upload
				continue
		placeholder
			if name == "image" || strings.HasPrefix(name, "image[") {
				info.Uploads = append(info.Uploads, upload)
		placeholder
			continue
	placeholder

		value := strings.TrimSpace(string(data))
		switch name {
		case "model":
			info.Model = value
		case "prompt":
			info.Prompt = value
		case "size":
			info.Size = value
		case "n":
			if n, err := strconv.Atoi(value); err == nil {
				info.N = n
		placeholder
		case "image", "image_url":
			if value != "" {
				info.InputImageURLs = append(info.InputImageURLs, value)
		placeholder
		case "mask", "mask_image_url":
			info.MaskImageURL = value
	placeholder
placeholder
placeholder

func GrokMediaVideoRequestSessionHash(requestID string) string {
	requestID = strings.TrimSpace(requestID)
	if requestID == "" {
		return ""
placeholder
	return "grok-video:" + DeriveSessionHashFromSeed(requestID)
placeholder

func (s *OpenAIGatewayService) BindGrokMediaVideoRequestAccount(ctx context.Context, groupID *int64, requestID string, accountID int64) error {
	return s.BindStickySession(ctx, groupID, GrokMediaVideoRequestSessionHash(requestID), accountID)
placeholder

func (e GrokMediaEndpoint) upstreamURL(baseURL, requestID string) (string, error) {
	switch e {
	case GrokMediaEndpointImagesGenerations:
		return xai.BuildImagesGenerationsURL(baseURL)
	case GrokMediaEndpointImagesEdits:
		return xai.BuildImagesEditsURL(baseURL)
	case GrokMediaEndpointVideosGenerations:
		return xai.BuildVideosGenerationsURL(baseURL)
	case GrokMediaEndpointVideoStatus:
		return xai.BuildVideoURL(baseURL, requestID)
	default:
		return "", fmt.Errorf("unsupported grok media endpoint: %s", e)
placeholder
placeholder

func (s *OpenAIGatewayService) ForwardGrokMedia(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	endpoint GrokMediaEndpoint,
	requestID string,
	body []byte,
	contentType string,
) (*OpenAIForwardResult, error) {
	startTime := time.Now()
	if account == nil {
		return nil, fmt.Errorf("grok account is required")
placeholder
	if account.Platform != PlatformGrok {
		return nil, fmt.Errorf("account platform %s is not supported for grok media", account.Platform)
placeholder

	token, _, err := s.GetAccessToken(ctx, account)
	if err != nil {
		return nil, err
placeholder
	targetURL, err := endpoint.upstreamURL(account.GetGrokBaseURL(), requestID)
	if err != nil {
		return nil, err
placeholder

	body, contentType, err = prepareGrokMediaForwardBody(endpoint, body, contentType)
	if err != nil {
		return nil, err
placeholder
	body, contentType, err = normalizeGrokMediaForwardBody(endpoint, body, contentType)
	if err != nil {
		return nil, err
placeholder

	var bodyReader io.Reader
	if endpoint.RequiresRequestBody() {
		bodyReader = bytes.NewReader(body)
placeholder
	upstreamCtx, releaseUpstreamCtx := detachUpstreamContext(ctx)
	defer releaseUpstreamCtx()
	upstreamReq, err := http.NewRequestWithContext(upstreamCtx, endpoint.httpMethod(), targetURL, bodyReader)
	if err != nil {
		return nil, err
placeholder
	upstreamReq.Header.Set("Authorization", "Bearer "+token)
	upstreamReq.Header.Set("Accept", "application/json")
	upstreamReq.Header.Set("User-Agent", "sub2api-grok/1.0")
	if endpoint.RequiresRequestBody() {
		contentType = strings.TrimSpace(contentType)
		if contentType == "" {
			contentType = "application/json"
	placeholder
		upstreamReq.Header.Set("Content-Type", contentType)
placeholder

	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
placeholder
	upstreamStart := time.Now()
	resp, err := s.httpUpstream.Do(upstreamReq, proxyURL, account.ID, account.Concurrency)
	SetOpsLatencyMs(c, OpsUpstreamLatencyMsKey, time.Since(upstreamStart).Milliseconds())
	if err != nil {
		return nil, s.handleOpenAIUpstreamTransportError(ctx, c, account, err, false)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	requestIDHeader := firstNonEmpty(resp.Header.Get("x-request-id"), resp.Header.Get("xai-request-id"))
	requestInfo := ParseGrokMediaRequest(contentType, body)
	requestModel := requestInfo.Model
	if resp.StatusCode >= 400 {
		s.updateGrokUsageSnapshot(ctx, account.ID, xai.ParseQuotaHeaders(resp.Header, resp.StatusCode))
		return s.handleGrokMediaErrorResponse(ctx, resp, c, account, requestIDHeader, requestModel)
placeholder

	s.updateGrokUsageSnapshot(ctx, account.ID, xai.ParseQuotaHeaders(resp.Header, resp.StatusCode))
	respBody, err := ReadUpstreamResponseBody(resp.Body, s.cfg, c, openAITooLargeError)
	if err != nil {
		return nil, err
placeholder
	writeGrokMediaResponse(c, resp, respBody, s.responseHeaderFilter)
	usage := grokMediaUsageFromResponse(endpoint, requestInfo, respBody)
	return &OpenAIForwardResult{
		RequestID:        requestIDHeader,
		ResponseID:       usage.ResponseID,
		Usage:            usage.Usage,
		Model:            requestModel,
		BillingModel:     requestModel,
		UpstreamModel:    requestModel,
		ResponseHeaders:  resp.Header.Clone(),
		Duration:         time.Since(startTime),
		ImageCount:       usage.ImageCount,
		ImageSize:        usage.ImageSize,
		ImageInputSize:   usage.ImageInputSize,
		ImageOutputSizes: usage.ImageOutputSizes,
placeholder, nil
placeholder

func prepareGrokMediaForwardBody(endpoint GrokMediaEndpoint, body []byte, contentType string) ([]byte, string, error) {
	if endpoint != GrokMediaEndpointImagesEdits || gjson.ValidBytes(body) {
		return body, contentType, nil
placeholder
	mediaType, _, err := mime.ParseMediaType(strings.TrimSpace(contentType))
	if err != nil || !strings.EqualFold(mediaType, "multipart/form-data") {
		return body, contentType, nil
placeholder

	info := ParseGrokMediaRequest(contentType, body)
	payload := make(map[string]any)
	if info.Model != "" {
		payload["model"] = info.Model
placeholder
	if info.Prompt != "" {
		payload["prompt"] = info.Prompt
placeholder
	if info.N > 1 {
		payload["n"] = info.N
placeholder
	if info.Size != "" {
		payload["size"] = info.Size
placeholder

	images := make([]map[string]string, 0, len(info.InputImageURLs)+len(info.Uploads))
	for _, imageURL := range info.InputImageURLs {
		if imageURL = strings.TrimSpace(imageURL); imageURL != "" {
			images = append(images, map[string]string{"image_url": imageURLplaceholder)
	placeholder
placeholder
	for _, upload := range info.Uploads {
		dataURL, err := openAIImageUploadToDataURL(upload)
		if err != nil {
			return nil, "", err
	placeholder
		images = append(images, map[string]string{"image_url": dataURLplaceholder)
placeholder
	if len(images) > 0 {
		payload["image"] = images[0]
		if len(images) > 1 {
			payload["images"] = images
	placeholder
placeholder

	maskImageURL := strings.TrimSpace(info.MaskImageURL)
	if info.MaskUpload != nil {
		dataURL, err := openAIImageUploadToDataURL(*info.MaskUpload)
		if err != nil {
			return nil, "", err
	placeholder
		maskImageURL = dataURL
placeholder
	if maskImageURL != "" {
		payload["mask"] = map[string]string{"image_url": maskImageURLplaceholder
placeholder

	out, err := marshalOpenAIUpstreamJSON(payload)
	if err != nil {
		return nil, "", err
placeholder
	return out, "application/json", nil
placeholder

func normalizeGrokMediaForwardBody(endpoint GrokMediaEndpoint, body []byte, contentType string) ([]byte, string, error) {
	if !endpoint.RequiresRequestBody() || !gjson.ValidBytes(body) {
		return body, contentType, nil
placeholder
	model := strings.TrimSpace(gjson.GetBytes(body, "model").String())
	upstreamModel := normalizeGrokMediaModelForEndpoint(endpoint, model)
	if upstreamModel == "" || upstreamModel == model {
		return body, contentType, nil
placeholder
	out, err := sjson.SetBytes(body, "model", upstreamModel)
	if err != nil {
		return nil, "", fmt.Errorf("rewrite grok media model: %w", err)
placeholder
	return out, contentType, nil
placeholder

func normalizeGrokMediaModelForEndpoint(endpoint GrokMediaEndpoint, model string) string {
	model = strings.TrimSpace(model)
	switch endpoint {
	case GrokMediaEndpointImagesGenerations, GrokMediaEndpointImagesEdits:
		if model == "grok-imagine" {
			return "grok-imagine-image-quality"
	placeholder
placeholder
	return model
placeholder

type grokMediaUsageMetadata struct {
	ResponseID       string
	Usage            OpenAIUsage
	ImageCount       int
	ImageSize        string
	ImageInputSize   string
	ImageOutputSizes []string
placeholder

func grokMediaUsageFromResponse(endpoint GrokMediaEndpoint, requestInfo GrokMediaRequestInfo, responseBody []byte) grokMediaUsageMetadata {
	usage, _ := extractOpenAIUsageFromJSONBytes(responseBody)
	meta := grokMediaUsageMetadata{Usage: usageplaceholder
	switch endpoint {
	case GrokMediaEndpointImagesGenerations, GrokMediaEndpointImagesEdits:
		imageCount := countOpenAIResponseImageOutputsFromJSONBytes(responseBody)
		if imageCount <= 0 {
			imageCount = requestInfo.N
	placeholder
		if imageCount <= 0 {
			imageCount = 1
	placeholder
		meta.ImageCount = imageCount
		meta.ImageSize = requestInfo.SizeTier
		meta.ImageInputSize = requestInfo.Size
		meta.ImageOutputSizes = collectOpenAIResponseImageOutputSizesFromJSONBytes(responseBody)
	case GrokMediaEndpointVideosGenerations:
		meta.ResponseID = extractGrokMediaVideoRequestID(responseBody)
		// Video generation is one billable media unit; the legacy usage schema stores it in ImageCount.
		meta.ImageCount = 1
		meta.ImageSize = requestInfo.SizeTier
		meta.ImageInputSize = requestInfo.Size
placeholder
	return meta
placeholder

func extractGrokMediaVideoRequestID(body []byte) string {
	if len(body) == 0 || !gjson.ValidBytes(body) {
		return ""
placeholder
	for _, path := range []string{"request_id", "id", "data.request_id", "data.id", "video.request_id", "video.id"placeholder {
		if id := strings.TrimSpace(gjson.GetBytes(body, path).String()); id != "" {
			return id
	placeholder
placeholder
	return ""
placeholder

func (s *OpenAIGatewayService) handleGrokMediaErrorResponse(
	ctx context.Context,
	resp *http.Response,
	c *gin.Context,
	account *Account,
	requestIDHeader string,
	requestedModel string,
) (*OpenAIForwardResult, error) {
	body := s.readUpstreamErrorBody(resp)
	upstreamMsg := sanitizeUpstreamErrorMessage(strings.TrimSpace(extractUpstreamErrorMessage(body)))
	if upstreamMsg == "" {
		upstreamMsg = fmt.Sprintf("xAI upstream returned status %d", resp.StatusCode)
placeholder

	upstreamDetail := ""
	if s.cfg != nil && s.cfg.Gateway.LogUpstreamErrorBody {
		maxBytes := s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes
		if maxBytes <= 0 {
			maxBytes = 2048
	placeholder
		upstreamDetail = truncateString(string(body), maxBytes)
placeholder
	setOpsUpstreamError(c, resp.StatusCode, upstreamMsg, upstreamDetail)

	if status, errType, errMsg, matched := applyErrorPassthroughRule(
		c,
		account.Platform,
		resp.StatusCode,
		body,
		http.StatusBadGateway,
		"upstream_error",
		"Upstream request failed",
	); matched {
		MarkResponseCommitted(c)
		writeGrokMediaErrorResponse(c, status, errType, errMsg)
		return nil, fmt.Errorf("upstream error: %d (passthrough rule matched) message=%s", resp.StatusCode, upstreamMsg)
placeholder

	if !account.ShouldHandleErrorCode(resp.StatusCode) {
		appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
			Platform:           account.Platform,
			AccountID:          account.ID,
			AccountName:        account.Name,
			UpstreamStatusCode: resp.StatusCode,
			UpstreamRequestID:  requestIDHeader,
			Kind:               "http_error",
			Message:            upstreamMsg,
			Detail:             upstreamDetail,
	placeholder)
		MarkResponseCommitted(c)
		writeGrokMediaErrorResponse(c, http.StatusInternalServerError, "upstream_error", "Upstream gateway error")
		return nil, fmt.Errorf("upstream error: %d (not in custom error codes) message=%s", resp.StatusCode, upstreamMsg)
placeholder

	s.handleGrokAccountUpstreamError(ctx, account, resp.StatusCode, resp.Header, body)
	kind := "http_error"
	if s.shouldFailoverUpstreamError(resp.StatusCode) {
		kind = "failover"
placeholder
	appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
		Platform:           account.Platform,
		AccountID:          account.ID,
		AccountName:        account.Name,
		UpstreamStatusCode: resp.StatusCode,
		UpstreamRequestID:  requestIDHeader,
		Kind:               kind,
		Message:            upstreamMsg,
		Detail:             upstreamDetail,
placeholder)
	if kind == "failover" {
		return nil, &UpstreamFailoverError{
			StatusCode:             resp.StatusCode,
			ResponseBody:           body,
			RetryableOnSameAccount: account.IsPoolMode() && account.IsPoolModeRetryableStatus(resp.StatusCode),
	placeholder
placeholder

	MarkResponseCommitted(c)
	writeGrokMediaErrorResponse(c, resp.StatusCode, grokMediaErrorType(resp.StatusCode), upstreamMsg)
	return nil, fmt.Errorf("upstream error: %d %s", resp.StatusCode, upstreamMsg)
placeholder

func grokMediaErrorType(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "invalid_request_error"
	case http.StatusNotFound:
		return "not_found_error"
	case http.StatusTooManyRequests:
		return "rate_limit_error"
	default:
		return "upstream_error"
placeholder
placeholder

func writeGrokMediaErrorResponse(c *gin.Context, statusCode int, errType, message string) {
	if c == nil || c.Writer == nil || c.Writer.Written() {
		return
placeholder
	c.JSON(statusCode, gin.H{
		"error": gin.H{
			"type":    strings.TrimSpace(errType),
			"message": strings.TrimSpace(message),
	placeholder,
placeholder)
placeholder

func writeGrokMediaResponse(c *gin.Context, resp *http.Response, body []byte, filter *responseheaders.CompiledHeaderFilter) {
	if c == nil || resp == nil {
		return
placeholder
	writeOpenAIPassthroughResponseHeaders(c.Writer.Header(), resp.Header, filter)
	contentType := strings.TrimSpace(resp.Header.Get("Content-Type"))
	if contentType == "" {
		contentType = "application/json"
placeholder
	c.Data(resp.StatusCode, contentType, body)
placeholder
