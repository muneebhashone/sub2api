package repository

import (
	"context"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/redis/go-redis/v9"
)

const dashboardStatsCacheKey = "dashboard:stats:v1"

type dashboardCache struct {
	rdb       *redis.Client
	keyPrefix string
placeholder

func NewDashboardCache(rdb *redis.Client, cfg *config.Config) service.DashboardStatsCache {
	prefix := "sub2api:"
	if cfg != nil {
		prefix = strings.TrimSpace(cfg.Dashboard.KeyPrefix)
placeholder
	if prefix != "" && !strings.HasSuffix(prefix, ":") {
		prefix += ":"
placeholder
	return &dashboardCache{
		rdb:       rdb,
		keyPrefix: prefix,
placeholder
placeholder

func (c *dashboardCache) GetDashboardStats(ctx context.Context) (string, error) {
	val, err := c.rdb.Get(ctx, c.buildKey()).Result()
	if err != nil {
		if err == redis.Nil {
			return "", service.ErrDashboardStatsCacheMiss
	placeholder
		return "", err
placeholder
	return val, nil
placeholder

func (c *dashboardCache) SetDashboardStats(ctx context.Context, data string, ttl time.Duration) error {
	return c.rdb.Set(ctx, c.buildKey(), data, ttl).Err()
placeholder

func (c *dashboardCache) buildKey() string {
	if c.keyPrefix == "" {
		return dashboardStatsCacheKey
placeholder
	return c.keyPrefix + dashboardStatsCacheKey
placeholder

func (c *dashboardCache) DeleteDashboardStats(ctx context.Context) error {
	return c.rdb.Del(ctx, c.buildKey()).Err()
placeholder
