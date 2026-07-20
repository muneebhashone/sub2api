//go:build unit

package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type stubSettingRepo struct {
	mu     sync.Mutex
	values map[string]string
placeholder

func newStubSettingRepo() *stubSettingRepo {
	return &stubSettingRepo{values: map[string]string{placeholderplaceholder
placeholder

func (r *stubSettingRepo) Get(context.Context, string) (*Setting, error) { return nil, nil placeholder
func (r *stubSettingRepo) GetValue(_ context.Context, key string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.values[key], nil
placeholder

func (r *stubSettingRepo) Set(_ context.Context, key, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.values[key] = value
	return nil
placeholder
func (r *stubSettingRepo) GetMultiple(context.Context, []string) (map[string]string, error) {
	return map[string]string{placeholder, nil
placeholder
func (r *stubSettingRepo) SetMultiple(context.Context, map[string]string) error { return nil placeholder
func (r *stubSettingRepo) GetAll(context.Context) (map[string]string, error) {
	return map[string]string{placeholder, nil
placeholder
func (r *stubSettingRepo) Delete(context.Context, string) error { return nil placeholder

// reversibleEncryptor stands in for AES: prefixed so a test can tell ciphertext
// from plaintext, and so decrypting a plaintext value fails like the real one.
type reversibleEncryptor struct{placeholder

func (reversibleEncryptor) Encrypt(plaintext string) (string, error) {
	return "enc:" + plaintext, nil
placeholder

func (reversibleEncryptor) Decrypt(ciphertext string) (string, error) {
	rest, ok := strings.CutPrefix(ciphertext, "enc:")
	if !ok {
		return "", errors.New("not encrypted")
placeholder
	return rest, nil
placeholder

type recordingStorage struct{ saved []string placeholder

func (s *recordingStorage) Save(_ context.Context, key, _ string, _ []byte) (string, error) {
	s.saved = append(s.saved, key)
	return "https://cdn.example.com/" + key, nil
placeholder

func newImageStorageFixture(t *testing.T, fallback config.ImageStorageConfig) (*ImageStorageSettingService, *stubSettingRepo, *[]config.ImageStorageConfig) {
	return newImageStorageFixtureWithKey(t, fallback, true)
placeholder

func newImageStorageFixtureWithKey(t *testing.T, fallback config.ImageStorageConfig, encryptionKeyConfigured bool) (*ImageStorageSettingService, *stubSettingRepo, *[]config.ImageStorageConfig) {
placeholder
	repo := newStubSettingRepo()
	encryptor := reversibleEncryptor{placeholder
	backup := NewBackupService(repo, &config.Config{
		Totp: config.TotpConfig{EncryptionKeyConfigured: encryptionKeyConfiguredplaceholder,
placeholder, encryptor, nil, nil)

	var built []config.ImageStorageConfig
	factory := func(_ context.Context, cfg *config.ImageStorageConfig) (ImageStorage, error) {
		built = append(built, *cfg)
		return &recordingStorage{placeholder, nil
placeholder
	return NewImageStorageSettingService(repo, encryptor, backup, factory, fallback), repo, &built
placeholder

func seedBackupS3(t *testing.T, repo *stubSettingRepo, cfg BackupS3Config) {
placeholder
	cfg.SecretAccessKey = "enc:" + cfg.SecretAccessKey
	data, err := json.Marshal(cfg)
placeholder
	require.NoError(t, repo.Set(context.Background(), settingKeyBackupS3Config, string(data)))
placeholder

// The admin switch must take effect without a restart: that is the entire point
// of moving image_storage out of config.yaml (#4542).
func TestImageStorageSettingsToggleTakesEffectWithoutRestart(t *testing.T) {
	svc, repo, built := newImageStorageFixture(t, config.ImageStorageConfig{placeholder)
	ctx := context.Background()
	seedBackupS3(t, repo, BackupS3Config{
		Endpoint: "https://acct.r2.cloudflarestorage.com", Region: "auto",
		Bucket: "backup-bucket", AccessKeyID: "ak", SecretAccessKey: "sk",
		Prefix: "backups/",
placeholder)

	uploader, enabled := svc.resolve()
	require.False(t, enabled, "disabled until an admin turns it on")
	require.Nil(t, uploader)

	_, err := svc.Update(ctx, ImageStorageSettings{Enabled: true, ReuseBackupS3: trueplaceholder)
placeholder

	uploader, enabled = svc.resolve()
	require.True(t, enabled, "saving the setting must enable the feature immediately")
	require.NotNil(t, uploader)

	_, err = svc.Update(ctx, ImageStorageSettings{Enabled: false, ReuseBackupS3: trueplaceholder)
placeholder
	_, enabled = svc.resolve()
	require.False(t, enabled, "turning it back off must also apply immediately")

	require.Len(t, *built, 1, "the S3 client is built only when the feature is on")
placeholder

func TestImageStorageSettingsReuseBackupCredentials(t *testing.T) {
	svc, repo, built := newImageStorageFixture(t, config.ImageStorageConfig{placeholder)
	ctx := context.Background()
	seedBackupS3(t, repo, BackupS3Config{
		Endpoint: "https://acct.r2.cloudflarestorage.com", Region: "wnam",
		Bucket: "backup-bucket", AccessKeyID: "backup-ak", SecretAccessKey: "backup-sk",
		Prefix: "backups/", ForcePathStyle: true,
placeholder)

	_, err := svc.Update(ctx, ImageStorageSettings{Enabled: true, ReuseBackupS3: true, Prefix: "images"placeholder)
placeholder
	_, enabled := svc.resolve()
	require.True(t, enabled)

	require.Len(t, *built, 1)
	got := (*built)[0]
	require.Equal(t, "https://acct.r2.cloudflarestorage.com", got.Endpoint)
	require.Equal(t, "wnam", got.Region)
	require.Equal(t, "backup-ak", got.AccessKeyID)
	require.Equal(t, "backup-sk", got.SecretAccessKey, "the backup secret must be decrypted before use")
	require.True(t, got.ForcePathStyle)
	require.Equal(t, "backup-bucket", got.Bucket, "an empty bucket falls back to the backup bucket")
	require.Equal(t, "images/", got.Prefix, "images stay under their own prefix so they never collide with backups/")

	// Reusing must not duplicate the secret into a second row.
	raw, err := repo.GetValue(ctx, settingKeyImageStorageConfig)
placeholder
	require.NotContains(t, raw, "backup-sk")
	require.NotContains(t, raw, "enc:")
placeholder

func TestImageStorageSettingsOwnCredentialsAreEncryptedAndMasked(t *testing.T) {
	svc, repo, built := newImageStorageFixture(t, config.ImageStorageConfig{placeholder)
	ctx := context.Background()

	saved, err := svc.Update(ctx, ImageStorageSettings{
		Enabled: true, Bucket: "my-images",
		Endpoint:    "https://acct.r2.cloudflarestorage.com",
		AccessKeyID: "ak", SecretAccessKey: "super-secret",
placeholder)
placeholder
	require.Empty(t, saved.SecretAccessKey, "the response must never echo the secret back")

	raw, err := repo.GetValue(ctx, settingKeyImageStorageConfig)
placeholder
	require.NotContains(t, raw, `"secret_access_key":"super-secret"`, "the secret must be encrypted at rest")
	require.Contains(t, raw, "enc:super-secret")

	fetched, err := svc.Get(ctx)
placeholder
	require.Empty(t, fetched.SecretAccessKey)
	require.True(t, svc.SecretConfigured(ctx))

	_, enabled := svc.resolve()
	require.True(t, enabled)
	require.Equal(t, "super-secret", (*built)[0].SecretAccessKey, "the stored secret must be decrypted before use")

	// An update that omits the secret keeps the stored one rather than wiping it.
	_, err = svc.Update(ctx, ImageStorageSettings{
		Enabled: true, Bucket: "my-images",
		Endpoint: "https://acct.r2.cloudflarestorage.com", AccessKeyID: "ak",
placeholder)
placeholder
	svc.resolve()
	require.Equal(t, "super-secret", (*built)[1].SecretAccessKey)
placeholder

// Persisting the service's own S3 secret must be refused when the encryption key
// is auto-generated, otherwise the ciphertext cannot be decrypted after a
// restart (#4524). Reusing the backup credentials stays allowed because it does
// not persist a second copy of the secret.
func TestImageStorageSettingsRejectSecretWithEphemeralKey(t *testing.T) {
	svc, repo, built := newImageStorageFixtureWithKey(t, config.ImageStorageConfig{placeholder, false)
	ctx := context.Background()

	_, err := svc.Update(ctx, ImageStorageSettings{
		Enabled: true, Bucket: "my-images",
		Endpoint:    "https://acct.r2.cloudflarestorage.com",
		AccessKeyID: "ak", SecretAccessKey: "super-secret",
placeholder)
	require.ErrorIs(t, err, ErrSecretEncryptionKeyNotConfigured)

	raw, _ := repo.GetValue(ctx, settingKeyImageStorageConfig)
	require.Empty(t, raw, "nothing must be persisted when the secret is rejected")
	require.Empty(t, *built)

	// Reusing backup credentials does not persist a secret, so it stays allowed.
	seedBackupS3(t, repo, BackupS3Config{
		Endpoint: "https://acct.r2.cloudflarestorage.com", Region: "auto",
		Bucket: "backup-bucket", AccessKeyID: "ak", SecretAccessKey: "sk", Prefix: "backups/",
placeholder)
	_, err = svc.Update(ctx, ImageStorageSettings{Enabled: true, ReuseBackupS3: trueplaceholder)
placeholder
placeholder

func TestImageStorageSettingsIncompleteStaysDisabled(t *testing.T) {
	svc, _, built := newImageStorageFixture(t, config.ImageStorageConfig{placeholder)
	ctx := context.Background()

	_, err := svc.Update(ctx, ImageStorageSettings{Enabled: true, Bucket: "my-images"placeholder)
placeholder

	_, enabled := svc.resolve()
	require.False(t, enabled, "missing credentials must not enable the feature")
	require.Empty(t, *built, "no client is built from an incomplete configuration")
placeholder

// Deployments that already enabled the feature through config.yaml must keep
// working after the setting moves into the database.
func TestImageStorageSettingsFallBackToConfigFile(t *testing.T) {
	svc, _, built := newImageStorageFixture(t, config.ImageStorageConfig{
		Enabled: true, Endpoint: "https://acct.r2.cloudflarestorage.com", Region: "auto",
		Bucket: "yaml-bucket", AccessKeyID: "yaml-ak", SecretAccessKey: "yaml-sk",
		Prefix: "images/", MaxDownloadByte: 1024,
placeholder)

	_, enabled := svc.resolve()
	require.True(t, enabled, "config.yaml still enables the feature when nothing is stored yet")
	require.Equal(t, "yaml-bucket", (*built)[0].Bucket)

	fetched, err := svc.Get(context.Background())
placeholder
	require.True(t, fetched.Enabled)
	require.Equal(t, "yaml-bucket", fetched.Bucket)
	require.Empty(t, fetched.SecretAccessKey)
placeholder
