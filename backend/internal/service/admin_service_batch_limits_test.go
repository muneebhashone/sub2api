//go:build unit

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type batchLimitsUserRepoStub struct {
	*userRepoStub
	calls       int
	userIDs     []int64
	concurrency *int
	rpmLimit    *int
	affected    int
	err         error
placeholder

func (s *batchLimitsUserRepoStub) BatchUpdateLimits(_ context.Context, userIDs []int64, concurrency, rpmLimit *int) (int, error) {
	s.calls++
	s.userIDs = append([]int64(nil), userIDs...)
	s.concurrency = cloneBatchLimitValue(concurrency)
	s.rpmLimit = cloneBatchLimitValue(rpmLimit)
	return s.affected, s.err
placeholder

func cloneBatchLimitValue(value *int) *int {
	if value == nil {
		return nil
placeholder
	cloned := *value
	return &cloned
placeholder

func TestAdminServiceBatchUpdateLimitsPassesOnlyProvidedFields(t *testing.T) {
	concurrency := 0
	repo := &batchLimitsUserRepoStub{
		userRepoStub: &userRepoStub{placeholder,
		affected:     2,
placeholder
	invalidator := &authCacheInvalidatorStub{placeholder
	service := &adminServiceImpl{userRepo: repo, authCacheInvalidator: invalidatorplaceholder

	affected, err := service.BatchUpdateLimits(
		context.Background(),
		[]int64{3, 0, 3, 7, -1placeholder,
		&concurrency,
		nil,
	)

placeholder
	require.Equal(t, 2, affected)
	require.Equal(t, []int64{3, 7placeholder, repo.userIDs)
	require.Equal(t, pointerToInt(0), repo.concurrency)
	require.Nil(t, repo.rpmLimit)
	require.Equal(t, []int64{3, 7placeholder, invalidator.userIDs)
placeholder

func TestAdminServiceBatchUpdateLimitsDoesNotInvalidateCacheOnRepositoryError(t *testing.T) {
	rpmLimit := 60
	repo := &batchLimitsUserRepoStub{
		userRepoStub: &userRepoStub{placeholder,
		err:          errors.New("database unavailable"),
placeholder
	invalidator := &authCacheInvalidatorStub{placeholder
	service := &adminServiceImpl{userRepo: repo, authCacheInvalidator: invalidatorplaceholder

	affected, err := service.BatchUpdateLimits(context.Background(), []int64{1, 2placeholder, nil, &rpmLimit)

	require.EqualError(t, err, "database unavailable")
	require.Zero(t, affected)
	require.Empty(t, invalidator.userIDs)
placeholder

func TestAdminServiceBatchUpdateLimitsRequiresAField(t *testing.T) {
	repo := &batchLimitsUserRepoStub{userRepoStub: &userRepoStub{placeholderplaceholder
	service := &adminServiceImpl{userRepo: repo, authCacheInvalidator: &authCacheInvalidatorStub{placeholderplaceholder

	affected, err := service.BatchUpdateLimits(context.Background(), []int64{1placeholder, nil, nil)

placeholder
	require.Zero(t, affected)
	require.Zero(t, repo.calls)
placeholder

func pointerToInt(value int) *int {
	return &value
placeholder
