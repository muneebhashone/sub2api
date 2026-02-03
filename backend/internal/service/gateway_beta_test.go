package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMergeAnthropicBeta(t *testing.T) {
	got := mergeAnthropicBeta(
		[]string{"oauth-2025-04-20", "placeholder"placeholder,
		"foo, oauth-2025-04-20,bar, foo",
	)
	require.Equal(t, "oauth-2025-04-20,placeholder,foo,bar", got)
placeholder

func TestMergeAnthropicBeta_EmptyIncoming(t *testing.T) {
	got := mergeAnthropicBeta(
		[]string{"oauth-2025-04-20", "placeholder"placeholder,
		"",
	)
	require.Equal(t, "oauth-2025-04-20,placeholder", got)
placeholder
