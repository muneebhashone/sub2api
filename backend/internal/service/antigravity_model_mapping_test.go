//go:build unit

package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsAntigravityModelSupported(t *testing.T) {
	tests := []struct {
		name     string
		model    string
		expected bool
placeholder{
		// 直接支持的模型
		{"直接支持 - claude-sonnet-4-5", "claude-sonnet-4-5", trueplaceholder,
		{"直接支持 - claude-opus-4-5-thinking", "claude-opus-4-5-thinking", trueplaceholder,
		{"直接支持 - claude-sonnet-4-5-thinking", "claude-sonnet-4-5-thinking", trueplaceholder,
		{"直接支持 - gemini-2.5-flash", "gemini-2.5-flash", trueplaceholder,
		{"直接支持 - gemini-2.5-flash-lite", "gemini-2.5-flash-lite", trueplaceholder,
		{"直接支持 - gemini-3-pro-high", "gemini-3-pro-high", trueplaceholder,

		// 可映射的模型
		{"可映射 - claude-3-5-sonnet-20241022", "claude-3-5-sonnet-20241022", trueplaceholder,
		{"可映射 - claude-3-5-sonnet-20240620", "claude-3-5-sonnet-20240620", trueplaceholder,
		{"可映射 - claude-opus-4", "claude-opus-4", trueplaceholder,
		{"可映射 - claude-haiku-4", "claude-haiku-4", trueplaceholder,
		{"可映射 - claude-3-haiku-20240307", "claude-3-haiku-20240307", trueplaceholder,

		// Gemini 前缀透传
		{"Gemini前缀 - gemini-1.5-pro", "gemini-1.5-pro", trueplaceholder,
		{"Gemini前缀 - gemini-unknown-model", "gemini-unknown-model", trueplaceholder,
		{"Gemini前缀 - gemini-future-version", "gemini-future-version", trueplaceholder,

		// Claude 前缀兜底
		{"Claude前缀 - claude-unknown-model", "claude-unknown-model", trueplaceholder,
		{"Claude前缀 - claude-3-opus-20240229", "claude-3-opus-20240229", trueplaceholder,
		{"Claude前缀 - claude-future-version", "claude-future-version", trueplaceholder,

		// 不支持的模型
		{"不支持 - gpt-4", "gpt-4", falseplaceholder,
		{"不支持 - gpt-4o", "gpt-4o", falseplaceholder,
		{"不支持 - llama-3", "llama-3", falseplaceholder,
		{"不支持 - mistral-7b", "mistral-7b", falseplaceholder,
		{"不支持 - 空字符串", "", falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAntigravityModelSupported(tt.model)
			require.Equal(t, tt.expected, got, "model: %s", tt.model)
	placeholder)
placeholder
placeholder

func TestAntigravityGatewayService_GetMappedModel(t *testing.T) {
	svc := &AntigravityGatewayService{placeholder

	tests := []struct {
		name           string
		requestedModel string
		accountMapping map[string]string
		expected       string
placeholder{
		// 1. 账户级映射优先（注意：model_mapping 在 credentials 中存储为 map[string]any）
		{
			name:           "账户映射优先",
			requestedModel: "claude-3-5-sonnet-20241022",
			accountMapping: map[string]string{"claude-3-5-sonnet-20241022": "custom-model"placeholder,
			expected:       "custom-model",
	placeholder,
		{
			name:           "账户映射覆盖系统映射",
			requestedModel: "claude-opus-4",
			accountMapping: map[string]string{"claude-opus-4": "my-opus"placeholder,
			expected:       "my-opus",
	placeholder,

		// 2. 系统默认映射
		{
			name:           "系统映射 - claude-3-5-sonnet-20241022",
			requestedModel: "claude-3-5-sonnet-20241022",
			accountMapping: nil,
			expected:       "claude-sonnet-4-5",
	placeholder,
		{
			name:           "系统映射 - claude-3-5-sonnet-20240620",
			requestedModel: "claude-3-5-sonnet-20240620",
			accountMapping: nil,
			expected:       "claude-sonnet-4-5",
	placeholder,
		{
			name:           "系统映射 - claude-opus-4",
			requestedModel: "claude-opus-4",
			accountMapping: nil,
			expected:       "claude-opus-4-5-thinking",
	placeholder,
		{
			name:           "系统映射 - claude-opus-4-5-20251101",
			requestedModel: "claude-opus-4-5-20251101",
			accountMapping: nil,
			expected:       "claude-opus-4-5-thinking",
	placeholder,
		{
			name:           "系统映射 - claude-haiku-4 → gemini-3-flash",
			requestedModel: "claude-haiku-4",
			accountMapping: nil,
			expected:       "gemini-3-flash",
	placeholder,
		{
			name:           "系统映射 - claude-haiku-4-5 → gemini-3-flash",
			requestedModel: "claude-haiku-4-5",
			accountMapping: nil,
			expected:       "gemini-3-flash",
	placeholder,
		{
			name:           "系统映射 - claude-3-haiku-20240307 → gemini-3-flash",
			requestedModel: "claude-3-haiku-20240307",
			accountMapping: nil,
			expected:       "gemini-3-flash",
	placeholder,
		{
			name:           "系统映射 - placeholder → gemini-3-flash",
			requestedModel: "placeholder",
			accountMapping: nil,
			expected:       "gemini-3-flash",
	placeholder,
		{
			name:           "系统映射 - claude-sonnet-4-5-20250929",
			requestedModel: "claude-sonnet-4-5-20250929",
			accountMapping: nil,
			expected:       "claude-sonnet-4-5",
	placeholder,

		// 3. Gemini 透传
		{
			name:           "Gemini透传 - gemini-2.5-flash",
			requestedModel: "gemini-2.5-flash",
			accountMapping: nil,
			expected:       "gemini-2.5-flash",
	placeholder,
		{
			name:           "Gemini透传 - gemini-1.5-pro",
			requestedModel: "gemini-1.5-pro",
			accountMapping: nil,
			expected:       "gemini-1.5-pro",
	placeholder,
		{
			name:           "Gemini透传 - gemini-future-model",
			requestedModel: "gemini-future-model",
			accountMapping: nil,
			expected:       "gemini-future-model",
	placeholder,

		// 4. 直接支持的模型
		{
			name:           "直接支持 - claude-sonnet-4-5",
			requestedModel: "claude-sonnet-4-5",
			accountMapping: nil,
			expected:       "claude-sonnet-4-5",
	placeholder,
		{
			name:           "直接支持 - claude-opus-4-5-thinking",
			requestedModel: "claude-opus-4-5-thinking",
			accountMapping: nil,
			expected:       "claude-opus-4-5-thinking",
	placeholder,
		{
			name:           "直接支持 - claude-sonnet-4-5-thinking",
			requestedModel: "claude-sonnet-4-5-thinking",
			accountMapping: nil,
			expected:       "claude-sonnet-4-5-thinking",
	placeholder,

		// 5. 默认值 fallback（未知 claude 模型）
		{
			name:           "默认值 - claude-unknown",
			requestedModel: "claude-unknown",
			accountMapping: nil,
			expected:       "claude-sonnet-4-5",
	placeholder,
		{
			name:           "默认值 - claude-3-opus-20240229",
			requestedModel: "claude-3-opus-20240229",
			accountMapping: nil,
			expected:       "claude-sonnet-4-5",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := &Account{
				Platform: PlatformAntigravity,
		placeholder
			if tt.accountMapping != nil {
				// GetModelMapping 期望 model_mapping 是 map[string]any 格式
				mappingAny := make(map[string]any)
				for k, v := range tt.accountMapping {
					mappingAny[k] = v
			placeholder
				account.Credentials = map[string]any{
					"model_mapping": mappingAny,
			placeholder
		placeholder

			got := svc.getMappedModel(account, tt.requestedModel)
			require.Equal(t, tt.expected, got, "model: %s", tt.requestedModel)
	placeholder)
placeholder
placeholder

func TestAntigravityGatewayService_GetMappedModel_EdgeCases(t *testing.T) {
	svc := &AntigravityGatewayService{placeholder

	tests := []struct {
		name           string
		requestedModel string
		expected       string
placeholder{
		// 空字符串回退到默认值
		{"空字符串", "", "claude-sonnet-4-5"placeholder,

		// 非 claude/gemini 前缀回退到默认值
		{"非claude/gemini前缀 - gpt", "gpt-4", "claude-sonnet-4-5"placeholder,
		{"非claude/gemini前缀 - llama", "llama-3", "claude-sonnet-4-5"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := &Account{Platform: PlatformAntigravityplaceholder
			got := svc.getMappedModel(account, tt.requestedModel)
			require.Equal(t, tt.expected, got)
	placeholder)
placeholder
placeholder

func TestAntigravityGatewayService_IsModelSupported(t *testing.T) {
	svc := &AntigravityGatewayService{placeholder

	tests := []struct {
		name     string
		model    string
		expected bool
placeholder{
		// 直接支持
		{"直接支持 - claude-sonnet-4-5", "claude-sonnet-4-5", trueplaceholder,
		{"直接支持 - gemini-3-flash", "gemini-3-flash", trueplaceholder,

		// 可映射
		{"可映射 - claude-opus-4", "claude-opus-4", trueplaceholder,

		// 前缀透传
		{"Gemini前缀", "gemini-unknown", trueplaceholder,
		{"Claude前缀", "claude-unknown", trueplaceholder,

		// 不支持
		{"不支持 - gpt-4", "gpt-4", falseplaceholder,
		{"不支持 - 空字符串", "", falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.IsModelSupported(tt.model)
			require.Equal(t, tt.expected, got)
	placeholder)
placeholder
placeholder
