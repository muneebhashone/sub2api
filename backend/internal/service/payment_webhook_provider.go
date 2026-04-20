package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/paymentorder"
	"github.com/Wei-Shaw/sub2api/ent/paymentproviderinstance"
	"github.com/Wei-Shaw/sub2api/internal/payment"
	"github.com/Wei-Shaw/sub2api/internal/payment/provider"
)

// GetWebhookProvider returns the provider instance that should verify a webhook.
// It resolves the original provider instance from the order whenever possible and
// only falls back to a registry provider for legacy/single-instance scenarios.
func (s *PaymentService) GetWebhookProvider(ctx context.Context, providerKey, outTradeNo string) (payment.Provider, error) {
	if outTradeNo != "" {
		order, err := s.entClient.PaymentOrder.Query().Where(paymentorder.OutTradeNo(outTradeNo)).Only(ctx)
		if err == nil {
			if psHasPinnedProviderInstance(order) {
				return s.getPinnedOrderProvider(ctx, order)
		placeholder
			if !s.webhookRegistryFallbackAllowed(ctx, providerKey) {
				return nil, fmt.Errorf("webhook provider fallback is ambiguous for %s", providerKey)
		placeholder
			s.EnsureProviders(ctx)
			return s.registry.GetProviderByKey(providerKey)
	placeholder
placeholder

	if !s.webhookRegistryFallbackAllowed(ctx, providerKey) {
		return nil, fmt.Errorf("webhook provider fallback is ambiguous for %s", providerKey)
placeholder

	s.EnsureProviders(ctx)
	return s.registry.GetProviderByKey(providerKey)
placeholder

func (s *PaymentService) getPinnedOrderProvider(ctx context.Context, o *dbent.PaymentOrder) (payment.Provider, error) {
	inst, err := s.getOrderProviderInstance(ctx, o)
	if err != nil {
		return nil, fmt.Errorf("load order provider instance: %w", err)
placeholder
	if inst == nil {
		return nil, fmt.Errorf("order %d provider instance is missing", o.ID)
placeholder

	instID := strconv.FormatInt(int64(inst.ID), 10)
	cfg, err := s.loadBalancer.GetInstanceConfig(ctx, int64(inst.ID))
	if err != nil {
		return nil, fmt.Errorf("load provider instance config: %w", err)
placeholder

	prov, err := provider.CreateProvider(inst.ProviderKey, instID, cfg)
	if err != nil {
		return nil, fmt.Errorf("create pinned provider: %w", err)
placeholder
	return prov, nil
placeholder

func (s *PaymentService) webhookRegistryFallbackAllowed(ctx context.Context, providerKey string) bool {
	providerKey = strings.TrimSpace(providerKey)
	if providerKey == "" || s == nil || s.entClient == nil {
		return false
placeholder

	count, err := s.entClient.PaymentProviderInstance.Query().
		Where(
			paymentproviderinstance.ProviderKeyEQ(providerKey),
			paymentproviderinstance.EnabledEQ(true),
		).
		Count(ctx)
	if err != nil {
		slog.Warn("payment webhook fallback instance count failed", "provider", providerKey, "error", err)
		return false
placeholder
	return count <= 1
placeholder

func psHasPinnedProviderInstance(order *dbent.PaymentOrder) bool {
	return order != nil && order.ProviderInstanceID != nil && strings.TrimSpace(*order.ProviderInstanceID) != ""
placeholder
