package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
)

const (
	grokQuotaUpstreamTimeout = 20 * time.Second
	grokQuotaProbeInput      = "."
	grokQuotaDefaultModel    = "grok-4.3"
)

type GrokQuotaProbeResult struct {
	Source          string             `json:"source"`
	Model           string             `json:"model"`
	Snapshot        *xai.QuotaSnapshot `json:"snapshot,omitempty"`
	StatusCode      int                `json:"status_code,omitempty"`
	HeadersObserved bool               `json:"headers_observed"`
	ResetSupported  bool               `json:"reset_supported"`
	FetchedAt       int64              `json:"fetched_at"`
placeholder

type GrokQuotaResetResult struct {
	Supported bool   `json:"supported"`
	Code      string `json:"code"`
	Message   string `json:"message"`
placeholder

type GrokQuotaService struct {
	accountRepo   AccountRepository
	proxyRepo     ProxyRepository
	tokenProvider *GrokTokenProvider
	httpUpstream  HTTPUpstream
placeholder

func NewGrokQuotaService(
	accountRepo AccountRepository,
	proxyRepo ProxyRepository,
	tokenProvider *GrokTokenProvider,
	httpUpstream HTTPUpstream,
) *GrokQuotaService {
	return &GrokQuotaService{
		accountRepo:   accountRepo,
		proxyRepo:     proxyRepo,
		tokenProvider: tokenProvider,
		httpUpstream:  httpUpstream,
placeholder
placeholder

func (s *GrokQuotaService) ProbeUsage(ctx context.Context, accountID int64) (*GrokQuotaProbeResult, error) {
	account, token, proxyURL, err := s.prepareProbe(ctx, accountID)
	if err != nil {
		return nil, err
placeholder

	probeModel := grokQuotaProbeModel()
	body, err := buildGrokQuotaProbeBody(probeModel)
	if err != nil {
		return nil, infraerrors.Newf(http.StatusBadRequest, "GROK_QUOTA_PROBE_BODY_ERROR", "failed to build probe body: %v", err)
placeholder
	targetURL, err := xai.BuildResponsesURL(account.GetGrokBaseURL())
	if err != nil {
		return nil, infraerrors.Newf(http.StatusBadRequest, "GROK_QUOTA_BASE_URL_INVALID", "invalid Grok base_url: %v", err)
placeholder

	callCtx, cancel := context.WithTimeout(ctx, grokQuotaUpstreamTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(callCtx, http.MethodPost, targetURL, bytes.NewReader(body))
	if err != nil {
		return nil, infraerrors.Newf(http.StatusInternalServerError, "GROK_QUOTA_PROBE_REQUEST_BUILD_FAILED", "failed to build upstream request: %v", err)
placeholder
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	applyGrokCLIHeaders(req.Header)

	resp, err := s.httpUpstream.Do(req, proxyURL, account.ID, maxInt(account.Concurrency, 1))
	if err != nil {
		return nil, infraerrors.Newf(http.StatusBadGateway, "GROK_QUOTA_PROBE_REQUEST_FAILED", "upstream probe failed: %v", err)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	snapshot := xai.ObserveQuotaHeaders(resp.Header, resp.StatusCode, "active_probe")
	_ = s.accountRepo.UpdateExtra(ctx, account.ID, map[string]any{
		grokQuotaSnapshotExtraKey: snapshot,
placeholder)

	result := &GrokQuotaProbeResult{
		Source:          "active_probe",
		Model:           probeModel,
		Snapshot:        snapshot,
		StatusCode:      resp.StatusCode,
		HeadersObserved: snapshot.HeadersObserved,
		ResetSupported:  false,
		FetchedAt:       time.Now().Unix(),
placeholder
	if resp.StatusCode == http.StatusTooManyRequests {
		return result, nil
placeholder
	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 240))
		bodyText := truncate(strings.TrimSpace(string(bodyBytes)), 240)
		slog.Warn("grok_quota_probe_failed", "account_id", account.ID, "model", probeModel, "status", resp.StatusCode, "body", bodyText)
		return nil, infraerrors.Newf(mapUpstreamStatus(resp.StatusCode), "GROK_QUOTA_PROBE_UPSTREAM_ERROR", "upstream returned %d for probe model %q: %s", resp.StatusCode, probeModel, bodyText)
placeholder
	return result, nil
placeholder

func (s *GrokQuotaService) ResetQuota(ctx context.Context, accountID int64) (*GrokQuotaResetResult, error) {
	if _, err := s.loadGrokOAuthAccount(ctx, accountID); err != nil {
		return nil, err
placeholder
	return nil, infraerrors.New(http.StatusNotImplemented, "GROK_QUOTA_RESET_UNSUPPORTED", "xAI does not expose a Grok subscription quota reset endpoint for OAuth accounts")
placeholder

func (s *GrokQuotaService) prepareProbe(ctx context.Context, accountID int64) (*Account, string, string, error) {
	if s == nil || s.tokenProvider == nil || s.httpUpstream == nil {
		return nil, "", "", infraerrors.New(http.StatusInternalServerError, "GROK_QUOTA_NOT_CONFIGURED", "grok quota service is not configured")
placeholder
	account, err := s.loadGrokOAuthAccount(ctx, accountID)
	if err != nil {
		return nil, "", "", err
placeholder

	token, err := s.tokenProvider.GetAccessToken(ctx, account)
	if err != nil {
		return nil, "", "", infraerrors.Newf(http.StatusBadGateway, "GROK_QUOTA_TOKEN_UNAVAILABLE", "failed to acquire access token: %v", err)
placeholder
	if strings.TrimSpace(token) == "" {
		return nil, "", "", infraerrors.New(http.StatusBadGateway, "GROK_QUOTA_TOKEN_UNAVAILABLE", "access token is empty")
placeholder

	return account, token, s.resolveProxyURL(ctx, account), nil
placeholder

func (s *GrokQuotaService) resolveProxyURL(ctx context.Context, account *Account) string {
	if account == nil || account.ProxyID == nil {
		return ""
placeholder
	switch {
	case account.Proxy != nil:
		return account.Proxy.URL()
	case s != nil && s.proxyRepo != nil:
		if proxy, err := s.proxyRepo.GetByID(ctx, *account.ProxyID); err == nil && proxy != nil {
			return proxy.URL()
	placeholder
placeholder
	return ""
placeholder

func (s *GrokQuotaService) loadGrokOAuthAccount(ctx context.Context, accountID int64) (*Account, error) {
	if s == nil || s.accountRepo == nil {
		return nil, infraerrors.New(http.StatusInternalServerError, "GROK_QUOTA_NOT_CONFIGURED", "grok quota service is not configured")
placeholder
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, infraerrors.Newf(http.StatusNotFound, "GROK_QUOTA_ACCOUNT_NOT_FOUND", "account not found: %v", err)
placeholder
	if account == nil {
		return nil, infraerrors.New(http.StatusNotFound, "GROK_QUOTA_ACCOUNT_NOT_FOUND", "account not found")
placeholder
	if account.Platform != PlatformGrok {
		return nil, infraerrors.New(http.StatusBadRequest, "GROK_QUOTA_INVALID_PLATFORM", "account is not a Grok account")
placeholder
	if account.Type != AccountTypeOAuth {
		return nil, infraerrors.New(http.StatusBadRequest, "GROK_QUOTA_INVALID_TYPE", "account is not an OAuth account")
placeholder
	return account, nil
placeholder

func grokQuotaProbeModel() string {
	return grokQuotaDefaultModel
placeholder

func buildGrokQuotaProbeBody(model string) ([]byte, error) {
	model = strings.TrimSpace(model)
	if model == "" {
		model = grokQuotaDefaultModel
placeholder
	return json.Marshal(map[string]any{
		"model":             model,
		"input":             grokQuotaProbeInput,
		"max_output_tokens": 1,
		"store":             false,
placeholder)
placeholder

func maxInt(a, b int) int {
	if a > b {
		return a
placeholder
	return b
placeholder
