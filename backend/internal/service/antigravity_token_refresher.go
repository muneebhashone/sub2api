package service

import (
	"context"
	"fmt"
	"log"
	"time"
)

const (
	// antigravityRefreshWindow Antigravity token 提前刷新窗口：15分钟
	// Google OAuth token 有效期55分钟，提前15分钟刷新
	antigravityRefreshWindow = 15 * time.Minute
)

// AntigravityTokenRefresher 实现 TokenRefresher 接口
type AntigravityTokenRefresher struct {
	antigravityOAuthService *AntigravityOAuthService
placeholder

func NewAntigravityTokenRefresher(antigravityOAuthService *AntigravityOAuthService) *AntigravityTokenRefresher {
	return &AntigravityTokenRefresher{
		antigravityOAuthService: antigravityOAuthService,
placeholder
placeholder

// CanRefresh 检查是否可以刷新此账户
func (r *AntigravityTokenRefresher) CanRefresh(account *Account) bool {
	return account.Platform == PlatformAntigravity && account.Type == AccountTypeOAuth
placeholder

// NeedsRefresh 检查账户是否需要刷新
// Antigravity 使用固定的15分钟刷新窗口，忽略全局配置
func (r *AntigravityTokenRefresher) NeedsRefresh(account *Account, _ time.Duration) bool {
	if !r.CanRefresh(account) {
		return false
placeholder
	expiresAt := account.GetCredentialAsTime("expires_at")
	if expiresAt == nil {
		return false
placeholder
	timeUntilExpiry := time.Until(*expiresAt)
	needsRefresh := timeUntilExpiry < antigravityRefreshWindow
	if needsRefresh {
		fmt.Printf("[AntigravityTokenRefresher] Account %d needs refresh: expires_at=%s, time_until_expiry=%v, window=%v\n",
			account.ID, expiresAt.Format("2006-01-02 15:04:05"), timeUntilExpiry, antigravityRefreshWindow)
placeholder
	return needsRefresh
placeholder

// Refresh 执行 token 刷新
func (r *AntigravityTokenRefresher) Refresh(ctx context.Context, account *Account) (map[string]any, error) {
	tokenInfo, err := r.antigravityOAuthService.RefreshAccountToken(ctx, account)
	if err != nil {
		return nil, err
placeholder

	newCredentials := r.antigravityOAuthService.BuildAccountCredentials(tokenInfo)
	for k, v := range account.Credentials {
		if _, exists := newCredentials[k]; !exists {
			newCredentials[k] = v
	placeholder
placeholder

	// 如果 project_id 获取失败但之前有 project_id，不返回错误（只是临时网络故障）
	// 只有真正缺失 project_id（从未获取过）时才返回错误
	if tokenInfo.ProjectIDMissing {
		// 检查是否保留了旧的 project_id
		if tokenInfo.ProjectID != "" {
			// 有旧的 project_id，只是本次获取失败，记录警告但不返回错误
			log.Printf("[AntigravityTokenRefresher] Account %d: LoadCodeAssist 临时失败，保留旧 project_id", account.ID)
	placeholder else {
			// 真正缺失 project_id，返回错误
			return newCredentials, fmt.Errorf("missing_project_id: 账户缺少project id，可能无法使用Antigravity")
	placeholder
placeholder

	return newCredentials, nil
placeholder
