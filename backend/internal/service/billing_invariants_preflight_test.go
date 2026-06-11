//go:build unit

// TASK-003 preflight 配额/余额拒绝不变量测试（INVARIANTS.md I-3.6）。
//
// CheckBillingEligibility 是网关入口的统一计费资格检查；handler 在
// AcquireUserSlotWithWait 等待结束后会用同一函数做"二次检查"
// （internal/handler/gateway_handler.go:255），因此这里锁定其无状态语义：
//   - 余额模式：余额 <= 0 → ErrInsufficientBalance
//   - 订阅模式：日/周/月用量达到分组限额 → ErrDaily/Weekly/MonthlyLimitExceeded；
//     订阅过期或非 active → ErrSubscriptionInvalid
//   - user×platform 配额：日限额耗尽 → ErrUserPlatformDailyQuotaExhausted（429）；
//     订阅模式豁免该检查
//   - 并发等待后二次检查：第一次放行后用量/余额变化，再次调用即拒绝
//   - simple 运行模式跳过所有计费检查
package service

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

// billInvBillingCacheStub 只实现 CheckBillingEligibility 路径会触达的读方法，
// 其余方法走嵌入接口（不会被调用）。
type billInvBillingCacheStub struct {
	BillingCache

	balance    float64
	sub        *SubscriptionCacheData
	quotaEntry *UserPlatformQuotaCacheEntry
placeholder

func (s *billInvBillingCacheStub) GetUserBalance(ctx context.Context, userID int64) (float64, error) {
	return s.balance, nil
placeholder

func (s *billInvBillingCacheStub) GetSubscriptionCache(ctx context.Context, userID, groupID int64) (*SubscriptionCacheData, error) {
	return s.sub, nil
placeholder

func (s *billInvBillingCacheStub) GetUserPlatformQuotaCache(ctx context.Context, userID int64, platform string) (*UserPlatformQuotaCacheEntry, bool, error) {
	if s.quotaEntry == nil {
		return nil, false, nil
placeholder
	return s.quotaEntry, true, nil
placeholder

// billInvQuotaRepoStub 仅用于让 userPlatformQuotaRepo 非 nil（cache HIT 路径不
// 会触达 DB），嵌入接口的其余方法不会被调用。
type billInvQuotaRepoStub struct {
	UserPlatformQuotaRepository
placeholder

func (s *billInvQuotaRepoStub) GetByUserPlatform(ctx context.Context, userID int64, platform string) (*UserPlatformQuotaRecord, error) {
	return nil, nil
placeholder

func billInvNewBillingCacheService(t *testing.T, cache BillingCache, cfg *config.Config) *BillingCacheService {
placeholder
	if cfg == nil {
		cfg = &config.Config{placeholder
placeholder
	svc := NewBillingCacheService(cache, nil, nil, nil, nil, nil, cfg, &billInvQuotaRepoStub{placeholder)
	t.Cleanup(svc.Stop)
	return svc
placeholder

// billInvQuotaEntryV1 构造当前窗口内的 SchemaV1 user×platform 配额缓存条目。
func billInvQuotaEntryV1(dailyLimit, dailyUsage float64) *UserPlatformQuotaCacheEntry {
	now := time.Now()
	return &UserPlatformQuotaCacheEntry{
		SchemaVersion:      UserPlatformQuotaCacheSchemaV1,
		DailyLimitUSD:      &dailyLimit,
		DailyUsageUSD:      dailyUsage,
		DailyWindowStart:   &now,
		WeeklyWindowStart:  &now,
		MonthlyWindowStart: &now,
placeholder
placeholder

// TestBillingInvariant_PreflightBalanceEligibility 锁定余额模式 preflight 语义。
func TestBillingInvariant_PreflightBalanceEligibility(t *testing.T) {
	user := &User{ID: 601placeholder

	t.Run("余额耗尽拒绝", func(t *testing.T) {
		svc := billInvNewBillingCacheService(t, &billInvBillingCacheStub{balance: 0placeholder, nil)
		err := svc.CheckBillingEligibility(context.Background(), user, nil, nil, nil, "")
		require.ErrorIs(t, err, ErrInsufficientBalance)
placeholder)

	t.Run("余额为正放行", func(t *testing.T) {
		svc := billInvNewBillingCacheService(t, &billInvBillingCacheStub{balance: 5.0placeholder, nil)
		err := svc.CheckBillingEligibility(context.Background(), user, nil, nil, nil, "")
	placeholder
placeholder)

	t.Run("并发等待后二次检查反映最新余额", func(t *testing.T) {
		cache := &billInvBillingCacheStub{balance: 0.01placeholder
		svc := billInvNewBillingCacheService(t, cache, nil)

		// 第一次检查（获取并发槽前）：余额尚存 → 放行
		require.NoError(t, svc.CheckBillingEligibility(context.Background(), user, nil, nil, nil, ""))

		// 等待期间其他请求把余额扣到 0 → 等待结束后的二次检查必须拒绝
		cache.balance = 0
		err := svc.CheckBillingEligibility(context.Background(), user, nil, nil, nil, "")
		require.ErrorIs(t, err, ErrInsufficientBalance)
placeholder)
placeholder

// TestBillingInvariant_PreflightSubscriptionLimits 锁定订阅模式 preflight 语义：
// 日/周/月任一窗口用量达到分组限额即拒绝；订阅非 active 或已过期拒绝。
func TestBillingInvariant_PreflightSubscriptionLimits(t *testing.T) {
	user := &User{ID: 601placeholder
	subscription := &UserSubscription{ID: 42placeholder
	group := &Group{
		ID:               7,
		SubscriptionType: SubscriptionTypeSubscription,
		DailyLimitUSD:    billInvF64Ptr(10),
		WeeklyLimitUSD:   billInvF64Ptr(50),
		MonthlyLimitUSD:  billInvF64Ptr(100),
placeholder
	activeFuture := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name    string
		sub     *SubscriptionCacheData
		wantErr error
placeholder{
		{
			name:    "限额内放行",
			sub:     &SubscriptionCacheData{Status: SubscriptionStatusActive, ExpiresAt: activeFuture, DailyUsage: 9.99, WeeklyUsage: 49.99, MonthlyUsage: 99.99placeholder,
			wantErr: nil,
	placeholder,
		{
			name:    "日限额达到拒绝",
			sub:     &SubscriptionCacheData{Status: SubscriptionStatusActive, ExpiresAt: activeFuture, DailyUsage: 10placeholder,
			wantErr: ErrDailyLimitExceeded,
	placeholder,
		{
			name:    "周限额达到拒绝",
			sub:     &SubscriptionCacheData{Status: SubscriptionStatusActive, ExpiresAt: activeFuture, WeeklyUsage: 50placeholder,
			wantErr: ErrWeeklyLimitExceeded,
	placeholder,
		{
			name:    "月限额达到拒绝",
			sub:     &SubscriptionCacheData{Status: SubscriptionStatusActive, ExpiresAt: activeFuture, MonthlyUsage: 100placeholder,
			wantErr: ErrMonthlyLimitExceeded,
	placeholder,
		{
			name:    "订阅过期拒绝",
			sub:     &SubscriptionCacheData{Status: SubscriptionStatusActive, ExpiresAt: time.Now().Add(-time.Minute)placeholder,
			wantErr: ErrSubscriptionInvalid,
	placeholder,
		{
			name:    "订阅非active拒绝",
			sub:     &SubscriptionCacheData{Status: "cancelled", ExpiresAt: activeFutureplaceholder,
			wantErr: ErrSubscriptionInvalid,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := billInvNewBillingCacheService(t, &billInvBillingCacheStub{sub: tt.subplaceholder, nil)
			err := svc.CheckBillingEligibility(context.Background(), user, nil, group, subscription, "")
			if tt.wantErr == nil {
			placeholder
		placeholder else {
				require.ErrorIs(t, err, tt.wantErr)
		placeholder
	placeholder)
placeholder
placeholder

// TestBillingInvariant_PreflightUserPlatformQuota 锁定 user×platform 配额的
// preflight 语义：余额模式下日限额耗尽 → 429 拒绝；订阅模式豁免该检查。
func TestBillingInvariant_PreflightUserPlatformQuota(t *testing.T) {
	user := &User{ID: 601placeholder

	t.Run("日配额耗尽拒绝", func(t *testing.T) {
		cache := &billInvBillingCacheStub{
			balance:    5.0, // 余额充足，确保拒绝来自 platform quota
			quotaEntry: billInvQuotaEntryV1(5.0, 5.0),
	placeholder
		svc := billInvNewBillingCacheService(t, cache, nil)
		err := svc.CheckBillingEligibility(context.Background(), user, nil, nil, nil, PlatformAnthropic)
		require.ErrorIs(t, err, ErrUserPlatformDailyQuotaExhausted)
placeholder)

	t.Run("日配额未满放行", func(t *testing.T) {
		cache := &billInvBillingCacheStub{
			balance:    5.0,
			quotaEntry: billInvQuotaEntryV1(5.0, 4.99),
	placeholder
		svc := billInvNewBillingCacheService(t, cache, nil)
		err := svc.CheckBillingEligibility(context.Background(), user, nil, nil, nil, PlatformAnthropic)
	placeholder
placeholder)

	t.Run("订阅模式豁免platform配额检查", func(t *testing.T) {
		group := &Group{
			ID:               7,
			SubscriptionType: SubscriptionTypeSubscription,
			DailyLimitUSD:    billInvF64Ptr(10),
	placeholder
		cache := &billInvBillingCacheStub{
			sub:        &SubscriptionCacheData{Status: SubscriptionStatusActive, ExpiresAt: time.Now().Add(24 * time.Hour)placeholder,
			quotaEntry: billInvQuotaEntryV1(5.0, 999), // platform 配额早已超限
	placeholder
		svc := billInvNewBillingCacheService(t, cache, nil)
		err := svc.CheckBillingEligibility(context.Background(), user, nil, group, &UserSubscription{ID: 42placeholder, PlatformAnthropic)
		require.NoError(t, err, "订阅模式下 user×platform 配额不应生效")
placeholder)
placeholder

// TestBillingInvariant_PreflightSimpleModeBypass 锁定 simple 运行模式跳过所有
// 计费检查（余额为 0 也放行）。
func TestBillingInvariant_PreflightSimpleModeBypass(t *testing.T) {
	cfg := &config.Config{RunMode: config.RunModeSimpleplaceholder
	svc := billInvNewBillingCacheService(t, &billInvBillingCacheStub{balance: 0placeholder, cfg)
	err := svc.CheckBillingEligibility(context.Background(), &User{ID: 601placeholder, nil, nil, nil, PlatformAnthropic)
placeholder
placeholder
