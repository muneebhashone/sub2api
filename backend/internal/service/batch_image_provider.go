package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

type BatchImageProvider interface {
	Name() string
	SupportsAccount(account *Account) bool
	Submit(ctx context.Context, job *BatchImageJob, account *Account, input BatchImageInput) (*BatchProviderJob, error)
	Get(ctx context.Context, job *BatchImageJob, account *Account) (*BatchProviderStatus, error)
	Cancel(ctx context.Context, job *BatchImageJob, account *Account) error
	OpenResult(ctx context.Context, job *BatchImageJob, account *Account) (io.ReadCloser, string, error)
	Cleanup(ctx context.Context, job *BatchImageJob, account *Account, target CleanupTarget) error
placeholder

type BatchImageProviderRegistry struct {
	providers map[string]BatchImageProvider
placeholder

func NewBatchImageProviderRegistry(providers ...BatchImageProvider) *BatchImageProviderRegistry {
	r := &BatchImageProviderRegistry{providers: make(map[string]BatchImageProvider, len(providers))placeholder
	for _, provider := range providers {
		if provider == nil || strings.TrimSpace(provider.Name()) == "" {
			continue
	placeholder
		r.providers[provider.Name()] = provider
placeholder
	return r
placeholder

func NewDefaultBatchImageProviderRegistry() *BatchImageProviderRegistry {
	return NewBatchImageProviderRegistry(
		NewGeminiAPIBatchImageProvider(nil),
		NewVertexBatchImageProvider(VertexBatchImageProviderOptions{placeholder, nil, nil, nil),
	)
placeholder

func NewBatchImageProviderRegistryFromConfig(cfg *config.Config) *BatchImageProviderRegistry {
	return NewBatchImageProviderRegistry(
		NewGeminiAPIBatchImageProvider(nil),
		NewVertexBatchImageProviderFromConfig(cfg, nil, nil, nil),
	)
placeholder

func (r *BatchImageProviderRegistry) Get(provider string) (BatchImageProvider, bool) {
	if r == nil {
		return nil, false
placeholder
	p, ok := r.providers[provider]
	return p, ok
placeholder

func (r *BatchImageProviderRegistry) MustGet(provider string) (BatchImageProvider, error) {
	p, ok := r.Get(provider)
	if !ok {
		return nil, ErrBatchImageInvalidProvider
placeholder
	return p, nil
placeholder

type BatchImageInput struct {
	BatchID     string
	Model       string
	DisplayName string
	Items       []BatchImageInputItem

	ResponseMimeType string
	AspectRatio      string
	ImageSize        string

	Metadata map[string]string
placeholder

type BatchImageInputItem struct {
	CustomID string
	Prompt   string

	ReferenceImages []BatchImageReference
placeholder

type BatchImageReference struct {
	MimeType string
	Data     []byte
placeholder

type BatchProviderJob struct {
	ProviderJobName   string
	ProviderInputRef  string
	ProviderOutputRef string
	RawState          string
placeholder

type BatchProviderInternalState string

const (
	BatchProviderStateQueued    BatchProviderInternalState = "queued"
	BatchProviderStateRunning   BatchProviderInternalState = "running"
	BatchProviderStateSucceeded BatchProviderInternalState = "succeeded"
	BatchProviderStateFailed    BatchProviderInternalState = "failed"
	BatchProviderStateCancelled BatchProviderInternalState = "cancelled"
	BatchProviderStateExpired   BatchProviderInternalState = "expired"
)

type BatchProviderStatus struct {
	RawState string

	InternalState BatchProviderInternalState
	Done          bool

	ProviderOutputRef string

	ErrorCode    string
	ErrorMessage string

	SuggestedRequeueAfter time.Duration
placeholder

type CleanupTarget string

const (
	CleanupTargetInput  CleanupTarget = "input"
	CleanupTargetOutput CleanupTarget = "output"
	CleanupTargetAll    CleanupTarget = "all"
)

var (
	ErrBatchImageProviderUnsupportedAccount      = infraerrors.New(http.StatusBadRequest, "BATCH_IMAGE_PROVIDER_UNSUPPORTED_ACCOUNT", "batch image provider does not support this account")
	ErrBatchImageProviderMissingAPIKey           = infraerrors.New(http.StatusBadRequest, "BATCH_IMAGE_PROVIDER_MISSING_API_KEY", "batch image provider account is missing api key")
	ErrBatchImageProviderMissingServiceAccount   = infraerrors.New(http.StatusBadRequest, "BATCH_IMAGE_PROVIDER_MISSING_SERVICE_ACCOUNT", "batch image provider account is missing service account credentials")
	ErrBatchImageProviderMissingJobName          = infraerrors.New(http.StatusBadRequest, "BATCH_IMAGE_PROVIDER_MISSING_JOB_NAME", "batch image provider job name is missing")
	ErrBatchImageProviderMissingResultRef        = infraerrors.New(http.StatusBadRequest, "BATCH_IMAGE_PROVIDER_MISSING_RESULT_REF", "batch image provider result reference is missing")
	ErrBatchImageProviderInlineResultUnsupported = infraerrors.New(http.StatusBadRequest, "GEMINI_INLINE_BATCH_RESULT_UNSUPPORTED", "Gemini inline batch result is not supported")
	ErrBatchImageProviderInvalidInput            = infraerrors.New(http.StatusBadRequest, "BATCH_IMAGE_PROVIDER_INVALID_INPUT", "invalid batch image provider input")
	ErrBatchImageProviderUnsafeCleanupPath       = infraerrors.New(http.StatusBadRequest, "VERTEX_UNSAFE_CLEANUP_PATH", "unsafe batch image cleanup path")
	ErrUnsupportedCleanupTarget                  = infraerrors.New(http.StatusBadRequest, "BATCH_IMAGE_PROVIDER_UNSUPPORTED_CLEANUP_TARGET", "unsupported batch image cleanup target")
)

func batchImageProviderJobName(job *BatchImageJob) string {
	if job == nil || job.ProviderJobName == nil {
		return ""
placeholder
	return strings.TrimSpace(*job.ProviderJobName)
placeholder

func batchImageProviderInputRef(job *BatchImageJob) string {
	if job == nil || job.ProviderInputRef == nil {
		return ""
placeholder
	return strings.TrimSpace(*job.ProviderInputRef)
placeholder

func batchImageProviderOutputRef(job *BatchImageJob) string {
	if job == nil || job.ProviderOutputRef == nil {
		return ""
placeholder
	return strings.TrimSpace(*job.ProviderOutputRef)
placeholder

func batchImageProviderAPIKey(account *Account) string {
	if account == nil {
		return ""
placeholder
	return strings.TrimSpace(account.GetCredential("api_key"))
placeholder

func batchImageProviderInputError(format string, args ...any) error {
	return ErrBatchImageProviderInvalidInput.WithCause(fmt.Errorf(format, args...))
placeholder
