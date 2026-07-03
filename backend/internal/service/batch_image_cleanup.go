package service

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

const (
	defaultBatchImageInputRetentionAfterTerminal  = 24 * time.Hour
	defaultBatchImageOutputRetentionAfterTerminal = 72 * time.Hour
	defaultBatchImageCleanupInterval              = 30 * time.Minute
	defaultBatchImageCleanupBatchSize             = 100
)

type BatchImageCleanupService struct {
	Repo             BatchImageRepository
	ProviderRegistry *BatchImageProviderRegistry
	AccountResolver  BatchImageAccountResolver
	Config           *config.Config

	cancel context.CancelFunc
	done   chan struct{placeholder
	mu     sync.Mutex
placeholder

func NewBatchImageCleanupService(repo BatchImageRepository, accountRepo AccountRepository, cfg *config.Config) *BatchImageCleanupService {
	return &BatchImageCleanupService{
		Repo:             repo,
		ProviderRegistry: NewDefaultBatchImageProviderRegistry(),
		AccountResolver:  &BatchImageAccountRepositoryResolver{Repo: accountRepoplaceholder,
		Config:           cfg,
placeholder
placeholder

func (s *BatchImageCleanupService) DeleteOutputsForOwner(ctx context.Context, owner BatchImageOwner, batchID string) (*BatchImagePublicBatch, error) {
	job, err := s.Repo.GetBatchImageJobByBatchIDForOwner(ctx, owner.UserID, owner.APIKeyID, batchID)
	if err != nil {
		return nil, err
placeholder
	if job.Status == BatchImageJobStatusOutputDeleted || job.OutputDeletedAt != nil {
		return BatchImageJobToPublic(job), nil
placeholder
	if job.Status != BatchImageJobStatusCompleted {
		return nil, ErrBatchImageOutputDeleteNotReady
placeholder
	_ = s.Repo.AppendBatchImageEvent(ctx, job.BatchID, "manual_output_delete_requested", map[string]any{
		"batch_id":       job.BatchID,
		"cleanup_target": "output",
		"reason":         "manual",
placeholder)
	if err := s.cleanupJob(ctx, job, CleanupTargetOutput, "manual"); err != nil {
		return nil, err
placeholder
	updated, err := s.Repo.GetBatchImageJobByBatchIDForOwner(ctx, owner.UserID, owner.APIKeyID, batchID)
	if err != nil {
		return nil, err
placeholder
	return BatchImageJobToPublic(updated), nil
placeholder

func (s *BatchImageCleanupService) CleanupInput(ctx context.Context, batchID string) error {
	job, err := s.Repo.GetBatchImageJobByBatchID(ctx, batchID)
	if err != nil {
		return err
placeholder
	return s.cleanupJob(ctx, job, CleanupTargetInput, "ttl")
placeholder

func (s *BatchImageCleanupService) CleanupOutput(ctx context.Context, batchID string, reason string) error {
	job, err := s.Repo.GetBatchImageJobByBatchID(ctx, batchID)
	if err != nil {
		return err
placeholder
	return s.cleanupJob(ctx, job, CleanupTargetOutput, reason)
placeholder

func (s *BatchImageCleanupService) RunOnce(ctx context.Context, now time.Time) (BatchImageCleanupRunResult, error) {
	if s == nil || s.Repo == nil {
		return BatchImageCleanupRunResult{placeholder, ErrBatchImageCleanupFailed
placeholder
	if now.IsZero() {
		now = time.Now()
placeholder
	limit := s.cleanupBatchSize()
	result := BatchImageCleanupRunResult{placeholder
	inputCutoff := now.Add(-s.inputRetentionAfterTerminal())
	inputJobs, err := s.Repo.ListBatchImageJobsDueForInputCleanup(ctx, inputCutoff, limit)
	if err != nil {
		return result, err
placeholder
	for _, job := range inputJobs {
		if job == nil {
			continue
	placeholder
		if err := s.cleanupJob(ctx, job, CleanupTargetInput, "ttl"); err != nil {
			result.Failures++
			continue
	placeholder
		result.InputCleaned++
placeholder
	outputJobs, err := s.Repo.ListBatchImageJobsDueForOutputCleanup(ctx, now, limit)
	if err != nil {
		return result, err
placeholder
	for _, job := range outputJobs {
		if job == nil {
			continue
	placeholder
		if err := s.cleanupJob(ctx, job, CleanupTargetOutput, "expired"); err != nil {
			result.Failures++
			continue
	placeholder
		result.OutputCleaned++
placeholder
	return result, nil
placeholder

func (s *BatchImageCleanupService) Start() {
	if s == nil || s.Repo == nil || s.Config == nil || !s.Config.BatchImage.Enabled || s.cleanupInterval() <= 0 {
		return
placeholder
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cancel != nil {
		return
placeholder
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.done = make(chan struct{placeholder)
	go func() {
		defer close(s.done)
		ticker := time.NewTicker(s.cleanupInterval())
		defer ticker.Stop()
		for {
			_, _ = s.RunOnce(ctx, time.Now())
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
		placeholder
	placeholder
placeholder()
placeholder

func (s *BatchImageCleanupService) Stop() {
	if s == nil {
		return
placeholder
	s.mu.Lock()
	cancel := s.cancel
	done := s.done
	s.cancel = nil
	s.done = nil
	s.mu.Unlock()
	if cancel != nil {
		cancel()
placeholder
	if done != nil {
		<-done
placeholder
placeholder

func (s *BatchImageCleanupService) cleanupJob(ctx context.Context, job *BatchImageJob, target CleanupTarget, reason string) error {
	if job == nil {
		return ErrBatchImageJobNotFound
placeholder
	switch target {
	case CleanupTargetInput:
		if job.InputDeletedAt != nil {
			return nil
	placeholder
		if !IsTerminalBatchImageJobStatus(job.Status) {
			return ErrBatchImageCleanupFailed
	placeholder
		_ = s.Repo.AppendBatchImageEvent(ctx, job.BatchID, "input_cleanup_started", cleanupEventPayload(job.BatchID, target, reason, nil))
	case CleanupTargetOutput:
		if job.OutputDeletedAt != nil || job.Status == BatchImageJobStatusOutputDeleted {
			return nil
	placeholder
		if job.Status != BatchImageJobStatusCompleted && job.Status != BatchImageJobStatusFailed && job.Status != BatchImageJobStatusCancelled {
			return ErrBatchImageOutputDeleteNotReady
	placeholder
		_ = s.Repo.AppendBatchImageEvent(ctx, job.BatchID, "output_cleanup_started", cleanupEventPayload(job.BatchID, target, reason, nil))
	default:
		return ErrUnsupportedCleanupTarget
placeholder

	if err := s.callProviderCleanup(ctx, job, target); err != nil {
		code := cleanupFailureCode(err)
		msg := sanitizeBatchImagePublicMessage(err.Error())
		_ = s.Repo.RecordBatchImageCleanupFailure(ctx, job.BatchID, code, msg)
		event := string(target) + "_cleanup_failed"
		_ = s.Repo.AppendBatchImageEvent(ctx, job.BatchID, event, map[string]any{"batch_id": job.BatchID, "cleanup_target": string(target), "reason": reason, "error_code": codeplaceholder)
		if errors.Is(err, ErrBatchImageProviderUnsafeCleanupPath) {
			return ErrBatchImageCleanupUnsafePath
	placeholder
		return ErrBatchImageProviderCleanupFailed
placeholder

	deletedAt := time.Now()
	if target == CleanupTargetInput {
		return s.Repo.MarkBatchImageInputDeleted(ctx, job.BatchID, deletedAt)
placeholder
	return s.Repo.MarkBatchImageOutputDeleted(ctx, job.BatchID, deletedAt)
placeholder

func (s *BatchImageCleanupService) callProviderCleanup(ctx context.Context, job *BatchImageJob, target CleanupTarget) error {
	if s == nil || s.ProviderRegistry == nil || s.AccountResolver == nil {
		return ErrBatchImageCleanupFailed
placeholder
	provider, ok := s.ProviderRegistry.Get(job.Provider)
	if !ok || provider == nil {
		return ErrBatchImageUnsupportedProvider
placeholder
	if job.AccountID == nil || *job.AccountID <= 0 {
		return ErrBatchImageMissingAccountID
placeholder
	account, err := s.AccountResolver.ResolveBatchImageAccount(ctx, *job.AccountID)
	if err != nil {
		return err
placeholder
	if err := provider.Cleanup(ctx, job, account, target); err != nil {
		if cleanupErrorIsNotFound(err) {
			return nil
	placeholder
		return err
placeholder
	return nil
placeholder

func (s *BatchImageCleanupService) inputRetentionAfterTerminal() time.Duration {
	if s != nil && s.Config != nil && s.Config.BatchImage.InputRetentionAfterTerminalHours > 0 {
		return time.Duration(s.Config.BatchImage.InputRetentionAfterTerminalHours) * time.Hour
placeholder
	return defaultBatchImageInputRetentionAfterTerminal
placeholder

func (s *BatchImageCleanupService) cleanupInterval() time.Duration {
	if s != nil && s.Config != nil && s.Config.BatchImage.CleanupIntervalMinutes > 0 {
		return time.Duration(s.Config.BatchImage.CleanupIntervalMinutes) * time.Minute
placeholder
	return defaultBatchImageCleanupInterval
placeholder

func (s *BatchImageCleanupService) cleanupBatchSize() int {
	if s != nil && s.Config != nil && s.Config.BatchImage.CleanupBatchSize > 0 {
		return s.Config.BatchImage.CleanupBatchSize
placeholder
	return defaultBatchImageCleanupBatchSize
placeholder

type BatchImageCleanupRunResult struct {
	InputCleaned  int
	OutputCleaned int
	Failures      int
placeholder

func cleanupEventPayload(batchID string, target CleanupTarget, reason string, deletedAt *time.Time) map[string]any {
	payload := map[string]any{
		"batch_id":       batchID,
		"cleanup_target": string(target),
		"reason":         reason,
placeholder
	if deletedAt != nil {
		payload["deleted_at"] = deletedAt.UTC().Format(time.RFC3339)
placeholder
	return payload
placeholder

func cleanupErrorIsNotFound(err error) bool {
	if err == nil {
		return false
placeholder
	reason := strings.ToUpper(infraerrors.Reason(err))
	msg := strings.ToUpper(err.Error())
	return strings.Contains(reason, "NOT_FOUND") || strings.Contains(msg, "NOT FOUND") || strings.Contains(msg, "404")
placeholder

func cleanupFailureCode(err error) string {
	if errors.Is(err, ErrBatchImageProviderUnsafeCleanupPath) {
		return "BATCH_IMAGE_CLEANUP_UNSAFE_PATH"
placeholder
	reason := strings.TrimSpace(infraerrors.Reason(err))
	if reason != "" {
		return reason
placeholder
	return "BATCH_IMAGE_PROVIDER_CLEANUP_FAILED"
placeholder
