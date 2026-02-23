//go:build unit

package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func TestSoraMediaStorage_StoreFromURLs(t *testing.T) {
	tmpDir := t.TempDir()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("data"))
placeholder))
	defer server.Close()

	cfg := &config.Config{
		Sora: config.SoraConfig{
			Storage: config.SoraStorageConfig{
				Type:                   "local",
				LocalPath:              tmpDir,
				MaxConcurrentDownloads: 1,
		placeholder,
	placeholder,
placeholder

	storage := NewSoraMediaStorage(cfg)
	urls, err := storage.StoreFromURLs(context.Background(), "image", []string{server.URL + "/img.png"placeholder)
placeholder
	require.Len(t, urls, 1)
	require.True(t, strings.HasPrefix(urls[0], "/image/"))
	require.True(t, strings.HasSuffix(urls[0], ".png"))

	localPath := filepath.Join(tmpDir, filepath.FromSlash(strings.TrimPrefix(urls[0], "/")))
	require.FileExists(t, localPath)
placeholder

func TestSoraMediaStorage_FallbackToUpstream(t *testing.T) {
	tmpDir := t.TempDir()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
placeholder))
	defer server.Close()

	cfg := &config.Config{
		Sora: config.SoraConfig{
			Storage: config.SoraStorageConfig{
				Type:               "local",
				LocalPath:          tmpDir,
				FallbackToUpstream: true,
		placeholder,
	placeholder,
placeholder

	storage := NewSoraMediaStorage(cfg)
	url := server.URL + "/broken.png"
	urls, err := storage.StoreFromURLs(context.Background(), "image", []string{urlplaceholder)
placeholder
	require.Equal(t, []string{urlplaceholder, urls)
placeholder

func TestSoraMediaStorage_MaxDownloadBytes(t *testing.T) {
	tmpDir := t.TempDir()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("too-large"))
placeholder))
	defer server.Close()

	cfg := &config.Config{
		Sora: config.SoraConfig{
			Storage: config.SoraStorageConfig{
				Type:             "local",
				LocalPath:        tmpDir,
				MaxDownloadBytes: 1,
		placeholder,
	placeholder,
placeholder

	storage := NewSoraMediaStorage(cfg)
	_, err := storage.StoreFromURLs(context.Background(), "image", []string{server.URL + "/img.png"placeholder)
placeholder
placeholder

func TestNormalizeSoraFileExt(t *testing.T) {
	require.Equal(t, ".png", normalizeSoraFileExt(".PNG"))
	require.Equal(t, ".mp4", normalizeSoraFileExt(".mp4"))
	require.Equal(t, "", normalizeSoraFileExt("../../etc/passwd"))
	require.Equal(t, "", normalizeSoraFileExt(".php"))
placeholder

func TestRemovePartialDownload(t *testing.T) {
	tmpDir := t.TempDir()
	root, err := os.OpenRoot(tmpDir)
placeholder
	defer func() { _ = root.Close() placeholder()

	filePath := "partial.bin"
	f, err := root.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
placeholder
	_, _ = f.WriteString("partial")
	_ = f.Close()

	removePartialDownload(root, filePath)
	_, err = root.Stat(filePath)
placeholder
	require.True(t, os.IsNotExist(err))
placeholder
