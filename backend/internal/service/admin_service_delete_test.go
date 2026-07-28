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
	user           *User
	getErr         error
	createErr      error
	deleteErr      error
	exists         bool
	existsErr      error
	aliasExists    bool
	aliasErr       error
	guardedCreates int
	nextID         int64
	created        []*User
	updated        []*User
	deletedIDs     []int64
	usersByEmail   map[string]*User
	getByEmailErr  error
placeholder

func (s *userRepoStub) Create(ctx context.Context, user *User) error {
	if s.createErr != nil {
		return s.createErr
placeholder
	if s.nextID != 0 && user.ID == 0 {
		user.ID = s.nextID
placeholder
	s.created = append(s.created, user)
	if s.usersByEmail == nil {
		s.usersByEmail = make(map[string]*User)
placeholder
	s.usersByEmail[user.Email] = user
	s.user = user
	return nil
placeholder

func (s *userRepoStub) CreateWithEmailAliasGuard(ctx context.Context, user *User) error {
	s.guardedCreates++
	if s.aliasErr != nil {
		return s.aliasErr
placeholder
	if s.aliasExists {
		return ErrEmailExists
placeholder
	return s.Create(ctx, user)
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
	if s.getByEmailErr != nil {
		return nil, s.getByEmailErr
placeholder
	if s.usersByEmail != nil {
		if user, ok := s.usersByEmail[email]; ok {
			return user, nil
	placeholder
placeholder
	if s.user != nil && s.user.Email == email {
		return s.user, nil
placeholder
	return nil, ErrUserNotFound
placeholder

func (s *userRepoStub) GetFirstAdmin(ctx context.Context) (*User, error) {
	panic("unexpected GetFirstAdmin call")
placeholder

func (s *userRepoStub) Update(ctx context.Context, user *User, fields UserUpdateFields) error {
	s.updated = append(s.updated, user)
	if s.usersByEmail == nil {
		s.usersByEmail = make(map[string]*User)
placeholder
	s.usersByEmail[user.Email] = user
	s.user = user
	return nil
placeholder

func (s *userRepoStub) Delete(ctx context.Context, id int64) error {
	s.deletedIDs = append(s.deletedIDs, id)
	return s.deleteErr
placeholder

func (s *userRepoStub) GetUserAvatar(ctx context.Context, userID int64) (*UserAvatar, error) {
	panic("unexpected GetUserAvatar call")
placeholder

func (s *userRepoStub) UpsertUserAvatar(ctx context.Context, userID int64, input UpsertUserAvatarInput) (*UserAvatar, error) {
	panic("unexpected UpsertUserAvatar call")
placeholder

func (s *userRepoStub) DeleteUserAvatar(ctx context.Context, userID int64) error {
	panic("unexpected DeleteUserAvatar call")
placeholder

func (s *userRepoStub) List(ctx context.Context, params pagination.PaginationParams) ([]User, *pagination.PaginationResult, error) {
	panic("unexpected List call")
placeholder

func (s *userRepoStub) ListWithFilters(ctx context.Context, params pagination.PaginationParams, filters UserListFilters) ([]User, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
placeholder

func (s *userRepoStub) GetLatestUsedAtByUserIDs(ctx context.Context, userIDs []int64) (map[int64]*time.Time, error) {
	panic("unexpected GetLatestUsedAtByUserIDs call")
placeholder

func (s *userRepoStub) GetLatestUsedAtByUserID(ctx context.Context, userID int64) (*time.Time, error) {
	panic("unexpected GetLatestUsedAtByUserID call")
placeholder

func (s *userRepoStub) UpdateUserLastActiveAt(ctx context.Context, userID int64, activeAt time.Time) error {
	panic("unexpected UpdateUserLastActiveAt call")
placeholder

func (s *userRepoStub) UpdateBalance(ctx context.Context, id int64, amount float64) error {
	panic("unexpected UpdateBalance call")
placeholder

func (s *userRepoStub) DeductBalance(ctx context.Context, id int64, amount float64) error {
	panic("unexpected DeductBalance call")
placeholder

func (s *userRepoStub) AdjustBalance(ctx context.Context, id int64, delta float64) (BalanceChange, error) {
	panic("unexpected AdjustBalance call")
placeholder

func (s *userRepoStub) SetBalance(ctx context.Context, id int64, value float64) (BalanceChange, error) {
	panic("unexpected SetBalance call")
placeholder

func (s *userRepoStub) UpdateConcurrency(ctx context.Context, id int64, amount int) error {
	panic("unexpected UpdateConcurrency call")
placeholder

func (s *userRepoStub) BatchSetConcurrency(context.Context, []int64, int) (int, error) { return 0, nil placeholder
func (s *userRepoStub) BatchAddConcurrency(context.Context, []int64, int) (int, error) { return 0, nil placeholder
func (s *userRepoStub) BatchUpdateLimits(context.Context, []int64, *int, *int) (int, error) {
	return 0, nil
placeholder

func (s *userRepoStub) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if s.existsErr != nil {
		return false, s.existsErr
placeholder
	return s.exists, nil
placeholder

func (s *userRepoStub) ExistsByEmailAlias(ctx context.Context, email string) (bool, error) {
	if s.aliasErr != nil {
		return false, s.aliasErr
placeholder
	return s.aliasExists, nil
placeholder

func (s *userRepoStub) RemoveGroupFromAllowedGroups(ctx context.Context, groupID int64) (int64, error) {
	panic("unexpected RemoveGroupFromAllowedGroups call")
placeholder

func (s *userRepoStub) RemoveGroupFromUserAllowedGroups(ctx context.Context, userID int64, groupID int64) error {
	panic("unexpected RemoveGroupFromUserAllowedGroups call")
placeholder

func (s *userRepoStub) AddGroupToAllowedGroups(ctx context.Context, userID int64, groupID int64) error {
	panic("unexpected AddGroupToAllowedGroups call")
placeholder

func (s *userRepoStub) ListUserAuthIdentities(ctx context.Context, userID int64) ([]UserAuthIdentityRecord, error) {
	panic("unexpected ListUserAuthIdentities call")
placeholder

func (s *userRepoStub) UnbindUserAuthProvider(context.Context, int64, string) error {
	panic("unexpected UnbindUserAuthProvider call")
placeholder

func (s *userRepoStub) UpdateTotpSecret(ctx context.Context, userID int64, encryptedSecret *string) error {
	panic("unexpected UpdateTotpSecret call")
placeholder

func (s *userRepoStub) EnableTotp(ctx context.Context, userID int64) error {
	panic("unexpected EnableTotp call")
placeholder

func (s *userRepoStub) DisableTotp(ctx context.Context, userID int64) error {
	panic("unexpected DisableTotp call")
placeholder

func (s *userRepoStub) GetByIDIncludeDeleted(ctx context.Context, id int64) (*User, error) {
	return s.GetByID(ctx, id)
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

func (s *groupRepoStub) GetByIDLite(ctx context.Context, id int64) (*Group, error) {
	panic("unexpected GetByIDLite call")
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

func (s *groupRepoStub) ListWithFilters(ctx context.Context, params pagination.PaginationParams, platform, status, search string, isExclusive *bool) ([]Group, *pagination.PaginationResult, error) {
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

func (s *groupRepoStub) GetAccountCount(ctx context.Context, groupID int64) (int64, int64, error) {
	panic("unexpected GetAccountCount call")
placeholder

func (s *groupRepoStub) DeleteAccountGroupsByGroupID(ctx context.Context, groupID int64) (int64, error) {
	panic("unexpected DeleteAccountGroupsByGroupID call")
placeholder

func (s *groupRepoStub) BindAccountsToGroup(ctx context.Context, groupID int64, accountIDs []int64) error {
	panic("unexpected BindAccountsToGroup call")
placeholder

func (s *groupRepoStub) GetAccountIDsByGroupIDs(ctx context.Context, groupIDs []int64) ([]int64, error) {
	panic("unexpected GetAccountIDsByGroupIDs call")
placeholder

func (s *groupRepoStub) UpdateSortOrders(ctx context.Context, updates []GroupSortOrderUpdate) error {
	return nil
placeholder

type deleteGroupAPIKeyRepoStub struct {
	apiKeyRepoStubForGroupUpdate
	keys         []string
	listErr      error
	listGroupIDs []int64
placeholder

func (s *deleteGroupAPIKeyRepoStub) ListKeysByGroupID(ctx context.Context, groupID int64) ([]string, error) {
	s.listGroupIDs = append(s.listGroupIDs, groupID)
	if s.listErr != nil {
		return nil, s.listErr
placeholder
	return s.keys, nil
placeholder

type proxyRepoStub struct {
	deleteErr    error
	countErr     error
	accountCount int64
	deletedIDs   []int64
placeholder

func (s *proxyRepoStub) Create(ctx context.Context, proxy *Proxy) error {
	panic("unexpected Create call")
placeholder

func (s *proxyRepoStub) GetByID(ctx context.Context, id int64) (*Proxy, error) {
	panic("unexpected GetByID call")
placeholder

func (s *proxyRepoStub) ListByIDs(ctx context.Context, ids []int64) ([]Proxy, error) {
	panic("unexpected ListByIDs call")
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

func (s *proxyRepoStub) ListWithFiltersAndAccountCount(ctx context.Context, params pagination.PaginationParams, protocol, status, search string) ([]ProxyWithAccountCount, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFiltersAndAccountCount call")
placeholder

func (s *proxyRepoStub) ExistsByHostPortAuth(ctx context.Context, host string, port int, username, password string) (bool, error) {
	panic("unexpected ExistsByHostPortAuth call")
placeholder

func (s *proxyRepoStub) CountAccountsByProxyID(ctx context.Context, proxyID int64) (int64, error) {
	if s.countErr != nil {
		return 0, s.countErr
placeholder
	return s.accountCount, nil
placeholder

func (s *proxyRepoStub) ListAccountSummariesByProxyID(ctx context.Context, proxyID int64) ([]ProxyAccountSummary, error) {
	panic("unexpected ListAccountSummariesByProxyID call")
placeholder
func (s *proxyRepoStub) SweepExpiredProxies(_ context.Context, _ time.Time) (int64, error) {
	return 0, nil
placeholder
func (s *proxyRepoStub) ListAllForFallback(_ context.Context) ([]Proxy, error) {
	return nil, nil
placeholder
func (s *proxyRepoStub) CountExpired(_ context.Context) (int64, error) {
	return 0, nil
placeholder
func (s *proxyRepoStub) CountExpiringSoon(_ context.Context, _ time.Time) (int64, error) {
	return 0, nil
placeholder

type redeemRepoStub struct {
	deleteErrByID map[int64]error
	deletedIDs    []int64

	batchUpdateIDs    []int64
	batchUpdateFields RedeemCodeBatchUpdateFields
	batchUpdateResult int64
	batchUpdateErr    error
	batchUpdateCalled bool
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

func (s *redeemRepoStub) BatchUpdate(ctx context.Context, ids []int64, fields RedeemCodeBatchUpdateFields) (int64, error) {
	s.batchUpdateCalled = true
	s.batchUpdateIDs = append([]int64(nil), ids...)
	s.batchUpdateFields = fields
	if s.batchUpdateErr != nil {
		return 0, s.batchUpdateErr
placeholder
	if s.batchUpdateResult != 0 {
		return s.batchUpdateResult, nil
placeholder
	return int64(len(ids)), nil
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

func (s *redeemRepoStub) ListByUserPaginated(ctx context.Context, userID int64, params pagination.PaginationParams, codeType string) ([]RedeemCode, *pagination.PaginationResult, error) {
	panic("unexpected ListByUserPaginated call")
placeholder

func (s *redeemRepoStub) SumPositiveBalanceByUser(ctx context.Context, userID int64) (float64, error) {
	panic("unexpected SumPositiveBalanceByUser call")
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

func (s *billingCacheStub) GetAPIKeyRateLimit(ctx context.Context, keyID int64) (*APIKeyRateLimitCacheData, error) {
	panic("unexpected GetAPIKeyRateLimit call")
placeholder
func (s *billingCacheStub) SetAPIKeyRateLimit(ctx context.Context, keyID int64, data *APIKeyRateLimitCacheData) error {
	panic("unexpected SetAPIKeyRateLimit call")
placeholder
func (s *billingCacheStub) UpdateAPIKeyRateLimitUsage(ctx context.Context, keyID int64, cost float64) error {
	panic("unexpected UpdateAPIKeyRateLimitUsage call")
placeholder
func (s *billingCacheStub) InvalidateAPIKeyRateLimit(ctx context.Context, keyID int64) error {
	panic("unexpected InvalidateAPIKeyRateLimit call")
placeholder

func (s *billingCacheStub) GetUserPlatformQuotaCache(ctx context.Context, userID int64, platform string) (*UserPlatformQuotaCacheEntry, bool, error) {
	panic("unexpected GetUserPlatformQuotaCache call")
placeholder

func (s *billingCacheStub) SetUserPlatformQuotaCache(ctx context.Context, userID int64, platform string, entry *UserPlatformQuotaCacheEntry, ttl time.Duration) error {
	panic("unexpected SetUserPlatformQuotaCache call")
placeholder

func (s *billingCacheStub) DeleteUserPlatformQuotaCache(ctx context.Context, userID int64, platform string) error {
	panic("unexpected DeleteUserPlatformQuotaCache call")
placeholder

func (s *billingCacheStub) IncrUserPlatformQuotaUsageCache(ctx context.Context, userID int64, platform string, cost float64, ttl time.Duration, markDirty bool) error {
	panic("unexpected IncrUserPlatformQuotaUsageCache call")
placeholder

func (s *billingCacheStub) PopDirtyUserPlatformQuotaKeys(ctx context.Context, n int) ([]UserPlatformQuotaKey, error) {
	panic("unexpected PopDirtyUserPlatformQuotaKeys call")
placeholder

func (s *billingCacheStub) ReaddDirtyUserPlatformQuotaKeys(ctx context.Context, keys []UserPlatformQuotaKey) error {
	panic("unexpected ReaddDirtyUserPlatformQuotaKeys call")
placeholder

func (s *billingCacheStub) BatchGetUserPlatformQuotaCache(ctx context.Context, keys []UserPlatformQuotaKey) ([]*UserPlatformQuotaCacheEntry, error) {
	panic("unexpected BatchGetUserPlatformQuotaCache call")
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

func TestAdminService_DeleteUser_DeletesOwnedAPIKeys(t *testing.T) {
	repo := &userRepoStub{user: &User{ID: 7, Role: RoleUserplaceholderplaceholder
	apiKeyRepo := &apiKeyRepoStub{
		allowListByUserID: true,
		listByUserIDKeys: []APIKey{
			{ID: 11, UserID: 7, Key: "sk-user-1"placeholder,
			{ID: 12, UserID: 7, Key: "sk-user-2"placeholder,
	placeholder,
placeholder
	invalidator := &authCacheInvalidatorStub{placeholder
	svc := &adminServiceImpl{
		userRepo:             repo,
		apiKeyRepo:           apiKeyRepo,
		authCacheInvalidator: invalidator,
placeholder

	err := svc.DeleteUser(context.Background(), 7)
placeholder
	require.Equal(t, []int64{7placeholder, repo.deletedIDs)
	require.Equal(t, []int64{7placeholder, apiKeyRepo.listByUserIDCalls)
	require.Equal(t, []int64{11, 12placeholder, apiKeyRepo.deletedIDs)
	require.ElementsMatch(t, []string{"sk-user-1", "sk-user-2"placeholder, invalidator.keys)
	require.Equal(t, []int64{7placeholder, invalidator.userIDs)
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

func TestAdminService_DeleteGroup_InvalidatesAuthCacheForBoundKeys(t *testing.T) {
	repo := &groupRepoStub{placeholder
	apiKeyRepo := &deleteGroupAPIKeyRepoStub{keys: []string{"k1", "k2"placeholderplaceholder
	invalidator := &authCacheInvalidatorStub{placeholder
	svc := &adminServiceImpl{
		groupRepo:            repo,
		apiKeyRepo:           apiKeyRepo,
		authCacheInvalidator: invalidator,
placeholder

	err := svc.DeleteGroup(context.Background(), 5)
placeholder
	require.Equal(t, []int64{5placeholder, repo.deleteCalls)
	require.Equal(t, []int64{5placeholder, apiKeyRepo.listGroupIDs)
	require.Equal(t, []string{"k1", "k2"placeholder, invalidator.keys)
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

func TestAdminService_DeleteProxy_InUse(t *testing.T) {
	repo := &proxyRepoStub{accountCount: 2placeholder
	svc := &adminServiceImpl{proxyRepo: repoplaceholder

	err := svc.DeleteProxy(context.Background(), 77)
	require.ErrorIs(t, err, ErrProxyInUse)
	require.Empty(t, repo.deletedIDs)
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
