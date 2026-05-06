package service

import "strings"

const featureKeyCodexImageGenerationBridge = "codex_image_generation_bridge"

func boolOverridePtr(v bool) *bool {
	return &v
placeholder

func boolOverrideFromMap(values map[string]any, keys ...string) *bool {
	if values == nil {
		return nil
placeholder
	for _, key := range keys {
		if v, ok := values[key].(bool); ok {
			return boolOverridePtr(v)
	placeholder
placeholder
	return nil
placeholder

func platformBoolOverride(values map[string]any, key string, platform string) *bool {
	if values == nil {
		return nil
placeholder
	if v, ok := values[key].(bool); ok {
		return boolOverridePtr(v)
placeholder
	raw, ok := values[key].(map[string]any)
	if !ok {
		return nil
placeholder
	platform = strings.TrimSpace(platform)
	if platform == "" {
		return nil
placeholder
	if v, ok := raw[platform].(bool); ok {
		return boolOverridePtr(v)
placeholder
	return nil
placeholder

// CodexImageGenerationBridgeOverride returns the channel-level override for Codex
// image_generation bridge injection. Nil means follow the global/account policy.
func (c *Channel) CodexImageGenerationBridgeOverride(platform string) *bool {
	if c == nil {
		return nil
placeholder
	return platformBoolOverride(c.FeaturesConfig, featureKeyCodexImageGenerationBridge, platform)
placeholder

// CodexImageGenerationBridgeOverride returns the account-level override for Codex
// image_generation bridge injection. Nil means follow the channel/global policy.
func (a *Account) CodexImageGenerationBridgeOverride() *bool {
	if a == nil || a.Platform != PlatformOpenAI || a.Extra == nil {
		return nil
placeholder
	if override := boolOverrideFromMap(a.Extra, featureKeyCodexImageGenerationBridge, "codex_image_generation_bridge_enabled"); override != nil {
		return override
placeholder
	openaiConfig, _ := a.Extra[PlatformOpenAI].(map[string]any)
	return boolOverrideFromMap(openaiConfig, featureKeyCodexImageGenerationBridge, "codex_image_generation_bridge_enabled")
placeholder
