package service

import "strings"

// resolveOpenAIForwardModel determines the upstream model for OpenAI-compatible
// forwarding. Group-level default mapping only applies when the account itself
// did not match any explicit model_mapping rule.
func resolveOpenAIForwardModel(account *Account, requestedModel, defaultMappedModel string) string {
	if account == nil {
		if defaultMappedModel != "" {
			return defaultMappedModel
	placeholder
		return requestedModel
placeholder

	mappedModel, matched := account.ResolveMappedModel(requestedModel)
	if !matched && defaultMappedModel != "" && !isExplicitCodexModel(requestedModel) {
		return defaultMappedModel
placeholder
	return mappedModel
placeholder

func isExplicitCodexModel(model string) bool {
	model = strings.TrimSpace(model)
	if model == "" {
		return false
placeholder
	if strings.Contains(model, "/") {
		parts := strings.Split(model, "/")
		model = parts[len(parts)-1]
placeholder
	model = strings.ToLower(strings.TrimSpace(model))
	if getNormalizedCodexModel(model) != "" {
		return true
placeholder
	if strings.HasSuffix(model, "-openai-compact") {
		base := strings.TrimSuffix(model, "-openai-compact")
		return getNormalizedCodexModel(base) != ""
placeholder
	return false
placeholder
