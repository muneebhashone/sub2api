package service

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
)

type stubOpenAIAccountRepo struct {
	AccountRepository
	accounts []Account
placeholder

func (r stubOpenAIAccountRepo) ListSchedulableByGroupIDAndPlatform(ctx context.Context, groupID int64, platform string) ([]Account, error) {
	return append([]Account(nil), r.accounts...), nil
placeholder

func (r stubOpenAIAccountRepo) ListSchedulableByPlatform(ctx context.Context, platform string) ([]Account, error) {
	return append([]Account(nil), r.accounts...), nil
placeholder

type stubConcurrencyCache struct {
	ConcurrencyCache
placeholder

func (c stubConcurrencyCache) AcquireAccountSlot(ctx context.Context, accountID int64, maxConcurrency int, requestID string) (bool, error) {
	return true, nil
placeholder

func (c stubConcurrencyCache) ReleaseAccountSlot(ctx context.Context, accountID int64, requestID string) error {
	return nil
placeholder

func (c stubConcurrencyCache) GetAccountsLoadBatch(ctx context.Context, accounts []AccountWithConcurrency) (map[int64]*AccountLoadInfo, error) {
	out := make(map[int64]*AccountLoadInfo, len(accounts))
	for _, acc := range accounts {
		out[acc.ID] = &AccountLoadInfo{AccountID: acc.ID, LoadRate: 0placeholder
placeholder
	return out, nil
placeholder

func TestOpenAISelectAccountWithLoadAwareness_FiltersUnschedulable(t *testing.T) {
	now := time.Now()
	resetAt := now.Add(10 * time.Minute)
	groupID := int64(1)

	rateLimited := Account{
		ID:               1,
		Platform:         PlatformOpenAI,
		Type:             AccountTypeAPIKey,
		Status:           StatusActive,
		Schedulable:      true,
		Concurrency:      1,
		Priority:         0,
		RateLimitResetAt: &resetAt,
placeholder
	available := Account{
		ID:          2,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
		Priority:    1,
placeholder

	svc := &OpenAIGatewayService{
		accountRepo:        stubOpenAIAccountRepo{accounts: []Account{rateLimited, availableplaceholderplaceholder,
		concurrencyService: NewConcurrencyService(stubConcurrencyCache{placeholder),
placeholder

	selection, err := svc.SelectAccountWithLoadAwareness(context.Background(), &groupID, "", "gpt-5.2", nil)
	if err != nil {
		t.Fatalf("SelectAccountWithLoadAwareness error: %v", err)
placeholder
	if selection == nil || selection.Account == nil {
		t.Fatalf("expected selection with account")
placeholder
	if selection.Account.ID != available.ID {
		t.Fatalf("expected account %d, got %d", available.ID, selection.Account.ID)
placeholder
	if selection.ReleaseFunc != nil {
		selection.ReleaseFunc()
placeholder
placeholder

func TestOpenAISelectAccountWithLoadAwareness_FiltersUnschedulableWhenNoConcurrencyService(t *testing.T) {
	now := time.Now()
	resetAt := now.Add(10 * time.Minute)
	groupID := int64(1)

	rateLimited := Account{
		ID:               1,
		Platform:         PlatformOpenAI,
		Type:             AccountTypeAPIKey,
		Status:           StatusActive,
		Schedulable:      true,
		Concurrency:      1,
		Priority:         0,
		RateLimitResetAt: &resetAt,
placeholder
	available := Account{
		ID:          2,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
		Priority:    1,
placeholder

	svc := &OpenAIGatewayService{
		accountRepo: stubOpenAIAccountRepo{accounts: []Account{rateLimited, availableplaceholderplaceholder,
		// concurrencyService is nil, forcing the non-load-batch selection path.
placeholder

	selection, err := svc.SelectAccountWithLoadAwareness(context.Background(), &groupID, "", "gpt-5.2", nil)
	if err != nil {
		t.Fatalf("SelectAccountWithLoadAwareness error: %v", err)
placeholder
	if selection == nil || selection.Account == nil {
		t.Fatalf("expected selection with account")
placeholder
	if selection.Account.ID != available.ID {
		t.Fatalf("expected account %d, got %d", available.ID, selection.Account.ID)
placeholder
	if selection.ReleaseFunc != nil {
		selection.ReleaseFunc()
placeholder
placeholder

func TestOpenAIStreamingTimeout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 1,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       pr,
		Header:     http.Header{placeholder,
placeholder

	start := time.Now()
	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1placeholder, start, "model", "model")
	_ = pw.Close()
	_ = pr.Close()

	if err == nil || !strings.Contains(err.Error(), "stream data interval timeout") {
		t.Fatalf("expected stream timeout error, got %v", err)
placeholder
	if !strings.Contains(rec.Body.String(), "stream_timeout") {
		t.Fatalf("expected stream_timeout SSE error, got %q", rec.Body.String())
placeholder
placeholder

func TestOpenAIStreamingTooLong(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               64 * 1024,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       pr,
		Header:     http.Header{placeholder,
placeholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		// 写入超过 MaxLineSize 的单行数据，触发 ErrTooLong
		payload := "data: " + strings.Repeat("a", 128*1024) + "\n"
		_, _ = pw.Write([]byte(payload))
placeholder()

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 2placeholder, time.Now(), "model", "model")
	_ = pr.Close()

	if !errors.Is(err, bufio.ErrTooLong) {
		t.Fatalf("expected ErrTooLong, got %v", err)
placeholder
	if !strings.Contains(rec.Body.String(), "response_too_large") {
		t.Fatalf("expected response_too_large SSE error, got %q", rec.Body.String())
placeholder
placeholder

func TestOpenAINonStreamingContentTypePassThrough(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Security: config.SecurityConfig{
			ResponseHeaders: config.ResponseHeaderConfig{Enabled: falseplaceholder,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	body := []byte(`{"usage":{"input_tokens":1,"output_tokens":2,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholder`)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(body)),
		Header:     http.Header{"Content-Type": []string{"application/vnd.test+json"placeholderplaceholder,
placeholder

	_, err := svc.handleNonStreamingResponse(c.Request.Context(), resp, c, &Account{placeholder, "model", "model")
	if err != nil {
		t.Fatalf("handleNonStreamingResponse error: %v", err)
placeholder

	if !strings.Contains(rec.Header().Get("Content-Type"), "application/vnd.test+json") {
		t.Fatalf("expected Content-Type passthrough, got %q", rec.Header().Get("Content-Type"))
placeholder
placeholder

func TestOpenAINonStreamingContentTypeDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Security: config.SecurityConfig{
			ResponseHeaders: config.ResponseHeaderConfig{Enabled: falseplaceholder,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	body := []byte(`{"usage":{"input_tokens":1,"output_tokens":2,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholder`)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(body)),
		Header:     http.Header{placeholder,
placeholder

	_, err := svc.handleNonStreamingResponse(c.Request.Context(), resp, c, &Account{placeholder, "model", "model")
	if err != nil {
		t.Fatalf("handleNonStreamingResponse error: %v", err)
placeholder

	if !strings.Contains(rec.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("expected default Content-Type, got %q", rec.Header().Get("Content-Type"))
placeholder
placeholder

func TestOpenAIStreamingHeadersOverride(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Security: config.SecurityConfig{
			ResponseHeaders: config.ResponseHeaderConfig{Enabled: falseplaceholder,
	placeholder,
		Gateway: config.GatewayConfig{
			StreamDataIntervalTimeout: 0,
			StreamKeepaliveInterval:   0,
			MaxLineSize:               defaultMaxLineSize,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	pr, pw := io.Pipe()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       pr,
		Header: http.Header{
			"Cache-Control": []string{"upstream"placeholder,
			"X-Request-Id":  []string{"req-123"placeholder,
			"Content-Type":  []string{"application/custom"placeholder,
	placeholder,
placeholder

	go func() {
		defer func() { _ = pw.Close() placeholder()
		_, _ = pw.Write([]byte("data: {placeholder\n\n"))
placeholder()

	_, err := svc.handleStreamingResponse(c.Request.Context(), resp, c, &Account{ID: 1placeholder, time.Now(), "model", "model")
	_ = pr.Close()
	if err != nil {
		t.Fatalf("handleStreamingResponse error: %v", err)
placeholder

	if rec.Header().Get("Cache-Control") != "no-cache" {
		t.Fatalf("expected Cache-Control override, got %q", rec.Header().Get("Cache-Control"))
placeholder
	if rec.Header().Get("Content-Type") != "text/event-stream" {
		t.Fatalf("expected Content-Type override, got %q", rec.Header().Get("Content-Type"))
placeholder
	if rec.Header().Get("X-Request-Id") != "req-123" {
		t.Fatalf("expected X-Request-Id passthrough, got %q", rec.Header().Get("X-Request-Id"))
placeholder
placeholder

func TestOpenAIInvalidBaseURLWhenAllowlistDisabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		Security: config.SecurityConfig{
			URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholder,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	account := &Account{
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
placeholder"base_url": "://invalid-url"placeholder,
placeholder

	_, err := svc.buildUpstreamRequest(c.Request.Context(), c, account, []byte("{placeholder"), "token", false, "", false)
	if err == nil {
		t.Fatalf("expected error for invalid base_url when allowlist disabled")
placeholder
placeholder

func TestOpenAIValidateUpstreamBaseURLDisabledRequiresHTTPS(t *testing.T) {
	cfg := &config.Config{
		Security: config.SecurityConfig{
			URLAllowlist: config.URLAllowlistConfig{Enabled: falseplaceholder,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	if _, err := svc.validateUpstreamBaseURL("http://not-https.example.com"); err == nil {
		t.Fatalf("expected http to be rejected when allow_insecure_http is false")
placeholder
	normalized, err := svc.validateUpstreamBaseURL("https://example.com")
	if err != nil {
		t.Fatalf("expected https to be allowed when allowlist disabled, got %v", err)
placeholder
	if normalized != "https://example.com" {
		t.Fatalf("expected raw url passthrough, got %q", normalized)
placeholder
placeholder

func TestOpenAIValidateUpstreamBaseURLDisabledAllowsHTTP(t *testing.T) {
	cfg := &config.Config{
		Security: config.SecurityConfig{
			URLAllowlist: config.URLAllowlistConfig{
				Enabled:           false,
				AllowInsecureHTTP: true,
		placeholder,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	normalized, err := svc.validateUpstreamBaseURL("http://not-https.example.com")
	if err != nil {
		t.Fatalf("expected http allowed when allow_insecure_http is true, got %v", err)
placeholder
	if normalized != "http://not-https.example.com" {
		t.Fatalf("expected raw url passthrough, got %q", normalized)
placeholder
placeholder

func TestOpenAIValidateUpstreamBaseURLEnabledEnforcesAllowlist(t *testing.T) {
	cfg := &config.Config{
		Security: config.SecurityConfig{
			URLAllowlist: config.URLAllowlistConfig{
				Enabled:       true,
				UpstreamHosts: []string{"example.com"placeholder,
		placeholder,
	placeholder,
placeholder
	svc := &OpenAIGatewayService{cfg: cfgplaceholder

	if _, err := svc.validateUpstreamBaseURL("https://example.com"); err != nil {
		t.Fatalf("expected allowlisted host to pass, got %v", err)
placeholder
	if _, err := svc.validateUpstreamBaseURL("https://evil.com"); err == nil {
		t.Fatalf("expected non-allowlisted host to fail")
placeholder
placeholder
