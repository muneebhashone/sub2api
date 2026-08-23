package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveBillingServiceTier(t *testing.T) {
	tests := []struct {
		name       string
		requested  string
		observed   string
		billing    string
		downgraded bool
placeholder{
		{name: "openai priority served as default", requested: "priority", observed: "default", billing: "default", downgraded: trueplaceholder,
		{name: "anthropic fast served as standard", requested: "fast", observed: "standard", billing: "standard", downgraded: trueplaceholder,
		{name: "priority honoured", requested: "priority", observed: "priority", billing: "priority"placeholder,
		{name: "no declaration keeps request", requested: "priority", observed: "", billing: "priority"placeholder,
		{name: "no request no declaration", requested: "", observed: "", billing: ""placeholder,
		{name: "response never raises the tier", requested: "", observed: "priority", billing: ""placeholder,
		{name: "flex never raised to default", requested: "flex", observed: "default", billing: "flex"placeholder,
		{name: "default echoed for untiered request", requested: "", observed: "default", billing: ""placeholder,
		{name: "unknown response tier ignored", requested: "priority", observed: "turbo", billing: "priority"placeholder,
		{name: "case and whitespace normalised", requested: " Priority ", observed: "DEFAULT", billing: "default", downgraded: trueplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveBillingServiceTier(tt.requested, tt.observed)
			require.Equal(t, tt.billing, got.Billing)
			require.Equal(t, tt.downgraded, got.Downgraded)
	placeholder)
placeholder
placeholder

func TestApplyServiceTierBillingResolutionOnlyRewritesDowngrades(t *testing.T) {
	t.Run("openai downgrade rewrites tier", func(t *testing.T) {
		requested := "priority"
		result := &OpenAIForwardResult{ServiceTier: &requested, UpstreamResponseServiceTier: "default"placeholder
		resolution := ApplyOpenAIServiceTierBillingResolution(result)
		require.True(t, resolution.Downgraded)
		require.NotNil(t, result.ServiceTier)
		require.Equal(t, "default", *result.ServiceTier)
placeholder)

	t.Run("openai honoured tier keeps pointer", func(t *testing.T) {
		requested := "priority"
		result := &OpenAIForwardResult{ServiceTier: &requested, UpstreamResponseServiceTier: "priority"placeholder
		require.False(t, ApplyOpenAIServiceTierBillingResolution(result).Downgraded)
		require.Same(t, &requested, result.ServiceTier)
placeholder)

	t.Run("openai untiered request stays nil", func(t *testing.T) {
		result := &OpenAIForwardResult{UpstreamResponseServiceTier: "priority"placeholder
		require.False(t, ApplyOpenAIServiceTierBillingResolution(result).Downgraded)
		require.Nil(t, result.ServiceTier)
placeholder)

	t.Run("anthropic standard speed rewrites fast", func(t *testing.T) {
		requested := "fast"
		result := &ForwardResult{ServiceTier: &requested, UpstreamResponseServiceTier: "standard"placeholder
		require.True(t, ApplyForwardServiceTierBillingResolution(result).Downgraded)
		require.Equal(t, "standard", *result.ServiceTier)
placeholder)

	t.Run("nil results are ignored", func(t *testing.T) {
		require.False(t, ApplyOpenAIServiceTierBillingResolution(nil).Downgraded)
		require.False(t, ApplyForwardServiceTierBillingResolution(nil).Downgraded)
placeholder)
placeholder
