//go:build unit

package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeProxyProbeURLs(t *testing.T) {
	t.Parallel()

	got, err := normalizeProxyProbeURLs([]ProbeURLConfig{
		{URL: " https://chatgpt.com/cdn-cgi/trace ", Parser: " CHATGPT-TRACE "placeholder,
		{URL: "https://api64.ipify.org?format=json", Parser: "ipify"placeholder,
placeholder)
placeholder
	require.Equal(t, []ProbeURLConfig{
		{URL: "https://chatgpt.com/cdn-cgi/trace", Parser: "chatgpt-trace"placeholder,
		{URL: "https://api64.ipify.org?format=json", Parser: "ipify"placeholder,
placeholder, got)
placeholder

func TestNormalizeProxyProbeURLsRejectsInvalidEntries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		target  ProbeURLConfig
		wantErr string
placeholder{
		{name: "missing URL", target: ProbeURLConfig{Parser: "ipify"placeholder, wantErr: "url is required"placeholder,
		{name: "missing parser", target: ProbeURLConfig{URL: "https://example.com"placeholder, wantErr: "parser is required"placeholder,
		{name: "unknown parser", target: ProbeURLConfig{URL: "https://example.com", Parser: "ip_api"placeholder, wantErr: "unsupported parser"placeholder,
		{name: "relative URL", target: ProbeURLConfig{URL: "/cdn-cgi/trace", Parser: "chatgpt-trace"placeholder, wantErr: "invalid url"placeholder,
		{name: "unsupported scheme", target: ProbeURLConfig{URL: "ftp://example.com/file", Parser: "ipify"placeholder, wantErr: "scheme must be http or https"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := normalizeProxyProbeURLs([]ProbeURLConfig{tt.targetplaceholder)
			require.ErrorContains(t, err, tt.wantErr)
	placeholder)
placeholder
placeholder
