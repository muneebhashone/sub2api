//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Stubs
// ---------------------------------------------------------------------------

// apiKeyRepoStubForGroupUpdate implements APIKeyRepository for AdminUpdateAPIKeyGroupID tests.
type apiKeyRepoStubForGroupUpdate struct {
	key       *APIKey
	getErr    error
	updateErr error
	updated   *APIKey // captures what was passed to Update
placeholder

func (s *apiKeyRepoStubForGroupUpdate) GetByID(_ context.Context, _ int64) (*APIKey, error) {
	if s.getErr != nil {
		return nil, s.getErr
placeholder
	clone := *s.key
	return &clone, nil
placeholder
func (s *apiKeyRepoStubForGroupUpdate) Update(_ context.Context, key *APIKey) error {
	if s.updateErr != nil {
		return s.updateErr
placeholder
	clone := *key
	s.updated = &clone
	return nil
placeholder

// Unused methods – panic on unexpected call.
func (s *apiKeyRepoStubForGroupUpdate) Create(context.Context, *APIKey) error { panic("unexpected") placeholder
func (s *apiKeyRepoStubForGroupUpdate) GetKeyAndOwnerID(context.Context, int64) (string, int64, error) {
	panic("unexpected")
placeholder
func (s *apiKeyRepoStubForGroupUpdate) GetByKey(context.Context, string) (*APIKey, error) {
	panic("unexpected")
placeholder
func (s *apiKeyRepoStubForGroupUpdate) GetByKeyForAuth(context.Context, string) (*APIKey, error) {
	panic("unexpected")
placeholder
func (s *apiKeyRepoStubForGroupUpdate) Delete(context.Context, int64) error { panic("unexpected") placeholder
func (s *apiKeyRepoStubForGroupUpdate) ListByUserID(context.Context, int64, pagination.PaginationParams) ([]APIKey, *pagination.PaginationResult, error) {
	panic("unexpected")
placeholder
func (s *apiKeyRepoStubForGroupUpdate) VerifyOwnership(context.Context, int64, []int64) ([]int64, error) {
	panic("unexpected")
placeholder
func (s *apiKeyRepoStubForGroupUpdate) CountByUserID(context.Context, int64) (int64, error) {
	panic("unexpected")
placeholder
func (s *apiKeyRepoStubForGroupUpdate) ExistsByKey(context.Context, string) (bool, error) {
	panic("unexpected")
placeholder
func (s *apiKeyRepoStubForGroupUpdate) ListByGroupID(context.Context, int64, pagination.PaginationParams) ([]APIKey, *pagination.PaginationResult, error) {
	panic("unexpected")
placeholder
func (s *apiKeyRepoStubForGroupUpdate) SearchAPIKeys(context.Context, int64, string, int) ([]APIKey, error) {
	panic("unexpected")
placeholder
func (s *apiKeyRepoStubForGroupUpdate) ClearGroupIDByGroupID(context.Context, int64) (int64, error) {
	panic("unexpected")
placeholder
func (s *apiKeyRepoStubForGroupUpdate) CountByGroupID(context.Context, int64) (int64, error) {
	panic("unexpected")
placeholder
func (s *apiKeyRepoStubForGroupUpdate) ListKeysByUserID(context.Context, int64) ([]string, error) {
	panic("unexpected")
placeholder
func (s *apiKeyRepoStubForGroupUpdate) ListKeysByGroupID(context.Context, int64) ([]string, error) {
	panic("unexpected")
placeholder
func (s *apiKeyRepoStubForGroupUpdate) IncrementQuotaUsed(context.Context, int64, float64) (float64, error) {
	panic("unexpected")
placeholder
func (s *apiKeyRepoStubForGroupUpdate) UpdateLastUsed(context.Context, int64, time.Time) error {
	panic("unexpected")
placeholder

// groupRepoStubForGroupUpdate implements GroupRepository for AdminUpdateAPIKeyGroupID tests.
type groupRepoStubForGroupUpdate struct {
	group          *Group
	getErr         error
	lastGetByIDArg int64
placeholder

func (s *groupRepoStubForGroupUpdate) GetByID(_ context.Context, id int64) (*Group, error) {
	s.lastGetByIDArg = id
	if s.getErr != nil {
		return nil, s.getErr
placeholder
	return s.group, nil
placeholder

// Unused methods – panic on unexpected call.
func (s *groupRepoStubForGroupUpdate) Create(context.Context, *Group) error { panic("unexpected") placeholder
func (s *groupRepoStubForGroupUpdate) GetByIDLite(context.Context, int64) (*Group, error) {
	panic("unexpected")
placeholder
func (s *groupRepoStubForGroupUpdate) Update(context.Context, *Group) error { panic("unexpected") placeholder
func (s *groupRepoStubForGroupUpdate) Delete(context.Context, int64) error  { panic("unexpected") placeholder
func (s *groupRepoStubForGroupUpdate) DeleteCascade(context.Context, int64) ([]int64, error) {
	panic("unexpected")
placeholder
func (s *groupRepoStubForGroupUpdate) List(context.Context, pagination.PaginationParams) ([]Group, *pagination.PaginationResult, error) {
	panic("unexpected")
placeholder
func (s *groupRepoStubForGroupUpdate) ListWithFilters(context.Context, pagination.PaginationParams, string, string, string, *bool) ([]Group, *pagination.PaginationResult, error) {
	panic("unexpected")
placeholder
func (s *groupRepoStubForGroupUpdate) ListActive(context.Context) ([]Group, error) {
	panic("unexpected")
placeholder
func (s *groupRepoStubForGroupUpdate) ListActiveByPlatform(context.Context, string) ([]Group, error) {
	panic("unexpected")
placeholder
func (s *groupRepoStubForGroupUpdate) ExistsByName(context.Context, string) (bool, error) {
	panic("unexpected")
placeholder
func (s *groupRepoStubForGroupUpdate) GetAccountCount(context.Context, int64) (int64, error) {
	panic("unexpected")
placeholder
func (s *groupRepoStubForGroupUpdate) DeleteAccountGroupsByGroupID(context.Context, int64) (int64, error) {
	panic("unexpected")
placeholder
func (s *groupRepoStubForGroupUpdate) GetAccountIDsByGroupIDs(context.Context, []int64) ([]int64, error) {
	panic("unexpected")
placeholder
func (s *groupRepoStubForGroupUpdate) BindAccountsToGroup(context.Context, int64, []int64) error {
	panic("unexpected")
placeholder
func (s *groupRepoStubForGroupUpdate) UpdateSortOrders(context.Context, []GroupSortOrderUpdate) error {
	panic("unexpected")
placeholder

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestAdminService_AdminUpdateAPIKeyGroupID_KeyNotFound(t *testing.T) {
	repo := &apiKeyRepoStubForGroupUpdate{getErr: ErrAPIKeyNotFoundplaceholder
	svc := &adminServiceImpl{apiKeyRepo: repoplaceholder

	_, err := svc.AdminUpdateAPIKeyGroupID(context.Background(), 999, int64Ptr(1))
	require.ErrorIs(t, err, ErrAPIKeyNotFound)
placeholder

func TestAdminService_AdminUpdateAPIKeyGroupID_NilGroupID_NoOp(t *testing.T) {
	existing := &APIKey{ID: 1, Key: "sk-test", GroupID: int64Ptr(5)placeholder
	repo := &apiKeyRepoStubForGroupUpdate{key: existingplaceholder
	svc := &adminServiceImpl{apiKeyRepo: repoplaceholder

	got, err := svc.AdminUpdateAPIKeyGroupID(context.Background(), 1, nil)
placeholder
	require.Equal(t, int64(1), got.ID)
	// Update should NOT have been called (updated stays nil)
	require.Nil(t, repo.updated)
placeholder

func TestAdminService_AdminUpdateAPIKeyGroupID_Unbind(t *testing.T) {
	existing := &APIKey{ID: 1, Key: "sk-test", GroupID: int64Ptr(5), Group: &Group{ID: 5, Name: "Old"placeholderplaceholder
	repo := &apiKeyRepoStubForGroupUpdate{key: existingplaceholder
	cache := &authCacheInvalidatorStub{placeholder
	svc := &adminServiceImpl{apiKeyRepo: repo, authCacheInvalidator: cacheplaceholder

	got, err := svc.AdminUpdateAPIKeyGroupID(context.Background(), 1, int64Ptr(0))
placeholder
	require.Nil(t, got.GroupID, "group_id should be nil after unbind")
	require.Nil(t, got.Group, "group object should be nil after unbind")
	require.NotNil(t, repo.updated, "Update should have been called")
	require.Nil(t, repo.updated.GroupID)
	require.Equal(t, []string{"sk-test"placeholder, cache.keys, "cache should be invalidated")
placeholder

func TestAdminService_AdminUpdateAPIKeyGroupID_BindActiveGroup(t *testing.T) {
	existing := &APIKey{ID: 1, Key: "sk-test", GroupID: nilplaceholder
	apiKeyRepo := &apiKeyRepoStubForGroupUpdate{key: existingplaceholder
	groupRepo := &groupRepoStubForGroupUpdate{group: &Group{ID: 10, Name: "Pro", Status: StatusActiveplaceholderplaceholder
	cache := &authCacheInvalidatorStub{placeholder
	svc := &adminServiceImpl{apiKeyRepo: apiKeyRepo, groupRepo: groupRepo, authCacheInvalidator: cacheplaceholder

	got, err := svc.AdminUpdateAPIKeyGroupID(context.Background(), 1, int64Ptr(10))
placeholder
	require.NotNil(t, got.GroupID)
	require.Equal(t, int64(10), *got.GroupID)
	require.Equal(t, int64(10), *apiKeyRepo.updated.GroupID)
	require.Equal(t, []string{"sk-test"placeholder, cache.keys)
	// M3: verify correct group ID was passed to repo
	require.Equal(t, int64(10), groupRepo.lastGetByIDArg)
	// C1 fix: verify Group object is populated
	require.NotNil(t, got.Group)
	require.Equal(t, "Pro", got.Group.Name)
placeholder

func TestAdminService_AdminUpdateAPIKeyGroupID_SameGroup_Idempotent(t *testing.T) {
	existing := &APIKey{ID: 1, Key: "sk-test", GroupID: int64Ptr(10), Group: &Group{ID: 10, Name: "Pro"placeholderplaceholder
	apiKeyRepo := &apiKeyRepoStubForGroupUpdate{key: existingplaceholder
	groupRepo := &groupRepoStubForGroupUpdate{group: &Group{ID: 10, Name: "Pro", Status: StatusActiveplaceholderplaceholder
	cache := &authCacheInvalidatorStub{placeholder
	svc := &adminServiceImpl{apiKeyRepo: apiKeyRepo, groupRepo: groupRepo, authCacheInvalidator: cacheplaceholder

	got, err := svc.AdminUpdateAPIKeyGroupID(context.Background(), 1, int64Ptr(10))
placeholder
	require.NotNil(t, got.GroupID)
	require.Equal(t, int64(10), *got.GroupID)
	// Update is still called (current impl doesn't short-circuit on same group)
	require.NotNil(t, apiKeyRepo.updated)
	require.Equal(t, []string{"sk-test"placeholder, cache.keys)
placeholder

func TestAdminService_AdminUpdateAPIKeyGroupID_GroupNotFound(t *testing.T) {
	existing := &APIKey{ID: 1, Key: "sk-test"placeholder
	apiKeyRepo := &apiKeyRepoStubForGroupUpdate{key: existingplaceholder
	groupRepo := &groupRepoStubForGroupUpdate{getErr: ErrGroupNotFoundplaceholder
	svc := &adminServiceImpl{apiKeyRepo: apiKeyRepo, groupRepo: groupRepoplaceholder

	_, err := svc.AdminUpdateAPIKeyGroupID(context.Background(), 1, int64Ptr(99))
	require.ErrorIs(t, err, ErrGroupNotFound)
placeholder

func TestAdminService_AdminUpdateAPIKeyGroupID_GroupNotActive(t *testing.T) {
	existing := &APIKey{ID: 1, Key: "sk-test"placeholder
	apiKeyRepo := &apiKeyRepoStubForGroupUpdate{key: existingplaceholder
	groupRepo := &groupRepoStubForGroupUpdate{group: &Group{ID: 5, Status: StatusDisabledplaceholderplaceholder
	svc := &adminServiceImpl{apiKeyRepo: apiKeyRepo, groupRepo: groupRepoplaceholder

	_, err := svc.AdminUpdateAPIKeyGroupID(context.Background(), 1, int64Ptr(5))
placeholder
	require.Equal(t, "GROUP_NOT_ACTIVE", infraerrors.Reason(err))
placeholder

func TestAdminService_AdminUpdateAPIKeyGroupID_UpdateFails(t *testing.T) {
	existing := &APIKey{ID: 1, Key: "sk-test", GroupID: int64Ptr(3)placeholder
	repo := &apiKeyRepoStubForGroupUpdate{key: existing, updateErr: errors.New("db write error")placeholder
	svc := &adminServiceImpl{apiKeyRepo: repoplaceholder

	_, err := svc.AdminUpdateAPIKeyGroupID(context.Background(), 1, int64Ptr(0))
placeholder
	require.Contains(t, err.Error(), "update api key")
placeholder

func TestAdminService_AdminUpdateAPIKeyGroupID_NegativeGroupID(t *testing.T) {
	existing := &APIKey{ID: 1, Key: "sk-test"placeholder
	apiKeyRepo := &apiKeyRepoStubForGroupUpdate{key: existingplaceholder
	svc := &adminServiceImpl{apiKeyRepo: apiKeyRepoplaceholder

	_, err := svc.AdminUpdateAPIKeyGroupID(context.Background(), 1, int64Ptr(-5))
placeholder
	require.Equal(t, "INVALID_GROUP_ID", infraerrors.Reason(err))
placeholder

func TestAdminService_AdminUpdateAPIKeyGroupID_PointerIsolation(t *testing.T) {
	existing := &APIKey{ID: 1, Key: "sk-test", GroupID: nilplaceholder
	apiKeyRepo := &apiKeyRepoStubForGroupUpdate{key: existingplaceholder
	groupRepo := &groupRepoStubForGroupUpdate{group: &Group{ID: 10, Name: "Pro", Status: StatusActiveplaceholderplaceholder
	cache := &authCacheInvalidatorStub{placeholder
	svc := &adminServiceImpl{apiKeyRepo: apiKeyRepo, groupRepo: groupRepo, authCacheInvalidator: cacheplaceholder

	inputGID := int64(10)
	got, err := svc.AdminUpdateAPIKeyGroupID(context.Background(), 1, &inputGID)
placeholder
	require.NotNil(t, got.GroupID)
	// Mutating the input pointer must NOT affect the stored value
	inputGID = 999
	require.Equal(t, int64(10), *got.GroupID)
	require.Equal(t, int64(10), *apiKeyRepo.updated.GroupID)
placeholder

func TestAdminService_AdminUpdateAPIKeyGroupID_NilCacheInvalidator(t *testing.T) {
	existing := &APIKey{ID: 1, Key: "sk-test"placeholder
	apiKeyRepo := &apiKeyRepoStubForGroupUpdate{key: existingplaceholder
	groupRepo := &groupRepoStubForGroupUpdate{group: &Group{ID: 7, Status: StatusActiveplaceholderplaceholder
	// authCacheInvalidator is nil – should not panic
	svc := &adminServiceImpl{apiKeyRepo: apiKeyRepo, groupRepo: groupRepoplaceholder

	got, err := svc.AdminUpdateAPIKeyGroupID(context.Background(), 1, int64Ptr(7))
placeholder
	require.NotNil(t, got.GroupID)
	require.Equal(t, int64(7), *got.GroupID)
placeholder
