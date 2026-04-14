//go:build unit

package service

import (
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
