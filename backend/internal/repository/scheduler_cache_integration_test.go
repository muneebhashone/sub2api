//go:build integration

package repository

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestSchedulerCacheSnapshotUsesSlimMetadataButKeepsFullAccount(t *testing.T) {
	ctx := context.Background()
	rdb := testRedis(t)
	cache := NewSchedulerCache(rdb)

	bucket := service.SchedulerBucket{GroupID: 2, Platform: service.PlatformGemini, Mode: service.SchedulerModeSingleplaceholder
	now := time.Now().UTC().Truncate(time.Second)
	limitReset := now.Add(10 * time.Minute)
	overloadUntil := now.Add(2 * time.Minute)
	tempUnschedUntil := now.Add(3 * time.Minute)
	windowEnd := now.Add(5 * time.Hour)

	account := service.Account{
		ID:          101,
		Name:        "gemini-heavy",
		Platform:    service.PlatformGemini,
		Type:        service.AccountTypeOAuth,
		Status:      service.StatusActive,
		Schedulable: true,
		Concurrency: 3,
		Priority:    7,
		LastUsedAt:  &now,
placeholder
			"api_key":       "gemini-api-key",
			"access_token":  "secret-access-token",
			"project_id":    "proj-1",
			"oauth_type":    "ai_studio",
			"model_mapping": map[string]any{"gemini-2.5-pro": "gemini-2.5-pro"placeholder,
			"huge_blob":     strings.Repeat("x", 4096),
	placeholder,
		Extra: map[string]any{
			"mixed_scheduling":             true,
			"window_cost_limit":            12.5,
			"window_cost_sticky_reserve":   8.0,
			"max_sessions":                 4,
			"session_idle_timeout_minutes": 11,
			"unused_large_field":           strings.Repeat("y", 4096),
	placeholder,
		RateLimitResetAt:       &limitReset,
		OverloadUntil:          &overloadUntil,
		TempUnschedulableUntil: &tempUnschedUntil,
		SessionWindowStart:     &now,
		SessionWindowEnd:       &windowEnd,
		SessionWindowStatus:    "active",
		GroupIDs:               []int64{bucket.GroupIDplaceholder,
		AccountGroups: []service.AccountGroup{
			{
				AccountID: 101,
				GroupID:   bucket.GroupID,
				Priority:  5,
				Group:     &service.Group{ID: bucket.GroupID, Name: "gemini-group"placeholder,
		placeholder,
	placeholder,
placeholder

	token, err := cache.CaptureBucketWriteToken(ctx, bucket)
placeholder
	require.NoError(t, cache.SetSnapshot(ctx, bucket, token, []service.Account{accountplaceholder))

	snapshot, hit, err := cache.GetSnapshot(ctx, bucket)
placeholder
	require.True(t, hit)
	require.Len(t, snapshot, 1)

	got := snapshot[0]
	require.NotNil(t, got)
	require.Equal(t, "gemini-api-key", got.GetCredential("api_key"))
	require.Equal(t, "proj-1", got.GetCredential("project_id"))
	require.Equal(t, "ai_studio", got.GetCredential("oauth_type"))
	require.NotEmpty(t, got.GetModelMapping())
	require.Empty(t, got.GetCredential("access_token"))
	require.Empty(t, got.GetCredential("huge_blob"))
	require.Equal(t, true, got.Extra["mixed_scheduling"])
	require.Equal(t, 12.5, got.GetWindowCostLimit())
	require.Equal(t, 8.0, got.GetWindowCostStickyReserve())
	require.Equal(t, 4, got.GetMaxSessions())
	require.Equal(t, 11, got.GetSessionIdleTimeoutMinutes())
	require.Nil(t, got.Extra["unused_large_field"])
	require.Equal(t, []int64{bucket.GroupIDplaceholder, got.GroupIDs)
	require.Len(t, got.AccountGroups, 1)
	require.Equal(t, account.ID, got.AccountGroups[0].AccountID)
	require.Equal(t, bucket.GroupID, got.AccountGroups[0].GroupID)
	require.Nil(t, got.AccountGroups[0].Group)

	full, err := cache.GetAccount(ctx, account.ID)
placeholder
	require.NotNil(t, full)
	require.Equal(t, "secret-access-token", full.GetCredential("access_token"))
	require.Equal(t, strings.Repeat("x", 4096), full.GetCredential("huge_blob"))
	require.Len(t, full.AccountGroups, 1)
	require.NotNil(t, full.AccountGroups[0].Group)
placeholder

func TestSchedulerCacheRetireAndReopenFencesOldEpochIntegration(t *testing.T) {
	ctx := context.Background()
	rdb := testRedis(t)
	cache := NewSchedulerCache(rdb)
	bucket := service.SchedulerBucket{GroupID: 77, Platform: service.PlatformAntigravity, Mode: service.SchedulerModeForcedplaceholder
	account := service.Account{ID: 7701, Platform: service.PlatformAntigravity, Type: service.AccountTypeOAuthplaceholder

	oldToken, err := cache.CaptureBucketWriteToken(ctx, bucket)
placeholder
	require.NoError(t, cache.SetSnapshot(ctx, bucket, oldToken, []service.Account{accountplaceholder))
	require.NoError(t, cache.RetireBucket(ctx, bucket))
	require.NoError(t, cache.RetireBucket(ctx, bucket))

	_, hit, err := cache.GetSnapshot(ctx, bucket)
placeholder
	require.False(t, hit)
	_, err = cache.CaptureBucketWriteToken(ctx, bucket)
	require.ErrorIs(t, err, service.ErrSchedulerBucketRetired)
	require.ErrorIs(t, cache.SetSnapshot(ctx, bucket, oldToken, []service.Account{accountplaceholder), service.ErrSchedulerBucketRetired)

	newToken, err := cache.ReopenBucket(ctx, bucket)
placeholder
	require.Greater(t, newToken.Epoch, oldToken.Epoch)
	require.ErrorIs(t, cache.SetSnapshot(ctx, bucket, oldToken, []service.Account{accountplaceholder), service.ErrSchedulerBucketWriteFenced)
	require.NoError(t, cache.SetSnapshot(ctx, bucket, newToken, []service.Account{accountplaceholder))

	snapshot, hit, err := cache.GetSnapshot(ctx, bucket)
placeholder
	require.True(t, hit)
	require.Len(t, snapshot, 1)
	require.Equal(t, account.ID, snapshot[0].ID)
placeholder

func TestSchedulerCacheGroupLifecycleLeaseOwnerAndTTLIntegration(t *testing.T) {
	ctx := context.Background()
	rdb := testRedis(t)
	cache := NewSchedulerCache(rdb)
	const groupID int64 = 78
	const ttl = 500 * time.Millisecond

	first, acquired, err := cache.TryAcquireGroupLifecycleLease(ctx, groupID, ttl)
placeholder
	require.True(t, acquired)
	pttl, err := rdb.PTTL(ctx, schedulerGroupLifecycleLockKey(groupID)).Result()
placeholder
	require.Positive(t, pttl)
	require.LessOrEqual(t, pttl, ttl)

	var second service.SchedulerGroupLifecycleLease
	require.Eventually(t, func() bool {
		var acquireErr error
		second, acquired, acquireErr = cache.TryAcquireGroupLifecycleLease(ctx, groupID, time.Minute)
		return acquireErr == nil && acquired
placeholder, 5*time.Second, 20*time.Millisecond)
	require.NotEqual(t, first.OwnerToken, second.OwnerToken)

	require.ErrorIs(t, cache.ReleaseGroupLifecycleLease(ctx, first), service.ErrSchedulerGroupLifecycleLeaseLost)
	_, acquired, err = cache.TryAcquireGroupLifecycleLease(ctx, groupID, time.Minute)
placeholder
	require.False(t, acquired, "a stale release must not delete the successor lease")

	require.NoError(t, cache.ReleaseGroupLifecycleLease(ctx, second))
	require.ErrorIs(t, cache.ReleaseGroupLifecycleLease(ctx, second), service.ErrSchedulerGroupLifecycleLeaseLost)
	third, acquired, err := cache.TryAcquireGroupLifecycleLease(ctx, groupID, time.Minute)
placeholder
	require.True(t, acquired)
	require.True(t, third.ValidFor(groupID))
	require.NoError(t, cache.ReleaseGroupLifecycleLease(ctx, third))
placeholder
