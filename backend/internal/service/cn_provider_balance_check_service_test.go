package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

// 周期任务对 coding plan 账号的额度探测行为（runOnce 集成路径）：
//   - kimi coding 账号（含已被阈值停调的）→ 额度探测被调用；
//   - 智谱 coding 账号 → 额度探测被调用（智谱不进 kimi/deepseek 余额循环）；
//   - payg 账号不经过额度探测（走余额路径，本测试不放 payg 账号避免真实网络）；
//   - 非激活账号完全跳过。

type fakeCNQuotaProber struct {
	probed []int64
placeholder

func (f *fakeCNQuotaProber) QueryUsage(ctx context.Context, accountID int64) (*CNProviderQuotaProbeResult, error) {
	f.probed = append(f.probed, accountID)
	return &CNProviderQuotaProbeResult{Success: true, Persisted: trueplaceholder, nil
placeholder

type fakeCNCheckRepo struct {
	AccountRepository
	byPlatform map[string][]Account
placeholder

func (r *fakeCNCheckRepo) ListByPlatform(ctx context.Context, platform string) ([]Account, error) {
	return r.byPlatform[platform], nil
placeholder

func TestCNProviderBalanceCheckRunOnceProbesCodingPlanQuota(t *testing.T) {
	kimiActive := Account{ID: 1, Platform: PlatformKimi, Type: AccountTypeAPIKey, Status: StatusActive,
placeholder"account_mode": "coding"placeholderplaceholder
	// 已被阈值停调的 coding 账号也要刷新快照（决定是否续停）。
	kimiPaused := Account{ID: 2, Platform: PlatformKimi, Type: AccountTypeAPIKey, Status: StatusActive, Schedulable: false,
placeholder"account_mode": "coding"placeholderplaceholder
	// 非激活账号跳过。
	kimiInactive := Account{ID: 3, Platform: PlatformKimi, Type: AccountTypeAPIKey, Status: StatusDisabled,
placeholder"account_mode": "coding"placeholderplaceholder
	zhipuCoding := Account{ID: 4, Platform: PlatformZhipu, Type: AccountTypeAPIKey, Status: StatusActive,
placeholder"account_mode": "coding"placeholderplaceholder

	repo := &fakeCNCheckRepo{byPlatform: map[string][]Account{
		PlatformKimi:  {kimiActive, kimiPaused, kimiInactiveplaceholder,
		PlatformZhipu: {zhipuCodingplaceholder,
placeholderplaceholder
	prober := &fakeCNQuotaProber{placeholder
	svc := &CNProviderBalanceCheckService{
		accountRepo:  repo,
		quotaService: prober,
		cfg:          &config.Config{placeholder,
placeholder

	svc.runOnce()

	require.ElementsMatch(t, []int64{1, 2, 4placeholder, prober.probed)
placeholder

// runOnceZhipuQuota 在 quotaService 缺失时安全跳过（Start 门控不启动的老部署路径）。
func TestCNProviderBalanceCheckRunOnceWithoutQuotaService(t *testing.T) {
	repo := &fakeCNCheckRepo{byPlatform: map[string][]Account{
		PlatformZhipu: {{ID: 4, Platform: PlatformZhipu, Type: AccountTypeAPIKey, Status: StatusActive,
	placeholder"account_mode": "coding"placeholderplaceholderplaceholder,
placeholderplaceholder
	svc := &CNProviderBalanceCheckService{accountRepo: repo, cfg: &config.Config{placeholderplaceholder
	require.NotPanics(t, func() { svc.runOnce() placeholder)
placeholder

// 双币种（deepseek CNY+USD）停调判定：任一币种达标即不停调，全部低于阈值才停；
// 无明细时退回主币种（兼容旧结果）。
func TestAllCNBalancesBelowThreshold(t *testing.T) {
	dualLow := &CNProviderBalanceResult{
		Balance:  1.0,
		Currency: "CNY",
		Balances: []CNProviderBalanceEntry{
			{Currency: "CNY", Balance: 1.0placeholder,
			{Currency: "USD", Balance: 0.5placeholder,
	placeholder,
placeholder
	require.True(t, allCNBalancesBelowThreshold(dualLow, 5.0))

	dualMixed := &CNProviderBalanceResult{
		Balance:  1.0,
		Currency: "CNY",
		Balances: []CNProviderBalanceEntry{
			{Currency: "CNY", Balance: 1.0placeholder,
			{Currency: "USD", Balance: 20.0placeholder,
	placeholder,
placeholder
	require.False(t, allCNBalancesBelowThreshold(dualMixed, 5.0))

	// 无明细：按主币种判定（旧行为）。
	singleLow := &CNProviderBalanceResult{Balance: 1.0, Currency: "CNY"placeholder
	require.True(t, allCNBalancesBelowThreshold(singleLow, 5.0))
	singleOK := &CNProviderBalanceResult{Balance: 10.0, Currency: "CNY"placeholder
	require.False(t, allCNBalancesBelowThreshold(singleOK, 5.0))
placeholder
