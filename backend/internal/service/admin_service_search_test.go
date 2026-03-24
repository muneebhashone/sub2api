//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

type accountRepoStubForAdminList struct {
	accountRepoStub

	listWithFiltersCalls    int
	listWithFiltersParams   pagination.PaginationParams
	listWithFiltersPlatform string
	listWithFiltersType     string
	listWithFiltersStatus   string
	listWithFiltersSearch   string
	listWithFiltersPrivacy  string
	listWithFiltersAccounts []Account
	listWithFiltersResult   *pagination.PaginationResult
	listWithFiltersErr      error
placeholder

func (s *accountRepoStubForAdminList) ListWithFilters(_ context.Context, params pagination.PaginationParams, platform, accountType, status, search string, groupID int64, privacyMode string) ([]Account, *pagination.PaginationResult, error) {
	s.listWithFiltersCalls++
	s.listWithFiltersParams = params
	s.listWithFiltersPlatform = platform
	s.listWithFiltersType = accountType
	s.listWithFiltersStatus = status
	s.listWithFiltersSearch = search
	s.listWithFiltersPrivacy = privacyMode

	if s.listWithFiltersErr != nil {
		return nil, nil, s.listWithFiltersErr
placeholder

	result := s.listWithFiltersResult
	if result == nil {
		result = &pagination.PaginationResult{
			Total:    int64(len(s.listWithFiltersAccounts)),
			Page:     params.Page,
			PageSize: params.PageSize,
	placeholder
placeholder

	return s.listWithFiltersAccounts, result, nil
placeholder

type proxyRepoStubForAdminList struct {
	proxyRepoStub

	listWithFiltersCalls    int
	listWithFiltersParams   pagination.PaginationParams
	listWithFiltersProtocol string
	listWithFiltersStatus   string
	listWithFiltersSearch   string
	listWithFiltersProxies  []Proxy
	listWithFiltersResult   *pagination.PaginationResult
	listWithFiltersErr      error

	listWithFiltersAndAccountCountCalls    int
	listWithFiltersAndAccountCountParams   pagination.PaginationParams
	listWithFiltersAndAccountCountProtocol string
	listWithFiltersAndAccountCountStatus   string
	listWithFiltersAndAccountCountSearch   string
	listWithFiltersAndAccountCountProxies  []ProxyWithAccountCount
	listWithFiltersAndAccountCountResult   *pagination.PaginationResult
	listWithFiltersAndAccountCountErr      error
placeholder

func (s *proxyRepoStubForAdminList) ListWithFilters(_ context.Context, params pagination.PaginationParams, protocol, status, search string) ([]Proxy, *pagination.PaginationResult, error) {
	s.listWithFiltersCalls++
	s.listWithFiltersParams = params
	s.listWithFiltersProtocol = protocol
	s.listWithFiltersStatus = status
	s.listWithFiltersSearch = search

	if s.listWithFiltersErr != nil {
		return nil, nil, s.listWithFiltersErr
placeholder

	result := s.listWithFiltersResult
	if result == nil {
		result = &pagination.PaginationResult{
			Total:    int64(len(s.listWithFiltersProxies)),
			Page:     params.Page,
			PageSize: params.PageSize,
	placeholder
placeholder

	return s.listWithFiltersProxies, result, nil
placeholder

func (s *proxyRepoStubForAdminList) ListWithFiltersAndAccountCount(_ context.Context, params pagination.PaginationParams, protocol, status, search string) ([]ProxyWithAccountCount, *pagination.PaginationResult, error) {
	s.listWithFiltersAndAccountCountCalls++
	s.listWithFiltersAndAccountCountParams = params
	s.listWithFiltersAndAccountCountProtocol = protocol
	s.listWithFiltersAndAccountCountStatus = status
	s.listWithFiltersAndAccountCountSearch = search

	if s.listWithFiltersAndAccountCountErr != nil {
		return nil, nil, s.listWithFiltersAndAccountCountErr
placeholder

	result := s.listWithFiltersAndAccountCountResult
	if result == nil {
		result = &pagination.PaginationResult{
			Total:    int64(len(s.listWithFiltersAndAccountCountProxies)),
			Page:     params.Page,
			PageSize: params.PageSize,
	placeholder
placeholder

	return s.listWithFiltersAndAccountCountProxies, result, nil
placeholder

type redeemRepoStubForAdminList struct {
	redeemRepoStub

	listWithFiltersCalls  int
	listWithFiltersParams pagination.PaginationParams
	listWithFiltersType   string
	listWithFiltersStatus string
	listWithFiltersSearch string
	listWithFiltersCodes  []RedeemCode
	listWithFiltersResult *pagination.PaginationResult
	listWithFiltersErr    error
placeholder

func (s *redeemRepoStubForAdminList) ListWithFilters(_ context.Context, params pagination.PaginationParams, codeType, status, search string) ([]RedeemCode, *pagination.PaginationResult, error) {
	s.listWithFiltersCalls++
	s.listWithFiltersParams = params
	s.listWithFiltersType = codeType
	s.listWithFiltersStatus = status
	s.listWithFiltersSearch = search

	if s.listWithFiltersErr != nil {
		return nil, nil, s.listWithFiltersErr
placeholder

	result := s.listWithFiltersResult
	if result == nil {
		result = &pagination.PaginationResult{
			Total:    int64(len(s.listWithFiltersCodes)),
			Page:     params.Page,
			PageSize: params.PageSize,
	placeholder
placeholder

	return s.listWithFiltersCodes, result, nil
placeholder

func (s *redeemRepoStubForAdminList) ListByUserPaginated(_ context.Context, userID int64, params pagination.PaginationParams, codeType string) ([]RedeemCode, *pagination.PaginationResult, error) {
	panic("unexpected ListByUserPaginated call")
placeholder

func (s *redeemRepoStubForAdminList) SumPositiveBalanceByUser(_ context.Context, userID int64) (float64, error) {
	panic("unexpected SumPositiveBalanceByUser call")
placeholder

func TestAdminService_ListAccounts_WithSearch(t *testing.T) {
	t.Run("search 参数正常传递到 repository 层", func(t *testing.T) {
		repo := &accountRepoStubForAdminList{
			listWithFiltersAccounts: []Account{{ID: 1, Name: "acc"placeholderplaceholder,
			listWithFiltersResult:   &pagination.PaginationResult{Total: 10placeholder,
	placeholder
		svc := &adminServiceImpl{accountRepo: repoplaceholder

		accounts, total, err := svc.ListAccounts(context.Background(), 1, 20, PlatformGemini, AccountTypeOAuth, StatusActive, "acc", 0, "")
	placeholder
		require.Equal(t, int64(10), total)
		require.Equal(t, []Account{{ID: 1, Name: "acc"placeholderplaceholder, accounts)

		require.Equal(t, 1, repo.listWithFiltersCalls)
		require.Equal(t, pagination.PaginationParams{Page: 1, PageSize: 20placeholder, repo.listWithFiltersParams)
		require.Equal(t, PlatformGemini, repo.listWithFiltersPlatform)
		require.Equal(t, AccountTypeOAuth, repo.listWithFiltersType)
		require.Equal(t, StatusActive, repo.listWithFiltersStatus)
		require.Equal(t, "acc", repo.listWithFiltersSearch)
placeholder)
placeholder

func TestAdminService_ListAccounts_WithPrivacyMode(t *testing.T) {
	t.Run("privacy_mode 参数正常传递到 repository 层", func(t *testing.T) {
		repo := &accountRepoStubForAdminList{
			listWithFiltersAccounts: []Account{{ID: 2, Name: "acc2"placeholderplaceholder,
			listWithFiltersResult:   &pagination.PaginationResult{Total: 1placeholder,
	placeholder
		svc := &adminServiceImpl{accountRepo: repoplaceholder

		accounts, total, err := svc.ListAccounts(context.Background(), 1, 20, PlatformOpenAI, AccountTypeOAuth, StatusActive, "acc2", 0, PrivacyModeCFBlocked)
	placeholder
		require.Equal(t, int64(1), total)
		require.Equal(t, []Account{{ID: 2, Name: "acc2"placeholderplaceholder, accounts)
		require.Equal(t, PrivacyModeCFBlocked, repo.listWithFiltersPrivacy)
placeholder)
placeholder

func TestAdminService_ListProxies_WithSearch(t *testing.T) {
	t.Run("search 参数正常传递到 repository 层", func(t *testing.T) {
		repo := &proxyRepoStubForAdminList{
			listWithFiltersProxies: []Proxy{{ID: 2, Name: "p1"placeholderplaceholder,
			listWithFiltersResult:  &pagination.PaginationResult{Total: 7placeholder,
	placeholder
		svc := &adminServiceImpl{proxyRepo: repoplaceholder

		proxies, total, err := svc.ListProxies(context.Background(), 3, 50, "http", StatusActive, "p1")
	placeholder
		require.Equal(t, int64(7), total)
		require.Equal(t, []Proxy{{ID: 2, Name: "p1"placeholderplaceholder, proxies)

		require.Equal(t, 1, repo.listWithFiltersCalls)
		require.Equal(t, pagination.PaginationParams{Page: 3, PageSize: 50placeholder, repo.listWithFiltersParams)
		require.Equal(t, "http", repo.listWithFiltersProtocol)
		require.Equal(t, StatusActive, repo.listWithFiltersStatus)
		require.Equal(t, "p1", repo.listWithFiltersSearch)
placeholder)
placeholder

func TestAdminService_ListProxiesWithAccountCount_WithSearch(t *testing.T) {
	t.Run("search 参数正常传递到 repository 层", func(t *testing.T) {
		repo := &proxyRepoStubForAdminList{
			listWithFiltersAndAccountCountProxies: []ProxyWithAccountCount{{Proxy: Proxy{ID: 3, Name: "p2"placeholder, AccountCount: 5placeholderplaceholder,
			listWithFiltersAndAccountCountResult:  &pagination.PaginationResult{Total: 9placeholder,
	placeholder
		svc := &adminServiceImpl{proxyRepo: repoplaceholder

		proxies, total, err := svc.ListProxiesWithAccountCount(context.Background(), 2, 10, "socks5", StatusDisabled, "p2")
	placeholder
		require.Equal(t, int64(9), total)
		require.Equal(t, []ProxyWithAccountCount{{Proxy: Proxy{ID: 3, Name: "p2"placeholder, AccountCount: 5placeholderplaceholder, proxies)

		require.Equal(t, 1, repo.listWithFiltersAndAccountCountCalls)
		require.Equal(t, pagination.PaginationParams{Page: 2, PageSize: 10placeholder, repo.listWithFiltersAndAccountCountParams)
		require.Equal(t, "socks5", repo.listWithFiltersAndAccountCountProtocol)
		require.Equal(t, StatusDisabled, repo.listWithFiltersAndAccountCountStatus)
		require.Equal(t, "p2", repo.listWithFiltersAndAccountCountSearch)
placeholder)
placeholder

func TestAdminService_ListRedeemCodes_WithSearch(t *testing.T) {
	t.Run("search 参数正常传递到 repository 层", func(t *testing.T) {
		repo := &redeemRepoStubForAdminList{
			listWithFiltersCodes:  []RedeemCode{{ID: 4, Code: "ABC"placeholderplaceholder,
			listWithFiltersResult: &pagination.PaginationResult{Total: 3placeholder,
	placeholder
		svc := &adminServiceImpl{redeemCodeRepo: repoplaceholder

		codes, total, err := svc.ListRedeemCodes(context.Background(), 1, 20, RedeemTypeBalance, StatusUnused, "ABC")
	placeholder
		require.Equal(t, int64(3), total)
		require.Equal(t, []RedeemCode{{ID: 4, Code: "ABC"placeholderplaceholder, codes)

		require.Equal(t, 1, repo.listWithFiltersCalls)
		require.Equal(t, pagination.PaginationParams{Page: 1, PageSize: 20placeholder, repo.listWithFiltersParams)
		require.Equal(t, RedeemTypeBalance, repo.listWithFiltersType)
		require.Equal(t, StatusUnused, repo.listWithFiltersStatus)
		require.Equal(t, "ABC", repo.listWithFiltersSearch)
placeholder)
placeholder
