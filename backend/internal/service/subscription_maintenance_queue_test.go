//go:build unit

package service

import (
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
