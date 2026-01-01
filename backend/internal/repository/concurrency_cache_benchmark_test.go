package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// 基准测试用 TTL 配置
const benchSlotTTLMinutes = 15

var benchSlotTTL = time.Duration(benchSlotTTLMinutes) * time.Minute

// BenchmarkAccountConcurrency 用于对比 SCAN 与有序集合的计数性能。
func BenchmarkAccountConcurrency(b *testing.B) {
	rdb := newBenchmarkRedisClient(b)
	defer func() {
		_ = rdb.Close()
placeholder()

	cache, _ := NewConcurrencyCache(rdb, benchSlotTTLMinutes).(*concurrencyCache)
	ctx := context.Background()

	for _, size := range []int{10, 100, 1000placeholder {
		size := size
		b.Run(fmt.Sprintf("zset/slots=%d", size), func(b *testing.B) {
			accountID := time.Now().UnixNano()
			key := accountSlotKey(accountID)

			b.StopTimer()
			members := make([]redis.Z, 0, size)
			now := float64(time.Now().Unix())
			for i := 0; i < size; i++ {
				members = append(members, redis.Z{
					Score:  now,
					Member: fmt.Sprintf("req_%d", i),
			placeholder)
		placeholder
			if err := rdb.ZAdd(ctx, key, members...).Err(); err != nil {
				b.Fatalf("初始化有序集合失败: %v", err)
		placeholder
			if err := rdb.Expire(ctx, key, benchSlotTTL).Err(); err != nil {
				b.Fatalf("设置有序集合 TTL 失败: %v", err)
		placeholder
			b.StartTimer()

			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := cache.GetAccountConcurrency(ctx, accountID); err != nil {
					b.Fatalf("获取并发数量失败: %v", err)
			placeholder
		placeholder

			b.StopTimer()
			if err := rdb.Del(ctx, key).Err(); err != nil {
				b.Fatalf("清理有序集合失败: %v", err)
		placeholder
	placeholder)

		b.Run(fmt.Sprintf("scan/slots=%d", size), func(b *testing.B) {
			accountID := time.Now().UnixNano()
			pattern := fmt.Sprintf("%s%d:*", accountSlotKeyPrefix, accountID)
			keys := make([]string, 0, size)

			b.StopTimer()
			pipe := rdb.Pipeline()
			for i := 0; i < size; i++ {
				key := fmt.Sprintf("%s%d:req_%d", accountSlotKeyPrefix, accountID, i)
				keys = append(keys, key)
				pipe.Set(ctx, key, "1", benchSlotTTL)
		placeholder
			if _, err := pipe.Exec(ctx); err != nil {
				b.Fatalf("初始化扫描键失败: %v", err)
		placeholder
			b.StartTimer()

			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := scanSlotCount(ctx, rdb, pattern); err != nil {
					b.Fatalf("SCAN 计数失败: %v", err)
			placeholder
		placeholder

			b.StopTimer()
			if err := rdb.Del(ctx, keys...).Err(); err != nil {
				b.Fatalf("清理扫描键失败: %v", err)
		placeholder
	placeholder)
placeholder
placeholder

func scanSlotCount(ctx context.Context, rdb *redis.Client, pattern string) (int, error) {
	var cursor uint64
	count := 0
	for {
		keys, nextCursor, err := rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return 0, err
	placeholder
		count += len(keys)
		if nextCursor == 0 {
			break
	placeholder
		cursor = nextCursor
placeholder
	return count, nil
placeholder

func newBenchmarkRedisClient(b *testing.B) *redis.Client {
	b.Helper()

	redisURL := os.Getenv("TEST_REDIS_URL")
	if redisURL == "" {
		b.Skip("未设置 TEST_REDIS_URL，跳过 Redis 基准测试")
placeholder

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		b.Fatalf("解析 TEST_REDIS_URL 失败: %v", err)
placeholder

	client := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		b.Fatalf("Redis 连接失败: %v", err)
placeholder

	return client
placeholder
