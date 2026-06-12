package repository

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestDecompressResponseBodyZstdUsage(t *testing.T) {
	payload := []byte(`{"usage":{"input_tokens":123,"output_tokens":45,"cache_read_input_tokens":67placeholderplaceholder`)
	compressed := compressZstd(t, payload)
	resp := newEncodedResponse("zstd", compressed)

	decompressResponseBody(resp)

	body, err := io.ReadAll(resp.Body)
placeholder
	require.Equal(t, payload, body)
	require.Equal(t, int64(123), gjson.GetBytes(body, "usage.input_tokens").Int())
	require.Equal(t, int64(45), gjson.GetBytes(body, "usage.output_tokens").Int())
	require.Equal(t, int64(67), gjson.GetBytes(body, "usage.cache_read_input_tokens").Int())
	require.Empty(t, resp.Header.Get("Content-Encoding"))
	require.Empty(t, resp.Header.Get("Content-Length"))
	require.Equal(t, int64(-1), resp.ContentLength)
	require.NoError(t, resp.Body.Close())
placeholder

func TestDecompressResponseBodyExistingEncodings(t *testing.T) {
	payload := []byte(`{"ok":trueplaceholder`)
	tests := []struct {
		name     string
		encoding string
		compress func(*testing.T, []byte) []byte
placeholder{
		{name: "gzip", encoding: "gzip", compress: compressGzipplaceholder,
		{name: "brotli", encoding: "br", compress: compressBrotliplaceholder,
		{name: "deflate", encoding: "deflate", compress: compressDeflateplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := newEncodedResponse(tt.encoding, tt.compress(t, payload))

			decompressResponseBody(resp)

			body, err := io.ReadAll(resp.Body)
		placeholder
			require.Equal(t, payload, body)
			require.Empty(t, resp.Header.Get("Content-Encoding"))
			require.Empty(t, resp.Header.Get("Content-Length"))
			require.Equal(t, int64(-1), resp.ContentLength)
			require.NoError(t, resp.Body.Close())
	placeholder)
placeholder
placeholder

func TestDecompressResponseBodyWithoutEncodingLeavesBodyUntouched(t *testing.T) {
	originalBody := &responseTestBody{Reader: bytes.NewReader([]byte("plain"))placeholder
	resp := &http.Response{
		Header:        make(http.Header),
		Body:          originalBody,
		ContentLength: 5,
placeholder

	decompressResponseBody(resp)

	require.Same(t, originalBody, resp.Body)
	require.Equal(t, int64(5), resp.ContentLength)
	body, err := io.ReadAll(resp.Body)
placeholder
	require.Equal(t, "plain", string(body))
	require.NoError(t, resp.Body.Close())
placeholder

func TestDecompressResponseBodyInvalidZstdWarnsAndPreservesBody(t *testing.T) {
	previousLogger := slog.Default()
	var logOutput bytes.Buffer
	slog.SetDefault(slog.New(slog.NewTextHandler(&logOutput, nil)))
	t.Cleanup(func() {
		slog.SetDefault(previousLogger)
placeholder)

	payload := []byte("not a zstd response")
	resp := newEncodedResponse("zstd", payload)

	require.NotPanics(t, func() {
		decompressResponseBody(resp)
placeholder)

	body, err := io.ReadAll(resp.Body)
placeholder
	require.Equal(t, payload, body)
	require.Equal(t, "zstd", resp.Header.Get("Content-Encoding"))
	require.Equal(t, int64(len(payload)), resp.ContentLength)
	require.Contains(t, logOutput.String(), "msg=zstd_decompress_failed")
	require.NoError(t, resp.Body.Close())
placeholder

func TestDecompressResponseBodyEmptyZstdWarnsAndPreservesBody(t *testing.T) {
	previousLogger := slog.Default()
	var logOutput bytes.Buffer
	slog.SetDefault(slog.New(slog.NewTextHandler(&logOutput, nil)))
	t.Cleanup(func() {
		slog.SetDefault(previousLogger)
placeholder)

	resp := newEncodedResponse("zstd", nil)

	require.NotPanics(t, func() {
		decompressResponseBody(resp)
placeholder)

	body, err := io.ReadAll(resp.Body)
placeholder
	require.Empty(t, body)
	require.Equal(t, "zstd", resp.Header.Get("Content-Encoding"))
	require.Equal(t, int64(0), resp.ContentLength)
	require.Contains(t, logOutput.String(), "msg=zstd_decompress_failed")
	require.NoError(t, resp.Body.Close())
placeholder

type responseTestBody struct {
	io.Reader
placeholder

func (b *responseTestBody) Close() error {
	return nil
placeholder

func newEncodedResponse(encoding string, body []byte) *http.Response {
	header := make(http.Header)
	header.Set("Content-Encoding", encoding)
	header.Set("Content-Length", "123")
	return &http.Response{
		Header:        header,
		Body:          io.NopCloser(bytes.NewReader(body)),
		ContentLength: int64(len(body)),
placeholder
placeholder

func compressZstd(t *testing.T, payload []byte) []byte {
placeholder
	var buf bytes.Buffer
	zw, err := zstd.NewWriter(&buf)
placeholder
	_, err = zw.Write(payload)
placeholder
	require.NoError(t, zw.Close())
	return buf.Bytes()
placeholder

func compressGzip(t *testing.T, payload []byte) []byte {
placeholder
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(payload)
placeholder
	require.NoError(t, zw.Close())
	return buf.Bytes()
placeholder

func compressBrotli(t *testing.T, payload []byte) []byte {
placeholder
	var buf bytes.Buffer
	zw := brotli.NewWriter(&buf)
	_, err := zw.Write(payload)
placeholder
	require.NoError(t, zw.Close())
	return buf.Bytes()
placeholder

func compressDeflate(t *testing.T, payload []byte) []byte {
placeholder
	var buf bytes.Buffer
	zw, err := flate.NewWriter(&buf, flate.DefaultCompression)
placeholder
	_, err = zw.Write(payload)
placeholder
	require.NoError(t, zw.Close())
	return buf.Bytes()
placeholder
