package service

import (
	"context"
	"encoding/json"
	"fmt"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/paymentproviderinstance"
	"github.com/Wei-Shaw/sub2api/internal/payment"
)

// GetAvailableMethodLimits collects all payment types from enabled provider
// instances and returns limits for each, plus the global widest range.
// Stripe sub-types (card, link) are aggregated under "stripe".
func (s *PaymentConfigService) GetAvailableMethodLimits(ctx context.Context) (*MethodLimitsResponse, error) {
	instances, err := s.entClient.PaymentProviderInstance.Query().
		Where(paymentproviderinstance.EnabledEQ(true)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query provider instances: %w", err)
placeholder
	typeInstances := pcGroupByPaymentType(instances)
	resp := &MethodLimitsResponse{
		Methods: make(map[string]MethodLimits, len(typeInstances)),
placeholder
	for pt, insts := range typeInstances {
		ml := pcAggregateMethodLimits(pt, insts)
		resp.Methods[ml.PaymentType] = ml
placeholder
	resp.GlobalMin, resp.GlobalMax = pcComputeGlobalRange(resp.Methods)
	return resp, nil
placeholder

// GetMethodLimits returns per-payment-type limits from enabled provider instances.
func (s *PaymentConfigService) GetMethodLimits(ctx context.Context, types []string) ([]MethodLimits, error) {
	instances, err := s.entClient.PaymentProviderInstance.Query().
		Where(paymentproviderinstance.EnabledEQ(true)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query provider instances: %w", err)
placeholder
	result := make([]MethodLimits, 0, len(types))
	for _, pt := range types {
		var matching []*dbent.PaymentProviderInstance
		for _, inst := range instances {
			if payment.InstanceSupportsType(inst.SupportedTypes, pt) {
				matching = append(matching, inst)
		placeholder
	placeholder
		result = append(result, pcAggregateMethodLimits(pt, matching))
placeholder
	return result, nil
placeholder

// pcGroupByPaymentType groups instances by user-facing payment type.
// For Stripe providers, ALL sub-types (card, link, alipay, wxpay) map to "stripe"
// because the user sees a single "Stripe" button, not individual sub-methods.
// Uses a seen set to avoid counting one instance twice.
func pcGroupByPaymentType(instances []*dbent.PaymentProviderInstance) map[string][]*dbent.PaymentProviderInstance {
	typeInstances := make(map[string][]*dbent.PaymentProviderInstance)
	seen := make(map[string]map[int64]bool)
	add := func(key string, inst *dbent.PaymentProviderInstance) {
		if seen[key] == nil {
			seen[key] = make(map[int64]bool)
	placeholder
		if !seen[key][int64(inst.ID)] {
			seen[key][int64(inst.ID)] = true
			typeInstances[key] = append(typeInstances[key], inst)
	placeholder
placeholder
	for _, inst := range instances {
		// Stripe provider: all sub-types → single "stripe" group
		if inst.ProviderKey == payment.TypeStripe {
			add(payment.TypeStripe, inst)
			continue
	placeholder
		for _, t := range splitTypes(inst.SupportedTypes) {
			add(t, inst)
	placeholder
placeholder
	return typeInstances
placeholder

// pcInstanceTypeLimits extracts per-type limits from a provider instance.
// Returns (limits, true) if configured; (zero, false) if unlimited.
// For Stripe instances, limits are stored under "stripe" key regardless of sub-types.
func pcInstanceTypeLimits(inst *dbent.PaymentProviderInstance, pt string) (payment.ChannelLimits, bool) {
	if inst.Limits == "" {
		return payment.ChannelLimits{placeholder, false
placeholder
	var limits payment.InstanceLimits
	if err := json.Unmarshal([]byte(inst.Limits), &limits); err != nil {
		return payment.ChannelLimits{placeholder, false
placeholder
	cl, ok := limits[pt]
	return cl, ok
placeholder

// unionFloat merges a single limit value into the aggregate using UNION semantics.
//   - For "min" fields (wantMin=true): keeps the lowest non-zero value
//   - For "max"/"cap" fields (wantMin=false): keeps the highest non-zero value
//   - If any value is 0 (unlimited), the result is unlimited.
//
// Returns (aggregated value, still limited).
func unionFloat(agg float64, limited bool, val float64, wantMin bool) (float64, bool) {
	if val == 0 {
		return agg, false
placeholder
	if !limited {
		return agg, false
placeholder
	if agg == 0 {
		return val, true
placeholder
	if wantMin && val < agg {
		return val, true
placeholder
	if !wantMin && val > agg {
		return val, true
placeholder
	return agg, true
placeholder

// pcAggregateMethodLimits computes the UNION (least restrictive) of limits
// across all provider instances for a given payment type.
//
// Since the load balancer can route an order to any available instance,
// the user should see the widest possible range:
//   - SingleMin: lowest floor across instances; 0 if any is unlimited
//   - SingleMax: highest ceiling across instances; 0 if any is unlimited
//   - DailyLimit: highest cap across instances; 0 if any is unlimited
func pcAggregateMethodLimits(pt string, instances []*dbent.PaymentProviderInstance) MethodLimits {
	ml := MethodLimits{PaymentType: ptplaceholder
	minLimited, maxLimited, dailyLimited := true, true, true

	for _, inst := range instances {
		cl, hasLimits := pcInstanceTypeLimits(inst, pt)
		if !hasLimits {
			return MethodLimits{PaymentType: ptplaceholder // any unlimited instance → all zeros
	placeholder
		ml.SingleMin, minLimited = unionFloat(ml.SingleMin, minLimited, cl.SingleMin, true)
		ml.SingleMax, maxLimited = unionFloat(ml.SingleMax, maxLimited, cl.SingleMax, false)
		ml.DailyLimit, dailyLimited = unionFloat(ml.DailyLimit, dailyLimited, cl.DailyLimit, false)
placeholder

	if !minLimited {
		ml.SingleMin = 0
placeholder
	if !maxLimited {
		ml.SingleMax = 0
placeholder
	if !dailyLimited {
		ml.DailyLimit = 0
placeholder
	return ml
placeholder

// pcComputeGlobalRange computes the widest [min, max] across all methods.
// Uses the same union logic: lowest min, highest max, 0 if any is unlimited.
func pcComputeGlobalRange(methods map[string]MethodLimits) (globalMin, globalMax float64) {
	minLimited, maxLimited := true, true
	for _, ml := range methods {
		globalMin, minLimited = unionFloat(globalMin, minLimited, ml.SingleMin, true)
		globalMax, maxLimited = unionFloat(globalMax, maxLimited, ml.SingleMax, false)
placeholder
	if !minLimited {
		globalMin = 0
placeholder
	if !maxLimited {
		globalMax = 0
placeholder
	return globalMin, globalMax
placeholder
