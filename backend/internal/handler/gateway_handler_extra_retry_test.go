//go:build unit

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// --- sleepFixedDelay ---

func TestSleepFixedDelay_ZeroDelay(t *testing.T) {
	got := sleepFixedDelay(context.Background(), 0)
	require.True(t, got, "zero delay should return true immediately")
placeholder

func TestSleepFixedDelay_NegativeDelay(t *testing.T) {
	got := sleepFixedDelay(context.Background(), -1*time.Second)
	require.True(t, got, "negative delay should return true immediately")
placeholder

func TestSleepFixedDelay_NormalDelay(t *testing.T) {
	start := time.Now()
	got := sleepFixedDelay(context.Background(), 50*time.Millisecond)
	elapsed := time.Since(start)
	require.True(t, got, "normal delay should return true")
	require.GreaterOrEqual(t, elapsed, 40*time.Millisecond, "should sleep at least ~50ms")
placeholder

func TestSleepFixedDelay_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	got := sleepFixedDelay(ctx, 10*time.Second)
	require.False(t, got, "cancelled context should return false")
placeholder

func TestSleepFixedDelay_ContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	got := sleepFixedDelay(ctx, 5*time.Second)
	require.False(t, got, "context timeout should return false before delay completes")
placeholder

// --- antigravityExtraRetryDelay constant ---

func TestAntigravityExtraRetryDelayValue(t *testing.T) {
	require.Equal(t, 500*time.Millisecond, antigravityExtraRetryDelay)
placeholder

// --- NewGatewayHandler antigravityExtraRetries field ---

func TestNewGatewayHandler_AntigravityExtraRetries_Default(t *testing.T) {
	h := NewGatewayHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	require.Equal(t, 10, h.antigravityExtraRetries, "default should be 10 when cfg is nil")
placeholder

func TestNewGatewayHandler_AntigravityExtraRetries_FromConfig(t *testing.T) {
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			AntigravityExtraRetries: 5,
	placeholder,
placeholder
	h := NewGatewayHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil, cfg)
	require.Equal(t, 5, h.antigravityExtraRetries, "should use config value")
placeholder

func TestNewGatewayHandler_AntigravityExtraRetries_ZeroDisables(t *testing.T) {
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			AntigravityExtraRetries: 0,
	placeholder,
placeholder
	h := NewGatewayHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil, cfg)
	require.Equal(t, 0, h.antigravityExtraRetries, "zero should disable extra retries")
placeholder

// --- handleFailoverAllAccountsExhausted (renamed: using handleFailoverExhausted) ---
// We test the error response format helpers that the extra retry path uses.

func TestHandleFailoverExhausted_JSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	h := &GatewayHandler{placeholder
	failoverErr := &service.UpstreamFailoverError{StatusCode: 429placeholder
	h.handleFailoverExhausted(c, failoverErr, service.PlatformAntigravity, false)

	require.Equal(t, http.StatusTooManyRequests, rec.Code)

	var body map[string]any
	err := json.Unmarshal(rec.Body.Bytes(), &body)
placeholder
	errObj, ok := body["error"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "rate_limit_error", errObj["type"])
placeholder

func TestHandleFailoverExhaustedSimple_JSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	h := &GatewayHandler{placeholder
	h.handleFailoverExhaustedSimple(c, 502, false)

	require.Equal(t, http.StatusBadGateway, rec.Code)

	var body map[string]any
	err := json.Unmarshal(rec.Body.Bytes(), &body)
placeholder
	errObj, ok := body["error"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "upstream_error", errObj["type"])
placeholder

// --- Extra retry platform filter logic ---

func TestExtraRetryPlatformFilter(t *testing.T) {
	tests := []struct {
		name             string
		switchCount      int
		maxAccountSwitch int
		platform         string
		expectSkip       bool
placeholder{
		{
			name:             "default_retry_phase_antigravity_not_skipped",
			switchCount:      1,
			maxAccountSwitch: 3,
			platform:         service.PlatformAntigravity,
			expectSkip:       false,
	placeholder,
		{
			name:             "default_retry_phase_gemini_not_skipped",
			switchCount:      1,
			maxAccountSwitch: 3,
			platform:         service.PlatformGemini,
			expectSkip:       false,
	placeholder,
		{
			name:             "extra_retry_phase_antigravity_not_skipped",
			switchCount:      3,
			maxAccountSwitch: 3,
			platform:         service.PlatformAntigravity,
			expectSkip:       false,
	placeholder,
		{
			name:             "extra_retry_phase_gemini_skipped",
			switchCount:      3,
			maxAccountSwitch: 3,
			platform:         service.PlatformGemini,
			expectSkip:       true,
	placeholder,
		{
			name:             "extra_retry_phase_anthropic_skipped",
			switchCount:      3,
			maxAccountSwitch: 3,
			platform:         service.PlatformAnthropic,
			expectSkip:       true,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Replicate the filter condition from the handler
			shouldSkip := tt.switchCount >= tt.maxAccountSwitch && tt.platform != service.PlatformAntigravity
			require.Equal(t, tt.expectSkip, shouldSkip)
	placeholder)
placeholder
placeholder

// --- Extra retry counter logic ---

func TestExtraRetryCounterExhaustion(t *testing.T) {
	tests := []struct {
		name               string
		maxExtraRetries    int
		currentExtraCount  int
		expectExhausted    bool
placeholder{
		{
			name:              "first_extra_retry",
			maxExtraRetries:   10,
			currentExtraCount: 1,
			expectExhausted:   false,
	placeholder,
		{
			name:              "at_limit",
			maxExtraRetries:   10,
			currentExtraCount: 10,
			expectExhausted:   false,
	placeholder,
		{
			name:              "exceeds_limit",
			maxExtraRetries:   10,
			currentExtraCount: 11,
			expectExhausted:   true,
	placeholder,
		{
			name:              "zero_disables_extra_retry",
			maxExtraRetries:   0,
			currentExtraCount: 1,
			expectExhausted:   true,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Replicate the exhaustion condition: antigravityExtraCount > h.antigravityExtraRetries
			exhausted := tt.currentExtraCount > tt.maxExtraRetries
			require.Equal(t, tt.expectExhausted, exhausted)
	placeholder)
placeholder
placeholder

// --- mapUpstreamError (used by handleFailoverExhausted) ---

func TestMapUpstreamError(t *testing.T) {
	h := &GatewayHandler{placeholder
	tests := []struct {
		name           string
		statusCode     int
		expectedStatus int
		expectedType   string
placeholder{
		{"429", 429, http.StatusTooManyRequests, "rate_limit_error"placeholder,
		{"529", 529, http.StatusServiceUnavailable, "overloaded_error"placeholder,
		{"500", 500, http.StatusBadGateway, "upstream_error"placeholder,
		{"502", 502, http.StatusBadGateway, "upstream_error"placeholder,
		{"503", 503, http.StatusBadGateway, "upstream_error"placeholder,
		{"504", 504, http.StatusBadGateway, "upstream_error"placeholder,
		{"401", 401, http.StatusBadGateway, "upstream_error"placeholder,
		{"403", 403, http.StatusBadGateway, "upstream_error"placeholder,
		{"unknown", 418, http.StatusBadGateway, "upstream_error"placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, errType, _ := h.mapUpstreamError(tt.statusCode)
			require.Equal(t, tt.expectedStatus, status)
			require.Equal(t, tt.expectedType, errType)
	placeholder)
placeholder
placeholder

// --- Gemini native path: handleGeminiFailoverExhausted ---

func TestHandleGeminiFailoverExhausted_NilError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	h := &GatewayHandler{placeholder
	h.handleGeminiFailoverExhausted(c, nil)

	require.Equal(t, http.StatusBadGateway, rec.Code)
	var body map[string]any
	err := json.Unmarshal(rec.Body.Bytes(), &body)
placeholder
	errObj, ok := body["error"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "Upstream request failed", errObj["message"])
placeholder

func TestHandleGeminiFailoverExhausted_429(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	h := &GatewayHandler{placeholder
	failoverErr := &service.UpstreamFailoverError{StatusCode: 429placeholder
	h.handleGeminiFailoverExhausted(c, failoverErr)

	require.Equal(t, http.StatusTooManyRequests, rec.Code)
placeholder

// --- handleStreamingAwareError streaming mode ---

func TestHandleStreamingAwareError_StreamStarted(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	// Simulate stream already started: set content type and write initial data
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.WriteHeaderNow()

	h := &GatewayHandler{placeholder
	h.handleStreamingAwareError(c, http.StatusTooManyRequests, "rate_limit_error", "test error", true)

	body := rec.Body.String()
	require.Contains(t, body, "rate_limit_error")
	require.Contains(t, body, "test error")
	require.Contains(t, body, "data: ")
placeholder

func TestHandleStreamingAwareError_NotStreaming(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	h := &GatewayHandler{placeholder
	h.handleStreamingAwareError(c, http.StatusServiceUnavailable, "api_error", "no model", false)

	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
	var body map[string]any
	err := json.Unmarshal(rec.Body.Bytes(), &body)
placeholder
	errObj, ok := body["error"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "api_error", errObj["type"])
	require.Equal(t, "no model", errObj["message"])
placeholder

// --- Integration: extra retry flow simulation ---

func TestExtraRetryFlowSimulation(t *testing.T) {
	// Simulate the full extra retry flow logic
	maxAccountSwitches := 3
	maxExtraRetries := 2
	switchCount := 0
	antigravityExtraCount := 0

	type attempt struct {
		platform   string
		isFailover bool
placeholder

	// Simulate: 3 default retries (all fail), then 2 extra retries (all fail), then exhausted
	attempts := []attempt{
		{service.PlatformAntigravity, trueplaceholder,  // switchCount 0 -> 1
		{service.PlatformGemini, trueplaceholder,        // switchCount 1 -> 2
		{service.PlatformAntigravity, trueplaceholder,  // switchCount 2 -> 3 (reaches max)
		{service.PlatformAntigravity, trueplaceholder,  // extra retry 1
		{service.PlatformAntigravity, trueplaceholder,  // extra retry 2
		{service.PlatformAntigravity, trueplaceholder,  // extra retry 3 -> exhausted
placeholder

	var exhausted bool
	var skipped int

	for _, a := range attempts {
		if exhausted {
			break
	placeholder

		// Extra retry phase: skip non-Antigravity
		if switchCount >= maxAccountSwitches && a.platform != service.PlatformAntigravity {
			skipped++
			continue
	placeholder

		if a.isFailover {
			if switchCount >= maxAccountSwitches {
				antigravityExtraCount++
				if antigravityExtraCount > maxExtraRetries {
					exhausted = true
					continue
			placeholder
				// extra retry delay + continue
				continue
		placeholder
			switchCount++
	placeholder
placeholder

	require.Equal(t, 3, switchCount, "should have 3 default retries")
	require.Equal(t, 3, antigravityExtraCount, "counter incremented 3 times")
	require.True(t, exhausted, "should be exhausted after exceeding max extra retries")
	require.Equal(t, 0, skipped, "no non-antigravity accounts in this simulation")
placeholder

func TestExtraRetryFlowSimulation_SkipsNonAntigravity(t *testing.T) {
	maxAccountSwitches := 2
	switchCount := 2 // already past default retries
	antigravityExtraCount := 0
	maxExtraRetries := 5

	type accountSelection struct {
		platform string
placeholder

	selections := []accountSelection{
		{service.PlatformGeminiplaceholder,       // should be skipped
		{service.PlatformAnthropicplaceholder,    // should be skipped
		{service.PlatformAntigravityplaceholder,  // should be attempted
placeholder

	var skippedCount int
	var attemptedCount int

	for _, sel := range selections {
		if switchCount >= maxAccountSwitches && sel.platform != service.PlatformAntigravity {
			skippedCount++
			continue
	placeholder
		// Simulate failover
		antigravityExtraCount++
		if antigravityExtraCount > maxExtraRetries {
			break
	placeholder
		attemptedCount++
placeholder

	require.Equal(t, 2, skippedCount, "gemini and anthropic accounts should be skipped")
	require.Equal(t, 1, attemptedCount, "only antigravity account should be attempted")
	require.Equal(t, 1, antigravityExtraCount)
placeholder
