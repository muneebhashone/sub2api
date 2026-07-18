package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

type grokMediaEligibilityProberStub struct {
	eligible bool
	reason   string
	err      error
	calls    int
placeholder

func (s *grokMediaEligibilityProberStub) ProbeMediaEligibility(context.Context, int64) (bool, string, error) {
	s.calls++
	return s.eligible, s.reason, s.err
placeholder

func TestShouldRecordGrokMediaUsage(t *testing.T) {
	tests := []struct {
		name     string
		endpoint service.GrokMediaEndpoint
		model    string
		want     bool
placeholder{
		{
			name:     "image generation records usage",
			endpoint: service.GrokMediaEndpointImagesGenerations,
			model:    "grok-imagine",
			want:     true,
	placeholder,
		{
			name:     "image edit records usage",
			endpoint: service.GrokMediaEndpointImagesEdits,
			model:    "grok-imagine-edit",
			want:     true,
	placeholder,
		{
			name:     "video generation records usage",
			endpoint: service.GrokMediaEndpointVideosGenerations,
			model:    "grok-imagine-video-1.5",
			want:     true,
	placeholder,
		{
			name:     "video status skips empty model usage",
			endpoint: service.GrokMediaEndpointVideoStatus,
			model:    "",
			want:     false,
	placeholder,
		{
			name:     "video content skips usage",
			endpoint: service.GrokMediaEndpointVideoContent,
			model:    "",
			want:     false,
	placeholder,
		{
			name:     "generation skips usage without model",
			endpoint: service.GrokMediaEndpointImagesGenerations,
			model:    " ",
			want:     false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, shouldRecordGrokMediaUsage(tt.endpoint, tt.model))
	placeholder)
placeholder
placeholder

func TestGrokMediaRequiredCapability(t *testing.T) {
	tests := []struct {
		name     string
		endpoint service.GrokMediaEndpoint
		want     service.OpenAIEndpointCapability
placeholder{
		{name: "image generation", endpoint: service.GrokMediaEndpointImagesGenerations, want: service.OpenAIEndpointCapabilityGrokMediaGenerationplaceholder,
		{name: "image edit", endpoint: service.GrokMediaEndpointImagesEdits, want: service.OpenAIEndpointCapabilityGrokMediaGenerationplaceholder,
		{name: "video generation", endpoint: service.GrokMediaEndpointVideosGenerations, want: service.OpenAIEndpointCapabilityGrokMediaGenerationplaceholder,
		{name: "video edit", endpoint: service.GrokMediaEndpointVideosEdits, want: service.OpenAIEndpointCapabilityGrokMediaGenerationplaceholder,
		{name: "video extension", endpoint: service.GrokMediaEndpointVideosExtensions, want: service.OpenAIEndpointCapabilityGrokMediaGenerationplaceholder,
		{name: "video status preserves lookup", endpoint: service.GrokMediaEndpointVideoStatus, want: ""placeholder,
		{name: "video content preserves lookup", endpoint: service.GrokMediaEndpointVideoContent, want: ""placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, grokMediaRequiredCapability(tt.endpoint))
	placeholder)
placeholder
placeholder

func TestEnsureGrokMediaAccountEligibility(t *testing.T) {
	t.Run("non oauth account does not probe", func(t *testing.T) {
		prober := &grokMediaEligibilityProberStub{placeholder
		h := &OpenAIGatewayHandler{grokMediaEligibilityProber: proberplaceholder
		account := &service.Account{Platform: service.PlatformGrok, Type: service.AccountTypeAPIKeyplaceholder

		eligible, reason, err := h.ensureGrokMediaAccountEligibility(context.Background(), account)

	placeholder
		require.True(t, eligible)
		require.Equal(t, "non_oauth", reason)
		require.Zero(t, prober.calls)
placeholder)

	t.Run("unobserved oauth is probed before forwarding", func(t *testing.T) {
		prober := &grokMediaEligibilityProberStub{eligible: true, reason: "eligible"placeholder
		h := &OpenAIGatewayHandler{grokMediaEligibilityProber: proberplaceholder
		account := &service.Account{ID: 7, Platform: service.PlatformGrok, Type: service.AccountTypeOAuthplaceholder

		eligible, reason, err := h.ensureGrokMediaAccountEligibility(context.Background(), account)

	placeholder
		require.True(t, eligible)
		require.Equal(t, "eligible", reason)
		require.Equal(t, 1, prober.calls)
placeholder)

	t.Run("missing prober fails closed", func(t *testing.T) {
		h := &OpenAIGatewayHandler{placeholder
		account := &service.Account{ID: 8, Platform: service.PlatformGrok, Type: service.AccountTypeOAuthplaceholder

		eligible, reason, err := h.ensureGrokMediaAccountEligibility(context.Background(), account)

	placeholder
		require.False(t, eligible)
		require.Equal(t, "billing_probe_unavailable", reason)
placeholder)

	t.Run("probe failure fails closed", func(t *testing.T) {
		probeErr := errors.New("probe failed")
		prober := &grokMediaEligibilityProberStub{reason: "billing_unobserved", err: probeErrplaceholder
		h := &OpenAIGatewayHandler{grokMediaEligibilityProber: proberplaceholder
		account := &service.Account{ID: 9, Platform: service.PlatformGrok, Type: service.AccountTypeOAuthplaceholder

		eligible, reason, err := h.ensureGrokMediaAccountEligibility(context.Background(), account)

		require.ErrorIs(t, err, probeErr)
		require.False(t, eligible)
		require.Equal(t, "billing_unobserved", reason)
placeholder)
placeholder
