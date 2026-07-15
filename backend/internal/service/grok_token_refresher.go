package service

import (
	"context"
	"errors"
	"strings"
	"time"
)

const grokTokenRefreshSkew = time.Hour

type GrokTokenRefresher struct {
	grokOAuthService GrokOAuthTokenService
placeholder

func NewGrokTokenRefresher(grokOAuthService GrokOAuthTokenService) *GrokTokenRefresher {
	return &GrokTokenRefresher{grokOAuthService: grokOAuthServiceplaceholder
placeholder

func (r *GrokTokenRefresher) CacheKey(account *Account) string {
	return GrokTokenCacheKey(account)
placeholder

func (r *GrokTokenRefresher) CanRefresh(account *Account) bool {
	return account != nil && account.Platform == PlatformGrok && account.Type == AccountTypeOAuth &&
		strings.TrimSpace(account.GetGrokRefreshToken()) != ""
placeholder

func (r *GrokTokenRefresher) NeedsRefresh(account *Account, refreshWindow time.Duration) bool {
	if account == nil || strings.TrimSpace(account.GetGrokRefreshToken()) == "" {
		return false
placeholder
	expiresAt := account.GetCredentialAsTime("expires_at")
	if expiresAt == nil {
		return true
placeholder
	if refreshWindow < grokTokenRefreshSkew {
		refreshWindow = grokTokenRefreshSkew
placeholder
	return time.Until(*expiresAt) < refreshWindow
placeholder

func (r *GrokTokenRefresher) Refresh(ctx context.Context, account *Account) (map[string]any, error) {
	if r == nil || r.grokOAuthService == nil {
		return nil, errors.New("grok oauth service is not configured")
placeholder
	tokenInfo, err := r.grokOAuthService.RefreshAccountToken(ctx, account)
	if err != nil {
		return nil, err
placeholder
	newCredentials := r.grokOAuthService.BuildAccountCredentials(tokenInfo)
	newCredentials = MergeCredentials(account.Credentials, newCredentials)
	if baseURL := strings.TrimSpace(account.GetCredential("base_url")); baseURL != "" {
		newCredentials["base_url"] = baseURL
placeholder
	return newCredentials, nil
placeholder
