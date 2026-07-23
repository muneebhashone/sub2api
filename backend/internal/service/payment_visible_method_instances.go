package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/paymentproviderinstance"
	"github.com/Wei-Shaw/sub2api/internal/payment"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

func enabledVisibleMethodsForProvider(providerKey, supportedTypes string) []string {
	methodSet := make(map[string]struct{placeholder, 2)
	addMethod := func(method string) {
		method = NormalizeVisibleMethod(method)
		if method != "" {
			methodSet[method] = struct{placeholder{placeholder
	placeholder
placeholder

	switch strings.TrimSpace(providerKey) {
	case payment.TypeAlipay:
		if strings.TrimSpace(supportedTypes) == "" {
			addMethod(payment.TypeAlipay)
			break
	placeholder
		for _, supportedType := range splitTypes(supportedTypes) {
			if NormalizeVisibleMethod(supportedType) == payment.TypeAlipay {
				addMethod(payment.TypeAlipay)
				break
		placeholder
	placeholder
	case payment.TypeWxpay:
		if strings.TrimSpace(supportedTypes) == "" {
			addMethod(payment.TypeWxpay)
			break
	placeholder
		for _, supportedType := range splitTypes(supportedTypes) {
			if NormalizeVisibleMethod(supportedType) == payment.TypeWxpay {
				addMethod(payment.TypeWxpay)
				break
		placeholder
	placeholder
	case payment.TypeEasyPay:
		for _, supportedType := range splitTypes(supportedTypes) {
			addMethod(supportedType)
	placeholder
placeholder

	methods := make([]string, 0, len(methodSet))
	for _, method := range []string{payment.TypeAlipay, payment.TypeWxpayplaceholder {
		if _, ok := methodSet[method]; ok {
			methods = append(methods, method)
			delete(methodSet, method)
	placeholder
placeholder
	for _, supportedType := range splitTypes(supportedTypes) {
		method := NormalizeVisibleMethod(supportedType)
		if _, ok := methodSet[method]; ok {
			methods = append(methods, method)
			delete(methodSet, method)
	placeholder
placeholder
	return methods
placeholder

func providerSupportsVisibleMethod(inst *dbent.PaymentProviderInstance, method string) bool {
	if inst == nil || !inst.Enabled {
		return false
placeholder
	method = NormalizeVisibleMethod(method)
	for _, candidate := range enabledVisibleMethodsForProvider(inst.ProviderKey, inst.SupportedTypes) {
		if candidate == method {
			return true
	placeholder
placeholder
	return false
placeholder

func filterEnabledVisibleMethodInstances(instances []*dbent.PaymentProviderInstance, method string) []*dbent.PaymentProviderInstance {
	filtered := make([]*dbent.PaymentProviderInstance, 0, len(instances))
	for _, inst := range instances {
		if providerSupportsVisibleMethod(inst, method) {
			filtered = append(filtered, inst)
	placeholder
placeholder
	return filtered
placeholder

func filterVisibleMethodInstancesByProviderKey(instances []*dbent.PaymentProviderInstance, method string, providerKey string) []*dbent.PaymentProviderInstance {
	filtered := make([]*dbent.PaymentProviderInstance, 0, len(instances))
	for _, inst := range instances {
		if !providerSupportsVisibleMethod(inst, method) {
			continue
	placeholder
		if !strings.EqualFold(strings.TrimSpace(inst.ProviderKey), strings.TrimSpace(providerKey)) {
			continue
	placeholder
		filtered = append(filtered, inst)
placeholder
	return filtered
placeholder

func distinctVisibleMethodProviderKeys(instances []*dbent.PaymentProviderInstance) []string {
	seen := make(map[string]struct{placeholder, len(instances))
	keys := make([]string, 0, len(instances))
	for _, inst := range instances {
		if inst == nil {
			continue
	placeholder
		key := strings.TrimSpace(inst.ProviderKey)
		if key == "" {
			continue
	placeholder
		normalized := strings.ToLower(key)
		if _, ok := seen[normalized]; ok {
			continue
	placeholder
		seen[normalized] = struct{placeholder{placeholder
		keys = append(keys, key)
placeholder
	return keys
placeholder

func selectVisibleMethodInstanceByProviderKey(instances []*dbent.PaymentProviderInstance, providerKey string) *dbent.PaymentProviderInstance {
	providerKey = strings.TrimSpace(providerKey)
	if providerKey == "" {
		return nil
placeholder
	for _, inst := range instances {
		if strings.EqualFold(strings.TrimSpace(inst.ProviderKey), providerKey) {
			return inst
	placeholder
placeholder
	return nil
placeholder

func (s *PaymentConfigService) validateVisibleMethodEnablementConflicts(
	ctx context.Context,
	excludeID int64,
	providerKey string,
	supportedTypes string,
	enabled bool,
) error {
	// Visible methods are selected by configured source (official/easypay),
	// so multiple enabled providers can intentionally claim the same user-facing
	// method. Order creation and limits will route through the configured source.
	_, _, _, _, _ = ctx, excludeID, providerKey, supportedTypes, enabled
	return nil
placeholder

func (s *PaymentConfigService) resolveVisibleMethodSourceProviderKey(ctx context.Context, method string) (string, error) {
	method = NormalizeVisibleMethod(method)
	sourceKey := visibleMethodSourceSettingKey(method)
	rawSource := ""
	if s != nil && s.settingRepo != nil && sourceKey != "" {
		value, err := s.settingRepo.GetValue(ctx, sourceKey)
		if err != nil {
			if !errors.Is(err, ErrSettingNotFound) {
				return "", fmt.Errorf("get %s: %w", sourceKey, err)
		placeholder
	placeholder else {
			rawSource = value
	placeholder
placeholder

	normalizedSource, err := normalizeVisibleMethodSettingSource(method, rawSource, true)
	if err != nil {
		return "", err
placeholder
	if normalizedSource == "" {
		return "", nil
placeholder
	providerKey, ok := VisibleMethodProviderKeyForSource(method, normalizedSource)
	if !ok {
		return "", infraerrors.BadRequest(
			"INVALID_PAYMENT_VISIBLE_METHOD_SOURCE",
			fmt.Sprintf("%s source must be one of the supported payment providers", method),
		)
placeholder
	return providerKey, nil
placeholder

func (s *PaymentConfigService) resolveVisibleMethodProviderKey(
	ctx context.Context,
	method string,
	matching []*dbent.PaymentProviderInstance,
) (string, error) {
	switch providerKeys := distinctVisibleMethodProviderKeys(matching); len(providerKeys) {
	case 0:
		return "", nil
	case 1:
		return strings.TrimSpace(providerKeys[0]), nil
	default:
		providerKey, err := s.resolveVisibleMethodSourceProviderKey(ctx, method)
		if err != nil {
			return "", err
	placeholder
		if providerKey == "" {
			return "", nil
	placeholder
		selected := selectVisibleMethodInstanceByProviderKey(matching, providerKey)
		if selected == nil {
			return "", infraerrors.BadRequest(
				"INVALID_PAYMENT_VISIBLE_METHOD_SOURCE",
				fmt.Sprintf("%s source has no enabled provider instance", method),
			)
	placeholder
		return strings.TrimSpace(selected.ProviderKey), nil
placeholder
placeholder

func (s *PaymentConfigService) resolveEnabledVisibleMethodInstance(
	ctx context.Context,
	method string,
) (*dbent.PaymentProviderInstance, error) {
	if s == nil || s.entClient == nil {
		return nil, nil
placeholder

	method = NormalizeVisibleMethod(method)
	if method == "" {
		return nil, nil
placeholder

	instances, err := s.entClient.PaymentProviderInstance.Query().
		Where(paymentproviderinstance.EnabledEQ(true)).
		Order(paymentproviderinstance.BySortOrder()).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query enabled payment providers: %w", err)
placeholder

	matching := filterEnabledVisibleMethodInstances(instances, method)
	providerKey, err := s.resolveVisibleMethodProviderKey(ctx, method, matching)
	if err != nil {
		return nil, err
placeholder
	if providerKey == "" {
		if len(matching) == 0 {
			return nil, nil
	placeholder
		return &dbent.PaymentProviderInstance{ProviderKey: ""placeholder, nil
placeholder
	return selectVisibleMethodInstanceByProviderKey(matching, providerKey), nil
placeholder

// UsesOfficialAlipayVisibleMethod reports whether the user-facing Alipay method
// currently resolves to an enabled official Alipay provider instance.
func (s *PaymentConfigService) UsesOfficialAlipayVisibleMethod(ctx context.Context) (bool, error) {
	instance, err := s.resolveEnabledVisibleMethodInstance(ctx, payment.TypeAlipay)
	if err != nil {
		return false, err
placeholder
	return isOfficialAlipayProviderInstance(instance), nil
placeholder

func isOfficialAlipayProviderInstance(instance *dbent.PaymentProviderInstance) bool {
	return instance != nil && strings.EqualFold(strings.TrimSpace(instance.ProviderKey), payment.TypeAlipay)
placeholder
