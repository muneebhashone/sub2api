package websearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	tavilySearchEndpoint   = "https://api.tavily.com/search"
	tavilyProviderName     = "tavily"
	tavilySearchDepthBasic = "basic"
)

// TavilyProvider implements web search via the Tavily Search API.
type TavilyProvider struct {
	apiKey     string
	httpClient *http.Client
placeholder

// NewTavilyProvider creates a Tavily Search provider.
// The caller is responsible for configuring the http.Client with proxy/timeouts.
func NewTavilyProvider(apiKey string, httpClient *http.Client) *TavilyProvider {
	if httpClient == nil {
		httpClient = http.DefaultClient
placeholder
	return &TavilyProvider{apiKey: apiKey, httpClient: httpClientplaceholder
placeholder

func (t *TavilyProvider) Name() string { return tavilyProviderName placeholder

func (t *TavilyProvider) Search(ctx context.Context, req SearchRequest) (*SearchResponse, error) {
	maxResults := req.MaxResults
	if maxResults <= 0 {
		maxResults = defaultMaxResults
placeholder

	payload := tavilyRequest{
		APIKey:      t.apiKey,
		Query:       req.Query,
		MaxResults:  maxResults,
		SearchDepth: tavilySearchDepthBasic,
placeholder

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("tavily: encode request: %w", err)
placeholder

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, tavilySearchEndpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("tavily: build request: %w", err)
placeholder
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := t.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("tavily: request failed: %w", err)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return nil, fmt.Errorf("tavily: read body: %w", err)
placeholder

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tavily: status %d: %s", resp.StatusCode, truncateBody(body))
placeholder

	var raw tavilyResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("tavily: decode response: %w", err)
placeholder

	results := make([]SearchResult, 0, len(raw.Results))
	for _, r := range raw.Results {
		results = append(results, SearchResult{
			URL:     r.URL,
			Title:   r.Title,
			Snippet: r.Content,
	placeholder)
placeholder

	return &SearchResponse{Results: results, Query: req.Queryplaceholder, nil
placeholder

type tavilyRequest struct {
	APIKey      string `json:"api_key"`
	Query       string `json:"query"`
	MaxResults  int    `json:"max_results"`
	SearchDepth string `json:"search_depth"`
placeholder

type tavilyResponse struct {
	Results []tavilyResult `json:"results"`
placeholder

type tavilyResult struct {
	URL     string  `json:"url"`
	Title   string  `json:"title"`
	Content string  `json:"content"`
	Score   float64 `json:"score"`
placeholder
