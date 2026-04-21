//go:build unit

package provider

import (
	"context"
	"net/url"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/payment"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
)

func TestMapWxState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
placeholder{
		{
			name:  "SUCCESS maps to paid",
			input: wxpayTradeStateSuccess,
			want:  payment.ProviderStatusPaid,
	placeholder,
		{
			name:  "REFUND maps to refunded",
			input: wxpayTradeStateRefund,
			want:  payment.ProviderStatusRefunded,
	placeholder,
		{
			name:  "CLOSED maps to failed",
			input: wxpayTradeStateClosed,
			want:  payment.ProviderStatusFailed,
	placeholder,
		{
			name:  "PAYERROR maps to failed",
			input: wxpayTradeStatePayError,
			want:  payment.ProviderStatusFailed,
	placeholder,
		{
			name:  "unknown state maps to pending",
			input: "NOTPAY",
			want:  payment.ProviderStatusPending,
	placeholder,
		{
			name:  "empty string maps to pending",
			input: "",
			want:  payment.ProviderStatusPending,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := mapWxState(tt.input)
			if got != tt.want {
				t.Errorf("mapWxState(%q) = %q, want %q", tt.input, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestWxSV(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *string
		want  string
placeholder{
		{
			name:  "nil pointer returns empty string",
			input: nil,
			want:  "",
	placeholder,
		{
			name:  "non-nil pointer returns value",
			input: strPtr("hello"),
			want:  "hello",
	placeholder,
		{
			name:  "pointer to empty string returns empty string",
			input: strPtr(""),
			want:  "",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := wxSV(tt.input)
			if got != tt.want {
				t.Errorf("wxSV() = %q, want %q", got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestBuildWxpayTransactionMetadata(t *testing.T) {
	t.Parallel()

	tx := &payments.Transaction{
		Appid:      strPtr("wx-app-id"),
		Mchid:      strPtr("mch-id"),
		TradeState: strPtr(wxpayTradeStateSuccess),
		Amount: &payments.Amount{
			Currency: strPtr(wxpayCurrency),
	placeholder,
placeholder

	metadata := buildWxpayTransactionMetadata(tx)
	if metadata[wxpayMetadataAppID] != "wx-app-id" {
		t.Fatalf("appid = %q", metadata[wxpayMetadataAppID])
placeholder
	if metadata[wxpayMetadataMerchantID] != "mch-id" {
		t.Fatalf("mchid = %q", metadata[wxpayMetadataMerchantID])
placeholder
	if metadata[wxpayMetadataCurrency] != wxpayCurrency {
		t.Fatalf("currency = %q", metadata[wxpayMetadataCurrency])
placeholder
	if metadata[wxpayMetadataTradeState] != wxpayTradeStateSuccess {
		t.Fatalf("trade_state = %q", metadata[wxpayMetadataTradeState])
placeholder
placeholder

func strPtr(s string) *string {
	return &s
placeholder

func TestFormatPEM(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     string
		keyType string
		want    string
placeholder{
		{
			name:    "raw key gets wrapped with headers",
			key:     "MIIBIjANBgkqhki...",
			keyType: "PUBLIC KEY",
			want:    "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhki...\n-----END PUBLIC KEY-----",
	placeholder,
		{
			name:    "already formatted key is returned as-is",
			key:     "placeholder
placeholder
placeholder\nMIIEvQIBADANBg...\n-----END PRIVATE KEY-----",
	placeholder,
		{
			name:    "key with leading/trailing whitespace is trimmed before check",
			key:     "  \n MIIBIjANBgkqhki...  \n ",
			keyType: "PUBLIC KEY",
			want:    "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhki...\n-----END PUBLIC KEY-----",
	placeholder,
		{
			name:    "already formatted key with whitespace is trimmed and returned",
			key:     "  placeholder
placeholder
placeholder\ndata\n-----END RSA PRIVATE KEY-----",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := formatPEM(tt.key, tt.keyType)
			if got != tt.want {
				t.Errorf("formatPEM(%q, %q) =\n%s\nwant:\n%s", tt.key, tt.keyType, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestNewWxpay(t *testing.T) {
	t.Parallel()

	validConfig := map[string]string{
		"appId":       "wx1234567890",
		"mchId":       "1234567890",
		"privateKey":  "fake-private-key",
		"apiV3Key":    "12345678901234567890123456789012", // exactly 32 bytes
		"publicKey":   "fake-public-key",
		"publicKeyId": "key-id-001",
		"certSerial":  "SERIAL001",
placeholder

	// helper to clone and override config fields
	withOverride := func(overrides map[string]string) map[string]string {
		cfg := make(map[string]string, len(validConfig))
		for k, v := range validConfig {
			cfg[k] = v
	placeholder
		for k, v := range overrides {
			cfg[k] = v
	placeholder
		return cfg
placeholder

	tests := []struct {
		name      string
		config    map[string]string
		wantErr   bool
		errSubstr string
placeholder{
		{
			name:    "valid config succeeds",
			config:  validConfig,
			wantErr: false,
	placeholder,
		{
			name:      "missing appId",
			config:    withOverride(map[string]string{"appId": ""placeholder),
			wantErr:   true,
			errSubstr: "appId",
	placeholder,
		{
			name:      "missing mchId",
			config:    withOverride(map[string]string{"mchId": ""placeholder),
			wantErr:   true,
			errSubstr: "mchId",
	placeholder,
		{
			name:      "missing privateKey",
			config:    withOverride(map[string]string{"privateKey": ""placeholder),
			wantErr:   true,
			errSubstr: "privateKey",
	placeholder,
		{
			name:      "missing apiV3Key",
			config:    withOverride(map[string]string{"apiV3Key": ""placeholder),
			wantErr:   true,
			errSubstr: "apiV3Key",
	placeholder,
		{
			name:      "missing publicKey",
			config:    withOverride(map[string]string{"publicKey": ""placeholder),
			wantErr:   true,
			errSubstr: "publicKey",
	placeholder,
		{
			name:      "missing publicKeyId",
			config:    withOverride(map[string]string{"publicKeyId": ""placeholder),
			wantErr:   true,
			errSubstr: "publicKeyId",
	placeholder,
		{
			name:      "apiV3Key too short",
			config:    withOverride(map[string]string{"apiV3Key": "short"placeholder),
			wantErr:   true,
			errSubstr: "exactly 32 bytes",
	placeholder,
		{
			name:      "apiV3Key too long",
			config:    withOverride(map[string]string{"apiV3Key": "123456789012345678901234567890123"placeholder), // 33 bytes
			wantErr:   true,
			errSubstr: "exactly 32 bytes",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewWxpay("test-instance", tt.config)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
			placeholder
				if tt.errSubstr != "" && !strings.Contains(err.Error(), tt.errSubstr) {
					t.Errorf("error %q should contain %q", err.Error(), tt.errSubstr)
			placeholder
				return
		placeholder
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
		placeholder
			if got == nil {
				t.Fatal("expected non-nil Wxpay instance")
		placeholder
			if got.instanceID != "test-instance" {
				t.Errorf("instanceID = %q, want %q", got.instanceID, "test-instance")
		placeholder
	placeholder)
placeholder
placeholder

func TestBuildWxpayResultURLPreservesResumeToken(t *testing.T) {
	t.Parallel()

	resultURL, err := buildWxpayResultURL("https://app.example.com/payment/result?order_id=42&resume_token=resume-42&status=success", payment.CreatePaymentRequest{
		OrderID:     "sub2_42",
		PaymentType: payment.TypeWxpay,
placeholder)
	if err != nil {
		t.Fatalf("buildWxpayResultURL returned error: %v", err)
placeholder

	parsed, err := url.Parse(resultURL)
	if err != nil {
		t.Fatalf("url.Parse returned error: %v", err)
placeholder
	query := parsed.Query()
	if parsed.Path != wxpayResultPath {
		t.Fatalf("path = %q, want %q", parsed.Path, wxpayResultPath)
placeholder
	if query.Get("resume_token") != "resume-42" {
		t.Fatalf("resume_token = %q, want %q", query.Get("resume_token"), "resume-42")
placeholder
	if query.Get("order_id") != "42" {
		t.Fatalf("order_id = %q, want %q", query.Get("order_id"), "42")
placeholder
	if query.Get("out_trade_no") != "sub2_42" {
		t.Fatalf("out_trade_no = %q, want %q", query.Get("out_trade_no"), "sub2_42")
placeholder
placeholder

func TestResolveWxpayJSAPIAppID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		config map[string]string
		want   string
placeholder{
		{
			name: "prefers dedicated mp app id",
			config: map[string]string{
				"mpAppId": "wx-mp-app",
				"appId":   "wx-merchant-app",
		placeholder,
			want: "wx-mp-app",
	placeholder,
		{
			name: "falls back to merchant app id",
			config: map[string]string{
				"appId": "wx-merchant-app",
		placeholder,
			want: "wx-merchant-app",
	placeholder,
		{
			name:   "missing app ids returns empty",
			config: map[string]string{placeholder,
			want:   "",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ResolveWxpayJSAPIAppID(tt.config); got != tt.want {
				t.Fatalf("ResolveWxpayJSAPIAppID() = %q, want %q", got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestResolveWxpayCreateMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		req      payment.CreatePaymentRequest
		wantMode string
		wantErr  string
placeholder{
		{
			name:     "desktop uses native",
			req:      payment.CreatePaymentRequest{placeholder,
			wantMode: wxpayModeNative,
	placeholder,
		{
			name: "mobile uses h5 when client ip is present",
			req: payment.CreatePaymentRequest{
				IsMobile: true,
				ClientIP: "203.0.113.10",
		placeholder,
			wantMode: wxpayModeH5,
	placeholder,
		{
			name: "mobile without client ip returns clear error",
			req: payment.CreatePaymentRequest{
				IsMobile: true,
		placeholder,
			wantErr: "requires client IP",
	placeholder,
		{
			name: "openid uses jsapi mode",
			req: payment.CreatePaymentRequest{
				OpenID: "openid-123",
		placeholder,
			wantMode: wxpayModeJSAPI,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := resolveWxpayCreateMode(tt.req)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
			placeholder
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error %q should contain %q", err.Error(), tt.wantErr)
			placeholder
				return
		placeholder
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
		placeholder
			if got != tt.wantMode {
				t.Fatalf("resolveWxpayCreateMode() = %q, want %q", got, tt.wantMode)
		placeholder
	placeholder)
placeholder
placeholder

func TestCreatePaymentWithOpenIDReturnsJSAPIResult(t *testing.T) {
	origJSAPIPrepay := wxpayJSAPIPrepayWithRequestPayment
	origNativePrepay := wxpayNativePrepay
	origH5Prepay := wxpayH5Prepay
	t.Cleanup(func() {
		wxpayJSAPIPrepayWithRequestPayment = origJSAPIPrepay
		wxpayNativePrepay = origNativePrepay
		wxpayH5Prepay = origH5Prepay
placeholder)

	jsapiCalls := 0
	nativeCalls := 0
	h5Calls := 0
	wxpayJSAPIPrepayWithRequestPayment = func(ctx context.Context, svc jsapi.JsapiApiService, req jsapi.PrepayRequest) (*jsapi.PrepayWithRequestPaymentResponse, *core.APIResult, error) {
		jsapiCalls++
		if got := wxSV(req.Payer.Openid); got != "openid-123" {
			t.Fatalf("openid = %q, want %q", got, "openid-123")
	placeholder
		if req.SceneInfo == nil || wxSV(req.SceneInfo.PayerClientIp) != "203.0.113.10" {
			t.Fatalf("scene_info payer_client_ip = %q, want %q", wxSV(req.SceneInfo.PayerClientIp), "203.0.113.10")
	placeholder
		return &jsapi.PrepayWithRequestPaymentResponse{
			Appid:     core.String("wx123"),
			TimeStamp: core.String("1712345678"),
			NonceStr:  core.String("nonce-123"),
			Package:   core.String("placeholder"),
			SignType:  core.String("RSA"),
			PaySign:   core.String("signed-payload"),
	placeholder, nil, nil
placeholder
	wxpayNativePrepay = func(ctx context.Context, svc native.NativeApiService, req native.PrepayRequest) (*native.PrepayResponse, *core.APIResult, error) {
		nativeCalls++
		return &native.PrepayResponse{placeholder, nil, nil
placeholder
	wxpayH5Prepay = func(ctx context.Context, svc h5.H5ApiService, req h5.PrepayRequest) (*h5.PrepayResponse, *core.APIResult, error) {
		h5Calls++
		return &h5.PrepayResponse{placeholder, nil, nil
placeholder

	provider := &Wxpay{
		config: map[string]string{
			"appId": "wx123",
			"mchId": "mch123",
	placeholder,
		coreClient: &core.Client{placeholder,
placeholder

	resp, err := provider.CreatePayment(context.Background(), payment.CreatePaymentRequest{
		OrderID:     "sub2_88",
		Amount:      "66.88",
		PaymentType: payment.TypeWxpay,
		NotifyURL:   "https://merchant.example/payment/notify",
		OpenID:      "openid-123",
		ClientIP:    "203.0.113.10",
placeholder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
placeholder
	if jsapiCalls != 1 {
		t.Fatalf("jsapi prepay calls = %d, want 1", jsapiCalls)
placeholder
	if nativeCalls != 0 {
		t.Fatalf("native prepay calls = %d, want 0", nativeCalls)
placeholder
	if h5Calls != 0 {
		t.Fatalf("h5 prepay calls = %d, want 0", h5Calls)
placeholder
	if resp.ResultType != payment.CreatePaymentResultJSAPIReady {
		t.Fatalf("result type = %q, want %q", resp.ResultType, payment.CreatePaymentResultJSAPIReady)
placeholder
	if resp.JSAPI == nil {
		t.Fatal("expected jsapi payload, got nil")
placeholder
	if resp.JSAPI.AppID != "wx123" {
		t.Fatalf("jsapi appId = %q, want %q", resp.JSAPI.AppID, "wx123")
placeholder
	if resp.JSAPI.TimeStamp != "1712345678" {
		t.Fatalf("jsapi timeStamp = %q, want %q", resp.JSAPI.TimeStamp, "1712345678")
placeholder
	if resp.JSAPI.NonceStr != "nonce-123" {
		t.Fatalf("jsapi nonceStr = %q, want %q", resp.JSAPI.NonceStr, "nonce-123")
placeholder
	if resp.JSAPI.Package != "placeholder" {
		t.Fatalf("jsapi package = %q, want %q", resp.JSAPI.Package, "placeholder")
placeholder
	if resp.JSAPI.SignType != "RSA" {
		t.Fatalf("jsapi signType = %q, want %q", resp.JSAPI.SignType, "RSA")
placeholder
	if resp.JSAPI.PaySign != "signed-payload" {
		t.Fatalf("jsapi paySign = %q, want %q", resp.JSAPI.PaySign, "signed-payload")
placeholder
placeholder

func TestCreatePaymentMobileH5IncludesConfiguredSceneInfo(t *testing.T) {
	origJSAPIPrepay := wxpayJSAPIPrepayWithRequestPayment
	origNativePrepay := wxpayNativePrepay
	origH5Prepay := wxpayH5Prepay
	t.Cleanup(func() {
		wxpayJSAPIPrepayWithRequestPayment = origJSAPIPrepay
		wxpayNativePrepay = origNativePrepay
		wxpayH5Prepay = origH5Prepay
placeholder)

	jsapiCalls := 0
	nativeCalls := 0
	h5Calls := 0
	wxpayJSAPIPrepayWithRequestPayment = func(ctx context.Context, svc jsapi.JsapiApiService, req jsapi.PrepayRequest) (*jsapi.PrepayWithRequestPaymentResponse, *core.APIResult, error) {
		jsapiCalls++
		return &jsapi.PrepayWithRequestPaymentResponse{placeholder, nil, nil
placeholder
	wxpayNativePrepay = func(ctx context.Context, svc native.NativeApiService, req native.PrepayRequest) (*native.PrepayResponse, *core.APIResult, error) {
		nativeCalls++
		return &native.PrepayResponse{placeholder, nil, nil
placeholder
	wxpayH5Prepay = func(ctx context.Context, svc h5.H5ApiService, req h5.PrepayRequest) (*h5.PrepayResponse, *core.APIResult, error) {
		h5Calls++
		if req.SceneInfo == nil {
			t.Fatal("expected scene_info, got nil")
	placeholder
		if got := wxSV(req.SceneInfo.PayerClientIp); got != "203.0.113.10" {
			t.Fatalf("scene_info payer_client_ip = %q, want %q", got, "203.0.113.10")
	placeholder
		if req.SceneInfo.H5Info == nil {
			t.Fatal("expected scene_info.h5_info, got nil")
	placeholder
		if got := wxSV(req.SceneInfo.H5Info.Type); got != wxpayH5Type {
			t.Fatalf("scene_info.h5_info.type = %q, want %q", got, wxpayH5Type)
	placeholder
		if got := wxSV(req.SceneInfo.H5Info.AppName); got != "Sub2API" {
			t.Fatalf("scene_info.h5_info.app_name = %q, want %q", got, "Sub2API")
	placeholder
		if got := wxSV(req.SceneInfo.H5Info.AppUrl); got != "https://app.example.com" {
			t.Fatalf("scene_info.h5_info.app_url = %q, want %q", got, "https://app.example.com")
	placeholder
		return &h5.PrepayResponse{
			H5Url: core.String("https://wx.tenpay.example/h5pay?prepay_id=1"),
	placeholder, nil, nil
placeholder

	provider := &Wxpay{
		config: map[string]string{
			"appId":     "wx123",
			"mchId":     "mch123",
			"h5AppName": "Sub2API",
			"h5AppUrl":  "https://app.example.com",
	placeholder,
		coreClient: &core.Client{placeholder,
placeholder

	resp, err := provider.CreatePayment(context.Background(), payment.CreatePaymentRequest{
		OrderID:     "sub2_99",
		Amount:      "66.88",
		PaymentType: payment.TypeWxpay,
		Subject:     "Balance Recharge",
		NotifyURL:   "https://merchant.example/payment/notify",
		ReturnURL:   "https://merchant.example/payment/result?resume_token=resume-99",
		ClientIP:    "203.0.113.10",
		IsMobile:    true,
placeholder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
placeholder
	if jsapiCalls != 0 {
		t.Fatalf("jsapi prepay calls = %d, want 0", jsapiCalls)
placeholder
	if nativeCalls != 0 {
		t.Fatalf("native prepay calls = %d, want 0", nativeCalls)
placeholder
	if h5Calls != 1 {
		t.Fatalf("h5 prepay calls = %d, want 1", h5Calls)
placeholder
	if !strings.Contains(resp.PayURL, "redirect_url=") {
		t.Fatalf("pay_url = %q, want redirect_url query appended", resp.PayURL)
placeholder
placeholder
