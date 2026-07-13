package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/imroc/req/v3"
	"github.com/stretchr/testify/require"
)

type tokenRefreshCandidateRepo struct {
	AccountRepository
	accounts              []Account
	updatedCredentialIDs  []int64
	setErrorCalls         int
	setTempUnschedCalls   int
	clearTempCalls        int
	lastTempUnschedReason string
	listActiveCalls       int
placeholder

func (r *tokenRefreshCandidateRepo) ListActive(context.Context) ([]Account, error) {
	r.listActiveCalls++
	return r.accounts, nil
placeholder

func (r *tokenRefreshCandidateRepo) ListOAuthRefreshCandidatePage(_ context.Context, options OAuthRefreshPageOptions) (*OAuthRefreshCandidatePage, error) {
	candidates := make([]Account, 0, len(r.accounts))
	now := time.Now()
	for _, account := range r.accounts {
		if account.ID <= options.AfterID {
			continue
	placeholder
		refreshToken, _ := account.Credentials["refresh_token"].(string)
		inRetryCooldown := account.TempUnschedulableUntil != nil &&
			account.TempUnschedulableUntil.After(now) &&
			strings.HasPrefix(account.TempUnschedulableReason, "token refresh retry exhausted:")
		platformAllowed := false
		for _, platform := range options.Platforms {
			if account.Platform == platform {
				platformAllowed = true
				break
		placeholder
	placeholder
		if options.ActiveOnly && account.Status != StatusActive ||
			account.Type != AccountTypeOAuth ||
			!platformAllowed ||
			options.RequireRefreshToken && strings.TrimSpace(refreshToken) == "" ||
			options.ExcludeRetryCooldown && inRetryCooldown {
			continue
	placeholder
		candidates = append(candidates, account)
		if len(candidates) == options.Limit {
			break
	placeholder
placeholder
	page := &OAuthRefreshCandidatePage{Accounts: candidates, HasMore: len(candidates) == options.Limitplaceholder
	if len(candidates) > 0 {
		page.NextAfterID = candidates[len(candidates)-1].ID
placeholder
	return page, nil
placeholder

func (r *tokenRefreshCandidateRepo) UpdateCredentials(_ context.Context, id int64, _ map[string]any) error {
	r.updatedCredentialIDs = append(r.updatedCredentialIDs, id)
	return nil
placeholder

func (r *tokenRefreshCandidateRepo) SetError(context.Context, int64, string) error {
	r.setErrorCalls++
	return nil
placeholder

func (r *tokenRefreshCandidateRepo) SetTempUnschedulable(_ context.Context, _ int64, _ time.Time, reason string) error {
	r.setTempUnschedCalls++
	r.lastTempUnschedReason = reason
	return nil
placeholder

func (r *tokenRefreshCandidateRepo) ClearTempUnschedulable(context.Context, int64) error {
	r.clearTempCalls++
	return nil
placeholder

func isOAuthRefreshPlatform(platform string) bool {
	switch platform {
	case PlatformAnthropic, PlatformOpenAI, PlatformGemini, PlatformAntigravity, PlatformGrok:
		return true
	default:
		return false
placeholder
placeholder

type tokenRefreshTestRefresher struct {
	err error
placeholder

func (r *tokenRefreshTestRefresher) CanRefresh(*Account) bool { return true placeholder

func (r *tokenRefreshTestRefresher) NeedsRefresh(*Account, time.Duration) bool { return true placeholder

func (r *tokenRefreshTestRefresher) Refresh(context.Context, *Account) (map[string]any, error) {
	if r.err != nil {
		return nil, r.err
placeholder
	return map[string]any{"access_token": "new-access-token", "refresh_token": "new-refresh-token"placeholder, nil
placeholder

func TestTokenRefreshService_ProcessRefreshUsesOAuthRefreshCandidates(t *testing.T) {
	future := time.Now().Add(10 * time.Minute)
	repo := &tokenRefreshCandidateRepo{
		accounts: []Account{
			{
				ID:          1,
				Platform:    PlatformOpenAI,
				Type:        AccountTypeOAuth,
				Status:      StatusActive,
		placeholder"refresh_token": "refresh-token"placeholder,
		placeholder,
			{
				ID:          2,
				Platform:    PlatformOpenAI,
				Type:        AccountTypeOAuth,
				Status:      StatusActive,
		placeholderplaceholder,
		placeholder,
			{
				ID:          3,
				Platform:    PlatformGemini,
				Type:        AccountTypeAPIKey,
				Status:      StatusActive,
		placeholder"refresh_token": "refresh-token"placeholder,
		placeholder,
			{
				ID:                      4,
				Platform:                PlatformAntigravity,
				Type:                    AccountTypeOAuth,
				Status:                  StatusActive,
				Credentials:             map[string]any{"refresh_token": "refresh-token"placeholder,
				TempUnschedulableUntil:  &future,
				TempUnschedulableReason: "token refresh retry exhausted: network timeout",
		placeholder,
			{
				ID:          5,
				Platform:    "other",
				Type:        AccountTypeOAuth,
				Status:      StatusActive,
		placeholder"refresh_token": "refresh-token"placeholder,
		placeholder,
			{
				ID:                      6,
				Platform:                PlatformAntigravity,
				Type:                    AccountTypeOAuth,
				Status:                  StatusActive,
				Credentials:             map[string]any{"refresh_token": "refresh-token"placeholder,
				Extra:                   map[string]any{"privacy_mode": AntigravityPrivacySetplaceholder,
				TempUnschedulableUntil:  &future,
				TempUnschedulableReason: "OAuth 401: unauthorized",
		placeholder,
	placeholder,
placeholder
	svc := &TokenRefreshService{
		accountRepo:    repo,
		candidatePager: repo,
		registrations: []tokenRefreshRegistration{
			{platform: PlatformOpenAI, refresher: &tokenRefreshTestRefresher{placeholderplaceholder,
			{platform: PlatformGemini, refresher: &tokenRefreshTestRefresher{placeholderplaceholder,
			{platform: PlatformAntigravity, refresher: &tokenRefreshTestRefresher{placeholderplaceholder,
	placeholder,
		refreshPolicy: DefaultBackgroundRefreshPolicy(),
		cfg:           &config.TokenRefreshConfig{RefreshBeforeExpiryHours: 1, MaxRetries: 1placeholder,
placeholder

	svc.processRefresh()

	require.Zero(t, repo.listActiveCalls, "TokenRefreshService should not use the broad active-account query")
	require.ElementsMatch(t, []int64{1, 6placeholder, repo.updatedCredentialIDs)
	require.Equal(t, 1, repo.clearTempCalls, "successful refresh should clear the OAuth 401 temp-unschedulable state")
placeholder

func TestTokenRefreshService_RefreshFailureDoesNotCallPrivacy(t *testing.T) {
	tests := []struct {
		name string
		err  error
placeholder{
		{name: "retry exhausted", err: errors.New("temporary upstream timeout")placeholder,
		{name: "non retryable", err: errors.New("invalid_grant: token revoked")placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &tokenRefreshCandidateRepo{placeholder
			svc := &TokenRefreshService{
				accountRepo:   repo,
				refreshPolicy: DefaultBackgroundRefreshPolicy(),
				cfg:           &config.TokenRefreshConfig{MaxRetries: 1, RetryBackoffSeconds: 0placeholder,
				privacyClientFactory: func(string) (*req.Client, error) {
					t.Fatalf("privacy client factory must not be called on refresh failure")
					return nil, errors.New("unexpected privacy call")
			placeholder,
		placeholder
			account := &Account{
				ID:       11,
				Platform: PlatformOpenAI,
				Type:     AccountTypeOAuth,
		placeholder
					"access_token":  "old-access-token",
					"refresh_token": "refresh-token",
			placeholder,
		placeholder

			err := svc.refreshWithRetry(context.Background(), account, &tokenRefreshTestRefresher{err: tt.errplaceholder, nil, time.Hour)

		placeholder
			if isNonRetryableRefreshError(tt.err) {
				require.Equal(t, 1, repo.setErrorCalls)
				require.Zero(t, repo.setTempUnschedCalls)
		placeholder else {
				require.Zero(t, repo.setErrorCalls)
				require.Equal(t, 1, repo.setTempUnschedCalls)
				require.True(t, strings.HasPrefix(repo.lastTempUnschedReason, "token refresh retry exhausted:"))
		placeholder
	placeholder)
placeholder
placeholder
