package service

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"go.uber.org/zap"
)

const (
	BatchImageParsedStatusSucceeded = "succeeded"
	BatchImageParsedStatusFailed    = "failed"

	defaultBatchImageProcessorRequeue = 30 * time.Second
	batchImageProviderErrorRequeue    = time.Minute
	batchImageMaxErrorMessageLength   = 1000
)

type BatchImageAccountResolver interface {
	ResolveBatchImageAccount(ctx context.Context, accountID int64) (*Account, error)
placeholder

type BatchImageAccountLookup interface {
	GetByID(ctx context.Context, id int64) (*Account, error)
placeholder

type BatchImageAccountRepositoryResolver struct {
	Repo BatchImageAccountLookup
placeholder

func (r *BatchImageAccountRepositoryResolver) ResolveBatchImageAccount(ctx context.Context, accountID int64) (*Account, error) {
	if r == nil || r.Repo == nil {
		return nil, ErrAccountNotFound
placeholder
	return r.Repo.GetByID(ctx, accountID)
placeholder

type BatchImageProviderProcessor struct {
	Repo             BatchImageRepository
	ProviderRegistry *BatchImageProviderRegistry
	AccountResolver  BatchImageAccountResolver
	Indexer          *BatchImageResultIndexer
	BillingRepo      UsageBillingRepository
	AuthCache        APIKeyAuthCacheInvalidator
	DefaultRequeue   time.Duration
placeholder

func (p *BatchImageProviderProcessor) Process(ctx context.Context, batchID string) (BatchImageProcessResult, error) {
	if p == nil || p.Repo == nil || p.ProviderRegistry == nil || p.AccountResolver == nil {
		return BatchImageProcessResult{placeholder, infraerrors.New(http.StatusInternalServerError, "BATCH_IMAGE_PROCESSOR_NOT_CONFIGURED", "batch image processor is not configured")
placeholder

	job, err := p.Repo.GetBatchImageJobByBatchID(ctx, batchID)
	if err != nil {
		return BatchImageProcessResult{placeholder, err
placeholder
	if isBatchImageProcessorDoneStatus(job.Status) {
		if err := p.releaseTerminalHold(ctx, job); err != nil {
			return BatchImageProcessResult{placeholder, err
	placeholder
		return BatchImageProcessResult{Terminal: trueplaceholder, nil
placeholder

	provider, ok := p.ProviderRegistry.Get(job.Provider)
	if !ok || provider == nil {
		return BatchImageProcessResult{placeholder, ErrBatchImageUnsupportedProvider
placeholder
	if job.AccountID == nil || *job.AccountID <= 0 {
		return BatchImageProcessResult{placeholder, ErrBatchImageMissingAccountID
placeholder
	account, err := p.AccountResolver.ResolveBatchImageAccount(ctx, *job.AccountID)
	if err != nil {
		return BatchImageProcessResult{placeholder, err
placeholder
	if !provider.SupportsAccount(account) {
		return BatchImageProcessResult{placeholder, ErrBatchImageProviderUnsupportedAccount
placeholder
	if strings.TrimSpace(batchImageDerefString(job.ProviderJobName)) == "" {
		return BatchImageProcessResult{placeholder, ErrBatchImageMissingProviderJobName
placeholder

	if job.Status == BatchImageJobStatusIndexing {
		return p.indexAndSettle(ctx, job, provider, account)
placeholder

	status, err := provider.Get(ctx, job, account)
	if err != nil {
		logger.L().Warn("batch_image.provider_status_check_failed",
			zap.String("batch_id", job.BatchID),
			zap.String("provider", job.Provider),
			zap.String("provider_job_name", batchImageDerefString(job.ProviderJobName)),
			zap.Error(err),
		)
		return BatchImageProcessResult{RequeueAfter: batchImageProviderErrorRequeueplaceholder, nil
placeholder
	if status == nil {
		return BatchImageProcessResult{RequeueAfter: p.requeueDelay(0)placeholder, nil
placeholder
	if err := p.persistProviderOutputRef(ctx, job, status.ProviderOutputRef); err != nil {
		return BatchImageProcessResult{placeholder, err
placeholder

	switch status.InternalState {
	case BatchProviderStateQueued:
		return BatchImageProcessResult{RequeueAfter: p.requeueDelay(status.SuggestedRequeueAfter)placeholder, nil
	case BatchProviderStateRunning:
		if job.Status != BatchImageJobStatusRunning {
			if err := p.Repo.TransitionBatchImageJobStatus(ctx, job.BatchID, BatchImageJobStatusRunning, BatchImageTransitionOptions{
				EventType:    "provider_status_checked",
				EventPayload: map[string]any{"provider_state": status.RawStateplaceholder,
		placeholder); err != nil {
				return BatchImageProcessResult{placeholder, err
		placeholder
			job.Status = BatchImageJobStatusRunning
	placeholder
		return BatchImageProcessResult{RequeueAfter: p.requeueDelay(status.SuggestedRequeueAfter)placeholder, nil
	case BatchProviderStateSucceeded:
		if job.Status != BatchImageJobStatusIndexing {
			if err := p.Repo.TransitionBatchImageJobStatus(ctx, job.BatchID, BatchImageJobStatusIndexing, BatchImageTransitionOptions{
				EventType:    "indexing_started",
				EventPayload: map[string]any{"provider_state": status.RawStateplaceholder,
		placeholder); err != nil {
				return BatchImageProcessResult{placeholder, err
		placeholder
			job.Status = BatchImageJobStatusIndexing
	placeholder
		return p.indexAndSettle(ctx, job, provider, account)
	case BatchProviderStateFailed, BatchProviderStateExpired:
		code := strings.TrimSpace(status.ErrorCode)
		if code == "" && status.InternalState == BatchProviderStateExpired {
			code = "PROVIDER_BATCH_EXPIRED"
	placeholder
		if code == "" {
			code = "PROVIDER_BATCH_FAILED"
	placeholder
		msg := truncateBatchImageMessage(status.ErrorMessage, batchImageMaxErrorMessageLength)
		if err := p.Repo.TransitionBatchImageJobStatus(ctx, job.BatchID, BatchImageJobStatusFailed, BatchImageTransitionOptions{
			EventType:    "job_failed",
			EventPayload: map[string]any{"provider_state": status.RawState, "error_code": codeplaceholder,
			ErrorCode:    batchImageStringPtr(code),
			ErrorMessage: batchImageOptionalStringPtr(msg),
	placeholder); err != nil {
			return BatchImageProcessResult{placeholder, err
	placeholder
		job.Status = BatchImageJobStatusFailed
		if err := p.releaseTerminalHold(ctx, job); err != nil {
			return BatchImageProcessResult{placeholder, err
	placeholder
		return BatchImageProcessResult{Terminal: trueplaceholder, nil
	case BatchProviderStateCancelled:
		if err := p.Repo.TransitionBatchImageJobStatus(ctx, job.BatchID, BatchImageJobStatusCancelled, BatchImageTransitionOptions{
			EventType:    "job_failed",
			EventPayload: map[string]any{"provider_state": status.RawState, "error_code": "PROVIDER_BATCH_CANCELLED"placeholder,
	placeholder); err != nil {
			return BatchImageProcessResult{placeholder, err
	placeholder
		job.Status = BatchImageJobStatusCancelled
		if err := p.releaseTerminalHold(ctx, job); err != nil {
			return BatchImageProcessResult{placeholder, err
	placeholder
		return BatchImageProcessResult{Terminal: trueplaceholder, nil
	default:
		return BatchImageProcessResult{RequeueAfter: p.requeueDelay(status.SuggestedRequeueAfter)placeholder, nil
placeholder
placeholder

func (p *BatchImageProviderProcessor) indexAndSettle(ctx context.Context, job *BatchImageJob, provider BatchImageProvider, account *Account) (BatchImageProcessResult, error) {
	indexer := p.Indexer
	if indexer == nil {
		indexer = &BatchImageResultIndexer{Repo: p.Repoplaceholder
placeholder
	if indexer.Repo == nil {
		indexer.Repo = p.Repo
placeholder

	result, err := indexer.Index(ctx, job, provider, account)
	if err != nil {
		if errors.Is(err, ErrBatchImageIndexOutputMissing) {
			return BatchImageProcessResult{placeholder, err
	placeholder
		code := "INDEX_PARSE_FAILED"
		if errors.Is(err, ErrBatchImageDuplicateCustomID) {
			code = "DUPLICATE_CUSTOM_ID_IN_OUTPUT"
	placeholder
		msg := truncateBatchImageMessage(err.Error(), batchImageMaxErrorMessageLength)
		transitionErr := p.Repo.TransitionBatchImageJobStatus(ctx, job.BatchID, BatchImageJobStatusFailed, BatchImageTransitionOptions{
			EventType:    "indexing_failed",
			EventPayload: map[string]any{"error_code": codeplaceholder,
			ErrorCode:    batchImageStringPtr(code),
			ErrorMessage: batchImageOptionalStringPtr(msg),
	placeholder)
		if transitionErr != nil {
			return BatchImageProcessResult{placeholder, transitionErr
	placeholder
		job.Status = BatchImageJobStatusFailed
		if err := p.releaseTerminalHold(ctx, job); err != nil {
			return BatchImageProcessResult{placeholder, err
	placeholder
		return BatchImageProcessResult{Terminal: trueplaceholder, nil
placeholder

	if err := p.Repo.TransitionBatchImageJobStatus(ctx, job.BatchID, BatchImageJobStatusSettling, BatchImageTransitionOptions{
		EventType: "indexing_completed",
		EventPayload: map[string]any{
			"success_count": result.SuccessCount,
			"fail_count":    result.FailCount,
			"total_count":   result.TotalCount,
	placeholder,
placeholder); err != nil {
		return BatchImageProcessResult{placeholder, err
placeholder
	return BatchImageProcessResult{RequeueAfter: time.Millisecondplaceholder, nil
placeholder

func (p *BatchImageProviderProcessor) releaseTerminalHold(ctx context.Context, job *BatchImageJob) error {
	if p == nil || job == nil {
		return nil
placeholder
	if job.Status != BatchImageJobStatusFailed && job.Status != BatchImageJobStatusCancelled {
		return nil
placeholder
	if err := releaseBatchImageBalanceHold(ctx, p.BillingRepo, job, batchImageDerefString(job.RequestHash)); err != nil {
		return err
placeholder
	if p.AuthCache != nil && job.UserID > 0 {
		p.AuthCache.InvalidateAuthCacheByUserID(ctx, job.UserID)
placeholder
	return nil
placeholder

func (p *BatchImageProviderProcessor) persistProviderOutputRef(ctx context.Context, job *BatchImageJob, ref string) error {
	ref = strings.TrimSpace(ref)
	if ref == "" || job == nil || batchImageDerefString(job.ProviderOutputRef) == ref {
		return nil
placeholder
	if err := p.Repo.UpdateBatchImageJobProviderOutputRef(ctx, job.BatchID, ref); err != nil {
		return err
placeholder
	job.ProviderOutputRef = &ref
	return nil
placeholder

func (p *BatchImageProviderProcessor) requeueDelay(suggested time.Duration) time.Duration {
	if suggested > 0 {
		return suggested
placeholder
	if p != nil && p.DefaultRequeue > 0 {
		return p.DefaultRequeue
placeholder
	return defaultBatchImageProcessorRequeue
placeholder

func isBatchImageProcessorDoneStatus(status string) bool {
	if status == BatchImageJobStatusSettling {
		return true
placeholder
	return IsTerminalBatchImageJobStatus(status)
placeholder

type BatchImageIndexResult struct {
	SuccessCount int
	FailCount    int
	TotalCount   int
placeholder

type BatchImageResultIndexer struct {
	Repo BatchImageRepository
placeholder

func (i *BatchImageResultIndexer) Index(ctx context.Context, job *BatchImageJob, provider BatchImageProvider, account *Account) (*BatchImageIndexResult, error) {
	if i == nil || i.Repo == nil || job == nil || provider == nil {
		return nil, ErrBatchImageIndexOutputMissing
placeholder
	r, _, err := provider.OpenResult(ctx, job, account)
	if err != nil {
		return nil, ErrBatchImageIndexOutputMissing.WithCause(err)
placeholder
	defer func() { _ = r.Close() placeholder()

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)

	seen := make(map[string]int)
	var items []CreateBatchImageItemParams
	result := &BatchImageIndexResult{placeholder
	lineNumber := 0
	now := time.Now()
	sourceObject := batchImageDerefString(job.ProviderOutputRef)
	if sourceObject == "" {
		sourceObject = batchImageDerefString(job.ProviderJobName)
placeholder

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
	placeholder
		parsed, err := ParseBatchImageResultLine([]byte(line), lineNumber)
		if err != nil {
			return nil, err
	placeholder
		if firstLine, ok := seen[parsed.CustomID]; ok {
			return nil, ErrBatchImageDuplicateCustomID.WithCause(fmt.Errorf("custom id %q duplicated at lines %d and %d", parsed.CustomID, firstLine, lineNumber))
	placeholder
		seen[parsed.CustomID] = lineNumber

		lineNo := parsed.SourceLineNumber
		item := CreateBatchImageItemParams{
			JobID:                job.BatchID,
			CustomID:             parsed.CustomID,
			Status:               BatchImageItemStatusFailed,
			ProviderSourceObject: batchImageOptionalStringPtr(sourceObject),
			SourceLineNumber:     &lineNo,
			ImageCount:           parsed.ImageCount,
			IndexedAt:            &now,
	placeholder
		if parsed.Status == BatchImageParsedStatusSucceeded {
			item.Status = BatchImageItemStatusSuccess
			item.MimeType = batchImageOptionalStringPtr(parsed.MimeType)
			item.FileExtension = batchImageOptionalStringPtr(parsed.FileExtension)
			result.SuccessCount++
	placeholder else {
			item.ErrorCode = batchImageOptionalStringPtr(parsed.ErrorCode)
			item.ErrorMessage = batchImageOptionalStringPtr(parsed.ErrorMessage)
			result.FailCount++
	placeholder
		items = append(items, item)
		result.TotalCount++
placeholder
	if err := scanner.Err(); err != nil {
		if errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, ErrBatchImageIndexParseFailed.WithCause(err)
	placeholder
		return nil, err
placeholder
	if result.TotalCount == 0 {
		return nil, ErrBatchImageIndexNoResultLines
placeholder
	if err := i.Repo.ReplaceBatchImageItemsForJob(ctx, job.BatchID, items, BatchImageCounts{
		SuccessCount: result.SuccessCount,
		FailCount:    result.FailCount,
placeholder); err != nil {
		return nil, err
placeholder
	return result, nil
placeholder

type ParsedBatchImageResult struct {
	CustomID      string
	Status        string
	MimeType      string
	FileExtension string
	ImageCount    int

	ErrorCode    string
	ErrorMessage string

	SourceLineNumber int
placeholder

func ParseBatchImageResultLine(line []byte, lineNumber int) (*ParsedBatchImageResult, error) {
	var obj map[string]any
	if err := json.Unmarshal(line, &obj); err != nil {
		return nil, ErrBatchImageIndexParseFailed.WithCause(fmt.Errorf("line %d: %w", lineNumber, err))
placeholder

	customID := batchImageFirstNonEmptyString(
		batchImageMapString(obj, "key"),
		batchImageMapString(obj, "custom_id"),
		batchImageMapString(obj, "customId"),
		batchImageNestedString(obj, "request", "key"),
	)
	if customID == "" {
		return nil, ErrBatchImageIndexParseFailed.WithCause(fmt.Errorf("line %d: missing custom id", lineNumber))
placeholder

	parsed := &ParsedBatchImageResult{
		CustomID:         customID,
		SourceLineNumber: lineNumber,
placeholder
	imageCount, mimeType := batchImageFindImageParts(obj)
	if imageCount > 0 {
		parsed.Status = BatchImageParsedStatusSucceeded
		parsed.ImageCount = imageCount
		parsed.MimeType = mimeType
		parsed.FileExtension = batchImageFileExtension(mimeType)
		return parsed, nil
placeholder

	if code, message, ok := batchImageFailureFromProviderFields(obj); ok {
		parsed.Status = BatchImageParsedStatusFailed
		parsed.ErrorCode = code
		parsed.ErrorMessage = truncateBatchImageMessage(message, batchImageMaxErrorMessageLength)
		return parsed, nil
placeholder

	if _, hasResponse := obj["response"]; hasResponse || batchImageHasCandidates(obj) {
		parsed.Status = BatchImageParsedStatusFailed
		parsed.ErrorCode = "EMPTY_IMAGE_OUTPUT"
		parsed.ErrorMessage = "provider response contained no image output"
		return parsed, nil
placeholder

	parsed.Status = BatchImageParsedStatusFailed
	parsed.ErrorCode = "PROVIDER_ITEM_FAILED"
	parsed.ErrorMessage = "provider result line contained no image output"
	return parsed, nil
placeholder

func batchImageFindImageParts(obj map[string]any) (int, string) {
	count, mimeType := batchImageFindImagePartsInCandidates(batchImageNestedAny(obj, "response", "candidates"))
	if count > 0 {
		return count, mimeType
placeholder
	return batchImageFindImagePartsInCandidates(obj["candidates"])
placeholder

func batchImageFindImagePartsInCandidates(raw any) (int, string) {
	candidates, ok := raw.([]any)
	if !ok {
		return 0, ""
placeholder
	count := 0
	firstMime := ""
	for _, candidateRaw := range candidates {
		candidate, ok := candidateRaw.(map[string]any)
		if !ok {
			continue
	placeholder
		partsRaw := batchImageNestedAny(candidate, "content", "parts")
		parts, ok := partsRaw.([]any)
		if !ok {
			continue
	placeholder
		for _, partRaw := range parts {
			part, ok := partRaw.(map[string]any)
			if !ok {
				continue
		placeholder
			inline, ok := firstMap(part["inlineData"], part["inline_data"])
			if !ok {
				continue
		placeholder
			data := strings.TrimSpace(batchImageMapString(inline, "data"))
			mime := batchImageFirstNonEmptyString(batchImageMapString(inline, "mimeType"), batchImageMapString(inline, "mime_type"))
			if data == "" || !strings.HasPrefix(strings.ToLower(strings.TrimSpace(mime)), "image/") {
				continue
		placeholder
			count++
			if firstMime == "" {
				firstMime = strings.TrimSpace(mime)
		placeholder
	placeholder
placeholder
	return count, firstMime
placeholder

func batchImageFailureFromProviderFields(obj map[string]any) (string, string, bool) {
	if status, ok := obj["status"].(map[string]any); ok {
		message := batchImageFirstNonEmptyString(batchImageMapString(status, "message"), batchImageMapString(status, "details"))
		code := batchImageFirstNonEmptyString(batchImageMapString(status, "code"), batchImageMapString(status, "status"))
		return batchImageMapFailureCode(code, message), message, true
placeholder
	if errObj, ok := obj["error"].(map[string]any); ok {
		message := batchImageFirstNonEmptyString(batchImageMapString(errObj, "message"), batchImageMapString(errObj, "details"))
		code := batchImageFirstNonEmptyString(batchImageMapString(errObj, "code"), batchImageMapString(errObj, "status"))
		return batchImageMapFailureCode(code, message), message, true
placeholder
	return "", "", false
placeholder

func batchImageMapFailureCode(code, message string) string {
	text := strings.ToLower(strings.TrimSpace(code + " " + message))
	switch {
	case strings.Contains(text, "safety"), strings.Contains(text, "policy"), strings.Contains(text, "blocked"), strings.Contains(text, "prohibited"):
		return "SAFETY_BLOCKED"
	case strings.Contains(text, "invalid_argument"), strings.Contains(text, "invalid argument"), strings.Contains(text, "bad request"):
		return "INVALID_ARGUMENT"
	case strings.Contains(text, "quota"), strings.Contains(text, "rate"), strings.Contains(text, "resource_exhausted"), strings.Contains(text, "too many requests"):
		return "PROVIDER_RATE_LIMITED"
	default:
		return "PROVIDER_ITEM_FAILED"
placeholder
placeholder

func batchImageFileExtension(mimeType string) string {
	switch strings.ToLower(strings.TrimSpace(mimeType)) {
	case "image/png":
		return "png"
	case "image/jpeg", "image/jpg":
		return "jpg"
	case "image/webp":
		return "webp"
	default:
		return ""
placeholder
placeholder

func batchImageHasCandidates(obj map[string]any) bool {
	if _, ok := obj["candidates"]; ok {
		return true
placeholder
	_, ok := batchImageNestedAny(obj, "response", "candidates").([]any)
	return ok
placeholder

func batchImageMapString(m map[string]any, key string) string {
	if m == nil {
		return ""
placeholder
	switch v := m[key].(type) {
	case string:
		return strings.TrimSpace(v)
	case json.Number:
		return v.String()
	case float64:
		return strconv.FormatInt(int64(v), 10)
	default:
		return ""
placeholder
placeholder

func batchImageNestedString(m map[string]any, keys ...string) string {
	if nested, ok := batchImageNestedAny(m, keys...).(string); ok {
		return strings.TrimSpace(nested)
placeholder
	return ""
placeholder

func batchImageNestedAny(m map[string]any, keys ...string) any {
	var current any = m
	for _, key := range keys {
		cm, ok := current.(map[string]any)
		if !ok {
			return nil
	placeholder
		current = cm[key]
placeholder
	return current
placeholder

func firstMap(values ...any) (map[string]any, bool) {
	for _, value := range values {
		if m, ok := value.(map[string]any); ok {
			return m, true
	placeholder
placeholder
	return nil, false
placeholder

func batchImageFirstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
	placeholder
placeholder
	return ""
placeholder

func batchImageDerefString(v *string) string {
	if v == nil {
		return ""
placeholder
	return strings.TrimSpace(*v)
placeholder

func batchImageStringPtr(v string) *string {
	return &v
placeholder

func batchImageOptionalStringPtr(v string) *string {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
placeholder
	return &v
placeholder

func truncateBatchImageMessage(message string, limit int) string {
	message = strings.TrimSpace(message)
	if limit <= 0 || len(message) <= limit {
		return message
placeholder
	return message[:limit]
placeholder
