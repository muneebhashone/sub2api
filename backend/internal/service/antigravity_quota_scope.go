package service

import (
	"strings"
	"time"
)

const antigravityQuotaScopesKey = "antigravity_quota_scopes"

// AntigravityQuotaScope 表示 Antigravity 的配额域
type AntigravityQuotaScope string

const (
	AntigravityQuotaScopeClaude      AntigravityQuotaScope = "claude"
	AntigravityQuotaScopeGeminiText  AntigravityQuotaScope = "gemini_text"
	AntigravityQuotaScopeGeminiImage AntigravityQuotaScope = "gemini_image"
)

// resolveAntigravityQuotaScope 根据模型名称解析配额域
func resolveAntigravityQuotaScope(requestedModel string) (AntigravityQuotaScope, bool) {
	model := normalizeAntigravityModelName(requestedModel)
	if model == "" {
		return "", false
placeholder
	switch {
	case strings.HasPrefix(model, "claude-"):
		return AntigravityQuotaScopeClaude, true
	case strings.HasPrefix(model, "gemini-"):
		if isImageGenerationModel(model) {
			return AntigravityQuotaScopeGeminiImage, true
	placeholder
		return AntigravityQuotaScopeGeminiText, true
	default:
		return "", false
placeholder
placeholder

func normalizeAntigravityModelName(model string) string {
	normalized := strings.ToLower(strings.TrimSpace(model))
	normalized = strings.TrimPrefix(normalized, "models/")
	return normalized
placeholder

// IsSchedulableForModel 结合 Antigravity 配额域限流判断是否可调度
func (a *Account) IsSchedulableForModel(requestedModel string) bool {
	if a == nil {
		return false
placeholder
	if !a.IsSchedulable() {
		return false
placeholder
	if a.isModelRateLimited(requestedModel) {
		return false
placeholder
	if a.Platform != PlatformAntigravity {
		return true
placeholder
	scope, ok := resolveAntigravityQuotaScope(requestedModel)
	if !ok {
		return true
placeholder
	resetAt := a.antigravityQuotaScopeResetAt(scope)
	if resetAt == nil {
		return true
placeholder
	now := time.Now()
	return !now.Before(*resetAt)
placeholder

func (a *Account) antigravityQuotaScopeResetAt(scope AntigravityQuotaScope) *time.Time {
	if a == nil || a.Extra == nil || scope == "" {
		return nil
placeholder
	rawScopes, ok := a.Extra[antigravityQuotaScopesKey].(map[string]any)
	if !ok {
		return nil
placeholder
	rawScope, ok := rawScopes[string(scope)].(map[string]any)
	if !ok {
		return nil
placeholder
	resetAtRaw, ok := rawScope["rate_limit_reset_at"].(string)
	if !ok || strings.TrimSpace(resetAtRaw) == "" {
		return nil
placeholder
	resetAt, err := time.Parse(time.RFC3339, resetAtRaw)
	if err != nil {
		return nil
placeholder
	return &resetAt
placeholder

var antigravityAllScopes = []AntigravityQuotaScope{
	AntigravityQuotaScopeClaude,
	AntigravityQuotaScopeGeminiText,
	AntigravityQuotaScopeGeminiImage,
placeholder

func (a *Account) GetAntigravityScopeRateLimits() map[string]int64 {
	if a == nil || a.Platform != PlatformAntigravity {
		return nil
placeholder
	now := time.Now()
	result := make(map[string]int64)
	for _, scope := range antigravityAllScopes {
		resetAt := a.antigravityQuotaScopeResetAt(scope)
		if resetAt != nil && now.Before(*resetAt) {
			remainingSec := int64(time.Until(*resetAt).Seconds())
			if remainingSec > 0 {
				result[string(scope)] = remainingSec
		placeholder
	placeholder
placeholder
	if len(result) == 0 {
		return nil
placeholder
	return result
placeholder
