package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOpenAIRateLimitResetCreditDetails_PreservesAvailableCreditOrder(t *testing.T) {
	body := []byte(`{
		"availableCount":"2",
		"credits":[
			{"reset_type":"codex_rate_limits","status":"redeemed","expires_at":"2026-07-01T04:05:06Z"placeholder,
			{"id":"credit-later","reset_type":"codex_rate_limits","status":"available","expires_at":"2026-07-04T04:05:06Z"placeholder,
			{"creditId":"credit-earlier","resetType":"codex_rate_limits","status":"available","expiresAt":"2026-07-03T04:05:06Z"placeholder,
			{"reset_type":"other","status":"available","expires_at":"2026-07-02T04:05:06Z"placeholder
		]
placeholder`)

	details, err := parseOpenAIRateLimitResetCreditDetails(body)
placeholder
	require.NotNil(t, details.AvailableCount)
	require.Equal(t, 2, *details.AvailableCount)
	require.Equal(t, []OpenAIRateLimitResetCreditDetail{
		{ExpiresAt: "2026-07-04T04:05:06Z"placeholder,
		{ExpiresAt: "2026-07-03T04:05:06Z"placeholder,
placeholder, details.Credits)
	require.Equal(t, []openAIAutoResetCreditCandidate{
		{ID: "credit-later", ExpiresAt: "2026-07-04T04:05:06Z"placeholder,
		{ID: "credit-earlier", ExpiresAt: "2026-07-03T04:05:06Z"placeholder,
placeholder, details.AutoResetCandidates)
placeholder

func TestQueryUsageResetCreditCountPrecedence(t *testing.T) {
	tests := []struct {
		name        string
		usageBody   string
		detailBody  string
		wantCount   int
		wantCredits int
		wantNil     bool
placeholder{
		{
			name:       "detail count creates missing usage credits",
			usageBody:  `{placeholder`,
			detailBody: `{"available_count":3,"credits":[{"expires_at":"2026-07-03T04:05:06Z"placeholder]placeholder`,
			wantCount:  3, wantCredits: 1,
	placeholder,
		{
			name:       "explicit detail zero overrides usage and records",
			usageBody:  `{"rate_limit_reset_credits":{"available_count":4placeholderplaceholder`,
			detailBody: `{"available_count":0,"credits":[{"expires_at":"2026-07-03T04:05:06Z"placeholder]placeholder`,
			wantCount:  0, wantCredits: 1,
	placeholder,
		{
			name:       "available records override usage when detail count is absent",
			usageBody:  `{"rate_limit_reset_credits":{"available_count":7placeholderplaceholder`,
			detailBody: `{"credits":[{"expires_at":"2026-07-03T04:05:06Z"placeholder,{"expiresAt":"2026-07-04T04:05:06Z"placeholder]placeholder`,
			wantCount:  2, wantCredits: 2,
	placeholder,
		{
			name:       "empty detail list overrides usage with zero",
			usageBody:  `{"rate_limit_reset_credits":{"available_count":7placeholderplaceholder`,
			detailBody: `{"credits":[]placeholder`,
			wantCount:  0,
	placeholder,
		{
			name:       "fully filtered list overrides usage with zero",
			usageBody:  `{"rate_limit_reset_credits":{"available_count":7placeholderplaceholder`,
			detailBody: `{"credits":[{"reset_type":"codex_rate_limits","status":"redeemed","expires_at":"2026-07-03T04:05:06Z"placeholder,{"reset_type":"other","status":"available","expires_at":"2026-07-04T04:05:06Z"placeholder]placeholder`,
			wantCount:  0,
	placeholder,
		{
			name:       "available records without expiry still count",
			usageBody:  `{"rate_limit_reset_credits":{"available_count":7placeholderplaceholder`,
			detailBody: `{"credits":[{"status":"available"placeholder,{"status":"available","expires_at":"2026-07-04T04:05:06Z"placeholder]placeholder`,
			wantCount:  2, wantCredits: 1,
	placeholder,
		{
			name:        "shape without count or list preserves usage details",
			usageBody:   `{"rate_limit_reset_credits":{"available_count":5,"credits":[{"expires_at":"usage-expiry"placeholder]placeholderplaceholder`,
			detailBody:  `{placeholder`,
			wantCount:   5,
			wantCredits: 1,
	placeholder,
		{
			name:        "valid detail count survives malformed authoritative list",
			usageBody:   `{"rate_limit_reset_credits":{"available_count":7,"credits":[{"expires_at":"usage-expiry"placeholder]placeholderplaceholder`,
			detailBody:  `{"available_count":2,"credits":"malformed"placeholder`,
			wantCount:   2,
			wantCredits: 1,
	placeholder,
		{
			name:       "valid detail count creates quota despite malformed authoritative list",
			usageBody:  `{placeholder`,
			detailBody: `{"available_count":2,"credits":"malformed"placeholder`,
			wantCount:  2,
	placeholder,
		{
			name:       "negative detail count without list preserves usage",
			usageBody:  `{"rate_limit_reset_credits":{"available_count":4placeholderplaceholder`,
			detailBody: `{"available_count":-1placeholder`,
			wantCount:  4,
	placeholder,
		{
			name:       "negative detail count falls back to available records",
			usageBody:  `{"rate_limit_reset_credits":{"available_count":4placeholderplaceholder`,
			detailBody: `{"available_count":-1,"credits":[{"status":"available","expires_at":"2026-07-04T04:05:06Z"placeholder]placeholder`,
			wantCount:  1, wantCredits: 1,
	placeholder,
		{
			name:       "empty object preserves missing usage credits",
			usageBody:  `{placeholder`,
			detailBody: `{placeholder`,
			wantNil:    true,
	placeholder,
		{
			name:       "null body preserves missing usage credits",
			usageBody:  `{placeholder`,
			detailBody: `null`,
			wantNil:    true,
	placeholder,
		{
			name:       "empty body preserves missing usage credits",
			usageBody:  `{placeholder`,
			detailBody: ``,
			wantNil:    true,
	placeholder,
		{
			name:       "null object record is not counted",
			usageBody:  `{"rate_limit_reset_credits":{"available_count":7placeholderplaceholder`,
			detailBody: `{"credits":[null]placeholder`,
			wantCount:  0,
	placeholder,
		{
			name:       "null top level record is not counted",
			usageBody:  `{"rate_limit_reset_credits":{"available_count":7placeholderplaceholder`,
			detailBody: `[null]`,
			wantCount:  0,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := &Account{
				ID:       100,
				Platform: PlatformOpenAI,
				Type:     AccountTypeOAuth,
				Status:   StatusActive,
		placeholder
					"chatgpt_account_id": "org-parent123",
			placeholder,
		placeholder
			repo := &stubQuotaAccountRepo{accounts: map[int64]*Account{100: accountplaceholderplaceholder
			tokenCache := &stubQuotaTokenCache{tokens: map[string]string{
				OpenAITokenCacheKey(account): "fake-token",
		placeholderplaceholder
			tokenProvider := NewOpenAITokenProvider(repo, tokenCache, nil)

			var detailCalls int
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("content-type", "application/json")
				switch r.URL.Path {
				case "/backend-api/wham/usage":
					_, _ = w.Write([]byte(tt.usageBody))
				case "/backend-api/wham/rate-limit-reset-credits":
					detailCalls++
					_, _ = w.Write([]byte(tt.detailBody))
				default:
					http.NotFound(w, r)
			placeholder
		placeholder))
			defer srv.Close()

			svc := NewOpenAIQuotaService(repo, nil, tokenProvider, newQuotaRedirectingFactory(srv))
			usage, err := svc.QueryUsage(context.Background(), 100)
		placeholder
			require.NotNil(t, usage)
			require.Equal(t, 1, detailCalls)
			if tt.wantNil {
				require.Nil(t, usage.RateLimitResetCredits)
				return
		placeholder
			require.NotNil(t, usage.RateLimitResetCredits)
			require.Equal(t, tt.wantCount, usage.RateLimitResetCredits.AvailableCount)
			require.Len(t, usage.RateLimitResetCredits.Credits, tt.wantCredits)
	placeholder)
placeholder
placeholder
