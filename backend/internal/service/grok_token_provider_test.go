//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/stretchr/testify/require"
)

type grokTokenCacheForProviderTest struct {
	token        string
	setKey       string
	setToken     string
	setTTL       time.Duration
	lockResult   bool
	releaseCalls int
placeholder

func (c *grokTokenCacheForProviderTest) GetAccessToken(context.Context, string) (string, error) {
	if c.token == "" {
		return "", errors.New("not cached")
placeholder
	return c.token, nil
placeholder

func (c *grokTokenCacheForProviderTest) SetAccessToken(_ context.Context, key string, token string, ttl time.Duration) error {
	c.setKey = key
	c.setToken = token
	c.setTTL = ttl
	return nil
placeholder

func (c *grokTokenCacheForProviderTest) DeleteAccessToken(context.Context, string) error {
	return nil
placeholder

func (c *grokTokenCacheForProviderTest) AcquireRefreshLock(context.Context, string, time.Duration) (bool, error) {
	return c.lockResult, nil
placeholder

func (c *grokTokenCacheForProviderTest) ReleaseRefreshLock(context.Context, string) error {
	c.releaseCalls++
	return nil
placeholder

func TestGrokTokenProviderRefreshesExpiredTokenOnRequestPath(t *testing.T) {
	t.Setenv(xai.EnvBaseURL, xai.DefaultCLIBaseURL)

	expiredAt := time.Now().Add(-time.Minute).UTC().Format(time.RFC3339)
	account := &Account{
		ID:       54,
		Platform: PlatformGrok,
		Type:     AccountTypeOAuth,
placeholder
			"access_token":  "expired-access-token",
			"refresh_token": "refresh-token",
			"expires_at":    expiredAt,
			"base_url":      xai.DefaultCLIBaseURL,
			"client_id":     "client-id",
	placeholder,
placeholder
	repo := &tokenRefreshAccountRepo{placeholder
	repo.accountsByID = map[int64]*Account{54: accountplaceholder
	cache := &grokTokenCacheForProviderTest{lockResult: trueplaceholder
	oauthSvc := NewGrokOAuthService(nil, &grokOAuthClientStub{
		refreshResponse: &xai.TokenResponse{
			AccessToken: "new-access-token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
	placeholder,
placeholder)
	defer oauthSvc.Stop()

	provider := NewGrokTokenProvider(repo, cache)
	provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), NewGrokTokenRefresher(oauthSvc))

	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	require.Equal(t, "new-access-token", token)
	require.Equal(t, 1, repo.updateCredentialsCalls)
	require.Equal(t, "new-access-token", repo.accountsByID[54].GetGrokAccessToken())
	require.Equal(t, "refresh-token", repo.accountsByID[54].GetGrokRefreshToken())
	require.Equal(t, xai.DefaultCLIBaseURL, repo.accountsByID[54].GetGrokBaseURL())
	require.Equal(t, "grok:account:54", cache.setKey)
	require.Equal(t, "new-access-token", cache.setToken)
	require.Greater(t, cache.setTTL, time.Duration(0))
	require.Equal(t, 1, cache.releaseCalls)
placeholder

func TestGrokTokenProviderRefreshFailureUnschedulesWithRedactedReason(t *testing.T) {
	expiredAt := time.Now().Add(-time.Minute).UTC().Format(time.RFC3339)
	account := &Account{
		ID:       55,
		Platform: PlatformGrok,
		Type:     AccountTypeOAuth,
placeholder
			"access_token":  "expired-access-token",
			"refresh_token": "refresh-token",
			"expires_at":    expiredAt,
			"base_url":      xai.DefaultCLIBaseURL,
	placeholder,
placeholder
	repo := &tokenRefreshAccountRepo{placeholder
	repo.accountsByID = map[int64]*Account{55: accountplaceholder
	cache := &grokTokenCacheForProviderTest{lockResult: trueplaceholder
	tempCache := &tempUnschedCacheStub{placeholder
	provider := NewGrokTokenProvider(repo, cache)
	provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), &tokenRefresherStub{
		err: errors.New("temporary refresh failure access_token=leaked-access refresh_token=leaked-refresh"),
placeholder)
	provider.SetTempUnschedCache(tempCache)

	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	require.Empty(t, token)
	require.Equal(t, 1, repo.setTempUnschedCalls)
	require.Equal(t, 0, repo.setErrorCalls)
	require.Contains(t, repo.lastTempUnschedReason, "access_token=***")
	require.Contains(t, repo.lastTempUnschedReason, "refresh_token=***")
	require.NotContains(t, repo.lastTempUnschedReason, "leaked-access")
	require.NotContains(t, repo.lastTempUnschedReason, "leaked-refresh")
	require.Equal(t, 1, tempCache.setCalls)
	require.NotNil(t, tempCache.lastState)
	require.NotContains(t, tempCache.lastState.ErrorMessage, "leaked-access")
	require.NotContains(t, tempCache.lastState.ErrorMessage, "leaked-refresh")
placeholder
