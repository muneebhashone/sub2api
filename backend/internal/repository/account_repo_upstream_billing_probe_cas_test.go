package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
)

func TestUpdateUpstreamBillingProbeSnapshotRequiresSameIdentityAndSnapshot(t *testing.T) {
	tests := []struct {
		name     string
		affected int64
		wantErr  error
placeholder{
		{name: "same identity and snapshot", affected: 1placeholder,
		{name: "identity or snapshot changed", affected: 0, wantErr: service.ErrUpstreamBillingProbeIdentityChangedplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
		placeholder
			t.Cleanup(func() { _ = db.Close() placeholder)
			driver := entsql.OpenDB(dialect.Postgres, db)
			client := dbent.NewClient(dbent.Driver(driver))
			t.Cleanup(func() { _ = client.Close() placeholder)

			mock.ExpectBegin()
			tx, err := client.Tx(context.Background())
		placeholder
			mock.ExpectQuery(`(?s)` + regexp.QuoteMeta("SELECT protocol, host, port") + `.*` + regexp.QuoteMeta("FOR SHARE")).
				WithArgs(int64(9)).
				WillReturnRows(sqlmock.NewRows([]string{"protocol", "host", "port", "username", "password", "status"placeholder).
					AddRow("http", "127.0.0.1", 3128, "user", "pass", service.StatusActive))
			mock.ExpectExec(`(?s)`+regexp.QuoteMeta("UPDATE accounts")+`.*`+regexp.QuoteMeta("WHERE id = $2")+`.*`+regexp.QuoteMeta("AND platform = $3")+`.*`+regexp.QuoteMeta("AND type = $4")+`.*`+regexp.QuoteMeta("AND credentials = $5::jsonb")+`.*`+regexp.QuoteMeta("AND proxy_id IS NOT DISTINCT FROM $6")+`.*`+regexp.QuoteMeta("COALESCE(extra -> 'upstream_billing_probe', 'null'::jsonb) = $7::jsonb")+`.*`+regexp.QuoteMeta("COALESCE(extra -> 'upstream_billing_probe_enabled', 'null'::jsonb) = $8::jsonb")+`.*`+regexp.QuoteMeta("COALESCE(extra -> 'upstream_billing_rate_sync_enabled', 'null'::jsonb) = $9::jsonb")).
				WithArgs(sqlmock.AnyArg(), int64(17), service.PlatformOpenAI, service.AccountTypeAPIKey, `{"api_key":"sk-test","base_url":"http://127.0.0.1:8080"placeholder`, int64(9), `{"status":"stale"placeholder`, "null", "null", nil).
				WillReturnResult(sqlmock.NewResult(0, tt.affected))
			if tt.affected > 0 {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO scheduler_outbox")).
					WithArgs(service.SchedulerOutboxEventAccountChanged, int64(17), nil, nil, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
		placeholder
			repo := newAccountRepositoryWithSQL(client, &recordingSQLExecutor{err: errors.New("must use transaction client")placeholder, nil)
			proxyID := int64(9)
			account := &service.Account{
				ID:       17,
				Platform: service.PlatformOpenAI,
				Type:     service.AccountTypeAPIKey,
		placeholder
					"api_key":  "sk-test",
					"base_url": "http://127.0.0.1:8080",
			placeholder,
				ProxyID: &proxyID,
				Proxy: &service.Proxy{
					ID:       proxyID,
					Protocol: "http",
					Host:     "127.0.0.1",
					Port:     3128,
					Username: "user",
					Password: "pass",
					Status:   service.StatusActive,
			placeholder,
				Extra: map[string]any{
					service.UpstreamBillingProbeExtraKey: map[string]any{"status": "stale"placeholder,
			placeholder,
		placeholder

			txCtx := dbent.NewTxContext(context.Background(), tx)
			err = repo.UpdateUpstreamBillingProbeSnapshot(txCtx, account, &service.UpstreamBillingProbeSnapshot{Status: service.UpstreamBillingProbeStatusOKplaceholder, nil)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
		placeholder else {
			placeholder
		placeholder
			mock.ExpectRollback()
			require.NoError(t, tx.Rollback())
			require.NoError(t, mock.ExpectationsWereMet())
	placeholder)
placeholder
placeholder

func TestUpdateUpstreamBillingProbeSnapshotCommitsSnapshotAndOutboxAtomically(t *testing.T) {
	db, mock, err := sqlmock.New()
placeholder
	t.Cleanup(func() { _ = db.Close() placeholder)
	driver := entsql.OpenDB(dialect.Postgres, db)
	client := dbent.NewClient(dbent.Driver(driver))
	t.Cleanup(func() { _ = client.Close() placeholder)

	mock.ExpectBegin()
	mock.ExpectExec(`(?s)`+regexp.QuoteMeta("UPDATE accounts")+`.*`+regexp.QuoteMeta("rate_multiplier = CASE")+`.*`+regexp.QuoteMeta("THEN $10::numeric")+`.*`+regexp.QuoteMeta("AND credentials = $5::jsonb")+`.*`+regexp.QuoteMeta("AND proxy_id IS NOT DISTINCT FROM $6")+`.*`+regexp.QuoteMeta("COALESCE(extra -> 'upstream_billing_probe', 'null'::jsonb) = $7::jsonb")).
		WithArgs(sqlmock.AnyArg(), int64(17), service.PlatformOpenAI, service.AccountTypeAPIKey, `{"api_key":"sk-test"placeholder`, nil, "null", "true", "true", 0.065).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO scheduler_outbox")).
		WithArgs(service.SchedulerOutboxEventAccountChanged, int64(17), nil, nil, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := newAccountRepositoryWithSQL(client, db, nil)
	account := &service.Account{
		ID:          17,
		Platform:    service.PlatformOpenAI,
		Type:        service.AccountTypeAPIKey,
placeholder"api_key": "sk-test"placeholder,
		Extra: map[string]any{
			service.UpstreamBillingProbeEnabledExtraKey:    true,
			service.UpstreamBillingRateSyncEnabledExtraKey: true,
	placeholder,
placeholder
	rateMultiplier := 0.065

	err = repo.UpdateUpstreamBillingProbeSnapshot(
		context.Background(),
		account,
		&service.UpstreamBillingProbeSnapshot{Status: service.UpstreamBillingProbeStatusOKplaceholder,
		&rateMultiplier,
	)

placeholder
	require.NoError(t, mock.ExpectationsWereMet())
placeholder

func TestUpdateUpstreamBillingProbeSnapshotRejectsChangedProxyIdentity(t *testing.T) {
	db, mock, err := sqlmock.New()
placeholder
	t.Cleanup(func() { _ = db.Close() placeholder)
	client := dbent.NewClient(dbent.Driver(entsql.OpenDB(dialect.Postgres, db)))
	t.Cleanup(func() { _ = client.Close() placeholder)

	mock.ExpectBegin()
	tx, err := client.Tx(context.Background())
placeholder
	mock.ExpectQuery(`(?s)` + regexp.QuoteMeta("SELECT protocol, host, port") + `.*` + regexp.QuoteMeta("FOR SHARE")).
		WithArgs(int64(9)).
		WillReturnRows(sqlmock.NewRows([]string{"protocol", "host", "port", "username", "password", "status"placeholder).
			AddRow("http", "new.example", 3128, "user", "pass", service.StatusActive))

	proxyID := int64(9)
	account := &service.Account{
		ID:          17,
		Platform:    service.PlatformOpenAI,
		Type:        service.AccountTypeAPIKey,
placeholder"api_key": "sk-test"placeholder,
		ProxyID:     &proxyID,
		Proxy: &service.Proxy{
			ID: proxyID, Protocol: "http", Host: "old.example", Port: 3128,
			Username: "user", Password: "pass", Status: service.StatusActive,
	placeholder,
placeholder
	repo := newAccountRepositoryWithSQL(client, db, nil)
	err = repo.UpdateUpstreamBillingProbeSnapshot(dbent.NewTxContext(context.Background(), tx), account, &service.UpstreamBillingProbeSnapshot{Status: service.UpstreamBillingProbeStatusOKplaceholder, nil)

	require.ErrorIs(t, err, service.ErrUpstreamBillingProbeIdentityChanged)
	mock.ExpectRollback()
	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
placeholder

func TestUpdateUpstreamBillingProbeSnapshotRollsBackWhenOutboxFails(t *testing.T) {
	db, mock, err := sqlmock.New()
placeholder
	t.Cleanup(func() { _ = db.Close() placeholder)
	driver := entsql.OpenDB(dialect.Postgres, db)
	client := dbent.NewClient(dbent.Driver(driver))
	t.Cleanup(func() { _ = client.Close() placeholder)

	mock.ExpectBegin()
	mock.ExpectExec(`(?s)`+regexp.QuoteMeta("UPDATE accounts")+`.*`+regexp.QuoteMeta("AND proxy_id IS NOT DISTINCT FROM $6")+`.*`+regexp.QuoteMeta("COALESCE(extra -> 'upstream_billing_probe', 'null'::jsonb) = $7::jsonb")).
		WithArgs(sqlmock.AnyArg(), int64(18), service.PlatformOpenAI, service.AccountTypeAPIKey, `{"api_key":"sk-test"placeholder`, nil, "null", "true", "true", 0.7).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO scheduler_outbox")).WillReturnError(errors.New("outbox failed"))
	mock.ExpectRollback()

	repo := newAccountRepositoryWithSQL(client, db, nil)
	account := &service.Account{
		ID:          18,
		Platform:    service.PlatformOpenAI,
		Type:        service.AccountTypeAPIKey,
placeholder"api_key": "sk-test"placeholder,
		Extra: map[string]any{
			service.UpstreamBillingProbeEnabledExtraKey:    true,
			service.UpstreamBillingRateSyncEnabledExtraKey: true,
	placeholder,
placeholder
	rateMultiplier := 0.7

	err = repo.UpdateUpstreamBillingProbeSnapshot(
		context.Background(),
		account,
		&service.UpstreamBillingProbeSnapshot{Status: service.UpstreamBillingProbeStatusOKplaceholder,
		&rateMultiplier,
	)

	require.EqualError(t, err, "outbox failed")
	require.NoError(t, mock.ExpectationsWereMet())
placeholder
