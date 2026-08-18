//go:build unit

package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// TestShouldFailoverGeminiUpstreamError — verifies the failover decision
// for the ErrorPolicyNone path (original logic preserved).
// ---------------------------------------------------------------------------

func TestShouldFailoverGeminiUpstreamError(t *testing.T) {
	svc := &GeminiMessagesCompatService{placeholder

	tests := []struct {
		name       string
		statusCode int
		expected   bool
placeholder{
		{"401_failover", 401, trueplaceholder,
		{"403_failover", 403, trueplaceholder,
		{"429_failover", 429, trueplaceholder,
		{"529_failover", 529, trueplaceholder,
		{"500_failover", 500, trueplaceholder,
		{"502_failover", 502, trueplaceholder,
		{"503_failover", 503, trueplaceholder,
		{"400_no_failover", 400, falseplaceholder,
		{"404_no_failover", 404, falseplaceholder,
		{"422_no_failover", 422, falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.shouldFailoverGeminiUpstreamError(tt.statusCode)
			require.Equal(t, tt.expected, got)
	placeholder)
placeholder
placeholder

// ---------------------------------------------------------------------------
// TestCheckErrorPolicy_GeminiAccounts — verifies CheckErrorPolicy works
// correctly for Gemini platform accounts (API Key type).
// ---------------------------------------------------------------------------

func TestCheckErrorPolicy_GeminiAccounts(t *testing.T) {
	tests := []struct {
		name       string
		account    *Account
		statusCode int
		body       []byte
		expected   ErrorPolicyResult
placeholder{
		{
			name: "gemini_apikey_custom_codes_hit",
			account: &Account{
				ID:       100,
				Type:     AccountTypeAPIKey,
		placeholder
		placeholder
					"custom_error_codes_enabled": true,
					"custom_error_codes":         []any{float64(429), float64(500)placeholder,
			placeholder,
		placeholder,
			statusCode: 429,
			body:       []byte(`{"error":"rate limited"placeholder`),
			expected:   ErrorPolicyMatched,
	placeholder,
		{
			name: "gemini_apikey_custom_codes_miss",
			account: &Account{
				ID:       101,
				Type:     AccountTypeAPIKey,
		placeholder
		placeholder
					"custom_error_codes_enabled": true,
					"custom_error_codes":         []any{float64(429)placeholder,
			placeholder,
		placeholder,
			statusCode: 500,
			body:       []byte(`{"error":"internal"placeholder`),
			expected:   ErrorPolicySkipped,
	placeholder,
		{
			name: "gemini_apikey_no_custom_codes_returns_none",
			account: &Account{
				ID:       102,
				Type:     AccountTypeAPIKey,
		placeholder
		placeholder,
			statusCode: 500,
			body:       []byte(`{"error":"internal"placeholder`),
			expected:   ErrorPolicyNone,
	placeholder,
		{
			name: "gemini_apikey_temp_unschedulable_hit",
			account: &Account{
				ID:       103,
				Type:     AccountTypeAPIKey,
		placeholder
		placeholder
					"temp_unschedulable_enabled": true,
					"temp_unschedulable_rules": []any{
						map[string]any{
							"error_code":       float64(503),
							"keywords":         []any{"overloaded"placeholder,
							"duration_minutes": float64(10),
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			statusCode: 503,
			body:       []byte(`overloaded service`),
			expected:   ErrorPolicyTempUnscheduled,
	placeholder,
		{
			name: "gemini_apikey_temp_unschedulable_401_second_hit_returns_none",
			account: &Account{
				ID:                      105,
				Type:                    AccountTypeAPIKey,
				Platform:                PlatformGemini,
				TempUnschedulableReason: `{"status_code":401,"until_unix":1735689600placeholder`,
		placeholder
					"temp_unschedulable_enabled": true,
					"temp_unschedulable_rules": []any{
						map[string]any{
							"error_code":       float64(401),
							"keywords":         []any{"unauthorized"placeholder,
							"duration_minutes": float64(10),
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			statusCode: 401,
			body:       []byte(`unauthorized`),
			expected:   ErrorPolicyNone,
	placeholder,
		{
			name: "gemini_custom_codes_override_temp_unschedulable",
			account: &Account{
				ID:       104,
				Type:     AccountTypeAPIKey,
		placeholder
		placeholder
					"custom_error_codes_enabled": true,
					"custom_error_codes":         []any{float64(503)placeholder,
					"temp_unschedulable_enabled": true,
					"temp_unschedulable_rules": []any{
						map[string]any{
							"error_code":       float64(503),
							"keywords":         []any{"overloaded"placeholder,
							"duration_minutes": float64(10),
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			statusCode: 503,
			body:       []byte(`overloaded`),
			expected:   ErrorPolicyMatched, // custom codes take precedence
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &errorPolicyRepoStub{placeholder
			svc := NewRateLimitService(repo, nil, &config.Config{placeholder, nil, nil)

			result := svc.CheckErrorPolicy(context.Background(), tt.account, tt.statusCode, tt.body)
			require.Equal(t, tt.expected, result)
	placeholder)
placeholder
placeholder

// ---------------------------------------------------------------------------
// TestGeminiErrorPolicyIntegration — verifies the Gemini error handling
// paths produce the correct behavior for each ErrorPolicyResult.
//
// These tests simulate the inline error policy switch in handleClaudeCompat
// and forwardNativeGemini by calling the same methods in the same order.
// ---------------------------------------------------------------------------

func TestGeminiErrorPolicyIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name                 string
		account              *Account
		statusCode           int
		respBody             []byte
		expectFailover       bool // expect UpstreamFailoverError
		expectHandleError    bool // expect handleGeminiUpstreamError to be called
		expectShouldFailover bool // for None path, whether shouldFailover triggers
		expectModelScope     string
placeholder{
		{
			name: "custom_codes_matched_429_failover",
			account: &Account{
				ID:       200,
				Type:     AccountTypeAPIKey,
		placeholder
		placeholder
					"custom_error_codes_enabled": true,
					"custom_error_codes":         []any{float64(429)placeholder,
			placeholder,
		placeholder,
			statusCode:        429,
			respBody:          []byte(`{"error":"rate limited"placeholder`),
			expectFailover:    true,
			expectHandleError: true,
	placeholder,
		{
			name: "custom_codes_skipped_500_failover",
			account: &Account{
				ID:       201,
				Type:     AccountTypeAPIKey,
		placeholder
		placeholder
					"custom_error_codes_enabled": true,
					"custom_error_codes":         []any{float64(429)placeholder,
			placeholder,
		placeholder,
			statusCode:        500,
			respBody:          []byte(`{"error":"internal"placeholder`),
			expectFailover:    true,
			expectHandleError: false,
	placeholder,
		{
			name: "custom_codes_skipped_400_no_failover",
			account: &Account{
				ID:       205,
				Type:     AccountTypeAPIKey,
		placeholder
		placeholder
					"custom_error_codes_enabled": true,
					"custom_error_codes":         []any{float64(429)placeholder,
			placeholder,
		placeholder,
			statusCode:        400,
			respBody:          []byte(`{"error":"bad request"placeholder`),
			expectFailover:    false,
			expectHandleError: false,
	placeholder,
		{
			name: "temp_unschedulable_matched_failover",
			account: &Account{
				ID:       202,
				Type:     AccountTypeAPIKey,
		placeholder
		placeholder
					"temp_unschedulable_enabled": true,
					"temp_unschedulable_rules": []any{
						map[string]any{
							"error_code":       float64(503),
							"keywords":         []any{"overloaded"placeholder,
							"duration_minutes": float64(10),
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			statusCode:        503,
			respBody:          []byte(`overloaded`),
			expectFailover:    true,
			expectHandleError: false,
			expectModelScope:  "gemini-2.5-pro",
	placeholder,
		{
			name: "no_policy_429_failover_via_shouldFailover",
			account: &Account{
				ID:       203,
				Type:     AccountTypeAPIKey,
		placeholder
		placeholder,
			statusCode:           429,
			respBody:             []byte(`{"error":"rate limited"placeholder`),
			expectFailover:       true,
			expectHandleError:    true,
			expectShouldFailover: true,
	placeholder,
		{
			name: "no_policy_400_no_failover",
			account: &Account{
				ID:       204,
				Type:     AccountTypeAPIKey,
		placeholder
		placeholder,
			statusCode:        400,
			respBody:          []byte(`{"error":"bad request"placeholder`),
			expectFailover:    false,
			expectHandleError: true,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &geminiErrorPolicyRepo{placeholder
			rlSvc := NewRateLimitService(repo, nil, &config.Config{placeholder, nil, nil)
			svc := &GeminiMessagesCompatService{
				accountRepo:      repo,
				rateLimitService: rlSvc,
		placeholder

			writer := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(writer)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

			// Simulate the Claude compat error handling path (same logic as native).
			// This mirrors the inline switch in handleClaudeCompat.
			var handleErrorCalled bool
			var gotFailover bool

			ctx := context.Background()
			statusCode := tt.statusCode
			respBody := tt.respBody
			account := tt.account
			headers := http.Header{placeholder

			if svc.rateLimitService != nil {
				policy := svc.rateLimitService.CheckErrorPolicy(ctx, account, statusCode, respBody, "gemini-2.5-pro")
				switch policy {
				case ErrorPolicySkipped:
					// Skipped → 不标记账号状态；可 failover 的状态码仍换号
					handleErrorCalled = false
					gotFailover = svc.skippedErrorPolicyFailoverError(c, account, statusCode, respBody, "req-test") != nil
					goto verify
				case ErrorPolicyMatched:
					svc.handleGeminiUpstreamError(ctx, account, statusCode, headers, respBody)
					handleErrorCalled = true
					gotFailover = true
					goto verify
				case ErrorPolicyTempUnscheduled:
					handleErrorCalled = false
					gotFailover = true
					goto verify
			placeholder
		placeholder

			// ErrorPolicyNone → original logic
			svc.handleGeminiUpstreamError(ctx, account, statusCode, headers, respBody)
			handleErrorCalled = true
			if svc.shouldFailoverGeminiUpstreamError(statusCode) {
				gotFailover = true
		placeholder

		verify:
			require.Equal(t, tt.expectFailover, gotFailover, "failover mismatch")
			require.Equal(t, tt.expectHandleError, handleErrorCalled, "handleGeminiUpstreamError call mismatch")
			if tt.expectModelScope != "" {
				require.Equal(t, 1, repo.setModelRateLimitedCalls)
				require.Equal(t, tt.expectModelScope, repo.lastModelScope)
				require.Zero(t, repo.setTempCalls)
				require.Zero(t, repo.setRateLimitedCalls, "model temp rule must not be widened into an account rate limit")
		placeholder

			if tt.expectShouldFailover {
				require.True(t, svc.shouldFailoverGeminiUpstreamError(statusCode),
					"shouldFailoverGeminiUpstreamError should return true for status %d", statusCode)
		placeholder
	placeholder)
placeholder
placeholder

// ---------------------------------------------------------------------------
// TestSkippedErrorPolicyFailoverError — ErrorPolicySkipped（池模式、或自定义
// 错误码未命中）不豁免换号：可 failover 的状态码返回 UpstreamFailoverError，
// 仅池模式账号可携带同账号重试标记。
// ---------------------------------------------------------------------------

func TestSkippedErrorPolicyFailoverError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &GeminiMessagesCompatService{placeholder

	poolAccount := func(extra map[string]any) *Account {
		creds := map[string]any{"pool_mode": trueplaceholder
		for k, v := range extra {
			creds[k] = v
	placeholder
	placeholderID: 300, Type: AccountTypeAPIKey, Platform: PlatformGemini, Credentials: credsplaceholder
placeholder
	customCodesAccount := &Account{
		ID: 301, Type: AccountTypeAPIKey, Platform: PlatformGemini,
placeholder
			"custom_error_codes_enabled": true,
			"custom_error_codes":         []any{float64(429)placeholder,
	placeholder,
placeholder

	tests := []struct {
		name              string
		account           *Account
		statusCode        int
		expectFailover    bool
		expectSameAccount bool
placeholder{
		{"pool_500_failover_no_same_account_retry", poolAccount(nil), 500, true, falseplaceholder,
		{"pool_429_failover_with_same_account_retry", poolAccount(nil), 429, true, trueplaceholder,
		{"pool_custom_retry_codes_500", poolAccount(map[string]any{
			"pool_mode_retry_status_codes": []any{float64(500)placeholder,
	placeholder), 500, true, trueplaceholder,
		{"pool_400_not_failover_worthy", poolAccount(nil), 400, false, falseplaceholder,
		{"custom_codes_miss_500_failover_no_same_account_retry", customCodesAccount, 500, true, falseplaceholder,
		{"custom_codes_miss_400_not_failover_worthy", customCodesAccount, 400, false, falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(writer)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

			body := []byte(`{"error":{"code":"bad_response_status_code","message":"openai_error"placeholderplaceholder`)
			failoverErr := svc.skippedErrorPolicyFailoverError(c, tt.account, tt.statusCode, body, "req-1")

			if !tt.expectFailover {
				require.Nil(t, failoverErr)
				return
		placeholder
			require.NotNil(t, failoverErr)
			require.Equal(t, tt.statusCode, failoverErr.StatusCode)
			require.Equal(t, body, failoverErr.ResponseBody)
			require.Equal(t, tt.expectSameAccount, failoverErr.RetryableOnSameAccount)
			require.True(t, failoverErr.ShouldRetryNextAccount())
	placeholder)
placeholder
placeholder

// ---------------------------------------------------------------------------
// TestGeminiErrorPolicy_NilRateLimitService — verifies nil safety
// ---------------------------------------------------------------------------

func TestGeminiErrorPolicy_NilRateLimitService(t *testing.T) {
	svc := &GeminiMessagesCompatService{
		rateLimitService: nil,
placeholder

	// When rateLimitService is nil, error policy is skipped → falls through to
	// shouldFailoverGeminiUpstreamError (original logic).
	// Verify this doesn't panic and follows expected behavior.

	ctx := context.Background()
	account := &Account{
		ID:       300,
		Type:     AccountTypeAPIKey,
placeholder
placeholder
			"custom_error_codes_enabled": true,
			"custom_error_codes":         []any{float64(429)placeholder,
	placeholder,
placeholder

	// The nil check should prevent CheckErrorPolicy from being called
	if svc.rateLimitService != nil {
		t.Fatal("rateLimitService should be nil for this test")
placeholder

	// shouldFailoverGeminiUpstreamError still works
	require.True(t, svc.shouldFailoverGeminiUpstreamError(429))
	require.False(t, svc.shouldFailoverGeminiUpstreamError(400))

	// handleGeminiUpstreamError should not panic with nil rateLimitService
	require.NotPanics(t, func() {
		svc.handleGeminiUpstreamError(ctx, account, 500, http.Header{placeholder, []byte(`error`))
placeholder)
placeholder

// ---------------------------------------------------------------------------
// geminiErrorPolicyRepo — minimal AccountRepository stub for Gemini error
// policy tests. Embeds mockAccountRepoForGemini and adds tracking.
// ---------------------------------------------------------------------------

func TestHandleGeminiUpstreamError_GoogleOneCapacityExhaustedUsesTierCooldown(t *testing.T) {
	repo := &rateLimit429AccountRepoStub{placeholder
	quotaSvc := NewGeminiQuotaService(&config.Config{placeholder, nil)
	rlSvc := NewRateLimitService(repo, nil, &config.Config{placeholder, quotaSvc, nil)
	svc := &GeminiMessagesCompatService{
		accountRepo:      repo,
		rateLimitService: rlSvc,
placeholder

	account := &Account{
		ID:       511,
placeholder
		Type:     AccountTypeOAuth,
placeholder
			"oauth_type": "google_one",
			"tier_id":    "google_ai_pro",
	placeholder,
placeholder
	body := []byte(`{"error":{"code":429,"details":[{"@type":"type.googleapis.com/google.rpc.ErrorInfo","domain":"cloudcode-pa.googleapis.com","metadata":{"model":"gemini-3.1-pro-preview"placeholder,"reason":"MODEL_CAPACITY_EXHAUSTED"placeholder],"message":"No capacity available for model gemini-3.1-pro-preview on the server","status":"RESOURCE_EXHAUSTED"placeholderplaceholder`)

	before := time.Now()
	svc.handleGeminiUpstreamError(context.Background(), account, http.StatusTooManyRequests, http.Header{placeholder, body)
	after := time.Now()

	require.Equal(t, 1, repo.rateLimitCalls)
	require.Equal(t, int64(511), repo.lastRateLimitID)
	require.WithinDuration(t, before.Add(5*time.Minute), repo.lastRateLimitReset, 2*time.Second)
	require.True(t, repo.lastRateLimitReset.After(before))
	require.True(t, repo.lastRateLimitReset.Before(after.Add(5*time.Minute).Add(2*time.Second)))
placeholder

// ---------------------------------------------------------------------------
// TestHandleGeminiUpstreamError_PoolMode429 — 池模式账号的 429 不写账号级限流。
//
// 429 的标记点在重试循环内（handleClaudeCompat / forwardNativeGemini /
// chat completions 三条路径），先于 CheckErrorPolicy 执行，池模式豁免只能落在
// handleGeminiUpstreamError 自身；否则一次上游 429 会把账号锁到 PST 午夜，
// 即便重试已经成功返回客户端。
// ---------------------------------------------------------------------------

func TestHandleGeminiUpstreamError_PoolMode429(t *testing.T) {
	// 中转上游的真实 429 文案：不含 "per day"，也没有 quotaResetDelay，
	// 解析失败后 apikey 账号会落到 PST 午夜兜底。
	body := []byte(`{"error":{"code":429,"message":"You have exhausted your capacity on this model. Your quota will reset after 6h53m10s."placeholderplaceholder`)

	tests := []struct {
		name              string
		account           *Account
		expectRateLimited bool
placeholder{
		{
			name: "pool_mode_apikey_stays_in_pool",
			account: &Account{
				ID:          600,
				Platform:    PlatformGemini,
				Type:        AccountTypeAPIKey,
		placeholder"pool_mode": trueplaceholder,
		placeholder,
			expectRateLimited: false,
	placeholder,
		{
			name: "custom_error_codes_hit_overrides_pool_mode",
			account: &Account{
				ID:       601,
		placeholder
				Type:     AccountTypeAPIKey,
		placeholder
					"pool_mode":                  true,
					"custom_error_codes_enabled": true,
					"custom_error_codes":         []any{float64(429)placeholder,
			placeholder,
		placeholder,
			expectRateLimited: true,
	placeholder,
		{
			name: "custom_error_codes_miss_skips",
			account: &Account{
				ID:       602,
		placeholder
				Type:     AccountTypeAPIKey,
		placeholder
					"pool_mode":                  true,
					"custom_error_codes_enabled": true,
					"custom_error_codes":         []any{float64(500)placeholder,
			placeholder,
		placeholder,
			expectRateLimited: false,
	placeholder,
		{
			name: "non_pool_apikey_still_rate_limited",
			account: &Account{
				ID:       603,
		placeholder
				Type:     AccountTypeAPIKey,
		placeholder,
			expectRateLimited: true,
	placeholder,
		{
			name: "oauth_account_ignores_pool_mode_flag",
			account: &Account{
				ID:          604,
				Platform:    PlatformGemini,
				Type:        AccountTypeOAuth,
		placeholder"pool_mode": trueplaceholder,
		placeholder,
			expectRateLimited: true,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &rateLimit429AccountRepoStub{placeholder
			svc := &GeminiMessagesCompatService{
				accountRepo:      repo,
				rateLimitService: NewRateLimitService(repo, nil, &config.Config{placeholder, nil, nil),
		placeholder

			svc.handleGeminiUpstreamError(context.Background(), tt.account, http.StatusTooManyRequests, http.Header{placeholder, body)

			if !tt.expectRateLimited {
				require.Zero(t, repo.rateLimitCalls, "池模式账号不应被标记账号级限流")
				return
		placeholder
			require.Equal(t, 1, repo.rateLimitCalls)
			require.Equal(t, tt.account.ID, repo.lastRateLimitID)
			require.True(t, repo.lastRateLimitReset.After(time.Now()))
	placeholder)
placeholder
placeholder

type geminiErrorPolicyRepo struct {
	mockAccountRepoForGemini
	setErrorCalls            int
	setRateLimitedCalls      int
	setTempCalls             int
	setModelRateLimitedCalls int
	lastModelScope           string
placeholder

func (r *geminiErrorPolicyRepo) SetError(_ context.Context, _ int64, _ string) error {
	r.setErrorCalls++
	return nil
placeholder

func (r *geminiErrorPolicyRepo) SetRateLimited(_ context.Context, _ int64, _ time.Time) error {
	r.setRateLimitedCalls++
	return nil
placeholder

func (r *geminiErrorPolicyRepo) SetTempUnschedulable(_ context.Context, _ int64, _ time.Time, _ string) error {
	r.setTempCalls++
	return nil
placeholder

func (r *geminiErrorPolicyRepo) SetModelRateLimit(_ context.Context, _ int64, scope string, _ time.Time, _ ...string) error {
	r.setModelRateLimitedCalls++
	r.lastModelScope = scope
	return nil
placeholder
