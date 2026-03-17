package service

import (
	"context"
	"net/http"
	"testing"
	"time"
)

type accountUsageCodexProbeRepo struct {
	stubOpenAIAccountRepo
	updateExtraCh chan map[string]any
	rateLimitCh   chan time.Time
placeholder

func (r *accountUsageCodexProbeRepo) UpdateExtra(_ context.Context, _ int64, updates map[string]any) error {
	if r.updateExtraCh != nil {
		copied := make(map[string]any, len(updates))
		for k, v := range updates {
			copied[k] = v
	placeholder
		r.updateExtraCh <- copied
placeholder
	return nil
placeholder

func (r *accountUsageCodexProbeRepo) SetRateLimited(_ context.Context, _ int64, resetAt time.Time) error {
	if r.rateLimitCh != nil {
		r.rateLimitCh <- resetAt
placeholder
	return nil
placeholder

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

func TestExtractOpenAICodexProbeSnapshotAccepts429WithResetAt(t *testing.T) {
	t.Parallel()

	headers := make(http.Header)
	headers.Set("x-codex-primary-used-percent", "100")
	headers.Set("x-codex-primary-reset-after-seconds", "604800")
	headers.Set("x-codex-primary-window-minutes", "10080")
	headers.Set("x-codex-secondary-used-percent", "100")
	headers.Set("x-codex-secondary-reset-after-seconds", "18000")
	headers.Set("x-codex-secondary-window-minutes", "300")

	updates, resetAt, err := extractOpenAICodexProbeSnapshot(&http.Response{StatusCode: http.StatusTooManyRequests, Header: headersplaceholder)
	if err != nil {
		t.Fatalf("extractOpenAICodexProbeSnapshot() error = %v", err)
placeholder
	if len(updates) == 0 {
		t.Fatal("expected codex probe updates from 429 headers")
placeholder
	if resetAt == nil {
		t.Fatal("expected resetAt from exhausted codex headers")
placeholder
placeholder

func TestAccountUsageService_PersistOpenAICodexProbeSnapshotSetsRateLimit(t *testing.T) {
	t.Parallel()

	repo := &accountUsageCodexProbeRepo{
		updateExtraCh: make(chan map[string]any, 1),
		rateLimitCh:   make(chan time.Time, 1),
placeholder
	svc := &AccountUsageService{accountRepo: repoplaceholder
	resetAt := time.Now().Add(2 * time.Hour).UTC().Truncate(time.Second)

	svc.persistOpenAICodexProbeSnapshot(321, map[string]any{
		"codex_7d_used_percent": 100.0,
		"codex_7d_reset_at":     resetAt.Format(time.RFC3339),
placeholder, &resetAt)

	select {
	case updates := <-repo.updateExtraCh:
		if got := updates["codex_7d_used_percent"]; got != 100.0 {
			t.Fatalf("codex_7d_used_percent = %v, want 100", got)
	placeholder
	case <-time.After(2 * time.Second):
		t.Fatal("waiting for codex probe extra persistence timed out")
placeholder

	select {
	case got := <-repo.rateLimitCh:
		if got.Before(resetAt.Add(-time.Second)) || got.After(resetAt.Add(time.Second)) {
			t.Fatalf("rate limit resetAt = %v, want around %v", got, resetAt)
	placeholder
	case <-time.After(2 * time.Second):
		t.Fatal("waiting for codex probe rate limit persistence timed out")
placeholder
placeholder

func TestBuildCodexUsageProgressFromExtra_ZerosExpiredWindow(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 3, 16, 12, 0, 0, 0, time.UTC)

	t.Run("expired 5h window zeroes utilization", func(t *testing.T) {
		extra := map[string]any{
			"codex_5h_used_percent": 42.0,
			"codex_5h_reset_at":     "2026-03-16T10:00:00Z", // 2h ago
	placeholder
		progress := buildCodexUsageProgressFromExtra(extra, "5h", now)
		if progress == nil {
			t.Fatal("expected non-nil progress")
	placeholder
		if progress.Utilization != 0 {
			t.Fatalf("expected Utilization=0 for expired window, got %v", progress.Utilization)
	placeholder
		if progress.RemainingSeconds != 0 {
			t.Fatalf("expected RemainingSeconds=0, got %v", progress.RemainingSeconds)
	placeholder
placeholder)

	t.Run("active 5h window keeps utilization", func(t *testing.T) {
		resetAt := now.Add(2 * time.Hour).Format(time.RFC3339)
		extra := map[string]any{
			"codex_5h_used_percent": 42.0,
			"codex_5h_reset_at":     resetAt,
	placeholder
		progress := buildCodexUsageProgressFromExtra(extra, "5h", now)
		if progress == nil {
			t.Fatal("expected non-nil progress")
	placeholder
		if progress.Utilization != 42.0 {
			t.Fatalf("expected Utilization=42, got %v", progress.Utilization)
	placeholder
placeholder)

	t.Run("expired 7d window zeroes utilization", func(t *testing.T) {
		extra := map[string]any{
			"codex_7d_used_percent": 88.0,
			"codex_7d_reset_at":     "2026-03-15T00:00:00Z", // yesterday
	placeholder
		progress := buildCodexUsageProgressFromExtra(extra, "7d", now)
		if progress == nil {
			t.Fatal("expected non-nil progress")
	placeholder
		if progress.Utilization != 0 {
			t.Fatalf("expected Utilization=0 for expired 7d window, got %v", progress.Utilization)
	placeholder
placeholder)
placeholder
