package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/httpclient"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

func NewProxyExitInfoProber(cfg *config.Config) service.ProxyExitInfoProber {
	insecure := false
	allowPrivate := false
	if cfg != nil {
		insecure = cfg.Security.ProxyProbe.InsecureSkipVerify
		allowPrivate = cfg.Security.URLAllowlist.AllowPrivateHosts
placeholder
	if insecure {
		log.Printf("[ProxyProbe] Warning: TLS verification is disabled for proxy probing.")
placeholder
	return &proxyProbeService{
		ipInfoURL:          defaultIPInfoURL,
		insecureSkipVerify: insecure,
		allowPrivateHosts:  allowPrivate,
placeholder
placeholder

const defaultIPInfoURL = "https://ipinfo.io/json"

type proxyProbeService struct {
	ipInfoURL          string
	insecureSkipVerify bool
	allowPrivateHosts  bool
placeholder

func (s *proxyProbeService) ProbeProxy(ctx context.Context, proxyURL string) (*service.ProxyExitInfo, int64, error) {
	client, err := httpclient.GetClient(httpclient.Options{
		ProxyURL:           proxyURL,
		Timeout:            15 * time.Second,
		InsecureSkipVerify: s.insecureSkipVerify,
		ProxyStrict:        true,
		ValidateResolvedIP: true,
		AllowPrivateHosts:  s.allowPrivateHosts,
placeholder)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create proxy client: %w", err)
placeholder

	startTime := time.Now()
	req, err := http.NewRequestWithContext(ctx, "GET", s.ipInfoURL, nil)
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
