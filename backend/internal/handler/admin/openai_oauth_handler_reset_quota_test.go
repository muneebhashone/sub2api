//go:build unit

package admin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

type openAIQuotaWorkflowStub struct {
	resetResult *service.OpenAIQuotaResetResult
	resetErr    error
	queryResult *service.OpenAIQuotaUsage
	queryErr    error
	cacheErr    error

	resetCalls int
	queryCalls int
	cacheCalls int
placeholder

func (s *openAIQuotaWorkflowStub) ResetCredit(context.Context, int64) (*service.OpenAIQuotaResetResult, error) {
	s.resetCalls++
	return s.resetResult, s.resetErr
placeholder

func (s *openAIQuotaWorkflowStub) QueryUsage(context.Context, int64) (*service.OpenAIQuotaUsage, error) {
	s.queryCalls++
	return s.queryResult, s.queryErr
placeholder

func (s *openAIQuotaWorkflowStub) CacheResetCreditsSnapshot(context.Context, int64, *service.OpenAIRateLimitResetCredits) error {
	s.cacheCalls++
	return s.cacheErr
placeholder

type openAIAccountStateRecovererStub struct {
	err         error
	calls       int
	accountID   int64
	lastOptions service.AccountRecoveryOptions
placeholder

func (s *openAIAccountStateRecovererStub) RecoverAccountState(_ context.Context, accountID int64, options service.AccountRecoveryOptions) (*service.SuccessfulTestRecoveryResult, error) {
	s.calls++
	s.accountID = accountID
	s.lastOptions = options
	return &service.SuccessfulTestRecoveryResult{placeholder, s.err
placeholder

type openAIResetAdminServiceStub struct {
	service.AdminService
	account *service.Account
	err     error
	calls   int
placeholder

func (s *openAIResetAdminServiceStub) GetAccount(context.Context, int64) (*service.Account, error) {
	s.calls++
	return s.account, s.err
placeholder

type openAIQuotaResetEnvelope struct {
	Code int                      `json:"code"`
	Data openAIQuotaResetResponse `json:"data"`
placeholder

func performOpenAIQuotaResetRequest(t *testing.T, handler *OpenAIOAuthHandler) (int, openAIQuotaResetEnvelope) {
placeholder
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/api/v1/admin/openai/accounts/:id/reset-quota", handler.ResetQuota)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/openai/accounts/42/reset-quota", nil)
	router.ServeHTTP(recorder, request)

	var envelope openAIQuotaResetEnvelope
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &envelope))
	return recorder.Code, envelope
placeholder

func successfulOpenAIQuotaWorkflowStub() *openAIQuotaWorkflowStub {
	return &openAIQuotaWorkflowStub{
		resetResult: &service.OpenAIQuotaResetResult{
			Code:         "success",
			WindowsReset: 1,
	placeholder,
		queryResult: &service.OpenAIQuotaUsage{
			FetchedAt: 123,
			RateLimitResetCredits: &service.OpenAIRateLimitResetCredits{
				AvailableCount: 0,
				Credits:        []service.OpenAIRateLimitResetCreditDetail{placeholder,
		placeholder,
	placeholder,
placeholder
placeholder

func TestOpenAIResetQuota_ResetFailureStopsWorkflow(t *testing.T) {
	quota := &openAIQuotaWorkflowStub{resetErr: errors.New("upstream reset failed")placeholder
	recoverer := &openAIAccountStateRecovererStub{placeholder
	handler := &OpenAIOAuthHandler{
		adminService:     &openAIResetAdminServiceStub{placeholder,
		quotaService:     quota,
		rateLimitService: recoverer,
placeholder

	status, _ := performOpenAIQuotaResetRequest(t, handler)

	require.Equal(t, http.StatusInternalServerError, status)
	require.Equal(t, 1, quota.resetCalls)
	require.Zero(t, quota.queryCalls)
	require.Zero(t, quota.cacheCalls)
	require.Zero(t, recoverer.calls)
placeholder

func TestOpenAIResetQuota_QueryFailureReturnsPartialSuccessAndStops(t *testing.T) {
	quota := successfulOpenAIQuotaWorkflowStub()
	quota.queryResult = nil
	quota.queryErr = errors.New("upstream query failed")
	recoverer := &openAIAccountStateRecovererStub{placeholder
	handler := &OpenAIOAuthHandler{
		adminService:     &openAIResetAdminServiceStub{placeholder,
		quotaService:     quota,
		rateLimitService: recoverer,
placeholder

	status, envelope := performOpenAIQuotaResetRequest(t, handler)

	require.Equal(t, http.StatusOK, status)
	require.Equal(t, openAIQuotaResetWarningCacheRefreshFailed, envelope.Data.WarningCode)
	require.False(t, envelope.Data.CacheRefreshed)
	require.False(t, envelope.Data.AccountStateRecovered)
	require.Nil(t, envelope.Data.Quota)
	require.Equal(t, 1, quota.resetCalls)
	require.Equal(t, 1, quota.queryCalls)
	require.Zero(t, quota.cacheCalls)
	require.Zero(t, recoverer.calls)
placeholder

func TestOpenAIResetQuota_CacheFailureReturnsPartialSuccessAndStops(t *testing.T) {
	quota := successfulOpenAIQuotaWorkflowStub()
	quota.cacheErr = errors.New("cache write failed")
	recoverer := &openAIAccountStateRecovererStub{placeholder
	handler := &OpenAIOAuthHandler{
		adminService:     &openAIResetAdminServiceStub{placeholder,
		quotaService:     quota,
		rateLimitService: recoverer,
placeholder

	status, envelope := performOpenAIQuotaResetRequest(t, handler)

	require.Equal(t, http.StatusOK, status)
	require.Equal(t, openAIQuotaResetWarningCacheRefreshFailed, envelope.Data.WarningCode)
	require.False(t, envelope.Data.CacheRefreshed)
	require.Nil(t, envelope.Data.Quota)
	require.Equal(t, 1, quota.resetCalls)
	require.Equal(t, 1, quota.queryCalls)
	require.Equal(t, 1, quota.cacheCalls)
	require.Zero(t, recoverer.calls)
placeholder

func TestOpenAIResetQuota_RecoveryFailureKeepsRefreshedQuota(t *testing.T) {
	quota := successfulOpenAIQuotaWorkflowStub()
	recoverer := &openAIAccountStateRecovererStub{err: errors.New("recovery failed")placeholder
	adminService := &openAIResetAdminServiceStub{placeholder
	handler := &OpenAIOAuthHandler{
		adminService:     adminService,
		quotaService:     quota,
		rateLimitService: recoverer,
placeholder

	status, envelope := performOpenAIQuotaResetRequest(t, handler)

	require.Equal(t, http.StatusOK, status)
	require.Equal(t, openAIQuotaResetWarningAccountRecoveryFailed, envelope.Data.WarningCode)
	require.True(t, envelope.Data.CacheRefreshed)
	require.False(t, envelope.Data.AccountStateRecovered)
	require.NotNil(t, envelope.Data.Quota)
	require.Equal(t, 1, quota.resetCalls)
	require.Equal(t, 1, quota.queryCalls)
	require.Equal(t, 1, quota.cacheCalls)
	require.Equal(t, 1, recoverer.calls)
	require.Zero(t, adminService.calls)
placeholder

func TestOpenAIResetQuota_SuccessReturnsQuotaAndRecoveredAccount(t *testing.T) {
	quota := successfulOpenAIQuotaWorkflowStub()
	recoverer := &openAIAccountStateRecovererStub{placeholder
	adminService := &openAIResetAdminServiceStub{account: &service.Account{
		ID:          42,
		Name:        "recovered",
		Platform:    service.PlatformOpenAI,
		Type:        service.AccountTypeOAuth,
		Status:      service.StatusActive,
		Schedulable: false,
placeholderplaceholder
	handler := &OpenAIOAuthHandler{
		adminService:     adminService,
		quotaService:     quota,
		rateLimitService: recoverer,
placeholder

	status, envelope := performOpenAIQuotaResetRequest(t, handler)

	require.Equal(t, http.StatusOK, status)
	require.Empty(t, envelope.Data.WarningCode)
	require.True(t, envelope.Data.CacheRefreshed)
	require.True(t, envelope.Data.AccountStateRecovered)
	require.NotNil(t, envelope.Data.Quota)
	require.NotNil(t, envelope.Data.Account)
	require.Equal(t, int64(42), envelope.Data.Account.ID)
	require.False(t, envelope.Data.Account.Schedulable)
	require.Equal(t, int64(42), recoverer.accountID)
	require.True(t, recoverer.lastOptions.InvalidateToken)
	require.Equal(t, 1, quota.resetCalls)
	require.Equal(t, 1, quota.queryCalls)
	require.Equal(t, 1, quota.cacheCalls)
	require.Equal(t, 1, recoverer.calls)
	require.Equal(t, 1, adminService.calls)
placeholder

func TestOpenAIResetQuota_AccountRefreshFailureReportsRecoveredState(t *testing.T) {
	quota := successfulOpenAIQuotaWorkflowStub()
	recoverer := &openAIAccountStateRecovererStub{placeholder
	adminService := &openAIResetAdminServiceStub{err: errors.New("account refresh failed")placeholder
	handler := &OpenAIOAuthHandler{
		adminService:     adminService,
		quotaService:     quota,
		rateLimitService: recoverer,
placeholder

	status, envelope := performOpenAIQuotaResetRequest(t, handler)

	require.Equal(t, http.StatusOK, status)
	require.Equal(t, openAIQuotaResetWarningAccountRefreshFailed, envelope.Data.WarningCode)
	require.True(t, envelope.Data.CacheRefreshed)
	require.True(t, envelope.Data.AccountStateRecovered)
	require.NotNil(t, envelope.Data.Quota)
	require.Nil(t, envelope.Data.Account)
	require.Equal(t, 1, quota.resetCalls)
	require.Equal(t, 1, quota.queryCalls)
	require.Equal(t, 1, quota.cacheCalls)
	require.Equal(t, 1, recoverer.calls)
	require.Equal(t, 1, adminService.calls)
placeholder
