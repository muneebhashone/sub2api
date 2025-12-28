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

// mockAccountRepoForMultiplatform 多平台测试用的 mock
type mockAccountRepoForMultiplatform struct {
	accounts          []Account
	accountsByID      map[int64]*Account
	listPlatformsFunc func(ctx context.Context, platforms []string) ([]Account, error)
placeholder

func (m *mockAccountRepoForMultiplatform) GetByID(ctx context.Context, id int64) (*Account, error) {
	if acc, ok := m.accountsByID[id]; ok {
		return acc, nil
placeholder
	return nil, errors.New("account not found")
placeholder

func (m *mockAccountRepoForMultiplatform) ListSchedulableByPlatforms(ctx context.Context, platforms []string) ([]Account, error) {
	if m.listPlatformsFunc != nil {
		return m.listPlatformsFunc(ctx, platforms)
placeholder
	// 过滤符合平台的账户
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

func (m *mockAccountRepoForMultiplatform) ListSchedulableByGroupIDAndPlatforms(ctx context.Context, groupID int64, platforms []string) ([]Account, error) {
	return m.ListSchedulableByPlatforms(ctx, platforms)
placeholder

// Stub methods to implement AccountRepository interface
func (m *mockAccountRepoForMultiplatform) Create(ctx context.Context, account *Account) error {
	return nil
placeholder
func (m *mockAccountRepoForMultiplatform) GetByCRSAccountID(ctx context.Context, crsAccountID string) (*Account, error) {
	return nil, nil
placeholder
func (m *mockAccountRepoForMultiplatform) Update(ctx context.Context, account *Account) error {
	return nil
placeholder
func (m *mockAccountRepoForMultiplatform) Delete(ctx context.Context, id int64) error { return nil placeholder
func (m *mockAccountRepoForMultiplatform) List(ctx context.Context, params pagination.PaginationParams) ([]Account, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (m *mockAccountRepoForMultiplatform) ListWithFilters(ctx context.Context, params pagination.PaginationParams, platform, accountType, status, search string) ([]Account, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (m *mockAccountRepoForMultiplatform) ListByGroup(ctx context.Context, groupID int64) ([]Account, error) {
	return nil, nil
placeholder
func (m *mockAccountRepoForMultiplatform) ListActive(ctx context.Context) ([]Account, error) {
	return nil, nil
placeholder
func (m *mockAccountRepoForMultiplatform) ListByPlatform(ctx context.Context, platform string) ([]Account, error) {
	return nil, nil
placeholder
func (m *mockAccountRepoForMultiplatform) UpdateLastUsed(ctx context.Context, id int64) error {
	return nil
placeholder
func (m *mockAccountRepoForMultiplatform) BatchUpdateLastUsed(ctx context.Context, updates map[int64]time.Time) error {
	return nil
placeholder
func (m *mockAccountRepoForMultiplatform) SetError(ctx context.Context, id int64, errorMsg string) error {
	return nil
placeholder
func (m *mockAccountRepoForMultiplatform) SetSchedulable(ctx context.Context, id int64, schedulable bool) error {
	return nil
placeholder
func (m *mockAccountRepoForMultiplatform) BindGroups(ctx context.Context, accountID int64, groupIDs []int64) error {
	return nil
placeholder
func (m *mockAccountRepoForMultiplatform) ListSchedulable(ctx context.Context) ([]Account, error) {
	return nil, nil
placeholder
func (m *mockAccountRepoForMultiplatform) ListSchedulableByGroupID(ctx context.Context, groupID int64) ([]Account, error) {
	return nil, nil
placeholder
func (m *mockAccountRepoForMultiplatform) ListSchedulableByPlatform(ctx context.Context, platform string) ([]Account, error) {
	return nil, nil
placeholder
func (m *mockAccountRepoForMultiplatform) ListSchedulableByGroupIDAndPlatform(ctx context.Context, groupID int64, platform string) ([]Account, error) {
	return nil, nil
placeholder
func (m *mockAccountRepoForMultiplatform) SetRateLimited(ctx context.Context, id int64, resetAt time.Time) error {
	return nil
placeholder
func (m *mockAccountRepoForMultiplatform) SetOverloaded(ctx context.Context, id int64, until time.Time) error {
	return nil
placeholder
func (m *mockAccountRepoForMultiplatform) ClearRateLimit(ctx context.Context, id int64) error {
	return nil
placeholder
func (m *mockAccountRepoForMultiplatform) UpdateSessionWindow(ctx context.Context, id int64, start, end *time.Time, status string) error {
	return nil
placeholder
func (m *mockAccountRepoForMultiplatform) UpdateExtra(ctx context.Context, id int64, updates map[string]any) error {
	return nil
placeholder
func (m *mockAccountRepoForMultiplatform) BulkUpdate(ctx context.Context, ids []int64, updates AccountBulkUpdate) (int64, error) {
	return 0, nil
placeholder

// Verify interface implementation
var _ AccountRepository = (*mockAccountRepoForMultiplatform)(nil)

// mockGatewayCacheForMultiplatform 多平台测试用的 cache mock
type mockGatewayCacheForMultiplatform struct {
	sessionBindings map[string]int64
placeholder

func (m *mockGatewayCacheForMultiplatform) GetSessionAccountID(ctx context.Context, sessionHash string) (int64, error) {
	if id, ok := m.sessionBindings[sessionHash]; ok {
		return id, nil
placeholder
	return 0, errors.New("not found")
placeholder

func (m *mockGatewayCacheForMultiplatform) SetSessionAccountID(ctx context.Context, sessionHash string, accountID int64, ttl time.Duration) error {
	if m.sessionBindings == nil {
		m.sessionBindings = make(map[string]int64)
placeholder
	m.sessionBindings[sessionHash] = accountID
	return nil
placeholder

func (m *mockGatewayCacheForMultiplatform) RefreshSessionTTL(ctx context.Context, sessionHash string, ttl time.Duration) error {
	return nil
placeholder

func ptr[T any](v T) *T {
	return &v
placeholder

func TestGatewayService_SelectAccountForModelWithExclusions_OnlyAnthropic(t *testing.T) {
	ctx := context.Background()

	repo := &mockAccountRepoForMultiplatform{
		accounts: []Account{
			{ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
			{ID: 2, Platform: PlatformAnthropic, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
	placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder
	for i := range repo.accounts {
		repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
placeholder

	cache := &mockGatewayCacheForMultiplatform{placeholder

	svc := &GatewayService{
		accountRepo: repo,
		cache:       cache,
placeholder

	acc, err := svc.selectAccountForModelWithPlatforms(ctx, nil, "", "claude-3-5-sonnet-20241022", nil, []string{PlatformAnthropic, PlatformAntigravityplaceholder)
placeholder
	require.NotNil(t, acc)
	require.Equal(t, int64(1), acc.ID, "应选择优先级最高的账户")
placeholder

func TestGatewayService_SelectAccountForModelWithExclusions_OnlyAntigravity(t *testing.T) {
	ctx := context.Background()

	repo := &mockAccountRepoForMultiplatform{
		accounts: []Account{
			{ID: 1, Platform: PlatformAntigravity, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
	placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder
	for i := range repo.accounts {
		repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
placeholder

	cache := &mockGatewayCacheForMultiplatform{placeholder

	svc := &GatewayService{
		accountRepo: repo,
		cache:       cache,
placeholder

	acc, err := svc.selectAccountForModelWithPlatforms(ctx, nil, "", "claude-3-5-sonnet-20241022", nil, []string{PlatformAnthropic, PlatformAntigravityplaceholder)
placeholder
	require.NotNil(t, acc)
	require.Equal(t, int64(1), acc.ID)
	require.Equal(t, PlatformAntigravity, acc.Platform)
placeholder

func TestGatewayService_SelectAccountForModelWithExclusions_MixedPlatforms_SamePriority(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	repo := &mockAccountRepoForMultiplatform{
		accounts: []Account{
			{ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: true, LastUsedAt: ptr(now.Add(-1 * time.Hour))placeholder,
			{ID: 2, Platform: PlatformAntigravity, Priority: 1, Status: StatusActive, Schedulable: true, LastUsedAt: ptr(now.Add(-2 * time.Hour))placeholder,
	placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder
	for i := range repo.accounts {
		repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
placeholder

	cache := &mockGatewayCacheForMultiplatform{placeholder

	svc := &GatewayService{
		accountRepo: repo,
		cache:       cache,
placeholder

	acc, err := svc.selectAccountForModelWithPlatforms(ctx, nil, "", "claude-3-5-sonnet-20241022", nil, []string{PlatformAnthropic, PlatformAntigravityplaceholder)
placeholder
	require.NotNil(t, acc)
	require.Equal(t, int64(2), acc.ID, "应选择最久未用的账户（Antigravity）")
placeholder

func TestGatewayService_SelectAccountForModelWithExclusions_MixedPlatforms_DiffPriority(t *testing.T) {
	ctx := context.Background()

	repo := &mockAccountRepoForMultiplatform{
		accounts: []Account{
			{ID: 1, Platform: PlatformAnthropic, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
			{ID: 2, Platform: PlatformAntigravity, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
	placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder
	for i := range repo.accounts {
		repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
placeholder

	cache := &mockGatewayCacheForMultiplatform{placeholder

	svc := &GatewayService{
		accountRepo: repo,
		cache:       cache,
placeholder

	acc, err := svc.selectAccountForModelWithPlatforms(ctx, nil, "", "claude-3-5-sonnet-20241022", nil, []string{PlatformAnthropic, PlatformAntigravityplaceholder)
placeholder
	require.NotNil(t, acc)
	require.Equal(t, int64(2), acc.ID, "应选择优先级更高的账户（Antigravity, priority=1）")
placeholder

func TestGatewayService_SelectAccountForModelWithExclusions_ModelNotSupported(t *testing.T) {
	ctx := context.Background()

	repo := &mockAccountRepoForMultiplatform{
		accounts: []Account{
			// Anthropic 账户配置了模型映射，只支持 other-model
			// 注意：model_mapping 需要是 map[string]any 格式
			{
				ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: true,
		placeholder"model_mapping": map[string]any{"other-model": "x"placeholderplaceholder,
		placeholder,
			// Antigravity 账户支持所有 claude 模型
			{ID: 2, Platform: PlatformAntigravity, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
	placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder
	for i := range repo.accounts {
		repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
placeholder

	cache := &mockGatewayCacheForMultiplatform{placeholder

	svc := &GatewayService{
		accountRepo: repo,
		cache:       cache,
placeholder

	acc, err := svc.selectAccountForModelWithPlatforms(ctx, nil, "", "claude-3-5-sonnet-20241022", nil, []string{PlatformAnthropic, PlatformAntigravityplaceholder)
placeholder
	require.NotNil(t, acc)
	require.Equal(t, int64(2), acc.ID, "Anthropic 不支持该模型，应选择 Antigravity")
placeholder

func TestGatewayService_SelectAccountForModelWithExclusions_NoAvailableAccounts(t *testing.T) {
	ctx := context.Background()

	repo := &mockAccountRepoForMultiplatform{
		accounts:     []Account{placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder

	cache := &mockGatewayCacheForMultiplatform{placeholder

	svc := &GatewayService{
		accountRepo: repo,
		cache:       cache,
placeholder

	acc, err := svc.selectAccountForModelWithPlatforms(ctx, nil, "", "claude-3-5-sonnet-20241022", nil, []string{PlatformAnthropic, PlatformAntigravityplaceholder)
placeholder
	require.Nil(t, acc)
	require.Contains(t, err.Error(), "no available accounts")
placeholder

func TestGatewayService_SelectAccountForModelWithExclusions_AllExcluded(t *testing.T) {
	ctx := context.Background()

	repo := &mockAccountRepoForMultiplatform{
		accounts: []Account{
			{ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
			{ID: 2, Platform: PlatformAntigravity, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
	placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder
	for i := range repo.accounts {
		repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
placeholder

	cache := &mockGatewayCacheForMultiplatform{placeholder

	svc := &GatewayService{
		accountRepo: repo,
		cache:       cache,
placeholder

	excludedIDs := map[int64]struct{placeholder{1: {placeholder, 2: {placeholderplaceholder
	acc, err := svc.selectAccountForModelWithPlatforms(ctx, nil, "", "claude-3-5-sonnet-20241022", excludedIDs, []string{PlatformAnthropic, PlatformAntigravityplaceholder)
placeholder
	require.Nil(t, acc)
placeholder

func TestGatewayService_SelectAccountForModelWithExclusions_Schedulability(t *testing.T) {
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
				{ID: 2, Platform: PlatformAntigravity, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
		placeholder,
			expectedID: 2,
	placeholder,
		{
			name: "限流账户被跳过",
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: true, RateLimitResetAt: ptr(now.Add(1 * time.Hour))placeholder,
				{ID: 2, Platform: PlatformAntigravity, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
		placeholder,
			expectedID: 2,
	placeholder,
		{
			name: "非active账户被跳过",
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: "error", Schedulable: trueplaceholder,
				{ID: 2, Platform: PlatformAntigravity, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
		placeholder,
			expectedID: 2,
	placeholder,
		{
			name: "schedulable=false被跳过",
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: falseplaceholder,
				{ID: 2, Platform: PlatformAntigravity, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
		placeholder,
			expectedID: 2,
	placeholder,
		{
			name: "过期的过载账户可调度",
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 1, Status: StatusActive, Schedulable: true, OverloadUntil: ptr(now.Add(-1 * time.Hour))placeholder,
				{ID: 2, Platform: PlatformAntigravity, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
		placeholder,
			expectedID: 1,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAccountRepoForMultiplatform{
				accounts:     tt.accounts,
				accountsByID: map[int64]*Account{placeholder,
		placeholder
			for i := range repo.accounts {
				repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
		placeholder

			cache := &mockGatewayCacheForMultiplatform{placeholder

			svc := &GatewayService{
				accountRepo: repo,
				cache:       cache,
		placeholder

			acc, err := svc.selectAccountForModelWithPlatforms(ctx, nil, "", "claude-3-5-sonnet-20241022", nil, []string{PlatformAnthropic, PlatformAntigravityplaceholder)
		placeholder
			require.NotNil(t, acc)
			require.Equal(t, tt.expectedID, acc.ID)
	placeholder)
placeholder
placeholder

func TestGatewayService_SelectAccountForModelWithExclusions_StickySession(t *testing.T) {
	ctx := context.Background()

	t.Run("粘性会话命中", func(t *testing.T) {
		repo := &mockAccountRepoForMultiplatform{
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
				{ID: 2, Platform: PlatformAntigravity, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
		placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder
		for i := range repo.accounts {
			repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
	placeholder

		cache := &mockGatewayCacheForMultiplatform{
			sessionBindings: map[string]int64{"session-123": 1placeholder,
	placeholder

		svc := &GatewayService{
			accountRepo: repo,
			cache:       cache,
	placeholder

		acc, err := svc.selectAccountForModelWithPlatforms(ctx, nil, "session-123", "claude-3-5-sonnet-20241022", nil, []string{PlatformAnthropic, PlatformAntigravityplaceholder)
	placeholder
		require.NotNil(t, acc)
		require.Equal(t, int64(1), acc.ID, "应返回粘性会话绑定的账户")
placeholder)

	t.Run("粘性会话账户被排除-降级选择", func(t *testing.T) {
		repo := &mockAccountRepoForMultiplatform{
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 2, Status: StatusActive, Schedulable: trueplaceholder,
				{ID: 2, Platform: PlatformAntigravity, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
		placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder
		for i := range repo.accounts {
			repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
	placeholder

		cache := &mockGatewayCacheForMultiplatform{
			sessionBindings: map[string]int64{"session-123": 1placeholder,
	placeholder

		svc := &GatewayService{
			accountRepo: repo,
			cache:       cache,
	placeholder

		excludedIDs := map[int64]struct{placeholder{1: {placeholderplaceholder
		acc, err := svc.selectAccountForModelWithPlatforms(ctx, nil, "session-123", "claude-3-5-sonnet-20241022", excludedIDs, []string{PlatformAnthropic, PlatformAntigravityplaceholder)
	placeholder
		require.NotNil(t, acc)
		require.Equal(t, int64(2), acc.ID, "粘性会话账户被排除，应选择其他账户")
placeholder)

	t.Run("粘性会话账户不可调度-降级选择", func(t *testing.T) {
		repo := &mockAccountRepoForMultiplatform{
			accounts: []Account{
				{ID: 1, Platform: PlatformAnthropic, Priority: 2, Status: "error", Schedulable: trueplaceholder,
				{ID: 2, Platform: PlatformAntigravity, Priority: 1, Status: StatusActive, Schedulable: trueplaceholder,
		placeholder,
			accountsByID: map[int64]*Account{placeholder,
	placeholder
		for i := range repo.accounts {
			repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
	placeholder

		cache := &mockGatewayCacheForMultiplatform{
			sessionBindings: map[string]int64{"session-123": 1placeholder,
	placeholder

		svc := &GatewayService{
			accountRepo: repo,
			cache:       cache,
	placeholder

		acc, err := svc.selectAccountForModelWithPlatforms(ctx, nil, "session-123", "claude-3-5-sonnet-20241022", nil, []string{PlatformAnthropic, PlatformAntigravityplaceholder)
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
