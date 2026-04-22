package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/paymentauditlog"
	"github.com/Wei-Shaw/sub2api/ent/paymentorder"
	"github.com/Wei-Shaw/sub2api/internal/payment"
	"github.com/Wei-Shaw/sub2api/internal/payment/provider"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

// --- Cancel & Expire ---

// Cancel rate limit configuration constants.
const (
	rateLimitUnitDay           = "day"
	rateLimitUnitMinute        = "minute"
	rateLimitUnitHour          = "hour"
	rateLimitModeFixed         = "fixed"
	checkPaidResultAlreadyPaid = "already_paid"
	checkPaidResultCancelled   = "cancelled"
)

func (s *PaymentService) checkCancelRateLimit(ctx context.Context, userID int64, cfg *PaymentConfig) error {
	if !cfg.CancelRateLimitEnabled || cfg.CancelRateLimitMax <= 0 {
		return nil
placeholder
	windowStart := cancelRateLimitWindowStart(cfg)
	operator := fmt.Sprintf("user:%d", userID)
	count, err := s.entClient.PaymentAuditLog.Query().
		Where(
			paymentauditlog.ActionEQ("ORDER_CANCELLED"),
			paymentauditlog.OperatorEQ(operator),
			paymentauditlog.CreatedAtGTE(windowStart),
		).Count(ctx)
	if err != nil {
		slog.Error("check cancel rate limit failed", "userID", userID, "error", err)
		return nil // fail open
placeholder
	if count >= cfg.CancelRateLimitMax {
		return infraerrors.TooManyRequests("CANCEL_RATE_LIMITED", "cancel rate limited").
			WithMetadata(map[string]string{
				"max":    strconv.Itoa(cfg.CancelRateLimitMax),
				"window": strconv.Itoa(cfg.CancelRateLimitWindow),
				"unit":   cfg.CancelRateLimitUnit,
		placeholder)
placeholder
	return nil
placeholder

func cancelRateLimitWindowStart(cfg *PaymentConfig) time.Time {
	now := time.Now()
	w := cfg.CancelRateLimitWindow
	if w <= 0 {
		w = 1
placeholder
	unit := cfg.CancelRateLimitUnit
	if unit == "" {
		unit = rateLimitUnitDay
placeholder
	if cfg.CancelRateLimitMode == rateLimitModeFixed {
		switch unit {
		case rateLimitUnitMinute:
			t := now.Truncate(time.Minute)
			return t.Add(-time.Duration(w-1) * time.Minute)
		case rateLimitUnitDay:
			y, m, d := now.Date()
			t := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
			return t.AddDate(0, 0, -(w - 1))
		default: // hour
			t := now.Truncate(time.Hour)
			return t.Add(-time.Duration(w-1) * time.Hour)
	placeholder
placeholder
	// rolling window
	switch unit {
	case rateLimitUnitMinute:
		return now.Add(-time.Duration(w) * time.Minute)
	case rateLimitUnitDay:
		return now.AddDate(0, 0, -w)
	default: // hour
		return now.Add(-time.Duration(w) * time.Hour)
placeholder
placeholder

func (s *PaymentService) CancelOrder(ctx context.Context, orderID, userID int64) (string, error) {
	o, err := s.entClient.PaymentOrder.Get(ctx, orderID)
	if err != nil {
		return "", infraerrors.NotFound("NOT_FOUND", "order not found")
placeholder
	if o.UserID != userID {
		return "", infraerrors.Forbidden("FORBIDDEN", "no permission for this order")
placeholder
	if o.Status != OrderStatusPending {
		return "", infraerrors.BadRequest("INVALID_STATUS", "order cannot be cancelled in current status")
placeholder
	return s.cancelCore(ctx, o, OrderStatusCancelled, fmt.Sprintf("user:%d", userID), "user cancelled order")
placeholder

func (s *PaymentService) AdminCancelOrder(ctx context.Context, orderID int64) (string, error) {
	o, err := s.entClient.PaymentOrder.Get(ctx, orderID)
	if err != nil {
		return "", infraerrors.NotFound("NOT_FOUND", "order not found")
placeholder
	if o.Status != OrderStatusPending {
		return "", infraerrors.BadRequest("INVALID_STATUS", "order cannot be cancelled in current status")
placeholder
	return s.cancelCore(ctx, o, OrderStatusCancelled, "admin", "admin cancelled order")
placeholder

func (s *PaymentService) cancelCore(ctx context.Context, o *dbent.PaymentOrder, fs, op, ad string) (string, error) {
	if o.PaymentTradeNo != "" || o.PaymentType != "" {
		if s.checkPaid(ctx, o) == checkPaidResultAlreadyPaid {
			return checkPaidResultAlreadyPaid, nil
	placeholder
placeholder
	c, err := s.entClient.PaymentOrder.Update().Where(paymentorder.IDEQ(o.ID), paymentorder.StatusEQ(OrderStatusPending)).SetStatus(fs).Save(ctx)
	if err != nil {
		return "", fmt.Errorf("update order status: %w", err)
placeholder
	if c > 0 {
		auditAction := "ORDER_CANCELLED"
		if fs == OrderStatusExpired {
			auditAction = "ORDER_EXPIRED"
	placeholder
		s.writeAuditLog(ctx, o.ID, auditAction, op, map[string]any{"detail": adplaceholder)
placeholder
	return checkPaidResultCancelled, nil
placeholder

func (s *PaymentService) checkPaid(ctx context.Context, o *dbent.PaymentOrder) string {
	prov, err := s.getOrderProvider(ctx, o)
	if err != nil {
		return ""
placeholder
	queryRef := paymentOrderQueryReference(o, prov)
	if queryRef == "" {
		return ""
placeholder
	resp, err := prov.QueryOrder(ctx, queryRef)
	if err != nil {
		slog.Warn("query upstream failed", "orderID", o.ID, "error", err)
		return ""
placeholder
	if resp.Status == payment.ProviderStatusPaid {
		if !isValidProviderAmount(resp.Amount) {
			s.writeAuditLog(ctx, o.ID, "PAYMENT_INVALID_AMOUNT", prov.ProviderKey(), map[string]any{
				"expected": o.PayAmount,
				"paid":     resp.Amount,
				"tradeNo":  resp.TradeNo,
				"queryRef": queryRef,
		placeholder)
			slog.Warn("query upstream returned invalid paid amount", "orderID", o.ID, "queryRef", queryRef, "paid", resp.Amount)
			retriedResp, retryOK := requeryPaidOrderOnce(ctx, prov, queryRef)
			if !retryOK {
				return ""
		placeholder
			resp = retriedResp
	placeholder
		notificationTradeNo := o.PaymentTradeNo
		if upstreamTradeNo := strings.TrimSpace(resp.TradeNo); paymentOrderShouldPersistUpstreamTradeNo(queryRef, upstreamTradeNo, notificationTradeNo) {
			if _, updateErr := s.entClient.PaymentOrder.Update().
				Where(paymentorder.IDEQ(o.ID)).
				SetPaymentTradeNo(upstreamTradeNo).
				Save(ctx); updateErr != nil {
				slog.Error("persist upstream trade no during checkPaid failed", "orderID", o.ID, "tradeNo", upstreamTradeNo, "error", updateErr)
		placeholder else {
				o.PaymentTradeNo = upstreamTradeNo
		placeholder
			notificationTradeNo = upstreamTradeNo
	placeholder
		if err := s.HandlePaymentNotification(ctx, &payment.PaymentNotification{TradeNo: notificationTradeNo, OrderID: o.OutTradeNo, Amount: resp.Amount, Status: payment.ProviderStatusSuccess, Metadata: resp.Metadataplaceholder, prov.ProviderKey()); err != nil {
			slog.Error("fulfillment failed during checkPaid", "orderID", o.ID, "error", err)
			// Still return already_paid — order was paid, fulfillment can be retried
	placeholder
		return checkPaidResultAlreadyPaid
placeholder
	if cp, ok := prov.(payment.CancelableProvider); ok {
		_ = cp.CancelPayment(ctx, queryRef)
placeholder
	return ""
placeholder

func requeryPaidOrderOnce(ctx context.Context, prov payment.Provider, queryRef string) (*payment.QueryOrderResponse, bool) {
	if prov == nil || strings.TrimSpace(queryRef) == "" {
		return nil, false
placeholder
	resp, err := prov.QueryOrder(ctx, queryRef)
	if err != nil {
		slog.Warn("query upstream retry failed", "queryRef", queryRef, "error", err)
		return nil, false
placeholder
	if resp == nil || resp.Status != payment.ProviderStatusPaid || !isValidProviderAmount(resp.Amount) {
		return nil, false
placeholder
	return resp, true
placeholder

func paymentOrderQueryReference(order *dbent.PaymentOrder, prov payment.Provider) string {
	if order == nil {
		return ""
placeholder

	providerKey := ""
	if prov != nil {
		providerKey = strings.TrimSpace(prov.ProviderKey())
placeholder
	if providerKey == "" {
		if snapshot := psOrderProviderSnapshot(order); snapshot != nil {
			providerKey = strings.TrimSpace(snapshot.ProviderKey)
	placeholder
placeholder
	if providerKey == "" {
		providerKey = strings.TrimSpace(psStringValue(order.ProviderKey))
placeholder
	if providerKey == "" {
		providerKey = strings.TrimSpace(order.PaymentType)
placeholder

	switch payment.GetBasePaymentType(providerKey) {
	case payment.TypeAlipay, payment.TypeEasyPay, payment.TypeWxpay:
		return strings.TrimSpace(order.OutTradeNo)
	default:
		if tradeNo := strings.TrimSpace(order.PaymentTradeNo); tradeNo != "" {
			return tradeNo
	placeholder
		return strings.TrimSpace(order.OutTradeNo)
placeholder
placeholder

func paymentOrderShouldPersistUpstreamTradeNo(queryRef, upstreamTradeNo, currentTradeNo string) bool {
	upstreamTradeNo = strings.TrimSpace(upstreamTradeNo)
	if upstreamTradeNo == "" {
		return false
placeholder
	if strings.EqualFold(upstreamTradeNo, strings.TrimSpace(currentTradeNo)) {
		return false
placeholder
	if strings.EqualFold(upstreamTradeNo, strings.TrimSpace(queryRef)) {
		return false
placeholder
	return true
placeholder

// VerifyOrderByOutTradeNo actively queries the upstream provider to check
// if a payment was made, and processes it if so. This handles the case where
// the provider's notify callback was missed (e.g. EasyPay popup mode).
func (s *PaymentService) VerifyOrderByOutTradeNo(ctx context.Context, outTradeNo string, userID int64) (*dbent.PaymentOrder, error) {
	outTradeNo, err := normalizeOrderLookupOutTradeNo(outTradeNo)
	if err != nil {
		return nil, err
placeholder
	o, err := s.entClient.PaymentOrder.Query().
		Where(paymentorder.OutTradeNo(outTradeNo)).
		Only(ctx)
	if err != nil {
		return nil, infraerrors.NotFound("NOT_FOUND", "order not found")
placeholder
	if o.UserID != userID {
		return nil, infraerrors.Forbidden("FORBIDDEN", "no permission for this order")
placeholder
	// Only verify orders that are still pending or recently expired
	if o.Status == OrderStatusPending || o.Status == OrderStatusExpired {
		result := s.checkPaid(ctx, o)
		if result == checkPaidResultAlreadyPaid {
			// Reload order to get updated status
			o, err = s.entClient.PaymentOrder.Get(ctx, o.ID)
			if err != nil {
				return nil, fmt.Errorf("reload order: %w", err)
		placeholder
	placeholder
placeholder
	return o, nil
placeholder

// VerifyOrderPublic returns the currently persisted public order state without
// triggering any upstream reconciliation. Signed resume-token recovery is the
// only public recovery path allowed to query upstream state.
func (s *PaymentService) VerifyOrderPublic(ctx context.Context, outTradeNo string) (*dbent.PaymentOrder, error) {
	outTradeNo, err := normalizeOrderLookupOutTradeNo(outTradeNo)
	if err != nil {
		return nil, err
placeholder
	o, err := s.entClient.PaymentOrder.Query().
		Where(paymentorder.OutTradeNo(outTradeNo)).
		Only(ctx)
	if err != nil {
		return nil, infraerrors.NotFound("NOT_FOUND", "order not found")
placeholder
	return o, nil
placeholder

func normalizeOrderLookupOutTradeNo(raw string) (string, error) {
	outTradeNo := strings.TrimSpace(raw)
	if outTradeNo == "" {
		return "", infraerrors.BadRequest("INVALID_OUT_TRADE_NO", "out_trade_no is required")
placeholder
	if len(outTradeNo) > 64 {
		return "", infraerrors.BadRequest("INVALID_OUT_TRADE_NO", "out_trade_no is invalid")
placeholder
	for _, ch := range outTradeNo {
		switch {
		case ch >= 'a' && ch <= 'z':
		case ch >= 'A' && ch <= 'Z':
		case ch >= '0' && ch <= '9':
		case ch == '_' || ch == '-':
		default:
			return "", infraerrors.BadRequest("INVALID_OUT_TRADE_NO", "out_trade_no is invalid")
	placeholder
placeholder
	return outTradeNo, nil
placeholder

func (s *PaymentService) ExpireTimedOutOrders(ctx context.Context) (int, error) {
	now := time.Now()
	orders, err := s.entClient.PaymentOrder.Query().Where(paymentorder.StatusEQ(OrderStatusPending), paymentorder.ExpiresAtLTE(now)).All(ctx)
	if err != nil {
		return 0, fmt.Errorf("query expired: %w", err)
placeholder
	n := 0
	for _, o := range orders {
		// Check upstream payment status before expiring — the user may have
		// paid just before timeout and the webhook hasn't arrived yet.
		outcome, _ := s.cancelCore(ctx, o, OrderStatusExpired, "system", "order expired")
		if outcome == checkPaidResultAlreadyPaid {
			slog.Info("order was paid during expiry", "orderID", o.ID)
			continue
	placeholder
		if outcome != "" {
			n++
	placeholder
placeholder
	return n, nil
placeholder

// getOrderProvider creates a provider using the order's original instance config.
// Falls back to registry lookup if instance ID is missing (legacy orders).
func (s *PaymentService) getOrderProvider(ctx context.Context, o *dbent.PaymentOrder) (payment.Provider, error) {
	inst, err := s.getOrderProviderInstance(ctx, o)
	if err != nil {
		return nil, fmt.Errorf("load order provider instance: %w", err)
placeholder
	if inst != nil {
		return s.createProviderFromInstance(ctx, inst)
placeholder
	if !paymentOrderAllowsRegistryFallback(o) {
		return nil, fmt.Errorf("order %d provider instance is unresolved", o.ID)
placeholder
	providerKey := paymentOrderFallbackProviderKey(s.registry, o)
	if providerKey == "" {
		return nil, fmt.Errorf("order %d provider fallback key is missing", o.ID)
placeholder
	if !s.webhookRegistryFallbackAllowed(ctx, providerKey) {
		return nil, fmt.Errorf("order %d provider fallback is ambiguous for %s", o.ID, providerKey)
placeholder
	s.EnsureProviders(ctx)
	return s.registry.GetProvider(o.PaymentType)
placeholder

func paymentOrderAllowsRegistryFallback(order *dbent.PaymentOrder) bool {
	if order == nil {
		return false
placeholder
	if psOrderProviderSnapshot(order) != nil {
		return false
placeholder
	if strings.TrimSpace(psStringValue(order.ProviderInstanceID)) != "" {
		return false
placeholder
	if strings.TrimSpace(psStringValue(order.ProviderKey)) != "" {
		return false
placeholder
	return true
placeholder

func paymentOrderFallbackProviderKey(registry *payment.Registry, order *dbent.PaymentOrder) string {
	if order == nil {
		return ""
placeholder
	if registry != nil {
		if key := strings.TrimSpace(registry.GetProviderKey(payment.PaymentType(order.PaymentType))); key != "" {
			return key
	placeholder
placeholder
	return strings.TrimSpace(payment.GetBasePaymentType(strings.TrimSpace(order.PaymentType)))
placeholder

func (s *PaymentService) createProviderFromInstance(ctx context.Context, inst *dbent.PaymentProviderInstance) (payment.Provider, error) {
	if inst == nil {
		return nil, fmt.Errorf("payment provider instance is missing")
placeholder

	cfg, err := s.loadBalancer.GetInstanceConfig(ctx, int64(inst.ID))
	if err != nil {
		return nil, fmt.Errorf("load provider instance config: %w", err)
placeholder
	if inst.PaymentMode != "" {
		cfg["paymentMode"] = inst.PaymentMode
placeholder

	instID := strconv.FormatInt(int64(inst.ID), 10)
	prov, err := provider.CreateProvider(inst.ProviderKey, instID, cfg)
	if err != nil {
		return nil, fmt.Errorf("create provider from instance: %w", err)
placeholder
	return prov, nil
placeholder

func psStringValue(value *string) string {
	if value == nil {
		return ""
placeholder
	return *value
placeholder
