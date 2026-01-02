package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/redis/go-redis/v9"
)

const (
	opsLatestMetricsKey = "ops:metrics:latest"

	opsDashboardOverviewKeyPrefix = "ops:dashboard:overview:"

	opsLatestMetricsTTL = 10 * time.Second
)

func (r *OpsRepository) GetCachedLatestSystemMetric(ctx context.Context) (*service.OpsMetrics, error) {
	if ctx == nil {
		ctx = context.Background()
placeholder
	if r == nil || r.rdb == nil {
		return nil, nil
placeholder

	data, err := r.rdb.Get(ctx, opsLatestMetricsKey).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
placeholder
	if err != nil {
		return nil, fmt.Errorf("redis get cached latest system metric: %w", err)
placeholder

	var metric service.OpsMetrics
	if err := json.Unmarshal(data, &metric); err != nil {
		return nil, fmt.Errorf("unmarshal cached latest system metric: %w", err)
placeholder
	return &metric, nil
placeholder

func (r *OpsRepository) SetCachedLatestSystemMetric(ctx context.Context, metric *service.OpsMetrics) error {
	if metric == nil {
		return nil
placeholder
	if ctx == nil {
		ctx = context.Background()
placeholder
	if r == nil || r.rdb == nil {
		return nil
placeholder

	data, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("marshal cached latest system metric: %w", err)
placeholder
	return r.rdb.Set(ctx, opsLatestMetricsKey, data, opsLatestMetricsTTL).Err()
placeholder

func (r *OpsRepository) GetCachedDashboardOverview(ctx context.Context, timeRange string) (*service.DashboardOverviewData, error) {
	if ctx == nil {
		ctx = context.Background()
placeholder
	if r == nil || r.rdb == nil {
		return nil, nil
placeholder
	rangeKey := strings.TrimSpace(timeRange)
	if rangeKey == "" {
		rangeKey = "1h"
placeholder

	key := opsDashboardOverviewKeyPrefix + rangeKey
	data, err := r.rdb.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
placeholder
	if err != nil {
		return nil, fmt.Errorf("redis get cached dashboard overview: %w", err)
placeholder

	var overview service.DashboardOverviewData
	if err := json.Unmarshal(data, &overview); err != nil {
		return nil, fmt.Errorf("unmarshal cached dashboard overview: %w", err)
placeholder
	return &overview, nil
placeholder

func (r *OpsRepository) SetCachedDashboardOverview(ctx context.Context, timeRange string, data *service.DashboardOverviewData, ttl time.Duration) error {
	if data == nil {
		return nil
placeholder
	if ttl <= 0 {
		ttl = 10 * time.Second
placeholder
	if ctx == nil {
		ctx = context.Background()
placeholder
	if r == nil || r.rdb == nil {
		return nil
placeholder

	rangeKey := strings.TrimSpace(timeRange)
	if rangeKey == "" {
		rangeKey = "1h"
placeholder

	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal cached dashboard overview: %w", err)
placeholder
	key := opsDashboardOverviewKeyPrefix + rangeKey
	return r.rdb.Set(ctx, key, payload, ttl).Err()
placeholder

func (r *OpsRepository) PingRedis(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
placeholder
	if r == nil || r.rdb == nil {
		return errors.New("redis client is nil")
placeholder
	return r.rdb.Ping(ctx).Err()
placeholder
