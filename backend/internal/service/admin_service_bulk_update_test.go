//go:build unit

package service

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

type accountRepoStubForBulkUpdate struct {
	accountRepoStub
	bulkUpdateErr       error
	bulkUpdateIDs       []int64
	bindGroupErrByID    map[int64]error
	bindGroupsCalls     []int64
	bindGroupsByAccount map[int64][]int64
	createAccount       *Account
	createID            int64
	createErr           error
	updatedAccounts     []*Account
	updateErr           error
	getByIDsAccounts    []*Account
	getByIDsErr         error
	getByIDsCalled      bool
	getByIDsIDs         []int64
	getByIDAccounts     map[int64]*Account
	getByIDErrByID      map[int64]error
	getByIDCalled       []int64
	listByGroupData     map[int64][]Account
	listByGroupErr      map[int64]error
	listData            []Account
	listResult          *pagination.PaginationResult
	listErr             error
	listCalled          bool
	lastListParams      pagination.PaginationParams
	lastListFilters     struct {
		platform    string
		accountType string
		status      string
		search      string
		groupID     int64
		privacyMode string
placeholder
placeholder

func (s *accountRepoStubForBulkUpdate) BulkUpdate(_ context.Context, ids []int64, _ AccountBulkUpdate) (int64, error) {
	s.bulkUpdateIDs = append([]int64{placeholder, ids...)
	if s.bulkUpdateErr != nil {
		return 0, s.bulkUpdateErr
placeholder
	return int64(len(ids)), nil
placeholder

func (s *accountRepoStubForBulkUpdate) Create(_ context.Context, account *Account) error {
	s.createAccount = account
	if s.createID > 0 {
		account.ID = s.createID
placeholder
	return s.createErr
placeholder

func (s *accountRepoStubForBulkUpdate) Update(_ context.Context, account *Account) error {
	s.updatedAccounts = append(s.updatedAccounts, account)
	return s.updateErr
placeholder

func (s *accountRepoStubForBulkUpdate) BindGroups(_ context.Context, accountID int64, groupIDs []int64) error {
	s.bindGroupsCalls = append(s.bindGroupsCalls, accountID)
	if s.bindGroupsByAccount == nil {
		s.bindGroupsByAccount = make(map[int64][]int64)
placeholder
	s.bindGroupsByAccount[accountID] = append([]int64{placeholder, groupIDs...)
	if err, ok := s.bindGroupErrByID[accountID]; ok {
		return err
placeholder
	return nil
placeholder

func (s *accountRepoStubForBulkUpdate) GetByIDs(_ context.Context, ids []int64) ([]*Account, error) {
	s.getByIDsCalled = true
	s.getByIDsIDs = append([]int64{placeholder, ids...)
	if s.getByIDsErr != nil {
		return nil, s.getByIDsErr
placeholder
	return s.getByIDsAccounts, nil
placeholder

func (s *accountRepoStubForBulkUpdate) GetByID(_ context.Context, id int64) (*Account, error) {
	s.getByIDCalled = append(s.getByIDCalled, id)
	if err, ok := s.getByIDErrByID[id]; ok {
		return nil, err
placeholder
	if account, ok := s.getByIDAccounts[id]; ok {
		return account, nil
placeholder
	return nil, errors.New("account not found")
placeholder

func (s *accountRepoStubForBulkUpdate) ListByGroup(_ context.Context, groupID int64) ([]Account, error) {
	if err, ok := s.listByGroupErr[groupID]; ok {
		return nil, err
placeholder
	if rows, ok := s.listByGroupData[groupID]; ok {
		return rows, nil
placeholder
	return nil, nil
placeholder

func (s *accountRepoStubForBulkUpdate) ListAllWithFilters(context.Context, string, string, string, string, int64, string) ([]Account, error) {
	return nil, nil
placeholder

func (s *accountRepoStubForBulkUpdate) ListWithFilters(_ context.Context, params pagination.PaginationParams, platform, accountType, status, search string, groupID int64, privacyMode string) ([]Account, *pagination.PaginationResult, error) {
	s.listCalled = true
	s.lastListParams = params
	s.lastListFilters.platform = platform
	s.lastListFilters.accountType = accountType
	s.lastListFilters.status = status
	s.lastListFilters.search = search
	s.lastListFilters.groupID = groupID
	s.lastListFilters.privacyMode = privacyMode
	if s.listErr != nil {
		return nil, nil, s.listErr
placeholder
	if s.listResult != nil {
		return s.listData, s.listResult, nil
placeholder
	return s.listData, &pagination.PaginationResult{Total: int64(len(s.listData))placeholder, nil
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

func TestAdminService_BulkUpdateAccounts_RejectsRateChangeForSyncedAccounts(t *testing.T) {
	repo := &accountRepoStubForBulkUpdate{
		getByIDsAccounts: []*Account{
			{
				ID: 1,
				Extra: map[string]any{
					UpstreamBillingProbeEnabledExtraKey:    true,
					UpstreamBillingRateSyncEnabledExtraKey: true,
			placeholder,
		placeholder,
			{ID: 2, Extra: map[string]any{placeholderplaceholder,
	placeholder,
placeholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder
	rateMultiplier := 0.5

	result, err := svc.BulkUpdateAccounts(context.Background(), &BulkUpdateAccountsInput{
		AccountIDs:     []int64{1, 2placeholder,
		RateMultiplier: &rateMultiplier,
placeholder)

	require.Nil(t, result)
placeholder
	var appErr *infraerrors.ApplicationError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, int32(http.StatusConflict), appErr.Code)
	require.Equal(t, "UPSTREAM_BILLING_RATE_SYNC_BULK_CONFLICT", appErr.Reason)
	require.Equal(t, "1", appErr.Metadata["count"])
	require.True(t, repo.getByIDsCalled)
	require.Empty(t, repo.bulkUpdateIDs, "rate conflict must be rejected before any write")
placeholder

// TestAdminService_BulkUpdateAccounts_PartialFailureIDs 验证部分失败时 success_ids/failed_ids 正确。
func TestAdminService_BulkUpdateAccounts_PartialFailureIDs(t *testing.T) {
	repo := &accountRepoStubForBulkUpdate{
		bindGroupErrByID: map[int64]error{
			2: errors.New("bind failed"),
	placeholder,
placeholder
	svc := &adminServiceImpl{
		accountRepo: repo,
		groupRepo:   &groupRepoStubForAdmin{getByID: &Group{ID: 10, Name: "g10"placeholderplaceholder,
placeholder

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

func TestAdminService_BulkUpdateAccounts_NilGroupRepoReturnsError(t *testing.T) {
	repo := &accountRepoStubForBulkUpdate{placeholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder

	groupIDs := []int64{10placeholder
	input := &BulkUpdateAccountsInput{
		AccountIDs: []int64{1placeholder,
		GroupIDs:   &groupIDs,
placeholder

	result, err := svc.BulkUpdateAccounts(context.Background(), input)
	require.Nil(t, result)
placeholder
	require.Contains(t, err.Error(), "group repository not configured")
placeholder

// TestAdminService_BulkUpdateAccounts_MixedChannelPreCheckBlocksOnExistingConflict verifies
// that the global pre-check detects a conflict with existing group members and returns an
// error before any DB write is performed.
func TestAdminService_BulkUpdateAccounts_MixedChannelPreCheckBlocksOnExistingConflict(t *testing.T) {
	repo := &accountRepoStubForBulkUpdate{
		getByIDsAccounts: []*Account{
			{ID: 1, Platform: PlatformAntigravityplaceholder,
	placeholder,
		// Group 10 already contains an Anthropic account.
		listByGroupData: map[int64][]Account{
			10: {{ID: 99, Platform: PlatformAnthropicplaceholderplaceholder,
	placeholder,
placeholder
	svc := &adminServiceImpl{
		accountRepo: repo,
		groupRepo:   &groupRepoStubForAdmin{getByID: &Group{ID: 10, Name: "target-group"placeholderplaceholder,
placeholder

	groupIDs := []int64{10placeholder
	input := &BulkUpdateAccountsInput{
		AccountIDs: []int64{1placeholder,
		GroupIDs:   &groupIDs,
placeholder

	result, err := svc.BulkUpdateAccounts(context.Background(), input)
	require.Nil(t, result)
placeholder
	require.Contains(t, err.Error(), "mixed channel")
	// No BindGroups should have been called since the check runs before any write.
	require.Empty(t, repo.bindGroupsCalls)
placeholder

func TestAdminServiceBulkUpdateAccounts_ResolvesIDsFromFilters(t *testing.T) {
	repo := &accountRepoStubForBulkUpdate{
		listData: []Account{
			{ID: 7placeholder,
			{ID: 11placeholder,
	placeholder,
		listResult: &pagination.PaginationResult{Total: 2placeholder,
placeholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder

	schedulable := true
	input := &BulkUpdateAccountsInput{
		Schedulable: &schedulable,
placeholder

	filtersField := reflect.ValueOf(input).Elem().FieldByName("Filters")
	require.True(t, filtersField.IsValid(), "BulkUpdateAccountsInput should expose Filters for filter-target bulk update")
	require.Equal(t, reflect.Ptr, filtersField.Kind(), "BulkUpdateAccountsInput.Filters should be a pointer field")

	filtersValue := reflect.New(filtersField.Type().Elem())
	filtersValue.Elem().FieldByName("Platform").SetString(PlatformOpenAI)
	filtersValue.Elem().FieldByName("Type").SetString(AccountTypeOAuth)
	filtersValue.Elem().FieldByName("Status").SetString(StatusActive)
	filtersValue.Elem().FieldByName("Group").SetString("12")
	filtersValue.Elem().FieldByName("PrivacyMode").SetString(PrivacyModeCFBlocked)
	filtersValue.Elem().FieldByName("Search").SetString("bulk-target")
	filtersField.Set(filtersValue)

	result, err := svc.BulkUpdateAccounts(context.Background(), input)
placeholder
	require.True(t, repo.listCalled, "expected filter-target bulk update to resolve matching IDs via account list filters")
	require.Equal(t, PlatformOpenAI, repo.lastListFilters.platform)
	require.Equal(t, AccountTypeOAuth, repo.lastListFilters.accountType)
	require.Equal(t, StatusActive, repo.lastListFilters.status)
	require.Equal(t, "bulk-target", repo.lastListFilters.search)
	require.Equal(t, int64(12), repo.lastListFilters.groupID)
	require.Equal(t, PrivacyModeCFBlocked, repo.lastListFilters.privacyMode)
	require.Equal(t, []int64{7, 11placeholder, repo.bulkUpdateIDs)
	require.Equal(t, 2, result.Success)
	require.Equal(t, 0, result.Failed)
	require.Equal(t, []int64{7, 11placeholder, result.SuccessIDs)
placeholder
