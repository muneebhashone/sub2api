package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// 非法 service_tier 必须在两个 OpenAI 端点（/v1/responses、/v1/chat/completions）
// 上以 OpenAI 兼容错误结构返回 HTTP 400。这些用例在 handler 的 service_tier
// 校验处短路，不会进入账号选择/重试。
//
// 合法值（fast/priority/flex/auto/default/scale）与省略/null 的接受语义由
// service 层纯校验函数 TestValidateOpenAIServiceTierField 覆盖，避免 handler
// 测试走入真实账号选择/重试路径。

func newServiceTierHandlerTest(t *testing.T) *OpenAIGatewayHandler {
placeholder
	return &OpenAIGatewayHandler{
		gatewayService:      &service.OpenAIGatewayService{placeholder,
		billingCacheService: service.NewBillingCacheService(nil, nil, nil, nil, nil, nil, &config.Config{RunMode: config.RunModeSimpleplaceholder, nil),
		apiKeyService:       &service.APIKeyService{placeholder,
		concurrencyHelper: &ConcurrencyHelper{concurrencyService: service.NewConcurrencyService(
			&helperConcurrencyCacheStub{userSeq: []bool{trueplaceholderplaceholder,
		)placeholder,
		cfg:          &config.Config{placeholder,
		imageLimiter: &imageConcurrencyLimiter{placeholder,
placeholder
placeholder

func runOpenAIHandlerServiceTierTest(t *testing.T, path, body string, handler func(h *OpenAIGatewayHandler, c *gin.Context)) *httptest.ResponseRecorder {
placeholder
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	groupID := int64(6401)
	userID := int64(6402)
	c.Set(string(middleware2.ContextKeyAPIKey), &service.APIKey{
		ID:      6403,
		GroupID: &groupID,
		Group: &service.Group{
			ID:       groupID,
			Platform: service.PlatformOpenAI,
	placeholder,
		User: &service.User{ID: userID, Status: service.StatusActiveplaceholder,
placeholder)
	c.Set(string(middleware2.ContextKeyUser), middleware2.AuthSubject{UserID: userID, Concurrency: 1placeholder)

	handler(newServiceTierHandlerTest(t), c)
	return rec
placeholder

func TestOpenAIGatewayHandlerResponses_InvalidServiceTierRejected400(t *testing.T) {
	for _, body := range []string{
		`{"model":"gpt-5.5","input":"hi","service_tier":"turbo"placeholder`,
		`{"model":"gpt-5.5","input":"hi","service_tier":"SPEED"placeholder`,
		`{"model":"gpt-5.5","input":"hi","service_tier":""placeholder`,
		`{"model":"gpt-5.5","input":"hi","service_tier":placeholder`,
		`{"model":"gpt-5.5","input":"hi","service_tier":{placeholderplaceholder`,
placeholder {
		rec := runOpenAIHandlerServiceTierTest(t, "/v1/responses", body, func(h *OpenAIGatewayHandler, c *gin.Context) {
			h.Responses(c)
	placeholder)
		require.Equal(t, http.StatusBadRequest, rec.Code, "body=%s", body)
		require.Contains(t, rec.Body.String(), "invalid_request_error", "body=%s", body)
		require.Contains(t, rec.Body.String(), "invalid service_tier", "body=%s", body)
placeholder
placeholder

func TestOpenAIGatewayHandlerChatCompletions_InvalidServiceTierRejected400(t *testing.T) {
	for _, body := range []string{
		`{"model":"gpt-5.5","messages":[{"role":"user","content":"hi"placeholder],"service_tier":"turbo"placeholder`,
		`{"model":"gpt-5.5","messages":[{"role":"user","content":"hi"placeholder],"service_tier":"ultra"placeholder`,
		`{"model":"gpt-5.5","messages":[{"role":"user","content":"hi"placeholder],"service_tier":""placeholder`,
		`{"model":"gpt-5.5","messages":[{"role":"user","content":"hi"placeholder],"service_tier":["priority"]placeholder`,
placeholder {
		rec := runOpenAIHandlerServiceTierTest(t, "/v1/chat/completions", body, func(h *OpenAIGatewayHandler, c *gin.Context) {
			h.ChatCompletions(c)
	placeholder)
		require.Equal(t, http.StatusBadRequest, rec.Code, "body=%s", body)
		require.Contains(t, rec.Body.String(), "invalid_request_error", "body=%s", body)
		require.Contains(t, rec.Body.String(), "invalid service_tier", "body=%s", body)
placeholder
placeholder
