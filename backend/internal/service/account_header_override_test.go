//go:build unit

package service

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func headerOverrideTestAccount(platform, accountType string, credentials map[string]any) *Account {
placeholder
		Platform:    platform,
		Type:        accountType,
		Credentials: credentials,
placeholder
placeholder

func TestIsHeaderOverrideEligible(t *testing.T) {
	tests := []struct {
		name     string
		platform string
		accType  string
		want     bool
placeholder{
		{"anthropic apikey", PlatformAnthropic, AccountTypeAPIKey, trueplaceholder,
		{"openai apikey", PlatformOpenAI, AccountTypeAPIKey, trueplaceholder,
		{"kimi apikey", PlatformKimi, AccountTypeAPIKey, trueplaceholder,
		{"zhipu apikey", PlatformZhipu, AccountTypeAPIKey, trueplaceholder,
		{"deepseek apikey", PlatformDeepseek, AccountTypeAPIKey, trueplaceholder,
		{"anthropic oauth", PlatformAnthropic, AccountTypeOAuth, falseplaceholder,
		{"openai oauth", PlatformOpenAI, AccountTypeOAuth, falseplaceholder,
		{"kimi oauth", PlatformKimi, AccountTypeOAuth, falseplaceholder,
		{"zhipu oauth", PlatformZhipu, AccountTypeOAuth, falseplaceholder,
		{"deepseek oauth", PlatformDeepseek, AccountTypeOAuth, falseplaceholder,
		{"gemini apikey", PlatformGemini, AccountTypeAPIKey, falseplaceholder,
		{"grok apikey", PlatformGrok, AccountTypeAPIKey, trueplaceholder,
		{"grok oauth", PlatformGrok, AccountTypeOAuth, trueplaceholder,
		{"antigravity apikey", PlatformAntigravity, AccountTypeAPIKey, falseplaceholder,
		{"anthropic bedrock", PlatformAnthropic, AccountTypeBedrock, falseplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := headerOverrideTestAccount(tt.platform, tt.accType, nil)
			require.Equal(t, tt.want, acc.IsHeaderOverrideEligible())
	placeholder)
placeholder

	var nilAccount *Account
	require.False(t, nilAccount.IsHeaderOverrideEligible())
	require.False(t, nilAccount.IsHeaderOverrideEnabled())
	require.Nil(t, nilAccount.GetHeaderOverrides())
placeholder

func TestIsHeaderOverrideEnabled(t *testing.T) {
	acc := headerOverrideTestAccount(PlatformAnthropic, AccountTypeAPIKey, map[string]any{
		credKeyHeaderOverrideEnabled: true,
placeholder)
	require.True(t, acc.IsHeaderOverrideEnabled())

	// 未配置 / 非 bool / false 均视为未启用
	require.False(t, headerOverrideTestAccount(PlatformAnthropic, AccountTypeAPIKey, nil).IsHeaderOverrideEnabled())
	require.False(t, headerOverrideTestAccount(PlatformAnthropic, AccountTypeAPIKey, map[string]any{
		credKeyHeaderOverrideEnabled: "true",
placeholder).IsHeaderOverrideEnabled())
	require.False(t, headerOverrideTestAccount(PlatformAnthropic, AccountTypeAPIKey, map[string]any{
		credKeyHeaderOverrideEnabled: false,
placeholder).IsHeaderOverrideEnabled())

	// 不符合平台/类型条件时即使配置了 true 也不启用
	require.False(t, headerOverrideTestAccount(PlatformAnthropic, AccountTypeOAuth, map[string]any{
		credKeyHeaderOverrideEnabled: true,
placeholder).IsHeaderOverrideEnabled())
	require.False(t, headerOverrideTestAccount(PlatformGemini, AccountTypeAPIKey, map[string]any{
		credKeyHeaderOverrideEnabled: true,
placeholder).IsHeaderOverrideEnabled())
placeholder

func TestGetHeaderOverrides(t *testing.T) {
	acc := headerOverrideTestAccount(PlatformOpenAI, AccountTypeAPIKey, map[string]any{
		credKeyHeaderOverrideEnabled: true,
		credKeyHeaderOverrides: map[string]any{
			"User-Agent":    "my-agent/1.0",  // 大写 key 归一化为小写
			" X-App ":       "cli",           // 名称去空白
			"x-empty":       "",              // 空 value（模板占位）跳过
			"authorization": "Bearer leaked", // 禁止覆写的头跳过
			"bad name":      "value",         // 非法 header 名跳过
			"x-padded":      "  padded  ",    // value 去空白
	placeholder,
placeholder)
	overrides := acc.GetHeaderOverrides()
	require.Equal(t, map[string]string{
		"user-agent": "my-agent/1.0",
		"x-app":      "cli",
		"x-padded":   "padded",
placeholder, overrides)

	// 未启用时返回 nil
	disabled := headerOverrideTestAccount(PlatformOpenAI, AccountTypeAPIKey, map[string]any{
		credKeyHeaderOverrides: map[string]any{"user-agent": "x"placeholder,
placeholder)
	require.Nil(t, disabled.GetHeaderOverrides())

	// 启用但全部为空 value 时返回 nil
	empty := headerOverrideTestAccount(PlatformOpenAI, AccountTypeAPIKey, map[string]any{
		credKeyHeaderOverrideEnabled: true,
		credKeyHeaderOverrides:       map[string]any{"user-agent": ""placeholder,
placeholder)
	require.Nil(t, empty.GetHeaderOverrides())

	// 未经 Normalize 落库的超长数据 / WebSocket 握手头在应用时被防御性跳过
	oversizedValue := strings.Repeat("a", maxHeaderOverrideValueLength+1)
	defensive := headerOverrideTestAccount(PlatformOpenAI, AccountTypeAPIKey, map[string]any{
		credKeyHeaderOverrideEnabled: true,
		credKeyHeaderOverrides: map[string]any{
			"x-big":                    oversizedValue,
			"sec-websocket-key":        "forged",
			"content-type":             "application/json", // 名单扩充前落库的数据也要被拦截
			"x-claude-code-session-id": "pinned-session",
			"x-ok":                     "ok",
	placeholder,
placeholder)
	require.Equal(t, map[string]string{"x-ok": "ok"placeholder, defensive.GetHeaderOverrides())
placeholder

func TestApplyHeaderOverrides(t *testing.T) {
	acc := headerOverrideTestAccount(PlatformAnthropic, AccountTypeAPIKey, map[string]any{
		credKeyHeaderOverrideEnabled: true,
		credKeyHeaderOverrides: map[string]any{
			"user-agent":     "override-agent/2.0",
			"anthropic-beta": "custom-beta-1",
			"x-custom":       "custom-value",
	placeholder,
placeholder)

	h := http.Header{placeholder
	// 模拟转发链路：canonical key 与 wire casing 原样 key 混合存在
	h.Set("User-Agent", "claude-cli/2.1.161 (external, cli)")
	h["anthropic-beta"] = []string{"claude-code-20250219,oauth-2025-04-20"placeholder // 非 canonical 原样 key
	h.Set("Content-Type", "application/json")

	acc.ApplyHeaderOverrides(h)

	// user-agent 覆盖且只有一个值（已知头恢复 wire casing）
	require.Equal(t, []string{"override-agent/2.0"placeholder, h["User-Agent"])
	// anthropic-beta：非 canonical 旧值被清除，写入 wire casing（小写）
	require.Equal(t, []string{"custom-beta-1"placeholder, h["anthropic-beta"])
	require.Empty(t, h["Anthropic-Beta"])
	// 新增头（未知头以小写原样键写入，与转发链路 wire casing 约定一致）
	require.Equal(t, []string{"custom-value"placeholder, h["x-custom"])
	require.Equal(t, "custom-value", getHeaderRaw(h, "x-custom"))
	// 未覆写的头不受影响
	require.Equal(t, "application/json", h.Get("Content-Type"))

	// 覆盖后不存在任何大小写重复
	count := 0
	for k := range h {
		if k == "anthropic-beta" || k == "Anthropic-Beta" {
			count++
	placeholder
placeholder
	require.Equal(t, 1, count)
placeholder

func TestApplyHeaderOverridesNoOpPaths(t *testing.T) {
	baseline := func() http.Header {
		h := http.Header{placeholder
		h.Set("User-Agent", "orig")
		return h
placeholder

	// OAuth 账号：即使配置了覆写也不生效
	oauth := headerOverrideTestAccount(PlatformAnthropic, AccountTypeOAuth, map[string]any{
		credKeyHeaderOverrideEnabled: true,
		credKeyHeaderOverrides:       map[string]any{"user-agent": "hacked"placeholder,
placeholder)
	h := baseline()
	oauth.ApplyHeaderOverrides(h)
	require.Equal(t, "orig", h.Get("User-Agent"))

	// 未启用开关
	off := headerOverrideTestAccount(PlatformAnthropic, AccountTypeAPIKey, map[string]any{
		credKeyHeaderOverrides: map[string]any{"user-agent": "hacked"placeholder,
placeholder)
	h = baseline()
	off.ApplyHeaderOverrides(h)
	require.Equal(t, "orig", h.Get("User-Agent"))

	// 禁止覆写的头（authorization / x-api-key / host 等）不会被应用
	blocked := headerOverrideTestAccount(PlatformOpenAI, AccountTypeAPIKey, map[string]any{
		credKeyHeaderOverrideEnabled: true,
		credKeyHeaderOverrides: map[string]any{
			"Authorization":  "Bearer evil",
			"X-Api-Key":      "evil",
			"Host":           "evil.example.com",
			"Content-Length": "0",
	placeholder,
placeholder)
	h = http.Header{placeholder
	h.Set("Authorization", "Bearer real-key")
	blocked.ApplyHeaderOverrides(h)
	require.Equal(t, "Bearer real-key", h.Get("Authorization"))
	require.Empty(t, h.Get("X-Api-Key"))
	require.Empty(t, h.Get("Host"))

	// nil header 不 panic
	blocked.ApplyHeaderOverrides(nil)
placeholder

func TestNormalizeHeaderOverrideCredentials(t *testing.T) {
	t.Run("nil credentials no-op", func(t *testing.T) {
		require.NoError(t, NormalizeHeaderOverrideCredentials(nil))
placeholder)

	t.Run("missing keys no-op", func(t *testing.T) {
		creds := map[string]any{"api_key": "sk-xxx"placeholder
		require.NoError(t, NormalizeHeaderOverrideCredentials(creds))
		_, exists := creds[credKeyHeaderOverrides]
		require.False(t, exists)
placeholder)

	t.Run("normalizes names and values", func(t *testing.T) {
		creds := map[string]any{
			credKeyHeaderOverrideEnabled: true,
			credKeyHeaderOverrides: map[string]any{
				" User-Agent ": " my-agent ",
				"X-App":        "",
				"":             "", // 完全空行被丢弃
		placeholder,
	placeholder
		require.NoError(t, NormalizeHeaderOverrideCredentials(creds))
		require.Equal(t, map[string]any{
			"user-agent": "my-agent",
			"x-app":      "",
	placeholder, creds[credKeyHeaderOverrides])
placeholder)

	t.Run("accepts map[string]string input", func(t *testing.T) {
		creds := map[string]any{
			credKeyHeaderOverrides: map[string]string{"X-App": "cli"placeholder,
	placeholder
		require.NoError(t, NormalizeHeaderOverrideCredentials(creds))
		require.Equal(t, map[string]any{"x-app": "cli"placeholder, creds[credKeyHeaderOverrides])
placeholder)

	t.Run("rejects non-bool enabled", func(t *testing.T) {
		err := NormalizeHeaderOverrideCredentials(map[string]any{
			credKeyHeaderOverrideEnabled: "yes",
	placeholder)
	placeholder
placeholder)

	t.Run("rejects non-object overrides", func(t *testing.T) {
		err := NormalizeHeaderOverrideCredentials(map[string]any{
			credKeyHeaderOverrides: []any{"user-agent"placeholder,
	placeholder)
	placeholder
placeholder)

	t.Run("rejects non-string value", func(t *testing.T) {
		err := NormalizeHeaderOverrideCredentials(map[string]any{
			credKeyHeaderOverrides: map[string]any{"x-app": placeholder,
	placeholder)
	placeholder
placeholder)

	t.Run("rejects invalid header name", func(t *testing.T) {
		for _, name := range []string{"bad name", "bad:name", "bad\nname", "值"placeholder {
			err := NormalizeHeaderOverrideCredentials(map[string]any{
				credKeyHeaderOverrides: map[string]any{name: "v"placeholder,
		placeholder)
			require.Error(t, err, "name %q should be rejected", name)
	placeholder
placeholder)

	t.Run("rejects empty name with value", func(t *testing.T) {
		err := NormalizeHeaderOverrideCredentials(map[string]any{
			credKeyHeaderOverrides: map[string]any{"  ": "v"placeholder,
	placeholder)
	placeholder
placeholder)

	t.Run("rejects blocked headers", func(t *testing.T) {
		for _, name := range []string{
			"Authorization", "x-api-key", "Host", "content-length", "Transfer-Encoding",
			"connection", "accept-encoding", "Sec-WebSocket-Key", "session_id",
			"conversation_id", "x-codex-turn-state", "chatgpt-account-id",
			"Content-Type", "Cookie", "x-goog-api-key",
			"X-Claude-Code-Session-Id", "x-client-request-id",
	placeholder {
			err := NormalizeHeaderOverrideCredentials(map[string]any{
				credKeyHeaderOverrides: map[string]any{name: "v"placeholder,
		placeholder)
			require.Error(t, err, "blocked header %q should be rejected", name)
	placeholder
placeholder)

	t.Run("allows tab inside value", func(t *testing.T) {
		creds := map[string]any{
			credKeyHeaderOverrides: map[string]any{"x-app": "a\tb"placeholder,
	placeholder
		require.NoError(t, NormalizeHeaderOverrideCredentials(creds))
		require.Equal(t, map[string]any{"x-app": "a\tb"placeholder, creds[credKeyHeaderOverrides])
placeholder)

	t.Run("rejects invalid value", func(t *testing.T) {
		err := NormalizeHeaderOverrideCredentials(map[string]any{
			credKeyHeaderOverrides: map[string]any{"x-app": "bad\nvalue"placeholder,
	placeholder)
	placeholder
placeholder)

	t.Run("rejects duplicate names case-insensitively", func(t *testing.T) {
		err := NormalizeHeaderOverrideCredentials(map[string]any{
			credKeyHeaderOverrides: map[string]any{
				"User-Agent": "a",
				"user-agent": "b",
		placeholder,
	placeholder)
	placeholder
placeholder)

	t.Run("rejects too many entries", func(t *testing.T) {
		entries := make(map[string]any, maxHeaderOverrideEntries+1)
		for i := 0; i <= maxHeaderOverrideEntries; i++ {
			entries["x-h-"+string(rune('a'+i%26))+string(rune('a'+(i/26)%26))+string(rune('a'+(i/676)%26))] = "v"
	placeholder
		err := NormalizeHeaderOverrideCredentials(map[string]any{
			credKeyHeaderOverrides: entries,
	placeholder)
	placeholder
placeholder)

	t.Run("rejects oversized value", func(t *testing.T) {
		big := make([]byte, maxHeaderOverrideValueLength+1)
		for i := range big {
			big[i] = 'a'
	placeholder
		err := NormalizeHeaderOverrideCredentials(map[string]any{
			credKeyHeaderOverrides: map[string]any{"x-app": string(big)placeholder,
	placeholder)
	placeholder
placeholder)
placeholder
