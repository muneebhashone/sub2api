package handler

import (
	"context"
	"strings"
	"sync"
	"time"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	opsModelKey  = "ops_model"
	opsStreamKey = "ops_stream"
)

const (
	opsErrorLogWorkerCount = 10
	opsErrorLogQueueSize   = 256
	opsErrorLogTimeout     = 2 * time.Second
)

type opsErrorLogJob struct {
	ops   *service.OpsService
	entry *service.OpsErrorLog
placeholder

var (
	opsErrorLogOnce  sync.Once
	opsErrorLogQueue chan opsErrorLogJob
)

func startOpsErrorLogWorkers() {
	opsErrorLogQueue = make(chan opsErrorLogJob, opsErrorLogQueueSize)
	for i := 0; i < opsErrorLogWorkerCount; i++ {
		go func() {
			for job := range opsErrorLogQueue {
				if job.ops == nil || job.entry == nil {
					continue
			placeholder
				ctx, cancel := context.WithTimeout(context.Background(), opsErrorLogTimeout)
				_ = job.ops.RecordError(ctx, job.entry)
				cancel()
		placeholder
	placeholder()
placeholder
placeholder

func enqueueOpsErrorLog(ops *service.OpsService, entry *service.OpsErrorLog) {
	if ops == nil || entry == nil {
		return
placeholder

	opsErrorLogOnce.Do(startOpsErrorLogWorkers)

	select {
	case opsErrorLogQueue <- opsErrorLogJob{ops: ops, entry: entryplaceholder:
	default:
		// Queue is full; drop to avoid blocking request handling.
placeholder
placeholder

func setOpsRequestContext(c *gin.Context, model string, stream bool) {
	c.Set(opsModelKey, model)
	c.Set(opsStreamKey, stream)
placeholder

func recordOpsError(c *gin.Context, ops *service.OpsService, status int, errType, message, fallbackPlatform string) {
	if ops == nil || c == nil {
		return
placeholder

	model, _ := c.Get(opsModelKey)
	stream, _ := c.Get(opsStreamKey)

	var modelName string
	if m, ok := model.(string); ok {
		modelName = m
placeholder
	streaming, _ := stream.(bool)

	apiKey, _ := middleware2.GetAPIKeyFromContext(c)

	logEntry := &service.OpsErrorLog{
		Phase:      classifyOpsPhase(errType, message),
		Type:       errType,
		Severity:   classifyOpsSeverity(errType, status),
		StatusCode: status,
		Platform:   resolveOpsPlatform(apiKey, fallbackPlatform),
		Model:      modelName,
		RequestID:  c.Writer.Header().Get("x-request-id"),
		Message:    message,
		ClientIP:   c.ClientIP(),
		RequestPath: func() string {
			if c.Request != nil && c.Request.URL != nil {
				return c.Request.URL.Path
		placeholder
			return ""
	placeholder(),
		Stream: streaming,
placeholder

	if apiKey != nil {
		logEntry.APIKeyID = &apiKey.ID
		if apiKey.User != nil {
			logEntry.UserID = &apiKey.User.ID
	placeholder
		if apiKey.GroupID != nil {
			logEntry.GroupID = apiKey.GroupID
	placeholder
placeholder

	enqueueOpsErrorLog(ops, logEntry)
placeholder

func resolveOpsPlatform(apiKey *service.APIKey, fallback string) string {
	if apiKey != nil && apiKey.Group != nil && apiKey.Group.Platform != "" {
		return apiKey.Group.Platform
placeholder
	return fallback
placeholder

func classifyOpsPhase(errType, message string) string {
	msg := strings.ToLower(message)
	switch errType {
	case "authentication_error":
		return "auth"
	case "billing_error", "subscription_error":
		return "billing"
	case "rate_limit_error":
		if strings.Contains(msg, "concurrency") || strings.Contains(msg, "pending") {
			return "concurrency"
	placeholder
		return "upstream"
	case "invalid_request_error":
		return "response"
	case "upstream_error", "overloaded_error":
		return "upstream"
	case "api_error":
		if strings.Contains(msg, "no available accounts") {
			return "scheduling"
	placeholder
		return "internal"
	default:
		return "internal"
placeholder
placeholder

func classifyOpsSeverity(errType string, status int) string {
	switch errType {
	case "invalid_request_error", "authentication_error", "billing_error", "subscription_error":
		return "P3"
placeholder
	if status >= 500 {
		return "P1"
placeholder
	if status == 429 {
		return "P1"
placeholder
	if status >= 400 {
		return "P2"
placeholder
	return "P3"
placeholder
