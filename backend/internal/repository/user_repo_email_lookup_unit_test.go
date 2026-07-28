package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/enttest"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "modernc.org/sqlite"
)

func newUserEntRepo(t *testing.T) (*userRepository, *dbent.Client) {
placeholder

	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()))
placeholder
	t.Cleanup(func() { _ = db.Close() placeholder)
	db.SetMaxOpenConns(10)

	_, err = db.Exec("PRAGMA foreign_keys = ON")
placeholder

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := enttest.NewClient(t, enttest.WithOptions(dbent.Driver(drv)))
	t.Cleanup(func() { _ = client.Close() placeholder)

	return newUserRepositoryWithSQL(client, db), client
placeholder

func TestUserRepositoryGetByEmailNormalizesLegacySpacingAndCase(t *testing.T) {
	repo, _ := newUserEntRepo(t)
	ctx := context.Background()

	err := repo.Create(ctx, &service.User{
		Email:        " Legacy@Example.com ",
		Username:     "legacy-user",
		PasswordHash: "hash",
		Role:         service.RoleUser,
		Status:       service.StatusActive,
placeholder)
placeholder

	got, err := repo.GetByEmail(ctx, "legacy@example.com")
placeholder
	require.Equal(t, " Legacy@Example.com ", got.Email)
placeholder

func TestUserRepositoryExistsByEmailNormalizesLegacySpacingAndCase(t *testing.T) {
	repo, _ := newUserEntRepo(t)
	ctx := context.Background()

	err := repo.Create(ctx, &service.User{
		Email:        " Legacy@Example.com ",
		Username:     "legacy-user",
		PasswordHash: "hash",
		Role:         service.RoleUser,
		Status:       service.StatusActive,
placeholder)
placeholder

	exists, err := repo.ExistsByEmail(ctx, "  LEGACY@example.com  ")
placeholder
	require.True(t, exists)
placeholder

func TestUserRepositoryCreateRejectsNormalizedEmailDuplicate(t *testing.T) {
	repo, _ := newUserEntRepo(t)
	ctx := context.Background()

	err := repo.Create(ctx, &service.User{
		Email:        " Existing@Example.com ",
		Username:     "existing-user",
		PasswordHash: "hash",
		Role:         service.RoleUser,
		Status:       service.StatusActive,
placeholder)
placeholder

	err = repo.Create(ctx, &service.User{
		Email:        "existing@example.com",
		Username:     "duplicate-user",
		PasswordHash: "hash",
		Role:         service.RoleUser,
		Status:       service.StatusActive,
placeholder)
	require.ErrorIs(t, err, service.ErrEmailExists)
placeholder

func TestUserRepositoryUpdateRejectsNormalizedEmailDuplicate(t *testing.T) {
	repo, _ := newUserEntRepo(t)
	ctx := context.Background()

	first := &service.User{
		Email:        " Existing@Example.com ",
		Username:     "existing-user",
		PasswordHash: "hash",
		Role:         service.RoleUser,
		Status:       service.StatusActive,
placeholder
	require.NoError(t, repo.Create(ctx, first))

	second := &service.User{
		Email:        "second@example.com",
		Username:     "second-user",
		PasswordHash: "hash",
		Role:         service.RoleUser,
		Status:       service.StatusActive,
placeholder
	require.NoError(t, repo.Create(ctx, second))

	second.Email = " existing@example.com "
	err := repo.Update(ctx, second, service.UserUpdateFields{Email: trueplaceholder)
	require.ErrorIs(t, err, service.ErrEmailExists)
placeholder

func TestUserRepositoryGetByEmailReportsNormalizedEmailConflict(t *testing.T) {
	repo, client := newUserEntRepo(t)
	ctx := context.Background()

	_, err := client.User.Create().
		SetEmail("Conflict@Example.com").
		SetUsername("conflict-user-1").
		SetPasswordHash("hash").
		SetRole(service.RoleUser).
		SetStatus(service.StatusActive).
		Save(ctx)
placeholder

	_, err = client.User.Create().
		SetEmail(" conflict@example.com ").
		SetUsername("conflict-user-2").
		SetPasswordHash("hash").
		SetRole(service.RoleUser).
		SetStatus(service.StatusActive).
		Save(ctx)
placeholder

	_, err = repo.GetByEmail(ctx, "conflict@example.com")
placeholder
	require.ErrorContains(t, err, "normalized email lookup matched multiple users")
placeholder

func TestUserRepositoryCreateSerializesNormalizedEmailConflictsUnderConcurrency(t *testing.T) {
	repo, client := newUserEntRepo(t)
	ctx := context.Background()

	firstCreateStarted := make(chan struct{placeholder)
	releaseFirstCreate := make(chan struct{placeholder)
	var firstCreate sync.Once
	client.User.Use(func(next dbent.Mutator) dbent.Mutator {
		return dbent.MutateFunc(func(ctx context.Context, m dbent.Mutation) (dbent.Value, error) {
			blocked := false
			if m.Op().Is(dbent.OpCreate) {
				firstCreate.Do(func() {
					blocked = true
					close(firstCreateStarted)
			placeholder)
		placeholder
			if blocked {
				<-releaseFirstCreate
		placeholder
			return next.Mutate(ctx, m)
	placeholder)
placeholder)

	type createResult struct {
		err error
placeholder

	results := make(chan createResult, 2)
	go func() {
		results <- createResult{err: repo.Create(ctx, &service.User{
			Email:        " Race@Example.com ",
			Username:     "race-user-1",
			PasswordHash: "hash",
			Role:         service.RoleUser,
			Status:       service.StatusActive,
	placeholder)placeholder
placeholder()

	<-firstCreateStarted

	go func() {
		results <- createResult{err: repo.Create(ctx, &service.User{
			Email:        "race@example.com",
			Username:     "race-user-2",
			PasswordHash: "hash",
			Role:         service.RoleUser,
			Status:       service.StatusActive,
	placeholder)placeholder
placeholder()

	time.Sleep(100 * time.Millisecond)
	close(releaseFirstCreate)

	first := <-results
	second := <-results

	errors := []error{first.err, second.errplaceholder
	successes := 0
	conflicts := 0
	for _, err := range errors {
		switch err {
		case nil:
			successes++
		case service.ErrEmailExists:
			conflicts++
		default:
			t.Fatalf("unexpected create error: %v", err)
	placeholder
placeholder
	require.Equal(t, 1, successes)
	require.Equal(t, 1, conflicts)

	count, err := client.User.Query().Where(userEmailLookupPredicate("race@example.com")).Count(ctx)
placeholder
	require.Equal(t, 1, count)
placeholder
