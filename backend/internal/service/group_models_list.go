package service

import "strings"

func normalizeGroupModelsListConfig(cfg GroupModelsListConfig) GroupModelsListConfig {
	out := GroupModelsListConfig{Enabled: cfg.Enabledplaceholder
	if len(cfg.Models) == 0 {
		return out
placeholder

	seen := make(map[string]struct{placeholder, len(cfg.Models))
	out.Models = make([]string, 0, len(cfg.Models))
	for _, model := range cfg.Models {
		model = strings.TrimSpace(model)
		if model == "" {
			continue
	placeholder
		if _, ok := seen[model]; ok {
			continue
	placeholder
		seen[model] = struct{placeholder{placeholder
		out.Models = append(out.Models, model)
placeholder
	if len(out.Models) == 0 {
		out.Models = nil
placeholder
	return out
placeholder

func (g *Group) CustomModelsListEnabled() bool {
	return g != nil && g.ModelsListConfig.Enabled && len(g.ModelsListConfig.Models) > 0
placeholder
