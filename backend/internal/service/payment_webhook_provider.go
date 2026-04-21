package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/paymentorder"
	"github.com/Wei-Shaw/sub2api/ent/paymentproviderinstance"
	"github.com/Wei-Shaw/sub2api/internal/payment"
)

// GetWebhookProvider returns the provider instance that should verify a webhook.
// It resolves the original provider instance from the order whenever possible and
// only falls back to a registry provider for legacy/single-instance scenarios.
func (s *PaymentService) GetWebhookProvider(ctx context.Context, providerKey, outTradeNo string) (payment.Provider, error) {
	providers, err := s.GetWebhookProviders(ctx, providerKey, outTradeNo)
	if err != nil {
		return nil, err
placeholder
	if len(providers) == 0 {
		return nil, payment.ErrProviderNotFound
placeholder
	return providers[0], nil
placeholder

// GetWebhookProviders returns provider candidates that can verify the webhook.
// Official WeChat Pay may require multiple candidates because the callback body
// cannot be bound to a merchant before decryption.
func (s *PaymentService) GetWebhookProviders(ctx context.Context, providerKey, outTradeNo string) ([]payment.Provider, error) {
	if outTradeNo != "" {
		order, err := s.entClient.PaymentOrder.Query().Where(paymentorder.OutTradeNo(outTradeNo)).Only(ctx)
		if err == nil {
			if psHasPinnedProviderInstance(order) {
				prov, err := s.getPinnedOrderProvider(ctx, order)
				if err != nil {
					return nil, err
			placeholder
				return []payment.Provider{provplaceholder, nil
		placeholder
			inst, err := s.getOrderProviderInstance(ctx, order)
			if err != nil {
				return nil, fmt.Errorf("load order provider instance: %w", err)
		placeholder
			if inst != nil {
				prov, err := s.createProviderFromInstance(ctx, inst)
				if err != nil {
					return nil, err
			placeholder
				return []payment.Provider{provplaceholder, nil
		placeholder
			if strings.TrimSpace(providerKey) == payment.TypeWxpay {
				return s.getEnabledWebhookProvidersByKey(ctx, providerKey)
		placeholder
			if !s.webhookRegistryFallbackAllowed(ctx, providerKey) {
				return nil, fmt.Errorf("webhook provider fallback is ambiguous for %s", providerKey)
		placeholder
			s.EnsureProviders(ctx)
			prov, err := s.registry.GetProviderByKey(providerKey)
			if err != nil {
				return nil, err
		placeholder
			return []payment.Provider{provplaceholder, nil
	placeholder
placeholder

	if strings.TrimSpace(providerKey) == payment.TypeWxpay {
		return s.getEnabledWebhookProvidersByKey(ctx, providerKey)
placeholder

	if !s.webhookRegistryFallbackAllowed(ctx, providerKey) {
		return nil, fmt.Errorf("webhook provider fallback is ambiguous for %s", providerKey)
placeholder

	s.EnsureProviders(ctx)
	prov, err := s.registry.GetProviderByKey(providerKey)
	if err != nil {
		return nil, err
placeholder
	return []payment.Provider{provplaceholder, nil
placeholder

func (s *PaymentService) getPinnedOrderProvider(ctx context.Context, o *dbent.PaymentOrder) (payment.Provider, error) {
	inst, err := s.getOrderProviderInstance(ctx, o)
	if err != nil {
		return nil, fmt.Errorf("load order provider instance: %w", err)
placeholder
	if inst == nil {
		return nil, fmt.Errorf("order %d provider instance is missing", o.ID)
placeholder
	return s.createProviderFromInstance(ctx, inst)
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
	return order != nil && (psOrderProviderSnapshot(order) != nil || (order.ProviderInstanceID != nil && strings.TrimSpace(*order.ProviderInstanceID) != ""))
placeholder

func (s *PaymentService) getEnabledWebhookProvidersByKey(ctx context.Context, providerKey string) ([]payment.Provider, error) {
	providerKey = strings.TrimSpace(providerKey)
	instances, err := s.entClient.PaymentProviderInstance.Query().
		Where(
			paymentproviderinstance.ProviderKeyEQ(providerKey),
			paymentproviderinstance.EnabledEQ(true),
		).
		Order(dbent.Asc(paymentproviderinstance.FieldSortOrder)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query webhook provider instances: %w", err)
placeholder
	if len(instances) == 0 {
		return nil, payment.ErrProviderNotFound
placeholder

	providers := make([]payment.Provider, 0, len(instances))
	for _, inst := range instances {
		prov, provErr := s.createProviderFromInstance(ctx, inst)
		if provErr != nil {
			slog.Warn("skip webhook provider instance", "provider", providerKey, "instanceID", inst.ID, "error", provErr)
			continue
	placeholder
		providers = append(providers, prov)
placeholder
	if len(providers) == 0 {
		return nil, payment.ErrProviderNotFound
placeholder
	return providers, nil
placeholder
