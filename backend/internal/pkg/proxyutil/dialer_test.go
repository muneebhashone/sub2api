package proxyutil

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigureTransportProxy_Nil(t *testing.T) {
	transport := &http.Transport{placeholder
	err := ConfigureTransportProxy(transport, nil)

placeholder
	assert.Nil(t, transport.Proxy, "nil proxy should not set Proxy")
	assert.Nil(t, transport.DialContext, "nil proxy should not set DialContext")
placeholder

func TestConfigureTransportProxy_HTTP(t *testing.T) {
	transport := &http.Transport{placeholder
	proxyURL, _ := url.Parse("http://proxy.example.com:8080")

	err := ConfigureTransportProxy(transport, proxyURL)

placeholder
	assert.NotNil(t, transport.Proxy, "HTTP proxy should set Proxy")
	assert.Nil(t, transport.DialContext, "HTTP proxy should not set DialContext")
placeholder

func TestConfigureTransportProxy_HTTPS(t *testing.T) {
	transport := &http.Transport{placeholder
	proxyURL, _ := url.Parse("https://secure-proxy.example.com:8443")

	err := ConfigureTransportProxy(transport, proxyURL)

placeholder
	assert.NotNil(t, transport.Proxy, "HTTPS proxy should set Proxy")
	assert.Nil(t, transport.DialContext, "HTTPS proxy should not set DialContext")
placeholder

func TestConfigureTransportProxy_SOCKS5(t *testing.T) {
	transport := &http.Transport{placeholder
	proxyURL, _ := url.Parse("socks5://socks.example.com:1080")

	err := ConfigureTransportProxy(transport, proxyURL)

placeholder
	assert.Nil(t, transport.Proxy, "SOCKS5 proxy should not set Proxy")
	assert.NotNil(t, transport.DialContext, "SOCKS5 proxy should set DialContext")
placeholder

func TestConfigureTransportProxy_SOCKS5H(t *testing.T) {
	transport := &http.Transport{placeholder
	proxyURL, _ := url.Parse("socks5h://socks.example.com:1080")

	err := ConfigureTransportProxy(transport, proxyURL)

placeholder
	assert.Nil(t, transport.Proxy, "SOCKS5H proxy should not set Proxy")
	assert.NotNil(t, transport.DialContext, "SOCKS5H proxy should set DialContext")
placeholder

func TestConfigureTransportProxy_CaseInsensitive(t *testing.T) {
	testCases := []struct {
		scheme   string
		useProxy bool // true = uses Transport.Proxy, false = uses DialContext
placeholder{
		{"HTTP://proxy.example.com:8080", trueplaceholder,
		{"Http://proxy.example.com:8080", trueplaceholder,
		{"HTTPS://proxy.example.com:8443", trueplaceholder,
		{"Https://proxy.example.com:8443", trueplaceholder,
		{"SOCKS5://socks.example.com:1080", falseplaceholder,
		{"Socks5://socks.example.com:1080", falseplaceholder,
		{"SOCKS5H://socks.example.com:1080", falseplaceholder,
		{"Socks5h://socks.example.com:1080", falseplaceholder,
placeholder

	for _, tc := range testCases {
		t.Run(tc.scheme, func(t *testing.T) {
			transport := &http.Transport{placeholder
			proxyURL, _ := url.Parse(tc.scheme)

			err := ConfigureTransportProxy(transport, proxyURL)

		placeholder
			if tc.useProxy {
				assert.NotNil(t, transport.Proxy)
				assert.Nil(t, transport.DialContext)
		placeholder else {
				assert.Nil(t, transport.Proxy)
				assert.NotNil(t, transport.DialContext)
		placeholder
	placeholder)
placeholder
placeholder

func TestConfigureTransportProxy_Unsupported(t *testing.T) {
	testCases := []string{
		"ftp://ftp.example.com",
		"file:///path/to/file",
		"unknown://example.com",
placeholder

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			transport := &http.Transport{placeholder
			proxyURL, _ := url.Parse(tc)

			err := ConfigureTransportProxy(transport, proxyURL)

		placeholder
			assert.Contains(t, err.Error(), "unsupported proxy scheme")
	placeholder)
placeholder
placeholder

func TestConfigureTransportProxy_WithAuth(t *testing.T) {
	transport := &http.Transport{placeholder
	proxyURL, _ := url.Parse("socks5://user:password@socks.example.com:1080")

	err := ConfigureTransportProxy(transport, proxyURL)

placeholder
	assert.NotNil(t, transport.DialContext, "SOCKS5 with auth should set DialContext")
placeholder

func TestConfigureTransportProxy_EmptyScheme(t *testing.T) {
	transport := &http.Transport{placeholder
	// 空 scheme 的 URL
	proxyURL := &url.URL{Host: "proxy.example.com:8080"placeholder

	err := ConfigureTransportProxy(transport, proxyURL)

placeholder
	assert.Contains(t, err.Error(), "unsupported proxy scheme")
placeholder

func TestConfigureTransportProxy_PreservesExistingConfig(t *testing.T) {
	// 验证代理配置不会覆盖 Transport 的其他配置
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
placeholder
	proxyURL, _ := url.Parse("socks5://socks.example.com:1080")

	err := ConfigureTransportProxy(transport, proxyURL)

placeholder
	assert.Equal(t, 100, transport.MaxIdleConns, "MaxIdleConns should be preserved")
	assert.Equal(t, 10, transport.MaxIdleConnsPerHost, "MaxIdleConnsPerHost should be preserved")
	assert.NotNil(t, transport.DialContext, "DialContext should be set")
placeholder

func TestConfigureTransportProxy_IPv6(t *testing.T) {
	testCases := []struct {
		name     string
		proxyURL string
placeholder{
		{"SOCKS5H with IPv6 loopback", "socks5h://[::1]:1080"placeholder,
		{"SOCKS5 with full IPv6", "socks5://[2001:db8::1]:1080"placeholder,
		{"HTTP with IPv6", "http://[::1]:8080"placeholder,
placeholder

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			transport := &http.Transport{placeholder
			proxyURL, err := url.Parse(tc.proxyURL)
			require.NoError(t, err, "URL should be parseable")

			err = ConfigureTransportProxy(transport, proxyURL)
		placeholder
	placeholder)
placeholder
placeholder

func TestConfigureTransportProxy_SpecialCharsInPassword(t *testing.T) {
	testCases := []struct {
		name     string
		proxyURL string
placeholder{
		// 密码包含 @ 符号（URL 编码为 %40）
		{"password with @", "socks5://user:p%40ssword@proxy.example.com:1080"placeholder,
		// 密码包含 : 符号（URL 编码为 %3A）
		{"password with :", "socks5://user:pass%3Aword@proxy.example.com:1080"placeholder,
		// 密码包含 / 符号（URL 编码为 %2F）
		{"password with /", "socks5://user:pass%2Fword@proxy.example.com:1080"placeholder,
		// 复杂密码
		{"complex password", "socks5h://admin:P%40ss%3Aw0rd%2F123@proxy.example.com:1080"placeholder,
placeholder

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			transport := &http.Transport{placeholder
			proxyURL, err := url.Parse(tc.proxyURL)
			require.NoError(t, err, "URL should be parseable")

			err = ConfigureTransportProxy(transport, proxyURL)
		placeholder
			assert.NotNil(t, transport.DialContext, "SOCKS5 should set DialContext")
	placeholder)
placeholder
placeholder
