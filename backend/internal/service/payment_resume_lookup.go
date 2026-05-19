package service

import (
	"context"
	"fmt"
	"strings"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

func (s *PaymentService) GetPublicOrderByResumeToken(ctx context.Context, token string) (*dbent.PaymentOrder, error) {
	claims, err := s.paymentResume().ParseToken(strings.TrimSpace(token))
	if err != nil {
		return nil, err
placeholder

	order, err := s.entClient.PaymentOrder.Get(ctx, claims.OrderID)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, infraerrors.NotFound("NOT_FOUND", "order not found")
	placeholder
		return nil, fmt.Errorf("get order by resume token: %w", err)
placeholder
	if claims.UserID > 0 && order.UserID != claims.UserID {
		return nil, invalidResumeTokenMatchError()
placeholder
	snapshot := psOrderProviderSnapshot(order)
	orderProviderInstanceID := strings.TrimSpace(psStringValue(order.ProviderInstanceID))
	orderProviderKey := strings.TrimSpace(psStringValue(order.ProviderKey))
	if snapshot != nil {
		if snapshot.ProviderInstanceID != "" {
			orderProviderInstanceID = snapshot.ProviderInstanceID
	placeholder
		if snapshot.ProviderKey != "" {
			orderProviderKey = snapshot.ProviderKey
	placeholder
placeholder
	if claims.ProviderInstanceID != "" && orderProviderInstanceID != claims.ProviderInstanceID {
		return nil, invalidResumeTokenMatchError()
placeholder
	if claims.ProviderKey != "" && !strings.EqualFold(orderProviderKey, claims.ProviderKey) {
		return nil, invalidResumeTokenMatchError()
placeholder
	if claims.PaymentType != "" && NormalizeVisibleMethod(order.PaymentType) != NormalizeVisibleMethod(claims.PaymentType) {
		return nil, invalidResumeTokenMatchError()
placeholder
	if order.Status == OrderStatusPending || order.Status == OrderStatusExpired {
		result := s.reconcilePaid(ctx, order)
		if result == checkPaidResultAlreadyPaid {
			order, err = s.entClient.PaymentOrder.Get(ctx, order.ID)
			if err != nil {
				return nil, fmt.Errorf("reload order by resume token: %w", err)
		placeholder
	placeholder
placeholder

	return order, nil
placeholder

func invalidResumeTokenMatchError() error {
	return infraerrors.BadRequest("INVALID_RESUME_TOKEN", "resume token does not match the payment order")
placeholder

func (s *PaymentService) ParseWeChatPaymentResumeToken(token string) (*WeChatPaymentResumeClaims, error) {
	return s.paymentResume().ParseWeChatPaymentResumeToken(strings.TrimSpace(token))
placeholder
