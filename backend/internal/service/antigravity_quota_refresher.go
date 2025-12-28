package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/antigravity"
)

// AntigravityQuotaRefresher 定时刷新 Antigravity 账户的配额信息
type AntigravityQuotaRefresher struct {
	accountRepo AccountRepository
	proxyRepo   ProxyRepository
	oauthSvc    *AntigravityOAuthService
	cfg         *config.TokenRefreshConfig

	stopCh chan struct{placeholder
	wg     sync.WaitGroup
placeholder

// NewAntigravityQuotaRefresher 创建配额刷新器
func NewAntigravityQuotaRefresher(
	accountRepo AccountRepository,
	proxyRepo ProxyRepository,
	oauthSvc *AntigravityOAuthService,
	cfg *config.Config,
) *AntigravityQuotaRefresher {
	return &AntigravityQuotaRefresher{
		accountRepo: accountRepo,
		proxyRepo:   proxyRepo,
		oauthSvc:    oauthSvc,
		cfg:         &cfg.TokenRefresh,
		stopCh:      make(chan struct{placeholder),
placeholder
placeholder

// Start 启动后台配额刷新服务
func (r *AntigravityQuotaRefresher) Start() {
	if !r.cfg.Enabled {
		log.Println("[AntigravityQuota] Service disabled by configuration")
		return
placeholder

	r.wg.Add(1)
	go r.refreshLoop()

	log.Printf("[AntigravityQuota] Service started (check every %d minutes)", r.cfg.CheckIntervalMinutes)
placeholder

// Stop 停止服务
func (r *AntigravityQuotaRefresher) Stop() {
	close(r.stopCh)
	r.wg.Wait()
	log.Println("[AntigravityQuota] Service stopped")
placeholder

// refreshLoop 刷新循环
func (r *AntigravityQuotaRefresher) refreshLoop() {
	defer r.wg.Done()

	checkInterval := time.Duration(r.cfg.CheckIntervalMinutes) * time.Minute
	if checkInterval < time.Minute {
		checkInterval = 5 * time.Minute
placeholder

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	// 启动时立即执行一次
	r.processRefresh()

	for {
		select {
		case <-ticker.C:
			r.processRefresh()
		case <-r.stopCh:
			return
	placeholder
placeholder
placeholder

// processRefresh 执行一次刷新
func (r *AntigravityQuotaRefresher) processRefresh() {
	ctx := context.Background()

	// 查询所有 active 的账户，然后过滤 antigravity 平台
	allAccounts, err := r.accountRepo.ListActive(ctx)
	if err != nil {
		log.Printf("[AntigravityQuota] Failed to list accounts: %v", err)
		return
placeholder

	// 过滤 antigravity 平台账户
	var accounts []Account
	for _, acc := range allAccounts {
		if acc.Platform == PlatformAntigravity {
			accounts = append(accounts, acc)
	placeholder
placeholder

	if len(accounts) == 0 {
		return
placeholder

	refreshed, failed := 0, 0

	for i := range accounts {
		account := &accounts[i]

		if err := r.refreshAccountQuota(ctx, account); err != nil {
			log.Printf("[AntigravityQuota] Account %d (%s) failed: %v", account.ID, account.Name, err)
			failed++
	placeholder else {
			refreshed++
	placeholder
placeholder

	log.Printf("[AntigravityQuota] Cycle complete: total=%d, refreshed=%d, failed=%d",
		len(accounts), refreshed, failed)
placeholder

// refreshAccountQuota 刷新单个账户的配额
func (r *AntigravityQuotaRefresher) refreshAccountQuota(ctx context.Context, account *Account) error {
	accessToken := account.GetCredential("access_token")
	projectID := account.GetCredential("project_id")

	if accessToken == "" || projectID == "" {
		return nil // 没有有效凭证，跳过
placeholder

	// 检查 token 是否过期，过期则刷新
	if r.isTokenExpired(account) {
		tokenInfo, err := r.oauthSvc.RefreshAccountToken(ctx, account)
		if err != nil {
			return err
	placeholder
		accessToken = tokenInfo.AccessToken
		// 更新凭证
		account.Credentials = r.oauthSvc.BuildAccountCredentials(tokenInfo)
placeholder

	// 获取代理 URL
	var proxyURL string
	if account.ProxyID != nil {
		proxy, err := r.proxyRepo.GetByID(ctx, *account.ProxyID)
		if err == nil && proxy != nil {
			proxyURL = proxy.URL()
	placeholder
placeholder

	// 调用 API 获取配额
	client := antigravity.NewClient(proxyURL)
	modelsResp, err := client.FetchAvailableModels(ctx, accessToken, projectID)
	if err != nil {
		return err
placeholder

	// 解析配额数据并更新 extra 字段
	r.updateAccountQuota(account, modelsResp)

	// 保存到数据库
	return r.accountRepo.Update(ctx, account)
placeholder

// isTokenExpired 检查 token 是否过期
func (r *AntigravityQuotaRefresher) isTokenExpired(account *Account) bool {
	expiresAt := parseAntigravityExpiresAt(account)
	if expiresAt == nil {
		return false
placeholder

	// 提前 5 分钟认为过期
	return time.Now().Add(5 * time.Minute).After(*expiresAt)
placeholder

// updateAccountQuota 更新账户的配额信息
func (r *AntigravityQuotaRefresher) updateAccountQuota(account *Account, modelsResp *antigravity.FetchAvailableModelsResponse) {
	if account.Extra == nil {
		account.Extra = make(map[string]any)
placeholder

	quota := make(map[string]any)

	for modelName, modelInfo := range modelsResp.Models {
		if modelInfo.QuotaInfo == nil {
			continue
	placeholder

		// 转换 remainingFraction (0.0-1.0) 为百分比 (0-100)
		remaining := int(modelInfo.QuotaInfo.RemainingFraction * 100)

		quota[modelName] = map[string]any{
			"remaining":  remaining,
			"reset_time": modelInfo.QuotaInfo.ResetTime,
	placeholder
placeholder

	account.Extra["quota"] = quota
	account.Extra["last_quota_check"] = time.Now().Format(time.RFC3339)
placeholder
