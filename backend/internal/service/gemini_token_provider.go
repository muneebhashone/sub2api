package service

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/model"
	"github.com/Wei-Shaw/sub2api/internal/service/ports"
)

const (
	geminiTokenRefreshSkew = 3 * time.Minute
	geminiTokenCacheSkew   = 5 * time.Minute
)

type GeminiTokenProvider struct {
	accountRepo        ports.AccountRepository
	tokenCache         ports.GeminiTokenCache
	geminiOAuthService *GeminiOAuthService
placeholder

func NewGeminiTokenProvider(
	accountRepo ports.AccountRepository,
	tokenCache ports.GeminiTokenCache,
	geminiOAuthService *GeminiOAuthService,
) *GeminiTokenProvider {
	return &GeminiTokenProvider{
		accountRepo:        accountRepo,
		tokenCache:         tokenCache,
		geminiOAuthService: geminiOAuthService,
placeholder
placeholder

func (p *GeminiTokenProvider) GetAccessToken(ctx context.Context, account *model.Account) (string, error) {
	if account == nil {
		return "", errors.New("account is nil")
placeholder
	if account.Platform != model.PlatformGemini || account.Type != model.AccountTypeOAuth {
		return "", errors.New("not a gemini oauth account")
placeholder

	cacheKey := geminiTokenCacheKey(account)

	// 1) Try cache first.
	if p.tokenCache != nil {
		if token, err := p.tokenCache.GetAccessToken(ctx, cacheKey); err == nil && strings.TrimSpace(token) != "" {
			return token, nil
	placeholder
placeholder

	// 2) Refresh if needed (pre-expiry skew).
	expiresAt := parseExpiresAt(account)
	needsRefresh := expiresAt == nil || time.Until(*expiresAt) <= geminiTokenRefreshSkew
	if needsRefresh && p.tokenCache != nil {
		locked, err := p.tokenCache.AcquireRefreshLock(ctx, cacheKey, 30*time.Second)
		if err == nil && locked {
			defer func() { _ = p.tokenCache.ReleaseRefreshLock(ctx, cacheKey) placeholder()

			// Re-check after lock (another worker may have refreshed).
			if token, err := p.tokenCache.GetAccessToken(ctx, cacheKey); err == nil && strings.TrimSpace(token) != "" {
				return token, nil
		placeholder

			fresh, err := p.accountRepo.GetByID(ctx, account.ID)
			if err == nil && fresh != nil {
				account = fresh
		placeholder
			expiresAt = parseExpiresAt(account)
			if expiresAt == nil || time.Until(*expiresAt) <= geminiTokenRefreshSkew {
				if p.geminiOAuthService == nil {
					return "", errors.New("gemini oauth service not configured")
			placeholder
				tokenInfo, err := p.geminiOAuthService.RefreshAccountToken(ctx, account)
				if err != nil {
					return "", err
			placeholder
				newCredentials := p.geminiOAuthService.BuildAccountCredentials(tokenInfo)
				for k, v := range account.Credentials {
					if _, exists := newCredentials[k]; !exists {
						newCredentials[k] = v
				placeholder
			placeholder
				account.Credentials = model.JSONB(newCredentials)
				_ = p.accountRepo.Update(ctx, account)
				expiresAt = parseExpiresAt(account)
		placeholder
	placeholder
placeholder

	accessToken := account.GetCredential("access_token")
	if strings.TrimSpace(accessToken) == "" {
		return "", errors.New("access_token not found in credentials")
placeholder

	// project_id is optional now:
	// - If present: will use Code Assist API (requires project_id)
	// - If absent: will use AI Studio API with OAuth token (like regular API key mode)
	// Auto-detect project_id only if explicitly enabled via a credential flag
	projectID := strings.TrimSpace(account.GetCredential("project_id"))
	autoDetectProjectID := account.GetCredential("auto_detect_project_id") == "true"

	if projectID == "" && autoDetectProjectID {
		if p.geminiOAuthService == nil {
			return accessToken, nil // Fallback to AI Studio API mode
	placeholder

		var proxyURL string
		if account.ProxyID != nil && p.geminiOAuthService.proxyRepo != nil {
			if proxy, err := p.geminiOAuthService.proxyRepo.GetByID(ctx, *account.ProxyID); err == nil && proxy != nil {
				proxyURL = proxy.URL()
		placeholder
	placeholder

		detected, err := p.geminiOAuthService.fetchProjectID(ctx, accessToken, proxyURL)
		if err != nil {
			log.Printf("[GeminiTokenProvider] Auto-detect project_id failed: %v, fallback to AI Studio API mode", err)
			return accessToken, nil
	placeholder
		detected = strings.TrimSpace(detected)
		if detected != "" {
			if account.Credentials == nil {
				account.Credentials = model.JSONB{placeholder
		placeholder
			account.Credentials["project_id"] = detected
			_ = p.accountRepo.Update(ctx, account)
	placeholder
placeholder

	// 3) Populate cache with TTL.
	if p.tokenCache != nil {
		ttl := 30 * time.Minute
		if expiresAt != nil {
			until := time.Until(*expiresAt)
			switch {
			case until > geminiTokenCacheSkew:
				ttl = until - geminiTokenCacheSkew
			case until > 0:
				ttl = until
			default:
				ttl = time.Minute
		placeholder
	placeholder
		_ = p.tokenCache.SetAccessToken(ctx, cacheKey, accessToken, ttl)
placeholder

	return accessToken, nil
placeholder

func geminiTokenCacheKey(account *model.Account) string {
	projectID := strings.TrimSpace(account.GetCredential("project_id"))
	if projectID != "" {
		return projectID
placeholder
	return "account:" + strconv.FormatInt(account.ID, 10)
placeholder

func parseExpiresAt(account *model.Account) *time.Time {
	raw := strings.TrimSpace(account.GetCredential("expires_at"))
	if raw == "" {
		return nil
placeholder
	if unixSec, err := strconv.ParseInt(raw, 10, 64); err == nil && unixSec > 0 {
		t := time.Unix(unixSec, 0)
		return &t
placeholder
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return &t
placeholder
	return nil
placeholder
