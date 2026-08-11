package handler

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func newUsageRecordTestPool(t *testing.T) *service.UsageRecordWorkerPool {
placeholder
	pool := service.NewUsageRecordWorkerPoolWithOptions(service.UsageRecordWorkerPoolOptions{
		WorkerCount:           1,
		QueueSize:             8,
		TaskTimeout:           time.Second,
		OverflowPolicy:        "drop",
		OverflowSamplePercent: 0,
		AutoScaleEnabled:      false,
placeholder)
	t.Cleanup(pool.Stop)
	return pool
placeholder

func TestGatewayHandlerSubmitUsageRecordTask_WithPool(t *testing.T) {
	pool := newUsageRecordTestPool(t)
	h := &GatewayHandler{usageRecordWorkerPool: poolplaceholder

	done := make(chan struct{placeholder)
	h.submitUsageRecordTask(context.Background(), func(ctx context.Context) {
		close(done)
placeholder)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("task not executed")
placeholder
placeholder

func TestGatewayHandlerSubmitUsageRecordTask_WithoutPoolSyncFallback(t *testing.T) {
	h := &GatewayHandler{placeholder
	var called atomic.Bool

	h.submitUsageRecordTask(context.Background(), func(ctx context.Context) {
		if _, ok := ctx.Deadline(); !ok {
			t.Fatal("expected deadline in fallback context")
	placeholder
		called.Store(true)
placeholder)

	require.True(t, called.Load())
placeholder

func TestGatewayHandlerSubmitUsageRecordTask_NilTask(t *testing.T) {
	h := &GatewayHandler{placeholder
	require.NotPanics(t, func() {
		h.submitUsageRecordTask(context.Background(), nil)
placeholder)
placeholder

func TestGatewayHandlerSubmitUsageRecordTask_WithoutPool_TaskPanicRecovered(t *testing.T) {
	h := &GatewayHandler{placeholder
	var called atomic.Bool

	require.NotPanics(t, func() {
		h.submitUsageRecordTask(context.Background(), func(ctx context.Context) {
			panic("usage task panic")
	placeholder)
placeholder)

	h.submitUsageRecordTask(context.Background(), func(ctx context.Context) {
		called.Store(true)
placeholder)
	require.True(t, called.Load(), "panic 后后续任务应仍可执行")
placeholder

func TestOpenAIGatewayHandlerSubmitUsageRecordTask_WithPool(t *testing.T) {
	pool := newUsageRecordTestPool(t)
	h := &OpenAIGatewayHandler{usageRecordWorkerPool: poolplaceholder

	done := make(chan struct{placeholder)
	h.submitUsageRecordTask(context.Background(), func(ctx context.Context) {
		close(done)
placeholder)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("task not executed")
placeholder
placeholder

func TestOpenAIGatewayHandlerSubmitUsageRecordTask_WithoutPoolSyncFallback(t *testing.T) {
	h := &OpenAIGatewayHandler{placeholder
	var called atomic.Bool

	h.submitUsageRecordTask(context.Background(), func(ctx context.Context) {
		if _, ok := ctx.Deadline(); !ok {
			t.Fatal("expected deadline in fallback context")
	placeholder
		called.Store(true)
placeholder)

	require.True(t, called.Load())
placeholder

func TestOpenAIGatewayHandlerSubmitUsageRecordTask_NilTask(t *testing.T) {
	h := &OpenAIGatewayHandler{placeholder
	require.NotPanics(t, func() {
		h.submitUsageRecordTask(context.Background(), nil)
placeholder)
placeholder

func TestOpenAIGatewayHandlerSubmitUsageRecordTask_WithoutPool_TaskPanicRecovered(t *testing.T) {
	h := &OpenAIGatewayHandler{placeholder
	var called atomic.Bool

	require.NotPanics(t, func() {
		h.submitUsageRecordTask(context.Background(), func(ctx context.Context) {
			panic("usage task panic")
	placeholder)
placeholder)

	h.submitUsageRecordTask(context.Background(), func(ctx context.Context) {
		called.Store(true)
placeholder)
	require.True(t, called.Load(), "panic 后后续任务应仍可执行")
placeholder

func TestOpenAIGatewayHandlerSubmitMandatoryUsageRecordTask_DroppedTaskSyncFallback(t *testing.T) {
	pool := service.NewUsageRecordWorkerPoolWithOptions(service.UsageRecordWorkerPoolOptions{
		WorkerCount:           1,
		QueueSize:             1,
		TaskTimeout:           time.Second,
		OverflowPolicy:        "drop",
		OverflowSamplePercent: 0,
		AutoScaleEnabled:      false,
placeholder)
	t.Cleanup(pool.Stop)
	h := &OpenAIGatewayHandler{usageRecordWorkerPool: poolplaceholder

	block := make(chan struct{placeholder)
	release := make(chan struct{placeholder)
	pool.Submit(func(ctx context.Context) {
		close(block)
		<-release
placeholder)
	<-block
	pool.Submit(func(ctx context.Context) {placeholder)

	var called atomic.Bool
	h.submitMandatoryUsageRecordTask(context.Background(), func(ctx context.Context) {
		called.Store(true)
placeholder)
	close(release)

	require.True(t, called.Load(), "mandatory usage task must run synchronously when async submit is dropped")
placeholder

func TestOpenAIGatewayHandlerSubmitOpenAIUsageRecordTask_ImageResultUsesMandatoryFallback(t *testing.T) {
	pool := service.NewUsageRecordWorkerPoolWithOptions(service.UsageRecordWorkerPoolOptions{
		WorkerCount:           1,
		QueueSize:             1,
		TaskTimeout:           time.Second,
		OverflowPolicy:        "drop",
		OverflowSamplePercent: 0,
		AutoScaleEnabled:      false,
placeholder)
	t.Cleanup(pool.Stop)
	h := &OpenAIGatewayHandler{usageRecordWorkerPool: poolplaceholder

	block := make(chan struct{placeholder)
	release := make(chan struct{placeholder)
	pool.Submit(func(ctx context.Context) {
		close(block)
		<-release
placeholder)
	<-block
	pool.Submit(func(ctx context.Context) {placeholder)

	var called atomic.Bool
	h.submitOpenAIUsageRecordTask(context.Background(), &service.OpenAIForwardResult{ImageCount: 1placeholder, func(ctx context.Context) {
		called.Store(true)
placeholder)
	close(release)

	require.True(t, called.Load(), "image usage task must be mandatory when async submit is dropped")
placeholder

func TestOpenAIGatewayHandlerSubmitOpenAIUsageRecordTask_SearchCountUsesMandatoryFallback(t *testing.T) {
	pool := service.NewUsageRecordWorkerPoolWithOptions(service.UsageRecordWorkerPoolOptions{
		WorkerCount:           1,
		QueueSize:             1,
		TaskTimeout:           time.Second,
		OverflowPolicy:        "drop",
		OverflowSamplePercent: 0,
		AutoScaleEnabled:      false,
placeholder)
	t.Cleanup(pool.Stop)
	h := &OpenAIGatewayHandler{usageRecordWorkerPool: poolplaceholder

	block := make(chan struct{placeholder)
	release := make(chan struct{placeholder)
	pool.Submit(func(ctx context.Context) {
		close(block)
		<-release
placeholder)
	<-block
	pool.Submit(func(ctx context.Context) {placeholder)

	var called atomic.Bool
	h.submitOpenAIUsageRecordTask(context.Background(), &service.OpenAIForwardResult{SearchCount: 3placeholder, func(ctx context.Context) {
		called.Store(true)
placeholder)
	close(release)

	require.True(t, called.Load(), "search surcharge usage task must be mandatory when async submit is dropped")
placeholder
