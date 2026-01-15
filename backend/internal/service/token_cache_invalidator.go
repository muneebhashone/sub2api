package service

import "context"

type TokenCacheInvalidator interface {
	InvalidateToken(ctx context.Context, account *Account) error
placeholder

type CompositeTokenCacheInvalidator struct {
	cache GeminiTokenCache // 统一使用一个缓存接口，通过缓存键前缀区分平台
placeholder

func NewCompositeTokenCacheInvalidator(cache GeminiTokenCache) *CompositeTokenCacheInvalidator {
	return &CompositeTokenCacheInvalidator{
		cache: cache,
placeholder
placeholder

func (c *CompositeTokenCacheInvalidator) InvalidateToken(ctx context.Context, account *Account) error {
	if c == nil || c.cache == nil || account == nil {
		return nil
placeholder
	if account.Type != AccountTypeOAuth {
		return nil
placeholder

	var cacheKey string
	switch account.Platform {
	case PlatformGemini:
		cacheKey = GeminiTokenCacheKey(account)
	case PlatformAntigravity:
		cacheKey = AntigravityTokenCacheKey(account)
	case PlatformOpenAI:
		cacheKey = OpenAITokenCacheKey(account)
	case PlatformAnthropic:
		cacheKey = ClaudeTokenCacheKey(account)
	default:
		return nil
placeholder
	return c.cache.DeleteAccessToken(ctx, cacheKey)
placeholder
