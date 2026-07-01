package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/util/urlvalidator"
	"github.com/gin-gonic/gin"
)

const (
	llmTesterMaxRequestBytes  = 12 << 20
	llmTesterMaxResponseBytes = 12 << 20
)

var (
	llmTesterVersionPathPattern = regexp.MustCompile(`/v\d+$`)
	llmTesterHTTPClient         = &http.Client{
		Timeout: 300 * time.Second,
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           llmTesterSafeDialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 240 * time.Second,
			IdleConnTimeout:       30 * time.Second,
	placeholder,
placeholder
	llmTesterDialer = &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
placeholder
	llmTesterBlockedCIDRs = mustParseLLMTesterCIDRs([]string{
		"0.0.0.0/8",
		"10.0.0.0/8",
		"100.64.0.0/10",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"::/128",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
placeholder)
)

type llmTesterProxyRequest struct {
	BaseURL string          `json:"base_url"`
	APIKey  string          `json:"api_key"`
	Payload json.RawMessage `json:"payload,omitempty"`
placeholder

func RegisterLLMTesterRoutes(v1 *gin.RouterGroup) {
	tester := v1.Group("/llm-tester")
	{
		tester.POST("/models", llmTesterProxyModels)
		tester.POST("/chat/completions", llmTesterProxyChatCompletions)
		tester.POST("/images/generations", llmTesterProxyImageGenerations)
		tester.POST("/videos/generations", llmTesterProxyVideoGenerations)
		tester.POST("/responses", llmTesterProxyResponses)
placeholder
placeholder

func llmTesterProxyModels(c *gin.Context) {
	var req llmTesterProxyRequest
	if !bindLLMTesterProxyRequest(c, &req) {
		return
placeholder
	forwardLLMTesterRequest(c, req, http.MethodGet, "models", nil)
placeholder

func llmTesterProxyChatCompletions(c *gin.Context) {
	var req llmTesterProxyRequest
	if !bindLLMTesterProxyRequest(c, &req) {
		return
placeholder
	if len(bytes.TrimSpace(req.Payload)) == 0 {
		response.BadRequest(c, "payload is required")
		return
placeholder
	forwardLLMTesterRequest(c, req, http.MethodPost, "chat/completions", bytes.NewReader(req.Payload))
placeholder

func llmTesterProxyImageGenerations(c *gin.Context) {
	var req llmTesterProxyRequest
	if !bindLLMTesterProxyRequest(c, &req) {
		return
placeholder
	if len(bytes.TrimSpace(req.Payload)) == 0 {
		response.BadRequest(c, "payload is required")
		return
placeholder
	forwardLLMTesterRequest(c, req, http.MethodPost, "images/generations", bytes.NewReader(req.Payload))
placeholder

func llmTesterProxyVideoGenerations(c *gin.Context) {
	var req llmTesterProxyRequest
	if !bindLLMTesterProxyRequest(c, &req) {
		return
placeholder
	if len(bytes.TrimSpace(req.Payload)) == 0 {
		response.BadRequest(c, "payload is required")
		return
placeholder
	forwardLLMTesterRequest(c, req, http.MethodPost, "videos/generations", bytes.NewReader(req.Payload))
placeholder

func llmTesterProxyResponses(c *gin.Context) {
	var req llmTesterProxyRequest
	if !bindLLMTesterProxyRequest(c, &req) {
		return
placeholder
	if len(bytes.TrimSpace(req.Payload)) == 0 {
		response.BadRequest(c, "payload is required")
		return
placeholder
	forwardLLMTesterRequest(c, req, http.MethodPost, "responses", bytes.NewReader(req.Payload))
placeholder

func bindLLMTesterProxyRequest(c *gin.Context, req *llmTesterProxyRequest) bool {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, llmTesterMaxRequestBytes)
	if err := json.NewDecoder(c.Request.Body).Decode(req); err != nil {
		response.BadRequest(c, "invalid request body")
		return false
placeholder
	if strings.TrimSpace(req.BaseURL) == "" {
		response.BadRequest(c, "base_url is required")
		return false
placeholder
	if strings.TrimSpace(req.APIKey) == "" {
		response.BadRequest(c, "api_key is required")
		return false
placeholder
	if len(req.APIKey) > 8192 {
		response.BadRequest(c, "api_key is too long")
		return false
placeholder
	return true
placeholder

func forwardLLMTesterRequest(c *gin.Context, req llmTesterProxyRequest, method, resource string, body io.Reader) {
	endpoint, err := buildLLMTesterEndpoint(req.BaseURL, resource)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
placeholder

	upstreamReq, err := http.NewRequestWithContext(c.Request.Context(), method, endpoint, body)
	if err != nil {
		response.BadRequest(c, "invalid upstream request")
		return
placeholder
	upstreamReq.Header.Set("Authorization", "Bearer "+strings.TrimSpace(req.APIKey))
	upstreamReq.Header.Set("Accept", "application/json")
	upstreamReq.Header.Set("User-Agent", "Sub2API-LLM-Tester/1.0")
	upstreamReq.Header.Set("X-Title", "Sub2API LLM Tester")
	if method == http.MethodPost {
		upstreamReq.Header.Set("Content-Type", "application/json")
placeholder
	if origin := c.GetHeader("Origin"); origin != "" {
		upstreamReq.Header.Set("HTTP-Referer", origin)
placeholder

	upstreamResp, err := llmTesterHTTPClient.Do(upstreamReq)
	if err != nil {
		response.Error(c, http.StatusBadGateway, fmt.Sprintf("upstream request failed: %s", err.Error()))
		return
placeholder
	defer upstreamResp.Body.Close()

	payload, err := readLLMTesterResponseBody(upstreamResp.Body)
	if err != nil {
		response.Error(c, http.StatusBadGateway, err.Error())
		return
placeholder

	contentType := upstreamResp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
placeholder
	c.Data(upstreamResp.StatusCode, contentType, payload)
placeholder

func buildLLMTesterEndpoint(baseURL, resource string) (string, error) {
	normalized, err := urlvalidator.ValidateHTTPSURL(baseURL, urlvalidator.ValidationOptions{placeholder)
	if err != nil {
		return "", err
placeholder
	parsed, err := url.Parse(normalized)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", errors.New("invalid base_url")
placeholder
	if parsed.User != nil {
		return "", errors.New("base_url must not include user info")
placeholder
	if err := urlvalidator.ValidateResolvedIP(parsed.Hostname()); err != nil {
		return "", err
placeholder
	parsed.RawQuery = ""
	parsed.Fragment = ""
	parsed.Path = strings.TrimRight(parsed.Path, "/")
	if !llmTesterVersionPathPattern.MatchString(parsed.Path) {
		parsed.Path = strings.TrimRight(parsed.Path, "/") + "/v1"
placeholder
	parsed.Path = strings.TrimRight(parsed.Path, "/") + "/" + strings.TrimLeft(resource, "/")
	return parsed.String(), nil
placeholder

func readLLMTesterResponseBody(body io.Reader) ([]byte, error) {
	limited := io.LimitReader(body, llmTesterMaxResponseBytes+1)
	payload, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("failed to read upstream response: %w", err)
placeholder
	if len(payload) > llmTesterMaxResponseBytes {
		return nil, errors.New("upstream response is too large")
placeholder
	return payload, nil
placeholder

func llmTesterSafeDialContext(ctx context.Context, network, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
placeholder
	if llmTesterBlockedHost(host) {
		return nil, &net.AddrError{Err: "blocked by SSRF policy", Addr: addressplaceholder
placeholder
	if ip := net.ParseIP(host); ip != nil {
		if llmTesterBlockedIP(ip) {
			return nil, &net.AddrError{Err: "blocked by SSRF policy", Addr: addressplaceholder
	placeholder
		return llmTesterDialer.DialContext(ctx, network, address)
placeholder

	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
placeholder
	if len(addrs) == 0 {
		return nil, &net.AddrError{Err: "no addresses for host", Addr: hostplaceholder
placeholder

	var lastErr error
	for _, addr := range addrs {
		if llmTesterBlockedIP(addr.IP) {
			lastErr = &net.AddrError{Err: "blocked by SSRF policy", Addr: addr.IP.String()placeholder
			continue
	placeholder
		conn, err := llmTesterDialer.DialContext(ctx, network, net.JoinHostPort(addr.IP.String(), port))
		if err == nil {
			return conn, nil
	placeholder
		lastErr = err
placeholder
	if lastErr == nil {
		lastErr = &net.AddrError{Err: "no usable addresses", Addr: hostplaceholder
placeholder
	return nil, lastErr
placeholder

func llmTesterBlockedHost(host string) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	return host == "" ||
		host == "localhost" ||
		strings.HasSuffix(host, ".localhost") ||
		host == "metadata" ||
		host == "metadata.google.internal" ||
		host == "metadata.goog" ||
		host == "instance-data" ||
		host == "instance-data.ec2.internal"
placeholder

func llmTesterBlockedIP(ip net.IP) bool {
	if ip == nil {
		return true
placeholder
	if ip.IsUnspecified() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsInterfaceLocalMulticast() || ip.IsPrivate() {
		return true
placeholder
	for _, cidr := range llmTesterBlockedCIDRs {
		if cidr.Contains(ip) {
			return true
	placeholder
placeholder
	return false
placeholder

func mustParseLLMTesterCIDRs(raw []string) []*net.IPNet {
	out := make([]*net.IPNet, 0, len(raw))
	for _, value := range raw {
		_, cidr, err := net.ParseCIDR(value)
		if err != nil {
			panic("llm_tester: invalid blocked CIDR " + value + ": " + err.Error())
	placeholder
		out = append(out, cidr)
placeholder
	return out
placeholder
