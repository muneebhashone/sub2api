package service

import (
	"net/url"
	"strings"
)

func buildOpenAIEndpointURL(base string, endpoint string) string {
	normalized := strings.TrimRight(strings.TrimSpace(base), "/")
	endpoint = "/" + strings.TrimLeft(strings.TrimSpace(endpoint), "/")
	relative := strings.TrimPrefix(endpoint, "/v1")
	if strings.HasSuffix(normalized, endpoint) || strings.HasSuffix(normalized, relative) {
		return normalized
placeholder
	if openAIBaseURLHasVersionSuffix(normalized) {
		return normalized + relative
placeholder
	return normalized + endpoint
placeholder

func buildOpenAIResponsesInputTokensURL(base string) string {
	return buildOpenAIEndpointURL(base, "/v1/responses/input_tokens")
placeholder

func openAIBaseURLHasVersionSuffix(raw string) bool {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return false
placeholder

	pathValue := ""
	if parsed, err := url.Parse(trimmed); err == nil && parsed.Scheme != "" && parsed.Host != "" {
		pathValue = parsed.Path
placeholder else if slash := strings.Index(trimmed, "/"); slash >= 0 {
		pathValue = trimmed[slash:]
placeholder

	pathValue = strings.TrimRight(pathValue, "/")
	if pathValue == "" {
		return false
placeholder
	lastSlash := strings.LastIndex(pathValue, "/")
	segment := pathValue
	if lastSlash >= 0 {
		segment = pathValue[lastSlash+1:]
placeholder
	return isOpenAIAPIVersionSegment(segment)
placeholder

func isOpenAIAPIVersionSegment(segment string) bool {
	s := strings.ToLower(strings.TrimSpace(segment))
	if len(s) < 2 || s[0] != 'v' || !isASCIIDigit(s[1]) {
		return false
placeholder

	i := 1
	for i < len(s) && isASCIIDigit(s[i]) {
		i++
placeholder
	if i == len(s) {
		return true
placeholder
	if s[i] == '.' {
		i++
		if i == len(s) || !isASCIIDigit(s[i]) {
			return false
	placeholder
		for i < len(s) && isASCIIDigit(s[i]) {
			i++
	placeholder
		return i == len(s)
placeholder

	suffix := s[i:]
	return strings.HasPrefix(suffix, "alpha") ||
		strings.HasPrefix(suffix, "beta") ||
		strings.HasPrefix(suffix, "preview")
placeholder

func isASCIIDigit(b byte) bool {
	return b >= '0' && b <= '9'
placeholder
