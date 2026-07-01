//go:build integration

package repository

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

func TestAccountRepoSparkShadowRoundTrip(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	repo := newAccountRepositoryWithSQL(tx.Client(), tx, nil)

	parent := &service.Account{
		Name:     "parent",
		Platform: service.PlatformOpenAI,
		Type:     service.AccountTypeOAuth,
		Status:   service.StatusActive,
placeholder
	if err := repo.Create(ctx, parent); err != nil {
		t.Fatalf("create parent: %v", err)
placeholder
	pid := parent.ID
	shadow := &service.Account{
		Name:            "shadow",
		Platform:        service.PlatformOpenAI,
		Type:            service.AccountTypeOAuth,
		Status:          service.StatusActive,
		ParentAccountID: &pid,
		QuotaDimension:  service.QuotaDimensionSpark,
placeholder
	if err := repo.Create(ctx, shadow); err != nil {
		t.Fatalf("create shadow: %v", err)
placeholder
	got, err := repo.GetByID(ctx, shadow.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
placeholder
	if got.ParentAccountID == nil || *got.ParentAccountID != pid {
		t.Fatalf("ParentAccountID round-trip: %v", got.ParentAccountID)
placeholder
	if got.QuotaDimension != service.QuotaDimensionSpark {
		t.Fatalf("QuotaDimension: %q", got.QuotaDimension)
placeholder
placeholder

func TestListShadowsByParent(t *testing.T) {
	// Schema enforces at most one spark shadow per parent (uq_accounts_spark_shadow_per_parent).
	// Test strategy: create 2 parents each with 1 spark shadow + 1 unrelated account;
	// assert ListShadowsByParent(parent1.ID) returns exactly 1 (filtering by both
	// parent_account_id and quota_dimension='spark', excluding parent2's shadow and unrelated).
	ctx := context.Background()
	tx := testEntTx(t)
	repo := newAccountRepositoryWithSQL(tx.Client(), tx, nil)

	// Create parent1 and its spark shadow
	parent1 := &service.Account{
		Name:     "list-parent1",
		Platform: service.PlatformOpenAI,
		Type:     service.AccountTypeOAuth,
		Status:   service.StatusActive,
placeholder
	if err := repo.Create(ctx, parent1); err != nil {
		t.Fatalf("create parent1: %v", err)
placeholder
	pid1 := parent1.ID

	shadow1 := &service.Account{
		Name:            "shadow1",
		Platform:        service.PlatformOpenAI,
		Type:            service.AccountTypeOAuth,
		Status:          service.StatusActive,
		ParentAccountID: &pid1,
		QuotaDimension:  service.QuotaDimensionSpark,
placeholder
	if err := repo.Create(ctx, shadow1); err != nil {
		t.Fatalf("create shadow1: %v", err)
placeholder

	// Create parent2 and its spark shadow (must NOT appear in parent1's list)
	parent2 := &service.Account{
		Name:     "list-parent2",
		Platform: service.PlatformOpenAI,
		Type:     service.AccountTypeOAuth,
		Status:   service.StatusActive,
placeholder
	if err := repo.Create(ctx, parent2); err != nil {
		t.Fatalf("create parent2: %v", err)
placeholder
	pid2 := parent2.ID

	shadow2 := &service.Account{
		Name:            "shadow2",
		Platform:        service.PlatformOpenAI,
		Type:            service.AccountTypeOAuth,
		Status:          service.StatusActive,
		ParentAccountID: &pid2,
		QuotaDimension:  service.QuotaDimensionSpark,
placeholder
	if err := repo.Create(ctx, shadow2); err != nil {
		t.Fatalf("create shadow2: %v", err)
placeholder

	// Create 1 unrelated normal account (no parent, global dimension)
	unrelated := &service.Account{
		Name:     "unrelated",
		Platform: service.PlatformOpenAI,
		Type:     service.AccountTypeOAuth,
		Status:   service.StatusActive,
placeholder
	if err := repo.Create(ctx, unrelated); err != nil {
		t.Fatalf("create unrelated: %v", err)
placeholder

	// Assert ListShadowsByParent returns exactly 1 for parent1
	got, err := repo.ListShadowsByParent(ctx, pid1)
	if err != nil {
		t.Fatalf("ListShadowsByParent: %v", err)
placeholder
	if len(got) != 1 {
		t.Fatalf("expected 1 spark shadow for parent1, got %d", len(got))
placeholder
	acc := got[0]
	if acc.ParentAccountID == nil || *acc.ParentAccountID != pid1 {
		t.Errorf("unexpected ParentAccountID: %v", acc.ParentAccountID)
placeholder
	if acc.QuotaDimension != service.QuotaDimensionSpark {
		t.Errorf("unexpected QuotaDimension: %q", acc.QuotaDimension)
placeholder
	if acc.ID != shadow1.ID {
		t.Errorf("expected shadow1.ID=%d, got %d", shadow1.ID, acc.ID)
placeholder
placeholder
