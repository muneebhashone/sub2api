//go:build unit

package provider

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/payment"
	"github.com/smartwalle/alipay/v3"
)

func TestIsTradeNotExist(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
placeholder{
		{
			name: "nil error returns false",
			err:  nil,
			want: false,
	placeholder,
		{
			name: "error containing ACQ.TRADE_NOT_EXIST returns true",
			err:  errors.New("alipay: sub_code=ACQ.TRADE_NOT_EXIST, sub_msg=交易不存在"),
			want: true,
	placeholder,
		{
			name: "error not containing the code returns false",
			err:  errors.New("alipay: sub_code=ACQ.SYSTEM_ERROR, sub_msg=系统错误"),
			want: false,
	placeholder,
		{
			name: "error with only partial match returns false",
			err:  errors.New("ACQ.TRADE_NOT"),
			want: false,
	placeholder,
		{
			name: "error with exact constant value returns true",
			err:  errors.New(alipayErrTradeNotExist),
			want: true,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := isTradeNotExist(tt.err)
			if got != tt.want {
				t.Errorf("isTradeNotExist(%v) = %v, want %v", tt.err, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestNewAlipay(t *testing.T) {
	t.Parallel()

	validConfig := map[string]string{
		"appId":      "2021001234567890",
		"privateKey": "placeholder",
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
			name:      "missing privateKey",
			config:    withOverride(map[string]string{"privateKey": ""placeholder),
			wantErr:   true,
			errSubstr: "privateKey",
	placeholder,
		{
			name:      "nil config map returns error for appId",
			config:    map[string]string{placeholder,
			wantErr:   true,
			errSubstr: "appId",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewAlipay("test-instance", tt.config)
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
				t.Fatal("expected non-nil Alipay instance")
		placeholder
			if got.instanceID != "test-instance" {
				t.Errorf("instanceID = %q, want %q", got.instanceID, "test-instance")
		placeholder
	placeholder)
placeholder
placeholder

func TestCreateTradeUsesPreCreateForDesktop(t *testing.T) {
	origPreCreate := alipayTradePreCreate
	origPagePay := alipayTradePagePay
	origWapPay := alipayTradeWapPay
	t.Cleanup(func() {
		alipayTradePreCreate = origPreCreate
		alipayTradePagePay = origPagePay
		alipayTradeWapPay = origWapPay
placeholder)

	preCreateCalls := 0
	pagePayCalls := 0
	wapPayCalls := 0
	alipayTradePreCreate = func(ctx context.Context, client *alipay.Client, param alipay.TradePreCreate) (*alipay.TradePreCreateRsp, error) {
		preCreateCalls++
		if param.OutTradeNo != "sub2_100" {
			t.Fatalf("out_trade_no = %q, want %q", param.OutTradeNo, "sub2_100")
	placeholder
		if param.NotifyURL != "https://merchant.example.com/api/v1/payment/webhook/alipay" {
			t.Fatalf("notify_url = %q", param.NotifyURL)
	placeholder
		return &alipay.TradePreCreateRsp{
			OutTradeNo: "sub2_100",
			QRCode:     "https://qr.alipay.example.com/precreate-token",
	placeholder, nil
placeholder
	alipayTradePagePay = func(client *alipay.Client, param alipay.TradePagePay) (*url.URL, error) {
		pagePayCalls++
		return url.Parse("https://openapi.alipay.com/gateway.do?page-pay")
placeholder
	alipayTradeWapPay = func(client *alipay.Client, param alipay.TradeWapPay) (*url.URL, error) {
		wapPayCalls++
		return url.Parse("https://openapi.alipay.com/gateway.do?wap-pay")
placeholder

	provider := &Alipay{placeholder
	resp, err := provider.createTrade(context.Background(), &alipay.Client{placeholder, payment.CreatePaymentRequest{
		OrderID: "sub2_100",
		Amount:  "88.00",
		Subject: "Balance recharge",
placeholder, "https://merchant.example.com/api/v1/payment/webhook/alipay", "https://merchant.example.com/payment/result", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
placeholder
	if preCreateCalls != 1 {
		t.Fatalf("precreate calls = %d, want 1", preCreateCalls)
placeholder
	if pagePayCalls != 0 {
		t.Fatalf("page pay calls = %d, want 0", pagePayCalls)
placeholder
	if wapPayCalls != 0 {
		t.Fatalf("wap pay calls = %d, want 0", wapPayCalls)
placeholder
	if resp.QRCode != "https://qr.alipay.example.com/precreate-token" {
		t.Fatalf("qr_code = %q", resp.QRCode)
placeholder
	if resp.PayURL != "" {
		t.Fatalf("pay_url = %q, want empty", resp.PayURL)
placeholder
placeholder

func TestCreateTradeUsesWapPayForMobile(t *testing.T) {
	origPreCreate := alipayTradePreCreate
	origWapPay := alipayTradeWapPay
	t.Cleanup(func() {
		alipayTradePreCreate = origPreCreate
		alipayTradeWapPay = origWapPay
placeholder)

	preCreateCalls := 0
	alipayTradePreCreate = func(ctx context.Context, client *alipay.Client, param alipay.TradePreCreate) (*alipay.TradePreCreateRsp, error) {
		preCreateCalls++
		return &alipay.TradePreCreateRsp{placeholder, nil
placeholder

	wapPayCalls := 0
	alipayTradeWapPay = func(client *alipay.Client, param alipay.TradeWapPay) (*url.URL, error) {
		wapPayCalls++
		if param.ReturnURL != "https://merchant.example.com/payment/result" {
			t.Fatalf("return_url = %q", param.ReturnURL)
	placeholder
		return url.Parse("https://openapi.alipay.com/gateway.do?wap-pay")
placeholder

	provider := &Alipay{placeholder
	resp, err := provider.createTrade(context.Background(), &alipay.Client{placeholder, payment.CreatePaymentRequest{
		OrderID:  "sub2_101",
		Amount:   "18.00",
		Subject:  "Balance recharge",
		IsMobile: true,
placeholder, "https://merchant.example.com/api/v1/payment/webhook/alipay", "https://merchant.example.com/payment/result", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
placeholder
	if preCreateCalls != 0 {
		t.Fatalf("precreate calls = %d, want 0", preCreateCalls)
placeholder
	if wapPayCalls != 1 {
		t.Fatalf("wap pay calls = %d, want 1", wapPayCalls)
placeholder
	if resp.PayURL == "" {
		t.Fatal("expected pay_url for mobile wap pay")
placeholder
	if resp.QRCode != "" {
		t.Fatalf("qr_code = %q, want empty", resp.QRCode)
placeholder
placeholder

func TestAlipayMerchantIdentityMetadata(t *testing.T) {
	t.Parallel()

	provider := &Alipay{
		config: map[string]string{
			"appId": "2021001234567890",
	placeholder,
placeholder

	metadata := provider.MerchantIdentityMetadata()
	if metadata["app_id"] != "2021001234567890" {
		t.Fatalf("app_id = %q, want %q", metadata["app_id"], "2021001234567890")
placeholder
placeholder
