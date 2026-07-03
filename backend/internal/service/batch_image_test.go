//go:build unit

package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCanTransitionBatchImageJob(t *testing.T) {
	tests := []struct {
		name string
		from string
		to   string
		want bool
placeholder{
		{name: "created_to_uploading", from: BatchImageJobStatusCreated, to: BatchImageJobStatusUploading, want: trueplaceholder,
		{name: "uploading_to_submitted", from: BatchImageJobStatusUploading, to: BatchImageJobStatusSubmitted, want: trueplaceholder,
		{name: "submitted_to_running", from: BatchImageJobStatusSubmitted, to: BatchImageJobStatusRunning, want: trueplaceholder,
		{name: "running_self_poll", from: BatchImageJobStatusRunning, to: BatchImageJobStatusRunning, want: trueplaceholder,
		{name: "running_to_indexing", from: BatchImageJobStatusRunning, to: BatchImageJobStatusIndexing, want: trueplaceholder,
		{name: "indexing_to_settling", from: BatchImageJobStatusIndexing, to: BatchImageJobStatusSettling, want: trueplaceholder,
		{name: "settling_to_completed", from: BatchImageJobStatusSettling, to: BatchImageJobStatusCompleted, want: trueplaceholder,
		{name: "submitted_to_cancelled", from: BatchImageJobStatusSubmitted, to: BatchImageJobStatusCancelled, want: trueplaceholder,
		{name: "non_terminal_to_failed", from: BatchImageJobStatusCreated, to: BatchImageJobStatusFailed, want: trueplaceholder,
		{name: "completed_to_output_deleted", from: BatchImageJobStatusCompleted, to: BatchImageJobStatusOutputDeleted, want: trueplaceholder,
		{name: "failed_to_output_deleted", from: BatchImageJobStatusFailed, to: BatchImageJobStatusOutputDeleted, want: trueplaceholder,
		{name: "cancelled_to_output_deleted", from: BatchImageJobStatusCancelled, to: BatchImageJobStatusOutputDeleted, want: trueplaceholder,
		{name: "created_to_running_invalid", from: BatchImageJobStatusCreated, to: BatchImageJobStatusRunning, want: falseplaceholder,
		{name: "completed_to_running_invalid", from: BatchImageJobStatusCompleted, to: BatchImageJobStatusRunning, want: falseplaceholder,
		{name: "output_deleted_to_failed_invalid", from: BatchImageJobStatusOutputDeleted, to: BatchImageJobStatusFailed, want: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, CanTransitionBatchImageJob(tt.from, tt.to))
	placeholder)
placeholder
placeholder

func TestIsTerminalBatchImageJobStatus(t *testing.T) {
	require.True(t, IsTerminalBatchImageJobStatus(BatchImageJobStatusCompleted))
	require.True(t, IsTerminalBatchImageJobStatus(BatchImageJobStatusFailed))
	require.True(t, IsTerminalBatchImageJobStatus(BatchImageJobStatusCancelled))
	require.True(t, IsTerminalBatchImageJobStatus(BatchImageJobStatusOutputDeleted))
	require.False(t, IsTerminalBatchImageJobStatus(BatchImageJobStatusRunning))
placeholder

func TestIsSupportedBatchImageProvider(t *testing.T) {
	require.True(t, IsSupportedBatchImageProvider(BatchImageProviderGeminiAPI))
	require.True(t, IsSupportedBatchImageProvider(BatchImageProviderVertex))
	require.False(t, IsSupportedBatchImageProvider("gemini_oauth"))
	require.False(t, IsSupportedBatchImageProvider(""))
placeholder

func TestNewBatchImageID(t *testing.T) {
	id, err := NewBatchImageID()
placeholder
	require.True(t, strings.HasPrefix(id, "imgbatch_"))
	require.Len(t, id, len("imgbatch_")+32)
placeholder
