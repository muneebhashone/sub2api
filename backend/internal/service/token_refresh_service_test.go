//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type tokenRefreshAccountRepo struct {
	mockAccountRepoForGemini
	updateCalls            int
	fullUpdateCalls        int
	updateCredentialsCalls int
	setErrorCalls          int
	clearTempCalls         int
	setTempUnschedCalls    int
	updateExtraCalls       int
	lastErrorMessage       string
	lastTempUnschedReason  string
	lastExtraUpdates       map[string]any
	lastAccount            *Account
	updateErr              error
	setErrorErr            error
	setTempUnschedErr      error
	beforeConditionalState func()
placeholder

func (r *tokenRefreshAccountRepo) Update(ctx context.Context, account *Account) error {
	r.updateCalls++
	r.fullUpdateCalls++
	r.lastAccount = account
	return r.updateErr
placeholder

func (r *tokenRefreshAccountRepo) UpdateCredentials(ctx context.Context, id int64, credentials map[string]any) error {
	r.updateCalls++
	r.updateCredentialsCalls++
	if r.updateErr != nil {
		return r.updateErr
placeholder
	cloned := shallowCopyMap(credentials)
	if r.accountsByID != nil {
		if acc, ok := r.accountsByID[id]; ok && acc != nil {
			acc.Credentials = cloned
			r.lastAccount = acc
			return nil
	placeholder
placeholder
	r.lastAccount = &Account{ID: id, Credentials: clonedplaceholder
	return nil
placeholder

func (r *tokenRefreshAccountRepo) SetError(ctx context.Context, id int64, errorMsg string) error {
	r.setErrorCalls++
	r.lastErrorMessage = errorMsg
	return r.setErrorErr
placeholder

func (r *tokenRefreshAccountRepo) ClearTempUnschedulable(ctx context.Context, id int64) error {
	r.clearTempCalls++
	return nil
placeholder

func (r *tokenRefreshAccountRepo) SetTempUnschedulable(ctx context.Context, id int64, until time.Time, reason string) error {
	r.setTempUnschedCalls++
	r.lastTempUnschedReason = reason
	return r.setTempUnschedErr
placeholder

func (r *tokenRefreshAccountRepo) SetGrokCredentialErrorIfMatch(
	_ context.Context,
	id int64,
	snapshot GrokCredentialMutationSnapshot,
	errorMsg string,
) (bool, error) {
	if r.beforeConditionalState != nil {
		hook := r.beforeConditionalState
		r.beforeConditionalState = nil
		hook()
placeholder
	account := r.accountsByID[id]
	if !grokCredentialSnapshotMatchesAccount(account, snapshot) ||
		(errorMsg == string(GrokCredentialReasonProxyInvalid) && account.Proxy != nil) {
		return false, nil
placeholder
	r.setErrorCalls++
	r.lastErrorMessage = errorMsg
	if r.setErrorErr != nil {
		return false, r.setErrorErr
placeholder
	account.Status = StatusError
	account.Schedulable = false
	account.ErrorMessage = errorMsg
	return true, nil
placeholder

func (r *tokenRefreshAccountRepo) SetGrokCredentialTempUnschedulableIfMatch(
	_ context.Context,
	id int64,
	snapshot GrokCredentialMutationSnapshot,
	until time.Time,
	reason string,
) (bool, error) {
	if r.beforeConditionalState != nil {
		hook := r.beforeConditionalState
		r.beforeConditionalState = nil
		hook()
placeholder
	account := r.accountsByID[id]
	if !grokCredentialSnapshotMatchesAccount(account, snapshot) {
		return false, nil
placeholder
	r.setTempUnschedCalls++
	r.lastTempUnschedReason = reason
	if r.setTempUnschedErr != nil {
		return false, r.setTempUnschedErr
placeholder
	value := until
	account.TempUnschedulableUntil = &value
	return true, nil
placeholder

func grokCredentialSnapshotMatchesAccount(account *Account, snapshot GrokCredentialMutationSnapshot) bool {
	return account != nil && account.IsGrokOAuth() && account.IsSchedulable() &&
		grokCredentialMutationSnapshot(account).CredentialsJSON == snapshot.CredentialsJSON &&
		grokCredentialProxyIDsEqual(account.ProxyID, snapshot.ProxyID)
placeholder

func (r *tokenRefreshAccountRepo) UpdateExtra(ctx context.Context, id int64, updates map[string]any) error {
	r.updateExtraCalls++
	r.lastExtraUpdates = shallowCopyMap(updates)
	if r.accountsByID != nil {
		if acc, ok := r.accountsByID[id]; ok && acc != nil {
			if acc.Extra == nil {
				acc.Extra = make(map[string]any, len(updates))
		placeholder
			for k, v := range updates {
				acc.Extra[k] = v
		placeholder
	placeholder
placeholder
	return nil
placeholder

type tokenCacheInvalidatorStub struct {
	calls int
	err   error
placeholder

func (s *tokenCacheInvalidatorStub) InvalidateToken(ctx context.Context, account *Account) error {
	s.calls++
	return s.err
placeholder

type tempUnschedCacheStub struct {
	deleteCalls int
	setCalls    int
	lastState   *TempUnschedState
placeholder

func (s *tempUnschedCacheStub) SetTempUnsched(ctx context.Context, accountID int64, state *TempUnschedState) error {
	s.setCalls++
	s.lastState = state
	return nil
placeholder

func (s *tempUnschedCacheStub) GetTempUnsched(ctx context.Context, accountID int64) (*TempUnschedState, error) {
	return nil, nil
placeholder

func (s *tempUnschedCacheStub) DeleteTempUnsched(ctx context.Context, accountID int64) error {
	s.deleteCalls++
	return nil
placeholder

type tokenRefresherStub struct {
	credentials map[string]any
	err         error
placeholder

func (r *tokenRefresherStub) CanRefresh(account *Account) bool {
	return true
placeholder

func (r *tokenRefresherStub) NeedsRefresh(account *Account, refreshWindowDuration time.Duration) bool {
	return true
placeholder

func (r *tokenRefresherStub) Refresh(ctx context.Context, account *Account) (map[string]any, error) {
	if r.err != nil {
		return nil, r.err
placeholder
	return r.credentials, nil
placeholder

func (r *tokenRefresherStub) CacheKey(account *Account) string {
	return "test:stub:" + account.Platform
placeholder

func TestTokenRefreshService_RefreshWithRetry_InvalidatesCache(t *testing.T) {
	repo := &tokenRefreshAccountRepo{placeholder
	invalidator := &tokenCacheInvalidatorStub{placeholder
	cfg := &config.Config{
		TokenRefresh: config.TokenRefreshConfig{
			MaxRetries:          1,
			RetryBackoffSeconds: 0,
	placeholder,
placeholder
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg, nil)
	account := &Account{
		ID:       5,
placeholder
		Type:     AccountTypeOAuth,
placeholder
	refresher := &tokenRefresherStub{
		credentials: map[string]any{
			"access_token": "new-token",
	placeholder,
placeholder

	err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
placeholder
	require.Equal(t, 1, repo.updateCalls)
	require.Equal(t, 1, repo.updateCredentialsCalls)
	require.Equal(t, 0, repo.fullUpdateCalls)
	require.Equal(t, 1, invalidator.calls)
	require.Equal(t, "new-token", account.GetCredential("access_token"))
placeholder

func TestTokenRefreshService_RefreshWithRetry_InvalidatorErrorIgnored(t *testing.T) {
	repo := &tokenRefreshAccountRepo{placeholder
	invalidator := &tokenCacheInvalidatorStub{err: errors.New("invalidate failed")placeholder
	cfg := &config.Config{
		TokenRefresh: config.TokenRefreshConfig{
			MaxRetries:          1,
			RetryBackoffSeconds: 0,
	placeholder,
placeholder
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg, nil)
	account := &Account{
		ID:       6,
placeholder
		Type:     AccountTypeOAuth,
placeholder
	refresher := &tokenRefresherStub{
		credentials: map[string]any{
			"access_token": "token",
	placeholder,
placeholder

	err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
placeholder
	require.Equal(t, 1, repo.updateCalls)
	require.Equal(t, 1, invalidator.calls)
placeholder

func TestTokenRefreshService_RefreshWithRetry_NilInvalidator(t *testing.T) {
	repo := &tokenRefreshAccountRepo{placeholder
	cfg := &config.Config{
		TokenRefresh: config.TokenRefreshConfig{
			MaxRetries:          1,
			RetryBackoffSeconds: 0,
	placeholder,
placeholder
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, nil, nil, cfg, nil)
	account := &Account{
		ID:       7,
placeholder
		Type:     AccountTypeOAuth,
placeholder
	refresher := &tokenRefresherStub{
		credentials: map[string]any{
			"access_token": "token",
	placeholder,
placeholder

	err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
placeholder
	require.Equal(t, 1, repo.updateCalls)
placeholder

// TestTokenRefreshService_RefreshWithRetry_Antigravity 测试 Antigravity 平台的缓存失效
func TestTokenRefreshService_RefreshWithRetry_Antigravity(t *testing.T) {
	repo := &tokenRefreshAccountRepo{placeholder
	invalidator := &tokenCacheInvalidatorStub{placeholder
	cfg := &config.Config{
		TokenRefresh: config.TokenRefreshConfig{
			MaxRetries:          1,
			RetryBackoffSeconds: 0,
	placeholder,
placeholder
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg, nil)
	account := &Account{
		ID:       8,
		Platform: PlatformAntigravity,
		Type:     AccountTypeOAuth,
placeholder
	refresher := &tokenRefresherStub{
		credentials: map[string]any{
			"access_token": "ag-token",
	placeholder,
placeholder

	err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
placeholder
	require.Equal(t, 1, repo.updateCalls)
	require.Equal(t, 1, invalidator.calls) // Antigravity 也应触发缓存失效
placeholder

func TestAntigravityTokenRefresher_NeedsRefresh_ForceRefreshMarker(t *testing.T) {
	refresher := NewAntigravityTokenRefresher(nil)
	account := &Account{
		ID:       3675,
		Platform: PlatformAntigravity,
		Type:     AccountTypeOAuth,
placeholder
			"expires_at": time.Now().Add(time.Hour).Format(time.RFC3339),
	placeholder,
		Extra: map[string]any{
			antigravityForceTokenRefreshExtraKey: true,
	placeholder,
placeholder

	require.True(t, refresher.NeedsRefresh(account, 0), "server-invalidated token must refresh even before expires_at")
placeholder

func TestAntigravityTokenRefresher_NeedsRefresh_NormalExpiryRulesUnchanged(t *testing.T) {
	refresher := NewAntigravityTokenRefresher(nil)

	t.Run("normal_unexpired_without_marker_does_not_refresh", func(t *testing.T) {
		account := &Account{
			ID:       3707,
			Platform: PlatformAntigravity,
			Type:     AccountTypeOAuth,
	placeholder
				"expires_at": time.Now().Add(time.Hour).Format(time.RFC3339),
		placeholder,
	placeholder

		require.False(t, refresher.NeedsRefresh(account, 0))
placeholder)

	t.Run("normal_expiring_refreshes", func(t *testing.T) {
		account := &Account{
			ID:       3708,
			Platform: PlatformAntigravity,
			Type:     AccountTypeOAuth,
	placeholder
				"expires_at": time.Now().Add(5 * time.Minute).Format(time.RFC3339),
		placeholder,
	placeholder

		require.True(t, refresher.NeedsRefresh(account, 0))
placeholder)
placeholder

func TestTokenRefreshService_RefreshWithRetry_AntigravityClearsForceRefreshOnSuccess(t *testing.T) {
	repo := &tokenRefreshAccountRepo{placeholder
	cfg := &config.Config{
		TokenRefresh: config.TokenRefreshConfig{
			MaxRetries:          1,
			RetryBackoffSeconds: 0,
	placeholder,
placeholder
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, nil, nil, cfg, nil)
	until := time.Now().Add(10 * time.Minute)
	account := &Account{
		ID:                     3709,
		Platform:               PlatformAntigravity,
		Type:                   AccountTypeOAuth,
		TempUnschedulableUntil: &until,
		Extra: map[string]any{
			antigravityForceTokenRefreshExtraKey:       true,
			antigravityForceTokenRefreshReasonExtraKey: "401_invalid",
			"privacy_mode": AntigravityPrivacySet,
	placeholder,
placeholder
	refresher := &tokenRefresherStub{
		credentials: map[string]any{
			"access_token": "new-ag-token",
	placeholder,
placeholder

	err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
placeholder
	require.Equal(t, 1, repo.updateCredentialsCalls)
	require.Equal(t, 1, repo.updateExtraCalls)
	require.Equal(t, false, repo.lastExtraUpdates[antigravityForceTokenRefreshExtraKey])
	require.Equal(t, "", repo.lastExtraUpdates[antigravityForceTokenRefreshReasonExtraKey])
	require.Equal(t, false, account.Extra[antigravityForceTokenRefreshExtraKey])
	require.Equal(t, 1, repo.clearTempCalls, "successful refresh should restore schedulability")
placeholder

func TestTokenRefreshService_RefreshWithRetry_AntigravityForceRefreshInvalidGrantSetsError(t *testing.T) {
	repo := &tokenRefreshAccountRepo{placeholder
	cfg := &config.Config{
		TokenRefresh: config.TokenRefreshConfig{
			MaxRetries:          3,
			RetryBackoffSeconds: 0,
	placeholder,
placeholder
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, nil, nil, cfg, nil)
	account := &Account{
		ID:       3710,
		Platform: PlatformAntigravity,
		Type:     AccountTypeOAuth,
		Extra: map[string]any{
			antigravityForceTokenRefreshExtraKey:       true,
			antigravityForceTokenRefreshReasonExtraKey: "401_invalid",
	placeholder,
placeholder
	refresher := &tokenRefresherStub{
		err: errors.New("invalid_grant: token revoked"),
placeholder

	err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
placeholder
	require.Equal(t, 1, repo.setErrorCalls)
	require.Equal(t, 0, repo.setTempUnschedCalls)
	require.Equal(t, 1, repo.updateExtraCalls)
	require.Equal(t, false, repo.lastExtraUpdates[antigravityForceTokenRefreshExtraKey])
	require.Contains(t, repo.lastErrorMessage, "non-retryable")
placeholder

// TestTokenRefreshService_RefreshWithRetry_NonOAuthAccount 测试非 OAuth 账号不触发缓存失效
func TestTokenRefreshService_RefreshWithRetry_NonOAuthAccount(t *testing.T) {
	repo := &tokenRefreshAccountRepo{placeholder
	invalidator := &tokenCacheInvalidatorStub{placeholder
	cfg := &config.Config{
		TokenRefresh: config.TokenRefreshConfig{
			MaxRetries:          1,
			RetryBackoffSeconds: 0,
	placeholder,
placeholder
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg, nil)
	account := &Account{
		ID:       9,
placeholder
		Type:     AccountTypeAPIKey, // 非 OAuth
placeholder
	refresher := &tokenRefresherStub{
		credentials: map[string]any{
			"access_token": "token",
	placeholder,
placeholder

	err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
placeholder
	require.Equal(t, 1, repo.updateCalls)
	require.Equal(t, 0, invalidator.calls) // 非 OAuth 不触发缓存失效
placeholder

// TestTokenRefreshService_RefreshWithRetry_OtherPlatformOAuth 测试所有 OAuth 平台都触发缓存失效
func TestTokenRefreshService_RefreshWithRetry_OtherPlatformOAuth(t *testing.T) {
	repo := &tokenRefreshAccountRepo{placeholder
	invalidator := &tokenCacheInvalidatorStub{placeholder
	cfg := &config.Config{
		TokenRefresh: config.TokenRefreshConfig{
			MaxRetries:          1,
			RetryBackoffSeconds: 0,
	placeholder,
placeholder
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg, nil)
	account := &Account{
		ID:       10,
		Platform: PlatformOpenAI, // OpenAI OAuth 账户
		Type:     AccountTypeOAuth,
placeholder
	refresher := &tokenRefresherStub{
		credentials: map[string]any{
			"access_token": "token",
	placeholder,
placeholder

	err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
placeholder
	require.Equal(t, 1, repo.updateCalls)
	require.Equal(t, 1, repo.updateCredentialsCalls)
	require.Equal(t, 1, invalidator.calls) // 所有 OAuth 账户刷新后触发缓存失效
placeholder

func TestTokenRefreshService_RefreshWithRetry_UsesCredentialsUpdater(t *testing.T) {
	repo := &tokenRefreshAccountRepo{placeholder
	cfg := &config.Config{
		TokenRefresh: config.TokenRefreshConfig{
			MaxRetries:          1,
			RetryBackoffSeconds: 0,
	placeholder,
placeholder
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, nil, nil, cfg, nil)
	resetAt := time.Now().Add(30 * time.Minute)
	account := &Account{
		ID:               17,
		Platform:         PlatformOpenAI,
		Type:             AccountTypeOAuth,
		RateLimitResetAt: &resetAt,
placeholder
			"access_token": "old-token",
	placeholder,
placeholder
	refresher := &tokenRefresherStub{
		credentials: map[string]any{
			"access_token": "new-token",
	placeholder,
placeholder

	err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
placeholder
	require.Equal(t, 1, repo.updateCredentialsCalls)
	require.Equal(t, 0, repo.fullUpdateCalls)
	require.NotNil(t, account.RateLimitResetAt)
	require.WithinDuration(t, resetAt, *account.RateLimitResetAt, time.Second)
placeholder

// TestTokenRefreshService_RefreshWithRetry_UpdateFailed 测试更新失败的情况
func TestTokenRefreshService_RefreshWithRetry_UpdateFailed(t *testing.T) {
	repo := &tokenRefreshAccountRepo{updateErr: errors.New("update failed")placeholder
	invalidator := &tokenCacheInvalidatorStub{placeholder
	cfg := &config.Config{
		TokenRefresh: config.TokenRefreshConfig{
			MaxRetries:          1,
			RetryBackoffSeconds: 0,
	placeholder,
placeholder
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg, nil)
	account := &Account{
		ID:       11,
placeholder
		Type:     AccountTypeOAuth,
placeholder
	refresher := &tokenRefresherStub{
		credentials: map[string]any{
			"access_token": "token",
	placeholder,
placeholder

	err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
placeholder
	require.Contains(t, err.Error(), "failed to save credentials")
	require.Equal(t, 1, repo.updateCalls)
	require.Equal(t, 0, invalidator.calls) // 更新失败时不应触发缓存失效
placeholder

// TestTokenRefreshService_RefreshWithRetry_RefreshFailed 测试可重试错误耗尽不标记 error
func TestTokenRefreshService_RefreshWithRetry_RefreshFailed(t *testing.T) {
	repo := &tokenRefreshAccountRepo{placeholder
	invalidator := &tokenCacheInvalidatorStub{placeholder
	cfg := &config.Config{
		TokenRefresh: config.TokenRefreshConfig{
			MaxRetries:          2,
			RetryBackoffSeconds: 0,
	placeholder,
placeholder
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg, nil)
	account := &Account{
		ID:       12,
placeholder
		Type:     AccountTypeOAuth,
placeholder
	refresher := &tokenRefresherStub{
		err: errors.New("refresh failed"),
placeholder

	err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
placeholder
	require.Equal(t, 0, repo.updateCalls)   // 刷新失败不应更新
	require.Equal(t, 0, invalidator.calls)  // 刷新失败不应触发缓存失效
	require.Equal(t, 0, repo.setErrorCalls) // 可重试错误耗尽不标记 error，下个周期继续重试
placeholder

// TestTokenRefreshService_RefreshWithRetry_AntigravityRefreshFailed 测试 Antigravity 刷新失败不设置错误状态
func TestTokenRefreshService_RefreshWithRetry_AntigravityRefreshFailed(t *testing.T) {
	repo := &tokenRefreshAccountRepo{placeholder
	invalidator := &tokenCacheInvalidatorStub{placeholder
	cfg := &config.Config{
		TokenRefresh: config.TokenRefreshConfig{
			MaxRetries:          1,
			RetryBackoffSeconds: 0,
	placeholder,
placeholder
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg, nil)
	account := &Account{
		ID:       13,
		Platform: PlatformAntigravity,
		Type:     AccountTypeOAuth,
placeholder
	refresher := &tokenRefresherStub{
		err: errors.New("network error"), // 可重试错误
placeholder

	err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
placeholder
	require.Equal(t, 0, repo.updateCalls)
	require.Equal(t, 0, invalidator.calls)
	require.Equal(t, 0, repo.setErrorCalls) // Antigravity 可重试错误不设置错误状态
placeholder

// TestTokenRefreshService_RefreshWithRetry_AntigravityNonRetryableError 测试 Antigravity 不可重试错误
func TestTokenRefreshService_RefreshWithRetry_AntigravityNonRetryableError(t *testing.T) {
	repo := &tokenRefreshAccountRepo{placeholder
	invalidator := &tokenCacheInvalidatorStub{placeholder
	cfg := &config.Config{
		TokenRefresh: config.TokenRefreshConfig{
			MaxRetries:          3,
			RetryBackoffSeconds: 0,
	placeholder,
placeholder
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg, nil)
	account := &Account{
		ID:       14,
		Platform: PlatformAntigravity,
		Type:     AccountTypeOAuth,
placeholder
	refresher := &tokenRefresherStub{
		err: errors.New("invalid_grant: token revoked"), // 不可重试错误
placeholder

	err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
placeholder
	require.Equal(t, 0, repo.updateCalls)
	require.Equal(t, 0, invalidator.calls)
	require.Equal(t, 1, repo.setErrorCalls) // 不可重试错误应设置错误状态
placeholder

// TestTokenRefreshService_RefreshWithRetry_ClearsTempUnschedulable 测试刷新成功后清除临时不可调度（DB + Redis）
func TestTokenRefreshService_RefreshWithRetry_ClearsTempUnschedulable(t *testing.T) {
	repo := &tokenRefreshAccountRepo{placeholder
	invalidator := &tokenCacheInvalidatorStub{placeholder
	tempCache := &tempUnschedCacheStub{placeholder
	cfg := &config.Config{
		TokenRefresh: config.TokenRefreshConfig{
			MaxRetries:          1,
			RetryBackoffSeconds: 0,
	placeholder,
placeholder
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg, tempCache)
	until := time.Now().Add(10 * time.Minute)
	account := &Account{
		ID:                     15,
		Platform:               PlatformGemini,
		Type:                   AccountTypeOAuth,
		TempUnschedulableUntil: &until,
placeholder
	refresher := &tokenRefresherStub{
		credentials: map[string]any{
			"access_token": "new-token",
	placeholder,
placeholder

	err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
placeholder
	require.Equal(t, 1, repo.updateCalls)
	require.Equal(t, 1, repo.clearTempCalls)   // DB 清除
	require.Equal(t, 1, tempCache.deleteCalls) // Redis 缓存也应清除
placeholder

// TestTokenRefreshService_RefreshWithRetry_NonRetryableErrorAllPlatforms 测试所有平台不可重试错误都 SetError
func TestTokenRefreshService_RefreshWithRetry_NonRetryableErrorAllPlatforms(t *testing.T) {
	tests := []struct {
		name     string
		platform string
placeholder{
		{name: "gemini", platform: PlatformGeminiplaceholder,
		{name: "anthropic", platform: PlatformAnthropicplaceholder,
		{name: "openai", platform: PlatformOpenAIplaceholder,
		{name: "antigravity", platform: PlatformAntigravityplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &tokenRefreshAccountRepo{placeholder
			invalidator := &tokenCacheInvalidatorStub{placeholder
			cfg := &config.Config{
				TokenRefresh: config.TokenRefreshConfig{
					MaxRetries:          3,
					RetryBackoffSeconds: 0,
			placeholder,
		placeholder
			service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg, nil)
			account := &Account{
				ID:       16,
				Platform: tt.platform,
				Type:     AccountTypeOAuth,
		placeholder
			refresher := &tokenRefresherStub{
				err: errors.New("invalid_grant: token revoked"),
		placeholder

			err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
		placeholder
			require.Equal(t, 1, repo.setErrorCalls) // 所有平台不可重试错误都应 SetError
	placeholder)
placeholder
placeholder

func TestTokenRefreshService_RefreshWithRetry_NoRefreshTokenDoesNotTempUnschedule(t *testing.T) {
	repo := &tokenRefreshAccountRepo{placeholder
	cfg := &config.Config{
		TokenRefresh: config.TokenRefreshConfig{
			MaxRetries:          2,
			RetryBackoffSeconds: 0,
	placeholder,
placeholder
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, nil, nil, cfg, nil)
	account := &Account{
		ID:       18,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
	refresher := &tokenRefresherStub{
		err: errors.New("no refresh token available"),
placeholder

	err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
placeholder
	require.Equal(t, 0, repo.updateCalls)
	require.Equal(t, 0, repo.setTempUnschedCalls, "missing refresh token should not mark the account temp unschedulable")
	require.Equal(t, 1, repo.setErrorCalls, "missing refresh token should be treated as a non-retryable credential state")
placeholder

// TestIsNonRetryableRefreshError 测试不可重试错误判断
func TestIsNonRetryableRefreshError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
placeholder{
		{name: "nil_error", err: nil, expected: falseplaceholder,
		{name: "network_error", err: errors.New("network timeout"), expected: falseplaceholder,
		{name: "invalid_grant", err: errors.New("invalid_grant"), expected: trueplaceholder,
		{name: "invalid_client", err: errors.New("invalid_client"), expected: trueplaceholder,
		{name: "invalid_refresh_token", err: errors.New(`OPENAI_OAUTH_TOKEN_REFRESH_FAILED: token refresh failed: status 401, body: {"error":{"code":"invalid_refresh_token"placeholderplaceholder`), expected: trueplaceholder,
		{name: "token_expired", err: errors.New(`OPENAI_OAUTH_TOKEN_REFRESH_FAILED: token refresh failed: status 401, body: {"error":{"code":"token_expired"placeholderplaceholder`), expected: trueplaceholder,
		{name: "refresh_token_reused", err: errors.New(`OPENAI_OAUTH_TOKEN_REFRESH_FAILED: token refresh failed: status 401, body: {"error":{"code":"refresh_token_reused"placeholderplaceholder`), expected: trueplaceholder,
		{name: "app_session_terminated", err: errors.New(`OPENAI_OAUTH_TOKEN_REFRESH_FAILED: token refresh failed: status 401, body: {"error": {"code": "app_session_terminated"placeholderplaceholder`), expected: trueplaceholder,
		{name: "unauthorized_client", err: errors.New("unauthorized_client"), expected: trueplaceholder,
		{name: "access_denied", err: errors.New("access_denied"), expected: trueplaceholder,
		{name: "no_refresh_token", err: errors.New("no refresh token available"), expected: trueplaceholder,
		{name: "grok_entitlement_denied", err: errors.New("GROK_OAUTH_ENTITLEMENT_DENIED: subscription required"), expected: trueplaceholder,
		{name: "invalid_scope", err: errors.New("invalid_scope: requested scope is not allowed"), expected: trueplaceholder,
		{name: "invalid_grant_with_desc", err: errors.New("Error: invalid_grant - token revoked"), expected: trueplaceholder,
		{name: "case_insensitive", err: errors.New("INVALID_GRANT"), expected: trueplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNonRetryableRefreshError(tt.err)
			require.Equal(t, tt.expected, result)
	placeholder)
placeholder
placeholder

// ========== Path A (refreshAPI) 测试用例 ==========

// mockTokenCacheForRefreshAPI 用于 Path A 测试的 GeminiTokenCache mock
type mockTokenCacheForRefreshAPI struct {
	lockResult   bool
	lockErr      error
	releaseCalls int
placeholder

func (m *mockTokenCacheForRefreshAPI) GetAccessToken(_ context.Context, _ string) (string, error) {
	return "", errors.New("not cached")
placeholder

func (m *mockTokenCacheForRefreshAPI) SetAccessToken(_ context.Context, _ string, _ string, _ time.Duration) error {
	return nil
placeholder

func (m *mockTokenCacheForRefreshAPI) DeleteAccessToken(_ context.Context, _ string) error {
	return nil
placeholder

func (m *mockTokenCacheForRefreshAPI) AcquireRefreshLock(_ context.Context, _ string, _ time.Duration) (bool, error) {
	return m.lockResult, m.lockErr
placeholder

func (m *mockTokenCacheForRefreshAPI) ReleaseRefreshLock(_ context.Context, _ string) error {
	m.releaseCalls++
	return nil
placeholder

// buildPathAService 构建注入了 refreshAPI 的 service（Path A 测试辅助）
func buildPathAService(repo *tokenRefreshAccountRepo, cache GeminiTokenCache, invalidator TokenCacheInvalidator) (*TokenRefreshService, *tokenRefresherStub) {
	cfg := &config.Config{
		TokenRefresh: config.TokenRefreshConfig{
			MaxRetries:          1,
			RetryBackoffSeconds: 0,
	placeholder,
placeholder
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg, nil)
	refreshAPI := NewOAuthRefreshAPI(repo, cache)
	service.SetRefreshAPI(refreshAPI)

	refresher := &tokenRefresherStub{
		credentials: map[string]any{
			"access_token": "refreshed-token",
	placeholder,
placeholder
	return service, refresher
placeholder

// TestPathA_Success 统一 API 路径正常成功：刷新 + DB 更新 + postRefreshActions
func TestPathA_Success(t *testing.T) {
	account := &Account{
		ID:       100,
placeholder
		Type:     AccountTypeOAuth,
		Status:   StatusActive,
placeholder
	repo := &tokenRefreshAccountRepo{placeholder
	repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
	invalidator := &tokenCacheInvalidatorStub{placeholder
	cache := &mockTokenCacheForRefreshAPI{lockResult: trueplaceholder

	service, refresher := buildPathAService(repo, cache, invalidator)

	err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
placeholder
	require.Equal(t, 1, repo.updateCalls)   // DB 更新被调用
	require.Equal(t, 1, invalidator.calls)  // 缓存失效被调用
	require.Equal(t, 1, cache.releaseCalls) // 锁被释放
placeholder

// TestPathA_LockHeld 锁被其他 worker 持有 → 返回 errRefreshSkipped
func TestPathA_LockHeld(t *testing.T) {
	account := &Account{
		ID:       101,
placeholder
		Type:     AccountTypeOAuth,
		Status:   StatusActive,
placeholder
	repo := &tokenRefreshAccountRepo{placeholder
	invalidator := &tokenCacheInvalidatorStub{placeholder
	cache := &mockTokenCacheForRefreshAPI{lockResult: falseplaceholder // 锁获取失败（被占）

	service, refresher := buildPathAService(repo, cache, invalidator)

	err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
	require.ErrorIs(t, err, errRefreshSkipped)
	require.Equal(t, 0, repo.updateCalls)  // 不应更新 DB
	require.Equal(t, 0, invalidator.calls) // 不应触发缓存失效
placeholder

// TestPathA_AlreadyRefreshed 二次检查发现已被其他路径刷新 → 返回 errRefreshSkipped
func TestPathA_AlreadyRefreshed(t *testing.T) {
	// NeedsRefresh 返回 false → RefreshIfNeeded 返回 {Refreshed: falseplaceholder
	account := &Account{
		ID:       102,
placeholder
		Type:     AccountTypeOAuth,
		Status:   StatusActive,
placeholder
	repo := &tokenRefreshAccountRepo{placeholder
	repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
	invalidator := &tokenCacheInvalidatorStub{placeholder
	cache := &mockTokenCacheForRefreshAPI{lockResult: trueplaceholder

	service, _ := buildPathAService(repo, cache, invalidator)

	// 使用一个 NeedsRefresh 返回 false 的 stub
	noRefreshNeeded := &tokenRefresherStub{
		credentials: map[string]any{"access_token": "token"placeholder,
placeholder
	// 覆盖 NeedsRefresh 行为 — 我们需要一个新的 stub 类型
	alwaysFreshStub := &alwaysFreshRefresherStub{placeholder

	err := service.refreshWithRetry(context.Background(), account, noRefreshNeeded, alwaysFreshStub, time.Hour)
	require.ErrorIs(t, err, errRefreshSkipped)
	require.Equal(t, 0, repo.updateCalls)
	require.Equal(t, 0, invalidator.calls)
placeholder

// alwaysFreshRefresherStub 二次检查时认为不需要刷新（模拟已被其他路径刷新）
type alwaysFreshRefresherStub struct{placeholder

func (r *alwaysFreshRefresherStub) CanRefresh(_ *Account) bool                    { return true placeholder
func (r *alwaysFreshRefresherStub) NeedsRefresh(_ *Account, _ time.Duration) bool { return false placeholder
func (r *alwaysFreshRefresherStub) Refresh(_ context.Context, _ *Account) (map[string]any, error) {
	return nil, errors.New("should not be called")
placeholder
func (r *alwaysFreshRefresherStub) CacheKey(account *Account) string {
	return "test:fresh:" + account.Platform
placeholder

// TestPathA_NonRetryableError 统一 API 路径返回不可重试错误 → SetError
func TestPathA_NonRetryableError(t *testing.T) {
	account := &Account{
		ID:       103,
placeholder
		Type:     AccountTypeOAuth,
		Status:   StatusActive,
placeholder
	repo := &tokenRefreshAccountRepo{placeholder
	repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
	invalidator := &tokenCacheInvalidatorStub{placeholder
	cache := &mockTokenCacheForRefreshAPI{lockResult: trueplaceholder

	service, _ := buildPathAService(repo, cache, invalidator)

	refresher := &tokenRefresherStub{
		err: errors.New("invalid_grant: token revoked"),
placeholder

	err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
placeholder
	require.Equal(t, 1, repo.setErrorCalls) // 应标记 error 状态
	require.Equal(t, 0, repo.updateCalls)   // 不应更新 credentials
	require.Equal(t, 0, invalidator.calls)  // 不应触发缓存失效
placeholder

// TestPathA_RetryableErrorExhausted 统一 API 路径可重试错误耗尽 → 不标记 error
func TestPathA_RetryableErrorExhausted(t *testing.T) {
	account := &Account{
		ID:       104,
placeholder
		Type:     AccountTypeOAuth,
		Status:   StatusActive,
placeholder
	repo := &tokenRefreshAccountRepo{placeholder
	repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
	invalidator := &tokenCacheInvalidatorStub{placeholder
	cache := &mockTokenCacheForRefreshAPI{lockResult: trueplaceholder

	cfg := &config.Config{
		TokenRefresh: config.TokenRefreshConfig{
			MaxRetries:          2,
			RetryBackoffSeconds: 0,
	placeholder,
placeholder
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg, nil)
	refreshAPI := NewOAuthRefreshAPI(repo, cache)
	service.SetRefreshAPI(refreshAPI)

	refresher := &tokenRefresherStub{
		err: errors.New("network timeout"),
placeholder

	err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
placeholder
	require.Equal(t, 0, repo.setErrorCalls) // 可重试错误不标记 error
	require.Equal(t, 0, repo.updateCalls)   // 刷新失败不应更新
	require.Equal(t, 0, invalidator.calls)  // 不应触发缓存失效
placeholder

// TestPathA_DBUpdateFailed 统一 API 路径 DB 更新失败 → 返回 error，不执行 postRefreshActions
func TestPathA_DBUpdateFailed(t *testing.T) {
	account := &Account{
		ID:       105,
placeholder
		Type:     AccountTypeOAuth,
		Status:   StatusActive,
placeholder
	repo := &tokenRefreshAccountRepo{updateErr: errors.New("db connection lost")placeholder
	repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
	invalidator := &tokenCacheInvalidatorStub{placeholder
	cache := &mockTokenCacheForRefreshAPI{lockResult: trueplaceholder

	service, refresher := buildPathAService(repo, cache, invalidator)

	err := service.refreshWithRetry(context.Background(), account, refresher, refresher, time.Hour)
placeholder
	require.ErrorIs(t, err, errOAuthRefreshCredentialPersist)
	require.Equal(t, 1, repo.updateCalls)  // DB 更新被尝试
	require.Equal(t, 0, invalidator.calls) // DB 失败时不应触发缓存失效
placeholder
