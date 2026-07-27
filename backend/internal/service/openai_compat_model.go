package service

import (
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/apicompat"
)

func NormalizeOpenAICompatRequestedModel(model string) string {
	trimmed := strings.TrimSpace(model)
	if trimmed == "" {
		return ""
placeholder

	normalized, _, ok := splitOpenAICompatReasoningModel(trimmed)
	if !ok || normalized == "" {
		return trimmed
placeholder
	return normalized
placeholder

func applyOpenAICompatModelNormalization(req *apicompat.AnthropicRequest) {
	if req == nil {
		return
placeholder

	originalModel := strings.TrimSpace(req.Model)
	if originalModel == "" {
		return
placeholder

	normalizedModel, derivedEffort, hasReasoningSuffix := splitOpenAICompatReasoningModel(originalModel)
	if hasReasoningSuffix && normalizedModel != "" {
		req.Model = normalizedModel
placeholder

	if req.OutputConfig != nil && strings.TrimSpace(req.OutputConfig.Effort) != "" {
		return
placeholder

	claudeEffort := openAIReasoningEffortToClaudeOutputEffort(derivedEffort)
	if claudeEffort == "" {
		return
placeholder

	if req.OutputConfig == nil {
		req.OutputConfig = &apicompat.AnthropicOutputConfig{placeholder
placeholder
	req.OutputConfig.Effort = claudeEffort
placeholder

func splitOpenAICompatReasoningModel(model string) (normalizedModel string, reasoningEffort string, ok bool) {
	trimmed := strings.TrimSpace(model)
	if trimmed == "" {
		return "", "", false
placeholder

	modelID := trimmed
	if strings.Contains(modelID, "/") {
		parts := strings.Split(modelID, "/")
		modelID = parts[len(parts)-1]
placeholder
	modelID = strings.TrimSpace(modelID)
	if !strings.HasPrefix(strings.ToLower(modelID), "gpt-") {
		return trimmed, "", false
placeholder

	parts := strings.FieldsFunc(strings.ToLower(modelID), func(r rune) bool {
		switch r {
		case '-', '_', ' ':
			return true
		default:
			return false
	placeholder
placeholder)
	if len(parts) == 0 {
		return trimmed, "", false
placeholder

	last := strings.NewReplacer("-", "", "_", "", " ", "").Replace(parts[len(parts)-1])
	switch last {
	case "none", "minimal":
	case "low", "medium", "high":
		reasoningEffort = last
	case "xhigh", "extrahigh":
		reasoningEffort = "xhigh"
	default:
		return trimmed, "", false
placeholder

	return normalizeCodexModel(modelID), reasoningEffort, true
placeholder

func openAIReasoningEffortToClaudeOutputEffort(effort string) string {
	switch strings.TrimSpace(effort) {
	case "low", "medium", "high":
		return effort
	case "xhigh":
		return "max"
	default:
		return ""
placeholder
placeholder

// openAICompatAnthropicReasoningEffort resolves the effort emitted by the
// Anthropic bridge after the final upstream model is known. Anthropic's max is
// normally translated to OpenAI xhigh, but GPT-5.6 accepts the original max
// value on Responses and Chat Completions.
func openAICompatAnthropicReasoningEffort(req *apicompat.AnthropicRequest, upstreamModel, convertedEffort string) string {
	if req == nil || req.OutputConfig == nil || !strings.EqualFold(strings.TrimSpace(req.OutputConfig.Effort), "max") {
		return convertedEffort
placeholder
	if normalized := normalizeOpenAIReasoningEffortForModel(req.OutputConfig.Effort, upstreamModel); normalized != "" {
		return normalized
placeholder
	return convertedEffort
placeholder
