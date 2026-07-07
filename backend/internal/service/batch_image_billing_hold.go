package service

import (
	"context"
	"errors"
	"strings"
)

const (
	batchImageHoldRequestPrefix    = "batch_image_hold:"
	batchImageCaptureRequestPrefix = "batch_image_capture:"
	batchImageReleaseRequestPrefix = "batch_image_release:"
)

func BatchImageHoldRequestID(batchID string) string {
	return batchImageHoldRequestPrefix + strings.TrimSpace(batchID)
placeholder

func BatchImageCaptureRequestID(batchID string) string {
	return batchImageCaptureRequestPrefix + strings.TrimSpace(batchID)
placeholder

func BatchImageReleaseRequestID(batchID string) string {
	return batchImageReleaseRequestPrefix + strings.TrimSpace(batchID)
placeholder

func buildBatchImageHoldCommand(job *BatchImageJob, requestID string, actualAmount float64, payloadHash string) (*BatchImageBalanceHoldCommand, error) {
	if job == nil {
		return nil, ErrBatchImageBillingHoldFailed
placeholder
	if job.APIKeyID == nil || *job.APIKeyID <= 0 {
		return nil, ErrBatchImageSettlementMissingAPIKeyID
placeholder
	holdAmount := job.EstimatedCost
	if job.HoldAmount != nil {
		holdAmount = *job.HoldAmount
placeholder
	if holdAmount < 0 {
		holdAmount = 0
placeholder
	if actualAmount < 0 {
		actualAmount = 0
placeholder
	return &BatchImageBalanceHoldCommand{
		RequestID:          requestID,
		APIKeyID:           *job.APIKeyID,
		UserID:             job.UserID,
		BatchID:            job.BatchID,
		HoldAmount:         holdAmount,
		ActualAmount:       actualAmount,
		RequestPayloadHash: strings.TrimSpace(payloadHash),
placeholder, nil
placeholder

func reserveBatchImageBalanceHold(ctx context.Context, repo UsageBillingRepository, job *BatchImageJob, payloadHash string) error {
	if repo == nil {
		return ErrBatchImageBillingHoldFailed.WithCause(errors.New("batch image billing repository is not configured"))
placeholder
	cmd, err := buildBatchImageHoldCommand(job, BatchImageHoldRequestID(job.BatchID), 0, payloadHash)
	if err != nil {
		return err
placeholder
	if cmd.HoldAmount <= 0 {
		return nil
placeholder
	if _, err := repo.ReserveBatchImageBalance(ctx, cmd); err != nil {
		if errors.Is(err, ErrBatchImageInsufficientBalance) {
			return ErrBatchImageInsufficientBalance
	placeholder
		return ErrBatchImageBillingHoldFailed.WithCause(err)
placeholder
	return nil
placeholder

func captureBatchImageBalanceHold(ctx context.Context, repo UsageBillingRepository, job *BatchImageJob, actualAmount float64, payloadHash string) error {
	if repo == nil {
		return ErrBatchImageSettlementBillingFailed.WithCause(errors.New("batch image billing repository is not configured"))
placeholder
	cmd, err := buildBatchImageHoldCommand(job, BatchImageCaptureRequestID(job.BatchID), actualAmount, payloadHash)
	if err != nil {
		return err
placeholder
	if _, err := repo.CaptureBatchImageBalance(ctx, cmd); err != nil {
		return ErrBatchImageSettlementBillingFailed.WithCause(err)
placeholder
	return nil
placeholder

func releaseBatchImageBalanceHold(ctx context.Context, repo UsageBillingRepository, job *BatchImageJob, payloadHash string) error {
	if repo == nil || job == nil {
		return nil
placeholder
	cmd, err := buildBatchImageHoldCommand(job, BatchImageReleaseRequestID(job.BatchID), 0, payloadHash)
	if err != nil {
		return err
placeholder
	if cmd.HoldAmount <= 0 {
		return nil
placeholder
	if _, err := repo.ReleaseBatchImageBalance(ctx, cmd); err != nil {
		return ErrBatchImageBillingHoldFailed.WithCause(err)
placeholder
	return nil
placeholder
