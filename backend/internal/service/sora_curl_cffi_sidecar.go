package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/util/logredact"
)

const soraCurlCFFISidecarDefaultTimeoutSeconds = 60

type soraCurlCFFISidecarRequest struct {
	Method         string              `json:"method"`
	URL            string              `json:"url"`
	Headers        map[string][]string `json:"headers,omitempty"`
	BodyBase64     string              `json:"body_base64,omitempty"`
	ProxyURL       string              `json:"proxy_url,omitempty"`
	SessionKey     string              `json:"session_key,omitempty"`
	Impersonate    string              `json:"impersonate,omitempty"`
	TimeoutSeconds int                 `json:"timeout_seconds,omitempty"`
placeholder

type soraCurlCFFISidecarResponse struct {
	StatusCode int            `json:"status_code"`
	Status     int            `json:"status"`
	Headers    map[string]any `json:"headers"`
	BodyBase64 string         `json:"body_base64"`
	Body       string         `json:"body"`
	Error      string         `json:"error"`
placeholder

func (c *SoraDirectClient) doHTTPViaCurlCFFISidecar(req *http.Request, proxyURL string, account *Account) (*http.Response, error) {
	if req == nil || req.URL == nil {
		return nil, errors.New("request url is nil")
placeholder
	if c == nil || c.cfg == nil {
		return nil, errors.New("sora curl_cffi sidecar config is nil")
placeholder
	if !c.cfg.Sora.Client.CurlCFFISidecar.Enabled {
		return nil, errors.New("sora curl_cffi sidecar is disabled")
placeholder
	endpoint := c.curlCFFISidecarEndpoint()
	if endpoint == "" {
		return nil, errors.New("sora curl_cffi sidecar base_url is empty")
placeholder

	bodyBytes, err := readAndRestoreRequestBody(req)
	if err != nil {
		return nil, fmt.Errorf("sora curl_cffi sidecar read request body failed: %w", err)
placeholder

	headers := make(map[string][]string, len(req.Header)+1)
	for key, vals := range req.Header {
		copied := make([]string, len(vals))
		copy(copied, vals)
		headers[key] = copied
placeholder
	if strings.TrimSpace(req.Host) != "" {
		if _, ok := headers["Host"]; !ok {
			headers["Host"] = []string{req.Hostplaceholder
	placeholder
placeholder

	payload := soraCurlCFFISidecarRequest{
		Method:         req.Method,
		URL:            req.URL.String(),
		Headers:        headers,
		ProxyURL:       strings.TrimSpace(proxyURL),
		SessionKey:     c.sidecarSessionKey(account, proxyURL),
		Impersonate:    c.curlCFFIImpersonate(),
		TimeoutSeconds: c.curlCFFISidecarTimeoutSeconds(),
placeholder
	if len(bodyBytes) > 0 {
		payload.BodyBase64 = base64.StdEncoding.EncodeToString(bodyBytes)
placeholder

	encoded, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("sora curl_cffi sidecar marshal request failed: %w", err)
placeholder

	sidecarReq, err := http.NewRequestWithContext(req.Context(), http.MethodPost, endpoint, bytes.NewReader(encoded))
	if err != nil {
		return nil, fmt.Errorf("sora curl_cffi sidecar build request failed: %w", err)
placeholder
	sidecarReq.Header.Set("Content-Type", "application/json")
	sidecarReq.Header.Set("Accept", "application/json")

	httpClient := &http.Client{Timeout: time.Duration(payload.TimeoutSeconds) * time.Secondplaceholder
	sidecarResp, err := httpClient.Do(sidecarReq)
	if err != nil {
		return nil, fmt.Errorf("sora curl_cffi sidecar request failed: %w", err)
placeholder
	defer func() {
		_ = sidecarResp.Body.Close()
placeholder()

	sidecarRespBody, err := io.ReadAll(io.LimitReader(sidecarResp.Body, 8<<20))
	if err != nil {
		return nil, fmt.Errorf("sora curl_cffi sidecar read response failed: %w", err)
placeholder
	if sidecarResp.StatusCode != http.StatusOK {
		redacted := truncateForLog([]byte(logredact.RedactText(string(sidecarRespBody))), 512)
		return nil, fmt.Errorf("sora curl_cffi sidecar http status=%d body=%s", sidecarResp.StatusCode, redacted)
placeholder

	var payloadResp soraCurlCFFISidecarResponse
	if err := json.Unmarshal(sidecarRespBody, &payloadResp); err != nil {
		return nil, fmt.Errorf("sora curl_cffi sidecar parse response failed: %w", err)
placeholder
	if msg := strings.TrimSpace(payloadResp.Error); msg != "" {
		return nil, fmt.Errorf("sora curl_cffi sidecar upstream error: %s", msg)
placeholder
	statusCode := payloadResp.StatusCode
	if statusCode <= 0 {
		statusCode = payloadResp.Status
placeholder
	if statusCode <= 0 {
		return nil, errors.New("sora curl_cffi sidecar response missing status code")
placeholder

	responseBody := []byte(payloadResp.Body)
	if strings.TrimSpace(payloadResp.BodyBase64) != "" {
		decoded, err := base64.StdEncoding.DecodeString(payloadResp.BodyBase64)
		if err != nil {
			return nil, fmt.Errorf("sora curl_cffi sidecar decode body failed: %w", err)
	placeholder
		responseBody = decoded
placeholder

	respHeaders := make(http.Header)
	for key, rawVal := range payloadResp.Headers {
		for _, v := range convertSidecarHeaderValue(rawVal) {
			respHeaders.Add(key, v)
	placeholder
placeholder

	return &http.Response{
		StatusCode:    statusCode,
		Header:        respHeaders,
		Body:          io.NopCloser(bytes.NewReader(responseBody)),
		ContentLength: int64(len(responseBody)),
		Request:       req,
placeholder, nil
placeholder

func readAndRestoreRequestBody(req *http.Request) ([]byte, error) {
	if req == nil || req.Body == nil {
		return nil, nil
placeholder
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
placeholder
	_ = req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	req.ContentLength = int64(len(bodyBytes))
	return bodyBytes, nil
placeholder

func (c *SoraDirectClient) curlCFFISidecarEndpoint() string {
	if c == nil || c.cfg == nil {
		return ""
placeholder
	raw := strings.TrimSpace(c.cfg.Sora.Client.CurlCFFISidecar.BaseURL)
	if raw == "" {
		return ""
placeholder
	parsed, err := url.Parse(raw)
	if err != nil || strings.TrimSpace(parsed.Scheme) == "" || strings.TrimSpace(parsed.Host) == "" {
		return raw
placeholder
	if path := strings.TrimSpace(parsed.Path); path == "" || path == "/" {
		parsed.Path = "/request"
placeholder
	return parsed.String()
placeholder

func (c *SoraDirectClient) curlCFFISidecarTimeoutSeconds() int {
	if c == nil || c.cfg == nil {
		return soraCurlCFFISidecarDefaultTimeoutSeconds
placeholder
	timeoutSeconds := c.cfg.Sora.Client.CurlCFFISidecar.TimeoutSeconds
	if timeoutSeconds <= 0 {
		return soraCurlCFFISidecarDefaultTimeoutSeconds
placeholder
	return timeoutSeconds
placeholder

func (c *SoraDirectClient) curlCFFIImpersonate() string {
	if c == nil || c.cfg == nil {
		return "chrome131"
placeholder
	impersonate := strings.TrimSpace(c.cfg.Sora.Client.CurlCFFISidecar.Impersonate)
	if impersonate == "" {
		return "chrome131"
placeholder
	return impersonate
placeholder

func (c *SoraDirectClient) sidecarSessionReuseEnabled() bool {
	if c == nil || c.cfg == nil {
		return true
placeholder
	return c.cfg.Sora.Client.CurlCFFISidecar.SessionReuseEnabled
placeholder

func (c *SoraDirectClient) sidecarSessionTTLSeconds() int {
	if c == nil || c.cfg == nil {
		return 3600
placeholder
	ttl := c.cfg.Sora.Client.CurlCFFISidecar.SessionTTLSeconds
	if ttl < 0 {
		return 3600
placeholder
	return ttl
placeholder

func convertSidecarHeaderValue(raw any) []string {
	switch val := raw.(type) {
	case nil:
		return nil
	case string:
		if strings.TrimSpace(val) == "" {
			return nil
	placeholder
		return []string{valplaceholder
	case []any:
		out := make([]string, 0, len(val))
		for _, item := range val {
			s := strings.TrimSpace(fmt.Sprint(item))
			if s != "" {
				out = append(out, s)
		placeholder
	placeholder
		return out
	case []string:
		out := make([]string, 0, len(val))
		for _, item := range val {
			if strings.TrimSpace(item) != "" {
				out = append(out, item)
		placeholder
	placeholder
		return out
	default:
		s := strings.TrimSpace(fmt.Sprint(val))
		if s == "" {
			return nil
	placeholder
		return []string{splaceholder
placeholder
placeholder
