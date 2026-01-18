package admin

import (
	"context"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

type stubAdminService struct {
	users       []service.User
	apiKeys     []service.APIKey
	groups      []service.Group
	accounts    []service.Account
	proxies     []service.Proxy
	proxyCounts []service.ProxyWithAccountCount
	redeems     []service.RedeemCode
placeholder

func newStubAdminService() *stubAdminService {
	now := time.Now().UTC()
	user := service.User{
		ID:        1,
		Email:     "user@example.com",
		Role:      service.RoleUser,
		Status:    service.StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
placeholder
	apiKey := service.APIKey{
		ID:        10,
		UserID:    user.ID,
		Key:       "sk-test",
		Name:      "test",
		Status:    service.StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
placeholder
	group := service.Group{
		ID:        2,
		Name:      "group",
		Platform:  service.PlatformAnthropic,
		Status:    service.StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
placeholder
	account := service.Account{
		ID:        3,
		Name:      "account",
		Platform:  service.PlatformAnthropic,
		Type:      service.AccountTypeOAuth,
		Status:    service.StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
placeholder
	proxy := service.Proxy{
		ID:        4,
		Name:      "proxy",
		Protocol:  "http",
		Host:      "127.0.0.1",
		Port:      8080,
		Status:    service.StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
placeholder
	redeem := service.RedeemCode{
		ID:        5,
		Code:      "R-TEST",
		Type:      service.RedeemTypeBalance,
		Value:     10,
		Status:    service.StatusUnused,
		CreatedAt: now,
placeholder
	return &stubAdminService{
		users:       []service.User{userplaceholder,
		apiKeys:     []service.APIKey{apiKeyplaceholder,
		groups:      []service.Group{groupplaceholder,
		accounts:    []service.Account{accountplaceholder,
		proxies:     []service.Proxy{proxyplaceholder,
		proxyCounts: []service.ProxyWithAccountCount{{Proxy: proxy, AccountCount: 1placeholderplaceholder,
		redeems:     []service.RedeemCode{redeemplaceholder,
placeholder
placeholder

func (s *stubAdminService) ListUsers(ctx context.Context, page, pageSize int, filters service.UserListFilters) ([]service.User, int64, error) {
	return s.users, int64(len(s.users)), nil
placeholder

func (s *stubAdminService) GetUser(ctx context.Context, id int64) (*service.User, error) {
	for i := range s.users {
		if s.users[i].ID == id {
			return &s.users[i], nil
	placeholder
placeholder
	user := service.User{ID: id, Email: "user@example.com", Status: service.StatusActiveplaceholder
	return &user, nil
placeholder

func (s *stubAdminService) CreateUser(ctx context.Context, input *service.CreateUserInput) (*service.User, error) {
	user := service.User{ID: 100, Email: input.Email, Status: service.StatusActiveplaceholder
	return &user, nil
placeholder

func (s *stubAdminService) UpdateUser(ctx context.Context, id int64, input *service.UpdateUserInput) (*service.User, error) {
	user := service.User{ID: id, Email: "updated@example.com", Status: service.StatusActiveplaceholder
	return &user, nil
placeholder

func (s *stubAdminService) DeleteUser(ctx context.Context, id int64) error {
	return nil
placeholder

func (s *stubAdminService) UpdateUserBalance(ctx context.Context, userID int64, balance float64, operation string, notes string) (*service.User, error) {
	user := service.User{ID: userID, Balance: balance, Status: service.StatusActiveplaceholder
	return &user, nil
placeholder

func (s *stubAdminService) GetUserAPIKeys(ctx context.Context, userID int64, page, pageSize int) ([]service.APIKey, int64, error) {
	return s.apiKeys, int64(len(s.apiKeys)), nil
placeholder

func (s *stubAdminService) GetUserUsageStats(ctx context.Context, userID int64, period string) (any, error) {
	return map[string]any{"user_id": userIDplaceholder, nil
placeholder

func (s *stubAdminService) ListGroups(ctx context.Context, page, pageSize int, platform, status, search string, isExclusive *bool) ([]service.Group, int64, error) {
	return s.groups, int64(len(s.groups)), nil
placeholder

func (s *stubAdminService) GetAllGroups(ctx context.Context) ([]service.Group, error) {
	return s.groups, nil
placeholder

func (s *stubAdminService) GetAllGroupsByPlatform(ctx context.Context, platform string) ([]service.Group, error) {
	return s.groups, nil
placeholder

func (s *stubAdminService) GetGroup(ctx context.Context, id int64) (*service.Group, error) {
	group := service.Group{ID: id, Name: "group", Status: service.StatusActiveplaceholder
	return &group, nil
placeholder

func (s *stubAdminService) CreateGroup(ctx context.Context, input *service.CreateGroupInput) (*service.Group, error) {
	group := service.Group{ID: 200, Name: input.Name, Status: service.StatusActiveplaceholder
	return &group, nil
placeholder

func (s *stubAdminService) UpdateGroup(ctx context.Context, id int64, input *service.UpdateGroupInput) (*service.Group, error) {
	group := service.Group{ID: id, Name: input.Name, Status: service.StatusActiveplaceholder
	return &group, nil
placeholder

func (s *stubAdminService) DeleteGroup(ctx context.Context, id int64) error {
	return nil
placeholder

func (s *stubAdminService) GetGroupAPIKeys(ctx context.Context, groupID int64, page, pageSize int) ([]service.APIKey, int64, error) {
	return s.apiKeys, int64(len(s.apiKeys)), nil
placeholder

func (s *stubAdminService) ListAccounts(ctx context.Context, page, pageSize int, platform, accountType, status, search string) ([]service.Account, int64, error) {
	return s.accounts, int64(len(s.accounts)), nil
placeholder

func (s *stubAdminService) GetAccount(ctx context.Context, id int64) (*service.Account, error) {
	account := service.Account{ID: id, Name: "account", Status: service.StatusActiveplaceholder
	return &account, nil
placeholder

func (s *stubAdminService) GetAccountsByIDs(ctx context.Context, ids []int64) ([]*service.Account, error) {
	out := make([]*service.Account, 0, len(ids))
	for _, id := range ids {
		account := service.Account{ID: id, Name: "account", Status: service.StatusActiveplaceholder
		out = append(out, &account)
placeholder
	return out, nil
placeholder

func (s *stubAdminService) CreateAccount(ctx context.Context, input *service.CreateAccountInput) (*service.Account, error) {
	account := service.Account{ID: 300, Name: input.Name, Status: service.StatusActiveplaceholder
	return &account, nil
placeholder

func (s *stubAdminService) UpdateAccount(ctx context.Context, id int64, input *service.UpdateAccountInput) (*service.Account, error) {
	account := service.Account{ID: id, Name: input.Name, Status: service.StatusActiveplaceholder
	return &account, nil
placeholder

func (s *stubAdminService) DeleteAccount(ctx context.Context, id int64) error {
	return nil
placeholder

func (s *stubAdminService) RefreshAccountCredentials(ctx context.Context, id int64) (*service.Account, error) {
	account := service.Account{ID: id, Name: "account", Status: service.StatusActiveplaceholder
	return &account, nil
placeholder

func (s *stubAdminService) ClearAccountError(ctx context.Context, id int64) (*service.Account, error) {
	account := service.Account{ID: id, Name: "account", Status: service.StatusActiveplaceholder
	return &account, nil
placeholder

func (s *stubAdminService) SetAccountSchedulable(ctx context.Context, id int64, schedulable bool) (*service.Account, error) {
	account := service.Account{ID: id, Name: "account", Status: service.StatusActive, Schedulable: schedulableplaceholder
	return &account, nil
placeholder

func (s *stubAdminService) BulkUpdateAccounts(ctx context.Context, input *service.BulkUpdateAccountsInput) (*service.BulkUpdateAccountsResult, error) {
	return &service.BulkUpdateAccountsResult{Success: 1, Failed: 0, SuccessIDs: []int64{1placeholderplaceholder, nil
placeholder

func (s *stubAdminService) ListProxies(ctx context.Context, page, pageSize int, protocol, status, search string) ([]service.Proxy, int64, error) {
	return s.proxies, int64(len(s.proxies)), nil
placeholder

func (s *stubAdminService) ListProxiesWithAccountCount(ctx context.Context, page, pageSize int, protocol, status, search string) ([]service.ProxyWithAccountCount, int64, error) {
	return s.proxyCounts, int64(len(s.proxyCounts)), nil
placeholder

func (s *stubAdminService) GetAllProxies(ctx context.Context) ([]service.Proxy, error) {
	return s.proxies, nil
placeholder

func (s *stubAdminService) GetAllProxiesWithAccountCount(ctx context.Context) ([]service.ProxyWithAccountCount, error) {
	return s.proxyCounts, nil
placeholder

func (s *stubAdminService) GetProxy(ctx context.Context, id int64) (*service.Proxy, error) {
	proxy := service.Proxy{ID: id, Name: "proxy", Status: service.StatusActiveplaceholder
	return &proxy, nil
placeholder

func (s *stubAdminService) CreateProxy(ctx context.Context, input *service.CreateProxyInput) (*service.Proxy, error) {
	proxy := service.Proxy{ID: 400, Name: input.Name, Status: service.StatusActiveplaceholder
	return &proxy, nil
placeholder

func (s *stubAdminService) UpdateProxy(ctx context.Context, id int64, input *service.UpdateProxyInput) (*service.Proxy, error) {
	proxy := service.Proxy{ID: id, Name: input.Name, Status: service.StatusActiveplaceholder
	return &proxy, nil
placeholder

func (s *stubAdminService) DeleteProxy(ctx context.Context, id int64) error {
	return nil
placeholder

func (s *stubAdminService) BatchDeleteProxies(ctx context.Context, ids []int64) (*service.ProxyBatchDeleteResult, error) {
	return &service.ProxyBatchDeleteResult{DeletedIDs: idsplaceholder, nil
placeholder

func (s *stubAdminService) GetProxyAccounts(ctx context.Context, proxyID int64) ([]service.ProxyAccountSummary, error) {
	return []service.ProxyAccountSummary{{ID: 1, Name: "account"placeholderplaceholder, nil
placeholder

func (s *stubAdminService) CheckProxyExists(ctx context.Context, host string, port int, username, password string) (bool, error) {
	return false, nil
placeholder

func (s *stubAdminService) TestProxy(ctx context.Context, id int64) (*service.ProxyTestResult, error) {
	return &service.ProxyTestResult{Success: true, Message: "ok"placeholder, nil
placeholder

func (s *stubAdminService) ListRedeemCodes(ctx context.Context, page, pageSize int, codeType, status, search string) ([]service.RedeemCode, int64, error) {
	return s.redeems, int64(len(s.redeems)), nil
placeholder

func (s *stubAdminService) GetRedeemCode(ctx context.Context, id int64) (*service.RedeemCode, error) {
	code := service.RedeemCode{ID: id, Code: "R-TEST", Status: service.StatusUnusedplaceholder
	return &code, nil
placeholder

func (s *stubAdminService) GenerateRedeemCodes(ctx context.Context, input *service.GenerateRedeemCodesInput) ([]service.RedeemCode, error) {
	return s.redeems, nil
placeholder

func (s *stubAdminService) DeleteRedeemCode(ctx context.Context, id int64) error {
	return nil
placeholder

func (s *stubAdminService) BatchDeleteRedeemCodes(ctx context.Context, ids []int64) (int64, error) {
	return int64(len(ids)), nil
placeholder

func (s *stubAdminService) ExpireRedeemCode(ctx context.Context, id int64) (*service.RedeemCode, error) {
	code := service.RedeemCode{ID: id, Code: "R-TEST", Status: service.StatusUsedplaceholder
	return &code, nil
placeholder

// Ensure stub implements interface.
var _ service.AdminService = (*stubAdminService)(nil)
