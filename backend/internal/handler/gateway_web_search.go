package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/pkg/ip"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/Wei-Shaw/sub2api/internal/pkg/websearch"
	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

const (
	defaultGrokWebSearchResults = 5
	maxGrokWebSearchResults     = 20
)

func (h *GatewayHandler) WebSearch(c *gin.Context) {
	type webSearchReq struct {
		Query      string `json:"query" binding:"required"`
		MaxResults int    `json:"max_results"`
placeholder

	var req webSearchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"type":    "invalid_request_error",
			"message": err.Error(),
	placeholderplaceholder)
		return
placeholder
	req.MaxResults = normalizeGrokWebSearchMaxResults(req.MaxResults)

	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok || apiKey == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
			"type":    "authentication_error",
			"message": "API key required",
	placeholderplaceholder)
		return
placeholder

	if apiKey.Group == nil || apiKey.Group.Platform != "grok" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"type":    "invalid_request_error",
			"message": "web search is only supported for grok groups",
	placeholderplaceholder)
		return
placeholder

	// Billing eligibility (same as other requests)
	subscription, _ := middleware2.GetSubscriptionFromContext(c)
	if err := h.billingCacheService.CheckBillingEligibility(c.Request.Context(), apiKey.User, apiKey, apiKey.Group, subscription, service.QuotaPlatform(c.Request.Context(), apiKey)); err != nil {
		status, code, message, retryAfter := billingErrorDetails(err)
		if retryAfter > 0 {
			c.Header("Retry-After", strconv.Itoa(retryAfter))
	placeholder
		c.JSON(status, gin.H{"error": gin.H{"type": code, "message": messageplaceholderplaceholder)
		return
placeholder

	// Use exactly the same scheduling as other requests (SelectAccountWithLoadAwareness handles load, rate limit, sticky, etc.)
	groupID := apiKey.GroupID
	if groupID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"type":    "invalid_request_error",
			"message": "group required",
	placeholderplaceholder)
		return
placeholder

	selected, err := h.gatewayService.SelectAccountWithLoadAwareness(c.Request.Context(), groupID, "", xai.DefaultTextModel, nil, "", 0)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": gin.H{
			"type":    "scheduling_error",
			"message": err.Error(),
	placeholderplaceholder)
		return
placeholder
	if selected == nil || selected.Account == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": gin.H{
			"type":    "scheduling_error",
			"message": "No available accounts",
	placeholderplaceholder)
		return
placeholder
	account := selected.Account
	accountReleaseFunc := selected.ReleaseFunc
	if !selected.Acquired {
		if selected.WaitPlan == nil || h.concurrencyHelper == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": gin.H{
				"type":    "scheduling_error",
				"message": "No available accounts",
		placeholderplaceholder)
			return
	placeholder
		accountWaitCounted := false
		canWait, waitErr := h.concurrencyHelper.IncrementAccountWaitCount(c.Request.Context(), account.ID, selected.WaitPlan.MaxWaiting)
		if waitErr != nil {
			logger.L().Warn("gateway.web_search.account_wait_counter_increment_failed",
				zap.Int64("account_id", account.ID),
				zap.Error(waitErr),
			)
	placeholder else if !canWait {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": gin.H{
				"type":    "rate_limit_error",
				"message": "Too many pending requests, please retry later",
		placeholderplaceholder)
			return
	placeholder else {
			accountWaitCounted = true
	placeholder
		releaseWait := func() {
			if accountWaitCounted {
				h.concurrencyHelper.DecrementAccountWaitCount(c.Request.Context(), account.ID)
				accountWaitCounted = false
		placeholder
	placeholder
		streamStarted := false
		release, acquireErr := h.concurrencyHelper.AcquireAccountSlotWithWaitTimeout(
			c,
			account.ID,
			selected.WaitPlan.MaxConcurrency,
			selected.WaitPlan.Timeout,
			false,
			&streamStarted,
		)
		releaseWait()
		if acquireErr != nil {
			h.handleConcurrencyError(c, acquireErr, "account", streamStarted)
			return
	placeholder
		accountReleaseFunc = release
placeholder
	if accountReleaseFunc != nil {
		defer accountReleaseFunc()
placeholder

	// Scheduling is 100% the same as other requests:
	// SelectAccountWithLoadAwareness handles load balancing, rate limits, failover, sticky sessions, concurrency, proxies etc.
	// Downstream rate limiting, billing etc. can be wired the same way.

	// Use Grok *native* web search via the selected Grok account + responses API + web_search tool.
	// This ensures results come from Grok's own search (not third-party emulation like Tavily/Brave).
	// Output is normalized to the same unified format for clients/agents/MCP.

	nativeResp, providerName, err := h.doGrokNativeWebSearch(c.Request.Context(), c, account, req.Query, req.MaxResults)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{
			"type":    "web_search_error",
			"message": err.Error(),
	placeholderplaceholder)
		return
placeholder

	userAgent := c.GetHeader("User-Agent")
	clientIP := ip.GetClientIP(c)
	inboundEndpoint := GetInboundEndpoint(c)
	upstreamEndpoint := GetUpstreamEndpoint(c, account.Platform)
	requestPayloadHash := service.HashUsageRequestPayload([]byte(req.Query))
	quotaPlatform := service.QuotaPlatform(c.Request.Context(), apiKey)
	h.submitUsageRecordTask(c.Request.Context(), func(ctx context.Context) {
		if err := h.gatewayService.RecordUsage(ctx, &service.RecordUsageInput{
			Result: &service.ForwardResult{
				RequestID:   "web_search:" + service.HashUsageRequestPayload([]byte(req.Query)),
				Model:       "grok-web-search",
				SearchCount: 1,
				Duration:    0,
		placeholder,
			APIKey:             apiKey,
			User:               apiKey.User,
			Account:            account,
			Subscription:       subscription,
			InboundEndpoint:    inboundEndpoint,
			UpstreamEndpoint:   upstreamEndpoint,
			UserAgent:          userAgent,
			IPAddress:          clientIP,
			RequestPayloadHash: requestPayloadHash,
			APIKeyService:      h.apiKeyService,
			QuotaPlatform:      quotaPlatform,
	placeholder); err != nil {
			logger.L().With(
				zap.String("component", "handler.gateway.web_search"),
				zap.Int64("user_id", apiKey.User.ID),
				zap.Int64("api_key_id", apiKey.ID),
				zap.Int64("account_id", account.ID),
			).Error("gateway.web_search.record_usage_failed", zap.Error(err))
	placeholder
placeholder)

	c.JSON(http.StatusOK, gin.H{
		"query":       req.Query,
		"results":     nativeResp.Results,
		"provider":    providerName,
		"max_results": req.MaxResults,
placeholder)
placeholder

// doGrokNativeWebSearch executes web search using the Grok account's native capability
// by calling the responses endpoint with web_search tool, then normalizes sources to unified format.
func (h *GatewayHandler) doGrokNativeWebSearch(ctx context.Context, c *gin.Context, account *service.Account, query string, maxResults int) (*websearch.SearchResponse, string, error) {
	maxResults = normalizeGrokWebSearchMaxResults(maxResults)

	// Build a minimal responses request that triggers Grok web search tool.
	// Ask for structured metadata because xAI action.sources commonly contains URLs only.
	searchBody := map[string]any{
		"model":   xai.DefaultTextModel,
		"input":   buildGrokWebSearchPrompt(query, maxResults),
		"tools":   []map[string]any{{"type": "web_search"placeholderplaceholder,
		"include": []string{"web_search_call.action.sources"placeholder,
		"store":   false,
		"stream":  false,
placeholder
	bodyBytes, _ := json.Marshal(searchBody)

	respBytes, err := h.gatewayService.DoGrokNativeResponsesJSON(ctx, c, account, bodyBytes)
	if err != nil {
		return nil, "", err
placeholder

	// Extract sources from Grok responses output.
	// Prefer web_search_call.action.sources (standardized), fallback to annotations or text links.
	results := extractGrokWebSearchSources(respBytes, maxResults)

	return &websearch.SearchResponse{
		Results: results,
		Query:   query,
placeholder, "grok-native", nil
placeholder

func normalizeGrokWebSearchMaxResults(maxResults int) int {
	if maxResults <= 0 {
		return defaultGrokWebSearchResults
placeholder
	if maxResults > maxGrokWebSearchResults {
		return maxGrokWebSearchResults
placeholder
	return maxResults
placeholder

func buildGrokWebSearchPrompt(query string, maxResults int) string {
	return fmt.Sprintf(`Search the web for the user query below. Return ONLY valid JSON with this exact shape: {"results":[{"url":"https://...","title":"page title","snippet":"concise factual summary"placeholder]placeholder. Return at most %d unique results. Every URL must be an actual web_search source. Populate a non-empty title and snippet for every result. Do not wrap the JSON in markdown.

User query:
%s`, normalizeGrokWebSearchMaxResults(maxResults), query)
placeholder

// extractGrokWebSearchSources returns model-enriched results only when their URLs
// are present in the actual web_search sources, then falls back to raw sources.
func extractGrokWebSearchSources(body []byte, maxResults int) []websearch.SearchResult {
	if len(body) == 0 || !gjson.ValidBytes(body) {
		return nil
placeholder
	maxResults = normalizeGrokWebSearchMaxResults(maxResults)

	sources := make(map[string]websearch.SearchResult)
	var sourceOrder []string
	addSource := func(rawURL, title, snippet string) {
		key, ok := normalizeGrokWebSearchURL(rawURL)
		if !ok {
			return
	placeholder
		result, exists := sources[key]
		if !exists {
			result.URL = strings.TrimSpace(rawURL)
			sourceOrder = append(sourceOrder, key)
	placeholder
		if result.Title == "" {
			result.Title = usableGrokWebSearchTitle(title, result.URL)
	placeholder
		if result.Snippet == "" {
			result.Snippet = strings.TrimSpace(snippet)
	placeholder
		sources[key] = result
placeholder

	output := gjson.GetBytes(body, "output")
	output.ForEach(func(_, item gjson.Result) bool {
		if item.Get("type").String() == "web_search_call" {
			sources := item.Get("action.sources")
			if sources.IsArray() {
				sources.ForEach(func(_, src gjson.Result) bool {
					addSource(src.Get("url").String(), src.Get("title").String(), src.Get("snippet").String())
					return true
			placeholder)
		placeholder
	placeholder
		if item.Get("type").String() == "message" {
			item.Get("content").ForEach(func(_, part gjson.Result) bool {
				if part.Get("type").String() != "output_text" {
					return true
			placeholder
				part.Get("annotations").ForEach(func(_, ann gjson.Result) bool {
					if ann.Get("type").String() == "url_citation" || ann.Get("type").String() == "web" {
						addSource(ann.Get("url").String(), ann.Get("title").String(), "")
				placeholder
					return true
			placeholder)
				return true
		placeholder)
	placeholder
		return true
placeholder)

	var out []websearch.SearchResult
	seen := make(map[string]bool)
	output.ForEach(func(_, item gjson.Result) bool {
		if item.Get("type").String() != "message" {
			return true
	placeholder
		item.Get("content").ForEach(func(_, part gjson.Result) bool {
			if part.Get("type").String() != "output_text" || len(out) >= maxResults {
				return true
		placeholder
			for _, result := range parseGrokWebSearchStructuredResults(part.Get("text").String()) {
				key, ok := normalizeGrokWebSearchURL(result.URL)
				if !ok || seen[key] {
					continue
			placeholder
				source, allowed := sources[key]
				if !allowed {
					continue
			placeholder
				seen[key] = true
				result.URL = source.URL
				result.Title = usableGrokWebSearchTitle(result.Title, result.URL)
				if result.Title == "" {
					result.Title = source.Title
			placeholder
				result.Snippet = strings.TrimSpace(result.Snippet)
				if result.Snippet == "" {
					result.Snippet = source.Snippet
			placeholder
				out = append(out, result)
				if len(out) >= maxResults {
					break
			placeholder
		placeholder
			return true
	placeholder)
		return len(out) < maxResults
placeholder)

	for _, key := range sourceOrder {
		if len(out) >= maxResults {
			break
	placeholder
		if seen[key] {
			continue
	placeholder
		result := sources[key]
		if result.Title == "" {
			result.Title = grokWebSearchTitleFromURL(result.URL)
	placeholder
		seen[key] = true
		out = append(out, result)
placeholder
	return out
placeholder

func parseGrokWebSearchStructuredResults(text string) []websearch.SearchResult {
	text = strings.TrimSpace(text)
	start := strings.IndexByte(text, '{')
	end := strings.LastIndexByte(text, 'placeholder')
	if start < 0 || end < start {
		return nil
placeholder
	var payload struct {
		Results []websearch.SearchResult `json:"results"`
placeholder
	if err := json.Unmarshal([]byte(text[start:end+1]), &payload); err != nil {
		return nil
placeholder
	return payload.Results
placeholder

func normalizeGrokWebSearchURL(rawURL string) (string, bool) {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return "", false
placeholder
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)
	u.Fragment = ""
	if u.Path == "" {
		u.Path = "/"
placeholder
	return u.String(), true
placeholder

func usableGrokWebSearchTitle(title, rawURL string) string {
	title = strings.TrimSpace(title)
	if title == "" || title == rawURL {
		return ""
placeholder
	if _, err := strconv.Atoi(title); err == nil {
		return ""
placeholder
	return title
placeholder

func grokWebSearchTitleFromURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		return rawURL
placeholder
	return strings.TrimPrefix(strings.ToLower(u.Host), "www.")
placeholder
