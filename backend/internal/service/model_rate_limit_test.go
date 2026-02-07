package service

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
)

func TestIsModelRateLimited(t *testing.T) {
	now := time.Now()
	future := now.Add(10 * time.Minute).Format(time.RFC3339)
	past := now.Add(-10 * time.Minute).Format(time.RFC3339)

	tests := []struct {
		name           string
		account        *Account
		requestedModel string
		expected       bool
placeholder{
		{
			name: "official model ID hit - claude-sonnet-4-5",
			account: &Account{
				Extra: map[string]any{
					modelRateLimitsKey: map[string]any{
						"claude-sonnet-4-5": map[string]any{
							"rate_limit_reset_at": future,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "claude-sonnet-4-5",
			expected:       true,
	placeholder,
		{
			name: "official model ID hit via mapping - request claude-3-5-sonnet, mapped to claude-sonnet-4-5",
			account: &Account{
		placeholder
					"model_mapping": map[string]any{
						"claude-3-5-sonnet": "claude-sonnet-4-5",
				placeholder,
			placeholder,
				Extra: map[string]any{
					modelRateLimitsKey: map[string]any{
						"claude-sonnet-4-5": map[string]any{
							"rate_limit_reset_at": future,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "claude-3-5-sonnet",
			expected:       true,
	placeholder,
		{
			name: "no rate limit - expired",
			account: &Account{
				Extra: map[string]any{
					modelRateLimitsKey: map[string]any{
						"claude-sonnet-4-5": map[string]any{
							"rate_limit_reset_at": past,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "claude-sonnet-4-5",
			expected:       false,
	placeholder,
		{
			name: "no rate limit - no matching key",
			account: &Account{
				Extra: map[string]any{
					modelRateLimitsKey: map[string]any{
						"gemini-3-flash": map[string]any{
							"rate_limit_reset_at": future,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "claude-sonnet-4-5",
			expected:       false,
	placeholder,
		{
			name:           "no rate limit - unsupported model",
			account:        &Account{placeholder,
			requestedModel: "gpt-4",
			expected:       false,
	placeholder,
		{
			name:           "no rate limit - empty model",
			account:        &Account{placeholder,
			requestedModel: "",
			expected:       false,
	placeholder,
		{
			name: "gemini model hit",
			account: &Account{
				Extra: map[string]any{
					modelRateLimitsKey: map[string]any{
						"gemini-3-pro-high": map[string]any{
							"rate_limit_reset_at": future,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "gemini-3-pro-high",
			expected:       true,
	placeholder,
		{
			name: "antigravity platform - gemini-3-pro-preview mapped to gemini-3-pro-high",
			account: &Account{
				Platform: PlatformAntigravity,
				Extra: map[string]any{
					modelRateLimitsKey: map[string]any{
						"gemini-3-pro-high": map[string]any{
							"rate_limit_reset_at": future,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "gemini-3-pro-preview",
			expected:       true,
	placeholder,
		{
			name: "non-antigravity platform - gemini-3-pro-preview NOT mapped",
			account: &Account{
		placeholder
				Extra: map[string]any{
					modelRateLimitsKey: map[string]any{
						"gemini-3-pro-high": map[string]any{
							"rate_limit_reset_at": future,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "gemini-3-pro-preview",
			expected:       false, // gemini 平台不走 antigravity 映射
	placeholder,
		{
			name: "antigravity platform - claude-opus-4-5-thinking mapped to opus-4-6-thinking",
			account: &Account{
				Platform: PlatformAntigravity,
				Extra: map[string]any{
					modelRateLimitsKey: map[string]any{
						"claude-opus-4-6-thinking": map[string]any{
							"rate_limit_reset_at": future,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "claude-opus-4-5-thinking",
			expected:       true,
	placeholder,
		{
			name: "no scope fallback - claude_sonnet should not match",
			account: &Account{
				Extra: map[string]any{
					modelRateLimitsKey: map[string]any{
						"claude_sonnet": map[string]any{
							"rate_limit_reset_at": future,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "claude-3-5-sonnet-20241022",
			expected:       false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.account.isModelRateLimitedWithContext(context.Background(), tt.requestedModel)
			if result != tt.expected {
				t.Errorf("isModelRateLimited(%q) = %v, want %v", tt.requestedModel, result, tt.expected)
		placeholder
	placeholder)
placeholder
placeholder

func TestIsModelRateLimited_Antigravity_ThinkingAffectsModelKey(t *testing.T) {
	now := time.Now()
	future := now.Add(10 * time.Minute).Format(time.RFC3339)

	account := &Account{
		Platform: PlatformAntigravity,
		Extra: map[string]any{
			modelRateLimitsKey: map[string]any{
				"claude-sonnet-4-5-thinking": map[string]any{
					"rate_limit_reset_at": future,
			placeholder,
		placeholder,
	placeholder,
placeholder

	ctx := context.WithValue(context.Background(), ctxkey.ThinkingEnabled, true)
	if !account.isModelRateLimitedWithContext(ctx, "claude-sonnet-4-5") {
		t.Errorf("expected model to be rate limited")
placeholder
placeholder

func TestGetModelRateLimitRemainingTime(t *testing.T) {
	now := time.Now()
	future10m := now.Add(10 * time.Minute).Format(time.RFC3339)
	future5m := now.Add(5 * time.Minute).Format(time.RFC3339)
	past := now.Add(-10 * time.Minute).Format(time.RFC3339)

	tests := []struct {
		name           string
		account        *Account
		requestedModel string
		minExpected    time.Duration
		maxExpected    time.Duration
placeholder{
		{
			name:           "nil account",
			account:        nil,
			requestedModel: "claude-sonnet-4-5",
			minExpected:    0,
			maxExpected:    0,
	placeholder,
		{
			name: "model rate limited - direct hit",
			account: &Account{
				Extra: map[string]any{
					modelRateLimitsKey: map[string]any{
						"claude-sonnet-4-5": map[string]any{
							"rate_limit_reset_at": future10m,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "claude-sonnet-4-5",
			minExpected:    9 * time.Minute,
			maxExpected:    11 * time.Minute,
	placeholder,
		{
			name: "model rate limited - via mapping",
			account: &Account{
		placeholder
					"model_mapping": map[string]any{
						"claude-3-5-sonnet": "claude-sonnet-4-5",
				placeholder,
			placeholder,
				Extra: map[string]any{
					modelRateLimitsKey: map[string]any{
						"claude-sonnet-4-5": map[string]any{
							"rate_limit_reset_at": future5m,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "claude-3-5-sonnet",
			minExpected:    4 * time.Minute,
			maxExpected:    6 * time.Minute,
	placeholder,
		{
			name: "expired rate limit",
			account: &Account{
				Extra: map[string]any{
					modelRateLimitsKey: map[string]any{
						"claude-sonnet-4-5": map[string]any{
							"rate_limit_reset_at": past,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "claude-sonnet-4-5",
			minExpected:    0,
			maxExpected:    0,
	placeholder,
		{
			name:           "no rate limit data",
			account:        &Account{placeholder,
			requestedModel: "claude-sonnet-4-5",
			minExpected:    0,
			maxExpected:    0,
	placeholder,
		{
			name: "no scope fallback",
			account: &Account{
				Extra: map[string]any{
					modelRateLimitsKey: map[string]any{
						"claude_sonnet": map[string]any{
							"rate_limit_reset_at": future5m,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "claude-3-5-sonnet-20241022",
			minExpected:    0,
			maxExpected:    0,
	placeholder,
		{
			name: "antigravity platform - claude-opus-4-5-thinking mapped to opus-4-6-thinking",
			account: &Account{
				Platform: PlatformAntigravity,
				Extra: map[string]any{
					modelRateLimitsKey: map[string]any{
						"claude-opus-4-6-thinking": map[string]any{
							"rate_limit_reset_at": future5m,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "claude-opus-4-5-thinking",
			minExpected:    4 * time.Minute,
			maxExpected:    6 * time.Minute,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.account.GetModelRateLimitRemainingTimeWithContext(context.Background(), tt.requestedModel)
			if result < tt.minExpected || result > tt.maxExpected {
				t.Errorf("GetModelRateLimitRemainingTime() = %v, want between %v and %v", result, tt.minExpected, tt.maxExpected)
		placeholder
	placeholder)
placeholder
placeholder

func TestGetQuotaScopeRateLimitRemainingTime(t *testing.T) {
	now := time.Now()
	future10m := now.Add(10 * time.Minute).Format(time.RFC3339)
	past := now.Add(-10 * time.Minute).Format(time.RFC3339)

	tests := []struct {
		name           string
		account        *Account
		requestedModel string
		minExpected    time.Duration
		maxExpected    time.Duration
placeholder{
		{
			name:           "nil account",
			account:        nil,
			requestedModel: "claude-sonnet-4-5",
			minExpected:    0,
			maxExpected:    0,
	placeholder,
		{
			name: "non-antigravity platform",
			account: &Account{
				Platform: PlatformAnthropic,
				Extra: map[string]any{
					antigravityQuotaScopesKey: map[string]any{
						"claude": map[string]any{
							"rate_limit_reset_at": future10m,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "claude-sonnet-4-5",
			minExpected:    0,
			maxExpected:    0,
	placeholder,
		{
			name: "claude scope rate limited",
			account: &Account{
				Platform: PlatformAntigravity,
				Extra: map[string]any{
					antigravityQuotaScopesKey: map[string]any{
						"claude": map[string]any{
							"rate_limit_reset_at": future10m,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "claude-sonnet-4-5",
			minExpected:    9 * time.Minute,
			maxExpected:    11 * time.Minute,
	placeholder,
		{
			name: "gemini_text scope rate limited",
			account: &Account{
				Platform: PlatformAntigravity,
				Extra: map[string]any{
					antigravityQuotaScopesKey: map[string]any{
						"gemini_text": map[string]any{
							"rate_limit_reset_at": future10m,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "gemini-3-flash",
			minExpected:    9 * time.Minute,
			maxExpected:    11 * time.Minute,
	placeholder,
		{
			name: "expired scope rate limit",
			account: &Account{
				Platform: PlatformAntigravity,
				Extra: map[string]any{
					antigravityQuotaScopesKey: map[string]any{
						"claude": map[string]any{
							"rate_limit_reset_at": past,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "claude-sonnet-4-5",
			minExpected:    0,
			maxExpected:    0,
	placeholder,
		{
			name: "unsupported model",
			account: &Account{
				Platform: PlatformAntigravity,
		placeholder,
			requestedModel: "gpt-4",
			minExpected:    0,
			maxExpected:    0,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.account.GetQuotaScopeRateLimitRemainingTime(tt.requestedModel)
			if result < tt.minExpected || result > tt.maxExpected {
				t.Errorf("GetQuotaScopeRateLimitRemainingTime() = %v, want between %v and %v", result, tt.minExpected, tt.maxExpected)
		placeholder
	placeholder)
placeholder
placeholder

func TestGetRateLimitRemainingTime(t *testing.T) {
	now := time.Now()
	future15m := now.Add(15 * time.Minute).Format(time.RFC3339)
	future5m := now.Add(5 * time.Minute).Format(time.RFC3339)

	tests := []struct {
		name           string
		account        *Account
		requestedModel string
		minExpected    time.Duration
		maxExpected    time.Duration
placeholder{
		{
			name:           "nil account",
			account:        nil,
			requestedModel: "claude-sonnet-4-5",
			minExpected:    0,
			maxExpected:    0,
	placeholder,
		{
			name: "model remaining > scope remaining - returns model",
			account: &Account{
				Platform: PlatformAntigravity,
				Extra: map[string]any{
					modelRateLimitsKey: map[string]any{
						"claude-sonnet-4-5": map[string]any{
							"rate_limit_reset_at": future15m, // 15 分钟
					placeholder,
				placeholder,
					antigravityQuotaScopesKey: map[string]any{
						"claude": map[string]any{
							"rate_limit_reset_at": future5m, // 5 分钟
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "claude-sonnet-4-5",
			minExpected:    14 * time.Minute, // 应返回较大的 15 分钟
			maxExpected:    16 * time.Minute,
	placeholder,
		{
			name: "scope remaining > model remaining - returns scope",
			account: &Account{
				Platform: PlatformAntigravity,
				Extra: map[string]any{
					modelRateLimitsKey: map[string]any{
						"claude-sonnet-4-5": map[string]any{
							"rate_limit_reset_at": future5m, // 5 分钟
					placeholder,
				placeholder,
					antigravityQuotaScopesKey: map[string]any{
						"claude": map[string]any{
							"rate_limit_reset_at": future15m, // 15 分钟
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "claude-sonnet-4-5",
			minExpected:    14 * time.Minute, // 应返回较大的 15 分钟
			maxExpected:    16 * time.Minute,
	placeholder,
		{
			name: "only model rate limited",
			account: &Account{
				Platform: PlatformAntigravity,
				Extra: map[string]any{
					modelRateLimitsKey: map[string]any{
						"claude-sonnet-4-5": map[string]any{
							"rate_limit_reset_at": future5m,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "claude-sonnet-4-5",
			minExpected:    4 * time.Minute,
			maxExpected:    6 * time.Minute,
	placeholder,
		{
			name: "only scope rate limited",
			account: &Account{
				Platform: PlatformAntigravity,
				Extra: map[string]any{
					antigravityQuotaScopesKey: map[string]any{
						"claude": map[string]any{
							"rate_limit_reset_at": future5m,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			requestedModel: "claude-sonnet-4-5",
			minExpected:    4 * time.Minute,
			maxExpected:    6 * time.Minute,
	placeholder,
		{
			name: "neither rate limited",
			account: &Account{
				Platform: PlatformAntigravity,
		placeholder,
			requestedModel: "claude-sonnet-4-5",
			minExpected:    0,
			maxExpected:    0,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.account.GetRateLimitRemainingTimeWithContext(context.Background(), tt.requestedModel)
			if result < tt.minExpected || result > tt.maxExpected {
				t.Errorf("GetRateLimitRemainingTime() = %v, want between %v and %v", result, tt.minExpected, tt.maxExpected)
		placeholder
	placeholder)
placeholder
placeholder
