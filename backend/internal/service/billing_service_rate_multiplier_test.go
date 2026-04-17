//go:build unit

package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCalculateCost_RateMultiplier_NegativeClampedToZero 锁定负数倍率被
// 钳制为 0（而非历史上的 1.0），避免配置异常导致静默按标准价扣费。
func TestCalculateCost_RateMultiplier_NegativeClampedToZero(t *testing.T) {
	svc := newTestBillingService()
	tokens := UsageTokens{InputTokens: 1000, OutputTokens: 500placeholder

	tests := []struct {
		name       string
		multiplier float64
		wantRatio  float64 // ActualCost / TotalCost
placeholder{
		{"negative clamped to 0", -1.5, 0placeholder,
		{"zero passes through as 0 (defense in depth)", 0, 0placeholder,
		{"positive 2x applied", 2.0, 2.0placeholder,
		{"positive 0.5x applied", 0.5, 0.5placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost, err := svc.CalculateCost("claude-sonnet-4", tokens, tt.multiplier)
		placeholder
			require.Greater(t, cost.TotalCost, 0.0, "TotalCost should be non-zero")
			require.InDelta(t, tt.wantRatio*cost.TotalCost, cost.ActualCost, 1e-9)
	placeholder)
placeholder
placeholder

// TestCalculateImageCost_RateMultiplier_NegativeClampedToZero 图片按次计费路径
// 同样遵循"负数 → 0"语义。
func TestCalculateImageCost_RateMultiplier_NegativeClampedToZero(t *testing.T) {
	svc := newTestBillingService()
	price := 0.04
	cfg := &ImagePriceConfig{Price1K: &priceplaceholder

	tests := []struct {
		name       string
		multiplier float64
		wantRatio  float64
placeholder{
		{"negative clamped to 0", -0.5, 0placeholder,
		{"zero passes through", 0, 0placeholder,
		{"positive 3x applied", 3.0, 3.0placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := svc.CalculateImageCost("imagen-3", "1K", 2, cfg, tt.multiplier)
			require.NotNil(t, cost)
			require.Greater(t, cost.TotalCost, 0.0)
			require.InDelta(t, tt.wantRatio*cost.TotalCost, cost.ActualCost, 1e-9)
	placeholder)
placeholder
placeholder
