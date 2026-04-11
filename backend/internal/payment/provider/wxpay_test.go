//go:build unit

package provider

import (
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/payment"
)

func TestMapWxState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
placeholder{
		{
			name:  "SUCCESS maps to paid",
			input: wxpayTradeStateSuccess,
			want:  payment.ProviderStatusPaid,
	placeholder,
		{
			name:  "REFUND maps to refunded",
			input: wxpayTradeStateRefund,
			want:  payment.ProviderStatusRefunded,
	placeholder,
		{
			name:  "CLOSED maps to failed",
			input: wxpayTradeStateClosed,
			want:  payment.ProviderStatusFailed,
	placeholder,
		{
			name:  "PAYERROR maps to failed",
			input: wxpayTradeStatePayError,
			want:  payment.ProviderStatusFailed,
	placeholder,
		{
			name:  "unknown state maps to pending",
			input: "NOTPAY",
			want:  payment.ProviderStatusPending,
	placeholder,
		{
			name:  "empty string maps to pending",
			input: "",
			want:  payment.ProviderStatusPending,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := mapWxState(tt.input)
			if got != tt.want {
				t.Errorf("mapWxState(%q) = %q, want %q", tt.input, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestWxSV(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *string
		want  string
placeholder{
		{
			name:  "nil pointer returns empty string",
			input: nil,
			want:  "",
	placeholder,
		{
			name:  "non-nil pointer returns value",
			input: strPtr("hello"),
			want:  "hello",
	placeholder,
		{
			name:  "pointer to empty string returns empty string",
			input: strPtr(""),
			want:  "",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := wxSV(tt.input)
			if got != tt.want {
				t.Errorf("wxSV() = %q, want %q", got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func strPtr(s string) *string {
	return &s
placeholder

func TestFormatPEM(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     string
		keyType string
		want    string
placeholder{
		{
			name:    "raw key gets wrapped with headers",
			key:     "MIIBIjANBgkqhki...",
			keyType: "PUBLIC KEY",
			want:    "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhki...\n-----END PUBLIC KEY-----",
	placeholder,
		{
			name:    "already formatted key is returned as-is",
			key:     "placeholder
placeholder
placeholder\nMIIEvQIBADANBg...\n-----END PRIVATE KEY-----",
	placeholder,
		{
			name:    "key with leading/trailing whitespace is trimmed before check",
			key:     "  \n MIIBIjANBgkqhki...  \n ",
			keyType: "PUBLIC KEY",
			want:    "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhki...\n-----END PUBLIC KEY-----",
	placeholder,
		{
			name:    "already formatted key with whitespace is trimmed and returned",
			key:     "  placeholder
placeholder
placeholder\ndata\n-----END RSA PRIVATE KEY-----",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := formatPEM(tt.key, tt.keyType)
			if got != tt.want {
				t.Errorf("formatPEM(%q, %q) =\n%s\nwant:\n%s", tt.key, tt.keyType, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestNewWxpay(t *testing.T) {
	t.Parallel()

	validConfig := map[string]string{
		"appId":       "wx1234567890",
		"mchId":       "1234567890",
		"privateKey":  "fake-private-key",
		"apiV3Key":    "12345678901234567890123456789012", // exactly 32 bytes
		"publicKey":   "fake-public-key",
		"publicKeyId": "key-id-001",
		"certSerial":  "SERIAL001",
placeholder

	// helper to clone and override config fields
	withOverride := func(overrides map[string]string) map[string]string {
		cfg := make(map[string]string, len(validConfig))
		for k, v := range validConfig {
			cfg[k] = v
	placeholder
		for k, v := range overrides {
			cfg[k] = v
	placeholder
		return cfg
placeholder

	tests := []struct {
		name      string
		config    map[string]string
		wantErr   bool
		errSubstr string
placeholder{
		{
			name:    "valid config succeeds",
			config:  validConfig,
			wantErr: false,
	placeholder,
		{
			name:      "missing appId",
			config:    withOverride(map[string]string{"appId": ""placeholder),
			wantErr:   true,
			errSubstr: "appId",
	placeholder,
		{
			name:      "missing mchId",
			config:    withOverride(map[string]string{"mchId": ""placeholder),
			wantErr:   true,
			errSubstr: "mchId",
	placeholder,
		{
			name:      "missing privateKey",
			config:    withOverride(map[string]string{"privateKey": ""placeholder),
			wantErr:   true,
			errSubstr: "privateKey",
	placeholder,
		{
			name:      "missing apiV3Key",
			config:    withOverride(map[string]string{"apiV3Key": ""placeholder),
			wantErr:   true,
			errSubstr: "apiV3Key",
	placeholder,
		{
			name:      "missing publicKey",
			config:    withOverride(map[string]string{"publicKey": ""placeholder),
			wantErr:   true,
			errSubstr: "publicKey",
	placeholder,
		{
			name:      "missing publicKeyId",
			config:    withOverride(map[string]string{"publicKeyId": ""placeholder),
			wantErr:   true,
			errSubstr: "publicKeyId",
	placeholder,
		{
			name:      "apiV3Key too short",
			config:    withOverride(map[string]string{"apiV3Key": "short"placeholder),
			wantErr:   true,
			errSubstr: "exactly 32 bytes",
	placeholder,
		{
			name:      "apiV3Key too long",
			config:    withOverride(map[string]string{"apiV3Key": "123456789012345678901234567890123"placeholder), // 33 bytes
			wantErr:   true,
			errSubstr: "exactly 32 bytes",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewWxpay("test-instance", tt.config)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
			placeholder
				if tt.errSubstr != "" && !strings.Contains(err.Error(), tt.errSubstr) {
					t.Errorf("error %q should contain %q", err.Error(), tt.errSubstr)
			placeholder
				return
		placeholder
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
		placeholder
			if got == nil {
				t.Fatal("expected non-nil Wxpay instance")
		placeholder
			if got.instanceID != "test-instance" {
				t.Errorf("instanceID = %q, want %q", got.instanceID, "test-instance")
		placeholder
	placeholder)
placeholder
placeholder
