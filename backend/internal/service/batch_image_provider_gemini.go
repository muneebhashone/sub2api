package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/geminicli"
)

const defaultGeminiBatchRequeueAfter = 30 * time.Second

type GeminiBatchClient interface {
	UploadJSONL(ctx context.Context, apiKey string, displayName string, r io.Reader) (*GeminiUploadedFile, error)
	CreateBatch(ctx context.Context, apiKey string, model string, fileName string, displayName string) (*GeminiBatchJob, error)
	GetBatch(ctx context.Context, apiKey string, batchName string) (*GeminiBatchJob, error)
	CancelBatch(ctx context.Context, apiKey string, batchName string) error
	DownloadFile(ctx context.Context, apiKey string, fileName string) (io.ReadCloser, string, error)
	DeleteFile(ctx context.Context, apiKey string, fileName string) error
placeholder

type GeminiUploadedFile struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	URI         string `json:"uri"`
	MimeType    string `json:"mimeType"`
placeholder

type GeminiBatchJob struct {
	Name     string               `json:"name"`
	State    string               `json:"state"`
	Dest     *GeminiBatchDest     `json:"dest"`
	Response *GeminiBatchResponse `json:"response"`
	Error    *GeminiBatchError    `json:"error"`
	Raw      map[string]any       `json:"-"`
placeholder

type GeminiBatchDest struct {
	FileName      string `json:"fileName"`
	FileNameSnake string `json:"file_name"`
placeholder

type GeminiBatchResponse struct {
	ResponsesFile       string `json:"responsesFile"`
	ResponsesFileSnake  string `json:"responses_file"`
	InlinedResponses    []any  `json:"inlinedResponses"`
	InlinedResponsesAlt []any  `json:"inlined_responses"`
placeholder

type GeminiBatchError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
placeholder

type GeminiAPIBatchImageProvider struct {
	client GeminiBatchClient
placeholder

func NewGeminiAPIBatchImageProvider(client GeminiBatchClient) *GeminiAPIBatchImageProvider {
	if client == nil {
		client = NewGeminiBatchHTTPClient("", nil)
placeholder
	return &GeminiAPIBatchImageProvider{client: clientplaceholder
placeholder

func (p *GeminiAPIBatchImageProvider) Name() string {
	return BatchImageProviderGeminiAPI
placeholder

func (p *GeminiAPIBatchImageProvider) SupportsAccount(account *Account) bool {
	return account != nil &&
		account.Platform == PlatformGemini &&
		account.Type == AccountTypeAPIKey &&
		batchImageProviderAPIKey(account) != ""
placeholder

func (p *GeminiAPIBatchImageProvider) Submit(ctx context.Context, job *BatchImageJob, account *Account, input BatchImageInput) (*BatchProviderJob, error) {
	if account == nil || account.Platform != PlatformGemini || account.Type != AccountTypeAPIKey {
		return nil, ErrBatchImageProviderUnsupportedAccount
placeholder
	apiKey := batchImageProviderAPIKey(account)
	if apiKey == "" {
		return nil, ErrBatchImageProviderMissingAPIKey
placeholder
	if input.BatchID == "" && job != nil {
		input.BatchID = job.BatchID
placeholder
	if input.Model == "" && job != nil {
		input.Model = job.Model
placeholder

	jsonl, err := BuildGeminiBatchJSONL(input)
	if err != nil {
		return nil, err
placeholder

	displayName := strings.TrimSpace(input.DisplayName)
	if displayName == "" {
		displayName = strings.TrimSpace(input.BatchID)
placeholder

	uploaded, err := p.client.UploadJSONL(ctx, apiKey, displayName, bytes.NewReader(jsonl))
	if err != nil {
		return nil, mapGeminiClientError(err)
placeholder
	if uploaded == nil || strings.TrimSpace(uploaded.Name) == "" {
		return nil, geminiProviderError("GEMINI_INVALID_RESPONSE", "Gemini upload response is missing file name", nil)
placeholder

	batch, err := p.client.CreateBatch(ctx, apiKey, input.Model, uploaded.Name, displayName)
	if err != nil {
		return nil, mapGeminiClientError(err)
placeholder
	if batch == nil || strings.TrimSpace(batch.Name) == "" {
		return nil, geminiProviderError("GEMINI_INVALID_RESPONSE", "Gemini batch response is missing job name", nil)
placeholder

	return &BatchProviderJob{
		ProviderJobName:  batch.Name,
		ProviderInputRef: uploaded.Name,
		RawState:         batch.State,
placeholder, nil
placeholder

func (p *GeminiAPIBatchImageProvider) Get(ctx context.Context, job *BatchImageJob, account *Account) (*BatchProviderStatus, error) {
	if account == nil || account.Platform != PlatformGemini || account.Type != AccountTypeAPIKey {
		return nil, ErrBatchImageProviderUnsupportedAccount
placeholder
	apiKey := batchImageProviderAPIKey(account)
	if apiKey == "" {
		return nil, ErrBatchImageProviderMissingAPIKey
placeholder
	jobName := batchImageProviderJobName(job)
	if jobName == "" {
		return nil, ErrBatchImageProviderMissingJobName
placeholder

	batch, err := p.client.GetBatch(ctx, apiKey, jobName)
	if err != nil {
		return nil, mapGeminiClientError(err)
placeholder
	if batch == nil {
		return nil, geminiProviderError("GEMINI_INVALID_RESPONSE", "Gemini batch response is empty", nil)
placeholder

	status := mapGeminiBatchState(batch)
	if status.InternalState == BatchProviderStateSucceeded {
		if geminiBatchHasInlineResults(batch) {
			return nil, ErrBatchImageProviderInlineResultUnsupported
	placeholder
		outputRef := geminiBatchOutputRef(batch)
		if outputRef == "" {
			status.InternalState = BatchProviderStateFailed
			status.Done = true
			status.ErrorCode = "GEMINI_RESULT_FILE_MISSING"
			status.ErrorMessage = "Gemini batch succeeded without a result file reference"
	placeholder
		status.ProviderOutputRef = outputRef
placeholder
	return status, nil
placeholder

func (p *GeminiAPIBatchImageProvider) Cancel(ctx context.Context, job *BatchImageJob, account *Account) error {
	if account == nil || account.Platform != PlatformGemini || account.Type != AccountTypeAPIKey {
		return ErrBatchImageProviderUnsupportedAccount
placeholder
	apiKey := batchImageProviderAPIKey(account)
	if apiKey == "" {
		return ErrBatchImageProviderMissingAPIKey
placeholder
	jobName := batchImageProviderJobName(job)
	if jobName == "" {
		return ErrBatchImageProviderMissingJobName
placeholder
	return mapGeminiClientError(p.client.CancelBatch(ctx, apiKey, jobName))
placeholder

func (p *GeminiAPIBatchImageProvider) OpenResult(ctx context.Context, job *BatchImageJob, account *Account) (io.ReadCloser, string, error) {
	if account == nil || account.Platform != PlatformGemini || account.Type != AccountTypeAPIKey {
		return nil, "", ErrBatchImageProviderUnsupportedAccount
placeholder
	apiKey := batchImageProviderAPIKey(account)
	if apiKey == "" {
		return nil, "", ErrBatchImageProviderMissingAPIKey
placeholder
	outputRef := batchImageProviderOutputRef(job)
	if outputRef == "" {
		return nil, "", ErrBatchImageProviderMissingResultRef
placeholder
	r, contentType, err := p.client.DownloadFile(ctx, apiKey, outputRef)
	return r, contentType, mapGeminiClientError(err)
placeholder

func (p *GeminiAPIBatchImageProvider) Cleanup(ctx context.Context, job *BatchImageJob, account *Account, target CleanupTarget) error {
	if account == nil || account.Platform != PlatformGemini || account.Type != AccountTypeAPIKey {
		return ErrBatchImageProviderUnsupportedAccount
placeholder
	apiKey := batchImageProviderAPIKey(account)
	if apiKey == "" {
		return ErrBatchImageProviderMissingAPIKey
placeholder

	switch target {
	case CleanupTargetInput:
		return p.deleteGeminiFileIfPresent(ctx, apiKey, batchImageProviderInputRef(job))
	case CleanupTargetOutput:
		return p.deleteGeminiFileIfPresent(ctx, apiKey, batchImageProviderOutputRef(job))
	case CleanupTargetAll:
		if err := p.deleteGeminiFileIfPresent(ctx, apiKey, batchImageProviderInputRef(job)); err != nil {
			return err
	placeholder
		return p.deleteGeminiFileIfPresent(ctx, apiKey, batchImageProviderOutputRef(job))
	default:
		return ErrUnsupportedCleanupTarget
placeholder
placeholder

func (p *GeminiAPIBatchImageProvider) deleteGeminiFileIfPresent(ctx context.Context, apiKey, fileName string) error {
	if strings.TrimSpace(fileName) == "" {
		return nil
placeholder
	return mapGeminiClientError(p.client.DeleteFile(ctx, apiKey, fileName))
placeholder

type geminiJSONLLine struct {
	Key     string                `json:"key"`
	Request geminiGenerateRequest `json:"request"`
placeholder

type geminiGenerateRequest struct {
	Contents         []geminiContent        `json:"contents"`
	GenerationConfig geminiGenerationConfig `json:"generationConfig"`
placeholder

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
placeholder

type geminiPart struct {
	Text string `json:"text,omitempty"`
placeholder

type geminiGenerationConfig struct {
	ResponseModalities []string `json:"responseModalities"`
placeholder

func BuildGeminiBatchJSONL(input BatchImageInput) ([]byte, error) {
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
		if len(item.ReferenceImages) > 0 {
			return nil, batchImageProviderInputError("reference images are not supported in PR3")
	placeholder

		// TODO(batch-image): add response_mime_type/aspect_ratio/image_size once the
		// Gemini batch image REST shape is stabilized for those options.
		line := geminiJSONLLine{
			Key: customID,
			Request: geminiGenerateRequest{
				Contents: []geminiContent{{
					Parts: []geminiPart{{Text: promptplaceholderplaceholder,
		placeholder
				GenerationConfig: geminiGenerationConfig{
					ResponseModalities: []string{"TEXT", "IMAGE"placeholder,
			placeholder,
		placeholder,
	placeholder
		if err := enc.Encode(line); err != nil {
			return nil, err
	placeholder
placeholder
	return buf.Bytes(), nil
placeholder

func mapGeminiBatchState(batch *GeminiBatchJob) *BatchProviderStatus {
	state := strings.TrimSpace(batch.State)
	normalized := strings.ToUpper(state)
	status := &BatchProviderStatus{
		RawState:              state,
		InternalState:         BatchProviderStateRunning,
		SuggestedRequeueAfter: defaultGeminiBatchRequeueAfter,
placeholder

	switch normalized {
	case "JOB_STATE_PENDING", "JOB_STATE_QUEUED":
		status.InternalState = BatchProviderStateQueued
	case "JOB_STATE_RUNNING":
		status.InternalState = BatchProviderStateRunning
	case "JOB_STATE_SUCCEEDED":
		status.InternalState = BatchProviderStateSucceeded
		status.Done = true
	case "JOB_STATE_FAILED":
		status.InternalState = BatchProviderStateFailed
		status.Done = true
		status.ErrorCode = "GEMINI_BATCH_FAILED"
	case "JOB_STATE_CANCELLED":
		status.InternalState = BatchProviderStateCancelled
		status.Done = true
		status.ErrorCode = "GEMINI_BATCH_CANCELLED"
	case "JOB_STATE_EXPIRED":
		status.InternalState = BatchProviderStateExpired
		status.Done = true
		status.ErrorCode = "GEMINI_BATCH_EXPIRED"
	default:
		if batch.Error != nil && (strings.TrimSpace(batch.Error.Message) != "" || strings.TrimSpace(batch.Error.Code) != "") {
			status.InternalState = BatchProviderStateFailed
			status.Done = true
			status.ErrorCode = "GEMINI_BATCH_FAILED"
	placeholder
placeholder

	if batch.Error != nil {
		if code := strings.TrimSpace(batch.Error.Code); code != "" {
			status.ErrorCode = code
	placeholder else if status.ErrorCode == "" && strings.TrimSpace(batch.Error.Status) != "" {
			status.ErrorCode = strings.TrimSpace(batch.Error.Status)
	placeholder
		status.ErrorMessage = strings.TrimSpace(batch.Error.Message)
placeholder
	return status
placeholder

func geminiBatchOutputRef(batch *GeminiBatchJob) string {
	if batch == nil {
		return ""
placeholder
	if batch.Dest != nil {
		if v := strings.TrimSpace(batch.Dest.FileName); v != "" {
			return v
	placeholder
		if v := strings.TrimSpace(batch.Dest.FileNameSnake); v != "" {
			return v
	placeholder
placeholder
	if batch.Response != nil {
		if v := strings.TrimSpace(batch.Response.ResponsesFile); v != "" {
			return v
	placeholder
		if v := strings.TrimSpace(batch.Response.ResponsesFileSnake); v != "" {
			return v
	placeholder
placeholder
	return ""
placeholder

func geminiBatchHasInlineResults(batch *GeminiBatchJob) bool {
	return batch != nil &&
		batch.Response != nil &&
		(len(batch.Response.InlinedResponses) > 0 || len(batch.Response.InlinedResponsesAlt) > 0)
placeholder

func geminiProviderError(reason, message string, cause error) error {
	err := infraerrors.New(http.StatusBadGateway, reason, message)
	if cause != nil {
		return err.WithCause(cause)
placeholder
	return err
placeholder

func mapGeminiClientError(err error) error {
	if err == nil {
		return nil
placeholder
	var apiErr *GeminiAPIError
	if errors.As(err, &apiErr) {
		switch apiErr.StatusCode {
		case http.StatusUnauthorized, http.StatusForbidden:
			return geminiProviderError("GEMINI_AUTH_FAILED", "Gemini authentication failed", nil)
		case http.StatusTooManyRequests:
			return geminiProviderError("GEMINI_RATE_LIMITED", "Gemini rate limit exceeded", nil)
		case http.StatusNotFound:
			return geminiProviderError("GEMINI_BATCH_NOT_FOUND", "Gemini batch resource was not found", nil)
		default:
			return geminiProviderError("GEMINI_INVALID_RESPONSE", "Gemini API request failed", nil)
	placeholder
placeholder
	return geminiProviderError("GEMINI_INVALID_RESPONSE", "Gemini API request failed", nil)
placeholder

type GeminiBatchHTTPClient struct {
	baseURL string
	client  *http.Client
placeholder

func NewGeminiBatchHTTPClient(baseURL string, client *http.Client) *GeminiBatchHTTPClient {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = geminicli.AIStudioBaseURL
placeholder
	if client == nil {
		client = http.DefaultClient
placeholder
	return &GeminiBatchHTTPClient{baseURL: baseURL, client: clientplaceholder
placeholder

func (c *GeminiBatchHTTPClient) UploadJSONL(ctx context.Context, apiKey string, displayName string, r io.Reader) (*GeminiUploadedFile, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	metadataHeader := textproto.MIMEHeader{placeholder
	metadataHeader.Set("Content-Disposition", `form-data; name="metadata"`)
	metadataHeader.Set("Content-Type", "application/json; charset=utf-8")
	metadataPart, err := writer.CreatePart(metadataHeader)
	if err != nil {
		return nil, err
placeholder
	metadata := map[string]any{"file": map[string]any{"displayName": displayName, "mimeType": "application/jsonl"placeholderplaceholder
	if err := json.NewEncoder(metadataPart).Encode(metadata); err != nil {
		return nil, err
placeholder
	fileHeader := textproto.MIMEHeader{placeholder
	fileHeader.Set("Content-Disposition", `form-data; name="file"; filename="batch.jsonl"`)
	fileHeader.Set("Content-Type", "application/jsonl")
	filePart, err := writer.CreatePart(fileHeader)
	if err != nil {
		return nil, err
placeholder
	if _, err := io.Copy(filePart, r); err != nil {
		return nil, err
placeholder
	if err := writer.Close(); err != nil {
		return nil, err
placeholder

	req, err := c.newRequest(ctx, http.MethodPost, "/upload/v1beta/files?uploadType=multipart", apiKey, &body)
	if err != nil {
		return nil, err
placeholder
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var resp struct {
		File *GeminiUploadedFile `json:"file"`
		*GeminiUploadedFile
placeholder
	if err := c.doJSON(req, &resp); err != nil {
		return nil, err
placeholder
	if resp.File != nil {
		return resp.File, nil
placeholder
	return resp.GeminiUploadedFile, nil
placeholder

func (c *GeminiBatchHTTPClient) CreateBatch(ctx context.Context, apiKey string, model string, fileName string, displayName string) (*GeminiBatchJob, error) {
	body := map[string]any{
		"batch": map[string]any{
			"displayName": displayName,
			"inputConfig": map[string]any{
				"fileName": fileName,
		placeholder,
	placeholder,
placeholder
	payload, _ := json.Marshal(body)
	path := fmt.Sprintf("/v1beta/models/%s:batchGenerateContent", url.PathEscape(strings.TrimSpace(model)))
	req, err := c.newRequest(ctx, http.MethodPost, path, apiKey, bytes.NewReader(payload))
	if err != nil {
		return nil, err
placeholder
	req.Header.Set("Content-Type", "application/json")
	return c.doBatchJob(req)
placeholder

func (c *GeminiBatchHTTPClient) GetBatch(ctx context.Context, apiKey string, batchName string) (*GeminiBatchJob, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v1beta/"+strings.TrimLeft(batchName, "/"), apiKey, nil)
	if err != nil {
		return nil, err
placeholder
	return c.doBatchJob(req)
placeholder

func (c *GeminiBatchHTTPClient) CancelBatch(ctx context.Context, apiKey string, batchName string) error {
	req, err := c.newRequest(ctx, http.MethodPost, "/v1beta/"+strings.TrimLeft(batchName, "/")+":cancel", apiKey, nil)
	if err != nil {
		return err
placeholder
	return c.doNoBody(req)
placeholder

func (c *GeminiBatchHTTPClient) DownloadFile(ctx context.Context, apiKey string, fileName string) (io.ReadCloser, string, error) {
	metaReq, err := c.newRequest(ctx, http.MethodGet, "/v1beta/"+strings.TrimLeft(fileName, "/"), apiKey, nil)
	if err != nil {
		return nil, "", err
placeholder
	var metadata struct {
		DownloadURI string `json:"downloadUri"`
		DownloadURL string `json:"download_url"`
		MimeType    string `json:"mimeType"`
placeholder
	if err := c.doJSON(metaReq, &metadata); err != nil {
		return nil, "", err
placeholder
	downloadURL := strings.TrimSpace(metadata.DownloadURI)
	if downloadURL == "" {
		downloadURL = strings.TrimSpace(metadata.DownloadURL)
placeholder
	if downloadURL == "" {
		downloadURL = c.baseURL + "/v1beta/" + strings.TrimLeft(fileName, "/") + ":download"
placeholder
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return nil, "", err
placeholder
	req.Header.Set("x-goog-api-key", apiKey)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, "", err
placeholder
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		return nil, "", readGeminiAPIError(resp)
placeholder
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = metadata.MimeType
placeholder
	if contentType == "" {
		contentType = "application/octet-stream"
placeholder
	return resp.Body, contentType, nil
placeholder

func (c *GeminiBatchHTTPClient) DeleteFile(ctx context.Context, apiKey string, fileName string) error {
	req, err := c.newRequest(ctx, http.MethodDelete, "/v1beta/"+strings.TrimLeft(fileName, "/"), apiKey, nil)
	if err != nil {
		return err
placeholder
	return c.doNoBody(req)
placeholder

func (c *GeminiBatchHTTPClient) doBatchJob(req *http.Request) (*GeminiBatchJob, error) {
	var job GeminiBatchJob
	if err := c.doJSON(req, &job); err != nil {
		return nil, err
placeholder
	job.Raw = map[string]any{placeholder
	return &job, nil
placeholder

func (c *GeminiBatchHTTPClient) doNoBody(req *http.Request) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
placeholder
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return readGeminiAPIError(resp)
placeholder
	return nil
placeholder

func (c *GeminiBatchHTTPClient) doJSON(req *http.Request, out any) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
placeholder
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return readGeminiAPIError(resp)
placeholder
	return json.NewDecoder(resp.Body).Decode(out)
placeholder

func (c *GeminiBatchHTTPClient) newRequest(ctx context.Context, method, path, apiKey string, body io.Reader) (*http.Request, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, ErrBatchImageProviderMissingAPIKey
placeholder
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
placeholder
	req.Header.Set("x-goog-api-key", apiKey)
	return req, nil
placeholder

type GeminiAPIError struct {
	StatusCode int
	Code       string
	Message    string
placeholder

func (e *GeminiAPIError) Error() string {
	if e == nil {
		return "<nil>"
placeholder
	if e.Code != "" {
		return fmt.Sprintf("gemini api error: status=%d code=%s message=%s", e.StatusCode, e.Code, e.Message)
placeholder
	return fmt.Sprintf("gemini api error: status=%d message=%s", e.StatusCode, e.Message)
placeholder

func readGeminiAPIError(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
	message := string(body)
	var parsed struct {
		Error struct {
			Code    any    `json:"code"`
			Message string `json:"message"`
			Status  string `json:"status"`
	placeholder `json:"error"`
placeholder
	if err := json.Unmarshal(body, &parsed); err == nil && parsed.Error.Message != "" {
		message = parsed.Error.Message
		return &GeminiAPIError{StatusCode: resp.StatusCode, Code: parsed.Error.Status, Message: messageplaceholder
placeholder
	return &GeminiAPIError{StatusCode: resp.StatusCode, Message: messageplaceholder
placeholder

var _ BatchImageProvider = (*GeminiAPIBatchImageProvider)(nil)
var _ GeminiBatchClient = (*GeminiBatchHTTPClient)(nil)
