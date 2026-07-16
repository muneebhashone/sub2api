package repository

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/Wei-Shaw/sub2api/internal/pkg/servertiming"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

// S3BackupStore implements service.BackupObjectStore using AWS S3 compatible storage
type S3BackupStore struct {
	client *s3.Client
	bucket string
placeholder

// NewS3BackupStoreFactory returns a BackupObjectStoreFactory that creates S3-backed stores
func NewS3BackupStoreFactory() service.BackupObjectStoreFactory {
	return func(ctx context.Context, cfg *service.BackupS3Config) (service.BackupObjectStore, error) {
		client, err := newS3Client(ctx, s3ClientParams{
			Endpoint:        cfg.Endpoint,
			Region:          cfg.Region,
			AccessKeyID:     cfg.AccessKeyID,
			SecretAccessKey: cfg.SecretAccessKey,
			ForcePathStyle:  cfg.ForcePathStyle,
	placeholder)
		if err != nil {
			return nil, err
	placeholder
		return &S3BackupStore{client: client, bucket: cfg.Bucketplaceholder, nil
placeholder
placeholder

func (s *S3BackupStore) Upload(ctx context.Context, key string, body io.Reader, contentType string) (int64, error) {
	// 读取全部内容以获取大小（S3 PutObject 需要知道内容长度）
	// 注意：阿里云 OSS 不兼容 s3manager 分片上传的签名方式，因此使用 PutObject
	data, err := io.ReadAll(body)
	if err != nil {
		return 0, fmt.Errorf("read body: %w", err)
placeholder

	finish := servertiming.ObserveDependency(ctx, "s3")
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         &key,
		Body:        bytes.NewReader(data),
		ContentType: &contentType,
placeholder)
	finish()
	if err != nil {
		return 0, fmt.Errorf("S3 PutObject: %w", err)
placeholder
	return int64(len(data)), nil
placeholder

func (s *S3BackupStore) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	finish := servertiming.ObserveDependency(ctx, "s3")
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
placeholder)
	finish()
	if err != nil {
		return nil, fmt.Errorf("S3 GetObject: %w", err)
placeholder
	return result.Body, nil
placeholder

func (s *S3BackupStore) Delete(ctx context.Context, key string) error {
	finish := servertiming.ObserveDependency(ctx, "s3")
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
placeholder)
	finish()
	return err
placeholder

func (s *S3BackupStore) PresignURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)
	result, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
placeholder, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", fmt.Errorf("presign url: %w", err)
placeholder
	return result.URL, nil
placeholder

func (s *S3BackupStore) HeadBucket(ctx context.Context) error {
	finish := servertiming.ObserveDependency(ctx, "s3")
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: &s.bucket,
placeholder)
	finish()
	if err != nil {
		return fmt.Errorf("S3 HeadBucket failed: %w", err)
placeholder
	return nil
placeholder
