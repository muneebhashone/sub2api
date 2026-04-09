package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/redis/go-redis/v9"
)

const (
	schedulerBucketSetKey       = "sched:buckets"
	schedulerOutboxWatermarkKey = "sched:outbox:watermark"
	schedulerAccountPrefix      = "sched:acc:"
	schedulerAccountMetaPrefix  = "sched:meta:"
	schedulerActivePrefix       = "sched:active:"
	schedulerReadyPrefix        = "sched:ready:"
	schedulerVersionPrefix      = "sched:ver:"
	schedulerSnapshotPrefix     = "sched:"
	schedulerLockPrefix         = "sched:lock:"

	defaultSchedulerSnapshotMGetChunkSize  = 128
	defaultSchedulerSnapshotWriteChunkSize = 256
)

type schedulerCache struct {
	rdb            *redis.Client
	mgetChunkSize  int
	writeChunkSize int
placeholder

func NewSchedulerCache(rdb *redis.Client) service.SchedulerCache {
	return newSchedulerCacheWithChunkSizes(rdb, defaultSchedulerSnapshotMGetChunkSize, defaultSchedulerSnapshotWriteChunkSize)
placeholder

func newSchedulerCacheWithChunkSizes(rdb *redis.Client, mgetChunkSize, writeChunkSize int) service.SchedulerCache {
	if mgetChunkSize <= 0 {
		mgetChunkSize = defaultSchedulerSnapshotMGetChunkSize
placeholder
	if writeChunkSize <= 0 {
		writeChunkSize = defaultSchedulerSnapshotWriteChunkSize
placeholder
	return &schedulerCache{
		rdb:            rdb,
		mgetChunkSize:  mgetChunkSize,
		writeChunkSize: writeChunkSize,
placeholder
placeholder

func (c *schedulerCache) GetSnapshot(ctx context.Context, bucket service.SchedulerBucket) ([]*service.Account, bool, error) {
	readyKey := schedulerBucketKey(schedulerReadyPrefix, bucket)
	readyVal, err := c.rdb.Get(ctx, readyKey).Result()
	if err == redis.Nil {
		return nil, false, nil
placeholder
	if err != nil {
		return nil, false, err
placeholder
	if readyVal != "1" {
		return nil, false, nil
placeholder

	activeKey := schedulerBucketKey(schedulerActivePrefix, bucket)
	activeVal, err := c.rdb.Get(ctx, activeKey).Result()
	if err == redis.Nil {
		return nil, false, nil
placeholder
	if err != nil {
		return nil, false, err
placeholder

	snapshotKey := schedulerSnapshotKey(bucket, activeVal)
	ids, err := c.rdb.ZRange(ctx, snapshotKey, 0, -1).Result()
	if err != nil {
		return nil, false, err
placeholder
	if len(ids) == 0 {
		// 空快照视为缓存未命中，触发数据库回退查询
		// 这解决了新分组创建后立即绑定账号时的竞态条件问题
		return nil, false, nil
placeholder

	keys := make([]string, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, schedulerAccountMetaKey(id))
placeholder
	values, err := c.mgetChunked(ctx, keys)
	if err != nil {
		return nil, false, err
placeholder

	accounts := make([]*service.Account, 0, len(values))
	for _, val := range values {
		if val == nil {
			return nil, false, nil
	placeholder
		account, err := decodeCachedAccount(val)
		if err != nil {
			return nil, false, err
	placeholder
		accounts = append(accounts, account)
placeholder

	return accounts, true, nil
placeholder

func (c *schedulerCache) SetSnapshot(ctx context.Context, bucket service.SchedulerBucket, accounts []service.Account) error {
	activeKey := schedulerBucketKey(schedulerActivePrefix, bucket)
	oldActive, _ := c.rdb.Get(ctx, activeKey).Result()

	versionKey := schedulerBucketKey(schedulerVersionPrefix, bucket)
	version, err := c.rdb.Incr(ctx, versionKey).Result()
	if err != nil {
		return err
placeholder

	versionStr := strconv.FormatInt(version, 10)
	snapshotKey := schedulerSnapshotKey(bucket, versionStr)

	if err := c.writeAccounts(ctx, accounts); err != nil {
		return err
placeholder

	pipe := c.rdb.Pipeline()
	if len(accounts) > 0 {
		// 使用序号作为 score，保持数据库返回的排序语义。
		members := make([]redis.Z, 0, len(accounts))
		for idx, account := range accounts {
			members = append(members, redis.Z{
				Score:  float64(idx),
				Member: strconv.FormatInt(account.ID, 10),
		placeholder)
	placeholder
		for start := 0; start < len(members); start += c.writeChunkSize {
			end := start + c.writeChunkSize
			if end > len(members) {
				end = len(members)
		placeholder
			pipe.ZAdd(ctx, snapshotKey, members[start:end]...)
	placeholder
placeholder else {
		pipe.Del(ctx, snapshotKey)
placeholder
	pipe.Set(ctx, activeKey, versionStr, 0)
	pipe.Set(ctx, schedulerBucketKey(schedulerReadyPrefix, bucket), "1", 0)
	pipe.SAdd(ctx, schedulerBucketSetKey, bucket.String())
	if _, err := pipe.Exec(ctx); err != nil {
		return err
placeholder

	if oldActive != "" && oldActive != versionStr {
		_ = c.rdb.Del(ctx, schedulerSnapshotKey(bucket, oldActive)).Err()
placeholder

	return nil
placeholder

func (c *schedulerCache) GetAccount(ctx context.Context, accountID int64) (*service.Account, error) {
	key := schedulerAccountKey(strconv.FormatInt(accountID, 10))
	val, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
placeholder
	if err != nil {
		return nil, err
placeholder
	return decodeCachedAccount(val)
placeholder

func (c *schedulerCache) SetAccount(ctx context.Context, account *service.Account) error {
	if account == nil || account.ID <= 0 {
		return nil
placeholder
	return c.writeAccounts(ctx, []service.Account{*accountplaceholder)
placeholder

func (c *schedulerCache) DeleteAccount(ctx context.Context, accountID int64) error {
	if accountID <= 0 {
		return nil
placeholder
	id := strconv.FormatInt(accountID, 10)
	return c.rdb.Del(ctx, schedulerAccountKey(id), schedulerAccountMetaKey(id)).Err()
placeholder

func (c *schedulerCache) UpdateLastUsed(ctx context.Context, updates map[int64]time.Time) error {
	if len(updates) == 0 {
		return nil
placeholder

	keys := make([]string, 0, len(updates))
	ids := make([]int64, 0, len(updates))
	for id := range updates {
		keys = append(keys, schedulerAccountKey(strconv.FormatInt(id, 10)))
		ids = append(ids, id)
placeholder

	values, err := c.mgetChunked(ctx, keys)
	if err != nil {
		return err
placeholder

	pipe := c.rdb.Pipeline()
	for i, val := range values {
		if val == nil {
			continue
	placeholder
		account, err := decodeCachedAccount(val)
		if err != nil {
			return err
	placeholder
		account.LastUsedAt = ptrTime(updates[ids[i]])
		updated, err := json.Marshal(account)
		if err != nil {
			return err
	placeholder
		metaPayload, err := json.Marshal(buildSchedulerMetadataAccount(*account))
		if err != nil {
			return err
	placeholder
		pipe.Set(ctx, keys[i], updated, 0)
		pipe.Set(ctx, schedulerAccountMetaKey(strconv.FormatInt(ids[i], 10)), metaPayload, 0)
placeholder
	_, err = pipe.Exec(ctx)
	return err
placeholder

func (c *schedulerCache) TryLockBucket(ctx context.Context, bucket service.SchedulerBucket, ttl time.Duration) (bool, error) {
	key := schedulerBucketKey(schedulerLockPrefix, bucket)
	return c.rdb.SetNX(ctx, key, time.Now().UnixNano(), ttl).Result()
placeholder

func (c *schedulerCache) ListBuckets(ctx context.Context) ([]service.SchedulerBucket, error) {
	raw, err := c.rdb.SMembers(ctx, schedulerBucketSetKey).Result()
	if err != nil {
		return nil, err
placeholder
	out := make([]service.SchedulerBucket, 0, len(raw))
	for _, entry := range raw {
		bucket, ok := service.ParseSchedulerBucket(entry)
		if !ok {
			continue
	placeholder
		out = append(out, bucket)
placeholder
	return out, nil
placeholder

func (c *schedulerCache) GetOutboxWatermark(ctx context.Context) (int64, error) {
	val, err := c.rdb.Get(ctx, schedulerOutboxWatermarkKey).Result()
	if err == redis.Nil {
		return 0, nil
placeholder
	if err != nil {
		return 0, err
placeholder
	id, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
placeholder
	return id, nil
placeholder

func (c *schedulerCache) SetOutboxWatermark(ctx context.Context, id int64) error {
	return c.rdb.Set(ctx, schedulerOutboxWatermarkKey, strconv.FormatInt(id, 10), 0).Err()
placeholder

func schedulerBucketKey(prefix string, bucket service.SchedulerBucket) string {
	return fmt.Sprintf("%s%d:%s:%s", prefix, bucket.GroupID, bucket.Platform, bucket.Mode)
placeholder

func schedulerSnapshotKey(bucket service.SchedulerBucket, version string) string {
	return fmt.Sprintf("%s%d:%s:%s:v%s", schedulerSnapshotPrefix, bucket.GroupID, bucket.Platform, bucket.Mode, version)
placeholder

func schedulerAccountKey(id string) string {
	return schedulerAccountPrefix + id
placeholder

func schedulerAccountMetaKey(id string) string {
	return schedulerAccountMetaPrefix + id
placeholder

func ptrTime(t time.Time) *time.Time {
	return &t
placeholder

func decodeCachedAccount(val any) (*service.Account, error) {
	var payload []byte
	switch raw := val.(type) {
	case string:
		payload = []byte(raw)
	case []byte:
		payload = raw
	default:
		return nil, fmt.Errorf("unexpected account cache type: %T", val)
placeholder
	var account service.Account
	if err := json.Unmarshal(payload, &account); err != nil {
		return nil, err
placeholder
	return &account, nil
placeholder

func (c *schedulerCache) writeAccounts(ctx context.Context, accounts []service.Account) error {
	if len(accounts) == 0 {
		return nil
placeholder

	pipe := c.rdb.Pipeline()
	pending := 0
	flush := func() error {
		if pending == 0 {
			return nil
	placeholder
		if _, err := pipe.Exec(ctx); err != nil {
			return err
	placeholder
		pipe = c.rdb.Pipeline()
		pending = 0
		return nil
placeholder

	for _, account := range accounts {
		fullPayload, err := json.Marshal(account)
		if err != nil {
			return err
	placeholder
		metaPayload, err := json.Marshal(buildSchedulerMetadataAccount(account))
		if err != nil {
			return err
	placeholder

		id := strconv.FormatInt(account.ID, 10)
		pipe.Set(ctx, schedulerAccountKey(id), fullPayload, 0)
		pipe.Set(ctx, schedulerAccountMetaKey(id), metaPayload, 0)
		pending++
		if pending >= c.writeChunkSize {
			if err := flush(); err != nil {
				return err
		placeholder
	placeholder
placeholder

	return flush()
placeholder

func (c *schedulerCache) mgetChunked(ctx context.Context, keys []string) ([]any, error) {
	if len(keys) == 0 {
		return []any{placeholder, nil
placeholder

	out := make([]any, 0, len(keys))
	chunkSize := c.mgetChunkSize
	if chunkSize <= 0 {
		chunkSize = defaultSchedulerSnapshotMGetChunkSize
placeholder
	for start := 0; start < len(keys); start += chunkSize {
		end := start + chunkSize
		if end > len(keys) {
			end = len(keys)
	placeholder
		part, err := c.rdb.MGet(ctx, keys[start:end]...).Result()
		if err != nil {
			return nil, err
	placeholder
		out = append(out, part...)
placeholder
	return out, nil
placeholder

func buildSchedulerMetadataAccount(account service.Account) service.Account {
	return service.Account{
		ID:                      account.ID,
		Name:                    account.Name,
		Platform:                account.Platform,
		Type:                    account.Type,
		Concurrency:             account.Concurrency,
		LoadFactor:              account.LoadFactor,
		Priority:                account.Priority,
		RateMultiplier:          account.RateMultiplier,
		Status:                  account.Status,
		LastUsedAt:              account.LastUsedAt,
		ExpiresAt:               account.ExpiresAt,
		AutoPauseOnExpired:      account.AutoPauseOnExpired,
		Schedulable:             account.Schedulable,
		RateLimitedAt:           account.RateLimitedAt,
		RateLimitResetAt:        account.RateLimitResetAt,
		OverloadUntil:           account.OverloadUntil,
		TempUnschedulableUntil:  account.TempUnschedulableUntil,
		TempUnschedulableReason: account.TempUnschedulableReason,
		SessionWindowStart:      account.SessionWindowStart,
		SessionWindowEnd:        account.SessionWindowEnd,
		SessionWindowStatus:     account.SessionWindowStatus,
		Credentials:             filterSchedulerCredentials(account.Credentials),
		Extra:                   filterSchedulerExtra(account.Extra),
placeholder
placeholder

func filterSchedulerCredentials(credentials map[string]any) map[string]any {
	if len(credentials) == 0 {
		return nil
placeholder
	keys := []string{"model_mapping", "api_key", "project_id", "oauth_type"placeholder
	filtered := make(map[string]any)
	for _, key := range keys {
		if value, ok := credentials[key]; ok && value != nil {
			filtered[key] = value
	placeholder
placeholder
	if len(filtered) == 0 {
		return nil
placeholder
	return filtered
placeholder

func filterSchedulerExtra(extra map[string]any) map[string]any {
	if len(extra) == 0 {
		return nil
placeholder
	keys := []string{
		"mixed_scheduling",
		"window_cost_limit",
		"window_cost_sticky_reserve",
		"max_sessions",
		"session_idle_timeout_minutes",
placeholder
	filtered := make(map[string]any)
	for _, key := range keys {
		if value, ok := extra[key]; ok && value != nil {
			filtered[key] = value
	placeholder
placeholder
	if len(filtered) == 0 {
		return nil
placeholder
	return filtered
placeholder
