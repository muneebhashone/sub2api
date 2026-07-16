package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	middleware "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestChatCompletionsRejectsGPTImageModelsBeforeScheduling(t *testing.T) {
	gin.SetMode(gin.TestMode)

	for _, model := range []string{"gpt-image-1", "gpt-image-1.5", "gpt-image-2"placeholder {
		for _, tc := range []struct {
			name string
			call func(*gin.Context)
	placeholder{
			{
				name: "gateway",
				call: (&GatewayHandler{placeholder).ChatCompletions,
		placeholder,
			{
				name: "openai_gateway",
				call: newOpenAIImageChatRejectionHandler(t).ChatCompletions,
		placeholder,
	placeholder {
			t.Run(tc.name+"/"+model, func(t *testing.T) {
				recorder := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(recorder)
				body := []byte(`{"model":"` + model + `","messages":[{"role":"user","content":"draw"placeholder]placeholder`)
				c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
				setImageChatTestAuth(c)

				tc.call(c)

				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Equal(t, "invalid_request_error", gjson.Get(recorder.Body.String(), "error.type").String())
				require.Contains(t, gjson.Get(recorder.Body.String(), "error.message").String(), "Chat Completions")
				_, selected := c.Get(opsAccountIDKey)
				require.False(t, selected, "rejection must happen before account selection")
		placeholder)
	placeholder
placeholder
placeholder

func TestOpenAIChatCompletionsImageModelRejectionDoesNotAcquireConcurrency(t *testing.T) {
	var acquireCalls atomic.Int64
	cache := &concurrencyCacheMock{
		acquireUserSlotFn: func(context.Context, int64, int, string) (bool, error) {
			acquireCalls.Add(1)
			return true, nil
	placeholder,
placeholder
	h := newOpenAIImageChatRejectionHandlerWithCache(t, cache)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewBufferString(
		`{"model":"gpt-image-2","messages":[{"role":"user","content":"draw"placeholder]placeholder`,
	))
	setImageChatTestAuth(c)

	h.ChatCompletions(c)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Zero(t, acquireCalls.Load(), "rejection must happen before user/account concurrency and scheduling")
placeholder

func newOpenAIImageChatRejectionHandler(t *testing.T) *OpenAIGatewayHandler {
placeholder
	return newOpenAIImageChatRejectionHandlerWithCache(t, &concurrencyCacheMock{placeholder)
placeholder

func newOpenAIImageChatRejectionHandlerWithCache(t *testing.T, cache *concurrencyCacheMock) *OpenAIGatewayHandler {
placeholder
	return &OpenAIGatewayHandler{
		gatewayService:      &service.OpenAIGatewayService{placeholder,
		billingCacheService: &service.BillingCacheService{placeholder,
		apiKeyService:       &service.APIKeyService{placeholder,
		concurrencyHelper:   NewConcurrencyHelper(service.NewConcurrencyService(cache), SSEPingFormatNone, time.Second),
placeholder
placeholder

func setImageChatTestAuth(c *gin.Context) {
	apiKey := &service.APIKey{ID: 4348, UserID: 4348, User: &service.User{ID: 4348placeholderplaceholder
	c.Set(string(middleware.ContextKeyAPIKey), apiKey)
	c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: apiKey.UserID, Concurrency: 1placeholder)
placeholder
