package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccount_IsOpenAIPassthroughEnabled(t *testing.T) {
	t.Run("新字段开启", func(t *testing.T) {
		account := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeAPIKey,
			Extra: map[string]any{
				"openai_passthrough": true,
		placeholder,
	placeholder
		require.True(t, account.IsOpenAIPassthroughEnabled())
placeholder)

	t.Run("兼容旧字段", func(t *testing.T) {
		account := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeOAuth,
			Extra: map[string]any{
				"openai_oauth_passthrough": true,
		placeholder,
	placeholder
		require.True(t, account.IsOpenAIPassthroughEnabled())
placeholder)

	t.Run("非OpenAI账号始终关闭", func(t *testing.T) {
		account := &Account{
			Platform: PlatformAnthropic,
			Type:     AccountTypeOAuth,
			Extra: map[string]any{
				"openai_passthrough": true,
		placeholder,
	placeholder
		require.False(t, account.IsOpenAIPassthroughEnabled())
placeholder)

	t.Run("空额外配置默认关闭", func(t *testing.T) {
		account := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeOAuth,
	placeholder
		require.False(t, account.IsOpenAIPassthroughEnabled())
placeholder)
placeholder

func TestAccount_IsOpenAIOAuthPassthroughEnabled(t *testing.T) {
	t.Run("仅OAuth类型允许返回开启", func(t *testing.T) {
		oauthAccount := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeOAuth,
			Extra: map[string]any{
				"openai_passthrough": true,
		placeholder,
	placeholder
		require.True(t, oauthAccount.IsOpenAIOAuthPassthroughEnabled())

		apiKeyAccount := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeAPIKey,
			Extra: map[string]any{
				"openai_passthrough": true,
		placeholder,
	placeholder
		require.False(t, apiKeyAccount.IsOpenAIOAuthPassthroughEnabled())
placeholder)
placeholder

func TestAccount_IsCodexCLIOnlyEnabled(t *testing.T) {
	t.Run("OpenAI OAuth 开启", func(t *testing.T) {
		account := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeOAuth,
			Extra: map[string]any{
				"codex_cli_only": true,
		placeholder,
	placeholder
		require.True(t, account.IsCodexCLIOnlyEnabled())
placeholder)

	t.Run("OpenAI OAuth 关闭", func(t *testing.T) {
		account := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeOAuth,
			Extra: map[string]any{
				"codex_cli_only": false,
		placeholder,
	placeholder
		require.False(t, account.IsCodexCLIOnlyEnabled())
placeholder)

	t.Run("字段缺失默认关闭", func(t *testing.T) {
		account := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeOAuth,
			Extra:    map[string]any{placeholder,
	placeholder
		require.False(t, account.IsCodexCLIOnlyEnabled())
placeholder)

	t.Run("类型非法默认关闭", func(t *testing.T) {
		account := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeOAuth,
			Extra: map[string]any{
				"codex_cli_only": "true",
		placeholder,
	placeholder
		require.False(t, account.IsCodexCLIOnlyEnabled())
placeholder)

	t.Run("非 OAuth 账号始终关闭", func(t *testing.T) {
		apiKeyAccount := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeAPIKey,
			Extra: map[string]any{
				"codex_cli_only": true,
		placeholder,
	placeholder
		require.False(t, apiKeyAccount.IsCodexCLIOnlyEnabled())

		otherPlatform := &Account{
			Platform: PlatformAnthropic,
			Type:     AccountTypeOAuth,
			Extra: map[string]any{
				"codex_cli_only": true,
		placeholder,
	placeholder
		require.False(t, otherPlatform.IsCodexCLIOnlyEnabled())
placeholder)
placeholder
