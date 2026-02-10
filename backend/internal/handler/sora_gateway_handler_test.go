//go:build unit

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/usagestats"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/Wei-Shaw/sub2api/internal/testutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// 编译期接口断言
var _ service.SoraClient = (*stubSoraClient)(nil)
var _ service.AccountRepository = (*stubAccountRepo)(nil)
var _ service.GroupRepository = (*stubGroupRepo)(nil)
var _ service.UsageLogRepository = (*stubUsageLogRepo)(nil)

type stubSoraClient struct {
	imageURLs []string
placeholder

func (s *stubSoraClient) Enabled() bool { return true placeholder
func (s *stubSoraClient) UploadImage(ctx context.Context, account *service.Account, data []byte, filename string) (string, error) {
	return "upload", nil
placeholder
func (s *stubSoraClient) CreateImageTask(ctx context.Context, account *service.Account, req service.SoraImageRequest) (string, error) {
	return "task-image", nil
placeholder
func (s *stubSoraClient) CreateVideoTask(ctx context.Context, account *service.Account, req service.SoraVideoRequest) (string, error) {
	return "task-video", nil
placeholder
func (s *stubSoraClient) GetImageTask(ctx context.Context, account *service.Account, taskID string) (*service.SoraImageTaskStatus, error) {
	return &service.SoraImageTaskStatus{ID: taskID, Status: "completed", URLs: s.imageURLsplaceholder, nil
placeholder
func (s *stubSoraClient) GetVideoTask(ctx context.Context, account *service.Account, taskID string) (*service.SoraVideoTaskStatus, error) {
	return &service.SoraVideoTaskStatus{ID: taskID, Status: "completed", URLs: s.imageURLsplaceholder, nil
placeholder

type stubAccountRepo struct {
	accounts map[int64]*service.Account
placeholder

func (r *stubAccountRepo) Create(ctx context.Context, account *service.Account) error { return nil placeholder
func (r *stubAccountRepo) GetByID(ctx context.Context, id int64) (*service.Account, error) {
	if acc, ok := r.accounts[id]; ok {
		return acc, nil
placeholder
	return nil, service.ErrAccountNotFound
placeholder
func (r *stubAccountRepo) GetByIDs(ctx context.Context, ids []int64) ([]*service.Account, error) {
	var result []*service.Account
	for _, id := range ids {
		if acc, ok := r.accounts[id]; ok {
			result = append(result, acc)
	placeholder
placeholder
	return result, nil
placeholder
func (r *stubAccountRepo) ExistsByID(ctx context.Context, id int64) (bool, error) {
	_, ok := r.accounts[id]
	return ok, nil
placeholder
func (r *stubAccountRepo) GetByCRSAccountID(ctx context.Context, crsAccountID string) (*service.Account, error) {
	return nil, nil
placeholder
func (r *stubAccountRepo) FindByExtraField(ctx context.Context, key string, value any) ([]service.Account, error) {
	return nil, nil
placeholder
func (r *stubAccountRepo) ListCRSAccountIDs(ctx context.Context) (map[string]int64, error) {
	return map[string]int64{placeholder, nil
placeholder
func (r *stubAccountRepo) Update(ctx context.Context, account *service.Account) error { return nil placeholder
func (r *stubAccountRepo) Delete(ctx context.Context, id int64) error                 { return nil placeholder
func (r *stubAccountRepo) List(ctx context.Context, params pagination.PaginationParams) ([]service.Account, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (r *stubAccountRepo) ListWithFilters(ctx context.Context, params pagination.PaginationParams, platform, accountType, status, search string) ([]service.Account, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (r *stubAccountRepo) ListByGroup(ctx context.Context, groupID int64) ([]service.Account, error) {
	return nil, nil
placeholder
func (r *stubAccountRepo) ListActive(ctx context.Context) ([]service.Account, error) { return nil, nil placeholder
func (r *stubAccountRepo) ListByPlatform(ctx context.Context, platform string) ([]service.Account, error) {
	return r.listSchedulableByPlatform(platform), nil
placeholder
func (r *stubAccountRepo) UpdateLastUsed(ctx context.Context, id int64) error { return nil placeholder
func (r *stubAccountRepo) BatchUpdateLastUsed(ctx context.Context, updates map[int64]time.Time) error {
	return nil
placeholder
func (r *stubAccountRepo) SetError(ctx context.Context, id int64, errorMsg string) error { return nil placeholder
func (r *stubAccountRepo) ClearError(ctx context.Context, id int64) error                { return nil placeholder
func (r *stubAccountRepo) SetSchedulable(ctx context.Context, id int64, schedulable bool) error {
	return nil
placeholder
func (r *stubAccountRepo) AutoPauseExpiredAccounts(ctx context.Context, now time.Time) (int64, error) {
	return 0, nil
placeholder
func (r *stubAccountRepo) BindGroups(ctx context.Context, accountID int64, groupIDs []int64) error {
	return nil
placeholder
func (r *stubAccountRepo) ListSchedulable(ctx context.Context) ([]service.Account, error) {
	return r.listSchedulable(), nil
placeholder
func (r *stubAccountRepo) ListSchedulableByGroupID(ctx context.Context, groupID int64) ([]service.Account, error) {
	return r.listSchedulable(), nil
placeholder
func (r *stubAccountRepo) ListSchedulableByPlatform(ctx context.Context, platform string) ([]service.Account, error) {
	return r.listSchedulableByPlatform(platform), nil
placeholder
func (r *stubAccountRepo) ListSchedulableByGroupIDAndPlatform(ctx context.Context, groupID int64, platform string) ([]service.Account, error) {
	return r.listSchedulableByPlatform(platform), nil
placeholder
func (r *stubAccountRepo) ListSchedulableByPlatforms(ctx context.Context, platforms []string) ([]service.Account, error) {
	var result []service.Account
	for _, acc := range r.accounts {
		for _, platform := range platforms {
			if acc.Platform == platform && acc.IsSchedulable() {
				result = append(result, *acc)
				break
		placeholder
	placeholder
placeholder
	return result, nil
placeholder
func (r *stubAccountRepo) ListSchedulableByGroupIDAndPlatforms(ctx context.Context, groupID int64, platforms []string) ([]service.Account, error) {
	return r.ListSchedulableByPlatforms(ctx, platforms)
placeholder
func (r *stubAccountRepo) SetRateLimited(ctx context.Context, id int64, resetAt time.Time) error {
	return nil
placeholder
func (r *stubAccountRepo) SetModelRateLimit(ctx context.Context, id int64, scope string, resetAt time.Time) error {
	return nil
placeholder
func (r *stubAccountRepo) SetOverloaded(ctx context.Context, id int64, until time.Time) error {
	return nil
placeholder
func (r *stubAccountRepo) SetTempUnschedulable(ctx context.Context, id int64, until time.Time, reason string) error {
	return nil
placeholder
func (r *stubAccountRepo) ClearTempUnschedulable(ctx context.Context, id int64) error { return nil placeholder
func (r *stubAccountRepo) ClearRateLimit(ctx context.Context, id int64) error         { return nil placeholder
func (r *stubAccountRepo) ClearAntigravityQuotaScopes(ctx context.Context, id int64) error {
	return nil
placeholder
func (r *stubAccountRepo) ClearModelRateLimits(ctx context.Context, id int64) error { return nil placeholder
func (r *stubAccountRepo) UpdateSessionWindow(ctx context.Context, id int64, start, end *time.Time, status string) error {
	return nil
placeholder
func (r *stubAccountRepo) UpdateExtra(ctx context.Context, id int64, updates map[string]any) error {
	return nil
placeholder
func (r *stubAccountRepo) BulkUpdate(ctx context.Context, ids []int64, updates service.AccountBulkUpdate) (int64, error) {
	return 0, nil
placeholder

func (r *stubAccountRepo) listSchedulable() []service.Account {
	var result []service.Account
	for _, acc := range r.accounts {
		if acc.IsSchedulable() {
			result = append(result, *acc)
	placeholder
placeholder
	return result
placeholder

func (r *stubAccountRepo) listSchedulableByPlatform(platform string) []service.Account {
	var result []service.Account
	for _, acc := range r.accounts {
		if acc.Platform == platform && acc.IsSchedulable() {
			result = append(result, *acc)
	placeholder
placeholder
	return result
placeholder

type stubGroupRepo struct {
	group *service.Group
placeholder

func (r *stubGroupRepo) Create(ctx context.Context, group *service.Group) error { return nil placeholder
func (r *stubGroupRepo) GetByID(ctx context.Context, id int64) (*service.Group, error) {
	return r.group, nil
placeholder
func (r *stubGroupRepo) GetByIDLite(ctx context.Context, id int64) (*service.Group, error) {
	return r.group, nil
placeholder
func (r *stubGroupRepo) Update(ctx context.Context, group *service.Group) error { return nil placeholder
func (r *stubGroupRepo) Delete(ctx context.Context, id int64) error             { return nil placeholder
func (r *stubGroupRepo) DeleteCascade(ctx context.Context, id int64) ([]int64, error) {
	return nil, nil
placeholder
func (r *stubGroupRepo) List(ctx context.Context, params pagination.PaginationParams) ([]service.Group, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (r *stubGroupRepo) ListWithFilters(ctx context.Context, params pagination.PaginationParams, platform, status, search string, isExclusive *bool) ([]service.Group, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (r *stubGroupRepo) ListActive(ctx context.Context) ([]service.Group, error) { return nil, nil placeholder
func (r *stubGroupRepo) ListActiveByPlatform(ctx context.Context, platform string) ([]service.Group, error) {
	return nil, nil
placeholder
func (r *stubGroupRepo) ExistsByName(ctx context.Context, name string) (bool, error) {
	return false, nil
placeholder
func (r *stubGroupRepo) GetAccountCount(ctx context.Context, groupID int64) (int64, error) {
	return 0, nil
placeholder
func (r *stubGroupRepo) DeleteAccountGroupsByGroupID(ctx context.Context, groupID int64) (int64, error) {
	return 0, nil
placeholder
func (r *stubGroupRepo) GetAccountIDsByGroupIDs(ctx context.Context, groupIDs []int64) ([]int64, error) {
	return nil, nil
placeholder
func (r *stubGroupRepo) BindAccountsToGroup(ctx context.Context, groupID int64, accountIDs []int64) error {
	return nil
placeholder
func (r *stubGroupRepo) UpdateSortOrders(ctx context.Context, updates []service.GroupSortOrderUpdate) error {
	return nil
placeholder

type stubUsageLogRepo struct{placeholder

func (s *stubUsageLogRepo) Create(ctx context.Context, log *service.UsageLog) (bool, error) {
	return true, nil
placeholder
func (s *stubUsageLogRepo) GetByID(ctx context.Context, id int64) (*service.UsageLog, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) Delete(ctx context.Context, id int64) error { return nil placeholder
func (s *stubUsageLogRepo) ListByUser(ctx context.Context, userID int64, params pagination.PaginationParams) ([]service.UsageLog, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (s *stubUsageLogRepo) ListByAPIKey(ctx context.Context, apiKeyID int64, params pagination.PaginationParams) ([]service.UsageLog, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (s *stubUsageLogRepo) ListByAccount(ctx context.Context, accountID int64, params pagination.PaginationParams) ([]service.UsageLog, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (s *stubUsageLogRepo) ListByUserAndTimeRange(ctx context.Context, userID int64, startTime, endTime time.Time) ([]service.UsageLog, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (s *stubUsageLogRepo) ListByAPIKeyAndTimeRange(ctx context.Context, apiKeyID int64, startTime, endTime time.Time) ([]service.UsageLog, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (s *stubUsageLogRepo) ListByAccountAndTimeRange(ctx context.Context, accountID int64, startTime, endTime time.Time) ([]service.UsageLog, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (s *stubUsageLogRepo) ListByModelAndTimeRange(ctx context.Context, modelName string, startTime, endTime time.Time) ([]service.UsageLog, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (s *stubUsageLogRepo) GetAccountWindowStats(ctx context.Context, accountID int64, startTime time.Time) (*usagestats.AccountStats, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) GetAccountTodayStats(ctx context.Context, accountID int64) (*usagestats.AccountStats, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) GetDashboardStats(ctx context.Context) (*usagestats.DashboardStats, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) GetUsageTrendWithFilters(ctx context.Context, startTime, endTime time.Time, granularity string, userID, apiKeyID, accountID, groupID int64, model string, stream *bool, billingType *int8) ([]usagestats.TrendDataPoint, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) GetModelStatsWithFilters(ctx context.Context, startTime, endTime time.Time, userID, apiKeyID, accountID, groupID int64, stream *bool, billingType *int8) ([]usagestats.ModelStat, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) GetAPIKeyUsageTrend(ctx context.Context, startTime, endTime time.Time, granularity string, limit int) ([]usagestats.APIKeyUsageTrendPoint, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) GetUserUsageTrend(ctx context.Context, startTime, endTime time.Time, granularity string, limit int) ([]usagestats.UserUsageTrendPoint, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) GetBatchUserUsageStats(ctx context.Context, userIDs []int64, startTime, endTime time.Time) (map[int64]*usagestats.BatchUserUsageStats, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) GetBatchAPIKeyUsageStats(ctx context.Context, apiKeyIDs []int64, startTime, endTime time.Time) (map[int64]*usagestats.BatchAPIKeyUsageStats, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) GetUserDashboardStats(ctx context.Context, userID int64) (*usagestats.UserDashboardStats, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) GetAPIKeyDashboardStats(ctx context.Context, apiKeyID int64) (*usagestats.UserDashboardStats, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) GetUserUsageTrendByUserID(ctx context.Context, userID int64, startTime, endTime time.Time, granularity string) ([]usagestats.TrendDataPoint, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) GetUserModelStats(ctx context.Context, userID int64, startTime, endTime time.Time) ([]usagestats.ModelStat, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) ListWithFilters(ctx context.Context, params pagination.PaginationParams, filters usagestats.UsageLogFilters) ([]service.UsageLog, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (s *stubUsageLogRepo) GetGlobalStats(ctx context.Context, startTime, endTime time.Time) (*usagestats.UsageStats, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) GetStatsWithFilters(ctx context.Context, filters usagestats.UsageLogFilters) (*usagestats.UsageStats, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) GetAccountUsageStats(ctx context.Context, accountID int64, startTime, endTime time.Time) (*usagestats.AccountUsageStatsResponse, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) GetUserStatsAggregated(ctx context.Context, userID int64, startTime, endTime time.Time) (*usagestats.UsageStats, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) GetAPIKeyStatsAggregated(ctx context.Context, apiKeyID int64, startTime, endTime time.Time) (*usagestats.UsageStats, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) GetAccountStatsAggregated(ctx context.Context, accountID int64, startTime, endTime time.Time) (*usagestats.UsageStats, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) GetModelStatsAggregated(ctx context.Context, modelName string, startTime, endTime time.Time) (*usagestats.UsageStats, error) {
	return nil, nil
placeholder
func (s *stubUsageLogRepo) GetDailyStatsAggregated(ctx context.Context, userID int64, startTime, endTime time.Time) ([]map[string]any, error) {
	return nil, nil
placeholder

func TestSoraGatewayHandler_ChatCompletions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		RunMode: config.RunModeSimple,
		Gateway: config.GatewayConfig{
			SoraStreamMode:     "force",
			MaxAccountSwitches: 1,
			Scheduling: config.GatewaySchedulingConfig{
				LoadBatchEnabled: false,
		placeholder,
	placeholder,
		Concurrency: config.ConcurrencyConfig{PingInterval: 0placeholder,
		Sora: config.SoraConfig{
			Client: config.SoraClientConfig{
				BaseURL:             "https://sora.test",
				PollIntervalSeconds: 1,
				MaxPollAttempts:     1,
		placeholder,
	placeholder,
placeholder

	account := &service.Account{ID: 1, Platform: service.PlatformSora, Status: service.StatusActive, Schedulable: true, Concurrency: 1, Priority: 1placeholder
	accountRepo := &stubAccountRepo{accounts: map[int64]*service.Account{account.ID: accountplaceholderplaceholder
	group := &service.Group{ID: 1, Platform: service.PlatformSora, Status: service.StatusActive, Hydrated: trueplaceholder
	groupRepo := &stubGroupRepo{group: groupplaceholder

	usageLogRepo := &stubUsageLogRepo{placeholder
	deferredService := service.NewDeferredService(accountRepo, nil, 0)
	billingService := service.NewBillingService(cfg, nil)
	concurrencyService := service.NewConcurrencyService(testutil.StubConcurrencyCache{placeholder)
	billingCacheService := service.NewBillingCacheService(nil, nil, nil, cfg)
	t.Cleanup(func() {
		billingCacheService.Stop()
placeholder)

	gatewayService := service.NewGatewayService(
		accountRepo,
		groupRepo,
		usageLogRepo,
		nil,
		nil,
		nil,
		testutil.StubGatewayCache{placeholder,
		cfg,
		nil,
		concurrencyService,
		billingService,
		nil,
		billingCacheService,
		nil,
		nil,
		deferredService,
		nil,
		testutil.StubSessionLimitCache{placeholder,
		nil,
	)

	soraClient := &stubSoraClient{imageURLs: []string{"https://example.com/a.png"placeholderplaceholder
	soraGatewayService := service.NewSoraGatewayService(soraClient, nil, nil, cfg)

	handler := NewSoraGatewayHandler(gatewayService, soraGatewayService, concurrencyService, billingCacheService, cfg)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	body := `{"model":"gpt-image","messages":[{"role":"user","content":"hello"placeholder]placeholder`
	c.Request = httptest.NewRequest(http.MethodPost, "/sora/v1/chat/completions", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	apiKey := &service.APIKey{
		ID:      1,
		UserID:  1,
		Status:  service.StatusActive,
		GroupID: &group.ID,
		User:    &service.User{ID: 1, Concurrency: 1, Status: service.StatusActiveplaceholder,
		Group:   group,
placeholder
	c.Set(string(middleware.ContextKeyAPIKey), apiKey)
	c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: apiKey.UserID, Concurrency: apiKey.User.Concurrencyplaceholder)

	handler.ChatCompletions(c)

	require.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.NotEmpty(t, resp["media_url"])
placeholder

// TestSoraHandler_StreamForcing 验证 sora handler 的 stream 强制逻辑
func TestSoraHandler_StreamForcing(t *testing.T) {
	// 测试 1：stream=false 时 sjson 强制修改为 true
	body := []byte(`{"model":"sora","messages":[{"role":"user","content":"test"placeholder],"stream":falseplaceholder`)
	clientStream := gjson.GetBytes(body, "stream").Bool()
	require.False(t, clientStream)
	newBody, err := sjson.SetBytes(body, "stream", true)
placeholder
	require.True(t, gjson.GetBytes(newBody, "stream").Bool())

	// 测试 2：stream=true 时不修改
	body2 := []byte(`{"model":"sora","messages":[{"role":"user","content":"test"placeholder],"stream":trueplaceholder`)
	require.True(t, gjson.GetBytes(body2, "stream").Bool())

	// 测试 3：无 stream 字段时 gjson 返回 false（零值）
	body3 := []byte(`{"model":"sora","messages":[{"role":"user","content":"test"placeholder]placeholder`)
	require.False(t, gjson.GetBytes(body3, "stream").Bool())
placeholder

// TestSoraHandler_ValidationExtraction 验证 sora handler 中 gjson 字段校验逻辑
func TestSoraHandler_ValidationExtraction(t *testing.T) {
	// model 缺失
	body := []byte(`{"messages":[{"role":"user","content":"test"placeholder]placeholder`)
	modelResult := gjson.GetBytes(body, "model")
	require.True(t, !modelResult.Exists() || modelResult.Type != gjson.String || modelResult.String() == "")

	// model 为数字 → 类型不是 gjson.String，应被拒绝
	body1b := []byte(`{"model":123,"messages":[{"role":"user","content":"test"placeholder]placeholder`)
	modelResult1b := gjson.GetBytes(body1b, "model")
	require.True(t, modelResult1b.Exists())
	require.NotEqual(t, gjson.String, modelResult1b.Type)

	// messages 缺失
	body2 := []byte(`{"model":"sora"placeholder`)
	require.False(t, gjson.GetBytes(body2, "messages").IsArray())

	// messages 不是 JSON 数组（字符串）
	body3 := []byte(`{"model":"sora","messages":"not array"placeholder`)
	require.False(t, gjson.GetBytes(body3, "messages").IsArray())

	// messages 是对象而非数组 → IsArray 返回 false
	body4 := []byte(`{"model":"sora","messages":{placeholderplaceholder`)
	require.False(t, gjson.GetBytes(body4, "messages").IsArray())

	// messages 是空数组 → IsArray 为 true 但 len==0，应被拒绝
	body5 := []byte(`{"model":"sora","messages":[]placeholder`)
	msgsResult := gjson.GetBytes(body5, "messages")
	require.True(t, msgsResult.IsArray())
	require.Equal(t, 0, len(msgsResult.Array()))

	// 非法 JSON 被 gjson.ValidBytes 拦截
	require.False(t, gjson.ValidBytes([]byte(`{invalid`)))
placeholder

// TestGenerateOpenAISessionHash_WithBody 验证 generateOpenAISessionHash 的 body/header 解析逻辑
func TestGenerateOpenAISessionHash_WithBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 从 body 提取 prompt_cache_key
	body := []byte(`{"model":"sora","prompt_cache_key":"session-abc"placeholder`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", nil)

	hash := generateOpenAISessionHash(c, body)
	require.NotEmpty(t, hash)

	// 无 prompt_cache_key 且无 header → 空 hash
	body2 := []byte(`{"model":"sora"placeholder`)
	hash2 := generateOpenAISessionHash(c, body2)
	require.Empty(t, hash2)

	// header 优先于 body
	c.Request.Header.Set("session_id", "from-header")
	hash3 := generateOpenAISessionHash(c, body)
	require.NotEmpty(t, hash3)
	require.NotEqual(t, hash, hash3) // 不同来源应产生不同 hash
placeholder
