package handler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApplySuggestedProfileToCompletionResponse(t *testing.T) {
	payload := map[string]any{
		"access_token": "token",
placeholder
	upstream := map[string]any{
		"suggested_display_name": "Alice",
		"suggested_avatar_url":   "https://cdn.example/avatar.png",
placeholder

	applySuggestedProfileToCompletionResponse(payload, upstream)

	require.Equal(t, "Alice", payload["suggested_display_name"])
	require.Equal(t, "https://cdn.example/avatar.png", payload["suggested_avatar_url"])
	require.Equal(t, true, payload["adoption_required"])
placeholder

func TestApplySuggestedProfileToCompletionResponseKeepsExistingPayloadValues(t *testing.T) {
	payload := map[string]any{
		"suggested_display_name": "Existing",
		"adoption_required":      false,
placeholder
	upstream := map[string]any{
		"suggested_display_name": "Alice",
		"suggested_avatar_url":   "https://cdn.example/avatar.png",
placeholder

	applySuggestedProfileToCompletionResponse(payload, upstream)

	require.Equal(t, "Existing", payload["suggested_display_name"])
	require.Equal(t, "https://cdn.example/avatar.png", payload["suggested_avatar_url"])
	require.Equal(t, true, payload["adoption_required"])
placeholder
