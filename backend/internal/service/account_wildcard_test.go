//go:build unit

package service

import (
	"testing"
)

func TestMatchWildcard(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		str      string
		expected bool
placeholder{
		// 精确匹配
		{"exact match", "claude-sonnet-4-5", "claude-sonnet-4-5", trueplaceholder,
		{"exact mismatch", "claude-sonnet-4-5", "claude-opus-4-5", falseplaceholder,

		// 通配符匹配
		{"wildcard prefix match", "claude-*", "claude-sonnet-4-5", trueplaceholder,
		{"wildcard prefix match 2", "claude-*", "claude-opus-4-5-thinking", trueplaceholder,
		{"wildcard prefix mismatch", "claude-*", "gemini-3-flash", falseplaceholder,
		{"wildcard partial match", "gemini-3*", "gemini-3-flash", trueplaceholder,
		{"wildcard partial match 2", "gemini-3*", "gemini-3-pro-image", trueplaceholder,
		{"wildcard partial mismatch", "gemini-3*", "gemini-2.5-flash", falseplaceholder,

		// 边界情况
		{"empty pattern exact", "", "", trueplaceholder,
		{"empty pattern mismatch", "", "claude", falseplaceholder,
		{"single star", "*", "anything", trueplaceholder,
		{"star at end only", "abc*", "abcdef", trueplaceholder,
		{"star at end empty suffix", "abc*", "abc", trueplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchWildcard(tt.pattern, tt.str)
			if result != tt.expected {
				t.Errorf("matchWildcard(%q, %q) = %v, want %v", tt.pattern, tt.str, result, tt.expected)
		placeholder
	placeholder)
placeholder
placeholder

func TestMatchWildcardMapping(t *testing.T) {
	tests := []struct {
		name           string
		mapping        map[string]string
		requestedModel string
		expected       string
placeholder{
		// 精确匹配优先于通配符
		{
			name: "exact match takes precedence",
			mapping: map[string]string{
				"claude-sonnet-4-5": "claude-sonnet-4-5-exact",
				"claude-*":          "claude-default",
		placeholder,
			requestedModel: "claude-sonnet-4-5",
			expected:       "claude-sonnet-4-5-exact",
	placeholder,

		// 最长通配符优先
		{
			name: "longer wildcard takes precedence",
			mapping: map[string]string{
				"claude-*":         "claude-default",
				"claude-sonnet-*":  "claude-sonnet-default",
				"claude-sonnet-4*": "claude-sonnet-4-series",
		placeholder,
			requestedModel: "claude-sonnet-4-5",
			expected:       "claude-sonnet-4-series",
	placeholder,

		// 单个通配符
		{
			name: "single wildcard",
			mapping: map[string]string{
				"claude-*": "claude-mapped",
		placeholder,
			requestedModel: "claude-opus-4-5",
			expected:       "claude-mapped",
	placeholder,

		// 无匹配返回原始模型
		{
			name: "no match returns original",
			mapping: map[string]string{
				"claude-*": "claude-mapped",
		placeholder,
			requestedModel: "gemini-3-flash",
			expected:       "gemini-3-flash",
	placeholder,

		// 空映射返回原始模型
		{
			name:           "empty mapping returns original",
			mapping:        map[string]string{placeholder,
			requestedModel: "claude-sonnet-4-5",
			expected:       "claude-sonnet-4-5",
	placeholder,

		// Gemini 模型映射
		{
			name: "gemini wildcard mapping",
			mapping: map[string]string{
				"gemini-3*":   "gemini-3-pro-high",
				"gemini-2.5*": "gemini-2.5-flash",
		placeholder,
			requestedModel: "gemini-3-flash-preview",
			expected:       "gemini-3-pro-high",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchWildcardMapping(tt.mapping, tt.requestedModel)
			if result != tt.expected {
				t.Errorf("matchWildcardMapping(%v, %q) = %q, want %q", tt.mapping, tt.requestedModel, result, tt.expected)
		placeholder
	placeholder)
placeholder
placeholder

func TestAccountIsModelSupported(t *testing.T) {
	tests := []struct {
		name           string
		credentials    map[string]any
		requestedModel string
		expected       bool
placeholder{
		// 无映射 = 允许所有
		{
			name:           "no mapping allows all",
			credentials:    nil,
			requestedModel: "any-model",
			expected:       true,
	placeholder,
		{
			name:           "empty mapping allows all",
			credentials:    map[string]any{placeholder,
			requestedModel: "any-model",
			expected:       true,
	placeholder,

		// 精确匹配
		{
			name: "exact match supported",
			credentials: map[string]any{
				"model_mapping": map[string]any{
					"claude-sonnet-4-5": "target-model",
			placeholder,
		placeholder,
			requestedModel: "claude-sonnet-4-5",
			expected:       true,
	placeholder,
		{
			name: "exact match not supported",
			credentials: map[string]any{
				"model_mapping": map[string]any{
					"claude-sonnet-4-5": "target-model",
			placeholder,
		placeholder,
			requestedModel: "claude-opus-4-5",
			expected:       false,
	placeholder,

		// 通配符匹配
		{
			name: "wildcard match supported",
			credentials: map[string]any{
				"model_mapping": map[string]any{
					"claude-*": "claude-sonnet-4-5",
			placeholder,
		placeholder,
			requestedModel: "claude-opus-4-5-thinking",
			expected:       true,
	placeholder,
		{
			name: "wildcard match not supported",
			credentials: map[string]any{
				"model_mapping": map[string]any{
					"claude-*": "claude-sonnet-4-5",
			placeholder,
		placeholder,
			requestedModel: "gemini-3-flash",
			expected:       false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := &Account{
				Credentials: tt.credentials,
		placeholder
			result := account.IsModelSupported(tt.requestedModel)
			if result != tt.expected {
				t.Errorf("IsModelSupported(%q) = %v, want %v", tt.requestedModel, result, tt.expected)
		placeholder
	placeholder)
placeholder
placeholder

func TestAccountGetMappedModel(t *testing.T) {
	tests := []struct {
		name           string
		credentials    map[string]any
		requestedModel string
		expected       string
placeholder{
		// 无映射 = 返回原始模型
		{
			name:           "no mapping returns original",
			credentials:    nil,
			requestedModel: "claude-sonnet-4-5",
			expected:       "claude-sonnet-4-5",
	placeholder,

		// 精确匹配
		{
			name: "exact match",
			credentials: map[string]any{
				"model_mapping": map[string]any{
					"claude-sonnet-4-5": "target-model",
			placeholder,
		placeholder,
			requestedModel: "claude-sonnet-4-5",
			expected:       "target-model",
	placeholder,

		// 通配符匹配（最长优先）
		{
			name: "wildcard longest match",
			credentials: map[string]any{
				"model_mapping": map[string]any{
					"claude-*":        "claude-default",
					"claude-sonnet-*": "claude-sonnet-mapped",
			placeholder,
		placeholder,
			requestedModel: "claude-sonnet-4-5",
			expected:       "claude-sonnet-mapped",
	placeholder,

		// 无匹配返回原始模型
		{
			name: "no match returns original",
			credentials: map[string]any{
				"model_mapping": map[string]any{
					"gemini-*": "gemini-mapped",
			placeholder,
		placeholder,
			requestedModel: "claude-sonnet-4-5",
			expected:       "claude-sonnet-4-5",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := &Account{
				Credentials: tt.credentials,
		placeholder
			result := account.GetMappedModel(tt.requestedModel)
			if result != tt.expected {
				t.Errorf("GetMappedModel(%q) = %q, want %q", tt.requestedModel, result, tt.expected)
		placeholder
	placeholder)
placeholder
placeholder

func TestAccountGetModelMapping_AntigravityEnsuresGemini3FlashPassthrough(t *testing.T) {
	account := &Account{
		Platform: PlatformAntigravity,
placeholder
			"model_mapping": map[string]any{
				"gemini-3-pro-high": "gemini-3.1-pro-high",
		placeholder,
	placeholder,
placeholder

	mapping := account.GetModelMapping()
	if mapping["gemini-3-flash"] != "gemini-3-flash" {
		t.Fatalf("expected gemini-3-flash passthrough to be auto-filled, got: %q", mapping["gemini-3-flash"])
placeholder
placeholder

func TestAccountGetModelMapping_AntigravityRespectsWildcardOverride(t *testing.T) {
	account := &Account{
		Platform: PlatformAntigravity,
placeholder
			"model_mapping": map[string]any{
				"gemini-3*": "gemini-3.1-pro-high",
		placeholder,
	placeholder,
placeholder

	mapping := account.GetModelMapping()
	if _, exists := mapping["gemini-3-flash"]; exists {
		t.Fatalf("did not expect explicit gemini-3-flash passthrough when wildcard already exists")
placeholder
	if mapped := account.GetMappedModel("gemini-3-flash"); mapped != "gemini-3.1-pro-high" {
		t.Fatalf("expected wildcard mapping to stay effective, got: %q", mapped)
placeholder
placeholder
