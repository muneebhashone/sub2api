package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
)

// SoraS3Storage 负责 Sora 媒体文件的 S3 存储操作。
// 从 Settings 表读取 S3 配置，初始化并缓存 S3 客户端。
type SoraS3Storage struct {
	settingService *SettingService

	mu     sync.RWMutex
	client *s3.Client
	cfg    *SoraS3Settings // 上次加载的配置快照

	healthCheckedAt time.Time
	healthErr       error
	healthTTL       time.Duration
placeholder

const defaultSoraS3HealthTTL = 30 * time.Second

// UpstreamDownloadError 表示从上游下载媒体失败（包含 HTTP 状态码）。
type UpstreamDownloadError struct {
	StatusCode int
placeholder

func (e *UpstreamDownloadError) Error() string {
	if e == nil {
		return "upstream download failed"
placeholder
	return fmt.Sprintf("upstream returned %d", e.StatusCode)
placeholder

// NewSoraS3Storage 创建 S3 存储服务实例。
func NewSoraS3Storage(settingService *SettingService) *SoraS3Storage {
	return &SoraS3Storage{
		settingService: settingService,
		healthTTL:      defaultSoraS3HealthTTL,
placeholder
placeholder

// Enabled 返回 S3 存储是否已启用且配置有效。
func (s *SoraS3Storage) Enabled(ctx context.Context) bool {
	cfg, err := s.getConfig(ctx)
	if err != nil || cfg == nil {
		return false
placeholder
	return cfg.Enabled && cfg.Bucket != ""
placeholder

// getConfig 获取当前 S3 配置（从 settings 表读取）。
func (s *SoraS3Storage) getConfig(ctx context.Context) (*SoraS3Settings, error) {
	if s.settingService == nil {
		return nil, fmt.Errorf("setting service not available")
placeholder
	return s.settingService.GetSoraS3Settings(ctx)
placeholder

// getClient 获取或初始化 S3 客户端（带缓存）。
// 配置变更时调用 RefreshClient 清除缓存。
func (s *SoraS3Storage) getClient(ctx context.Context) (*s3.Client, *SoraS3Settings, error) {
	s.mu.RLock()
	if s.client != nil && s.cfg != nil {
		client, cfg := s.client, s.cfg
		s.mu.RUnlock()
		return client, cfg, nil
placeholder
	s.mu.RUnlock()

	return s.initClient(ctx)
placeholder

func (s *SoraS3Storage) initClient(ctx context.Context) (*s3.Client, *SoraS3Settings, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 双重检查
	if s.client != nil && s.cfg != nil {
		return s.client, s.cfg, nil
placeholder

	cfg, err := s.getConfig(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("load s3 config: %w", err)
placeholder
	if !cfg.Enabled {
		return nil, nil, fmt.Errorf("sora s3 storage is disabled")
placeholder
	if cfg.Bucket == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
		return nil, nil, fmt.Errorf("sora s3 config incomplete: bucket, access_key_id, secret_access_key are required")
placeholder

	client, region, err := buildSoraS3Client(ctx, cfg)
	if err != nil {
		return nil, nil, err
placeholder

	s.client = client
	s.cfg = cfg
	logger.LegacyPrintf("service.sora_s3", "[SoraS3] 客户端已初始化 bucket=%s endpoint=%s region=%s", cfg.Bucket, cfg.Endpoint, region)
	return client, cfg, nil
placeholder

// RefreshClient 清除缓存的 S3 客户端，下次使用时重新初始化。
// 应在系统设置中 S3 配置变更时调用。
func (s *SoraS3Storage) RefreshClient() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.client = nil
	s.cfg = nil
	s.healthCheckedAt = time.Time{placeholder
	s.healthErr = nil
	logger.LegacyPrintf("service.sora_s3", "[SoraS3] 客户端缓存已清除，下次使用将重新初始化")
placeholder

// TestConnection 测试 S3 连接（HeadBucket）。
func (s *SoraS3Storage) TestConnection(ctx context.Context) error {
	client, cfg, err := s.getClient(ctx)
	if err != nil {
		return err
placeholder
	_, err = client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: &cfg.Bucket,
placeholder)
	if err != nil {
		return fmt.Errorf("s3 HeadBucket failed: %w", err)
placeholder
	return nil
placeholder

// IsHealthy 返回 S3 健康状态（带短缓存，避免每次请求都触发 HeadBucket）。
func (s *SoraS3Storage) IsHealthy(ctx context.Context) bool {
	if s == nil {
		return false
placeholder
	now := time.Now()
	s.mu.RLock()
	lastCheck := s.healthCheckedAt
	lastErr := s.healthErr
	ttl := s.healthTTL
	s.mu.RUnlock()

	if ttl <= 0 {
		ttl = defaultSoraS3HealthTTL
placeholder
	if !lastCheck.IsZero() && now.Sub(lastCheck) < ttl {
		return lastErr == nil
placeholder

	err := s.TestConnection(ctx)
	s.mu.Lock()
	s.healthCheckedAt = time.Now()
	s.healthErr = err
	s.mu.Unlock()
	return err == nil
placeholder

// TestConnectionWithSettings 使用临时配置测试连接，不污染缓存的客户端。
func (s *SoraS3Storage) TestConnectionWithSettings(ctx context.Context, cfg *SoraS3Settings) error {
	if cfg == nil {
		return fmt.Errorf("s3 config is required")
placeholder
	if !cfg.Enabled {
		return fmt.Errorf("sora s3 storage is disabled")
placeholder
	if cfg.Endpoint == "" || cfg.Bucket == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
		return fmt.Errorf("sora s3 config incomplete: endpoint, bucket, access_key_id, secret_access_key are required")
placeholder
	client, _, err := buildSoraS3Client(ctx, cfg)
	if err != nil {
		return err
placeholder
	_, err = client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: &cfg.Bucket,
placeholder)
	if err != nil {
		return fmt.Errorf("s3 HeadBucket failed: %w", err)
placeholder
	return nil
placeholder

// GenerateObjectKey 生成 S3 object key。
// 格式: {prefixplaceholdersora/{userIDplaceholder/{YYYY/MM/DDplaceholder/{uuidplaceholder.{extplaceholder
func (s *SoraS3Storage) GenerateObjectKey(prefix string, userID int64, ext string) string {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
placeholder
	datePath := time.Now().Format("2006/01/02")
	key := fmt.Sprintf("sora/%d/%s/%s%s", userID, datePath, uuid.NewString(), ext)
	if prefix != "" {
		prefix = strings.TrimRight(prefix, "/") + "/"
		key = prefix + key
placeholder
	return key
placeholder

// UploadFromURL 从上游 URL 下载并流式上传到 S3。
// 返回 S3 object key。
func (s *SoraS3Storage) UploadFromURL(ctx context.Context, userID int64, sourceURL string) (string, int64, error) {
	client, cfg, err := s.getClient(ctx)
	if err != nil {
		return "", 0, err
placeholder

	// 下载源文件
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return "", 0, fmt.Errorf("create download request: %w", err)
placeholder
	httpClient := &http.Client{Timeout: 5 * time.Minuteplaceholder
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("download from upstream: %w", err)
placeholder
	defer func() {
		_ = resp.Body.Close()
placeholder()

	if resp.StatusCode != http.StatusOK {
		return "", 0, &UpstreamDownloadError{StatusCode: resp.StatusCodeplaceholder
placeholder

	// 推断文件扩展名
	ext := fileExtFromURL(sourceURL)
	if ext == "" {
		ext = fileExtFromContentType(resp.Header.Get("Content-Type"))
placeholder
	if ext == "" {
		ext = ".bin"
placeholder

	objectKey := s.GenerateObjectKey(cfg.Prefix, userID, ext)

	// 检测 Content-Type
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
placeholder

	reader, writer := io.Pipe()
	uploadErrCh := make(chan error, 1)
	go func() {
		defer close(uploadErrCh)
		input := &s3.PutObjectInput{
			Bucket:      &cfg.Bucket,
			Key:         &objectKey,
			Body:        reader,
			ContentType: &contentType,
	placeholder
		if resp.ContentLength >= 0 {
			input.ContentLength = &resp.ContentLength
	placeholder
		_, uploadErr := client.PutObject(ctx, input)
		uploadErrCh <- uploadErr
placeholder()

	written, copyErr := io.CopyBuffer(writer, resp.Body, make([]byte, 1024*1024))
	_ = writer.CloseWithError(copyErr)
	uploadErr := <-uploadErrCh
	if copyErr != nil {
		return "", 0, fmt.Errorf("stream upload copy failed: %w", copyErr)
placeholder
	if uploadErr != nil {
		return "", 0, fmt.Errorf("s3 upload: %w", uploadErr)
placeholder

	logger.LegacyPrintf("service.sora_s3", "[SoraS3] 上传完成 key=%s size=%d", objectKey, written)
	return objectKey, written, nil
placeholder

func buildSoraS3Client(ctx context.Context, cfg *SoraS3Settings) (*s3.Client, string, error) {
	if cfg == nil {
		return nil, "", fmt.Errorf("s3 config is required")
placeholder
	region := cfg.Region
	if region == "" {
		region = "us-east-1"
placeholder

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		),
	)
	if err != nil {
		return nil, "", fmt.Errorf("load aws config: %w", err)
placeholder

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = &cfg.Endpoint
	placeholder
		if cfg.ForcePathStyle {
			o.UsePathStyle = true
	placeholder
		o.APIOptions = append(o.APIOptions, v4.SwapComputePayloadSHA256ForUnsignedPayloadMiddleware)
		// 兼容非 TLS 连接（如 MinIO）的流式上传，避免 io.Pipe checksum 校验失败
		o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
placeholder)
	return client, region, nil
placeholder

// DeleteObjects 删除一组 S3 object（遍历逐一删除）。
func (s *SoraS3Storage) DeleteObjects(ctx context.Context, objectKeys []string) error {
	if len(objectKeys) == 0 {
		return nil
placeholder

	client, cfg, err := s.getClient(ctx)
	if err != nil {
		return err
placeholder

	var lastErr error
	for _, key := range objectKeys {
		k := key
		_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: &cfg.Bucket,
			Key:    &k,
	placeholder)
		if err != nil {
			logger.LegacyPrintf("service.sora_s3", "[SoraS3] 删除失败 key=%s err=%v", key, err)
			lastErr = err
	placeholder
placeholder
	return lastErr
placeholder

// GetAccessURL 获取 S3 文件的访问 URL。
// CDN URL 优先，否则生成 24h 预签名 URL。
func (s *SoraS3Storage) GetAccessURL(ctx context.Context, objectKey string) (string, error) {
	_, cfg, err := s.getClient(ctx)
	if err != nil {
		return "", err
placeholder

	// CDN URL 优先
	if cfg.CDNURL != "" {
		cdnBase := strings.TrimRight(cfg.CDNURL, "/")
		return cdnBase + "/" + objectKey, nil
placeholder

	// 生成 24h 预签名 URL
	return s.GeneratePresignedURL(ctx, objectKey, 24*time.Hour)
placeholder

// GeneratePresignedURL 生成预签名 URL。
func (s *SoraS3Storage) GeneratePresignedURL(ctx context.Context, objectKey string, ttl time.Duration) (string, error) {
	client, cfg, err := s.getClient(ctx)
	if err != nil {
		return "", err
placeholder

	presignClient := s3.NewPresignClient(client)
	result, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &cfg.Bucket,
		Key:    &objectKey,
placeholder, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", fmt.Errorf("presign url: %w", err)
placeholder
	return result.URL, nil
placeholder

// GetMediaType 从 object key 推断媒体类型（image/video）。
func GetMediaTypeFromKey(objectKey string) string {
	ext := strings.ToLower(path.Ext(objectKey))
	switch ext {
	case ".mp4", ".mov", ".webm", ".m4v", ".avi", ".mkv", ".3gp", ".flv":
		return "video"
	default:
		return "image"
placeholder
placeholder
