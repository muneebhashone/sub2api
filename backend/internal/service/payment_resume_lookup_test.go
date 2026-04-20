//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/payment"
	"github.com/stretchr/testify/require"
)

type paymentResumeLookupProvider struct {
	queryCount int
placeholder

func (p *paymentResumeLookupProvider) Name() string { return "resume-lookup-provider" placeholder

func (p *paymentResumeLookupProvider) ProviderKey() string { return payment.TypeAlipay placeholder

func (p *paymentResumeLookupProvider) SupportedTypes() []payment.PaymentType {
	return []payment.PaymentType{payment.TypeAlipayplaceholder
placeholder

func (p *paymentResumeLookupProvider) CreatePayment(context.Context, payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	panic("unexpected call")
placeholder

func (p *paymentResumeLookupProvider) QueryOrder(context.Context, string) (*payment.QueryOrderResponse, error) {
	p.queryCount++
	return &payment.QueryOrderResponse{Status: payment.ProviderStatusPendingplaceholder, nil
placeholder

func (p *paymentResumeLookupProvider) VerifyNotification(context.Context, string, map[string]string) (*payment.PaymentNotification, error) {
	panic("unexpected call")
placeholder

func (p *paymentResumeLookupProvider) Refund(context.Context, payment.RefundRequest) (*payment.RefundResponse, error) {
	panic("unexpected call")
placeholder

func TestGetPublicOrderByResumeTokenReturnsMatchingOrder(t *testing.T) {
	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	user, err := client.User.Create().
		SetEmail("resume@example.com").
		SetPasswordHash("hash").
		SetUsername("resume-user").
		Save(ctx)
placeholder

	instanceID := "12"
	providerKey := payment.TypeEasyPay
	order, err := client.PaymentOrder.Create().
		SetUserID(user.ID).
		SetUserEmail(user.Email).
		SetUserName(user.Username).
		SetAmount(88).
		SetPayAmount(88).
		SetFeeRate(0).
		SetRechargeCode("RESUME-ORDER").
		SetOutTradeNo("sub2_resume_lookup").
		SetPaymentType(payment.TypeAlipay).
		SetPaymentTradeNo("trade-1").
		SetOrderType(payment.OrderTypeBalance).
		SetStatus(OrderStatusPending).
		SetExpiresAt(time.Now().Add(time.Hour)).
		SetClientIP("127.0.0.1").
		SetSrcHost("api.example.com").
		SetProviderInstanceID(instanceID).
		SetProviderKey(providerKey).
		Save(ctx)
placeholder

	resumeSvc := NewPaymentResumeService([]byte("placeholder"))
	token, err := resumeSvc.CreateToken(ResumeTokenClaims{
		OrderID:            order.ID,
		UserID:             user.ID,
		ProviderInstanceID: instanceID,
		ProviderKey:        providerKey,
		PaymentType:        payment.TypeAlipay,
		CanonicalReturnURL: "https://app.example.com/payment/result",
placeholder)
placeholder

	svc := &PaymentService{
		entClient:     client,
		resumeService: resumeSvc,
placeholder

	got, err := svc.GetPublicOrderByResumeToken(ctx, token)
placeholder
	require.Equal(t, order.ID, got.ID)
placeholder

func TestGetPublicOrderByResumeTokenRejectsSnapshotMismatch(t *testing.T) {
	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	user, err := client.User.Create().
		SetEmail("resume-mismatch@example.com").
		SetPasswordHash("hash").
		SetUsername("resume-mismatch-user").
		Save(ctx)
placeholder

	order, err := client.PaymentOrder.Create().
		SetUserID(user.ID).
		SetUserEmail(user.Email).
		SetUserName(user.Username).
		SetAmount(88).
		SetPayAmount(88).
		SetFeeRate(0).
		SetRechargeCode("RESUME-MISMATCH").
		SetOutTradeNo("sub2_resume_lookup_mismatch").
		SetPaymentType(payment.TypeAlipay).
		SetPaymentTradeNo("trade-2").
		SetOrderType(payment.OrderTypeBalance).
		SetStatus(OrderStatusPending).
		SetExpiresAt(time.Now().Add(time.Hour)).
		SetClientIP("127.0.0.1").
		SetSrcHost("api.example.com").
		SetProviderInstanceID("12").
		SetProviderKey(payment.TypeEasyPay).
		Save(ctx)
placeholder

	resumeSvc := NewPaymentResumeService([]byte("placeholder"))
	token, err := resumeSvc.CreateToken(ResumeTokenClaims{
		OrderID:            order.ID,
		UserID:             user.ID,
		ProviderInstanceID: "99",
		ProviderKey:        payment.TypeEasyPay,
		PaymentType:        payment.TypeAlipay,
		CanonicalReturnURL: "https://app.example.com/payment/result",
placeholder)
placeholder

	svc := &PaymentService{
		entClient:     client,
		resumeService: resumeSvc,
placeholder

	_, err = svc.GetPublicOrderByResumeToken(ctx, token)
placeholder
	require.Contains(t, err.Error(), "resume token")
placeholder

func TestGetPublicOrderByResumeTokenChecksUpstreamForPendingOrder(t *testing.T) {
	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	user, err := client.User.Create().
		SetEmail("resume-refresh@example.com").
		SetPasswordHash("hash").
		SetUsername("resume-refresh-user").
		Save(ctx)
placeholder

	order, err := client.PaymentOrder.Create().
		SetUserID(user.ID).
		SetUserEmail(user.Email).
		SetUserName(user.Username).
		SetAmount(88).
		SetPayAmount(88).
		SetFeeRate(0).
		SetRechargeCode("RESUME-PENDING").
		SetOutTradeNo("sub2_resume_lookup_pending").
		SetPaymentType(payment.TypeAlipay).
		SetPaymentTradeNo("trade-pending").
		SetOrderType(payment.OrderTypeBalance).
		SetStatus(OrderStatusPending).
		SetExpiresAt(time.Now().Add(time.Hour)).
		SetClientIP("127.0.0.1").
		SetSrcHost("api.example.com").
		Save(ctx)
placeholder

	resumeSvc := NewPaymentResumeService([]byte("placeholder"))
	token, err := resumeSvc.CreateToken(ResumeTokenClaims{
		OrderID:            order.ID,
		UserID:             user.ID,
		PaymentType:        payment.TypeAlipay,
		CanonicalReturnURL: "https://app.example.com/payment/result",
placeholder)
placeholder

	registry := payment.NewRegistry()
	provider := &paymentResumeLookupProvider{placeholder
	registry.Register(provider)

	svc := &PaymentService{
		entClient:       client,
		registry:        registry,
		resumeService:   resumeSvc,
		providersLoaded: true,
placeholder

	got, err := svc.GetPublicOrderByResumeToken(ctx, token)
placeholder
	require.Equal(t, order.ID, got.ID)
	require.Equal(t, 1, provider.queryCount)
placeholder
