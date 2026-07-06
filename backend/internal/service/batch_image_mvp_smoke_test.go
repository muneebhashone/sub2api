//go:build unit

package service

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func TestBatchImageMVPFlow(t *testing.T) {
	ctx := context.Background()
	repo := newFakeBatchImageRepository()
	queue := &publicBatchImageQueue{placeholder
	provider := &batchImageSmokeProvider{
		name: BatchImageProviderGeminiAPI,
		states: []BatchProviderInternalState{
			BatchProviderStateRunning,
			BatchProviderStateSucceeded,
	placeholder,
		result: batchImageSmokeResultJSONL(),
placeholder
	accountID := int64(101)
	accountRepo := &publicBatchImageAccountRepo{accounts: []Account{testBatchImageAccount(accountID, AccountTypeAPIKey)placeholderplaceholder
	cfg := &config.Config{BatchImage: config.BatchImageConfig{
		Enabled:                           true,
		MaxItemsPerJobDefault:             10,
		MaxPromptCharsPerItem:             8000,
		DefaultResponseMimeType:           "image/png",
		DefaultImageSize:                  "1K",
		MaxDownloadItemsZip:               10,
		MaxDownloadDurationSeconds:        60,
		OutputRetentionAfterTerminalHours: 72,
placeholderplaceholder
	registry := NewBatchImageProviderRegistry(provider)
	billing := &fakeBatchImageBillingRepo{placeholder
	pricing := &fakeBatchImagePricingResolver{unitPrice: 0.25placeholder
	owner := testBatchImageOwner()

	publicSvc := &BatchImagePublicService{
		Repo:             repo,
		AccountRepo:      accountRepo,
		Queue:            queue,
		ProviderRegistry: registry,
		Pricing:          pricing,
		BillingRepo:      billing,
		Config:           cfg,
placeholder
	processor := &BatchImagePipelineProcessor{
		ProviderProcessor: &BatchImageProviderProcessor{
			Repo:             repo,
			ProviderRegistry: registry,
			AccountResolver:  &fakeBatchImageAccountResolver{account: &accountRepo.accounts[0]placeholder,
			BillingRepo:      billing,
	placeholder,
		SettlementService: &BatchImageSettlementService{
			Repo:        repo,
			BillingRepo: billing,
			Pricing:     pricing,
			Config:      cfg,
	placeholder,
placeholder
	downloadSvc := &BatchImageDownloadService{
		Repo:             repo,
		ProviderRegistry: registry,
		AccountResolver:  &fakeBatchImageAccountResolver{account: &accountRepo.accounts[0]placeholder,
		Limiter:          &fakeBatchImageDownloadLimiter{placeholder,
		Config:           cfg,
placeholder
	cleanupSvc := &BatchImageCleanupService{
		Repo:             repo,
		ProviderRegistry: registry,
		AccountResolver:  &fakeBatchImageAccountResolver{account: &accountRepo.accounts[0]placeholder,
		Config:           cfg,
placeholder

	submitted, err := publicSvc.Submit(ctx, owner, validBatchImageSubmitRequest(), "")
placeholder
	require.Equal(t, "image.batch", submitted.Object)
	require.True(t, strings.HasPrefix(submitted.ID, "imgbatch_"))
	require.Equal(t, "queued", submitted.Status)
	require.Equal(t, 2, submitted.ItemCount)
	require.Equal(t, []string{submitted.IDplaceholder, queue.enqueued)
	require.Len(t, provider.submits, 1)
	require.Len(t, billing.reserves, 1)
	require.Equal(t, BatchImageHoldRequestID(submitted.ID), billing.reserves[0].RequestID)
	require.InDelta(t, 0.3, billing.reserves[0].HoldAmount, 1e-12)
	requireBatchImagePublicJSONHasNoInternals(t, mustMarshalBatchImageSmokeJSON(t, submitted))

	firstProcess, err := processor.Process(ctx, submitted.ID)
placeholder
	require.False(t, firstProcess.Terminal)
	require.Equal(t, BatchImageJobStatusRunning, repo.jobs[submitted.ID].Status)

	indexProcess, err := processor.Process(ctx, submitted.ID)
placeholder
	require.False(t, indexProcess.Terminal)
	require.Equal(t, time.Millisecond, indexProcess.RequeueAfter)
	require.Equal(t, BatchImageJobStatusSettling, repo.jobs[submitted.ID].Status)
	require.Equal(t, BatchImageCounts{SuccessCount: 1, FailCount: 1placeholder, repo.counts[submitted.ID])

	settleProcess, err := processor.Process(ctx, submitted.ID)
placeholder
	require.True(t, settleProcess.Terminal)
	job := repo.jobs[submitted.ID]
	require.Equal(t, BatchImageJobStatusCompleted, job.Status)
	require.NotNil(t, job.OutputExpiresAt)
	require.Equal(t, 1, job.SuccessCount)
	require.Equal(t, 1, job.FailCount)
	require.Len(t, billing.captures, 1)
	require.Equal(t, BatchImageCaptureRequestID(submitted.ID), billing.captures[0].RequestID)
	require.InDelta(t, 0.3, billing.captures[0].HoldAmount, 1e-12)
	require.InDelta(t, 0.125, billing.captures[0].ActualAmount, 1e-12)

	secondSettlement, err := processor.SettlementService.Settle(ctx, submitted.ID)
placeholder
	require.True(t, secondSettlement.AlreadySettled)
	require.Len(t, billing.captures, 1)

	status, err := publicSvc.Get(ctx, owner, submitted.ID)
placeholder
	require.Equal(t, "completed", status.Status)
	require.Equal(t, 1, status.SuccessCount)
	require.Equal(t, 1, status.FailCount)
	require.NotNil(t, status.ActualCost)
	requireBatchImagePublicJSONHasNoInternals(t, mustMarshalBatchImageSmokeJSON(t, status))

	items, err := publicSvc.ListItems(ctx, owner, submitted.ID, BatchImageItemsQuery{Limit: 100placeholder)
placeholder
	require.False(t, items.HasMore)
	require.Len(t, items.Data, 2)
	require.Equal(t, "cover_001", items.Data[0].CustomID)
	require.Equal(t, "succeeded", items.Data[0].Status)
	require.Equal(t, "cover_002", items.Data[1].CustomID)
	require.Equal(t, "failed", items.Data[1].Status)
	require.NotNil(t, items.Data[1].Error)
	require.Nil(t, repo.items[submitted.ID][1].BilledAmount)
	requireBatchImagePublicJSONHasNoInternals(t, mustMarshalBatchImageSmokeJSON(t, items))

	stream, err := downloadSvc.OpenItemContent(ctx, owner, submitted.ID, "cover_001", 0)
placeholder
	body, err := io.ReadAll(stream.Reader)
placeholder
	require.NoError(t, stream.Reader.Close())
	require.Equal(t, []byte("smoke-png"), body)
	require.Equal(t, "image/png", stream.ContentType)
	require.Equal(t, "cover_001.png", stream.Filename)

	var zipBuf bytes.Buffer
	zipResult, err := downloadSvc.StreamZip(ctx, owner, submitted.ID, BatchImageZipOptions{placeholder, &zipBuf)
placeholder
	require.Equal(t, 1, zipResult.FileCount)
	require.Equal(t, 1, zipResult.ErrorCount)
	zipFiles := readZipFiles(t, zipBuf.Bytes())
	require.Equal(t, []byte("smoke-png"), zipFiles["images/cover_001.png"])
	require.Contains(t, zipFiles, "manifest.json")
	require.Contains(t, zipFiles, "errors.json")
	requireBatchImagePublicJSONHasNoInternals(t, string(bytes.Join(mapValues(zipFiles), []byte("\n"))))

	zipReader, err := zip.NewReader(bytes.NewReader(zipBuf.Bytes()), int64(zipBuf.Len()))
placeholder
	require.ElementsMatch(t, []string{"images/cover_001.png", "manifest.json", "errors.json"placeholder, batchImageSmokeZipNames(zipReader))

	deleted, err := cleanupSvc.DeleteOutputsForOwner(ctx, owner, submitted.ID)
placeholder
	require.Equal(t, "output_deleted", deleted.Status)
	require.Equal(t, []CleanupTarget{CleanupTargetOutputplaceholder, provider.cleanupTargets)
	requireBatchImagePublicJSONHasNoInternals(t, mustMarshalBatchImageSmokeJSON(t, deleted))

	deletedAgain, err := cleanupSvc.DeleteOutputsForOwner(ctx, owner, submitted.ID)
placeholder
	require.Equal(t, "output_deleted", deletedAgain.Status)
	require.Equal(t, []CleanupTarget{CleanupTargetOutputplaceholder, provider.cleanupTargets)

	stream, err = downloadSvc.OpenItemContent(ctx, owner, submitted.ID, "cover_001", 0)
	require.Nil(t, stream)
	require.ErrorIs(t, err, ErrBatchImageOutputDeleted)
	var afterDelete bytes.Buffer
	zipResult, err = downloadSvc.StreamZip(ctx, owner, submitted.ID, BatchImageZipOptions{placeholder, &afterDelete)
	require.Nil(t, zipResult)
	require.ErrorIs(t, err, ErrBatchImageOutputDeleted)
	require.Empty(t, afterDelete.Bytes())
placeholder

func batchImageSmokeResultJSONL() string {
	return strings.Join([]string{
		`{"key":"cover_001","response":{"candidates":[{"content":{"parts":[{"inlineData":{"mimeType":"image/png","data":"c21va2UtcG5n"placeholderplaceholder]placeholderplaceholder]placeholderplaceholder`,
		`{"key":"cover_002","status":{"code":3,"message":"blocked by safety policy"placeholderplaceholder`,
placeholder, "\n") + "\n"
placeholder

func mustMarshalBatchImageSmokeJSON(t *testing.T, value any) string {
placeholder
	body, err := json.Marshal(value)
placeholder
	return string(body)
placeholder

func batchImageSmokeZipNames(reader *zip.Reader) []string {
	names := make([]string, 0, len(reader.File))
	for _, file := range reader.File {
		names = append(names, file.Name)
placeholder
	return names
placeholder

type batchImageSmokeProvider struct {
	name           string
	states         []BatchProviderInternalState
	submits        []BatchImageInput
	result         string
	cleanupTargets []CleanupTarget
placeholder

func (p *batchImageSmokeProvider) Name() string { return p.name placeholder

func (p *batchImageSmokeProvider) SupportsAccount(account *Account) bool {
	return account != nil && account.IsSchedulable()
placeholder

func (p *batchImageSmokeProvider) Submit(_ context.Context, _ *BatchImageJob, _ *Account, input BatchImageInput) (*BatchProviderJob, error) {
	p.submits = append(p.submits, input)
	return &BatchProviderJob{
		ProviderJobName:   "providers/fake-provider-job/raw-id",
		ProviderInputRef:  "files/fake-provider-job/input.jsonl",
		ProviderOutputRef: "files/fake-provider-job/output.jsonl",
placeholder, nil
placeholder

func (p *batchImageSmokeProvider) Get(context.Context, *BatchImageJob, *Account) (*BatchProviderStatus, error) {
	state := BatchProviderStateSucceeded
	if len(p.states) > 0 {
		state = p.states[0]
		p.states = p.states[1:]
placeholder
	return &BatchProviderStatus{
		RawState:          strings.ToUpper(string(state)),
		InternalState:     state,
		Done:              state == BatchProviderStateSucceeded,
		ProviderOutputRef: "files/fake-provider-job/output.jsonl",
placeholder, nil
placeholder

func (p *batchImageSmokeProvider) Cancel(context.Context, *BatchImageJob, *Account) error {
	return nil
placeholder

func (p *batchImageSmokeProvider) OpenResult(context.Context, *BatchImageJob, *Account) (io.ReadCloser, string, error) {
	return io.NopCloser(strings.NewReader(p.result)), "application/jsonl", nil
placeholder

func (p *batchImageSmokeProvider) Cleanup(_ context.Context, _ *BatchImageJob, _ *Account, target CleanupTarget) error {
	p.cleanupTargets = append(p.cleanupTargets, target)
	return nil
placeholder

var _ BatchImageProvider = (*batchImageSmokeProvider)(nil)
