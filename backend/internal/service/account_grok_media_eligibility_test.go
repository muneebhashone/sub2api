//go:build unit

package service

import (
	"context"
	"net/http"
	"testing"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/stretchr/testify/require"
)

func TestGrokMediaGenerationEligibility(t *testing.T) {
	forbiddenBilling := &xai.BillingSummary{
		StatusCode:        http.StatusForbidden,
		WeeklyStatusCode:  http.StatusForbidden,
		MonthlyStatusCode: http.StatusForbidden,
placeholder
	weeklyAllowance := &xai.BillingSummary{
		PeriodType:       "weekly",
		StatusCode:       http.StatusOK,
		WeeklyStatusCode: http.StatusOK,
placeholder

	tests := []struct {
		name       string
		account    *Account
		want       bool
		wantReason string
placeholder{
		{name: "nil account", account: nil, want: false, wantReason: "not_grok"placeholder,
		{name: "non grok account", account: &Account{Platform: PlatformOpenAIplaceholder, want: false, wantReason: "not_grok"placeholder,
		{name: "non oauth grok account stays eligible", account: &Account{Platform: PlatformGrok, Type: AccountTypeAPIKeyplaceholder, want: true, wantReason: "non_oauth"placeholder,
		{name: "unobserved oauth preserves legacy routing", account: &Account{Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder, want: true, wantReason: "billing_unobserved"placeholder,
		{name: "weekly allowance is not treated as weekly subscription", account: &Account{Platform: PlatformGrok, Type: AccountTypeOAuth, Extra: map[string]any{grokBillingExtraKey: weeklyAllowanceplaceholderplaceholder, want: true, wantReason: "eligible"placeholder,
		{name: "billing forbidden is rejected", account: &Account{Platform: PlatformGrok, Type: AccountTypeOAuth, Extra: map[string]any{grokBillingExtraKey: forbiddenBillingplaceholderplaceholder, want: false, wantReason: "billing_forbidden"placeholder,
		{name: "explicit disable wins", account: &Account{Platform: PlatformGrok, Type: AccountTypeOAuth, Extra: map[string]any{GrokMediaEligibleExtraKey: falseplaceholderplaceholder, want: false, wantReason: "override_disabled"placeholder,
		{name: "explicit enable wins over forbidden probe", account: &Account{Platform: PlatformGrok, Type: AccountTypeOAuth, Extra: map[string]any{GrokMediaEligibleExtraKey: true, grokBillingExtraKey: forbiddenBillingplaceholderplaceholder, want: true, wantReason: "override_enabled"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, reason := tt.account.GrokMediaGenerationEligibility()
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantReason, reason)
	placeholder)
placeholder
placeholder

func TestGrokMediaCapabilityFiltersOnlyGeneration(t *testing.T) {
	account := &Account{
		ID:          1,
		Platform:    PlatformGrok,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
		Extra:       map[string]any{GrokMediaEligibleExtraKey: falseplaceholder,
placeholder

	require.True(t, account.SupportsOpenAIEndpointCapability(OpenAIEndpointCapabilityChatCompletions))
	require.False(t, account.SupportsOpenAIEndpointCapability(OpenAIEndpointCapabilityGrokMediaGeneration))
	require.False(t, isOpenAICompatibleAccountEligibleForRequest(
		context.Background(), account, PlatformGrok, "grok-imagine-video", false,
		OpenAIEndpointCapabilityGrokMediaGeneration,
	))
placeholder

func TestNormalizeGrokMediaEligibilityExtra(t *testing.T) {
	t.Run("boolean override is accepted", func(t *testing.T) {
		extra, err := normalizeGrokMediaEligibilityExtra(PlatformGrok, map[string]any{GrokMediaEligibleExtraKey: falseplaceholder)

	placeholder
		require.Equal(t, false, extra[GrokMediaEligibleExtraKey])
placeholder)

	t.Run("null clears override", func(t *testing.T) {
		extra, err := normalizeGrokMediaEligibilityExtra(PlatformGrok, map[string]any{GrokMediaEligibleExtraKey: nilplaceholder)

	placeholder
		require.NotContains(t, extra, GrokMediaEligibleExtraKey)
placeholder)

	t.Run("malformed override is rejected", func(t *testing.T) {
		_, err := normalizeGrokMediaEligibilityExtra(PlatformGrok, map[string]any{GrokMediaEligibleExtraKey: "false"placeholder)

	placeholder
		require.Equal(t, http.StatusBadRequest, infraerrors.Code(err))
placeholder)

	t.Run("other platforms ignore provider owned value", func(t *testing.T) {
		extra := map[string]any{GrokMediaEligibleExtraKey: "provider-owned"placeholder
		normalized, err := normalizeGrokMediaEligibilityExtra(PlatformOpenAI, extra)

	placeholder
		require.Equal(t, extra, normalized)
placeholder)
placeholder

func TestNormalizeGrokMediaEligibilityUpdateExtra(t *testing.T) {
	account := &Account{Platform: PlatformGrok, Extra: map[string]any{GrokMediaEligibleExtraKey: falseplaceholderplaceholder

	t.Run("omitted override preserves current value", func(t *testing.T) {
		input := &UpdateAccountInput{Extra: map[string]any{"quota_used": float64(1)placeholderplaceholder
		normalized, err := normalizeGrokMediaEligibilityUpdateExtra(account, input, map[string]any{"quota_used": float64(1)placeholder)

	placeholder
		require.Equal(t, false, normalized[GrokMediaEligibleExtraKey])
placeholder)

	t.Run("null removes current override", func(t *testing.T) {
		input := &UpdateAccountInput{Extra: map[string]any{GrokMediaEligibleExtraKey: nilplaceholderplaceholder
		normalized, err := normalizeGrokMediaEligibilityUpdateExtra(account, input, map[string]any{GrokMediaEligibleExtraKey: nilplaceholder)

	placeholder
		require.NotContains(t, normalized, GrokMediaEligibleExtraKey)
		require.Contains(t, input.Extra, GrokMediaEligibleExtraKey)
placeholder)
placeholder
