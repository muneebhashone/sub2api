package service

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/usagestats"
	"github.com/stretchr/testify/require"
)

type usageRepoStub struct {
	UsageLogRepository
	stats      *usagestats.DashboardStats
	rangeStats *usagestats.DashboardStats
	err        error
	rangeErr   error
	calls      int32
	rangeCalls int32
	rangeStart time.Time
	rangeEnd   time.Time
	onCall     chan struct{placeholder
placeholder

func (s *usageRepoStub) GetDashboardStats(ctx context.Context) (*usagestats.DashboardStats, error) {
	atomic.AddInt32(&s.calls, 1)
	if s.onCall != nil {
		select {
		case s.onCall <- struct{placeholder{placeholder:
		default:
	placeholder
placeholder
	if s.err != nil {
		return nil, s.err
placeholder
	return s.stats, nil
placeholder

func (s *usageRepoStub) GetDashboardStatsWithRange(ctx context.Context, start, end time.Time) (*usagestats.DashboardStats, error) {
	atomic.AddInt32(&s.rangeCalls, 1)
	s.rangeStart = start
	s.rangeEnd = end
	if s.rangeErr != nil {
		return nil, s.rangeErr
placeholder
	if s.rangeStats != nil {
		return s.rangeStats, nil
placeholder
	return s.stats, nil
placeholder

type dashboardCacheStub struct {
	get       func(ctx context.Context) (string, error)
	set       func(ctx context.Context, data string, ttl time.Duration) error
	del       func(ctx context.Context) error
	getCalls  int32
	setCalls  int32
	delCalls  int32
	lastSetMu sync.Mutex
	lastSet   string
placeholder

func (c *dashboardCacheStub) GetDashboardStats(ctx context.Context) (string, error) {
	atomic.AddInt32(&c.getCalls, 1)
	if c.get != nil {
		return c.get(ctx)
placeholder
	return "", ErrDashboardStatsCacheMiss
placeholder

func (c *dashboardCacheStub) SetDashboardStats(ctx context.Context, data string, ttl time.Duration) error {
	atomic.AddInt32(&c.setCalls, 1)
	c.lastSetMu.Lock()
	c.lastSet = data
	c.lastSetMu.Unlock()
	if c.set != nil {
		return c.set(ctx, data, ttl)
placeholder
	return nil
placeholder

func (c *dashboardCacheStub) DeleteDashboardStats(ctx context.Context) error {
	atomic.AddInt32(&c.delCalls, 1)
	if c.del != nil {
		return c.del(ctx)
placeholder
	return nil
placeholder

type dashboardAggregationRepoStub struct {
	watermark time.Time
	err       error
placeholder

func (s *dashboardAggregationRepoStub) AggregateRange(ctx context.Context, start, end time.Time) error {
	return nil
placeholder

func (s *dashboardAggregationRepoStub) RecomputeRange(ctx context.Context, start, end time.Time) error {
	return nil
placeholder

func (s *dashboardAggregationRepoStub) GetAggregationWatermark(ctx context.Context) (time.Time, error) {
	if s.err != nil {
		return time.Time{placeholder, s.err
placeholder
	return s.watermark, nil
placeholder

func (s *dashboardAggregationRepoStub) UpdateAggregationWatermark(ctx context.Context, aggregatedAt time.Time) error {
	return nil
placeholder

func (s *dashboardAggregationRepoStub) CleanupAggregates(ctx context.Context, hourlyCutoff, dailyCutoff time.Time) error {
	return nil
placeholder

func (s *dashboardAggregationRepoStub) CleanupUsageLogs(ctx context.Context, cutoff time.Time) error {
	return nil
placeholder

func (s *dashboardAggregationRepoStub) EnsureUsageLogsPartitions(ctx context.Context, now time.Time) error {
	return nil
placeholder

func (c *dashboardCacheStub) readLastEntry(t *testing.T) dashboardStatsCacheEntry {
placeholder
	c.lastSetMu.Lock()
	data := c.lastSet
	c.lastSetMu.Unlock()

	var entry dashboardStatsCacheEntry
	err := json.Unmarshal([]byte(data), &entry)
placeholder
	return entry
placeholder

func TestDashboardService_CacheHitFresh(t *testing.T) {
	stats := &usagestats.DashboardStats{
		TotalUsers:     10,
		StatsUpdatedAt: time.Unix(0, 0).UTC().Format(time.RFC3339),
		StatsStale:     true,
placeholder
	entry := dashboardStatsCacheEntry{
		Stats:     stats,
		UpdatedAt: time.Now().Unix(),
placeholder
	payload, err := json.Marshal(entry)
placeholder

	cache := &dashboardCacheStub{
		get: func(ctx context.Context) (string, error) {
			return string(payload), nil
	placeholder,
placeholder
	repo := &usageRepoStub{
		stats: &usagestats.DashboardStats{TotalUsers: 99placeholder,
placeholder
	aggRepo := &dashboardAggregationRepoStub{watermark: time.Unix(0, 0).UTC()placeholder
	cfg := &config.Config{
		Dashboard: config.DashboardCacheConfig{Enabled: trueplaceholder,
		DashboardAgg: config.DashboardAggregationConfig{
			Enabled: true,
	placeholder,
placeholder
	svc := NewDashboardService(repo, aggRepo, cache, cfg)

	got, err := svc.GetDashboardStats(context.Background())
placeholder
	require.Equal(t, stats, got)
	require.Equal(t, int32(0), atomic.LoadInt32(&repo.calls))
	require.Equal(t, int32(1), atomic.LoadInt32(&cache.getCalls))
	require.Equal(t, int32(0), atomic.LoadInt32(&cache.setCalls))
placeholder

func TestDashboardService_CacheMiss_StoresCache(t *testing.T) {
	stats := &usagestats.DashboardStats{
		TotalUsers:     7,
		StatsUpdatedAt: time.Unix(0, 0).UTC().Format(time.RFC3339),
		StatsStale:     true,
placeholder
	cache := &dashboardCacheStub{
		get: func(ctx context.Context) (string, error) {
			return "", ErrDashboardStatsCacheMiss
	placeholder,
placeholder
	repo := &usageRepoStub{stats: statsplaceholder
	aggRepo := &dashboardAggregationRepoStub{watermark: time.Unix(0, 0).UTC()placeholder
	cfg := &config.Config{
		Dashboard: config.DashboardCacheConfig{Enabled: trueplaceholder,
		DashboardAgg: config.DashboardAggregationConfig{
			Enabled: true,
	placeholder,
placeholder
	svc := NewDashboardService(repo, aggRepo, cache, cfg)

	got, err := svc.GetDashboardStats(context.Background())
placeholder
	require.Equal(t, stats, got)
	require.Equal(t, int32(1), atomic.LoadInt32(&repo.calls))
	require.Equal(t, int32(1), atomic.LoadInt32(&cache.getCalls))
	require.Equal(t, int32(1), atomic.LoadInt32(&cache.setCalls))
	entry := cache.readLastEntry(t)
	require.Equal(t, stats, entry.Stats)
	require.WithinDuration(t, time.Now(), time.Unix(entry.UpdatedAt, 0), time.Second)
placeholder

func TestDashboardService_CacheDisabled_SkipsCache(t *testing.T) {
	stats := &usagestats.DashboardStats{
		TotalUsers:     3,
		StatsUpdatedAt: time.Unix(0, 0).UTC().Format(time.RFC3339),
		StatsStale:     true,
placeholder
	cache := &dashboardCacheStub{
		get: func(ctx context.Context) (string, error) {
			return "", nil
	placeholder,
placeholder
	repo := &usageRepoStub{stats: statsplaceholder
	aggRepo := &dashboardAggregationRepoStub{watermark: time.Unix(0, 0).UTC()placeholder
	cfg := &config.Config{
		Dashboard: config.DashboardCacheConfig{Enabled: falseplaceholder,
		DashboardAgg: config.DashboardAggregationConfig{
			Enabled: true,
	placeholder,
placeholder
	svc := NewDashboardService(repo, aggRepo, cache, cfg)

	got, err := svc.GetDashboardStats(context.Background())
placeholder
	require.Equal(t, stats, got)
	require.Equal(t, int32(1), atomic.LoadInt32(&repo.calls))
	require.Equal(t, int32(0), atomic.LoadInt32(&cache.getCalls))
	require.Equal(t, int32(0), atomic.LoadInt32(&cache.setCalls))
placeholder

func TestDashboardService_CacheHitStale_TriggersAsyncRefresh(t *testing.T) {
	staleStats := &usagestats.DashboardStats{
		TotalUsers:     11,
		StatsUpdatedAt: time.Unix(0, 0).UTC().Format(time.RFC3339),
		StatsStale:     true,
placeholder
	entry := dashboardStatsCacheEntry{
		Stats:     staleStats,
		UpdatedAt: time.Now().Add(-defaultDashboardStatsFreshTTL * 2).Unix(),
placeholder
	payload, err := json.Marshal(entry)
placeholder

	cache := &dashboardCacheStub{
		get: func(ctx context.Context) (string, error) {
			return string(payload), nil
	placeholder,
placeholder
	refreshCh := make(chan struct{placeholder, 1)
	repo := &usageRepoStub{
		stats:  &usagestats.DashboardStats{TotalUsers: 22placeholder,
		onCall: refreshCh,
placeholder
	aggRepo := &dashboardAggregationRepoStub{watermark: time.Unix(0, 0).UTC()placeholder
	cfg := &config.Config{
		Dashboard: config.DashboardCacheConfig{Enabled: trueplaceholder,
		DashboardAgg: config.DashboardAggregationConfig{
			Enabled: true,
	placeholder,
placeholder
	svc := NewDashboardService(repo, aggRepo, cache, cfg)

	got, err := svc.GetDashboardStats(context.Background())
placeholder
	require.Equal(t, staleStats, got)

	select {
	case <-refreshCh:
	case <-time.After(1 * time.Second):
		t.Fatal("等待异步刷新超时")
placeholder
	require.Eventually(t, func() bool {
		return atomic.LoadInt32(&cache.setCalls) >= 1
placeholder, 1*time.Second, 10*time.Millisecond)
placeholder

func TestDashboardService_CacheParseError_EvictsAndRefetches(t *testing.T) {
	cache := &dashboardCacheStub{
		get: func(ctx context.Context) (string, error) {
			return "not-json", nil
	placeholder,
placeholder
	stats := &usagestats.DashboardStats{TotalUsers: 9placeholder
	repo := &usageRepoStub{stats: statsplaceholder
	aggRepo := &dashboardAggregationRepoStub{watermark: time.Unix(0, 0).UTC()placeholder
	cfg := &config.Config{
		Dashboard: config.DashboardCacheConfig{Enabled: trueplaceholder,
		DashboardAgg: config.DashboardAggregationConfig{
			Enabled: true,
	placeholder,
placeholder
	svc := NewDashboardService(repo, aggRepo, cache, cfg)

	got, err := svc.GetDashboardStats(context.Background())
placeholder
	require.Equal(t, stats, got)
	require.Equal(t, int32(1), atomic.LoadInt32(&cache.delCalls))
	require.Equal(t, int32(1), atomic.LoadInt32(&repo.calls))
placeholder

func TestDashboardService_CacheParseError_RepoFailure(t *testing.T) {
	cache := &dashboardCacheStub{
		get: func(ctx context.Context) (string, error) {
			return "not-json", nil
	placeholder,
placeholder
	repo := &usageRepoStub{err: errors.New("db down")placeholder
	aggRepo := &dashboardAggregationRepoStub{watermark: time.Unix(0, 0).UTC()placeholder
	cfg := &config.Config{
		Dashboard: config.DashboardCacheConfig{Enabled: trueplaceholder,
		DashboardAgg: config.DashboardAggregationConfig{
			Enabled: true,
	placeholder,
placeholder
	svc := NewDashboardService(repo, aggRepo, cache, cfg)

	_, err := svc.GetDashboardStats(context.Background())
placeholder
	require.Equal(t, int32(1), atomic.LoadInt32(&cache.delCalls))
placeholder

func TestDashboardService_StatsUpdatedAtEpochWhenMissing(t *testing.T) {
	stats := &usagestats.DashboardStats{placeholder
	repo := &usageRepoStub{stats: statsplaceholder
	aggRepo := &dashboardAggregationRepoStub{watermark: time.Unix(0, 0).UTC()placeholder
	cfg := &config.Config{Dashboard: config.DashboardCacheConfig{Enabled: falseplaceholderplaceholder
	svc := NewDashboardService(repo, aggRepo, nil, cfg)

	got, err := svc.GetDashboardStats(context.Background())
placeholder
	require.Equal(t, "1970-01-01T00:00:00Z", got.StatsUpdatedAt)
	require.True(t, got.StatsStale)
placeholder

func TestDashboardService_StatsStaleFalseWhenFresh(t *testing.T) {
	aggNow := time.Now().UTC().Truncate(time.Second)
	stats := &usagestats.DashboardStats{placeholder
	repo := &usageRepoStub{stats: statsplaceholder
	aggRepo := &dashboardAggregationRepoStub{watermark: aggNowplaceholder
	cfg := &config.Config{
		Dashboard: config.DashboardCacheConfig{Enabled: falseplaceholder,
		DashboardAgg: config.DashboardAggregationConfig{
			Enabled:         true,
			IntervalSeconds: 60,
			LookbackSeconds: 120,
	placeholder,
placeholder
	svc := NewDashboardService(repo, aggRepo, nil, cfg)

	got, err := svc.GetDashboardStats(context.Background())
placeholder
	require.Equal(t, aggNow.Format(time.RFC3339), got.StatsUpdatedAt)
	require.False(t, got.StatsStale)
placeholder

func TestDashboardService_AggDisabled_UsesUsageLogsFallback(t *testing.T) {
	expected := &usagestats.DashboardStats{TotalUsers: 42placeholder
	repo := &usageRepoStub{
		rangeStats: expected,
		err:        errors.New("should not call aggregated stats"),
placeholder
	cfg := &config.Config{
		Dashboard: config.DashboardCacheConfig{Enabled: falseplaceholder,
		DashboardAgg: config.DashboardAggregationConfig{
			Enabled: false,
			Retention: config.DashboardAggregationRetentionConfig{
				UsageLogsDays: 7,
		placeholder,
	placeholder,
placeholder
	svc := NewDashboardService(repo, nil, nil, cfg)

	got, err := svc.GetDashboardStats(context.Background())
placeholder
	require.Equal(t, int64(42), got.TotalUsers)
	require.Equal(t, int32(0), atomic.LoadInt32(&repo.calls))
	require.Equal(t, int32(1), atomic.LoadInt32(&repo.rangeCalls))
	require.False(t, repo.rangeEnd.IsZero())
	require.Equal(t, truncateToDayUTC(repo.rangeEnd.AddDate(0, 0, -7)), repo.rangeStart)
placeholder
