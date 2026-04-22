//go:build unit

package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/payment"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

func TestNormalizeVisibleMethods(t *testing.T) {
	t.Parallel()

	got := NormalizeVisibleMethods([]string{
		"alipay_direct",
		"alipay",
		" wxpay_direct ",
		"wxpay",
		"stripe",
placeholder)

	want := []string{"alipay", "wxpay", "stripe"placeholder
	if len(got) != len(want) {
		t.Fatalf("NormalizeVisibleMethods len = %d, want %d (%v)", len(got), len(want), got)
placeholder
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("NormalizeVisibleMethods[%d] = %q, want %q (full=%v)", i, got[i], want[i], got)
	placeholder
placeholder
placeholder

func TestNormalizePaymentSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		expect string
placeholder{
		{name: "empty uses default", input: "", expect: PaymentSourceHostedRedirectplaceholder,
		{name: "wechat alias normalized", input: "wechat_in_app", expect: PaymentSourceWechatInAppResumeplaceholder,
		{name: "canonical value preserved", input: PaymentSourceWechatInAppResume, expect: PaymentSourceWechatInAppResumeplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NormalizePaymentSource(tt.input); got != tt.expect {
				t.Fatalf("NormalizePaymentSource(%q) = %q, want %q", tt.input, got, tt.expect)
		placeholder
	placeholder)
placeholder
placeholder

func TestCanonicalizeReturnURL(t *testing.T) {
	t.Parallel()

	got, err := CanonicalizeReturnURL("https://example.com/payment/result?b=2#a", "example.com", "")
	if err != nil {
		t.Fatalf("CanonicalizeReturnURL returned error: %v", err)
placeholder
	if got != "https://example.com/payment/result?b=2" {
		t.Fatalf("CanonicalizeReturnURL = %q, want %q", got, "https://example.com/payment/result?b=2")
placeholder
placeholder

func TestCanonicalizeReturnURLRejectsRelativeURL(t *testing.T) {
	t.Parallel()

	if _, err := CanonicalizeReturnURL("/payment/result", "example.com", ""); err == nil {
		t.Fatal("CanonicalizeReturnURL should reject relative URLs")
placeholder
placeholder

func TestCanonicalizeReturnURLRejectsExternalHost(t *testing.T) {
	t.Parallel()

	if _, err := CanonicalizeReturnURL("https://evil.example/payment/result", "app.example.com", ""); err == nil {
		t.Fatal("CanonicalizeReturnURL should reject external hosts")
placeholder
placeholder

func TestCanonicalizeReturnURLAllowsConfiguredFrontendHost(t *testing.T) {
	t.Parallel()

	got, err := CanonicalizeReturnURL(
		"https://app.example.com/payment/result?from=checkout",
		"api.example.com",
		"https://app.example.com/purchase",
	)
	if err != nil {
		t.Fatalf("CanonicalizeReturnURL returned error: %v", err)
placeholder
	if got != "https://app.example.com/payment/result?from=checkout" {
		t.Fatalf("CanonicalizeReturnURL = %q, want %q", got, "https://app.example.com/payment/result?from=checkout")
placeholder
placeholder

func TestCanonicalizeReturnURLRejectsNonCanonicalPath(t *testing.T) {
	t.Parallel()

	if _, err := CanonicalizeReturnURL("https://app.example.com/orders/42", "app.example.com", ""); err == nil {
		t.Fatal("CanonicalizeReturnURL should reject non-canonical result paths")
placeholder
placeholder

func TestBuildPaymentReturnURL(t *testing.T) {
	t.Parallel()

	got, err := buildPaymentReturnURL("https://example.com/payment/result?from=checkout#fragment", 42, "sub2_42", "resume-token")
	if err != nil {
		t.Fatalf("buildPaymentReturnURL returned error: %v", err)
placeholder

	parsed, err := url.Parse(got)
	if err != nil {
		t.Fatalf("url.Parse returned error: %v", err)
placeholder
	if parsed.Fragment != "" {
		t.Fatalf("buildPaymentReturnURL should strip fragments, got %q", parsed.Fragment)
placeholder
	query := parsed.Query()
	if query.Get("from") != "checkout" {
		t.Fatalf("expected original query to be preserved, got %q", query.Get("from"))
placeholder
	if query.Get("order_id") != strconv.FormatInt(42, 10) {
		t.Fatalf("order_id = %q", query.Get("order_id"))
placeholder
	if query.Get("out_trade_no") != "sub2_42" {
		t.Fatalf("out_trade_no = %q", query.Get("out_trade_no"))
placeholder
	if query.Get("resume_token") != "resume-token" {
		t.Fatalf("resume_token = %q", query.Get("resume_token"))
placeholder
	if query.Get("status") != "success" {
		t.Fatalf("status = %q", query.Get("status"))
placeholder
placeholder

func TestBuildPaymentReturnURLWithoutResumeTokenStillIncludesOutTradeNo(t *testing.T) {
	t.Parallel()

	got, err := buildPaymentReturnURL("https://example.com/payment/result", 42, "sub2_42", "")
	if err != nil {
		t.Fatalf("buildPaymentReturnURL returned error: %v", err)
placeholder

	parsed, err := url.Parse(got)
	if err != nil {
		t.Fatalf("url.Parse returned error: %v", err)
placeholder
	query := parsed.Query()
	if query.Get("order_id") != "42" {
		t.Fatalf("order_id = %q", query.Get("order_id"))
placeholder
	if query.Get("out_trade_no") != "sub2_42" {
		t.Fatalf("out_trade_no = %q", query.Get("out_trade_no"))
placeholder
	if query.Get("resume_token") != "" {
		t.Fatalf("resume_token = %q, want empty", query.Get("resume_token"))
placeholder
placeholder

func TestBuildPaymentReturnURLEmptyBase(t *testing.T) {
	t.Parallel()

	got, err := buildPaymentReturnURL("", 42, "sub2_42", "resume-token")
	if err != nil {
		t.Fatalf("buildPaymentReturnURL returned error: %v", err)
placeholder
	if got != "" {
		t.Fatalf("buildPaymentReturnURL = %q, want empty string", got)
placeholder
placeholder

func TestPaymentResumeTokenRoundTrip(t *testing.T) {
	t.Parallel()

	svc := NewPaymentResumeService([]byte("placeholder"))
	token, err := svc.CreateToken(ResumeTokenClaims{
		OrderID:            42,
		UserID:             7,
		ProviderInstanceID: "19",
		ProviderKey:        "easypay",
		PaymentType:        "wxpay",
		CanonicalReturnURL: "https://example.com/payment/result",
		IssuedAt:           1234567890,
placeholder)
	if err != nil {
		t.Fatalf("CreateToken returned error: %v", err)
placeholder

	claims, err := svc.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken returned error: %v", err)
placeholder
	if claims.OrderID != 42 || claims.UserID != 7 {
		t.Fatalf("claims mismatch: %+v", claims)
placeholder
	if claims.ProviderInstanceID != "19" || claims.ProviderKey != "easypay" || claims.PaymentType != "wxpay" {
		t.Fatalf("claims provider snapshot mismatch: %+v", claims)
placeholder
	if claims.CanonicalReturnURL != "https://example.com/payment/result" {
		t.Fatalf("claims return URL = %q", claims.CanonicalReturnURL)
placeholder
placeholder

func TestCreateTokenRejectsMissingSigningKey(t *testing.T) {
	t.Parallel()

	svc := NewPaymentResumeService(nil)
	_, err := svc.CreateToken(ResumeTokenClaims{OrderID: 42placeholder)
	if err == nil {
		t.Fatal("CreateToken should reject missing signing key")
placeholder
placeholder

func TestParseTokenRejectsFallbackSignedTokenWhenSigningKeyMissing(t *testing.T) {
	t.Parallel()

	token := mustCreateFallbackSignedToken(t, ResumeTokenClaims{OrderID: 42, UserID: 7placeholder)
	svc := NewPaymentResumeService(nil)
	_, err := svc.ParseToken(token)
	if err == nil {
		t.Fatal("ParseToken should reject tokens when signing key is missing")
placeholder
placeholder

func TestParseTokenRejectsExpiredToken(t *testing.T) {
	t.Parallel()

	svc := NewPaymentResumeService([]byte("placeholder"))
	token, err := svc.CreateToken(ResumeTokenClaims{
		OrderID:   42,
		UserID:    7,
		IssuedAt:  time.Now().Add(-25 * time.Hour).Unix(),
		ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(),
placeholder)
	if err != nil {
		t.Fatalf("CreateToken returned error: %v", err)
placeholder

	_, err = svc.ParseToken(token)
	if err == nil {
		t.Fatal("ParseToken should reject expired tokens")
placeholder
placeholder

func TestWeChatPaymentResumeTokenRoundTrip(t *testing.T) {
	t.Parallel()

	svc := NewPaymentResumeService([]byte("placeholder"))
	token, err := svc.CreateWeChatPaymentResumeToken(WeChatPaymentResumeClaims{
		OpenID:      "openid-123",
		PaymentType: payment.TypeWxpay,
		Amount:      "12.50",
		OrderType:   payment.OrderTypeSubscription,
		PlanID:      7,
		RedirectTo:  "/purchase?from=wechat",
		Scope:       "snsapi_base",
		IssuedAt:    1234567890,
placeholder)
	if err != nil {
		t.Fatalf("CreateWeChatPaymentResumeToken returned error: %v", err)
placeholder

	claims, err := svc.ParseWeChatPaymentResumeToken(token)
	if err != nil {
		t.Fatalf("ParseWeChatPaymentResumeToken returned error: %v", err)
placeholder
	if claims.OpenID != "openid-123" || claims.PaymentType != payment.TypeWxpay {
		t.Fatalf("claims mismatch: %+v", claims)
placeholder
	if claims.Amount != "12.50" || claims.OrderType != payment.OrderTypeSubscription || claims.PlanID != 7 {
		t.Fatalf("claims payment context mismatch: %+v", claims)
placeholder
	if claims.RedirectTo != "/purchase?from=wechat" || claims.Scope != "snsapi_base" {
		t.Fatalf("claims redirect/scope mismatch: %+v", claims)
placeholder
placeholder

func TestCreateWeChatPaymentResumeTokenRejectsMissingSigningKey(t *testing.T) {
	t.Parallel()

	svc := NewPaymentResumeService(nil)
	_, err := svc.CreateWeChatPaymentResumeToken(WeChatPaymentResumeClaims{OpenID: "openid-123"placeholder)
	if err == nil {
		t.Fatal("CreateWeChatPaymentResumeToken should reject missing signing key")
placeholder
placeholder

func TestParseWeChatPaymentResumeTokenRejectsFallbackSignedTokenWhenSigningKeyMissing(t *testing.T) {
	t.Parallel()

	token := mustCreateFallbackSignedToken(t, WeChatPaymentResumeClaims{
		TokenType:   wechatPaymentResumeTokenType,
		OpenID:      "openid-123",
		PaymentType: payment.TypeWxpay,
placeholder)
	svc := NewPaymentResumeService(nil)
	_, err := svc.ParseWeChatPaymentResumeToken(token)
	if err == nil {
		t.Fatal("ParseWeChatPaymentResumeToken should reject tokens when signing key is missing")
placeholder
placeholder

func TestParseWeChatPaymentResumeTokenRejectsExpiredToken(t *testing.T) {
	t.Parallel()

	svc := NewPaymentResumeService([]byte("placeholder"))
	token, err := svc.CreateWeChatPaymentResumeToken(WeChatPaymentResumeClaims{
		OpenID:      "openid-123",
		PaymentType: payment.TypeWxpay,
		IssuedAt:    time.Now().Add(-30 * time.Minute).Unix(),
		ExpiresAt:   time.Now().Add(-1 * time.Minute).Unix(),
placeholder)
	if err != nil {
		t.Fatalf("CreateWeChatPaymentResumeToken returned error: %v", err)
placeholder

	_, err = svc.ParseWeChatPaymentResumeToken(token)
	if err == nil {
		t.Fatal("ParseWeChatPaymentResumeToken should reject expired tokens")
placeholder
placeholder

func TestPaymentServiceParseWeChatPaymentResumeTokenUsesExplicitSigningKey(t *testing.T) {
	t.Setenv("PAYMENT_RESUME_SIGNING_KEY", "explicit-payment-resume-signing-key")

	token, err := NewPaymentResumeService([]byte("explicit-payment-resume-signing-key")).CreateWeChatPaymentResumeToken(WeChatPaymentResumeClaims{
		OpenID:      "openid-explicit-key",
		PaymentType: payment.TypeWxpay,
placeholder)
	if err != nil {
		t.Fatalf("CreateWeChatPaymentResumeToken returned error: %v", err)
placeholder

	svc := &PaymentService{
		configService: &PaymentConfigService{
			encryptionKey: []byte("placeholder"),
	placeholder,
placeholder

	claims, err := svc.ParseWeChatPaymentResumeToken(token)
	if err != nil {
		t.Fatalf("ParseWeChatPaymentResumeToken returned error: %v", err)
placeholder
	if claims.OpenID != "openid-explicit-key" {
		t.Fatalf("openid = %q, want %q", claims.OpenID, "openid-explicit-key")
placeholder
placeholder

func TestPaymentServiceParseWeChatPaymentResumeTokenAcceptsLegacyEncryptionKeyDuringMigration(t *testing.T) {
	t.Setenv("PAYMENT_RESUME_SIGNING_KEY", "explicit-payment-resume-signing-key")

	legacyKey := []byte("placeholder")
	token, err := NewPaymentResumeService(legacyKey).CreateWeChatPaymentResumeToken(WeChatPaymentResumeClaims{
		OpenID:      "openid-legacy-key",
		PaymentType: payment.TypeWxpay,
placeholder)
	if err != nil {
		t.Fatalf("CreateWeChatPaymentResumeToken returned error: %v", err)
placeholder

	svc := &PaymentService{
		configService: &PaymentConfigService{
			encryptionKey: legacyKey,
	placeholder,
placeholder

	claims, err := svc.ParseWeChatPaymentResumeToken(token)
	if err != nil {
		t.Fatalf("ParseWeChatPaymentResumeToken returned error: %v", err)
placeholder
	if claims.OpenID != "openid-legacy-key" {
		t.Fatalf("openid = %q, want %q", claims.OpenID, "openid-legacy-key")
placeholder
placeholder

func TestNormalizeVisibleMethodSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		method string
		input  string
		want   string
placeholder{
		{name: "alipay official alias", method: payment.TypeAlipay, input: "alipay", want: VisibleMethodSourceOfficialAlipayplaceholder,
		{name: "alipay easypay alias", method: payment.TypeAlipay, input: "easypay", want: VisibleMethodSourceEasyPayAlipayplaceholder,
		{name: "wxpay official alias", method: payment.TypeWxpay, input: "wxpay", want: VisibleMethodSourceOfficialWechatplaceholder,
		{name: "wxpay easypay alias", method: payment.TypeWxpay, input: "easypay", want: VisibleMethodSourceEasyPayWechatplaceholder,
		{name: "unsupported source", method: payment.TypeWxpay, input: "stripe", want: ""placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NormalizeVisibleMethodSource(tt.method, tt.input); got != tt.want {
				t.Fatalf("NormalizeVisibleMethodSource(%q, %q) = %q, want %q", tt.method, tt.input, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestVisibleMethodProviderKeyForSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		method string
		source string
		want   string
		ok     bool
placeholder{
		{name: "official alipay", method: payment.TypeAlipay, source: VisibleMethodSourceOfficialAlipay, want: payment.TypeAlipay, ok: trueplaceholder,
		{name: "easypay alipay", method: payment.TypeAlipay, source: VisibleMethodSourceEasyPayAlipay, want: payment.TypeEasyPay, ok: trueplaceholder,
		{name: "official wechat", method: payment.TypeWxpay, source: VisibleMethodSourceOfficialWechat, want: payment.TypeWxpay, ok: trueplaceholder,
		{name: "easypay wechat", method: payment.TypeWxpay, source: VisibleMethodSourceEasyPayWechat, want: payment.TypeEasyPay, ok: trueplaceholder,
		{name: "mismatched method and source", method: payment.TypeAlipay, source: VisibleMethodSourceOfficialWechat, want: "", ok: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, ok := VisibleMethodProviderKeyForSource(tt.method, tt.source)
			if got != tt.want || ok != tt.ok {
				t.Fatalf("VisibleMethodProviderKeyForSource(%q, %q) = (%q, %v), want (%q, %v)", tt.method, tt.source, got, ok, tt.want, tt.ok)
		placeholder
	placeholder)
placeholder
placeholder

func TestVisibleMethodLoadBalancerUsesEnabledProviderInstance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	_, err := client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeAlipay).
		SetName("Official Alipay").
		SetConfig("{placeholder").
		SetSupportedTypes("alipay").
		SetEnabled(true).
		SetSortOrder(1).
		Save(ctx)
	if err != nil {
		t.Fatalf("create alipay provider: %v", err)
placeholder

	inner := &captureLoadBalancer{placeholder
	configService := &PaymentConfigService{
		entClient: client,
placeholder
	lb := newVisibleMethodLoadBalancer(inner, configService)

	_, err = lb.SelectInstance(ctx, "", payment.TypeAlipay, payment.StrategyRoundRobin, 12.5)
	if err != nil {
		t.Fatalf("SelectInstance returned error: %v", err)
placeholder
	if inner.lastProviderKey != payment.TypeAlipay {
		t.Fatalf("lastProviderKey = %q, want %q", inner.lastProviderKey, payment.TypeAlipay)
placeholder
placeholder

func TestVisibleMethodLoadBalancerUsesConfiguredSourceWhenMultipleProvidersEnabled(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		method        payment.PaymentType
		officialName  string
		officialTypes string
		easyPayName   string
		easyPayTypes  string
		sourceSetting string
		wantProvider  string
placeholder{
		{
			name:          "alipay uses official source",
			method:        payment.TypeAlipay,
			officialName:  "Official Alipay",
			officialTypes: "alipay",
			easyPayName:   "EasyPay Alipay",
			easyPayTypes:  "alipay",
			sourceSetting: VisibleMethodSourceOfficialAlipay,
			wantProvider:  payment.TypeAlipay,
	placeholder,
		{
			name:          "alipay uses easypay source",
			method:        payment.TypeAlipay,
			officialName:  "Official Alipay",
			officialTypes: "alipay",
			easyPayName:   "EasyPay Alipay",
			easyPayTypes:  "alipay",
			sourceSetting: VisibleMethodSourceEasyPayAlipay,
			wantProvider:  payment.TypeEasyPay,
	placeholder,
		{
			name:          "wxpay uses official source",
			method:        payment.TypeWxpay,
			officialName:  "Official WeChat",
			officialTypes: "wxpay",
			easyPayName:   "EasyPay WeChat",
			easyPayTypes:  "wxpay",
			sourceSetting: VisibleMethodSourceOfficialWechat,
			wantProvider:  payment.TypeWxpay,
	placeholder,
		{
			name:          "wxpay uses easypay source",
			method:        payment.TypeWxpay,
			officialName:  "Official WeChat",
			officialTypes: "wxpay",
			easyPayName:   "EasyPay WeChat",
			easyPayTypes:  "wxpay",
			sourceSetting: VisibleMethodSourceEasyPayWechat,
			wantProvider:  payment.TypeEasyPay,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			client := newPaymentConfigServiceTestClient(t)

			officialProviderKey := payment.TypeAlipay
			if tt.method == payment.TypeWxpay {
				officialProviderKey = payment.TypeWxpay
		placeholder

			_, err := client.PaymentProviderInstance.Create().
				SetProviderKey(officialProviderKey).
				SetName(tt.officialName).
				SetConfig("{placeholder").
				SetSupportedTypes(tt.officialTypes).
				SetEnabled(true).
				SetSortOrder(1).
				Save(ctx)
			if err != nil {
				t.Fatalf("create official provider: %v", err)
		placeholder

			_, err = client.PaymentProviderInstance.Create().
				SetProviderKey(payment.TypeEasyPay).
				SetName(tt.easyPayName).
				SetConfig("{placeholder").
				SetSupportedTypes(tt.easyPayTypes).
				SetEnabled(true).
				SetSortOrder(2).
				Save(ctx)
			if err != nil {
				t.Fatalf("create easypay provider: %v", err)
		placeholder

			inner := &captureLoadBalancer{placeholder
			configService := &PaymentConfigService{
				entClient: client,
				settingRepo: &paymentConfigSettingRepoStub{
					values: map[string]string{
						visibleMethodSourceSettingKey(tt.method): tt.sourceSetting,
				placeholder,
			placeholder,
		placeholder
			lb := newVisibleMethodLoadBalancer(inner, configService)

			_, err = lb.SelectInstance(ctx, "", tt.method, payment.StrategyRoundRobin, 12.5)
			if err != nil {
				t.Fatalf("SelectInstance returned error: %v", err)
		placeholder
			if inner.lastProviderKey != tt.wantProvider {
				t.Fatalf("lastProviderKey = %q, want %q", inner.lastProviderKey, tt.wantProvider)
		placeholder
	placeholder)
placeholder
placeholder

func TestVisibleMethodLoadBalancerPreservesLegacyCrossProviderRoutingWhenSourceMissing(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)

	_, err := client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeAlipay).
		SetName("Official Alipay").
		SetConfig("{placeholder").
		SetSupportedTypes("alipay").
		SetEnabled(true).
		SetSortOrder(1).
		Save(ctx)
	if err != nil {
		t.Fatalf("create official provider: %v", err)
placeholder

	_, err = client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeEasyPay).
		SetName("EasyPay Alipay").
		SetConfig("{placeholder").
		SetSupportedTypes("alipay").
		SetEnabled(true).
		SetSortOrder(2).
		Save(ctx)
	if err != nil {
		t.Fatalf("create easypay provider: %v", err)
placeholder

	inner := &captureLoadBalancer{placeholder
	configService := &PaymentConfigService{
		entClient: client,
		settingRepo: &paymentConfigSettingRepoStub{
			values: map[string]string{
				visibleMethodSourceSettingKey(payment.TypeAlipay): "",
		placeholder,
	placeholder,
placeholder
	lb := newVisibleMethodLoadBalancer(inner, configService)

	_, err = lb.SelectInstance(ctx, "", payment.TypeAlipay, payment.StrategyRoundRobin, 9.9)
	if err != nil {
		t.Fatalf("SelectInstance returned error: %v", err)
placeholder
	if inner.lastProviderKey != "" {
		t.Fatalf("lastProviderKey = %q, want legacy cross-provider empty key", inner.lastProviderKey)
placeholder
	if inner.lastPaymentType != payment.TypeAlipay {
		t.Fatalf("lastPaymentType = %q, want %q", inner.lastPaymentType, payment.TypeAlipay)
placeholder
placeholder

func TestVisibleMethodLoadBalancerRejectsInvalidSourceWhenMultipleProvidersEnabled(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		method      payment.PaymentType
		sourceValue string
		wantMessage string
placeholder{
		{
			name:        "invalid wxpay source",
			method:      payment.TypeWxpay,
			sourceValue: "stripe",
			wantMessage: "wxpay source must be one of the supported payment providers",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			client := newPaymentConfigServiceTestClient(t)

			officialProviderKey := payment.TypeAlipay
			officialSupportedTypes := "alipay"
			officialName := "Official Alipay"
			easyPaySupportedTypes := "alipay"
			easyPayName := "EasyPay Alipay"
			if tt.method == payment.TypeWxpay {
				officialProviderKey = payment.TypeWxpay
				officialSupportedTypes = "wxpay"
				officialName = "Official WeChat"
				easyPaySupportedTypes = "wxpay"
				easyPayName = "EasyPay WeChat"
		placeholder

			_, err := client.PaymentProviderInstance.Create().
				SetProviderKey(officialProviderKey).
				SetName(officialName).
				SetConfig("{placeholder").
				SetSupportedTypes(officialSupportedTypes).
				SetEnabled(true).
				SetSortOrder(1).
				Save(ctx)
			if err != nil {
				t.Fatalf("create official provider: %v", err)
		placeholder

			_, err = client.PaymentProviderInstance.Create().
				SetProviderKey(payment.TypeEasyPay).
				SetName(easyPayName).
				SetConfig("{placeholder").
				SetSupportedTypes(easyPaySupportedTypes).
				SetEnabled(true).
				SetSortOrder(2).
				Save(ctx)
			if err != nil {
				t.Fatalf("create easypay provider: %v", err)
		placeholder

			inner := &captureLoadBalancer{placeholder
			configService := &PaymentConfigService{
				entClient: client,
				settingRepo: &paymentConfigSettingRepoStub{
					values: map[string]string{
						visibleMethodSourceSettingKey(tt.method): tt.sourceValue,
				placeholder,
			placeholder,
		placeholder
			lb := newVisibleMethodLoadBalancer(inner, configService)

			_, err = lb.SelectInstance(ctx, "", tt.method, payment.StrategyRoundRobin, 9.9)
			if err == nil {
				t.Fatal("SelectInstance should reject invalid visible method source configuration")
		placeholder
			if infraerrors.Reason(err) != "INVALID_PAYMENT_VISIBLE_METHOD_SOURCE" {
				t.Fatalf("Reason(err) = %q, want %q", infraerrors.Reason(err), "INVALID_PAYMENT_VISIBLE_METHOD_SOURCE")
		placeholder
			if infraerrors.Message(err) != tt.wantMessage {
				t.Fatalf("Message(err) = %q, want %q", infraerrors.Message(err), tt.wantMessage)
		placeholder
	placeholder)
placeholder
placeholder

func TestVisibleMethodLoadBalancerRejectsMissingEnabledVisibleMethodProvider(t *testing.T) {
	t.Parallel()

	inner := &captureLoadBalancer{placeholder
	configService := &PaymentConfigService{
		entClient: newPaymentConfigServiceTestClient(t),
placeholder
	lb := newVisibleMethodLoadBalancer(inner, configService)

	if _, err := lb.SelectInstance(context.Background(), "", payment.TypeWxpay, payment.StrategyRoundRobin, 9.9); err == nil {
		t.Fatal("SelectInstance should reject when no enabled provider instance exists")
placeholder
placeholder

type captureLoadBalancer struct {
	lastProviderKey string
	lastPaymentType string
placeholder

func (c *captureLoadBalancer) GetInstanceConfig(context.Context, int64) (map[string]string, error) {
	return map[string]string{placeholder, nil
placeholder

func (c *captureLoadBalancer) SelectInstance(_ context.Context, providerKey string, paymentType payment.PaymentType, _ payment.Strategy, _ float64) (*payment.InstanceSelection, error) {
	c.lastProviderKey = providerKey
	c.lastPaymentType = paymentType
	return &payment.InstanceSelection{ProviderKey: providerKey, SupportedTypes: paymentTypeplaceholder, nil
placeholder

func mustCreateFallbackSignedToken(t *testing.T, claims any) string {
placeholder

	payload, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("marshal claims: %v", err)
placeholder
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, []byte("sub2api-payment-resume"))
	_, _ = mac.Write([]byte(encodedPayload))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return encodedPayload + "." + signature
placeholder
