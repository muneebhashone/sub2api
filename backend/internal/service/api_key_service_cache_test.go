//go:build unit

package service

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

type authRepoStub struct {
	getByKeyForAuth   func(ctx context.Context, key string) (*APIKey, error)
	listKeysByUserID  func(ctx context.Context, userID int64) ([]string, error)
	listKeysByGroupID func(ctx context.Context, groupID int64) ([]string, error)
placeholder

func (s *authRepoStub) Create(ctx context.Context, key *APIKey) error {
	panic("unexpected Create call")
placeholder

func (s *authRepoStub) GetByID(ctx context.Context, id int64) (*APIKey, error) {
	panic("unexpected GetByID call")
placeholder

func (s *authRepoStub) GetKeyAndOwnerID(ctx context.Context, id int64) (string, int64, error) {
	panic("unexpected GetKeyAndOwnerID call")
placeholder

func (s *authRepoStub) GetByKey(ctx context.Context, key string) (*APIKey, error) {
	panic("unexpected GetByKey call")
placeholder

func (s *authRepoStub) GetByKeyForAuth(ctx context.Context, key string) (*APIKey, error) {
	if s.getByKeyForAuth == nil {
		panic("unexpected GetByKeyForAuth call")
placeholder
	return s.getByKeyForAuth(ctx, key)
placeholder

func (s *authRepoStub) Update(ctx context.Context, key *APIKey) error {
	panic("unexpected Update call")
placeholder

func (s *authRepoStub) Delete(ctx context.Context, id int64) error {
	panic("unexpected Delete call")
placeholder

func (s *authRepoStub) DeleteWithAudit(ctx context.Context, id int64) error {
	panic("unexpected DeleteWithAudit call")
placeholder

func (s *authRepoStub) ListByUserID(ctx context.Context, userID int64, params pagination.PaginationParams, filters APIKeyListFilters) ([]APIKey, *pagination.PaginationResult, error) {
	panic("unexpected ListByUserID call")
placeholder

func (s *authRepoStub) VerifyOwnership(ctx context.Context, userID int64, apiKeyIDs []int64) ([]int64, error) {
	panic("unexpected VerifyOwnership call")
placeholder

func (s *authRepoStub) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	panic("unexpected CountByUserID call")
placeholder

func (s *authRepoStub) ExistsByKey(ctx context.Context, key string) (bool, error) {
	panic("unexpected ExistsByKey call")
placeholder

func (s *authRepoStub) ListByGroupID(ctx context.Context, groupID int64, params pagination.PaginationParams) ([]APIKey, *pagination.PaginationResult, error) {
	panic("unexpected ListByGroupID call")
placeholder

func (s *authRepoStub) SearchAPIKeys(ctx context.Context, userID int64, keyword string, limit int) ([]APIKey, error) {
	panic("unexpected SearchAPIKeys call")
placeholder

func (s *authRepoStub) ClearGroupIDByGroupID(ctx context.Context, groupID int64) (int64, error) {
	panic("unexpected ClearGroupIDByGroupID call")
placeholder
func (s *authRepoStub) UpdateGroupIDByUserAndGroup(ctx context.Context, userID, oldGroupID, newGroupID int64) (int64, error) {
	panic("unexpected UpdateGroupIDByUserAndGroup call")
placeholder

func (s *authRepoStub) CountByGroupID(ctx context.Context, groupID int64) (int64, error) {
	panic("unexpected CountByGroupID call")
placeholder

func (s *authRepoStub) ListKeysByUserID(ctx context.Context, userID int64) ([]string, error) {
	if s.listKeysByUserID == nil {
		panic("unexpected ListKeysByUserID call")
placeholder
	return s.listKeysByUserID(ctx, userID)
placeholder

func (s *authRepoStub) ListKeysByGroupID(ctx context.Context, groupID int64) ([]string, error) {
	if s.listKeysByGroupID == nil {
		panic("unexpected ListKeysByGroupID call")
placeholder
	return s.listKeysByGroupID(ctx, groupID)
placeholder

func (s *authRepoStub) IncrementQuotaUsed(ctx context.Context, id int64, amount float64) (float64, error) {
	panic("unexpected IncrementQuotaUsed call")
placeholder

func (s *authRepoStub) UpdateLastUsed(ctx context.Context, id int64, usedAt time.Time) error {
	panic("unexpected UpdateLastUsed call")
placeholder
func (s *authRepoStub) IncrementRateLimitUsage(ctx context.Context, id int64, cost float64) error {
	panic("unexpected IncrementRateLimitUsage call")
placeholder
func (s *authRepoStub) ResetRateLimitWindows(ctx context.Context, id int64) error {
	panic("unexpected ResetRateLimitWindows call")
placeholder
func (s *authRepoStub) GetRateLimitData(ctx context.Context, id int64) (*APIKeyRateLimitData, error) {
	panic("unexpected GetRateLimitData call")
placeholder

type authCacheStub struct {
	getAuthCache   func(ctx context.Context, key string) (*APIKeyAuthCacheEntry, error)
	setAuthKeys    []string
	deleteAuthKeys []string
placeholder

func (s *authCacheStub) GetCreateAttemptCount(ctx context.Context, userID int64) (int, error) {
	return 0, nil
placeholder

func (s *authCacheStub) IncrementCreateAttemptCount(ctx context.Context, userID int64) error {
	return nil
placeholder

func (s *authCacheStub) DeleteCreateAttemptCount(ctx context.Context, userID int64) error {
	return nil
placeholder

func (s *authCacheStub) IncrementDailyUsage(ctx context.Context, apiKey string) error {
	return nil
placeholder

func (s *authCacheStub) SetDailyUsageExpiry(ctx context.Context, apiKey string, ttl time.Duration) error {
	return nil
placeholder

func (s *authCacheStub) GetAuthCache(ctx context.Context, key string) (*APIKeyAuthCacheEntry, error) {
	if s.getAuthCache == nil {
		return nil, redis.Nil
placeholder
	return s.getAuthCache(ctx, key)
placeholder

func (s *authCacheStub) SetAuthCache(ctx context.Context, key string, entry *APIKeyAuthCacheEntry, ttl time.Duration) error {
	s.setAuthKeys = append(s.setAuthKeys, key)
	return nil
placeholder

func (s *authCacheStub) DeleteAuthCache(ctx context.Context, key string) error {
	s.deleteAuthKeys = append(s.deleteAuthKeys, key)
	return nil
placeholder

func (s *authCacheStub) PublishAuthCacheInvalidation(ctx context.Context, cacheKey string) error {
	return nil
placeholder

func (s *authCacheStub) SubscribeAuthCacheInvalidation(ctx context.Context, handler func(cacheKey string)) error {
	return nil
placeholder

func TestAPIKeyService_GetByKey_UsesL2Cache(t *testing.T) {
	cache := &authCacheStub{placeholder
	repo := &authRepoStub{
		getByKeyForAuth: func(ctx context.Context, key string) (*APIKey, error) {
			return nil, errors.New("unexpected repo call")
	placeholder,
placeholder
	cfg := &config.Config{
		APIKeyAuth: config.APIKeyAuthCacheConfig{
			L2TTLSeconds:       60,
			NegativeTTLSeconds: 30,
	placeholder,
placeholder
	svc := NewAPIKeyService(repo, nil, nil, nil, nil, cache, cfg)

	groupID := int64(9)
	cacheEntry := &APIKeyAuthCacheEntry{
		Snapshot: &APIKeyAuthSnapshot{
			Version:  apiKeyAuthSnapshotVersion,
			APIKeyID: 1,
			UserID:   2,
			GroupID:  &groupID,
			Status:   StatusActive,
			User: APIKeyAuthUserSnapshot{
				ID:          2,
				Status:      StatusActive,
				Role:        RoleUser,
				Balance:     10,
				Concurrency: 3,
		placeholder,
			Group: &APIKeyAuthGroupSnapshot{
				ID:                  groupID,
				Name:                "g",
				Platform:            PlatformAnthropic,
				Status:              StatusActive,
				SubscriptionType:    SubscriptionTypeStandard,
				RateMultiplier:      1,
				ModelRoutingEnabled: true,
				ModelRouting: map[string][]int64{
					"claude-opus-*": {1, 2placeholder,
			placeholder,
		placeholder,
	placeholder,
placeholder
	cache.getAuthCache = func(ctx context.Context, key string) (*APIKeyAuthCacheEntry, error) {
		return cacheEntry, nil
placeholder

	apiKey, err := svc.GetByKey(context.Background(), "k1")
placeholder
	require.Equal(t, int64(1), apiKey.ID)
	require.Equal(t, int64(2), apiKey.User.ID)
	require.Equal(t, groupID, apiKey.Group.ID)
	require.True(t, apiKey.Group.ModelRoutingEnabled)
	require.Equal(t, map[string][]int64{"claude-opus-*": {1, 2placeholderplaceholder, apiKey.Group.ModelRouting)
placeholder

func TestAPIKeyService_SnapshotRoundTrip_PreservesMessagesDispatchModelConfig(t *testing.T) {
	svc := NewAPIKeyService(nil, nil, nil, nil, nil, nil, &config.Config{placeholder)
	groupID := int64(9)
	apiKey := &APIKey{
		ID:      1,
		UserID:  2,
		GroupID: &groupID,
		Key:     "k-roundtrip",
		Name:    "Audit Key",
		Status:  StatusActive,
		User: &User{
			ID:          2,
			Status:      StatusActive,
			Role:        RoleUser,
			Balance:     10,
			Concurrency: 3,
	placeholder,
		Group: &Group{
			ID:                    groupID,
			Name:                  "openai",
			Platform:              PlatformOpenAI,
			Status:                StatusActive,
			SubscriptionType:      SubscriptionTypeStandard,
			RateMultiplier:        1,
			AllowMessagesDispatch: true,
			DefaultMappedModel:    "gpt-5.4",
			MessagesDispatchModelConfig: OpenAIMessagesDispatchModelConfig{
				OpusMappedModel:   "gpt-5.4-nano",
				SonnetMappedModel: "gpt-5.3-codex",
				HaikuMappedModel:  "gpt-5.4-mini",
				ExactModelMappings: map[string]string{
					"claude-sonnet-4.5": "gpt-5.4-nano",
			placeholder,
		placeholder,
	placeholder,
placeholder

	snapshot := svc.snapshotFromAPIKey(context.Background(), apiKey)
	roundTrip := svc.snapshotToAPIKey(apiKey.Key, snapshot)

	require.NotNil(t, roundTrip)
	require.Equal(t, apiKey.Name, roundTrip.Name)
	require.NotNil(t, roundTrip.Group)
	require.Equal(t, apiKey.Group.MessagesDispatchModelConfig, roundTrip.Group.MessagesDispatchModelConfig)
placeholder

func TestAPIKeyService_SnapshotRoundTrip_PreservesReasoningEffortPolicy(t *testing.T) {
	svc := NewAPIKeyService(nil, nil, nil, nil, nil, nil, &config.Config{placeholder)
	groupID := int64(9)
	apiKey := &APIKey{
		ID:      1,
		UserID:  2,
		GroupID: &groupID,
		Key:     "k-reasoning-policy",
		Status:  StatusActive,
		User: &User{
			ID:          2,
			Status:      StatusActive,
			Role:        RoleUser,
			Balance:     10,
			Concurrency: 3,
	placeholder,
		Group: &Group{
			ID:                 groupID,
			Name:               "openai",
			Platform:           PlatformOpenAI,
			Status:             StatusActive,
			SubscriptionType:   SubscriptionTypeStandard,
			RateMultiplier:     1,
			MaxReasoningEffort: "medium",
			ReasoningEffortMappings: []ReasoningEffortMapping{
				{From: "max", To: "xhigh"placeholder,
		placeholder,
	placeholder,
placeholder

	snapshot := svc.snapshotFromAPIKey(context.Background(), apiKey)
	roundTrip := svc.snapshotToAPIKey(apiKey.Key, snapshot)

	require.NotNil(t, roundTrip)
	require.NotNil(t, roundTrip.Group)
	require.Equal(t, "medium", roundTrip.Group.MaxReasoningEffort)
	require.Equal(t, apiKey.Group.ReasoningEffortMappings, roundTrip.Group.ReasoningEffortMappings)
placeholder

func TestAPIKeyService_GetByKey_IgnoresLegacyAuthCacheSnapshotWithoutMessagesDispatchConfig(t *testing.T) {
	cache := &authCacheStub{placeholder
	var repoCalls int32
	repo := &authRepoStub{
		getByKeyForAuth: func(ctx context.Context, key string) (*APIKey, error) {
			atomic.AddInt32(&repoCalls, 1)
			groupID := int64(9)
			return &APIKey{
				ID:      1,
				UserID:  2,
				GroupID: &groupID,
				Status:  StatusActive,
				User: &User{
					ID:          2,
					Status:      StatusActive,
					Role:        RoleUser,
					Balance:     10,
					Concurrency: 3,
			placeholder,
				Group: &Group{
					ID:                    groupID,
					Name:                  "openai",
					Platform:              PlatformOpenAI,
					Status:                StatusActive,
					Hydrated:              true,
					SubscriptionType:      SubscriptionTypeStandard,
					RateMultiplier:        1,
					AllowMessagesDispatch: true,
					DefaultMappedModel:    "gpt-5.4",
					MessagesDispatchModelConfig: OpenAIMessagesDispatchModelConfig{
						OpusMappedModel: "gpt-5.4-nano",
				placeholder,
			placeholder,
		placeholder, nil
	placeholder,
placeholder
	cfg := &config.Config{
		APIKeyAuth: config.APIKeyAuthCacheConfig{
			L2TTLSeconds: 60,
	placeholder,
placeholder
	svc := NewAPIKeyService(repo, nil, nil, nil, nil, cache, cfg)

	groupID := int64(9)
	cache.getAuthCache = func(ctx context.Context, key string) (*APIKeyAuthCacheEntry, error) {
		return &APIKeyAuthCacheEntry{
			Snapshot: &APIKeyAuthSnapshot{
				APIKeyID: 1,
				UserID:   2,
				GroupID:  &groupID,
				Status:   StatusActive,
				User: APIKeyAuthUserSnapshot{
					ID:          2,
					Status:      StatusActive,
					Role:        RoleUser,
					Balance:     10,
					Concurrency: 3,
			placeholder,
				Group: &APIKeyAuthGroupSnapshot{
					ID:                    groupID,
					Name:                  "openai",
					Platform:              PlatformOpenAI,
					Status:                StatusActive,
					SubscriptionType:      SubscriptionTypeStandard,
					RateMultiplier:        1,
					AllowMessagesDispatch: true,
					DefaultMappedModel:    "gpt-5.4",
			placeholder,
		placeholder,
	placeholder, nil
placeholder

	apiKey, err := svc.GetByKey(context.Background(), "k-legacy")
placeholder
	require.Equal(t, int32(1), atomic.LoadInt32(&repoCalls))
	require.NotNil(t, apiKey.Group)
	require.Equal(t, "gpt-5.4-nano", apiKey.Group.MessagesDispatchModelConfig.OpusMappedModel)
placeholder

func TestAPIKeyService_GetByKey_NegativeCache(t *testing.T) {
	cache := &authCacheStub{placeholder
	repo := &authRepoStub{
		getByKeyForAuth: func(ctx context.Context, key string) (*APIKey, error) {
			return nil, errors.New("unexpected repo call")
	placeholder,
placeholder
	cfg := &config.Config{
		APIKeyAuth: config.APIKeyAuthCacheConfig{
			L2TTLSeconds:       60,
			NegativeTTLSeconds: 30,
	placeholder,
placeholder
	svc := NewAPIKeyService(repo, nil, nil, nil, nil, cache, cfg)
	cache.getAuthCache = func(ctx context.Context, key string) (*APIKeyAuthCacheEntry, error) {
		return &APIKeyAuthCacheEntry{NotFound: trueplaceholder, nil
placeholder

	_, err := svc.GetByKey(context.Background(), "missing")
	require.ErrorIs(t, err, ErrAPIKeyNotFound)
placeholder

func TestAPIKeyService_GetByKey_CacheMissStoresL2(t *testing.T) {
	cache := &authCacheStub{placeholder
	repo := &authRepoStub{
		getByKeyForAuth: func(ctx context.Context, key string) (*APIKey, error) {
			return &APIKey{
				ID:     5,
				UserID: 7,
				Status: StatusActive,
				User: &User{
					ID:          7,
					Status:      StatusActive,
					Role:        RoleUser,
					Balance:     12,
					Concurrency: 2,
			placeholder,
		placeholder, nil
	placeholder,
placeholder
	cfg := &config.Config{
		APIKeyAuth: config.APIKeyAuthCacheConfig{
			L2TTLSeconds:       60,
			NegativeTTLSeconds: 30,
	placeholder,
placeholder
	svc := NewAPIKeyService(repo, nil, nil, nil, nil, cache, cfg)
	cache.getAuthCache = func(ctx context.Context, key string) (*APIKeyAuthCacheEntry, error) {
		return nil, redis.Nil
placeholder

	apiKey, err := svc.GetByKey(context.Background(), "k2")
placeholder
	require.Equal(t, int64(5), apiKey.ID)
	require.Len(t, cache.setAuthKeys, 1)
placeholder

func TestAPIKeyService_GetByKey_UsesL1Cache(t *testing.T) {
	var calls int32
	cache := &authCacheStub{placeholder
	repo := &authRepoStub{
		getByKeyForAuth: func(ctx context.Context, key string) (*APIKey, error) {
			atomic.AddInt32(&calls, 1)
			return &APIKey{
				ID:     21,
				UserID: 3,
				Status: StatusActive,
				User: &User{
					ID:          3,
					Status:      StatusActive,
					Role:        RoleUser,
					Balance:     5,
					Concurrency: 2,
			placeholder,
		placeholder, nil
	placeholder,
placeholder
	cfg := &config.Config{
		APIKeyAuth: config.APIKeyAuthCacheConfig{
			L1Size:       1000,
			L1TTLSeconds: 60,
	placeholder,
placeholder
	svc := NewAPIKeyService(repo, nil, nil, nil, nil, cache, cfg)
	require.NotNil(t, svc.authCacheL1)

	_, err := svc.GetByKey(context.Background(), "k-l1")
placeholder
	svc.authCacheL1.Wait()
	cacheKey := svc.authCacheKey("k-l1")
	_, ok := svc.authCacheL1.Get(cacheKey)
	require.True(t, ok)
	_, err = svc.GetByKey(context.Background(), "k-l1")
placeholder
	require.Equal(t, int32(1), atomic.LoadInt32(&calls))
placeholder

func TestAPIKeyService_InvalidateAuthCacheByUserID(t *testing.T) {
	cache := &authCacheStub{placeholder
	repo := &authRepoStub{
		listKeysByUserID: func(ctx context.Context, userID int64) ([]string, error) {
			return []string{"k1", "k2"placeholder, nil
	placeholder,
placeholder
	cfg := &config.Config{
		APIKeyAuth: config.APIKeyAuthCacheConfig{
			L2TTLSeconds:       60,
			NegativeTTLSeconds: 30,
	placeholder,
placeholder
	svc := NewAPIKeyService(repo, nil, nil, nil, nil, cache, cfg)

	svc.InvalidateAuthCacheByUserID(context.Background(), 7)
	require.Len(t, cache.deleteAuthKeys, 2)
placeholder

func TestAPIKeyService_InvalidateAuthCacheByGroupID(t *testing.T) {
	cache := &authCacheStub{placeholder
	repo := &authRepoStub{
		listKeysByGroupID: func(ctx context.Context, groupID int64) ([]string, error) {
			return []string{"k1", "k2"placeholder, nil
	placeholder,
placeholder
	cfg := &config.Config{
		APIKeyAuth: config.APIKeyAuthCacheConfig{
			L2TTLSeconds: 60,
	placeholder,
placeholder
	svc := NewAPIKeyService(repo, nil, nil, nil, nil, cache, cfg)

	svc.InvalidateAuthCacheByGroupID(context.Background(), 9)
	require.Len(t, cache.deleteAuthKeys, 2)
placeholder

func TestAPIKeyService_InvalidateAuthCacheByKey(t *testing.T) {
	cache := &authCacheStub{placeholder
	repo := &authRepoStub{
		listKeysByUserID: func(ctx context.Context, userID int64) ([]string, error) {
			return nil, nil
	placeholder,
placeholder
	cfg := &config.Config{
		APIKeyAuth: config.APIKeyAuthCacheConfig{
			L2TTLSeconds: 60,
	placeholder,
placeholder
	svc := NewAPIKeyService(repo, nil, nil, nil, nil, cache, cfg)

	svc.InvalidateAuthCacheByKey(context.Background(), "k1")
	require.Len(t, cache.deleteAuthKeys, 1)
placeholder

func TestAPIKeyService_GetByKey_CachesNegativeOnRepoMiss(t *testing.T) {
	var repoCalls atomic.Int32
	cache := &authCacheStub{placeholder
	repo := &authRepoStub{
		getByKeyForAuth: func(ctx context.Context, key string) (*APIKey, error) {
			repoCalls.Add(1)
			return nil, ErrAPIKeyNotFound
	placeholder,
placeholder
	cfg := &config.Config{
		APIKeyAuth: config.APIKeyAuthCacheConfig{
			L1Size:             100,
			L1TTLSeconds:       60,
			L2TTLSeconds:       60,
			NegativeTTLSeconds: 30,
	placeholder,
placeholder
	svc := NewAPIKeyService(repo, nil, nil, nil, nil, cache, cfg)
	cache.getAuthCache = func(ctx context.Context, key string) (*APIKeyAuthCacheEntry, error) {
		return nil, redis.Nil
placeholder

	_, err := svc.GetByKey(context.Background(), "missing")
	require.ErrorIs(t, err, ErrAPIKeyNotFound)
	require.Empty(t, cache.setAuthKeys, "attacker-controlled misses must not be written to Redis")
	svc.authNegativeCacheL1.Wait()
	_, err = svc.GetByKey(context.Background(), "missing")
	require.ErrorIs(t, err, ErrAPIKeyNotFound)
	require.Equal(t, int32(1), repoCalls.Load())
placeholder

func TestAPIKeyService_GetByKeyRejectsInvalidLengthBeforeCaches(t *testing.T) {
	var cacheCalls atomic.Int32
	cache := &authCacheStub{getAuthCache: func(context.Context, string) (*APIKeyAuthCacheEntry, error) {
		cacheCalls.Add(1)
		return nil, redis.Nil
placeholderplaceholder
	repo := &authRepoStub{getByKeyForAuth: func(context.Context, string) (*APIKey, error) {
		t.Fatal("invalid credential reached repository")
		return nil, nil
placeholderplaceholder
	svc := NewAPIKeyService(repo, nil, nil, nil, nil, cache, &config.Config{APIKeyAuth: config.APIKeyAuthCacheConfig{L2TTLSeconds: 60placeholderplaceholder)

	for _, key := range []string{"", strings.Repeat("x", MaxAPIKeyCredentialBytes+1)placeholder {
		_, err := svc.GetByKey(context.Background(), key)
		require.ErrorIs(t, err, ErrAPIKeyNotFound)
placeholder
	require.Zero(t, cacheCalls.Load())
placeholder

func TestAPIKeyService_GetByKeyAllowsMaximumLength(t *testing.T) {
	key := strings.Repeat("x", MaxAPIKeyCredentialBytes)
	var repoCalls atomic.Int32
	repo := &authRepoStub{getByKeyForAuth: func(_ context.Context, got string) (*APIKey, error) {
		repoCalls.Add(1)
		require.Equal(t, key, got)
		return nil, ErrAPIKeyNotFound
placeholderplaceholder
	svc := NewAPIKeyService(repo, nil, nil, nil, nil, nil, &config.Config{placeholder)
	_, err := svc.GetByKey(context.Background(), key)
	require.ErrorIs(t, err, ErrAPIKeyNotFound)
	require.Equal(t, int32(1), repoCalls.Load())
placeholder

func TestAPIKeyService_AuthLookupBulkheadRejectsExcessMisses(t *testing.T) {
	entered := make(chan struct{placeholder)
	release := make(chan struct{placeholder)
	repo := &authRepoStub{getByKeyForAuth: func(context.Context, string) (*APIKey, error) {
		close(entered)
		<-release
		return nil, ErrAPIKeyNotFound
placeholderplaceholder
	svc := NewAPIKeyService(repo, nil, nil, nil, nil, nil, &config.Config{APIKeyAuth: config.APIKeyAuthCacheConfig{LookupConcurrency: 1placeholderplaceholder)

	done := make(chan error, 1)
	go func() {
		_, err := svc.GetByKey(context.Background(), "first")
		done <- err
placeholder()
	<-entered

	_, err := svc.GetByKey(context.Background(), "second")
	require.ErrorIs(t, err, ErrAPIKeyAuthOverloaded)
	metrics := svc.AuthLookupMetrics()
	require.Equal(t, uint64(2), metrics.Total)
	require.Equal(t, uint64(1), metrics.Rejected)
	require.Equal(t, int64(1), metrics.InFlight)
	require.Equal(t, 1, metrics.Capacity)

	close(release)
	require.ErrorIs(t, <-done, ErrAPIKeyNotFound)
placeholder

func TestAPIKeyService_GetByKey_SingleflightCollapses(t *testing.T) {
	var calls int32
	cache := &authCacheStub{placeholder
	repo := &authRepoStub{
		getByKeyForAuth: func(ctx context.Context, key string) (*APIKey, error) {
			atomic.AddInt32(&calls, 1)
			time.Sleep(50 * time.Millisecond)
			return &APIKey{
				ID:     11,
				UserID: 2,
				Status: StatusActive,
				User: &User{
					ID:          2,
					Status:      StatusActive,
					Role:        RoleUser,
					Balance:     1,
					Concurrency: 1,
			placeholder,
		placeholder, nil
	placeholder,
placeholder
	cfg := &config.Config{
		APIKeyAuth: config.APIKeyAuthCacheConfig{
			Singleflight: true,
	placeholder,
placeholder
	svc := NewAPIKeyService(repo, nil, nil, nil, nil, cache, cfg)

	start := make(chan struct{placeholder)
	wg := sync.WaitGroup{placeholder
	errs := make([]error, 5)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			<-start
			_, err := svc.GetByKey(context.Background(), "k1")
			errs[idx] = err
	placeholder(i)
placeholder
	close(start)
	wg.Wait()

	for _, err := range errs {
	placeholder
placeholder
	require.Equal(t, int32(1), atomic.LoadInt32(&calls))
placeholder
