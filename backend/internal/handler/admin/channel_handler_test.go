//go:build unit

package admin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func float64Ptr(v float64) *float64 { return &v placeholder
func intPtr(v int) *int             { return &v placeholder

// ---------------------------------------------------------------------------
// 1. channelToResponse
// ---------------------------------------------------------------------------

func TestChannelToResponse_NilInput(t *testing.T) {
	require.Nil(t, channelToResponse(nil))
placeholder

func TestChannelToResponse_FullChannel(t *testing.T) {
	now := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)
	ch := &service.Channel{
		ID:                 42,
		Name:               "test-channel",
		Description:        "desc",
		Status:             "active",
		BillingModelSource: "upstream",
		RestrictModels:     true,
		CreatedAt:          now,
		UpdatedAt:          now.Add(time.Hour),
		GroupIDs:           []int64{1, 2, 3placeholder,
		ModelPricing: []service.ChannelModelPricing{
			{
				ID:              10,
				Platform:        "openai",
				Models:          []string{"gpt-4"placeholder,
				BillingMode:     service.BillingModeToken,
				InputPrice:      float64Ptr(0.01),
				OutputPrice:     float64Ptr(0.03),
				CacheWritePrice: float64Ptr(0.005),
				CacheReadPrice:  float64Ptr(0.002),
				PerRequestPrice: float64Ptr(0.5),
		placeholder,
	placeholder,
		ModelMapping: map[string]map[string]string{
			"anthropic": {"claude-3-haiku": "claude-haiku-3"placeholder,
	placeholder,
placeholder

	resp := channelToResponse(ch)
	require.NotNil(t, resp)
	require.Equal(t, int64(42), resp.ID)
	require.Equal(t, "test-channel", resp.Name)
	require.Equal(t, "desc", resp.Description)
	require.Equal(t, "active", resp.Status)
	require.Equal(t, "upstream", resp.BillingModelSource)
	require.True(t, resp.RestrictModels)
	require.Equal(t, []int64{1, 2, 3placeholder, resp.GroupIDs)
	require.Equal(t, "2025-06-01T12:00:00Z", resp.CreatedAt)
	require.Equal(t, "2025-06-01T13:00:00Z", resp.UpdatedAt)

	// model mapping
	require.Len(t, resp.ModelMapping, 1)
	require.Equal(t, "claude-haiku-3", resp.ModelMapping["anthropic"]["claude-3-haiku"])

	// pricing
	require.Len(t, resp.ModelPricing, 1)
	p := resp.ModelPricing[0]
	require.Equal(t, int64(10), p.ID)
	require.Equal(t, "openai", p.Platform)
	require.Equal(t, []string{"gpt-4"placeholder, p.Models)
	require.Equal(t, "token", p.BillingMode)
	require.Equal(t, float64Ptr(0.01), p.InputPrice)
	require.Equal(t, float64Ptr(0.03), p.OutputPrice)
	require.Equal(t, float64Ptr(0.005), p.CacheWritePrice)
	require.Equal(t, float64Ptr(0.002), p.CacheReadPrice)
	require.Equal(t, float64Ptr(0.5), p.PerRequestPrice)
	require.Empty(t, p.Intervals)
placeholder

func TestChannelToResponse_EmptyDefaults(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	ch := &service.Channel{
		ID:                 1,
		Name:               "ch",
		BillingModelSource: service.BillingModelSourceChannelMapped,
		CreatedAt:          now,
		UpdatedAt:          now,
		GroupIDs:           nil,
		ModelMapping:       nil,
		ModelPricing: []service.ChannelModelPricing{
			{
				Platform:    "",
				BillingMode: "",
				Models:      []string{"m1"placeholder,
		placeholder,
	placeholder,
placeholder

	// handler 层 channelToResponse 现在是纯透传：BillingModelSource 的空值兜底
	// 已下放到 service 层（Create/GetByID/List/Update/ListAvailable 出口统一处理），
	// 因此这里构造 fixture 时直接传入归一化后的值。
	resp := channelToResponse(ch)
	require.Equal(t, "channel_mapped", resp.BillingModelSource)
	require.NotNil(t, resp.GroupIDs)
	require.Empty(t, resp.GroupIDs)
	require.NotNil(t, resp.ModelMapping)
	require.Empty(t, resp.ModelMapping)

	require.Len(t, resp.ModelPricing, 1)
	require.Equal(t, "anthropic", resp.ModelPricing[0].Platform)
	require.Equal(t, "token", resp.ModelPricing[0].BillingMode)
placeholder

func TestChannelToResponse_BillingModelSourcePassthrough(t *testing.T) {
	// handler 不再兜底 BillingModelSource：空值应原样透传（由 service 层负责默认回填）。
	ch := &service.Channel{
		ID:                 1,
		Name:               "ch",
		BillingModelSource: "",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
placeholder
	resp := channelToResponse(ch)
	require.Equal(t, "", resp.BillingModelSource, "handler 应纯透传，默认值由 service.normalizeBillingModelSource 负责")
placeholder

func TestChannelToResponse_NilModels(t *testing.T) {
	now := time.Now()
	ch := &service.Channel{
		ID:        1,
		Name:      "ch",
		CreatedAt: now,
		UpdatedAt: now,
		ModelPricing: []service.ChannelModelPricing{
			{
				Models: nil,
		placeholder,
	placeholder,
placeholder

	resp := channelToResponse(ch)
	require.Len(t, resp.ModelPricing, 1)
	require.NotNil(t, resp.ModelPricing[0].Models)
	require.Empty(t, resp.ModelPricing[0].Models)
placeholder

func TestChannelToResponse_WithIntervals(t *testing.T) {
	now := time.Now()
	ch := &service.Channel{
		ID:        1,
		Name:      "ch",
		CreatedAt: now,
		UpdatedAt: now,
		ModelPricing: []service.ChannelModelPricing{
			{
				Models:      []string{"m1"placeholder,
				BillingMode: service.BillingModePerRequest,
				Intervals: []service.PricingInterval{
					{
						ID:              100,
						MinTokens:       0,
						MaxTokens:       intPtr(1000),
						TierLabel:       "1K",
						InputPrice:      float64Ptr(0.01),
						OutputPrice:     float64Ptr(0.02),
						CacheWritePrice: float64Ptr(0.003),
						CacheReadPrice:  float64Ptr(0.001),
						PerRequestPrice: float64Ptr(0.1),
						SortOrder:       1,
				placeholder,
					{
						ID:        101,
						MinTokens: 1000,
						MaxTokens: nil,
						TierLabel: "unlimited",
						SortOrder: 2,
				placeholder,
			placeholder,
		placeholder,
	placeholder,
placeholder

	resp := channelToResponse(ch)
	require.Len(t, resp.ModelPricing, 1)
	intervals := resp.ModelPricing[0].Intervals
	require.Len(t, intervals, 2)

	iv0 := intervals[0]
	require.Equal(t, int64(100), iv0.ID)
	require.Equal(t, 0, iv0.MinTokens)
	require.Equal(t, intPtr(1000), iv0.MaxTokens)
	require.Equal(t, "1K", iv0.TierLabel)
	require.Equal(t, float64Ptr(0.01), iv0.InputPrice)
	require.Equal(t, float64Ptr(0.02), iv0.OutputPrice)
	require.Equal(t, float64Ptr(0.003), iv0.CacheWritePrice)
	require.Equal(t, float64Ptr(0.001), iv0.CacheReadPrice)
	require.Equal(t, float64Ptr(0.1), iv0.PerRequestPrice)
	require.Equal(t, 1, iv0.SortOrder)

	iv1 := intervals[1]
	require.Equal(t, int64(101), iv1.ID)
	require.Equal(t, 1000, iv1.MinTokens)
	require.Nil(t, iv1.MaxTokens)
	require.Equal(t, "unlimited", iv1.TierLabel)
	require.Equal(t, 2, iv1.SortOrder)
placeholder

func TestChannelToResponse_MultipleEntries(t *testing.T) {
	now := time.Now()
	ch := &service.Channel{
		ID:        1,
		Name:      "multi",
		CreatedAt: now,
		UpdatedAt: now,
		ModelPricing: []service.ChannelModelPricing{
			{
				ID:          1,
				Platform:    "anthropic",
				Models:      []string{"claude-sonnet-4"placeholder,
				BillingMode: service.BillingModeToken,
				InputPrice:  float64Ptr(0.003),
				OutputPrice: float64Ptr(0.015),
		placeholder,
			{
				ID:              2,
				Platform:        "openai",
				Models:          []string{"gpt-4", "gpt-4o"placeholder,
				BillingMode:     service.BillingModePerRequest,
				PerRequestPrice: float64Ptr(1.0),
		placeholder,
			{
				ID:               3,
				Platform:         "gemini",
				Models:           []string{"gemini-2.5-pro"placeholder,
				BillingMode:      service.BillingModeImage,
				ImageOutputPrice: float64Ptr(0.05),
				PerRequestPrice:  float64Ptr(0.2),
		placeholder,
	placeholder,
placeholder

	resp := channelToResponse(ch)
	require.Len(t, resp.ModelPricing, 3)

	require.Equal(t, int64(1), resp.ModelPricing[0].ID)
	require.Equal(t, "anthropic", resp.ModelPricing[0].Platform)
	require.Equal(t, []string{"claude-sonnet-4"placeholder, resp.ModelPricing[0].Models)
	require.Equal(t, "token", resp.ModelPricing[0].BillingMode)

	require.Equal(t, int64(2), resp.ModelPricing[1].ID)
	require.Equal(t, "openai", resp.ModelPricing[1].Platform)
	require.Equal(t, []string{"gpt-4", "gpt-4o"placeholder, resp.ModelPricing[1].Models)
	require.Equal(t, "per_request", resp.ModelPricing[1].BillingMode)

	require.Equal(t, int64(3), resp.ModelPricing[2].ID)
	require.Equal(t, "gemini", resp.ModelPricing[2].Platform)
	require.Equal(t, []string{"gemini-2.5-pro"placeholder, resp.ModelPricing[2].Models)
	require.Equal(t, "image", resp.ModelPricing[2].BillingMode)
	require.Equal(t, float64Ptr(0.05), resp.ModelPricing[2].ImageOutputPrice)
placeholder

// ---------------------------------------------------------------------------
// 2. pricingRequestToService
// ---------------------------------------------------------------------------

func TestPricingRequestToService_Defaults(t *testing.T) {
	tests := []struct {
		name      string
		req       channelModelPricingRequest
		wantField string // which default field to check
		wantValue string
placeholder{
		{
			name: "empty billing mode defaults to token",
			req: channelModelPricingRequest{
				Models:      []string{"m1"placeholder,
				BillingMode: "",
		placeholder,
			wantField: "BillingMode",
			wantValue: string(service.BillingModeToken),
	placeholder,
		{
			name: "empty platform stays empty",
			req: channelModelPricingRequest{
				Models:   []string{"m1"placeholder,
				Platform: "",
		placeholder,
			wantField: "Platform",
			wantValue: "",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricingRequestToService([]channelModelPricingRequest{tt.reqplaceholder, true)
			require.Len(t, result, 1)
			switch tt.wantField {
			case "BillingMode":
				require.Equal(t, service.BillingMode(tt.wantValue), result[0].BillingMode)
			case "Platform":
				require.Equal(t, tt.wantValue, result[0].Platform)
		placeholder
	placeholder)
placeholder
placeholder

func TestPricingRequestToService_WithAllFields(t *testing.T) {
	reqs := []channelModelPricingRequest{
		{
			Platform:         "openai",
			Models:           []string{"gpt-4", "gpt-4o"placeholder,
			BillingMode:      "per_request",
			InputPrice:       float64Ptr(0.01),
			OutputPrice:      float64Ptr(0.03),
			CacheWritePrice:  float64Ptr(0.005),
			CacheReadPrice:   float64Ptr(0.002),
			ImageOutputPrice: float64Ptr(0.04),
			PerRequestPrice:  float64Ptr(0.5),
	placeholder,
placeholder

	result := pricingRequestToService(reqs, true)
	require.Len(t, result, 1)
	r := result[0]
	require.Equal(t, "openai", r.Platform)
	require.Equal(t, []string{"gpt-4", "gpt-4o"placeholder, r.Models)
	require.Equal(t, service.BillingModePerRequest, r.BillingMode)
	require.Equal(t, float64Ptr(0.01), r.InputPrice)
	require.Equal(t, float64Ptr(0.03), r.OutputPrice)
	require.Equal(t, float64Ptr(0.005), r.CacheWritePrice)
	require.Equal(t, float64Ptr(0.002), r.CacheReadPrice)
	require.Equal(t, float64Ptr(0.04), r.ImageOutputPrice)
	require.Equal(t, float64Ptr(0.5), r.PerRequestPrice)
placeholder

func TestPricingRequestToService_WithIntervals(t *testing.T) {
	reqs := []channelModelPricingRequest{
		{
			Models:      []string{"m1"placeholder,
			BillingMode: "per_request",
			Intervals: []pricingIntervalRequest{
				{
					MinTokens:       0,
					MaxTokens:       intPtr(2000),
					TierLabel:       "small",
					InputPrice:      float64Ptr(0.01),
					OutputPrice:     float64Ptr(0.02),
					CacheWritePrice: float64Ptr(0.003),
					CacheReadPrice:  float64Ptr(0.001),
					PerRequestPrice: float64Ptr(0.1),
					SortOrder:       1,
			placeholder,
				{
					MinTokens: 2000,
					MaxTokens: nil,
					TierLabel: "large",
					SortOrder: 2,
			placeholder,
		placeholder,
	placeholder,
placeholder

	result := pricingRequestToService(reqs, true)
	require.Len(t, result, 1)
	require.Len(t, result[0].Intervals, 2)

	iv0 := result[0].Intervals[0]
	require.Equal(t, 0, iv0.MinTokens)
	require.Equal(t, intPtr(2000), iv0.MaxTokens)
	require.Equal(t, "small", iv0.TierLabel)
	require.Equal(t, float64Ptr(0.01), iv0.InputPrice)
	require.Equal(t, float64Ptr(0.02), iv0.OutputPrice)
	require.Equal(t, float64Ptr(0.003), iv0.CacheWritePrice)
	require.Equal(t, float64Ptr(0.001), iv0.CacheReadPrice)
	require.Equal(t, float64Ptr(0.1), iv0.PerRequestPrice)
	require.Equal(t, 1, iv0.SortOrder)

	iv1 := result[0].Intervals[1]
	require.Equal(t, 2000, iv1.MinTokens)
	require.Nil(t, iv1.MaxTokens)
	require.Equal(t, "large", iv1.TierLabel)
	require.Equal(t, 2, iv1.SortOrder)
placeholder

func TestPricingRequestToService_EmptySlice(t *testing.T) {
	result := pricingRequestToService([]channelModelPricingRequest{placeholder, true)
	require.NotNil(t, result)
	require.Empty(t, result)
placeholder

func TestPricingRequestToService_NilPriceFields(t *testing.T) {
	reqs := []channelModelPricingRequest{
		{
			Models:      []string{"m1"placeholder,
			BillingMode: "token",
			// all price fields are nil by default
	placeholder,
placeholder

	result := pricingRequestToService(reqs, true)
	require.Len(t, result, 1)
	r := result[0]
	require.Nil(t, r.InputPrice)
	require.Nil(t, r.OutputPrice)
	require.Nil(t, r.CacheWritePrice)
	require.Nil(t, r.CacheReadPrice)
	require.Nil(t, r.ImageOutputPrice)
	require.Nil(t, r.PerRequestPrice)
placeholder

func TestPricingRequestToService_TimePricing(t *testing.T) {
	req := channelModelPricingRequest{
		Models:      []string{"gpt-5"placeholder,
		BillingMode: "token",
		TimePricing: &channelTimePricingRequest{
			Timezone:     "Asia/Shanghai",
			WeekdaysOnly: true,
			Periods: []channelTimePricingPeriodRequest{{
				StartTime: "09:00", EndTime: "12:00", Multiplier: 2,
	placeholder
	placeholder,
placeholder

	got := pricingRequestToService([]channelModelPricingRequest{reqplaceholder, true)
	require.Equal(t, "Asia/Shanghai", got[0].TimePricing.Timezone)
	require.True(t, got[0].TimePricing.WeekdaysOnly)
	require.Equal(t, 2.0, got[0].TimePricing.Periods[0].Multiplier)
placeholder

func TestPricingRequestToService_TimePricingNil(t *testing.T) {
	got := pricingRequestToService([]channelModelPricingRequest{{Models: []string{"gpt-5"placeholderplaceholderplaceholder, true)
	require.Nil(t, got[0].TimePricing)
placeholder

// 账号成本统计规则不支持倍率：allowChannelMultipliers=false 时必须丢弃，
// 避免渠道倍率意外污染账号成本口径。
func TestPricingRequestToService_MultipliersGatedByFlag(t *testing.T) {
	req := channelModelPricingRequest{
		Models:         []string{"gpt-5"placeholder,
		BillingMode:    "token",
		FastMultiplier: float64Ptr(2.5),
		FlexMultiplier: float64Ptr(0.5),
		Intervals: []pricingIntervalRequest{{
			MinTokens:            272000,
			InputMultiplier:      float64Ptr(2),
			OutputMultiplier:     float64Ptr(1.5),
			CacheWriteMultiplier: float64Ptr(2),
			CacheReadMultiplier:  float64Ptr(2),
placeholder
placeholder

	allowed := pricingRequestToService([]channelModelPricingRequest{reqplaceholder, true)
	require.Equal(t, float64Ptr(2.5), allowed[0].FastMultiplier)
	require.Equal(t, float64Ptr(0.5), allowed[0].FlexMultiplier)
	require.Equal(t, float64Ptr(2), allowed[0].Intervals[0].InputMultiplier)
	require.Equal(t, float64Ptr(1.5), allowed[0].Intervals[0].OutputMultiplier)
	require.Equal(t, float64Ptr(2), allowed[0].Intervals[0].CacheWriteMultiplier)
	require.Equal(t, float64Ptr(2), allowed[0].Intervals[0].CacheReadMultiplier)

	dropped := pricingRequestToService([]channelModelPricingRequest{reqplaceholder, false)
	require.Nil(t, dropped[0].FastMultiplier)
	require.Nil(t, dropped[0].FlexMultiplier)
	require.Nil(t, dropped[0].Intervals[0].InputMultiplier)
	require.Nil(t, dropped[0].Intervals[0].OutputMultiplier)
	require.Nil(t, dropped[0].Intervals[0].CacheWriteMultiplier)
	require.Nil(t, dropped[0].Intervals[0].CacheReadMultiplier)
	// 非倍率字段不受开关影响
	require.Equal(t, 272000, dropped[0].Intervals[0].MinTokens)
placeholder

func TestPricingToResponse_TimePricing(t *testing.T) {
	got := pricingToResponse(&service.ChannelModelPricing{
		BillingMode: service.BillingModeToken,
		TimePricing: &service.ChannelTimePricing{
			Timezone:     "Asia/Shanghai",
			WeekdaysOnly: true,
			Periods: []service.ChannelTimePricingPeriod{{
				StartTime: "14:00", EndTime: "18:00", Multiplier: 1.25,
	placeholder
	placeholder,
placeholder)

	require.NotNil(t, got.TimePricing)
	require.Equal(t, "Asia/Shanghai", got.TimePricing.Timezone)
	require.True(t, got.TimePricing.WeekdaysOnly)
	require.Equal(t, 1.25, got.TimePricing.Periods[0].Multiplier)
placeholder

func TestPricingToResponse_TimePricingNil(t *testing.T) {
	got := pricingToResponse(&service.ChannelModelPricing{placeholder)
	require.Nil(t, got.TimePricing)
placeholder

// ---------------------------------------------------------------------------
// 3. SyncPricingModels handler
// ---------------------------------------------------------------------------

func setupSyncPricingModelsRouter(pricingSvc *service.PricingService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &ChannelHandler{pricingService: pricingSvcplaceholder
	router.GET("/channels/pricing/sync-models", h.SyncPricingModels)
	return router
placeholder

func TestSyncPricingModels_MissingPlatform(t *testing.T) {
	svc := service.NewPricingService(nil, nil)
	router := setupSyncPricingModelsRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/channels/pricing/sync-models", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
placeholder

func TestSyncPricingModels_UnsupportedPlatform(t *testing.T) {
	svc := service.NewPricingService(nil, nil)
	router := setupSyncPricingModelsRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/channels/pricing/sync-models?platform=unknown", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
placeholder

func TestSyncPricingModels_ValidPlatform_EmptyService(t *testing.T) {
	svc := service.NewPricingService(nil, nil)
	router := setupSyncPricingModelsRouter(svc)

	for _, platform := range []string{"anthropic", "openai", "gemini", "antigravity", "grok", "kimi", "zhipu", "deepseek"placeholder {
		req := httptest.NewRequest(http.MethodGet, "/channels/pricing/sync-models?platform="+platform, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, "platform=%s", platform)

		var body struct {
			Data struct {
				Models []string `json:"models"`
		placeholder `json:"data"`
	placeholder
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		require.NotNil(t, body.Data.Models, "models must not be null for platform=%s", platform)
placeholder
placeholder
