package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestApplyMigrations_NilDB(t *testing.T) {
	err := ApplyMigrations(context.Background(), nil)
placeholder
	require.Contains(t, err.Error(), "nil sql db")
placeholder

func TestApplyMigrations_DelegatesToApplyMigrationsFS(t *testing.T) {
	db, mock, err := sqlmock.New()
placeholder
	defer func() { _ = db.Close() placeholder()

	mock.ExpectQuery("SELECT pg_try_advisory_lock\\(\\$1\\)").
		WithArgs(migrationsAdvisoryLockID).
		WillReturnError(errors.New("lock failed"))

	err = ApplyMigrations(context.Background(), db)
placeholder
	require.Contains(t, err.Error(), "acquire migrations lock")
	require.NoError(t, mock.ExpectationsWereMet())
placeholder

func TestLatestMigrationBaseline(t *testing.T) {
	t.Run("empty_fs_returns_baseline", func(t *testing.T) {
		version, description, hash, err := latestMigrationBaseline(fstest.MapFS{placeholder)
	placeholder
		require.Equal(t, "baseline", version)
		require.Equal(t, "baseline", description)
		require.Equal(t, "", hash)
placeholder)

	t.Run("uses_latest_sorted_sql_file", func(t *testing.T) {
		fsys := fstest.MapFS{
			"001_init.sql": &fstest.MapFile{Data: []byte("CREATE TABLE t1(id int);")placeholder,
			"010_final.sql": &fstest.MapFile{
				Data: []byte("CREATE TABLE t2(id int);"),
		placeholder,
	placeholder
		version, description, hash, err := latestMigrationBaseline(fsys)
	placeholder
		require.Equal(t, "010_final", version)
		require.Equal(t, "010_final", description)
		require.Len(t, hash, 64)
placeholder)

	t.Run("read_file_error", func(t *testing.T) {
		fsys := fstest.MapFS{
			"010_bad.sql": &fstest.MapFile{Mode: fs.ModeDirplaceholder,
	placeholder
		_, _, _, err := latestMigrationBaseline(fsys)
	placeholder
placeholder)
placeholder

func TestIsMigrationChecksumCompatible_AdditionalCases(t *testing.T) {
	require.False(t, isMigrationChecksumCompatible("unknown.sql", "db", "file"))

	var (
		name string
		rule migrationChecksumCompatibilityRule
	)
	for n, r := range migrationChecksumCompatibilityRules {
		name = n
		rule = r
		break
placeholder
	require.NotEmpty(t, name)

	require.False(t, isMigrationChecksumCompatible(name, "db-not-accepted", "file-not-match"))
	require.False(t, isMigrationChecksumCompatible(name, "db-not-accepted", rule.fileChecksum))

	var accepted string
	for checksum := range rule.acceptedDBChecksum {
		accepted = checksum
		break
placeholder
	require.NotEmpty(t, accepted)
	require.True(t, isMigrationChecksumCompatible(name, accepted, rule.fileChecksum))
placeholder

func TestEnsureAtlasBaselineAligned(t *testing.T) {
	t.Run("skip_when_no_legacy_table", func(t *testing.T) {
		db, mock, err := sqlmock.New()
	placeholder
		defer func() { _ = db.Close() placeholder()

		mock.ExpectQuery("SELECT EXISTS \\(").
			WithArgs("schema_migrations").
			WillReturnRows(sqlmock.NewRows([]string{"exists"placeholder).AddRow(false))

		err = ensureAtlasBaselineAligned(context.Background(), db, fstest.MapFS{placeholder)
	placeholder
		require.NoError(t, mock.ExpectationsWereMet())
placeholder)

	t.Run("create_atlas_and_insert_baseline_when_empty", func(t *testing.T) {
		db, mock, err := sqlmock.New()
	placeholder
		defer func() { _ = db.Close() placeholder()

		mock.ExpectQuery("SELECT EXISTS \\(").
			WithArgs("schema_migrations").
			WillReturnRows(sqlmock.NewRows([]string{"exists"placeholder).AddRow(true))
		mock.ExpectQuery("SELECT EXISTS \\(").
			WithArgs("atlas_schema_revisions").
			WillReturnRows(sqlmock.NewRows([]string{"exists"placeholder).AddRow(false))
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS atlas_schema_revisions").
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM atlas_schema_revisions").
			WillReturnRows(sqlmock.NewRows([]string{"count"placeholder).AddRow(0))
		mock.ExpectExec("INSERT INTO atlas_schema_revisions").
			WithArgs("002_next", "002_next", 1, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		fsys := fstest.MapFS{
			"001_init.sql": &fstest.MapFile{Data: []byte("CREATE TABLE t1(id int);")placeholder,
			"002_next.sql": &fstest.MapFile{Data: []byte("CREATE TABLE t2(id int);")placeholder,
	placeholder
		err = ensureAtlasBaselineAligned(context.Background(), db, fsys)
	placeholder
		require.NoError(t, mock.ExpectationsWereMet())
placeholder)

	t.Run("error_when_checking_legacy_table", func(t *testing.T) {
		db, mock, err := sqlmock.New()
	placeholder
		defer func() { _ = db.Close() placeholder()

		mock.ExpectQuery("SELECT EXISTS \\(").
			WithArgs("schema_migrations").
			WillReturnError(errors.New("exists failed"))

		err = ensureAtlasBaselineAligned(context.Background(), db, fstest.MapFS{placeholder)
	placeholder
		require.Contains(t, err.Error(), "check schema_migrations")
		require.NoError(t, mock.ExpectationsWereMet())
placeholder)

	t.Run("error_when_counting_atlas_rows", func(t *testing.T) {
		db, mock, err := sqlmock.New()
	placeholder
		defer func() { _ = db.Close() placeholder()

		mock.ExpectQuery("SELECT EXISTS \\(").
			WithArgs("schema_migrations").
			WillReturnRows(sqlmock.NewRows([]string{"exists"placeholder).AddRow(true))
		mock.ExpectQuery("SELECT EXISTS \\(").
			WithArgs("atlas_schema_revisions").
			WillReturnRows(sqlmock.NewRows([]string{"exists"placeholder).AddRow(true))
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM atlas_schema_revisions").
			WillReturnError(errors.New("count failed"))

		err = ensureAtlasBaselineAligned(context.Background(), db, fstest.MapFS{placeholder)
	placeholder
		require.Contains(t, err.Error(), "count atlas_schema_revisions")
		require.NoError(t, mock.ExpectationsWereMet())
placeholder)

	t.Run("error_when_creating_atlas_table", func(t *testing.T) {
		db, mock, err := sqlmock.New()
	placeholder
		defer func() { _ = db.Close() placeholder()

		mock.ExpectQuery("SELECT EXISTS \\(").
			WithArgs("schema_migrations").
			WillReturnRows(sqlmock.NewRows([]string{"exists"placeholder).AddRow(true))
		mock.ExpectQuery("SELECT EXISTS \\(").
			WithArgs("atlas_schema_revisions").
			WillReturnRows(sqlmock.NewRows([]string{"exists"placeholder).AddRow(false))
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS atlas_schema_revisions").
			WillReturnError(errors.New("create failed"))

		err = ensureAtlasBaselineAligned(context.Background(), db, fstest.MapFS{placeholder)
	placeholder
		require.Contains(t, err.Error(), "create atlas_schema_revisions")
		require.NoError(t, mock.ExpectationsWereMet())
placeholder)

	t.Run("error_when_inserting_baseline", func(t *testing.T) {
		db, mock, err := sqlmock.New()
	placeholder
		defer func() { _ = db.Close() placeholder()

		mock.ExpectQuery("SELECT EXISTS \\(").
			WithArgs("schema_migrations").
			WillReturnRows(sqlmock.NewRows([]string{"exists"placeholder).AddRow(true))
		mock.ExpectQuery("SELECT EXISTS \\(").
			WithArgs("atlas_schema_revisions").
			WillReturnRows(sqlmock.NewRows([]string{"exists"placeholder).AddRow(true))
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM atlas_schema_revisions").
			WillReturnRows(sqlmock.NewRows([]string{"count"placeholder).AddRow(0))
		mock.ExpectExec("INSERT INTO atlas_schema_revisions").
			WithArgs("001_init", "001_init", 1, sqlmock.AnyArg()).
			WillReturnError(errors.New("insert failed"))

		fsys := fstest.MapFS{
			"001_init.sql": &fstest.MapFile{Data: []byte("CREATE TABLE t(id int);")placeholder,
	placeholder
		err = ensureAtlasBaselineAligned(context.Background(), db, fsys)
	placeholder
		require.Contains(t, err.Error(), "insert atlas baseline")
		require.NoError(t, mock.ExpectationsWereMet())
placeholder)
placeholder

func TestApplyMigrationsFS_ChecksumMismatchRejected(t *testing.T) {
	db, mock, err := sqlmock.New()
placeholder
	defer func() { _ = db.Close() placeholder()

	prepareMigrationsBootstrapExpectations(mock)
	mock.ExpectQuery("SELECT checksum FROM schema_migrations WHERE filename = \\$1").
		WithArgs("001_init.sql").
		WillReturnRows(sqlmock.NewRows([]string{"checksum"placeholder).AddRow("mismatched-checksum"))
	mock.ExpectExec("SELECT pg_advisory_unlock\\(\\$1\\)").
		WithArgs(migrationsAdvisoryLockID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	fsys := fstest.MapFS{
		"001_init.sql": &fstest.MapFile{Data: []byte("CREATE TABLE t(id int);")placeholder,
placeholder
	err = applyMigrationsFS(context.Background(), db, fsys)
placeholder
	require.Contains(t, err.Error(), "checksum mismatch")
	require.NoError(t, mock.ExpectationsWereMet())
placeholder

func TestApplyMigrationsFS_CheckMigrationQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
placeholder
	defer func() { _ = db.Close() placeholder()

	prepareMigrationsBootstrapExpectations(mock)
	mock.ExpectQuery("SELECT checksum FROM schema_migrations WHERE filename = \\$1").
		WithArgs("001_err.sql").
		WillReturnError(errors.New("query failed"))
	mock.ExpectExec("SELECT pg_advisory_unlock\\(\\$1\\)").
		WithArgs(migrationsAdvisoryLockID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	fsys := fstest.MapFS{
		"001_err.sql": &fstest.MapFile{Data: []byte("SELECT 1;")placeholder,
placeholder
	err = applyMigrationsFS(context.Background(), db, fsys)
placeholder
	require.Contains(t, err.Error(), "check migration 001_err.sql")
	require.NoError(t, mock.ExpectationsWereMet())
placeholder

func TestApplyMigrationsFS_SkipEmptyAndAlreadyApplied(t *testing.T) {
	db, mock, err := sqlmock.New()
placeholder
	defer func() { _ = db.Close() placeholder()

	prepareMigrationsBootstrapExpectations(mock)

	alreadySQL := "CREATE TABLE t(id int);"
	checksum := migrationChecksum(alreadySQL)
	mock.ExpectQuery("SELECT checksum FROM schema_migrations WHERE filename = \\$1").
		WithArgs("001_already.sql").
		WillReturnRows(sqlmock.NewRows([]string{"checksum"placeholder).AddRow(checksum))
	mock.ExpectExec("SELECT pg_advisory_unlock\\(\\$1\\)").
		WithArgs(migrationsAdvisoryLockID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	fsys := fstest.MapFS{
		"000_empty.sql":   &fstest.MapFile{Data: []byte("   \n\t ")placeholder,
		"001_already.sql": &fstest.MapFile{Data: []byte(alreadySQL)placeholder,
placeholder
	err = applyMigrationsFS(context.Background(), db, fsys)
placeholder
	require.NoError(t, mock.ExpectationsWereMet())
placeholder

func TestApplyMigrationsFS_ReadMigrationError(t *testing.T) {
	db, mock, err := sqlmock.New()
placeholder
	defer func() { _ = db.Close() placeholder()

	prepareMigrationsBootstrapExpectations(mock)
	mock.ExpectExec("SELECT pg_advisory_unlock\\(\\$1\\)").
		WithArgs(migrationsAdvisoryLockID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	fsys := fstest.MapFS{
		"001_bad.sql": &fstest.MapFile{Mode: fs.ModeDirplaceholder,
placeholder
	err = applyMigrationsFS(context.Background(), db, fsys)
placeholder
	require.Contains(t, err.Error(), "read migration 001_bad.sql")
	require.NoError(t, mock.ExpectationsWereMet())
placeholder

func TestPgAdvisoryLockAndUnlock_ErrorBranches(t *testing.T) {
	t.Run("context_cancelled_while_not_locked", func(t *testing.T) {
		db, mock, err := sqlmock.New()
	placeholder
		defer func() { _ = db.Close() placeholder()

		mock.ExpectQuery("SELECT pg_try_advisory_lock\\(\\$1\\)").
			WithArgs(migrationsAdvisoryLockID).
			WillReturnRows(sqlmock.NewRows([]string{"pg_try_advisory_lock"placeholder).AddRow(false))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
		defer cancel()
		err = pgAdvisoryLock(ctx, db)
	placeholder
		require.Contains(t, err.Error(), "acquire migrations lock")
		require.NoError(t, mock.ExpectationsWereMet())
placeholder)

	t.Run("unlock_exec_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
	placeholder
		defer func() { _ = db.Close() placeholder()

		mock.ExpectExec("SELECT pg_advisory_unlock\\(\\$1\\)").
			WithArgs(migrationsAdvisoryLockID).
			WillReturnError(errors.New("unlock failed"))

		err = pgAdvisoryUnlock(context.Background(), db)
	placeholder
		require.Contains(t, err.Error(), "release migrations lock")
		require.NoError(t, mock.ExpectationsWereMet())
placeholder)

	t.Run("acquire_lock_after_retry", func(t *testing.T) {
		db, mock, err := sqlmock.New()
	placeholder
		defer func() { _ = db.Close() placeholder()

		mock.ExpectQuery("SELECT pg_try_advisory_lock\\(\\$1\\)").
			WithArgs(migrationsAdvisoryLockID).
			WillReturnRows(sqlmock.NewRows([]string{"pg_try_advisory_lock"placeholder).AddRow(false))
		mock.ExpectQuery("SELECT pg_try_advisory_lock\\(\\$1\\)").
			WithArgs(migrationsAdvisoryLockID).
			WillReturnRows(sqlmock.NewRows([]string{"pg_try_advisory_lock"placeholder).AddRow(true))

		ctx, cancel := context.WithTimeout(context.Background(), migrationsLockRetryInterval*3)
		defer cancel()
		start := time.Now()
		err = pgAdvisoryLock(ctx, db)
	placeholder
		require.GreaterOrEqual(t, time.Since(start), migrationsLockRetryInterval)
		require.NoError(t, mock.ExpectationsWereMet())
placeholder)
placeholder

func migrationChecksum(content string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(content)))
	return hex.EncodeToString(sum[:])
placeholder
