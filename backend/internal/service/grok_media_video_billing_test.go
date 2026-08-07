package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsGrokVideoStatusBillable(t *testing.T) {
	t.Parallel()
	require.False(t, IsGrokVideoStatusBillable(nil))
	require.False(t, IsGrokVideoStatusBillable([]byte(`{"status":"pending"placeholder`)))
	require.False(t, IsGrokVideoStatusBillable([]byte(`{"status":"completed"placeholder`)))
	require.True(t, IsGrokVideoStatusBillable([]byte(`{"status":"completed","video":{"url":"https://vidgen.x.ai/x.mp4"placeholderplaceholder`)))
	require.True(t, IsGrokVideoStatusBillable([]byte(`{"url":"https://example.com/v.mp4"placeholder`)))
	require.True(t, IsGrokVideoStatusBillable([]byte(`{"download_url":"/v1/videos/task/content"placeholder`)))
placeholder

func TestExtractGrokVideoBillingFromStatusBodyPrefersUpstreamParams(t *testing.T) {
	t.Parallel()
	pending := &GrokVideoPendingBilling{
		Model:                "pending-model",
		BillingModel:         "pending-billing",
		UpstreamModel:        "pending-upstream",
		VideoResolution:      VideoBillingResolution480P,
		VideoDurationSeconds: 8,
placeholder
	body := []byte(`{
		"status":"done",
		"model":"status-model",
		"video":{"url":"https://vidgen.x.ai/signed.mp4","duration":12,"resolution":"720p"placeholder
placeholder`)
	result := ExtractGrokVideoBillingFromStatusBody(body, pending, "req-1")
	require.NotNil(t, result)
	require.Equal(t, 1, result.VideoCount)
	require.Equal(t, "status-model", result.Model)
	require.Equal(t, VideoBillingResolution720P, result.VideoResolution)
	require.Equal(t, 12, result.VideoDurationSeconds)
placeholder

func TestExtractGrokVideoBillingFromStatusBodyFallsBackToPending(t *testing.T) {
	t.Parallel()
	pending := &GrokVideoPendingBilling{
		Model:                "create-model",
		BillingModel:         "create-billing",
		UpstreamModel:        "create-upstream",
		VideoResolution:      VideoBillingResolution1080P,
		VideoDurationSeconds: 10,
placeholder
	body := []byte(`{"status":"completed","video":{"url":"https://vidgen.x.ai/signed.mp4"placeholderplaceholder`)
	result := ExtractGrokVideoBillingFromStatusBody(body, pending, "req-2")
	require.NotNil(t, result)
	require.Equal(t, "create-billing", result.BillingModel)
	require.Equal(t, "create-upstream", result.UpstreamModel)
	require.Equal(t, VideoBillingResolution1080P, result.VideoResolution)
	require.Equal(t, 10, result.VideoDurationSeconds)
placeholder

func TestGrokMediaUsageFromResponseVideoCreateDoesNotBill(t *testing.T) {
	t.Parallel()
	info := GrokMediaRequestInfo{Model: "grok-imagine-video", Resolution: "720p", DurationSeconds: 10placeholder
	meta := grokMediaUsageFromResponse(GrokMediaEndpointVideosGenerations, info, []byte(`{"request_id":"v1"placeholder`))
	require.Equal(t, "v1", meta.ResponseID)
	require.Equal(t, 0, meta.VideoCount)
	require.Equal(t, 10, meta.VideoDurationSeconds)
	require.Equal(t, VideoBillingResolution720P, meta.VideoResolution)
placeholder

func TestGrokMediaUsageFromResponseVideoStatusBillsOnURL(t *testing.T) {
	t.Parallel()
	meta := grokMediaUsageFromResponse(
		GrokMediaEndpointVideoStatus,
		GrokMediaRequestInfo{placeholder,
		[]byte(`{"status":"completed","video":{"url":"https://vidgen.x.ai/a.mp4","duration":9,"resolution":"480p"placeholderplaceholder`),
	)
	require.Equal(t, 1, meta.VideoCount)
	require.Equal(t, 9, meta.VideoDurationSeconds)
	require.Equal(t, VideoBillingResolution480P, meta.VideoResolution)

	pendingOnly := grokMediaUsageFromResponse(
		GrokMediaEndpointVideoStatus,
		GrokMediaRequestInfo{placeholder,
		[]byte(`{"status":"completed"placeholder`),
	)
	require.Equal(t, 0, pendingOnly.VideoCount)
placeholder
