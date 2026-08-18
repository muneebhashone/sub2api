package service

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/domain"
	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"golang.org/x/sync/singleflight"
)

// 渠道监控「配额模式」的配额抓取器。
//
// 不直接对接上游，而是把账号侧现成的用量服务归一成 domain.MonitorQuotaSnapshot：
//   - 海外 5 家（anthropic/openai/gemini/antigravity/grok）→ AccountUsageService.GetUsage
//   - 国产 coding plan（kimi/zhipu/deepseek）→ CNProviderQuotaService.QueryUsage
//   - 国产 payg（kimi/deepseek）→ CNProviderBalanceService.QueryBalance
//     （zhipu payg 无公开余额端点，QueryBalance 会返回该错误，原样透出）
//
// Fetch 永不返回 error：所有失败都降级为 Success=false 的快照照常入库，
// 由 deriveQuotaCheckResult 推导为 failed/error 状态。
//
// 多个监控可能关联同一账号，而 interval 最小 15s 且国产配额服务自身无缓存，
// 所以快照统一带 TTL 缓存（成功 monitorQuotaFetchCacheTTL、失败
// monitorQuotaErrorCacheTTL 负缓存），防止打爆上游配额端点；同账号的并发
// 抓取由 singleflight 合并为一次上游查询。

// monitorUsageSource 海外平台账号用量查询（AccountUsageService 天然满足）。
type monitorUsageSource interface {
	GetUsage(ctx context.Context, accountID int64, force ...bool) (*UsageInfo, error)
placeholder

// monitorCNQuotaSource 国产 coding plan 滚动窗口额度探测（CNProviderQuotaService 天然满足）。
type monitorCNQuotaSource interface {
	QueryUsage(ctx context.Context, accountID int64) (*CNProviderQuotaProbeResult, error)
placeholder

// monitorCNBalanceSource 国产 payg 余额探测（CNProviderBalanceService 天然满足）。
type monitorCNBalanceSource interface {
	QueryBalance(ctx context.Context, accountID int64) (*CNProviderBalanceResult, error)
placeholder

// monitorAccountSource 账号加载（AccountRepository 天然满足）。
type monitorAccountSource interface {
	GetByID(ctx context.Context, id int64) (*Account, error)
placeholder

// ChannelMonitorQuotaFetcher 配额抓取器（成功/失败快照均带 TTL 缓存，
// 同账号并发抓取由 singleflight 合并）。
type ChannelMonitorQuotaFetcher struct {
	usage     monitorUsageSource
	cnQuota   monitorCNQuotaSource
	cnBalance monitorCNBalanceSource
	accounts  monitorAccountSource

	mu     sync.Mutex
	cache  map[int64]monitorQuotaCacheEntry
	flight singleflight.Group
placeholder

type monitorQuotaCacheEntry struct {
	snapshot *domain.MonitorQuotaSnapshot
	expiry   time.Time
placeholder

// NewChannelMonitorQuotaFetcher 构造配额抓取器。
// 参数取具体服务类型以便 wire 直连；单元测试在同包内用 struct 字面量注入 stub。
func NewChannelMonitorQuotaFetcher(
	usage *AccountUsageService,
	cnQuota *CNProviderQuotaService,
	cnBalance *CNProviderBalanceService,
	accounts AccountRepository,
) *ChannelMonitorQuotaFetcher {
	f := &ChannelMonitorQuotaFetcher{cache: make(map[int64]monitorQuotaCacheEntry)placeholder
	if usage != nil {
		f.usage = usage
placeholder
	if cnQuota != nil {
		f.cnQuota = cnQuota
placeholder
	if cnBalance != nil {
		f.cnBalance = cnBalance
placeholder
	if accounts != nil {
		f.accounts = accounts
placeholder
	return f
placeholder

// LoadAccount 加载账号（不走缓存）。供 Create/Update 时校验
// provider 与 account.platform 一致；账号不存在时返回错误。
func (f *ChannelMonitorQuotaFetcher) LoadAccount(ctx context.Context, id int64) (*Account, error) {
	if f == nil || f.accounts == nil {
		return nil, fmt.Errorf("quota fetcher is not configured")
placeholder
	return f.accounts.GetByID(ctx, id)
placeholder

// Fetch 抓取账号的最新配额快照。永不返回 error：失败降级为
// Success=false 快照（Error 带摘要），保证检测历史的时间线连续。
func (f *ChannelMonitorQuotaFetcher) Fetch(ctx context.Context, accountID int64) *domain.MonitorQuotaSnapshot {
	if f == nil {
		// fail-closed：fetcher 未注入（存量测试构造）时不 panic，降级为错误快照。
		return quotaErrorSnapshot("usage", "quota fetcher is not configured", time.Now())
placeholder

	now := time.Now()

	if cached, ok := f.cachedSnapshot(accountID, now); ok {
		return cached
placeholder

	// singleflight 合并同账号并发抓取；脱离调用方 ctx（仿 CN 配额服务），
	// 避免某个监控的取消波及共享同一账号的其他监控。
	key := "monitor-quota:" + strconv.FormatInt(accountID, 10)
	ch := f.flight.DoChan(key, func() (any, error) {
		fetchCtx, cancel := context.WithTimeout(context.Background(), monitorQuotaFetchTimeout)
		defer cancel()
		snapshot := f.fetchUncached(fetchCtx, accountID, time.Now())
		// 失败也进短 TTL 负缓存：凭据失效/故障期间不必每次调度都打上游。
		ttl := monitorQuotaFetchCacheTTL
		if !snapshot.Success {
			ttl = monitorQuotaErrorCacheTTL
	placeholder
		f.storeSnapshot(accountID, snapshot, time.Now().Add(ttl))
		return snapshot, nil
placeholder)
	select {
	case <-ctx.Done():
		return quotaErrorSnapshot("usage", "context canceled", now)
	case res := <-ch:
		snapshot, ok := res.Val.(*domain.MonitorQuotaSnapshot)
		if res.Err != nil || !ok || snapshot == nil {
			return quotaErrorSnapshot("usage", "quota fetch failed", now)
	placeholder
		return snapshot
placeholder
placeholder

func (f *ChannelMonitorQuotaFetcher) cachedSnapshot(accountID int64, now time.Time) (*domain.MonitorQuotaSnapshot, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	entry, ok := f.cache[accountID]
	if !ok || now.After(entry.expiry) {
		return nil, false
placeholder
	return entry.snapshot, true
placeholder

func (f *ChannelMonitorQuotaFetcher) storeSnapshot(accountID int64, snapshot *domain.MonitorQuotaSnapshot, expiry time.Time) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.cache[accountID] = monitorQuotaCacheEntry{snapshot: snapshot, expiry: expiryplaceholder
placeholder

func (f *ChannelMonitorQuotaFetcher) fetchUncached(ctx context.Context, accountID int64, now time.Time) *domain.MonitorQuotaSnapshot {
	if f == nil {
		return quotaErrorSnapshot("usage", "quota fetcher is not configured", now)
placeholder

	account, err := f.LoadAccount(ctx, accountID)
	if err != nil || account == nil {
		// FK ON DELETE SET NULL 后 account_id 可能为空/失效；显式报「账号未关联」，
		// 推导为 degraded（配置问题，不是渠道故障）。
		slog.Warn("channel_monitor: load linked account failed",
			"account_id", accountID, "error", err)
		return quotaErrorSnapshot("usage", "linked account not found", now)
placeholder

	switch account.Platform {
	case domain.PlatformKimi, domain.PlatformZhipu, domain.PlatformDeepseek:
		if account.IsCodingPlan() {
			return f.fetchCNQuota(ctx, accountID, now)
	placeholder
		return f.fetchCNBalance(ctx, accountID, now)
	default:
		return f.fetchUsage(ctx, accountID, now)
placeholder
placeholder

// fetchUsage 海外平台：AccountUsageService.GetUsage → 快照。
func (f *ChannelMonitorQuotaFetcher) fetchUsage(ctx context.Context, accountID int64, now time.Time) *domain.MonitorQuotaSnapshot {
	if f.usage == nil {
		return quotaErrorSnapshot("usage", "usage service is not configured", now)
placeholder
	usage, err := f.usage.GetUsage(ctx, accountID)
	if err != nil {
		msg := truncateMessage(sanitizeErrorMessage(err.Error()))
		return &domain.MonitorQuotaSnapshot{
			Source:            "usage",
			Success:           false,
			CredentialInvalid: isCredentialErrorMessage(msg),
			Error:             msg,
			FetchedAt:         now,
	placeholder
placeholder
	if usage == nil {
		return quotaErrorSnapshot("usage", "usage service returned no data", now)
placeholder
	// openai/gemini/antigravity/grok 的失败多走「值通道」（err==nil 但错误
	// 降级在 UsageInfo 字段里），必须显式识别，否则会被误判为 operational。
	if failed, credInvalid, msg := usageFailureInfo(usage); failed {
		return &domain.MonitorQuotaSnapshot{
			Source:            "usage",
			Success:           false,
			CredentialInvalid: credInvalid,
			Error:             truncateMessage(sanitizeErrorMessage(msg)),
			FetchedAt:         now,
	placeholder
placeholder
	snapshot := &domain.MonitorQuotaSnapshot{
		Source:    "usage",
		Success:   true,
		PlanLevel: usage.SubscriptionTier,
		Tiers:     usageQuotaTiers(usage),
		FetchedAt: now,
placeholder
	if snapshot.PlanLevel == "" {
		snapshot.PlanLevel = usage.SubscriptionTierRaw
placeholder
	return snapshot
placeholder

// usageQuotaTiers 把 UsageInfo 的各平台窗口归一为 tier 列表（无数据的窗口跳过）。
func usageQuotaTiers(usage *UsageInfo) []domain.MonitorQuotaTier {
	if usage == nil {
		return nil
placeholder
	tiers := make([]domain.MonitorQuotaTier, 0, 8)
	appendProgressTier(&tiers, "5h", "", usage.FiveHour)
	appendProgressTier(&tiers, "7d", "", usage.SevenDay)
	appendProgressTier(&tiers, "7d-sonnet", "", usage.SevenDaySonnet)
	appendProgressTier(&tiers, "7d-fable", "", usage.SevenDayFable)
	appendProgressTier(&tiers, "30d", "", usage.ThirtyDay)
	// Gemini 多档日配额：同 Window 不同 Label。
	appendProgressTier(&tiers, "daily", "shared", usage.GeminiSharedDaily)
	appendProgressTier(&tiers, "daily", "pro", usage.GeminiProDaily)
	appendProgressTier(&tiers, "daily", "flash", usage.GeminiFlashDaily)
	// Grok requests/tokens 两个日窗口 + 月度计费窗口。
	appendQuotaWindowTier(&tiers, "daily", "requests", usage.GrokRequestQuota)
	appendQuotaWindowTier(&tiers, "daily", "tokens", usage.GrokTokenQuota)
	// Antigravity per-model 总量额度，Label = 模型名（按名排序保证输出稳定）。
	for _, model := range sortedQuotaModelNames(usage.AntigravityQuota) {
		q := usage.AntigravityQuota[model]
		if q == nil {
			continue
	placeholder
		tiers = append(tiers, domain.MonitorQuotaTier{
			Window:      "total",
			Label:       model,
			UsedPercent: float64(q.Utilization),
			ResetAt:     q.ResetTime,
	placeholder)
placeholder
	if len(tiers) == 0 {
		return nil
placeholder
	return tiers
placeholder

func appendProgressTier(tiers *[]domain.MonitorQuotaTier, window, label string, p *UsageProgress) {
	if p == nil {
		return
placeholder
	tier := domain.MonitorQuotaTier{
		Window:      window,
		Label:       label,
		UsedPercent: p.Utilization,
placeholder
	if p.ResetsAt != nil {
		tier.ResetAt = p.ResetsAt.UTC().Format(time.RFC3339)
placeholder
	if p.LimitRequests > 0 {
		tier.Used = float64(p.UsedRequests)
		tier.Limit = float64(p.LimitRequests)
placeholder
	*tiers = append(*tiers, tier)
placeholder

func appendQuotaWindowTier(tiers *[]domain.MonitorQuotaTier, window, label string, q *xai.QuotaWindow) {
	if q == nil || q.Limit == nil || *q.Limit <= 0 {
		return
placeholder
	used := *q.Limit
	if q.Remaining != nil {
		used = *q.Limit - *q.Remaining
		if used < 0 {
			used = 0
	placeholder
placeholder
	tier := domain.MonitorQuotaTier{
		Window:      window,
		Label:       label,
		Used:        float64(used),
		Limit:       float64(*q.Limit),
		UsedPercent: float64(used) / float64(*q.Limit) * 100,
placeholder
	if q.ResetAt != "" {
		tier.ResetAt = q.ResetAt
placeholder else if q.ResetUnix != nil && *q.ResetUnix > 0 {
		tier.ResetAt = time.Unix(*q.ResetUnix, 0).UTC().Format(time.RFC3339)
placeholder
	*tiers = append(*tiers, tier)
placeholder

func sortedQuotaModelNames(quotas map[string]*AntigravityModelQuota) []string {
	names := make([]string, 0, len(quotas))
	for name := range quotas {
		names = append(names, name)
placeholder
	sort.Strings(names)
	return names
placeholder

// fetchCNQuota 国产 coding plan：CNProviderQuotaService.QueryUsage → 快照。
func (f *ChannelMonitorQuotaFetcher) fetchCNQuota(ctx context.Context, accountID int64, now time.Time) *domain.MonitorQuotaSnapshot {
	if f.cnQuota == nil {
		return quotaErrorSnapshot("cn_quota", "cn quota service is not configured", now)
placeholder
	result, err := f.cnQuota.QueryUsage(ctx, accountID)
	if err != nil {
		msg := truncateMessage(sanitizeErrorMessage(err.Error()))
		return &domain.MonitorQuotaSnapshot{
			Source:            "cn_quota",
			Success:           false,
			CredentialInvalid: isCredentialErrorMessage(msg),
			Error:             msg,
			FetchedAt:         now,
	placeholder
placeholder
	snapshot := &domain.MonitorQuotaSnapshot{
		Source:    "cn_quota",
		Success:   result.Success,
		PlanLevel: result.PlanLevel,
		Error:     result.Error,
		FetchedAt: now,
placeholder
	if !result.Success && !result.CredentialValid {
		snapshot.CredentialInvalid = true
placeholder
	if len(result.Tiers) > 0 {
		snapshot.Tiers = make([]domain.MonitorQuotaTier, 0, len(result.Tiers))
		for _, t := range result.Tiers {
			snapshot.Tiers = append(snapshot.Tiers, domain.MonitorQuotaTier{
				Window:      t.Window,
				UsedPercent: t.UsedPercent,
				ResetAt:     t.ResetAt,
		placeholder)
	placeholder
placeholder
	if !snapshot.Success {
		snapshot.Error = firstNonEmpty(snapshot.Error, "cn quota probe failed")
placeholder
	return snapshot
placeholder

// fetchCNBalance 国产 payg：CNProviderBalanceService.QueryBalance → 快照。
func (f *ChannelMonitorQuotaFetcher) fetchCNBalance(ctx context.Context, accountID int64, now time.Time) *domain.MonitorQuotaSnapshot {
	if f.cnBalance == nil {
		return quotaErrorSnapshot("cn_balance", "cn balance service is not configured", now)
placeholder
	result, err := f.cnBalance.QueryBalance(ctx, accountID)
	if err != nil {
		msg := truncateMessage(sanitizeErrorMessage(err.Error()))
		return &domain.MonitorQuotaSnapshot{
			Source:            "cn_balance",
			Success:           false,
			CredentialInvalid: isCredentialErrorMessage(msg),
			Error:             msg,
			FetchedAt:         now,
	placeholder
placeholder
	snapshot := &domain.MonitorQuotaSnapshot{
		Source:    "cn_balance",
		Success:   result.Success,
		Currency:  result.Currency,
		Error:     result.Error,
		FetchedAt: now,
placeholder
	if result.Success {
		balance := result.Balance
		snapshot.Balance = &balance
placeholder else if result.StatusCode == 401 || result.StatusCode == 403 {
		snapshot.CredentialInvalid = true
placeholder
	if len(result.Balances) > 0 {
		snapshot.Balances = make([]domain.MonitorBalance, 0, len(result.Balances))
		for _, b := range result.Balances {
			snapshot.Balances = append(snapshot.Balances, domain.MonitorBalance{
				Currency: b.Currency,
				Balance:  b.Balance,
		placeholder)
	placeholder
placeholder
	if !snapshot.Success {
		snapshot.Error = firstNonEmpty(snapshot.Error, "cn balance probe failed")
placeholder
	return snapshot
placeholder

// quotaErrorSnapshot 构造统一错误快照。
func quotaErrorSnapshot(source, message string, now time.Time) *domain.MonitorQuotaSnapshot {
	return &domain.MonitorQuotaSnapshot{
		Source:    source,
		Success:   false,
		Error:     truncateMessage(sanitizeErrorMessage(message)),
		FetchedAt: now,
placeholder
placeholder

// isCredentialErrorMessage 上游 401/403 鉴权失败的启发式识别
// （海外 GetUsage 的错误没有结构化状态码，只能看文本）。
func isCredentialErrorMessage(msg string) bool {
	msg = strings.ToLower(msg)
	return strings.Contains(msg, "401") ||
		strings.Contains(msg, "403") ||
		strings.Contains(msg, "unauthorized") ||
		strings.Contains(msg, "forbidden") ||
		strings.Contains(msg, "invalid_api_key") ||
		strings.Contains(msg, "authentication")
placeholder

// usageFailureInfo 识别 GetUsage 经「值通道」返回的失败：antigravity/grok
// 等平台 err==nil 但把错误降级在 UsageInfo 字段里（Error/ErrorCode/状态标记）。
// 返回 failed=false 表示可用；credentialInvalid 表示凭据失效（401/403 语义，
// 推导为 failed 状态）；msg 为失败摘要。
//
// grok 的 ErrorCode=quota_unknown 是「尚未观测到计费快照/限流头」的已知未知态，
// 不是失败（严格按 ErrorCode 判会把健康 grok 账号永久判 error），显式豁免。
func usageFailureInfo(usage *UsageInfo) (failed, credentialInvalid bool, msg string) {
	if usage == nil {
		return false, false, ""
placeholder
	if usage.ErrorCode == "quota_unknown" {
		return false, false, ""
placeholder
	failed = usage.Error != "" || usage.NeedsReauth || usage.IsBanned ||
		usage.IsForbidden || usage.ErrorCode != ""
	if !failed {
		return false, false, ""
placeholder
	credentialInvalid = usage.NeedsReauth || usage.IsBanned || usage.IsForbidden ||
		usage.ErrorCode == errorCodeUnauthenticated || usage.ErrorCode == errorCodeForbidden
	msg = firstNonEmpty(usage.Error, usage.ForbiddenReason, usage.ErrorCode, "usage fetch failed")
	return failed, credentialInvalid, msg
placeholder

// deriveQuotaCheckResult 把配额快照推导为检测状态（复用既有 status 枚举，
// 时间线/可用率机制自动生效）：
//   - 查询成功且无告警        → operational
//   - 任一窗口使用率 >= 阈值或余额耗尽 → degraded
//   - 账号未关联（配置问题）    → degraded
//   - 凭据失效（401/403）     → failed
//   - 网络/解析等其他错误      → error
func deriveQuotaCheckResult(snapshot *domain.MonitorQuotaSnapshot, model string, checkedAt time.Time) *CheckResult {
	res := &CheckResult{Model: model, CheckedAt: checkedAtplaceholder
	if snapshot == nil {
		res.Status = MonitorStatusError
		res.Message = "quota snapshot missing"
		return res
placeholder

	switch {
	case !snapshot.Success && snapshot.CredentialInvalid:
		res.Status = MonitorStatusFailed
		res.Message = snapshot.Error
	case !snapshot.Success && strings.Contains(snapshot.Error, "linked account not found"):
		res.Status = MonitorStatusDegraded
		res.Message = snapshot.Error
	case !snapshot.Success:
		res.Status = MonitorStatusError
		res.Message = snapshot.Error
	default:
		if hint := quotaDegradedHint(snapshot); hint != "" {
			res.Status = MonitorStatusDegraded
			res.Message = hint
	placeholder else {
			res.Status = MonitorStatusOperational
	placeholder
placeholder
	return res
placeholder

// quotaDegradedHint 生成 degraded 的 message（指出触发告警的窗口/余额）；
// 空串表示无告警。
func quotaDegradedHint(snapshot *domain.MonitorQuotaSnapshot) string {
	for _, tier := range snapshot.Tiers {
		if tier.UsedPercent >= monitorQuotaDegradedUsedPercent {
			name := tier.Window
			if tier.Label != "" {
				name = tier.Label + "/" + tier.Window
		placeholder
			return fmt.Sprintf("quota high: %s at %s%%", name, strconv.FormatFloat(tier.UsedPercent, 'f', 1, 64))
	placeholder
placeholder
	if snapshot.Balance != nil && *snapshot.Balance <= 0 {
		return fmt.Sprintf("balance depleted (%s)", firstNonEmpty(snapshot.Currency, "?"))
placeholder
	return ""
placeholder
