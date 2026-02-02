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
	updateCalls   int
	setErrorCalls int
	lastAccount   *Account
	updateErr     error
placeholder

func (r *tokenRefreshAccountRepo) Update(ctx context.Context, account *Account) error {
	r.updateCalls++
	r.lastAccount = account
	return r.updateErr
placeholder

func (r *tokenRefreshAccountRepo) SetError(ctx context.Context, id int64, errorMsg string) error {
	r.setErrorCalls++
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

func TestTokenRefreshService_RefreshWithRetry_InvalidatesCache(t *testing.T) {
	repo := &tokenRefreshAccountRepo{placeholder
	invalidator := &tokenCacheInvalidatorStub{placeholder
	cfg := &config.Config{
		TokenRefresh: config.TokenRefreshConfig{
			MaxRetries:          1,
			RetryBackoffSeconds: 0,
	placeholder,
placeholder
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg)
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

	err := service.refreshWithRetry(context.Background(), account, refresher)
placeholder
	require.Equal(t, 1, repo.updateCalls)
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
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg)
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

	err := service.refreshWithRetry(context.Background(), account, refresher)
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
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, nil, nil, cfg)
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

	err := service.refreshWithRetry(context.Background(), account, refresher)
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
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg)
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

	err := service.refreshWithRetry(context.Background(), account, refresher)
placeholder
	require.Equal(t, 1, repo.updateCalls)
	require.Equal(t, 1, invalidator.calls) // Antigravity 也应触发缓存失效
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
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg)
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

	err := service.refreshWithRetry(context.Background(), account, refresher)
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
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg)
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

	err := service.refreshWithRetry(context.Background(), account, refresher)
placeholder
	require.Equal(t, 1, repo.updateCalls)
	require.Equal(t, 1, invalidator.calls) // 所有 OAuth 账户刷新后触发缓存失效
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
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg)
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

	err := service.refreshWithRetry(context.Background(), account, refresher)
placeholder
	require.Contains(t, err.Error(), "failed to save credentials")
	require.Equal(t, 1, repo.updateCalls)
	require.Equal(t, 0, invalidator.calls) // 更新失败时不应触发缓存失效
placeholder

// TestTokenRefreshService_RefreshWithRetry_RefreshFailed 测试刷新失败的情况
func TestTokenRefreshService_RefreshWithRetry_RefreshFailed(t *testing.T) {
	repo := &tokenRefreshAccountRepo{placeholder
	invalidator := &tokenCacheInvalidatorStub{placeholder
	cfg := &config.Config{
		TokenRefresh: config.TokenRefreshConfig{
			MaxRetries:          2,
			RetryBackoffSeconds: 0,
	placeholder,
placeholder
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg)
	account := &Account{
		ID:       12,
placeholder
		Type:     AccountTypeOAuth,
placeholder
	refresher := &tokenRefresherStub{
		err: errors.New("refresh failed"),
placeholder

	err := service.refreshWithRetry(context.Background(), account, refresher)
placeholder
	require.Equal(t, 0, repo.updateCalls)   // 刷新失败不应更新
	require.Equal(t, 0, invalidator.calls)  // 刷新失败不应触发缓存失效
	require.Equal(t, 1, repo.setErrorCalls) // 应设置错误状态
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
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg)
	account := &Account{
		ID:       13,
		Platform: PlatformAntigravity,
		Type:     AccountTypeOAuth,
placeholder
	refresher := &tokenRefresherStub{
		err: errors.New("network error"), // 可重试错误
placeholder

	err := service.refreshWithRetry(context.Background(), account, refresher)
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
	service := NewTokenRefreshService(repo, nil, nil, nil, nil, invalidator, nil, cfg)
	account := &Account{
		ID:       14,
		Platform: PlatformAntigravity,
		Type:     AccountTypeOAuth,
placeholder
	refresher := &tokenRefresherStub{
		err: errors.New("invalid_grant: token revoked"), // 不可重试错误
placeholder

	err := service.refreshWithRetry(context.Background(), account, refresher)
placeholder
	require.Equal(t, 0, repo.updateCalls)
	require.Equal(t, 0, invalidator.calls)
	require.Equal(t, 1, repo.setErrorCalls) // 不可重试错误应设置错误状态
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
		{name: "unauthorized_client", err: errors.New("unauthorized_client"), expected: trueplaceholder,
		{name: "access_denied", err: errors.New("access_denied"), expected: trueplaceholder,
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
