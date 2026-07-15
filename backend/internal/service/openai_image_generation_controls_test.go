package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestOpenAIGatewayServiceForward_RejectsDisabledImageGenerationIntents(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string
		body []byte
placeholder{
		{
			name: "image model",
			body: []byte(`{"model":"gpt-image-2","input":"draw"placeholder`),
	placeholder,
		{
			name: "image tool",
			body: []byte(`{"model":"gpt-5.4","input":"draw","tools":[{"type":"image_generation"placeholder]placeholder`),
	placeholder,
		{
			name: "image tool choice",
			body: []byte(`{"model":"gpt-5.4","input":"draw","tool_choice":{"type":"image_generation"placeholderplaceholder`),
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upstream := &httpUpstreamRecorder{placeholder
			svc := newOpenAIImageGenerationControlTestService(upstream)
			c, recorder := newOpenAIImageGenerationControlTestContext(false, "unit-test-agent/1.0")
			account := newOpenAIImageGenerationControlTestAccount()

			result, err := svc.Forward(context.Background(), c, account, tt.body)

		placeholder
			require.Nil(t, result)
			require.Equal(t, http.StatusForbidden, recorder.Code)
			require.Equal(t, "permission_error", gjson.GetBytes(recorder.Body.Bytes(), "error.type").String())
			require.Nil(t, upstream.lastReq, "disabled image request must not reach upstream")
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayServiceForward_DisabledGroupAllowsTextOnlyResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"id":"resp_text","model":"gpt-5.4","usage":{"input_tokens":3,"output_tokens":2placeholderplaceholder`)),
	placeholder,
placeholder
	svc := newOpenAIImageGenerationControlTestService(upstream)
	c, recorder := newOpenAIImageGenerationControlTestContext(false, "unit-test-agent/1.0")
	account := newOpenAIImageGenerationControlTestAccount()

	result, err := svc.Forward(context.Background(), c, account, []byte(`{"model":"gpt-5.4","input":"write code","stream":falseplaceholder`))

placeholder
	require.NotNil(t, result)
	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, 3, result.Usage.InputTokens)
	require.Equal(t, 2, result.Usage.OutputTokens)
	require.Equal(t, 0, result.ImageCount)
	require.NotNil(t, upstream.lastReq)
placeholder

func TestOpenAIGatewayServiceForward_CodexImageInjectionRespectsGroupCapability(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		allowImages   bool
		bridgeEnabled bool
		responsesLite bool
		wantInjected  bool
placeholder{
		{name: "disabled group skips injection", allowImages: false, bridgeEnabled: true, wantInjected: falseplaceholder,
		{name: "enabled group skips injection by default", allowImages: true, bridgeEnabled: false, wantInjected: falseplaceholder,
		{name: "enabled group injects image tool when bridge enabled", allowImages: true, bridgeEnabled: true, wantInjected: trueplaceholder,
		{name: "responses lite skips hosted image bridge", allowImages: true, bridgeEnabled: true, responsesLite: true, wantInjected: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upstream := &httpUpstreamRecorder{
				resp: &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
					Body:       io.NopCloser(strings.NewReader(`{"id":"resp_codex","model":"gpt-5.4","usage":{"input_tokens":1,"output_tokens":1placeholderplaceholder`)),
			placeholder,
		placeholder
			svc := newOpenAIImageGenerationControlTestService(upstream)
			svc.cfg.Gateway.CodexImageGenerationBridgeEnabled = tt.bridgeEnabled
			c, _ := newOpenAIImageGenerationControlTestContext(tt.allowImages, "codex_cli_rs/0.98.0")
			if tt.responsesLite {
				c.Request.Header.Set(responsesLiteHeader, "true")
		placeholder
			account := newOpenAIImageGenerationControlTestAccount()

			result, err := svc.Forward(context.Background(), c, account, []byte(`{"model":"gpt-5.4","input":"write code","stream":falseplaceholder`))

		placeholder
			require.NotNil(t, result)
			require.NotNil(t, upstream.lastReq)
			hasImageTool := gjson.GetBytes(upstream.lastBody, `tools.#(type=="image_generation")`).Exists()
			require.Equal(t, tt.wantInjected, hasImageTool)
			expectedLiteHeader := ""
			if tt.responsesLite {
				expectedLiteHeader = "true"
		placeholder
			require.Equal(t, expectedLiteHeader, upstream.lastReq.Header.Get(responsesLiteHeader))
			instructions := gjson.GetBytes(upstream.lastBody, "instructions").String()
			require.Equal(t, tt.wantInjected, strings.Contains(instructions, "image_generation"))
			toolChoice := gjson.GetBytes(upstream.lastBody, "tool_choice")
			require.Equal(t, tt.wantInjected, toolChoice.Exists())
			if tt.wantInjected {
				require.Equal(t, "auto", toolChoice.String())
		placeholder
	placeholder)
placeholder
placeholder

func TestOpenAIBuildUpstreamRequestOpenAIPassthroughForwardsResponsesLiteHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := newOpenAIImageGenerationControlTestContext(true, "codex_cli_rs/0.98.0")
	c.Request.Header.Set(responsesLiteHeader, "true")

	svc := newOpenAIImageGenerationControlTestService(&httpUpstreamRecorder{placeholder)
	req, err := svc.buildUpstreamRequestOpenAIPassthrough(
		c.Request.Context(),
		c,
		newOpenAIImageGenerationControlTestAccount(),
		[]byte(`{"model":"gpt-5.4","input":"write code"placeholder`),
		"test-token",
	)

placeholder
	require.Equal(t, "true", req.Header.Get(responsesLiteHeader))
placeholder

func TestOpenAIGatewayServiceForward_ExplicitImageToolWorksWithBridgeDisabled(t *testing.T) {
	gin.SetMode(gin.TestMode)

	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"id":"resp_explicit_image","model":"gpt-5.4","usage":{"input_tokens":2,"output_tokens":1placeholderplaceholder`)),
	placeholder,
placeholder
	svc := newOpenAIImageGenerationControlTestService(upstream)
	c, _ := newOpenAIImageGenerationControlTestContext(true, "codex_cli_rs/0.98.0")
	account := newOpenAIImageGenerationControlTestAccount()
	body := []byte(`{"model":"gpt-5.4","input":"draw","stream":false,"tools":[{"type":"image_generation","format":"jpeg"placeholder]placeholder`)

	result, err := svc.Forward(context.Background(), c, account, body)

placeholder
	require.NotNil(t, result)
	require.NotNil(t, upstream.lastReq)
	require.True(t, gjson.GetBytes(upstream.lastBody, `tools.#(type=="image_generation")`).Exists())
	require.Equal(t, "jpeg", gjson.GetBytes(upstream.lastBody, `tools.#(type=="image_generation").output_format`).String())
	require.False(t, gjson.GetBytes(upstream.lastBody, `tools.#(type=="image_generation").format`).Exists())
	instructions := gjson.GetBytes(upstream.lastBody, "instructions").String()
	require.NotContains(t, instructions, "image_generation")
placeholder

func TestOpenAIGatewayServiceForward_AccountPolicyStripsExplicitImageTool(t *testing.T) {
	gin.SetMode(gin.TestMode)

	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"id":"resp_stripped_image","model":"gpt-5.4","usage":{"input_tokens":2,"output_tokens":1placeholderplaceholder`)),
	placeholder,
placeholder
	svc := newOpenAIImageGenerationControlTestService(upstream)
	c, _ := newOpenAIImageGenerationControlTestContext(true, "codex_cli_rs/0.98.0")
	account := newOpenAIImageGenerationControlTestAccount()
	account.Extra = map[string]any{
		featureKeyCodexImageGenerationExplicitToolPolicy: codexImageGenerationExplicitToolPolicyStrip,
placeholder
	body := []byte(`{
		"model":"gpt-5.4",
		"input":"draw",
		"stream":false,
		"tools":[
			{"type":"function","name":"shell","parameters":{"type":"object"placeholderplaceholder,
			{"type":"image_generation","format":"jpeg"placeholder
		],
		"tool_choice":{"type":"image_generation"placeholder
placeholder`)

	result, err := svc.Forward(context.Background(), c, account, body)

placeholder
	require.NotNil(t, result)
	require.NotNil(t, upstream.lastReq)
	require.False(t, gjson.GetBytes(upstream.lastBody, `tools.#(type=="image_generation")`).Exists())
	require.True(t, gjson.GetBytes(upstream.lastBody, `tools.#(type=="function")`).Exists())
	require.False(t, gjson.GetBytes(upstream.lastBody, "tool_choice").Exists())
	instructions := gjson.GetBytes(upstream.lastBody, "instructions").String()
	require.NotContains(t, instructions, "image_generation")
placeholder

func TestOpenAIGatewayServiceForward_AccountPolicyStripsImageNamespaceTools(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		passthrough bool
placeholder{
		{name: "managed forwarding"placeholder,
		{name: "passthrough forwarding", passthrough: trueplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upstream := &httpUpstreamRecorder{
				resp: &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
					Body:       io.NopCloser(strings.NewReader(`{"id":"resp_stripped_namespace","model":"gpt-5.5","usage":{"input_tokens":2,"output_tokens":1placeholderplaceholder`)),
			placeholder,
		placeholder
			svc := newOpenAIImageGenerationControlTestService(upstream)
			c, _ := newOpenAIImageGenerationControlTestContext(false, "codex_cli_rs/0.144.1")
			SetOpenAIClientTransport(c, OpenAIClientTransportHTTP)
			account := newOpenAIImageGenerationControlTestAccount()
			account.Extra = map[string]any{
				featureKeyCodexImageGenerationExplicitToolPolicy: codexImageGenerationExplicitToolPolicyStrip,
				"openai_passthrough":                             tt.passthrough,
		placeholder
			body := []byte(`{
				"model":"gpt-5.5",
				"stream":false,
				"tools":[
					{"type":"function","name":"shell","parameters":{"type":"object"placeholderplaceholder,
					{"type":"namespace","name":"image_gen","tools":[{"type":"function","name":"imagegen"placeholder]placeholder,
					{"type":"namespace","name":"code_tools","tools":[{"type":"function","name":"run"placeholder]placeholder
				],
				"input":[
					{"type":"message","role":"user","content":[{"type":"input_text","text":"write code"placeholder]placeholder,
					{"type":"additional_tools","tools":[{"type":"namespace","name":"image_gen","tools":[{"type":"function","name":"imagegen"placeholder]placeholder]placeholder
				],
				"tool_choice":"auto"
		placeholder`)

			result, err := svc.Forward(context.Background(), c, account, body)

		placeholder
			require.NotNil(t, result)
			require.NotNil(t, upstream.lastReq)
			var forwarded map[string]any
			require.NoError(t, json.Unmarshal(upstream.lastBody, &forwarded))
			require.False(t, hasOpenAIImageGenerationTool(forwarded))
			require.Equal(t, "auto", forwarded["tool_choice"])
			require.True(t, gjson.GetBytes(upstream.lastBody, `tools.#(name=="shell")`).Exists())
			require.True(t, gjson.GetBytes(upstream.lastBody, `tools.#(name=="code_tools")`).Exists())
			require.Equal(t, "write code", gjson.GetBytes(upstream.lastBody, "input.0.content.0.text").String())
			cached, known := getOpenAIImageIntentHint(c)
			require.True(t, known)
			require.True(t, cached)
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayServiceForward_ChannelBridgeOverrideEnablesCodexInjection(t *testing.T) {
	gin.SetMode(gin.TestMode)

	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"id":"resp_channel_bridge","model":"gpt-5.4","usage":{"input_tokens":1,"output_tokens":1placeholderplaceholder`)),
	placeholder,
placeholder
	svc := newOpenAIImageGenerationControlTestService(upstream)
	groupID := int64(4242)
	svc.channelService = newOpenAIImageGenerationControlChannelService(groupID, &Channel{
		ID:     9001,
		Status: StatusActive,
		FeaturesConfig: map[string]any{
			featureKeyCodexImageGenerationBridge: map[string]any{PlatformOpenAI: trueplaceholder,
	placeholder,
placeholder)
	c, _ := newOpenAIImageGenerationControlTestContext(true, "codex_cli_rs/0.98.0")
	account := newOpenAIImageGenerationControlTestAccount()

	result, err := svc.Forward(context.Background(), c, account, []byte(`{"model":"gpt-5.4","input":"write code","stream":falseplaceholder`))

placeholder
	require.NotNil(t, result)
	require.NotNil(t, upstream.lastReq)
	require.True(t, gjson.GetBytes(upstream.lastBody, `tools.#(type=="image_generation")`).Exists())
	require.Equal(t, "auto", gjson.GetBytes(upstream.lastBody, "tool_choice").String())
	instructions := gjson.GetBytes(upstream.lastBody, "instructions").String()
	require.Contains(t, instructions, "image_generation")
placeholder

func TestOpenAIGatewayServiceForward_CodexBridgeDoesNotInjectHostedToolAlongsideImageGenNamespace(t *testing.T) {
	gin.SetMode(gin.TestMode)

	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"id":"resp_namespace_image","model":"gpt-5.5","usage":{"input_tokens":1,"output_tokens":1placeholderplaceholder`)),
	placeholder,
placeholder
	svc := newOpenAIImageGenerationControlTestService(upstream)
	svc.cfg.Gateway.CodexImageGenerationBridgeEnabled = true
	c, _ := newOpenAIImageGenerationControlTestContext(true, "codex_cli_rs/0.144.1")
	account := newOpenAIImageGenerationControlTestAccount()
	body := []byte(`{
		"model":"gpt-5.5",
		"stream":false,
		"tools":[
			{"type":"function","name":"shell","parameters":{"type":"object"placeholderplaceholder,
			{"type":"namespace","name":"image_gen","tools":[{"type":"function","name":"imagegen"placeholder]placeholder
		],
		"input":[
			{"type":"message","role":"user","content":[{"type":"input_text","text":"draw a cat"placeholder]placeholder,
			{"type":"additional_tools","tools":[{"type":"namespace","name":"image_gen","tools":[{"type":"function","name":"imagegen"placeholder]placeholder]placeholder
		],
		"tool_choice":"auto"
placeholder`)

	result, err := svc.Forward(context.Background(), c, account, body)

placeholder
	require.NotNil(t, result)
	require.NotNil(t, upstream.lastReq)
	require.False(t, gjson.GetBytes(upstream.lastBody, `tools.#(type=="image_generation")`).Exists())
	require.Equal(t, "namespace", gjson.GetBytes(upstream.lastBody, `tools.#(name=="image_gen").type`).String())
	require.Equal(t, "namespace", gjson.GetBytes(upstream.lastBody, `input.#(type=="additional_tools").tools.#(name=="image_gen").type`).String())
placeholder

func TestOpenAIGatewayServiceForward_CodexBridgePreservesImageGenFunction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string
		tool string
placeholder{
		{
			name: "flat function",
			tool: `{"type":"function","name":"image_gen.imagegen","parameters":{"type":"object"placeholderplaceholder`,
	placeholder,
		{
			name: "nested function",
			tool: `{"type":"function","function":{"name":"image_gen.imagegen","parameters":{"type":"object"placeholderplaceholderplaceholder`,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upstream := &httpUpstreamRecorder{
				resp: &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
					Body:       io.NopCloser(strings.NewReader(`{"id":"resp_function_image","model":"gpt-5.5","usage":{"input_tokens":1,"output_tokens":1placeholderplaceholder`)),
			placeholder,
		placeholder
			svc := newOpenAIImageGenerationControlTestService(upstream)
			svc.cfg.Gateway.CodexImageGenerationBridgeEnabled = true
			c, _ := newOpenAIImageGenerationControlTestContext(true, "codex_cli_rs/0.144.1")
			account := newOpenAIImageGenerationControlTestAccount()
			body := []byte(`{"model":"gpt-5.5","input":"draw a cat","stream":false,"tools":[` + tt.tool + `]placeholder`)

			result, err := svc.Forward(context.Background(), c, account, body)

		placeholder
			require.NotNil(t, result)
			require.NotNil(t, upstream.lastReq)

			var forwarded map[string]any
			require.NoError(t, json.Unmarshal(upstream.lastBody, &forwarded))
			require.True(t, hasCodexImageGenerationFunctionTool(forwarded))
			require.False(t, gjson.GetBytes(upstream.lastBody, `tools.#(type=="image_generation")`).Exists())
			require.False(t, gjson.GetBytes(upstream.lastBody, "tool_choice").Exists())
			require.NotContains(t, gjson.GetBytes(upstream.lastBody, "instructions").String(), codexImageGenerationBridgeMarker)
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayServiceForward_CodexBridgePreservesExistingToolChoice(t *testing.T) {
	gin.SetMode(gin.TestMode)

	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"id":"resp_codex_tool_choice","model":"gpt-5.4","usage":{"input_tokens":1,"output_tokens":1placeholderplaceholder`)),
	placeholder,
placeholder
	svc := newOpenAIImageGenerationControlTestService(upstream)
	svc.cfg.Gateway.CodexImageGenerationBridgeEnabled = true
	c, _ := newOpenAIImageGenerationControlTestContext(true, "codex_cli_rs/0.98.0")
	account := newOpenAIImageGenerationControlTestAccount()

	result, err := svc.Forward(context.Background(), c, account, []byte(`{"model":"gpt-5.4","input":"draw","stream":false,"tools":[{"type":"image_generation"placeholder],"tool_choice":{"type":"image_generation"placeholderplaceholder`))

placeholder
	require.NotNil(t, result)
	require.NotNil(t, upstream.lastReq)
	require.Equal(t, "image_generation", gjson.GetBytes(upstream.lastBody, "tool_choice.type").String())
placeholder

func TestOpenAIGatewayServiceForward_CodexBridgeSkipsCompactRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	upstream := &httpUpstreamRecorder{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"id":"resp_codex_compact","model":"gpt-5.4","usage":{"input_tokens":1,"output_tokens":1placeholderplaceholder`)),
	placeholder,
placeholder
	svc := newOpenAIImageGenerationControlTestService(upstream)
	svc.cfg.Gateway.CodexImageGenerationBridgeEnabled = true
	c, _ := newOpenAIImageGenerationControlTestContext(true, "codex_cli_rs/0.98.0")
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses/compact", nil)
	c.Request.Header.Set("User-Agent", "codex_cli_rs/0.98.0")
	account := newOpenAIImageGenerationControlTestAccount()

	// /responses/compact 上游不接受 tool_choice，bridge 注入必须整体豁免 compact 请求。
	result, err := svc.Forward(context.Background(), c, account, []byte(`{"model":"gpt-5.4","input":"summarize the conversation","stream":falseplaceholder`))

placeholder
	require.NotNil(t, result)
	require.NotNil(t, upstream.lastReq)
	require.False(t, gjson.GetBytes(upstream.lastBody, "tool_choice").Exists())
	require.False(t, gjson.GetBytes(upstream.lastBody, `tools.#(type=="image_generation")`).Exists())
	instructions := gjson.GetBytes(upstream.lastBody, "instructions").String()
	require.NotContains(t, instructions, "image_generation")
placeholder

func TestOpenAIGatewayService_CodexImageGenerationBridgeOverridePrecedence(t *testing.T) {
	groupID := int64(4242)

	tests := []struct {
		name    string
		global  bool
		channel *Channel
		account *Account
		want    bool
placeholder{
		{
			name:   "global default enables bridge",
			global: true,
			account: &Account{
				Platform: PlatformOpenAI,
		placeholder,
			want: true,
	placeholder,
		{
			name:   "channel true overrides disabled global",
			global: false,
			channel: &Channel{ID: 1, Status: StatusActive, FeaturesConfig: map[string]any{
				featureKeyCodexImageGenerationBridge: map[string]any{PlatformOpenAI: trueplaceholder,
	placeholder
			account: &Account{Platform: PlatformOpenAIplaceholder,
			want:    true,
	placeholder,
		{
			name:   "channel false overrides enabled global",
			global: true,
			channel: &Channel{ID: 1, Status: StatusActive, FeaturesConfig: map[string]any{
				featureKeyCodexImageGenerationBridge: map[string]any{PlatformOpenAI: falseplaceholder,
	placeholder
			account: &Account{Platform: PlatformOpenAIplaceholder,
			want:    false,
	placeholder,
		{
			name:   "account false overrides channel and global true",
			global: true,
			channel: &Channel{ID: 1, Status: StatusActive, FeaturesConfig: map[string]any{
				featureKeyCodexImageGenerationBridge: map[string]any{PlatformOpenAI: trueplaceholder,
	placeholder
			account: &Account{
				Platform: PlatformOpenAI,
				Extra:    map[string]any{featureKeyCodexImageGenerationBridge: falseplaceholder,
		placeholder,
			want: false,
	placeholder,
		{
			name:   "nested account true overrides channel false",
			global: false,
			channel: &Channel{ID: 1, Status: StatusActive, FeaturesConfig: map[string]any{
				featureKeyCodexImageGenerationBridge: map[string]any{PlatformOpenAI: falseplaceholder,
	placeholder
			account: &Account{
				Platform: PlatformOpenAI,
				Extra: map[string]any{
					PlatformOpenAI: map[string]any{"codex_image_generation_bridge_enabled": trueplaceholder,
			placeholder,
		placeholder,
			want: true,
	placeholder,
		{
			name:   "non openai account extra is ignored",
			global: false,
			account: &Account{
				Platform: PlatformAnthropic,
				Extra:    map[string]any{featureKeyCodexImageGenerationBridge: trueplaceholder,
		placeholder,
			want: false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newOpenAIImageGenerationControlTestService(&httpUpstreamRecorder{placeholder)
			svc.cfg.Gateway.CodexImageGenerationBridgeEnabled = tt.global
			if tt.channel != nil {
				svc.channelService = newOpenAIImageGenerationControlChannelService(groupID, tt.channel)
		placeholder
			apiKey := &APIKey{GroupID: &groupIDplaceholder

			got := svc.isCodexImageGenerationBridgeEnabled(context.Background(), tt.account, apiKey)

			require.Equal(t, tt.want, got)
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayServiceHandleResponsesImageOutputs_NonStreaming(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newOpenAIImageGenerationControlTestService(&httpUpstreamRecorder{placeholder)
	c, _ := newOpenAIImageGenerationControlTestContext(true, "unit-test-agent/1.0")
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(`{
			"id":"resp_image_json",
			"model":"gpt-5.4",
			"output":[{"id":"ig_json_1","type":"image_generation_call","result":"final-image"placeholder],
			"usage":{"input_tokens":7,"output_tokens":3,"output_tokens_details":{"image_tokens":2placeholderplaceholder
	placeholder`)),
placeholder

	result, err := svc.handleNonStreamingResponse(context.Background(), resp, c, &Account{ID: 1, Type: AccountTypeAPIKeyplaceholder, "gpt-5.4", "gpt-5.4")

placeholder
	require.NotNil(t, result)
	require.Equal(t, 1, result.imageCount)
	require.NotNil(t, result.usage)
	require.Equal(t, 7, result.usage.InputTokens)
	require.Equal(t, 3, result.usage.OutputTokens)
	require.Equal(t, 2, result.usage.ImageOutputTokens)
placeholder

func TestOpenAIGatewayServiceHandleResponsesImageOutputs_Streaming(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newOpenAIImageGenerationControlTestService(&httpUpstreamRecorder{placeholder)
	c, recorder := newOpenAIImageGenerationControlTestContext(true, "unit-test-agent/1.0")
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(
			"data: {\"type\":\"response.output_item.done\",\"item\":{\"id\":\"ig_stream_1\",\"type\":\"image_generation_call\",\"status\":\"generating\",\"result\":\"final-image\"placeholderplaceholder\n\n" +
				"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_image_stream\",\"model\":\"gpt-5.5\",\"output\":[{\"id\":\"ig_stream_1\",\"type\":\"image_generation_call\",\"status\":\"generating\",\"result\":\"final-image\"placeholder],\"usage\":{\"input_tokens\":11,\"output_tokens\":5,\"output_tokens_details\":{\"image_tokens\":4placeholderplaceholderplaceholderplaceholder\n\n",
		)),
placeholder

	result, err := svc.handleStreamingResponse(context.Background(), resp, c, &Account{ID: 1placeholder, time.Now(), "gpt-5.5", "gpt-5.5")

placeholder
	require.NotNil(t, result)
	require.Equal(t, 1, result.imageCount)
	require.NotNil(t, result.usage)
	require.Equal(t, 11, result.usage.InputTokens)
	require.Equal(t, 5, result.usage.OutputTokens)
	require.Equal(t, 4, result.usage.ImageOutputTokens)
	require.NotContains(t, recorder.Body.String(), `"status":"generating"`)
	require.Equal(t, 2, strings.Count(recorder.Body.String(), `"status":"completed"`))
placeholder

func TestOpenAIGatewayServiceHandleResponsesImageOutputs_StreamingPassthrough(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newOpenAIImageGenerationControlTestService(&httpUpstreamRecorder{placeholder)
	c, recorder := newOpenAIImageGenerationControlTestContext(true, "unit-test-agent/1.0")
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(
			"data: {\"type\":\"response.output_item.done\",\"item\":{\"id\":\"ig_stream_1\",\"type\":\"image_generation_call\",\"status\":\"in_progress\",\"result\":\"final-image\"placeholderplaceholder\n\n" +
				"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_image_stream\",\"model\":\"gpt-5.5\",\"output\":[{\"id\":\"ig_stream_1\",\"type\":\"image_generation_call\",\"status\":\"in_progress\",\"result\":\"final-image\"placeholder],\"usage\":{\"input_tokens\":11,\"output_tokens\":5,\"output_tokens_details\":{\"image_tokens\":4placeholderplaceholderplaceholderplaceholder\n\n",
		)),
placeholder

	result, err := svc.handleStreamingResponsePassthrough(context.Background(), resp, c, &Account{ID: 1placeholder, time.Now(), "gpt-5.5", "gpt-5.5")

placeholder
	require.NotNil(t, result)
	require.NotContains(t, recorder.Body.String(), `"status":"in_progress"`)
	require.Equal(t, 2, strings.Count(recorder.Body.String(), `"status":"completed"`))
placeholder

func TestNormalizeCompletedImageGenerationStatus(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        string
		wantChanged bool
placeholder{
		{
			name:        "output item done with result",
			input:       `{"type":"response.output_item.done","item":{"type":"image_generation_call","status":"generating","result":"image-data"placeholderplaceholder`,
			want:        `{"type":"response.output_item.done","item":{"type":"image_generation_call","status":"completed","result":"image-data"placeholderplaceholder`,
			wantChanged: true,
	placeholder,
		{
			name:        "terminal response only changes completed image result",
			input:       `{"type":"response.completed","response":{"output":[{"type":"image_generation_call","status":"in_progress","result":"image-data"placeholder,{"type":"image_generation_call","status":"failed","result":"partial-data"placeholder]placeholderplaceholder`,
			want:        `{"type":"response.completed","response":{"output":[{"type":"image_generation_call","status":"completed","result":"image-data"placeholder,{"type":"image_generation_call","status":"failed","result":"partial-data"placeholder]placeholderplaceholder`,
			wantChanged: true,
	placeholder,
		{
			name:        "done item without result",
			input:       `{"type":"response.output_item.done","item":{"type":"image_generation_call","status":"generating"placeholderplaceholder`,
			want:        `{"type":"response.output_item.done","item":{"type":"image_generation_call","status":"generating"placeholderplaceholder`,
			wantChanged: false,
	placeholder,
		{
			name:        "non-final image event",
			input:       `{"type":"response.output_item.added","item":{"type":"image_generation_call","status":"generating","result":"image-data"placeholderplaceholder`,
			want:        `{"type":"response.output_item.added","item":{"type":"image_generation_call","status":"generating","result":"image-data"placeholderplaceholder`,
			wantChanged: false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, changed := normalizeCompletedImageGenerationStatus([]byte(tt.input))

			require.Equal(t, tt.wantChanged, changed)
			require.JSONEq(t, tt.want, string(got))
	placeholder)
placeholder
placeholder

// TestHandleStreamingResponse_CyberPolicyCapturesRealUpstreamTokens 锁定流式
// /v1/responses 命中 cyber_policy 的计费正确性：response.failed 自带的真实 usage
// 必须在打 cyber 标记前被解析进 mark；否则计费走 mark.UpstreamInTok 会按 0 token
// 漏记真实用量（该路径返回错误，handler 仅经 RecordCyberPolicyUsageLog 计费）。
func TestHandleStreamingResponse_CyberPolicyCapturesRealUpstreamTokens(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newOpenAIImageGenerationControlTestService(&httpUpstreamRecorder{placeholder)
	c, _ := newOpenAIImageGenerationControlTestContext(false, "unit-test-agent/1.0")
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(
			"data: {\"type\":\"response.created\",\"response\":{\"id\":\"resp_cyber\"placeholderplaceholder\n\n" +
				"data: {\"type\":\"response.failed\",\"response\":{\"id\":\"resp_cyber\",\"error\":{\"code\":\"cyber_policy\",\"message\":\"blocked by network policy\"placeholder,\"usage\":{\"input_tokens\":1234,\"output_tokens\":7placeholderplaceholderplaceholder\n\n",
		)),
placeholder

	_, err := svc.handleStreamingResponse(context.Background(), resp, c, &Account{ID: 1placeholder, time.Now(), "gpt-5.5", "gpt-5.5")
	require.Error(t, err, "cyber 命中的流式响应应返回错误（sawFailedEvent）")

	mark := GetOpsCyberPolicy(c)
	require.NotNil(t, mark, "必须打上 cyber 标记")
	require.Equal(t, "cyber_policy", mark.Code)
	require.Equal(t, 1234, mark.UpstreamInTok, "必须捕获 response.failed 自带真实 input token，而非解析前的 0")
	require.Equal(t, 7, mark.UpstreamOutTok)
placeholder

func newOpenAIImageGenerationControlTestService(upstream *httpUpstreamRecorder) *OpenAIGatewayService {
	cfg := &config.Config{placeholder
	return &OpenAIGatewayService{
		cfg:              cfg,
		httpUpstream:     upstream,
		cache:            &stubGatewayCache{placeholder,
		openaiWSResolver: NewOpenAIWSProtocolResolver(cfg),
		toolCorrector:    NewCodexToolCorrector(),
placeholder
placeholder

func newOpenAIImageGenerationControlChannelService(groupID int64, ch *Channel) *ChannelService {
	svc := &ChannelService{placeholder
	cache := newEmptyChannelCache()
	if ch != nil {
		cache.channelByGroupID[groupID] = ch
		cache.byID[ch.ID] = ch
placeholder
	cache.loadedAt = time.Now()
	svc.cache.Store(cache)
	return svc
placeholder

func newOpenAIImageGenerationControlTestContext(allowImages bool, userAgent string) (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	c.Request.Header.Set("User-Agent", userAgent)
	groupID := int64(4242)
	c.Set("api_key", &APIKey{
		ID:      2424,
		GroupID: &groupID,
		Group: &Group{
			ID:                   groupID,
			AllowImageGeneration: allowImages,
			RateMultiplier:       1,
			ImageRateMultiplier:  1,
	placeholder,
placeholder)
	return c, recorder
placeholder

func newOpenAIImageGenerationControlTestAccount() *Account {
placeholder
		ID:          5151,
		Name:        "openai-image-controls",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
placeholder
			"api_key": "sk-test",
	placeholder,
placeholder
placeholder
