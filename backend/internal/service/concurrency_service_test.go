//go:build unit

package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// stubConcurrencyCacheForTest 用于并发服务单元测试的缓存桩
type stubConcurrencyCacheForTest struct {
	acquireResult  bool
	acquireErr     error
	releaseErr     error
	concurrency    int
	concurrencyErr error
	waitAllowed    bool
	waitErr        error
	waitCount      int
	waitCountErr   error
	loadBatch      map[int64]*AccountLoadInfo
	loadBatchErr   error
	usersLoadBatch map[int64]*UserLoadInfo
	usersLoadErr   error
	cleanupErr     error

	// 记录调用
	releasedAccountIDs []int64
	releasedRequestIDs []string
placeholder

var _ ConcurrencyCache = (*stubConcurrencyCacheForTest)(nil)

func (c *stubConcurrencyCacheForTest) AcquireAccountSlot(_ context.Context, _ int64, _ int, _ string) (bool, error) {
	return c.acquireResult, c.acquireErr
placeholder
func (c *stubConcurrencyCacheForTest) ReleaseAccountSlot(_ context.Context, accountID int64, requestID string) error {
	c.releasedAccountIDs = append(c.releasedAccountIDs, accountID)
	c.releasedRequestIDs = append(c.releasedRequestIDs, requestID)
	return c.releaseErr
placeholder
func (c *stubConcurrencyCacheForTest) GetAccountConcurrency(_ context.Context, _ int64) (int, error) {
	return c.concurrency, c.concurrencyErr
placeholder
func (c *stubConcurrencyCacheForTest) GetAccountConcurrencyBatch(_ context.Context, accountIDs []int64) (map[int64]int, error) {
	result := make(map[int64]int, len(accountIDs))
	for _, accountID := range accountIDs {
		if c.concurrencyErr != nil {
			return nil, c.concurrencyErr
	placeholder
		result[accountID] = c.concurrency
placeholder
	return result, nil
placeholder
func (c *stubConcurrencyCacheForTest) IncrementAccountWaitCount(_ context.Context, _ int64, _ int) (bool, error) {
	return c.waitAllowed, c.waitErr
placeholder
func (c *stubConcurrencyCacheForTest) DecrementAccountWaitCount(_ context.Context, _ int64) error {
	return nil
placeholder
func (c *stubConcurrencyCacheForTest) GetAccountWaitingCount(_ context.Context, _ int64) (int, error) {
	return c.waitCount, c.waitCountErr
placeholder
func (c *stubConcurrencyCacheForTest) AcquireUserSlot(_ context.Context, _ int64, _ int, _ string) (bool, error) {
	return c.acquireResult, c.acquireErr
placeholder
func (c *stubConcurrencyCacheForTest) ReleaseUserSlot(_ context.Context, _ int64, _ string) error {
	return c.releaseErr
placeholder
func (c *stubConcurrencyCacheForTest) GetUserConcurrency(_ context.Context, _ int64) (int, error) {
	return c.concurrency, c.concurrencyErr
placeholder
func (c *stubConcurrencyCacheForTest) IncrementWaitCount(_ context.Context, _ int64, _ int) (bool, error) {
	return c.waitAllowed, c.waitErr
placeholder
func (c *stubConcurrencyCacheForTest) DecrementWaitCount(_ context.Context, _ int64) error {
	return nil
placeholder
func (c *stubConcurrencyCacheForTest) GetAccountsLoadBatch(_ context.Context, _ []AccountWithConcurrency) (map[int64]*AccountLoadInfo, error) {
	return c.loadBatch, c.loadBatchErr
placeholder
func (c *stubConcurrencyCacheForTest) GetUsersLoadBatch(_ context.Context, _ []UserWithConcurrency) (map[int64]*UserLoadInfo, error) {
	return c.usersLoadBatch, c.usersLoadErr
placeholder
func (c *stubConcurrencyCacheForTest) CleanupExpiredAccountSlots(_ context.Context, _ int64) error {
	return c.cleanupErr
placeholder

func (c *stubConcurrencyCacheForTest) CleanupStaleProcessSlots(_ context.Context, _ string) error {
	return c.cleanupErr
placeholder

type trackingConcurrencyCache struct {
	stubConcurrencyCacheForTest
	cleanupPrefix string
placeholder

func (c *trackingConcurrencyCache) CleanupStaleProcessSlots(_ context.Context, prefix string) error {
	c.cleanupPrefix = prefix
	return c.cleanupErr
placeholder

func TestCleanupStaleProcessSlots_NilCache(t *testing.T) {
	svc := &ConcurrencyService{cache: nilplaceholder
	require.NoError(t, svc.CleanupStaleProcessSlots(context.Background()))
placeholder

func TestCleanupStaleProcessSlots_DelegatesPrefix(t *testing.T) {
	cache := &trackingConcurrencyCache{placeholder
	svc := NewConcurrencyService(cache)
	require.NoError(t, svc.CleanupStaleProcessSlots(context.Background()))
	require.Equal(t, RequestIDPrefix(), cache.cleanupPrefix)
placeholder

func TestAcquireAccountSlot_Success(t *testing.T) {
	cache := &stubConcurrencyCacheForTest{acquireResult: trueplaceholder
	svc := NewConcurrencyService(cache)

	result, err := svc.AcquireAccountSlot(context.Background(), 1, 5)
placeholder
	require.True(t, result.Acquired)
	require.NotNil(t, result.ReleaseFunc)
placeholder

func TestAcquireAccountSlot_Failure(t *testing.T) {
	cache := &stubConcurrencyCacheForTest{acquireResult: falseplaceholder
	svc := NewConcurrencyService(cache)

	result, err := svc.AcquireAccountSlot(context.Background(), 1, 5)
placeholder
	require.False(t, result.Acquired)
	require.Nil(t, result.ReleaseFunc)
placeholder

func TestAcquireAccountSlot_UnlimitedConcurrency(t *testing.T) {
	svc := NewConcurrencyService(&stubConcurrencyCacheForTest{placeholder)

	for _, maxConcurrency := range []int{0, -1placeholder {
		result, err := svc.AcquireAccountSlot(context.Background(), 1, maxConcurrency)
	placeholder
		require.True(t, result.Acquired, "maxConcurrency=%d 应无限制通过", maxConcurrency)
		require.NotNil(t, result.ReleaseFunc, "ReleaseFunc 应为 no-op 函数")
placeholder
placeholder

func TestAcquireAccountSlot_CacheError(t *testing.T) {
	cache := &stubConcurrencyCacheForTest{acquireErr: errors.New("redis down")placeholder
	svc := NewConcurrencyService(cache)

	result, err := svc.AcquireAccountSlot(context.Background(), 1, 5)
placeholder
	require.Nil(t, result)
placeholder

func TestAcquireAccountSlot_ReleaseDecrements(t *testing.T) {
	cache := &stubConcurrencyCacheForTest{acquireResult: trueplaceholder
	svc := NewConcurrencyService(cache)

	result, err := svc.AcquireAccountSlot(context.Background(), 42, 5)
placeholder
	require.True(t, result.Acquired)

	// 调用 ReleaseFunc 应释放槽位
	result.ReleaseFunc()

	require.Len(t, cache.releasedAccountIDs, 1)
	require.Equal(t, int64(42), cache.releasedAccountIDs[0])
	require.Len(t, cache.releasedRequestIDs, 1)
	require.NotEmpty(t, cache.releasedRequestIDs[0], "requestID 不应为空")
placeholder

func TestAcquireUserSlot_IndependentFromAccount(t *testing.T) {
	cache := &stubConcurrencyCacheForTest{acquireResult: trueplaceholder
	svc := NewConcurrencyService(cache)

	// 用户槽位获取应独立于账户槽位
	result, err := svc.AcquireUserSlot(context.Background(), 100, 3)
placeholder
	require.True(t, result.Acquired)
	require.NotNil(t, result.ReleaseFunc)
placeholder

func TestAcquireUserSlot_UnlimitedConcurrency(t *testing.T) {
	svc := NewConcurrencyService(&stubConcurrencyCacheForTest{placeholder)

	result, err := svc.AcquireUserSlot(context.Background(), 1, 0)
placeholder
	require.True(t, result.Acquired)
placeholder

func TestGenerateRequestID_UsesStablePrefixAndMonotonicCounter(t *testing.T) {
	id1 := generateRequestID()
	id2 := generateRequestID()
	require.NotEmpty(t, id1)
	require.NotEmpty(t, id2)

	p1 := strings.Split(id1, "-")
	p2 := strings.Split(id2, "-")
	require.Len(t, p1, 2)
	require.Len(t, p2, 2)
	require.Equal(t, p1[0], p2[0], "同一进程前缀应保持一致")

	n1, err := strconv.ParseUint(p1[1], 36, 64)
placeholder
	n2, err := strconv.ParseUint(p2[1], 36, 64)
placeholder
	require.Equal(t, n1+1, n2, "计数器应单调递增")
placeholder

func TestGetAccountsLoadBatch_ReturnsCorrectData(t *testing.T) {
	expected := map[int64]*AccountLoadInfo{
		1: {AccountID: 1, CurrentConcurrency: 3, WaitingCount: 0, LoadRate: 60placeholder,
		2: {AccountID: 2, CurrentConcurrency: 5, WaitingCount: 2, LoadRate: 100placeholder,
placeholder
	cache := &stubConcurrencyCacheForTest{loadBatch: expectedplaceholder
	svc := NewConcurrencyService(cache)

	accounts := []AccountWithConcurrency{
		{ID: 1, MaxConcurrency: 5placeholder,
		{ID: 2, MaxConcurrency: 5placeholder,
placeholder
	result, err := svc.GetAccountsLoadBatch(context.Background(), accounts)
placeholder
	require.Equal(t, expected, result)
placeholder

func TestGetAccountsLoadBatch_NilCache(t *testing.T) {
	svc := &ConcurrencyService{cache: nilplaceholder

	result, err := svc.GetAccountsLoadBatch(context.Background(), nil)
placeholder
	require.Empty(t, result)
placeholder

func TestIncrementWaitCount_Success(t *testing.T) {
	cache := &stubConcurrencyCacheForTest{waitAllowed: trueplaceholder
	svc := NewConcurrencyService(cache)

	allowed, err := svc.IncrementWaitCount(context.Background(), 1, 25)
placeholder
	require.True(t, allowed)
placeholder

func TestIncrementWaitCount_QueueFull(t *testing.T) {
	cache := &stubConcurrencyCacheForTest{waitAllowed: falseplaceholder
	svc := NewConcurrencyService(cache)

	allowed, err := svc.IncrementWaitCount(context.Background(), 1, 25)
placeholder
	require.False(t, allowed)
placeholder

func TestIncrementWaitCount_FailOpen(t *testing.T) {
	// Redis 错误时应 fail-open（允许请求通过）
	cache := &stubConcurrencyCacheForTest{waitErr: errors.New("redis timeout")placeholder
	svc := NewConcurrencyService(cache)

	allowed, err := svc.IncrementWaitCount(context.Background(), 1, 25)
	require.NoError(t, err, "Redis 错误不应传播")
	require.True(t, allowed, "Redis 错误时应 fail-open")
placeholder

func TestIncrementWaitCount_NilCache(t *testing.T) {
	svc := &ConcurrencyService{cache: nilplaceholder

	allowed, err := svc.IncrementWaitCount(context.Background(), 1, 25)
placeholder
	require.True(t, allowed, "nil cache 应 fail-open")
placeholder

func TestCalculateMaxWait(t *testing.T) {
	tests := []struct {
		concurrency int
		expected    int
placeholder{
		{5, 25placeholder,  // 5 + 20
		{1, 21placeholder,  // 1 + 20
		{0, 21placeholder,  // min(1) + 20
		{-1, 21placeholder, // min(1) + 20
		{10, 30placeholder, // 10 + 20
placeholder
	for _, tt := range tests {
		result := CalculateMaxWait(tt.concurrency)
		require.Equal(t, tt.expected, result, "CalculateMaxWait(%d)", tt.concurrency)
placeholder
placeholder

func TestGetAccountWaitingCount(t *testing.T) {
	cache := &stubConcurrencyCacheForTest{waitCount: 5placeholder
	svc := NewConcurrencyService(cache)

	count, err := svc.GetAccountWaitingCount(context.Background(), 1)
placeholder
	require.Equal(t, 5, count)
placeholder

func TestGetAccountWaitingCount_NilCache(t *testing.T) {
	svc := &ConcurrencyService{cache: nilplaceholder

	count, err := svc.GetAccountWaitingCount(context.Background(), 1)
placeholder
	require.Equal(t, 0, count)
placeholder

func TestGetAccountConcurrencyBatch(t *testing.T) {
	cache := &stubConcurrencyCacheForTest{concurrency: 3placeholder
	svc := NewConcurrencyService(cache)

	result, err := svc.GetAccountConcurrencyBatch(context.Background(), []int64{1, 2, 3placeholder)
placeholder
	require.Len(t, result, 3)
	for _, id := range []int64{1, 2, 3placeholder {
		require.Equal(t, 3, result[id])
placeholder
placeholder

func TestIncrementAccountWaitCount_FailOpen(t *testing.T) {
	cache := &stubConcurrencyCacheForTest{waitErr: errors.New("redis error")placeholder
	svc := NewConcurrencyService(cache)

	allowed, err := svc.IncrementAccountWaitCount(context.Background(), 1, 10)
	require.NoError(t, err, "Redis 错误不应传播")
	require.True(t, allowed, "Redis 错误时应 fail-open")
placeholder

func TestIncrementAccountWaitCount_NilCache(t *testing.T) {
	svc := &ConcurrencyService{cache: nilplaceholder

	allowed, err := svc.IncrementAccountWaitCount(context.Background(), 1, 10)
placeholder
	require.True(t, allowed)
placeholder
