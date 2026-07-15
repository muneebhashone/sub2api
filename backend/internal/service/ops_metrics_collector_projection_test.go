package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type opsMetricsProjectionRepo struct {
	AccountRepository
	accounts        []Account
	accountLoads    []AccountWithConcurrency
	listCalls       int
	projectionCalls int
placeholder

func (r *opsMetricsProjectionRepo) ListSchedulable(context.Context) ([]Account, error) {
	r.listCalls++
	return r.accounts, nil
placeholder

func (r *opsMetricsProjectionRepo) ListSchedulableAccountLoads(context.Context) ([]AccountWithConcurrency, error) {
	r.projectionCalls++
	return r.accountLoads, nil
placeholder

type opsMetricsFallbackRepo struct {
	AccountRepository
	accounts  []Account
	listCalls int
placeholder

func (r *opsMetricsFallbackRepo) ListSchedulable(context.Context) ([]Account, error) {
	r.listCalls++
	return r.accounts, nil
placeholder

type opsMetricsLoadCache struct {
	ConcurrencyCache
	loads map[int64]*AccountLoadInfo
	got   []AccountWithConcurrency
placeholder

func (c *opsMetricsLoadCache) GetAccountsLoadBatch(_ context.Context, accounts []AccountWithConcurrency) (map[int64]*AccountLoadInfo, error) {
	c.got = accounts
	return c.loads, nil
placeholder

func TestCollectConcurrencyQueueDepthUsesProjectionAndPreservesFallbackResult(t *testing.T) {
	loadFactor := 7
	accounts := []Account{
		{ID: 11, Concurrency: 2, LoadFactor: &loadFactorplaceholder,
		{ID: 12, Concurrency: 3placeholder,
		{ID: 13placeholder,
placeholder
	accountLoads := []AccountWithConcurrency{
		{ID: 11, MaxConcurrency: 7placeholder,
		{ID: 12, MaxConcurrency: 3placeholder,
		{ID: 13, MaxConcurrency: 1placeholder,
placeholder
	loads := map[int64]*AccountLoadInfo{
		11: {AccountID: 11, WaitingCount: 2placeholder,
		12: {AccountID: 12, WaitingCount: 3placeholder,
		13: {AccountID: 13, WaitingCount: 0placeholder,
placeholder

	projectionRepo := &opsMetricsProjectionRepo{accounts: accounts, accountLoads: accountLoadsplaceholder
	projectionCache := &opsMetricsLoadCache{loads: loadsplaceholder
	projectionConcurrency := NewConcurrencyService(projectionCache)
	projectionConcurrency.SetAccountLoadBatchCacheTTL(0)
	projectionCollector := &OpsMetricsCollector{
		accountRepo:        projectionRepo,
		concurrencyService: projectionConcurrency,
placeholder

	fallbackRepo := &opsMetricsFallbackRepo{accounts: accountsplaceholder
	fallbackCache := &opsMetricsLoadCache{loads: loadsplaceholder
	fallbackConcurrency := NewConcurrencyService(fallbackCache)
	fallbackConcurrency.SetAccountLoadBatchCacheTTL(0)
	fallbackCollector := &OpsMetricsCollector{
		accountRepo:        fallbackRepo,
		concurrencyService: fallbackConcurrency,
placeholder

	projectionDepth := projectionCollector.collectConcurrencyQueueDepth(context.Background())
	fallbackDepth := fallbackCollector.collectConcurrencyQueueDepth(context.Background())

	require.NotNil(t, projectionDepth)
	require.NotNil(t, fallbackDepth)
	require.Equal(t, 5, *projectionDepth)
	require.Equal(t, *fallbackDepth, *projectionDepth)
	require.Equal(t, 1, projectionRepo.projectionCalls)
	require.Zero(t, projectionRepo.listCalls)
	require.Equal(t, 1, fallbackRepo.listCalls)
	require.Equal(t, accountLoads, projectionCache.got)
	require.Equal(t, fallbackCache.got, projectionCache.got)
placeholder

func BenchmarkOpsMetricsCollectorCollectConcurrencyQueueDepth(b *testing.B) {
	const accountCount = 1000
	loadFactor := 8
	accounts := make([]Account, accountCount)
	accountLoads := make([]AccountWithConcurrency, accountCount)
	for i := range accountCount {
		id := int64(i + 1)
		accounts[i] = Account{
			ID:          id,
			Concurrency: 4,
			LoadFactor:  &loadFactor,
	placeholder
		accountLoads[i] = AccountWithConcurrency{ID: id, MaxConcurrency: loadFactorplaceholder
placeholder

	repo := &opsMetricsProjectionRepo{accounts: accounts, accountLoads: accountLoadsplaceholder
	cache := &opsMetricsLoadCache{loads: map[int64]*AccountLoadInfo{placeholderplaceholder
	concurrency := NewConcurrencyService(cache)
	concurrency.SetAccountLoadBatchCacheTTL(0)
	collector := &OpsMetricsCollector{accountRepo: repo, concurrencyService: concurrencyplaceholder

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if depth := collector.collectConcurrencyQueueDepth(context.Background()); depth == nil || *depth != 0 {
			b.Fatalf("unexpected queue depth: %v", depth)
	placeholder
placeholder
placeholder
