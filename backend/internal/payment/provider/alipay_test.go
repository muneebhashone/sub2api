//go:build unit

package provider

import (
	"errors"
	"strings"
	"testing"
)

func TestIsTradeNotExist(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
placeholder{
		{
			name: "nil error returns false",
			err:  nil,
			want: false,
	placeholder,
		{
			name: "error containing ACQ.TRADE_NOT_EXIST returns true",
			err:  errors.New("alipay: sub_code=ACQ.TRADE_NOT_EXIST, sub_msg=交易不存在"),
			want: true,
	placeholder,
		{
			name: "error not containing the code returns false",
			err:  errors.New("alipay: sub_code=ACQ.SYSTEM_ERROR, sub_msg=系统错误"),
			want: false,
	placeholder,
		{
			name: "error with only partial match returns false",
			err:  errors.New("ACQ.TRADE_NOT"),
			want: false,
	placeholder,
		{
			name: "error with exact constant value returns true",
			err:  errors.New(alipayErrTradeNotExist),
			want: true,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := isTradeNotExist(tt.err)
			if got != tt.want {
				t.Errorf("isTradeNotExist(%v) = %v, want %v", tt.err, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestNewAlipay(t *testing.T) {
	t.Parallel()

	validConfig := map[string]string{
		"appId":      "2021001234567890",
		"privateKey": "placeholder",
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
			name:      "missing privateKey",
			config:    withOverride(map[string]string{"privateKey": ""placeholder),
			wantErr:   true,
			errSubstr: "privateKey",
	placeholder,
		{
			name:    "nil config map returns error for appId",
			config:  map[string]string{placeholder,
			wantErr: true,
			errSubstr: "appId",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewAlipay("test-instance", tt.config)
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
				t.Fatal("expected non-nil Alipay instance")
		placeholder
			if got.instanceID != "test-instance" {
				t.Errorf("instanceID = %q, want %q", got.instanceID, "test-instance")
		placeholder
	placeholder)
placeholder
placeholder
