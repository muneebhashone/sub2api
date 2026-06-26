//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiagnoseModelAvailabilityForPlatform_NoModel_AlwaysAvailable(t *testing.T) {
	repo := &mockAccountRepoForPlatform{accounts: nil, accountsByID: map[int64]*Account{placeholderplaceholder
	svc := &GatewayService{accountRepo: repo, cfg: testConfig()placeholder

	diag := svc.DiagnoseModelAvailabilityForPlatform(context.Background(), nil, "", PlatformOpenAI)

	require.True(t, diag.HasAccountsInPool, "empty model must return HasAccountsInPool=true so caller stays on 503")
	require.True(t, diag.HasModelSupport, "empty model must return HasModelSupport=true so caller stays on 503")
placeholder

func TestDiagnoseModelAvailabilityForPlatform_EmptyPlatform_AlwaysAvailable(t *testing.T) {
	repo := &mockAccountRepoForPlatform{accounts: nil, accountsByID: map[int64]*Account{placeholderplaceholder
	svc := &GatewayService{accountRepo: repo, cfg: testConfig()placeholder

	diag := svc.DiagnoseModelAvailabilityForPlatform(context.Background(), nil, "gpt-5", "")

	require.True(t, diag.HasAccountsInPool)
	require.True(t, diag.HasModelSupport, "empty platform must fall back to {true,trueplaceholder so caller stays on 503")
placeholder

func TestDiagnoseModelAvailabilityForPlatform_NilReceiver(t *testing.T) {
	var svc *GatewayService

	diag := svc.DiagnoseModelAvailabilityForPlatform(context.Background(), nil, "gpt-5", PlatformOpenAI)

	require.True(t, diag.HasAccountsInPool)
	require.True(t, diag.HasModelSupport)
placeholder

func TestDiagnoseModelAvailabilityForPlatform_NoAccountsInPool(t *testing.T) {
	repo := &mockAccountRepoForPlatform{accounts: nil, accountsByID: map[int64]*Account{placeholderplaceholder
	svc := &GatewayService{accountRepo: repo, cfg: testConfig()placeholder

	diag := svc.DiagnoseModelAvailabilityForPlatform(context.Background(), nil, "gpt-5", PlatformOpenAI)

	require.False(t, diag.HasAccountsInPool)
	require.False(t, diag.HasModelSupport, "no accounts means no support; caller stays on 503 (empty-pool branch)")
placeholder

func TestDiagnoseModelAvailabilityForPlatform_ExplicitMappingMatches(t *testing.T) {
	repo := &mockAccountRepoForPlatform{
		accounts: []Account{
			{
				ID:          1,
				Platform:    PlatformOpenAI,
				Status:      StatusActive,
				Schedulable: true,
		placeholder
					"model_mapping": map[string]any{"gpt-5.1-codex-mini": "gpt-5.1-codex-mini"placeholder,
			placeholder,
		placeholder,
	placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder
	for i := range repo.accounts {
		repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
placeholder
	svc := &GatewayService{accountRepo: repo, cfg: testConfig()placeholder

	diag := svc.DiagnoseModelAvailabilityForPlatform(context.Background(), nil, "gpt-5.1-codex-mini", PlatformOpenAI)

	require.True(t, diag.HasAccountsInPool)
	require.True(t, diag.HasModelSupport)
placeholder

func TestDiagnoseModelAvailabilityForPlatform_EmptyMappingAllowsAll(t *testing.T) {
	repo := &mockAccountRepoForPlatform{
		accounts: []Account{
			{ID: 1, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true /* no ModelMapping = allow all */placeholder,
	placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder
	for i := range repo.accounts {
		repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
placeholder
	svc := &GatewayService{accountRepo: repo, cfg: testConfig()placeholder

	diag := svc.DiagnoseModelAvailabilityForPlatform(context.Background(), nil, "gpt-5.1-codex-mini", PlatformOpenAI)

	require.True(t, diag.HasModelSupport, "empty model_mapping must be treated as 'allow all' (Account.IsModelSupported semantics)")
placeholder

func TestDiagnoseModelAvailabilityForPlatform_WildcardMappingMatches(t *testing.T) {
	repo := &mockAccountRepoForPlatform{
		accounts: []Account{
			{
				ID:          1,
				Platform:    PlatformOpenAI,
				Status:      StatusActive,
				Schedulable: true,
		placeholder
					"model_mapping": map[string]any{"*": "gpt-5"placeholder,
			placeholder,
		placeholder,
	placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder
	for i := range repo.accounts {
		repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
placeholder
	svc := &GatewayService{accountRepo: repo, cfg: testConfig()placeholder

	diag := svc.DiagnoseModelAvailabilityForPlatform(context.Background(), nil, "gpt-5.1-codex-mini", PlatformOpenAI)

	require.True(t, diag.HasModelSupport, "wildcard mapping must classify the request as 'serviceable'")
placeholder

func TestDiagnoseModelAvailabilityForPlatform_NoMatchingModel_ReturnsNotFoundSignal(t *testing.T) {
	repo := &mockAccountRepoForPlatform{
		accounts: []Account{
			{
				ID:          1,
				Platform:    PlatformOpenAI,
				Status:      StatusActive,
				Schedulable: true,
		placeholder"model_mapping": map[string]any{"gpt-5": "gpt-5"placeholderplaceholder,
		placeholder,
			{
				ID:          2,
				Platform:    PlatformOpenAI,
				Status:      StatusActive,
				Schedulable: true,
		placeholder"model_mapping": map[string]any{"gpt-5-mini": "gpt-5-mini"placeholderplaceholder,
		placeholder,
	placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder
	for i := range repo.accounts {
		repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
placeholder
	svc := &GatewayService{accountRepo: repo, cfg: testConfig()placeholder

	diag := svc.DiagnoseModelAvailabilityForPlatform(context.Background(), nil, "gpt-5.1-codex-mini", PlatformOpenAI)

	require.True(t, diag.HasAccountsInPool, "group has OpenAI accounts")
	require.False(t, diag.HasModelSupport, "no account mapping admits the requested model — handler should return 404")
placeholder

func TestDiagnoseModelAvailabilityForPlatform_WrongPlatformFiltersOut(t *testing.T) {
	// Group has only Anthropic accounts; user routes to OpenAI gateway.
	// Diagnosis must NOT see Anthropic accounts (listSchedulableAccounts filters
	// by platform), so HasAccountsInPool is false and the caller stays on 503.
	repo := &mockAccountRepoForPlatform{
		accounts: []Account{
			{
				ID:          1,
				Platform:    PlatformAnthropic,
				Status:      StatusActive,
				Schedulable: true,
		placeholder"model_mapping": map[string]any{"claude-sonnet-4-5": "claude-sonnet-4-5"placeholderplaceholder,
		placeholder,
	placeholder,
		accountsByID: map[int64]*Account{placeholder,
placeholder
	for i := range repo.accounts {
		repo.accountsByID[repo.accounts[i].ID] = &repo.accounts[i]
placeholder
	svc := &GatewayService{accountRepo: repo, cfg: testConfig()placeholder

	diag := svc.DiagnoseModelAvailabilityForPlatform(context.Background(), nil, "gpt-5", PlatformOpenAI)

	require.False(t, diag.HasAccountsInPool, "OpenAI route must not see Anthropic accounts in pool")
	require.False(t, diag.HasModelSupport)
placeholder
