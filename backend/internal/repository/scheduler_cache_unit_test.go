//go:build unit

package repository

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func newSchedulerCacheUnit(t *testing.T) *schedulerCache {
	cache, _ := newSchedulerCacheUnitWithRedis(t)
	return cache
placeholder

func newSchedulerCacheUnitWithRedis(t *testing.T) (*schedulerCache, *miniredis.Miniredis) {
placeholder
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()placeholder)
	t.Cleanup(func() { _ = rdb.Close() placeholder)
	cache, ok := newSchedulerCacheWithChunkSizes(rdb, defaultSchedulerSnapshotMGetChunkSize, defaultSchedulerSnapshotWriteChunkSize).(*schedulerCache)
	require.True(t, ok)
	return cache, mr
placeholder

func TestSchedulerCacheWriteAccountIDsSkipsUnencodableTimes(t *testing.T) {
	ctx := context.Background()
	cache := newSchedulerCacheUnit(t)
	invalidTime := time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC)

	accountIDs, err := cache.writeAccountIDs(ctx, []service.Account{
		{ID: 111, Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKeyplaceholder,
		{ID: 112, Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey, ExpiresAt: &invalidTimeplaceholder,
placeholder)
placeholder
	require.Equal(t, []int64{111placeholder, accountIDs)

	cached, err := cache.GetAccount(ctx, 111)
placeholder
	require.NotNil(t, cached)

	invalid, err := cache.GetAccount(ctx, 112)
placeholder
	require.Nil(t, invalid)
placeholder

func TestSchedulerCacheSetAccountClearsUnencodablePayload(t *testing.T) {
	ctx := context.Background()
	cache := newSchedulerCacheUnit(t)

	account := service.Account{ID: 113, Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKeyplaceholder
	require.NoError(t, cache.SetAccount(ctx, &account))

	invalidTime := time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC)
	account.ExpiresAt = &invalidTime
	require.NoError(t, cache.SetAccount(ctx, &account))

	cached, err := cache.GetAccount(ctx, account.ID)
placeholder
	require.Nil(t, cached)
placeholder

func TestSchedulerCacheUpdateLastUsedClearsUnencodablePayload(t *testing.T) {
	ctx := context.Background()
	cache := newSchedulerCacheUnit(t)
	account := service.Account{ID: 114, Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKeyplaceholder
	require.NoError(t, cache.SetAccount(ctx, &account))

	invalidTime := time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC)
	require.NoError(t, cache.UpdateLastUsed(ctx, map[int64]time.Time{account.ID: invalidTimeplaceholder))

	cached, err := cache.GetAccount(ctx, account.ID)
placeholder
	require.Nil(t, cached)
placeholder

func TestSchedulerCacheSnapshotAccountIDReusePreservesPayloadAndMembers(t *testing.T) {
	ctx := context.Background()
	cache, _ := newSchedulerCacheUnitWithRedis(t)
	invalidTime := time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC)
	validOne := service.Account{
		ID:          701,
		Name:        "first",
		Platform:    service.PlatformOpenAI,
		Type:        service.AccountTypeOAuth,
placeholder"model_mapping": map[string]any{"z": "last", "a": "first"placeholderplaceholder,
		Extra:       map[string]any{"mixed_scheduling": trueplaceholder,
		GroupIDs:    []int64{17placeholder,
placeholder
	validTwo := service.Account{ID: 702, Name: "second", Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKeyplaceholder
	invalid := service.Account{ID: 799, Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey, ExpiresAt: &invalidTimeplaceholder
	accounts := []service.Account{validOne, invalid, validTwo, validOneplaceholder

	single := service.SchedulerBucket{GroupID: 17, Platform: service.PlatformOpenAI, Mode: service.SchedulerModeSingleplaceholder
	singleToken, err := cache.CaptureBucketWriteToken(ctx, single)
placeholder
	accountIDs, err := cache.SetSnapshotAndReturnAccountIDs(ctx, single, singleToken, accounts)
placeholder
	require.Equal(t, []int64{701, 702, 701placeholder, accountIDs, "应保留可编码账号的原顺序和重复项")

	wantFull, err := json.Marshal(validOne)
placeholder
	wantMeta, err := json.Marshal(buildSchedulerMetadataAccount(validOne))
placeholder
	fullBefore, err := cache.rdb.Get(ctx, schedulerAccountKey("701")).Bytes()
placeholder
	metaBefore, err := cache.rdb.Get(ctx, schedulerAccountMetaKey("701")).Bytes()
placeholder
	require.Equal(t, wantFull, fullBefore)
	require.Equal(t, wantMeta, metaBefore)

	forced := service.SchedulerBucket{GroupID: 17, Platform: service.PlatformOpenAI, Mode: service.SchedulerModeForcedplaceholder
	forcedToken, err := cache.CaptureBucketWriteToken(ctx, forced)
placeholder
	require.NoError(t, cache.SetSnapshotByAccountIDs(ctx, forced, forcedToken, accountIDs))

	fullAfter, err := cache.rdb.Get(ctx, schedulerAccountKey("701")).Bytes()
placeholder
	metaAfter, err := cache.rdb.Get(ctx, schedulerAccountMetaKey("701")).Bytes()
placeholder
	require.Equal(t, fullBefore, fullAfter, "ID-only 路径不得重写完整账号键")
	require.Equal(t, metaBefore, metaAfter, "ID-only 路径不得重写调度元数据键")

	for _, bucket := range []service.SchedulerBucket{single, forcedplaceholder {
		version, err := cache.rdb.Get(ctx, schedulerBucketKey(schedulerActivePrefix, bucket)).Result()
	placeholder
		members, err := cache.rdb.ZRange(ctx, schedulerSnapshotKey(bucket, version), 0, -1).Result()
	placeholder
		require.Equal(t, []string{"702", "701"placeholder, members, bucket.String())
placeholder
	missing, err := cache.GetAccount(ctx, invalid.ID)
placeholder
	require.Nil(t, missing)
placeholder

func TestSchedulerCacheSetSnapshotMatchesIDPublishing(t *testing.T) {
	ctx := context.Background()
	cache, _ := newSchedulerCacheUnitWithRedis(t)
	invalidTime := time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC)
	validOne := service.Account{
		ID:          721,
		Name:        "first",
		Platform:    service.PlatformOpenAI,
		Type:        service.AccountTypeOAuth,
placeholder"model_mapping": map[string]any{"source": "target"placeholderplaceholder,
		Extra:       map[string]any{"mixed_scheduling": trueplaceholder,
		GroupIDs:    []int64{21placeholder,
placeholder
	validTwo := service.Account{ID: 722, Name: "second", Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKeyplaceholder
	invalid := service.Account{ID: 799, Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey, ExpiresAt: &invalidTimeplaceholder
	accounts := []service.Account{validOne, invalid, validTwo, validOneplaceholder

	normal := service.SchedulerBucket{GroupID: 21, Platform: service.PlatformOpenAI, Mode: service.SchedulerModeSingleplaceholder
	normalToken, err := cache.CaptureBucketWriteToken(ctx, normal)
placeholder
	require.NoError(t, cache.SetSnapshot(ctx, normal, normalToken, accounts))

	fullBefore, err := cache.rdb.Get(ctx, schedulerAccountKey("721")).Bytes()
placeholder
	metaBefore, err := cache.rdb.Get(ctx, schedulerAccountMetaKey("721")).Bytes()
placeholder

	idOnly := service.SchedulerBucket{GroupID: 21, Platform: service.PlatformOpenAI, Mode: service.SchedulerModeForcedplaceholder
	idOnlyToken, err := cache.CaptureBucketWriteToken(ctx, idOnly)
placeholder
	accountIDs, err := cache.SetSnapshotAndReturnAccountIDs(ctx, idOnly, idOnlyToken, accounts)
placeholder
	require.Equal(t, []int64{721, 722, 721placeholder, accountIDs)

	fullAfter, err := cache.rdb.Get(ctx, schedulerAccountKey("721")).Bytes()
placeholder
	metaAfter, err := cache.rdb.Get(ctx, schedulerAccountMetaKey("721")).Bytes()
placeholder
	require.Equal(t, fullBefore, fullAfter, "普通快照和 ID 发布必须写入相同完整账号 payload")
	require.Equal(t, metaBefore, metaAfter, "普通快照和 ID 发布必须写入相同元数据 payload")

	for _, bucket := range []service.SchedulerBucket{normal, idOnlyplaceholder {
		version, err := cache.rdb.Get(ctx, schedulerBucketKey(schedulerActivePrefix, bucket)).Result()
	placeholder
		members, err := cache.rdb.ZRange(ctx, schedulerSnapshotKey(bucket, version), 0, -1).Result()
	placeholder
		require.Equal(t, []string{"722", "721"placeholder, members, bucket.String())
placeholder
placeholder

func TestSchedulerCacheSnapshotAccountIDReuseKeepsEmptySnapshotSemantics(t *testing.T) {
	ctx := context.Background()
	cache := newSchedulerCacheUnit(t)
	invalidTime := time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC)
	accounts := []service.Account{{ID: 811, Platform: service.PlatformOpenAI, ExpiresAt: &invalidTimeplaceholderplaceholder

	single := service.SchedulerBucket{GroupID: 18, Platform: service.PlatformOpenAI, Mode: service.SchedulerModeSingleplaceholder
	singleToken, err := cache.CaptureBucketWriteToken(ctx, single)
placeholder
	accountIDs, err := cache.SetSnapshotAndReturnAccountIDs(ctx, single, singleToken, accounts)
placeholder
	require.Empty(t, accountIDs)

	forced := service.SchedulerBucket{GroupID: 18, Platform: service.PlatformOpenAI, Mode: service.SchedulerModeForcedplaceholder
	forcedToken, err := cache.CaptureBucketWriteToken(ctx, forced)
placeholder
	require.NoError(t, cache.SetSnapshotByAccountIDs(ctx, forced, forcedToken, accountIDs))

	for _, bucket := range []service.SchedulerBucket{single, forcedplaceholder {
		ready, err := cache.rdb.Get(ctx, schedulerBucketKey(schedulerReadyPrefix, bucket)).Result()
	placeholder
		require.Equal(t, "1", ready)
		snapshot, hit, err := cache.GetSnapshot(ctx, bucket)
	placeholder
		require.False(t, hit, bucket.String())
		require.Nil(t, snapshot)
placeholder
placeholder

func TestSchedulerCacheSetSnapshotByAccountIDsKeepsFencing(t *testing.T) {
	ctx := context.Background()
	cache := newSchedulerCacheUnit(t)
	bucket := service.SchedulerBucket{GroupID: 19, Platform: service.PlatformOpenAI, Mode: service.SchedulerModeForcedplaceholder

	err := cache.SetSnapshotByAccountIDs(ctx, bucket, service.SchedulerBucketWriteToken{placeholder, []int64{901placeholder)
	require.ErrorIs(t, err, service.ErrSchedulerBucketWriteFenced)
	_, err = cache.rdb.Get(ctx, schedulerBucketKey(schedulerVersionPrefix, bucket)).Result()
	require.ErrorIs(t, err, redis.Nil)

	token, err := cache.CaptureBucketWriteToken(ctx, bucket)
placeholder
	require.NoError(t, cache.RetireBucket(ctx, bucket))
	err = cache.SetSnapshotByAccountIDs(ctx, bucket, token, []int64{901placeholder)
	require.ErrorIs(t, err, service.ErrSchedulerBucketRetired)
placeholder

func TestSchedulerCacheSetSnapshotByAccountIDsDoesNotResurrectDeletedAccount(t *testing.T) {
	ctx := context.Background()
	cache := newSchedulerCacheUnit(t)
	account := service.Account{ID: 902, Platform: service.PlatformOpenAI, Type: service.AccountTypeOAuthplaceholder
	single := service.SchedulerBucket{GroupID: 20, Platform: service.PlatformOpenAI, Mode: service.SchedulerModeSingleplaceholder
	singleToken, err := cache.CaptureBucketWriteToken(ctx, single)
placeholder
	accountIDs, err := cache.SetSnapshotAndReturnAccountIDs(ctx, single, singleToken, []service.Account{accountplaceholder)
placeholder
	require.Equal(t, []int64{account.IDplaceholder, accountIDs)
	require.NoError(t, cache.DeleteAccount(ctx, account.ID))

	forced := service.SchedulerBucket{GroupID: 20, Platform: service.PlatformOpenAI, Mode: service.SchedulerModeForcedplaceholder
	forcedToken, err := cache.CaptureBucketWriteToken(ctx, forced)
placeholder
	require.NoError(t, cache.SetSnapshotByAccountIDs(ctx, forced, forcedToken, accountIDs))

	full, err := cache.GetAccount(ctx, account.ID)
placeholder
	require.Nil(t, full, "ID-only 发布不得复活已删除的完整账号键")
	snapshot, hit, err := cache.GetSnapshot(ctx, forced)
placeholder
	require.False(t, hit, "元数据缺失时必须安全回源，而不是返回残缺快照")
	require.Nil(t, snapshot)
placeholder

func TestMarshalSchedulerCacheAccountKeepsEncodingJSONWireFormat(t *testing.T) {
	cases := []struct {
		name    string
		account service.Account
placeholder{
		{name: "nil collections", account: service.Account{ID: 801placeholderplaceholder,
		{name: "empty collections", account: service.Account{
			ID:          802,
	placeholderplaceholder,
			Extra:       map[string]any{placeholder,
			GroupIDs:    []int64{placeholder,
			Groups:      []*service.Group{placeholder,
placeholder
		{name: "nested maps and escaping", account: service.Account{
			ID:          803,
	placeholder"model_mapping": map[string]any{"z": "<last>", "a": "&first"placeholderplaceholder,
			Extra:       map[string]any{"mixed_scheduling": trueplaceholder,
placeholder
placeholder
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			full, meta, err := marshalSchedulerCacheAccount(tc.account)
		placeholder
			wantFull, err := json.Marshal(tc.account)
		placeholder
			wantMeta, err := json.Marshal(buildSchedulerMetadataAccount(tc.account))
		placeholder
			require.Equal(t, wantFull, full)
			require.Equal(t, wantMeta, meta)
	placeholder)
placeholder
placeholder

func TestBuildSchedulerMetadataAccount_KeepsOpenAIWSFlags(t *testing.T) {
	account := service.Account{
		ID:       42,
		Platform: service.PlatformOpenAI,
		Type:     service.AccountTypeOAuth,
		Extra: map[string]any{
			"openai_oauth_responses_websockets_v2_enabled": true,
			"openai_oauth_responses_websockets_v2_mode":    service.OpenAIWSIngressModePassthrough,
			"openai_ws_force_http":                         true,
			"openai_responses_mode":                        "force_chat_completions",
			"openai_responses_supported":                   false,
			"mixed_scheduling":                             true,
			"unused_large_field":                           "drop-me",
	placeholder,
placeholder

	got := buildSchedulerMetadataAccount(account)

	require.Equal(t, true, got.Extra["openai_oauth_responses_websockets_v2_enabled"])
	require.Equal(t, service.OpenAIWSIngressModePassthrough, got.Extra["openai_oauth_responses_websockets_v2_mode"])
	require.Equal(t, true, got.Extra["openai_ws_force_http"])
	require.Equal(t, "force_chat_completions", got.Extra["openai_responses_mode"])
	require.Equal(t, false, got.Extra["openai_responses_supported"])
	require.Equal(t, true, got.Extra["mixed_scheduling"])
	require.Nil(t, got.Extra["unused_large_field"])
placeholder

func TestBuildSchedulerMetadataAccount_KeepsGrokMediaEligibility(t *testing.T) {
	t.Run("explicit override", func(t *testing.T) {
		account := service.Account{
			ID:       43,
			Platform: service.PlatformGrok,
			Type:     service.AccountTypeOAuth,
			Extra: map[string]any{
				service.GrokMediaEligibleExtraKey: false,
				"unused_large_field":              "drop-me",
		placeholder,
	placeholder

		got := buildSchedulerMetadataAccount(account)

		eligible, reason := got.GrokMediaGenerationEligibility()
		require.False(t, eligible)
		require.Equal(t, "override_disabled", reason)
		require.Equal(t, false, got.Extra[service.GrokMediaEligibleExtraKey])
		require.Nil(t, got.Extra["unused_large_field"])
placeholder)

	t.Run("forbidden billing observation", func(t *testing.T) {
		account := service.Account{
			ID:       44,
			Platform: service.PlatformGrok,
			Type:     service.AccountTypeOAuth,
			Extra: map[string]any{
				"grok_billing_snapshot": map[string]any{
					"status_code":         200,
					"weekly_status_code":  403,
					"monthly_status_code": 200,
			placeholder,
		placeholder,
	placeholder

		got := buildSchedulerMetadataAccount(account)

		eligible, reason := got.GrokMediaGenerationEligibility()
		require.False(t, eligible)
		require.Equal(t, "billing_forbidden", reason)
		require.NotNil(t, got.Extra["grok_billing_snapshot"])
placeholder)
placeholder

func TestBuildSchedulerMetadataAccount_KeepsSlimGroupMembership(t *testing.T) {
	account := service.Account{
		ID:       42,
		Platform: service.PlatformAnthropic,
		GroupIDs: []int64{7, 9, 7, 0placeholder,
		AccountGroups: []service.AccountGroup{
			{
				AccountID: 42,
				GroupID:   7,
				Priority:  2,
				Account:   &service.Account{ID: 42, Name: "drop-from-metadata"placeholder,
				Group:     &service.Group{ID: 7, Name: "drop-from-metadata"placeholder,
		placeholder,
			{
				AccountID: 42,
				GroupID:   11,
				Priority:  3,
				Group:     &service.Group{ID: 11, Name: "drop-from-metadata"placeholder,
		placeholder,
			{
				AccountID: 42,
				GroupID:   0,
				Priority:  4,
		placeholder,
	placeholder,
placeholder

	got := buildSchedulerMetadataAccount(account)

	require.Equal(t, []int64{7, 9, 11placeholder, got.GroupIDs)
	require.Len(t, got.AccountGroups, 2)
	require.Equal(t, int64(42), got.AccountGroups[0].AccountID)
	require.Equal(t, int64(7), got.AccountGroups[0].GroupID)
	require.Equal(t, 2, got.AccountGroups[0].Priority)
	require.Nil(t, got.AccountGroups[0].Account)
	require.Nil(t, got.AccountGroups[0].Group)
	require.Equal(t, int64(11), got.AccountGroups[1].GroupID)
	require.Nil(t, got.Groups)
placeholder

func TestBuildSchedulerMetadataAccount_KeepsQuotaAutoPauseFields(t *testing.T) {
	account := service.Account{
		ID: 88,
		Extra: map[string]any{
			"codex_5h_used_percent":        12.34,
			"codex_7d_used_percent":        56.78,
			"codex_5h_reset_at":            "2026-05-29T10:00:00Z",
			"codex_7d_reset_at":            "2026-06-01T10:00:00Z",
			"codex_5h_reset_after_seconds": 300,
			"codex_7d_reset_after_seconds": 600,
			"codex_usage_updated_at":       "2026-05-29T09:00:00Z",
			"auto_pause_5h_threshold":      0.95,
			"auto_pause_7d_threshold":      0.96,
			"auto_pause_5h_disabled":       true,
			"auto_pause_7d_disabled":       false,
	placeholder,
placeholder

	got := buildSchedulerMetadataAccount(account)

	require.Equal(t, 12.34, got.Extra["codex_5h_used_percent"])
	require.Equal(t, 56.78, got.Extra["codex_7d_used_percent"])
	require.Equal(t, "2026-05-29T10:00:00Z", got.Extra["codex_5h_reset_at"])
	require.Equal(t, "2026-06-01T10:00:00Z", got.Extra["codex_7d_reset_at"])
	require.Equal(t, 300, got.Extra["codex_5h_reset_after_seconds"])
	require.Equal(t, 600, got.Extra["codex_7d_reset_after_seconds"])
	require.Equal(t, "2026-05-29T09:00:00Z", got.Extra["codex_usage_updated_at"])
	require.Equal(t, 0.95, got.Extra["auto_pause_5h_threshold"])
	require.Equal(t, 0.96, got.Extra["auto_pause_7d_threshold"])
	require.Equal(t, true, got.Extra["auto_pause_5h_disabled"])
	require.Equal(t, false, got.Extra["auto_pause_7d_disabled"])
placeholder

func TestBuildSchedulerMetadataAccount_KeepsModelRateLimits(t *testing.T) {
	account := service.Account{
		ID:       90,
		Platform: service.PlatformAntigravity,
		Extra: map[string]any{
			"model_rate_limits": map[string]any{
				"gemini-3-flash": map[string]any{
					"rate_limit_reset_at": "2026-05-30T10:10:00Z",
			placeholder,
				"antigravity:gemini": map[string]any{
					"rate_limit_reset_at": "2026-05-30T10:10:00Z",
			placeholder,
		placeholder,
			"unused_large_field": "drop-me",
	placeholder,
placeholder

	got := buildSchedulerMetadataAccount(account)

	limits, ok := got.Extra["model_rate_limits"].(map[string]any)
	require.True(t, ok)
	require.Contains(t, limits, "gemini-3-flash")
	require.Contains(t, limits, "antigravity:gemini")
	require.Nil(t, got.Extra["unused_large_field"])
placeholder

func TestBuildSchedulerMetadataAccount_KeepsSparkShadowRoutingIdentity(t *testing.T) {
	parentID := int64(100)
	account := service.Account{
		ID:              200,
		Platform:        service.PlatformOpenAI,
		Type:            service.AccountTypeOAuth,
		ParentAccountID: &parentID,
		QuotaDimension:  service.QuotaDimensionSpark,
placeholder
			"model_mapping": map[string]any{
				"gpt-5.3-codex-spark": "gpt-5.3-codex-spark",
		placeholder,
			"compact_model_mapping": map[string]any{
				"gpt-5.4": "gpt-5.4-openai-compact",
		placeholder,
			"access_token": "drop-me",
	placeholder,
placeholder

	got := buildSchedulerMetadataAccount(account)

	require.NotNil(t, got.ParentAccountID)
	require.Equal(t, parentID, *got.ParentAccountID)
	require.Equal(t, service.QuotaDimensionSpark, got.QuotaDimension)
	require.Equal(t, map[string]any{"gpt-5.3-codex-spark": "gpt-5.3-codex-spark"placeholder, got.Credentials["model_mapping"])
	require.Equal(t, map[string]any{"gpt-5.4": "gpt-5.4-openai-compact"placeholder, got.Credentials["compact_model_mapping"])
	require.Nil(t, got.Credentials["access_token"])
placeholder

func TestSchedulerCacheBucketRetirementFencesWritersAndReopen(t *testing.T) {
	ctx := context.Background()
	cache, mr := newSchedulerCacheUnitWithRedis(t)
	bucket := service.SchedulerBucket{GroupID: 41, Platform: service.PlatformOpenAI, Mode: service.SchedulerModeSingleplaceholder
	otherBucket := service.SchedulerBucket{GroupID: 42, Platform: service.PlatformOpenAI, Mode: service.SchedulerModeSingleplaceholder
	account := service.Account{ID: 4101, Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKeyplaceholder

	token, err := cache.CaptureBucketWriteToken(ctx, bucket)
placeholder
	require.True(t, token.ValidFor(bucket))
	require.NoError(t, cache.SetSnapshot(ctx, bucket, token, []service.Account{accountplaceholder))

	// A token is bound to the full bucket identity, not just an epoch number.
	err = cache.SetSnapshot(ctx, otherBucket, token, []service.Account{accountplaceholder)
	require.ErrorIs(t, err, service.ErrSchedulerBucketWriteFenced)
	_, err = cache.rdb.Get(ctx, schedulerBucketKey(schedulerVersionPrefix, otherBucket)).Result()
	require.ErrorIs(t, err, redis.Nil)
	otherAccount := service.Account{ID: 4201, Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKeyplaceholder
	otherToken, err := cache.CaptureBucketWriteToken(ctx, otherBucket)
placeholder
	require.NoError(t, cache.SetSnapshot(ctx, otherBucket, otherToken, []service.Account{otherAccountplaceholder))
	otherEpoch := otherToken.Epoch

	activeVersion, err := cache.rdb.Get(ctx, schedulerBucketKey(schedulerActivePrefix, bucket)).Result()
placeholder
	require.NoError(t, cache.RetireBucket(ctx, bucket))
	retiredEpoch, err := cache.rdb.Get(ctx, schedulerBucketKey(schedulerEpochPrefix, bucket)).Int64()
placeholder
	require.Greater(t, retiredEpoch, token.Epoch)

	// Retirement is idempotent and does not advance the epoch again.
	require.NoError(t, cache.RetireBucket(ctx, bucket))
	retiredEpochAgain, err := cache.rdb.Get(ctx, schedulerBucketKey(schedulerEpochPrefix, bucket)).Int64()
placeholder
	require.Equal(t, retiredEpoch, retiredEpochAgain)

	// New readers miss because ready/active were removed atomically. A reader that
	// captured activeVersion before retirement may still finish against that version.
	_, hit, err := cache.GetSnapshot(ctx, bucket)
placeholder
	require.False(t, hit)
	ids, err := cache.rdb.ZRange(ctx, schedulerSnapshotKey(bucket, activeVersion), 0, -1).Result()
placeholder
	require.Equal(t, []string{"4101"placeholder, ids)
	ttl, err := cache.rdb.TTL(ctx, schedulerSnapshotKey(bucket, activeVersion)).Result()
placeholder
	require.Positive(t, ttl)
	require.LessOrEqual(t, ttl, time.Duration(snapshotGraceTTLSeconds)*time.Second)

	buckets, err := cache.ListBuckets(ctx)
placeholder
	require.NotContains(t, buckets, bucket)
	require.Contains(t, buckets, otherBucket)
	otherSnapshot, otherHit, err := cache.GetSnapshot(ctx, otherBucket)
placeholder
	require.True(t, otherHit)
	require.Len(t, otherSnapshot, 1)
	require.Equal(t, otherAccount.ID, otherSnapshot[0].ID)
	otherEpochAfter, err := cache.rdb.Get(ctx, schedulerBucketKey(schedulerEpochPrefix, otherBucket)).Int64()
placeholder
	require.Equal(t, otherEpoch, otherEpochAfter)

	_, err = cache.CaptureBucketWriteToken(ctx, bucket)
	require.ErrorIs(t, err, service.ErrSchedulerBucketRetired)
	versionBeforeRejectedWrite, err := cache.rdb.Get(ctx, schedulerBucketKey(schedulerVersionPrefix, bucket)).Int64()
placeholder
	err = cache.SetSnapshot(ctx, bucket, token, []service.Account{accountplaceholder)
	require.ErrorIs(t, err, service.ErrSchedulerBucketRetired)
	versionAfterRejectedWrite, err := cache.rdb.Get(ctx, schedulerBucketKey(schedulerVersionPrefix, bucket)).Int64()
placeholder
	require.Equal(t, versionBeforeRejectedWrite, versionAfterRejectedWrite, "fenced writers must not allocate a new version")
	retired, err := cache.rdb.Exists(ctx, schedulerBucketKey(schedulerRetiredPrefix, bucket)).Result()
placeholder
	require.EqualValues(t, 1, retired, "ordinary writers must never clear the tombstone")
	mr.FastForward(time.Duration(snapshotGraceTTLSeconds+1) * time.Second)
	exists, err := cache.rdb.Exists(ctx, schedulerSnapshotKey(bucket, activeVersion)).Result()
placeholder
	require.Zero(t, exists, "retired active snapshot must expire after the in-flight grace period")

	newToken, err := cache.ReopenBucket(ctx, bucket)
placeholder
	require.True(t, newToken.ValidFor(bucket))
	require.Equal(t, retiredEpoch, newToken.Epoch)
	reopenedAgain, err := cache.ReopenBucket(ctx, bucket)
placeholder
	require.Equal(t, newToken, reopenedAgain, "reopen must be idempotent within one retirement generation")
	err = cache.SetSnapshot(ctx, bucket, token, []service.Account{accountplaceholder)
	require.ErrorIs(t, err, service.ErrSchedulerBucketWriteFenced)
	require.NoError(t, cache.SetSnapshot(ctx, bucket, newToken, []service.Account{accountplaceholder))
	reopenedWhileOpen, err := cache.ReopenBucket(ctx, bucket)
placeholder
	require.Equal(t, newToken, reopenedWhileOpen)

	snapshot, hit, err := cache.GetSnapshot(ctx, bucket)
placeholder
	require.True(t, hit)
	require.Len(t, snapshot, 1)
	require.Equal(t, account.ID, snapshot[0].ID)
placeholder

func TestSchedulerCacheActivationIsFencedAfterRetire(t *testing.T) {
	ctx := context.Background()
	cache := newSchedulerCacheUnit(t)
	bucket := service.SchedulerBucket{GroupID: 51, Platform: service.PlatformAnthropic, Mode: service.SchedulerModeMixedplaceholder
	account := service.Account{ID: 5101, Platform: service.PlatformAnthropic, Type: service.AccountTypeAPIKeyplaceholder

	token, err := cache.CaptureBucketWriteToken(ctx, bucket)
placeholder
	version, err := cache.allocateSnapshotVersion(ctx, bucket, token)
placeholder
	_, err = cache.writeSnapshotVersionAndReturnAccountIDs(ctx, bucket, version, []service.Account{accountplaceholder)
placeholder

	// Deterministic race C: retirement and authoritative reopen both happen after
	// INCR/write but before the old writer activates.
	require.NoError(t, cache.RetireBucket(ctx, bucket))
	_, err = cache.ReopenBucket(ctx, bucket)
placeholder
	err = cache.activateSnapshotVersion(ctx, bucket, token, version)
	require.ErrorIs(t, err, service.ErrSchedulerBucketWriteFenced)

	exists, err := cache.rdb.Exists(ctx, schedulerSnapshotKey(bucket, version)).Result()
placeholder
	require.Zero(t, exists, "fenced activation must delete its unpublished snapshot")
	exists, err = cache.rdb.Exists(
		ctx,
		schedulerBucketKey(schedulerReadyPrefix, bucket),
		schedulerBucketKey(schedulerActivePrefix, bucket),
	).Result()
placeholder
	require.Zero(t, exists)
	buckets, err := cache.ListBuckets(ctx)
placeholder
	require.NotContains(t, buckets, bucket)
placeholder

func TestSchedulerCacheConcurrentReopenReturnsSameToken(t *testing.T) {
	ctx := context.Background()
	cache := newSchedulerCacheUnit(t)
	bucket := service.SchedulerBucket{GroupID: 53, Platform: service.PlatformOpenAI, Mode: service.SchedulerModeForcedplaceholder
	account := service.Account{ID: 5301, Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKeyplaceholder

	oldToken, err := cache.CaptureBucketWriteToken(ctx, bucket)
placeholder
	require.NoError(t, cache.RetireBucket(ctx, bucket))

	type reopenResult struct {
		token service.SchedulerBucketWriteToken
		err   error
placeholder
	start := make(chan struct{placeholder)
	results := make(chan reopenResult, 2)
	for range 2 {
		go func() {
			<-start
			token, err := cache.ReopenBucket(ctx, bucket)
			results <- reopenResult{token: token, err: errplaceholder
	placeholder()
placeholder
	close(start)
	first := <-results
	second := <-results
	require.NoError(t, first.err)
	require.NoError(t, second.err)
	require.Equal(t, first.token, second.token)
	require.Greater(t, first.token.Epoch, oldToken.Epoch)

	require.ErrorIs(t, cache.SetSnapshot(ctx, bucket, oldToken, []service.Account{accountplaceholder), service.ErrSchedulerBucketWriteFenced)
	require.NoError(t, cache.SetSnapshot(ctx, bucket, first.token, []service.Account{accountplaceholder))
placeholder

func TestSchedulerCacheReopenExpiresPreviousActiveSnapshot(t *testing.T) {
	ctx := context.Background()
	cache, mr := newSchedulerCacheUnitWithRedis(t)
	bucket := service.SchedulerBucket{GroupID: 52, Platform: service.PlatformGemini, Mode: service.SchedulerModeForcedplaceholder
	account := service.Account{ID: 5201, Platform: service.PlatformGemini, Type: service.AccountTypeAPIKeyplaceholder

	oldToken, err := cache.CaptureBucketWriteToken(ctx, bucket)
placeholder
	require.NoError(t, cache.SetSnapshot(ctx, bucket, oldToken, []service.Account{accountplaceholder))
	oldVersion, err := cache.rdb.Get(ctx, schedulerBucketKey(schedulerActivePrefix, bucket)).Result()
placeholder
	retiredEpoch := oldToken.Epoch + 1
	require.NoError(t, cache.rdb.Set(ctx, schedulerBucketKey(schedulerEpochPrefix, bucket), retiredEpoch, 0).Err())
	require.NoError(t, cache.rdb.Set(ctx, schedulerBucketKey(schedulerRetiredPrefix, bucket), retiredEpoch, 0).Err())

	newToken, err := cache.ReopenBucket(ctx, bucket)
placeholder
	require.Equal(t, retiredEpoch, newToken.Epoch)
	_, hit, err := cache.GetSnapshot(ctx, bucket)
placeholder
	require.False(t, hit)
	ttl, err := cache.rdb.TTL(ctx, schedulerSnapshotKey(bucket, oldVersion)).Result()
placeholder
	require.Positive(t, ttl)
	require.LessOrEqual(t, ttl, time.Duration(snapshotGraceTTLSeconds)*time.Second)

	require.ErrorIs(t, cache.SetSnapshot(ctx, bucket, oldToken, []service.Account{accountplaceholder), service.ErrSchedulerBucketWriteFenced)
	mr.FastForward(time.Duration(snapshotGraceTTLSeconds+1) * time.Second)
	exists, err := cache.rdb.Exists(ctx, schedulerSnapshotKey(bucket, oldVersion)).Result()
placeholder
	require.Zero(t, exists)
	require.NoError(t, cache.SetSnapshot(ctx, bucket, newToken, []service.Account{accountplaceholder))
placeholder

func TestSchedulerCacheGroupLifecycleLeaseConcurrentAcquireSingleOwner(t *testing.T) {
	ctx := context.Background()
	cache := newSchedulerCacheUnit(t)
	const groupID int64 = 71

	type result struct {
		lease    service.SchedulerGroupLifecycleLease
		acquired bool
		err      error
placeholder
	start := make(chan struct{placeholder)
	results := make(chan result, 32)
	for range 32 {
		go func() {
			<-start
			lease, acquired, err := cache.TryAcquireGroupLifecycleLease(ctx, groupID, time.Minute)
			results <- result{lease: lease, acquired: acquired, err: errplaceholder
	placeholder()
placeholder
	close(start)

	var owner service.SchedulerGroupLifecycleLease
	acquiredCount := 0
	for range 32 {
		got := <-results
		require.NoError(t, got.err)
		if got.acquired {
			acquiredCount++
			owner = got.lease
			require.True(t, got.lease.ValidFor(groupID))
	placeholder else {
			require.Equal(t, service.SchedulerGroupLifecycleLease{placeholder, got.lease)
	placeholder
placeholder
	require.Equal(t, 1, acquiredCount)
	require.Len(t, owner.OwnerToken, schedulerGroupLifecycleOwnerTokenBytes*2)
	require.Equal(t, strings.ToLower(owner.OwnerToken), owner.OwnerToken)
	decodedOwner, err := hex.DecodeString(owner.OwnerToken)
placeholder
	require.Len(t, decodedOwner, schedulerGroupLifecycleOwnerTokenBytes)

	require.NoError(t, cache.ReleaseGroupLifecycleLease(ctx, owner))
	next, acquired, err := cache.TryAcquireGroupLifecycleLease(ctx, groupID, time.Minute)
placeholder
	require.True(t, acquired)
	require.True(t, next.ValidFor(groupID))
	require.NotEqual(t, owner.OwnerToken, next.OwnerToken)
	require.NoError(t, cache.ReleaseGroupLifecycleLease(ctx, next))
placeholder

func TestSchedulerCacheGroupLifecycleLeaseStaleReleaseCannotDeleteSuccessor(t *testing.T) {
	ctx := context.Background()
	cache, mr := newSchedulerCacheUnitWithRedis(t)
	const groupID int64 = 72
	const ttl = time.Minute

	first, acquired, err := cache.TryAcquireGroupLifecycleLease(ctx, groupID, ttl)
placeholder
	require.True(t, acquired)

	mr.FastForward(ttl + time.Second)
	second, acquired, err := cache.TryAcquireGroupLifecycleLease(ctx, groupID, ttl)
placeholder
	require.True(t, acquired)
	require.NotEqual(t, first.OwnerToken, second.OwnerToken)

	require.ErrorIs(t, cache.ReleaseGroupLifecycleLease(ctx, first), service.ErrSchedulerGroupLifecycleLeaseLost)
	owner, err := cache.rdb.Get(ctx, schedulerGroupLifecycleLockKey(groupID)).Result()
placeholder
	require.Equal(t, second.OwnerToken, owner)

	_, acquired, err = cache.TryAcquireGroupLifecycleLease(ctx, groupID, ttl)
placeholder
	require.False(t, acquired)
	require.NoError(t, cache.ReleaseGroupLifecycleLease(ctx, second))
	require.ErrorIs(t, cache.ReleaseGroupLifecycleLease(ctx, second), service.ErrSchedulerGroupLifecycleLeaseLost)
placeholder

func TestSchedulerCacheGroupLifecycleLeaseExpiredReleaseIsLost(t *testing.T) {
	ctx := context.Background()
	cache, mr := newSchedulerCacheUnitWithRedis(t)
	const groupID int64 = 73
	const ttl = time.Minute

	lease, acquired, err := cache.TryAcquireGroupLifecycleLease(ctx, groupID, ttl)
placeholder
	require.True(t, acquired)
	mr.FastForward(ttl + time.Second)

	require.ErrorIs(t, cache.ReleaseGroupLifecycleLease(ctx, lease), service.ErrSchedulerGroupLifecycleLeaseLost)
placeholder

func TestSchedulerCacheGroupLifecycleLeaseWrongOwnerAndCrossGroupAreLost(t *testing.T) {
	ctx := context.Background()
	cache := newSchedulerCacheUnit(t)
	const firstGroupID int64 = 74
	const secondGroupID int64 = 75

	first, acquired, err := cache.TryAcquireGroupLifecycleLease(ctx, firstGroupID, time.Minute)
placeholder
	require.True(t, acquired)
	second, acquired, err := cache.TryAcquireGroupLifecycleLease(ctx, secondGroupID, time.Minute)
placeholder
	require.True(t, acquired, "different groups must acquire independently")
	require.NotEqual(t, first.OwnerToken, second.OwnerToken)

	wrongOwner := first
	wrongOwner.OwnerToken = strings.Repeat("0", schedulerGroupLifecycleOwnerTokenBytes*2)
	if wrongOwner.OwnerToken == first.OwnerToken {
		wrongOwner.OwnerToken = strings.Repeat("1", schedulerGroupLifecycleOwnerTokenBytes*2)
placeholder
	require.ErrorIs(t, cache.ReleaseGroupLifecycleLease(ctx, wrongOwner), service.ErrSchedulerGroupLifecycleLeaseLost)

	crossGroup := first
	crossGroup.GroupID = secondGroupID
	require.ErrorIs(t, cache.ReleaseGroupLifecycleLease(ctx, crossGroup), service.ErrSchedulerGroupLifecycleLeaseLost)

	firstOwner, err := cache.rdb.Get(ctx, schedulerGroupLifecycleLockKey(firstGroupID)).Result()
placeholder
	require.Equal(t, first.OwnerToken, firstOwner)
	secondOwner, err := cache.rdb.Get(ctx, schedulerGroupLifecycleLockKey(secondGroupID)).Result()
placeholder
	require.Equal(t, second.OwnerToken, secondOwner)

	require.NoError(t, cache.ReleaseGroupLifecycleLease(ctx, first))
	require.NoError(t, cache.ReleaseGroupLifecycleLease(ctx, second))
placeholder

func TestSchedulerCacheGroupLifecycleLeaseCanceledContextFailsClosed(t *testing.T) {
	cache := newSchedulerCacheUnit(t)
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	lease, acquired, err := cache.TryAcquireGroupLifecycleLease(canceledCtx, 76, time.Minute)
	require.ErrorIs(t, err, context.Canceled)
	require.False(t, acquired)
	require.Equal(t, service.SchedulerGroupLifecycleLease{placeholder, lease)

	ctx := context.Background()
	lease, acquired, err = cache.TryAcquireGroupLifecycleLease(ctx, 76, time.Minute)
placeholder
	require.True(t, acquired)
	require.ErrorIs(t, cache.ReleaseGroupLifecycleLease(canceledCtx, lease), context.Canceled)
	owner, err := cache.rdb.Get(ctx, schedulerGroupLifecycleLockKey(lease.GroupID)).Result()
placeholder
	require.Equal(t, lease.OwnerToken, owner)
	require.NoError(t, cache.ReleaseGroupLifecycleLease(ctx, lease))
placeholder

func TestSchedulerCacheGroupLifecycleLeaseRejectsInvalidInput(t *testing.T) {
	ctx := context.Background()
	cache := newSchedulerCacheUnit(t)

	lease, acquired, err := cache.TryAcquireGroupLifecycleLease(ctx, 0, time.Minute)
	require.ErrorIs(t, err, service.ErrSchedulerGroupLifecycleLeaseInvalid)
	require.False(t, acquired)
	require.Equal(t, service.SchedulerGroupLifecycleLease{placeholder, lease)

	lease, acquired, err = cache.TryAcquireGroupLifecycleLease(ctx, 73, 0)
	require.ErrorIs(t, err, service.ErrSchedulerGroupLifecycleLeaseInvalid)
	require.False(t, acquired)
	require.Equal(t, service.SchedulerGroupLifecycleLease{placeholder, lease)

	canceledCtx, cancel := context.WithCancel(ctx)
	cancel()
	require.ErrorIs(t, cache.ReleaseGroupLifecycleLease(canceledCtx, service.SchedulerGroupLifecycleLease{placeholder), service.ErrSchedulerGroupLifecycleLeaseInvalid)
	keys, err := cache.rdb.DBSize(ctx).Result()
placeholder
	require.Zero(t, keys)
placeholder

var schedulerCachePayloadBenchmarkSink int

func BenchmarkSchedulerCacheAccountPayloadReuse(b *testing.B) {
	for _, size := range []int{1, 100, 10_000placeholder {
		accounts := schedulerCacheBenchmarkAccounts(size)
		b.Run(fmt.Sprintf("pair_baseline_%d_accounts", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				first, err := benchmarkSchedulerLegacySnapshotPayload(accounts)
				if err != nil {
					b.Fatal(err)
			placeholder
				second, err := benchmarkSchedulerLegacySnapshotPayload(accounts)
				if err != nil {
					b.Fatal(err)
			placeholder
				schedulerCachePayloadBenchmarkSink = first + second
		placeholder
	placeholder)
		b.Run(fmt.Sprintf("pair_reuse_%d_accounts", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				ids, total, err := benchmarkSchedulerReusableSnapshotPayload(accounts)
				if err != nil {
					b.Fatal(err)
			placeholder
				// 第二个桶仍构造成员，只跳过账号 JSON 与全局账号键。
				total += len(schedulerSnapshotMembers(ids))
				schedulerCachePayloadBenchmarkSink = total
		placeholder
	placeholder)
		b.Run(fmt.Sprintf("first_baseline_%d_accounts", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				total, err := benchmarkSchedulerLegacySnapshotPayload(accounts)
				if err != nil {
					b.Fatal(err)
			placeholder
				schedulerCachePayloadBenchmarkSink = total
		placeholder
	placeholder)
		b.Run(fmt.Sprintf("first_reuse_%d_accounts", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				ids, total, err := benchmarkSchedulerReusableSnapshotPayload(accounts)
				if err != nil {
					b.Fatal(err)
			placeholder
				total += len(ids)
				schedulerCachePayloadBenchmarkSink = total
		placeholder
	placeholder)
placeholder
placeholder

func benchmarkSchedulerLegacySnapshotPayload(accounts []service.Account) (int, error) {
	cacheable := make([]service.Account, 0, len(accounts))
	total := 0
	for _, account := range accounts {
		full, meta, err := marshalSchedulerCacheAccount(account)
		if err != nil {
			continue
	placeholder
		total += len(full) + len(meta)
		cacheable = append(cacheable, account)
placeholder
	members := make([]redis.Z, 0, len(cacheable))
	for idx, account := range cacheable {
		members = append(members, redis.Z{Score: float64(idx), Member: strconv.FormatInt(account.ID, 10)placeholder)
placeholder
	return total + len(members), nil
placeholder

func benchmarkSchedulerReusableSnapshotPayload(accounts []service.Account) ([]int64, int, error) {
	accountIDs := make([]int64, 0, len(accounts))
	total := 0
	for _, account := range accounts {
		full, meta, err := marshalSchedulerCacheAccount(account)
		if err != nil {
			continue
	placeholder
		total += len(full) + len(meta)
		accountIDs = append(accountIDs, account.ID)
placeholder
	total += len(schedulerSnapshotMembers(accountIDs))
	return accountIDs, total, nil
placeholder

func schedulerCacheBenchmarkAccounts(size int) []service.Account {
	largeValue := strings.Repeat("x", 4096)
	credentials := map[string]any{
		"api_key":       "benchmark-key",
		"model_mapping": map[string]any{"z-model": "z-target", "a-model": "a-target"placeholder,
		"large_value":   largeValue,
placeholder
	extra := map[string]any{
		"mixed_scheduling": true,
		"model_rate_limits": map[string]any{
			"z-model": map[string]any{"rate_limit_reset_at": "2026-07-16T00:00:00Z"placeholder,
			"a-model": map[string]any{"rate_limit_reset_at": "2026-07-16T00:00:00Z"placeholder,
	placeholder,
		"large_value": largeValue,
placeholder
	accounts := make([]service.Account, size)
	for i := range accounts {
		id := int64(i + 1)
		accounts[i] = service.Account{
			ID:          id,
			Name:        "benchmark-account",
			Platform:    service.PlatformOpenAI,
			Type:        service.AccountTypeOAuth,
			Credentials: credentials,
			Extra:       extra,
			GroupIDs:    []int64{7, 9placeholder,
			AccountGroups: []service.AccountGroup{
				{AccountID: id, GroupID: 7, Priority: 1placeholder,
				{AccountID: id, GroupID: 9, Priority: 2placeholder,
		placeholder,
	placeholder
placeholder
	return accounts
placeholder
