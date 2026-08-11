//go:build unit

package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type modelNotFoundRateLimitCall struct {
	accountID int64
	scope     string
	resetAt   time.Time
	reason    string
placeholder

type modelNotFoundAccountRepoStub struct {
	mockAccountRepoForGemini
	tempCalls           int
	modelRateLimitCalls []modelNotFoundRateLimitCall
	modelRateLimitErr   error
placeholder

func (r *modelNotFoundAccountRepoStub) SetTempUnschedulable(ctx context.Context, id int64, until time.Time, reason string) error {
	r.tempCalls++
	return nil
placeholder

func (r *modelNotFoundAccountRepoStub) SetModelRateLimit(ctx context.Context, id int64, scope string, resetAt time.Time, reason ...string) error {
	call := modelNotFoundRateLimitCall{
		accountID: id,
		scope:     scope,
		resetAt:   resetAt,
placeholder
	if len(reason) > 0 {
		call.reason = reason[0]
placeholder
	r.modelRateLimitCalls = append(r.modelRateLimitCalls, call)
	return r.modelRateLimitErr
placeholder

func TestRateLimitService_HandleUpstreamError_ModelNotFoundUsesModelRateLimit(t *testing.T) {
	repo := &modelNotFoundAccountRepoStub{placeholder
	svc := &RateLimitService{accountRepo: repoplaceholder
	account := openAIModelNotFoundTempAccount()

	handled := svc.HandleUpstreamError(
		context.Background(),
		account,
		http.StatusNotFound,
		http.Header{placeholder,
		[]byte(`{"error":{"code":"model_not_found","message":"model not found"placeholderplaceholder`),
		"gpt-5.4",
	)

	require.True(t, handled)
	require.Zero(t, repo.tempCalls)
	require.Len(t, repo.modelRateLimitCalls, 1)
	call := repo.modelRateLimitCalls[0]
	require.Equal(t, account.ID, call.accountID)
	require.Equal(t, "gpt-5.4", call.scope)
	require.Equal(t, upstreamModelNotFoundReason, call.reason)
	require.WithinDuration(t, time.Now().Add(upstreamModelNotFoundCooldown), call.resetAt, 5*time.Second)
placeholder

func TestRateLimitService_HandleUpstreamError_ModelNotFoundWriteFailureDoesNotTempUnschedule(t *testing.T) {
	repo := &modelNotFoundAccountRepoStub{modelRateLimitErr: errors.New("write failed")placeholder
	svc := &RateLimitService{accountRepo: repoplaceholder
	account := openAIModelNotFoundTempAccount()

	handled := svc.HandleUpstreamError(
		context.Background(),
		account,
		http.StatusNotFound,
		http.Header{placeholder,
		[]byte(`{"error":{"code":"model_not_found","message":"model not found"placeholderplaceholder`),
		"gpt-5.4",
	)

	require.True(t, handled)
	require.Zero(t, repo.tempCalls)
	require.Len(t, repo.modelRateLimitCalls, 1)
placeholder

func TestRateLimitService_HandleUpstreamError_Bare404UsesModelScopedTempUnschedulableWhenModelKnown(t *testing.T) {
	repo := &modelNotFoundAccountRepoStub{placeholder
	svc := &RateLimitService{accountRepo: repoplaceholder
	account := openAIModelNotFoundTempAccount()

	handled := svc.HandleUpstreamError(
		context.Background(),
		account,
		http.StatusNotFound,
		http.Header{placeholder,
		[]byte(`{"error":{"message":"endpoint not found"placeholderplaceholder`),
		"gpt-5.4",
	)

	require.True(t, handled)
	require.Zero(t, repo.tempCalls)
	require.Len(t, repo.modelRateLimitCalls, 1)
	call := repo.modelRateLimitCalls[0]
	require.Equal(t, account.ID, call.accountID)
	require.Equal(t, "gpt-5.4", call.scope)
	require.WithinDuration(t, time.Now().Add(10*time.Minute), call.resetAt, 5*time.Second)

	var state TempUnschedState
	require.NoError(t, json.Unmarshal([]byte(call.reason), &state))
	require.Equal(t, http.StatusNotFound, state.StatusCode)
	require.Equal(t, "not found", state.MatchedKeyword)
placeholder

func TestRateLimitService_HandleUpstreamError_Bare404WithoutModelKeepsAccountTempUnschedulable(t *testing.T) {
	repo := &modelNotFoundAccountRepoStub{placeholder
	svc := &RateLimitService{accountRepo: repoplaceholder
	account := openAIModelNotFoundTempAccount()

	handled := svc.HandleUpstreamError(
		context.Background(),
		account,
		http.StatusNotFound,
		http.Header{placeholder,
		[]byte(`{"error":{"message":"endpoint not found"placeholderplaceholder`),
	)

	require.True(t, handled)
	require.Equal(t, 1, repo.tempCalls)
	require.Empty(t, repo.modelRateLimitCalls)
placeholder

func TestRateLimitService_HandleUpstreamError_ModelTempWriteFailureNeverWidensToAccount(t *testing.T) {
	repo := &modelNotFoundAccountRepoStub{modelRateLimitErr: errors.New("write failed")placeholder
	svc := &RateLimitService{accountRepo: repoplaceholder
	account := openAIModelNotFoundTempAccount()

	handled := svc.HandleUpstreamError(
		context.Background(),
		account,
		http.StatusNotFound,
		http.Header{placeholder,
		[]byte(`{"error":{"message":"endpoint not found"placeholderplaceholder`),
		"gpt-5.4",
	)

	require.True(t, handled)
	require.Zero(t, repo.tempCalls)
	require.Len(t, repo.modelRateLimitCalls, 1)
placeholder

func TestRateLimitService_HandleTempUnschedulable_PoolModeWithoutCustomPolicySkipsState(t *testing.T) {
	repo := &modelNotFoundAccountRepoStub{placeholder
	svc := &RateLimitService{accountRepo: repoplaceholder
	account := openAIModelNotFoundTempAccount()
	account.Credentials["pool_mode"] = true

	handled := svc.HandleTempUnschedulable(
		context.Background(),
		account,
		http.StatusNotFound,
		[]byte(`{"error":{"message":"endpoint not found"placeholderplaceholder`),
		"gpt-5.4",
	)

	require.False(t, handled)
	require.Zero(t, repo.tempCalls)
	require.Empty(t, repo.modelRateLimitCalls)
placeholder

func TestRateLimitService_HandleTempUnschedulable_PoolModeCustomPolicyUsesModelScope(t *testing.T) {
	repo := &modelNotFoundAccountRepoStub{placeholder
	svc := &RateLimitService{accountRepo: repoplaceholder
	account := openAIModelNotFoundTempAccount()
	account.Credentials["pool_mode"] = true
	account.Credentials["custom_error_codes_enabled"] = true
	account.Credentials["custom_error_codes"] = []any{float64(http.StatusNotFound)placeholder

	handled := svc.HandleTempUnschedulable(
		context.Background(),
		account,
		http.StatusNotFound,
		[]byte(`{"error":{"message":"endpoint not found"placeholderplaceholder`),
		"gpt-5.4",
	)

	require.True(t, handled)
	require.Zero(t, repo.tempCalls)
	require.Len(t, repo.modelRateLimitCalls, 1)
	require.Equal(t, "gpt-5.4", repo.modelRateLimitCalls[0].scope)
placeholder

func TestRateLimitService_TempUnschedulableContextPreservesModelForPoolDependency(t *testing.T) {
	repo := &modelNotFoundAccountRepoStub{placeholder
	svc := &RateLimitService{accountRepo: repoplaceholder
	account := openAIModelNotFoundTempAccount()
	account.Credentials["pool_mode"] = true
	ctx := withTempUnschedulableModel(context.Background(), []string{"gpt-5.4"placeholder)

	// #4496 calls tryTempUnschedulable from the pool-mode branch without an
	// explicit model argument. The request context must preserve the canonical
	// model so that combined behavior remains model-scoped after it lands.
	handled := svc.tryTempUnschedulable(
		ctx,
		account,
		http.StatusNotFound,
		[]byte(`{"error":{"message":"endpoint not found"placeholderplaceholder`),
	)

	require.True(t, handled)
	require.Zero(t, repo.tempCalls)
	require.Len(t, repo.modelRateLimitCalls, 1)
	require.Equal(t, "gpt-5.4", repo.modelRateLimitCalls[0].scope)
placeholder

func TestRateLimitService_HandleUpstreamError_CustomPolicyExclusionSkipsAllState(t *testing.T) {
	repo := &modelNotFoundAccountRepoStub{placeholder
	svc := &RateLimitService{accountRepo: repoplaceholder
	account := openAIModelNotFoundTempAccount()
	account.Credentials["custom_error_codes_enabled"] = true
	account.Credentials["custom_error_codes"] = []any{float64(http.StatusServiceUnavailable)placeholder

	handled := svc.HandleUpstreamError(
		context.Background(),
		account,
		http.StatusNotFound,
		http.Header{placeholder,
		[]byte(`{"error":{"message":"endpoint not found"placeholderplaceholder`),
		"gpt-5.4",
	)

	require.False(t, handled)
	require.Zero(t, repo.tempCalls)
	require.Empty(t, repo.modelRateLimitCalls)
placeholder

func TestRateLimitService_HandleTempUnschedulable_AuthenticationFailureStaysAccountScoped(t *testing.T) {
	repo := &modelNotFoundAccountRepoStub{placeholder
	svc := &RateLimitService{accountRepo: repoplaceholder
	account := openAIModelNotFoundTempAccount()
	account.TempUnschedulableReason = "legacy non-JSON reason"
	account.Credentials["temp_unschedulable_rules"] = []any{
		map[string]any{
			"error_code":       float64(http.StatusUnauthorized),
			"keywords":         []any{"unauthorized"placeholder,
			"duration_minutes": float64(10),
	placeholder,
placeholder

	handled := svc.HandleTempUnschedulable(
		context.Background(),
		account,
		http.StatusUnauthorized,
		[]byte(`{"error":{"message":"unauthorized"placeholderplaceholder`),
		"gpt-5.4",
	)

	require.True(t, handled)
	require.Equal(t, 1, repo.tempCalls)
	require.Empty(t, repo.modelRateLimitCalls)
placeholder

func TestRateLimitService_ModelTempUnschedulableIsolatesSchedulerByModel(t *testing.T) {
	repo := &modelNotFoundAccountRepoStub{placeholder
	svc := &RateLimitService{accountRepo: repoplaceholder
	account := openAIModelNotFoundTempAccount()
	account.Credentials["model_mapping"] = map[string]any{
		"public-a":   "upstream-a",
		"upstream-a": "upstream-b",
placeholder

	handled := svc.HandleUpstreamError(
		context.Background(),
		account,
		http.StatusNotFound,
		http.Header{placeholder,
		[]byte(`{"error":{"message":"endpoint not found"placeholderplaceholder`),
		"upstream-a",
	)

	require.True(t, handled)
	require.Len(t, repo.modelRateLimitCalls, 1)
	call := repo.modelRateLimitCalls[0]
	require.Equal(t, "upstream-a", call.scope, "canonical upstream model must not be mapped a second time")

	account.Extra = map[string]any{
		modelRateLimitsKey: map[string]any{
			call.scope: map[string]any{
				"rate_limit_reset_at": call.resetAt.UTC().Format(time.RFC3339),
		placeholder,
	placeholder,
placeholder

	require.False(t, account.IsSchedulableForModelWithContext(context.Background(), "public-a"))
	require.True(t, account.IsSchedulableForModelWithContext(context.Background(), "gpt-5.6-sol"))
placeholder

func openAIModelNotFoundTempAccount() *Account {
placeholder
		ID:          101,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Schedulable: true,
placeholder
			"temp_unschedulable_enabled": true,
			"temp_unschedulable_rules": []any{
				map[string]any{
					"error_code":       float64(http.StatusNotFound),
					"keywords":         []any{"not found"placeholder,
					"duration_minutes": float64(10),
			placeholder,
		placeholder,
	placeholder,
placeholder
placeholder

func TestRateLimitService_HandleUpstreamError_CodexPlanGatedModelUsesModelRateLimit(t *testing.T) {
	repo := &modelNotFoundAccountRepoStub{placeholder
	svc := &RateLimitService{accountRepo: repoplaceholder
	account := openAICodexPlanGatedOAuthAccount()

	handled := svc.HandleUpstreamError(
		context.Background(),
		account,
		http.StatusBadRequest,
		http.Header{placeholder,
		[]byte(`{"detail":"The 'gpt-5.6-sol' model is not supported when using Codex with a ChatGPT account."placeholder`),
		"gpt-5.6-sol",
	)

	require.True(t, handled)
	require.Zero(t, repo.tempCalls)
	require.Len(t, repo.modelRateLimitCalls, 1)
	call := repo.modelRateLimitCalls[0]
	require.Equal(t, account.ID, call.accountID)
	require.Equal(t, "gpt-5.6-sol", call.scope)
	require.Equal(t, upstreamCodexPlanGatedModelReason, call.reason)
	require.WithinDuration(t, time.Now().Add(upstreamCodexPlanGatedModelCooldown), call.resetAt, 5*time.Second)
placeholder

func TestRateLimitService_HandleUpstreamError_CodexPlanGatedModelRespectsModelMapping(t *testing.T) {
	repo := &modelNotFoundAccountRepoStub{placeholder
	svc := &RateLimitService{accountRepo: repoplaceholder
	account := openAICodexPlanGatedOAuthAccount()
	account.Credentials["model_mapping"] = map[string]any{"gpt-5.6-sol": "gpt-5.6-sol-upstream"placeholder

	handled := svc.HandleUpstreamError(
		context.Background(),
		account,
		http.StatusBadRequest,
		http.Header{placeholder,
		[]byte(`{"detail":"The 'gpt-5.6-sol-upstream' model is not supported when using Codex with a ChatGPT account."placeholder`),
		"gpt-5.6-sol",
	)

	require.True(t, handled)
	require.Len(t, repo.modelRateLimitCalls, 1)
	require.Equal(t, "gpt-5.6-sol-upstream", repo.modelRateLimitCalls[0].scope)
placeholder

func TestRateLimitService_HandleUpstreamError_CodexPlanGatedModelIgnoresAPIKeyAccount(t *testing.T) {
	repo := &modelNotFoundAccountRepoStub{placeholder
	svc := &RateLimitService{accountRepo: repoplaceholder
	account := openAICodexPlanGatedOAuthAccount()
	account.Type = AccountTypeAPIKey

	handled := svc.HandleUpstreamError(
		context.Background(),
		account,
		http.StatusBadRequest,
		http.Header{placeholder,
		[]byte(`{"detail":"The 'gpt-5.6-sol' model is not supported when using Codex with a ChatGPT account."placeholder`),
		"gpt-5.6-sol",
	)

	require.False(t, handled)
	require.Empty(t, repo.modelRateLimitCalls)
placeholder

func TestRateLimitService_HandleUpstreamError_CodexPlanGatedImageModelSkipsCooldown(t *testing.T) {
	for _, model := range []string{"gpt-image-1", "gpt-image-1.5", "gpt-image-2"placeholder {
		t.Run(model, func(t *testing.T) {
			repo := &modelNotFoundAccountRepoStub{placeholder
			svc := &RateLimitService{accountRepo: repoplaceholder
			account := openAICodexPlanGatedOAuthAccount()

			handled := svc.HandleUpstreamError(
				context.Background(),
				account,
				http.StatusBadRequest,
				http.Header{placeholder,
				[]byte(`{"detail":"The '`+model+`' model is not supported when using Codex with a ChatGPT account."placeholder`),
				model,
			)

			require.True(t, handled, "attempt should still fail over")
			require.Empty(t, repo.modelRateLimitCalls,
				"image models must not be cooled down: the account still serves them over /v1/images/*")
			require.Zero(t, repo.tempCalls)
	placeholder)
placeholder
placeholder

func TestRateLimitService_HandleUpstreamError_CodexPlanGatedTextModelStillCoolsDown(t *testing.T) {
	repo := &modelNotFoundAccountRepoStub{placeholder
	svc := &RateLimitService{accountRepo: repoplaceholder
	account := openAICodexPlanGatedOAuthAccount()

	handled := svc.HandleUpstreamError(
		context.Background(),
		account,
		http.StatusBadRequest,
		http.Header{placeholder,
		[]byte(`{"detail":"The 'gpt-5.6-sol' model is not supported when using Codex with a ChatGPT account."placeholder`),
		"gpt-5.6-sol",
	)

	require.True(t, handled)
	require.Len(t, repo.modelRateLimitCalls, 1, "non-image plan-gated models keep the existing cooldown")
	require.Equal(t, upstreamCodexPlanGatedModelReason, repo.modelRateLimitCalls[0].reason)
placeholder

func openAICodexPlanGatedOAuthAccount() *Account {
placeholder
		ID:          202,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Schedulable: true,
placeholderplaceholder,
placeholder
placeholder

// 请求本身就走 /v1/images/* 时必须保留冷却。
//
// OAuth 账号的 /v1/images/* 上游同样是 Codex Responses（openai_images_responses.go
// → handleOpenAIImagesErrorResponse → handleOpenAIAccountUpstreamError →
// HandleUpstreamModelNotFound），所以这条路径也会命中 plan-gated 分支。账号确实
// 不具备生图能力时，冷却是唯一的刹车：调度层靠 model_rate_limits 跳过该账号后
// 快速 503；一旦跳过冷却，每个请求都会完整走一遍号池，对上游形成无上界的 400 放大。
func TestRateLimitService_HandleUpstreamError_CodexPlanGatedImageModelKeepsCooldownOnImagesEndpoint(t *testing.T) {
	repo := &modelNotFoundAccountRepoStub{placeholder
	svc := &RateLimitService{accountRepo: repoplaceholder
	account := openAICodexPlanGatedOAuthAccount()

	handled := svc.HandleUpstreamError(
		WithOpenAIImagesEndpoint(context.Background()),
		account,
		http.StatusBadRequest,
		http.Header{placeholder,
		[]byte(`{"detail":"The 'gpt-image-2' model is not supported when using Codex with a ChatGPT account."placeholder`),
		"gpt-image-2",
	)

	require.True(t, handled)
	require.Len(t, repo.modelRateLimitCalls, 1,
		"/v1/images/* 上的 plan-gated 拒绝是真实的能力缺失，必须保留冷却刹车")
	require.Equal(t, "gpt-image-2", repo.modelRateLimitCalls[0].scope)
	require.Equal(t, upstreamCodexPlanGatedModelReason, repo.modelRateLimitCalls[0].reason)
placeholder

// 仅 WithOpenAIImageGenerationIntent（/v1/responses 因模型名自动置位）不算专用生图
// 端点，仍按"用错端点"处理。
func TestRateLimitService_HandleUpstreamError_CodexPlanGatedImageModelSkipsCooldownOnIntentOnly(t *testing.T) {
	repo := &modelNotFoundAccountRepoStub{placeholder
	svc := &RateLimitService{accountRepo: repoplaceholder
	account := openAICodexPlanGatedOAuthAccount()

	handled := svc.HandleUpstreamError(
		WithOpenAIImageGenerationIntent(context.Background()),
		account,
		http.StatusBadRequest,
		http.Header{placeholder,
		[]byte(`{"detail":"The 'gpt-image-2' model is not supported when using Codex with a ChatGPT account."placeholder`),
		"gpt-image-2",
	)

	require.True(t, handled)
	require.Empty(t, repo.modelRateLimitCalls)
placeholder

// 守卫口径必须与冷却键一致：冷却键走 account.GetMappedModel，账号可以把文本别名
// 映射到 gpt-image-*，只判请求模型会漏掉这种形态，原 bug 原样复现。
func TestRateLimitService_HandleUpstreamError_CodexPlanGatedImageModelSkipsCooldownViaModelMapping(t *testing.T) {
	repo := &modelNotFoundAccountRepoStub{placeholder
	svc := &RateLimitService{accountRepo: repoplaceholder
	account := openAICodexPlanGatedOAuthAccount()
	account.Credentials["model_mapping"] = map[string]any{"my-draw-alias": "gpt-image-2"placeholder

	handled := svc.HandleUpstreamError(
		context.Background(),
		account,
		http.StatusBadRequest,
		http.Header{placeholder,
		[]byte(`{"detail":"The 'gpt-image-2' model is not supported when using Codex with a ChatGPT account."placeholder`),
		"my-draw-alias",
	)

	require.True(t, handled)
	require.Empty(t, repo.modelRateLimitCalls,
		"映射后的上游模型是图片模型，冷却键会写到 gpt-image-2 上，守卫必须一并识别")
placeholder

// 404 model-not-found 分支不受守卫影响：即使是图片模型也照常冷却。
func TestRateLimitService_HandleUpstreamError_ModelNotFoundImageModelStillCoolsDown(t *testing.T) {
	repo := &modelNotFoundAccountRepoStub{placeholder
	svc := &RateLimitService{accountRepo: repoplaceholder
	account := openAICodexPlanGatedOAuthAccount()

	handled := svc.HandleUpstreamError(
		context.Background(),
		account,
		http.StatusNotFound,
		http.Header{placeholder,
		[]byte(`{"error":{"message":"The model 'gpt-image-2' does not exist","code":"model_not_found"placeholderplaceholder`),
		"gpt-image-2",
	)

	require.True(t, handled)
	require.Len(t, repo.modelRateLimitCalls, 1, "守卫只作用于 codex plan-gated 分支")
	require.Equal(t, upstreamModelNotFoundReason, repo.modelRateLimitCalls[0].reason)
placeholder
