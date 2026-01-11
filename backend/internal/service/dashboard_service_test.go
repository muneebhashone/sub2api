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
	stats  *usagestats.DashboardStats
	err    error
	calls  int32
	onCall chan struct{placeholder
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
		TotalUsers: 10,
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
	cfg := &config.Config{Dashboard: config.DashboardCacheConfig{Enabled: trueplaceholderplaceholder
	svc := NewDashboardService(repo, cache, cfg)

	got, err := svc.GetDashboardStats(context.Background())
placeholder
	require.Equal(t, stats, got)
	require.Equal(t, int32(0), atomic.LoadInt32(&repo.calls))
	require.Equal(t, int32(1), atomic.LoadInt32(&cache.getCalls))
	require.Equal(t, int32(0), atomic.LoadInt32(&cache.setCalls))
placeholder

func TestDashboardService_CacheMiss_StoresCache(t *testing.T) {
	stats := &usagestats.DashboardStats{
		TotalUsers: 7,
placeholder
	cache := &dashboardCacheStub{
		get: func(ctx context.Context) (string, error) {
			return "", ErrDashboardStatsCacheMiss
	placeholder,
placeholder
	repo := &usageRepoStub{stats: statsplaceholder
	cfg := &config.Config{Dashboard: config.DashboardCacheConfig{Enabled: trueplaceholderplaceholder
	svc := NewDashboardService(repo, cache, cfg)

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
		TotalUsers: 3,
placeholder
	cache := &dashboardCacheStub{
		get: func(ctx context.Context) (string, error) {
			return "", nil
	placeholder,
placeholder
	repo := &usageRepoStub{stats: statsplaceholder
	cfg := &config.Config{Dashboard: config.DashboardCacheConfig{Enabled: falseplaceholderplaceholder
	svc := NewDashboardService(repo, cache, cfg)

	got, err := svc.GetDashboardStats(context.Background())
placeholder
	require.Equal(t, stats, got)
	require.Equal(t, int32(1), atomic.LoadInt32(&repo.calls))
	require.Equal(t, int32(0), atomic.LoadInt32(&cache.getCalls))
	require.Equal(t, int32(0), atomic.LoadInt32(&cache.setCalls))
placeholder

func TestDashboardService_CacheHitStale_TriggersAsyncRefresh(t *testing.T) {
	staleStats := &usagestats.DashboardStats{
		TotalUsers: 11,
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
	cfg := &config.Config{Dashboard: config.DashboardCacheConfig{Enabled: trueplaceholderplaceholder
	svc := NewDashboardService(repo, cache, cfg)

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
	cfg := &config.Config{Dashboard: config.DashboardCacheConfig{Enabled: trueplaceholderplaceholder
	svc := NewDashboardService(repo, cache, cfg)

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
	cfg := &config.Config{Dashboard: config.DashboardCacheConfig{Enabled: trueplaceholderplaceholder
	svc := NewDashboardService(repo, cache, cfg)

	_, err := svc.GetDashboardStats(context.Background())
placeholder
	require.Equal(t, int32(1), atomic.LoadInt32(&cache.delCalls))
placeholder
