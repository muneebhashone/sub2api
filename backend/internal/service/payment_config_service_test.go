package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/enttest"
	"github.com/Wei-Shaw/sub2api/internal/payment"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "modernc.org/sqlite"
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

func TestAlipayMobilePrecreateEnvironmentOverride(t *testing.T) {
	svc := &PaymentConfigService{placeholder

	t.Setenv(SettingAlipayMobilePrecreateDeepLink, "true")
	if !svc.parsePaymentConfig(map[string]string{SettingAlipayMobilePrecreateDeepLink: "false"placeholder).AlipayMobilePrecreateDeepLink {
		t.Fatal("expected environment variable to enable mobile Alipay precreate")
placeholder

	t.Setenv(SettingAlipayMobilePrecreateDeepLink, "false")
	if svc.parsePaymentConfig(map[string]string{SettingAlipayMobilePrecreateDeepLink: "true"placeholder).AlipayMobilePrecreateDeepLink {
		t.Fatal("expected environment variable to disable mobile Alipay precreate")
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
		if cfg.AlipayMobilePrecreateDeepLink {
			t.Fatal("expected AlipayMobilePrecreateDeepLink=false by default")
	placeholder
placeholder)

	t.Run("all values populated", func(t *testing.T) {
		t.Parallel()
		vals := map[string]string{
			SettingPaymentEnabled:                "true",
			SettingMinRechargeAmount:             "5.00",
			SettingMaxRechargeAmount:             "1000.00",
			SettingDailyRechargeLimit:            "5000.00",
			SettingOrderTimeoutMinutes:           "15",
			SettingMaxPendingOrders:              "5",
			SettingEnabledPaymentTypes:           "alipay,wxpay,stripe",
			SettingBalancePayDisabled:            "true",
			SettingLoadBalanceStrategy:           "least_amount",
			SettingProductNamePrefix:             "PRE",
			SettingProductNameSuffix:             "SUF",
			SettingAlipayMobilePrecreateDeepLink: "true",
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
		if !cfg.AlipayMobilePrecreateDeepLink {
			t.Fatal("expected AlipayMobilePrecreateDeepLink=true")
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

	t.Run("enabled types are normalized to visible methods and deduplicated", func(t *testing.T) {
		t.Parallel()
		vals := map[string]string{
			SettingEnabledPaymentTypes: "alipay_direct, alipay, wxpay_direct, wxpay",
	placeholder
		cfg := svc.parsePaymentConfig(vals)
		if len(cfg.EnabledTypes) != 2 {
			t.Fatalf("EnabledTypes len = %d, want 2", len(cfg.EnabledTypes))
	placeholder
		if cfg.EnabledTypes[0] != "alipay" || cfg.EnabledTypes[1] != "wxpay" {
			t.Fatalf("EnabledTypes = %v, want [alipay wxpay]", cfg.EnabledTypes)
	placeholder
placeholder)

	t.Run("custom enabled types are preserved", func(t *testing.T) {
		t.Parallel()
		vals := map[string]string{
			SettingEnabledPaymentTypes: "alipay,ldc,usdt_trc20",
	placeholder
		cfg := svc.parsePaymentConfig(vals)
		want := []string{"alipay", "ldc", "usdt_trc20"placeholder
		if len(cfg.EnabledTypes) != len(want) {
			t.Fatalf("EnabledTypes len = %d, want %d (%v)", len(cfg.EnabledTypes), len(want), cfg.EnabledTypes)
	placeholder
		for i := range want {
			if cfg.EnabledTypes[i] != want[i] {
				t.Fatalf("EnabledTypes[%d] = %q, want %q (full=%v)", i, cfg.EnabledTypes[i], want[i], cfg.EnabledTypes)
		placeholder
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

func TestApplyVisibleMethodRoutingToEnabledTypes(t *testing.T) {
	t.Parallel()

	base := []string{"alipay", "wxpay", "stripe"placeholder
	vals := map[string]string{
		SettingPaymentVisibleMethodAlipayEnabled: "true",
		SettingPaymentVisibleMethodAlipaySource:  VisibleMethodSourceOfficialAlipay,
		SettingPaymentVisibleMethodWxpayEnabled:  "true",
		SettingPaymentVisibleMethodWxpaySource:   VisibleMethodSourceOfficialWechat,
placeholder
	available := map[string]bool{
		VisibleMethodSourceOfficialAlipay: true,
		VisibleMethodSourceOfficialWechat: false,
placeholder

	got := applyVisibleMethodRoutingToEnabledTypes(base, vals, available)
	want := []string{"alipay", "stripe"placeholder
	if len(got) != len(want) {
		t.Fatalf("applyVisibleMethodRoutingToEnabledTypes len = %d, want %d (%v)", len(got), len(want), got)
placeholder
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("applyVisibleMethodRoutingToEnabledTypes[%d] = %q, want %q (full=%v)", i, got[i], want[i], got)
	placeholder
placeholder
placeholder

func TestApplyVisibleMethodRoutingAddsConfiguredVisibleMethod(t *testing.T) {
	t.Parallel()

	base := []string{"stripe"placeholder
	vals := map[string]string{
		SettingPaymentVisibleMethodAlipayEnabled: "true",
		SettingPaymentVisibleMethodAlipaySource:  VisibleMethodSourceEasyPayAlipay,
placeholder
	available := map[string]bool{
		VisibleMethodSourceEasyPayAlipay: true,
placeholder

	got := applyVisibleMethodRoutingToEnabledTypes(base, vals, available)
	want := []string{"stripe", "alipay"placeholder
	if len(got) != len(want) {
		t.Fatalf("applyVisibleMethodRoutingToEnabledTypes len = %d, want %d (%v)", len(got), len(want), got)
placeholder
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("applyVisibleMethodRoutingToEnabledTypes[%d] = %q, want %q (full=%v)", i, got[i], want[i], got)
	placeholder
placeholder
placeholder

func TestBuildVisibleMethodSourceAvailability(t *testing.T) {
	t.Parallel()

	instances := []*dbent.PaymentProviderInstance{
		{ProviderKey: payment.TypeAlipay, SupportedTypes: "alipay"placeholder,
		{ProviderKey: payment.TypeEasyPay, SupportedTypes: "wxpay_direct, alipay"placeholder,
		{ProviderKey: payment.TypeWxpay, SupportedTypes: "wxpay_direct"placeholder,
placeholder

	got := buildVisibleMethodSourceAvailability(instances)
	if !got[VisibleMethodSourceOfficialAlipay] {
		t.Fatalf("expected %q to be available", VisibleMethodSourceOfficialAlipay)
placeholder
	if !got[VisibleMethodSourceEasyPayAlipay] {
		t.Fatalf("expected %q to be available", VisibleMethodSourceEasyPayAlipay)
placeholder
	if !got[VisibleMethodSourceOfficialWechat] {
		t.Fatalf("expected %q to be available", VisibleMethodSourceOfficialWechat)
placeholder
	if !got[VisibleMethodSourceEasyPayWechat] {
		t.Fatalf("expected %q to be available", VisibleMethodSourceEasyPayWechat)
placeholder
placeholder

func TestGetPaymentConfigKeepsStoredEnabledTypes(t *testing.T) {
	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)

	_, err := client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeEasyPay).
		SetName("EasyPay Alipay").
		SetConfig("{placeholder").
		SetSupportedTypes("alipay").
		SetEnabled(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("create easypay instance: %v", err)
placeholder

	svc := &PaymentConfigService{
		entClient: client,
		settingRepo: &paymentConfigSettingRepoStub{
			values: map[string]string{
				SettingEnabledPaymentTypes: "alipay,wxpay,stripe",
		placeholder,
	placeholder,
placeholder

	cfg, err := svc.GetPaymentConfig(ctx)
	if err != nil {
		t.Fatalf("GetPaymentConfig returned error: %v", err)
placeholder

	want := []string{payment.TypeAlipay, payment.TypeWxpay, payment.TypeStripeplaceholder
	if len(cfg.EnabledTypes) != len(want) {
		t.Fatalf("EnabledTypes len = %d, want %d (%v)", len(cfg.EnabledTypes), len(want), cfg.EnabledTypes)
placeholder
	for i := range want {
		if cfg.EnabledTypes[i] != want[i] {
			t.Fatalf("EnabledTypes[%d] = %q, want %q (full=%v)", i, cfg.EnabledTypes[i], want[i], cfg.EnabledTypes)
	placeholder
placeholder
placeholder

func newPaymentConfigServiceTestClient(t *testing.T) *dbent.Client {
placeholder

	dbName := fmt.Sprintf(
		"file:%s?mode=memory&cache=shared",
		strings.NewReplacer("/", "_", " ", "_").Replace(t.Name()),
	)
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
placeholder
	t.Cleanup(func() { _ = db.Close() placeholder)

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatalf("enable foreign keys: %v", err)
placeholder

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := enttest.NewClient(t, enttest.WithOptions(dbent.Driver(drv)))
	t.Cleanup(func() { _ = client.Close() placeholder)
	return client
placeholder

type paymentConfigSettingRepoStub struct {
	values  map[string]string
	updates map[string]string
placeholder

func (s *paymentConfigSettingRepoStub) Get(context.Context, string) (*Setting, error) {
	return nil, nil
placeholder
func (s *paymentConfigSettingRepoStub) GetValue(_ context.Context, key string) (string, error) {
	return s.values[key], nil
placeholder
func (s *paymentConfigSettingRepoStub) Set(context.Context, string, string) error { return nil placeholder
func (s *paymentConfigSettingRepoStub) GetMultiple(_ context.Context, keys []string) (map[string]string, error) {
	out := make(map[string]string, len(keys))
	for _, key := range keys {
		out[key] = s.values[key]
placeholder
	return out, nil
placeholder
func (s *paymentConfigSettingRepoStub) SetMultiple(_ context.Context, values map[string]string) error {
	s.updates = make(map[string]string, len(values))
	for key, value := range values {
		s.updates[key] = value
		if s.values == nil {
			s.values = map[string]string{placeholder
	placeholder
		s.values[key] = value
placeholder
	return nil
placeholder
func (s *paymentConfigSettingRepoStub) GetAll(context.Context) (map[string]string, error) {
	return s.values, nil
placeholder
func (s *paymentConfigSettingRepoStub) Delete(context.Context, string) error { return nil placeholder

func TestUpdatePaymentConfig_PersistsVisibleMethodRouting(t *testing.T) {
	repo := &paymentConfigSettingRepoStub{values: map[string]string{placeholderplaceholder
	svc := &PaymentConfigService{settingRepo: repoplaceholder

	alipayEnabled := true
	wxpayEnabled := false
	err := svc.UpdatePaymentConfig(context.Background(), UpdatePaymentConfigRequest{
		VisibleMethodAlipayEnabled: &alipayEnabled,
		VisibleMethodAlipaySource:  paymentConfigStrPtr(VisibleMethodSourceEasyPayAlipay),
		VisibleMethodWxpayEnabled:  &wxpayEnabled,
		VisibleMethodWxpaySource:   paymentConfigStrPtr(VisibleMethodSourceOfficialWechat),
placeholder)
	if err != nil {
		t.Fatalf("UpdatePaymentConfig returned error: %v", err)
placeholder

	if repo.values[SettingPaymentVisibleMethodAlipayEnabled] != "true" {
		t.Fatalf("alipay enabled = %q, want true", repo.values[SettingPaymentVisibleMethodAlipayEnabled])
placeholder
	if repo.values[SettingPaymentVisibleMethodAlipaySource] != VisibleMethodSourceEasyPayAlipay {
		t.Fatalf("alipay source = %q, want %q", repo.values[SettingPaymentVisibleMethodAlipaySource], VisibleMethodSourceEasyPayAlipay)
placeholder
	if repo.values[SettingPaymentVisibleMethodWxpayEnabled] != "false" {
		t.Fatalf("wxpay enabled = %q, want false", repo.values[SettingPaymentVisibleMethodWxpayEnabled])
placeholder
	if repo.values[SettingPaymentVisibleMethodWxpaySource] != VisibleMethodSourceOfficialWechat {
		t.Fatalf("wxpay source = %q, want %q", repo.values[SettingPaymentVisibleMethodWxpaySource], VisibleMethodSourceOfficialWechat)
placeholder
placeholder

func paymentConfigStrPtr(value string) *string {
	return &value
placeholder
