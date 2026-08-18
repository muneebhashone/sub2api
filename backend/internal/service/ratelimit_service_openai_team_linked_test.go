//go:build unit

package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

const teamLinkedDeactivatedBody = `{"detail":{"code":"deactivated_workspace","message":"This workspace has been deactivated."placeholderplaceholder`

type teamLinkedAccountRepoStub struct {
	mockAccountRepoForGemini
	teamAccounts []Account
	listErr      error
	listCalls    int
	setErrorIDs  []int64
	setErrorMsgs map[int64]string
	failSetError map[int64]error
placeholder

// ListByPlatform 镜像真实仓库语义：仅返回该平台的 active 账户。
func (r *teamLinkedAccountRepoStub) ListByPlatform(ctx context.Context, platform string) ([]Account, error) {
	r.listCalls++
	if r.listErr != nil {
		return nil, r.listErr
placeholder
	out := make([]Account, 0, len(r.teamAccounts))
	for _, acc := range r.teamAccounts {
		if acc.Platform == platform && acc.Status == StatusActive {
			out = append(out, acc)
	placeholder
placeholder
	return out, nil
placeholder

func (r *teamLinkedAccountRepoStub) SetError(ctx context.Context, id int64, errorMsg string) error {
	if err, ok := r.failSetError[id]; ok {
		return err
placeholder
	r.setErrorIDs = append(r.setErrorIDs, id)
	if r.setErrorMsgs == nil {
		r.setErrorMsgs = make(map[int64]string)
placeholder
	r.setErrorMsgs[id] = errorMsg
	return nil
placeholder

func newTeamLinkedAccount(id int64, teamID string) Account {
	return Account{
		ID:          id,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
placeholder"chatgpt_account_id": teamIDplaceholder,
placeholder
placeholder

// newTeamLinkedFixture: #1 触发者(team-A) #2 同队 #3 异队 #4 apikey #5 影子 #6 同队 #7 同队但已 error
func newTeamLinkedFixture() []Account {
	parentID := int64(1)
	shadow := Account{
		ID:              5,
		Platform:        PlatformOpenAI,
		Type:            AccountTypeOAuth,
		Status:          StatusActive,
		ParentAccountID: &parentID,
placeholder
	apikey := newTeamLinkedAccount(4, "team-A")
	apikey.Type = AccountTypeAPIKey
	erroredSibling := newTeamLinkedAccount(7, "team-A")
	erroredSibling.Status = StatusError
	return []Account{
		newTeamLinkedAccount(1, "team-A"),
		newTeamLinkedAccount(2, "team-A"),
		newTeamLinkedAccount(3, "team-B"),
		apikey,
		shadow,
		newTeamLinkedAccount(6, "team-A"),
		erroredSibling,
placeholder
placeholder

func newTeamLinkedTestService(repo *teamLinkedAccountRepoStub) (*RateLimitService, *runtimeBlockRecorder) {
	rl := NewRateLimitService(repo, nil, &config.Config{placeholder, nil, nil)
	blocker := &runtimeBlockRecorder{placeholder
	rl.SetAccountRuntimeBlocker(blocker)
	return rl, blocker
placeholder

func TestTeamLinkedError_FanoutMarksSameTeamAccounts(t *testing.T) {
	repo := &teamLinkedAccountRepoStub{teamAccounts: newTeamLinkedFixture()placeholder
	rl, blocker := newTeamLinkedTestService(repo)
	trigger := newTeamLinkedAccount(1, "team-A")

	shouldDisable := rl.HandleUpstreamError(context.Background(), &trigger, http.StatusPaymentRequired, http.Header{placeholder, []byte(teamLinkedDeactivatedBody))

	require.True(t, shouldDisable)
	// fan-out 先标记同队兄弟（#2、#6），触发账户 #1 随后由常规 case 402 标记
	require.Equal(t, []int64{2, 6, 1placeholder, repo.setErrorIDs)
	require.Contains(t, repo.setErrorMsgs[2], "team-linked error triggered by account #1")
	require.Contains(t, repo.setErrorMsgs[6], "team-linked error triggered by account #1")
	require.Contains(t, repo.setErrorMsgs[1], "Workspace deactivated (402)")
	require.NotContains(t, repo.setErrorMsgs[1], "team-linked")
	// 熔断顺序：兄弟账户先于落库全部进程内熔断，触发账户走 auth_error
	require.Equal(t, []string{openAITeamLinkedErrorBlockReason, openAITeamLinkedErrorBlockReason, "auth_error"placeholder, blocker.reasons)
	require.Equal(t, int64(2), blocker.accounts[0].ID)
	require.Equal(t, int64(6), blocker.accounts[1].ID)
	require.Equal(t, int64(1), blocker.accounts[2].ID)
placeholder

func TestTeamLinkedError_GenericPaymentErrorDoesNotFanout(t *testing.T) {
	repo := &teamLinkedAccountRepoStub{teamAccounts: newTeamLinkedFixture()placeholder
	rl, _ := newTeamLinkedTestService(repo)
	trigger := newTeamLinkedAccount(1, "team-A")

	rl.HandleUpstreamError(context.Background(), &trigger, http.StatusPaymentRequired, http.Header{placeholder, []byte(`{"error":{"message":"insufficient balance"placeholderplaceholder`))

	require.Equal(t, []int64{1placeholder, repo.setErrorIDs)
	require.Contains(t, repo.setErrorMsgs[1], "Payment required (402)")
	require.Zero(t, repo.listCalls)
placeholder

func TestTeamLinkedError_DedupWithinTTL(t *testing.T) {
	repo := &teamLinkedAccountRepoStub{teamAccounts: newTeamLinkedFixture()placeholder
	rl, _ := newTeamLinkedTestService(repo)
	first := newTeamLinkedAccount(1, "team-A")
	second := newTeamLinkedAccount(2, "team-A")

	rl.HandleUpstreamError(context.Background(), &first, http.StatusPaymentRequired, http.Header{placeholder, []byte(teamLinkedDeactivatedBody))
	rl.HandleUpstreamError(context.Background(), &second, http.StatusPaymentRequired, http.Header{placeholder, []byte(teamLinkedDeactivatedBody))

	// 第二次触发被去重：只有 #2 自身经 case 402 标记，未再次 fan-out
	require.Equal(t, []int64{2, 6, 1, 2placeholder, repo.setErrorIDs)
	require.Equal(t, 1, repo.listCalls)
placeholder

func TestTeamLinkedError_APIKeyTriggerDoesNotFanout(t *testing.T) {
	repo := &teamLinkedAccountRepoStub{teamAccounts: newTeamLinkedFixture()placeholder
	rl, _ := newTeamLinkedTestService(repo)
	trigger := newTeamLinkedAccount(4, "team-A")
	trigger.Type = AccountTypeAPIKey

	rl.HandleUpstreamError(context.Background(), &trigger, http.StatusPaymentRequired, http.Header{placeholder, []byte(teamLinkedDeactivatedBody))

	require.Equal(t, []int64{4placeholder, repo.setErrorIDs)
	require.Zero(t, repo.listCalls)
placeholder

func TestTeamLinkedError_DirectCallSkipsTriggerAccount(t *testing.T) {
	// 直调对应 fastpath 调用点：账户级临时不可调度规则短路时联动仍然生效
	repo := &teamLinkedAccountRepoStub{teamAccounts: newTeamLinkedFixture()placeholder
	rl, blocker := newTeamLinkedTestService(repo)
	trigger := newTeamLinkedAccount(1, "team-A")

	rl.maybeHandleOpenAITeamLinkedError(context.Background(), &trigger, http.StatusPaymentRequired, []byte(teamLinkedDeactivatedBody))

	require.Equal(t, []int64{2, 6placeholder, repo.setErrorIDs)
	require.Equal(t, []string{openAITeamLinkedErrorBlockReason, openAITeamLinkedErrorBlockReasonplaceholder, blocker.reasons)
placeholder

func TestTeamLinkedError_MissingTeamIDDoesNothing(t *testing.T) {
	repo := &teamLinkedAccountRepoStub{teamAccounts: newTeamLinkedFixture()placeholder
	rl, blocker := newTeamLinkedTestService(repo)
	trigger := Account{ID: 9, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusActiveplaceholder

	rl.maybeHandleOpenAITeamLinkedError(context.Background(), &trigger, http.StatusPaymentRequired, []byte(teamLinkedDeactivatedBody))

	require.Empty(t, repo.setErrorIDs)
	require.Empty(t, blocker.reasons)
	require.Zero(t, repo.listCalls)
placeholder

func TestTeamLinkedError_SetErrorFailureDoesNotAbortRemaining(t *testing.T) {
	repo := &teamLinkedAccountRepoStub{
		teamAccounts: newTeamLinkedFixture(),
		failSetError: map[int64]error{2: errors.New("db down")placeholder,
placeholder
	rl, blocker := newTeamLinkedTestService(repo)
	trigger := newTeamLinkedAccount(1, "team-A")

	rl.maybeHandleOpenAITeamLinkedError(context.Background(), &trigger, http.StatusPaymentRequired, []byte(teamLinkedDeactivatedBody))

	require.Equal(t, []int64{6placeholder, repo.setErrorIDs)
	// 进程内熔断先于落库执行，两个账户都已被熔断
	require.Len(t, blocker.reasons, 2)
placeholder
