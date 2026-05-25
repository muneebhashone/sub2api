//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
)

func TestPlatformFromAPIKey_NilSafe(t *testing.T) {
	if got := PlatformFromAPIKey(nil); got != "" {
		t.Errorf("nil APIKey should yield empty string, got %q", got)
placeholder
placeholder

func TestPlatformFromAPIKey_NilGroup(t *testing.T) {
	k := &APIKey{Group: nilplaceholder
	if got := PlatformFromAPIKey(k); got != "" {
		t.Errorf("APIKey with nil Group should yield empty string, got %q", got)
placeholder
placeholder

func TestPlatformFromAPIKey_DerivesFromGroup(t *testing.T) {
	tests := []struct {
		name     string
		platform string
placeholder{
		{"anthropic", "anthropic"placeholder,
		{"openai", "openai"placeholder,
		{"gemini", "gemini"placeholder,
		{"antigravity", "antigravity"placeholder,
		{"empty", ""placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &APIKey{
				Group: &Group{Platform: tt.platformplaceholder,
		placeholder
			got := PlatformFromAPIKey(k)
			if got != tt.platform {
				t.Errorf("PlatformFromAPIKey(%q) = %q, want %q", tt.platform, got, tt.platform)
		placeholder
	placeholder)
placeholder
placeholder

// TestQuotaPlatform 锁定配额计量口径：ForcePlatform 路由（如 /antigravity）按 ForcePlatform 计，
// 否则回退到 Group 平台。preflight 与 post-billing 共用此口径，保证一致。
func TestQuotaPlatform(t *testing.T) {
	apiKey := &APIKey{Group: &Group{Platform: PlatformAnthropicplaceholderplaceholder

	t.Run("no force platform falls back to group platform", func(t *testing.T) {
		if got := QuotaPlatform(context.Background(), apiKey); got != PlatformAnthropic {
			t.Errorf("QuotaPlatform without force = %q, want %q", got, PlatformAnthropic)
	placeholder
placeholder)

	t.Run("force platform overrides group platform", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ctxkey.ForcePlatform, PlatformAntigravity)
		if got := QuotaPlatform(ctx, apiKey); got != PlatformAntigravity {
			t.Errorf("QuotaPlatform with force = %q, want %q", got, PlatformAntigravity)
	placeholder
placeholder)

	t.Run("empty force platform falls back to group platform", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ctxkey.ForcePlatform, "")
		if got := QuotaPlatform(ctx, apiKey); got != PlatformAnthropic {
			t.Errorf("QuotaPlatform with empty force = %q, want %q", got, PlatformAnthropic)
	placeholder
placeholder)

	t.Run("nil api key with force platform returns force platform", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ctxkey.ForcePlatform, PlatformAntigravity)
		if got := QuotaPlatform(ctx, nil); got != PlatformAntigravity {
			t.Errorf("QuotaPlatform(nil) with force = %q, want %q", got, PlatformAntigravity)
	placeholder
placeholder)
placeholder
