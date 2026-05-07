package repository

import (
	"context"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/redis/go-redis/v9"
)

const contentModerationFlaggedHashSetKey = "content_moderation:flagged_hashes"

type contentModerationHashCache struct {
	rdb *redis.Client
placeholder

func NewContentModerationHashCache(rdb *redis.Client) service.ContentModerationHashCache {
	return &contentModerationHashCache{rdb: rdbplaceholder
placeholder

func (c *contentModerationHashCache) RecordFlaggedInputHash(ctx context.Context, inputHash string) error {
	inputHash = strings.TrimSpace(inputHash)
	if c == nil || c.rdb == nil || inputHash == "" {
		return nil
placeholder
	return c.rdb.SAdd(ctx, contentModerationFlaggedHashSetKey, inputHash).Err()
placeholder

func (c *contentModerationHashCache) HasFlaggedInputHash(ctx context.Context, inputHash string) (bool, error) {
	inputHash = strings.TrimSpace(inputHash)
	if c == nil || c.rdb == nil || inputHash == "" {
		return false, nil
placeholder
	return c.rdb.SIsMember(ctx, contentModerationFlaggedHashSetKey, inputHash).Result()
placeholder

func (c *contentModerationHashCache) DeleteFlaggedInputHash(ctx context.Context, inputHash string) (bool, error) {
	inputHash = strings.TrimSpace(inputHash)
	if c == nil || c.rdb == nil || inputHash == "" {
		return false, nil
placeholder
	deleted, err := c.rdb.SRem(ctx, contentModerationFlaggedHashSetKey, inputHash).Result()
	if err != nil {
		return false, err
placeholder
	return deleted > 0, nil
placeholder

func (c *contentModerationHashCache) ClearFlaggedInputHashes(ctx context.Context) (int64, error) {
	if c == nil || c.rdb == nil {
		return 0, nil
placeholder
	deleted, err := c.rdb.SCard(ctx, contentModerationFlaggedHashSetKey).Result()
	if err != nil {
		return 0, err
placeholder
	if deleted == 0 {
		return 0, nil
placeholder
	if err := c.rdb.Del(ctx, contentModerationFlaggedHashSetKey).Err(); err != nil {
		return 0, err
placeholder
	return deleted, nil
placeholder

func (c *contentModerationHashCache) CountFlaggedInputHashes(ctx context.Context) (int64, error) {
	if c == nil || c.rdb == nil {
		return 0, nil
placeholder
	return c.rdb.SCard(ctx, contentModerationFlaggedHashSetKey).Result()
placeholder
