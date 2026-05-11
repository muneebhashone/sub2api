//go:build unit

package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/enttest"
	"github.com/Wei-Shaw/sub2api/internal/payment"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "modernc.org/sqlite"
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

func TestVerifyOrderPublicReturnsLegacyOrderState(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	db, err := sql.Open("sqlite", "file:payment_handler_public_verify?mode=memory&cache=shared")
placeholder
	t.Cleanup(func() { _ = db.Close() placeholder)

	_, err = db.Exec("PRAGMA foreign_keys = ON")
placeholder

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := enttest.NewClient(t, enttest.WithOptions(dbent.Driver(drv)))
	t.Cleanup(func() { _ = client.Close() placeholder)

	user, err := client.User.Create().
		SetEmail("public-verify@example.com").
		SetPasswordHash("hash").
		SetUsername("public-verify-user").
		Save(context.Background())
placeholder

	order, err := client.PaymentOrder.Create().
		SetUserID(user.ID).
		SetUserEmail(user.Email).
		SetUserName(user.Username).
		SetAmount(88).
		SetPayAmount(90.64).
		SetFeeRate(0.03).
		SetRechargeCode("PUBLIC-VERIFY").
		SetOutTradeNo("legacy-order-no").
		SetPaymentType(payment.TypeAlipay).
		SetPaymentTradeNo("trade-public-verify").
		SetOrderType(payment.OrderTypeBalance).
		SetStatus(service.OrderStatusPending).
		SetExpiresAt(time.Now().Add(time.Hour)).
		SetClientIP("127.0.0.1").
		SetSrcHost("api.example.com").
		SetProviderSnapshot(map[string]any{"currency": "HKD"placeholder).
		Save(context.Background())
placeholder

	paymentSvc := service.NewPaymentService(client, payment.NewRegistry(), nil, nil, nil, nil, nil, nil, nil)
	h := NewPaymentHandler(paymentSvc, nil, nil)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(
		http.MethodPost,
		"/api/v1/payment/public/orders/verify",
		bytes.NewBufferString(`{"out_trade_no":"legacy-order-no"placeholder`),
	)
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.VerifyOrderPublic(ctx)

	require.Equal(t, http.StatusOK, recorder.Code)

	var resp struct {
		Code int `json:"code"`
		Data struct {
			ID           int64   `json:"id"`
			OutTradeNo   string  `json:"out_trade_no"`
			Amount       float64 `json:"amount"`
			PayAmount    float64 `json:"pay_amount"`
			FeeRate      float64 `json:"fee_rate"`
			Currency     string  `json:"currency"`
			PaymentType  string  `json:"payment_type"`
			OrderType    string  `json:"order_type"`
			Status       string  `json:"status"`
			RefundAmount float64 `json:"refund_amount"`
			CreatedAt    string  `json:"created_at"`
			ExpiresAt    string  `json:"expires_at"`
	placeholder `json:"data"`
placeholder
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	require.Equal(t, 0, resp.Code)
	require.Equal(t, order.ID, resp.Data.ID)
	require.Equal(t, "legacy-order-no", resp.Data.OutTradeNo)
	require.Equal(t, 90.64, resp.Data.PayAmount)
	require.Equal(t, 0.03, resp.Data.FeeRate)
	require.Equal(t, "HKD", resp.Data.Currency)
	require.Equal(t, payment.TypeAlipay, resp.Data.PaymentType)
	require.Equal(t, payment.OrderTypeBalance, resp.Data.OrderType)
	require.Equal(t, service.OrderStatusPending, resp.Data.Status)
	require.Equal(t, 0.0, resp.Data.RefundAmount)
	require.NotEmpty(t, resp.Data.CreatedAt)
	require.NotEmpty(t, resp.Data.ExpiresAt)
placeholder

func TestResolveOrderPublicByResumeTokenReturnsFrontendContractFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("PAYMENT_RESUME_SIGNING_KEY", "placeholder")

	db, err := sql.Open("sqlite", "file:payment_handler_public_resolve?mode=memory&cache=shared")
placeholder
	t.Cleanup(func() { _ = db.Close() placeholder)

	_, err = db.Exec("PRAGMA foreign_keys = ON")
placeholder

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := enttest.NewClient(t, enttest.WithOptions(dbent.Driver(drv)))
	t.Cleanup(func() { _ = client.Close() placeholder)

	user, err := client.User.Create().
		SetEmail("public-resolve@example.com").
		SetPasswordHash("hash").
		SetUsername("public-resolve-user").
		Save(context.Background())
placeholder

	order, err := client.PaymentOrder.Create().
		SetUserID(user.ID).
		SetUserEmail(user.Email).
		SetUserName(user.Username).
		SetAmount(100).
		SetPayAmount(103).
		SetFeeRate(0.03).
		SetRechargeCode("PUBLIC-RESOLVE").
		SetOutTradeNo("resolve-order-no").
		SetPaymentType(payment.TypeAlipay).
		SetPaymentTradeNo("trade-public-resolve").
		SetOrderType(payment.OrderTypeBalance).
		SetStatus(service.OrderStatusPaid).
		SetExpiresAt(time.Now().Add(time.Hour)).
		SetPaidAt(time.Now()).
		SetClientIP("127.0.0.1").
		SetSrcHost("api.example.com").
		SetProviderSnapshot(map[string]any{"currency": "USD"placeholder).
		Save(context.Background())
placeholder

	resumeSvc := service.NewPaymentResumeService([]byte("placeholder"))
	token, err := resumeSvc.CreateToken(service.ResumeTokenClaims{
		OrderID:            order.ID,
		UserID:             user.ID,
		PaymentType:        payment.TypeAlipay,
		CanonicalReturnURL: "https://app.example.com/payment/result",
placeholder)
placeholder

	configSvc := service.NewPaymentConfigService(client, nil, []byte("placeholder"))
	paymentSvc := service.NewPaymentService(client, payment.NewRegistry(), nil, nil, nil, configSvc, nil, nil, nil)
	h := NewPaymentHandler(paymentSvc, nil, nil)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(
		http.MethodPost,
		"/api/v1/payment/public/orders/resolve",
		bytes.NewBufferString(`{"resume_token":"`+token+`"placeholder`),
	)
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.ResolveOrderPublicByResumeToken(ctx)

	require.Equal(t, http.StatusOK, recorder.Code)

	var resp struct {
		Code int            `json:"code"`
		Data map[string]any `json:"data"`
placeholder
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	require.Equal(t, 0, resp.Code)
	require.Equal(t, float64(order.ID), resp.Data["id"])
	require.Equal(t, "resolve-order-no", resp.Data["out_trade_no"])
	require.Equal(t, 100.0, resp.Data["amount"])
	require.Equal(t, 103.0, resp.Data["pay_amount"])
	require.Equal(t, 0.03, resp.Data["fee_rate"])
	require.Equal(t, "USD", resp.Data["currency"])
	require.Equal(t, payment.TypeAlipay, resp.Data["payment_type"])
	require.Equal(t, payment.OrderTypeBalance, resp.Data["order_type"])
	require.Equal(t, service.OrderStatusPaid, resp.Data["status"])
	require.Contains(t, resp.Data, "created_at")
	require.Contains(t, resp.Data, "expires_at")
	require.Contains(t, resp.Data, "refund_amount")
placeholder

func TestResolveOrderPublicByResumeTokenReturnsBadRequestForMismatchedToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("PAYMENT_RESUME_SIGNING_KEY", "placeholder")

	db, err := sql.Open("sqlite", "file:payment_handler_public_resolve_mismatch?mode=memory&cache=shared")
placeholder
	t.Cleanup(func() { _ = db.Close() placeholder)

	_, err = db.Exec("PRAGMA foreign_keys = ON")
placeholder

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := enttest.NewClient(t, enttest.WithOptions(dbent.Driver(drv)))
	t.Cleanup(func() { _ = client.Close() placeholder)

	user, err := client.User.Create().
		SetEmail("public-resolve-mismatch@example.com").
		SetPasswordHash("hash").
		SetUsername("public-resolve-mismatch-user").
		Save(context.Background())
placeholder

	order, err := client.PaymentOrder.Create().
		SetUserID(user.ID).
		SetUserEmail(user.Email).
		SetUserName(user.Username).
		SetAmount(100).
		SetPayAmount(103).
		SetFeeRate(0.03).
		SetRechargeCode("PUBLIC-RESOLVE-MISMATCH").
		SetOutTradeNo("resolve-order-mismatch-no").
		SetPaymentType(payment.TypeAlipay).
		SetPaymentTradeNo("trade-public-resolve-mismatch").
		SetOrderType(payment.OrderTypeBalance).
		SetStatus(service.OrderStatusPaid).
		SetExpiresAt(time.Now().Add(time.Hour)).
		SetPaidAt(time.Now()).
		SetClientIP("127.0.0.1").
		SetSrcHost("api.example.com").
		Save(context.Background())
placeholder

	resumeSvc := service.NewPaymentResumeService([]byte("placeholder"))
	token, err := resumeSvc.CreateToken(service.ResumeTokenClaims{
		OrderID:            order.ID,
		UserID:             user.ID + 999,
		PaymentType:        payment.TypeAlipay,
		CanonicalReturnURL: "https://app.example.com/payment/result",
placeholder)
placeholder

	configSvc := service.NewPaymentConfigService(client, nil, []byte("placeholder"))
	paymentSvc := service.NewPaymentService(client, payment.NewRegistry(), nil, nil, nil, configSvc, nil, nil, nil)
	h := NewPaymentHandler(paymentSvc, nil, nil)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(
		http.MethodPost,
		"/api/v1/payment/public/orders/resolve",
		bytes.NewBufferString(`{"resume_token":"`+token+`"placeholder`),
	)
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.ResolveOrderPublicByResumeToken(ctx)

	require.Equal(t, http.StatusBadRequest, recorder.Code)

	var resp struct {
		Code    int    `json:"code"`
		Reason  string `json:"reason"`
		Message string `json:"message"`
placeholder
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Equal(t, "INVALID_RESUME_TOKEN", resp.Reason)
placeholder

func TestVerifyOrderPublicRejectsBlankOutTradeNo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := sql.Open("sqlite", "file:payment_handler_public_verify_blank?mode=memory&cache=shared")
placeholder
	t.Cleanup(func() { _ = db.Close() placeholder)

	_, err = db.Exec("PRAGMA foreign_keys = ON")
placeholder

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := enttest.NewClient(t, enttest.WithOptions(dbent.Driver(drv)))
	t.Cleanup(func() { _ = client.Close() placeholder)

	paymentSvc := service.NewPaymentService(client, payment.NewRegistry(), nil, nil, nil, nil, nil, nil, nil)
	h := NewPaymentHandler(paymentSvc, nil, nil)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(
		http.MethodPost,
		"/api/v1/payment/public/orders/verify",
		bytes.NewBufferString(`{"out_trade_no":"   "placeholder`),
	)
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.VerifyOrderPublic(ctx)

	require.Equal(t, http.StatusBadRequest, recorder.Code)

	var resp struct {
		Code   int    `json:"code"`
		Reason string `json:"reason"`
placeholder
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Equal(t, "INVALID_OUT_TRADE_NO", resp.Reason)
placeholder
