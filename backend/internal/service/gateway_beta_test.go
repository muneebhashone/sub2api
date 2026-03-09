package service

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/claude"

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

func TestStripBetaTokens(t *testing.T) {
	tests := []struct {
		name   string
		header string
		tokens []string
		want   string
placeholder{
		{
			name:   "single token in middle",
			header: "oauth-2025-04-20,context-1m-2025-08-07,placeholder",
			tokens: []string{"context-1m-2025-08-07"placeholder,
			want:   "oauth-2025-04-20,placeholder",
	placeholder,
		{
			name:   "single token at start",
			header: "context-1m-2025-08-07,oauth-2025-04-20,placeholder",
			tokens: []string{"context-1m-2025-08-07"placeholder,
			want:   "oauth-2025-04-20,placeholder",
	placeholder,
		{
			name:   "single token at end",
			header: "oauth-2025-04-20,placeholder,context-1m-2025-08-07",
			tokens: []string{"context-1m-2025-08-07"placeholder,
			want:   "oauth-2025-04-20,placeholder",
	placeholder,
		{
			name:   "token not present",
			header: "oauth-2025-04-20,placeholder",
			tokens: []string{"context-1m-2025-08-07"placeholder,
			want:   "oauth-2025-04-20,placeholder",
	placeholder,
		{
			name:   "empty header",
			header: "",
			tokens: []string{"context-1m-2025-08-07"placeholder,
			want:   "",
	placeholder,
		{
			name:   "with spaces",
			header: "oauth-2025-04-20, context-1m-2025-08-07 , placeholder",
			tokens: []string{"context-1m-2025-08-07"placeholder,
			want:   "oauth-2025-04-20,placeholder",
	placeholder,
		{
			name:   "only token",
			header: "context-1m-2025-08-07",
			tokens: []string{"context-1m-2025-08-07"placeholder,
			want:   "",
	placeholder,
		{
			name:   "nil tokens",
			header: "oauth-2025-04-20,placeholder",
			tokens: nil,
			want:   "oauth-2025-04-20,placeholder",
	placeholder,
		{
			name:   "multiple tokens removed",
			header: "oauth-2025-04-20,context-1m-2025-08-07,placeholder,fast-mode-2026-02-01",
			tokens: []string{"context-1m-2025-08-07", "fast-mode-2026-02-01"placeholder,
			want:   "oauth-2025-04-20,placeholder",
	placeholder,
		{
			name:   "DroppedBetas removes fast-mode only",
			header: "oauth-2025-04-20,context-1m-2025-08-07,fast-mode-2026-02-01,placeholder",
			tokens: claude.DroppedBetas,
			want:   "oauth-2025-04-20,context-1m-2025-08-07,placeholder",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripBetaTokens(tt.header, tt.tokens)
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

func TestMergeAnthropicBetaDropping_DroppedBetas(t *testing.T) {
	required := []string{"oauth-2025-04-20", "placeholder"placeholder
	incoming := "context-1m-2025-08-07,fast-mode-2026-02-01,foo-beta,oauth-2025-04-20"
	drop := droppedBetaSet()

	got := mergeAnthropicBetaDropping(required, incoming, drop)
	require.Equal(t, "oauth-2025-04-20,placeholder,context-1m-2025-08-07,foo-beta", got)
	require.Contains(t, got, "context-1m-2025-08-07")
	require.NotContains(t, got, "fast-mode-2026-02-01")
placeholder

func TestDroppedBetaSet(t *testing.T) {
	// Base set contains DroppedBetas
	base := droppedBetaSet()
	require.NotContains(t, base, claude.BetaContext1M)
	require.Contains(t, base, claude.BetaFastMode)
	require.Len(t, base, len(claude.DroppedBetas))

	// With extra tokens
	extended := droppedBetaSet(claude.BetaClaudeCode)
	require.NotContains(t, extended, claude.BetaContext1M)
	require.Contains(t, extended, claude.BetaFastMode)
	require.Contains(t, extended, claude.BetaClaudeCode)
	require.Len(t, extended, len(claude.DroppedBetas)+1)
placeholder

func TestBuildBetaTokenSet(t *testing.T) {
	got := buildBetaTokenSet([]string{"foo", "", "bar", "foo"placeholder)
	require.Len(t, got, 2)
	require.Contains(t, got, "foo")
	require.Contains(t, got, "bar")
	require.NotContains(t, got, "")

	empty := buildBetaTokenSet(nil)
	require.Empty(t, empty)
placeholder

func TestContainsBetaToken(t *testing.T) {
	tests := []struct {
		name   string
		header string
		token  string
		want   bool
placeholder{
		{"present in middle", "oauth-2025-04-20,fast-mode-2026-02-01,placeholder", "fast-mode-2026-02-01", trueplaceholder,
		{"present at start", "fast-mode-2026-02-01,oauth-2025-04-20", "fast-mode-2026-02-01", trueplaceholder,
		{"present at end", "oauth-2025-04-20,fast-mode-2026-02-01", "fast-mode-2026-02-01", trueplaceholder,
		{"only token", "fast-mode-2026-02-01", "fast-mode-2026-02-01", trueplaceholder,
		{"not present", "oauth-2025-04-20,placeholder", "fast-mode-2026-02-01", falseplaceholder,
		{"with spaces", "oauth-2025-04-20, fast-mode-2026-02-01 , placeholder", "fast-mode-2026-02-01", trueplaceholder,
		{"empty header", "", "fast-mode-2026-02-01", falseplaceholder,
		{"empty token", "fast-mode-2026-02-01", "", falseplaceholder,
		{"partial match", "fast-mode-2026-02-01-extra", "fast-mode-2026-02-01", falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsBetaToken(tt.header, tt.token)
			require.Equal(t, tt.want, got)
	placeholder)
placeholder
placeholder

func TestStripBetaTokensWithSet_EmptyDropSet(t *testing.T) {
	header := "oauth-2025-04-20,placeholder"
	got := stripBetaTokensWithSet(header, map[string]struct{placeholder{placeholder)
	require.Equal(t, header, got)
placeholder

func TestIsCountTokensUnsupported404(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		want       bool
placeholder{
		{
			name:       "exact endpoint not found",
			statusCode: 404,
			body:       `{"error":{"message":"Not found: /v1/messages/count_tokens","type":"not_found_error"placeholderplaceholder`,
			want:       true,
	placeholder,
		{
			name:       "contains count_tokens and not found",
			statusCode: 404,
			body:       `{"error":{"message":"count_tokens route not found","type":"not_found_error"placeholderplaceholder`,
			want:       true,
	placeholder,
		{
			name:       "generic 404",
			statusCode: 404,
			body:       `{"error":{"message":"resource not found","type":"not_found_error"placeholderplaceholder`,
			want:       false,
	placeholder,
		{
			name:       "404 with empty error message",
			statusCode: 404,
			body:       `{"error":{"message":"","type":"not_found_error"placeholderplaceholder`,
			want:       false,
	placeholder,
		{
			name:       "non-404 status",
			statusCode: 400,
			body:       `{"error":{"message":"Not found: /v1/messages/count_tokens","type":"invalid_request_error"placeholderplaceholder`,
			want:       false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isCountTokensUnsupported404(tt.statusCode, []byte(tt.body))
			require.Equal(t, tt.want, got)
	placeholder)
placeholder
placeholder
