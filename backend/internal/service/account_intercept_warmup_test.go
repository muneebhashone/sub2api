//go:build unit

package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccount_IsInterceptWarmupEnabled(t *testing.T) {
	tests := []struct {
		name        string
		credentials map[string]any
		expected    bool
placeholder{
		{
			name:        "nil credentials",
			credentials: nil,
			expected:    false,
	placeholder,
		{
			name:        "empty map",
			credentials: map[string]any{placeholder,
			expected:    false,
	placeholder,
		{
			name:        "field not present",
			credentials: map[string]any{"access_token": "tok"placeholder,
			expected:    false,
	placeholder,
		{
			name:        "field is true",
			credentials: map[string]any{"intercept_warmup_requests": trueplaceholder,
			expected:    true,
	placeholder,
		{
			name:        "field is false",
			credentials: map[string]any{"intercept_warmup_requests": falseplaceholder,
			expected:    false,
	placeholder,
		{
			name:        "field is string true",
			credentials: map[string]any{"intercept_warmup_requests": "true"placeholder,
			expected:    false,
	placeholder,
		{
			name:        "field is int 1",
			credentials: map[string]any{"intercept_warmup_requests": 1placeholder,
			expected:    false,
	placeholder,
		{
			name:        "field is nil",
			credentials: map[string]any{"intercept_warmup_requests": nilplaceholder,
			expected:    false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{Credentials: tt.credentialsplaceholder
			result := a.IsInterceptWarmupEnabled()
			require.Equal(t, tt.expected, result)
	placeholder)
placeholder
placeholder
