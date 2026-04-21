//go:build unit

package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"strconv"
	"testing"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/payment"
	"github.com/stretchr/testify/require"
)

const webhookProviderTestEncryptionKey = "placeholder"

type webhookProviderTestDouble struct {
	key   string
	types []payment.PaymentType
placeholder

func (p webhookProviderTestDouble) Name() string                          { return p.key placeholder
func (p webhookProviderTestDouble) ProviderKey() string                   { return p.key placeholder
func (p webhookProviderTestDouble) SupportedTypes() []payment.PaymentType { return p.types placeholder
func (p webhookProviderTestDouble) CreatePayment(context.Context, payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	panic("unexpected call")
placeholder
func (p webhookProviderTestDouble) QueryOrder(context.Context, string) (*payment.QueryOrderResponse, error) {
	panic("unexpected call")
placeholder
func (p webhookProviderTestDouble) VerifyNotification(context.Context, string, map[string]string) (*payment.PaymentNotification, error) {
	panic("unexpected call")
placeholder
func (p webhookProviderTestDouble) Refund(context.Context, payment.RefundRequest) (*payment.RefundResponse, error) {
	panic("unexpected call")
placeholder

func encryptWebhookProviderConfig(t *testing.T, config map[string]string) string {
placeholder

	data, err := json.Marshal(config)
placeholder

	encrypted, err := payment.Encrypt(string(data), []byte(webhookProviderTestEncryptionKey))
placeholder
	return encrypted
placeholder

func newWebhookProviderTestLoadBalancer(client *dbent.Client) payment.LoadBalancer {
	return payment.NewDefaultLoadBalancer(client, []byte(webhookProviderTestEncryptionKey))
placeholder

func encryptValidWebhookWxpayConfig(t *testing.T, suffix string) string {
placeholder

	key, err := rsa.GenerateKey(rand.Reader, 2048)
placeholder

	privDER, err := x509.MarshalPKCS8PrivateKey(key)
placeholder
	pubDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
placeholder

	return encryptWebhookProviderConfig(t, map[string]string{
		"appId":       "wx-app-" + suffix,
		"mchId":       "mch-" + suffix,
		"privateKey":  string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDERplaceholder)),
		"apiV3Key":    webhookProviderTestEncryptionKey,
		"publicKey":   string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDERplaceholder)),
		"publicKeyId": "public-key-id-" + suffix,
		"certSerial":  "cert-serial-" + suffix,
placeholder)
placeholder

func TestGetOrderProviderInstanceResolvesUniqueLegacyProviderKey(t *testing.T) {
	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	inst, err := client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeStripe).
		SetName("stripe-a").
		SetConfig(encryptWebhookProviderConfig(t, map[string]string{"secretKey": "sk_test_legacy_provider_key"placeholder)).
		SetSupportedTypes("stripe").
		SetEnabled(true).
		Save(ctx)
placeholder

	providerKey := payment.TypeStripe
	order := &dbent.PaymentOrder{
		PaymentType: payment.TypeStripe,
		ProviderKey: &providerKey,
placeholder

	svc := &PaymentService{
		entClient:    client,
		loadBalancer: newWebhookProviderTestLoadBalancer(client),
placeholder

	got, err := svc.getOrderProviderInstance(ctx, order)
placeholder
	require.NotNil(t, got)
	require.Equal(t, inst.ID, got.ID)
placeholder

func TestGetOrderProviderInstanceResolvesUniqueLegacyPaymentType(t *testing.T) {
	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	inst, err := client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeWxpay).
		SetName("wxpay-a").
		SetConfig("{placeholder").
		SetSupportedTypes("wxpay").
		SetEnabled(true).
		Save(ctx)
placeholder

	order := &dbent.PaymentOrder{
		PaymentType: payment.TypeWxpayDirect,
placeholder

	svc := &PaymentService{
		entClient:    client,
		loadBalancer: newWebhookProviderTestLoadBalancer(client),
placeholder

	got, err := svc.getOrderProviderInstance(ctx, order)
placeholder
	require.NotNil(t, got)
	require.Equal(t, inst.ID, got.ID)
placeholder

func TestGetOrderProviderInstanceLeavesAmbiguousLegacyOrderUnresolved(t *testing.T) {
	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	_, err := client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeEasyPay).
		SetName("easypay-a").
		SetConfig("{placeholder").
		SetSupportedTypes("wxpay").
		SetEnabled(true).
		Save(ctx)
placeholder
	_, err = client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeWxpay).
		SetName("wxpay-a").
		SetConfig("{placeholder").
		SetSupportedTypes("wxpay").
		SetEnabled(true).
		Save(ctx)
placeholder

	order := &dbent.PaymentOrder{
		PaymentType: payment.TypeWxpay,
placeholder

	svc := &PaymentService{
		entClient:    client,
		loadBalancer: newWebhookProviderTestLoadBalancer(client),
placeholder

	got, err := svc.getOrderProviderInstance(ctx, order)
placeholder
	require.Nil(t, got)
placeholder

func TestGetOrderProviderInstanceLeavesLegacyProviderKeyUnresolvedWhenHistoricalInstancesConflict(t *testing.T) {
	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	_, err := client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeStripe).
		SetName("stripe-disabled-legacy").
		SetConfig("{placeholder").
		SetSupportedTypes("stripe").
		SetEnabled(false).
		Save(ctx)
placeholder
	_, err = client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeStripe).
		SetName("stripe-enabled-current").
		SetConfig("{placeholder").
		SetSupportedTypes("stripe").
		SetEnabled(true).
		Save(ctx)
placeholder

	providerKey := payment.TypeStripe
	order := &dbent.PaymentOrder{
		PaymentType: payment.TypeStripe,
		ProviderKey: &providerKey,
placeholder

	svc := &PaymentService{
		entClient:    client,
		loadBalancer: newWebhookProviderTestLoadBalancer(client),
placeholder

	got, err := svc.getOrderProviderInstance(ctx, order)
placeholder
	require.Nil(t, got)
placeholder

func TestGetOrderProviderInstanceLeavesProviderKeyMatchUnresolvedWhenTypeNotSupported(t *testing.T) {
	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	_, err := client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeWxpay).
		SetName("wxpay-only").
		SetConfig("{placeholder").
		SetSupportedTypes("wxpay").
		SetEnabled(true).
		Save(ctx)
placeholder

	providerKey := payment.TypeWxpay
	order := &dbent.PaymentOrder{
		PaymentType: payment.TypeAlipayDirect,
		ProviderKey: &providerKey,
placeholder

	svc := &PaymentService{
		entClient:    client,
		loadBalancer: newWebhookProviderTestLoadBalancer(client),
placeholder

	got, err := svc.getOrderProviderInstance(ctx, order)
placeholder
	require.Nil(t, got)
placeholder

func TestGetOrderProviderInstanceUsesProviderSnapshotWhenPinnedColumnMissing(t *testing.T) {
	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	inst, err := client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeStripe).
		SetName("stripe-snapshot").
		SetConfig(encryptWebhookProviderConfig(t, map[string]string{"secretKey": "sk_snapshot"placeholder)).
		SetSupportedTypes("stripe").
		SetEnabled(true).
		Save(ctx)
placeholder

	order := &dbent.PaymentOrder{
		ID:          42,
		PaymentType: payment.TypeStripe,
		ProviderSnapshot: map[string]any{
			"schema_version":       1,
			"provider_instance_id": strconv.FormatInt(inst.ID, 10),
			"provider_key":         payment.TypeStripe,
	placeholder,
placeholder

	svc := &PaymentService{
		entClient:    client,
		loadBalancer: newWebhookProviderTestLoadBalancer(client),
placeholder

	got, err := svc.getOrderProviderInstance(ctx, order)
placeholder
	require.NotNil(t, got)
	require.Equal(t, inst.ID, got.ID)
placeholder

func TestGetOrderProviderInstanceRejectsMissingSnapshotInstanceWithoutLegacyFallback(t *testing.T) {
	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	_, err := client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeStripe).
		SetName("stripe-legacy-fallback").
		SetConfig(encryptWebhookProviderConfig(t, map[string]string{"secretKey": "sk_legacy"placeholder)).
		SetSupportedTypes("stripe").
		SetEnabled(true).
		Save(ctx)
placeholder

	order := &dbent.PaymentOrder{
		ID:          43,
		PaymentType: payment.TypeStripe,
		ProviderSnapshot: map[string]any{
			"schema_version":       1,
			"provider_instance_id": "999999",
			"provider_key":         payment.TypeStripe,
	placeholder,
placeholder

	svc := &PaymentService{
		entClient:    client,
		loadBalancer: newWebhookProviderTestLoadBalancer(client),
placeholder

	got, err := svc.getOrderProviderInstance(ctx, order)
	require.Nil(t, got)
placeholder
	require.Contains(t, err.Error(), "provider snapshot instance 999999 is missing")
placeholder

func TestGetWebhookProviderRejectsAmbiguousRegistryFallback(t *testing.T) {
	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	wxpayConfigA := encryptValidWebhookWxpayConfig(t, "a")
	wxpayConfigB := encryptValidWebhookWxpayConfig(t, "b")
	_, err := client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeWxpay).
		SetName("wxpay-a").
		SetConfig(wxpayConfigA).
		SetSupportedTypes("wxpay").
		SetEnabled(true).
		Save(ctx)
placeholder
	_, err = client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeWxpay).
		SetName("wxpay-b").
		SetConfig(wxpayConfigB).
		SetSupportedTypes("wxpay").
		SetEnabled(true).
		Save(ctx)
placeholder

	svc := &PaymentService{
		entClient:       client,
		loadBalancer:    newWebhookProviderTestLoadBalancer(client),
		registry:        payment.NewRegistry(),
		providersLoaded: true,
placeholder

	providers, err := svc.GetWebhookProviders(ctx, payment.TypeWxpay, "")
placeholder
	require.Len(t, providers, 2)
placeholder

func TestGetWebhookProvidersRejectAmbiguousFallbackForNonWxpay(t *testing.T) {
	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	_, err := client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeAlipay).
		SetName("alipay-a").
		SetConfig("{placeholder").
		SetSupportedTypes("alipay").
		SetEnabled(true).
		Save(ctx)
placeholder
	_, err = client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeAlipay).
		SetName("alipay-b").
		SetConfig("{placeholder").
		SetSupportedTypes("alipay").
		SetEnabled(true).
		Save(ctx)
placeholder

	svc := &PaymentService{
		entClient:       client,
		registry:        payment.NewRegistry(),
		providersLoaded: true,
placeholder

	_, err = svc.GetWebhookProviders(ctx, payment.TypeAlipay, "")
placeholder
	require.Contains(t, err.Error(), "ambiguous")
placeholder

func TestGetWebhookProviderAllowsSingleInstanceRegistryFallback(t *testing.T) {
	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	_, err := client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeStripe).
		SetName("stripe-a").
		SetConfig("{placeholder").
		SetSupportedTypes("stripe").
		SetEnabled(true).
		Save(ctx)
placeholder

	registry := payment.NewRegistry()
	registry.Register(webhookProviderTestDouble{
		key:   payment.TypeStripe,
		types: []payment.PaymentType{payment.TypeStripeplaceholder,
placeholder)

	svc := &PaymentService{
		entClient:       client,
		registry:        registry,
		providersLoaded: true,
placeholder

	providers, err := svc.GetWebhookProviders(ctx, payment.TypeStripe, "")
placeholder
	require.Len(t, providers, 1)
	prov := providers[0]
	require.Equal(t, payment.TypeStripe, prov.ProviderKey())
placeholder

func TestGetWebhookProviderRejectsRegistryFallbackForPinnedOrder(t *testing.T) {
	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	user, err := client.User.Create().
		SetEmail("webhook@example.com").
		SetPasswordHash("hash").
		SetUsername("webhook").
		Save(ctx)
placeholder

	pinnedInstanceID := "999"
	_, err = client.PaymentOrder.Create().
		SetUserID(user.ID).
		SetUserEmail(user.Email).
		SetUserName(user.Username).
		SetAmount(88).
		SetPayAmount(88).
		SetFeeRate(0).
		SetRechargeCode("TEST-RECHARGE").
		SetOutTradeNo("sub2_test_pinned_order").
		SetPaymentType(payment.TypeWxpay).
		SetPaymentTradeNo("").
		SetOrderType(payment.OrderTypeBalance).
		SetStatus(OrderStatusPending).
		SetExpiresAt(time.Now().Add(time.Hour)).
		SetClientIP("127.0.0.1").
		SetSrcHost("api.example.com").
		SetProviderInstanceID(pinnedInstanceID).
		Save(ctx)
placeholder

	registry := payment.NewRegistry()
	registry.Register(webhookProviderTestDouble{
		key:   payment.TypeWxpay,
		types: []payment.PaymentType{payment.TypeWxpayplaceholder,
placeholder)

	svc := &PaymentService{
		entClient:       client,
		registry:        registry,
		providersLoaded: true,
placeholder

	_, err = svc.GetWebhookProviders(ctx, payment.TypeWxpay, "sub2_test_pinned_order")
placeholder
	require.Contains(t, err.Error(), "provider instance")
placeholder

func TestGetWebhookProviderUsesProviderSnapshotBeforeWxpayFallback(t *testing.T) {
	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	user, err := client.User.Create().
		SetEmail("snapshot-webhook@example.com").
		SetPasswordHash("hash").
		SetUsername("snapshot-webhook").
		Save(ctx)
placeholder

	wxpayConfigA := encryptValidWebhookWxpayConfig(t, "snapshot-a")
	wxpayConfigB := encryptValidWebhookWxpayConfig(t, "snapshot-b")
	instA, err := client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeWxpay).
		SetName("wxpay-snapshot-a").
		SetConfig(wxpayConfigA).
		SetSupportedTypes("wxpay").
		SetEnabled(true).
		Save(ctx)
placeholder
	_, err = client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeWxpay).
		SetName("wxpay-snapshot-b").
		SetConfig(wxpayConfigB).
		SetSupportedTypes("wxpay").
		SetEnabled(true).
		Save(ctx)
placeholder

	_, err = client.PaymentOrder.Create().
		SetUserID(user.ID).
		SetUserEmail(user.Email).
		SetUserName(user.Username).
		SetAmount(66).
		SetPayAmount(66).
		SetFeeRate(0).
		SetRechargeCode("SNAPSHOT-WEBHOOK").
		SetOutTradeNo("sub2_test_snapshot_webhook_order").
		SetPaymentType(payment.TypeWxpay).
		SetPaymentTradeNo("").
		SetOrderType(payment.OrderTypeBalance).
		SetStatus(OrderStatusPending).
		SetExpiresAt(time.Now().Add(time.Hour)).
		SetClientIP("127.0.0.1").
		SetSrcHost("api.example.com").
		SetProviderSnapshot(map[string]any{
			"schema_version":       1,
			"provider_instance_id": strconv.FormatInt(instA.ID, 10),
			"provider_key":         payment.TypeWxpay,
			"payment_mode":         "native",
	placeholder).
		Save(ctx)
placeholder

	svc := &PaymentService{
		entClient:       client,
		loadBalancer:    newWebhookProviderTestLoadBalancer(client),
		registry:        payment.NewRegistry(),
		providersLoaded: true,
placeholder

	providers, err := svc.GetWebhookProviders(ctx, payment.TypeWxpay, "sub2_test_snapshot_webhook_order")
placeholder
	require.Len(t, providers, 1)
	require.Equal(t, payment.TypeWxpay, providers[0].ProviderKey())
placeholder
