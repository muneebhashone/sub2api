//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
)

type snapshotHydrationCache struct {
	snapshot []*Account
	accounts map[int64]*Account
placeholder

func (c *snapshotHydrationCache) GetSnapshot(ctx context.Context, bucket SchedulerBucket) ([]*Account, bool, error) {
	return c.snapshot, true, nil
placeholder

func (c *snapshotHydrationCache) CaptureBucketWriteToken(ctx context.Context, bucket SchedulerBucket) (SchedulerBucketWriteToken, error) {
	return SchedulerBucketWriteToken{Bucket: bucket, Epoch: 1placeholder, nil
placeholder

func (c *snapshotHydrationCache) SetSnapshot(ctx context.Context, bucket SchedulerBucket, token SchedulerBucketWriteToken, accounts []Account) error {
	return nil
placeholder

func (c *snapshotHydrationCache) RetireBucket(ctx context.Context, bucket SchedulerBucket) error {
	return nil
placeholder

func (c *snapshotHydrationCache) ReopenBucket(ctx context.Context, bucket SchedulerBucket) (SchedulerBucketWriteToken, error) {
	return SchedulerBucketWriteToken{Bucket: bucket, Epoch: 1placeholder, nil
placeholder

func (c *snapshotHydrationCache) TryAcquireGroupLifecycleLease(context.Context, int64, time.Duration) (SchedulerGroupLifecycleLease, bool, error) {
	return SchedulerGroupLifecycleLease{placeholder, false, nil
placeholder

func (c *snapshotHydrationCache) ReleaseGroupLifecycleLease(context.Context, SchedulerGroupLifecycleLease) error {
	return nil
placeholder

func (c *snapshotHydrationCache) GetAccount(ctx context.Context, accountID int64) (*Account, error) {
	if c.accounts == nil {
		return nil, nil
placeholder
	return c.accounts[accountID], nil
placeholder

func (c *snapshotHydrationCache) SetAccount(ctx context.Context, account *Account) error {
	return nil
placeholder

func (c *snapshotHydrationCache) DeleteAccount(ctx context.Context, accountID int64) error {
	return nil
placeholder

func (c *snapshotHydrationCache) UpdateLastUsed(ctx context.Context, updates map[int64]time.Time) error {
	return nil
placeholder

func (c *snapshotHydrationCache) TryLockBucket(ctx context.Context, bucket SchedulerBucket, ttl time.Duration) (bool, error) {
	return true, nil
placeholder

func (c *snapshotHydrationCache) UnlockBucket(ctx context.Context, bucket SchedulerBucket) error {
	return nil
placeholder

func (c *snapshotHydrationCache) ListBuckets(ctx context.Context) ([]SchedulerBucket, error) {
	return nil, nil
placeholder

func (c *snapshotHydrationCache) GetOutboxWatermark(ctx context.Context) (int64, error) {
	return 0, nil
placeholder

func (c *snapshotHydrationCache) SetOutboxWatermark(ctx context.Context, id int64) error {
	return nil
placeholder

func TestOpenAISelectAccountWithLoadAwareness_HydratesSelectedAccountFromSchedulerSnapshot(t *testing.T) {
	cache := &snapshotHydrationCache{
		snapshot: []*Account{
			{
				ID:          1,
				Platform:    PlatformOpenAI,
				Type:        AccountTypeAPIKey,
				Status:      StatusActive,
				Schedulable: true,
				Concurrency: 1,
				Priority:    1,
		placeholder
					"model_mapping": map[string]any{
						"gpt-4": "gpt-4",
				placeholder,
			placeholder,
		placeholder,
	placeholder,
		accounts: map[int64]*Account{
			1: {
				ID:          1,
				Platform:    PlatformOpenAI,
				Type:        AccountTypeAPIKey,
				Status:      StatusActive,
				Schedulable: true,
				Concurrency: 1,
				Priority:    1,
		placeholder
					"api_key":       "sk-live",
					"model_mapping": map[string]any{"gpt-4": "gpt-4"placeholder,
			placeholder,
		placeholder,
	placeholder,
placeholder

	schedulerSnapshot := NewSchedulerSnapshotService(cache, nil, nil, nil, nil)
	groupID := int64(2)
	svc := &OpenAIGatewayService{
		schedulerSnapshot: schedulerSnapshot,
		cache:             &stubGatewayCache{placeholder,
placeholder

	selection, err := svc.SelectAccountWithLoadAwareness(context.Background(), &groupID, "", "gpt-4", nil)
	if err != nil {
		t.Fatalf("SelectAccountWithLoadAwareness error: %v", err)
placeholder
	if selection == nil || selection.Account == nil {
		t.Fatalf("expected selected account")
placeholder
	if got := selection.Account.GetOpenAIApiKey(); got != "sk-live" {
		t.Fatalf("expected hydrated api key, got %q", got)
placeholder
placeholder

func TestOpenAINewAcquiredSelectionResult_ReleasesSlotWhenHydrationFails(t *testing.T) {
	cache := &snapshotHydrationCache{
		accounts: map[int64]*Account{placeholder,
placeholder
	schedulerSnapshot := NewSchedulerSnapshotService(cache, nil, stubOpenAIAccountRepo{placeholder, nil, nil)
	svc := &OpenAIGatewayService{
		schedulerSnapshot: schedulerSnapshot,
placeholder
	releaseCalls := 0

	selection, err := svc.newAcquiredSelectionResult(context.Background(), &Account{ID: 1001placeholder, func() {
		releaseCalls++
placeholder)

	if err == nil {
		t.Fatalf("expected hydration error")
placeholder
	if selection != nil {
		t.Fatalf("expected nil selection on hydration error")
placeholder
	if releaseCalls != 1 {
		t.Fatalf("expected release to be called once, got %d", releaseCalls)
placeholder
placeholder

func TestGatewaySelectAccountWithLoadAwareness_HydratesSelectedAccountFromSchedulerSnapshot(t *testing.T) {
	cache := &snapshotHydrationCache{
		snapshot: []*Account{
			{
				ID:          9,
				Platform:    PlatformAnthropic,
				Type:        AccountTypeAPIKey,
				Status:      StatusActive,
				Schedulable: true,
				Concurrency: 1,
				Priority:    1,
		placeholder,
	placeholder,
		accounts: map[int64]*Account{
			9: {
				ID:          9,
				Platform:    PlatformAnthropic,
				Type:        AccountTypeAPIKey,
				Status:      StatusActive,
				Schedulable: true,
				Concurrency: 1,
				Priority:    1,
		placeholder
					"api_key": "anthropic-live-key",
			placeholder,
		placeholder,
	placeholder,
placeholder

	schedulerSnapshot := NewSchedulerSnapshotService(cache, nil, nil, nil, nil)
	svc := &GatewayService{
		schedulerSnapshot: schedulerSnapshot,
		cache:             &mockGatewayCacheForPlatform{placeholder,
		cfg:               testConfig(),
placeholder

	result, err := svc.SelectAccountWithLoadAwareness(context.Background(), nil, "", "claude-3-5-sonnet-20241022", nil, "", 0)
	if err != nil {
		t.Fatalf("SelectAccountWithLoadAwareness error: %v", err)
placeholder
	if result == nil || result.Account == nil {
		t.Fatalf("expected selected account")
placeholder
	if got := result.Account.GetCredential("api_key"); got != "anthropic-live-key" {
		t.Fatalf("expected hydrated api key, got %q", got)
placeholder
placeholder

func TestGatewaySelectAccountWithLoadAwareness_SkipsAntigravityGeminiFamilyRateLimitedSnapshot(t *testing.T) {
	resetAt := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
	cache := &snapshotHydrationCache{
		snapshot: []*Account{
			{
				ID:          1,
				Platform:    PlatformAntigravity,
				Type:        AccountTypeOAuth,
				Status:      StatusActive,
				Schedulable: true,
				Concurrency: 1,
				Priority:    1,
				AccountGroups: []AccountGroup{
					{AccountID: 1, GroupID: 22placeholder,
			placeholder,
				GroupIDs: []int64{22placeholder,
				Extra: map[string]any{
					"mixed_scheduling": true,
					modelRateLimitsKey: map[string]any{
						antigravityGeminiModelRateLimitKey: map[string]any{
							"rate_limit_reset_at": resetAt,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			{
				ID:          2,
				Platform:    PlatformAntigravity,
				Type:        AccountTypeOAuth,
				Status:      StatusActive,
				Schedulable: true,
				Concurrency: 1,
				Priority:    2,
				AccountGroups: []AccountGroup{
					{AccountID: 2, GroupID: 22placeholder,
			placeholder,
				GroupIDs: []int64{22placeholder,
				Extra: map[string]any{
					"mixed_scheduling": true,
			placeholder,
		placeholder,
	placeholder,
		accounts: map[int64]*Account{
			1: {ID: 1, Platform: PlatformAntigravity, Type: AccountTypeOAuthplaceholder,
			2: {ID: 2, Platform: PlatformAntigravity, Type: AccountTypeOAuthplaceholder,
	placeholder,
placeholder
	groupID := int64(22)
	svc := &GatewayService{
		schedulerSnapshot: NewSchedulerSnapshotService(cache, nil, nil, nil, nil),
		groupRepo: &mockGroupRepoForGateway{
			groups: map[int64]*Group{
				groupID: {
					ID:       groupID,
			placeholder
					Status:   StatusActive,
					Hydrated: true,
			placeholder,
		placeholder,
	placeholder,
		concurrencyService: NewConcurrencyService(&mockConcurrencyCache{placeholder),
		cfg: &config.Config{
			Gateway: config.GatewayConfig{
				Scheduling: config.GatewaySchedulingConfig{
					LoadBatchEnabled:         true,
					StickySessionMaxWaiting:  3,
					StickySessionWaitTimeout: time.Second,
					FallbackWaitTimeout:      time.Second,
					FallbackMaxWaiting:       10,
			placeholder,
		placeholder,
	placeholder,
placeholder

	result, err := svc.SelectAccountWithLoadAwareness(context.Background(), &groupID, "", "gemini-3-flash-preview", nil, "", 0)
	if err != nil {
		t.Fatalf("SelectAccountWithLoadAwareness error: %v", err)
placeholder
	if result == nil || result.Account == nil {
		t.Fatalf("expected selected account")
placeholder
	if result.Account.ID != 2 {
		t.Fatalf("expected scheduler to skip Gemini-family limited antigravity account 1, got %d", result.Account.ID)
placeholder
placeholder
