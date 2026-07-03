//go:build unit

package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBatchImageSettlementService_SettlesAndChargesSuccessfulImagesOnly(t *testing.T) {
	repo := newFakeBatchImageRepository()
	job := testSettlingBatchImageJob("imgbatch_settle")
	job.SuccessCount = 3
	job.FailCount = 2
	job.ItemCount = 5
	repo.jobs[job.BatchID] = job
	billing := &fakeBatchImageBillingRepo{placeholder
	svc := &BatchImageSettlementService{Repo: repo, BillingRepo: billing, Pricing: &fakeBatchImagePricingResolver{unitPrice: 0.25placeholderplaceholder

	result, err := svc.Settle(context.Background(), job.BatchID)
placeholder
	require.Equal(t, 0.75, result.ActualCost)
	require.Equal(t, "batch_image_settlement:"+job.BatchID, result.RequestID)
	require.False(t, result.AlreadySettled)
	require.Equal(t, BatchImageJobStatusCompleted, repo.jobs[job.BatchID].Status)
	require.NotNil(t, repo.jobs[job.BatchID].ActualCost)
	require.Equal(t, 0.75, *repo.jobs[job.BatchID].ActualCost)
	require.NotEmpty(t, batchImageDerefString(repo.jobs[job.BatchID].ManifestHash))
	require.NotNil(t, repo.jobs[job.BatchID].SettledAt)
	require.Len(t, billing.commands, 1)
	require.Equal(t, int64(321), billing.commands[0].APIKeyID)
	require.Equal(t, job.UserID, billing.commands[0].UserID)
	require.Equal(t, int64(654), billing.commands[0].AccountID)
	require.Equal(t, job.Model, billing.commands[0].Model)
	require.Equal(t, 3, billing.commands[0].ImageCount)
	require.Equal(t, 0.75, billing.commands[0].BalanceCost)
	require.Equal(t, "image", billing.commands[0].MediaType)
	require.NotContains(t, fmt.Sprintf("%+v", billing.commands[0]), batchImageTestData)
	require.NotContains(t, fmt.Sprintf("%+v", billing.commands[0]), "gs://")
	require.NotContains(t, fmt.Sprintf("%+v", billing.commands[0]), "prompt")
placeholder

func TestBatchImageSettlementService_ZeroSuccessCanComplete(t *testing.T) {
	repo := newFakeBatchImageRepository()
	job := testSettlingBatchImageJob("imgbatch_zero")
	job.SuccessCount = 0
	job.FailCount = 4
	job.ItemCount = 4
	repo.jobs[job.BatchID] = job
	billing := &fakeBatchImageBillingRepo{placeholder
	svc := &BatchImageSettlementService{Repo: repo, BillingRepo: billing, Pricing: &fakeBatchImagePricingResolver{unitPrice: 0.25placeholderplaceholder

	result, err := svc.Settle(context.Background(), job.BatchID)
placeholder
	require.Equal(t, 0.0, result.ActualCost)
	require.Equal(t, BatchImageJobStatusCompleted, repo.jobs[job.BatchID].Status)
	require.Len(t, billing.commands, 1)
	require.Equal(t, 0.0, billing.commands[0].BalanceCost)
placeholder

func TestBatchImageSettlementService_CompletedJobReturnsAlreadySettledWithoutBilling(t *testing.T) {
	repo := newFakeBatchImageRepository()
	job := testSettlingBatchImageJob("imgbatch_done")
	job.Status = BatchImageJobStatusCompleted
	cost := 0.5
	job.ActualCost = &cost
	repo.jobs[job.BatchID] = job
	billing := &fakeBatchImageBillingRepo{placeholder
	svc := &BatchImageSettlementService{Repo: repo, BillingRepo: billing, Pricing: &fakeBatchImagePricingResolver{unitPrice: 0.25placeholderplaceholder

	result, err := svc.Settle(context.Background(), job.BatchID)
placeholder
	require.True(t, result.AlreadySettled)
	require.Equal(t, 0.5, result.ActualCost)
	require.Empty(t, billing.commands)
placeholder

func TestBatchImageSettlementService_IdempotentAfterBillingCrash(t *testing.T) {
	repo := newFakeBatchImageRepository()
	job := testSettlingBatchImageJob("imgbatch_crash")
	repo.jobs[job.BatchID] = job
	billing := &fakeBatchImageBillingRepo{alreadyApplied: map[string]bool{BatchImageSettlementRequestID(job.BatchID): trueplaceholderplaceholder
	svc := &BatchImageSettlementService{Repo: repo, BillingRepo: billing, Pricing: &fakeBatchImagePricingResolver{unitPrice: 0.25placeholderplaceholder

	result, err := svc.Settle(context.Background(), job.BatchID)
placeholder
	require.Equal(t, 0.5, result.ActualCost)
	require.Equal(t, BatchImageJobStatusCompleted, repo.jobs[job.BatchID].Status)
	require.Len(t, billing.commands, 1)
placeholder

func TestBatchImageSettlementService_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*BatchImageJob)
		pricing BatchImagePricingResolver
		want    error
placeholder{
		{name: "invalid_status", mutate: func(j *BatchImageJob) { j.Status = BatchImageJobStatusRunning placeholder, want: ErrBatchImageSettlementInvalidStatusplaceholder,
		{name: "negative_success_count", mutate: func(j *BatchImageJob) { j.SuccessCount = -1 placeholder, want: ErrBatchImageSettlementInvalidCountsplaceholder,
		{name: "negative_fail_count", mutate: func(j *BatchImageJob) { j.FailCount = -1 placeholder, want: ErrBatchImageSettlementInvalidCountsplaceholder,
		{name: "missing_api_key", mutate: func(j *BatchImageJob) { j.APIKeyID = nil placeholder, want: ErrBatchImageSettlementMissingAPIKeyIDplaceholder,
		{name: "missing_account", mutate: func(j *BatchImageJob) { j.AccountID = nil placeholder, want: ErrBatchImageSettlementMissingAccountIDplaceholder,
		{name: "pricing_missing", pricing: &fakeBatchImagePricingResolver{err: ErrBatchImageSettlementPricingMissingplaceholder, want: ErrBatchImageSettlementPricingMissingplaceholder,
		{name: "manifest_conflict", mutate: func(j *BatchImageJob) { v := "different"; j.ManifestHash = &v placeholder, want: ErrBatchImageSettlementManifestConflictplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeBatchImageRepository()
			job := testSettlingBatchImageJob("imgbatch_" + tt.name)
			if tt.mutate != nil {
				tt.mutate(job)
		placeholder
			repo.jobs[job.BatchID] = job
			pricing := tt.pricing
			if pricing == nil {
				pricing = &fakeBatchImagePricingResolver{unitPrice: 0.25placeholder
		placeholder
			billing := &fakeBatchImageBillingRepo{placeholder
			svc := &BatchImageSettlementService{Repo: repo, BillingRepo: billing, Pricing: pricingplaceholder

			_, err := svc.Settle(context.Background(), job.BatchID)
			require.ErrorIs(t, err, tt.want)
			require.Empty(t, billing.commands)
			require.NotEqual(t, BatchImageJobStatusCompleted, repo.jobs[job.BatchID].Status)
	placeholder)
placeholder
placeholder

func TestBatchImageSettlementService_BillingFailureLeavesSettlingAndRecordsError(t *testing.T) {
	repo := newFakeBatchImageRepository()
	job := testSettlingBatchImageJob("imgbatch_billing_fail")
	repo.jobs[job.BatchID] = job
	billing := &fakeBatchImageBillingRepo{err: errors.New("temporary billing timeout with gs://hidden-output")placeholder
	svc := &BatchImageSettlementService{Repo: repo, BillingRepo: billing, Pricing: &fakeBatchImagePricingResolver{unitPrice: 0.25placeholderplaceholder

	_, err := svc.Settle(context.Background(), job.BatchID)
	require.ErrorIs(t, err, ErrBatchImageSettlementBillingFailed)
	require.Equal(t, BatchImageJobStatusSettling, repo.jobs[job.BatchID].Status)
	require.Equal(t, "SETTLEMENT_BILLING_FAILED", batchImageDerefString(repo.jobs[job.BatchID].LastErrorCode))
	require.Contains(t, batchImageDerefString(repo.jobs[job.BatchID].LastErrorMessage), "temporary billing timeout")
	require.NotNil(t, billing.commands[0])
placeholder

func TestBatchImagePipelineProcessor_SettlesQueuedSettlingJob(t *testing.T) {
	repo := newFakeBatchImageRepository()
	job := testSettlingBatchImageJob("imgbatch_pipeline")
	repo.jobs[job.BatchID] = job
	billing := &fakeBatchImageBillingRepo{placeholder
	settlement := &BatchImageSettlementService{Repo: repo, BillingRepo: billing, Pricing: &fakeBatchImagePricingResolver{unitPrice: 0.25placeholderplaceholder
	processor := &BatchImagePipelineProcessor{
		ProviderProcessor: &BatchImageProviderProcessor{Repo: repo, ProviderRegistry: NewBatchImageProviderRegistry(&fakeProcessorProvider{placeholder), AccountResolver: &fakeBatchImageAccountResolver{account: &Account{placeholderplaceholderplaceholder,
		SettlementService: settlement,
placeholder

	result, err := processor.Process(context.Background(), job.BatchID)
placeholder
	require.True(t, result.Terminal)
	require.Equal(t, BatchImageJobStatusCompleted, repo.jobs[job.BatchID].Status)
	require.Len(t, billing.commands, 1)
placeholder

func TestBatchImagePipelineProcessor_RequeuesTransientSettlementFailure(t *testing.T) {
	repo := newFakeBatchImageRepository()
	job := testSettlingBatchImageJob("imgbatch_pipeline_retry")
	repo.jobs[job.BatchID] = job
	settlement := &BatchImageSettlementService{Repo: repo, BillingRepo: &fakeBatchImageBillingRepo{err: errors.New("temporary")placeholder, Pricing: &fakeBatchImagePricingResolver{unitPrice: 0.25placeholderplaceholder
	processor := &BatchImagePipelineProcessor{
		ProviderProcessor: &BatchImageProviderProcessor{Repo: repo, ProviderRegistry: NewBatchImageProviderRegistry(&fakeProcessorProvider{placeholder), AccountResolver: &fakeBatchImageAccountResolver{account: &Account{placeholderplaceholderplaceholder,
		SettlementService: settlement,
placeholder

	result, err := processor.Process(context.Background(), job.BatchID)
placeholder
	require.False(t, result.Terminal)
	require.Equal(t, batchImageSettlementRetryDelay, result.RequeueAfter)
	require.Equal(t, BatchImageJobStatusSettling, repo.jobs[job.BatchID].Status)
placeholder

func TestBatchImageSettlementManifestHash(t *testing.T) {
	job := testSettlingBatchImageJob("imgbatch_hash")
	first := BuildBatchImageSettlementManifestHash(job)
	job.CreatedAt = job.CreatedAt.AddDate(0, 0, 1)
	job.UpdatedAt = job.UpdatedAt.AddDate(0, 0, 1)
	require.Equal(t, first, BuildBatchImageSettlementManifestHash(job))

	job.SuccessCount++
	require.NotEqual(t, first, BuildBatchImageSettlementManifestHash(job))

	job.SuccessCount--
	promptOrBase64 := first + " prompt " + batchImageTestData
	require.NotContains(t, BuildBatchImageSettlementManifestHash(job), promptOrBase64)
placeholder

func TestBatchImageSettlementBillingRequestIDs(t *testing.T) {
	repo := newFakeBatchImageRepository()
	first := testSettlingBatchImageJob("imgbatch_unique_1")
	second := testSettlingBatchImageJob("imgbatch_unique_2")
	repo.jobs[first.BatchID] = first
	repo.jobs[second.BatchID] = second
	billing := &fakeBatchImageBillingRepo{placeholder
	svc := &BatchImageSettlementService{Repo: repo, BillingRepo: billing, Pricing: &fakeBatchImagePricingResolver{unitPrice: 0.25placeholderplaceholder

	_, err := svc.Settle(context.Background(), first.BatchID)
placeholder
	_, err = svc.Settle(context.Background(), first.BatchID)
placeholder
	_, err = svc.Settle(context.Background(), second.BatchID)
placeholder

	require.Len(t, billing.commands, 2)
	require.Equal(t, "batch_image_settlement:"+first.BatchID, billing.commands[0].RequestID)
	require.Equal(t, "batch_image_settlement:"+second.BatchID, billing.commands[1].RequestID)
	require.NotEqual(t, billing.commands[0].RequestID, billing.commands[1].RequestID)
	require.Len(t, billing.seen, 2)
placeholder

func testSettlingBatchImageJob(batchID string) *BatchImageJob {
	apiKeyID := int64(321)
	accountID := int64(654)
	providerJobName := "providers/job"
	outputRef := "files/output"
	return &BatchImageJob{
		BatchID:           batchID,
		UserID:            123,
		APIKeyID:          &apiKeyID,
		AccountID:         &accountID,
		Provider:          BatchImageProviderGeminiAPI,
		Model:             "gemini-image",
		Status:            BatchImageJobStatusSettling,
		ProviderJobName:   &providerJobName,
		ProviderOutputRef: &outputRef,
		ItemCount:         3,
		SuccessCount:      2,
		FailCount:         1,
placeholder
placeholder

type fakeBatchImagePricingResolver struct {
	unitPrice float64
	err       error
placeholder

func (r *fakeBatchImagePricingResolver) BatchImageUnitPrice(context.Context, *BatchImageJob) (float64, error) {
	if r.err != nil {
		return 0, r.err
placeholder
	return r.unitPrice, nil
placeholder

type fakeBatchImageBillingRepo struct {
	commands       []*UsageBillingCommand
	seen           map[string]struct{placeholder
	alreadyApplied map[string]bool
	err            error
placeholder

func (r *fakeBatchImageBillingRepo) Apply(_ context.Context, cmd *UsageBillingCommand) (*UsageBillingApplyResult, error) {
	if r.seen == nil {
		r.seen = make(map[string]struct{placeholder)
placeholder
	if r.err != nil {
		r.commands = append(r.commands, cmd)
		return nil, r.err
placeholder
	if cmd != nil {
		cmd.Normalize()
		if _, ok := r.seen[cmd.RequestID]; ok || r.alreadyApplied[cmd.RequestID] {
			r.commands = append(r.commands, cmd)
			return &UsageBillingApplyResult{Applied: falseplaceholder, nil
	placeholder
		r.seen[cmd.RequestID] = struct{placeholder{placeholder
placeholder
	r.commands = append(r.commands, cmd)
	return &UsageBillingApplyResult{Applied: trueplaceholder, nil
placeholder

var _ UsageBillingRepository = (*fakeBatchImageBillingRepo)(nil)
var _ BatchImagePricingResolver = (*fakeBatchImagePricingResolver)(nil)
var _ = strings.TrimSpace
