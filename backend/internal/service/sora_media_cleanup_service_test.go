//go:build unit

package service

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func TestSoraMediaCleanupService_RunCleanup_NilCfg(t *testing.T) {
	storage := &SoraMediaStorage{placeholder
	svc := &SoraMediaCleanupService{storage: storage, cfg: nilplaceholder
	// 不应 panic
	svc.runCleanup()
placeholder

func TestSoraMediaCleanupService_RunCleanup_NilStorage(t *testing.T) {
	cfg := &config.Config{placeholder
	svc := &SoraMediaCleanupService{storage: nil, cfg: cfgplaceholder
	// 不应 panic
	svc.runCleanup()
placeholder

func TestSoraMediaCleanupService_RunCleanup_ZeroRetention(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Sora: config.SoraConfig{
			Storage: config.SoraStorageConfig{
				Type:      "local",
				LocalPath: tmpDir,
				Cleanup: config.SoraStorageCleanupConfig{
					Enabled:       true,
					RetentionDays: 0,
			placeholder,
		placeholder,
	placeholder,
placeholder
	storage := NewSoraMediaStorage(cfg)
	svc := NewSoraMediaCleanupService(storage, cfg)
	// retention=0 应跳过清理
	svc.runCleanup()
placeholder

func TestSoraMediaCleanupService_Start_NilCfg(t *testing.T) {
	svc := NewSoraMediaCleanupService(nil, nil)
	svc.Start() // cfg == nil 时应直接返回
placeholder

func TestSoraMediaCleanupService_Start_StorageDisabled(t *testing.T) {
	cfg := &config.Config{
		Sora: config.SoraConfig{
			Storage: config.SoraStorageConfig{
				Cleanup: config.SoraStorageCleanupConfig{
					Enabled: true,
			placeholder,
		placeholder,
	placeholder,
placeholder
	svc := NewSoraMediaCleanupService(nil, cfg)
	svc.Start() // storage == nil 时应直接返回
placeholder

func TestSoraMediaCleanupService_Start_WithTimezone(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Timezone: "Asia/Shanghai",
		Sora: config.SoraConfig{
			Storage: config.SoraStorageConfig{
				Type:      "local",
				LocalPath: tmpDir,
				Cleanup: config.SoraStorageCleanupConfig{
					Enabled:  true,
					Schedule: "0 3 * * *",
			placeholder,
		placeholder,
	placeholder,
placeholder
	storage := NewSoraMediaStorage(cfg)
	svc := NewSoraMediaCleanupService(storage, cfg)
	svc.Start()
	t.Cleanup(svc.Stop)
placeholder

func TestSoraMediaCleanupService_Start_Disabled(t *testing.T) {
	cfg := &config.Config{
		Sora: config.SoraConfig{
			Storage: config.SoraStorageConfig{
				Cleanup: config.SoraStorageCleanupConfig{
					Enabled: false,
			placeholder,
		placeholder,
	placeholder,
placeholder
	svc := NewSoraMediaCleanupService(nil, cfg)
	svc.Start() // 不应 panic，也不应启动 cron
placeholder

func TestSoraMediaCleanupService_Start_NilSelf(t *testing.T) {
	var svc *SoraMediaCleanupService
	svc.Start() // 不应 panic
placeholder

func TestSoraMediaCleanupService_Start_EmptySchedule(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Sora: config.SoraConfig{
			Storage: config.SoraStorageConfig{
				Type:      "local",
				LocalPath: tmpDir,
				Cleanup: config.SoraStorageCleanupConfig{
					Enabled:  true,
					Schedule: "",
			placeholder,
		placeholder,
	placeholder,
placeholder
	storage := NewSoraMediaStorage(cfg)
	svc := NewSoraMediaCleanupService(storage, cfg)
	svc.Start() // 空 schedule 不应启动
placeholder

func TestSoraMediaCleanupService_Start_InvalidSchedule(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Sora: config.SoraConfig{
			Storage: config.SoraStorageConfig{
				Type:      "local",
				LocalPath: tmpDir,
				Cleanup: config.SoraStorageCleanupConfig{
					Enabled:  true,
					Schedule: "invalid-cron",
			placeholder,
		placeholder,
	placeholder,
placeholder
	storage := NewSoraMediaStorage(cfg)
	svc := NewSoraMediaCleanupService(storage, cfg)
	svc.Start() // 无效 schedule 不应 panic
placeholder

func TestSoraMediaCleanupService_Start_ValidSchedule(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Sora: config.SoraConfig{
			Storage: config.SoraStorageConfig{
				Type:      "local",
				LocalPath: tmpDir,
				Cleanup: config.SoraStorageCleanupConfig{
					Enabled:  true,
					Schedule: "0 3 * * *",
			placeholder,
		placeholder,
	placeholder,
placeholder
	storage := NewSoraMediaStorage(cfg)
	svc := NewSoraMediaCleanupService(storage, cfg)
	svc.Start()
	t.Cleanup(svc.Stop)
placeholder

func TestSoraMediaCleanupService_Stop_NilSelf(t *testing.T) {
	var svc *SoraMediaCleanupService
	svc.Stop() // 不应 panic
placeholder

func TestSoraMediaCleanupService_Stop_WithoutStart(t *testing.T) {
	svc := NewSoraMediaCleanupService(nil, &config.Config{placeholder)
	svc.Stop() // cron 未启动时 Stop 不应 panic
placeholder

func TestSoraMediaCleanupService_RunCleanup(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Sora: config.SoraConfig{
			Storage: config.SoraStorageConfig{
				Type:      "local",
				LocalPath: tmpDir,
				Cleanup: config.SoraStorageCleanupConfig{
					Enabled:       true,
					RetentionDays: 1,
			placeholder,
		placeholder,
	placeholder,
placeholder

	storage := NewSoraMediaStorage(cfg)
	require.NoError(t, storage.EnsureLocalDirs())

	oldImage := filepath.Join(storage.ImageRoot(), "old.png")
	newVideo := filepath.Join(storage.VideoRoot(), "new.mp4")
	require.NoError(t, os.WriteFile(oldImage, []byte("old"), 0o644))
	require.NoError(t, os.WriteFile(newVideo, []byte("new"), 0o644))

	oldTime := time.Now().Add(-48 * time.Hour)
	require.NoError(t, os.Chtimes(oldImage, oldTime, oldTime))

	cleanup := NewSoraMediaCleanupService(storage, cfg)
	cleanup.runCleanup()

	require.NoFileExists(t, oldImage)
	require.FileExists(t, newVideo)
placeholder
