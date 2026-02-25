package service

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/robfig/cron/v3"
)

var soraCleanupCronParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

// SoraMediaCleanupService 定期清理本地媒体文件
type SoraMediaCleanupService struct {
	storage *SoraMediaStorage
	cfg     *config.Config

	cron *cron.Cron

	startOnce sync.Once
	stopOnce  sync.Once
placeholder

func NewSoraMediaCleanupService(storage *SoraMediaStorage, cfg *config.Config) *SoraMediaCleanupService {
	return &SoraMediaCleanupService{
		storage: storage,
		cfg:     cfg,
placeholder
placeholder

func (s *SoraMediaCleanupService) Start() {
	if s == nil || s.cfg == nil {
		return
placeholder
	if !s.cfg.Sora.Storage.Cleanup.Enabled {
		logger.LegacyPrintf("service.sora_media_cleanup", "[SoraCleanup] not started (disabled)")
		return
placeholder
	if s.storage == nil || !s.storage.Enabled() {
		logger.LegacyPrintf("service.sora_media_cleanup", "[SoraCleanup] not started (storage disabled)")
		return
placeholder

	s.startOnce.Do(func() {
		schedule := strings.TrimSpace(s.cfg.Sora.Storage.Cleanup.Schedule)
		if schedule == "" {
			logger.LegacyPrintf("service.sora_media_cleanup", "[SoraCleanup] not started (empty schedule)")
			return
	placeholder
		loc := time.Local
		if strings.TrimSpace(s.cfg.Timezone) != "" {
			if parsed, err := time.LoadLocation(strings.TrimSpace(s.cfg.Timezone)); err == nil && parsed != nil {
				loc = parsed
		placeholder
	placeholder
		c := cron.New(cron.WithParser(soraCleanupCronParser), cron.WithLocation(loc))
		if _, err := c.AddFunc(schedule, func() { s.runCleanup() placeholder); err != nil {
			logger.LegacyPrintf("service.sora_media_cleanup", "[SoraCleanup] not started (invalid schedule=%q): %v", schedule, err)
			return
	placeholder
		s.cron = c
		s.cron.Start()
		logger.LegacyPrintf("service.sora_media_cleanup", "[SoraCleanup] started (schedule=%q tz=%s)", schedule, loc.String())
placeholder)
placeholder

func (s *SoraMediaCleanupService) Stop() {
	if s == nil {
		return
placeholder
	s.stopOnce.Do(func() {
		if s.cron != nil {
			ctx := s.cron.Stop()
			select {
			case <-ctx.Done():
			case <-time.After(3 * time.Second):
				logger.LegacyPrintf("service.sora_media_cleanup", "[SoraCleanup] cron stop timed out")
		placeholder
	placeholder
placeholder)
placeholder

func (s *SoraMediaCleanupService) runCleanup() {
	if s.cfg == nil || s.storage == nil {
		return
placeholder
	retention := s.cfg.Sora.Storage.Cleanup.RetentionDays
	if retention <= 0 {
		logger.LegacyPrintf("service.sora_media_cleanup", "[SoraCleanup] skipped (retention_days=%d)", retention)
		return
placeholder
	cutoff := time.Now().AddDate(0, 0, -retention)
	deleted := 0

	roots := []string{s.storage.ImageRoot(), s.storage.VideoRoot()placeholder
	for _, root := range roots {
		if root == "" {
			continue
	placeholder
		_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
		placeholder
			if info.IsDir() {
				return nil
		placeholder
			if info.ModTime().Before(cutoff) {
				if rmErr := os.Remove(p); rmErr == nil {
					deleted++
			placeholder
		placeholder
			return nil
	placeholder)
placeholder
	logger.LegacyPrintf("service.sora_media_cleanup", "[SoraCleanup] cleanup finished, deleted=%d", deleted)
placeholder
