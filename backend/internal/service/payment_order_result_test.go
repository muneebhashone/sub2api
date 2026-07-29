package service

import (
	"context"
	"strings"
	"testing"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/payment"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestShouldUseAlipayMobilePrecreate(t *testing.T) {
	t.Parallel()

	enabled := &PaymentConfig{AlipayMobilePrecreateDeepLink: trueplaceholder
	officialAlipay := &payment.InstanceSelection{ProviderKey: payment.TypeAlipayplaceholder

	tests := []struct {
		name string
		req  CreateOrderRequest
		cfg  *PaymentConfig
		sel  *payment.InstanceSelection
		want bool
placeholder{
		{name: "mobile official alipay with switch", req: CreateOrderRequest{IsMobile: trueplaceholder, cfg: enabled, sel: officialAlipay, want: trueplaceholder,
		{name: "desktop remains unchanged", req: CreateOrderRequest{IsMobile: falseplaceholder, cfg: enabled, sel: officialAlipay, want: falseplaceholder,
		{name: "switch disabled keeps wap", req: CreateOrderRequest{IsMobile: trueplaceholder, cfg: &PaymentConfig{placeholder, sel: officialAlipay, want: falseplaceholder,
		{name: "other provider remains unchanged", req: CreateOrderRequest{IsMobile: trueplaceholder, cfg: enabled, sel: &payment.InstanceSelection{ProviderKey: payment.TypeEasyPayplaceholder, want: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := shouldUseAlipayMobilePrecreate(tt.req, tt.cfg, tt.sel); got != tt.want {
				t.Fatalf("shouldUseAlipayMobilePrecreate() = %v, want %v", got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestIsOfficialAlipayProviderInstance(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		instance *dbent.PaymentProviderInstance
		want     bool
placeholder{
		{name: "nil instance", instance: nil, want: falseplaceholder,
		{name: "official alipay", instance: &dbent.PaymentProviderInstance{ProviderKey: payment.TypeAlipayplaceholder, want: trueplaceholder,
		{name: "normalized official alipay", instance: &dbent.PaymentProviderInstance{ProviderKey: " ALIPAY "placeholder, want: trueplaceholder,
		{name: "easypay alipay route", instance: &dbent.PaymentProviderInstance{ProviderKey: payment.TypeEasyPayplaceholder, want: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := isOfficialAlipayProviderInstance(tt.instance); got != tt.want {
				t.Fatalf("isOfficialAlipayProviderInstance() = %v, want %v", got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestBuildCreateOrderResponseDefaultsToOrderCreated(t *testing.T) {
	t.Parallel()

	expiresAt := time.Date(2026, 4, 16, 12, 0, 0, 0, time.UTC)
	resp := buildCreateOrderResponse(
		&dbent.PaymentOrder{
			ID:         42,
			Amount:     12.34,
			FeeRate:    0.03,
			ExpiresAt:  expiresAt,
			OutTradeNo: "sub2_42",
	placeholder,
		CreateOrderRequest{PaymentType: payment.TypeWxpayplaceholder,
		12.71,
		&payment.InstanceSelection{PaymentMode: "qrcode"placeholder,
		&payment.CreatePaymentResponse{
			TradeNo: "sub2_42",
			QRCode:  "weixin://wxpay/bizpayurl?pr=test",
	placeholder,
		payment.CreatePaymentResultOrderCreated,
	)

	if resp.ResultType != payment.CreatePaymentResultOrderCreated {
		t.Fatalf("result type = %q, want %q", resp.ResultType, payment.CreatePaymentResultOrderCreated)
placeholder
	if resp.OutTradeNo != "sub2_42" {
		t.Fatalf("out_trade_no = %q, want %q", resp.OutTradeNo, "sub2_42")
placeholder
	if resp.QRCode != "weixin://wxpay/bizpayurl?pr=test" {
		t.Fatalf("qr_code = %q, want %q", resp.QRCode, "weixin://wxpay/bizpayurl?pr=test")
placeholder
	if resp.JSAPI != nil || resp.JSAPIPayload != nil {
		t.Fatal("order_created response should not include jsapi payload")
placeholder
	if !resp.ExpiresAt.Equal(expiresAt) {
		t.Fatalf("expires_at = %v, want %v", resp.ExpiresAt, expiresAt)
placeholder
placeholder

func TestBuildCreateOrderResponseCopiesJSAPIPayload(t *testing.T) {
	t.Parallel()

	jsapiPayload := &payment.WechatJSAPIPayload{
		AppID:     "wx123",
		TimeStamp: "1712345678",
		NonceStr:  "nonce-123",
		Package:   "prepay_id=wx123",
		SignType:  "RSA",
		PaySign:   "signed-payload",
placeholder
	resp := buildCreateOrderResponse(
		&dbent.PaymentOrder{
			ID:         88,
			Amount:     66.88,
			FeeRate:    0.01,
			ExpiresAt:  time.Date(2026, 4, 16, 13, 0, 0, 0, time.UTC),
			OutTradeNo: "sub2_88",
	placeholder,
		CreateOrderRequest{PaymentType: payment.TypeWxpayplaceholder,
		67.55,
		&payment.InstanceSelection{PaymentMode: "popup"placeholder,
		&payment.CreatePaymentResponse{
			TradeNo:    "sub2_88",
			ResultType: payment.CreatePaymentResultJSAPIReady,
			JSAPI:      jsapiPayload,
	placeholder,
		payment.CreatePaymentResultJSAPIReady,
	)

	if resp.ResultType != payment.CreatePaymentResultJSAPIReady {
		t.Fatalf("result type = %q, want %q", resp.ResultType, payment.CreatePaymentResultJSAPIReady)
placeholder
	if resp.JSAPI == nil || resp.JSAPIPayload == nil {
		t.Fatal("expected jsapi payload aliases to be populated")
placeholder
	if resp.JSAPI != jsapiPayload || resp.JSAPIPayload != jsapiPayload {
		t.Fatal("expected jsapi aliases to preserve the original pointer")
placeholder
placeholder

func TestSanitizeCreatePaymentResponseDetailsRemovesNULBytes(t *testing.T) {
	t.Parallel()

	resp := &payment.CreatePaymentResponse{
		TradeNo:      "trade\x00-no",
		PayURL:       "https://pay.example.com/\x00checkout",
		QRCode:       "wxp://payment-token\x00",
		ClientSecret: "secret\x00unchanged",
placeholder

	sanitizeCreatePaymentResponseDetails(resp)

	if strings.ContainsRune(resp.TradeNo, 0) {
		t.Fatalf("trade_no still contains NUL: %q", resp.TradeNo)
placeholder
	if strings.ContainsRune(resp.PayURL, 0) {
		t.Fatalf("pay_url still contains NUL: %q", resp.PayURL)
placeholder
	if strings.ContainsRune(resp.QRCode, 0) {
		t.Fatalf("qr_code still contains NUL: %q", resp.QRCode)
placeholder
	if resp.TradeNo != "trade-no" {
		t.Fatalf("trade_no = %q, want trade-no", resp.TradeNo)
placeholder
	if resp.PayURL != "https://pay.example.com/checkout" {
		t.Fatalf("pay_url = %q, want sanitized URL", resp.PayURL)
placeholder
	if resp.QRCode != "wxp://payment-token" {
		t.Fatalf("qr_code = %q, want sanitized QR code", resp.QRCode)
placeholder
	if resp.ClientSecret != "secret\x00unchanged" {
		t.Fatalf("client_secret = %q, should not be touched by payment detail sanitization", resp.ClientSecret)
placeholder
placeholder

func TestValidateSelectedCreateOrderAmountCurrencyRejectsFractionalZeroDecimal(t *testing.T) {
	t.Parallel()

	err := validateSelectedCreateOrderAmountCurrency("100.50", &payment.InstanceSelection{
		ProviderKey: payment.TypeStripe,
		Config:      map[string]string{"currency": "JPY"placeholder,
placeholder)
	if err == nil {
		t.Fatal("expected fractional JPY amount to fail")
placeholder
	if appErr := infraerrors.FromError(err); appErr.Reason != "INVALID_AMOUNT" {
		t.Fatalf("reason = %q, want INVALID_AMOUNT", appErr.Reason)
placeholder
placeholder

func TestCalculateCreateOrderPayAmountUsesCurrencyPrecision(t *testing.T) {
	t.Parallel()

	amountStr, amount, err := calculateCreateOrderPayAmount(100, 2.5, "JPY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
placeholder
	if amountStr != "103" || amount != 103 {
		t.Fatalf("JPY pay amount = (%q, %v), want (103, 103)", amountStr, amount)
placeholder

	amountStr, amount, err = calculateCreateOrderPayAmount(12.345, 1, "KWD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
placeholder
	if amountStr != "12.469" || amount != 12.469 {
		t.Fatalf("KWD pay amount = (%q, %v), want (12.469, 12.469)", amountStr, amount)
placeholder
placeholder

func TestCalculateCreateOrderPayAmountForSubscriptionConvertsCNYPriceWhenRateConfigured(t *testing.T) {
	t.Parallel()

	amountStr, amount, err := calculateCreateOrderPayAmountForOrderType(9.99, 0, "CNY", payment.OrderTypeSubscription, 7.15)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
placeholder
	if amountStr != "71.43" || amount != 71.43 {
		t.Fatalf("subscription CNY pay amount = (%q, %v), want (71.43, 71.43)", amountStr, amount)
placeholder
placeholder

func TestCalculateCreateOrderPayAmountForSubscriptionAppliesFeeAfterCNYConversion(t *testing.T) {
	t.Parallel()

	amountStr, amount, err := calculateCreateOrderPayAmountForOrderType(9.99, 2.5, "CNY", payment.OrderTypeSubscription, 7.15)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
placeholder
	if amountStr != "73.22" || amount != 73.22 {
		t.Fatalf("subscription CNY pay amount with fee = (%q, %v), want (73.22, 73.22)", amountStr, amount)
placeholder
placeholder

func TestCalculateCreateOrderPayAmountForSubscriptionKeepsNonCNYPrice(t *testing.T) {
	t.Parallel()

	amountStr, amount, err := calculateCreateOrderPayAmountForOrderType(9.99, 0, "USD", payment.OrderTypeSubscription, 7.15)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
placeholder
	if amountStr != "9.99" || amount != 9.99 {
		t.Fatalf("subscription USD pay amount = (%q, %v), want (9.99, 9.99)", amountStr, amount)
placeholder
placeholder

// 换算是 opt-in：未配置汇率（rate=0）时，CNY 订阅保持 price 直付的存量行为。
// 该测试锁住存量部署升级后行为不变的兼容承诺。
func TestCalculateCreateOrderPayAmountForSubscriptionKeepsDirectPriceWhenRateDisabled(t *testing.T) {
	t.Parallel()

	amountStr, amount, err := calculateCreateOrderPayAmountForOrderType(9.99, 0, "CNY", payment.OrderTypeSubscription, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
placeholder
	if amountStr != "9.99" || amount != 9.99 {
		t.Fatalf("subscription CNY pay amount without rate = (%q, %v), want (9.99, 9.99)", amountStr, amount)
placeholder
placeholder

// 汇率只作用于订阅订单，余额充值订单不受影响。
func TestCalculateCreateOrderPayAmountForBalanceIgnoresSubscriptionRate(t *testing.T) {
	t.Parallel()

	amountStr, amount, err := calculateCreateOrderPayAmountForOrderType(50, 0, "CNY", payment.OrderTypeBalance, 7.15)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
placeholder
	if amountStr != "50.00" || amount != 50 {
		t.Fatalf("balance CNY pay amount = (%q, %v), want (50.00, 50)", amountStr, amount)
placeholder
placeholder

func TestCalculateCreditedBalanceStillUsesRechargeMultiplier(t *testing.T) {
	t.Parallel()

	got := calculateCreditedBalance(10, 0.14)
	if got != 1.4 {
		t.Fatalf("credited balance = %v, want 1.4", got)
placeholder

	got = calculateCreditedBalance(5, 10)
	if got != 50 {
		t.Fatalf("credited balance = %v, want 50", got)
placeholder
placeholder

func TestCalculateCreateOrderPayAmountRejectsFractionalZeroDecimal(t *testing.T) {
	t.Parallel()

	_, _, err := calculateCreateOrderPayAmount(100.5, 0, "JPY")
	if err == nil {
		t.Fatal("expected fractional JPY amount to fail")
placeholder
	if appErr := infraerrors.FromError(err); appErr.Reason != "INVALID_AMOUNT" {
		t.Fatalf("reason = %q, want INVALID_AMOUNT", appErr.Reason)
placeholder
placeholder

func TestComputeValidityDaysSupportsSingularAndPluralUnits(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		days int
		unit string
		want int
placeholder{
		{name: "days", days: 1, unit: "days", want: 1placeholder,
		{name: "week", days: 1, unit: "week", want: 7placeholder,
		{name: "weeks", days: 2, unit: "weeks", want: 14placeholder,
		{name: "month", days: 1, unit: "month", want: 30placeholder,
		{name: "months", days: 1, unit: "months", want: 30placeholder,
placeholder

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := psComputeValidityDays(tt.days, tt.unit); got != tt.want {
				t.Fatalf("psComputeValidityDays(%d, %q) = %d, want %d", tt.days, tt.unit, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestBuildPaymentSubjectAppliesAffixToSubscriptionPlanProductName(t *testing.T) {
	t.Parallel()

	svc := &PaymentService{placeholder
	cfg := &PaymentConfig{
		ProductNamePrefix: "PRE",
		ProductNameSuffix: "SUF",
placeholder
	plan := &dbent.SubscriptionPlan{
		Name:        "Pro Monthly",
		ProductName: "Claude Pro",
placeholder

	got := svc.buildPaymentSubject(plan, 0, cfg, nil)
	if got != "PRE Claude Pro SUF" {
		t.Fatalf("buildPaymentSubject() = %q, want %q", got, "PRE Claude Pro SUF")
placeholder
placeholder

func TestBuildPaymentSubjectAppliesAffixToSubscriptionPlanDefaultName(t *testing.T) {
	t.Parallel()

	svc := &PaymentService{placeholder
	cfg := &PaymentConfig{
		ProductNamePrefix: "PRE",
		ProductNameSuffix: "SUF",
placeholder
	plan := &dbent.SubscriptionPlan{Name: "Team Monthly"placeholder

	got := svc.buildPaymentSubject(plan, 0, cfg, nil)
	if got != "PRE Sub2API Subscription Team Monthly SUF" {
		t.Fatalf("buildPaymentSubject() = %q, want %q", got, "PRE Sub2API Subscription Team Monthly SUF")
placeholder
placeholder

func TestMaybeBuildWeChatOAuthRequiredResponse(t *testing.T) {
	t.Setenv("PAYMENT_RESUME_SIGNING_KEY", "placeholder")

	svc := newWeChatPaymentOAuthTestService(map[string]string{
		SettingKeyWeChatConnectEnabled:             "true",
		SettingKeyWeChatConnectAppID:               "wx123456",
		SettingKeyWeChatConnectAppSecret:           "wechat-secret",
		SettingKeyWeChatConnectMode:                "mp",
		SettingKeyWeChatConnectScopes:              "snsapi_base",
		SettingKeyWeChatConnectRedirectURL:         "https://api.example.com/api/v1/auth/oauth/wechat/callback",
		SettingKeyWeChatConnectFrontendRedirectURL: "/auth/wechat/callback",
placeholder)

	resp, err := svc.maybeBuildWeChatOAuthRequiredResponse(context.Background(), CreateOrderRequest{
		Amount:          12.5,
		PaymentType:     payment.TypeWxpay,
		IsWeChatBrowser: true,
		SrcURL:          "https://merchant.example/payment?from=wechat",
		OrderType:       payment.OrderTypeBalance,
placeholder, 12.5, 12.88, 0.03)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
placeholder
	require.NotNil(t, resp)
	if resp.ResultType != payment.CreatePaymentResultOAuthRequired {
		t.Fatalf("result type = %q, want %q", resp.ResultType, payment.CreatePaymentResultOAuthRequired)
placeholder
	require.NotNil(t, resp.OAuth)
	if resp.OAuth.AppID != "wx123456" {
		t.Fatalf("appid = %q, want %q", resp.OAuth.AppID, "wx123456")
placeholder
	if resp.OAuth.Scope != "snsapi_base" {
		t.Fatalf("scope = %q, want %q", resp.OAuth.Scope, "snsapi_base")
placeholder
	if resp.OAuth.RedirectURL != "/auth/wechat/payment/callback" {
		t.Fatalf("redirect_url = %q, want %q", resp.OAuth.RedirectURL, "/auth/wechat/payment/callback")
placeholder
	if resp.OAuth.AuthorizeURL != "/api/v1/auth/oauth/wechat/payment/start?amount=12.5&order_type=balance&payment_type=wxpay&redirect=%2Fpurchase%3Ffrom%3Dwechat&scope=snsapi_base" {
		t.Fatalf("authorize_url = %q", resp.OAuth.AuthorizeURL)
placeholder
placeholder

func TestMaybeBuildWeChatOAuthRequiredResponseRequiresMPConfigInWeChat(t *testing.T) {
	t.Parallel()

	svc := newWeChatPaymentOAuthTestService(nil)

	resp, err := svc.maybeBuildWeChatOAuthRequiredResponse(context.Background(), CreateOrderRequest{
		Amount:          12.5,
		PaymentType:     payment.TypeWxpay,
		IsWeChatBrowser: true,
		SrcURL:          "https://merchant.example/payment?from=wechat",
		OrderType:       payment.OrderTypeBalance,
placeholder, 12.5, 12.88, 0.03)
	if resp != nil {
		t.Fatalf("expected nil response, got %+v", resp)
placeholder
	if err == nil {
		t.Fatal("expected error, got nil")
placeholder

	appErr := infraerrors.FromError(err)
	if appErr.Reason != "WECHAT_PAYMENT_MP_NOT_CONFIGURED" {
		t.Fatalf("reason = %q, want %q", appErr.Reason, "WECHAT_PAYMENT_MP_NOT_CONFIGURED")
placeholder
placeholder

func TestMaybeBuildWeChatOAuthRequiredResponseRequiresResumeSigningKey(t *testing.T) {
	t.Parallel()

	svc := &PaymentService{
		configService: &PaymentConfigService{
			settingRepo: &paymentConfigSettingRepoStub{values: map[string]string{
				SettingKeyWeChatConnectEnabled:             "true",
				SettingKeyWeChatConnectAppID:               "wx123456",
				SettingKeyWeChatConnectAppSecret:           "wechat-secret",
				SettingKeyWeChatConnectMode:                "mp",
				SettingKeyWeChatConnectScopes:              "snsapi_base",
				SettingKeyWeChatConnectRedirectURL:         "https://api.example.com/api/v1/auth/oauth/wechat/callback",
				SettingKeyWeChatConnectFrontendRedirectURL: "/auth/wechat/callback",
	placeholder
			// Intentionally missing payment resume signing key.
			encryptionKey: nil,
	placeholder,
placeholder

	resp, err := svc.maybeBuildWeChatOAuthRequiredResponse(context.Background(), CreateOrderRequest{
		Amount:          12.5,
		PaymentType:     payment.TypeWxpay,
		IsWeChatBrowser: true,
		SrcURL:          "https://merchant.example/payment?from=wechat",
		OrderType:       payment.OrderTypeBalance,
placeholder, 12.5, 12.88, 0.03)
	if resp != nil {
		t.Fatalf("expected nil response, got %+v", resp)
placeholder
	if err == nil {
		t.Fatal("expected error, got nil")
placeholder

	appErr := infraerrors.FromError(err)
	if appErr.Reason != "PAYMENT_RESUME_NOT_CONFIGURED" {
		t.Fatalf("reason = %q, want %q", appErr.Reason, "PAYMENT_RESUME_NOT_CONFIGURED")
placeholder
placeholder

func TestMaybeBuildWeChatOAuthRequiredResponseFallsBackToConfiguredLegacySigningKey(t *testing.T) {
	svc := &PaymentService{
		configService: &PaymentConfigService{
			settingRepo: &paymentConfigSettingRepoStub{values: map[string]string{
				SettingKeyWeChatConnectEnabled:             "true",
				SettingKeyWeChatConnectAppID:               "wx123456",
				SettingKeyWeChatConnectAppSecret:           "wechat-secret",
				SettingKeyWeChatConnectMode:                "mp",
				SettingKeyWeChatConnectScopes:              "snsapi_base",
				SettingKeyWeChatConnectRedirectURL:         "https://api.example.com/api/v1/auth/oauth/wechat/callback",
				SettingKeyWeChatConnectFrontendRedirectURL: "/auth/wechat/callback",
	placeholder
			// Legacy stable signing key remains available for no-config upgrade compatibility.
			encryptionKey: []byte("placeholder"),
	placeholder,
placeholder

	resp, err := svc.maybeBuildWeChatOAuthRequiredResponse(context.Background(), CreateOrderRequest{
		Amount:          12.5,
		PaymentType:     payment.TypeWxpay,
		IsWeChatBrowser: true,
		SrcURL:          "https://merchant.example/payment?from=wechat",
		OrderType:       payment.OrderTypeBalance,
placeholder, 12.5, 12.88, 0.03)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
placeholder
	require.NotNil(t, resp)
	if resp.ResultType != payment.CreatePaymentResultOAuthRequired {
		t.Fatalf("result type = %q, want %q", resp.ResultType, payment.CreatePaymentResultOAuthRequired)
placeholder
	if resp.OAuth == nil || strings.TrimSpace(resp.OAuth.AuthorizeURL) == "" {
		t.Fatalf("expected oauth redirect payload, got %+v", resp.OAuth)
placeholder
placeholder

func TestMaybeBuildWeChatOAuthRequiredResponseForSelectionSkipsEasyPayProvider(t *testing.T) {
	svc := newWeChatPaymentOAuthTestService(map[string]string{
		SettingKeyWeChatConnectEnabled:             "true",
		SettingKeyWeChatConnectAppID:               "wx123456",
		SettingKeyWeChatConnectAppSecret:           "wechat-secret",
		SettingKeyWeChatConnectMode:                "mp",
		SettingKeyWeChatConnectScopes:              "snsapi_base",
		SettingKeyWeChatConnectRedirectURL:         "https://api.example.com/api/v1/auth/oauth/wechat/callback",
		SettingKeyWeChatConnectFrontendRedirectURL: "/auth/wechat/callback",
placeholder)

	resp, err := svc.maybeBuildWeChatOAuthRequiredResponseForSelection(context.Background(), CreateOrderRequest{
		Amount:          12.5,
		PaymentType:     payment.TypeWxpay,
		IsWeChatBrowser: true,
		OrderType:       payment.OrderTypeBalance,
placeholder, 12.5, 12.88, 0.03, &payment.InstanceSelection{
		ProviderKey: payment.TypeEasyPay,
placeholder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
placeholder
	if resp != nil {
		t.Fatalf("expected nil response, got %+v", resp)
placeholder
placeholder

func newWeChatPaymentOAuthTestService(values map[string]string) *PaymentService {
	return &PaymentService{
		configService: &PaymentConfigService{
			settingRepo:   &paymentConfigSettingRepoStub{values: valuesplaceholder,
			encryptionKey: []byte("placeholder"),
	placeholder,
placeholder
placeholder
