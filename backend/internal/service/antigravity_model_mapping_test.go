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
		// 在默认映射中的模型（支持）
		{"默认映射 - claude-sonnet-4-5", "claude-sonnet-4-5", trueplaceholder,
		{"默认映射 - claude-opus-4-6-thinking", "claude-opus-4-6-thinking", trueplaceholder,
		{"默认映射 - claude-opus-4-6", "claude-opus-4-6", trueplaceholder,
		{"默认映射 - claude-opus-4-5-thinking", "claude-opus-4-5-thinking", trueplaceholder,
		{"默认映射 - claude-sonnet-4-5-thinking", "claude-sonnet-4-5-thinking", trueplaceholder,
		{"默认映射 - gemini-2.5-flash", "gemini-2.5-flash", trueplaceholder,
		{"默认映射 - gemini-2.5-flash-lite", "gemini-2.5-flash-lite", trueplaceholder,
		{"默认映射 - gemini-3-pro-high", "gemini-3-pro-high", trueplaceholder,
		{"默认映射 - claude-haiku-4-5", "claude-haiku-4-5", trueplaceholder,

		// 不在默认映射中的模型（不支持）
		{"未配置 - claude-3-5-sonnet-20241022", "claude-3-5-sonnet-20241022", falseplaceholder,
		{"未配置 - claude-3-5-sonnet-20240620", "claude-3-5-sonnet-20240620", falseplaceholder,
		{"未配置 - claude-3-haiku-20240307", "claude-3-haiku-20240307", falseplaceholder,
		{"未配置 - gemini-unknown-model", "gemini-unknown-model", falseplaceholder,
		{"未配置 - gemini-future-version", "gemini-future-version", falseplaceholder,
		{"未配置 - claude-unknown-model", "claude-unknown-model", falseplaceholder,
		{"未配置 - claude-3-opus-20240229", "claude-3-opus-20240229", falseplaceholder,
		{"未配置 - claude-future-version", "claude-future-version", falseplaceholder,

		// 非 Claude/Gemini 模型（不支持）
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
		// 1. 账户级映射优先
		{
			name:           "账户映射优先",
			requestedModel: "claude-3-5-sonnet-20241022",
			accountMapping: map[string]string{"claude-3-5-sonnet-20241022": "custom-model"placeholder,
			expected:       "custom-model",
	placeholder,
		{
			name:           "账户映射 - 可覆盖默认映射的模型",
			requestedModel: "claude-sonnet-4-5",
			accountMapping: map[string]string{"claude-sonnet-4-5": "my-custom-sonnet"placeholder,
			expected:       "my-custom-sonnet",
	placeholder,
		{
			name:           "账户映射 - 可覆盖未知模型",
			requestedModel: "claude-opus-4",
			accountMapping: map[string]string{"claude-opus-4": "my-opus"placeholder,
			expected:       "my-opus",
	placeholder,

		// 2. 默认映射（DefaultAntigravityModelMapping）
		{
			name:           "默认映射 - claude-opus-4-6 → claude-opus-4-6-thinking",
			requestedModel: "claude-opus-4-6",
			accountMapping: nil,
			expected:       "claude-opus-4-6-thinking",
	placeholder,
		{
			name:           "默认映射 - claude-opus-4-5-20251101 → claude-opus-4-6-thinking",
			requestedModel: "claude-opus-4-5-20251101",
			accountMapping: nil,
			expected:       "claude-opus-4-6-thinking",
	placeholder,
		{
			name:           "默认映射 - claude-opus-4-5-thinking → claude-opus-4-6-thinking",
			requestedModel: "claude-opus-4-5-thinking",
			accountMapping: nil,
			expected:       "claude-opus-4-6-thinking",
	placeholder,
		{
			name:           "默认映射 - claude-haiku-4-5 → claude-sonnet-4-5",
			requestedModel: "claude-haiku-4-5",
			accountMapping: nil,
			expected:       "claude-sonnet-4-5",
	placeholder,
		{
			name:           "默认映射 - placeholder → claude-sonnet-4-5",
			requestedModel: "placeholder",
			accountMapping: nil,
			expected:       "claude-sonnet-4-5",
	placeholder,
		{
			name:           "默认映射 - claude-sonnet-4-5-20250929 → claude-sonnet-4-5",
			requestedModel: "claude-sonnet-4-5-20250929",
			accountMapping: nil,
			expected:       "claude-sonnet-4-5",
	placeholder,

		// 3. 默认映射中的透传（映射到自己）
		{
			name:           "默认映射透传 - claude-sonnet-4-5",
			requestedModel: "claude-sonnet-4-5",
			accountMapping: nil,
			expected:       "claude-sonnet-4-5",
	placeholder,
		{
			name:           "默认映射透传 - claude-opus-4-6-thinking",
			requestedModel: "claude-opus-4-6-thinking",
			accountMapping: nil,
			expected:       "claude-opus-4-6-thinking",
	placeholder,
		{
			name:           "默认映射透传 - claude-sonnet-4-5-thinking",
			requestedModel: "claude-sonnet-4-5-thinking",
			accountMapping: nil,
			expected:       "claude-sonnet-4-5-thinking",
	placeholder,
		{
			name:           "默认映射透传 - gemini-2.5-flash",
			requestedModel: "gemini-2.5-flash",
			accountMapping: nil,
			expected:       "gemini-2.5-flash",
	placeholder,
		{
			name:           "默认映射透传 - gemini-2.5-pro",
			requestedModel: "gemini-2.5-pro",
			accountMapping: nil,
			expected:       "gemini-2.5-pro",
	placeholder,
		{
			name:           "默认映射透传 - gemini-3-flash",
			requestedModel: "gemini-3-flash",
			accountMapping: nil,
			expected:       "gemini-3-flash",
	placeholder,

		// 4. 未在默认映射中的模型返回空字符串（不支持）
		{
			name:           "未知模型 - claude-unknown 返回空",
			requestedModel: "claude-unknown",
			accountMapping: nil,
			expected:       "",
	placeholder,
		{
			name:           "未知模型 - claude-3-5-sonnet-20241022 返回空（未在默认映射）",
			requestedModel: "claude-3-5-sonnet-20241022",
			accountMapping: nil,
			expected:       "",
	placeholder,
		{
			name:           "未知模型 - claude-3-opus-20240229 返回空",
			requestedModel: "claude-3-opus-20240229",
			accountMapping: nil,
			expected:       "",
	placeholder,
		{
			name:           "未知模型 - claude-opus-4 返回空",
			requestedModel: "claude-opus-4",
			accountMapping: nil,
			expected:       "",
	placeholder,
		{
			name:           "未知模型 - gemini-future-model 返回空",
			requestedModel: "gemini-future-model",
			accountMapping: nil,
			expected:       "",
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
		// 空字符串和非 claude/gemini 前缀返回空字符串
		{"空字符串", "", ""placeholder,
		{"非claude/gemini前缀 - gpt", "gpt-4", ""placeholder,
		{"非claude/gemini前缀 - llama", "llama-3", ""placeholder,
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

		// 可映射（有明确前缀映射）
		{"可映射 - claude-opus-4-6", "claude-opus-4-6", trueplaceholder,

		// 前缀透传（claude 和 gemini 前缀）
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
