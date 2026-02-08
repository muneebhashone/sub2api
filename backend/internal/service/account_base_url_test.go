//go:build unit

package service

import (
	"testing"
)

func TestGetBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		account  Account
		expected string
placeholder{
		{
			name: "non-apikey type returns empty",
			account: Account{
				Type:     AccountTypeOAuth,
				Platform: PlatformAnthropic,
		placeholder,
			expected: "",
	placeholder,
		{
			name: "apikey without base_url returns default anthropic",
			account: Account{
				Type:        AccountTypeAPIKey,
				Platform:    PlatformAnthropic,
		placeholderplaceholder,
		placeholder,
			expected: "https://api.anthropic.com",
	placeholder,
		{
			name: "apikey with custom base_url",
			account: Account{
				Type:        AccountTypeAPIKey,
				Platform:    PlatformAnthropic,
		placeholder"base_url": "https://custom.example.com"placeholder,
		placeholder,
			expected: "https://custom.example.com",
	placeholder,
		{
			name: "antigravity apikey auto-appends /antigravity",
			account: Account{
				Type:        AccountTypeAPIKey,
				Platform:    PlatformAntigravity,
		placeholder"base_url": "https://upstream.example.com"placeholder,
		placeholder,
			expected: "https://upstream.example.com/antigravity",
	placeholder,
		{
			name: "antigravity apikey trims trailing slash before appending",
			account: Account{
				Type:        AccountTypeAPIKey,
				Platform:    PlatformAntigravity,
		placeholder"base_url": "https://upstream.example.com/"placeholder,
		placeholder,
			expected: "https://upstream.example.com/antigravity",
	placeholder,
		{
			name: "antigravity non-apikey returns empty",
			account: Account{
				Type:        AccountTypeOAuth,
				Platform:    PlatformAntigravity,
		placeholder"base_url": "https://upstream.example.com"placeholder,
		placeholder,
			expected: "",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.account.GetBaseURL()
			if result != tt.expected {
				t.Errorf("GetBaseURL() = %q, want %q", result, tt.expected)
		placeholder
	placeholder)
placeholder
placeholder

func TestGetGeminiBaseURL(t *testing.T) {
	const defaultGeminiURL = "https://generativelanguage.googleapis.com"

	tests := []struct {
		name     string
		account  Account
		expected string
placeholder{
		{
			name: "apikey without base_url returns default",
			account: Account{
				Type:        AccountTypeAPIKey,
				Platform:    PlatformGemini,
		placeholderplaceholder,
		placeholder,
			expected: defaultGeminiURL,
	placeholder,
		{
			name: "apikey with custom base_url",
			account: Account{
				Type:        AccountTypeAPIKey,
				Platform:    PlatformGemini,
		placeholder"base_url": "https://custom-gemini.example.com"placeholder,
		placeholder,
			expected: "https://custom-gemini.example.com",
	placeholder,
		{
			name: "antigravity apikey auto-appends /antigravity",
			account: Account{
				Type:        AccountTypeAPIKey,
				Platform:    PlatformAntigravity,
		placeholder"base_url": "https://upstream.example.com"placeholder,
		placeholder,
			expected: "https://upstream.example.com/antigravity",
	placeholder,
		{
			name: "antigravity apikey trims trailing slash",
			account: Account{
				Type:        AccountTypeAPIKey,
				Platform:    PlatformAntigravity,
		placeholder"base_url": "https://upstream.example.com/"placeholder,
		placeholder,
			expected: "https://upstream.example.com/antigravity",
	placeholder,
		{
			name: "antigravity oauth does NOT append /antigravity",
			account: Account{
				Type:        AccountTypeOAuth,
				Platform:    PlatformAntigravity,
		placeholder"base_url": "https://upstream.example.com"placeholder,
		placeholder,
			expected: "https://upstream.example.com",
	placeholder,
		{
			name: "oauth without base_url returns default",
			account: Account{
				Type:        AccountTypeOAuth,
				Platform:    PlatformAntigravity,
		placeholderplaceholder,
		placeholder,
			expected: defaultGeminiURL,
	placeholder,
		{
			name: "nil credentials returns default",
			account: Account{
				Type:     AccountTypeAPIKey,
		placeholder
		placeholder,
			expected: defaultGeminiURL,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.account.GetGeminiBaseURL(defaultGeminiURL)
			if result != tt.expected {
				t.Errorf("GetGeminiBaseURL() = %q, want %q", result, tt.expected)
		placeholder
	placeholder)
placeholder
placeholder
