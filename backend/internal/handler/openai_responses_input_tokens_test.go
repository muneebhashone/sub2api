package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestResponsesInputTokensIsClassifiedAsTokenCountRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses/input_tokens", nil)

	require.True(t, isCountTokensRequest(c))
	require.True(t, isTokenCountRequestPath("/responses/input_tokens"))
	require.False(t, isTokenCountRequestPath("/v1/responses"))
placeholder

func TestCountTokensAccountSlot_CancellationStopsBeforeForward(t *testing.T) {
	gin.SetMode(gin.TestMode)

	for _, tt := range []struct {
		name      string
		anthropic bool
placeholder{
		{name: "responses input tokens", anthropic: falseplaceholder,
		{name: "anthropic count tokens", anthropic: trueplaceholder,
placeholder {
		t.Run(tt.name, func(t *testing.T) {
			cache := &concurrencyCacheMock{
				acquireAccountSlotFn: func(context.Context, int64, int, string) (bool, error) {
					return false, nil
			placeholder,
		placeholder
			h := &OpenAIGatewayHandler{
				gatewayService:    &service.OpenAIGatewayService{placeholder,
				concurrencyHelper: NewConcurrencyHelper(service.NewConcurrencyService(cache), SSEPingFormatNone, time.Second),
		placeholder
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/count_tokens", nil).WithContext(ctx)
			groupID := int64(41)
			selection := &service.AccountSelectionResult{
				Account: &service.Account{ID: 42, Platform: service.PlatformOpenAIplaceholder,
				WaitPlan: &service.AccountWaitPlan{
					AccountID:      42,
					MaxConcurrency: 1,
					MaxWaiting:     1,
					Timeout:        time.Second,
			placeholder,
		placeholder

			release, acquired := h.acquireCountTokensAccountSlot(c, &groupID, "", selection, tt.anthropic, zap.NewNop())
			forwarded := false
			if acquired {
				forwarded = true
				if release != nil {
					release()
			placeholder
		placeholder

			require.False(t, acquired)
			require.Nil(t, release)
			require.False(t, forwarded, "a canceled WaitPlan must stop before count-token forwarding")
			require.Zero(t, atomic.LoadInt32(&cache.releaseAccountCalled))
	placeholder)
placeholder
placeholder

func TestCountTokensAccountSlot_SelectionReleaseRunsExactlyOnce(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var released atomic.Int32
	ctx, cancel := context.WithCancel(context.Background())
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses/input_tokens", nil).WithContext(ctx)
	h := &OpenAIGatewayHandler{gatewayService: &service.OpenAIGatewayService{placeholderplaceholder
	groupID := int64(51)
	selection := &service.AccountSelectionResult{
		Account:     &service.Account{ID: 52, Platform: service.PlatformOpenAIplaceholder,
		Acquired:    true,
		ReleaseFunc: func() { released.Add(1) placeholder,
placeholder

	release, acquired := h.acquireCountTokensAccountSlot(c, &groupID, "", selection, false, zap.NewNop())
	require.True(t, acquired)
	require.NotNil(t, release)
	release()
	cancel()
	require.Eventually(t, func() bool { return released.Load() == 1 placeholder, time.Second, 10*time.Millisecond)
	require.Equal(t, int32(1), released.Load())
placeholder
