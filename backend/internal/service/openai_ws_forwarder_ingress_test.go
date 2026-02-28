package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	coderws "github.com/coder/websocket"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestIsOpenAIWSClientDisconnectError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
placeholder{
		{name: "nil", err: nil, want: falseplaceholder,
		{name: "io_eof", err: io.EOF, want: trueplaceholder,
		{name: "net_closed", err: net.ErrClosed, want: trueplaceholder,
		{name: "context_canceled", err: context.Canceled, want: trueplaceholder,
		{name: "ws_normal_closure", err: coderws.CloseError{Code: coderws.StatusNormalClosureplaceholder, want: trueplaceholder,
		{name: "ws_going_away", err: coderws.CloseError{Code: coderws.StatusGoingAwayplaceholder, want: trueplaceholder,
		{name: "ws_no_status", err: coderws.CloseError{Code: coderws.StatusNoStatusRcvdplaceholder, want: trueplaceholder,
		{name: "ws_abnormal_1006", err: coderws.CloseError{Code: coderws.StatusAbnormalClosureplaceholder, want: trueplaceholder,
		{name: "ws_policy_violation", err: coderws.CloseError{Code: coderws.StatusPolicyViolationplaceholder, want: falseplaceholder,
		{name: "wrapped_eof_message", err: errors.New("failed to get reader: failed to read frame header: EOF"), want: trueplaceholder,
		{name: "connection_reset_by_peer", err: errors.New("failed to read frame header: read tcp 127.0.0.1:1234->127.0.0.1:5678: read: connection reset by peer"), want: trueplaceholder,
		{name: "broken_pipe", err: errors.New("write tcp 127.0.0.1:1234->127.0.0.1:5678: write: broken pipe"), want: trueplaceholder,
placeholder

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, isOpenAIWSClientDisconnectError(tt.err))
	placeholder)
placeholder
placeholder

func TestIsOpenAIWSIngressPreviousResponseNotFound(t *testing.T) {
	t.Parallel()

	require.False(t, isOpenAIWSIngressPreviousResponseNotFound(nil))
	require.False(t, isOpenAIWSIngressPreviousResponseNotFound(errors.New("plain error")))
	require.False(t, isOpenAIWSIngressPreviousResponseNotFound(
		wrapOpenAIWSIngressTurnError("read_upstream", errors.New("upstream read failed"), false),
	))
	require.False(t, isOpenAIWSIngressPreviousResponseNotFound(
		wrapOpenAIWSIngressTurnError(openAIWSIngressStagePreviousResponseNotFound, errors.New("previous response not found"), true),
	))
	require.True(t, isOpenAIWSIngressPreviousResponseNotFound(
		wrapOpenAIWSIngressTurnError(openAIWSIngressStagePreviousResponseNotFound, errors.New("previous response not found"), false),
	))
placeholder

func TestOpenAIWSIngressPreviousResponseRecoveryEnabled(t *testing.T) {
	t.Parallel()

	var nilService *OpenAIGatewayService
	require.True(t, nilService.openAIWSIngressPreviousResponseRecoveryEnabled(), "nil service should default to enabled")

	svcWithNilCfg := &OpenAIGatewayService{placeholder
	require.True(t, svcWithNilCfg.openAIWSIngressPreviousResponseRecoveryEnabled(), "nil config should default to enabled")

	svc := &OpenAIGatewayService{
		cfg: &config.Config{placeholder,
placeholder
	require.False(t, svc.openAIWSIngressPreviousResponseRecoveryEnabled(), "explicit config default should be false")

	svc.cfg.Gateway.OpenAIWS.IngressPreviousResponseRecoveryEnabled = true
	require.True(t, svc.openAIWSIngressPreviousResponseRecoveryEnabled())
placeholder

func TestDropPreviousResponseIDFromRawPayload(t *testing.T) {
	t.Parallel()

	t.Run("empty_payload", func(t *testing.T) {
		updated, removed, err := dropPreviousResponseIDFromRawPayload(nil)
	placeholder
		require.False(t, removed)
		require.Empty(t, updated)
placeholder)

	t.Run("payload_without_previous_response_id", func(t *testing.T) {
		payload := []byte(`{"type":"response.create","model":"gpt-5.1"placeholder`)
		updated, removed, err := dropPreviousResponseIDFromRawPayload(payload)
	placeholder
		require.False(t, removed)
		require.Equal(t, string(payload), string(updated))
placeholder)

	t.Run("normal_delete_success", func(t *testing.T) {
		payload := []byte(`{"type":"response.create","model":"gpt-5.1","previous_response_id":"resp_abc"placeholder`)
		updated, removed, err := dropPreviousResponseIDFromRawPayload(payload)
	placeholder
		require.True(t, removed)
		require.False(t, gjson.GetBytes(updated, "previous_response_id").Exists())
placeholder)

	t.Run("duplicate_keys_are_removed", func(t *testing.T) {
		payload := []byte(`{"type":"response.create","previous_response_id":"resp_a","input":[],"previous_response_id":"resp_b"placeholder`)
		updated, removed, err := dropPreviousResponseIDFromRawPayload(payload)
	placeholder
		require.True(t, removed)
		require.False(t, gjson.GetBytes(updated, "previous_response_id").Exists())
placeholder)

	t.Run("nil_delete_fn_uses_default_delete_logic", func(t *testing.T) {
		payload := []byte(`{"type":"response.create","model":"gpt-5.1","previous_response_id":"resp_abc"placeholder`)
		updated, removed, err := dropPreviousResponseIDFromRawPayloadWithDeleteFn(payload, nil)
	placeholder
		require.True(t, removed)
		require.False(t, gjson.GetBytes(updated, "previous_response_id").Exists())
placeholder)

	t.Run("delete_error", func(t *testing.T) {
		payload := []byte(`{"type":"response.create","model":"gpt-5.1","previous_response_id":"resp_abc"placeholder`)
		updated, removed, err := dropPreviousResponseIDFromRawPayloadWithDeleteFn(payload, func(_ []byte, _ string) ([]byte, error) {
			return nil, errors.New("delete failed")
	placeholder)
	placeholder
		require.False(t, removed)
		require.Equal(t, string(payload), string(updated))
placeholder)

	t.Run("malformed_json_is_still_best_effort_deleted", func(t *testing.T) {
		payload := []byte(`{"type":"response.create","previous_response_id":"resp_abc"`)
		require.True(t, gjson.GetBytes(payload, "previous_response_id").Exists())

		updated, removed, err := dropPreviousResponseIDFromRawPayload(payload)
	placeholder
		require.True(t, removed)
		require.False(t, gjson.GetBytes(updated, "previous_response_id").Exists())
placeholder)
placeholder

func TestAlignStoreDisabledPreviousResponseID(t *testing.T) {
	t.Parallel()

	t.Run("empty_payload", func(t *testing.T) {
		updated, changed, err := alignStoreDisabledPreviousResponseID(nil, "resp_target")
	placeholder
		require.False(t, changed)
		require.Empty(t, updated)
placeholder)

	t.Run("empty_expected", func(t *testing.T) {
		payload := []byte(`{"type":"response.create","previous_response_id":"resp_old"placeholder`)
		updated, changed, err := alignStoreDisabledPreviousResponseID(payload, "")
	placeholder
		require.False(t, changed)
		require.Equal(t, string(payload), string(updated))
placeholder)

	t.Run("missing_previous_response_id", func(t *testing.T) {
		payload := []byte(`{"type":"response.create","model":"gpt-5.1"placeholder`)
		updated, changed, err := alignStoreDisabledPreviousResponseID(payload, "resp_target")
	placeholder
		require.False(t, changed)
		require.Equal(t, string(payload), string(updated))
placeholder)

	t.Run("already_aligned", func(t *testing.T) {
		payload := []byte(`{"type":"response.create","previous_response_id":"resp_target"placeholder`)
		updated, changed, err := alignStoreDisabledPreviousResponseID(payload, "resp_target")
	placeholder
		require.False(t, changed)
		require.Equal(t, "resp_target", gjson.GetBytes(updated, "previous_response_id").String())
placeholder)

	t.Run("mismatch_rewrites_to_expected", func(t *testing.T) {
		payload := []byte(`{"type":"response.create","previous_response_id":"resp_old","input":[]placeholder`)
		updated, changed, err := alignStoreDisabledPreviousResponseID(payload, "resp_target")
	placeholder
		require.True(t, changed)
		require.Equal(t, "resp_target", gjson.GetBytes(updated, "previous_response_id").String())
placeholder)

	t.Run("duplicate_keys_rewrites_to_single_expected", func(t *testing.T) {
		payload := []byte(`{"type":"response.create","previous_response_id":"resp_old_1","input":[],"previous_response_id":"resp_old_2"placeholder`)
		updated, changed, err := alignStoreDisabledPreviousResponseID(payload, "resp_target")
	placeholder
		require.True(t, changed)
		require.Equal(t, "resp_target", gjson.GetBytes(updated, "previous_response_id").String())
placeholder)
placeholder

func TestSetPreviousResponseIDToRawPayload(t *testing.T) {
	t.Parallel()

	t.Run("empty_payload", func(t *testing.T) {
		updated, err := setPreviousResponseIDToRawPayload(nil, "resp_target")
	placeholder
		require.Empty(t, updated)
placeholder)

	t.Run("empty_previous_response_id", func(t *testing.T) {
		payload := []byte(`{"type":"response.create","model":"gpt-5.1"placeholder`)
		updated, err := setPreviousResponseIDToRawPayload(payload, "")
	placeholder
		require.Equal(t, string(payload), string(updated))
placeholder)

	t.Run("set_previous_response_id_when_missing", func(t *testing.T) {
		payload := []byte(`{"type":"response.create","model":"gpt-5.1"placeholder`)
		updated, err := setPreviousResponseIDToRawPayload(payload, "resp_target")
	placeholder
		require.Equal(t, "resp_target", gjson.GetBytes(updated, "previous_response_id").String())
		require.Equal(t, "gpt-5.1", gjson.GetBytes(updated, "model").String())
placeholder)

	t.Run("overwrite_existing_previous_response_id", func(t *testing.T) {
		payload := []byte(`{"type":"response.create","model":"gpt-5.1","previous_response_id":"resp_old"placeholder`)
		updated, err := setPreviousResponseIDToRawPayload(payload, "resp_new")
	placeholder
		require.Equal(t, "resp_new", gjson.GetBytes(updated, "previous_response_id").String())
placeholder)
placeholder

func TestShouldInferIngressFunctionCallOutputPreviousResponseID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                    string
		storeDisabled           bool
		turn                    int
		hasFunctionCallOutput   bool
		currentPreviousResponse string
		expectedPrevious        string
		want                    bool
placeholder{
		{
			name:                  "infer_when_all_conditions_match",
			storeDisabled:         true,
			turn:                  2,
			hasFunctionCallOutput: true,
			expectedPrevious:      "resp_1",
			want:                  true,
	placeholder,
		{
			name:                  "skip_when_store_enabled",
			storeDisabled:         false,
			turn:                  2,
			hasFunctionCallOutput: true,
			expectedPrevious:      "resp_1",
			want:                  false,
	placeholder,
		{
			name:                  "skip_on_first_turn",
			storeDisabled:         true,
			turn:                  1,
			hasFunctionCallOutput: true,
			expectedPrevious:      "resp_1",
			want:                  false,
	placeholder,
		{
			name:                  "skip_without_function_call_output",
			storeDisabled:         true,
			turn:                  2,
			hasFunctionCallOutput: false,
			expectedPrevious:      "resp_1",
			want:                  false,
	placeholder,
		{
			name:                    "skip_when_request_already_has_previous_response_id",
			storeDisabled:           true,
			turn:                    2,
			hasFunctionCallOutput:   true,
			currentPreviousResponse: "resp_client",
			expectedPrevious:        "resp_1",
			want:                    false,
	placeholder,
		{
			name:                  "skip_when_last_turn_response_id_missing",
			storeDisabled:         true,
			turn:                  2,
			hasFunctionCallOutput: true,
			expectedPrevious:      "",
			want:                  false,
	placeholder,
		{
			name:                  "trim_whitespace_before_judgement",
			storeDisabled:         true,
			turn:                  2,
			hasFunctionCallOutput: true,
			expectedPrevious:      "   resp_2   ",
			want:                  true,
	placeholder,
placeholder

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := shouldInferIngressFunctionCallOutputPreviousResponseID(
				tt.storeDisabled,
				tt.turn,
				tt.hasFunctionCallOutput,
				tt.currentPreviousResponse,
				tt.expectedPrevious,
			)
			require.Equal(t, tt.want, got)
	placeholder)
placeholder
placeholder

func TestOpenAIWSInputIsPrefixExtended(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		previous  []byte
		current   []byte
		want      bool
		expectErr bool
placeholder{
		{
			name:     "both_missing_input",
			previous: []byte(`{"type":"response.create","model":"gpt-5.1"placeholder`),
			current:  []byte(`{"type":"response.create","model":"gpt-5.1","previous_response_id":"resp_1"placeholder`),
			want:     true,
	placeholder,
		{
			name:     "previous_missing_current_empty_array",
			previous: []byte(`{"type":"response.create","model":"gpt-5.1"placeholder`),
			current:  []byte(`{"type":"response.create","model":"gpt-5.1","input":[]placeholder`),
			want:     true,
	placeholder,
		{
			name:     "previous_missing_current_non_empty_array",
			previous: []byte(`{"type":"response.create","model":"gpt-5.1"placeholder`),
			current:  []byte(`{"type":"response.create","model":"gpt-5.1","input":[{"type":"input_text","text":"hello"placeholder]placeholder`),
			want:     false,
	placeholder,
		{
			name:     "array_prefix_match",
			previous: []byte(`{"input":[{"type":"input_text","text":"hello"placeholder]placeholder`),
			current:  []byte(`{"input":[{"text":"hello","type":"input_text"placeholder,{"type":"input_text","text":"world"placeholder]placeholder`),
			want:     true,
	placeholder,
		{
			name:     "array_prefix_mismatch",
			previous: []byte(`{"input":[{"type":"input_text","text":"hello"placeholder]placeholder`),
			current:  []byte(`{"input":[{"type":"input_text","text":"different"placeholder]placeholder`),
			want:     false,
	placeholder,
		{
			name:     "current_shorter_than_previous",
			previous: []byte(`{"input":[{"type":"input_text","text":"a"placeholder,{"type":"input_text","text":"b"placeholder]placeholder`),
			current:  []byte(`{"input":[{"type":"input_text","text":"a"placeholder]placeholder`),
			want:     false,
	placeholder,
		{
			name:     "previous_has_input_current_missing",
			previous: []byte(`{"input":[{"type":"input_text","text":"a"placeholder]placeholder`),
			current:  []byte(`{"model":"gpt-5.1"placeholder`),
			want:     false,
	placeholder,
		{
			name:     "input_string_treated_as_single_item",
			previous: []byte(`{"input":"hello"placeholder`),
			current:  []byte(`{"input":"hello"placeholder`),
			want:     true,
	placeholder,
		{
			name:      "current_invalid_input_json",
			previous:  []byte(`{"input":[]placeholder`),
			current:   []byte(`{"input":[placeholder`),
			expectErr: true,
	placeholder,
		{
			name:      "invalid_input_json",
			previous:  []byte(`{"input":[placeholder`),
			current:   []byte(`{"input":[]placeholder`),
			expectErr: true,
	placeholder,
placeholder

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := openAIWSInputIsPrefixExtended(tt.previous, tt.current)
			if tt.expectErr {
			placeholder
				return
		placeholder
		placeholder
			require.Equal(t, tt.want, got)
	placeholder)
placeholder
placeholder

func TestNormalizeOpenAIWSJSONForCompare(t *testing.T) {
	t.Parallel()

	normalized, err := normalizeOpenAIWSJSONForCompare([]byte(`{"b":2,"a":1placeholder`))
placeholder
	require.Equal(t, `{"a":1,"b":2placeholder`, string(normalized))

	_, err = normalizeOpenAIWSJSONForCompare([]byte("   "))
placeholder

	_, err = normalizeOpenAIWSJSONForCompare([]byte(`{"a":`))
placeholder
placeholder

func TestNormalizeOpenAIWSJSONForCompareOrRaw(t *testing.T) {
	t.Parallel()

	require.Equal(t, `{"a":1,"b":2placeholder`, string(normalizeOpenAIWSJSONForCompareOrRaw([]byte(`{"b":2,"a":1placeholder`))))
	require.Equal(t, `{"a":`, string(normalizeOpenAIWSJSONForCompareOrRaw([]byte(`{"a":`))))
placeholder

func TestNormalizeOpenAIWSPayloadWithoutInputAndPreviousResponseID(t *testing.T) {
	t.Parallel()

	normalized, err := normalizeOpenAIWSPayloadWithoutInputAndPreviousResponseID(
		[]byte(`{"model":"gpt-5.1","input":[1],"previous_response_id":"resp_x","metadata":{"b":2,"a":1placeholderplaceholder`),
	)
placeholder
	require.False(t, gjson.GetBytes(normalized, "input").Exists())
	require.False(t, gjson.GetBytes(normalized, "previous_response_id").Exists())
	require.Equal(t, float64(1), gjson.GetBytes(normalized, "metadata.a").Float())

	_, err = normalizeOpenAIWSPayloadWithoutInputAndPreviousResponseID(nil)
placeholder

	_, err = normalizeOpenAIWSPayloadWithoutInputAndPreviousResponseID([]byte(`[]`))
placeholder
placeholder

func TestOpenAIWSExtractNormalizedInputSequence(t *testing.T) {
	t.Parallel()

	t.Run("empty_payload", func(t *testing.T) {
		items, exists, err := openAIWSExtractNormalizedInputSequence(nil)
	placeholder
		require.False(t, exists)
		require.Nil(t, items)
placeholder)

	t.Run("input_missing", func(t *testing.T) {
		items, exists, err := openAIWSExtractNormalizedInputSequence([]byte(`{"type":"response.create"placeholder`))
	placeholder
		require.False(t, exists)
		require.Nil(t, items)
placeholder)

	t.Run("input_array", func(t *testing.T) {
		items, exists, err := openAIWSExtractNormalizedInputSequence([]byte(`{"input":[{"type":"input_text","text":"hello"placeholder]placeholder`))
	placeholder
		require.True(t, exists)
		require.Len(t, items, 1)
placeholder)

	t.Run("input_object", func(t *testing.T) {
		items, exists, err := openAIWSExtractNormalizedInputSequence([]byte(`{"input":{"type":"input_text","text":"hello"placeholderplaceholder`))
	placeholder
		require.True(t, exists)
		require.Len(t, items, 1)
placeholder)

	t.Run("input_string", func(t *testing.T) {
		items, exists, err := openAIWSExtractNormalizedInputSequence([]byte(`{"input":"hello"placeholder`))
	placeholder
		require.True(t, exists)
		require.Len(t, items, 1)
		require.Equal(t, `"hello"`, string(items[0]))
placeholder)

	t.Run("input_number", func(t *testing.T) {
		items, exists, err := openAIWSExtractNormalizedInputSequence([]byte(`{"input":42placeholder`))
	placeholder
		require.True(t, exists)
		require.Len(t, items, 1)
		require.Equal(t, "42", string(items[0]))
placeholder)

	t.Run("input_bool", func(t *testing.T) {
		items, exists, err := openAIWSExtractNormalizedInputSequence([]byte(`{"input":trueplaceholder`))
	placeholder
		require.True(t, exists)
		require.Len(t, items, 1)
		require.Equal(t, "true", string(items[0]))
placeholder)

	t.Run("input_null", func(t *testing.T) {
		items, exists, err := openAIWSExtractNormalizedInputSequence([]byte(`{"input":nullplaceholder`))
	placeholder
		require.True(t, exists)
		require.Len(t, items, 1)
		require.Equal(t, "null", string(items[0]))
placeholder)

	t.Run("input_invalid_array_json", func(t *testing.T) {
		items, exists, err := openAIWSExtractNormalizedInputSequence([]byte(`{"input":[placeholder`))
	placeholder
		require.True(t, exists)
		require.Nil(t, items)
placeholder)
placeholder

func TestShouldKeepIngressPreviousResponseID(t *testing.T) {
	t.Parallel()

	previousPayload := []byte(`{
		"type":"response.create",
		"model":"gpt-5.1",
		"store":false,
		"tools":[{"type":"function","name":"tool_a"placeholder],
		"input":[{"type":"input_text","text":"hello"placeholder]
placeholder`)
	currentStrictPayload := []byte(`{
		"type":"response.create",
		"model":"gpt-5.1",
		"store":false,
		"tools":[{"name":"tool_a","type":"function"placeholder],
		"previous_response_id":"resp_turn_1",
		"input":[{"text":"hello","type":"input_text"placeholder,{"type":"input_text","text":"world"placeholder]
placeholder`)

	t.Run("strict_incremental_keep", func(t *testing.T) {
		keep, reason, err := shouldKeepIngressPreviousResponseID(previousPayload, currentStrictPayload, "resp_turn_1", false)
	placeholder
		require.True(t, keep)
		require.Equal(t, "strict_incremental_ok", reason)
placeholder)

	t.Run("missing_previous_response_id", func(t *testing.T) {
		payload := []byte(`{"type":"response.create","model":"gpt-5.1","input":[]placeholder`)
		keep, reason, err := shouldKeepIngressPreviousResponseID(previousPayload, payload, "resp_turn_1", false)
	placeholder
		require.False(t, keep)
		require.Equal(t, "missing_previous_response_id", reason)
placeholder)

	t.Run("missing_last_turn_response_id", func(t *testing.T) {
		keep, reason, err := shouldKeepIngressPreviousResponseID(previousPayload, currentStrictPayload, "", false)
	placeholder
		require.False(t, keep)
		require.Equal(t, "missing_last_turn_response_id", reason)
placeholder)

	t.Run("previous_response_id_mismatch", func(t *testing.T) {
		keep, reason, err := shouldKeepIngressPreviousResponseID(previousPayload, currentStrictPayload, "resp_turn_other", false)
	placeholder
		require.False(t, keep)
		require.Equal(t, "previous_response_id_mismatch", reason)
placeholder)

	t.Run("missing_previous_turn_payload", func(t *testing.T) {
		keep, reason, err := shouldKeepIngressPreviousResponseID(nil, currentStrictPayload, "resp_turn_1", false)
	placeholder
		require.False(t, keep)
		require.Equal(t, "missing_previous_turn_payload", reason)
placeholder)

	t.Run("non_input_changed", func(t *testing.T) {
		payload := []byte(`{
			"type":"response.create",
			"model":"gpt-5.1-mini",
			"store":false,
			"tools":[{"type":"function","name":"tool_a"placeholder],
			"previous_response_id":"resp_turn_1",
			"input":[{"type":"input_text","text":"hello"placeholder,{"type":"input_text","text":"world"placeholder]
	placeholder`)
		keep, reason, err := shouldKeepIngressPreviousResponseID(previousPayload, payload, "resp_turn_1", false)
	placeholder
		require.False(t, keep)
		require.Equal(t, "non_input_changed", reason)
placeholder)

	t.Run("delta_input_keeps_previous_response_id", func(t *testing.T) {
		payload := []byte(`{
			"type":"response.create",
			"model":"gpt-5.1",
			"store":false,
			"tools":[{"type":"function","name":"tool_a"placeholder],
			"previous_response_id":"resp_turn_1",
			"input":[{"type":"input_text","text":"different"placeholder]
	placeholder`)
		keep, reason, err := shouldKeepIngressPreviousResponseID(previousPayload, payload, "resp_turn_1", false)
	placeholder
		require.True(t, keep)
		require.Equal(t, "strict_incremental_ok", reason)
placeholder)

	t.Run("function_call_output_keeps_previous_response_id", func(t *testing.T) {
		payload := []byte(`{
			"type":"response.create",
			"model":"gpt-5.1",
			"store":false,
			"previous_response_id":"resp_external",
			"input":[{"type":"function_call_output","call_id":"call_1","output":"ok"placeholder]
	placeholder`)
		keep, reason, err := shouldKeepIngressPreviousResponseID(previousPayload, payload, "resp_turn_1", true)
	placeholder
		require.True(t, keep)
		require.Equal(t, "has_function_call_output", reason)
placeholder)

	t.Run("non_input_compare_error", func(t *testing.T) {
		keep, reason, err := shouldKeepIngressPreviousResponseID([]byte(`[]`), currentStrictPayload, "resp_turn_1", false)
	placeholder
		require.False(t, keep)
		require.Equal(t, "non_input_compare_error", reason)
placeholder)

	t.Run("current_payload_compare_error", func(t *testing.T) {
		keep, reason, err := shouldKeepIngressPreviousResponseID(previousPayload, []byte(`{"previous_response_id":"resp_turn_1","input":[placeholder`), "resp_turn_1", false)
	placeholder
		require.False(t, keep)
		require.Equal(t, "non_input_compare_error", reason)
placeholder)
placeholder

func TestBuildOpenAIWSReplayInputSequence(t *testing.T) {
	t.Parallel()

	lastFull := []json.RawMessage{
		json.RawMessage(`{"type":"input_text","text":"hello"placeholder`),
placeholder

	t.Run("no_previous_response_id_use_current", func(t *testing.T) {
		items, exists, err := buildOpenAIWSReplayInputSequence(
			lastFull,
			true,
			[]byte(`{"input":[{"type":"input_text","text":"new"placeholder]placeholder`),
			false,
		)
	placeholder
		require.True(t, exists)
		require.Len(t, items, 1)
		require.Equal(t, "new", gjson.GetBytes(items[0], "text").String())
placeholder)

	t.Run("previous_response_id_delta_append", func(t *testing.T) {
		items, exists, err := buildOpenAIWSReplayInputSequence(
			lastFull,
			true,
			[]byte(`{"previous_response_id":"resp_1","input":[{"type":"input_text","text":"world"placeholder]placeholder`),
			true,
		)
	placeholder
		require.True(t, exists)
		require.Len(t, items, 2)
		require.Equal(t, "hello", gjson.GetBytes(items[0], "text").String())
		require.Equal(t, "world", gjson.GetBytes(items[1], "text").String())
placeholder)

	t.Run("previous_response_id_full_input_replace", func(t *testing.T) {
		items, exists, err := buildOpenAIWSReplayInputSequence(
			lastFull,
			true,
			[]byte(`{"previous_response_id":"resp_1","input":[{"type":"input_text","text":"hello"placeholder,{"type":"input_text","text":"world"placeholder]placeholder`),
			true,
		)
	placeholder
		require.True(t, exists)
		require.Len(t, items, 2)
		require.Equal(t, "hello", gjson.GetBytes(items[0], "text").String())
		require.Equal(t, "world", gjson.GetBytes(items[1], "text").String())
placeholder)
placeholder

func TestSetOpenAIWSPayloadInputSequence(t *testing.T) {
	t.Parallel()

	t.Run("set_items", func(t *testing.T) {
		original := []byte(`{"type":"response.create","previous_response_id":"resp_1"placeholder`)
		items := []json.RawMessage{
			json.RawMessage(`{"type":"input_text","text":"hello"placeholder`),
			json.RawMessage(`{"type":"input_text","text":"world"placeholder`),
	placeholder
		updated, err := setOpenAIWSPayloadInputSequence(original, items, true)
	placeholder
		require.Equal(t, "hello", gjson.GetBytes(updated, "input.0.text").String())
		require.Equal(t, "world", gjson.GetBytes(updated, "input.1.text").String())
placeholder)

	t.Run("preserve_empty_array_not_null", func(t *testing.T) {
		original := []byte(`{"type":"response.create","previous_response_id":"resp_1"placeholder`)
		updated, err := setOpenAIWSPayloadInputSequence(original, nil, true)
	placeholder
		require.True(t, gjson.GetBytes(updated, "input").IsArray())
		require.Len(t, gjson.GetBytes(updated, "input").Array(), 0)
		require.False(t, gjson.GetBytes(updated, "input").Type == gjson.Null)
placeholder)
placeholder

func TestCloneOpenAIWSRawMessages(t *testing.T) {
	t.Parallel()

	t.Run("nil_slice", func(t *testing.T) {
		cloned := cloneOpenAIWSRawMessages(nil)
		require.Nil(t, cloned)
placeholder)

	t.Run("empty_slice", func(t *testing.T) {
		items := make([]json.RawMessage, 0)
		cloned := cloneOpenAIWSRawMessages(items)
		require.NotNil(t, cloned)
		require.Len(t, cloned, 0)
placeholder)
placeholder
