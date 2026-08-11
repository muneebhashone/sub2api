//go:build unit

package redissession

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestStoreRoundTripAndSingleUse(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()placeholder)
	t.Cleanup(func() { _ = rdb.Close() placeholder)
	store := New(rdb, "oauth:test", time.Minute)
	ctx := context.Background()

	require.NoError(t, store.Set(ctx, "sid", map[string]string{"state": "state"placeholder))
	var got map[string]string
	ok, err := store.Get(ctx, "sid", &got)
placeholder
	require.True(t, ok)
	require.Equal(t, "state", got["state"])

	ok, err = store.TryConsume(ctx, "sid")
placeholder
	require.True(t, ok)
	ok, err = store.TryConsume(ctx, "sid")
placeholder
	require.False(t, ok)

	require.NoError(t, store.Delete(ctx, "sid"))
	ok, err = store.Get(ctx, "sid", &got)
placeholder
	require.False(t, ok)
placeholder
