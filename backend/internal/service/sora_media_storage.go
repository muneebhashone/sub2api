package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/google/uuid"
)

const (
	soraStorageDefaultRoot = "/app/data/sora"
)

// SoraMediaStorage 负责下载并落地 Sora 媒体
type SoraMediaStorage struct {
	cfg                *config.Config
	root               string
	imageRoot          string
	videoRoot          string
	downloadTimeout    time.Duration
	maxDownloadBytes   int64
	fallbackToUpstream bool
	debug              bool
	sem                chan struct{placeholder
	ready              bool
placeholder

func NewSoraMediaStorage(cfg *config.Config) *SoraMediaStorage {
	storage := &SoraMediaStorage{cfg: cfgplaceholder
	storage.refreshConfig()
	if storage.Enabled() {
		if err := storage.EnsureLocalDirs(); err != nil {
			log.Printf("[SoraStorage] 初始化失败: %v", err)
	placeholder
placeholder
	return storage
placeholder

func (s *SoraMediaStorage) Enabled() bool {
	if s == nil || s.cfg == nil {
		return false
placeholder
	return strings.ToLower(strings.TrimSpace(s.cfg.Sora.Storage.Type)) == "local"
placeholder

func (s *SoraMediaStorage) Root() string {
	if s == nil {
		return ""
placeholder
	return s.root
placeholder

func (s *SoraMediaStorage) ImageRoot() string {
	if s == nil {
		return ""
placeholder
	return s.imageRoot
placeholder

func (s *SoraMediaStorage) VideoRoot() string {
	if s == nil {
		return ""
placeholder
	return s.videoRoot
placeholder

func (s *SoraMediaStorage) refreshConfig() {
	if s == nil || s.cfg == nil {
		return
placeholder
	root := strings.TrimSpace(s.cfg.Sora.Storage.LocalPath)
	if root == "" {
		root = soraStorageDefaultRoot
placeholder
	root = filepath.Clean(root)
	if !filepath.IsAbs(root) {
		if absRoot, err := filepath.Abs(root); err == nil {
			root = absRoot
	placeholder
placeholder
	s.root = root
	s.imageRoot = filepath.Join(root, "image")
	s.videoRoot = filepath.Join(root, "video")

	maxConcurrent := s.cfg.Sora.Storage.MaxConcurrentDownloads
	if maxConcurrent <= 0 {
		maxConcurrent = 4
placeholder
	timeoutSeconds := s.cfg.Sora.Storage.DownloadTimeoutSeconds
	if timeoutSeconds <= 0 {
		timeoutSeconds = 120
placeholder
	s.downloadTimeout = time.Duration(timeoutSeconds) * time.Second

	maxBytes := s.cfg.Sora.Storage.MaxDownloadBytes
	if maxBytes <= 0 {
		maxBytes = 200 << 20
placeholder
	s.maxDownloadBytes = maxBytes
	s.fallbackToUpstream = s.cfg.Sora.Storage.FallbackToUpstream
	s.debug = s.cfg.Sora.Storage.Debug
	s.sem = make(chan struct{placeholder, maxConcurrent)
placeholder

// EnsureLocalDirs 创建并校验本地目录
func (s *SoraMediaStorage) EnsureLocalDirs() error {
	if s == nil || !s.Enabled() {
		return nil
placeholder
	if err := os.MkdirAll(s.imageRoot, 0o755); err != nil {
		return fmt.Errorf("create image dir: %w", err)
placeholder
	if err := os.MkdirAll(s.videoRoot, 0o755); err != nil {
		return fmt.Errorf("create video dir: %w", err)
placeholder
	s.ready = true
	return nil
placeholder

// StoreFromURLs 下载并存储媒体，返回相对路径或回退 URL
func (s *SoraMediaStorage) StoreFromURLs(ctx context.Context, mediaType string, urls []string) ([]string, error) {
	if len(urls) == 0 {
		return nil, nil
placeholder
	if s == nil || !s.Enabled() {
		return urls, nil
placeholder
	if !s.ready {
		if err := s.EnsureLocalDirs(); err != nil {
			return nil, err
	placeholder
placeholder
	results := make([]string, 0, len(urls))
	for _, raw := range urls {
		relative, err := s.downloadAndStore(ctx, mediaType, raw)
		if err != nil {
			if s.fallbackToUpstream {
				results = append(results, raw)
				continue
		placeholder
			return nil, err
	placeholder
		results = append(results, relative)
placeholder
	return results, nil
placeholder

func (s *SoraMediaStorage) downloadAndStore(ctx context.Context, mediaType, rawURL string) (string, error) {
	if strings.TrimSpace(rawURL) == "" {
		return "", errors.New("empty url")
placeholder
	root := s.imageRoot
	if mediaType == "video" {
		root = s.videoRoot
placeholder
	if root == "" {
		return "", errors.New("storage root not configured")
placeholder

	retries := 3
	for attempt := 1; attempt <= retries; attempt++ {
		release, err := s.acquire(ctx)
		if err != nil {
			return "", err
	placeholder
		relative, err := s.downloadOnce(ctx, root, mediaType, rawURL)
		release()
		if err == nil {
			return relative, nil
	placeholder
		if s.debug {
			log.Printf("[SoraStorage] 下载失败(%d/%d): %s err=%v", attempt, retries, sanitizeMediaLogURL(rawURL), err)
	placeholder
		if attempt < retries {
			time.Sleep(time.Duration(attempt*attempt) * time.Second)
			continue
	placeholder
		return "", err
placeholder
	return "", errors.New("download retries exhausted")
placeholder

func (s *SoraMediaStorage) downloadOnce(ctx context.Context, root, mediaType, rawURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return "", err
placeholder
	client := &http.Client{Timeout: s.downloadTimeoutplaceholder
	resp, err := client.Do(req)
	if err != nil {
		return "", err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
		return "", fmt.Errorf("download failed: %d %s", resp.StatusCode, string(body))
placeholder

	ext := normalizeSoraFileExt(fileExtFromURL(rawURL))
	if ext == "" {
		ext = normalizeSoraFileExt(fileExtFromContentType(resp.Header.Get("Content-Type")))
placeholder
	if ext == "" {
		ext = ".bin"
placeholder
	if s.maxDownloadBytes > 0 && resp.ContentLength > s.maxDownloadBytes {
		return "", fmt.Errorf("download size exceeds limit: %d", resp.ContentLength)
placeholder

	storageRoot, err := os.OpenRoot(root)
	if err != nil {
		return "", err
placeholder
	defer func() { _ = storageRoot.Close() placeholder()

	datePath := time.Now().Format("2006/01/02")
	datePathFS := filepath.FromSlash(datePath)
	if err := storageRoot.MkdirAll(datePathFS, 0o755); err != nil {
		return "", err
placeholder
	filename := uuid.NewString() + ext
	filePath := filepath.Join(datePathFS, filename)
	out, err := storageRoot.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return "", err
placeholder
	defer func() { _ = out.Close() placeholder()

	limited := io.LimitReader(resp.Body, s.maxDownloadBytes+1)
	written, err := io.Copy(out, limited)
	if err != nil {
		removePartialDownload(storageRoot, filePath)
		return "", err
placeholder
	if s.maxDownloadBytes > 0 && written > s.maxDownloadBytes {
		removePartialDownload(storageRoot, filePath)
		return "", fmt.Errorf("download size exceeds limit: %d", written)
placeholder

	relative := path.Join("/", mediaType, datePath, filename)
	if s.debug {
		log.Printf("[SoraStorage] 已落地 %s -> %s", sanitizeMediaLogURL(rawURL), relative)
placeholder
	return relative, nil
placeholder

func (s *SoraMediaStorage) acquire(ctx context.Context) (func(), error) {
	if s.sem == nil {
		return func() {placeholder, nil
placeholder
	select {
	case s.sem <- struct{placeholder{placeholder:
		return func() { <-s.sem placeholder, nil
	case <-ctx.Done():
		return nil, ctx.Err()
placeholder
placeholder

func fileExtFromURL(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
placeholder
	ext := path.Ext(parsed.Path)
	return strings.ToLower(ext)
placeholder

func fileExtFromContentType(ct string) string {
	if ct == "" {
		return ""
placeholder
	if exts, err := mime.ExtensionsByType(ct); err == nil && len(exts) > 0 {
		return strings.ToLower(exts[0])
placeholder
	return ""
placeholder

func normalizeSoraFileExt(ext string) string {
	ext = strings.ToLower(strings.TrimSpace(ext))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".bmp", ".svg", ".tif", ".tiff", ".heic",
		".mp4", ".mov", ".webm", ".m4v", ".avi", ".mkv", ".3gp", ".flv":
		return ext
	default:
		return ""
placeholder
placeholder

func removePartialDownload(root *os.Root, filePath string) {
	if root == nil || strings.TrimSpace(filePath) == "" {
		return
placeholder
	_ = root.Remove(filePath)
placeholder

// sanitizeMediaLogURL 脱敏 URL 用于日志记录（去除 query 参数中可能的 token 信息）
func sanitizeMediaLogURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		if len(rawURL) > 80 {
			return rawURL[:80] + "..."
	placeholder
		return rawURL
placeholder
	safe := parsed.Scheme + "://" + parsed.Host + parsed.Path
	if len(safe) > 120 {
		return safe[:120] + "..."
placeholder
	return safe
placeholder
