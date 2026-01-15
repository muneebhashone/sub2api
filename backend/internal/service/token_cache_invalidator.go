package service

import "context"

type TokenCacheInvalidator interface {
	InvalidateToken(ctx context.Context, account *Account) error
placeholder

type CompositeTokenCacheInvalidator struct {
	geminiCache GeminiTokenCache
placeholder

func NewCompositeTokenCacheInvalidator(geminiCache GeminiTokenCache) *CompositeTokenCacheInvalidator {
	return &CompositeTokenCacheInvalidator{
		geminiCache: geminiCache,
placeholder
placeholder

func (c *CompositeTokenCacheInvalidator) InvalidateToken(ctx context.Context, account *Account) error {
	if c == nil || c.geminiCache == nil || account == nil {
		return nil
placeholder
	if account.Type != AccountTypeOAuth {
		return nil
placeholder

	switch account.Platform {
	case PlatformGemini:
		return c.geminiCache.DeleteAccessToken(ctx, GeminiTokenCacheKey(account))
	case PlatformAntigravity:
		return c.geminiCache.DeleteAccessToken(ctx, AntigravityTokenCacheKey(account))
	default:
		return nil
placeholder
placeholder
