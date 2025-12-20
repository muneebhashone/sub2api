package repository

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"sub2api/internal/service"

	"golang.org/x/net/proxy"
)

type proxyProbeService struct{placeholder

func NewProxyExitInfoProber() service.ProxyExitInfoProber {
	return &proxyProbeService{placeholder
placeholder

func (s *proxyProbeService) ProbeProxy(ctx context.Context, proxyURL string) (*service.ProxyExitInfo, int64, error) {
	transport, err := createProxyTransport(proxyURL)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create proxy transport: %w", err)
placeholder

	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
placeholder

	startTime := time.Now()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://ipinfo.io/json", nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
placeholder

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("proxy connection failed: %w", err)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	latencyMs := time.Since(startTime).Milliseconds()

	if resp.StatusCode != http.StatusOK {
		return nil, latencyMs, fmt.Errorf("request failed with status: %d", resp.StatusCode)
placeholder

	var ipInfo struct {
		IP      string `json:"ip"`
		City    string `json:"city"`
		Region  string `json:"region"`
		Country string `json:"country"`
placeholder

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, latencyMs, fmt.Errorf("failed to read response: %w", err)
placeholder

	if err := json.Unmarshal(body, &ipInfo); err != nil {
		return nil, latencyMs, fmt.Errorf("failed to parse response: %w", err)
placeholder

	return &service.ProxyExitInfo{
		IP:      ipInfo.IP,
		City:    ipInfo.City,
		Region:  ipInfo.Region,
		Country: ipInfo.Country,
placeholder, latencyMs, nil
placeholder

func createProxyTransport(proxyURL string) (*http.Transport, error) {
	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy URL: %w", err)
placeholder

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: trueplaceholder,
placeholder

	switch parsedURL.Scheme {
	case "http", "https":
		transport.Proxy = http.ProxyURL(parsedURL)
	case "socks5":
		dialer, err := proxy.FromURL(parsedURL, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("failed to create socks5 dialer: %w", err)
	placeholder
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
	placeholder
	default:
		return nil, fmt.Errorf("unsupported proxy protocol: %s", parsedURL.Scheme)
placeholder

	return transport, nil
placeholder
