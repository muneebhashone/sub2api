package middleware

import (
	"context"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

const (
	opsAuthErrorLogWorkerCount = 10
	opsAuthErrorLogQueueSize   = 256
	opsAuthErrorLogTimeout     = 2 * time.Second
)

type opsAuthErrorLogJob struct {
	ops   *service.OpsService
	entry *service.OpsErrorLog
placeholder

var (
	opsAuthErrorLogOnce  sync.Once
	opsAuthErrorLogQueue chan opsAuthErrorLogJob
)

func startOpsAuthErrorLogWorkers() {
	opsAuthErrorLogQueue = make(chan opsAuthErrorLogJob, opsAuthErrorLogQueueSize)
	for i := 0; i < opsAuthErrorLogWorkerCount; i++ {
		go func() {
			for job := range opsAuthErrorLogQueue {
				if job.ops == nil || job.entry == nil {
					continue
			placeholder
				ctx, cancel := context.WithTimeout(context.Background(), opsAuthErrorLogTimeout)
				_ = job.ops.RecordError(ctx, job.entry)
				cancel()
		placeholder
	placeholder()
placeholder
placeholder

func enqueueOpsAuthErrorLog(ops *service.OpsService, entry *service.OpsErrorLog) {
	if ops == nil || entry == nil {
		return
placeholder

	opsAuthErrorLogOnce.Do(startOpsAuthErrorLogWorkers)

	select {
	case opsAuthErrorLogQueue <- opsAuthErrorLogJob{ops: ops, entry: entryplaceholder:
	default:
		// Queue is full; drop to avoid blocking request handling.
placeholder
placeholder
