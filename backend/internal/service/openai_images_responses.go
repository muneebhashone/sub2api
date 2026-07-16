package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"errors"
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

// IsOpenAIImagesRetryableUpstreamError reports whether an Images error is an
// upstream server failure that may be retried on another account.
func IsOpenAIImagesRetryableUpstreamError(err *OpenAIImagesUpstreamError) bool {
	return err != nil && err.StatusCode >= http.StatusInternalServerError
placeholder

func openAIImagesSSEErrorStatus(errType, code string) int {
	errType = strings.ToLower(strings.TrimSpace(errType))
	code = strings.ToLower(strings.TrimSpace(code))

	switch {
	case strings.Contains(errType, "rate_limit"), strings.Contains(code, "rate_limit"):
		return http.StatusTooManyRequests
	case strings.Contains(errType, "authentication"), strings.Contains(code, "invalid_api_key"), code == "unauthorized":
		return http.StatusUnauthorized
	case strings.Contains(errType, "permission"), code == "forbidden":
		return http.StatusForbidden
	case strings.Contains(errType, "not_found"), strings.Contains(code, "not_found"):
		return http.StatusNotFound
	case strings.Contains(errType, "invalid_request"),
		errType == "image_generation_user_error",
		code == "moderation_blocked",
		strings.Contains(code, "content_policy"),
		strings.Contains(code, "policy_violation"),
		strings.Contains(code, "safety_violation"):
		return http.StatusBadRequest
	default:
		return http.StatusBadGateway
placeholder
placeholder

func openAIImagesUpstreamErrorResponseBody(err *OpenAIImagesUpstreamError) []byte {
	if err == nil {
		return nil
placeholder
	body := []byte(`{"error":{"type":"","message":""placeholderplaceholder`)
	body, _ = sjson.SetBytes(body, "error.type", err.clientErrorType())
	body, _ = sjson.SetBytes(body, "error.message", err.clientMessage())
	if code := strings.TrimSpace(err.Code); code != "" {
		body, _ = sjson.SetBytes(body, "error.code", code)
placeholder
	if param := strings.TrimSpace(err.Param); param != "" {
		body, _ = sjson.SetBytes(body, "error.param", param)
placeholder
	return body
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
		reconcileOpenAIResponsesImageResultSizes(finalResults, &finalMeta)
		return finalResults, createdAt, usageRaw, finalMeta, true, nil
placeholder

	if len(fallbackResults) > 0 {
		firstMeta := fallbackResults[0]
		mergeOpenAIResponsesImageMeta(&firstMeta, responseMeta)
		reconcileOpenAIResponsesImageResultSizes(fallbackResults, &firstMeta)
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
	case "response.incomplete":
		// 上游在生成预算内未产出图片（超时/被截断），返回 response.incomplete 而非 error。
		// 旧逻辑识别不到，统一报成模糊的 "upstream did not return image output" + 502，
		// 且不触发 failover。这里把它显式建模为可重试的上游错误，使其能换账号重试。
		return openAIImagesIncompleteUpstreamError(gjson.GetBytes(payload, "response"))
	default:
		return nil
placeholder
placeholder

// extractOpenAIImagesModelRefusal 从上游 SSE 响应体提取「模型未出图、改用文字拒绝」
// 的拒绝文本（内容审核场景）。
//
// 上游 response.completed 无图时，模型常以 output_text / message 形式输出拒绝说明
// （如“被安全系统判定为不适合生成”）。这类失败是内容策略拦截，重试/换账号均无效，
// 应把该文本作为内容策略错误透传给客户端。返回空串表示无文字输出（真空响应）。
func extractOpenAIImagesModelRefusal(body []byte) string {
	var b strings.Builder
	collect := func(s string) {
		if s = strings.TrimSpace(s); s != "" {
			if b.Len() > 0 {
				_ = b.WriteByte(' ')
		placeholder
			_, _ = b.WriteString(s)
	placeholder
placeholder
	forEachOpenAISSEDataPayload(string(body), func(payload []byte) {
		if !gjson.ValidBytes(payload) {
			return
	placeholder
		switch gjson.GetBytes(payload, "type").String() {
		case "response.output_text.delta":
			// 流式文本增量。
			collect(gjson.GetBytes(payload, "delta").String())
		case "response.completed", "response.output_item.done":
			// 终态里的 message/output_text。
			gjson.GetBytes(payload, "response.output").ForEach(func(_, item gjson.Result) bool {
				if item.Get("type").String() == "message" {
					item.Get("content").ForEach(func(_, part gjson.Result) bool {
						if part.Get("type").String() == "output_text" {
							collect(part.Get("text").String())
					placeholder
						return true
				placeholder)
			placeholder
				return true
		placeholder)
			if item := gjson.GetBytes(payload, "item"); item.Get("type").String() == "message" {
				item.Get("content").ForEach(func(_, part gjson.Result) bool {
					if part.Get("type").String() == "output_text" {
						collect(part.Get("text").String())
				placeholder
					return true
			placeholder)
		placeholder
	placeholder
placeholder)
	refusal := strings.TrimSpace(b.String())
	// 截断过长文本，避免把整段模型输出塞进错误响应。
	const maxRefusal = 600
	if len(refusal) > maxRefusal {
		refusal = refusal[:maxRefusal]
placeholder
	return refusal
placeholder

// summarizeOpenAIImagesNoOutputBody 从上游 SSE 响应体提取诊断摘要，用于软失败时
// 记录到 ops 日志（上游无图、无标准错误的场景）。提取最终事件类型、response.status、
// incomplete_details.reason，并附 body 截断片段，便于事后定位上游到底返回了什么。
func summarizeOpenAIImagesNoOutputBody(body []byte) string {
	var lastType, status, incompleteReason string
	forEachOpenAISSEDataPayload(string(body), func(payload []byte) {
		if !gjson.ValidBytes(payload) {
			return
	placeholder
		if t := strings.TrimSpace(gjson.GetBytes(payload, "type").String()); t != "" {
			lastType = t
	placeholder
		if resp := gjson.GetBytes(payload, "response"); resp.Exists() {
			if s := strings.TrimSpace(resp.Get("status").String()); s != "" {
				status = s
		placeholder
			if r := strings.TrimSpace(resp.Get("incomplete_details.reason").String()); r != "" {
				incompleteReason = r
		placeholder
	placeholder
placeholder)
	var b strings.Builder
	_, _ = b.WriteString("no_image_output")
	if lastType != "" {
		fmt.Fprintf(&b, " last_event=%s", lastType)
placeholder
	if status != "" {
		fmt.Fprintf(&b, " status=%s", status)
placeholder
	if incompleteReason != "" {
		fmt.Fprintf(&b, " incomplete_reason=%s", incompleteReason)
placeholder
	// 附 body 截断片段（脱敏后），上限 1KB，避免日志膨胀。
	snippet := strings.TrimSpace(string(body))
	const maxSnippet = 1024
	if len(snippet) > maxSnippet {
		snippet = snippet[:maxSnippet] + "...(truncated)"
placeholder
	if snippet != "" {
		fmt.Fprintf(&b, " body=%s", snippet)
placeholder
	return b.String()
placeholder

// openAIImagesIncompleteUpstreamError 从 response.incomplete 事件构建可重试的上游错误。
// incomplete_details.reason 常见取值：max_output_tokens / content_filter 等。
// content_filter 视为客户端错误（400，重试无意义）；其余（生成超时/截断）视为
// 可重试的 502，触发 failover 换账号重试。
func openAIImagesIncompleteUpstreamError(response gjson.Result) *OpenAIImagesUpstreamError {
	if !response.Exists() {
		return nil
placeholder
	reason := strings.TrimSpace(response.Get("incomplete_details.reason").String())
	statusCode := http.StatusBadGateway // 默认可重试（生成未完成）
	errType := "incomplete_error"
	if strings.Contains(strings.ToLower(reason), "content_filter") ||
		strings.Contains(strings.ToLower(reason), "moderation") {
		statusCode = http.StatusBadRequest // 内容过滤，重试无意义
		errType = "image_generation_user_error"
placeholder
	message := "Upstream did not complete image generation"
	if reason != "" {
		message = fmt.Sprintf("Upstream image generation incomplete: %s", reason)
placeholder
	return &OpenAIImagesUpstreamError{
		StatusCode:        statusCode,
		ErrorType:         errType,
		Code:              "response_incomplete",
		Message:           sanitizeUpstreamErrorMessage(message),
		UpstreamRequestID: strings.TrimSpace(response.Get("id").String()),
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
	statusCode := openAIImagesSSEErrorStatus(errType, code)
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

// openAIImagesErrorTypeForStatus returns an OpenAI-style error type when the
// upstream body does not provide one of its own.
func openAIImagesErrorTypeForStatus(status int) string {
	switch {
	case status == http.StatusBadRequest:
		return "invalid_request_error"
	case status == http.StatusUnauthorized:
		return "authentication_error"
	case status == http.StatusForbidden:
		return "permission_error"
	case status == http.StatusNotFound:
		return "not_found_error"
	case status == http.StatusTooManyRequests:
		return "rate_limit_error"
	case status >= 500:
		return "api_error"
	default:
		return "upstream_error"
placeholder
placeholder

// openAIImagesUpstreamErrorFromHTTP builds an OpenAIImagesUpstreamError from a
// non-2xx upstream HTTP response, preserving the real status code, type, code,
// message and param so the client sees the actual upstream error instead of a
// generic 502.
func openAIImagesUpstreamErrorFromHTTP(statusCode int, header http.Header, body []byte) *OpenAIImagesUpstreamError {
	errType := strings.TrimSpace(gjson.GetBytes(body, "error.type").String())
	code := strings.TrimSpace(extractUpstreamErrorCode(body))
	param := strings.TrimSpace(gjson.GetBytes(body, "error.param").String())
	message := sanitizeUpstreamErrorMessage(strings.TrimSpace(extractUpstreamErrorMessage(body)))
	if message == "" {
		message = fmt.Sprintf("Upstream request failed (status %d)", statusCode)
placeholder
	if errType == "" {
		errType = openAIImagesErrorTypeForStatus(statusCode)
placeholder
	requestID := ""
	if header != nil {
		requestID = strings.TrimSpace(header.Get("x-request-id"))
placeholder
	return &OpenAIImagesUpstreamError{
		StatusCode:        statusCode,
		ErrorType:         errType,
		Code:              code,
		Message:           message,
		Param:             param,
		UpstreamRequestID: requestID,
placeholder
placeholder

// handleOpenAIImagesErrorResponse is the non-failover error handler for the
// images endpoints (/v1/images/generations and /v1/images/edits). Unlike the
// generic handleErrorResponse — which collapses every non-failover upstream
// error into a generic 502 "Upstream request failed" — it surfaces the real
// upstream status code and error message/type/code/param to the client. This
// mirrors how the Chat Completions and Messages compat paths use
// handleCompatErrorResponse.
//
// It returns an *OpenAIImagesUpstreamError (already written to the client) so
// the images handler treats it as a terminal user-facing error rather than
// re-writing a fallback response.
func (s *OpenAIGatewayService) handleOpenAIImagesErrorResponse(
	ctx context.Context,
	resp *http.Response,
	c *gin.Context,
	account *Account,
	requestedModel ...string,
) (*OpenAIForwardResult, error) {
	body := s.readUpstreamErrorBody(resp)

	upstreamMsg := sanitizeUpstreamErrorMessage(strings.TrimSpace(extractUpstreamErrorMessage(body)))
	upstreamDetail := ""
	if s.cfg != nil && s.cfg.Gateway.LogUpstreamErrorBody {
		maxBytes := s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes
		if maxBytes <= 0 {
			maxBytes = 2048
	placeholder
		upstreamDetail = truncateString(string(body), maxBytes)
placeholder
	setOpsUpstreamError(c, resp.StatusCode, upstreamMsg, upstreamDetail)

	if s.cfg != nil && s.cfg.Gateway.LogUpstreamErrorBody {
		logger.LegacyPrintf("service.openai_gateway",
			"OpenAI images upstream error %d (account=%d platform=%s type=%s): %s",
			resp.StatusCode,
			account.ID,
			account.Platform,
			account.Type,
			truncateForLog(body, s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes),
		)
placeholder

	// Honor admin-configured error passthrough rules first.
	if status, errType, errMsg, matched := applyErrorPassthroughRule(
		c,
		account.Platform,
		resp.StatusCode,
		body,
		http.StatusBadGateway,
		"upstream_error",
		"Upstream request failed",
	); matched {
		upErr := &OpenAIImagesUpstreamError{
			StatusCode:        status,
			ErrorType:         errType,
			Message:           errMsg,
			UpstreamRequestID: strings.TrimSpace(resp.Header.Get("x-request-id")),
	placeholder
		writeOpenAIImagesUpstreamErrorResponse(c, upErr)
		return nil, upErr
placeholder

	// If the account is not configured to handle this status code, fall back to
	// a generic gateway error without exposing upstream internals (mirrors
	// handleCompatErrorResponse).
	if !account.ShouldHandleErrorCode(resp.StatusCode) {
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
		upErr := &OpenAIImagesUpstreamError{
			StatusCode:        http.StatusInternalServerError,
			ErrorType:         "upstream_error",
			Message:           "Upstream gateway error",
			UpstreamRequestID: strings.TrimSpace(resp.Header.Get("x-request-id")),
	placeholder
		writeOpenAIImagesUpstreamErrorResponse(c, upErr)
		return nil, upErr
placeholder

	// Track rate limits / decide whether to disable the account (secondary failover).
	var modelForCooldown string
	if len(requestedModel) > 0 {
		modelForCooldown = strings.TrimSpace(requestedModel[0])
placeholder
	shouldDisable := s.handleOpenAIAccountUpstreamError(ctx, account, resp.StatusCode, resp.Header, body, modelForCooldown)
	kind := "http_error"
	if shouldDisable {
		kind = "failover"
placeholder
	appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
		Platform:           account.Platform,
		AccountID:          account.ID,
		AccountName:        account.Name,
		UpstreamStatusCode: resp.StatusCode,
		UpstreamRequestID:  resp.Header.Get("x-request-id"),
		Kind:               kind,
		Message:            upstreamMsg,
		Detail:             upstreamDetail,
placeholder)
	if shouldDisable {
		return nil, &UpstreamFailoverError{
			StatusCode:             resp.StatusCode,
			ResponseBody:           body,
			RetryableOnSameAccount: account.IsPoolMode() && account.IsPoolModeRetryableStatus(resp.StatusCode),
	placeholder
placeholder

	// Surface the real upstream error to the client.
	upErr := openAIImagesUpstreamErrorFromHTTP(resp.StatusCode, resp.Header, body)
	writeOpenAIImagesUpstreamErrorResponse(c, upErr)
	return nil, upErr
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
	if c == nil || c.Writer == nil || err == nil {
		return false
placeholder
	if c.Writer.Written() && OpenAIImagesJSONKeepaliveAdjustedWrittenSize(c) >= 0 {
		return false
placeholder
	StopOpenAIImagesJSONKeepaliveCommitted(c)
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

func (s *OpenAIGatewayService) parseOpenAIImagesSSEUsageBytes(data []byte, usage *OpenAIUsage) {
	s.parseSSEUsageBytes(data, usage)
	if usage == nil || !gjson.ValidBytes(data) || gjson.GetBytes(data, "type").String() != "response.completed" {
		return
placeholder
	if toolUsage, ok := openAIImagesToolUsageFromGJSON(gjson.GetBytes(data, "response.tool_usage.image_gen")); ok {
		*usage = toolUsage
placeholder
placeholder

func openAIImagesToolUsageFromGJSON(value gjson.Result) (OpenAIUsage, bool) {
	if !value.Exists() || !value.IsObject() {
		return OpenAIUsage{placeholder, false
placeholder
	inputTokens, inputOK := boundedJSONNonNegativeInt(value.Get("input_tokens"))
	outputTokens, outputOK := boundedJSONNonNegativeInt(value.Get("output_tokens"))
	imageOutputTokens, imageOutputOK := boundedJSONNonNegativeInt(value.Get("output_tokens_details.image_tokens"))
	if !inputOK || !outputOK || !imageOutputOK {
		return OpenAIUsage{placeholder, false
placeholder
	return OpenAIUsage{
		InputTokens:       inputTokens,
		OutputTokens:      outputTokens,
		ImageOutputTokens: imageOutputTokens,
placeholder, true
placeholder

// boundedJSONNonNegativeInt parses integral JSON exponent notation without
// invoking an arbitrary-precision parser on an upstream-controlled exponent.
func boundedJSONNonNegativeInt(value gjson.Result) (int, bool) {
	if !value.Exists() || value.Type != gjson.Number {
		return 0, false
placeholder
	raw := value.Raw
	if len(raw) == 0 || len(raw) > 64 || raw[0] == '-' {
		return 0, false
placeholder

	mantissaEnd := len(raw)
	for i, c := range raw {
		if c != 'e' && c != 'E' {
			continue
	placeholder
		mantissaEnd = i
		break
placeholder

	digits := raw[:mantissaEnd]
	fractionDigits := 0
	digitCount := 0
	dotSeen := false
	mantissaIsZero := true
	for _, c := range digits {
		switch {
		case c == '.' && !dotSeen:
			dotSeen = true
		case c >= '0' && c <= '9':
			digitCount++
			mantissaIsZero = mantissaIsZero && c == '0'
			if dotSeen {
				fractionDigits++
		placeholder
		default:
			return 0, false
	placeholder
placeholder

	exponent := 0
	if mantissaEnd < len(raw) {
		exponentRaw := raw[mantissaEnd+1:]
		negative := false
		if len(exponentRaw) > 0 && (exponentRaw[0] == '+' || exponentRaw[0] == '-') {
			negative = exponentRaw[0] == '-'
			exponentRaw = exponentRaw[1:]
	placeholder
		if len(exponentRaw) == 0 {
			return 0, false
	placeholder
		for len(exponentRaw) > 1 && exponentRaw[0] == '0' {
			exponentRaw = exponentRaw[1:]
	placeholder
		for _, digit := range exponentRaw {
			if digit < '0' || digit > '9' {
				return 0, false
		placeholder
	placeholder
		if mantissaIsZero {
			return 0, true
	placeholder
		if len(exponentRaw) > 3 {
			return 0, false
	placeholder
		for _, digit := range exponentRaw {
			exponent = exponent*10 + int(digit-'0')
	placeholder
		if exponent > 100 {
			return 0, false
	placeholder
		if negative {
			exponent = -exponent
	placeholder
placeholder

	trailingZeros := exponent - fractionDigits
	scaleReduction := 0
	if trailingZeros < 0 {
		scaleReduction = -trailingZeros
		remaining := scaleReduction
		allZeros := true
		for i := len(digits) - 1; i >= 0; i-- {
			if digits[i] == '.' {
				continue
		placeholder
			if digits[i] != '0' {
				allZeros = false
				if remaining > 0 {
					return 0, false
			placeholder
		placeholder
			if remaining > 0 {
				remaining--
		placeholder
	placeholder
		if remaining > 0 {
			if allZeros {
				return 0, true
		placeholder
			return 0, false
	placeholder
placeholder

	maxInt := int(^uint(0) >> 1)
	parsed := 0
	digitsToAccumulate := digitCount - scaleReduction
	for _, c := range digits {
		if c == '.' {
			continue
	placeholder
		if digitsToAccumulate <= 0 {
			break
	placeholder
		if parsed > (maxInt-int(c-'0'))/10 {
			return 0, false
	placeholder
		parsed = parsed*10 + int(c-'0')
		digitsToAccumulate--
placeholder
	if trailingZeros < 0 {
		return parsed, true
placeholder
	for ; trailingZeros > 0; trailingZeros-- {
		if parsed > maxInt/10 {
			return 0, false
	placeholder
		parsed *= 10
placeholder
	return parsed, true
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
		s.parseOpenAIImagesSSEUsageBytes(data, &usage)
placeholder)
	results, createdAt, usageRaw, firstMeta, _, err := collectOpenAIImagesFromResponsesBody(body)
	if err != nil {
		return OpenAIUsage{placeholder, 0, nil, err
placeholder
	if len(results) == 0 {
		if upstreamErr := extractOpenAIImagesUpstreamError(body); upstreamErr != nil {
			setOpsUpstreamError(c, upstreamErr.clientStatusCode(), upstreamErr.clientMessage(), "")
			if !IsOpenAIImagesRetryableUpstreamError(upstreamErr) {
				writeOpenAIImagesUpstreamErrorResponse(c, upstreamErr)
		placeholder
			return OpenAIUsage{placeholder, 0, nil, upstreamErr
	placeholder
		// 软失败兜底：上游无图。先区分两种情形（实测真因，见下）：
		//
		// (A) 内容审核拒绝：模型未出图，但输出了文字拒绝（response.completed 里带
		//     output_text / message，内容如“被安全系统判定为不适合生成”）。这是用户
		//     prompt 触发 OpenAI 内容策略，模型主动拒绝改用文字回应。**换账号/重试均无效**
		//     （内容层拦截，与账号/承载模型无关），应把拒绝理由作为 400 透传给客户端，
		//     避免无谓地重试 + 消耗其它账号配额，且让客户端拿到可读的拒绝原因。
		// (B) 真空响应：既无图也无任何文字输出（罕见，如偶发路由到 gpt-5.x-mini、
		//     image_gen 工具未执行）。这是上游的概率性失败，此时才按可重试处理。
		if refusal := extractOpenAIImagesModelRefusal(body); refusal != "" {
			refusalErr := &OpenAIImagesUpstreamError{
				StatusCode: http.StatusBadRequest,
				ErrorType:  "image_generation_user_error",
				Code:       "content_policy_violation",
				Message:    sanitizeUpstreamErrorMessage(refusal),
		placeholder
			setOpsUpstreamError(c, http.StatusBadRequest, refusalErr.clientMessage(), summarizeOpenAIImagesNoOutputBody(body))
			writeOpenAIImagesUpstreamErrorResponse(c, refusalErr)
			return OpenAIUsage{placeholder, 0, nil, refusalErr
	placeholder
		// (B) 真空响应：记录上游诊断摘要到 ops（last_event/status/model/body 片段）便于
		// 排查，并返回 UpstreamFailoverError 触发重试。因实测为「同账号概率性失败」，优先
		// RetryableOnSameAccount 同账号快速重试（默认 3 次，大概率某次正常出图），用尽后
		// 由 handler 自然换账号 failover（switchCount 上限保护），既提高成功率又不无谓
		// 消耗其它账号配额。
		setOpsUpstreamError(c, http.StatusBadGateway, "upstream did not return image output", summarizeOpenAIImagesNoOutputBody(body))
		return OpenAIUsage{placeholder, 0, nil, &UpstreamFailoverError{
			StatusCode:             http.StatusBadGateway,
			ResponseBody:           body,
			RetryableOnSameAccount: true,
	placeholder
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
	writerSizeBeforeResponse := OpenAIImagesJSONKeepaliveAdjustedWrittenSize(c)

	processData := func(dataBytes []byte) {
		if processDataDone || processDataErr != nil {
			return
	placeholder
		if firstTokenMs == nil {
			ms := int(time.Since(startTime).Milliseconds())
			firstTokenMs = &ms
	placeholder
		s.parseOpenAIImagesSSEUsageBytes(dataBytes, &usage)
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
			reconcileOpenAIResponsesImageResultSizes(finalResults, nil)
			if len(finalResults) == 0 {
				outputErr := fmt.Errorf("upstream did not return image output")
				// 软失败：response.completed 事件里没有图片。记录上游诊断摘要到 ops，
				// 与非流式路径保持一致，避免上游响应信息丢失。
				setOpsUpstreamError(c, http.StatusBadGateway, "upstream did not return image output", summarizeOpenAIImagesNoOutputBody(dataBytes))
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
				retryable := IsOpenAIImagesRetryableUpstreamError(upstreamErr)
				if !clientDisconnected && (!retryable || c.Writer.Size() != writerSizeBeforeResponse) {
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
			finalResults := append([]openAIResponsesImageResult(nil), pendingResults...)
			for i := range finalResults {
				mergeOpenAIResponsesImageMeta(&finalResults[i], streamMeta)
		placeholder
			reconcileOpenAIResponsesImageResultSizes(finalResults, nil)
			for _, img := range finalResults {
				key := openAIResponsesImageResultKey("", img)
				if _, exists := emitted[key]; exists {
					continue
			placeholder
				payload := buildOpenAIImagesStreamCompletedPayload(eventName, img, format, createdAt, nil)
				emitted[key] = struct{placeholder{placeholder
				s.tryWriteOpenAIImagesStreamEvent(c, flusher, &clientDisconnected, &lastDownstreamWriteAt, eventName, payload)
		placeholder
			imageCount = len(emitted)
			imageOutputSizes = openAIResponsesImageResultSizes(finalResults)
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
		respBody := s.readUpstreamErrorBody(resp)
		_ = resp.Body.Close()
		respBody = s.redactAgentIdentitySensitiveBody(upstreamCtx, account, respBody)
		if !agentIdentityTaskRecoveryWasTried(ctx) && s.isAgentIdentityAccount(ctx, account) && isAgentIdentityTaskInvalidHTTPResponse(resp.StatusCode, respBody) {
			expectedTaskID := account.GetCredential("task_id")
			if err := s.recoverAgentIdentityTask(ctx, account, expectedTaskID); err != nil {
				return nil, fmt.Errorf("agent identity task recovery failed: %w", err)
		placeholder
			return s.forwardOpenAIImagesOAuth(markAgentIdentityTaskRecoveryTried(ctx), c, account, parsed, channelMappedModel)
	placeholder
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
			s.handleFailoverSideEffects(upstreamCtx, resp, account, respBody, requestModel)
			return nil, &UpstreamFailoverError{
				StatusCode:             resp.StatusCode,
				ResponseBody:           respBody,
				RetryableOnSameAccount: account.IsPoolMode() && account.IsPoolModeRetryableStatus(resp.StatusCode),
		placeholder
	placeholder
		return s.handleOpenAIImagesErrorResponse(upstreamCtx, resp, c, account, requestModel)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	var (
		usage            OpenAIUsage
		imageCount       int
		imageOutputSizes []string
		firstTokenMs     *int
	)
	// 与 handleOpenAIImagesOAuthResponseError 的比较端同口径：排除非流式 JSON
	// keepalive 心跳字节，避免 failover 第 2 轮起把上一轮心跳残留误判为已写响应。
	writerSizeBeforeResponse := OpenAIImagesJSONKeepaliveAdjustedWrittenSize(c)
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
			return nil, s.handleOpenAIImagesOAuthResponseError(
				upstreamCtx,
				c,
				account,
				requestModel,
				safeUpstreamURL(upstreamReq.URL.String()),
				resp,
				writerSizeBeforeResponse,
				err,
			)
	placeholder
placeholder else {
		usage, imageCount, imageOutputSizes, err = s.handleOpenAIImagesOAuthNonStreamingResponse(resp, c, parsed.ResponseFormat, requestModel)
		if err != nil {
			return nil, s.handleOpenAIImagesOAuthResponseError(
				upstreamCtx,
				c,
				account,
				requestModel,
				safeUpstreamURL(upstreamReq.URL.String()),
				resp,
				writerSizeBeforeResponse,
				err,
			)
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

func (s *OpenAIGatewayService) handleOpenAIImagesOAuthResponseError(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	requestedModel string,
	upstreamURL string,
	resp *http.Response,
	writerSizeBeforeResponse int,
	err error,
) error {
	var upstreamErr *OpenAIImagesUpstreamError
	if !errors.As(err, &upstreamErr) {
		return err
placeholder

	retryable := IsOpenAIImagesRetryableUpstreamError(upstreamErr)
	responseWritten := c != nil && c.Writer != nil && OpenAIImagesJSONKeepaliveAdjustedWrittenSize(c) != writerSizeBeforeResponse
	kind := "http_error"
	if retryable {
		kind = "failover"
		if responseWritten {
			kind = "retry_exhausted_failover"
	placeholder
placeholder

	requestID := strings.TrimSpace(upstreamErr.UpstreamRequestID)
	headers := http.Header(nil)
	if resp != nil {
		headers = resp.Header.Clone()
		if requestID == "" {
			requestID = strings.TrimSpace(resp.Header.Get("x-request-id"))
	placeholder
placeholder
	appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
		Platform:           account.Platform,
		AccountID:          account.ID,
		AccountName:        account.Name,
		UpstreamStatusCode: upstreamErr.StatusCode,
		UpstreamRequestID:  requestID,
		UpstreamURL:        upstreamURL,
		Kind:               kind,
		Message:            upstreamErr.clientMessage(),
placeholder)

	if !retryable || responseWritten {
		return err
placeholder

	responseBody := openAIImagesUpstreamErrorResponseBody(upstreamErr)
	s.handleOpenAIAccountUpstreamError(ctx, account, upstreamErr.StatusCode, headers, responseBody, requestedModel)
	return &UpstreamFailoverError{
		StatusCode:             upstreamErr.StatusCode,
		ResponseBody:           responseBody,
		ResponseHeaders:        headers,
		RetryableOnSameAccount: account.IsPoolMode() && account.IsPoolModeRetryableStatus(upstreamErr.StatusCode),
placeholder
placeholder
