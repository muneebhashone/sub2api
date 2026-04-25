//go:build unit

package service

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestResolveRebateRatePercent_PerUserOverride verifies that per-inviter
// AffRebateRatePercent overrides the global rate, that NULL falls back to the
// global rate, and that out-of-range exclusive rates are clamped silently.
//
// SettingService is left nil here so globalRebateRatePercent returns the
// documented default (AffiliateRebateRateDefault = 20%) — this exercises the
// fallback path without spinning up a settings stub.
func TestResolveRebateRatePercent_PerUserOverride(t *testing.T) {
	t.Parallel()
	svc := &AffiliateService{placeholder

	// nil exclusive rate → falls back to global default (20%)
	require.InDelta(t, AffiliateRebateRateDefault,
		svc.resolveRebateRatePercent(context.Background(), &AffiliateSummary{placeholder), 1e-9)

	// exclusive rate set → overrides global
	rate := 50.0
	require.InDelta(t, 50.0,
		svc.resolveRebateRatePercent(context.Background(), &AffiliateSummary{AffRebateRatePercent: &rateplaceholder), 1e-9)

	// exclusive rate 0 → returns 0 (no rebate, intentional)
	zero := 0.0
	require.InDelta(t, 0.0,
		svc.resolveRebateRatePercent(context.Background(), &AffiliateSummary{AffRebateRatePercent: &zeroplaceholder), 1e-9)

	// exclusive rate above max → clamped to Max
	tooHigh := 250.0
	require.InDelta(t, AffiliateRebateRateMax,
		svc.resolveRebateRatePercent(context.Background(), &AffiliateSummary{AffRebateRatePercent: &tooHighplaceholder), 1e-9)

	// exclusive rate below min → clamped to Min
	tooLow := -5.0
	require.InDelta(t, AffiliateRebateRateMin,
		svc.resolveRebateRatePercent(context.Background(), &AffiliateSummary{AffRebateRatePercent: &tooLowplaceholder), 1e-9)
placeholder

// TestIsEnabled_NilSettingServiceReturnsDefault verifies that IsEnabled
// safely handles a nil settingService dependency by returning the default
// (off). This protects callers from nil-pointer crashes in misconfigured
// environments.
func TestIsEnabled_NilSettingServiceReturnsDefault(t *testing.T) {
	t.Parallel()
	svc := &AffiliateService{placeholder
	require.False(t, svc.IsEnabled(context.Background()))
	require.Equal(t, AffiliateEnabledDefault, svc.IsEnabled(context.Background()))
placeholder

// TestValidateExclusiveRate_BoundaryAndInvalid covers the validator used by
// admin-facing rate setters: nil is always valid (clear), in-range values
// are accepted, NaN/Inf and out-of-range values produce a typed BadRequest.
func TestValidateExclusiveRate_BoundaryAndInvalid(t *testing.T) {
	t.Parallel()
	require.NoError(t, validateExclusiveRate(nil))

	for _, v := range []float64{0, 0.01, 50, 99.99, 100placeholder {
		v := v
		require.NoError(t, validateExclusiveRate(&v), "value %v should be valid", v)
placeholder

	for _, v := range []float64{-0.01, 100.01, -100, 200placeholder {
		v := v
		require.Error(t, validateExclusiveRate(&v), "value %v should be rejected", v)
placeholder

	nan := math.NaN()
	require.Error(t, validateExclusiveRate(&nan))
	posInf := math.Inf(1)
	require.Error(t, validateExclusiveRate(&posInf))
	negInf := math.Inf(-1)
	require.Error(t, validateExclusiveRate(&negInf))
placeholder

func TestMaskEmail(t *testing.T) {
	t.Parallel()
	require.Equal(t, "a***@g***.com", maskEmail("alice@gmail.com"))
	require.Equal(t, "x***@d***", maskEmail("x@domain"))
	require.Equal(t, "", maskEmail(""))
placeholder

func TestIsValidAffiliateCodeFormat(t *testing.T) {
	t.Parallel()

	// 邀请码格式校验同时服务于：
	// 1) 系统自动生成的 12 位随机码（A-Z 去 I/O，2-9 去 0/1）
	// 2) 管理员设置的自定义专属码（如 "VIP2026"、"NEW_USER-1"）
	// 因此校验放宽到 [A-Z0-9_-]{4,32placeholder（要求调用方先 ToUpper）。
	cases := []struct {
		name string
		in   string
		want bool
placeholder{
		{"valid canonical 12-char", "ABCDEFGHJKLM", trueplaceholder,
		{"valid all digits 2-9", "234567892345", trueplaceholder,
		{"valid mixed", "A2B3C4D5E6F7", trueplaceholder,
		{"valid admin custom short", "VIP1", trueplaceholder,
		{"valid admin custom with hyphen", "NEW-USER", trueplaceholder,
		{"valid admin custom with underscore", "VIP_2026", trueplaceholder,
		{"valid 32-char max", "ABCDEFGHIJKLMNOPQRSTUVWXYZ012345", trueplaceholder,
		// Previously-excluded chars (I/O/0/1) are now allowed since admins may use them.
		{"letter I now allowed", "IBCDEFGHJKLM", trueplaceholder,
		{"letter O now allowed", "OBCDEFGHJKLM", trueplaceholder,
		{"digit 0 now allowed", "0BCDEFGHJKLM", trueplaceholder,
		{"digit 1 now allowed", "1BCDEFGHJKLM", trueplaceholder,
		{"too short (3 chars)", "ABC", falseplaceholder,
		{"too long (33 chars)", "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456", falseplaceholder,
		{"lowercase rejected (caller must ToUpper first)", "abcdefghjklm", falseplaceholder,
		{"empty", "", falseplaceholder,
		{"utf8 non-ascii", "ÄÄÄÄÄÄ", falseplaceholder, // bytes out of charset
		{"ascii punctuation .", "ABCDEFGHJK.M", falseplaceholder,
		{"whitespace", "ABCDEFGHJK M", falseplaceholder,
placeholder
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.want, isValidAffiliateCodeFormat(tc.in))
	placeholder)
placeholder
placeholder
