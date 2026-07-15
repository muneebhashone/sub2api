//go:build unit

package service

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// ---------- mock helpers ----------

// refreshAPIAccountRepo implements AccountRepository for OAuthRefreshAPI tests.
type refreshAPIAccountRepo struct {
	mockAccountRepoForGemini
	account                 *Account // returned by GetByID
	getByIDErr              error
	getByIDCalls            int
	getByIDErrAfterCall     int
	getByIDErrAfterCallErr  error
	updateErr               error
	updateCalls             int
	updateCredentialsCalls  int
	successCASCalls         int
	beforeSuccessCAS        func(*refreshAPIAccountRepo)
	lastExpectedCredentials map[string]any
	lastExpectedProxyID     *int64
placeholder

func (r *refreshAPIAccountRepo) GetByID(_ context.Context, _ int64) (*Account, error) {
	r.getByIDCalls++
	if r.getByIDErrAfterCall > 0 && r.getByIDCalls >= r.getByIDErrAfterCall {
		return nil, r.getByIDErrAfterCallErr
placeholder
	if r.getByIDErr != nil {
		return nil, r.getByIDErr
placeholder
	return activeRefreshAPITestAccount(r.account), nil
placeholder

func activeRefreshAPITestAccount(account *Account) *Account {
	if account == nil || account.Status != "" {
		return account
placeholder
	copy := *account
	copy.Status = StatusActive
	return &copy
placeholder

func (r *refreshAPIAccountRepo) Update(_ context.Context, _ *Account) error {
	r.updateCalls++
	return r.updateErr
placeholder

func (r *refreshAPIAccountRepo) UpdateCredentials(_ context.Context, id int64, credentials map[string]any) error {
	r.updateCalls++
	r.updateCredentialsCalls++
	if r.updateErr != nil {
		return r.updateErr
placeholder
	if r.account == nil || r.account.ID != id {
		r.account = &Account{ID: idplaceholder
placeholder
	r.account.Credentials = shallowCopyMap(credentials)
	return nil
placeholder

func (r *refreshAPIAccountRepo) UpdateGrokOAuthCredentialsIfUnchanged(
	_ context.Context,
	id int64,
	expectedCredentials map[string]any,
	expectedProxyID *int64,
	credentials map[string]any,
) (bool, error) {
	r.successCASCalls++
	r.lastExpectedCredentials = shallowCopyMap(expectedCredentials)
	if expectedProxyID != nil {
		proxyID := *expectedProxyID
		r.lastExpectedProxyID = &proxyID
placeholder else {
		r.lastExpectedProxyID = nil
placeholder
	if r.beforeSuccessCAS != nil {
		r.beforeSuccessCAS(r)
placeholder
	if r.updateErr != nil {
		return false, r.updateErr
placeholder
	if r.account == nil || r.account.ID != id || r.account.Platform != PlatformGrok ||
		r.account.Type != AccountTypeOAuth ||
		!reflect.DeepEqual(r.account.Credentials, expectedCredentials) ||
		!reflect.DeepEqual(r.account.ProxyID, expectedProxyID) {
		return false, nil
placeholder
	r.updateCalls++
	r.updateCredentialsCalls++
	r.account.Credentials = shallowCopyMap(credentials)
	return true, nil
placeholder

// refreshAPIExecutorStub implements OAuthRefreshExecutor for tests.
type refreshAPIExecutorStub struct {
	needsRefresh  bool
	cannotRefresh bool
	credentials   map[string]any
	err           error
	refreshCalls  int
	canRefresh    func(*Account) bool
	onRefresh     func()
	delay         time.Duration
placeholder

func (e *refreshAPIExecutorStub) CanRefresh(account *Account) bool {
	if e.cannotRefresh {
		return false
placeholder
	if e.canRefresh != nil {
		return e.canRefresh(account)
placeholder
	return true
placeholder

func (e *refreshAPIExecutorStub) NeedsRefresh(_ *Account, _ time.Duration) bool {
	return e.needsRefresh
placeholder

func (e *refreshAPIExecutorStub) Refresh(_ context.Context, _ *Account) (map[string]any, error) {
	e.refreshCalls++
	if e.delay > 0 {
		time.Sleep(e.delay)
placeholder
	if e.onRefresh != nil {
		e.onRefresh()
placeholder
	if e.err != nil {
		return nil, e.err
placeholder
	return e.credentials, nil
placeholder

func (e *refreshAPIExecutorStub) CacheKey(account *Account) string {
	return "test:api:" + account.Platform
placeholder

// refreshAPICacheStub implements GeminiTokenCache for OAuthRefreshAPI tests.
type refreshAPICacheStub struct {
	lockResult    bool
	lockErr       error
	releaseCalls  int
	releaseCtxErr error
	deleteCalls   int
	deleteKey     string
	deleteCtxErr  error
placeholder

func (c *refreshAPICacheStub) GetAccessToken(context.Context, string) (string, error) {
	return "", nil
placeholder

func (c *refreshAPICacheStub) SetAccessToken(context.Context, string, string, time.Duration) error {
	return nil
placeholder

func (c *refreshAPICacheStub) DeleteAccessToken(ctx context.Context, key string) error {
	c.deleteCalls++
	c.deleteKey = key
	c.deleteCtxErr = ctx.Err()
	return nil
placeholder

func (c *refreshAPICacheStub) AcquireRefreshLock(context.Context, string, time.Duration) (bool, error) {
	return c.lockResult, c.lockErr
placeholder

func (c *refreshAPICacheStub) ReleaseRefreshLock(ctx context.Context, _ string) error {
	c.releaseCalls++
	c.releaseCtxErr = ctx.Err()
	return nil
placeholder

// ========== RefreshIfNeeded tests ==========

func TestRefreshIfNeeded_Success(t *testing.T) {
	account := &Account{ID: 1, Platform: PlatformAnthropic, Type: AccountTypeOAuth, Status: StatusActiveplaceholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		credentials:  map[string]any{"access_token": "new-token"placeholder,
placeholder

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

placeholder
	require.True(t, result.Refreshed)
	require.NotNil(t, result.NewCredentials)
	require.Equal(t, "new-token", result.NewCredentials["access_token"])
	require.NotNil(t, result.NewCredentials["_token_version"]) // version stamp set
	require.Equal(t, 1, repo.updateCalls)                      // DB updated
	require.Equal(t, 1, repo.updateCredentialsCalls)
	require.Equal(t, 1, cache.releaseCalls) // lock released
	require.Equal(t, 1, executor.refreshCalls)
placeholder

func TestRefreshIfNeeded_UpdateCredentialsPreservesRateLimitState(t *testing.T) {
	resetAt := time.Now().Add(45 * time.Minute)
	account := &Account{
		ID:               11,
		Platform:         PlatformGemini,
		Type:             AccountTypeOAuth,
		Status:           StatusActive,
		RateLimitResetAt: &resetAt,
placeholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		credentials:  map[string]any{"access_token": "safe-token"placeholder,
placeholder

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

placeholder
	require.True(t, result.Refreshed)
	require.Equal(t, 1, repo.updateCredentialsCalls)
	require.NotNil(t, repo.account.RateLimitResetAt)
	require.WithinDuration(t, resetAt, *repo.account.RateLimitResetAt, time.Second)
placeholder

func TestRefreshIfNeeded_LockHeld(t *testing.T) {
	account := &Account{ID: 2, Platform: PlatformAnthropic, Status: StatusActiveplaceholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	cache := &refreshAPICacheStub{lockResult: falseplaceholder // lock not acquired
	executor := &refreshAPIExecutorStub{needsRefresh: trueplaceholder

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

placeholder
	require.True(t, result.LockHeld)
	require.False(t, result.Refreshed)
	require.Equal(t, 0, repo.updateCalls)
	require.Equal(t, 0, executor.refreshCalls)
placeholder

func TestRefreshIfNeeded_LockErrorDegrades(t *testing.T) {
	account := &Account{ID: 3, Platform: PlatformGemini, Type: AccountTypeOAuth, Status: StatusActiveplaceholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	cache := &refreshAPICacheStub{lockErr: errors.New("redis down")placeholder // lock error
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		credentials:  map[string]any{"access_token": "degraded-token"placeholder,
placeholder

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

placeholder
	require.True(t, result.Refreshed)       // still refreshed (degraded mode)
	require.Equal(t, 1, repo.updateCalls)   // DB updated
	require.Equal(t, 0, cache.releaseCalls) // no lock to release
	require.Equal(t, 1, executor.refreshCalls)
placeholder

func TestRefreshIfNeeded_NoCacheNoLock(t *testing.T) {
	account := &Account{ID: 4, Platform: PlatformGemini, Type: AccountTypeOAuth, Status: StatusActiveplaceholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		credentials:  map[string]any{"access_token": "no-cache-token"placeholder,
placeholder

	api := NewOAuthRefreshAPI(repo, nil) // no cache = no lock
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

placeholder
	require.True(t, result.Refreshed)
	require.Equal(t, 1, repo.updateCalls)
placeholder

func TestRefreshIfNeeded_AlreadyRefreshed(t *testing.T) {
	account := &Account{ID: 5, Platform: PlatformAnthropic, Status: StatusActiveplaceholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	executor := &refreshAPIExecutorStub{needsRefresh: falseplaceholder // already refreshed

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

placeholder
	require.False(t, result.Refreshed)
	require.False(t, result.LockHeld)
	require.NotNil(t, result.Account) // returns fresh account
	require.Equal(t, 0, repo.updateCalls)
	require.Equal(t, 0, executor.refreshCalls)
placeholder

func TestRefreshIfNeeded_RefreshError(t *testing.T) {
	account := &Account{ID: 6, Platform: PlatformAnthropic, Status: StatusActiveplaceholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		err:          errors.New("invalid_grant: token revoked"),
placeholder

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

placeholder
	require.NotNil(t, result)
	require.NotNil(t, result.Account)
	require.Equal(t, account.ID, result.Account.ID)
	require.Contains(t, err.Error(), "invalid_grant")
	require.Equal(t, 0, repo.updateCalls)   // no DB update on refresh error
	require.Equal(t, 1, cache.releaseCalls) // lock still released via defer
placeholder

func TestRefreshIfNeeded_DBUpdateError(t *testing.T) {
	account := &Account{ID: 7, Platform: PlatformGemini, Type: AccountTypeOAuth, Status: StatusActiveplaceholder
	repo := &refreshAPIAccountRepo{
		account:   account,
		updateErr: errors.New("db connection lost"),
placeholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		credentials:  map[string]any{"access_token": "token"placeholder,
placeholder

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

placeholder
	require.Nil(t, result)
	require.ErrorIs(t, err, errOAuthRefreshCredentialPersist)
	require.Equal(t, 1, repo.updateCalls) // attempted
placeholder

func TestRefreshIfNeeded_GrokSuccessCASLetsConcurrentReauthorizationWin(t *testing.T) {
	proxyID := int64(17)
	account := &Account{
		ID:       70,
		Platform: PlatformGrok,
		Type:     AccountTypeOAuth,
		Status:   StatusActive,
		ProxyID:  &proxyID,
placeholder
			"access_token":   "attempted-access",
			"refresh_token":  "attempted-refresh",
			"_token_version": int64(1),
	placeholder,
placeholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	repo.beforeSuccessCAS = func(r *refreshAPIAccountRepo) {
		repairedProxyID := int64(23)
		r.account.ProxyID = &repairedProxyID
		r.account.Credentials = map[string]any{
			"access_token":   "reauthorized-access",
			"refresh_token":  "reauthorized-refresh",
			"_token_version": int64(2),
	placeholder
placeholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		credentials: map[string]any{
			"access_token":  "provider-access",
			"refresh_token": "provider-refresh",
	placeholder,
placeholder

	result, err := NewOAuthRefreshAPI(repo, nil).RefreshIfNeeded(context.Background(), account, executor, time.Hour)

placeholder
	require.NotNil(t, result)
	require.False(t, result.Refreshed, "a lost success CAS is an already-refreshed skip")
	require.Nil(t, result.NewCredentials)
	require.Equal(t, "reauthorized-refresh", result.Account.GetGrokRefreshToken())
	require.NotNil(t, result.Account.ProxyID)
	require.Equal(t, int64(23), *result.Account.ProxyID)
	require.Equal(t, 1, repo.successCASCalls)
	require.Equal(t, "attempted-refresh", repo.lastExpectedCredentials["refresh_token"])
	require.NotNil(t, repo.lastExpectedProxyID)
	require.Equal(t, proxyID, *repo.lastExpectedProxyID)
	require.Zero(t, repo.updateCredentialsCalls, "the provider result must not overwrite a concurrent repair")
placeholder

func TestRefreshIfNeeded_GrokSuccessPersistenceFailureIsProviderContainment(t *testing.T) {
	account := &Account{
		ID:       71,
		Platform: PlatformGrok,
		Type:     AccountTypeOAuth,
		Status:   StatusActive,
placeholder
			"access_token":  "attempted-access",
			"refresh_token": "attempted-refresh",
	placeholder,
placeholder
	repo := &refreshAPIAccountRepo{account: account, updateErr: errors.New("database unavailable")placeholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		credentials: map[string]any{
			"access_token":  "provider-access",
			"refresh_token": "provider-refresh",
	placeholder,
placeholder

	result, err := NewOAuthRefreshAPI(repo, nil).RefreshIfNeeded(context.Background(), account, executor, time.Hour)

placeholder
	require.Nil(t, result)
	var containmentErr *providerCycleContainmentRefreshError
	require.ErrorAs(t, err, &containmentErr)
	require.Equal(t, "attempted-refresh", account.GetGrokRefreshToken(),
		"an ambiguous persistence result must not mutate the in-memory account")
	require.Equal(t, 1, repo.successCASCalls)
	require.Zero(t, repo.updateCredentialsCalls)
placeholder

func TestRefreshIfNeeded_GrokSuccessDurableRereadFailureIsProviderContainment(t *testing.T) {
	account := &Account{
		ID:       72,
		Platform: PlatformGrok,
		Type:     AccountTypeOAuth,
		Status:   StatusActive,
placeholder
			"access_token":  "attempted-access",
			"refresh_token": "attempted-refresh",
	placeholder,
placeholder
	repo := &refreshAPIAccountRepo{
		account:                account,
		getByIDErrAfterCall:    2,
		getByIDErrAfterCallErr: errors.New("durable state unavailable"),
placeholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		credentials: map[string]any{
			"access_token":  "provider-access",
			"refresh_token": "provider-refresh",
	placeholder,
placeholder

	result, err := NewOAuthRefreshAPI(repo, cache).RefreshIfNeeded(context.Background(), account, executor, time.Hour)

placeholder
	require.Nil(t, result)
	var containmentErr *providerCycleContainmentRefreshError
	require.ErrorAs(t, err, &containmentErr)
	require.Equal(t, 2, repo.getByIDCalls)
	require.Equal(t, 1, repo.successCASCalls)
	require.Equal(t, 1, cache.deleteCalls, "a committed credential rotation must invalidate the pre-rotation access-token cache")
	require.NoError(t, cache.deleteCtxErr)
placeholder

func TestRefreshIfNeeded_DBRereadFails(t *testing.T) {
	account := &Account{ID: 8, Platform: PlatformAnthropic, Type: AccountTypeOAuth, Status: StatusActiveplaceholder
	repo := &refreshAPIAccountRepo{
		account:    nil, // GetByID returns nil
		getByIDErr: errors.New("db timeout"),
placeholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		credentials:  map[string]any{"access_token": "fallback-token"placeholder,
placeholder

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

placeholder
	var stateUnavailable *oauthRefreshStateUnavailableError
	require.ErrorAs(t, err, &stateUnavailable)
	require.Nil(t, result)
	require.Zero(t, executor.refreshCalls, "a failed DB reread must not refresh stale credentials")
	require.Zero(t, repo.updateCalls)
	require.Equal(t, 1, cache.releaseCalls)
placeholder

func TestRefreshIfNeeded_RequestPathDBRereadNilFailsClosed(t *testing.T) {
	account := &Account{ID: 81, Platform: PlatformGrok, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: trueplaceholder
	repo := &refreshAPIAccountRepo{placeholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	executor := &refreshAPIExecutorStub{needsRefresh: trueplaceholder

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(withOAuthRefreshRequestPath(context.Background()), account, executor, 3*time.Minute)

	require.ErrorIs(t, err, errOAuthRefreshAccountStateChanged)
	require.Nil(t, result)
	require.Zero(t, executor.refreshCalls)
	require.Zero(t, repo.updateCalls)
	require.Equal(t, 1, cache.releaseCalls)
placeholder

func TestRefreshIfNeeded_RequestPathDBRereadInactiveFailsClosed(t *testing.T) {
	account := &Account{ID: 82, Platform: PlatformGrok, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: trueplaceholder
	freshAccount := &Account{ID: account.ID, Platform: PlatformGrok, Type: AccountTypeOAuth, Status: StatusDisabledplaceholder
	repo := &refreshAPIAccountRepo{account: freshAccountplaceholder
	executor := &refreshAPIExecutorStub{needsRefresh: trueplaceholder

	api := NewOAuthRefreshAPI(repo, nil)
	result, err := api.RefreshIfNeeded(withOAuthRefreshRequestPath(context.Background()), account, executor, 3*time.Minute)

	require.ErrorContains(t, err, "account is not active")
	require.Nil(t, result)
	require.Zero(t, executor.refreshCalls)
	require.Zero(t, repo.updateCalls)
placeholder

func TestRefreshIfNeeded_RequestPathDBRereadRevalidatesExecutorContract(t *testing.T) {
	tests := []struct {
		name          string
		freshPlatform string
		freshType     string
placeholder{
		{name: "platform changed", freshPlatform: PlatformAnthropic, freshType: AccountTypeOAuthplaceholder,
		{name: "type changed", freshPlatform: PlatformGrok, freshType: AccountTypeUpstreamplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := &Account{ID: 83, Platform: PlatformGrok, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: trueplaceholder
			freshAccount := &Account{ID: account.ID, Platform: tt.freshPlatform, Type: tt.freshType, Status: StatusActive, Schedulable: trueplaceholder
			repo := &refreshAPIAccountRepo{account: freshAccountplaceholder
			executor := NewGrokTokenRefresher(nil)

			api := NewOAuthRefreshAPI(repo, nil)
			result, err := api.RefreshIfNeeded(withOAuthRefreshRequestPath(context.Background()), account, executor, 3*time.Minute)

			require.ErrorIs(t, err, errOAuthRefreshAccountStateChanged)
			require.Nil(t, result)
			require.Zero(t, repo.updateCalls)
	placeholder)
placeholder
placeholder

func TestRefreshIfNeeded_LocalLockWaitHonorsContext(t *testing.T) {
	account := &Account{ID: 80, Platform: PlatformGrok, Type: AccountTypeOAuth, Status: StatusActiveplaceholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	executor := &refreshAPIExecutorStub{needsRefresh: trueplaceholder
	api := NewOAuthRefreshAPI(repo, nil)
	lock := api.getLocalLock(executor.CacheKey(account))
	require.NoError(t, lock.Lock(context.Background()))
	defer lock.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	result, err := api.RefreshIfNeeded(ctx, account, executor, time.Hour)

	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Nil(t, result)
	require.Zero(t, executor.refreshCalls)
placeholder

func TestRefreshIfNeeded_ReleasesDistributedLockAfterParentCancellation(t *testing.T) {
	account := &Account{ID: 81, Platform: PlatformGrok, Type: AccountTypeOAuth, Status: StatusActiveplaceholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	ctx, cancel := context.WithCancel(context.Background())
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		err:          errors.New("temporary provider error"),
		onRefresh:    cancel,
placeholder
	api := NewOAuthRefreshAPI(repo, cache)

	_, err := api.RefreshIfNeeded(ctx, account, executor, time.Hour)

placeholder
	require.Equal(t, 1, cache.releaseCalls)
	require.NoError(t, cache.releaseCtxErr, "lock cleanup must not reuse the canceled attempt context")
placeholder

func TestRefreshIfNeeded_RevalidatesFreshAccountBeforeRefresh(t *testing.T) {
	selected := &Account{ID: 82, Platform: PlatformGrok, Type: AccountTypeOAuth, Status: StatusActiveplaceholder
	tests := []struct {
		name  string
		fresh *Account
placeholder{
		{name: "converted to API key", fresh: &Account{ID: 82, Platform: PlatformGrok, Type: AccountTypeAPIKey, Status: StatusActiveplaceholderplaceholder,
		{name: "disabled", fresh: &Account{ID: 82, Platform: PlatformGrok, Type: AccountTypeOAuth, Status: StatusDisabledplaceholderplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &refreshAPIAccountRepo{account: tt.freshplaceholder
			executor := &refreshAPIExecutorStub{
				needsRefresh: true,
				canRefresh: func(account *Account) bool {
					return account.Platform == PlatformGrok && account.Type == AccountTypeOAuth
			placeholder,
		placeholder
			api := NewOAuthRefreshAPI(repo, nil)

			result, err := api.RefreshIfNeeded(context.Background(), selected, executor, time.Hour)

		placeholder
			require.False(t, result.Refreshed)
			require.Zero(t, executor.refreshCalls)
			require.Zero(t, repo.updateCalls)
	placeholder)
placeholder
placeholder

func TestRefreshIfNeeded_RequestPathDBRereadMissingGrokRefreshCredentialReturnsPermanentSignal(t *testing.T) {
	account := &Account{
		ID:          84,
		Platform:    PlatformGrok,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Schedulable: true,
placeholder
			"refresh_token": "caller-snapshot-refresh-token",
	placeholder,
placeholder
	freshAccount := &Account{ID: account.ID, Platform: PlatformGrok, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: trueplaceholder
	repo := &refreshAPIAccountRepo{account: freshAccountplaceholder
	executor := NewGrokTokenRefresher(nil)

	api := NewOAuthRefreshAPI(repo, nil)
	result, err := api.RefreshIfNeeded(withOAuthRefreshRequestPath(context.Background()), account, executor, 3*time.Minute)

	require.ErrorIs(t, err, errGrokOAuthRefreshTokenMissing)
	require.Nil(t, result)
	require.Zero(t, repo.updateCalls)
placeholder

func TestRefreshIfNeeded_LateSuccessAfterDeadlineDoesNotPersist(t *testing.T) {
	account := &Account{ID: 85, Platform: PlatformGrok, Type: AccountTypeOAuth, Status: StatusActiveplaceholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		credentials:  map[string]any{"access_token": "late-token"placeholder,
		delay:        30 * time.Millisecond,
placeholder
	api := NewOAuthRefreshAPI(repo, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	result, err := api.RefreshIfNeeded(ctx, account, executor, time.Hour)

	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Nil(t, result)
	require.Zero(t, repo.updateCredentialsCalls, "late credentials must not cross the unified API persistence boundary")
placeholder

func TestRefreshIfNeeded_NilCredentials(t *testing.T) {
	account := &Account{ID: 9, Platform: PlatformGemini, Type: AccountTypeOAuth, Status: StatusActiveplaceholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		credentials:  nil, // Refresh returns nil credentials
placeholder

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

placeholder
	require.True(t, result.Refreshed)
	require.Nil(t, result.NewCredentials)
	require.Equal(t, 0, repo.updateCalls) // no DB update when credentials are nil
placeholder

// ========== MergeCredentials tests ==========

func TestMergeCredentials_Basic(t *testing.T) {
	old := map[string]any{"a": "1", "b": "2", "c": "3"placeholder
	new := map[string]any{"a": "new", "d": "4"placeholder

	result := MergeCredentials(old, new)

	require.Equal(t, "new", result["a"]) // new value preserved
	require.Equal(t, "2", result["b"])   // old value kept
	require.Equal(t, "3", result["c"])   // old value kept
	require.Equal(t, "4", result["d"])   // new value preserved
placeholder

func TestMergeCredentials_NilNew(t *testing.T) {
	old := map[string]any{"a": "1"placeholder

	result := MergeCredentials(old, nil)

	require.NotNil(t, result)
	require.Equal(t, "1", result["a"])
placeholder

func TestMergeCredentials_NilOld(t *testing.T) {
	new := map[string]any{"a": "1"placeholder

	result := MergeCredentials(nil, new)

	require.Equal(t, "1", result["a"])
placeholder

func TestMergeCredentials_BothNil(t *testing.T) {
	result := MergeCredentials(nil, nil)
	require.NotNil(t, result)
	require.Empty(t, result)
placeholder

func TestMergeCredentials_NewOverridesOld(t *testing.T) {
	old := map[string]any{"access_token": "old-token", "refresh_token": "old-refresh"placeholder
	new := map[string]any{"access_token": "new-token"placeholder

	result := MergeCredentials(old, new)

	require.Equal(t, "new-token", result["access_token"])    // overridden
	require.Equal(t, "old-refresh", result["refresh_token"]) // preserved
placeholder

// ========== BuildClaudeAccountCredentials tests ==========

func TestBuildClaudeAccountCredentials_Full(t *testing.T) {
	tokenInfo := &TokenInfo{
		AccessToken:  "at-123",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		ExpiresAt:    1700000000,
		RefreshToken: "rt-456",
		Scope:        "openid",
placeholder

	creds := BuildClaudeAccountCredentials(tokenInfo)

	require.Equal(t, "at-123", creds["access_token"])
	require.Equal(t, "Bearer", creds["token_type"])
	require.Equal(t, "3600", creds["expires_in"])
	require.Equal(t, "1700000000", creds["expires_at"])
	require.Equal(t, "rt-456", creds["refresh_token"])
	require.Equal(t, "openid", creds["scope"])
placeholder

func TestBuildClaudeAccountCredentials_Minimal(t *testing.T) {
	tokenInfo := &TokenInfo{
		AccessToken: "at-789",
		TokenType:   "Bearer",
		ExpiresIn:   7200,
		ExpiresAt:   1700003600,
placeholder

	creds := BuildClaudeAccountCredentials(tokenInfo)

	require.Equal(t, "at-789", creds["access_token"])
	require.Equal(t, "Bearer", creds["token_type"])
	require.Equal(t, "7200", creds["expires_in"])
	require.Equal(t, "1700003600", creds["expires_at"])
	_, hasRefresh := creds["refresh_token"]
	_, hasScope := creds["scope"]
	require.False(t, hasRefresh, "refresh_token should not be set when empty")
	require.False(t, hasScope, "scope should not be set when empty")
placeholder

// refreshAPIAccountRepoWithRace supports returning a different account on subsequent GetByID calls
// to simulate race conditions where another worker has refreshed the token.
type refreshAPIAccountRepoWithRace struct {
	refreshAPIAccountRepo
	raceAccount  *Account // returned on 2nd+ GetByID call
	getByIDCalls int
placeholder

func (r *refreshAPIAccountRepoWithRace) GetByID(_ context.Context, _ int64) (*Account, error) {
	r.getByIDCalls++
	if r.getByIDCalls > 1 && r.raceAccount != nil {
		return activeRefreshAPITestAccount(r.raceAccount), nil
placeholder
	if r.getByIDErr != nil {
		return nil, r.getByIDErr
placeholder
	return activeRefreshAPITestAccount(r.account), nil
placeholder

// ========== Race recovery tests ==========

func TestRefreshIfNeeded_InvalidGrantRaceRecovered(t *testing.T) {
	// Account with old refresh token
	account := &Account{
		ID:          10,
		Platform:    PlatformAnthropic,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
placeholder"refresh_token": "old-rt", "access_token": "old-at"placeholder,
placeholder
	// After race, DB has new refresh token from another worker
	racedAccount := &Account{
		ID:          10,
		Platform:    PlatformAnthropic,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
placeholder"refresh_token": "new-rt", "access_token": "new-at"placeholder,
placeholder
	repo := &refreshAPIAccountRepoWithRace{
		refreshAPIAccountRepo: refreshAPIAccountRepo{account: accountplaceholder,
		raceAccount:           racedAccount,
placeholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		err:          errors.New("invalid_grant: refresh token not found or invalid"),
placeholder

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

	require.NoError(t, err, "race-recovered invalid_grant should not return error")
	require.False(t, result.Refreshed)
	require.False(t, result.LockHeld)
	require.NotNil(t, result.Account)
	require.Equal(t, "new-rt", result.Account.GetCredential("refresh_token"))
	require.Equal(t, 0, repo.updateCalls) // no DB update needed, another worker did it
placeholder

func TestRefreshIfNeeded_InvalidGrantGenuine(t *testing.T) {
	// Account with revoked refresh token - DB still has the same token
	account := &Account{
		ID:          11,
		Platform:    PlatformAnthropic,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
placeholder"refresh_token": "revoked-rt", "access_token": "old-at"placeholder,
placeholder
	repo := &refreshAPIAccountRepoWithRace{
		refreshAPIAccountRepo: refreshAPIAccountRepo{account: accountplaceholder,
		raceAccount:           account, // same refresh_token on re-read
placeholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		err:          errors.New("invalid_grant: refresh token revoked"),
placeholder

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

	require.Error(t, err, "genuine invalid_grant should propagate error")
	require.NotNil(t, result)
	require.NotNil(t, result.Account)
	require.Equal(t, "revoked-rt", result.Account.GetCredential("refresh_token"))
	require.Contains(t, err.Error(), "invalid_grant")
placeholder

func TestRefreshIfNeeded_InvalidGrantDBRereadFailsOnRecovery(t *testing.T) {
	account := &Account{
		ID:          12,
		Platform:    PlatformAnthropic,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
placeholder"refresh_token": "old-rt"placeholder,
placeholder
	repo := &refreshAPIAccountRepoWithRace{
		refreshAPIAccountRepo: refreshAPIAccountRepo{account: accountplaceholder,
		raceAccount:           nil, // GetByID returns nil on recovery attempt
placeholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	executor := &refreshAPIExecutorStub{
		needsRefresh: true,
		err:          errors.New("invalid_grant"),
placeholder

	api := NewOAuthRefreshAPI(repo, cache)
	result, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)

	require.Error(t, err, "should propagate error when recovery DB re-read fails")
	require.NotNil(t, result)
	require.NotNil(t, result.Account)
	require.Equal(t, "old-rt", result.Account.GetCredential("refresh_token"))
placeholder

func TestRefreshIfNeeded_LocalMutexSerializesConcurrent(t *testing.T) {
	// Test that two goroutines for the same account are serialized by the local mutex.
	// The first goroutine refreshes successfully; the second sees NeedsRefresh=false.
	refreshed := &Account{
		ID:          20,
		Platform:    PlatformAnthropic,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
placeholder"refresh_token": "new-rt", "access_token": "new-at"placeholder,
placeholder
	callCount := 0
	repo := &refreshAPIAccountRepo{account: &Account{
		ID:          20,
		Platform:    PlatformAnthropic,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
placeholder"refresh_token": "old-rt"placeholder,
placeholderplaceholder

	// After first refresh, NeedsRefresh should return false
	// We simulate this by using an executor that decrements needsRefresh after first call
	var mu sync.Mutex
	dynamicExecutor := &dynamicRefreshExecutor{
		canRefresh: true,
		cacheKey:   "test:mutex:anthropic",
		refreshFunc: func(_ context.Context, _ *Account) (map[string]any, error) {
			mu.Lock()
			callCount++
			mu.Unlock()
			time.Sleep(50 * time.Millisecond) // slow refresh
			return map[string]any{"access_token": "new-at"placeholder, nil
	placeholder,
		needsRefreshFunc: func() bool {
			mu.Lock()
			defer mu.Unlock()
			return callCount == 0 // only first call needs refresh
	placeholder,
placeholder

	_ = refreshed

	api := NewOAuthRefreshAPI(repo, nil) // no distributed lock, only local mutex

	var wg sync.WaitGroup
	results := make([]*OAuthRefreshResult, 2)
	errs := make([]error, 2)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx], errs[idx] = api.RefreshIfNeeded(context.Background(), repo.account, dynamicExecutor, 3*time.Minute)
	placeholder(i)
placeholder
	wg.Wait()

	require.NoError(t, errs[0])
	require.NoError(t, errs[1])

	// Only one goroutine should have actually called Refresh
	mu.Lock()
	require.Equal(t, 1, callCount, "only one refresh call should have been made")
	mu.Unlock()
placeholder

func TestRefreshIfNeeded_LocalLockWaitHonorsContextCancellation(t *testing.T) {
	account := &Account{ID: 21, Platform: PlatformGrok, Type: AccountTypeOAuth, Status: StatusActiveplaceholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	refreshStarted := make(chan struct{placeholder)
	releaseRefresh := make(chan struct{placeholder)
	var once sync.Once
	executor := &dynamicRefreshExecutor{
		canRefresh:       true,
		cacheKey:         "test:context-lock:grok",
		needsRefreshFunc: func() bool { return true placeholder,
		refreshFunc: func(context.Context, *Account) (map[string]any, error) {
			once.Do(func() { close(refreshStarted) placeholder)
			<-releaseRefresh
			return map[string]any{"access_token": "new-at"placeholder, nil
	placeholder,
placeholder
	api := NewOAuthRefreshAPI(repo, nil)
	firstDone := make(chan error, 1)
	go func() {
		_, err := api.RefreshIfNeeded(context.Background(), account, executor, 3*time.Minute)
		firstDone <- err
placeholder()
	<-refreshStarted

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	startedAt := time.Now()
	result, err := api.RefreshIfNeeded(ctx, account, executor, 3*time.Minute)

	require.Nil(t, result)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Less(t, time.Since(startedAt), 500*time.Millisecond)
	close(releaseRefresh)
	require.NoError(t, <-firstDone)
placeholder

func TestRefreshIfNeeded_ReleasesDistributedLockWithCleanupContext(t *testing.T) {
	account := &Account{
		ID:       22,
		Platform: PlatformGrok,
		Type:     AccountTypeOAuth,
		Status:   StatusActive,
placeholder
			"access_token":  "old-access",
			"refresh_token": "old-refresh",
	placeholder,
placeholder
	repo := &refreshAPIAccountRepo{account: accountplaceholder
	cache := &refreshAPICacheStub{lockResult: trueplaceholder
	ctx, cancel := context.WithCancel(context.Background())
	executor := &dynamicRefreshExecutor{
		canRefresh:       true,
		cacheKey:         "test:cleanup:grok",
		needsRefreshFunc: func() bool { return true placeholder,
		refreshFunc: func(context.Context, *Account) (map[string]any, error) {
			cancel()
			return map[string]any{"access_token": "new-at"placeholder, nil
	placeholder,
placeholder
	api := NewOAuthRefreshAPI(repo, cache)

	result, err := api.RefreshIfNeeded(ctx, account, executor, 3*time.Minute)

	require.ErrorIs(t, err, context.Canceled)
	require.Nil(t, result)
	require.Zero(t, repo.updateCalls)
	require.Equal(t, "old-access", account.GetGrokAccessToken())
	require.Zero(t, account.GetCredentialAsInt64("_token_version"))
	require.Equal(t, 1, cache.releaseCalls)
	require.NoError(t, cache.releaseCtxErr)
placeholder

// dynamicRefreshExecutor is a test helper with function-based NeedsRefresh and Refresh.
type dynamicRefreshExecutor struct {
	canRefresh       bool
	cacheKey         string
	needsRefreshFunc func() bool
	refreshFunc      func(context.Context, *Account) (map[string]any, error)
placeholder

func (e *dynamicRefreshExecutor) CanRefresh(_ *Account) bool { return e.canRefresh placeholder

func (e *dynamicRefreshExecutor) NeedsRefresh(_ *Account, _ time.Duration) bool {
	return e.needsRefreshFunc()
placeholder

func (e *dynamicRefreshExecutor) Refresh(ctx context.Context, account *Account) (map[string]any, error) {
	return e.refreshFunc(ctx, account)
placeholder

func (e *dynamicRefreshExecutor) CacheKey(_ *Account) string {
	return e.cacheKey
placeholder

// ========== NewOAuthRefreshAPI TTL tests ==========

func TestNewOAuthRefreshAPI_DefaultTTL(t *testing.T) {
	api := NewOAuthRefreshAPI(nil, nil)
	require.Equal(t, defaultRefreshLockTTL, api.lockTTL)
placeholder

func TestNewOAuthRefreshAPI_CustomTTL(t *testing.T) {
	api := NewOAuthRefreshAPI(nil, nil, 90*time.Second)
	require.Equal(t, 90*time.Second, api.lockTTL)
placeholder

func TestNewOAuthRefreshAPI_ZeroTTLUsesDefault(t *testing.T) {
	api := NewOAuthRefreshAPI(nil, nil, 0)
	require.Equal(t, defaultRefreshLockTTL, api.lockTTL)
placeholder

// ========== isInvalidGrantError tests ==========

func TestIsInvalidGrantError(t *testing.T) {
	require.True(t, isInvalidGrantError(errors.New("invalid_grant: token revoked")))
	require.True(t, isInvalidGrantError(errors.New("INVALID_GRANT")))
	require.False(t, isInvalidGrantError(errors.New("invalid_client")))
	require.False(t, isInvalidGrantError(nil))
placeholder

// ========== BackgroundRefreshPolicy tests ==========

func TestBackgroundRefreshPolicy_DefaultSkips(t *testing.T) {
	p := DefaultBackgroundRefreshPolicy()

	require.ErrorIs(t, p.handleLockHeld(), errRefreshSkipped)
	require.ErrorIs(t, p.handleAlreadyRefreshed(), errRefreshSkipped)
placeholder

func TestBackgroundRefreshPolicy_SuccessOverride(t *testing.T) {
	p := BackgroundRefreshPolicy{
		OnLockHeld:       BackgroundSkipAsSuccess,
		OnAlreadyRefresh: BackgroundSkipAsSuccess,
placeholder

	require.NoError(t, p.handleLockHeld())
	require.NoError(t, p.handleAlreadyRefreshed())
placeholder

// ========== ProviderRefreshPolicy tests ==========

func TestClaudeProviderRefreshPolicy(t *testing.T) {
	p := ClaudeProviderRefreshPolicy()
	require.Equal(t, ProviderRefreshErrorUseExistingToken, p.OnRefreshError)
	require.Equal(t, ProviderLockHeldWaitForCache, p.OnLockHeld)
	require.Equal(t, time.Minute, p.FailureTTL)
placeholder

func TestOpenAIProviderRefreshPolicy(t *testing.T) {
	p := OpenAIProviderRefreshPolicy()
	require.Equal(t, ProviderRefreshErrorUseExistingToken, p.OnRefreshError)
	require.Equal(t, ProviderLockHeldWaitForCache, p.OnLockHeld)
	require.Equal(t, time.Minute, p.FailureTTL)
placeholder

func TestGeminiProviderRefreshPolicy(t *testing.T) {
	p := GeminiProviderRefreshPolicy()
	require.Equal(t, ProviderRefreshErrorReturn, p.OnRefreshError)
	require.Equal(t, ProviderLockHeldUseExistingToken, p.OnLockHeld)
	require.Equal(t, time.Duration(0), p.FailureTTL)
placeholder

func TestAntigravityProviderRefreshPolicy(t *testing.T) {
	p := AntigravityProviderRefreshPolicy()
	require.Equal(t, ProviderRefreshErrorReturn, p.OnRefreshError)
	require.Equal(t, ProviderLockHeldUseExistingToken, p.OnLockHeld)
	require.Equal(t, time.Duration(0), p.FailureTTL)
placeholder
