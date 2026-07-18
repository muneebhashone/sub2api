package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	grokConversationIDHeader        = "X-Grok-Conv-Id"
	grokFreeCacheNativeToolsJSON    = `[{"type":"web_search"placeholder,{"type":"x_search"placeholder]`
	grokFreeCacheDisabledToolChoice = "none"
	grokFreeRolling24hTokenLimit    = int64(2_000_000)
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
	if gjson.GetBytes(intentSourceBody, "tools").Exists() || gjson.GetBytes(intentSourceBody, "tool_choice").Exists() {
		return out, nil
placeholder
	out, err = sjson.SetRawBytes(out, "tools", []byte(grokFreeCacheNativeToolsJSON))
	if err != nil {
		return nil, err
placeholder
	return sjson.SetBytes(out, "tool_choice", grokFreeCacheDisabledToolChoice)
placeholder

// applyGrokFreeMessagesFunctionToolCacheRoute enables xAI's cache-capable
// mixed-tools route only for the Anthropic Messages bridge and only when the
// selected account is known to be Free. Native tools become eligible under
// auto selection, so callers must not apply this policy to paid accounts or
// other ingress protocols implicitly.
func applyGrokFreeMessagesFunctionToolCacheRoute(body, intentSourceBody []byte, account *Account, cacheIdentity string) ([]byte, error) {
	if strings.TrimSpace(cacheIdentity) == "" || !isKnownGrokFreeAccount(account) {
		return body, nil
placeholder
	intentTools := gjson.GetBytes(intentSourceBody, "tools")
	intentToolChoice := gjson.GetBytes(intentSourceBody, "tool_choice")
	if !isGrokFreeCacheFunctionToolIntent(intentTools, intentToolChoice) {
		return body, nil
placeholder
	return appendMissingGrokFreeCacheNativeTools(body)
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
			*snapshot.Tokens.Limit == grokFreeRolling24hTokenLimit {
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
		if !tool.IsObject() || strings.TrimSpace(tool.Get("type").String()) != "function" {
			return false
	placeholder
		// Responses function declarations keep name at the top level. Reject
		// Chat Completions' nested function shape and incomplete declarations.
		if strings.TrimSpace(tool.Get("name").String()) == "" || tool.Get("function").Exists() {
			return false
	placeholder
placeholder
	if !toolChoice.Exists() {
		return true
placeholder
	return toolChoice.Type == gjson.String && strings.TrimSpace(toolChoice.String()) == "auto"
placeholder

func appendMissingGrokFreeCacheNativeTools(body []byte) ([]byte, error) {
	tools := gjson.GetBytes(body, "tools")
	if !tools.Exists() || !tools.IsArray() {
		return body, nil
placeholder

	items := tools.Array()
	if len(items) == 0 {
		return body, nil
placeholder
	merged := make([]json.RawMessage, 0, len(items)+2)
	present := make(map[string]bool, 2)
	hasFunction := false
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
			if name == "web_search" || name == "x_search" {
				if present[name] {
					continue
			placeholder
				raw, err := json.Marshal(map[string]string{"type": nameplaceholder)
				if err != nil {
					return nil, err
			placeholder
				merged = append(merged, raw)
				present[name] = true
				continue
		placeholder
			hasFunction = true
			merged = append(merged, json.RawMessage(tool.Raw))
		case "web_search", "x_search":
			if present[toolType] {
				continue
		placeholder
			merged = append(merged, json.RawMessage(tool.Raw))
			present[toolType] = true
		default:
			return body, nil
	placeholder
placeholder
	if !hasFunction {
		return body, nil
placeholder
	// Only complement missing native search tools when the request already contains
	// at least one search tool (native or function-form). Pure client function tools
	// (e.g. view_image) must not trigger injection to avoid biasing model tool
	// selection (#4486).
	if !present["web_search"] && !present["x_search"] {
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
