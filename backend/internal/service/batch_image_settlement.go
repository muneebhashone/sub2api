package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
)

const (
	batchImageSettlementRequestPrefix = "batch_image_settlement:"
	batchImageSettlementRetryDelay    = time.Minute
)

type BatchImagePricingResolver interface {
	BatchImageUnitPrice(ctx context.Context, job *BatchImageJob) (float64, error)
placeholder

type BatchImageModelPricingResolver struct {
	Resolver *ModelPricingResolver
placeholder

func (r *BatchImageModelPricingResolver) BatchImageUnitPrice(ctx context.Context, job *BatchImageJob) (float64, error) {
	if r == nil || r.Resolver == nil || job == nil || strings.TrimSpace(job.Model) == "" {
		return 0, ErrBatchImageSettlementPricingMissing
placeholder
	resolved := r.Resolver.Resolve(ctx, PricingInput{Model: job.Modelplaceholder)
	if resolved == nil {
		return 0, ErrBatchImageSettlementPricingMissing
placeholder
	switch resolved.Mode {
	case BillingModeImage, BillingModePerRequest:
		if resolved.DefaultPerRequestPrice > 0 {
			return resolved.DefaultPerRequestPrice, nil
	placeholder
		if len(resolved.RequestTiers) == 1 && resolved.RequestTiers[0].PerRequestPrice != nil && *resolved.RequestTiers[0].PerRequestPrice >= 0 {
			return *resolved.RequestTiers[0].PerRequestPrice, nil
	placeholder
	case BillingModeToken:
		if resolved.BasePricing != nil && (resolved.BasePricing.ImageOutputPriceExplicit || resolved.BasePricing.ImageOutputPricePerToken > 0) {
			return resolved.BasePricing.ImageOutputPricePerToken, nil
	placeholder
placeholder
	return 0, ErrBatchImageSettlementPricingMissing
placeholder

type BatchImageSettlementService struct {
	Repo        BatchImageRepository
	BillingRepo UsageBillingRepository
	Pricing     BatchImagePricingResolver
	Config      *config.Config
placeholder

type BatchImageSettlementResult struct {
	BatchID        string
	SuccessCount   int
	FailCount      int
	ActualCost     float64
	ManifestHash   string
	RequestID      string
	AlreadySettled bool
placeholder

func (s *BatchImageSettlementService) Settle(ctx context.Context, batchID string) (*BatchImageSettlementResult, error) {
	if s == nil || s.Repo == nil || s.BillingRepo == nil || s.Pricing == nil {
		return nil, ErrBatchImageSettlementBillingFailed.WithCause(errors.New("batch image settlement service is not configured"))
placeholder
	job, err := s.Repo.GetBatchImageJobByBatchID(ctx, batchID)
	if err != nil {
		return nil, err
placeholder

	manifestHash := BuildBatchImageSettlementManifestHash(job)
	result := &BatchImageSettlementResult{
		BatchID:      job.BatchID,
		SuccessCount: job.SuccessCount,
		FailCount:    job.FailCount,
		ManifestHash: manifestHash,
		RequestID:    BatchImageSettlementRequestID(job.BatchID),
placeholder
	if job.ActualCost != nil {
		result.ActualCost = *job.ActualCost
placeholder
	if job.Status == BatchImageJobStatusCompleted {
		result.AlreadySettled = true
		return result, nil
placeholder
	if job.Status != BatchImageJobStatusSettling {
		return nil, ErrBatchImageSettlementInvalidStatus
placeholder
	if job.SuccessCount < 0 || job.FailCount < 0 || job.ItemCount < 0 {
		return nil, ErrBatchImageSettlementInvalidCounts
placeholder
	if strings.TrimSpace(batchImageDerefString(job.ManifestHash)) != "" && batchImageDerefString(job.ManifestHash) != manifestHash {
		return nil, ErrBatchImageSettlementManifestConflict
placeholder
	if job.APIKeyID == nil || *job.APIKeyID <= 0 {
		return nil, ErrBatchImageSettlementMissingAPIKeyID
placeholder
	if job.AccountID == nil || *job.AccountID <= 0 {
		return nil, ErrBatchImageSettlementMissingAccountID
placeholder

	unitPrice, err := s.Pricing.BatchImageUnitPrice(ctx, job)
	if err != nil {
		return nil, err
placeholder
	if unitPrice < 0 {
		return nil, ErrBatchImageSettlementPricingMissing
placeholder
	actualCost := float64(job.SuccessCount) * unitPrice
	result.ActualCost = actualCost

	cmd := &UsageBillingCommand{
		RequestID:          result.RequestID,
		APIKeyID:           *job.APIKeyID,
		RequestPayloadHash: manifestHash,
		UserID:             job.UserID,
		AccountID:          *job.AccountID,
		Model:              job.Model,
		BillingType:        BillingTypeBalance,
		ImageCount:         job.SuccessCount,
		MediaType:          "image",
		BalanceCost:        actualCost,
placeholder
	if _, err := s.BillingRepo.Apply(ctx, cmd); err != nil {
		msg := truncateBatchImageMessage(err.Error(), batchImageMaxErrorMessageLength)
		_ = s.Repo.SetBatchImageJobSettlementFailed(ctx, job.BatchID, "SETTLEMENT_BILLING_FAILED", msg)
		return nil, ErrBatchImageSettlementBillingFailed.WithCause(err)
placeholder

	now := time.Now()
	outputExpiresAt := now.Add(s.outputRetentionAfterTerminal())
	if err := s.Repo.MarkBatchImageJobSettled(ctx, MarkBatchImageJobSettledParams{
		BatchID:         job.BatchID,
		ActualCost:      actualCost,
		ManifestHash:    manifestHash,
		Now:             &now,
		OutputExpiresAt: &outputExpiresAt,
		EventPayload: map[string]any{
			"batch_id":      job.BatchID,
			"request_id":    result.RequestID,
			"success_count": job.SuccessCount,
			"fail_count":    job.FailCount,
			"actual_cost":   actualCost,
			"manifest_hash": manifestHash,
	placeholder,
placeholder); err != nil {
		return nil, err
placeholder

	return result, nil
placeholder

func (s *BatchImageSettlementService) outputRetentionAfterTerminal() time.Duration {
	if s != nil && s.Config != nil && s.Config.BatchImage.OutputRetentionAfterTerminalHours > 0 {
		return time.Duration(s.Config.BatchImage.OutputRetentionAfterTerminalHours) * time.Hour
placeholder
	return 72 * time.Hour
placeholder

func BatchImageSettlementRequestID(batchID string) string {
	return batchImageSettlementRequestPrefix + strings.TrimSpace(batchID)
placeholder

func BuildBatchImageSettlementManifestHash(job *BatchImageJob) string {
	if job == nil {
		return ""
placeholder
	parts := []string{
		strings.TrimSpace(job.BatchID),
		strings.TrimSpace(job.Provider),
		strings.TrimSpace(job.Model),
		batchImageDerefString(job.ProviderJobName),
		batchImageDerefString(job.ProviderOutputRef),
		strconv.Itoa(job.SuccessCount),
		strconv.Itoa(job.FailCount),
		strconv.Itoa(job.ItemCount),
placeholder
	sum := sha256.Sum256([]byte(strings.Join(parts, "\x00")))
	return hex.EncodeToString(sum[:])
placeholder

type BatchImagePipelineProcessor struct {
	ProviderProcessor *BatchImageProviderProcessor
	SettlementService *BatchImageSettlementService
	RetryDelay        time.Duration
placeholder

func (p *BatchImagePipelineProcessor) Process(ctx context.Context, batchID string) (BatchImageProcessResult, error) {
	if p == nil || p.ProviderProcessor == nil {
		return BatchImageProcessResult{placeholder, errors.New("batch image pipeline processor is not configured")
placeholder
	job, err := p.ProviderProcessor.Repo.GetBatchImageJobByBatchID(ctx, batchID)
	if err != nil {
		return BatchImageProcessResult{placeholder, err
placeholder
	if job.Status == BatchImageJobStatusSettling {
		if p.SettlementService == nil {
			return BatchImageProcessResult{Terminal: trueplaceholder, nil
	placeholder
		_, err := p.SettlementService.Settle(ctx, batchID)
		if err != nil {
			if errors.Is(err, ErrBatchImageSettlementBillingFailed) {
				delay := p.RetryDelay
				if delay <= 0 {
					delay = batchImageSettlementRetryDelay
			placeholder
				return BatchImageProcessResult{RequeueAfter: delayplaceholder, nil
		placeholder
			return BatchImageProcessResult{placeholder, err
	placeholder
		return BatchImageProcessResult{Terminal: trueplaceholder, nil
placeholder
	return p.ProviderProcessor.Process(ctx, batchID)
placeholder

func (r *BatchImageSettlementResult) String() string {
	if r == nil {
		return ""
placeholder
	return fmt.Sprintf("batch_id=%s success=%d fail=%d actual_cost=%0.10f already_settled=%t",
		r.BatchID, r.SuccessCount, r.FailCount, r.ActualCost, r.AlreadySettled)
placeholder
