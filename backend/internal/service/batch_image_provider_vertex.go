package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

const (
	defaultVertexBatchRequeueAfter = 30 * time.Second
	defaultVertexBatchLocation     = "global"
	defaultVertexManagedGCSPrefix  = "batch-image/{envplaceholder/{batch_idplaceholder"
)

type VertexBatchImageProviderOptions struct {
	Enabled                bool
	ProjectID              string
	Location               string
	ManagedGCSBucket       string
	ManagedGCSPrefix       string
	Environment            string
	InputRetentionHours    int
	OutputRetentionHours   int
	BatchPredictionBaseURL string
	GCSBaseURL             string
placeholder

func NewVertexBatchImageProviderOptionsFromConfig(cfg *config.Config) VertexBatchImageProviderOptions {
	if cfg == nil {
		return VertexBatchImageProviderOptions{placeholder
placeholder
	return VertexBatchImageProviderOptions{
		Enabled:                cfg.BatchImage.VertexEnabled,
		ProjectID:              cfg.BatchImage.VertexProjectID,
		Location:               cfg.BatchImage.VertexLocation,
		ManagedGCSBucket:       cfg.BatchImage.VertexManagedGCSBucket,
		ManagedGCSPrefix:       cfg.BatchImage.VertexManagedGCSPrefix,
		Environment:            cfg.Log.Environment,
		InputRetentionHours:    cfg.BatchImage.VertexInputRetentionHours,
		OutputRetentionHours:   cfg.BatchImage.VertexOutputRetentionHours,
		BatchPredictionBaseURL: cfg.BatchImage.VertexBatchPredictionBaseURL,
		GCSBaseURL:             cfg.BatchImage.VertexGCSBaseURL,
placeholder
placeholder

type VertexBatchClient interface {
	CreateBatchPredictionJob(ctx context.Context, accessToken string, req VertexCreateBatchPredictionJobRequest) (*VertexBatchPredictionJob, error)
	GetBatchPredictionJob(ctx context.Context, accessToken string, name string) (*VertexBatchPredictionJob, error)
	CancelBatchPredictionJob(ctx context.Context, accessToken string, name string) error
placeholder

type VertexBatchObjectStore interface {
	UploadJSONL(ctx context.Context, accessToken string, uri string, r io.Reader) error
	ListJSONLObjects(ctx context.Context, accessToken string, prefixURI string) ([]string, error)
	OpenObject(ctx context.Context, accessToken string, uri string) (io.ReadCloser, string, error)
	DeleteObject(ctx context.Context, accessToken string, uri string) error
	DeletePrefix(ctx context.Context, accessToken string, prefixURI string) error
placeholder

type VertexCreateBatchPredictionJobRequest struct {
	ProjectID      string                     `json:"-"`
	Location       string                     `json:"-"`
	DisplayName    string                     `json:"displayName"`
	Model          string                     `json:"model"`
	InputConfig    VertexBatchInputConfig     `json:"inputConfig"`
	OutputConfig   VertexBatchOutputConfig    `json:"outputConfig"`
	InstanceConfig *VertexBatchInstanceConfig `json:"instanceConfig,omitempty"`
placeholder

type VertexBatchInputConfig struct {
	InstancesFormat string               `json:"instancesFormat"`
	GCSSource       VertexBatchGCSSource `json:"gcsSource"`
placeholder

type VertexBatchGCSSource struct {
	URIs []string `json:"uris"`
placeholder

type VertexBatchOutputConfig struct {
	PredictionsFormat string                    `json:"predictionsFormat"`
	GCSDestination    VertexBatchGCSDestination `json:"gcsDestination"`
placeholder

type VertexBatchGCSDestination struct {
	OutputURIPrefix string `json:"outputUriPrefix"`
placeholder

type VertexBatchInstanceConfig struct {
	KeyField string `json:"keyField"`
placeholder

type VertexBatchPredictionJob struct {
	Name         string                  `json:"name"`
	DisplayName  string                  `json:"displayName"`
	State        string                  `json:"state"`
	OutputConfig VertexBatchOutputConfig `json:"outputConfig"`
	Error        *VertexBatchJobError    `json:"error"`
placeholder

type VertexBatchJobError struct {
	Code    any    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
placeholder

type VertexBatchImageProvider struct {
	opts        VertexBatchImageProviderOptions
	client      VertexBatchClient
	objectStore VertexBatchObjectStore
	tokenCache  GeminiTokenCache
placeholder

func NewVertexBatchImageProvider(opts VertexBatchImageProviderOptions, client VertexBatchClient, objectStore VertexBatchObjectStore, tokenCache GeminiTokenCache) *VertexBatchImageProvider {
	opts = normalizeVertexBatchImageProviderOptions(opts)
	if client == nil {
		client = NewVertexBatchHTTPClient(opts.BatchPredictionBaseURL, nil)
placeholder
	if objectStore == nil {
		objectStore = NewVertexGCSObjectStore(opts.GCSBaseURL, nil)
placeholder
	return &VertexBatchImageProvider{
		opts:        opts,
		client:      client,
		objectStore: objectStore,
		tokenCache:  tokenCache,
placeholder
placeholder

func NewVertexBatchImageProviderFromConfig(cfg *config.Config, client VertexBatchClient, objectStore VertexBatchObjectStore, tokenCache GeminiTokenCache) *VertexBatchImageProvider {
	return NewVertexBatchImageProvider(NewVertexBatchImageProviderOptionsFromConfig(cfg), client, objectStore, tokenCache)
placeholder

func normalizeVertexBatchImageProviderOptions(opts VertexBatchImageProviderOptions) VertexBatchImageProviderOptions {
	opts.ProjectID = strings.TrimSpace(opts.ProjectID)
	opts.Location = strings.TrimSpace(opts.Location)
	if opts.Location == "" {
		opts.Location = defaultVertexBatchLocation
placeholder
	opts.ManagedGCSBucket = strings.Trim(strings.TrimSpace(opts.ManagedGCSBucket), "/")
	opts.ManagedGCSPrefix = strings.Trim(strings.TrimSpace(opts.ManagedGCSPrefix), "/")
	if opts.ManagedGCSPrefix == "" {
		opts.ManagedGCSPrefix = defaultVertexManagedGCSPrefix
placeholder
	opts.Environment = strings.TrimSpace(opts.Environment)
	if opts.Environment == "" {
		opts.Environment = "default"
placeholder
	opts.BatchPredictionBaseURL = strings.TrimRight(strings.TrimSpace(opts.BatchPredictionBaseURL), "/")
	opts.GCSBaseURL = strings.TrimRight(strings.TrimSpace(opts.GCSBaseURL), "/")
	return opts
placeholder

func (p *VertexBatchImageProvider) Name() string {
	return BatchImageProviderVertex
placeholder

func (p *VertexBatchImageProvider) SupportsAccount(account *Account) bool {
	if account == nil || account.Platform != PlatformGemini || account.Type != AccountTypeServiceAccount {
		return false
placeholder
	_, err := parseVertexServiceAccountKey(account)
	return err == nil
placeholder

func (p *VertexBatchImageProvider) Submit(ctx context.Context, job *BatchImageJob, account *Account, input BatchImageInput) (*BatchProviderJob, error) {
	if err := p.validateAccount(account); err != nil {
		return nil, err
placeholder
	if strings.TrimSpace(p.opts.ManagedGCSBucket) == "" {
		return nil, vertexProviderError("VERTEX_MANAGED_GCS_BUCKET_MISSING", "Vertex managed GCS bucket is not configured", nil)
placeholder
	if input.BatchID == "" && job != nil {
		input.BatchID = job.BatchID
placeholder
	if input.Model == "" && job != nil {
		input.Model = job.Model
placeholder

	jsonl, err := BuildVertexBatchJSONL(input)
	if err != nil {
		return nil, err
placeholder
	refs, err := p.managedRefs(input.BatchID)
	if err != nil {
		return nil, err
placeholder

	accessToken, err := p.accessToken(ctx, account)
	if err != nil {
		return nil, mapVertexClientError(err)
placeholder
	if err := p.objectStore.UploadJSONL(ctx, accessToken, refs.InputURI, bytes.NewReader(jsonl)); err != nil {
		return nil, vertexProviderError("VERTEX_GCS_UPLOAD_FAILED", "Vertex managed GCS upload failed", nil)
placeholder

	projectID := strings.TrimSpace(p.opts.ProjectID)
	if projectID == "" {
		projectID = account.VertexProjectID()
placeholder
	if projectID == "" {
		return nil, vertexProviderError("VERTEX_PROJECT_ID_MISSING", "Vertex project id is not configured", nil)
placeholder
	location := strings.TrimSpace(p.opts.Location)
	if location == "" {
		location = account.VertexLocation(input.Model)
placeholder

	req := VertexCreateBatchPredictionJobRequest{
		ProjectID:      projectID,
		Location:       location,
		DisplayName:    vertexBatchDisplayName(input),
		Model:          NormalizeVertexBatchModelPath(input.Model),
		InputConfig:    VertexBatchInputConfig{InstancesFormat: "jsonl", GCSSource: VertexBatchGCSSource{URIs: []string{refs.InputURIplaceholderplaceholderplaceholder,
		OutputConfig:   VertexBatchOutputConfig{PredictionsFormat: "jsonl", GCSDestination: VertexBatchGCSDestination{OutputURIPrefix: refs.OutputPrefixURIplaceholderplaceholder,
		InstanceConfig: &VertexBatchInstanceConfig{KeyField: "key"placeholder,
placeholder
	created, err := p.client.CreateBatchPredictionJob(ctx, accessToken, req)
	if err != nil {
		return nil, mapVertexClientError(err)
placeholder
	if created == nil || strings.TrimSpace(created.Name) == "" {
		return nil, vertexProviderError("VERTEX_INVALID_RESPONSE", "Vertex batch response is missing job name", nil)
placeholder
	return &BatchProviderJob{
		ProviderJobName:   created.Name,
		ProviderInputRef:  refs.InputURI,
		ProviderOutputRef: refs.OutputPrefixURI,
		RawState:          created.State,
placeholder, nil
placeholder

func (p *VertexBatchImageProvider) Get(ctx context.Context, job *BatchImageJob, account *Account) (*BatchProviderStatus, error) {
	if err := p.validateAccount(account); err != nil {
		return nil, err
placeholder
	jobName := batchImageProviderJobName(job)
	if jobName == "" {
		return nil, ErrBatchImageProviderMissingJobName
placeholder
	accessToken, err := p.accessToken(ctx, account)
	if err != nil {
		return nil, mapVertexClientError(err)
placeholder
	vertexJob, err := p.client.GetBatchPredictionJob(ctx, accessToken, jobName)
	if err != nil {
		return nil, mapVertexClientError(err)
placeholder
	if vertexJob == nil {
		return nil, vertexProviderError("VERTEX_INVALID_RESPONSE", "Vertex batch response is empty", nil)
placeholder
	status := mapVertexBatchState(vertexJob)
	outputRef := strings.TrimSpace(vertexJob.OutputConfig.GCSDestination.OutputURIPrefix)
	if outputRef == "" {
		outputRef = batchImageProviderOutputRef(job)
placeholder
	if outputRef == "" && job != nil && job.GCSOutputURI != nil {
		outputRef = strings.TrimSpace(*job.GCSOutputURI)
placeholder
	status.ProviderOutputRef = outputRef
	return status, nil
placeholder

func (p *VertexBatchImageProvider) Cancel(ctx context.Context, job *BatchImageJob, account *Account) error {
	if err := p.validateAccount(account); err != nil {
		return err
placeholder
	jobName := batchImageProviderJobName(job)
	if jobName == "" {
		return ErrBatchImageProviderMissingJobName
placeholder
	accessToken, err := p.accessToken(ctx, account)
	if err != nil {
		return mapVertexClientError(err)
placeholder
	return mapVertexClientError(p.client.CancelBatchPredictionJob(ctx, accessToken, jobName))
placeholder

func (p *VertexBatchImageProvider) OpenResult(ctx context.Context, job *BatchImageJob, account *Account) (io.ReadCloser, string, error) {
	if err := p.validateAccount(account); err != nil {
		return nil, "", err
placeholder
	outputRef := batchImageProviderOutputRef(job)
	if outputRef == "" && job != nil && job.GCSOutputURI != nil {
		outputRef = strings.TrimSpace(*job.GCSOutputURI)
placeholder
	if outputRef == "" {
		return nil, "", ErrBatchImageProviderMissingResultRef
placeholder
	accessToken, err := p.accessToken(ctx, account)
	if err != nil {
		return nil, "", mapVertexClientError(err)
placeholder
	objects, err := p.objectStore.ListJSONLObjects(ctx, accessToken, outputRef)
	if err != nil {
		return nil, "", vertexProviderError("VERTEX_GCS_LIST_FAILED", "Vertex managed GCS list failed", nil)
placeholder
	sort.Strings(objects)
	if len(objects) == 0 {
		return nil, "", vertexProviderError("VERTEX_RESULT_OBJECTS_MISSING", "Vertex result objects are missing", nil)
placeholder
	return &vertexCombinedJSONLReadCloser{
		ctx:         ctx,
		accessToken: accessToken,
		objects:     objects,
		store:       p.objectStore,
placeholder, "application/jsonl", nil
placeholder

func (p *VertexBatchImageProvider) Cleanup(ctx context.Context, job *BatchImageJob, account *Account, target CleanupTarget) error {
	if err := p.validateAccount(account); err != nil {
		return err
placeholder
	accessToken, err := p.accessToken(ctx, account)
	if err != nil {
		return mapVertexClientError(err)
placeholder
	inputRef := batchImageProviderInputRef(job)
	outputRef := batchImageProviderOutputRef(job)
	if job != nil {
		if inputRef == "" && job.GCSInputURI != nil {
			inputRef = strings.TrimSpace(*job.GCSInputURI)
	placeholder
		if outputRef == "" && job.GCSOutputURI != nil {
			outputRef = strings.TrimSpace(*job.GCSOutputURI)
	placeholder
placeholder

	switch target {
	case CleanupTargetInput:
		return p.deleteManagedInput(ctx, accessToken, job, inputRef)
	case CleanupTargetOutput:
		return p.deleteManagedOutput(ctx, accessToken, job, outputRef)
	case CleanupTargetAll:
		if err := p.deleteManagedInput(ctx, accessToken, job, inputRef); err != nil {
			return err
	placeholder
		return p.deleteManagedOutput(ctx, accessToken, job, outputRef)
	default:
		return ErrUnsupportedCleanupTarget
placeholder
placeholder

func (p *VertexBatchImageProvider) validateAccount(account *Account) error {
	if account == nil || account.Platform != PlatformGemini || account.Type != AccountTypeServiceAccount {
		return ErrBatchImageProviderUnsupportedAccount
placeholder
	if _, err := parseVertexServiceAccountKey(account); err != nil {
		return ErrBatchImageProviderMissingServiceAccount
placeholder
	return nil
placeholder

func (p *VertexBatchImageProvider) accessToken(ctx context.Context, account *Account) (string, error) {
	return getVertexServiceAccountAccessToken(ctx, p.tokenCache, account)
placeholder

func (p *VertexBatchImageProvider) deleteManagedInput(ctx context.Context, accessToken string, job *BatchImageJob, uri string) error {
	if strings.TrimSpace(uri) == "" {
		return nil
placeholder
	if !p.isSafeManagedInput(job, uri) {
		return ErrBatchImageProviderUnsafeCleanupPath
placeholder
	return mapVertexClientError(p.objectStore.DeleteObject(ctx, accessToken, uri))
placeholder

func (p *VertexBatchImageProvider) deleteManagedOutput(ctx context.Context, accessToken string, job *BatchImageJob, uri string) error {
	if strings.TrimSpace(uri) == "" {
		return nil
placeholder
	if !p.isSafeManagedOutput(job, uri) {
		return ErrBatchImageProviderUnsafeCleanupPath
placeholder
	return mapVertexClientError(p.objectStore.DeletePrefix(ctx, accessToken, uri))
placeholder

func (p *VertexBatchImageProvider) isSafeManagedInput(job *BatchImageJob, uri string) bool {
	if job == nil || strings.TrimSpace(job.BatchID) == "" {
		return false
placeholder
	refs, err := p.managedRefs(job.BatchID)
	return err == nil && strings.TrimSpace(uri) == refs.InputURI
placeholder

func (p *VertexBatchImageProvider) isSafeManagedOutput(job *BatchImageJob, uri string) bool {
	if job == nil || strings.TrimSpace(job.BatchID) == "" {
		return false
placeholder
	refs, err := p.managedRefs(job.BatchID)
	return err == nil && strings.HasPrefix(strings.TrimSpace(uri), refs.OutputPrefixURI)
placeholder

type vertexManagedRefs struct {
	Prefix          string
	InputURI        string
	OutputPrefixURI string
placeholder

func (p *VertexBatchImageProvider) managedRefs(batchID string) (vertexManagedRefs, error) {
	batchID = strings.TrimSpace(batchID)
	if !IsValidBatchImageID(batchID) {
		return vertexManagedRefs{placeholder, batchImageProviderInputError("valid batch_id is required")
placeholder
	bucket := strings.Trim(strings.TrimSpace(p.opts.ManagedGCSBucket), "/")
	if bucket == "" || strings.Contains(bucket, "://") {
		return vertexManagedRefs{placeholder, vertexProviderError("VERTEX_MANAGED_GCS_BUCKET_MISSING", "Vertex managed GCS bucket is not configured", nil)
placeholder
	prefix := buildVertexManagedGCSPrefix(p.opts.ManagedGCSPrefix, p.opts.Environment, batchID)
	if !strings.Contains(prefix, batchID) {
		return vertexManagedRefs{placeholder, batchImageProviderInputError("managed GCS prefix must contain batch_id")
placeholder
	base := "gs://" + bucket + "/" + strings.Trim(prefix, "/")
	return vertexManagedRefs{
		Prefix:          strings.Trim(prefix, "/"),
		InputURI:        base + "/input/requests.jsonl",
		OutputPrefixURI: base + "/output/",
placeholder, nil
placeholder

func buildVertexManagedGCSPrefix(template, env, batchID string) string {
	template = strings.Trim(strings.TrimSpace(template), "/")
	if template == "" {
		template = defaultVertexManagedGCSPrefix
placeholder
	env = sanitizeVertexGCSPathSegment(env)
	batchID = sanitizeVertexGCSPathSegment(batchID)
	prefix := strings.ReplaceAll(template, "{envplaceholder", env)
	prefix = strings.ReplaceAll(prefix, "{batch_idplaceholder", batchID)
	return strings.Trim(prefix, "/")
placeholder

func sanitizeVertexGCSPathSegment(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return "default"
placeholder
	var b strings.Builder
	for _, r := range v {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			_, _ = b.WriteRune(r)
		case r == '-', r == '_', r == '.':
			_, _ = b.WriteRune(r)
		default:
			_ = b.WriteByte('-')
	placeholder
placeholder
	return b.String()
placeholder

func vertexBatchDisplayName(input BatchImageInput) string {
	if v := strings.TrimSpace(input.DisplayName); v != "" {
		return v
placeholder
	if v := strings.TrimSpace(input.BatchID); v != "" {
		return "sub2api-" + v
placeholder
	return "sub2api-image-batch"
placeholder

func BuildVertexBatchJSONL(input BatchImageInput) ([]byte, error) {
	if strings.TrimSpace(input.Model) == "" {
		return nil, batchImageProviderInputError("model is required")
placeholder
	if len(input.Items) == 0 {
		return nil, batchImageProviderInputError("at least one item is required")
placeholder
	seen := make(map[string]struct{placeholder, len(input.Items))
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, item := range input.Items {
		customID := strings.TrimSpace(item.CustomID)
		if customID == "" {
			return nil, batchImageProviderInputError("custom_id is required")
	placeholder
		if _, ok := seen[customID]; ok {
			return nil, batchImageProviderInputError("duplicate custom_id %q", customID)
	placeholder
		seen[customID] = struct{placeholder{placeholder
		prompt := strings.TrimSpace(item.Prompt)
		if prompt == "" {
			return nil, batchImageProviderInputError("prompt is required for custom_id %q", customID)
	placeholder
		parts, err := vertexBatchImageParts(prompt, item.ReferenceImages)
		if err != nil {
			return nil, err
	placeholder
		line := map[string]any{
			"key": customID,
			"request": map[string]any{
				"contents": []any{map[string]any{
					"role":  "user",
					"parts": parts,
		placeholder
				"generationConfig": map[string]any{
					"responseModalities": []string{"TEXT", "IMAGE"placeholder,
			placeholder,
		placeholder,
	placeholder
		if err := enc.Encode(line); err != nil {
			return nil, err
	placeholder
placeholder
	return buf.Bytes(), nil
placeholder

func vertexBatchImageParts(prompt string, refs []BatchImageReference) ([]any, error) {
	parts := []any{map[string]any{"text": promptplaceholderplaceholder
	for _, ref := range refs {
		mimeType := normalizeBatchImageReferenceMimeType(ref.MimeType)
		if mimeType == "" {
			return nil, batchImageProviderInputError("reference image mime_type is required")
	placeholder
		fileURI := strings.TrimSpace(ref.FileURI)
		switch {
		case len(ref.Data) > 0 && fileURI == "":
			parts = append(parts, map[string]any{
				"inlineData": map[string]any{
					"mimeType": mimeType,
					"data":     base64.StdEncoding.EncodeToString(ref.Data),
			placeholder,
		placeholder)
		case len(ref.Data) == 0 && fileURI != "":
			parts = append(parts, map[string]any{
				"fileData": map[string]any{
					"mimeType": mimeType,
					"fileUri":  fileURI,
			placeholder,
		placeholder)
		default:
			return nil, batchImageProviderInputError("reference image must contain exactly one of data or file_uri")
	placeholder
placeholder
	return parts, nil
placeholder

func NormalizeVertexBatchModelPath(model string) string {
	model = strings.Trim(strings.TrimSpace(model), "/")
	if strings.HasPrefix(model, "publishers/") || strings.HasPrefix(model, "projects/") {
		return model
placeholder
	return "publishers/google/models/" + model
placeholder

func BuildVertexBatchPredictionJobsEndpoint(baseURL, projectID, location string) (string, error) {
	projectID = strings.TrimSpace(projectID)
	location = strings.TrimSpace(location)
	if projectID == "" {
		return "", errors.New("vertex project_id is required")
placeholder
	if location == "" {
		location = defaultVertexBatchLocation
placeholder
	if !vertexLocationPattern.MatchString(location) {
		return "", fmt.Errorf("invalid vertex location: %s", location)
placeholder
	if strings.TrimSpace(baseURL) != "" {
		return strings.TrimRight(strings.TrimSpace(baseURL), "/") + "/v1/projects/" + url.PathEscape(projectID) + "/locations/" + url.PathEscape(location) + "/batchPredictionJobs", nil
placeholder
	host := fmt.Sprintf("%s-aiplatform.googleapis.com", location)
	if location == "global" {
		host = "aiplatform.googleapis.com"
placeholder
	return fmt.Sprintf("https://%s/v1/projects/%s/locations/%s/batchPredictionJobs", host, url.PathEscape(projectID), url.PathEscape(location)), nil
placeholder

func mapVertexBatchState(job *VertexBatchPredictionJob) *BatchProviderStatus {
	state := strings.TrimSpace(job.State)
	status := &BatchProviderStatus{
		RawState:              state,
		InternalState:         BatchProviderStateRunning,
		SuggestedRequeueAfter: defaultVertexBatchRequeueAfter,
placeholder
	switch strings.ToUpper(state) {
	case "JOB_STATE_PENDING", "JOB_STATE_QUEUED":
		status.InternalState = BatchProviderStateQueued
	case "JOB_STATE_RUNNING", "JOB_STATE_PAUSED":
		status.InternalState = BatchProviderStateRunning
	case "JOB_STATE_SUCCEEDED":
		status.InternalState = BatchProviderStateSucceeded
		status.Done = true
		status.SuggestedRequeueAfter = 0
	case "JOB_STATE_FAILED":
		status.InternalState = BatchProviderStateFailed
		status.Done = true
		status.ErrorCode = "VERTEX_BATCH_FAILED"
		status.SuggestedRequeueAfter = 0
	case "JOB_STATE_CANCELLED":
		status.InternalState = BatchProviderStateCancelled
		status.Done = true
		status.ErrorCode = "VERTEX_BATCH_CANCELLED"
		status.SuggestedRequeueAfter = 0
	case "JOB_STATE_EXPIRED":
		status.InternalState = BatchProviderStateExpired
		status.Done = true
		status.ErrorCode = "VERTEX_BATCH_EXPIRED"
		status.SuggestedRequeueAfter = 0
	default:
		if job.Error != nil && strings.TrimSpace(job.Error.Message) != "" {
			status.InternalState = BatchProviderStateFailed
			status.Done = true
			status.ErrorCode = "VERTEX_BATCH_FAILED"
			status.SuggestedRequeueAfter = 0
	placeholder
placeholder
	if job.Error != nil {
		if code := strings.TrimSpace(job.Error.Status); code != "" {
			status.ErrorCode = code
	placeholder
		status.ErrorMessage = strings.TrimSpace(job.Error.Message)
placeholder
	return status
placeholder

func vertexProviderError(reason, message string, cause error) error {
	err := infraerrors.New(http.StatusBadGateway, reason, message)
	if cause != nil {
		return err.WithCause(cause)
placeholder
	return err
placeholder

func mapVertexClientError(err error) error {
	if err == nil {
		return nil
placeholder
	if errors.Is(err, ErrBatchImageProviderMissingServiceAccount) ||
		errors.Is(err, ErrBatchImageProviderMissingJobName) ||
		errors.Is(err, ErrBatchImageProviderMissingResultRef) ||
		errors.Is(err, ErrBatchImageProviderUnsafeCleanupPath) ||
		errors.Is(err, ErrUnsupportedCleanupTarget) {
		return err
placeholder
	var apiErr *VertexAPIError
	if errors.As(err, &apiErr) {
		switch apiErr.StatusCode {
		case http.StatusUnauthorized:
			return vertexProviderError("VERTEX_AUTH_FAILED", "Vertex authentication failed", nil)
		case http.StatusForbidden:
			return vertexProviderError("VERTEX_PERMISSION_DENIED", "Vertex permission denied", nil)
		case http.StatusTooManyRequests:
			return vertexProviderError("VERTEX_RATE_LIMITED", "Vertex rate limit exceeded", nil)
		case http.StatusNotFound:
			return vertexProviderError("VERTEX_BATCH_NOT_FOUND", "Vertex batch resource was not found", nil)
		default:
			return vertexProviderError("VERTEX_INVALID_RESPONSE", "Vertex API request failed", nil)
	placeholder
placeholder
	return vertexProviderError("VERTEX_INVALID_RESPONSE", "Vertex API request failed", err)
placeholder

type vertexCombinedJSONLReadCloser struct {
	ctx          context.Context
	accessToken  string
	objects      []string
	store        VertexBatchObjectStore
	index        int
	current      io.ReadCloser
	needBoundary bool
	closed       bool
placeholder

func (r *vertexCombinedJSONLReadCloser) Read(p []byte) (int, error) {
	if r.closed {
		return 0, io.ErrClosedPipe
placeholder
	if r.needBoundary {
		if len(p) == 0 {
			return 0, nil
	placeholder
		p[0] = '\n'
		r.needBoundary = false
		return 1, nil
placeholder
	for {
		if r.current == nil {
			if r.index >= len(r.objects) {
				return 0, io.EOF
		placeholder
			obj := r.objects[r.index]
			r.index++
			rc, _, err := r.store.OpenObject(r.ctx, r.accessToken, obj)
			if err != nil {
				return 0, err
		placeholder
			r.current = rc
	placeholder
		n, err := r.current.Read(p)
		if err == io.EOF {
			_ = r.current.Close()
			r.current = nil
			if r.index < len(r.objects) {
				if n > 0 {
					r.needBoundary = true
					return n, nil
			placeholder
				if len(p) == 0 {
					return 0, nil
			placeholder
				p[0] = '\n'
				return 1, nil
		placeholder
			if n > 0 {
				return n, nil
		placeholder
			continue
	placeholder
		return n, err
placeholder
placeholder

func (r *vertexCombinedJSONLReadCloser) Close() error {
	r.closed = true
	if r.current != nil {
		return r.current.Close()
placeholder
	return nil
placeholder

type VertexBatchHTTPClient struct {
	baseURL string
	client  *http.Client
placeholder

func NewVertexBatchHTTPClient(baseURL string, client *http.Client) *VertexBatchHTTPClient {
	if client == nil {
		client = batchImageDefaultHTTPClient()
placeholder
	return &VertexBatchHTTPClient{baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"), client: clientplaceholder
placeholder

func (c *VertexBatchHTTPClient) CreateBatchPredictionJob(ctx context.Context, accessToken string, req VertexCreateBatchPredictionJobRequest) (*VertexBatchPredictionJob, error) {
	endpoint, err := BuildVertexBatchPredictionJobsEndpoint(c.baseURL, req.ProjectID, req.Location)
	if err != nil {
		return nil, err
placeholder
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
placeholder
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, err
placeholder
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+accessToken)
	return doVertexJSON[VertexBatchPredictionJob](c.client, httpReq)
placeholder

func (c *VertexBatchHTTPClient) GetBatchPredictionJob(ctx context.Context, accessToken string, name string) (*VertexBatchPredictionJob, error) {
	endpoint := c.vertexResourceURL(name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
placeholder
	req.Header.Set("Authorization", "Bearer "+accessToken)
	return doVertexJSON[VertexBatchPredictionJob](c.client, req)
placeholder

func (c *VertexBatchHTTPClient) CancelBatchPredictionJob(ctx context.Context, accessToken string, name string) error {
	endpoint := c.vertexResourceURL(name) + ":cancel"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return err
placeholder
	req.Header.Set("Authorization", "Bearer "+accessToken)
	return doVertexNoBody(c.client, req)
placeholder

func (c *VertexBatchHTTPClient) vertexResourceURL(name string) string {
	name = strings.TrimLeft(strings.TrimSpace(name), "/")
	if c.baseURL != "" {
		return c.baseURL + "/v1/" + name
placeholder
	return "https://aiplatform.googleapis.com/v1/" + name
placeholder

type VertexGCSObjectStore struct {
	baseURL string
	client  *http.Client
placeholder

func NewVertexGCSObjectStore(baseURL string, client *http.Client) *VertexGCSObjectStore {
	if client == nil {
		client = batchImageDefaultHTTPClient()
placeholder
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = "https://storage.googleapis.com"
placeholder
	return &VertexGCSObjectStore{baseURL: baseURL, client: clientplaceholder
placeholder

func (s *VertexGCSObjectStore) UploadJSONL(ctx context.Context, accessToken string, uri string, r io.Reader) error {
	bucket, object, err := parseGCSURI(uri)
	if err != nil {
		return err
placeholder
	endpoint := fmt.Sprintf("%s/upload/storage/v1/b/%s/o?uploadType=media&name=%s", s.baseURL, url.PathEscape(bucket), url.QueryEscape(object))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, r)
	if err != nil {
		return err
placeholder
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/jsonl")
	return doVertexNoBody(s.client, req)
placeholder

func (s *VertexGCSObjectStore) ListJSONLObjects(ctx context.Context, accessToken string, prefixURI string) ([]string, error) {
	return s.listObjects(ctx, accessToken, prefixURI, true)
placeholder

func (s *VertexGCSObjectStore) listObjects(ctx context.Context, accessToken string, prefixURI string, jsonlOnly bool) ([]string, error) {
	bucket, prefix, err := parseGCSURI(prefixURI)
	if err != nil {
		return nil, err
placeholder
	var objects []string
	pageToken := ""
	for {
		endpoint := fmt.Sprintf("%s/storage/v1/b/%s/o?prefix=%s", s.baseURL, url.PathEscape(bucket), url.QueryEscape(prefix))
		if pageToken != "" {
			endpoint += "&pageToken=" + url.QueryEscape(pageToken)
	placeholder
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, err
	placeholder
		req.Header.Set("Authorization", "Bearer "+accessToken)
		var page struct {
			Items []struct {
				Name string `json:"name"`
		placeholder `json:"items"`
			NextPageToken string `json:"nextPageToken"`
	placeholder
		if err := doVertexDecodeJSON(s.client, req, &page); err != nil {
			return nil, err
	placeholder
		for _, item := range page.Items {
			if !jsonlOnly || strings.HasSuffix(item.Name, ".jsonl") {
				objects = append(objects, "gs://"+bucket+"/"+item.Name)
		placeholder
	placeholder
		if page.NextPageToken == "" {
			return objects, nil
	placeholder
		pageToken = page.NextPageToken
placeholder
placeholder

func (s *VertexGCSObjectStore) OpenObject(ctx context.Context, accessToken string, uri string) (io.ReadCloser, string, error) {
	bucket, object, err := parseGCSURI(uri)
	if err != nil {
		return nil, "", err
placeholder
	endpoint := fmt.Sprintf("%s/storage/v1/b/%s/o/%s?alt=media", s.baseURL, url.PathEscape(bucket), url.PathEscape(object))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, "", err
placeholder
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, "", err
placeholder
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer func() { _ = resp.Body.Close() placeholder()
		return nil, "", readVertexAPIError(resp)
placeholder
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/jsonl"
placeholder
	return resp.Body, contentType, nil
placeholder

func (s *VertexGCSObjectStore) DeleteObject(ctx context.Context, accessToken string, uri string) error {
	bucket, object, err := parseGCSURI(uri)
	if err != nil {
		return err
placeholder
	endpoint := fmt.Sprintf("%s/storage/v1/b/%s/o/%s", s.baseURL, url.PathEscape(bucket), url.PathEscape(object))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
placeholder
	req.Header.Set("Authorization", "Bearer "+accessToken)
	return doVertexNoBody(s.client, req)
placeholder

func (s *VertexGCSObjectStore) DeletePrefix(ctx context.Context, accessToken string, prefixURI string) error {
	objects, err := s.listObjects(ctx, accessToken, prefixURI, false)
	if err != nil {
		return err
placeholder
	for _, object := range objects {
		if err := s.DeleteObject(ctx, accessToken, object); err != nil {
			return err
	placeholder
placeholder
	return nil
placeholder

func parseGCSURI(uri string) (bucket, object string, err error) {
	uri = strings.TrimSpace(uri)
	if !strings.HasPrefix(uri, "gs://") {
		return "", "", fmt.Errorf("invalid gcs uri")
placeholder
	rest := strings.TrimPrefix(uri, "gs://")
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
		return "", "", fmt.Errorf("invalid gcs uri")
placeholder
	return parts[0], parts[1], nil
placeholder

type VertexAPIError struct {
	StatusCode int
	Code       string
	Message    string
placeholder

func (e *VertexAPIError) Error() string {
	if e == nil {
		return "<nil>"
placeholder
	if e.Code != "" {
		return fmt.Sprintf("vertex api error: status=%d code=%s message=%s", e.StatusCode, e.Code, e.Message)
placeholder
	return fmt.Sprintf("vertex api error: status=%d message=%s", e.StatusCode, e.Message)
placeholder

func doVertexJSON[T any](client *http.Client, req *http.Request) (*T, error) {
	var out T
	if err := doVertexDecodeJSON(client, req, &out); err != nil {
		return nil, err
placeholder
	return &out, nil
placeholder

func doVertexDecodeJSON(client *http.Client, req *http.Request, out any) error {
	resp, err := client.Do(req)
	if err != nil {
		return err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return readVertexAPIError(resp)
placeholder
	return json.NewDecoder(resp.Body).Decode(out)
placeholder

func doVertexNoBody(client *http.Client, req *http.Request) error {
	resp, err := client.Do(req)
	if err != nil {
		return err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return readVertexAPIError(resp)
placeholder
	return nil
placeholder

func readVertexAPIError(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
	message := string(body)
	code := ""
	var parsed struct {
		Error struct {
			Code    any    `json:"code"`
			Message string `json:"message"`
			Status  string `json:"status"`
	placeholder `json:"error"`
placeholder
	if err := json.Unmarshal(body, &parsed); err == nil && parsed.Error.Message != "" {
		message = parsed.Error.Message
		code = parsed.Error.Status
placeholder
	return &VertexAPIError{StatusCode: resp.StatusCode, Code: code, Message: messageplaceholder
placeholder

var _ BatchImageProvider = (*VertexBatchImageProvider)(nil)
var _ VertexBatchClient = (*VertexBatchHTTPClient)(nil)
var _ VertexBatchObjectStore = (*VertexGCSObjectStore)(nil)
