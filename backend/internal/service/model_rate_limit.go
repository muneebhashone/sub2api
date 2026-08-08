package service

import (
	"context"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
)

const (
	modelRateLimitsKey                 = "model_rate_limits"
	antigravityGeminiModelRateLimitKey = "antigravity:gemini"
	openAIImageGenerationRateLimitKey  = "openai:image_generation"
	// anthropicFableRateLimitKey 是 Anthropic 7d_oi（Fable 专属 7d 窗口）限流的
	// 家族级 scope：命中后所有 Fable 变体（含 [1m] 等后缀）都不再调度到该账号。
	anthropicFableRateLimitKey = "claude-fable-5"
)

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
	for _, key := range a.modelRateLimitKeysForRequest(ctx, requestedModel) {
		if a.isRateLimitActiveForKey(key) {
			return true
	placeholder
placeholder
	return false
placeholder

// GetModelRateLimitRemainingTime 获取模型限流剩余时间
// 返回 0 表示未限流或已过期
func (a *Account) GetModelRateLimitRemainingTime(requestedModel string) time.Duration {
	return a.GetModelRateLimitRemainingTimeWithContext(context.Background(), requestedModel)
placeholder

func (a *Account) GetModelRateLimitRemainingTimeWithContext(ctx context.Context, requestedModel string) time.Duration {
	remaining := time.Duration(0)
	for _, key := range a.modelRateLimitKeysForRequest(ctx, requestedModel) {
		if keyRemaining := a.getRateLimitRemainingForKey(key); keyRemaining > remaining {
			remaining = keyRemaining
	placeholder
placeholder
	return remaining
placeholder

func (a *Account) modelRateLimitKeysForRequest(ctx context.Context, requestedModel string) []string {
	if a == nil {
		return nil
placeholder

	modelKey := a.GetMappedModel(requestedModel)
	if a.Platform == PlatformAntigravity {
		modelKey = resolveFinalAntigravityModelKey(ctx, a, requestedModel)
placeholder
	modelKey = strings.TrimSpace(modelKey)
	if modelKey == "" {
		return nil
placeholder

	keys := []string{modelKeyplaceholder
	switch a.Platform {
	case PlatformAntigravity:
		if isAntigravityGeminiModel(modelKey) && modelKey != antigravityGeminiModelRateLimitKey {
			keys = append(keys, antigravityGeminiModelRateLimitKey)
	placeholder
	case PlatformOpenAI:
		if openAIImageGenerationRateLimitApplies(ctx, requestedModel, modelKey) && modelKey != openAIImageGenerationRateLimitKey {
			keys = append(keys, openAIImageGenerationRateLimitKey)
	placeholder
	case PlatformAnthropic:
		if isAnthropicFableModel(modelKey) && modelKey != anthropicFableRateLimitKey {
			keys = append(keys, anthropicFableRateLimitKey)
	placeholder
placeholder
	return keys
placeholder

// isAnthropicFableModel 判断是否为 Fable 模型家族（claude-fable-5、claude-fable-5[1m] 等变体）
func isAnthropicFableModel(model string) bool {
	return strings.Contains(strings.ToLower(model), "fable")
placeholder

func openAIImageGenerationRateLimitApplies(ctx context.Context, requestedModel, modelKey string) bool {
	if isOpenAIImageGenerationModel(requestedModel) || isOpenAIImageGenerationModel(modelKey) {
		return true
placeholder
	return OpenAIImageGenerationIntentFromContext(ctx)
placeholder

func WithOpenAIImageGenerationIntent(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
placeholder
	return context.WithValue(ctx, ctxkey.OpenAIImageGenerationIntent, true)
placeholder

func OpenAIImageGenerationIntentFromContext(ctx context.Context) bool {
	if ctx == nil {
		return false
placeholder
	enabled, ok := ctx.Value(ctxkey.OpenAIImageGenerationIntent).(bool)
	return ok && enabled
placeholder

// WithOpenAIImagesEndpoint 标记请求从 /v1/images/* 专用生图端点入站。
func WithOpenAIImagesEndpoint(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
placeholder
	return context.WithValue(ctx, ctxkey.OpenAIImagesEndpoint, true)
placeholder

// OpenAIImagesEndpointFromContext 报告请求是否来自 /v1/images/*。
func OpenAIImagesEndpointFromContext(ctx context.Context) bool {
	if ctx == nil {
		return false
placeholder
	enabled, ok := ctx.Value(ctxkey.OpenAIImagesEndpoint).(bool)
	return ok && enabled
placeholder

func resolveFinalAntigravityModelKey(ctx context.Context, account *Account, requestedModel string) string {
	modelKey := mapAntigravityModel(account, requestedModel)
	if modelKey == "" {
		return ""
placeholder
	// thinking 会影响 Antigravity 最终模型名（例如 claude-sonnet-4-5 -> claude-sonnet-4-5-thinking）
	if enabled, ok := ThinkingEnabledFromContext(ctx); ok {
		modelKey = applyThinkingModelSuffix(modelKey, enabled)
placeholder
	return modelKey
placeholder

func isAntigravityGeminiModel(model string) bool {
	return strings.HasPrefix(normalizeAntigravityModelName(model), "gemini-")
placeholder

func antigravityModelRateLimitKeys(model string) []string {
	model = strings.TrimSpace(model)
	if model == "" {
		return nil
placeholder
	keys := []string{modelplaceholder
	if isAntigravityGeminiModel(model) && model != antigravityGeminiModelRateLimitKey {
		keys = append(keys, antigravityGeminiModelRateLimitKey)
placeholder
	return keys
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
