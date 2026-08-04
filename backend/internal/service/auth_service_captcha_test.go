//go:build unit

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func newAuthServiceForCaptchaRepoTest(repo *settingRepoStub, required bool, turnstileVerifier TurnstileVerifier, tencentVerifier TencentCaptchaVerifier) *AuthService {
	cfg := &config.Config{
		Server:    config.ServerConfig{Mode: "release"placeholder,
		Turnstile: config.TurnstileConfig{Required: requiredplaceholder,
placeholder
	settingService := NewSettingService(repo, cfg)
	turnstileService := NewTurnstileService(settingService, turnstileVerifier)
	tencentService := NewTencentCaptchaService(settingService, tencentVerifier)
	svc := NewAuthService(nil, &userRepoStub{placeholder, nil, nil, cfg, settingService, nil, turnstileService, nil, nil, nil, nil, nil)
	svc.SetTencentCaptchaService(tencentService)
	return svc
placeholder

func newAuthServiceForCaptchaTest(settings map[string]string, required bool, turnstileVerifier TurnstileVerifier, tencentVerifier TencentCaptchaVerifier) *AuthService {
	cfg := &config.Config{
		Server:    config.ServerConfig{Mode: "release"placeholder,
		Turnstile: config.TurnstileConfig{Required: requiredplaceholder,
placeholder
	settingService := NewSettingService(&settingRepoStub{values: settingsplaceholder, cfg)
	var turnstileService *TurnstileService
	if turnstileVerifier != nil {
		turnstileService = NewTurnstileService(settingService, turnstileVerifier)
placeholder
	svc := NewAuthService(nil, &userRepoStub{placeholder, nil, nil, cfg, settingService, nil, turnstileService, nil, nil, nil, nil, nil)
	if tencentVerifier != nil {
		svc.SetTencentCaptchaService(NewTencentCaptchaService(settingService, tencentVerifier))
placeholder
	return svc
placeholder

func tencentCaptchaSettings() map[string]string {
	return map[string]string{
		SettingKeyTencentCaptchaEnabled:        "true",
		SettingKeyTencentCaptchaAppID:          "123456789",
		SettingKeyTencentCaptchaAppSecretKey:   "app-secret",
		SettingKeyTencentCaptchaCloudSecretID:  "cloud-secret-id",
		SettingKeyTencentCaptchaCloudSecretKey: "cloud-secret-key",
placeholder
placeholder

func TestVerifyCaptchaUsesTencentWhenEnabled(t *testing.T) {
	verifier := &tencentCaptchaVerifierStub{response: &TencentCaptchaVerifyResponse{CaptchaCode: 1placeholderplaceholder
	svc := newAuthServiceForCaptchaTest(tencentCaptchaSettings(), false, nil, verifier)

	err := svc.VerifyCaptcha(context.Background(), CaptchaProof{
		TencentTicket:  "ticket",
		TencentRandstr: "@rand",
placeholder, "203.0.113.10")

placeholder
	require.Equal(t, 1, verifier.calls)
placeholder

func TestVerifyCaptchaRejectsDirtyDoubleEnabledSettings(t *testing.T) {
	settings := tencentCaptchaSettings()
	settings[SettingKeyTurnstileEnabled] = "true"
	settings[SettingKeyTurnstileSecretKey] = "turnstile-secret"
	turnstileVerifier := &turnstileVerifierSpy{placeholder
	tencentVerifier := &tencentCaptchaVerifierStub{response: &TencentCaptchaVerifyResponse{CaptchaCode: 1placeholderplaceholder
	svc := newAuthServiceForCaptchaTest(settings, false, turnstileVerifier, tencentVerifier)

	err := svc.VerifyCaptcha(context.Background(), CaptchaProof{
		TurnstileToken: "turnstile-token",
		TencentTicket:  "ticket",
		TencentRandstr: "@rand",
placeholder, "203.0.113.10")

	require.ErrorIs(t, err, ErrCaptchaProviderConflict)
	require.Zero(t, turnstileVerifier.called)
	require.Zero(t, tencentVerifier.calls)
placeholder

func TestVerifyCaptchaRequiredModeAcceptsCompleteTencentProvider(t *testing.T) {
	verifier := &tencentCaptchaVerifierStub{response: &TencentCaptchaVerifyResponse{CaptchaCode: 1placeholderplaceholder
	svc := newAuthServiceForCaptchaTest(tencentCaptchaSettings(), true, nil, verifier)

	err := svc.VerifyCaptcha(context.Background(), CaptchaProof{
		TencentTicket:  "ticket",
		TencentRandstr: "@rand",
placeholder, "203.0.113.10")

placeholder
placeholder

func TestVerifyCaptchaForRegisterSkipsDuplicateTencentTicketAfterEmailCode(t *testing.T) {
	settings := tencentCaptchaSettings()
	settings[SettingKeyEmailVerifyEnabled] = "true"
	verifier := &tencentCaptchaVerifierStub{response: &TencentCaptchaVerifyResponse{CaptchaCode: 1placeholderplaceholder
	svc := newAuthServiceForCaptchaTest(settings, true, nil, verifier)

	err := svc.VerifyCaptchaForRegister(context.Background(), CaptchaProof{placeholder, "203.0.113.10", "123456")

placeholder
	require.Zero(t, verifier.calls)
placeholder

func TestVerifyCaptchaFailsClosedWhenProviderSettingsCannotBeRead(t *testing.T) {
	repo := &settingRepoStub{err: errors.New("settings unavailable")placeholder
	svc := newAuthServiceForCaptchaRepoTest(repo, false, &turnstileVerifierSpy{placeholder, &tencentCaptchaVerifierStub{placeholder)

	err := svc.VerifyCaptcha(context.Background(), CaptchaProof{placeholder, "203.0.113.10")

	require.ErrorIs(t, err, ErrServiceUnavailable)
placeholder

func TestVerifyCaptchaReadsProviderConfigurationOnce(t *testing.T) {
	repo := &settingRepoStub{values: tencentCaptchaSettings()placeholder
	verifier := &tencentCaptchaVerifierStub{response: &TencentCaptchaVerifyResponse{CaptchaCode: 1placeholderplaceholder
	svc := newAuthServiceForCaptchaRepoTest(repo, false, &turnstileVerifierSpy{placeholder, verifier)

	err := svc.VerifyCaptcha(context.Background(), CaptchaProof{
		TencentTicket:  "ticket",
		TencentRandstr: "@rand",
placeholder, "203.0.113.10")

placeholder
	require.Equal(t, 1, repo.getMultipleCalls)
	require.Zero(t, repo.getValueCalls)
	require.Equal(t, 1, verifier.calls)
placeholder

func TestVerifyCaptchaRejectsEnabledTencentProviderWithIncompleteCredentials(t *testing.T) {
	repo := &settingRepoStub{values: map[string]string{
		SettingKeyTencentCaptchaEnabled: "true",
		SettingKeyTencentCaptchaAppID:   "123456789",
placeholderplaceholder
	verifier := &tencentCaptchaVerifierStub{response: &TencentCaptchaVerifyResponse{CaptchaCode: 1placeholderplaceholder
	svc := newAuthServiceForCaptchaRepoTest(repo, false, &turnstileVerifierSpy{placeholder, verifier)

	err := svc.VerifyCaptcha(context.Background(), CaptchaProof{
		TencentTicket:  "ticket",
		TencentRandstr: "@rand",
placeholder, "203.0.113.10")

	require.ErrorIs(t, err, ErrTencentCaptchaNotConfigured)
	require.Equal(t, 1, repo.getMultipleCalls)
	require.Zero(t, verifier.calls)
placeholder

func TestVerifyTencentCaptchaIfEnabledVerifiesTencentProof(t *testing.T) {
	verifier := &tencentCaptchaVerifierStub{response: &TencentCaptchaVerifyResponse{CaptchaCode: 1placeholderplaceholder
	svc := newAuthServiceForCaptchaTest(tencentCaptchaSettings(), false, nil, verifier)

	err := svc.VerifyTencentCaptchaIfEnabled(context.Background(), CaptchaProof{
		TencentTicket:  "ticket",
		TencentRandstr: "@rand",
placeholder, "203.0.113.10")

placeholder
	require.Equal(t, 1, verifier.calls)
	require.Equal(t, TencentCaptchaProof{Ticket: "ticket", Randstr: "@rand"placeholder, verifier.proof)
placeholder

func TestVerifyTencentCaptchaIfEnabledDoesNotExpandTurnstileCoverage(t *testing.T) {
	settings := map[string]string{
		SettingKeyTurnstileEnabled:   "true",
		SettingKeyTurnstileSecretKey: "turnstile-secret",
placeholder
	turnstileVerifier := &turnstileVerifierSpy{placeholder
	svc := newAuthServiceForCaptchaTest(settings, false, turnstileVerifier, nil)

	err := svc.VerifyTencentCaptchaIfEnabled(context.Background(), CaptchaProof{placeholder, "203.0.113.10")

placeholder
	require.Zero(t, turnstileVerifier.called)
placeholder

func TestVerifyTencentCaptchaIfEnabledFailsClosedOnSettingReadError(t *testing.T) {
	repo := &settingRepoStub{err: errors.New("settings unavailable")placeholder
	svc := newAuthServiceForCaptchaRepoTest(repo, false, &turnstileVerifierSpy{placeholder, &tencentCaptchaVerifierStub{placeholder)

	err := svc.VerifyTencentCaptchaIfEnabled(context.Background(), CaptchaProof{placeholder, "203.0.113.10")

	require.ErrorIs(t, err, ErrServiceUnavailable)
placeholder
