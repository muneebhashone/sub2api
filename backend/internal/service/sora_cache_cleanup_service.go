package service

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
)

const (
	soraCacheCleanupInterval = time.Hour
	soraCacheCleanupBatch    = 200
)

// SoraCacheCleanupService 负责清理 Sora 视频缓存文件。
type SoraCacheCleanupService struct {
	cacheRepo      SoraCacheFileRepository
	settingService *SettingService
	cfg            *config.Config
	stopCh         chan struct{placeholder
	stopOnce       sync.Once
placeholder

func NewSoraCacheCleanupService(cacheRepo SoraCacheFileRepository, settingService *SettingService, cfg *config.Config) *SoraCacheCleanupService {
	return &SoraCacheCleanupService{
		cacheRepo:      cacheRepo,
		settingService: settingService,
		cfg:            cfg,
		stopCh:         make(chan struct{placeholder),
placeholder
placeholder

func (s *SoraCacheCleanupService) Start() {
	if s == nil || s.cacheRepo == nil {
		return
placeholder
	go s.cleanupLoop()
placeholder

func (s *SoraCacheCleanupService) Stop() {
	if s == nil {
		return
placeholder
	s.stopOnce.Do(func() {
		close(s.stopCh)
placeholder)
placeholder

func (s *SoraCacheCleanupService) cleanupLoop() {
	ticker := time.NewTicker(soraCacheCleanupInterval)
	defer ticker.Stop()

	s.cleanupOnce()
	for {
		select {
		case <-ticker.C:
			s.cleanupOnce()
		case <-s.stopCh:
			return
	placeholder
placeholder
placeholder

func (s *SoraCacheCleanupService) cleanupOnce() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	if s.cacheRepo == nil {
		return
placeholder

	cfg := s.getSoraConfig(ctx)
	videoDir := strings.TrimSpace(cfg.Cache.VideoDir)
	if videoDir == "" {
		return
placeholder
	maxBytes := cfg.Cache.MaxBytes
	if maxBytes <= 0 {
		return
placeholder

	size, err := dirSize(videoDir)
	if err != nil {
		log.Printf("[SoraCacheCleanup] 计算目录大小失败: %v", err)
		return
placeholder
	if size <= maxBytes {
		return
placeholder

	for size > maxBytes {
		entries, err := s.cacheRepo.ListOldest(ctx, soraCacheCleanupBatch)
		if err != nil {
			log.Printf("[SoraCacheCleanup] 读取缓存记录失败: %v", err)
			return
	placeholder
		if len(entries) == 0 {
			log.Printf("[SoraCacheCleanup] 无缓存记录但目录仍超限: size=%d max=%d", size, maxBytes)
			return
	placeholder

		ids := make([]int64, 0, len(entries))
		for _, entry := range entries {
			if entry == nil {
				continue
		placeholder
			removedSize := entry.SizeBytes
			if entry.CachePath != "" {
				if info, err := os.Stat(entry.CachePath); err == nil {
					if removedSize <= 0 {
						removedSize = info.Size()
				placeholder
			placeholder
				if err := os.Remove(entry.CachePath); err != nil && !os.IsNotExist(err) {
					log.Printf("[SoraCacheCleanup] 删除缓存文件失败: path=%s err=%v", entry.CachePath, err)
			placeholder
		placeholder

			if entry.ID > 0 {
				ids = append(ids, entry.ID)
		placeholder
			if removedSize > 0 {
				size -= removedSize
				if size < 0 {
					size = 0
			placeholder
		placeholder
	placeholder

		if len(ids) > 0 {
			if err := s.cacheRepo.DeleteByIDs(ctx, ids); err != nil {
				log.Printf("[SoraCacheCleanup] 删除缓存记录失败: %v", err)
		placeholder
	placeholder

		if size > maxBytes {
			if refreshed, err := dirSize(videoDir); err == nil {
				size = refreshed
		placeholder
	placeholder
placeholder
placeholder

func (s *SoraCacheCleanupService) getSoraConfig(ctx context.Context) config.SoraConfig {
	if s.settingService != nil {
		return s.settingService.GetSoraConfig(ctx)
placeholder
	if s.cfg != nil {
		return s.cfg.Sora
placeholder
	return config.SoraConfig{placeholder
placeholder
