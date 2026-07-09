package httputil

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tidwall/gjson"
)

func TestNormalizeLenientJSONRequestBody_accepts_client_control_chars_in_strings(t *testing.T) {
	tests := []struct {
		name    string
		body    []byte
		path    string
		want    string
		wantRaw string
placeholder{
		{
			name:    "null byte in message content",
			body:    []byte("{\"messages\":[{\"content\":\"hello\x00world\"placeholder]placeholder"),
			path:    "messages.0.content",
			want:    "hello\x00world",
			wantRaw: `"hello\u0000world"`,
	placeholder,
		{
			name:    "ansi escape in message content",
			body:    []byte("{\"messages\":[{\"content\":\"hello\x1b[31mred\x1b[0m\"placeholder]placeholder"),
			path:    "messages.0.content",
			want:    "hello\x1b[31mred\x1b[0m",
			wantRaw: `"hello\u001b[31mred\u001b[0m"`,
	placeholder,
		{
			name:    "leading UTF-8 BOM",
			body:    []byte("\xef\xbb\xbf{\"input\":\"hello\"placeholder"),
			path:    "input",
			want:    "hello",
			wantRaw: `"hello"`,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			if gjson.ValidBytes(tt.body) {
				t.Fatalf("test payload should reproduce strict JSON rejection: %q", tt.body)
		placeholder

			// When
			got, err := NormalizeLenientJSONRequestBody(tt.body, 1024)
			if err != nil {
				t.Fatalf("NormalizeLenientJSONRequestBody: %v", err)
		placeholder

			// Then
			if !gjson.ValidBytes(got) {
				t.Fatalf("normalized body should be valid JSON: %q", got)
		placeholder
			result := gjson.GetBytes(got, tt.path)
			if result.String() != tt.want {
				t.Fatalf("value mismatch: got %q want %q", result.String(), tt.want)
		placeholder
			if result.Raw != tt.wantRaw {
				t.Fatalf("raw value mismatch: got %q want %q", result.Raw, tt.wantRaw)
		placeholder
	placeholder)
placeholder
placeholder

func TestNormalizeLenientJSONRequestBody_keeps_invalid_structure_invalid(t *testing.T) {
	tests := []struct {
		name string
		body []byte
placeholder{
		{
			name: "truncated JSON",
			body: []byte("{\"messages\":[{\"content\":\"hello\"placeholder]"),
	placeholder,
		{
			name: "control character outside string",
			body: []byte("{\"input\":\"hello\"placeholder\x00"),
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			got, err := NormalizeLenientJSONRequestBody(tt.body, 1024)
			if err != nil {
				t.Fatalf("NormalizeLenientJSONRequestBody: %v", err)
		placeholder

			// Then
			if gjson.ValidBytes(got) {
				t.Fatalf("normalization must not repair invalid JSON structure: %q", got)
		placeholder
	placeholder)
placeholder
placeholder

func TestNormalizeLenientJSONRequestBody_allows_http_requests_with_client_control_chars(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Given
		body, err := ReadLenientJSONRequestBodyWithPrealloc(r, 1024)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
	placeholder

		// When
		if !gjson.ValidBytes(body) {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
	placeholder
		w.WriteHeader(http.StatusAccepted)
placeholder))
	defer server.Close()

	tests := []struct {
		name string
		body []byte
		want int
placeholder{
		{
			name: "null byte in JSON string",
			body: []byte("{\"model\":\"gpt-5.5\",\"messages\":[{\"role\":\"user\",\"content\":\"hello\x00world\"placeholder]placeholder"),
			want: http.StatusAccepted,
	placeholder,
		{
			name: "ANSI escape in JSON string",
			body: []byte("{\"model\":\"gpt-5.5\",\"messages\":[{\"role\":\"user\",\"content\":\"hello\x1b[31mred\x1b[0m\"placeholder]placeholder"),
			want: http.StatusAccepted,
	placeholder,
		{
			name: "leading UTF-8 BOM",
			body: []byte("\xef\xbb\xbf{\"model\":\"gpt-5.5\",\"messages\":[{\"role\":\"user\",\"content\":\"hello\"placeholder]placeholder"),
			want: http.StatusAccepted,
	placeholder,
		{
			name: "truncated JSON",
			body: []byte("{\"model\":\"gpt-5.5\",\"messages\":[{\"role\":\"user\",\"content\":\"hello\"placeholder]"),
			want: http.StatusBadRequest,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, server.URL+"/v1/chat/completions", bytes.NewReader(tt.body))
			if err != nil {
				t.Fatalf("NewRequest: %v", err)
		placeholder
			req.Header.Set("Content-Type", "application/json")

			resp, err := server.Client().Do(req)
			if err != nil {
				t.Fatalf("Do: %v", err)
		placeholder
			defer func() { _ = resp.Body.Close() placeholder()

			if resp.StatusCode != tt.want {
				t.Fatalf("status mismatch: got %d want %d", resp.StatusCode, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestNormalizeLenientJSONRequestBody_rejects_expansion_past_limit(t *testing.T) {
	// Given
	body := []byte("{\"input\":\"\x00\x00\"placeholder")

	// When
	_, err := NormalizeLenientJSONRequestBody(body, int64(len(body)+5))

	// Then
	var maxErr *http.MaxBytesError
	if !errors.As(err, &maxErr) {
		t.Fatalf("expected MaxBytesError, got %T %v", err, err)
placeholder
	if maxErr.Limit != int64(len(body)+5) {
		t.Fatalf("limit mismatch: got %d want %d", maxErr.Limit, len(body)+5)
placeholder
placeholder
