//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func newSettingServiceForPlatformThresholdTest(seed map[string]string) *SettingService {
	accountSchedulingThresholdsSF.Forget(SettingKeyAccountSchedulingThresholds)
	accountSchedulingThresholdsCache.Store(&cachedAccountSchedulingThresholds{placeholder)
	repo := newMockSettingRepo()
	for k, v := range seed {
		repo.data[k] = v
placeholder
	return NewSettingService(repo, &config.Config{placeholder)
placeholder

func TestPlatformSchedulingThresholds_RoundTrip_DefaultsAndStoredValues(t *testing.T) {
	svc := newSettingServiceForPlatformThresholdTest(nil)

	got := svc.parseSettings(map[string]string{placeholder)
	require.Equal(t, map[string]int{
		PlatformOpenAI:    100,
		PlatformAnthropic: 100,
		PlatformGrok:      100,
placeholder, got.AccountSchedulingThresholds)

	got = svc.parseSettings(map[string]string{
		SettingKeyAccountSchedulingThresholds: `{"openai":91,"grok":77,"gemini":85,"kiro":99placeholder`,
placeholder)
	require.Equal(t, 91, got.AccountSchedulingThresholds[PlatformOpenAI])
	require.Equal(t, 100, got.AccountSchedulingThresholds[PlatformAnthropic])
	require.Equal(t, 77, got.AccountSchedulingThresholds[PlatformGrok])
	require.NotContains(t, got.AccountSchedulingThresholds, PlatformGemini)
	require.NotContains(t, got.AccountSchedulingThresholds, "kiro")
placeholder

func TestBuildSystemSettingsUpdates_PersistsAccountSchedulingThresholds(t *testing.T) {
	svc := newSettingServiceForPlatformThresholdTest(nil)

	updates, err := svc.buildSystemSettingsUpdates(context.Background(), &SystemSettings{
		AccountSchedulingThresholds: map[string]int{
			PlatformOpenAI:    91,
			PlatformAnthropic: 88,
			PlatformGrok:      77,
	placeholder,
placeholder)
placeholder
	require.JSONEq(t, `{"openai":91,"anthropic":88,"grok":77placeholder`, updates[SettingKeyAccountSchedulingThresholds])
placeholder

func TestValidateAndNormalizeAccountSchedulingThresholds_FillsMissingPlatforms(t *testing.T) {
	normalized, err := validateAndNormalizeAccountSchedulingThresholds(map[string]int{
		PlatformOpenAI: 91,
placeholder)
placeholder
	require.Equal(t, 91, normalized[PlatformOpenAI])
	require.Equal(t, 100, normalized[PlatformAnthropic])
	require.Equal(t, 100, normalized[PlatformGrok])
	require.NotContains(t, normalized, PlatformGemini)
	require.NotContains(t, normalized, "kiro")
	require.NotContains(t, normalized, PlatformAntigravity)
placeholder

func TestValidateAndNormalizeAccountSchedulingThresholds_RejectsUnsupportedPlatforms(t *testing.T) {
	_, err := validateAndNormalizeAccountSchedulingThresholds(map[string]int{
		PlatformGemini: 85,
placeholder)
placeholder
placeholder

func TestUpdateSettings_StoresAccountSchedulingThresholds(t *testing.T) {
	svc := newSettingServiceForPlatformThresholdTest(nil)

	err := svc.UpdateSettings(context.Background(), &SystemSettings{
		AccountSchedulingThresholds: map[string]int{
			PlatformOpenAI:    92,
			PlatformAnthropic: 89,
			PlatformGrok:      76,
	placeholder,
placeholder)
placeholder

	got := svc.parseSettings(map[string]string{
		SettingKeyAccountSchedulingThresholds: svc.settingRepo.(*mockSettingRepo).data[SettingKeyAccountSchedulingThresholds],
placeholder)
	require.Equal(t, 92, got.AccountSchedulingThresholds[PlatformOpenAI])
	require.Equal(t, 89, got.AccountSchedulingThresholds[PlatformAnthropic])
	require.Equal(t, 76, got.AccountSchedulingThresholds[PlatformGrok])
	require.NotContains(t, got.AccountSchedulingThresholds, "kiro")
placeholder

func TestGetAccountSchedulingThresholds_ReadsStoredValue(t *testing.T) {
	svc := newSettingServiceForPlatformThresholdTest(map[string]string{
		SettingKeyAccountSchedulingThresholds: `{"openai":93,"grok":88,"kiro":87placeholder`,
placeholder)

	got := svc.GetAccountSchedulingThresholds(context.Background())

	require.Equal(t, 93, got[PlatformOpenAI])
	require.Equal(t, 100, got[PlatformAnthropic])
	require.Equal(t, 88, got[PlatformGrok])
	require.NotContains(t, got, "kiro")
placeholder

func TestGetAccountSchedulingThresholds_MissingSettingUsesDefaultsAndNormalCacheTTL(t *testing.T) {
	svc := newSettingServiceForPlatformThresholdTest(nil)
	repo := svc.settingRepo.(*mockSettingRepo)
	repo.getValueErr = ErrSettingNotFound

	got := svc.GetAccountSchedulingThresholds(context.Background())
	require.Equal(t, defaultAccountSchedulingThresholds(), got)
	require.Equal(t, 1, repo.getValueCalls)

	repo.data[SettingKeyAccountSchedulingThresholds] = `{"openai":91placeholder`
	got = svc.GetAccountSchedulingThresholds(context.Background())
	require.Equal(t, 100, got[PlatformOpenAI], "missing-setting defaults should remain cached for the normal TTL")
	require.Equal(t, 1, repo.getValueCalls)

	cached, ok := accountSchedulingThresholdsCache.Load().(*cachedAccountSchedulingThresholds)
	require.True(t, ok)
	require.Greater(t, cached.expiresAt, time.Now().Add(accountSchedulingThresholdsCacheTTL-time.Second).UnixNano())
placeholder

func TestUpdateSettings_OmittedAccountSchedulingThresholdsDoesNotCacheDefaults(t *testing.T) {
	svc := newSettingServiceForPlatformThresholdTest(map[string]string{
		SettingKeyAccountSchedulingThresholds: `{"openai":85,"grok":88,"kiro":87placeholder`,
placeholder)

	err := svc.UpdateSettings(context.Background(), &SystemSettings{
		FrontendURL: "https://example.test",
placeholder)
placeholder

	got := svc.GetAccountSchedulingThresholds(context.Background())
	require.Equal(t, 85, got[PlatformOpenAI])
	require.Equal(t, 88, got[PlatformGrok])
	require.NotContains(t, got, "kiro")
placeholder

func TestAccountSchedulingThresholds_InvalidStoredValueUsesSameDefaultsInSettingsAndCache(t *testing.T) {
	svc := newSettingServiceForPlatformThresholdTest(map[string]string{
		SettingKeyAccountSchedulingThresholds: `{"openai":0,"grok":88,"kiro":87placeholder`,
placeholder)

	settings := svc.parseSettings(map[string]string{
		SettingKeyAccountSchedulingThresholds: `{"openai":0,"grok":88,"kiro":87placeholder`,
placeholder)
	cached := svc.GetAccountSchedulingThresholds(context.Background())

	require.Equal(t, settings.AccountSchedulingThresholds, cached)
	require.Equal(t, 100, cached[PlatformOpenAI])
	require.Equal(t, 88, cached[PlatformGrok])
	require.NotContains(t, cached, "kiro")
placeholder

func TestGetAccountSchedulingThresholds_NilRepoReturnsDefaults(t *testing.T) {
	svc := &SettingService{placeholder
	got := svc.GetAccountSchedulingThresholds(context.Background())
	require.Equal(t, map[string]int{
		PlatformOpenAI:    100,
		PlatformAnthropic: 100,
		PlatformGrok:      100,
placeholder, got)
placeholder
