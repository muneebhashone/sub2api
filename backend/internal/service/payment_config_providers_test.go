//go:build unit

package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

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

func TestIsSensitiveProviderConfigField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		providerKey string
		field       string
		wantSen     bool
placeholder{
		// Stripe: publishableKey is public, only secretKey/webhookSecret are secrets
		{"stripe", "secretKey", trueplaceholder,
		{"stripe", "webhookSecret", trueplaceholder,
		{"stripe", "SecretKey", trueplaceholder, // case-insensitive
		{"stripe", "publishableKey", falseplaceholder,
		{"stripe", "appId", falseplaceholder,

		// Alipay
		{"alipay", "privateKey", trueplaceholder,
		{"alipay", "publicKey", trueplaceholder,
		{"alipay", "alipayPublicKey", trueplaceholder,
		{"alipay", "appId", falseplaceholder,
		{"alipay", "notifyUrl", falseplaceholder,

		// Wxpay
		{"wxpay", "privateKey", trueplaceholder,
		{"wxpay", "apiV3Key", trueplaceholder,
		{"wxpay", "publicKey", trueplaceholder,
		{"wxpay", "publicKeyId", falseplaceholder,
		{"wxpay", "certSerial", falseplaceholder,
		{"wxpay", "mchId", falseplaceholder,

		// EasyPay
		{"easypay", "pkey", trueplaceholder,
		{"easypay", "pid", falseplaceholder,
		{"easypay", "apiBase", falseplaceholder,

		// Unknown provider: never sensitive
		{"unknown", "secretKey", falseplaceholder,
placeholder

	for _, tc := range tests {
		tc := tc
		t.Run(tc.providerKey+"/"+tc.field, func(t *testing.T) {
			t.Parallel()

			got := isSensitiveProviderConfigField(tc.providerKey, tc.field)
			assert.Equal(t, tc.wantSen, got, "isSensitiveProviderConfigField(%q, %q)", tc.providerKey, tc.field)
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

func TestCreateProviderInstanceAllowsVisibleMethodProvidersFromDifferentSources(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	svc := &PaymentConfigService{
		entClient:     client,
		encryptionKey: []byte("placeholder"),
placeholder

	_, err := svc.CreateProviderInstance(ctx, CreateProviderInstanceRequest{
		ProviderKey: "easypay",
		Name:        "EasyPay Alipay",
		Config: map[string]string{
			"pid":       "1001",
			"pkey":      "pkey-1001",
			"apiBase":   "https://pay.example.com",
			"notifyUrl": "https://merchant.example.com/notify",
			"returnUrl": "https://merchant.example.com/return",
	placeholder,
		SupportedTypes: []string{"alipay"placeholder,
		Enabled:        true,
placeholder)
placeholder

	_, err = svc.CreateProviderInstance(ctx, CreateProviderInstanceRequest{
		ProviderKey:    "alipay",
		Name:           "Official Alipay",
		Config:         map[string]string{"appId": "app-1", "privateKey": "private-key"placeholder,
		SupportedTypes: []string{"alipay"placeholder,
		Enabled:        true,
placeholder)
placeholder
placeholder

func TestUpdateProviderInstanceAllowsEnablingVisibleMethodProviderFromDifferentSource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	svc := &PaymentConfigService{
		entClient:     client,
		encryptionKey: []byte("placeholder"),
placeholder

	existing, err := svc.CreateProviderInstance(ctx, CreateProviderInstanceRequest{
		ProviderKey: "easypay",
		Name:        "EasyPay WeChat",
		Config: map[string]string{
			"pid":       "2001",
			"pkey":      "pkey-2001",
			"apiBase":   "https://pay.example.com",
			"notifyUrl": "https://merchant.example.com/notify",
			"returnUrl": "https://merchant.example.com/return",
	placeholder,
		SupportedTypes: []string{"wxpay"placeholder,
		Enabled:        true,
placeholder)
placeholder
	require.NotNil(t, existing)

	candidate, err := svc.CreateProviderInstance(ctx, CreateProviderInstanceRequest{
		ProviderKey:    "wxpay",
		Name:           "Official WeChat",
		Config:         validWxpayProviderConfig(t),
		SupportedTypes: []string{"wxpay"placeholder,
		Enabled:        false,
placeholder)
placeholder

	_, err = svc.UpdateProviderInstance(ctx, candidate.ID, UpdateProviderInstanceRequest{
		Enabled: boolPtrValue(true),
placeholder)
placeholder
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
		ProviderKey: "easypay",
		Name:        "EasyPay",
		Config: map[string]string{
			"pid":       "3001",
			"pkey":      "pkey-3001",
			"apiBase":   "https://pay.example.com",
			"notifyUrl": "https://merchant.example.com/notify",
			"returnUrl": "https://merchant.example.com/return",
	placeholder,
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

func validWxpayProviderConfig(t *testing.T) map[string]string {
placeholder

	key, err := rsa.GenerateKey(rand.Reader, 2048)
placeholder

	privDER, err := x509.MarshalPKCS8PrivateKey(key)
placeholder
	pubDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
placeholder

	return map[string]string{
		"appId":       "wx-app-test",
		"mchId":       "mch-test",
		"privateKey":  string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDERplaceholder)),
		"apiV3Key":    "12345678901234567890123456789012",
		"publicKey":   string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDERplaceholder)),
		"publicKeyId": "public-key-id-test",
		"certSerial":  "cert-serial-test",
placeholder
placeholder
