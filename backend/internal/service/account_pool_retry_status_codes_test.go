//go:build unit

package service

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPoolModeRetryStatusCodes(t *testing.T) {
	tests := []struct {
		name     string
		account  *Account
		expected []int
placeholder{
		{
			name:     "nil_account_returns_nil",
			account:  nil,
			expected: nil,
	placeholder,
		{
			name: "nil_credentials_returns_nil",
			account: &Account{
				Type:     AccountTypeAPIKey,
				Platform: PlatformOpenAI,
		placeholder,
			expected: nil,
	placeholder,
		{
			name: "missing_key_returns_nil",
			account: &Account{
				Type:        AccountTypeAPIKey,
				Platform:    PlatformOpenAI,
		placeholder"pool_mode": trueplaceholder,
		placeholder,
			expected: nil,
	placeholder,
		{
			name: "empty_slice_is_preserved",
			account: &Account{
		placeholder
					"pool_mode_retry_status_codes": []any{placeholder,
			placeholder,
		placeholder,
			expected: []int{placeholder,
	placeholder,
		{
			name: "float64_values_from_json_are_normalized",
			account: &Account{
		placeholder
					"pool_mode_retry_status_codes": []any{float64(429), float64(401), float64(403)placeholder,
			placeholder,
		placeholder,
			expected: []int{401, 403, 429placeholder,
	placeholder,
		{
			name: "json_number_values_supported",
			account: &Account{
		placeholder
					"pool_mode_retry_status_codes": []any{json.Number("502"), json.Number("503")placeholder,
			placeholder,
		placeholder,
			expected: []int{502, 503placeholder,
	placeholder,
		{
			name: "string_values_supported",
			account: &Account{
		placeholder
					"pool_mode_retry_status_codes": []any{"520", "529"placeholder,
			placeholder,
		placeholder,
			expected: []int{520, 529placeholder,
	placeholder,
		{
			name: "duplicates_are_deduped",
			account: &Account{
		placeholder
					"pool_mode_retry_status_codes": []any{float64(429), float64(429), float64(401)placeholder,
			placeholder,
		placeholder,
			expected: []int{401, 429placeholder,
	placeholder,
		{
			name: "out_of_range_values_dropped",
			account: &Account{
		placeholder
					"pool_mode_retry_status_codes": []any{float64(99), float64(600), float64(429)placeholder,
			placeholder,
		placeholder,
			expected: []int{429placeholder,
	placeholder,
		{
			name: "invalid_string_dropped",
			account: &Account{
		placeholder
					"pool_mode_retry_status_codes": []any{"oops", float64(429)placeholder,
			placeholder,
		placeholder,
			expected: []int{429placeholder,
	placeholder,
		{
			name: "non_array_value_returns_nil",
			account: &Account{
		placeholder
					"pool_mode_retry_status_codes": "not-an-array",
			placeholder,
		placeholder,
			expected: nil,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.account.GetPoolModeRetryStatusCodes())
	placeholder)
placeholder
placeholder

func TestIsPoolModeRetryableStatus_Account(t *testing.T) {
	tests := []struct {
		name       string
		account    *Account
		statusCode int
		expected   bool
placeholder{
		{
			name:       "nil_account_falls_back_to_default_401",
			account:    nil,
			statusCode: 401,
			expected:   true,
	placeholder,
		{
			name:       "nil_account_falls_back_to_default_500",
			account:    nil,
			statusCode: 500,
			expected:   false,
	placeholder,
		{
			name: "unconfigured_uses_default_403",
			account: &Account{
		placeholder"pool_mode": trueplaceholder,
		placeholder,
			statusCode: 403,
			expected:   true,
	placeholder,
		{
			name: "unconfigured_uses_default_502_false",
			account: &Account{
		placeholder"pool_mode": trueplaceholder,
		placeholder,
			statusCode: 502,
			expected:   false,
	placeholder,
		{
			name: "configured_list_overrides_default_401_dropped",
			account: &Account{
		placeholder
					"pool_mode_retry_status_codes": []any{float64(502), float64(503)placeholder,
			placeholder,
		placeholder,
			statusCode: 401,
			expected:   false,
	placeholder,
		{
			name: "configured_list_overrides_default_502_added",
			account: &Account{
		placeholder
					"pool_mode_retry_status_codes": []any{float64(502), float64(503)placeholder,
			placeholder,
		placeholder,
			statusCode: 502,
			expected:   true,
	placeholder,
		{
			name: "empty_list_disables_all_default_codes",
			account: &Account{
		placeholder
					"pool_mode_retry_status_codes": []any{placeholder,
			placeholder,
		placeholder,
			statusCode: 429,
			expected:   false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.account.IsPoolModeRetryableStatus(tt.statusCode))
	placeholder)
placeholder
placeholder
