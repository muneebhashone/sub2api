package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/stretchr/testify/require"
)

type upstreamBillingProbeAdminRepo struct {
	*upstreamBillingProbeAccountRepo
placeholder

func (r *upstreamBillingProbeAdminRepo) ListShadowsByParent(context.Context, int64) ([]*Account, error) {
	return nil, nil
placeholder

func TestCreateAccountDropsManagedUpstreamBillingProbeState(t *testing.T) {
	repo := &upstreamBillingProbeAccountRepo{placeholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder

	created, err := svc.CreateAccount(context.Background(), &CreateAccountInput{
		Name:                 "upstream",
		Platform:             PlatformOpenAI,
		Type:                 AccountTypeAPIKey,
		Credentials:          map[string]any{"api_key": "sk-test"placeholder,
		SkipDefaultGroupBind: true,
		Extra: map[string]any{
			UpstreamBillingProbeEnabledExtraKey: true,
			UpstreamBillingProbeExtraKey:        map[string]any{"status": "ok"placeholder,
	placeholder,
placeholder)

placeholder
	require.NotContains(t, created.Extra, UpstreamBillingProbeEnabledExtraKey)
	require.NotContains(t, created.Extra, UpstreamBillingProbeExtraKey)
placeholder

func TestUpdateAccountPreservesManagedUpstreamBillingProbeStateForUnrelatedEdit(t *testing.T) {
	accountID := int64(110)
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{
		accountID: {
			ID:       accountID,
			Platform: PlatformOpenAI,
			Type:     AccountTypeAPIKey,
			Status:   StatusActive,
			Extra: map[string]any{
				UpstreamBillingProbeEnabledExtraKey: true,
				UpstreamBillingProbeExtraKey:        map[string]any{"status": "ok"placeholder,
		placeholder,
	placeholder,
placeholderplaceholder

	svc := &adminServiceImpl{accountRepo: repoplaceholder
	updated, err := svc.UpdateAccount(context.Background(), accountID, &UpdateAccountInput{
		Extra: map[string]any{"custom": "value"placeholder,
placeholder)

placeholder
	require.Equal(t, true, updated.Extra[UpstreamBillingProbeEnabledExtraKey])
	require.Contains(t, updated.Extra, UpstreamBillingProbeExtraKey)
	require.Equal(t, "value", updated.Extra["custom"])
placeholder

func TestUpdateAccountPreservesGrokBillingSnapshotForUnrelatedEdit(t *testing.T) {
	accountID := int64(112)
	billing := &xai.BillingSummary{
		StatusCode:       http.StatusForbidden,
		WeeklyStatusCode: http.StatusForbidden,
placeholder
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{
		accountID: {
			ID:       accountID,
			Platform: PlatformGrok,
			Type:     AccountTypeOAuth,
			Status:   StatusActive,
			Extra:    map[string]any{grokBillingExtraKey: billingplaceholder,
	placeholder,
placeholderplaceholder

	updated, err := (&adminServiceImpl{accountRepo: repoplaceholder).UpdateAccount(context.Background(), accountID, &UpdateAccountInput{
		Extra: map[string]any{"custom": "value"placeholder,
placeholder)

placeholder
	require.Equal(t, billing, updated.Extra[grokBillingExtraKey])
	require.Equal(t, "value", updated.Extra["custom"])
	eligible, reason := updated.GrokMediaGenerationEligibility()
	require.False(t, eligible)
	require.Equal(t, "billing_forbidden", reason)
placeholder

func TestUpdateAccountPreservesProbeSnapshotWhenIdentityValuesAreUnchanged(t *testing.T) {
	accountID := int64(119)
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{
		accountID: {
			ID:       accountID,
			Platform: PlatformOpenAI,
			Type:     AccountTypeAPIKey,
			Status:   StatusActive,
	placeholder
				"api_key":                    "sk-existing",
				"base_url":                   "https://upstream.example",
				credKeyHeaderOverrideEnabled: true,
				credKeyHeaderOverrides:       map[string]any{"x-route": "stable"placeholder,
		placeholder,
			Extra: map[string]any{
				UpstreamBillingProbeEnabledExtraKey: true,
				UpstreamBillingProbeExtraKey:        map[string]any{"status": "ok"placeholder,
		placeholder,
	placeholder,
placeholderplaceholder

	updated, err := (&adminServiceImpl{accountRepo: repoplaceholder).UpdateAccount(context.Background(), accountID, &UpdateAccountInput{
placeholder
			"base_url":                   "https://upstream.example",
			credKeyHeaderOverrideEnabled: true,
			credKeyHeaderOverrides:       map[string]any{"x-route": "stable"placeholder,
	placeholder,
placeholder)

placeholder
	require.Contains(t, updated.Extra, UpstreamBillingProbeExtraKey)
placeholder

func TestUpdateAccountInvalidatesProbeSnapshotWhenUpstreamIdentityChanges(t *testing.T) {
	tests := []struct {
		name        string
		input       *UpdateAccountInput
		wantEnabled bool
placeholder{
		{
			name:        "api key",
			input:       &UpdateAccountInput{Credentials: map[string]any{"api_key": "sk-new"placeholderplaceholder,
			wantEnabled: true,
	placeholder,
		{
			name:        "base url",
			input:       &UpdateAccountInput{Credentials: map[string]any{"base_url": "https://new.example"placeholderplaceholder,
			wantEnabled: true,
	placeholder,
		{
			name: "header override",
			input: &UpdateAccountInput{Credentials: map[string]any{
				credKeyHeaderOverrideEnabled: true,
				credKeyHeaderOverrides:       map[string]any{"x-route": "new"placeholder,
	placeholder
			wantEnabled: true,
	placeholder,
		{
			name:        "account type",
			input:       &UpdateAccountInput{Type: AccountTypeOAuthplaceholder,
			wantEnabled: false,
	placeholder,
placeholder

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountID := int64(120 + i)
			repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{
				accountID: {
					ID:       accountID,
					Platform: PlatformOpenAI,
					Type:     AccountTypeAPIKey,
					Status:   StatusActive,
			placeholder
						"api_key":  "sk-old",
						"base_url": "https://old.example",
				placeholder,
					Extra: map[string]any{
						UpstreamBillingProbeEnabledExtraKey: true,
						UpstreamBillingProbeExtraKey:        map[string]any{"status": "ok"placeholder,
				placeholder,
			placeholder,
		placeholderplaceholder

			updated, err := (&adminServiceImpl{accountRepo: repoplaceholder).UpdateAccount(context.Background(), accountID, tt.input)

		placeholder
			require.NotContains(t, updated.Extra, UpstreamBillingProbeExtraKey)
			if tt.wantEnabled {
				require.Equal(t, true, updated.Extra[UpstreamBillingProbeEnabledExtraKey])
		placeholder else {
				require.NotContains(t, updated.Extra, UpstreamBillingProbeEnabledExtraKey)
		placeholder
	placeholder)
placeholder
placeholder

func TestUpdateAccountInvalidatesProbeSnapshotWhenProxyChanges(t *testing.T) {
	accountID := int64(140)
	oldProxyID := int64(7)
	newProxyID := int64(8)
	baseRepo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{
		accountID: {
			ID:          accountID,
			Platform:    PlatformOpenAI,
			Type:        AccountTypeAPIKey,
			Status:      StatusActive,
	placeholder"api_key": "sk-test"placeholder,
			ProxyID:     &oldProxyID,
			Extra: map[string]any{
				UpstreamBillingProbeEnabledExtraKey: true,
				UpstreamBillingProbeExtraKey:        map[string]any{"status": "ok"placeholder,
		placeholder,
	placeholder,
placeholderplaceholder

	updated, err := (&adminServiceImpl{accountRepo: &upstreamBillingProbeAdminRepo{baseRepoplaceholderplaceholder).UpdateAccount(
		context.Background(),
		accountID,
		&UpdateAccountInput{ProxyID: &newProxyIDplaceholder,
	)

placeholder
	require.Equal(t, newProxyID, *updated.ProxyID)
	require.NotContains(t, updated.Extra, UpstreamBillingProbeExtraKey)
placeholder

func TestUpdateAccountPreservesProbeSnapshotWhenProxyIsUnchanged(t *testing.T) {
	accountID := int64(141)
	existingProxyID := int64(7)
	unchangedProxyID := int64(7)
	baseRepo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{
		accountID: {
			ID:          accountID,
			Platform:    PlatformOpenAI,
			Type:        AccountTypeAPIKey,
			Status:      StatusActive,
	placeholder"api_key": "sk-test"placeholder,
			ProxyID:     &existingProxyID,
			Extra: map[string]any{
				UpstreamBillingProbeEnabledExtraKey: true,
				UpstreamBillingProbeExtraKey:        map[string]any{"status": "ok"placeholder,
		placeholder,
	placeholder,
placeholderplaceholder

	updated, err := (&adminServiceImpl{accountRepo: &upstreamBillingProbeAdminRepo{baseRepoplaceholderplaceholder).UpdateAccount(
		context.Background(),
		accountID,
		&UpdateAccountInput{ProxyID: &unchangedProxyIDplaceholder,
	)

placeholder
	require.Contains(t, updated.Extra, UpstreamBillingProbeExtraKey)
placeholder

func TestUpdateAccountAcceptsProbeEnabledAndRejectsInjectedSnapshot(t *testing.T) {
	accountID := int64(111)
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{
		accountID: {
			ID:       accountID,
			Platform: PlatformOpenAI,
			Type:     AccountTypeAPIKey,
			Status:   StatusActive,
			Extra:    map[string]any{placeholder,
	placeholder,
placeholderplaceholder

	svc := &adminServiceImpl{accountRepo: repoplaceholder
	updated, err := svc.UpdateAccount(context.Background(), accountID, &UpdateAccountInput{
		Extra: map[string]any{
			UpstreamBillingProbeEnabledExtraKey: true,
			UpstreamBillingProbeExtraKey:        map[string]any{"status": "ok"placeholder,
	placeholder,
placeholder)

placeholder
	require.Equal(t, true, updated.Extra[UpstreamBillingProbeEnabledExtraKey])
	require.NotContains(t, updated.Extra, UpstreamBillingProbeExtraKey)
placeholder

func TestUpdateAccountExplicitProbeDisableUsesDedicatedExtraUpdate(t *testing.T) {
	accountID := int64(113)
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{
		accountID: {
			ID:       accountID,
			Platform: PlatformOpenAI,
			Type:     AccountTypeAPIKey,
			Status:   StatusActive,
			Extra: map[string]any{
				UpstreamBillingProbeEnabledExtraKey: true,
				UpstreamBillingProbeExtraKey:        map[string]any{"status": "ok"placeholder,
		placeholder,
	placeholder,
placeholderplaceholder

	_, err := (&adminServiceImpl{accountRepo: repoplaceholder).UpdateAccount(context.Background(), accountID, &UpdateAccountInput{
		Extra: map[string]any{UpstreamBillingProbeEnabledExtraKey: falseplaceholder,
placeholder)

placeholder
	require.Len(t, repo.updates[accountID], 1)
	require.Equal(t, false, repo.updates[accountID][0][UpstreamBillingProbeEnabledExtraKey])
placeholder

func TestUpdateAccountExplicitUnchangedProbeEnabledStillUsesDedicatedExtraUpdate(t *testing.T) {
	accountID := int64(114)
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{
		accountID: {
			ID:       accountID,
			Platform: PlatformOpenAI,
			Type:     AccountTypeAPIKey,
			Status:   StatusActive,
			Extra:    map[string]any{UpstreamBillingProbeEnabledExtraKey: trueplaceholder,
	placeholder,
placeholderplaceholder

	_, err := (&adminServiceImpl{accountRepo: repoplaceholder).UpdateAccount(context.Background(), accountID, &UpdateAccountInput{
		Extra: map[string]any{UpstreamBillingProbeEnabledExtraKey: trueplaceholder,
placeholder)

placeholder
	require.Len(t, repo.updates[accountID], 1)
	require.Equal(t, true, repo.updates[accountID][0][UpstreamBillingProbeEnabledExtraKey])
placeholder

func TestUpdateAccountRejectsInvalidProbeEnabled(t *testing.T) {
	accountID := int64(112)
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{
		accountID: {
			ID:       accountID,
			Platform: PlatformOpenAI,
			Type:     AccountTypeAPIKey,
			Status:   StatusActive,
			Extra:    map[string]any{placeholder,
	placeholder,
placeholderplaceholder

	svc := &adminServiceImpl{accountRepo: repoplaceholder
	_, err := svc.UpdateAccount(context.Background(), accountID, &UpdateAccountInput{
		Extra: map[string]any{UpstreamBillingProbeEnabledExtraKey: "true"placeholder,
placeholder)

placeholder
placeholder

func TestBulkUpdateAccountsDropsManagedUpstreamBillingProbeState(t *testing.T) {
	repo := &upstreamBillingProbeAccountRepo{placeholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder
	input := &BulkUpdateAccountsInput{
		AccountIDs: []int64{1placeholder,
		Extra: map[string]any{
			"custom":                            "value",
			UpstreamBillingProbeEnabledExtraKey: true,
			UpstreamBillingProbeExtraKey:        map[string]any{"status": "ok"placeholder,
	placeholder,
placeholder

	result, err := svc.BulkUpdateAccounts(context.Background(), input)

placeholder
	require.Equal(t, 1, result.Success)
	require.Len(t, repo.bulkUpdates, 1)
	require.Equal(t, "value", repo.bulkUpdates[0].Extra["custom"])
	require.NotContains(t, repo.bulkUpdates[0].Extra, UpstreamBillingProbeEnabledExtraKey)
	require.NotContains(t, repo.bulkUpdates[0].Extra, UpstreamBillingProbeExtraKey)
placeholder

func TestBulkUpdateAccountsInvalidatesProbeSnapshotForIdentityCredentials(t *testing.T) {
	repo := &upstreamBillingProbeAccountRepo{placeholder
	input := &BulkUpdateAccountsInput{
		AccountIDs:  []int64{1placeholder,
placeholder"api_key": "sk-new"placeholder,
placeholder

	result, err := (&adminServiceImpl{accountRepo: repoplaceholder).BulkUpdateAccounts(context.Background(), input)

placeholder
	require.Equal(t, 1, result.Success)
	require.Len(t, repo.bulkUpdates, 1)
	require.Contains(t, repo.bulkUpdates[0].Extra, UpstreamBillingProbeExtraKey)
	require.Nil(t, repo.bulkUpdates[0].Extra[UpstreamBillingProbeExtraKey])
placeholder

func TestBulkUpdateAccountsInvalidatesProbeSnapshotForProxyUpdate(t *testing.T) {
	proxyID := int64(9)
	baseRepo := &upstreamBillingProbeAccountRepo{placeholder
	input := &BulkUpdateAccountsInput{
		AccountIDs: []int64{1placeholder,
		ProxyID:    &proxyID,
placeholder

	result, err := (&adminServiceImpl{accountRepo: &upstreamBillingProbeAdminRepo{baseRepoplaceholderplaceholder).BulkUpdateAccounts(context.Background(), input)

placeholder
	require.Equal(t, 1, result.Success)
	require.Len(t, baseRepo.bulkUpdates, 1)
	require.Contains(t, baseRepo.bulkUpdates[0].Extra, UpstreamBillingProbeExtraKey)
	require.Nil(t, baseRepo.bulkUpdates[0].Extra[UpstreamBillingProbeExtraKey])
placeholder

func TestBulkUpdateAccountsKeepsProbeSnapshotForUnrelatedCredentials(t *testing.T) {
	repo := &upstreamBillingProbeAccountRepo{placeholder
	input := &BulkUpdateAccountsInput{
		AccountIDs:  []int64{1placeholder,
placeholder"model_mapping": map[string]any{"gpt-old": "gpt-new"placeholderplaceholder,
placeholder

	_, err := (&adminServiceImpl{accountRepo: repoplaceholder).BulkUpdateAccounts(context.Background(), input)

placeholder
	require.Len(t, repo.bulkUpdates, 1)
	require.NotContains(t, repo.bulkUpdates[0].Extra, UpstreamBillingProbeExtraKey)
placeholder
