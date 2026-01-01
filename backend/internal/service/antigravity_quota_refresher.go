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
	cfg         *config.TokenRefreshConfig

	stopCh chan struct{placeholder
	wg     sync.WaitGroup
placeholder

// NewAntigravityQuotaRefresher 创建配额刷新器
func NewAntigravityQuotaRefresher(
	accountRepo AccountRepository,
	proxyRepo ProxyRepository,
	_ *AntigravityOAuthService,
	cfg *config.Config,
) *AntigravityQuotaRefresher {
	return &AntigravityQuotaRefresher{
		accountRepo: accountRepo,
		proxyRepo:   proxyRepo,
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

	if accessToken == "" {
		return nil // 没有 access_token，跳过
placeholder

	// token 过期则跳过，由 TokenRefreshService 负责刷新
	if r.isTokenExpired(account) {
		return nil
placeholder

	// 获取代理 URL
	var proxyURL string
	if account.ProxyID != nil {
		proxy, err := r.proxyRepo.GetByID(ctx, *account.ProxyID)
		if err == nil && proxy != nil {
			proxyURL = proxy.URL()
	placeholder
placeholder

	client := antigravity.NewClient(proxyURL)

	if account.Extra == nil {
		account.Extra = make(map[string]any)
placeholder

	// 获取账户信息（tier、project_id 等）
	loadResp, loadRaw, _ := client.LoadCodeAssist(ctx, accessToken)
	if loadRaw != nil {
		account.Extra["load_code_assist"] = loadRaw
placeholder
	if loadResp != nil {
		// 尝试从 API 获取 project_id
		if projectID == "" && loadResp.CloudAICompanionProject != "" {
			projectID = loadResp.CloudAICompanionProject
			account.Credentials["project_id"] = projectID
	placeholder
placeholder

	// 如果仍然没有 project_id，随机生成一个并保存
	if projectID == "" {
		projectID = antigravity.GenerateMockProjectID()
		account.Credentials["project_id"] = projectID
		log.Printf("[AntigravityQuotaRefresher] 为账户 %d 生成随机 project_id: %s", account.ID, projectID)
placeholder

	// 调用 API 获取配额
	modelsResp, modelsRaw, err := client.FetchAvailableModels(ctx, accessToken, projectID)
	if err != nil {
		return r.accountRepo.Update(ctx, account) // 保存已有的 load_code_assist 信息
placeholder

	// 保存完整的配额响应
	if modelsRaw != nil {
		account.Extra["available_models"] = modelsRaw
placeholder

	// 解析配额数据为前端使用的格式
	r.updateAccountQuota(account, modelsResp)

	account.Extra["last_refresh"] = time.Now().Format(time.RFC3339)

	// 保存到数据库
	return r.accountRepo.Update(ctx, account)
placeholder

// isTokenExpired 检查 token 是否过期
func (r *AntigravityQuotaRefresher) isTokenExpired(account *Account) bool {
	expiresAt := account.GetCredentialAsTime("expires_at")
	if expiresAt == nil {
		return false
placeholder

	// 提前 5 分钟认为过期
	return time.Now().Add(5 * time.Minute).After(*expiresAt)
placeholder

// updateAccountQuota 更新账户的配额信息（前端使用的格式）
func (r *AntigravityQuotaRefresher) updateAccountQuota(account *Account, modelsResp *antigravity.FetchAvailableModelsResponse) {
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
placeholder
