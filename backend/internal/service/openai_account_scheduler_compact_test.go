package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

// TestOpenAIGatewayService_SelectAccountWithScheduler_CompactPrefersSupportedOverUnknown
// 验证 compact 调度时显式支持 (tier=2) 优先于未探测 (tier=1)。
func TestOpenAIGatewayService_SelectAccountWithScheduler_CompactPrefersSupportedOverUnknown(t *testing.T) {
	resetOpenAIAdvancedSchedulerSettingCacheForTest()

	ctx := context.Background()
	groupID := int64(91001)
	accounts := []Account{
		{
			ID:          71001,
			Platform:    PlatformOpenAI,
			Type:        AccountTypeAPIKey,
			Status:      StatusActive,
			Schedulable: true,
			Concurrency: 1,
			Priority:    0,
			Extra:       map[string]any{placeholder, // unknown
	placeholder,
		{
			ID:          71002,
			Platform:    PlatformOpenAI,
			Type:        AccountTypeAPIKey,
			Status:      StatusActive,
			Schedulable: true,
			Concurrency: 1,
			Priority:    0,
			Extra:       map[string]any{"openai_compact_supported": trueplaceholder, // tier=2
	placeholder,
placeholder
	cfg := &config.Config{placeholder
	cfg.Gateway.Scheduling.LoadBatchEnabled = false
	svc := &OpenAIGatewayService{
		accountRepo:        schedulerTestOpenAIAccountRepo{accounts: accountsplaceholder,
		cache:              &schedulerTestGatewayCache{placeholder,
		cfg:                cfg,
		concurrencyService: NewConcurrencyService(schedulerTestConcurrencyCache{placeholder),
placeholder

	selection, _, err := svc.SelectAccountWithScheduler(
		ctx,
		&groupID,
		"",
		"",
		"gpt-5.4",
		nil,
		OpenAIUpstreamTransportAny,
		true,
	)
placeholder
	require.NotNil(t, selection)
	require.NotNil(t, selection.Account)
	require.Equal(t, int64(71002), selection.Account.ID, "compact-supported account should win over unknown")
placeholder

// TestOpenAIGatewayService_SelectAccountWithScheduler_CompactRejectsExplicitlyUnsupported
// 验证 force_off / 已探测不支持 (tier=0) 的账号不会被 compact 请求选中。
func TestOpenAIGatewayService_SelectAccountWithScheduler_CompactRejectsExplicitlyUnsupported(t *testing.T) {
	resetOpenAIAdvancedSchedulerSettingCacheForTest()

	ctx := context.Background()
	groupID := int64(91002)
	accounts := []Account{
		{
			ID:          71010,
			Platform:    PlatformOpenAI,
			Type:        AccountTypeAPIKey,
			Status:      StatusActive,
			Schedulable: true,
			Concurrency: 1,
			Priority:    0,
			Extra:       map[string]any{"openai_compact_mode": OpenAICompactModeForceOffplaceholder,
	placeholder,
		{
			ID:          71011,
			Platform:    PlatformOpenAI,
			Type:        AccountTypeAPIKey,
			Status:      StatusActive,
			Schedulable: true,
			Concurrency: 1,
			Priority:    0,
			Extra:       map[string]any{"openai_compact_supported": falseplaceholder,
	placeholder,
placeholder
	cfg := &config.Config{placeholder
	cfg.Gateway.Scheduling.LoadBatchEnabled = false
	svc := &OpenAIGatewayService{
		accountRepo:        schedulerTestOpenAIAccountRepo{accounts: accountsplaceholder,
		cache:              &schedulerTestGatewayCache{placeholder,
		cfg:                cfg,
		concurrencyService: NewConcurrencyService(schedulerTestConcurrencyCache{placeholder),
placeholder

	selection, _, err := svc.SelectAccountWithScheduler(
		ctx,
		&groupID,
		"",
		"",
		"gpt-5.4",
		nil,
		OpenAIUpstreamTransportAny,
		true,
	)
placeholder
	require.True(t, errors.Is(err, ErrNoAvailableCompactAccounts), "compact-only accounts should rejected explicitly unsupported and return compact error")
	require.Nil(t, selection)
placeholder

func newOpenAICompactionSchedulerTestService(accounts []Account, advanced bool) *OpenAIGatewayService {
	cfg := &config.Config{placeholder
	cfg.Gateway.Scheduling.LoadBatchEnabled = false
	svc := &OpenAIGatewayService{
		accountRepo:        schedulerTestOpenAIAccountRepo{accounts: accountsplaceholder,
		cache:              &schedulerTestGatewayCache{placeholder,
		cfg:                cfg,
		concurrencyService: NewConcurrencyService(schedulerTestConcurrencyCache{placeholder),
placeholder
	if advanced {
		svc.rateLimitService = newOpenAIAdvancedSchedulerRateLimitService("true")
placeholder
	return svc
placeholder

func selectOpenAICompactionSchedulerTestAccount(t *testing.T, svc *OpenAIGatewayService, groupID int64, requireCompact bool) (*AccountSelectionResult, error) {
placeholder
	selection, _, err := svc.SelectAccountWithSchedulerForCapability(
		context.Background(),
		&groupID,
		"",
		"",
		"gpt-5.6-sol",
		nil,
		OpenAIUpstreamTransportAny,
		OpenAIEndpointCapabilityResponses,
		requireCompact,
		false,
		false,
	)
	return selection, err
placeholder

func TestOpenAIGatewayService_SelectAccountWithScheduler_NativeCompactionIgnoresLegacyCompactProbe(t *testing.T) {
	for _, advanced := range []bool{false, trueplaceholder {
		t.Run(map[bool]string{false: "legacy_scheduler", true: "advanced_scheduler"placeholder[advanced], func(t *testing.T) {
			resetOpenAIAdvancedSchedulerSettingCacheForTest()
			svc := newOpenAICompactionSchedulerTestService([]Account{{
				ID:          71012,
				Platform:    PlatformOpenAI,
				Type:        AccountTypeAPIKey,
				Status:      StatusActive,
				Schedulable: true,
				Concurrency: 1,
				Extra: map[string]any{
					"openai_compact_supported":   false,
					"openai_responses_supported": true,
			placeholder,
	placeholder advanced)

			selection, err := selectOpenAICompactionSchedulerTestAccount(t, svc, 91007, false)
		placeholder
			require.NotNil(t, selection)
			require.Equal(t, int64(71012), selection.Account.ID)
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayService_SelectAccountWithScheduler_NativeCompactionAllowsForceOff(t *testing.T) {
	for _, advanced := range []bool{false, trueplaceholder {
		t.Run(map[bool]string{false: "legacy_scheduler", true: "advanced_scheduler"placeholder[advanced], func(t *testing.T) {
			resetOpenAIAdvancedSchedulerSettingCacheForTest()
			svc := newOpenAICompactionSchedulerTestService([]Account{{
				ID:          71013,
				Platform:    PlatformOpenAI,
				Type:        AccountTypeAPIKey,
				Status:      StatusActive,
				Schedulable: true,
				Concurrency: 1,
				Extra: map[string]any{
					"openai_compact_mode":        OpenAICompactModeForceOff,
					"openai_responses_supported": true,
			placeholder,
	placeholder advanced)

			selection, err := selectOpenAICompactionSchedulerTestAccount(t, svc, 91008, false)
		placeholder
			require.NotNil(t, selection)
			require.Equal(t, int64(71013), selection.Account.ID)
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayService_SelectAccountWithScheduler_NativeCompactionRequiresResponsesCapability(t *testing.T) {
	for _, advanced := range []bool{false, trueplaceholder {
		t.Run(map[bool]string{false: "legacy_scheduler", true: "advanced_scheduler"placeholder[advanced], func(t *testing.T) {
			resetOpenAIAdvancedSchedulerSettingCacheForTest()
			svc := newOpenAICompactionSchedulerTestService([]Account{{
				ID:          71014,
				Platform:    PlatformOpenAI,
				Type:        AccountTypeAPIKey,
				Status:      StatusActive,
				Schedulable: true,
				Concurrency: 1,
				Extra: map[string]any{
					"openai_compact_mode":        OpenAICompactModeForceOn,
					"openai_responses_supported": false,
			placeholder,
	placeholder advanced)

			selection, err := selectOpenAICompactionSchedulerTestAccount(t, svc, 91009, false)
			require.ErrorIs(t, err, ErrNoAvailableAccounts)
			require.NotErrorIs(t, err, ErrNoAvailableCompactAccounts)
			require.NotContains(t, err.Error(), "/responses/compact")
			require.Nil(t, selection)
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayService_SelectAccountWithScheduler_LegacyCompactionKeepsCompactEligibility(t *testing.T) {
	for _, advanced := range []bool{false, trueplaceholder {
		t.Run(map[bool]string{false: "legacy_scheduler", true: "advanced_scheduler"placeholder[advanced], func(t *testing.T) {
			resetOpenAIAdvancedSchedulerSettingCacheForTest()
			svc := newOpenAICompactionSchedulerTestService([]Account{
				{
					ID:          71015,
					Platform:    PlatformOpenAI,
					Type:        AccountTypeAPIKey,
					Status:      StatusActive,
					Schedulable: true,
					Concurrency: 1,
					Extra: map[string]any{
						"openai_compact_supported":   false,
						"openai_responses_supported": true,
				placeholder,
			placeholder,
				{
					ID:          71016,
					Platform:    PlatformOpenAI,
					Type:        AccountTypeAPIKey,
					Status:      StatusActive,
					Schedulable: true,
					Concurrency: 1,
					Extra: map[string]any{
						"openai_compact_mode":        OpenAICompactModeForceOff,
						"openai_responses_supported": true,
				placeholder,
			placeholder,
		placeholder, advanced)

			selection, err := selectOpenAICompactionSchedulerTestAccount(t, svc, 91010, true)
			require.ErrorIs(t, err, ErrNoAvailableCompactAccounts)
			require.Contains(t, err.Error(), "/responses/compact")
			require.Nil(t, selection)
	placeholder)
placeholder
placeholder

// TestOpenAIGatewayService_SelectAccountWithScheduler_CompactRequiresResponsesCapability
// prevents a confirmed Chat Completions-only API key from silently dropping the
// compaction trigger through the Responses-to-Chat fallback.
func TestOpenAIGatewayService_SelectAccountWithScheduler_CompactRequiresResponsesCapability(t *testing.T) {
	resetOpenAIAdvancedSchedulerSettingCacheForTest()

	ctx := context.Background()
	groupID := int64(91005)
	accounts := []Account{{
		ID:          71050,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Schedulable: true,
		Concurrency: 1,
		Extra: map[string]any{
			"openai_compact_supported":   true,
			"openai_responses_supported": false,
	placeholder,
placeholderplaceholder
	cfg := &config.Config{placeholder
	cfg.Gateway.Scheduling.LoadBatchEnabled = false
	svc := &OpenAIGatewayService{
		accountRepo:        schedulerTestOpenAIAccountRepo{accounts: accountsplaceholder,
		cache:              &schedulerTestGatewayCache{placeholder,
		cfg:                cfg,
		concurrencyService: NewConcurrencyService(schedulerTestConcurrencyCache{placeholder),
placeholder

	selection, _, err := svc.SelectAccountWithSchedulerForCapability(
		ctx,
		&groupID,
		"",
		"",
		"gpt-5.4",
		nil,
		OpenAIUpstreamTransportAny,
		OpenAIEndpointCapabilityResponses,
		true,
		false,
		false,
	)
placeholder
	require.True(t, errors.Is(err, ErrNoAvailableAccounts))
	require.Nil(t, selection)
placeholder

func TestOpenAIGatewayService_SelectAccountWithScheduler_CompactSkipsChatOnlyAccount(t *testing.T) {
	resetOpenAIAdvancedSchedulerSettingCacheForTest()

	ctx := context.Background()
	groupID := int64(91006)
	accounts := []Account{
		{
			ID:          71060,
			Platform:    PlatformOpenAI,
			Type:        AccountTypeAPIKey,
			Status:      StatusActive,
			Schedulable: true,
			Concurrency: 1,
			Priority:    10,
			Extra: map[string]any{
				"openai_compact_supported":   true,
				"openai_responses_supported": false,
		placeholder,
	placeholder,
		{
			ID:          71061,
			Platform:    PlatformOpenAI,
			Type:        AccountTypeAPIKey,
			Status:      StatusActive,
			Schedulable: true,
			Concurrency: 1,
			Priority:    0,
			Extra: map[string]any{
				"openai_compact_supported":   true,
				"openai_responses_supported": true,
		placeholder,
	placeholder,
placeholder
	cfg := &config.Config{placeholder
	cfg.Gateway.Scheduling.LoadBatchEnabled = false
	svc := &OpenAIGatewayService{
		accountRepo:        schedulerTestOpenAIAccountRepo{accounts: accountsplaceholder,
		cache:              &schedulerTestGatewayCache{placeholder,
		cfg:                cfg,
		concurrencyService: NewConcurrencyService(schedulerTestConcurrencyCache{placeholder),
placeholder

	selection, _, err := svc.SelectAccountWithSchedulerForCapability(
		ctx,
		&groupID,
		"",
		"",
		"gpt-5.6-sol",
		nil,
		OpenAIUpstreamTransportAny,
		OpenAIEndpointCapabilityResponses,
		true,
		false,
		false,
	)
placeholder
	require.NotNil(t, selection)
	require.NotNil(t, selection.Account)
	require.Equal(t, int64(71061), selection.Account.ID)
placeholder

// TestOpenAIGatewayService_SelectAccountWithScheduler_CompactFallsBackToUnknown
// 验证当没有"已知支持"账号时，compact 请求会回退到"未探测"账号。
func TestOpenAIGatewayService_SelectAccountWithScheduler_CompactFallsBackToUnknown(t *testing.T) {
	resetOpenAIAdvancedSchedulerSettingCacheForTest()

	ctx := context.Background()
	groupID := int64(91003)
	accounts := []Account{
		{
			ID:          71020,
			Platform:    PlatformOpenAI,
			Type:        AccountTypeAPIKey,
			Status:      StatusActive,
			Schedulable: true,
			Concurrency: 1,
			Priority:    0,
			Extra:       map[string]any{"openai_compact_supported": falseplaceholder, // tier=0
	placeholder,
		{
			ID:          71021,
			Platform:    PlatformOpenAI,
			Type:        AccountTypeAPIKey,
			Status:      StatusActive,
			Schedulable: true,
			Concurrency: 1,
			Priority:    0,
			Extra:       map[string]any{placeholder, // unknown -> tier=1
	placeholder,
placeholder
	cfg := &config.Config{placeholder
	cfg.Gateway.Scheduling.LoadBatchEnabled = false
	svc := &OpenAIGatewayService{
		accountRepo:        schedulerTestOpenAIAccountRepo{accounts: accountsplaceholder,
		cache:              &schedulerTestGatewayCache{placeholder,
		cfg:                cfg,
		concurrencyService: NewConcurrencyService(schedulerTestConcurrencyCache{placeholder),
placeholder

	selection, _, err := svc.SelectAccountWithScheduler(
		ctx,
		&groupID,
		"",
		"",
		"gpt-5.4",
		nil,
		OpenAIUpstreamTransportAny,
		true,
	)
placeholder
	require.NotNil(t, selection)
	require.NotNil(t, selection.Account)
	require.Equal(t, int64(71021), selection.Account.ID, "unknown account should be picked when no supported account available")
placeholder

// TestOpenAIGatewayService_SelectAccountWithScheduler_CompactAllowsGrok verifies
// that OpenAI-compatible compact routing does not reject Grok accounts.
func TestOpenAIGatewayService_SelectAccountWithScheduler_CompactAllowsGrok(t *testing.T) {
	resetOpenAIAdvancedSchedulerSettingCacheForTest()

	ctx := context.Background()
	groupID := int64(91004)
	accounts := []Account{
		{
			ID:          71030,
			Platform:    PlatformGrok,
			Type:        AccountTypeOAuth,
			Status:      StatusActive,
			Schedulable: true,
			Concurrency: 1,
			Priority:    0,
	placeholder
				"model_mapping": map[string]any{"grok-4.5": "grok-4.5"placeholder,
		placeholder,
	placeholder,
placeholder
	cfg := &config.Config{placeholder
	cfg.Gateway.Scheduling.LoadBatchEnabled = false
	svc := &OpenAIGatewayService{
		accountRepo:        schedulerTestOpenAIAccountRepo{accounts: accountsplaceholder,
		cache:              &schedulerTestGatewayCache{placeholder,
		cfg:                cfg,
		concurrencyService: NewConcurrencyService(schedulerTestConcurrencyCache{placeholder),
placeholder

	selection, _, err := svc.SelectAccountWithSchedulerForCapability(
		ctx,
		&groupID,
		"",
		"",
		"grok-4.5",
		nil,
		OpenAIUpstreamTransportAny,
		OpenAIEndpointCapabilityChatCompletions,
		true,
		false,
		true,
		PlatformGrok,
	)
placeholder
	require.NotNil(t, selection)
	require.NotNil(t, selection.Account)
	require.Equal(t, int64(71030), selection.Account.ID)
placeholder

// TestOpenAICompactSupportTier 验证 tier 分类逻辑。
func TestOpenAICompactSupportTier(t *testing.T) {
	tests := []struct {
		name    string
		account *Account
		want    int
placeholder{
		{name: "nil", account: nil, want: 0placeholder,
		{name: "non openai", account: &Account{Platform: PlatformAnthropicplaceholder, want: 0placeholder,
		{name: "grok", account: &Account{Platform: PlatformGrokplaceholder, want: 2placeholder,
		{name: "openai unknown", account: &Account{Platform: PlatformOpenAI, Extra: map[string]any{placeholderplaceholder, want: 1placeholder,
		{name: "openai supported", account: &Account{Platform: PlatformOpenAI, Extra: map[string]any{"openai_compact_supported": trueplaceholderplaceholder, want: 2placeholder,
		{name: "openai unsupported", account: &Account{Platform: PlatformOpenAI, Extra: map[string]any{"openai_compact_supported": falseplaceholderplaceholder, want: 0placeholder,
		{name: "force on", account: &Account{Platform: PlatformOpenAI, Extra: map[string]any{"openai_compact_mode": OpenAICompactModeForceOnplaceholderplaceholder, want: 2placeholder,
		{name: "force off overrides probe true", account: &Account{Platform: PlatformOpenAI, Extra: map[string]any{"openai_compact_mode": OpenAICompactModeForceOff, "openai_compact_supported": trueplaceholderplaceholder, want: 0placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := openAICompactSupportTier(tt.account); got != tt.want {
				t.Fatalf("openAICompactSupportTier(...) = %d, want %d", got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder
