//go:build unit

package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractOrigin(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
placeholder{
		{"empty string", "", ""placeholder,
		{"whitespace only", "   ", ""placeholder,
		{"valid https", "https://pay.example.com/checkout", "https://pay.example.com"placeholder,
		{"valid http", "http://pay.example.com/checkout", "http://pay.example.com"placeholder,
		{"https with port", "https://pay.example.com:8443/checkout", "https://pay.example.com:8443"placeholder,
		{"protocol-relative //host", "//pay.example.com/path", ""placeholder,
		{"no scheme", "pay.example.com/path", ""placeholder,
		{"ftp scheme rejected", "ftp://pay.example.com/file", ""placeholder,
		{"empty host after parse", "https:///path", ""placeholder,
		{"invalid url", "://bad url", ""placeholder,
		{"only scheme", "https://", ""placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractOrigin(tt.input)
			assert.Equal(t, tt.want, got)
	placeholder)
placeholder
placeholder
