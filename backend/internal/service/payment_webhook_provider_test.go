//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/payment"
	"github.com/stretchr/testify/require"
)

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

func TestGetWebhookProviderRejectsAmbiguousRegistryFallback(t *testing.T) {
	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	_, err := client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeWxpay).
		SetName("wxpay-a").
		SetConfig("{placeholder").
		SetSupportedTypes("wxpay").
		SetEnabled(true).
		Save(ctx)
placeholder
	_, err = client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeWxpay).
		SetName("wxpay-b").
		SetConfig("{placeholder").
		SetSupportedTypes("wxpay").
		SetEnabled(true).
		Save(ctx)
placeholder

	svc := &PaymentService{
		entClient:       client,
		registry:        payment.NewRegistry(),
		providersLoaded: true,
placeholder

	_, err = svc.GetWebhookProvider(ctx, payment.TypeWxpay, "")
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

	prov, err := svc.GetWebhookProvider(ctx, payment.TypeStripe, "")
placeholder
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

	_, err = svc.GetWebhookProvider(ctx, payment.TypeWxpay, "sub2_test_pinned_order")
placeholder
	require.Contains(t, err.Error(), "provider instance")
placeholder
