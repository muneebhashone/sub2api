//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Part 1: isAccountInGroup 单元测试
// ============================================================================

func TestIsAccountInGroup(t *testing.T) {
	svc := &GatewayService{placeholder
	groupID100 := int64(100)
	groupID200 := int64(200)

	tests := []struct {
		name     string
		account  *Account
		groupID  *int64
		expected bool
placeholder{
		// groupID == nil（无分组 API Key）
		{
			"nil_groupID_ungrouped_account_nil_groups",
			&Account{ID: 1, AccountGroups: nilplaceholder,
			nil, true,
	placeholder,
		{
			"nil_groupID_ungrouped_account_empty_slice",
			&Account{ID: 2, AccountGroups: []AccountGroup{placeholderplaceholder,
			nil, true,
	placeholder,
		{
			"nil_groupID_grouped_account_single",
			&Account{ID: 3, AccountGroups: []AccountGroup{{GroupID: 100placeholderplaceholderplaceholder,
			nil, false,
	placeholder,
		{
			"nil_groupID_grouped_account_multiple",
			&Account{ID: 4, AccountGroups: []AccountGroup{{GroupID: 100placeholder, {GroupID: 200placeholderplaceholderplaceholder,
			nil, false,
	placeholder,
		// groupID != nil（有分组 API Key）
		{
			"with_groupID_account_in_group",
			&Account{ID: 5, AccountGroups: []AccountGroup{{GroupID: 100placeholderplaceholderplaceholder,
			&groupID100, true,
	placeholder,
		{
			"with_groupID_account_not_in_group",
			&Account{ID: 6, AccountGroups: []AccountGroup{{GroupID: 200placeholderplaceholderplaceholder,
			&groupID100, false,
	placeholder,
		{
			"with_groupID_ungrouped_account",
			&Account{ID: 7, AccountGroups: nilplaceholder,
			&groupID100, false,
	placeholder,
		{
			"with_groupID_multi_group_account_match_one",
			&Account{ID: 8, AccountGroups: []AccountGroup{{GroupID: 100placeholder, {GroupID: 200placeholderplaceholderplaceholder,
			&groupID200, true,
	placeholder,
		{
			"with_groupID_multi_group_account_no_match",
			&Account{ID: 9, AccountGroups: []AccountGroup{{GroupID: 300placeholder, {GroupID: 400placeholderplaceholderplaceholder,
			&groupID100, false,
	placeholder,
		// 防御性边界
		{
			"nil_account_nil_groupID",
			nil,
			nil, false,
	placeholder,
		{
			"nil_account_with_groupID",
			nil,
			&groupID100, false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.isAccountInGroup(tt.account, tt.groupID)
			require.Equal(t, tt.expected, got, "isAccountInGroup 结果不符预期")
	placeholder)
placeholder
placeholder

// ============================================================================
// Part 2: 分组隔离端到端调度测试
// ============================================================================

// groupAwareMockAccountRepo 嵌入 mockAccountRepoForPlatform，覆写分组隔离相关方法。
// allAccounts 存储所有账号，分组查询方法按 AccountGroups 字段进行真实过滤。
type groupAwareMockAccountRepo struct {
	*mockAccountRepoForPlatform
	allAccounts []Account
placeholder

// ListSchedulableUngroupedByPlatform 仅返回未分组账号（AccountGroups 为空）
func (m *groupAwareMockAccountRepo) ListSchedulableUngroupedByPlatform(ctx context.Context, platform string) ([]Account, error) {
	var result []Account
	for _, acc := range m.allAccounts {
		if acc.Platform == platform && acc.IsSchedulable() && len(acc.AccountGroups) == 0 {
			result = append(result, acc)
	placeholder
placeholder
	return result, nil
placeholder

// ListSchedulableUngroupedByPlatforms 仅返回未分组账号（多平台版本）
func (m *groupAwareMockAccountRepo) ListSchedulableUngroupedByPlatforms(ctx context.Context, platforms []string) ([]Account, error) {
	platformSet := make(map[string]bool, len(platforms))
	for _, p := range platforms {
		platformSet[p] = true
placeholder
	var result []Account
	for _, acc := range m.allAccounts {
		if platformSet[acc.Platform] && acc.IsSchedulable() && len(acc.AccountGroups) == 0 {
			result = append(result, acc)
	placeholder
placeholder
	return result, nil
placeholder

// ListSchedulableByGroupIDAndPlatform 返回属于指定分组的账号
func (m *groupAwareMockAccountRepo) ListSchedulableByGroupIDAndPlatform(ctx context.Context, groupID int64, platform string) ([]Account, error) {
	var result []Account
	for _, acc := range m.allAccounts {
		if acc.Platform == platform && acc.IsSchedulable() && accountBelongsToGroup(acc, groupID) {
			result = append(result, acc)
	placeholder
placeholder
	return result, nil
placeholder

// ListSchedulableByGroupIDAndPlatforms 返回属于指定分组的账号（多平台版本）
func (m *groupAwareMockAccountRepo) ListSchedulableByGroupIDAndPlatforms(ctx context.Context, groupID int64, platforms []string) ([]Account, error) {
	platformSet := make(map[string]bool, len(platforms))
	for _, p := range platforms {
		platformSet[p] = true
placeholder
	var result []Account
	for _, acc := range m.allAccounts {
		if platformSet[acc.Platform] && acc.IsSchedulable() && accountBelongsToGroup(acc, groupID) {
			result = append(result, acc)
	placeholder
placeholder
	return result, nil
placeholder

// accountBelongsToGroup 检查账号是否属于指定分组
func accountBelongsToGroup(acc Account, groupID int64) bool {
	for _, ag := range acc.AccountGroups {
		if ag.GroupID == groupID {
			return true
	placeholder
placeholder
	return false
placeholder

// Verify interface implementation
var _ AccountRepository = (*groupAwareMockAccountRepo)(nil)

// newGroupAwareMockRepo 创建分组感知的 mock repo
func newGroupAwareMockRepo(accounts []Account) *groupAwareMockAccountRepo {
	byID := make(map[int64]*Account, len(accounts))
	for i := range accounts {
		byID[accounts[i].ID] = &accounts[i]
placeholder
	return &groupAwareMockAccountRepo{
		mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
			accounts:     accounts,
			accountsByID: byID,
	placeholder,
		allAccounts: accounts,
placeholder
placeholder

func TestGroupIsolation_UngroupedKey_ShouldNotScheduleGroupedAccounts(t *testing.T) {
	// 场景：无分组 API Key（groupID=nil），池中只有已分组账号 → 应返回错误
	ctx := context.Background()

	accounts := []Account{
		{ID: 1, Platform: PlatformOpenAI, Priority: 1, Status: StatusActive, Schedulable: true,
			AccountGroups: []AccountGroup{{GroupID: 100placeholderplaceholderplaceholder,
		{ID: 2, Platform: PlatformOpenAI, Priority: 2, Status: StatusActive, Schedulable: true,
			AccountGroups: []AccountGroup{{GroupID: 200placeholderplaceholderplaceholder,
placeholder
	repo := newGroupAwareMockRepo(accounts)
	cache := &mockGatewayCacheForPlatform{placeholder

	svc := &GatewayService{
		accountRepo: repo,
		cache:       cache,
		cfg:         testConfig(),
placeholder

	acc, err := svc.selectAccountForModelWithPlatform(ctx, nil, "", "", nil, PlatformOpenAI)
	require.Error(t, err, "无分组 Key 不应调度到已分组账号")
	require.Nil(t, acc)
placeholder

func TestGroupIsolation_GroupedKey_ShouldNotScheduleUngroupedAccounts(t *testing.T) {
	// 场景：有分组 API Key（groupID=100），池中只有未分组账号 → 应返回错误
	ctx := context.Background()
	groupID := int64(100)

	accounts := []Account{
		{ID: 1, Platform: PlatformOpenAI, Priority: 1, Status: StatusActive, Schedulable: true,
			AccountGroups: nilplaceholder,
		{ID: 2, Platform: PlatformOpenAI, Priority: 2, Status: StatusActive, Schedulable: true,
			AccountGroups: []AccountGroup{placeholderplaceholder,
placeholder
	repo := newGroupAwareMockRepo(accounts)
	cache := &mockGatewayCacheForPlatform{placeholder

	svc := &GatewayService{
		accountRepo: repo,
		cache:       cache,
		cfg:         testConfig(),
placeholder

	acc, err := svc.selectAccountForModelWithPlatform(ctx, &groupID, "", "", nil, PlatformOpenAI)
	require.Error(t, err, "有分组 Key 不应调度到未分组账号")
	require.Nil(t, acc)
placeholder

func TestGroupIsolation_UngroupedKey_ShouldOnlyScheduleUngroupedAccounts(t *testing.T) {
	// 场景：无分组 API Key（groupID=nil），池中有未分组和已分组账号 → 应只选中未分组的
	ctx := context.Background()

	accounts := []Account{
		{ID: 1, Platform: PlatformOpenAI, Priority: 1, Status: StatusActive, Schedulable: true,
			AccountGroups: []AccountGroup{{GroupID: 100placeholderplaceholderplaceholder, // 已分组，不应被选中
		{ID: 2, Platform: PlatformOpenAI, Priority: 2, Status: StatusActive, Schedulable: true,
			AccountGroups: nilplaceholder, // 未分组，应被选中
		{ID: 3, Platform: PlatformOpenAI, Priority: 3, Status: StatusActive, Schedulable: true,
			AccountGroups: []AccountGroup{{GroupID: 200placeholderplaceholderplaceholder, // 已分组，不应被选中
placeholder
	repo := newGroupAwareMockRepo(accounts)
	cache := &mockGatewayCacheForPlatform{placeholder

	svc := &GatewayService{
		accountRepo: repo,
		cache:       cache,
		cfg:         testConfig(),
placeholder

	acc, err := svc.selectAccountForModelWithPlatform(ctx, nil, "", "", nil, PlatformOpenAI)
	require.NoError(t, err, "应成功调度未分组账号")
	require.NotNil(t, acc)
	require.Equal(t, int64(2), acc.ID, "应选中未分组的账号 ID=2")
placeholder

func TestGroupIsolation_GroupedKey_ShouldOnlyScheduleMatchingGroupAccounts(t *testing.T) {
	// 场景：有分组 API Key（groupID=100），池中有未分组和多个分组账号 → 应只选中分组 100 内的
	ctx := context.Background()
	groupID := int64(100)

	accounts := []Account{
		{ID: 1, Platform: PlatformOpenAI, Priority: 1, Status: StatusActive, Schedulable: true,
			AccountGroups: nilplaceholder, // 未分组，不应被选中
		{ID: 2, Platform: PlatformOpenAI, Priority: 2, Status: StatusActive, Schedulable: true,
			AccountGroups: []AccountGroup{{GroupID: 200placeholderplaceholderplaceholder, // 属于分组 200，不应被选中
		{ID: 3, Platform: PlatformOpenAI, Priority: 3, Status: StatusActive, Schedulable: true,
			AccountGroups: []AccountGroup{{GroupID: 100placeholderplaceholderplaceholder, // 属于分组 100，应被选中
placeholder
	repo := newGroupAwareMockRepo(accounts)
	cache := &mockGatewayCacheForPlatform{placeholder

	svc := &GatewayService{
		accountRepo: repo,
		cache:       cache,
		cfg:         testConfig(),
placeholder

	acc, err := svc.selectAccountForModelWithPlatform(ctx, &groupID, "", "", nil, PlatformOpenAI)
	require.NoError(t, err, "应成功调度分组内账号")
	require.NotNil(t, acc)
	require.Equal(t, int64(3), acc.ID, "应选中分组 100 内的账号 ID=3")
placeholder

// ============================================================================
// Part 3: SimpleMode 旁路测试
// ============================================================================

func TestGroupIsolation_SimpleMode_SkipsGroupIsolation(t *testing.T) {
	// SimpleMode 应跳过分组隔离，使用 ListSchedulableByPlatform 返回所有账号。
	// 测试非 useMixed 路径（platform=openai，不会触发 mixed 调度逻辑）。
	ctx := context.Background()

	// 混合未分组和已分组账号，SimpleMode 下应全部可调度
	accounts := []Account{
		{ID: 1, Platform: PlatformOpenAI, Priority: 2, Status: StatusActive, Schedulable: true,
			AccountGroups: []AccountGroup{{GroupID: 100placeholderplaceholderplaceholder, // 已分组
		{ID: 2, Platform: PlatformOpenAI, Priority: 1, Status: StatusActive, Schedulable: true,
			AccountGroups: nilplaceholder, // 未分组
placeholder

	// 使用基础 mock（ListSchedulableByPlatform 返回所有匹配平台的账号，不做分组过滤）
	byID := make(map[int64]*Account, len(accounts))
	for i := range accounts {
		byID[accounts[i].ID] = &accounts[i]
placeholder
	repo := &mockAccountRepoForPlatform{
		accounts:     accounts,
		accountsByID: byID,
placeholder
	cache := &mockGatewayCacheForPlatform{placeholder

	svc := &GatewayService{
		accountRepo: repo,
		cache:       cache,
		cfg:         &config.Config{RunMode: config.RunModeSimpleplaceholder,
placeholder

	// groupID=nil 时，SimpleMode 应使用 ListSchedulableByPlatform（不过滤分组）
	acc, err := svc.selectAccountForModelWithPlatform(ctx, nil, "", "", nil, PlatformOpenAI)
	require.NoError(t, err, "SimpleMode 应跳过分组隔离直接返回账号")
	require.NotNil(t, acc)
	// 应选择优先级最高的账号（Priority=1, ID=2），即使它未分组
	require.Equal(t, int64(2), acc.ID, "SimpleMode 应按优先级选择，不考虑分组")
placeholder

func TestGroupIsolation_SimpleMode_GroupedAccountAlsoSchedulable(t *testing.T) {
	// SimpleMode + groupID=nil 时，已分组账号也应该可被调度
	ctx := context.Background()

	// 只有已分组账号，在 standard 模式下 groupID=nil 会报错，但 simple 模式应正常
	accounts := []Account{
		{ID: 1, Platform: PlatformOpenAI, Priority: 1, Status: StatusActive, Schedulable: true,
			AccountGroups: []AccountGroup{{GroupID: 100placeholderplaceholderplaceholder,
placeholder

	byID := make(map[int64]*Account, len(accounts))
	for i := range accounts {
		byID[accounts[i].ID] = &accounts[i]
placeholder
	repo := &mockAccountRepoForPlatform{
		accounts:     accounts,
		accountsByID: byID,
placeholder
	cache := &mockGatewayCacheForPlatform{placeholder

	svc := &GatewayService{
		accountRepo: repo,
		cache:       cache,
		cfg:         &config.Config{RunMode: config.RunModeSimpleplaceholder,
placeholder

	acc, err := svc.selectAccountForModelWithPlatform(ctx, nil, "", "", nil, PlatformOpenAI)
	require.NoError(t, err, "SimpleMode 下已分组账号也应可调度")
	require.NotNil(t, acc)
	require.Equal(t, int64(1), acc.ID, "SimpleMode 应能调度已分组账号")
placeholder
