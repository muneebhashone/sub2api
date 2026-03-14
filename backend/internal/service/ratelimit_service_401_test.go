//go:build unit

package service

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type rateLimitAccountRepoStub struct {
	mockAccountRepoForGemini
	setErrorCalls int
	tempCalls     int
	lastErrorMsg  string
placeholder

func (r *rateLimitAccountRepoStub) SetError(ctx context.Context, id int64, errorMsg string) error {
	r.setErrorCalls++
	r.lastErrorMsg = errorMsg
	return nil
placeholder

func (r *rateLimitAccountRepoStub) SetTempUnschedulable(ctx context.Context, id int64, until time.Time, reason string) error {
	r.tempCalls++
	return nil
placeholder

type tokenCacheInvalidatorRecorder struct {
	accounts []*Account
	err      error
placeholder

func (r *tokenCacheInvalidatorRecorder) InvalidateToken(ctx context.Context, account *Account) error {
	r.accounts = append(r.accounts, account)
	return r.err
placeholder

func TestRateLimitService_HandleUpstreamError_OAuth401SetsTempUnschedulable(t *testing.T) {
	t.Run("gemini", func(t *testing.T) {
		repo := &rateLimitAccountRepoStub{placeholder
		invalidator := &tokenCacheInvalidatorRecorder{placeholder
		service := NewRateLimitService(repo, nil, &config.Config{placeholder, nil, nil)
		service.SetTokenCacheInvalidator(invalidator)
		account := &Account{
			ID:       100,
	placeholder
			Type:     AccountTypeOAuth,
	placeholder
				"temp_unschedulable_enabled": true,
				"temp_unschedulable_rules": []any{
					map[string]any{
						"error_code":       401,
						"keywords":         []any{"unauthorized"placeholder,
						"duration_minutes": 30,
						"description":      "custom rule",
				placeholder,
			placeholder,
		placeholder,
	placeholder

		shouldDisable := service.HandleUpstreamError(context.Background(), account, 401, http.Header{placeholder, []byte("unauthorized"))

		require.True(t, shouldDisable)
		require.Equal(t, 0, repo.setErrorCalls)
		require.Equal(t, 1, repo.tempCalls)
		require.Len(t, invalidator.accounts, 1)
placeholder)

	t.Run("antigravity_401_uses_SetError", func(t *testing.T) {
		// Antigravity 401 由 applyErrorPolicy 的 temp_unschedulable_rules 控制，
		// HandleUpstreamError 中走 SetError 路径。
		repo := &rateLimitAccountRepoStub{placeholder
		invalidator := &tokenCacheInvalidatorRecorder{placeholder
		service := NewRateLimitService(repo, nil, &config.Config{placeholder, nil, nil)
		service.SetTokenCacheInvalidator(invalidator)
		account := &Account{
			ID:       100,
			Platform: PlatformAntigravity,
			Type:     AccountTypeOAuth,
	placeholder

		shouldDisable := service.HandleUpstreamError(context.Background(), account, 401, http.Header{placeholder, []byte("unauthorized"))

		require.True(t, shouldDisable)
		require.Equal(t, 1, repo.setErrorCalls)
		require.Equal(t, 0, repo.tempCalls)
		require.Empty(t, invalidator.accounts)
placeholder)
placeholder

func TestRateLimitService_HandleUpstreamError_OAuth401InvalidatorError(t *testing.T) {
	repo := &rateLimitAccountRepoStub{placeholder
	invalidator := &tokenCacheInvalidatorRecorder{err: errors.New("boom")placeholder
	service := NewRateLimitService(repo, nil, &config.Config{placeholder, nil, nil)
	service.SetTokenCacheInvalidator(invalidator)
	account := &Account{
		ID:       101,
placeholder
		Type:     AccountTypeOAuth,
placeholder

	shouldDisable := service.HandleUpstreamError(context.Background(), account, 401, http.Header{placeholder, []byte("unauthorized"))

	require.True(t, shouldDisable)
	require.Equal(t, 0, repo.setErrorCalls)
	require.Equal(t, 1, repo.tempCalls)
	require.Len(t, invalidator.accounts, 1)
placeholder

func TestRateLimitService_HandleUpstreamError_NonOAuth401(t *testing.T) {
	repo := &rateLimitAccountRepoStub{placeholder
	invalidator := &tokenCacheInvalidatorRecorder{placeholder
	service := NewRateLimitService(repo, nil, &config.Config{placeholder, nil, nil)
	service.SetTokenCacheInvalidator(invalidator)
	account := &Account{
		ID:       102,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
placeholder

	shouldDisable := service.HandleUpstreamError(context.Background(), account, 401, http.Header{placeholder, []byte("unauthorized"))

	require.True(t, shouldDisable)
	require.Equal(t, 1, repo.setErrorCalls)
	require.Empty(t, invalidator.accounts)
placeholder
