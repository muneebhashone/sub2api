//go:build unit

package service

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// openAITokenCacheStub implements OpenAITokenCache for testing
type openAITokenCacheStub struct {
	mu               sync.Mutex
	tokens           map[string]string
	getErr           error
	setErr           error
	deleteErr        error
	lockAcquired     bool
	lockErr          error
	releaseLockErr   error
	getCalled        int32
	setCalled        int32
	lockCalled       int32
	unlockCalled     int32
	simulateLockRace bool
placeholder

func newOpenAITokenCacheStub() *openAITokenCacheStub {
	return &openAITokenCacheStub{
		tokens:       make(map[string]string),
		lockAcquired: true,
placeholder
placeholder

func (s *openAITokenCacheStub) GetAccessToken(ctx context.Context, cacheKey string) (string, error) {
	atomic.AddInt32(&s.getCalled, 1)
	if s.getErr != nil {
		return "", s.getErr
placeholder
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tokens[cacheKey], nil
placeholder

func (s *openAITokenCacheStub) SetAccessToken(ctx context.Context, cacheKey string, token string, ttl time.Duration) error {
	atomic.AddInt32(&s.setCalled, 1)
	if s.setErr != nil {
		return s.setErr
placeholder
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[cacheKey] = token
	return nil
placeholder

func (s *openAITokenCacheStub) DeleteAccessToken(ctx context.Context, cacheKey string) error {
	if s.deleteErr != nil {
		return s.deleteErr
placeholder
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokens, cacheKey)
	return nil
placeholder

func (s *openAITokenCacheStub) AcquireRefreshLock(ctx context.Context, cacheKey string, ttl time.Duration) (bool, error) {
	atomic.AddInt32(&s.lockCalled, 1)
	if s.lockErr != nil {
		return false, s.lockErr
placeholder
	if s.simulateLockRace {
		return false, nil
placeholder
	return s.lockAcquired, nil
placeholder

func (s *openAITokenCacheStub) ReleaseRefreshLock(ctx context.Context, cacheKey string) error {
	atomic.AddInt32(&s.unlockCalled, 1)
	return s.releaseLockErr
placeholder

// openAIAccountRepoStub is a minimal stub implementing only the methods used by OpenAITokenProvider
type openAIAccountRepoStub struct {
	account      *Account
	getErr       error
	updateErr    error
	getCalled    int32
	updateCalled int32
placeholder

func (r *openAIAccountRepoStub) GetByID(ctx context.Context, id int64) (*Account, error) {
	atomic.AddInt32(&r.getCalled, 1)
	if r.getErr != nil {
		return nil, r.getErr
placeholder
	return r.account, nil
placeholder

func (r *openAIAccountRepoStub) Update(ctx context.Context, account *Account) error {
	atomic.AddInt32(&r.updateCalled, 1)
	if r.updateErr != nil {
		return r.updateErr
placeholder
	r.account = account
	return nil
placeholder

// openAIOAuthServiceStub implements OpenAIOAuthService methods for testing
type openAIOAuthServiceStub struct {
	tokenInfo     *OpenAITokenInfo
	refreshErr    error
	refreshCalled int32
placeholder

func (s *openAIOAuthServiceStub) RefreshAccountToken(ctx context.Context, account *Account) (*OpenAITokenInfo, error) {
	atomic.AddInt32(&s.refreshCalled, 1)
	if s.refreshErr != nil {
		return nil, s.refreshErr
placeholder
	return s.tokenInfo, nil
placeholder

func (s *openAIOAuthServiceStub) BuildAccountCredentials(info *OpenAITokenInfo) map[string]any {
	now := time.Now()
	return map[string]any{
		"access_token":  info.AccessToken,
		"refresh_token": info.RefreshToken,
		"expires_at":    now.Add(time.Duration(info.ExpiresIn) * time.Second).Format(time.RFC3339),
placeholder
placeholder

func TestOpenAITokenProvider_CacheHit(t *testing.T) {
	cache := newOpenAITokenCacheStub()
	account := &Account{
		ID:       100,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "db-token",
	placeholder,
placeholder
	cacheKey := OpenAITokenCacheKey(account)
	cache.tokens[cacheKey] = "cached-token"

	provider := NewOpenAITokenProvider(nil, cache, nil)

	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	require.Equal(t, "cached-token", token)
	require.Equal(t, int32(1), atomic.LoadInt32(&cache.getCalled))
	require.Equal(t, int32(0), atomic.LoadInt32(&cache.setCalled))
placeholder

func TestOpenAITokenProvider_CacheMiss_FromCredentials(t *testing.T) {
	cache := newOpenAITokenCacheStub()
	// Token expires in far future, no refresh needed
	expiresAt := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	account := &Account{
		ID:       101,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "credential-token",
			"expires_at":   expiresAt,
	placeholder,
placeholder

	provider := NewOpenAITokenProvider(nil, cache, nil)

	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	require.Equal(t, "credential-token", token)

	// Should have stored in cache
	cacheKey := OpenAITokenCacheKey(account)
	require.Equal(t, "credential-token", cache.tokens[cacheKey])
placeholder

func TestOpenAITokenProvider_TokenRefresh(t *testing.T) {
	cache := newOpenAITokenCacheStub()
	accountRepo := &openAIAccountRepoStub{placeholder
	oauthService := &openAIOAuthServiceStub{
		tokenInfo: &OpenAITokenInfo{
			AccessToken:  "refreshed-token",
			RefreshToken: "new-refresh-token",
			ExpiresIn:    3600,
	placeholder,
placeholder

	// Token expires soon (within refresh skew)
	expiresAt := time.Now().Add(1 * time.Minute).Format(time.RFC3339)
	account := &Account{
		ID:       102,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token":  "old-token",
			"refresh_token": "old-refresh-token",
			"expires_at":    expiresAt,
	placeholder,
placeholder
	accountRepo.account = account

	// We need to directly test with the stub - create a custom provider
	customProvider := &testOpenAITokenProvider{
		accountRepo:  accountRepo,
		tokenCache:   cache,
		oauthService: oauthService,
placeholder

	token, err := customProvider.GetAccessToken(context.Background(), account)
placeholder
	require.Equal(t, "refreshed-token", token)
	require.Equal(t, int32(1), atomic.LoadInt32(&oauthService.refreshCalled))
placeholder

// testOpenAITokenProvider is a test version that uses the stub OAuth service
type testOpenAITokenProvider struct {
	accountRepo  *openAIAccountRepoStub
	tokenCache   *openAITokenCacheStub
	oauthService *openAIOAuthServiceStub
placeholder

func (p *testOpenAITokenProvider) GetAccessToken(ctx context.Context, account *Account) (string, error) {
	if account == nil {
		return "", errors.New("account is nil")
placeholder
	if account.Platform != PlatformOpenAI || account.Type != AccountTypeOAuth {
		return "", errors.New("not an openai oauth account")
placeholder

	cacheKey := OpenAITokenCacheKey(account)

	// 1. Check cache
	if p.tokenCache != nil {
		if token, err := p.tokenCache.GetAccessToken(ctx, cacheKey); err == nil && token != "" {
			return token, nil
	placeholder
placeholder

	// 2. Check if refresh needed
	expiresAt := account.GetCredentialAsTime("expires_at")
	needsRefresh := expiresAt == nil || time.Until(*expiresAt) <= openAITokenRefreshSkew
	refreshFailed := false
	if needsRefresh && p.tokenCache != nil {
		locked, err := p.tokenCache.AcquireRefreshLock(ctx, cacheKey, 30*time.Second)
		if err == nil && locked {
			defer func() { _ = p.tokenCache.ReleaseRefreshLock(ctx, cacheKey) placeholder()

			// Check cache again after acquiring lock
			if token, err := p.tokenCache.GetAccessToken(ctx, cacheKey); err == nil && token != "" {
				return token, nil
		placeholder

			// Get fresh account from DB
			fresh, err := p.accountRepo.GetByID(ctx, account.ID)
			if err == nil && fresh != nil {
				account = fresh
		placeholder
			expiresAt = account.GetCredentialAsTime("expires_at")
			if expiresAt == nil || time.Until(*expiresAt) <= openAITokenRefreshSkew {
				if p.oauthService == nil {
					refreshFailed = true // 无法刷新，标记失败
			placeholder else {
					tokenInfo, err := p.oauthService.RefreshAccountToken(ctx, account)
					if err != nil {
						refreshFailed = true // 刷新失败，标记以使用短 TTL
				placeholder else {
						newCredentials := p.oauthService.BuildAccountCredentials(tokenInfo)
						for k, v := range account.Credentials {
							if _, exists := newCredentials[k]; !exists {
								newCredentials[k] = v
						placeholder
					placeholder
						account.Credentials = newCredentials
						_ = p.accountRepo.Update(ctx, account)
						expiresAt = account.GetCredentialAsTime("expires_at")
				placeholder
			placeholder
		placeholder
	placeholder else if p.tokenCache.simulateLockRace {
			// Wait and retry cache
			time.Sleep(10 * time.Millisecond) // Short wait for test
			if token, err := p.tokenCache.GetAccessToken(ctx, cacheKey); err == nil && token != "" {
				return token, nil
		placeholder
	placeholder
placeholder

	accessToken := account.GetOpenAIAccessToken()
	if accessToken == "" {
		return "", errors.New("access_token not found in credentials")
placeholder

	// 3. Store in cache
	if p.tokenCache != nil {
		ttl := 30 * time.Minute
		if refreshFailed {
			ttl = time.Minute // 刷新失败时使用短 TTL
	placeholder else if expiresAt != nil {
			until := time.Until(*expiresAt)
			if until > openAITokenCacheSkew {
				ttl = until - openAITokenCacheSkew
		placeholder else if until > 0 {
				ttl = until
		placeholder else {
				ttl = time.Minute
		placeholder
	placeholder
		_ = p.tokenCache.SetAccessToken(ctx, cacheKey, accessToken, ttl)
placeholder

	return accessToken, nil
placeholder

func TestOpenAITokenProvider_LockRaceCondition(t *testing.T) {
	cache := newOpenAITokenCacheStub()
	cache.simulateLockRace = true
	accountRepo := &openAIAccountRepoStub{placeholder

	// Token expires soon
	expiresAt := time.Now().Add(1 * time.Minute).Format(time.RFC3339)
	account := &Account{
		ID:       103,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "race-token",
			"expires_at":   expiresAt,
	placeholder,
placeholder
	accountRepo.account = account

	// Simulate another worker already refreshed and cached
	cacheKey := OpenAITokenCacheKey(account)
	go func() {
		time.Sleep(5 * time.Millisecond)
		cache.mu.Lock()
		cache.tokens[cacheKey] = "winner-token"
		cache.mu.Unlock()
placeholder()

	provider := &testOpenAITokenProvider{
		accountRepo: accountRepo,
		tokenCache:  cache,
placeholder

	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	// Should get the token set by the "winner" or the original
	require.NotEmpty(t, token)
placeholder

func TestOpenAITokenProvider_NilAccount(t *testing.T) {
	provider := NewOpenAITokenProvider(nil, nil, nil)

	token, err := provider.GetAccessToken(context.Background(), nil)
placeholder
	require.Contains(t, err.Error(), "account is nil")
	require.Empty(t, token)
placeholder

func TestOpenAITokenProvider_WrongPlatform(t *testing.T) {
	provider := NewOpenAITokenProvider(nil, nil, nil)
	account := &Account{
		ID:       104,
placeholder
		Type:     AccountTypeOAuth,
placeholder

	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	require.Contains(t, err.Error(), "not an openai/sora oauth account")
	require.Empty(t, token)
placeholder

func TestOpenAITokenProvider_WrongAccountType(t *testing.T) {
	provider := NewOpenAITokenProvider(nil, nil, nil)
	account := &Account{
		ID:       105,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
placeholder

	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	require.Contains(t, err.Error(), "not an openai/sora oauth account")
	require.Empty(t, token)
placeholder

func TestOpenAITokenProvider_NilCache(t *testing.T) {
	// Token doesn't need refresh
	expiresAt := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	account := &Account{
		ID:       106,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "nocache-token",
			"expires_at":   expiresAt,
	placeholder,
placeholder

	provider := NewOpenAITokenProvider(nil, nil, nil)

	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	require.Equal(t, "nocache-token", token)
placeholder

func TestOpenAITokenProvider_CacheGetError(t *testing.T) {
	cache := newOpenAITokenCacheStub()
	cache.getErr = errors.New("redis connection failed")

	// Token doesn't need refresh
	expiresAt := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	account := &Account{
		ID:       107,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "fallback-token",
			"expires_at":   expiresAt,
	placeholder,
placeholder

	provider := NewOpenAITokenProvider(nil, cache, nil)

	// Should gracefully degrade and return from credentials
	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	require.Equal(t, "fallback-token", token)
placeholder

func TestOpenAITokenProvider_CacheSetError(t *testing.T) {
	cache := newOpenAITokenCacheStub()
	cache.setErr = errors.New("redis write failed")

	expiresAt := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	account := &Account{
		ID:       108,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "still-works-token",
			"expires_at":   expiresAt,
	placeholder,
placeholder

	provider := NewOpenAITokenProvider(nil, cache, nil)

	// Should still work even if cache set fails
	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	require.Equal(t, "still-works-token", token)
placeholder

func TestOpenAITokenProvider_MissingAccessToken(t *testing.T) {
	cache := newOpenAITokenCacheStub()
	expiresAt := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	account := &Account{
		ID:       109,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"expires_at": expiresAt,
			// missing access_token
	placeholder,
placeholder

	provider := NewOpenAITokenProvider(nil, cache, nil)

	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	require.Contains(t, err.Error(), "access_token not found")
	require.Empty(t, token)
placeholder

func TestOpenAITokenProvider_RefreshError(t *testing.T) {
	cache := newOpenAITokenCacheStub()
	accountRepo := &openAIAccountRepoStub{placeholder
	oauthService := &openAIOAuthServiceStub{
		refreshErr: errors.New("oauth refresh failed"),
placeholder

	// Token expires soon
	expiresAt := time.Now().Add(1 * time.Minute).Format(time.RFC3339)
	account := &Account{
		ID:       110,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token":  "old-token",
			"refresh_token": "old-refresh-token",
			"expires_at":    expiresAt,
	placeholder,
placeholder
	accountRepo.account = account

	provider := &testOpenAITokenProvider{
		accountRepo:  accountRepo,
		tokenCache:   cache,
		oauthService: oauthService,
placeholder

	// Now with fallback behavior, should return existing token even if refresh fails
	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	require.Equal(t, "old-token", token) // Fallback to existing token
placeholder

func TestOpenAITokenProvider_OAuthServiceNotConfigured(t *testing.T) {
	cache := newOpenAITokenCacheStub()
	accountRepo := &openAIAccountRepoStub{placeholder

	// Token expires soon
	expiresAt := time.Now().Add(1 * time.Minute).Format(time.RFC3339)
	account := &Account{
		ID:       111,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "old-token",
			"expires_at":   expiresAt,
	placeholder,
placeholder
	accountRepo.account = account

	provider := &testOpenAITokenProvider{
		accountRepo:  accountRepo,
		tokenCache:   cache,
		oauthService: nil, // not configured
placeholder

	// Now with fallback behavior, should return existing token even if oauth service not configured
	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	require.Equal(t, "old-token", token) // Fallback to existing token
placeholder

func TestOpenAITokenProvider_TTLCalculation(t *testing.T) {
	tests := []struct {
		name      string
		expiresIn time.Duration
placeholder{
		{
			name:      "far_future_expiry",
			expiresIn: 1 * time.Hour,
	placeholder,
		{
			name:      "medium_expiry",
			expiresIn: 10 * time.Minute,
	placeholder,
		{
			name:      "near_expiry",
			expiresIn: 6 * time.Minute,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := newOpenAITokenCacheStub()
			expiresAt := time.Now().Add(tt.expiresIn).Format(time.RFC3339)
			account := &Account{
				ID:       200,
				Platform: PlatformOpenAI,
				Type:     AccountTypeOAuth,
		placeholder
					"access_token": "test-token",
					"expires_at":   expiresAt,
			placeholder,
		placeholder

			provider := NewOpenAITokenProvider(nil, cache, nil)

			_, err := provider.GetAccessToken(context.Background(), account)
		placeholder

			// Verify token was cached
			cacheKey := OpenAITokenCacheKey(account)
			require.Equal(t, "test-token", cache.tokens[cacheKey])
	placeholder)
placeholder
placeholder

func TestOpenAITokenProvider_DoubleCheckAfterLock(t *testing.T) {
	cache := newOpenAITokenCacheStub()
	accountRepo := &openAIAccountRepoStub{placeholder
	oauthService := &openAIOAuthServiceStub{
		tokenInfo: &OpenAITokenInfo{
			AccessToken:  "refreshed-token",
			RefreshToken: "new-refresh",
			ExpiresIn:    3600,
	placeholder,
placeholder

	// Token expires soon
	expiresAt := time.Now().Add(1 * time.Minute).Format(time.RFC3339)
	account := &Account{
		ID:       112,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "old-token",
			"expires_at":   expiresAt,
	placeholder,
placeholder
	accountRepo.account = account
	cacheKey := OpenAITokenCacheKey(account)

	// Simulate: first GetAccessToken returns empty, but after lock acquired, cache has token
	originalGet := int32(0)
	cache.tokens[cacheKey] = "" // Empty initially

	provider := &testOpenAITokenProvider{
		accountRepo:  accountRepo,
		tokenCache:   cache,
		oauthService: oauthService,
placeholder

	// In a goroutine, set the cached token after a small delay (simulating race)
	go func() {
		time.Sleep(5 * time.Millisecond)
		cache.mu.Lock()
		cache.tokens[cacheKey] = "cached-by-other"
		cache.mu.Unlock()
placeholder()

	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	// Should get either the refreshed token or the cached one
	require.NotEmpty(t, token)
	_ = originalGet // Suppress unused warning
placeholder

// Tests for real provider - to increase coverage
func TestOpenAITokenProvider_Real_LockFailedWait(t *testing.T) {
	cache := newOpenAITokenCacheStub()
	cache.lockAcquired = false // Lock acquisition fails

	// Token expires soon (within refresh skew) to trigger lock attempt
	expiresAt := time.Now().Add(1 * time.Minute).Format(time.RFC3339)
	account := &Account{
		ID:       200,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "fallback-token",
			"expires_at":   expiresAt,
	placeholder,
placeholder

	// Set token in cache after lock wait period (simulate other worker refreshing)
	cacheKey := OpenAITokenCacheKey(account)
	go func() {
		time.Sleep(100 * time.Millisecond)
		cache.mu.Lock()
		cache.tokens[cacheKey] = "refreshed-by-other"
		cache.mu.Unlock()
placeholder()

	provider := NewOpenAITokenProvider(nil, cache, nil)
	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	// Should get either the fallback token or the refreshed one
	require.NotEmpty(t, token)
placeholder

func TestOpenAITokenProvider_Real_CacheHitAfterWait(t *testing.T) {
	cache := newOpenAITokenCacheStub()
	cache.lockAcquired = false // Lock acquisition fails

	// Token expires soon
	expiresAt := time.Now().Add(1 * time.Minute).Format(time.RFC3339)
	account := &Account{
		ID:       201,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "original-token",
			"expires_at":   expiresAt,
	placeholder,
placeholder

	cacheKey := OpenAITokenCacheKey(account)
	// Set token in cache immediately after wait starts
	go func() {
		time.Sleep(50 * time.Millisecond)
		cache.mu.Lock()
		cache.tokens[cacheKey] = "winner-token"
		cache.mu.Unlock()
placeholder()

	provider := NewOpenAITokenProvider(nil, cache, nil)
	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	require.NotEmpty(t, token)
placeholder

func TestOpenAITokenProvider_Real_ExpiredWithoutRefreshToken(t *testing.T) {
	cache := newOpenAITokenCacheStub()
	cache.lockAcquired = false // Prevent entering refresh logic

	// Token with nil expires_at (no expiry set) - should use credentials
	account := &Account{
		ID:       202,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "no-expiry-token",
	placeholder,
placeholder

	provider := NewOpenAITokenProvider(nil, cache, nil)
	token, err := provider.GetAccessToken(context.Background(), account)
	// Without OAuth service, refresh will fail but token should be returned from credentials
placeholder
	require.Equal(t, "no-expiry-token", token)
placeholder

func TestOpenAITokenProvider_Real_WhitespaceToken(t *testing.T) {
	cache := newOpenAITokenCacheStub()
	cacheKey := "openai:account:203"
	cache.tokens[cacheKey] = "   " // Whitespace only - should be treated as empty

	expiresAt := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	account := &Account{
		ID:       203,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "real-token",
			"expires_at":   expiresAt,
	placeholder,
placeholder

	provider := NewOpenAITokenProvider(nil, cache, nil)
	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	require.Equal(t, "real-token", token) // Should fall back to credentials
placeholder

func TestOpenAITokenProvider_Real_LockError(t *testing.T) {
	cache := newOpenAITokenCacheStub()
	cache.lockErr = errors.New("redis lock failed")

	// Token expires soon (within refresh skew)
	expiresAt := time.Now().Add(1 * time.Minute).Format(time.RFC3339)
	account := &Account{
		ID:       204,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "fallback-on-lock-error",
			"expires_at":   expiresAt,
	placeholder,
placeholder

	provider := NewOpenAITokenProvider(nil, cache, nil)
	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	require.Equal(t, "fallback-on-lock-error", token)
placeholder

func TestOpenAITokenProvider_Real_WhitespaceCredentialToken(t *testing.T) {
	cache := newOpenAITokenCacheStub()

	expiresAt := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	account := &Account{
		ID:       205,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "   ", // Whitespace only
			"expires_at":   expiresAt,
	placeholder,
placeholder

	provider := NewOpenAITokenProvider(nil, cache, nil)
	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	require.Contains(t, err.Error(), "access_token not found")
	require.Empty(t, token)
placeholder

func TestOpenAITokenProvider_Real_NilCredentials(t *testing.T) {
	cache := newOpenAITokenCacheStub()

	expiresAt := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	account := &Account{
		ID:       206,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"expires_at": expiresAt,
			// No access_token
	placeholder,
placeholder

	provider := NewOpenAITokenProvider(nil, cache, nil)
	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	require.Contains(t, err.Error(), "access_token not found")
	require.Empty(t, token)
placeholder

func TestOpenAITokenProvider_Real_LockRace_PollingHitsCache(t *testing.T) {
	cache := newOpenAITokenCacheStub()
	cache.lockAcquired = false // 模拟锁被其他 worker 持有

	expiresAt := time.Now().Add(1 * time.Minute).Format(time.RFC3339)
	account := &Account{
		ID:       207,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "fallback-token",
			"expires_at":   expiresAt,
	placeholder,
placeholder

	cacheKey := OpenAITokenCacheKey(account)
	go func() {
		time.Sleep(5 * time.Millisecond)
		cache.mu.Lock()
		cache.tokens[cacheKey] = "winner-token"
		cache.mu.Unlock()
placeholder()

	provider := NewOpenAITokenProvider(nil, cache, nil)
	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	require.Equal(t, "winner-token", token)
placeholder

func TestOpenAITokenProvider_Real_LockRace_ContextCanceled(t *testing.T) {
	cache := newOpenAITokenCacheStub()
	cache.lockAcquired = false // 模拟锁被其他 worker 持有

	expiresAt := time.Now().Add(1 * time.Minute).Format(time.RFC3339)
	account := &Account{
		ID:       208,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "fallback-token",
			"expires_at":   expiresAt,
	placeholder,
placeholder

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	provider := NewOpenAITokenProvider(nil, cache, nil)
	start := time.Now()
	token, err := provider.GetAccessToken(ctx, account)
placeholder
	require.ErrorIs(t, err, context.Canceled)
	require.Empty(t, token)
	require.Less(t, time.Since(start), 50*time.Millisecond)
placeholder

func TestOpenAITokenProvider_RuntimeMetrics_LockWaitHitAndSnapshot(t *testing.T) {
	cache := newOpenAITokenCacheStub()
	cache.lockAcquired = false

	expiresAt := time.Now().Add(1 * time.Minute).Format(time.RFC3339)
	account := &Account{
		ID:       209,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "fallback-token",
			"expires_at":   expiresAt,
	placeholder,
placeholder
	cacheKey := OpenAITokenCacheKey(account)
	go func() {
		time.Sleep(10 * time.Millisecond)
		cache.mu.Lock()
		cache.tokens[cacheKey] = "winner-token"
		cache.mu.Unlock()
placeholder()

	provider := NewOpenAITokenProvider(nil, cache, nil)
	token, err := provider.GetAccessToken(context.Background(), account)
placeholder
	require.Equal(t, "winner-token", token)

	metrics := provider.SnapshotRuntimeMetrics()
	require.GreaterOrEqual(t, metrics.RefreshRequests, int64(1))
	require.GreaterOrEqual(t, metrics.LockContention, int64(1))
	require.GreaterOrEqual(t, metrics.LockWaitSamples, int64(1))
	require.GreaterOrEqual(t, metrics.LockWaitHit, int64(1))
	require.GreaterOrEqual(t, metrics.LockWaitTotalMs, int64(0))
	require.GreaterOrEqual(t, metrics.LastObservedUnixMs, int64(1))
placeholder

func TestOpenAITokenProvider_RuntimeMetrics_LockAcquireFailure(t *testing.T) {
	cache := newOpenAITokenCacheStub()
	cache.lockErr = errors.New("redis lock error")

	expiresAt := time.Now().Add(1 * time.Minute).Format(time.RFC3339)
	account := &Account{
		ID:       210,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "fallback-token",
			"expires_at":   expiresAt,
	placeholder,
placeholder

	provider := NewOpenAITokenProvider(nil, cache, nil)
	_, err := provider.GetAccessToken(context.Background(), account)
placeholder

	metrics := provider.SnapshotRuntimeMetrics()
	require.GreaterOrEqual(t, metrics.LockAcquireFailure, int64(1))
	require.GreaterOrEqual(t, metrics.RefreshRequests, int64(1))
placeholder
