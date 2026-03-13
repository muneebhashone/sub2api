//go:build unit

package admin

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSnapshotCache_SetAndGet(t *testing.T) {
	c := newSnapshotCache(5 * time.Second)

	entry := c.Set("key1", map[string]string{"hello": "world"placeholder)
	require.NotEmpty(t, entry.ETag)
	require.NotNil(t, entry.Payload)

	got, ok := c.Get("key1")
	require.True(t, ok)
	require.Equal(t, entry.ETag, got.ETag)
placeholder

func TestSnapshotCache_Expiration(t *testing.T) {
	c := newSnapshotCache(1 * time.Millisecond)

	c.Set("key1", "value")
	time.Sleep(5 * time.Millisecond)

	_, ok := c.Get("key1")
	require.False(t, ok, "expired entry should not be returned")
placeholder

func TestSnapshotCache_GetEmptyKey(t *testing.T) {
	c := newSnapshotCache(5 * time.Second)
	_, ok := c.Get("")
	require.False(t, ok)
placeholder

func TestSnapshotCache_GetMiss(t *testing.T) {
	c := newSnapshotCache(5 * time.Second)
	_, ok := c.Get("nonexistent")
	require.False(t, ok)
placeholder

func TestSnapshotCache_NilReceiver(t *testing.T) {
	var c *snapshotCache
	_, ok := c.Get("key")
	require.False(t, ok)

	entry := c.Set("key", "value")
	require.Empty(t, entry.ETag)
placeholder

func TestSnapshotCache_SetEmptyKey(t *testing.T) {
	c := newSnapshotCache(5 * time.Second)

	// Set with empty key should return entry but not store it
	entry := c.Set("", "value")
	require.NotEmpty(t, entry.ETag)

	_, ok := c.Get("")
	require.False(t, ok)
placeholder

func TestSnapshotCache_DefaultTTL(t *testing.T) {
	c := newSnapshotCache(0)
	require.Equal(t, 30*time.Second, c.ttl)

	c2 := newSnapshotCache(-1 * time.Second)
	require.Equal(t, 30*time.Second, c2.ttl)
placeholder

func TestSnapshotCache_ETagDeterministic(t *testing.T) {
	c := newSnapshotCache(5 * time.Second)
	payload := map[string]int{"a": 1, "b": 2placeholder

	entry1 := c.Set("k1", payload)
	entry2 := c.Set("k2", payload)
	require.Equal(t, entry1.ETag, entry2.ETag, "same payload should produce same ETag")
placeholder

func TestSnapshotCache_ETagFormat(t *testing.T) {
	c := newSnapshotCache(5 * time.Second)
	entry := c.Set("k", "test")
	// ETag should be quoted hex string: "abcdef..."
	require.True(t, len(entry.ETag) > 2)
	require.Equal(t, byte('"'), entry.ETag[0])
	require.Equal(t, byte('"'), entry.ETag[len(entry.ETag)-1])
placeholder

func TestBuildETagFromAny_UnmarshalablePayload(t *testing.T) {
	// channels are not JSON-serializable
	etag := buildETagFromAny(make(chan int))
	require.Empty(t, etag)
placeholder

func TestSnapshotCache_GetOrLoad_MissThenHit(t *testing.T) {
	c := newSnapshotCache(5 * time.Second)
	var loads atomic.Int32

	entry, hit, err := c.GetOrLoad("key1", func() (any, error) {
		loads.Add(1)
		return map[string]string{"hello": "world"placeholder, nil
placeholder)
placeholder
	require.False(t, hit)
	require.NotEmpty(t, entry.ETag)
	require.Equal(t, int32(1), loads.Load())

	entry2, hit, err := c.GetOrLoad("key1", func() (any, error) {
		loads.Add(1)
		return map[string]string{"unexpected": "value"placeholder, nil
placeholder)
placeholder
	require.True(t, hit)
	require.Equal(t, entry.ETag, entry2.ETag)
	require.Equal(t, int32(1), loads.Load())
placeholder

func TestSnapshotCache_GetOrLoad_ConcurrentSingleflight(t *testing.T) {
	c := newSnapshotCache(5 * time.Second)
	var loads atomic.Int32
	start := make(chan struct{placeholder)
	const callers = 8
	errCh := make(chan error, callers)

	var wg sync.WaitGroup
	wg.Add(callers)
	for range callers {
		go func() {
			defer wg.Done()
			<-start
			_, _, err := c.GetOrLoad("shared", func() (any, error) {
				loads.Add(1)
				time.Sleep(20 * time.Millisecond)
				return "value", nil
		placeholder)
			errCh <- err
	placeholder()
placeholder
	close(start)
	wg.Wait()
	close(errCh)

	for err := range errCh {
	placeholder
placeholder

	require.Equal(t, int32(1), loads.Load())
placeholder

func TestParseBoolQueryWithDefault(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		def  bool
		want bool
placeholder{
		{"empty returns default true", "", true, trueplaceholder,
		{"empty returns default false", "", false, falseplaceholder,
		{"1", "1", false, trueplaceholder,
		{"true", "true", false, trueplaceholder,
		{"TRUE", "TRUE", false, trueplaceholder,
		{"yes", "yes", false, trueplaceholder,
		{"on", "on", false, trueplaceholder,
		{"0", "0", true, falseplaceholder,
		{"false", "false", true, falseplaceholder,
		{"FALSE", "FALSE", true, falseplaceholder,
		{"no", "no", true, falseplaceholder,
		{"off", "off", true, falseplaceholder,
		{"whitespace trimmed", "  true  ", false, trueplaceholder,
		{"unknown returns default true", "maybe", true, trueplaceholder,
		{"unknown returns default false", "maybe", false, falseplaceholder,
placeholder

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := parseBoolQueryWithDefault(tc.raw, tc.def)
			require.Equal(t, tc.want, got)
	placeholder)
placeholder
placeholder
