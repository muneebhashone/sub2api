package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	grokTokenCacheSkew          = 5 * time.Minute
	grokRequestRefreshTimeout   = 8 * time.Second
	grokRefreshLockWaitTimeout  = 2 * time.Second
	grokRefreshLockPollInterval = 25 * time.Millisecond
)

var (
	errGrokOAuthRefreshNotConfigured = errors.New("grok oauth refresh is not configured")
	errGrokOAuthRefreshTokenMissing  = errors.New("grok oauth refresh token is missing")
	errGrokOAuthAccessTokenMissing   = errors.New("grok oauth access token is missing")
	errGrokOAuthAccessTokenExpired   = errors.New("grok oauth access token is expired")
	errGrokOAuthConfiguredProxyMiss  = errors.New("grok oauth configured proxy is missing")
)

type GrokTokenCache = GeminiTokenCache

type GrokTokenProvider struct {
	accountRepo      AccountRepository
	tokenCache       GrokTokenCache
	refreshAPI       *OAuthRefreshAPI
	executor         OAuthRefreshExecutor
	refreshPolicy    ProviderRefreshPolicy
	tempUnschedCache TempUnschedCache
placeholder

func NewGrokTokenProvider(
	accountRepo AccountRepository,
	tokenCache GrokTokenCache,
) *GrokTokenProvider {
	return &GrokTokenProvider{
		accountRepo:   accountRepo,
		tokenCache:    tokenCache,
		refreshPolicy: GrokProviderRefreshPolicy(),
placeholder
placeholder

func (p *GrokTokenProvider) SetRefreshAPI(api *OAuthRefreshAPI, executor OAuthRefreshExecutor) {
	p.refreshAPI = api
	p.executor = executor
placeholder

func (p *GrokTokenProvider) SetRefreshPolicy(policy ProviderRefreshPolicy) {
	p.refreshPolicy = policy
placeholder

func (p *GrokTokenProvider) SetTempUnschedCache(cache TempUnschedCache) {
	p.tempUnschedCache = cache
placeholder

func (p *GrokTokenProvider) GetAccessToken(ctx context.Context, account *Account) (string, error) {
	if account == nil {
		return "", errors.New("account is nil")
placeholder
	if account.Platform != PlatformGrok || account.Type != AccountTypeOAuth {
		return "", errors.New("not a grok oauth account")
placeholder
	selectedProxyID := cloneGrokProxyID(account.ProxyID)
	if eligibilityErr := grokOAuthRequestAccountEligibilityError(account); eligibilityErr != nil {
		return "", withGrokCredentialFailureSnapshot(eligibilityErr, account)
placeholder

	expiresAt := account.GetCredentialAsTime("expires_at")
	accountAccessToken := strings.TrimSpace(account.GetGrokAccessToken())
	if accountAccessToken == "" {
		return "", withGrokCredentialFailureSnapshot(errGrokOAuthAccessTokenMissing, account)
placeholder
	if strings.TrimSpace(account.GetGrokRefreshToken()) == "" {
		return "", withGrokCredentialFailureSnapshot(errGrokOAuthRefreshTokenMissing, account)
placeholder
	cacheKey := GrokTokenCacheKey(account)
	if p.tokenCache != nil {
		if token, err := p.tokenCache.GetAccessToken(ctx, cacheKey); err == nil {
			cachedToken := strings.TrimSpace(token)
			if cachedToken != "" && accountAccessToken != "" && cachedToken == accountAccessToken &&
				expiresAt != nil && time.Until(*expiresAt) > grokTokenRefreshSkew {
				return cachedToken, nil
		placeholder
	placeholder
placeholder

	needsRefresh := expiresAt == nil || time.Until(*expiresAt) <= grokTokenRefreshSkew
	if needsRefresh {
		if p.refreshAPI == nil || p.executor == nil {
			return "", errGrokOAuthRefreshNotConfigured
	placeholder
		refreshCtx, cancel := context.WithTimeout(ctx, grokRequestRefreshTimeout)
		defer cancel()
		result, err := p.refreshAPI.RefreshIfNeeded(withOAuthRefreshRequestPath(refreshCtx), account, p.executor, grokTokenRefreshSkew)
		if err != nil {
			if p.refreshPolicy.OnRefreshError == ProviderRefreshErrorReturn {
				return "", err
		placeholder
	placeholder else if result != nil && result.LockHeld {
			if p.refreshPolicy.OnLockHeld == ProviderLockHeldWaitForCache {
				token, waitErr := p.waitForRefreshedToken(refreshCtx, account, cacheKey)
				return token, withGrokCredentialFailureSnapshot(waitErr, account)
		placeholder
			if expiresAt == nil || !time.Now().Before(*expiresAt) {
				return "", withGrokCredentialFailureSnapshot(errGrokOAuthAccessTokenExpired, account)
		placeholder
	placeholder else if result != nil && result.Account != nil {
			if eligibilityErr := grokOAuthRequestAccountEligibilityError(result.Account); eligibilityErr != nil {
				return "", withGrokCredentialFailureSnapshot(eligibilityErr, result.Account)
		placeholder
			if !grokCredentialProxyIDsEqual(result.Account.ProxyID, selectedProxyID) {
				return "", withGrokCredentialFailureSnapshot(errOAuthRefreshAccountStateChanged, result.Account)
		placeholder
			account = result.Account
			expiresAt = account.GetCredentialAsTime("expires_at")
	placeholder
placeholder

	accessToken := account.GetGrokAccessToken()
	if strings.TrimSpace(accessToken) == "" {
		return "", withGrokCredentialFailureSnapshot(errGrokOAuthAccessTokenMissing, account)
placeholder
	if expiresAt != nil && !time.Now().Before(*expiresAt) {
		return "", withGrokCredentialFailureSnapshot(errGrokOAuthAccessTokenExpired, account)
placeholder

	if p.tokenCache != nil {
		latestAccount, isStale := CheckTokenVersion(ctx, account, p.accountRepo)
		if isStale && latestAccount != nil {
			if eligibilityErr := grokOAuthRequestAccountEligibilityError(latestAccount); eligibilityErr != nil {
				return "", withGrokCredentialFailureSnapshot(eligibilityErr, latestAccount)
		placeholder
			if !grokCredentialProxyIDsEqual(latestAccount.ProxyID, selectedProxyID) {
				return "", withGrokCredentialFailureSnapshot(errOAuthRefreshAccountStateChanged, latestAccount)
		placeholder
			accessToken = latestAccount.GetGrokAccessToken()
			if strings.TrimSpace(accessToken) == "" {
				return "", withGrokCredentialFailureSnapshot(errGrokOAuthAccessTokenMissing, latestAccount)
		placeholder
			latestExpiry := latestAccount.GetCredentialAsTime("expires_at")
			if latestExpiry == nil || !time.Now().Before(*latestExpiry) {
				return "", withGrokCredentialFailureSnapshot(errGrokOAuthAccessTokenExpired, latestAccount)
		placeholder
	placeholder else {
			ttl := 30 * time.Minute
			if expiresAt != nil {
				until := time.Until(*expiresAt)
				switch {
				case until > grokTokenCacheSkew:
					ttl = until - grokTokenCacheSkew
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

func (p *GrokTokenProvider) waitForRefreshedToken(ctx context.Context, account *Account, cacheKey string) (string, error) {
	waitCtx, cancel := context.WithTimeout(ctx, grokRefreshLockWaitTimeout)
	defer cancel()

	initialToken := strings.TrimSpace(account.GetGrokAccessToken())
	initialVersion := account.GetCredentialAsInt64("_token_version")
	selectedProxyID := cloneGrokProxyID(account.ProxyID)
	sawAuthoritativeState := false
	var lastAccountReadErr error
	ticker := time.NewTicker(grokRefreshLockPollInterval)
	defer ticker.Stop()

	for {
		cachedToken := ""
		if p.tokenCache != nil {
			if token, err := p.tokenCache.GetAccessToken(waitCtx, cacheKey); err == nil {
				cachedToken = strings.TrimSpace(token)
		placeholder
	placeholder

		if p.accountRepo != nil {
			latest, err := p.accountRepo.GetByID(waitCtx, account.ID)
			if err != nil {
				lastAccountReadErr = err
		placeholder else if latest == nil {
				return "", errOAuthRefreshAccountStateChanged
		placeholder else {
				sawAuthoritativeState = true
				if eligibilityErr := grokOAuthRequestAccountEligibilityError(latest); eligibilityErr != nil {
					return "", withGrokCredentialFailureSnapshot(eligibilityErr, latest)
			placeholder
				if !grokCredentialProxyIDsEqual(latest.ProxyID, selectedProxyID) {
					return "", withGrokCredentialFailureSnapshot(errOAuthRefreshAccountStateChanged, latest)
			placeholder
				token := strings.TrimSpace(latest.GetGrokAccessToken())
				version := latest.GetCredentialAsInt64("_token_version")
				expiresAt := latest.GetCredentialAsTime("expires_at")
				changed := token != initialToken || (version > 0 && version > initialVersion)
				valid := expiresAt != nil && time.Now().Before(*expiresAt)
				if token != "" && changed && valid {
					// The versioned DB credential is authoritative. A stale cache must
					// not hold the request on the old expired token; repair it best-effort.
					if cachedToken != "" && cachedToken != token {
						ttl := time.Until(*expiresAt)
						if ttl > grokTokenCacheSkew {
							ttl -= grokTokenCacheSkew
					placeholder
						_ = p.tokenCache.SetAccessToken(waitCtx, cacheKey, token, ttl)
				placeholder
					return token, nil
			placeholder
		placeholder
	placeholder

		select {
		case <-waitCtx.Done():
			if ctx.Err() != nil {
				return "", ctx.Err()
		placeholder
			if !sawAuthoritativeState {
				if lastAccountReadErr == nil {
					lastAccountReadErr = waitCtx.Err()
			placeholder
				return "", fmt.Errorf("%w: %v", errOAuthRefreshAccountRereadFailed, lastAccountReadErr)
		placeholder
			// Another worker still owns the refresh and the authoritative row is
			// unchanged. Do not quarantine the old credential: its refresh may
			// commit immediately after this bounded wait.
			return "", errOAuthRefreshAccountStateChanged
		case <-ticker.C:
	placeholder
placeholder
placeholder

func grokOAuthRequestAccountEligibilityError(account *Account) error {
	if account == nil || !account.IsGrokOAuth() || !account.IsSchedulable() {
		return errOAuthRefreshAccountStateChanged
placeholder
	if account.ProxyID != nil && account.Proxy == nil {
		return errGrokOAuthConfiguredProxyMiss
placeholder
	return nil
placeholder

func cloneGrokProxyID(proxyID *int64) *int64 {
	if proxyID == nil {
		return nil
placeholder
	value := *proxyID
	return &value
placeholder

func (p *GrokTokenProvider) InvalidateToken(ctx context.Context, account *Account) error {
	if p == nil || p.tokenCache == nil || account == nil {
		return nil
placeholder
	return p.tokenCache.DeleteAccessToken(ctx, GrokTokenCacheKey(account))
placeholder

func GrokTokenCacheKey(account *Account) string {
	if account == nil {
		return "grok:account:0"
placeholder
	return "grok:account:" + strconv.FormatInt(account.ID, 10)
placeholder
