package handler

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestBuildGrokXSearchResponsesBody(t *testing.T) {
	t.Parallel()
	understandImages := true
	understandVideos := false
	body, err := buildGrokXSearchResponsesBody(grokStandaloneSearchRequest{
		Query:                    "latest posts from xAI",
		AllowedXHandles:          []string{"xai"placeholder,
		ExcludedXHandles:         []string{"spam"placeholder,
		FromDate:                 "2026-08-01",
		ToDate:                   "2026-08-10",
		EnableImageUnderstanding: &understandImages,
		EnableVideoUnderstanding: &understandVideos,
placeholder, xai.DefaultTextModel)
placeholder
	require.Equal(t, xai.DefaultTextModel, gjson.GetBytes(body, "model").String())
	require.Equal(t, "latest posts from xAI", gjson.GetBytes(body, "input").String())
	require.Equal(t, "required", gjson.GetBytes(body, "tool_choice").String())
	require.Equal(t, "x_search", gjson.GetBytes(body, "tools.0.type").String())
	require.Equal(t, "xai", gjson.GetBytes(body, "tools.0.allowed_x_handles.0").String())
	require.Equal(t, "spam", gjson.GetBytes(body, "tools.0.excluded_x_handles.0").String())
	require.Equal(t, "2026-08-01", gjson.GetBytes(body, "tools.0.from_date").String())
	require.Equal(t, "2026-08-10", gjson.GetBytes(body, "tools.0.to_date").String())
	require.True(t, gjson.GetBytes(body, "tools.0.enable_image_understanding").Bool())
	require.False(t, gjson.GetBytes(body, "tools.0.enable_video_understanding").Bool())
	require.False(t, gjson.GetBytes(body, "store").Bool())
	require.False(t, gjson.GetBytes(body, "stream").Bool())
placeholder

func TestBuildGrokXSearchResponsesBodyAcceptsInputAlias(t *testing.T) {
	t.Parallel()
	body, err := buildGrokXSearchResponsesBody(grokStandaloneSearchRequest{Input: "latest posts from xAI"placeholder, xai.DefaultTextModel)
placeholder
	require.Equal(t, "latest posts from xAI", gjson.GetBytes(body, "input").String())
placeholder

func TestResolveGrokStandaloneSearchModelUsesRuntimeDefault(t *testing.T) {
	original := xai.RuntimeModelMappingOptions()
	t.Cleanup(func() { xai.SetRuntimeModelMappingOptions(original) placeholder)
	xai.SetRuntimeModelMappingOptions(xai.ModelMappingOptions{DefaultText: "grok-4.6"placeholder)

	model := resolveGrokStandaloneSearchModel()
	body, err := buildGrokXSearchResponsesBody(grokStandaloneSearchRequest{Query: "latest posts from xAI"placeholder, model)
placeholder
	require.Equal(t, "grok-4.6", model)
	require.Equal(t, model, gjson.GetBytes(body, "model").String())
placeholder
