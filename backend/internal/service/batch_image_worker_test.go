//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBatchImageWorker_ProcessesJobOnce(t *testing.T) {
	queue := newFakeBatchImageQueue("imgbatch_worker_once")
	processor := &fakeBatchImageProcessor{placeholder
	worker := NewBatchImageWorker(queue, processor, BatchImageWorkerOptions{ReserveBlockTimeout: time.Millisecondplaceholder)

	require.NoError(t, worker.RunOnce(context.Background()))
	require.Equal(t, []string{"imgbatch_worker_once"placeholder, processor.processed)
	require.Len(t, queue.requeued, 1)
	require.Equal(t, defaultBatchImageWorkerRequeueDelay, queue.requeued[0].delay)
	require.Equal(t, 1, queue.releaseCount)
placeholder

func TestBatchImageWorker_RequeuesNonTerminalResultWithRequestedDelay(t *testing.T) {
	queue := newFakeBatchImageQueue("imgbatch_worker_requeue")
	processor := &fakeBatchImageProcessor{result: BatchImageProcessResult{RequeueAfter: 42 * time.Secondplaceholderplaceholder
	worker := NewBatchImageWorker(queue, processor, BatchImageWorkerOptions{placeholder)

	require.NoError(t, worker.RunOnce(context.Background()))
	require.Len(t, queue.requeued, 1)
	require.Equal(t, "imgbatch_worker_requeue", queue.requeued[0].batchID)
	require.Equal(t, 42*time.Second, queue.requeued[0].delay)
	require.Empty(t, queue.acked)
placeholder

func TestBatchImageWorker_AcksTerminalResult(t *testing.T) {
	queue := newFakeBatchImageQueue("imgbatch_worker_terminal")
	processor := &fakeBatchImageProcessor{result: BatchImageProcessResult{Terminal: trueplaceholderplaceholder
	worker := NewBatchImageWorker(queue, processor, BatchImageWorkerOptions{placeholder)

	require.NoError(t, worker.RunOnce(context.Background()))
	require.Equal(t, []string{"imgbatch_worker_terminal"placeholder, queue.acked)
	require.Empty(t, queue.requeued)
placeholder

func TestBatchImageWorker_RequeuesOnProcessorError(t *testing.T) {
	queue := newFakeBatchImageQueue("imgbatch_worker_error")
	processor := &fakeBatchImageProcessor{err: errors.New("processor failed")placeholder
	worker := NewBatchImageWorker(queue, processor, BatchImageWorkerOptions{ErrorRetryDelay: 7 * time.Secondplaceholder)

	require.NoError(t, worker.RunOnce(context.Background()))
	require.Len(t, queue.requeued, 1)
	require.Equal(t, 7*time.Second, queue.requeued[0].delay)
	require.Empty(t, queue.acked)
placeholder

func TestBatchImageWorker_SkipsWhenJobLockNotAcquired(t *testing.T) {
	queue := newFakeBatchImageQueue("imgbatch_worker_locked")
	queue.lockAcquired = false
	processor := &fakeBatchImageProcessor{placeholder
	worker := NewBatchImageWorker(queue, processor, BatchImageWorkerOptions{placeholder)

	require.NoError(t, worker.RunOnce(context.Background()))
	require.Empty(t, processor.processed)
	require.Empty(t, queue.requeued)
	require.Empty(t, queue.acked)
placeholder

func TestNewBatchImageWorkerOptionsFromConfig_UsesFiniteReserveTimeout(t *testing.T) {
	opts := NewBatchImageWorkerOptionsFromConfig(nil)
	require.Equal(t, defaultBatchImageWorkerReserveBlockTimeout, opts.ReserveBlockTimeout)
	require.Positive(t, opts.ReserveBlockTimeout)
placeholder

type fakeBatchImageQueue struct {
	reserved     ReservedBatchImageJob
	lockAcquired bool
	acked        []string
	requeued     []fakeBatchImageRequeue
	releaseCount int
placeholder

type fakeBatchImageRequeue struct {
	batchID string
	delay   time.Duration
placeholder

func newFakeBatchImageQueue(batchID string) *fakeBatchImageQueue {
	return &fakeBatchImageQueue{
		reserved:     ReservedBatchImageJob{BatchID: batchIDplaceholder,
		lockAcquired: true,
placeholder
placeholder

func (q *fakeBatchImageQueue) Enqueue(context.Context, string) error {
	return nil
placeholder

func (q *fakeBatchImageQueue) Reserve(context.Context, time.Duration) (ReservedBatchImageJob, error) {
	return q.reserved, nil
placeholder

func (q *fakeBatchImageQueue) RequeueAfter(_ context.Context, batchID string, delay time.Duration) error {
	q.requeued = append(q.requeued, fakeBatchImageRequeue{batchID: batchID, delay: delayplaceholder)
	return nil
placeholder

func (q *fakeBatchImageQueue) Ack(_ context.Context, batchID string) error {
	q.acked = append(q.acked, batchID)
	return nil
placeholder

func (q *fakeBatchImageQueue) Heartbeat(context.Context, string) error {
	return nil
placeholder

func (q *fakeBatchImageQueue) MoveDueDelayedToReady(context.Context, int) (int, error) {
	return 0, nil
placeholder

func (q *fakeBatchImageQueue) RecoverStaleActive(context.Context, time.Duration, int) (int, error) {
	return 0, nil
placeholder

func (q *fakeBatchImageQueue) TryAcquireJobLock(context.Context, string, time.Duration) (BatchImageJobLock, bool, error) {
	if !q.lockAcquired {
		return nil, false, nil
placeholder
	return fakeBatchImageLock{release: func() { q.releaseCount++ placeholderplaceholder, true, nil
placeholder

type fakeBatchImageLock struct {
	release func()
placeholder

func (l fakeBatchImageLock) Release(context.Context) error {
	if l.release != nil {
		l.release()
placeholder
	return nil
placeholder

type fakeBatchImageProcessor struct {
	result    BatchImageProcessResult
	err       error
	processed []string
placeholder

func (p *fakeBatchImageProcessor) Process(_ context.Context, batchID string) (BatchImageProcessResult, error) {
	p.processed = append(p.processed, batchID)
	return p.result, p.err
placeholder
