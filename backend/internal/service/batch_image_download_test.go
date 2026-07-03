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

func TestBatchImageDownloadService_OpenItemContent(t *testing.T) {
	ctx := context.Background()

	t.Run("streams image bytes with safe headers data", func(t *testing.T) {
		svc, _, limiter := newTestBatchImageDownloadService()

		stream, err := svc.OpenItemContent(ctx, testBatchImageOwner(), "imgbatch_download", "cover/../001", 1)
	placeholder
		defer stream.Reader.Close()

		body, err := io.ReadAll(stream.Reader)
	placeholder
		require.Equal(t, []byte("second"), body)
		require.Equal(t, "image/jpeg", stream.ContentType)
		require.Equal(t, "cover___001.jpg", stream.Filename)
		require.Equal(t, 1, limiter.acquireCount)
		require.Zero(t, limiter.releaseCount)
		require.NoError(t, stream.Reader.Close())
		require.Equal(t, 1, limiter.releaseCount)
placeholder)

	tests := []struct {
		name   string
		mutate func(*fakeBatchImageRepository)
		id     string
		item   string
		index  int
		want   error
placeholder{
		{name: "non_owner", id: "imgbatch_download", item: "cover/../001", mutate: func(r *fakeBatchImageRepository) {
			v := int64(999)
			r.jobs["imgbatch_download"].APIKeyID = &v
	placeholder, want: ErrBatchImageJobNotFoundplaceholder,
		{name: "not_completed", id: "imgbatch_download", item: "cover/../001", mutate: func(r *fakeBatchImageRepository) {
			r.jobs["imgbatch_download"].Status = BatchImageJobStatusRunning
	placeholder, want: ErrBatchImageNotReadyplaceholder,
		{name: "output_deleted", id: "imgbatch_download", item: "cover/../001", mutate: func(r *fakeBatchImageRepository) {
			r.jobs["imgbatch_download"].Status = BatchImageJobStatusOutputDeleted
	placeholder, want: ErrBatchImageOutputDeletedplaceholder,
		{name: "missing_item", id: "imgbatch_download", item: "missing", want: ErrBatchImageItemNotFoundplaceholder,
		{name: "failed_item", id: "imgbatch_download", item: "bad", want: ErrBatchImageItemFailedplaceholder,
		{name: "out_of_range", id: "imgbatch_download", item: "cover/../001", index: 2, want: ErrBatchImageItemImageIndexOutOfRangeplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _ := newTestBatchImageDownloadService()
			if tt.mutate != nil {
				tt.mutate(repo)
		placeholder

			got, err := svc.OpenItemContent(ctx, testBatchImageOwner(), tt.id, tt.item, tt.index)
			require.Nil(t, got)
			require.ErrorIs(t, err, tt.want)
			require.NotContains(t, err.Error(), batchImageDownloadTestBase64)
			require.NotContains(t, err.Error(), "providers/")
			require.NotContains(t, err.Error(), "gs://")
	placeholder)
placeholder
placeholder

func TestBatchImageDownloadService_StreamZip(t *testing.T) {
	ctx := context.Background()

	t.Run("streams zip with images manifest and errors", func(t *testing.T) {
		svc, _, limiter := newTestBatchImageDownloadService()
		var buf bytes.Buffer

		result, err := svc.StreamZip(ctx, testBatchImageOwner(), "imgbatch_download", BatchImageZipOptions{placeholder, &buf)
	placeholder
		require.Equal(t, 3, result.FileCount)
		require.Equal(t, 1, limiter.acquireCount)
		require.Equal(t, 1, limiter.releaseCount)

		files := readZipFiles(t, buf.Bytes())
		require.Equal(t, []byte("first"), files["images/cover___001.png"])
		require.Equal(t, []byte("second"), files["images/cover___001_2.jpg"])
		require.Equal(t, []byte("third"), files["images/ok_2.webp"])
		require.Contains(t, files, "manifest.json")
		require.Contains(t, files, "errors.json")

		zipText := string(bytes.Join(mapValues(files), []byte("\n")))
		require.NotContains(t, zipText, batchImageDownloadTestBase64)
		require.NotContains(t, zipText, "provider_job_name")
		require.NotContains(t, zipText, "provider_input_ref")
		require.NotContains(t, zipText, "gcs_output_uri")
		require.NotContains(t, zipText, "account_id")
		require.NotContains(t, zipText, "providers/")
		require.NotContains(t, zipText, "gs://")

		var manifest struct {
			Files []struct {
				CustomID   string `json:"custom_id"`
				Filename   string `json:"filename"`
				MimeType   string `json:"mime_type"`
				ImageIndex int    `json:"image_index"`
		placeholder `json:"files"`
	placeholder
		require.NoError(t, json.Unmarshal(files["manifest.json"], &manifest))
		require.Len(t, manifest.Files, 3)
		require.Equal(t, "images/cover___001_2.jpg", manifest.Files[1].Filename)
		require.Equal(t, 1, manifest.Files[1].ImageIndex)

		var errorsJSON []map[string]string
		require.NoError(t, json.Unmarshal(files["errors.json"], &errorsJSON))
		require.Len(t, errorsJSON, 1)
		require.Equal(t, "bad", errorsJSON[0]["custom_id"])
		require.Equal(t, "SAFETY_BLOCKED", errorsJSON[0]["code"])
placeholder)

	t.Run("limiter denial returns public limit error", func(t *testing.T) {
		svc, _, limiter := newTestBatchImageDownloadService()
		limiter.deny = true
		var buf bytes.Buffer

		result, err := svc.StreamZip(ctx, testBatchImageOwner(), "imgbatch_download", BatchImageZipOptions{placeholder, &buf)
		require.Nil(t, result)
		require.ErrorIs(t, err, ErrBatchImageDownloadLimited)
		require.Empty(t, buf.Bytes())
placeholder)

	t.Run("rejects too many zip items before opening output", func(t *testing.T) {
		svc, repo, _ := newTestBatchImageDownloadService()
		repo.jobs["imgbatch_download"].SuccessCount = 3
		svc.Config.BatchImage.MaxDownloadItemsZip = 1
		var buf bytes.Buffer

		result, err := svc.StreamZip(ctx, testBatchImageOwner(), "imgbatch_download", BatchImageZipOptions{placeholder, &buf)
		require.Nil(t, result)
		require.ErrorIs(t, err, ErrBatchImageZipTooManyItems)
		require.Empty(t, buf.Bytes())
placeholder)
placeholder

func TestExtractBatchImagePartsFromResultLine(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		wantID    string
		wantMime  string
		wantError string
placeholder{
		{name: "inlineData_mimeType_response", line: `{"key":"a","response":{"candidates":[{"content":{"parts":[{"inlineData":{"mimeType":"image/png","data":"` + batchImageDownloadTestBase64 + `"placeholderplaceholder]placeholderplaceholder]placeholderplaceholder`, wantID: "a", wantMime: "image/png"placeholder,
		{name: "inline_data_mime_type_top_level", line: `{"custom_id":"b","candidates":[{"content":{"parts":[{"inline_data":{"mime_type":"image/jpeg","data":"` + batchImageDownloadTestBase64 + `"placeholderplaceholder]placeholderplaceholder]placeholder`, wantID: "b", wantMime: "image/jpeg"placeholder,
		{name: "status_failure", line: `{"key":"c","status":{"code":"INVALID_ARGUMENT","message":"bad prompt"placeholderplaceholder`, wantID: "c", wantError: "INVALID_ARGUMENT"placeholder,
		{name: "error_failure", line: `{"key":"d","error":{"code":"SAFETY","message":"blocked"placeholderplaceholder`, wantID: "d", wantError: "SAFETY_BLOCKED"placeholder,
		{name: "empty_output", line: `{"key":"e","response":{"candidates":[{"content":{"parts":[{"text":"none"placeholder]placeholderplaceholder]placeholderplaceholder`, wantID: "e", wantError: "EMPTY_IMAGE_OUTPUT"placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractBatchImagePartsFromResultLine([]byte(tt.line))
		placeholder
			require.Equal(t, tt.wantID, got.CustomID)
			if tt.wantMime != "" {
				require.Len(t, got.Images, 1)
				require.Equal(t, tt.wantMime, got.Images[0].MimeType)
				require.NotEmpty(t, got.Images[0].Base64Data)
		placeholder
			if tt.wantError != "" {
				require.Equal(t, tt.wantError, got.ErrorCode)
		placeholder
	placeholder)
placeholder

	_, err := ExtractBatchImagePartsFromResultLine([]byte(`{"response":{"candidates":[{"content":{"parts":[{"inlineData":{"mimeType":"image/png","data":"` + batchImageDownloadTestBase64 + `"placeholderplaceholder]placeholderplaceholder]placeholderplaceholder`))
placeholder
	require.NotContains(t, err.Error(), batchImageDownloadTestBase64)
placeholder

func TestBatchImageDownloadFilenames(t *testing.T) {
	require.Equal(t, "___secret_name.png", BatchImageSafeDownloadFilename("../../secret\nname", "png"))
	require.Equal(t, `attachment; filename="cover_001.png"`, BatchImageContentDispositionAttachment(`cover"001.png`))
placeholder

func newTestBatchImageDownloadService() (*BatchImageDownloadService, *fakeBatchImageRepository, *fakeBatchImageDownloadLimiter) {
	repo := newFakeBatchImageRepository()
	apiKeyID := int64(22)
	accountID := int64(101)
	repo.jobs["imgbatch_download"] = &BatchImageJob{
		BatchID:           "imgbatch_download",
		UserID:            11,
		APIKeyID:          &apiKeyID,
		AccountID:         &accountID,
		Provider:          BatchImageProviderGeminiAPI,
		Model:             "gemini-2.5-flash-image",
		Status:            BatchImageJobStatusCompleted,
		ProviderJobName:   batchImageStringPtr("providers/internal/job"),
		ProviderOutputRef: batchImageStringPtr("gs://bucket/internal/output.jsonl"),
		ItemCount:         3,
		SuccessCount:      2,
		FailCount:         1,
		CreatedAt:         time.Now(),
placeholder
	mime := "image/png"
	ext := "png"
	webp := "image/webp"
	webpExt := "webp"
	code := "SAFETY_BLOCKED"
	msg := "blocked in gs://bucket/internal/output.jsonl"
	repo.items["imgbatch_download"] = []CreateBatchImageItemParams{
		{JobID: "imgbatch_download", CustomID: "cover/../001", Status: BatchImageItemStatusSuccess, MimeType: &mime, FileExtension: &ext, ImageCount: 2placeholder,
		{JobID: "imgbatch_download", CustomID: "bad", Status: BatchImageItemStatusFailed, ErrorCode: &code, ErrorMessage: &msgplaceholder,
		{JobID: "imgbatch_download", CustomID: "ok_2", Status: BatchImageItemStatusSuccess, MimeType: &webp, FileExtension: &webpExt, ImageCount: 1placeholder,
placeholder
	provider := &publicBatchImageProvider{name: BatchImageProviderGeminiAPI, result: batchImageDownloadResultJSONL()placeholder
	limiter := &fakeBatchImageDownloadLimiter{placeholder
	svc := &BatchImageDownloadService{
		Repo:             repo,
		ProviderRegistry: NewBatchImageProviderRegistry(provider),
		AccountResolver:  &fakeBatchImageAccountResolver{account: &Account{ID: accountID, Platform: PlatformGemini, Type: AccountTypeAPIKey, Status: StatusActive, Schedulable: trueplaceholderplaceholder,
		Limiter:          limiter,
		Config:           &config.Config{BatchImage: config.BatchImageConfig{MaxDownloadItemsZip: 10, MaxDownloadDurationSeconds: 60placeholderplaceholder,
placeholder
	return svc, repo, limiter
placeholder

const batchImageDownloadTestBase64 = "Zmlyc3Q="

func batchImageDownloadResultJSONL() string {
	return strings.Join([]string{
		`{"key":"cover/../001","response":{"candidates":[{"content":{"parts":[{"inlineData":{"mimeType":"image/png","data":"Zmlyc3Q="placeholderplaceholder,{"inlineData":{"mimeType":"image/jpeg","data":"c2Vjb25k"placeholderplaceholder]placeholderplaceholder]placeholderplaceholder`,
		`{"key":"bad","error":{"code":"SAFETY","message":"blocked"placeholderplaceholder`,
		`{"key":"ok_2","candidates":[{"content":{"parts":[{"inline_data":{"mime_type":"image/webp","data":"dGhpcmQ="placeholderplaceholder]placeholderplaceholder]placeholder`,
placeholder, "\n") + "\n"
placeholder

func readZipFiles(t *testing.T, data []byte) map[string][]byte {
placeholder
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
placeholder
	out := make(map[string][]byte, len(reader.File))
	for _, file := range reader.File {
		rc, err := file.Open()
	placeholder
		body, err := io.ReadAll(rc)
	placeholder
		require.NoError(t, rc.Close())
		out[file.Name] = body
placeholder
	return out
placeholder

func mapValues(in map[string][]byte) [][]byte {
	out := make([][]byte, 0, len(in))
	for _, value := range in {
		out = append(out, value)
placeholder
	return out
placeholder

type fakeBatchImageDownloadLimiter struct {
	acquireCount int
	releaseCount int
	deny         bool
placeholder

func (l *fakeBatchImageDownloadLimiter) Acquire(context.Context, string, string) (BatchImageDownloadPermit, error) {
	l.acquireCount++
	if l.deny {
		return nil, ErrBatchImageDownloadLimited
placeholder
	return &fakeBatchImageDownloadPermit{release: func() { l.releaseCount++ placeholderplaceholder, nil
placeholder

type fakeBatchImageDownloadPermit struct {
	once    bool
	release func()
placeholder

func (p *fakeBatchImageDownloadPermit) Release(context.Context) error {
	if p.once {
		return nil
placeholder
	p.once = true
	if p.release != nil {
		p.release()
placeholder
	return nil
placeholder

var _ BatchImageDownloadLimiter = (*fakeBatchImageDownloadLimiter)(nil)
var _ BatchImageDownloadPermit = (*fakeBatchImageDownloadPermit)(nil)
