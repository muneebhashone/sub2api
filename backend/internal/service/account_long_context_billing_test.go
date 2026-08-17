//go:build unit

package service

import (
	"context"
	"net/http"
	"testing"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestAccountIsOpenAILongContextBillingEnabled(t *testing.T) {
	tests := []struct {
		name    string
		account *Account
		want    bool
placeholder{
		{name: "nil account is disabled", account: nil, want: falseplaceholder,
		{name: "non OpenAI account is disabled", account: &Account{Platform: PlatformGrokplaceholder, want: falseplaceholder,
		{name: "missing extra defaults disabled", account: &Account{Platform: PlatformOpenAIplaceholder, want: falseplaceholder,
		{name: "missing key defaults disabled", account: &Account{Platform: PlatformOpenAI, Extra: map[string]any{placeholderplaceholder, want: falseplaceholder,
		{name: "explicit true is enabled", account: &Account{Platform: PlatformOpenAI, Extra: map[string]any{"openai_long_context_billing_enabled": trueplaceholderplaceholder, want: trueplaceholder,
		{name: "explicit false is disabled", account: &Account{Platform: PlatformOpenAI, Extra: map[string]any{"openai_long_context_billing_enabled": falseplaceholderplaceholder, want: falseplaceholder,
		{name: "malformed value is disabled", account: &Account{Platform: PlatformOpenAI, Extra: map[string]any{"openai_long_context_billing_enabled": "false"placeholderplaceholder, want: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.account.IsOpenAILongContextBillingEnabled())
	placeholder)
placeholder
placeholder

func TestNormalizeOpenAILongContextBillingExtra(t *testing.T) {
	t.Run("OpenAI missing key persists disabled default", func(t *testing.T) {
		extra, err := normalizeOpenAILongContextBillingExtra(PlatformOpenAI, nil)

	placeholder
		require.Equal(t, false, extra["openai_long_context_billing_enabled"])
placeholder)

	t.Run("OpenAI explicit false is preserved", func(t *testing.T) {
		extra, err := normalizeOpenAILongContextBillingExtra(PlatformOpenAI, map[string]any{"openai_long_context_billing_enabled": falseplaceholder)

	placeholder
		require.Equal(t, false, extra["openai_long_context_billing_enabled"])
placeholder)

	t.Run("OpenAI malformed value is rejected", func(t *testing.T) {
		_, err := normalizeOpenAILongContextBillingExtra(PlatformOpenAI, map[string]any{"openai_long_context_billing_enabled": "false"placeholder)

	placeholder
		require.Equal(t, http.StatusBadRequest, infraerrors.Code(err))
placeholder)

	t.Run("non OpenAI extra is unchanged", func(t *testing.T) {
		extra, err := normalizeOpenAILongContextBillingExtra(PlatformGrok, nil)

	placeholder
		require.Nil(t, extra)
placeholder)

	t.Run("non OpenAI malformed value is ignored", func(t *testing.T) {
		extra := map[string]any{openAILongContextBillingEnabledKey: "provider-owned"placeholder
		normalized, err := normalizeOpenAILongContextBillingExtra(PlatformAnthropic, extra)

	placeholder
		require.Equal(t, extra, normalized)
placeholder)
placeholder

type longContextBillingRepoStub struct {
	accountRepoStub
	account          *Account
	accounts         []*Account
	createdAccount   *Account
	updateExtraCalls int
	bulkUpdateCalls  int
placeholder

func (r *longContextBillingRepoStub) Create(_ context.Context, account *Account) error {
	account.ID = 1
	r.account = account
	r.createdAccount = account
	return nil
placeholder

func (r *longContextBillingRepoStub) GetByID(_ context.Context, _ int64) (*Account, error) {
	return r.account, nil
placeholder

func (r *longContextBillingRepoStub) GetByIDs(_ context.Context, _ []int64) ([]*Account, error) {
	if r.accounts != nil {
		return r.accounts, nil
placeholder
	if r.account == nil {
		return nil, nil
placeholder
	return []*Account{r.accountplaceholder, nil
placeholder

func (r *longContextBillingRepoStub) Update(_ context.Context, account *Account) error {
	r.account = account
	return nil
placeholder

func (r *longContextBillingRepoStub) UpdateExtra(_ context.Context, _ int64, _ map[string]any) error {
	r.updateExtraCalls++
	return nil
placeholder

func (r *longContextBillingRepoStub) BulkUpdate(_ context.Context, _ []int64, _ AccountBulkUpdate) (int64, error) {
	r.bulkUpdateCalls++
	return 1, nil
placeholder

func TestAdminServiceCreateAccountDefaultsOpenAILongContextBillingDisabled(t *testing.T) {
	repo := &longContextBillingRepoStub{placeholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder

	account, err := svc.CreateAccount(context.Background(), &CreateAccountInput{
		Name:                 "openai-account",
		Platform:             PlatformOpenAI,
		Type:                 AccountTypeAPIKey,
		Credentials:          map[string]any{"api_key": "test"placeholder,
		SkipDefaultGroupBind: true,
placeholder)

placeholder
	require.Same(t, account, repo.createdAccount)
	require.Equal(t, false, account.Extra[openAILongContextBillingEnabledKey])
placeholder

func TestAdminServiceCreateAccountRejectsMalformedOpenAILongContextBillingValue(t *testing.T) {
	repo := &longContextBillingRepoStub{placeholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder

	account, err := svc.CreateAccount(context.Background(), &CreateAccountInput{
		Platform: PlatformOpenAI,
		Extra:    map[string]any{openAILongContextBillingEnabledKey: "false"placeholder,
placeholder)

	require.Nil(t, account)
	require.Equal(t, http.StatusBadRequest, infraerrors.Code(err))
	require.Nil(t, repo.createdAccount)
placeholder

func TestAdminServiceUpdateAccountPreservesOpenAILongContextBillingOptOutWhenOmitted(t *testing.T) {
	repo := &longContextBillingRepoStub{account: &Account{
		ID:       1,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
		Extra:    map[string]any{openAILongContextBillingEnabledKey: falseplaceholder,
placeholderplaceholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder

	account, err := svc.UpdateAccount(context.Background(), 1, &UpdateAccountInput{Extra: map[string]any{placeholderplaceholder)

placeholder
	require.Equal(t, false, account.Extra[openAILongContextBillingEnabledKey])
placeholder

func TestAdminServiceUpdateAccountAllowsExplicitCodexImportOptIn(t *testing.T) {
	repo := &longContextBillingRepoStub{account: &Account{
		ID:          1,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
placeholder"access_token": "old-token"placeholder,
		Extra: map[string]any{
			openAILongContextBillingEnabledKey: false,
			"import_source":                    "codex_session",
	placeholder,
placeholderplaceholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder

	account, err := svc.UpdateAccount(context.Background(), 1, &UpdateAccountInput{
placeholder"access_token": "new-token"placeholder,
		Extra: map[string]any{
			openAILongContextBillingEnabledKey: true,
			"import_source":                    "codex_session",
	placeholder,
placeholder)

placeholder
	require.Equal(t, true, account.Extra[openAILongContextBillingEnabledKey])
placeholder

func TestAdminServiceUpdateAccountAllowsExplicitOptInOutsideCodexImport(t *testing.T) {
	repo := &longContextBillingRepoStub{account: &Account{
		ID:       1,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
		Extra: map[string]any{
			openAILongContextBillingEnabledKey: false,
			"import_source":                    "codex_session",
	placeholder,
placeholderplaceholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder

	account, err := svc.UpdateAccount(context.Background(), 1, &UpdateAccountInput{Extra: map[string]any{
		openAILongContextBillingEnabledKey: true,
		"import_source":                    "codex_session",
placeholderplaceholder)

placeholder
	require.Equal(t, true, account.Extra[openAILongContextBillingEnabledKey])
placeholder

func TestAdminServiceUpdateAccountRejectsMalformedOpenAILongContextBillingValue(t *testing.T) {
	repo := &longContextBillingRepoStub{account: &Account{ID: 1, Platform: PlatformOpenAIplaceholderplaceholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder

	account, err := svc.UpdateAccount(context.Background(), 1, &UpdateAccountInput{Extra: map[string]any{
		openAILongContextBillingEnabledKey: 1,
placeholderplaceholder)

	require.Nil(t, account)
	require.Equal(t, http.StatusBadRequest, infraerrors.Code(err))
placeholder

func TestAdminServiceUpdateAccountExtraRejectsMalformedOpenAILongContextBillingValue(t *testing.T) {
	repo := &longContextBillingRepoStub{account: &Account{ID: 1, Platform: PlatformOpenAIplaceholderplaceholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder

	err := svc.UpdateAccountExtra(context.Background(), 1, map[string]any{
		openAILongContextBillingEnabledKey: "true",
placeholder)

	require.Equal(t, http.StatusBadRequest, infraerrors.Code(err))
	require.Zero(t, repo.updateExtraCalls)
placeholder

func TestAdminServiceUpdateAccountExtraAllowsProviderOwnedValueForNonOpenAIAccount(t *testing.T) {
	repo := &longContextBillingRepoStub{account: &Account{ID: 1, Platform: PlatformAnthropicplaceholderplaceholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder

	err := svc.UpdateAccountExtra(context.Background(), 1, map[string]any{
		openAILongContextBillingEnabledKey: "provider-owned",
placeholder)

placeholder
	require.Equal(t, 1, repo.updateExtraCalls)
placeholder

func TestAdminServiceBulkUpdateAccountsRejectsMalformedOpenAILongContextBillingValue(t *testing.T) {
	repo := &longContextBillingRepoStub{account: &Account{ID: 1, Platform: PlatformOpenAIplaceholderplaceholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder

	result, err := svc.BulkUpdateAccounts(context.Background(), &BulkUpdateAccountsInput{
		AccountIDs: []int64{1placeholder,
		Extra:      map[string]any{openAILongContextBillingEnabledKey: []bool{trueplaceholderplaceholder,
placeholder)

	require.Nil(t, result)
	require.Equal(t, http.StatusBadRequest, infraerrors.Code(err))
	require.Zero(t, repo.bulkUpdateCalls)
placeholder

func TestAdminServiceBulkUpdateAccountsRejectsOpenAILongContextKeyForNonOpenAIAccounts(t *testing.T) {
	repo := &longContextBillingRepoStub{account: &Account{ID: 1, Platform: PlatformGrokplaceholderplaceholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder

	result, err := svc.BulkUpdateAccounts(context.Background(), &BulkUpdateAccountsInput{
		AccountIDs: []int64{1placeholder,
		Extra:      map[string]any{openAILongContextBillingEnabledKey: trueplaceholder,
placeholder)

	require.Nil(t, result)
	var appErr *infraerrors.ApplicationError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, "OPENAI_BULK_TARGET_INVALID", appErr.Reason)
	require.Zero(t, repo.bulkUpdateCalls)
placeholder

func TestAdminServiceBulkUpdateAccountsRejectsMalformedValueForMixedTargetsIncludingOpenAI(t *testing.T) {
	repo := &longContextBillingRepoStub{accounts: []*Account{
		{ID: 1, Platform: PlatformGrokplaceholder,
		{ID: 2, Platform: PlatformOpenAIplaceholder,
placeholderplaceholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder

	result, err := svc.BulkUpdateAccounts(context.Background(), &BulkUpdateAccountsInput{
		AccountIDs: []int64{1, 2placeholder,
		Extra:      map[string]any{openAILongContextBillingEnabledKey: "malformed"placeholder,
placeholder)

	require.Nil(t, result)
	require.Equal(t, http.StatusBadRequest, infraerrors.Code(err))
	require.Zero(t, repo.bulkUpdateCalls)
placeholder
