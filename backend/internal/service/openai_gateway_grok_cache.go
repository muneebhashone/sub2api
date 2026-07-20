package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	grokConversationIDHeader         = "X-Grok-Conv-Id"
	grokClientToolCacheOptInHeader   = "X-Sub2API-Grok-Client-Tool-Cache"
	grokFreeCacheNativeToolsJSON     = `[{"type":"web_search"placeholder,{"type":"x_search"placeholder]`
	grokFreeCacheDisabledToolChoice  = "none"
	grokClientToolCacheOptInExtraKey = "grok_client_tool_cache_enabled"
)

// resolveGrokCacheIdentity derives one stable, tenant-isolated routing identity
// for xAI's server-side prompt cache. The returned value is safe to expose to
// the upstream: it never contains the client's raw session identifier.
//
// A valid downstream API key is required. This intentionally fails closed on
// internal probes and incomplete request contexts instead of creating a cache
// identity that could be shared by unrelated tenants.
func resolveGrokCacheIdentity(c *gin.Context, body []byte, explicitKey, upstreamModel string) string {
	apiKeyID := getAPIKeyIDFromContext(c)
	if apiKeyID <= 0 {
		return ""
placeholder
	// /responses/compact rejects tool_choice and does not represent a normal
	// conversation turn. Keep both cache identity and Free-tier routing
	// augmentation out of this path.
	if isOpenAIResponsesCompactPath(c) {
		return ""
placeholder

	model := strings.ToLower(strings.TrimSpace(upstreamModel))
	if model == "" {
		return ""
placeholder

	seed := explicitGrokCacheSeed(c, body, explicitKey)
	if seed == "" {
		seed = deriveOpenAIStablePrefixSessionSeed(body)
		if seed == "" {
			// A model alone is too broad for cache routing. Preserve the
			// existing first-user-derived identity when no reusable prefix is
			// available so unrelated prompts do not share one tenant-wide key.
			seed = deriveOpenAIAnchoredContentSessionSeed(body)
	placeholder
placeholder
	if seed == "" {
		return ""
placeholder

	// generateSessionUUID hashes the whole seed before formatting it as a UUID.
	// Include a versioned namespace so this identity cannot collide with other
	// upstream session identifiers derived by sub2api.
	isolatedSeed := fmt.Sprintf("grok-prompt-cache:v1:%d:%s:%s", apiKeyID, model, seed)
	return generateSessionUUID(isolatedSeed)
placeholder

func explicitGrokCacheSeed(c *gin.Context, body []byte, explicitKey string) string {
	seed := ""
	if c != nil {
		seed = strings.TrimSpace(c.GetHeader("session_id"))
		if seed == "" {
			seed = strings.TrimSpace(c.GetHeader("conversation_id"))
	placeholder
		if seed == "" {
			seed = strings.TrimSpace(c.GetHeader(grokConversationIDHeader))
	placeholder
placeholder
	if seed == "" && len(body) > 0 {
		seed = strings.TrimSpace(gjson.GetBytes(body, "prompt_cache_key").String())
placeholder
	if seed == "" {
		seed = strings.TrimSpace(explicitKey)
placeholder
	return seed
placeholder

func isGrokRequestContext(c *gin.Context) bool {
	if c == nil {
		return false
placeholder
	v, exists := c.Get("api_key")
	if !exists {
		return false
placeholder
	apiKey, ok := v.(*APIKey)
	return ok && apiKey != nil && apiKey.Group != nil && apiKey.Group.Platform == PlatformGrok
placeholder

// applyGrokResponsesCacheIdentity writes the cache routing identity into an
// xAI Responses request. Existing client values are deliberately replaced by
// the tenant-isolated value to prevent collisions on shared OAuth accounts.
//
// Free OAuth requests without native search tools are routed by xAI to the
// non-cacheable build-free model. For otherwise tool-free requests, add the
// native tools with tool_choice=none: this selects the cache-capable tier
// without allowing an actual search. Explicit client function tools are handled by
// applyGrokFreeMessagesFunctionToolCacheRoute (Messages bridge and native Responses).
func applyGrokResponsesCacheIdentity(body, intentSourceBody []byte, identity string, injectFreeTierTools bool) ([]byte, error) {
	identity = strings.TrimSpace(identity)
	if identity == "" {
		if gjson.GetBytes(body, "prompt_cache_key").Exists() {
			return sjson.DeleteBytes(body, "prompt_cache_key")
	placeholder
		return body, nil
placeholder
	out, err := sjson.SetBytes(body, "prompt_cache_key", identity)
	if err != nil {
		return nil, err
placeholder
	if !injectFreeTierTools {
		return out, nil
placeholder
	// Inspect the pre-sanitization source. patchGrokResponsesBody may remove an
	// unsupported client tool and its tool_choice; that must not turn an
	// explicit client tool intent into an eligible native-tool request.
	if hasGrokResponsesToolIntent(intentSourceBody) {
		return out, nil
placeholder
	out, err = sjson.SetRawBytes(out, "tools", []byte(grokFreeCacheNativeToolsJSON))
	if err != nil {
		return nil, err
placeholder
	return sjson.SetBytes(out, "tool_choice", grokFreeCacheDisabledToolChoice)
placeholder

func hasGrokResponsesToolIntent(body []byte) bool {
	if gjson.GetBytes(body, "tools").Exists() || gjson.GetBytes(body, "tool_choice").Exists() {
		return true
placeholder
	input := gjson.GetBytes(body, "input")
	if !input.IsArray() {
		return false
placeholder
	for _, item := range input.Array() {
		if strings.TrimSpace(item.Get("type").String()) != "additional_tools" {
			continue
	placeholder
		tools := item.Get("tools")
		if !tools.Exists() || !tools.IsArray() || len(tools.Array()) > 0 {
			return true
	placeholder
placeholder
	return false
placeholder

// applyGrokFreeMessagesFunctionToolCacheRoute enables xAI's cache-capable
// mixed-tools route only for known Free accounts. Pure client tools default to
// the cache-capable route so an intermediate sub2api does not need to preserve
// client-specific opt-in headers. Operators can explicitly disable this per
// account when native search tools would change the desired behavior (#4486).
func applyGrokFreeMessagesFunctionToolCacheRoute(body, intentSourceBody []byte, account *Account, cacheIdentity string) ([]byte, error) {
	allowPureClientTools, _ := grokClientToolCacheAccountPolicy(account)
	return applyGrokFreeToolCacheRoute(body, intentSourceBody, account, cacheIdentity, allowPureClientTools, true)
placeholder

// applyGrokFreeRequestToolCacheRoute also accepts a request-scoped opt-in. The
// sub2api header is consumed locally because buildGrokResponsesRequest only
// forwards the explicitly supported OpenAI-Beta header from downstream.
func applyGrokFreeRequestToolCacheRoute(c *gin.Context, body, intentSourceBody []byte, account *Account, cacheIdentity string) ([]byte, error) {
	allowPureClientTools, accountPolicyExplicit := grokClientToolCacheAccountPolicy(account)
	requestOptOut := false
	if c != nil {
		switch strings.ToLower(strings.TrimSpace(c.GetHeader(grokClientToolCacheOptInHeader))) {
		case "1", "true", "yes", "on", "prefer-cache":
			allowPureClientTools = true
		case "0", "false", "no", "off":
			allowPureClientTools = false
			requestOptOut = true
	placeholder
placeholder
	if !allowPureClientTools && !accountPolicyExplicit && !requestOptOut && isGrokClaudeDesktopResponsesCacheRequest(c) {
		allowPureClientTools = true
placeholder
	// A function merely named web_search/x_search is still a client function.
	// Known Free OAuth accounts use the cache route by default; a request-scoped
	// opt-in may override an account opt-out, while an explicit request opt-out
	// always wins. The legacy Claude fingerprint remains only as a compatibility
	// fallback when no account policy has been recorded (#4486).
	return applyGrokFreeToolCacheRoute(body, intentSourceBody, account, cacheIdentity, allowPureClientTools, allowPureClientTools)
placeholder

// grokClientToolCacheAccountPolicy is intentionally strict for configured
// values: only a JSON boolean is accepted. A missing key defaults on solely for
// accounts positively identified as Grok Free OAuth; paid, API-key, and unknown
// accounts remain fail-closed.
func grokClientToolCacheAccountPolicy(account *Account) (enabled, explicit bool) {
	if !isKnownGrokFreeAccount(account) {
		return false, false
placeholder
	if account.Extra == nil {
		return true, false
placeholder
	value, exists := account.Extra[grokClientToolCacheOptInExtraKey]
	if !exists {
		return true, false
placeholder
	enabled, valid := value.(bool)
	if !valid {
		return false, true
placeholder
	return enabled, true
placeholder

// isGrokClaudeDesktopResponsesCacheRequest recognizes the strict wire
// fingerprint emitted when Claude Desktop's local agent is translated by
// CC Switch into an OpenAI Responses request. Requiring every independent
// signal prevents a generic Claude-compatible client (or the Chat bridge)
// from silently opting into the mixed native/client tool route.
func isGrokClaudeDesktopResponsesCacheRequest(c *gin.Context) bool {
	if c == nil || c.Request == nil || c.Request.URL == nil || isOpenAIResponsesCompactPath(c) {
		return false
placeholder
	path := strings.TrimRight(strings.TrimSpace(c.Request.URL.Path), "/")
	if !strings.HasSuffix(path, "/responses") {
		return false
placeholder

	if !claudeCodeUAPattern.MatchString(strings.TrimSpace(c.GetHeader("User-Agent"))) {
		return false
placeholder
	switch strings.ToLower(strings.TrimSpace(c.GetHeader("X-App"))) {
	case "cli", "cli-bg":
	default:
		return false
placeholder
	if !strings.EqualFold(strings.TrimSpace(c.GetHeader("anthropic-client-platform")), "desktop_app") {
		return false
placeholder
	return strings.TrimSpace(c.GetHeader("X-Claude-Code-Session-Id")) != ""
placeholder

func applyGrokFreeToolCacheRoute(body, intentSourceBody []byte, account *Account, cacheIdentity string, allowPureClientTools, allowFunctionSearch bool) ([]byte, error) {
	if strings.TrimSpace(cacheIdentity) == "" || !isKnownGrokFreeAccount(account) {
		return body, nil
placeholder
	intentTools := gjson.GetBytes(intentSourceBody, "tools")
	intentToolChoice := gjson.GetBytes(intentSourceBody, "tool_choice")
	if !isGrokFreeCacheFunctionToolIntent(intentTools, intentToolChoice) {
		return body, nil
placeholder
	if intentToolChoice.Type == gjson.String && strings.TrimSpace(intentToolChoice.String()) == grokFreeCacheDisabledToolChoice {
		// Adding native cache-routing tools cannot change behavior when the
		// client has explicitly disabled all tool execution.
		return appendGrokFreeCacheNativeToolsWithPolicy(body, true, false)
placeholder
	return appendGrokFreeCacheNativeToolsWithPolicy(body, allowPureClientTools, allowFunctionSearch)
placeholder

func isKnownGrokFreeAccount(account *Account) bool {
	if account == nil || !account.IsGrokOAuth() {
		return false
placeholder
	freeSignal := false
	paidSignal := false
	inferredFreeSignal := false
	if billing, err := grokBillingSnapshotFromExtra(account.Extra); err == nil && billing != nil {
		if tier := strings.TrimSpace(billing.Plan); tier != "" {
			if isGrokFreeSubscriptionTier(tier) {
				freeSignal = true
		placeholder else if !isGrokUnknownSubscriptionTier(tier) {
				paidSignal = true
		placeholder
	placeholder
		if billing.UsagePercent != nil || billing.UsedPercent != nil ||
			(billing.MonthlyLimitCents != nil && *billing.MonthlyLimitCents > 0) {
			paidSignal = true
	placeholder
		// xAI deliberately reports an empty plan for Free accounts; only paid
		// subscriptions receive a SuperGrok plan/monthly limit. A successful
		// monthly billing observation with no paid signal is therefore positive
		// Free evidence, not an unknown tier. Keep partial probes fail-closed.
		if strings.TrimSpace(billing.MonthlyUpdatedAt) != "" ||
			(billing.StatusCode >= http.StatusOK && billing.StatusCode < http.StatusMultipleChoices &&
				!billing.Partial && len(billing.FailedWindows) == 0) {
			inferredFreeSignal = true
	placeholder
placeholder
	if snapshot, err := grokQuotaSnapshotFromExtra(account.Extra); err == nil && snapshot != nil {
		if tier := strings.TrimSpace(snapshot.SubscriptionTier); tier != "" {
			if isGrokFreeSubscriptionTier(tier) {
				freeSignal = true
		placeholder else if !isGrokUnknownSubscriptionTier(tier) {
				paidSignal = true
		placeholder
	placeholder
		if snapshot.Tokens != nil && snapshot.Tokens.Limit != nil &&
			xai.IsGrokFreeRolling24hTokenLimit(*snapshot.Tokens.Limit) {
			inferredFreeSignal = true
	placeholder
placeholder
	if tier := strings.TrimSpace(account.GetCredential("subscription_tier")); tier != "" {
		if isGrokFreeSubscriptionTier(tier) {
			freeSignal = true
	placeholder else if !isGrokUnknownSubscriptionTier(tier) {
			paidSignal = true
	placeholder
placeholder
	// Explicit paid evidence always wins over an inferred Free signal. This
	// protects upgraded/stale accounts whose previous quota snapshot still
	// carries the historical 2M Free token limit.
	return !paidSignal && (freeSignal || inferredFreeSignal)
placeholder

func isGrokFreeSubscriptionTier(tier string) bool {
	switch strings.ToLower(strings.TrimSpace(tier)) {
	case "free", "grok-free", "grok_free", "free-tier", "free_tier", "basic", "grok-basic", "grok_basic":
		return true
	default:
		return false
placeholder
placeholder

func isGrokUnknownSubscriptionTier(tier string) bool {
	switch strings.ToLower(strings.TrimSpace(tier)) {
	case "", "unknown", "n/a", "none":
		return true
	default:
		return false
placeholder
placeholder

func isGrokFreeCacheFunctionToolIntent(tools, toolChoice gjson.Result) bool {
	if !tools.IsArray() {
		return false
placeholder
	items := tools.Array()
	if len(items) == 0 {
		return false
placeholder
	for _, tool := range items {
		if !tool.IsObject() {
			return false
	placeholder
		toolType := strings.TrimSpace(tool.Get("type").String())
		if _, ok := grokResponsesSupportedToolTypes[toolType]; !ok {
			return false
	placeholder
		if toolType == "function" {
			// Responses function declarations keep name at the top level. Reject
			// Chat Completions' nested function shape and incomplete declarations.
			if strings.TrimSpace(tool.Get("name").String()) == "" || tool.Get("function").Exists() {
				return false
		placeholder
	placeholder
placeholder
	if !toolChoice.Exists() {
		return true
placeholder
	if toolChoice.Type != gjson.String {
		return false
placeholder
	switch strings.TrimSpace(toolChoice.String()) {
	case "auto", grokFreeCacheDisabledToolChoice:
		return true
	default:
		return false
placeholder
placeholder

func appendMissingGrokFreeCacheNativeTools(body []byte) ([]byte, error) {
	return appendGrokFreeCacheNativeTools(body, false)
placeholder

func appendGrokFreeCacheNativeTools(body []byte, allowPureClientTools bool) ([]byte, error) {
	return appendGrokFreeCacheNativeToolsWithPolicy(body, allowPureClientTools, true)
placeholder

func appendGrokFreeCacheNativeToolsWithPolicy(body []byte, allowPureClientTools, allowFunctionSearch bool) ([]byte, error) {
	tools := gjson.GetBytes(body, "tools")
	if !tools.Exists() || !tools.IsArray() {
		return body, nil
placeholder

	items := tools.Array()
	if len(items) == 0 {
		return body, nil
placeholder
	hasNativeSearch := false
	for _, tool := range items {
		switch strings.TrimSpace(tool.Get("type").String()) {
		case "web_search", "x_search":
			hasNativeSearch = true
	placeholder
placeholder
	if !allowPureClientTools && !allowFunctionSearch && !hasNativeSearch {
		return body, nil
placeholder
	merged := make([]json.RawMessage, 0, len(items)+2)
	present := make(map[string]bool, 2)
	hasCompanionTool := false
	for _, tool := range items {
		toolType := strings.TrimSpace(tool.Get("type").String())
		switch toolType {
		case "function":
			name := strings.TrimSpace(tool.Get("name").String())
			if !tool.IsObject() || name == "" || tool.Get("function").Exists() {
				return body, nil
		placeholder
			// Grok Build may declare search as function tools. Convert to native
			// entries so Free OAuth stays cache-capable without duplicate names.
			if (name == "web_search" || name == "x_search") && allowFunctionSearch {
				if present[name] {
					continue
			placeholder
				raw, err := json.Marshal(map[string]string{"type": nameplaceholder)
				if err != nil {
					return nil, err
			placeholder
				merged = append(merged, raw)
				present[name] = true
				if allowPureClientTools {
					hasCompanionTool = true
			placeholder
				continue
		placeholder
			if name == "web_search" || name == "x_search" {
				// Keep the client function intact and avoid adding a same-named
				// native tool unless conversion was explicitly enabled.
				present[name] = true
		placeholder
			hasCompanionTool = true
			merged = append(merged, json.RawMessage(tool.Raw))
		case "web_search", "x_search":
			if present[toolType] {
				continue
		placeholder
			merged = append(merged, json.RawMessage(tool.Raw))
			present[toolType] = true
		default:
			if _, ok := grokResponsesSupportedToolTypes[toolType]; !ok {
				return body, nil
		placeholder
			hasCompanionTool = true
			merged = append(merged, json.RawMessage(tool.Raw))
	placeholder
placeholder
	if !hasCompanionTool {
		return body, nil
placeholder
	// Only complement missing native search tools when the request already contains
	// at least one search tool (native or function-form). Pure client function tools
	// (e.g. view_image) must not trigger injection to avoid biasing model tool
	// selection (#4486).
	if !allowPureClientTools && !present["web_search"] && !present["x_search"] {
		return body, nil
placeholder
	for _, toolType := range []string{"web_search", "x_search"placeholder {
		if present[toolType] {
			continue
	placeholder
		raw, err := json.Marshal(map[string]string{"type": toolTypeplaceholder)
		if err != nil {
			return nil, err
	placeholder
		merged = append(merged, raw)
placeholder
	encoded, err := json.Marshal(merged)
	if err != nil {
		return nil, err
placeholder
	return sjson.SetRawBytes(body, "tools", encoded)
placeholder

// applyGrokCacheHeaders applies the documented Chat Completions conversation
// routing header. The request is built from a fresh header map, so client
// supplied x-grok headers cannot override this server-derived value.
func applyGrokCacheHeaders(headers http.Header, identity string) {
	if headers == nil {
		return
placeholder
	identity = strings.TrimSpace(identity)
	if identity == "" {
		headers.Del(grokConversationIDHeader)
		return
placeholder
	headers.Set(grokConversationIDHeader, identity)
placeholder

// stripGrokChatPromptCacheKey removes the Responses-only body field after it
// has been used as an identity seed. Chat Completions routes cache by header.
func stripGrokChatPromptCacheKey(body []byte) ([]byte, error) {
	if !gjson.GetBytes(body, "prompt_cache_key").Exists() {
		return body, nil
placeholder
	return sjson.DeleteBytes(body, "prompt_cache_key")
placeholder
