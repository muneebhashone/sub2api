package service

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestSanitizeGrokResponsesModelInputNormalizesReplayHistory(t *testing.T) {
	body := []byte(`{
		"input":[
			{"role":"system","content":"rules"placeholder,
			{"type":"message","role":"assistant","id":"msg_partial","content":[{"type":"output_text","text":"done"placeholder]placeholder,
			{"type":"reasoning","id":"rs_1","status":"completed","content":null,"summary":[{"type":"summary_text","text":"keep"placeholder]placeholder,
			{"type":"message","role":"assistant","id":"msg_complete","status":"completed","content":[{"type":"output_text","text":"kept"placeholder]placeholder
		]placeholder`)

	patched, err := sanitizeGrokResponsesModelInput(body)
placeholder
	require.Equal(t, "message", gjson.GetBytes(patched, "input.0.type").String())
	require.Equal(t, "done", gjson.GetBytes(patched, "input.1.content").String())
	require.False(t, gjson.GetBytes(patched, "input.1.id").Exists())
	require.False(t, gjson.GetBytes(patched, "input.2.status").Exists())
	require.False(t, gjson.GetBytes(patched, "input.2.content").Exists())
	require.Equal(t, "keep", gjson.GetBytes(patched, "input.2.summary.0.text").String())
	require.Equal(t, "msg_complete", gjson.GetBytes(patched, "input.3.id").String())
	require.Equal(t, "completed", gjson.GetBytes(patched, "input.3.status").String())
placeholder

func TestSanitizeGrokResponsesModelInputStripsOnlyNonPairCallIDs(t *testing.T) {
	body := []byte(`{"input":[
		{"type":"message","role":"user","call_id":"remove_message","content":"continue"placeholder,
		{"type":"reasoning","call_id":"remove_reasoning","summary":[]placeholder,
		{"type":"image_generation_call","call_id":"remove_image","id":"ig_1","status":"completed"placeholder,
		{"type":"function_call","call_id":"keep_function","name":"lookup","arguments":"{placeholder"placeholder,
		{"type":"function_call_output","call_id":"keep_function","output":"ok"placeholder
	]placeholder`)

	patched, err := sanitizeGrokResponsesModelInput(body)
placeholder
	for i := 0; i < 3; i++ {
		require.False(t, gjson.GetBytes(patched, "input."+strconv.Itoa(i)+".call_id").Exists())
placeholder
	require.Equal(t, "keep_function", gjson.GetBytes(patched, "input.3.call_id").String())
	require.Equal(t, "keep_function", gjson.GetBytes(patched, "input.4.call_id").String())
placeholder

func TestSanitizeGrokResponsesModelInputPairsCallsInTwoPasses(t *testing.T) {
	body := []byte(`{
		"input":[
			{"type":"function_call_output","output":{"before":true,"large_id":9007199254740993placeholderplaceholder,
			{"type":"function_call","function":{"name":"lookup","arguments":{"query":"x","large_id":9007199254740993placeholderplaceholderplaceholder,
			{"type":"custom_tool_call","id":"custom_1","name":"apply_patch","input":"patch"placeholder,
			{"type":"custom_tool_call_output","tool_call_id":"custom_1","output":["ok"]placeholder,
			{"type":"tool_search_call","id":"search_1","arguments":{"query":"repo"placeholderplaceholder,
			{"type":"tool_search_output","call_id":"search_1","results":{"groups":["repo"]placeholderplaceholder,
			{"role":"tool","tool_call_id":"orphan","content":{"value":1placeholderplaceholder
		]placeholder`)

	patched, err := sanitizeGrokResponsesModelInput(body)
placeholder
	generated := gjson.GetBytes(patched, "input.1.call_id").String()
	require.NotEmpty(t, generated)
	require.Equal(t, generated, gjson.GetBytes(patched, "input.0.call_id").String())
	require.JSONEq(t, `{"before":true,"large_id":9007199254740993placeholder`, gjson.GetBytes(patched, "input.0.output").String())
	require.Equal(t, "9007199254740993", gjson.Get(gjson.GetBytes(patched, "input.0.output").String(), "large_id").Raw)
	require.Equal(t, "lookup", gjson.GetBytes(patched, "input.1.name").String())
	require.JSONEq(t, `{"query":"x","large_id":9007199254740993placeholder`, gjson.GetBytes(patched, "input.1.arguments").String())
	require.Equal(t, "9007199254740993", gjson.Get(gjson.GetBytes(patched, "input.1.arguments").String(), "large_id").Raw)
	customCallID := gjson.GetBytes(patched, "input.2.call_id").String()
	require.NotEmpty(t, customCallID)
	require.NotEqual(t, "custom_1", customCallID)
	require.JSONEq(t, `{"input":"patch"placeholder`, gjson.GetBytes(patched, "input.2.arguments").String())
	require.Equal(t, customCallID, gjson.GetBytes(patched, "input.3.call_id").String())
	require.JSONEq(t, `["ok"]`, gjson.GetBytes(patched, "input.3.output").String())
	require.Equal(t, "tool_search", gjson.GetBytes(patched, "input.4.name").String())
	searchCallID := gjson.GetBytes(patched, "input.4.call_id").String()
	require.NotEmpty(t, searchCallID)
	require.NotEqual(t, "search_1", searchCallID)
	require.Equal(t, searchCallID, gjson.GetBytes(patched, "input.5.call_id").String())
	require.JSONEq(t, `{"groups":["repo"]placeholder`, gjson.GetBytes(patched, "input.5.output").String())
	require.Equal(t, "function_call_output", gjson.GetBytes(patched, "input.6.type").String())
	require.Equal(t, "orphan", gjson.GetBytes(patched, "input.6.call_id").String())
placeholder

func TestSanitizeGrokResponsesModelInputMapsAllItemAliasesAndRejectsConflicts(t *testing.T) {
	body := []byte(`{"input":[
		{"type":"function_call","id":"fc_alias","call_id":"call_canonical","name":"one","arguments":"{placeholder"placeholder,
		{"type":"function_call_output","id":"fc_alias","output":"ok"placeholder,
		{"type":"function_call","id":"duplicate_item","name":"two","arguments":"{placeholder"placeholder,
		{"type":"function_call","id":"duplicate_item","name":"three","arguments":"{placeholder"placeholder,
		{"type":"function_call_output","id":"duplicate_item","output":"ambiguous"placeholder
	]placeholder`)

	patched, err := sanitizeGrokResponsesModelInput(body)
placeholder
	require.Equal(t, "call_canonical", gjson.GetBytes(patched, "input.0.call_id").String())
	require.Equal(t, "call_canonical", gjson.GetBytes(patched, "input.1.call_id").String())
	require.False(t, gjson.GetBytes(patched, "input.0.id").Exists())
	secondCallID := gjson.GetBytes(patched, "input.2.call_id").String()
	thirdCallID := gjson.GetBytes(patched, "input.3.call_id").String()
	ambiguousOutputID := gjson.GetBytes(patched, "input.4.call_id").String()
	require.NotEmpty(t, secondCallID)
	require.NotEmpty(t, thirdCallID)
	require.NotEmpty(t, ambiguousOutputID)
	require.NotEqual(t, "duplicate_item", secondCallID)
	require.NotEqual(t, secondCallID, thirdCallID)
	require.NotEqual(t, secondCallID, ambiguousOutputID)
	require.NotEqual(t, thirdCallID, ambiguousOutputID)
placeholder

func TestPatchGrokResponsesBodyCombinesAdditionalToolsAndDropsOrphanControls(t *testing.T) {
	body := []byte(`{
		"input":[{"type":"additional_tools","tools":[{"type":"function","name":"lookup"placeholder]placeholder,{"role":"user","content":"hi"placeholder],
		"tools":[{"type":"function","name":"lookup"placeholder],
		"tool_choice":"auto","parallel_tool_calls":true
placeholder`)
	patched, err := patchGrokResponsesBody(body, "grok-4.5")
placeholder
	require.Equal(t, int64(1), gjson.GetBytes(patched, "tools.#").Int())
	require.Equal(t, "lookup", gjson.GetBytes(patched, "tools.0.name").String())
	require.Equal(t, "message", gjson.GetBytes(patched, "input.0.type").String())

	withoutTools, err := patchGrokResponsesBody([]byte(`{"input":"hi","tool_choice":"auto","parallel_tool_calls":trueplaceholder`), "grok-4.5")
placeholder
	require.False(t, gjson.GetBytes(withoutTools, "tool_choice").Exists())
	require.False(t, gjson.GetBytes(withoutTools, "parallel_tool_calls").Exists())

	for _, malformedTools := range []string{"null", `{placeholder`, `"invalid"`placeholder {
		malformedBody := []byte(`{"input":"hi","tools":` + malformedTools + `,"tool_choice":"auto","parallel_tool_calls":trueplaceholder`)
		patched, err := patchGrokResponsesBody(malformedBody, "grok-4.5")
	placeholder
		require.False(t, gjson.GetBytes(patched, "tool_choice").Exists(), string(patched))
		require.False(t, gjson.GetBytes(patched, "parallel_tool_calls").Exists(), string(patched))
placeholder
placeholder

func TestSanitizeGrokCompactionReplayBodyPreservesVisibleSummary(t *testing.T) {
	body := []byte(`{
		"previous_response_id":"resp_stale",
		"input":[{"type":"compaction","encrypted_content":"opaque","summary":[{"type":"summary_text","text":"visible history"placeholder]placeholder,{"type":"message","role":"user","content":"continue"placeholder]
placeholder`)
	patched, changed, err := sanitizeGrokCompactionReplayBody(body)
placeholder
	require.True(t, changed)
	require.False(t, gjson.GetBytes(patched, "previous_response_id").Exists())
	require.False(t, gjson.GetBytes(patched, `input.#(type=="reasoning")`).Exists())
	require.Contains(t, gjson.GetBytes(patched, "input.0.content.0.text").String(), "visible history")
	require.Equal(t, "continue", gjson.GetBytes(patched, "input.1.content").String())
	require.True(t, isGrokCompactionReplayDecodeError(http.StatusBadRequest, []byte(`{"error":{"message":"failed to deserialize compaction summary history"placeholderplaceholder`)))
	require.True(t, isGrokCompactionReplayDecodeError(http.StatusBadRequest, []byte(`{"error":{"type":"invalid_request_error"placeholder,"detail":"could not decode compaction blob"placeholder`)))
placeholder

func TestGrokStructuredErrorCandidatesDoNotShadowTopLevelMessages(t *testing.T) {
	shadowedInvalidEncrypted := []byte(`{"code":"invalid-argument","error":{"type":"invalid_request_error"placeholder,"message":"Could not decrypt encrypted_content because it was modified"placeholder`)
	require.True(t, isGrokInvalidEncryptedContentResponse(http.StatusBadRequest, shadowedInvalidEncrypted))
	require.True(t, isGrokCompactionReplayDecodeError(http.StatusBadRequest, []byte(`{"error":{"type":"invalid_request_error"placeholder,"message":"could not decode compaction history"placeholder`)))
placeholder

func TestGrokDecoderCompatibility422FailsOverWithoutCooldown(t *testing.T) {
	body := []byte(`{"detail":"data did not match any variant of untagged enum ModelInput at input[3]"placeholder`)
	require.True(t, isGrokDecoderCompatibilityError(http.StatusUnprocessableEntity, body))
	require.True(t, isGrokDecoderCompatibilityError(http.StatusUnprocessableEntity, []byte(`{"message":"could not decode ModelInput at input.3"placeholder`)))
	require.True(t, isGrokDecoderCompatibilityError(http.StatusUnprocessableEntity, []byte(`{"error":{"type":"invalid_request_error"placeholder,"message":"could not deserialize ModelInput at input[3]"placeholder`)))
	require.True(t, isGrokDecoderCompatibilityError(http.StatusUnprocessableEntity, []byte(`{"error":"Failed to deserialize the JSON body into the target type: messages[1]: data did not match any variant of untagged enum Content at line 1 column 6577"placeholder`)))
	require.True(t, (&OpenAIGatewayService{placeholder).shouldFailoverGrokUpstreamError(http.StatusUnprocessableEntity, body))
	decision := classifyGrokUpstreamFailure(http.StatusUnprocessableEntity, body, "grok-4.5")
	require.False(t, decision.ShouldCooldown)
	require.Equal(t, GrokFailureNone, decision.Class)

	require.False(t, isGrokDecoderCompatibilityError(http.StatusUnprocessableEntity, []byte(`{"error":{"message":"invalid user parameter"placeholderplaceholder`)))
	require.False(t, isGrokDecoderCompatibilityError(http.StatusUnprocessableEntity, []byte(`{"error":"messages[1].content is required"placeholder`)))
	require.False(t, isGrokDecoderCompatibilityError(http.StatusBadRequest, body))
placeholder
