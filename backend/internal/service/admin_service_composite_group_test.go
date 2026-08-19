//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/claude"
	"github.com/stretchr/testify/require"
)

type accountRepoStubForCompositeModelsList struct {
	accountRepoStub
	accounts []Account
placeholder

func (s *accountRepoStubForCompositeModelsList) ListSchedulableByGroupID(_ context.Context, _ int64) ([]Account, error) {
	return s.accounts, nil
placeholder

func TestAdminService_CreateCompositeGroupCopiesAccountsFromConcreteGroups(t *testing.T) {
	var copiedFrom []int64
	var boundGroupID int64
	var boundAccountIDs []int64
	groupRepo := &groupRepoStubForAdmin{
		createID: 99,
		getByIDByID: map[int64]*Group{
			10: {ID: 10, Platform: PlatformOpenAIplaceholder,
			20: {ID: 20, Platform: PlatformGeminiplaceholder,
	placeholder,
		getAccountIDsByGroupIDsFn: func(groupIDs []int64) ([]int64, error) {
			copiedFrom = append([]int64{placeholder, groupIDs...)
			return []int64{101, 202placeholder, nil
	placeholder,
		bindAccountsToGroupFn: func(groupID int64, accountIDs []int64) error {
			boundGroupID = groupID
			boundAccountIDs = append([]int64{placeholder, accountIDs...)
			return nil
	placeholder,
placeholder
	svc := &adminServiceImpl{groupRepo: groupRepoplaceholder

	group, err := svc.CreateGroup(context.Background(), &CreateGroupInput{
		Name:               "Composite",
		Platform:           PlatformComposite,
		RateMultiplier:     1,
		MaxReasoningEffort: "medium",
		ReasoningEffortMappings: []ReasoningEffortMapping{
			{From: "max", To: "xhigh"placeholder,
	placeholder,
		CopyAccountsFromGroupIDs: []int64{10, 20, 10placeholder,
placeholder)

placeholder
	require.Equal(t, PlatformComposite, groupRepo.created.Platform)
	require.Equal(t, "medium", groupRepo.created.MaxReasoningEffort)
	require.Equal(t, []ReasoningEffortMapping{{From: "max", To: "xhigh"placeholderplaceholder, groupRepo.created.ReasoningEffortMappings)
	require.Equal(t, int64(99), group.ID)
	require.Equal(t, int64(2), group.AccountCount)
	require.ElementsMatch(t, []int64{10, 20placeholder, copiedFrom)
	require.Equal(t, int64(99), boundGroupID)
	require.ElementsMatch(t, []int64{101, 202placeholder, boundAccountIDs)
placeholder

func TestAdminService_UpdateCompositeGroupCopiesAccountsFromConcreteGroups(t *testing.T) {
	var clearedGroupID int64
	var copiedFrom []int64
	var boundGroupID int64
	var boundAccountIDs []int64
	groupRepo := &groupRepoStubForAdmin{
		getByIDByID: map[int64]*Group{
			10: {ID: 10, Platform: PlatformOpenAIplaceholder,
			20: {ID: 20, Platform: PlatformGrokplaceholder,
			99: {ID: 99, Platform: PlatformComposite, RateMultiplier: 1, SubscriptionType: SubscriptionTypeStandardplaceholder,
	placeholder,
		deleteAccountGroupsByGroupIDFn: func(groupID int64) (int64, error) {
			clearedGroupID = groupID
			return 2, nil
	placeholder,
		getAccountIDsByGroupIDsFn: func(groupIDs []int64) ([]int64, error) {
			copiedFrom = append([]int64{placeholder, groupIDs...)
			return []int64{301, 302placeholder, nil
	placeholder,
		bindAccountsToGroupFn: func(groupID int64, accountIDs []int64) error {
			boundGroupID = groupID
			boundAccountIDs = append([]int64{placeholder, accountIDs...)
			return nil
	placeholder,
placeholder
	svc := &adminServiceImpl{groupRepo: groupRepoplaceholder
	maxReasoningEffort := "low"
	reasoningEffortMappings := []ReasoningEffortMapping{{From: "max", To: "high"placeholderplaceholder

	group, err := svc.UpdateGroup(context.Background(), 99, &UpdateGroupInput{
		MaxReasoningEffort:       &maxReasoningEffort,
		ReasoningEffortMappings:  &reasoningEffortMappings,
		CopyAccountsFromGroupIDs: []int64{10, 20placeholder,
placeholder)

placeholder
	require.Equal(t, PlatformComposite, group.Platform)
	require.Equal(t, "low", group.MaxReasoningEffort)
	require.Equal(t, reasoningEffortMappings, group.ReasoningEffortMappings)
	require.Equal(t, int64(99), clearedGroupID)
	require.ElementsMatch(t, []int64{10, 20placeholder, copiedFrom)
	require.Equal(t, int64(99), boundGroupID)
	require.ElementsMatch(t, []int64{301, 302placeholder, boundAccountIDs)
placeholder

func TestAdminService_CreateAccountAllowsCompositeGroupAssignment(t *testing.T) {
	accountRepo := &accountRepoStubForBulkUpdate{createID: 7placeholder
	groupRepo := &groupRepoStubForAdmin{
		getByIDByID: map[int64]*Group{
			99: {ID: 99, Platform: PlatformCompositeplaceholder,
	placeholder,
placeholder
	svc := &adminServiceImpl{accountRepo: accountRepo, groupRepo: groupRepoplaceholder

	account, err := svc.CreateAccount(context.Background(), &CreateAccountInput{
		Name:                  "OpenAI account",
		Platform:              PlatformOpenAI,
		Type:                  AccountTypeAPIKey,
		Concurrency:           1,
		GroupIDs:              []int64{99placeholder,
		SkipDefaultGroupBind:  true,
		SkipMixedChannelCheck: true,
placeholder)

placeholder
	require.Equal(t, int64(7), account.ID)
	require.Equal(t, PlatformOpenAI, accountRepo.createAccount.Platform)
	require.ElementsMatch(t, []int64{99placeholder, accountRepo.bindGroupsByAccount[7])
placeholder

func TestAdminService_UpdateAccountAllowsCompositeGroupAssignment(t *testing.T) {
	accountRepo := &accountRepoStubForBulkUpdate{
		getByIDAccounts: map[int64]*Account{
			7: {ID: 7, Platform: PlatformGemini, Type: AccountTypeAPIKey, Status: StatusActive, Extra: map[string]any{placeholderplaceholder,
	placeholder,
placeholder
	groupRepo := &groupRepoStubForAdmin{
		getByIDByID: map[int64]*Group{
			99: {ID: 99, Platform: PlatformCompositeplaceholder,
	placeholder,
placeholder
	svc := &adminServiceImpl{accountRepo: accountRepo, groupRepo: groupRepoplaceholder
	groupIDs := []int64{99placeholder

	account, err := svc.UpdateAccount(context.Background(), 7, &UpdateAccountInput{
		GroupIDs:              &groupIDs,
		SkipMixedChannelCheck: true,
placeholder)

placeholder
	require.Equal(t, int64(7), account.ID)
	require.Len(t, accountRepo.updatedAccounts, 1)
	require.ElementsMatch(t, []int64{99placeholder, accountRepo.bindGroupsByAccount[7])
placeholder

func TestAdminService_CompositeModelsListCandidatesIncludeConcreteAccountMappings(t *testing.T) {
	accountRepo := &accountRepoStubForCompositeModelsList{
		accounts: []Account{
			{
				ID:       1,
				Platform: PlatformOpenAI,
		placeholder
					"model_mapping": map[string]any{"gpt-custom": "gpt-5"placeholder,
			placeholder,
		placeholder,
			{
				ID:       2,
		placeholder
		placeholder
					"model_mapping": map[string]any{"gemini-custom": "gemini-2.5-flash"placeholder,
			placeholder,
		placeholder,
			{
				ID:       3,
				Platform: PlatformKimi,
		placeholder
					"model_mapping": map[string]any{"kimi-custom": "kimi-k2"placeholder,
			placeholder,
		placeholder,
	placeholder,
placeholder
	groupRepo := &groupRepoStubForAdmin{
		getByIDByID: map[int64]*Group{
			99: {ID: 99, Platform: PlatformCompositeplaceholder,
	placeholder,
placeholder
	svc := &adminServiceImpl{accountRepo: accountRepo, groupRepo: groupRepoplaceholder

	candidates, err := svc.GetGroupModelsListCandidates(context.Background(), 99, PlatformComposite)

placeholder
	require.Contains(t, candidates, "gpt-custom")
	require.Contains(t, candidates, "gemini-custom")
	require.Contains(t, candidates, "kimi-custom")
	require.Contains(t, candidates, "gpt-5.5")
	require.Contains(t, candidates, "gemini-2.5-flash")
placeholder

// 独立 CN 分组的模型列表候选沿用 default 分支的 Claude 默认列表；
// composite 支持不得改变独立分组的候选语义。
func TestAdminService_CNProviderModelsListCandidatesKeepClaudeDefaults(t *testing.T) {
	want := make([]string, 0, len(claude.DefaultModels))
	for _, model := range claude.DefaultModels {
		want = append(want, model.ID)
placeholder
	for _, platform := range []string{PlatformKimi, PlatformZhipu, PlatformDeepseekplaceholder {
		require.Equal(t, want, defaultModelsListCandidateIDs(platform), "platform=%s", platform)
placeholder
placeholder
