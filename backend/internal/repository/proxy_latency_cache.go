package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/redis/go-redis/v9"
)

const proxyLatencyKeyPrefix = "proxy:latency:"

func proxyLatencyKey(proxyID int64) string {
	return fmt.Sprintf("%s%d", proxyLatencyKeyPrefix, proxyID)
placeholder

type proxyLatencyCache struct {
	rdb *redis.Client
placeholder

func NewProxyLatencyCache(rdb *redis.Client) service.ProxyLatencyCache {
	return &proxyLatencyCache{rdb: rdbplaceholder
placeholder

func (c *proxyLatencyCache) GetProxyLatencies(ctx context.Context, proxyIDs []int64) (map[int64]*service.ProxyLatencyInfo, error) {
	results := make(map[int64]*service.ProxyLatencyInfo)
	if len(proxyIDs) == 0 {
		return results, nil
placeholder

	keys := make([]string, 0, len(proxyIDs))
	for _, id := range proxyIDs {
		keys = append(keys, proxyLatencyKey(id))
placeholder

	values, err := c.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return results, err
placeholder

	for i, raw := range values {
		if raw == nil {
			continue
	placeholder
		var payload []byte
		switch v := raw.(type) {
		case string:
			payload = []byte(v)
		case []byte:
			payload = v
		default:
			continue
	placeholder
		var info service.ProxyLatencyInfo
		if err := json.Unmarshal(payload, &info); err != nil {
			continue
	placeholder
		results[proxyIDs[i]] = &info
placeholder

	return results, nil
placeholder

func (c *proxyLatencyCache) SetProxyLatency(ctx context.Context, proxyID int64, info *service.ProxyLatencyInfo) error {
	if info == nil {
		return nil
placeholder
	payload, err := json.Marshal(info)
	if err != nil {
		return err
placeholder
	return c.rdb.Set(ctx, proxyLatencyKey(proxyID), payload, 0).Err()
placeholder
