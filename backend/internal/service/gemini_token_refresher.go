package service

import (
	"context"
	"strconv"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/model"
)

type GeminiTokenRefresher struct {
	geminiOAuthService *GeminiOAuthService
placeholder

func NewGeminiTokenRefresher(geminiOAuthService *GeminiOAuthService) *GeminiTokenRefresher {
	return &GeminiTokenRefresher{geminiOAuthService: geminiOAuthServiceplaceholder
placeholder

func (r *GeminiTokenRefresher) CanRefresh(account *model.Account) bool {
	return account.Platform == model.PlatformGemini && account.Type == model.AccountTypeOAuth
placeholder

func (r *GeminiTokenRefresher) NeedsRefresh(account *model.Account, refreshWindow time.Duration) bool {
	if !r.CanRefresh(account) {
		return false
placeholder
	expiresAtStr := account.GetCredential("expires_at")
	if expiresAtStr == "" {
		return false
placeholder
	expiresAt, err := strconv.ParseInt(expiresAtStr, 10, 64)
	if err != nil {
		return false
placeholder
	expiryTime := time.Unix(expiresAt, 0)
	return time.Until(expiryTime) < refreshWindow
placeholder

func (r *GeminiTokenRefresher) Refresh(ctx context.Context, account *model.Account) (map[string]any, error) {
	tokenInfo, err := r.geminiOAuthService.RefreshAccountToken(ctx, account)
	if err != nil {
		return nil, err
placeholder

	newCredentials := r.geminiOAuthService.BuildAccountCredentials(tokenInfo)
	for k, v := range account.Credentials {
		if _, exists := newCredentials[k]; !exists {
			newCredentials[k] = v
	placeholder
placeholder

	return newCredentials, nil
placeholder
