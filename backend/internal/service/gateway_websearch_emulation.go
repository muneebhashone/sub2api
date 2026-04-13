package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/websearch"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

// Web search emulation constants
const (
	toolTypeWebSearchPrefix = "web_search"
	toolTypeGoogleSearch    = "google_search"
	toolNameWebSearch       = "web_search"
	toolNameGoogleSearch    = "google_search"
	toolNameWebSearch2025   = "web_search_20250305"

	webSearchDefaultMaxResults = 5
	defaultWebSearchModel      = "claude-sonnet-4-6"
	webSearchMsgIDPrefix       = "msg_ws_"
	webSearchToolUseIDPrefix   = "srvtoolu_ws_"
	tokenEstimateDivisor       = 4

	// featureKeyWebSearchEmulation is the key used in Account.Extra and Channel.FeaturesConfig.
	featureKeyWebSearchEmulation = "web_search_emulation"
)

// webSearchManagerPtr stores *websearch.Manager atomically for concurrent safety.
var webSearchManagerPtr atomic.Pointer[websearch.Manager]

// SetWebSearchManager wires the websearch.Manager into the gateway (goroutine-safe).
func SetWebSearchManager(m *websearch.Manager) {
	webSearchManagerPtr.Store(m)
placeholder

func getWebSearchManager() *websearch.Manager {
	return webSearchManagerPtr.Load()
placeholder

// shouldEmulateWebSearch checks whether a request should be intercepted.
//
// Judgment chain: manager exists → only web_search tool → global enabled → account/channel enabled.
// Account-level mode: "enabled" (force on), "disabled" (force off), "default" (follow channel).
func (s *GatewayService) shouldEmulateWebSearch(ctx context.Context, account *Account, groupID *int64, body []byte) bool {
	if getWebSearchManager() == nil {
		return false
placeholder
	if !isOnlyWebSearchToolInBody(body) {
		return false
placeholder
	if !s.settingService.IsWebSearchEmulationEnabled(ctx) {
		return false
placeholder

	mode := account.GetWebSearchEmulationMode()
	switch mode {
	case WebSearchModeEnabled:
		return true
	case WebSearchModeDisabled:
		return false
	default: // "default" → follow channel config
		if groupID == nil || s.channelService == nil {
			return false
	placeholder
		ch, err := s.channelService.GetChannelForGroup(ctx, *groupID)
		if err != nil || ch == nil {
			return false
	placeholder
		return ch.IsWebSearchEmulationEnabled(account.Platform)
placeholder
placeholder

// isOnlyWebSearchToolInBody checks if the body contains exactly one web_search tool.
func isOnlyWebSearchToolInBody(body []byte) bool {
	tools := gjson.GetBytes(body, "tools")
	if !tools.IsArray() {
		return false
placeholder
	arr := tools.Array()
	if len(arr) != 1 {
		return false
placeholder
	return isWebSearchToolJSON(arr[0])
placeholder

func isWebSearchToolJSON(tool gjson.Result) bool {
	toolType := tool.Get("type").String()
	if strings.HasPrefix(toolType, toolTypeWebSearchPrefix) || toolType == toolTypeGoogleSearch {
		return true
placeholder
	switch tool.Get("name").String() {
	case toolNameWebSearch, toolNameGoogleSearch, toolNameWebSearch2025:
		return true
placeholder
	return false
placeholder

// extractSearchQueryFromBody extracts the last user message text as the search query.
func extractSearchQueryFromBody(body []byte) string {
	messages := gjson.GetBytes(body, "messages")
	if !messages.IsArray() {
		return ""
placeholder
	arr := messages.Array()
	if len(arr) == 0 {
		return ""
placeholder
	lastMsg := arr[len(arr)-1]
	if lastMsg.Get("role").String() != "user" {
		return ""
placeholder
	return extractWebSearchTextFromContent(lastMsg.Get("content"))
placeholder

func extractWebSearchTextFromContent(content gjson.Result) string {
	if content.Type == gjson.String {
		return content.String()
placeholder
	if content.IsArray() {
		for _, block := range content.Array() {
			if block.Get("type").String() == "text" {
				if text := block.Get("text").String(); text != "" {
					return text
			placeholder
		placeholder
	placeholder
placeholder
	return ""
placeholder

// handleWebSearchEmulation intercepts a web-search-only request,
// calls a third-party search API, and constructs an Anthropic-format response.
func (s *GatewayService) handleWebSearchEmulation(
	ctx context.Context, c *gin.Context, account *Account, parsed *ParsedRequest,
) (*ForwardResult, error) {
	startTime := time.Now()

	// Release the serial queue lock immediately — we don't need upstream.
	if parsed.OnUpstreamAccepted != nil {
		parsed.OnUpstreamAccepted()
placeholder

	query := extractSearchQueryFromBody(parsed.Body)
	if query == "" {
		return nil, fmt.Errorf("web search emulation: no query found in messages")
placeholder

	slog.Info("web search emulation: executing search",
		"account_id", account.ID, "account_name", account.Name, "query", query)

	resp, providerName, err := doWebSearch(ctx, account, query)
	if err != nil {
		// Proxy unavailable → trigger account switch via UpstreamFailoverError
		if errors.Is(err, websearch.ErrProxyUnavailable) {
			return nil, &UpstreamFailoverError{
				StatusCode:   http.StatusBadGateway,
				ResponseBody: []byte(err.Error()),
		placeholder
	placeholder
		return nil, err
placeholder

	slog.Info("web search emulation: search completed",
		"provider", providerName, "results_count", len(resp.Results))

	model := parsed.Model
	if model == "" {
		model = defaultWebSearchModel
placeholder

	if parsed.Stream {
		return writeWebSearchStreamResponse(c, query, resp, model, startTime)
placeholder
	return writeWebSearchNonStreamResponse(c, query, resp, model, startTime)
placeholder

func doWebSearch(ctx context.Context, account *Account, query string) (*websearch.SearchResponse, string, error) {
	proxyURL := resolveAccountProxyURL(account)
	mgr := getWebSearchManager()
	if mgr == nil {
		return nil, "", fmt.Errorf("web search emulation: manager not initialized")
placeholder
	resp, providerName, err := mgr.SearchWithBestProvider(ctx, websearch.SearchRequest{
		Query: query, MaxResults: webSearchDefaultMaxResults, ProxyURL: proxyURL,
placeholder)
	if err != nil {
		slog.Error("web search emulation: search failed", "error", err)
		return nil, "", fmt.Errorf("web search emulation: %w", err)
placeholder
	return resp, providerName, nil
placeholder

func resolveAccountProxyURL(account *Account) string {
	if account.ProxyID != nil && account.Proxy != nil {
		return account.Proxy.URL()
placeholder
	return ""
placeholder

// --- SSE streaming response ---

func writeWebSearchStreamResponse(
	c *gin.Context, query string, resp *websearch.SearchResponse, model string, startTime time.Time,
) (*ForwardResult, error) {
	msgID := webSearchMsgIDPrefix + uuid.New().String()
	toolUseID := webSearchToolUseIDPrefix + uuid.New().String()[:16]

	setSSEHeaders(c)
	if err := writeSSEMessageStart(c.Writer, msgID, model); err != nil {
		return nil, fmt.Errorf("web search emulation: SSE write: %w", err)
placeholder
	writeSSEServerToolUse(c.Writer, toolUseID, query, 0)
	writeSSEToolResult(c.Writer, toolUseID, resp.Results, 1)
	textSummary := buildTextSummary(query, resp.Results)
	writeSSETextBlock(c.Writer, textSummary, 2)
	writeSSEMessageEnd(c.Writer, len(textSummary)/tokenEstimateDivisor)
	c.Writer.Flush()

	return &ForwardResult{Model: model, Duration: time.Since(startTime), Usage: ClaudeUsage{placeholderplaceholder, nil
placeholder

func setSSEHeaders(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.WriteHeader(http.StatusOK)
placeholder

func writeSSEMessageStart(w http.ResponseWriter, msgID, model string) error {
	evt := map[string]any{
		"type": "message_start",
		"message": map[string]any{
			"id": msgID, "type": "message", "role": "assistant", "model": model,
			"content": []any{placeholder, "stop_reason": nil, "stop_sequence": nil,
			"usage": map[string]int{"input_tokens": 0, "output_tokens": 0placeholder,
	placeholder,
placeholder
	return flushSSEJSON(w, "message_start", evt)
placeholder

func writeSSEServerToolUse(w http.ResponseWriter, toolUseID, query string, index int) {
	start := map[string]any{
		"type": "content_block_start", "index": index,
		"content_block": map[string]any{
			"type": "server_tool_use", "id": toolUseID,
			"name": toolNameWebSearch, "input": map[string]string{"query": queryplaceholder,
	placeholder,
placeholder
	_ = flushSSEJSON(w, "content_block_start", start)
	_ = flushSSEJSON(w, "content_block_stop", map[string]any{"type": "content_block_stop", "index": indexplaceholder)
placeholder

func writeSSEToolResult(w http.ResponseWriter, toolUseID string, results []websearch.SearchResult, index int) {
	start := map[string]any{
		"type": "content_block_start", "index": index,
		"content_block": map[string]any{
			"type": "web_search_tool_result", "tool_use_id": toolUseID,
			"content": buildSearchResultBlocks(results),
	placeholder,
placeholder
	_ = flushSSEJSON(w, "content_block_start", start)
	_ = flushSSEJSON(w, "content_block_stop", map[string]any{"type": "content_block_stop", "index": indexplaceholder)
placeholder

func writeSSETextBlock(w http.ResponseWriter, text string, index int) {
	_ = flushSSEJSON(w, "content_block_start", map[string]any{
		"type": "content_block_start", "index": index,
		"content_block": map[string]any{"type": "text", "text": ""placeholder,
placeholder)
	_ = flushSSEJSON(w, "content_block_delta", map[string]any{
		"type": "content_block_delta", "index": index,
		"delta": map[string]string{"type": "text_delta", "text": textplaceholder,
placeholder)
	_ = flushSSEJSON(w, "content_block_stop", map[string]any{"type": "content_block_stop", "index": indexplaceholder)
placeholder

func writeSSEMessageEnd(w http.ResponseWriter, outputTokens int) {
	_ = flushSSEJSON(w, "message_delta", map[string]any{
		"type":  "message_delta",
		"delta": map[string]any{"stop_reason": "end_turn", "stop_sequence": nilplaceholder,
		"usage": map[string]int{"output_tokens": outputTokensplaceholder,
placeholder)
	_ = flushSSEJSON(w, "message_stop", map[string]string{"type": "message_stop"placeholder)
placeholder

// flushSSEJSON marshals data to JSON and writes an SSE event. Returns error on marshal failure.
func flushSSEJSON(w http.ResponseWriter, event string, data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		slog.Error("web search emulation: failed to marshal SSE event",
			"event", event, "error", err)
		return err
placeholder
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, b)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
placeholder
	return nil
placeholder

// --- Non-streaming JSON response ---

func writeWebSearchNonStreamResponse(
	c *gin.Context, query string, resp *websearch.SearchResponse, model string, startTime time.Time,
) (*ForwardResult, error) {
	msgID := webSearchMsgIDPrefix + uuid.New().String()
	toolUseID := webSearchToolUseIDPrefix + uuid.New().String()[:16]
	textSummary := buildTextSummary(query, resp.Results)

	msg := map[string]any{
		"id": msgID, "type": "message", "role": "assistant", "model": model,
		"content": []any{
			map[string]any{
				"type": "server_tool_use", "id": toolUseID,
				"name": toolNameWebSearch, "input": map[string]string{"query": queryplaceholder,
		placeholder,
			map[string]any{
				"type": "web_search_tool_result", "tool_use_id": toolUseID,
				"content": buildSearchResultBlocks(resp.Results),
		placeholder,
			map[string]any{"type": "text", "text": textSummaryplaceholder,
	placeholder,
		"stop_reason": "end_turn", "stop_sequence": nil,
		"usage": map[string]int{"input_tokens": 0, "output_tokens": len(textSummary) / tokenEstimateDivisorplaceholder,
placeholder

	body, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("web search emulation: marshal response: %w", err)
placeholder
	c.Data(http.StatusOK, "application/json", body)

	return &ForwardResult{Model: model, Duration: time.Since(startTime), Usage: ClaudeUsage{placeholderplaceholder, nil
placeholder

// --- Helpers ---

func buildSearchResultBlocks(results []websearch.SearchResult) []map[string]string {
	blocks := make([]map[string]string, 0, len(results))
	for _, r := range results {
		block := map[string]string{
			"type":  "web_search_result",
			"url":   r.URL,
			"title": r.Title,
	placeholder
		if r.Snippet != "" {
			block["page_content"] = r.Snippet
	placeholder
		if r.PageAge != "" {
			block["page_age"] = r.PageAge
	placeholder
		blocks = append(blocks, block)
placeholder
	return blocks
placeholder

func buildTextSummary(query string, results []websearch.SearchResult) string {
	if len(results) == 0 {
		return "No search results found for: " + query
placeholder
	var sb strings.Builder
	fmt.Fprintf(&sb, "Here are the search results for \"%s\":\n\n", query)
	for i, r := range results {
		fmt.Fprintf(&sb, "%d. **%s**\n   %s\n   %s\n\n", i+1, r.Title, r.URL, r.Snippet)
placeholder
	return sb.String()
placeholder
