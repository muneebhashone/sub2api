//go:build unit

package ip

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
placeholder{
		// 私有 IPv4
		{"10.x 私有地址", "10.0.0.1", trueplaceholder,
		{"10.x 私有地址段末", "10.255.255.255", trueplaceholder,
		{"172.16.x 私有地址", "172.16.0.1", trueplaceholder,
		{"172.31.x 私有地址", "172.31.255.255", trueplaceholder,
		{"192.168.x 私有地址", "192.168.1.1", trueplaceholder,
		{"127.0.0.1 本地回环", "127.0.0.1", trueplaceholder,
		{"127.x 回环段", "127.255.255.255", trueplaceholder,

		// 公网 IPv4
		{"8.8.8.8 公网 DNS", "8.8.8.8", falseplaceholder,
		{"1.1.1.1 公网", "1.1.1.1", falseplaceholder,
		{"172.15.255.255 非私有", "172.15.255.255", falseplaceholder,
		{"172.32.0.0 非私有", "172.32.0.0", falseplaceholder,
		{"11.0.0.1 公网", "11.0.0.1", falseplaceholder,

		// IPv6
		{"::1 IPv6 回环", "::1", trueplaceholder,
		{"fc00:: IPv6 私有", "fc00::1", trueplaceholder,
		{"fd00:: IPv6 私有", "fd00::1", trueplaceholder,
		{"2001:db8::1 IPv6 公网", "2001:db8::1", falseplaceholder,

		// 无效输入
		{"空字符串", "", falseplaceholder,
		{"非法字符串", "not-an-ip", falseplaceholder,
		{"不完整 IP", "192.168", falseplaceholder,
placeholder

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := isPrivateIP(tc.ip)
			require.Equal(t, tc.expected, got, "isPrivateIP(%q)", tc.ip)
	placeholder)
placeholder
placeholder
