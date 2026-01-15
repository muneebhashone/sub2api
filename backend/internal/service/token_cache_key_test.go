//go:build unit

package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGeminiTokenCacheKey(t *testing.T) {
	tests := []struct {
		name     string
		account  *Account
		expected string
placeholder{
		{
			name: "with_project_id",
			account: &Account{
				ID: 100,
		placeholder
					"project_id": "my-project-123",
			placeholder,
		placeholder,
			expected: "my-project-123",
	placeholder,
		{
			name: "project_id_with_whitespace",
			account: &Account{
				ID: 101,
		placeholder
					"project_id": "  project-with-spaces  ",
			placeholder,
		placeholder,
			expected: "project-with-spaces",
	placeholder,
		{
			name: "empty_project_id_fallback_to_account_id",
			account: &Account{
				ID: 102,
		placeholder
					"project_id": "",
			placeholder,
		placeholder,
			expected: "account:102",
	placeholder,
		{
			name: "whitespace_only_project_id_fallback_to_account_id",
			account: &Account{
				ID: 103,
		placeholder
					"project_id": "   ",
			placeholder,
		placeholder,
			expected: "account:103",
	placeholder,
		{
			name: "no_project_id_key_fallback_to_account_id",
			account: &Account{
				ID:          104,
		placeholderplaceholder,
		placeholder,
			expected: "account:104",
	placeholder,
		{
			name: "nil_credentials_fallback_to_account_id",
			account: &Account{
				ID:          105,
				Credentials: nil,
		placeholder,
			expected: "account:105",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GeminiTokenCacheKey(tt.account)
			require.Equal(t, tt.expected, result)
	placeholder)
placeholder
placeholder

func TestAntigravityTokenCacheKey(t *testing.T) {
	tests := []struct {
		name     string
		account  *Account
		expected string
placeholder{
		{
			name: "with_project_id",
			account: &Account{
				ID: 200,
		placeholder
					"project_id": "ag-project-456",
			placeholder,
		placeholder,
			expected: "ag:ag-project-456",
	placeholder,
		{
			name: "project_id_with_whitespace",
			account: &Account{
				ID: 201,
		placeholder
					"project_id": "  ag-project-spaces  ",
			placeholder,
		placeholder,
			expected: "ag:ag-project-spaces",
	placeholder,
		{
			name: "empty_project_id_fallback_to_account_id",
			account: &Account{
				ID: 202,
		placeholder
					"project_id": "",
			placeholder,
		placeholder,
			expected: "ag:account:202",
	placeholder,
		{
			name: "whitespace_only_project_id_fallback_to_account_id",
			account: &Account{
				ID: 203,
		placeholder
					"project_id": "   ",
			placeholder,
		placeholder,
			expected: "ag:account:203",
	placeholder,
		{
			name: "no_project_id_key_fallback_to_account_id",
			account: &Account{
				ID:          204,
		placeholderplaceholder,
		placeholder,
			expected: "ag:account:204",
	placeholder,
		{
			name: "nil_credentials_fallback_to_account_id",
			account: &Account{
				ID:          205,
				Credentials: nil,
		placeholder,
			expected: "ag:account:205",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AntigravityTokenCacheKey(tt.account)
			require.Equal(t, tt.expected, result)
	placeholder)
placeholder
placeholder
