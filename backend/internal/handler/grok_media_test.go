package handler

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

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
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, grokMediaRequiredCapability(tt.endpoint))
	placeholder)
placeholder
placeholder

func TestGrokMediaScheduleModelUsesNormalizedMappedUpstream(t *testing.T) {
	account := &service.Account{
		Platform: service.PlatformGrok,
placeholder
			"model_mapping": map[string]any{
				"grok-imagine-video-1.5": "wrong-raw-model",
				"grok-imagine-video":     "mapped-video-model",
		placeholder,
	placeholder,
placeholder

	require.Equal(t, "mapped-video-model", grokMediaScheduleModel(account, "grok-imagine-video", nil))
	require.Equal(t, "actual-upstream-model", grokMediaScheduleModel(account, "grok-imagine-video", &service.OpenAIForwardResult{
		UpstreamModel: "actual-upstream-model",
placeholder))
	require.Equal(t, "mapped-video-model", grokMediaScheduleModel(account, "grok-imagine-video", &service.OpenAIForwardResult{placeholder))
	require.Equal(t, "grok-imagine-video", grokMediaScheduleModel(nil, " grok-imagine-video ", nil))
placeholder
