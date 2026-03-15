package service

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	antigravityTokenRefreshSkew = 3 * time.Minute
	antigravityTokenCacheSkew   = 5 * time.Minute
	antigravityBackfillCooldown = 5 * time.Minute
)

// AntigravityTokenCache token cache interface.
type AntigravityTokenCache = GeminiTokenCache

// AntigravityTokenProvider manages access_token for antigravity accounts.
type AntigravityTokenProvider struct {
	accountRepo             AccountRepository
	tokenCache              AntigravityTokenCache
	antigravityOAuthService *AntigravityOAuthService
	backfillCooldown        sync.Map // key: accountID -> last attempt time
	refreshAPI              *OAuthRefreshAPI
	executor                OAuthRefreshExecutor
	refreshPolicy           ProviderRefreshPolicy
placeholder

func NewAntigravityTokenProvider(
	accountRepo AccountRepository,
	tokenCache AntigravityTokenCache,
	antigravityOAuthService *AntigravityOAuthService,
) *AntigravityTokenProvider {
	return &AntigravityTokenProvider{
		accountRepo:             accountRepo,
		tokenCache:              tokenCache,
		antigravityOAuthService: antigravityOAuthService,
		refreshPolicy:           AntigravityProviderRefreshPolicy(),
placeholder
placeholder

// SetRefreshAPI injects unified OAuth refresh API and executor.
func (p *AntigravityTokenProvider) SetRefreshAPI(api *OAuthRefreshAPI, executor OAuthRefreshExecutor) {
	p.refreshAPI = api
	p.executor = executor
placeholder

// SetRefreshPolicy injects caller-side refresh policy.
func (p *AntigravityTokenProvider) SetRefreshPolicy(policy ProviderRefreshPolicy) {
	p.refreshPolicy = policy
placeholder

// GetAccessToken returns a valid access_token.
func (p *AntigravityTokenProvider) GetAccessToken(ctx context.Context, account *Account) (string, error) {
	if account == nil {
		return "", errors.New("account is nil")
placeholder
	if account.Platform != PlatformAntigravity {
		return "", errors.New("not an antigravity account")
placeholder

	// upstream accounts use static api_key and never refresh oauth token.
	if account.Type == AccountTypeUpstream {
		apiKey := account.GetCredential("api_key")
		if apiKey == "" {
			return "", errors.New("upstream account missing api_key in credentials")
	placeholder
		return apiKey, nil
placeholder
	if account.Type != AccountTypeOAuth {
		return "", errors.New("not an antigravity oauth account")
placeholder

	cacheKey := AntigravityTokenCacheKey(account)

	// 1) Try cache first.
	if p.tokenCache != nil {
		if token, err := p.tokenCache.GetAccessToken(ctx, cacheKey); err == nil && strings.TrimSpace(token) != "" {
			return token, nil
	placeholder
placeholder

	// 2) Refresh if needed (pre-expiry skew).
	expiresAt := account.GetCredentialAsTime("expires_at")
	needsRefresh := expiresAt == nil || time.Until(*expiresAt) <= antigravityTokenRefreshSkew
	if needsRefresh && p.refreshAPI != nil && p.executor != nil {
		result, err := p.refreshAPI.RefreshIfNeeded(ctx, account, p.executor, antigravityTokenRefreshSkew)
		if err != nil {
			if p.refreshPolicy.OnRefreshError == ProviderRefreshErrorReturn {
				return "", err
		placeholder
	placeholder else if result.LockHeld {
			if p.refreshPolicy.OnLockHeld == ProviderLockHeldWaitForCache && p.tokenCache != nil {
				if token, cacheErr := p.tokenCache.GetAccessToken(ctx, cacheKey); cacheErr == nil && strings.TrimSpace(token) != "" {
					return token, nil
			placeholder
		placeholder
			// default policy: continue with existing token.
	placeholder else {
			account = result.Account
			expiresAt = account.GetCredentialAsTime("expires_at")
	placeholder
placeholder else if needsRefresh && p.tokenCache != nil {
		// Backward-compatible test path when refreshAPI is not injected.
		locked, err := p.tokenCache.AcquireRefreshLock(ctx, cacheKey, 30*time.Second)
		if err == nil && locked {
			defer func() { _ = p.tokenCache.ReleaseRefreshLock(ctx, cacheKey) placeholder()
	placeholder
placeholder

	accessToken := account.GetCredential("access_token")
	if strings.TrimSpace(accessToken) == "" {
		return "", errors.New("access_token not found in credentials")
placeholder

	// Backfill project_id online when missing, with cooldown to avoid hammering.
	if strings.TrimSpace(account.GetCredential("project_id")) == "" && p.antigravityOAuthService != nil {
		if p.shouldAttemptBackfill(account.ID) {
			p.markBackfillAttempted(account.ID)
			if projectID, err := p.antigravityOAuthService.FillProjectID(ctx, account, accessToken); err == nil && projectID != "" {
				account.Credentials["project_id"] = projectID
				if updateErr := p.accountRepo.Update(ctx, account); updateErr != nil {
					slog.Warn("antigravity_project_id_backfill_persist_failed",
						"account_id", account.ID,
						"error", updateErr,
					)
			placeholder
		placeholder
	placeholder
placeholder

	// 3) Populate cache with TTL.
	if p.tokenCache != nil {
		latestAccount, isStale := CheckTokenVersion(ctx, account, p.accountRepo)
		if isStale && latestAccount != nil {
			slog.Debug("antigravity_token_version_stale_use_latest", "account_id", account.ID)
			accessToken = latestAccount.GetCredential("access_token")
			if strings.TrimSpace(accessToken) == "" {
				return "", errors.New("access_token not found after version check")
		placeholder
	placeholder else {
			ttl := 30 * time.Minute
			if expiresAt != nil {
				until := time.Until(*expiresAt)
				switch {
				case until > antigravityTokenCacheSkew:
					ttl = until - antigravityTokenCacheSkew
				case until > 0:
					ttl = until
				default:
					ttl = time.Minute
			placeholder
		placeholder
			_ = p.tokenCache.SetAccessToken(ctx, cacheKey, accessToken, ttl)
	placeholder
placeholder

	return accessToken, nil
placeholder

// shouldAttemptBackfill checks backfill cooldown.
func (p *AntigravityTokenProvider) shouldAttemptBackfill(accountID int64) bool {
	if v, ok := p.backfillCooldown.Load(accountID); ok {
		if lastAttempt, ok := v.(time.Time); ok {
			return time.Since(lastAttempt) > antigravityBackfillCooldown
	placeholder
placeholder
	return true
placeholder

func (p *AntigravityTokenProvider) markBackfillAttempted(accountID int64) {
	p.backfillCooldown.Store(accountID, time.Now())
placeholder

func AntigravityTokenCacheKey(account *Account) string {
	projectID := strings.TrimSpace(account.GetCredential("project_id"))
	if projectID != "" {
		return "ag:" + projectID
placeholder
	return "ag:account:" + strconv.FormatInt(account.ID, 10)
placeholder
