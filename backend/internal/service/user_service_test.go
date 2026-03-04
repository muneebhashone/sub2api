//go:build unit

package service

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

// --- mock: UserRepository ---

type mockUserRepo struct {
	updateBalanceErr error
	updateBalanceFn  func(ctx context.Context, id int64, amount float64) error
placeholder

func (m *mockUserRepo) Create(context.Context, *User) error               { return nil placeholder
func (m *mockUserRepo) GetByID(context.Context, int64) (*User, error)     { return &User{placeholder, nil placeholder
func (m *mockUserRepo) GetByEmail(context.Context, string) (*User, error) { return &User{placeholder, nil placeholder
func (m *mockUserRepo) GetFirstAdmin(context.Context) (*User, error)      { return &User{placeholder, nil placeholder
func (m *mockUserRepo) Update(context.Context, *User) error               { return nil placeholder
func (m *mockUserRepo) Delete(context.Context, int64) error               { return nil placeholder
func (m *mockUserRepo) List(context.Context, pagination.PaginationParams) ([]User, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (m *mockUserRepo) ListWithFilters(context.Context, pagination.PaginationParams, UserListFilters) ([]User, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (m *mockUserRepo) UpdateBalance(ctx context.Context, id int64, amount float64) error {
	if m.updateBalanceFn != nil {
		return m.updateBalanceFn(ctx, id, amount)
placeholder
	return m.updateBalanceErr
placeholder
func (m *mockUserRepo) DeductBalance(context.Context, int64, float64) error { return nil placeholder
func (m *mockUserRepo) UpdateConcurrency(context.Context, int64, int) error { return nil placeholder
func (m *mockUserRepo) ExistsByEmail(context.Context, string) (bool, error) { return false, nil placeholder
func (m *mockUserRepo) RemoveGroupFromAllowedGroups(context.Context, int64) (int64, error) {
	return 0, nil
placeholder
func (m *mockUserRepo) AddGroupToAllowedGroups(context.Context, int64, int64) error { return nil placeholder
func (m *mockUserRepo) UpdateTotpSecret(context.Context, int64, *string) error      { return nil placeholder
func (m *mockUserRepo) EnableTotp(context.Context, int64) error                { return nil placeholder
func (m *mockUserRepo) DisableTotp(context.Context, int64) error               { return nil placeholder

// --- mock: APIKeyAuthCacheInvalidator ---

type mockAuthCacheInvalidator struct {
	invalidatedUserIDs []int64
	mu                 sync.Mutex
placeholder

func (m *mockAuthCacheInvalidator) InvalidateAuthCacheByKey(context.Context, string)    {placeholder
func (m *mockAuthCacheInvalidator) InvalidateAuthCacheByGroupID(context.Context, int64) {placeholder
func (m *mockAuthCacheInvalidator) InvalidateAuthCacheByUserID(_ context.Context, userID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.invalidatedUserIDs = append(m.invalidatedUserIDs, userID)
placeholder

// --- mock: BillingCache ---

type mockBillingCache struct {
	invalidateErr       error
	invalidateCallCount atomic.Int64
	invalidatedUserIDs  []int64
	mu                  sync.Mutex
placeholder

func (m *mockBillingCache) GetUserBalance(context.Context, int64) (float64, error)  { return 0, nil placeholder
func (m *mockBillingCache) SetUserBalance(context.Context, int64, float64) error    { return nil placeholder
func (m *mockBillingCache) DeductUserBalance(context.Context, int64, float64) error { return nil placeholder
func (m *mockBillingCache) InvalidateUserBalance(_ context.Context, userID int64) error {
	m.invalidateCallCount.Add(1)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.invalidatedUserIDs = append(m.invalidatedUserIDs, userID)
	return m.invalidateErr
placeholder
func (m *mockBillingCache) GetSubscriptionCache(context.Context, int64, int64) (*SubscriptionCacheData, error) {
	return nil, nil
placeholder
func (m *mockBillingCache) SetSubscriptionCache(context.Context, int64, int64, *SubscriptionCacheData) error {
	return nil
placeholder
func (m *mockBillingCache) UpdateSubscriptionUsage(context.Context, int64, int64, float64) error {
	return nil
placeholder
func (m *mockBillingCache) InvalidateSubscriptionCache(context.Context, int64, int64) error {
	return nil
placeholder
func (m *mockBillingCache) GetAPIKeyRateLimit(context.Context, int64) (*APIKeyRateLimitCacheData, error) {
	return nil, nil
placeholder
func (m *mockBillingCache) SetAPIKeyRateLimit(context.Context, int64, *APIKeyRateLimitCacheData) error {
	return nil
placeholder
func (m *mockBillingCache) UpdateAPIKeyRateLimitUsage(context.Context, int64, float64) error {
	return nil
placeholder
func (m *mockBillingCache) InvalidateAPIKeyRateLimit(context.Context, int64) error {
	return nil
placeholder

// --- 测试 ---

func TestUpdateBalance_Success(t *testing.T) {
	repo := &mockUserRepo{placeholder
	cache := &mockBillingCache{placeholder
	svc := NewUserService(repo, nil, cache)

	err := svc.UpdateBalance(context.Background(), 42, 100.0)
placeholder

	// 等待异步 goroutine 完成
	require.Eventually(t, func() bool {
		return cache.invalidateCallCount.Load() == 1
placeholder, 2*time.Second, 10*time.Millisecond, "应异步调用 InvalidateUserBalance")

	cache.mu.Lock()
	defer cache.mu.Unlock()
	require.Equal(t, []int64{42placeholder, cache.invalidatedUserIDs, "应对 userID=42 失效缓存")
placeholder

func TestUpdateBalance_NilBillingCache_NoPanic(t *testing.T) {
	repo := &mockUserRepo{placeholder
	svc := NewUserService(repo, nil, nil) // billingCache = nil

	err := svc.UpdateBalance(context.Background(), 1, 50.0)
	require.NoError(t, err, "billingCache 为 nil 时不应 panic")
placeholder

func TestUpdateBalance_CacheFailure_DoesNotAffectReturn(t *testing.T) {
	repo := &mockUserRepo{placeholder
	cache := &mockBillingCache{invalidateErr: errors.New("redis connection refused")placeholder
	svc := NewUserService(repo, nil, cache)

	err := svc.UpdateBalance(context.Background(), 99, 200.0)
	require.NoError(t, err, "缓存失效失败不应影响主流程返回值")

	// 等待异步 goroutine 完成（即使失败也应调用）
	require.Eventually(t, func() bool {
		return cache.invalidateCallCount.Load() == 1
placeholder, 2*time.Second, 10*time.Millisecond, "即使失败也应调用 InvalidateUserBalance")
placeholder

func TestUpdateBalance_RepoError_ReturnsError(t *testing.T) {
	repo := &mockUserRepo{updateBalanceErr: errors.New("database error")placeholder
	cache := &mockBillingCache{placeholder
	svc := NewUserService(repo, nil, cache)

	err := svc.UpdateBalance(context.Background(), 1, 100.0)
	require.Error(t, err, "repo 失败时应返回错误")
	require.Contains(t, err.Error(), "update balance")

	// repo 失败时不应触发缓存失效
	time.Sleep(100 * time.Millisecond)
	require.Equal(t, int64(0), cache.invalidateCallCount.Load(),
		"repo 失败时不应调用 InvalidateUserBalance")
placeholder

func TestUpdateBalance_WithAuthCacheInvalidator(t *testing.T) {
	repo := &mockUserRepo{placeholder
	auth := &mockAuthCacheInvalidator{placeholder
	cache := &mockBillingCache{placeholder
	svc := NewUserService(repo, auth, cache)

	err := svc.UpdateBalance(context.Background(), 77, 300.0)
placeholder

	// 验证 auth cache 同步失效
	auth.mu.Lock()
	require.Equal(t, []int64{77placeholder, auth.invalidatedUserIDs)
	auth.mu.Unlock()

	// 验证 billing cache 异步失效
	require.Eventually(t, func() bool {
		return cache.invalidateCallCount.Load() == 1
placeholder, 2*time.Second, 10*time.Millisecond)
placeholder

func TestNewUserService_FieldsAssignment(t *testing.T) {
	repo := &mockUserRepo{placeholder
	auth := &mockAuthCacheInvalidator{placeholder
	cache := &mockBillingCache{placeholder

	svc := NewUserService(repo, auth, cache)
	require.NotNil(t, svc)
	require.Equal(t, repo, svc.userRepo)
	require.Equal(t, auth, svc.authCacheInvalidator)
	require.Equal(t, cache, svc.billingCache)
placeholder
