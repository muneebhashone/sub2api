package service

import (
	_ "embed"
	"strings"
)

//go:embed prompts/codex_cli_instructions.md
var codexCLIInstructions string

var codexModelMap = map[string]string{
	"gpt-5.4":                    "gpt-5.4",
	"gpt-5.4-none":               "gpt-5.4",
	"gpt-5.4-low":                "gpt-5.4",
	"gpt-5.4-medium":             "gpt-5.4",
	"gpt-5.4-high":               "gpt-5.4",
	"gpt-5.4-xhigh":              "gpt-5.4",
	"gpt-5.4-chat-latest":        "gpt-5.4",
	"gpt-5.3":                    "gpt-5.3-codex",
	"gpt-5.3-none":               "gpt-5.3-codex",
	"gpt-5.3-low":                "gpt-5.3-codex",
	"gpt-5.3-medium":             "gpt-5.3-codex",
	"gpt-5.3-high":               "gpt-5.3-codex",
	"gpt-5.3-xhigh":              "gpt-5.3-codex",
	"gpt-5.3-codex":              "gpt-5.3-codex",
	"gpt-5.3-codex-spark":        "gpt-5.3-codex",
	"gpt-5.3-codex-spark-low":    "gpt-5.3-codex",
	"gpt-5.3-codex-spark-medium": "gpt-5.3-codex",
	"gpt-5.3-codex-spark-high":   "gpt-5.3-codex",
	"gpt-5.3-codex-spark-xhigh":  "gpt-5.3-codex",
	"gpt-5.3-codex-low":          "gpt-5.3-codex",
	"gpt-5.3-codex-medium":       "gpt-5.3-codex",
	"gpt-5.3-codex-high":         "gpt-5.3-codex",
	"gpt-5.3-codex-xhigh":        "gpt-5.3-codex",
	"gpt-5.1-codex":              "gpt-5.1-codex",
	"gpt-5.1-codex-low":          "gpt-5.1-codex",
	"gpt-5.1-codex-medium":       "gpt-5.1-codex",
	"gpt-5.1-codex-high":         "gpt-5.1-codex",
	"gpt-5.1-codex-max":          "gpt-5.1-codex-max",
	"gpt-5.1-codex-max-low":      "gpt-5.1-codex-max",
	"gpt-5.1-codex-max-medium":   "gpt-5.1-codex-max",
	"gpt-5.1-codex-max-high":     "gpt-5.1-codex-max",
	"gpt-5.1-codex-max-xhigh":    "gpt-5.1-codex-max",
	"gpt-5.2":                    "gpt-5.2",
	"gpt-5.2-none":               "gpt-5.2",
	"gpt-5.2-low":                "gpt-5.2",
	"gpt-5.2-medium":             "gpt-5.2",
	"gpt-5.2-high":               "gpt-5.2",
	"gpt-5.2-xhigh":              "gpt-5.2",
	"gpt-5.2-codex":              "gpt-5.2-codex",
	"gpt-5.2-codex-low":          "gpt-5.2-codex",
	"gpt-5.2-codex-medium":       "gpt-5.2-codex",
	"gpt-5.2-codex-high":         "gpt-5.2-codex",
	"gpt-5.2-codex-xhigh":        "gpt-5.2-codex",
	"gpt-5.1-codex-mini":         "gpt-5.1-codex-mini",
	"gpt-5.1-codex-mini-medium":  "gpt-5.1-codex-mini",
	"gpt-5.1-codex-mini-high":    "gpt-5.1-codex-mini",
	"gpt-5.1":                    "gpt-5.1",
	"gpt-5.1-none":               "gpt-5.1",
	"gpt-5.1-low":                "gpt-5.1",
	"gpt-5.1-medium":             "gpt-5.1",
	"gpt-5.1-high":               "gpt-5.1",
	"gpt-5.1-chat-latest":        "gpt-5.1",
	"gpt-5-codex":                "gpt-5.1-codex",
	"codex-mini-latest":          "gpt-5.1-codex-mini",
	"gpt-5-codex-mini":           "gpt-5.1-codex-mini",
	"gpt-5-codex-mini-medium":    "gpt-5.1-codex-mini",
	"gpt-5-codex-mini-high":      "gpt-5.1-codex-mini",
	"gpt-5":                      "gpt-5.1",
	"gpt-5-mini":                 "gpt-5.1",
	"gpt-5-nano":                 "gpt-5.1",
placeholder

type codexTransformResult struct {
	Modified        bool
	NormalizedModel string
	PromptCacheKey  string
placeholder

func applyCodexOAuthTransform(reqBody map[string]any, isCodexCLI bool, isCompact bool) codexTransformResult {
	result := codexTransformResult{placeholder
	// 工具续链需求会影响存储策略与 input 过滤逻辑。
	needsToolContinuation := NeedsToolContinuation(reqBody)

	model := ""
	if v, ok := reqBody["model"].(string); ok {
		model = v
placeholder
	normalizedModel := normalizeCodexModel(model)
	if normalizedModel != "" {
		if model != normalizedModel {
			reqBody["model"] = normalizedModel
			result.Modified = true
	placeholder
		result.NormalizedModel = normalizedModel
placeholder

	if isCompact {
		if _, ok := reqBody["store"]; ok {
			delete(reqBody, "store")
			result.Modified = true
	placeholder
		if _, ok := reqBody["stream"]; ok {
			delete(reqBody, "stream")
			result.Modified = true
	placeholder
placeholder else {
		// OAuth 走 ChatGPT internal API 时，store 必须为 false；显式 true 也会强制覆盖。
		// 避免上游返回 "Store must be set to false"。
		if v, ok := reqBody["store"].(bool); !ok || v {
			reqBody["store"] = false
			result.Modified = true
	placeholder
		if v, ok := reqBody["stream"].(bool); !ok || !v {
			reqBody["stream"] = true
			result.Modified = true
	placeholder
placeholder

	// Strip parameters unsupported by codex models via the Responses API.
	for _, key := range []string{
		"max_output_tokens",
		"max_completion_tokens",
		"temperature",
		"top_p",
		"frequency_penalty",
		"presence_penalty",
placeholder {
		if _, ok := reqBody[key]; ok {
			delete(reqBody, key)
			result.Modified = true
	placeholder
placeholder

	if normalizeCodexTools(reqBody) {
		result.Modified = true
placeholder

	if v, ok := reqBody["prompt_cache_key"].(string); ok {
		result.PromptCacheKey = strings.TrimSpace(v)
placeholder

	// instructions 处理逻辑：根据是否是 Codex CLI 分别调用不同方法
	if applyInstructions(reqBody, isCodexCLI) {
		result.Modified = true
placeholder

	// 续链场景保留 item_reference 与 id，避免 call_id 上下文丢失。
	if input, ok := reqBody["input"].([]any); ok {
		input = filterCodexInput(input, needsToolContinuation)
		reqBody["input"] = input
		result.Modified = true
placeholder

	return result
placeholder

func normalizeCodexModel(model string) string {
	if model == "" {
		return "gpt-5.1"
placeholder

	modelID := model
	if strings.Contains(modelID, "/") {
		parts := strings.Split(modelID, "/")
		modelID = parts[len(parts)-1]
placeholder

	if mapped := getNormalizedCodexModel(modelID); mapped != "" {
		return mapped
placeholder

	normalized := strings.ToLower(modelID)

	if strings.Contains(normalized, "gpt-5.4") || strings.Contains(normalized, "gpt 5.4") {
		return "gpt-5.4"
placeholder
	if strings.Contains(normalized, "gpt-5.2-codex") || strings.Contains(normalized, "gpt 5.2 codex") {
		return "gpt-5.2-codex"
placeholder
	if strings.Contains(normalized, "gpt-5.2") || strings.Contains(normalized, "gpt 5.2") {
		return "gpt-5.2"
placeholder
	if strings.Contains(normalized, "gpt-5.3-codex") || strings.Contains(normalized, "gpt 5.3 codex") {
		return "gpt-5.3-codex"
placeholder
	if strings.Contains(normalized, "gpt-5.3") || strings.Contains(normalized, "gpt 5.3") {
		return "gpt-5.3-codex"
placeholder
	if strings.Contains(normalized, "gpt-5.1-codex-max") || strings.Contains(normalized, "gpt 5.1 codex max") {
		return "gpt-5.1-codex-max"
placeholder
	if strings.Contains(normalized, "gpt-5.1-codex-mini") || strings.Contains(normalized, "gpt 5.1 codex mini") {
		return "gpt-5.1-codex-mini"
placeholder
	if strings.Contains(normalized, "codex-mini-latest") ||
		strings.Contains(normalized, "gpt-5-codex-mini") ||
		strings.Contains(normalized, "gpt 5 codex mini") {
		return "codex-mini-latest"
placeholder
	if strings.Contains(normalized, "gpt-5.1-codex") || strings.Contains(normalized, "gpt 5.1 codex") {
		return "gpt-5.1-codex"
placeholder
	if strings.Contains(normalized, "gpt-5.1") || strings.Contains(normalized, "gpt 5.1") {
		return "gpt-5.1"
placeholder
	if strings.Contains(normalized, "codex") {
		return "gpt-5.1-codex"
placeholder
	if strings.Contains(normalized, "gpt-5") || strings.Contains(normalized, "gpt 5") {
		return "gpt-5.1"
placeholder

	return "gpt-5.1"
placeholder

func getNormalizedCodexModel(modelID string) string {
	if modelID == "" {
		return ""
placeholder
	if mapped, ok := codexModelMap[modelID]; ok {
		return mapped
placeholder
	lower := strings.ToLower(modelID)
	for key, value := range codexModelMap {
		if strings.ToLower(key) == lower {
			return value
	placeholder
placeholder
	return ""
placeholder

func getOpenCodeCodexHeader() string {
	// 兼容保留：历史上这里会从 opencode 仓库拉取 codex_header.txt。
	// 现在我们与 Codex CLI 一致，直接使用仓库内置的 instructions，避免读写缓存与外网依赖。
	return getCodexCLIInstructions()
placeholder

func getCodexCLIInstructions() string {
	return codexCLIInstructions
placeholder

func GetOpenCodeInstructions() string {
	return getOpenCodeCodexHeader()
placeholder

// GetCodexCLIInstructions 返回内置的 Codex CLI 指令内容。
func GetCodexCLIInstructions() string {
	return getCodexCLIInstructions()
placeholder

// applyInstructions 处理 instructions 字段
// isCodexCLI=true: 仅补充缺失的 instructions（使用内置 Codex CLI 指令）
// isCodexCLI=false: 优先使用内置 Codex CLI 指令覆盖
func applyInstructions(reqBody map[string]any, isCodexCLI bool) bool {
	if isCodexCLI {
		return applyCodexCLIInstructions(reqBody)
placeholder
	return applyOpenCodeInstructions(reqBody)
placeholder

// applyCodexCLIInstructions 为 Codex CLI 请求补充缺失的 instructions
// 仅在 instructions 为空时添加内置 Codex CLI 指令（不依赖 opencode 缓存/回源）
func applyCodexCLIInstructions(reqBody map[string]any) bool {
	if !isInstructionsEmpty(reqBody) {
		return false // 已有有效 instructions，不修改
placeholder

	instructions := strings.TrimSpace(getCodexCLIInstructions())
	if instructions != "" {
		reqBody["instructions"] = instructions
		return true
placeholder

	return false
placeholder

// applyOpenCodeInstructions 为非 Codex CLI 请求应用内置 Codex CLI 指令（兼容历史函数名）
// 优先使用内置 Codex CLI 指令覆盖
func applyOpenCodeInstructions(reqBody map[string]any) bool {
	instructions := strings.TrimSpace(getOpenCodeCodexHeader())
	existingInstructions, _ := reqBody["instructions"].(string)
	existingInstructions = strings.TrimSpace(existingInstructions)

	if instructions != "" {
		if existingInstructions != instructions {
			reqBody["instructions"] = instructions
			return true
	placeholder
placeholder else if existingInstructions == "" {
		codexInstructions := strings.TrimSpace(getCodexCLIInstructions())
		if codexInstructions != "" {
			reqBody["instructions"] = codexInstructions
			return true
	placeholder
placeholder

	return false
placeholder

// isInstructionsEmpty 检查 instructions 字段是否为空
// 处理以下情况：字段不存在、nil、空字符串、纯空白字符串
func isInstructionsEmpty(reqBody map[string]any) bool {
	val, exists := reqBody["instructions"]
	if !exists {
		return true
placeholder
	if val == nil {
		return true
placeholder
	str, ok := val.(string)
	if !ok {
		return true
placeholder
	return strings.TrimSpace(str) == ""
placeholder

// filterCodexInput 按需过滤 item_reference 与 id。
// preserveReferences 为 true 时保持引用与 id，以满足续链请求对上下文的依赖。
func filterCodexInput(input []any, preserveReferences bool) []any {
	filtered := make([]any, 0, len(input))
	for _, item := range input {
		m, ok := item.(map[string]any)
		if !ok {
			filtered = append(filtered, item)
			continue
	placeholder
		typ, _ := m["type"].(string)
		if typ == "item_reference" {
			if !preserveReferences {
				continue
		placeholder
			newItem := make(map[string]any, len(m))
			for key, value := range m {
				newItem[key] = value
		placeholder
			filtered = append(filtered, newItem)
			continue
	placeholder

		newItem := m
		copied := false
		// 仅在需要修改字段时创建副本，避免直接改写原始输入。
		ensureCopy := func() {
			if copied {
				return
		placeholder
			newItem = make(map[string]any, len(m))
			for key, value := range m {
				newItem[key] = value
		placeholder
			copied = true
	placeholder

		if isCodexToolCallItemType(typ) {
			if callID, ok := m["call_id"].(string); !ok || strings.TrimSpace(callID) == "" {
				if id, ok := m["id"].(string); ok && strings.TrimSpace(id) != "" {
					ensureCopy()
					newItem["call_id"] = id
			placeholder
		placeholder
	placeholder

		if !preserveReferences {
			ensureCopy()
			delete(newItem, "id")
			if !isCodexToolCallItemType(typ) {
				delete(newItem, "call_id")
		placeholder
	placeholder

		filtered = append(filtered, newItem)
placeholder
	return filtered
placeholder

func isCodexToolCallItemType(typ string) bool {
	if typ == "" {
		return false
placeholder
	return strings.HasSuffix(typ, "_call") || strings.HasSuffix(typ, "_call_output")
placeholder

func normalizeCodexTools(reqBody map[string]any) bool {
	rawTools, ok := reqBody["tools"]
	if !ok || rawTools == nil {
		return false
placeholder
	tools, ok := rawTools.([]any)
	if !ok {
		return false
placeholder

	modified := false
	validTools := make([]any, 0, len(tools))

	for _, tool := range tools {
		toolMap, ok := tool.(map[string]any)
		if !ok {
			// Keep unknown structure as-is to avoid breaking upstream behavior.
			validTools = append(validTools, tool)
			continue
	placeholder

		toolType, _ := toolMap["type"].(string)
		toolType = strings.TrimSpace(toolType)
		if toolType != "function" {
			validTools = append(validTools, toolMap)
			continue
	placeholder

		// OpenAI Responses-style tools use top-level name/parameters.
		if name, ok := toolMap["name"].(string); ok && strings.TrimSpace(name) != "" {
			validTools = append(validTools, toolMap)
			continue
	placeholder

		// ChatCompletions-style tools use {type:"function", function:{...placeholderplaceholder.
		functionValue, hasFunction := toolMap["function"]
		function, ok := functionValue.(map[string]any)
		if !hasFunction || functionValue == nil || !ok || function == nil {
			// Drop invalid function tools.
			modified = true
			continue
	placeholder

		if _, ok := toolMap["name"]; !ok {
			if name, ok := function["name"].(string); ok && strings.TrimSpace(name) != "" {
				toolMap["name"] = name
				modified = true
		placeholder
	placeholder
		if _, ok := toolMap["description"]; !ok {
			if desc, ok := function["description"].(string); ok && strings.TrimSpace(desc) != "" {
				toolMap["description"] = desc
				modified = true
		placeholder
	placeholder
		if _, ok := toolMap["parameters"]; !ok {
			if params, ok := function["parameters"]; ok {
				toolMap["parameters"] = params
				modified = true
		placeholder
	placeholder
		if _, ok := toolMap["strict"]; !ok {
			if strict, ok := function["strict"]; ok {
				toolMap["strict"] = strict
				modified = true
		placeholder
	placeholder

		validTools = append(validTools, toolMap)
placeholder

	if modified {
		reqBody["tools"] = validTools
placeholder

	return modified
placeholder
