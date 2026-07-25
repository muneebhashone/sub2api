//go:build integration

package repository

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestListDueOllamaCloudUsageAccountsOrderingLimitAndProxyHydration(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	repo := newAccountRepositoryWithSQL(tx.Client(), tx, nil)
	now := time.Date(2026, time.July, 22, 12, 0, 0, 0, time.UTC)
	proxy := mustCreateProxy(t, tx.Client(), &service.Proxy{
		Name: "ollama-due-proxy", Protocol: "http", Host: "127.0.0.1", Port: 3128,
		Username: "user", Password: "pass", Status: service.StatusActive,
placeholder)

	createAccount := func(name, baseURL string, proxyID *int64, snapshot map[string]any, lastUsed *time.Time) *service.Account {
	placeholder
		extra := map[string]any{
			service.OllamaCloudUsageSessionExtraKey:     "cipher:wos-session=fixture",
			service.OllamaCloudUsageAutoRefreshExtraKey: true,
	placeholder
		if snapshot != nil {
			extra[service.OllamaCloudUsageSnapshotExtraKey] = snapshot
	placeholder
		return mustCreateAccount(t, tx.Client(), &service.Account{
			Name: name, Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey,
	placeholder"api_key": name, "base_url": baseURLplaceholder,
			Extra:       extra, ProxyID: proxyID, LastUsedAt: lastUsed,
	placeholder)
placeholder

	uppercasePath := createAccount("ollama-uppercase-path", "https://ollama.com/V1", nil, nil, nil)
	missingSnapshot := createAccount("ollama-due-missing", "HTTPS://WWW.OLLAMA.COM:443/v1", &proxy.ID, nil, nil)
	fetched := now.Add(-2 * time.Hour)
	activity := now.Add(-5 * time.Minute)
	due := createAccount("ollama-due-activity", "https://ollama.com", nil, map[string]any{
		"status":          service.OllamaCloudUsageStatusOK,
		"fetched_at":      fetched.UTC().Format(time.RFC3339Nano),
		"last_attempt_at": fetched.UTC().Format(time.RFC3339Nano),
		"next_refresh_at": fetched.Add(time.Hour).UTC().Format(time.RFC3339Nano),
placeholder, &activity)
	// Success snapshot without newer activity must not be listed.
	_ = createAccount("ollama-not-due-idle", "https://ollama.com", nil, map[string]any{
		"status":          service.OllamaCloudUsageStatusOK,
		"fetched_at":      now.Add(-time.Hour).UTC().Format(time.RFC3339Nano),
		"last_attempt_at": now.Add(-time.Hour).UTC().Format(time.RFC3339Nano),
		"next_refresh_at": now.Add(-time.Minute).UTC().Format(time.RFC3339Nano),
placeholder, nil)
	_ = createAccount("ollama-ineligible", "https://ollama.com.evil.test", nil, nil, nil)

	accounts, err := repo.ListDueOllamaCloudUsageAccounts(ctx, now, 2)

placeholder
	require.Len(t, accounts, 2)
	require.Equal(t, missingSnapshot.ID, accounts[0].ID)
	require.Equal(t, due.ID, accounts[1].ID)
	require.NotNil(t, accounts[1].LastUsedAt)
	require.WithinDuration(t, activity.UTC(), accounts[1].LastUsedAt.UTC(), time.Second)
	require.NotContains(t, accountIDs(accounts), uppercasePath.ID)
	require.NotNil(t, accounts[0].Proxy)
	require.Equal(t, proxy.ID, accounts[0].Proxy.ID)
	require.Equal(t, proxy.URL(), accounts[0].Proxy.URL())
placeholder

func TestListDueOllamaCloudUsageAccountsUsesGroupMaxLastUsedAndFailsOpen(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	repo := newAccountRepositoryWithSQL(tx.Client(), tx, nil)
	now := time.Date(2026, time.July, 22, 14, 0, 0, 0, time.UTC)
	fetched := now.Add(-30 * time.Minute)
	older := now.Add(-10 * time.Minute)
	newer := now.Add(-2 * time.Minute)

	leader := mustCreateAccount(t, tx.Client(), &service.Account{
		Name: "ollama-group-leader", Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey,
placeholder"api_key": "shared-key", "base_url": "https://ollama.com"placeholder,
		Extra: map[string]any{
			service.OllamaCloudUsageSessionExtraKey:     "cipher:wos-session=fixture",
			service.OllamaCloudUsageAutoRefreshExtraKey: true,
			service.OllamaCloudUsageSnapshotExtraKey: map[string]any{
				"status":          service.OllamaCloudUsageStatusOK,
				"fetched_at":      fetched.UTC().Format(time.RFC3339Nano),
				"last_attempt_at": fetched.UTC().Format(time.RFC3339Nano),
				"next_refresh_at": fetched.Add(time.Hour).UTC().Format(time.RFC3339Nano),
		placeholder,
	placeholder,
		LastUsedAt: &older,
placeholder)
	_ = mustCreateAccount(t, tx.Client(), &service.Account{
		Name: "ollama-group-sibling", Platform: service.PlatformAnthropic, Type: service.AccountTypeAPIKey,
placeholder"api_key": "shared-key", "base_url": "https://www.ollama.com/v1"placeholder,
		LastUsedAt:  &newer,
placeholder)
	invalid := mustCreateAccount(t, tx.Client(), &service.Account{
		Name: "ollama-invalid-snapshot", Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey,
placeholder"api_key": "invalid-key", "base_url": "https://ollama.com"placeholder,
		Extra: map[string]any{
			service.OllamaCloudUsageSessionExtraKey:     "cipher:wos-session=fixture",
			service.OllamaCloudUsageAutoRefreshExtraKey: true,
			service.OllamaCloudUsageSnapshotExtraKey: map[string]any{
				"status": service.OllamaCloudUsageStatusOK, "fetched_at": "2026-02-30T09:00:00.123456789Z",
		placeholder,
	placeholder,
placeholder)
	idle := mustCreateAccount(t, tx.Client(), &service.Account{
		Name: "ollama-idle-ok", Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey,
placeholder"api_key": "idle-key", "base_url": "https://ollama.com"placeholder,
		Extra: map[string]any{
			service.OllamaCloudUsageSessionExtraKey:     "cipher:wos-session=fixture",
			service.OllamaCloudUsageAutoRefreshExtraKey: true,
			service.OllamaCloudUsageSnapshotExtraKey: map[string]any{
				"status":          service.OllamaCloudUsageStatusOK,
				"fetched_at":      fetched.UTC().Format(time.RFC3339Nano),
				"last_attempt_at": fetched.UTC().Format(time.RFC3339Nano),
				"next_refresh_at": fetched.Add(time.Hour).UTC().Format(time.RFC3339Nano),
		placeholder,
	placeholder,
placeholder)

	accounts, err := repo.ListDueOllamaCloudUsageAccounts(ctx, now, 10)

	require.NoError(t, err, "invalid stored values must not abort the query")
	ids := accountIDs(accounts)
	require.Contains(t, ids, invalid.ID)
	require.Contains(t, ids, leader.ID)
	require.NotContains(t, ids, idle.ID)
	for _, account := range accounts {
		if account.ID == leader.ID {
			require.NotNil(t, account.LastUsedAt)
			require.WithinDuration(t, newer.UTC(), account.LastUsedAt.UTC(), time.Second,
				"group MAX(last_used_at) must come from the sibling")
	placeholder
placeholder
placeholder

func TestLockAndMergeAccountProbeExtraCoalescesNullableOllamaGroupIdentity(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	account := mustCreateAccount(t, tx.Client(), &service.Account{
		Name: "ordinary-openai-without-base-url", Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey,
placeholder"api_key": "sk-no-base-url"placeholder,
		Extra:       map[string]any{service.UpstreamBillingProbeEnabledExtraKey: trueplaceholder,
placeholder)
	loaded, err := newAccountRepositoryWithSQL(tx.Client(), tx, nil).GetByID(ctx, account.ID)
placeholder

	merged, err := lockAndMergeAccountProbeExtra(ctx, tx.Client(), loaded, nil)

	require.NoError(t, err, "a NULL Ollama eligibility expression must scan as false")
	require.NotContains(t, merged, service.OllamaCloudUsageSessionExtraKey)
	require.Equal(t, true, merged[service.UpstreamBillingProbeEnabledExtraKey])
placeholder

func TestOllamaCloudUsageGroupWritesAreAtomicAcrossPlatformsAndURLVariants(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	repo := newAccountRepositoryWithSQL(tx.Client(), tx, nil)
	create := func(name, platform, apiKey, baseURL string) *service.Account {
	placeholder
		return mustCreateAccount(t, tx.Client(), &service.Account{
			Name: name, Platform: platform, Type: service.AccountTypeAPIKey,
	placeholder"api_key": apiKey, "base_url": baseURLplaceholder,
			Extra:       map[string]any{placeholder,
	placeholder)
placeholder
	first := create("ollama-group-openai", service.PlatformOpenAI, "shared-key", "https://ollama.com")
	second := create("ollama-group-anthropic", service.PlatformAnthropic, "shared-key", "HTTPS://WWW.OLLAMA.COM:443/v1")
	different := create("ollama-group-different", service.PlatformOpenAI, "different-key", "https://ollama.com")

	require.NoError(t, repo.SaveOllamaCloudUsageSession(ctx, first, "cipher:shared", false))
	for _, id := range []int64{first.ID, second.IDplaceholder {
		account, err := repo.GetByID(ctx, id)
	placeholder
		require.Equal(t, "cipher:shared", account.Extra[service.OllamaCloudUsageSessionExtraKey])
		require.Equal(t, false, account.Extra[service.OllamaCloudUsageAutoRefreshExtraKey])
placeholder
	differentLoaded, err := repo.GetByID(ctx, different.ID)
placeholder
	require.NotContains(t, differentLoaded.Extra, service.OllamaCloudUsageSessionExtraKey)

	secondLoaded, err := repo.GetByID(ctx, second.ID)
placeholder
	require.NoError(t, repo.SetOllamaCloudUsageAutoRefresh(ctx, secondLoaded, true))
	firstLoaded, err := repo.GetByID(ctx, first.ID)
placeholder
	secondLoaded, err = repo.GetByID(ctx, second.ID)
placeholder
	require.Equal(t, true, firstLoaded.Extra[service.OllamaCloudUsageAutoRefreshExtraKey])
	require.Equal(t, true, secondLoaded.Extra[service.OllamaCloudUsageAutoRefreshExtraKey])

	now := time.Now().UTC()
	snapshot := &service.OllamaCloudUsageSnapshot{
		Status: service.OllamaCloudUsageStatusOK, LastAttemptAt: now, NextRefreshAt: now.Add(time.Hour),
placeholder
	require.NoError(t, repo.UpdateOllamaCloudUsageSnapshot(ctx, firstLoaded, snapshot))
	secondLoaded, err = repo.GetByID(ctx, second.ID)
placeholder
	require.Equal(t, service.OllamaCloudUsageStatusOK,
		secondLoaded.Extra[service.OllamaCloudUsageSnapshotExtraKey].(map[string]any)["status"])

	staleSecond := secondLoaded
	require.NoError(t, repo.UpdateCredentials(ctx, second.ID, map[string]any{
		"api_key": "rotated-key", "base_url": "https://ollama.com",
placeholder))
	require.ErrorIs(t, repo.DisableOllamaCloudUsageAutoRefresh(ctx, staleSecond), service.ErrOllamaCloudUsageIdentityChanged)
	firstLoaded, err = repo.GetByID(ctx, first.ID)
placeholder
	secondLoaded, err = repo.GetByID(ctx, second.ID)
placeholder
	require.Equal(t, "cipher:shared", firstLoaded.Extra[service.OllamaCloudUsageSessionExtraKey])
	require.Equal(t, true, firstLoaded.Extra[service.OllamaCloudUsageAutoRefreshExtraKey])
	require.NotContains(t, secondLoaded.Extra, service.OllamaCloudUsageSessionExtraKey)
	require.NotContains(t, secondLoaded.Extra, service.OllamaCloudUsageAutoRefreshExtraKey)

	require.NoError(t, repo.DeleteOllamaCloudUsageSession(ctx, firstLoaded))
	firstLoaded, err = repo.GetByID(ctx, first.ID)
placeholder
	require.NotContains(t, firstLoaded.Extra, service.OllamaCloudUsageSessionExtraKey)
placeholder

func TestConcurrentOllamaCloudUsageSaveAndDeleteSerializeGroupState(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := testEntClient(t)
	repo := newAccountRepositoryWithSQL(client, integrationDB, nil)
	suffix := time.Now().UnixNano()
	apiKey := fmt.Sprintf("ollama-concurrent-%d", suffix)
	create := func(platform string) *service.Account {
	placeholder
		return mustCreateAccount(t, client, &service.Account{
			Name: fmt.Sprintf("%s-%s", apiKey, platform), Platform: platform, Type: service.AccountTypeAPIKey,
	placeholder"api_key": apiKey, "base_url": "https://ollama.com"placeholder,
			Extra: map[string]any{
				service.OllamaCloudUsageSessionExtraKey:     "cipher:initial",
				service.OllamaCloudUsageAutoRefreshExtraKey: true,
		placeholder,
	placeholder)
placeholder
	first := create(service.PlatformOpenAI)
	second := create(service.PlatformAnthropic)
	t.Cleanup(func() {
		_, _ = integrationDB.ExecContext(context.Background(), "DELETE FROM accounts WHERE id IN ($1, $2)", first.ID, second.ID)
placeholder)
	anchor, err := repo.GetByID(ctx, first.ID)
placeholder

	start := make(chan struct{placeholder)
	errs := make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		errs <- repo.SaveOllamaCloudUsageSession(ctx, anchor, "cipher:replacement", true)
placeholder()
	go func() {
		defer wg.Done()
		<-start
		errs <- repo.DeleteOllamaCloudUsageSession(ctx, anchor)
placeholder()
	close(start)
	wg.Wait()
	close(errs)
	for writeErr := range errs {
		require.NoError(t, writeErr)
placeholder

	firstLoaded, err := repo.GetByID(ctx, first.ID)
placeholder
	secondLoaded, err := repo.GetByID(ctx, second.ID)
placeholder
	managedState := func(account *service.Account) map[string]any {
		state := make(map[string]any)
		for _, key := range []string{
			service.OllamaCloudUsageSessionExtraKey,
			service.OllamaCloudUsageAutoRefreshExtraKey,
			service.OllamaCloudUsageSnapshotExtraKey,
	placeholder {
			if value, ok := account.Extra[key]; ok {
				state[key] = value
		placeholder
	placeholder
		return state
placeholder
	firstState := managedState(firstLoaded)
	require.Equal(t, firstState, managedState(secondLoaded), "a serialized last commit must own the whole group")
	if len(firstState) > 0 {
		require.Equal(t, "cipher:replacement", firstState[service.OllamaCloudUsageSessionExtraKey])
		require.Equal(t, true, firstState[service.OllamaCloudUsageAutoRefreshExtraKey])
		require.NotContains(t, firstState, service.OllamaCloudUsageSnapshotExtraKey)
placeholder
placeholder

func accountIDs(accounts []service.Account) []int64 {
	ids := make([]int64, len(accounts))
	for index := range accounts {
		ids[index] = accounts[index].ID
placeholder
	return ids
placeholder

func TestOllamaCloudUsageCredentialAndBulkUpdatesPreserveManagedStateOnlyWhenSafe(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	repo := newAccountRepositoryWithSQL(tx.Client(), tx, nil)
	now := time.Now().UTC()
	newAccount := func(name string) *service.Account {
		return mustCreateAccount(t, tx.Client(), &service.Account{
			Name: name, Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey,
	placeholder"api_key": "old-key", "base_url": "https://ollama.com"placeholder,
			Extra: map[string]any{
				service.OllamaCloudUsageSessionExtraKey:     "cipher:wos-session=fixture",
				service.OllamaCloudUsageAutoRefreshExtraKey: true,
				service.OllamaCloudUsageSnapshotExtraKey: map[string]any{
					"status": service.OllamaCloudUsageStatusOK, "last_attempt_at": now, "next_refresh_at": now.Add(time.Hour),
			placeholder,
		placeholder,
	placeholder)
placeholder

	rawAccount := newAccount("ollama-raw-credentials")
	require.NoError(t, repo.UpdateCredentials(ctx, rawAccount.ID, map[string]any{
		"api_key": "old-key", "base_url": "https://ollama.com/V1",
placeholder))
	rawUpdated, err := repo.GetByID(ctx, rawAccount.ID)
placeholder
	require.NotContains(t, rawUpdated.Extra, service.OllamaCloudUsageSessionExtraKey)
	require.NotContains(t, rawUpdated.Extra, service.OllamaCloudUsageAutoRefreshExtraKey)
	require.NotContains(t, rawUpdated.Extra, service.OllamaCloudUsageSnapshotExtraKey)

	bulkAccount := newAccount("ollama-bulk-credentials")
	rows, err := repo.BulkUpdate(ctx, []int64{bulkAccount.IDplaceholder, service.AccountBulkUpdate{
placeholder"base_url": "HTTPS://WWW.OLLAMA.COM:443/v1"placeholder,
placeholder)
placeholder
	require.Equal(t, int64(1), rows)
	bulkUnchanged, err := repo.GetByID(ctx, bulkAccount.ID)
placeholder
	require.Contains(t, bulkUnchanged.Extra, service.OllamaCloudUsageSnapshotExtraKey)

	rows, err = repo.BulkUpdate(ctx, []int64{bulkAccount.IDplaceholder, service.AccountBulkUpdate{
placeholder"base_url": "https://ollama.com/V1"placeholder,
placeholder)
placeholder
	require.Equal(t, int64(1), rows)
	bulkIneligible, err := repo.GetByID(ctx, bulkAccount.ID)
placeholder
	require.NotContains(t, bulkIneligible.Extra, service.OllamaCloudUsageSessionExtraKey)
	require.NotContains(t, bulkIneligible.Extra, service.OllamaCloudUsageAutoRefreshExtraKey)
	require.NotContains(t, bulkIneligible.Extra, service.OllamaCloudUsageSnapshotExtraKey)
placeholder

func TestProxyIdentityUpdateInvalidatesOllamaSnapshotAndRejectsInFlightCAS(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	accountRepo := newAccountRepositoryWithSQL(tx.Client(), tx, nil)
	proxyRepo := newProxyRepositoryWithSQL(tx.Client(), tx)
	proxy := mustCreateProxy(t, tx.Client(), &service.Proxy{
		Name: "ollama-identity-proxy", Protocol: "http", Host: "old.example", Port: 8080,
		Username: "old-user", Password: "old-pass", Status: service.StatusActive,
placeholder)
	now := time.Now().UTC()
	account := mustCreateAccount(t, tx.Client(), &service.Account{
		Name: "ollama-proxy-account", Platform: service.PlatformAnthropic, Type: service.AccountTypeAPIKey,
placeholder"api_key": "key", "base_url": "https://ollama.com"placeholder,
		ProxyID:     &proxy.ID,
		Extra: map[string]any{
			service.OllamaCloudUsageSessionExtraKey:     "cipher:wos-session=fixture",
			service.OllamaCloudUsageAutoRefreshExtraKey: true,
			service.OllamaCloudUsageSnapshotExtraKey: map[string]any{
				"status": service.OllamaCloudUsageStatusOK, "last_attempt_at": now, "next_refresh_at": now.Add(time.Hour),
		placeholder,
	placeholder,
placeholder)
	inFlight, err := accountRepo.GetByID(ctx, account.ID)
placeholder
	require.NotNil(t, inFlight.Proxy)
	require.Equal(t, "old.example", inFlight.Proxy.Host)

	proxyToUpdate, err := proxyRepo.GetByID(ctx, proxy.ID)
placeholder
	proxyToUpdate.Host = "new.example"
	require.NoError(t, proxyRepo.Update(ctx, proxyToUpdate))

	got, err := accountRepo.GetByID(ctx, account.ID)
placeholder
	require.NotContains(t, got.Extra, service.OllamaCloudUsageSnapshotExtraKey)
	require.Equal(t, "cipher:wos-session=fixture", got.Extra[service.OllamaCloudUsageSessionExtraKey])
	require.Equal(t, true, got.Extra[service.OllamaCloudUsageAutoRefreshExtraKey])

	err = accountRepo.UpdateOllamaCloudUsageSnapshot(ctx, inFlight, &service.OllamaCloudUsageSnapshot{
		Status: service.OllamaCloudUsageStatusOK, LastAttemptAt: now, NextRefreshAt: now.Add(time.Hour),
placeholder)
	require.ErrorIs(t, err, service.ErrOllamaCloudUsageIdentityChanged)
placeholder

// 无变化的凭证持久化（如 CRS 同步重放同一凭证）不得触发任何 extra 清理；
// 真实变化仍必须按旧语义清 openai 探测快照。
func TestUpdateCredentialsUnchangedCredentialsPreserveManagedExtra(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	repo := newAccountRepositoryWithSQL(tx.Client(), tx, nil)

	probeAccount := mustCreateAccount(t, tx.Client(), &service.Account{
		Name: "openai-probe-unchanged", Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey,
placeholder"api_key": "sk-probe", "base_url": "https://relay.example.com/v1"placeholder,
		Extra: map[string]any{
			service.UpstreamBillingProbeEnabledExtraKey: true,
			service.UpstreamBillingProbeExtraKey:        map[string]any{"status": "ok"placeholder,
	placeholder,
placeholder)
	require.NoError(t, repo.UpdateCredentials(ctx, probeAccount.ID, map[string]any{
		"api_key": "sk-probe", "base_url": "https://relay.example.com/v1",
placeholder))
	probeLoaded, err := repo.GetByID(ctx, probeAccount.ID)
placeholder
	require.Contains(t, probeLoaded.Extra, service.UpstreamBillingProbeExtraKey,
		"unchanged credentials must not clear the probe snapshot")

	now := time.Now().UTC()
	ollamaAccount := mustCreateAccount(t, tx.Client(), &service.Account{
		Name: "ollama-unchanged", Platform: service.PlatformAnthropic, Type: service.AccountTypeAPIKey,
placeholder"api_key": "ollama-key", "base_url": "https://ollama.com"placeholder,
		Extra: map[string]any{
			service.OllamaCloudUsageSessionExtraKey:     "cipher:wos-session=fixture",
			service.OllamaCloudUsageAutoRefreshExtraKey: true,
			service.OllamaCloudUsageSnapshotExtraKey: map[string]any{
				"status": service.OllamaCloudUsageStatusOK, "last_attempt_at": now, "next_refresh_at": now.Add(time.Hour),
		placeholder,
	placeholder,
placeholder)
	require.NoError(t, repo.UpdateCredentials(ctx, ollamaAccount.ID, map[string]any{
		"api_key": "ollama-key", "base_url": "https://ollama.com",
placeholder))
	ollamaLoaded, err := repo.GetByID(ctx, ollamaAccount.ID)
placeholder
	require.Equal(t, "cipher:wos-session=fixture", ollamaLoaded.Extra[service.OllamaCloudUsageSessionExtraKey])
	require.Equal(t, true, ollamaLoaded.Extra[service.OllamaCloudUsageAutoRefreshExtraKey])
	require.Contains(t, ollamaLoaded.Extra, service.OllamaCloudUsageSnapshotExtraKey)

	require.NoError(t, repo.UpdateCredentials(ctx, probeAccount.ID, map[string]any{
		"api_key": "sk-probe", "base_url": "https://relay.example.org/v1",
placeholder))
	probeLoaded, err = repo.GetByID(ctx, probeAccount.ID)
placeholder
	require.NotContains(t, probeLoaded.Extra, service.UpstreamBillingProbeExtraKey,
		"changed credentials must keep clearing the probe snapshot")
placeholder
