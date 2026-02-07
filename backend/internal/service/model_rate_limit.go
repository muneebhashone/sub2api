package service

import (
	"context"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
)

const modelRateLimitsKey = "model_rate_limits"

// isRateLimitActiveForKey 检查指定 key 的限流是否生效
func (a *Account) isRateLimitActiveForKey(key string) bool {
	resetAt := a.modelRateLimitResetAt(key)
	return resetAt != nil && time.Now().Before(*resetAt)
placeholder

// getRateLimitRemainingForKey 获取指定 key 的限流剩余时间，0 表示未限流或已过期
func (a *Account) getRateLimitRemainingForKey(key string) time.Duration {
	resetAt := a.modelRateLimitResetAt(key)
	if resetAt == nil {
		return 0
placeholder
	remaining := time.Until(*resetAt)
	if remaining > 0 {
		return remaining
placeholder
	return 0
placeholder

func (a *Account) isModelRateLimitedWithContext(ctx context.Context, requestedModel string) bool {
	if a == nil {
		return false
placeholder

	modelKey := a.GetMappedModel(requestedModel)
	if a.Platform == PlatformAntigravity {
		modelKey = resolveFinalAntigravityModelKey(ctx, a, requestedModel)
placeholder
	modelKey = strings.TrimSpace(modelKey)
	if modelKey == "" {
		return false
placeholder
	return a.isRateLimitActiveForKey(modelKey)
placeholder

// GetModelRateLimitRemainingTime 获取模型限流剩余时间
// 返回 0 表示未限流或已过期
func (a *Account) GetModelRateLimitRemainingTime(requestedModel string) time.Duration {
	return a.GetModelRateLimitRemainingTimeWithContext(context.Background(), requestedModel)
placeholder

func (a *Account) GetModelRateLimitRemainingTimeWithContext(ctx context.Context, requestedModel string) time.Duration {
	if a == nil {
		return 0
placeholder

	modelKey := a.GetMappedModel(requestedModel)
	if a.Platform == PlatformAntigravity {
		modelKey = resolveFinalAntigravityModelKey(ctx, a, requestedModel)
placeholder
	modelKey = strings.TrimSpace(modelKey)
	if modelKey == "" {
		return 0
placeholder
	return a.getRateLimitRemainingForKey(modelKey)
placeholder

func resolveFinalAntigravityModelKey(ctx context.Context, account *Account, requestedModel string) string {
	modelKey := mapAntigravityModel(account, requestedModel)
	if modelKey == "" {
		return ""
placeholder
	// thinking 会影响 Antigravity 最终模型名（例如 claude-sonnet-4-5 -> claude-sonnet-4-5-thinking）
	if enabled, ok := ctx.Value(ctxkey.ThinkingEnabled).(bool); ok {
		modelKey = applyThinkingModelSuffix(modelKey, enabled)
placeholder
	return modelKey
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
