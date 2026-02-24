package urlvalidator

import "testing"

func TestValidateURLFormat(t *testing.T) {
	if _, err := ValidateURLFormat("", false); err == nil {
		t.Fatalf("expected empty url to fail")
placeholder
	if _, err := ValidateURLFormat("://bad", false); err == nil {
		t.Fatalf("expected invalid url to fail")
placeholder
	if _, err := ValidateURLFormat("http://example.com", false); err == nil {
		t.Fatalf("expected http to fail when allow_insecure_http is false")
placeholder
	if _, err := ValidateURLFormat("https://example.com", false); err != nil {
		t.Fatalf("expected https to pass, got %v", err)
placeholder
	if _, err := ValidateURLFormat("http://example.com", true); err != nil {
		t.Fatalf("expected http to pass when allow_insecure_http is true, got %v", err)
placeholder
	if _, err := ValidateURLFormat("https://example.com:bad", true); err == nil {
		t.Fatalf("expected invalid port to fail")
placeholder

	// 验证末尾斜杠被移除
	normalized, err := ValidateURLFormat("https://example.com/", false)
	if err != nil {
		t.Fatalf("expected trailing slash url to pass, got %v", err)
placeholder
	if normalized != "https://example.com" {
		t.Fatalf("expected trailing slash to be removed, got %s", normalized)
placeholder

	// 验证多个末尾斜杠被移除
	normalized, err = ValidateURLFormat("https://example.com///", false)
	if err != nil {
		t.Fatalf("expected multiple trailing slashes to pass, got %v", err)
placeholder
	if normalized != "https://example.com" {
		t.Fatalf("expected all trailing slashes to be removed, got %s", normalized)
placeholder

	// 验证带路径的 URL 末尾斜杠被移除
	normalized, err = ValidateURLFormat("https://example.com/api/v1/", false)
	if err != nil {
		t.Fatalf("expected trailing slash url with path to pass, got %v", err)
placeholder
	if normalized != "https://example.com/api/v1" {
		t.Fatalf("expected trailing slash to be removed from path, got %s", normalized)
placeholder
placeholder

func TestValidateHTTPURL(t *testing.T) {
	if _, err := ValidateHTTPURL("http://example.com", false, ValidationOptions{placeholder); err == nil {
		t.Fatalf("expected http to fail when allow_insecure_http is false")
placeholder
	if _, err := ValidateHTTPURL("http://example.com", true, ValidationOptions{placeholder); err != nil {
		t.Fatalf("expected http to pass when allow_insecure_http is true, got %v", err)
placeholder
	if _, err := ValidateHTTPURL("https://example.com", false, ValidationOptions{RequireAllowlist: trueplaceholder); err == nil {
		t.Fatalf("expected require allowlist to fail when empty")
placeholder
	if _, err := ValidateHTTPURL("https://example.com", false, ValidationOptions{AllowedHosts: []string{"api.example.com"placeholderplaceholder); err == nil {
		t.Fatalf("expected host not in allowlist to fail")
placeholder
	if _, err := ValidateHTTPURL("https://api.example.com", false, ValidationOptions{AllowedHosts: []string{"api.example.com"placeholderplaceholder); err != nil {
		t.Fatalf("expected allowlisted host to pass, got %v", err)
placeholder
	if _, err := ValidateHTTPURL("https://sub.api.example.com", false, ValidationOptions{AllowedHosts: []string{"*.example.com"placeholderplaceholder); err != nil {
		t.Fatalf("expected wildcard allowlist to pass, got %v", err)
placeholder
	if _, err := ValidateHTTPURL("https://localhost", false, ValidationOptions{AllowPrivate: falseplaceholder); err == nil {
		t.Fatalf("expected localhost to be blocked when allow_private_hosts is false")
placeholder
placeholder
