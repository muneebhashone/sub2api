//go:build unit

package handler

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/payment"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

func TestApplyWeChatPaymentResumeClaims(t *testing.T) {
	t.Parallel()

	req := CreateOrderRequest{
		Amount:      0,
		PaymentType: payment.TypeWxpay,
		OrderType:   payment.OrderTypeBalance,
placeholder

	err := applyWeChatPaymentResumeClaims(&req, &service.WeChatPaymentResumeClaims{
		OpenID:      "openid-123",
		PaymentType: payment.TypeWxpay,
		Amount:      "12.50",
		OrderType:   payment.OrderTypeSubscription,
		PlanID:      7,
placeholder)
	if err != nil {
		t.Fatalf("applyWeChatPaymentResumeClaims returned error: %v", err)
placeholder
	if req.OpenID != "openid-123" {
		t.Fatalf("openid = %q, want %q", req.OpenID, "openid-123")
placeholder
	if req.Amount != 12.5 {
		t.Fatalf("amount = %v, want 12.5", req.Amount)
placeholder
	if req.OrderType != payment.OrderTypeSubscription {
		t.Fatalf("order_type = %q, want %q", req.OrderType, payment.OrderTypeSubscription)
placeholder
	if req.PlanID != 7 {
		t.Fatalf("plan_id = %d, want 7", req.PlanID)
placeholder
placeholder

func TestApplyWeChatPaymentResumeClaimsRejectsPaymentTypeMismatch(t *testing.T) {
	t.Parallel()

	req := CreateOrderRequest{
		PaymentType: payment.TypeAlipay,
placeholder

	err := applyWeChatPaymentResumeClaims(&req, &service.WeChatPaymentResumeClaims{
		OpenID:      "openid-123",
		PaymentType: payment.TypeWxpay,
		Amount:      "12.50",
		OrderType:   payment.OrderTypeBalance,
placeholder)
	if err == nil {
		t.Fatal("applyWeChatPaymentResumeClaims should reject mismatched payment types")
placeholder
placeholder
