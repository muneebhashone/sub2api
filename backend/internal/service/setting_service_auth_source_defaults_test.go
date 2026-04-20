//go:build unit

package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type authSourceDefaultsRepoStub struct {
	values  map[string]string
	updates map[string]string
placeholder

func (s *authSourceDefaultsRepoStub) Get(ctx context.Context, key string) (*Setting, error) {
	panic("unexpected Get call")
placeholder

func (s *authSourceDefaultsRepoStub) GetValue(ctx context.Context, key string) (string, error) {
	panic("unexpected GetValue call")
placeholder

func (s *authSourceDefaultsRepoStub) Set(ctx context.Context, key, value string) error {
	panic("unexpected Set call")
placeholder

func (s *authSourceDefaultsRepoStub) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	out := make(map[string]string, len(keys))
	for _, key := range keys {
		if value, ok := s.values[key]; ok {
			out[key] = value
	placeholder
placeholder
	return out, nil
placeholder

func (s *authSourceDefaultsRepoStub) SetMultiple(ctx context.Context, settings map[string]string) error {
	s.updates = make(map[string]string, len(settings))
	for key, value := range settings {
		s.updates[key] = value
		if s.values == nil {
			s.values = map[string]string{placeholder
	placeholder
		s.values[key] = value
placeholder
	return nil
placeholder

func (s *authSourceDefaultsRepoStub) GetAll(ctx context.Context) (map[string]string, error) {
	panic("unexpected GetAll call")
placeholder

func (s *authSourceDefaultsRepoStub) Delete(ctx context.Context, key string) error {
	panic("unexpected Delete call")
placeholder

func TestSettingService_GetAuthSourceDefaultSettings_ParsesValuesAndDefaults(t *testing.T) {
	repo := &authSourceDefaultsRepoStub{
		values: map[string]string{
			SettingKeyAuthSourceDefaultEmailBalance:            "12.5",
			SettingKeyAuthSourceDefaultEmailConcurrency:        "7",
			SettingKeyAuthSourceDefaultEmailSubscriptions:      `[{"group_id":11,"validity_days":30placeholder]`,
			SettingKeyAuthSourceDefaultEmailGrantOnSignup:      "false",
			SettingKeyAuthSourceDefaultLinuxDoGrantOnFirstBind: "true",
			SettingKeyForceEmailOnThirdPartySignup:             "true",
	placeholder,
placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)

	got, err := svc.GetAuthSourceDefaultSettings(context.Background())
placeholder
	require.Equal(t, 12.5, got.Email.Balance)
	require.Equal(t, 7, got.Email.Concurrency)
	require.Equal(t, []DefaultSubscriptionSetting{{GroupID: 11, ValidityDays: 30placeholderplaceholder, got.Email.Subscriptions)
	require.False(t, got.Email.GrantOnSignup)
	require.False(t, got.Email.GrantOnFirstBind)
	require.Equal(t, 0.0, got.LinuxDo.Balance)
	require.Equal(t, 5, got.LinuxDo.Concurrency)
	require.Equal(t, []DefaultSubscriptionSetting{placeholder, got.LinuxDo.Subscriptions)
	require.True(t, got.LinuxDo.GrantOnSignup)
	require.True(t, got.LinuxDo.GrantOnFirstBind)
	require.Equal(t, 5, got.OIDC.Concurrency)
	require.Equal(t, 5, got.WeChat.Concurrency)
	require.True(t, got.ForceEmailOnThirdPartySignup)
placeholder

func TestSettingService_UpdateAuthSourceDefaultSettings_PersistsAllKeys(t *testing.T) {
	repo := &authSourceDefaultsRepoStub{placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)

	err := svc.UpdateAuthSourceDefaultSettings(context.Background(), &AuthSourceDefaultSettings{
		Email: ProviderDefaultGrantSettings{
			Balance:          1.25,
			Concurrency:      3,
			Subscriptions:    []DefaultSubscriptionSetting{{GroupID: 21, ValidityDays: 14placeholderplaceholder,
			GrantOnSignup:    false,
			GrantOnFirstBind: true,
	placeholder,
		LinuxDo: ProviderDefaultGrantSettings{
			Balance:          2,
			Concurrency:      4,
			Subscriptions:    []DefaultSubscriptionSetting{{GroupID: 22, ValidityDays: 30placeholderplaceholder,
			GrantOnSignup:    true,
			GrantOnFirstBind: false,
	placeholder,
		OIDC: ProviderDefaultGrantSettings{
			Balance:          3,
			Concurrency:      5,
			Subscriptions:    []DefaultSubscriptionSetting{{GroupID: 23, ValidityDays: 60placeholderplaceholder,
			GrantOnSignup:    true,
			GrantOnFirstBind: true,
	placeholder,
		WeChat: ProviderDefaultGrantSettings{
			Balance:          4,
			Concurrency:      6,
			Subscriptions:    []DefaultSubscriptionSetting{{GroupID: 24, ValidityDays: 90placeholderplaceholder,
			GrantOnSignup:    false,
			GrantOnFirstBind: false,
	placeholder,
		ForceEmailOnThirdPartySignup: true,
placeholder)
placeholder
	require.Equal(t, "1.25000000", repo.updates[SettingKeyAuthSourceDefaultEmailBalance])
	require.Equal(t, "3", repo.updates[SettingKeyAuthSourceDefaultEmailConcurrency])
	require.Equal(t, "false", repo.updates[SettingKeyAuthSourceDefaultEmailGrantOnSignup])
	require.Equal(t, "true", repo.updates[SettingKeyAuthSourceDefaultEmailGrantOnFirstBind])
	require.Equal(t, "true", repo.updates[SettingKeyForceEmailOnThirdPartySignup])

	var got []DefaultSubscriptionSetting
	require.NoError(t, json.Unmarshal([]byte(repo.updates[SettingKeyAuthSourceDefaultWeChatSubscriptions]), &got))
	require.Equal(t, []DefaultSubscriptionSetting{{GroupID: 24, ValidityDays: 90placeholderplaceholder, got)
placeholder
