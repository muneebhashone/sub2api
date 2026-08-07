//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type thresholdSelectionAccountRepoStub struct {
	rateLimitAccountRepoStub
	accounts []Account
placeholder

func (r *thresholdSelectionAccountRepoStub) ListSchedulableByPlatform(_ context.Context, platform string) ([]Account, error) {
	filtered := make([]Account, 0, len(r.accounts))
	for _, account := range r.accounts {
		if account.Platform == platform {
			filtered = append(filtered, account)
	placeholder
placeholder
	return filtered, nil
placeholder

func (r *thresholdSelectionAccountRepoStub) ListSchedulableByGroupIDAndPlatform(ctx context.Context, _ int64, platform string) ([]Account, error) {
	return r.ListSchedulableByPlatform(ctx, platform)
placeholder

func (r *thresholdSelectionAccountRepoStub) ListSchedulableUngroupedByPlatform(ctx context.Context, platform string) ([]Account, error) {
	return r.ListSchedulableByPlatform(ctx, platform)
placeholder

func TestGatewayService_ListSchedulableAccounts_DoesNotFilterUnsupportedThresholdPlatforms(t *testing.T) {
	accountSchedulingThresholdsSF.Forget(SettingKeyAccountSchedulingThresholds)
	accountSchedulingThresholdsCache.Store(&cachedAccountSchedulingThresholds{placeholder)

	settingsRepo := newMockSettingRepo()
	settingsRepo.data[SettingKeyAccountSchedulingThresholds] = `{"openai":90placeholder`

	accountRepo := &thresholdSelectionAccountRepoStub{
		accounts: []Account{
			{
				ID:          3101,
				Platform:    PlatformKiro,
				Status:      StatusActive,
				Schedulable: true,
		placeholder
					"account_scheduling_threshold": 1,
			placeholder,
				Extra: map[string]any{
					"kiro_sched_utilization": 95.0,
					"kiro_sched_reset_at":    time.Now().UTC().Add(2 * time.Hour).Format(time.RFC3339),
			placeholder,
		placeholder,
			{
				ID:          3102,
				Platform:    PlatformKiro,
				Status:      StatusActive,
				Schedulable: true,
				Extra: map[string]any{
					"kiro_sched_utilization": 42.0,
					"kiro_sched_reset_at":    time.Now().UTC().Add(2 * time.Hour).Format(time.RFC3339),
			placeholder,
		placeholder,
	placeholder,
placeholder

	rateLimitService := NewRateLimitService(accountRepo, nil, &config.Config{placeholder, nil, nil)
	rateLimitService.SetSettingService(NewSettingService(settingsRepo, &config.Config{placeholder))
	svc := &GatewayService{
		accountRepo:      accountRepo,
		cfg:              &config.Config{placeholder,
		rateLimitService: rateLimitService,
placeholder

	accounts, useMixed, err := svc.listSchedulableAccounts(context.Background(), nil, PlatformKiro, false)

placeholder
	require.False(t, useMixed)
	require.Len(t, accounts, 2)
	require.Equal(t, int64(3101), accounts[0].ID)
	require.Equal(t, int64(3102), accounts[1].ID)
	require.Equal(t, 0, accountRepo.tempCalls)
placeholder

func TestOpenAIGatewayService_ListSchedulableAccounts_FiltersThresholdBlockedAccounts(t *testing.T) {
	accountSchedulingThresholdsSF.Forget(SettingKeyAccountSchedulingThresholds)
	accountSchedulingThresholdsCache.Store(&cachedAccountSchedulingThresholds{placeholder)

	settingsRepo := newMockSettingRepo()
	settingsRepo.data[SettingKeyAccountSchedulingThresholds] = `{"openai":85placeholder`

	accountRepo := &thresholdSelectionAccountRepoStub{
		accounts: []Account{
			{
				ID:          4101,
				Platform:    PlatformOpenAI,
				Status:      StatusActive,
				Schedulable: true,
				Extra: map[string]any{
					"codex_7d_used_percent": 91.0,
					"codex_7d_reset_at":     time.Now().UTC().Add(12 * time.Hour).Format(time.RFC3339),
			placeholder,
		placeholder,
			{
				ID:          4102,
				Platform:    PlatformOpenAI,
				Status:      StatusActive,
				Schedulable: true,
				Extra: map[string]any{
					"codex_7d_used_percent": 40.0,
					"codex_7d_reset_at":     time.Now().UTC().Add(12 * time.Hour).Format(time.RFC3339),
			placeholder,
		placeholder,
	placeholder,
placeholder

	rateLimitService := NewRateLimitService(accountRepo, nil, &config.Config{placeholder, nil, nil)
	rateLimitService.SetSettingService(NewSettingService(settingsRepo, &config.Config{placeholder))
	svc := &OpenAIGatewayService{
		accountRepo:      accountRepo,
		cfg:              &config.Config{placeholder,
		rateLimitService: rateLimitService,
placeholder

	accounts, err := svc.listSchedulableAccounts(context.Background(), nil, PlatformOpenAI)

placeholder
	require.Len(t, accounts, 1)
	require.Equal(t, int64(4102), accounts[0].ID)
	require.Equal(t, 1, accountRepo.tempCalls)
placeholder
