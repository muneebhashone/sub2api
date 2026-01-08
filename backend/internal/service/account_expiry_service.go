package service

import (
	"context"
	"log"
	"sync"
	"time"
)

// AccountExpiryService periodically pauses expired accounts when auto-pause is enabled.
type AccountExpiryService struct {
	accountRepo AccountRepository
	interval    time.Duration
	stopCh      chan struct{placeholder
	stopOnce    sync.Once
	wg          sync.WaitGroup
placeholder

func NewAccountExpiryService(accountRepo AccountRepository, interval time.Duration) *AccountExpiryService {
	return &AccountExpiryService{
		accountRepo: accountRepo,
		interval:    interval,
		stopCh:      make(chan struct{placeholder),
placeholder
placeholder

func (s *AccountExpiryService) Start() {
	if s == nil || s.accountRepo == nil || s.interval <= 0 {
		return
placeholder
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		s.runOnce()
		for {
			select {
			case <-ticker.C:
				s.runOnce()
			case <-s.stopCh:
				return
		placeholder
	placeholder
placeholder()
placeholder

func (s *AccountExpiryService) Stop() {
	if s == nil {
		return
placeholder
	s.stopOnce.Do(func() {
		close(s.stopCh)
placeholder)
	s.wg.Wait()
placeholder

func (s *AccountExpiryService) runOnce() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updated, err := s.accountRepo.AutoPauseExpiredAccounts(ctx, time.Now())
	if err != nil {
		log.Printf("[AccountExpiry] Auto pause expired accounts failed: %v", err)
		return
placeholder
	if updated > 0 {
		log.Printf("[AccountExpiry] Auto paused %d expired accounts", updated)
placeholder
placeholder
