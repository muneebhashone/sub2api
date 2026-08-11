//go:build integration

package repository

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestAccountUpdatePreservesConcurrentProbeSnapshot(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	repo := newAccountRepositoryWithSQL(tx.Client(), tx, nil)
	account := mustCreateAccount(t, tx.Client(), &service.Account{
		Name:        "probe-update-preserve",
		Platform:    service.PlatformOpenAI,
		Type:        service.AccountTypeAPIKey,
placeholder"api_key": "sk-old"placeholder,
		Extra:       map[string]any{service.UpstreamBillingProbeEnabledExtraKey: trueplaceholder,
placeholder)

	stale, err := repo.GetByID(ctx, account.ID)
placeholder
	require.NotContains(t, stale.Extra, service.UpstreamBillingProbeExtraKey)
	require.NoError(t, repo.UpdateUpstreamBillingProbeSnapshot(ctx, stale, &service.UpstreamBillingProbeSnapshot{
		Status:        service.UpstreamBillingProbeStatusOK,
		LastAttemptAt: time.Now().UTC(),
placeholder, nil))

	stale.Name = "ordinary-edit"
	require.NoError(t, repo.Update(ctx, stale))
	got, err := repo.GetByID(ctx, account.ID)
placeholder
	snapshot, ok := got.Extra[service.UpstreamBillingProbeExtraKey].(map[string]any)
	require.True(t, ok)
	require.Equal(t, service.UpstreamBillingProbeStatusOK, snapshot["status"])

	require.NoError(t, repo.UpdateExtra(ctx, got.ID, map[string]any{service.UpstreamBillingProbeEnabledExtraKey: falseplaceholder))
	disabled, err := repo.GetByID(ctx, account.ID)
placeholder
	require.NotContains(t, disabled.Extra, service.UpstreamBillingProbeExtraKey)
placeholder

func TestAdminAccountEditPreservesRateSynchronizedAfterLoad(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	repo := newAccountRepositoryWithSQL(tx.Client(), tx, nil)
	initialRate := 0.1
	account := mustCreateAccount(t, tx.Client(), &service.Account{
		Name:           "probe-rate-concurrent-edit",
		Platform:       service.PlatformOpenAI,
		Type:           service.AccountTypeAPIKey,
		RateMultiplier: &initialRate,
		Credentials:    map[string]any{"api_key": "sk-test"placeholder,
		Extra: map[string]any{
			service.UpstreamBillingProbeEnabledExtraKey:    true,
			service.UpstreamBillingRateSyncEnabledExtraKey: true,
	placeholder,
placeholder)

	staleAdminEdit, err := repo.GetByID(ctx, account.ID)
placeholder
	probeAccount, err := repo.GetByID(ctx, account.ID)
placeholder

	synchronizedRate := 0.2
	require.NoError(t, repo.UpdateUpstreamBillingProbeSnapshot(ctx, probeAccount, &service.UpstreamBillingProbeSnapshot{
		Status:        service.UpstreamBillingProbeStatusOK,
		LastAttemptAt: time.Now().UTC(),
placeholder, &synchronizedRate))

	staleAdminEdit.Name = "name-only-edit"
	require.NoError(t, repo.UpdateWithAccountBillingSettings(ctx, staleAdminEdit, nil, nil, nil))

	got, err := repo.GetByID(ctx, account.ID)
placeholder
	require.Equal(t, "name-only-edit", got.Name)
	require.NotNil(t, got.RateMultiplier)
	require.Equal(t, synchronizedRate, *got.RateMultiplier)
placeholder

func TestProbeSnapshotSyncsRateOnlyForSuccessfulEnabledAccount(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	repo := newAccountRepositoryWithSQL(tx.Client(), tx, nil)
	initialRate := 0.25
	account := mustCreateAccount(t, tx.Client(), &service.Account{
		Name:           "probe-rate-sync",
		Platform:       service.PlatformGemini,
		Type:           service.AccountTypeAPIKey,
		RateMultiplier: &initialRate,
		Credentials:    map[string]any{"api_key": "sk-test"placeholder,
		Extra: map[string]any{
			service.UpstreamBillingProbeEnabledExtraKey:    true,
			service.UpstreamBillingRateSyncEnabledExtraKey: true,
	placeholder,
placeholder)

	loaded, err := repo.GetByID(ctx, account.ID)
placeholder
	syncedRate := 0.065
	require.NoError(t, repo.UpdateUpstreamBillingProbeSnapshot(ctx, loaded, &service.UpstreamBillingProbeSnapshot{
		Status:        service.UpstreamBillingProbeStatusOK,
		LastAttemptAt: time.Now().UTC(),
placeholder, &syncedRate))

	got, err := repo.GetByID(ctx, account.ID)
placeholder
	require.NotNil(t, got.RateMultiplier)
	require.Equal(t, syncedRate, *got.RateMultiplier)

	failedRate := 0.9
	require.NoError(t, repo.UpdateUpstreamBillingProbeSnapshot(ctx, got, &service.UpstreamBillingProbeSnapshot{
		Status:        service.UpstreamBillingProbeStatusFailed,
		LastAttemptAt: time.Now().UTC(),
placeholder, &failedRate))
	got, err = repo.GetByID(ctx, account.ID)
placeholder
	require.NotNil(t, got.RateMultiplier)
	require.Equal(t, syncedRate, *got.RateMultiplier)

	require.NoError(t, repo.UpdateExtra(ctx, account.ID, map[string]any{
		service.UpstreamBillingRateSyncEnabledExtraKey: false,
placeholder))
	manual, err := repo.GetByID(ctx, account.ID)
placeholder
	manualProbeRate := 0.4
	require.NoError(t, repo.UpdateUpstreamBillingProbeSnapshot(ctx, manual, &service.UpstreamBillingProbeSnapshot{
		Status:        service.UpstreamBillingProbeStatusOK,
		LastAttemptAt: time.Now().UTC(),
placeholder, &manualProbeRate))
	got, err = repo.GetByID(ctx, account.ID)
placeholder
	require.NotNil(t, got.RateMultiplier)
	require.Equal(t, syncedRate, *got.RateMultiplier)

	require.NoError(t, repo.UpdateExtra(ctx, account.ID, map[string]any{
		service.UpstreamBillingProbeEnabledExtraKey:    false,
		service.UpstreamBillingRateSyncEnabledExtraKey: false,
placeholder))
	disabled, err := repo.GetByID(ctx, account.ID)
placeholder
	require.NoError(t, repo.UpdateUpstreamBillingProbeSnapshot(ctx, disabled, &service.UpstreamBillingProbeSnapshot{
		Status:        service.UpstreamBillingProbeStatusOK,
		LastAttemptAt: time.Now().UTC(),
placeholder, &manualProbeRate))
	got, err = repo.GetByID(ctx, account.ID)
placeholder
	require.NotNil(t, got.RateMultiplier)
	require.Equal(t, syncedRate, *got.RateMultiplier)
placeholder

func TestAccountUpdatePreservesConcurrentProbeEnableFlag(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	repo := newAccountRepositoryWithSQL(tx.Client(), tx, nil)
	account := mustCreateAccount(t, tx.Client(), &service.Account{
		Name:        "probe-update-enable",
		Platform:    service.PlatformOpenAI,
		Type:        service.AccountTypeAPIKey,
placeholder"api_key": "sk-test"placeholder,
		Extra: map[string]any{
			service.UpstreamBillingProbeEnabledExtraKey: true,
			service.UpstreamBillingProbeExtraKey:        map[string]any{"status": service.UpstreamBillingProbeStatusOKplaceholder,
	placeholder,
placeholder)

	stale, err := repo.GetByID(ctx, account.ID)
placeholder
	require.NoError(t, repo.UpdateExtra(ctx, account.ID, map[string]any{service.UpstreamBillingProbeEnabledExtraKey: falseplaceholder))
	stale.Name = "ordinary-edit"
	require.NoError(t, repo.Update(ctx, stale))

	got, err := repo.GetByID(ctx, account.ID)
placeholder
	require.Equal(t, false, got.Extra[service.UpstreamBillingProbeEnabledExtraKey])
	require.NotContains(t, got.Extra, service.UpstreamBillingProbeExtraKey)
placeholder

func TestAccountUpdateClearsProbeSnapshotWhenIdentityChanges(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	repo := newAccountRepositoryWithSQL(tx.Client(), tx, nil)
	account := mustCreateAccount(t, tx.Client(), &service.Account{
		Name:        "probe-update-identity",
		Platform:    service.PlatformOpenAI,
		Type:        service.AccountTypeAPIKey,
placeholder"api_key": "sk-old"placeholder,
		Extra: map[string]any{
			service.UpstreamBillingProbeEnabledExtraKey: true,
			service.UpstreamBillingProbeExtraKey:        map[string]any{"status": service.UpstreamBillingProbeStatusOKplaceholder,
	placeholder,
placeholder)

	loaded, err := repo.GetByID(ctx, account.ID)
placeholder
	loaded.Credentials["api_key"] = "sk-new"
	require.NoError(t, repo.Update(ctx, loaded))

	got, err := repo.GetByID(ctx, account.ID)
placeholder
	require.NotContains(t, got.Extra, service.UpstreamBillingProbeExtraKey)
placeholder

func TestBulkUpdateAndCredentialUpdateDeleteProbeKey(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	repo := newAccountRepositoryWithSQL(tx.Client(), tx, nil)
	newAccount := func(name string) *service.Account {
		return mustCreateAccount(t, tx.Client(), &service.Account{
			Name:        name,
			Platform:    service.PlatformOpenAI,
			Type:        service.AccountTypeAPIKey,
	placeholder"api_key": "sk-old"placeholder,
			Extra: map[string]any{
				service.UpstreamBillingProbeEnabledExtraKey: true,
				service.UpstreamBillingProbeExtraKey:        map[string]any{"status": service.UpstreamBillingProbeStatusOKplaceholder,
		placeholder,
	placeholder)
placeholder

	bulkAccount := newAccount("probe-bulk-clear")
	_, err := repo.BulkUpdate(ctx, []int64{bulkAccount.IDplaceholder, service.AccountBulkUpdate{
		Extra: map[string]any{service.UpstreamBillingProbeExtraKey: nilplaceholder,
placeholder)
placeholder
	got, err := repo.GetByID(ctx, bulkAccount.ID)
placeholder
	require.NotContains(t, got.Extra, service.UpstreamBillingProbeExtraKey)

	credentialAccount := newAccount("probe-credentials-clear")
	require.NoError(t, repo.UpdateCredentials(ctx, credentialAccount.ID, map[string]any{"api_key": "sk-new"placeholder))
	got, err = repo.GetByID(ctx, credentialAccount.ID)
placeholder
	require.NotContains(t, got.Extra, service.UpstreamBillingProbeExtraKey)
placeholder

func TestProbeSnapshotCASIncludesLoadedEnabledState(t *testing.T) {
	tests := []struct {
		name           string
		loadedEnabled  bool
		concurrentFlip *bool
		wantConflict   bool
placeholder{
		{name: "manual_false_stays_false", loadedEnabled: falseplaceholder,
		{name: "periodic_true_disabled_in_flight", loadedEnabled: true, concurrentFlip: boolPtr(false), wantConflict: trueplaceholder,
		{name: "manual_false_enabled_in_flight", loadedEnabled: false, concurrentFlip: boolPtr(true), wantConflict: trueplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tx := testEntTx(t)
			repo := newAccountRepositoryWithSQL(tx.Client(), tx, nil)
			account := mustCreateAccount(t, tx.Client(), &service.Account{
				Name:        "probe-enabled-cas-" + tt.name,
				Platform:    service.PlatformOpenAI,
				Type:        service.AccountTypeAPIKey,
		placeholder"api_key": "sk-test"placeholder,
				Extra:       map[string]any{service.UpstreamBillingProbeEnabledExtraKey: tt.loadedEnabledplaceholder,
		placeholder)
			inFlight, err := repo.GetByID(ctx, account.ID)
		placeholder
			if tt.concurrentFlip != nil {
				require.NoError(t, repo.UpdateExtra(ctx, account.ID, map[string]any{
					service.UpstreamBillingProbeEnabledExtraKey: *tt.concurrentFlip,
			placeholder))
		placeholder

			err = repo.UpdateUpstreamBillingProbeSnapshot(ctx, inFlight, &service.UpstreamBillingProbeSnapshot{
				Status:        service.UpstreamBillingProbeStatusOK,
				LastAttemptAt: time.Now().UTC(),
		placeholder, nil)
			if tt.wantConflict {
				require.ErrorIs(t, err, service.ErrUpstreamBillingProbeIdentityChanged)
		placeholder else {
			placeholder
		placeholder
			got, err := repo.GetByID(ctx, account.ID)
		placeholder
			if tt.wantConflict {
				require.NotContains(t, got.Extra, service.UpstreamBillingProbeExtraKey)
		placeholder else {
				require.Contains(t, got.Extra, service.UpstreamBillingProbeExtraKey)
		placeholder
	placeholder)
placeholder
placeholder

func TestProbeSnapshotCASProtectsManualRateAfterSyncDisabled(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	repo := newAccountRepositoryWithSQL(tx.Client(), tx, nil)
	initialRate := 0.25
	account := mustCreateAccount(t, tx.Client(), &service.Account{
		Name:           "probe-sync-cas",
		Platform:       service.PlatformAnthropic,
		Type:           service.AccountTypeAPIKey,
		RateMultiplier: &initialRate,
		Credentials:    map[string]any{"api_key": "sk-test"placeholder,
		Extra: map[string]any{
			service.UpstreamBillingProbeEnabledExtraKey:    true,
			service.UpstreamBillingRateSyncEnabledExtraKey: true,
	placeholder,
placeholder)

	inFlight, err := repo.GetByID(ctx, account.ID)
placeholder
	manual, err := repo.GetByID(ctx, account.ID)
placeholder
	manualRate := 0.8
	syncDisabled := false
	require.NoError(t, repo.UpdateWithAccountBillingSettings(ctx, manual, nil, &syncDisabled, &manualRate))

	probedRate := 0.1
	err = repo.UpdateUpstreamBillingProbeSnapshot(ctx, inFlight, &service.UpstreamBillingProbeSnapshot{
		Status:        service.UpstreamBillingProbeStatusOK,
		LastAttemptAt: time.Now().UTC(),
placeholder, &probedRate)
	require.ErrorIs(t, err, service.ErrUpstreamBillingProbeIdentityChanged)

	got, err := repo.GetByID(ctx, account.ID)
placeholder
	require.NotNil(t, got.RateMultiplier)
	require.Equal(t, manualRate, *got.RateMultiplier)
	require.Equal(t, false, got.Extra[service.UpstreamBillingRateSyncEnabledExtraKey])
	require.NotContains(t, got.Extra, service.UpstreamBillingProbeExtraKey)
placeholder

func boolPtr(value bool) *bool {
	return &value
placeholder

func TestProxyIdentityUpdateInvalidatesProbeAndRejectsInFlightSnapshot(t *testing.T) {
	tests := []struct {
		name             string
		includeProbeKey  bool
		probeValue       any
		wantInvalidation bool
placeholder{
		{name: "missing_snapshot"placeholder,
		{name: "json_null_snapshot", includeProbeKey: trueplaceholder,
		{name: "existing_snapshot", includeProbeKey: true, probeValue: map[string]any{"status": service.UpstreamBillingProbeStatusOKplaceholder, wantInvalidation: trueplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tx := testEntTx(t)
			accountRepo := newAccountRepositoryWithSQL(tx.Client(), tx, nil)
			proxyRepo := newProxyRepositoryWithSQL(tx.Client(), tx)
			proxy := mustCreateProxy(t, tx.Client(), &service.Proxy{
				Name:     "probe-proxy",
				Protocol: "http",
				Host:     "old.example",
				Port:     8080,
				Username: "old-user",
				Password: "old-pass",
				Status:   service.StatusActive,
		placeholder)
			extra := map[string]any{service.UpstreamBillingProbeEnabledExtraKey: trueplaceholder
			if tt.includeProbeKey {
				extra[service.UpstreamBillingProbeExtraKey] = tt.probeValue
		placeholder
			account := mustCreateAccount(t, tx.Client(), &service.Account{
				Name:        "proxy-probe-account",
				Platform:    service.PlatformOpenAI,
				Type:        service.AccountTypeAPIKey,
		placeholder"api_key": "sk-test"placeholder,
				Extra:       extra,
				ProxyID:     &proxy.ID,
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
			if tt.wantInvalidation || !tt.includeProbeKey {
				require.NotContains(t, got.Extra, service.UpstreamBillingProbeExtraKey)
		placeholder else {
				require.Contains(t, got.Extra, service.UpstreamBillingProbeExtraKey)
				require.Nil(t, got.Extra[service.UpstreamBillingProbeExtraKey])
		placeholder
			if !tt.wantInvalidation {
				require.Equal(t, inFlight.UpdatedAt, got.UpdatedAt, "missing/null snapshots must not cause an account row write")
		placeholder
			err = accountRepo.UpdateUpstreamBillingProbeSnapshot(ctx, inFlight, &service.UpstreamBillingProbeSnapshot{
				Status:        service.UpstreamBillingProbeStatusOK,
				LastAttemptAt: time.Now().UTC(),
		placeholder, nil)
			require.ErrorIs(t, err, service.ErrUpstreamBillingProbeIdentityChanged)

			rows, err := tx.QueryContext(ctx, `
				SELECT COUNT(*), COALESCE(MAX(payload::text), '')
				FROM scheduler_outbox
				WHERE event_type = $1
			`, service.SchedulerOutboxEventAccountBulkChanged)
		placeholder
			require.True(t, rows.Next())
			var (
				outboxCount int
				payloadJSON string
			)
			require.NoError(t, rows.Scan(&outboxCount, &payloadJSON))
			require.NoError(t, rows.Close())
			if tt.wantInvalidation {
				require.Equal(t, 1, outboxCount)
				var payload struct {
					AccountIDs []int64 `json:"account_ids"`
			placeholder
				require.NoError(t, json.Unmarshal([]byte(payloadJSON), &payload))
				require.Equal(t, []int64{account.IDplaceholder, payload.AccountIDs)
		placeholder else {
				require.Zero(t, outboxCount, "no snapshot change means no PR2 cache invalidation event")
		placeholder
	placeholder)
placeholder
placeholder

func TestSweepExpiredProxyWithoutFallbackInvalidatesOnlyExistingProbeSnapshot(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	proxyRepo := newProxyRepositoryWithSQL(tx.Client(), tx)
	accountRepo := newAccountRepositoryWithSQL(tx.Client(), tx, nil)
	past := time.Now().Add(-time.Hour)
	proxy := &service.Proxy{
		Name:           "expired-probe-proxy-none",
		Protocol:       "http",
		Host:           "127.0.0.1",
		Port:           8080,
		Status:         service.StatusActive,
		ExpiresAt:      &past,
		FallbackMode:   service.FallbackModeNone,
		ExpiryWarnDays: 7,
placeholder
	require.NoError(t, proxyRepo.Create(ctx, proxy))
	newAccount := func(name string, probe any, includeProbe bool) *service.Account {
		extra := map[string]any{service.UpstreamBillingProbeEnabledExtraKey: trueplaceholder
		if includeProbe {
			extra[service.UpstreamBillingProbeExtraKey] = probe
	placeholder
		return mustCreateAccount(t, tx.Client(), &service.Account{
			Name:        name,
			Platform:    service.PlatformOpenAI,
			Type:        service.AccountTypeAPIKey,
	placeholder"api_key": "sk-test"placeholder,
			Extra:       extra,
			ProxyID:     &proxy.ID,
	placeholder)
placeholder
	withSnapshot := newAccount("expired-proxy-with-snapshot", map[string]any{"status": service.UpstreamBillingProbeStatusOKplaceholder, true)
	withoutSnapshot := newAccount("expired-proxy-without-snapshot", nil, false)
	withJSONNull := newAccount("expired-proxy-null-snapshot", nil, true)
	untouchedUpdatedAt := make(map[int64]time.Time, 2)
	for _, untouched := range []*service.Account{withoutSnapshot, withJSONNullplaceholder {
		loaded, err := accountRepo.GetByID(ctx, untouched.ID)
	placeholder
		untouchedUpdatedAt[untouched.ID] = loaded.UpdatedAt
placeholder

	changed, err := proxyRepo.SweepExpiredProxies(ctx, time.Now())
placeholder
	require.Zero(t, changed, "probe invalidation must not inflate the rerouted account count")

	got, err := accountRepo.GetByID(ctx, withSnapshot.ID)
placeholder
	require.NotContains(t, got.Extra, service.UpstreamBillingProbeExtraKey)
	for _, untouched := range []*service.Account{withoutSnapshot, withJSONNullplaceholder {
		got, err = accountRepo.GetByID(ctx, untouched.ID)
	placeholder
		require.Equal(t, untouchedUpdatedAt[untouched.ID], got.UpdatedAt)
placeholder

	payload := latestBulkAccountOutboxPayload(t, ctx, tx)
	require.Equal(t, []int64{withSnapshot.IDplaceholder, payload)
placeholder

func TestSweepExpiredProxyFallbackRerouteDeletesProbeSnapshot(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	proxyRepo := newProxyRepositoryWithSQL(tx.Client(), tx)
	accountRepo := newAccountRepositoryWithSQL(tx.Client(), tx, nil)
	past := time.Now().Add(-time.Hour)
	proxy := &service.Proxy{
		Name:           "expired-probe-proxy-direct",
		Protocol:       "http",
		Host:           "127.0.0.1",
		Port:           8080,
		Status:         service.StatusActive,
		ExpiresAt:      &past,
		FallbackMode:   service.FallbackModeDirect,
		ExpiryWarnDays: 7,
placeholder
	require.NoError(t, proxyRepo.Create(ctx, proxy))
	account := mustCreateAccount(t, tx.Client(), &service.Account{
		Name:        "expired-proxy-rerouted-snapshot",
		Platform:    service.PlatformOpenAI,
		Type:        service.AccountTypeAPIKey,
placeholder"api_key": "sk-test"placeholder,
		Extra: map[string]any{
			service.UpstreamBillingProbeEnabledExtraKey: true,
			service.UpstreamBillingProbeExtraKey:        map[string]any{"status": service.UpstreamBillingProbeStatusOKplaceholder,
	placeholder,
		ProxyID: &proxy.ID,
placeholder)

	changed, err := proxyRepo.SweepExpiredProxies(ctx, time.Now())
placeholder
	require.EqualValues(t, 1, changed)

	got, err := accountRepo.GetByID(ctx, account.ID)
placeholder
	require.Nil(t, got.ProxyID)
	require.NotContains(t, got.Extra, service.UpstreamBillingProbeExtraKey)
	require.Equal(t, []int64{account.IDplaceholder, latestBulkAccountOutboxPayload(t, ctx, tx))
placeholder

func latestBulkAccountOutboxPayload(t *testing.T, ctx context.Context, tx sqlQueryer) []int64 {
placeholder
	var payloadJSON []byte
	require.NoError(t, scanSingleRow(ctx, tx, `
		SELECT payload
		FROM scheduler_outbox
		WHERE event_type = $1
		ORDER BY id DESC
		LIMIT 1
	`, []any{service.SchedulerOutboxEventAccountBulkChangedplaceholder, &payloadJSON))
	var payload struct {
		AccountIDs []int64 `json:"account_ids"`
placeholder
	require.NoError(t, json.Unmarshal(payloadJSON, &payload))
	return payload.AccountIDs
placeholder
