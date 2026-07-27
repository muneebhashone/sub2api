//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// matchAccountStatsRule
// ---------------------------------------------------------------------------

func TestMatchAccountStatsRule_BothEmpty_NoMatch(t *testing.T) {
	rule := &AccountStatsPricingRule{placeholder
	require.False(t, matchAccountStatsRule(rule, 1, 10))
placeholder

func TestMatchAccountStatsRule_AccountIDMatch(t *testing.T) {
	rule := &AccountStatsPricingRule{AccountIDs: []int64{1, 2, 3placeholderplaceholder
	require.True(t, matchAccountStatsRule(rule, 2, 999))
placeholder

func TestMatchAccountStatsRule_GroupIDMatch(t *testing.T) {
	rule := &AccountStatsPricingRule{GroupIDs: []int64{10, 20placeholderplaceholder
	require.True(t, matchAccountStatsRule(rule, 999, 20))
placeholder

func TestMatchAccountStatsRule_BothConfigured_AccountMatch(t *testing.T) {
	rule := &AccountStatsPricingRule{
		AccountIDs: []int64{1, 2placeholder,
		GroupIDs:   []int64{10, 20placeholder,
placeholder
	require.True(t, matchAccountStatsRule(rule, 2, 999))
placeholder

func TestMatchAccountStatsRule_BothConfigured_GroupMatch(t *testing.T) {
	rule := &AccountStatsPricingRule{
		AccountIDs: []int64{1, 2placeholder,
		GroupIDs:   []int64{10, 20placeholder,
placeholder
	require.True(t, matchAccountStatsRule(rule, 999, 10))
placeholder

func TestMatchAccountStatsRule_BothConfigured_NeitherMatch(t *testing.T) {
	rule := &AccountStatsPricingRule{
		AccountIDs: []int64{1, 2placeholder,
		GroupIDs:   []int64{10, 20placeholder,
placeholder
	require.False(t, matchAccountStatsRule(rule, 999, 999))
placeholder

// ---------------------------------------------------------------------------
// findPricingForModel
// ---------------------------------------------------------------------------

func TestFindPricingForModel(t *testing.T) {
	exactPricing := ChannelModelPricing{
		ID:     1,
		Models: []string{"claude-opus-4"placeholder,
placeholder
	wildcardPricing := ChannelModelPricing{
		ID:     2,
		Models: []string{"claude-*"placeholder,
placeholder
	platformPricing := ChannelModelPricing{
		ID:       3,
		Platform: "openai",
		Models:   []string{"gpt-4o"placeholder,
placeholder
	emptyPlatformPricing := ChannelModelPricing{
		ID:     4,
		Models: []string{"gemini-2.5-pro"placeholder,
placeholder

	tests := []struct {
		name     string
		list     []ChannelModelPricing
		platform string
		model    string
		wantID   int64
		wantNil  bool
placeholder{
		{
			name:     "exact match",
			list:     []ChannelModelPricing{exactPricingplaceholder,
			platform: "anthropic",
			model:    "claude-opus-4",
			wantID:   1,
	placeholder,
		{
			name:     "exact match case insensitive",
			list:     []ChannelModelPricing{{ID: 5, Models: []string{"Claude-Opus-4"placeholderplaceholderplaceholder,
			platform: "",
			model:    "claude-opus-4",
			wantID:   5,
	placeholder,
		{
			name:     "wildcard match",
			list:     []ChannelModelPricing{wildcardPricingplaceholder,
			platform: "anthropic",
			model:    "claude-opus-4",
			wantID:   2,
	placeholder,
		{
			name:     "exact match takes priority over wildcard",
			list:     []ChannelModelPricing{wildcardPricing, exactPricingplaceholder,
			platform: "anthropic",
			model:    "claude-opus-4",
			wantID:   1,
	placeholder,
		{
			name:     "platform mismatch skipped",
			list:     []ChannelModelPricing{platformPricingplaceholder,
			platform: "anthropic",
			model:    "gpt-4o",
			wantNil:  true,
	placeholder,
		{
			name:     "empty platform in pricing matches any",
			list:     []ChannelModelPricing{emptyPlatformPricingplaceholder,
			platform: "gemini",
			model:    "gemini-2.5-pro",
			wantID:   4,
	placeholder,
		{
			name:     "empty platform in query matches any pricing platform",
			list:     []ChannelModelPricing{platformPricingplaceholder,
			platform: "",
			model:    "gpt-4o",
			wantID:   3,
	placeholder,
		{
			name:     "no match at all",
			list:     []ChannelModelPricing{exactPricing, wildcardPricingplaceholder,
			platform: "anthropic",
			model:    "gpt-4o",
			wantNil:  true,
	placeholder,
		{
			name:    "empty list returns nil",
			list:    nil,
			model:   "claude-opus-4",
			wantNil: true,
	placeholder,
		{
			name: "wildcard matches by config order (first match wins)",
			list: []ChannelModelPricing{
				{ID: 10, Models: []string{"claude-*"placeholderplaceholder,
				{ID: 11, Models: []string{"claude-opus-*"placeholderplaceholder,
		placeholder,
			platform: "",
			model:    "claude-opus-4",
			wantID:   10, // config order: "claude-*" is first and matches, so it wins
	placeholder,
		{
			name: "shorter wildcard used when longer does not match",
			list: []ChannelModelPricing{
				{ID: 10, Models: []string{"claude-*"placeholderplaceholder,
				{ID: 11, Models: []string{"claude-opus-*"placeholderplaceholder,
		placeholder,
			platform: "",
			model:    "claude-sonnet-4",
			wantID:   10, // only "claude-*" matches
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findPricingForModel(tt.list, tt.platform, tt.model)
			if tt.wantNil {
				require.Nil(t, result)
				return
		placeholder
			require.NotNil(t, result)
			require.Equal(t, tt.wantID, result.ID)
	placeholder)
placeholder
placeholder

// ---------------------------------------------------------------------------
// calculateStatsCost
// ---------------------------------------------------------------------------

func TestCalculateStatsCost_NilPricing(t *testing.T) {
	result := calculateStatsCost(nil, UsageTokens{placeholder, 1)
	require.Nil(t, result)
placeholder

func TestCalculateStatsCost_TokenBilling(t *testing.T) {
	pricing := &ChannelModelPricing{
		BillingMode: BillingModeToken,
		InputPrice:  testPtrFloat64(0.001),
		OutputPrice: testPtrFloat64(0.002),
placeholder
	tokens := UsageTokens{
		InputTokens:  100,
		OutputTokens: 50,
placeholder
	result := calculateStatsCost(pricing, tokens, 1)
	require.NotNil(t, result)
	// 100*0.001 + 50*0.002 = 0.1 + 0.1 = 0.2
	require.InDelta(t, 0.2, *result, 1e-12)
placeholder

func TestCalculateStatsCost_TokenBilling_WithCache(t *testing.T) {
	pricing := &ChannelModelPricing{
		BillingMode:     BillingModeToken,
		InputPrice:      testPtrFloat64(0.001),
		OutputPrice:     testPtrFloat64(0.002),
		CacheWritePrice: testPtrFloat64(0.003),
		CacheReadPrice:  testPtrFloat64(0.0005),
placeholder
	tokens := UsageTokens{
		InputTokens:         100,
		OutputTokens:        50,
		CacheCreationTokens: 200,
		CacheReadTokens:     300,
placeholder
	result := calculateStatsCost(pricing, tokens, 1)
	require.NotNil(t, result)
	// 100*0.001 + 50*0.002 + 200*0.003 + 300*0.0005
	// = 0.1 + 0.1 + 0.6 + 0.15 = 0.95
	require.InDelta(t, 0.95, *result, 1e-12)
placeholder

func TestCalculateStatsCost_TokenBilling_WithImageOutput(t *testing.T) {
	pricing := &ChannelModelPricing{
		BillingMode:      BillingModeToken,
		InputPrice:       testPtrFloat64(0.001),
		OutputPrice:      testPtrFloat64(0.002),
		ImageOutputPrice: testPtrFloat64(0.01),
placeholder
	tokens := UsageTokens{
		InputTokens:       100,
		OutputTokens:      50,
		ImageOutputTokens: 10,
placeholder
	result := calculateStatsCost(pricing, tokens, 1)
	require.NotNil(t, result)
	// 100*0.001 + 50*0.002 + 10*0.01 = 0.1 + 0.1 + 0.1 = 0.3
	require.InDelta(t, 0.3, *result, 1e-12)
placeholder

func TestCalculateStatsCost_TokenBilling_PartialPricesNil(t *testing.T) {
	pricing := &ChannelModelPricing{
		BillingMode: BillingModeToken,
		InputPrice:  testPtrFloat64(0.001),
		// OutputPrice, CacheWritePrice, etc. are all nil → treated as 0
placeholder
	tokens := UsageTokens{
		InputTokens:         100,
		OutputTokens:        50,
		CacheCreationTokens: 200,
placeholder
	result := calculateStatsCost(pricing, tokens, 1)
	require.NotNil(t, result)
	// Only input contributes: 100*0.001 = 0.1
	require.InDelta(t, 0.1, *result, 1e-12)
placeholder

func TestCalculateStatsCost_TokenBilling_AllTokensZero(t *testing.T) {
	pricing := &ChannelModelPricing{
		BillingMode: BillingModeToken,
		InputPrice:  testPtrFloat64(0.001),
		OutputPrice: testPtrFloat64(0.002),
placeholder
	tokens := UsageTokens{placeholder // all zeros
	result := calculateStatsCost(pricing, tokens, 1)
	// totalCost == 0 → returns nil (does not override, falls back to default formula)
	require.Nil(t, result)
placeholder

func TestCalculateStatsCost_PerRequestBilling(t *testing.T) {
	pricing := &ChannelModelPricing{
		BillingMode:     BillingModePerRequest,
		PerRequestPrice: testPtrFloat64(0.05),
placeholder
	tokens := UsageTokens{InputTokens: 999, OutputTokens: 999placeholder
	result := calculateStatsCost(pricing, tokens, 3)
	require.NotNil(t, result)
	// 0.05 * 3 = 0.15
	require.InDelta(t, 0.15, *result, 1e-12)
placeholder

func TestCalculateStatsCost_PerRequestBilling_PriceNil(t *testing.T) {
	pricing := &ChannelModelPricing{
		BillingMode: BillingModePerRequest,
		// PerRequestPrice is nil
placeholder
	result := calculateStatsCost(pricing, UsageTokens{placeholder, 1)
	require.Nil(t, result)
placeholder

func TestCalculateStatsCost_PerRequestBilling_PriceZero(t *testing.T) {
	pricing := &ChannelModelPricing{
		BillingMode:     BillingModePerRequest,
		PerRequestPrice: testPtrFloat64(0),
placeholder
	result := calculateStatsCost(pricing, UsageTokens{placeholder, 1)
	// price == 0 → condition *pricing.PerRequestPrice > 0 is false → returns nil
	require.Nil(t, result)
placeholder

func TestCalculateStatsCost_ImageBilling(t *testing.T) {
	pricing := &ChannelModelPricing{
		BillingMode:     BillingModeImage,
		PerRequestPrice: testPtrFloat64(0.10),
placeholder
	result := calculateStatsCost(pricing, UsageTokens{placeholder, 2)
	require.NotNil(t, result)
	// 0.10 * 2 = 0.20
	require.InDelta(t, 0.20, *result, 1e-12)
placeholder

func TestCalculateStatsCost_ImageBilling_PriceNil(t *testing.T) {
	pricing := &ChannelModelPricing{
		BillingMode: BillingModeImage,
		// PerRequestPrice is nil
placeholder
	result := calculateStatsCost(pricing, UsageTokens{placeholder, 1)
	require.Nil(t, result)
placeholder

func TestCalculateStatsCost_DefaultBillingMode_FallsToToken(t *testing.T) {
	// BillingMode is empty string (default) → falls into token billing
	pricing := &ChannelModelPricing{
		InputPrice:  testPtrFloat64(0.001),
		OutputPrice: testPtrFloat64(0.002),
placeholder
	tokens := UsageTokens{
		InputTokens:  100,
		OutputTokens: 50,
placeholder
	result := calculateStatsCost(pricing, tokens, 1)
	require.NotNil(t, result)
	require.InDelta(t, 0.2, *result, 1e-12)
placeholder

// ---------------------------------------------------------------------------
// tryCustomRules — 多规则顺序测试
// ---------------------------------------------------------------------------

func TestTryCustomRules_FirstMatchWins(t *testing.T) {
	channel := &Channel{
		AccountStatsPricingRules: []AccountStatsPricingRule{
			{
				GroupIDs: []int64{1placeholder,
				Pricing: []ChannelModelPricing{
					{ID: 100, Models: []string{"claude-opus-4"placeholder, InputPrice: testPtrFloat64(0.01), OutputPrice: testPtrFloat64(0.02)placeholder,
			placeholder,
		placeholder,
			{
				GroupIDs: []int64{1placeholder,
				Pricing: []ChannelModelPricing{
					{ID: 200, Models: []string{"claude-opus-4"placeholder, InputPrice: testPtrFloat64(0.99), OutputPrice: testPtrFloat64(0.99)placeholder,
			placeholder,
		placeholder,
	placeholder,
placeholder
	tokens := UsageTokens{InputTokens: 100, OutputTokens: 50placeholder
	result := tryCustomRules(channel, 999, 1, "", "claude-opus-4", tokens, 1)
	require.NotNil(t, result)
	// 应使用第一条规则的价格：100*0.01 + 50*0.02 = 2.0
	require.InDelta(t, 2.0, *result, 1e-12)
placeholder

func TestTryCustomRules_SkipsNonMatchingRules(t *testing.T) {
	channel := &Channel{
		AccountStatsPricingRules: []AccountStatsPricingRule{
			{
				AccountIDs: []int64{888placeholder, // 不匹配
				Pricing: []ChannelModelPricing{
					{ID: 100, Models: []string{"claude-opus-4"placeholder, InputPrice: testPtrFloat64(0.99)placeholder,
			placeholder,
		placeholder,
			{
				GroupIDs: []int64{1placeholder, // 匹配
				Pricing: []ChannelModelPricing{
					{ID: 200, Models: []string{"claude-opus-4"placeholder, InputPrice: testPtrFloat64(0.05)placeholder,
			placeholder,
		placeholder,
	placeholder,
placeholder
	tokens := UsageTokens{InputTokens: 100placeholder
	result := tryCustomRules(channel, 999, 1, "", "claude-opus-4", tokens, 1)
	require.NotNil(t, result)
	// 跳过规则1（账号不匹配），使用规则2：100*0.05 = 5.0
	require.InDelta(t, 5.0, *result, 1e-12)
placeholder

func TestTryCustomRules_NoMatch_ReturnsNil(t *testing.T) {
	channel := &Channel{
		AccountStatsPricingRules: []AccountStatsPricingRule{
			{
				AccountIDs: []int64{888placeholder,
				Pricing: []ChannelModelPricing{
					{ID: 100, Models: []string{"claude-opus-4"placeholder, InputPrice: testPtrFloat64(0.01)placeholder,
			placeholder,
		placeholder,
	placeholder,
placeholder
	tokens := UsageTokens{InputTokens: 100placeholder
	result := tryCustomRules(channel, 999, 2, "", "claude-opus-4", tokens, 1)
	require.Nil(t, result) // 账号和分组都不匹配
placeholder

func TestTryCustomRules_RuleMatchesButModelNot_ContinuesToNext(t *testing.T) {
	channel := &Channel{
		AccountStatsPricingRules: []AccountStatsPricingRule{
			{
				GroupIDs: []int64{1placeholder,
				Pricing: []ChannelModelPricing{
					{ID: 100, Models: []string{"gpt-4o"placeholder, InputPrice: testPtrFloat64(0.01)placeholder, // 模型不匹配
			placeholder,
		placeholder,
			{
				GroupIDs: []int64{1placeholder,
				Pricing: []ChannelModelPricing{
					{ID: 200, Models: []string{"claude-opus-4"placeholder, InputPrice: testPtrFloat64(0.05)placeholder, // 模型匹配
			placeholder,
		placeholder,
	placeholder,
placeholder
	tokens := UsageTokens{InputTokens: 100placeholder
	result := tryCustomRules(channel, 999, 1, "", "claude-opus-4", tokens, 1)
	require.NotNil(t, result)
	require.InDelta(t, 5.0, *result, 1e-12) // 使用规则2
placeholder

// ---------------------------------------------------------------------------
// tryModelFilePricing
// ---------------------------------------------------------------------------

// newTestBillingServiceWithPrices creates a BillingService with pre-populated
// fallback prices for testing. No config or pricing service is needed.
// The key must match what getFallbackPricing resolves to for a given model name.
// E.g., model "claude-sonnet-4" resolves to key "claude-sonnet-4".
func newTestBillingServiceWithPrices(prices map[string]*ModelPricing) *BillingService {
	return &BillingService{
		fallbackPrices: prices,
placeholder
placeholder

func TestTryModelFilePricing_Success(t *testing.T) {
	bs := newTestBillingServiceWithPrices(map[string]*ModelPricing{
		"claude-sonnet-4": {
			InputPricePerToken:  0.001,
			OutputPricePerToken: 0.002,
	placeholder,
placeholder)
	tokens := UsageTokens{InputTokens: 100, OutputTokens: 50placeholder
	result := tryModelFilePricing(bs, "claude-sonnet-4", tokens)
	require.NotNil(t, result)
	// 100*0.001 + 50*0.002 = 0.1 + 0.1 = 0.2
	require.InDelta(t, 0.2, *result, 1e-12)
placeholder

func TestTryModelFilePricing_AppliesLongContextPricing(t *testing.T) {
	bs := newTestBillingServiceWithPrices(map[string]*ModelPricing{
		"gpt-5.6-sol": {
			InputPricePerToken:          0.001,
			OutputPricePerToken:         0.002,
			CacheReadPricePerToken:      0.0001,
			LongContextInputThreshold:   100,
			LongContextInputMultiplier:  2,
			LongContextOutputMultiplier: 1.5,
	placeholder,
placeholder)
	tokens := UsageTokens{InputTokens: 101, OutputTokens: 10, CacheReadTokens: 5placeholder

	result := tryModelFilePricing(bs, "gpt-5.6-sol", tokens)

	require.NotNil(t, result)
	// Input and cache-read use the 2x input tier; output uses the 1.5x tier.
	require.InDelta(t, 0.233, *result, 1e-12)
placeholder

func TestTryModelFilePricing_PricingNotFound(t *testing.T) {
	// "nonexistent-model" does not match any fallback pattern
	bs := newTestBillingServiceWithPrices(map[string]*ModelPricing{placeholder)
	tokens := UsageTokens{InputTokens: 100, OutputTokens: 50placeholder
	result := tryModelFilePricing(bs, "nonexistent-model", tokens)
	require.Nil(t, result)
placeholder

func TestTryModelFilePricing_NilFallback(t *testing.T) {
	// getFallbackPricing returns nil when key maps to nil
	bs := newTestBillingServiceWithPrices(map[string]*ModelPricing{
		"claude-sonnet-4": nil,
placeholder)
	tokens := UsageTokens{InputTokens: 100placeholder
	result := tryModelFilePricing(bs, "claude-sonnet-4", tokens)
	require.Nil(t, result)
placeholder

func TestTryModelFilePricing_ZeroCost(t *testing.T) {
	bs := newTestBillingServiceWithPrices(map[string]*ModelPricing{
		"claude-sonnet-4": {
			InputPricePerToken:  0.001,
			OutputPricePerToken: 0.002,
	placeholder,
placeholder)
	tokens := UsageTokens{placeholder // all zero tokens → cost = 0 → nil
	result := tryModelFilePricing(bs, "claude-sonnet-4", tokens)
	require.Nil(t, result)
placeholder

func TestTryModelFilePricing_WithImageOutput(t *testing.T) {
	bs := newTestBillingServiceWithPrices(map[string]*ModelPricing{
		"claude-sonnet-4": {
			InputPricePerToken:       0.001,
			OutputPricePerToken:      0.002,
			ImageOutputPricePerToken: 0.01,
	placeholder,
placeholder)
	tokens := UsageTokens{
		InputTokens:       100,
		OutputTokens:      50,
		ImageOutputTokens: 10,
placeholder
	result := tryModelFilePricing(bs, "claude-sonnet-4", tokens)
	require.NotNil(t, result)
	// 100*0.001 + 50*0.002 + 10*0.01 = 0.1 + 0.1 + 0.1 = 0.3
	require.InDelta(t, 0.3, *result, 1e-12)
placeholder

func TestTryModelFilePricing_WithCacheTokens(t *testing.T) {
	bs := newTestBillingServiceWithPrices(map[string]*ModelPricing{
		"claude-sonnet-4": {
			InputPricePerToken:         0.001,
			OutputPricePerToken:        0.002,
			CacheCreationPricePerToken: 0.003,
			CacheReadPricePerToken:     0.0005,
	placeholder,
placeholder)
	tokens := UsageTokens{
		InputTokens:         100,
		OutputTokens:        50,
		CacheCreationTokens: 200,
		CacheReadTokens:     300,
placeholder
	result := tryModelFilePricing(bs, "claude-sonnet-4", tokens)
	require.NotNil(t, result)
	// 100*0.001 + 50*0.002 + 200*0.003 + 300*0.0005
	// = 0.1 + 0.1 + 0.6 + 0.15 = 0.95
	require.InDelta(t, 0.95, *result, 1e-12)
placeholder

// ---------------------------------------------------------------------------
// resolveAccountStatsCost — integration tests covering the 4-level priority chain
// ---------------------------------------------------------------------------

func TestResolveAccountStatsCost_NilChannelService(t *testing.T) {
	result := resolveAccountStatsCost(
		context.Background(),
		nil, // channelService is nil
		newTestBillingServiceWithPrices(map[string]*ModelPricing{placeholder),
		1, 1, "claude-sonnet-4",
		UsageTokens{InputTokens: 100placeholder, 1, 0.5,
	)
	require.Nil(t, result)
placeholder

func TestResolveAccountStatsCost_EmptyUpstreamModel(t *testing.T) {
	cs := newTestChannelServiceForStats(t, &Channel{
		ID:     1,
		Status: StatusActive,
placeholder, 1, "")

	result := resolveAccountStatsCost(
		context.Background(),
		cs,
		newTestBillingServiceWithPrices(map[string]*ModelPricing{placeholder),
		1, 1, "", // empty upstream model
		UsageTokens{InputTokens: 100placeholder, 1, 0.5,
	)
	require.Nil(t, result)
placeholder

func TestResolveAccountStatsCost_GetChannelForGroupReturnsNil(t *testing.T) {
	// Group 99 is NOT in the cache, so GetChannelForGroup returns nil
	cs := newTestChannelServiceForStats(t, &Channel{
		ID:     1,
		Status: StatusActive,
placeholder, 1, "")

	result := resolveAccountStatsCost(
		context.Background(),
		cs,
		newTestBillingServiceWithPrices(map[string]*ModelPricing{placeholder),
		1, 99, "claude-sonnet-4", // groupID 99 has no channel
		UsageTokens{InputTokens: 100placeholder, 1, 0.5,
	)
	require.Nil(t, result)
placeholder

func TestResolveAccountStatsCost_HitsCustomRule(t *testing.T) {
	channel := &Channel{
		ID:     1,
		Status: StatusActive,
		AccountStatsPricingRules: []AccountStatsPricingRule{
			{
				GroupIDs: []int64{10placeholder,
				Pricing: []ChannelModelPricing{
					{
						ID:          100,
						Models:      []string{"claude-sonnet-4"placeholder,
						InputPrice:  testPtrFloat64(0.01),
						OutputPrice: testPtrFloat64(0.02),
				placeholder,
			placeholder,
		placeholder,
	placeholder,
placeholder
	cs := newTestChannelServiceForStats(t, channel, 10, "anthropic")

	tokens := UsageTokens{InputTokens: 100, OutputTokens: 50placeholder

	result := resolveAccountStatsCost(
		context.Background(),
		cs, nil, // billingService not needed when custom rule hits
		1, 10, "claude-sonnet-4",
		tokens, 1, 999.0, // totalCost ignored because custom rule hits
	)
	require.NotNil(t, result)
	// 100*0.01 + 50*0.02 = 1.0 + 1.0 = 2.0
	require.InDelta(t, 2.0, *result, 1e-12)
placeholder

func TestResolveAccountStatsCost_ApplyPricingToAccountStats_UsesTotalCost(t *testing.T) {
	channel := &Channel{
		ID:                         1,
		Status:                     StatusActive,
		ApplyPricingToAccountStats: true,
		// No custom rules
placeholder
	cs := newTestChannelServiceForStats(t, channel, 10, "anthropic")

	tokens := UsageTokens{InputTokens: 100, OutputTokens: 50placeholder

	result := resolveAccountStatsCost(
		context.Background(),
		cs, nil,
		1, 10, "claude-sonnet-4",
		tokens, 1, 0.75, // totalCost = 0.75
	)
	require.NotNil(t, result)
	require.InDelta(t, 0.75, *result, 1e-12)
placeholder

func TestResolveAccountStatsCost_ApplyPricingToAccountStats_ZeroTotalCost_ReturnsNil(t *testing.T) {
	channel := &Channel{
		ID:                         1,
		Status:                     StatusActive,
		ApplyPricingToAccountStats: true,
placeholder
	cs := newTestChannelServiceForStats(t, channel, 10, "anthropic")

	result := resolveAccountStatsCost(
		context.Background(),
		cs, nil,
		1, 10, "claude-sonnet-4",
		UsageTokens{placeholder, 1, 0.0, // totalCost = 0
	)
	require.Nil(t, result)
placeholder

func TestResolveAccountStatsCost_FallsBackToLiteLLM(t *testing.T) {
	channel := &Channel{
		ID:                         1,
		Status:                     StatusActive,
		ApplyPricingToAccountStats: false, // not enabled
		// No custom rules
placeholder
	cs := newTestChannelServiceForStats(t, channel, 10, "anthropic")

	bs := newTestBillingServiceWithPrices(map[string]*ModelPricing{
		"claude-sonnet-4": {
			InputPricePerToken:  0.001,
			OutputPricePerToken: 0.002,
	placeholder,
placeholder)

	tokens := UsageTokens{InputTokens: 100, OutputTokens: 50placeholder

	result := resolveAccountStatsCost(
		context.Background(),
		cs, bs,
		1, 10, "claude-sonnet-4",
		tokens, 1, 999.0, // totalCost ignored
	)
	require.NotNil(t, result)
	// 100*0.001 + 50*0.002 = 0.1 + 0.1 = 0.2
	require.InDelta(t, 0.2, *result, 1e-12)
placeholder

func TestResolveAccountStatsCost_Gemini36FlashTierUsesFallbackPricing(t *testing.T) {
	channel := &Channel{
		ID:                         1,
		Status:                     StatusActive,
		ApplyPricingToAccountStats: false,
placeholder
	cs := newTestChannelServiceForStats(t, channel, 10, "antigravity")
	bs := NewBillingService(&config.Config{placeholder, nil)

	result := resolveAccountStatsCost(
		context.Background(),
		cs, bs,
		1, 10, "gemini-3.6-flash-low",
		UsageTokens{InputTokens: 1_000_000, OutputTokens: 1_000_000, CacheReadTokens: 1_000_000placeholder, 1, 0,
	)
	require.NotNil(t, result)
	require.InDelta(t, 9.15, *result, 1e-12)
placeholder

func TestResolveAccountStatsCost_AllMiss_ReturnsNil(t *testing.T) {
	channel := &Channel{
		ID:                         1,
		Status:                     StatusActive,
		ApplyPricingToAccountStats: false,
		// No custom rules
placeholder
	cs := newTestChannelServiceForStats(t, channel, 10, "anthropic")

	// BillingService with no pricing for the model
	bs := newTestBillingServiceWithPrices(map[string]*ModelPricing{placeholder)

	tokens := UsageTokens{InputTokens: 100, OutputTokens: 50placeholder

	result := resolveAccountStatsCost(
		context.Background(),
		cs, bs,
		1, 10, "totally-unknown-model",
		tokens, 1, 0.0,
	)
	require.Nil(t, result)
placeholder

func TestResolveAccountStatsCost_NilBillingService_SkipsLiteLLM(t *testing.T) {
	channel := &Channel{
		ID:                         1,
		Status:                     StatusActive,
		ApplyPricingToAccountStats: false,
placeholder
	cs := newTestChannelServiceForStats(t, channel, 10, "anthropic")

	result := resolveAccountStatsCost(
		context.Background(),
		cs, nil, // billingService is nil
		1, 10, "claude-sonnet-4",
		UsageTokens{InputTokens: 100placeholder, 1, 0.0,
	)
	require.Nil(t, result)
placeholder

func TestResolveAccountStatsCost_CustomRulePriorityOverApplyPricing(t *testing.T) {
	// Both custom rule and ApplyPricingToAccountStats are configured;
	// custom rule should take precedence.
	channel := &Channel{
		ID:                         1,
		Status:                     StatusActive,
		ApplyPricingToAccountStats: true,
		AccountStatsPricingRules: []AccountStatsPricingRule{
			{
				GroupIDs: []int64{10placeholder,
				Pricing: []ChannelModelPricing{
					{
						ID:         100,
						Models:     []string{"claude-sonnet-4"placeholder,
						InputPrice: testPtrFloat64(0.05),
				placeholder,
			placeholder,
		placeholder,
	placeholder,
placeholder
	cs := newTestChannelServiceForStats(t, channel, 10, "anthropic")

	tokens := UsageTokens{InputTokens: 100placeholder

	result := resolveAccountStatsCost(
		context.Background(),
		cs, nil,
		1, 10, "claude-sonnet-4",
		tokens, 1, 99.0, // totalCost = 99.0 (would be used if ApplyPricing wins)
	)
	require.NotNil(t, result)
	// Custom rule: 100*0.05 = 5.0 (NOT 99.0 from totalCost)
	require.InDelta(t, 5.0, *result, 1e-12)
placeholder

// ---------------------------------------------------------------------------
// helpers for resolveAccountStatsCost tests
// ---------------------------------------------------------------------------

// newTestChannelServiceForStats creates a ChannelService with a single channel
// mapped to the given groupID, suitable for resolveAccountStatsCost tests.
func newTestChannelServiceForStats(t *testing.T, channel *Channel, groupID int64, platform string) *ChannelService {
placeholder
	cache := newEmptyChannelCache()
	cache.channelByGroupID[groupID] = channel
	cache.groupPlatform[groupID] = platform
	cs := &ChannelService{placeholder
	cache.loadedAt = time.Now()
	cs.cache.Store(cache)
	return cs
placeholder
