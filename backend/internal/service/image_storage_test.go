package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// pngBytes is a minimal payload whose signature makes http.DetectContentType
// report image/png.
var pngBytes = []byte("\x89PNG\r\n\x1a\nfake-png-payload")

type savedImage struct {
	key         string
	contentType string
	data        []byte
placeholder

type fakeImageStorage struct {
	saved []savedImage
	url   string
	err   error
placeholder

func (f *fakeImageStorage) Save(_ context.Context, key, contentType string, data []byte) (string, error) {
	if f.err != nil {
		return "", f.err
placeholder
	f.saved = append(f.saved, savedImage{key: key, contentType: contentType, data: append([]byte(nil), data...)placeholder)
	if f.url != "" {
		return f.url, nil
placeholder
	return "https://cdn.test/" + key, nil
placeholder

func TestImageResultUploaderRewritesB64JSON(t *testing.T) {
	storage := &fakeImageStorage{placeholder
	uploader := NewImageResultUploader(storage, "images/", 0, nil)

	b64 := base64.StdEncoding.EncodeToString(pngBytes)
	result := json.RawMessage(`{"created":1,"data":[{"b64_json":"` + b64 + `","revised_prompt":"a cat"placeholder]placeholder`)

	out, err := uploader.Rewrite(context.Background(), "imgtask_abc", result)
placeholder

	require.Len(t, storage.saved, 1)
	require.Equal(t, "images/imgtask_abc-0.png", storage.saved[0].key)
	require.Equal(t, "image/png", storage.saved[0].contentType)
	require.Equal(t, pngBytes, storage.saved[0].data)

	var parsed struct {
		Data []map[string]json.RawMessage `json:"data"`
placeholder
	require.NoError(t, json.Unmarshal(out, &parsed))
	require.Len(t, parsed.Data, 1)
	require.JSONEq(t, `"https://cdn.test/images/imgtask_abc-0.png"`, string(parsed.Data[0]["url"]))
	_, hasB64 := parsed.Data[0]["b64_json"]
	require.False(t, hasB64, "b64_json must be stripped after offload")
	require.JSONEq(t, `"a cat"`, string(parsed.Data[0]["revised_prompt"]), "unrelated fields preserved")
placeholder

func TestImageResultUploaderRewritesURL(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(pngBytes)
placeholder))
	defer upstream.Close()

	storage := &fakeImageStorage{placeholder
	uploader := NewImageResultUploader(storage, "images/", 0, nil)

	result := json.RawMessage(`{"created":1,"data":[{"url":"` + upstream.URL + `/pic.png"placeholder]placeholder`)
	out, err := uploader.Rewrite(context.Background(), "imgtask_xyz", result)
placeholder

	require.Len(t, storage.saved, 1)
	require.Equal(t, pngBytes, storage.saved[0].data)
	require.Equal(t, "image/png", storage.saved[0].contentType)

	var parsed struct {
		Data []map[string]json.RawMessage `json:"data"`
placeholder
	require.NoError(t, json.Unmarshal(out, &parsed))
	require.JSONEq(t, `"https://cdn.test/images/imgtask_xyz-0.png"`, string(parsed.Data[0]["url"]))
placeholder

func TestImageResultUploaderPropagatesStorageError(t *testing.T) {
	storage := &fakeImageStorage{err: errors.New("bucket unreachable")placeholder
	uploader := NewImageResultUploader(storage, "images/", 0, nil)

	b64 := base64.StdEncoding.EncodeToString(pngBytes)
	result := json.RawMessage(`{"data":[{"b64_json":"` + b64 + `"placeholder]placeholder`)

	_, err := uploader.Rewrite(context.Background(), "imgtask_err", result)
placeholder
	require.Contains(t, err.Error(), "bucket unreachable")
placeholder

func TestImageResultUploaderNilStoragePassthrough(t *testing.T) {
	var uploader *ImageResultUploader
	result := json.RawMessage(`{"data":[{"url":"https://example.test/x.png"placeholder]placeholder`)
	out, err := uploader.Rewrite(context.Background(), "imgtask_nil", result)
placeholder
	require.JSONEq(t, string(result), string(out))
placeholder

func TestImageTaskServiceCompleteOffloadsToStorage(t *testing.T) {
	store := &imageTaskMemoryStore{placeholder
	storage := &fakeImageStorage{placeholder
	uploader := NewImageResultUploader(storage, "images/", 0, nil)
	svc := NewImageTaskServiceWithUploader(store, uploader, time.Hour, time.Minute)
	require.True(t, svc.Enabled())

	owner := ImageTaskOwner{UserID: 1, APIKeyID: 2placeholder
	created, err := svc.Create(context.Background(), owner)
placeholder

	b64 := base64.StdEncoding.EncodeToString(pngBytes)
	result := json.RawMessage(`{"created":1,"data":[{"b64_json":"` + b64 + `"placeholder]placeholder`)
	require.NoError(t, svc.Complete(context.Background(), created.ID, http.StatusOK, result))

	got, err := svc.Get(context.Background(), owner, created.ID)
placeholder
	require.Equal(t, ImageTaskStatusCompleted, got.Status)
	require.Equal(t, "https://cdn.test/images/"+created.ID+"-0.png", got.ImageURL)
	require.NotContains(t, string(got.Result), "b64_json", "large base64 must not be persisted to Redis")
	require.Len(t, storage.saved, 1)
placeholder

func TestImageTaskServiceCompleteOffloadFailureMarksFailed(t *testing.T) {
	store := &imageTaskMemoryStore{placeholder
	storage := &fakeImageStorage{err: errors.New("bucket unreachable")placeholder
	uploader := NewImageResultUploader(storage, "images/", 0, nil)
	svc := NewImageTaskServiceWithUploader(store, uploader, time.Hour, time.Minute)

	owner := ImageTaskOwner{UserID: 1, APIKeyID: 2placeholder
	created, err := svc.Create(context.Background(), owner)
placeholder

	b64 := base64.StdEncoding.EncodeToString(pngBytes)
	result := json.RawMessage(`{"data":[{"b64_json":"` + b64 + `"placeholder]placeholder`)
	require.NoError(t, svc.Complete(context.Background(), created.ID, http.StatusOK, result))

	got, err := svc.Get(context.Background(), owner, created.ID)
placeholder
	require.Equal(t, ImageTaskStatusFailed, got.Status)
	require.Equal(t, http.StatusBadGateway, got.HTTPStatus)
	require.Contains(t, string(got.Error), "object storage")
	require.NotContains(t, string(got.Result), "b64_json", "failed offload must not persist base64 to Redis")
placeholder
