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
// CalculateCostUnified
// ---------------------------------------------------------------------------

func TestCalculateCostUnified_NilResolver_FallsBackToOldPath(t *testing.T) {
	svc := newTestBillingService()

	tokens := UsageTokens{InputTokens: 1000, OutputTokens: 500placeholder
	input := CostInput{
		Model:          "claude-sonnet-4",
		Tokens:         tokens,
		RateMultiplier: 1.0,
		Resolver:       nil, // no resolver
placeholder
	cost, err := svc.CalculateCostUnified(input)
placeholder

	// Should match the old-path result exactly
	expected, err := svc.calculateCostInternal("claude-sonnet-4", tokens, 1.0, "", nil)
placeholder
	require.InDelta(t, expected.TotalCost, cost.TotalCost, 1e-10)
	require.InDelta(t, expected.ActualCost, cost.ActualCost, 1e-10)
	// BillingMode is NOT set by old path through CalculateCostUnified (resolver == nil)
	require.Empty(t, cost.BillingMode)
placeholder

func TestCalculateCostUnified_TokenMode(t *testing.T) {
	bs := newTestBillingService()
	resolver := NewModelPricingResolver(nil, bs)

	tokens := UsageTokens{InputTokens: 1000, OutputTokens: 500placeholder
	input := CostInput{
		Ctx:            context.Background(),
		Model:          "claude-sonnet-4",
		Tokens:         tokens,
		RateMultiplier: 1.5,
		Resolver:       resolver,
placeholder
	cost, err := bs.CalculateCostUnified(input)
placeholder
	require.NotNil(t, cost)

	// Verify token billing: Input: 1000*3e-6=0.003, Output: 500*15e-6=0.0075
	expectedTotal := 1000*3e-6 + 500*15e-6
	require.InDelta(t, expectedTotal, cost.TotalCost, 1e-10)
	require.InDelta(t, expectedTotal*1.5, cost.ActualCost, 1e-10)
	require.Equal(t, string(BillingModeToken), cost.BillingMode)
placeholder

func TestCalculateCostUnified_PerRequestMode(t *testing.T) {
	// Set up a ChannelService with a per-request pricing channel
	cs := newTestChannelServiceWithCache(t, &channelCache{
		pricingByGroupModel: map[channelModelKey]*ChannelModelPricing{
			{groupID: 1, model: "claude-sonnet-4"placeholder: {
				BillingMode:     BillingModePerRequest,
				PerRequestPrice: testPtrFloat64(0.05),
		placeholder,
	placeholder,
		channelByGroupID: map[int64]*Channel{
			1: {ID: 1, Status: StatusActiveplaceholder,
	placeholder,
		groupPlatform:           map[int64]string{1: ""placeholder,
		wildcardByGroupPlatform: map[channelGroupPlatformKey][]*wildcardPricingEntry{placeholder,
		mappingByGroupModel:     map[channelModelKey]string{placeholder,
		wildcardMappingByGP:     map[channelGroupPlatformKey][]*wildcardMappingEntry{placeholder,
		byID:                    map[int64]*Channel{placeholder,
placeholder)

	bs := newTestBillingService()
	resolver := NewModelPricingResolver(cs, bs)
	groupID := int64(1)

	input := CostInput{
		Ctx:            context.Background(),
		Model:          "claude-sonnet-4",
		GroupID:        &groupID,
		Tokens:         UsageTokens{InputTokens: 100, OutputTokens: 50placeholder,
		RequestCount:   3,
		RateMultiplier: 2.0,
		Resolver:       resolver,
placeholder
	cost, err := bs.CalculateCostUnified(input)
placeholder
	require.NotNil(t, cost)

	// 3 requests * $0.05 = $0.15
	require.InDelta(t, 0.15, cost.TotalCost, 1e-10)
	// ActualCost = 0.15 * 2.0 = 0.30
	require.InDelta(t, 0.30, cost.ActualCost, 1e-10)
	require.Equal(t, string(BillingModePerRequest), cost.BillingMode)
placeholder

func TestCalculateCostUnified_ImageMode(t *testing.T) {
	cs := newTestChannelServiceWithCache(t, &channelCache{
		pricingByGroupModel: map[channelModelKey]*ChannelModelPricing{
			{groupID: 2, model: "gemini-image"placeholder: {
				BillingMode:     BillingModeImage,
				PerRequestPrice: testPtrFloat64(0.10),
		placeholder,
	placeholder,
		channelByGroupID: map[int64]*Channel{
			2: {ID: 2, Status: StatusActiveplaceholder,
	placeholder,
		groupPlatform:           map[int64]string{2: ""placeholder,
		wildcardByGroupPlatform: map[channelGroupPlatformKey][]*wildcardPricingEntry{placeholder,
		mappingByGroupModel:     map[channelModelKey]string{placeholder,
		wildcardMappingByGP:     map[channelGroupPlatformKey][]*wildcardMappingEntry{placeholder,
		byID:                    map[int64]*Channel{placeholder,
placeholder)

	bs := &BillingService{
		cfg:            &config.Config{placeholder,
		fallbackPrices: map[string]*ModelPricing{placeholder,
placeholder
	resolver := NewModelPricingResolver(cs, bs)
	groupID := int64(2)

	input := CostInput{
		Ctx:            context.Background(),
		Model:          "gemini-image",
		GroupID:        &groupID,
		Tokens:         UsageTokens{placeholder,
		RequestCount:   2,
		RateMultiplier: 1.0,
		Resolver:       resolver,
placeholder
	cost, err := bs.CalculateCostUnified(input)
placeholder
	require.NotNil(t, cost)

	// 2 * $0.10 = $0.20
	require.InDelta(t, 0.20, cost.TotalCost, 1e-10)
	require.InDelta(t, 0.20, cost.ActualCost, 1e-10)
	require.Equal(t, string(BillingModeImage), cost.BillingMode)
placeholder

func TestCalculateCostUnified_RateMultiplierZeroDefaultsToOne(t *testing.T) {
	bs := newTestBillingService()
	resolver := NewModelPricingResolver(nil, bs)

	tokens := UsageTokens{InputTokens: 1000, OutputTokens: 500placeholder

	costZero, err := bs.CalculateCostUnified(CostInput{
		Ctx:            context.Background(),
		Model:          "claude-sonnet-4",
		Tokens:         tokens,
		RateMultiplier: 0, // should default to 1.0
		Resolver:       resolver,
placeholder)
placeholder

	costOne, err := bs.CalculateCostUnified(CostInput{
		Ctx:            context.Background(),
		Model:          "claude-sonnet-4",
		Tokens:         tokens,
		RateMultiplier: 1.0,
		Resolver:       resolver,
placeholder)
placeholder

	require.InDelta(t, costOne.ActualCost, costZero.ActualCost, 1e-10)
placeholder

func TestCalculateCostUnified_NegativeRateMultiplierDefaultsToOne(t *testing.T) {
	bs := newTestBillingService()
	resolver := NewModelPricingResolver(nil, bs)

	tokens := UsageTokens{InputTokens: 1000placeholder

	costNeg, err := bs.CalculateCostUnified(CostInput{
		Ctx:            context.Background(),
		Model:          "claude-sonnet-4",
		Tokens:         tokens,
		RateMultiplier: -5.0,
		Resolver:       resolver,
placeholder)
placeholder

	costOne, err := bs.CalculateCostUnified(CostInput{
		Ctx:            context.Background(),
		Model:          "claude-sonnet-4",
		Tokens:         tokens,
		RateMultiplier: 1.0,
		Resolver:       resolver,
placeholder)
placeholder

	require.InDelta(t, costOne.ActualCost, costNeg.ActualCost, 1e-10)
placeholder

func TestCalculateCostUnified_BillingModeFieldFilled(t *testing.T) {
	bs := newTestBillingService()
	resolver := NewModelPricingResolver(nil, bs)

	cost, err := bs.CalculateCostUnified(CostInput{
		Ctx:            context.Background(),
		Model:          "claude-sonnet-4",
		Tokens:         UsageTokens{InputTokens: 100placeholder,
		RateMultiplier: 1.0,
		Resolver:       resolver,
placeholder)
placeholder
	require.Equal(t, "token", cost.BillingMode)
placeholder

func TestCalculateCostUnified_UsesPreResolvedPricing(t *testing.T) {
	bs := newTestBillingService()
	resolver := NewModelPricingResolver(nil, bs)

	// Pre-resolve with per_request mode to verify it's used instead of re-resolving
	preResolved := &ResolvedPricing{
		Mode:                   BillingModePerRequest,
		DefaultPerRequestPrice: 0.07,
placeholder

	cost, err := bs.CalculateCostUnified(CostInput{
		Ctx:            context.Background(),
		Model:          "claude-sonnet-4",
		Tokens:         UsageTokens{InputTokens: 100placeholder,
		RequestCount:   2,
		RateMultiplier: 1.0,
		Resolver:       resolver,
		Resolved:       preResolved,
placeholder)
placeholder
	require.NotNil(t, cost)

	// 2 * $0.07 = $0.14
	require.InDelta(t, 0.14, cost.TotalCost, 1e-10)
	require.Equal(t, string(BillingModePerRequest), cost.BillingMode)
placeholder

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// newTestChannelServiceWithCache creates a ChannelService with a pre-populated
// cache snapshot, bypassing the repository layer entirely.
func newTestChannelServiceWithCache(t *testing.T, cache *channelCache) *ChannelService {
placeholder
	cs := &ChannelService{placeholder
	cache.loadedAt = time.Now()
	cs.cache.Store(cache)
	return cs
placeholder
