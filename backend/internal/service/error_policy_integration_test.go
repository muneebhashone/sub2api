//go:build unit

package service

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/antigravity"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Mocks (scoped to this file by naming convention)
// ---------------------------------------------------------------------------

// epFixedUpstream returns a fixed response for every request.
type epFixedUpstream struct {
	statusCode int
	body       string
	calls      int
placeholder

func (u *epFixedUpstream) Do(req *http.Request, proxyURL string, accountID int64, accountConcurrency int) (*http.Response, error) {
	u.calls++
	return &http.Response{
		StatusCode: u.statusCode,
		Header:     http.Header{placeholder,
		Body:       io.NopCloser(strings.NewReader(u.body)),
placeholder, nil
placeholder

func (u *epFixedUpstream) DoWithTLS(req *http.Request, proxyURL string, accountID int64, accountConcurrency int, enableTLSFingerprint bool) (*http.Response, error) {
	return u.Do(req, proxyURL, accountID, accountConcurrency)
placeholder

// epAccountRepo records SetTempUnschedulable / SetError calls.
type epAccountRepo struct {
	mockAccountRepoForGemini
	tempCalls   int
	setErrCalls int
placeholder

func (r *epAccountRepo) SetTempUnschedulable(_ context.Context, _ int64, _ time.Time, _ string) error {
	r.tempCalls++
	return nil
placeholder

func (r *epAccountRepo) SetError(_ context.Context, _ int64, _ string) error {
	r.setErrCalls++
	return nil
placeholder

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func saveAndSetBaseURLs(t *testing.T) {
placeholder
	oldBaseURLs := append([]string(nil), antigravity.BaseURLs...)
	oldAvail := antigravity.DefaultURLAvailability
	antigravity.BaseURLs = []string{"https://ep-test.example"placeholder
	antigravity.DefaultURLAvailability = antigravity.NewURLAvailability(time.Minute)
	t.Cleanup(func() {
		antigravity.BaseURLs = oldBaseURLs
		antigravity.DefaultURLAvailability = oldAvail
placeholder)
placeholder

func newRetryParams(account *Account, upstream HTTPUpstream, handleError func(context.Context, string, *Account, int, http.Header, []byte, string, int64, string, bool) *handleModelRateLimitResult) antigravityRetryLoopParams {
	return antigravityRetryLoopParams{
		ctx:            context.Background(),
		prefix:         "[ep-test]",
		account:        account,
		accessToken:    "token",
		action:         "generateContent",
		body:           []byte(`{"input":"test"placeholder`),
		httpUpstream:   upstream,
		requestedModel: "claude-sonnet-4-5",
		handleError:    handleError,
placeholder
placeholder

// ---------------------------------------------------------------------------
// TestRetryLoop_ErrorPolicy_CustomErrorCodes
// ---------------------------------------------------------------------------

func TestRetryLoop_ErrorPolicy_CustomErrorCodes(t *testing.T) {
	tests := []struct {
		name              string
		upstreamStatus    int
		upstreamBody      string
		customCodes       []any
		expectHandleError int
		expectUpstream    int
		expectStatusCode  int
placeholder{
		{
			name:              "429_in_custom_codes_matched",
			upstreamStatus:    429,
			upstreamBody:      `{"error":"rate limited"placeholder`,
			customCodes:       []any{float64(429)placeholder,
			expectHandleError: 1,
			expectUpstream:    1,
			expectStatusCode:  429,
	placeholder,
		{
			name:              "429_not_in_custom_codes_skipped",
			upstreamStatus:    429,
			upstreamBody:      `{"error":"rate limited"placeholder`,
			customCodes:       []any{float64(500)placeholder,
			expectHandleError: 0,
			expectUpstream:    1,
			expectStatusCode:  429,
	placeholder,
		{
			name:              "500_in_custom_codes_matched",
			upstreamStatus:    500,
			upstreamBody:      `{"error":"internal"placeholder`,
			customCodes:       []any{float64(500)placeholder,
			expectHandleError: 1,
			expectUpstream:    1,
			expectStatusCode:  500,
	placeholder,
		{
			name:              "500_not_in_custom_codes_skipped",
			upstreamStatus:    500,
			upstreamBody:      `{"error":"internal"placeholder`,
			customCodes:       []any{float64(429)placeholder,
			expectHandleError: 0,
			expectUpstream:    1,
			expectStatusCode:  500,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saveAndSetBaseURLs(t)

			upstream := &epFixedUpstream{statusCode: tt.upstreamStatus, body: tt.upstreamBodyplaceholder
			repo := &epAccountRepo{placeholder
			rlSvc := NewRateLimitService(repo, nil, &config.Config{placeholder, nil, nil)

			account := &Account{
				ID:          100,
				Type:        AccountTypeAPIKey,
				Platform:    PlatformAntigravity,
				Schedulable: true,
				Status:      StatusActive,
				Concurrency: 1,
		placeholder
					"custom_error_codes_enabled": true,
					"custom_error_codes":         tt.customCodes,
			placeholder,
		placeholder

			svc := &AntigravityGatewayService{rateLimitService: rlSvcplaceholder

			var handleErrorCount int
			p := newRetryParams(account, upstream, func(_ context.Context, _ string, _ *Account, _ int, _ http.Header, _ []byte, _ string, _ int64, _ string, _ bool) *handleModelRateLimitResult {
				handleErrorCount++
				return nil
		placeholder)

			result, err := svc.antigravityRetryLoop(p)

		placeholder
			require.NotNil(t, result)
			require.NotNil(t, result.resp)
			defer func() { _ = result.resp.Body.Close() placeholder()

			require.Equal(t, tt.expectStatusCode, result.resp.StatusCode)
			require.Equal(t, tt.expectHandleError, handleErrorCount, "handleError call count")
			require.Equal(t, tt.expectUpstream, upstream.calls, "upstream call count")
	placeholder)
placeholder
placeholder

// ---------------------------------------------------------------------------
// TestRetryLoop_ErrorPolicy_TempUnschedulable
// ---------------------------------------------------------------------------

func TestRetryLoop_ErrorPolicy_TempUnschedulable(t *testing.T) {
	tempRulesAccount := func(rules []any) *Account {
	placeholder
			ID:          200,
			Type:        AccountTypeOAuth,
			Platform:    PlatformAntigravity,
			Schedulable: true,
			Status:      StatusActive,
			Concurrency: 1,
	placeholder
				"temp_unschedulable_enabled": true,
				"temp_unschedulable_rules":   rules,
		placeholder,
	placeholder
placeholder

	overloadedRule := map[string]any{
		"error_code":       float64(503),
		"keywords":         []any{"overloaded"placeholder,
		"duration_minutes": float64(10),
placeholder

	rateLimitRule := map[string]any{
		"error_code":       float64(429),
		"keywords":         []any{"rate limited keyword"placeholder,
		"duration_minutes": float64(5),
placeholder

	t.Run("503_overloaded_matches_rule", func(t *testing.T) {
		saveAndSetBaseURLs(t)

		upstream := &epFixedUpstream{statusCode: 503, body: `overloaded`placeholder
		repo := &epAccountRepo{placeholder
		rlSvc := NewRateLimitService(repo, nil, &config.Config{placeholder, nil, nil)
		svc := &AntigravityGatewayService{rateLimitService: rlSvcplaceholder

		account := tempRulesAccount([]any{overloadedRuleplaceholder)
		p := newRetryParams(account, upstream, func(_ context.Context, _ string, _ *Account, _ int, _ http.Header, _ []byte, _ string, _ int64, _ string, _ bool) *handleModelRateLimitResult {
			t.Error("handleError should not be called for temp unschedulable")
			return nil
	placeholder)

		result, err := svc.antigravityRetryLoop(p)

		require.Nil(t, result)
		var switchErr *AntigravityAccountSwitchError
		require.ErrorAs(t, err, &switchErr)
		require.Equal(t, account.ID, switchErr.OriginalAccountID)
		require.Equal(t, 1, upstream.calls, "should not retry")
placeholder)

	t.Run("429_rate_limited_keyword_matches_rule", func(t *testing.T) {
		saveAndSetBaseURLs(t)

		upstream := &epFixedUpstream{statusCode: 429, body: `rate limited keyword`placeholder
		repo := &epAccountRepo{placeholder
		rlSvc := NewRateLimitService(repo, nil, &config.Config{placeholder, nil, nil)
		svc := &AntigravityGatewayService{rateLimitService: rlSvcplaceholder

		account := tempRulesAccount([]any{rateLimitRuleplaceholder)
		p := newRetryParams(account, upstream, func(_ context.Context, _ string, _ *Account, _ int, _ http.Header, _ []byte, _ string, _ int64, _ string, _ bool) *handleModelRateLimitResult {
			t.Error("handleError should not be called for temp unschedulable")
			return nil
	placeholder)

		result, err := svc.antigravityRetryLoop(p)

		require.Nil(t, result)
		var switchErr *AntigravityAccountSwitchError
		require.ErrorAs(t, err, &switchErr)
		require.Equal(t, account.ID, switchErr.OriginalAccountID)
		require.Equal(t, 1, upstream.calls, "should not retry")
placeholder)

	t.Run("503_body_no_match_continues_default_retry", func(t *testing.T) {
		saveAndSetBaseURLs(t)

		upstream := &epFixedUpstream{statusCode: 503, body: `random`placeholder
		repo := &epAccountRepo{placeholder
		rlSvc := NewRateLimitService(repo, nil, &config.Config{placeholder, nil, nil)
		svc := &AntigravityGatewayService{rateLimitService: rlSvcplaceholder

		account := tempRulesAccount([]any{overloadedRuleplaceholder)

		// Use a short-lived context: the backoff sleep (~1s) will be
		// interrupted, proving the code entered the default retry path
		// instead of breaking early via error policy.
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		p := newRetryParams(account, upstream, func(_ context.Context, _ string, _ *Account, _ int, _ http.Header, _ []byte, _ string, _ int64, _ string, _ bool) *handleModelRateLimitResult {
			return nil
	placeholder)
		p.ctx = ctx

		result, err := svc.antigravityRetryLoop(p)

		// Context cancellation during backoff proves default retry was entered
		require.Nil(t, result)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.GreaterOrEqual(t, upstream.calls, 1, "should have called upstream at least once")
placeholder)
placeholder

// ---------------------------------------------------------------------------
// TestRetryLoop_ErrorPolicy_NilRateLimitService
// ---------------------------------------------------------------------------

func TestRetryLoop_ErrorPolicy_NilRateLimitService(t *testing.T) {
	saveAndSetBaseURLs(t)

	upstream := &epFixedUpstream{statusCode: 429, body: `{"error":"rate limited"placeholder`placeholder
	// rateLimitService is nil — must not panic
	svc := &AntigravityGatewayService{rateLimitService: nilplaceholder

	account := &Account{
		ID:          300,
		Type:        AccountTypeOAuth,
		Platform:    PlatformAntigravity,
		Schedulable: true,
		Status:      StatusActive,
		Concurrency: 1,
placeholder

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	p := newRetryParams(account, upstream, func(_ context.Context, _ string, _ *Account, _ int, _ http.Header, _ []byte, _ string, _ int64, _ string, _ bool) *handleModelRateLimitResult {
		return nil
placeholder)
	p.ctx = ctx

	// Should not panic; enters the default retry path (eventually times out)
	result, err := svc.antigravityRetryLoop(p)

	require.Nil(t, result)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.GreaterOrEqual(t, upstream.calls, 1)
placeholder

// ---------------------------------------------------------------------------
// TestRetryLoop_ErrorPolicy_NoPolicy_OriginalBehavior
// ---------------------------------------------------------------------------

func TestRetryLoop_ErrorPolicy_NoPolicy_OriginalBehavior(t *testing.T) {
	saveAndSetBaseURLs(t)

	upstream := &epFixedUpstream{statusCode: 429, body: `{"error":"rate limited"placeholder`placeholder
	repo := &epAccountRepo{placeholder
	rlSvc := NewRateLimitService(repo, nil, &config.Config{placeholder, nil, nil)
	svc := &AntigravityGatewayService{rateLimitService: rlSvcplaceholder

	// Plain OAuth account with no error policy configured
	account := &Account{
		ID:          400,
		Type:        AccountTypeOAuth,
		Platform:    PlatformAntigravity,
		Schedulable: true,
		Status:      StatusActive,
		Concurrency: 1,
placeholder

	var handleErrorCount int
	p := newRetryParams(account, upstream, func(_ context.Context, _ string, _ *Account, _ int, _ http.Header, _ []byte, _ string, _ int64, _ string, _ bool) *handleModelRateLimitResult {
		handleErrorCount++
		return nil
placeholder)

	result, err := svc.antigravityRetryLoop(p)

placeholder
	require.NotNil(t, result)
	require.NotNil(t, result.resp)
	defer func() { _ = result.resp.Body.Close() placeholder()

	require.Equal(t, http.StatusTooManyRequests, result.resp.StatusCode)
	require.Equal(t, antigravityMaxRetries, upstream.calls, "should exhaust all retries")
	require.Equal(t, 1, handleErrorCount, "handleError should be called once after retries exhausted")
placeholder
