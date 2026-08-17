//go:build unit

package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

// upsertCapturingQuotaRepo 实现 service.UserPlatformQuotaRepository，捕获 UpsertForUser 调用。
type upsertCapturingQuotaRepo struct {
	service.UserPlatformQuotaRepository
	listRecords []service.UserPlatformQuotaRecord
	listErr     error
	upsertCalls []upsertCall
	upsertErr   error
	resetCalls  []resetCall
	resetErr    error
placeholder

type upsertCall struct {
	userID  int64
	records []service.UserPlatformQuotaRecord
placeholder
type resetCall struct {
	userID   int64
	platform string
	window   string
	newStart time.Time
placeholder

func (r *upsertCapturingQuotaRepo) ListByUser(_ context.Context, _ int64) ([]service.UserPlatformQuotaRecord, error) {
	return r.listRecords, r.listErr
placeholder
func (r *upsertCapturingQuotaRepo) UpsertForUser(_ context.Context, userID int64, records []service.UserPlatformQuotaRecord) error {
	cloned := make([]service.UserPlatformQuotaRecord, len(records))
	copy(cloned, records)
	r.upsertCalls = append(r.upsertCalls, upsertCall{userID: userID, records: clonedplaceholder)
	return r.upsertErr
placeholder
func (r *upsertCapturingQuotaRepo) ResetExpiredWindow(_ context.Context, userID int64, platform string, window string, newStart time.Time) error {
	r.resetCalls = append(r.resetCalls, resetCall{userID, platform, window, newStartplaceholder)
	return r.resetErr
placeholder

// billingCacheStub 实现 service.BillingCache 中本测试关心的 Delete 方法；其他方法 panic。
type billingCacheStub struct {
	service.BillingCache
	deleteCalls []deleteCall
	deleteErr   error
placeholder

type deleteCall struct {
	userID   int64
	platform string
placeholder

func (b *billingCacheStub) DeleteUserPlatformQuotaCache(_ context.Context, userID int64, platform string) error {
	b.deleteCalls = append(b.deleteCalls, deleteCall{userID, platformplaceholder)
	return b.deleteErr
placeholder

func buildTestHandler(repo service.UserPlatformQuotaRepository, cache service.BillingCache) *UserHandler {
	return &UserHandler{
		userPlatformQuotaRepo: repo,
		billingCache:          cache,
		adminService:          newStubAdminService(),
placeholder
placeholder

func putReq(t *testing.T, body string) (*gin.Context, *httptest.ResponseRecorder) {
placeholder
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodPut, "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "42"placeholderplaceholder
	return c, w
placeholder

func TestUpdateUserPlatformQuotas_Success(t *testing.T) {
	repo := &upsertCapturingQuotaRepo{placeholder
	cache := &billingCacheStub{placeholder
	h := buildTestHandler(repo, cache)

	body := `{"quotas":[
		{"platform":"anthropic","daily_limit_usd":10.0,"weekly_limit_usd":null,"monthly_limit_usd":100.0placeholder,
		{"platform":"openai","daily_limit_usd":80.0,"weekly_limit_usd":300.0,"monthly_limit_usd":nullplaceholder,
		{"platform":"gemini","daily_limit_usd":null,"weekly_limit_usd":null,"monthly_limit_usd":nullplaceholder,
		{"platform":"antigravity","daily_limit_usd":null,"weekly_limit_usd":null,"monthly_limit_usd":nullplaceholder,
		{"platform":"grok","daily_limit_usd":null,"weekly_limit_usd":null,"monthly_limit_usd":nullplaceholder
	]placeholder`
	c, w := putReq(t, body)
	h.UpdateUserPlatformQuotas(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
placeholder
	if len(repo.upsertCalls) != 1 {
		t.Fatalf("UpsertForUser should be called once, got %d", len(repo.upsertCalls))
placeholder
	// upsert 记录数 = 请求体中给出的平台数（未给出的平台不落库）。
	if repo.upsertCalls[0].userID != 42 || len(repo.upsertCalls[0].records) != 5 {
		t.Errorf("unexpected upsert call: %+v", repo.upsertCalls[0])
placeholder
	// 缓存失效：按全部允许平台统一失效（含 kimi/zhipu/deepseek）。
	if len(cache.deleteCalls) != len(service.AllowedQuotaPlatforms) {
		t.Errorf("expected %d cache delete calls, got %d: %+v", len(service.AllowedQuotaPlatforms), len(cache.deleteCalls), cache.deleteCalls)
placeholder
placeholder

func TestUpdateUserPlatformQuotas_RejectsDuplicatePlatform(t *testing.T) {
	h := buildTestHandler(&upsertCapturingQuotaRepo{placeholder, &billingCacheStub{placeholder)
	body := `{"quotas":[
		{"platform":"anthropic","daily_limit_usd":1placeholder,
		{"platform":"anthropic","daily_limit_usd":2placeholder
	]placeholder`
	c, w := putReq(t, body)
	h.UpdateUserPlatformQuotas(c)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
placeholder
placeholder

func TestUpdateUserPlatformQuotas_RejectsInvalidPlatform(t *testing.T) {
	h := buildTestHandler(&upsertCapturingQuotaRepo{placeholder, &billingCacheStub{placeholder)
	body := `{"quotas":[{"platform":"unknown","daily_limit_usd":1placeholder]placeholder`
	c, w := putReq(t, body)
	h.UpdateUserPlatformQuotas(c)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
placeholder
placeholder

func TestUpdateUserPlatformQuotas_RejectsNegativeLimit(t *testing.T) {
	h := buildTestHandler(&upsertCapturingQuotaRepo{placeholder, &billingCacheStub{placeholder)
	body := `{"quotas":[{"platform":"anthropic","daily_limit_usd":-1placeholder]placeholder`
	c, w := putReq(t, body)
	h.UpdateUserPlatformQuotas(c)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
placeholder
placeholder

func TestUpdateUserPlatformQuotas_RejectsTooManyEntries(t *testing.T) {
	h := buildTestHandler(&upsertCapturingQuotaRepo{placeholder, &billingCacheStub{placeholder)
	body := `{"quotas":[
		{"platform":"anthropic"placeholder,{"platform":"openai"placeholder,{"platform":"gemini"placeholder,{"platform":"antigravity"placeholder,{"platform":"grok"placeholder,{"platform":"anthropic"placeholder
	]placeholder`
	c, w := putReq(t, body)
	h.UpdateUserPlatformQuotas(c)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
placeholder
placeholder

func TestUpdateUserPlatformQuotas_ReturnsLatestState(t *testing.T) {
	repo := &upsertCapturingQuotaRepo{
		listRecords: []service.UserPlatformQuotaRecord{
			{UserID: 42, Platform: "anthropic"placeholder,
	placeholder,
placeholder
	cache := &billingCacheStub{placeholder
	h := buildTestHandler(repo, cache)

	body := `{"quotas":[{"platform":"anthropic","daily_limit_usd":10placeholder]placeholder`
	c, w := putReq(t, body)
	h.UpdateUserPlatformQuotas(c)
	if !strings.Contains(w.Body.String(), `"platform_quotas"`) {
		t.Errorf("response should contain platform_quotas array: %s", w.Body.String())
placeholder
placeholder

// ───────── T4: Reset 测试 ─────────

func postReq(t *testing.T, body string) (*gin.Context, *httptest.ResponseRecorder) {
placeholder
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "42"placeholderplaceholder
	return c, w
placeholder

func TestResetUserPlatformQuotaWindow_Success(t *testing.T) {
	repo := &upsertCapturingQuotaRepo{placeholder
	cache := &billingCacheStub{placeholder
	h := buildTestHandler(repo, cache)
	body := `{"platform":"anthropic","window":"daily"placeholder`
	c, w := postReq(t, body)
	h.ResetUserPlatformQuotaWindow(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
placeholder
	if len(repo.resetCalls) != 1 {
		t.Fatalf("ResetExpiredWindow should be called once, got %d", len(repo.resetCalls))
placeholder
	if repo.resetCalls[0].userID != 42 ||
		repo.resetCalls[0].platform != "anthropic" ||
		repo.resetCalls[0].window != "daily" {
		t.Errorf("unexpected reset call: %+v", repo.resetCalls[0])
placeholder
	if len(cache.deleteCalls) != 1 ||
		cache.deleteCalls[0].userID != 42 ||
		cache.deleteCalls[0].platform != "anthropic" {
		t.Errorf("expected 1 cache delete for anthropic, got %+v", cache.deleteCalls)
placeholder
placeholder

func TestResetUserPlatformQuotaWindow_RejectsInvalidWindow(t *testing.T) {
	h := buildTestHandler(&upsertCapturingQuotaRepo{placeholder, &billingCacheStub{placeholder)
	c, w := postReq(t, `{"platform":"anthropic","window":"yearly"placeholder`)
	h.ResetUserPlatformQuotaWindow(c)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
placeholder
placeholder

func TestResetUserPlatformQuotaWindow_RejectsInvalidPlatform(t *testing.T) {
	h := buildTestHandler(&upsertCapturingQuotaRepo{placeholder, &billingCacheStub{placeholder)
	c, w := postReq(t, `{"platform":"unknown","window":"daily"placeholder`)
	h.ResetUserPlatformQuotaWindow(c)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
placeholder
placeholder

func TestResetUserPlatformQuotaWindow_NotFound(t *testing.T) {
	// handler 检查 service.ErrUserPlatformQuotaNotFound（由 adapter 包装而来）
	repo := &upsertCapturingQuotaRepo{resetErr: service.ErrUserPlatformQuotaNotFoundplaceholder
	h := buildTestHandler(repo, &billingCacheStub{placeholder)
	c, w := postReq(t, `{"platform":"anthropic","window":"daily"placeholder`)
	h.ResetUserPlatformQuotaWindow(c)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d: %s", w.Code, w.Body.String())
placeholder
placeholder

func TestUpdateUserPlatformQuotas_JSONErrorOnRepoFailure(t *testing.T) {
	repo := &upsertCapturingQuotaRepo{upsertErr: errors.New("db down")placeholder
	cache := &billingCacheStub{placeholder
	h := buildTestHandler(repo, cache)
	body := `{"quotas":[{"platform":"anthropic","daily_limit_usd":10placeholder]placeholder`
	c, w := putReq(t, body)
	h.UpdateUserPlatformQuotas(c)
	if w.Code < 500 {
		t.Errorf("expected 5xx, got %d", w.Code)
placeholder
	// 返回 JSON 错误响应
	var body2 map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body2); err != nil {
		t.Errorf("expected JSON error body, got: %s", w.Body.String())
placeholder
placeholder

func TestUpdateUserPlatformQuotas_UserNotFound(t *testing.T) {
	repo := &upsertCapturingQuotaRepo{placeholder
	cache := &billingCacheStub{placeholder
	adminSvc := newStubAdminService()
	adminSvc.getUserErr = service.ErrUserNotFound
	h := &UserHandler{
		userPlatformQuotaRepo: repo,
		billingCache:          cache,
		adminService:          adminSvc,
placeholder
	body := `{"quotas":[{"platform":"anthropic","daily_limit_usd":10placeholder]placeholder`
	c, w := putReq(t, body)
	h.UpdateUserPlatformQuotas(c)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 when user not found, got %d: %s", w.Code, w.Body.String())
placeholder
placeholder

func TestResetUserPlatformQuotaWindow_UserNotFound(t *testing.T) {
	repo := &upsertCapturingQuotaRepo{placeholder
	cache := &billingCacheStub{placeholder
	adminSvc := newStubAdminService()
	adminSvc.getUserErr = service.ErrUserNotFound
	h := &UserHandler{
		userPlatformQuotaRepo: repo,
		billingCache:          cache,
		adminService:          adminSvc,
placeholder
	c, w := postReq(t, `{"platform":"anthropic","window":"daily"placeholder`)
	h.ResetUserPlatformQuotaWindow(c)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 when user not found, got %d: %s", w.Code, w.Body.String())
placeholder
placeholder
