package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// sparkShadowUsageTestRepo is a minimal AccountRepository stub for spark shadow
// usage tests.  GetByID serves both shadow and parent accounts from a map;
// UpdateExtra records the persisted updates for assertion.
type sparkShadowUsageTestRepo struct {
	AccountRepository
	accounts      map[int64]*Account
	updateExtraCh chan map[string]any
placeholder

func (r *sparkShadowUsageTestRepo) GetByID(_ context.Context, id int64) (*Account, error) {
	if acc, ok := r.accounts[id]; ok {
		return acc, nil
placeholder
	return nil, fmt.Errorf("account %d not found", id)
placeholder

func (r *sparkShadowUsageTestRepo) UpdateExtra(_ context.Context, _ int64, updates map[string]any) error {
	if r.updateExtraCh != nil {
		copied := make(map[string]any, len(updates))
		for k, v := range updates {
			copied[k] = v
	placeholder
		r.updateExtraCh <- copied
placeholder
	return nil
placeholder

// TestGetOpenAIUsage_SparkShadow_WritesExtraAndReturnsNonEmptyWindows covers
// two assertions required by Task 3.2:
//
// A) After getOpenAIUsage on a spark shadow account the shadow row's
// Extra["codex_5h_used_percent"] is persisted, and the upstream call carried
// the PARENT account's chatgpt-account-id (not the shadow's empty one).
//
// B) (P1-b regression guard) The UsageInfo RETURNED by the same call has
// non-nil FiveHour AND SevenDay windows — proving that the rebuild happened
// and not just the DB write.
func TestGetOpenAIUsage_SparkShadow_WritesExtraAndReturnsNonEmptyWindows(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	pid := int64(100)
	shadow := &Account{
		ID:              200,
		ParentAccountID: &pid,
		Platform:        PlatformOpenAI,
		Type:            AccountTypeOAuth,
		Status:          StatusActive,
		QuotaDimension:  QuotaDimensionSpark,
placeholder
	parent := &Account{
		ID:       100,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
		Status:   StatusActive,
placeholder
			"chatgpt_account_id": "org-spark-parent",
	placeholder,
placeholder

	// Repo shared by both the OpenAIQuotaService (needs shadow+parent for resolve)
	// and the AccountUsageService (needs UpdateExtra for persist).
	updateExtraCh := make(chan map[string]any, 1)
	repo := &sparkShadowUsageTestRepo{
		accounts:      map[int64]*Account{200: shadow, 100: parentplaceholder,
		updateExtraCh: updateExtraCh,
placeholder

	// Token cache: return a fake token for the parent account key.
	tokenCache := &stubQuotaTokenCache{tokens: map[string]string{
		OpenAITokenCacheKey(parent): "fake-access-token",
placeholderplaceholder
	tokenProvider := NewOpenAITokenProvider(repo, tokenCache, nil)

	// httptest server: records the chatgpt-account-id header and returns a
	// synthetic OpenAIQuotaUsage with codex_bengalfox 5h+7d windows.
	var capturedAccountID string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAccountID = r.Header.Get("chatgpt-account-id")
		w.Header().Set("content-type", "application/json")
		resp := OpenAIQuotaUsage{
			AdditionalRateLimits: []OpenAIAdditionalRateLimit{
				{
					MeteredFeature: "codex_bengalfox",
					RateLimit: &OpenAIRateLimit{
						// Primary window → 5h (18000 s = 300 min)
						PrimaryWindow: &OpenAIRateLimitWindow{
							UsedPercent:        42.5,
							ResetAfterSeconds:  3600,
							LimitWindowSeconds: 18000,
					placeholder,
						// Secondary window → 7d (604800 s = 10080 min)
						SecondaryWindow: &OpenAIRateLimitWindow{
							UsedPercent:        10.0,
							ResetAfterSeconds:  86400,
							LimitWindowSeconds: 604800,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
	placeholder
		_ = json.NewEncoder(w).Encode(resp)
placeholder))
	defer srv.Close()

	quotaService := NewOpenAIQuotaService(repo, nil, tokenProvider, newQuotaRedirectingFactory(srv))
	svc := &AccountUsageService{
		accountRepo:        repo,
		openAIQuotaService: quotaService,
placeholder

	usage, err := svc.getOpenAIUsage(ctx, shadow, true /*force*/)
placeholder

	// Assertion A-1: upstream received the PARENT's chatgpt-account-id.
	require.Equal(t, "org-spark-parent", capturedAccountID,
		"QueryUsage must use parent's chatgpt-account-id for spark shadow accounts")

	// Assertion A-2: shadow Extra was persisted with codex_5h_used_percent.
	select {
	case updates := <-updateExtraCh:
		require.Contains(t, updates, "codex_5h_used_percent",
			"persisted extra must contain codex_5h_used_percent")
		require.InDelta(t, 42.5, updates["codex_5h_used_percent"], 0.01,
			"codex_5h_used_percent must match the upstream value")
	case <-time.After(2 * time.Second):
		t.Fatal("UpdateExtra was not called within timeout — spark shadow persist did not happen")
placeholder

	// Assertion B (P1-b regression guard): returned UsageInfo must have
	// non-nil windows. This FAILS if the code only writes Extra without
	// rebuilding the returned UsageInfo.
	require.NotNil(t, usage.FiveHour,
		"returned UsageInfo.FiveHour must be non-nil (rebuild from merged Extra must happen)")
	require.NotNil(t, usage.SevenDay,
		"returned UsageInfo.SevenDay must be non-nil (rebuild from merged Extra must happen)")
placeholder
