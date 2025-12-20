package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"sub2api/internal/service"
)

type claudeUsageService struct{placeholder

func NewClaudeUsageFetcher() service.ClaudeUsageFetcher {
	return &claudeUsageService{placeholder
placeholder

func (s *claudeUsageService) FetchUsage(ctx context.Context, accessToken, proxyURL string) (*service.ClaudeUsageResponse, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if proxyURL != "" {
		if parsedURL, err := url.Parse(proxyURL); err == nil {
			transport.Proxy = http.ProxyURL(parsedURL)
	placeholder
placeholder

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
placeholder

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.anthropic.com/api/oauth/usage", nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
placeholder

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("anthropic-beta", "oauth-2025-04-20")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
placeholder
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
placeholder

	var usageResp service.ClaudeUsageResponse
	if err := json.NewDecoder(resp.Body).Decode(&usageResp); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
placeholder

	return &usageResp, nil
placeholder
