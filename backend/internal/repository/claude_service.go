package repository

import (
	"net/http"
	"net/url"
	"time"

	"sub2api/internal/config"
	"sub2api/internal/service"
)

type claudeUpstreamService struct {
	defaultClient *http.Client
	cfg           *config.Config
placeholder

func NewClaudeUpstream(cfg *config.Config) service.ClaudeUpstream {
	responseHeaderTimeout := time.Duration(cfg.Gateway.ResponseHeaderTimeout) * time.Second
	if responseHeaderTimeout == 0 {
		responseHeaderTimeout = 300 * time.Second
placeholder

	transport := &http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		ResponseHeaderTimeout: responseHeaderTimeout,
placeholder

	return &claudeUpstreamService{
		defaultClient: &http.Client{Transport: transportplaceholder,
		cfg:           cfg,
placeholder
placeholder

func (s *claudeUpstreamService) Do(req *http.Request, proxyURL string) (*http.Response, error) {
	if proxyURL == "" {
		return s.defaultClient.Do(req)
placeholder
	client := s.createProxyClient(proxyURL)
	return client.Do(req)
placeholder

func (s *claudeUpstreamService) createProxyClient(proxyURL string) *http.Client {
	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return s.defaultClient
placeholder

	responseHeaderTimeout := time.Duration(s.cfg.Gateway.ResponseHeaderTimeout) * time.Second
	if responseHeaderTimeout == 0 {
		responseHeaderTimeout = 300 * time.Second
placeholder

	transport := &http.Transport{
		Proxy:                 http.ProxyURL(parsedURL),
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		ResponseHeaderTimeout: responseHeaderTimeout,
placeholder

	return &http.Client{Transport: transportplaceholder
placeholder
