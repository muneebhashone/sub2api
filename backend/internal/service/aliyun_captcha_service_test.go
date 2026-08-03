//go:build unit

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type aliyunVerifierSpy struct {
	called    int
	lastCred  AliyunCaptchaCredentials
	lastParam string
	result    *AliyunCaptchaVerifyResult
	err       error
placeholder

func (s *aliyunVerifierSpy) VerifyCaptcha(_ context.Context, cred AliyunCaptchaCredentials, param string) (*AliyunCaptchaVerifyResult, error) {
	s.called++
	s.lastCred = cred
	s.lastParam = param
	if s.err != nil {
		return nil, s.err
placeholder
	if s.result != nil {
		return s.result, nil
placeholder
	return &AliyunCaptchaVerifyResult{VerifyResult: trueplaceholder, nil
placeholder

func aliyunEnabledSettings() map[string]string {
	return map[string]string{
		SettingKeyAliyunCaptchaEnabled:         "true",
		SettingKeyAliyunCaptchaAccessKeyID:     "ak-id",
		SettingKeyAliyunCaptchaAccessKeySecret: "ak-secret",
		SettingKeyAliyunCaptchaSceneID:         "scene-1",
		SettingKeyAliyunCaptchaPrefix:          "prefix-1",
placeholder
placeholder

func aliyunTestConfig() AliyunCaptchaConfig {
	return AliyunCaptchaConfig{
		Enabled:         true,
		AccessKeyID:     "ak-id",
		AccessKeySecret: "ak-secret",
		SceneID:         "scene-1",
		Region:          AliyunCaptchaRegionCN,
placeholder
placeholder

func newAliyunAuthServiceForTest(cfg *config.Config, settings map[string]string, aliyunSpy *aliyunVerifierSpy) *AuthService {
	settingService := NewSettingService(&settingPublicRepoStub{values: settingsplaceholder, cfg)
	authService := NewAuthService(
		nil, // entClient
		nil, // userRepo
		nil, // redeemRepo
		nil, // refreshTokenCache
		cfg,
		settingService,
		nil, // emailService
		NewTurnstileService(settingService, &turnstileVerifierSpy{placeholder),
		nil, // emailQueueService
		nil, // promoService
		nil, // defaultSubAssigner
		nil, // affiliateService
		nil, // userPlatformQuotaRepo
	)
	authService.SetAliyunCaptchaService(NewAliyunCaptchaService(settingService, aliyunSpy))
	return authService
placeholder

func TestAliyunCaptchaServiceVerifyParamDispatch(t *testing.T) {
	spy := &aliyunVerifierSpy{placeholder
	svc := NewAliyunCaptchaService(nil, spy)

	err := svc.VerifyParamWithConfig(context.Background(), aliyunTestConfig(), "captcha-verify-param")

placeholder
	require.Equal(t, 1, spy.called)
	require.Equal(t, "captcha-verify-param", spy.lastParam)
	require.Equal(t, "ak-id", spy.lastCred.AccessKeyID)
	require.Equal(t, "scene-1", spy.lastCred.SceneID)
	require.Equal(t, "captcha.cn-shanghai.aliyuncs.com", spy.lastCred.Endpoint)
placeholder

func TestAliyunCaptchaServiceSgpEndpoint(t *testing.T) {
	spy := &aliyunVerifierSpy{placeholder
	svc := NewAliyunCaptchaService(nil, spy)
	cfg := aliyunTestConfig()
	cfg.Region = AliyunCaptchaRegionSGP

	err := svc.VerifyParamWithConfig(context.Background(), cfg, "captcha-verify-param")

placeholder
	require.Equal(t, "captcha.ap-southeast-1.aliyuncs.com", spy.lastCred.Endpoint)
placeholder

func TestAliyunCaptchaServiceFailsClosedOnVerifierError(t *testing.T) {
	spy := &aliyunVerifierSpy{err: errors.New("network down")placeholder
	svc := NewAliyunCaptchaService(nil, spy)

	err := svc.VerifyParamWithConfig(context.Background(), aliyunTestConfig(), "captcha-verify-param")

	require.ErrorIs(t, err, ErrAliyunCaptchaVerificationFailed)
placeholder

func TestAliyunCaptchaServiceRejectsVerifyResultFalse(t *testing.T) {
	spy := &aliyunVerifierSpy{result: &AliyunCaptchaVerifyResult{VerifyResult: false, VerifyCode: "F001"placeholderplaceholder
	svc := NewAliyunCaptchaService(nil, spy)

	err := svc.VerifyParamWithConfig(context.Background(), aliyunTestConfig(), "captcha-verify-param")

	require.ErrorIs(t, err, ErrAliyunCaptchaVerificationFailed)
placeholder

func TestAliyunCaptchaServiceRejectsIncompleteCredentials(t *testing.T) {
	spy := &aliyunVerifierSpy{placeholder
	svc := NewAliyunCaptchaService(nil, spy)
	cfg := aliyunTestConfig()
	cfg.AccessKeySecret = ""

	err := svc.VerifyParamWithConfig(context.Background(), cfg, "captcha-verify-param")

	require.ErrorIs(t, err, ErrAliyunCaptchaNotConfigured)
	require.Zero(t, spy.called)
placeholder

func TestAliyunCaptchaServiceRejectsEmptyParam(t *testing.T) {
	spy := &aliyunVerifierSpy{placeholder
	svc := NewAliyunCaptchaService(nil, spy)

	err := svc.VerifyParamWithConfig(context.Background(), aliyunTestConfig(), "")

	require.ErrorIs(t, err, ErrAliyunCaptchaVerificationFailed)
	require.Zero(t, spy.called)
placeholder

func TestAliyunCaptchaServiceValidateCredentials(t *testing.T) {
	t.Run("invalid credential code", func(t *testing.T) {
		spy := &aliyunVerifierSpy{err: &AliyunCaptchaAPIError{Code: "SignatureDoesNotMatch", Message: "bad sk"placeholderplaceholder
		svc := NewAliyunCaptchaService(nil, spy)

		err := svc.ValidateCredentials(context.Background(), "id", "sk", "scene", "cn")
		require.ErrorIs(t, err, ErrCaptchaInvalidCredentials)
placeholder)

	t.Run("network error surfaces", func(t *testing.T) {
		spy := &aliyunVerifierSpy{err: errors.New("timeout")placeholder
		svc := NewAliyunCaptchaService(nil, spy)

		err := svc.ValidateCredentials(context.Background(), "id", "sk", "scene", "cn")
	placeholder
		require.NotErrorIs(t, err, ErrCaptchaInvalidCredentials)
placeholder)

	t.Run("verify result false means credentials valid", func(t *testing.T) {
		spy := &aliyunVerifierSpy{result: &AliyunCaptchaVerifyResult{VerifyResult: falseplaceholderplaceholder
		svc := NewAliyunCaptchaService(nil, spy)

		err := svc.ValidateCredentials(context.Background(), "id", "sk", "scene", "sgp")
	placeholder
		require.Equal(t, "captcha.ap-southeast-1.aliyuncs.com", spy.lastCred.Endpoint)
placeholder)
placeholder

func TestAuthServiceVerifyCaptchaDispatchesAliyun(t *testing.T) {
	spy := &aliyunVerifierSpy{placeholder
	authService := newAliyunAuthServiceForTest(&config.Config{placeholder, aliyunEnabledSettings(), spy)

	// 阿里云 captchaVerifyParam 复用 turnstile_token 请求字段
	err := authService.VerifyCaptcha(context.Background(), CaptchaProof{TurnstileToken: "captcha-verify-param"placeholder, "127.0.0.1")

placeholder
	require.Equal(t, 1, spy.called)
	require.Equal(t, "captcha-verify-param", spy.lastParam)
placeholder

func TestAuthServiceVerifyCaptchaRejectsProviderConflict(t *testing.T) {
	settings := aliyunEnabledSettings()
	settings[SettingKeyTurnstileEnabled] = "true"
	settings[SettingKeyTurnstileSecretKey] = "secret"
	spy := &aliyunVerifierSpy{placeholder
	authService := newAliyunAuthServiceForTest(&config.Config{placeholder, settings, spy)

	err := authService.VerifyCaptcha(context.Background(), CaptchaProof{TurnstileToken: "param"placeholder, "127.0.0.1")

	require.ErrorIs(t, err, ErrCaptchaProviderConflict)
	require.Zero(t, spy.called)
placeholder

func TestAuthServiceVerifyCaptchaRequiredModeWithAliyun(t *testing.T) {
	cfg := &config.Config{
		Server:    config.ServerConfig{Mode: "release"placeholder,
		Turnstile: config.TurnstileConfig{Required: trueplaceholder,
placeholder
	spy := &aliyunVerifierSpy{placeholder
	authService := newAliyunAuthServiceForTest(cfg, aliyunEnabledSettings(), spy)

	// required 模式 + 阿里云启用且凭证齐全：不误报 NOT_CONFIGURED，正常走阿里云校验
	err := authService.VerifyCaptcha(context.Background(), CaptchaProof{TurnstileToken: "captcha-verify-param"placeholder, "127.0.0.1")

placeholder
	require.Equal(t, 1, spy.called)
placeholder

func TestAuthServiceVerifyActionCaptchaIfEnabledDispatchesAliyun(t *testing.T) {
	spy := &aliyunVerifierSpy{placeholder
	authService := newAliyunAuthServiceForTest(&config.Config{placeholder, aliyunEnabledSettings(), spy)

	err := authService.VerifyActionCaptchaIfEnabled(context.Background(), CaptchaProof{TurnstileToken: "captcha-verify-param"placeholder, "127.0.0.1")

placeholder
	require.Equal(t, 1, spy.called)
	require.Equal(t, "captcha-verify-param", spy.lastParam)
placeholder

func TestAuthServiceVerifyActionCaptchaIfEnabledSkipsWhenOnlyTurnstile(t *testing.T) {
	spy := &aliyunVerifierSpy{placeholder
	authService := newAliyunAuthServiceForTest(&config.Config{placeholder, map[string]string{
		SettingKeyTurnstileEnabled:   "true",
		SettingKeyTurnstileSecretKey: "secret",
placeholder, spy)

	// Turnstile 不扩大既有覆盖：扩展入口不拦截
	err := authService.VerifyActionCaptchaIfEnabled(context.Background(), CaptchaProof{placeholder, "127.0.0.1")

placeholder
	require.Zero(t, spy.called)
placeholder
