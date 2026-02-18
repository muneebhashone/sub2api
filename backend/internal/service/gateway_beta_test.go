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

func TestStripBetaToken(t *testing.T) {
	tests := []struct {
		name   string
		header string
		token  string
		want   string
placeholder{
		{
			name:   "token in middle",
			header: "oauth-2025-04-20,context-1m-2025-08-07,placeholder",
			token:  "context-1m-2025-08-07",
			want:   "oauth-2025-04-20,placeholder",
	placeholder,
		{
			name:   "token at start",
			header: "context-1m-2025-08-07,oauth-2025-04-20,placeholder",
			token:  "context-1m-2025-08-07",
			want:   "oauth-2025-04-20,placeholder",
	placeholder,
		{
			name:   "token at end",
			header: "oauth-2025-04-20,placeholder,context-1m-2025-08-07",
			token:  "context-1m-2025-08-07",
			want:   "oauth-2025-04-20,placeholder",
	placeholder,
		{
			name:   "token not present",
			header: "oauth-2025-04-20,placeholder",
			token:  "context-1m-2025-08-07",
			want:   "oauth-2025-04-20,placeholder",
	placeholder,
		{
			name:   "empty header",
			header: "",
			token:  "context-1m-2025-08-07",
			want:   "",
	placeholder,
		{
			name:   "with spaces",
			header: "oauth-2025-04-20, context-1m-2025-08-07 , placeholder",
			token:  "context-1m-2025-08-07",
			want:   "oauth-2025-04-20,placeholder",
	placeholder,
		{
			name:   "only token",
			header: "context-1m-2025-08-07",
			token:  "context-1m-2025-08-07",
			want:   "",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripBetaToken(tt.header, tt.token)
			require.Equal(t, tt.want, got)
	placeholder)
placeholder
placeholder

func TestMergeAnthropicBetaDropping_Context1M(t *testing.T) {
	required := []string{"oauth-2025-04-20", "placeholder"placeholder
	incoming := "context-1m-2025-08-07,foo-beta,oauth-2025-04-20"
	drop := map[string]struct{placeholder{"context-1m-2025-08-07": {placeholderplaceholder

	got := mergeAnthropicBetaDropping(required, incoming, drop)
	require.Equal(t, "oauth-2025-04-20,placeholder,foo-beta", got)
	require.NotContains(t, got, "context-1m-2025-08-07")
placeholder
