//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type schedulerCancellationCache struct {
	SchedulerCache
	cancel        context.CancelFunc
	tokenCaptures int
placeholder

func (c *schedulerCancellationCache) GetSnapshot(ctx context.Context, _ SchedulerBucket) ([]*Account, bool, error) {
	c.cancel()
	return nil, false, ctx.Err()
placeholder

func (c *schedulerCancellationCache) CaptureBucketWriteToken(ctx context.Context, _ SchedulerBucket) (SchedulerBucketWriteToken, error) {
	c.tokenCaptures++
	return SchedulerBucketWriteToken{placeholder, ctx.Err()
placeholder

func (c *schedulerCancellationCache) GetAccount(ctx context.Context, _ int64) (*Account, error) {
	c.cancel()
	return nil, ctx.Err()
placeholder

type schedulerCancellationAccountRepo struct {
	AccountRepository
	listCalls    int
	getByIDCalls int
placeholder

func (r *schedulerCancellationAccountRepo) ListSchedulableUngroupedByPlatform(ctx context.Context, _ string) ([]Account, error) {
	r.listCalls++
	return nil, ctx.Err()
placeholder

func (r *schedulerCancellationAccountRepo) GetByID(ctx context.Context, _ int64) (*Account, error) {
	r.getByIDCalls++
	return nil, ctx.Err()
placeholder

func TestSchedulerSnapshotListStopsAfterRequestCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	cache := &schedulerCancellationCache{cancel: cancelplaceholder
	repo := &schedulerCancellationAccountRepo{placeholder
	svc := NewSchedulerSnapshotService(cache, nil, repo, nil, nil)

	accounts, useMixed, err := svc.ListSchedulableAccounts(ctx, nil, PlatformOpenAI, false)

	require.ErrorIs(t, err, context.Canceled)
	require.Nil(t, accounts)
	require.False(t, useMixed)
	require.Zero(t, cache.tokenCaptures, "canceled requests must not capture a cache publish token")
	require.Zero(t, repo.listCalls, "canceled requests must not fall back to the database")
placeholder

func TestSchedulerSnapshotGetAccountStopsAfterRequestCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	cache := &schedulerCancellationCache{cancel: cancelplaceholder
	repo := &schedulerCancellationAccountRepo{placeholder
	svc := NewSchedulerSnapshotService(cache, nil, repo, nil, nil)

	account, err := svc.GetAccount(ctx, 42)

	require.ErrorIs(t, err, context.Canceled)
	require.Nil(t, account)
	require.Zero(t, repo.getByIDCalls, "canceled requests must not fall back to the database")
placeholder
