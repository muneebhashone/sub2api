//go:build unit

package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	middleware "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type countingGatewaySchedulerCache struct {
	*fakeSchedulerCache
	snapshotCalls atomic.Int64
placeholder

func (c *countingGatewaySchedulerCache) GetSnapshot(ctx context.Context, bucket service.SchedulerBucket) ([]*service.Account, bool, error) {
	c.snapshotCalls.Add(1)
	return c.fakeSchedulerCache.GetSnapshot(ctx, bucket)
placeholder

func TestGatewayHandlerPreCancelledCompatibleRequestsDoNotSelectAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	groupID := int64(9100)
	group := &service.Group{ID: groupID, Hydrated: true, Platform: service.PlatformAnthropic, Status: service.StatusActiveplaceholder
	account := &service.Account{
		ID: 9101, Platform: service.PlatformAnthropic, Type: service.AccountTypeAPIKey,
		Status: service.StatusActive, Schedulable: true, Concurrency: 1,
		AccountGroups: []service.AccountGroup{{AccountID: 9101, GroupID: groupIDplaceholderplaceholder,
placeholder
	schedulerCache := &countingGatewaySchedulerCache{fakeSchedulerCache: &fakeSchedulerCache{accounts: []*service.Account{accountplaceholderplaceholderplaceholder
	schedulerSnapshot := service.NewSchedulerSnapshotService(schedulerCache, nil, nil, nil, nil)
	gatewayService := service.NewGatewayService(
		nil, &fakeGroupRepo{group: groupplaceholder, nil, nil, nil, nil, nil, nil, nil,
		schedulerSnapshot, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)
	cfg := &config.Config{RunMode: config.RunModeSimpleplaceholder
	billingCacheService := service.NewBillingCacheService(nil, nil, nil, nil, nil, nil, cfg, nil)
	t.Cleanup(billingCacheService.Stop)
	h := &GatewayHandler{
		gatewayService:      gatewayService,
		billingCacheService: billingCacheService,
		concurrencyHelper:   NewConcurrencyHelper(service.NewConcurrencyService(&fakeConcurrencyCache{placeholder), SSEPingFormatClaude, 0),
		maxAccountSwitches:  1,
		cfg:                 cfg,
placeholder
	apiKey := &service.APIKey{
		ID: 9102, UserID: 9103, GroupID: &groupID, Group: group, Status: service.StatusActive,
		User: &service.User{ID: 9103, Concurrency: 10, Balance: 100placeholder,
placeholder

	tests := []struct {
		name string
		path string
		body string
		call func(*gin.Context)
placeholder{
		{
			name: "responses", path: "/v1/responses", body: `{"model":"claude-test","input":"hello","stream":falseplaceholder`,
			call: h.Responses,
	placeholder,
		{
			name: "chat completions", path: "/v1/chat/completions", body: `{"model":"claude-test","messages":[{"role":"user","content":"hello"placeholder],"stream":falseplaceholder`,
			call: h.ChatCompletions,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedulerCache.snapshotCalls.Store(0)
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			ctx = context.WithValue(ctx, ctxkey.Group, group)
			req := httptest.NewRequest(http.MethodPost, tt.path, bytes.NewBufferString(tt.body)).WithContext(ctx)
			req.Header.Set("Content-Type", "application/json")
			c.Request = req
			c.Set(string(middleware.ContextKeyAPIKey), apiKey)
			c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: apiKey.UserID, Concurrency: 10placeholder)

			tt.call(c)

			require.Zero(t, schedulerCache.snapshotCalls.Load(), "a cancelled request must stop before the account selector")
			_, selected := c.Get(opsAccountIDKey)
			require.False(t, selected)
	placeholder)
placeholder
placeholder
