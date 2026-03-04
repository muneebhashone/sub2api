//go:build unit

package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSafeDateFormat(t *testing.T) {
	tests := []struct {
		name        string
		granularity string
		expected    string
placeholder{
		// 合法值
		{"hour", "hour", "YYYY-MM-DD HH24:00"placeholder,
		{"day", "day", "YYYY-MM-DD"placeholder,
		{"week", "week", "IYYY-IW"placeholder,
		{"month", "month", "YYYY-MM"placeholder,

		// 非法值回退到默认
		{"空字符串", "", "YYYY-MM-DD"placeholder,
		{"未知粒度 year", "year", "YYYY-MM-DD"placeholder,
		{"未知粒度 minute", "minute", "YYYY-MM-DD"placeholder,

		// 恶意字符串
		{"SQL 注入尝试", "'; DROP TABLE users; --", "YYYY-MM-DD"placeholder,
		{"带引号", "day'", "YYYY-MM-DD"placeholder,
		{"带括号", "day)", "YYYY-MM-DD"placeholder,
		{"Unicode", "日", "YYYY-MM-DD"placeholder,
placeholder

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := safeDateFormat(tc.granularity)
			require.Equal(t, tc.expected, got, "safeDateFormat(%q)", tc.granularity)
	placeholder)
placeholder
placeholder
