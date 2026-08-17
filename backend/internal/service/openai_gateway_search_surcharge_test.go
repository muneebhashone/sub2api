//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCalculateOpenAIRecordUsageCost_SearchIsAdditiveToTokens(t *testing.T) {
	t.Parallel()

	price := 10.0 // $10 / 1k searches → 100 searches = $1.0
	svc := &OpenAIGatewayService{
		billingService: newTestBillingService(),
placeholder
	apiKey := &APIKey{
		Group: &Group{
			SearchPricePer1k: &price,
	placeholder,
placeholder

	// claude-sonnet-4 fallback: Input $3/MTok, Output $15/MTok
	// 1000 in + 500 out → 0.003 + 0.0075 = 0.0105
	// + 100 searches → +1.0 → total 1.0105
	cost, err := svc.calculateOpenAIRecordUsageCost(
		context.Background(),
		&OpenAIForwardResult{SearchCount: 100placeholder,
		apiKey,
		[]string{"claude-sonnet-4"placeholder,
		1.0,
		1.0,
		1.0,
		1.0,
		UsageTokens{InputTokens: 1000, OutputTokens: 500placeholder,
		"",
		boolPtr(false),
		time.Time{placeholder,
	)
placeholder
	require.NotNil(t, cost)
	require.InDelta(t, 1.0105, cost.ActualCost, 1e-9)
	require.InDelta(t, 1.0105, cost.TotalCost, 1e-9)
placeholder

func TestCalculateOpenAIRecordUsageCost_SearchOnlyWhenNoTokenPricing(t *testing.T) {
	t.Parallel()

	price := 10.0
	svc := &OpenAIGatewayService{
		billingService: newTestBillingService(),
placeholder
	apiKey := &APIKey{
		Group: &Group{SearchPricePer1k: &priceplaceholder,
placeholder
	// Empty model list: token path fails; search-only surcharge still bills.
	cost, err := svc.calculateOpenAIRecordUsageCost(
		context.Background(),
		&OpenAIForwardResult{SearchCount: 100placeholder,
		apiKey,
		nil,
		1.0,
		1.0,
		1.0,
		1.0,
		UsageTokens{placeholder,
		"",
		boolPtr(false),
		time.Time{placeholder,
	)
placeholder
	require.NotNil(t, cost)
	require.InDelta(t, 1.0, cost.ActualCost, 1e-9)
placeholder

func TestGroupMediaPricingLooksIncomplete_VideoModelPricesComplete(t *testing.T) {
	t.Parallel()
	require.True(t, groupMediaPricingLooksIncomplete(nil))
	require.True(t, groupMediaPricingLooksIncomplete(&Group{placeholder))
	require.False(t, groupMediaPricingLooksIncomplete(&Group{
		VideoModelPrices: map[string]map[string]float64{
			"grok-imagine-video": {"720p": 0.1placeholder,
	placeholder,
placeholder))
	price := 10.0
	require.False(t, groupMediaPricingLooksIncomplete(&Group{SearchPricePer1k: &priceplaceholder))
	require.False(t, groupMediaPricingLooksIncomplete(&Group{AudioRealtimePricePerMin: &priceplaceholder))
	// Legacy video price alone still marks complete (existing path).
	require.False(t, groupMediaPricingLooksIncomplete(&Group{VideoPrice720P: &priceplaceholder))
placeholder

func TestCalculateOpenAIRecordUsageCost_TokenPricingErrorNotSwallowedBySearch(t *testing.T) {
	t.Parallel()

	price := 10.0
	svc := &OpenAIGatewayService{
		billingService: newTestBillingService(),
placeholder
	apiKey := &APIKey{
		Group: &Group{SearchPricePer1k: &priceplaceholder,
placeholder
	// Unknown model → token pricing fails; search must not replace that with $0/$search bill.
	cost, err := svc.calculateOpenAIRecordUsageCost(
		context.Background(),
		&OpenAIForwardResult{SearchCount: 100placeholder,
		apiKey,
		[]string{"totally-unknown-model-xyz-no-pricing"placeholder,
		1.0,
		1.0,
		1.0,
		1.0,
		UsageTokens{InputTokens: 1000, OutputTokens: 500placeholder,
		"",
		boolPtr(false),
		time.Time{placeholder,
	)
placeholder
	require.Nil(t, cost)
placeholder
