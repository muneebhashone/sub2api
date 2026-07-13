//go:build unit

package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestPatchGrokResponsesBodySetsMappedModelAndDropsUnsupportedFields(t *testing.T) {
	t.Parallel()

	body := []byte(`{
		"model": "grok",
		"input": "hello",
		"prompt_cache_retention": "24h",
		"safety_identifier": "user-1",
		"reasoning": {"effort": "high"placeholder
placeholder`)

	patched, err := patchGrokResponsesBody(body, "grok-4.3")
placeholder
	require.True(t, json.Valid(patched))
	require.Equal(t, "grok-4.3", gjson.GetBytes(patched, "model").String())
	require.False(t, gjson.GetBytes(patched, "prompt_cache_retention").Exists())
	require.False(t, gjson.GetBytes(patched, "safety_identifier").Exists())
	require.Equal(t, "high", gjson.GetBytes(patched, "reasoning.effort").String())
placeholder

func TestExtractGrokResponsesReasoningEffortSupportsOpenAICompatibleField(t *testing.T) {
	t.Parallel()

	effort := extractOpenAIReasoningEffortFromBody(
		[]byte(`{"model":"grok-4.3","reasoning_effort":"high"placeholder`),
		"grok-4.3",
	)
	require.NotNil(t, effort)
	require.Equal(t, "high", *effort)
placeholder

func TestPatchGrokResponsesBodyDropsGrok45ReasoningUnsupportedFields(t *testing.T) {
	t.Parallel()

	body := []byte(`{
		"model": "grok-latest",
		"input": "hello",
		"presence_penalty": 0.1,
		"presencePenalty": 0.2,
		"frequency_penalty": 0.3,
		"frequencyPenalty": 0.4,
		"stop": ["done"]
placeholder`)

	patched, err := patchGrokResponsesBody(body, "grok-4.5")
placeholder
	require.True(t, json.Valid(patched))
	require.Equal(t, "grok-4.5", gjson.GetBytes(patched, "model").String())
	require.False(t, gjson.GetBytes(patched, "presence_penalty").Exists())
	require.False(t, gjson.GetBytes(patched, "presencePenalty").Exists())
	require.False(t, gjson.GetBytes(patched, "frequency_penalty").Exists())
	require.False(t, gjson.GetBytes(patched, "frequencyPenalty").Exists())
	require.False(t, gjson.GetBytes(patched, "stop").Exists())
placeholder

func TestPatchGrokResponsesBodyKeepsPenaltyAndStopFieldsForNon45Models(t *testing.T) {
	t.Parallel()

	body := []byte(`{
		"model": "grok-4.3",
		"input": "hello",
		"presence_penalty": 0.1,
		"frequency_penalty": 0.2,
		"stop": ["done"]
placeholder`)

	patched, err := patchGrokResponsesBody(body, "grok-4.3")
placeholder
	require.True(t, json.Valid(patched))
	require.Equal(t, "grok-4.3", gjson.GetBytes(patched, "model").String())
	require.Equal(t, 0.1, gjson.GetBytes(patched, "presence_penalty").Float())
	require.Equal(t, 0.2, gjson.GetBytes(patched, "frequency_penalty").Float())
	require.Len(t, gjson.GetBytes(patched, "stop").Array(), 1)
placeholder

func TestPatchGrokResponsesBodyDropsNestedUnsupportedFields(t *testing.T) {
	t.Parallel()

	body := []byte(`{
		"model": "grok",
		"input": "hello",
		"external_web_access": true,
		"tools": [
			{"type": "function", "name": "kept_fn", "external_web_access": true, "parameters": {"type": "object", "properties": {"q": {"type": "string", "external_web_access": trueplaceholderplaceholderplaceholderplaceholder
		],
		"metadata": {"external_web_access": falseplaceholder
placeholder`)

	patched, err := patchGrokResponsesBody(body, "grok-4.3")
placeholder
	require.True(t, json.Valid(patched))
	require.False(t, strings.Contains(string(patched), "external_web_access"))
	require.Equal(t, "kept_fn", gjson.GetBytes(patched, "tools.0.name").String())
placeholder

func TestPatchGrokResponsesBodyDropsUnsupportedNamespaceTools(t *testing.T) {
	t.Parallel()

	body := []byte(`{
		"model": "grok",
		"input": "hello",
		"tools": [
			{"type": "namespace", "namespace": "functions", "tools": [{"type": "function", "name": "inner"placeholder]placeholder,
			{"type": "function", "name": "kept_fn", "parameters": {"type": "object"placeholderplaceholder,
			{"type": "shell", "name": "kept_shell"placeholder
		],
		"tool_choice": {"type": "function", "name": "kept_fn"placeholder
placeholder`)

	patched, err := patchGrokResponsesBody(body, "grok-4.3")
placeholder
	require.True(t, json.Valid(patched))
	require.Equal(t, "grok-4.3", gjson.GetBytes(patched, "model").String())
	require.Len(t, gjson.GetBytes(patched, "tools").Array(), 2)
	require.False(t, gjson.GetBytes(patched, `tools.#(type=="namespace")`).Exists())
	require.True(t, gjson.GetBytes(patched, `tools.#(type=="function")`).Exists())
	require.True(t, gjson.GetBytes(patched, `tools.#(type=="shell")`).Exists())
	require.Equal(t, "kept_fn", gjson.GetBytes(patched, "tool_choice.name").String())
placeholder

func TestPatchGrokResponsesBodyDropsToolChoiceWhenNoSupportedToolsRemain(t *testing.T) {
	t.Parallel()

	body := []byte(`{
		"model": "grok",
		"input": "hello",
		"tools": [
			{"type": "namespace", "namespace": "functions"placeholder,
			{"type": "image_generation", "model": "gpt-image-2"placeholder
		],
		"tool_choice": {"type": "namespace", "namespace": "functions"placeholder
placeholder`)

	patched, err := patchGrokResponsesBody(body, "grok-4.3")
placeholder
	require.True(t, json.Valid(patched))
	require.False(t, gjson.GetBytes(patched, "tools").Exists())
	require.False(t, gjson.GetBytes(patched, "tool_choice").Exists())
placeholder

func TestBuildGrokResponsesRequestUsesAccountBaseURLAndBearerToken(t *testing.T) {
	t.Setenv(xai.EnvAllowUnsafeURLOverrides, "true")

	account := &Account{
		Platform: PlatformGrok,
		Type:     AccountTypeOAuth,
placeholder
			"base_url": "https://xai.test/v1/",
	placeholder,
placeholder

	req, err := buildGrokResponsesRequest(context.Background(), nil, account, []byte(`{"model":"grok-4.3"placeholder`), "access-token", "isolated-cache-id")
placeholder
	require.Equal(t, http.MethodPost, req.Method)
	require.Equal(t, "https://xai.test/v1/responses", req.URL.String())
	require.Equal(t, "Bearer access-token", req.Header.Get("Authorization"))
	require.Equal(t, "application/json", req.Header.Get("Content-Type"))
	require.Contains(t, req.Header.Get("Accept"), "text/event-stream")
	require.Equal(t, grokCLIVersion, req.Header.Get("X-Grok-Client-Version"))
	require.Equal(t, "isolated-cache-id", req.Header.Get(grokConversationIDHeader))

	data, err := io.ReadAll(req.Body)
placeholder
	require.Equal(t, `{"model":"grok-4.3"placeholder`, strings.TrimSpace(string(data)))
placeholder

func TestBuildGrokResponsesRequestRejectsUnsafeAccountBaseURL(t *testing.T) {
	t.Parallel()

	account := &Account{
		Platform: PlatformGrok,
		Type:     AccountTypeOAuth,
placeholder
			"base_url": "https://xai.test/v1",
	placeholder,
placeholder

	_, err := buildGrokResponsesRequest(context.Background(), nil, account, []byte(`{"model":"grok-4.3"placeholder`), "access-token", "")
placeholder
	require.Contains(t, err.Error(), "invalid base url")
placeholder

func TestGrokMediaGenerationGateCoversImagesAndVideo(t *testing.T) {
	tests := []struct {
		name     string
		endpoint GrokMediaEndpoint
		want     bool
placeholder{
		{name: "image generation", endpoint: GrokMediaEndpointImagesGenerations, want: trueplaceholder,
		{name: "image edit", endpoint: GrokMediaEndpointImagesEdits, want: trueplaceholder,
		{name: "video generation", endpoint: GrokMediaEndpointVideosGenerations, want: trueplaceholder,
		{name: "video status", endpoint: GrokMediaEndpointVideoStatus, want: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.endpoint.IsGenerationRequest())
	placeholder)
placeholder
placeholder

func TestExtractGrokMediaModelSupportsJSONAndMultipart(t *testing.T) {
	require.Equal(t, "grok-imagine", ExtractGrokMediaModel("application/json", []byte(`{"model":"grok-imagine"placeholder`)))

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	require.NoError(t, writer.WriteField("prompt", "draw a cat"))
	require.NoError(t, writer.WriteField("model", "grok-imagine-edit"))
	require.NoError(t, writer.Close())

	require.Equal(t, "grok-imagine-edit", ExtractGrokMediaModel(writer.FormDataContentType(), buf.Bytes()))
placeholder

func TestParseGrokMediaRequestBuildsMultipartModerationBody(t *testing.T) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	require.NoError(t, writer.WriteField("prompt", "edit this private image"))
	require.NoError(t, writer.WriteField("model", "grok-imagine-edit"))
	partHeader := textproto.MIMEHeader{placeholder
	partHeader.Set("Content-Disposition", `form-data; name="image"; filename="input.png"`)
	partHeader.Set("Content-Type", "image/png")
	part, err := writer.CreatePart(partHeader)
placeholder
	_, err = part.Write([]byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0aplaceholder)
placeholder
	require.NoError(t, writer.Close())

	info := ParseGrokMediaRequest(writer.FormDataContentType(), buf.Bytes())
	require.Equal(t, "grok-imagine-edit", info.Model)
	require.Equal(t, "edit this private image", info.Prompt)

	moderationBody := info.ModerationBody()
	require.NotEmpty(t, moderationBody)
	require.Equal(t, "edit this private image", gjson.GetBytes(moderationBody, "prompt").String())
	require.True(t, strings.HasPrefix(gjson.GetBytes(moderationBody, "images.0.image_url").String(), "data:image/"))
placeholder

func TestParseGrokMediaVideoRequestResolution(t *testing.T) {
	info := ParseGrokMediaRequest("application/json", []byte(`{"model":"grok-imagine-video","prompt":"waves","resolution":"720p"placeholder`))

	require.Equal(t, "grok-imagine-video", info.Model)
	require.Equal(t, "720p", info.Resolution)
placeholder

func TestNormalizeGrokMediaModelForEndpoint(t *testing.T) {
	tests := []struct {
		name          string
		endpoint      GrokMediaEndpoint
		model         string
		hasInputImage bool
		want          string
placeholder{
		{name: "image generation alias", endpoint: GrokMediaEndpointImagesGenerations, model: "grok-imagine", want: "grok-imagine-image-quality"placeholder,
		{name: "image edit alias", endpoint: GrokMediaEndpointImagesEdits, model: "grok-imagine", want: "grok-imagine-image-quality"placeholder,
		{name: "image quality passthrough", endpoint: GrokMediaEndpointImagesGenerations, model: "grok-imagine-image-quality", want: "grok-imagine-image-quality"placeholder,
		{name: "image fast passthrough", endpoint: GrokMediaEndpointImagesGenerations, model: "grok-imagine-image", want: "grok-imagine-image"placeholder,
		{name: "video passthrough", endpoint: GrokMediaEndpointVideosGenerations, model: "grok-imagine-video", want: "grok-imagine-video"placeholder,
		{name: "video 1.5 text-only fallback", endpoint: GrokMediaEndpointVideosGenerations, model: "grok-imagine-video-1.5", want: "grok-imagine-video"placeholder,
		{name: "video 1.5 image-to-video passthrough", endpoint: GrokMediaEndpointVideosGenerations, model: "grok-imagine-video-1.5", hasInputImage: true, want: "grok-imagine-video-1.5"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, normalizeGrokMediaModelForEndpoint(tt.endpoint, tt.model, tt.hasInputImage))
	placeholder)
placeholder
placeholder

func TestForwardGrokMediaImagesGenerationNormalizesImagineAlias(t *testing.T) {
	t.Setenv(xai.EnvAllowUnsafeURLOverrides, "true")
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok-imagine","prompt":"draw a cat"placeholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	account := &Account{
		ID:          61,
		Name:        "grok",
		Platform:    PlatformGrok,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "api-key",
			"base_url": "https://xai.test/v1",
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type":   []string{"application/json"placeholder,
			"Xai-Request-Id": []string{"xai-image-req"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(`{"data":[]placeholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder

	result, err := svc.ForwardGrokMedia(context.Background(), c, account, GrokMediaEndpointImagesGenerations, "", body, "application/json")
placeholder
	require.Equal(t, "https://xai.test/v1/images/generations", upstream.lastReq.URL.String())
	require.Equal(t, http.MethodPost, upstream.lastReq.Method)
	require.Equal(t, "Bearer api-key", upstream.lastReq.Header.Get("Authorization"))
	require.Equal(t, "application/json", upstream.lastReq.Header.Get("Content-Type"))
	require.Equal(t, grokCLIVersion, upstream.lastReq.Header.Get("X-Grok-Client-Version"))
	require.JSONEq(t, `{"model":"grok-imagine-image-quality","prompt":"draw a cat"placeholder`, string(upstream.lastBody))
	require.Equal(t, http.StatusOK, recorder.Code)
	require.JSONEq(t, `{"data":[]placeholder`, recorder.Body.String())
	require.Equal(t, "xai-image-req", result.RequestID)
	require.Equal(t, "grok-imagine-image-quality", result.Model)
	require.Equal(t, "grok-imagine-image-quality", result.BillingModel)
	require.Equal(t, 1, result.ImageCount)
	require.Equal(t, ImageBillingSize2K, result.ImageSize)
placeholder

func TestForwardGrokMediaImagesGenerationStripsUnsupportedSize(t *testing.T) {
	t.Setenv(xai.EnvAllowUnsafeURLOverrides, "true")
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok-imagine-image","prompt":"draw a cat","size":"1024x1024"placeholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	account := &Account{
		ID:          65,
		Name:        "grok",
		Platform:    PlatformGrok,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "api-key",
			"base_url": "https://xai.test/v1",
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"application/json"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(`{"data":[]placeholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder

	result, err := svc.ForwardGrokMedia(context.Background(), c, account, GrokMediaEndpointImagesGenerations, "", body, "application/json")
placeholder
	require.JSONEq(t, `{"model":"grok-imagine-image","prompt":"draw a cat"placeholder`, string(upstream.lastBody))
	require.Equal(t, ImageBillingSize1K, result.ImageSize)
	require.Equal(t, "1024x1024", result.ImageInputSize)
placeholder

func TestForwardGrokMediaImagesEditMultipartConvertsToJSON(t *testing.T) {
	t.Setenv(xai.EnvAllowUnsafeURLOverrides, "true")
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	require.NoError(t, writer.WriteField("model", "grok-imagine-edit"))
	require.NoError(t, writer.WriteField("prompt", "edit this private image"))
	partHeader := textproto.MIMEHeader{placeholder
	partHeader.Set("Content-Disposition", `form-data; name="image"; filename="input.png"`)
	partHeader.Set("Content-Type", "image/png")
	part, err := writer.CreatePart(partHeader)
placeholder
	_, err = part.Write([]byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0aplaceholder)
placeholder
	require.NoError(t, writer.Close())

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/edits", bytes.NewReader(buf.Bytes()))
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	account := &Account{
		ID:          62,
		Name:        "grok",
		Platform:    PlatformGrok,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "api-key",
			"base_url": "https://xai.test/v1",
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"application/json"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(`{"data":[]placeholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder

	_, err = svc.ForwardGrokMedia(context.Background(), c, account, GrokMediaEndpointImagesEdits, "", buf.Bytes(), writer.FormDataContentType())
placeholder
	require.Equal(t, "https://xai.test/v1/images/edits", upstream.lastReq.URL.String())
	require.Equal(t, "application/json", upstream.lastReq.Header.Get("Content-Type"))
	require.True(t, json.Valid(upstream.lastBody))
	require.Equal(t, "grok-imagine-edit", gjson.GetBytes(upstream.lastBody, "model").String())
	require.Equal(t, "edit this private image", gjson.GetBytes(upstream.lastBody, "prompt").String())
	require.True(t, strings.HasPrefix(gjson.GetBytes(upstream.lastBody, "image.image_url").String(), "data:image/png;base64,"))
placeholder

func TestForwardGrokMediaVideoGenerationReturnsUsageAndResponseID(t *testing.T) {
	t.Setenv(xai.EnvAllowUnsafeURLOverrides, "true")
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok-imagine-video-1.5","prompt":"waves","resolution":"720p","duration":10placeholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/videos/generations", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	account := &Account{
		ID:          63,
		Name:        "grok",
		Platform:    PlatformGrok,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "api-key",
			"base_url": "https://xai.test/v1",
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type":   []string{"application/json"placeholder,
			"Xai-Request-Id": []string{"xai-video-generate-req"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(`{"request_id":"video-request-123","usage":{"prompt_tokens":3,"completion_tokens":4placeholderplaceholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder

	result, err := svc.ForwardGrokMedia(context.Background(), c, account, GrokMediaEndpointVideosGenerations, "", body, "application/json")
placeholder
	require.Equal(t, "https://xai.test/v1/videos/generations", upstream.lastReq.URL.String())
	require.JSONEq(t, `{"model":"grok-imagine-video","prompt":"waves","resolution":"720p","duration":10placeholder`, string(upstream.lastBody))
	require.Equal(t, "video-request-123", result.ResponseID)
	require.Equal(t, "grok-imagine-video", result.BillingModel)
	require.Equal(t, 3, result.Usage.InputTokens)
	require.Equal(t, 4, result.Usage.OutputTokens)
	require.Equal(t, 1, result.ImageCount)
	require.Empty(t, result.ImageSize)
	require.Equal(t, 1, result.VideoCount)
	require.Equal(t, VideoBillingResolution720P, result.VideoResolution)
	require.Equal(t, 10, result.VideoDurationSeconds)
placeholder

func TestForwardGrokMediaVideoGenerationPreservesImageToVideoModel(t *testing.T) {
	t.Setenv(xai.EnvAllowUnsafeURLOverrides, "true")
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok-imagine-video-1.5","prompt":"animate","image":{"image_url":"data:image/png;base64,aW1n"placeholderplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/videos/generations", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	account := &Account{
		ID:          63,
		Name:        "grok",
		Platform:    PlatformGrok,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "api-key",
			"base_url": "https://xai.test/v1",
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"application/json"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(`{"request_id":"video-request-456"placeholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder

	result, err := svc.ForwardGrokMedia(context.Background(), c, account, GrokMediaEndpointVideosGenerations, "", body, "application/json")
placeholder
	require.Equal(t, "https://xai.test/v1/videos/generations", upstream.lastReq.URL.String())
	require.JSONEq(t, `{"model":"grok-imagine-video-1.5","prompt":"animate","image":{"image_url":"data:image/png;base64,aW1n"placeholderplaceholder`, string(upstream.lastBody))
	require.Equal(t, "video-request-456", result.ResponseID)
	require.Equal(t, "grok-imagine-video-1.5", result.BillingModel)
	// 未指定 duration 时按上游默认 8 秒计费。
	require.Equal(t, VideoBillingDefaultDurationSeconds, result.VideoDurationSeconds)
placeholder

func TestForwardGrokMediaVideoStatusUsesGETWithoutBody(t *testing.T) {
	t.Setenv(xai.EnvAllowUnsafeURLOverrides, "true")
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/videos/request-123", nil)

	account := &Account{
		ID:          62,
		Name:        "grok",
		Platform:    PlatformGrok,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "api-key",
			"base_url": "https://xai.test/v1",
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type":   []string{"application/json"placeholder,
			"Xai-Request-Id": []string{"xai-video-req"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(`{"id":"request-123","status":"completed"placeholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder

	result, err := svc.ForwardGrokMedia(context.Background(), c, account, GrokMediaEndpointVideoStatus, "request-123", nil, "")
placeholder
	require.Equal(t, "https://xai.test/v1/videos/request-123", upstream.lastReq.URL.String())
	require.Equal(t, http.MethodGet, upstream.lastReq.Method)
	require.Equal(t, "Bearer api-key", upstream.lastReq.Header.Get("Authorization"))
	require.Empty(t, upstream.lastReq.Header.Get("Content-Type"))
	require.Empty(t, upstream.lastBody)
	require.Equal(t, http.StatusOK, recorder.Code)
	require.JSONEq(t, `{"id":"request-123","status":"completed"placeholder`, recorder.Body.String())
	require.Equal(t, "xai-video-req", result.RequestID)
placeholder

func TestBindGrokMediaVideoRequestAccountUsesRequestIDStickyHash(t *testing.T) {
	ctx := context.Background()
	groupID := int64(7)
	cache := &stubGatewayCache{placeholder
	svc := &OpenAIGatewayService{cache: cacheplaceholder

	hash := GrokMediaVideoRequestSessionHash("video-request-123")
	require.NotEmpty(t, hash)
	require.NoError(t, svc.BindGrokMediaVideoRequestAccount(ctx, &groupID, "video-request-123", 63))

	accountID, err := svc.getStickySessionAccountID(ctx, &groupID, hash)
placeholder
	require.Equal(t, int64(63), accountID)
placeholder

func TestForwardGrokMedia429ReconcilesRateLimitBeforeCustomErrorBypass(t *testing.T) {
	t.Setenv(xai.EnvAllowUnsafeURLOverrides, "true")
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok-imagine","prompt":"draw a cat"placeholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	account := &Account{
		ID:          64,
		Name:        "grok",
		Platform:    PlatformGrok,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":                    "api-key",
			"base_url":                   "https://xai.test/v1",
			"custom_error_codes_enabled": true,
			"custom_error_codes":         []any{float64(http.StatusBadRequest)placeholder,
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header: http.Header{
			"Content-Type":   []string{"application/json"placeholder,
			"Xai-Request-Id": []string{"xai-error-req"placeholder,
			"Retry-After":    []string{"45"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(`{"error":{"message":"do not expose this upstream detail"placeholderplaceholder`)),
placeholderplaceholder
	repo := &grokQuotaAccountRepo{placeholder
	svc := &OpenAIGatewayService{httpUpstream: upstream, accountRepo: repoplaceholder

	result, err := svc.ForwardGrokMedia(context.Background(), c, account, GrokMediaEndpointImagesGenerations, "", body, "application/json")
placeholder
	require.Nil(t, result)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Contains(t, recorder.Body.String(), "Upstream gateway error")
	require.NotContains(t, recorder.Body.String(), "do not expose")
	require.Equal(t, 1, repo.rateLimitedCalls)
	require.Zero(t, repo.tempUnschedCalls)
	require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestForwardAsChatCompletionsForGrokStopFallsBackToXAIChatCompletions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok","messages":[{"role":"user","content":"hi"placeholder],"stream":false,"stop":"done","prompt_cache_key":"raw-client-cache-key"placeholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	c.Set("api_key", &APIKey{ID: 5101placeholder)

	account := &Account{
		ID:          51,
		Name:        "grok",
		Platform:    PlatformGrok,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token": "access-token",
			"expires_at":   time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
			"base_url":     xai.DefaultCLIBaseURL,
	placeholder,
placeholder
	repo := &grokQuotaAccountRepo{
		mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
			accountsByID: map[int64]*Account{51: accountplaceholder,
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type":                   []string{"application/json"placeholder,
			"Xai-Request-Id":                 []string{"xai-req"placeholder,
			"X-Ratelimit-Limit-Requests":     []string{"10"placeholder,
			"X-Ratelimit-Remaining-Requests": []string{"9"placeholder,
			"X-Ratelimit-Limit-Tokens":       []string{"1000"placeholder,
			"X-Ratelimit-Remaining-Tokens":   []string{"990"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(`{"id":"chatcmpl","object":"chat.completion","model":"grok-4.3","choices":[],"usage":{"prompt_tokens":1,"completion_tokens":2,"prompt_tokens_details":{"cached_tokens":1placeholderplaceholderplaceholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{
		httpUpstream:      upstream,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		accountRepo:       repo,
placeholder

	result, err := svc.ForwardAsChatCompletions(context.Background(), c, account, body, "", "")
placeholder
	require.Equal(t, xai.DefaultCLIBaseURL+"/chat/completions", upstream.lastReq.URL.String())
	require.Equal(t, "Bearer access-token", upstream.lastReq.Header.Get("Authorization"))
	require.NotEmpty(t, upstream.lastReq.Header.Get(grokConversationIDHeader))
	require.NotEqual(t, "raw-client-cache-key", upstream.lastReq.Header.Get(grokConversationIDHeader))
	require.Equal(t, "grok-4.5", gjson.GetBytes(upstream.lastBody, "model").String())
	require.False(t, gjson.GetBytes(upstream.lastBody, "prompt_cache_key").Exists())
	require.Equal(t, "grok", result.Model)
	require.Equal(t, "grok-4.5", result.UpstreamModel)
	require.Equal(t, 1, result.Usage.InputTokens)
	require.Equal(t, 2, result.Usage.OutputTokens)
	require.Equal(t, 1, result.Usage.CacheReadInputTokens)
	require.NotNil(t, repo.updates[51][grokQuotaSnapshotExtraKey])
	require.Equal(t, http.StatusOK, recorder.Code)
placeholder

func TestForwardGrokResponsesStreamingUsesXAIResponsesAndSnapshots(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok","input":"hi","stream":true,"reasoning_effort":"high"placeholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("OpenAI-Beta", "responses=experimental")
	c.Set("api_key", &APIKey{ID: 5201placeholder)

	account := &Account{
		ID:          52,
		Name:        "grok",
		Platform:    PlatformGrok,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token": "access-token",
			"expires_at":   time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
			"base_url":     xai.DefaultCLIBaseURL,
	placeholder,
placeholder
	repo := &grokQuotaAccountRepo{
		mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
			accountsByID: map[int64]*Account{52: accountplaceholder,
	placeholder,
placeholder
	upstreamBody := strings.Join([]string{
		`data: {"type":"response.output_text.delta","sequence_number":0,"delta":"ok"placeholder`,
		"",
		`data: {"type":"response.completed","sequence_number":1,"response":{"id":"resp_grok","model":"grok-4.3","usage":{"input_tokens":5,"output_tokens":3,"input_tokens_details":{"cached_tokens":2placeholderplaceholderplaceholderplaceholder`,
		"",
placeholder, "\n")
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type":                   []string{"text/event-stream"placeholder,
			"Xai-Request-Id":                 []string{"xai-stream-req"placeholder,
			"X-Ratelimit-Limit-Requests":     []string{"10"placeholder,
			"X-Ratelimit-Remaining-Requests": []string{"8"placeholder,
			"X-Ratelimit-Limit-Tokens":       []string{"1000"placeholder,
			"X-Ratelimit-Remaining-Tokens":   []string{"990"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(upstreamBody)),
placeholderplaceholder
	svc := &OpenAIGatewayService{
		httpUpstream:      upstream,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		accountRepo:       repo,
placeholder

	result, err := svc.forwardGrokResponses(context.Background(), c, account, body, "grok", true, time.Now())
placeholder
	require.Equal(t, xai.DefaultCLIBaseURL+"/responses", upstream.lastReq.URL.String())
	require.Equal(t, "Bearer access-token", upstream.lastReq.Header.Get("Authorization"))
	require.Equal(t, "responses=experimental", upstream.lastReq.Header.Get("OpenAI-Beta"))
	require.Equal(t, "grok-4.5", gjson.GetBytes(upstream.lastBody, "model").String())
	require.NotEmpty(t, gjson.GetBytes(upstream.lastBody, "prompt_cache_key").String())
	require.Equal(t, gjson.GetBytes(upstream.lastBody, "prompt_cache_key").String(), upstream.lastReq.Header.Get(grokConversationIDHeader))
	require.Equal(t, "web_search", gjson.GetBytes(upstream.lastBody, "tools.0.type").String())
	require.Equal(t, "x_search", gjson.GetBytes(upstream.lastBody, "tools.1.type").String())
	require.Equal(t, "none", gjson.GetBytes(upstream.lastBody, "tool_choice").String())
	require.Equal(t, "high", gjson.GetBytes(upstream.lastBody, "reasoning_effort").String())
	require.True(t, gjson.GetBytes(upstream.lastBody, "stream").Bool())
	require.True(t, result.Stream)
	require.Equal(t, "resp_grok", result.ResponseID)
	require.Equal(t, "xai-stream-req", result.RequestID)
	require.Equal(t, 5, result.Usage.InputTokens)
	require.Equal(t, 3, result.Usage.OutputTokens)
	require.Equal(t, 2, result.Usage.CacheReadInputTokens)
	require.NotNil(t, result.ReasoningEffort)
	require.Equal(t, "high", *result.ReasoningEffort)
	require.Contains(t, recorder.Header().Get("Content-Type"), "text/event-stream")
	require.Contains(t, recorder.Body.String(), "response.output_text.delta")
	require.NotNil(t, repo.updates[52][grokQuotaSnapshotExtraKey])
placeholder

func TestForwardGrokResponsesNonStreamingUsesCacheIdentityAndCachedUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok","input":"hi","stream":false,"tools":[{"type":"namespace","name":"client_tools"placeholder],"tool_choice":{"type":"namespace","name":"client_tools"placeholderplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("api_key", &APIKey{ID: 5202placeholder)

	account := &Account{
		ID:          56,
		Name:        "grok",
		Platform:    PlatformGrok,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token": "access-token",
			"expires_at":   time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
			"base_url":     xai.DefaultCLIBaseURL,
	placeholder,
placeholder
	repo := &grokQuotaAccountRepo{
		mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
			accountsByID: map[int64]*Account{56: accountplaceholder,
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type":   []string{"application/json"placeholder,
			"Xai-Request-Id": []string{"xai-non-stream-req"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(`{"id":"resp_grok_non_stream","object":"response","model":"grok-4.3","status":"completed","output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"ok"placeholder]placeholder],"usage":{"input_tokens":7,"output_tokens":2,"total_tokens":9,"input_tokens_details":{"cached_tokens":4placeholderplaceholderplaceholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{
		httpUpstream:      upstream,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		accountRepo:       repo,
placeholder

	result, err := svc.forwardGrokResponses(context.Background(), c, account, body, "grok", false, time.Now())
placeholder
	require.NotNil(t, result)
	require.False(t, result.Stream)
	require.Equal(t, "resp_grok_non_stream", result.ResponseID)
	require.Equal(t, 7, result.Usage.InputTokens)
	require.Equal(t, 2, result.Usage.OutputTokens)
	require.Equal(t, 4, result.Usage.CacheReadInputTokens)
	require.Equal(t, "grok-4.5", gjson.GetBytes(upstream.lastBody, "model").String())
	identity := gjson.GetBytes(upstream.lastBody, "prompt_cache_key").String()
	require.NotEmpty(t, identity)
	require.Equal(t, identity, upstream.lastReq.Header.Get(grokConversationIDHeader))
	// The sanitizer drops this unsupported client tool, but its explicit intent
	// must still prevent native cache-routing tools from being injected.
	require.False(t, gjson.GetBytes(upstream.lastBody, "tools").Exists())
	require.False(t, gjson.GetBytes(upstream.lastBody, "tool_choice").Exists())
	require.Equal(t, "resp_grok_non_stream", gjson.Get(recorder.Body.String(), "id").String())
placeholder

func TestForwardGrokResponsesFailoverKeepsCacheIdentityAcrossAccounts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok","input":[{"role":"user","content":"stable prefix"placeholder],"stream":falseplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))
	c.Set("api_key", &APIKey{ID: 5203placeholder)

	newAccount := func(id int64, token string) *Account {
	placeholder
			ID:          id,
			Name:        fmt.Sprintf("grok-%d", id),
			Platform:    PlatformGrok,
			Type:        AccountTypeOAuth,
			Concurrency: 1,
	placeholder
				"access_token": token,
				"expires_at":   time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
				"base_url":     xai.DefaultCLIBaseURL,
		placeholder,
	placeholder
placeholder
	firstAccount := newAccount(58, "access-token-a")
	secondAccount := newAccount(59, "access-token-b")
	repo := &grokQuotaAccountRepo{
		mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
			accountsByID: map[int64]*Account{58: firstAccount, 59: secondAccountplaceholder,
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{responses: []*http.Response{
		{
			StatusCode: http.StatusServiceUnavailable,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"temporary"placeholderplaceholder`)),
	placeholder,
		{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"id":"resp_after_failover","object":"response","model":"grok-4.3","status":"completed","output":[],"usage":{"input_tokens":5,"output_tokens":1placeholderplaceholder`)),
	placeholder,
placeholderplaceholder
	svc := &OpenAIGatewayService{
		httpUpstream:      upstream,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		accountRepo:       repo,
placeholder

	_, err := svc.forwardGrokResponses(context.Background(), c, firstAccount, body, "grok", false, time.Now())
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)

	result, err := svc.forwardGrokResponses(context.Background(), c, secondAccount, body, "grok", false, time.Now())
placeholder
	require.NotNil(t, result)
	require.Len(t, upstream.requests, 2)
	require.Len(t, upstream.bodies, 2)
	firstIdentity := gjson.GetBytes(upstream.bodies[0], "prompt_cache_key").String()
	secondIdentity := gjson.GetBytes(upstream.bodies[1], "prompt_cache_key").String()
	require.NotEmpty(t, firstIdentity)
	require.Equal(t, firstIdentity, secondIdentity)
	require.Equal(t, firstIdentity, upstream.requests[0].Header.Get(grokConversationIDHeader))
	require.Equal(t, secondIdentity, upstream.requests[1].Header.Get(grokConversationIDHeader))
	require.Equal(t, "Bearer access-token-a", upstream.requests[0].Header.Get("Authorization"))
	require.Equal(t, "Bearer access-token-b", upstream.requests[1].Header.Get("Authorization"))
placeholder

func TestForwardAsChatCompletionsForGrokStreamingStopFallsBackToRawXAIChatCompletions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok","messages":[{"role":"user","content":"hi"placeholder],"stream":true,"stop":"done"placeholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set(grokConversationIDHeader, "native-client-conversation")
	c.Set("api_key", &APIKey{ID: 5301placeholder)

	account := &Account{
		ID:          53,
		Name:        "grok",
		Platform:    PlatformGrok,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token": "access-token",
			"expires_at":   time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
			"base_url":     xai.DefaultCLIBaseURL,
	placeholder,
placeholder
	repo := &grokQuotaAccountRepo{
		mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
			accountsByID: map[int64]*Account{53: accountplaceholder,
	placeholder,
placeholder
	upstreamBody := strings.Join([]string{
		`data: {"id":"chatcmpl_grok","object":"chat.completion.chunk","model":"grok-4.3","choices":[{"index":0,"delta":{"content":"ok"placeholderplaceholder]placeholder`,
		"",
		`data: {"id":"chatcmpl_grok","object":"chat.completion.chunk","model":"grok-4.3","choices":[],"usage":{"prompt_tokens":6,"completion_tokens":4,"total_tokens":10,"prompt_tokens_details":{"cached_tokens":1placeholderplaceholderplaceholder`,
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type":                   []string{"text/event-stream"placeholder,
			"X-Request-Id":                   []string{"chat-stream-req"placeholder,
			"X-Ratelimit-Limit-Requests":     []string{"10"placeholder,
			"X-Ratelimit-Remaining-Requests": []string{"7"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(upstreamBody)),
placeholderplaceholder
	svc := &OpenAIGatewayService{
		cfg:               rawChatCompletionsTestConfig(),
		httpUpstream:      upstream,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		accountRepo:       repo,
placeholder

	result, err := svc.ForwardAsChatCompletions(context.Background(), c, account, body, "", "")
placeholder
	require.Equal(t, xai.DefaultCLIBaseURL+"/chat/completions", upstream.lastReq.URL.String())
	require.Equal(t, "Bearer access-token", upstream.lastReq.Header.Get("Authorization"))
	require.Equal(t, "text/event-stream", upstream.lastReq.Header.Get("Accept"))
	require.Equal(t, "sub2api-grok/1.0", upstream.lastReq.Header.Get("User-Agent"))
	require.Equal(t, grokCLIVersion, upstream.lastReq.Header.Get("X-Grok-Client-Version"))
	require.NotEmpty(t, upstream.lastReq.Header.Get(grokConversationIDHeader))
	require.NotEqual(t, "native-client-conversation", upstream.lastReq.Header.Get(grokConversationIDHeader))
	require.Equal(t, "grok-4.5", gjson.GetBytes(upstream.lastBody, "model").String())
	require.True(t, gjson.GetBytes(upstream.lastBody, "stream_options.include_usage").Bool())
	require.True(t, result.Stream)
	require.Equal(t, 6, result.Usage.InputTokens)
	require.Equal(t, 4, result.Usage.OutputTokens)
	require.Equal(t, 1, result.Usage.CacheReadInputTokens)
	require.Contains(t, recorder.Body.String(), "data: [DONE]")
	require.NotNil(t, repo.updates[53][grokQuotaSnapshotExtraKey])
placeholder

func TestForwardAsChatCompletionsForGrokComposerBridgesImageInput(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok-composer-2.5-fast","messages":[{"role":"system","content":"You are concise."placeholder,{"role":"user","content":[{"type":"text","text":"What is shown?"placeholder,{"type":"image_url","image_url":{"url":"data:image/png;base64,QUJD"placeholderplaceholder]placeholder],"stream":falseplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("api_key", &APIKey{ID: 5501placeholder)

	account := &Account{
		ID:          55,
		Name:        "grok",
		Platform:    PlatformGrok,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token": "access-token",
			"expires_at":   time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
			"base_url":     xai.DefaultCLIBaseURL,
	placeholder,
placeholder
	repo := &grokQuotaAccountRepo{
		mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
			accountsByID: map[int64]*Account{55: accountplaceholder,
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{responses: []*http.Response{
		{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholder, "xai-request-id": []string{"vision-req"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"id":"resp_vision","object":"response","model":"grok-build-0.1","output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"A small diagram with ABC letters."placeholder]placeholder],"usage":{"input_tokens":11,"output_tokens":7,"total_tokens":18placeholderplaceholder`)),
	placeholder,
		{
			StatusCode: http.StatusOK,
			Header: http.Header{
				"Content-Type":                   []string{"application/json"placeholder,
				"X-Request-Id":                   []string{"composer-req"placeholder,
				"X-Ratelimit-Limit-Requests":     []string{"10"placeholder,
				"X-Ratelimit-Remaining-Requests": []string{"9"placeholder,
				"X-Ratelimit-Limit-Tokens":       []string{"1000"placeholder,
				"X-Ratelimit-Remaining-Tokens":   []string{"980"placeholder,
		placeholder,
			Body: io.NopCloser(strings.NewReader(`{"id":"chatcmpl_composer","object":"chat.completion","model":"grok-composer-2.5-fast","choices":[{"index":0,"message":{"role":"assistant","content":"It shows ABC."placeholder,"finish_reason":"stop"placeholder],"usage":{"prompt_tokens":3,"completion_tokens":5,"total_tokens":8placeholderplaceholder`)),
	placeholder,
placeholderplaceholder
	svc := &OpenAIGatewayService{
		cfg:               rawChatCompletionsTestConfig(),
		httpUpstream:      upstream,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		accountRepo:       repo,
placeholder

	result, err := svc.ForwardAsChatCompletions(context.Background(), c, account, body, "", "")
placeholder
	require.NotNil(t, result)
	require.Len(t, upstream.requests, 2)
	require.Equal(t, xai.DefaultCLIBaseURL+"/responses", upstream.requests[0].URL.String())
	require.Empty(t, upstream.requests[0].Header.Get(grokConversationIDHeader))
	require.Equal(t, "grok-build-0.1", gjson.GetBytes(upstream.bodies[0], "model").String())
	require.Equal(t, "input_image", gjson.GetBytes(upstream.bodies[0], "input.0.content.1.type").String())
	require.Equal(t, xai.DefaultCLIBaseURL+"/chat/completions", upstream.requests[1].URL.String())
	require.NotEmpty(t, upstream.requests[1].Header.Get(grokConversationIDHeader))
	require.Equal(t, "grok-composer-2.5-fast", gjson.GetBytes(upstream.bodies[1], "model").String())
	require.False(t, strings.Contains(string(upstream.bodies[1]), "image_url"))
	require.Contains(t, gjson.GetBytes(upstream.bodies[1], "messages.1.content").String(), "Image 1 description")
	require.Contains(t, gjson.GetBytes(upstream.bodies[1], "messages.1.content").String(), "A small diagram with ABC letters.")
	require.Equal(t, 14, result.Usage.InputTokens)
	require.Equal(t, 12, result.Usage.OutputTokens)
	require.Equal(t, "It shows ABC.", gjson.Get(recorder.Body.String(), "choices.0.message.content").String())
	require.NotNil(t, repo.updates[55][grokQuotaSnapshotExtraKey])
placeholder

func TestForwardAsAnthropicForGrokUsesXAIResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok","max_tokens":32,"stream":false,"messages":[{"role":"user","content":"hi"placeholder]placeholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	c.Set("api_key", &APIKey{ID: 5401placeholder)
	c.Request.Header.Set("OpenAI-Beta", "grok-experimental")
	c.Request.Header.Set("originator", "opencode")

	account := &Account{
		ID:          54,
		Name:        "grok",
		Platform:    PlatformGrok,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token": "access-token",
			"expires_at":   time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
			"base_url":     xai.DefaultCLIBaseURL,
	placeholder,
placeholder
	repo := &grokQuotaAccountRepo{
		mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
			accountsByID: map[int64]*Account{54: accountplaceholder,
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: grokMessagesSSECompletedResponse("resp_grok_messages", 3)placeholder
	svc := &OpenAIGatewayService{
		httpUpstream:      upstream,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		accountRepo:       repo,
placeholder

	result, err := svc.ForwardAsAnthropic(context.Background(), c, account, body, "", "")
placeholder
	require.Equal(t, xai.DefaultCLIBaseURL+"/responses", upstream.lastReq.URL.String())
	require.Equal(t, "Bearer access-token", upstream.lastReq.Header.Get("Authorization"))
	require.Equal(t, "sub2api-grok/1.0", upstream.lastReq.Header.Get("User-Agent"))
	require.Equal(t, grokCLIVersion, upstream.lastReq.Header.Get("X-Grok-Client-Version"))
	require.Equal(t, "grok-experimental", upstream.lastReq.Header.Get("OpenAI-Beta"))
	require.Empty(t, upstream.lastReq.Header.Get("originator"))
	require.Empty(t, upstream.lastReq.Header.Get("version"))
	require.Equal(t, "grok-4.5", gjson.GetBytes(upstream.lastBody, "model").String())
	require.NotEmpty(t, gjson.GetBytes(upstream.lastBody, "prompt_cache_key").String())
	require.Equal(t, gjson.GetBytes(upstream.lastBody, "prompt_cache_key").String(), upstream.lastReq.Header.Get(grokConversationIDHeader))
	require.Equal(t, "web_search", gjson.GetBytes(upstream.lastBody, "tools.0.type").String())
	require.Equal(t, "x_search", gjson.GetBytes(upstream.lastBody, "tools.1.type").String())
	require.Equal(t, "none", gjson.GetBytes(upstream.lastBody, "tool_choice").String())
	require.Empty(t, upstream.lastReq.Header.Get("session_id"))
	require.True(t, gjson.GetBytes(upstream.lastBody, "stream").Bool())
	require.NotContains(t, string(upstream.lastBody), "chatgpt.com")
	require.Equal(t, "grok", result.Model)
	require.Equal(t, "grok-4.5", result.UpstreamModel)
	require.Equal(t, 5, result.Usage.InputTokens)
	require.Equal(t, 2, result.Usage.OutputTokens)
	require.Equal(t, 3, result.Usage.CacheReadInputTokens)
	require.Contains(t, recorder.Body.String(), `"type":"message"`)
	require.Equal(t, int64(3), gjson.Get(recorder.Body.String(), "usage.cache_read_input_tokens").Int())
	require.Contains(t, recorder.Body.String(), "ok")
placeholder

func TestForwardAsAnthropicForGrokStreamingPreservesCacheUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok","max_tokens":32,"stream":true,"messages":[{"role":"user","content":"hi"placeholder]placeholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	c.Set("api_key", &APIKey{ID: 5402placeholder)

	account := &Account{
		ID:          57,
		Name:        "grok",
		Platform:    PlatformGrok,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token": "access-token",
			"expires_at":   time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
			"base_url":     xai.DefaultCLIBaseURL,
	placeholder,
placeholder
	repo := &grokQuotaAccountRepo{
		mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
			accountsByID: map[int64]*Account{57: accountplaceholder,
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: grokMessagesSSECompletedResponse("resp_grok_messages_stream", 2)placeholder
	svc := &OpenAIGatewayService{
		httpUpstream:      upstream,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		accountRepo:       repo,
placeholder

	result, err := svc.ForwardAsAnthropic(context.Background(), c, account, body, "", "")
placeholder
	require.NotNil(t, result)
	require.Equal(t, 2, result.Usage.CacheReadInputTokens)
	identity := gjson.GetBytes(upstream.lastBody, "prompt_cache_key").String()
	require.NotEmpty(t, identity)
	require.Equal(t, identity, upstream.lastReq.Header.Get(grokConversationIDHeader))
	require.Contains(t, recorder.Header().Get("Content-Type"), "text/event-stream")
	require.Contains(t, recorder.Body.String(), `"cache_read_input_tokens":2`)
placeholder

func grokMessagesSSECompletedResponse(responseID string, cachedTokens int) *http.Response {
	body := strings.Join([]string{
		fmt.Sprintf(`data: {"type":"response.completed","response":{"id":%q,"object":"response","model":"grok-4.3","status":"completed","output":[{"type":"message","id":"msg_1","role":"assistant","status":"completed","content":[{"type":"output_text","text":"ok"placeholder]placeholder],"usage":{"input_tokens":5,"output_tokens":2,"total_tokens":7,"input_tokens_details":{"cached_tokens":%dplaceholderplaceholderplaceholderplaceholder`, responseID, cachedTokens),
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(body)),
placeholder
placeholder

func TestHandleGrokAccountUpstreamErrorTempUnschedulesNonRateLimitStates(t *testing.T) {
	tests := []struct {
		name            string
		status          int
		headers         http.Header
		wantReason      string
		wantMinCooldown time.Duration
		wantMaxCooldown time.Duration
placeholder{
		{
			name:            "unauthorized reauth",
			status:          http.StatusUnauthorized,
			wantReason:      "grok oauth token unauthorized",
			wantMinCooldown: 10*time.Minute - time.Second,
			wantMaxCooldown: 10*time.Minute + time.Second,
	placeholder,
		{
			name:            "forbidden entitlement",
			status:          http.StatusForbidden,
			wantReason:      "grok entitlement or subscription tier denied",
			wantMinCooldown: 30*time.Minute - time.Second,
			wantMaxCooldown: 30*time.Minute + time.Second,
	placeholder,
		{
			name:            "upstream temporary error",
			status:          http.StatusInternalServerError,
			wantReason:      "grok upstream temporary error",
			wantMinCooldown: 2*time.Minute - time.Second,
			wantMaxCooldown: 2*time.Minute + time.Second,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := &Account{ID: 61, Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder
			repo := &grokQuotaAccountRepo{placeholder
			svc := &OpenAIGatewayService{accountRepo: repoplaceholder
			before := time.Now()

			svc.handleGrokAccountUpstreamError(context.Background(), account, tt.status, tt.headers, nil)

			require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
			require.Equal(t, 1, repo.tempUnschedCalls)
			require.Zero(t, repo.rateLimitedCalls)
			require.Equal(t, account.ID, repo.lastTempUnschedID)
			require.Equal(t, tt.wantReason, repo.lastTempUnschedReason)
			require.True(t, repo.lastTempUnschedUntil.After(before.Add(tt.wantMinCooldown)))
			require.True(t, repo.lastTempUnschedUntil.Before(before.Add(tt.wantMaxCooldown)))
	placeholder)
placeholder
placeholder

func TestHandleGrokAccountUpstreamError429SetsRateLimitedFromRetryAfter(t *testing.T) {
	account := &Account{ID: 61, Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder
	repo := &grokQuotaAccountRepo{placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder
	before := time.Now()

	svc.handleGrokAccountUpstreamError(context.Background(), account, http.StatusTooManyRequests, http.Header{"Retry-After": []string{"45"placeholderplaceholder, nil)

	require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
	require.Equal(t, 1, repo.rateLimitedCalls)
	require.Equal(t, account.ID, repo.lastRateLimitedID)
	require.WithinDuration(t, before.Add(45*time.Second), repo.lastRateLimitResetAt, time.Second)
	require.Zero(t, repo.tempUnschedCalls)
placeholder

func TestHandleGrokAccountUpstreamError429UsesLatestExhaustedWindowReset(t *testing.T) {
	now := time.Now()
	requestReset := now.Add(10 * time.Minute).Truncate(time.Second)
	tokenReset := now.Add(20 * time.Minute).Truncate(time.Second)
	headers := http.Header{
		"X-Ratelimit-Limit-Requests":     []string{"10"placeholder,
		"X-Ratelimit-Remaining-Requests": []string{"0"placeholder,
		"X-Ratelimit-Reset-Requests":     []string{fmt.Sprintf("%d", requestReset.Unix())placeholder,
		"X-Ratelimit-Limit-Tokens":       []string{"1000"placeholder,
		"X-Ratelimit-Remaining-Tokens":   []string{"0"placeholder,
		"X-Ratelimit-Reset-Tokens":       []string{fmt.Sprintf("%d", tokenReset.Unix())placeholder,
placeholder
	account := &Account{ID: 62, Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder
	repo := &grokQuotaAccountRepo{placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder

	svc.handleGrokAccountUpstreamError(context.Background(), account, http.StatusTooManyRequests, headers, nil)

	require.Equal(t, 1, repo.rateLimitedCalls)
	require.WithinDuration(t, tokenReset, repo.lastRateLimitResetAt, time.Second)
	require.Zero(t, repo.tempUnschedCalls)
placeholder

func TestHandleGrokAccountUpstreamError429UsesFallbackReset(t *testing.T) {
	account := &Account{ID: 63, Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder
	repo := &grokQuotaAccountRepo{placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder
	before := time.Now()

	svc.handleGrokAccountUpstreamError(context.Background(), account, http.StatusTooManyRequests, nil, nil)

	require.Equal(t, 1, repo.rateLimitedCalls)
	require.WithinDuration(t, before.Add(grokRateLimitFallbackCooldown), repo.lastRateLimitResetAt, time.Second)
	require.Zero(t, repo.tempUnschedCalls)
placeholder

func TestGrokRateLimitResetAtUsesFutureWindowAfterRetryAfterExpires(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	observedAt := now.Add(-2 * time.Minute)
	windowReset := now.Add(15 * time.Minute)
	retryAfter := 30
	snapshot := &xai.QuotaSnapshot{
		StatusCode:        http.StatusTooManyRequests,
		UpdatedAt:         observedAt.Format(time.RFC3339),
		RetryAfterSeconds: &retryAfter,
		Requests: &xai.QuotaWindow{
			Limit:     grokInt64PtrForTest(10),
			Remaining: grokInt64PtrForTest(0),
			ResetUnix: grokInt64PtrForTest(windowReset.Unix()),
	placeholder,
placeholder

	resetAt, limited := grokRateLimitResetAt(snapshot, now)

	require.True(t, limited)
	require.WithinDuration(t, windowReset, resetAt, time.Second)
placeholder

func TestHandleGrokAccountUpstreamError429DoesNotShortenExistingPause(t *testing.T) {
	existingUntil := time.Now().Add(15 * time.Minute)
	account := &Account{
		ID:                      64,
		Platform:                PlatformGrok,
		Type:                    AccountTypeOAuth,
		TempUnschedulableUntil:  &existingUntil,
		TempUnschedulableReason: "existing pause",
placeholder
	repo := &grokQuotaAccountRepo{placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder

	svc.handleGrokAccountUpstreamError(context.Background(), account, http.StatusTooManyRequests, http.Header{"Retry-After": []string{"45"placeholderplaceholder, nil)

	require.Equal(t, 1, repo.rateLimitedCalls)
	require.WithinDuration(t, time.Now().Add(45*time.Second), repo.lastRateLimitResetAt, time.Second)
	require.Zero(t, repo.tempUnschedCalls)
	value, ok := svc.openaiAccountRuntimeBlockUntil.Load(account.ID)
	require.True(t, ok)
	runtimeUntil, ok := value.(time.Time)
	require.True(t, ok)
	require.WithinDuration(t, existingUntil, runtimeUntil, time.Second)
placeholder

func TestUpdateGrokUsageSnapshotExhaustedSuccessBypassesThrottleAndSetsRateLimited(t *testing.T) {
	account := &Account{ID: 65, Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder
	repo := &grokQuotaAccountRepo{placeholder
	svc := &OpenAIGatewayService{
		accountRepo:           repo,
		codexSnapshotThrottle: newAccountWriteThrottle(time.Hour),
placeholder
	now := time.Now()

	// Consume the normal snapshot write allowance first.
	svc.updateGrokUsageSnapshot(context.Background(), account, &xai.QuotaSnapshot{
		StatusCode: http.StatusOK,
		Requests: &xai.QuotaWindow{
			Limit:     grokInt64PtrForTest(10),
			Remaining: grokInt64PtrForTest(9),
	placeholder,
		UpdatedAt: now.UTC().Format(time.RFC3339),
placeholder)
	resetAt := now.Add(30 * time.Minute).Truncate(time.Second)
	svc.updateGrokUsageSnapshot(context.Background(), account, &xai.QuotaSnapshot{
		StatusCode: http.StatusOK,
		Requests: &xai.QuotaWindow{
			Limit:     grokInt64PtrForTest(10),
			Remaining: grokInt64PtrForTest(0),
			ResetUnix: grokInt64PtrForTest(resetAt.Unix()),
			ResetAt:   resetAt.UTC().Format(time.RFC3339),
	placeholder,
		UpdatedAt: now.UTC().Format(time.RFC3339),
placeholder)

	require.Equal(t, 2, repo.updateCalls)
	require.Equal(t, 1, repo.rateLimitedCalls)
	require.Equal(t, account.ID, repo.lastRateLimitedID)
	require.WithinDuration(t, resetAt, repo.lastRateLimitResetAt, time.Second)
	require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestUpdateGrokUsageSnapshotAvailableSuccessDoesNotSetRateLimited(t *testing.T) {
	repo := &grokQuotaAccountRepo{placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder
	account := &Account{ID: 66, Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder

	svc.updateGrokUsageSnapshot(context.Background(), account, &xai.QuotaSnapshot{
		StatusCode: http.StatusOK,
		Requests: &xai.QuotaWindow{
			Limit:     grokInt64PtrForTest(10),
			Remaining: grokInt64PtrForTest(1),
	placeholder,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
placeholder)

	require.Equal(t, 1, repo.updateCalls)
	require.Zero(t, repo.rateLimitedCalls)
placeholder

func TestUpdateGrokUsageSnapshotExhaustedSuccessWithoutResetUsesFallback(t *testing.T) {
	repo := &grokQuotaAccountRepo{placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder
	account := &Account{ID: 67, Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder
	before := time.Now()

	svc.updateGrokUsageSnapshot(context.Background(), account, &xai.QuotaSnapshot{
		StatusCode: http.StatusOK,
		Tokens: &xai.QuotaWindow{
			Limit:     grokInt64PtrForTest(2_000_000),
			Remaining: grokInt64PtrForTest(0),
	placeholder,
		UpdatedAt: before.UTC().Format(time.RFC3339),
placeholder)

	require.Equal(t, 1, repo.rateLimitedCalls)
	require.WithinDuration(t, before.Add(grokRateLimitFallbackCooldown), repo.lastRateLimitResetAt, time.Second)
	stored, ok := repo.updates[account.ID][grokQuotaSnapshotExtraKey].(*xai.QuotaSnapshot)
	require.True(t, ok)
	require.NotNil(t, stored.Tokens.ResetUnix)
	paused, _ := shouldAutoPauseGrokQuotaWindow("tokens", stored.Tokens, before.Add(time.Second))
	require.True(t, paused)
	paused, _ = shouldAutoPauseGrokQuotaWindow("tokens", stored.Tokens, repo.lastRateLimitResetAt.Add(time.Second))
	require.False(t, paused)
placeholder

func TestOpenAIWSHTTPBridgeGrok429PersistsRateLimit(t *testing.T) {
	repo := &grokQuotaAccountRepo{placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{"Retry-After": []string{"45"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"rate limited"placeholderplaceholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{accountRepo: repo, httpUpstream: upstreamplaceholder
	account := &Account{ID: 68, Platform: PlatformGrok, Type: AccountTypeOAuth, Concurrency: 1placeholder
	before := time.Now()

	result, err := svc.proxyOpenAIWSHTTPBridgeTurn(
		context.Background(), nil, account, "token",
		[]byte(`{"type":"response.create","model":"grok-4.3","input":"hi"placeholder`),
		64, "grok-4.3", "", "", "", "cache-id", 1,
		func([]byte) error { return nil placeholder,
	)

placeholder
	require.Nil(t, result)
	require.Equal(t, 1, repo.rateLimitedCalls)
	require.WithinDuration(t, before.Add(45*time.Second), repo.lastRateLimitResetAt, time.Second)
	require.Zero(t, repo.tempUnschedCalls)
	require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestOpenAIWSHTTPBridgeGrokExhaustedSuccessPersistsRateLimit(t *testing.T) {
	repo := &grokQuotaAccountRepo{placeholder
	resetAt := time.Now().Add(20 * time.Minute).UTC().Truncate(time.Second)
	resp := grokMessagesSSECompletedResponse("resp_ws_limited", 0)
	resp.Header.Set("X-Ratelimit-Limit-Requests", "10")
	resp.Header.Set("X-Ratelimit-Remaining-Requests", "0")
	resp.Header.Set("X-Ratelimit-Reset-Requests", fmt.Sprintf("%d", resetAt.Unix()))
	upstream := &httpUpstreamRecorder{resp: respplaceholder
	svc := &OpenAIGatewayService{accountRepo: repo, httpUpstream: upstreamplaceholder
	account := &Account{ID: 69, Platform: PlatformGrok, Type: AccountTypeOAuth, Concurrency: 1placeholder

	result, err := svc.proxyOpenAIWSHTTPBridgeTurn(
		context.Background(), nil, account, "token",
		[]byte(`{"type":"response.create","model":"grok-4.3","input":"hi"placeholder`),
		64, "grok-4.3", "", "", "", "cache-id", 1,
		func([]byte) error { return nil placeholder,
	)

placeholder
	require.NotNil(t, result)
	require.Equal(t, 1, repo.rateLimitedCalls)
	require.WithinDuration(t, resetAt, repo.lastRateLimitResetAt, time.Second)
	require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestFailoverOpenAIUpstreamHTTPErrorUsesOnlyGrokRateLimitPolicy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &grokQuotaAccountRepo{placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder
	account := &Account{ID: 70, Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{"Retry-After": []string{"45"placeholderplaceholder,
placeholder
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	failoverErr := svc.failoverOpenAIUpstreamHTTPError(
		context.Background(), c, account, resp,
		[]byte(`{"error":{"message":"rate limited"placeholderplaceholder`), "rate limited", "grok-4.3",
	)

	require.NotNil(t, failoverErr)
	require.Equal(t, 1, repo.rateLimitedCalls)
	require.Zero(t, repo.tempUnschedCalls)
placeholder
