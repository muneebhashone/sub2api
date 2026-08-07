//go:build unit

package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/stretchr/testify/require"
)

type grokBaseURLSettingRepoStub struct{ values map[string]string placeholder

func (r *grokBaseURLSettingRepoStub) GetValue(_ context.Context, key string) (string, error) {
	if value, ok := r.values[key]; ok {
		return value, nil
placeholder
	return "", fmt.Errorf("setting %s not found", key)
placeholder
func (r *grokBaseURLSettingRepoStub) Get(context.Context, string) (*Setting, error) {
	return nil, fmt.Errorf("unused")
placeholder
func (r *grokBaseURLSettingRepoStub) Set(context.Context, string, string) error { return nil placeholder
func (r *grokBaseURLSettingRepoStub) GetMultiple(context.Context, []string) (map[string]string, error) {
	return r.values, nil
placeholder
func (r *grokBaseURLSettingRepoStub) SetMultiple(context.Context, map[string]string) error {
	return nil
placeholder
func (r *grokBaseURLSettingRepoStub) GetAll(context.Context) (map[string]string, error) {
	return r.values, nil
placeholder
func (r *grokBaseURLSettingRepoStub) Delete(context.Context, string) error { return nil placeholder

func TestGrokBaseURLForMode(t *testing.T) {
	for _, tc := range []struct {
		mode string
		want string
placeholder{
		{"api", xai.DefaultBaseURLplaceholder,
		{"us-east-1", xai.DefaultUSEast1BaseURLplaceholder,
		{"us-west-2", xai.DefaultUSWest2BaseURLplaceholder,
		{"eu-west-1", xai.DefaultEUWest1BaseURLplaceholder,
		{"cli", xai.DefaultCLIBaseURLplaceholder,
		{"invalid", xai.DefaultCLIBaseURLplaceholder,
placeholder {
		t.Run(tc.mode, func(t *testing.T) {
			require.Equal(t, tc.want, GrokBaseURLForMode(tc.mode))
	placeholder)
placeholder
placeholder

func TestSettingServiceResolveGrokBaseURLHonorsModeAndExplicitPins(t *testing.T) {
	repo := &grokBaseURLSettingRepoStub{values: map[string]string{SettingKeyGrokDefaultBaseURLMode: GrokDefaultBaseURLModeUSWest2placeholderplaceholder
	svc := NewSettingService(repo, nil)
	account := &Account{Platform: PlatformGrok, Type: AccountTypeOAuth, Credentials: map[string]any{placeholderplaceholder
	require.Equal(t, xai.DefaultUSWest2BaseURL, svc.ResolveGrokBaseURL(context.Background(), account))

	// An explicit official endpoint remains pinned.
	account.Credentials["base_url"] = xai.DefaultBaseURL
	require.Equal(t, xai.DefaultBaseURL, svc.ResolveGrokBaseURL(context.Background(), account))

	// An explicit regional pin remains authoritative.
	account.Credentials["base_url"] = xai.DefaultEUWest1BaseURL
	require.Equal(t, xai.DefaultEUWest1BaseURL, svc.ResolveGrokBaseURL(context.Background(), account))
placeholder

func TestAccountGetGrokBaseURLOrPreservesCustomOAuthURLForPolicyValidation(t *testing.T) {
	account := &Account{Platform: PlatformGrok, Type: AccountTypeOAuth, Credentials: map[string]any{
		"base_url": "https://attacker.invalid/v1",
placeholderplaceholder
	require.Equal(t, "https://attacker.invalid/v1", account.GetGrokBaseURLOr(xai.DefaultCLIBaseURL))
placeholder
