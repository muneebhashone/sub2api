//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

// testConfig 返回一个用于测试的默认配置
func testConfig() *config.Config {
	return &config.Config{RunMode: config.RunModeStandardplaceholder
placeholder

// mockAccountRepoForPlatform 单平台测试用的 mock
type mockAccountRepoForPlatform struct {
	accounts         []Account
	accountsByID     map[int64]*Account
	listPlatformFunc func(ctx context.Context, platform string) ([]Account, error)
placeholder

func (m *mockAccountRepoForPlatform) GetByID(ctx context.Context, id int64) (*Account, error) {
	if acc, ok := m.accountsByID[id]; ok {
		return acc, nil
placeholder
	return nil, errors.New("account not found")
placeholder

func (m *mockAccountRepoForPlatform) GetByIDs(ctx context.Context, ids []int64) ([]*Account, error) {
	var result []*Account
	for _, id := range ids {
		if acc, ok := m.accountsByID[id]; ok {
			result = append(result, acc)
	placeholder
placeholder
	return result, nil
placeholder

func (m *mockAccountRepoForPlatform) ExistsByID(ctx context.Context, id int64) (bool, error) {
	if m.accountsByID == nil {
		return false, nil
placeholder
	_, ok := m.accountsByID[id]
	return ok, nil
placeholder

func (m *mockAccountRepoForPlatform) ListSchedulableByPlatform(ctx context.Context, platform string) ([]Account, error) {
	if m.listPlatformFunc != nil {
		return m.listPlatformFunc(ctx, platform)
placeholder
	var result []Account
	for _, acc := range m.accounts {
		if acc.Platform == platform && acc.IsSchedulable() {
			result = append(result, acc)
	placeholder
placeholder
	return result, nil
placeholder

func (m *mockAccountRepoForPlatform) ListSchedulableByGroupIDAndPlatform(ctx context.Context, groupID int64, platform string) ([]Account, error) {
	return m.ListSchedulableByPlatform(ctx, platform)
placeholder

// Stub methods to implement AccountRepository interface
func (m *mockAccountRepoForPlatform) Create(ctx context.Context, account *Account) error {
	return nil
placeholder
func (m *mockAccountRepoForPlatform) GetByCRSAccountID(ctx context.Context, crsAccountID string) (*Account, error) {
	return nil, nil
placeholder
func (m *mockAccountRepoForPlatform) Update(ctx context.Context, account *Account) error {
	return nil
placeholder
func (m *mockAccountRepoForPlatform) Delete(ctx context.Context, id int64) error { return nil placeholder
func (m *mockAccountRepoForPlatform) List(ctx context.Context, params pagination.PaginationParams) ([]Account, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (m *mockAccountRepoForPlatform) ListWithFilters(ctx context.Context, params pagination.PaginationParams, platform, accountType, status, search string) ([]Account, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (m *mockAccountRepoForPlatform) ListByGroup(ctx context.Context, groupID int64) ([]Account, error) {
	return nil, nil
placeholder
func (m *mockAccountRepoForPlatform) ListActive(ctx context.Context) ([]Account, error) {
	return nil, nil
placeholder
func (m *mockAccountRepoForPlatform) ListByPlatform(ctx context.Context, platform string) ([]Account, error) {
	return nil, nil
placeholder
func (m *mockAccountRepoForPlatform) UpdateLastUsed(ctx context.Context, id int64) error {
	return nil
placeholder
func (m *mockAccountRepoForPlatform) BatchUpdateLastUsed(ctx context.Context, updates map[int64]time.Time) error {
	return nil
placeholder
func (m *mockAccountRepoForPlatform) SetError(ctx context.Context, id int64, errorMsg string) error {
	return nil
placeholder
func (m *mockAccountRepoForPlatform) SetSchedulable(ctx context.Context, id int64, schedulable bool) error {
	return nil
placeholder
func (m *mockAccountRepoForPlatform) BindGroups(ctx context.Context, accountID int64, groupIDs []int64) error {
	return nil
placeholder
func (m *mockAccountRepoForPlatform) ListSchedulable(ctx context.Context) ([]Account, error) {
	return nil, nil
placeholder
func (m *mockAccountRepoForPlatform) ListSchedulableByGroupID(ctx context.Context, groupID int64) ([]Account, error) {
	return nil, nil
placeholder
func (m *mockAccountRepoForPlatform) ListSchedulableByPlatforms(ctx context.Context, platforms []string) ([]Account, error) {
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
func (m *mockAccountRepoForPlatform) ListSchedulableByGroupIDAndPlatforms(ctx context.Context, groupID int64, platforms []string) ([]Account, error) {
	return m.ListSchedulableByPlatforms(ctx, platforms)
placeholder
func (m *mockAccountRepoForPlatform) SetRateLimited(ctx context.Context, id int64, resetAt time.Time) error {
	return nil
placeholder
func (m *mockAccountRepoForPlatform) SetOverloaded(ctx context.Context, id int64, until time.Time) error {
	return nil
placeholder
func (m *mockAccountRepoForPlatform) ClearRateLimit(ctx context.Context, id int64) error {
	return nil
placeholder
func (m *mockAccountRepoForPlatform) UpdateSessionWindow(ctx context.Context, id int64, start, end *time.Time, status string) error {
	return nil
placeholder
func (m *mockAccountRepoForPlatform) UpdateExtra(ctx context.Context, id int64, updates map[string]any) error {
	return nil
placeholder
func (m *mockAccountRepoForPlatform) BulkUpdate(ctx context.Context, ids []int64, updates AccountBulkUpdate) (int64, error) {
	return 0, nil
placeholder

// Verify interface implementation
var _ AccountRepository = (*mockAccountRepoForPlatform)(nil)

// mockGatewayCacheForPlatform 单平台测试用的 cache mock
type mockGatewayCacheForPlatform struct {
	sessionBindings map[string]int64
placeholder

func (m *mockGatewayCacheForPlatform) GetSessionAccountID(ctx context.Context, sessionHash string) (int64, error) {
	if id, ok := m.sessionBindings[sessionHash]; ok {
		return id, nil
placeholder
	return 0, errors.New("not found")
placeholder

func (m *mockGatewayCacheForPlatform) SetSessionAccountID(ctx context.Context, sessionHash string, accountID int64, ttl time.Duration) error {
	if m.sessionBindings == nil {
		m.sessionBindings = make(map[string]int64)
placeholder
	m.sessionBindings[sessionHash] = accountID
	return nil
placeholder

func (m *mockGatewayCacheForPlatform) RefreshSessionTTL(ctx context.Context, sessionHash string, ttl time.Duration) error {
	return nil
placeholder

func ptr[T any](v T) *T {
	return &v
placeholder

// TestGatewayService_SelectAccountForModelWithPlatform_Anthropic 测试 anthropic 单平台选择
func TestGatewayService_SelectAccountForModelWithPlatform_Anthropic(t *testing.T) {
	ctx := context.Background()

	repo := &mockAccountRepoForPlatform{
		accounts: []Account{
			{ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
			{ID: 2, Platform: PlatformAnthropic, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
			{ID: 3, Platform: PlatformAntigravity, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder, // 应被隔离
	placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder
	for i := range repo.accounts {
		repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
placeholder

	cache := &mockGatewayCacheForPlatform{placeholder

	svc := &GatewayService{
		accountRepo: repo,
		cache:       cache,
		cfg:         testConfig(),
placeholder

	acc, err := svc.selectAccountForModelWithPlatform(ctx, nil, "", "claude-3-5-sonnet-20241022", nil, PlatformAnthropic)
placeholder
	require.NotNil(t, acc)
	require.Equal(t, int64(1), acc.ID, "应选择优先级最高的 anthropic 账户")
	require.Equal(t, PlatformAnthropic, acc.Platform, "应只返回 anthropic 平台账户")
placeholder

// TestGatewayService_SelectAccountForModelWithPlatform_Antigravity 测试 antigravity 单平台选择
func TestGatewayService_SelectAccountForModelWithPlatform_Antigravity(t *testing.T) {
	ctx := context.Background()

	repo := &mockAccountRepoForPlatform{
		accounts: []Account{
			{ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder, // 应被隔离
			{ID: 2, Platform: PlatformAntigravity, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
	placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder
	for i := range repo.accounts {
		repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
placeholder

	cache := &mockGatewayCacheForPlatform{placeholder

	svc := &GatewayService{
		accountRepo: repo,
		cache:       cache,
		cfg:         testConfig(),
placeholder

	acc, err := svc.selectAccountForModelWithPlatform(ctx, nil, "", "claude-3-5-sonnet-20241022", nil, PlatformAntigravity)
placeholder
	require.NotNil(t, acc)
	require.Equal(t, int64(2), acc.ID)
	require.Equal(t, PlatformAntigravity, acc.Platform, "应只返回 antigravity 平台账户")
placeholder

// TestGatewayService_SelectAccountForModelWithPlatform_PriorityAndLastUsed 测试优先级和最后使用时间
func TestGatewayService_SelectAccountForModelWithPlatform_PriorityAndLastUsed(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	repo := &mockAccountRepoForPlatform{
		accounts: []Account{
			{ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: true, LastUsedAt: ptr(now.Add(-1 * time.Hour))placeholder,
			{ID: 2, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: true, LastUsedAt: ptr(now.Add(-2 * time.Hour))placeholder,
	placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder
	for i := range repo.accounts {
		repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
placeholder

	cache := &mockGatewayCacheForPlatform{placeholder

	svc := &GatewayService{
		accountRepo: repo,
		cache:       cache,
		cfg:         testConfig(),
placeholder

	acc, err := svc.selectAccountForModelWithPlatform(ctx, nil, "", "claude-3-5-sonnet-20241022", nil, PlatformAnthropic)
placeholder
	require.NotNil(t, acc)
	require.Equal(t, int64(2), acc.ID, "同优先级应选择最久未用的账户")
placeholder

func TestGatewayService_SelectAccountForModelWithPlatform_GeminiOAuthPreference(t *testing.T) {
	ctx := context.Background()

	repo := &mockAccountRepoForPlatform{
		accounts: []Account{
			{ID: 1, Platform: PlatformGemini, Priority: 1, Status: StatusActive, Schedulable: true, Type: AccountTypeApiKeyplaceholder,
			{ID: 2, Platform: PlatformGemini, Priority: 1, Status: StatusActive, Schedulable: true, Type: AccountTypeOAuthplaceholder,
	placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder
	for i := range repo.accounts {
		repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
placeholder

	cache := &mockGatewayCacheForPlatform{placeholder

	svc := &GatewayService{
		accountRepo: repo,
		cache:       cache,
		cfg:         testConfig(),
placeholder

	acc, err := svc.selectAccountForModelWithPlatform(ctx, nil, "", "gemini-2.5-pro", nil, PlatformGemini)
placeholder
	require.NotNil(t, acc)
	require.Equal(t, int64(2), acc.ID, "同优先级且未使用时应优先选择OAuth账户")
placeholder

// TestGatewayService_SelectAccountForModelWithPlatform_NoAvailableAccounts 测试无可用账户
func TestGatewayService_SelectAccountForModelWithPlatform_NoAvailableAccounts(t *testing.T) {
	ctx := context.Background()

	repo := &mockAccountRepoForPlatform{
		accounts:     []Account{placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder

	cache := &mockGatewayCacheForPlatform{placeholder

	svc := &GatewayService{
		accountRepo: repo,
		cache:       cache,
		cfg:         testConfig(),
placeholder

	acc, err := svc.selectAccountForModelWithPlatform(ctx, nil, "", "claude-3-5-sonnet-20241022", nil, PlatformAnthropic)
placeholder
	require.Nil(t, acc)
	require.Contains(t, err.Error(), "no available accounts")
placeholder

// TestGatewayService_SelectAccountForModelWithPlatform_AllExcluded 测试所有账户被排除
func TestGatewayService_SelectAccountForModelWithPlatform_AllExcluded(t *testing.T) {
	ctx := context.Background()

	repo := &mockAccountRepoForPlatform{
		accounts: []Account{
			{ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
			{ID: 2, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
	placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder
	for i := range repo.accounts {
		repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
placeholder

	cache := &mockGatewayCacheForPlatform{placeholder

	svc := &GatewayService{
		accountRepo: repo,
		cache:       cache,
		cfg:         testConfig(),
placeholder

	excludedIDs := map[int64]struct{placeholder{1: {placeholder, 2: {placeholderplaceholder
	acc, err := svc.selectAccountForModelWithPlatform(ctx, nil, "", "claude-3-5-sonnet-20241022", excludedIDs, PlatformAnthropic)
placeholder
	require.Nil(t, acc)
placeholder

// TestGatewayService_SelectAccountForModelWithPlatform_Schedulability 测试账户可调度性检查
func TestGatewayService_SelectAccountForModelWithPlatform_Schedulability(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name       string
		accounts   []Account
		expectedID int64
placeholder{
		{
			name: "过载账户被跳过",
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: true, OverloadUntil: ptr(now.Add(1 * time.Hour))placeholder,
				{ID: 2, Platform: PlatformAnthropic, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
		placeholder,
			expectedID: 2,
	placeholder,
		{
			name: "限流账户被跳过",
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: true, RateLimitResetAt: ptr(now.Add(1 * time.Hour))placeholder,
				{ID: 2, Platform: PlatformAnthropic, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
		placeholder,
			expectedID: 2,
	placeholder,
		{
			name: "非active账户被跳过",
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: "error", Schedulable: trueplaceholder,
				{ID: 2, Platform: PlatformAnthropic, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
		placeholder,
			expectedID: 2,
	placeholder,
		{
			name: "schedulable=false被跳过",
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: falseplaceholder,
				{ID: 2, Platform: PlatformAnthropic, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
		placeholder,
			expectedID: 2,
	placeholder,
		{
			name: "过期的过载账户可调度",
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: true, OverloadUntil: ptr(now.Add(-1 * time.Hour))placeholder,
				{ID: 2, Platform: PlatformAnthropic, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
		placeholder,
			expectedID: 1,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAccountRepoForPlatform{
				accounts:     tt.accounts,
				accountsByID: map[int64]*Account{placeholder,
		placeholder
			for i := range repo.accounts {
				repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
		placeholder

			cache := &mockGatewayCacheForPlatform{placeholder

			svc := &GatewayService{
				accountRepo: repo,
				cache:       cache,
				cfg:         testConfig(),
		placeholder

			acc, err := svc.selectAccountForModelWithPlatform(ctx, nil, "", "claude-3-5-sonnet-20241022", nil, PlatformAnthropic)
		placeholder
			require.NotNil(t, acc)
			require.Equal(t, tt.expectedID, acc.ID)
	placeholder)
placeholder
placeholder

// TestGatewayService_SelectAccountForModelWithPlatform_StickySession 测试粘性会话
func TestGatewayService_SelectAccountForModelWithPlatform_StickySession(t *testing.T) {
	ctx := context.Background()

	t.Run("粘性会话命中-同平台", func(t *testing.T) {
		repo := &mockAccountRepoForPlatform{
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
				{ID: 2, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
		placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder
		for i := range repo.accounts {
			repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
	placeholder

		cache := &mockGatewayCacheForPlatform{
			sessionBindings: map[string]int64{"session-123": 1placeholder,
	placeholder

		svc := &GatewayService{
			accountRepo: repo,
			cache:       cache,
			cfg:         testConfig(),
	placeholder

		acc, err := svc.selectAccountForModelWithPlatform(ctx, nil, "session-123", "claude-3-5-sonnet-20241022", nil, PlatformAnthropic)
	placeholder
		require.NotNil(t, acc)
		require.Equal(t, int64(1), acc.ID, "应返回粘性会话绑定的账户")
placeholder)

	t.Run("粘性会话不匹配平台-降级选择", func(t *testing.T) {
		repo := &mockAccountRepoForPlatform{
			accounts: []Account{
				{ID: 1, Platform: PlatformAntigravity, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder, // 粘性会话绑定但平台不匹配
				{ID: 2, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
		placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder
		for i := range repo.accounts {
			repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
	placeholder

		cache := &mockGatewayCacheForPlatform{
			sessionBindings: map[string]int64{"session-123": 1placeholder, // 绑定 antigravity 账户
	placeholder

		svc := &GatewayService{
			accountRepo: repo,
			cache:       cache,
			cfg:         testConfig(),
	placeholder

		// 请求 anthropic 平台，但粘性会话绑定的是 antigravity 账户
		acc, err := svc.selectAccountForModelWithPlatform(ctx, nil, "session-123", "claude-3-5-sonnet-20241022", nil, PlatformAnthropic)
	placeholder
		require.NotNil(t, acc)
		require.Equal(t, int64(2), acc.ID, "粘性会话账户平台不匹配，应降级选择同平台账户")
		require.Equal(t, PlatformAnthropic, acc.Platform)
placeholder)

	t.Run("粘性会话账户被排除-降级选择", func(t *testing.T) {
		repo := &mockAccountRepoForPlatform{
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
				{ID: 2, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
		placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder
		for i := range repo.accounts {
			repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
	placeholder

		cache := &mockGatewayCacheForPlatform{
			sessionBindings: map[string]int64{"session-123": 1placeholder,
	placeholder

		svc := &GatewayService{
			accountRepo: repo,
			cache:       cache,
			cfg:         testConfig(),
	placeholder

		excludedIDs := map[int64]struct{placeholder{1: {placeholderplaceholder
		acc, err := svc.selectAccountForModelWithPlatform(ctx, nil, "session-123", "claude-3-5-sonnet-20241022", excludedIDs, PlatformAnthropic)
	placeholder
		require.NotNil(t, acc)
		require.Equal(t, int64(2), acc.ID, "粘性会话账户被排除，应选择其他账户")
placeholder)

	t.Run("粘性会话账户不可调度-降级选择", func(t *testing.T) {
		repo := &mockAccountRepoForPlatform{
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 2, Status: "error", Schedulable: trueplaceholder,
				{ID: 2, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
		placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder
		for i := range repo.accounts {
			repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
	placeholder

		cache := &mockGatewayCacheForPlatform{
			sessionBindings: map[string]int64{"session-123": 1placeholder,
	placeholder

		svc := &GatewayService{
			accountRepo: repo,
			cache:       cache,
			cfg:         testConfig(),
	placeholder

		acc, err := svc.selectAccountForModelWithPlatform(ctx, nil, "session-123", "claude-3-5-sonnet-20241022", nil, PlatformAnthropic)
	placeholder
		require.NotNil(t, acc)
		require.Equal(t, int64(2), acc.ID, "粘性会话账户不可调度，应选择其他账户")
placeholder)
placeholder

func TestGatewayService_isModelSupportedByAccount(t *testing.T) {
	svc := &GatewayService{placeholder

	tests := []struct {
		name     string
		account  *Account
		model    string
		expected bool
placeholder{
		{
			name:     "Antigravity平台-支持claude模型",
			account:  &Account{Platform: PlatformAntigravityplaceholder,
			model:    "claude-3-5-sonnet-20241022",
			expected: true,
	placeholder,
		{
			name:     "Antigravity平台-支持gemini模型",
			account:  &Account{Platform: PlatformAntigravityplaceholder,
			model:    "gemini-2.5-flash",
			expected: true,
	placeholder,
		{
			name:     "Antigravity平台-不支持gpt模型",
			account:  &Account{Platform: PlatformAntigravityplaceholder,
			model:    "gpt-4",
			expected: false,
	placeholder,
		{
			name:     "Anthropic平台-无映射配置-支持所有模型",
			account:  &Account{Platform: PlatformAnthropicplaceholder,
			model:    "claude-3-5-sonnet-20241022",
			expected: true,
	placeholder,
		{
			name: "Anthropic平台-有映射配置-只支持配置的模型",
			account: &Account{
				Platform:    PlatformAnthropic,
		placeholder"model_mapping": map[string]any{"claude-opus-4": "x"placeholderplaceholder,
		placeholder,
			model:    "claude-3-5-sonnet-20241022",
			expected: false,
	placeholder,
		{
			name: "Anthropic平台-有映射配置-支持配置的模型",
			account: &Account{
				Platform:    PlatformAnthropic,
		placeholder"model_mapping": map[string]any{"claude-3-5-sonnet-20241022": "x"placeholderplaceholder,
		placeholder,
			model:    "claude-3-5-sonnet-20241022",
			expected: true,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.isModelSupportedByAccount(tt.account, tt.model)
			require.Equal(t, tt.expected, got)
	placeholder)
placeholder
placeholder

// TestGatewayService_selectAccountWithMixedScheduling 测试混合调度
func TestGatewayService_selectAccountWithMixedScheduling(t *testing.T) {
	ctx := context.Background()

	t.Run("混合调度-Gemini优先选择OAuth账户", func(t *testing.T) {
		repo := &mockAccountRepoForPlatform{
			accounts: []Account{
				{ID: 1, Platform: PlatformGemini, Priority: 1, Status: StatusActive, Schedulable: true, Type: AccountTypeApiKeyplaceholder,
				{ID: 2, Platform: PlatformGemini, Priority: 1, Status: StatusActive, Schedulable: true, Type: AccountTypeOAuthplaceholder,
		placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder
		for i := range repo.accounts {
			repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
	placeholder

		cache := &mockGatewayCacheForPlatform{placeholder

		svc := &GatewayService{
			accountRepo: repo,
			cache:       cache,
			cfg:         testConfig(),
	placeholder

		acc, err := svc.selectAccountWithMixedScheduling(ctx, nil, "", "gemini-2.5-pro", nil, PlatformGemini)
	placeholder
		require.NotNil(t, acc)
		require.Equal(t, int64(2), acc.ID, "同优先级且未使用时应优先选择OAuth账户")
placeholder)

	t.Run("混合调度-包含启用mixed_scheduling的antigravity账户", func(t *testing.T) {
		repo := &mockAccountRepoForPlatform{
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
				{ID: 2, Platform: PlatformAntigravity, Priority: 1, Status: StatusActive, Schedulable: true, Extra: map[string]any{"mixed_scheduling": trueplaceholderplaceholder,
		placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder
		for i := range repo.accounts {
			repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
	placeholder

		cache := &mockGatewayCacheForPlatform{placeholder

		svc := &GatewayService{
			accountRepo: repo,
			cache:       cache,
			cfg:         testConfig(),
	placeholder

		acc, err := svc.selectAccountWithMixedScheduling(ctx, nil, "", "claude-3-5-sonnet-20241022", nil, PlatformAnthropic)
	placeholder
		require.NotNil(t, acc)
		require.Equal(t, int64(2), acc.ID, "应选择优先级最高的账户（包含启用混合调度的antigravity）")
placeholder)

	t.Run("混合调度-过滤未启用mixed_scheduling的antigravity账户", func(t *testing.T) {
		repo := &mockAccountRepoForPlatform{
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
				{ID: 2, Platform: PlatformAntigravity, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder, // 未启用 mixed_scheduling
		placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder
		for i := range repo.accounts {
			repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
	placeholder

		cache := &mockGatewayCacheForPlatform{placeholder

		svc := &GatewayService{
			accountRepo: repo,
			cache:       cache,
			cfg:         testConfig(),
	placeholder

		acc, err := svc.selectAccountWithMixedScheduling(ctx, nil, "", "claude-3-5-sonnet-20241022", nil, PlatformAnthropic)
	placeholder
		require.NotNil(t, acc)
		require.Equal(t, int64(1), acc.ID, "未启用mixed_scheduling的antigravity账户应被过滤")
		require.Equal(t, PlatformAnthropic, acc.Platform)
placeholder)

	t.Run("混合调度-粘性会话命中启用mixed_scheduling的antigravity账户", func(t *testing.T) {
		repo := &mockAccountRepoForPlatform{
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
				{ID: 2, Platform: PlatformAntigravity, Priority: 2, Status: StatusActive, Schedulable: true, Extra: map[string]any{"mixed_scheduling": trueplaceholderplaceholder,
		placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder
		for i := range repo.accounts {
			repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
	placeholder

		cache := &mockGatewayCacheForPlatform{
			sessionBindings: map[string]int64{"session-123": 2placeholder,
	placeholder

		svc := &GatewayService{
			accountRepo: repo,
			cache:       cache,
			cfg:         testConfig(),
	placeholder

		acc, err := svc.selectAccountWithMixedScheduling(ctx, nil, "session-123", "claude-3-5-sonnet-20241022", nil, PlatformAnthropic)
	placeholder
		require.NotNil(t, acc)
		require.Equal(t, int64(2), acc.ID, "应返回粘性会话绑定的启用mixed_scheduling的antigravity账户")
placeholder)

	t.Run("混合调度-粘性会话命中未启用mixed_scheduling的antigravity账户-降级选择", func(t *testing.T) {
		repo := &mockAccountRepoForPlatform{
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
				{ID: 2, Platform: PlatformAntigravity, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder, // 未启用 mixed_scheduling
		placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder
		for i := range repo.accounts {
			repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
	placeholder

		cache := &mockGatewayCacheForPlatform{
			sessionBindings: map[string]int64{"session-123": 2placeholder,
	placeholder

		svc := &GatewayService{
			accountRepo: repo,
			cache:       cache,
			cfg:         testConfig(),
	placeholder

		acc, err := svc.selectAccountWithMixedScheduling(ctx, nil, "session-123", "claude-3-5-sonnet-20241022", nil, PlatformAnthropic)
	placeholder
		require.NotNil(t, acc)
		require.Equal(t, int64(1), acc.ID, "粘性会话绑定的账户未启用mixed_scheduling，应降级选择anthropic账户")
placeholder)

	t.Run("混合调度-仅有启用mixed_scheduling的antigravity账户", func(t *testing.T) {
		repo := &mockAccountRepoForPlatform{
			accounts: []Account{
				{ID: 1, Platform: PlatformAntigravity, Priority: 1, Status: StatusActive, Schedulable: true, Extra: map[string]any{"mixed_scheduling": trueplaceholderplaceholder,
		placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder
		for i := range repo.accounts {
			repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
	placeholder

		cache := &mockGatewayCacheForPlatform{placeholder

		svc := &GatewayService{
			accountRepo: repo,
			cache:       cache,
			cfg:         testConfig(),
	placeholder

		acc, err := svc.selectAccountWithMixedScheduling(ctx, nil, "", "claude-3-5-sonnet-20241022", nil, PlatformAnthropic)
	placeholder
		require.NotNil(t, acc)
		require.Equal(t, int64(1), acc.ID)
		require.Equal(t, PlatformAntigravity, acc.Platform)
placeholder)

	t.Run("混合调度-无可用账户", func(t *testing.T) {
		repo := &mockAccountRepoForPlatform{
			accounts: []Account{
				{ID: 1, Platform: PlatformAntigravity, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder, // 未启用 mixed_scheduling
		placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder
		for i := range repo.accounts {
			repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
	placeholder

		cache := &mockGatewayCacheForPlatform{placeholder

		svc := &GatewayService{
			accountRepo: repo,
			cache:       cache,
			cfg:         testConfig(),
	placeholder

		acc, err := svc.selectAccountWithMixedScheduling(ctx, nil, "", "claude-3-5-sonnet-20241022", nil, PlatformAnthropic)
	placeholder
		require.Nil(t, acc)
		require.Contains(t, err.Error(), "no available accounts")
placeholder)
placeholder

// TestAccount_IsMixedSchedulingEnabled 测试混合调度开关检查
func TestAccount_IsMixedSchedulingEnabled(t *testing.T) {
	tests := []struct {
		name     string
		account  Account
		expected bool
placeholder{
		{
			name:     "非antigravity平台-返回false",
			account:  Account{Platform: PlatformAnthropicplaceholder,
			expected: false,
	placeholder,
		{
			name:     "antigravity平台-无extra-返回false",
			account:  Account{Platform: PlatformAntigravityplaceholder,
			expected: false,
	placeholder,
		{
			name:     "antigravity平台-extra无mixed_scheduling-返回false",
			account:  Account{Platform: PlatformAntigravity, Extra: map[string]any{placeholderplaceholder,
			expected: false,
	placeholder,
		{
			name:     "antigravity平台-mixed_scheduling=false-返回false",
			account:  Account{Platform: PlatformAntigravity, Extra: map[string]any{"mixed_scheduling": falseplaceholderplaceholder,
			expected: false,
	placeholder,
		{
			name:     "antigravity平台-mixed_scheduling=true-返回true",
			account:  Account{Platform: PlatformAntigravity, Extra: map[string]any{"mixed_scheduling": trueplaceholderplaceholder,
			expected: true,
	placeholder,
		{
			name:     "antigravity平台-mixed_scheduling非bool类型-返回false",
			account:  Account{Platform: PlatformAntigravity, Extra: map[string]any{"mixed_scheduling": "true"placeholderplaceholder,
			expected: false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.account.IsMixedSchedulingEnabled()
			require.Equal(t, tt.expected, got)
	placeholder)
placeholder
placeholder

// mockConcurrencyService for testing
type mockConcurrencyService struct {
	accountLoads      map[int64]*AccountLoadInfo
	accountWaitCounts map[int64]int
	acquireResults    map[int64]bool
placeholder

func (m *mockConcurrencyService) GetAccountsLoadBatch(ctx context.Context, accounts []AccountWithConcurrency) (map[int64]*AccountLoadInfo, error) {
	if m.accountLoads == nil {
		return map[int64]*AccountLoadInfo{placeholder, nil
placeholder
	result := make(map[int64]*AccountLoadInfo)
	for _, acc := range accounts {
		if load, ok := m.accountLoads[acc.ID]; ok {
			result[acc.ID] = load
	placeholder else {
			result[acc.ID] = &AccountLoadInfo{
				AccountID:          acc.ID,
				CurrentConcurrency: 0,
				WaitingCount:       0,
				LoadRate:           0,
		placeholder
	placeholder
placeholder
	return result, nil
placeholder

func (m *mockConcurrencyService) GetAccountWaitingCount(ctx context.Context, accountID int64) (int, error) {
	if m.accountWaitCounts == nil {
		return 0, nil
placeholder
	return m.accountWaitCounts[accountID], nil
placeholder

// TestGatewayService_SelectAccountWithLoadAwareness tests load-aware account selection
func TestGatewayService_SelectAccountWithLoadAwareness(t *testing.T) {
	ctx := context.Background()

	t.Run("禁用负载批量查询-降级到传统选择", func(t *testing.T) {
		repo := &mockAccountRepoForPlatform{
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: true, Concurrency: 5placeholder,
				{ID: 2, Platform: PlatformAnthropic, Priority: 2, Status: StatusActive, Schedulable: true, Concurrency: 5placeholder,
		placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder
		for i := range repo.accounts {
			repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
	placeholder

		cache := &mockGatewayCacheForPlatform{placeholder

		cfg := testConfig()
		cfg.Gateway.Scheduling.LoadBatchEnabled = false

		svc := &GatewayService{
			accountRepo:        repo,
			cache:              cache,
			cfg:                cfg,
			concurrencyService: nil, // No concurrency service
	placeholder

		result, err := svc.SelectAccountWithLoadAwareness(ctx, nil, "", "claude-3-5-sonnet-20241022", nil)
	placeholder
		require.NotNil(t, result)
		require.NotNil(t, result.Account)
		require.Equal(t, int64(1), result.Account.ID, "应选择优先级最高的账号")
placeholder)

	t.Run("无ConcurrencyService-降级到传统选择", func(t *testing.T) {
		repo := &mockAccountRepoForPlatform{
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 2, Status: StatusActive, Schedulable: true, Concurrency: 5placeholder,
				{ID: 2, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: true, Concurrency: 5placeholder,
		placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder
		for i := range repo.accounts {
			repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
	placeholder

		cache := &mockGatewayCacheForPlatform{placeholder

		cfg := testConfig()
		cfg.Gateway.Scheduling.LoadBatchEnabled = true

		svc := &GatewayService{
			accountRepo:        repo,
			cache:              cache,
			cfg:                cfg,
			concurrencyService: nil,
	placeholder

		result, err := svc.SelectAccountWithLoadAwareness(ctx, nil, "", "claude-3-5-sonnet-20241022", nil)
	placeholder
		require.NotNil(t, result)
		require.NotNil(t, result.Account)
		require.Equal(t, int64(2), result.Account.ID, "应选择优先级最高的账号")
placeholder)

	t.Run("排除账号-不选择被排除的账号", func(t *testing.T) {
		repo := &mockAccountRepoForPlatform{
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: true, Concurrency: 5placeholder,
				{ID: 2, Platform: PlatformAnthropic, Priority: 2, Status: StatusActive, Schedulable: true, Concurrency: 5placeholder,
		placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder
		for i := range repo.accounts {
			repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
	placeholder

		cache := &mockGatewayCacheForPlatform{placeholder

		cfg := testConfig()
		cfg.Gateway.Scheduling.LoadBatchEnabled = false

		svc := &GatewayService{
			accountRepo:        repo,
			cache:              cache,
			cfg:                cfg,
			concurrencyService: nil,
	placeholder

		excludedIDs := map[int64]struct{placeholder{1: {placeholderplaceholder
		result, err := svc.SelectAccountWithLoadAwareness(ctx, nil, "", "claude-3-5-sonnet-20241022", excludedIDs)
	placeholder
		require.NotNil(t, result)
		require.NotNil(t, result.Account)
		require.Equal(t, int64(2), result.Account.ID, "不应选择被排除的账号")
placeholder)

	t.Run("无可用账号-返回错误", func(t *testing.T) {
		repo := &mockAccountRepoForPlatform{
			accounts:     []Account{placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder

		cache := &mockGatewayCacheForPlatform{placeholder

		cfg := testConfig()
		cfg.Gateway.Scheduling.LoadBatchEnabled = false

		svc := &GatewayService{
			accountRepo:        repo,
			cache:              cache,
			cfg:                cfg,
			concurrencyService: nil,
	placeholder

		result, err := svc.SelectAccountWithLoadAwareness(ctx, nil, "", "claude-3-5-sonnet-20241022", nil)
	placeholder
		require.Nil(t, result)
		require.Contains(t, err.Error(), "no available accounts")
placeholder)
placeholder
