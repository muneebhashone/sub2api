package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestPrepareBedrockRequestBody_BasicFields(t *testing.T) {
	input := `{"model":"claude-opus-4-6","stream":true,"max_tokens":1024,"messages":[{"role":"user","content":"hi"placeholder]placeholder`
	result, err := PrepareBedrockRequestBody([]byte(input), "us.anthropic.claude-opus-4-6-v1", "")
placeholder

	// anthropic_version 应被注入
	assert.Equal(t, "bedrock-2023-05-31", gjson.GetBytes(result, "anthropic_version").String())
	// model 和 stream 应被移除
	assert.False(t, gjson.GetBytes(result, "model").Exists())
	assert.False(t, gjson.GetBytes(result, "stream").Exists())
	// max_tokens 应保留
	assert.Equal(t, int64(1024), gjson.GetBytes(result, "max_tokens").Int())
placeholder

func TestPrepareBedrockRequestBody_OutputFormatInlineSchema(t *testing.T) {
	t.Run("schema inlined into last user message array content", func(t *testing.T) {
		input := `{"model":"claude-sonnet-4-5","output_format":{"type":"json","schema":{"name":"string"placeholderplaceholder,"messages":[{"role":"user","content":[{"type":"text","text":"hello"placeholder]placeholder]placeholder`
		result, err := PrepareBedrockRequestBody([]byte(input), "us.anthropic.claude-sonnet-4-5-v1", "")
	placeholder

		assert.False(t, gjson.GetBytes(result, "output_format").Exists())
		// schema 应内联到最后一条 user message 的 content 数组末尾
		contentArr := gjson.GetBytes(result, "messages.0.content").Array()
		require.Len(t, contentArr, 2)
		assert.Equal(t, "text", contentArr[1].Get("type").String())
		assert.Contains(t, contentArr[1].Get("text").String(), `"name":"string"`)
placeholder)

	t.Run("schema inlined into string content", func(t *testing.T) {
		input := `{"model":"claude-sonnet-4-5","output_format":{"type":"json","schema":{"result":"number"placeholderplaceholder,"messages":[{"role":"user","content":"compute this"placeholder]placeholder`
		result, err := PrepareBedrockRequestBody([]byte(input), "us.anthropic.claude-sonnet-4-5-v1", "")
	placeholder

		assert.False(t, gjson.GetBytes(result, "output_format").Exists())
		contentArr := gjson.GetBytes(result, "messages.0.content").Array()
		require.Len(t, contentArr, 2)
		assert.Equal(t, "compute this", contentArr[0].Get("text").String())
		assert.Contains(t, contentArr[1].Get("text").String(), `"result":"number"`)
placeholder)

	t.Run("no schema field just removes output_format", func(t *testing.T) {
		input := `{"model":"claude-sonnet-4-5","output_format":{"type":"json"placeholder,"messages":[{"role":"user","content":"hi"placeholder]placeholder`
		result, err := PrepareBedrockRequestBody([]byte(input), "us.anthropic.claude-sonnet-4-5-v1", "")
	placeholder

		assert.False(t, gjson.GetBytes(result, "output_format").Exists())
placeholder)

	t.Run("no messages just removes output_format", func(t *testing.T) {
		input := `{"model":"claude-sonnet-4-5","output_format":{"type":"json","schema":{"name":"string"placeholderplaceholderplaceholder`
		result, err := PrepareBedrockRequestBody([]byte(input), "us.anthropic.claude-sonnet-4-5-v1", "")
	placeholder

		assert.False(t, gjson.GetBytes(result, "output_format").Exists())
placeholder)
placeholder

func TestPrepareBedrockRequestBody_RemoveOutputConfig(t *testing.T) {
	input := `{"model":"claude-sonnet-4-5","output_config":{"max_tokens":100placeholder,"messages":[]placeholder`
	result, err := PrepareBedrockRequestBody([]byte(input), "us.anthropic.claude-sonnet-4-5-v1", "")
placeholder

	assert.False(t, gjson.GetBytes(result, "output_config").Exists())
placeholder

func TestRemoveCustomFieldFromTools(t *testing.T) {
	input := `{
		"tools": [
			{"name":"tool1","custom":{"defer_loading":trueplaceholder,"description":"desc1"placeholder,
			{"name":"tool2","description":"desc2"placeholder,
			{"name":"tool3","custom":{"defer_loading":true,"other":placeholder,"description":"desc3"placeholder
		]
placeholder`
	result := removeCustomFieldFromTools([]byte(input))

	tools := gjson.GetBytes(result, "tools").Array()
	require.Len(t, tools, 3)
	// custom 应被移除
	assert.False(t, tools[0].Get("custom").Exists())
	// name/description 应保留
	assert.Equal(t, "tool1", tools[0].Get("name").String())
	assert.Equal(t, "desc1", tools[0].Get("description").String())
	// 没有 custom 的工具不受影响
	assert.Equal(t, "tool2", tools[1].Get("name").String())
	// 第三个工具的 custom 也应被移除
	assert.False(t, tools[2].Get("custom").Exists())
	assert.Equal(t, "tool3", tools[2].Get("name").String())
placeholder

func TestRemoveCustomFieldFromTools_NoTools(t *testing.T) {
	input := `{"messages":[{"role":"user","content":"hi"placeholder]placeholder`
	result := removeCustomFieldFromTools([]byte(input))
	// 无 tools 时不改变原始数据
	assert.JSONEq(t, input, string(result))
placeholder

func TestSanitizeBedrockCacheControl_RemoveScope(t *testing.T) {
	input := `{
		"system": [{"type":"text","text":"sys","cache_control":{"type":"ephemeral","scope":"global"placeholderplaceholder],
		"messages": [{"role":"user","content":[{"type":"text","text":"hi","cache_control":{"type":"ephemeral","scope":"global"placeholderplaceholder]placeholder]
placeholder`
	result := sanitizeBedrockCacheControl([]byte(input), "us.anthropic.claude-opus-4-6-v1")

	// scope 应被移除
	assert.False(t, gjson.GetBytes(result, "system.0.cache_control.scope").Exists())
	assert.False(t, gjson.GetBytes(result, "messages.0.content.0.cache_control.scope").Exists())
	// type 应保留
	assert.Equal(t, "ephemeral", gjson.GetBytes(result, "system.0.cache_control.type").String())
	assert.Equal(t, "ephemeral", gjson.GetBytes(result, "messages.0.content.0.cache_control.type").String())
placeholder

func TestSanitizeBedrockCacheControl_TTL_OldModel(t *testing.T) {
	input := `{
		"system": [{"type":"text","text":"sys","cache_control":{"type":"ephemeral","ttl":"5m"placeholderplaceholder]
placeholder`
	// 旧模型（Claude 3.5）不支持 ttl
	result := sanitizeBedrockCacheControl([]byte(input), "anthropic.claude-3-5-sonnet-20241022-v2:0")

	assert.False(t, gjson.GetBytes(result, "system.0.cache_control.ttl").Exists())
	assert.Equal(t, "ephemeral", gjson.GetBytes(result, "system.0.cache_control.type").String())
placeholder

func TestSanitizeBedrockCacheControl_TTL_Claude45_Supported(t *testing.T) {
	input := `{
		"system": [{"type":"text","text":"sys","cache_control":{"type":"ephemeral","ttl":"5m"placeholderplaceholder]
placeholder`
	// Claude 4.5+ 支持 "5m" 和 "1h"
	result := sanitizeBedrockCacheControl([]byte(input), "us.anthropic.claude-sonnet-4-5-20250929-v1:0")

	assert.True(t, gjson.GetBytes(result, "system.0.cache_control.ttl").Exists())
	assert.Equal(t, "5m", gjson.GetBytes(result, "system.0.cache_control.ttl").String())
placeholder

func TestSanitizeBedrockCacheControl_TTL_Claude45_UnsupportedValue(t *testing.T) {
	input := `{
		"system": [{"type":"text","text":"sys","cache_control":{"type":"ephemeral","ttl":"10m"placeholderplaceholder]
placeholder`
	// Claude 4.5 不支持 "10m"
	result := sanitizeBedrockCacheControl([]byte(input), "us.anthropic.claude-sonnet-4-5-20250929-v1:0")

	assert.False(t, gjson.GetBytes(result, "system.0.cache_control.ttl").Exists())
placeholder

func TestSanitizeBedrockCacheControl_TTL_Claude46(t *testing.T) {
	input := `{
		"messages": [{"role":"user","content":[{"type":"text","text":"hi","cache_control":{"type":"ephemeral","ttl":"1h"placeholderplaceholder]placeholder]
placeholder`
	result := sanitizeBedrockCacheControl([]byte(input), "us.anthropic.claude-opus-4-6-v1")

	assert.True(t, gjson.GetBytes(result, "messages.0.content.0.cache_control.ttl").Exists())
	assert.Equal(t, "1h", gjson.GetBytes(result, "messages.0.content.0.cache_control.ttl").String())
placeholder

func TestSanitizeBedrockCacheControl_NoCacheControl(t *testing.T) {
	input := `{"system":[{"type":"text","text":"sys"placeholder],"messages":[{"role":"user","content":[{"type":"text","text":"hi"placeholder]placeholder]placeholder`
	result := sanitizeBedrockCacheControl([]byte(input), "us.anthropic.claude-opus-4-6-v1")
	// 无 cache_control 时不改变原始数据
	assert.JSONEq(t, input, string(result))
placeholder

func TestIsBedrockClaude45OrNewer(t *testing.T) {
	tests := []struct {
		modelID string
		expect  bool
placeholder{
		{"us.anthropic.claude-opus-4-6-v1", trueplaceholder,
		{"us.anthropic.claude-sonnet-4-6", trueplaceholder,
		{"us.anthropic.claude-sonnet-4-5-20250929-v1:0", trueplaceholder,
		{"us.anthropic.claude-opus-4-5-20251101-v1:0", trueplaceholder,
		{"us.anthropic.placeholder-v1:0", trueplaceholder,
		{"anthropic.claude-3-5-sonnet-20241022-v2:0", falseplaceholder,
		{"anthropic.claude-3-opus-20240229-v1:0", falseplaceholder,
		{"anthropic.claude-3-haiku-20240307-v1:0", falseplaceholder,
		// 未来版本应自动支持
		{"us.anthropic.claude-sonnet-5-0-v1", trueplaceholder,
		{"us.anthropic.claude-opus-4-7-v1", trueplaceholder,
		// 旧版本
		{"anthropic.claude-opus-4-1-v1", falseplaceholder,
		{"anthropic.claude-sonnet-4-0-v1", falseplaceholder,
		// 非 Claude 模型
		{"amazon.nova-pro-v1", falseplaceholder,
		{"meta.llama3-70b", falseplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.modelID, func(t *testing.T) {
			assert.Equal(t, tt.expect, isBedrockClaude45OrNewer(tt.modelID))
	placeholder)
placeholder
placeholder

func TestPrepareBedrockRequestBody_FullIntegration(t *testing.T) {
	// 模拟一个完整的 Claude Code 请求
	input := `{
		"model": "claude-opus-4-6",
		"stream": true,
		"max_tokens": 16384,
		"output_format": {"type": "json", "schema": {"result": "string"placeholderplaceholder,
		"output_config": {"max_tokens": 100placeholder,
		"system": [{"type": "text", "text": "You are helpful", "cache_control": {"type": "ephemeral", "scope": "global", "ttl": "5m"placeholderplaceholder],
		"messages": [
			{"role": "user", "content": [{"type": "text", "text": "hello", "cache_control": {"type": "ephemeral", "ttl": "1h"placeholderplaceholder]placeholder
		],
		"tools": [
			{"name": "bash", "description": "Run bash", "custom": {"defer_loading": trueplaceholder, "input_schema": {"type": "object"placeholderplaceholder,
			{"name": "read", "description": "Read file", "input_schema": {"type": "object"placeholderplaceholder
		]
placeholder`

	betaHeader := "placeholder, context-1m-2025-08-07, compact-2026-01-12"
	result, err := PrepareBedrockRequestBody([]byte(input), "us.anthropic.claude-opus-4-6-v1", betaHeader)
placeholder

	// 基本字段
	assert.Equal(t, "bedrock-2023-05-31", gjson.GetBytes(result, "anthropic_version").String())
	assert.False(t, gjson.GetBytes(result, "model").Exists())
	assert.False(t, gjson.GetBytes(result, "stream").Exists())
	assert.Equal(t, int64(16384), gjson.GetBytes(result, "max_tokens").Int())

	// anthropic_beta 应包含所有 beta tokens
	betaArr := gjson.GetBytes(result, "anthropic_beta").Array()
	require.Len(t, betaArr, 3)
	assert.Equal(t, "placeholder", betaArr[0].String())
	assert.Equal(t, "context-1m-2025-08-07", betaArr[1].String())
	assert.Equal(t, "compact-2026-01-12", betaArr[2].String())

	// output_format 应被移除，schema 内联到最后一条 user message
	assert.False(t, gjson.GetBytes(result, "output_format").Exists())
	assert.False(t, gjson.GetBytes(result, "output_config").Exists())
	// content 数组：原始 text block + 内联 schema block
	contentArr := gjson.GetBytes(result, "messages.0.content").Array()
	require.Len(t, contentArr, 2)
	assert.Equal(t, "hello", contentArr[0].Get("text").String())
	assert.Contains(t, contentArr[1].Get("text").String(), `"result":"string"`)

	// tools 中的 custom 应被移除
	assert.False(t, gjson.GetBytes(result, "tools.0.custom").Exists())
	assert.Equal(t, "bash", gjson.GetBytes(result, "tools.0.name").String())
	assert.Equal(t, "read", gjson.GetBytes(result, "tools.1.name").String())

	// cache_control: scope 应被移除，ttl 在 Claude 4.6 上保留合法值
	assert.False(t, gjson.GetBytes(result, "system.0.cache_control.scope").Exists())
	assert.Equal(t, "ephemeral", gjson.GetBytes(result, "system.0.cache_control.type").String())
	assert.Equal(t, "5m", gjson.GetBytes(result, "system.0.cache_control.ttl").String())
	assert.Equal(t, "1h", gjson.GetBytes(result, "messages.0.content.0.cache_control.ttl").String())
placeholder

func TestPrepareBedrockRequestBody_BetaHeader(t *testing.T) {
	input := `{"messages":[{"role":"user","content":"hi"placeholder],"max_tokens":100placeholder`

	t.Run("empty beta header", func(t *testing.T) {
		result, err := PrepareBedrockRequestBody([]byte(input), "us.anthropic.claude-opus-4-6-v1", "")
	placeholder
		assert.False(t, gjson.GetBytes(result, "anthropic_beta").Exists())
placeholder)

	t.Run("single beta token", func(t *testing.T) {
		result, err := PrepareBedrockRequestBody([]byte(input), "us.anthropic.claude-opus-4-6-v1", "placeholder")
	placeholder
		arr := gjson.GetBytes(result, "anthropic_beta").Array()
		require.Len(t, arr, 1)
		assert.Equal(t, "placeholder", arr[0].String())
placeholder)

	t.Run("multiple beta tokens with spaces", func(t *testing.T) {
		result, err := PrepareBedrockRequestBody([]byte(input), "us.anthropic.claude-opus-4-6-v1", "placeholder , context-1m-2025-08-07 ")
	placeholder
		arr := gjson.GetBytes(result, "anthropic_beta").Array()
		require.Len(t, arr, 2)
		assert.Equal(t, "placeholder", arr[0].String())
		assert.Equal(t, "context-1m-2025-08-07", arr[1].String())
placeholder)

	t.Run("json array beta header", func(t *testing.T) {
		result, err := PrepareBedrockRequestBody([]byte(input), "us.anthropic.claude-opus-4-6-v1", `["placeholder","context-1m-2025-08-07"]`)
	placeholder
		arr := gjson.GetBytes(result, "anthropic_beta").Array()
		require.Len(t, arr, 2)
		assert.Equal(t, "placeholder", arr[0].String())
		assert.Equal(t, "context-1m-2025-08-07", arr[1].String())
placeholder)
placeholder

func TestParseAnthropicBetaHeader(t *testing.T) {
	assert.Nil(t, parseAnthropicBetaHeader(""))
	assert.Equal(t, []string{"a"placeholder, parseAnthropicBetaHeader("a"))
	assert.Equal(t, []string{"a", "b"placeholder, parseAnthropicBetaHeader("a,b"))
	assert.Equal(t, []string{"a", "b"placeholder, parseAnthropicBetaHeader("a , b "))
	assert.Equal(t, []string{"a", "b", "c"placeholder, parseAnthropicBetaHeader("a,b,c"))
	assert.Equal(t, []string{"a", "b"placeholder, parseAnthropicBetaHeader(`["a","b"]`))
placeholder

func TestFilterBedrockBetaTokens(t *testing.T) {
	t.Run("supported tokens pass through", func(t *testing.T) {
		tokens := []string{"placeholder", "context-1m-2025-08-07", "compact-2026-01-12"placeholder
		result := filterBedrockBetaTokens(tokens)
		assert.Equal(t, tokens, result)
placeholder)

	t.Run("unsupported tokens are filtered out", func(t *testing.T) {
		tokens := []string{"placeholder", "output-128k-2025-02-19", "files-api-2025-04-14", "structured-outputs-2025-11-13"placeholder
		result := filterBedrockBetaTokens(tokens)
		assert.Equal(t, []string{"placeholder"placeholder, result)
placeholder)

	t.Run("advanced-tool-use transforms to tool-search-tool", func(t *testing.T) {
		tokens := []string{"advanced-tool-use-2025-11-20"placeholder
		result := filterBedrockBetaTokens(tokens)
		assert.Contains(t, result, "tool-search-tool-2025-10-19")
		// tool-examples 自动关联
		assert.Contains(t, result, "tool-examples-2025-10-29")
placeholder)

	t.Run("tool-search-tool auto-associates tool-examples", func(t *testing.T) {
		tokens := []string{"tool-search-tool-2025-10-19"placeholder
		result := filterBedrockBetaTokens(tokens)
		assert.Contains(t, result, "tool-search-tool-2025-10-19")
		assert.Contains(t, result, "tool-examples-2025-10-29")
placeholder)

	t.Run("no duplication when tool-examples already present", func(t *testing.T) {
		tokens := []string{"tool-search-tool-2025-10-19", "tool-examples-2025-10-29"placeholder
		result := filterBedrockBetaTokens(tokens)
		count := 0
		for _, t := range result {
			if t == "tool-examples-2025-10-29" {
				count++
		placeholder
	placeholder
		assert.Equal(t, 1, count)
placeholder)

	t.Run("empty input returns nil", func(t *testing.T) {
		result := filterBedrockBetaTokens(nil)
		assert.Nil(t, result)
placeholder)

	t.Run("all unsupported returns nil", func(t *testing.T) {
		result := filterBedrockBetaTokens([]string{"output-128k-2025-02-19", "effort-2025-11-24"placeholder)
		assert.Nil(t, result)
placeholder)

	t.Run("duplicate tokens are deduplicated", func(t *testing.T) {
		tokens := []string{"context-1m-2025-08-07", "context-1m-2025-08-07"placeholder
		result := filterBedrockBetaTokens(tokens)
		assert.Equal(t, []string{"context-1m-2025-08-07"placeholder, result)
placeholder)
placeholder

func TestPrepareBedrockRequestBody_BetaFiltering(t *testing.T) {
	input := `{"messages":[{"role":"user","content":"hi"placeholder],"max_tokens":100placeholder`

	t.Run("unsupported beta tokens are filtered", func(t *testing.T) {
		result, err := PrepareBedrockRequestBody([]byte(input), "us.anthropic.claude-opus-4-6-v1",
			"placeholder, output-128k-2025-02-19, files-api-2025-04-14")
	placeholder
		arr := gjson.GetBytes(result, "anthropic_beta").Array()
		require.Len(t, arr, 1)
		assert.Equal(t, "placeholder", arr[0].String())
placeholder)

	t.Run("advanced-tool-use transformed in full pipeline", func(t *testing.T) {
		result, err := PrepareBedrockRequestBody([]byte(input), "us.anthropic.claude-opus-4-6-v1",
			"advanced-tool-use-2025-11-20")
	placeholder
		arr := gjson.GetBytes(result, "anthropic_beta").Array()
		require.Len(t, arr, 2)
		assert.Equal(t, "tool-search-tool-2025-10-19", arr[0].String())
		assert.Equal(t, "tool-examples-2025-10-29", arr[1].String())
placeholder)
placeholder

func TestBedrockCrossRegionPrefix(t *testing.T) {
	tests := []struct {
		region string
		expect string
placeholder{
		// US regions
		{"us-east-1", "us"placeholder,
		{"us-east-2", "us"placeholder,
		{"us-west-1", "us"placeholder,
		{"us-west-2", "us"placeholder,
		// GovCloud
		{"us-gov-east-1", "us-gov"placeholder,
		{"us-gov-west-1", "us-gov"placeholder,
		// EU regions
		{"eu-west-1", "eu"placeholder,
		{"eu-west-2", "eu"placeholder,
		{"eu-west-3", "eu"placeholder,
		{"eu-central-1", "eu"placeholder,
		{"eu-central-2", "eu"placeholder,
		{"eu-north-1", "eu"placeholder,
		{"eu-south-1", "eu"placeholder,
		// APAC regions
		{"ap-northeast-1", "jp"placeholder,
		{"ap-northeast-2", "apac"placeholder,
		{"ap-southeast-1", "apac"placeholder,
		{"ap-southeast-2", "au"placeholder,
		{"ap-south-1", "apac"placeholder,
		// Canada / South America fallback to us
		{"ca-central-1", "us"placeholder,
		{"sa-east-1", "us"placeholder,
		// Unknown defaults to us
		{"me-south-1", "us"placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.region, func(t *testing.T) {
			assert.Equal(t, tt.expect, BedrockCrossRegionPrefix(tt.region))
	placeholder)
placeholder
placeholder

func TestResolveBedrockModelID(t *testing.T) {
	t.Run("default alias resolves and adjusts region", func(t *testing.T) {
		account := &Account{
			Platform: PlatformAnthropic,
			Type:     AccountTypeBedrock,
	placeholder
				"aws_region": "eu-west-1",
		placeholder,
	placeholder

		modelID, ok := ResolveBedrockModelID(account, "claude-sonnet-4-5")
		require.True(t, ok)
		assert.Equal(t, "eu.anthropic.claude-sonnet-4-5-20250929-v1:0", modelID)
placeholder)

	t.Run("custom alias mapping reuses default bedrock mapping", func(t *testing.T) {
		account := &Account{
			Platform: PlatformAnthropic,
			Type:     AccountTypeBedrock,
	placeholder
				"aws_region": "ap-southeast-2",
				"model_mapping": map[string]any{
					"claude-*": "claude-opus-4-6",
			placeholder,
		placeholder,
	placeholder

		modelID, ok := ResolveBedrockModelID(account, "claude-opus-4-6-thinking")
		require.True(t, ok)
		assert.Equal(t, "au.anthropic.claude-opus-4-6-v1", modelID)
placeholder)

	t.Run("force global rewrites anthropic regional model id", func(t *testing.T) {
		account := &Account{
			Platform: PlatformAnthropic,
			Type:     AccountTypeBedrock,
	placeholder
				"aws_region":       "us-east-1",
				"aws_force_global": "true",
				"model_mapping": map[string]any{
					"claude-sonnet-4-6": "us.anthropic.claude-sonnet-4-6",
			placeholder,
		placeholder,
	placeholder

		modelID, ok := ResolveBedrockModelID(account, "claude-sonnet-4-6")
		require.True(t, ok)
		assert.Equal(t, "global.anthropic.claude-sonnet-4-6", modelID)
placeholder)

	t.Run("direct bedrock model id passes through", func(t *testing.T) {
		account := &Account{
			Platform: PlatformAnthropic,
			Type:     AccountTypeBedrock,
	placeholder
				"aws_region": "us-east-1",
		placeholder,
	placeholder

		modelID, ok := ResolveBedrockModelID(account, "anthropic.placeholder-v1:0")
		require.True(t, ok)
		assert.Equal(t, "anthropic.placeholder-v1:0", modelID)
placeholder)

	t.Run("unsupported alias returns false", func(t *testing.T) {
		account := &Account{
			Platform: PlatformAnthropic,
			Type:     AccountTypeBedrock,
	placeholder
				"aws_region": "us-east-1",
		placeholder,
	placeholder

		_, ok := ResolveBedrockModelID(account, "claude-3-5-sonnet-20241022")
		assert.False(t, ok)
placeholder)
placeholder

func TestAutoInjectBedrockBetaTokens(t *testing.T) {
	t.Run("inject interleaved-thinking when thinking present", func(t *testing.T) {
		body := []byte(`{"thinking":{"type":"enabled","budget_tokens":10000placeholder,"messages":[{"role":"user","content":"hi"placeholder]placeholder`)
		result := autoInjectBedrockBetaTokens(nil, body, "us.anthropic.claude-opus-4-6-v1")
		assert.Contains(t, result, "placeholder")
placeholder)

	t.Run("no duplicate when already present", func(t *testing.T) {
		body := []byte(`{"thinking":{"type":"enabled","budget_tokens":10000placeholder,"messages":[{"role":"user","content":"hi"placeholder]placeholder`)
		result := autoInjectBedrockBetaTokens([]string{"placeholder"placeholder, body, "us.anthropic.claude-opus-4-6-v1")
		count := 0
		for _, t := range result {
			if t == "placeholder" {
				count++
		placeholder
	placeholder
		assert.Equal(t, 1, count)
placeholder)

	t.Run("inject computer-use when computer tool present", func(t *testing.T) {
		body := []byte(`{"tools":[{"type":"computer_20250124","name":"computer","display_width_px":placeholder],"messages":[{"role":"user","content":"hi"placeholder]placeholder`)
		result := autoInjectBedrockBetaTokens(nil, body, "us.anthropic.claude-opus-4-6-v1")
		assert.Contains(t, result, "computer-use-2025-11-24")
placeholder)

	t.Run("inject advanced-tool-use for programmatic tool calling", func(t *testing.T) {
		body := []byte(`{"tools":[{"name":"bash","allowed_callers":["code_execution_20250825"]placeholder],"messages":[{"role":"user","content":"hi"placeholder]placeholder`)
		result := autoInjectBedrockBetaTokens(nil, body, "us.anthropic.claude-opus-4-6-v1")
		assert.Contains(t, result, "advanced-tool-use-2025-11-20")
placeholder)

	t.Run("inject advanced-tool-use for input examples", func(t *testing.T) {
		body := []byte(`{"tools":[{"name":"bash","input_examples":[{"cmd":"ls"placeholder]placeholder],"messages":[{"role":"user","content":"hi"placeholder]placeholder`)
		result := autoInjectBedrockBetaTokens(nil, body, "us.anthropic.claude-opus-4-6-v1")
		assert.Contains(t, result, "advanced-tool-use-2025-11-20")
placeholder)

	t.Run("inject tool-search-tool directly for pure tool search (no programmatic/inputExamples)", func(t *testing.T) {
		body := []byte(`{"tools":[{"type":"tool_search_tool_regex_20251119","name":"search"placeholder],"messages":[{"role":"user","content":"hi"placeholder]placeholder`)
		result := autoInjectBedrockBetaTokens(nil, body, "us.anthropic.claude-sonnet-4-6")
		// 纯 tool search 场景直接注入 Bedrock 特定头，不走 advanced-tool-use 转换
		assert.Contains(t, result, "tool-search-tool-2025-10-19")
		assert.NotContains(t, result, "advanced-tool-use-2025-11-20")
placeholder)

	t.Run("inject advanced-tool-use when tool search combined with programmatic calling", func(t *testing.T) {
		body := []byte(`{"tools":[{"type":"tool_search_tool_regex_20251119","name":"search"placeholder,{"name":"bash","allowed_callers":["code_execution_20250825"]placeholder],"messages":[{"role":"user","content":"hi"placeholder]placeholder`)
		result := autoInjectBedrockBetaTokens(nil, body, "us.anthropic.claude-sonnet-4-6")
		// 混合场景使用 advanced-tool-use（后续由 filter 转换为 tool-search-tool）
		assert.Contains(t, result, "advanced-tool-use-2025-11-20")
placeholder)

	t.Run("do not inject tool-search beta for unsupported models", func(t *testing.T) {
		body := []byte(`{"tools":[{"type":"tool_search_tool_regex_20251119","name":"search"placeholder],"messages":[{"role":"user","content":"hi"placeholder]placeholder`)
		result := autoInjectBedrockBetaTokens(nil, body, "anthropic.claude-3-5-sonnet-20241022-v2:0")
		assert.NotContains(t, result, "advanced-tool-use-2025-11-20")
		assert.NotContains(t, result, "tool-search-tool-2025-10-19")
placeholder)

	t.Run("no injection for regular tools", func(t *testing.T) {
		body := []byte(`{"tools":[{"name":"bash","description":"run bash","input_schema":{"type":"object"placeholderplaceholder],"messages":[{"role":"user","content":"hi"placeholder]placeholder`)
		result := autoInjectBedrockBetaTokens(nil, body, "us.anthropic.claude-opus-4-6-v1")
		assert.Empty(t, result)
placeholder)

	t.Run("no injection when no features detected", func(t *testing.T) {
		body := []byte(`{"messages":[{"role":"user","content":"hi"placeholder],"max_tokens":100placeholder`)
		result := autoInjectBedrockBetaTokens(nil, body, "us.anthropic.claude-opus-4-6-v1")
		assert.Empty(t, result)
placeholder)

	t.Run("preserves existing tokens", func(t *testing.T) {
		body := []byte(`{"thinking":{"type":"enabled"placeholder,"messages":[{"role":"user","content":"hi"placeholder]placeholder`)
		existing := []string{"context-1m-2025-08-07", "compact-2026-01-12"placeholder
		result := autoInjectBedrockBetaTokens(existing, body, "us.anthropic.claude-opus-4-6-v1")
		assert.Contains(t, result, "context-1m-2025-08-07")
		assert.Contains(t, result, "compact-2026-01-12")
		assert.Contains(t, result, "placeholder")
placeholder)
placeholder

func TestResolveBedrockBetaTokens(t *testing.T) {
	t.Run("body-only tool features resolve to final bedrock tokens", func(t *testing.T) {
		body := []byte(`{"tools":[{"name":"bash","allowed_callers":["code_execution_20250825"]placeholder],"messages":[{"role":"user","content":"hi"placeholder]placeholder`)
		result := ResolveBedrockBetaTokens("", body, "us.anthropic.claude-opus-4-6-v1")
		assert.Contains(t, result, "tool-search-tool-2025-10-19")
		assert.Contains(t, result, "tool-examples-2025-10-29")
placeholder)

	t.Run("unsupported client beta tokens are filtered out", func(t *testing.T) {
		body := []byte(`{"messages":[{"role":"user","content":"hi"placeholder]placeholder`)
		result := ResolveBedrockBetaTokens("placeholder,files-api-2025-04-14", body, "us.anthropic.claude-opus-4-6-v1")
		assert.Equal(t, []string{"placeholder"placeholder, result)
placeholder)
placeholder

func TestPrepareBedrockRequestBody_AutoBetaInjection(t *testing.T) {
	t.Run("thinking in body auto-injects beta without header", func(t *testing.T) {
		input := `{"messages":[{"role":"user","content":"hi"placeholder],"max_tokens":100,"thinking":{"type":"enabled","budget_tokens":10000placeholderplaceholder`
		result, err := PrepareBedrockRequestBody([]byte(input), "us.anthropic.claude-opus-4-6-v1", "")
	placeholder
		arr := gjson.GetBytes(result, "anthropic_beta").Array()
		found := false
		for _, v := range arr {
			if v.String() == "placeholder" {
				found = true
		placeholder
	placeholder
		assert.True(t, found, "interleaved-thinking should be auto-injected")
placeholder)

	t.Run("header tokens merged with auto-injected tokens", func(t *testing.T) {
		input := `{"messages":[{"role":"user","content":"hi"placeholder],"max_tokens":100,"thinking":{"type":"enabled","budget_tokens":10000placeholderplaceholder`
		result, err := PrepareBedrockRequestBody([]byte(input), "us.anthropic.claude-opus-4-6-v1", "context-1m-2025-08-07")
	placeholder
		arr := gjson.GetBytes(result, "anthropic_beta").Array()
		names := make([]string, len(arr))
		for i, v := range arr {
			names[i] = v.String()
	placeholder
		assert.Contains(t, names, "context-1m-2025-08-07")
		assert.Contains(t, names, "placeholder")
placeholder)
placeholder

func TestAdjustBedrockModelRegionPrefix(t *testing.T) {
	tests := []struct {
		name    string
		modelID string
		region  string
		expect  string
placeholder{
		// US region — no change needed
		{"us region keeps us prefix", "us.anthropic.claude-opus-4-6-v1", "us-east-1", "us.anthropic.claude-opus-4-6-v1"placeholder,
		// EU region — replace us → eu
		{"eu region replaces prefix", "us.anthropic.claude-opus-4-6-v1", "eu-west-1", "eu.anthropic.claude-opus-4-6-v1"placeholder,
		{"eu region sonnet", "us.anthropic.claude-sonnet-4-6", "eu-central-1", "eu.anthropic.claude-sonnet-4-6"placeholder,
		// APAC region — jp and au have dedicated prefixes per AWS docs
		{"jp region (ap-northeast-1)", "us.anthropic.claude-sonnet-4-5-20250929-v1:0", "ap-northeast-1", "jp.anthropic.claude-sonnet-4-5-20250929-v1:0"placeholder,
		{"au region (ap-southeast-2)", "us.anthropic.placeholder-v1:0", "ap-southeast-2", "au.anthropic.placeholder-v1:0"placeholder,
		{"apac region (ap-southeast-1)", "us.anthropic.claude-sonnet-4-5-20250929-v1:0", "ap-southeast-1", "apac.anthropic.claude-sonnet-4-5-20250929-v1:0"placeholder,
		// eu → us (user manually set eu prefix, moved to us region)
		{"eu to us", "eu.anthropic.claude-opus-4-6-v1", "us-west-2", "us.anthropic.claude-opus-4-6-v1"placeholder,
		// global prefix — replace to match region
		{"global to eu", "global.anthropic.claude-opus-4-6-v1", "eu-west-1", "eu.anthropic.claude-opus-4-6-v1"placeholder,
		// No known prefix — leave unchanged
		{"no prefix unchanged", "anthropic.claude-3-5-sonnet-20241022-v2:0", "eu-west-1", "anthropic.claude-3-5-sonnet-20241022-v2:0"placeholder,
		// GovCloud — uses independent us-gov prefix
		{"govcloud from us", "us.anthropic.claude-opus-4-6-v1", "us-gov-east-1", "us-gov.anthropic.claude-opus-4-6-v1"placeholder,
		{"govcloud already correct", "us-gov.anthropic.claude-opus-4-6-v1", "us-gov-west-1", "us-gov.anthropic.claude-opus-4-6-v1"placeholder,
		// Force global (special region value)
		{"force global from us", "us.anthropic.claude-opus-4-6-v1", "global", "global.anthropic.claude-opus-4-6-v1"placeholder,
		{"force global from eu", "eu.anthropic.claude-sonnet-4-6", "global", "global.anthropic.claude-sonnet-4-6"placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, AdjustBedrockModelRegionPrefix(tt.modelID, tt.region))
	placeholder)
placeholder
placeholder
