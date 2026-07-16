package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildOpenAIEndpointURLPreservesURLComponents(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		endpoint string
		want     string
placeholder{
		{name: "root", base: "https://upstream.example", endpoint: "/v1/models", want: "https://upstream.example/v1/models"placeholder,
		{name: "v1", base: "https://upstream.example/v1", endpoint: "/v1/responses", want: "https://upstream.example/v1/responses"placeholder,
		{name: "prefix", base: "https://upstream.example/openai", endpoint: "/v1/chat/completions", want: "https://upstream.example/openai/v1/chat/completions"placeholder,
		{name: "version", base: "https://upstream.example/openai/v2", endpoint: "/v1/embeddings", want: "https://upstream.example/openai/v2/embeddings"placeholder,
		{name: "query", base: "https://upstream.example/v1?redirect=/", endpoint: "/v1/sub2api/billing", want: "https://upstream.example/v1/sub2api/billing?redirect=/"placeholder,
		{name: "fragment is removed", base: "https://upstream.example/v1#stale", endpoint: "/v1/alpha/search", want: "https://upstream.example/v1/alpha/search"placeholder,
		{name: "ipv6", base: "http://[2001:db8::1]:8080/v1?tenant=a#stale", endpoint: "/v1/responses/input_tokens", want: "http://[2001:db8::1]:8080/v1/responses/input_tokens?tenant=a"placeholder,
		{name: "already complete", base: "https://upstream.example/v1/images/generations?tenant=a", endpoint: "/v1/images/generations", want: "https://upstream.example/v1/images/generations?tenant=a"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, buildOpenAIEndpointURL(tt.base, tt.endpoint))
	placeholder)
placeholder
placeholder
