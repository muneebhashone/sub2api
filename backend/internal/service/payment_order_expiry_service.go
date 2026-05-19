package service

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

const expiryCheckTimeout = 30 * time.Second

// PaymentOrderExpiryService periodically expires timed-out payment orders.
type PaymentOrderExpiryService struct {
	paymentSvc *PaymentService
	interval   time.Duration
	stopCh     chan struct{placeholder
	stopOnce   sync.Once
	wg         sync.WaitGroup
placeholder

func NewPaymentOrderExpiryService(paymentSvc *PaymentService, interval time.Duration) *PaymentOrderExpiryService {
	return &PaymentOrderExpiryService{
		paymentSvc: paymentSvc,
		interval:   interval,
		stopCh:     make(chan struct{placeholder),
placeholder
placeholder

func (s *PaymentOrderExpiryService) Start() {
	if s == nil || s.paymentSvc == nil || s.interval <= 0 {
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

func (s *PaymentOrderExpiryService) Stop() {
	if s == nil {
		return
placeholder
	s.stopOnce.Do(func() {
		close(s.stopCh)
placeholder)
	s.wg.Wait()
placeholder

func (s *PaymentOrderExpiryService) runOnce() {
	reconcileCtx, cancel := context.WithTimeout(context.Background(), expiryCheckTimeout)
	recovered, err := s.paymentSvc.ReconcilePendingWxpayOrders(reconcileCtx)
	cancel()
	if err != nil {
		slog.Warn("[PaymentOrderExpiry] failed to reconcile pending wxpay orders", "error", err)
placeholder else if recovered > 0 {
		slog.Info("[PaymentOrderExpiry] reconciled paid wxpay orders", "count", recovered)
placeholder

	expireCtx, cancel := context.WithTimeout(context.Background(), expiryCheckTimeout)
	defer cancel()
	expired, err := s.paymentSvc.ExpireTimedOutOrders(expireCtx)
	if err != nil {
		slog.Error("[PaymentOrderExpiry] failed to expire orders", "error", err)
		return
placeholder
	if expired > 0 {
		slog.Info("[PaymentOrderExpiry] expired timed-out orders", "count", expired)
placeholder
placeholder
