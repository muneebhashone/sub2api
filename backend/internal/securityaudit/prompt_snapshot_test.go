package securityaudit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

func TestExtractPromptSnapshotProtocols(t *testing.T) {
	tests := []struct {
		protocol, body, first string
		count                 int
placeholder{
		{"openai_chat_completions", `{"messages":[{"role":"user","content":"old"placeholder,{"role":"assistant","content":"assistant turn"placeholder,{"role":"user","content":[{"type":"text","text":"最新😀"placeholder]placeholder]placeholder`, "最新😀", 3placeholder,
		{"openai_responses", `{"input":[{"role":"user","content":[{"type":"input_text","text":"response text"placeholder]placeholder]placeholder`, "response text", 1placeholder,
		{"anthropic_messages", `{"messages":[{"role":"user","content":[{"type":"text","text":"claude"placeholder]placeholder]placeholder`, "claude", 1placeholder,
		{"gemini", `{"contents":[{"role":"user","parts":[{"text":"gemini"placeholder,{"inline_data":{"data":"BASE64"placeholderplaceholder]placeholder]placeholder`, "gemini", 1placeholder,
		{"openai_images", `{"prompt":"draw a cat","image":"BASE64SECRET"placeholder`, "draw a cat", 1placeholder,
		{"responses_websocket", `{"type":"response.create","response":{"input":"turn two"placeholderplaceholder`, "turn two", 1placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.protocol, func(t *testing.T) {
			snapshot, err := ExtractPromptSnapshot(Request{Protocol: tt.protocol, Body: []byte(tt.body), Stage: "http"placeholder)
		placeholder
			require.True(t, strings.HasPrefix(snapshot.ScanText, tt.first))
			require.Equal(t, tt.count, snapshot.MessageCount)
			require.Equal(t, utf8.RuneCountInString(metadataTextForTest(snapshot.ScanText)), snapshot.PromptLength)
			require.NotEmpty(t, snapshot.PromptHash)
			require.NotContains(t, snapshot.ScanText, "BASE64SECRET")
	placeholder)
placeholder
placeholder

func TestSnapshotRedactsCanariesAndPreservesHashOfScanText(t *testing.T) {
	body := `{"messages":[{"role":"user","content":"PROMPT_CANARY_ABC123 email@example.com +86 138 0013 8000 Bearer AUTH_CANARY_XYZ sk-secretvalue123 password=supersecret123"placeholder]placeholder`
	snapshot, err := ExtractPromptSnapshot(Request{Protocol: "openai_chat_completions", Body: []byte(body)placeholder)
placeholder
	require.NotContains(t, snapshot.RedactedPreview, "ABC123")
	require.NotContains(t, snapshot.RedactedPreview, "email@example.com")
	require.NotContains(t, snapshot.RedactedPreview, "AUTH_CANARY_XYZ")
	require.NotContains(t, snapshot.RedactedPreview, "secretvalue123")
	require.NotContains(t, snapshot.RedactedPreview, "supersecret123")
	require.NotContains(t, snapshot.RedactedPreview, "138 0013 8000")
	require.Contains(t, snapshot.ScanText, "PROMPT_CANARY_ABC123")
	require.NotEqual(t, snapshot.ScanText, snapshot.RedactedPreview)
	digest := sha256.Sum256([]byte(metadataTextForTest(snapshot.ScanText)))
	require.Equal(t, hex.EncodeToString(digest[:]), snapshot.PromptHash)
	require.Empty(t, snapshot.Redacted().ScanText)
placeholder

func TestSplitRunesDoesNotSplitUTF8(t *testing.T) {
	chunks := SplitRunes("中文😀éabc", 2)
	require.Equal(t, []string{"中文", "😀e", "́a", "bc"placeholder, chunks)
	for _, chunk := range chunks {
		require.True(t, utf8.ValidString(chunk))
placeholder
	require.Equal(t, "中文😀éabc", strings.Join(chunks, ""))
placeholder

func TestSplitRunesKeepsPrioritySegmentIndependent(t *testing.T) {
	latest := "请帮我编写一篇黄色小说 名字你来取"
	history := strings.Repeat("AGENTS.md 项目约束。", 40)
	chunks := SplitRunes(latest+promptAuditPrioritySeparator+history, 128)
	require.Greater(t, len(chunks), 2)
	require.Equal(t, latest, chunks[0])
	require.Equal(t, history, strings.Join(chunks[1:], ""))
	for _, chunk := range chunks {
		require.NotContains(t, chunk, promptAuditPrioritySeparator)
placeholder
placeholder

func TestPromptSnapshotLatestUserTextBlockIsOnePrioritizedSegment(t *testing.T) {
	body := []byte(`{
		"messages":[
			{"role":"user","content":"历史输入"placeholder,
			{"role":"assistant","content":"assistant client injection"placeholder,
			{"role":"tool","content":"tool client injection"placeholder,
			{"role":"user","content":[
				{"type":"text","text":"最新第一块😀"placeholder,
				{"type":"image_url","image_url":{"url":"data:image/png;base64,IMAGE_CANARY_BASE64"placeholderplaceholder,
				{"type":"text","text":"最新第二块é"placeholder
			]placeholder
		]
placeholder`)
	snapshot, err := ExtractPromptSnapshot(Request{Protocol: "openai_chat_completions", Body: bodyplaceholder)
placeholder
	require.Equal(t, 5, snapshot.MessageCount)
	require.True(t, strings.HasPrefix(snapshot.ScanText, "最新第二块é"+promptAuditPrioritySeparator))
	require.Contains(t, snapshot.ScanText, "最新第一块😀")
	require.Contains(t, snapshot.ScanText, "历史输入")
	require.Contains(t, snapshot.ScanText, "assistant client injection")
	require.Contains(t, snapshot.ScanText, "tool client injection")
	require.NotContains(t, snapshot.ScanText, "IMAGE_CANARY_BASE64")
	require.Equal(t, utf8.RuneCountInString(metadataTextForTest(snapshot.ScanText)), snapshot.PromptLength)
placeholder

func TestPromptSnapshotSeparatesAnthropicUserPromptFromHarnessBlocks(t *testing.T) {
	latest := "请帮我编写一篇黄色小说 名字你来取"
	agents := "# AGENTS.md instructions\n<INSTRUCTIONS>" + strings.Repeat("安全约束。", 80) + "</INSTRUCTIONS>"
	environment := "<environment_context><cwd>/workspace</cwd></environment_context>"
	body := []byte(`{"system":"system policy","messages":[{"role":"user","content":[` +
		`{"type":"text","text":` + string(mustJSON(t, agents)) + `placeholder,` +
		`{"type":"text","text":` + string(mustJSON(t, environment)) + `placeholder,` +
		`{"type":"text","text":` + string(mustJSON(t, latest)) + `placeholder` +
		`]placeholder]placeholder`)

	snapshot, err := ExtractPromptSnapshot(Request{Protocol: "anthropic_messages", Body: bodyplaceholder)
placeholder
	require.Equal(t, 4, snapshot.MessageCount)
	require.True(t, strings.HasPrefix(snapshot.ScanText, latest+promptAuditPrioritySeparator))
	require.True(t, strings.HasPrefix(snapshot.RedactedPreview, "请帮我编写一篇黄色小说"))

	chunks := SplitRunes(snapshot.ScanText, 128)
	require.Equal(t, latest, chunks[0])
	require.Contains(t, strings.Join(chunks[1:], ""), "# AGENTS.md instructions")
	require.Contains(t, strings.Join(chunks[1:], ""), "<environment_context>")
	require.NotContains(t, strings.Join(chunks, ""), promptAuditPrioritySeparator)
placeholder

func TestPromptSnapshotResponsesShapes(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
placeholder{
		{name: "string", body: `{"input":"plain response input"placeholder`, want: "plain response input"placeholder,
		{name: "message array", body: `{"input":[{"role":"assistant","content":"assistant turn"placeholder,{"role":"user","content":[{"type":"input_text","text":"message block"placeholder]placeholder]placeholder`, want: "message block\n\nassistant turn"placeholder,
		{name: "direct input text", body: `{"input":[{"type":"input_text","text":"direct block"placeholder]placeholder`, want: "direct block"placeholder,
		{name: "single object", body: `{"input":{"role":"user","content":[{"type":"input_text","text":"single object"placeholder]placeholderplaceholder`, want: "single object"placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshot, err := ExtractPromptSnapshot(Request{Protocol: "openai_responses", Body: []byte(tt.body)placeholder)
		placeholder
			require.Equal(t, tt.want, metadataTextForTest(snapshot.ScanText))
	placeholder)
placeholder
placeholder

func TestPromptSnapshotGeminiBatchShapesAndMediaExclusion(t *testing.T) {
	body := []byte(`{
		"contents":{"role":"user","parts":[{"text":"root content"placeholder,{"inlineData":{"data":"ROOT_BASE64"placeholderplaceholder]placeholder,
		"instances":[{"prompt":"instance prompt"placeholder],
		"requests":[
			{"contents":[{"role":"model","parts":[{"text":"ignore model"placeholder]placeholder,{"role":"user","parts":[{"text":"nested user"placeholder]placeholder]placeholder,
			{"instances":[{"prompt":"nested instance"placeholder]placeholder
		]
placeholder`)
	snapshot, err := ExtractPromptSnapshot(Request{Protocol: "gemini", Body: bodyplaceholder)
placeholder
	require.True(t, strings.HasPrefix(snapshot.ScanText, "nested instance"))
	for _, expected := range []string{"root content", "instance prompt", "nested user", "nested instance"placeholder {
		require.Contains(t, snapshot.ScanText, expected)
placeholder
	require.NotContains(t, snapshot.ScanText, "ROOT_BASE64")
	require.Contains(t, snapshot.ScanText, "ignore model")
placeholder

func TestPromptSnapshotMediaOnlyExtractsDeterministicTextPrompts(t *testing.T) {
	body := []byte(`{
		"prompt":"draw a lighthouse",
		"image":"data:image/png;base64,IMAGE_CANARY",
		"input":{"negative_prompt":"no fog","image_prompt":"https://example.test/input.png","prompt":"draw a lighthouse"placeholder,
		"request":{"lyrics":"ocean song","input":"` + strings.Repeat("A", 300) + `"placeholder,
		"images":[{"description":"nested textual direction","image_url":"https://example.test/image.png"placeholder]
placeholder`)
	snapshot, err := ExtractPromptSnapshot(Request{Protocol: "grok_media", Body: bodyplaceholder)
placeholder
	require.Equal(t, 4, snapshot.MessageCount)
	for _, expected := range []string{"draw a lighthouse", "no fog", "ocean song", "nested textual direction"placeholder {
		require.Contains(t, snapshot.ScanText, expected)
placeholder
	require.Equal(t, 1, strings.Count(snapshot.ScanText, "draw a lighthouse"))
	require.NotContains(t, snapshot.ScanText, "IMAGE_CANARY")
	require.NotContains(t, snapshot.ScanText, "example.test")
	require.NotContains(t, snapshot.ScanText, strings.Repeat("A", 100))
placeholder

func TestResponsesWebSocketOnlyAuditsResponseCreateAndPreservesStage(t *testing.T) {
	for _, stage := range []string{"first_turn", "subsequent_turn"placeholder {
		snapshot, err := ExtractPromptSnapshot(Request{
			Protocol: "openai_responses", Stage: stage,
			Body: []byte(`{"type":"response.create","response":{"model":"gpt-test","input":[{"role":"user","content":[{"type":"input_text","text":"ws turn"placeholder]placeholder]placeholderplaceholder`),
	placeholder)
	placeholder
		require.Equal(t, "ws turn", snapshot.ScanText)
		require.Equal(t, stage, snapshot.Stage)
placeholder
	_, err := ExtractPromptSnapshot(Request{
		Protocol: "openai_responses", Stage: "subsequent_turn",
		Body: []byte(`{"type":"conversation.item.create","response":{"input":"must not scan this frame"placeholderplaceholder`),
placeholder)
	require.True(t, errors.Is(err, ErrNoPromptText))
placeholder

func TestPromptSnapshotEmptyAndLongUnicodeInput(t *testing.T) {
	_, err := ExtractPromptSnapshot(Request{Protocol: "openai_chat_completions", Body: []byte(`{"messages":[{"role":"function","content":"not audited role"placeholder,{"role":"user","content":"  "placeholder]placeholder`)placeholder)
	require.True(t, errors.Is(err, ErrNoPromptText))

	latest := strings.Repeat("最新😀é", 80)
	history := strings.Repeat("历史中文", 80)
	body := []byte(`{"messages":[{"role":"user","content":` + string(mustJSON(t, history)) + `placeholder,{"role":"user","content":` + string(mustJSON(t, latest)) + `placeholder]placeholder`)
	snapshot, err := ExtractPromptSnapshot(Request{Protocol: "openai_chat_completions", Body: bodyplaceholder)
placeholder
	require.True(t, strings.HasPrefix(snapshot.ScanText, latest))
	chunks := SplitRunes(snapshot.ScanText, 127)
	require.Equal(t, strings.Replace(snapshot.ScanText, promptAuditPrioritySeparator, "", 1), strings.Join(chunks, ""))
	require.Equal(t, latest, chunks[0]+strings.Join(chunks[1:len(SplitRunes(latest, 127))], ""))
	for _, chunk := range chunks {
		require.LessOrEqual(t, len([]rune(chunk)), 127)
		require.True(t, utf8.ValidString(chunk))
placeholder
placeholder

func TestPromptSnapshotIncludesClientControlledInstructions(t *testing.T) {
	tests := []struct {
		name, protocol, body string
		want                 []string
placeholder{
		{
			name:     "openai system developer assistant tool",
			protocol: "openai_chat_completions",
			body:     `{"messages":[{"role":"system","content":"system jailbreak"placeholder,{"role":"developer","content":"developer policy"placeholder,{"role":"assistant","content":"assistant jailbreak"placeholder,{"role":"tool","content":"tool payload"placeholder,{"role":"user","content":"hello"placeholder]placeholder`,
			want:     []string{"system jailbreak", "developer policy", "assistant jailbreak", "tool payload", "hello"placeholder,
	placeholder,
		{
			name:     "openai system only",
			protocol: "openai_chat_completions",
			body:     `{"messages":[{"role":"system","content":"only system instruction"placeholder]placeholder`,
			want:     []string{"only system instruction"placeholder,
	placeholder,
		{
			name:     "responses instructions",
			protocol: "openai_responses",
			body:     `{"instructions":"response instructions","input":[{"role":"user","content":[{"type":"input_text","text":"user turn"placeholder]placeholder]placeholder`,
			want:     []string{"response instructions", "user turn"placeholder,
	placeholder,
		{
			name:     "anthropic system",
			protocol: "anthropic_messages",
			body:     `{"system":"claude system","messages":[{"role":"user","content":[{"type":"text","text":"claude user"placeholder]placeholder]placeholder`,
			want:     []string{"claude system", "claude user"placeholder,
	placeholder,
		{
			name:     "gemini systemInstruction",
			protocol: "gemini",
			body:     `{"systemInstruction":{"parts":[{"text":"gemini system"placeholder]placeholder,"contents":[{"role":"user","parts":[{"text":"gemini user"placeholder]placeholder]placeholder`,
			want:     []string{"gemini system", "gemini user"placeholder,
	placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshot, err := ExtractPromptSnapshot(Request{Protocol: tt.protocol, Body: []byte(tt.body)placeholder)
		placeholder
			for _, expected := range tt.want {
				require.Contains(t, snapshot.ScanText, expected)
		placeholder
	placeholder)
placeholder
placeholder

func TestBuildPromptPreviewWithholdsMajorityOfOrdinaryText(t *testing.T) {
	prompt := strings.Repeat("机密业务提示词内容", 40)
	preview := BuildPromptPreview(prompt, DefaultPromptPreviewMaxRunes)
	require.NotEmpty(t, preview)
	require.Contains(t, preview, "***")
	require.LessOrEqual(t, utf8.RuneCountInString(strings.TrimSuffix(strings.TrimSuffix(preview, "…"), "***")), 24)
	require.Less(t, utf8.RuneCountInString(preview), utf8.RuneCountInString(prompt)/2)
	require.NotContains(t, preview, prompt)
placeholder

func TestBuildPromptPreviewFullyMasksShortUnlabelledSecrets(t *testing.T) {
	require.Equal(t, "***", BuildPromptPreview("short-secret-value!!", DefaultPromptPreviewMaxRunes))
	require.Equal(t, "***", BuildPromptPreview(strings.Repeat("a", 31), DefaultPromptPreviewMaxRunes))
	partial := BuildPromptPreview(strings.Repeat("b", 32), DefaultPromptPreviewMaxRunes)
	require.True(t, strings.HasPrefix(partial, "b"))
	require.Contains(t, partial, "***")
placeholder

func mustJSON(t *testing.T, value string) []byte {
placeholder
	raw, err := json.Marshal(value)
placeholder
	return raw
placeholder

func metadataTextForTest(scanText string) string {
	return strings.Replace(scanText, promptAuditPrioritySeparator, "\n\n", 1)
placeholder
