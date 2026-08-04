package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
)

var (
	ErrTencentCaptchaVerificationFailed = infraerrors.BadRequest("TENCENT_CAPTCHA_VERIFICATION_FAILED", "tencent captcha verification failed")
	ErrTencentCaptchaNotConfigured      = infraerrors.ServiceUnavailable("TENCENT_CAPTCHA_NOT_CONFIGURED", "tencent captcha not configured")
)

type TencentCaptchaProof struct {
	Ticket  string
	Randstr string
placeholder

type TencentCaptchaCredentials struct {
	AppID          uint64
	AppSecretKey   string
	CloudSecretID  string
	CloudSecretKey string
placeholder

type TencentCaptchaVerifyResponse struct {
	CaptchaCode int64
	CaptchaMsg  string
	RequestID   string
placeholder

type TencentCaptchaVerifier interface {
	VerifyTicket(context.Context, TencentCaptchaCredentials, TencentCaptchaProof, string) (*TencentCaptchaVerifyResponse, error)
placeholder

type TencentCaptchaService struct {
	settingService *SettingService
	verifier       TencentCaptchaVerifier
placeholder

func NewTencentCaptchaService(settingService *SettingService, verifier TencentCaptchaVerifier) *TencentCaptchaService {
	return &TencentCaptchaService{settingService: settingService, verifier: verifierplaceholder
placeholder

func (s *TencentCaptchaService) VerifyTicket(ctx context.Context, ticket, randstr, remoteIP string) error {
	if s == nil || s.settingService == nil {
		return ErrTencentCaptchaNotConfigured
placeholder
	providerConfig, err := s.settingService.GetCaptchaProviderConfig(ctx)
	if err != nil {
		logger.LegacyPrintf("service.tencent_captcha", "%s", "[TencentCaptcha] failed to read captcha provider settings")
		return ErrServiceUnavailable
placeholder
	config := providerConfig.Tencent
	if !config.Enabled {
		return nil
placeholder
	return s.VerifyTicketWithConfig(ctx, config, ticket, randstr, remoteIP)
placeholder

func (s *TencentCaptchaService) VerifyTicketWithConfig(ctx context.Context, config TencentCaptchaConfig, ticket, randstr, remoteIP string) error {
	credentials, ok := parseTencentCaptchaCredentials(config)
	if !ok || s.verifier == nil {
		return ErrTencentCaptchaNotConfigured
placeholder

	proof := TencentCaptchaProof{
		Ticket:  strings.TrimSpace(ticket),
		Randstr: strings.TrimSpace(randstr),
placeholder
	if proof.Ticket == "" || proof.Randstr == "" || strings.HasPrefix(proof.Ticket, "trerror_") {
		return ErrTencentCaptchaVerificationFailed
placeholder

	result, err := s.verifier.VerifyTicket(ctx, credentials, proof, remoteIP)
	if err != nil {
		logger.LegacyPrintf("service.tencent_captcha", "%s", "[TencentCaptcha] verification request failed")
		return fmt.Errorf("%w: verifier request failed", ErrTencentCaptchaVerificationFailed)
placeholder
	if result == nil || result.CaptchaCode != 1 {
		if result != nil {
			logger.LegacyPrintf("service.tencent_captcha", "[TencentCaptcha] rejected code=%d request_id=%s", result.CaptchaCode, result.RequestID)
	placeholder
		return ErrTencentCaptchaVerificationFailed
placeholder
	return nil
placeholder

func parseTencentCaptchaCredentials(config TencentCaptchaConfig) (TencentCaptchaCredentials, bool) {
	appID, err := strconv.ParseUint(strings.TrimSpace(config.AppID), 10, 64)
	if err != nil || appID == 0 {
		return TencentCaptchaCredentials{placeholder, false
placeholder
	credentials := TencentCaptchaCredentials{
		AppID:          appID,
		AppSecretKey:   strings.TrimSpace(config.AppSecretKey),
		CloudSecretID:  strings.TrimSpace(config.CloudSecretID),
		CloudSecretKey: strings.TrimSpace(config.CloudSecretKey),
placeholder
	if credentials.AppSecretKey == "" || credentials.CloudSecretID == "" || credentials.CloudSecretKey == "" {
		return TencentCaptchaCredentials{placeholder, false
placeholder
	return credentials, true
placeholder
