package service

import (
	"context"
	"log"
	"sync"
	"time"
)

// DeferredService provides deferred batch update functionality
type DeferredService struct {
	accountRepo AccountRepository
	timingWheel *TimingWheelService
	interval    time.Duration

	lastUsedUpdates sync.Map
placeholder

// NewDeferredService creates a new DeferredService instance
func NewDeferredService(accountRepo AccountRepository, timingWheel *TimingWheelService, interval time.Duration) *DeferredService {
	return &DeferredService{
		accountRepo: accountRepo,
		timingWheel: timingWheel,
		interval:    interval,
placeholder
placeholder

// Start starts the deferred service
func (s *DeferredService) Start() {
	s.timingWheel.ScheduleRecurring("deferred:last_used", s.interval, s.flushLastUsed)
	log.Printf("[DeferredService] Started (interval: %v)", s.interval)
placeholder

// Stop stops the deferred service
func (s *DeferredService) Stop() {
	s.timingWheel.Cancel("deferred:last_used")
	s.flushLastUsed()
	log.Printf("[DeferredService] Service stopped")
placeholder

func (s *DeferredService) ScheduleLastUsedUpdate(accountID int64) {
	s.lastUsedUpdates.Store(accountID, time.Now())
placeholder

func (s *DeferredService) flushLastUsed() {
	updates := make(map[int64]time.Time)
	s.lastUsedUpdates.Range(func(key, value any) bool {
		id, ok := key.(int64)
		if !ok {
			return true
	placeholder
		ts, ok := value.(time.Time)
		if !ok {
			return true
	placeholder
		updates[id] = ts
		s.lastUsedUpdates.Delete(key)
		return true
placeholder)

	if len(updates) == 0 {
		return
placeholder

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.accountRepo.BatchUpdateLastUsed(ctx, updates); err != nil {
		log.Printf("[DeferredService] BatchUpdateLastUsed failed (%d accounts): %v", len(updates), err)
		for id, ts := range updates {
			s.lastUsedUpdates.Store(id, ts)
	placeholder
placeholder else {
		log.Printf("[DeferredService] BatchUpdateLastUsed flushed %d accounts", len(updates))
placeholder
placeholder
