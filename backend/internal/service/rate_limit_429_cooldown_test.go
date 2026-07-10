//go:build unit

package service

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type rateLimit429AccountRepoStub struct {
	mockAccountRepoForGemini
	rateLimitCalls     int
	lastRateLimitID    int64
	lastRateLimitReset time.Time
placeholder

func (r *rateLimit429AccountRepoStub) SetRateLimited(_ context.Context, id int64, resetAt time.Time) error {
	r.rateLimitCalls++
	r.lastRateLimitID = id
	r.lastRateLimitReset = resetAt
	return nil
placeholder

func TestGetRateLimit429CooldownSettings_DefaultsWhenNotSet(t *testing.T) {
	repo := newMockSettingRepo()
	svc := NewSettingService(repo, &config.Config{placeholder)

	settings, err := svc.GetRateLimit429CooldownSettings(context.Background())
placeholder
	require.True(t, settings.Enabled)
	require.Equal(t, 5, settings.CooldownSeconds)
placeholder

func TestGetRateLimit429CooldownSettings_ReadsFromDB(t *testing.T) {
	repo := newMockSettingRepo()
	data, _ := json.Marshal(RateLimit429CooldownSettings{Enabled: false, CooldownSeconds: 12placeholder)
	repo.data[SettingKeyRateLimit429CooldownSettings] = string(data)
	svc := NewSettingService(repo, &config.Config{placeholder)

	settings, err := svc.GetRateLimit429CooldownSettings(context.Background())
placeholder
	require.False(t, settings.Enabled)
	require.Equal(t, 12, settings.CooldownSeconds)
placeholder

func TestSetRateLimit429CooldownSettings_EnabledRejectsOutOfRange(t *testing.T) {
	svc := NewSettingService(newMockSettingRepo(), &config.Config{placeholder)

	for _, seconds := range []int{0, -1, 7201, 99999placeholder {
		err := svc.SetRateLimit429CooldownSettings(context.Background(), &RateLimit429CooldownSettings{
			Enabled: true, CooldownSeconds: seconds,
	placeholder)
		require.Error(t, err, "should reject enabled=true + cooldown_seconds=%d", seconds)
		require.Contains(t, err.Error(), "cooldown_seconds must be between 1-7200")
placeholder
placeholder

func TestHandle429_FallbackUsesDBSeconds(t *testing.T) {
	accountRepo := &rateLimit429AccountRepoStub{placeholder
	settingRepo := newMockSettingRepo()
	data, _ := json.Marshal(RateLimit429CooldownSettings{Enabled: true, CooldownSeconds: 12placeholder)
	settingRepo.data[SettingKeyRateLimit429CooldownSettings] = string(data)

	settingSvc := NewSettingService(settingRepo, &config.Config{placeholder)
	svc := NewRateLimitService(accountRepo, nil, &config.Config{placeholder, nil, nil)
	svc.SetSettingService(settingSvc)

	account := &Account{ID: 42, Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder
	before := time.Now()
	svc.handle429(context.Background(), account, http.Header{placeholder, []byte(`{"error":{"type":"rate_limit_error","message":"slow down"placeholderplaceholder`))
	after := time.Now()

	require.Equal(t, 1, accountRepo.rateLimitCalls)
	require.Equal(t, int64(42), accountRepo.lastRateLimitID)
	require.True(t, !accountRepo.lastRateLimitReset.Before(before.Add(12*time.Second)) && !accountRepo.lastRateLimitReset.After(after.Add(12*time.Second)))
placeholder

func TestHandle429_FallbackDisabledSkipsLocalMark(t *testing.T) {
	accountRepo := &rateLimit429AccountRepoStub{placeholder
	settingRepo := newMockSettingRepo()
	data, _ := json.Marshal(RateLimit429CooldownSettings{Enabled: false, CooldownSeconds: 12placeholder)
	settingRepo.data[SettingKeyRateLimit429CooldownSettings] = string(data)

	settingSvc := NewSettingService(settingRepo, &config.Config{placeholder)
	svc := NewRateLimitService(accountRepo, nil, &config.Config{placeholder, nil, nil)
	svc.SetSettingService(settingSvc)

	account := &Account{ID: 43, Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder
	svc.handle429(context.Background(), account, http.Header{placeholder, []byte(`{"error":{"type":"rate_limit_error","message":"slow down"placeholderplaceholder`))

	require.Zero(t, accountRepo.rateLimitCalls)
placeholder

// Anthropic 无 reset 头的 429（如 Extra usage required）也应走兜底冷却，
// 否则账号永不冷却，调度器会让每个请求反复撞同一批 429 账号（旋转木马）。
func TestHandle429_AnthropicNoResetTimeUsesFallbackCooldown(t *testing.T) {
	accountRepo := &rateLimit429AccountRepoStub{placeholder
	settingRepo := newMockSettingRepo()
	data, _ := json.Marshal(RateLimit429CooldownSettings{Enabled: true, CooldownSeconds: 12placeholder)
	settingRepo.data[SettingKeyRateLimit429CooldownSettings] = string(data)

	settingSvc := NewSettingService(settingRepo, &config.Config{placeholder)
	svc := NewRateLimitService(accountRepo, nil, &config.Config{placeholder, nil, nil)
	svc.SetSettingService(settingSvc)

	account := &Account{ID: 45, Platform: PlatformAnthropic, Type: AccountTypeOAuthplaceholder
	before := time.Now()
	svc.handle429(context.Background(), account, http.Header{placeholder, []byte(`{"error":{"type":"rate_limit_error","message":"Extra usage required"placeholderplaceholder`))
	after := time.Now()

	require.Equal(t, 1, accountRepo.rateLimitCalls)
	require.Equal(t, int64(45), accountRepo.lastRateLimitID)
	require.True(t, !accountRepo.lastRateLimitReset.Before(before.Add(12*time.Second)) && !accountRepo.lastRateLimitReset.After(after.Add(12*time.Second)))
placeholder

// 管理端关闭兜底冷却时，Anthropic 无 reset 头的 429 保持旧行为：不标记账号。
func TestHandle429_AnthropicNoResetTimeFallbackDisabledSkipsMark(t *testing.T) {
	accountRepo := &rateLimit429AccountRepoStub{placeholder
	settingRepo := newMockSettingRepo()
	data, _ := json.Marshal(RateLimit429CooldownSettings{Enabled: false, CooldownSeconds: 12placeholder)
	settingRepo.data[SettingKeyRateLimit429CooldownSettings] = string(data)

	settingSvc := NewSettingService(settingRepo, &config.Config{placeholder)
	svc := NewRateLimitService(accountRepo, nil, &config.Config{placeholder, nil, nil)
	svc.SetSettingService(settingSvc)

	account := &Account{ID: 46, Platform: PlatformAnthropic, Type: AccountTypeOAuthplaceholder
	svc.handle429(context.Background(), account, http.Header{placeholder, []byte(`{"error":{"type":"rate_limit_error","message":"Extra usage required"placeholderplaceholder`))

	require.Zero(t, accountRepo.rateLimitCalls)
placeholder

func TestHandle429_FallbackUsesDefaultSecondsWhenSettingServiceMissing(t *testing.T) {
	accountRepo := &rateLimit429AccountRepoStub{placeholder
	cfg := &config.Config{placeholder
	svc := NewRateLimitService(accountRepo, nil, cfg, nil, nil)

	account := &Account{ID: 44, Platform: PlatformGemini, Type: AccountTypeAPIKeyplaceholder
	before := time.Now()
	svc.handle429(context.Background(), account, http.Header{placeholder, []byte(`{"error":{"message":"slow down"placeholderplaceholder`))
	after := time.Now()

	require.Equal(t, 1, accountRepo.rateLimitCalls)
	require.Equal(t, int64(44), accountRepo.lastRateLimitID)
	require.True(t, !accountRepo.lastRateLimitReset.Before(before.Add(5*time.Second)) && !accountRepo.lastRateLimitReset.After(after.Add(5*time.Second)))
placeholder
