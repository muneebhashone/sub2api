package service

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestDeriveOpenAIContentSessionSeed_EmptyInputs(t *testing.T) {
	require.Empty(t, deriveOpenAIContentSessionSeed(nil))
	require.Empty(t, deriveOpenAIContentSessionSeed([]byte{placeholder))
	require.Empty(t, deriveOpenAIContentSessionSeed([]byte(`{placeholder`)))
placeholder

func TestDeriveOpenAIContentSessionSeed_ModelOnly(t *testing.T) {
	seed := deriveOpenAIContentSessionSeed([]byte(`{"model":"gpt-5.4"placeholder`))
	require.Contains(t, seed, contentSessionSeedPrefix)
	require.Contains(t, seed, "model=gpt-5.4")
placeholder

func TestDeriveOpenAIContentSessionSeed_ChatCompletions_StableAcrossTurns(t *testing.T) {
	turn1 := []byte(`{
		"model": "gpt-5.4",
		"messages": [
			{"role": "system", "content": "You are helpful."placeholder,
			{"role": "user", "content": "Hello"placeholder
		]
placeholder`)
	turn2 := []byte(`{
		"model": "gpt-5.4",
		"messages": [
			{"role": "system", "content": "You are helpful."placeholder,
			{"role": "user", "content": "Hello"placeholder,
			{"role": "assistant", "content": "Hi there!"placeholder,
			{"role": "user", "content": "How are you?"placeholder
		]
placeholder`)
	s1 := deriveOpenAIContentSessionSeed(turn1)
	s2 := deriveOpenAIContentSessionSeed(turn2)
	require.Equal(t, s1, s2, "seed should be stable across later turns")
	require.NotEmpty(t, s1)
placeholder

func TestDeriveOpenAIContentSessionSeed_ChatCompletions_DifferentFirstUserDiffers(t *testing.T) {
	req1 := []byte(`{"model":"gpt-5.4","messages":[{"role":"user","content":"Question A"placeholder]placeholder`)
	req2 := []byte(`{"model":"gpt-5.4","messages":[{"role":"user","content":"Question B"placeholder]placeholder`)
	s1 := deriveOpenAIContentSessionSeed(req1)
	s2 := deriveOpenAIContentSessionSeed(req2)
	require.NotEqual(t, s1, s2)
placeholder

func TestDeriveOpenAIContentSessionSeed_ChatCompletions_DifferentSystemDiffers(t *testing.T) {
	req1 := []byte(`{"model":"gpt-5.4","messages":[{"role":"system","content":"A"placeholder,{"role":"user","content":"Hi"placeholder]placeholder`)
	req2 := []byte(`{"model":"gpt-5.4","messages":[{"role":"system","content":"B"placeholder,{"role":"user","content":"Hi"placeholder]placeholder`)
	s1 := deriveOpenAIContentSessionSeed(req1)
	s2 := deriveOpenAIContentSessionSeed(req2)
	require.NotEqual(t, s1, s2)
placeholder

func TestDeriveOpenAIContentSessionSeed_ChatCompletions_DifferentModelDiffers(t *testing.T) {
	req1 := []byte(`{"model":"gpt-5.4","messages":[{"role":"user","content":"Hi"placeholder]placeholder`)
	req2 := []byte(`{"model":"gpt-4o","messages":[{"role":"user","content":"Hi"placeholder]placeholder`)
	s1 := deriveOpenAIContentSessionSeed(req1)
	s2 := deriveOpenAIContentSessionSeed(req2)
	require.NotEqual(t, s1, s2)
placeholder

func TestDeriveOpenAIContentSessionSeed_ChatCompletions_WithTools(t *testing.T) {
	withTools := []byte(`{
		"model": "gpt-5.4",
		"tools": [{"type":"function","function":{"name":"get_weather"placeholderplaceholder],
		"messages": [{"role": "user", "content": "Hello"placeholder]
placeholder`)
	withoutTools := []byte(`{
		"model": "gpt-5.4",
		"messages": [{"role": "user", "content": "Hello"placeholder]
placeholder`)
	s1 := deriveOpenAIContentSessionSeed(withTools)
	s2 := deriveOpenAIContentSessionSeed(withoutTools)
	require.NotEqual(t, s1, s2, "tools should affect the seed")
	require.Contains(t, s1, "|tools=")
placeholder

func TestDeriveOpenAIContentSessionSeed_ChatCompletions_WithFunctions(t *testing.T) {
	body := []byte(`{
		"model": "gpt-5.4",
		"functions": [{"name":"get_weather","parameters":{placeholderplaceholder],
		"messages": [{"role": "user", "content": "Hello"placeholder]
placeholder`)
	seed := deriveOpenAIContentSessionSeed(body)
	require.Contains(t, seed, "|functions=")
placeholder

func TestDeriveOpenAIContentSessionSeed_ChatCompletions_DeveloperRole(t *testing.T) {
	body := []byte(`{
		"model": "gpt-5.4",
		"messages": [
			{"role": "developer", "content": "You are helpful."placeholder,
			{"role": "user", "content": "Hello"placeholder
		]
placeholder`)
	seed := deriveOpenAIContentSessionSeed(body)
	require.Contains(t, seed, "|system=")
	require.Contains(t, seed, "|first_user=")
placeholder

func TestDeriveOpenAIContentSessionSeed_ChatCompletions_StructuredContent(t *testing.T) {
	body := []byte(`{
		"model": "gpt-5.4",
		"messages": [
			{"role": "user", "content": [{"type":"text","text":"Hello"placeholder]placeholder
		]
placeholder`)
	seed := deriveOpenAIContentSessionSeed(body)
	require.NotEmpty(t, seed)
	require.Contains(t, seed, "|first_user=")
placeholder

func TestDeriveOpenAIContentSessionSeed_ResponsesAPI_InputString(t *testing.T) {
	body := []byte(`{"model":"gpt-5.4","input":"Hello, how are you?"placeholder`)
	seed := deriveOpenAIContentSessionSeed(body)
	require.Contains(t, seed, "|input=Hello, how are you?")
placeholder

func TestDeriveOpenAIContentSessionSeed_ResponsesAPI_InputArray(t *testing.T) {
	body := []byte(`{
		"model": "gpt-5.4",
		"input": [
			{"role": "system", "content": "You are helpful."placeholder,
			{"role": "user", "content": "Hello"placeholder
		]
placeholder`)
	seed := deriveOpenAIContentSessionSeed(body)
	require.Contains(t, seed, "|system=")
	require.Contains(t, seed, "|first_user=")
placeholder

func TestDeriveOpenAIContentSessionSeed_ResponsesAPI_WithInstructions(t *testing.T) {
	body := []byte(`{
		"model": "gpt-5.4",
		"instructions": "You are a coding assistant.",
		"input": "Write a hello world"
placeholder`)
	seed := deriveOpenAIContentSessionSeed(body)
	require.Contains(t, seed, "|instructions=You are a coding assistant.")
	require.Contains(t, seed, "|input=Write a hello world")
placeholder

func TestDeriveOpenAIContentSessionSeed_Deterministic(t *testing.T) {
	body := []byte(`{
		"model": "gpt-5.4",
		"messages": [
			{"role": "system", "content": "You are helpful."placeholder,
			{"role": "user", "content": "Hello"placeholder
		]
placeholder`)
	s1 := deriveOpenAIContentSessionSeed(body)
	s2 := deriveOpenAIContentSessionSeed(body)
	require.Equal(t, s1, s2, "seed must be deterministic")
placeholder

func TestDeriveOpenAIContentSessionSeed_PrefixPresent(t *testing.T) {
	body := []byte(`{"model":"gpt-5.4","messages":[{"role":"user","content":"Hi"placeholder]placeholder`)
	seed := deriveOpenAIContentSessionSeed(body)
	require.True(t, len(seed) > len(contentSessionSeedPrefix))
	require.Equal(t, contentSessionSeedPrefix, seed[:len(contentSessionSeedPrefix)])
placeholder

func TestDeriveOpenAIContentSessionSeed_EmptyToolsIgnored(t *testing.T) {
	body := []byte(`{"model":"gpt-5.4","tools":[],"messages":[{"role":"user","content":"Hi"placeholder]placeholder`)
	seed := deriveOpenAIContentSessionSeed(body)
	require.NotContains(t, seed, "|tools=")
placeholder

func TestDeriveOpenAIContentSessionSeed_MessagesPreferredOverInput(t *testing.T) {
	body := []byte(`{
		"model": "gpt-5.4",
		"messages": [{"role": "user", "content": "from messages"placeholder],
		"input": "from input"
placeholder`)
	seed := deriveOpenAIContentSessionSeed(body)
	require.Contains(t, seed, "|first_user=")
	require.NotContains(t, seed, "|input=")
placeholder

func TestDeriveOpenAIContentSessionSeed_JSONCanonicalisation(t *testing.T) {
	compact := []byte(`{"model":"gpt-5.4","tools":[{"type":"function","function":{"name":"get_weather","description":"Get weather"placeholderplaceholder],"messages":[{"role":"user","content":"Hi"placeholder]placeholder`)
	spaced := []byte(`{
		"model": "gpt-5.4",
		"tools": [
			{ "type" : "function", "function": { "description": "Get weather", "name": "get_weather" placeholder placeholder
		],
		"messages": [ { "role": "user", "content": "Hi" placeholder ]
placeholder`)
	s1 := deriveOpenAIContentSessionSeed(compact)
	s2 := deriveOpenAIContentSessionSeed(spaced)
	require.Equal(t, s1, s2, "different formatting of identical JSON should produce the same seed")
placeholder

func TestDeriveOpenAIContentSessionSeed_SingleScanMatchesLegacyBytes(t *testing.T) {
	largeValue := strings.Repeat("payload", 1<<17)
	tests := []struct {
		name string
		body []byte
placeholder{
		{
			name: "large chat completions",
			body: []byte(`{"metadata":"` + largeValue + `","model":"gpt-5.4","tools":[{"type":"function","function":{"name":"lookup"placeholderplaceholder],"functions":[{"name":"legacy_lookup"placeholder],"messages":[{"role":"system","content":"System prompt"placeholder,{"role":"developer","content":[{"type":"text","text":"Developer prompt"placeholder]placeholder,{"role":"user","content":"Hello"placeholder]placeholder`),
	placeholder,
		{
			name: "large responses",
			body: []byte(`{"metadata":"` + largeValue + `","model":"gpt-5.4","instructions":"Be concise.","tools":[{"type":"function","name":"lookup"placeholder],"input":[{"role":"system","content":"System prompt"placeholder,{"role":"user","content":[{"type":"input_text","text":"Hello"placeholder]placeholder]placeholder`),
	placeholder,
		{
			name: "fields in reverse order",
			body: []byte(`{"input":"fallback input","messages":[{"role":"user","content":"chat wins"placeholder],"instructions":"Be concise.","functions":[{"name":"lookup"placeholder],"tools":[{"type":"function","name":"lookup"placeholder],"model":"gpt-5.4"placeholder`),
	placeholder,
		{
			name: "missing and wrong type fields",
			body: []byte(`{"tools":[],"functions":null,"instructions":0,"messages":{placeholder,"input":[{"type":"input_text","text":"fallback"placeholder]placeholder`),
	placeholder,
		{
			name: "duplicate fields keep first value",
			body: []byte(`{"model":"first","model":"second","tools":[{"name":"first"placeholder],"tools":[{"name":"second"placeholder],"functions":[],"functions":[{"name":"second"placeholder],"instructions":"first","instructions":"second","messages":null,"messages":[{"role":"user","content":"second"placeholder],"input":"first input","input":"second input"placeholder`),
	placeholder,
		{
			name: "escaped field names",
			body: []byte(`{"mo\u0064el":"gpt-5.4","mess\u0061ges":[{"role":"user","content":"Hello"placeholder]placeholder`),
	placeholder,
		{
			name: "trailing object fields are outside the root",
			body: []byte(`{"foo":1placeholder{"model":"trailing","input":"trailing input"placeholder`),
	placeholder,
		{
			name: "trailing quoted fields are outside the root",
			body: []byte(`{"model":"root"placeholder"input":"trailing input"`),
	placeholder,
		{
			name: "leading garbage before the root",
			body: []byte(`garbage{"model":"gpt-5.4","input":"Hello"placeholder`),
	placeholder,
		{
			name: "escaped braces remain inside string values",
			body: []byte(`{"metadata":"escaped placeholder and [ and \" quote","model":"root","input":"Hello"placeholder{"model":"trailing"placeholder`),
	placeholder,
		{
			name: "nested braces do not end the root",
			body: []byte(`{"metadata":{"nested":"placeholder ]"placeholder,"model":"root","input":"Hello"placeholder{"model":"trailing"placeholder`),
	placeholder,
		{
			name: "root array does not expose nested or trailing object fields",
			body: []byte(`[{"model":"nested"placeholder]{"model":"trailing","input":"trailing input"placeholder`),
	placeholder,
		{
			name: "trailing messages do not override root input",
			body: []byte(`{"input":"root input"placeholder{"messages":[{"role":"user","content":"trailing"placeholder]placeholder`),
	placeholder,
		{
			name: "truncated string containing a closing brace",
			body: []byte(`{"model":"root","metadata":"still placeholder inside`),
	placeholder,
		{
			name: "lenient truncated body",
			body: []byte(`{"model":"gpt-5.4","messages":[{"role":"user","content":"Hello"placeholder]`),
	placeholder,
placeholder

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, legacyDeriveOpenAIContentSessionSeed(test.body), deriveOpenAIContentSessionSeed(test.body))
	placeholder)
placeholder
placeholder

func TestDeriveOpenAIContentSessionSeed_AllTruncationOffsetsMatchLegacyBytes(t *testing.T) {
	bodies := []string{
		`{"model":"gpt-5.4","tools":[{"type":"function","function":{"name":"lookup"placeholderplaceholder],"functions":[{"name":"legacy"placeholder],"instructions":"escaped \" placeholder text","messages":[{"role":"system","content":"System"placeholder,{"role":"user","content":[{"type":"text","text":"Hello"placeholder]placeholder],"input":"fallback"placeholder`,
		`{"model":"gpt-5.4","instructions":"Be concise.","tools":[{"type":"function","name":"lookup"placeholder],"input":[{"role":"system","content":"System"placeholder,{"role":"user","content":[{"type":"input_text","text":"Hello"placeholder]placeholder]placeholder`,
placeholder
	for bodyIndex, body := range bodies {
		for end := 1; end < len(body); end++ {
			truncated := []byte(body[:end])
			require.Equalf(t, legacyDeriveOpenAIContentSessionSeed(truncated), deriveOpenAIContentSessionSeed(truncated), "body %d truncated at byte %d", bodyIndex, end)
	placeholder
placeholder
placeholder

func legacyDeriveOpenAIContentSessionSeed(body []byte) string {
	if len(body) == 0 {
		return ""
placeholder

	var b strings.Builder

	if model := gjson.GetBytes(body, "model").String(); model != "" {
		_, _ = b.WriteString("model=")
		_, _ = b.WriteString(model)
placeholder

	if tools := gjson.GetBytes(body, "tools"); tools.Exists() && tools.IsArray() && tools.Raw != "[]" {
		_, _ = b.WriteString("|tools=")
		_, _ = b.WriteString(normalizeCompatSeedJSON(json.RawMessage(tools.Raw)))
placeholder

	if funcs := gjson.GetBytes(body, "functions"); funcs.Exists() && funcs.IsArray() && funcs.Raw != "[]" {
		_, _ = b.WriteString("|functions=")
		_, _ = b.WriteString(normalizeCompatSeedJSON(json.RawMessage(funcs.Raw)))
placeholder

	if instr := gjson.GetBytes(body, "instructions").String(); instr != "" {
		_, _ = b.WriteString("|instructions=")
		_, _ = b.WriteString(instr)
placeholder

	firstUserCaptured := false

	msgs := gjson.GetBytes(body, "messages")
	if msgs.Exists() && msgs.IsArray() {
		msgs.ForEach(func(_, msg gjson.Result) bool {
			role := msg.Get("role").String()
			switch role {
			case "system", "developer":
				_, _ = b.WriteString("|system=")
				if c := msg.Get("content"); c.Exists() {
					_, _ = b.WriteString(normalizeCompatSeedJSON(json.RawMessage(c.Raw)))
			placeholder
			case "user":
				if !firstUserCaptured {
					_, _ = b.WriteString("|first_user=")
					if c := msg.Get("content"); c.Exists() {
						_, _ = b.WriteString(normalizeCompatSeedJSON(json.RawMessage(c.Raw)))
				placeholder
					firstUserCaptured = true
			placeholder
		placeholder
			return true
	placeholder)
placeholder else if inp := gjson.GetBytes(body, "input"); inp.Exists() {
		if inp.Type == gjson.String {
			_, _ = b.WriteString("|input=")
			_, _ = b.WriteString(inp.String())
	placeholder else if inp.IsArray() {
			inp.ForEach(func(_, item gjson.Result) bool {
				role := item.Get("role").String()
				switch role {
				case "system", "developer":
					_, _ = b.WriteString("|system=")
					if c := item.Get("content"); c.Exists() {
						_, _ = b.WriteString(normalizeCompatSeedJSON(json.RawMessage(c.Raw)))
				placeholder
				case "user":
					if !firstUserCaptured {
						_, _ = b.WriteString("|first_user=")
						if c := item.Get("content"); c.Exists() {
							_, _ = b.WriteString(normalizeCompatSeedJSON(json.RawMessage(c.Raw)))
					placeholder
						firstUserCaptured = true
				placeholder
			placeholder
				if !firstUserCaptured && item.Get("type").String() == "input_text" {
					_, _ = b.WriteString("|first_user=")
					if text := item.Get("text").String(); text != "" {
						_, _ = b.WriteString(text)
				placeholder
					firstUserCaptured = true
			placeholder
				return true
		placeholder)
	placeholder
placeholder

	if b.Len() == 0 {
		return ""
placeholder
	return contentSessionSeedPrefix + b.String()
placeholder

func TestDeriveOpenAIContentSessionSeed_ResponsesAPI_InputTextTypedItem(t *testing.T) {
	body := []byte(`{
		"model": "gpt-5.4",
		"input": [{"type": "input_text", "text": "Hello world"placeholder]
placeholder`)
	seed := deriveOpenAIContentSessionSeed(body)
	require.Contains(t, seed, "|first_user=")
	require.Contains(t, seed, "Hello world")
placeholder

func TestDeriveOpenAIContentSessionSeed_ResponsesAPI_TypedMessageItem(t *testing.T) {
	body := []byte(`{
		"model": "gpt-5.4",
		"input": [{"type": "message", "role": "user", "content": "Hello from typed message"placeholder]
placeholder`)
	seed := deriveOpenAIContentSessionSeed(body)
	require.Contains(t, seed, "|first_user=")
	require.Contains(t, seed, "Hello from typed message")
placeholder

func TestDeriveOpenAIStablePrefixSessionSeed_IgnoresUserContent(t *testing.T) {
	first := []byte(`{
		"model": "grok",
		"instructions": "Be concise.",
		"tools": [{"type":"function","name":"lookup","parameters":{"type":"object"placeholderplaceholder],
		"input": [{"role":"user","content":"Question A"placeholder]
placeholder`)
	second := []byte(`{
		"model": "grok",
		"instructions": "Be concise.",
		"tools": [{"parameters":{"type":"object"placeholder,"name":"lookup","type":"function"placeholder],
		"input": [{"role":"user","content":"Question B"placeholder]
placeholder`)

	firstSeed := deriveOpenAIStablePrefixSessionSeed(first)
	secondSeed := deriveOpenAIStablePrefixSessionSeed(second)

	require.NotEmpty(t, firstSeed)
	require.Equal(t, firstSeed, secondSeed)
	require.NotContains(t, firstSeed, "Question A")
	require.NotContains(t, firstSeed, "first_user")
placeholder

func TestDeriveOpenAIStablePrefixSessionSeed_IsolatesStablePrefixFields(t *testing.T) {
	base := []byte(`{
		"instructions":"Be concise.",
		"tools":[{"type":"function","name":"lookup"placeholder],
		"input":[{"role":"system","content":"System A"placeholder,{"role":"user","content":"Question"placeholder]
placeholder`)
	differentInstructions := []byte(`{
		"instructions":"Be detailed.",
		"tools":[{"type":"function","name":"lookup"placeholder],
		"input":[{"role":"system","content":"System A"placeholder,{"role":"user","content":"Question"placeholder]
placeholder`)
	differentTools := []byte(`{
		"instructions":"Be concise.",
		"tools":[{"type":"function","name":"search"placeholder],
		"input":[{"role":"system","content":"System A"placeholder,{"role":"user","content":"Question"placeholder]
placeholder`)
	differentSystem := []byte(`{
		"instructions":"Be concise.",
		"tools":[{"type":"function","name":"lookup"placeholder],
		"input":[{"role":"system","content":"System B"placeholder,{"role":"user","content":"Question"placeholder]
placeholder`)

	baseSeed := deriveOpenAIStablePrefixSessionSeed(base)
	require.NotEqual(t, baseSeed, deriveOpenAIStablePrefixSessionSeed(differentInstructions))
	require.NotEqual(t, baseSeed, deriveOpenAIStablePrefixSessionSeed(differentTools))
	require.NotEqual(t, baseSeed, deriveOpenAIStablePrefixSessionSeed(differentSystem))
placeholder

func TestDeriveOpenAIStablePrefixSessionSeed_ChatSystemAndDeveloper(t *testing.T) {
	first := []byte(`{
		"messages":[
			{"role":"system","content":"System prompt"placeholder,
			{"role":"developer","content":[{"type":"text","text":"Developer prompt"placeholder]placeholder,
			{"role":"user","content":"Question A"placeholder
		]
placeholder`)
	second := []byte(`{
		"messages":[
			{"role":"system","content":"System prompt"placeholder,
			{"role":"developer","content":[{"text":"Developer prompt","type":"text"placeholder]placeholder,
			{"role":"user","content":"Question B"placeholder
		]
placeholder`)

	firstSeed := deriveOpenAIStablePrefixSessionSeed(first)
	require.Equal(t, firstSeed, deriveOpenAIStablePrefixSessionSeed(second))
	require.Contains(t, firstSeed, "System prompt")
	require.Contains(t, firstSeed, "Developer prompt")
placeholder

func TestDeriveOpenAIStablePrefixSessionSeed_EncodesSystemAndDeveloperRoles(t *testing.T) {
	systemThenDeveloper := []byte(`{
		"messages":[
			{"role":"system","content":"Prompt A"placeholder,
			{"role":"developer","content":"Prompt B"placeholder
		]
placeholder`)
	developerThenSystem := []byte(`{
		"messages":[
			{"role":"developer","content":"Prompt A"placeholder,
			{"role":"system","content":"Prompt B"placeholder
		]
placeholder`)

	firstSeed := deriveOpenAIStablePrefixSessionSeed(systemThenDeveloper)
	secondSeed := deriveOpenAIStablePrefixSessionSeed(developerThenSystem)

	require.NotEqual(t, firstSeed, secondSeed)
	require.Contains(t, firstSeed, "|system=")
	require.Contains(t, firstSeed, "|developer=")
placeholder

func TestDeriveOpenAIStablePrefixSessionSeed_EncodesInstructionDelimiters(t *testing.T) {
	instructionOnly := []byte(`{
		"instructions":"foo|system=\"bar\""
placeholder`)
	instructionAndSystem := []byte(`{
		"instructions":"foo",
		"input":[{"role":"system","content":"bar"placeholder]
placeholder`)

	firstSeed := deriveOpenAIStablePrefixSessionSeed(instructionOnly)
	secondSeed := deriveOpenAIStablePrefixSessionSeed(instructionAndSystem)

	require.NotEmpty(t, firstSeed)
	require.NotEmpty(t, secondSeed)
	require.NotEqual(t, firstSeed, secondSeed)
placeholder

func TestDeriveOpenAIAnchoredContentSessionSeed_RequiresMeaningfulAnchor(t *testing.T) {
	emptyAnchors := [][]byte{
		nil,
		[]byte(`{"model":"grok"placeholder`),
		[]byte(`{"model":"grok","messages":[{"role":"assistant","content":"answer"placeholder]placeholder`),
		[]byte(`{"model":"grok","messages":[{"role":"user","content":"  "placeholder]placeholder`),
		[]byte(`{"model":"grok","messages":[{"role":"user","content":[{"type":"text","text":""placeholder]placeholder]placeholder`),
		[]byte(`{"model":"grok","input":"  "placeholder`),
		[]byte(`{"model":"grok","input":[{"type":"input_text","text":""placeholder]placeholder`),
placeholder
	for _, body := range emptyAnchors {
		require.Empty(t, deriveOpenAIAnchoredContentSessionSeed(body))
placeholder

	meaningfulAnchors := [][]byte{
		[]byte(`{"model":"grok","messages":[{"role":"user","content":"question"placeholder]placeholder`),
		[]byte(`{"model":"grok","messages":[{"role":"user","content":[{"type":"text","text":"question"placeholder]placeholder]placeholder`),
		[]byte(`{"model":"grok","input":"question"placeholder`),
		[]byte(`{"model":"grok","input":[{"type":"input_text","text":"question"placeholder]placeholder`),
placeholder
	for _, body := range meaningfulAnchors {
		require.NotEmpty(t, deriveOpenAIAnchoredContentSessionSeed(body))
placeholder
placeholder

func TestDeriveOpenAIStablePrefixSessionSeed_RequiresMeaningfulPrefix(t *testing.T) {
	tests := [][]byte{
		nil,
		[]byte(`{placeholder`),
		[]byte(`{"model":"grok","input":"Question A"placeholder`),
		[]byte(`{"model":"grok","tools":[],"input":"Question A"placeholder`),
		[]byte(`{"model":"grok","functions":[],"instructions":"  ","messages":[{"role":"system","content":""placeholder,{"role":"user","content":"Question A"placeholder]placeholder`),
placeholder

	for _, body := range tests {
		require.Empty(t, deriveOpenAIStablePrefixSessionSeed(body))
placeholder
placeholder
