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

	"github.com/Wei-Shaw/sub2api/internal/config"
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

func TestPatchGrokResponsesBodyDropsRedundantViewImageForCurrentInlineImage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
placeholder{
		{
			name: "top-level tools",
			body: `{
				"model":"grok-4.6",
				"input":[{"type":"message","role":"user","content":[
					{"type":"input_text","text":"What text is in this image?"placeholder,
					{"type":"input_image","image_url":"data:image/png;base64,AA=="placeholder
				]placeholder],
				"tools":[
					{"type":"function","name":"view_image","parameters":{"type":"object"placeholderplaceholder,
					{"type":"function","name":"shell_command","parameters":{"type":"object"placeholderplaceholder
				]
		placeholder`,
	placeholder,
		{
			name: "Responses Lite additional tools",
			body: `{
				"model":"grok-4.6",
				"input":[
					{"type":"additional_tools","role":"developer","tools":[
						{"type":"function","name":"view_image","parameters":{"type":"object"placeholderplaceholder,
						{"type":"function","name":"shell_command","parameters":{"type":"object"placeholderplaceholder
					]placeholder,
					{"type":"message","role":"user","content":[
						{"type":"input_text","text":"What text is in this image?"placeholder,
						{"type":"input_image","image_url":"data:image/png;base64,AA=="placeholder
					]placeholder
				]
		placeholder`,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			patched, err := patchGrokResponsesBody([]byte(tt.body), "grok-4.6")
		placeholder
			require.False(t, gjson.GetBytes(patched, `tools.#(name=="view_image")`).Exists())
			require.Equal(t, "shell_command", gjson.GetBytes(patched, "tools.0.name").String())
	placeholder)
placeholder
placeholder

func TestPatchGrokResponsesBodyKeepsNonRedundantViewImage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
placeholder{
		{
			name: "current turn has no inline image",
			body: `{"input":[{"role":"user","content":[{"type":"input_text","text":"Inspect a local image"placeholder]placeholder],"tools":[{"type":"function","name":"view_image"placeholder]placeholder`,
	placeholder,
		{
			name: "inline image is only historical",
			body: `{"input":[{"role":"user","content":[{"type":"input_image","image_url":"data:image/png;base64,AA=="placeholder]placeholder,{"role":"assistant","content":[{"type":"output_text","text":"Done"placeholder]placeholder,{"role":"user","content":[{"type":"input_text","text":"Inspect another local image"placeholder]placeholder],"tools":[{"type":"function","name":"view_image"placeholder]placeholder`,
	placeholder,
		{
			name: "view image is explicitly selected",
			body: `{"input":[{"role":"user","content":[{"type":"input_image","image_url":"data:image/png;base64,AA=="placeholder]placeholder],"tools":[{"type":"function","name":"view_image"placeholder],"tool_choice":{"type":"function","name":"view_image"placeholderplaceholder`,
	placeholder,
		{
			name: "required with view image as the only tool",
			body: `{"input":[{"role":"user","content":[{"type":"input_image","image_url":"data:image/png;base64,AA=="placeholder]placeholder],"tools":[{"type":"function","name":"view_image"placeholder],"tool_choice":"required"placeholder`,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			patched, err := patchGrokResponsesBody([]byte(tt.body), "grok-4.6")
		placeholder
			require.Equal(t, "view_image", gjson.GetBytes(patched, "tools.0.name").String())
	placeholder)
placeholder
placeholder

func TestPatchGrokResponsesBodyDropsViewImageOnlyToolMetadata(t *testing.T) {
	t.Parallel()

	body := []byte(`{
		"input":[{"role":"user","content":[{"type":"input_image","image_url":"data:image/png;base64,AA=="placeholder]placeholder],
		"tools":[{"type":"function","name":"view_image"placeholder],
		"tool_choice":"auto",
		"parallel_tool_calls":true
placeholder`)
	patched, err := patchGrokResponsesBody(body, "grok-4.6")
placeholder
	require.False(t, gjson.GetBytes(patched, "tools").Exists())
	require.False(t, gjson.GetBytes(patched, "tool_choice").Exists())
	require.False(t, gjson.GetBytes(patched, "parallel_tool_calls").Exists())
placeholder

func TestPatchGrokResponsesBodySanitizesComposerReasoningParameters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		upstreamModel string
		wantReasoning bool
placeholder{
		{name: "composer fast", upstreamModel: "grok-composer-2.5-fast"placeholder,
		{name: "composer shorthand", upstreamModel: "grok-composer"placeholder,
		{name: "composer legacy alias", upstreamModel: "composer-2.5"placeholder,
		{name: "provider-prefixed composer", upstreamModel: "xai/grok-composer-2.5-fast"placeholder,
		{name: "grok 4.5", upstreamModel: "grok-4.5", wantReasoning: trueplaceholder,
		{name: "grok 4.6", upstreamModel: "grok-4.6", wantReasoning: trueplaceholder,
		{name: "grok 4.6 latest", upstreamModel: "grok-4.6-latest", wantReasoning: trueplaceholder,
placeholder

	bodyTemplate := []byte(`{
		"model": "grok",
		"input": "hello",
		"reasoning": {"effort": "medium", "summary": "auto"placeholder,
		"reasoning_effort": "medium",
		"reasoningEffort": "medium"
placeholder`)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patched, err := patchGrokResponsesBody(append([]byte(nil), bodyTemplate...), tt.upstreamModel)
		placeholder
			require.True(t, json.Valid(patched))
			require.Equal(t, tt.upstreamModel, gjson.GetBytes(patched, "model").String())

			if tt.wantReasoning {
				require.Equal(t, "medium", gjson.GetBytes(patched, "reasoning.effort").String())
				require.Equal(t, "medium", gjson.GetBytes(patched, "reasoning_effort").String())
				require.False(t, gjson.GetBytes(patched, "reasoningEffort").Exists())
				return
		placeholder

			require.False(t, gjson.GetBytes(patched, "reasoning").Exists())
			require.False(t, gjson.GetBytes(patched, "reasoning_effort").Exists())
			require.False(t, gjson.GetBytes(patched, "reasoningEffort").Exists())
	placeholder)
placeholder
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

func TestPatchGrokResponsesBodyDropsLogprobsForGrok420Family(t *testing.T) {
	t.Parallel()
	body := []byte(`{"model":"grok-4.20-0309-reasoning","input":"hello","logprobs":true,"top_logprobs":5placeholder`)
	patched, err := patchGrokResponsesBody(body, "grok-4.20-0309-reasoning")
placeholder
	require.False(t, gjson.GetBytes(patched, "logprobs").Exists())
	require.False(t, gjson.GetBytes(patched, "top_logprobs").Exists())
placeholder

func TestPatchGrokResponsesBodyNormalizesReasoningEffortAliases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		body          string
		upstreamModel string
		path          string
		want          string
placeholder{
		{name: "minimal nested", body: `{"input":"hi","reasoning":{"effort":"minimal"placeholderplaceholder`, upstreamModel: "grok-4.5", path: "reasoning.effort", want: "low"placeholder,
		{name: "xhigh stays high for 4.5", body: `{"input":"hi","reasoning_effort":"xhigh"placeholder`, upstreamModel: "grok-4.5", path: "reasoning_effort", want: "high"placeholder,
		{name: "xhigh nested for 4.6", body: `{"input":"hi","reasoning":{"effort":"xhigh"placeholderplaceholder`, upstreamModel: "grok-4.6", path: "reasoning.effort", want: "xhigh"placeholder,
		{name: "xhigh snake for 4.6 latest", body: `{"input":"hi","reasoning_effort":"xhigh"placeholder`, upstreamModel: "grok-4.6-latest", path: "reasoning_effort", want: "xhigh"placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patched, err := patchGrokResponsesBody([]byte(tt.body), tt.upstreamModel)
		placeholder
			require.Equal(t, tt.want, gjson.GetBytes(patched, tt.path).String(), string(patched))
			require.False(t, gjson.GetBytes(patched, "reasoningEffort").Exists())
	placeholder)
placeholder
placeholder

func TestPatchGrokResponsesBodyAddsDefaultFunctionParameters(t *testing.T) {
	patched, err := patchGrokResponsesBody(
		[]byte(`{"input":"hi","tools":[{"type":"function","name":"lookup"placeholder,{"type":"function","name":"wait","parameters":nullplaceholder]placeholder`),
		"grok-4.5",
	)
placeholder
	for _, tool := range gjson.GetBytes(patched, "tools").Array() {
		require.Equal(t, "object", tool.Get("parameters.type").String(), string(patched))
		require.True(t, tool.Get("parameters.properties").IsObject(), string(patched))
placeholder
placeholder

func TestNormalizeGrokChatReasoningEffort(t *testing.T) {
	patched, err := normalizeGrokChatReasoningEffort([]byte(`{"reasoningEffort":"ultra"placeholder`), "grok-4.3")
placeholder
	require.Equal(t, "high", gjson.GetBytes(patched, "reasoning_effort").String())
	require.False(t, gjson.GetBytes(patched, "reasoningEffort").Exists())

	patched, err = normalizeGrokChatReasoningEffort([]byte(`{"reasoning_effort":"xhigh"placeholder`), "grok-4.6")
placeholder
	require.Equal(t, "xhigh", gjson.GetBytes(patched, "reasoning_effort").String())

	patched, err = normalizeGrokChatReasoningEffort([]byte(`{"reasoning_effort":"high"placeholder`), "grok-composer-2.5-fast")
placeholder
	require.False(t, gjson.GetBytes(patched, "reasoning_effort").Exists())
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

func TestPatchGrokResponsesBodyFlattensNamespaceTools(t *testing.T) {
	t.Parallel()

	body := []byte(`{
		"model": "grok",
		"input": "hello",
		"tools": [
			{"type": "namespace", "name": "functions", "tools": [{"type": "function", "name": "inner"placeholder]placeholder,
			{"type": "function", "name": "kept_fn", "parameters": {"type": "object"placeholderplaceholder,
			{"type": "shell", "name": "kept_shell"placeholder
		],
		"tool_choice": {"type": "function", "namespace": "functions", "name": "inner"placeholder
placeholder`)

	patched, _, err := patchGrokResponsesBodyWithClientTools(body, "grok-4.3")
placeholder
	require.True(t, json.Valid(patched))
	require.Equal(t, "grok-4.3", gjson.GetBytes(patched, "model").String())
	require.Len(t, gjson.GetBytes(patched, "tools").Array(), 3)
	require.False(t, gjson.GetBytes(patched, `tools.#(type=="namespace")`).Exists())
	require.True(t, gjson.GetBytes(patched, `tools.#(type=="function")`).Exists())
	require.True(t, gjson.GetBytes(patched, `tools.#(type=="shell")`).Exists())
	require.Equal(t, "functions__inner", gjson.GetBytes(patched, "tools.0.name").String())
	require.Equal(t, "functions__inner", gjson.GetBytes(patched, "tool_choice.name").String())
	require.False(t, gjson.GetBytes(patched, "tool_choice.namespace").Exists())
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

func TestSanitizeGrokResponsesToolsKeepsToolChoiceOnlyWithSupportedTools(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		body           string
		wantTools      bool
		wantToolChoice bool
placeholder{
		{
			name: "missing tools with string tool choice",
			body: `{"input":"hello","tool_choice":"auto"placeholder`,
	placeholder,
		{
			name: "missing tools with object tool choice",
			body: `{"input":"hello","tool_choice":{"type":"function","name":"lookup"placeholderplaceholder`,
	placeholder,
		{
			name:      "empty tools",
			body:      `{"input":"hello","tools":[],"tool_choice":"auto"placeholder`,
			wantTools: true,
	placeholder,
		{
			name: "all tools unsupported",
			body: `{"input":"hello","tools":[{"type":"namespace","name":"client_tools"placeholder],"tool_choice":"auto"placeholder`,
	placeholder,
		{
			name:           "supported tool",
			body:           `{"input":"hello","tools":[{"type":"function","name":"lookup"placeholder],"tool_choice":"auto"placeholder`,
			wantTools:      true,
			wantToolChoice: true,
	placeholder,
		{
			name:           "malformed non-array tools are removed",
			body:           `{"input":"hello","tools":{"type":"function","name":"lookup"placeholder,"tool_choice":"auto"placeholder`,
			wantTools:      false,
			wantToolChoice: false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patched, err := sanitizeGrokResponsesTools([]byte(tt.body))
		placeholder
			require.True(t, json.Valid(patched))
			require.Equal(t, tt.wantTools, gjson.GetBytes(patched, "tools").Exists())
			require.Equal(t, tt.wantToolChoice, gjson.GetBytes(patched, "tool_choice").Exists())
			if tt.wantToolChoice {
				require.Equal(t, "auto", gjson.GetBytes(patched, "tool_choice").String())
		placeholder
	placeholder)
placeholder
placeholder

func TestPatchGrokResponsesBodyPromotesCodexAdditionalTools(t *testing.T) {
	t.Parallel()

	body := []byte(`{
		"model": "grok",
		"tools": [
			{"type": "function", "name": "existing", "description": "top-level wins"placeholder,
			{"type": "web_search"placeholder
		],
		"tool_choice": "auto",
		"input": [
			{
				"type": "additional_tools",
				"role": "developer",
				"tools": [
					{"type": "function", "name": "existing", "description": "duplicate carrier definition"placeholder,
					{"type": "function", "name": "wait"placeholder,
					{"type": "web_search"placeholder,
					{"type": "shell"placeholder,
					{"type": "custom", "name": "apply_patch"placeholder,
					{"type": "namespace", "name": "collaboration"placeholder
				]
		placeholder,
			{
				"type": "message",
				"role": "developer",
				"content": [{"type": "input_text", "text": "system prompt"placeholder]
		placeholder,
			{
				"type": "message",
				"role": "user",
				"content": [{"type": "input_text", "text": "hello"placeholder]
		placeholder
		]
placeholder`)

	patched, _, err := patchGrokResponsesBodyWithClientTools(body, "grok-4.5")
placeholder
	require.True(t, json.Valid(patched))
	require.Equal(t, "grok-4.5", gjson.GetBytes(patched, "model").String())
	require.Equal(t, 2, len(gjson.GetBytes(patched, "input").Array()))
	require.False(t, gjson.GetBytes(patched, `input.#(type=="additional_tools")`).Exists())
	tools := gjson.GetBytes(patched, "tools").Array()
	require.Len(t, tools, 5)
	require.Equal(t, "existing", tools[0].Get("name").String())
	require.Equal(t, "top-level wins", tools[0].Get("description").String())
	require.Equal(t, "web_search", tools[1].Get("type").String())
	require.Equal(t, "wait", tools[2].Get("name").String())
	require.Equal(t, "shell", tools[3].Get("type").String())
	require.Equal(t, "function", tools[4].Get("type").String())
	require.Equal(t, "apply_patch", tools[4].Get("name").String())
	require.Equal(t, "string", tools[4].Get("parameters.properties.input.type").String())
	require.False(t, gjson.GetBytes(patched, `tools.#(type=="custom")`).Exists())
	require.False(t, gjson.GetBytes(patched, `tools.#(type=="namespace")`).Exists())
	require.Equal(t, "auto", gjson.GetBytes(patched, "tool_choice").String())
	require.Equal(t, "developer", gjson.GetBytes(patched, "input.0.role").String())
	require.Equal(t, "system prompt", gjson.GetBytes(patched, "input.0.content.0.text").String())
	require.Equal(t, "user", gjson.GetBytes(patched, "input.1.role").String())
	require.Equal(t, "hello", gjson.GetBytes(patched, "input.1.content.0.text").String())
placeholder

func TestForwardGrokResponsesCodexAdditionalToolsUsesMixedCacheIntent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{
		"model":"grok",
		"stream":false,
		"prompt_cache_key":"codex-session",
		"input":[
			{"type":"additional_tools","role":"developer","tools":[
				{"type":"function","name":"lookup","description":"look up a key","parameters":{"type":"object"placeholderplaceholder,
				{"type":"function","name":"web_search","description":"search","parameters":{"type":"object"placeholderplaceholder,
				{"type":"custom","name":"apply_patch"placeholder,
				{"type":"namespace","name":"collaboration"placeholder
			]placeholder,
			{"type":"message","role":"user","content":[{"type":"input_text","text":"hello"placeholder]placeholder
		]
placeholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))
	c.Request.Header.Set(grokClientToolCacheOptInHeader, "prefer-cache")
	c.Set("api_key", &APIKey{ID: 4501placeholder)

	account := healthyGrokOAuthGatewayTestAccount(4501, "access-token")
	account.Credentials["subscription_tier"] = "free"
	repo := &grokQuotaAccountRepo{
		mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
			accountsByID: map[int64]*Account{account.ID: accountplaceholder,
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(`{
			"id":"resp_codex_lite","object":"response","model":"grok-4.5","status":"completed",
			"output":[],"usage":{"input_tokens":10,"output_tokens":1placeholder
	placeholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{
		httpUpstream:      upstream,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		accountRepo:       repo,
placeholder

	result, err := svc.forwardGrokResponses(context.Background(), c, account, body, "grok", false, time.Now())

placeholder
	require.NotNil(t, result)
	require.Equal(t, "resp_codex_lite", result.ResponseID)
	require.False(t, gjson.GetBytes(upstream.lastBody, `input.#(type=="additional_tools")`).Exists())
	tools := gjson.GetBytes(upstream.lastBody, "tools").Array()
	require.Len(t, tools, 4)
	require.Equal(t, "function", tools[0].Get("type").String())
	require.Equal(t, "lookup", tools[0].Get("name").String())
	require.Equal(t, "web_search", tools[1].Get("type").String())
	require.Equal(t, "function", tools[2].Get("type").String())
	require.Equal(t, "apply_patch", tools[2].Get("name").String())
	require.Equal(t, "x_search", tools[3].Get("type").String())
	require.False(t, gjson.GetBytes(upstream.lastBody, "tool_choice").Exists())
	require.False(t, gjson.GetBytes(upstream.lastBody, `tools.#(type=="custom")`).Exists())
	require.False(t, gjson.GetBytes(upstream.lastBody, `tools.#(type=="namespace")`).Exists())
	identity := gjson.GetBytes(upstream.lastBody, "prompt_cache_key").String()
	require.NotEmpty(t, identity)
	require.Equal(t, identity, upstream.lastReq.Header.Get(grokConversationIDHeader))
	require.Empty(t, upstream.lastReq.Header.Get(grokClientToolCacheOptInHeader))
placeholder

func TestForwardGrokResponsesClaudeDesktopClientToolsUseCacheRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	firstBody := []byte(`{
		"model":"grok","stream":false,"instructions":"You are Claude Desktop.",
		"tools":[
			{"type":"function","name":"Read","parameters":{"type":"object"placeholderplaceholder,
			{"type":"function","name":"Edit","parameters":{"type":"object"placeholderplaceholder,
			{"type":"function","name":"WebSearch","parameters":{"type":"object"placeholderplaceholder,
			{"type":"function","name":"mcp__workspace__bash","parameters":{"type":"object"placeholderplaceholder
		],
		"input":[{"role":"user","content":[{"type":"input_text","text":"first turn"placeholder]placeholder]
placeholder`)
	secondBody := []byte(`{
		"model":"grok","stream":false,"instructions":"You are Claude Desktop.",
		"tools":[
			{"type":"function","name":"Read","parameters":{"type":"object"placeholderplaceholder,
			{"type":"function","name":"Edit","parameters":{"type":"object"placeholderplaceholder,
			{"type":"function","name":"WebSearch","parameters":{"type":"object"placeholderplaceholder,
			{"type":"function","name":"mcp__workspace__bash","parameters":{"type":"object"placeholderplaceholder
		],
		"input":[
			{"role":"user","content":[{"type":"input_text","text":"first turn"placeholder]placeholder,
			{"role":"assistant","content":[{"type":"output_text","text":"first answer"placeholder]placeholder,
			{"role":"user","content":[{"type":"input_text","text":"second turn"placeholder]placeholder
		]
placeholder`)

	account := healthyGrokOAuthGatewayTestAccount(4504, "access-token")
	account.Credentials["subscription_tier"] = "free"
	repo := &grokQuotaAccountRepo{
		mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
			accountsByID: map[int64]*Account{account.ID: accountplaceholder,
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{responses: []*http.Response{
		{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body: io.NopCloser(strings.NewReader(`{
				"id":"resp_claude_desktop_1","object":"response","model":"grok-4.5","status":"completed",
				"output":[],"usage":{"input_tokens":30000,"output_tokens":10,"input_tokens_details":{"cached_tokens":0placeholderplaceholder
		placeholder`)),
	placeholder,
		{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body: io.NopCloser(strings.NewReader(`{
				"id":"resp_claude_desktop_2","object":"response","model":"grok-4.5","status":"completed",
				"output":[],"usage":{"input_tokens":30100,"output_tokens":12,"input_tokens_details":{"cached_tokens":placeholder
		placeholder`)),
	placeholder,
placeholderplaceholder
	svc := &OpenAIGatewayService{
		httpUpstream:      upstream,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		accountRepo:       repo,
placeholder

	newContext := func(body []byte) *gin.Context {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))
		c.Request.Header.Set("User-Agent", "claude-cli/2.1.215 (external, claude-desktop-3p, agent-sdk/0.3.215)")
		c.Request.Header.Set("X-App", "cli")
		c.Request.Header.Set("anthropic-client-platform", "desktop_app")
		c.Request.Header.Set("X-Claude-Code-Session-Id", "claude-desktop-session")
		c.Set("api_key", &APIKey{ID: 4504placeholder)
		return c
placeholder

	first, err := svc.forwardGrokResponses(context.Background(), newContext(firstBody), account, firstBody, "grok", false, time.Now())
placeholder
	second, err := svc.forwardGrokResponses(context.Background(), newContext(secondBody), account, secondBody, "grok", false, time.Now())
placeholder

	require.Equal(t, 0, first.Usage.CacheReadInputTokens)
	require.Equal(t, 28672, second.Usage.CacheReadInputTokens)
	require.Len(t, upstream.bodies, 2)
	require.Len(t, upstream.requests, 2)
	for i := range upstream.bodies {
		tools := gjson.GetBytes(upstream.bodies[i], "tools").Array()
		require.Len(t, tools, 6)
		require.Equal(t, "Read", tools[0].Get("name").String())
		require.Equal(t, "Edit", tools[1].Get("name").String())
		require.Equal(t, "WebSearch", tools[2].Get("name").String())
		require.Equal(t, "mcp__workspace__bash", tools[3].Get("name").String())
		require.Equal(t, "web_search", tools[4].Get("type").String())
		require.Equal(t, "x_search", tools[5].Get("type").String())
		require.False(t, gjson.GetBytes(upstream.bodies[i], "tool_choice").Exists())
		require.Empty(t, upstream.requests[i].Header.Get("X-App"))
		require.Empty(t, upstream.requests[i].Header.Get("anthropic-client-platform"))
		require.Empty(t, upstream.requests[i].Header.Get("X-Claude-Code-Session-Id"))
placeholder
	firstIdentity := gjson.GetBytes(upstream.bodies[0], "prompt_cache_key").String()
	secondIdentity := gjson.GetBytes(upstream.bodies[1], "prompt_cache_key").String()
	require.NotEmpty(t, firstIdentity)
	require.Equal(t, firstIdentity, secondIdentity)
	require.Equal(t, firstIdentity, upstream.requests[0].Header.Get(grokConversationIDHeader))
	require.Equal(t, secondIdentity, upstream.requests[1].Header.Get(grokConversationIDHeader))
placeholder

func TestGrokResponsesCacheIdentityIncludesPromotedCodexTools(t *testing.T) {
	c := newGrokCacheTestContext(4503)
	lookupBody := []byte(`{"model":"grok","input":[{"type":"additional_tools","tools":[{"type":"function","name":"lookup","parameters":{"type":"object"placeholderplaceholder]placeholder,{"type":"message","role":"user","content":"same prompt"placeholder]placeholder`)
	readBody := []byte(`{"model":"grok","input":[{"type":"additional_tools","tools":[{"type":"function","name":"read_file","parameters":{"type":"object"placeholderplaceholder]placeholder,{"type":"message","role":"user","content":"same prompt"placeholder]placeholder`)

	patchedLookup, err := patchGrokResponsesBody(lookupBody, "grok-4.5")
placeholder
	patchedRead, err := patchGrokResponsesBody(readBody, "grok-4.5")
placeholder

	lookupIdentity := resolveGrokCacheIdentity(c, patchedLookup, "", "grok-4.5")
	readIdentity := resolveGrokCacheIdentity(c, patchedRead, "", "grok-4.5")
	require.NotEmpty(t, lookupIdentity)
	require.NotEmpty(t, readIdentity)
	require.NotEqual(t, lookupIdentity, readIdentity)
placeholder

func TestCodexUnsupportedAdditionalToolsDoNotBecomeToolFreeCacheIntent(t *testing.T) {
	body := []byte(`{
		"model":"grok","tool_choice":"auto",
		"input":[
			{"type":"additional_tools","role":"developer","tools":[
				{"type":"custom","name":"apply_patch"placeholder,
				{"type":"namespace","name":"collaboration"placeholder
			]placeholder,
			{"type":"message","role":"user","content":"hello"placeholder
		]
placeholder`)
	patched, err := patchGrokResponsesBody(body, "grok-4.5")
placeholder
	require.False(t, gjson.GetBytes(patched, "tools").Exists())
	require.False(t, gjson.GetBytes(patched, "tool_choice").Exists())

	mixedCacheIntent := patched
	patched, err = applyGrokResponsesCacheIdentity(patched, body, "isolated-id", true)
placeholder
	account := healthyGrokOAuthGatewayTestAccount(4502, "access-token")
	account.Credentials["subscription_tier"] = "free"
	patched, err = applyGrokFreeRequestToolCacheRoute(nil, patched, mixedCacheIntent, account, "isolated-id")

placeholder
	require.False(t, gjson.GetBytes(patched, "tools").Exists())
	require.False(t, gjson.GetBytes(patched, "tool_choice").Exists())
	require.Equal(t, "isolated-id", gjson.GetBytes(patched, "prompt_cache_key").String())
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

	req, err := buildGrokResponsesRequest(context.Background(), nil, account, []byte(`{"model":"grok-4.3"placeholder`), "access-token", "isolated-cache-id", nil)
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

func TestBuildGrokCompactRequestBodyUsesResponsesCompactionTurn(t *testing.T) {
	body := []byte(`{"model":"grok-4.5","input":[{"type":"message","role":"user","content":[{"type":"input_text","text":"hello"placeholder]placeholder],"tools":[{"type":"function","name":"shell"placeholder],"stream":trueplaceholder`)

	patched, err := buildGrokCompactRequestBody(body)
placeholder
	require.False(t, gjson.GetBytes(patched, "stream").Bool())
	require.False(t, gjson.GetBytes(patched, "store").Bool())
	require.Equal(t, "none", gjson.GetBytes(patched, "tool_choice").String())
	require.Equal(t, "reasoning.encrypted_content", gjson.GetBytes(patched, "include.0").String())
	require.Equal(t, "hello", gjson.GetBytes(patched, "input.0.content.0.text").String())
	prompt := gjson.GetBytes(patched, "input.1.content.0.text").String()
	require.Contains(t, prompt, "1. Primary Request and Intent")
	require.Contains(t, prompt, "9. Optional Next Step")
	require.Contains(t, prompt, "Respond with ONLY the <summary>...</summary> block")
	require.NotContains(t, prompt, "<summary_request>")
placeholder

func TestConvertGrokResponseToOpenAICompact(t *testing.T) {
	body := []byte(`{
		"id":"resp_grok_1",
		"object":"response",
		"status":"completed",
		"model":"grok-4.5",
		"output":[
			{"id":"rs_1","type":"reasoning","summary":[],"encrypted_content":"grok-encrypted-state"placeholder,
			{"id":"msg_1","type":"message","role":"assistant","content":[{"type":"output_text","text":"summary text"placeholder]placeholder
		],
		"usage":{"input_tokens":10,"output_tokens":4,"total_tokens":14placeholder
placeholder`)

	converted, err := convertGrokResponseToOpenAICompact(body)
placeholder
	require.Equal(t, "resp_grok_1", gjson.GetBytes(converted, "id").String())
	require.Len(t, gjson.GetBytes(converted, "output").Array(), 1)
	require.Equal(t, "compaction", gjson.GetBytes(converted, "output.0.type").String())
	require.Equal(t, "grok-encrypted-state", gjson.GetBytes(converted, "output.0.encrypted_content").String())
	require.Equal(t, "summary text", gjson.GetBytes(converted, "output.0.summary.0.text").String())
	require.Equal(t, int64(14), gjson.GetBytes(converted, "usage.total_tokens").Int())
placeholder

func TestPatchGrokResponsesBodyRestoresCompactInput(t *testing.T) {
	body := []byte(`{
		"model":"grok-4.5",
		"input":[
			{"id":"cmp_1","type":"compaction","status":"completed","encrypted_content":"grok-encrypted-state","summary":[{"type":"summary_text","text":"summary text"placeholder]placeholder,
			{"type":"message","role":"user","content":[{"type":"input_text","text":"continue"placeholder]placeholder
		]
placeholder`)

	patched, err := patchGrokResponsesBody(body, "grok-4.5")
placeholder
	require.Equal(t, "reasoning", gjson.GetBytes(patched, "input.0.type").String())
	require.Equal(t, "grok-encrypted-state", gjson.GetBytes(patched, "input.0.encrypted_content").String())
	require.Equal(t, "message", gjson.GetBytes(patched, "input.1.type").String())
	require.Contains(t, gjson.GetBytes(patched, "input.1.content.0.text").String(), "summary text")
	require.Equal(t, "continue", gjson.GetBytes(patched, "input.2.content.0.text").String())
placeholder

func TestConvertGrokResponseToOpenAICompactRequiresEncryptedContent(t *testing.T) {
	_, err := convertGrokResponseToOpenAICompact([]byte(`{"output":[{"type":"message","content":[{"type":"output_text","text":"summary"placeholder]placeholder]placeholder`))
	require.ErrorContains(t, err, "reasoning.encrypted_content")
placeholder

func TestBuildGrokResponsesRequestAllowsPublicAPIKeyBaseURLByDefault(t *testing.T) {
	account := &Account{
		Platform: PlatformGrok,
		Type:     AccountTypeAPIKey,
placeholder
			"base_url": "https://grok.example.test/v1/",
	placeholder,
placeholder

	req, err := buildGrokResponsesRequest(context.Background(), nil, account, []byte(`{"model":"grok-4.3"placeholder`), "api-key", "", nil)
placeholder
	require.Equal(t, "https://grok.example.test/v1/responses", req.URL.String())
	require.Equal(t, "Bearer api-key", req.Header.Get("Authorization"))
	require.Empty(t, req.Header.Get("X-Grok-Client-Version"))
	require.NotEqual(t, grokUpstreamUserAgent, req.Header.Get("User-Agent"))
placeholder

func TestBuildGrokResponsesRequestHonorsOAuthOfficialEndpointSwitch(t *testing.T) {
	t.Parallel()

	account := &Account{
		Platform: PlatformGrok,
		Type:     AccountTypeOAuth,
placeholder
			"base_url": xai.DefaultBaseURL,
	placeholder,
placeholder

	req, err := buildGrokResponsesRequest(context.Background(), nil, account, []byte(`{"model":"grok-4.3"placeholder`), "access-token", "", nil)
placeholder
	require.Equal(t, xai.DefaultBaseURL+"/responses", req.URL.String())
placeholder

func TestBuildGrokResponsesRequestAppliesHeaderOverridesLast(t *testing.T) {
	t.Parallel()

	account := &Account{
		Platform: PlatformGrok,
		Type:     AccountTypeOAuth,
placeholder
			"base_url":                "https://relay.example.test/v1",
			"header_override_enabled": true,
			"header_overrides": map[string]any{
				"User-Agent":            "relay-client/2.0",
				"X-Grok-Client-Version": "9.9.9",
				"X-Relay-Token":         "relay-secret",
		placeholder,
	placeholder,
placeholder

	req, err := buildGrokResponsesRequest(context.Background(), nil, account, []byte(`{"model":"grok-4.3"placeholder`), "access-token", "conv-1", nil)
placeholder
	require.Equal(t, "https://relay.example.test/v1/responses", req.URL.String())
	// 覆写值优先于内置 CLI 身份头。名字不在 wire casing 映射中的覆写头
	// 以小写键直写（HTTP/2 线上语义），需按写入形态断言。
	require.Equal(t, "relay-client/2.0", req.Header.Get("User-Agent"))
	require.Equal(t, []string{"9.9.9"placeholder, req.Header["x-grok-client-version"])
	require.Empty(t, req.Header.Get("X-Grok-Client-Version"))
	require.Equal(t, []string{"relay-secret"placeholder, req.Header["x-relay-token"])
	// 会话路由头与认证头不受覆写影响。
	require.Equal(t, "conv-1", req.Header.Get(grokConversationIDHeader))
	require.Equal(t, "Bearer access-token", req.Header.Get("Authorization"))
placeholder

func TestBuildGrokResponsesRequestIgnoresBlockedHeaderOverrides(t *testing.T) {
	t.Parallel()

	account := &Account{
		Platform: PlatformGrok,
		Type:     AccountTypeAPIKey,
placeholder
			"header_override_enabled": true,
			"header_overrides": map[string]any{
				"Authorization":  "Bearer stolen",
				"x-grok-conv-id": "pinned-conversation",
		placeholder,
	placeholder,
placeholder

	req, err := buildGrokResponsesRequest(context.Background(), nil, account, []byte(`{"model":"grok-4.3"placeholder`), "api-key", "conv-2", nil)
placeholder
	require.Equal(t, "Bearer api-key", req.Header.Get("Authorization"))
	require.Equal(t, "conv-2", req.Header.Get(grokConversationIDHeader))
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

func TestParseGrokMediaRequestAcceptsOfficialImageURLFields(t *testing.T) {
	body := []byte(`{
		"model":"grok-imagine-video-1.5",
		"image":{"url":"https://example.com/source.png"placeholder,
		"reference_images":[{"url":"https://example.com/reference.png"placeholder]
placeholder`)

	info := ParseGrokMediaRequest("application/json", body)

	require.Equal(t, []string{
		"https://example.com/source.png",
		"https://example.com/reference.png",
placeholder, info.InputImageURLs)
	require.True(t, info.HasInputImage())
placeholder

func TestNormalizeGrokMediaForwardBodyCanonicalizesImageURLAlias(t *testing.T) {
	body := []byte(`{
		"model":"grok-imagine-video-1.5",
		"prompt":"animate",
		"image":{"image_url":"https://example.com/source.png"placeholder,
		"duration":8
placeholder`)

	out, contentType, err := normalizeGrokMediaForwardBody(GrokMediaEndpointVideosGenerations, body, "application/json")

placeholder
	require.Equal(t, "application/json", contentType)
	require.Equal(t, "grok-imagine-video-1.5", gjson.GetBytes(out, "model").String())
	require.Equal(t, "https://example.com/source.png", gjson.GetBytes(out, "image.url").String())
	require.False(t, gjson.GetBytes(out, "image.image_url").Exists())
placeholder

func TestNormalizeGrokMediaForwardBodyPreservesImageToVideoModelForOfficialURL(t *testing.T) {
	body := []byte(`{
		"model":"grok-imagine-video-1.5",
		"prompt":"animate",
		"image":{"url":"https://example.com/source.png"placeholder
placeholder`)

	out, _, err := normalizeGrokMediaForwardBody(GrokMediaEndpointVideosGenerations, body, "application/json")

placeholder
	require.Equal(t, "grok-imagine-video-1.5", gjson.GetBytes(out, "model").String())
	require.Equal(t, "https://example.com/source.png", gjson.GetBytes(out, "image.url").String())
placeholder

func TestCanonicalizeGrokMediaImageURLFieldsPreservesOfficialURL(t *testing.T) {
	body := []byte(`{
		"image":{"url":"https://example.com/official.png","image_url":"https://example.com/legacy.png"placeholder,
		"images":[
			{"image_url":"https://example.com/first.png"placeholder,
			{"url":"https://example.com/second.png"placeholder
		],
		"reference_images":[{"image_url":"https://example.com/reference.png"placeholder],
		"mask":{"image_url":"https://example.com/mask.png"placeholder
placeholder`)

	out, err := canonicalizeGrokMediaImageURLFields(body, "image", "images", "reference_images", "mask")

placeholder
	require.Equal(t, "https://example.com/official.png", gjson.GetBytes(out, "image.url").String())
	require.False(t, gjson.GetBytes(out, "image.image_url").Exists())
	require.Equal(t, "https://example.com/first.png", gjson.GetBytes(out, "images.0.url").String())
	require.False(t, gjson.GetBytes(out, "images.0.image_url").Exists())
	require.Equal(t, "https://example.com/second.png", gjson.GetBytes(out, "images.1.url").String())
	require.Equal(t, "https://example.com/reference.png", gjson.GetBytes(out, "reference_images.0.url").String())
	require.False(t, gjson.GetBytes(out, "reference_images.0.image_url").Exists())
	require.Equal(t, "https://example.com/mask.png", gjson.GetBytes(out, "mask.url").String())
	require.False(t, gjson.GetBytes(out, "mask.image_url").Exists())
placeholder

func TestCanonicalizeGrokMediaImageURLFieldsReplacesEmptyOfficialURL(t *testing.T) {
	body := []byte(`{"image":{"url":" ","image_url":"https://example.com/legacy.png"placeholderplaceholder`)

	out, err := canonicalizeGrokMediaImageURLFields(body, "image")

placeholder
	require.Equal(t, "https://example.com/legacy.png", gjson.GetBytes(out, "image.url").String())
	require.False(t, gjson.GetBytes(out, "image.image_url").Exists())
placeholder

func TestPrepareGrokImageEditNormalizesOfficialImageObjects(t *testing.T) {
	body := []byte(`{
		"model":"grok-imagine-image-quality",
		"image":{"image_url":{"url":"https://example.com/first.png"placeholderplaceholder,
		"images":["https://example.com/second.png"],
		"mask":{"image_url":"https://example.com/mask.png"placeholder
placeholder`)

	out, contentType, err := prepareGrokMediaForwardBody(GrokMediaEndpointImagesEdits, body, "application/json")
placeholder
	require.Equal(t, "application/json", contentType)
	for _, path := range []string{"image", "images.0", "mask"placeholder {
		require.Equal(t, "image_url", gjson.GetBytes(out, path+".type").String())
		require.NotEmpty(t, gjson.GetBytes(out, path+".url").String())
		require.False(t, gjson.GetBytes(out, path+".image_url").Exists())
placeholder
placeholder

func TestPrepareGrokImageEditRejectsMoreThanThreeSources(t *testing.T) {
	body := []byte(`{"images":["https://example.com/1.png","https://example.com/2.png","https://example.com/3.png","https://example.com/4.png"]placeholder`)

	out, _, err := prepareGrokMediaForwardBody(GrokMediaEndpointImagesEdits, body, "application/json")
placeholder
	require.Nil(t, out)
	require.Contains(t, err.Error(), "maximum of 3 source images")
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
		{name: "video 1.5 text-only remains explicit", endpoint: GrokMediaEndpointVideosGenerations, model: "grok-imagine-video-1.5", want: "grok-imagine-video-1.5"placeholder,
		{name: "video 1.5 image-to-video passthrough", endpoint: GrokMediaEndpointVideosGenerations, model: "grok-imagine-video-1.5", hasInputImage: true, want: "grok-imagine-video-1.5"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, NormalizeGrokMediaModelForEndpoint(tt.endpoint, tt.model, tt.hasInputImage))
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
		Body: io.NopCloser(strings.NewReader(`{"data":[{"url":"https://images.test/cat.png"placeholder]placeholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder

	result, err := svc.ForwardGrokMedia(context.Background(), c, account, GrokMediaEndpointImagesGenerations, "", body, "application/json")
placeholder
	require.Equal(t, "https://xai.test/v1/images/generations", upstream.lastReq.URL.String())
	require.Equal(t, http.MethodPost, upstream.lastReq.Method)
	require.Equal(t, "Bearer api-key", upstream.lastReq.Header.Get("Authorization"))
	require.Equal(t, "application/json", upstream.lastReq.Header.Get("Content-Type"))
	require.Empty(t, upstream.lastReq.Header.Get("X-Grok-Client-Version"))
	require.NotEqual(t, grokUpstreamUserAgent, upstream.lastReq.Header.Get("User-Agent"))
	require.JSONEq(t, `{"model":"grok-imagine-image-quality","prompt":"draw a cat"placeholder`, string(upstream.lastBody))
	require.Equal(t, http.StatusOK, recorder.Code)
	require.JSONEq(t, `{"data":[{"url":"https://images.test/cat.png"placeholder]placeholder`, recorder.Body.String())
	require.Equal(t, "xai-image-req", result.RequestID)
	require.Equal(t, "grok-imagine-image-quality", result.Model)
	require.Equal(t, "grok-imagine-image-quality", result.BillingModel)
	require.Equal(t, 1, result.ImageCount)
	require.Equal(t, ImageBillingSize2K, result.ImageSize)
placeholder

func TestForwardGrokMediaAppliesAccountModelMappingAfterEndpointNormalization(t *testing.T) {
	t.Setenv(xai.EnvAllowUnsafeURLOverrides, "true")
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		endpoint         GrokMediaEndpoint
		path             string
		body             string
		modelMapping     map[string]any
		wantRequestModel string
		wantUpstream     string
		wantBody         string
		responseBody     string
placeholder{
		{
			name:             "image generation maps normalized image alias",
			endpoint:         GrokMediaEndpointImagesGenerations,
			path:             "/v1/images/generations",
			body:             `{"model":"grok-imagine","prompt":"draw a cat"placeholder`,
			modelMapping:     map[string]any{"grok-imagine-image-quality": "vendor-image-model"placeholder,
			wantRequestModel: "grok-imagine-image-quality",
			wantUpstream:     "vendor-image-model",
			wantBody:         `{"model":"vendor-image-model","prompt":"draw a cat"placeholder`,
			responseBody:     `{"data":[{"url":"https://images.test/mapped.png"placeholder]placeholder`,
	placeholder,
		{
			name:             "video generation maps text-only fallback model",
			endpoint:         GrokMediaEndpointVideosGenerations,
			path:             "/v1/videos/generations",
			body:             `{"model":"grok-imagine-video-1.5","prompt":"waves"placeholder`,
			modelMapping:     map[string]any{"grok-imagine-video": "grok-image-video"placeholder,
			wantRequestModel: "grok-imagine-video-1.5",
			wantUpstream:     "grok-imagine-video-1.5",
			wantBody:         `{"model":"grok-imagine-video-1.5","prompt":"waves"placeholder`,
			responseBody:     `{"request_id":"video-request-mapped"placeholder`,
	placeholder,
		{
			name:             "image-to-video preserves then maps the requested model",
			endpoint:         GrokMediaEndpointVideosGenerations,
			path:             "/v1/videos/generations",
			body:             `{"model":"grok-imagine-video-1.5","prompt":"animate","image":{"url":"https://example.com/input.png"placeholderplaceholder`,
			modelMapping:     map[string]any{"grok-imagine-video-1.5": "vendor-image-video"placeholder,
			wantRequestModel: "grok-imagine-video-1.5",
			wantUpstream:     "vendor-image-video",
			wantBody:         `{"model":"vendor-image-video","prompt":"animate","image":{"url":"https://example.com/input.png"placeholderplaceholder`,
			responseBody:     `{"request_id":"image-video-request-mapped"placeholder`,
	placeholder,
		{
			name:             "mapping and image sanitization compose",
			endpoint:         GrokMediaEndpointImagesGenerations,
			path:             "/v1/images/generations",
			body:             `{"model":"grok-imagine","prompt":"draw","size":"1024x1024"placeholder`,
			modelMapping:     map[string]any{"grok-imagine-image-quality": "vendor-image-model"placeholder,
			wantRequestModel: "grok-imagine-image-quality",
			wantUpstream:     "vendor-image-model",
			wantBody:         `{"model":"vendor-image-model","prompt":"draw"placeholder`,
			responseBody:     `{"data":[{"url":"https://images.test/mapped.png"placeholder]placeholder`,
	placeholder,
		{
			name:             "whitespace mapping target safely preserves normalized model",
			endpoint:         GrokMediaEndpointImagesGenerations,
			path:             "/v1/images/generations",
			body:             `{"model":"grok-imagine","prompt":"draw"placeholder`,
			modelMapping:     map[string]any{"grok-imagine-image-quality": "   "placeholder,
			wantRequestModel: "grok-imagine-image-quality",
			wantUpstream:     "grok-imagine-image-quality",
			wantBody:         `{"model":"grok-imagine-image-quality","prompt":"draw"placeholder`,
			responseBody:     `{"data":[{"url":"https://images.test/mapped.png"placeholder]placeholder`,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request = httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")

			account := &Account{
				ID:          66,
				Name:        "grok-mapped",
				Platform:    PlatformGrok,
				Type:        AccountTypeAPIKey,
				Concurrency: 1,
		placeholder
					"api_key":       "api-key",
					"base_url":      "https://xai.test/v1",
					"model_mapping": tt.modelMapping,
			placeholder,
		placeholder
			upstream := &httpUpstreamRecorder{resp: &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
				Body:       io.NopCloser(strings.NewReader(tt.responseBody)),
		placeholderplaceholder
			svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder

			result, err := svc.ForwardGrokMedia(context.Background(), c, account, tt.endpoint, "", []byte(tt.body), "application/json")

		placeholder
			require.JSONEq(t, tt.wantBody, string(upstream.lastBody))
			require.Equal(t, tt.wantRequestModel, result.Model)
			require.Equal(t, tt.wantRequestModel, result.BillingModel)
			require.Equal(t, tt.wantUpstream, result.UpstreamModel)
	placeholder)
placeholder
placeholder

func TestForwardGrokMediaImagesGenerationRejectsEmptySuccessfulResponse(t *testing.T) {
	t.Setenv(xai.EnvAllowUnsafeURLOverrides, "true")
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok-imagine-image","prompt":"draw a cat"placeholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	account := &Account{
		ID:          66,
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
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"data":[]placeholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder

	result, err := svc.ForwardGrokMedia(context.Background(), c, account, GrokMediaEndpointImagesGenerations, "", body, "application/json")
	require.Nil(t, result)
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusBadGateway, failoverErr.StatusCode)
	require.JSONEq(t, `{"data":[]placeholder`, string(failoverErr.ResponseBody))
	require.Empty(t, recorder.Body.String())
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
			"api_key":       "api-key",
			"base_url":      "https://xai.test/v1",
			"model_mapping": map[string]any{"grok-imagine-edit": "vendor-image-edit"placeholder,
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"application/json"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(`{"data":[{"url":"https://images.test/cat.png"placeholder]placeholder`)),
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
			"api_key":       "api-key",
			"base_url":      "https://xai.test/v1",
			"model_mapping": map[string]any{"grok-imagine-edit": "vendor-image-edit"placeholder,
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"application/json"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(`{"data":[{"url":"https://images.test/edited.png"placeholder]placeholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder

	result, err := svc.ForwardGrokMedia(context.Background(), c, account, GrokMediaEndpointImagesEdits, "", buf.Bytes(), writer.FormDataContentType())
placeholder
	require.Equal(t, "https://xai.test/v1/images/edits", upstream.lastReq.URL.String())
	require.Equal(t, "application/json", upstream.lastReq.Header.Get("Content-Type"))
	require.True(t, json.Valid(upstream.lastBody))
	require.Equal(t, "vendor-image-edit", gjson.GetBytes(upstream.lastBody, "model").String())
	require.Equal(t, "edit this private image", gjson.GetBytes(upstream.lastBody, "prompt").String())
	require.True(t, strings.HasPrefix(gjson.GetBytes(upstream.lastBody, "image.url").String(), "data:image/png;base64,"))
	require.False(t, gjson.GetBytes(upstream.lastBody, "image.image_url").Exists())
	require.Equal(t, "grok-imagine-edit", result.BillingModel)
	require.Equal(t, "vendor-image-edit", result.UpstreamModel)
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
	require.JSONEq(t, `{"model":"grok-imagine-video-1.5","prompt":"waves","resolution":"720p","duration":10placeholder`, string(upstream.lastBody))
	require.Equal(t, "video-request-123", result.ResponseID)
	require.Equal(t, "grok-imagine-video-1.5", result.BillingModel)
	require.Equal(t, 3, result.Usage.InputTokens)
	require.Equal(t, 4, result.Usage.OutputTokens)
	// Create accepts the job only — VideoCount stays 0 until status returns video.url.
	require.Equal(t, 0, result.ImageCount)
	require.Empty(t, result.ImageSize)
	require.Equal(t, 0, result.VideoCount)
	require.Equal(t, VideoBillingResolution720P, result.VideoResolution)
	require.Equal(t, 10, result.VideoDurationSeconds)
placeholder

func TestForwardGrokMediaVideoGenerationReturnsTaskIDAsResponseID(t *testing.T) {
	t.Setenv(xai.EnvAllowUnsafeURLOverrides, "true")
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok-imagine-video","prompt":"waves"placeholder`)
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
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"task_id":"video-task-123"placeholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder

	result, err := svc.ForwardGrokMedia(context.Background(), c, account, GrokMediaEndpointVideosGenerations, "", body, "application/json")
placeholder
	require.Equal(t, "video-task-123", result.ResponseID)
placeholder

func TestExtractGrokMediaVideoRequestIDPreservesExistingPrecedence(t *testing.T) {
	body := []byte(`{
		"request_id":"request-id",
		"id":"id",
		"task_id":"task-id",
		"data":{"request_id":"data-request-id","id":"data-id","task_id":"data-task-id"placeholder,
		"video":{"request_id":"video-request-id","id":"video-id","task_id":"video-task-id"placeholder
placeholder`)

	require.Equal(t, "request-id", extractGrokMediaVideoRequestID(body))
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
	require.JSONEq(t, `{"model":"grok-imagine-video-1.5","prompt":"animate","image":{"url":"data:image/png;base64,aW1n"placeholderplaceholder`, string(upstream.lastBody))
	require.Equal(t, "video-request-456", result.ResponseID)
	require.Equal(t, "grok-imagine-video-1.5", result.BillingModel)
	// 未指定 duration 时按上游默认 8 秒计费。
	require.Equal(t, VideoBillingDefaultDurationSeconds, result.VideoDurationSeconds)
placeholder

func TestForwardGrokMediaOAuthImageToVideoUsesOfficialAPIForLargeBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	imageData := strings.Repeat("A", 2*1024*1024)
	body := []byte(`{"model":"grok-imagine-video-1.5","prompt":"animate","image":{"image_url":"data:image/png;base64,` + imageData + `"placeholderplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/videos/generations", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	account := &Account{
		ID:          66,
		Name:        "grok-oauth",
		Platform:    PlatformGrok,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
placeholder
			"access_token":  "oauth-access-token",
			"refresh_token": "oauth-refresh-token",
			"expires_at":    time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339),
			"base_url":      xai.DefaultCLIBaseURL,
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"application/json"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(`{"request_id":"video-request-oauth"placeholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{httpUpstream: upstream, grokTokenProvider: NewGrokTokenProvider(nil, nil)placeholder

	_, err := svc.ForwardGrokMedia(context.Background(), c, account, GrokMediaEndpointVideosGenerations, "", body, "application/json")
placeholder
	require.Equal(t, xai.DefaultBaseURL+"/videos/generations", upstream.lastReq.URL.String())
	require.Empty(t, upstream.lastReq.Header.Get("X-XAI-Token-Auth"))
	require.Empty(t, upstream.lastReq.Header.Get("x-grok-client-version"))
	require.Equal(t, "data:image/png;base64,"+imageData, gjson.GetBytes(upstream.lastBody, "image.url").String())
	require.False(t, gjson.GetBytes(upstream.lastBody, "image.image_url").Exists())
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
	require.Empty(t, upstream.lastReq.Header.Get("X-Grok-Client-Version"))
	require.NotEqual(t, grokUpstreamUserAgent, upstream.lastReq.Header.Get("User-Agent"))
	require.Empty(t, upstream.lastReq.Header.Get("Content-Type"))
	require.Empty(t, upstream.lastBody)
	require.Equal(t, http.StatusOK, recorder.Code)
	require.JSONEq(t, `{"id":"request-123","status":"completed"placeholder`, recorder.Body.String())
	require.Equal(t, "xai-video-req", result.RequestID)
placeholder

func TestForwardGrokMediaVideoMutationEndpoints(t *testing.T) {
	t.Setenv(xai.EnvAllowUnsafeURLOverrides, "true")
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		endpoint GrokMediaEndpoint
		path     string
placeholder{
		{name: "edit", endpoint: GrokMediaEndpointVideosEdits, path: "/videos/edits"placeholder,
		{name: "extension", endpoint: GrokMediaEndpointVideosExtensions, path: "/videos/extensions"placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			body := []byte(`{"model":"grok-imagine-video","prompt":"continue","video":{"url":"https://example.com/in.mp4"placeholder,"duration":6placeholder`)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1"+tt.path, bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			account := &Account{
				ID: 71, Name: "grok", Platform: PlatformGrok, Type: AccountTypeAPIKey, Concurrency: 1,
		placeholder
					"api_key":       "api-key",
					"base_url":      "https://xai.test/v1",
					"model_mapping": map[string]any{"grok-imagine-video": "vendor-video-mutation"placeholder,
			placeholder,
		placeholder
			upstream := &httpUpstreamRecorder{resp: &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
				Body:       io.NopCloser(strings.NewReader(`{"request_id":"video-mutation-123"placeholder`)),
		placeholderplaceholder
			svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder

			result, err := svc.ForwardGrokMedia(context.Background(), c, account, tt.endpoint, "", body, "application/json")
		placeholder
			require.Equal(t, "https://xai.test/v1"+tt.path, upstream.lastReq.URL.String())
			require.Equal(t, http.MethodPost, upstream.lastReq.Method)
			require.JSONEq(t, `{"model":"vendor-video-mutation","prompt":"continue","video":{"url":"https://example.com/in.mp4"placeholder,"duration":6placeholder`, string(upstream.lastBody))
			require.Equal(t, "video-mutation-123", result.ResponseID)
			require.Equal(t, 0, result.VideoCount)
			require.Equal(t, 6, result.VideoDurationSeconds)
			require.Equal(t, "grok-imagine-video", result.BillingModel)
			require.Equal(t, "vendor-video-mutation", result.UpstreamModel)
	placeholder)
placeholder
placeholder

func TestGrokMediaVideoRequestBindingIsScopedToUserAndAPIKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/videos/video-request-123", nil)
	c.Request.Header.Set("session_id", "shared-client-session")
	groupID := int64(7)
	cache := &stubGatewayCache{placeholder
	svc := &OpenAIGatewayService{cache: cacheplaceholder
	const userID int64 = 41
	const apiKeyID int64 = 51
	require.NotEmpty(t, svc.GenerateExplicitSessionHash(c, nil))
	ctx := c.Request.Context()

	hash := GrokMediaVideoRequestSessionHash("video-request-123", userID, apiKeyID)
	require.NotEmpty(t, hash)
	require.NoError(t, svc.BindGrokMediaVideoRequestAccount(ctx, &groupID, "video-request-123", userID, apiKeyID, 63))

	accountID, err := svc.ResolveGrokMediaVideoRequestAccount(ctx, &groupID, "video-request-123", userID, apiKeyID)
placeholder
	require.Equal(t, int64(63), accountID)

	accountID, err = svc.ResolveGrokMediaVideoRequestAccount(ctx, &groupID, "video-request-123", userID+1, apiKeyID)
placeholder
	require.Zero(t, accountID)

	accountID, err = svc.ResolveGrokMediaVideoRequestAccount(ctx, &groupID, "video-request-123", userID, apiKeyID+1)
placeholder
	require.Zero(t, accountID)
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

func TestGrokMedia429FailoverPreservesRetryAfter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", nil)
	account := &Account{
		ID: 641, Name: "grok-oauth", Platform: PlatformGrok, Type: AccountTypeOAuth,
		Status: StatusActive, Schedulable: true,
placeholder
			"custom_error_codes_enabled": true,
			"custom_error_codes":         []any{float64(http.StatusTooManyRequests)placeholder,
	placeholder,
placeholder
	repo := &grokQuotaAccountRepo{placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{"Retry-After": []string{"45"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"rate limited"placeholderplaceholder`)),
placeholder

	result, err := svc.handleGrokMediaErrorResponse(context.Background(), resp, c, account, "request-id", "grok-imagine")

	require.Nil(t, result)
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusTooManyRequests, failoverErr.StatusCode)
	require.Equal(t, "45", failoverErr.ResponseHeaders.Get("Retry-After"))
placeholder

func healthyGrokOAuthGatewayTestAccount(id int64, token string) *Account {
placeholder
		ID:          id,
		Name:        "grok",
		Platform:    PlatformGrok,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
placeholder
			"access_token":  token,
			"refresh_token": "refresh-token",
			"expires_at":    time.Now().Add(2 * grokTokenRefreshSkew).UTC().Format(time.RFC3339),
			"base_url":      xai.DefaultCLIBaseURL,
	placeholder,
placeholder
placeholder

func TestForwardAsChatCompletionsForGrokStopFallsBackToXAIChatCompletions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok","messages":[{"role":"user","content":"hi"placeholder],"stream":false,"stop":"done","prompt_cache_key":"raw-client-cache-key"placeholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	c.Set("api_key", &APIKey{ID: 5101placeholder)

	account := healthyGrokOAuthGatewayTestAccount(51, "access-token")
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
	require.Equal(t, "grok-4.6", gjson.GetBytes(upstream.lastBody, "model").String())
	require.False(t, gjson.GetBytes(upstream.lastBody, "prompt_cache_key").Exists())
	require.Equal(t, "grok", result.Model)
	require.Equal(t, "grok-4.6", result.UpstreamModel)
	require.Equal(t, 1, result.Usage.InputTokens)
	require.Equal(t, 2, result.Usage.OutputTokens)
	require.Equal(t, 1, result.Usage.CacheReadInputTokens)
	require.NotNil(t, repo.updates[51][grokQuotaSnapshotExtraKey])
	require.Equal(t, http.StatusOK, recorder.Code)
placeholder

func TestForwardGrokResponsesStreamingDefaultsEmptyModelTo45AndSnapshots(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"input":"hi","stream":true,"reasoning_effort":"high"placeholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("OpenAI-Beta", "responses=experimental")
	c.Set("api_key", &APIKey{ID: 5201placeholder)

	account := healthyGrokOAuthGatewayTestAccount(52, "access-token")
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

	result, err := svc.forwardGrokResponses(context.Background(), c, account, body, "", true, time.Now())
placeholder
	require.Equal(t, xai.DefaultCLIBaseURL+"/responses", upstream.lastReq.URL.String())
	require.Equal(t, "Bearer access-token", upstream.lastReq.Header.Get("Authorization"))
	require.Equal(t, "responses=experimental", upstream.lastReq.Header.Get("OpenAI-Beta"))
	require.Equal(t, "grok-4.6", gjson.GetBytes(upstream.lastBody, "model").String())
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

func TestForwardGrokResponsesAPIKeyUsesXAIResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok","input":"hi","stream":trueplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	account := &Account{
		ID:          53,
		Name:        "grok-api-key",
		Platform:    PlatformGrok,
		Type:        AccountTypeAPIKey,
		Concurrency: 2,
placeholder
			"api_key":  "xai-test-key",
			"base_url": "https://api.x.ai/v1",
	placeholder,
placeholder
	upstreamBody := strings.Join([]string{
		`data: {"type":"response.output_text.delta","sequence_number":0,"delta":"ok"placeholder`,
		"",
		`data: {"type":"response.completed","sequence_number":1,"response":{"id":"resp_grok_api_key","model":"grok-4.5","usage":{"input_tokens":2,"output_tokens":1placeholderplaceholderplaceholder`,
		"",
placeholder, "\n")
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(upstreamBody)),
placeholderplaceholder
	svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder

	result, err := svc.forwardGrokResponses(context.Background(), c, account, body, "grok", true, time.Now())
placeholder
	require.Equal(t, "https://api.x.ai/v1/responses", upstream.lastReq.URL.String())
	require.Equal(t, "Bearer xai-test-key", upstream.lastReq.Header.Get("Authorization"))
	require.Empty(t, upstream.lastReq.Header.Get("X-Grok-Client-Version"))
	require.NotEqual(t, grokUpstreamUserAgent, upstream.lastReq.Header.Get("User-Agent"))
	require.Equal(t, "grok-4.6", gjson.GetBytes(upstream.lastBody, "model").String())
	require.Equal(t, "resp_grok_api_key", result.ResponseID)
	require.Equal(t, 2, result.Usage.InputTokens)
	require.Equal(t, 1, result.Usage.OutputTokens)
placeholder

func TestForwardGrokResponsesRetriesInvalidEncryptedContentOnce(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{
		"model":"grok",
		"input":[
			{"type":"reasoning","summary":[{"type":"summary_text","text":"keep this summary"placeholder],"encrypted_content":"encrypted-reasoning"placeholder,
			{"type":"message","role":"user","content":[{"type":"input_text","text":"hi"placeholder]placeholder
		],
		"metadata":{"large_id":9007199254740993placeholder,
		"stream":false
placeholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))
	c.Set("api_key", &APIKey{ID: 4535placeholder)

	account := &Account{
		ID:          4535,
		Name:        "grok-api-key",
		Platform:    PlatformGrok,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 2,
placeholder
			"api_key":  "same-token",
			"base_url": "https://api.x.ai/v1",
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{responses: []*http.Response{
		{
			StatusCode: http.StatusBadRequest,
			Header: http.Header{
				"Content-Type":   []string{"application/json"placeholder,
				"Xai-Request-Id": []string{"recoverable-first"placeholder,
		placeholder,
			Body: io.NopCloser(strings.NewReader(`{"code":"invalid-argument","error":"Could not decrypt the provided encrypted_content. Ensure the value is unmodified."placeholder`)),
	placeholder,
		{
			StatusCode: http.StatusOK,
			Header: http.Header{
				"Content-Type":   []string{"application/json"placeholder,
				"Xai-Request-Id": []string{"recovered-second"placeholder,
		placeholder,
			Body: io.NopCloser(strings.NewReader(`{"id":"resp_recovered","object":"response","model":"grok-4.5","status":"completed","output":[],"usage":{"input_tokens":2,"output_tokens":1placeholderplaceholder`)),
	placeholder,
placeholderplaceholder
	svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder

	result, err := svc.forwardGrokResponses(context.Background(), c, account, body, "grok", false, time.Now())
placeholder
	require.NotNil(t, result)
	require.Equal(t, "resp_recovered", result.ResponseID)
	require.Equal(t, "recovered-second", result.RequestID)
	require.Len(t, upstream.requests, 2)
	require.Len(t, upstream.bodies, 2)

	require.Equal(t, "reasoning", gjson.GetBytes(upstream.bodies[0], "input.0.type").String())
	require.Equal(t, "encrypted-reasoning", gjson.GetBytes(upstream.bodies[0], "input.0.encrypted_content").String())
	require.Equal(t, "reasoning", gjson.GetBytes(upstream.bodies[1], "input.0.type").String())
	require.False(t, gjson.GetBytes(upstream.bodies[1], "input.0.encrypted_content").Exists())
	require.Equal(t, "keep this summary", gjson.GetBytes(upstream.bodies[1], "input.0.summary.0.text").String())
	require.Equal(t, "message", gjson.GetBytes(upstream.bodies[1], "input.1.type").String())
	require.Equal(t, "9007199254740993", gjson.GetBytes(upstream.bodies[0], "metadata.large_id").Raw)
	require.Equal(t, "9007199254740993", gjson.GetBytes(upstream.bodies[1], "metadata.large_id").Raw)

	firstIdentity := gjson.GetBytes(upstream.bodies[0], "prompt_cache_key").String()
	secondIdentity := gjson.GetBytes(upstream.bodies[1], "prompt_cache_key").String()
	require.NotEmpty(t, firstIdentity)
	require.Equal(t, firstIdentity, secondIdentity)
	for _, req := range upstream.requests {
		require.Equal(t, "Bearer same-token", req.Header.Get("Authorization"))
		require.Equal(t, firstIdentity, req.Header.Get(grokConversationIDHeader))
placeholder
	require.Equal(t, StatusActive, account.Status)
	_, hasUpstreamErrors := c.Get(OpsUpstreamErrorsKey)
	require.False(t, hasUpstreamErrors)
	_, hasTerminalStatus := c.Get(OpsUpstreamStatusCodeKey)
	require.False(t, hasTerminalStatus)
placeholder

func TestForwardGrokResponsesInvalidEncryptedContentRecoveryDoesNotOvermatch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	matchingError := `{"code":"invalid-argument","error":"Could not decrypt the provided encrypted_content."placeholder`
	tests := []struct {
		name         string
		requestBody  string
		responseBody string
placeholder{
		{
			name:         "different top-level code",
			requestBody:  `{"model":"grok","input":[{"type":"reasoning","encrypted_content":"cipher"placeholder],"stream":falseplaceholder`,
			responseBody: `{"code":"bad-request","error":"Could not decrypt the provided encrypted_content."placeholder`,
	placeholder,
		{
			name:         "message does not mention decryption",
			requestBody:  `{"model":"grok","input":[{"type":"reasoning","encrypted_content":"cipher"placeholder],"stream":falseplaceholder`,
			responseBody: `{"code":"invalid-argument","error":"The provided encrypted_content is invalid."placeholder`,
	placeholder,
		{
			name:         "request has no encrypted reasoning",
			requestBody:  `{"model":"grok","input":[{"type":"message","role":"user","content":"hi"placeholder],"stream":falseplaceholder`,
			responseBody: matchingError,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			body := []byte(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))

			account := &Account{
				ID:          4536,
				Name:        "grok-api-key",
				Platform:    PlatformGrok,
				Type:        AccountTypeAPIKey,
				Concurrency: 1,
		placeholder"api_key": "token", "base_url": "https://api.x.ai/v1"placeholder,
		placeholder
			upstream := &httpUpstreamRecorder{responses: []*http.Response{{
				StatusCode: http.StatusBadRequest,
				Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
				Body:       io.NopCloser(strings.NewReader(tt.responseBody)),
		placeholderplaceholderplaceholder
			svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder

			result, err := svc.forwardGrokResponses(context.Background(), c, account, body, "grok", false, time.Now())
			require.Nil(t, result)
		placeholder
			require.Len(t, upstream.requests, 1)
			require.Len(t, upstream.bodies, 1)
	placeholder)
placeholder
placeholder

func TestGrokCompactionBlobRecoveryStripsCompactionItem(t *testing.T) {
	body := []byte(`{"model":"grok","input":[{"type":"compaction","id":"cmp_1","encrypted_content":"blob"placeholder,{"type":"message","role":"user","content":"hi"placeholder]placeholder`)
	require.True(t, isGrokInvalidEncryptedContentResponse(http.StatusUnprocessableEntity, []byte(`{"code":"invalid_compaction","error":"could not decode the compaction blob"placeholder`)))
	retry, changed, err := trimGrokInvalidEncryptedContentRetryBody(body)
placeholder
	require.True(t, changed)
	require.Equal(t, "message", gjson.GetBytes(retry, "input.0.type").String())
	require.False(t, gjson.GetBytes(retry, "input.#(type==\"compaction\")").Exists())
placeholder

func TestForwardGrokResponsesInvalidEncryptedContentRecoveryNestedErrorShape(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok","input":[{"type":"reasoning","encrypted_content":"cipher"placeholder,{"type":"message","role":"user","content":"hi"placeholder],"stream":falseplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))

	account := &Account{
		ID:          4538,
		Name:        "grok-api-key",
		Platform:    PlatformGrok,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder"api_key": "token", "base_url": "https://api.x.ai/v1"placeholder,
placeholder
	upstream := &httpUpstreamRecorder{responses: []*http.Response{
		{
			StatusCode: http.StatusBadRequest,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"code":"invalid-argument","error":{"message":"Could not decrypt the provided encrypted_content."placeholderplaceholder`)),
	placeholder,
		{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(`{"id":"resp_ok","output":[],"usage":{"input_tokens":1,"output_tokens":1placeholderplaceholder`)),
	placeholder,
placeholderplaceholder
	svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder

	result, err := svc.forwardGrokResponses(context.Background(), c, account, body, "grok", false, time.Now())
placeholder
	require.NotNil(t, result)
	require.Len(t, upstream.requests, 2)
	require.True(t, gjson.GetBytes(upstream.bodies[0], "input.0.encrypted_content").Exists())
	require.False(t, gjson.GetBytes(upstream.bodies[1], "input.0.encrypted_content").Exists())
placeholder

func TestForwardGrokResponsesInvalidEncryptedContentRetryFailureIsTerminal(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok","input":[{"type":"reasoning","encrypted_content":"cipher"placeholder,{"type":"message","role":"user","content":"hi"placeholder],"stream":falseplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))

	account := &Account{
		ID:          4537,
		Name:        "grok-api-key",
		Platform:    PlatformGrok,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder"api_key": "same-token", "base_url": "https://api.x.ai/v1"placeholder,
placeholder
	newInvalidEncryptedResponse := func(requestID string) *http.Response {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Header: http.Header{
				"Content-Type":   []string{"application/json"placeholder,
				"Xai-Request-Id": []string{requestIDplaceholder,
		placeholder,
			Body: io.NopCloser(strings.NewReader(`{"code":"invalid-argument","error":"Could not decrypt the provided encrypted_content."placeholder`)),
	placeholder
placeholder
	upstream := &httpUpstreamRecorder{responses: []*http.Response{
		newInvalidEncryptedResponse("recoverable-first"),
		newInvalidEncryptedResponse("terminal-second"),
placeholderplaceholder
	svc := &OpenAIGatewayService{httpUpstream: upstreamplaceholder

	result, err := svc.forwardGrokResponses(context.Background(), c, account, body, "grok", false, time.Now())
	require.Nil(t, result)
placeholder
	require.Len(t, upstream.requests, 2)
	require.Len(t, upstream.bodies, 2)
	require.True(t, gjson.GetBytes(upstream.bodies[0], "input.0.encrypted_content").Exists())
	require.False(t, gjson.GetBytes(upstream.bodies[1], `input.#(type=="reasoning")`).Exists())

	rawEvents, ok := c.Get(OpsUpstreamErrorsKey)
	require.True(t, ok)
	events, ok := rawEvents.([]*OpsUpstreamErrorEvent)
	require.True(t, ok)
	require.NotEmpty(t, events)
	for _, event := range events {
		require.NotEqual(t, "recoverable-first", event.UpstreamRequestID)
placeholder
	require.Equal(t, http.StatusBadRequest, c.GetInt(OpsUpstreamStatusCodeKey))
placeholder

func TestForwardAsChatCompletionsForGrokAPIKeyUsesConfiguredRawEndpointWithoutOAuthIdentity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok","messages":[{"role":"user","content":"hi"placeholder],"stream":falseplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	account := &Account{
		ID:          706,
		Platform:    PlatformGrok,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "third-party-key",
			"base_url": "https://grok.example.test/v1",
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"id":"chatcmpl","object":"chat.completion","model":"grok-4.5","choices":[],"usage":{"prompt_tokens":1,"completion_tokens":1placeholderplaceholder`)),
placeholderplaceholder
	svc := &OpenAIGatewayService{cfg: rawChatCompletionsTestConfig(), httpUpstream: upstreamplaceholder

	_, err := svc.ForwardAsChatCompletions(context.Background(), c, account, body, "", "")
placeholder
	require.Equal(t, "https://grok.example.test/v1/chat/completions", upstream.lastReq.URL.String())
	require.Equal(t, "Bearer third-party-key", upstream.lastReq.Header.Get("Authorization"))
	require.Empty(t, upstream.lastReq.Header.Get("X-Grok-Client-Version"))
	require.NotEqual(t, grokUpstreamUserAgent, upstream.lastReq.Header.Get("User-Agent"))
placeholder

func TestForwardAsChatCompletionsForGrokAPIKeyRejectsNonStreamingResponseWithoutUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok","messages":[{"role":"user","content":"hi"placeholder],"stream":falseplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	account := &Account{
		ID:          707,
		Platform:    PlatformGrok,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "third-party-key",
			"base_url": "https://grok.example.test/v1",
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(
			`{"id":"resp_missing_usage","object":"chat.completion","model":"grok-4.5","choices":[{"index":0,"message":{"role":"assistant","content":"ok"placeholder,"finish_reason":"stop"placeholder]placeholder`,
		)),
placeholderplaceholder
	svc := &OpenAIGatewayService{cfg: rawChatCompletionsTestConfig(), httpUpstream: upstreamplaceholder

	result, err := svc.ForwardAsChatCompletions(context.Background(), c, account, body, "", "")

	require.Nil(t, result)
	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, http.StatusBadGateway, failoverErr.StatusCode)
	require.Equal(t, grokMissingUsageErrorCode, gjson.GetBytes(failoverErr.ResponseBody, "error.code").String())
	require.False(t, c.Writer.Written(), "an unbillable response must not be committed to the client")
	require.Empty(t, recorder.Body.String())
placeholder

func TestAccountTestServiceGrokAPIKeyUsesXAIResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	account := &Account{
		ID:          54,
		Name:        "grok-api-key",
		Platform:    PlatformGrok,
		Type:        AccountTypeAPIKey,
		Concurrency: 2,
placeholder
			"api_key":  "xai-test-key",
			"base_url": "https://api.x.ai/v1",
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(
			"data: {\"type\":\"response.output_text.delta\",\"delta\":\"ok\"placeholder\n\n" +
				"data: {\"type\":\"response.completed\"placeholder\n\n",
		)),
placeholderplaceholder
	svc := &AccountTestService{httpUpstream: upstreamplaceholder
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/accounts/54/test", nil)

	err := svc.testGrokAccountConnection(c, account, "grok", "", AccountTestModeDefault, AccountTestOptions{placeholder)
placeholder
	require.Equal(t, "https://api.x.ai/v1/responses", upstream.lastReq.URL.String())
	require.Equal(t, "Bearer xai-test-key", upstream.lastReq.Header.Get("Authorization"))
	require.Contains(t, recorder.Body.String(), `"type":"test_complete"`)
placeholder

func TestAccountTestServiceGrokAPIKeyAllowsConfiguredHTTPWhenGlobalPolicyDoes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	account := &Account{
		ID:          55,
		Name:        "grok-api-key-http",
		Platform:    PlatformGrok,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
placeholder
			"api_key":  "third-party-key",
			"base_url": "http://grok.example.test/v1",
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body: io.NopCloser(strings.NewReader(
			"data: {\"type\":\"response.output_text.delta\",\"delta\":\"ok\"placeholder\n\n" +
				"data: {\"type\":\"response.completed\"placeholder\n\n",
		)),
placeholderplaceholder
	svc := &AccountTestService{cfg: rawChatCompletionsTestConfig(), httpUpstream: upstreamplaceholder
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/accounts/55/test", nil)

	err := svc.testGrokAccountConnection(c, account, "grok", "", AccountTestModeDefault, AccountTestOptions{placeholder)
placeholder
	require.Equal(t, "http://grok.example.test/v1/responses", upstream.lastReq.URL.String())
	require.Equal(t, "Bearer third-party-key", upstream.lastReq.Header.Get("Authorization"))
	require.Empty(t, upstream.lastReq.Header.Get("X-Grok-Client-Version"))
	require.Contains(t, recorder.Body.String(), `"type":"test_complete"`)
placeholder

func TestAccountTestServiceGrokOAuthPaymentRequiredTemporarilyUnschedulesAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	account := healthyGrokOAuthGatewayTestAccount(56, "access-token")
	repo := &grokQuotaAccountRepo{placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusPaymentRequired,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"code":"personal-team-blocked:spending-limit"placeholder`)),
placeholderplaceholder
	svc := &AccountTestService{
		accountRepo:       repo,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		httpUpstream:      upstream,
placeholder
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/accounts/56/test", nil)
	before := time.Now()

	err := svc.testGrokAccountConnection(c, account, "grok", "", AccountTestModeDefault, AccountTestOptions{placeholder)

placeholder
	require.Zero(t, repo.tempUnschedCalls)
	require.Equal(t, 1, repo.rateLimitedCalls)
	require.Equal(t, account.ID, repo.lastRateLimitedID)
	require.WithinDuration(t, before.Add(grokSpendingLimitProbeCooldown), repo.lastRateLimitResetAt, time.Second)
	require.Contains(t, recorder.Body.String(), `"type":"error"`)
	require.Contains(t, recorder.Body.String(), "Grok Responses API returned 402")
placeholder

func TestForwardAsChatCompletionsForGrokStreamingUsesRawXAIChatCompletions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok","messages":[{"role":"user","content":"hi"placeholder],"stream":trueplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	account := healthyGrokOAuthGatewayTestAccount(53, "access-token")
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
	require.Equal(t, xai.CLIUserAgent(xai.CLIClientVersion), upstream.lastReq.Header.Get("User-Agent"))
	require.Equal(t, "grok-4.6", gjson.GetBytes(upstream.lastBody, "model").String())
	require.True(t, gjson.GetBytes(upstream.lastBody, "stream_options.include_usage").Bool())
	require.True(t, result.Stream)
	require.Equal(t, 6, result.Usage.InputTokens)
	require.Equal(t, 4, result.Usage.OutputTokens)
	require.Equal(t, 1, result.Usage.CacheReadInputTokens)
	require.Contains(t, recorder.Body.String(), "data: [DONE]")
	require.NotNil(t, repo.updates[53][grokQuotaSnapshotExtraKey])
placeholder

func TestForwardGrokResponsesNonStreamingUsesCacheIdentityAndCachedUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok","input":"hi","stream":false,"tools":[{"type":"namespace","name":"client_tools"placeholder],"tool_choice":{"type":"namespace","name":"client_tools"placeholderplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("api_key", &APIKey{ID: 5202placeholder)

	account := healthyGrokOAuthGatewayTestAccount(56, "access-token")
	observedResetAt := time.Now().Add(-time.Second).UTC().Truncate(time.Second)
	observedLimitedAt := observedResetAt.Add(-grokRateLimitRepeatCooldown)
	account.RateLimitedAt = &observedLimitedAt
	account.RateLimitResetAt = &observedResetAt
	repo := &grokQuotaAccountRepo{
		mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
			accountsByID: map[int64]*Account{56: accountplaceholder,
	placeholder,
		recoveryClearResult: true,
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
	require.Equal(t, "grok-4.6", gjson.GetBytes(upstream.lastBody, "model").String())
	identity := gjson.GetBytes(upstream.lastBody, "prompt_cache_key").String()
	require.NotEmpty(t, identity)
	require.Equal(t, identity, upstream.lastReq.Header.Get(grokConversationIDHeader))
	// The sanitizer drops this unsupported client tool, but its explicit intent
	// must still prevent native cache-routing tools from being injected.
	require.False(t, gjson.GetBytes(upstream.lastBody, "tools").Exists())
	require.False(t, gjson.GetBytes(upstream.lastBody, "tool_choice").Exists())
	require.Equal(t, "resp_grok_non_stream", gjson.Get(recorder.Body.String(), "id").String())
	require.Equal(t, 1, repo.recoveryClearCalls)
	require.Equal(t, observedLimitedAt, repo.recoveryObservedAt)
	require.Equal(t, observedResetAt, repo.recoveryObservedReset)
placeholder

func TestForwardGrokResponsesFailoverKeepsCacheIdentityAcrossAccounts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok","input":[{"role":"user","content":"stable prefix"placeholder],"stream":falseplaceholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))
	c.Set("api_key", &APIKey{ID: 5203placeholder)

	newAccount := func(id int64, token string) *Account {
		account := healthyGrokOAuthGatewayTestAccount(id, token)
		account.Name = fmt.Sprintf("grok-%d", id)
		return account
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

	account := healthyGrokOAuthGatewayTestAccount(53, "access-token")
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
	require.Equal(t, xai.CLIUserAgent(xai.CLIClientVersion), upstream.lastReq.Header.Get("User-Agent"))
	require.Equal(t, grokCLIVersion, upstream.lastReq.Header.Get("X-Grok-Client-Version"))
	require.NotEmpty(t, upstream.lastReq.Header.Get(grokConversationIDHeader))
	require.NotEqual(t, "native-client-conversation", upstream.lastReq.Header.Get(grokConversationIDHeader))
	require.Equal(t, "grok-4.6", gjson.GetBytes(upstream.lastBody, "model").String())
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

	account := healthyGrokOAuthGatewayTestAccount(55, "access-token")
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

	account := healthyGrokOAuthGatewayTestAccount(54, "access-token")
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
	require.Equal(t, xai.CLIUserAgent(xai.CLIClientVersion), upstream.lastReq.Header.Get("User-Agent"))
	require.Equal(t, grokCLIVersion, upstream.lastReq.Header.Get("X-Grok-Client-Version"))
	require.Equal(t, "grok-experimental", upstream.lastReq.Header.Get("OpenAI-Beta"))
	require.Empty(t, upstream.lastReq.Header.Get("originator"))
	require.Empty(t, upstream.lastReq.Header.Get("version"))
	require.Equal(t, "grok-4.6", gjson.GetBytes(upstream.lastBody, "model").String())
	require.NotEmpty(t, gjson.GetBytes(upstream.lastBody, "prompt_cache_key").String())
	require.Equal(t, gjson.GetBytes(upstream.lastBody, "prompt_cache_key").String(), upstream.lastReq.Header.Get(grokConversationIDHeader))
	require.Equal(t, "web_search", gjson.GetBytes(upstream.lastBody, "tools.0.type").String())
	require.Equal(t, "x_search", gjson.GetBytes(upstream.lastBody, "tools.1.type").String())
	require.Equal(t, "none", gjson.GetBytes(upstream.lastBody, "tool_choice").String())
	require.Empty(t, upstream.lastReq.Header.Get("session_id"))
	require.True(t, gjson.GetBytes(upstream.lastBody, "stream").Bool())
	require.NotContains(t, string(upstream.lastBody), "chatgpt.com")
	require.Equal(t, "grok", result.Model)
	require.Equal(t, "grok-4.6", result.UpstreamModel)
	require.Equal(t, 5, result.Usage.InputTokens)
	require.Equal(t, 2, result.Usage.OutputTokens)
	require.Equal(t, 3, result.Usage.CacheReadInputTokens)
	require.Contains(t, recorder.Body.String(), `"type":"message"`)
	require.Equal(t, int64(3), gjson.Get(recorder.Body.String(), "usage.cache_read_input_tokens").Int())
	require.Contains(t, recorder.Body.String(), "ok")
placeholder

func TestForwardAsAnthropicForGrokFunctionToolUsesCacheCapableMixedRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{
		"model":"grok","max_tokens":32,"stream":false,
		"messages":[{"role":"user","content":"look up alpha"placeholder],
		"tools":[{"name":"lookup","description":"look up a key","input_schema":{"type":"object","properties":{"key":{"type":"string"placeholderplaceholder,"required":["key"]placeholderplaceholder,{"name":"web_search","description":"search the web","input_schema":{"type":"object","properties":{"query":{"type":"string"placeholderplaceholder,"required":["query"]placeholderplaceholder],
		"tool_choice":{"type":"auto"placeholder
placeholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	c.Set("api_key", &APIKey{ID: 5403placeholder)

	account := healthyGrokOAuthGatewayTestAccount(58, "access-token")
	account.Extra = map[string]any{grokBillingExtraKey: map[string]any{
		"status_code":        http.StatusOK,
		"source":             "billing_probe",
		"monthly_updated_at": "2026-07-15T05:00:00Z",
placeholderplaceholder
	repo := &grokQuotaAccountRepo{
		mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
			accountsByID: map[int64]*Account{58: accountplaceholder,
	placeholder,
placeholder
	responseBody := strings.Join([]string{
		`data: {"type":"response.completed","response":{"id":"resp_grok_function","object":"response","model":"grok-4.5","status":"completed","output":[{"type":"function_call","id":"fc_lookup","call_id":"call_lookup","name":"lookup","arguments":"{\"key\":\"alpha\"placeholder","status":"completed"placeholder],"usage":{"input_tokens":7000,"output_tokens":2,"total_tokens":7002,"input_tokens_details":{"cached_tokens":6144placeholderplaceholderplaceholderplaceholder`,
		"",
		"data: [DONE]",
		"",
placeholder, "\n")
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(responseBody)),
placeholderplaceholder
	svc := &OpenAIGatewayService{
		httpUpstream:      upstream,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		accountRepo:       repo,
placeholder

	result, err := svc.ForwardAsAnthropic(context.Background(), c, account, body, "", "")

placeholder
	require.NotNil(t, result)
	require.Equal(t, xai.DefaultCLIBaseURL+"/responses", upstream.lastReq.URL.String())
	identity := gjson.GetBytes(upstream.lastBody, "prompt_cache_key").String()
	require.NotEmpty(t, identity)
	require.Equal(t, identity, upstream.lastReq.Header.Get(grokConversationIDHeader))
	tools := gjson.GetBytes(upstream.lastBody, "tools").Array()
	require.Len(t, tools, 3)
	require.Equal(t, "function", tools[0].Get("type").String())
	require.Equal(t, "lookup", tools[0].Get("name").String())
	require.Equal(t, "object", tools[0].Get("parameters.type").String())
	require.Equal(t, "web_search", tools[1].Get("type").String())
	require.Equal(t, "x_search", tools[2].Get("type").String())
	require.Equal(t, "auto", gjson.GetBytes(upstream.lastBody, "tool_choice").String())

	require.Equal(t, 7000, result.Usage.InputTokens)
	require.Equal(t, 6144, result.Usage.CacheReadInputTokens)
	clientBody := recorder.Body.String()
	require.Equal(t, "tool_use", gjson.Get(clientBody, "content.0.type").String())
	require.Equal(t, "call_lookup", gjson.Get(clientBody, "content.0.id").String())
	require.Equal(t, "lookup", gjson.Get(clientBody, "content.0.name").String())
	require.Equal(t, "alpha", gjson.Get(clientBody, "content.0.input.key").String())
	require.Equal(t, "tool_use", gjson.Get(clientBody, "stop_reason").String())
	require.Equal(t, int64(856), gjson.Get(clientBody, "usage.input_tokens").Int())
	require.Equal(t, int64(6144), gjson.Get(clientBody, "usage.cache_read_input_tokens").Int())
placeholder

func TestForwardAsAnthropicForGrokStreamingPreservesCacheUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"model":"grok","max_tokens":32,"stream":true,"messages":[{"role":"user","content":"hi"placeholder]placeholder`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
	c.Set("api_key", &APIKey{ID: 5402placeholder)

	account := healthyGrokOAuthGatewayTestAccount(57, "access-token")
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

func TestHandleGrokAccountUpstreamErrorSpendingLimitUsesRecoverableProbeCool(t *testing.T) {
	repo := &grokQuotaAccountRepo{placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder
	account := &Account{ID: 2570, Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder
	before := time.Now()
	body := []byte(`{"code":"personal-team-blocked:spending-limit","error":"You have run out of credits"placeholder`)

	svc.handleGrokAccountUpstreamError(context.Background(), account, http.StatusForbidden, nil, body)

	require.Equal(t, 1, repo.rateLimitedCalls)
	require.WithinDuration(t, before.Add(grokSpendingLimitProbeCooldown), repo.lastRateLimitResetAt, 2*time.Second)
	require.Zero(t, repo.tempUnschedCalls)
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
			wantReason:      "grok credentials unauthorized",
			wantMinCooldown: 10*time.Minute - time.Second,
			wantMaxCooldown: 10*time.Minute + time.Second,
	placeholder,
		{
			name:            "forbidden entitlement",
			status:          http.StatusForbidden,
			wantReason:      "grok access or entitlement denied",
			wantMinCooldown: 30*time.Minute - time.Second,
			wantMaxCooldown: 30*time.Minute + time.Second,
	placeholder,
		{
			name:            "payment required",
			status:          http.StatusPaymentRequired,
			wantReason:      "grok payment required",
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

func TestHandleGrokAccountUpstreamErrorSpendingLimit403RateLimits(t *testing.T) {
	account := &Account{ID: 614, Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder
	repo := &grokQuotaAccountRepo{placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder
	before := time.Now()
	body := []byte(`{"code":"personal-team-blocked:spending-limit","error":"You have run out of credits"placeholder`)

	svc.handleGrokAccountUpstreamError(context.Background(), account, http.StatusForbidden, nil, body)

	require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
	require.Equal(t, 1, repo.rateLimitedCalls)
	require.Equal(t, account.ID, repo.lastRateLimitedID)
	require.WithinDuration(t, before.Add(grokSpendingLimitProbeCooldown), repo.lastRateLimitResetAt, 2*time.Second)
	require.Zero(t, repo.tempUnschedCalls)
	require.True(t, isGrokSpendingLimitError(body))
placeholder

func TestHandleGrokAccountUpstreamError5xxRespectsPoolMode(t *testing.T) {
	t.Run("pool mode keeps scheduling state", func(t *testing.T) {
		account := &Account{
			ID:       611,
			Platform: PlatformGrok,
			Type:     AccountTypeAPIKey,
	placeholder
				"pool_mode": true,
		placeholder,
	placeholder
		repo := &grokQuotaAccountRepo{placeholder
		svc := &OpenAIGatewayService{accountRepo: repoplaceholder

		svc.handleGrokAccountUpstreamError(context.Background(), account, http.StatusBadGateway, nil, nil)

		require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
		require.Zero(t, repo.tempUnschedCalls)
		require.Nil(t, account.TempUnschedulableUntil)
		require.Empty(t, account.TempUnschedulableReason)
placeholder)

	t.Run("non-pool mode keeps two minute cooldown", func(t *testing.T) {
		account := &Account{ID: 612, Platform: PlatformGrok, Type: AccountTypeAPIKeyplaceholder
		repo := &grokQuotaAccountRepo{placeholder
		svc := &OpenAIGatewayService{accountRepo: repoplaceholder
		before := time.Now()

		svc.handleGrokAccountUpstreamError(context.Background(), account, http.StatusBadGateway, nil, nil)

		require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
		require.Equal(t, 1, repo.tempUnschedCalls)
		require.Equal(t, account.ID, repo.lastTempUnschedID)
		require.Equal(t, "grok upstream temporary error", repo.lastTempUnschedReason)
		require.WithinDuration(t, before.Add(2*time.Minute), repo.lastTempUnschedUntil, time.Second)
placeholder)
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

func TestHandleGrokAccountUpstreamError429PoolModeKeepsSchedulingState(t *testing.T) {
	account := &Account{
		ID:       613,
		Platform: PlatformGrok,
		Type:     AccountTypeAPIKey,
placeholder
			"pool_mode": true,
	placeholder,
placeholder
	repo := &grokQuotaAccountRepo{placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder

	svc.handleGrokAccountUpstreamError(
		context.Background(), account, http.StatusTooManyRequests,
		http.Header{"Retry-After": []string{"45"placeholderplaceholder, nil,
	)

	require.Equal(t, 1, repo.updateCalls, "pool mode should retain the quota snapshot for observability")
	require.Zero(t, repo.rateLimitedCalls)
	require.Zero(t, repo.tempUnschedCalls)
	require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
	require.Nil(t, account.RateLimitResetAt)
placeholder

func TestHandleGrokAccountUpstreamError402RecoversAfterCooldownExpiry(t *testing.T) {
	account := &Account{
		ID: 610, Platform: PlatformGrok, Type: AccountTypeOAuth,
		Status: StatusActive, Schedulable: true,
placeholder
	repo := &grokQuotaAccountRepo{placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder

	svc.handleGrokAccountUpstreamError(context.Background(), account, http.StatusPaymentRequired, nil, nil)
	require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
	require.Equal(t, 1, repo.tempUnschedCalls)

	expired := time.Now().Add(-time.Second)
	account.TempUnschedulableUntil = &expired
	svc.openaiAccountRuntimeBlockUntil.Store(account.ID, expired)

	require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
	require.True(t, account.IsSchedulable())
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

func TestGrokRateLimitResetAtForAccountEscalatesRepeated429s(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	retryAfter := 45
	snapshot := &xai.QuotaSnapshot{
		StatusCode:        http.StatusTooManyRequests,
		RetryAfterSeconds: &retryAfter,
		UpdatedAt:         now.Format(time.RFC3339),
placeholder
	tests := []struct {
		name             string
		previousCooldown time.Duration
		wantCooldown     time.Duration
placeholder{
		{name: "repeat after short boundary", previousCooldown: 45 * time.Second, wantCooldown: grokRateLimitRepeatCooldownplaceholder,
		{name: "sustained repeat", previousCooldown: grokRateLimitRepeatCooldown, wantCooldown: grokRateLimitSustainedCooldownplaceholder,
		{name: "capped repeat", previousCooldown: grokRateLimitSustainedCooldown, wantCooldown: grokRateLimitMaxAdaptiveCooldownplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			previousReset := now.Add(-time.Second)
			previousLimited := previousReset.Add(-tt.previousCooldown)
			account := &Account{
				ID:               630,
				Platform:         PlatformGrok,
				Type:             AccountTypeOAuth,
				RateLimitedAt:    &previousLimited,
				RateLimitResetAt: &previousReset,
		placeholder

			resetAt, limited := grokRateLimitResetAtForAccount(account, snapshot, now)

			require.True(t, limited)
			require.WithinDuration(t, now.Add(tt.wantCooldown), resetAt, time.Second)
	placeholder)
placeholder
placeholder

func TestGrokRateLimitResetAtForAccountPreservesAuthoritativeAndQuietRecovery(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	retryAfter := 45
	previousReset := now.Add(-grokRateLimitBackoffQuietPeriod - time.Second)
	previousLimited := previousReset.Add(-grokRateLimitSustainedCooldown)
	account := &Account{
		ID:               631,
		Platform:         PlatformGrok,
		Type:             AccountTypeOAuth,
		RateLimitedAt:    &previousLimited,
		RateLimitResetAt: &previousReset,
placeholder
	snapshot := &xai.QuotaSnapshot{
		StatusCode:        http.StatusTooManyRequests,
		RetryAfterSeconds: &retryAfter,
		UpdatedAt:         now.Format(time.RFC3339),
placeholder

	resetAt, limited := grokRateLimitResetAtForAccount(account, snapshot, now)
	require.True(t, limited)
	require.WithinDuration(t, now.Add(45*time.Second), resetAt, time.Second)

	authoritativeReset := now.Add(2 * time.Hour)
	remaining := int64(0)
	snapshot.Requests = &xai.QuotaWindow{Remaining: &remaining, ResetUnix: grokInt64PtrForTest(authoritativeReset.Unix())placeholder
	recentReset := now.Add(-time.Second)
	recentLimited := recentReset.Add(-grokRateLimitSustainedCooldown)
	account.RateLimitResetAt = &recentReset
	account.RateLimitedAt = &recentLimited

	resetAt, limited = grokRateLimitResetAtForAccount(account, snapshot, now)
	require.True(t, limited)
	require.WithinDuration(t, authoritativeReset, resetAt, time.Second)
placeholder

func TestGrokRateLimitResetAtForAccountLeavesAPIKey429PolicyUnchanged(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	retryAfter := 45
	previousReset := now.Add(-time.Second)
	previousLimited := previousReset.Add(-grokRateLimitSustainedCooldown)
	account := &Account{
		ID:               632,
		Platform:         PlatformGrok,
		Type:             AccountTypeAPIKey,
		RateLimitedAt:    &previousLimited,
		RateLimitResetAt: &previousReset,
placeholder
	snapshot := &xai.QuotaSnapshot{
		StatusCode:        http.StatusTooManyRequests,
		RetryAfterSeconds: &retryAfter,
		UpdatedAt:         now.Format(time.RFC3339),
placeholder

	resetAt, limited := grokRateLimitResetAtForAccount(account, snapshot, now)
	require.True(t, limited)
	require.WithinDuration(t, now.Add(45*time.Second), resetAt, time.Second)
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

func TestUpdateGrokUsageFromResponseHeaderlessSuccessClearsObservedCooldown(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	limitedAt := now.Add(-grokRateLimitRepeatCooldown)
	observedResetAt := now.Add(-time.Second)
	account := &Account{
		ID:               660,
		Platform:         PlatformGrok,
		Type:             AccountTypeOAuth,
		RateLimitedAt:    &limitedAt,
		RateLimitResetAt: &observedResetAt,
placeholder
	repo := &grokQuotaAccountRepo{recoveryClearResult: trueplaceholder
	svc := &OpenAIGatewayService{
		accountRepo:           repo,
		codexSnapshotThrottle: newAccountWriteThrottle(time.Hour),
placeholder

	svc.updateGrokUsageFromResponse(context.Background(), account, nil, http.StatusOK)

	require.Zero(t, repo.updateCalls, "headerless success must not overwrite an informative quota snapshot")
	require.Equal(t, 1, repo.recoveryClearCalls)
	require.Equal(t, limitedAt, repo.recoveryObservedAt)
	require.Equal(t, observedResetAt, repo.recoveryObservedReset)
	require.Same(t, &observedResetAt, account.RateLimitResetAt, "shared account snapshots must not be mutated in place")
placeholder

func TestUpdateGrokUsageFromResponseRecoveryRespectsCancellationAndAPIKeyBoundary(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	observedResetAt := now.Add(-time.Second)
	observedLimitedAt := observedResetAt.Add(-grokRateLimitRepeatCooldown)

	t.Run("parent cancellation does not mutate account state", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		account := &Account{
			ID:               661,
			Platform:         PlatformGrok,
			Type:             AccountTypeOAuth,
			RateLimitedAt:    &observedLimitedAt,
			RateLimitResetAt: &observedResetAt,
	placeholder
		repo := &grokQuotaAccountRepo{recoveryClearResult: trueplaceholder
		svc := &OpenAIGatewayService{accountRepo: repoplaceholder

		svc.updateGrokUsageFromResponse(ctx, account, nil, http.StatusOK)

		require.Zero(t, repo.recoveryClearCalls)
placeholder)

	t.Run("API key success does not alter OAuth cooldown state", func(t *testing.T) {
		account := &Account{
			ID:               662,
			Platform:         PlatformGrok,
			Type:             AccountTypeAPIKey,
			RateLimitedAt:    &observedLimitedAt,
			RateLimitResetAt: &observedResetAt,
	placeholder
		repo := &grokQuotaAccountRepo{recoveryClearResult: trueplaceholder
		svc := &OpenAIGatewayService{accountRepo: repoplaceholder

		svc.updateGrokUsageFromResponse(context.Background(), account, nil, http.StatusOK)

		require.Zero(t, repo.recoveryClearCalls)
placeholder)
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

func TestOpenAIWSHTTPBridgeSSEErrorSideEffectsRunOncePerPlatform(t *testing.T) {
	gin.SetMode(gin.TestMode)

	for _, platform := range []string{PlatformOpenAI, PlatformGrokplaceholder {
		t.Run(platform, func(t *testing.T) {
			repo := &grokQuotaAccountRepo{placeholder
			cfg := &config.Config{placeholder
			upstream := &httpUpstreamRecorder{resp: &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body: io.NopCloser(strings.NewReader(
					"data: {\"type\":\"error\",\"error\":{\"type\":\"rate_limit_error\",\"code\":\"rate_limit_exceeded\",\"message\":\"limited\"placeholderplaceholder\n\n",
				)),
		placeholderplaceholder
			svc := &OpenAIGatewayService{
				cfg:          cfg,
				accountRepo:  repo,
				httpUpstream: upstream,
		placeholder
			if platform == PlatformOpenAI {
				svc.rateLimitService = NewRateLimitService(repo, nil, cfg, nil, nil)
		placeholder
			account := &Account{ID: 70, Platform: platform, Type: AccountTypeOAuth, Concurrency: 1placeholder
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request = httptest.NewRequest(http.MethodGet, "/v1/responses", nil)
			payload := []byte(`{"type":"response.create","model":"gpt-5","input":"hi"placeholder`)
			writes := 0

			result, err := svc.proxyOpenAIWSHTTPBridgeTurn(
				context.Background(), c, account, "sk-test", payload, len(payload),
				"gpt-5", "", "", "", "", 1,
				func([]byte) error {
					writes++
					return nil
			placeholder,
			)

			require.Nil(t, result)
			var failoverErr *UpstreamFailoverError
			require.ErrorAs(t, err, &failoverErr)
			require.Equal(t, http.StatusTooManyRequests, failoverErr.StatusCode)
			require.Zero(t, writes)
			require.Equal(t, 1, repo.rateLimitedCalls)
	placeholder)
placeholder
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

func TestPatchGrokResponsesBody_StripsReasoningContentNull(t *testing.T) {
	t.Parallel()

	body := []byte(`{
		"model": "grok-latest",
		"input": [
			{"type":"message","role":"user","content":[{"type":"input_text","text":"hi"placeholder]placeholder,
			{"type":"reasoning","summary":[{"type":"summary_text","text":"thinking..."placeholder],"content":null,"encrypted_content":nullplaceholder,
			{"type":"message","role":"assistant","content":[{"type":"output_text","text":"Hello!"placeholder]placeholder
		]
placeholder`)

	patched, err := patchGrokResponsesBody(body, "grok-4.5")
placeholder
	require.True(t, json.Valid(patched))

	input := gjson.GetBytes(patched, "input")
	require.True(t, input.IsArray())

	items := input.Array()
	require.Len(t, items, 3)

	reasoning := items[1]
	require.Equal(t, "reasoning", reasoning.Get("type").String())
	require.True(t, reasoning.Get("summary").Exists(), "summary should be preserved")
	require.False(t, reasoning.Get("content").Exists(), "content: null should be stripped")
placeholder

func TestPatchGrokResponsesBody_KeepsReasoningContentNonNull(t *testing.T) {
	t.Parallel()

	body := []byte(`{
		"model": "grok-latest",
		"input": [
			{"type":"reasoning","summary":[{"type":"summary_text","text":"ok"placeholder],"content":"real content"placeholder
		]
placeholder`)

	patched, err := patchGrokResponsesBody(body, "grok-4.5")
placeholder

	reasoning := gjson.GetBytes(patched, "input.0")
	require.Equal(t, "real content", reasoning.Get("content").String(), "non-null content must not be stripped")
placeholder

func TestPatchGrokResponsesBody_MultipleReasoningContentNull(t *testing.T) {
	t.Parallel()

	body := []byte(`{
		"model": "grok-latest",
		"input": [
			{"type":"reasoning","summary":[{"type":"summary_text","text":"r1"placeholder],"content":nullplaceholder,
			{"type":"message","role":"user","content":"hi"placeholder,
			{"type":"reasoning","summary":[{"type":"summary_text","text":"r2"placeholder],"content":nullplaceholder
		]
placeholder`)

	patched, err := patchGrokResponsesBody(body, "grok-4.5")
placeholder

	items := gjson.GetBytes(patched, "input").Array()
	require.Len(t, items, 3)

	require.False(t, items[0].Get("content").Exists())
	require.False(t, items[2].Get("content").Exists())
placeholder

func TestIsGrokImageGenerationModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		model string
		want  bool
placeholder{
		{"grok-imagine", trueplaceholder,
		{"grok-imagine-image-quality", trueplaceholder,
		{"grok-imagine-edit", trueplaceholder,
		{"grok-imagine-image-hd", trueplaceholder,
		{"grok-4.5", falseplaceholder,
		{"grok-composer", falseplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			require.Equal(t, tt.want, isGrokImageGenerationModel(tt.model))
	placeholder)
placeholder
placeholder

func TestBuildGrokSchedulerExtraUpdates_FeedsThresholdEvaluator(t *testing.T) {
	int64p := func(v int64) *int64 { return &v placeholder
	resetUnix := time.Now().Add(90 * time.Minute).Unix()
	snapshot := &xai.QuotaSnapshot{
		Requests: &xai.QuotaWindow{Limit: int64p(100), Remaining: int64p(30)placeholder,                         // 70% used
		Tokens:   &xai.QuotaWindow{Limit: int64p(1000), Remaining: int64p(50), ResetUnix: &resetUnixplaceholder, // 95% used (most constrained)
placeholder

	updates := buildGrokSchedulerExtraUpdates(snapshot)
	require.NotNil(t, updates)
	require.InDelta(t, 95.0, updates["grok_sched_utilization"], 0.001, "picks the most-constrained window")
	require.Contains(t, updates, "grok_sched_reset_at")

	// The written extras must actually drive EvaluateAccountSchedulingThreshold
	// (proves the previously-dead read side is now fed).
	account := &Account{Platform: PlatformGrok, Extra: updatesplaceholder
	decision := EvaluateAccountSchedulingThreshold(account, map[string]int{PlatformGrok: 90placeholder, time.Now())
	require.True(t, decision.ShouldPause)
	require.InDelta(t, 95.0, decision.UsedPercent, 0.001)
	require.NotNil(t, decision.Until)
placeholder

func TestBuildGrokSchedulerExtraUpdates_NilWhenNoQuotaWindows(t *testing.T) {
	require.Nil(t, buildGrokSchedulerExtraUpdates(&xai.QuotaSnapshot{placeholder))
	require.Nil(t, buildGrokSchedulerExtraUpdates(nil))
placeholder
