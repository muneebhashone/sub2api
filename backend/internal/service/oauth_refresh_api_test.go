//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// ---------- mock helpers ----------

// refreshAPIAccountRepo implements AccountRepository for OAuthRefreshAPI tests.
type refreshAPIAccountRepo struct {
	mockAccountRepoForGemini
	account                *Account // returned by GetByID
	getByIDErr             error
	updateErr              error
	updateCalls            int
	updateCredentialsCalls int
placeholder

func (r *refreshAPIAccountRepo) GetByID(_ context.Context, _ int64) (*Account, error) {
	if r.getByIDErr != nil {
		return nil, r.getByIDErr
placeholder
	return r.account, nil
placeholder

func (r *refreshAPIAccountRepo) Update(_ context.Context, _ *Account) error {
	r.updateCalls++
	return r.updateErr
placeholder

func (r *refreshAPIAccountRepo) UpdateCredentials(_ context.Context, id int64, credentials map[string]any) error {
	r.updateCalls++
	r.updateCredentialsCalls++
	if r.updateErr != nil {
		return r.updateErr
placeholder
	if r.account == nil || r.account.ID != id {
		r.account = &Account{ID: idplaceholder
placeholder
	r.account.Credentials = cloneCredentials(credentials)
	return nil
placeholder

// refreshAPIExecutorStub implements OAuthRefreshExecutor for tests.
type refreshAPIExecutorStub struct {
	needsRefresh bool
	credentials  map[string]any
	err          error
	refreshCalls int
placeholder

func (e *refreshAPIExecutorStub) CanRefresh(_ *Account) bool { return true placeholder

func (e *refreshAPIExecutorStub) NeedsRefresh(_ *Account, _ time.Duration) bool {
	return e.needsRefresh
placeholder

func (e *refreshAPIExecutorStub) Refresh(_ context.Context, _ *Account) (map[string]any, error) {
	e.refreshCalls++
	if e.err != nil {
		return nil, e.err
placeholder
	return e.credentials, nil
placeholder

func (e *refreshAPIExecutorStub) CacheKey(account *Account) string {
	return "test:api:" + account.Platform
placeholder

// refreshAPICacheStub implements GeminiTokenCache for OAuthRefreshAPI tests.
type refreshAPICacheStub struct {
	lockResult   bool
	lockErr      error
	releaseCalls int
placeholder

func (c *refreshAPICacheStub) GetAccessToken(context.Context, string) (string, error) {
	return "", nil
placeholder

func (c *refreshAPICacheStub) SetAccessToken(context.Context, string, string, time.Duration) error {
	return nil
placeholder

func (c *refreshAPICacheStub) DeleteAccessToken(context.Context, string) error { return nil placeholder

func (c *refreshAPICacheStub) AcquireRefreshLock(context.Context, string, time.Duration) (bool, error) {
	return c.lockResult, c.lockErr
placeholder

func (c *refreshAPICacheStub) ReleaseRefreshLock(context.Context, string) error {
	c.releaseCalls++
	return nil
placeholder

// ========== RefreshIfNeeded tests ==========

func TestRefreshIfNeeded_Success(t *testing.T) {
	account := &Account{ID: 1, Platform: PlatformAnthropic, Type: AccountTypeOAuthplaceholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		credentials:  map[string]any{"access_token": "new-token"placeholder,
placeholder

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

placeholder
	require.True(t, result.Refreshed)
	require.NotNil(t, result.NewCredentials)
	require.Equal(t, "new-token", result.NewCredentials["access_token"])
	require.NotNil(t, result.NewCredentials["_token_version"]) // version stamp set
	require.Equal(t, 1, repo.updateCalls)                      // DB updated
	require.Equal(t, 1, repo.updateCredentialsCalls)
	require.Equal(t, 1, cache.releaseCalls) // lock released
	require.Equal(t, 1, executor.refreshCalls)
placeholder

func TestRefreshIfNeeded_UpdateCredentialsPreservesRateLimitState(t *testing.T) {
	resetAt := time.Now().Add(45 * time.Minute)
	account := &Account{
		ID:               11,
		Platform:         PlatformGemini,
		Type:             AccountTypeOAuth,
		RateLimitResetAt: &resetAt,
placeholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		credentials:  map[string]any{"access_token": "safe-token"placeholder,
placeholder

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

placeholder
	require.True(t, result.Refreshed)
	require.Equal(t, 1, repo.updateCredentialsCalls)
	require.NotNil(t, repo.account.RateLimitResetAt)
	require.WithinDuration(t, resetAt, *repo.account.RateLimitResetAt, time.Second)
placeholder

func TestRefreshIfNeeded_LockHeld(t *testing.T) {
	account := &Account{ID: 2, Platform: PlatformAnthropicplaceholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	cache := &refreshAPICacheStub{lockResult: falseplaceholder // lock not acquired
	executor := &refreshAPIExecutorStub{needsRefresh: trueplaceholder

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

placeholder
	require.True(t, result.LockHeld)
	require.False(t, result.Refreshed)
	require.Equal(t, 0, repo.updateCalls)
	require.Equal(t, 0, executor.refreshCalls)
placeholder

func TestRefreshIfNeeded_LockErrorDegrades(t *testing.T) {
	account := &Account{ID: 3, Platform: PlatformGemini, Type: AccountTypeOAuthplaceholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	cache := &refreshAPICacheStub{lockErr: errors.New("redis down")placeholder // lock error
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		credentials:  map[string]any{"access_token": "degraded-token"placeholder,
placeholder

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

placeholder
	require.True(t, result.Refreshed)       // still refreshed (degraded mode)
	require.Equal(t, 1, repo.updateCalls)   // DB updated
	require.Equal(t, 0, cache.releaseCalls) // no lock to release
	require.Equal(t, 1, executor.refreshCalls)
placeholder

func TestRefreshIfNeeded_NoCacheNoLock(t *testing.T) {
	account := &Account{ID: 4, Platform: PlatformGemini, Type: AccountTypeOAuthplaceholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		credentials:  map[string]any{"access_token": "no-cache-token"placeholder,
placeholder

	api := NewOAuthRefreshAPI(repo, nil) // no cache = no lock
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

placeholder
	require.True(t, result.Refreshed)
	require.Equal(t, 1, repo.updateCalls)
placeholder

func TestRefreshIfNeeded_AlreadyRefreshed(t *testing.T) {
	account := &Account{ID: 5, Platform: PlatformAnthropicplaceholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	executor := &refreshAPIExecutorStub{needsRefresh: falseplaceholder // already refreshed

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

placeholder
	require.False(t, result.Refreshed)
	require.False(t, result.LockHeld)
	require.NotNil(t, result.Account) // returns fresh account
	require.Equal(t, 0, repo.updateCalls)
	require.Equal(t, 0, executor.refreshCalls)
placeholder

func TestRefreshIfNeeded_RefreshError(t *testing.T) {
	account := &Account{ID: 6, Platform: PlatformAnthropicplaceholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		err:          errors.New("invalid_grant: token revoked"),
placeholder

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

placeholder
	require.Nil(t, result)
	require.Contains(t, err.Error(), "invalid_grant")
	require.Equal(t, 0, repo.updateCalls)   // no DB update on refresh error
	require.Equal(t, 1, cache.releaseCalls) // lock still released via defer
placeholder

func TestRefreshIfNeeded_DBUpdateError(t *testing.T) {
	account := &Account{ID: 7, Platform: PlatformGemini, Type: AccountTypeOAuthplaceholder
	repo := &refreshAPIAccountRepo{
		account:   account,
		updateErr: errors.New("db connection lost"),
placeholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		credentials:  map[string]any{"access_token": "token"placeholder,
placeholder

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

placeholder
	require.Nil(t, result)
	require.Contains(t, err.Error(), "DB update failed")
	require.Equal(t, 1, repo.updateCalls) // attempted
placeholder

func TestRefreshIfNeeded_DBRereadFails(t *testing.T) {
	account := &Account{ID: 8, Platform: PlatformAnthropic, Type: AccountTypeOAuthplaceholder
	repo := &refreshAPIAccountRepo{
		account:    nil, // GetByID returns nil
		getByIDErr: errors.New("db timeout"),
placeholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		credentials:  map[string]any{"access_token": "fallback-token"placeholder,
placeholder

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

placeholder
	require.True(t, result.Refreshed)
	require.Equal(t, 1, executor.refreshCalls) // still refreshes using passed-in account
placeholder

func TestRefreshIfNeeded_NilCredentials(t *testing.T) {
	account := &Account{ID: 9, Platform: PlatformGemini, Type: AccountTypeOAuthplaceholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		credentials:  nil, // Refresh returns nil credentials
placeholder

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

placeholder
	require.True(t, result.Refreshed)
	require.Nil(t, result.NewCredentials)
	require.Equal(t, 0, repo.updateCalls) // no DB update when credentials are nil
placeholder

// ========== MergeCredentials tests ==========

func TestMergeCredentials_Basic(t *testing.T) {
	old := map[string]any{"a": "1", "b": "2", "c": "3"placeholder
	new := map[string]any{"a": "new", "d": "4"placeholder

	result := MergeCredentials(old, new)

	require.Equal(t, "new", result["a"]) // new value preserved
	require.Equal(t, "2", result["b"])   // old value kept
	require.Equal(t, "3", result["c"])   // old value kept
	require.Equal(t, "4", result["d"])   // new value preserved
placeholder

func TestMergeCredentials_NilNew(t *testing.T) {
	old := map[string]any{"a": "1"placeholder

	result := MergeCredentials(old, nil)

	require.NotNil(t, result)
	require.Equal(t, "1", result["a"])
placeholder

func TestMergeCredentials_NilOld(t *testing.T) {
	new := map[string]any{"a": "1"placeholder

	result := MergeCredentials(nil, new)

	require.Equal(t, "1", result["a"])
placeholder

func TestMergeCredentials_BothNil(t *testing.T) {
	result := MergeCredentials(nil, nil)
	require.NotNil(t, result)
	require.Empty(t, result)
placeholder

func TestMergeCredentials_NewOverridesOld(t *testing.T) {
	old := map[string]any{"access_token": "old-token", "refresh_token": "old-refresh"placeholder
	new := map[string]any{"access_token": "new-token"placeholder

	result := MergeCredentials(old, new)

	require.Equal(t, "new-token", result["access_token"])    // overridden
	require.Equal(t, "old-refresh", result["refresh_token"]) // preserved
placeholder

// ========== BuildClaudeAccountCredentials tests ==========

func TestBuildClaudeAccountCredentials_Full(t *testing.T) {
	tokenInfo := &TokenInfo{
		AccessToken:  "at-123",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		ExpiresAt:    1700000000,
		RefreshToken: "rt-456",
		Scope:        "openid",
placeholder

	creds := BuildClaudeAccountCredentials(tokenInfo)

	require.Equal(t, "at-123", creds["access_token"])
	require.Equal(t, "Bearer", creds["token_type"])
	require.Equal(t, "3600", creds["expires_in"])
	require.Equal(t, "1700000000", creds["expires_at"])
	require.Equal(t, "rt-456", creds["refresh_token"])
	require.Equal(t, "openid", creds["scope"])
placeholder

func TestBuildClaudeAccountCredentials_Minimal(t *testing.T) {
	tokenInfo := &TokenInfo{
		AccessToken: "at-789",
		TokenType:   "Bearer",
		ExpiresIn:   7200,
		ExpiresAt:   1700003600,
placeholder

	creds := BuildClaudeAccountCredentials(tokenInfo)

	require.Equal(t, "at-789", creds["access_token"])
	require.Equal(t, "Bearer", creds["token_type"])
	require.Equal(t, "7200", creds["expires_in"])
	require.Equal(t, "1700003600", creds["expires_at"])
	_, hasRefresh := creds["refresh_token"]
	_, hasScope := creds["scope"]
	require.False(t, hasRefresh, "refresh_token should not be set when empty")
	require.False(t, hasScope, "scope should not be set when empty")
placeholder

// ========== BackgroundRefreshPolicy tests ==========

func TestBackgroundRefreshPolicy_DefaultSkips(t *testing.T) {
	p := DefaultBackgroundRefreshPolicy()

	require.ErrorIs(t, p.handleLockHeld(), errRefreshSkipped)
	require.ErrorIs(t, p.handleAlreadyRefreshed(), errRefreshSkipped)
placeholder

func TestBackgroundRefreshPolicy_SuccessOverride(t *testing.T) {
	p := BackgroundRefreshPolicy{
		OnLockHeld:       BackgroundSkipAsSuccess,
		OnAlreadyRefresh: BackgroundSkipAsSuccess,
placeholder

	require.NoError(t, p.handleLockHeld())
	require.NoError(t, p.handleAlreadyRefreshed())
placeholder

// ========== ProviderRefreshPolicy tests ==========

func TestClaudeProviderRefreshPolicy(t *testing.T) {
	p := ClaudeProviderRefreshPolicy()
	require.Equal(t, ProviderRefreshErrorUseExistingToken, p.OnRefreshError)
	require.Equal(t, ProviderLockHeldWaitForCache, p.OnLockHeld)
	require.Equal(t, time.Minute, p.FailureTTL)
placeholder

func TestOpenAIProviderRefreshPolicy(t *testing.T) {
	p := OpenAIProviderRefreshPolicy()
	require.Equal(t, ProviderRefreshErrorUseExistingToken, p.OnRefreshError)
	require.Equal(t, ProviderLockHeldWaitForCache, p.OnLockHeld)
	require.Equal(t, time.Minute, p.FailureTTL)
placeholder

func TestGeminiProviderRefreshPolicy(t *testing.T) {
	p := GeminiProviderRefreshPolicy()
	require.Equal(t, ProviderRefreshErrorReturn, p.OnRefreshError)
	require.Equal(t, ProviderLockHeldUseExistingToken, p.OnLockHeld)
	require.Equal(t, time.Duration(0), p.FailureTTL)
placeholder

func TestAntigravityProviderRefreshPolicy(t *testing.T) {
	p := AntigravityProviderRefreshPolicy()
	require.Equal(t, ProviderRefreshErrorReturn, p.OnRefreshError)
	require.Equal(t, ProviderLockHeldUseExistingToken, p.OnLockHeld)
	require.Equal(t, time.Duration(0), p.FailureTTL)
placeholder
