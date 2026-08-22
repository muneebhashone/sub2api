package service

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type guardianAffinityGroupRepo struct {
	GroupRepository
	group *Group
	err   error
placeholder

type guardianAffinityAccountRepo struct {
	schedulerGroupAwareOpenAIAccountRepo
	setErrorCalls int
placeholder

func (r *guardianAffinityAccountRepo) SetError(context.Context, int64, string) error {
	r.setErrorCalls++
	return nil
placeholder

func (r guardianAffinityGroupRepo) GetByID(context.Context, int64) (*Group, error) {
	if r.err != nil {
		return nil, r.err
placeholder
	return r.group, nil
placeholder

func (r guardianAffinityGroupRepo) GetByIDLite(context.Context, int64) (*Group, error) {
	if r.err != nil {
		return nil, r.err
placeholder
	return r.group, nil
placeholder

func guardianAffinityTestContext(t *testing.T, model, subagent, parentHeader, metadata string) context.Context {
placeholder
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", nil)
	c.Request.Header.Set(openAISubagentHeader, subagent)
	if parentHeader != "" {
		c.Request.Header.Set(codexParentThreadIDHeader, parentHeader)
placeholder
	if metadata != "" {
		c.Request.Header.Set(codexTurnMetadataHeader, metadata)
placeholder
	return WithOpenAIGuardianParentAffinity(context.Background(), c, nil, model)
placeholder

func TestWithOpenAIGuardianParentAffinity_RequiresUnambiguousReviewLineage(t *testing.T) {
	parentID := "11111111-1111-4111-8111-111111111111"
	wantHash := DeriveSessionHashFromSeed(parentID)

	for _, subagent := range []string{"guardian", "review", "GUARDIAN"placeholder {
		t.Run(subagent, func(t *testing.T) {
			ctx := guardianAffinityTestContext(t, codexAutoReviewModel, subagent, parentID, `{"parent_thread_id":"`+parentID+`"placeholder`)
			affinity, ok := openAIGuardianParentAffinityFromContext(ctx)
			require.True(t, ok)
			require.Equal(t, wantHash, affinity.currentSessionHash)
	placeholder)
placeholder

	t.Run("metadata only", func(t *testing.T) {
		ctx := guardianAffinityTestContext(t, codexAutoReviewModel, "guardian", "", `{"parent_thread_id":"`+parentID+`"placeholder`)
		_, ok := openAIGuardianParentAffinityFromContext(ctx)
		require.True(t, ok)
placeholder)

	t.Run("websocket envelope metadata", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodGet, "/openai/v1/responses", nil)
		body := []byte(`{"type":"response.create","response":{"model":"codex-auto-review","client_metadata":{"x-codex-turn-metadata":"{\"parent_thread_id\":\"` + parentID + `\",\"subagent_kind\":\"guardian\"placeholder"placeholderplaceholderplaceholder`)
		ctx := WithOpenAIGuardianParentAffinity(context.Background(), c, body, codexAutoReviewModel)
		affinity, ok := openAIGuardianParentAffinityFromContext(ctx)
		require.True(t, ok)
		require.Equal(t, wantHash, affinity.currentSessionHash)
placeholder)

	for name, ctx := range map[string]context.Context{
		"ordinary model":       guardianAffinityTestContext(t, "gpt-5.6-sol", "guardian", parentID, ""),
		"ordinary subagent":    guardianAffinityTestContext(t, codexAutoReviewModel, "collab_spawn", parentID, ""),
		"missing parent":       guardianAffinityTestContext(t, codexAutoReviewModel, "guardian", "", ""),
		"conflicting lineage":  guardianAffinityTestContext(t, codexAutoReviewModel, "guardian", parentID, `{"parent_thread_id":"different-parent"placeholder`),
		"conflicting subagent": guardianAffinityTestContext(t, codexAutoReviewModel, "guardian", parentID, `{"parent_thread_id":"`+parentID+`","subagent_kind":"collab_spawn"placeholder`),
placeholder {
		t.Run(name, func(t *testing.T) {
			_, ok := openAIGuardianParentAffinityFromContext(ctx)
			require.False(t, ok)
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayService_GuardianParentAffinitySelectsParentAccountAcrossSchedulers(t *testing.T) {
	parentID := "22222222-2222-4222-8222-222222222222"
	parentHash := DeriveSessionHashFromSeed(parentID)
	groupID := int64(102001)

	for _, mode := range []struct {
		name           string
		advanced       string
		stickyWeighted string
placeholder{
		{name: "legacy", advanced: "false"placeholder,
		{name: "advanced", advanced: "true"placeholder,
		{name: "advanced sticky weighted", advanced: "true", stickyWeighted: "true"placeholder,
placeholder {
		t.Run(mode.name, func(t *testing.T) {
			accounts := []Account{
				{
					ID: 39001, Platform: PlatformOpenAI, Type: AccountTypeOAuth,
					Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 10,
					GroupIDs: []int64{groupIDplaceholder, Credentials: map[string]any{"plan_type": "team"placeholder,
			placeholder,
				{
					ID: 39002, Platform: PlatformOpenAI, Type: AccountTypeOAuth,
					Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 0,
					GroupIDs: []int64{groupIDplaceholder, Credentials: map[string]any{"plan_type": "team"placeholder,
			placeholder,
		placeholder
			cfg := &config.Config{placeholder
			cfg.Gateway.OpenAIWS.LBTopK = 2
			cache := &schedulerTestGatewayCache{sessionBindings: map[string]int64{"openai:" + parentHash: 39001placeholderplaceholder
			svc := &OpenAIGatewayService{
				accountRepo:        schedulerGroupAwareOpenAIAccountRepo{schedulerTestOpenAIAccountRepo{accounts: accountsplaceholderplaceholder,
				cache:              cache,
				cfg:                cfg,
				rateLimitService:   newOpenAIAdvancedSchedulerRateLimitService(mode.advanced, mode.stickyWeighted),
				concurrencyService: NewConcurrencyService(schedulerTestConcurrencyCache{acquireResults: map[int64]bool{39001: true, 39002: trueplaceholderplaceholder),
		placeholder

			ctx := guardianAffinityTestContext(t, codexAutoReviewModel, "guardian", parentID, "")
			selection, decision, err := svc.SelectAccountWithScheduler(
				ctx, &groupID, "", "guardian-child-session", codexAutoReviewModel,
				nil, OpenAIUpstreamTransportAny, false,
			)
		placeholder
			require.NotNil(t, selection)
			require.Equal(t, int64(39001), selection.Account.ID)
			require.Equal(t, openAIAccountScheduleLayerGuardianParent, decision.Layer)
			require.Zero(t, cache.deletedSessions["openai:"+parentHash])
			if selection.ReleaseFunc != nil {
				selection.ReleaseFunc()
		placeholder
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayService_GuardianParentAffinityFallsBackWithoutCrossGroupOrFailoverBypass(t *testing.T) {
	parentID := "33333333-3333-4333-8333-333333333333"
	parentHash := DeriveSessionHashFromSeed(parentID)
	groupID := int64(102011)
	otherGroupID := int64(102012)

	for name, excluded := range map[string]map[int64]struct{placeholder{
		"parent moved out of group":              nil,
		"parent excluded after upstream failure": {39011: {placeholderplaceholder,
placeholder {
		t.Run(name, func(t *testing.T) {
			parentGroups := []int64{groupIDplaceholder
			if excluded == nil {
				parentGroups = []int64{otherGroupIDplaceholder
		placeholder
			accounts := []Account{
				{ID: 39011, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 0, GroupIDs: parentGroups, Credentials: map[string]any{"plan_type": "team"placeholderplaceholder,
				{ID: 39012, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 5, GroupIDs: []int64{groupIDplaceholder, Credentials: map[string]any{"plan_type": "team"placeholderplaceholder,
		placeholder
			cache := &schedulerTestGatewayCache{sessionBindings: map[string]int64{"openai:" + parentHash: 39011placeholderplaceholder
			svc := &OpenAIGatewayService{
				accountRepo:        schedulerGroupAwareOpenAIAccountRepo{schedulerTestOpenAIAccountRepo{accounts: accountsplaceholderplaceholder,
				cache:              cache,
				cfg:                &config.Config{placeholder,
				rateLimitService:   newOpenAIAdvancedSchedulerRateLimitService("true"),
				concurrencyService: NewConcurrencyService(schedulerTestConcurrencyCache{acquireResults: map[int64]bool{39011: true, 39012: trueplaceholderplaceholder),
		placeholder

			ctx := guardianAffinityTestContext(t, codexAutoReviewModel, "guardian", parentID, "")
			selection, _, err := svc.SelectAccountWithScheduler(
				ctx, &groupID, "", "guardian-fallback-child", codexAutoReviewModel,
				excluded, OpenAIUpstreamTransportAny, false,
			)
		placeholder
			require.NotNil(t, selection)
			require.Equal(t, int64(39012), selection.Account.ID)
			require.Zero(t, cache.deletedSessions["openai:"+parentHash], "a child request must never delete its parent's binding")
			if selection.ReleaseFunc != nil {
				selection.ReleaseFunc()
		placeholder
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayService_GuardianParentHashCollisionPreservesParentBinding(t *testing.T) {
	parentID := "44444444-4444-4444-8444-444444444444"
	parentHash := DeriveSessionHashFromSeed(parentID)
	groupID := int64(102021)
	otherGroupID := int64(102022)

	for _, advanced := range []string{"false", "true"placeholder {
		t.Run(map[string]string{"false": "legacy", "true": "advanced"placeholder[advanced], func(t *testing.T) {
			accounts := []Account{
				{ID: 39021, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 0, GroupIDs: []int64{otherGroupIDplaceholder, Credentials: map[string]any{"plan_type": "team"placeholderplaceholder,
				{ID: 39022, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 5, GroupIDs: []int64{groupIDplaceholder, Credentials: map[string]any{"plan_type": "team"placeholderplaceholder,
		placeholder
			cache := &schedulerTestGatewayCache{sessionBindings: map[string]int64{"openai:" + parentHash: 39021placeholderplaceholder
			svc := &OpenAIGatewayService{
				accountRepo:        schedulerGroupAwareOpenAIAccountRepo{schedulerTestOpenAIAccountRepo{accounts: accountsplaceholderplaceholder,
				cache:              cache,
				cfg:                &config.Config{placeholder,
				rateLimitService:   newOpenAIAdvancedSchedulerRateLimitService(advanced),
				concurrencyService: NewConcurrencyService(schedulerTestConcurrencyCache{acquireResults: map[int64]bool{39021: true, 39022: trueplaceholderplaceholder),
		placeholder

			ctx := guardianAffinityTestContext(t, codexAutoReviewModel, "guardian", parentID, "")
			selection, _, err := svc.SelectAccountWithScheduler(
				ctx, &groupID, "", parentHash, codexAutoReviewModel,
				nil, OpenAIUpstreamTransportAny, false,
			)
		placeholder
			require.NotNil(t, selection)
			require.Equal(t, int64(39022), selection.Account.ID)
			require.NoError(t, svc.BindStickySessionAfterProfitAdmission(ctx, &groupID, parentHash, selection.Account.ID))
			require.Equal(t, int64(39021), cache.sessionBindings["openai:"+parentHash])
			require.Zero(t, cache.deletedSessions["openai:"+parentHash])
			if selection.ReleaseFunc != nil {
				selection.ReleaseFunc()
		placeholder
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayService_GuardianParentAffinityHonorsRequiredPrivacy(t *testing.T) {
	parentID := "55555555-5555-4555-8555-555555555555"
	parentHash := DeriveSessionHashFromSeed(parentID)
	groupID := int64(102031)

	for _, advanced := range []string{"false", "true"placeholder {
		t.Run(map[string]string{"false": "legacy", "true": "advanced"placeholder[advanced], func(t *testing.T) {
			accounts := []Account{
				{ID: 39031, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 0, GroupIDs: []int64{groupIDplaceholder, Credentials: map[string]any{"plan_type": "team"placeholderplaceholder,
				{ID: 39032, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 5, GroupIDs: []int64{groupIDplaceholder, Credentials: map[string]any{"plan_type": "team"placeholder, Extra: map[string]any{"privacy_mode": PrivacyModeTrainingOffplaceholderplaceholder,
		placeholder
			repo := &guardianAffinityAccountRepo{schedulerGroupAwareOpenAIAccountRepo: schedulerGroupAwareOpenAIAccountRepo{schedulerTestOpenAIAccountRepo{accounts: accountsplaceholderplaceholderplaceholder
			cache := &schedulerTestGatewayCache{sessionBindings: map[string]int64{"openai:" + parentHash: 39031placeholderplaceholder
			svc := &OpenAIGatewayService{
				accountRepo:        repo,
				cache:              cache,
				cfg:                &config.Config{placeholder,
				rateLimitService:   newOpenAIAdvancedSchedulerRateLimitService(advanced),
				concurrencyService: NewConcurrencyService(schedulerTestConcurrencyCache{acquireResults: map[int64]bool{39031: true, 39032: trueplaceholderplaceholder),
				schedulerSnapshot: &SchedulerSnapshotService{
					accountRepo: repo,
					groupRepo:   guardianAffinityGroupRepo{group: &Group{ID: groupID, Name: "privacy", RequirePrivacySet: trueplaceholderplaceholder,
			placeholder,
		placeholder

			ctx := guardianAffinityTestContext(t, codexAutoReviewModel, "guardian", parentID, "")
			selection, _, err := svc.SelectAccountWithScheduler(
				ctx, &groupID, "", "guardian-privacy-child", codexAutoReviewModel,
				nil, OpenAIUpstreamTransportAny, false,
			)
		placeholder
			require.NotNil(t, selection)
			require.Equal(t, int64(39032), selection.Account.ID)
			require.Zero(t, repo.setErrorCalls, "a group-scoped privacy gate must not globally error a shared account")
			require.False(t, svc.isOpenAIAccountRequestRuntimeBlocked(&accounts[0], codexAutoReviewModel))
			if selection.ReleaseFunc != nil {
				selection.ReleaseFunc()
		placeholder
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayService_PreviousResponseHonorsGroupAndRequiredPrivacy(t *testing.T) {
	groupID := int64(3904)

	tests := []struct {
		name         string
		boundAccount Account
		groupErr     error
placeholder{
		{
			name: "privacy unset",
			boundAccount: Account{
				ID: 39041, Platform: PlatformOpenAI, Type: AccountTypeAPIKey,
				Status: StatusActive, Schedulable: true, Concurrency: 1,
				GroupIDs: []int64{groupIDplaceholder,
				Extra:    map[string]any{"openai_apikey_responses_websockets_v2_enabled": trueplaceholder,
		placeholder,
	placeholder,
		{
			name: "privacy policy lookup error fails closed",
			boundAccount: Account{
				ID: 39041, Platform: PlatformOpenAI, Type: AccountTypeAPIKey,
				Status: StatusActive, Schedulable: true, Concurrency: 1,
				GroupIDs: []int64{groupIDplaceholder,
				Extra:    map[string]any{"openai_apikey_responses_websockets_v2_enabled": trueplaceholder,
		placeholder,
			groupErr: errors.New("group repository unavailable"),
	placeholder,
		{
			name: "different group",
			boundAccount: Account{
				ID: 39041, Platform: PlatformOpenAI, Type: AccountTypeAPIKey,
				Status: StatusActive, Schedulable: true, Concurrency: 1,
				GroupIDs: []int64{groupID + 1placeholder,
				Extra: map[string]any{
					"openai_apikey_responses_websockets_v2_enabled": true,
					"privacy_mode": PrivacyModeTrainingOff,
			placeholder,
		placeholder,
	placeholder,
placeholder

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fallback := Account{
				ID: 39042, Platform: PlatformOpenAI, Type: AccountTypeAPIKey,
				Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 5,
				GroupIDs: []int64{groupIDplaceholder,
				Extra: map[string]any{
					"openai_apikey_responses_websockets_v2_enabled": true,
					"privacy_mode": PrivacyModeTrainingOff,
			placeholder,
		placeholder
			accounts := []Account{tc.boundAccount, fallbackplaceholder
			repo := &guardianAffinityAccountRepo{schedulerGroupAwareOpenAIAccountRepo: schedulerGroupAwareOpenAIAccountRepo{schedulerTestOpenAIAccountRepo{accounts: accountsplaceholderplaceholderplaceholder
			cache := &schedulerTestGatewayCache{placeholder
			store := NewOpenAIWSStateStore(cache)
			groupRepo := guardianAffinityGroupRepo{
				group: &Group{
					ID: groupID, Name: "privacy-required", Platform: PlatformOpenAI,
					Status: StatusActive, RequirePrivacySet: true,
			placeholder,
				err: tc.groupErr,
		placeholder
			svc := &OpenAIGatewayService{
				accountRepo:        repo,
				cache:              cache,
				cfg:                &config.Config{placeholder,
				rateLimitService:   newOpenAIAdvancedSchedulerRateLimitService("true"),
				concurrencyService: NewConcurrencyService(&schedulerTestConcurrencyCache{placeholder),
				openaiWSStateStore: store,
				schedulerSnapshot: &SchedulerSnapshotService{
					accountRepo: repo,
					groupRepo:   groupRepo,
			placeholder,
		placeholder
			responseID := "resp_privacy_guard"
			require.NoError(t, store.BindResponseAccount(context.Background(), groupID, responseID, tc.boundAccount.ID, time.Hour))

			directSelection, directErr := svc.SelectAccountByPreviousResponseID(
				context.Background(), &groupID, responseID, codexAutoReviewModel, nil, false,
			)
			require.NoError(t, directErr)
			require.Nil(t, directSelection, "the previous-response helper must enforce fresh group/privacy state")

			selection, decision, err := svc.SelectAccountWithSchedulerForCapability(
				context.Background(), &groupID, responseID, "", codexAutoReviewModel,
				nil, OpenAIUpstreamTransportAny, OpenAIEndpointCapabilityResponses,
				false, false, true,
			)
		placeholder
			require.NotNil(t, selection)
			require.Equal(t, fallback.ID, selection.Account.ID)
			require.NotEqual(t, openAIAccountScheduleLayerPreviousResponse, decision.Layer)
			require.Zero(t, repo.setErrorCalls)
			require.False(t, svc.isOpenAIAccountRequestRuntimeBlocked(&accounts[0], codexAutoReviewModel))
			boundAccountID, getErr := store.GetResponseAccount(context.Background(), groupID, responseID)
			require.NoError(t, getErr)
			require.Equal(t, tc.boundAccount.ID, boundAccountID, "transient policy misses must preserve the response binding")
			if selection.ReleaseFunc != nil {
				selection.ReleaseFunc()
		placeholder
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayService_PreviousResponseSimpleModeIgnoresGroupMembership(t *testing.T) {
	groupID := int64(3905)
	bound := Account{
		ID: 39051, Platform: PlatformOpenAI, Type: AccountTypeAPIKey,
		Status: StatusActive, Schedulable: true, Concurrency: 1,
		GroupIDs: []int64{groupID + 1placeholder,
		Extra:    map[string]any{"openai_apikey_responses_websockets_v2_enabled": trueplaceholder,
placeholder
	fallback := Account{
		ID: 39052, Platform: PlatformOpenAI, Type: AccountTypeAPIKey,
		Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 10,
		GroupIDs: []int64{groupIDplaceholder,
		Extra:    map[string]any{"openai_apikey_responses_websockets_v2_enabled": trueplaceholder,
placeholder
	accounts := []Account{bound, fallbackplaceholder
	repo := &guardianAffinityAccountRepo{schedulerGroupAwareOpenAIAccountRepo: schedulerGroupAwareOpenAIAccountRepo{schedulerTestOpenAIAccountRepo{accounts: accountsplaceholderplaceholderplaceholder
	cache := &schedulerTestGatewayCache{placeholder
	store := NewOpenAIWSStateStore(cache)
	cfg := &config.Config{RunMode: config.RunModeSimpleplaceholder
	svc := &OpenAIGatewayService{
		accountRepo:        repo,
		cache:              cache,
		cfg:                cfg,
		rateLimitService:   newOpenAIAdvancedSchedulerRateLimitService("true"),
		concurrencyService: NewConcurrencyService(&schedulerTestConcurrencyCache{placeholder),
		openaiWSStateStore: store,
		schedulerSnapshot: &SchedulerSnapshotService{
			accountRepo: repo,
			groupRepo: guardianAffinityGroupRepo{group: &Group{
				ID: groupID, Name: "simple-mode", Platform: PlatformOpenAI, Status: StatusActive,
	placeholder
	placeholder,
placeholder
	responseID := "resp_simple_mode_cross_group"
	require.NoError(t, store.BindResponseAccount(context.Background(), groupID, responseID, bound.ID, time.Hour))

	directSelection, err := svc.SelectAccountByPreviousResponseID(
		context.Background(), &groupID, responseID, codexAutoReviewModel, nil, false,
	)
placeholder
	require.NotNil(t, directSelection)
	require.Equal(t, bound.ID, directSelection.Account.ID)
	if directSelection.ReleaseFunc != nil {
		directSelection.ReleaseFunc()
placeholder

	selection, decision, err := svc.SelectAccountWithSchedulerForCapability(
		context.Background(), &groupID, responseID, "", codexAutoReviewModel,
		nil, OpenAIUpstreamTransportAny, OpenAIEndpointCapabilityResponses,
		false, false, true,
	)
placeholder
	require.NotNil(t, selection)
	require.Equal(t, bound.ID, selection.Account.ID)
	require.Equal(t, openAIAccountScheduleLayerPreviousResponse, decision.Layer)
	if selection.ReleaseFunc != nil {
		selection.ReleaseFunc()
placeholder
placeholder
