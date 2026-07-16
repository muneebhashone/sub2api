package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestOpenAIRequestView_ExtractsRawScalars(t *testing.T) {
	view := newOpenAIRequestView([]byte(`{"model":" gpt-5 ","stream":true,"prompt_cache_key":" ses-1 ","previous_response_id":" resp-1 ","service_tier":" fast ","reasoning":{"effort":" medium "placeholderplaceholder`))

	require.Equal(t, "gpt-5", view.Model)
	require.True(t, view.Stream)
	require.Equal(t, "ses-1", view.PromptCacheKey)
	require.Equal(t, "resp-1", view.PreviousResponseID)
	require.Equal(t, "fast", view.ServiceTier)
	require.Equal(t, "medium", view.ReasoningEffort)
placeholder

func TestOpenAIRequestView_ExtractsFieldsAfterLargeInput(t *testing.T) {
	body := []byte(`{"model":"gpt-5","input":[{"content":"` + strings.Repeat("payload", 1024) + `"placeholder],"stream":true,"prompt_cache_key":"session-1","previous_response_id":"resp-1","service_tier":"flex","reasoning":{"effort":"high"placeholderplaceholder`)

	view := newOpenAIRequestView(body)

	require.Equal(t, "gpt-5", view.Model)
	require.True(t, view.Stream)
	require.Equal(t, "session-1", view.PromptCacheKey)
	require.Equal(t, "resp-1", view.PreviousResponseID)
	require.Equal(t, "flex", view.ServiceTier)
	require.Equal(t, "high", view.ReasoningEffort)
placeholder

func TestOpenAIRequestView_KeepsFirstDuplicateField(t *testing.T) {
	view := newOpenAIRequestView([]byte(`{"model":"gpt-5","model":"gpt-5.1","reasoning":{"effort":"low"placeholder,"reasoning":{"effort":"high"placeholderplaceholder`))

	require.Equal(t, "gpt-5", view.Model)
	require.Equal(t, "low", view.ReasoningEffort)
placeholder

func TestOpenAIRequestView_KeepsLenientPrefixExtraction(t *testing.T) {
	view := newOpenAIRequestView([]byte(`{"model":"gpt-5","stream":true,"input":[`))

	require.Equal(t, "gpt-5", view.Model)
	require.True(t, view.Stream)
placeholder

func TestOpenAIRequestView_DecodeKeepsFullMapBehavior(t *testing.T) {
	view := newOpenAIRequestView([]byte(`{"model":"gpt-5","stream":true,"input":[{"type":"message","content":"hi"placeholder]placeholder`))

	reqBody, err := view.Decode(nil)
placeholder
	require.Equal(t, "gpt-5", reqBody["model"])
	require.IsType(t, []any{placeholder, reqBody["input"])
placeholder

func TestOpenAIRequestView_ApplyPatches(t *testing.T) {
	view := newOpenAIRequestView([]byte(`{"model":"gpt-5","previous_response_id":"resp_1","reasoning":{"effort":"minimal"placeholder,"input":[{"type":"message","content":"hi"placeholder]placeholder`))
	view.MarkPatchSet("model", "gpt-5.1")
	view.MarkPatchDelete("previous_response_id")
	view.MarkPatchSet("reasoning.effort", "none")

	patched, err := view.ApplyPatches()
placeholder
	require.JSONEq(t, `{"model":"gpt-5.1","reasoning":{"effort":"none"placeholder,"input":[{"type":"message","content":"hi"placeholder]placeholder`, string(patched))
placeholder

func TestOpenAIRequestView_RejectsEscapedPatchPath(t *testing.T) {
	view := newOpenAIRequestView([]byte(`{"metadata":{"user.id":"old"placeholderplaceholder`))
	view.MarkPatchSet(`metadata.user\.id`, "new")

	require.False(t, view.HasPatches())
	_, err := view.ApplyPatches()
placeholder
placeholder

func TestOpenAIRequestView_ApplyPatchesDisabled(t *testing.T) {
	view := newOpenAIRequestView([]byte(`{"model":"gpt-5"placeholder`))
	view.MarkPatchSet("model", "gpt-5.1")
	view.DisablePatches()

	_, err := view.ApplyPatches()
placeholder
placeholder

func TestOpenAIRequestView_HasPatches(t *testing.T) {
	view := newOpenAIRequestView([]byte(`{"model":"gpt-5"placeholder`))
	require.False(t, view.HasPatches())

	view.MarkPatchSet("model", "gpt-5.1")
	require.True(t, view.HasPatches())

	view.DisablePatches()
	require.False(t, view.HasPatches())
placeholder

func TestOpenAIGatewayService_Forward_HTTPPatchPathKeepsLargeInputRaw(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body: io.NopCloser(strings.NewReader(
				`{"usage":{"input_tokens":1,"output_tokens":2,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholder`,
			)),
	placeholder,
placeholder
	cfg := &config.Config{placeholder
	cfg.Security.URLAllowlist.Enabled = false
	svc := &OpenAIGatewayService{cfg: cfg, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          1,
		Name:        "openai-apikey",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "https://example.com",
	placeholder,
		Extra: map[string]any{"use_responses_api": trueplaceholder,
placeholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	SetOpenAIClientTransport(c, OpenAIClientTransportHTTP)

	body := []byte(`{"model":"gpt-5","stream":false,"reasoning":{"effort":"minimal"placeholder,"input":[{"type":"message","content":[{"type":"input_text","text":"hi","nonce":9007199254740993placeholder]placeholder]placeholder`)
	result, err := svc.Forward(context.Background(), c, account, body)
placeholder
	require.NotNil(t, result)
	require.NotNil(t, upstream.lastReq)
	// 合成路径默认 instructions 现按模型填入真实 Codex base prompt（此处 inbound model=gpt-5）。
	encodedInstr, _ := json.Marshal(defaultCodexSynthInstructions("gpt-5"))
	expectedBody := fmt.Sprintf(`{"model":"gpt-5","stream":false,"reasoning":{"effort":"none"placeholder,"instructions":%s,"input":[{"type":"message","content":[{"type":"input_text","text":"hi","nonce":9007199254740993placeholder]placeholder]placeholder`, string(encodedInstr))
	require.JSONEq(t, expectedBody, string(upstream.lastBody))
	require.Equal(t, "9007199254740993", gjson.GetBytes(upstream.lastBody, "input.0.content.0.nonce").Raw)
placeholder

func TestOpenAIGatewayService_Forward_DecodedMutationKeepsLaterFieldDeletes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"usage":{"input_tokens":1,"output_tokens":2placeholderplaceholder`)),
	placeholder,
placeholder
	cfg := &config.Config{placeholder
	cfg.Security.URLAllowlist.Enabled = false
	svc := &OpenAIGatewayService{cfg: cfg, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          2,
		Name:        "openai-apikey",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "https://example.com",
	placeholder,
		Extra: map[string]any{"use_responses_api": trueplaceholder,
placeholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	SetOpenAIClientTransport(c, OpenAIClientTransportHTTP)

	body := []byte(`{"model":"gpt-5.4","stream":false,"max_completion_tokens":12,"tools":[{"type":"image_generation","format":"png"placeholder],"input":[{"type":"message","content":"draw"placeholder]placeholder`)
	result, err := svc.Forward(context.Background(), c, account, body)
placeholder
	require.NotNil(t, result)
	require.False(t, gjson.GetBytes(upstream.lastBody, "max_completion_tokens").Exists())
	require.False(t, gjson.GetBytes(upstream.lastBody, "tools.0.format").Exists())
	require.Equal(t, "png", gjson.GetBytes(upstream.lastBody, "tools.0.output_format").String())
placeholder

// #4417：/v1/responses 原生转发路径需将 Chat-Completions 风格的 max_tokens 归一化为
// max_output_tokens，并移除兼容上游不接受的 prompt_cache_options。
func TestOpenAIGatewayService_Forward_NormalizesMaxTokensAndStripsPromptCacheOptions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	runForward := func(t *testing.T, body []byte) []byte {
	placeholder
		upstream := &httpUpstreamRecorder{
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
				Body:       io.NopCloser(strings.NewReader(`{"usage":{"input_tokens":1,"output_tokens":2placeholderplaceholder`)),
		placeholder,
	placeholder
		cfg := &config.Config{placeholder
		cfg.Security.URLAllowlist.Enabled = false
		svc := &OpenAIGatewayService{cfg: cfg, httpUpstream: upstreamplaceholder
		account := &Account{
			ID:          4,
			Name:        "openai-apikey",
			Platform:    PlatformOpenAI,
			Type:        AccountTypeAPIKey,
			Concurrency: 1,
	placeholder
				"api_key":  "sk-test",
				"base_url": "https://example.com",
		placeholder,
			Extra: map[string]any{"openai_responses_supported": trueplaceholder,
	placeholder
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
		SetOpenAIClientTransport(c, OpenAIClientTransportHTTP)

		result, err := svc.Forward(context.Background(), c, account, body)
	placeholder
		require.NotNil(t, result)
		return upstream.lastBody
placeholder

	t.Run("max_tokens 归一化为 max_output_tokens 并移除 prompt_cache_options", func(t *testing.T) {
		out := runForward(t, []byte(`{"model":"gpt-5.4","stream":false,"max_tokens":256,"prompt_cache_options":{"enabled":trueplaceholder,"input":[{"type":"message","content":"hi"placeholder]placeholder`))
		require.Equal(t, int64(256), gjson.GetBytes(out, "max_output_tokens").Int())
		require.False(t, gjson.GetBytes(out, "max_tokens").Exists())
		require.False(t, gjson.GetBytes(out, "prompt_cache_options").Exists())
placeholder)

	t.Run("同时存在时保留 max_output_tokens 丢弃 max_tokens", func(t *testing.T) {
		out := runForward(t, []byte(`{"model":"gpt-5.4","stream":false,"max_tokens":256,"max_output_tokens":512,"input":[{"type":"message","content":"hi"placeholder]placeholder`))
		require.Equal(t, int64(512), gjson.GetBytes(out, "max_output_tokens").Int())
		require.False(t, gjson.GetBytes(out, "max_tokens").Exists())
placeholder)
placeholder

func TestOpenAIGatewayService_Forward_MappedImageModelUsesImageGate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"usage":{"input_tokens":1,"output_tokens":2placeholderplaceholder`)),
	placeholder,
placeholder
	cfg := &config.Config{placeholder
	cfg.Security.URLAllowlist.Enabled = false
	svc := &OpenAIGatewayService{cfg: cfg, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          3,
		Name:        "openai-apikey",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":       "sk-test",
			"base_url":      "https://example.com",
			"model_mapping": map[string]any{"draw-alias": "gpt-image-2"placeholder,
	placeholder,
		Extra: map[string]any{"use_responses_api": trueplaceholder,
placeholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	c.Set("api_key", &APIKey{Group: &Group{AllowImageGeneration: falseplaceholderplaceholder)
	SetOpenAIClientTransport(c, OpenAIClientTransportHTTP)

	body := []byte(`{"model":"draw-alias","stream":false,"input":"draw"placeholder`)
	result, err := svc.Forward(context.Background(), c, account, body)
placeholder
	require.Nil(t, result)
	require.Nil(t, upstream.lastReq)
	require.Equal(t, http.StatusForbidden, rec.Code)
	cached, known := getOpenAIImageIntentHint(c)
	require.True(t, known)
	require.False(t, cached)

	textAccount := *account
	textAccount.ID = 4
	textAccount.Credentials = map[string]any{
		"api_key":  "sk-test",
		"base_url": "https://example.com",
placeholder
	result, err = svc.Forward(context.Background(), c, &textAccount, body)
placeholder
	require.NotNil(t, result)
	require.NotNil(t, upstream.lastReq)
	require.Len(t, upstream.bodies, 1)
	cached, known = getOpenAIImageIntentHint(c)
	require.True(t, known)
	require.False(t, cached)
placeholder

func TestOpenAIGatewayService_Forward_TextResponsesSetsBillingModelToMappedModel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholder, "x-request-id": []string{"rid_text_mapped_billing"placeholderplaceholder,
			Body: io.NopCloser(strings.NewReader(
				`{"id":"resp_text_mapped","object":"response","model":"gpt-5.5","status":"completed","output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"ok"placeholder]placeholder],"usage":{"input_tokens":20,"output_tokens":10,"total_tokens":30placeholderplaceholder`,
			)),
	placeholder,
placeholder
	cfg := &config.Config{placeholder
	cfg.Security.URLAllowlist.Enabled = false
	svc := &OpenAIGatewayService{cfg: cfg, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          4,
		Name:        "openai-apikey",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":       "sk-test",
			"base_url":      "https://example.com",
			"model_mapping": map[string]any{"gpt-5.4": "gpt-5.5"placeholder,
	placeholder,
		Extra: map[string]any{"use_responses_api": trueplaceholder,
placeholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	SetOpenAIClientTransport(c, OpenAIClientTransportHTTP)

	body := []byte(`{"model":"gpt-5.4","stream":false,"input":"hello"placeholder`)
	result, err := svc.Forward(context.Background(), c, account, body)
placeholder
	require.NotNil(t, result)
	require.Equal(t, "gpt-5.4", result.Model)
	require.Equal(t, "gpt-5.5", result.BillingModel)
	require.Equal(t, "gpt-5.5", result.UpstreamModel)
	require.Equal(t, "gpt-5.5", gjson.GetBytes(upstream.lastBody, "model").String())
	require.Equal(t, 0, result.ImageCount)
placeholder

func TestOpenAIGatewayService_Forward_TextResponsesWithoutMappingKeepsRequestedBillingModel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholder, "x-request-id": []string{"rid_text_unmapped_billing"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"id":"resp_text_unmapped","object":"response","model":"gpt-5.4","status":"completed","usage":{"input_tokens":20,"output_tokens":10,"total_tokens":30placeholderplaceholder`)),
	placeholder,
placeholder
	cfg := &config.Config{placeholder
	cfg.Security.URLAllowlist.Enabled = false
	svc := &OpenAIGatewayService{cfg: cfg, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          4,
		Name:        "openai-apikey",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "https://example.com",
	placeholder,
		Extra: map[string]any{"use_responses_api": trueplaceholder,
placeholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	SetOpenAIClientTransport(c, OpenAIClientTransportHTTP)

	result, err := svc.Forward(context.Background(), c, account, []byte(`{"model":"gpt-5.4","stream":false,"input":"hello"placeholder`))
placeholder
	require.NotNil(t, result)
	require.Equal(t, "gpt-5.4", result.Model)
	require.Equal(t, "gpt-5.4", result.BillingModel)
	require.Equal(t, "gpt-5.4", result.UpstreamModel)
placeholder

func TestOpenAIGatewayService_Forward_TextResponsesBillingModelMatchesChatCompletions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{placeholder
	cfg.Security.URLAllowlist.Enabled = false
	account := &Account{
		ID:          5,
		Name:        "openai-apikey",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":       "sk-test",
			"base_url":      "https://example.com",
			"model_mapping": map[string]any{"gpt-5.4": "gpt-5.5"placeholder,
	placeholder,
		Extra: map[string]any{"use_responses_api": trueplaceholder,
placeholder

	responsesUpstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholder, "x-request-id": []string{"rid_responses_mapped_billing"placeholderplaceholder,
			Body: io.NopCloser(strings.NewReader(
				`{"id":"resp_native","object":"response","model":"gpt-5.5","status":"completed","output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"ok"placeholder]placeholder],"usage":{"input_tokens":20,"output_tokens":10,"total_tokens":30placeholderplaceholder`,
			)),
	placeholder,
placeholder
	responsesSvc := &OpenAIGatewayService{cfg: cfg, httpUpstream: responsesUpstreamplaceholder
	responsesRecorder := httptest.NewRecorder()
	responsesCtx, _ := gin.CreateTestContext(responsesRecorder)
	responsesCtx.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	SetOpenAIClientTransport(responsesCtx, OpenAIClientTransportHTTP)
	responsesResult, err := responsesSvc.Forward(context.Background(), responsesCtx, account, []byte(`{"model":"gpt-5.4","stream":false,"input":"hello"placeholder`))
placeholder
	require.NotNil(t, responsesResult)

	chatUpstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "x-request-id": []string{"rid_chat_mapped_billing"placeholderplaceholder,
			Body: io.NopCloser(strings.NewReader(
				`data: {"type":"response.completed","response":{"id":"resp_chat","object":"response","model":"gpt-5.5","status":"completed","output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"ok"placeholder]placeholder],"usage":{"input_tokens":20,"output_tokens":10,"total_tokens":30placeholderplaceholderplaceholder` + "\n\n",
			)),
	placeholder,
placeholder
	chatSvc := &OpenAIGatewayService{cfg: cfg, httpUpstream: chatUpstreamplaceholder
	chatRecorder := httptest.NewRecorder()
	chatCtx, _ := gin.CreateTestContext(chatRecorder)
	chatCtx.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/chat/completions", nil)
	chatResult, err := chatSvc.ForwardAsChatCompletions(context.Background(), chatCtx, account, []byte(`{"model":"gpt-5.4","stream":false,"messages":[{"role":"user","content":"hello"placeholder]placeholder`), "", "")
placeholder
	require.NotNil(t, chatResult)

	require.Equal(t, chatResult.BillingModel, responsesResult.BillingModel)
	require.Equal(t, "gpt-5.5", responsesResult.BillingModel)
	require.Equal(t, "gpt-5.5", chatResult.BillingModel)
placeholder

func TestOpenAIGatewayService_Forward_TextDataImageDoesNotForceMapMarshal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"usage":{"input_tokens":1,"output_tokens":2placeholderplaceholder`)),
	placeholder,
placeholder
	cfg := &config.Config{placeholder
	cfg.Security.URLAllowlist.Enabled = false
	svc := &OpenAIGatewayService{cfg: cfg, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          4,
		Name:        "openai-apikey",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "https://example.com",
	placeholder,
		Extra: map[string]any{"use_responses_api": trueplaceholder,
placeholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	SetOpenAIClientTransport(c, OpenAIClientTransportHTTP)

	body := []byte(`{"model":"gpt-5","stream":false,"input":[{"type":"message","content":[{"type":"input_text","text":"literal data:image/png;base64, only","nonce":1e1000000placeholder]placeholder]placeholder`)
	result, err := svc.Forward(context.Background(), c, account, body)
placeholder
	require.NotNil(t, result)
	require.Equal(t, "1e1000000", gjson.GetBytes(upstream.lastBody, "input.0.content.0.nonce").Raw)
placeholder

func TestOpenAIGatewayService_Forward_ImageToolBillingDoesNotForceFullDecode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body: io.NopCloser(strings.NewReader(
				`{"output":[{"id":"ig_1","type":"image_generation_call","result":"final-image"placeholder],"usage":{"input_tokens":1,"output_tokens":2placeholderplaceholder`,
			)),
	placeholder,
placeholder
	cfg := &config.Config{placeholder
	cfg.Security.URLAllowlist.Enabled = false
	svc := &OpenAIGatewayService{cfg: cfg, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          9,
		Name:        "openai-apikey",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "https://example.com",
	placeholder,
		Extra: map[string]any{"use_responses_api": trueplaceholder,
placeholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	SetOpenAIClientTransport(c, OpenAIClientTransportHTTP)

	body := []byte(`{"model":"gpt-5","stream":false,"tools":[{"type":"image_generation","model":"gpt-image-2","size":"2048x1152"placeholder],"input":[{"type":"message","content":[{"type":"input_text","text":"draw","nonce":1e1000000placeholder]placeholder]placeholder`)
	result, err := svc.Forward(context.Background(), c, account, body)
placeholder
	require.NotNil(t, result)
	require.Equal(t, "1e1000000", gjson.GetBytes(upstream.lastBody, "input.0.content.0.nonce").Raw)
	require.Equal(t, 1, result.ImageCount)
	require.Equal(t, "2K", result.ImageSize)
	require.Equal(t, "gpt-image-2", result.BillingModel)
placeholder

func TestOpenAIGatewayService_Forward_ImageToolWithImageOnlyModelIsNormalized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"usage":{"input_tokens":1,"output_tokens":2placeholderplaceholder`)),
	placeholder,
placeholder
	cfg := &config.Config{placeholder
	cfg.Security.URLAllowlist.Enabled = false
	svc := &OpenAIGatewayService{cfg: cfg, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          11,
		Name:        "openai-apikey",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "https://example.com",
	placeholder,
		Extra: map[string]any{"use_responses_api": trueplaceholder,
placeholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	SetOpenAIClientTransport(c, OpenAIClientTransportHTTP)

	body := []byte(`{"model":"gpt-image-2","stream":false,"tools":[{"type":"image_generation","model":"gpt-image-2"placeholder],"input":"draw"placeholder`)
	result, err := svc.Forward(context.Background(), c, account, body)
placeholder
	require.NotNil(t, result)
	require.Equal(t, openAIImagesResponsesMainModel, gjson.GetBytes(upstream.lastBody, "model").String())
placeholder

func TestOpenAIGatewayService_Forward_HTTPRetryRecoveryDoesNotDecodeBeforeError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstream := &httpUpstreamRecorder{
		responses: []*http.Response{
			{
				StatusCode: http.StatusBadRequest,
				Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
				Body:       io.NopCloser(strings.NewReader(`{"error":{"code":"invalid_encrypted_content","type":"invalid_request_error","message":"bad encrypted content"placeholderplaceholder`)),
		placeholder,
			{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
				Body:       io.NopCloser(strings.NewReader(`{"usage":{"input_tokens":1,"output_tokens":2placeholderplaceholder`)),
		placeholder,
	placeholder,
placeholder
	cfg := &config.Config{placeholder
	cfg.Security.URLAllowlist.Enabled = false
	svc := &OpenAIGatewayService{cfg: cfg, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          10,
		Name:        "openai-apikey",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "https://example.com",
	placeholder,
		Extra: map[string]any{"use_responses_api": trueplaceholder,
placeholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	SetOpenAIClientTransport(c, OpenAIClientTransportHTTP)

	body := []byte(`{"model":"gpt-5","stream":false,"input":[{"type":"reasoning","encrypted_content":"gAAA","summary":[{"type":"summary_text","text":"keep me"placeholder]placeholder,{"type":"message","content":[{"type":"input_text","text":"hi","nonce":9007199254740993placeholder]placeholder]placeholder`)
	result, err := svc.Forward(context.Background(), c, account, body)
placeholder
	require.NotNil(t, result)
	require.Len(t, upstream.bodies, 2)
	require.Equal(t, "gAAA", gjson.GetBytes(upstream.bodies[0], "input.0.encrypted_content").String())
	require.Equal(t, "9007199254740993", gjson.GetBytes(upstream.bodies[0], "input.1.content.0.nonce").Raw)
	require.False(t, gjson.GetBytes(upstream.bodies[1], "input.0.encrypted_content").Exists())
	require.Equal(t, "summary_text", gjson.GetBytes(upstream.bodies[1], "input.0.summary.0.type").String())
placeholder

func TestOpenAIGatewayService_Forward_CodexSparkRejectsEscapedInputImage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"usage":{"input_tokens":1,"output_tokens":2placeholderplaceholder`)),
	placeholder,
placeholder
	cfg := &config.Config{placeholder
	cfg.Security.URLAllowlist.Enabled = false
	svc := &OpenAIGatewayService{cfg: cfg, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          5,
		Name:        "openai-apikey",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "https://example.com",
	placeholder,
		Extra: map[string]any{"use_responses_api": trueplaceholder,
placeholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	SetOpenAIClientTransport(c, OpenAIClientTransportHTTP)

	body := []byte(`{"model":"gpt-5.3-codex-spark","stream":false,"input":[{"type":"input_` + "\\u0069" + `mage","file_id":"file_1"placeholder]placeholder`)
	result, err := svc.Forward(context.Background(), c, account, body)
placeholder
	require.Nil(t, result)
	require.Nil(t, upstream.lastReq)
	require.Equal(t, http.StatusBadRequest, rec.Code)
placeholder

func TestOpenAIGatewayService_Forward_CodexBridgeInjectionSetsImageBilling(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body: io.NopCloser(strings.NewReader(
				`{"output":[{"id":"ig_1","type":"image_generation_call","result":"final-image","size":"1024x1024"placeholder],"usage":{"input_tokens":1,"output_tokens":2placeholderplaceholder`,
			)),
	placeholder,
placeholder
	cfg := &config.Config{placeholder
	cfg.Security.URLAllowlist.Enabled = false
	cfg.Gateway.ForceCodexCLI = true
	cfg.Gateway.CodexImageGenerationBridgeEnabled = true
	svc := &OpenAIGatewayService{cfg: cfg, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          7,
		Name:        "openai-apikey",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "https://example.com",
	placeholder,
		Extra: map[string]any{"use_responses_api": trueplaceholder,
placeholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	c.Set("api_key", &APIKey{Group: &Group{AllowImageGeneration: trueplaceholderplaceholder)
	SetOpenAIClientTransport(c, OpenAIClientTransportHTTP)

	body := []byte(`{"model":"gpt-5","stream":false,"input":"draw if needed"placeholder`)
	result, err := svc.Forward(context.Background(), c, account, body)
placeholder
	require.NotNil(t, result)
	require.Equal(t, 1, result.ImageCount)
	require.Equal(t, "2K", result.ImageSize)
	require.Equal(t, "gpt-image-2", result.BillingModel)
placeholder

func TestOpenAIGatewayService_Forward_HTTPDeletesPreviousResponseIDWhenPresent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{placeholder
	cfg.Security.URLAllowlist.Enabled = false
	account := &Account{
		ID:          8,
		Name:        "openai-apikey",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "https://example.com",
	placeholder,
		Extra: map[string]any{"use_responses_api": trueplaceholder,
placeholder

	for _, body := range [][]byte{
		[]byte(`{"model":"gpt-5","stream":false,"previous_response_id":"","input":"hi"placeholder`),
		[]byte(`{"model":"gpt-5","stream":false,"previous_response_id":null,"input":"hi"placeholder`),
placeholder {
		upstream := &httpUpstreamRecorder{
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
				Body:       io.NopCloser(strings.NewReader(`{"usage":{"input_tokens":1,"output_tokens":2placeholderplaceholder`)),
		placeholder,
	placeholder
		svc := &OpenAIGatewayService{cfg: cfg, httpUpstream: upstreamplaceholder
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
		SetOpenAIClientTransport(c, OpenAIClientTransportHTTP)

		result, err := svc.Forward(context.Background(), c, account, body)
	placeholder
		require.NotNil(t, result)
		require.False(t, gjson.GetBytes(upstream.lastBody, "previous_response_id").Exists())
placeholder
placeholder

func TestOpenAIGatewayService_Forward_StripsImageGenerationToolForSparkAPIKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"usage":{"input_tokens":1,"output_tokens":2placeholderplaceholder`)),
	placeholder,
placeholder
	cfg := &config.Config{placeholder
	cfg.Security.URLAllowlist.Enabled = false
	svc := &OpenAIGatewayService{cfg: cfg, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          11,
		Name:        "openai-apikey",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "https://example.com",
	placeholder,
		Extra: map[string]any{"use_responses_api": trueplaceholder,
placeholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	// Allow image generation so the tool is normalized (not gated out), reproducing
	// the leak the strip must override.
	c.Set("api_key", &APIKey{Group: &Group{AllowImageGeneration: trueplaceholderplaceholder)
	SetOpenAIClientTransport(c, OpenAIClientTransportHTTP)

	body := []byte(`{"model":"gpt-5.3-codex-spark","stream":false,"input":"hi","tools":[{"type":"function","name":"shell"placeholder,{"type":"image_generation","output_format":"png"placeholder]placeholder`)
	result, err := svc.Forward(context.Background(), c, account, body)
placeholder
	require.NotNil(t, result)
	require.NotNil(t, upstream.lastReq)
	require.False(t, gjson.GetBytes(upstream.lastBody, `tools.#(type=="image_generation")`).Exists())
	require.True(t, gjson.GetBytes(upstream.lastBody, `tools.#(type=="function")`).Exists())
placeholder

func TestOpenAIRequestBodyMayContainEmptyBase64InputImageSeesEscapedJSON(t *testing.T) {
	body := []byte(`{"input":[{"type":"message","content":[{"type":"input_image","image_` + "\\u0075" + `rl":"data:image/png;base64` + "\\u002c" + `   "placeholder]placeholder]placeholder`)

	require.True(t, openAIRequestBodyMayContainEmptyBase64InputImage(body))
placeholder

func TestOpenAIRequestBodyMayContainEmptyBase64InputImageSeesEscapedImageType(t *testing.T) {
	body := []byte(`{"input":[{"type":"message","content":[{"type":"input_` + "\\u0069" + `mage","image_url":"data:image/png;base64,   "placeholder]placeholder]placeholder`)

	require.True(t, openAIRequestBodyMayContainEmptyBase64InputImage(body))
placeholder

func TestOpenAIRequestBodyMayContainEmptyBase64InputImageSeesEscapedInputPrefix(t *testing.T) {
	body := []byte(`{"input":[{"type":"message","content":[{"type":"inp` + "\\u0075" + `t_image","image_url":"data:image/png;base64,   "placeholder]placeholder]placeholder`)

	require.True(t, openAIRequestBodyMayContainEmptyBase64InputImage(body))
placeholder

func TestOpenAIGatewayService_Forward_ImageOnlyModelKeepsSupportedVerbosity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"usage":{"input_tokens":1,"output_tokens":2placeholderplaceholder`)),
	placeholder,
placeholder
	cfg := &config.Config{placeholder
	cfg.Security.URLAllowlist.Enabled = false
	svc := &OpenAIGatewayService{cfg: cfg, httpUpstream: upstreamplaceholder
	account := &Account{
		ID:          6,
		Name:        "openai-apikey",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "sk-test",
			"base_url": "https://example.com",
	placeholder,
		Extra: map[string]any{"use_responses_api": trueplaceholder,
placeholder
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	SetOpenAIClientTransport(c, OpenAIClientTransportHTTP)

	body := []byte(`{"model":"gpt-image-2","stream":false,"text":{"verbosity":"low"placeholder,"input":"draw"placeholder`)
	result, err := svc.Forward(context.Background(), c, account, body)
placeholder
	require.NotNil(t, result)
	require.Equal(t, "low", gjson.GetBytes(upstream.lastBody, "text.verbosity").String())
	require.Equal(t, openAIImagesResponsesMainModel, gjson.GetBytes(upstream.lastBody, "model").String())
placeholder

func TestExtractOpenAIRequestMetaFromBody(t *testing.T) {
	tests := []struct {
		name          string
		body          []byte
		wantModel     string
		wantStream    bool
		wantPromptKey string
placeholder{
		{
			name:          "完整字段",
			body:          []byte(`{"model":"gpt-5","stream":true,"prompt_cache_key":" ses-1 "placeholder`),
			wantModel:     "gpt-5",
			wantStream:    true,
			wantPromptKey: "ses-1",
	placeholder,
		{
			name:          "缺失可选字段",
			body:          []byte(`{"model":"gpt-4"placeholder`),
			wantModel:     "gpt-4",
			wantStream:    false,
			wantPromptKey: "",
	placeholder,
		{
			name:          "空请求体",
			body:          nil,
			wantModel:     "",
			wantStream:    false,
			wantPromptKey: "",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, stream, promptKey := extractOpenAIRequestMetaFromBody(tt.body)
			require.Equal(t, tt.wantModel, model)
			require.Equal(t, tt.wantStream, stream)
			require.Equal(t, tt.wantPromptKey, promptKey)
	placeholder)
placeholder
placeholder

func TestExtractOpenAIReasoningEffortFromBody(t *testing.T) {
	tests := []struct {
		name      string
		body      []byte
		model     string
		wantNil   bool
		wantValue string
placeholder{
		{
			name:      "优先读取 reasoning.effort",
			body:      []byte(`{"reasoning":{"effort":"medium"placeholderplaceholder`),
			model:     "gpt-5-high",
			wantNil:   false,
			wantValue: "medium",
	placeholder,
		{
			name:      "兼容 reasoning_effort",
			body:      []byte(`{"reasoning_effort":"x-high"placeholder`),
			model:     "",
			wantNil:   false,
			wantValue: "xhigh",
	placeholder,
		{
			name:      "DeepSeek max 归一化为 xhigh",
			body:      []byte(`{"reasoning_effort":"max"placeholder`),
			model:     "deepseek-v4-pro",
			wantNil:   false,
			wantValue: "xhigh",
	placeholder,
		{
			name:    "minimal 归一化为空",
			body:    []byte(`{"reasoning":{"effort":"minimal"placeholderplaceholder`),
			model:   "gpt-5-high",
			wantNil: true,
	placeholder,
		{
			name:      "缺失字段时从模型后缀推导",
			body:      []byte(`{"input":"hi"placeholder`),
			model:     "gpt-5-high",
			wantNil:   false,
			wantValue: "high",
	placeholder,
		{
			name:    "未知后缀不返回",
			body:    []byte(`{"input":"hi"placeholder`),
			model:   "gpt-5-unknown",
			wantNil: true,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractOpenAIReasoningEffortFromBody(tt.body, tt.model)
			if tt.wantNil {
				require.Nil(t, got)
				return
		placeholder
			require.NotNil(t, got)
			require.Equal(t, tt.wantValue, *got)
	placeholder)
placeholder
placeholder

func TestGetOpenAIRequestBodyMap_ParseError(t *testing.T) {
	_, err := getOpenAIRequestBodyMap(nil, []byte(`{invalid-json`))
placeholder
	require.Contains(t, err.Error(), "parse request")
placeholder

func TestGetOpenAIRequestBodyMap_DoesNotWriteContextCache(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	got, err := getOpenAIRequestBodyMap(c, []byte(`{"model":"gpt-5","stream":trueplaceholder`))
placeholder
	require.Equal(t, "gpt-5", got["model"])
	require.Empty(t, c.Keys)
placeholder

func TestSanitizeEmptyBase64InputImagesInOpenAIRequestBodyMap(t *testing.T) {
	var reqBody map[string]any
	require.NoError(t, json.Unmarshal([]byte(`{
		"model":"gpt-5.4",
		"input":[
			{"role":"user","content":[
				{"type":"input_text","text":"Describe this"placeholder,
				{"type":"input_image","image_url":"data:image/png;base64,   "placeholder,
				{"type":"input_image","image_url":"data:image/png;base64,abc123"placeholder
			]placeholder,
			{"role":"user","content":[
				{"type":"input_image","image_url":"data:image/png;base64,"placeholder
			]placeholder,
			{"type":"input_image","image_url":"data:image/png;base64,"placeholder,
			{"type":"input_image","image_url":"data:image/png;base64,top-level-valid"placeholder
		]
placeholder`), &reqBody))

	require.True(t, sanitizeEmptyBase64InputImagesInOpenAIRequestBodyMap(reqBody))

	normalized, err := json.Marshal(reqBody)
placeholder
	require.JSONEq(t, `{
		"model":"gpt-5.4",
		"input":[
			{"role":"user","content":[
				{"type":"input_text","text":"Describe this"placeholder,
				{"type":"input_image","image_url":"data:image/png;base64,abc123"placeholder
			]placeholder,
			{"type":"input_image","image_url":"data:image/png;base64,top-level-valid"placeholder
		]
placeholder`, string(normalized))
placeholder

func TestSanitizeEmptyBase64InputImagesInOpenAIBody(t *testing.T) {
	body, changed, err := sanitizeEmptyBase64InputImagesInOpenAIBody([]byte(`{
		"model":"gpt-5.4",
		"stream":true,
		"input":[
			{"role":"user","content":[
				{"type":"input_text","text":"Describe this"placeholder,
				{"type":"input_image","image_url":"data:image/png;base64,"placeholder
			]placeholder
		]
placeholder`))
placeholder
	require.True(t, changed)
	require.JSONEq(t, `{
		"model":"gpt-5.4",
		"stream":true,
		"input":[
			{"role":"user","content":[
				{"type":"input_text","text":"Describe this"placeholder
			]placeholder
		]
placeholder`, string(body))
placeholder
