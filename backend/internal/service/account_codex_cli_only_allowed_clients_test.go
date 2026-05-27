package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccount_GetCodexCLIOnlyAllowedClients(t *testing.T) {
	t.Run("OAuth 账号读取 []any 字符串列表", func(t *testing.T) {
		account := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeOAuth,
			Extra:    map[string]any{"codex_cli_only_allowed_clients": []any{"claude_code"placeholderplaceholder,
	placeholder
		require.Equal(t, []string{"claude_code"placeholder, account.GetCodexCLIOnlyAllowedClients())
placeholder)

	t.Run("OAuth 账号读取 []string 列表", func(t *testing.T) {
		account := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeOAuth,
			Extra:    map[string]any{"codex_cli_only_allowed_clients": []string{"claude_code"placeholderplaceholder,
	placeholder
		require.Equal(t, []string{"claude_code"placeholder, account.GetCodexCLIOnlyAllowedClients())
placeholder)

	t.Run("[]string 跳过空白元素", func(t *testing.T) {
		account := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeOAuth,
			Extra:    map[string]any{"codex_cli_only_allowed_clients": []string{"claude_code", "", "  "placeholderplaceholder,
	placeholder
		require.Equal(t, []string{"claude_code"placeholder, account.GetCodexCLIOnlyAllowedClients())
placeholder)

	t.Run("跳过非字符串与空白元素", func(t *testing.T) {
		account := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeOAuth,
			Extra:    map[string]any{"codex_cli_only_allowed_clients": []any{"claude_code", 123, "", "  "placeholderplaceholder,
	placeholder
		require.Equal(t, []string{"claude_code"placeholder, account.GetCodexCLIOnlyAllowedClients())
placeholder)

	t.Run("非 OAuth 账号返回空", func(t *testing.T) {
		account := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeAPIKey,
			Extra:    map[string]any{"codex_cli_only_allowed_clients": []any{"claude_code"placeholderplaceholder,
	placeholder
		require.Empty(t, account.GetCodexCLIOnlyAllowedClients())
placeholder)

	t.Run("Extra 为空返回空", func(t *testing.T) {
		account := &Account{Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder
		require.Empty(t, account.GetCodexCLIOnlyAllowedClients())
placeholder)

	t.Run("字段缺失返回空", func(t *testing.T) {
		account := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeOAuth,
			Extra:    map[string]any{placeholder,
	placeholder
		require.Empty(t, account.GetCodexCLIOnlyAllowedClients())
placeholder)
placeholder
