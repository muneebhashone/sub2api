//go:build unit

package service

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type grokCredentialPersistingRepo struct {
	*tokenRefreshAccountRepo
placeholder

func (r *grokCredentialPersistingRepo) SetError(ctx context.Context, id int64, message string) error {
	if err := ctx.Err(); err != nil {
		return err
placeholder
	if err := r.tokenRefreshAccountRepo.SetError(ctx, id, message); err != nil {
		return err
placeholder
	if account := r.accountsByID[id]; account != nil {
		account.Status = StatusError
		account.Schedulable = false
		account.ErrorMessage = message
placeholder
	return nil
placeholder

type grokCredentialProxyRepoStub struct {
	ProxyRepository
	proxy *Proxy
	err   error
placeholder

func (r *grokCredentialProxyRepoStub) GetByID(context.Context, int64) (*Proxy, error) {
	return r.proxy, r.err
placeholder

type grokCredentialBlockingRepo struct {
	*tokenRefreshAccountRepo
	setErrorStarted chan struct{placeholder
	setTempStarted  chan struct{placeholder
	onceError       sync.Once
	onceTemp        sync.Once
placeholder

type grokCredentialCommitThenCancelRepo struct {
	*tokenRefreshAccountRepo
	returnErr error
placeholder

type grokCredentialUncommittedDeadlineRepo struct {
	*tokenRefreshAccountRepo
placeholder

func (r *grokCredentialUncommittedDeadlineRepo) SetGrokCredentialErrorIfMatch(
	context.Context,
	int64,
	GrokCredentialMutationSnapshot,
	string,
) (bool, error) {
	return false, context.DeadlineExceeded
placeholder

func (r *grokCredentialUncommittedDeadlineRepo) SetGrokCredentialTempUnschedulableIfMatch(
	context.Context,
	int64,
	GrokCredentialMutationSnapshot,
	time.Time,
	string,
) (bool, error) {
	return false, context.DeadlineExceeded
placeholder

func (r *grokCredentialCommitThenCancelRepo) SetGrokCredentialErrorIfMatch(
	ctx context.Context,
	id int64,
	_ GrokCredentialMutationSnapshot,
	reason string,
) (bool, error) {
	account := r.accountsByID[id]
	account.Status = StatusError
	account.Schedulable = false
	account.ErrorMessage = reason
	if r.returnErr != nil {
		return false, r.returnErr
placeholder
	<-ctx.Done()
	return false, ctx.Err()
placeholder

func (r *grokCredentialCommitThenCancelRepo) SetGrokCredentialTempUnschedulableIfMatch(
	ctx context.Context,
	id int64,
	_ GrokCredentialMutationSnapshot,
	until time.Time,
	reason string,
) (bool, error) {
	account := r.accountsByID[id]
	account.TempUnschedulableUntil = &until
	account.TempUnschedulableReason = reason
	if r.returnErr != nil {
		return false, r.returnErr
placeholder
	<-ctx.Done()
	return false, ctx.Err()
placeholder

func (r *grokCredentialBlockingRepo) SetError(ctx context.Context, _ int64, _ string) error {
	r.onceError.Do(func() { close(r.setErrorStarted) placeholder)
	<-ctx.Done()
	return ctx.Err()
placeholder

func (r *grokCredentialBlockingRepo) SetTempUnschedulable(ctx context.Context, _ int64, _ time.Time, _ string) error {
	r.onceTemp.Do(func() { close(r.setTempStarted) placeholder)
	<-ctx.Done()
	return ctx.Err()
placeholder

func (r *grokCredentialBlockingRepo) SetGrokCredentialErrorIfMatch(
	ctx context.Context,
	_ int64,
	_ GrokCredentialMutationSnapshot,
	_ string,
) (bool, error) {
	r.onceError.Do(func() { close(r.setErrorStarted) placeholder)
	<-ctx.Done()
	return false, ctx.Err()
placeholder

func (r *grokCredentialBlockingRepo) SetGrokCredentialTempUnschedulableIfMatch(
	ctx context.Context,
	_ int64,
	_ GrokCredentialMutationSnapshot,
	_ time.Time,
	_ string,
) (bool, error) {
	r.onceTemp.Do(func() { close(r.setTempStarted) placeholder)
	<-ctx.Done()
	return false, ctx.Err()
placeholder

type grokCredentialBlockingCache struct {
	GrokTokenCache
	deleteStarted chan struct{placeholder
	releaseDelete chan struct{placeholder
	once          sync.Once
	mu            sync.Mutex
	deleted       bool
placeholder

type grokCredentialSequencedRepo struct {
	*tokenRefreshAccountRepo
	mu      sync.Mutex
	latest  *Account
	getCall int
placeholder

type grokCredentialRereadFailureRepo struct {
	*tokenRefreshAccountRepo
	account *Account
	err     error
placeholder

type grokCredentialCountingRefresher struct {
	refreshCalls int
placeholder

func (r *grokCredentialCountingRefresher) CacheKey(account *Account) string {
	return GrokTokenCacheKey(account)
placeholder

func (r *grokCredentialCountingRefresher) CanRefresh(*Account) bool { return true placeholder

func (r *grokCredentialCountingRefresher) NeedsRefresh(*Account, time.Duration) bool { return true placeholder

func (r *grokCredentialCountingRefresher) Refresh(context.Context, *Account) (map[string]any, error) {
	r.refreshCalls++
	return map[string]any{"access_token": "must-not-be-used"placeholder, nil
placeholder

func (r *grokCredentialRereadFailureRepo) GetByID(context.Context, int64) (*Account, error) {
	return r.account, r.err
placeholder

func (r *grokCredentialSequencedRepo) GetByID(ctx context.Context, id int64) (*Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.getCall++
	if r.getCall > 1 && r.latest != nil {
		return r.latest, nil
placeholder
	return r.tokenRefreshAccountRepo.GetByID(ctx, id)
placeholder

func (c *grokCredentialBlockingCache) DeleteAccessToken(ctx context.Context, _ string) error {
	c.once.Do(func() { close(c.deleteStarted) placeholder)
	select {
	case <-c.releaseDelete:
		c.mu.Lock()
		c.deleted = true
		c.mu.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
placeholder
placeholder

func (c *grokCredentialBlockingCache) wasDeleted() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.deleted
placeholder

func TestUpstreamFailoverErrorNextAccountActionPreservesLegacyRetry(t *testing.T) {
	t.Parallel()

	require.True(t, (&UpstreamFailoverError{placeholder).ShouldRetryNextAccount())
	require.True(t, (&UpstreamFailoverError{NextAccountAction: NextAccountRetryplaceholder).ShouldRetryNextAccount())
	require.False(t, (&UpstreamFailoverError{NextAccountAction: NextAccountStopplaceholder).ShouldRetryNextAccount())
placeholder

func TestGetRequestCredentialMapsPermanentGrokOAuthFailureAndRedactsSecrets(t *testing.T) {
	gin.SetMode(gin.TestMode)
	account := expiredGrokOAuthAccountForCredentialTest(701)
	repo := &tokenRefreshAccountRepo{placeholder
	repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
	cache := &grokTokenCacheForProviderTest{lockResult: trueplaceholder
	provider := NewGrokTokenProvider(repo, cache)
	provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), &tokenRefresherStub{
		err: infraerrors.New(http.StatusBadGateway, "GROK_OAUTH_TOKEN_REFRESH_FAILED", "invalid_grant access_token=leaked-access refresh_token=leaked-refresh"),
placeholder)
	svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	token, kind, err := svc.getRequestCredential(context.Background(), c, account)
placeholder
	require.Empty(t, token)
	require.Empty(t, kind)

	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, GatewayFailureStageAccountAuth, failoverErr.Stage)
	require.Equal(t, GatewayFailureScopeAccount, failoverErr.Scope)
	require.Equal(t, GrokCredentialReasonRevoked, failoverErr.Reason)
	require.True(t, failoverErr.ShouldRetryNextAccount())
	require.Equal(t, 0, failoverErr.StatusCode)
	require.Equal(t, http.StatusServiceUnavailable, failoverErr.ClientStatusCode)
	require.NotContains(t, err.Error(), "leaked-access")
	require.NotContains(t, err.Error(), "leaked-refresh")

	require.Equal(t, 1, repo.setErrorCalls)
	require.Zero(t, repo.setTempUnschedCalls)
	require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
	require.Equal(t, []string{GrokTokenCacheKey(account)placeholder, cache.deletedKeys)
	require.NotContains(t, repo.lastErrorMessage, "leaked-access")
	require.NotContains(t, repo.lastErrorMessage, "leaked-refresh")

	rawEvents, ok := c.Get(OpsUpstreamErrorsKey)
	require.True(t, ok)
	events, ok := rawEvents.([]*OpsUpstreamErrorEvent)
	require.True(t, ok)
	require.Len(t, events, 1)
	require.Equal(t, string(GatewayFailureStageAccountAuth), events[0].Stage)
	require.Equal(t, string(GatewayFailureScopeAccount), events[0].Scope)
	require.Equal(t, string(GrokCredentialReasonRevoked), events[0].Reason)
	require.Zero(t, events[0].UpstreamStatusCode)
	require.NotContains(t, events[0].Message, "leaked-access")
	require.NotContains(t, events[0].Message, "leaked-refresh")
placeholder

func TestGetRequestCredentialPermanentMappingsPersistAndInvalidate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name        string
		prepare     func(*Account)
		refreshErr  error
		wantReason  GatewayFailureReason
		cachedToken string
placeholder{
		{
			name: "missing refresh credential",
			prepare: func(account *Account) {
				delete(account.Credentials, "refresh_token")
		placeholder,
			wantReason: GrokCredentialReasonMissing,
	placeholder,
		{
			name: "missing access credential",
			prepare: func(account *Account) {
				delete(account.Credentials, "access_token")
				account.Credentials["expires_at"] = time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339)
		placeholder,
			wantReason:  GrokCredentialReasonMissing,
			cachedToken: "stale-cached-access",
	placeholder,
		{
			name:       "explicit entitlement action required",
			prepare:    func(*Account) {placeholder,
			refreshErr: infraerrors.New(http.StatusForbidden, "GROK_OAUTH_ENTITLEMENT_DENIED", "access_denied"),
			wantReason: GrokCredentialReasonEntitlement,
	placeholder,
placeholder

	for index, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := expiredGrokOAuthAccountForCredentialTest(int64(720 + index))
			account.Status = StatusActive
			account.Schedulable = true
			tt.prepare(account)
			baseRepo := &tokenRefreshAccountRepo{placeholder
			baseRepo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
			repo := &grokCredentialPersistingRepo{tokenRefreshAccountRepo: baseRepoplaceholder
			cache := &grokTokenCacheForProviderTest{lockResult: true, token: tt.cachedTokenplaceholder
			provider := NewGrokTokenProvider(repo, cache)
			provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), &tokenRefresherStub{err: tt.refreshErrplaceholder)
			svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
			c, _ := gin.CreateTestContext(httptest.NewRecorder())

			_, _, err := svc.getRequestCredential(context.Background(), c, account)
			var failoverErr *UpstreamFailoverError
			require.ErrorAs(t, err, &failoverErr)
			require.Equal(t, tt.wantReason, failoverErr.Reason)
			require.Equal(t, GatewayFailureScopeAccount, failoverErr.Scope)
			require.Equal(t, 1, baseRepo.setErrorCalls)
			require.Equal(t, StatusError, account.Status)
			require.False(t, account.Schedulable)
			require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
			require.Equal(t, []string{GrokTokenCacheKey(account)placeholder, cache.deletedKeys)
	placeholder)
placeholder
placeholder

func TestGetRequestCredentialMissingAccessNeverRefreshesAndPermanentlyFailsOver(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name      string
		expiresAt *time.Time
placeholder{
		{name: "expiry missing"placeholder,
		{name: "expired", expiresAt: func() *time.Time { value := time.Now().Add(-time.Minute); return &value placeholder()placeholder,
		{name: "near expiry", expiresAt: func() *time.Time { value := time.Now().Add(30 * time.Minute); return &value placeholder()placeholder,
placeholder

	for index, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := expiredGrokOAuthAccountForCredentialTest(int64(760 + index))
			account.Schedulable = true
			delete(account.Credentials, "access_token")
			if tt.expiresAt == nil {
				delete(account.Credentials, "expires_at")
		placeholder else {
				account.Credentials["expires_at"] = tt.expiresAt.UTC().Format(time.RFC3339)
		placeholder
			baseRepo := &tokenRefreshAccountRepo{placeholder
			baseRepo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
			repo := &grokCredentialPersistingRepo{tokenRefreshAccountRepo: baseRepoplaceholder
			cache := &grokTokenCacheForProviderTest{lockResult: true, token: "stale-cache-must-not-win"placeholder
			refresher := &grokCredentialCountingRefresher{placeholder
			provider := NewGrokTokenProvider(repo, cache)
			provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), refresher)
			svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
			c, _ := gin.CreateTestContext(httptest.NewRecorder())

			token, kind, err := svc.getRequestCredential(context.Background(), c, account)
			var failoverErr *UpstreamFailoverError
			require.ErrorAs(t, err, &failoverErr)
			require.Empty(t, token)
			require.Empty(t, kind)
			require.Equal(t, GrokCredentialReasonMissing, failoverErr.Reason)
			require.Equal(t, NextAccountRetry, failoverErr.NextAccountAction)
			require.Zero(t, refresher.refreshCalls, "structurally missing access credentials must not reach the token endpoint")
			require.Equal(t, 1, baseRepo.setErrorCalls)
			require.Equal(t, StatusError, account.Status)
			require.False(t, account.Schedulable)
			require.Equal(t, []string{GrokTokenCacheKey(account)placeholder, cache.deletedKeys)
			require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
	placeholder)
placeholder
placeholder

func TestGetRequestCredentialWarmCachedAccessWithMissingRefreshPermanentlyFailsOver(t *testing.T) {
	account := expiredGrokOAuthAccountForCredentialTest(764)
	account.Credentials["access_token"] = "valid-access"
	account.Credentials["expires_at"] = time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339)
	delete(account.Credentials, "refresh_token")
	baseRepo := &tokenRefreshAccountRepo{placeholder
	baseRepo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
	repo := &grokCredentialPersistingRepo{tokenRefreshAccountRepo: baseRepoplaceholder
	cache := &grokTokenCacheForProviderTest{lockResult: true, token: "valid-access"placeholder
	refresher := &grokCredentialCountingRefresher{placeholder
	provider := NewGrokTokenProvider(repo, cache)
	provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), refresher)
	svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	_, _, err := svc.getRequestCredential(context.Background(), c, account)

	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, GrokCredentialReasonMissing, failoverErr.Reason)
	require.Equal(t, NextAccountRetry, failoverErr.NextAccountAction)
	require.Zero(t, refresher.refreshCalls)
	require.Equal(t, 1, baseRepo.setErrorCalls)
	require.Equal(t, []string{GrokTokenCacheKey(account)placeholder, cache.deletedKeys)
	require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestGetRequestCredentialMapsTransientAndProviderFailuresSeparately(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("account transient temporarily unschedules", func(t *testing.T) {
		account := expiredGrokOAuthAccountForCredentialTest(702)
		repo := &tokenRefreshAccountRepo{placeholder
		repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
		cache := &grokTokenCacheForProviderTest{lockResult: trueplaceholder
		provider := NewGrokTokenProvider(repo, cache)
		provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), &tokenRefresherStub{err: errors.New("temporary refresh transport failure")placeholder)
		svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		_, _, err := svc.getRequestCredential(context.Background(), c, account)
		var failoverErr *UpstreamFailoverError
		require.ErrorAs(t, err, &failoverErr)
		require.Equal(t, GatewayFailureScopeAccount, failoverErr.Scope)
		require.Equal(t, GrokCredentialReasonRefreshTransient, failoverErr.Reason)
		require.True(t, failoverErr.ShouldRetryNextAccount())
		require.Zero(t, repo.setErrorCalls)
		require.Equal(t, 1, repo.setTempUnschedCalls)
		require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder)

	t.Run("shared provider configuration stops without mutation", func(t *testing.T) {
		account := expiredGrokOAuthAccountForCredentialTest(703)
		repo := &tokenRefreshAccountRepo{placeholder
		repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
		provider := NewGrokTokenProvider(repo, nil)
		svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		_, _, err := svc.getRequestCredential(context.Background(), c, account)
		var failoverErr *UpstreamFailoverError
		require.ErrorAs(t, err, &failoverErr)
		require.Equal(t, GatewayFailureScopeProvider, failoverErr.Scope)
		require.Equal(t, NextAccountStop, failoverErr.NextAccountAction)
		require.Zero(t, repo.setErrorCalls)
		require.Zero(t, repo.setTempUnschedCalls)
		require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder)

	t.Run("account reread failures preserve shared versus missing-row scope", func(t *testing.T) {
		for _, tt := range []struct {
			name    string
			account *Account
			err     error
	placeholder{
			{name: "repository error", err: errors.New("database temporarily unavailable")placeholder,
			{name: "missing row"placeholder,
	placeholder {
			t.Run(tt.name, func(t *testing.T) {
				account := expiredGrokOAuthAccountForCredentialTest(712)
				baseRepo := &tokenRefreshAccountRepo{placeholder
				baseRepo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
				repo := &grokCredentialRereadFailureRepo{tokenRefreshAccountRepo: baseRepo, account: tt.account, err: tt.errplaceholder
				cache := &grokTokenCacheForProviderTest{lockResult: trueplaceholder
				provider := NewGrokTokenProvider(repo, cache)
				provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), NewGrokTokenRefresher(nil))
				svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
				c, _ := gin.CreateTestContext(httptest.NewRecorder())

				_, _, err := svc.getRequestCredential(context.Background(), c, account)
				var failoverErr *UpstreamFailoverError
				require.ErrorAs(t, err, &failoverErr)
				if tt.err != nil {
					require.Equal(t, GatewayFailureScopeProvider, failoverErr.Scope)
					require.Equal(t, GrokCredentialReasonProviderDown, failoverErr.Reason)
					require.Equal(t, NextAccountStop, failoverErr.NextAccountAction)
			placeholder else {
					require.Equal(t, GatewayFailureScopeAccount, failoverErr.Scope)
					require.Equal(t, GrokCredentialReasonAccountChanged, failoverErr.Reason)
					require.Equal(t, NextAccountRetry, failoverErr.NextAccountAction)
			placeholder
				require.Zero(t, baseRepo.setErrorCalls)
				require.Zero(t, baseRepo.setTempUnschedCalls)
				require.Empty(t, cache.deletedKeys)
				require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
		placeholder)
	placeholder
placeholder)

	t.Run("fresh account eligibility changes retry without mutating stale state", func(t *testing.T) {
		for _, tt := range []struct {
			name   string
			mutate func(*Account)
	placeholder{
			{
				name: "account disabled",
				mutate: func(account *Account) {
					account.Status = StatusDisabled
			placeholder,
		placeholder,
			{
				name: "account converted",
				mutate: func(account *Account) {
					account.Type = AccountTypeUpstream
			placeholder,
		placeholder,
			{
				name: "account manually unschedulable",
				mutate: func(account *Account) {
					account.Schedulable = false
			placeholder,
		placeholder,
			{
				name: "account temporarily unschedulable",
				mutate: func(account *Account) {
					until := time.Now().Add(time.Minute)
					account.TempUnschedulableUntil = &until
			placeholder,
		placeholder,
	placeholder {
			t.Run(tt.name, func(t *testing.T) {
				staleAccount := expiredGrokOAuthAccountForCredentialTest(713)
				freshAccount := *staleAccount
				freshAccount.Credentials = shallowCopyMap(staleAccount.Credentials)
				tt.mutate(&freshAccount)
				repo := &tokenRefreshAccountRepo{placeholder
				repo.accountsByID = map[int64]*Account{staleAccount.ID: &freshAccountplaceholder
				cache := &grokTokenCacheForProviderTest{lockResult: trueplaceholder
				provider := NewGrokTokenProvider(repo, cache)
				provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), NewGrokTokenRefresher(nil))
				svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
				c, _ := gin.CreateTestContext(httptest.NewRecorder())

				_, _, err := svc.getRequestCredential(context.Background(), c, staleAccount)
				var failoverErr *UpstreamFailoverError
				require.ErrorAs(t, err, &failoverErr)
				require.Equal(t, GatewayFailureScopeAccount, failoverErr.Scope)
				require.Equal(t, GrokCredentialReasonAccountChanged, failoverErr.Reason)
				require.Equal(t, NextAccountRetry, failoverErr.NextAccountAction)
				require.Zero(t, repo.setErrorCalls)
				require.Zero(t, repo.setTempUnschedCalls)
				require.Empty(t, cache.deletedKeys)
				require.False(t, svc.isOpenAIAccountRuntimeBlocked(staleAccount))
		placeholder)
	placeholder
placeholder)

	t.Run("fresh missing refresh credential permanently blocks the account", func(t *testing.T) {
		staleAccount := expiredGrokOAuthAccountForCredentialTest(714)
		staleAccount.Schedulable = true
		freshAccount := *staleAccount
		freshAccount.Credentials = shallowCopyMap(staleAccount.Credentials)
		delete(freshAccount.Credentials, "refresh_token")
		baseRepo := &tokenRefreshAccountRepo{placeholder
		baseRepo.accountsByID = map[int64]*Account{staleAccount.ID: &freshAccountplaceholder
		repo := &grokCredentialPersistingRepo{tokenRefreshAccountRepo: baseRepoplaceholder
		cache := &grokTokenCacheForProviderTest{lockResult: trueplaceholder
		provider := NewGrokTokenProvider(repo, cache)
		provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), NewGrokTokenRefresher(nil))
		svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		_, _, err := svc.getRequestCredential(context.Background(), c, staleAccount)
		var failoverErr *UpstreamFailoverError
		require.ErrorAs(t, err, &failoverErr)
		require.Equal(t, GatewayFailureScopeAccount, failoverErr.Scope)
		require.Equal(t, GrokCredentialReasonMissing, failoverErr.Reason)
		require.Equal(t, NextAccountRetry, failoverErr.NextAccountAction)
		require.Equal(t, 1, baseRepo.setErrorCalls)
		require.Zero(t, baseRepo.setTempUnschedCalls)
		require.Equal(t, StatusError, freshAccount.Status)
		require.False(t, freshAccount.Schedulable)
		require.Equal(t, []string{GrokTokenCacheKey(staleAccount)placeholder, cache.deletedKeys)
		require.True(t, svc.isOpenAIAccountRuntimeBlocked(staleAccount))
placeholder)

	t.Run("refresh added after locked structural failure wins conditional mutation", func(t *testing.T) {
		staleAccount := expiredGrokOAuthAccountForCredentialTest(717)
		freshAccount := *staleAccount
		freshAccount.Credentials = shallowCopyMap(staleAccount.Credentials)
		delete(freshAccount.Credentials, "refresh_token")
		repo := &tokenRefreshAccountRepo{placeholder
		repo.accountsByID = map[int64]*Account{staleAccount.ID: &freshAccountplaceholder
		repo.beforeConditionalState = func() {
			repaired := freshAccount
			repaired.Credentials = shallowCopyMap(freshAccount.Credentials)
			repaired.Credentials["refresh_token"] = "repaired-refresh-token"
			repaired.Credentials["expires_at"] = time.Now().Add(time.Hour).UTC().Format(time.RFC3339)
			repo.accountsByID[staleAccount.ID] = &repaired
	placeholder
		cache := &grokTokenCacheForProviderTest{lockResult: trueplaceholder
		provider := NewGrokTokenProvider(repo, cache)
		provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), NewGrokTokenRefresher(nil))
		svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		token, kind, err := svc.getRequestCredential(context.Background(), c, staleAccount)

	placeholder
		require.Equal(t, "expired-access-token", token)
		require.Equal(t, "oauth", kind)
		require.Zero(t, repo.setErrorCalls)
		require.Zero(t, repo.setTempUnschedCalls)
		require.Empty(t, cache.deletedKeys)
		require.False(t, svc.isOpenAIAccountRuntimeBlocked(staleAccount))
placeholder)

	t.Run("expiry-only repair wins full credential fingerprint CAS", func(t *testing.T) {
		account := expiredGrokOAuthAccountForCredentialTest(718)
		repo := &tokenRefreshAccountRepo{placeholder
		repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
		repo.beforeConditionalState = func() {
			repaired := *account
			repaired.Credentials = shallowCopyMap(account.Credentials)
			repaired.Credentials["expires_at"] = time.Now().Add(time.Hour).UTC().Format(time.RFC3339)
			repo.accountsByID[account.ID] = &repaired
	placeholder
		cache := &grokTokenCacheForProviderTest{lockResult: trueplaceholder
		provider := NewGrokTokenProvider(repo, cache)
		provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), &tokenRefresherStub{
			err: infraerrors.New(http.StatusBadGateway, "GROK_OAUTH_TOKEN_REFRESH_FAILED", "invalid_grant"),
	placeholder)
		svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		token, kind, err := svc.getRequestCredential(context.Background(), c, account)

	placeholder
		require.Equal(t, "expired-access-token", token)
		require.Equal(t, "oauth", kind)
		require.Zero(t, repo.setErrorCalls)
		require.Empty(t, cache.deletedKeys)
		require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder)

	t.Run("generic token endpoint 403 stops as shared provider failure", func(t *testing.T) {
		account := expiredGrokOAuthAccountForCredentialTest(708)
		repo := &tokenRefreshAccountRepo{placeholder
		repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
		cache := &grokTokenCacheForProviderTest{lockResult: trueplaceholder
		provider := NewGrokTokenProvider(repo, cache)
		provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), &tokenRefresherStub{
			err: infraerrors.New(http.StatusBadGateway, "GROK_OAUTH_TOKEN_REFRESH_FAILED", "token refresh failed: status 403, body: forbidden"),
	placeholder)
		svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		_, _, err := svc.getRequestCredential(context.Background(), c, account)
		var failoverErr *UpstreamFailoverError
		require.ErrorAs(t, err, &failoverErr)
		require.Equal(t, GatewayFailureScopeProvider, failoverErr.Scope)
		require.Equal(t, GrokCredentialReasonProviderDown, failoverErr.Reason)
		require.Equal(t, NextAccountStop, failoverErr.NextAccountAction)
		require.Zero(t, repo.setErrorCalls)
		require.Zero(t, repo.setTempUnschedCalls)
		require.Empty(t, cache.deletedKeys)
		require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder)

	t.Run("account proxy generic 403 remains bounded account transient", func(t *testing.T) {
		account := expiredGrokOAuthAccountForCredentialTest(711)
		proxyID := int64(43)
		account.ProxyID = &proxyID
		account.Proxy = &Proxy{placeholder
		repo := &tokenRefreshAccountRepo{placeholder
		repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
		cache := &grokTokenCacheForProviderTest{lockResult: trueplaceholder
		provider := NewGrokTokenProvider(repo, cache)
		provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), &tokenRefresherStub{
			err: infraerrors.New(http.StatusBadGateway, "GROK_OAUTH_TOKEN_REFRESH_FAILED", "token refresh failed: status 403, body: forbidden"),
	placeholder)
		svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		_, _, err := svc.getRequestCredential(context.Background(), c, account)
		var failoverErr *UpstreamFailoverError
		require.ErrorAs(t, err, &failoverErr)
		require.Equal(t, GatewayFailureScopeAccount, failoverErr.Scope)
		require.Equal(t, GrokCredentialReasonRefreshTransient, failoverErr.Reason)
		require.Equal(t, NextAccountRetry, failoverErr.NextAccountAction)
		require.Zero(t, repo.setErrorCalls)
		require.Equal(t, 1, repo.setTempUnschedCalls)
		require.Empty(t, cache.deletedKeys)
		require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder)

	t.Run("proxy repository read failure stops without account mutation", func(t *testing.T) {
		account := expiredGrokOAuthAccountForCredentialTest(709)
		proxyID := int64(41)
		account.ProxyID = &proxyID
		account.Proxy = &Proxy{placeholder
		repo := &tokenRefreshAccountRepo{placeholder
		repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
		cache := &grokTokenCacheForProviderTest{lockResult: trueplaceholder
		oauthSvc := NewGrokOAuthService(&grokCredentialProxyRepoStub{err: errors.New("database temporarily unavailable")placeholder, &grokOAuthClientStub{placeholder)
		defer oauthSvc.Stop()
		provider := NewGrokTokenProvider(repo, cache)
		provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), NewGrokTokenRefresher(oauthSvc))
		svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		_, _, err := svc.getRequestCredential(context.Background(), c, account)
		var failoverErr *UpstreamFailoverError
		require.ErrorAs(t, err, &failoverErr)
		require.Equal(t, GatewayFailureScopeProvider, failoverErr.Scope)
		require.Equal(t, GrokCredentialReasonProviderDown, failoverErr.Reason)
		require.Equal(t, NextAccountStop, failoverErr.NextAccountAction)
		require.Zero(t, repo.setErrorCalls)
		require.Zero(t, repo.setTempUnschedCalls)
		require.Empty(t, cache.deletedKeys)
		require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder)

	t.Run("structurally missing configured proxy permanently blocks only that account", func(t *testing.T) {
		account := expiredGrokOAuthAccountForCredentialTest(710)
		account.Status = StatusActive
		account.Schedulable = true
		proxyID := int64(42)
		account.ProxyID = &proxyID
		baseRepo := &tokenRefreshAccountRepo{placeholder
		baseRepo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
		repo := &grokCredentialPersistingRepo{tokenRefreshAccountRepo: baseRepoplaceholder
		cache := &grokTokenCacheForProviderTest{lockResult: trueplaceholder
		oauthSvc := NewGrokOAuthService(&grokCredentialProxyRepoStub{err: ErrProxyNotFoundplaceholder, &grokOAuthClientStub{placeholder)
		defer oauthSvc.Stop()
		provider := NewGrokTokenProvider(repo, cache)
		provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), NewGrokTokenRefresher(oauthSvc))
		svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		_, _, err := svc.getRequestCredential(context.Background(), c, account)
		var failoverErr *UpstreamFailoverError
		require.ErrorAs(t, err, &failoverErr)
		require.Equal(t, GatewayFailureScopeAccount, failoverErr.Scope)
		require.Equal(t, GrokCredentialReasonProxyInvalid, failoverErr.Reason)
		require.Equal(t, NextAccountRetry, failoverErr.NextAccountAction)
		require.Equal(t, 1, baseRepo.setErrorCalls)
		require.Equal(t, StatusError, account.Status)
		require.False(t, account.Schedulable)
		require.Equal(t, []string{GrokTokenCacheKey(account)placeholder, cache.deletedKeys)
		require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder)
placeholder

func TestGetRequestCredentialRuntimeBlockWinsBeforeWarmTokenCache(t *testing.T) {
	account := expiredGrokOAuthAccountForCredentialTest(716)
	account.Credentials["access_token"] = "valid-access"
	account.Credentials["expires_at"] = time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339)
	repo := &tokenRefreshAccountRepo{placeholder
	repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
	cache := &grokTokenCacheForProviderTest{token: "valid-access"placeholder
	provider := NewGrokTokenProvider(repo, cache)
	svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
	svc.BlockAccountScheduling(account, time.Now().Add(time.Minute), "independent")
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	_, _, err := svc.getRequestCredential(context.Background(), c, account)

	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, GrokCredentialReasonAccountChanged, failoverErr.Reason)
	require.Equal(t, NextAccountRetry, failoverErr.NextAccountAction)
	require.Zero(t, cache.getCalls)
	require.Zero(t, repo.setErrorCalls)
	require.Zero(t, repo.setTempUnschedCalls)
placeholder

func TestGetRequestCredentialWarmCachedAccessWithMissingConfiguredProxyPermanentlyFailsOver(t *testing.T) {
	account := expiredGrokOAuthAccountForCredentialTest(715)
	account.Credentials["access_token"] = "valid-access"
	account.Credentials["expires_at"] = time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339)
	proxyID := int64(44)
	account.ProxyID = &proxyID
	account.Proxy = nil
	baseRepo := &tokenRefreshAccountRepo{placeholder
	baseRepo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
	repo := &grokCredentialPersistingRepo{tokenRefreshAccountRepo: baseRepoplaceholder
	cache := &grokTokenCacheForProviderTest{lockResult: true, token: "valid-access"placeholder
	refresher := &grokCredentialCountingRefresher{placeholder
	provider := NewGrokTokenProvider(repo, cache)
	provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), refresher)
	svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	_, _, err := svc.getRequestCredential(context.Background(), c, account)

	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, GrokCredentialReasonProxyInvalid, failoverErr.Reason)
	require.Equal(t, NextAccountRetry, failoverErr.NextAccountAction)
	require.Zero(t, refresher.refreshCalls)
	require.Equal(t, 1, baseRepo.setErrorCalls)
	require.Equal(t, []string{GrokTokenCacheKey(account)placeholder, cache.deletedKeys)
	require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestGetRequestCredentialCancellationAndBudgetDoNotMutateAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	account := expiredGrokOAuthAccountForCredentialTest(704)
	repo := &tokenRefreshAccountRepo{placeholder
	repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
	provider := NewGrokTokenProvider(repo, nil)
	svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder

	t.Run("parent cancellation is returned directly", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		_, _, err := svc.getRequestCredential(ctx, c, account)
		require.ErrorIs(t, err, context.Canceled)
		var failoverErr *UpstreamFailoverError
		require.False(t, errors.As(err, &failoverErr))
placeholder)

	t.Run("request credential budget stops safely", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set(grokCredentialFailoverDeadlineKey, time.Now().Add(-time.Second))

		_, _, err := svc.getRequestCredential(context.Background(), c, account)
		var failoverErr *UpstreamFailoverError
		require.ErrorAs(t, err, &failoverErr)
		require.Equal(t, GatewayFailureScopeRequest, failoverErr.Scope)
		require.Equal(t, GrokCredentialReasonFailoverTimeout, failoverErr.Reason)
		require.False(t, failoverErr.ShouldRetryNextAccount())
placeholder)

	require.Zero(t, repo.setErrorCalls)
	require.Zero(t, repo.setTempUnschedCalls)
	require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestGetRequestCredentialStateMutationFailureStopsAndKeepsRuntimeBlock(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name            string
		refreshErr      error
		configure       func(*tokenRefreshAccountRepo, *grokTokenCacheForProviderTest)
		wantSetError    int
		wantSetTemp     int
		wantCacheDelete int
placeholder{
		{
			name:       "permanent state persistence",
			refreshErr: infraerrors.New(http.StatusBadGateway, "GROK_OAUTH_TOKEN_REFRESH_FAILED", "invalid_grant"),
			configure: func(repo *tokenRefreshAccountRepo, _ *grokTokenCacheForProviderTest) {
				repo.setErrorErr = errors.New("database write failed")
		placeholder,
			wantSetError: 1,
	placeholder,
		{
			name:       "transient state persistence",
			refreshErr: errors.New("temporary refresh transport failure"),
			configure: func(repo *tokenRefreshAccountRepo, _ *grokTokenCacheForProviderTest) {
				repo.setTempUnschedErr = errors.New("database write failed")
		placeholder,
			wantSetTemp: 1,
	placeholder,
		{
			name:       "permanent token cache invalidation",
			refreshErr: infraerrors.New(http.StatusBadGateway, "GROK_OAUTH_TOKEN_REFRESH_FAILED", "invalid_grant"),
			configure: func(_ *tokenRefreshAccountRepo, cache *grokTokenCacheForProviderTest) {
				cache.deleteErr = errors.New("cache delete failed")
		placeholder,
			wantSetError:    1,
			wantCacheDelete: 1,
	placeholder,
placeholder

	for index, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := expiredGrokOAuthAccountForCredentialTest(int64(740 + index))
			repo := &tokenRefreshAccountRepo{placeholder
			repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
			cache := &grokTokenCacheForProviderTest{lockResult: trueplaceholder
			tt.configure(repo, cache)
			provider := NewGrokTokenProvider(repo, cache)
			provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), &tokenRefresherStub{err: tt.refreshErrplaceholder)
			svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
			c, _ := gin.CreateTestContext(httptest.NewRecorder())

			_, _, err := svc.getRequestCredential(context.Background(), c, account)
			var failoverErr *UpstreamFailoverError
			require.ErrorAs(t, err, &failoverErr)
			require.Equal(t, GatewayFailureScopeProvider, failoverErr.Scope)
			require.Equal(t, GrokCredentialReasonStateUpdate, failoverErr.Reason)
			require.Equal(t, NextAccountStop, failoverErr.NextAccountAction)
			require.Equal(t, tt.wantSetError, repo.setErrorCalls)
			require.Equal(t, tt.wantSetTemp, repo.setTempUnschedCalls)
			require.Len(t, cache.deletedKeys, tt.wantCacheDelete)
			require.True(t, svc.isOpenAIAccountRuntimeBlocked(account), "failed mutation must retain the immediate local block")
	placeholder)
placeholder
placeholder

func TestGrokCredentialMutationBoundariesHonorParentCancellation(t *testing.T) {
	t.Run("blocked SetError cancellation prevents cache and runtime mutation", func(t *testing.T) {
		account := expiredGrokOAuthAccountForCredentialTest(730)
		baseRepo := &tokenRefreshAccountRepo{placeholder
		baseRepo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
		repo := &grokCredentialBlockingRepo{
			tokenRefreshAccountRepo: baseRepo,
			setErrorStarted:         make(chan struct{placeholder),
			setTempStarted:          make(chan struct{placeholder),
	placeholder
		cache := &grokTokenCacheForProviderTest{placeholder
		svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: NewGrokTokenProvider(repo, cache)placeholder
		ctx, cancel := context.WithCancel(context.Background())
		result := make(chan error, 1)
		go func() {
			_, err := svc.applyGrokCredentialAccountFailure(ctx, account, grokCredentialFailureClass{
				reason: GrokCredentialReasonRevoked, permanent: true,
		placeholder)
			result <- err
	placeholder()
		<-repo.setErrorStarted
		require.True(t, svc.isOpenAIAccountRuntimeBlocked(account), "runtime block must precede persistent SetError")
		cancel()

		require.ErrorIs(t, <-result, context.Canceled)
		require.Zero(t, baseRepo.setErrorCalls)
		require.Empty(t, cache.deletedKeys)
		require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder)

	t.Run("blocked temporary unschedule cancellation prevents runtime mutation", func(t *testing.T) {
		account := expiredGrokOAuthAccountForCredentialTest(731)
		baseRepo := &tokenRefreshAccountRepo{placeholder
		baseRepo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
		repo := &grokCredentialBlockingRepo{
			tokenRefreshAccountRepo: baseRepo,
			setErrorStarted:         make(chan struct{placeholder),
			setTempStarted:          make(chan struct{placeholder),
	placeholder
		svc := &OpenAIGatewayService{accountRepo: repoplaceholder
		ctx, cancel := context.WithCancel(context.Background())
		result := make(chan error, 1)
		go func() {
			_, err := svc.applyGrokCredentialAccountFailure(ctx, account, grokCredentialFailureClass{
				reason: GrokCredentialReasonRefreshTransient, transient: true,
		placeholder)
			result <- err
	placeholder()
		<-repo.setTempStarted
		require.True(t, svc.isOpenAIAccountRuntimeBlocked(account), "runtime block must precede temporary unscheduling")
		cancel()

		require.ErrorIs(t, <-result, context.Canceled)
		require.Zero(t, baseRepo.setTempUnschedCalls)
		require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder)

	t.Run("post-commit cancellation finishes cache cleanup and retains quarantine", func(t *testing.T) {
		account := expiredGrokOAuthAccountForCredentialTest(732)
		repo := &tokenRefreshAccountRepo{placeholder
		repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
		cache := &grokCredentialBlockingCache{deleteStarted: make(chan struct{placeholder), releaseDelete: make(chan struct{placeholder)placeholder
		svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: NewGrokTokenProvider(repo, cache)placeholder
		ctx, cancel := context.WithCancel(context.Background())
		result := make(chan error, 1)
		go func() {
			_, err := svc.applyGrokCredentialAccountFailure(ctx, account, grokCredentialFailureClass{
				reason: GrokCredentialReasonRevoked, permanent: true,
		placeholder)
			result <- err
	placeholder()
		<-cache.deleteStarted
		require.True(t, svc.isOpenAIAccountRuntimeBlocked(account), "runtime block must precede cache invalidation")
		cancel()
		close(cache.releaseDelete)

		require.ErrorIs(t, <-result, context.Canceled)
		require.Equal(t, 1, repo.setErrorCalls)
		require.True(t, cache.wasDeleted())
		require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder)
placeholder

func TestGrokCredentialMutationLockWaitHonorsCredentialBudget(t *testing.T) {
	account := expiredGrokOAuthAccountForCredentialTest(735)
	repo := &tokenRefreshAccountRepo{placeholder
	repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
	svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: NewGrokTokenProvider(repo, &grokTokenCacheForProviderTest{placeholder)placeholder
	mutationLock := svc.grokCredentialMutationLock(account.ID)
	require.NoError(t, mutationLock.Lock(context.Background()))
	defer mutationLock.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	startedAt := time.Now()
	token, err := svc.applyGrokCredentialAccountFailure(ctx, account, grokCredentialFailureClass{
		reason: GrokCredentialReasonRevoked, permanent: true,
placeholder)

	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Empty(t, token)
	require.Less(t, time.Since(startedAt), 500*time.Millisecond)
	require.Zero(t, repo.setErrorCalls)
	require.Zero(t, repo.setTempUnschedCalls)
	require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestGetRequestCredentialBudgetBoundsBlockedConditionalMutation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	account := expiredGrokOAuthAccountForCredentialTest(736)
	baseRepo := &tokenRefreshAccountRepo{placeholder
	baseRepo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
	repo := &grokCredentialBlockingRepo{
		tokenRefreshAccountRepo: baseRepo,
		setErrorStarted:         make(chan struct{placeholder),
		setTempStarted:          make(chan struct{placeholder),
placeholder
	cache := &grokTokenCacheForProviderTest{lockResult: trueplaceholder
	provider := NewGrokTokenProvider(repo, cache)
	provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), &tokenRefresherStub{
		err: infraerrors.New(http.StatusBadGateway, "GROK_OAUTH_TOKEN_REFRESH_FAILED", "invalid_grant"),
placeholder)
	svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(grokCredentialFailoverDeadlineKey, time.Now().Add(40*time.Millisecond))

	startedAt := time.Now()
	token, kind, err := svc.getRequestCredential(context.Background(), c, account)

	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, GatewayFailureScopeRequest, failoverErr.Scope)
	require.Equal(t, GrokCredentialReasonFailoverTimeout, failoverErr.Reason)
	require.Equal(t, NextAccountStop, failoverErr.NextAccountAction)
	require.Empty(t, token)
	require.Empty(t, kind)
	require.Less(t, time.Since(startedAt), 500*time.Millisecond)
	require.Zero(t, baseRepo.setErrorCalls)
	require.Zero(t, baseRepo.setTempUnschedCalls)
	require.Empty(t, cache.deletedKeys)
	require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestGetRequestCredentialLockHeldTimeoutDoesNotQuarantineAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name       string
		buildRepo  func(*Account) AccountRepository
		wantScope  GatewayFailureScope
		wantReason GatewayFailureReason
		wantAction NextAccountAction
placeholder{
		{
			name: "authoritative row unchanged",
			buildRepo: func(account *Account) AccountRepository {
				repo := &tokenRefreshAccountRepo{placeholder
				repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
				return repo
		placeholder,
			wantScope:  GatewayFailureScopeAccount,
			wantReason: GrokCredentialReasonAccountChanged,
			wantAction: NextAccountRetry,
	placeholder,
		{
			name: "selected account was deleted",
			buildRepo: func(account *Account) AccountRepository {
				base := &tokenRefreshAccountRepo{placeholder
				base.accountsByID = map[int64]*Account{account.ID: accountplaceholder
				return &grokCredentialRereadFailureRepo{tokenRefreshAccountRepo: baseplaceholder
		placeholder,
			wantScope:  GatewayFailureScopeAccount,
			wantReason: GrokCredentialReasonAccountChanged,
			wantAction: NextAccountRetry,
	placeholder,
		{
			name: "shared account store unavailable",
			buildRepo: func(account *Account) AccountRepository {
				base := &tokenRefreshAccountRepo{placeholder
				base.accountsByID = map[int64]*Account{account.ID: accountplaceholder
				return &grokCredentialRereadFailureRepo{tokenRefreshAccountRepo: base, err: errors.New("database unavailable")placeholder
		placeholder,
			wantScope:  GatewayFailureScopeProvider,
			wantReason: GrokCredentialReasonProviderDown,
			wantAction: NextAccountStop,
	placeholder,
placeholder

	for index, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := expiredGrokOAuthAccountForCredentialTest(int64(7400 + index))
			repo := tt.buildRepo(account)
			cache := &grokTokenCacheForProviderTest{lockResult: falseplaceholder
			provider := NewGrokTokenProvider(repo, cache)
			provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), &tokenRefresherStub{placeholder)
			svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
			c, _ := gin.CreateTestContext(httptest.NewRecorder())

			startedAt := time.Now()
			_, _, err := svc.getRequestCredential(context.Background(), c, account)

			var failoverErr *UpstreamFailoverError
			require.ErrorAs(t, err, &failoverErr)
			require.Equal(t, tt.wantScope, failoverErr.Scope)
			require.Equal(t, tt.wantReason, failoverErr.Reason)
			require.Equal(t, tt.wantAction, failoverErr.NextAccountAction)
			require.Less(t, time.Since(startedAt), 3*time.Second)
			switch countingRepo := repo.(type) {
			case *tokenRefreshAccountRepo:
				require.Zero(t, countingRepo.setErrorCalls)
				require.Zero(t, countingRepo.setTempUnschedCalls)
			case *grokCredentialRereadFailureRepo:
				require.Zero(t, countingRepo.tokenRefreshAccountRepo.setErrorCalls)
				require.Zero(t, countingRepo.tokenRefreshAccountRepo.setTempUnschedCalls)
		placeholder
			require.Empty(t, cache.deletedKeys)
			require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
	placeholder)
placeholder
placeholder

func TestGrokCredentialMutationCancellationAmbiguityConfirmsDurableCommit(t *testing.T) {
	tests := []struct {
		name      string
		class     grokCredentialFailureClass
		committed func(*Account) bool
placeholder{
		{
			name:  "permanent quarantine",
			class: grokCredentialFailureClass{reason: GrokCredentialReasonRevoked, permanent: trueplaceholder,
			committed: func(account *Account) bool {
				return account.Status == StatusError && !account.Schedulable && account.ErrorMessage == string(GrokCredentialReasonRevoked)
		placeholder,
	placeholder,
		{
			name:  "temporary quarantine",
			class: grokCredentialFailureClass{reason: GrokCredentialReasonRefreshTransient, transient: trueplaceholder,
			committed: func(account *Account) bool {
				return account.TempUnschedulableUntil != nil && account.TempUnschedulableReason == string(GrokCredentialReasonRefreshTransient)
		placeholder,
	placeholder,
placeholder

	for index, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := expiredGrokOAuthAccountForCredentialTest(int64(737 + index))
			baseRepo := &tokenRefreshAccountRepo{placeholder
			baseRepo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
			repo := &grokCredentialCommitThenCancelRepo{tokenRefreshAccountRepo: baseRepoplaceholder
			cache := &grokTokenCacheForProviderTest{placeholder
			svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: NewGrokTokenProvider(repo, cache)placeholder
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
			defer cancel()

			token, err := svc.applyGrokCredentialAccountFailure(ctx, account, tt.class)

			require.ErrorIs(t, err, context.DeadlineExceeded)
			require.Empty(t, token)
			require.True(t, tt.committed(account), "the detached confirmation must recognize the durable mutation")
			require.True(t, svc.isOpenAIAccountRuntimeBlocked(account), "a confirmed durable quarantine must retain its runtime block")
			if tt.class.permanent {
				require.Equal(t, []string{GrokTokenCacheKey(account)placeholder, cache.deletedKeys)
		placeholder else {
				require.Empty(t, cache.deletedKeys)
		placeholder
	placeholder)
placeholder
placeholder

func TestGrokCredentialInnerStateDeadlineAmbiguityConfirmsDurableCommit(t *testing.T) {
	tests := []struct {
		name  string
		class grokCredentialFailureClass
placeholder{
		{name: "permanent quarantine", class: grokCredentialFailureClass{reason: GrokCredentialReasonRevoked, permanent: trueplaceholderplaceholder,
		{name: "temporary quarantine", class: grokCredentialFailureClass{reason: GrokCredentialReasonRefreshTransient, transient: trueplaceholderplaceholder,
placeholder

	for index, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := expiredGrokOAuthAccountForCredentialTest(int64(7500 + index))
			baseRepo := &tokenRefreshAccountRepo{placeholder
			baseRepo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
			repo := &grokCredentialCommitThenCancelRepo{
				tokenRefreshAccountRepo: baseRepo,
				returnErr:               context.DeadlineExceeded,
		placeholder
			cache := &grokTokenCacheForProviderTest{placeholder
			svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: NewGrokTokenProvider(repo, cache)placeholder

			token, err := svc.applyGrokCredentialAccountFailure(context.Background(), account, tt.class)

			require.NoError(t, err, "the detached readback must resolve the inner timeout's commit ambiguity")
			require.Empty(t, token)
			require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
			if tt.class.permanent {
				require.Equal(t, StatusError, account.Status)
				require.False(t, account.Schedulable)
				require.Equal(t, []string{GrokTokenCacheKey(account)placeholder, cache.deletedKeys)
		placeholder else {
				require.NotNil(t, account.TempUnschedulableUntil)
				require.Equal(t, string(GrokCredentialReasonRefreshTransient), account.TempUnschedulableReason)
				require.Empty(t, cache.deletedKeys)
		placeholder
	placeholder)
placeholder
placeholder

func TestGrokCredentialUnconfirmedInnerStateDeadlineStopsAndRetainsSafetyBlock(t *testing.T) {
	tests := []struct {
		name  string
		class grokCredentialFailureClass
placeholder{
		{name: "permanent quarantine", class: grokCredentialFailureClass{reason: GrokCredentialReasonRevoked, permanent: trueplaceholderplaceholder,
		{name: "temporary quarantine", class: grokCredentialFailureClass{reason: GrokCredentialReasonRefreshTransient, transient: trueplaceholderplaceholder,
placeholder

	for index, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := expiredGrokOAuthAccountForCredentialTest(int64(7600 + index))
			baseRepo := &tokenRefreshAccountRepo{placeholder
			baseRepo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
			repo := &grokCredentialUncommittedDeadlineRepo{tokenRefreshAccountRepo: baseRepoplaceholder
			cache := &grokTokenCacheForProviderTest{placeholder
			svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: NewGrokTokenProvider(repo, cache)placeholder

			token, err := svc.applyGrokCredentialAccountFailure(context.Background(), account, tt.class)

			require.ErrorIs(t, err, errGrokCredentialStateUpdateFailed)
			require.Empty(t, token)
			require.True(t, svc.isOpenAIAccountRuntimeBlocked(account), "an unknown commit outcome must retain the local safety block")
			require.Equal(t, StatusActive, account.Status)
			require.True(t, account.Schedulable)
			require.Nil(t, account.TempUnschedulableUntil)
			require.Empty(t, cache.deletedKeys)
	placeholder)
placeholder
placeholder

func TestGrokCredentialRuntimeRollbackOwnership(t *testing.T) {
	account := expiredGrokOAuthAccountForCredentialTest(734)
	t.Run("later extending block survives", func(t *testing.T) {
		svc := &OpenAIGatewayService{placeholder
		until := time.Now().Add(time.Minute)
		rollbackFirst := svc.blockGrokCredentialRuntime(account, until, "first")

		secondInstalled := make(chan struct{placeholder)
		go func() {
			svc.BlockAccountScheduling(account, until.Add(time.Minute), "independent")
			close(secondInstalled)
	placeholder()
		<-secondInstalled
		rollbackFirst()

		require.True(t, svc.isOpenAIAccountRuntimeBlocked(account),
			"rollback owned by the first invocation must not remove a later extending block")
placeholder)

	t.Run("independent shorter block steals rollback ownership", func(t *testing.T) {
		svc := &OpenAIGatewayService{placeholder
		until := time.Now().Add(2 * time.Minute)
		rollbackFirst := svc.blockGrokCredentialRuntime(account, until, "first")
		svc.BlockAccountScheduling(account, until.Add(-time.Minute), "shorter-no-op")

		rollbackFirst()

		require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder)

	t.Run("serialized tentative rollbacks leave no block", func(t *testing.T) {
		svc := &OpenAIGatewayService{placeholder
		for i := 0; i < 2; i++ {
			mu := svc.grokCredentialMutationLock(account.ID)
			require.NoError(t, mu.Lock(context.Background()))
			rollback := svc.blockGrokCredentialRuntime(account, time.Now().Add(time.Minute), "tentative")
			rollback()
			mu.Unlock()
	placeholder
		require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder)
placeholder

func TestGetRequestCredentialAPIKeyBypassesOAuthFailureMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	account := &Account{
		ID:       705,
		Platform: PlatformGrok,
		Type:     AccountTypeAPIKey,
placeholder
			"api_key":  "third-party-key",
			"base_url": "https://grok.example.test/v1",
	placeholder,
placeholder
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	svc := &OpenAIGatewayService{placeholder

	token, kind, err := svc.getRequestCredential(context.Background(), c, account)
placeholder
	require.Equal(t, "third-party-key", token)
	require.Equal(t, "apikey", kind)
	_, hasEvents := c.Get(OpsUpstreamErrorsKey)
	require.False(t, hasEvents)
placeholder

func TestPermanentCredentialFailureDoesNotDisableConcurrentlyRefreshedAccount(t *testing.T) {
	account := expiredGrokOAuthAccountForCredentialTest(707)
	latest := *account
	latest.Credentials = shallowCopyMap(account.Credentials)
	latest.Credentials["access_token"] = "fresh-access-token"
	latest.Credentials["refresh_token"] = "rotated-refresh-token"
	latest.Credentials["expires_at"] = time.Now().Add(time.Hour).UTC().Format(time.RFC3339)
	latest.Credentials["_token_version"] = time.Now().UnixMilli()
	repo := &tokenRefreshAccountRepo{placeholder
	repo.accountsByID = map[int64]*Account{account.ID: &latestplaceholder
	cache := &grokTokenCacheForProviderTest{placeholder
	svc := &OpenAIGatewayService{
		accountRepo:       repo,
		grokTokenProvider: NewGrokTokenProvider(repo, cache),
placeholder

	_, mutationErr := svc.applyGrokCredentialAccountFailure(context.Background(), account, grokCredentialFailureClass{
		scope:     GatewayFailureScopeAccount,
		reason:    GrokCredentialReasonRevoked,
		action:    NextAccountRetry,
		permanent: true,
placeholder)

	require.NoError(t, mutationErr)
	require.Zero(t, repo.setErrorCalls)
	require.Empty(t, cache.deletedKeys)
	require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestCredentialFailureConditionalMutationLosesToConcurrentRefresh(t *testing.T) {
	for index, tt := range []struct {
		name       string
		refreshErr error
placeholder{
		{name: "permanent", refreshErr: infraerrors.New(http.StatusBadGateway, "GROK_OAUTH_TOKEN_REFRESH_FAILED", "invalid_grant")placeholder,
		{name: "transient", refreshErr: errors.New("temporary refresh transport failure")placeholder,
placeholder {
		t.Run(tt.name, func(t *testing.T) {
			account := expiredGrokOAuthAccountForCredentialTest(int64(770 + index))
			repo := &tokenRefreshAccountRepo{placeholder
			repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
			repo.beforeConditionalState = func() {
				fresh := *account
				fresh.Credentials = shallowCopyMap(account.Credentials)
				fresh.Credentials["access_token"] = "refresh-won-token"
				fresh.Credentials["refresh_token"] = "refresh-won-refresh"
				fresh.Credentials["expires_at"] = time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339)
				fresh.Credentials["_token_version"] = time.Now().UnixMilli()
				repo.accountsByID[account.ID] = &fresh
		placeholder
			cache := &grokTokenCacheForProviderTest{lockResult: trueplaceholder
			provider := NewGrokTokenProvider(repo, cache)
			provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), &tokenRefresherStub{err: tt.refreshErrplaceholder)
			svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
			c, _ := gin.CreateTestContext(httptest.NewRecorder())

			token, kind, err := svc.getRequestCredential(context.Background(), c, account)

		placeholder
			require.Equal(t, "refresh-won-token", token)
			require.Equal(t, "oauth", kind)
			require.Zero(t, repo.setErrorCalls)
			require.Zero(t, repo.setTempUnschedCalls)
			require.Empty(t, cache.deletedKeys)
			require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
	placeholder)
placeholder
placeholder

func TestCredentialFailureConditionalMutationLosesToConcurrentProxyRepair(t *testing.T) {
	for index, tt := range []struct {
		name       string
		refreshErr error
placeholder{
		{name: "permanent", refreshErr: infraerrors.New(http.StatusBadGateway, "GROK_OAUTH_TOKEN_REFRESH_FAILED", "invalid_grant")placeholder,
		{name: "transient", refreshErr: errors.New("temporary refresh transport failure")placeholder,
placeholder {
		t.Run(tt.name, func(t *testing.T) {
			account := expiredGrokOAuthAccountForCredentialTest(int64(780 + index))
			oldProxyID := int64(10)
			account.ProxyID = &oldProxyID
			account.Proxy = &Proxy{placeholder
			repo := &tokenRefreshAccountRepo{placeholder
			repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
			repo.beforeConditionalState = func() {
				fresh := *account
				fresh.Credentials = shallowCopyMap(account.Credentials)
				repairedProxyID := int64(11)
				fresh.ProxyID = &repairedProxyID
				repo.accountsByID[account.ID] = &fresh
		placeholder
			cache := &grokTokenCacheForProviderTest{lockResult: trueplaceholder
			provider := NewGrokTokenProvider(repo, cache)
			provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), &tokenRefresherStub{err: tt.refreshErrplaceholder)
			svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
			c, _ := gin.CreateTestContext(httptest.NewRecorder())

			_, _, err := svc.getRequestCredential(context.Background(), c, account)

			var failoverErr *UpstreamFailoverError
			require.ErrorAs(t, err, &failoverErr)
			require.Equal(t, GrokCredentialReasonAccountChanged, failoverErr.Reason)
			require.Equal(t, NextAccountRetry, failoverErr.NextAccountAction)
			require.Zero(t, repo.setErrorCalls)
			require.Zero(t, repo.setTempUnschedCalls)
			require.Empty(t, cache.deletedKeys)
			require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
	placeholder)
placeholder
placeholder

func TestCredentialFailureConditionalMutationLosesToSameIDProxyRestoration(t *testing.T) {
	account := expiredGrokOAuthAccountForCredentialTest(790)
	proxyID := int64(10)
	account.ProxyID = &proxyID
	account.Proxy = nil
	repo := &tokenRefreshAccountRepo{placeholder
	repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
	repo.beforeConditionalState = func() {
		account.Proxy = &Proxy{ID: proxyIDplaceholder
placeholder
	cache := &grokTokenCacheForProviderTest{lockResult: trueplaceholder
	provider := NewGrokTokenProvider(repo, cache)
	svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	_, _, err := svc.getRequestCredential(context.Background(), c, account)

	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, GrokCredentialReasonAccountChanged, failoverErr.Reason)
	require.Equal(t, NextAccountRetry, failoverErr.NextAccountAction)
	require.Zero(t, repo.setErrorCalls)
	require.Zero(t, repo.setTempUnschedCalls)
	require.Empty(t, cache.deletedKeys)
	require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestCredentialFailureConditionalMutationLosesToConcurrentUnschedulableState(t *testing.T) {
	future := time.Now().Add(time.Hour)
	states := []struct {
		name   string
		mutate func(*Account)
placeholder{
		{name: "admin schedulable false", mutate: func(account *Account) { account.Schedulable = false placeholderplaceholder,
		{name: "temporary cooldown", mutate: func(account *Account) { account.TempUnschedulableUntil = &future placeholderplaceholder,
		{name: "rate limit cooldown", mutate: func(account *Account) { account.RateLimitResetAt = &future placeholderplaceholder,
		{name: "overload cooldown", mutate: func(account *Account) { account.OverloadUntil = &future placeholderplaceholder,
placeholder
	classes := []struct {
		name  string
		class grokCredentialFailureClass
placeholder{
		{name: "permanent", class: grokCredentialFailureClass{reason: GrokCredentialReasonRevoked, permanent: trueplaceholderplaceholder,
		{name: "transient", class: grokCredentialFailureClass{reason: GrokCredentialReasonRefreshTransient, transient: trueplaceholderplaceholder,
placeholder

	for classIndex, classCase := range classes {
		for stateIndex, stateCase := range states {
			t.Run(classCase.name+"/"+stateCase.name, func(t *testing.T) {
				account := expiredGrokOAuthAccountForCredentialTest(int64(791 + classIndex*10 + stateIndex))
				repo := &tokenRefreshAccountRepo{placeholder
				repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
				repo.beforeConditionalState = func() { stateCase.mutate(account) placeholder
				svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: NewGrokTokenProvider(repo, &grokTokenCacheForProviderTest{placeholder)placeholder

				token, err := svc.applyGrokCredentialAccountFailure(context.Background(), account, classCase.class)

				require.ErrorIs(t, err, errOAuthRefreshAccountStateChanged)
				require.Empty(t, token)
				require.Zero(t, repo.setErrorCalls)
				require.Zero(t, repo.setTempUnschedCalls)
				require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
		placeholder)
	placeholder
placeholder
placeholder

func TestCredentialFailureCASMissDoesNotRecoverIneligibleLatestCredential(t *testing.T) {
	future := time.Now().Add(time.Hour)
	states := []struct {
		name               string
		mutate             func(*OpenAIGatewayService, *Account)
		wantRuntimeBlocked bool
placeholder{
		{name: "disabled", mutate: func(_ *OpenAIGatewayService, account *Account) { account.Status = StatusDisabled placeholderplaceholder,
		{name: "not schedulable", mutate: func(_ *OpenAIGatewayService, account *Account) { account.Schedulable = false placeholderplaceholder,
		{name: "temporarily unschedulable", mutate: func(_ *OpenAIGatewayService, account *Account) { account.TempUnschedulableUntil = &future placeholderplaceholder,
		{name: "rate limited", mutate: func(_ *OpenAIGatewayService, account *Account) { account.RateLimitResetAt = &future placeholderplaceholder,
		{name: "overloaded", mutate: func(_ *OpenAIGatewayService, account *Account) { account.OverloadUntil = &future placeholderplaceholder,
		{
			name: "independently runtime blocked",
			mutate: func(svc *OpenAIGatewayService, account *Account) {
				svc.BlockAccountScheduling(account, time.Now().Add(24*time.Hour), "independent")
		placeholder,
			wantRuntimeBlocked: true,
	placeholder,
placeholder
	classes := []struct {
		name  string
		class grokCredentialFailureClass
placeholder{
		{name: "permanent", class: grokCredentialFailureClass{reason: GrokCredentialReasonRevoked, permanent: trueplaceholderplaceholder,
		{name: "transient", class: grokCredentialFailureClass{reason: GrokCredentialReasonRefreshTransient, transient: trueplaceholderplaceholder,
placeholder

	for classIndex, classCase := range classes {
		for stateIndex, stateCase := range states {
			t.Run(classCase.name+"/"+stateCase.name, func(t *testing.T) {
				account := expiredGrokOAuthAccountForCredentialTest(int64(800 + classIndex*20 + stateIndex))
				repo := &tokenRefreshAccountRepo{placeholder
				repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
				svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: NewGrokTokenProvider(repo, &grokTokenCacheForProviderTest{placeholder)placeholder
				repo.beforeConditionalState = func() {
					latest := *account
					latest.Credentials = shallowCopyMap(account.Credentials)
					latest.Credentials["access_token"] = "fresh-but-ineligible-token"
					latest.Credentials["refresh_token"] = "fresh-but-ineligible-refresh"
					latest.Credentials["expires_at"] = time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339)
					latest.Credentials["_token_version"] = time.Now().UnixMilli()
					stateCase.mutate(svc, &latest)
					repo.accountsByID[account.ID] = &latest
			placeholder

				token, err := svc.applyGrokCredentialAccountFailure(context.Background(), account, classCase.class)

				require.ErrorIs(t, err, errOAuthRefreshAccountStateChanged)
				require.Empty(t, token)
				require.Zero(t, repo.setErrorCalls)
				require.Zero(t, repo.setTempUnschedCalls)
				require.Equal(t, stateCase.wantRuntimeBlocked, svc.isOpenAIAccountRuntimeBlocked(account))
		placeholder)
	placeholder
placeholder
placeholder

func TestGetRequestCredentialSharedCredentialPersistenceFailureStopsWithoutAccountMutation(t *testing.T) {
	account := expiredGrokOAuthAccountForCredentialTest(782)
	repo := &tokenRefreshAccountRepo{conditionalSuccessErr: errors.New("database unavailable")placeholder
	repo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
	cache := &grokTokenCacheForProviderTest{lockResult: trueplaceholder
	provider := NewGrokTokenProvider(repo, cache)
	provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), &tokenRefresherStub{credentials: map[string]any{
		"access_token":  "new-access-token",
		"refresh_token": "new-refresh-token",
		"expires_at":    time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
placeholderplaceholder)
	svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	_, _, err := svc.getRequestCredential(context.Background(), c, account)

	var failoverErr *UpstreamFailoverError
	require.ErrorAs(t, err, &failoverErr)
	require.Equal(t, GatewayFailureScopeProvider, failoverErr.Scope)
	require.Equal(t, GrokCredentialReasonProviderDown, failoverErr.Reason)
	require.Equal(t, NextAccountStop, failoverErr.NextAccountAction)
	require.Equal(t, 1, repo.conditionalSuccessCalls)
	require.Zero(t, repo.updateCredentialsCalls)
	require.Zero(t, repo.setErrorCalls)
	require.Zero(t, repo.setTempUnschedCalls)
	require.Empty(t, cache.deletedKeys)
	require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestGetRequestCredentialRecoversConcurrentRefreshWithoutFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)
	account := expiredGrokOAuthAccountForCredentialTest(733)
	latest := *account
	latest.Credentials = shallowCopyMap(account.Credentials)
	latest.Credentials["access_token"] = "fresh-concurrent-access"
	latest.Credentials["refresh_token"] = "fresh-concurrent-refresh"
	latest.Credentials["expires_at"] = time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339)
	latest.Credentials["_token_version"] = time.Now().UnixMilli()
	baseRepo := &tokenRefreshAccountRepo{placeholder
	baseRepo.accountsByID = map[int64]*Account{account.ID: accountplaceholder
	repo := &grokCredentialSequencedRepo{tokenRefreshAccountRepo: baseRepo, latest: &latestplaceholder
	cache := &grokTokenCacheForProviderTest{lockResult: trueplaceholder
	provider := NewGrokTokenProvider(repo, cache)
	provider.SetRefreshAPI(NewOAuthRefreshAPI(repo, cache), &tokenRefresherStub{
		err: infraerrors.New(http.StatusForbidden, "GROK_OAUTH_ENTITLEMENT_DENIED", "access_denied"),
placeholder)
	svc := &OpenAIGatewayService{accountRepo: repo, grokTokenProvider: providerplaceholder
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	token, kind, err := svc.getRequestCredential(context.Background(), c, account)

placeholder
	require.Equal(t, "fresh-concurrent-access", token)
	require.Equal(t, "oauth", kind)
	require.Zero(t, baseRepo.setErrorCalls)
	require.Zero(t, baseRepo.setTempUnschedCalls)
	require.Empty(t, cache.deletedKeys)
	require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
	_, hasEvents := c.Get(OpsUpstreamErrorsKey)
	require.False(t, hasEvents)
placeholder

func expiredGrokOAuthAccountForCredentialTest(id int64) *Account {
placeholder
		ID:          id,
		Platform:    PlatformGrok,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Schedulable: true,
placeholder
			"access_token":  "expired-access-token",
			"refresh_token": "refresh-token",
			"expires_at":    time.Now().Add(-time.Minute).UTC().Format(time.RFC3339),
			"base_url":      xai.DefaultCLIBaseURL,
	placeholder,
placeholder
placeholder
