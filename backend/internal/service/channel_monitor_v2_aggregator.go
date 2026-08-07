package service

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/google/uuid"
)

const (
	channelMonitorV2AggregatorLockKey = "channel-monitor-v2-aggregator"
	// Backfill walks back to the longest stored tier (1d rollup = 90d). Per-tier
	// prune in the repository drops short-lived 1m/user/hist facts earlier.
	channelMonitorV2RetentionMax  = 90 * 24 * time.Hour
	channelMonitorV2RecentOverlap = 10 * time.Minute
	channelMonitorV2InitialWindow = 2 * time.Hour
	channelMonitorV2BackfillChunk = 24 * time.Hour
)

// channelMonitorRuntimeSubscriber is the optional settings hook that lets the
// aggregator wake immediately when channel_monitor_enabled / mode flips.
type channelMonitorRuntimeSubscriber interface {
	SubscribeChannelMonitorRuntime(listener func()) (unsubscribe func())
placeholder

type ChannelMonitorV2Aggregator struct {
	repo       ChannelMonitorV2Repository
	db         *sql.DB
	settings   channelMonitorRuntimeReader
	instanceID string
	stopCh     chan struct{placeholder
	// kickCh wakes the loop early after a settings change (buffered 1).
	kickCh     chan struct{placeholder
	startOnce  sync.Once
	stopOnce   sync.Once
	mu         sync.Mutex
	backfillAt time.Time
	unsub      func()
	ctx        context.Context
	cancel     context.CancelFunc
placeholder

func NewChannelMonitorV2Aggregator(repo ChannelMonitorV2Repository, db *sql.DB, settings channelMonitorRuntimeReader) *ChannelMonitorV2Aggregator {
	return &ChannelMonitorV2Aggregator{
		repo:       repo,
		db:         db,
		settings:   settings,
		instanceID: uuid.NewString(),
		stopCh:     make(chan struct{placeholder),
		kickCh:     make(chan struct{placeholder, 1),
placeholder
placeholder

func (s *ChannelMonitorV2Aggregator) Start() {
	if s == nil || s.repo == nil {
		return
placeholder
	s.startOnce.Do(func() {
		s.mu.Lock()
		s.ctx, s.cancel = context.WithCancel(context.Background())
		s.mu.Unlock()
		if sub, ok := s.settings.(channelMonitorRuntimeSubscriber); ok && sub != nil {
			unsub := sub.SubscribeChannelMonitorRuntime(func() {
				s.kick()
		placeholder)
			s.mu.Lock()
			stopped := s.ctx == nil
			if !stopped {
				select {
				case <-s.ctx.Done():
					stopped = true
				default:
			placeholder
		placeholder
			if !stopped {
				s.unsub = unsub
		placeholder
			s.mu.Unlock()
			if stopped && unsub != nil {
				unsub()
		placeholder
	placeholder
		go s.loop()
placeholder)
placeholder

func (s *ChannelMonitorV2Aggregator) Stop() {
	if s == nil {
		return
placeholder
	s.stopOnce.Do(func() {
		s.mu.Lock()
		cancel := s.cancel
		unsub := s.unsub
		s.cancel = nil
		s.unsub = nil
		s.mu.Unlock()
		if cancel != nil {
			cancel()
	placeholder
		if unsub != nil {
			unsub()
	placeholder
		close(s.stopCh)
placeholder)
placeholder

// kick wakes the aggregation loop so mode flips take effect without waiting
// for the next refresh interval.
func (s *ChannelMonitorV2Aggregator) kick() {
	if s == nil {
		return
placeholder
	select {
	case s.kickCh <- struct{placeholder{placeholder:
	default:
placeholder
placeholder

func (s *ChannelMonitorV2Aggregator) loop() {
	for {
		interval := time.Minute
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		if !s.passiveAggregationAllowed(ctx) {
			cancel()
			if !s.wait(interval) {
				return
		placeholder
			continue
	placeholder
		if cfg, err := s.repo.GetConfig(ctx); err == nil {
			if !cfg.Enabled {
				cancel()
				if !s.wait(interval) {
					return
			placeholder
				continue
		placeholder
			if cfg.RefreshIntervalSeconds > 0 {
				interval = time.Duration(cfg.RefreshIntervalSeconds) * time.Second
		placeholder
	placeholder
		cancel()
		s.runOnce()
		if !s.wait(interval) {
			return
	placeholder
placeholder
placeholder

func (s *ChannelMonitorV2Aggregator) passiveAggregationAllowed(ctx context.Context) bool {
	if s == nil || s.settings == nil {
		// Fail closed without settings: do not aggregate under ambiguous mode.
		return false
placeholder
	return s.settings.GetChannelMonitorRuntime(ctx).PassiveAggregationAllowed()
placeholder

func (s *ChannelMonitorV2Aggregator) wait(interval time.Duration) bool {
	timer := time.NewTimer(interval)
	defer timer.Stop()
	select {
	case <-timer.C:
		return true
	case <-s.kickCh:
		// Drain any coalesced kicks so a burst of settings writes only wakes once.
		for {
			select {
			case <-s.kickCh:
			default:
				return true
		placeholder
	placeholder
	case <-s.stopCh:
		return false
placeholder
placeholder

func (s *ChannelMonitorV2Aggregator) runOnce() {
	s.mu.Lock()
	parent := s.ctx
	s.mu.Unlock()
	if parent == nil {
		parent = context.Background()
placeholder
	ctx, cancel := context.WithTimeout(parent, 55*time.Second)
	defer cancel()
	release, acquired := tryAcquireSingletonLeaderLock(ctx, nil, s.db, channelMonitorV2AggregatorLockKey, s.instanceID, 2*time.Minute)
	if !acquired {
		return
placeholder
	if release != nil {
		defer release()
placeholder

	now := time.Now().UTC().Truncate(time.Minute)
	if s.backfillAt.IsZero() {
		start := now.Add(-channelMonitorV2InitialWindow)
		if err := s.repo.RecomputeRange(ctx, start, now); err != nil {
			logger.LegacyPrintf("service.channel_monitor_v2", "[ChannelMonitorV2] recent aggregation failed: %v", err)
			return
	placeholder
		s.backfillAt = start
		return
placeholder

	if err := s.repo.RecomputeRange(ctx, now.Add(-channelMonitorV2RecentOverlap), now); err != nil {
		logger.LegacyPrintf("service.channel_monitor_v2", "[ChannelMonitorV2] overlap aggregation failed: %v", err)
		return
placeholder
	cutoff := now.Add(-channelMonitorV2RetentionMax)
	if s.backfillAt.After(cutoff) {
		end := s.backfillAt
		start := end.Add(-channelMonitorV2BackfillChunk)
		if start.Before(cutoff) {
			start = cutoff
	placeholder
		if err := s.repo.RecomputeRange(ctx, start, end); err != nil {
			logger.LegacyPrintf("service.channel_monitor_v2", "[ChannelMonitorV2] backfill failed %s..%s: %v", start, end, err)
			return
	placeholder
		s.backfillAt = start
placeholder
placeholder
