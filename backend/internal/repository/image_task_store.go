package repository

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/redis/go-redis/v9"
)

const imageTaskKeyPrefix = "image_task:"

type imageTaskStore struct {
	rdb *redis.Client
placeholder

func NewImageTaskStore(rdb *redis.Client) service.ImageTaskStore {
	return &imageTaskStore{rdb: rdbplaceholder
placeholder

func (s *imageTaskStore) Save(ctx context.Context, task *service.ImageTaskRecord, ttl time.Duration) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
placeholder
	return s.rdb.Set(ctx, imageTaskKey(task.ID), data, ttl).Err()
placeholder

func (s *imageTaskStore) Get(ctx context.Context, id string) (*service.ImageTaskRecord, error) {
	data, err := s.rdb.Get(ctx, imageTaskKey(id)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, service.ErrImageTaskNotFound
	placeholder
		return nil, err
placeholder
	var task service.ImageTaskRecord
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, err
placeholder
	return &task, nil
placeholder

func imageTaskKey(id string) string {
	return imageTaskKeyPrefix + strings.TrimSpace(id)
placeholder
