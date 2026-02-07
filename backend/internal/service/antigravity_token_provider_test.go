//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAntigravityTokenProvider_GetAccessToken_Upstream(t *testing.T) {
	provider := &AntigravityTokenProvider{placeholder

	t.Run("upstream account with valid api_key", func(t *testing.T) {
		account := &Account{
			Platform: PlatformAntigravity,
			Type:     AccountTypeUpstream,
	placeholder
				"api_key": "sk-test-key-12345",
		placeholder,
	placeholder
		token, err := provider.GetAccessToken(context.Background(), account)
	placeholder
		require.Equal(t, "sk-test-key-12345", token)
placeholder)

	t.Run("upstream account missing api_key", func(t *testing.T) {
		account := &Account{
			Platform:    PlatformAntigravity,
			Type:        AccountTypeUpstream,
	placeholderplaceholder,
	placeholder
		token, err := provider.GetAccessToken(context.Background(), account)
	placeholder
		require.Contains(t, err.Error(), "upstream account missing api_key")
		require.Empty(t, token)
placeholder)

	t.Run("upstream account with empty api_key", func(t *testing.T) {
		account := &Account{
			Platform: PlatformAntigravity,
			Type:     AccountTypeUpstream,
	placeholder
				"api_key": "",
		placeholder,
	placeholder
		token, err := provider.GetAccessToken(context.Background(), account)
	placeholder
		require.Contains(t, err.Error(), "upstream account missing api_key")
		require.Empty(t, token)
placeholder)

	t.Run("upstream account with nil credentials", func(t *testing.T) {
		account := &Account{
			Platform: PlatformAntigravity,
			Type:     AccountTypeUpstream,
	placeholder
		token, err := provider.GetAccessToken(context.Background(), account)
	placeholder
		require.Contains(t, err.Error(), "upstream account missing api_key")
		require.Empty(t, token)
placeholder)
placeholder

func TestAntigravityTokenProvider_GetAccessToken_Guards(t *testing.T) {
	provider := &AntigravityTokenProvider{placeholder

	t.Run("nil account", func(t *testing.T) {
		token, err := provider.GetAccessToken(context.Background(), nil)
	placeholder
		require.Contains(t, err.Error(), "account is nil")
		require.Empty(t, token)
placeholder)

	t.Run("non-antigravity platform", func(t *testing.T) {
		account := &Account{
			Platform: PlatformAnthropic,
			Type:     AccountTypeOAuth,
	placeholder
		token, err := provider.GetAccessToken(context.Background(), account)
	placeholder
		require.Contains(t, err.Error(), "not an antigravity account")
		require.Empty(t, token)
placeholder)

	t.Run("unsupported account type", func(t *testing.T) {
		account := &Account{
			Platform: PlatformAntigravity,
			Type:     AccountTypeAPIKey,
	placeholder
		token, err := provider.GetAccessToken(context.Background(), account)
	placeholder
		require.Contains(t, err.Error(), "not an antigravity oauth account")
		require.Empty(t, token)
placeholder)
placeholder
