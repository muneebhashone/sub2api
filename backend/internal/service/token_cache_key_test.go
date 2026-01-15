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

func TestOpenAITokenCacheKey(t *testing.T) {
	tests := []struct {
		name     string
		account  *Account
		expected string
placeholder{
		{
			name: "basic_account",
			account: &Account{
				ID: 300,
		placeholder,
			expected: "openai:account:300",
	placeholder,
		{
			name: "account_with_credentials",
			account: &Account{
				ID: 301,
		placeholder
					"access_token": "test-token",
			placeholder,
		placeholder,
			expected: "openai:account:301",
	placeholder,
		{
			name: "account_id_zero",
			account: &Account{
				ID: 0,
		placeholder,
			expected: "openai:account:0",
	placeholder,
		{
			name: "large_account_id",
			account: &Account{
				ID: 9999999999,
		placeholder,
			expected: "openai:account:9999999999",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := OpenAITokenCacheKey(tt.account)
			require.Equal(t, tt.expected, result)
	placeholder)
placeholder
placeholder

func TestClaudeTokenCacheKey(t *testing.T) {
	tests := []struct {
		name     string
		account  *Account
		expected string
placeholder{
		{
			name: "basic_account",
			account: &Account{
				ID: 400,
		placeholder,
			expected: "claude:account:400",
	placeholder,
		{
			name: "account_with_credentials",
			account: &Account{
				ID: 401,
		placeholder
					"access_token": "claude-token",
			placeholder,
		placeholder,
			expected: "claude:account:401",
	placeholder,
		{
			name: "account_id_zero",
			account: &Account{
				ID: 0,
		placeholder,
			expected: "claude:account:0",
	placeholder,
		{
			name: "large_account_id",
			account: &Account{
				ID: 9999999999,
		placeholder,
			expected: "claude:account:9999999999",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClaudeTokenCacheKey(tt.account)
			require.Equal(t, tt.expected, result)
	placeholder)
placeholder
placeholder

func TestCacheKeyUniqueness(t *testing.T) {
	// 确保不同平台的缓存键不会冲突
	account := &Account{ID: placeholder

	openaiKey := OpenAITokenCacheKey(account)
	claudeKey := ClaudeTokenCacheKey(account)

	require.NotEqual(t, openaiKey, claudeKey, "OpenAI and Claude cache keys should be different")
	require.Contains(t, openaiKey, "openai:")
	require.Contains(t, claudeKey, "claude:")
placeholder
