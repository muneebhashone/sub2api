package handler

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

// 本文件覆盖：worker 池已停止（进程关停窗口）时，计费任务不得静默丢失，
// 必须降级为内联同步执行；显式配置的 drop/sample 溢出丢弃仍按配置语义保留。

func newStoppedUsageRecordPoolForTest() *service.UsageRecordWorkerPool {
	pool := service.NewUsageRecordWorkerPoolWithOptions(service.UsageRecordWorkerPoolOptions{
		WorkerCount:    1,
		QueueSize:      1,
		TaskTimeout:    time.Second,
		OverflowPolicy: config.UsageRecordOverflowPolicySync,
placeholder)
	pool.Stop()
	return pool
placeholder

func TestGatewayHandlerSubmitUsageRecordTask_StoppedPoolFallsBackToSync(t *testing.T) {
	h := &GatewayHandler{usageRecordWorkerPool: newStoppedUsageRecordPoolForTest()placeholder

	executed := false
	h.submitUsageRecordTask(context.Background(), func(ctx context.Context) {
		executed = true
placeholder)
	require.True(t, executed, "池已停止时计费任务必须内联同步执行")
placeholder

func TestOpenAIGatewayHandlerSubmitUsageRecordTask_StoppedPoolFallsBackToSync(t *testing.T) {
	h := &OpenAIGatewayHandler{usageRecordWorkerPool: newStoppedUsageRecordPoolForTest()placeholder

	executed := false
	h.submitUsageRecordTask(context.Background(), func(ctx context.Context) {
		executed = true
placeholder)
	require.True(t, executed, "池已停止时计费任务必须内联同步执行")
placeholder

func TestGatewayHandlerSubmitUsageRecordTask_DropPolicyOverflowStillDrops(t *testing.T) {
	pool := service.NewUsageRecordWorkerPoolWithOptions(service.UsageRecordWorkerPoolOptions{
		WorkerCount:    1,
		QueueSize:      1,
		TaskTimeout:    time.Minute,
		OverflowPolicy: config.UsageRecordOverflowPolicyDrop,
placeholder)
	t.Cleanup(pool.Stop)
	h := &GatewayHandler{usageRecordWorkerPool: poolplaceholder

	started := make(chan struct{placeholder)
	block := make(chan struct{placeholder)
	t.Cleanup(func() { close(block) placeholder)
	// 占满 worker 槽位后再填满队列，保证第三个任务触发溢出。
	require.Equal(t, service.UsageRecordSubmitModeEnqueued, pool.Submit(func(ctx context.Context) {
		close(started)
		<-block
placeholder))
	<-started
	require.Equal(t, service.UsageRecordSubmitModeEnqueued, pool.Submit(func(ctx context.Context) {
		<-block
placeholder))

	var executed atomic.Bool
	h.submitUsageRecordTask(context.Background(), func(ctx context.Context) {
		executed.Store(true)
placeholder)
	time.Sleep(50 * time.Millisecond)
	require.False(t, executed.Load(), "drop 溢出策略是运维显式配置的取舍，不应被同步兜底覆盖")
placeholder
