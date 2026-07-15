package service

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type grokReconcileRepo struct {
	AccountRepository

	mu                      sync.Mutex
	accounts                []Account
	requests                []OAuthRefreshPageOptions
	setErrorIDs             []int64
	updatedCredIDs          []int64
	setErrorMessage         []string
	getByIDOverrides        map[int64]Account
	pageOverride            *OAuthRefreshCandidatePage
	reauthorizeOnCAS        bool
	reauthorizeOnRefreshCAS bool
	conditionalCalls        int
placeholder

func (r *grokReconcileRepo) GetByID(_ context.Context, id int64) (*Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if override, ok := r.getByIDOverrides[id]; ok {
		account := override
		return &account, nil
placeholder
	for i := range r.accounts {
		if r.accounts[i].ID == id {
			account := r.accounts[i]
			return &account, nil
	placeholder
placeholder
	return nil, ErrAccountNotFound
placeholder

func (r *grokReconcileRepo) ListOAuthRefreshCandidatePage(_ context.Context, options OAuthRefreshPageOptions) (*OAuthRefreshCandidatePage, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.requests = append(r.requests, options)
	if r.pageOverride != nil {
		page := *r.pageOverride
		page.Accounts = append([]Account(nil), r.pageOverride.Accounts...)
		return &page, nil
placeholder
	accounts := append([]Account(nil), r.accounts...)
	sort.Slice(accounts, func(i, j int) bool { return accounts[i].ID < accounts[j].ID placeholder)
	page := make([]Account, 0, options.Limit)
	for _, account := range accounts {
		if account.ID <= options.AfterID {
			continue
	placeholder
		platformAllowed := false
		for _, platform := range options.Platforms {
			if account.Platform == platform {
				platformAllowed = true
				break
		placeholder
	placeholder
		if !platformAllowed || options.ActiveOnly && account.Status != StatusActive {
			continue
	placeholder
		if options.IncludeSetupToken {
			if account.Type != AccountTypeOAuth && account.Type != AccountTypeSetupToken {
				continue
		placeholder
	placeholder else if account.Type != AccountTypeOAuth {
			continue
	placeholder
		if options.RequireRefreshToken && strings.TrimSpace(account.GetGrokRefreshToken()) == "" {
			continue
	placeholder
		page = append(page, account)
		if len(page) == options.Limit {
			break
	placeholder
placeholder
	result := &OAuthRefreshCandidatePage{Accounts: page, HasMore: len(page) == options.Limitplaceholder
	if len(page) > 0 {
		result.NextAfterID = page[len(page)-1].ID
placeholder
	return result, nil
placeholder

func (r *grokReconcileRepo) UpdateCredentials(_ context.Context, id int64, credentials map[string]any) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.updatedCredIDs = append(r.updatedCredIDs, id)
	for i := range r.accounts {
		if r.accounts[i].ID == id {
			r.accounts[i].Credentials = MergeCredentials(r.accounts[i].Credentials, credentials)
	placeholder
placeholder
	return nil
placeholder

func (r *grokReconcileRepo) SetError(_ context.Context, id int64, message string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.setErrorIDs = append(r.setErrorIDs, id)
	r.setErrorMessage = append(r.setErrorMessage, message)
	for i := range r.accounts {
		if r.accounts[i].ID == id {
			r.accounts[i].Status = StatusError
			r.accounts[i].Schedulable = false
			r.accounts[i].ErrorMessage = message
	placeholder
placeholder
	return nil
placeholder

func (r *grokReconcileRepo) SetGrokOAuthErrorIfCredentialsUnchanged(
	_ context.Context,
	id int64,
	expectedCredentials map[string]any,
	message string,
) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.conditionalCalls++
	for i := range r.accounts {
		account := &r.accounts[i]
		if account.ID != id {
			continue
	placeholder
		if r.reauthorizeOnCAS {
			r.reauthorizeOnCAS = false
			account.Credentials = map[string]any{
				"access_token":   "fresh-access",
				"refresh_token":  "fresh-refresh",
				"expires_at":     time.Now().UTC().Add(4 * time.Hour).Format(time.RFC3339),
				"_token_version": int64(2),
		placeholder
	placeholder
		if account.Platform != PlatformGrok || account.Type != AccountTypeOAuth || account.Status != StatusActive ||
			strings.TrimSpace(account.GetGrokRefreshToken()) != "" || !reflect.DeepEqual(account.Credentials, expectedCredentials) {
			return false, nil
	placeholder
		r.setErrorIDs = append(r.setErrorIDs, id)
		r.setErrorMessage = append(r.setErrorMessage, message)
		account.Status = StatusError
		account.Schedulable = false
		account.ErrorMessage = message
		return true, nil
placeholder
	return false, nil
placeholder

func (r *grokReconcileRepo) SetGrokOAuthRefreshErrorIfCredentialsUnchanged(
	_ context.Context,
	id int64,
	expectedCredentials map[string]any,
	expectedProxyID *int64,
	message string,
) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.accounts {
		account := &r.accounts[i]
		if account.ID != id {
			continue
	placeholder
		if r.reauthorizeOnRefreshCAS {
			r.reauthorizeOnRefreshCAS = false
			account.Credentials = map[string]any{
				"access_token":   "fresh-access",
				"refresh_token":  "fresh-refresh",
				"expires_at":     time.Now().UTC().Add(4 * time.Hour).Format(time.RFC3339),
				"_token_version": int64(3),
		placeholder
	placeholder
		if account.Platform != PlatformGrok || account.Type != AccountTypeOAuth || account.Status != StatusActive ||
			!reflect.DeepEqual(account.ProxyID, expectedProxyID) ||
			!reflect.DeepEqual(account.Credentials, expectedCredentials) {
			return false, nil
	placeholder
		r.setErrorIDs = append(r.setErrorIDs, id)
		r.setErrorMessage = append(r.setErrorMessage, message)
		account.Status = StatusError
		account.Schedulable = false
		account.ErrorMessage = message
		return true, nil
placeholder
	return false, nil
placeholder

func (r *grokReconcileRepo) SetGrokOAuthRefreshTempUnschedulableIfCredentialsUnchanged(
	_ context.Context,
	id int64,
	expectedCredentials map[string]any,
	expectedProxyID *int64,
	until time.Time,
	reason string,
) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.accounts {
		account := &r.accounts[i]
		if account.ID != id {
			continue
	placeholder
		if account.Platform != PlatformGrok || account.Type != AccountTypeOAuth || account.Status != StatusActive ||
			!reflect.DeepEqual(account.ProxyID, expectedProxyID) ||
			!reflect.DeepEqual(account.Credentials, expectedCredentials) {
			return false, nil
	placeholder
		account.TempUnschedulableUntil = &until
		account.TempUnschedulableReason = reason
		return true, nil
placeholder
	return false, nil
placeholder

func (r *grokReconcileRepo) snapshot() ([]OAuthRefreshPageOptions, []int64, []int64, []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]OAuthRefreshPageOptions(nil), r.requests...), append([]int64(nil), r.setErrorIDs...), append([]int64(nil), r.updatedCredIDs...), append([]string(nil), r.setErrorMessage...)
placeholder

type reconcileInvalidator struct {
	mu  sync.Mutex
	ids []int64
	err error
placeholder

type reconcileRuntimeBlocker struct {
	mu      sync.Mutex
	blocked []int64
	cleared []int64
placeholder

func (b *reconcileRuntimeBlocker) BlockAccountScheduling(account *Account, _ time.Time, _ string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if account != nil {
		b.blocked = append(b.blocked, account.ID)
placeholder
placeholder

func (b *reconcileRuntimeBlocker) ClearAccountSchedulingBlock(accountID int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.cleared = append(b.cleared, accountID)
placeholder

func (b *reconcileRuntimeBlocker) snapshot() (blocked, cleared []int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]int64(nil), b.blocked...), append([]int64(nil), b.cleared...)
placeholder

func (i *reconcileInvalidator) InvalidateToken(_ context.Context, account *Account) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.ids = append(i.ids, account.ID)
	return i.err
placeholder

func (i *reconcileInvalidator) count() int {
	i.mu.Lock()
	defer i.mu.Unlock()
	return len(i.ids)
placeholder

func newGrokReconcileService(repo *grokReconcileRepo, refresher *poolHealthRefresher, invalidator TokenCacheInvalidator) *TokenRefreshService {
	return &TokenRefreshService{
		accountRepo:      repo,
		candidatePager:   repo,
		cacheInvalidator: invalidator,
		registrations: []tokenRefreshRegistration{{
			platform:  PlatformGrok,
			refresher: refresher,
			executor:  refresher,
placeholder
		refreshPolicy: DefaultBackgroundRefreshPolicy(),
		cfg: &config.TokenRefreshConfig{
			MaxRetries:               1,
			CandidatePageSize:        50,
			ProviderConcurrency:      2,
			ProviderQPS:              100,
			ProviderFailureThreshold: 3,
			AttemptTimeoutSeconds:    1,
	placeholder,
placeholder
placeholder

func grokReconcileFixtures() []Account {
	now := time.Now().UTC()
	return []Account{
		{
			ID:          1,
			Platform:    PlatformGrok,
			Type:        AccountTypeOAuth,
			Status:      StatusActive,
			Schedulable: true,
	placeholder"access_token": "access-secret"placeholder,
	placeholder,
		{
			ID:          2,
			Platform:    PlatformGrok,
			Type:        AccountTypeOAuth,
			Status:      StatusActive,
			Schedulable: true,
	placeholder"refresh_token": "refresh-secret", "expires_at": now.Add(10 * time.Minute).Format(time.RFC3339)placeholder,
	placeholder,
		{
			ID:          3,
			Platform:    PlatformGrok,
			Type:        AccountTypeOAuth,
			Status:      StatusActive,
			Schedulable: true,
	placeholder"access_token": "access-secret", "refresh_token": "refresh-secret", "expires_at": now.Add(30 * time.Minute).Format(time.RFC3339)placeholder,
	placeholder,
		{
			ID:          4,
			Platform:    PlatformGrok,
			Type:        AccountTypeOAuth,
			Status:      StatusActive,
			Schedulable: true,
	placeholder"access_token": "access-secret", "refresh_token": "refresh-secret", "expires_at": now.Add(4 * time.Hour).Format(time.RFC3339)placeholder,
	placeholder,
		{
			ID:          5,
			Platform:    PlatformGrok,
			Type:        AccountTypeAPIKey,
			Status:      StatusActive,
			Schedulable: true,
	placeholder"api_key": "api-key-secret"placeholder,
	placeholder,
placeholder
placeholder

func TestTokenRefreshService_ReconcileGrokOAuthDefaultsToDryRunAndSanitizedPlan(t *testing.T) {
	repo := &grokReconcileRepo{accounts: grokReconcileFixtures()placeholder
	refresher := &poolHealthRefresher{placeholder
	svc := newGrokReconcileService(repo, refresher, nil)

	result, err := svc.ReconcileGrokOAuth(context.Background(), GrokOAuthReconcileInput{placeholder)

placeholder
	require.True(t, result.DryRun)
	require.Equal(t, 4, result.Scanned, "Grok API-key rows must not enter the OAuth reconciliation page")
	require.Equal(t, 3, result.Actionable)
	require.Equal(t, 1, result.WouldBlock)
	require.Equal(t, 2, result.WouldRefresh)
	require.Zero(t, result.Blocked)
	require.Zero(t, result.Refreshed)
	require.Zero(t, refresher.calls.Load())
	_, setErrorIDs, updatedIDs, _ := repo.snapshot()
	require.Empty(t, setErrorIDs)
	require.Empty(t, updatedIDs)

	payload, err := json.Marshal(result)
placeholder
	text := string(payload)
	require.NotContains(t, text, "access-secret")
	require.NotContains(t, text, "refresh-secret")
	require.NotContains(t, text, "api-key-secret")
	require.NotContains(t, text, `"credentials":`)
placeholder

func TestGrokTokenRefresher_NeedsRefreshWhenAccessTokenMissingDespiteFarFutureExpiry(t *testing.T) {
	refresher := NewGrokTokenRefresher(nil)
	account := grokPoolAccount(99)
	delete(account.Credentials, "access_token")
	account.Credentials["expires_at"] = time.Now().UTC().Add(12 * time.Hour).Format(time.RFC3339)

	require.True(t, refresher.NeedsRefresh(&account, time.Hour))
placeholder

func TestTokenRefreshService_ReconcileGrokOAuthApplyIsIdempotent(t *testing.T) {
	repo := &grokReconcileRepo{accounts: grokReconcileFixtures()placeholder
	invalidator := &reconcileInvalidator{placeholder
	refresher := &poolHealthRefresher{newCredentials: map[string]any{
		"access_token":  "rotated-access-secret",
		"refresh_token": "rotated-refresh-secret",
		"expires_at":    time.Now().UTC().Add(4 * time.Hour).Format(time.RFC3339),
placeholderplaceholder
	svc := newGrokReconcileService(repo, refresher, invalidator)

	first, err := svc.ReconcileGrokOAuth(context.Background(), GrokOAuthReconcileInput{Apply: true, Limit: 50placeholder)
placeholder
	require.False(t, first.DryRun)
	require.Equal(t, 1, first.Blocked)
	require.Equal(t, 2, first.Refreshed)
	require.Zero(t, first.Failed)
	requests, setErrorIDs, updatedIDs, messages := repo.snapshot()
	require.Equal(t, []int64{1placeholder, setErrorIDs)
	sort.Slice(updatedIDs, func(i, j int) bool { return updatedIDs[i] < updatedIDs[j] placeholder)
	require.Equal(t, []int64{2, 3placeholder, updatedIDs)
	require.Len(t, messages, 1)
	require.NotContains(t, messages[0], "secret")
	require.False(t, requests[0].RequireRefreshToken, "structurally invalid rows must remain discoverable")
	require.False(t, requests[0].IncludeSetupToken)
	require.Equal(t, 3, invalidator.count(), "block and refresh actions must invalidate token cache state")

	second, err := svc.ReconcileGrokOAuth(context.Background(), GrokOAuthReconcileInput{Apply: true, Limit: 50placeholder)
placeholder
	require.Zero(t, second.Actionable)
	require.Equal(t, int64(2), refresher.calls.Load(), "already refreshed rows must not be refreshed again")
	_, setErrorIDs, updatedIDs, _ = repo.snapshot()
	require.Equal(t, []int64{1placeholder, setErrorIDs, "already blocked invalid rows must not transition twice")
	require.Len(t, updatedIDs, 2)
placeholder

func TestTokenRefreshService_ReconcileGrokOAuthCursorResumesWithoutDuplicates(t *testing.T) {
	fixtures := grokReconcileFixtures()[:3]
	repo := &grokReconcileRepo{accounts: fixturesplaceholder
	svc := newGrokReconcileService(repo, &poolHealthRefresher{placeholder, nil)

	first, err := svc.ReconcileGrokOAuth(context.Background(), GrokOAuthReconcileInput{Limit: 2placeholder)
placeholder
	require.True(t, first.HasMore)
	require.Equal(t, int64(2), first.NextAfterID)
	require.Len(t, first.Items, 2)

	second, err := svc.ReconcileGrokOAuth(context.Background(), GrokOAuthReconcileInput{AfterID: first.NextAfterID, Limit: 2placeholder)
placeholder
	require.False(t, second.HasMore)
	require.Zero(t, second.NextAfterID)
	require.Len(t, second.Items, 1)
	require.NotEqual(t, first.Items[0].AccountID, second.Items[0].AccountID)
	require.NotEqual(t, first.Items[1].AccountID, second.Items[0].AccountID)
placeholder

func TestTokenRefreshService_ReconcileGrokOAuthCursorUsesRawPageAfterHydrationGap(t *testing.T) {
	account := grokReconcileFixtures()[0]
	repo := &grokReconcileRepo{pageOverride: &OAuthRefreshCandidatePage{
		Accounts:    []Account{accountplaceholder,
		NextAfterID: account.ID + 1,
		HasMore:     true,
placeholderplaceholder
	svc := newGrokReconcileService(repo, &poolHealthRefresher{placeholder, nil)

	result, err := svc.ReconcileGrokOAuth(context.Background(), GrokOAuthReconcileInput{Limit: 2placeholder)

placeholder
	require.True(t, result.HasMore)
	require.Equal(t, account.ID+1, result.NextAfterID,
		"cursor must advance past a raw selected ID that disappeared during hydration")
placeholder

func TestTokenRefreshService_ReconcileGrokOAuthRejectsConflictingApplyMode(t *testing.T) {
	svc := newGrokReconcileService(&grokReconcileRepo{placeholder, &poolHealthRefresher{placeholder, nil)

	_, err := svc.ReconcileGrokOAuth(context.Background(), GrokOAuthReconcileInput{Apply: true, DryRun: trueplaceholder)

	require.ErrorIs(t, err, ErrGrokOAuthReconcileMode)
placeholder

func TestTokenRefreshService_ReconcileGrokOAuthSkipsStaleBlockAfterConcurrentReauthorization(t *testing.T) {
	stale := grokReconcileFixtures()[0]
	latest := stale
	latest.Credentials = map[string]any{
		"access_token":  "fresh-access",
		"refresh_token": "fresh-refresh",
		"expires_at":    time.Now().UTC().Add(4 * time.Hour).Format(time.RFC3339),
placeholder
	repo := &grokReconcileRepo{
		accounts:         []Account{staleplaceholder,
		getByIDOverrides: map[int64]Account{stale.ID: latestplaceholder,
placeholder
	svc := newGrokReconcileService(repo, &poolHealthRefresher{placeholder, &reconcileInvalidator{placeholder)

	result, err := svc.ReconcileGrokOAuth(context.Background(), GrokOAuthReconcileInput{Apply: true, Limit: 50placeholder)

placeholder
	require.Zero(t, result.Blocked)
	require.Equal(t, 1, result.Skipped)
	require.Equal(t, GrokOAuthReconcileOutcomeSkipped, result.Items[0].Outcome)
	_, setErrorIDs, _, _ := repo.snapshot()
	require.Empty(t, setErrorIDs, "a concurrently reauthorized account must not be disabled from stale page state")
placeholder

func TestTokenRefreshService_ReconcileGrokOAuthDoesNotRuntimeBlockWhenReauthorizationWinsConditionalMutation(t *testing.T) {
	account := grokReconcileFixtures()[0]
	account.Credentials["_token_version"] = int64(1)
	repo := &grokReconcileRepo{
		accounts:         []Account{accountplaceholder,
		reauthorizeOnCAS: true,
placeholder
	invalidator := &reconcileInvalidator{placeholder
	blocker := &reconcileRuntimeBlocker{placeholder
	svc := newGrokReconcileService(repo, &poolHealthRefresher{placeholder, invalidator)
	svc.SetAccountRuntimeBlocker(blocker)

	result, err := svc.ReconcileGrokOAuth(context.Background(), GrokOAuthReconcileInput{Apply: true, Limit: 50placeholder)

placeholder
	require.Zero(t, result.Blocked)
	require.Equal(t, 1, result.Skipped)
	require.Equal(t, GrokOAuthReconcileOutcomeSkipped, result.Items[0].Outcome)
	require.Zero(t, invalidator.count(), "a lost compare-and-set race must not invalidate fresh credentials")
	_, setErrorIDs, _, _ := repo.snapshot()
	require.Empty(t, setErrorIDs)
	require.Equal(t, 1, repo.conditionalCalls)
	blocked, cleared := blocker.snapshot()
	require.Empty(t, blocked, "a lost compare-and-set race must never install a runtime block")
	require.Empty(t, cleared, "reconciliation must not clear a block it does not own")

	latest, getErr := repo.GetByID(context.Background(), account.ID)
	require.NoError(t, getErr)
	require.Equal(t, StatusActive, latest.Status)
	require.True(t, latest.Schedulable)
	require.Equal(t, "fresh-refresh", latest.GetGrokRefreshToken())
placeholder

func TestTokenRefreshService_ReconcileGrokOAuthReportsPermanentRefreshMutationAsBlocked(t *testing.T) {
	account := grokReconcileFixtures()[2]
	repo := &grokReconcileRepo{accounts: []Account{accountplaceholderplaceholder
	refresher := &poolHealthRefresher{err: errors.New(`GROK_OAUTH_ENTITLEMENT_DENIED: subscription required`)placeholder
	svc := newGrokReconcileService(repo, refresher, &reconcileInvalidator{placeholder)

	result, err := svc.ReconcileGrokOAuth(context.Background(), GrokOAuthReconcileInput{Apply: true, Limit: 50placeholder)

placeholder
	require.Equal(t, 1, result.Blocked)
	require.Zero(t, result.Failed)
	require.Zero(t, result.Partial)
	require.Equal(t, GrokOAuthReconcileActionBlock, result.Items[0].Action)
	require.Equal(t, GrokOAuthReconcileReasonCredentialRejected, result.Items[0].Reason)
	require.Equal(t, GrokOAuthReconcileOutcomeApplied, result.Items[0].Outcome)
	_, setErrorIDs, _, _ := repo.snapshot()
	require.Equal(t, []int64{account.IDplaceholder, setErrorIDs)
placeholder

func TestTokenRefreshService_ReconcileGrokOAuthReportsConcurrentRefreshReauthorizationAsSkipped(t *testing.T) {
	account := grokReconcileFixtures()[2]
	repo := &grokReconcileRepo{
		accounts:                []Account{accountplaceholder,
		reauthorizeOnRefreshCAS: true,
placeholder
	invalidator := &reconcileInvalidator{placeholder
	refresher := &poolHealthRefresher{err: errors.New("invalid_grant: revoked")placeholder
	svc := newGrokReconcileService(repo, refresher, invalidator)

	result, err := svc.ReconcileGrokOAuth(context.Background(), GrokOAuthReconcileInput{Apply: true, Limit: 50placeholder)

placeholder
	require.Equal(t, 1, result.Skipped)
	require.Zero(t, result.Failed)
	require.Zero(t, result.Blocked)
	require.Equal(t, GrokOAuthReconcileOutcomeSkipped, result.Items[0].Outcome)
	require.Zero(t, invalidator.count())
	_, setErrorIDs, _, _ := repo.snapshot()
	require.Empty(t, setErrorIDs)
	latest, getErr := repo.GetByID(context.Background(), account.ID)
	require.NoError(t, getErr)
	require.Equal(t, StatusActive, latest.Status)
	require.Equal(t, "fresh-refresh", latest.GetGrokRefreshToken())
placeholder

func TestTokenRefreshService_ReconcileGrokOAuthReportsInvalidationFailureAsPartial(t *testing.T) {
	account := grokReconcileFixtures()[0]
	repo := &grokReconcileRepo{accounts: []Account{accountplaceholderplaceholder
	svc := newGrokReconcileService(repo, &poolHealthRefresher{placeholder, &reconcileInvalidator{err: errors.New("cache unavailable")placeholder)

	result, err := svc.ReconcileGrokOAuth(context.Background(), GrokOAuthReconcileInput{Apply: true, Limit: 50placeholder)

placeholder
	require.Equal(t, 1, result.Blocked)
	require.Equal(t, 1, result.Partial)
	require.Zero(t, result.Failed)
	require.Equal(t, GrokOAuthReconcileOutcomePartial, result.Items[0].Outcome)
placeholder
