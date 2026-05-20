package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSafeUpstreamURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
placeholder{
		{"strips query", "https://api.anthropic.com/v1/messages?beta=true", "https://api.anthropic.com/v1/messages"placeholder,
		{"strips fragment", "https://api.openai.com/v1/responses#frag", "https://api.openai.com/v1/responses"placeholder,
		{"strips both", "https://host/path?token=secret#x", "https://host/path"placeholder,
		{"no query or fragment", "https://host/path", "https://host/path"placeholder,
		{"empty string", "", ""placeholder,
		{"whitespace only", "  ", ""placeholder,
		{"query before fragment", "https://h/p?a=1#f", "https://h/p"placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, safeUpstreamURL(tt.input))
	placeholder)
placeholder
placeholder
