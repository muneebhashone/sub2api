//go:build unit

package service

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func newGrokCacheTestContext(apiKeyID int64) *gin.Context {
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	if apiKeyID > 0 {
		c.Set("api_key", &APIKey{ID: apiKeyID, Group: &Group{Platform: PlatformGrokplaceholderplaceholder)
placeholder
	return c
placeholder

func TestGrokPreviousResponseSessionSeed(t *testing.T) {
	require.Equal(t, "grok-prev-resp:resp_abc123", grokPreviousResponseSessionSeed([]byte(`{"previous_response_id":"resp_abc123"placeholder`)))
	require.Empty(t, grokPreviousResponseSessionSeed([]byte(`{"previous_response_id":"msg_abc123"placeholder`)))
	require.Empty(t, grokPreviousResponseSessionSeed([]byte(`{"previous_response_id":""placeholder`)))
	require.Empty(t, grokPreviousResponseSessionSeed([]byte(`{placeholder`)))
placeholder

func TestResolveGrokCacheIdentityUsesPreviousResponseIDWhenNoOtherSeed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := newGrokCacheTestContext(301)
	// No prompt_cache_key / headers / reusable prefix — only previous_response_id.
	body := []byte(`{"model":"grok","input":[{"role":"user","content":"follow up"placeholder],"previous_response_id":"resp_chain_001"placeholder`)
	got := resolveGrokCacheIdentity(c, body, "", "grok-4.5")
	require.NotEmpty(t, got)

	// Same previous_response_id → same identity (model already in isolated seed).
	again := resolveGrokCacheIdentity(c, body, "", "grok-4.5")
	require.Equal(t, got, again)

	// Different model → different identity (model scope).
	otherModel := resolveGrokCacheIdentity(c, body, "", "grok-4.3")
	require.NotEqual(t, got, otherModel)

	// prompt_cache_key still wins over previous_response_id.
	withCache := []byte(`{"model":"grok","prompt_cache_key":"client-session","previous_response_id":"resp_chain_001","input":[{"role":"user","content":"x"placeholder]placeholder`)
	cacheID := resolveGrokCacheIdentity(c, withCache, "", "grok-4.5")
	require.NotEmpty(t, cacheID)
	require.NotEqual(t, got, cacheID)
placeholder

func TestResolveGrokCacheIdentityStableAcrossAppendOnlyTurns(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := newGrokCacheTestContext(101)
	round1 := []byte(`{"model":"grok","instructions":"be concise","tools":[{"type":"function","name":"lookup","parameters":{"type":"object"placeholderplaceholder],"input":[{"role":"user","content":"first question"placeholder]placeholder`)
	round2 := []byte(`{"model":"grok","instructions":"be concise","tools":[{"type":"function","name":"lookup","parameters":{"type":"object"placeholderplaceholder],"input":[{"role":"user","content":"first question"placeholder,{"role":"assistant","content":"first answer"placeholder,{"role":"user","content":"second question"placeholder]placeholder`)

	first := resolveGrokCacheIdentity(c, round1, "", "grok-4.5")
	second := resolveGrokCacheIdentity(c, round2, "", "grok-4.5")

	require.NotEmpty(t, first)
	require.Len(t, first, 36)
	require.Equal(t, first, second)
placeholder

func TestResolveGrokCacheIdentityStableAcrossIndependentPromptsWithSamePrefix(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := newGrokCacheTestContext(102)
	firstBody := []byte(`{"model":"grok","instructions":"be concise","tools":[{"type":"function","name":"lookup"placeholder],"input":[{"role":"user","content":"Question A"placeholder]placeholder`)
	secondBody := []byte(`{"model":"grok","instructions":"be concise","tools":[{"type":"function","name":"lookup"placeholder],"input":[{"role":"user","content":"Question B"placeholder]placeholder`)

	first := resolveGrokCacheIdentity(c, firstBody, "", "grok-4.5")
	second := resolveGrokCacheIdentity(c, secondBody, "", "grok-4.5")

	require.NotEmpty(t, first)
	require.Equal(t, first, second)
placeholder

func TestResolveGrokCacheIdentityStablePrefixIsolation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	baseBody := []byte(`{"model":"grok","instructions":"be concise","tools":[{"type":"function","name":"lookup"placeholder],"input":[{"role":"system","content":"System A"placeholder,{"role":"user","content":"Question A"placeholder]placeholder`)
	differentInstructions := []byte(`{"model":"grok","instructions":"be detailed","tools":[{"type":"function","name":"lookup"placeholder],"input":[{"role":"system","content":"System A"placeholder,{"role":"user","content":"Question B"placeholder]placeholder`)
	differentSystem := []byte(`{"model":"grok","instructions":"be concise","tools":[{"type":"function","name":"lookup"placeholder],"input":[{"role":"system","content":"System B"placeholder,{"role":"user","content":"Question B"placeholder]placeholder`)
	differentTools := []byte(`{"model":"grok","instructions":"be concise","tools":[{"type":"function","name":"search"placeholder],"input":[{"role":"system","content":"System A"placeholder,{"role":"user","content":"Question B"placeholder]placeholder`)

	base := resolveGrokCacheIdentity(newGrokCacheTestContext(103), baseBody, "", "grok-4.5")
	require.NotEqual(t, base, resolveGrokCacheIdentity(newGrokCacheTestContext(104), baseBody, "", "grok-4.5"))
	require.NotEqual(t, base, resolveGrokCacheIdentity(newGrokCacheTestContext(103), baseBody, "", "grok-4.3"))
	require.NotEqual(t, base, resolveGrokCacheIdentity(newGrokCacheTestContext(103), differentInstructions, "", "grok-4.5"))
	require.NotEqual(t, base, resolveGrokCacheIdentity(newGrokCacheTestContext(103), differentSystem, "", "grok-4.5"))
	require.NotEqual(t, base, resolveGrokCacheIdentity(newGrokCacheTestContext(103), differentTools, "", "grok-4.5"))
placeholder

func TestResolveGrokCacheIdentityFallsBackWhenStablePrefixIsEmpty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := newGrokCacheTestContext(105)
	firstBody := []byte(`{"model":"grok","tools":[],"input":"Question A"placeholder`)
	secondBody := []byte(`{"model":"grok","tools":[],"input":"Question B"placeholder`)

	first := resolveGrokCacheIdentity(c, firstBody, "", "grok-4.5")
	second := resolveGrokCacheIdentity(c, secondBody, "", "grok-4.5")

	require.NotEmpty(t, first)
	require.NotEmpty(t, second)
	require.NotEqual(t, first, second)
placeholder

func TestResolveGrokCacheIdentitySkipsUnanchoredFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := newGrokCacheTestContext(106)
	tests := [][]byte{
		[]byte(`{"model":"grok"placeholder`),
		[]byte(`{"model":"grok","messages":[{"role":"assistant","content":"answer"placeholder]placeholder`),
		[]byte(`{"model":"grok","messages":[{"role":"user","content":""placeholder]placeholder`),
		[]byte(`{"model":"grok","input":"  "placeholder`),
placeholder

	for _, body := range tests {
		require.Empty(t, resolveGrokCacheIdentity(c, body, "", "grok-4.5"))
placeholder
placeholder

func TestResolveGrokCacheIdentityIsolatesAPIKeyAndMappedModel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"grok","input":"same prompt"placeholder`)

	base := resolveGrokCacheIdentity(newGrokCacheTestContext(201), body, "", "grok-4.5")
	otherTenant := resolveGrokCacheIdentity(newGrokCacheTestContext(202), body, "", "grok-4.5")
	otherModel := resolveGrokCacheIdentity(newGrokCacheTestContext(201), body, "", "grok-4.3")

	require.NotEmpty(t, base)
	require.NotEqual(t, base, otherTenant)
	require.NotEqual(t, base, otherModel)
placeholder

func TestResolveGrokCacheIdentityUsesAndIsolatesNativeConversationHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := newGrokCacheTestContext(301)
	c.Request.Header.Set(grokConversationIDHeader, "raw-native-conversation")
	body1 := []byte(`{"model":"grok","input":"one"placeholder`)
	body2 := []byte(`{"model":"grok","input":"different body that must not replace the explicit session"placeholder`)

	first := resolveGrokCacheIdentity(c, body1, "body-cache-key", "grok-4.5")
	second := resolveGrokCacheIdentity(c, body2, "another-body-cache-key", "grok-4.5")

	require.Equal(t, "raw-native-conversation", (&OpenAIGatewayService{placeholder).ExtractSessionID(c, body1))
	require.Equal(t, first, second)
	require.NotEqual(t, "raw-native-conversation", first)
	require.NotContains(t, first, "raw-native-conversation")
placeholder

func TestResolveGrokCacheIdentityExplicitHeaderPriority(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"grok","prompt_cache_key":"body-key","input":"hi"placeholder`)
	c := newGrokCacheTestContext(401)
	c.Request.Header.Set(grokConversationIDHeader, "grok-key")
	c.Request.Header.Set("conversation_id", "conversation-key")
	c.Request.Header.Set("session_id", "session-key")

	got := resolveGrokCacheIdentity(c, body, "explicit-argument", "grok-4.5")
	onlySession := newGrokCacheTestContext(401)
	onlySession.Request.Header.Set("session_id", "session-key")
	want := resolveGrokCacheIdentity(onlySession, []byte(`{"model":"grok","input":"unrelated"placeholder`), "", "grok-4.5")

	require.Equal(t, want, got)
placeholder

func TestResolveGrokCacheIdentityIDEHeaderPriority(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"grok","prompt_cache_key":"body-key","input":"hi"placeholder`)
	headers := []struct {
		name  string
		value string
placeholder{
		{name: openCodeSessionAffinityHeader, value: "opencode-affinity"placeholder,
		{name: openCodeSessionIDHeader, value: "opencode-session-id"placeholder,
		{name: openCodeNativeSessionHeader, value: "opencode-native-session"placeholder,
		{name: codeBuddyConversationHeader, value: "codebuddy-conversation"placeholder,
		{name: grokConversationIDHeader, value: "grok-conversation"placeholder,
placeholder

	c := newGrokCacheTestContext(402)
	for _, header := range headers {
		c.Request.Header.Set(header.name, header.value)
placeholder
	for _, header := range headers {
		got := resolveGrokCacheIdentity(c, body, "explicit-argument", "grok-4.5")
		onlyCurrent := newGrokCacheTestContext(402)
		onlyCurrent.Request.Header.Set(header.name, header.value)
		want := resolveGrokCacheIdentity(onlyCurrent, []byte(`{"model":"grok","input":"unrelated"placeholder`), "", "grok-4.5")
		require.Equal(t, want, got, header.name)
		c.Request.Header.Del(header.name)
placeholder
placeholder

func TestExplicitGrokCacheSeedPriority(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := newGrokCacheTestContext(403)
	headers := []struct {
		name  string
		value string
placeholder{
		{name: claudeCodeSessionHeader, value: "claude-session"placeholder,
		{name: "session_id", value: "generic-session"placeholder,
		{name: "conversation_id", value: "generic-conversation"placeholder,
		{name: openCodeSessionAffinityHeader, value: "opencode-affinity"placeholder,
		{name: openCodeSessionIDHeader, value: "opencode-session-id"placeholder,
		{name: openCodeNativeSessionHeader, value: "opencode-native-session"placeholder,
		{name: codeBuddyConversationHeader, value: "codebuddy-conversation"placeholder,
		{name: grokConversationIDHeader, value: "grok-conversation"placeholder,
placeholder
	for _, header := range headers {
		c.Request.Header.Set(header.name, header.value)
placeholder

	body := []byte(`{"model":"grok","prompt_cache_key":"body-key","input":"hi"placeholder`)
	for _, header := range headers {
		require.Equal(t, header.value, explicitGrokCacheSeed(c, body, "explicit-argument"), header.name)
		c.Request.Header.Del(header.name)
placeholder
	require.Equal(t, "body-key", explicitGrokCacheSeed(c, body, "explicit-argument"))
	require.Equal(t, "explicit-argument", explicitGrokCacheSeed(c, []byte(`{"model":"grok"placeholder`), "explicit-argument"))
placeholder

func TestResolveGrokCacheIdentityIDEHeadersAreStableIsolatedAndOpaque(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name   string
		header string
placeholder{
		{name: "OpenCode affinity", header: openCodeSessionAffinityHeaderplaceholder,
		{name: "OpenCode session ID", header: openCodeSessionIDHeaderplaceholder,
		{name: "OpenCode native session", header: openCodeNativeSessionHeaderplaceholder,
		{name: "CodeBuddy conversation", header: codeBuddyConversationHeaderplaceholder,
placeholder

	for index, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawSession := "raw-ide-session-" + tt.name
			apiKeyID := int64(800 + index)
			c := newGrokCacheTestContext(apiKeyID)
			c.Request.Header.Set(tt.header, rawSession)
			firstBody := []byte(`{"model":"grok","prompt_cache_key":"turn-one-body-key","input":"first turn"placeholder`)
			secondBody := []byte(`{"model":"grok","prompt_cache_key":"turn-two-body-key","input":"different second turn"placeholder`)

			first := resolveGrokCacheIdentity(c, firstBody, "first-explicit-key", "grok-4.5")
			second := resolveGrokCacheIdentity(c, secondBody, "second-explicit-key", "grok-4.5")
			require.NotEmpty(t, first)
			require.Equal(t, first, second)
			require.NotEqual(t, rawSession, first)
			require.NotContains(t, first, rawSession)

			otherTenant := newGrokCacheTestContext(apiKeyID + 100)
			otherTenant.Request.Header.Set(tt.header, rawSession)
			require.NotEqual(t, first, resolveGrokCacheIdentity(otherTenant, firstBody, "", "grok-4.5"))
			require.NotEqual(t, first, resolveGrokCacheIdentity(c, firstBody, "", "grok-4.3"))
	placeholder)
placeholder
placeholder

func TestOpenCodeResponsesHeaderAndBodyCacheSignalsConverge(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const rawSession = "opencode-session-42"
	c := newGrokCacheTestContext(901)
	c.Request.Header.Set(openCodeSessionAffinityHeader, rawSession)
	c.Request.Header.Set(openCodeSessionIDHeader, rawSession)
	firstBody := []byte(`{"model":"grok","prompt_cache_key":"opencode-session-42","input":"first turn"placeholder`)
	secondBody := []byte(`{"model":"grok","prompt_cache_key":"opencode-session-42","input":"different second turn"placeholder`)

	first := resolveGrokCacheIdentity(c, firstBody, "", "grok-4.5")
	second := resolveGrokCacheIdentity(c, secondBody, "", "grok-4.5")
	bodyOnly := resolveGrokCacheIdentity(newGrokCacheTestContext(901), secondBody, "", "grok-4.5")
	require.NotEmpty(t, first)
	require.Equal(t, first, second)
	require.Equal(t, first, bodyOnly)

	patched, err := applyGrokResponsesCacheIdentity(secondBody, secondBody, second, false)
placeholder
	require.Equal(t, second, gjson.GetBytes(patched, "prompt_cache_key").String())
	headers := make(http.Header)
	applyGrokCacheHeaders(headers, second)
	require.Equal(t, second, headers.Get(grokConversationIDHeader))
	require.NotContains(t, string(patched), rawSession)
placeholder

func TestResolveGrokCacheIdentityPrefersClaudeCodeSession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := newGrokCacheTestContext(701)
	c.Request.Header.Set(claudeCodeSessionHeader, "cc-session-abc")
	c.Request.Header.Set("session_id", "session-key")
	body1 := []byte(`{"model":"grok","input":"turn-1"placeholder`)
	body2 := []byte(`{"model":"grok","input":"turn-2-different-user-text"placeholder`)

	first := resolveGrokCacheIdentity(c, body1, "", "grok-4.5")
	second := resolveGrokCacheIdentity(c, body2, "unrelated-explicit", "grok-4.5")
	require.NotEmpty(t, first)
	require.Equal(t, first, second, "same Claude Code session must keep stable Grok cache identity across turns")

	// metadata.user_id JSON form used by Claude Code clients
	metaBody := []byte(`{"model":"grok","metadata":{"user_id":"{\"session_id\":\"meta-session-xyz\"placeholder"placeholder,"input":"hi"placeholder`)
	metaOnly := newGrokCacheTestContext(702)
	metaID := resolveGrokCacheIdentity(metaOnly, metaBody, "", "grok-4.5")
	require.NotEmpty(t, metaID)
	metaBody2 := []byte(`{"model":"grok","metadata":{"user_id":"{\"session_id\":\"meta-session-xyz\"placeholder"placeholder,"input":"later turn"placeholder`)
	require.Equal(t, metaID, resolveGrokCacheIdentity(metaOnly, metaBody2, "", "grok-4.5"))
placeholder

func TestResolveGrokCacheIdentityFailsClosedWithoutAPIKeyContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := newGrokCacheTestContext(0)
	c.Request.Header.Set(grokConversationIDHeader, "native-session")

	require.Empty(t, resolveGrokCacheIdentity(c, []byte(`{"model":"grok","input":"hi"placeholder`), "", "grok-4.5"))
	require.Empty(t, resolveGrokCacheIdentity(nil, []byte(`{"model":"grok","prompt_cache_key":"key"placeholder`), "key", "grok-4.5"))
placeholder

func TestGrokConversationHeaderIsScopedToGrokRequestScheduling(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"grok","prompt_cache_key":"body-session","input":"hi"placeholder`)

	grokContext := newGrokCacheTestContext(601)
	grokContext.Request.Header.Set(grokConversationIDHeader, "native-grok-session")
	require.Equal(t, "native-grok-session", (&OpenAIGatewayService{placeholder).ExtractSessionID(grokContext, body))

	openAIContext := newGrokCacheTestContext(601)
	openAIContext.Set("api_key", &APIKey{ID: 601, Group: &Group{Platform: PlatformOpenAIplaceholderplaceholder)
	openAIContext.Request.Header.Set(grokConversationIDHeader, "must-be-ignored")
	require.Equal(t, "body-session", (&OpenAIGatewayService{placeholder).ExtractSessionID(openAIContext, body))

	withoutGrokHeader := newGrokCacheTestContext(601)
	withoutGrokHeader.Set("api_key", &APIKey{ID: 601, Group: &Group{Platform: PlatformOpenAIplaceholderplaceholder)
	require.Equal(t,
		(&OpenAIGatewayService{placeholder).GenerateSessionHash(withoutGrokHeader, body),
		(&OpenAIGatewayService{placeholder).GenerateSessionHash(openAIContext, body),
	)
placeholder

func TestApplyGrokCacheIdentityWritesResponsesBodyAndHeader(t *testing.T) {
	sourceBody := []byte(`{"model":"grok-4.5","prompt_cache_key":"raw-client-key"placeholder`)
	body, err := applyGrokResponsesCacheIdentity(sourceBody, sourceBody, "isolated-id", true)
placeholder
	require.Equal(t, "isolated-id", gjson.GetBytes(body, "prompt_cache_key").String())
	require.Equal(t, "web_search", gjson.GetBytes(body, "tools.0.type").String())
	require.Equal(t, "x_search", gjson.GetBytes(body, "tools.1.type").String())
	require.Equal(t, grokFreeCacheDisabledToolChoice, gjson.GetBytes(body, "tool_choice").String())

	headers := make(http.Header)
	headers.Set(grokConversationIDHeader, "spoofed-client-value")
	applyGrokCacheHeaders(headers, "isolated-id")
	require.Equal(t, "isolated-id", headers.Get(grokConversationIDHeader))
	applyGrokCacheHeaders(headers, "")
	require.Empty(t, headers.Get(grokConversationIDHeader))

	chatBody, err := stripGrokChatPromptCacheKey(body)
placeholder
	require.False(t, gjson.GetBytes(chatBody, "prompt_cache_key").Exists())

	unscopedSourceBody := []byte(`{"model":"grok","prompt_cache_key":"raw-client-key"placeholder`)
	unscopedBody, err := applyGrokResponsesCacheIdentity(unscopedSourceBody, unscopedSourceBody, "", true)
placeholder
	require.False(t, gjson.GetBytes(unscopedBody, "prompt_cache_key").Exists())
	require.False(t, gjson.GetBytes(unscopedBody, "tools").Exists())
	require.False(t, gjson.GetBytes(unscopedBody, "tool_choice").Exists())
placeholder

func TestGrokFreeMessagesClientToolCacheDefaultsOnForKnownFree(t *testing.T) {
	account := healthyGrokOAuthGatewayTestAccount(901, "access-token")
	account.Credentials["subscription_tier"] = " FREE "
	tests := []struct {
		name           string
		toolChoiceJSON string
		wantChoice     bool
placeholder{
		{name: "missing tool choice"placeholder,
		{name: "automatic tool choice", toolChoiceJSON: `,"tool_choice":"auto"`, wantChoice: trueplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intentBody := []byte(`{"model":"grok","tools":[{"type":"function","name":"lookup","description":"look up a value","parameters":{"type":"object"placeholderplaceholder,{"type":"function","name":"save","parameters":{"type":"object"placeholderplaceholder]` + tt.toolChoiceJSON + `placeholder`)
			body, err := applyGrokResponsesCacheIdentity(intentBody, intentBody, "isolated-id", true)
		placeholder
			body, err = applyGrokFreeMessagesFunctionToolCacheRoute(body, intentBody, account, "isolated-id")

		placeholder
			require.Equal(t, "isolated-id", gjson.GetBytes(body, "prompt_cache_key").String())
			tools := gjson.GetBytes(body, "tools").Array()
			require.Len(t, tools, 4)
			require.Equal(t, "function", tools[0].Get("type").String())
			require.Equal(t, "lookup", tools[0].Get("name").String())
			require.Equal(t, "function", tools[1].Get("type").String())
			require.Equal(t, "save", tools[1].Get("name").String())
			require.Equal(t, "web_search", tools[2].Get("type").String())
			require.Equal(t, "x_search", tools[3].Get("type").String())
			require.Equal(t, tt.wantChoice, gjson.GetBytes(body, "tool_choice").Exists())
	placeholder)
placeholder
placeholder

func TestGrokFreeMessagesClientToolCacheDefaultsWithMissingAccountSetting(t *testing.T) {
	body := []byte(`{"model":"grok","tools":[{"type":"function","name":"view_image","parameters":{"type":"object"placeholderplaceholder],"tool_choice":"auto"placeholder`)
	tests := []struct {
		name  string
		extra map[string]any
placeholder{
		{name: "nil extra"placeholder,
		{name: "empty extra", extra: map[string]any{placeholderplaceholder,
		{name: "unrelated extra", extra: map[string]any{"other_setting": trueplaceholderplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := healthyGrokOAuthGatewayTestAccount(90101, "access-token")
			account.Credentials["subscription_tier"] = "free"
			account.Extra = tt.extra

			patched, err := applyGrokFreeMessagesFunctionToolCacheRoute(body, body, account, "isolated-id")

		placeholder
			tools := gjson.GetBytes(patched, "tools").Array()
			require.Len(t, tools, 3)
			require.Equal(t, "web_search", tools[1].Get("type").String())
			require.Equal(t, "x_search", tools[2].Get("type").String())
	placeholder)
placeholder
placeholder

func TestApplyGrokCacheIdentityAppendsNativeToolsWhenSearchPresent(t *testing.T) {
	account := healthyGrokOAuthGatewayTestAccount(901, "access-token")
	account.Credentials["subscription_tier"] = " FREE "

	// Function tools INCLUDING web_search → convert + complement with x_search.
	intentBody := []byte(`{"model":"grok","tools":[{"type":"function","name":"lookup","description":"look up a value","parameters":{"type":"object"placeholderplaceholder,{"type":"function","name":"web_search","description":"search","parameters":{"type":"object"placeholderplaceholder]placeholder`)
	body, err := applyGrokResponsesCacheIdentity(intentBody, intentBody, "isolated-id", true)
placeholder
	body, err = applyGrokFreeMessagesFunctionToolCacheRoute(body, intentBody, account, "isolated-id")
placeholder

	tools := gjson.GetBytes(body, "tools").Array()
	require.Len(t, tools, 3, "lookup(function) + web_search(native) + x_search(native)")
	require.Equal(t, "function", tools[0].Get("type").String())
	require.Equal(t, "lookup", tools[0].Get("name").String())
	require.Equal(t, "web_search", tools[1].Get("type").String())
	require.Equal(t, "x_search", tools[2].Get("type").String())
placeholder

func TestGrokFreeClientToolCacheAccountOptIn(t *testing.T) {
	account := healthyGrokOAuthGatewayTestAccount(9011, "access-token")
	account.Credentials["subscription_tier"] = "free"
	account.Extra = map[string]any{grokClientToolCacheOptInExtraKey: trueplaceholder
	intentBody := []byte(`{"model":"grok","tools":[{"type":"function","name":"view_image","parameters":{"type":"object"placeholderplaceholder,{"type":"function","name":"read_file","parameters":{"type":"object"placeholderplaceholder],"tool_choice":"auto"placeholder`)

	body, err := applyGrokResponsesCacheIdentity(intentBody, intentBody, "isolated-id", true)
placeholder
	body, err = applyGrokFreeMessagesFunctionToolCacheRoute(body, intentBody, account, "isolated-id")
placeholder

	tools := gjson.GetBytes(body, "tools").Array()
	require.Len(t, tools, 4)
	require.Equal(t, "view_image", tools[0].Get("name").String())
	require.Equal(t, "read_file", tools[1].Get("name").String())
	require.Equal(t, "web_search", tools[2].Get("type").String())
	require.Equal(t, "x_search", tools[3].Get("type").String())
placeholder

func TestGrokFreeMessagesClientToolCacheAccountOptOut(t *testing.T) {
	body := []byte(`{"model":"grok","tools":[{"type":"function","name":"view_image","parameters":{"type":"object"placeholderplaceholder],"tool_choice":"auto"placeholder`)
	tests := []struct {
		name  string
		value any
placeholder{
		{name: "explicit false", value: falseplaceholder,
		{name: "string false is malformed", value: "false"placeholder,
		{name: "numeric value is malformed", value: 1placeholder,
		{name: "null value is malformed", value: nilplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := healthyGrokOAuthGatewayTestAccount(90111, "access-token")
			account.Credentials["subscription_tier"] = "free"
			account.Extra = map[string]any{grokClientToolCacheOptInExtraKey: tt.valueplaceholder

			patched, err := applyGrokFreeMessagesFunctionToolCacheRoute(body, body, account, "isolated-id")

		placeholder
			require.JSONEq(t, string(body), string(patched))
	placeholder)
placeholder
placeholder

func TestGrokFreeClientToolCacheRequestOptInOverridesAccountOptOut(t *testing.T) {
	account := healthyGrokOAuthGatewayTestAccount(9014, "access-token")
	account.Credentials["subscription_tier"] = "free"
	account.Extra = map[string]any{grokClientToolCacheOptInExtraKey: falseplaceholder
	c := newGrokCacheTestContext(9014)
	c.Request.Header.Set(grokClientToolCacheOptInHeader, "prefer-cache")
	intentBody := []byte(`{"model":"grok","tools":[{"type":"function","name":"view_image","parameters":{"type":"object"placeholderplaceholder],"tool_choice":"auto"placeholder`)

	body, err := applyGrokResponsesCacheIdentity(intentBody, intentBody, "isolated-id", true)
placeholder
	body, err = applyGrokFreeRequestToolCacheRoute(c, body, intentBody, account, "isolated-id")
placeholder

	tools := gjson.GetBytes(body, "tools").Array()
	require.Len(t, tools, 3)
	require.Equal(t, "view_image", tools[0].Get("name").String())
	require.Equal(t, "web_search", tools[1].Get("type").String())
	require.Equal(t, "x_search", tools[2].Get("type").String())
placeholder

func TestGrokFreeChatRequestClientToolCacheDefaultsOnWithoutClientFingerprint(t *testing.T) {
	account := healthyGrokOAuthGatewayTestAccount(90140, "access-token")
	account.Credentials["subscription_tier"] = "free"
	c := newGrokCacheTestContext(90140)
	c.Request.URL.Path = "/v1/chat/completions"
	body := []byte(`{"model":"grok","tools":[{"type":"function","name":"view_image","parameters":{"type":"object"placeholderplaceholder],"tool_choice":"auto"placeholder`)

	patched, err := applyGrokFreeRequestToolCacheRoute(c, body, body, account, "isolated-id")

placeholder
	tools := gjson.GetBytes(patched, "tools").Array()
	require.Len(t, tools, 3)
	require.Equal(t, "view_image", tools[0].Get("name").String())
	require.Equal(t, "web_search", tools[1].Get("type").String())
	require.Equal(t, "x_search", tools[2].Get("type").String())
placeholder

func TestGrokFreeClientToolCacheClaudeDesktopResponsesAutoOptIn(t *testing.T) {
	account := healthyGrokOAuthGatewayTestAccount(90141, "access-token")
	account.Credentials["subscription_tier"] = "free"
	body := []byte(`{"model":"grok","tools":[{"type":"function","name":"Read","parameters":{"type":"object"placeholderplaceholder,{"type":"function","name":"Edit","parameters":{"type":"object"placeholderplaceholder],"tool_choice":"auto"placeholder`)

	for _, xApp := range []string{"cli", "cli-bg"placeholder {
		t.Run(xApp, func(t *testing.T) {
			c := newGrokCacheTestContext(90141)
			// The desktop marker text is intentionally not required; CC Switch and
			// Claude Desktop may change the descriptive User-Agent suffix.
			c.Request.Header.Set("User-Agent", "claude-cli/2.1.215 (external, future-desktop, agent-sdk/0.3.215)")
			c.Request.Header.Set("X-App", xApp)
			c.Request.Header.Set("anthropic-client-platform", "desktop_app")
			c.Request.Header.Set("X-Claude-Code-Session-Id", "desktop-session-1")

			patched, err := applyGrokFreeRequestToolCacheRoute(c, body, body, account, "isolated-id")

		placeholder
			tools := gjson.GetBytes(patched, "tools").Array()
			require.Len(t, tools, 4)
			require.Equal(t, "Read", tools[0].Get("name").String())
			require.Equal(t, "Edit", tools[1].Get("name").String())
			require.Equal(t, "web_search", tools[2].Get("type").String())
			require.Equal(t, "x_search", tools[3].Get("type").String())
	placeholder)
placeholder
placeholder

func TestGrokFreeClientToolCacheClaudeDesktopFingerprintRequiresAllSignals(t *testing.T) {
	account := healthyGrokOAuthGatewayTestAccount(90142, "access-token")
	account.Credentials["subscription_tier"] = "free"
	account.Extra = map[string]any{grokClientToolCacheOptInExtraKey: falseplaceholder
	body := []byte(`{"model":"grok","tools":[{"type":"function","name":"Read","parameters":{"type":"object"placeholderplaceholder],"tool_choice":"auto"placeholder`)

	tests := []struct {
		name     string
		path     string
		ua       string
		xApp     string
		platform string
		session  string
placeholder{
		{
			name:     "chat path",
			path:     "/v1/chat/completions",
			ua:       "claude-cli/2.1.215 (external, claude-desktop-3p, agent-sdk/0.3.215)",
			xApp:     "cli",
			platform: "desktop_app",
			session:  "desktop-session-1",
	placeholder,
		{
			name:     "compact responses path",
			path:     "/v1/responses/compact",
			ua:       "claude-cli/2.1.215 (external, claude-desktop-3p, agent-sdk/0.3.215)",
			xApp:     "cli",
			platform: "desktop_app",
			session:  "desktop-session-1",
	placeholder,
		{
			name:     "non claude cli user agent",
			path:     "/v1/responses",
			ua:       "Mozilla/5.0 (claude-desktop-3p)",
			xApp:     "cli",
			platform: "desktop_app",
			session:  "desktop-session-1",
	placeholder,
		{
			name:     "missing x app",
			path:     "/v1/responses",
			ua:       "claude-cli/2.1.215 (external, claude-desktop-3p, agent-sdk/0.3.215)",
			platform: "desktop_app",
			session:  "desktop-session-1",
	placeholder,
		{
			name:     "wrong x app",
			path:     "/v1/responses",
			ua:       "claude-cli/2.1.215 (external, claude-desktop-3p, agent-sdk/0.3.215)",
			xApp:     "desktop",
			platform: "desktop_app",
			session:  "desktop-session-1",
	placeholder,
		{
			name:    "missing client platform",
			path:    "/v1/responses",
			ua:      "claude-cli/2.1.215 (external, claude-desktop-3p, agent-sdk/0.3.215)",
			xApp:    "cli",
			session: "desktop-session-1",
	placeholder,
		{
			name:     "wrong client platform",
			path:     "/v1/responses",
			ua:       "claude-cli/2.1.215 (external, claude-desktop-3p, agent-sdk/0.3.215)",
			xApp:     "cli",
			platform: "web",
			session:  "desktop-session-1",
	placeholder,
		{
			name:     "missing session header",
			path:     "/v1/responses",
			ua:       "claude-cli/2.1.215 (external, claude-desktop-3p, agent-sdk/0.3.215)",
			xApp:     "cli",
			platform: "desktop_app",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newGrokCacheTestContext(90142)
			c.Request.URL.Path = tt.path
			if tt.ua != "" {
				c.Request.Header.Set("User-Agent", tt.ua)
		placeholder
			if tt.xApp != "" {
				c.Request.Header.Set("X-App", tt.xApp)
		placeholder
			if tt.platform != "" {
				c.Request.Header.Set("anthropic-client-platform", tt.platform)
		placeholder
			if tt.session != "" {
				c.Request.Header.Set("X-Claude-Code-Session-Id", tt.session)
		placeholder

			patched, err := applyGrokFreeRequestToolCacheRoute(c, body, body, account, "isolated-id")

		placeholder
			require.JSONEq(t, string(body), string(patched))
	placeholder)
placeholder
placeholder

func TestGrokFreeClientToolCacheExplicitRequestOptOut(t *testing.T) {
	account := healthyGrokOAuthGatewayTestAccount(90144, "access-token")
	account.Credentials["subscription_tier"] = "free"
	body := []byte(`{"model":"grok","tools":[{"type":"function","name":"Read","parameters":{"type":"object"placeholderplaceholder],"tool_choice":"auto"placeholder`)

	for _, value := range []string{"0", "false", "no", "off"placeholder {
		t.Run(value, func(t *testing.T) {
			c := newGrokCacheTestContext(90144)
			c.Request.URL.Path = "/v1/chat/completions"
			c.Request.Header.Set(grokClientToolCacheOptInHeader, value)

			patched, err := applyGrokFreeRequestToolCacheRoute(c, body, body, account, "isolated-id")

		placeholder
			require.JSONEq(t, string(body), string(patched))
	placeholder)
placeholder
placeholder

func TestGrokFreeClientToolCacheClaudeDesktopAutoOptInDoesNotOverridePaidTier(t *testing.T) {
	account := healthyGrokOAuthGatewayTestAccount(90143, "access-token")
	account.Credentials["subscription_tier"] = "supergrok"
	c := newGrokCacheTestContext(90143)
	c.Request.Header.Set("User-Agent", "claude-cli/2.1.215 (external, claude-desktop-3p, agent-sdk/0.3.215)")
	c.Request.Header.Set("X-App", "cli")
	c.Request.Header.Set("anthropic-client-platform", "desktop_app")
	c.Request.Header.Set("X-Claude-Code-Session-Id", "desktop-session-1")
	body := []byte(`{"model":"grok","tools":[{"type":"function","name":"Read","parameters":{"type":"object"placeholderplaceholder],"tool_choice":"auto"placeholder`)

	patched, err := applyGrokFreeRequestToolCacheRoute(c, body, body, account, "isolated-id")

placeholder
	require.JSONEq(t, string(body), string(patched))
placeholder

func TestGrokFreeRequestClientSearchFunctionUsesDefaultAccountPolicy(t *testing.T) {
	account := healthyGrokOAuthGatewayTestAccount(9015, "access-token")
	account.Credentials["subscription_tier"] = "free"
	c := newGrokCacheTestContext(9015)
	body := []byte(`{"model":"grok","tools":[{"type":"function","name":"view_image","parameters":{"type":"object"placeholderplaceholder,{"type":"function","name":"web_search","parameters":{"type":"object"placeholderplaceholder],"tool_choice":"auto"placeholder`)

	patched, err := applyGrokFreeRequestToolCacheRoute(c, body, body, account, "isolated-id")

placeholder
	tools := gjson.GetBytes(patched, "tools").Array()
	require.Len(t, tools, 3)
	require.Equal(t, "view_image", tools[0].Get("name").String())
	require.Equal(t, "web_search", tools[1].Get("type").String())
	require.Equal(t, "x_search", tools[2].Get("type").String())
placeholder

func TestGrokFreeRequestToolChoiceNoneUsesSafeCacheRoute(t *testing.T) {
	account := healthyGrokOAuthGatewayTestAccount(9016, "access-token")
	account.Credentials["subscription_tier"] = "free"
	c := newGrokCacheTestContext(9016)
	body := []byte(`{"model":"grok","tools":[{"type":"function","name":"view_image","parameters":{"type":"object"placeholderplaceholder],"tool_choice":"none"placeholder`)

	patched, err := applyGrokFreeRequestToolCacheRoute(c, body, body, account, "isolated-id")

placeholder
	tools := gjson.GetBytes(patched, "tools").Array()
	require.Len(t, tools, 3)
	require.Equal(t, "view_image", tools[0].Get("name").String())
	require.Equal(t, "web_search", tools[1].Get("type").String())
	require.Equal(t, "x_search", tools[2].Get("type").String())
	require.Equal(t, "none", gjson.GetBytes(patched, "tool_choice").String())
placeholder

func TestApplyGrokCacheIdentityRecognizesResponsesLiteAdditionalTools(t *testing.T) {
	intentBody := []byte(`{"model":"grok","input":[{"type":"additional_tools","tools":[{"type":"function","name":"lookup","parameters":{"type":"object"placeholderplaceholder]placeholder,{"type":"message","role":"user","content":"hello"placeholder]placeholder`)
	patchedBody := []byte(`{"model":"grok-4.5","tools":[{"type":"function","name":"lookup","parameters":{"type":"object"placeholderplaceholder],"input":[{"type":"message","role":"user","content":"hello"placeholder]placeholder`)

	patched, err := applyGrokResponsesCacheIdentity(patchedBody, intentBody, "isolated-id", true)

placeholder
	tools := gjson.GetBytes(patched, "tools").Array()
	require.Len(t, tools, 1)
	require.Equal(t, "lookup", tools[0].Get("name").String())
	require.False(t, gjson.GetBytes(patched, "tool_choice").Exists())
	require.Equal(t, "isolated-id", gjson.GetBytes(patched, "prompt_cache_key").String())
placeholder

func TestGrokFreeCacheRoutePreservesMixedSupportedToolsWithSearchIntent(t *testing.T) {
	account := healthyGrokOAuthGatewayTestAccount(9012, "access-token")
	account.Credentials["subscription_tier"] = "free"
	intentBody := []byte(`{"model":"grok","tools":[{"type":"function","name":"view_image","parameters":{"type":"object"placeholderplaceholder,{"type":"shell"placeholder,{"type":"web_search"placeholder],"tool_choice":"auto"placeholder`)

	body, err := applyGrokResponsesCacheIdentity(intentBody, intentBody, "isolated-id", true)
placeholder
	body, err = applyGrokFreeMessagesFunctionToolCacheRoute(body, intentBody, account, "isolated-id")
placeholder

	tools := gjson.GetBytes(body, "tools").Array()
	require.Len(t, tools, 4)
	require.Equal(t, "function", tools[0].Get("type").String())
	require.Equal(t, "shell", tools[1].Get("type").String())
	require.Equal(t, "web_search", tools[2].Get("type").String())
	require.Equal(t, "x_search", tools[3].Get("type").String())
placeholder

func TestGrokClientToolCacheOptInDoesNotOverridePaidTier(t *testing.T) {
	account := healthyGrokOAuthGatewayTestAccount(9013, "access-token")
	account.Credentials["subscription_tier"] = "supergrok"
	account.Extra = map[string]any{grokClientToolCacheOptInExtraKey: trueplaceholder
	body := []byte(`{"model":"grok","tools":[{"type":"function","name":"view_image","parameters":{"type":"object"placeholderplaceholder],"tool_choice":"auto"placeholder`)

	patched, err := applyGrokFreeMessagesFunctionToolCacheRoute(body, body, account, "isolated-id")

placeholder
	require.JSONEq(t, string(body), string(patched))
placeholder

func TestApplyGrokCacheIdentityRequiresPatchedFunctionTools(t *testing.T) {
	account := healthyGrokOAuthGatewayTestAccount(902, "access-token")
	account.Credentials["subscription_tier"] = "free"
	intentBody := []byte(`{"model":"grok","tools":[{"type":"function","name":"lookup"placeholder],"tool_choice":"auto"placeholder`)
	tests := []struct {
		name        string
		patchedBody string
placeholder{
		{name: "missing tools", patchedBody: `{"model":"grok-4.5"placeholder`placeholder,
		{name: "empty tools", patchedBody: `{"model":"grok-4.5","tools":[]placeholder`placeholder,
		{name: "native tools only", patchedBody: `{"model":"grok-4.5","tools":[{"type":"web_search"placeholder]placeholder`placeholder,
		{name: "unexpected patched tool", patchedBody: `{"model":"grok-4.5","tools":[{"type":"function","name":"lookup"placeholder,{"type":"namespace","name":"server"placeholder]placeholder`placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeTools := gjson.Get(tt.patchedBody, "tools")
			body, err := applyGrokResponsesCacheIdentity([]byte(tt.patchedBody), intentBody, "isolated-id", true)
		placeholder
			body, err = applyGrokFreeMessagesFunctionToolCacheRoute(body, intentBody, account, "isolated-id")

		placeholder
			require.Equal(t, "isolated-id", gjson.GetBytes(body, "prompt_cache_key").String())
			afterTools := gjson.GetBytes(body, "tools")
			require.Equal(t, beforeTools.Exists(), afterTools.Exists())
			require.Equal(t, beforeTools.Raw, afterTools.Raw)
	placeholder)
placeholder
placeholder

func TestGrokFreeMessagesFunctionToolCacheRouteRequiresKnownFreeTier(t *testing.T) {
	// Include web_search as a function to also cover its conversion to a native tool.
	intentBody := []byte(`{"model":"grok","tools":[{"type":"function","name":"lookup"placeholder,{"type":"function","name":"web_search"placeholder],"tool_choice":"auto"placeholder`)
	tests := []struct {
		name    string
		account *Account
		wantMix bool
placeholder{
		{
			name: "free credential tier",
			account: func() *Account {
				a := healthyGrokOAuthGatewayTestAccount(910, "access-token")
				a.Credentials["subscription_tier"] = "free"
				return a
		placeholder(),
			wantMix: true,
	placeholder,
		{
			name: "free billing tier",
			account: func() *Account {
				a := healthyGrokOAuthGatewayTestAccount(911, "access-token")
				a.Extra = map[string]any{grokBillingExtraKey: map[string]any{"plan": "FREE"placeholderplaceholder
				return a
		placeholder(),
			wantMix: true,
	placeholder,
		{
			name: "free successful billing has blank plan",
			account: func() *Account {
				a := healthyGrokOAuthGatewayTestAccount(9111, "access-token")
				a.Extra = map[string]any{grokBillingExtraKey: map[string]any{
					"status_code":        http.StatusOK,
					"source":             "billing_probe",
					"monthly_updated_at": "2026-07-15T05:00:00Z",
			placeholderplaceholder
				return a
		placeholder(),
			wantMix: true,
	placeholder,
		{
			name: "free rolling token quota",
			account: func() *Account {
				a := healthyGrokOAuthGatewayTestAccount(9112, "access-token")
				a.Extra = map[string]any{grokQuotaSnapshotExtraKey: map[string]any{
					"headers_observed": true,
					"tokens":           map[string]any{"limit": xai.GrokFreeRolling24hTokenLimitplaceholder,
			placeholderplaceholder
				return a
		placeholder(),
			wantMix: true,
	placeholder,
		{
			name: "legacy free rolling token quota",
			account: func() *Account {
				a := healthyGrokOAuthGatewayTestAccount(9113, "access-token")
				a.Extra = map[string]any{grokQuotaSnapshotExtraKey: map[string]any{
					"headers_observed": true,
					"tokens":           map[string]any{"limit": int64(2_000_000)placeholder,
			placeholderplaceholder
				return a
		placeholder(),
			wantMix: true,
	placeholder,
		{
			name: "supergrok remains unchanged",
			account: func() *Account {
				a := healthyGrokOAuthGatewayTestAccount(912, "access-token")
				a.Credentials["subscription_tier"] = "supergrok"
				return a
		placeholder(),
	placeholder,
		{
			name: "paid billing overrides stale free quota",
			account: func() *Account {
				a := healthyGrokOAuthGatewayTestAccount(9121, "access-token")
				a.Extra = map[string]any{
					grokBillingExtraKey: map[string]any{"plan": "SuperGrok", "status_code": http.StatusOKplaceholder,
					grokQuotaSnapshotExtraKey: map[string]any{
						"headers_observed": true,
						"tokens":           map[string]any{"limit": int64(2_000_000)placeholder,
				placeholder,
			placeholder
				return a
		placeholder(),
	placeholder,
		{
			name: "paid billing overrides stale free credentials",
			account: func() *Account {
				a := healthyGrokOAuthGatewayTestAccount(9123, "access-token")
				a.Credentials["subscription_tier"] = "free"
				a.Extra = map[string]any{
					grokBillingExtraKey: map[string]any{"plan": "SuperGrok", "status_code": http.StatusOKplaceholder,
			placeholder
				return a
		placeholder(),
	placeholder,
		{
			name: "partial billing without monthly evidence remains unknown",
			account: func() *Account {
				a := healthyGrokOAuthGatewayTestAccount(9122, "access-token")
				a.Extra = map[string]any{grokBillingExtraKey: map[string]any{
					"status_code":    http.StatusOK,
					"source":         "billing_probe",
					"partial":        true,
					"failed_windows": []string{"monthly"placeholder,
			placeholderplaceholder
				return a
		placeholder(),
	placeholder,
		{
			name:    "unknown tier remains unchanged",
			account: healthyGrokOAuthGatewayTestAccount(913, "access-token"),
	placeholder,
		{
			name: "api key remains unchanged",
			account: &Account{
				ID:       914,
				Platform: PlatformGrok,
				Type:     AccountTypeAPIKey,
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := applyGrokFreeMessagesFunctionToolCacheRoute(intentBody, intentBody, tt.account, "isolated-id")

		placeholder
			tools := gjson.GetBytes(body, "tools").Array()
			if tt.wantMix {
				require.Len(t, tools, 3)
				require.Equal(t, "web_search", tools[1].Get("type").String())
				require.Equal(t, "x_search", tools[2].Get("type").String())
				return
		placeholder
			require.Len(t, tools, 2, "non-free accounts should not get native search injected")
	placeholder)
placeholder
placeholder

func TestGrokFreeMessagesFunctionToolCacheRouteRequiresIdentity(t *testing.T) {
	account := healthyGrokOAuthGatewayTestAccount(915, "access-token")
	account.Credentials["subscription_tier"] = "free"
	body := []byte(`{"model":"grok","tools":[{"type":"function","name":"lookup"placeholder],"tool_choice":"auto"placeholder`)

	patched, err := applyGrokFreeMessagesFunctionToolCacheRoute(body, body, account, "")

placeholder
	require.JSONEq(t, string(body), string(patched))
	require.Len(t, gjson.GetBytes(patched, "tools").Array(), 1)
placeholder

func TestApplyGrokCacheIdentityPreservesIneligibleClientToolFields(t *testing.T) {
	tests := []struct {
		name string
		body string
placeholder{
		{
			name: "empty tools array",
			body: `{"model":"grok","tools":[]placeholder`,
	placeholder,
		{
			name: "null tools",
			body: `{"model":"grok","tools":nullplaceholder`,
	placeholder,
		{
			name: "tool choice only",
			body: `{"model":"grok","tool_choice":{"type":"function","name":"lookup"placeholderplaceholder`,
	placeholder,
		{
			name: "null tool choice",
			body: `{"model":"grok","tool_choice":nullplaceholder`,
	placeholder,
		{
			name: "native tool with auto choice",
			body: `{"model":"grok","tools":[{"type":"web_search"placeholder],"tool_choice":"auto"placeholder`,
	placeholder,
		{
			name: "function with required choice",
			body: `{"model":"grok","tools":[{"type":"function","name":"lookup"placeholder],"tool_choice":"required"placeholder`,
	placeholder,
		{
			name: "function with none choice",
			body: `{"model":"grok","tools":[{"type":"function","name":"lookup"placeholder],"tool_choice":"none"placeholder`,
	placeholder,
		{
			name: "function with specific choice",
			body: `{"model":"grok","tools":[{"type":"function","name":"lookup"placeholder],"tool_choice":{"type":"function","name":"lookup"placeholderplaceholder`,
	placeholder,
		{
			name: "function with object auto choice",
			body: `{"model":"grok","tools":[{"type":"function","name":"lookup"placeholder],"tool_choice":{"type":"auto"placeholderplaceholder`,
	placeholder,
		{
			name: "function mixed with unsupported tool",
			body: `{"model":"grok","tools":[{"type":"function","name":"lookup"placeholder,{"type":"namespace","name":"client_tools"placeholder],"tool_choice":"auto"placeholder`,
	placeholder,
		{
			name: "unsupported tool only",
			body: `{"model":"grok","tools":[{"type":"namespace","name":"client_tools"placeholder]placeholder`,
	placeholder,
		{
			name: "chat completions function shape",
			body: `{"model":"grok","tools":[{"type":"function","function":{"name":"lookup","parameters":{"type":"object"placeholderplaceholderplaceholder],"tool_choice":"auto"placeholder`,
	placeholder,
		{
			name: "incomplete responses function",
			body: `{"model":"grok","tools":[{"type":"function","parameters":{"type":"object"placeholderplaceholder]placeholder`,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeTools := gjson.Get(tt.body, "tools")
			beforeChoice := gjson.Get(tt.body, "tool_choice")
			body, err := applyGrokResponsesCacheIdentity([]byte(tt.body), []byte(tt.body), "isolated-id", true)
		placeholder
			require.Equal(t, "isolated-id", gjson.GetBytes(body, "prompt_cache_key").String())
			require.Equal(t, beforeTools.Exists(), gjson.GetBytes(body, "tools").Exists())
			require.Equal(t, beforeTools.Raw, gjson.GetBytes(body, "tools").Raw)
			require.Equal(t, beforeChoice.Exists(), gjson.GetBytes(body, "tool_choice").Exists())
			require.Equal(t, beforeChoice.Raw, gjson.GetBytes(body, "tool_choice").Raw)
	placeholder)
placeholder
placeholder

func TestApplyGrokCacheIdentityUsesPreSanitizationToolIntent(t *testing.T) {
	tests := []struct {
		name       string
		intentBody string
placeholder{
		{
			name:       "unsupported tools removed by sanitizer",
			intentBody: `{"model":"grok","tools":[{"type":"namespace","name":"client_tools"placeholder]placeholder`,
	placeholder,
		{
			name:       "tool choice removed with unsupported tool",
			intentBody: `{"model":"grok","tools":[{"type":"namespace","name":"client_tools"placeholder],"tool_choice":{"type":"namespace","name":"client_tools"placeholderplaceholder`,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is the shape apply receives after patchGrokResponsesBody has
			// removed unsupported tools and their associated tool_choice.
			patchedBody := []byte(`{"model":"grok-4.5","input":"hello"placeholder`)
			body, err := applyGrokResponsesCacheIdentity(patchedBody, []byte(tt.intentBody), "isolated-id", true)

		placeholder
			require.Equal(t, "isolated-id", gjson.GetBytes(body, "prompt_cache_key").String())
			require.False(t, gjson.GetBytes(body, "tools").Exists())
			require.False(t, gjson.GetBytes(body, "tool_choice").Exists())
	placeholder)
placeholder
placeholder

func TestApplyGrokCacheIdentityWithoutFreeTierRoutingOnlyWritesIdentity(t *testing.T) {
	sourceBody := []byte(`{"model":"grok-4.5","input":"hello"placeholder`)
	body, err := applyGrokResponsesCacheIdentity(sourceBody, sourceBody, "isolated-id", false)

placeholder
	require.Equal(t, "isolated-id", gjson.GetBytes(body, "prompt_cache_key").String())
	require.False(t, gjson.GetBytes(body, "tools").Exists())
	require.False(t, gjson.GetBytes(body, "tool_choice").Exists())
placeholder

func TestGrokCompactRequestSkipsCacheIdentityAndNativeTools(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := newGrokCacheTestContext(701)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses/compact", nil)
	body := []byte(`{"model":"grok","input":"compact this","prompt_cache_key":"raw-client-key"placeholder`)

	identity := resolveGrokCacheIdentity(c, body, "", "grok-4.5")
	patched, err := applyGrokResponsesCacheIdentity(body, body, identity, true)

placeholder
	require.Empty(t, identity)
	require.False(t, gjson.GetBytes(patched, "prompt_cache_key").Exists())
	require.False(t, gjson.GetBytes(patched, "tools").Exists())
	require.False(t, gjson.GetBytes(patched, "tool_choice").Exists())
placeholder

func TestResolveGrokCacheIdentityConcurrentDeterminism(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const workers = 50
	body := []byte(`{"model":"grok","messages":[{"role":"system","content":"stable"placeholder,{"role":"user","content":"hello"placeholder]placeholder`)
	identities := make(chan string, workers)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			identities <- resolveGrokCacheIdentity(newGrokCacheTestContext(501), body, "", "grok-4.5")
	placeholder()
placeholder
	wg.Wait()
	close(identities)

	var first string
	for identity := range identities {
		if first == "" {
			first = identity
	placeholder
		require.Equal(t, first, identity)
placeholder
	require.NotEmpty(t, first)
placeholder
