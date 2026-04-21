//go:build unit

package service

import (
	"context"
	"testing"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateProviderRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		providerKey    string
		providerName   string
		supportedTypes string
		wantErr        bool
		errContains    string
placeholder{
		{
			name:           "valid easypay with types",
			providerKey:    "easypay",
			providerName:   "MyProvider",
			supportedTypes: "alipay,wxpay",
			wantErr:        false,
	placeholder,
		{
			name:           "valid stripe with empty types",
			providerKey:    "stripe",
			providerName:   "Stripe Provider",
			supportedTypes: "",
			wantErr:        false,
	placeholder,
		{
			name:           "valid alipay provider",
			providerKey:    "alipay",
			providerName:   "Alipay Direct",
			supportedTypes: "alipay",
			wantErr:        false,
	placeholder,
		{
			name:           "valid wxpay provider",
			providerKey:    "wxpay",
			providerName:   "WeChat Pay",
			supportedTypes: "wxpay",
			wantErr:        false,
	placeholder,
		{
			name:           "invalid provider key",
			providerKey:    "invalid",
			providerName:   "Name",
			supportedTypes: "alipay",
			wantErr:        true,
			errContains:    "invalid provider key",
	placeholder,
		{
			name:           "empty name",
			providerKey:    "easypay",
			providerName:   "",
			supportedTypes: "alipay",
			wantErr:        true,
			errContains:    "provider name is required",
	placeholder,
		{
			name:           "whitespace-only name",
			providerKey:    "easypay",
			providerName:   "  ",
			supportedTypes: "alipay",
			wantErr:        true,
			errContains:    "provider name is required",
	placeholder,
		{
			name:           "tab-only name",
			providerKey:    "easypay",
			providerName:   "\t",
			supportedTypes: "alipay",
			wantErr:        true,
			errContains:    "provider name is required",
	placeholder,
placeholder

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := validateProviderRequest(tc.providerKey, tc.providerName, tc.supportedTypes)
			if tc.wantErr {
			placeholder
				assert.Contains(t, err.Error(), tc.errContains)
		placeholder else {
			placeholder
		placeholder
	placeholder)
placeholder
placeholder

func TestIsSensitiveConfigField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		field   string
		wantSen bool
placeholder{
		// Sensitive fields (contain key/secret/private/password/pkey patterns)
		{"secretKey", trueplaceholder,
		{"apiSecret", trueplaceholder,
		{"pkey", trueplaceholder,
		{"privateKey", trueplaceholder,
		{"apiPassword", trueplaceholder,
		{"appKey", trueplaceholder,
		{"SECRET_TOKEN", trueplaceholder,
		{"PrivateData", trueplaceholder,
		{"PASSWORD", trueplaceholder,
		{"mySecretValue", trueplaceholder,

		// Non-sensitive fields
		{"appId", falseplaceholder,
		{"mchId", falseplaceholder,
		{"apiBase", falseplaceholder,
		{"endpoint", falseplaceholder,
		{"merchantNo", falseplaceholder,
		{"paymentMode", falseplaceholder,
		{"notifyUrl", falseplaceholder,
placeholder

	for _, tc := range tests {
		t.Run(tc.field, func(t *testing.T) {
			t.Parallel()

			got := isSensitiveConfigField(tc.field)
			assert.Equal(t, tc.wantSen, got, "isSensitiveConfigField(%q)", tc.field)
	placeholder)
placeholder
placeholder

func TestJoinTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []string
		want  string
placeholder{
		{
			name:  "multiple types",
			input: []string{"alipay", "wxpay"placeholder,
			want:  "alipay,wxpay",
	placeholder,
		{
			name:  "single type",
			input: []string{"stripe"placeholder,
			want:  "stripe",
	placeholder,
		{
			name:  "empty slice",
			input: []string{placeholder,
			want:  "",
	placeholder,
		{
			name:  "nil slice",
			input: nil,
			want:  "",
	placeholder,
		{
			name:  "three types",
			input: []string{"alipay", "wxpay", "stripe"placeholder,
			want:  "alipay,wxpay,stripe",
	placeholder,
		{
			name:  "types with spaces are not trimmed",
			input: []string{" alipay ", " wxpay "placeholder,
			want:  " alipay , wxpay ",
	placeholder,
placeholder

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := joinTypes(tc.input)
			assert.Equal(t, tc.want, got)
	placeholder)
placeholder
placeholder

func TestCreateProviderInstanceRejectsConflictingVisibleMethodEnablement(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	svc := &PaymentConfigService{
		entClient:     client,
		encryptionKey: []byte("placeholder"),
placeholder

	_, err := svc.CreateProviderInstance(ctx, CreateProviderInstanceRequest{
		ProviderKey:    "easypay",
		Name:           "EasyPay Alipay",
		Config:         map[string]string{"pid": "1001"placeholder,
		SupportedTypes: []string{"alipay"placeholder,
		Enabled:        true,
placeholder)
placeholder

	_, err = svc.CreateProviderInstance(ctx, CreateProviderInstanceRequest{
		ProviderKey:    "alipay",
		Name:           "Official Alipay",
		Config:         map[string]string{"appId": "app-1"placeholder,
		SupportedTypes: []string{"alipay"placeholder,
		Enabled:        true,
placeholder)
placeholder
	require.Equal(t, "PAYMENT_PROVIDER_CONFLICT", infraerrors.Reason(err))
placeholder

func TestUpdateProviderInstanceRejectsEnablingConflictingVisibleMethodProvider(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	svc := &PaymentConfigService{
		entClient:     client,
		encryptionKey: []byte("placeholder"),
placeholder

	existing, err := svc.CreateProviderInstance(ctx, CreateProviderInstanceRequest{
		ProviderKey:    "easypay",
		Name:           "EasyPay WeChat",
		Config:         map[string]string{"pid": "2001"placeholder,
		SupportedTypes: []string{"wxpay"placeholder,
		Enabled:        true,
placeholder)
placeholder
	require.NotNil(t, existing)

	candidate, err := svc.CreateProviderInstance(ctx, CreateProviderInstanceRequest{
		ProviderKey:    "wxpay",
		Name:           "Official WeChat",
		Config:         map[string]string{"appId": "wx-app"placeholder,
		SupportedTypes: []string{"wxpay"placeholder,
		Enabled:        false,
placeholder)
placeholder

	_, err = svc.UpdateProviderInstance(ctx, candidate.ID, UpdateProviderInstanceRequest{
		Enabled: boolPtrValue(true),
placeholder)
placeholder
	require.Equal(t, "PAYMENT_PROVIDER_CONFLICT", infraerrors.Reason(err))
placeholder

func TestUpdateProviderInstancePersistsEnabledAndSupportedTypes(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	svc := &PaymentConfigService{
		entClient:     client,
		encryptionKey: []byte("placeholder"),
placeholder

	instance, err := svc.CreateProviderInstance(ctx, CreateProviderInstanceRequest{
		ProviderKey:    "easypay",
		Name:           "EasyPay",
		Config:         map[string]string{"pid": "3001"placeholder,
		SupportedTypes: []string{"alipay"placeholder,
		Enabled:        false,
placeholder)
placeholder

	_, err = svc.UpdateProviderInstance(ctx, instance.ID, UpdateProviderInstanceRequest{
		Enabled:        boolPtrValue(true),
		SupportedTypes: []string{"alipay", "wxpay"placeholder,
placeholder)
placeholder

	saved, err := client.PaymentProviderInstance.Get(ctx, instance.ID)
placeholder
	require.True(t, saved.Enabled)
	require.Equal(t, "alipay,wxpay", saved.SupportedTypes)
placeholder

func boolPtrValue(v bool) *bool {
	return &v
placeholder
