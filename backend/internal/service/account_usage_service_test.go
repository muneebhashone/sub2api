package service

import (
	"net/http"
	"testing"
	"time"
)

func TestShouldRefreshOpenAICodexSnapshot(t *testing.T) {
	t.Parallel()

	rateLimitedUntil := time.Now().Add(5 * time.Minute)
	now := time.Now()
	usage := &UsageInfo{
		FiveHour: &UsageProgress{Utilization: 0placeholder,
		SevenDay: &UsageProgress{Utilization: 0placeholder,
placeholder

	if !shouldRefreshOpenAICodexSnapshot(&Account{RateLimitResetAt: &rateLimitedUntilplaceholder, usage, now) {
		t.Fatal("expected rate-limited account to force codex snapshot refresh")
placeholder

	if shouldRefreshOpenAICodexSnapshot(&Account{placeholder, usage, now) {
		t.Fatal("expected complete non-rate-limited usage to skip codex snapshot refresh")
placeholder

	if !shouldRefreshOpenAICodexSnapshot(&Account{placeholder, &UsageInfo{FiveHour: nil, SevenDay: &UsageProgress{placeholderplaceholder, now) {
		t.Fatal("expected missing 5h snapshot to require refresh")
placeholder

	staleAt := now.Add(-(openAIProbeCacheTTL + time.Minute)).Format(time.RFC3339)
	if !shouldRefreshOpenAICodexSnapshot(&Account{
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
		Extra: map[string]any{
			"openai_oauth_responses_websockets_v2_enabled": true,
			"codex_usage_updated_at":                       staleAt,
	placeholder,
placeholder, usage, now) {
		t.Fatal("expected stale ws snapshot to trigger refresh")
placeholder
placeholder

func TestExtractOpenAICodexProbeUpdatesAccepts429WithCodexHeaders(t *testing.T) {
	t.Parallel()

	headers := make(http.Header)
	headers.Set("x-codex-primary-used-percent", "100")
	headers.Set("x-codex-primary-reset-after-seconds", "604800")
	headers.Set("x-codex-primary-window-minutes", "10080")
	headers.Set("x-codex-secondary-used-percent", "100")
	headers.Set("x-codex-secondary-reset-after-seconds", "18000")
	headers.Set("x-codex-secondary-window-minutes", "300")

	updates, err := extractOpenAICodexProbeUpdates(&http.Response{StatusCode: http.StatusTooManyRequests, Header: headersplaceholder)
	if err != nil {
		t.Fatalf("extractOpenAICodexProbeUpdates() error = %v", err)
placeholder
	if len(updates) == 0 {
		t.Fatal("expected codex probe updates from 429 headers")
placeholder
	if got := updates["codex_5h_used_percent"]; got != 100.0 {
		t.Fatalf("codex_5h_used_percent = %v, want 100", got)
placeholder
	if got := updates["codex_7d_used_percent"]; got != 100.0 {
		t.Fatalf("codex_7d_used_percent = %v, want 100", got)
placeholder
placeholder
