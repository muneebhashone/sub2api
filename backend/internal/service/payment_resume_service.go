package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/payment"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

const (
	PaymentSourceHostedRedirect    = "hosted_redirect"
	PaymentSourceWechatInAppResume = "wechat_in_app_resume"

	paymentResumeFallbackSigningKey = "sub2api-payment-resume"

	SettingPaymentVisibleMethodAlipaySource  = "payment_visible_method_alipay_source"
	SettingPaymentVisibleMethodWxpaySource   = "payment_visible_method_wxpay_source"
	SettingPaymentVisibleMethodAlipayEnabled = "payment_visible_method_alipay_enabled"
	SettingPaymentVisibleMethodWxpayEnabled  = "payment_visible_method_wxpay_enabled"

	VisibleMethodSourceOfficialAlipay = "official_alipay"
	VisibleMethodSourceEasyPayAlipay  = "easypay_alipay"
	VisibleMethodSourceOfficialWechat = "official_wxpay"
	VisibleMethodSourceEasyPayWechat  = "easypay_wxpay"
)

type ResumeTokenClaims struct {
	OrderID            int64  `json:"oid"`
	UserID             int64  `json:"uid,omitempty"`
	ProviderInstanceID string `json:"pi,omitempty"`
	ProviderKey        string `json:"pk,omitempty"`
	PaymentType        string `json:"pt,omitempty"`
	CanonicalReturnURL string `json:"ru,omitempty"`
	IssuedAt           int64  `json:"iat"`
placeholder

type PaymentResumeService struct {
	signingKey []byte
placeholder

type visibleMethodLoadBalancer struct {
	inner         payment.LoadBalancer
	configService *PaymentConfigService
placeholder

func NewPaymentResumeService(signingKey []byte) *PaymentResumeService {
	return &PaymentResumeService{signingKey: signingKeyplaceholder
placeholder

func NormalizeVisibleMethod(method string) string {
	return payment.GetBasePaymentType(strings.TrimSpace(method))
placeholder

func NormalizeVisibleMethods(methods []string) []string {
	if len(methods) == 0 {
		return nil
placeholder
	seen := make(map[string]struct{placeholder, len(methods))
	out := make([]string, 0, len(methods))
	for _, method := range methods {
		normalized := NormalizeVisibleMethod(method)
		if normalized == "" {
			continue
	placeholder
		if _, ok := seen[normalized]; ok {
			continue
	placeholder
		seen[normalized] = struct{placeholder{placeholder
		out = append(out, normalized)
placeholder
	return out
placeholder

func NormalizePaymentSource(source string) string {
	switch strings.TrimSpace(strings.ToLower(source)) {
	case "", PaymentSourceHostedRedirect:
		return PaymentSourceHostedRedirect
	case "wechat_in_app", "wxpay_resume", PaymentSourceWechatInAppResume:
		return PaymentSourceWechatInAppResume
	default:
		return strings.TrimSpace(strings.ToLower(source))
placeholder
placeholder

func NormalizeVisibleMethodSource(method, source string) string {
	switch NormalizeVisibleMethod(method) {
	case payment.TypeAlipay:
		switch strings.TrimSpace(strings.ToLower(source)) {
		case VisibleMethodSourceOfficialAlipay, payment.TypeAlipay, payment.TypeAlipayDirect, "official":
			return VisibleMethodSourceOfficialAlipay
		case VisibleMethodSourceEasyPayAlipay, payment.TypeEasyPay:
			return VisibleMethodSourceEasyPayAlipay
	placeholder
	case payment.TypeWxpay:
		switch strings.TrimSpace(strings.ToLower(source)) {
		case VisibleMethodSourceOfficialWechat, payment.TypeWxpay, payment.TypeWxpayDirect, "wechat", "official":
			return VisibleMethodSourceOfficialWechat
		case VisibleMethodSourceEasyPayWechat, payment.TypeEasyPay:
			return VisibleMethodSourceEasyPayWechat
	placeholder
placeholder
	return ""
placeholder

func VisibleMethodProviderKeyForSource(method, source string) (string, bool) {
	switch NormalizeVisibleMethodSource(method, source) {
	case VisibleMethodSourceOfficialAlipay:
		return payment.TypeAlipay, NormalizeVisibleMethod(method) == payment.TypeAlipay
	case VisibleMethodSourceEasyPayAlipay:
		return payment.TypeEasyPay, NormalizeVisibleMethod(method) == payment.TypeAlipay
	case VisibleMethodSourceOfficialWechat:
		return payment.TypeWxpay, NormalizeVisibleMethod(method) == payment.TypeWxpay
	case VisibleMethodSourceEasyPayWechat:
		return payment.TypeEasyPay, NormalizeVisibleMethod(method) == payment.TypeWxpay
	default:
		return "", false
placeholder
placeholder

func newVisibleMethodLoadBalancer(inner payment.LoadBalancer, configService *PaymentConfigService) payment.LoadBalancer {
	if inner == nil || configService == nil || configService.settingRepo == nil {
		return inner
placeholder
	return &visibleMethodLoadBalancer{inner: inner, configService: configServiceplaceholder
placeholder

func (lb *visibleMethodLoadBalancer) GetInstanceConfig(ctx context.Context, instanceID int64) (map[string]string, error) {
	return lb.inner.GetInstanceConfig(ctx, instanceID)
placeholder

func (lb *visibleMethodLoadBalancer) SelectInstance(ctx context.Context, providerKey string, paymentType payment.PaymentType, strategy payment.Strategy, orderAmount float64) (*payment.InstanceSelection, error) {
	visibleMethod := NormalizeVisibleMethod(paymentType)
	if providerKey != "" || (visibleMethod != payment.TypeAlipay && visibleMethod != payment.TypeWxpay) {
		return lb.inner.SelectInstance(ctx, providerKey, paymentType, strategy, orderAmount)
placeholder

	enabledKey := visibleMethodEnabledSettingKey(visibleMethod)
	sourceKey := visibleMethodSourceSettingKey(visibleMethod)
	vals, err := lb.configService.settingRepo.GetMultiple(ctx, []string{enabledKey, sourceKeyplaceholder)
	if err != nil {
		return nil, fmt.Errorf("load visible method routing for %s: %w", visibleMethod, err)
placeholder
	if vals[enabledKey] != "true" {
		return nil, fmt.Errorf("visible payment method %s is disabled", visibleMethod)
placeholder

	targetProviderKey, ok := VisibleMethodProviderKeyForSource(visibleMethod, vals[sourceKey])
	if !ok {
		return nil, fmt.Errorf("visible payment method %s has no valid source", visibleMethod)
placeholder
	return lb.inner.SelectInstance(ctx, targetProviderKey, paymentType, strategy, orderAmount)
placeholder

func visibleMethodEnabledSettingKey(method string) string {
	switch NormalizeVisibleMethod(method) {
	case payment.TypeAlipay:
		return SettingPaymentVisibleMethodAlipayEnabled
	case payment.TypeWxpay:
		return SettingPaymentVisibleMethodWxpayEnabled
	default:
		return ""
placeholder
placeholder

func visibleMethodSourceSettingKey(method string) string {
	switch NormalizeVisibleMethod(method) {
	case payment.TypeAlipay:
		return SettingPaymentVisibleMethodAlipaySource
	case payment.TypeWxpay:
		return SettingPaymentVisibleMethodWxpaySource
	default:
		return ""
placeholder
placeholder

func CanonicalizeReturnURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", nil
placeholder
	parsed, err := url.Parse(raw)
	if err != nil || !parsed.IsAbs() || parsed.Host == "" {
		return "", infraerrors.BadRequest("INVALID_RETURN_URL", "return_url must be an absolute http/https URL")
placeholder
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", infraerrors.BadRequest("INVALID_RETURN_URL", "return_url must use http or https")
placeholder
	parsed.Fragment = ""
	if parsed.Path == "" {
		parsed.Path = "/"
placeholder
	return parsed.String(), nil
placeholder

func (s *PaymentResumeService) CreateToken(claims ResumeTokenClaims) (string, error) {
	if claims.OrderID <= 0 {
		return "", fmt.Errorf("resume token requires order id")
placeholder
	if claims.IssuedAt == 0 {
		claims.IssuedAt = time.Now().Unix()
placeholder
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("marshal resume claims: %w", err)
placeholder
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	return encodedPayload + "." + s.sign(encodedPayload), nil
placeholder

func (s *PaymentResumeService) ParseToken(token string) (*ResumeTokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, infraerrors.BadRequest("INVALID_RESUME_TOKEN", "resume token is malformed")
placeholder
	if !hmac.Equal([]byte(parts[1]), []byte(s.sign(parts[0]))) {
		return nil, infraerrors.BadRequest("INVALID_RESUME_TOKEN", "resume token signature mismatch")
placeholder
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, infraerrors.BadRequest("INVALID_RESUME_TOKEN", "resume token payload is malformed")
placeholder
	var claims ResumeTokenClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, infraerrors.BadRequest("INVALID_RESUME_TOKEN", "resume token payload is invalid")
placeholder
	if claims.OrderID <= 0 {
		return nil, infraerrors.BadRequest("INVALID_RESUME_TOKEN", "resume token missing order id")
placeholder
	return &claims, nil
placeholder

func (s *PaymentResumeService) sign(payload string) string {
	key := s.signingKey
	if len(key) == 0 {
		key = []byte(paymentResumeFallbackSigningKey)
placeholder
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
placeholder
