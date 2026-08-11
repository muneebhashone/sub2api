package service

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSetOpenAICodexRoutingHintCanonicalizesOfficialServiceTiers(t *testing.T) {
	oauthAccount := &Account{Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder
	tests := []struct {
		name        string
		model       string
		serviceTier string
		want        string
placeholder{
		{name: "fast alias", model: "gpt-5.6", serviceTier: " fast ", want: "model=gpt-5.6;tier=priority"placeholder,
		{name: "priority", model: "gpt-5.6", serviceTier: "priority", want: "model=gpt-5.6;tier=priority"placeholder,
		{name: "flex", model: "gpt-5.6", serviceTier: "flex", want: "model=gpt-5.6;tier=flex"placeholder,
		{name: "explicit default sentinel", model: "gpt-5.6", serviceTier: "default", want: "model=gpt-5.6"placeholder,
		{name: "omitted tier", model: "gpt-5.6", want: "model=gpt-5.6"placeholder,
		{name: "auto is not expanded without catalog support", model: "gpt-5.6", serviceTier: "auto", want: "model=gpt-5.6"placeholder,
		{name: "scale is not expanded without catalog support", model: "gpt-5.6", serviceTier: "scale", want: "model=gpt-5.6"placeholder,
		{name: "unknown tier does not expand protocol", model: "gpt-5.6", serviceTier: "turbo", want: "model=gpt-5.6"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := make(http.Header)
			setOpenAICodexRoutingHint(headers, oauthAccount, tt.model, tt.serviceTier)
			require.Equal(t, tt.want, headers.Get(openAICodexRoutingHintHeader))
	placeholder)
placeholder

	t.Run("invalid header value is omitted", func(t *testing.T) {
		headers := make(http.Header)
		setOpenAICodexRoutingHint(headers, oauthAccount, "gpt-5.6\ninvalid", "priority")
		require.Empty(t, headers.Get(openAICodexRoutingHintHeader))
placeholder)

	for _, model := range []string{"gpt-5.6;evil", "gpt=5.6"placeholder {
		t.Run("delimiter in model is omitted: "+model, func(t *testing.T) {
			headers := make(http.Header)
			setOpenAICodexRoutingHint(headers, oauthAccount, model, "priority")
			require.Empty(t, headers.Get(openAICodexRoutingHintHeader))
	placeholder)
placeholder

	t.Run("api key strips gateway-owned hint in every key casing", func(t *testing.T) {
		headers := make(http.Header)
		headers[openAICodexRoutingHintHeader] = []string{"lowercase-spoof"placeholder
		headers["X-Codex-Routing-Hint"] = []string{"canonical-spoof"placeholder
		setOpenAICodexRoutingHint(headers, &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder, "gpt-5.6", "priority")
		for key := range headers {
			require.False(t, strings.EqualFold(key, openAICodexRoutingHintHeader))
	placeholder
placeholder)

	t.Run("oauth replaces spoofed lowercase hint", func(t *testing.T) {
		headers := make(http.Header)
		headers[openAICodexRoutingHintHeader] = []string{"model=spoof;tier=flex"placeholder
		setOpenAICodexRoutingHint(headers, oauthAccount, "gpt-5.6", "priority")
		require.Equal(t, "model=gpt-5.6;tier=priority", headers.Get(openAICodexRoutingHintHeader))
		require.Len(t, headers, 1)
placeholder)
placeholder

func TestOpenAIOAuthHTTPBuildersSendRoutingHintFromFinalBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	oauthAccount := &Account{
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"chatgpt_account_id": "test-account",
	placeholder,
placeholder
	svc := &OpenAIGatewayService{placeholder

	tests := []struct {
		name string
		body []byte
		want string
placeholder{
		{name: "fast", body: []byte(`{"model":"gpt-5.6-codex","service_tier":"fast"placeholder`), want: "model=gpt-5.6-codex;tier=priority"placeholder,
		{name: "flex", body: []byte(`{"model":"gpt-5.6-codex","service_tier":"flex"placeholder`), want: "model=gpt-5.6-codex;tier=flex"placeholder,
		{name: "default", body: []byte(`{"model":"gpt-5.6-codex","service_tier":"default"placeholder`), want: "model=gpt-5.6-codex"placeholder,
		{name: "omitted", body: []byte(`{"model":"gpt-5.6-codex"placeholder`), want: "model=gpt-5.6-codex"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, passthrough := range []bool{false, trueplaceholder {
				mode := "ordinary"
				if passthrough {
					mode = "passthrough"
			placeholder
				t.Run(mode, func(t *testing.T) {
					recorder := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(recorder)
					c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(tt.body))

					var req *http.Request
					var err error
					if passthrough {
						req, err = svc.buildUpstreamRequestOpenAIPassthrough(context.Background(), c, oauthAccount, tt.body, "test-token")
				placeholder else {
						req, err = svc.buildUpstreamRequest(context.Background(), c, oauthAccount, tt.body, "test-token", false, "", true)
				placeholder
				placeholder
					require.Equal(t, tt.want, req.Header.Get(openAICodexRoutingHintHeader))
			placeholder)
		placeholder
	placeholder)
placeholder
placeholder

func TestOpenAIHTTPPassthroughStripsOnlyOAuthLegacyResponsesBeta(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &OpenAIGatewayService{cfg: &config.Config{
		Security: config.SecurityConfig{
			URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholder,
	placeholder,
placeholderplaceholder
	body := []byte(`{"model":"gpt-5.6-codex","service_tier":"priority"placeholder`)

	build := func(t *testing.T, account *Account, betaValues []string, rawLowercaseKey bool) http.Header {
	placeholder
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))
		if rawLowercaseKey {
			c.Request.Header["openai-beta"] = append([]string(nil), betaValues...)
	placeholder else {
			for _, value := range betaValues {
				c.Request.Header.Add("OpenAI-Beta", value)
		placeholder
	placeholder

		req, err := svc.buildUpstreamRequestOpenAIPassthrough(context.Background(), c, account, body, "test-token")
	placeholder
		return req.Header
placeholder

	oauth := &Account{
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"chatgpt_account_id": "test-account",
	placeholder,
placeholder
	apiKey := &Account{
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
placeholder
			"api_key": "test-api-key",
	placeholder,
placeholder

	t.Run("oauth legacy only is removed including raw lowercase key", func(t *testing.T) {
		headers := build(t, oauth, []string{"responses=experimental"placeholder, true)
		require.Empty(t, headers.Values("OpenAI-Beta"))
placeholder)

	t.Run("oauth mixed beta preserves independent tokens", func(t *testing.T) {
		headers := build(t, oauth, []string{
			"responses=experimental, future_feature=v1",
			"another_feature=v2, RESPONSES=EXPERIMENTAL",
	placeholder, false)
		require.Equal(t, []string{"future_feature=v1", "another_feature=v2"placeholder, headers.Values("OpenAI-Beta"))
placeholder)

	t.Run("api key explicit beta remains caller controlled", func(t *testing.T) {
		headers := build(t, apiKey, []string{"responses=experimental, future_feature=v1"placeholder, false)
		require.Equal(t, []string{"responses=experimental, future_feature=v1"placeholder, headers.Values("OpenAI-Beta"))
placeholder)
placeholder

func TestBuildOpenAIWSHeadersSendsOAuthRoutingHintOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/responses", nil)
	svc := &OpenAIGatewayService{placeholder
	decision := OpenAIWSProtocolDecision{Transport: OpenAIUpstreamTransportResponsesWebsocketV2placeholder

	build := func(t *testing.T, account *Account, tier string) http.Header {
		headers, _, err := svc.buildOpenAIWSHeaders(
			context.Background(),
			c,
			account,
			"test-token",
			decision,
			true,
			"",
			"",
			"",
			"gpt-5.6-codex",
			tier,
		)
	placeholder
		return headers
placeholder

	oauthAccount := &Account{
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"chatgpt_account_id": "test-account",
	placeholder,
placeholder
	require.Equal(t, "model=gpt-5.6-codex;tier=priority", build(t, oauthAccount, "fast").Get(openAICodexRoutingHintHeader))
	require.Equal(t, "model=gpt-5.6-codex", build(t, oauthAccount, "default").Get(openAICodexRoutingHintHeader))
	require.Empty(t, build(t, &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder, "priority").Get(openAICodexRoutingHintHeader))
placeholder

func TestOpenAIRoutingDiagnosticsUseFinalDerivedValuesOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logSink, restore := captureStructuredLog(t)
	defer restore()

	account := &Account{
		ID:       917,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"chatgpt_account_id": "chatgpt-account",
	placeholder,
placeholder
	body := []byte(`{"model":"gpt-5.6-codex","service_tier":"fast"placeholder`)
	svc := &OpenAIGatewayService{placeholder

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))
	c.Request.Header.Set("Authorization", "Bearer caller-secret")
	c.Request.Header.Set(openAICodexRoutingHintHeader, "model=caller-secret")
	_, err := svc.buildUpstreamRequest(context.Background(), c, account, body, "oauth-secret", false, "", true)
placeholder

	decision := OpenAIWSProtocolDecision{Transport: OpenAIUpstreamTransportResponsesWebsocketV2placeholder
	_, _, err = svc.buildOpenAIWSHeaders(
		context.Background(), c, account, "oauth-secret", decision, true,
		"", "", "", "gpt-5.6-codex", "fast",
	)
placeholder

	require.True(t, logSink.ContainsMessageAtLevel("openai routing decision", "debug"))
	require.True(t, logSink.ContainsFieldValue("account_id", "917"))
	require.True(t, logSink.ContainsFieldValue("final_model", "gpt-5.6-codex"))
	require.True(t, logSink.ContainsFieldValue("final_service_tier", "priority"))
	require.True(t, logSink.ContainsFieldValue("routing_hint_generated", "true"))
	require.True(t, logSink.ContainsFieldValue("transport", "http"))
	require.True(t, logSink.ContainsFieldValue("transport", string(OpenAIUpstreamTransportResponsesWebsocketV2)))
	require.True(t, logSink.ContainsFieldValue("ws_affinity_decision", "not_applicable"))
	require.True(t, logSink.ContainsFieldValue("ws_affinity_decision", "soft_routing_hint"))
	require.False(t, logSink.ContainsFieldValue("authorization", "caller-secret"))
	require.False(t, logSink.ContainsFieldValue("credentials", "oauth-secret"))
	require.False(t, logSink.ContainsFieldValue("routing_hint", "caller-secret"))
placeholder

func TestOpenAIWSConnPoolPreferredContinuationIgnoresRoutingHintChanges(t *testing.T) {
	cfg := &config.Config{placeholder
	cfg.Gateway.OpenAIWS.MaxConnsPerAccount = 2
	cfg.Gateway.OpenAIWS.MinIdlePerAccount = 0
	cfg.Gateway.OpenAIWS.MaxIdlePerAccount = 2

	pool := newOpenAIWSConnPool(cfg)
	dialer := &openAIWSCountingDialer{placeholder
	pool.setClientDialerForTest(dialer)
	account := &Account{ID: 913, Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder

	acquire := func(t *testing.T, hint, preferred string, forcePreferred bool) *openAIWSConnLease {
	placeholder
		headers := make(http.Header)
		if hint != "" {
			headers.Set(openAICodexRoutingHintHeader, hint)
	placeholder
		lease, err := pool.Acquire(context.Background(), openAIWSAcquireRequest{
			Account:            account,
			WSURL:              "wss://example.com/v1/responses",
			Headers:            headers,
			PreferredConnID:    preferred,
			ForcePreferredConn: forcePreferred,
	placeholder)
	placeholder
		require.NotNil(t, lease)
		return lease
placeholder

	standard := acquire(t, "model=gpt-5.6-codex", "", false)
	connID := standard.ConnID()
	standard.Release()

	priority := acquire(t, "model=gpt-5.6-codex;tier=priority", connID, true)
	require.True(t, priority.Reused())
	require.Equal(t, connID, priority.ConnID())
	priority.Release()

	standardAgain := acquire(t, "model=gpt-5.6-codex", connID, true)
	require.True(t, standardAgain.Reused())
	require.Equal(t, connID, standardAgain.ConnID())
	standardAgain.Release()

	require.Equal(t, 1, dialer.DialCount(), "routing hint is dial-time advisory, not continuation compatibility")
placeholder

func TestOpenAIWSConnPoolUsesRoutingHintAsSoftDialAffinity(t *testing.T) {
	cfg := &config.Config{placeholder
	cfg.Gateway.OpenAIWS.MaxConnsPerAccount = 4
	cfg.Gateway.OpenAIWS.MinIdlePerAccount = 0
	cfg.Gateway.OpenAIWS.MaxIdlePerAccount = 4

	pool := newOpenAIWSConnPool(cfg)
	dialer := &openAIWSCountingDialer{placeholder
	pool.setClientDialerForTest(dialer)
	account := &Account{ID: 913, Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder

	acquire := func(t *testing.T, hint string) *openAIWSConnLease {
		headers := make(http.Header)
		headers.Set(openAICodexRoutingHintHeader, hint)
		lease, err := pool.Acquire(context.Background(), openAIWSAcquireRequest{
			Account: account,
			WSURL:   "wss://example.com/v1/responses",
			Headers: headers,
	placeholder)
	placeholder
		require.NotNil(t, lease)
		return lease
placeholder

	priority := acquire(t, "model=gpt-5.6-codex;tier=priority")
	priorityConnID := priority.ConnID()
	priority.Release()

	priorityAgain := acquire(t, "model=gpt-5.6-codex;tier=priority")
	require.True(t, priorityAgain.Reused())
	require.Equal(t, priorityConnID, priorityAgain.ConnID())
	priorityAgain.Release()

	flex := acquire(t, "model=gpt-5.6-codex;tier=flex")
	require.False(t, flex.Reused())
	require.NotEqual(t, priorityConnID, flex.ConnID())
	flex.Release()

	otherModel := acquire(t, "model=gpt-5.5-codex;tier=priority")
	require.False(t, otherModel.Reused())
	require.NotEqual(t, priorityConnID, otherModel.ConnID())
	otherModel.Release()

	defaultTier := acquire(t, "model=gpt-5.6-codex")
	require.False(t, defaultTier.Reused())
	require.NotEqual(t, priorityConnID, defaultTier.ConnID())
	defaultConnID := defaultTier.ConnID()
	defaultTier.Release()

	defaultAgain := acquire(t, "model=gpt-5.6-codex")
	require.True(t, defaultAgain.Reused())
	require.Equal(t, defaultConnID, defaultAgain.ConnID())
	defaultAgain.Release()

	require.Equal(t, 4, dialer.DialCount())
placeholder
