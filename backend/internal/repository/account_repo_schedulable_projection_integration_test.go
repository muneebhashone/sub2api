//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestListSchedulableAccountLoadsMatchesListSchedulable(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	client := tx.Client()
	repo := newAccountRepositoryWithSQL(client, tx, nil)
	now := time.Now()
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	create := func(name string) *service.Account {
		return mustCreateAccount(t, client, &service.Account{Name: name, Schedulable: trueplaceholder)
placeholder

	positiveLoad := create("projection-positive-load")
	_, err := client.Account.UpdateOneID(positiveLoad.ID).SetConcurrency(2).SetLoadFactor(9).SetPriority(30).Save(ctx)
placeholder
	concurrencyFallback := create("projection-concurrency-fallback")
	_, err = client.Account.UpdateOneID(concurrencyFallback.ID).SetConcurrency(4).SetPriority(10).Save(ctx)
placeholder
	zeroFallback := create("projection-zero-fallback")
	_, err = client.Account.UpdateOneID(zeroFallback.ID).SetConcurrency(0).SetLoadFactor(0).SetPriority(20).Save(ctx)
placeholder

	disabled := create("projection-disabled")
	_, err = client.Account.UpdateOneID(disabled.ID).SetStatus(service.StatusDisabled).Save(ctx)
placeholder
	unschedulable := create("projection-unschedulable")
	_, err = client.Account.UpdateOneID(unschedulable.ID).SetSchedulable(false).Save(ctx)
placeholder
	expired := create("projection-expired")
	_, err = client.Account.UpdateOneID(expired.ID).SetExpiresAt(past).SetAutoPauseOnExpired(true).Save(ctx)
placeholder
	expiredAllowed := create("projection-expired-allowed")
	_, err = client.Account.UpdateOneID(expiredAllowed.ID).SetExpiresAt(past).SetAutoPauseOnExpired(false).Save(ctx)
placeholder
	overloaded := create("projection-overloaded")
	_, err = client.Account.UpdateOneID(overloaded.ID).SetOverloadUntil(future).Save(ctx)
placeholder
	overloadCleared := create("projection-overload-cleared")
	_, err = client.Account.UpdateOneID(overloadCleared.ID).SetOverloadUntil(past).Save(ctx)
placeholder
	rateLimited := create("projection-rate-limited")
	_, err = client.Account.UpdateOneID(rateLimited.ID).SetRateLimitResetAt(future).Save(ctx)
placeholder
	rateLimitCleared := create("projection-rate-limit-cleared")
	_, err = client.Account.UpdateOneID(rateLimitCleared.ID).SetRateLimitResetAt(past).Save(ctx)
placeholder
	tempBlocked := create("projection-temp-blocked")
	_, err = client.Account.UpdateOneID(tempBlocked.ID).SetTempUnschedulableUntil(future).Save(ctx)
placeholder
	tempCleared := create("projection-temp-cleared")
	_, err = client.Account.UpdateOneID(tempCleared.ID).SetTempUnschedulableUntil(past).Save(ctx)
placeholder

	accounts, err := repo.ListSchedulable(ctx)
placeholder
	loads, err := repo.ListSchedulableAccountLoads(ctx)
placeholder

	accountIDs := make([]int64, 0, len(accounts))
	wantByID := make(map[int64]int, len(accounts))
	for i := range accounts {
		accountIDs = append(accountIDs, accounts[i].ID)
		wantByID[accounts[i].ID] = accounts[i].EffectiveLoadFactor()
placeholder

	loadIDs := make([]int64, 0, len(loads))
	byID := make(map[int64]int, len(loads))
	for _, load := range loads {
		loadIDs = append(loadIDs, load.ID)
		byID[load.ID] = load.MaxConcurrency
placeholder
	require.Equal(t, accountIDs, loadIDs)
	targetIDs := map[int64]struct{placeholder{
		positiveLoad.ID: {placeholder, concurrencyFallback.ID: {placeholder, zeroFallback.ID: {placeholder,
placeholder
	targetOrder := make([]int64, 0, len(targetIDs))
	for _, id := range loadIDs {
		if _, ok := targetIDs[id]; ok {
			targetOrder = append(targetOrder, id)
	placeholder
placeholder
	require.Equal(t, []int64{concurrencyFallback.ID, zeroFallback.ID, positiveLoad.IDplaceholder, targetOrder)
	require.Equal(t, wantByID, byID)
	require.Equal(t, 9, byID[positiveLoad.ID])
	require.Equal(t, 4, byID[concurrencyFallback.ID])
	require.Equal(t, 1, byID[zeroFallback.ID])
	for _, included := range []*service.Account{expiredAllowed, overloadCleared, rateLimitCleared, tempClearedplaceholder {
		require.Contains(t, byID, included.ID)
placeholder
	for _, excluded := range []*service.Account{disabled, unschedulable, expired, overloaded, rateLimited, tempBlockedplaceholder {
		require.NotContains(t, byID, excluded.ID)
placeholder
placeholder
