//go:build unit

package service

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/stretchr/testify/require"
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

func TestGetGrokBaseURLUsesSubscriptionProxyForOAuth(t *testing.T) {
	tests := []struct {
		name     string
		account  Account
		expected string
placeholder{
		{
			name: "oauth without base_url uses CLI subscription proxy",
			account: Account{
				Type:        AccountTypeOAuth,
				Platform:    PlatformGrok,
		placeholderplaceholder,
		placeholder,
			expected: xai.DefaultCLIBaseURL,
	placeholder,
		{
			name: "oauth legacy API default is migrated at runtime to CLI subscription proxy",
			account: Account{
				Type:     AccountTypeOAuth,
				Platform: PlatformGrok,
		placeholder
					"base_url": xai.DefaultBaseURL,
			placeholder,
		placeholder,
			expected: xai.DefaultCLIBaseURL,
	placeholder,
		{
			name: "oauth legacy API default with trailing slash is migrated at runtime",
			account: Account{
				Type:     AccountTypeOAuth,
				Platform: PlatformGrok,
		placeholder
					"base_url": xai.DefaultBaseURL + "/",
			placeholder,
		placeholder,
			expected: xai.DefaultCLIBaseURL,
	placeholder,
		{
			name: "oauth legacy API root is migrated at runtime",
			account: Account{
				Type:     AccountTypeOAuth,
				Platform: PlatformGrok,
		placeholder
					"base_url": "https://api.x.ai",
			placeholder,
		placeholder,
			expected: xai.DefaultCLIBaseURL,
	placeholder,
		{
			name: "oauth legacy API root with canonical HTTPS port is migrated at runtime",
			account: Account{
				Type:     AccountTypeOAuth,
				Platform: PlatformGrok,
		placeholder
					"base_url": "HTTPS://API.X.AI:443/",
			placeholder,
		placeholder,
			expected: xai.DefaultCLIBaseURL,
	placeholder,
		{
			name: "oauth legacy API canonical port with leading zeroes is migrated at runtime",
			account: Account{
				Type:     AccountTypeOAuth,
				Platform: PlatformGrok,
		placeholder
					"base_url": "https://api.x.ai:0443/v1",
			placeholder,
		placeholder,
			expected: xai.DefaultCLIBaseURL,
	placeholder,
		{
			name: "oauth legacy API encoded version path is migrated at runtime",
			account: Account{
				Type:     AccountTypeOAuth,
				Platform: PlatformGrok,
		placeholder
					"base_url": "https://api.x.ai/%76%31",
			placeholder,
		placeholder,
			expected: xai.DefaultCLIBaseURL,
	placeholder,
		{
			name: "oauth legacy API encoded trailing slash is migrated at runtime",
			account: Account{
				Type:     AccountTypeOAuth,
				Platform: PlatformGrok,
		placeholder
					"base_url": "https://api.x.ai/v1%2F",
			placeholder,
		placeholder,
			expected: xai.DefaultCLIBaseURL,
	placeholder,
		{
			name: "oauth non-default API port remains pinned to CLI proxy",
			account: Account{
				Type:     AccountTypeOAuth,
				Platform: PlatformGrok,
		placeholder
					"base_url": "https://api.x.ai:8443/v1",
			placeholder,
		placeholder,
			expected: xai.DefaultCLIBaseURL,
	placeholder,
		{
			name: "oauth explicit custom base_url redirects forwarding traffic",
			account: Account{
				Type:     AccountTypeOAuth,
				Platform: PlatformGrok,
		placeholder
					"base_url": "https://custom.example.com/v1",
			placeholder,
		placeholder,
			expected: "https://custom.example.com/v1",
	placeholder,
		{
			name: "oauth custom base_url with path prefix redirects forwarding traffic",
			account: Account{
				Type:     AccountTypeOAuth,
				Platform: PlatformGrok,
		placeholder
					"base_url": "https://relay.example.com/xai/v1",
			placeholder,
		placeholder,
			expected: "https://relay.example.com/xai/v1",
	placeholder,
		{
			name: "API key without base_url uses official credit-backed API",
			account: Account{
				Type:        AccountTypeAPIKey,
				Platform:    PlatformGrok,
		placeholderplaceholder,
		placeholder,
			expected: xai.DefaultBaseURL,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.account.GetGrokBaseURL())
	placeholder)
placeholder
placeholder

func TestGetGrokBaseURLHonorsOAuthCustomRegardlessOfUnsafeOverrides(t *testing.T) {
	t.Setenv(xai.EnvAllowUnsafeURLOverrides, "true")
	account := Account{
		Type:     AccountTypeOAuth,
		Platform: PlatformGrok,
placeholder
			"base_url": "https://custom.example.com/v1",
	placeholder,
placeholder

	require.Equal(t, "https://custom.example.com/v1", account.GetGrokBaseURL())
placeholder

func TestGetGrokMediaBaseURLPinsOAuthMediaToCLIProxy(t *testing.T) {
	tests := []struct {
		name     string
		account  Account
		expected string
placeholder{
		{
			name: "oauth without base_url uses CLI subscription proxy",
			account: Account{
				Type:        AccountTypeOAuth,
				Platform:    PlatformGrok,
		placeholderplaceholder,
		placeholder,
			expected: xai.DefaultCLIBaseURL,
	placeholder,
		{
			name: "oauth stored CLI proxy stays on CLI subscription proxy",
			account: Account{
				Type:     AccountTypeOAuth,
				Platform: PlatformGrok,
		placeholder
					"base_url": xai.DefaultCLIBaseURL,
			placeholder,
		placeholder,
			expected: xai.DefaultCLIBaseURL,
	placeholder,
		{
			name: "oauth stored CLI proxy variant is canonicalized to CLI proxy",
			account: Account{
				Type:     AccountTypeOAuth,
				Platform: PlatformGrok,
		placeholder
					"base_url": "HTTPS://CLI-CHAT-PROXY.GROK.COM:443/%76%31/",
			placeholder,
		placeholder,
			expected: xai.DefaultCLIBaseURL,
	placeholder,
		{
			name: "oauth legacy official API is pinned to CLI proxy",
			account: Account{
				Type:     AccountTypeOAuth,
				Platform: PlatformGrok,
		placeholder
					"base_url": xai.DefaultBaseURL,
			placeholder,
		placeholder,
			expected: xai.DefaultCLIBaseURL,
	placeholder,
		{
			name: "oauth custom base_url redirects media traffic",
			account: Account{
				Type:     AccountTypeOAuth,
				Platform: PlatformGrok,
		placeholder
					"base_url": "https://custom.example.com/v1",
			placeholder,
		placeholder,
			expected: "https://custom.example.com/v1",
	placeholder,
		{
			name: "API key retains its configured media API",
			account: Account{
				Type:     AccountTypeAPIKey,
				Platform: PlatformGrok,
		placeholder
					"base_url": "https://grok.example.com/v1",
			placeholder,
		placeholder,
			expected: "https://grok.example.com/v1",
	placeholder,
		{
			name: "non-Grok account has no Grok media base URL",
			account: Account{
				Type:        AccountTypeOAuth,
				Platform:    PlatformOpenAI,
		placeholderplaceholder,
		placeholder,
			expected: "",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.account.GetGrokMediaBaseURL())
	placeholder)
placeholder
placeholder

func TestGetGrokMediaBaseURLHonorsOAuthCustomRegardlessOfUnsafeOverrides(t *testing.T) {
	t.Setenv(xai.EnvAllowUnsafeURLOverrides, "true")
	account := Account{
		Type:     AccountTypeOAuth,
		Platform: PlatformGrok,
placeholder
			"base_url": "https://custom.example.com/v1",
	placeholder,
placeholder

	require.Equal(t, "https://custom.example.com/v1", account.GetGrokMediaBaseURL())
placeholder
