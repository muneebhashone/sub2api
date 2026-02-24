package soraerror

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var (
	cfRayPattern  = regexp.MustCompile(`(?i)cf-ray[:\s=]+([a-z0-9-]+)`)
	cRayPattern   = regexp.MustCompile(`(?i)cRay:\s*'([a-z0-9-]+)'`)
	htmlChallenge = []string{
		"window._cf_chl_opt",
		"just a moment",
		"enable javascript and cookies to continue",
		"__cf_chl_",
		"challenge-platform",
placeholder
)

// IsCloudflareChallengeResponse reports whether the upstream response matches Cloudflare challenge behavior.
func IsCloudflareChallengeResponse(statusCode int, headers http.Header, body []byte) bool {
	if statusCode != http.StatusForbidden && statusCode != http.StatusTooManyRequests {
		return false
placeholder

	if headers != nil && strings.EqualFold(strings.TrimSpace(headers.Get("cf-mitigated")), "challenge") {
		return true
placeholder

	preview := strings.ToLower(TruncateBody(body, 4096))
	for _, marker := range htmlChallenge {
		if strings.Contains(preview, marker) {
			return true
	placeholder
placeholder

	contentType := ""
	if headers != nil {
		contentType = strings.ToLower(strings.TrimSpace(headers.Get("content-type")))
placeholder
	if strings.Contains(contentType, "text/html") &&
		(strings.Contains(preview, "<html") || strings.Contains(preview, "<!doctype html")) &&
		(strings.Contains(preview, "cloudflare") || strings.Contains(preview, "challenge")) {
		return true
placeholder

	return false
placeholder

// ExtractCloudflareRayID extracts cf-ray from headers or response body.
func ExtractCloudflareRayID(headers http.Header, body []byte) string {
	if headers != nil {
		rayID := strings.TrimSpace(headers.Get("cf-ray"))
		if rayID != "" {
			return rayID
	placeholder
		rayID = strings.TrimSpace(headers.Get("Cf-Ray"))
		if rayID != "" {
			return rayID
	placeholder
placeholder

	preview := TruncateBody(body, 8192)
	if matches := cfRayPattern.FindStringSubmatch(preview); len(matches) >= 2 {
		return strings.TrimSpace(matches[1])
placeholder
	if matches := cRayPattern.FindStringSubmatch(preview); len(matches) >= 2 {
		return strings.TrimSpace(matches[1])
placeholder
	return ""
placeholder

// FormatCloudflareChallengeMessage appends cf-ray info when available.
func FormatCloudflareChallengeMessage(base string, headers http.Header, body []byte) string {
	rayID := ExtractCloudflareRayID(headers, body)
	if rayID == "" {
		return base
placeholder
	return fmt.Sprintf("%s (cf-ray: %s)", base, rayID)
placeholder

// ExtractUpstreamErrorCodeAndMessage extracts structured error code/message from common JSON layouts.
func ExtractUpstreamErrorCodeAndMessage(body []byte) (string, string) {
	trimmed := strings.TrimSpace(string(body))
	if trimmed == "" {
		return "", ""
placeholder
	if !json.Valid([]byte(trimmed)) {
		return "", truncateMessage(trimmed, 256)
placeholder

	var payload map[string]any
	if err := json.Unmarshal([]byte(trimmed), &payload); err != nil {
		return "", truncateMessage(trimmed, 256)
placeholder

	code := firstNonEmpty(
		extractNestedString(payload, "error", "code"),
		extractRootString(payload, "code"),
	)
	message := firstNonEmpty(
		extractNestedString(payload, "error", "message"),
		extractRootString(payload, "message"),
		extractNestedString(payload, "error", "detail"),
		extractRootString(payload, "detail"),
	)
	return strings.TrimSpace(code), truncateMessage(strings.TrimSpace(message), 512)
placeholder

// TruncateBody truncates body text for logging/inspection.
func TruncateBody(body []byte, max int) string {
	if max <= 0 {
		max = 512
placeholder
	raw := strings.TrimSpace(string(body))
	if len(raw) <= max {
		return raw
placeholder
	return raw[:max] + "...(truncated)"
placeholder

func truncateMessage(s string, max int) string {
	if max <= 0 {
		return ""
placeholder
	if len(s) <= max {
		return s
placeholder
	return s[:max] + "...(truncated)"
placeholder

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
	placeholder
placeholder
	return ""
placeholder

func extractRootString(m map[string]any, key string) string {
	if m == nil {
		return ""
placeholder
	v, ok := m[key]
	if !ok {
		return ""
placeholder
	s, _ := v.(string)
	return s
placeholder

func extractNestedString(m map[string]any, parent, key string) string {
	if m == nil {
		return ""
placeholder
	node, ok := m[parent]
	if !ok {
		return ""
placeholder
	child, ok := node.(map[string]any)
	if !ok {
		return ""
placeholder
	s, _ := child[key].(string)
	return s
placeholder
