//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

type userRepoStub struct {
	user       *User
	getErr     error
	deleteErr  error
	deletedIDs []int64
placeholder

func (s *userRepoStub) Create(ctx context.Context, user *User) error {
	panic("unexpected Create call")
placeholder

func (s *userRepoStub) GetByID(ctx context.Context, id int64) (*User, error) {
	if s.getErr != nil {
		return nil, s.getErr
placeholder
	if s.user == nil {
		return nil, ErrUserNotFound
placeholder
	return s.user, nil
placeholder

func (s *userRepoStub) GetByEmail(ctx context.Context, email string) (*User, error) {
	panic("unexpected GetByEmail call")
placeholder

func (s *userRepoStub) GetFirstAdmin(ctx context.Context) (*User, error) {
	panic("unexpected GetFirstAdmin call")
placeholder

func (s *userRepoStub) Update(ctx context.Context, user *User) error {
	panic("unexpected Update call")
placeholder

func (s *userRepoStub) Delete(ctx context.Context, id int64) error {
	s.deletedIDs = append(s.deletedIDs, id)
	return s.deleteErr
placeholder

func (s *userRepoStub) List(ctx context.Context, params pagination.PaginationParams) ([]User, *pagination.PaginationResult, error) {
	panic("unexpected List call")
placeholder

func (s *userRepoStub) ListWithFilters(ctx context.Context, params pagination.PaginationParams, status, role, search string) ([]User, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
placeholder

func (s *userRepoStub) UpdateBalance(ctx context.Context, id int64, amount float64) error {
	panic("unexpected UpdateBalance call")
placeholder

func (s *userRepoStub) DeductBalance(ctx context.Context, id int64, amount float64) error {
	panic("unexpected DeductBalance call")
placeholder

func (s *userRepoStub) UpdateConcurrency(ctx context.Context, id int64, amount int) error {
	panic("unexpected UpdateConcurrency call")
placeholder

func (s *userRepoStub) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	panic("unexpected ExistsByEmail call")
placeholder

func (s *userRepoStub) RemoveGroupFromAllowedGroups(ctx context.Context, groupID int64) (int64, error) {
	panic("unexpected RemoveGroupFromAllowedGroups call")
placeholder

type groupRepoStub struct {
	affectedUserIDs []int64
	deleteErr       error
	deleteCalls     []int64
placeholder

func (s *groupRepoStub) Create(ctx context.Context, group *Group) error {
	panic("unexpected Create call")
placeholder

func (s *groupRepoStub) GetByID(ctx context.Context, id int64) (*Group, error) {
	panic("unexpected GetByID call")
placeholder

func (s *groupRepoStub) Update(ctx context.Context, group *Group) error {
	panic("unexpected Update call")
placeholder

func (s *groupRepoStub) Delete(ctx context.Context, id int64) error {
	panic("unexpected Delete call")
placeholder

func (s *groupRepoStub) DeleteCascade(ctx context.Context, id int64) ([]int64, error) {
	s.deleteCalls = append(s.deleteCalls, id)
	return s.affectedUserIDs, s.deleteErr
placeholder

func (s *groupRepoStub) List(ctx context.Context, params pagination.PaginationParams) ([]Group, *pagination.PaginationResult, error) {
	panic("unexpected List call")
placeholder

func (s *groupRepoStub) ListWithFilters(ctx context.Context, params pagination.PaginationParams, platform, status string, isExclusive *bool) ([]Group, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
placeholder

func (s *groupRepoStub) ListActive(ctx context.Context) ([]Group, error) {
	panic("unexpected ListActive call")
placeholder

func (s *groupRepoStub) ListActiveByPlatform(ctx context.Context, platform string) ([]Group, error) {
	panic("unexpected ListActiveByPlatform call")
placeholder

func (s *groupRepoStub) ExistsByName(ctx context.Context, name string) (bool, error) {
	panic("unexpected ExistsByName call")
placeholder

func (s *groupRepoStub) GetAccountCount(ctx context.Context, groupID int64) (int64, error) {
	panic("unexpected GetAccountCount call")
placeholder

func (s *groupRepoStub) DeleteAccountGroupsByGroupID(ctx context.Context, groupID int64) (int64, error) {
	panic("unexpected DeleteAccountGroupsByGroupID call")
placeholder

type proxyRepoStub struct {
	deleteErr  error
	deletedIDs []int64
placeholder

func (s *proxyRepoStub) Create(ctx context.Context, proxy *Proxy) error {
	panic("unexpected Create call")
placeholder

func (s *proxyRepoStub) GetByID(ctx context.Context, id int64) (*Proxy, error) {
	panic("unexpected GetByID call")
placeholder

func (s *proxyRepoStub) Update(ctx context.Context, proxy *Proxy) error {
	panic("unexpected Update call")
placeholder

func (s *proxyRepoStub) Delete(ctx context.Context, id int64) error {
	s.deletedIDs = append(s.deletedIDs, id)
	return s.deleteErr
placeholder

func (s *proxyRepoStub) List(ctx context.Context, params pagination.PaginationParams) ([]Proxy, *pagination.PaginationResult, error) {
	panic("unexpected List call")
placeholder

func (s *proxyRepoStub) ListWithFilters(ctx context.Context, params pagination.PaginationParams, protocol, status, search string) ([]Proxy, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
placeholder

func (s *proxyRepoStub) ListActive(ctx context.Context) ([]Proxy, error) {
	panic("unexpected ListActive call")
placeholder

func (s *proxyRepoStub) ListActiveWithAccountCount(ctx context.Context) ([]ProxyWithAccountCount, error) {
	panic("unexpected ListActiveWithAccountCount call")
placeholder

func (s *proxyRepoStub) ExistsByHostPortAuth(ctx context.Context, host string, port int, username, password string) (bool, error) {
	panic("unexpected ExistsByHostPortAuth call")
placeholder

func (s *proxyRepoStub) CountAccountsByProxyID(ctx context.Context, proxyID int64) (int64, error) {
	panic("unexpected CountAccountsByProxyID call")
placeholder

type redeemRepoStub struct {
	deleteErrByID map[int64]error
	deletedIDs    []int64
placeholder

func (s *redeemRepoStub) Create(ctx context.Context, code *RedeemCode) error {
	panic("unexpected Create call")
placeholder

func (s *redeemRepoStub) CreateBatch(ctx context.Context, codes []RedeemCode) error {
	panic("unexpected CreateBatch call")
placeholder

func (s *redeemRepoStub) GetByID(ctx context.Context, id int64) (*RedeemCode, error) {
	panic("unexpected GetByID call")
placeholder

func (s *redeemRepoStub) GetByCode(ctx context.Context, code string) (*RedeemCode, error) {
	panic("unexpected GetByCode call")
placeholder

func (s *redeemRepoStub) Update(ctx context.Context, code *RedeemCode) error {
	panic("unexpected Update call")
placeholder

func (s *redeemRepoStub) Delete(ctx context.Context, id int64) error {
	s.deletedIDs = append(s.deletedIDs, id)
	if s.deleteErrByID != nil {
		if err, ok := s.deleteErrByID[id]; ok {
			return err
	placeholder
placeholder
	return nil
placeholder

func (s *redeemRepoStub) Use(ctx context.Context, id, userID int64) error {
	panic("unexpected Use call")
placeholder

func (s *redeemRepoStub) List(ctx context.Context, params pagination.PaginationParams) ([]RedeemCode, *pagination.PaginationResult, error) {
	panic("unexpected List call")
placeholder

func (s *redeemRepoStub) ListWithFilters(ctx context.Context, params pagination.PaginationParams, codeType, status, search string) ([]RedeemCode, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
placeholder

func (s *redeemRepoStub) ListByUser(ctx context.Context, userID int64, limit int) ([]RedeemCode, error) {
	panic("unexpected ListByUser call")
placeholder

type subscriptionInvalidateCall struct {
	userID  int64
	groupID int64
placeholder

type billingCacheStub struct {
	invalidations chan subscriptionInvalidateCall
placeholder

func newBillingCacheStub(buffer int) *billingCacheStub {
	return &billingCacheStub{invalidations: make(chan subscriptionInvalidateCall, buffer)placeholder
placeholder

func (s *billingCacheStub) GetUserBalance(ctx context.Context, userID int64) (float64, error) {
	panic("unexpected GetUserBalance call")
placeholder

func (s *billingCacheStub) SetUserBalance(ctx context.Context, userID int64, balance float64) error {
	panic("unexpected SetUserBalance call")
placeholder

func (s *billingCacheStub) DeductUserBalance(ctx context.Context, userID int64, amount float64) error {
	panic("unexpected DeductUserBalance call")
placeholder

func (s *billingCacheStub) InvalidateUserBalance(ctx context.Context, userID int64) error {
	panic("unexpected InvalidateUserBalance call")
placeholder

func (s *billingCacheStub) GetSubscriptionCache(ctx context.Context, userID, groupID int64) (*SubscriptionCacheData, error) {
	panic("unexpected GetSubscriptionCache call")
placeholder

func (s *billingCacheStub) SetSubscriptionCache(ctx context.Context, userID, groupID int64, data *SubscriptionCacheData) error {
	panic("unexpected SetSubscriptionCache call")
placeholder

func (s *billingCacheStub) UpdateSubscriptionUsage(ctx context.Context, userID, groupID int64, cost float64) error {
	panic("unexpected UpdateSubscriptionUsage call")
placeholder

func (s *billingCacheStub) InvalidateSubscriptionCache(ctx context.Context, userID, groupID int64) error {
	s.invalidations <- subscriptionInvalidateCall{userID: userID, groupID: groupIDplaceholder
	return nil
placeholder

func waitForInvalidations(t *testing.T, ch <-chan subscriptionInvalidateCall, expected int) []subscriptionInvalidateCall {
placeholder
	calls := make([]subscriptionInvalidateCall, 0, expected)
	timeout := time.After(2 * time.Second)
	for len(calls) < expected {
		select {
		case call := <-ch:
			calls = append(calls, call)
		case <-timeout:
			t.Fatalf("timeout waiting for %d invalidations, got %d", expected, len(calls))
	placeholder
placeholder
	return calls
placeholder

func TestAdminService_DeleteUser_Success(t *testing.T) {
	repo := &userRepoStub{user: &User{ID: 7, Role: RoleUserplaceholderplaceholder
	svc := &adminServiceImpl{userRepo: repoplaceholder

	err := svc.DeleteUser(context.Background(), 7)
placeholder
	require.Equal(t, []int64{7placeholder, repo.deletedIDs)
placeholder

func TestAdminService_DeleteUser_NotFound(t *testing.T) {
	repo := &userRepoStub{getErr: ErrUserNotFoundplaceholder
	svc := &adminServiceImpl{userRepo: repoplaceholder

	err := svc.DeleteUser(context.Background(), 404)
	require.ErrorIs(t, err, ErrUserNotFound)
	require.Empty(t, repo.deletedIDs)
placeholder

func TestAdminService_DeleteUser_AdminGuard(t *testing.T) {
	repo := &userRepoStub{user: &User{ID: 1, Role: RoleAdminplaceholderplaceholder
	svc := &adminServiceImpl{userRepo: repoplaceholder

	err := svc.DeleteUser(context.Background(), 1)
placeholder
	require.ErrorContains(t, err, "cannot delete admin user")
	require.Empty(t, repo.deletedIDs)
placeholder

func TestAdminService_DeleteUser_DeleteError(t *testing.T) {
	deleteErr := errors.New("delete failed")
	repo := &userRepoStub{
		user:      &User{ID: 9, Role: RoleUserplaceholder,
		deleteErr: deleteErr,
placeholder
	svc := &adminServiceImpl{userRepo: repoplaceholder

	err := svc.DeleteUser(context.Background(), 9)
	require.ErrorIs(t, err, deleteErr)
	require.Equal(t, []int64{9placeholder, repo.deletedIDs)
placeholder

func TestAdminService_DeleteGroup_Success_WithCacheInvalidation(t *testing.T) {
	cache := newBillingCacheStub(2)
	repo := &groupRepoStub{affectedUserIDs: []int64{11, 12placeholderplaceholder
	svc := &adminServiceImpl{
		groupRepo:           repo,
		billingCacheService: &BillingCacheService{cache: cacheplaceholder,
placeholder

	err := svc.DeleteGroup(context.Background(), 5)
placeholder
	require.Equal(t, []int64{5placeholder, repo.deleteCalls)

	calls := waitForInvalidations(t, cache.invalidations, 2)
	require.ElementsMatch(t, []subscriptionInvalidateCall{
		{userID: 11, groupID: 5placeholder,
		{userID: 12, groupID: 5placeholder,
placeholder, calls)
placeholder

func TestAdminService_DeleteGroup_NotFound(t *testing.T) {
	repo := &groupRepoStub{deleteErr: ErrGroupNotFoundplaceholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	err := svc.DeleteGroup(context.Background(), 99)
	require.ErrorIs(t, err, ErrGroupNotFound)
placeholder

func TestAdminService_DeleteGroup_Error(t *testing.T) {
	deleteErr := errors.New("delete failed")
	repo := &groupRepoStub{deleteErr: deleteErrplaceholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	err := svc.DeleteGroup(context.Background(), 42)
	require.ErrorIs(t, err, deleteErr)
placeholder

func TestAdminService_DeleteProxy_Success(t *testing.T) {
	repo := &proxyRepoStub{placeholder
	svc := &adminServiceImpl{proxyRepo: repoplaceholder

	err := svc.DeleteProxy(context.Background(), 7)
placeholder
	require.Equal(t, []int64{7placeholder, repo.deletedIDs)
placeholder

func TestAdminService_DeleteProxy_Idempotent(t *testing.T) {
	repo := &proxyRepoStub{placeholder
	svc := &adminServiceImpl{proxyRepo: repoplaceholder

	err := svc.DeleteProxy(context.Background(), 404)
placeholder
	require.Equal(t, []int64{404placeholder, repo.deletedIDs)
placeholder

func TestAdminService_DeleteProxy_Error(t *testing.T) {
	deleteErr := errors.New("delete failed")
	repo := &proxyRepoStub{deleteErr: deleteErrplaceholder
	svc := &adminServiceImpl{proxyRepo: repoplaceholder

	err := svc.DeleteProxy(context.Background(), 33)
	require.ErrorIs(t, err, deleteErr)
placeholder

func TestAdminService_DeleteRedeemCode_Success(t *testing.T) {
	repo := &redeemRepoStub{placeholder
	svc := &adminServiceImpl{redeemCodeRepo: repoplaceholder

	err := svc.DeleteRedeemCode(context.Background(), 10)
placeholder
	require.Equal(t, []int64{10placeholder, repo.deletedIDs)
placeholder

func TestAdminService_DeleteRedeemCode_Idempotent(t *testing.T) {
	repo := &redeemRepoStub{placeholder
	svc := &adminServiceImpl{redeemCodeRepo: repoplaceholder

	err := svc.DeleteRedeemCode(context.Background(), 999)
placeholder
	require.Equal(t, []int64{999placeholder, repo.deletedIDs)
placeholder

func TestAdminService_DeleteRedeemCode_Error(t *testing.T) {
	deleteErr := errors.New("delete failed")
	repo := &redeemRepoStub{deleteErrByID: map[int64]error{1: deleteErrplaceholderplaceholder
	svc := &adminServiceImpl{redeemCodeRepo: repoplaceholder

	err := svc.DeleteRedeemCode(context.Background(), 1)
	require.ErrorIs(t, err, deleteErr)
	require.Equal(t, []int64{1placeholder, repo.deletedIDs)
placeholder

func TestAdminService_BatchDeleteRedeemCodes_Success(t *testing.T) {
	repo := &redeemRepoStub{placeholder
	svc := &adminServiceImpl{redeemCodeRepo: repoplaceholder

	deleted, err := svc.BatchDeleteRedeemCodes(context.Background(), []int64{1, 2, 3placeholder)
placeholder
	require.Equal(t, int64(3), deleted)
	require.Equal(t, []int64{1, 2, 3placeholder, repo.deletedIDs)
placeholder

func TestAdminService_BatchDeleteRedeemCodes_PartialFailures(t *testing.T) {
	repo := &redeemRepoStub{
		deleteErrByID: map[int64]error{
			2: errors.New("db error"),
	placeholder,
placeholder
	svc := &adminServiceImpl{redeemCodeRepo: repoplaceholder

	deleted, err := svc.BatchDeleteRedeemCodes(context.Background(), []int64{1, 2, 3placeholder)
placeholder
	require.Equal(t, int64(2), deleted)
	require.Equal(t, []int64{1, 2, 3placeholder, repo.deletedIDs)
placeholder
