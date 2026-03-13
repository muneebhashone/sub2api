//go:build unit

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// userGroupRateRepoStubForGroupRate implements UserGroupRateRepository for group rate tests.
type userGroupRateRepoStubForGroupRate struct {
	getByGroupIDData map[int64][]UserGroupRateEntry
	getByGroupIDErr  error

	deletedGroupIDs  []int64
	deleteByGroupErr error

	syncedGroupID int64
	syncedEntries []GroupRateMultiplierInput
	syncGroupErr  error
placeholder

func (s *userGroupRateRepoStubForGroupRate) GetByUserID(_ context.Context, _ int64) (map[int64]float64, error) {
	panic("unexpected GetByUserID call")
placeholder

func (s *userGroupRateRepoStubForGroupRate) GetByUserAndGroup(_ context.Context, _, _ int64) (*float64, error) {
	panic("unexpected GetByUserAndGroup call")
placeholder

func (s *userGroupRateRepoStubForGroupRate) GetByGroupID(_ context.Context, groupID int64) ([]UserGroupRateEntry, error) {
	if s.getByGroupIDErr != nil {
		return nil, s.getByGroupIDErr
placeholder
	return s.getByGroupIDData[groupID], nil
placeholder

func (s *userGroupRateRepoStubForGroupRate) SyncUserGroupRates(_ context.Context, _ int64, _ map[int64]*float64) error {
	panic("unexpected SyncUserGroupRates call")
placeholder

func (s *userGroupRateRepoStubForGroupRate) SyncGroupRateMultipliers(_ context.Context, groupID int64, entries []GroupRateMultiplierInput) error {
	s.syncedGroupID = groupID
	s.syncedEntries = entries
	return s.syncGroupErr
placeholder

func (s *userGroupRateRepoStubForGroupRate) DeleteByGroupID(_ context.Context, groupID int64) error {
	s.deletedGroupIDs = append(s.deletedGroupIDs, groupID)
	return s.deleteByGroupErr
placeholder

func (s *userGroupRateRepoStubForGroupRate) DeleteByUserID(_ context.Context, _ int64) error {
	panic("unexpected DeleteByUserID call")
placeholder

func TestAdminService_GetGroupRateMultipliers(t *testing.T) {
	t.Run("returns entries for group", func(t *testing.T) {
		repo := &userGroupRateRepoStubForGroupRate{
			getByGroupIDData: map[int64][]UserGroupRateEntry{
				10: {
					{UserID: 1, UserName: "alice", UserEmail: "alice@test.com", RateMultiplier: placeholder,
					{UserID: 2, UserName: "bob", UserEmail: "bob@test.com", RateMultiplier: 0.8placeholder,
			placeholder,
		placeholder,
	placeholder
		svc := &adminServiceImpl{userGroupRateRepo: repoplaceholder

		entries, err := svc.GetGroupRateMultipliers(context.Background(), 10)
	placeholder
		require.Len(t, entries, 2)
		require.Equal(t, int64(1), entries[0].UserID)
		require.Equal(t, "alice", entries[0].UserName)
		require.Equal(t, 1.5, entries[0].RateMultiplier)
		require.Equal(t, int64(2), entries[1].UserID)
		require.Equal(t, 0.8, entries[1].RateMultiplier)
placeholder)

	t.Run("returns nil when repo is nil", func(t *testing.T) {
		svc := &adminServiceImpl{userGroupRateRepo: nilplaceholder

		entries, err := svc.GetGroupRateMultipliers(context.Background(), 10)
	placeholder
		require.Nil(t, entries)
placeholder)

	t.Run("returns empty slice for group with no entries", func(t *testing.T) {
		repo := &userGroupRateRepoStubForGroupRate{
			getByGroupIDData: map[int64][]UserGroupRateEntry{placeholder,
	placeholder
		svc := &adminServiceImpl{userGroupRateRepo: repoplaceholder

		entries, err := svc.GetGroupRateMultipliers(context.Background(), 99)
	placeholder
		require.Nil(t, entries)
placeholder)

	t.Run("propagates repo error", func(t *testing.T) {
		repo := &userGroupRateRepoStubForGroupRate{
			getByGroupIDErr: errors.New("db error"),
	placeholder
		svc := &adminServiceImpl{userGroupRateRepo: repoplaceholder

		_, err := svc.GetGroupRateMultipliers(context.Background(), 10)
	placeholder
		require.Contains(t, err.Error(), "db error")
placeholder)
placeholder

func TestAdminService_ClearGroupRateMultipliers(t *testing.T) {
	t.Run("deletes by group ID", func(t *testing.T) {
		repo := &userGroupRateRepoStubForGroupRate{placeholder
		svc := &adminServiceImpl{userGroupRateRepo: repoplaceholder

		err := svc.ClearGroupRateMultipliers(context.Background(), 42)
	placeholder
		require.Equal(t, []int64{42placeholder, repo.deletedGroupIDs)
placeholder)

	t.Run("returns nil when repo is nil", func(t *testing.T) {
		svc := &adminServiceImpl{userGroupRateRepo: nilplaceholder

		err := svc.ClearGroupRateMultipliers(context.Background(), 42)
	placeholder
placeholder)

	t.Run("propagates repo error", func(t *testing.T) {
		repo := &userGroupRateRepoStubForGroupRate{
			deleteByGroupErr: errors.New("delete failed"),
	placeholder
		svc := &adminServiceImpl{userGroupRateRepo: repoplaceholder

		err := svc.ClearGroupRateMultipliers(context.Background(), 42)
	placeholder
		require.Contains(t, err.Error(), "delete failed")
placeholder)
placeholder

func TestAdminService_BatchSetGroupRateMultipliers(t *testing.T) {
	t.Run("syncs entries to repo", func(t *testing.T) {
		repo := &userGroupRateRepoStubForGroupRate{placeholder
		svc := &adminServiceImpl{userGroupRateRepo: repoplaceholder

		entries := []GroupRateMultiplierInput{
			{UserID: 1, RateMultiplier: placeholder,
			{UserID: 2, RateMultiplier: 0.8placeholder,
	placeholder
		err := svc.BatchSetGroupRateMultipliers(context.Background(), 10, entries)
	placeholder
		require.Equal(t, int64(10), repo.syncedGroupID)
		require.Equal(t, entries, repo.syncedEntries)
placeholder)

	t.Run("returns nil when repo is nil", func(t *testing.T) {
		svc := &adminServiceImpl{userGroupRateRepo: nilplaceholder

		err := svc.BatchSetGroupRateMultipliers(context.Background(), 10, nil)
	placeholder
placeholder)

	t.Run("propagates repo error", func(t *testing.T) {
		repo := &userGroupRateRepoStubForGroupRate{
			syncGroupErr: errors.New("sync failed"),
	placeholder
		svc := &adminServiceImpl{userGroupRateRepo: repoplaceholder

		err := svc.BatchSetGroupRateMultipliers(context.Background(), 10, []GroupRateMultiplierInput{
			{UserID: 1, RateMultiplier: 1.0placeholder,
	placeholder)
	placeholder
		require.Contains(t, err.Error(), "sync failed")
placeholder)
placeholder
