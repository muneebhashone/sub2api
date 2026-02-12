package repository

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/enttest"
	"github.com/Wei-Shaw/sub2api/ent/securitysecret"
	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "modernc.org/sqlite"
)

func newSecuritySecretTestClient(t *testing.T) *dbent.Client {
placeholder
	name := strings.ReplaceAll(t.Name(), "/", "_")
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", name)

	db, err := sql.Open("sqlite", dsn)
placeholder
	t.Cleanup(func() { _ = db.Close() placeholder)

	_, err = db.Exec("PRAGMA foreign_keys = ON")
placeholder

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := enttest.NewClient(t, enttest.WithOptions(dbent.Driver(drv)))
	t.Cleanup(func() { _ = client.Close() placeholder)
	return client
placeholder

func TestEnsureBootstrapSecretsNilInputs(t *testing.T) {
	err := ensureBootstrapSecrets(context.Background(), nil, &config.Config{placeholder)
placeholder
	require.Contains(t, err.Error(), "nil ent client")

	client := newSecuritySecretTestClient(t)
	err = ensureBootstrapSecrets(context.Background(), client, nil)
placeholder
	require.Contains(t, err.Error(), "nil config")
placeholder

func TestEnsureBootstrapSecretsGenerateAndPersistJWTSecret(t *testing.T) {
	client := newSecuritySecretTestClient(t)
	cfg := &config.Config{placeholder

	err := ensureBootstrapSecrets(context.Background(), client, cfg)
placeholder
	require.NotEmpty(t, cfg.JWT.Secret)
	require.GreaterOrEqual(t, len([]byte(cfg.JWT.Secret)), 32)

	stored, err := client.SecuritySecret.Query().Where(securitysecret.KeyEQ(securitySecretKeyJWT)).Only(context.Background())
placeholder
	require.Equal(t, cfg.JWT.Secret, stored.Value)
placeholder

func TestEnsureBootstrapSecretsLoadExistingJWTSecret(t *testing.T) {
	client := newSecuritySecretTestClient(t)
	_, err := client.SecuritySecret.Create().SetKey(securitySecretKeyJWT).SetValue("existing-jwt-secret-32bytes-long!!!!").Save(context.Background())
placeholder

	cfg := &config.Config{placeholder
	err = ensureBootstrapSecrets(context.Background(), client, cfg)
placeholder
	require.Equal(t, "existing-jwt-secret-32bytes-long!!!!", cfg.JWT.Secret)
placeholder

func TestEnsureBootstrapSecretsRejectInvalidStoredSecret(t *testing.T) {
	client := newSecuritySecretTestClient(t)
	_, err := client.SecuritySecret.Create().SetKey(securitySecretKeyJWT).SetValue("too-short").Save(context.Background())
placeholder

	cfg := &config.Config{placeholder
	err = ensureBootstrapSecrets(context.Background(), client, cfg)
placeholder
	require.Contains(t, err.Error(), "at least 32 bytes")
placeholder

func TestEnsureBootstrapSecretsPersistConfiguredJWTSecret(t *testing.T) {
	client := newSecuritySecretTestClient(t)
	cfg := &config.Config{
		JWT: config.JWTConfig{Secret: "configured-jwt-secret-32bytes-long!!"placeholder,
placeholder

	err := ensureBootstrapSecrets(context.Background(), client, cfg)
placeholder

	stored, err := client.SecuritySecret.Query().Where(securitysecret.KeyEQ(securitySecretKeyJWT)).Only(context.Background())
placeholder
	require.Equal(t, "configured-jwt-secret-32bytes-long!!", stored.Value)
placeholder

func TestEnsureBootstrapSecretsConfiguredSecretTooShort(t *testing.T) {
	client := newSecuritySecretTestClient(t)
	cfg := &config.Config{JWT: config.JWTConfig{Secret: "short"placeholderplaceholder

	err := ensureBootstrapSecrets(context.Background(), client, cfg)
placeholder
	require.Contains(t, err.Error(), "at least 32 bytes")
placeholder

func TestEnsureBootstrapSecretsConfiguredSecretDuplicateIgnored(t *testing.T) {
	client := newSecuritySecretTestClient(t)
	_, err := client.SecuritySecret.Create().
		SetKey(securitySecretKeyJWT).
		SetValue("existing-jwt-secret-32bytes-long!!!!").
		Save(context.Background())
placeholder

	cfg := &config.Config{JWT: config.JWTConfig{Secret: "another-configured-jwt-secret-32!!!!"placeholderplaceholder
	err = ensureBootstrapSecrets(context.Background(), client, cfg)
placeholder

	stored, err := client.SecuritySecret.Query().Where(securitysecret.KeyEQ(securitySecretKeyJWT)).Only(context.Background())
placeholder
	require.Equal(t, "existing-jwt-secret-32bytes-long!!!!", stored.Value)
placeholder

func TestGetOrCreateGeneratedSecuritySecretTrimmedExistingValue(t *testing.T) {
	client := newSecuritySecretTestClient(t)
	_, err := client.SecuritySecret.Create().
		SetKey("trimmed_key").
		SetValue("  existing-trimmed-secret-32bytes-long!!  ").
		Save(context.Background())
placeholder

	value, created, err := getOrCreateGeneratedSecuritySecret(context.Background(), client, "trimmed_key", 32)
placeholder
	require.False(t, created)
	require.Equal(t, "existing-trimmed-secret-32bytes-long!!", value)
placeholder

func TestGetOrCreateGeneratedSecuritySecretQueryError(t *testing.T) {
	client := newSecuritySecretTestClient(t)
	require.NoError(t, client.Close())

	_, _, err := getOrCreateGeneratedSecuritySecret(context.Background(), client, "closed_client_key", 32)
placeholder
placeholder

func TestGetOrCreateGeneratedSecuritySecretCreateValidationError(t *testing.T) {
	client := newSecuritySecretTestClient(t)
	tooLongKey := strings.Repeat("k", 101)

	_, _, err := getOrCreateGeneratedSecuritySecret(context.Background(), client, tooLongKey, 32)
placeholder
placeholder

func TestGetOrCreateGeneratedSecuritySecretConcurrentCreation(t *testing.T) {
	client := newSecuritySecretTestClient(t)
	const goroutines = 8
	key := "concurrent_bootstrap_key"

	values := make([]string, goroutines)
	createdFlags := make([]bool, goroutines)
	errs := make([]error, goroutines)

	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			values[idx], createdFlags[idx], errs[idx] = getOrCreateGeneratedSecuritySecret(context.Background(), client, key, 32)
	placeholder(i)
placeholder
	wg.Wait()

	for i := range errs {
		require.NoError(t, errs[i])
		require.NotEmpty(t, values[i])
placeholder
	for i := 1; i < len(values); i++ {
		require.Equal(t, values[0], values[i])
placeholder

	createdCount := 0
	for _, created := range createdFlags {
		if created {
			createdCount++
	placeholder
placeholder
	require.GreaterOrEqual(t, createdCount, 1)
	require.LessOrEqual(t, createdCount, 1)

	count, err := client.SecuritySecret.Query().Where(securitysecret.KeyEQ(key)).Count(context.Background())
placeholder
	require.Equal(t, 1, count)
placeholder

func TestGetOrCreateGeneratedSecuritySecretGenerateError(t *testing.T) {
	client := newSecuritySecretTestClient(t)
	originalRead := readRandomBytes
	readRandomBytes = func([]byte) (int, error) {
		return 0, errors.New("boom")
placeholder
	t.Cleanup(func() {
		readRandomBytes = originalRead
placeholder)

	_, _, err := getOrCreateGeneratedSecuritySecret(context.Background(), client, "gen_error_key", 32)
placeholder
	require.Contains(t, err.Error(), "boom")
placeholder

func TestCreateSecuritySecretIfAbsent(t *testing.T) {
	client := newSecuritySecretTestClient(t)

	err := createSecuritySecretIfAbsent(context.Background(), client, "abc", "short")
placeholder
	require.Contains(t, err.Error(), "at least 32 bytes")

	err = createSecuritySecretIfAbsent(context.Background(), client, "abc", "valid-jwt-secret-value-32bytes-long")
placeholder

	err = createSecuritySecretIfAbsent(context.Background(), client, "abc", "another-valid-secret-value-32bytes")
placeholder

	count, err := client.SecuritySecret.Query().Where(securitysecret.KeyEQ("abc")).Count(context.Background())
placeholder
	require.Equal(t, 1, count)
placeholder

func TestCreateSecuritySecretIfAbsentValidationError(t *testing.T) {
	client := newSecuritySecretTestClient(t)
	err := createSecuritySecretIfAbsent(
		context.Background(),
		client,
		strings.Repeat("k", 101),
		"valid-jwt-secret-value-32bytes-long",
	)
placeholder
placeholder

func TestGenerateHexSecretReadError(t *testing.T) {
	originalRead := readRandomBytes
	readRandomBytes = func([]byte) (int, error) {
		return 0, errors.New("read random failed")
placeholder
	t.Cleanup(func() {
		readRandomBytes = originalRead
placeholder)

	_, err := generateHexSecret(32)
placeholder
	require.Contains(t, err.Error(), "read random failed")
placeholder

func TestGenerateHexSecretLengths(t *testing.T) {
	v1, err := generateHexSecret(0)
placeholder
	require.Len(t, v1, 64)
	_, err = hex.DecodeString(v1)
placeholder

	v2, err := generateHexSecret(16)
placeholder
	require.Len(t, v2, 32)
	_, err = hex.DecodeString(v2)
placeholder

	require.NotEqual(t, v1, v2)
placeholder
