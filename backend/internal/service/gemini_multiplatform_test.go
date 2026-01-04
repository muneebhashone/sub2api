//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

// mockAccountRepoForGemini Gemini 测试用的 mock
type mockAccountRepoForGemini struct {
	accounts     []Account
	accountsByID map[int64]*Account
placeholder

func (m *mockAccountRepoForGemini) GetByID(ctx context.Context, id int64) (*Account, error) {
	if acc, ok := m.accountsByID[id]; ok {
		return acc, nil
placeholder
	return nil, errors.New("account not found")
placeholder

func (m *mockAccountRepoForGemini) GetByIDs(ctx context.Context, ids []int64) ([]*Account, error) {
	var result []*Account
	for _, id := range ids {
		if acc, ok := m.accountsByID[id]; ok {
			result = append(result, acc)
	placeholder
placeholder
	return result, nil
placeholder

func (m *mockAccountRepoForGemini) ExistsByID(ctx context.Context, id int64) (bool, error) {
	if m.accountsByID == nil {
		return false, nil
placeholder
	_, ok := m.accountsByID[id]
	return ok, nil
placeholder

func (m *mockAccountRepoForGemini) ListSchedulableByPlatform(ctx context.Context, platform string) ([]Account, error) {
	var result []Account
	for _, acc := range m.accounts {
		if acc.Platform == platform && acc.IsSchedulable() {
			result = append(result, acc)
	placeholder
placeholder
	return result, nil
placeholder

func (m *mockAccountRepoForGemini) ListSchedulableByGroupIDAndPlatform(ctx context.Context, groupID int64, platform string) ([]Account, error) {
	// 测试时不区分 groupID，直接按 platform 过滤
	return m.ListSchedulableByPlatform(ctx, platform)
placeholder

// Stub methods to implement AccountRepository interface
func (m *mockAccountRepoForGemini) Create(ctx context.Context, account *Account) error { return nil placeholder
func (m *mockAccountRepoForGemini) GetByCRSAccountID(ctx context.Context, crsAccountID string) (*Account, error) {
	return nil, nil
placeholder
func (m *mockAccountRepoForGemini) Update(ctx context.Context, account *Account) error { return nil placeholder
func (m *mockAccountRepoForGemini) Delete(ctx context.Context, id int64) error         { return nil placeholder
func (m *mockAccountRepoForGemini) List(ctx context.Context, params pagination.PaginationParams) ([]Account, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (m *mockAccountRepoForGemini) ListWithFilters(ctx context.Context, params pagination.PaginationParams, platform, accountType, status, search string) ([]Account, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (m *mockAccountRepoForGemini) ListByGroup(ctx context.Context, groupID int64) ([]Account, error) {
	return nil, nil
placeholder
func (m *mockAccountRepoForGemini) ListActive(ctx context.Context) ([]Account, error) {
	return nil, nil
placeholder
func (m *mockAccountRepoForGemini) ListByPlatform(ctx context.Context, platform string) ([]Account, error) {
	return nil, nil
placeholder
func (m *mockAccountRepoForGemini) UpdateLastUsed(ctx context.Context, id int64) error { return nil placeholder
func (m *mockAccountRepoForGemini) BatchUpdateLastUsed(ctx context.Context, updates map[int64]time.Time) error {
	return nil
placeholder
func (m *mockAccountRepoForGemini) SetError(ctx context.Context, id int64, errorMsg string) error {
	return nil
placeholder
func (m *mockAccountRepoForGemini) SetSchedulable(ctx context.Context, id int64, schedulable bool) error {
	return nil
placeholder
func (m *mockAccountRepoForGemini) BindGroups(ctx context.Context, accountID int64, groupIDs []int64) error {
	return nil
placeholder
func (m *mockAccountRepoForGemini) ListSchedulable(ctx context.Context) ([]Account, error) {
	return nil, nil
placeholder
func (m *mockAccountRepoForGemini) ListSchedulableByGroupID(ctx context.Context, groupID int64) ([]Account, error) {
	return nil, nil
placeholder
func (m *mockAccountRepoForGemini) ListSchedulableByPlatforms(ctx context.Context, platforms []string) ([]Account, error) {
	var result []Account
	platformSet := make(map[string]bool)
	for _, p := range platforms {
		platformSet[p] = true
placeholder
	for _, acc := range m.accounts {
		if platformSet[acc.Platform] && acc.IsSchedulable() {
			result = append(result, acc)
	placeholder
placeholder
	return result, nil
placeholder
func (m *mockAccountRepoForGemini) ListSchedulableByGroupIDAndPlatforms(ctx context.Context, groupID int64, platforms []string) ([]Account, error) {
	return m.ListSchedulableByPlatforms(ctx, platforms)
placeholder
func (m *mockAccountRepoForGemini) SetRateLimited(ctx context.Context, id int64, resetAt time.Time) error {
	return nil
placeholder
func (m *mockAccountRepoForGemini) SetOverloaded(ctx context.Context, id int64, until time.Time) error {
	return nil
placeholder
func (m *mockAccountRepoForGemini) SetTempUnschedulable(ctx context.Context, id int64, until time.Time, reason string) error {
	return nil
placeholder
func (m *mockAccountRepoForGemini) ClearTempUnschedulable(ctx context.Context, id int64) error {
	return nil
placeholder
func (m *mockAccountRepoForGemini) ClearRateLimit(ctx context.Context, id int64) error { return nil placeholder
func (m *mockAccountRepoForGemini) UpdateSessionWindow(ctx context.Context, id int64, start, end *time.Time, status string) error {
	return nil
placeholder
func (m *mockAccountRepoForGemini) UpdateExtra(ctx context.Context, id int64, updates map[string]any) error {
	return nil
placeholder
func (m *mockAccountRepoForGemini) BulkUpdate(ctx context.Context, ids []int64, updates AccountBulkUpdate) (int64, error) {
	return 0, nil
placeholder

// Verify interface implementation
var _ AccountRepository = (*mockAccountRepoForGemini)(nil)

// mockGroupRepoForGemini Gemini 测试用的 group repo mock
type mockGroupRepoForGemini struct {
	groups map[int64]*Group
placeholder

func (m *mockGroupRepoForGemini) GetByID(ctx context.Context, id int64) (*Group, error) {
	if g, ok := m.groups[id]; ok {
		return g, nil
placeholder
	return nil, errors.New("group not found")
placeholder

// Stub methods to implement GroupRepository interface
func (m *mockGroupRepoForGemini) Create(ctx context.Context, group *Group) error { return nil placeholder
func (m *mockGroupRepoForGemini) Update(ctx context.Context, group *Group) error { return nil placeholder
func (m *mockGroupRepoForGemini) Delete(ctx context.Context, id int64) error     { return nil placeholder
func (m *mockGroupRepoForGemini) DeleteCascade(ctx context.Context, id int64) ([]int64, error) {
	return nil, nil
placeholder
func (m *mockGroupRepoForGemini) List(ctx context.Context, params pagination.PaginationParams) ([]Group, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (m *mockGroupRepoForGemini) ListWithFilters(ctx context.Context, params pagination.PaginationParams, platform, status string, isExclusive *bool) ([]Group, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (m *mockGroupRepoForGemini) ListActive(ctx context.Context) ([]Group, error) { return nil, nil placeholder
func (m *mockGroupRepoForGemini) ListActiveByPlatform(ctx context.Context, platform string) ([]Group, error) {
	return nil, nil
placeholder
func (m *mockGroupRepoForGemini) ExistsByName(ctx context.Context, name string) (bool, error) {
	return false, nil
placeholder
func (m *mockGroupRepoForGemini) GetAccountCount(ctx context.Context, groupID int64) (int64, error) {
	return 0, nil
placeholder
func (m *mockGroupRepoForGemini) DeleteAccountGroupsByGroupID(ctx context.Context, groupID int64) (int64, error) {
	return 0, nil
placeholder

var _ GroupRepository = (*mockGroupRepoForGemini)(nil)

// mockGatewayCacheForGemini Gemini 测试用的 cache mock
type mockGatewayCacheForGemini struct {
	sessionBindings map[string]int64
placeholder

func (m *mockGatewayCacheForGemini) GetSessionAccountID(ctx context.Context, sessionHash string) (int64, error) {
	if id, ok := m.sessionBindings[sessionHash]; ok {
		return id, nil
placeholder
	return 0, errors.New("not found")
placeholder

func (m *mockGatewayCacheForGemini) SetSessionAccountID(ctx context.Context, sessionHash string, accountID int64, ttl time.Duration) error {
	if m.sessionBindings == nil {
		m.sessionBindings = make(map[string]int64)
placeholder
	m.sessionBindings[sessionHash] = accountID
	return nil
placeholder

func (m *mockGatewayCacheForGemini) RefreshSessionTTL(ctx context.Context, sessionHash string, ttl time.Duration) error {
	return nil
placeholder

// TestGeminiMessagesCompatService_SelectAccountForModelWithExclusions_GeminiPlatform 测试 Gemini 单平台选择
func TestGeminiMessagesCompatService_SelectAccountForModelWithExclusions_GeminiPlatform(t *testing.T) {
	ctx := context.Background()

	repo := &mockAccountRepoForGemini{
		accounts: []Account{
			{ID: 1, Platform: PlatformGemini, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
			{ID: 2, Platform: PlatformGemini, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
			{ID: 3, Platform: PlatformAntigravity, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder, // 应被隔离
	placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder
	for i := range repo.accounts {
		repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
placeholder

	cache := &mockGatewayCacheForGemini{placeholder
	groupRepo := &mockGroupRepoForGemini{groups: map[int64]*Group{placeholderplaceholder

	svc := &GeminiMessagesCompatService{
		accountRepo: repo,
		groupRepo:   groupRepo,
		cache:       cache,
placeholder

	// 无分组时使用 gemini 平台
	acc, err := svc.SelectAccountForModelWithExclusions(ctx, nil, "", "gemini-2.5-flash", nil)
placeholder
	require.NotNil(t, acc)
	require.Equal(t, int64(1), acc.ID, "应选择优先级最高的 gemini 账户")
	require.Equal(t, PlatformGemini, acc.Platform, "无分组时应只返回 gemini 平台账户")
placeholder

// TestGeminiMessagesCompatService_SelectAccountForModelWithExclusions_AntigravityGroup 测试 antigravity 分组
func TestGeminiMessagesCompatService_SelectAccountForModelWithExclusions_AntigravityGroup(t *testing.T) {
	ctx := context.Background()

	repo := &mockAccountRepoForGemini{
		accounts: []Account{
			{ID: 1, Platform: PlatformGemini, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,      // 应被隔离
			{ID: 2, Platform: PlatformAntigravity, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder, // 应被选择
	placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder
	for i := range repo.accounts {
		repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
placeholder

	cache := &mockGatewayCacheForGemini{placeholder
	groupRepo := &mockGroupRepoForGemini{
		groups: map[int64]*Group{
			1: {ID: 1, Platform: PlatformAntigravityplaceholder,
	placeholder,
placeholder

	svc := &GeminiMessagesCompatService{
		accountRepo: repo,
		groupRepo:   groupRepo,
		cache:       cache,
placeholder

	groupID := int64(1)
	acc, err := svc.SelectAccountForModelWithExclusions(ctx, &groupID, "", "gemini-2.5-flash", nil)
placeholder
	require.NotNil(t, acc)
	require.Equal(t, int64(2), acc.ID)
	require.Equal(t, PlatformAntigravity, acc.Platform, "antigravity 分组应只返回 antigravity 账户")
placeholder

// TestGeminiMessagesCompatService_SelectAccountForModelWithExclusions_OAuthPreferred 测试 OAuth 优先
func TestGeminiMessagesCompatService_SelectAccountForModelWithExclusions_OAuthPreferred(t *testing.T) {
	ctx := context.Background()

	repo := &mockAccountRepoForGemini{
		accounts: []Account{
			{ID: 1, Platform: PlatformGemini, Type: AccountTypeAPIKey, Priority: 1, Status: StatusActive, Schedulable: true, LastUsedAt: nilplaceholder,
			{ID: 2, Platform: PlatformGemini, Type: AccountTypeOAuth, Priority: 1, Status: StatusActive, Schedulable: true, LastUsedAt: nilplaceholder,
	placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder
	for i := range repo.accounts {
		repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
placeholder

	cache := &mockGatewayCacheForGemini{placeholder
	groupRepo := &mockGroupRepoForGemini{groups: map[int64]*Group{placeholderplaceholder

	svc := &GeminiMessagesCompatService{
		accountRepo: repo,
		groupRepo:   groupRepo,
		cache:       cache,
placeholder

	acc, err := svc.SelectAccountForModelWithExclusions(ctx, nil, "", "gemini-2.5-flash", nil)
placeholder
	require.NotNil(t, acc)
	require.Equal(t, int64(2), acc.ID, "同优先级且都未使用时，应优先选择 OAuth 账户")
	require.Equal(t, AccountTypeOAuth, acc.Type)
placeholder

// TestGeminiMessagesCompatService_SelectAccountForModelWithExclusions_NoAvailableAccounts 测试无可用账户
func TestGeminiMessagesCompatService_SelectAccountForModelWithExclusions_NoAvailableAccounts(t *testing.T) {
	ctx := context.Background()

	repo := &mockAccountRepoForGemini{
		accounts:     []Account{placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder

	cache := &mockGatewayCacheForGemini{placeholder
	groupRepo := &mockGroupRepoForGemini{groups: map[int64]*Group{placeholderplaceholder

	svc := &GeminiMessagesCompatService{
		accountRepo: repo,
		groupRepo:   groupRepo,
		cache:       cache,
placeholder

	acc, err := svc.SelectAccountForModelWithExclusions(ctx, nil, "", "gemini-2.5-flash", nil)
placeholder
	require.Nil(t, acc)
	require.Contains(t, err.Error(), "no available")
placeholder

// TestGeminiMessagesCompatService_SelectAccountForModelWithExclusions_StickySession 测试粘性会话
func TestGeminiMessagesCompatService_SelectAccountForModelWithExclusions_StickySession(t *testing.T) {
	ctx := context.Background()

	t.Run("粘性会话命中-同平台", func(t *testing.T) {
		repo := &mockAccountRepoForGemini{
			accounts: []Account{
				{ID: 1, Platform: PlatformGemini, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
				{ID: 2, Platform: PlatformGemini, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
		placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder
		for i := range repo.accounts {
			repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
	placeholder

		// 注意：缓存键使用 "gemini:" 前缀
		cache := &mockGatewayCacheForGemini{
			sessionBindings: map[string]int64{"gemini:session-123": 1placeholder,
	placeholder
		groupRepo := &mockGroupRepoForGemini{groups: map[int64]*Group{placeholderplaceholder

		svc := &GeminiMessagesCompatService{
			accountRepo: repo,
			groupRepo:   groupRepo,
			cache:       cache,
	placeholder

		acc, err := svc.SelectAccountForModelWithExclusions(ctx, nil, "session-123", "gemini-2.5-flash", nil)
	placeholder
		require.NotNil(t, acc)
		require.Equal(t, int64(1), acc.ID, "应返回粘性会话绑定的账户")
placeholder)

	t.Run("粘性会话平台不匹配-降级选择", func(t *testing.T) {
		repo := &mockAccountRepoForGemini{
			accounts: []Account{
				{ID: 1, Platform: PlatformAntigravity, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder, // 粘性会话绑定
				{ID: 2, Platform: PlatformGemini, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
		placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder
		for i := range repo.accounts {
			repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
	placeholder

		cache := &mockGatewayCacheForGemini{
			sessionBindings: map[string]int64{"gemini:session-123": 1placeholder, // 绑定 antigravity 账户
	placeholder
		groupRepo := &mockGroupRepoForGemini{groups: map[int64]*Group{placeholderplaceholder

		svc := &GeminiMessagesCompatService{
			accountRepo: repo,
			groupRepo:   groupRepo,
			cache:       cache,
	placeholder

		// 无分组时使用 gemini 平台，粘性会话绑定的 antigravity 账户平台不匹配
		acc, err := svc.SelectAccountForModelWithExclusions(ctx, nil, "session-123", "gemini-2.5-flash", nil)
	placeholder
		require.NotNil(t, acc)
		require.Equal(t, int64(2), acc.ID, "粘性会话账户平台不匹配，应降级选择 gemini 账户")
		require.Equal(t, PlatformGemini, acc.Platform)
placeholder)

	t.Run("粘性会话不命中无前缀缓存键", func(t *testing.T) {
		repo := &mockAccountRepoForGemini{
			accounts: []Account{
				{ID: 1, Platform: PlatformGemini, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
				{ID: 2, Platform: PlatformGemini, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
		placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder
		for i := range repo.accounts {
			repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
	placeholder

		// 缓存键没有 "gemini:" 前缀，不应命中
		cache := &mockGatewayCacheForGemini{
			sessionBindings: map[string]int64{"session-123": 1placeholder,
	placeholder
		groupRepo := &mockGroupRepoForGemini{groups: map[int64]*Group{placeholderplaceholder

		svc := &GeminiMessagesCompatService{
			accountRepo: repo,
			groupRepo:   groupRepo,
			cache:       cache,
	placeholder

		acc, err := svc.SelectAccountForModelWithExclusions(ctx, nil, "session-123", "gemini-2.5-flash", nil)
	placeholder
		require.NotNil(t, acc)
		// 粘性会话未命中，按优先级选择
		require.Equal(t, int64(2), acc.ID, "粘性会话未命中，应按优先级选择")
placeholder)
placeholder

// TestGeminiPlatformRouting_DocumentRouteDecision 测试平台路由决策逻辑
func TestGeminiPlatformRouting_DocumentRouteDecision(t *testing.T) {
	tests := []struct {
		name            string
		platform        string
		expectedService string // "gemini" 表示 ForwardNative, "antigravity" 表示 ForwardGemini
placeholder{
		{
			name:            "Gemini平台走ForwardNative",
			platform:        PlatformGemini,
			expectedService: "gemini",
	placeholder,
		{
			name:            "Antigravity平台走ForwardGemini",
			platform:        PlatformAntigravity,
			expectedService: "antigravity",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := &Account{Platform: tt.platformplaceholder

			// 模拟 Handler 层的路由逻辑
			var serviceName string
			if account.Platform == PlatformAntigravity {
				serviceName = "antigravity"
		placeholder else {
				serviceName = "gemini"
		placeholder

			require.Equal(t, tt.expectedService, serviceName,
				"平台 %s 应该路由到 %s 服务", tt.platform, tt.expectedService)
	placeholder)
placeholder
placeholder

func TestGeminiMessagesCompatService_isModelSupportedByAccount(t *testing.T) {
	svc := &GeminiMessagesCompatService{placeholder

	tests := []struct {
		name     string
		account  *Account
		model    string
		expected bool
placeholder{
		{
			name:     "Antigravity平台-支持gemini模型",
			account:  &Account{Platform: PlatformAntigravityplaceholder,
			model:    "gemini-2.5-flash",
			expected: true,
	placeholder,
		{
			name:     "Antigravity平台-支持claude模型",
			account:  &Account{Platform: PlatformAntigravityplaceholder,
			model:    "claude-3-5-sonnet-20241022",
			expected: true,
	placeholder,
		{
			name:     "Antigravity平台-不支持gpt模型",
			account:  &Account{Platform: PlatformAntigravityplaceholder,
			model:    "gpt-4",
			expected: false,
	placeholder,
		{
			name:     "Gemini平台-无映射配置-支持所有模型",
			account:  &Account{Platform: PlatformGeminiplaceholder,
			model:    "gemini-2.5-flash",
			expected: true,
	placeholder,
		{
			name: "Gemini平台-有映射配置-只支持配置的模型",
			account: &Account{
				Platform:    PlatformGemini,
		placeholder"model_mapping": map[string]any{"gemini-1.5-pro": "x"placeholderplaceholder,
		placeholder,
			model:    "gemini-2.5-flash",
			expected: false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.isModelSupportedByAccount(tt.account, tt.model)
			require.Equal(t, tt.expected, got)
	placeholder)
placeholder
placeholder
