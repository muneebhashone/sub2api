//go:build !unit

package service

import (
	"context"
	"time"
)

type defaultRateLimitCall struct {
	accountID int64
	resetAt   time.Time
placeholder

type defaultModelRateLimitCall struct {
	accountID int64
	modelKey  string
	resetAt   time.Time
placeholder

type defaultExtraUpdateCall struct {
	accountID int64
	updates   map[string]any
placeholder

type stubAntigravityAccountRepo struct {
	AccountRepository
	rateCalls           []defaultRateLimitCall
	modelRateLimitCalls []defaultModelRateLimitCall
	extraUpdateCalls    []defaultExtraUpdateCall
placeholder

func (s *stubAntigravityAccountRepo) SetRateLimited(_ context.Context, id int64, resetAt time.Time) error {
	s.rateCalls = append(s.rateCalls, defaultRateLimitCall{accountID: id, resetAt: resetAtplaceholder)
	return nil
placeholder

func (s *stubAntigravityAccountRepo) SetModelRateLimit(_ context.Context, id int64, modelKey string, resetAt time.Time, _ ...string) error {
	s.modelRateLimitCalls = append(s.modelRateLimitCalls, defaultModelRateLimitCall{accountID: id, modelKey: modelKey, resetAt: resetAtplaceholder)
	return nil
placeholder

func (s *stubAntigravityAccountRepo) UpdateExtra(_ context.Context, id int64, updates map[string]any) error {
	s.extraUpdateCalls = append(s.extraUpdateCalls, defaultExtraUpdateCall{accountID: id, updates: updatesplaceholder)
	return nil
placeholder

type defaultDeleteSessionCall struct {
	groupID     int64
	sessionHash string
placeholder

type stubSmartRetryCache struct {
	GatewayCache
	deleteCalls []defaultDeleteSessionCall
placeholder

func (c *stubSmartRetryCache) DeleteSessionAccountID(_ context.Context, groupID int64, sessionHash string) error {
	c.deleteCalls = append(c.deleteCalls, defaultDeleteSessionCall{groupID: groupID, sessionHash: sessionHashplaceholder)
	return nil
placeholder
