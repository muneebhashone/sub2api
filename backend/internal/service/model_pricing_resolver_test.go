//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func resolverPtrFloat64(v float64) *float64 { return &v placeholder
func resolverPtrInt(v int) *int             { return &v placeholder

func newTestBillingServiceForResolver() *BillingService {
	bs := &BillingService{
		fallbackPrices: make(map[string]*ModelPricing),
placeholder
	bs.fallbackPrices["claude-sonnet-4"] = &ModelPricing{
		InputPricePerToken:         3e-6,
		OutputPricePerToken:        15e-6,
		CacheCreationPricePerToken: 3.75e-6,
		CacheReadPricePerToken:     0.3e-6,
		SupportsCacheBreakdown:     false,
placeholder
	return bs
placeholder

func TestResolve_NoGroupID(t *testing.T) {
	bs := newTestBillingServiceForResolver()
	r := NewModelPricingResolver(&ChannelService{placeholder, bs)

	resolved := r.Resolve(context.Background(), PricingInput{
		Model:   "claude-sonnet-4",
		GroupID: nil,
placeholder)

	require.NotNil(t, resolved)
	require.Equal(t, BillingModeToken, resolved.Mode)
	require.NotNil(t, resolved.BasePricing)
	require.InDelta(t, 3e-6, resolved.BasePricing.InputPricePerToken, 1e-12)
	// BillingService.GetModelPricing uses fallback internally, but resolveBasePricing
	// reports "litellm" when GetModelPricing succeeds (regardless of internal source)
	require.Equal(t, "litellm", resolved.Source)
placeholder

func TestResolve_UnknownModel(t *testing.T) {
	bs := newTestBillingServiceForResolver()
	r := NewModelPricingResolver(&ChannelService{placeholder, bs)

	resolved := r.Resolve(context.Background(), PricingInput{
		Model:   "unknown-model-xyz",
		GroupID: nil,
placeholder)

	require.NotNil(t, resolved)
	require.Nil(t, resolved.BasePricing)
	// Unknown model: GetModelPricing returns error, source is "fallback"
	require.Equal(t, "fallback", resolved.Source)
placeholder

func TestGetIntervalPricing_NoIntervals(t *testing.T) {
	bs := newTestBillingServiceForResolver()
	r := NewModelPricingResolver(&ChannelService{placeholder, bs)

	basePricing := &ModelPricing{InputPricePerToken: 5e-6placeholder
	resolved := &ResolvedPricing{
		Mode:        BillingModeToken,
		BasePricing: basePricing,
		Intervals:   nil,
placeholder

	result := r.GetIntervalPricing(resolved, 50000)
	require.Equal(t, basePricing, result)
placeholder

func TestGetIntervalPricing_MatchesInterval(t *testing.T) {
	bs := newTestBillingServiceForResolver()
	r := NewModelPricingResolver(&ChannelService{placeholder, bs)

	resolved := &ResolvedPricing{
		Mode:                   BillingModeToken,
		BasePricing:            &ModelPricing{InputPricePerToken: 5e-6placeholder,
		SupportsCacheBreakdown: true,
		Intervals: []PricingInterval{
			{MinTokens: 0, MaxTokens: resolverPtrInt(128000), InputPrice: resolverPtrFloat64(1e-6), OutputPrice: resolverPtrFloat64(2e-6)placeholder,
			{MinTokens: 128000, MaxTokens: nil, InputPrice: resolverPtrFloat64(3e-6), OutputPrice: resolverPtrFloat64(6e-6)placeholder,
	placeholder,
placeholder

	result := r.GetIntervalPricing(resolved, 50000)
	require.NotNil(t, result)
	require.InDelta(t, 1e-6, result.InputPricePerToken, 1e-12)
	require.InDelta(t, 2e-6, result.OutputPricePerToken, 1e-12)
	require.True(t, result.SupportsCacheBreakdown)

	result2 := r.GetIntervalPricing(resolved, 200000)
	require.NotNil(t, result2)
	require.InDelta(t, 3e-6, result2.InputPricePerToken, 1e-12)
placeholder

func TestGetIntervalPricing_NoMatch_FallsBackToBase(t *testing.T) {
	bs := newTestBillingServiceForResolver()
	r := NewModelPricingResolver(&ChannelService{placeholder, bs)

	basePricing := &ModelPricing{InputPricePerToken: 99e-6placeholder
	resolved := &ResolvedPricing{
		Mode:        BillingModeToken,
		BasePricing: basePricing,
		Intervals: []PricingInterval{
			{MinTokens: 10000, MaxTokens: resolverPtrInt(50000), InputPrice: resolverPtrFloat64(1e-6)placeholder,
	placeholder,
placeholder

	result := r.GetIntervalPricing(resolved, 5000)
	require.Equal(t, basePricing, result)
placeholder

func TestGetRequestTierPrice(t *testing.T) {
	bs := newTestBillingServiceForResolver()
	r := NewModelPricingResolver(&ChannelService{placeholder, bs)

	resolved := &ResolvedPricing{
		Mode: BillingModePerRequest,
		RequestTiers: []PricingInterval{
			{TierLabel: "1K", PerRequestPrice: resolverPtrFloat64(0.04)placeholder,
			{TierLabel: "2K", PerRequestPrice: resolverPtrFloat64(0.08)placeholder,
	placeholder,
placeholder

	require.InDelta(t, 0.04, r.GetRequestTierPrice(resolved, "1K"), 1e-12)
	require.InDelta(t, 0.08, r.GetRequestTierPrice(resolved, "2K"), 1e-12)
	require.InDelta(t, 0.0, r.GetRequestTierPrice(resolved, "4K"), 1e-12)
placeholder

func TestGetRequestTierPriceByContext(t *testing.T) {
	bs := newTestBillingServiceForResolver()
	r := NewModelPricingResolver(&ChannelService{placeholder, bs)

	resolved := &ResolvedPricing{
		Mode: BillingModePerRequest,
		RequestTiers: []PricingInterval{
			{MinTokens: 0, MaxTokens: resolverPtrInt(128000), PerRequestPrice: resolverPtrFloat64(0.05)placeholder,
			{MinTokens: 128000, MaxTokens: nil, PerRequestPrice: resolverPtrFloat64(0.10)placeholder,
	placeholder,
placeholder

	require.InDelta(t, 0.05, r.GetRequestTierPriceByContext(resolved, 50000), 1e-12)
	require.InDelta(t, 0.10, r.GetRequestTierPriceByContext(resolved, 200000), 1e-12)
placeholder

func TestGetRequestTierPrice_NilPerRequestPrice(t *testing.T) {
	bs := newTestBillingServiceForResolver()
	r := NewModelPricingResolver(&ChannelService{placeholder, bs)

	resolved := &ResolvedPricing{
		Mode: BillingModePerRequest,
		RequestTiers: []PricingInterval{
			{TierLabel: "1K", PerRequestPrice: nilplaceholder,
	placeholder,
placeholder

	require.InDelta(t, 0.0, r.GetRequestTierPrice(resolved, "1K"), 1e-12)
placeholder
