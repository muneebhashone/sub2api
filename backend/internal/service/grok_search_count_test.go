package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCountGrokNativeSearchCallsFromJSONBytes(t *testing.T) {
	t.Parallel()
	require.Equal(t, 0, countGrokNativeSearchCallsFromJSONBytes(nil))
	require.Equal(t, 0, countGrokNativeSearchCallsFromJSONBytes([]byte(`{"output":[]placeholder`)))
	body := []byte(`{"output":[
		{"type":"web_search_call","id":"ws1","status":"completed"placeholder,
		{"type":"x_search_call","id":"xs1"placeholder,
		{"type":"function_call","name":"tool_search","call_id":"ts1"placeholder,
		{"type":"function_call","name":"lookup","call_id":"other"placeholder
	]placeholder`)
	require.Equal(t, 3, countGrokNativeSearchCallsFromJSONBytes(body))
placeholder

func TestCountGrokNativeSearchCallsFromSSEBodyDedups(t *testing.T) {
	t.Parallel()
	sse := stringsJoin(
		`data: {"type":"response.output_item.done","item":{"type":"web_search_call","id":"ws1","call_id":"c1"placeholderplaceholder`,
		`data: {"type":"response.output_item.done","item":{"type":"web_search_call","id":"ws1","call_id":"c1"placeholderplaceholder`,
		`data: {"type":"response.completed","response":{"output":[{"type":"web_search_call","id":"ws1","call_id":"c1"placeholder,{"type":"x_search_call","id":"xs1","call_id":"c2"placeholder]placeholderplaceholder`,
	)
	require.Equal(t, 2, countGrokNativeSearchCallsFromSSEBody(sse))
placeholder

func TestCountGrokNativeSearchCallsInSSEDataDedup_LiveStreamPath(t *testing.T) {
	t.Parallel()
	// Mirrors the live streaming accumulator: item.done then response.completed
	// for the same call_id must bill once (regression for ~2× surcharge).
	seen := make(map[string]struct{placeholder)
	done := []byte(`{"type":"response.output_item.done","item":{"type":"web_search_call","id":"ws1","call_id":"c1"placeholderplaceholder`)
	completed := []byte(`{"type":"response.completed","response":{"output":[{"type":"web_search_call","id":"ws1","call_id":"c1"placeholder,{"type":"x_search_call","id":"xs1","call_id":"c2"placeholder]placeholderplaceholder`)
	require.Equal(t, 1, countGrokNativeSearchCallsInSSEDataDedup(done, seen))
	require.Equal(t, 1, countGrokNativeSearchCallsInSSEDataDedup(completed, seen))
	// Raw (no-dedup) path still double-counts the same envelope pair.
	require.Equal(t, 1, countGrokNativeSearchCallsInSSEData(done))
	require.Equal(t, 2, countGrokNativeSearchCallsInSSEData(completed))
placeholder

func TestCountGrokNativeSearchCallsInSSEDataDedup_NoIDStillDedups(t *testing.T) {
	t.Parallel()
	// Upstream sometimes omits call_id/id; synthetic keys must still prevent 2×.
	seen := make(map[string]struct{placeholder)
	done := []byte(`{"type":"response.output_item.done","item":{"type":"web_search_call"placeholderplaceholder`)
	completed := []byte(`{"type":"response.completed","response":{"output":[{"type":"web_search_call"placeholder]placeholderplaceholder`)
	require.Equal(t, 1, countGrokNativeSearchCallsInSSEDataDedup(done, seen))
	require.Equal(t, 0, countGrokNativeSearchCallsInSSEDataDedup(completed, seen))
placeholder

func stringsJoin(lines ...string) string {
	out := ""
	for _, l := range lines {
		out += l + "\n\n"
placeholder
	return out
placeholder
