//go:build unit

package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestBatchImageProviderRegistry_ReturnsVertex(t *testing.T) {
	registry := NewDefaultBatchImageProviderRegistry()
	provider, ok := registry.Get(BatchImageProviderVertex)
	require.True(t, ok)
	require.Equal(t, BatchImageProviderVertex, provider.Name())
placeholder

func TestVertexProvider_SupportsOnlyGeminiServiceAccount(t *testing.T) {
	provider := newTestVertexProvider(&fakeVertexBatchClient{placeholder, &fakeVertexObjectStore{placeholder)

	require.True(t, provider.SupportsAccount(vertexServiceAccount()))
	require.False(t, provider.SupportsAccount(&Account{Platform: PlatformGemini, Type: AccountTypeAPIKey, Credentials: map[string]any{"api_key": "sk"placeholderplaceholder))
	require.False(t, provider.SupportsAccount(&Account{Platform: PlatformGemini, Type: AccountTypeOAuth, Credentials: map[string]any{"access_token": "tok"placeholderplaceholder))
	require.False(t, provider.SupportsAccount(&Account{Platform: PlatformAnthropic, Type: AccountTypeServiceAccount, Credentials: vertexServiceAccount().Credentialsplaceholder))
	require.False(t, provider.SupportsAccount(&Account{Platform: PlatformGemini, Type: AccountTypeServiceAccount, Credentials: map[string]any{placeholderplaceholder))
placeholder

func TestVertexProvider_MissingServiceAccountRejected(t *testing.T) {
	provider := newTestVertexProvider(&fakeVertexBatchClient{placeholder, &fakeVertexObjectStore{placeholder)
	_, err := provider.Submit(context.Background(), nil, &Account{Platform: PlatformGemini, Type: AccountTypeServiceAccount, Credentials: map[string]any{placeholderplaceholder, validVertexBatchInput())
	require.ErrorIs(t, err, ErrBatchImageProviderMissingServiceAccount)
placeholder

func TestVertexProvider_MissingManagedGCSBucketRejected(t *testing.T) {
	provider := NewVertexBatchImageProvider(VertexBatchImageProviderOptions{ProjectID: "proj", Environment: "test"placeholder, &fakeVertexBatchClient{placeholder, &fakeVertexObjectStore{placeholder, &fakeGeminiTokenCache{token: "token"placeholder)
	_, err := provider.Submit(context.Background(), nil, vertexServiceAccount(), validVertexBatchInput())
placeholder
	require.Equal(t, "VERTEX_MANAGED_GCS_BUCKET_MISSING", infraerrors.Reason(err))
placeholder

func TestBuildVertexBatchJSONL_WritesValidLinesAndPreservesCustomID(t *testing.T) {
	input := validVertexBatchInput()
	input.Items = append(input.Items, BatchImageInputItem{CustomID: "cover_002", Prompt: "Second prompt"placeholder)

	jsonl, err := BuildVertexBatchJSONL(input)
placeholder
	lines := strings.Split(strings.TrimSpace(string(jsonl)), "\n")
	require.Len(t, lines, 2)
	requireVertexJSONLLine(t, lines[0], "cover_001", "A clean product hero image")
	requireVertexJSONLLine(t, lines[1], "cover_002", "Second prompt")
placeholder

func TestBuildVertexBatchJSONL_RejectsDuplicateCustomIDs(t *testing.T) {
	input := validVertexBatchInput()
	input.Items = append(input.Items, BatchImageInputItem{CustomID: "cover_001", Prompt: "Duplicate"placeholder)
	_, err := BuildVertexBatchJSONL(input)
	require.ErrorIs(t, err, ErrBatchImageProviderInvalidInput)
placeholder

func TestBuildVertexBatchJSONL_RejectsEmptyPrompt(t *testing.T) {
	input := validVertexBatchInput()
	input.Items[0].Prompt = " "
	_, err := BuildVertexBatchJSONL(input)
	require.ErrorIs(t, err, ErrBatchImageProviderInvalidInput)
placeholder

func TestBuildVertexBatchJSONL_WritesReferenceImages(t *testing.T) {
	input := validVertexBatchInput()
	input.Items[0].ReferenceImages = []BatchImageReference{
		{MimeType: "image/png", Data: []byte("png-bytes")placeholder,
		{MimeType: "image/jpeg", FileURI: "gs://bucket/refs/style.jpg"placeholder,
placeholder

	jsonl, err := BuildVertexBatchJSONL(input)
placeholder
	lines := strings.Split(strings.TrimSpace(string(jsonl)), "\n")
	require.Len(t, lines, 1)

placeholder
	require.NoError(t, json.Unmarshal([]byte(lines[0]), &got))
placeholder
placeholder
placeholder
	require.Len(t, parts, 3)
	require.Equal(t, "A clean product hero image", parts[0].(map[string]any)["text"])
	inlineData := parts[1].(map[string]any)["inlineData"].(map[string]any)
	require.Equal(t, "image/png", inlineData["mimeType"])
	require.Equal(t, "cG5nLWJ5dGVz", inlineData["data"])
	fileData := parts[2].(map[string]any)["fileData"].(map[string]any)
	require.Equal(t, "image/jpeg", fileData["mimeType"])
	require.Equal(t, "gs://bucket/refs/style.jpg", fileData["fileUri"])
placeholder

func TestNormalizeVertexBatchModelPath(t *testing.T) {
	require.Equal(t, "publishers/google/models/gemini-3.1-flash-image", NormalizeVertexBatchModelPath("gemini-3.1-flash-image"))
	require.Equal(t, "publishers/google/models/gemini-2.5-flash-image", NormalizeVertexBatchModelPath("publishers/google/models/gemini-2.5-flash-image"))
	require.Equal(t, "projects/p/locations/global/models/m", NormalizeVertexBatchModelPath("projects/p/locations/global/models/m"))
placeholder

func TestBuildVertexBatchPredictionJobsEndpoint(t *testing.T) {
	global, err := BuildVertexBatchPredictionJobsEndpoint("", "my-project", "global")
placeholder
	require.Equal(t, "https://aiplatform.googleapis.com/v1/projects/my-project/locations/global/batchPredictionJobs", global)

	regional, err := BuildVertexBatchPredictionJobsEndpoint("", "my-project", "asia-northeast1")
placeholder
	require.Equal(t, "https://asia-northeast1-aiplatform.googleapis.com/v1/projects/my-project/locations/asia-northeast1/batchPredictionJobs", regional)
placeholder

func TestVertexProvider_SubmitUploadsJSONLAndCreatesBatchPredictionJob(t *testing.T) {
placeholder
	store := &fakeVertexObjectStore{placeholder
	provider := newTestVertexProvider(vertexClient, store)

	got, err := provider.Submit(context.Background(), &BatchImageJob{BatchID: "imgbatch_abc123", Model: "gemini-3.1-flash-image"placeholder, vertexServiceAccount(), validVertexBatchInput())
placeholder

	require.Equal(t, "gs://managed-bucket/batch-image/test/imgbatch_abc123/input/requests.jsonl", store.uploadURI)
	require.Equal(t, "projects/proj/locations/global/batchPredictionJobs/job-1", got.ProviderJobName)
	require.Equal(t, store.uploadURI, got.ProviderInputRef)
	require.Equal(t, "gs://managed-bucket/batch-image/test/imgbatch_abc123/output/", got.ProviderOutputRef)
	require.Equal(t, "jsonl", vertexClient.createdReq.InputConfig.InstancesFormat)
	require.Equal(t, "jsonl", vertexClient.createdReq.OutputConfig.PredictionsFormat)
	require.Equal(t, got.ProviderOutputRef, vertexClient.createdReq.OutputConfig.GCSDestination.OutputURIPrefix)
	require.Equal(t, "key", vertexClient.createdReq.InstanceConfig.KeyField)
	require.NotContains(t, string(vertexClient.createdPayloadForAssert(t)), "serviceAccount")
	require.NotContains(t, string(vertexClient.createdPayloadForAssert(t)), "encryptionSpec")
	require.NotContains(t, got.ProviderInputRef+got.ProviderOutputRef+got.ProviderJobName, "A clean product hero image")
	require.NotContains(t, string(store.uploadedJSONL), "private_key")
placeholder

func TestVertexProvider_GetMapsStates(t *testing.T) {
	tests := []struct {
		name      string
		state     string
		err       *VertexBatchJobError
		wantState BatchProviderInternalState
		wantDone  bool
		wantCode  string
placeholder{
		{name: "pending", state: "JOB_STATE_PENDING", wantState: BatchProviderStateQueuedplaceholder,
		{name: "queued", state: "JOB_STATE_QUEUED", wantState: BatchProviderStateQueuedplaceholder,
		{name: "running", state: "JOB_STATE_RUNNING", wantState: BatchProviderStateRunningplaceholder,
		{name: "succeeded", state: "JOB_STATE_SUCCEEDED", wantState: BatchProviderStateSucceeded, wantDone: trueplaceholder,
		{name: "failed", state: "JOB_STATE_FAILED", err: &VertexBatchJobError{Status: "INVALID_ARGUMENT", Message: "bad request"placeholder, wantState: BatchProviderStateFailed, wantDone: true, wantCode: "INVALID_ARGUMENT"placeholder,
		{name: "cancelled", state: "JOB_STATE_CANCELLED", wantState: BatchProviderStateCancelled, wantDone: true, wantCode: "VERTEX_BATCH_CANCELLED"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := "gs://managed-bucket/batch-image/test/imgbatch_abc123/output/"
			provider := newTestVertexProvider(&fakeVertexBatchClient{got: &VertexBatchPredictionJob{
				Name:         "projects/proj/locations/global/batchPredictionJobs/job-1",
				State:        tt.state,
				Error:        tt.err,
				OutputConfig: VertexBatchOutputConfig{GCSDestination: VertexBatchGCSDestination{OutputURIPrefix: outputplaceholderplaceholder,
	placeholder &fakeVertexObjectStore{placeholder)
			got, err := provider.Get(context.Background(), vertexJobWithName("projects/proj/locations/global/batchPredictionJobs/job-1"), vertexServiceAccount())
		placeholder
			require.Equal(t, tt.wantState, got.InternalState)
			require.Equal(t, tt.wantDone, got.Done)
			require.Equal(t, output, got.ProviderOutputRef)
			require.Equal(t, tt.wantCode, got.ErrorCode)
	placeholder)
placeholder
placeholder

func TestVertexProvider_OpenResultReturnsCombinedJSONLStream(t *testing.T) {
	output := "gs://managed-bucket/batch-image/test/imgbatch_abc123/output/"
	store := &fakeVertexObjectStore{
		listed: []string{
			output + "predictions_2.jsonl",
			output + "predictions_1.jsonl",
	placeholder,
		objects: map[string]string{
			output + "predictions_1.jsonl": `{"key":"1"placeholder` + "\n",
			output + "predictions_2.jsonl": `{"key":"2"placeholder` + "\n",
	placeholder,
placeholder
	provider := newTestVertexProvider(&fakeVertexBatchClient{placeholder, store)
	r, contentType, err := provider.OpenResult(context.Background(), &BatchImageJob{ProviderOutputRef: &outputplaceholder, vertexServiceAccount())
placeholder
	defer r.Close()

	body, err := io.ReadAll(r)
placeholder
	require.Equal(t, "application/jsonl", contentType)
	require.Equal(t, "{\"key\":\"1\"placeholder\n\n{\"key\":\"2\"placeholder\n", string(body))
placeholder

func TestVertexProvider_OpenResultMissingObjectsReturnsTypedError(t *testing.T) {
	output := "gs://managed-bucket/batch-image/test/imgbatch_abc123/output/"
	provider := newTestVertexProvider(&fakeVertexBatchClient{placeholder, &fakeVertexObjectStore{placeholder)
	_, _, err := provider.OpenResult(context.Background(), &BatchImageJob{ProviderOutputRef: &outputplaceholder, vertexServiceAccount())
placeholder
	require.Equal(t, "VERTEX_RESULT_OBJECTS_MISSING", infraerrors.Reason(err))
placeholder

func TestVertexProvider_CancelCallsClient(t *testing.T) {
	vertexClient := &fakeVertexBatchClient{placeholder
placeholder

	err := provider.Cancel(context.Background(), vertexJobWithName("projects/proj/locations/global/batchPredictionJobs/job-1"), vertexServiceAccount())
placeholder
	require.Equal(t, "projects/proj/locations/global/batchPredictionJobs/job-1", vertexClient.cancelledName)
placeholder

func TestVertexProvider_CleanupDeletesOnlyManagedPaths(t *testing.T) {
	input := "gs://managed-bucket/batch-image/test/imgbatch_abc123/input/requests.jsonl"
	output := "gs://managed-bucket/batch-image/test/imgbatch_abc123/output/"
	store := &fakeVertexObjectStore{placeholder
	provider := newTestVertexProvider(&fakeVertexBatchClient{placeholder, store)

	err := provider.Cleanup(context.Background(), &BatchImageJob{BatchID: "imgbatch_abc123", ProviderInputRef: &input, ProviderOutputRef: &outputplaceholder, vertexServiceAccount(), CleanupTargetAll)
placeholder
	require.Equal(t, []string{inputplaceholder, store.deletedObjects)
	require.Equal(t, []string{outputplaceholder, store.deletedPrefixes)
placeholder

func TestVertexProvider_CleanupRejectsUnsafePath(t *testing.T) {
	input := "gs://other-bucket/batch-image/test/imgbatch_abc123/input/requests.jsonl"
	provider := newTestVertexProvider(&fakeVertexBatchClient{placeholder, &fakeVertexObjectStore{placeholder)

	err := provider.Cleanup(context.Background(), &BatchImageJob{BatchID: "imgbatch_abc123", ProviderInputRef: &inputplaceholder, vertexServiceAccount(), CleanupTargetInput)
	require.ErrorIs(t, err, ErrBatchImageProviderUnsafeCleanupPath)
placeholder

func TestVertexProvider_ErrorsDoNotExposeServiceAccountSecrets(t *testing.T) {
	privateKey := "placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder

placeholder
placeholder
placeholder
placeholder
placeholder
placeholder

placeholder
placeholder
placeholder

placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder

placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder

placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder

placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder, client, store, &fakeGeminiTokenCache{token: "ya29.test-token"placeholder)
placeholder

placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder
placeholder\nabc\n-----END PRIVATE KEY-----\n",
		placeholder
		placeholder,
	placeholder,
placeholder
placeholder

func vertexJobWithName(name string) *BatchImageJob {
	return &BatchImageJob{ProviderJobName: &nameplaceholder
placeholder

type fakeVertexBatchClient struct {
	created       *VertexBatchPredictionJob
	got           *VertexBatchPredictionJob
	createErr     error
	getErr        error
	cancelErr     error
	createdReq    VertexCreateBatchPredictionJobRequest
	cancelledName string
placeholder

func (f *fakeVertexBatchClient) CreateBatchPredictionJob(_ context.Context, accessToken string, req VertexCreateBatchPredictionJobRequest) (*VertexBatchPredictionJob, error) {
	if strings.TrimSpace(accessToken) == "" {
		return nil, errors.New("missing token")
placeholder
	f.createdReq = req
	if f.createErr != nil {
		return nil, f.createErr
placeholder
	if f.created != nil {
		return f.created, nil
placeholder
	return &VertexBatchPredictionJob{Name: "projects/proj/locations/global/batchPredictionJobs/job-1", State: "JOB_STATE_PENDING"placeholder, nil
placeholder

func (f *fakeVertexBatchClient) GetBatchPredictionJob(_ context.Context, _ string, _ string) (*VertexBatchPredictionJob, error) {
	if f.getErr != nil {
		return nil, f.getErr
placeholder
	return f.got, nil
placeholder

func (f *fakeVertexBatchClient) CancelBatchPredictionJob(_ context.Context, _ string, name string) error {
	f.cancelledName = name
	return f.cancelErr
placeholder

func (f *fakeVertexBatchClient) createdPayloadForAssert(t *testing.T) []byte {
placeholder
	b, err := json.Marshal(f.createdReq)
placeholder
	return b
placeholder

type fakeVertexObjectStore struct {
	uploadURI       string
	uploadedJSONL   []byte
	uploadErr       error
	listed          []string
	objects         map[string]string
	listErr         error
	openErr         error
	deleteErr       error
	deletedObjects  []string
	deletedPrefixes []string
placeholder

func (f *fakeVertexObjectStore) UploadJSONL(_ context.Context, _ string, uri string, r io.Reader) error {
	f.uploadURI = uri
	f.uploadedJSONL, _ = io.ReadAll(r)
	return f.uploadErr
placeholder

func (f *fakeVertexObjectStore) ListJSONLObjects(_ context.Context, _ string, _ string) ([]string, error) {
	if f.listErr != nil {
		return nil, f.listErr
placeholder
	out := make([]string, 0, len(f.listed))
	for _, item := range f.listed {
		if strings.HasSuffix(item, ".jsonl") {
			out = append(out, item)
	placeholder
placeholder
	return out, nil
placeholder

func (f *fakeVertexObjectStore) OpenObject(_ context.Context, _ string, uri string) (io.ReadCloser, string, error) {
	if f.openErr != nil {
		return nil, "", f.openErr
placeholder
	return io.NopCloser(bytes.NewBufferString(f.objects[uri])), "application/jsonl", nil
placeholder

func (f *fakeVertexObjectStore) DeleteObject(_ context.Context, _ string, uri string) error {
	f.deletedObjects = append(f.deletedObjects, uri)
	return f.deleteErr
placeholder

func (f *fakeVertexObjectStore) DeletePrefix(_ context.Context, _ string, uri string) error {
	f.deletedPrefixes = append(f.deletedPrefixes, uri)
	return f.deleteErr
placeholder

type fakeGeminiTokenCache struct {
	token string
placeholder

func (f *fakeGeminiTokenCache) GetAccessToken(context.Context, string) (string, error) {
	if strings.TrimSpace(f.token) == "" {
		return "", errors.New("missing token")
placeholder
	return f.token, nil
placeholder

func (f *fakeGeminiTokenCache) SetAccessToken(context.Context, string, string, time.Duration) error {
	return nil
placeholder

func (f *fakeGeminiTokenCache) DeleteAccessToken(context.Context, string) error {
	return nil
placeholder

func (f *fakeGeminiTokenCache) AcquireRefreshLock(context.Context, string, time.Duration) (bool, error) {
	return false, nil
placeholder

func (f *fakeGeminiTokenCache) ReleaseRefreshLock(context.Context, string) error {
	return nil
placeholder
