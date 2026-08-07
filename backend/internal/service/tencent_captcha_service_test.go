//go:build unit

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type tencentCaptchaVerifierStub struct {
	response    *TencentCaptchaVerifyResponse
	err         error
	calls       int
	proof       TencentCaptchaProof
	remoteIP    string
	credentials TencentCaptchaCredentials
placeholder

func (s *tencentCaptchaVerifierStub) VerifyTicket(_ context.Context, credentials TencentCaptchaCredentials, proof TencentCaptchaProof, remoteIP string) (*TencentCaptchaVerifyResponse, error) {
	s.calls++
	s.credentials = credentials
	s.proof = proof
	s.remoteIP = remoteIP
	return s.response, s.err
placeholder

func newTencentCaptchaTestService(verifier TencentCaptchaVerifier) *TencentCaptchaService {
	return newTencentCaptchaTestServiceWithRegion(verifier, "")
placeholder

func newTencentCaptchaTestServiceWithRegion(verifier TencentCaptchaVerifier, region string) *TencentCaptchaService {
	values := map[string]string{
		SettingKeyTencentCaptchaEnabled:        "true",
		SettingKeyTencentCaptchaAppID:          "123456789",
		SettingKeyTencentCaptchaAppSecretKey:   "app-secret",
		SettingKeyTencentCaptchaCloudSecretID:  "cloud-secret-id",
		SettingKeyTencentCaptchaCloudSecretKey: "cloud-secret-key",
placeholder
	if region != "" {
		values[SettingKeyTencentCaptchaRegion] = region
placeholder
	settings := NewSettingService(&settingPublicRepoStub{values: valuesplaceholder, &config.Config{placeholder)
	return NewTencentCaptchaService(settings, verifier)
placeholder

// 站点决定服务端票据校验接入点：国际站账号的密钥在国内站接入点上无法通过鉴权，
// 因此这条映射一旦错位，国际站验证码会整体失效。
func TestTencentCaptchaServiceRoutesVerifyEndpointByRegion(t *testing.T) {
	cases := []struct {
		name         string
		region       string
		wantEndpoint string
placeholder{
		{"未配置回落中国站", "", "captcha.tencentcloudapi.com"placeholder,
		{"中国站", TencentCaptchaRegionCN, "captcha.tencentcloudapi.com"placeholder,
		{"国际站", TencentCaptchaRegionINTL, "captcha.intl.tencentcloudapi.com"placeholder,
		{"非法值回落中国站", "sgp", "captcha.tencentcloudapi.com"placeholder,
placeholder
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			verifier := &tencentCaptchaVerifierStub{response: &TencentCaptchaVerifyResponse{CaptchaCode: 1placeholderplaceholder
			svc := newTencentCaptchaTestServiceWithRegion(verifier, tc.region)

			require.NoError(t, svc.VerifyTicket(context.Background(), "ticket", "@rand", "203.0.113.10"))
			require.Equal(t, tc.wantEndpoint, verifier.credentials.Endpoint)
	placeholder)
placeholder
placeholder

func TestTencentCaptchaServiceAcceptsCaptchaCodeOne(t *testing.T) {
	verifier := &tencentCaptchaVerifierStub{response: &TencentCaptchaVerifyResponse{CaptchaCode: 1placeholderplaceholder
	svc := newTencentCaptchaTestService(verifier)

	err := svc.VerifyTicket(context.Background(), "ticket", "@rand", "203.0.113.10")

placeholder
	require.Equal(t, 1, verifier.calls)
	require.Equal(t, TencentCaptchaProof{Ticket: "ticket", Randstr: "@rand"placeholder, verifier.proof)
	require.Equal(t, "203.0.113.10", verifier.remoteIP)
placeholder

func TestTencentCaptchaServiceRejectsDisasterRecoveryTicketWithoutCallingVerifier(t *testing.T) {
	verifier := &tencentCaptchaVerifierStub{response: &TencentCaptchaVerifyResponse{CaptchaCode: 1placeholderplaceholder
	svc := newTencentCaptchaTestService(verifier)

	err := svc.VerifyTicket(context.Background(), "trerror_1001_123456789_1", "@rand", "203.0.113.10")

	require.ErrorIs(t, err, ErrTencentCaptchaVerificationFailed)
	require.Zero(t, verifier.calls)
placeholder

func TestTencentCaptchaServiceRejectsEveryNonOneCode(t *testing.T) {
	for _, code := range []int64{0, 7, 8, 9, 15, 16, 21, 100placeholder {
		t.Run(string(rune(code)), func(t *testing.T) {
			verifier := &tencentCaptchaVerifierStub{response: &TencentCaptchaVerifyResponse{CaptchaCode: codeplaceholderplaceholder
			svc := newTencentCaptchaTestService(verifier)

			err := svc.VerifyTicket(context.Background(), "ticket", "@rand", "203.0.113.10")

			require.ErrorIs(t, err, ErrTencentCaptchaVerificationFailed)
	placeholder)
placeholder
placeholder

func TestTencentCaptchaServiceFailsClosedOnVerifierError(t *testing.T) {
	verifier := &tencentCaptchaVerifierStub{err: errors.New("sdk unavailable")placeholder
	svc := newTencentCaptchaTestService(verifier)

	err := svc.VerifyTicket(context.Background(), "ticket", "@rand", "203.0.113.10")

placeholder
	require.ErrorIs(t, err, ErrTencentCaptchaVerificationFailed)
placeholder

func TestTencentCaptchaServiceRejectsIncompleteConfiguration(t *testing.T) {
	settings := NewSettingService(&settingPublicRepoStub{values: map[string]string{
		SettingKeyTencentCaptchaEnabled: "true",
		SettingKeyTencentCaptchaAppID:   "123456789",
placeholderplaceholder, &config.Config{placeholder)
	verifier := &tencentCaptchaVerifierStub{response: &TencentCaptchaVerifyResponse{CaptchaCode: 1placeholderplaceholder
	svc := NewTencentCaptchaService(settings, verifier)

	err := svc.VerifyTicket(context.Background(), "ticket", "@rand", "203.0.113.10")

	require.ErrorIs(t, err, ErrTencentCaptchaNotConfigured)
	require.Zero(t, verifier.calls)
placeholder

func TestTencentCaptchaServiceFailsClosedOnSettingsReadError(t *testing.T) {
	settings := NewSettingService(&settingPublicRepoStub{err: errors.New("settings unavailable")placeholder, &config.Config{placeholder)
	verifier := &tencentCaptchaVerifierStub{response: &TencentCaptchaVerifyResponse{CaptchaCode: 1placeholderplaceholder
	svc := NewTencentCaptchaService(settings, verifier)

	err := svc.VerifyTicket(context.Background(), "ticket", "@rand", "203.0.113.10")

	require.ErrorIs(t, err, ErrServiceUnavailable)
	require.Zero(t, verifier.calls)
placeholder
