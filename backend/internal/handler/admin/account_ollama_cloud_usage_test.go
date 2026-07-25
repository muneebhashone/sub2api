package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type ollamaCloudUsageHandlerTestRepo struct {
	service.AccountRepository
	account           *service.Account
	accounts          []*service.Account
	groupResolveCalls int
placeholder

func (r *ollamaCloudUsageHandlerTestRepo) GetByID(_ context.Context, id int64) (*service.Account, error) {
	if r.account != nil && r.account.ID == id {
		return r.account, nil
placeholder
	for _, account := range r.accounts {
		if account.ID == id {
			return account, nil
	placeholder
placeholder
	return nil, service.ErrAccountNotFound
placeholder

func (r *ollamaCloudUsageHandlerTestRepo) ListOllamaCloudUsageGroupAccounts(_ context.Context, _ []*service.Account) ([]service.Account, error) {
	r.groupResolveCalls++
	result := make([]service.Account, 0, len(r.accounts)+1)
	if r.account != nil {
		result = append(result, *r.account)
placeholder
	for _, account := range r.accounts {
		result = append(result, *account)
placeholder
	return result, nil
placeholder

func (r *ollamaCloudUsageHandlerTestRepo) SaveOllamaCloudUsageSession(context.Context, *service.Account, string, bool) error {
	return nil
placeholder
func (r *ollamaCloudUsageHandlerTestRepo) DeleteOllamaCloudUsageSession(context.Context, *service.Account) error {
	return nil
placeholder
func (r *ollamaCloudUsageHandlerTestRepo) SetOllamaCloudUsageAutoRefresh(context.Context, *service.Account, bool) error {
	return nil
placeholder
func (r *ollamaCloudUsageHandlerTestRepo) UpdateOllamaCloudUsageSnapshot(context.Context, *service.Account, *service.OllamaCloudUsageSnapshot) error {
	return nil
placeholder
func (r *ollamaCloudUsageHandlerTestRepo) DisableOllamaCloudUsageAutoRefresh(context.Context, *service.Account) error {
	return nil
placeholder
func (r *ollamaCloudUsageHandlerTestRepo) ListDueOllamaCloudUsageAccounts(context.Context, time.Time, int) ([]service.Account, error) {
	return nil, nil
placeholder

func newOllamaCloudUsageHandlerTestService(t *testing.T) *service.OllamaCloudUsageService {
placeholder
	svc := service.NewOllamaCloudUsageService(nil, nil, nil, nil, false)
	t.Cleanup(svc.Stop)
	return svc
placeholder

func newOllamaCloudUsageHandlerContext(method, target, body, id string) (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, target, bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request
	if id != "" {
		ctx.Params = gin.Params{{Key: "id", Value: idplaceholderplaceholder
placeholder
	return ctx, recorder
placeholder

func TestOllamaCloudUsageHandlersValidateRequestsAndDependencies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newOllamaCloudUsageHandlerTestService(t)

	t.Run("invalid account id", func(t *testing.T) {
		ctx, recorder := newOllamaCloudUsageHandlerContext(http.MethodGet, "/admin/accounts/not-an-id/ollama-cloud-usage", "", "not-an-id")
		(&AccountHandler{ollamaCloudUsage: svcplaceholder).GetOllamaCloudUsage(ctx)
		require.Equal(t, http.StatusBadRequest, recorder.Code)
placeholder)

	t.Run("empty session", func(t *testing.T) {
		ctx, recorder := newOllamaCloudUsageHandlerContext(http.MethodPut, "/admin/accounts/7/ollama-cloud-usage/session", `{"session":""placeholder`, "7")
		(&AccountHandler{ollamaCloudUsage: svcplaceholder).SaveOllamaCloudUsageSession(ctx)
		require.Equal(t, http.StatusBadRequest, recorder.Code)
placeholder)

	t.Run("missing enabled", func(t *testing.T) {
		ctx, recorder := newOllamaCloudUsageHandlerContext(http.MethodPut, "/admin/accounts/7/ollama-cloud-usage/auto-refresh", `{placeholder`, "7")
		(&AccountHandler{ollamaCloudUsage: svcplaceholder).SetOllamaCloudUsageAutoRefresh(ctx)
		require.Equal(t, http.StatusBadRequest, recorder.Code)
placeholder)

	t.Run("service unavailable", func(t *testing.T) {
		ctx, recorder := newOllamaCloudUsageHandlerContext(http.MethodGet, "/admin/accounts/7/ollama-cloud-usage", "", "7")
		(&AccountHandler{placeholder).GetOllamaCloudUsage(ctx)
		require.Equal(t, http.StatusServiceUnavailable, recorder.Code)
		require.Contains(t, recorder.Body.String(), "OLLAMA_CLOUD_USAGE_UNAVAILABLE")
placeholder)
placeholder

func TestOllamaCloudUsageEncryptionKeyStateConsistentAcrossAccountResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	for _, configured := range []bool{false, trueplaceholder {
		t.Run("configured="+strconv.FormatBool(configured), func(t *testing.T) {
			account := &service.Account{
				ID:          7,
				Name:        "ollama",
				Platform:    service.PlatformOpenAI,
				Type:        service.AccountTypeAPIKey,
		placeholder"base_url": "https://ollama.com", "api_key": "test-key"placeholder,
				Extra:       map[string]any{placeholder,
				Status:      service.StatusActive,
		placeholder
			adminService := newStubAdminService()
			adminService.accounts = []service.Account{*accountplaceholder
			adminService.getAccountResult = account
			usageService := service.NewOllamaCloudUsageService(
				&ollamaCloudUsageHandlerTestRepo{account: accountplaceholder, nil, nil, nil, configured,
			)
			t.Cleanup(usageService.Stop)

			handler := NewAccountHandler(adminService, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
			handler.SetOllamaCloudUsageService(usageService)
			router := gin.New()
			router.GET("/accounts", handler.List)
			router.GET("/accounts/:id", handler.GetByID)
			router.GET("/accounts/:id/ollama-cloud-usage", handler.GetOllamaCloudUsage)

			listRecorder := httptest.NewRecorder()
			router.ServeHTTP(listRecorder, httptest.NewRequest(http.MethodGet, "/accounts?page=1&page_size=20", nil))
			require.Equal(t, http.StatusOK, listRecorder.Code)
			var listPayload struct {
				Data struct {
					Items []struct {
						OllamaCloudUsage *service.OllamaCloudUsageState `json:"ollama_cloud_usage"`
				placeholder `json:"items"`
			placeholder `json:"data"`
		placeholder
			require.NoError(t, json.Unmarshal(listRecorder.Body.Bytes(), &listPayload))
			require.Len(t, listPayload.Data.Items, 1)
			require.NotNil(t, listPayload.Data.Items[0].OllamaCloudUsage)

			detailRecorder := httptest.NewRecorder()
			router.ServeHTTP(detailRecorder, httptest.NewRequest(http.MethodGet, "/accounts/7", nil))
			require.Equal(t, http.StatusOK, detailRecorder.Code)
			var detailPayload struct {
				Data struct {
					OllamaCloudUsage *service.OllamaCloudUsageState `json:"ollama_cloud_usage"`
			placeholder `json:"data"`
		placeholder
			require.NoError(t, json.Unmarshal(detailRecorder.Body.Bytes(), &detailPayload))
			require.NotNil(t, detailPayload.Data.OllamaCloudUsage)

			stateRecorder := httptest.NewRecorder()
			router.ServeHTTP(stateRecorder, httptest.NewRequest(http.MethodGet, "/accounts/7/ollama-cloud-usage", nil))
			require.Equal(t, http.StatusOK, stateRecorder.Code)
			var statePayload struct {
				Data service.OllamaCloudUsageState `json:"data"`
		placeholder
			require.NoError(t, json.Unmarshal(stateRecorder.Body.Bytes(), &statePayload))

			listConfigured := listPayload.Data.Items[0].OllamaCloudUsage.EncryptionKeyConfigured
			detailConfigured := detailPayload.Data.OllamaCloudUsage.EncryptionKeyConfigured
			require.Equal(t, configured, listConfigured)
			require.Equal(t, statePayload.Data.EncryptionKeyConfigured, listConfigured)
			require.Equal(t, statePayload.Data.EncryptionKeyConfigured, detailConfigured)
	placeholder)
placeholder
placeholder

func TestOllamaCloudUsageSharedStateMatchesListDetailAndSpecialEndpointWithoutListNPlusOne(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now().UTC()
	source := &service.Account{
		ID: 7, Name: "source", Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey,
placeholder"base_url": "https://ollama.com", "api_key": "shared-secret-key"placeholder,
		Extra: map[string]any{
			service.OllamaCloudUsageSessionExtraKey:     "ciphertext-secret",
			service.OllamaCloudUsageAutoRefreshExtraKey: true,
			service.OllamaCloudUsageSnapshotExtraKey: &service.OllamaCloudUsageSnapshot{
				Status: service.OllamaCloudUsageStatusOK, Data: &service.OllamaCloudUsageData{Plan: "pro"placeholder,
				LastAttemptAt: now, NextRefreshAt: now.Add(time.Hour),
		placeholder,
	placeholder,
		Status: service.StatusActive,
placeholder
	sibling := &service.Account{
		ID: 8, Name: "sibling", Platform: service.PlatformAnthropic, Type: service.AccountTypeAPIKey,
placeholder"base_url": "HTTPS://WWW.OLLAMA.COM:443/v1", "api_key": "shared-secret-key"placeholder,
		Extra:       map[string]any{placeholder, Status: service.StatusActive,
placeholder
	repo := &ollamaCloudUsageHandlerTestRepo{accounts: []*service.Account{source, siblingplaceholderplaceholder
	adminService := newStubAdminService()
	adminService.accounts = []service.Account{*source, *siblingplaceholder
	adminService.getAccountResult = sibling
	usageService := service.NewOllamaCloudUsageService(repo, nil, nil, nil, true)
	t.Cleanup(usageService.Stop)
	handler := NewAccountHandler(adminService, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	handler.SetOllamaCloudUsageService(usageService)
	router := gin.New()
	router.GET("/accounts", handler.List)
	router.GET("/accounts/:id", handler.GetByID)
	router.GET("/accounts/:id/ollama-cloud-usage", handler.GetOllamaCloudUsage)

	listRecorder := httptest.NewRecorder()
	router.ServeHTTP(listRecorder, httptest.NewRequest(http.MethodGet, "/accounts?page=1&page_size=20", nil))
	require.Equal(t, http.StatusOK, listRecorder.Code)
	require.Equal(t, 1, repo.groupResolveCalls, "the full list page must use one group-resolution batch")
	var listPayload struct {
		Data struct {
			Items []struct {
				ID               int64                          `json:"id"`
				OllamaCloudUsage *service.OllamaCloudUsageState `json:"ollama_cloud_usage"`
		placeholder `json:"items"`
	placeholder `json:"data"`
placeholder
	require.NoError(t, json.Unmarshal(listRecorder.Body.Bytes(), &listPayload))
	require.Len(t, listPayload.Data.Items, 2)
	for _, item := range listPayload.Data.Items {
		require.True(t, item.OllamaCloudUsage.Configured)
		require.Equal(t, "pro", item.OllamaCloudUsage.Snapshot.Data.Plan)
placeholder

	detailRecorder := httptest.NewRecorder()
	router.ServeHTTP(detailRecorder, httptest.NewRequest(http.MethodGet, "/accounts/8", nil))
	require.Equal(t, http.StatusOK, detailRecorder.Code)
	var detailPayload struct {
		Data struct {
			OllamaCloudUsage *service.OllamaCloudUsageState `json:"ollama_cloud_usage"`
	placeholder `json:"data"`
placeholder
	require.NoError(t, json.Unmarshal(detailRecorder.Body.Bytes(), &detailPayload))

	stateRecorder := httptest.NewRecorder()
	router.ServeHTTP(stateRecorder, httptest.NewRequest(http.MethodGet, "/accounts/8/ollama-cloud-usage", nil))
	require.Equal(t, http.StatusOK, stateRecorder.Code)
	var statePayload struct {
		Data service.OllamaCloudUsageState `json:"data"`
placeholder
	require.NoError(t, json.Unmarshal(stateRecorder.Body.Bytes(), &statePayload))
	require.Equal(t, statePayload.Data.Configured, detailPayload.Data.OllamaCloudUsage.Configured)
	require.Equal(t, statePayload.Data.Snapshot, detailPayload.Data.OllamaCloudUsage.Snapshot)
	for _, body := range []string{listRecorder.Body.String(), detailRecorder.Body.String(), stateRecorder.Body.String()placeholder {
		require.NotContains(t, body, "shared-secret-key")
		require.NotContains(t, body, "ciphertext-secret")
placeholder
placeholder

func TestGetOllamaCloudUsageSettingsHandlerSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, recorder := newOllamaCloudUsageHandlerContext(http.MethodGet, "/admin/accounts/ollama-cloud-usage/settings", "", "")
	handler := &AccountHandler{ollamaCloudUsage: newOllamaCloudUsageHandlerTestService(t)placeholder

	handler.GetOllamaCloudUsageSettings(ctx)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Contains(t, recorder.Body.String(), `"enabled":false`)
	require.Contains(t, recorder.Body.String(), `"interval_minutes":60`)
	require.Contains(t, recorder.Body.String(), `"debounce_minutes":1`)
placeholder
