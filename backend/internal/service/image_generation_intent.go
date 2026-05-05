package service

import (
	"encoding/json"
	"strings"

	"github.com/tidwall/gjson"
)

const (
	openAIResponsesEndpoint          = "/v1/responses"
	openAIResponsesCompactEndpoint   = "/v1/responses/compact"
	imageGenerationPermissionMessage = "Image generation is not enabled for this group"
)

// ImageGenerationPermissionMessage returns the stable end-user error text for disabled groups.
func ImageGenerationPermissionMessage() string {
	return imageGenerationPermissionMessage
placeholder

// GroupAllowsImageGeneration preserves ungrouped-key behavior and enforces the flag when a group is present.
func GroupAllowsImageGeneration(group *Group) bool {
	return group == nil || group.AllowImageGeneration
placeholder

// IsImageGenerationIntent classifies requests that can produce generated images.
func IsImageGenerationIntent(endpoint string, requestedModel string, body []byte) bool {
	if IsImageGenerationEndpoint(endpoint) {
		return true
placeholder
	if isOpenAIImageGenerationModel(requestedModel) {
		return true
placeholder
	if len(body) == 0 || !gjson.ValidBytes(body) {
		return false
placeholder
	if model := strings.TrimSpace(gjson.GetBytes(body, "model").String()); isOpenAIImageGenerationModel(model) {
		return true
placeholder
	if openAIJSONToolsContainImageGeneration(gjson.GetBytes(body, "tools")) {
		return true
placeholder
	return openAIJSONToolChoiceSelectsImageGeneration(gjson.GetBytes(body, "tool_choice"))
placeholder

// IsImageGenerationIntentMap is the map-backed variant used after service-side request mutation.
func IsImageGenerationIntentMap(endpoint string, requestedModel string, reqBody map[string]any) bool {
	if IsImageGenerationEndpoint(endpoint) {
		return true
placeholder
	if isOpenAIImageGenerationModel(requestedModel) {
		return true
placeholder
	if reqBody == nil {
		return false
placeholder
	if isOpenAIImageGenerationModel(firstNonEmptyString(reqBody["model"])) {
		return true
placeholder
	if hasOpenAIImageGenerationTool(reqBody) {
		return true
placeholder
	return openAIAnyToolChoiceSelectsImageGeneration(reqBody["tool_choice"])
placeholder

// IsImageGenerationEndpoint identifies dedicated generated-image endpoints.
func IsImageGenerationEndpoint(endpoint string) bool {
	switch normalizeImageGenerationEndpoint(endpoint) {
	case "/v1/images/generations", "/v1/images/edits", "/images/generations", "/images/edits":
		return true
	default:
		return false
placeholder
placeholder

func normalizeImageGenerationEndpoint(endpoint string) string {
	endpoint = strings.TrimSpace(strings.ToLower(endpoint))
	if endpoint == "" {
		return ""
placeholder
	endpoint = strings.TrimPrefix(endpoint, "https://api.openai.com")
	if idx := strings.IndexByte(endpoint, '?'); idx >= 0 {
		endpoint = endpoint[:idx]
placeholder
	return strings.TrimRight(endpoint, "/")
placeholder

func openAIJSONToolsContainImageGeneration(tools gjson.Result) bool {
	if !tools.IsArray() {
		return false
placeholder
	found := false
	tools.ForEach(func(_, item gjson.Result) bool {
		if strings.TrimSpace(item.Get("type").String()) == "image_generation" {
			found = true
			return false
	placeholder
		return true
placeholder)
	return found
placeholder

func openAIJSONToolChoiceSelectsImageGeneration(choice gjson.Result) bool {
	if !choice.Exists() {
		return false
placeholder
	if choice.Type == gjson.String {
		return strings.TrimSpace(choice.String()) == "image_generation"
placeholder
	if !choice.IsObject() {
		return false
placeholder
	if strings.TrimSpace(choice.Get("type").String()) == "image_generation" {
		return true
placeholder
	if strings.TrimSpace(choice.Get("tool.type").String()) == "image_generation" {
		return true
placeholder
	if strings.TrimSpace(choice.Get("function.name").String()) == "image_generation" {
		return true
placeholder
	return false
placeholder

func openAIAnyToolChoiceSelectsImageGeneration(choice any) bool {
	switch v := choice.(type) {
	case string:
		return strings.TrimSpace(v) == "image_generation"
	case map[string]any:
		if strings.TrimSpace(firstNonEmptyString(v["type"])) == "image_generation" {
			return true
	placeholder
		if tool, ok := v["tool"].(map[string]any); ok && strings.TrimSpace(firstNonEmptyString(tool["type"])) == "image_generation" {
			return true
	placeholder
		if fn, ok := v["function"].(map[string]any); ok && strings.TrimSpace(firstNonEmptyString(fn["name"])) == "image_generation" {
			return true
	placeholder
placeholder
	return false
placeholder

func getAPIKeyFromContext(c interface{ Get(string) (any, bool) placeholder) *APIKey {
	if c == nil {
		return nil
placeholder
	v, exists := c.Get("api_key")
	if !exists {
		return nil
placeholder
	apiKey, _ := v.(*APIKey)
	return apiKey
placeholder

func apiKeyGroup(apiKey *APIKey) *Group {
	if apiKey == nil {
		return nil
placeholder
	return apiKey.Group
placeholder

func cloneRequestMapForImageIntent(body []byte) map[string]any {
	if len(body) == 0 {
		return nil
placeholder
	var out map[string]any
	if err := json.Unmarshal(body, &out); err != nil {
		return nil
placeholder
	return out
placeholder

func resolveOpenAIResponsesImageBillingConfig(reqBody map[string]any, fallbackModel string) (string, string, error) {
	imageModel := ""
	imageSize := ""
	hasImageTool := false
	if reqBody != nil {
		rawTools, _ := reqBody["tools"].([]any)
		for _, rawTool := range rawTools {
			toolMap, ok := rawTool.(map[string]any)
			if !ok || strings.TrimSpace(firstNonEmptyString(toolMap["type"])) != "image_generation" {
				continue
		placeholder
			hasImageTool = true
			imageModel = strings.TrimSpace(firstNonEmptyString(toolMap["model"]))
			imageSize = strings.TrimSpace(firstNonEmptyString(toolMap["size"]))
			break
	placeholder
		if imageSize == "" {
			imageSize = strings.TrimSpace(firstNonEmptyString(reqBody["size"]))
	placeholder
placeholder
	if imageModel == "" && reqBody != nil {
		bodyModel := strings.TrimSpace(firstNonEmptyString(reqBody["model"]))
		if isOpenAIImageBillingModelAlias(bodyModel) || !hasImageTool {
			imageModel = bodyModel
	placeholder
placeholder
	if imageModel == "" && hasImageTool {
		imageModel = "gpt-image-2"
placeholder
	if imageModel == "" {
		imageModel = strings.TrimSpace(fallbackModel)
placeholder
	sizeTier := normalizeOpenAIImageSizeTier(imageSize)
	return imageModel, sizeTier, nil
placeholder

func resolveOpenAIResponsesImageBillingConfigFromBody(body []byte, fallbackModel string) (string, string, error) {
	reqBody := cloneRequestMapForImageIntent(body)
	return resolveOpenAIResponsesImageBillingConfig(reqBody, fallbackModel)
placeholder

func isOpenAIImageBillingModelAlias(model string) bool {
	normalized := strings.ToLower(strings.TrimSpace(model))
	if normalized == "" {
		return false
placeholder
	return isOpenAIImageGenerationModel(normalized) || strings.Contains(normalized, "image")
placeholder
