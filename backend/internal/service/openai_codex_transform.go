package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	opencodeCodexHeaderURL = "https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/session/prompt/codex_header.txt"
	codexCacheTTL          = 15 * time.Minute
)

var codexModelMap = map[string]string{
	"gpt-5.1-codex":             "gpt-5.1-codex",
	"gpt-5.1-codex-low":         "gpt-5.1-codex",
	"gpt-5.1-codex-medium":      "gpt-5.1-codex",
	"gpt-5.1-codex-high":        "gpt-5.1-codex",
	"gpt-5.1-codex-max":         "gpt-5.1-codex-max",
	"gpt-5.1-codex-max-low":     "gpt-5.1-codex-max",
	"gpt-5.1-codex-max-medium":  "gpt-5.1-codex-max",
	"gpt-5.1-codex-max-high":    "gpt-5.1-codex-max",
	"gpt-5.1-codex-max-xhigh":   "gpt-5.1-codex-max",
	"gpt-5.2":                   "gpt-5.2",
	"gpt-5.2-none":              "gpt-5.2",
	"gpt-5.2-low":               "gpt-5.2",
	"gpt-5.2-medium":            "gpt-5.2",
	"gpt-5.2-high":              "gpt-5.2",
	"gpt-5.2-xhigh":             "gpt-5.2",
	"gpt-5.2-codex":             "gpt-5.2-codex",
	"gpt-5.2-codex-low":         "gpt-5.2-codex",
	"gpt-5.2-codex-medium":      "gpt-5.2-codex",
	"gpt-5.2-codex-high":        "gpt-5.2-codex",
	"gpt-5.2-codex-xhigh":       "gpt-5.2-codex",
	"gpt-5.1-codex-mini":        "gpt-5.1-codex-mini",
	"gpt-5.1-codex-mini-medium": "gpt-5.1-codex-mini",
	"gpt-5.1-codex-mini-high":   "gpt-5.1-codex-mini",
	"gpt-5.1":                   "gpt-5.1",
	"gpt-5.1-none":              "gpt-5.1",
	"gpt-5.1-low":               "gpt-5.1",
	"gpt-5.1-medium":            "gpt-5.1",
	"gpt-5.1-high":              "gpt-5.1",
	"gpt-5.1-chat-latest":       "gpt-5.1",
	"gpt-5-codex":               "gpt-5.1-codex",
	"codex-mini-latest":         "gpt-5.1-codex-mini",
	"gpt-5-codex-mini":          "gpt-5.1-codex-mini",
	"gpt-5-codex-mini-medium":   "gpt-5.1-codex-mini",
	"gpt-5-codex-mini-high":     "gpt-5.1-codex-mini",
	"gpt-5":                     "gpt-5.1",
	"gpt-5-mini":                "gpt-5.1",
	"gpt-5-nano":                "gpt-5.1",
placeholder

type codexTransformResult struct {
	Modified        bool
	NormalizedModel string
	PromptCacheKey  string
placeholder

type opencodeCacheMetadata struct {
	ETag        string `json:"etag"`
	LastFetch   string `json:"lastFetch,omitempty"`
	LastChecked int64  `json:"lastChecked"`
placeholder

func applyCodexOAuthTransform(reqBody map[string]any) codexTransformResult {
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

	// 续链场景强制启用 store；非续链仍按原策略强制关闭存储。
	if needsToolContinuation {
		if v, ok := reqBody["store"].(bool); !ok || !v {
			reqBody["store"] = true
			result.Modified = true
	placeholder
placeholder else {
		if v, ok := reqBody["store"].(bool); !ok || v {
			reqBody["store"] = false
			result.Modified = true
	placeholder
placeholder
	if v, ok := reqBody["stream"].(bool); !ok || !v {
		reqBody["stream"] = true
		result.Modified = true
placeholder

	if _, ok := reqBody["max_output_tokens"]; ok {
		delete(reqBody, "max_output_tokens")
		result.Modified = true
placeholder
	if _, ok := reqBody["max_completion_tokens"]; ok {
		delete(reqBody, "max_completion_tokens")
		result.Modified = true
placeholder

	if normalizeCodexTools(reqBody) {
		result.Modified = true
placeholder

	if v, ok := reqBody["prompt_cache_key"].(string); ok {
		result.PromptCacheKey = strings.TrimSpace(v)
placeholder

	instructions := strings.TrimSpace(getOpenCodeCodexHeader())
	existingInstructions, _ := reqBody["instructions"].(string)
	existingInstructions = strings.TrimSpace(existingInstructions)

	if instructions != "" {
		if existingInstructions != instructions {
			reqBody["instructions"] = instructions
			result.Modified = true
	placeholder
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

	if strings.Contains(normalized, "gpt-5.2-codex") || strings.Contains(normalized, "gpt 5.2 codex") {
		return "gpt-5.2-codex"
placeholder
	if strings.Contains(normalized, "gpt-5.2") || strings.Contains(normalized, "gpt 5.2") {
		return "gpt-5.2"
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

func getOpenCodeCachedPrompt(url, cacheFileName, metaFileName string) string {
	cacheDir := codexCachePath("")
	if cacheDir == "" {
		return ""
placeholder
	cacheFile := filepath.Join(cacheDir, cacheFileName)
	metaFile := filepath.Join(cacheDir, metaFileName)

	var cachedContent string
	if content, ok := readFile(cacheFile); ok {
		cachedContent = content
placeholder

	var meta opencodeCacheMetadata
	if loadJSON(metaFile, &meta) && meta.LastChecked > 0 && cachedContent != "" {
		if time.Since(time.UnixMilli(meta.LastChecked)) < codexCacheTTL {
			return cachedContent
	placeholder
placeholder

	content, etag, status, err := fetchWithETag(url, meta.ETag)
	if err == nil && status == http.StatusNotModified && cachedContent != "" {
		return cachedContent
placeholder
	if err == nil && status >= 200 && status < 300 && content != "" {
		_ = writeFile(cacheFile, content)
		meta = opencodeCacheMetadata{
			ETag:        etag,
			LastFetch:   time.Now().UTC().Format(time.RFC3339),
			LastChecked: time.Now().UnixMilli(),
	placeholder
		_ = writeJSON(metaFile, meta)
		return content
placeholder

	return cachedContent
placeholder

func getOpenCodeCodexHeader() string {
	return getOpenCodeCachedPrompt(opencodeCodexHeaderURL, "opencode-codex-header.txt", "opencode-codex-header-meta.json")
placeholder

func GetOpenCodeInstructions() string {
	return getOpenCodeCodexHeader()
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
		if typ, ok := m["type"].(string); ok && typ == "item_reference" {
			if !preserveReferences {
				continue
		placeholder
	placeholder
		newItem := m
		if !preserveReferences {
			newItem = make(map[string]any, len(m))
			for key, value := range m {
				newItem[key] = value
		placeholder
			delete(newItem, "id")
	placeholder
		filtered = append(filtered, newItem)
placeholder
	return filtered
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
	for idx, tool := range tools {
		toolMap, ok := tool.(map[string]any)
		if !ok {
			continue
	placeholder

		toolType, _ := toolMap["type"].(string)
		if strings.TrimSpace(toolType) != "function" {
			continue
	placeholder

		function, ok := toolMap["function"].(map[string]any)
		if !ok {
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

		tools[idx] = toolMap
placeholder

	if modified {
		reqBody["tools"] = tools
placeholder

	return modified
placeholder

func codexCachePath(filename string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
placeholder
	cacheDir := filepath.Join(home, ".opencode", "cache")
	if filename == "" {
		return cacheDir
placeholder
	return filepath.Join(cacheDir, filename)
placeholder

func readFile(path string) (string, bool) {
	if path == "" {
		return "", false
placeholder
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
placeholder
	return string(data), true
placeholder

func writeFile(path, content string) error {
	if path == "" {
		return fmt.Errorf("empty cache path")
placeholder
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
placeholder
	return os.WriteFile(path, []byte(content), 0o644)
placeholder

func loadJSON(path string, target any) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
placeholder
	if err := json.Unmarshal(data, target); err != nil {
		return false
placeholder
	return true
placeholder

func writeJSON(path string, value any) error {
	if path == "" {
		return fmt.Errorf("empty json path")
placeholder
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
placeholder
	data, err := json.Marshal(value)
	if err != nil {
		return err
placeholder
	return os.WriteFile(path, data, 0o644)
placeholder

func fetchWithETag(url, etag string) (string, string, int, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", "", 0, err
placeholder
	req.Header.Set("User-Agent", "sub2api-codex")
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
placeholder
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", 0, err
placeholder
	defer func() {
		_ = resp.Body.Close()
placeholder()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", resp.StatusCode, err
placeholder
	return string(body), resp.Header.Get("etag"), resp.StatusCode, nil
placeholder
