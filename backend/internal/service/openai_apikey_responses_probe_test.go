package service

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/stretchr/testify/require"
)

func TestDecideResponsesProbeSupport(t *testing.T) {
	fnCall := []byte(`{"output":[{"type":"reasoning"placeholder,{"type":"function_call","name":"probe_ping"placeholder]placeholder`)
	reasoningOnly := []byte(`{"output":[{"type":"reasoning"placeholder]placeholder`)

	cases := []struct {
		name   string
		status int
		body   []byte
		want   bool
placeholder{
		// Endpoint clearly absent on third-party OpenAI-compatible upstreams.
		{"404 endpoint absent", 404, fnCall, falseplaceholder,
		{"405 method not allowed", 405, fnCall, falseplaceholder,
		// 2xx: tool capability is judged by presence of a function_call output item.
		{"200 with function_call", 200, fnCall, trueplaceholder,
		// Volcengine Ark coding/v3 × kimi-k2.6: reasoning only, no function_call.
		{"200 reasoning only", 200, reasoningOnly, falseplaceholder,
		{"200 invalid json", 200, []byte("not-json"), falseplaceholder,
		{"200 no output field", 200, []byte(`{"status":"completed"placeholder`), falseplaceholder,
		// Non-2xx (other than 404/405): endpoint exists, capability undecidable -> conservative true.
		{"400 conservative true", 400, reasoningOnly, trueplaceholder,
		{"401 conservative true", 401, nil, trueplaceholder,
		{"500 conservative true", 500, nil, trueplaceholder,
placeholder
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, decideResponsesProbeSupport(tc.status, tc.body))
	placeholder)
placeholder
placeholder

func TestResponsesProbeBodyHasFunctionCall(t *testing.T) {
	require.True(t, responsesProbeBodyHasFunctionCall([]byte(`{"output":[{"type":"function_call"placeholder]placeholder`)))
	require.True(t, responsesProbeBodyHasFunctionCall([]byte(`{"output":[{"type":"reasoning"placeholder,{"type":"function_call"placeholder]placeholder`)))
	require.False(t, responsesProbeBodyHasFunctionCall([]byte(`{"output":[{"type":"reasoning"placeholder]placeholder`)))
	require.False(t, responsesProbeBodyHasFunctionCall([]byte(`{"output":[]placeholder`)))
	require.False(t, responsesProbeBodyHasFunctionCall([]byte(`{placeholder`)))
	require.False(t, responsesProbeBodyHasFunctionCall([]byte(`garbage`)))
placeholder

func TestSelectResponsesProbeModel(t *testing.T) {
	// No model_mapping -> fall back to DefaultTestModel (OpenAI official APIKey).
	require.Equal(t, openai.DefaultTestModel, selectResponsesProbeModel(&Account{placeholder))

	// model_mapping values are upstream models; pick first by sort for reproducibility.
	acct := &Account{Credentials: map[string]any{
		"model_mapping": map[string]any{
			"client-b": "zeta-model",
			"client-a": "alpha-model",
	placeholder,
placeholderplaceholder
	require.Equal(t, "alpha-model", selectResponsesProbeModel(acct))

	// Wildcard / blank upstream values are skipped.
	acctWild := &Account{Credentials: map[string]any{
		"model_mapping": map[string]any{
			"a": "*",
			"b": "  ",
			"c": "real-model",
	placeholder,
placeholderplaceholder
	require.Equal(t, "real-model", selectResponsesProbeModel(acctWild))

	// Only wildcard mappings -> DefaultTestModel.
	acctAllWild := &Account{Credentials: map[string]any{
		"model_mapping": map[string]any{"a": "gpt-*"placeholder,
placeholderplaceholder
	require.Equal(t, openai.DefaultTestModel, selectResponsesProbeModel(acctAllWild))
placeholder
