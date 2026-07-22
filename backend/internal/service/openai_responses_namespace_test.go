package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShouldFlattenOpenAIResponsesNamespaces(t *testing.T) {
	oauth := &Account{Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder
	apiKey := &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder
	grokOAuth := &Account{Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder

	tests := []struct {
		name               string
		account            *Account
		transport          OpenAIUpstreamTransport
		passthroughEnabled bool
		want               bool
placeholder{
		{name: "oauth_http", account: oauth, transport: OpenAIUpstreamTransportHTTPSSE, want: trueplaceholder,
		{name: "oauth_http_passthrough", account: oauth, transport: OpenAIUpstreamTransportHTTPSSE, passthroughEnabled: true, want: trueplaceholder,
		// WSv2 出口原样转发上游事件、不做回程还原，摊平会让客户端收到无法匹配的平名。
		{name: "oauth_wsv2", account: oauth, transport: OpenAIUpstreamTransportResponsesWebsocketV2, want: falseplaceholder,
		// 透传账号先于 WSv2 分支经 HTTP 转发返回，仍需摊平。
		{name: "oauth_wsv2_passthrough", account: oauth, transport: OpenAIUpstreamTransportResponsesWebsocketV2, passthroughEnabled: true, want: trueplaceholder,
		{name: "apikey_http", account: apiKey, transport: OpenAIUpstreamTransportHTTPSSE, want: falseplaceholder,
		{name: "grok_oauth_http", account: grokOAuth, transport: OpenAIUpstreamTransportHTTPSSE, want: falseplaceholder,
		{name: "nil_account", account: nil, transport: OpenAIUpstreamTransportHTTPSSE, want: falseplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, shouldFlattenOpenAIResponsesNamespaces(tt.account, tt.transport, tt.passthroughEnabled))
	placeholder)
placeholder
placeholder
