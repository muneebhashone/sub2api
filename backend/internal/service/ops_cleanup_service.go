package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

const (
	opsCleanupJobName = "ops_cleanup"

	opsCleanupLeaderLockKeyDefault = "ops:cleanup:leader"
	opsCleanupLeaderLockTTLDefault = 30 * time.Minute
)

var opsCleanupCronParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

var opsCleanupReleaseScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("DEL", KEYS[1])
end
return 0
`)

// OpsCleanupService periodically deletes old ops data to prevent unbounded DB growth.
//
// - Scheduling: 5-field cron spec (minute hour dom month dow).
// - Multi-instance: best-effort Redis leader lock so only one node runs cleanup.
// - Safety: deletes in batches to avoid long transactions.
type OpsCleanupService struct {
	opsRepo     OpsRepository
	db          *sql.DB
	redisClient *redis.Client
	cfg         *config.Config

	instanceID string

	cron *cron.Cron

	startOnce sync.Once
	stopOnce  sync.Once

	warnNoRedisOnce sync.Once
placeholder

func NewOpsCleanupService(
	opsRepo OpsRepository,
	db *sql.DB,
	redisClient *redis.Client,
	cfg *config.Config,
) *OpsCleanupService {
	return &OpsCleanupService{
		opsRepo:     opsRepo,
		db:          db,
		redisClient: redisClient,
		cfg:         cfg,
		instanceID:  uuid.NewString(),
placeholder
placeholder

func (s *OpsCleanupService) Start() {
	if s == nil {
		return
placeholder
	if s.cfg != nil && !s.cfg.Ops.Enabled {
		return
placeholder
	if s.cfg != nil && !s.cfg.Ops.Cleanup.Enabled {
		log.Printf("[OpsCleanup] not started (disabled)")
		return
placeholder
	if s.opsRepo == nil || s.db == nil {
		log.Printf("[OpsCleanup] not started (missing deps)")
		return
placeholder

	s.startOnce.Do(func() {
		schedule := "0 2 * * *"
		if s.cfg != nil && strings.TrimSpace(s.cfg.Ops.Cleanup.Schedule) != "" {
			schedule = strings.TrimSpace(s.cfg.Ops.Cleanup.Schedule)
	placeholder

		loc := time.Local
		if s.cfg != nil && strings.TrimSpace(s.cfg.Timezone) != "" {
			if parsed, err := time.LoadLocation(strings.TrimSpace(s.cfg.Timezone)); err == nil && parsed != nil {
				loc = parsed
		placeholder
		placeholder
	
			c := cron.New(cron.WithParser(opsCleanupCronParser), cron.WithLocation(loc))
			_, err := c.AddFunc(schedule, func() { s.runScheduled() placeholder)
			if err != nil {
				log.Printf("[OpsCleanup] not started (invalid schedule=%q): %v", schedule, err)
				return
		placeholder
			s.cron = c
			s.cron.Start()
			log.Printf("[OpsCleanup] started (schedule=%q tz=%s)", schedule, loc.String())
	placeholder)
placeholder

func (s *OpsCleanupService) Stop() {
	if s == nil {
		return
placeholder
	s.stopOnce.Do(func() {
		if s.cron != nil {
			ctx := s.cron.Stop()
			select {
			case <-ctx.Done():
			case <-time.After(3 * time.Second):
				log.Printf("[OpsCleanup] cron stop timed out")
		placeholder
	placeholder
placeholder)
placeholder

func (s *OpsCleanupService) runScheduled() {
	if s == nil || s.db == nil || s.opsRepo == nil {
		return
placeholder

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	release, ok := s.tryAcquireLeaderLock(ctx)
	if !ok {
		return
placeholder
	if release != nil {
		defer release()
placeholder

	startedAt := time.Now().UTC()
	runAt := startedAt

	counts, err := s.runCleanupOnce(ctx)
	if err != nil {
		s.recordHeartbeatError(runAt, time.Since(startedAt), err)
		log.Printf("[OpsCleanup] cleanup failed: %v", err)
		return
placeholder
	s.recordHeartbeatSuccess(runAt, time.Since(startedAt))
	log.Printf("[OpsCleanup] cleanup complete: %s", counts)
placeholder

type opsCleanupDeletedCounts struct {
	errorLogs     int64
	retryAttempts int64
	alertEvents   int64
	systemMetrics int64
	hourlyPreagg  int64
	dailyPreagg   int64
placeholder

func (c opsCleanupDeletedCounts) String() string {
	return fmt.Sprintf(
		"error_logs=%d retry_attempts=%d alert_events=%d system_metrics=%d hourly_preagg=%d daily_preagg=%d",
		c.errorLogs,
		c.retryAttempts,
		c.alertEvents,
		c.systemMetrics,
		c.hourlyPreagg,
		c.dailyPreagg,
	)
placeholder

func (s *OpsCleanupService) runCleanupOnce(ctx context.Context) (opsCleanupDeletedCounts, error) {
	out := opsCleanupDeletedCounts{placeholder
	if s == nil || s.db == nil || s.cfg == nil {
		return out, nil
placeholder

	batchSize := 5000

	now := time.Now().UTC()

	// Error-like tables: error logs / retry attempts / alert events.
	if days := s.cfg.Ops.Cleanup.ErrorLogRetentionDays; days > 0 {
		cutoff := now.AddDate(0, 0, -days)
		n, err := deleteOldRowsByID(ctx, s.db, "ops_error_logs", "created_at", cutoff, batchSize, false)
		if err != nil {
			return out, err
	placeholder
		out.errorLogs = n

		n, err = deleteOldRowsByID(ctx, s.db, "ops_retry_attempts", "created_at", cutoff, batchSize, false)
		if err != nil {
			return out, err
	placeholder
		out.retryAttempts = n

		n, err = deleteOldRowsByID(ctx, s.db, "ops_alert_events", "created_at", cutoff, batchSize, false)
		if err != nil {
			return out, err
	placeholder
		out.alertEvents = n
placeholder

	// Minute-level metrics snapshots.
	if days := s.cfg.Ops.Cleanup.MinuteMetricsRetentionDays; days > 0 {
		cutoff := now.AddDate(0, 0, -days)
		n, err := deleteOldRowsByID(ctx, s.db, "ops_system_metrics", "created_at", cutoff, batchSize, false)
		if err != nil {
			return out, err
	placeholder
		out.systemMetrics = n
placeholder

	// Pre-aggregation tables (hourly/daily).
	if days := s.cfg.Ops.Cleanup.HourlyMetricsRetentionDays; days > 0 {
		cutoff := now.AddDate(0, 0, -days)
		n, err := deleteOldRowsByID(ctx, s.db, "ops_metrics_hourly", "bucket_start", cutoff, batchSize, false)
		if err != nil {
			return out, err
	placeholder
		out.hourlyPreagg = n

		n, err = deleteOldRowsByID(ctx, s.db, "ops_metrics_daily", "bucket_date", cutoff, batchSize, true)
		if err != nil {
			return out, err
	placeholder
		out.dailyPreagg = n
placeholder

	return out, nil
placeholder

func deleteOldRowsByID(
	ctx context.Context,
	db *sql.DB,
	table string,
	timeColumn string,
	cutoff time.Time,
	batchSize int,
	castCutoffToDate bool,
) (int64, error) {
	if db == nil {
		return 0, nil
placeholder
	if batchSize <= 0 {
		batchSize = 5000
placeholder

	where := fmt.Sprintf("%s < $1", timeColumn)
	if castCutoffToDate {
		where = fmt.Sprintf("%s < $1::date", timeColumn)
placeholder

	q := fmt.Sprintf(`
WITH batch AS (
  SELECT id FROM %s
  WHERE %s
  ORDER BY id
  LIMIT $2
)
DELETE FROM %s
WHERE id IN (SELECT id FROM batch)
`, table, where, table)

	var total int64
	for {
		res, err := db.ExecContext(ctx, q, cutoff, batchSize)
		if err != nil {
			// If ops tables aren't present yet (partial deployments), treat as no-op.
			if strings.Contains(strings.ToLower(err.Error()), "does not exist") && strings.Contains(strings.ToLower(err.Error()), "relation") {
				return total, nil
		placeholder
			return total, err
	placeholder
		affected, err := res.RowsAffected()
		if err != nil {
			return total, err
	placeholder
		total += affected
		if affected == 0 {
			break
	placeholder
placeholder
	return total, nil
placeholder

func (s *OpsCleanupService) tryAcquireLeaderLock(ctx context.Context) (func(), bool) {
	if s == nil {
		return nil, false
placeholder
	// In simple run mode, assume single instance.
	if s.cfg != nil && s.cfg.RunMode == config.RunModeSimple {
		return nil, true
placeholder

	key := opsCleanupLeaderLockKeyDefault
	ttl := opsCleanupLeaderLockTTLDefault

	// Prefer Redis leader lock when available, but avoid stampeding the DB when Redis is flaky by
	// falling back to a DB advisory lock.
	if s.redisClient != nil {
		ok, err := s.redisClient.SetNX(ctx, key, s.instanceID, ttl).Result()
		if err == nil {
			if !ok {
				return nil, false
		placeholder
			return func() {
				_, _ = opsCleanupReleaseScript.Run(ctx, s.redisClient, []string{keyplaceholder, s.instanceID).Result()
		placeholder, true
	placeholder
		// Redis error: fall back to DB advisory lock.
		s.warnNoRedisOnce.Do(func() {
			log.Printf("[OpsCleanup] leader lock SetNX failed; falling back to DB advisory lock: %v", err)
	placeholder)
placeholder else {
		s.warnNoRedisOnce.Do(func() {
			log.Printf("[OpsCleanup] redis not configured; using DB advisory lock")
	placeholder)
placeholder

	release, ok := tryAcquireDBAdvisoryLock(ctx, s.db, hashAdvisoryLockID(key))
	if !ok {
		return nil, false
placeholder
	return release, true
placeholder

func (s *OpsCleanupService) recordHeartbeatSuccess(runAt time.Time, duration time.Duration) {
	if s == nil || s.opsRepo == nil {
		return
placeholder
	now := time.Now().UTC()
	durMs := duration.Milliseconds()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = s.opsRepo.UpsertJobHeartbeat(ctx, &OpsUpsertJobHeartbeatInput{
		JobName:        opsCleanupJobName,
		LastRunAt:      &runAt,
		LastSuccessAt:  &now,
		LastDurationMs: &durMs,
placeholder)
placeholder

func (s *OpsCleanupService) recordHeartbeatError(runAt time.Time, duration time.Duration, err error) {
	if s == nil || s.opsRepo == nil || err == nil {
		return
placeholder
	now := time.Now().UTC()
	durMs := duration.Milliseconds()
	msg := truncateString(err.Error(), 2048)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = s.opsRepo.UpsertJobHeartbeat(ctx, &OpsUpsertJobHeartbeatInput{
		JobName:        opsCleanupJobName,
		LastRunAt:      &runAt,
		LastErrorAt:    &now,
		LastError:      &msg,
		LastDurationMs: &durMs,
placeholder)
placeholder
