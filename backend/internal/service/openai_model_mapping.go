package service

import "strings"

// resolveOpenAIForwardModel resolves the account/group mapping result for
// OpenAI-compatible forwarding. Group-level default mapping only applies when
// the account itself did not match any explicit model_mapping rule.
func resolveOpenAIForwardModel(account *Account, requestedModel, defaultMappedModel string) string {
	if account == nil {
		if defaultMappedModel != "" {
			return defaultMappedModel
	placeholder
		return requestedModel
placeholder

	mappedModel, matched := account.ResolveMappedModel(requestedModel)
	if !matched && defaultMappedModel != "" {
		return defaultMappedModel
placeholder
	return mappedModel
placeholder

func resolveOpenAIUpstreamModel(model string) string {
	if isBareGPT53CodexSparkModel(model) {
		return "gpt-5.3-codex-spark"
placeholder
	return normalizeCodexModel(strings.TrimSpace(model))
placeholder

func isBareGPT53CodexSparkModel(model string) bool {
	modelID := strings.TrimSpace(model)
	if modelID == "" {
		return false
placeholder
	if strings.Contains(modelID, "/") {
		parts := strings.Split(modelID, "/")
		modelID = parts[len(parts)-1]
placeholder
	normalized := strings.ToLower(strings.TrimSpace(modelID))
	return normalized == "gpt-5.3-codex-spark" || normalized == "gpt 5.3 codex spark"
placeholder
