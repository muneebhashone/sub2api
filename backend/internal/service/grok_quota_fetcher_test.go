//go:build unit

package service

import (
	"net/http"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/stretchr/testify/require"
)

func grokInt64PtrForTest(v int64) *int64 { return &v placeholder
func grokIntPtrForTest(v int) *int       { return &v placeholder

func TestGrokQuotaFetcherBuildUsageInfoUnknownUntilFirstSnapshot(t *testing.T) {
	t.Parallel()

	usage := NewGrokQuotaFetcher().BuildUsageInfo(&Account{Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder)
	require.Equal(t, "passive", usage.Source)
	require.Equal(t, xai.GrokFreeRolling24hTokenLimit, usage.GrokFreeTokenLimit)
	require.Equal(t, "quota_unknown", usage.ErrorCode)
	require.Contains(t, usage.Error, "unknown until billing is probed")
placeholder

func TestGrokQuotaFetcherPrefersLiveJWTTierOverStaleBillingPlan(t *testing.T) {
	t.Parallel()

	account := &Account{
		Platform: PlatformGrok,
		Type:     AccountTypeOAuth,
placeholder
			"access_token":      makeGrokOAuthJWT(map[string]any{"tier": 0placeholder),
			"subscription_tier": "supergrok_heavy",
	placeholder,
		Extra: map[string]any{
			grokBillingExtraKey: &xai.BillingSummary{
				Plan:       "SuperGrok Heavy",
				StatusCode: http.StatusOK,
				UpdatedAt:  "2030-01-01T00:00:00Z",
		placeholder,
	placeholder,
placeholder

	usage := NewGrokQuotaFetcher().BuildUsageInfo(account)
	require.Equal(t, "free", usage.SubscriptionTier)
	require.Equal(t, "free", usage.SubscriptionTierRaw)
placeholder

func TestGrokQuotaFetcherUsesCredentialTierWhenBillingHasNoPlan(t *testing.T) {
	t.Parallel()

	account := &Account{
		Platform: PlatformGrok,
		Type:     AccountTypeOAuth,
placeholder
			"subscription_tier":  " FREE ",
			"entitlement_status": " active ",
	placeholder,
		Extra: map[string]any{
			grokBillingExtraKey: &xai.BillingSummary{
				PeriodType: "weekly",
				StatusCode: http.StatusOK,
				UpdatedAt:  "2030-01-01T00:00:00Z",
		placeholder,
	placeholder,
placeholder

	usage := NewGrokQuotaFetcher().BuildUsageInfo(account)

	require.NotNil(t, usage.GrokBilling)
	require.Equal(t, "FREE", usage.SubscriptionTier)
	require.Equal(t, "FREE", usage.SubscriptionTierRaw)
	require.Equal(t, "active", usage.GrokEntitlementStatus)
placeholder

func TestGrokQuotaFetcherBuildUsageInfoFromSnapshot(t *testing.T) {
	t.Parallel()

	updatedAt := "2030-01-01T00:00:00Z"
	account := &Account{
		Platform: PlatformGrok,
		Type:     AccountTypeOAuth,
		Extra: map[string]any{
			grokQuotaSnapshotExtraKey: &xai.QuotaSnapshot{
				Requests: &xai.QuotaWindow{
					Limit:     grokInt64PtrForTest(100),
					Remaining: grokInt64PtrForTest(12),
					ResetAt:   updatedAt,
			placeholder,
				Tokens: &xai.QuotaWindow{
					Limit:     grokInt64PtrForTest(1000),
					Remaining: grokInt64PtrForTest(900),
			placeholder,
				RetryAfterSeconds: grokIntPtrForTest(30),
				SubscriptionTier:  "supergrok",
				EntitlementStatus: "active",
				StatusCode:        http.StatusTooManyRequests,
				LastProbeAt:       updatedAt,
				LastHeadersSeenAt: updatedAt,
				UpdatedAt:         updatedAt,
		placeholder,
	placeholder,
placeholder

	usage := NewGrokQuotaFetcher().BuildUsageInfo(account)
	require.Equal(t, "passive", usage.Source)
	require.Equal(t, "rate_limited", usage.ErrorCode)
	require.Equal(t, "observed", usage.GrokQuotaSnapshotState)
	require.Equal(t, "supergrok", usage.SubscriptionTier)
	require.Equal(t, "active", usage.GrokEntitlementStatus)
	require.Equal(t, int64(100), *usage.GrokRequestQuota.Limit)
	require.Equal(t, int64(12), *usage.GrokRequestQuota.Remaining)
	require.Equal(t, 30, *usage.GrokRetryAfterSeconds)
	require.NotNil(t, usage.UpdatedAt)
	require.Equal(t, updatedAt, usage.GrokLastQuotaProbeAt)
	require.Equal(t, updatedAt, usage.GrokLastHeadersSeenAt)
	require.Equal(t, http.StatusTooManyRequests, usage.GrokLastStatusCode)
	require.True(t, usage.UpdatedAt.Equal(time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)))
placeholder

func TestGrokQuotaFetcherSnapshotErrorOverridesSuccessfulBillingStatus(t *testing.T) {
	t.Parallel()

	updatedAt := "2030-01-01T00:00:00Z"
	account := &Account{
		Platform: PlatformGrok,
		Type:     AccountTypeOAuth,
		Extra: map[string]any{
			grokBillingExtraKey: &xai.BillingSummary{
				PeriodType: "weekly",
				StatusCode: http.StatusOK,
				UpdatedAt:  updatedAt,
		placeholder,
			grokQuotaSnapshotExtraKey: &xai.QuotaSnapshot{
				StatusCode: http.StatusTooManyRequests,
				UpdatedAt:  updatedAt,
		placeholder,
	placeholder,
placeholder

	usage := NewGrokQuotaFetcher().BuildUsageInfo(account)

	require.Equal(t, "rate_limited", usage.ErrorCode)
	require.Equal(t, http.StatusTooManyRequests, usage.GrokLastStatusCode)
placeholder

func TestGrokQuotaFetcherNewerSuccessfulActiveProbeClearsBillingForbidden(t *testing.T) {
	t.Parallel()

	billingAt := "2030-01-01T00:00:00Z"
	probeAt := "2030-01-01T00:05:00Z"
	account := &Account{
		Platform: PlatformGrok,
		Type:     AccountTypeOAuth,
placeholder
			"entitlement_status": "forbidden",
	placeholder,
		Extra: map[string]any{
			grokBillingExtraKey: &xai.BillingSummary{
				StatusCode: http.StatusForbidden,
				UpdatedAt:  billingAt,
		placeholder,
			grokQuotaSnapshotExtraKey: &xai.QuotaSnapshot{
				StatusCode:        http.StatusOK,
				ObservationSource: "active_probe",
				LastProbeAt:       probeAt,
				UpdatedAt:         probeAt,
		placeholder,
	placeholder,
placeholder

	usage := NewGrokQuotaFetcher().BuildUsageInfo(account)

	require.False(t, usage.IsForbidden)
	require.Empty(t, usage.ForbiddenType)
	require.Empty(t, usage.ErrorCode)
	require.Empty(t, usage.GrokEntitlementStatus)
	require.Equal(t, http.StatusOK, usage.GrokLastStatusCode)
	require.Equal(t, probeAt, usage.GrokLastQuotaProbeAt)
	require.NotNil(t, usage.UpdatedAt)
	require.True(t, usage.UpdatedAt.Equal(time.Date(2030, 1, 1, 0, 5, 0, 0, time.UTC)))
placeholder

func TestGrokQuotaFetcherSameSecondSuccessfulActiveProbeClearsBillingForbidden(t *testing.T) {
	t.Parallel()

	observedAt := "2030-01-01T00:05:00Z"
	account := &Account{
		Platform: PlatformGrok,
		Type:     AccountTypeOAuth,
		Extra: map[string]any{
			grokBillingExtraKey: &xai.BillingSummary{
				StatusCode: http.StatusForbidden,
				UpdatedAt:  observedAt,
		placeholder,
			grokQuotaSnapshotExtraKey: &xai.QuotaSnapshot{
				StatusCode:        http.StatusOK,
				ObservationSource: "active_probe",
				LastProbeAt:       observedAt,
				UpdatedAt:         observedAt,
		placeholder,
	placeholder,
placeholder

	usage := NewGrokQuotaFetcher().BuildUsageInfo(account)

	require.False(t, usage.IsForbidden)
	require.Empty(t, usage.ForbiddenType)
	require.Empty(t, usage.ErrorCode)
	require.Equal(t, http.StatusOK, usage.GrokLastStatusCode)
placeholder

func TestGrokQuotaFetcherDoesNotClearBillingForbiddenWithoutNewerSuccessfulActiveProbe(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		snapshot xai.QuotaSnapshot
placeholder{
		{
			name: "older active probe",
			snapshot: xai.QuotaSnapshot{
				StatusCode:        http.StatusOK,
				ObservationSource: "active_probe",
				LastProbeAt:       "2030-01-01T00:04:59Z",
				UpdatedAt:         "2030-01-01T00:04:59Z",
		placeholder,
	placeholder,
		{
			name: "newer passive response",
			snapshot: xai.QuotaSnapshot{
				StatusCode:        http.StatusOK,
				ObservationSource: "upstream_response",
				UpdatedAt:         "2030-01-01T00:05:01Z",
		placeholder,
	placeholder,
		{
			name: "newer failed active probe",
			snapshot: xai.QuotaSnapshot{
				StatusCode:        http.StatusTooManyRequests,
				ObservationSource: "active_probe",
				LastProbeAt:       "2030-01-01T00:05:01Z",
				UpdatedAt:         "2030-01-01T00:05:01Z",
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			account := &Account{
				Platform: PlatformGrok,
				Type:     AccountTypeOAuth,
				Extra: map[string]any{
					grokBillingExtraKey: &xai.BillingSummary{
						StatusCode: http.StatusForbidden,
						UpdatedAt:  "2030-01-01T00:05:00Z",
				placeholder,
					grokQuotaSnapshotExtraKey: tt.snapshot,
			placeholder,
		placeholder

			usage := NewGrokQuotaFetcher().BuildUsageInfo(account)

			require.True(t, usage.IsForbidden)
			require.Equal(t, "forbidden", usage.ForbiddenType)
			require.Equal(t, "forbidden", usage.ErrorCode)
	placeholder)
placeholder
placeholder

func TestGrokQuotaFetcherBuildUsageInfoFromNoHeadersProbe(t *testing.T) {
	t.Parallel()

	probedAt := "2030-01-01T00:00:00Z"
	account := &Account{
		Platform: PlatformGrok,
		Type:     AccountTypeOAuth,
		Extra: map[string]any{
			grokQuotaSnapshotExtraKey: xai.QuotaSnapshot{
				StatusCode:        http.StatusOK,
				HeadersObserved:   false,
				ObservationSource: "active_probe",
				LastProbeAt:       probedAt,
				UpdatedAt:         probedAt,
		placeholder,
	placeholder,
placeholder

	usage := NewGrokQuotaFetcher().BuildUsageInfo(account)
	require.Equal(t, "quota_unknown", usage.ErrorCode)
	require.Equal(t, "no_headers", usage.GrokQuotaSnapshotState)
	require.Contains(t, usage.Error, "No xAI quota headers observed")
	require.Equal(t, probedAt, usage.GrokLastQuotaProbeAt)
	require.Empty(t, usage.GrokLastHeadersSeenAt)
	require.Equal(t, http.StatusOK, usage.GrokLastStatusCode)
	require.Nil(t, usage.GrokRequestQuota)
	require.Nil(t, usage.GrokTokenQuota)
placeholder

func TestGrokQuotaFetcherClassifiesForbiddenAndReauth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		statusCode  int
		wantReauth  bool
		wantForbid  bool
		wantCode    string
		wantEntitle string
placeholder{
		{name: "reauth", statusCode: http.StatusUnauthorized, wantReauth: true, wantCode: "unauthenticated"placeholder,
		{name: "forbidden", statusCode: http.StatusForbidden, wantForbid: true, wantCode: "forbidden", wantEntitle: "forbidden"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			account := &Account{
				Platform: PlatformGrok,
				Type:     AccountTypeOAuth,
				Extra: map[string]any{
					grokQuotaSnapshotExtraKey: xai.QuotaSnapshot{
						StatusCode:      tt.statusCode,
						HeadersObserved: true,
						UpdatedAt:       time.Now().UTC().Format(time.RFC3339),
				placeholder,
			placeholder,
		placeholder
			usage := NewGrokQuotaFetcher().BuildUsageInfo(account)
			require.Equal(t, tt.wantReauth, usage.NeedsReauth)
			require.Equal(t, tt.wantForbid, usage.IsForbidden)
			require.Equal(t, tt.wantCode, usage.ErrorCode)
			require.Equal(t, tt.wantEntitle, usage.GrokEntitlementStatus)
	placeholder)
placeholder
placeholder
