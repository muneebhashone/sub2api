package service

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestOpenAIResponsesInputItemIDPrefixUsesObservedOutputContracts(t *testing.T) {
	tests := []struct {
		itemType string
		id       string
		strip    bool
placeholder{
		{itemType: "message", id: "msg_123", strip: falseplaceholder,
		{itemType: "message", id: "item_123", strip: trueplaceholder,
		{itemType: "reasoning", id: "rs_123", strip: falseplaceholder,
		{itemType: "reasoning", id: "item_123", strip: trueplaceholder,
		{itemType: "function_call", id: "fc_123", strip: falseplaceholder,
		{itemType: "function_call", id: "call_123", strip: trueplaceholder,
		{itemType: "tool_call", id: "fc_123", strip: falseplaceholder,
		{itemType: "local_shell_call", id: "fc_123", strip: falseplaceholder,
		{itemType: "mcp_tool_call", id: "fc_123", strip: falseplaceholder,
		{itemType: "custom_tool_call", id: "ctc_123", strip: falseplaceholder,
		{itemType: "custom_tool_call", id: "fc_123", strip: trueplaceholder,
		{itemType: "tool_search_call", id: "tsc_123", strip: falseplaceholder,
		{itemType: "tool_search_call", id: "fc_123", strip: trueplaceholder,
		{itemType: "web_search_call", id: "ws_123", strip: falseplaceholder,
		{itemType: "web_search_call", id: "item_123", strip: trueplaceholder,
		{itemType: "custom_tool_call_output", id: "fc_123", strip: falseplaceholder,
		{itemType: "custom_tool_call_output", id: "ctco_123", strip: trueplaceholder,
		// Do not impose an inferred contract on output types for which there is
		// no observed upstream prefix rejection.
		{itemType: "function_call_output", id: "fco_123", strip: falseplaceholder,
		{itemType: "tool_search_output", id: "tso_123", strip: falseplaceholder,
		{itemType: "mcp_tool_call_output", id: "mcpo_123", strip: falseplaceholder,
		{itemType: "future_item", id: "item_123", strip: falseplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.itemType+"/"+tt.id, func(t *testing.T) {
			require.Equal(t, tt.strip, shouldStripOpenAIResponsesInputItemID(tt.itemType, tt.id))
	placeholder)
placeholder
placeholder

func TestSanitizeOpenAIResponsesInputItemIDsDoesNotCascadeAcrossIDNamespaces(t *testing.T) {
	body := []byte(`{"input":[
		{"type":"function_call","id":"item_bad_call","call_id":"call_valid","name":"lookup","arguments":"{placeholder"placeholder,
		{"type":"function_call_output","call_id":"call_valid","output":"preserve paired output"placeholder,
		{"type":"function_call_output","call_id":"item_bad_call","output":"preserve opaque output"placeholder,
		{"type":"item_reference","id":"item_bad_call"placeholder,
		{"type":"item_reference","id":"remote_valid"placeholder,
		{"type":"custom_tool_call","id":"ctc_valid","call_id":"ctco_bad_output","name":"apply_patch","input":"patch"placeholder,
		{"type":"custom_tool_call_output","id":"ctco_bad_output","call_id":"ctco_bad_output","output":"preserve by call_id"placeholder
	]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesInputItemIDs(body)

placeholder
	require.True(t, changed)
	items := gjson.GetBytes(sanitized, "input").Array()
	require.Len(t, items, 7)
	require.False(t, items[0].Get("id").Exists())
	require.Equal(t, "call_valid", items[0].Get("call_id").String())
	require.Equal(t, "preserve paired output", items[1].Get("output").String())
	require.Equal(t, "preserve opaque output", items[2].Get("output").String())
	require.Equal(t, "item_bad_call", items[3].Get("id").String())
	require.Equal(t, "remote_valid", items[4].Get("id").String())
	require.Equal(t, "ctc_valid", items[5].Get("id").String())
	require.False(t, items[6].Get("id").Exists())
	require.Equal(t, "ctco_bad_output", items[6].Get("call_id").String())
placeholder

func TestSanitizeOpenAIResponsesInputItemIDsLeavesUnrelatedReferencesUntouched(t *testing.T) {
	body := []byte(`{"previous_response_id":"resp_1","input":[{"type":"item_reference","id":"remote_item"placeholder]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesInputItemIDs(body)

placeholder
	require.False(t, changed)
	require.Equal(t, body, sanitized)
placeholder

func TestSanitizeOpenAIResponsesInputItemIDsPreservesReferenceToDuplicateRetainedID(t *testing.T) {
	body := []byte(`{"input":[{"type":"function_call","id":"ctc_shared","call_id":"call_1"placeholder,{"type":"custom_tool_call","id":"ctc_shared","call_id":"call_2"placeholder,{"type":"item_reference","id":"ctc_shared"placeholder]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesInputItemIDs(body)

placeholder
	require.True(t, changed)
	require.False(t, gjson.GetBytes(sanitized, "input.0.id").Exists())
	require.Equal(t, "ctc_shared", gjson.GetBytes(sanitized, "input.1.id").String())
	require.Equal(t, "ctc_shared", gjson.GetBytes(sanitized, "input.2.id").String())
placeholder

func TestSanitizeOpenAIResponsesInputItemIDsPreservesOpaqueOutputsAndReferences(t *testing.T) {
	body := []byte(`{"input":[
		{"type":"function_call","id":"item_shared","call_id":"call_real"placeholder,
		{"type":"function_call_output","id":"item_shared","call_id":"item_shared","output":"dangling"placeholder,
		{"type":"item_reference","id":"item_shared"placeholder,
		{"type":"function_call_output","id":"kept_output","call_id":"call_real","output":"kept"placeholder,
		{"type":"item_reference","id":"kept_output"placeholder
	]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesInputItemIDs(body)

placeholder
	require.True(t, changed)
	require.Len(t, gjson.GetBytes(sanitized, "input").Array(), 5)
	require.False(t, gjson.GetBytes(sanitized, "input.0.id").Exists())
	require.Equal(t, "dangling", gjson.GetBytes(sanitized, "input.1.output").String())
	require.Equal(t, "item_shared", gjson.GetBytes(sanitized, "input.2.id").String())
	require.Equal(t, "kept_output", gjson.GetBytes(sanitized, "input.3.id").String())
	require.Equal(t, "kept_output", gjson.GetBytes(sanitized, "input.4.id").String())

	second, changedAgain, err := sanitizeOpenAIResponsesInputItemIDs(sanitized)
placeholder
	require.False(t, changedAgain)
	require.Equal(t, sanitized, second)
placeholder

func TestSanitizeOpenAIResponsesInputItemIDsStripsEmptyKnownIDsOnly(t *testing.T) {
	body := []byte(`{"input":[{"type":"message","id":"","content":"hello"placeholder,{"type":"future_item","id":""placeholder]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesInputItemIDs(body)

placeholder
	require.True(t, changed)
	require.False(t, gjson.GetBytes(sanitized, "input.0.id").Exists())
	require.True(t, gjson.GetBytes(sanitized, "input.1.id").Exists())
placeholder

func TestSanitizeOpenAIResponsesInputItemIDsStripsOnlyNonPairCallIDs(t *testing.T) {
	body := []byte(`{"input":[
		{"type":"message","call_id":"remove_message","content":"hi"placeholder,
		{"type":"reasoning","call_id":"remove_reasoning","id":"rs_keep","encrypted_content":"cipher","summary":[]placeholder,
		{"type":"image_generation_call","call_id":"remove_image","id":"ig_keep","status":"completed"placeholder,
		{"type":"function_call","call_id":"keep_function","name":"lookup","arguments":"{placeholder"placeholder,
		{"type":"function_call_output","call_id":"keep_function","output":"ok"placeholder,
		{"type":"custom_tool_call","call_id":"keep_custom","name":"patch","input":"x"placeholder,
		{"type":"custom_tool_call_output","call_id":"keep_custom","output":"ok"placeholder,
		{"type":"tool_search_call","call_id":"keep_search","arguments":"{placeholder"placeholder,
		{"type":"tool_search_output","call_id":"keep_search","output":"ok"placeholder,
		{"type":"local_shell_call","call_id":"keep_shell","name":"shell","arguments":"{placeholder"placeholder
	]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesInputItemIDs(body)
placeholder
	require.True(t, changed)
	for i := 0; i < 3; i++ {
		require.False(t, gjson.GetBytes(sanitized, "input."+strconv.Itoa(i)+".call_id").Exists())
placeholder
	for i := 3; i < 10; i++ {
		require.True(t, gjson.GetBytes(sanitized, "input."+strconv.Itoa(i)+".call_id").Exists())
placeholder
placeholder
