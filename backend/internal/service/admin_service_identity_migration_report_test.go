package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/enttest"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "modernc.org/sqlite"
)

func newAdminServiceMigrationReportTestClient(t *testing.T) *dbent.Client {
placeholder

	db, err := sql.Open("sqlite", "file:admin_service_migration_reports?mode=memory&cache=shared&_fk=1")
placeholder
	t.Cleanup(func() { _ = db.Close() placeholder)

	_, err = db.Exec("PRAGMA foreign_keys = ON")
placeholder

	_, err = db.Exec(`CREATE TABLE auth_identity_migration_reports (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		report_type TEXT NOT NULL,
		report_key TEXT NOT NULL,
		details TEXT NOT NULL DEFAULT '{placeholder',
		created_at DATETIME NOT NULL
	)`)
placeholder

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := enttest.NewClient(t, enttest.WithOptions(dbent.Driver(drv)))
	t.Cleanup(func() { _ = client.Close() placeholder)
	return client
placeholder

func TestAdminServiceListAuthIdentityMigrationReports(t *testing.T) {
	client := newAdminServiceMigrationReportTestClient(t)
	driver, ok := client.Driver().(*entsql.Driver)
	require.True(t, ok)

	now := time.Now().UTC()
	_, err := driver.DB().ExecContext(context.Background(), `
INSERT INTO auth_identity_migration_reports (report_type, report_key, details, created_at)
VALUES
	($1, $2, $3, $4),
	($5, $6, $7, $8)`,
		"oidc_synthetic_email_requires_manual_recovery", "u-1", `{"user_id":1placeholder`, now,
		"wechat_provider_key_conflict", "u-2", `{"user_id":2placeholder`, now.Add(-time.Minute),
	)
placeholder

	svc := &adminServiceImpl{entClient: clientplaceholder
	reports, total, err := svc.ListAuthIdentityMigrationReports(context.Background(), "oidc_synthetic_email_requires_manual_recovery", 1, 20)
placeholder
	require.Equal(t, int64(1), total)
	require.Len(t, reports, 1)
	require.Equal(t, "oidc_synthetic_email_requires_manual_recovery", reports[0].ReportType)
	require.Equal(t, float64(1), reports[0].Details["user_id"])
placeholder

func TestAdminServiceGetAuthIdentityMigrationReportSummary(t *testing.T) {
	client := newAdminServiceMigrationReportTestClient(t)
	driver, ok := client.Driver().(*entsql.Driver)
	require.True(t, ok)

	now := time.Now().UTC()
	_, err := driver.DB().ExecContext(context.Background(), `
INSERT INTO auth_identity_migration_reports (report_type, report_key, details, created_at)
VALUES
	($1, $2, $3, $4),
	($5, $6, $7, $8),
	($9, $10, $11, $12)`,
		"oidc_synthetic_email_requires_manual_recovery", "u-1", `{"user_id":1placeholder`, now,
		"wechat_provider_key_conflict", "u-2", `{"user_id":2placeholder`, now.Add(-time.Minute),
		"wechat_provider_key_conflict", "u-3", `{"user_id":3placeholder`, now.Add(-2*time.Minute),
	)
placeholder

	svc := &adminServiceImpl{entClient: clientplaceholder
	summary, err := svc.GetAuthIdentityMigrationReportSummary(context.Background())
placeholder
	require.Equal(t, int64(3), summary.Total)
	require.Equal(t, int64(1), summary.ByType["oidc_synthetic_email_requires_manual_recovery"])
	require.Equal(t, int64(2), summary.ByType["wechat_provider_key_conflict"])
placeholder
