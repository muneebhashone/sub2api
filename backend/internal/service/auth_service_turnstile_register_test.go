//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type turnstileVerifierSpy struct {
	called    int
	lastToken string
	result    *TurnstileVerifyResponse
	err       error
placeholder

func (s *turnstileVerifierSpy) VerifyToken(_ context.Context, _ string, token, _ string) (*TurnstileVerifyResponse, error) {
	s.called++
	s.lastToken = token
	if s.err != nil {
		return nil, s.err
placeholder
	if s.result != nil {
		return s.result, nil
placeholder
	return &TurnstileVerifyResponse{Success: trueplaceholder, nil
placeholder

func newAuthServiceForRegisterTurnstileTest(settings map[string]string, verifier TurnstileVerifier) *AuthService {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Mode: "release",
	placeholder,
		Turnstile: config.TurnstileConfig{
			Required: true,
	placeholder,
placeholder

	settingService := NewSettingService(&settingRepoStub{values: settingsplaceholder, cfg)
	turnstileService := NewTurnstileService(settingService, verifier)

	return NewAuthService(
		&userRepoStub{placeholder,
		nil, // redeemRepo
		nil, // refreshTokenCache
		cfg,
		settingService,
		nil, // emailService
		turnstileService,
		nil, // emailQueueService
		nil, // promoService
	)
placeholder

func TestAuthService_VerifyTurnstileForRegister_SkipWhenEmailVerifyCodeProvided(t *testing.T) {
	verifier := &turnstileVerifierSpy{placeholder
	service := newAuthServiceForRegisterTurnstileTest(map[string]string{
		SettingKeyEmailVerifyEnabled:  "true",
		SettingKeyTurnstileEnabled:    "true",
		SettingKeyTurnstileSecretKey:  "secret",
		SettingKeyRegistrationEnabled: "true",
placeholder, verifier)

	err := service.VerifyTurnstileForRegister(context.Background(), "", "127.0.0.1", "123456")
placeholder
	require.Equal(t, 0, verifier.called)
placeholder

func TestAuthService_VerifyTurnstileForRegister_RequireWhenVerifyCodeMissing(t *testing.T) {
	verifier := &turnstileVerifierSpy{placeholder
	service := newAuthServiceForRegisterTurnstileTest(map[string]string{
		SettingKeyEmailVerifyEnabled: "true",
		SettingKeyTurnstileEnabled:   "true",
		SettingKeyTurnstileSecretKey: "secret",
placeholder, verifier)

	err := service.VerifyTurnstileForRegister(context.Background(), "", "127.0.0.1", "")
	require.ErrorIs(t, err, ErrTurnstileVerificationFailed)
placeholder

func TestAuthService_VerifyTurnstileForRegister_NoSkipWhenEmailVerifyDisabled(t *testing.T) {
	verifier := &turnstileVerifierSpy{placeholder
	service := newAuthServiceForRegisterTurnstileTest(map[string]string{
		SettingKeyEmailVerifyEnabled: "false",
		SettingKeyTurnstileEnabled:   "true",
		SettingKeyTurnstileSecretKey: "secret",
placeholder, verifier)

	err := service.VerifyTurnstileForRegister(context.Background(), "turnstile-token", "127.0.0.1", "123456")
placeholder
	require.Equal(t, 1, verifier.called)
	require.Equal(t, "turnstile-token", verifier.lastToken)
placeholder
