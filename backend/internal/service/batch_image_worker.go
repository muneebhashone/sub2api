package service

import (
	"context"
	"errors"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"go.uber.org/zap"
)

const (
	defaultBatchImageWorkerLockTTL             = 5 * time.Minute
	defaultBatchImageWorkerLockConflictDelay   = 5 * time.Second
	defaultBatchImageWorkerErrorRetryDelay     = time.Minute
	defaultBatchImageWorkerRequeueDelay        = 30 * time.Second
	defaultBatchImageWorkerDelayedPollInterval = 5 * time.Second
	defaultBatchImageWorkerRecoveryInterval    = 5 * time.Minute
	defaultBatchImageWorkerStaleActiveAfter    = 10 * time.Minute
	defaultBatchImageWorkerDelayedMoveLimit    = 100
	defaultBatchImageWorkerRecoverLimit        = 100
	defaultBatchImageWorkerErrorBackoff        = time.Second
	defaultBatchImageWorkerReserveBlockTimeout = 5 * time.Second
)

type BatchImageProcessor interface {
	Process(ctx context.Context, batchID string) (BatchImageProcessResult, error)
placeholder

type BatchImageProcessResult struct {
	RequeueAfter time.Duration
	Terminal     bool
placeholder

type BatchImageWorkerOptions struct {
	ReserveBlockTimeout time.Duration
	JobLockTTL          time.Duration
	LockConflictDelay   time.Duration
	DefaultRequeueDelay time.Duration
	ErrorRetryDelay     time.Duration
	ErrorBackoff        time.Duration
	DelayedPollInterval time.Duration
	RecoveryInterval    time.Duration
	StaleActiveAfter    time.Duration
	DelayedMoveLimit    int
	RecoverLimit        int
placeholder

type BatchImageWorker struct {
	queue     BatchImageQueue
	processor BatchImageProcessor
	opts      BatchImageWorkerOptions
placeholder

func NewBatchImageWorker(queue BatchImageQueue, processor BatchImageProcessor, opts BatchImageWorkerOptions) *BatchImageWorker {
	return &BatchImageWorker{
		queue:     queue,
		processor: processor,
		opts:      normalizeBatchImageWorkerOptions(opts),
placeholder
placeholder

func NewBatchImageWorkerOptionsFromConfig(cfg *config.Config) BatchImageWorkerOptions {
	if cfg == nil {
		return normalizeBatchImageWorkerOptions(BatchImageWorkerOptions{placeholder)
placeholder
	return normalizeBatchImageWorkerOptions(BatchImageWorkerOptions{
		JobLockTTL:          time.Duration(cfg.BatchImage.JobLockTTLSeconds) * time.Second,
		LockConflictDelay:   time.Duration(cfg.BatchImage.LockConflictDelaySeconds) * time.Second,
		DefaultRequeueDelay: time.Duration(cfg.BatchImage.DefaultRequeueDelaySeconds) * time.Second,
		ErrorRetryDelay:     time.Duration(cfg.BatchImage.ErrorRetryDelaySeconds) * time.Second,
		DelayedPollInterval: time.Duration(cfg.BatchImage.DelayedMoverIntervalSeconds) * time.Second,
		RecoveryInterval:    time.Duration(cfg.BatchImage.RecoveryIntervalSeconds) * time.Second,
		StaleActiveAfter:    time.Duration(cfg.BatchImage.StaleActiveAfterSeconds) * time.Second,
		DelayedMoveLimit:    cfg.BatchImage.DelayedMoveLimit,
		RecoverLimit:        cfg.BatchImage.RecoverLimit,
placeholder)
placeholder

func normalizeBatchImageWorkerOptions(opts BatchImageWorkerOptions) BatchImageWorkerOptions {
	if opts.ReserveBlockTimeout <= 0 {
		opts.ReserveBlockTimeout = defaultBatchImageWorkerReserveBlockTimeout
placeholder
	if opts.JobLockTTL <= 0 {
		opts.JobLockTTL = defaultBatchImageWorkerLockTTL
placeholder
	if opts.LockConflictDelay <= 0 {
		opts.LockConflictDelay = defaultBatchImageWorkerLockConflictDelay
placeholder
	if opts.DefaultRequeueDelay <= 0 {
		opts.DefaultRequeueDelay = defaultBatchImageWorkerRequeueDelay
placeholder
	if opts.ErrorRetryDelay <= 0 {
		opts.ErrorRetryDelay = defaultBatchImageWorkerErrorRetryDelay
placeholder
	if opts.ErrorBackoff <= 0 {
		opts.ErrorBackoff = defaultBatchImageWorkerErrorBackoff
placeholder
	if opts.DelayedPollInterval <= 0 {
		opts.DelayedPollInterval = defaultBatchImageWorkerDelayedPollInterval
placeholder
	if opts.RecoveryInterval <= 0 {
		opts.RecoveryInterval = defaultBatchImageWorkerRecoveryInterval
placeholder
	if opts.StaleActiveAfter <= 0 {
		opts.StaleActiveAfter = defaultBatchImageWorkerStaleActiveAfter
placeholder
	if opts.DelayedMoveLimit <= 0 {
		opts.DelayedMoveLimit = defaultBatchImageWorkerDelayedMoveLimit
placeholder
	if opts.RecoverLimit <= 0 {
		opts.RecoverLimit = defaultBatchImageWorkerRecoverLimit
placeholder
	return opts
placeholder

func (w *BatchImageWorker) Run(ctx context.Context) {
	if w == nil {
		return
placeholder
	for {
		if err := ctx.Err(); err != nil {
			return
	placeholder
		if err := w.RunOnce(ctx); err != nil && ctx.Err() == nil {
			sleepOrDone(ctx, w.opts.ErrorBackoff)
	placeholder
placeholder
placeholder

func (w *BatchImageWorker) RunOnce(ctx context.Context) error {
	if w == nil || w.queue == nil || w.processor == nil {
		return nil
placeholder

	reserved, err := w.queue.Reserve(ctx, w.opts.ReserveBlockTimeout)
	if errors.Is(err, ErrBatchImageQueueEmpty) {
		return nil
placeholder
	if err != nil {
		return err
placeholder

	lock, ok, err := w.queue.TryAcquireJobLock(ctx, reserved.BatchID, w.opts.JobLockTTL)
	if err != nil {
		if requeueErr := w.queue.RequeueAfter(ctx, reserved.BatchID, w.opts.LockConflictDelay); requeueErr != nil {
			return requeueErr
	placeholder
		return err
placeholder
	if !ok {
		return nil
placeholder
	defer func() {
		_ = lock.Release(ctx)
placeholder()

	result, err := w.processor.Process(ctx, reserved.BatchID)
	if err != nil {
		logger.L().Warn("batch_image.worker_process_failed",
			zap.String("batch_id", reserved.BatchID),
			zap.Error(err),
		)
		return w.queue.RequeueAfter(ctx, reserved.BatchID, w.opts.ErrorRetryDelay)
placeholder
	if result.Terminal {
		return w.queue.Ack(ctx, reserved.BatchID)
placeholder
	delay := result.RequeueAfter
	if delay <= 0 {
		delay = w.opts.DefaultRequeueDelay
placeholder
	return w.queue.RequeueAfter(ctx, reserved.BatchID, delay)
placeholder

func (w *BatchImageWorker) MoveDueDelayedOnce(ctx context.Context) (int, error) {
	if w == nil || w.queue == nil {
		return 0, nil
placeholder
	return w.queue.MoveDueDelayedToReady(ctx, w.opts.DelayedMoveLimit)
placeholder

func (w *BatchImageWorker) RunDelayedMover(ctx context.Context) {
	if w == nil {
		return
placeholder
	for {
		if err := ctx.Err(); err != nil {
			return
	placeholder
		moved, _ := w.MoveDueDelayedOnce(ctx)
		if moved > 0 {
			continue
	placeholder
		sleepOrDone(ctx, w.opts.DelayedPollInterval)
placeholder
placeholder

func (w *BatchImageWorker) RecoverStaleActiveOnce(ctx context.Context) (int, error) {
	if w == nil || w.queue == nil {
		return 0, nil
placeholder
	return w.queue.RecoverStaleActive(ctx, w.opts.StaleActiveAfter, w.opts.RecoverLimit)
placeholder

func (w *BatchImageWorker) RunStaleActiveRecovery(ctx context.Context) {
	if w == nil {
		return
placeholder
	for {
		if err := ctx.Err(); err != nil {
			return
	placeholder
		_, _ = w.RecoverStaleActiveOnce(ctx)
		sleepOrDone(ctx, w.opts.RecoveryInterval)
placeholder
placeholder

func sleepOrDone(ctx context.Context, d time.Duration) {
	if d <= 0 {
		return
placeholder
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
	case <-timer.C:
placeholder
placeholder
