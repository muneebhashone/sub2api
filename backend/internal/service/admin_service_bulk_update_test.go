//go:build unit

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type accountRepoStubForBulkUpdate struct {
	accountRepoStub
	bulkUpdateErr     error
	bulkUpdateIDs     []int64
	bindGroupErrByID  map[int64]error
placeholder

func (s *accountRepoStubForBulkUpdate) BulkUpdate(_ context.Context, ids []int64, _ AccountBulkUpdate) (int64, error) {
	s.bulkUpdateIDs = append([]int64{placeholder, ids...)
	if s.bulkUpdateErr != nil {
		return 0, s.bulkUpdateErr
placeholder
	return int64(len(ids)), nil
placeholder

func (s *accountRepoStubForBulkUpdate) BindGroups(_ context.Context, accountID int64, _ []int64) error {
	if err, ok := s.bindGroupErrByID[accountID]; ok {
		return err
placeholder
	return nil
placeholder

// TestAdminService_BulkUpdateAccounts_AllSuccessIDs 验证批量更新成功时返回 success_ids/failed_ids。
func TestAdminService_BulkUpdateAccounts_AllSuccessIDs(t *testing.T) {
	repo := &accountRepoStubForBulkUpdate{placeholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder

	schedulable := true
	input := &BulkUpdateAccountsInput{
		AccountIDs:  []int64{1, 2, 3placeholder,
		Schedulable: &schedulable,
placeholder

	result, err := svc.BulkUpdateAccounts(context.Background(), input)
placeholder
	require.Equal(t, 3, result.Success)
	require.Equal(t, 0, result.Failed)
	require.ElementsMatch(t, []int64{1, 2, 3placeholder, result.SuccessIDs)
	require.Empty(t, result.FailedIDs)
	require.Len(t, result.Results, 3)
placeholder

// TestAdminService_BulkUpdateAccounts_PartialFailureIDs 验证部分失败时 success_ids/failed_ids 正确。
func TestAdminService_BulkUpdateAccounts_PartialFailureIDs(t *testing.T) {
	repo := &accountRepoStubForBulkUpdate{
		bindGroupErrByID: map[int64]error{
			2: errors.New("bind failed"),
	placeholder,
placeholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder

	groupIDs := []int64{10placeholder
	schedulable := false
	input := &BulkUpdateAccountsInput{
		AccountIDs:            []int64{1, 2, 3placeholder,
		GroupIDs:              &groupIDs,
		Schedulable:           &schedulable,
		SkipMixedChannelCheck: true,
placeholder

	result, err := svc.BulkUpdateAccounts(context.Background(), input)
placeholder
	require.Equal(t, 2, result.Success)
	require.Equal(t, 1, result.Failed)
	require.ElementsMatch(t, []int64{1, 3placeholder, result.SuccessIDs)
	require.ElementsMatch(t, []int64{2placeholder, result.FailedIDs)
	require.Len(t, result.Results, 3)
placeholder
