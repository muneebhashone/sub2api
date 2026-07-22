package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/redis/go-redis/v9"
)

const (
	schedulerBucketSetKey          = "sched:buckets"
	schedulerOutboxWatermarkKey    = "sched:outbox:watermark"
	schedulerAccountPrefix         = "sched:acc:"
	schedulerAccountMetaPrefix     = "sched:meta:"
	schedulerAccountLastUsedPrefix = "sched:acc:last_used:"
	schedulerActivePrefix          = "sched:active:"
	schedulerReadyPrefix           = "sched:ready:"
	schedulerVersionPrefix         = "sched:ver:"
	schedulerEpochPrefix           = "sched:epoch:"
	schedulerRetiredPrefix         = "sched:retired:"
	schedulerSnapshotPrefix        = "sched:"
	schedulerLockPrefix            = "sched:lock:"

	defaultSchedulerSnapshotMGetChunkSize  = 128
	defaultSchedulerSnapshotWriteChunkSize = 256
	schedulerLastUsedUpdateChunkSize       = 256

	// snapshotGraceTTLSeconds 旧快照过期的宽限期（秒）。
	// 替代立即 DEL，让正在读取旧版本的 reader 有足够时间完成 ZRANGE。
	snapshotGraceTTLSeconds = 60
)

const (
	schedulerGroupLifecycleLockPrefix      = "sched:group:lifecycle-lock:"
	schedulerGroupLifecycleOwnerTokenBytes = 16
)

var updateSchedulerLastUsedScript = redis.NewScript(`
local updated = 0
for index = 1, #ARGV do
    local key_index = (index - 1) * 2 + 1
    local candidate = tonumber(ARGV[index])
    if candidate == nil then
        return redis.error_reply('invalid last_used value')
    end
    if redis.call('EXISTS', KEYS[key_index]) == 1 then
        local current = tonumber(redis.call('GET', KEYS[key_index + 1]))
        if current == nil or candidate > current then
            redis.call('SET', KEYS[key_index + 1], ARGV[index])
            updated = updated + 1
        end
    end
end
return updated
`)

var (
	// epoch 标识 bucket writer 的代际，retired key 是持久退休标记。
	// Capture、allocate、activate 都在 Lua 内同时校验两者：-1 表示已退休，-2 表示 epoch 无效或与 token 代际不匹配；
	// allocate 与 activate 的双重校验可拦截快照写入期间发生的 Retire。
	// Retire 仅在首次退休时推进 epoch，Reopen 只清除标记并沿用该代际，因此重复调用保持幂等。
	captureBucketWriteTokenScript = redis.NewScript(`
if redis.call('EXISTS', KEYS[2]) == 1 then
    return -1
end

local currentEpoch = redis.call('GET', KEYS[1])
if currentEpoch == false then
    redis.call('SET', KEYS[1], '1')
    return 1
end

local parsedEpoch = tonumber(currentEpoch)
if parsedEpoch == nil or parsedEpoch < 1 then
    return -2
end
return parsedEpoch
`)

	allocateSnapshotVersionScript = redis.NewScript(`
if redis.call('EXISTS', KEYS[2]) == 1 then
    return -1
end

local currentEpoch = tonumber(redis.call('GET', KEYS[1]))
local expectedEpoch = tonumber(ARGV[1])
if currentEpoch == nil or expectedEpoch == nil or currentEpoch ~= expectedEpoch then
    return -2
end

return redis.call('INCR', KEYS[3])
`)

	retireBucketScript = redis.NewScript(`
local retired = redis.call('GET', KEYS[2])
local currentEpoch = tonumber(redis.call('GET', KEYS[1])) or 0

if retired == false then
    currentEpoch = currentEpoch + 1
    if currentEpoch < 1 then
        currentEpoch = 1
    end
    redis.call('SET', KEYS[1], tostring(currentEpoch))
    redis.call('SET', KEYS[2], tostring(currentEpoch))
elseif currentEpoch < 1 then
    currentEpoch = tonumber(retired) or 1
    redis.call('SET', KEYS[1], tostring(currentEpoch))
end

redis.call('SREM', KEYS[3], ARGV[1])
local currentActive = redis.call('GET', KEYS[5])
if currentActive ~= false then
    redis.call('EXPIRE', ARGV[2] .. currentActive, tonumber(ARGV[3]))
end
redis.call('DEL', KEYS[4], KEYS[5])
return currentEpoch
`)

	reopenBucketScript = redis.NewScript(`
local currentEpochRaw = redis.call('GET', KEYS[1])
local currentEpoch = tonumber(currentEpochRaw)
local retiredEpochRaw = redis.call('GET', KEYS[2])

if retiredEpochRaw == false then
    if currentEpochRaw == false then
        redis.call('SET', KEYS[1], '1')
        return 1
    end
    if currentEpoch == nil or currentEpoch < 1 then
        return -2
    end
    return currentEpoch
end

local retiredEpoch = tonumber(retiredEpochRaw)
if retiredEpoch == nil or retiredEpoch < 1 then
    return -2
end
if currentEpoch == nil or currentEpoch < retiredEpoch then
    currentEpoch = retiredEpoch
end

redis.call('SET', KEYS[1], tostring(currentEpoch))
redis.call('DEL', KEYS[2])
redis.call('SREM', KEYS[3], ARGV[1])
local currentActive = redis.call('GET', KEYS[5])
if currentActive ~= false then
    redis.call('EXPIRE', ARGV[2] .. currentActive, tonumber(ARGV[3]))
end
redis.call('DEL', KEYS[4], KEYS[5])
return currentEpoch
`)

	// 释放租约必须先比较所有者令牌再删除，过期持有者的延迟释放不能误删继任租约。
	releaseGroupLifecycleLeaseScript = redis.NewScript(`
if redis.call('GET', KEYS[1]) == ARGV[1] then
    return redis.call('DEL', KEYS[1])
end
return 0
`)

	// activateSnapshotScript 原子 CAS 切换快照版本。
	// 仅当新版本号 >= 当前激活版本时才切换，防止并发写入导致版本回滚。
	// 旧快照使用 EXPIRE 设置宽限期而非立即 DEL，避免与 reader 竞态。
	//
	// KEYS[1] = activeKey     (sched:active:{bucketplaceholder)
	// KEYS[2] = readyKey      (sched:ready:{bucketplaceholder)
	// KEYS[3] = bucketSetKey  (sched:buckets)
	// KEYS[4] = snapshotKey   (新写入的快照 key)
	// KEYS[5] = epochKey
	// KEYS[6] = retiredKey
	// ARGV[1] = 新版本号字符串
	// ARGV[2] = bucket 字符串 (用于 SADD)
	// ARGV[3] = 快照 key 前缀 (用于构造旧快照 key)
	// ARGV[4] = 宽限期 TTL 秒数
	// ARGV[5] = writer epoch
	//
	// 返回 1 = 已激活, 0 = 版本过旧未激活
	activateSnapshotScript = redis.NewScript(`
if redis.call('EXISTS', KEYS[6]) == 1 then
    redis.call('DEL', KEYS[4])
    return -1
end

local currentEpoch = tonumber(redis.call('GET', KEYS[5]))
local expectedEpoch = tonumber(ARGV[5])
if currentEpoch == nil or expectedEpoch == nil or currentEpoch ~= expectedEpoch then
    redis.call('DEL', KEYS[4])
    return -2
end

local currentActive = redis.call('GET', KEYS[1])
local newVersion = tonumber(ARGV[1])

if currentActive ~= false then
	local curVersion = tonumber(currentActive)
	if curVersion and newVersion < curVersion then
		redis.call('DEL', KEYS[4])
		return 0
	end
end

redis.call('SET', KEYS[1], ARGV[1])
redis.call('SET', KEYS[2], '1')
redis.call('SADD', KEYS[3], ARGV[2])

if currentActive ~= false and currentActive ~= ARGV[1] then
	redis.call('EXPIRE', ARGV[3] .. currentActive, tonumber(ARGV[4]))
end

return 1
`)
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
	lastUsedKeys := make([]string, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, schedulerAccountMetaKey(id))
		lastUsedKeys = append(lastUsedKeys, schedulerLastUsedKey(id))
placeholder
	values, err := c.mgetChunked(ctx, keys)
	if err != nil {
		return nil, false, err
placeholder
	lastUsedValues, err := c.mgetChunked(ctx, lastUsedKeys)
	if err != nil {
		return nil, false, err
placeholder

	accounts := make([]*service.Account, 0, len(values))
	for i, val := range values {
		if val == nil {
			return nil, false, nil
	placeholder
		account, err := decodeCachedAccount(val)
		if err != nil {
			return nil, false, err
	placeholder
		if err := applySchedulerLastUsed(account, lastUsedValues[i]); err != nil {
			return nil, false, err
	placeholder
		accounts = append(accounts, account)
placeholder

	return accounts, true, nil
placeholder

func (c *schedulerCache) CaptureBucketWriteToken(ctx context.Context, bucket service.SchedulerBucket) (service.SchedulerBucketWriteToken, error) {
	result, err := captureBucketWriteTokenScript.Run(ctx, c.rdb, []string{
		schedulerBucketKey(schedulerEpochPrefix, bucket),
		schedulerBucketKey(schedulerRetiredPrefix, bucket),
placeholder).Int64()
	if err != nil {
		return service.SchedulerBucketWriteToken{placeholder, err
placeholder
	if err := schedulerBucketWriteResultError(result, bucket); err != nil {
		return service.SchedulerBucketWriteToken{placeholder, err
placeholder
	return service.SchedulerBucketWriteToken{Bucket: bucket, Epoch: resultplaceholder, nil
placeholder

func (c *schedulerCache) RetireBucket(ctx context.Context, bucket service.SchedulerBucket) error {
	snapshotKeyPrefix := fmt.Sprintf("%s%d:%s:%s:v", schedulerSnapshotPrefix, bucket.GroupID, bucket.Platform, bucket.Mode)
	result, err := retireBucketScript.Run(ctx, c.rdb, []string{
		schedulerBucketKey(schedulerEpochPrefix, bucket),
		schedulerBucketKey(schedulerRetiredPrefix, bucket),
		schedulerBucketSetKey,
		schedulerBucketKey(schedulerReadyPrefix, bucket),
		schedulerBucketKey(schedulerActivePrefix, bucket),
placeholder, bucket.String(), snapshotKeyPrefix, snapshotGraceTTLSeconds).Int64()
	if err != nil {
		return err
placeholder
	if result < 1 {
		return fmt.Errorf("retire scheduler bucket %s returned invalid epoch %d", bucket.String(), result)
placeholder
	return nil
placeholder

func (c *schedulerCache) ReopenBucket(ctx context.Context, bucket service.SchedulerBucket) (service.SchedulerBucketWriteToken, error) {
	snapshotKeyPrefix := fmt.Sprintf("%s%d:%s:%s:v", schedulerSnapshotPrefix, bucket.GroupID, bucket.Platform, bucket.Mode)
	result, err := reopenBucketScript.Run(ctx, c.rdb, []string{
		schedulerBucketKey(schedulerEpochPrefix, bucket),
		schedulerBucketKey(schedulerRetiredPrefix, bucket),
		schedulerBucketSetKey,
		schedulerBucketKey(schedulerReadyPrefix, bucket),
		schedulerBucketKey(schedulerActivePrefix, bucket),
placeholder, bucket.String(), snapshotKeyPrefix, snapshotGraceTTLSeconds).Int64()
	if err != nil {
		return service.SchedulerBucketWriteToken{placeholder, err
placeholder
	if err := schedulerBucketWriteResultError(result, bucket); err != nil {
		return service.SchedulerBucketWriteToken{placeholder, err
placeholder
	return service.SchedulerBucketWriteToken{Bucket: bucket, Epoch: resultplaceholder, nil
placeholder

func (c *schedulerCache) TryAcquireGroupLifecycleLease(ctx context.Context, groupID int64, ttl time.Duration) (service.SchedulerGroupLifecycleLease, bool, error) {
	if groupID <= 0 {
		return service.SchedulerGroupLifecycleLease{placeholder, false, fmt.Errorf("%w: group id must be positive", service.ErrSchedulerGroupLifecycleLeaseInvalid)
placeholder
	if ttl <= 0 {
		return service.SchedulerGroupLifecycleLease{placeholder, false, fmt.Errorf("%w: ttl must be positive", service.ErrSchedulerGroupLifecycleLeaseInvalid)
placeholder
	ownerToken, err := newSchedulerGroupLifecycleOwnerToken()
	if err != nil {
		return service.SchedulerGroupLifecycleLease{placeholder, false, err
placeholder
	acquired, err := c.rdb.SetNX(ctx, schedulerGroupLifecycleLockKey(groupID), ownerToken, ttl).Result()
	if err != nil {
		return service.SchedulerGroupLifecycleLease{placeholder, false, err
placeholder
	if !acquired {
		return service.SchedulerGroupLifecycleLease{placeholder, false, nil
placeholder
	return service.SchedulerGroupLifecycleLease{GroupID: groupID, OwnerToken: ownerTokenplaceholder, true, nil
placeholder

func (c *schedulerCache) ReleaseGroupLifecycleLease(ctx context.Context, lease service.SchedulerGroupLifecycleLease) error {
	if !lease.ValidFor(lease.GroupID) {
		return service.ErrSchedulerGroupLifecycleLeaseInvalid
placeholder
	result, err := releaseGroupLifecycleLeaseScript.Run(
		ctx,
		c.rdb,
		[]string{schedulerGroupLifecycleLockKey(lease.GroupID)placeholder,
		lease.OwnerToken,
	).Int64()
	if err != nil {
		return err
placeholder
	if result == 0 {
		return fmt.Errorf("%w: group=%d", service.ErrSchedulerGroupLifecycleLeaseLost, lease.GroupID)
placeholder
	if result != 1 {
		return fmt.Errorf("release scheduler group lifecycle lease returned %d", result)
placeholder
	return nil
placeholder

func newSchedulerGroupLifecycleOwnerToken() (string, error) {
	raw := make([]byte, schedulerGroupLifecycleOwnerTokenBytes)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generate scheduler group lifecycle owner token: %w", err)
placeholder
	return hex.EncodeToString(raw), nil
placeholder

func (c *schedulerCache) SetSnapshot(ctx context.Context, bucket service.SchedulerBucket, token service.SchedulerBucketWriteToken, accounts []service.Account) error {
	if !token.ValidFor(bucket) {
		return fmt.Errorf("%w: bucket=%s", service.ErrSchedulerBucketWriteFenced, bucket.String())
placeholder
	// 分配版本与激活指针是两个 fencing 边界；中间写入的数据只有通过第二次校验才能发布。
	version, err := c.allocateSnapshotVersion(ctx, bucket, token)
	if err != nil {
		return err
placeholder
	// 快照成员最终只依赖可编码账号的有序 ID；直接复用 ID 路径，避免为
	// 随后立即丢弃的完整 Account 再分配一份临时切片。
	if _, err := c.writeSnapshotVersionAndReturnAccountIDs(ctx, bucket, version, accounts); err != nil {
		return err
placeholder
	return c.activateSnapshotVersion(ctx, bucket, token, version)
placeholder

// SetSnapshotAndReturnAccountIDs 完整发布快照，并返回实际成功编码并写入的有序账号 ID。
// 该可选能力只供同一重建批次复用，返回前仍会完成版本激活与 fencing 校验。
func (c *schedulerCache) SetSnapshotAndReturnAccountIDs(ctx context.Context, bucket service.SchedulerBucket, token service.SchedulerBucketWriteToken, accounts []service.Account) ([]int64, error) {
	if !token.ValidFor(bucket) {
		return nil, fmt.Errorf("%w: bucket=%s", service.ErrSchedulerBucketWriteFenced, bucket.String())
placeholder
	// 分配版本与激活指针是两个 fencing 边界；中间写入的数据只有通过第二次校验才能发布。
	version, err := c.allocateSnapshotVersion(ctx, bucket, token)
	if err != nil {
		return nil, err
placeholder
	accountIDs, err := c.writeSnapshotVersionAndReturnAccountIDs(ctx, bucket, version, accounts)
	if err != nil {
		return nil, err
placeholder
	if err := c.activateSnapshotVersion(ctx, bucket, token, version); err != nil {
		return nil, err
placeholder
	return accountIDs, nil
placeholder

// SetSnapshotByAccountIDs 复用同批次首次完整写入后得到的账号成员。
// 每个桶仍独立分配版本、写入有序集合并执行激活 fencing，只省略重复的账号 JSON 与全局键写入。
func (c *schedulerCache) SetSnapshotByAccountIDs(ctx context.Context, bucket service.SchedulerBucket, token service.SchedulerBucketWriteToken, accountIDs []int64) error {
	if !token.ValidFor(bucket) {
		return fmt.Errorf("%w: bucket=%s", service.ErrSchedulerBucketWriteFenced, bucket.String())
placeholder
	version, err := c.allocateSnapshotVersion(ctx, bucket, token)
	if err != nil {
		return err
placeholder
	if err := c.writeSnapshotAccountIDs(ctx, bucket, version, accountIDs); err != nil {
		return err
placeholder
	return c.activateSnapshotVersion(ctx, bucket, token, version)
placeholder

func (c *schedulerCache) allocateSnapshotVersion(ctx context.Context, bucket service.SchedulerBucket, token service.SchedulerBucketWriteToken) (string, error) {
	result, err := allocateSnapshotVersionScript.Run(ctx, c.rdb, []string{
		schedulerBucketKey(schedulerEpochPrefix, bucket),
		schedulerBucketKey(schedulerRetiredPrefix, bucket),
		schedulerBucketKey(schedulerVersionPrefix, bucket),
placeholder, token.Epoch).Int64()
	if err != nil {
		return "", err
placeholder
	if err := schedulerBucketWriteResultError(result, bucket); err != nil {
		return "", err
placeholder
	return strconv.FormatInt(result, 10), nil
placeholder

func (c *schedulerCache) writeSnapshotVersionAndReturnAccountIDs(ctx context.Context, bucket service.SchedulerBucket, version string, accounts []service.Account) ([]int64, error) {
	accountIDs, err := c.writeAccountIDs(ctx, accounts)
	if err != nil {
		return nil, err
placeholder
	if err := c.writeSnapshotAccountIDs(ctx, bucket, version, accountIDs); err != nil {
		return nil, err
placeholder
	return accountIDs, nil
placeholder

func (c *schedulerCache) writeSnapshotAccountIDs(ctx context.Context, bucket service.SchedulerBucket, version string, accountIDs []int64) error {
	members := schedulerSnapshotMembers(accountIDs)
	return c.writeSnapshotMembers(ctx, bucket, version, members)
placeholder

func schedulerSnapshotMembers(accountIDs []int64) []redis.Z {
	if len(accountIDs) == 0 {
		return nil
placeholder
	// 使用序号作为 score，保持数据库返回的排序语义；重复 ID 继续交由 Redis ZADD
	// 按最后一个 score 覆盖，与直接从账号切片构造成员时的行为一致。
	members := make([]redis.Z, 0, len(accountIDs))
	for idx, accountID := range accountIDs {
		members = append(members, redis.Z{
			Score:  float64(idx),
			Member: strconv.FormatInt(accountID, 10),
	placeholder)
placeholder
	return members
placeholder

func (c *schedulerCache) writeSnapshotMembers(ctx context.Context, bucket service.SchedulerBucket, version string, members []redis.Z) error {
	if len(members) == 0 {
		return nil
placeholder
	snapshotKey := schedulerSnapshotKey(bucket, version)
	pipe := c.rdb.Pipeline()
	for start := 0; start < len(members); start += c.writeChunkSize {
		end := start + c.writeChunkSize
		if end > len(members) {
			end = len(members)
	placeholder
		pipe.ZAdd(ctx, snapshotKey, members[start:end]...)
placeholder
	_, err := pipe.Exec(ctx)
	return err
placeholder

func (c *schedulerCache) activateSnapshotVersion(ctx context.Context, bucket service.SchedulerBucket, token service.SchedulerBucketWriteToken, version string) error {
	snapshotKey := schedulerSnapshotKey(bucket, version)
	// Phase 2: 原子 CAS 切换版本，同时再次校验退休状态与 writer epoch。
	// Lua 脚本保证：仅当新版本 >= 当前激活版本时才切换 active 指针，
	// 防止并发写入导致版本回滚。
	// 旧快照使用 EXPIRE 宽限期而非立即 DEL，避免 reader 竞态。
	activeKey := schedulerBucketKey(schedulerActivePrefix, bucket)
	readyKey := schedulerBucketKey(schedulerReadyPrefix, bucket)
	snapshotKeyPrefix := fmt.Sprintf("%s%d:%s:%s:v", schedulerSnapshotPrefix, bucket.GroupID, bucket.Platform, bucket.Mode)

	keys := []string{
		activeKey,
		readyKey,
		schedulerBucketSetKey,
		snapshotKey,
		schedulerBucketKey(schedulerEpochPrefix, bucket),
		schedulerBucketKey(schedulerRetiredPrefix, bucket),
placeholder
	args := []any{version, bucket.String(), snapshotKeyPrefix, snapshotGraceTTLSeconds, token.Epochplaceholder

	result, err := activateSnapshotScript.Run(ctx, c.rdb, keys, args...).Int64()
	if err != nil {
		return err
placeholder
	return schedulerBucketWriteResultError(result, bucket)
placeholder

func schedulerBucketWriteResultError(result int64, bucket service.SchedulerBucket) error {
	switch result {
	case -1:
		return fmt.Errorf("%w: bucket=%s", service.ErrSchedulerBucketRetired, bucket.String())
	case -2:
		return fmt.Errorf("%w: bucket=%s", service.ErrSchedulerBucketWriteFenced, bucket.String())
	default:
		return nil
placeholder
placeholder

func (c *schedulerCache) GetAccount(ctx context.Context, accountID int64) (*service.Account, error) {
	id := strconv.FormatInt(accountID, 10)
	values, err := c.rdb.MGet(ctx, schedulerAccountKey(id), schedulerLastUsedKey(id)).Result()
	if err != nil {
		return nil, err
placeholder
	if len(values) != 2 || values[0] == nil {
		return nil, nil
placeholder
	account, err := decodeCachedAccount(values[0])
	if err != nil {
		return nil, err
placeholder
	if err := applySchedulerLastUsed(account, values[1]); err != nil {
		return nil, err
placeholder
	return account, nil
placeholder

func (c *schedulerCache) SetAccount(ctx context.Context, account *service.Account) error {
	if account == nil || account.ID <= 0 {
		return nil
placeholder
	accountIDs, err := c.writeAccountIDs(ctx, []service.Account{*accountplaceholder)
	if err != nil {
		return err
placeholder
	if len(accountIDs) == 0 {
		return c.DeleteAccount(ctx, account.ID)
placeholder
	return nil
placeholder

func (c *schedulerCache) DeleteAccount(ctx context.Context, accountID int64) error {
	if accountID <= 0 {
		return nil
placeholder
	id := strconv.FormatInt(accountID, 10)
	return c.rdb.Del(ctx, schedulerAccountKey(id), schedulerAccountMetaKey(id), schedulerLastUsedKey(id)).Err()
placeholder

func (c *schedulerCache) UpdateLastUsed(ctx context.Context, updates map[int64]time.Time) error {
	if len(updates) == 0 {
		return nil
placeholder

	pipe := c.rdb.Pipeline()
	queued := 0
	keys := make([]string, 0, schedulerLastUsedUpdateChunkSize*2)
	args := make([]any, 0, schedulerLastUsedUpdateChunkSize)
	queueBatch := func() {
		if len(args) == 0 {
			return
	placeholder
		updateSchedulerLastUsedScript.Eval(ctx, pipe, keys, args...)
		queued++
		keys = make([]string, 0, schedulerLastUsedUpdateChunkSize*2)
		args = make([]any, 0, schedulerLastUsedUpdateChunkSize)
placeholder
	for id, usedAt := range updates {
		if id <= 0 {
			continue
	placeholder
		millis, err := schedulerLastUsedMillis(usedAt)
		if err != nil {
			slog.Warn("scheduler cache removes account with unencodable payload",
				"account_id", id,
				"error", err,
			)
			idText := strconv.FormatInt(id, 10)
			pipe.Del(ctx, schedulerAccountKey(idText), schedulerAccountMetaKey(idText), schedulerLastUsedKey(idText))
			queued++
			continue
	placeholder
		idText := strconv.FormatInt(id, 10)
		keys = append(keys, schedulerAccountKey(idText), schedulerLastUsedKey(idText))
		args = append(args, millis)
		if len(args) >= schedulerLastUsedUpdateChunkSize {
			queueBatch()
	placeholder
placeholder
	queueBatch()
	if queued == 0 {
		return nil
placeholder
	_, err := pipe.Exec(ctx)
	return err
placeholder

func (c *schedulerCache) TryLockBucket(ctx context.Context, bucket service.SchedulerBucket, ttl time.Duration) (bool, error) {
	key := schedulerBucketKey(schedulerLockPrefix, bucket)
	return c.rdb.SetNX(ctx, key, time.Now().UnixNano(), ttl).Result()
placeholder

func (c *schedulerCache) UnlockBucket(ctx context.Context, bucket service.SchedulerBucket) error {
	key := schedulerBucketKey(schedulerLockPrefix, bucket)
	return c.rdb.Del(ctx, key).Err()
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

func schedulerGroupLifecycleLockKey(groupID int64) string {
	return schedulerGroupLifecycleLockPrefix + strconv.FormatInt(groupID, 10)
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

func schedulerLastUsedKey(id string) string {
	return schedulerAccountLastUsedPrefix + id
placeholder

func ptrTime(t time.Time) *time.Time {
	return &t
placeholder

func schedulerLastUsedMillis(value time.Time) (int64, error) {
	if _, err := value.MarshalJSON(); err != nil {
		return 0, err
placeholder
	return value.UTC().UnixMilli(), nil
placeholder

func applySchedulerLastUsed(account *service.Account, value any) error {
	if account == nil || value == nil {
		return nil
placeholder
	var raw string
	switch typed := value.(type) {
	case string:
		raw = typed
	case []byte:
		raw = string(typed)
	default:
		return fmt.Errorf("unexpected last_used cache type: %T", value)
placeholder
	millis, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid last_used cache value %q: %w", raw, err)
placeholder
	lastUsedAt := time.UnixMilli(millis).UTC()
	if account.LastUsedAt == nil || lastUsedAt.After(*account.LastUsedAt) {
		account.LastUsedAt = ptrTime(lastUsedAt)
placeholder
	return nil
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

func (c *schedulerCache) writeAccountIDs(ctx context.Context, accounts []service.Account) ([]int64, error) {
	if len(accounts) == 0 {
		return nil, nil
placeholder

	pipe := c.rdb.Pipeline()
	accountIDs := make([]int64, 0, len(accounts))
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
		fullPayload, metaPayload, err := marshalSchedulerCacheAccount(account)
		if err != nil {
			slog.Warn("scheduler cache skips account with unencodable payload",
				"account_id", account.ID,
				"error", err,
			)
			continue
	placeholder

		id := strconv.FormatInt(account.ID, 10)
		pipe.Set(ctx, schedulerAccountKey(id), fullPayload, 0)
		pipe.Set(ctx, schedulerAccountMetaKey(id), metaPayload, 0)
		// Keep the hot LastUsedAt side key untouched: a lagging snapshot rebuild
		// must not overwrite a newer scheduler update.
		accountIDs = append(accountIDs, account.ID)
		pending++
		if pending >= c.writeChunkSize {
			if err := flush(); err != nil {
				return nil, err
		placeholder
	placeholder
placeholder

	if err := flush(); err != nil {
		return nil, err
placeholder
	return accountIDs, nil
placeholder

func marshalSchedulerCacheAccount(account service.Account) ([]byte, []byte, error) {
	fullPayload, err := json.Marshal(account)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal account: %w", err)
placeholder
	metaPayload, err := json.Marshal(buildSchedulerMetadataAccount(account))
	if err != nil {
		return nil, nil, fmt.Errorf("marshal account metadata: %w", err)
placeholder
	return fullPayload, metaPayload, nil
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
		ParentAccountID:         account.ParentAccountID,
		QuotaDimension:          account.QuotaDimension,
		AccountGroups:           filterSchedulerAccountGroups(account.AccountGroups),
		GroupIDs:                filterSchedulerGroupIDs(account.GroupIDs, account.AccountGroups),
		Credentials:             filterSchedulerCredentials(account.Credentials),
		Extra:                   filterSchedulerExtra(account.Extra),
placeholder
placeholder

func filterSchedulerAccountGroups(accountGroups []service.AccountGroup) []service.AccountGroup {
	if len(accountGroups) == 0 {
		return nil
placeholder

	filtered := make([]service.AccountGroup, 0, len(accountGroups))
	for _, ag := range accountGroups {
		if ag.GroupID <= 0 {
			continue
	placeholder
		filtered = append(filtered, service.AccountGroup{
			AccountID: ag.AccountID,
			GroupID:   ag.GroupID,
			Priority:  ag.Priority,
			CreatedAt: ag.CreatedAt,
	placeholder)
placeholder
	if len(filtered) == 0 {
		return nil
placeholder
	return filtered
placeholder

func filterSchedulerGroupIDs(groupIDs []int64, accountGroups []service.AccountGroup) []int64 {
	if len(groupIDs) == 0 && len(accountGroups) == 0 {
		return nil
placeholder

	seen := make(map[int64]struct{placeholder, len(groupIDs)+len(accountGroups))
	filtered := make([]int64, 0, len(groupIDs)+len(accountGroups))
	for _, id := range groupIDs {
		if id <= 0 {
			continue
	placeholder
		if _, ok := seen[id]; ok {
			continue
	placeholder
		seen[id] = struct{placeholder{placeholder
		filtered = append(filtered, id)
placeholder
	for _, ag := range accountGroups {
		if ag.GroupID <= 0 {
			continue
	placeholder
		if _, ok := seen[ag.GroupID]; ok {
			continue
	placeholder
		seen[ag.GroupID] = struct{placeholder{placeholder
		filtered = append(filtered, ag.GroupID)
placeholder
	if len(filtered) == 0 {
		return nil
placeholder
	return filtered
placeholder

func filterSchedulerCredentials(credentials map[string]any) map[string]any {
	if len(credentials) == 0 {
		return nil
placeholder
	keys := []string{"model_mapping", "compact_model_mapping", "api_key", "project_id", "oauth_type", "plan_type"placeholder
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
		"quota_limit",
		"quota_used",
		"quota_daily_limit",
		"quota_daily_used",
		"quota_daily_start",
		"quota_daily_reset_mode",
		"quota_daily_reset_hour",
		"quota_weekly_limit",
		"quota_weekly_used",
		"quota_weekly_start",
		"quota_weekly_reset_mode",
		"quota_weekly_reset_day",
		"quota_weekly_reset_hour",
		"quota_reset_timezone",
		"mixed_scheduling",
		"window_cost_limit",
		"window_cost_sticky_reserve",
		"max_sessions",
		"session_idle_timeout_minutes",
		"openai_oauth_responses_websockets_v2_enabled",
		"openai_oauth_responses_websockets_v2_mode",
		"openai_apikey_responses_websockets_v2_enabled",
		"openai_apikey_responses_websockets_v2_mode",
		"responses_websockets_v2_enabled",
		"openai_ws_enabled",
		"openai_ws_force_http",
		"openai_responses_mode",
		"openai_responses_supported",
		"codex_5h_used_percent",
		"codex_7d_used_percent",
		"codex_5h_reset_at",
		"codex_7d_reset_at",
		"codex_5h_reset_after_seconds",
		"codex_7d_reset_after_seconds",
		"codex_usage_updated_at",
		"auto_pause_5h_threshold",
		"auto_pause_7d_threshold",
		"auto_pause_5h_disabled",
		"auto_pause_7d_disabled",
		"model_rate_limits",
		service.UpstreamBillingProbeExtraKey,
		service.GrokMediaEligibleExtraKey,
		"grok_billing_snapshot",
placeholder
	filtered := make(map[string]any)
	for _, key := range keys {
		if value, ok := extra[key]; ok && value != nil {
			if key == service.UpstreamBillingProbeExtraKey {
				filteredProbe := filterSchedulerUpstreamBillingProbe(value)
				if filteredProbe == nil {
					continue
			placeholder
				value = filteredProbe
		placeholder
			filtered[key] = value
	placeholder
placeholder
	if len(filtered) == 0 {
		return nil
placeholder
	return filtered
placeholder

func filterSchedulerUpstreamBillingProbe(value any) map[string]any {
	source, ok := value.(map[string]any)
	if !ok {
		return nil
placeholder

	status, ok := source["status"].(string)
	if !ok || status == "" {
		return nil
placeholder
	filtered := map[string]any{"status": statusplaceholder
	for _, key := range []string{"received_at", "fresh_until", "next_probe_at"placeholder {
		if field, exists := source[key]; exists && field != nil {
			filtered[key] = field
	placeholder
placeholder
	data, ok := source["data"].(map[string]any)
	if !ok {
		return filtered
placeholder
	filteredData := make(map[string]any)
	for _, key := range []string{
		"billing_scope",
		"resolved_rate_multiplier",
		"peak_rate_enabled",
		"peak_start",
		"peak_end",
		"peak_rate_multiplier",
		"timezone",
placeholder {
		if field, exists := data[key]; exists && field != nil {
			filteredData[key] = field
	placeholder
placeholder
	if len(filteredData) > 0 {
		filtered["data"] = filteredData
placeholder
	return filtered
placeholder
