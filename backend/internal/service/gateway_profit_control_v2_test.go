//go:build unit

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	"github.com/stretchr/testify/require"
)

func gatewayProfitTestGroup(id int64, platform string) *Group {
	return &Group{
		ID:                   id,
		Name:                 "profit-" + platform,
		Platform:             platform,
		Status:               StatusActive,
		Hydrated:             true,
		RateMultiplier:       0.5,
		SubscriptionType:     SubscriptionTypeStandard,
		ProfitControlEnabled: true,
		ProfitMinMargin:      0,
		ProfitSafetyBuffer:   0,
placeholder
placeholder

func gatewayProfitTestContext(group *Group) context.Context {
	ctx := context.WithValue(context.Background(), ctxkey.Group, group)
	ctx, _ = WithGatewayTokenRequestPricing(ctx)
	return ctx
placeholder

func gatewayProfitTestAccount(id int64, platform string, rate float64, groupID int64) Account {
	return Account{
		ID:             id,
		Name:           "account",
		Platform:       platform,
		Type:           AccountTypeAPIKey,
		Status:         StatusActive,
		Schedulable:    true,
		Concurrency:    2,
		Priority:       1,
		RateMultiplier: &rate,
		AccountGroups:  []AccountGroup{{AccountID: id, GroupID: groupIDplaceholderplaceholder,
		GroupIDs:       []int64{groupIDplaceholder,
placeholder
placeholder

func TestGatewayProfitControlInstallsForFivePlatformsOnlyOnTokenRequests(t *testing.T) {
	for _, platform := range []string{
		PlatformOpenAI,
		PlatformAnthropic,
		PlatformGemini,
		PlatformGrok,
		PlatformAntigravity,
placeholder {
		t.Run(platform, func(t *testing.T) {
			group := gatewayProfitTestGroup(101, platform)
			groupID := group.ID
			svc := &GatewayService{placeholder

			tokenCtx := svc.withGatewayProfitControlGate(gatewayProfitTestContext(group), &groupID)
			gate, _ := tokenCtx.Value(openAIProfitControlGateCtxKey{placeholder).(*openAIProfitControlGate)
			require.NotNil(t, gate)
			require.Equal(t, platform, gate.platform)
			require.InDelta(t, 0.5, gate.threshold, 1e-12)

			metadataCtx := context.WithValue(context.Background(), ctxkey.Group, group)
			metadataCtx = svc.withGatewayProfitControlGate(metadataCtx, &groupID)
			gate, _ = metadataCtx.Value(openAIProfitControlGateCtxKey{placeholder).(*openAIProfitControlGate)
			require.Nil(t, gate, "未显式标记为 token 请求的入口不得装门")
	placeholder)
placeholder
placeholder

func TestGatewayProfitControlCompositeBillingUsesScheduledMemberConfig(t *testing.T) {
	billingGroup := &Group{
		ID:               201,
		Platform:         PlatformComposite,
		Status:           StatusActive,
		Hydrated:         true,
		RateMultiplier:   0.4,
		SubscriptionType: SubscriptionTypeStandard,
placeholder
	memberGroup := gatewayProfitTestGroup(202, PlatformAnthropic)
	memberGroup.RateMultiplier = 99
	memberGroup.ProfitMinMargin = 0.25

	ctx := context.WithValue(context.Background(), ctxkey.Group, billingGroup)
	ctx, pricingAt := WithGatewayTokenRequestPricing(ctx)
	svc := &GatewayService{
		schedulerSnapshot: NewSchedulerSnapshotService(
			nil,
			nil,
			nil,
			profitControlGroupRepo{group: memberGroupplaceholder,
			nil,
		),
placeholder
	ctx = svc.withGatewayProfitControlGate(ctx, &memberGroup.ID)
	gate, _ := ctx.Value(openAIProfitControlGateCtxKey{placeholder).(*openAIProfitControlGate)
	require.NotNil(t, gate)
	require.Equal(t, memberGroup.ID, gate.groupID)
	require.Equal(t, PlatformAnthropic, gate.platform)
	require.Equal(t, pricingAt, gate.pricingAt)
	require.InDelta(t, 0.4*(1-0.25), gate.threshold, 1e-12, "D 必须取 composite 计费父分组，margin 取被调度成员分组")
placeholder

func TestGatewayProfitControlGroupLoadFailureClearsForeignGate(t *testing.T) {
	billingGroup := &Group{
		ID:               211,
		Platform:         PlatformComposite,
		Status:           StatusActive,
		Hydrated:         true,
		RateMultiplier:   0.4,
		SubscriptionType: SubscriptionTypeStandard,
placeholder
	targetGroupID := int64(212)
	ctx := gatewayProfitTestContext(billingGroup)
	ctx = context.WithValue(ctx, openAIProfitControlGateCtxKey{placeholder, &openAIProfitControlGate{
		groupID:   210,
		platform:  PlatformAnthropic,
		threshold: 0.1,
placeholder)
	svc := &GatewayService{
		schedulerSnapshot: NewSchedulerSnapshotService(
			nil,
			nil,
			nil,
			profitControlFailingGroupRepo{placeholder,
			nil,
		),
placeholder

	ctx = svc.withGatewayProfitControlGate(ctx, &targetGroupID)
	gate, ok := ctx.Value(openAIProfitControlGateCtxKey{placeholder).(*openAIProfitControlGate)
	require.True(t, ok)
	require.Nil(t, gate, "加载新分组失败时必须清除其他分组遗留的门")

	account := gatewayProfitTestAccount(213, PlatformAnthropic, 0.8, targetGroupID)
	require.True(t, svc.isGatewayAccountProfitEligible(ctx, &account), "配置读取失败按既定语义 fail-open")
placeholder

type profitControlFailingGroupRepo struct {
	GroupRepository
placeholder

func (profitControlFailingGroupRepo) GetByID(context.Context, int64) (*Group, error) {
	return nil, errors.New("group cache unavailable")
placeholder

func TestGatewayProfitControlLegacyMixedAndRoutedSelection(t *testing.T) {
	t.Run("legacy single-platform selection", func(t *testing.T) {
		group := gatewayProfitTestGroup(111, PlatformGrok)
		cheap := gatewayProfitTestAccount(1, PlatformGrok, 0.2, group.ID)
		expensive := gatewayProfitTestAccount(2, PlatformGrok, 0.8, group.ID)
		repo := &mockAccountRepoForPlatform{
			accounts:     []Account{expensive, cheapplaceholder,
			accountsByID: map[int64]*Account{cheap.ID: &cheap, expensive.ID: &expensiveplaceholder,
	placeholder
		svc := &GatewayService{
			accountRepo: repo,
			cache:       &mockGatewayCacheForPlatform{placeholder,
			cfg:         testConfig(),
	placeholder

		selected, err := svc.SelectAccountForModelWithExclusions(
			gatewayProfitTestContext(group), &group.ID, "", "", nil,
		)
	placeholder
		require.Equal(t, cheap.ID, selected.ID)

		_, err = svc.SelectAccountForModelWithExclusions(
			gatewayProfitTestContext(group), &group.ID, "", "", map[int64]struct{placeholder{cheap.ID: {placeholderplaceholder,
		)
	placeholder
		require.ErrorIs(t, err, ErrNoAvailableAccounts)
placeholder)

	t.Run("mixed routing filters the routed account", func(t *testing.T) {
		group := gatewayProfitTestGroup(112, PlatformAnthropic)
		group.ModelRoutingEnabled = true
		group.ModelRouting = map[string][]int64{"claude-test": {2, 1placeholderplaceholder
		cheap := gatewayProfitTestAccount(1, PlatformAntigravity, 0.2, group.ID)
		cheap.Extra = map[string]any{"mixed_scheduling": trueplaceholder
		cheap.Credentials = map[string]any{"model_mapping": map[string]any{"claude-test": "claude-test"placeholderplaceholder
		expensive := gatewayProfitTestAccount(2, PlatformAnthropic, 0.8, group.ID)
		repo := &mockAccountRepoForPlatform{
			accounts:     []Account{expensive, cheapplaceholder,
			accountsByID: map[int64]*Account{cheap.ID: &cheap, expensive.ID: &expensiveplaceholder,
	placeholder
		svc := &GatewayService{
			accountRepo: repo,
			cache:       &mockGatewayCacheForPlatform{placeholder,
			cfg:         testConfig(),
	placeholder

		selected, err := svc.SelectAccountForModelWithExclusions(
			gatewayProfitTestContext(group), &group.ID, "", "claude-test", nil,
		)
	placeholder
		require.Equal(t, cheap.ID, selected.ID)
placeholder)
placeholder

func TestGatewayProfitControlLoadAwareSelectionAndFailover(t *testing.T) {
	group := gatewayProfitTestGroup(121, PlatformGrok)
	cheap := gatewayProfitTestAccount(1, PlatformGrok, 0.2, group.ID)
	expensive := gatewayProfitTestAccount(2, PlatformGrok, 0.8, group.ID)
	repo := &mockAccountRepoForPlatform{
		accounts:     []Account{expensive, cheapplaceholder,
		accountsByID: map[int64]*Account{cheap.ID: &cheap, expensive.ID: &expensiveplaceholder,
placeholder
	cfg := &config.Config{RunMode: config.RunModeStandardplaceholder
	cfg.Gateway.Scheduling.LoadBatchEnabled = true
	svc := &GatewayService{
		accountRepo:        repo,
		cache:              &mockGatewayCacheForPlatform{placeholder,
		cfg:                cfg,
		concurrencyService: NewConcurrencyService(stubConcurrencyCache{placeholder),
placeholder

	result, err := svc.SelectAccountWithLoadAwareness(
		gatewayProfitTestContext(group), &group.ID, "", "", nil, "", 0,
	)
placeholder
	require.NotNil(t, result)
	require.Equal(t, cheap.ID, result.Account.ID)
	if result.ReleaseFunc != nil {
		result.ReleaseFunc()
placeholder

	result, err = svc.SelectAccountWithLoadAwareness(
		gatewayProfitTestContext(group),
		&group.ID,
		"",
		"",
		map[int64]struct{placeholder{cheap.ID: {placeholderplaceholder,
		"",
		0,
	)
	require.Nil(t, result)
placeholder
	require.ErrorIs(t, err, ErrNoAvailableAccounts)
placeholder

func TestGatewayProfitControlStickyVetoKeepsBindingUntilRateRecovers(t *testing.T) {
	group := gatewayProfitTestGroup(131, PlatformAnthropic)
	expensive := gatewayProfitTestAccount(1, PlatformAnthropic, 0.8, group.ID)
	cheap := gatewayProfitTestAccount(2, PlatformAnthropic, 0.2, group.ID)
	repo := &mockAccountRepoForPlatform{
		accounts:     []Account{expensive, cheapplaceholder,
		accountsByID: map[int64]*Account{expensive.ID: &expensive, cheap.ID: &cheapplaceholder,
placeholder
	cache := &mockGatewayCacheForPlatform{
		sessionBindings: map[string]int64{"sticky-profit": expensive.IDplaceholder,
placeholder
	svc := &GatewayService{
		accountRepo: repo,
		cache:       cache,
		cfg:         testConfig(),
placeholder
	ctx := gatewayProfitTestContext(group)

	selected, err := svc.SelectAccountForModelWithExclusions(ctx, &group.ID, "sticky-profit", "", nil)
placeholder
	require.Equal(t, cheap.ID, selected.ID)
	require.Equal(t, expensive.ID, cache.sessionBindings["sticky-profit"], "候选过滤不得覆盖旧粘性绑定")

	require.NoError(t, svc.BindStickySessionAfterProfitAdmission(
		svc.withGatewayProfitControlGate(ctx, &group.ID),
		&group.ID,
		"sticky-profit",
		cheap.ID,
		expensive.ID,
	))
	require.Equal(t, expensive.ID, cache.sessionBindings["sticky-profit"], "终检通过的 fallback 账号也不得覆盖旧绑定")
	require.Zero(t, cache.deletedSessions["sticky-profit"])

	recovered := expensive
	recoveredRate := 0.2
	recovered.RateMultiplier = &recoveredRate
	repo.accounts[0] = recovered
	repo.accountsByID[recovered.ID] = &repo.accounts[0]

	selected, err = svc.SelectAccountForModelWithExclusions(ctx, &group.ID, "sticky-profit", "", nil)
placeholder
	require.Equal(t, recovered.ID, selected.ID, "倍率恢复后应重新命中原粘性账号")
	require.Zero(t, cache.deletedSessions["sticky-profit"])
placeholder

type gatewayProfitSnapshotCache struct {
	SchedulerCache
	account *Account
	err     error
placeholder

func (c *gatewayProfitSnapshotCache) GetAccount(context.Context, int64) (*Account, error) {
	return c.account, c.err
placeholder

type gatewayProfitAccountRepo struct {
	AccountRepository
	account *Account
	err     error
placeholder

func (r gatewayProfitAccountRepo) GetByID(context.Context, int64) (*Account, error) {
	return r.account, r.err
placeholder

func TestGatewayProfitControlTerminalRefreshUsesReplacementObject(t *testing.T) {
	selected := gatewayProfitTestAccount(141, PlatformGemini, 0.2, 1)
	replacement := selected
	expensiveRate := 0.8
	replacement.RateMultiplier = &expensiveRate

	snapshot := NewSchedulerSnapshotService(
		&gatewayProfitSnapshotCache{account: &replacementplaceholder,
		nil,
		gatewayProfitAccountRepo{placeholder,
		nil,
		nil,
	)
	ctx := context.WithValue(context.Background(), openAIProfitControlGateCtxKey{placeholder, &openAIProfitControlGate{
		groupID:   1,
		platform:  PlatformGemini,
		threshold: 0.5,
placeholder)

	latest, vetoed, reason := profitControlVetoLatest(ctx, &selected, snapshot)
	require.Same(t, &replacement, latest)
	require.True(t, vetoed)
	require.Equal(t, openAIProfitFilterReasonThreshold, reason)
	require.InDelta(t, 0.2, *selected.RateMultiplier, 1e-12, "测试必须替换缓存对象，不能原地修改旧指针")
placeholder

func TestGatewayProfitControlTerminalRefreshFallsBackFromCacheToDatabase(t *testing.T) {
	selected := gatewayProfitTestAccount(145, PlatformAnthropic, 0.2, 1)
	replacement := selected
	expensiveRate := 0.8
	replacement.RateMultiplier = &expensiveRate

	snapshot := NewSchedulerSnapshotService(
		&gatewayProfitSnapshotCache{err: errors.New("cache unavailable")placeholder,
		nil,
		gatewayProfitAccountRepo{account: &replacementplaceholder,
		nil,
		nil,
	)
	ctx := context.WithValue(context.Background(), openAIProfitControlGateCtxKey{placeholder, &openAIProfitControlGate{
		groupID:   1,
		platform:  PlatformAnthropic,
		threshold: 0.5,
placeholder)

	latest, vetoed, reason := profitControlVetoLatest(ctx, &selected, snapshot)
	require.Same(t, &replacement, latest)
	require.True(t, vetoed, "缓存读取失败时必须继续从数据库重读，不能直接使用选号旧对象")
	require.Equal(t, openAIProfitFilterReasonThreshold, reason)
placeholder

func TestGatewayProfitControlTerminalRefreshFailureFallsBackToSelectedObject(t *testing.T) {
	selected := gatewayProfitTestAccount(151, PlatformAntigravity, 0.2, 1)
	snapshot := NewSchedulerSnapshotService(
		&gatewayProfitSnapshotCache{err: errors.New("cache unavailable")placeholder,
		nil,
		gatewayProfitAccountRepo{err: errors.New("database unavailable")placeholder,
		nil,
		nil,
	)
	ctx := context.WithValue(context.Background(), openAIProfitControlGateCtxKey{placeholder, &openAIProfitControlGate{
		groupID:   1,
		platform:  PlatformAntigravity,
		threshold: 0.5,
placeholder)

	latest, vetoed, reason := profitControlVetoLatest(ctx, &selected, snapshot)
	require.Same(t, &selected, latest)
	require.False(t, vetoed)
	require.Empty(t, reason)
placeholder
