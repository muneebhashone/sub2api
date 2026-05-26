package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/Wei-Shaw/sub2api/internal/util/responseheaders"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type openAIResponsesImageResult struct {
	Result        string
	RevisedPrompt string
	OutputFormat  string
	Size          string
	Background    string
	Quality       string
	Model         string
placeholder

type OpenAIImagesUpstreamError struct {
	StatusCode        int
	ErrorType         string
	Code              string
	Message           string
	Param             string
	UpstreamRequestID string
placeholder

func (e *OpenAIImagesUpstreamError) Error() string {
	if e == nil {
		return ""
placeholder
	code := strings.TrimSpace(e.Code)
	if code == "" {
		code = strings.TrimSpace(e.ErrorType)
placeholder
	message := strings.TrimSpace(e.Message)
	if code != "" && message != "" {
		return fmt.Sprintf("openai images upstream error: %s: %s", code, message)
placeholder
	if message != "" {
		return "openai images upstream error: " + message
placeholder
	if code != "" {
		return "openai images upstream error: " + code
placeholder
	return "openai images upstream error"
placeholder

func (e *OpenAIImagesUpstreamError) clientStatusCode() int {
	if e == nil {
		return http.StatusBadGateway
placeholder
	if e.StatusCode > 0 {
		return e.StatusCode
placeholder
	return http.StatusBadGateway
placeholder

func (e *OpenAIImagesUpstreamError) clientErrorType() string {
	if e == nil {
		return "upstream_error"
placeholder
	if trimmed := strings.TrimSpace(e.ErrorType); trimmed != "" {
		return trimmed
placeholder
	return "upstream_error"
placeholder

func (e *OpenAIImagesUpstreamError) clientMessage() string {
	if e == nil {
		return "Upstream request failed"
placeholder
	if trimmed := strings.TrimSpace(e.Message); trimmed != "" {
		return trimmed
placeholder
	if trimmed := strings.TrimSpace(e.Code); trimmed != "" {
		return trimmed
placeholder
	return "Upstream request failed"
placeholder

func openAIResponsesImageResultKey(itemID string, result openAIResponsesImageResult) string {
	if strings.TrimSpace(result.Result) != "" {
		return strings.TrimSpace(result.OutputFormat) + "|" + strings.TrimSpace(result.Result)
placeholder
	return "item:" + strings.TrimSpace(itemID)
placeholder

func appendOpenAIResponsesImageResultDedup(results *[]openAIResponsesImageResult, seen map[string]struct{placeholder, itemID string, result openAIResponsesImageResult) bool {
	if results == nil {
		return false
placeholder
	key := openAIResponsesImageResultKey(itemID, result)
	if key != "" {
		if _, exists := seen[key]; exists {
			return false
	placeholder
		seen[key] = struct{placeholder{placeholder
placeholder
	*results = append(*results, result)
	return true
placeholder

func mergeOpenAIResponsesImageMeta(dst *openAIResponsesImageResult, src openAIResponsesImageResult) {
	if dst == nil {
		return
placeholder
	if trimmed := strings.TrimSpace(src.OutputFormat); trimmed != "" {
		dst.OutputFormat = trimmed
placeholder
	if trimmed := strings.TrimSpace(src.Size); trimmed != "" {
		dst.Size = trimmed
placeholder
	if trimmed := strings.TrimSpace(src.Background); trimmed != "" {
		dst.Background = trimmed
placeholder
	if trimmed := strings.TrimSpace(src.Quality); trimmed != "" {
		dst.Quality = trimmed
placeholder
	if trimmed := strings.TrimSpace(src.Model); trimmed != "" {
		dst.Model = trimmed
placeholder
placeholder

func openAIResponsesImageResultSizes(results []openAIResponsesImageResult) []string {
	if len(results) == 0 {
		return nil
placeholder
	sizes := make([]string, 0, len(results))
	for _, result := range results {
		if size := strings.TrimSpace(result.Size); size != "" {
			sizes = append(sizes, size)
	placeholder
placeholder
	if len(sizes) == 0 {
		return nil
placeholder
	return sizes
placeholder

func extractOpenAIResponsesImageMetaFromLifecycleEvent(payload []byte) (openAIResponsesImageResult, int64, bool) {
	switch gjson.GetBytes(payload, "type").String() {
	case "response.created", "response.in_progress", "response.completed":
	default:
		return openAIResponsesImageResult{placeholder, 0, false
placeholder

	response := gjson.GetBytes(payload, "response")
	if !response.Exists() {
		return openAIResponsesImageResult{placeholder, 0, false
placeholder

	meta := openAIResponsesImageResult{
		OutputFormat: strings.TrimSpace(response.Get("tools.0.output_format").String()),
		Size:         strings.TrimSpace(response.Get("tools.0.size").String()),
		Background:   strings.TrimSpace(response.Get("tools.0.background").String()),
		Quality:      strings.TrimSpace(response.Get("tools.0.quality").String()),
		Model:        strings.TrimSpace(response.Get("tools.0.model").String()),
placeholder
	return meta, response.Get("created_at").Int(), true
placeholder

func buildOpenAIImagesStreamPartialPayload(
	eventType string,
	b64 string,
	partialImageIndex int64,
	responseFormat string,
	createdAt int64,
	meta openAIResponsesImageResult,
) []byte {
	if createdAt <= 0 {
		createdAt = time.Now().Unix()
placeholder

	payload := []byte(`{"type":"","created_at":0,"partial_image_index":0,"b64_json":""placeholder`)
	payload, _ = sjson.SetBytes(payload, "type", eventType)
	payload, _ = sjson.SetBytes(payload, "created_at", createdAt)
	payload, _ = sjson.SetBytes(payload, "partial_image_index", partialImageIndex)
	payload, _ = sjson.SetBytes(payload, "b64_json", b64)
	if strings.EqualFold(strings.TrimSpace(responseFormat), "url") {
		payload, _ = sjson.SetBytes(payload, "url", "data:"+openAIImageOutputMIMEType(meta.OutputFormat)+";base64,"+b64)
placeholder
	if meta.Background != "" {
		payload, _ = sjson.SetBytes(payload, "background", meta.Background)
placeholder
	if meta.OutputFormat != "" {
		payload, _ = sjson.SetBytes(payload, "output_format", meta.OutputFormat)
placeholder
	if meta.Quality != "" {
		payload, _ = sjson.SetBytes(payload, "quality", meta.Quality)
placeholder
	if meta.Size != "" {
		payload, _ = sjson.SetBytes(payload, "size", meta.Size)
placeholder
	if meta.Model != "" {
		payload, _ = sjson.SetBytes(payload, "model", meta.Model)
placeholder
	return payload
placeholder

func buildOpenAIImagesStreamCompletedPayload(
	eventType string,
	img openAIResponsesImageResult,
	responseFormat string,
	createdAt int64,
	usageRaw []byte,
) []byte {
	if createdAt <= 0 {
		createdAt = time.Now().Unix()
placeholder

	payload := []byte(`{"type":"","created_at":0,"b64_json":""placeholder`)
	payload, _ = sjson.SetBytes(payload, "type", eventType)
	payload, _ = sjson.SetBytes(payload, "created_at", createdAt)
	payload, _ = sjson.SetBytes(payload, "b64_json", img.Result)
	if strings.EqualFold(strings.TrimSpace(responseFormat), "url") {
		payload, _ = sjson.SetBytes(payload, "url", "data:"+openAIImageOutputMIMEType(img.OutputFormat)+";base64,"+img.Result)
placeholder
	if img.Background != "" {
		payload, _ = sjson.SetBytes(payload, "background", img.Background)
placeholder
	if img.OutputFormat != "" {
		payload, _ = sjson.SetBytes(payload, "output_format", img.OutputFormat)
placeholder
	if img.Quality != "" {
		payload, _ = sjson.SetBytes(payload, "quality", img.Quality)
placeholder
	if img.Size != "" {
		payload, _ = sjson.SetBytes(payload, "size", img.Size)
placeholder
	if img.Model != "" {
		payload, _ = sjson.SetBytes(payload, "model", img.Model)
placeholder
	if len(usageRaw) > 0 && gjson.ValidBytes(usageRaw) {
		payload, _ = sjson.SetRawBytes(payload, "usage", usageRaw)
placeholder
	return payload
placeholder

func openAIImageOutputMIMEType(outputFormat string) string {
	if outputFormat == "" {
		return "image/png"
placeholder
	if strings.Contains(outputFormat, "/") {
		return outputFormat
placeholder
	switch strings.ToLower(strings.TrimSpace(outputFormat)) {
	case "png":
		return "image/png"
	case "jpg", "jpeg":
		return "image/jpeg"
	case "webp":
		return "image/webp"
	default:
		return "image/png"
placeholder
placeholder

func openAIImageUploadToDataURL(upload OpenAIImagesUpload) (string, error) {
	if len(upload.Data) == 0 {
		return "", fmt.Errorf("upload %q is empty", strings.TrimSpace(upload.FileName))
placeholder
	contentType := strings.TrimSpace(upload.ContentType)
	if contentType == "" {
		contentType = http.DetectContentType(upload.Data)
placeholder
	return "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(upload.Data), nil
placeholder

func buildOpenAIImagesResponsesRequest(parsed *OpenAIImagesRequest, toolModel string) ([]byte, error) {
	if parsed == nil {
		return nil, fmt.Errorf("parsed images request is required")
placeholder
	prompt := strings.TrimSpace(parsed.Prompt)
	if prompt == "" {
		return nil, fmt.Errorf("prompt is required")
placeholder

	inputImages := make([]string, 0, len(parsed.InputImageURLs)+len(parsed.Uploads))
	for _, imageURL := range parsed.InputImageURLs {
		if trimmed := strings.TrimSpace(imageURL); trimmed != "" {
			inputImages = append(inputImages, trimmed)
	placeholder
placeholder
	for _, upload := range parsed.Uploads {
		dataURL, err := openAIImageUploadToDataURL(upload)
		if err != nil {
			return nil, err
	placeholder
		inputImages = append(inputImages, dataURL)
placeholder
	if parsed.IsEdits() && len(inputImages) == 0 {
		return nil, fmt.Errorf("image input is required")
placeholder

	req := []byte(`{"instructions":"","stream":true,"reasoning":{"effort":"medium","summary":"auto"placeholder,"parallel_tool_calls":true,"include":["reasoning.encrypted_content"],"model":"","store":false,"tool_choice":{"type":"image_generation"placeholderplaceholder`)
	req, _ = sjson.SetBytes(req, "model", openAIImagesResponsesMainModel)

	input := []byte(`[{"type":"message","role":"user","content":[{"type":"input_text","text":""placeholder]placeholder]`)
	input, _ = sjson.SetBytes(input, "0.content.0.text", prompt)
	for index, imageURL := range inputImages {
		part := []byte(`{"type":"input_image","image_url":""placeholder`)
		part, _ = sjson.SetBytes(part, "image_url", imageURL)
		input, _ = sjson.SetRawBytes(input, fmt.Sprintf("0.content.%d", index+1), part)
placeholder
	req, _ = sjson.SetRawBytes(req, "input", input)

	action := "generate"
	if parsed.IsEdits() {
		action = "edit"
placeholder
	tool := []byte(`{"type":"image_generation","action":"","model":""placeholder`)
	tool, _ = sjson.SetBytes(tool, "action", action)
	tool, _ = sjson.SetBytes(tool, "model", strings.TrimSpace(toolModel))
	if shouldPassOpenAIImagesN(toolModel, parsed.N) {
		tool, _ = sjson.SetBytes(tool, "n", parsed.N)
placeholder

	for _, field := range []struct {
		path  string
		value string
placeholder{
		{path: "size", value: parsed.Sizeplaceholder,
		{path: "quality", value: parsed.Qualityplaceholder,
		{path: "background", value: parsed.Backgroundplaceholder,
		{path: "output_format", value: parsed.OutputFormatplaceholder,
		{path: "moderation", value: parsed.Moderationplaceholder,
		{path: "style", value: parsed.Styleplaceholder,
placeholder {
		if trimmed := strings.TrimSpace(field.value); trimmed != "" {
			tool, _ = sjson.SetBytes(tool, field.path, trimmed)
	placeholder
placeholder
	if parsed.OutputCompression != nil {
		tool, _ = sjson.SetBytes(tool, "output_compression", *parsed.OutputCompression)
placeholder
	if parsed.PartialImages != nil {
		tool, _ = sjson.SetBytes(tool, "partial_images", *parsed.PartialImages)
placeholder

	maskImageURL := strings.TrimSpace(parsed.MaskImageURL)
	if parsed.MaskUpload != nil {
		dataURL, err := openAIImageUploadToDataURL(*parsed.MaskUpload)
		if err != nil {
			return nil, err
	placeholder
		maskImageURL = dataURL
placeholder
	if maskImageURL != "" {
		tool, _ = sjson.SetBytes(tool, "input_image_mask.image_url", maskImageURL)
placeholder

	req, _ = sjson.SetRawBytes(req, "tools", []byte(`[]`))
	req, _ = sjson.SetRawBytes(req, "tools.-1", tool)
	return req, nil
placeholder

func shouldPassOpenAIImagesN(model string, n int) bool {
	if n <= 1 {
		return false
placeholder
	return !strings.EqualFold(strings.TrimSpace(model), "dall-e-3")
placeholder

func extractOpenAIImagesFromResponsesCompleted(payload []byte) ([]openAIResponsesImageResult, int64, []byte, openAIResponsesImageResult, error) {
	if gjson.GetBytes(payload, "type").String() != "response.completed" {
		return nil, 0, nil, openAIResponsesImageResult{placeholder, fmt.Errorf("unexpected event type")
placeholder

	createdAt := gjson.GetBytes(payload, "response.created_at").Int()
	if createdAt <= 0 {
		createdAt = time.Now().Unix()
placeholder

	var (
		results   []openAIResponsesImageResult
		firstMeta openAIResponsesImageResult
	)
	output := gjson.GetBytes(payload, "response.output")
	if output.IsArray() {
		for _, item := range output.Array() {
			if item.Get("type").String() != "image_generation_call" {
				continue
		placeholder
			result := strings.TrimSpace(item.Get("result").String())
			if result == "" {
				continue
		placeholder
			entry := openAIResponsesImageResult{
				Result:        result,
				RevisedPrompt: strings.TrimSpace(item.Get("revised_prompt").String()),
				OutputFormat:  strings.TrimSpace(item.Get("output_format").String()),
				Size:          strings.TrimSpace(item.Get("size").String()),
				Background:    strings.TrimSpace(item.Get("background").String()),
				Quality:       strings.TrimSpace(item.Get("quality").String()),
		placeholder
			if len(results) == 0 {
				firstMeta = entry
		placeholder
			results = append(results, entry)
	placeholder
placeholder

	var usageRaw []byte
	if usage := gjson.GetBytes(payload, "response.tool_usage.image_gen"); usage.Exists() && usage.IsObject() {
		usageRaw = []byte(usage.Raw)
placeholder
	return results, createdAt, usageRaw, firstMeta, nil
placeholder

func extractOpenAIImageFromResponsesOutputItemDone(payload []byte) (openAIResponsesImageResult, string, bool, error) {
	if gjson.GetBytes(payload, "type").String() != "response.output_item.done" {
		return openAIResponsesImageResult{placeholder, "", false, fmt.Errorf("unexpected event type")
placeholder

	item := gjson.GetBytes(payload, "item")
	if !item.Exists() || item.Get("type").String() != "image_generation_call" {
		return openAIResponsesImageResult{placeholder, "", false, nil
placeholder

	result := strings.TrimSpace(item.Get("result").String())
	if result == "" {
		return openAIResponsesImageResult{placeholder, "", false, nil
placeholder

	entry := openAIResponsesImageResult{
		Result:        result,
		RevisedPrompt: strings.TrimSpace(item.Get("revised_prompt").String()),
		OutputFormat:  strings.TrimSpace(item.Get("output_format").String()),
		Size:          strings.TrimSpace(item.Get("size").String()),
		Background:    strings.TrimSpace(item.Get("background").String()),
		Quality:       strings.TrimSpace(item.Get("quality").String()),
placeholder
	return entry, strings.TrimSpace(item.Get("id").String()), true, nil
placeholder

func collectOpenAIImagesFromResponsesBody(body []byte) ([]openAIResponsesImageResult, int64, []byte, openAIResponsesImageResult, bool, error) {
	var (
		fallbackResults []openAIResponsesImageResult
		fallbackSeen    = make(map[string]struct{placeholder)
		finalResults    []openAIResponsesImageResult
		finalMeta       openAIResponsesImageResult
		collectErr      error
		createdAt       int64
		usageRaw        []byte
		foundFinal      bool
		responseMeta    openAIResponsesImageResult
	)

	forEachOpenAISSEDataPayload(string(body), func(payload []byte) {
		if collectErr != nil || len(finalResults) > 0 {
			return
	placeholder
		if !gjson.ValidBytes(payload) {
			return
	placeholder
		if meta, eventCreatedAt, ok := extractOpenAIResponsesImageMetaFromLifecycleEvent(payload); ok {
			mergeOpenAIResponsesImageMeta(&responseMeta, meta)
			if eventCreatedAt > 0 {
				createdAt = eventCreatedAt
		placeholder
	placeholder

		switch gjson.GetBytes(payload, "type").String() {
		case "response.output_item.done":
			result, itemID, ok, err := extractOpenAIImageFromResponsesOutputItemDone(payload)
			if err != nil {
				collectErr = err
				return
		placeholder
			if ok {
				mergeOpenAIResponsesImageMeta(&result, responseMeta)
				appendOpenAIResponsesImageResultDedup(&fallbackResults, fallbackSeen, itemID, result)
		placeholder
		case "response.completed":
			results, completedAt, completedUsageRaw, firstMeta, err := extractOpenAIImagesFromResponsesCompleted(payload)
			if err != nil {
				collectErr = err
				return
		placeholder
			foundFinal = true
			if completedAt > 0 {
				createdAt = completedAt
		placeholder
			if len(completedUsageRaw) > 0 {
				usageRaw = completedUsageRaw
		placeholder
			if len(results) > 0 {
				mergeOpenAIResponsesImageMeta(&firstMeta, responseMeta)
				finalResults = results
				finalMeta = firstMeta
				return
		placeholder
			if len(fallbackResults) > 0 {
				firstMeta = fallbackResults[0]
				mergeOpenAIResponsesImageMeta(&firstMeta, responseMeta)
				finalResults = fallbackResults
				finalMeta = firstMeta
				return
		placeholder
	placeholder
placeholder)
	if collectErr != nil {
		return nil, 0, nil, openAIResponsesImageResult{placeholder, false, collectErr
placeholder
	if len(finalResults) > 0 {
		return finalResults, createdAt, usageRaw, finalMeta, true, nil
placeholder

	if len(fallbackResults) > 0 {
		firstMeta := fallbackResults[0]
		mergeOpenAIResponsesImageMeta(&firstMeta, responseMeta)
		return fallbackResults, createdAt, usageRaw, firstMeta, foundFinal, nil
placeholder
	return nil, createdAt, usageRaw, openAIResponsesImageResult{placeholder, foundFinal, nil
placeholder

func extractOpenAIImagesUpstreamError(body []byte) *OpenAIImagesUpstreamError {
	var upstreamErr *OpenAIImagesUpstreamError
	forEachOpenAISSEDataPayload(string(body), func(payload []byte) {
		if upstreamErr != nil || !gjson.ValidBytes(payload) {
			return
	placeholder
		upstreamErr = openAIImagesUpstreamErrorFromSSEPayload(payload)
placeholder)
	return upstreamErr
placeholder

func openAIImagesUpstreamErrorFromSSEPayload(payload []byte) *OpenAIImagesUpstreamError {
	if !gjson.ValidBytes(payload) {
		return nil
placeholder
	switch gjson.GetBytes(payload, "type").String() {
	case "error":
		return openAIImagesUpstreamErrorFromGJSON(gjson.GetBytes(payload, "error"), "")
	case "response.failed":
		response := gjson.GetBytes(payload, "response")
		return openAIImagesUpstreamErrorFromGJSON(response.Get("error"), response.Get("id").String())
	default:
		return nil
placeholder
placeholder

func openAIImagesUpstreamErrorFromGJSON(errorObj gjson.Result, upstreamRequestID string) *OpenAIImagesUpstreamError {
	if !errorObj.Exists() {
		return nil
placeholder
	code := strings.TrimSpace(errorObj.Get("code").String())
	errType := strings.TrimSpace(errorObj.Get("type").String())
	message := strings.TrimSpace(errorObj.Get("message").String())
	param := strings.TrimSpace(errorObj.Get("param").String())
	statusCode := http.StatusBadGateway
	if strings.EqualFold(code, "moderation_blocked") || strings.EqualFold(errType, "image_generation_user_error") {
		statusCode = http.StatusBadRequest
placeholder
	if message == "" {
		message = "Upstream request failed"
placeholder
	return &OpenAIImagesUpstreamError{
		StatusCode:        statusCode,
		ErrorType:         errType,
		Code:              code,
		Message:           sanitizeUpstreamErrorMessage(message),
		Param:             param,
		UpstreamRequestID: strings.TrimSpace(upstreamRequestID),
placeholder
placeholder

func buildOpenAIImagesAPIResponse(
	results []openAIResponsesImageResult,
	createdAt int64,
	usageRaw []byte,
	firstMeta openAIResponsesImageResult,
	responseFormat string,
) ([]byte, error) {
	if createdAt <= 0 {
		createdAt = time.Now().Unix()
placeholder
	out := []byte(`{"created":0,"data":[]placeholder`)
	out, _ = sjson.SetBytes(out, "created", createdAt)

	format := strings.ToLower(strings.TrimSpace(responseFormat))
	if format == "" {
		format = "b64_json"
placeholder
	for _, img := range results {
		item := []byte(`{placeholder`)
		if format == "url" {
			item, _ = sjson.SetBytes(item, "url", "data:"+openAIImageOutputMIMEType(img.OutputFormat)+";base64,"+img.Result)
	placeholder else {
			item, _ = sjson.SetBytes(item, "b64_json", img.Result)
	placeholder
		if img.RevisedPrompt != "" {
			item, _ = sjson.SetBytes(item, "revised_prompt", img.RevisedPrompt)
	placeholder
		out, _ = sjson.SetRawBytes(out, "data.-1", item)
placeholder
	if firstMeta.Background != "" {
		out, _ = sjson.SetBytes(out, "background", firstMeta.Background)
placeholder
	if firstMeta.OutputFormat != "" {
		out, _ = sjson.SetBytes(out, "output_format", firstMeta.OutputFormat)
placeholder
	if firstMeta.Quality != "" {
		out, _ = sjson.SetBytes(out, "quality", firstMeta.Quality)
placeholder
	if firstMeta.Size != "" {
		out, _ = sjson.SetBytes(out, "size", firstMeta.Size)
placeholder
	if firstMeta.Model != "" {
		out, _ = sjson.SetBytes(out, "model", firstMeta.Model)
placeholder
	if len(usageRaw) > 0 && gjson.ValidBytes(usageRaw) {
		out, _ = sjson.SetRawBytes(out, "usage", usageRaw)
placeholder
	return out, nil
placeholder

func openAIImagesStreamPrefix(parsed *OpenAIImagesRequest) string {
	if parsed != nil && parsed.IsEdits() {
		return "image_edit"
placeholder
	return "image_generation"
placeholder

func buildOpenAIImagesStreamErrorBody(message string) []byte {
	body := []byte(`{"type":"error","error":{"type":"upstream_error","message":""placeholderplaceholder`)
	if strings.TrimSpace(message) == "" {
		message = "upstream request failed"
placeholder
	body, _ = sjson.SetBytes(body, "error.message", message)
	return body
placeholder

func buildOpenAIImagesStreamErrorBodyFromUpstream(err *OpenAIImagesUpstreamError) []byte {
	if err == nil {
		return buildOpenAIImagesStreamErrorBody("")
placeholder
	body := buildOpenAIImagesStreamErrorBody(err.clientMessage())
	body, _ = sjson.SetBytes(body, "error.type", err.clientErrorType())
	if code := strings.TrimSpace(err.Code); code != "" {
		body, _ = sjson.SetBytes(body, "error.code", code)
placeholder
	if param := strings.TrimSpace(err.Param); param != "" {
		body, _ = sjson.SetBytes(body, "error.param", param)
placeholder
	return body
placeholder

func writeOpenAIImagesUpstreamErrorResponse(c *gin.Context, err *OpenAIImagesUpstreamError) bool {
	if c == nil || c.Writer == nil || c.Writer.Written() || err == nil {
		return false
placeholder
	errorObj := gin.H{
		"type":    err.clientErrorType(),
		"message": err.clientMessage(),
placeholder
	if code := strings.TrimSpace(err.Code); code != "" {
		errorObj["code"] = code
placeholder
	if param := strings.TrimSpace(err.Param); param != "" {
		errorObj["param"] = param
placeholder
	c.JSON(err.clientStatusCode(), gin.H{
		"error": errorObj,
placeholder)
	return true
placeholder

func (s *OpenAIGatewayService) writeOpenAIImagesStreamEvent(c *gin.Context, flusher http.Flusher, eventName string, payload []byte) error {
	if strings.TrimSpace(eventName) != "" {
		if _, err := fmt.Fprintf(c.Writer, "event: %s\n", eventName); err != nil {
			return err
	placeholder
placeholder
	if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", payload); err != nil {
		return err
placeholder
	flusher.Flush()
	return nil
placeholder

func (s *OpenAIGatewayService) tryWriteOpenAIImagesStreamEvent(
	c *gin.Context,
	flusher http.Flusher,
	clientDisconnected *bool,
	lastWriteAt *time.Time,
	eventName string,
	payload []byte,
) bool {
	if clientDisconnected != nil && *clientDisconnected {
		return false
placeholder
	if err := s.writeOpenAIImagesStreamEvent(c, flusher, eventName, payload); err != nil {
		if clientDisconnected != nil {
			*clientDisconnected = true
	placeholder
		logger.LegacyPrintf("service.openai_gateway", "[OpenAI] Images stream client disconnected, continue draining upstream for billing")
		return false
placeholder
	if lastWriteAt != nil {
		*lastWriteAt = time.Now()
placeholder
	return true
placeholder

func (s *OpenAIGatewayService) handleOpenAIImagesOAuthNonStreamingResponse(
	resp *http.Response,
	c *gin.Context,
	responseFormat string,
	fallbackModel string,
) (OpenAIUsage, int, []string, error) {
	body, err := ReadUpstreamResponseBody(resp.Body, s.cfg, c, openAITooLargeError)
	if err != nil {
		return OpenAIUsage{placeholder, 0, nil, err
placeholder

	var usage OpenAIUsage
	forEachOpenAISSEDataPayload(string(body), func(data []byte) {
		s.parseSSEUsageBytes(data, &usage)
placeholder)
	results, createdAt, usageRaw, firstMeta, _, err := collectOpenAIImagesFromResponsesBody(body)
	if err != nil {
		return OpenAIUsage{placeholder, 0, nil, err
placeholder
	if len(results) == 0 {
		if upstreamErr := extractOpenAIImagesUpstreamError(body); upstreamErr != nil {
			setOpsUpstreamError(c, upstreamErr.clientStatusCode(), upstreamErr.clientMessage(), "")
			writeOpenAIImagesUpstreamErrorResponse(c, upstreamErr)
			return OpenAIUsage{placeholder, 0, nil, upstreamErr
	placeholder
		return OpenAIUsage{placeholder, 0, nil, fmt.Errorf("upstream did not return image output")
placeholder
	if strings.TrimSpace(firstMeta.Model) == "" {
		firstMeta.Model = strings.TrimSpace(fallbackModel)
placeholder

	responseBody, err := buildOpenAIImagesAPIResponse(results, createdAt, usageRaw, firstMeta, responseFormat)
	if err != nil {
		return OpenAIUsage{placeholder, 0, nil, err
placeholder
	responseheaders.WriteFilteredHeaders(c.Writer.Header(), resp.Header, s.responseHeaderFilter)
	c.Data(resp.StatusCode, "application/json; charset=utf-8", responseBody)
	return usage, len(results), openAIResponsesImageResultSizes(results), nil
placeholder

func (s *OpenAIGatewayService) handleOpenAIImagesOAuthStreamingResponse(
	resp *http.Response,
	c *gin.Context,
	startTime time.Time,
	responseFormat string,
	streamPrefix string,
	fallbackModel string,
) (OpenAIUsage, int, []string, *int, error) {
	responseheaders.WriteFilteredHeaders(c.Writer.Header(), resp.Header, s.responseHeaderFilter)
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Status(resp.StatusCode)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return OpenAIUsage{placeholder, 0, nil, nil, fmt.Errorf("streaming is not supported by response writer")
placeholder

	format := strings.ToLower(strings.TrimSpace(responseFormat))
	if format == "" {
		format = "b64_json"
placeholder

	usage := OpenAIUsage{placeholder
	imageCount := 0
	var imageOutputSizes []string
	var firstTokenMs *int
	emitted := make(map[string]struct{placeholder)
	pendingResults := make([]openAIResponsesImageResult, 0, 1)
	pendingSeen := make(map[string]struct{placeholder)
	streamMeta := openAIResponsesImageResult{Model: strings.TrimSpace(fallbackModel)placeholder
	var createdAt int64
	clientDisconnected := false
	lastDownstreamWriteAt := time.Now()
	var sseData openAISSEDataAccumulator
	var processDataErr error
	processDataDone := false

	processData := func(dataBytes []byte) {
		if processDataDone || processDataErr != nil {
			return
	placeholder
		if firstTokenMs == nil {
			ms := int(time.Since(startTime).Milliseconds())
			firstTokenMs = &ms
	placeholder
		s.parseSSEUsageBytes(dataBytes, &usage)
		if !gjson.ValidBytes(dataBytes) {
			return
	placeholder
		if meta, eventCreatedAt, ok := extractOpenAIResponsesImageMetaFromLifecycleEvent(dataBytes); ok {
			mergeOpenAIResponsesImageMeta(&streamMeta, meta)
			if eventCreatedAt > 0 {
				createdAt = eventCreatedAt
		placeholder
	placeholder
		switch gjson.GetBytes(dataBytes, "type").String() {
		case "response.image_generation_call.partial_image":
			b64 := strings.TrimSpace(gjson.GetBytes(dataBytes, "partial_image_b64").String())
			if b64 == "" {
				return
		placeholder
			eventName := streamPrefix + ".partial_image"
			partialMeta := streamMeta
			mergeOpenAIResponsesImageMeta(&partialMeta, openAIResponsesImageResult{
				OutputFormat: strings.TrimSpace(gjson.GetBytes(dataBytes, "output_format").String()),
				Background:   strings.TrimSpace(gjson.GetBytes(dataBytes, "background").String()),
		placeholder)
			payload := buildOpenAIImagesStreamPartialPayload(
				eventName,
				b64,
				gjson.GetBytes(dataBytes, "partial_image_index").Int(),
				format,
				createdAt,
				partialMeta,
			)
			s.tryWriteOpenAIImagesStreamEvent(c, flusher, &clientDisconnected, &lastDownstreamWriteAt, eventName, payload)
		case "response.output_item.done":
			img, itemID, ok, extractErr := extractOpenAIImageFromResponsesOutputItemDone(dataBytes)
			if extractErr != nil {
				s.tryWriteOpenAIImagesStreamEvent(c, flusher, &clientDisconnected, &lastDownstreamWriteAt, "error", buildOpenAIImagesStreamErrorBody(extractErr.Error()))
				processDataErr = extractErr
				processDataDone = true
				return
		placeholder
			if !ok {
				return
		placeholder
			mergeOpenAIResponsesImageMeta(&streamMeta, img)
			mergeOpenAIResponsesImageMeta(&img, streamMeta)
			key := openAIResponsesImageResultKey(itemID, img)
			if _, exists := emitted[key]; exists {
				return
		placeholder
			if _, exists := pendingSeen[key]; exists {
				return
		placeholder
			pendingSeen[key] = struct{placeholder{placeholder
			pendingResults = append(pendingResults, img)
		case "response.completed":
			results, _, usageRaw, firstMeta, extractErr := extractOpenAIImagesFromResponsesCompleted(dataBytes)
			if extractErr != nil {
				s.tryWriteOpenAIImagesStreamEvent(c, flusher, &clientDisconnected, &lastDownstreamWriteAt, "error", buildOpenAIImagesStreamErrorBody(extractErr.Error()))
				processDataErr = extractErr
				processDataDone = true
				return
		placeholder
			mergeOpenAIResponsesImageMeta(&streamMeta, firstMeta)
			finalResults := make([]openAIResponsesImageResult, 0, len(results)+len(pendingResults))
			finalSeen := make(map[string]struct{placeholder)
			for _, img := range results {
				mergeOpenAIResponsesImageMeta(&img, streamMeta)
				appendOpenAIResponsesImageResultDedup(&finalResults, finalSeen, "", img)
		placeholder
			for _, img := range pendingResults {
				mergeOpenAIResponsesImageMeta(&img, streamMeta)
				appendOpenAIResponsesImageResultDedup(&finalResults, finalSeen, "", img)
		placeholder
			if len(finalResults) == 0 {
				outputErr := fmt.Errorf("upstream did not return image output")
				s.tryWriteOpenAIImagesStreamEvent(c, flusher, &clientDisconnected, &lastDownstreamWriteAt, "error", buildOpenAIImagesStreamErrorBody(outputErr.Error()))
				processDataErr = outputErr
				processDataDone = true
				return
		placeholder
			eventName := streamPrefix + ".completed"
			for _, img := range finalResults {
				key := openAIResponsesImageResultKey("", img)
				if _, exists := emitted[key]; exists {
					continue
			placeholder
				payload := buildOpenAIImagesStreamCompletedPayload(eventName, img, format, createdAt, usageRaw)
				emitted[key] = struct{placeholder{placeholder
				s.tryWriteOpenAIImagesStreamEvent(c, flusher, &clientDisconnected, &lastDownstreamWriteAt, eventName, payload)
		placeholder
			imageCount = len(emitted)
			imageOutputSizes = openAIResponsesImageResultSizes(finalResults)
			processDataDone = true
		case "error", "response.failed":
			if upstreamErr := openAIImagesUpstreamErrorFromSSEPayload(dataBytes); upstreamErr != nil {
				if !clientDisconnected {
					s.tryWriteOpenAIImagesStreamEvent(c, flusher, &clientDisconnected, &lastDownstreamWriteAt, "error", buildOpenAIImagesStreamErrorBodyFromUpstream(upstreamErr))
			placeholder
				setOpsUpstreamError(c, upstreamErr.clientStatusCode(), upstreamErr.clientMessage(), "")
				processDataErr = upstreamErr
				processDataDone = true
				return
		placeholder
	placeholder
placeholder

	processLine := func(line []byte) (bool, error) {
		if len(line) == 0 {
			return false, nil
	placeholder
		sseData.AddLine(string(line), processData)
		if processDataErr != nil {
			return true, processDataErr
	placeholder
		return processDataDone, nil
placeholder

	flushData := func() (bool, error) {
		sseData.Flush(processData)
		if processDataErr != nil {
			return true, processDataErr
	placeholder
		return processDataDone, nil
placeholder

	finalizePending := func() error {
		if imageCount > 0 {
			return nil
	placeholder
		if len(pendingResults) > 0 {
			eventName := streamPrefix + ".completed"
			for _, img := range pendingResults {
				mergeOpenAIResponsesImageMeta(&img, streamMeta)
				key := openAIResponsesImageResultKey("", img)
				if _, exists := emitted[key]; exists {
					continue
			placeholder
				payload := buildOpenAIImagesStreamCompletedPayload(eventName, img, format, createdAt, nil)
				emitted[key] = struct{placeholder{placeholder
				s.tryWriteOpenAIImagesStreamEvent(c, flusher, &clientDisconnected, &lastDownstreamWriteAt, eventName, payload)
		placeholder
			imageCount = len(emitted)
			imageOutputSizes = openAIResponsesImageResultSizes(pendingResults)
			return nil
	placeholder

		streamErr := fmt.Errorf("stream disconnected before image generation completed")
		s.tryWriteOpenAIImagesStreamEvent(c, flusher, &clientDisconnected, &lastDownstreamWriteAt, "error", buildOpenAIImagesStreamErrorBody(streamErr.Error()))
		return streamErr
placeholder

	streamInterval := s.openAIImageStreamDataInterval()
	keepaliveInterval := s.openAIImageStreamKeepaliveInterval()
	if streamInterval <= 0 && keepaliveInterval <= 0 {
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadBytes('\n')
			done, processErr := processLine(line)
			if processErr != nil {
				return usage, imageCount, imageOutputSizes, firstTokenMs, processErr
		placeholder
			if done {
				return usage, imageCount, imageOutputSizes, firstTokenMs, nil
		placeholder
			if err == io.EOF {
				break
		placeholder
			if err != nil {
				if done, processErr := flushData(); processErr != nil {
					return usage, imageCount, imageOutputSizes, firstTokenMs, processErr
			placeholder else if done {
					return usage, imageCount, imageOutputSizes, firstTokenMs, nil
			placeholder
				s.tryWriteOpenAIImagesStreamEvent(c, flusher, &clientDisconnected, &lastDownstreamWriteAt, "error", buildOpenAIImagesStreamErrorBody(err.Error()))
				return usage, imageCount, imageOutputSizes, firstTokenMs, err
		placeholder
	placeholder
		if done, processErr := flushData(); processErr != nil {
			return usage, imageCount, imageOutputSizes, firstTokenMs, processErr
	placeholder else if done {
			return usage, imageCount, imageOutputSizes, firstTokenMs, nil
	placeholder
		if err := finalizePending(); err != nil {
			return usage, imageCount, imageOutputSizes, firstTokenMs, err
	placeholder
		return usage, imageCount, imageOutputSizes, firstTokenMs, nil
placeholder

	type readEvent struct {
		line []byte
		err  error
placeholder
	events := make(chan readEvent, 16)
	done := make(chan struct{placeholder)
	sendEvent := func(ev readEvent) bool {
		select {
		case events <- ev:
			return true
		case <-done:
			return false
	placeholder
placeholder
	var lastReadAt int64
	atomic.StoreInt64(&lastReadAt, time.Now().UnixNano())
	go func() {
		defer close(events)
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if len(line) > 0 {
				atomic.StoreInt64(&lastReadAt, time.Now().UnixNano())
		placeholder
			if len(line) > 0 && !sendEvent(readEvent{line: lineplaceholder) {
				return
		placeholder
			if err == io.EOF {
				return
		placeholder
			if err != nil {
				_ = sendEvent(readEvent{err: errplaceholder)
				return
		placeholder
	placeholder
placeholder()
	defer close(done)

	var intervalTicker *time.Ticker
	if streamInterval > 0 {
		intervalTicker = time.NewTicker(streamInterval)
		defer intervalTicker.Stop()
placeholder
	var intervalCh <-chan time.Time
	if intervalTicker != nil {
		intervalCh = intervalTicker.C
placeholder

	var keepaliveTicker *time.Ticker
	if keepaliveInterval > 0 {
		keepaliveTicker = time.NewTicker(keepaliveInterval)
		defer keepaliveTicker.Stop()
placeholder
	var keepaliveCh <-chan time.Time
	if keepaliveTicker != nil {
		keepaliveCh = keepaliveTicker.C
placeholder

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				if done, processErr := flushData(); processErr != nil {
					return usage, imageCount, imageOutputSizes, firstTokenMs, processErr
			placeholder else if done {
					return usage, imageCount, imageOutputSizes, firstTokenMs, nil
			placeholder
				if err := finalizePending(); err != nil {
					return usage, imageCount, imageOutputSizes, firstTokenMs, err
			placeholder
				return usage, imageCount, imageOutputSizes, firstTokenMs, nil
		placeholder
			if ev.err != nil {
				if done, processErr := flushData(); processErr != nil {
					return usage, imageCount, imageOutputSizes, firstTokenMs, processErr
			placeholder else if done {
					return usage, imageCount, imageOutputSizes, firstTokenMs, nil
			placeholder
				s.tryWriteOpenAIImagesStreamEvent(c, flusher, &clientDisconnected, &lastDownstreamWriteAt, "error", buildOpenAIImagesStreamErrorBody(ev.err.Error()))
				return usage, imageCount, imageOutputSizes, firstTokenMs, ev.err
		placeholder
			done, processErr := processLine(ev.line)
			if processErr != nil {
				return usage, imageCount, imageOutputSizes, firstTokenMs, processErr
		placeholder
			if done {
				return usage, imageCount, imageOutputSizes, firstTokenMs, nil
		placeholder
		case <-intervalCh:
			lastRead := time.Unix(0, atomic.LoadInt64(&lastReadAt))
			if time.Since(lastRead) < streamInterval {
				continue
		placeholder
			if clientDisconnected {
				return usage, imageCount, imageOutputSizes, firstTokenMs, fmt.Errorf("image stream incomplete after timeout")
		placeholder
			logger.LegacyPrintf("service.openai_gateway", "[OpenAI] Images responses stream data interval timeout: interval=%s", streamInterval)
			s.tryWriteOpenAIImagesStreamEvent(c, flusher, &clientDisconnected, &lastDownstreamWriteAt, "error", buildOpenAIImagesStreamErrorBody(fmt.Sprintf("upstream image stream idle for %s", streamInterval)))
			return usage, imageCount, imageOutputSizes, firstTokenMs, fmt.Errorf("image stream data interval timeout")
		case <-keepaliveCh:
			if clientDisconnected || time.Since(lastDownstreamWriteAt) < keepaliveInterval {
				continue
		placeholder
			if _, writeErr := io.WriteString(c.Writer, ":\n\n"); writeErr != nil {
				clientDisconnected = true
				logger.LegacyPrintf("service.openai_gateway", "[OpenAI] Images responses stream client disconnected during keepalive, continue draining upstream for billing")
				continue
		placeholder
			flusher.Flush()
			lastDownstreamWriteAt = time.Now()
	placeholder
placeholder
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
	if requestModel == "" {
		requestModel = "gpt-image-2"
placeholder
	if err := validateOpenAIImagesModel(requestModel); err != nil {
		return nil, err
placeholder
	logger.LegacyPrintf(
		"service.openai_gateway",
		"[OpenAI] Images request routing request_model=%s endpoint=%s account_type=%s uploads=%d",
		requestModel,
		parsed.Endpoint,
		account.Type,
		len(parsed.Uploads),
	)
	upstreamCtx, releaseUpstreamCtx := detachUpstreamContext(ctx)
	defer releaseUpstreamCtx()

	token, _, err := s.GetAccessToken(upstreamCtx, account)
	if err != nil {
		return nil, err
placeholder

	responsesBody, err := buildOpenAIImagesResponsesRequest(parsed, requestModel)
	if err != nil {
		return nil, err
placeholder
	upstreamReq, err := s.buildUpstreamRequest(upstreamCtx, c, account, responsesBody, token, true, parsed.StickySessionSeed(), false)
	if err != nil {
		return nil, err
placeholder
	upstreamReq.Header.Set("Content-Type", "application/json")
	upstreamReq.Header.Set("Accept", "text/event-stream")

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
			s.handleFailoverSideEffects(upstreamCtx, resp, account)
			return nil, &UpstreamFailoverError{
				StatusCode:             resp.StatusCode,
				ResponseBody:           respBody,
				RetryableOnSameAccount: account.IsPoolMode() && account.IsPoolModeRetryableStatus(resp.StatusCode),
		placeholder
	placeholder
		return s.handleErrorResponse(upstreamCtx, resp, c, account, responsesBody)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	var (
		usage            OpenAIUsage
		imageCount       int
		imageOutputSizes []string
		firstTokenMs     *int
	)
	if parsed.Stream {
		usage, imageCount, imageOutputSizes, firstTokenMs, err = s.handleOpenAIImagesOAuthStreamingResponse(resp, c, startTime, parsed.ResponseFormat, openAIImagesStreamPrefix(parsed), requestModel)
		if err != nil {
			if imageCount > 0 {
				return &OpenAIForwardResult{
					RequestID:        resp.Header.Get("x-request-id"),
					Usage:            usage,
					Model:            requestModel,
					UpstreamModel:    requestModel,
					Stream:           parsed.Stream,
					ResponseHeaders:  resp.Header.Clone(),
					Duration:         time.Since(startTime),
					FirstTokenMs:     firstTokenMs,
					ImageCount:       imageCount,
					ImageSize:        parsed.SizeTier,
					ImageInputSize:   parsed.Size,
					ImageOutputSizes: imageOutputSizes,
			placeholder, err
		placeholder
			return nil, err
	placeholder
placeholder else {
		usage, imageCount, imageOutputSizes, err = s.handleOpenAIImagesOAuthNonStreamingResponse(resp, c, parsed.ResponseFormat, requestModel)
		if err != nil {
			return nil, err
	placeholder
placeholder
	if imageCount <= 0 {
		imageCount = parsed.N
placeholder
	return &OpenAIForwardResult{
		RequestID:        resp.Header.Get("x-request-id"),
		Usage:            usage,
		Model:            requestModel,
		UpstreamModel:    requestModel,
		Stream:           parsed.Stream,
		ResponseHeaders:  resp.Header.Clone(),
		Duration:         time.Since(startTime),
		FirstTokenMs:     firstTokenMs,
		ImageCount:       imageCount,
		ImageSize:        parsed.SizeTier,
		ImageInputSize:   parsed.Size,
		ImageOutputSizes: imageOutputSizes,
placeholder, nil
placeholder
