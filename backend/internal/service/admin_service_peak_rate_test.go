package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

type peakRateGroupRepoStub struct {
	getByID *Group
	updated *Group
placeholder

func (s *peakRateGroupRepoStub) Create(context.Context, *Group) error {
	panic("unexpected Create call")
placeholder

func (s *peakRateGroupRepoStub) GetByID(context.Context, int64) (*Group, error) {
	return s.getByID, nil
placeholder

func (s *peakRateGroupRepoStub) GetByIDLite(context.Context, int64) (*Group, error) {
	return s.getByID, nil
placeholder

func (s *peakRateGroupRepoStub) Update(_ context.Context, group *Group) error {
	s.updated = group
	return nil
placeholder

func (s *peakRateGroupRepoStub) Delete(context.Context, int64) error {
	panic("unexpected Delete call")
placeholder

func (s *peakRateGroupRepoStub) DeleteCascade(context.Context, int64) ([]int64, error) {
	panic("unexpected DeleteCascade call")
placeholder

func (s *peakRateGroupRepoStub) List(context.Context, pagination.PaginationParams) ([]Group, *pagination.PaginationResult, error) {
	panic("unexpected List call")
placeholder

func (s *peakRateGroupRepoStub) ListWithFilters(context.Context, pagination.PaginationParams, string, string, string, *bool) ([]Group, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
placeholder

func (s *peakRateGroupRepoStub) ListActive(context.Context) ([]Group, error) {
	panic("unexpected ListActive call")
placeholder

func (s *peakRateGroupRepoStub) ListActiveByPlatform(context.Context, string) ([]Group, error) {
	panic("unexpected ListActiveByPlatform call")
placeholder

func (s *peakRateGroupRepoStub) ExistsByName(context.Context, string) (bool, error) {
	panic("unexpected ExistsByName call")
placeholder

func (s *peakRateGroupRepoStub) GetAccountCount(context.Context, int64) (int64, int64, error) {
	panic("unexpected GetAccountCount call")
placeholder

func (s *peakRateGroupRepoStub) DeleteAccountGroupsByGroupID(context.Context, int64) (int64, error) {
	panic("unexpected DeleteAccountGroupsByGroupID call")
placeholder

func (s *peakRateGroupRepoStub) GetAccountIDsByGroupIDs(context.Context, []int64) ([]int64, error) {
	panic("unexpected GetAccountIDsByGroupIDs call")
placeholder

func (s *peakRateGroupRepoStub) BindAccountsToGroup(context.Context, int64, []int64) error {
	panic("unexpected BindAccountsToGroup call")
placeholder

func (s *peakRateGroupRepoStub) UpdateSortOrders(context.Context, []GroupSortOrderUpdate) error {
	panic("unexpected UpdateSortOrders call")
placeholder

func TestAdminService_UpdateGroup_ClearsPeakRateWhenChangingToStandardDefault(t *testing.T) {
	repo := &peakRateGroupRepoStub{getByID: &Group{
		ID:                 1,
		Name:               "existing-group",
		Platform:           PlatformOpenAI,
		Status:             StatusActive,
		SubscriptionType:   SubscriptionTypeSubscription,
		PeakRateEnabled:    true,
		PeakStart:          "14:00",
		PeakEnd:            "18:00",
		PeakRateMultiplier: 3,
placeholderplaceholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	group, err := svc.UpdateGroup(context.Background(), 1, &UpdateGroupInput{
		SubscriptionType: SubscriptionTypeStandard,
placeholder)
placeholder
	require.NotNil(t, group)
	require.NotNil(t, repo.updated)
	require.Equal(t, SubscriptionTypeStandard, repo.updated.SubscriptionType)
	require.False(t, repo.updated.PeakRateEnabled)
	require.Equal(t, "", repo.updated.PeakStart)
	require.Equal(t, "", repo.updated.PeakEnd)
	require.Equal(t, 1.0, repo.updated.PeakRateMultiplier)
placeholder
