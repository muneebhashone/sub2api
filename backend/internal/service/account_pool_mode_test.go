//go:build unit

package service

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPoolModeRetryCount(t *testing.T) {
	tests := []struct {
		name     string
		account  *Account
		expected int
placeholder{
		{
			name: "default_when_not_pool_mode",
			account: &Account{
				Type:        AccountTypeAPIKey,
				Platform:    PlatformOpenAI,
		placeholderplaceholder,
		placeholder,
			expected: defaultPoolModeRetryCount,
	placeholder,
		{
			name: "default_when_missing_retry_count",
			account: &Account{
				Type:     AccountTypeAPIKey,
				Platform: PlatformOpenAI,
		placeholder
					"pool_mode": true,
			placeholder,
		placeholder,
			expected: defaultPoolModeRetryCount,
	placeholder,
		{
			name: "supports_float64_from_json_credentials",
			account: &Account{
				Type:     AccountTypeAPIKey,
				Platform: PlatformOpenAI,
		placeholder
					"pool_mode":             true,
					"pool_mode_retry_count": float64(5),
			placeholder,
		placeholder,
			expected: 5,
	placeholder,
		{
			name: "supports_json_number",
			account: &Account{
				Type:     AccountTypeAPIKey,
				Platform: PlatformOpenAI,
		placeholder
					"pool_mode":             true,
					"pool_mode_retry_count": json.Number("4"),
			placeholder,
		placeholder,
			expected: 4,
	placeholder,
		{
			name: "supports_string_value",
			account: &Account{
				Type:     AccountTypeAPIKey,
				Platform: PlatformOpenAI,
		placeholder
					"pool_mode":             true,
					"pool_mode_retry_count": "2",
			placeholder,
		placeholder,
			expected: 2,
	placeholder,
		{
			name: "negative_value_is_clamped_to_zero",
			account: &Account{
				Type:     AccountTypeAPIKey,
				Platform: PlatformOpenAI,
		placeholder
					"pool_mode":             true,
					"pool_mode_retry_count": -1,
			placeholder,
		placeholder,
			expected: 0,
	placeholder,
		{
			name: "oversized_value_is_clamped_to_max",
			account: &Account{
				Type:     AccountTypeAPIKey,
				Platform: PlatformOpenAI,
		placeholder
					"pool_mode":             true,
					"pool_mode_retry_count": 99,
			placeholder,
		placeholder,
			expected: maxPoolModeRetryCount,
	placeholder,
		{
			name: "invalid_value_falls_back_to_default",
			account: &Account{
				Type:     AccountTypeAPIKey,
				Platform: PlatformOpenAI,
		placeholder
					"pool_mode":             true,
					"pool_mode_retry_count": "oops",
			placeholder,
		placeholder,
			expected: defaultPoolModeRetryCount,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.account.GetPoolModeRetryCount())
	placeholder)
placeholder
placeholder
