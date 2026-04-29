//go:build unit

package service

import (
	"context"
	"testing"
	"time"
)

type snapshotHydrationCache struct {
	snapshot []*Account
	accounts map[int64]*Account
placeholder

func (c *snapshotHydrationCache) GetSnapshot(ctx context.Context, bucket SchedulerBucket) ([]*Account, bool, error) {
	return c.snapshot, true, nil
placeholder

func (c *snapshotHydrationCache) SetSnapshot(ctx context.Context, bucket SchedulerBucket, accounts []Account) error {
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
