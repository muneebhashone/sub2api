package securityaudit

import (
	"context"
	"log/slog"
	"strings"
)

const (
	EventConfigUpdated        = "prompt_audit.config_updated"
	EventConfigLoaded         = "prompt_guard.config_loaded"
	EventConfigReloadDegraded = "prompt_guard.config_reload_degraded"
	EventProbeStarted         = "prompt_audit.endpoint_probe_started"
	EventProbeFinished        = "prompt_audit.endpoint_probe_finished"
	EventProbeFailed          = "prompt_audit.endpoint_probe_failed"
	EventJobEnqueued          = "prompt_audit.job_enqueued"
	EventEnqueueSkipped       = "prompt_audit.enqueue_skipped"
	EventEnqueueDropped       = "prompt_audit.enqueue_dropped"
	EventAuditStarted         = "prompt_audit.started"
	EventProcessingReclaimed  = "prompt_audit.processing_reclaimed"
	EventProcessed            = "prompt_audit.processed"
	EventProcessFailed        = "prompt_audit.process_failed"
	EventFindingRecorded      = "prompt_audit.finding_recorded"
	EventChunkStarted         = "prompt_audit.scan_chunk_started"
	EventChunkCompleted       = "prompt_audit.scan_chunk_completed"
	EventChunkFailed          = "prompt_audit.scan_chunk_failed"
	EventChunksAggregated     = "prompt_audit.scan_chunks_aggregated"
	EventEvaluationStarted    = "prompt_guard.evaluation_started"
	EventGuardAllowed         = "prompt_guard.allowed"
	EventGuardBlocked         = "prompt_guard.blocked"
	EventGuardFailed          = "prompt_guard.failed"
	EventResultRecordFailed   = "prompt_guard.result_record_failed"
	EventEventDeleted         = "prompt_audit.event_deleted"
	EventEventsDeleted        = "prompt_audit.events_deleted"
	EventDeletePreviewed      = "prompt_audit.events_delete_previewed"
	EventEventsFilterDeleted  = "prompt_audit.events_filter_deleted"
)

var knownLogEvents = map[string]struct{placeholder{
	EventConfigUpdated: {placeholder, EventConfigLoaded: {placeholder, EventConfigReloadDegraded: {placeholder,
	EventProbeStarted: {placeholder, EventProbeFinished: {placeholder, EventProbeFailed: {placeholder,
	EventJobEnqueued: {placeholder, EventEnqueueSkipped: {placeholder, EventEnqueueDropped: {placeholder,
	EventAuditStarted: {placeholder, EventProcessingReclaimed: {placeholder, EventProcessed: {placeholder, EventProcessFailed: {placeholder, EventFindingRecorded: {placeholder,
	EventChunkStarted: {placeholder, EventChunkCompleted: {placeholder, EventChunkFailed: {placeholder, EventChunksAggregated: {placeholder,
	EventEvaluationStarted: {placeholder, EventGuardAllowed: {placeholder, EventGuardBlocked: {placeholder, EventGuardFailed: {placeholder, EventResultRecordFailed: {placeholder,
	EventEventDeleted: {placeholder, EventEventsDeleted: {placeholder, EventDeletePreviewed: {placeholder, EventEventsFilterDeleted: {placeholder,
placeholder

var allowedLogFields = map[string]struct{placeholder{
	"request_id": {placeholder, "user_id": {placeholder, "api_key_id": {placeholder, "group_id": {placeholder, "provider": {placeholder,
	"protocol": {placeholder, "endpoint": {placeholder, "model": {placeholder, "job_id": {placeholder, "event_id": {placeholder,
	"config_version": {placeholder, "guard_endpoint_id": {placeholder, "decision": {placeholder, "risk_level": {placeholder,
	"action": {placeholder, "chunk_index": {placeholder, "chunk_total": {placeholder, "chunk_chars": {placeholder, "input_chars": {placeholder,
	"input_limit": {placeholder, "latency_ms": {placeholder, "status": {placeholder, "error_code": {placeholder, "error_kind": {placeholder,
	"queue_length": {placeholder, "queue_capacity": {placeholder, "stage": {placeholder, "upstream_dispatched": {placeholder,
	"billing_preconsumed": {placeholder, "worker_id": {placeholder, "reclaimed_total": {placeholder, "attempts": {placeholder,
	"max_attempts": {placeholder, "claim_version": {placeholder, "http_status": {placeholder, "retryable": {placeholder,
placeholder

func LogInfo(event string, fields map[string]any) {
	if _, ok := knownLogEvents[event]; !ok {
		return
placeholder
	slog.LogAttrs(context.Background(), slog.LevelInfo, event, safeAttrs(fields)...)
placeholder
func LogWarn(event string, fields map[string]any) {
	if _, ok := knownLogEvents[event]; !ok {
		return
placeholder
	slog.LogAttrs(context.Background(), slog.LevelWarn, event, safeAttrs(fields)...)
placeholder
func LogError(event string, fields map[string]any) {
	if _, ok := knownLogEvents[event]; !ok {
		return
placeholder
	slog.LogAttrs(context.Background(), slog.LevelError, event, safeAttrs(fields)...)
placeholder

func safeAttrs(fields map[string]any) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(fields))
	for key, value := range fields {
		key = strings.TrimSpace(key)
		if _, allowed := allowedLogFields[key]; !allowed {
			continue
	placeholder
		if text, ok := value.(string); ok {
			if key == "error_kind" || key == "error_code" {
				value = stableErrorCode(text)
		placeholder else {
				value = TrimRunes(strings.TrimSpace(text), 256)
		placeholder
	placeholder
		attrs = append(attrs, slog.Any(key, value))
placeholder
	return attrs
placeholder

func mergeLogFields(base map[string]any, extra map[string]any) map[string]any {
	result := make(map[string]any, len(base)+len(extra))
	for key, value := range base {
		result[key] = value
placeholder
	for key, value := range extra {
		result[key] = value
placeholder
	return result
placeholder

func requestLogFields(req Request) map[string]any {
	return map[string]any{
		"request_id": req.RequestID, "user_id": req.UserID, "api_key_id": req.APIKeyID,
		"group_id": pointerLogID(req.GroupID), "provider": req.Provider, "protocol": req.Protocol,
		"endpoint": req.Endpoint, "model": req.Model, "stage": req.Stage,
placeholder
placeholder

func snapshotLogFields(snapshot PromptSnapshot) map[string]any {
	return map[string]any{
		"request_id": snapshot.RequestID, "user_id": snapshot.UserID, "api_key_id": snapshot.APIKeyID,
		"group_id": pointerLogID(snapshot.GroupID), "provider": snapshot.Provider, "protocol": snapshot.Protocol,
		"endpoint": snapshot.Endpoint, "model": snapshot.Model, "stage": snapshot.Stage,
placeholder
placeholder

func jobLogFields(job *Job) map[string]any {
	if job == nil {
		return map[string]any{placeholder
placeholder
	fields := snapshotLogFields(job.Snapshot)
	fields["job_id"] = job.ID
	fields["config_version"] = job.ConfigVersion
	fields["claim_version"] = job.ClaimVersion
	return fields
placeholder

func stableErrorCode(code string) string {
	code = strings.ToLower(strings.TrimSpace(code))
	if code == "" {
		return "unknown_error"
placeholder
	for _, char := range code {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '_' || char == '-' || char == '.' {
			continue
	placeholder
		return "redacted_error"
placeholder
	return TrimRunes(code, 64)
placeholder

func stableErrorMessage(code string) string {
	switch stableErrorCode(code) {
	case ErrorCodeBlocked:
		return "Prompt Guard blocked the request"
	case ErrorCodeUnavailable, "payload_store_unavailable", "payload_missing":
		return "Prompt Audit dependency is unavailable"
	case ErrorCodeInvalidResponse:
		return "Prompt Guard returned an invalid response"
	case "queue_full", "queue_admission_busy":
		return "Prompt Audit queue is unavailable"
	case "worker_panic":
		return "Prompt Audit worker failed"
	case "config_load_failed", "config_ttl_reload_failed", "config_invalidation_reload_failed":
		return "Prompt Audit configuration could not be loaded"
	default:
		return "Prompt Audit operation failed"
placeholder
placeholder

func sanitizeStoredError(code string) (string, string) {
	stableCode := stableErrorCode(code)
	return stableCode, TrimRunes(stableErrorMessage(stableCode), 160)
placeholder
