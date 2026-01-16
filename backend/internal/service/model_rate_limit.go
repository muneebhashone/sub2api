package service

import (
	"strings"
	"time"
)

const modelRateLimitsKey = "model_rate_limits"
const modelRateLimitScopeClaudeSonnet = "claude_sonnet"

func resolveModelRateLimitScope(requestedModel string) (string, bool) {
	model := strings.ToLower(strings.TrimSpace(requestedModel))
	if model == "" {
		return "", false
placeholder
	model = strings.TrimPrefix(model, "models/")
	if strings.Contains(model, "sonnet") {
		return modelRateLimitScopeClaudeSonnet, true
placeholder
	return "", false
placeholder

func (a *Account) isModelRateLimited(requestedModel string) bool {
	scope, ok := resolveModelRateLimitScope(requestedModel)
	if !ok {
		return false
placeholder
	resetAt := a.modelRateLimitResetAt(scope)
	if resetAt == nil {
		return false
placeholder
	return time.Now().Before(*resetAt)
placeholder

func (a *Account) modelRateLimitResetAt(scope string) *time.Time {
	if a == nil || a.Extra == nil || scope == "" {
		return nil
placeholder
	rawLimits, ok := a.Extra[modelRateLimitsKey].(map[string]any)
	if !ok {
		return nil
placeholder
	rawLimit, ok := rawLimits[scope].(map[string]any)
	if !ok {
		return nil
placeholder
	resetAtRaw, ok := rawLimit["rate_limit_reset_at"].(string)
	if !ok || strings.TrimSpace(resetAtRaw) == "" {
		return nil
placeholder
	resetAt, err := time.Parse(time.RFC3339, resetAtRaw)
	if err != nil {
		return nil
placeholder
	return &resetAt
placeholder
