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

type accountBillingSettingsAdminRepo struct {
	*upstreamBillingProbeAccountRepo
	concurrentRate   *float64
	lastExplicitRate *float64
	updateCalls      int
placeholder

func (r *accountBillingSettingsAdminRepo) UpdateWithAccountBillingSettings(
	_ context.Context,
	account *Account,
	probeEnabled *bool,
	rateSyncEnabled *bool,
	rateMultiplier *float64,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	current := r.accounts[account.ID]
	if current == nil {
		return ErrAccountNotFound
placeholder
	updated := *account
	updated.Credentials = mergeMap(nil, account.Credentials)
	updated.Extra = mergeMap(nil, account.Extra)
	if updated.Extra == nil {
		updated.Extra = make(map[string]any)
placeholder
	if probeEnabled != nil {
		updated.Extra[UpstreamBillingProbeEnabledExtraKey] = *probeEnabled
placeholder
	if rateSyncEnabled != nil {
		updated.Extra[UpstreamBillingRateSyncEnabledExtraKey] = *rateSyncEnabled
placeholder
	switch {
	case rateMultiplier != nil:
		value := *rateMultiplier
		updated.RateMultiplier = &value
		r.lastExplicitRate = &value
	case r.concurrentRate != nil:
		value := *r.concurrentRate
		updated.RateMultiplier = &value
		r.lastExplicitRate = nil
	default:
		updated.RateMultiplier = cloneAccountValuePointer(current.RateMultiplier)
		r.lastExplicitRate = nil
placeholder
	r.accounts[account.ID] = &updated
	r.updateCalls++
	return nil
placeholder

func TestUpdateAccountRoutesRateIntentThroughAtomicBillingUpdater(t *testing.T) {
	accountID := int64(109)
	initialRate := 0.1
	concurrentRate := 0.2
	repo := &accountBillingSettingsAdminRepo{
		upstreamBillingProbeAccountRepo: &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{
			accountID: {
				ID:             accountID,
				Name:           "before",
				Platform:       PlatformOpenAI,
				Type:           AccountTypeAPIKey,
				Status:         StatusActive,
				RateMultiplier: &initialRate,
				Extra: map[string]any{
					UpstreamBillingProbeEnabledExtraKey:    true,
					UpstreamBillingRateSyncEnabledExtraKey: true,
			placeholder,
		placeholder,
placeholder
		concurrentRate: &concurrentRate,
placeholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder

	updated, err := svc.UpdateAccount(context.Background(), accountID, &UpdateAccountInput{Name: "after"placeholder)
placeholder
	require.Equal(t, 1, repo.updateCalls)
	require.Nil(t, repo.lastExplicitRate)
	require.Equal(t, concurrentRate, *updated.RateMultiplier)

	zero := 0.0
	updated, err = svc.UpdateAccount(context.Background(), accountID, &UpdateAccountInput{RateMultiplier: &zeroplaceholder)
placeholder
	require.Equal(t, 2, repo.updateCalls)
	require.NotNil(t, repo.lastExplicitRate)
	require.Zero(t, *repo.lastExplicitRate)
	require.Zero(t, *updated.RateMultiplier)
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
			UpstreamBillingProbeEnabledExtraKey:    true,
			UpstreamBillingRateSyncEnabledExtraKey: true,
			UpstreamBillingProbeExtraKey:           map[string]any{"status": "ok"placeholder,
	placeholder,
placeholder)

placeholder
	require.NotContains(t, created.Extra, UpstreamBillingProbeEnabledExtraKey)
	require.NotContains(t, created.Extra, UpstreamBillingRateSyncEnabledExtraKey)
	require.NotContains(t, created.Extra, UpstreamBillingProbeExtraKey)
placeholder

func TestCreateAccountAcceptsDedicatedUpstreamBillingProbeSetting(t *testing.T) {
	enabled := true
	repo := &upstreamBillingProbeAccountRepo{placeholder
	created, err := (&adminServiceImpl{accountRepo: repoplaceholder).CreateAccount(context.Background(), &CreateAccountInput{
		Name:                 "upstream",
		Platform:             PlatformOpenAI,
		Type:                 AccountTypeAPIKey,
		Credentials:          map[string]any{"api_key": "sk-test"placeholder,
		ProbeEnabled:         &enabled,
		SkipDefaultGroupBind: true,
placeholder)

placeholder
	require.Equal(t, true, created.Extra[UpstreamBillingProbeEnabledExtraKey])

	_, err = (&adminServiceImpl{accountRepo: repoplaceholder).CreateAccount(context.Background(), &CreateAccountInput{
		Name:                 "oauth",
		Platform:             PlatformOpenAI,
		Type:                 AccountTypeOAuth,
		Credentials:          map[string]any{"access_token": "token"placeholder,
		ProbeEnabled:         &enabled,
		SkipDefaultGroupBind: true,
placeholder)
	require.ErrorIs(t, err, ErrUpstreamBillingProbeAccountInvalid)
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
				UpstreamBillingProbeEnabledExtraKey:    true,
				UpstreamBillingRateSyncEnabledExtraKey: true,
				UpstreamBillingProbeExtraKey:           map[string]any{"status": "ok"placeholder,
		placeholder,
	placeholder,
placeholderplaceholder

	svc := &adminServiceImpl{accountRepo: repoplaceholder
	updated, err := svc.UpdateAccount(context.Background(), accountID, &UpdateAccountInput{
		Extra: map[string]any{"custom": "value"placeholder,
placeholder)

placeholder
	require.Equal(t, true, updated.Extra[UpstreamBillingProbeEnabledExtraKey])
	require.Equal(t, true, updated.Extra[UpstreamBillingRateSyncEnabledExtraKey])
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
						UpstreamBillingProbeEnabledExtraKey:    true,
						UpstreamBillingRateSyncEnabledExtraKey: true,
						UpstreamBillingProbeExtraKey:           map[string]any{"status": "ok"placeholder,
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
				require.NotContains(t, updated.Extra, UpstreamBillingRateSyncEnabledExtraKey)
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
			UpstreamBillingProbeEnabledExtraKey:    true,
			UpstreamBillingRateSyncEnabledExtraKey: true,
			UpstreamBillingProbeExtraKey:           map[string]any{"status": "ok"placeholder,
	placeholder,
placeholder)

placeholder
	require.Equal(t, true, updated.Extra[UpstreamBillingProbeEnabledExtraKey])
	require.NotContains(t, updated.Extra, UpstreamBillingRateSyncEnabledExtraKey)
	require.NotContains(t, updated.Extra, UpstreamBillingProbeExtraKey)
placeholder

func TestUpdateAccountRateSyncControlsProbeAndManualMode(t *testing.T) {
	accountID := int64(151)
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{
		accountID: {
			ID:       accountID,
	placeholder
			Type:     AccountTypeAPIKey,
			Status:   StatusActive,
			Extra:    map[string]any{placeholder,
	placeholder,
placeholderplaceholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder

	syncEnabled := true
	updated, err := svc.UpdateAccount(context.Background(), accountID, &UpdateAccountInput{
		RateSyncEnabled: &syncEnabled,
placeholder)
placeholder
	require.Equal(t, true, updated.Extra[UpstreamBillingProbeEnabledExtraKey])
	require.Equal(t, true, updated.Extra[UpstreamBillingRateSyncEnabledExtraKey])

	syncEnabled = false
	updated, err = svc.UpdateAccount(context.Background(), accountID, &UpdateAccountInput{
		RateSyncEnabled: &syncEnabled,
placeholder)
placeholder
	require.Equal(t, true, updated.Extra[UpstreamBillingProbeEnabledExtraKey])
	require.Equal(t, false, updated.Extra[UpstreamBillingRateSyncEnabledExtraKey])
placeholder

func TestUpdateAccountRejectsSyncWithExplicitlyDisabledProbe(t *testing.T) {
	accountID := int64(152)
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{
		accountID: {
			ID:       accountID,
			Platform: PlatformAnthropic,
			Type:     AccountTypeAPIKey,
			Status:   StatusActive,
	placeholder,
placeholderplaceholder
	probeEnabled := false
	syncEnabled := true

	_, err := (&adminServiceImpl{accountRepo: repoplaceholder).UpdateAccount(context.Background(), accountID, &UpdateAccountInput{
		ProbeEnabled:    &probeEnabled,
		RateSyncEnabled: &syncEnabled,
placeholder)

placeholder
	require.Empty(t, repo.updates[accountID])
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
	require.Equal(t, false, repo.updates[accountID][0][UpstreamBillingRateSyncEnabledExtraKey])
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

func TestUpdateAccountExtraDropsManagedBillingProbeFields(t *testing.T) {
	accountID := int64(153)
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{
		accountID: {ID: accountID, Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder,
placeholderplaceholder

	err := (&adminServiceImpl{accountRepo: repoplaceholder).UpdateAccountExtra(context.Background(), accountID, map[string]any{
		"custom":                               "value",
		UpstreamBillingProbeEnabledExtraKey:    true,
		UpstreamBillingRateSyncEnabledExtraKey: true,
		UpstreamBillingProbeExtraKey:           map[string]any{"status": "ok"placeholder,
placeholder)

placeholder
	require.Equal(t, "value", repo.accounts[accountID].Extra["custom"])
	require.NotContains(t, repo.accounts[accountID].Extra, UpstreamBillingProbeEnabledExtraKey)
	require.NotContains(t, repo.accounts[accountID].Extra, UpstreamBillingRateSyncEnabledExtraKey)
	require.NotContains(t, repo.accounts[accountID].Extra, UpstreamBillingProbeExtraKey)
placeholder

func TestBulkUpdateAccountsDropsManagedUpstreamBillingProbeState(t *testing.T) {
	repo := &upstreamBillingProbeAccountRepo{placeholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder
	input := &BulkUpdateAccountsInput{
		AccountIDs: []int64{1placeholder,
		Extra: map[string]any{
			"custom":                               "value",
			UpstreamBillingProbeEnabledExtraKey:    true,
			UpstreamBillingRateSyncEnabledExtraKey: true,
			UpstreamBillingProbeExtraKey:           map[string]any{"status": "ok"placeholder,
	placeholder,
placeholder

	result, err := svc.BulkUpdateAccounts(context.Background(), input)

placeholder
	require.Equal(t, 1, result.Success)
	require.Len(t, repo.bulkUpdates, 1)
	require.Equal(t, "value", repo.bulkUpdates[0].Extra["custom"])
	require.NotContains(t, repo.bulkUpdates[0].Extra, UpstreamBillingProbeEnabledExtraKey)
	require.NotContains(t, repo.bulkUpdates[0].Extra, UpstreamBillingRateSyncEnabledExtraKey)
	require.NotContains(t, repo.bulkUpdates[0].Extra, UpstreamBillingProbeExtraKey)
placeholder

func TestBulkUpdateAccountsAcceptsDedicatedUpstreamBillingProbeSetting(t *testing.T) {
	for _, enabled := range []bool{true, falseplaceholder {
		t.Run(map[bool]string{true: "enable", false: "disable"placeholder[enabled], func(t *testing.T) {
			repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{
				1: {ID: 1, Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder,
				2: {ID: 2, Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder,
		placeholderplaceholder

			result, err := (&adminServiceImpl{accountRepo: repoplaceholder).BulkUpdateAccounts(context.Background(), &BulkUpdateAccountsInput{
				AccountIDs:   []int64{1, 2placeholder,
				ProbeEnabled: &enabled,
		placeholder)

		placeholder
			require.Equal(t, 2, result.Success)
			require.Len(t, repo.bulkUpdates, 1)
			require.Equal(t, enabled, repo.bulkUpdates[0].Extra[UpstreamBillingProbeEnabledExtraKey])
			if !enabled {
				require.Equal(t, false, repo.bulkUpdates[0].Extra[UpstreamBillingRateSyncEnabledExtraKey])
		placeholder
			require.NotNil(t, repo.bulkUpdates[0].ProbeEnabled)
			require.Equal(t, enabled, *repo.bulkUpdates[0].ProbeEnabled)
	placeholder)
placeholder
placeholder

func TestBulkUpdateAccountsRejectsProbeSettingForIneligibleTargetBeforeWrite(t *testing.T) {
	for _, enabled := range []bool{true, falseplaceholder {
		t.Run(map[bool]string{true: "enable", false: "disable"placeholder[enabled], func(t *testing.T) {
			repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{
				1: {ID: 1, Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder,
				2: {ID: 2, Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder,
		placeholderplaceholder

			_, err := (&adminServiceImpl{accountRepo: repoplaceholder).BulkUpdateAccounts(context.Background(), &BulkUpdateAccountsInput{
				AccountIDs:   []int64{1, 2placeholder,
				ProbeEnabled: &enabled,
		placeholder)

			require.ErrorIs(t, err, ErrUpstreamBillingProbeAccountInvalid)
			require.Empty(t, repo.bulkUpdates)
	placeholder)
placeholder
placeholder

func TestBulkUpdateAccountsRejectsProbeSettingWhenTargetIsMissing(t *testing.T) {
	enabled := true
	repo := &upstreamBillingProbeAccountRepo{accounts: map[int64]*Account{
		1: {ID: 1, Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder,
placeholderplaceholder

	_, err := (&adminServiceImpl{accountRepo: repoplaceholder).BulkUpdateAccounts(context.Background(), &BulkUpdateAccountsInput{
		AccountIDs:   []int64{1, 2placeholder,
		ProbeEnabled: &enabled,
placeholder)

	require.ErrorIs(t, err, ErrAccountNotFound)
	require.Empty(t, repo.bulkUpdates)
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
