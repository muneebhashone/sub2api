package securityaudit

import (
	"context"
	"errors"
)

type Enqueuer struct {
	config  ConfigStore
	repo    JobRepository
	payload PayloadStore
	metrics Metrics
placeholder

func NewEnqueuer(config ConfigStore, repo JobRepository, payload PayloadStore, metrics ...Metrics) *Enqueuer {
	var metric Metrics
	if len(metrics) > 0 {
		metric = metrics[0]
placeholder
	return &Enqueuer{config: config, repo: repo, payload: payload, metrics: metricplaceholder
placeholder

func (e *Enqueuer) Enqueue(ctx context.Context, req Request) error {
	if e == nil || e.config == nil || e.repo == nil || e.payload == nil {
		return errors.New("prompt audit enqueuer unavailable")
placeholder
	cfg, ok := e.config.Active()
	baseFields := requestLogFields(req)
	if !ok || cfg.EffectiveMode() != ModeAsync {
		LogInfo(EventEnqueueSkipped, mergeLogFields(baseFields, map[string]any{"status": "skipped", "error_code": "mode_not_async"placeholder))
		return nil
placeholder
	baseFields["config_version"] = cfg.ConfigVersion
	if !cfg.IncludesGroup(req.GroupID) {
		LogInfo(EventEnqueueSkipped, mergeLogFields(baseFields, map[string]any{"status": "skipped", "error_code": "group_out_of_scope"placeholder))
		return nil
placeholder
	if len(cfg.EnabledEndpoints()) == 0 {
		e.recordDropped()
		LogWarn(EventEnqueueDropped, mergeLogFields(baseFields, map[string]any{"status": "dropped", "error_code": "no_enabled_endpoint"placeholder))
		return nil
placeholder
	snapshot, err := ExtractPromptSnapshot(req)
	if errors.Is(err, ErrNoPromptText) {
		LogInfo(EventEnqueueSkipped, mergeLogFields(baseFields, map[string]any{"status": "skipped", "error_code": "no_user_text"placeholder))
		return nil
placeholder
	if err != nil {
		e.recordDropped()
		LogWarn(EventEnqueueDropped, mergeLogFields(baseFields, map[string]any{"status": "dropped", "error_code": "snapshot_invalid"placeholder))
		return nil
placeholder
	job, err := e.repo.CreateStagingWithCapacity(ctx, snapshot.Redacted(), cfg.ConfigVersion, 3, cfg.QueueCapacity)
	if err != nil {
		code := "database_unavailable"
		if errors.Is(err, ErrQueueFull) {
			code = "queue_full"
	placeholder
		if errors.Is(err, ErrQueueAdmissionBusy) {
			code = "queue_admission_busy"
	placeholder
		LogWarn(EventEnqueueDropped, mergeLogFields(baseFields, map[string]any{
			"queue_capacity": cfg.QueueCapacity, "status": "dropped", "error_code": code,
	placeholder))
		e.recordDropped()
		return err
placeholder
	if err := e.payload.Set(ctx, job.ID, snapshot.ScanText, DefaultPayloadTTL); err != nil {
		_ = e.repo.MarkStagingFailed(ctx, job.ID, "payload_store_failed", "payload store unavailable")
		LogWarn(EventEnqueueDropped, mergeLogFields(baseFields, map[string]any{
			"job_id": job.ID, "status": "dropped", "error_code": "payload_store_failed",
	placeholder))
		e.recordDropped()
		return err
placeholder
	if err := e.repo.PublishQueued(ctx, job.ID); err != nil {
		_ = e.payload.Delete(ctx, job.ID)
		_ = e.repo.MarkStagingFailed(ctx, job.ID, "queue_publish_failed", "queue publish failed")
		LogWarn(EventEnqueueDropped, mergeLogFields(baseFields, map[string]any{
			"job_id": job.ID, "status": "dropped", "error_code": "queue_publish_failed",
	placeholder))
		e.recordDropped()
		return err
placeholder
	LogInfo(EventJobEnqueued, mergeLogFields(baseFields, map[string]any{
		"job_id":         job.ID,
		"queue_capacity": cfg.QueueCapacity, "status": "queued",
placeholder))
	if e.metrics != nil {
		e.metrics.IncEnqueued()
placeholder
	return nil
placeholder

func (e *Enqueuer) recordDropped() {
	if e != nil && e.metrics != nil {
		e.metrics.IncDropped()
placeholder
placeholder
