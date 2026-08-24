package service

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type pluginRoutingHTTPUpstream struct {
	doCalls        int
	doWithTLSCalls int
placeholder

func (u *pluginRoutingHTTPUpstream) Do(*http.Request, string, int64, int) (*http.Response, error) {
	u.doCalls++
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader("legacy")),
placeholder, nil
placeholder

func (u *pluginRoutingHTTPUpstream) DoWithTLS(
	request *http.Request,
	proxyURL string,
	accountID int64,
	accountConcurrency int,
	_ *tlsfingerprint.Profile,
) (*http.Response, error) {
	u.doWithTLSCalls++
	return u.Do(request, proxyURL, accountID, accountConcurrency)
placeholder

func TestPluginManagerRoutingDoesNotTouchAPIKeyOrOtherProviders(t *testing.T) {
	manager := &PluginManager{placeholder
	request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://example.com/v1/responses", nil)
placeholder

	accounts := []*Account{
		{ID: 1, Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder,
		{ID: 2, Platform: PlatformAnthropic, Type: AccountTypeOAuthplaceholder,
		{ID: 3, Platform: PlatformGemini, Type: AccountTypeOAuthplaceholder,
placeholder
	for _, account := range accounts {
		response, handled, routeErr := manager.RoundTripOpenAIOAuth(context.Background(), request, "", account)
		assert.Nil(t, response)
		assert.False(t, handled)
		assert.NoError(t, routeErr)
placeholder
placeholder

func TestPluginManagerRoutingKeepsOAuthOnLegacyPathWithoutEnabledBinding(t *testing.T) {
	manager := &PluginManager{placeholder
	request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://example.com/v1/responses", nil)
placeholder
	account := &Account{ID: 10, Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder

	response, handled, routeErr := manager.RoundTripOpenAIOAuth(context.Background(), request, "", account)

	assert.Nil(t, response)
	assert.False(t, handled)
	assert.NoError(t, routeErr)
placeholder

func TestPluginManagerRoutingSelectsOnlyEligibleOpenAIOAuthAccounts(t *testing.T) {
	manager := &PluginManager{placeholder
	manager.route.Store(&pluginRoute{pluginID: 1, rolloutPercent: 100, unavailable: "测试不可用"placeholder)

	assert.True(t, manager.ShouldRouteOpenAIOAuth(&Account{ID: 10, Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder))
	assert.False(t, manager.ShouldRouteOpenAIOAuth(&Account{ID: 10, Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder))
	assert.False(t, manager.ShouldRouteOpenAIOAuth(&Account{ID: 10, Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder))
	assert.False(t, manager.ShouldRouteOpenAIOAuth(nil))
placeholder

func TestOpenAIGatewayPluginRoutingPreservesAPIKeyAndFailsClosedForOAuth(t *testing.T) {
	manager := &PluginManager{placeholder
	manager.route.Store(&pluginRoute{pluginID: 1, rolloutPercent: 100, unavailable: "测试不可用"placeholder)
	upstream := &pluginRoutingHTTPUpstream{placeholder
	service := &OpenAIGatewayService{pluginManager: manager, httpUpstream: upstreamplaceholder

	apiKeyRequest, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://example.com/v1/responses", nil)
placeholder
	apiKeyResponse, err := service.doOpenAIUpstream(apiKeyRequest, "", &Account{
		ID: 1, Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Concurrency: 1,
placeholder)
placeholder
	require.NotNil(t, apiKeyResponse)
	_ = apiKeyResponse.Body.Close()
	assert.Equal(t, 1, upstream.doCalls)

	oauthRequest, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://example.com/v1/responses", nil)
placeholder
	oauthResponse, err := service.doOpenAIUpstream(oauthRequest, "", &Account{
		ID: 2, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Concurrency: 1,
placeholder)
placeholder
	assert.Nil(t, oauthResponse)
	assert.Contains(t, err.Error(), "插件不可用")
	assert.Equal(t, 1, upstream.doCalls)
placeholder

func TestStablePluginBucketIsDeterministicAndBounded(t *testing.T) {
	for id := int64(1); id <= 1000; id++ {
		first := stablePluginBucket(id)
		assert.Equal(t, first, stablePluginBucket(id))
		assert.Less(t, first, uint64(100))
placeholder
placeholder
