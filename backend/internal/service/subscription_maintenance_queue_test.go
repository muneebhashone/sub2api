//go:build unit

package service

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSubscriptionMaintenanceQueue_TryEnqueue_QueueFull(t *testing.T) {
	q := NewSubscriptionMaintenanceQueue(1, 1)
	t.Cleanup(q.Stop)

	block := make(chan struct{placeholder)
	var started atomic.Int32

	require.NoError(t, q.TryEnqueue(func() {
		started.Store(1)
		<-block
placeholder))

	// Wait until worker started consuming the first task.
	require.Eventually(t, func() bool { return started.Load() == 1 placeholder, time.Second, 10*time.Millisecond)

	// Queue size is 1; with the worker blocked, enqueueing one more should fill it.
	require.NoError(t, q.TryEnqueue(func() {placeholder))

	// Now the queue is full; next enqueue must fail.
	err := q.TryEnqueue(func() {placeholder)
placeholder
	require.Contains(t, err.Error(), "full")

	close(block)
placeholder

func TestSubscriptionMaintenanceQueue_TryEnqueue_PanicDoesNotKillWorker(t *testing.T) {
	q := NewSubscriptionMaintenanceQueue(1, 8)
	t.Cleanup(q.Stop)

	require.NoError(t, q.TryEnqueue(func() { panic("boom") placeholder))

	done := make(chan struct{placeholder)
	require.NoError(t, q.TryEnqueue(func() { close(done) placeholder))

	select {
	case <-done:
		// ok
	case <-time.After(time.Second):
		t.Fatalf("worker did not continue after panic")
placeholder
placeholder

func TestSubscriptionMaintenanceQueue_TryEnqueue_AfterStop(t *testing.T) {
	q := NewSubscriptionMaintenanceQueue(1, 8)
	q.Stop()

	err := q.TryEnqueue(func() {placeholder)
placeholder
	require.Contains(t, err.Error(), "stopped")
placeholder

func TestSubscriptionMaintenanceQueue_TryEnqueue_NilReceiver(t *testing.T) {
	var q *SubscriptionMaintenanceQueue
	err := q.TryEnqueue(func() {placeholder)
placeholder
	require.Contains(t, err.Error(), "nil")
placeholder

func TestSubscriptionMaintenanceQueue_TryEnqueue_NilTask(t *testing.T) {
	q := NewSubscriptionMaintenanceQueue(1, 8)
	t.Cleanup(q.Stop)

	err := q.TryEnqueue(nil)
placeholder
	require.Contains(t, err.Error(), "nil")
placeholder

func TestSubscriptionMaintenanceQueue_Stop_NilReceiver(t *testing.T) {
	var q *SubscriptionMaintenanceQueue
	// 不应 panic
	q.Stop()
placeholder

func TestSubscriptionMaintenanceQueue_Stop_Idempotent(t *testing.T) {
	q := NewSubscriptionMaintenanceQueue(1, 4)
	q.Stop()
	q.Stop() // 第二次调用不应 panic
placeholder

func TestNewSubscriptionMaintenanceQueue_ZeroParams(t *testing.T) {
	q := NewSubscriptionMaintenanceQueue(0, 0)
	t.Cleanup(q.Stop)

	// workerCount/queueSize 应被修正为 1
	err := q.TryEnqueue(func() {placeholder)
placeholder
placeholder

func TestNewSubscriptionMaintenanceQueue_NegativeParams(t *testing.T) {
	q := NewSubscriptionMaintenanceQueue(-1, -1)
	t.Cleanup(q.Stop)

	err := q.TryEnqueue(func() {placeholder)
placeholder
placeholder

func TestSubscriptionMaintenanceQueue_ConcurrentEnqueueAndStop(t *testing.T) {
	// 并发调用 TryEnqueue 和 Stop 不应 panic
	for i := 0; i < 100; i++ {
		q := NewSubscriptionMaintenanceQueue(2, 4)
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				_ = q.TryEnqueue(func() {placeholder)
		placeholder
	placeholder()

		go func() {
			defer wg.Done()
			time.Sleep(time.Microsecond * time.Duration(i%10))
			q.Stop()
	placeholder()

		wg.Wait()
placeholder
placeholder
