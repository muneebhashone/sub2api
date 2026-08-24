package service

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/tidwall/gjson"
)

const (
	grokObservedModelsExtraKey = "grok_observed_models"
	grokObservedModelsTTL      = 6 * time.Hour
	grokObservedModelsTimeout  = 15 * time.Second
)

type grokObservedModelsSnapshot struct {
	Models    []string `json:"models"`
	FetchedAt string   `json:"fetched_at"`
	Source    string   `json:"source,omitempty"`
placeholder

var grokObservedModelsFlight sync.Map // accountID -> *singleflight-ish in-flight

// scheduleGrokObservedModelsSync best-effort fetches upstream /v1/models for a
// Grok OAuth account and stores IDs in Extra. Never blocks request path long;
// callers should fire-and-forget after successful auth/probe.
func (s *GrokQuotaService) scheduleGrokObservedModelsSync(account *Account) {
	if s == nil || account == nil || !account.IsGrokOAuth() || s.accountRepo == nil {
		return
placeholder
	id := account.ID
	if _, loaded := grokObservedModelsFlight.LoadOrStore(id, struct{placeholder{placeholder); loaded {
		return
placeholder
	// Copy credentials for background use.
	acc := *account
	go func() {
		defer grokObservedModelsFlight.Delete(id)
		ctx, cancel := context.WithTimeout(context.Background(), grokObservedModelsTimeout)
		defer cancel()
		if err := s.syncGrokObservedModels(ctx, &acc); err != nil {
			slog.Debug("grok_observed_models_sync_failed", "account_id", id, "error", err)
	placeholder
placeholder()
placeholder

func (s *GrokQuotaService) syncGrokObservedModels(ctx context.Context, account *Account) error {
	if s == nil || account == nil {
		return nil
placeholder
	// Skip if snapshot is still fresh.
	if snap := parseGrokObservedModels(account.Extra); snap != nil {
		if t, err := time.Parse(time.RFC3339, snap.FetchedAt); err == nil && time.Since(t) < grokObservedModelsTTL {
			return nil
	placeholder
placeholder
	token := strings.TrimSpace(account.GetGrokAccessToken())
	if token == "" && s.tokenProvider != nil {
		// Best-effort warm; avoid forcing refresh storms.
		if at, err := s.tokenProvider.GetAccessToken(ctx, account); err == nil {
			token = strings.TrimSpace(at)
	placeholder
placeholder
	if token == "" {
		return nil
placeholder
	baseURL := strings.TrimSpace(account.GetGrokBaseURL())
	if s.settingService != nil {
		baseURL = strings.TrimSpace(s.settingService.ResolveGrokBaseURL(ctx, account))
placeholder
	if baseURL == "" {
		baseURL = xai.DefaultCLIBaseURL
placeholder
	validator, err := grokBaseURLValidator(account, s.cfg)
	if err != nil {
		return err
placeholder
	validatedBaseURL, err := validator(baseURL)
	if err != nil {
		return err
placeholder
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, buildOpenAIModelsURL(validatedBaseURL), nil)
	if err != nil {
		return err
placeholder
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", defaultGrokUpstreamUserAgent())
	if account.IsGrokOAuth() {
		applyGrokCLIHeaders(req.Header)
		if isGrokCLIProxyTarget(req.URL.String()) {
			if userID := strings.TrimSpace(account.GetCredential("sub")); userID != "" {
				req.Header.Set("X-UserID", userID)
		placeholder
			if email := strings.TrimSpace(account.GetCredential("email")); email != "" {
				req.Header.Set("X-Email", email)
		placeholder
	placeholder
placeholder
	account.ApplyHeaderOverrides(req.Header)

	proxyURL := ""
	if s.proxyRepo != nil && account.ProxyID != nil {
		if p, err := s.proxyRepo.GetByID(ctx, *account.ProxyID); err == nil && p != nil {
			proxyURL = p.URL()
	placeholder
placeholder
	if s.httpUpstream == nil {
		return nil
placeholder
	resp, err := s.httpUpstream.Do(req, proxyURL, account.ID, account.Concurrency)
	if err != nil {
		return err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
placeholder
	if resp.StatusCode >= 400 {
		return nil
placeholder
	ids := extractGrokModelIDsFromModelsBody(body)
	if len(ids) == 0 {
		return nil
placeholder
	snap := grokObservedModelsSnapshot{
		Models:    ids,
		FetchedAt: time.Now().UTC().Format(time.RFC3339),
		Source:    "upstream_v1_models",
placeholder
	raw, err := json.Marshal(snap)
	if err != nil {
		return err
placeholder
	var asMap map[string]any
	if err := json.Unmarshal(raw, &asMap); err != nil {
		return err
placeholder
	return s.accountRepo.UpdateExtra(ctx, account.ID, map[string]any{
		grokObservedModelsExtraKey: asMap,
placeholder)
placeholder

func extractGrokModelIDsFromModelsBody(body []byte) []string {
	data := gjson.GetBytes(body, "data")
	if !data.IsArray() {
		// Some gateways return a bare array.
		data = gjson.ParseBytes(body)
placeholder
	seen := make(map[string]struct{placeholder)
	var out []string
	data.ForEach(func(_, v gjson.Result) bool {
		id := strings.TrimSpace(v.Get("id").String())
		if id == "" {
			id = strings.TrimSpace(v.String())
	placeholder
		if id == "" {
			return true
	placeholder
		if _, ok := seen[id]; ok {
			return true
	placeholder
		seen[id] = struct{placeholder{placeholder
		out = append(out, id)
		return true
placeholder)
	return out
placeholder

func parseGrokObservedModels(extra map[string]any) *grokObservedModelsSnapshot {
	if extra == nil {
		return nil
placeholder
	raw, ok := extra[grokObservedModelsExtraKey]
	if !ok || raw == nil {
		return nil
placeholder
	b, err := json.Marshal(raw)
	if err != nil {
		return nil
placeholder
	var snap grokObservedModelsSnapshot
	if err := json.Unmarshal(b, &snap); err != nil {
		return nil
placeholder
	if len(snap.Models) == 0 {
		return nil
placeholder
	return &snap
placeholder
