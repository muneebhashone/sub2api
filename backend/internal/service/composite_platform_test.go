package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	"github.com/stretchr/testify/require"
)

func TestDetectModelPlatform(t *testing.T) {
	tests := []struct {
		name     string
		model    string
		platform string
		ok       bool
placeholder{
		{name: "claude", model: "claude-sonnet-4-5", platform: PlatformAnthropic, ok: trueplaceholder,
		{name: "anthropic prefix", model: "anthropic/claude-opus-4-5", platform: PlatformAnthropic, ok: trueplaceholder,
		{name: "gpt", model: "gpt-5.1", platform: PlatformOpenAI, ok: trueplaceholder,
		{name: "o series", model: "o3-mini", platform: PlatformOpenAI, ok: trueplaceholder,
		{name: "embedding", model: "text-embedding-3-large", platform: PlatformOpenAI, ok: trueplaceholder,
		{name: "gemini", model: "gemini-3-pro", platform: PlatformGemini, ok: trueplaceholder,
		{name: "gemini models prefix", model: "models/gemini-2.5-flash", platform: PlatformGemini, ok: trueplaceholder,
		{name: "learnlm", model: "learnlm-2.0-flash-experimental", platform: PlatformGemini, ok: trueplaceholder,
		{name: "grok", model: "grok-4", platform: PlatformGrok, ok: trueplaceholder,
		{name: "xai prefix", model: "xai/grok-4", platform: PlatformGrok, ok: trueplaceholder,
		{name: "kimi", model: "kimi-k2-thinking", platform: PlatformKimi, ok: trueplaceholder,
		{name: "moonshot prefix", model: "moonshot/moonshot-v1-32k", platform: PlatformKimi, ok: trueplaceholder,
		{name: "zhipu", model: "glm-5.2", platform: PlatformZhipu, ok: trueplaceholder,
		{name: "deepseek", model: "deepseek-v4-pro", platform: PlatformDeepseek, ok: trueplaceholder,
		{name: "unknown", model: "llama-4-maverick", ok: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			platform, ok := DetectModelPlatform(tt.model)
			require.Equal(t, tt.ok, ok)
			require.Equal(t, tt.platform, platform)
	placeholder)
placeholder
placeholder

func TestQuotaPlatformCompositeUsesResolvedOrForceOnly(t *testing.T) {
	apiKey := &APIKey{Group: &Group{Platform: PlatformCompositeplaceholderplaceholder

	require.Equal(t, "", QuotaPlatform(context.Background(), apiKey))
	require.Equal(t, PlatformGemini, QuotaPlatform(WithResolvedTargetPlatform(context.Background(), PlatformGemini), apiKey))
	require.Equal(t, PlatformAntigravity, QuotaPlatform(context.WithValue(context.Background(), ctxkey.ForcePlatform, PlatformAntigravity), apiKey))

	ctx := WithResolvedTargetPlatform(context.Background(), PlatformAnthropic)
	ctx = context.WithValue(ctx, ctxkey.ForcePlatform, PlatformAntigravity)
	require.Equal(t, PlatformAntigravity, QuotaPlatform(ctx, apiKey))
placeholder

func TestCompositeGroupSchedulerHasAllCanonicalPlatformBuckets(t *testing.T) {
	seen := make(map[string]struct{placeholder)
	for _, bucket := range schedulerCanonicalBuckets(99) {
		seen[bucket.Platform] = struct{placeholder{placeholder
placeholder
	platforms := make([]string, 0, len(seen))
	for platform := range seen {
		platforms = append(platforms, platform)
placeholder
	require.ElementsMatch(t,
		[]string{PlatformAnthropic, PlatformGemini, PlatformOpenAI, PlatformAntigravity, PlatformGrok, PlatformKimi, PlatformZhipu, PlatformDeepseekplaceholder,
		platforms,
	)
placeholder

func TestCompositeConcretePlatformsIncludeCNProviders(t *testing.T) {
	for _, platform := range []string{PlatformKimi, PlatformZhipu, PlatformDeepseekplaceholder {
		require.True(t, isConcreteRequestPlatform(platform))
		require.True(t, canCopyAccountsFromGroupPlatform(PlatformComposite, platform))
placeholder
placeholder
