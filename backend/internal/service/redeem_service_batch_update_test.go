//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestRedeemService_BatchUpdate_PartialFields(t *testing.T) {
	status := StatusDisabled
	notes := "maintenance window"
	expiresAt := time.Now().UTC().Add(24 * time.Hour)
	repo := &redeemRepoStub{placeholder
	svc := &RedeemService{redeemRepo: repoplaceholder

	result, err := svc.BatchUpdate(context.Background(), &RedeemCodeBatchUpdateInput{
		IDs: []int64{1, 2, 2placeholder,
		Fields: RedeemCodeBatchUpdateFields{
			Status:    &status,
			ExpiresAt: NullableTimeUpdate{Set: true, Value: &expiresAtplaceholder,
			Notes:     &notes,
	placeholder,
placeholder)

placeholder
	require.Equal(t, int64(2), result.Updated)
	require.True(t, repo.batchUpdateCalled)
	require.Equal(t, []int64{1, 2placeholder, repo.batchUpdateIDs)
	require.Equal(t, &status, repo.batchUpdateFields.Status)
	require.True(t, repo.batchUpdateFields.ExpiresAt.Set)
	require.WithinDuration(t, expiresAt, *repo.batchUpdateFields.ExpiresAt.Value, time.Second)
	require.Equal(t, &notes, repo.batchUpdateFields.Notes)
	require.False(t, repo.batchUpdateFields.GroupID.Set)
	require.Nil(t, repo.batchUpdateFields.Type)
	require.Nil(t, repo.batchUpdateFields.Value)
placeholder

func TestRedeemService_BatchUpdate_RejectsInvalidID(t *testing.T) {
	repo := &redeemRepoStub{placeholder
	svc := &RedeemService{redeemRepo: repoplaceholder
	notes := "bad id"

	result, err := svc.BatchUpdate(context.Background(), &RedeemCodeBatchUpdateInput{
		IDs:    []int64{1, 0placeholder,
		Fields: RedeemCodeBatchUpdateFields{Notes: &notesplaceholder,
placeholder)

	require.Nil(t, result)
placeholder
	require.True(t, infraerrors.IsBadRequest(err))
	require.False(t, repo.batchUpdateCalled)
placeholder

func TestRedeemService_BatchUpdate_RejectsCoreFieldsForUsedCodes(t *testing.T) {
	repo := &redeemRepoStub{placeholder
	svc := &RedeemService{redeemRepo: repoplaceholder
	newValue := 100.0

	result, err := svc.BatchUpdate(context.Background(), &RedeemCodeBatchUpdateInput{
		IDs: []int64{42placeholder,
		Fields: RedeemCodeBatchUpdateFields{
			Value: &newValue,
	placeholder,
placeholder)

	require.Nil(t, result)
placeholder
	require.True(t, infraerrors.IsBadRequest(err))
	require.False(t, repo.batchUpdateCalled)
placeholder
