package service

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/payment"
)

func TestPcParseFloat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      string
		defaultVal float64
		expected   float64
placeholder{
		{"empty string returns default", "", 1.0, 1.0placeholder,
		{"valid float", "3.14", 0, 3.14placeholder,
		{"valid integer as float", "42", 0, 42.0placeholder,
		{"invalid string returns default", "notanumber", 9.99, 9.99placeholder,
		{"zero value", "0", 5.0, 0placeholder,
		{"negative value", "-10.5", 0, -placeholder,
		{"very large value", "99999999.99", 0, 99999999.99placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := pcParseFloat(tt.input, tt.defaultVal)
			if got != tt.expected {
				t.Fatalf("pcParseFloat(%q, %v) = %v, want %v", tt.input, tt.defaultVal, got, tt.expected)
		placeholder
	placeholder)
placeholder
placeholder

func TestPcParseInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      string
		defaultVal int
		expected   int
placeholder{
		{"empty string returns default", "", 30, 30placeholder,
		{"valid int", "10", 0, 10placeholder,
		{"invalid string returns default", "abc", 5, 5placeholder,
		{"float string returns default", "3.14", 0, 0placeholder,
		{"zero value", "0", 99, 0placeholder,
		{"negative value", "-1", 0, -1placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := pcParseInt(tt.input, tt.defaultVal)
			if got != tt.expected {
				t.Fatalf("pcParseInt(%q, %v) = %v, want %v", tt.input, tt.defaultVal, got, tt.expected)
		placeholder
	placeholder)
placeholder
placeholder

func TestParsePaymentConfig(t *testing.T) {
	t.Parallel()

	svc := &PaymentConfigService{placeholder

	t.Run("empty vals uses defaults", func(t *testing.T) {
		t.Parallel()
		cfg := svc.parsePaymentConfig(map[string]string{placeholder)
		if cfg.Enabled {
			t.Fatal("expected Enabled=false by default")
	placeholder
		if cfg.MinAmount != 1 {
			t.Fatalf("expected MinAmount=1, got %v", cfg.MinAmount)
	placeholder
		if cfg.MaxAmount != 0 {
			t.Fatalf("expected MaxAmount=0 (no limit), got %v", cfg.MaxAmount)
	placeholder
		if cfg.OrderTimeoutMin != 30 {
			t.Fatalf("expected OrderTimeoutMin=30, got %v", cfg.OrderTimeoutMin)
	placeholder
		if cfg.MaxPendingOrders != 3 {
			t.Fatalf("expected MaxPendingOrders=3, got %v", cfg.MaxPendingOrders)
	placeholder
		if cfg.LoadBalanceStrategy != payment.DefaultLoadBalanceStrategy {
			t.Fatalf("expected LoadBalanceStrategy=%s, got %q", payment.DefaultLoadBalanceStrategy, cfg.LoadBalanceStrategy)
	placeholder
		if len(cfg.EnabledTypes) != 0 {
			t.Fatalf("expected empty EnabledTypes, got %v", cfg.EnabledTypes)
	placeholder
placeholder)

	t.Run("all values populated", func(t *testing.T) {
		t.Parallel()
		vals := map[string]string{
			SettingPaymentEnabled:      "true",
			SettingMinRechargeAmount:   "5.00",
			SettingMaxRechargeAmount:   "1000.00",
			SettingDailyRechargeLimit:  "5000.00",
			SettingOrderTimeoutMinutes: "15",
			SettingMaxPendingOrders:    "5",
			SettingEnabledPaymentTypes: "alipay,wxpay,stripe",
			SettingBalancePayDisabled:  "true",
			SettingLoadBalanceStrategy: "least_amount",
			SettingProductNamePrefix:   "PRE",
			SettingProductNameSuffix:   "SUF",
	placeholder
		cfg := svc.parsePaymentConfig(vals)

		if !cfg.Enabled {
			t.Fatal("expected Enabled=true")
	placeholder
		if cfg.MinAmount != 5 {
			t.Fatalf("MinAmount = %v, want 5", cfg.MinAmount)
	placeholder
		if cfg.MaxAmount != 1000 {
			t.Fatalf("MaxAmount = %v, want 1000", cfg.MaxAmount)
	placeholder
		if cfg.DailyLimit != 5000 {
			t.Fatalf("DailyLimit = %v, want 5000", cfg.DailyLimit)
	placeholder
		if cfg.OrderTimeoutMin != 15 {
			t.Fatalf("OrderTimeoutMin = %v, want 15", cfg.OrderTimeoutMin)
	placeholder
		if cfg.MaxPendingOrders != 5 {
			t.Fatalf("MaxPendingOrders = %v, want 5", cfg.MaxPendingOrders)
	placeholder
		if len(cfg.EnabledTypes) != 3 {
			t.Fatalf("EnabledTypes len = %d, want 3", len(cfg.EnabledTypes))
	placeholder
		if cfg.EnabledTypes[0] != "alipay" || cfg.EnabledTypes[1] != "wxpay" || cfg.EnabledTypes[2] != "stripe" {
			t.Fatalf("EnabledTypes = %v, want [alipay wxpay stripe]", cfg.EnabledTypes)
	placeholder
		if !cfg.BalanceDisabled {
			t.Fatal("expected BalanceDisabled=true")
	placeholder
		if cfg.LoadBalanceStrategy != "least_amount" {
			t.Fatalf("LoadBalanceStrategy = %q, want %q", cfg.LoadBalanceStrategy, "least_amount")
	placeholder
		if cfg.ProductNamePrefix != "PRE" {
			t.Fatalf("ProductNamePrefix = %q, want %q", cfg.ProductNamePrefix, "PRE")
	placeholder
		if cfg.ProductNameSuffix != "SUF" {
			t.Fatalf("ProductNameSuffix = %q, want %q", cfg.ProductNameSuffix, "SUF")
	placeholder
placeholder)

	t.Run("enabled types with spaces are trimmed", func(t *testing.T) {
		t.Parallel()
		vals := map[string]string{
			SettingEnabledPaymentTypes: " alipay , wxpay ",
	placeholder
		cfg := svc.parsePaymentConfig(vals)
		if len(cfg.EnabledTypes) != 2 {
			t.Fatalf("EnabledTypes len = %d, want 2", len(cfg.EnabledTypes))
	placeholder
		if cfg.EnabledTypes[0] != "alipay" || cfg.EnabledTypes[1] != "wxpay" {
			t.Fatalf("EnabledTypes = %v, want [alipay wxpay]", cfg.EnabledTypes)
	placeholder
placeholder)

	t.Run("empty enabled types string", func(t *testing.T) {
		t.Parallel()
		vals := map[string]string{
			SettingEnabledPaymentTypes: "",
	placeholder
		cfg := svc.parsePaymentConfig(vals)
		if len(cfg.EnabledTypes) != 0 {
			t.Fatalf("expected empty EnabledTypes for empty string, got %v", cfg.EnabledTypes)
	placeholder
placeholder)
placeholder

func TestGetBasePaymentType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
placeholder{
		{payment.TypeEasyPay, payment.TypeEasyPayplaceholder,
		{payment.TypeStripe, payment.TypeStripeplaceholder,
		{payment.TypeCard, payment.TypeStripeplaceholder,
		{payment.TypeLink, payment.TypeStripeplaceholder,
		{payment.TypeAlipay, payment.TypeAlipayplaceholder,
		{payment.TypeAlipayDirect, payment.TypeAlipayplaceholder,
		{payment.TypeWxpay, payment.TypeWxpayplaceholder,
		{payment.TypeWxpayDirect, payment.TypeWxpayplaceholder,
		{"unknown", "unknown"placeholder,
		{"", ""placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			got := payment.GetBasePaymentType(tt.input)
			if got != tt.expected {
				t.Fatalf("GetBasePaymentType(%q) = %q, want %q", tt.input, got, tt.expected)
		placeholder
	placeholder)
placeholder
placeholder
