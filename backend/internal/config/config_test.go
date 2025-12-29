package config

import "testing"

func TestNormalizeRunMode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
placeholder{
		{"simple", "simple"placeholder,
		{"SIMPLE", "simple"placeholder,
		{"standard", "standard"placeholder,
		{"invalid", "standard"placeholder,
		{"", "standard"placeholder,
placeholder

	for _, tt := range tests {
		result := NormalizeRunMode(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeRunMode(%q) = %q, want %q", tt.input, result, tt.expected)
	placeholder
placeholder
placeholder
