//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/ent/group"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestEnsureSimpleModeDefaultGroups_CreatesMissingDefaults(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	client := tx.Client()

	seedCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	require.NoError(t, ensureSimpleModeDefaultGroups(seedCtx, client))

	assertGroupExists := func(name string) {
		exists, err := client.Group.Query().Where(group.NameEQ(name), group.DeletedAtIsNil()).Exist(seedCtx)
	placeholder
		require.True(t, exists, "expected group %s to exist", name)
placeholder

	assertGroupExists(service.PlatformAnthropic + "-default")
	assertGroupExists(service.PlatformOpenAI + "-default")
	assertGroupExists(service.PlatformGemini + "-default")
	assertGroupExists(service.PlatformAntigravity + "-default-1")
	assertGroupExists(service.PlatformAntigravity + "-default-2")
placeholder

func TestEnsureSimpleModeDefaultGroups_IgnoresSoftDeletedGroups(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	client := tx.Client()

	seedCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Create and then soft-delete an anthropic default group.
	g, err := client.Group.Create().
		SetName(service.PlatformAnthropic + "-default").
		SetPlatform(service.PlatformAnthropic).
		SetStatus(service.StatusActive).
		SetSubscriptionType(service.SubscriptionTypeStandard).
		SetRateMultiplier(1.0).
		SetIsExclusive(false).
		Save(seedCtx)
placeholder

	_, err = client.Group.Delete().Where(group.IDEQ(g.ID)).Exec(seedCtx)
placeholder

	require.NoError(t, ensureSimpleModeDefaultGroups(seedCtx, client))

	// New active one should exist.
	count, err := client.Group.Query().Where(group.NameEQ(service.PlatformAnthropic+"-default"), group.DeletedAtIsNil()).Count(seedCtx)
placeholder
	require.Equal(t, 1, count)
placeholder

func TestEnsureSimpleModeDefaultGroups_AntigravityNeedsTwoGroupsOnlyByCount(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	client := tx.Client()

	seedCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	mustCreateGroup(t, client, &service.Group{Name: "ag-custom-1-" + time.Now().Format(time.RFC3339Nano), Platform: service.PlatformAntigravityplaceholder)
	mustCreateGroup(t, client, &service.Group{Name: "ag-custom-2-" + time.Now().Format(time.RFC3339Nano), Platform: service.PlatformAntigravityplaceholder)

	require.NoError(t, ensureSimpleModeDefaultGroups(seedCtx, client))

	count, err := client.Group.Query().Where(group.PlatformEQ(service.PlatformAntigravity), group.DeletedAtIsNil()).Count(seedCtx)
placeholder
	require.GreaterOrEqual(t, count, 2)
placeholder
