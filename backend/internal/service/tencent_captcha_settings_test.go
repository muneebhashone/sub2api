//go:build unit

package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func TestSettingService_ParseSettingsMasksTencentCaptchaCredentials(t *testing.T) {
	svc := NewSettingService(&settingGetAllRepoStub{values: map[string]string{
		SettingKeyTencentCaptchaEnabled:        "true",
		SettingKeyTencentCaptchaAppID:          "123456789",
		SettingKeyTencentCaptchaAppSecretKey:   "app-secret",
		SettingKeyTencentCaptchaCloudSecretID:  "cloud-secret-id",
		SettingKeyTencentCaptchaCloudSecretKey: "cloud-secret-key",
placeholderplaceholder, &config.Config{placeholder)

	settings, err := svc.GetAllSettings(context.Background())

placeholder
	require.True(t, settings.TencentCaptchaEnabled)
	require.Equal(t, "123456789", settings.TencentCaptchaAppID)
	require.True(t, settings.TencentCaptchaAppSecretKeyConfigured)
	require.True(t, settings.TencentCaptchaCloudSecretIDConfigured)
	require.True(t, settings.TencentCaptchaCloudSecretKeyConfigured)
	require.Equal(t, "app-secret", settings.TencentCaptchaAppSecretKey)
	require.Equal(t, "cloud-secret-id", settings.TencentCaptchaCloudSecretID)
	require.Equal(t, "cloud-secret-key", settings.TencentCaptchaCloudSecretKey)
placeholder

func TestSettingService_GetPublicSettingsExposesOnlyTencentCaptchaAppID(t *testing.T) {
	svc := NewSettingService(&settingPublicRepoStub{values: map[string]string{
		SettingKeyTencentCaptchaEnabled:        "true",
		SettingKeyTencentCaptchaAppID:          "123456789",
		SettingKeyTencentCaptchaAppSecretKey:   "app-secret",
		SettingKeyTencentCaptchaCloudSecretID:  "cloud-secret-id",
		SettingKeyTencentCaptchaCloudSecretKey: "cloud-secret-key",
		SettingKeyTencentCaptchaRegion:         TencentCaptchaRegionINTL,
placeholderplaceholder, &config.Config{placeholder)

	settings, err := svc.GetPublicSettings(context.Background())
placeholder
	require.True(t, settings.TencentCaptchaEnabled)
	require.Equal(t, "123456789", settings.TencentCaptchaAppID)
	// 站点必须原样公开下发：前端据此决定加载哪个站点的 SDK 脚本与构造函数形态。
	require.Equal(t, TencentCaptchaRegionINTL, settings.TencentCaptchaRegion)

	raw, err := json.Marshal(settings)
placeholder
	require.NotContains(t, string(raw), "app-secret")
	require.NotContains(t, string(raw), "cloud-secret-id")
	require.NotContains(t, string(raw), "cloud-secret-key")
placeholder

func TestSettingService_GetTencentCaptchaConfig(t *testing.T) {
	repo := &settingPublicRepoStub{values: map[string]string{
		SettingKeyTencentCaptchaEnabled:        "true",
		SettingKeyTencentCaptchaAppID:          "123456789",
		SettingKeyTencentCaptchaAppSecretKey:   "app-secret",
		SettingKeyTencentCaptchaCloudSecretID:  "cloud-secret-id",
		SettingKeyTencentCaptchaCloudSecretKey: "cloud-secret-key",
placeholderplaceholder
	svc := NewSettingService(repo, &config.Config{placeholder)

	got := svc.GetTencentCaptchaConfig(context.Background())

	require.Equal(t, TencentCaptchaConfig{
		Enabled:        true,
		AppID:          "123456789",
		AppSecretKey:   "app-secret",
		CloudSecretID:  "cloud-secret-id",
		CloudSecretKey: "cloud-secret-key",
		// 未配置站点时回落中国站，保持存量部署行为不变
		Region: TencentCaptchaRegionCN,
placeholder, got)
placeholder
