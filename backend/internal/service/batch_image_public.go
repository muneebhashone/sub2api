package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"go.uber.org/zap"
)

const (
	defaultBatchImageMaxItems           = 200
	defaultBatchImageMaxOutputImages    = 200
	defaultBatchImageMaxOutputCount     = 4
	defaultBatchImageMaxPromptChars     = 8000
	defaultBatchImageResponseMime       = "image/png"
	defaultBatchImageImageSize          = "1K"
	defaultBatchImageDiscountMultiplier = 0.5
	defaultBatchImageHoldMultiplier     = 0.6
	maxBatchImagePublicErrorChars       = 500
	maxBatchImageReferenceImageBytes    = 10 * 1024 * 1024
	defaultBatchImageMaxReferenceImages = 1000
	defaultBatchImageMaxReferenceBytes  = 128 * 1024 * 1024
)

type BatchImageAccountSelectionRepository interface {
	GetByID(ctx context.Context, id int64) (*Account, error)
	ListSchedulableByPlatform(ctx context.Context, platform string) ([]Account, error)
	ListSchedulableByGroupIDAndPlatform(ctx context.Context, groupID int64, platform string) ([]Account, error)
placeholder

type BatchImageGroupPricingRepository interface {
	GetByIDLite(ctx context.Context, id int64) (*Group, error)
placeholder

type BatchImageUserGroupRateRepository interface {
	GetByUserAndGroup(ctx context.Context, userID, groupID int64) (*float64, error)
placeholder

type BatchImageSubmitRequest struct {
	Model            string                 `json:"model"`
	TaskName         string                 `json:"task_name"`
	ParentBatchID    string                 `json:"parent_batch_id"`
	Provider         string                 `json:"provider"`
	Items            []BatchImageSubmitItem `json:"items"`
	ResponseMimeType string                 `json:"response_mime_type"`
	AspectRatio      string                 `json:"aspect_ratio"`
	ImageSize        string                 `json:"image_size"`
	Metadata         map[string]string      `json:"metadata"`
	SessionID        *string                `json:"-"`
placeholder

type BatchImageSubmitItem struct {
	CustomID        string                     `json:"custom_id"`
	Prompt          string                     `json:"prompt"`
	OutputCount     int                        `json:"output_count,omitempty"`
	ReferenceImages []BatchImageReferenceInput `json:"reference_images,omitempty"`
placeholder

type BatchImageReferenceInput struct {
	ID       string `json:"id,omitempty"`
	Type     string `json:"type,omitempty"`
	MimeType string `json:"mime_type"`
	Data     []byte `json:"data,omitempty"`
	FileURI  string `json:"file_uri,omitempty"`
placeholder

type BatchImageOwner struct {
	UserID   int64
	APIKeyID int64
	GroupID  *int64
placeholder

type BatchImagePublicService struct {
	Repo              BatchImageRepository
	AccountRepo       BatchImageAccountSelectionRepository
	GroupRepo         BatchImageGroupPricingRepository
	UserGroupRateRepo BatchImageUserGroupRateRepository
	Queue             BatchImageQueue
	ProviderRegistry  *BatchImageProviderRegistry
	Pricing           BatchImagePricingResolver
	BillingRepo       UsageBillingRepository
	AuthCache         APIKeyAuthCacheInvalidator
	Config            *config.Config
placeholder

type BatchImagePricingSnapshot struct {
	BaseUnitPrice           float64
	GroupRateMultiplier     float64
	AccountRateMultiplier   float64
	BatchDiscountMultiplier float64
	HoldMultiplier          float64
	BillableUnitPrice       float64
	HoldUnitPrice           float64
	EstimatedCost           float64
	HoldAmount              float64
placeholder

type BatchImagePublicBatch struct {
	ID              string   `json:"id"`
	Object          string   `json:"object"`
	TaskName        string   `json:"task_name"`
	ParentBatchID   *string  `json:"parent_batch_id,omitempty"`
	Status          string   `json:"status"`
	Model           string   `json:"model"`
	Provider        string   `json:"provider"`
	ItemCount       int      `json:"item_count"`
	SuccessCount    int      `json:"success_count"`
	FailCount       int      `json:"fail_count"`
	EstimatedCost   float64  `json:"estimated_cost"`
	HoldAmount      float64  `json:"hold_amount"`
	ActualCost      *float64 `json:"actual_cost"`
	CreatedAt       int64    `json:"created_at"`
	SubmittedAt     *int64   `json:"submitted_at"`
	SettledAt       *int64   `json:"settled_at"`
	DownloadedAt    *int64   `json:"downloaded_at,omitempty"`
	OutputDeletedAt *int64   `json:"output_deleted_at,omitempty"`
placeholder

type BatchImagePublicItem struct {
	CustomID      string                 `json:"custom_id"`
	Status        string                 `json:"status"`
	PromptPreview *string                `json:"prompt_preview,omitempty"`
	MimeType      *string                `json:"mime_type"`
	FileExtension *string                `json:"file_extension"`
	ImageCount    int                    `json:"image_count"`
	Error         *BatchImagePublicError `json:"error"`
placeholder

type BatchImagePublicError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Source  string `json:"source,omitempty"`
placeholder

type BatchImagePublicItemsResponse struct {
	Object  string                 `json:"object"`
	Data    []BatchImagePublicItem `json:"data"`
	HasMore bool                   `json:"has_more"`
placeholder

type BatchImagePublicListResponse struct {
	Object  string                   `json:"object"`
	Data    []*BatchImagePublicBatch `json:"data"`
	HasMore bool                     `json:"has_more"`
placeholder

type BatchImagePublicModel struct {
	ID       string `json:"id"`
	Object   string `json:"object"`
	Provider string `json:"provider"`
placeholder

type BatchImagePublicModelsResponse struct {
	Object string                  `json:"object"`
	Data   []BatchImagePublicModel `json:"data"`
placeholder

type BatchImageJobsQuery struct {
	Status     string
	TaskName   string
	Downloaded string
	From       string
	To         string
	Limit      int
	Cursor     string
placeholder

type BatchImageItemsQuery struct {
	Status string
	Limit  int
	Cursor string
placeholder

func NewBatchImagePublicService(repo BatchImageRepository, accountRepo AccountRepository, groupRepo GroupRepository, userGroupRateRepo UserGroupRateRepository, queue BatchImageQueue, pricing *BatchImageModelPricingResolver, billingRepo UsageBillingRepository, authCache APIKeyAuthCacheInvalidator, cfg *config.Config) *BatchImagePublicService {
	return &BatchImagePublicService{
		Repo:              repo,
		AccountRepo:       accountRepo,
		GroupRepo:         groupRepo,
		UserGroupRateRepo: userGroupRateRepo,
		Queue:             queue,
		ProviderRegistry:  NewBatchImageProviderRegistryFromConfig(cfg),
		Pricing:           pricing,
		BillingRepo:       billingRepo,
		AuthCache:         authCache,
		Config:            cfg,
placeholder
placeholder

func (s *BatchImagePublicService) Submit(ctx context.Context, owner BatchImageOwner, req BatchImageSubmitRequest, idempotencyKey string) (*BatchImagePublicBatch, error) {
	if !s.enabled() {
		return nil, ErrBatchImageDisabled
placeholder
	normalized, err := s.validateSubmitRequest(req)
	if err != nil {
		return nil, err
placeholder
	// 与 ListModels 使用同一鉴权谓词（AllowBatchImageGeneration + Platform==Gemini），
	// 避免两个入口校验口径不一致留下防御纵深缺口。
	if err := s.ensureGroupAllowsBatchImage(ctx, owner.GroupID); err != nil {
		return nil, err
placeholder
	requestHash := HashBatchImageSubmitRequest(normalized)
	idempotencyKey = strings.TrimSpace(idempotencyKey)
	if idempotencyKey != "" {
		existing, err := s.Repo.GetBatchImageJobByIdempotencyKey(ctx, owner.UserID, owner.APIKeyID, idempotencyKey)
		if err == nil {
			if batchImageDerefString(existing.RequestHash) != requestHash {
				return nil, ErrBatchImageIdempotencyConflict
		placeholder
			if existing.Status == BatchImageJobStatusSubmitted && s.Queue != nil {
				if enqueueErr := s.Queue.Enqueue(ctx, existing.BatchID); enqueueErr != nil && !errors.Is(enqueueErr, ErrBatchImageAlreadyQueued) {
					_ = s.Repo.RecordBatchImageJobSubmitFailure(ctx, existing.BatchID, "QUEUE_FAILED", sanitizeBatchImagePublicMessage(enqueueErr.Error()), false)
					return nil, ErrBatchImageQueueFailed
			placeholder
		placeholder
			return BatchImageJobToPublic(existing), nil
	placeholder
		if !errors.Is(err, ErrBatchImageJobNotFound) {
			return nil, err
	placeholder
placeholder

	provider, account, err := s.selectProviderAndAccount(ctx, owner, normalized.Provider, normalized.Model)
	if err != nil {
		return nil, err
placeholder
	pricingSnapshot, err := s.resolvePricingSnapshot(ctx, owner, normalized, provider.Name(), account)
	if err != nil {
		return nil, err
placeholder
	parentBatchID := batchImageOptionalStringPtr(normalized.ParentBatchID)
	if parentBatchID != nil {
		parent, parentErr := s.Repo.GetBatchImageJobByBatchIDForOwner(ctx, owner.UserID, owner.APIKeyID, *parentBatchID)
		if parentErr != nil {
			return nil, parentErr
	placeholder
		if parent.ParentBatchID != nil && strings.TrimSpace(*parent.ParentBatchID) != "" {
			parentBatchID = batchImageOptionalStringPtr(*parent.ParentBatchID)
	placeholder
placeholder
	batchID, err := NewBatchImageID()
	if err != nil {
		return nil, err
placeholder
	apiKeyID := owner.APIKeyID
	accountID := account.ID
	holdID := BatchImageHoldRequestID(batchID)
	holdAmount := pricingSnapshot.HoldAmount
	job, err := s.Repo.CreateBatchImageJob(ctx, CreateBatchImageJobParams{
		BatchID:                 batchID,
		UserID:                  owner.UserID,
		APIKeyID:                &apiKeyID,
		AccountID:               &accountID,
		Provider:                provider.Name(),
		Model:                   normalized.Model,
		TaskName:                normalized.TaskName,
		ParentBatchID:           parentBatchID,
		Status:                  BatchImageJobStatusCreated,
		ItemCount:               len(normalized.Items),
		EstimatedCost:           pricingSnapshot.EstimatedCost,
		HoldAmount:              &holdAmount,
		BaseUnitPrice:           pricingSnapshot.BaseUnitPrice,
		GroupRateMultiplier:     pricingSnapshot.GroupRateMultiplier,
		AccountRateMultiplier:   pricingSnapshot.AccountRateMultiplier,
		BatchDiscountMultiplier: pricingSnapshot.BatchDiscountMultiplier,
		HoldMultiplier:          pricingSnapshot.HoldMultiplier,
		BillableUnitPrice:       pricingSnapshot.BillableUnitPrice,
		HoldUnitPrice:           pricingSnapshot.HoldUnitPrice,
		PricingSnapshotVersion:  1,
		Currency:                "USD",
		HoldID:                  &holdID,
		IdempotencyKey:          batchImageOptionalStringPtr(idempotencyKey),
		RequestHash:             batchImageStringPtr(requestHash),
		SessionID:               normalized.SessionID,
placeholder)
	if err != nil {
		return nil, err
placeholder
	if err := reserveBatchImageBalanceHold(ctx, s.BillingRepo, job, requestHash); err != nil {
		code := "BILLING_HOLD_FAILED"
		if errors.Is(err, ErrBatchImageInsufficientBalance) {
			code = "INSUFFICIENT_BALANCE"
	placeholder
		_ = s.Repo.RecordBatchImageJobSubmitFailure(ctx, job.BatchID, code, sanitizeBatchImagePublicMessage(err.Error()), true)
		s.hidePreUpstreamSubmitFailure(ctx, owner, job)
		return nil, err
placeholder
	s.invalidateAuthCache(ctx, owner.UserID)
	if err := s.createPendingItems(ctx, job.BatchID, requestHash, normalized.Items); err != nil {
		if releaseErr := s.releaseFailedSubmitHold(ctx, job, requestHash); releaseErr != nil {
			return nil, releaseErr
	placeholder
		_ = s.Repo.RecordBatchImageJobSubmitFailure(ctx, job.BatchID, "ITEM_CREATE_FAILED", sanitizeBatchImagePublicMessage(err.Error()), true)
		s.hidePreUpstreamSubmitFailure(ctx, owner, job)
		return nil, ErrBatchImageQueueFailed
placeholder

	input := BatchImageInput{
		BatchID:          job.BatchID,
		Model:            normalized.Model,
		DisplayName:      job.BatchID,
		ResponseMimeType: normalized.ResponseMimeType,
		AspectRatio:      normalized.AspectRatio,
		ImageSize:        normalized.ImageSize,
		Metadata:         normalized.Metadata,
		Items:            make([]BatchImageInputItem, 0, len(normalized.Items)),
placeholder
	for _, item := range normalized.Items {
		refs := make([]BatchImageReference, 0, len(item.ReferenceImages))
		for _, ref := range item.ReferenceImages {
			refs = append(refs, BatchImageReference(ref))
	placeholder
		input.Items = append(input.Items, BatchImageInputItem{
			CustomID:        item.CustomID,
			Prompt:          item.Prompt,
			ReferenceImages: refs,
	placeholder)
placeholder

	// 上游提交（上传参考图 + 创建批任务）可能长达数分钟且不刷新 updated_at，
	// 会被 stale 恢复扫描误判为滞留并退款。提交前转入 uploading 刷新时间戳，
	// 提交期间用心跳持续续期。
	if err := s.Repo.TransitionBatchImageJobStatus(ctx, job.BatchID, BatchImageJobStatusUploading, BatchImageTransitionOptions{
		EventType:    "upload_started",
		EventPayload: map[string]any{"batch_id": job.BatchIDplaceholder,
placeholder); err != nil {
		if releaseErr := s.releaseFailedSubmitHold(ctx, job, requestHash); releaseErr != nil {
			return nil, releaseErr
	placeholder
		// 并发 Cancel 等导致的非法转换：job 已处于终态，不再覆盖其状态。
		if !errors.Is(err, ErrBatchImageInvalidTransition) {
			_ = s.Repo.RecordBatchImageJobSubmitFailure(ctx, job.BatchID, "UPLOAD_TRANSITION_FAILED", sanitizeBatchImagePublicMessage(err.Error()), true)
			s.hidePreUpstreamSubmitFailure(ctx, owner, job)
	placeholder
		return nil, err
placeholder
	job.Status = BatchImageJobStatusUploading

	hbCtx, hbCancel := context.WithCancel(ctx)
	hbDone := make(chan struct{placeholder)
	go s.runSubmitHeartbeat(hbCtx, job.BatchID, hbDone)
	providerJob, err := provider.Submit(ctx, job, account, input)
	hbCancel()
	<-hbDone
	if err != nil {
		if releaseErr := s.releaseFailedSubmitHold(ctx, job, requestHash); releaseErr != nil {
			return nil, releaseErr
	placeholder
		publicErr := batchImageProviderSubmitPublicError(err)
		reason := batchImageProviderSubmitRecordCode(publicErr)
		_ = s.Repo.RecordBatchImageJobSubmitFailure(ctx, job.BatchID, reason, sanitizeBatchImagePublicMessage(err.Error()), true)
		s.hidePreUpstreamSubmitFailure(ctx, owner, job)
		return nil, publicErr
placeholder
	if providerJob == nil || strings.TrimSpace(providerJob.ProviderJobName) == "" {
		if releaseErr := s.releaseFailedSubmitHold(ctx, job, requestHash); releaseErr != nil {
			return nil, releaseErr
	placeholder
		_ = s.Repo.RecordBatchImageJobSubmitFailure(ctx, job.BatchID, "PROVIDER_SUBMIT_FAILED", "provider job name missing", true)
		s.hidePreUpstreamSubmitFailure(ctx, owner, job)
		return nil, ErrBatchImageProviderSubmitFailed
placeholder

	if err := s.Repo.UpdateBatchImageJobProviderSubmit(ctx, UpdateBatchImageJobProviderSubmitParams{
		BatchID:           job.BatchID,
		ProviderJobName:   providerJob.ProviderJobName,
		ProviderInputRef:  providerJob.ProviderInputRef,
		ProviderOutputRef: providerJob.ProviderOutputRef,
		GCSInputURI:       batchImageGCSRef(provider.Name(), providerJob.ProviderInputRef),
		GCSOutputURI:      batchImageGCSRef(provider.Name(), providerJob.ProviderOutputRef),
		EventPayload:      map[string]any{"provider": provider.Name()placeholder,
placeholder); err != nil {
		// job 可能已被恢复扫描转 failed 并退款：上游批任务已创建成功，
		// 必须尽力取消并清理输入，否则上游照常产生成本（孤儿任务）。
		s.abortOrphanProviderJob(ctx, provider, job, account, providerJob)
		return nil, err
placeholder

	if s.Queue != nil {
		if err := s.Queue.Enqueue(ctx, job.BatchID); err != nil && !errors.Is(err, ErrBatchImageAlreadyQueued) {
			_ = s.Repo.RecordBatchImageJobSubmitFailure(ctx, job.BatchID, "QUEUE_FAILED", sanitizeBatchImagePublicMessage(err.Error()), false)
			return nil, ErrBatchImageQueueFailed
	placeholder
placeholder

	created, err := s.Repo.GetBatchImageJobByBatchID(ctx, job.BatchID)
	if err != nil {
		return nil, err
placeholder
	return BatchImageJobToPublic(created), nil
placeholder

func (s *BatchImagePublicService) releaseFailedSubmitHold(ctx context.Context, job *BatchImageJob, requestHash string) error {
	if err := releaseBatchImageBalanceHold(ctx, s.BillingRepo, job, requestHash); err != nil {
		_ = s.Repo.RecordBatchImageJobSubmitFailure(ctx, job.BatchID, "BILLING_RELEASE_FAILED", sanitizeBatchImagePublicMessage(err.Error()), true)
		s.enqueueBillingRetry(ctx, job.BatchID)
		return ErrBatchImageBillingHoldFailed
placeholder
	s.invalidateAuthCache(ctx, job.UserID)
	return nil
placeholder

// runSubmitHeartbeat 在 provider.Submit 期间周期性刷新 job 的 updated_at，
// 使 stale 恢复扫描能区分"仍在慢提交"与"进程死亡后的滞留"。
func (s *BatchImagePublicService) runSubmitHeartbeat(ctx context.Context, batchID string, done chan<- struct{placeholder) {
	defer close(done)
	interval := s.submitHeartbeatInterval()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.Repo.TouchBatchImageJobSubmitting(ctx, batchID); err != nil && ctx.Err() == nil {
				logger.L().Warn("batch_image.submit_heartbeat_failed",
					zap.String("batch_id", batchID),
					zap.Error(err),
				)
		placeholder
	placeholder
placeholder
placeholder

func (s *BatchImagePublicService) submitHeartbeatInterval() time.Duration {
	staleAfter := 10 * time.Minute
	if s != nil && s.Config != nil && s.Config.BatchImage.StaleActiveAfterSeconds > 0 {
		staleAfter = time.Duration(s.Config.BatchImage.StaleActiveAfterSeconds) * time.Second
placeholder
	interval := staleAfter / 3
	if interval < 15*time.Second {
		interval = 15 * time.Second
placeholder
	return interval
placeholder

// abortOrphanProviderJob 在上游任务创建成功但本地状态推进失败时，
// 尽力取消上游批任务并清理已上传的输入文件，避免孤儿任务持续产生成本。
func (s *BatchImagePublicService) abortOrphanProviderJob(ctx context.Context, provider BatchImageProvider, job *BatchImageJob, account *Account, providerJob *BatchProviderJob) {
	if s == nil || provider == nil || job == nil || providerJob == nil {
		return
placeholder
	orphan := *job
	orphan.ProviderJobName = batchImageOptionalStringPtr(providerJob.ProviderJobName)
	orphan.ProviderInputRef = batchImageOptionalStringPtr(providerJob.ProviderInputRef)
	orphan.GCSInputURI = batchImageOptionalStringPtr(batchImageGCSRef(provider.Name(), providerJob.ProviderInputRef))
	if err := provider.Cancel(ctx, &orphan, account); err != nil {
		logger.L().Warn("batch_image.orphan_provider_job_cancel_failed",
			zap.String("batch_id", job.BatchID),
			zap.String("provider", provider.Name()),
			zap.Error(err),
		)
placeholder
	if err := provider.Cleanup(ctx, &orphan, account, CleanupTargetInput); err != nil {
		logger.L().Warn("batch_image.orphan_provider_job_cleanup_failed",
			zap.String("batch_id", job.BatchID),
			zap.String("provider", provider.Name()),
			zap.Error(err),
		)
placeholder
	if err := s.Repo.AppendBatchImageEvent(ctx, job.BatchID, "provider_job_aborted_after_submit", map[string]any{
		"batch_id": job.BatchID,
		"provider": provider.Name(),
placeholder); err != nil {
		logger.L().Warn("batch_image.orphan_provider_job_event_failed",
			zap.String("batch_id", job.BatchID),
			zap.Error(err),
		)
placeholder
placeholder

func (s *BatchImagePublicService) createPendingItems(ctx context.Context, batchID, requestHash string, items []BatchImageSubmitItem) error {
	if s == nil || s.Repo == nil || len(items) == 0 {
		return nil
placeholder
	params := make([]CreateBatchImageItemParams, 0, len(items))
	for _, item := range items {
		preview := truncateBatchImageMessage(item.Prompt, s.maxPromptChars())
		params = append(params, CreateBatchImageItemParams{
			JobID:         batchID,
			CustomID:      item.CustomID,
			Status:        BatchImageItemStatusPending,
			RequestHash:   batchImageStringPtr(requestHash),
			PromptPreview: batchImageStringPtr(preview),
			ImageCount:    0,
	placeholder)
placeholder
	return s.Repo.BulkCreateBatchImageItems(ctx, params)
placeholder

func (s *BatchImagePublicService) enqueueBillingRetry(ctx context.Context, batchID string) {
	if s == nil || s.Queue == nil {
		return
placeholder
	if err := s.Queue.Enqueue(ctx, batchID); err != nil && !errors.Is(err, ErrBatchImageAlreadyQueued) {
		logger.L().Warn("batch_image.billing_retry_enqueue_failed",
			zap.String("batch_id", batchID),
			zap.Error(err),
		)
		if eventErr := s.Repo.AppendBatchImageEvent(ctx, batchID, "billing_retry_enqueue_failed", map[string]any{
			"batch_id": batchID,
			"error":    sanitizeBatchImagePublicMessage(err.Error()),
	placeholder); eventErr != nil {
			logger.L().Warn("batch_image.billing_retry_event_failed",
				zap.String("batch_id", batchID),
				zap.Error(eventErr),
			)
	placeholder
placeholder
placeholder

func (s *BatchImagePublicService) hidePreUpstreamSubmitFailure(ctx context.Context, owner BatchImageOwner, job *BatchImageJob) {
	if s == nil || s.Repo == nil || job == nil || job.ProviderJobName != nil {
		return
placeholder
	if err := s.Repo.MarkBatchImageJobUserDeleted(ctx, owner.UserID, owner.APIKeyID, job.BatchID, time.Now()); err != nil {
		logger.L().Warn("batch_image.hide_pre_upstream_failure_failed",
			zap.String("batch_id", job.BatchID),
			zap.Error(err),
		)
placeholder
placeholder

func (s *BatchImagePublicService) Get(ctx context.Context, owner BatchImageOwner, batchID string) (*BatchImagePublicBatch, error) {
	job, err := s.Repo.GetBatchImageJobByBatchIDForOwner(ctx, owner.UserID, owner.APIKeyID, batchID)
	if err != nil {
		return nil, err
placeholder
	return BatchImageJobToPublic(job), nil
placeholder

func (s *BatchImagePublicService) List(ctx context.Context, owner BatchImageOwner, query BatchImageJobsQuery) (*BatchImagePublicListResponse, error) {
	filter := BatchImageJobFilter{Limit: query.Limit, Offset: parseBatchImageCursor(query.Cursor), ExcludeDeleted: trueplaceholder
	filter.TaskNameLike = strings.TrimSpace(query.TaskName)
	switch strings.TrimSpace(query.Status) {
	case "", "all":
	case "queued":
		filter.Status = BatchImageJobStatusSubmitted
	case "processing_results":
		filter.Status = BatchImageJobStatusIndexing
	case "completed":
		filter.Status = BatchImageJobStatusCompleted
	case "failed":
		filter.Status = BatchImageJobStatusFailed
	case "cancelled":
		filter.Status = BatchImageJobStatusCancelled
	case "output_deleted":
		filter.Status = BatchImageJobStatusOutputDeleted
	default:
		filter.Status = strings.TrimSpace(query.Status)
placeholder
	switch strings.TrimSpace(strings.ToLower(query.Downloaded)) {
	case "", "all":
	case "true", "1", "yes", "downloaded":
		downloaded := true
		filter.Downloaded = &downloaded
	case "false", "0", "no", "not_downloaded":
		downloaded := false
		filter.Downloaded = &downloaded
	default:
		return nil, ErrBatchImageInvalidItems
placeholder
	if from := parseBatchImageListTime(query.From); from != nil {
		filter.CreatedAfter = from
placeholder
	if to := parseBatchImageListTime(query.To); to != nil {
		filter.CreatedBefore = to
placeholder
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 20
placeholder
	jobs, err := s.Repo.ListBatchImageJobsForOwner(ctx, owner.UserID, owner.APIKeyID, filter)
	if err != nil {
		return nil, err
placeholder
	data := make([]*BatchImagePublicBatch, 0, len(jobs))
	for _, job := range jobs {
		data = append(data, BatchImageJobToPublic(job))
placeholder
	return &BatchImagePublicListResponse{
		Object:  "list",
		Data:    data,
		HasMore: len(data) == filter.Limit,
placeholder, nil
placeholder

func (s *BatchImagePublicService) MarkDownloaded(ctx context.Context, owner BatchImageOwner, batchID string) error {
	job, err := s.Repo.GetBatchImageJobByBatchIDForOwner(ctx, owner.UserID, owner.APIKeyID, batchID)
	if err != nil {
		return err
placeholder
	return s.Repo.MarkBatchImageDownloaded(ctx, job.BatchID, time.Now())
placeholder

func (s *BatchImagePublicService) DeleteRecord(ctx context.Context, owner BatchImageOwner, batchID string) error {
	job, err := s.Repo.GetBatchImageJobByBatchIDForOwner(ctx, owner.UserID, owner.APIKeyID, batchID)
	if err != nil {
		return err
placeholder
	if !isBatchImageProcessorDoneStatus(job.Status) {
		return ErrBatchImageRecordDeleteNotReady
placeholder
	return s.Repo.MarkBatchImageJobUserDeleted(ctx, owner.UserID, owner.APIKeyID, job.BatchID, time.Now())
placeholder

func (s *BatchImagePublicService) ListModels(ctx context.Context, owner BatchImageOwner) (*BatchImagePublicModelsResponse, error) {
	if !s.enabled() {
		return nil, ErrBatchImageDisabled
placeholder
	if s.Pricing == nil {
		return nil, ErrBatchImageSettlementPricingMissing
placeholder
	if err := s.ensureGroupAllowsBatchImage(ctx, owner.GroupID); err != nil {
		return nil, err
placeholder

	modelsByProvider := make(map[string]map[string]struct{placeholder)
	for _, providerName := range batchImageProviderSelectionOrder("") {
		provider, ok := s.ProviderRegistry.Get(providerName)
		if !ok || provider == nil {
			continue
	placeholder
		accounts, err := s.listCandidateAccounts(ctx, owner.GroupID, batchImageProviderPlatform(providerName))
		if err != nil {
			return nil, err
	placeholder
		for i := range accounts {
			account := accounts[i]
			if !account.IsSchedulable() || !provider.SupportsAccount(&account) {
				continue
		placeholder
			for _, model := range batchImageModelsFromAccountMapping(&account) {
				if _, err := s.Pricing.BatchImageUnitPrice(ctx, &BatchImageJob{Provider: providerName, Model: modelplaceholder); err != nil {
					continue
			placeholder
				if !account.IsModelSupported(model) {
					continue
			placeholder
				if modelsByProvider[providerName] == nil {
					modelsByProvider[providerName] = make(map[string]struct{placeholder)
			placeholder
				modelsByProvider[providerName][model] = struct{placeholder{placeholder
		placeholder
	placeholder
placeholder

	out := make([]BatchImagePublicModel, 0)
	for _, providerName := range batchImageProviderSelectionOrder("") {
		models := make([]string, 0, len(modelsByProvider[providerName]))
		for model := range modelsByProvider[providerName] {
			models = append(models, model)
	placeholder
		sort.Strings(models)
		for _, model := range models {
			out = append(out, BatchImagePublicModel{
				ID:       model,
				Object:   "image.batch.model",
				Provider: providerName,
		placeholder)
	placeholder
placeholder
	return &BatchImagePublicModelsResponse{Object: "list", Data: outplaceholder, nil
placeholder

func (s *BatchImagePublicService) ListItems(ctx context.Context, owner BatchImageOwner, batchID string, query BatchImageItemsQuery) (*BatchImagePublicItemsResponse, error) {
	filter := BatchImageItemFilter{Limit: query.Limit, Offset: parseBatchImageCursor(query.Cursor)placeholder
	switch strings.TrimSpace(query.Status) {
	case "", "all":
	case "succeeded", "success":
		filter.Status = BatchImageItemStatusSuccess
	case "pending":
		filter.Status = BatchImageItemStatusPending
	case "failed":
		filter.Status = BatchImageItemStatusFailed
	default:
		return nil, ErrBatchImageInvalidItems
placeholder
	if filter.Limit <= 0 || filter.Limit > 500 {
		filter.Limit = 100
placeholder
	items, err := s.Repo.ListBatchImageItemsForOwner(ctx, owner.UserID, owner.APIKeyID, batchID, filter)
	if err != nil {
		return nil, err
placeholder
	data := make([]BatchImagePublicItem, 0, len(items))
	for _, item := range items {
		data = append(data, BatchImageItemToPublic(item))
placeholder
	return &BatchImagePublicItemsResponse{
		Object:  "list",
		Data:    data,
		HasMore: len(data) == filter.Limit,
placeholder, nil
placeholder

func (s *BatchImagePublicService) Cancel(ctx context.Context, owner BatchImageOwner, batchID string) (*BatchImagePublicBatch, error) {
	job, err := s.Repo.GetBatchImageJobByBatchIDForOwner(ctx, owner.UserID, owner.APIKeyID, batchID)
	if err != nil {
		return nil, err
placeholder
	if isBatchImageProcessorDoneStatus(job.Status) {
		if job.Status == BatchImageJobStatusFailed || job.Status == BatchImageJobStatusCancelled {
			if err := releaseBatchImageBalanceHold(ctx, s.BillingRepo, job, batchImageDerefString(job.RequestHash)); err != nil {
				s.enqueueBillingRetry(ctx, job.BatchID)
				return nil, ErrBatchImageCancelFailed
		placeholder
			s.invalidateAuthCache(ctx, owner.UserID)
	placeholder
		return BatchImageJobToPublic(job), nil
placeholder
	if job.ProviderJobName != nil && strings.TrimSpace(*job.ProviderJobName) != "" {
		provider, ok := s.ProviderRegistry.Get(job.Provider)
		if !ok || provider == nil {
			return nil, ErrBatchImageUnsupportedProvider
	placeholder
		if job.AccountID == nil {
			return nil, ErrBatchImageCancelFailed
	placeholder
		account, err := s.AccountRepo.GetByID(ctx, *job.AccountID)
		if err != nil {
			return nil, ErrBatchImageCancelFailed
	placeholder
		if err := provider.Cancel(ctx, job, account); err != nil {
			return nil, ErrBatchImageCancelFailed
	placeholder
		if eventErr := s.Repo.AppendBatchImageEvent(ctx, job.BatchID, "job_cancel_requested", map[string]any{"batch_id": job.BatchIDplaceholder); eventErr != nil {
			logger.L().Warn("batch_image.cancel_event_failed",
				zap.String("batch_id", job.BatchID),
				zap.Error(eventErr),
			)
	placeholder
		if s.Queue != nil {
			if err := s.Queue.Enqueue(ctx, job.BatchID); err != nil && !errors.Is(err, ErrBatchImageAlreadyQueued) {
				return nil, ErrBatchImageCancelFailed
		placeholder
	placeholder
		updated, err := s.Repo.GetBatchImageJobByBatchIDForOwner(ctx, owner.UserID, owner.APIKeyID, batchID)
		if err != nil {
			return nil, err
	placeholder
		return BatchImageJobToPublic(updated), nil
placeholder
	if err := s.Repo.TransitionBatchImageJobStatus(ctx, job.BatchID, BatchImageJobStatusCancelled, BatchImageTransitionOptions{
		EventType:    "job_cancelled",
		EventPayload: map[string]any{"batch_id": job.BatchIDplaceholder,
placeholder); err != nil {
		return nil, err
placeholder
	if err := releaseBatchImageBalanceHold(ctx, s.BillingRepo, job, batchImageDerefString(job.RequestHash)); err != nil {
		s.enqueueBillingRetry(ctx, job.BatchID)
		return nil, ErrBatchImageCancelFailed
placeholder
	s.invalidateAuthCache(ctx, owner.UserID)
	updated, err := s.Repo.GetBatchImageJobByBatchIDForOwner(ctx, owner.UserID, owner.APIKeyID, batchID)
	if err != nil {
		return nil, err
placeholder
	return BatchImageJobToPublic(updated), nil
placeholder

func (s *BatchImagePublicService) validateSubmitRequest(req BatchImageSubmitRequest) (BatchImageSubmitRequest, error) {
	req.Model = strings.TrimSpace(req.Model)
	req.TaskName = strings.TrimSpace(req.TaskName)
	req.ParentBatchID = strings.TrimSpace(req.ParentBatchID)
	req.Provider = strings.TrimSpace(req.Provider)
	req.ResponseMimeType = strings.TrimSpace(req.ResponseMimeType)
	req.AspectRatio = strings.TrimSpace(req.AspectRatio)
	req.ImageSize = strings.TrimSpace(req.ImageSize)
	if req.Model == "" {
		return req, ErrBatchImageInvalidModel
placeholder
	if req.TaskName == "" {
		req.TaskName = defaultBatchImageTaskName(time.Now())
placeholder
	if len(req.TaskName) > 255 {
		req.TaskName = truncateBatchImageMessage(req.TaskName, 255)
placeholder
	if req.Provider != "" && !IsSupportedBatchImageProvider(req.Provider) {
		return req, ErrBatchImageUnsupportedProvider
placeholder
	if len(req.Items) == 0 {
		return req, ErrBatchImageInvalidItems
placeholder
	maxItems := s.maxItems()
	if len(req.Items) > maxItems {
		return req, ErrBatchImageInvalidItems
placeholder
	if req.ResponseMimeType == "" {
		req.ResponseMimeType = s.defaultResponseMimeType()
placeholder
	if req.ImageSize == "" {
		req.ImageSize = s.defaultImageSize()
placeholder
	if !strings.EqualFold(req.ImageSize, defaultBatchImageImageSize) {
		return req, ErrBatchImageInvalidItems
placeholder
	req.ImageSize = defaultBatchImageImageSize
	req.Metadata = sanitizeBatchImageMetadata(req.Metadata)

	seen := make(map[string]struct{placeholder, len(req.Items))
	totalReferenceImages := 0
	totalInlineReferenceBytes := 0
	totalOutputImages := 0
	expandedItems := make([]BatchImageSubmitItem, 0, len(req.Items))
	for i := range req.Items {
		req.Items[i].CustomID = strings.TrimSpace(req.Items[i].CustomID)
		if req.Items[i].CustomID == "" {
			req.Items[i].CustomID = fmt.Sprintf("item_%06d", i+1)
	placeholder
		outputCount := req.Items[i].OutputCount
		if outputCount == 0 {
			outputCount = 1
	placeholder
		if outputCount < 1 || outputCount > s.maxOutputImagesPerItem() {
			return req, ErrBatchImageInvalidItems
	placeholder
		totalOutputImages += outputCount
		if totalOutputImages > s.maxOutputImagesPerJob() {
			return req, ErrBatchImageTooManyOutputImages
	placeholder
		req.Items[i].Prompt = strings.TrimSpace(req.Items[i].Prompt)
		if req.Items[i].Prompt == "" {
			return req, ErrBatchImageInvalidItems
	placeholder
		if len(req.Items[i].Prompt) > s.maxPromptChars() {
			return req, ErrBatchImagePromptTooLong
	placeholder
		referenceCount, inlineReferenceBytes, err := normalizeBatchImageReferenceInputs(req.Model, &req.Items[i])
		if err != nil {
			return req, err
	placeholder
		totalReferenceImages += referenceCount * outputCount
		if totalReferenceImages > s.maxReferenceImagesPerJob() {
			return req, ErrBatchImageTooManyReferenceImages
	placeholder
		totalInlineReferenceBytes += inlineReferenceBytes * outputCount
		if totalInlineReferenceBytes > s.maxReferenceInlineBytesPerJob() {
			return req, ErrBatchImageReferenceImagesTooLarge
	placeholder
		for repeatIndex := 1; repeatIndex <= outputCount; repeatIndex++ {
			expanded := req.Items[i]
			expanded.OutputCount = 0
			if outputCount > 1 {
				expanded.CustomID = fmt.Sprintf("%s_%0*d", req.Items[i].CustomID, batchImageRepeatSuffixWidth(outputCount), repeatIndex)
		placeholder
			if _, ok := seen[expanded.CustomID]; ok {
				return req, ErrBatchImageDuplicateCustomIDInRequest
		placeholder
			seen[expanded.CustomID] = struct{placeholder{placeholder
			expandedItems = append(expandedItems, expanded)
	placeholder
placeholder
	req.Items = expandedItems
	return req, nil
placeholder

func normalizeBatchImageReferenceInputs(model string, item *BatchImageSubmitItem) (int, int, error) {
	if item == nil || len(item.ReferenceImages) == 0 {
		return 0, 0, nil
placeholder
	maxRefs := maxBatchImageReferenceImagesForModel(model)
	if maxRefs <= 0 || len(item.ReferenceImages) > maxRefs {
		return 0, 0, ErrBatchImageTooManyReferenceImages
placeholder
	out := make([]BatchImageReferenceInput, 0, len(item.ReferenceImages))
	inlineBytes := 0
	for _, ref := range item.ReferenceImages {
		ref.ID = truncateBatchImageMessage(strings.TrimSpace(ref.ID), 80)
		ref.Type = truncateBatchImageMessage(strings.TrimSpace(ref.Type), 40)
		ref.MimeType = normalizeBatchImageReferenceMimeType(ref.MimeType)
		ref.FileURI = strings.TrimSpace(ref.FileURI)
		if ref.MimeType == "" {
			return 0, 0, ErrBatchImageInvalidReferenceImage
	placeholder
		if len(ref.Data) == 0 && ref.FileURI == "" {
			return 0, 0, ErrBatchImageInvalidReferenceImage
	placeholder
		if len(ref.Data) > 0 && ref.FileURI != "" {
			return 0, 0, ErrBatchImageInvalidReferenceImage
	placeholder
		if len(ref.Data) > maxBatchImageReferenceImageBytes {
			return 0, 0, ErrBatchImageInvalidReferenceImage
	placeholder
		if ref.FileURI != "" && !strings.HasPrefix(ref.FileURI, "gs://") {
			return 0, 0, ErrBatchImageInvalidReferenceImage
	placeholder
		inlineBytes += len(ref.Data)
		out = append(out, ref)
placeholder
	item.ReferenceImages = out
	return len(out), inlineBytes, nil
placeholder

func normalizeBatchImageReferenceMimeType(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "image/jpeg", "image/jpg":
		return "image/jpeg"
	case "image/png":
		return "image/png"
	case "image/webp":
		return "image/webp"
	default:
		return ""
placeholder
placeholder

func batchImageRepeatSuffixWidth(count int) int {
	if count < 10 {
		return 2
placeholder
	return len(strconv.Itoa(count))
placeholder

func maxBatchImageReferenceImagesForModel(model string) int {
	model = strings.ToLower(strings.TrimSpace(model))
	if strings.Contains(model, "pro-image") {
		return 14
placeholder
	if strings.Contains(model, "flash-image") {
		return 3
placeholder
	return 0
placeholder

func (s *BatchImagePublicService) selectProviderAndAccount(ctx context.Context, owner BatchImageOwner, requestedProvider, model string) (BatchImageProvider, *Account, error) {
	providers := batchImageProviderSelectionOrder(requestedProvider)
	for _, providerName := range providers {
		provider, ok := s.ProviderRegistry.Get(providerName)
		if !ok || provider == nil {
			continue
	placeholder
		accounts, err := s.listCandidateAccounts(ctx, owner.GroupID, batchImageProviderPlatform(providerName))
		if err != nil {
			return nil, nil, err
	placeholder
		sort.SliceStable(accounts, func(i, j int) bool {
			if accounts[i].Priority != accounts[j].Priority {
				return accounts[i].Priority > accounts[j].Priority
		placeholder
			return accounts[i].ID < accounts[j].ID
	placeholder)
		for i := range accounts {
			account := accounts[i]
			if !account.IsSchedulable() || !account.IsModelSupported(model) {
				continue
		placeholder
			if provider.SupportsAccount(&account) {
				return provider, &account, nil
		placeholder
	placeholder
placeholder
	if requestedProvider != "" {
		return nil, nil, ErrBatchImageNoAccountAvailable
placeholder
	return nil, nil, ErrBatchImageNoAccountAvailable
placeholder

func (s *BatchImagePublicService) listCandidateAccounts(ctx context.Context, groupID *int64, platform string) ([]Account, error) {
	if s.AccountRepo == nil {
		return nil, ErrBatchImageNoAccountAvailable
placeholder
	if groupID != nil && *groupID > 0 {
		return s.AccountRepo.ListSchedulableByGroupIDAndPlatform(ctx, *groupID, platform)
placeholder
	return s.AccountRepo.ListSchedulableByPlatform(ctx, platform)
placeholder

func (s *BatchImagePublicService) ensureGroupAllowsBatchImage(ctx context.Context, groupID *int64) error {
	if groupID == nil || *groupID <= 0 {
		return nil
placeholder
	if s.GroupRepo == nil {
		return ErrBatchImageSettlementPricingMissing
placeholder
	group, err := s.GroupRepo.GetByIDLite(ctx, *groupID)
	if err != nil || group == nil {
		return ErrBatchImageSettlementPricingMissing
placeholder
	if !group.AllowBatchImageGeneration {
		return ErrBatchImageGroupDisabled
placeholder
	if group.Platform != PlatformGemini {
		return ErrBatchImageGroupDisabled
placeholder
	return nil
placeholder

func (s *BatchImagePublicService) resolvePricingSnapshot(ctx context.Context, owner BatchImageOwner, req BatchImageSubmitRequest, provider string, account *Account) (*BatchImagePricingSnapshot, error) {
	unit := -1.0
	groupMultiplier := 1.0
	discountMultiplier := defaultBatchImageDiscountMultiplier
	holdMultiplier := defaultBatchImageHoldMultiplier
	if owner.GroupID != nil && *owner.GroupID > 0 {
		if s.GroupRepo == nil {
			return nil, ErrBatchImageSettlementPricingMissing
	placeholder
		group, err := s.GroupRepo.GetByIDLite(ctx, *owner.GroupID)
		if err != nil || group == nil {
			return nil, ErrBatchImageSettlementPricingMissing
	placeholder
		if !group.AllowBatchImageGeneration {
			return nil, ErrBatchImageGroupDisabled
	placeholder
		groupDefaultMultiplier := group.RateMultiplier
		if groupDefaultMultiplier < 0 {
			groupDefaultMultiplier = 0
	placeholder
		effectiveGroupMultiplier := groupDefaultMultiplier
		if s.UserGroupRateRepo != nil {
			userRate, rateErr := s.UserGroupRateRepo.GetByUserAndGroup(ctx, owner.UserID, group.ID)
			if rateErr != nil {
				return nil, ErrBatchImageSettlementPricingMissing
		placeholder
			if userRate != nil {
				effectiveGroupMultiplier = *userRate
		placeholder
	placeholder
		groupMultiplier = effectiveGroupMultiplier
		if group.ImageRateIndependent {
			groupMultiplier = group.ImageRateMultiplier
	placeholder
		if groupMultiplier < 0 {
			groupMultiplier = 0
	placeholder
		discountMultiplier = group.BatchImageDiscountMultiplier
		if discountMultiplier < 0 {
			discountMultiplier = 0
	placeholder
		if group.BatchImageHoldMultiplier >= 0 {
			holdMultiplier = group.BatchImageHoldMultiplier
	placeholder
		if configuredUnit := group.GetImagePrice(req.ImageSize); configuredUnit != nil && *configuredUnit >= 0 {
			unit = *configuredUnit
	placeholder
placeholder
	if unit < 0 {
		if s.Pricing == nil {
			return nil, ErrBatchImageSettlementPricingMissing
	placeholder
		resolvedUnit, err := s.Pricing.BatchImageUnitPrice(ctx, &BatchImageJob{Provider: provider, Model: req.Modelplaceholder)
		if err != nil || resolvedUnit < 0 {
			return nil, ErrBatchImageSettlementPricingMissing
	placeholder
		unit = resolvedUnit
placeholder
	// 定价不变式：hold 比例不得低于 discount 比例，否则成功率足够高时
	// actualCost > holdAmount，结算永远失败、冻结余额无法解冻。
	// 管理端已校验新配置，此处兜底钳制存量脏数据。
	if holdMultiplier < discountMultiplier {
		logger.L().Warn("batch_image.hold_multiplier_below_discount_clamped",
			zap.Float64("hold_multiplier", holdMultiplier),
			zap.Float64("discount_multiplier", discountMultiplier),
		)
		holdMultiplier = discountMultiplier
placeholder
	accountMultiplier := 1.0
	if account != nil {
		accountMultiplier = account.BillingRateMultiplier()
placeholder
	if accountMultiplier < 0 {
		accountMultiplier = 0
placeholder
	standardUnitPrice := unit * groupMultiplier * accountMultiplier
	billableUnitPrice := standardUnitPrice * discountMultiplier
	holdUnitPrice := standardUnitPrice * holdMultiplier
	return &BatchImagePricingSnapshot{
		BaseUnitPrice:           unit,
		GroupRateMultiplier:     groupMultiplier,
		AccountRateMultiplier:   accountMultiplier,
		BatchDiscountMultiplier: discountMultiplier,
		HoldMultiplier:          holdMultiplier,
		BillableUnitPrice:       billableUnitPrice,
		HoldUnitPrice:           holdUnitPrice,
		EstimatedCost:           billableUnitPrice * float64(len(req.Items)),
		HoldAmount:              holdUnitPrice * float64(len(req.Items)),
placeholder, nil
placeholder

func (s *BatchImagePublicService) enabled() bool {
	return s != nil && s.Repo != nil && s.AccountRepo != nil && s.Config != nil && s.Config.BatchImage.Enabled
placeholder

func (s *BatchImagePublicService) invalidateAuthCache(ctx context.Context, userID int64) {
	if s != nil && s.AuthCache != nil && userID > 0 {
		s.AuthCache.InvalidateAuthCacheByUserID(ctx, userID)
placeholder
placeholder

func (s *BatchImagePublicService) maxItems() int {
	if s != nil && s.Config != nil && s.Config.BatchImage.MaxItemsPerJobDefault > 0 {
		return s.Config.BatchImage.MaxItemsPerJobDefault
placeholder
	return defaultBatchImageMaxItems
placeholder

func (s *BatchImagePublicService) maxOutputImagesPerJob() int {
	if s != nil && s.Config != nil && s.Config.BatchImage.MaxOutputImagesPerJob > 0 {
		return s.Config.BatchImage.MaxOutputImagesPerJob
placeholder
	return defaultBatchImageMaxOutputImages
placeholder

func (s *BatchImagePublicService) maxOutputImagesPerItem() int {
	if s != nil && s.Config != nil && s.Config.BatchImage.MaxOutputImagesPerItem > 0 {
		return s.Config.BatchImage.MaxOutputImagesPerItem
placeholder
	return defaultBatchImageMaxOutputCount
placeholder

func (s *BatchImagePublicService) maxPromptChars() int {
	if s != nil && s.Config != nil && s.Config.BatchImage.MaxPromptCharsPerItem > 0 {
		return s.Config.BatchImage.MaxPromptCharsPerItem
placeholder
	return defaultBatchImageMaxPromptChars
placeholder

func (s *BatchImagePublicService) maxReferenceImagesPerJob() int {
	if s != nil && s.Config != nil && s.Config.BatchImage.MaxReferenceImagesPerJob > 0 {
		return s.Config.BatchImage.MaxReferenceImagesPerJob
placeholder
	return defaultBatchImageMaxReferenceImages
placeholder

func (s *BatchImagePublicService) maxReferenceInlineBytesPerJob() int {
	if s != nil && s.Config != nil && s.Config.BatchImage.MaxReferenceInlineBytesPerJob > 0 {
		return s.Config.BatchImage.MaxReferenceInlineBytesPerJob
placeholder
	return defaultBatchImageMaxReferenceBytes
placeholder

func (s *BatchImagePublicService) defaultResponseMimeType() string {
	if s != nil && s.Config != nil && strings.TrimSpace(s.Config.BatchImage.DefaultResponseMimeType) != "" {
		return strings.TrimSpace(s.Config.BatchImage.DefaultResponseMimeType)
placeholder
	return defaultBatchImageResponseMime
placeholder

func (s *BatchImagePublicService) defaultImageSize() string {
	if s != nil && s.Config != nil && strings.TrimSpace(s.Config.BatchImage.DefaultImageSize) != "" {
		return strings.TrimSpace(s.Config.BatchImage.DefaultImageSize)
placeholder
	return defaultBatchImageImageSize
placeholder

func BatchImageJobToPublic(job *BatchImageJob) *BatchImagePublicBatch {
	if job == nil {
		return nil
placeholder
	holdAmount := job.EstimatedCost
	if job.HoldAmount != nil {
		holdAmount = *job.HoldAmount
placeholder
	return &BatchImagePublicBatch{
		ID:              job.BatchID,
		Object:          "image.batch",
		TaskName:        batchImagePublicTaskName(job),
		ParentBatchID:   job.ParentBatchID,
		Status:          PublicBatchImageStatus(job.Status),
		Model:           job.Model,
		Provider:        job.Provider,
		ItemCount:       job.ItemCount,
		SuccessCount:    job.SuccessCount,
		FailCount:       job.FailCount,
		EstimatedCost:   job.EstimatedCost,
		HoldAmount:      holdAmount,
		ActualCost:      job.ActualCost,
		CreatedAt:       job.CreatedAt.Unix(),
		SubmittedAt:     batchImageUnixPtr(job.SubmittedAt),
		SettledAt:       batchImageUnixPtr(job.SettledAt),
		DownloadedAt:    batchImageUnixPtr(job.DownloadedAt),
		OutputDeletedAt: batchImageUnixPtr(job.OutputDeletedAt),
placeholder
placeholder

func BatchImageItemToPublic(item *BatchImageItem) BatchImagePublicItem {
	out := BatchImagePublicItem{
		CustomID:      item.CustomID,
		Status:        "failed",
		PromptPreview: item.PromptPreview,
		MimeType:      item.MimeType,
		FileExtension: item.FileExtension,
		ImageCount:    item.ImageCount,
placeholder
	if item.Status == BatchImageItemStatusPending {
		out.Status = "pending"
		return out
placeholder
	if item.Status == BatchImageItemStatusSuccess {
		out.Status = "succeeded"
		return out
placeholder
	out.Error = &BatchImagePublicError{
		Code:    batchImageDerefString(item.ErrorCode),
		Message: sanitizeBatchImagePublicMessage(batchImageDerefString(item.ErrorMessage)),
		Source:  batchImageItemErrorSource(item),
placeholder
	return out
placeholder

func batchImageItemErrorSource(item *BatchImageItem) string {
	if item == nil || item.ErrorCode == nil {
		return ""
placeholder
	code := strings.TrimSpace(*item.ErrorCode)
	if batchImageDerefString(item.ProviderSourceObject) != "" {
		return "provider"
placeholder
	switch code {
	case "EMPTY_IMAGE_OUTPUT", "PROVIDER_ITEM_FAILED":
		return "provider"
	case "INDEX_OUTPUT_MISSING", "INDEX_PARSE_FAILED", "DUPLICATE_CUSTOM_ID_IN_OUTPUT":
		return "system"
	default:
		return ""
placeholder
placeholder

func PublicBatchImageStatus(status string) string {
	switch status {
	case BatchImageJobStatusCreated, BatchImageJobStatusUploading, BatchImageJobStatusSubmitted:
		return "queued"
	case BatchImageJobStatusRunning:
		return "running"
	case BatchImageJobStatusIndexing:
		return "processing_results"
	case BatchImageJobStatusSettling:
		return "settling"
	case BatchImageJobStatusCompleted:
		return "completed"
	case BatchImageJobStatusFailed:
		return "failed"
	case BatchImageJobStatusCancelled:
		return "cancelled"
	case BatchImageJobStatusOutputDeleted:
		return "output_deleted"
	default:
		return status
placeholder
placeholder

func HashBatchImageSubmitRequest(req BatchImageSubmitRequest) string {
	req.Metadata = sanitizeBatchImageMetadata(req.Metadata)
	b, _ := json.Marshal(req)
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
placeholder

func batchImageProviderPlatform(provider string) string {
	switch provider {
	case BatchImageProviderGeminiAPI, BatchImageProviderVertex:
		return PlatformGemini
	default:
		return PlatformGemini
placeholder
placeholder

func batchImageProviderSelectionOrder(requestedProvider string) []string {
	if strings.TrimSpace(requestedProvider) != "" {
		return []string{strings.TrimSpace(requestedProvider)placeholder
placeholder
	return []string{BatchImageProviderGeminiAPI, BatchImageProviderVertexplaceholder
placeholder

func batchImageModelsFromAccountMapping(account *Account) []string {
	if account == nil {
		return nil
placeholder
	mapping := account.GetModelMapping()
	if len(mapping) == 0 {
		return nil
placeholder
	models := make(map[string]struct{placeholder)
	for model := range mapping {
		model = strings.TrimSpace(model)
		if model == "" {
			continue
	placeholder
		if strings.ContainsAny(model, "*?") {
			for _, candidate := range defaultBatchImageModelCandidates() {
				if matchWildcard(model, candidate) {
					models[candidate] = struct{placeholder{placeholder
			placeholder
		placeholder
			continue
	placeholder
		models[model] = struct{placeholder{placeholder
placeholder
	out := make([]string, 0, len(models))
	for model := range models {
		out = append(out, model)
placeholder
	sort.Strings(out)
	return out
placeholder

func defaultBatchImageModelCandidates() []string {
	return []string{
		"gemini-2.0-flash-exp-image-generation",
		"gemini-2.5-flash-image",
		"gemini-3-pro-image",
		"gemini-3-pro-image-preview",
		"gemini-3.1-flash-image",
		"gemini-3.1-flash-image-preview",
		"gemini-3.1-flash-lite-image",
placeholder
placeholder

func batchImageGCSRef(provider, ref string) string {
	if provider == BatchImageProviderVertex && strings.HasPrefix(strings.TrimSpace(ref), "gs://") {
		return strings.TrimSpace(ref)
placeholder
	return ""
placeholder

func batchImageProviderSubmitPublicError(err error) error {
	reason := strings.TrimSpace(infraerrors.Reason(err))
	switch reason {
	case "VERTEX_MANAGED_GCS_BUCKET_MISSING":
		return ErrBatchImageVertexGCSBucketMissing
	case "BATCH_IMAGE_PROVIDER_MISSING_API_KEY":
		return ErrBatchImageProviderMissingAPIKey
	case "BATCH_IMAGE_PROVIDER_MISSING_SERVICE_ACCOUNT":
		return ErrBatchImageProviderMissingServiceAccount
	case "BATCH_IMAGE_PROVIDER_UNSUPPORTED_ACCOUNT":
		return ErrBatchImageProviderUnsupportedAccount
	default:
		return ErrBatchImageProviderSubmitFailed
placeholder
placeholder

func batchImagePublicTaskName(job *BatchImageJob) string {
	if job == nil {
		return ""
placeholder
	if strings.TrimSpace(job.TaskName) != "" {
		return strings.TrimSpace(job.TaskName)
placeholder
	return defaultBatchImageTaskName(job.CreatedAt)
placeholder

func defaultBatchImageTaskName(now time.Time) string {
	if now.IsZero() {
		now = time.Now()
placeholder
	return now.Format("2006-01-02 15:04:05")
placeholder

func batchImageProviderSubmitRecordCode(err error) string {
	reason := strings.TrimSpace(infraerrors.Reason(err))
	if reason == "" || reason == "BATCH_IMAGE_PROVIDER_SUBMIT_FAILED" {
		return "PROVIDER_SUBMIT_FAILED"
placeholder
	return reason
placeholder

func parseBatchImageListTime(raw string) *time.Time {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
placeholder
	if unix, err := strconv.ParseInt(raw, 10, 64); err == nil && unix > 0 {
		t := time.Unix(unix, 0)
		return &t
placeholder
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return &t
placeholder
	if t, err := time.Parse("2006-01-02", raw); err == nil {
		return &t
placeholder
	return nil
placeholder

func sanitizeBatchImageMetadata(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
placeholder
	keys := make([]string, 0, len(in))
	for k := range in {
		keys = append(keys, k)
placeholder
	sort.Strings(keys)
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		key := strings.TrimSpace(k)
		if key == "" || len(key) > 64 {
			continue
	placeholder
		value := strings.TrimSpace(in[k])
		if len(value) > 256 {
			value = value[:256]
	placeholder
		out[key] = value
		if len(out) >= 20 {
			break
	placeholder
placeholder
	return out
placeholder

func sanitizeBatchImagePublicMessage(message string) string {
	message = strings.TrimSpace(message)
	for _, marker := range []string{"gs://", "files/", "projects/"placeholder {
		if strings.Contains(message, marker) {
			message = "upstream provider operation failed"
			break
	placeholder
placeholder
	if len(message) > maxBatchImagePublicErrorChars {
		message = message[:maxBatchImagePublicErrorChars]
placeholder
	return message
placeholder

func batchImageUnixPtr(t *time.Time) *int64 {
	if t == nil {
		return nil
placeholder
	v := t.Unix()
	return &v
placeholder

func parseBatchImageCursor(cursor string) int {
	offset, err := strconv.Atoi(strings.TrimSpace(cursor))
	if err != nil || offset < 0 {
		return 0
placeholder
	return offset
placeholder
